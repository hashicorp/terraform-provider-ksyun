package ksyun

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunLb() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLbCreate,
		Read:   resourceKsyunLbRead,
		Update: resourceKsyunLbUpdate,
		Delete: resourceKsyunLbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"load_balancer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunLbCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		req["VpcId"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		req["LoadBalancerName"] = fmt.Sprintf("%v", v)
	} else {
		req["LoadBalancerName"] = resource.PrefixedUniqueId("tf-lb-")
	}
	if v, ok := d.GetOk("type"); ok {
		req["Type"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		req["SubnetId"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		req["PrivateIpAddress"] = fmt.Sprintf("%v", v)
	}
	action := "CreateLoadBalancer"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := slbconn.CreateLoadBalancer(&req)
	if err != nil {
		return fmt.Errorf("Error CreateLoadBalancer : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	id, ok := (*resp)["LoadBalancerId"]
	if !ok {
		return fmt.Errorf("Error CreateLoadBalancer : no LoadBalancerId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("Error CreateLoadBalancer : no LoadBalancerId found")
	}
	if err := d.Set("load_balancer_id", idres); err != nil {
		return err
	}
	d.SetId(idres)

	return resourceKsyunLbRead(d, m)
}

func resourceKsyunLbRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerId.1"] = d.Id()
	action := "DescribeLoadBalancers"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeLoadBalancers(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeLoadBalancers : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	type describeLoadBalancersStruct struct {
		LoadBalancerDescriptions []struct {
			LoadBalancerId    string
			LoadBalancerName  string
			IsWaf             bool
			Type              string
			VpcId             string
			LoadBalancerState string
			CreateTime        string
			ListenersCount    int
			ProjectId         string
			State             string
			IpVersion         string
		}
	}
	by, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("Error DescribeLoadBalancers when marshal : %s", err)
	}
	var result describeLoadBalancersStruct
	err = json.Unmarshal(by, &result)
	if err != nil {
		return fmt.Errorf("Error DescribeLoadBalancers when unmarshal: %s", err)
	}
	if len(result.LoadBalancerDescriptions) == 0 {
		d.SetId("")
		return nil
	}
	lb0 := result.LoadBalancerDescriptions[0]
	if err := d.Set("load_balancer_id", lb0.LoadBalancerId); err != nil {
		return err
	}
	if err := d.Set("load_balancer_name", lb0.LoadBalancerName); err != nil {
		return err
	}
	if err := d.Set("is_waf", lb0.IsWaf); err != nil {
		return err
	}
	if err := d.Set("type", lb0.Type); err != nil {
		return err
	}
	if err := d.Set("vpc_id", lb0.VpcId); err != nil {
		return err
	}
	if err := d.Set("load_balancer_state", lb0.LoadBalancerState); err != nil {
		return err
	}
	if err := d.Set("create_time", lb0.CreateTime); err != nil {
		return err
	}
	if err := d.Set("project_id", lb0.ProjectId); err != nil {
		return err
	}
	if err := d.Set("state", lb0.State); err != nil {
		return err
	}
	if err := d.Set("ip_version", lb0.IpVersion); err != nil {
		return err
	}
	return nil
}

func resourceKsyunLbUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerId"] = d.Id()
	if v, ok := d.GetOk("load_balancer_name"); ok {
		req["LoadBalancerName"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("load_balancer_state"); ok {
		req["LoadBalancerState"] = fmt.Sprintf("%v", v)
	} else {
		return fmt.Errorf("cann't change load_balancer_state to empty string")
	}
	// Enable partial attribute modification
	d.Partial(true)
	// Whether the representative has any modifications
	attributeUpdate := false
	if d.HasChange("load_balancer_name") {
		attributeUpdate = true
	}
	if d.HasChange("load_balancer_state") {
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifyLoadBalancer"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.ModifyLoadBalancer(&req)
		if err != nil {
			logger.Debug(logger.AllFormat, action+" first", req, *resp, err)
			if strings.Contains(err.Error(), "400") {
				time.Sleep(time.Second * 2)
				resp, err = slbconn.ModifyLoadBalancer(&req)
				if err != nil {
					return fmt.Errorf("update Slb (%v)error twice:%v", req, err)
				}
			}
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		d.SetPartial("load_balancer_name")
		d.SetPartial("load_balancer_state")
	}
	d.Partial(false)
	return resourceKsyunLbRead(d, m)
}

func resourceKsyunLbDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerId"] = d.Id()
	/*
		_, err := slbconn.DeleteLoadBalancer(&req)
		if err != nil {
			return fmt.Errorf("release Slb error:%v", err)
		}
		return nil
	*/
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteLoadBalancer"
		logger.Debug(logger.ReqFormat, action, req)

		resp, err1 := slbconn.DeleteLoadBalancer(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["LoadBalancerId.1"] = d.Id()
		action = "DescribeLoadBalancers"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeLoadBalancers(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading lb when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)

		itemSet, ok := (*resp)["LoadBalancerDescriptions"]
		if !ok {
			return nil
		}
		items, ok := itemSet.([]interface{})
		if !ok || len(items) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified lb %q has not been deleted due to unknown error", d.Id()))
	})
}
