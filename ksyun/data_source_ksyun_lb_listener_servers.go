package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKsyunLbListenerServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunServersRead,
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
			"listener_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"real_server_ip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
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
						"real_server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"real_server_state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"register_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLbServersRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var servers []string
	if ids, ok := d.GetOk("ids"); ok {
		servers = SchemaSetToStringSlice(ids)
	}
	for k, v := range servers {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("RegisterId.%d", k+1)] = v
	}
	filters := []string{"listener_id", "real_server_ip"}
	req = *SchemaSetsToFilterMap(d, filters, &req)
	var allServers []interface{}

	resp, err := conn.DescribeInstancesWithListener(&req)
	if err != nil {
		return fmt.Errorf("error on reading Server list (%s):%s", req, err)
	}
	itemSet, ok := (*resp)["RealServerSet"]
	if !ok {
		return fmt.Errorf("error on reading Server set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allServers = append(allServers, items...)
	datas := GetSubSliceDByRep(allServers, serverKeys)
	err = dataSourceKscSave(d, "servers", servers, datas)
	if err != nil {
		return fmt.Errorf("error on save server list, %s", err)
	}
	return nil
}
