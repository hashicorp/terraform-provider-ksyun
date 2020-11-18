package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"log"
	"strconv"
	"time"
)

var krdsTfField = []string{
	"db_instance_identifier",
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
	"db_parameter_group_id",
	"preferred_backup_time",
	"availability_zone_1",
	"availability_zone_2",
	"project_id",
	"port",
	"vip",
	"instance_has_eip",
	"eip",
	"eip_port",
}

var getInTheCar = map[string]bool{
	"db_instance_identifier": true,
	"instance_create_time":   true,
	"port":                   true,
	"db_parameter_group_id":  true,
	"sub_order_id":           true,
	"availability_zone_1":    true,
	"availability_zone_2":    true,
	"region":                 true,
}

func resourceKsyunKrds() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunKrdsCreate,
		Update: resourceKsyunMysqlUpdate,
		Read:   resourceKsyunMysqlRead,
		Delete: resourceKsyunMysqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_identifier": {
				Computed:    true,
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
				Description: "HRDS",
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "engine is db type, only support mysql|percona",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db engine version only support 5.5|5.6|5.7|8.0",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"db_parameter_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "proprietary db parameter group id for mysql",
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
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Computed: true,
			},
			"sub_order_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_create_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_has_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
			},
			"eip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"eip_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
		},
	}
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceKsyunKrdsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
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
		if v1, ok := d.GetOk(Camel2Hungarian(v)); ok {
			createReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateDBInstance"
	logger.Debug(logger.RespFormat, action, createReq)
	resp, err = conn.CreateDBInstance(&createReq)
	logger.Debug(logger.AllFormat, action, createReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Instance(krds): %s", err)
	}

	if resp != nil {
		bodyData := (*resp)["Data"].(map[string]interface{})
		krdsInstance := bodyData["DBInstance"].(map[string]interface{})
		instanceId := krdsInstance["DBInstanceIdentifier"].(string)
		logger.DebugInfo("~*~*~*~*~ DBInstanceIdentifier : %v", instanceId)
		d.SetId(instanceId)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{tCreatingStatus},
		Target:     []string{tActiveStatus, tFailedStatus, tDeletedStatus, tStopedStatus},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    mysqlInstanceStateRefresh(conn, d.Id(), []string{tCreatingStatus}),
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on creating Instance(krds): %s", err)
	}

	err = resourceKsyunMysqlRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on ModifyDBParameterGroup Instance(krds): %s", err)
	}

	if d.Get("instance_has_eip") == true {
		modifyInstanceEip(d, meta)
	}
	return modifyParameters(d, meta)
}

func modifyParameters(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	paramsReq := make(map[string]interface{})
	paramsReq["DBParameterGroupId"] = d.Get("db_parameter_group_id").(string)
	documented := d.Get("parameters").(*schema.Set).List()
	if len(documented) > 0 {
		for paramIndex, i := range documented {
			num := paramIndex + 1
			paramsReq["Parameters.Name."+strconv.Itoa(num)] = i.(map[string]interface{})["name"].(string)
			paramsReq["Parameters.Value."+strconv.Itoa(num)] = i.(map[string]interface{})["value"].(string)
		}
		mdAction := "ModifyDBParameterGroup"
		logger.Debug(logger.RespFormat, mdAction, paramsReq)
		paramResp, err := conn.ModifyDBParameterGroup(&paramsReq)
		logger.Debug(logger.AllFormat, mdAction, paramsReq, *paramResp, err)
		return err
	}
	return nil
}

func modifyInstanceEip(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	_ = checkStatus(d, conn)
	req := map[string]interface{}{
		"DBInstanceIdentifier": d.Id(),
	}

	if d.Get("instance_has_eip") == true {
		action := "AllocateDBInstanceEip"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.AllocateDBInstanceEip(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)

		if err != nil {
			return err
		}
		d.SetPartial("eip")
		d.SetPartial("eip_port")
		return nil
	} else {
		action := "ReleaseDBInstanceEip"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ReleaseDBInstanceEip(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)

		if err != nil {
			return err
		}
		d.SetPartial("eip")
		d.SetPartial("eip_port")
		return nil
	}
}

