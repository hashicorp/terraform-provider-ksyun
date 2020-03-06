package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunLoadBalancerAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclEntryCreate,
		Delete: resourceKsyunLoadBalancerAclEntryDelete,
		Read:   resourceKsyunLoadBalancerAclEntryRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_acl_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rule_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}
func resourceKsyunLoadBalancerAclEntryRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return fmt.Errorf("error id:%v", d.Id())
	}
	req["LoadBalancerAclId.1"] = ids[0]
	action := "DescribeLoadBalancerAcls"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := slbconn.DescribeLoadBalancerAcls(&req)
	if err != nil {
		return fmt.Errorf(" read LoadBalancerAcls : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	resSet := (*resp)["LoadBalancerAclSet"]
	res, ok := resSet.([]interface{})
	if !ok || len(res) == 0 {
		d.SetId("")
		return nil
	}
	subPara, ok := res[0].(map[string]interface{})
	if !ok || len(subPara) == 0 {
		d.SetId("")
		return nil
	}
	lbes, ok := subPara["LoadBalancerAclEntrySet"].([]interface{})
	if !ok || len(lbes) == 0 {
		d.SetId("")
		return nil
	}
	for _, aclEntry := range lbes {
		aclEntryItem, ok := aclEntry.(map[string]interface{})
		if !ok || len(aclEntryItem) == 0 {
			d.SetId("")
			return nil
		}
		if aclEntryItem["LoadBalancerAclEntryId"] == ids[1] {
			for key, value := range aclEntryItem {
				if err := d.Set(Hump2Downline(key), value); err != nil {
					return err
				}
				return nil
			}
		}

	}
	return fmt.Errorf("no LoadBalancerAclEntrySet get")
}
func resourceKsyunLoadBalancerAclEntryCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"load_balancer_acl_id",
		"cidr_block",
		"rule_number",
		"rule_action",
		"protocol",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateLoadBalancerAclEntry"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.CreateLoadBalancerAclEntry(&req)
	if err != nil {
		return fmt.Errorf("create LoadBalancerAclEntry : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	loadBalancerAclEntry, ok := (*resp)["LoadBalancerAclEntry"]
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry found")
	}

	lbae, ok := loadBalancerAclEntry.(map[string]interface{})
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry data found")
	}

	id, ok := lbae["LoadBalancerAclEntryId"]
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry id found")
	}
	if err := d.Set("load_balancer_acl_entry_id", id); err != nil {
		return err
	}
	ids, ok := id.(string)
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry id found")
	}
	SetDByResp(d, lbae, lbAclEntryKeys, map[string]bool{})
	ids = fmt.Sprintf("%v:%v", d.Get("load_balancer_acl_id"), ids)
	d.SetId(ids)
	return nil
}

func resourceKsyunLoadBalancerAclEntryDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return fmt.Errorf("error id:%v", d.Id())
	}
	req := make(map[string]interface{})
	req["LoadBalancerAclEntryId"] = ids[1]
	req["LoadBalancerAclId"] = ids[0]
	/*
		_, err := slbconn.DeregisterInstancesFromListener(&req)
		if err != nil {
			return fmt.Errorf("delete LoadBalancerAclEntry error:%v", err)
		}
		return nil
	*/
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteLoadBalancerAclEntry"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteLoadBalancerAclEntry(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		return resource.NonRetryableError(fmt.Errorf("DeleteLoadBalancerAclEntry error:%v", err1))
	})
}
