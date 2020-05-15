package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
)

func resourceKsyunMongodbSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbSecurityRuleCreate,
		Delete: resourceMongodbSecurityRuleDelete,
		Update: resourceMongodbSecurityRuleUpdate,
		Read:   resourceMongodbSecurityRuleRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidrs": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"to_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongodbSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {

	cidrs := d.Get("cidrs")
	if cidrs.(string) == "" {
		return fmt.Errorf("error on set instance security rule: cidrs is empty")
	}
	createReq := make(map[string]interface{})
	createReq["InstanceId"] = d.Get("instance_id")
	createReq["cidrs"] = cidrs

	conn := meta.(*KsyunClient).mongodbconn
	logger.Debug(logger.ReqFormat, "AddSecurityGroupRule", createReq)
	resp, err := conn.AddSecurityGroupRule(&createReq)
	if err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	logger.Debug(logger.RespFormat, "AddSecurityGroupRule", createReq, *resp)

	d.SetId(createReq["InstanceId"].(string))

	return resourceMongodbSecurityRuleRead(d, meta)
}

func resourceMongodbSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*KsyunClient).mongodbconn

	rules := d.Get("rules").([]interface{})
	var cidrs []string
	for _, rule := range rules {
		r := rule.(map[string]interface{})
		cidrs = append(cidrs, r["cidr"].(string))
	}

	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId"] = d.Id()
	deleteReq["cidrs"] = strings.Join(cidrs, ",")
	logger.Debug(logger.ReqFormat, "DeleteSecurityGroupRules", deleteReq)
	resp, err := conn.DeleteSecurityGroupRules(&deleteReq)
	if err != nil {
		return fmt.Errorf("error on delete instance security rule: %s", err)
	}
	logger.Debug(logger.RespFormat, "DeleteSecurityGroupRules", deleteReq, *resp)

	return nil
}

func resourceMongodbSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)
	defer d.Partial(false)

	conn := meta.(*KsyunClient).mongodbconn

	if d.HasChange("cidrs") {
		d.SetPartial("cidrs")
		var addRules, deleteRules []string

		oldCidr, newCidr := d.GetChange("cidrs")
		if newCidr.(string) == "" {
			rules := d.Get("rules").([]interface{})
			for _, rule := range rules {
				r := rule.(map[string]interface{})
				deleteRules = append(deleteRules, r["cidr"].(string))
			}
		} else if oldCidr.(string) == "" {
			addRules = strings.Split(newCidr.(string), ",")
		} else {
			oldCidrs := strings.Split(oldCidr.(string), ",")
			newCidrs := strings.Split(newCidr.(string), ",")

			for _, oldRule := range oldCidrs {
				if !strings.Contains(newCidr.(string), oldRule) {
					deleteRules = append(deleteRules, oldRule)
				}
			}

			for _, newRule := range newCidrs {
				if !strings.Contains(oldCidr.(string), newRule) {
					addRules = append(addRules, newRule)
				}
			}
		}
		if len(addRules) > 0 {
			createReq := make(map[string]interface{})
			createReq["InstanceId"] = d.Id()
			createReq["cidrs"] = strings.Join(addRules, ",")
			logger.Debug(logger.ReqFormat, "AddSecurityGroupRule", createReq)
			resp, err := conn.AddSecurityGroupRule(&createReq)
			if err != nil {
				return fmt.Errorf("error on add instance security rule: %s", err)
			}
			logger.Debug(logger.RespFormat, "AddSecurityGroupRule", createReq, *resp)
		}

		if len(deleteRules) > 0 {
			deleteReq := make(map[string]interface{})
			deleteReq["InstanceId"] = d.Id()
			deleteReq["cidrs"] = strings.Join(deleteRules, ",")
			logger.Debug(logger.ReqFormat, "DeleteSecurityGroupRules", deleteReq)
			resp, err := conn.DeleteSecurityGroupRules(&deleteReq)
			if err != nil {
				return fmt.Errorf("error on delete instance security rule: %s", err)
			}
			logger.Debug(logger.RespFormat, "DeleteSecurityGroupRules", deleteReq, *resp)
		}
	}

	return resourceMongodbSecurityRuleRead(d, meta)
}

func resourceMongodbSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*KsyunClient).mongodbconn
	readReq := make(map[string]interface{})
	readReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "ListSecurityGroupRules", readReq)
	resp, err := conn.ListSecurityGroupRules(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading instance security rule %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, "ListSecurityGroupRules", readReq, *resp)

	rules := (*resp)["MongoDBSecurityGroupRule"].([]interface{})

	if err := d.Set("rules", rules); err != nil {
		return fmt.Errorf("error set data:%v", err)
	}

	return nil
}
