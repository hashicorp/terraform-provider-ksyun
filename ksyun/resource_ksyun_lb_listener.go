package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerCreate,
		Read:   resourceKsyunListenerRead,
		Update: resourceKsyunListenerUpdate,
		Delete: resourceKsyunListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"session_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"session_persistence_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"session": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_persistence_period": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"session_state": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				//			Set: resourceKscListenerSessionHash,
			},
			"real_server": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"real_server_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_server_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"real_server_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"register_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
func resourceKsyunListenerCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"load_balancer_id",
		"listener_state",
		"listener_name",
		"listener_protocol",
		"certificate_id",
		"listener_port",
		"method",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	if v, ok := d.GetOk("session"); ok {
		FlatternStruct(v, &req)
	}
	action := "CreateListeners"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.CreateListeners(&req)
	if err != nil {
		return fmt.Errorf("Error CreateListeners : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	id, ok := (*resp)["ListenerId"]
	if !ok {
		return fmt.Errorf("Error CreateListeners : no ListenerId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("Error CreateListeners : no ListenerId found")
	}
	d.SetId(idres)
	if err := d.Set("listener_id", idres); err != nil {
		return err
	}
	return resourceKsyunListenerRead(d, m)
}

func resourceKsyunListenerRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["ListenerId.1"] = d.Id()
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
	excludes := SetDByResp(d, items[0], listenerKeys, map[string]bool{
		"HealthCheck": true,
		"RealServer":  true,
		"Session":     true},
	)

	subSession := GetSubStructDByRep(excludes["Session"], map[string]bool{})
	if err := d.Set("session", []interface{}{subSession}); err != nil {
		return err
	}
	server, ok := excludes["RealServer"].([]interface{})
	if ok {
		subRes := GetSubSliceDByRep(server, serverKeys)
		if err := d.Set("real_server", subRes); err != nil {
			return err
		}
	}
	subHealth := GetSubStructDByRep(excludes["HealthCheck"], map[string]bool{})
	if err := d.Set("health_check", []interface{}{subHealth}); err != nil {
		return err
	}
	return nil
}

func resourceKsyunListenerUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["ListenerId"] = d.Id()
	allAttributes := []string{
		"certificate_id",
		"listener_name",
		"listener_state",
		"method",
		/*
			"session_state",
			"session_persistence_period",
			"cookie_type",
			"cookie_name",
		*/
	}
	// Whether the representative has any modifications
	attributeUpdate := false
	var updates []string
	//Get the property that needs to be modified
	for _, v := range allAttributes {
		if d.HasChange(v) {
			attributeUpdate = true
			updates = append(updates, v)
		}
	}
	if !attributeUpdate {
		return nil
	}
	//Create a modification request
	for _, v := range allAttributes {
		if v1, ok := d.GetOk(v); ok {
			req[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
		}
	}
	// Enable partial attribute modification
	d.Partial(true)
	action := "ModifyListeners"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.ModifyListeners(&req)
	if err != nil {
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if strings.Contains(err.Error(), "400") {
			time.Sleep(time.Second * 3)
			resp, err = slbconn.ModifyLoadBalancer(&req)
			if err != nil {
				return fmt.Errorf("update Listener (%v)error:%v", req, err)
			}
		}
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunListenerRead(d, m)
}

func resourceKsyunListenerDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["ListenerId"] = d.Id()
	/*
		req["LoadBalancerId"] = d.Id()
		_, err := slbconn.DeleteLoadBalancer(&req)
		if err != nil {
			return fmt.Errorf("release Listener error:%v", err)
		}
		return nil
	*/
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteListeners"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteListeners(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["ListenerId.1"] = d.Id()
		action = "DescribeListeners"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeListeners(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading Listener when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["ListenerSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified Listener %q has not been deleted due to unknown error", d.Id()))
	})
}
