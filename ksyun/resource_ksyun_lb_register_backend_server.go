package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunRegisterBackendServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRegisterBackendServerCreate,
		Read:   resourceKsyunRegisterBackendServerRead,
		Update: resourceKsyunRegisterBackendServerUpdate,
		Delete: resourceKsyunRegisterBackendServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_server_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_server_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"register_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_server_ip": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"real_server_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_server_state": {
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
func resourceKsyunRegisterBackendServerCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"backend_server_group_id",
		"backend_server_ip",
		"backend_server_port",
		"weight",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "RegisterBackendServer"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.RegisterBackendServer(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("Error RegisterBackendServer : %s", err)
	}
	item, ok := (*resp)["BackendServer"]
	if !ok {
		return fmt.Errorf("RegisterBackendServer Error  : no RegisterBackendServer found")
	}
	data, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("RegisterBackendServer Error  : no RegisterBackendServer found")
	}
	id, ok := data["RegisterId"]
	if !ok {
		return fmt.Errorf("RegisterBackendServer Error  : no id found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("RegisterBackendServer Error : no id found")
	}
	d.SetId(idres)
	return resourceKsyunRegisterBackendServerRead(d, m)
}

func resourceKsyunRegisterBackendServerRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RegisterId.1"] = d.Id()
	action := "DescribeBackendServers"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeBackendServers(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeBackendServers : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	itemset := (*resp)["BackendServerSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], registerBackendServerKeys, map[string]bool{})
	return nil
}

func resourceKsyunRegisterBackendServerUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RegisterId"] = d.Id()
	allAttributes := []string{
		"weight",
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
	action := "ModifyBackendServer"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.ModifyBackendServer(&req)
	if err != nil {
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if strings.Contains(err.Error(), "400") {
			time.Sleep(time.Second * 3)
			resp, err = slbconn.ModifyBackendServer(&req)
			if err != nil {
				return fmt.Errorf("update RegisterBackendServer (%v)error:%v", req, err)
			}
		}
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunRegisterBackendServerRead(d, m)
}

func resourceKsyunRegisterBackendServerDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RegisterId"] = d.Id()
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeregisterBackendServer"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeregisterBackendServer(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["RegisterId.1"] = d.Id()
		action = "DescribeBackendServers"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeBackendServers(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading RegisterBackendServer when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["BackendServerSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified RegisterBackendServer %q has not been deleted due to unknown error", d.Id()))
	})
}
