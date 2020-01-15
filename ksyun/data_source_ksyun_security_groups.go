package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKsyunSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSecurityGroupsRead,
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
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_entry_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_entry_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"icmp_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_range_from": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_range_to": {
										Type:     schema.TypeInt,
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

func dataSourceKsyunSecurityGroupsRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).vpcconn
	req := make(map[string]interface{})
	var SecurityGroups []string
	if ids, ok := d.GetOk("ids"); ok {
		SecurityGroups = SchemaSetToStringSlice(ids)
	}
	for k, v := range SecurityGroups {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("SecurityGroupId.%d", k+1)] = v
	}
	filters := []string{"vpc_id"}
	req = *SchemaSetsToFilterMap(d, filters, &req)
	var allSecurityGroups []interface{}

	resp, err := conn.DescribeSecurityGroups(&req)
	if err != nil {
		return fmt.Errorf("error on reading SecurityGroup list req(%v):%v", req, err)
	}
	itemSet, ok := (*resp)["SecurityGroupSet"]
	if !ok {
		return fmt.Errorf("error on reading SecurityGroup set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allSecurityGroups = append(allSecurityGroups, items...)
	datas := GetSubSliceDByRep(allSecurityGroups, vpcSecurityGroupKeys)
	for k, v := range datas {
		securitySet := v["security_group_entry_set"]
		if securitySets, ok := securitySet.([]interface{}); ok {
			securitys := GetSubSliceDByRep(securitySets, vpcSecurityGroupEntrySetKeys)
			datas[k]["security_group_entry_set"] = securitys
		}
	}
	err = dataSourceKscSave(d, "security_groups", SecurityGroups, datas)
	if err != nil {
		return fmt.Errorf("error on save SecurityGroup list, %s", err)
	}
	return nil
}
