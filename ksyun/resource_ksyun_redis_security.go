package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// instance security rule
func resourceRedisSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityRuleCreate,
		Delete: resourceRedisSecurityRuleDelete,
		Update: resourceRedisSecurityRuleUpdate,
		Read:   resourceRedisSecurityRuleRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceRedisSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		rules []string
		resp  *map[string]interface{}
		err   error
		//ok    bool
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	createReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		createReq["AvailableZone"] = az
	}
	rules = SchemaSetToStringSlice(d.Get("rules"))
	/*if rules, ok = d.Get("rules").([]interface{}); !ok {
		return fmt.Errorf("type of security_rule.rules must be array string")
	}*/
	for i, rule := range rules {
		createReq[fmt.Sprintf("%v%v", "SecurityRules.Cidr.", i+1)] = rule
	}
	action := "SetCacheSecurityRules"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.SetCacheSecurityRules(&createReq); err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	d.SetId(createReq["CacheId"].(string))
	resourceRedisSecurityRuleRead(d, meta)
	return nil
}

func resourceRedisSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	// query security rule
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	action := "DescribeCacheSecurityRules"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeCacheSecurityRules(&readReq); err != nil {
		return fmt.Errorf("error on reading instance security rule %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		logger.Info("instance security rule result size : 0")
		return nil
	}

	var rules []interface{}
	for _, v := range data {
		group := v.(map[string]interface{})
		rules = append(rules, group["securityRuleId"])
	}

	// delete security rules
	deleteReq := make(map[string]interface{})
	deleteReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		deleteReq["AvailableZone"] = az
	}
	action = "DeleteCacheSecurityRule"
	for _, rule := range rules {
		deleteReq["SecurityRuleId"] = fmt.Sprintf("%v", rule)
		logger.Debug(logger.ReqFormat, action, deleteReq)
		if resp, err = conn.DeleteCacheSecurityRule(&deleteReq); err != nil {
			return fmt.Errorf("error on delete instance security rule: %s", err)
		}
		logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	}
	return nil
}

func resourceRedisSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		deleteRulesResults []string
		addRulesResults    []string
		resp               *map[string]interface{}
		err                error
	)
	d.Partial(true)
	defer d.Partial(false)
	updateReq := make(map[string]interface{})
	updateReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		updateReq["AvailableZone"] = az
	}
	if d.HasChange("rules") {
		oldMc, newMc := d.GetChange("rules")
		oldRules := SchemaSetToStringSlice(oldMc)
		newRules := SchemaSetToStringSlice(newMc)
		//oldRules := oldMc.([]interface{})
		//newRules := newMc.([]interface{})
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
	}
	if len(addRulesResults) > 0 {
		for i, rule := range addRulesResults {
			updateReq[fmt.Sprintf("%v%v", "SecurityRules.Cidr.", i+1)] = rule
		}
		conn := meta.(*KsyunClient).kcsv1conn
		action := "SetCacheSecurityRules"
		logger.Debug(logger.ReqFormat, action, updateReq)
		if resp, err = conn.SetCacheSecurityRules(&updateReq); err != nil {
			return fmt.Errorf("error on add instance security rule: %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq, *resp)
	}
	if len(deleteRulesResults) > 0 {
		conn := meta.(*KsyunClient).kcsv1conn
		// query security rule
		readReq := make(map[string]interface{})
		readReq["CacheId"] = d.Get("cache_id")
		if az, ok := d.GetOk("available_zone"); ok {
			readReq["AvailableZone"] = az
		}
		action := "DescribeCacheSecurityRules"
		logger.Debug(logger.ReqFormat, action, readReq)
		if resp, err = conn.DescribeCacheSecurityRules(&readReq); err != nil {
			return fmt.Errorf("error on reading instance security rule %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, readReq, *resp)
		data := (*resp)["Data"].([]interface{})
		if len(data) > 0 {
			var rules []interface{}
		X:
			for _, v := range data {
				group := v.(map[string]interface{})
				for _, r := range deleteRulesResults {
					if group["cidr"].(string) == r {
						rules = append(rules, group["securityRuleId"])
						continue X
					}
				}
			}

			deleteRuleReq := make(map[string]interface{})
			deleteRuleReq["CacheId"] = d.Get("cache_id")
			for _, delRuleId := range rules {
				conn := meta.(*KsyunClient).kcsv1conn
				deleteRuleReq["SecurityRuleId"] = fmt.Sprintf("%v", delRuleId)
				if az, ok := d.GetOk("available_zone"); ok {
					deleteRuleReq["AvailableZone"] = az
				}
				action := "DeleteCacheSecurityRule"
				logger.Debug(logger.ReqFormat, action, deleteRuleReq)
				if resp, err = conn.DeleteCacheSecurityRule(&deleteRuleReq); err != nil {
					return fmt.Errorf("error on delete instance security rule: %s", err)
				}
				logger.Debug(logger.RespFormat, action, deleteRuleReq, *resp)
			}
		}
	}
	resourceRedisSecurityRuleRead(d, meta)
	return nil
}

func resourceRedisSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	action := "DescribeCacheSecurityRules"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeCacheSecurityRules(&readReq); err != nil {
		return fmt.Errorf("error on reading instance security rule %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		logger.Info("instance security rule result size : 0")
		return nil
	}
	result := make(map[string]interface{})
	var rulesTemp []string
	for _, v := range data {
		group := v.(map[string]interface{})
		rulesTemp = append(rulesTemp, group["cidr"].(string))
	}
	result["rules"] = rulesTemp
	for k, v := range result {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("error set data %v :%v", v, err)
		}
	}
	return nil
}
