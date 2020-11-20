package ksyun

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// instance
func resourceRedisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisInstanceCreate,
		Delete: resourceRedisInstanceDelete,
		Update: resourceRedisInstanceUpdate,
		Read:   resourceRedisInstanceRead,
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
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
					v, ok := i.(int)
					if !ok {
						es = append(es, fmt.Errorf("expected type of %s to be int", k))
						return
					}

					if v != 1 && v != 2 {
						es = append(es, fmt.Errorf("expected %s to be in 1 or 2, got %d", k, v))
						return
					}

					return
				},
			},
			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"slave_num": {
				Type:     schema.TypeInt,
				Required: true,
				//Optional: true,
				//Computed: true,
				//ValidateFunc: validation.IntBetween(0,8),
			},
			"net_type": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
					v, ok := i.(int)
					if !ok {
						es = append(es, fmt.Errorf("expected type of %s to be int", k))
						return
					}

					if v != 1 && v != 2 {
						es = append(es, fmt.Errorf("expected %s to be in 1 or 2, got %d", k, v))
						return
					}

					return
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bill_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pass_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cache_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
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
			"used_memory": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"sub_order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_use": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reset_all_parameters": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	// create instance
	var (
		resp *map[string]interface{}
		err  error
		az   string
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	createParam := []string{"available_zone", "name", "capacity", "net_type", "vpc_id", "vnet_id", "mode", "slave_num", "bill_type", "duration", "duration_unit", "pass_word", "iam_project_id", "protocol", "security_group_id"}
	for _, v := range createParam {
		if v1, ok := d.GetOk(v); ok {
			createReq[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateCacheCluster"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.CreateCacheCluster(&createReq); err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["CacheId"].(string))
	}
	if createReq["AvailableZone"] != nil {
		az = createReq["AvailableZone"].(string)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"2"},
		Refresh:    stateRefreshForCreateFunc(conn, az, d.Id(), []string{"2"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	_ = resourceRedisInstanceRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	// create instance parameter
	return resourceRedisInstanceParamCreate(d, meta)
}

func resourceRedisInstanceParamCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp  *map[string]interface{}
		param map[string]interface{}
		az    string
		err   error
	)
	v1, ok := d.GetOk("parameters")
	if !ok || v1 == nil {
		return nil
	}

	if v2, o := d.GetOk("reset_all_parameters"); o && !v2.(bool) {
		conn := meta.(*KsyunClient).kcsv1conn
		createReq := make(map[string]interface{})
		if az, ok := d.GetOk("available_zone"); ok {
			createReq["AvailableZone"] = az
		}
		createReq["CacheId"] = d.Get("cache_id")
		createReq["Protocol"] = d.Get("protocol")
		createReq["ResetAllParameters"] = fmt.Sprintf("%v", d.Get("reset_all_parameters"))

		if param, ok = v1.(map[string]interface{}); !ok {
			return fmt.Errorf("expected type of parameter to not be map")
		}
		if len(param) == 0 {
			return nil
		}
		if err = validParam(d); err != nil {
			return err
		}
		params := v1.(map[string]interface{})
		var i int
		for k, v := range params {
			i = i + 1
			createReq[fmt.Sprintf("%v%v", "Parameters.ParameterName.", i)] = fmt.Sprintf("%v", k)
			createReq[fmt.Sprintf("%v%v", "Parameters.ParameterValue.", i)] = fmt.Sprintf("%v", v)
		}
		action := "SetCacheParameters"
		logger.Debug(logger.ReqFormat, action, createReq)
		if resp, err = conn.SetCacheParameters(&createReq); err != nil {
			return fmt.Errorf("error on set instance parameter: %s", err)
		}
		logger.Debug(logger.RespFormat, action, createReq, *resp)
		if createReq["AvailableZone"] != nil {
			az = createReq["AvailableZone"].(string)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"2"},
			Refresh:    stateRefreshForOperateFunc(conn, az, d.Id(), []string{"2"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 1 * time.Minute,
		}
		_, err = stateConf.WaitForState()
		_ = resourceRedisInstanceParamRead(d, meta)
		if err != nil {
			return fmt.Errorf("error on set instance parameter: %s", err)
		}
	}
	return nil
}

func resourceRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	conn := meta.(*KsyunClient).kcsv1conn
	deleteReq := make(map[string]interface{})
	deleteReq["CacheId"] = d.Id()
	if az, ok := d.GetOk("available_zone"); ok {
		deleteReq["AvailableZone"] = az
	}
	action := "DeleteCacheCluster"
	logger.Debug(logger.ReqFormat, action, deleteReq)
	if resp, err = conn.DeleteCacheCluster(&deleteReq); err != nil {
		return fmt.Errorf("error on deleting instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		var (
			resp *map[string]interface{}
			err  error
		)

		queryReq := make(map[string]interface{})
		queryReq["CacheId"] = d.Id()
		if az, ok := d.GetOk("available_zone"); ok {
			queryReq["AvailableZone"] = az
		}
		action := "DescribeCacheCluster"
		logger.Debug(logger.ReqFormat, action, queryReq)
		if resp, err = conn.DescribeCacheCluster(&queryReq); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "cannot be found") {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		logger.Debug(logger.RespFormat, action, queryReq, *resp)
		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp   *map[string]interface{}
		v      interface{}
		params interface{}
		err    error
		ok     bool
		az     string
	)

	d.Partial(true)
	defer d.Partial(false)
	conn := meta.(*KsyunClient).kcsv1conn
	// rename
	if d.HasChange("name") {
		d.SetPartial("name")
		if v, ok = d.GetOk("name"); !ok {
			return fmt.Errorf("cann't change name to empty string")
		}
		rename := make(map[string]interface{})
		rename["CacheId"] = d.Id()
		rename["Name"] = v.(string)
		if az, ok := d.GetOk("available_zone"); ok {
			rename["AvailableZone"] = az
		}
		action := "RenameCacheCluster"
		logger.Debug(logger.ReqFormat, action, rename)
		if resp, err = conn.RenameCacheCluster(&rename); err != nil {
			return fmt.Errorf("error on rename instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, rename, *resp)
	}
	// update password
	if d.HasChange("pass_word") {
		d.SetPartial("pass_word")
		if v, ok = d.GetOk("pass_word"); !ok {
			return fmt.Errorf("cann't change password to empty string")
		}
		password := make(map[string]interface{})
		password["CacheId"] = d.Id()
		password["Password"] = v.(string)
		password["Mode"] = d.Get("mode")
		if az, ok := d.GetOk("available_zone"); ok {
			password["AvailableZone"] = az
		}
		action := "UpdatePassword"
		logger.Debug(logger.ReqFormat, action, password)
		if resp, err = conn.UpdatePassword(&password); err != nil {
			return fmt.Errorf("error on update instance password %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, password, *resp)
	}
	// resize mem
	if d.HasChange("capacity") {
		d.SetPartial("capacity")
		if v, ok = d.GetOk("capacity"); !ok {
			return fmt.Errorf("cann't resize capacity to empty string")
		}
		resize := make(map[string]interface{})
		resize["CacheId"] = d.Id()
		resize["Capacity"] = fmt.Sprintf("%v", v.(int))
		if az, ok := d.GetOk("available_zone"); ok {
			resize["AvailableZone"] = az
		}
		action := "ResizeCacheCluster"
		logger.Debug(logger.ReqFormat, action, resize)
		if resp, err = conn.ResizeCacheCluster(&resize); err != nil {
			return fmt.Errorf("error on resize instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, resize, *resp)
		if resize["AvailableZone"] != nil {
			az = resize["AvailableZone"].(string)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"2"},
			Refresh:    stateRefreshForOperateFunc(conn, az, d.Id(), []string{"2"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      20 * time.Second,
			MinTimeout: 1 * time.Minute,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error on update Instance: %s", err)
		}
	}
	_ = resourceRedisInstanceRead(d, meta)

	// update parameter
	// parameters no change
	if !d.HasChange("reset_all_parameters") && !d.HasChange("parameters") {
		logger.Info("instance parameters have not changed")
		return nil
	}

	updateReq := make(map[string]interface{})
	updateReq["CacheId"] = d.Get("cache_id")
	updateReq["protocol"] = d.Get("protocol")
	if az, ok := d.GetOk("available_zone"); ok {
		updateReq["AvailableZone"] = az
	}

	reset := d.Get("reset_all_parameters")
	updateReq["ResetAllParameters"] = fmt.Sprintf("%v", reset)
	if !reset.(bool) {
		if params, ok = d.GetOk("parameters"); !ok {
			logger.Info("instance parameters do not exist")
			return nil
		}
		param, ok1 := params.(map[string]interface{})
		if !ok1 {
			logger.Info("type of instance parameters must be map")
			return nil
		}
		if len(param) == 0 {
			logger.Info("instance parameters size : 0")
			return nil
		}
		if err := validParam(d); err != nil {
			return err
		}
		var i int
		for k, v := range param {
			i = i + 1
			updateReq[fmt.Sprintf("%v%v", "Parameters.ParameterName.", i)] = fmt.Sprintf("%v", k)
			updateReq[fmt.Sprintf("%v%v", "Parameters.ParameterValue.", i)] = fmt.Sprintf("%v", v)
		}
	}

	action := "SetCacheParameters"
	logger.Debug(logger.ReqFormat, action, updateReq)
	if resp, err = conn.SetCacheParameters(&updateReq); err != nil {
		return fmt.Errorf("error on set instance parameter: %s", err)
	}
	logger.Debug(logger.RespFormat, action, updateReq, *resp)
	if updateReq["AvailableZone"] != nil {
		az = updateReq["AvailableZone"].(string)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"2"},
		Refresh:    stateRefreshForOperateFunc(conn, az, d.Id(), []string{"2"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	_ = resourceRedisInstanceParamRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on set instance parameter: %s", err)
	}
	return nil
}

func resourceRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		item map[string]interface{}
		resp *map[string]interface{}
		ok   bool
		err  error
	)

	conn := meta.(*KsyunClient).kcsv1conn
	queryReq := make(map[string]interface{})
	queryReq["CacheId"] = d.Id()
	if az, ok := d.GetOk("available_zone"); ok {
		queryReq["AvailableZone"] = az
	}
	action := "DescribeCacheCluster"
	logger.Debug(logger.ReqFormat, action, queryReq)
	if resp, err = conn.DescribeCacheCluster(&queryReq); err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, queryReq, *resp)
	if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range item {
		if k == "protocol" || k == "slaveNum" || !redisInstanceKeys[k] {
			continue
		}
		result[Hump2Downline(k)] = v
	}
	for k, v := range result {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("error set data %v :%v", v, err)
		}
	}

	//resourceRedisInstanceParamRead(d, meta)
	return nil
}

func resourceRedisInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Id()
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	action := "DescribeCacheParameters"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeCacheParameters(&readReq); err != nil {
		return fmt.Errorf("error on reading instance parameter %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		logger.Info("instance parameters result size : 0")
		return nil
	}
	result := make(map[string]interface{})
	for _, d := range data {
		param := d.(map[string]interface{})
		result[param["name"].(string)] = fmt.Sprintf("%v", param["currentValue"])
	}
	if err := d.Set("parameters", result); err != nil {
		return fmt.Errorf("error set data %v :%v", result, err)
	}
	return nil
}

func stateRefreshForCreateFunc(client *kcsv1.Kcsv1, az, instanceId string, target []string) resource.StateRefreshFunc {
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
		serviceStatus := int(item["serviceStatus"].(float64))
		// instance status error
		if status == 0 || status == 99 {
			return nil, "", fmt.Errorf("instance create error,status:%v", status)
		}
		// trade instance status error
		if serviceStatus == 3 {
			return nil, "", fmt.Errorf("instance create error,serviceStatus:%v", serviceStatus)
		}
		state := strconv.Itoa(status)
		for k, v := range target {
			if v == state && serviceStatus == 2 {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func stateRefreshForOperateFunc(client *kcsv1.Kcsv1, az, instanceId string, target []string) resource.StateRefreshFunc {
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
			if v == state {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}

		return resp, state, nil
	}
}

func validParam(d *schema.ResourceData) error {
	var (
		param map[string]interface{}
		ok    bool
	)
	protocol := d.Get("protocol")
	v1, ok := d.GetOk("parameters")
	if !ok || v1 == nil {
		return fmt.Errorf("parameter to not be null")
	}
	if param, ok = v1.(map[string]interface{}); !ok {
		return fmt.Errorf("expected type of parameter to not be map")
	}
	var filter map[string]*ValidatorParam
	if protocol == "4.0" || protocol == "5.0" {
		filter = GetValidatorParamForProto4()
	} else {
		filter = GetValidatorParamForProto()
	}
	for key, value := range param {
		result := filter[key]
		if result == nil {
			return fmt.Errorf("expected %s to be in range %v", key, filter)
		}
		if result.Valid.Type == "enum" {
			var ok bool
			for _, val := range result.Valid.Values {
				if val == value.(string) {
					ok = true
				}
			}
			if !ok {
				return fmt.Errorf("expected value of %s to in %v", key, result.Valid.Values)
			}
		} else if result.Valid.Type == "range" {
			v, ok := value.(string)
			if !ok {
				return fmt.Errorf("expected type of %v to be int", key)
			}
			val, _ := strconv.ParseInt(v, 10, 64)
			if val < result.Valid.Min || val > result.Valid.Max {
				return fmt.Errorf("expected value of %s to be in range [%v-%v]", key, result.Valid.Min, result.Valid.Max)
			}
		} else {
			if ok := regexp.MustCompile(result.Valid.Value).MatchString(value.(string)); !ok {
				return fmt.Errorf("expected value of %s to match regular expression %v", key, result.Valid.Value)
			}
		}
	}
	return nil
}

type Validity struct {
	Type     string
	DataType string
	Value    string
	Values   []string
	Min      int64
	Max      int64
}

type ValidatorParam struct {
	Name  string
	Desc  string
	Valid *Validity
}

var filterForProto4 map[string]*ValidatorParam
var onceFor4 sync.Once

func GetValidatorParamForProto4() map[string]*ValidatorParam {
	onceFor4.Do(func() {
		filterForProto4 = make(map[string]*ValidatorParam)
		filterForProto4["appendonly"] = &ValidatorParam{
			Name: "appendonly",
			Desc: "是否开启AOF持久化功能",
			Valid: &Validity{
				Type:     "enum",
				DataType: "string",
				Value:    "",
				Values:   []string{"yes", "no"},
				Min:      0,
				Max:      0,
			},
		}
		filterForProto4["appendfsync"] = &ValidatorParam{
			Name: "appendfsync",
			Desc: "AOF文件同步方式",
			Valid: &Validity{
				Type:     "enum",
				DataType: "string",
				Value:    "",
				Values:   []string{"everysec", "always", "no"},
				Min:      0,
				Max:      0,
			},
		}
		filterForProto4["maxmemory-policy"] = &ValidatorParam{
			Name: "maxmemory-policy",
			Desc: "Redis淘汰策略",
			Valid: &Validity{
				Type:     "enum",
				DataType: "string",
				Value:    "",
				Values:   []string{"volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl", "allkeys-lru", "allkeys-lfu", "allkeys-random", "noeviction"},
				Min:      0,
				Max:      0,
			},
		}
		filterForProto4["maxmemory-samples"] = &ValidatorParam{
			Name: "maxmemory-samples",
			Desc: "淘汰算法运行时的采样数",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      1,
				Max:      10,
			},
		}
		filterForProto4["hash-max-ziplist-entries"] = &ValidatorParam{
			Name: "hash-max-ziplist-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto4["hash-max-ziplist-value"] = &ValidatorParam{
			Name: "hash-max-ziplist-value",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto4["list-max-ziplist-size"] = &ValidatorParam{
			Name: "list-max-ziplist-size",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      -2147483648,
				Max:      2147483647,
			},
		}
		filterForProto4["set-max-intset-entries"] = &ValidatorParam{
			Name: "set-max-intset-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto4["zset-max-ziplist-entries"] = &ValidatorParam{
			Name: "zset-max-ziplist-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto4["zset-max-ziplist-value"] = &ValidatorParam{
			Name: "zset-max-ziplist-value",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto4["notify-keyspace-events"] = &ValidatorParam{
			Name: "notify-keyspace-events",
			Desc: "键空间通知配置",
			Valid: &Validity{
				Type:     "regexp",
				DataType: "string",
				Value:    "[KEg$lshzxeA]*",
				Values:   nil,
				Min:      0,
				Max:      0,
			},
		}
		filterForProto4["timeout"] = &ValidatorParam{
			Name: "timeout",
			Desc: "连接空闲超时时间",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      86400,
			},
		}
	})

	return filterForProto4
}

var filterForProto map[string]*ValidatorParam
var once sync.Once

func GetValidatorParamForProto() map[string]*ValidatorParam {
	once.Do(func() {
		filterForProto = make(map[string]*ValidatorParam)
		filterForProto["appendonly"] = &ValidatorParam{
			Name: "appendonly",
			Desc: "是否开启AOF持久化功能",
			Valid: &Validity{
				Type:     "enum",
				DataType: "string",
				Value:    "",
				Values:   []string{"yes", "no"},
				Min:      0,
				Max:      0,
			},
		}
		filterForProto["appendfsync"] = &ValidatorParam{
			Name: "appendfsync",
			Desc: "AOF文件同步方式",
			Valid: &Validity{
				Type:     "enum",
				DataType: "string",
				Value:    "",
				Values:   []string{"everysec", "always", "no"},
				Min:      0,
				Max:      0,
			},
		}
		filterForProto["maxmemory-policy"] = &ValidatorParam{
			Name: "maxmemory-policy",
			Desc: "Redis淘汰策略",
			Valid: &Validity{
				Type:     "enum",
				DataType: "string",
				Value:    "",
				Values:   []string{"volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl", "allkeys-lru", "allkeys-lfu", "allkeys-random", "noeviction"},
				Min:      0,
				Max:      0,
			},
		}
		filterForProto["maxmemory-samples"] = &ValidatorParam{
			Name: "maxmemory-samples",
			Desc: "淘汰算法运行时的采样数",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      1,
				Max:      10,
			},
		}
		filterForProto["hash-max-ziplist-entries"] = &ValidatorParam{
			Name: "hash-max-ziplist-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto["hash-max-ziplist-value"] = &ValidatorParam{
			Name: "hash-max-ziplist-value",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto["list-max-ziplist-entries"] = &ValidatorParam{
			Name: "list-max-ziplist-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto["list-max-ziplist-value"] = &ValidatorParam{
			Name: "list-max-ziplist-value",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}

		filterForProto["set-max-intset-entries"] = &ValidatorParam{
			Name: "set-max-intset-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto["zset-max-ziplist-entries"] = &ValidatorParam{
			Name: "zset-max-ziplist-entries",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto["zset-max-ziplist-value"] = &ValidatorParam{
			Name: "zset-max-ziplist-value",
			Desc: "内部数据结构优化的阈值",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      9007199254740990,
			},
		}
		filterForProto["notify-keyspace-events"] = &ValidatorParam{
			Name: "notify-keyspace-events",
			Desc: "键空间通知配置",
			Valid: &Validity{
				Type:     "regexp",
				DataType: "string",
				Value:    "[KEg$lshzxeA]*",
				Values:   nil,
				Min:      0,
				Max:      0,
			},
		}
		filterForProto["timeout"] = &ValidatorParam{
			Name: "timeout",
			Desc: "连接空闲超时时间",
			Valid: &Validity{
				Type:     "range",
				DataType: "integer",
				Value:    "",
				Values:   nil,
				Min:      0,
				Max:      86400,
			},
		}
	})

	return filterForProto
}
