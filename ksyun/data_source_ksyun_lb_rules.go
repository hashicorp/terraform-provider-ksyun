package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunSlbRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSlbRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"host_header_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lb_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_header_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backend_server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_sync": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"method": {
							Type:     schema.TypeString,
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
										Computed: true,
									},
									"health_check_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"healthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"url_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"session": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"session_persistence_period": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"session_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSlbRulesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var slbRules []string
	if ids, ok := d.GetOk("ids"); ok {
		slbRules = SchemaSetToStringSlice(ids)
	}
	for k, v := range slbRules {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("RuleId.%d", k+1)] = v
	}
	filters := []string{"host_header_id"}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	var allSlbRules []interface{}

	resp, err := conn.DescribeRules(&req)
	if err != nil {
		return fmt.Errorf("error on reading slb rule list req (%v):%v", req, err)
	}
	itemSet, ok := (*resp)["RuleSet"]
	if !ok {
		return fmt.Errorf("error on reading slb rule set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allSlbRules = append(allSlbRules, items...)
	//	excludes:=[]string{"HealthCheck","RealServer","Session"}
	slbRuleDataKeys := slbRuleKeys
	slbRuleDataKeys["RuleId"] = true
	datas := GetSubSliceDByRep(allSlbRules, slbRuleDataKeys)
	dealListenrData(datas)
	err = dataSourceKscSave(d, "lb_rules", slbRules, datas)
	if err != nil {
		return fmt.Errorf("error on save lb rule list, %s", err)
	}
	return nil
}
