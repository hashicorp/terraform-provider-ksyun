package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunHealthChecks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunHealthChecksRead,
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
			"health_checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunHealthChecksRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var healthChecks []string
	if ids, ok := d.GetOk("ids"); ok {
		healthChecks = SchemaSetToStringSlice(ids)
	}
	for k, v := range healthChecks {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("HealthCheckId.%d", k+1)] = v
	}
	filters := []string{"listener_id"}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	var allHealthChecks []interface{}

	resp, err := conn.DescribeHealthChecks(&req)
	if err != nil {
		return fmt.Errorf("error on reading HealthCheck list, %s", req)
	}
	itemSet, ok := (*resp)["HealthCheckSet"]
	if !ok {
		return fmt.Errorf("error on reading HealthCheck set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allHealthChecks = append(allHealthChecks, items...)
	datas := GetSubSliceDByRep(allHealthChecks, healthCheckKeys)
	err = dataSourceKscSave(d, "health_checks", healthChecks, datas)
	if err != nil {
		return fmt.Errorf("error on save HealthCheck list, %s", err)
	}
	return nil
}
