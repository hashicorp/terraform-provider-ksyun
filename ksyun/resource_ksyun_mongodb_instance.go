package ksyun

import (
	"errors"
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/mongodb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunMongodbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbInstanceCreate,
		Delete: resourceMongodbInstanceDelete,
		Update: resourceMongodbInstanceUpdate,
		Read:   resourceMongodbInstanceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timing_switch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_what": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongos_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shard_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"area": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongodbInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	createReq := make(map[string]interface{})
	creates := []string{
		"name",
		"instance_account",
		"instance_password",
		"instance_class",
		"storage",
		"node_num",
		"vpc_id",
		"vnet_id",
		"db_version",
		"pay_type",
		"duration",
		"iam_project_id",
		"availability_zone",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok && v != "availability_zone" {
			createReq[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
		}
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		azs := strings.Split(v.(string), ",")
		if len(azs) > 1 {
			return fmt.Errorf("error on creating instance: %s", " not support multiple availableZone")
		}
		for i, az := range azs {
			createReq[fmt.Sprintf("AvailabilityZone.%v", i+1)] = az
		}
	}
	logger.Debug(logger.ReqFormat, "CreateMongoDBInstance", createReq)

	conn := meta.(*KsyunClient).mongodbconn
	resp, err := conn.CreateMongoDBInstance(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	logger.Debug(logger.RespFormat, "CreateMongoDBInstance", createReq, *resp)

	if resp != nil {
		d.SetId((*resp)["MongoDBInstanceResult"].(map[string]interface{})["InstanceId"].(string))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"running"},
		Refresh:    mongodbInstanceStateRefreshForCreateFunc(conn, d.Id(), []string{"creating", "running"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	return resourceMongodbInstanceRead(d, meta)
}

func resourceMongodbInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*KsyunClient).mongodbconn

	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DeleteMongoDBInstance", deleteReq)

	// wait for order update
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteMongoDBInstance(&deleteReq)
		if err != nil {
			return resource.RetryableError(errors.New(""))
		} else {
			return nil
		}
	})

	if err != nil {
		return fmt.Errorf("error on deleting instance %q, %s", d.Id(), err)
	}

	return resource.Retry(20*time.Minute, func() *resource.RetryError {

		queryReq := make(map[string]interface{})
		queryReq["InstanceId"] = d.Id()

		logger.Debug(logger.ReqFormat, "DescribeMongoDBInstance", queryReq)
		resp, err := conn.DescribeMongoDBInstance(&queryReq)
		logger.Debug(logger.RespFormat, "DescribeMongoDBInstance", queryReq, resp)

		if err != nil {
			if strings.Contains(err.Error(), "InstanceNotFound") {
				return nil
			} else {
				return resource.NonRetryableError(err)
			}
		}

		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceMongodbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)
	defer d.Partial(false)
	conn := meta.(*KsyunClient).mongodbconn
	// rename
	if d.HasChange("name") {
		d.SetPartial("name")
		v, ok := d.GetOk("name")
		if !ok {
			return fmt.Errorf("cann't change name to empty string")
		}
		rename := make(map[string]interface{})
		rename["InstanceId"] = d.Id()
		rename["Name"] = v.(string)
		logger.Debug(logger.ReqFormat, "RenameMongoDBInstance", rename)
		resp, err := conn.RenameMongoDBInstance(&rename)
		if err != nil {
			return fmt.Errorf("error on rename instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, "RenameMongoDBInstance", rename, *resp)
	}
	//resize
	if d.HasChange("node_num") {
		d.SetPartial("node_num")
		o, n := d.GetChange("node_num")
		if n.(int) <= o.(int) {
			return fmt.Errorf("only supports add node")
		}
		updateReq := make(map[string]interface{})
		updateReq["InstanceId"] = d.Id()
		updateReq["NodeNum"] = fmt.Sprintf("%v", n)
		logger.Debug(logger.ReqFormat, "AddSecondaryInstance", updateReq)

		conn := meta.(*KsyunClient).mongodbconn
		resp, err := conn.AddSecondaryInstance(&updateReq)
		if err != nil {
			return fmt.Errorf("error on add instance node: %s", err)
		}
		logger.Debug(logger.RespFormat, "AddSecondaryInstance", updateReq, *resp)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"resizing"},
			Target:     []string{"running"},
			Refresh:    mongodbInstanceStateRefreshForOperateFunc(conn, d.Id(), []string{"resizing", "running"}),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      20 * time.Second,
			MinTimeout: 1 * time.Minute,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error on add Instance node: %s", err)
		}
	}

	return resourceMongodbInstanceRead(d, meta)
}

func resourceMongodbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).mongodbconn

	readReq := make(map[string]interface{})
	readReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DescribeMongoDBInstance", readReq)
	resp, err := conn.DescribeMongoDBInstance(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, "DescribeMongoDBInstance", readReq, *resp)
	item, ok := (*resp)["MongoDBInstanceResult"].(map[string]interface{})

	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range item {
		if mongodbInstanceKeys[k] {
			if k == "IP" {
				result["ip"] = v
			} else {
				result[Hump2Downline(k)] = v
			}
		}
	}

	for k, v := range result {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("error set data %v :%v", v, err)
		}
	}

	return mongodbInstanceNodeRead(d, meta)
}

