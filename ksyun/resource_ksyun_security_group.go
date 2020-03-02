package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupCreate,
		Update: resourceKsyunSecurityGroupUpdate,
		Read:   resourceKsyunSecurityGroupRead,
		Delete: resourceKsyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_entry_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icmp_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"icmp_code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_from": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_to": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceKsyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	var resp *map[string]interface{}
	var err error
	createSecurityGroup := make(map[string]interface{})
	createSecurityGroup["VpcId"] = d.Get("vpc_id")
	createSecurityGroup["SecurityGroupName"] = d.Get("security_group_name")
	action := "CreateSecurityGroup"
	logger.Debug(logger.ReqFormat, action, createSecurityGroup)
	resp, err = conn.CreateSecurityGroup(&createSecurityGroup)
	logger.Debug(logger.AllFormat, action, createSecurityGroup, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating SecurityGroup, %s", err)
	}
	if resp != nil {
		securityGroup := (*resp)["SecurityGroup"].(map[string]interface{})
		securityGroupId := securityGroup["SecurityGroupId"].(string)
		d.SetId(securityGroupId)
	}
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn

	readSecurityGroup := make(map[string]interface{})
	readSecurityGroup["SecurityGroupId.1"] = d.Id()

	action := "DescribeSecurityGroups"
	logger.Debug(logger.ReqFormat, action, readSecurityGroup)
	resp, err := conn.DescribeSecurityGroups(&readSecurityGroup)
	logger.Debug(logger.AllFormat, action, readSecurityGroup, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading SecurityGroup %q, %s", d.Id(), err)
	}

	itemset := (*resp)["SecurityGroupSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	excludes := SetDByResp(d, items[0], vpcSecurityGroupKeys, map[string]bool{
		"SecurityGroupEntrySet": true,
	},
	)
	if excludes["SecurityGroupEntrySet"] == nil {
		return nil
	}
	securityGroupEntrySet := GetSubSliceDByRep(excludes["SecurityGroupEntrySet"].([]interface{}), vpcSecurityGroupEntrySetKeys)
	if err := d.Set("security_group_entry_set", securityGroupEntrySet); err != nil {
		return err
	}
	return nil
}

func resourceKsyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	d.Partial(true)
	attributeUpdate := false
	modifySecurityGroup := make(map[string]interface{})
	modifySecurityGroup["SecurityGroupId"] = d.Id()

	if d.HasChange("security_group_name") && !d.IsNewResource() {
		modifySecurityGroup["SecurityGroupName"] = fmt.Sprintf("%v", d.Get("security_group_name"))
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifySecurityGroup"
		logger.Debug(logger.ReqFormat, action, modifySecurityGroup)
		resp, err := conn.ModifySecurityGroup(&modifySecurityGroup)
		logger.Debug(logger.AllFormat, action, modifySecurityGroup, *resp, err)
		if err != nil {
			return fmt.Errorf("error on updating SecurityGroup, %s", err)
		}
		d.SetPartial("security_group_name")
	}
	d.Partial(false)
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	//delete
	deleteSecurityGroup := make(map[string]interface{})
	deleteSecurityGroup["SecurityGroupId"] = d.Id()
	action := "DeleteSecurityGroup"
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteSecurityGroup)
		resp, err := conn.DeleteSecurityGroup(&deleteSecurityGroup)
		logger.Debug(logger.AllFormat, action, deleteSecurityGroup, *resp, err)
		if err == nil || (err != nil && notFoundError(err)) {
			return nil
		}
		if err != nil && inUseError(err) {
			return resource.RetryableError(err)
		}

		//check
		readSecurityGroup := make(map[string]interface{})
		readSecurityGroup["SecurityGroupId.1"] = d.Id()
		action = "DescribeSecurityGroups"
		logger.Debug(logger.ReqFormat, action, readSecurityGroup)
		resp, err1 := conn.DescribeSecurityGroups(&readSecurityGroup)
		logger.Debug(logger.AllFormat, action, readSecurityGroup, *resp, err1)
		if err1 != nil && notFoundError(err1) {
			return nil
		}
		if err1 != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading SecurityGroups when deleting %q, %s", d.Id(), err))
		}
		items, ok := (*resp)["SecurityGroupSet"]
		if !ok {
			return nil
		}
		item, ok := items.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting SecurityGroups %q, %s", d.Id(), err))
	})

}
