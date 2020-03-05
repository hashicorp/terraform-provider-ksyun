package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKsyunListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunListenersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"method": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"listener_id": {
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
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"session_persistence_period": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"session_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_name": {
										Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func dataSourceKsyunListenersRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var Listeners []string
	if ids, ok := d.GetOk("ids"); ok {
		Listeners = SchemaSetToStringSlice(ids)
	}
	for k, v := range Listeners {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("ListenerId.%d", k+1)] = v
	}
	filters := []string{"load_balancer_id"}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	var allListeners []interface{}

	resp, err := conn.DescribeListeners(&req)
	if err != nil {
		return fmt.Errorf("error on reading listener list req (%v):%v", req, err)
	}
	itemSet, ok := (*resp)["ListenerSet"]
	if !ok {
		return fmt.Errorf("error on reading Listener set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allListeners = append(allListeners, items...)
	//	excludes:=[]string{"HealthCheck","RealServer","Session"}
	datas := GetSubSliceDByRep(allListeners, listenerKeys)
	dealListenrData(datas)
	err = dataSourceKscSave(d, "listeners", Listeners, datas)
	if err != nil {
		return fmt.Errorf("error on save Listener list, %s", err)
	}
	return nil
}

func dealListenrData(datas []map[string]interface{}) {
	for k, v := range datas {
		for k1, v1 := range v {
			switch k1 {
			case "health_check":
				datas[k]["health_check"] = GetSubSliceDByRep([]interface{}{v1}, healthCheckKeys)
			case "real_server":
				vv := v1.([]interface{})
				datas[k]["real_server"] = GetSubSliceDByRep(vv, serverKeys)
			case "session":
				datas[k]["session"] = GetSubSliceDByRep([]interface{}{v1}, sessionKeys)
			}
		}
	}
}