func mongodbInstanceStateRefreshForCreateFunc(client *mongodb.Mongodb, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		queryReq := map[string]interface{}{"InstanceId": instanceId}
		logger.Debug(logger.ReqFormat, "DescribeMongoDBInstance", queryReq)

		resp, err := client.DescribeMongoDBInstance(&queryReq)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, "DescribeMongoDBInstance", queryReq, *resp)

		item, ok := (*resp)["MongoDBInstanceResult"].(map[string]interface{})

		if !ok {
			return nil, "", fmt.Errorf("no instance information was queried. InstanceId:%s", instanceId)
		}
		status := item["Status"].(string)
		if status == "error" {
			return nil, "", fmt.Errorf("instance create error, status:%v", status)
		}

		for k, v := range target {
			if v == status {
				return resp, status, nil
			}
			if k == len(target)-1 {
				status = "creating"
			}
		}
		return resp, status, nil
	}
}

func mongodbInstanceStateRefreshForOperateFunc(client *mongodb.Mongodb, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		queryReq := map[string]interface{}{"InstanceId": instanceId}
		logger.Debug(logger.ReqFormat, "DescribeMongoDBInstance", queryReq)

		resp, err := client.DescribeMongoDBInstance(&queryReq)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, "DescribeMongoDBInstance", queryReq, *resp)

		item, ok := (*resp)["MongoDBInstanceResult"].(map[string]interface{})

		if !ok {
			return nil, "", fmt.Errorf("no instance information was queried. InstanceId:%s", instanceId)
		}
		status := item["Status"].(string)
		if status == "error" {
			return nil, "", fmt.Errorf("instance create error, status:%v", status)
		}

		for k, v := range target {
			if v == status {
				return resp, status, nil
			}
			if k == len(target)-1 {
				status = "resizing"
			}
		}
		return resp, status, nil
	}
}

func mongodbInstanceNodeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).mongodbconn

	readReq := make(map[string]interface{})
	readReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DescribeMongoDBInstanceNode", readReq)
	resp, err := conn.DescribeMongoDBInstanceNode(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading instance node %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, "DescribeMongoDBInstanceNode", readReq, *resp)
	itemSet, ok := (*resp)["MongoDBInstanceNodeResult"].([]interface{})

	if !ok || len(itemSet) == 0 {
		logger.Info("instance node result size : 0")
		return nil
	}

	nodes := GetSubSliceDByRep(itemSet, mongodbInstanceNodeKeys)
	for _, node := range nodes {
		node["ip"] = node["i_p"]
		delete(node, "i_p")
	}

	if err := d.Set("nodes", nodes); err != nil {
		return fmt.Errorf("error set data %v :%v", nodes, err)
	}

	return nil
}
