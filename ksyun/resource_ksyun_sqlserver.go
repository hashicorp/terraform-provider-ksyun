package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/sqlserver"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunSqlServer() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunSqlServerCreate,
		Update: resourceKsyunSqlServerUpdate,
		Read:   resourceKsyunSqlServerRead,
		Delete: resourceKsyunSqlServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance identifier",
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				Description: "this value regex db.ram.d{1,3}|db.disk.d{1,5} , " +
					"db.ram is rds random access memory size, db.disk is disk size",
			},
			"db_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HRDS_SS",
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "engine is db type, only support SQLServer",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db engine version only support 2008r2,2012,2016",
			},
			"master_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"master_user_password": {
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
			"bill_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "proprietary security group id for krds",
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"availability_zone_1": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"availability_zone_2": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
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
								},
							},
						},
					},
				},
			},
		},
	}
}
func resourceKsyunSqlServerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	var err error
	creates := []string{
		"DBInstanceClass",
		"DBInstanceName",
		"DBInstanceType",
		"Engine",
		"EngineVersion",
		"MasterUserName",
		"MasterUserPassword",
		"VpcId",
		"SubnetId",
		"BillType",
		"Duration",
		"SecurityGroupId",
		"PreferredBackupTime",
		"AvailabilityZone.1",
		"AvailabilityZone.2",
		"ProjectId",
		"Port",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(FuckHump2Downline(v)); ok {
			createReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateDBInstance"
	logger.Debug(logger.RespFormat, action, createReq)
	resp, err = conn.CreateDBInstance(&createReq)
	logger.Debug(logger.AllFormat, action, createReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Instance(sqlserver): %s", err)
	}

	if resp != nil {
		bodyData := (*resp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		instanceId := sqlserverInstance["DBInstanceIdentifier"].(string)
		logger.DebugInfo(" DBInstanceIdentifier : %v", instanceId)
		d.SetId(instanceId)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{tCreatingStatus},
		Target:     []string{tActiveStatus, tFailedStatus, tDeletedStatus, tStopedStatus},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    sqlserverInstanceStateRefreshForCreate(conn, d.Id(), []string{tCreatingStatus}),
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceKsyunSqlServerRead(d, meta)
}

func sqlserverInstanceStateRefreshForCreate(client *sqlserver.Sqlserver, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"DBInstanceIdentifier": instanceId}
		action := "DescribeDBInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeDBInstances(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return nil, "", err
		}
		bodyData := (*resp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)

		return resp, state, nil

	}
}

func resourceKsyunSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	req := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeDBInstances(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Instance(sqlserver) %q, %s", d.Id(), err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading Instance(sqlserver) body %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})

	sqlserverIds := make([]string, len(instances))
	sqlserverMap := make([]map[string]interface{}, len(instances))
	for k, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "DBInstanceClass" {
				dbclass := v.(map[string]interface{})
				dbinstanceclass := make(map[string]interface{})
				for j, q := range dbclass {
					dbinstanceclass[FuckHump2Downline(j)] = q
				}
				wtf := make([]interface{}, 1)
				wtf[0] = dbinstanceclass
				instanceInfo["db_instance_class"] = wtf
				delete(instanceInfo, "DBInstanceClass")
			} else {
				delete(instanceInfo, k)
				instanceInfo[FuckHump2Downline(k)] = v
			}
		}
		sqlserverMap[k] = instanceInfo
		logger.DebugInfo(" converted ---- %+v ", instanceInfo)

		sqlserverIds[k] = instanceInfo["db_instance_identifier"].(string)
		logger.DebugInfo("sqlserverIds fuck : %v", sqlserverIds)
	}

	logger.DebugInfo(" converted ---- %+v ", sqlserverMap)
	dataSourceDbSave(d, "sqlservers", sqlserverIds, sqlserverMap)

	return nil
}

func resourceKsyunSqlServerUpdate(d *schema.ResourceData, meta interface{}) error {
	// 关闭事务，允许部分属性被修改  d.Partial(true) d.Partial(false)
	updateField := []string{
		"db_instance_class",
		"db_instance_name",
		"db_instance_type",
		"engine",
		"engine_version",
		"master_user_name",
		"master_user_password",
		"vpc_id",
		"subnet_id",
		"bill_type",
		"duration",
		"security_group_id",
		"preferred_backup_time",
		"availability_zone_1",
		"availability_zone_2",
		"project_id",
		"port",
	}
	d.Partial(true)
	for _, v := range updateField {
		if d.HasChange(v) {
			return fmt.Errorf("error on updating instance , sqlserver is not support update")
		}
	}
	d.Partial(false)
	return nil
}

func resourceKsyunSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	deleteReq := make(map[string]interface{})
	deleteReq["DBInstanceIdentifier"] = d.Id()

	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		readReq := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
		discribeAction := "DescribeInstances"
		logger.Debug(logger.ReqFormat, discribeAction, readReq)
		desResp, desErr := conn.DescribeDBInstances(&readReq)
		logger.Debug(logger.AllFormat, discribeAction, readReq, *desResp, desErr)

		if desErr != nil {
			if notFoundError(desErr) {
				return nil
			} else {
				return resource.NonRetryableError(desErr)
			}
		}

		bodyData := (*desResp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)

		if state != tDeletedStatus {
			deleteAction := "DeleteDBInstance"
			logger.Debug(logger.ReqFormat, deleteAction, deleteReq)
			deleteResp, deleteErr := conn.DeleteDBInstance(&deleteReq)
			logger.Debug(logger.AllFormat, deleteAction, deleteReq, *deleteResp, deleteErr)
			if deleteErr == nil || notFoundError(deleteErr) {
				return nil
			}
			if deleteErr != nil {
				return resource.RetryableError(deleteErr)
			}

			logger.Debug(logger.ReqFormat, discribeAction, readReq)
			postDesResp, postDesErr := conn.DescribeDBInstances(&readReq)
			logger.Debug(logger.AllFormat, discribeAction, readReq, *postDesResp, postDesErr)

			if desErr != nil {
				if notFoundError(desErr) {
					return nil
				} else {
					return resource.NonRetryableError(fmt.Errorf("error on  reading kec when delete %q, %s", d.Id(), desErr))
				}
			}
		}

		return resource.RetryableError(desErr)
	})
}
