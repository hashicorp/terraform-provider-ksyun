package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func dataSourceKsyunKrds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKrdsRead,
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
			"db_instance_status": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "ACTIVE（运行中）/INVALID（请续费）",
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
			"krds": {
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
							Computed: true,
						},
						"db_instance_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db_instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"publicly_accessible": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"bill_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"multi_availability_zone": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"db_parameter_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_used": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"preferred_backup_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKrdsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	desReq := make(map[string]interface{})
	des := []string{
		"DBInstanceIdentifier",
		"DBInstanceType",
		"DBInstanceStatus",
		"Keyword",
		"Order",
		"ProjectId",
		"Marker",
		"MaxRecords",
	}
	for _, v := range des {
		if v1, ok := d.GetOk(Camel2Hungarian(v)); ok {
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
		return fmt.Errorf("error on reading Instance(krds) body %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})

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
			} else {
				dk := Camel2Hungarian(k)
				if _, ok := dataSourceKsyunKrds().Schema["krds"].Elem.(*schema.Resource).Schema[dk]; ok {
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
	_ = dataDbSave(d, "krds", krdsIds, krdsMapList)

	return nil
}
