package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
)

// instance List
func dataSourceKsyunMongodbs() *schema.Resource {
	return &schema.Resource{
		// Instance List Query Function
		Read: dataSourceMongodbInstancesRead,
		// Define input and output parameters
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"name": {
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
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeInt,
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
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnet_id": {
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
						"pay_type": {
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
						"iam_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iam_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
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
					},
				},
			},
		},
	}
}

func dataSourceMongodbInstancesRead(d *schema.ResourceData, meta interface{}) error {

	var (
		allInstances []interface{}
		limit        = 100
		nextToken    string
	)

	readReq := make(map[string]interface{})
	filters := []string{"iam_project_id", "instance_id", "vnet_id", "vpc_id", "name", "vip"}
	for _, v := range filters {
		if value, ok := d.GetOk(v); ok {
			readReq[Downline2Hump(v)] = fmt.Sprintf("%v", value)
		}
	}
	readReq["Limit"] = fmt.Sprintf("%v", limit)

	conn := meta.(*KsyunClient).mongodbconn

	for {
		if nextToken != "" {
			readReq["Offset"] = nextToken
		}
		logger.Debug(logger.ReqFormat, "DescribeMongoDBInstances", readReq)

		resp, err := conn.DescribeMongoDBInstances(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", readReq, err)
		}
		logger.Debug(logger.RespFormat, "DescribeMongoDBInstances", readReq, *resp)

		itemSet, ok := (*resp)["MongoDBInstancesResult"]
		if !ok {
			break
		}
		items, ok := itemSet.([]interface{})
		if !ok {
			break
		}
		if items == nil || len(items) < 1 {
			break
		}
		allInstances = append(allInstances, items...)
		if len(items) < limit {
			break
		}
		nextToken = strconv.Itoa(int((*resp)["limit"].(float64)) + int((*resp)["offset"].(float64)))
	}

	values := GetSubSliceDByRep(allInstances, mongodbInstanceKeys)
	for _, v := range values {
		v["ip"] = v["i_p"]
		delete(v, "i_p")
	}
	if err := dataSourceKscSave(d, "instances", []string{}, values); err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}

	return nil
}
