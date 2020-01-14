package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKsyunSlbAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSlbAclsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
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
			"lb_acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_entry_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_acl_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"load_balancer_acl_entry_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_number": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
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
func dataSourceKsyunSlbAclsRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var slbAcls []string
	if ids, ok := d.GetOk("ids"); ok {
		slbAcls = SchemaSetToStringSlice(ids)
	}
	for k, v := range slbAcls {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("LoadBalancerAclId.%d", k+1)] = v
	}
	var allSlbAcls []interface{}
	var limit int = 100
	var nextToken string
	for {
		req["MaxResults"] = limit
		if nextToken != "" {
			req["NextToken"] = nextToken
		}
		resp, err := conn.DescribeLoadBalancerAcls(&req)
		if err != nil {
			return fmt.Errorf("error on reading SlbAcl list, %s", req)
		}

		itemSet, ok := (*resp)["LoadBalancerAclSet"]
		if !ok {
			return fmt.Errorf("error on reading SlbAcl set")

		}
		//???????????????????????????
		items, ok := itemSet.([]interface{})
		if !ok {
			break
		}
		if items == nil || len(items) < 1 {
			break
		}
		allSlbAcls = append(allSlbAcls, items...)

		if len(items) < limit {
			break
		}

		if nextTokens, ok := (*resp)["NextToken"]; ok {
			nextToken = fmt.Sprintf("%v", nextTokens)
		} else {
			break
		}
	}

	datas := GetSubSliceDByRep(allSlbAcls, lbAclKeys)
	dealSlbAclData(datas)
	err := dataSourceKscSave(d, "lb_acls", slbAcls, datas)
	if err != nil {
		return fmt.Errorf("error on save SlbAcl list, %s", err)
	}
	return nil

}

func dealSlbAclData(datas []map[string]interface{}) {
	for k, v := range datas {
		for k1, v1 := range v {
			switch k1 {
			case "load_balancer_acl_entry_set":
				vv := v1.([]interface{})
				datas[k]["load_balancer_acl_entry_set"] = GetSubSliceDByRep(vv, lbAclEntryKeys)
			}
		}
	}
}
