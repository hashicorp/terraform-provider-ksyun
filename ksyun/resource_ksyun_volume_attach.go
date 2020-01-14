package ksyun

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-ksyun/logger"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKsyunVolumeAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeAttachCreate,
		Read:   resourceKsyunVolumeAttachRead,
		Delete: resourceKsyunVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKsyunVolumeAttachCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	var resp *map[string]interface{}
	attachReq := make(map[string]interface{})
	attach := []string{
		"volume_id",
		"instance_id",
	}
	for _, v := range attach {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			attachReq[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "AttachVolume"
	logger.Debug(logger.ReqFormat, action, attachReq)
	resp, err := conn.AttachVolume(&attachReq)
	if err != nil {
		return fmt.Errorf("error on attaching volume: %s", err)
	}
	logger.Debug(logger.RespFormat, action, attachReq, *resp)
	status, ok := (*resp)["Return"]
	if !ok {
		return fmt.Errorf("error on attaching volume")
	}
	status1, ok := status.(bool)
	if !ok || !status1 {
		return fmt.Errorf("error on attaching volume")
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("volume_id"), d.Get("instance_id")))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"in-use"},
		Refresh:    resourceKsyunVolumeStatusRefresh(conn, d.Get("volume_id").(string), []string{"in-use"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()
	_ = resourceKsyunVolumeAttachRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on waiting for volume %q complete attach, %s", d.Get("volume_id").(string), err)
	}
	return nil
}

func resourceKsyunVolumeAttachRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	readReq := make(map[string]interface{})
	readReq["VolumeId.1"] = strings.Split(d.Id(), ":")[0]
	action := "DescribeVolumes"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeVolumes(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	volumeList, ok := (*resp)["Volumes"]
	if !ok {
		d.SetId("")
		return nil
	}
	volumes, ok1 := volumeList.([]interface{})
	if !ok1 {
		d.SetId("")
		return nil
	}
	if volumes == nil || len(volumes) < 1 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, volumes[0], volumeKeys, map[string]bool{})
	return nil
}

func resourceKsyunVolumeAttachDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	detachReq := make(map[string]interface{})
	detachReq["VolumeId"] = strings.Split(d.Id(), ":")[0]
	action := "DetachVolume"
	logger.Debug(logger.ReqFormat, action, detachReq)
	resp, err := conn.DetachVolume(&detachReq)
	if err != nil {
		return fmt.Errorf("error on detach volume %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, detachReq, *resp)
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		readReq := make(map[string]interface{})
		readReq["VolumeId.1"] = strings.Split(d.Id(), ":")[0]
		action := "DescribeVolumes"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err := conn.DescribeVolumes(&readReq)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		logger.Debug(logger.RespFormat, action, readReq, *resp)
		volumeList, ok := (*resp)["Volumes"]
		if !ok {
			return nil
		}
		volumes, ok1 := volumeList.([]interface{})
		if !ok1 {
			return nil
		}
		if volumes == nil || len(volumes) < 1 {
			return nil
		}
		volume, ok2 := volumes[0].(map[string]interface{})
		if !ok2 {
			return nil
		}
		status, ok3 := volume["VolumeStatus"]
		if !ok3 {
			return nil
		}
		if status == "available" {
			return nil
		}
		return resource.RetryableError(errors.New("detaching"))
	})
}
