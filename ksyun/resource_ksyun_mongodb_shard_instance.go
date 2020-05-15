package ksyun

import (
	"errors"
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/mongodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunMongodbShardInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbShardInstanceCreate,
		Delete: resourceMongodbShardInstanceDelete,
		Update: resourceMongodbShardInstanceUpdate,
		Read:   resourceMongodbShardInstanceRead,
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
			"mongos_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"shard_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"mongos_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"shard_class": {
				Type:     schema.TypeString,
				Required: true,
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
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
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
			"instance_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_storage": {
				Type:     schema.TypeInt,
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
			"mongos_nodes": {
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
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"shard_nodes": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongodbShardInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	createReq := make(map[string]interface{})
	creates := []string{
		"name",
		"instance_account",
		"instance_password",
		"mongos_num",
		"shard_num",
		"shard_class",
		"mongos_class",
		"storage",
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
	resp, err := conn.CreateMongoDBShardInstance(&createReq)
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
		Refresh:    mongodbShardInstanceStateRefreshForCreateFunc(conn, d.Id(), []string{"creating", "running"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	return resourceMongodbShardInstanceRead(d, meta)
}

func resourceMongodbShardInstanceDelete(d *schema.ResourceData, meta interface{}) error {

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

func resourceMongodbShardInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	logger.Info("no support for shard instance update functions")
	return errors.New("no support for shard instance update functions")
}

func resourceMongodbShardInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).mongodbconn

	queryReq := make(map[string]interface{})
	queryReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DescribeMongoDBInstance", queryReq)
	resp, err := conn.DescribeMongoDBInstance(&queryReq)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, "DescribeMongoDBInstance", queryReq, *resp)
	item, ok := (*resp)["MongoDBInstanceResult"].(map[string]interface{})

	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range item {
		if mongodbInstanceKeys[k] {
			if k == "IP" {
				result["ip"] = v
			} else if k == "Storage" {
				result["total_storage"] = v
				delete(result, "Storage")
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

	return mongodbShardInstanceNodeRead(d, meta)
}

func mongodbShardInstanceStateRefreshForCreateFunc(client *mongodb.Mongodb, instanceId string, target []string) resource.StateRefreshFunc {
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

func mongodbShardInstanceNodeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).mongodbconn

	readReq := make(map[string]interface{})
	readReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DescribeMongoDBShardNode", readReq)
	resp, err := conn.DescribeMongoDBShardNode(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading shard instance node %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, "DescribeMongoDBShardNode", readReq, *resp)

	items, ok := (*resp)["MongosNodeResult"].([]interface{})
	if !ok || len(items) == 0 {
		logger.Info("shard instance mongos node result size : 0")
	} else {
		mongos := GetSubSliceDByRep(items, mongodbShardInstanceMongosNodeKeys)
		if err := d.Set("mongos_nodes", mongos); err != nil {
			return fmt.Errorf("error set data %v :%v", mongos, err)
		}
	}

	items, ok = (*resp)["ShardNodeResult"].([]interface{})
	if !ok || len(items) == 0 {
		logger.Info("shard instance shard node result size : 0")
	} else {
		shards := GetSubSliceDByRep(items, mongodbShardInstanceShardNodeKeys)
		if err := d.Set("shard_nodes", shards); err != nil {
			return fmt.Errorf("error set data %v :%v", shards, err)
		}
	}

	return nil
}
