package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

// instance node
func resourceRedisInstanceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisInstanceNodeCreate,
		Delete: resourceRedisInstanceNodeDelete,
		Update: resourceRedisInstanceNodeUpdate,
		Read:   resourceRedisInstanceNodeRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pre_node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRedisInstanceNodeCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
		az   string
	)

	// create
	conn := meta.(*KsyunClient).kcsv2conn
	createNodeReq := make(map[string]interface{})
	createNodeReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		createNodeReq["AvailableZone"] = az
	}
	action := "AddCacheSlaveNode"
	logger.Debug(logger.ReqFormat, action, createNodeReq)
	if resp, err = conn.AddCacheSlaveNode(&createNodeReq); err != nil {
		return fmt.Errorf("error on add instance node: %s", err)
	}
	if resp != nil {
		d.Set("instance_id", (*resp)["Data"].(map[string]interface{})["NodeId"].(string))
	}
	d.SetId(d.Get("instance_id").(string))
	logger.Debug(logger.RespFormat, action, createNodeReq, *resp)

	if createNodeReq["AvailableZone"] != nil {
		az = createNodeReq["AvailableZone"].(string)
	}
	refreshConn := meta.(*KsyunClient).kcsv1conn
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"2"},
		Refresh:    stateRefreshForOperateNodeFunc(refreshConn, az, d.Get("cache_id").(string), []string{"2"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	resourceRedisInstanceNodeRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on add Instance node: %s", err)
	}
	return nil
}

func resourceRedisInstanceNodeDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
		az   string
	)

	// delete
	conn := meta.(*KsyunClient).kcsv2conn
	deleteParamReq := make(map[string]interface{})
	deleteParamReq["CacheId"] = d.Get("cache_id")
	deleteParamReq["NodeId"] = d.Get("instance_id")
	if az, ok := d.GetOk("available_zone"); ok {
		deleteParamReq["AvailableZone"] = az
	}
	action := "DeleteCacheSlaveNode"
	logger.Debug(logger.ReqFormat, action, deleteParamReq)
	if resp, err = conn.DeleteCacheSlaveNode(&deleteParamReq); err != nil {
		return fmt.Errorf("error on delete instance node: %s", err)
	}
	logger.Debug(logger.RespFormat, action, deleteParamReq, *resp)

	if deleteParamReq["AvailableZone"] != nil {
		az = deleteParamReq["AvailableZone"].(string)
	}
	refreshConn := meta.(*KsyunClient).kcsv1conn
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"2"},
		Refresh:    stateRefreshForOperateNodeFunc(refreshConn, az, d.Get("cache_id").(string), []string{"2"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error on delete instance node: %s", err)
	}
	return nil
}

func resourceRedisInstanceNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	logger.Info("no support for instance readonly node update functions")
	return nil
}

func resourceRedisInstanceNodeRead(d *schema.ResourceData, meta interface{}) error {
	var (
		item interface{}
		resp *map[string]interface{}
		ok   bool
		err  error
	)

	conn := meta.(*KsyunClient).kcsv2conn
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	action := "DescribeCacheReadonlyNode"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeCacheReadonlyNode(&readReq); err != nil {
		return fmt.Errorf("error on reading instance node %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	if item, ok = (*resp)["Data"]; !ok {
		return nil
	}
	items, ok := item.([]interface{})
	if !ok || len(items) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	nodeId := d.Get("instance_id").(string)
	for _, v := range items {
		vMap := v.(map[string]interface{})
		if nodeId == vMap["instanceId"] {
			result["instance_id"] = vMap["instanceId"]
			result["name"] = vMap["name"]
			result["port"] = fmt.Sprintf("%v", vMap["port"])
			result["ip"] = vMap["ip"]
			result["status"] = vMap["status"]
			result["create_time"] = vMap["createTime"]
			result["proxy"] = vMap["proxy"]
		}
	}
	for k, v := range result {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("error set data %v :%v", v, err)
		}
	}
	return nil
}

func stateRefreshForOperateNodeFunc(client *kcsv1.Kcsv1, az, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			resp *map[string]interface{}
			item map[string]interface{}
			ok   bool
			err  error
		)

		queryReq := map[string]interface{}{"CacheId": instanceId}
		queryReq["AvailableZone"] = az
		action := "DescribeCacheCluster"
		logger.Debug(logger.ReqFormat, action, queryReq)
		if resp, err = client.DescribeCacheCluster(&queryReq); err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, action, queryReq, *resp)
		if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
			return nil, "", fmt.Errorf("no instance information was queried.%s", "")
		}
		status := int(item["status"].(float64))
		if status == 0 || status == 99 {
			return nil, "", fmt.Errorf("instance operate error,status:%v", status)
		}
		state := strconv.Itoa(status)
		for k, v := range target {
			if v == state && int(item["transition"].(float64)) != 1 {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}

		return resp, state, nil
	}
}
