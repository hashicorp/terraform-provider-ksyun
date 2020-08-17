package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunListenerHostHeader() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerHostHeaderCreate,
		Read:   resourceKsyunListenerHostHeaderRead,
		Update: resourceKsyunListenerHostHeaderUpdate,
		Delete: resourceKsyunListenerHostHeaderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_header": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_header_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunListenerHostHeaderCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"listener_id",
		"host_header",
		"certificate_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateHostHeader"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.CreateHostHeader(&req)
	if err != nil {
		return fmt.Errorf("Error CreateHostHeader : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	item, ok := (*resp)["HostHeader"]
	if !ok {
		return fmt.Errorf("Error CreateHostHeader : no HostHeader found")
	}
	data, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error CreateHostHeader : no HostHeader found")
	}
	id, ok := data["HostHeaderId"]
	if !ok {
		return fmt.Errorf("Error CreateHostHeader : no HostHeaderId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("Error CreateHostHeader : no HostHeaderId found")
	}
	d.SetId(idres)
	if err := d.Set("host_header_id", idres); err != nil {
		return err
	}
	return resourceKsyunListenerHostHeaderRead(d, m)
}

func resourceKsyunListenerHostHeaderRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["HostHeaderId.1"] = d.Id()
	action := "DescribeHostHeaders"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeHostHeaders(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeHostHeaders : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	itemset := (*resp)["HostHeaderSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], hostHeaderKeys, map[string]bool{})
	return nil
}

func resourceKsyunListenerHostHeaderUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["HostHeaderId"] = d.Id()
	allAttributes := []string{
		"certificate_id",
	}
	// Whether the representative has any modifications
	attributeUpdate := false
	var updates []string
	//Get the property that needs to be modified
	for _, v := range allAttributes {
		if d.HasChange(v) {
			attributeUpdate = true
			updates = append(updates, v)
		}
	}
	if !attributeUpdate {
		return nil
	}
	//Create a modification request
	for _, v := range allAttributes {
		if v1, ok := d.GetOk(v); ok {
			req[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
		}
	}
	// Enable partial attribute modification
	d.Partial(true)
	action := "ModifyHostHeader"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.ModifyHostHeader(&req)
	if err != nil {
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if strings.Contains(err.Error(), "400") {
			time.Sleep(time.Second * 3)
			resp, err = slbconn.ModifyHostHeader(&req)
			if err != nil {
				return fmt.Errorf("update HostHeader (%v)error:%v", req, err)
			}
		}
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunListenerHostHeaderRead(d, m)
}

func resourceKsyunListenerHostHeaderDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["HostHeaderId"] = d.Id()
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteHostHeader"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteHostHeader(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["HostHeaderId.1"] = d.Id()
		action = "DescribeHostHeaders"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeHostHeaders(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on DescribeHostHeaders when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["HostHeaderSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified HostHeader %q has not been deleted due to unknown error", d.Id()))
	})
}
