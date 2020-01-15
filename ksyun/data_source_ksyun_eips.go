package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceKsyunEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunEipsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"project_id": {
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
			"network_interface_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"instance_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"internet_gateway_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"band_width_share_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"line_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"public_ip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"band_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"band_width_share_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_band_width_share": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunEipsRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).eipconn
	var eips []string
	req := make(map[string]interface{})

	if ids, ok := d.GetOk("ids"); ok {
		eips = SchemaSetToStringSlice(ids)
	}
	for k, v := range eips {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("AllocationId.%d", k+1)] = v
	}
	var projectIds []string
	if ids, ok := d.GetOk("project_id"); ok {
		projectIds = SchemaSetToStringSlice(ids)
	}

	for k, v := range projectIds {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("ProjectId.%d", k+1)] = v
	}

	filters := []string{
		"network_interface_id",
		"instance_type",
		"internet_gateway_id",
		"band_width_share_id",
		"line_id",
		"public_ip",
		"project_id",
	}
	req = *SchemaSetsToFilterMap(d, filters, &req)

	var allEips []interface{}
	var limit int = 100
	var nextToken string
	for {
		req["MaxResults"] = limit
		if nextToken != "" {
			req["NextToken"] = nextToken
		}
		resp, err := conn.DescribeAddresses(&req)
		if err != nil {
			return fmt.Errorf("error on reading eip list req(%v):%v", req, err)
		}
		itemSet, ok := (*resp)["AddressesSet"]
		if !ok {
			//return fmt.Errorf("error on reading eip set")
			break
		}
		items, ok := itemSet.([]interface{})
		if !ok {
			break
		}
		if items == nil || len(items) < 1 {
			break
		}
		allEips = append(allEips, items...)

		if len(items) < limit {
			break
		}
		if nextTokens, ok := (*resp)["NextToken"]; ok {
			nextToken = fmt.Sprintf("%v", nextTokens)
		} else {
			break
		}
	}
	datas := GetSubSliceDByRep(allEips, eipKeys)
	err := dataSourceKscSave(d, "eips", eips, datas)
	if err != nil {
		return fmt.Errorf("error on save eip list, %s", err)
	}
	return nil
}
