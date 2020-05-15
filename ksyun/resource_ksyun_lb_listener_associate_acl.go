package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunListenerLBAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerLBAclCreate,
		Read:   resourceKsyunListenerLBAclRead,
		Delete: resourceKsyunListenerLBAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func resourceKsyunListenerLBAclCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"load_balancer_acl_id",
		"listener_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}

	action := "AssociateLoadBalancerAcl"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.AssociateLoadBalancerAcl(&req)
	if err != nil {
		return fmt.Errorf("Error CreateListenerLBAcls : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	status, ok := (*resp)["Return"]
	if !ok {
		return fmt.Errorf("Error CreateListenerLBAcls")
	}
	statu, ok := status.(bool)
	if !ok {
		return fmt.Errorf("Error CreateListenerLBAcls ")
	}
	if !statu {
		return fmt.Errorf("Error CreateListenerLBAcls : fail")
	}
	id := fmt.Sprintf("%s:%s", d.Get("listener_id"), d.Get("load_balancer_acl_id"))
	d.SetId(id)
	return resourceKsyunListenerLBAclRead(d, m)
}

func resourceKsyunListenerLBAclRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	p := strings.Split(d.Id(), ":")
	req["ListenerId.1"] = p[0]
	action := "DescribeListeners"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeListeners(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeListeners : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	itemset := (*resp)["ListenerSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	associatekeys := map[string]bool{
		"LoadBalancerAclId": true,
		"ListenerId":        true,
	}
	SetDByResp(d, items[0], associatekeys, map[string]bool{})
	return nil
}

func resourceKsyunListenerLBAclDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	p := strings.Split(d.Id(), ":")
	req := make(map[string]interface{})
	req["ListenerId"] = p[0]
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DisassociateLoadBalancerAcl"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DisassociateLoadBalancerAcl(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		return resource.NonRetryableError(fmt.Errorf("DisassociateLoadBalancerAcl error:%v", err1))
	})
}
