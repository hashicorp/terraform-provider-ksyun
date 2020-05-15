package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSubnetCreate,
		Update: resourceKsyunSubnetUpdate,
		Read:   resourceKsyunSubnetRead,
		Delete: resourceKsyunSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"subnet_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"subnet_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSubnetType,
			},

			"dhcp_ip_to": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIpAddress,
			},

			"dhcp_ip_from": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIpAddress,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"gateway_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"dns1": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
			},

			"dns2": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
			},
			"network_acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"availability_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availble_i_p_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	var resp *map[string]interface{}
	var err error
	creates := []string{
		"availability_zone",
		"subnet_name",
		"cidr_block",
		"subnet_type",
		"dhcp_ip_from",
		"dhcp_ip_to",
		"gateway_ip",
		"vpc_id",
		"dns1",
		"dns2",
	}
	createSubnet := make(map[string]interface{})
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createSubnet[vv] = fmt.Sprintf("%v", v1)
		}
	}

	if d.Get("subnet_type") != "Reserve" && (d.Get("gateway_ip") == nil || d.Get("gateway_ip") == "") {
		return fmt.Errorf("subnet_type not Reserve,Must set gateway_ip")
	}
	action := "CreateSubnet"
	logger.Debug(logger.ReqFormat, action, createSubnet)
	resp, err = conn.CreateSubnet(&createSubnet)
	logger.Debug(logger.AllFormat, action, createSubnet, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Subnet, %s", err)
	}
	if resp != nil {
		Subnet := (*resp)["Subnet"].(map[string]interface{})
		d.SetId(Subnet["SubnetId"].(string))
	}
	return resourceKsyunSubnetRead(d, meta)
}

func resourceKsyunSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	readSubnet := make(map[string]interface{})
	readSubnet["SubnetId.1"] = d.Id()
	action := "DescribeSubnets"
	logger.Debug(logger.ReqFormat, action, readSubnet)
	resp, err := conn.DescribeSubnets(&readSubnet)
	logger.Debug(logger.AllFormat, action, readSubnet, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Subnet %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["SubnetSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SetDByResp(d, items[0], subnetKeys, map[string]bool{})
	}
	return nil
}

func resourceKsyunSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	d.Partial(true)
	attributeUpdate := false
	modifySubnet := make(map[string]interface{})
	modifySubnet["SubnetId"] = d.Id()

	if d.HasChange("subnet_name") && !d.IsNewResource() {
		modifySubnet["SubnetName"] = fmt.Sprintf("%v", d.Get("subnet_name"))
		attributeUpdate = true
	}
	if d.HasChange("dns1") && !d.IsNewResource() {
		modifySubnet["Dns1"] = fmt.Sprintf("%v", d.Get("dns1"))
		attributeUpdate = true
	}
	if d.HasChange("dns2") && !d.IsNewResource() {
		modifySubnet["Dns2"] = fmt.Sprintf("%v", d.Get("dns2").(string))
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifySubnet"
		logger.Debug(logger.ReqFormat, action, modifySubnet)
		resp, err := conn.ModifySubnet(&modifySubnet)
		logger.Debug(logger.AllFormat, action, modifySubnet, *resp, err)
		if err != nil {
			return fmt.Errorf("error on updating Subnet, %s", err)
		}
		d.SetPartial("subnet_name")
		d.SetPartial("dns1")
		d.SetPartial("dns2")
	}
	d.Partial(false)
	return resourceKsyunSubnetRead(d, meta)
}

func resourceKsyunSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	deleteSubnet := make(map[string]interface{})
	deleteSubnet["SubnetId"] = d.Id()
	action := "DeleteSubnet"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteSubnet)
		resp, err1 := conn.DeleteSubnet(&deleteSubnet)
		logger.Debug(logger.AllFormat, action, deleteSubnet, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		readSubnet := make(map[string]interface{})
		readSubnet["SubnetId.1"] = d.Id()
		action = "DescribeSubnets"
		logger.Debug(logger.ReqFormat, action, readSubnet)
		resp, err := conn.DescribeSubnets(&readSubnet)
		logger.Debug(logger.AllFormat, action, readSubnet, *resp, err)
		if err != nil && notFoundError(err1) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on  reading SubnetS when delete %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["SubnetSet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting SubnetS %q, %s", d.Id(), err1))
	})

}
