package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func dataSourceKsyunSqlServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSqlServerRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db_instance_identifier": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "source instance identifier",
			},

			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "HRDS（高可用）,RR（只读实例）,TRDS（临时实例）",
			},
			"keyword": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"marker": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"max_records": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			// 与存入数据一致datakey
			"sqlservers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_class": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vcpus": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"disk": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"max_conn": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"mem": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"db_instance_identifier": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"db_instance_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"db_instance_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"db_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"instance_create_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"master_user_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"publicly_accessible": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"read_replica_db_instance_identifiers": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"read_replica_db_instance_identifier": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"bill_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_source": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"multi_availability_zone": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_use": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bill_type_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"db_parameter_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"datastore_version_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_used": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
						"preferred_backup_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"product_what": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"service_start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_order_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"audit": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"db_source": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_instance_identifier": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"db_instance_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"db_instance_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"point_in_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"service_end_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"eip_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"rip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	desReq := make(map[string]interface{})
	des := []string{
		"DBInstanceStatus",
		"DBInstanceType",
		"DBInstanceIdentifier",
		"Keyword",
		"ExpiryDateLessThan",
		"Marker",
		"MaxRecords",
	}
	for _, v := range des {
		if v1, ok := d.GetOk(strings.ToLower(v)); ok {
			desReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, desReq)
	resp, err := conn.DescribeDBInstances(&desReq)
	logger.Debug(logger.AllFormat, action, desReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Instance(sqlserver)  %s", err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading Instance(sqlserver) body %+v", (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})
	if len(instances) == 0 {
		return fmt.Errorf("empty on reading Instance(sqlserver) body %+v", *resp)
	}

	krdsIds := make([]string, len(instances))
	krdsMapList := make([]map[string]interface{}, len(instances))
	for num, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		krdsMap := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "DBInstanceClass" {
				dbclass := v.(map[string]interface{})
				dbinstanceclass := make(map[string]interface{})
				for j, q := range dbclass {
					dbinstanceclass[Camel2Hungarian(j)] = q
				}
				// shit 这里不传list会出现各种报错，我日了
				wtf := make([]interface{}, 1)
				wtf[0] = dbinstanceclass
				krdsMap["db_instance_class"] = wtf
			} else if k == "ReadReplicaDBInstanceIdentifiers" {
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
					krdsMap["read_replica_db_instance_identifiers"] = wtf
				}
			} else if k == "DBSource" {
				dbsource := v.(map[string]interface{})
				dbsourcemap := make(map[string]interface{})
				for j, q := range dbsource {
					dbsourcemap[Camel2Hungarian(j)] = q
				}
				wtf := make([]interface{}, 1)
				wtf[0] = dbsourcemap
				krdsMap["db_source"] = wtf
			} else {
				dk := Camel2Hungarian(k)
				if _, ok := dataSourceKsyunSqlServer().Schema["sqlservers"].Elem.(*schema.Resource).Schema[dk]; ok {
					krdsMap[dk] = v
				}
			}
		}
		logger.DebugInfo(" converted ---- %+v ", krdsMap)

		krdsIds[num] = krdsMap["db_instance_identifier"].(string)
		logger.DebugInfo("krdsIds fuck : %v", krdsIds)
		krdsMapList[num] = krdsMap
	}

	logger.DebugInfo(" converted ---- %+v ", krdsMapList)
	_ = dataDbSave(d, "sqlservers", krdsIds, krdsMapList)

	return nil
}
