package ksyun

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func dataSourceKsyunVPCs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVPCsRead,

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

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cidr_block": {
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

func dataSourceKsyunVPCsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}

	client := meta.(*KsyunClient)
	conn := client.vpcconn
	readVpc := make(map[string]interface{})
	result = []map[string]interface{}{}

	if ids, ok := d.GetOk("ids"); ok {
		SchemaSetToInstanceMap(ids, "VpcId", &readVpc)
	}
	resp, err := conn.DescribeVpcs(&readVpc)
	if err != nil {
		return fmt.Errorf("error on reading vpc list req(%v):%v", readVpc, err)
	}
	l := (*resp)["VpcSet"].([]interface{})
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range l {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["VpcName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, l)
	}
	err = dataSourceKsyunVPCsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading vpc list, %s", err)
	}
	return nil
}

func dataSourceKsyunVPCsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	str, _ := json.Marshal(&result)
	fmt.Printf("%+v\n", string(str))
	//fmt.Printf("%+v\n", len(result))
	var ids []string
	var data []map[string]interface{}
	var err error

	ids = []string{}
	data = []map[string]interface{}{}

	for _, item := range result {
		ids = append(ids, item["VpcId"].(string))

		data = append(data, map[string]interface{}{
			"id":          item["VpcId"],
			"name":        item["VpcName"],
			"create_time": item["CreateTime"],
			"cidr_block":  item["CidrBlock"],
		})
	}
	d.SetId(hashStringArray(ids))
	err = d.Set("total_count", len(result))
	if err != nil {
		return err
	}
	err = d.Set("vpcs", data)
	if err != nil {
		return err
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		err = writeToFile(outputFile.(string), data)
		if err != nil {
			return err
		}
	}

	return nil
}
