package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// redis security group allocate
func resourceRedisSecurityGroupAllocate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupAllocateCreate,
		Delete: resourceRedisSecurityGroupAllocateDelete,
		Read:   resourceRedisSecurityGroupAllocateRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cache_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				ForceNew: true,
			},
		},
	}
}

func resourceRedisSecurityGroupAllocateCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		createReq["AvailableZone"] = az
	}
	createReq["SecurityGroupId"] = d.Get("security_group_id")
	ids := SchemaSetToStringSlice(d.Get("cache_ids"))
	for i, id := range ids {
		createReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = id
	}
	action := "AllocateSecurityGroup"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.AllocateSecurityGroup(&createReq); err != nil {
		return fmt.Errorf("error on allocate redis security group: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	d.SetId(d.Get("security_group_id").(string))
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		return nil
	}

	return fmt.Errorf("error on allocate redis security group: %v", data)
}

func resourceRedisSecurityGroupAllocateDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		createReq["AvailableZone"] = az
	}
	createReq["SecurityGroupId"] = d.Get("security_group_id")
	ids := SchemaSetToStringSlice(d.Get("cache_ids"))
	for i, id := range ids {
		createReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = id
	}
	action := "DeallocateSecurityGroup"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.DeallocateSecurityGroup(&createReq); err != nil {
		return fmt.Errorf("error on deallocate redis security group: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		return nil
	}

	return fmt.Errorf("error on deallocate redis security group: %v", data)
}

func resourceRedisSecurityGroupAllocateRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	readReq["SecurityGroupId"] = d.Id()
	readReq["Limit"] = 1000
	readReq["FilterCache"] = true
	action := "DescribeInstances"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeInstances(&readReq); err != nil {
		return fmt.Errorf("error on reading redis security group allocate instances %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].(map[string]interface{})
	if len(data) == 0 {
		logger.Info("redis security group result size : 0")
		return nil
	}

	lists := data["list"].([]interface{})
	if len(lists) == 0 {
		logger.Info("redis security group rule result size : 0")
		return nil
	}

	var result []string
	for _, v := range lists {
		inst := v.(map[string]interface{})
		result = append(result, inst["id"].(string))
	}
	if err := d.Set("cache_ids", result); err != nil {
		return fmt.Errorf("error set data %v :%v", result, err)
	}
	return nil
}
