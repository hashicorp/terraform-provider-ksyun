package ksyun

import (
	"errors"
	"fmt"
	"time"

	"github.com/KscSDK/ksc-sdk-go/service/ebs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeCreate,
		Update: resourceKsyunVolumeUpdate,
		Read:   resourceKsyunVolumeRead,
		Delete: resourceKsyunVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceKsyunVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	creates := []string{
		"volume_name",
		"volume_type",
		"volume_desc",
		"size",
		"charge_type",
		"availability_zone",
		"project_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createReq[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateVolume"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err := conn.CreateVolume(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating volume: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	id, ok := (*resp)["VolumeId"]
	if !ok {
		return fmt.Errorf("error on creating volume : no id found")
	}
	idRes, ok := id.(string)
	if !ok {
		return fmt.Errorf("error on creating volume : no id found")
	}
	d.SetId(idRes)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"available"},
		Refresh:    resourceKsyunVolumeStatusRefresh(conn, d.Id(), []string{"available"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()
	_ = resourceKsyunVolumeRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on waiting for volume %q complete creating, %s", d.Id(), err)
	}
	return nil
}

func resourceKsyunVolumeStatusRefresh(conn *ebs.Ebs, volumeId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"VolumeId.1": volumeId}
		action := "DescribeVolumes"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.DescribeVolumes(&req)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		volumeList, ok := (*resp)["Volumes"]
		if !ok {
			return nil, "", fmt.Errorf("no volume get")
		}
		volumes, ok1 := volumeList.([]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no volume get")
		}
		if volumes == nil || len(volumes) < 1 {
			return nil, "", fmt.Errorf("no volume get")
		}
		volume, ok2 := volumes[0].(map[string]interface{})
		if !ok2 {
			return nil, "", fmt.Errorf("no volume get")
		}
		status, ok3 := volume["VolumeStatus"]
		if !ok3 {
			return nil, "", fmt.Errorf("no volume status get")
		}
		if status == "error" {
			return nil, "", fmt.Errorf("volume error")
		}
		for k, v := range target {
			if v == status {
				return resp, status.(string), nil
			}
			if k == len(target)-1 {
				status = statusPending
			}
		}
		return resp, status.(string), nil
	}
}

func resourceKsyunVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	readReq := make(map[string]interface{})
	readReq["VolumeId.1"] = d.Id()
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

func resourceKsyunVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	conn := meta.(*KsyunClient).ebsconn
	if d.HasChange("volume_name") && !d.IsNewResource() {
		d.SetPartial("volume_name")
		update := make(map[string]interface{})
		update["VolumeId"] = d.Id()
		if v, ok := d.GetOk("volume_name"); ok {
			update["VolumeName"] = v.(string)
		} else {
			return fmt.Errorf("cann't change volume_name to empty string")
		}
		action := "ModifyVolume"
		logger.Debug(logger.ReqFormat, action, update)
		resp, err := conn.ModifyVolume(&update)
		if err != nil {
			return fmt.Errorf("error on update volume name %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, update, *resp)
	}
	if d.HasChange("volume_desc") && !d.IsNewResource() {
		d.SetPartial("volume_desc")
		update := make(map[string]interface{})
		update["VolumeId"] = d.Id()
		if v, ok := d.GetOk("volume_desc"); ok {
			update["VolumeDesc"] = v.(string)
		} else {
			return fmt.Errorf("cann't change volume_desc to empty string")
		}
		action := "ModifyVolume"
		logger.Debug(logger.ReqFormat, action, update)
		resp, err := conn.ModifyVolume(&update)
		if err != nil {
			return fmt.Errorf("error on update volume desc %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, update, *resp)
	}
	if d.HasChange("size") && !d.IsNewResource() {
		d.SetPartial("size")
		update := make(map[string]interface{})
		update["VolumeId"] = d.Id()
		if v, ok := d.GetOk("size"); ok {
			update["Size"] = fmt.Sprintf("%v", v.(int))
		} else {
			return fmt.Errorf("cann't change size to empty")
		}
		action := "ResizeVolume"
		logger.Debug(logger.ReqFormat, action, update)
		resp, err := conn.ResizeVolume(&update)
		if err != nil {
			return fmt.Errorf("error on resize volume %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, update, *resp)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"available"},
			Refresh:    resourceKsyunVolumeStatusRefresh(conn, d.Id(), []string{"available"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForState()
		_ = resourceKsyunVolumeRead(d, meta)
		if err != nil {
			return fmt.Errorf("error on waiting for volume %q complete resize, %s", d.Id(), err)
		}
	}
	d.Partial(false)
	return resourceKsyunVolumeRead(d, meta)
}

func resourceKsyunVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	deleteReq := make(map[string]interface{})
	deleteReq["VolumeId"] = d.Id()
	action := "DeleteVolume"
	logger.Debug(logger.ReqFormat, action, deleteReq)
	resp, err := conn.DeleteVolume(&deleteReq)
	if err != nil {
		return fmt.Errorf("error on delete volume %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		readReq := make(map[string]interface{})
		readReq["VolumeId.1"] = d.Id()
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
		if status == "recycling" {
			return nil
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}
