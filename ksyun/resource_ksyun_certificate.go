package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunCertificateCreate,
		Read:   resourceKsyunCertificateRead,
		Update: resourceKsyunCertificateUpdate,
		Delete: resourceKsyunCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunCertificateCreate(d *schema.ResourceData, m interface{}) error {
	kcmconn := m.(*KsyunClient).kcmconn
	createCertificate := make(map[string]interface{})
	creates := []string{
		"certificate_name",
		"private_key",
		"public_key",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			if v == "private_key" || v == "public_key" {
				v1 = strings.Replace(fmt.Sprintf("%s", v1), "\n", "\\n", -1)
			}
			vv := Downline2Hump(v)
			createCertificate[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateCertificate"
	logger.Debug(logger.ReqFormat, action, createCertificate)
	resp, err := kcmconn.CreateCertificate(&createCertificate)
	logger.Debug(logger.AllFormat, action, createCertificate, *resp, err)
	if err != nil {
		return fmt.Errorf("createCertificate Error  : %s", err)
	}
	certificate, ok := (*resp)["Certificate"]
	if !ok {
		return fmt.Errorf("createCertificate Error  : no id found")
	}
	certificateItem, ok := certificate.(map[string]interface{})
	if !ok {
		return fmt.Errorf("createCertificate Error  : no id found")
	}
	id, ok := certificateItem["CertificateId"]
	if !ok {
		return fmt.Errorf("createCertificate Error  : no id found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("createCertificate Error : no id found")
	}
	d.Set("certificate_id", idres)
	d.SetId(idres)
	return resourceKsyunCertificateRead(d, m)
}

func resourceKsyunCertificateRead(d *schema.ResourceData, m interface{}) error {
	kcmconn := m.(*KsyunClient).kcmconn
	readCertificate := make(map[string]interface{})
	readCertificate["CertificateId.1"] = d.Id()

	action := "DescribeCertificates"
	logger.Debug(logger.ReqFormat, action, readCertificate)
	resp, err := kcmconn.DescribeCertificates(&readCertificate)
	logger.Debug(logger.AllFormat, action, readCertificate, *resp, err)
	if err != nil {
		return fmt.Errorf("Error  : %s", err)
	}
	itemset, ok := (*resp)["CertificateSet"]
	if !ok {
		d.SetId("")
		return nil
	}
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], certificateKeys, map[string]bool{})
	return nil
}

func resourceKsyunCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	kcmconn := m.(*KsyunClient).kcmconn
	// Enable partial attribute modification
	d.Partial(true)
	// Whether the representative has any modifications
	attributeUpdate := false
	updateReq := make(map[string]interface{})
	updateReq["CertificateId"] = d.Id()
	// modify
	if d.HasChange("certificate_name") && !d.IsNewResource() {
		if v, ok := d.GetOk("certificate_name"); ok {
			updateReq["CertificateName"] = fmt.Sprintf("%v", v)
			attributeUpdate = true
		}
	}
	if attributeUpdate {
		action := "ModifyCertificate"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := kcmconn.ModifyCertificate(&updateReq)
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		if err != nil {
			return fmt.Errorf("update Certificate (%v)error:%v", updateReq, err)
		}
		d.SetPartial("certificate_name")
	}
	d.Partial(false)
	return resourceKsyunCertificateRead(d, m)
}

func resourceKsyunCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kcmconn
	//delete
	deleteCertificate := make(map[string]interface{})
	deleteCertificate["CertificateId"] = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteCertificate"
		logger.Debug(logger.ReqFormat, action, deleteCertificate)
		resp, err1 := conn.DeleteCertificate(&deleteCertificate)
		logger.Debug(logger.AllFormat, action, deleteCertificate, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		//check
		readCertificate := make(map[string]interface{})
		readCertificate["CertificateId.1"] = d.Id()
		action = "DescribeCertificates"
		logger.Debug(logger.ReqFormat, action, readCertificate)
		resp, err := conn.DescribeCertificates(&readCertificate)
		logger.Debug(logger.AllFormat, action, readCertificate, *resp, err)
		if err != nil && notFoundError(err) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading Certificate when deleting %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["CertificateSet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting Certificate %v,%v", d.Id(), err1))
	})

}
