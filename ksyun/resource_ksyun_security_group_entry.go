package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
)

func resourceKsyunSecurityGroupEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupEntryCreate,
		Read:   resourceKsyunSecurityGroupEntryRead,
		Delete: resourceKsyunSecurityGroupEntryDelete,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"icmp_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"icmp_code": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port_range_from": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port_range_to": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunSecurityGroupEntryCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	creates := []string{
		"description",
		"security_group_id",
		"cidr_block",
		"direction",
		"protocol",
		"icmp_type",
		"icmp_code",
		"port_range_from",
		"port_range_to",
	}
	createSecurityGroupEntry := make(map[string]interface{})
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createSecurityGroupEntry[vv] = fmt.Sprintf("%v", v1)
		}
	}
	var resp *map[string]interface{}
	var err error
	action := "AuthorizeSecurityGroupEntry"
	logger.Debug(logger.ReqFormat, action, createSecurityGroupEntry)
	resp, err = conn.AuthorizeSecurityGroupEntry(&createSecurityGroupEntry)
	logger.Debug(logger.AllFormat, action, createSecurityGroupEntry, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating SecurityGroupEntry, %s", err)
	}
	idSet, ok := (*resp)["SecurityGroupEntryIdSet"]
	if !ok {
		return fmt.Errorf("AuthorizeSecurityGroupEntry Error  : no ids found")
	}
	ids, ok := idSet.([]interface{})
	if !ok || len(ids) == 0 {
		return fmt.Errorf("AuthorizeSecurityGroupEntry Error  : no ids found")
	}
	idres, ok := ids[0].(string)
	if !ok {
		return fmt.Errorf("AuthorizeSecurityGroupEntry Error : no id found")
	}
	if err := d.Set("security_group_entry_id", idres); err != nil {
		return err
	}
	d.SetId(idres)
	return nil
}
func resourceKsyunSecurityGroupEntryRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	readSecurityGroup := make(map[string]interface{})
	if _, ok := d.GetOk("security_group_id"); !ok {
		return fmt.Errorf("security_group_id can't be null")
	}
	sgId := d.Get("security_group_id")
	readSecurityGroup["SecurityGroupId.1"] = sgId

	action := "DescribeSecurityGroups"
	logger.Debug(logger.ReqFormat, action, readSecurityGroup)
	resp, err := conn.DescribeSecurityGroups(&readSecurityGroup)
	logger.Debug(logger.AllFormat, action, readSecurityGroup, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading SecurityGroup %q, %s", d.Id(), err)
	}

	itemset := (*resp)["SecurityGroupSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if item["security_group_entry_id"] == d.Id() {
			for sgk, sgv := range item {
				if vpcSecurityGroupEntrySetKeys[sgk] {
					if err := d.Set(Downline2Hump(sgk), sgv); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func resourceKsyunSecurityGroupEntryDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	//delete
	if _, ok := d.GetOk("security_group_id"); !ok {
		return fmt.Errorf("security_group_id can't be null")
	}
	deleteSecurityGroupEntry := make(map[string]interface{})
	deleteSecurityGroupEntry["SecurityGroupId"] = d.Get("security_group_id")
	deleteSecurityGroupEntry["SecurityGroupEntryId"] = d.Id()
	action := "RevokeSecurityGroupEntry"
	logger.Debug(logger.ReqFormat, action, deleteSecurityGroupEntry)
	resp, err0 := conn.RevokeSecurityGroupEntry(&deleteSecurityGroupEntry)
	logger.Debug(logger.AllFormat, action, deleteSecurityGroupEntry, *resp, err0)
	if err0 == nil || (err0 != nil && strings.Contains(err0.Error(), "NotFound")) {
		return nil
	}
	//check
	readSecurityGroup := make(map[string]interface{})
	readSecurityGroup["SecurityGroupId.1"] = d.Get("security_group_id")
	action = "DescribeSecurityGroups"
	logger.Debug(logger.ReqFormat, action, readSecurityGroup)
	resp, err1 := conn.DescribeSecurityGroups(&readSecurityGroup)
	logger.Debug(logger.AllFormat, action, readSecurityGroup, *resp, err1)
	if err1 != nil && strings.Contains(err1.Error(), "NotFound") {
		return nil
	}
	if err1 != nil {
		return nil
	}
	itemSet, ok := (*resp)["SecurityGroupSet"]
	if !ok {
		return nil
	}
	items, ok := itemSet.([]interface{})
	if !ok || len(items) == 0 {
		return nil
	}
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		sgeSet, ok := item["SecurityGroupEntrySet"]
		if !ok {
			continue
		}
		sges, ok := sgeSet.([]interface{})
		if !ok {
			continue
		}
		for k, sge := range sges {
			sgeMap, ok := sge.(map[string]interface{})
			if !ok {
				continue
			}
			if sgeMap["SecurityGroupEntryId"] == d.Id() {
				return fmt.Errorf("error on  deleting SecurityGroups %q, %s", d.Id(), err0)
			}
			if k == len(sges)-1 {
				return nil
			}
		}
	}
	return fmt.Errorf("error on  deleting SecurityGroupEntrys %q, %s", d.Id(), err0)
}
