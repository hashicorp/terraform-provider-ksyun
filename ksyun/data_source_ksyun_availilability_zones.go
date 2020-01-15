package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func dataSourceKsyunAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAvailabilityZonesRead,
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

			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunAvailabilityZonesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).vpcconn
	req := make(map[string]interface{})

	var allAvailabilityZones []interface{}
	action := "DescribeAvailabilityZones"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeAvailabilityZones(&req)
	if err != nil {
		return fmt.Errorf("error on reading AvailabilityZone list req(%v):%v", req, err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	itemSet, ok := (*resp)["AvailabilityZoneInfo"]
	if !ok {
		return fmt.Errorf("error on reading AvailabilityZone set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allAvailabilityZones = append(allAvailabilityZones, items...)
	datas := GetSubSliceDByRep(allAvailabilityZones, availabilityZoneKeys)
	err = dataSourceKscSave(d, "availability_zones", []string{time.Now().UTC().String()}, datas)
	if err != nil {
		return fmt.Errorf("error on save AvailabilityZone list, %s", err)
	}
	return nil
}
