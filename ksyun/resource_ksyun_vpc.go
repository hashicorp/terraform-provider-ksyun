package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVPCCreate,
		Update: resourceKsyunVPCUpdate,
		Read:   resourceKsyunVPCRead,
		Delete: resourceKsyunVPCDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"is_default": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceKsyunVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	var resp *map[string]interface{}
	var err error
	createVpc := make(map[string]interface{})

	creates := []string{
		"vpc_name",
		"cidr_block",
		"is_default",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createVpc[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateVpc"
	logger.Debug(logger.ReqFormat, action, createVpc)
	resp, err = conn.CreateVpc(&createVpc)
	logger.Debug(logger.AllFormat, action, createVpc, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating vpc, %s", err)
	}
	if resp != nil {
		vpc := (*resp)["Vpc"].(map[string]interface{})
		d.SetId(vpc["VpcId"].(string))
	}
	return resourceKsyunVPCRead(d, meta)
}

func resourceKsyunVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	readVpc := make(map[string]interface{})
	readVpc["VpcId.1"] = d.Id()
	action := "DescribeVpcs"
	logger.Debug(logger.ReqFormat, action, readVpc)
	resp, err := conn.DescribeVpcs(&readVpc)
	logger.Debug(logger.AllFormat, action, readVpc, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading vpc %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["VpcSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SetDByResp(d, items[0], vpcKeys, map[string]bool{})
	}
	return nil
}

func resourceKsyunVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	attributeUpdate := false
	d.Partial(true)
	modifyVpc := make(map[string]interface{})
	modifyVpc["VpcId"] = d.Id()

	if d.HasChange("vpc_name") && !d.IsNewResource() {
		modifyVpc["VpcName"] = fmt.Sprintf("%v", d.Get("vpc_name"))
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifyVpc"
		logger.Debug(logger.ReqFormat, action, modifyVpc)
		resp, err := conn.ModifyVpc(&modifyVpc)
		logger.Debug(logger.AllFormat, action, modifyVpc, *resp, err)
		if err != nil {
			return fmt.Errorf("error on updating vpc, %s", err)
		}
		d.SetPartial("vpc_name")
	}
	d.Partial(false)
	return resourceKsyunVPCRead(d, meta)
}

func resourceKsyunVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	deleteVpc := make(map[string]interface{})
	deleteVpc["VpcId"] = d.Id()
	action := "DeleteVpc"
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteVpc)
		resp, err1 := conn.DeleteVpc(&deleteVpc)
		logger.Debug(logger.AllFormat, action, deleteVpc, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		readVpc := make(map[string]interface{})
		readVpc["VpcId.1"] = d.Id()
		action = "DescribeVpcs"
		logger.Debug(logger.ReqFormat, action, deleteVpc)
		resp, err := conn.DescribeVpcs(&readVpc)
		logger.Debug(logger.AllFormat, action, readVpc, *resp)
		if err != nil && notFoundError(err) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on  reading VpcS when delete %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["VpcSet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting VpcS %q, %s", d.Id(), err1))
	})
}
