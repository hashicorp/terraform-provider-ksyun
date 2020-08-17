package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunBackendServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBackendServerGroupsRead,
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
			"backend_server_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_server_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backend_server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backend_server_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backend_server_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
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
					},
				},
			},
		},
	}
}

func dataSourceKsyunBackendServerGroupsRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var backendServerGroups []string
	if ids, ok := d.GetOk("ids"); ok {
		backendServerGroups = SchemaSetToStringSlice(ids)
	}
	for k, v := range backendServerGroups {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("BackendServerGroupId.%d", k+1)] = v
	}
	var allBackendServerGroups []interface{}
	resp, err := conn.DescribeBackendServerGroups(&req)
	if err != nil {
		return fmt.Errorf("error on reading backend server group list req (%v):%v", req, err)
	}
	itemSet, ok := (*resp)["BackendServerGroupSet"]
	if !ok {
		return fmt.Errorf("error on reading backend server group set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allBackendServerGroups = append(allBackendServerGroups, items...)
	datas := GetSubSliceDByRep(allBackendServerGroups, backendServerGroupKeys)
	dealListenrData(datas)
	err = dataSourceKscSave(d, "backend_server_groups", backendServerGroups, datas)
	if err != nil {
		return fmt.Errorf("error on save backend_server_group list, %s", err)
	}
	return nil
}