func mysqlInstanceStateRefresh(client *krds.Krds, instanceId string, target []string) resource.StateRefreshFunc {
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
		krdsInstance := instances[0].(map[string]interface{})
		state := krdsInstance["DBInstanceStatus"].(string)

		return resp, state, nil
	}
}

func resourceKsyunMysqlRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	req := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeDBInstances(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)

	if err != nil {
		return fmt.Errorf("error on reading Instance(krds) %q, %s", d.Id(), err)
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
			} else if k == "ReadReplicaDBInstanceIdentifiers" {
			} else if k == "DBSource" {
			} else {
				krdsMap[Camel2Hungarian(k)] = v
			}
			if k == "Eip" {
				krdsMap["instance_has_eip"] = true
			}
		}
		if krdsMap["instance_has_eip"] == nil {
			krdsMap["instance_has_eip"] = false
		}
		logger.DebugInfo(" converted ---- %+v ", krdsMap)

		krdsIds[num] = krdsMap["db_instance_identifier"].(string)
		logger.DebugInfo("krdsIds fuck : %v", krdsIds)
		krdsMapList[num] = krdsMap
	}
	logger.DebugInfo(" converted ---- %+v ", krdsMapList)
	_ = SetDByFkResp(d, krdsMapList[0], getInTheCar)

	return nil
}

func resourceKsyunMysqlUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	// 允许部分属性被修改  d.Partial(true) d.Partial(false)
	execModifyDBInstanceSpec := false
	execModifyDBInstance := false
	execModifyDBInstanceType := false
	execUpgradeDBInstanceEngineVersion := false
	execModifyDBBackupPolicy := false
	execModifyParameters := false
	execModifyDBInstanceAvailabilityZone := false
	execModifyDBInstanceEip := false
	d.Partial(true)
	for _, v := range krdsTfField {
		if d.HasChange(v) && !d.IsNewResource() {
			if v == "engine" || v == "master_user_name" || v == "vpc_id" || v == "subnet_id" || v == "bill_type" || v == "duration" || v == "db_instance_identifier" || v == "project_id" || v == "db_parameter_group_id" {
				return fmt.Errorf("error on updating instance , krds is not support update %s", v)
			}
			if v == "db_instance_class" {
				execModifyDBInstanceSpec = true
			}
			if v == "db_instance_name" {
				execModifyDBInstance = true
			}
			if v == "db_instance_type" {
				execModifyDBInstanceType = true
			}
			if v == "engine_version" {
				execUpgradeDBInstanceEngineVersion = true
			}
			if v == "master_user_password" {
				execModifyDBInstance = true
			}
			if v == "security_group_id" {
				execModifyDBInstance = true
			}
			if v == "preferred_backup_time" {
				execModifyDBBackupPolicy = true
			}
			if v == "availability_zone_1" || v == "availability_zone_2" {
				execModifyDBInstanceAvailabilityZone = true
			}
			if v == "instance_has_eip" {
				execModifyDBInstanceEip = true
			}

			oldType, _ := d.GetChange("db_instance_type")
			if "RR" == oldType {
				if v == "availability_zone_1" || v == "availability_zone_2" || v == "preferred_backup_time" || v == "db_instance_type" {
					return fmt.Errorf("error on updating instance , krds rr is not support update %s", v)
				}
			}
		}
	}
	if d.HasChange("parameters") {
		execModifyParameters = true
	}
	log.Printf(" if the response status code is 409, the instance is doing other things, " +
		"please wait several minutes and retry")

	if execModifyDBInstance {
		req := map[string]interface{}{
			"DBInstanceIdentifier": d.Id(),
		}
		if v, ok := d.GetOk("db_instance_name"); ok {
			req["DBInstanceName"] = v
		}
		if v, ok := d.GetOk("master_user_password"); ok {
			req["MasterUserPassword"] = v
		}
		if v, ok := d.GetOk("security_group_id"); ok {
			req["SecurityGroupId"] = v
		}
		action := "ModifyDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ModifyDBInstance(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return err
		}
		if d.HasChange("db_instance_name") {
			d.SetPartial("db_instance_name")
		}
		if d.HasChange("master_user_password") {
			d.SetPartial("master_user_password")
		}
		if d.HasChange("security_group_id") {
			d.SetPartial("security_group_id")
		}
	}
	if execModifyDBBackupPolicy {
		req := map[string]interface{}{
			"DBInstanceIdentifier": d.Id(),
			"PreferredBackupTime":  d.Get("preferred_backup_time"),
		}
		action := "ModifyDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ModifyDBInstance(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return err
		}
		d.SetPartial("preferred_backup_time")
	}
	if execModifyParameters {
		err := modifyParameters(d, meta)
		if err != nil {
			return err
		}
		d.SetPartial("parameters")
	}

	if execModifyDBInstanceType {
		_ = checkStatus(d, conn)

		oldType, newType := d.GetChange("db_instance_type")
		if "TRDS" == oldType && "HRDS" == newType {
			req := map[string]interface{}{
				"DBInstanceIdentifier": d.Id(),
				"DBInstanceType":       d.Get("db_instance_type"),
			}
			action := "ModifyDBInstanceType"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err := conn.ModifyDBInstanceType(&req)
			logger.Debug(logger.AllFormat, action, req, *resp, err)
			if err != nil {
				return err
			}
			if d.HasChange("db_instance_type") {
				d.SetPartial("db_instance_type")
			}
		} else {
			return fmt.Errorf("error on updating instance , krds is not support %s to %s", oldType, newType)
		}
	}
	if execModifyDBInstanceAvailabilityZone {
		_ = checkStatus(d, conn)

		req := map[string]interface{}{
			"DBInstanceIdentifier": d.Id(),
			"AvailabilityZone.1":   d.Get("availability_zone_1"),
			"AvailabilityZone.2":   d.Get("availability_zone_2"),
		}
		action := "ModifyDBInstanceAvailabilityZone"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ModifyDBInstanceAvailabilityZone(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return err
		}
		if d.HasChange("availability_zone_1") {
			d.SetPartial("availability_zone_1")
		}
		if d.HasChange("availability_zone_2") {
			d.SetPartial("availability_zone_2")
		}
	}
	if execUpgradeDBInstanceEngineVersion {
		_ = checkStatus(d, conn)

		req := map[string]interface{}{
			"DBInstanceIdentifier": d.Id(),
			"Engine":               d.Get("engine"),
			"EngineVersion":        d.Get("engine_version"),
		}
		action := "UpgradeDBInstanceEngineVersion"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.UpgradeDBInstanceEngineVersion(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return err
		}
		d.SetPartial("engine_version")
	}

	if execModifyDBInstanceSpec {
		_ = checkStatus(d, conn)

		req := map[string]interface{}{
			"DBInstanceIdentifier": d.Id(),
			"DBInstanceClass":      d.Get("db_instance_class"),
		}
		action := "ModifyDBInstanceSpec"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ModifyDBInstanceSpec(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)

		if err != nil {
			return err
		}
		d.SetPartial("db_instance_class")
	}
	if execModifyDBInstanceEip {
		modifyInstanceEip(d, meta)
	}
	d.Partial(false)

	_ = checkStatus(d, conn)
	return resourceKsyunMysqlRead(d, meta)
}

func checkStatus(d *schema.ResourceData, conn *krds.Krds) error {
	stateConf := &resource.StateChangeConf{
		Pending:    waitStatus,
		Target:     finalStatus,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    mysqlInstanceStateRefresh(conn, d.Id(), finalStatus),
	}
	_, err := stateConf.WaitForState()

	//fmt.Println("checkResp is ", checkResp)
	if err != nil {
		return fmt.Errorf("error on updating ModifyDBInstanceType , err = %s", err)
	} else {
		_, instFinalStatus, _ := mysqlInstanceStateRefresh(conn, d.Id(), finalStatus)()

		if instFinalStatus != tActiveStatus {
			return fmt.Errorf("error status : %s, ", instFinalStatus)
		} else {
			return nil
		}
	}
}

func resourceKsyunMysqlDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
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
					return resource.NonRetryableError(fmt.Errorf("error on  reading krds when delete %q, %s", d.Id(), desErr))
				}
			}
		}

		return resource.RetryableError(desErr)
	})
}
