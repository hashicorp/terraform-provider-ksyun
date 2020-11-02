package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// redis security group rule
func resourceRedisSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupRuleCreate,
		Delete: resourceRedisSecurityGroupRuleDelete,
		Update: resourceRedisSecurityGroupRuleUpdate,
		Read:   resourceRedisSecurityGroupRuleRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceRedisSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		createReq["AvailableZone"] = az
	}
	createReq["SecurityGroupId"] = d.Get("security_group_id")
	rules := SchemaSetToStringSlice(d.Get("rules"))
	for i, rule := range rules {
		createReq[fmt.Sprintf("%v%v", "Cidrs.", i+1)] = rule
	}
	action := "CreateSecurityGroupRule"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.CreateSecurityGroupRule(&createReq); err != nil {
		return fmt.Errorf("error on create redis security group rule: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	d.SetId(d.Get("security_group_id").(string))
	return nil
}

func resourceRedisSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	// read redis security group rule
	rules, err := readSecGroupRule(d, meta)
	if err != nil {
		return nil
	}

	// delete redis security group rule
	deleteReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		deleteReq["AvailableZone"] = az
	}
	deleteReq["SecurityGroupId"] = d.Id()
	for k, v := range rules {
		deleteReq[fmt.Sprintf("%v%v", "SecurityGroupRuleId.", k+1)] = v["id"]
	}
	action := "DeleteSecurityGroupRule"
	logger.Debug(logger.ReqFormat, action, deleteReq)
	if resp, err = conn.DeleteSecurityGroupRule(&deleteReq); err != nil {
		return fmt.Errorf("error on delete redis security group rule: %s", err)
	}
	logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	return nil
}

func resourceRedisSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		deleteRulesResults []string
		addRulesResults    []string
		resp               *map[string]interface{}
		err                error
	)
	d.Partial(true)
	defer d.Partial(false)
	if !d.HasChange("rules") {
		return nil
	}

	oldMc, newMc := d.GetChange("rules")
	oldRules := SchemaSetToStringSlice(oldMc)
	newRules := SchemaSetToStringSlice(newMc)
	for _, oldRule := range oldRules {
		oldR := oldRule
		exist := false
		for _, newRule := range newRules {
			newR := newRule
			if newR == oldR {
				exist = true
			}
		}
		if !exist {
			deleteRulesResults = append(deleteRulesResults, oldR)
		}
	}

	for _, newRule := range newRules {
		ip := newRule
		exist := false
		for _, oldRule := range oldRules {
			oldR := oldRule
			if oldR == ip {
				exist = true
			}
		}
		if !exist {
			addRulesResults = append(addRulesResults, ip)
		}
	}

	conn := meta.(*KsyunClient).kcsv1conn
	if len(addRulesResults) > 0 {
		updateReq := make(map[string]interface{})
		updateReq["SecurityGroupId"] = d.Id()
		if az, ok := d.GetOk("available_zone"); ok {
			updateReq["AvailableZone"] = az
		}
		for i, rule := range addRulesResults {
			updateReq[fmt.Sprintf("%v%v", "Cidrs.", i+1)] = rule
		}

		action := "CreateSecurityGroupRule"
		logger.Debug(logger.ReqFormat, action, updateReq)
		if resp, err = conn.CreateSecurityGroupRule(&updateReq); err != nil {
			return fmt.Errorf("error on add redis security group rule: %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq, *resp)
	}

	if len(deleteRulesResults) > 0 {
		rules, e := readSecGroupRule(d, meta)
		if e != nil {
			return e
		}

		if rules == nil {
			return nil
		}

		deleteRuleReq := make(map[string]interface{})
		if az, ok := d.GetOk("available_zone"); ok {
			deleteRuleReq["AvailableZone"] = az
		}
		deleteRuleReq["SecurityGroupId"] = d.Id()
		k := 1
	X:
		for _, v := range rules {
			for _, r := range deleteRulesResults {
				if v["cidr"] == r {
					deleteRuleReq[fmt.Sprintf("%v%v", "SecurityGroupRuleId.", k)] = v["id"]
					k++
					continue X
				}
			}
		}

		action := "DeleteSecurityGroupRule"
		logger.Debug(logger.ReqFormat, action, deleteRuleReq)
		if resp, err = conn.DeleteSecurityGroupRule(&deleteRuleReq); err != nil {
			return fmt.Errorf("error on delete redis security group rule: %s", err)
		}
		logger.Debug(logger.RespFormat, action, deleteRuleReq, *resp)
	}

	return nil
}

func resourceRedisSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	rules, err := readSecGroupRule(d, meta)
	if err != nil {
		return nil
	}

	var result []string
	for _, v := range rules {
		result = append(result, v["cidr"])
	}

	if err := d.Set("rules", result); err != nil {
		return fmt.Errorf("error set data %v :%v", result, err)
	}

	return nil
}

func readSecGroupRule(d *schema.ResourceData, meta interface{}) ([]map[string]string, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	readReq["SecurityGroupId"] = d.Id()
	action := "DescribeSecurityGroup"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeSecurityGroup(&readReq); err != nil {
		return nil, fmt.Errorf("error on reading redis security group %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].(map[string]interface{})
	if len(data) == 0 {
		logger.Info("redis security group result size : 0")
		return nil, nil
	}

	rules := data["rules"].([]interface{})
	if len(rules) == 0 {
		logger.Info("redis security group rule result size : 0")
		return nil, nil
	}

	var result []map[string]string
	for _, v := range rules {
		group := v.(map[string]interface{})
		rule := map[string]string{}
		rule["id"] = group["id"].(string)
		rule["cidr"] = group["cidr"].(string)
		result = append(result, rule)
	}

	return result, nil
}
