package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunEpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunEpcsRead,

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

			"host_name": {
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

			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cabinet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"host_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"epc_host_status": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"os_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"product_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cluster_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enable_container": {
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
			"epchosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"host_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"host_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"sn": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"raid": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_interface_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"security_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cloud_monitor_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"allow_modify_hyper_threading": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},

						"releasable_time": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"tor_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"cabinet_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"rack_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"cabinet_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"enable_bond": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},

						"product_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"os_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"memory": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"enable_container": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},

						"system_file_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"data_file_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"data_disk_catalogue": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"data_disk_catalogue_suffix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"hyper_threading": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"cpu": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"frequence": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"core_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"gpu": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"frequence": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"core_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"gpu_count": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"disk_set": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"raid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"space": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"system_disk_space": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"network_interface_attribute_set": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							MaxItems: 2,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_interface_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"network_interface_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"private_ip_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"d_n_s1": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"d_n_s2": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"mac": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"security_group_set": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"security_group_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
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

func dataSourceKsyunEpcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.epcconn
	readEpc := make(map[string]interface{})
	var hostIds []string
	if ids, ok := d.GetOk("ids"); ok {
		hostIds = SchemaSetToStringSlice(ids)
	}
	for k, v := range hostIds {
		if v == "" {
			continue
		}
		readEpc[fmt.Sprintf("HostId.%d", k+1)] = v
	}
	var projectIds []string
	if ids, ok := d.GetOk("project_id"); ok {
		projectIds = SchemaSetToStringSlice(ids)
	}

	for k, v := range projectIds {
		if v == "" {
			continue
		}
		readEpc[fmt.Sprintf("ProjectId.%d", k+1)] = v
	}

	filters := []string{
		"host_name",
		"vpc_id",
		"subnet_id",
		"cabinet_id",
		"host_type",
		"epc_host_status",
		"os_name",
		"product_type",
		"cluster_id",
		"enable_container",
	}
	SchemaSetsToFilterMap(d, filters, &readEpc)
	var maxResults int = 100
	var nextToken string
	var allHosts []interface{}

	for {
		readEpc["MaxResults"] = fmt.Sprintf("%v", maxResults)
		if nextToken != "" {
			readEpc["NextToken"] = nextToken
		}
		action := "DescribeEpcs"
		logger.Debug(logger.ReqFormat, action, readEpc)
		resp, err := conn.DescribeEpcs(&readEpc)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", readEpc, err)
		}
		logger.Debug(logger.RespFormat, action, readEpc, *resp)
		itemSet, ok := (*resp)["HostSet"]
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
		allHosts = append(allHosts, items...)

		if len(items) < maxResults {
			break
		}
		if nextTokens, ok := (*resp)["NextToken"]; ok {
			nextToken = fmt.Sprintf("%v", nextTokens)
		} else {
			break
		}
	}

	datas := GetSubSliceDByRep(allHosts, epcInstanceKeys)
	dealEpcInstanceData(datas)
	err := dataSourceKsyunEpcsSave(d, "epchosts", hostIds, datas)
	if err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}
	return nil
}

func dealEpcInstanceData(datas []map[string]interface{}) {
	for k, v := range datas {
		for k1, v1 := range v {
			logger.Debug(logger.ReqFormat, k1, v1)
			switch k1 {
			case "cpu":
				datas[k]["cpu"] = GetSubDByRep(v1, epcCpuKeys, map[string]bool{})
			case "gpu":
				datas[k]["gpu"] = GetSubDByRep(v1, epcGpuKeys, map[string]bool{})
			case "disk_set":
				datas[k]["disk_set"] = GetSubSliceDByRep(v1.([]interface{}), epcDiskSetKeys)
			case "network_interface_attribute_set":
				vv := v1.([]interface{})
				networkSet := GetSubSliceDByRep(vv, epcNetworkInterfaceKeys)
				for itemK, itemV := range networkSet {
					for k2, v2 := range itemV {
						switch k2 {
						case "security_group_set":
							vv := v2.([]interface{})
							networkSet[itemK]["security_group_set"] = GetSubSliceDByRep(vv, epcSecurityGroupKeys)
						}
					}
				}
				datas[k]["network_interface_attribute_set"] = networkSet
			}
		}
	}
}

func dataSourceKsyunEpcsSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {

	d.SetId(hashStringArray(ids))
	d.Set("total_count", len(datas))

	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), datas)
	}
	return nil
}
