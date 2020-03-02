package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunInstancesWithListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstancesWithListenerCreate,
		Read:   resourceKsyunInstancesWithListenerRead,
		Update: resourceKsyunInstancesWithListenerUpdate,
		Delete: resourceKsyunInstancesWithListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"real_server_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"real_server_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"real_server_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "1",
			},
			"real_server_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"register_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunInstancesWithListenerCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"listener_id",
		"real_server_ip",
		"real_server_port",
		"real_server_type",
		"instance_id",
		"weight",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "RegisterInstancesWithListener"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.RegisterInstancesWithListener(&req)
	if err != nil {
		return fmt.Errorf("create InstancesWithListener : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	id, ok := (*resp)["RegisterId"]
	if !ok {
		return fmt.Errorf("create InstancesWithListener : no HealthCheckId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("create InstancesWithListener : no HealthCheckId found")
	}
	if err:=d.Set("register_id", idres);err!=nil{
		return err
	}
	d.SetId(idres)
	return resourceKsyunInstancesWithListenerRead(d, m)
}

func resourceKsyunInstancesWithListenerRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RegisterId.1"] = d.Id()
	action := "DescribeInstancesWithListener"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := slbconn.DescribeInstancesWithListener(&req)
	if err != nil {
		return fmt.Errorf(" read InstancesWithListeners : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	resSet := (*resp)["RealServerSet"]
	res, ok := resSet.([]interface{})
	if !ok || len(res) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, res[0], serverKeys, map[string]bool{})
	return nil
}

func resourceKsyunInstancesWithListenerUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RegisterId"] = d.Id()
	allAttributes := []string{
		"real_server_port",
		"weight",
	}
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
	action := "ModifyInstancesWithListener"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := slbconn.ModifyInstancesWithListener(&req)
	if err != nil {
		return fmt.Errorf("update InstancesWithListener (%v)error:%v", req, err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	// Set partial modification properties
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunInstancesWithListenerRead(d, m)
}

func resourceKsyunInstancesWithListenerDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RegisterId"] = d.Id()
	/*
		_, err := slbconn.DeregisterInstancesFromListener(&req)
		if err != nil {
			return fmt.Errorf("delete InstancesWithListener error:%v", err)
		}
		return nil
	*/
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeregisterInstancesFromListener"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeregisterInstancesFromListener(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}

		req := make(map[string]interface{})
		req["RegisterId.1"] = d.Id()
		action = "DescribeInstancesWithListener"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeInstancesWithListener(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading Listener instance when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["RealServerSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified Listener instance %q has not been deleted due to unknown error", d.Id()))
	})
}
