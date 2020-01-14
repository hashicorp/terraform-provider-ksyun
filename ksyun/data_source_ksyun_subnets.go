package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"regexp"
)

func dataSourceKsyunSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSubnetssRead,

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

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"subnet_types": {
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
			"subnets": {
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

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_ip_from": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_ip_to": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dns1": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dns2": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"availability_zone_name": {
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

func dataSourceKsyunSubnetssRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var allSubnets []interface{}

	limit := 100
	offset := 1

	client := meta.(*KsyunClient)
	conn := client.vpcconn
	readSubnet := make(map[string]interface{})
	result = []map[string]interface{}{}
	allSubnets = []interface{}{}

	if ids, ok := d.GetOk("ids"); ok {
		SchemaSetToInstanceMap(ids, "SubnetId", &readSubnet)
	}
	//filter
	index := int(0)
	if vpcIds, ok := d.GetOk("vpc_ids"); ok {
		index = index + 1
		SchemaSetToFilterMap(vpcIds, "vpc-id", index, &readSubnet)
	}
	if subnetTypes, ok := d.GetOk("subnet_types"); ok {
		index = index + 1
		SchemaSetToFilterMap(subnetTypes, "subnet-type", index, &readSubnet)
	}

	for {
		readSubnet["MaxResults"] = limit
		readSubnet["NextToken"] = offset

		resp, err := conn.DescribeSubnets(&readSubnet)
		if err != nil {
			return fmt.Errorf("error on reading subnet list req(%v):%v", readSubnet, err)
		}
		l := (*resp)["SubnetSet"].([]interface{})
		allSubnets = append(allSubnets, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range allSubnets {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["SubnetName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, allSubnets)
	}

	err := dataSourceKsyunSubnetsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading subnet list, %s", err)
	}
	return nil
}

func dataSourceKsyunSubnetsSave(d *schema.ResourceData, result []map[string]interface{}) error {

	var ids []string
	var data []map[string]interface{}
	var err error

	ids = []string{}
	data = []map[string]interface{}{}

	for _, item := range result {
		ids = append(ids, item["SubnetId"].(string))

		data = append(data, map[string]interface{}{
			"id":                     item["SubnetId"],
			"name":                   item["SubnetName"],
			"cidr_block":             item["CidrBlock"],
			"vpc_id":                 item["VpcId"],
			"subnet_type":            item["SubnetType"],
			"dhcp_ip_from":           item["DhcpIpFrom"],
			"dhcp_ip_to":             item["DhcpIpTo"],
			"gateway_ip":             item["GatewayIp"],
			"dns1":                   item["Dns1"],
			"dns2":                   item["Dns2"],
			"network_acl_id":         item["NetworkAclId"],
			"nat_id":                 item["NatId"],
			"availability_zone_name": item["AvailabilityZoneName"],
			"create_time":            item["CreateTime"],
		})
	}
	d.SetId(hashStringArray(ids))
	err = d.Set("total_count", len(result))
	if err != nil {
		return err
	}
	err = d.Set("subnets", data)
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
