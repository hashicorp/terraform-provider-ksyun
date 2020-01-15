package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKsyunLines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLinesRead,
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
			"line_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"lines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLinesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).eipconn
	req := make(map[string]interface{})
	var lines []string

	var allLines []interface{}

	resp, err := conn.GetLines(&req)
	if err != nil {
		return fmt.Errorf("error on reading Line list req(%v):%v", req, err)
	}
	itemSet, ok := (*resp)["LineSet"]
	if !ok {
		return fmt.Errorf("error on reading Line set")

	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allLines = append(allLines, items...)
	datas := GetSubSliceDByRep(allLines, lineKeys)
	if name, ok := d.GetOk("line_name"); ok {
		var dataFilter []map[string]interface{}
		for k, v := range datas {
			if v["line_name"] == name {
				dataFilter = append(dataFilter, datas[k])
				break
			}
		}
		datas = dataFilter
	}

	err = dataSourceKscSave(d, "lines", lines, datas)
	if err != nil {
		return fmt.Errorf("error on save lines list, %s", err)
	}
	return nil
}
