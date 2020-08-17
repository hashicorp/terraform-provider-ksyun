package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunBackendServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBackendServerGroupCreate,
		Read:   resourceKsyunBackendServerGroupRead,
		Update: resourceKsyunBackendServerGroupUpdate,
		Delete: resourceKsyunBackendServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_server_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_server_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backend_server_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunBackendServerGroupCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"vpc_id",
		"backend_server_group_name",
		"backend_server_group_type",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateBackendServerGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.CreateBackendServerGroup(&req)
	if err != nil {
		return fmt.Errorf("Error CreateBackendServerGroup : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	item, ok := (*resp)["BackendServerGroup"]
	if !ok {
		return fmt.Errorf("Error CreateBackendServerGroup : no BackendServerGroup found")
	}
	data, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error CreateBackendServerGroup : no BackendServerGroup found")
	}
	id, ok := data["BackendServerGroupId"]
	if !ok {
		return fmt.Errorf("Error CreateBackendServerGroup : no BackendServerGroupId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("Error CreateBackendServerGroup : no BackendServerGroupId found")
	}
	d.SetId(idres)
	return resourceKsyunBackendServerGroupRead(d, m)
}

func resourceKsyunBackendServerGroupRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["BackendServerGroupId.1"] = d.Id()
	action := "DescribeBackendServerGroups"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeBackendServerGroups(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeBackendServerGroups : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	itemset := (*resp)["BackendServerGroupSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], backendServerGroupKeys, map[string]bool{})
	return nil
}

func resourceKsyunBackendServerGroupUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["BackendServerGroupId"] = d.Id()
	allAttributes := []string{
		"backend_server_group_name",
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
	action := "ModifyBackendServerGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.ModifyBackendServerGroup(&req)
	if err != nil {
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if strings.Contains(err.Error(), "400") {
			time.Sleep(time.Second * 3)
			resp, err = slbconn.ModifyBackendServerGroup(&req)
			if err != nil {
				return fmt.Errorf("update BackendServerGroup (%v)error:%v", req, err)
			}
		}
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunBackendServerGroupRead(d, m)
}

func resourceKsyunBackendServerGroupDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["BackendServerGroupId"] = d.Id()
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteBackendServerGroup"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteBackendServerGroup(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["BackendServerGroupId.1"] = d.Id()
		action = "DescribeBackendServerGroups"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeBackendServerGroups(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on DescribeBackendServerGroups when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["BackendServerGroupSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified BackendServerGroup %q has not been deleted due to unknown error", d.Id()))
	})
}
