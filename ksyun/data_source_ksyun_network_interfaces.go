package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKsyunNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunNetworkInterfacesRead,
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
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"securitygroup_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_ip_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"d_n_s1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"d_n_s2": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunNetworkInterfacesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).vpcconn
	req := make(map[string]interface{})
	var NetworkInterfaceIds []string
	if ids, ok := d.GetOk("ids"); ok {
		NetworkInterfaceIds = SchemaSetToStringSlice(ids)
	}
	for k, v := range NetworkInterfaceIds {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("NetworkInterfaceId.%d", k+1)] = v
	}
	filters := []string{
		"vpc_id",
		"subnet_id",
		"securitygroup_id",
		"instance_type",
		"instance_id",
		"private_ip_address",
	}
	req = *SchemaSetsToFilterMap(d, filters, &req)
	var allNetworkInterfaces []interface{}
	resp, err := conn.DescribeNetworkInterfaces(&req)
	if err != nil {
		return fmt.Errorf("error on reading NetworkInterface list, %s", req)
	}
	itemSet, ok := (*resp)["NetworkInterfaceSet"]
	if !ok {
		return fmt.Errorf("error on reading NetworkInterface set")
	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allNetworkInterfaces = append(allNetworkInterfaces, items...)
	datas := GetSubSliceDByRep(allNetworkInterfaces, vpcNetworkInterfaceKeys)
	for _, v := range datas {
		securitygroup := v["security_group_set"]
		if securityGroups, ok := securitygroup.([]interface{}); ok {
			securitygroupSets := GetSubSliceDByRep(securityGroups, groupIdentifierKeys)
			v["security_group_set"] = securitygroupSets
		}
	}
	err = dataSourceKscSave(d, "network_interfaces", NetworkInterfaceIds, datas)
	if err != nil {
		return fmt.Errorf("error on save NetworkInterface list, %s", err)
	}
	return nil
}
