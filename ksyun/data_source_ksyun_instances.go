package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_interface": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"network_interface_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"group_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"instance_state": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_configure": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"v_c_p_u": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"g_p_u": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"memory_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"data_disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"data_disk_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"root_disk_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_state": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitoring": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"sriov_net_support": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_set": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
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
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_set": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"security_group_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"group_set": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
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
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"system_disk": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stopped_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_what": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_scaling_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_show_sriov_net_support": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunInstancesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).kecconn
	var instances []string
	req := make(map[string]interface{})

	if ids, ok := d.GetOk("ids"); ok {
		instances = SchemaSetToStringSlice(ids)
	}
	for k, v := range instances {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("InstanceId.%d", k+1)] = v
	}
	var projectIds []string
	if ids, ok := d.GetOk("project_id"); ok {
		projectIds = SchemaSetToStringSlice(ids)
	}

	for k, v := range projectIds {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("ProjectId.%d", k+1)] = v
	}

	filters := []string{
		"instance_id",
		"subnet_id",
		"vpc_id",
		//"network_interface",
		//	"instance_state",
		//	"availability_zone",
	}
	index := len(req) + 1
	SchemaSetsToFilterMap(d, filters, &req)

	netWorkreq := make(map[string]interface{})
	if v, ok := d.GetOk("network_interface"); ok {
		ConvertFilterStructPrefix(v, &netWorkreq, "network-interface")
		for k1, v1 := range netWorkreq {
			v2 := SchemaSetToStringSlice(v1)
			if len(v2) == 0 {
				continue
			}
			for k3, v3 := range v2 {
				req[fmt.Sprintf("Filter.%v.Value.%d", index, k3+1)] = v3
			}
			req[fmt.Sprintf("Filter.%v.Name", index)] = k1
			index++
		}
	}
	statereq := make(map[string]interface{})
	if v, ok := d.GetOk("instance_state"); ok {
		ConvertFilterStructPrefix(v, &statereq, "instance-state")
		for k1, v1 := range statereq {
			v2 := SchemaSetToStringSlice(v1)
			if len(v2) == 0 {
				continue
			}
			for k3, v3 := range v2 {
				req[fmt.Sprintf("Filter.%v.Value.%d", index, k3+1)] = v3
			}
			req[fmt.Sprintf("Filter.%v.Name", index)] = k1
			index++
		}
	}
	zonereq := make(map[string]interface{})
	if v, ok := d.GetOk("availability_zone"); ok {
		ConvertFilterStructPrefix(v, &zonereq, "availability-zone")
		for k1, v1 := range zonereq {
			v2 := SchemaSetToStringSlice(v1)
			if len(v2) == 0 {
				continue
			}
			for k3, v3 := range v2 {
				req[fmt.Sprintf("Filter.%v.Value.%d", index, k3+1)] = v3
			}
			req[fmt.Sprintf("Filter.%v.Name", index)] = k1
			index++
		}
	}
	var allinstances []interface{}
	var limit int = 100
	var nextToken string
	for {
		req["MaxResults"] = fmt.Sprintf("%v", limit)
		if nextToken != "" {
			req["NextToken"] = nextToken
		}
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.DescribeInstances(&req)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", req, err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		itemSet, ok := (*resp)["InstancesSet"]
		if !ok {
			//return fmt.Errorf("error on reading instance set")
			break
		}
		items, ok := itemSet.([]interface{})
		if !ok {
			break
		}
		if items == nil || len(items) < 1 {
			break
		}
		allinstances = append(allinstances, items...)

		if len(items) < limit {
			break
		}
		if nextTokens, ok := (*resp)["Marker"]; ok {
			nextToken = fmt.Sprintf("%v", nextTokens)
		} else {
			break
		}
	}
	datas := GetSubSliceDByRep(allinstances, instanceKeys)
	dealInstanceData(datas)
	err := dataSourceKscSave(d, "instances", instances, datas)
	if err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}
	return nil
}
func dealInstanceData(datas []map[string]interface{}) {
	for k, v := range datas {
		for k1, v1 := range v {
			switch k1 {
			case "instance_configure":
				datas[k]["instance_configure"] = GetSubDByRep(v1, instanceConfigureKeys, map[string]bool{})
			case "system_disk":
				datas[k]["system_disk"] = GetSubDByRep(v1, systemDiskKeys, map[string]bool{})
			case "monitoring":
				datas[k]["monitoring"] = GetSubDByRep(v1, monitoringKeys, map[string]bool{})
			case "instance_state":
				datas[k]["instance_state"] = GetSubDByRep(v1, instanceStateKeys, map[string]bool{})
			case "key_set":
				datas[k]["key_id"] = v1
				delete(datas[k], "key_set")
			case "network_interface_set":
				vv := v1.([]interface{})
				networkSet := GetSubSliceDByRep(vv, kecNetworkInterfaceSetKeys)
				for itemK, itemV := range networkSet {
					for k2, v2 := range itemV {
						switch k2 {
						case "group_set":
							vv := v2.([]interface{})
							networkSet[itemK]["group_set"] = GetSubSliceDByRep(vv, groupSetKeys)
						case "security_group_set":
							vv := v2.([]interface{})
							networkSet[itemK]["security_group_set"] = GetSubSliceDByRep(vv, kecSecurityGroupSetKeys)
						}
					}
				}
				datas[k]["network_interface_set"] = networkSet
			}
		}
	}
}
