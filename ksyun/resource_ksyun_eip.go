package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipCreate,
		Read:   resourceKsyunEipRead,
		Update: resourceKsyunEipUpdate,
		Delete: resourceKsyunEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"line_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PrePaidByMonth",
					"PostPaidByPeak",
					"PostPaidByDay",
					"PostPaidByTransfer",
					"PostPaidByHour",
					"HourlyInstantSettlement",
				}, false),
			},
			"purchase_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allocation_id": {
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
func resourceKsyunEipCreate(d *schema.ResourceData, m interface{}) error {
	eipConn := m.(*KsyunClient).eipconn
	createEip := make(map[string]interface{})
	creates := []string{
		"line_id",
		"band_width",
		"charge_type",
		"purchase_time",
		"project_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createEip[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "AllocateAddress"
	logger.Debug(logger.ReqFormat, action, createEip)
	resp, err := eipConn.AllocateAddress(&createEip)
	logger.Debug(logger.AllFormat, action, createEip, *resp, err)
	if err != nil {
		return fmt.Errorf("createEip Error  : %s", err)
	}
	id, ok := (*resp)["AllocationId"]
	if !ok {
		return fmt.Errorf("createEip Error  : no id found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("createEip Error : no id found")
	}
	if err := d.Set("allocation_id", idres); err != nil {
		return err
	}
	d.SetId(idres)
	return resourceKsyunEipRead(d, m)
}

func resourceKsyunEipRead(d *schema.ResourceData, m interface{}) error {
	eipConn := m.(*KsyunClient).eipconn
	readEip := make(map[string]interface{})
	readEip["AllocationId.1"] = d.Id()
	if pd, ok := d.GetOk("project_id"); ok {
		readEip["project_id"] = fmt.Sprintf("%v", pd)
	}
	action := "DescribeAddresses"
	logger.Debug(logger.ReqFormat, action, readEip)
	resp, err := eipConn.DescribeAddresses(&readEip)
	logger.Debug(logger.AllFormat, action, readEip, *resp, err)
	if err != nil {
		return fmt.Errorf("Error  : %s", err)
	}
	itemset, ok := (*resp)["AddressesSet"]
	if !ok {
		d.SetId("")
		return nil
	}
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], eipKeys, map[string]bool{})
	return nil
}

func resourceKsyunEipUpdate(d *schema.ResourceData, m interface{}) error {
	eipConn := m.(*KsyunClient).eipconn
	// Enable partial attribute modification
	d.Partial(true)
	// Whether the representative has any modifications
	attributeUpdate := false
	updateReq := make(map[string]interface{})
	updateReq["AllocationId"] = d.Id()
	// modify
	if d.HasChange("band_width") && !d.IsNewResource() {
		if v, ok := d.GetOk("band_width"); ok {
			updateReq["BandWidth"] = fmt.Sprintf("%v", v)
		} else {
			return fmt.Errorf("cann't change bandwidth to empty string")
		}
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifyAddress"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := eipConn.ModifyAddress(&updateReq)
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		if err != nil {
			return fmt.Errorf("update eip (%v)error:%v", updateReq, err)
		}
		d.SetPartial("band_width")
	}
	d.Partial(false)
	return resourceKsyunEipRead(d, m)
}

func resourceKsyunEipDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).eipconn
	//delete
	deleteEip := make(map[string]interface{})
	deleteEip["AllocationId"] = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "ReleaseAddress"
		logger.Debug(logger.ReqFormat, action, deleteEip)
		_, err1 := conn.ReleaseAddress(&deleteEip)
		logger.Debug(logger.AllFormat, action, deleteEip, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}

		//check
		readEip := make(map[string]interface{})
		readEip["AllocationId.1"] = d.Id()
		if pd, ok := d.GetOk("project_id"); ok {
			readEip["project_id"] = fmt.Sprintf("%v", pd)
		}
		action = "DescribeAddresses"
		logger.Debug(logger.ReqFormat, action, readEip)
		resp, err := conn.DescribeAddresses(&readEip)
		logger.Debug(logger.AllFormat, action, readEip, *resp, err)
		if err != nil && notFoundError(err) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading eip when deleting %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["AddressesSet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting eip %v,%v", d.Id(), err1))
	})

}
