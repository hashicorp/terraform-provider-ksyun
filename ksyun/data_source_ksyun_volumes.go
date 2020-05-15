package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVolumesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				MaxItems: 100,
			},
			"volume_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunVolumesRead(d *schema.ResourceData, meta interface{}) error {
	action := "DescribeVolumes"
	conn := meta.(*KsyunClient).ebsconn
	readReq := make(map[string]interface{})
	var volumes []string
	if ids, ok := d.GetOk("ids"); ok {
		volumes = SchemaSetToStringSlice(ids)
	}
	for k, v := range volumes {
		if v == "" {
			continue
		}
		readReq[fmt.Sprintf("VolumeId.%d", k+1)] = v
	}
	if v, ok := d.GetOk("volume_category"); ok {
		readReq["VolumeCategory"] = v
	}
	if v, ok := d.GetOk("volume_status"); ok {
		readReq["VolumeStatus"] = v
	}
	if v, ok := d.GetOk("volume_type"); ok {
		readReq["VolumeType"] = v
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		readReq["AvailabilityZone"] = v
	}
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeVolumes(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading volume list req(%v):%s", readReq, err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	volumeList, ok := (*resp)["Volumes"]
	if !ok {
		return nil
	}
	items, ok := volumeList.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	var allVolumes []interface{}
	allVolumes = append(allVolumes, items...)
	datas := GetSubSliceDByRep(allVolumes, volumeKeys)
	err = dataSourceKscSave(d, "volumes", volumes, datas)
	if err != nil {
		return fmt.Errorf("error on save volume list, %s", err)
	}
	return nil
}
