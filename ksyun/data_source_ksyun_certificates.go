package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceKsyunCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunCertificatesRead,
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

			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunCertificatesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).kcmconn
	var Certificates []string
	req := make(map[string]interface{})

	if ids, ok := d.GetOk("ids"); ok {
		Certificates = SchemaSetToStringSlice(ids)
	}
	for k, v := range Certificates {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("CertificateId.%d", k+1)] = v
	}
	resp, err := conn.DescribeCertificates(&req)
	if err != nil {
		return fmt.Errorf("error on reading Certificate list req(%v):%v", req, err)
	}
	itemSet, ok := (*resp)["CertificateSet"]
	if !ok {
		return nil
	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}

	datas := GetSubSliceDByRep(items, certificateKeys)
	err = dataSourceKscSave(d, "certificates", Certificates, datas)
	if err != nil {
		return fmt.Errorf("error on save Certificate list, %s", err)
	}
	return nil
}
