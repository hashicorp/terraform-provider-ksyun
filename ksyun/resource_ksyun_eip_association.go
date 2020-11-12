package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipAssociationCreate,
		Read:   resourceKsyunEipAssociationRead,
		Delete: resourceKsyunEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"line_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_band_width_share": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunEipAssociationCreate(d *schema.ResourceData, m interface{}) error {
	eipConn := m.(*KsyunClient).eipconn

	req := make(map[string]interface{})
	creates := []string{
		"allocation_id",
		"instance_type",
		"instance_id",
		"network_interface_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "AssociateAddress"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := eipConn.AssociateAddress(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("Error AssociateAddress : %s", err)
	}
	status, ok := (*resp)["Return"]
	if !ok {
		return fmt.Errorf("Error AssociateAddress ")
	}
	status1, ok := status.(bool)
	if !ok || !status1 {
		return fmt.Errorf("Error AssociateAddress:fail ")
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("allocation_id"), d.Get("instance_id")))
	return resourceKsyunEipAssociationRead(d, m)
}

func resourceKsyunEipAssociationRead(d *schema.ResourceData, m interface{}) error {
	eipConn := m.(*KsyunClient).eipconn
	p := strings.Split(d.Id(), ":")
	req := make(map[string]interface{})
	req["AllocationId.1"] = p[0]
	if pd, ok := d.GetOk("project_id"); ok {
		req["ProjectId.1"] = fmt.Sprintf("%v", pd)
	}
	action := "DescribeAddresses"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := eipConn.DescribeAddresses(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("Error describeAddresses : %s", err)
	}
	itemset, ok := (*resp)["AddressesSet"]
	if !ok {
		d.SetId("")
		return nil
	}
	items := itemset.([]interface{})
	if len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], eipKeys, map[string]bool{})
	return nil
}

func resourceKsyunEipAssociationDelete(d *schema.ResourceData, m interface{}) error {
	eipConn := m.(*KsyunClient).eipconn
	deleteReq := make(map[string]interface{})
	p := strings.Split(d.Id(), ":")
	deleteReq["AllocationId"] = p[0]
	action := "DisassociateAddress"
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err1 := eipConn.DisassociateAddress(&deleteReq)
		logger.Debug(logger.AllFormat, action, deleteReq, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["AllocationId.1"] = p[0]
		if pd, ok := d.GetOk("project_id"); ok {
			req["ProjectId.1"] = fmt.Sprintf("%v", pd)
		}
		action = "DescribeAddresses"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := eipConn.DescribeAddresses(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil && notFoundError(err1) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on  reading eip when delete %q, %s", d.Id(), err))
		}
		addressesSets, ok := (*resp)["AddressesSet"]
		if !ok {
			return nil
		}
		addsets, ok := addressesSets.([]interface{})
		if !ok || len(addsets) == 0 {
			return nil
		}
		addset, ok := addsets[0].(map[string]interface{})
		if !ok {
			return nil
		}
		if instanceId, ok := addset["InstanceId"]; ok {
			if instanceId == p[1] {
				return resource.NonRetryableError(fmt.Errorf("the specified DisassociateAddress %q has not been deleted", d.Id()))
			}
		}
		return nil
	})
}
