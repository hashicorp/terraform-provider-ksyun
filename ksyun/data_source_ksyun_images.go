package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"regexp"
)

func dataSourceKsyunImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunImagesRead,
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
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_npe": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sys_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunImagesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).kecconn
	req := make(map[string]interface{})
	var imageIds []string
	var allImages []interface{}
	resp, err := conn.DescribeImages(&req)
	if err != nil {
		return fmt.Errorf("error on reading Image list req(%v):%v", req, err)
	}
	itemSet, ok := (*resp)["ImagesSet"]
	if !ok {
		return fmt.Errorf("error on reading Image set")
	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allImages = append(allImages, items...)
	datas := GetSubSliceDByRep(allImages, imageKeys)
	if name, ok := d.GetOk("platform"); ok {
		var dataFilter []map[string]interface{}
		for _, v := range datas {
			if v["platform"] == name {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	if standard, ok := d.GetOk("is_public"); ok {
		var dataFilter []map[string]interface{}
		for _, v := range datas {
			if v["is_public"] == standard {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	if imageSource, ok := d.GetOk("image_source"); ok {
		var dataFilter []map[string]interface{}
		for _, v := range datas {
			if v["image_source"] == imageSource {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		var dataFilter []map[string]interface{}
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range datas {
			if r == nil || r.MatchString(v["name"].(string)) {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	err = dataSourceKscSave(d, "images", imageIds, datas)
	if err != nil {
		return fmt.Errorf("error on save Images list, %s", err)
	}
	return nil
}
