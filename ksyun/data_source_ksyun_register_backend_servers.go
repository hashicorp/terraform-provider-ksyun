package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunRegisterBackendServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunRegisterBackendServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"backend_server_group_id": {
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
			"register_backend_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backend_server_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"register_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_server_ip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"real_server_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_server_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunRegisterBackendServersRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var backendServers []string
	if ids, ok := d.GetOk("register_id"); ok {
		backendServers = SchemaSetToStringSlice(ids)
	}
	for k, v := range backendServers {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("RegisterId.%d", k+1)] = v
	}
	var allRegisterBackendServers []interface{}
	resp, err := conn.DescribeBackendServers(&req)
	if err != nil {
		return fmt.Errorf("error on reading register_backend_servers req (%v):%v", req, err)
	}
	filters := []string{
		"backend_server_group_id",
	}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	itemSet, ok := (*resp)["BackendServerSet"]
	if !ok {
		return fmt.Errorf("error on reading register_backend_servers")
	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allRegisterBackendServers = append(allRegisterBackendServers, items...)
	datas := GetSubSliceDByRep(allRegisterBackendServers, registerBackendServerKeys)
	err = dataSourceKscSave(d, "register_backend_servers", backendServers, datas)
	if err != nil {
		return fmt.Errorf("error on save register_backend_servers list, %s", err)
	}
	return nil
}
