package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunSSHKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSSHKeysRead,
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
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key": {
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

func dataSourceKsyunSSHKeysRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).sksconn
	var sSHKeys []string
	req := make(map[string]interface{})

	if ids, ok := d.GetOk("ids"); ok {
		sSHKeys = SchemaSetToStringSlice(ids)
	}
	for k, v := range sSHKeys {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("SSHKeyId.%d", k+1)] = v
	}
	if name, ok := d.GetOk("key_name"); ok {
		req["Filter.1.Name"] = "key-name"
		req["Filter.1.Value.1"] = fmt.Sprintf("%v", name)
	}
	var alls []interface{}
	var limit int = 30
	var nextToken string
	for {
		req["MaxResults"] = fmt.Sprintf("%v", limit)
		if !(nextToken == "" || nextToken == "null") {
			req["NextToken"] = nextToken
		}

		resp, err := conn.DescribeKeys(&req)
		if err != nil {
			return fmt.Errorf("error on reading SSHKey list req(%v):%v", req, err)
		}
		logger.Debug(logger.RespFormat, "DescribeKeys", req, *resp)
		itemSet, ok := (*resp)["KeySet"]
		if !ok {
			break
		}
		items, ok := itemSet.([]interface{})
		if !ok {
			break
		}
		if items == nil || len(items) < 1 {
			break
		}
		alls = append(alls, items...)
		if nextTokens, ok := (*resp)["NextToken"]; ok {
			if nextTokens == "null" || nextTokens == "" || nextTokens == nil {
				break
			}
			nextToken = fmt.Sprintf("%v", nextTokens)
		} else {
			break
		}
	}
	datas := GetSubSliceDByRep(alls, sshKeyKeys)
	err := dataSourceKscSave(d, "keys", sSHKeys, datas)
	if err != nil {
		return fmt.Errorf("error on save SSHKey list, %s", err)
	}
	return nil
}
