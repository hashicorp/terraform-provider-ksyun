package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunSlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSlbRuleCreate,
		Read:   resourceKsyunSlbRuleRead,
		Update: resourceKsyunSlbRuleUpdate,
		Delete: resourceKsyunSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_header_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_server_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_sync": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"url_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"session": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_persistence_period": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"session_state": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func resourceKsyunSlbRuleCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"path",
		"host_header_id",
		"backend_server_group_id",
		"listener_sync",
		"method",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	createStructs := []string{
		"session",
		"health_check",
	}
	for _, v := range createStructs {
		if v1, ok := d.GetOk(v); ok {
			FlatternStruct(v1, &req)
		}
	}

	action := "CreateSlbRule"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.CreateSlbRule(&req)
	if err != nil {
		return fmt.Errorf("Error CreateSlbRules : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	item, ok := (*resp)["Rule"]
	if !ok {
		return fmt.Errorf("Error CreateSlbRules : no SlbRule found")
	}
	data, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error CreateSlbRules : no SlbRule found")
	}
	id, ok := data["RuleId"]
	if !ok {
		return fmt.Errorf("Error CreateSlbRules : no SlbRuleId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("Error CreateSlbRules : no SlbRuleId found")
	}
	d.SetId(idres)
	if err := d.Set("rule_id", idres); err != nil {
		return err
	}
	return resourceKsyunSlbRuleRead(d, m)
}

func resourceKsyunSlbRuleRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RuleId.1"] = d.Id()
	action := "DescribeRules"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeRules(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeSlbRules : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	itemset := (*resp)["RuleSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	excludes := SetDByResp(d, items[0], listenerKeys, map[string]bool{
		"HealthCheck": true,
		"Session":     true},
	)

	subSession := GetSubStructDByRep(excludes["Session"], map[string]bool{})
	if err := d.Set("session", []interface{}{subSession}); err != nil {
		return err
	}
	subHealth := GetSubStructDByRep(excludes["HealthCheck"], map[string]bool{})
	if err := d.Set("health_check", []interface{}{subHealth}); err != nil {
		return err
	}
	return nil
}

func resourceKsyunSlbRuleUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RuleId"] = d.Id()
	allAttributes := []string{
		"backend_server_group_id",
		"listener_sync",
		"method",
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
	if d.HasChange("session") && !d.IsNewResource() {
		FlatternStruct(d.Get("session"), &req)
		attributeUpdate = true
	}
	if d.HasChange("health_check") && !d.IsNewResource() {
		FlatternStruct(d.Get("health_check"), &req)
		attributeUpdate = true
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
	action := "ModifySlbRule"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.ModifySlbRule(&req)
	if err != nil {
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if strings.Contains(err.Error(), "400") {
			time.Sleep(time.Second * 3)
			resp, err = slbconn.ModifySlbRule(&req)
			if err != nil {
				return fmt.Errorf("update SlbRule (%v)error:%v", req, err)
			}
		}
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunSlbRuleRead(d, m)
}

func resourceKsyunSlbRuleDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["RuleId"] = d.Id()

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteRule"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteRule(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["RuleId.1"] = d.Id()
		action = "DescribeRules"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeRules(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading SlbRule when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["RuleSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified SlbRule %q has not been deleted due to unknown error", d.Id()))
	})
}
