package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

var getInTheSgCar = map[string]bool{
	"security_group_id":          true,
	"security_group_name":        true,
	"security_group_description": true,
	"created":                    true,
	//"instances": true,
	"security_group_rule": true,
}

func resourceKsyunKrdsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKrdsSecurityGroupCreate,
		Update: resourceKsyunKrdsSecurityGroupUpdate,
		Read:   resourceKsyunKrdsSecurityGroupRead,
		Delete: resourceKsyunKrdsSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_description": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"created": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_rule": {
				Type:     schema.TypeSet,
				Set:      secParameterToHash,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_rule_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"security_group_rule_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group_rule_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			//"instances": {
			//	Type:     schema.TypeSet,
			//	Optional: true,
			//	Computed: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"db_instance_identifier": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Computed: true,
			//			},
			//			"db_instance_name": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Computed: true,
			//			},
			//			"vip": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Computed: true,
			//			},
			//			"created": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Computed: true,
			//			},
			//			"db_instance_type": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Computed: true,
			//			},
			//		},
			//	},
			//},
		},
	}
}

func secParameterToHash(ruleMap interface{}) int {
	rule := ruleMap.(map[string]interface{})
	return hashcode.String(rule["security_group_rule_protocol"].(string) + "|" + rule["security_group_rule_name"].(string))
}

func resourceKsyunKrdsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var resp *map[string]interface{}
	var err error
	createReq := make(map[string]interface{})
	creates := []string{
		"SecurityGroupName",
		"SecurityGroupDescription",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(Camel2Hungarian(v)); ok {
			createReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	rules := d.Get("security_group_rule").(*schema.Set).List()
	for ruleindex, rule := range rules {
		num := ruleindex + 1
		createReq["SecurityGroupRule.SecurityGroupRuleProtocol."+strconv.Itoa(num)] = rule.(map[string]interface{})["security_group_rule_protocol"].(string)
		createReq["SecurityGroupRule.SecurityGroupRuleName."+strconv.Itoa(num)] = rule.(map[string]interface{})["security_group_rule_name"].(string)
	}

	action := "CreateSecurityGroup"
	logger.Debug(logger.RespFormat, action, createReq)
	resp, err = conn.CreateSecurityGroup(&createReq)
	logger.Debug(logger.AllFormat, action, createReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating instance security group(krds): %s", err)
	}

	if resp != nil {
		bodyData := (*resp)["Data"].(map[string]interface{})
		securityGroups := bodyData["SecurityGroups"].([]interface{})
		instanceId := securityGroups[0].(map[string]interface{})["SecurityGroupId"].(string)
		logger.DebugInfo("~*~*~*~*~ SecurityGroupId : %v", instanceId)
		d.SetId(instanceId)
	}

	return resourceKsyunKrdsSecurityGroupRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn

	logger.DebugInfo("update security group , IsNewReource : %v", d.IsNewResource())
	d.Partial(true)
	if d.HasChange("security_group_name") || d.HasChange("security_group_description") {
		req := map[string]interface{}{
			"SecurityGroupId":          d.Id(),
			"SecurityGroupName":        d.Get("security_group_name"),
			"SecurityGroupDescription": d.Get("security_group_description"),
			"SecurityGroupRules":       d.Get("security_group_rule").(*schema.Set).List(),
			"SecurityGroupRuleAction":  "Attach",
		}
		action := "ModifySecurityGroup"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ModifySecurityGroup(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return fmt.Errorf("error on update instance security group(krds): %s", err)
		}
		d.SetPartial("security_group_name")
		d.SetPartial("security_group_description")
	}
	if d.HasChange("security_group_rule") {
		old, new := d.GetChange("security_group_rule")

		addRules := []interface{}{}
		delRules := []interface{}{}
		oldRules := old.(*schema.Set).List()
		newRules := new.(*schema.Set).List()
		for _, oldRule := range oldRules {
			isExist := false
			for _, newRule := range newRules {
				if secParameterToHash(oldRule) == secParameterToHash(newRule) {
					isExist = true
				}
			}
			if !isExist {
				delRules = append(delRules, oldRule)
			}
		}
		for _, newRule := range newRules {
			isExist := false
			for _, oldRule := range oldRules {
				if secParameterToHash(oldRule) == secParameterToHash(newRule) {
					isExist = true
				}
			}
			if !isExist {
				addRules = append(addRules, newRule)
			}
		}

		addReq := map[string]interface{}{
			"SecurityGroupId":         d.Id(),
			"SecurityGroupRuleAction": "Attach",
		}
		for index, rule := range addRules {
			num := index + 1
			addReq["SecurityGroupRule.SecurityGroupRuleProtocol."+strconv.Itoa(num)] = rule.(map[string]interface{})["security_group_rule_protocol"].(string)
			addReq["SecurityGroupRule.SecurityGroupRuleName."+strconv.Itoa(num)] = rule.(map[string]interface{})["security_group_rule_name"].(string)
		}
		action := "ModifySecurityGroupRule"
		logger.Debug(logger.ReqFormat, action, addReq)
		addResp, err := conn.ModifySecurityGroupRule(&addReq)
		logger.Debug(logger.AllFormat, action, addReq, *addResp, err)

		rulesInfoMap := map[int]map[string]interface{}{}
		// terraform.tfstate读取的不对，所以用output_file里的数据
		rulesInfo := d.Get("security_group_rule").(*schema.Set).List()
		if err != nil {
			return err
		}
		for _, value := range rulesInfo {
			rulesInfoMap[secParameterToHash(value)] = value.(map[string]interface{})
		}
		logger.DebugInfo(" rulesInfoMap : %v", rulesInfoMap)

		delReq := map[string]interface{}{
			"SecurityGroupId":         d.Id(),
			"SecurityGroupRules":      delRules,
			"SecurityGroupRuleAction": "Delete",
		}
		for index, rule := range delRules {
			num := index + 1
			delReq["SecurityGroupRule.SecurityGroupRuleId."+strconv.Itoa(num)] = rulesInfoMap[secParameterToHash(rule)]["security_group_rule_id"]
		}
		action = "ModifySecurityGroupRule"
		logger.Debug(logger.ReqFormat, action, delReq)
		delResp, err := conn.ModifySecurityGroupRule(&delReq)
		logger.Debug(logger.AllFormat, action, delReq, *delResp, err)
		d.SetPartial("security_group_rule")
	}
	err := resourceKsyunKrdsSecurityGroupRead(d, meta)
	d.Partial(false)

	return err
}

func resourceKsyunKrdsSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	req := map[string]interface{}{"SecurityGroupId": d.Id()}
	action := "DescribeSecurityGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeSecurityGroup(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)

	if err != nil {
		return fmt.Errorf("error on reading instance security group id : %q, %s", d.Id(), err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading instance security group id %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["SecurityGroups"].([]interface{})

	krdsIds := make([]string, len(instances))
	krdsMapList := make([]map[string]interface{}, len(instances))
	for num, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		krdsMap := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "Instances" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[Camel2Hungarian(j)] = q
						}
						wtf[num] = rrmap
					}
					krdsMap["instances"] = wtf
				}
			} else if k == "SecurityGroupRules" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[Camel2Hungarian(j)] = q
						}
						wtf[num] = rrmap
					}
					krdsMap["security_group_rule"] = wtf
				}
			} else {
				krdsMap[Camel2Hungarian(k)] = v
			}
		}
		logger.DebugInfo(" converted ---- %+v ", krdsMap)

		krdsIds[num] = krdsMap["security_group_id"].(string)
		logger.DebugInfo("krdsIds fuck : %v", krdsIds)
		krdsMapList[num] = krdsMap
	}
	logger.DebugInfo(" converted ---- %+v ", krdsMapList)
	if len(krdsMapList) == 0 {
		return fmt.Errorf("no data on reading security group (%q) ", d.Id())
	}
	_ = SetDByFkResp(d, krdsMapList[0], getInTheSgCar)

	return nil
}

func resourceKsyunKrdsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	req := map[string]interface{}{"SecurityGroupId": d.Id()}
	action := "DeleteSecurityGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DeleteSecurityGroup(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)

	if err != nil {
		return fmt.Errorf("error on reading instance security group id : %q, %s", d.Id(), err)
	}
	return nil
}
