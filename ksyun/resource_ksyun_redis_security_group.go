package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// redis security group
func resourceRedisSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupCreate,
		Delete: resourceRedisSecurityGroupDelete,
		Update: resourceRedisSecurityGroupUpdate,
		Read:   resourceRedisSecurityGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRedisSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		createReq["AvailableZone"] = az
	}
	createReq["Name"] = d.Get("name")
	createReq["Description"] = d.Get("description")
	action := "CreateSecurityGroup"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.CreateSecurityGroup(&createReq); err != nil {
		return fmt.Errorf("error on create redis security group: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["securityGroupId"].(string))
	}
	return nil
}

func resourceRedisSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	// delete redis security group
	deleteReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		deleteReq["AvailableZone"] = az
	}
	deleteReq["SecurityGroupId.1"] = d.Id()
	action := "DeleteSecurityGroup"
	logger.Debug(logger.ReqFormat, action, deleteReq)
	if resp, err = conn.DeleteSecurityGroup(&deleteReq); err != nil {
		return fmt.Errorf("error on delete redis security group: %s", err)
	}
	logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	return nil
}

func resourceRedisSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	d.Partial(true)
	defer d.Partial(false)

	if !d.HasChange("name") && !d.HasChange("description") {
		return nil
	}

	updateReq := make(map[string]interface{})
	if az, ok := d.GetOk("available_zone"); ok {
		updateReq["AvailableZone"] = az
	}
	updateReq["SecurityGroupId"] = d.Id()
	updateReq["Name"] = d.Get("name")
	updateReq["Description"] = d.Get("description")
	conn := meta.(*KsyunClient).kcsv1conn
	action := "ModifySecurityGroup"
	logger.Debug(logger.ReqFormat, action, updateReq)
	if resp, err = conn.ModifySecurityGroup(&updateReq); err != nil {
		return fmt.Errorf("error on modify redis security group: %s", err)
	}
	logger.Debug(logger.RespFormat, action, updateReq, *resp)
	return nil
}

func resourceRedisSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeSecurityGroup"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeSecurityGroup(&readReq); err != nil {
		return fmt.Errorf("error on reading redis security group %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].(map[string]interface{})
	if len(data) == 0 {
		logger.Info("redis security group result size : 0")
		return nil
	}

	name := data["name"].(string)
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("error set data %v :%v", name, err)
	}

	description := data["description"].(string)
	if err := d.Set("description", description); err != nil {
		return fmt.Errorf("error set data %v :%v", description, err)
	}
	return nil
}
