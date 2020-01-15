package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceKsyunLbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLbsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"lbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLbsRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	var slbs []string
	if ids, ok := d.GetOk("ids"); ok {
		slbs = SchemaSetToStringSlice(ids)
	}
	for k, v := range slbs {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("LoadBalancerId.%d", k+1)] = v
	}

	if ids, ok := d.GetOk("state"); ok {
		req["State"] = ids
	}
	filters := []string{"vpc_id"}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	var allSlbs []interface{}

	resp, err := conn.DescribeLoadBalancers(&req)
	if err != nil {
		return fmt.Errorf("error on reading Slb list, %s", req)
	}
	itemSet, ok := (*resp)["LoadBalancerDescriptions"]
	if !ok {
		return fmt.Errorf("error on reading Slb set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allSlbs = append(allSlbs, items...)
	datas := GetSubSliceDByRep(allSlbs, slbKeys)
	err = dataSourceKscSave(d, "lbs", slbs, datas)
	if err != nil {
		return fmt.Errorf("error on save Slb list, %s", err)
	}
	return nil
}
