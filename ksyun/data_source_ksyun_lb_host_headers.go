package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunListenerHostHeaders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunListenerHostHeadersRead,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_header_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
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

func dataSourceKsyunListenerHostHeadersRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var ListenerHostHeaders []string
	if ids, ok := d.GetOk("ids"); ok {
		ListenerHostHeaders = SchemaSetToStringSlice(ids)
	}
	for k, v := range ListenerHostHeaders {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("HostHeaderId.%d", k+1)] = v
	}
	filters := []string{"listener_id"}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	var allListenerHostHeaders []interface{}

	resp, err := conn.DescribeHostHeaders(&req)
	if err != nil {
		return fmt.Errorf("error on reading host header list req (%v):%v", req, err)
	}
	itemSet, ok := (*resp)["HostHeaderSet"]
	if !ok {
		return fmt.Errorf("error on reading HostHeader set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allListenerHostHeaders = append(allListenerHostHeaders, items...)
	//	excludes:=[]string{"HealthCheck","RealServer","Session"}
	datas := GetSubSliceDByRep(allListenerHostHeaders, hostHeaderKeys)
	err = dataSourceKscSave(d, "host_headers", ListenerHostHeaders, datas)
	if err != nil {
		return fmt.Errorf("error on save HostHeader list, %s", err)
	}
	return nil
}
