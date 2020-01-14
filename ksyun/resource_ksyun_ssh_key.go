package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSSHKeyCreate,
		Read:   resourceKsyunSSHKeyRead,
		Update: resourceKsyunSSHKeyUpdate,
		Delete: resourceKsyunSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	sksconn := m.(*KsyunClient).sksconn
	createSSHKey := make(map[string]interface{})
	creates := []string{
		"key_name",
		"public_key",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createSSHKey[vv] = fmt.Sprintf("%v", v1)
		}
	}
	var action string
	var resp *map[string]interface{}
	var err error
	if _, ok := d.GetOk("public_key"); ok {
		action = "ImportKey"
		logger.Debug(logger.ReqFormat, action, createSSHKey)
		resp, err = sksconn.ImportKey(&createSSHKey)
	} else {
		action = "CreateKey"
		logger.Debug(logger.ReqFormat, action, createSSHKey)
		resp, err = sksconn.CreateKey(&createSSHKey)
	}
	logger.Debug(logger.AllFormat, action, createSSHKey, *resp, err)
	if err != nil {
		return fmt.Errorf("createSSHKey Error  : %s", err)
	}
	if _, ok := d.GetOk("public_key"); !ok {
		privateKey, ok := (*resp)["PrivateKey"]
		if !ok {
			return fmt.Errorf("createSSHKey Error  : no PrivateKey found")
		}
		d.Set("private_key", privateKey)
	}
	sSHKey, ok := (*resp)["Key"]
	if !ok {
		return fmt.Errorf("createSSHKey Error  : no id found")
	}
	sSHKeyItem, ok := sSHKey.(map[string]interface{})
	if !ok {
		return fmt.Errorf("createSSHKey Error  : no id found")
	}
	id, ok := sSHKeyItem["KeyId"]
	if !ok {
		return fmt.Errorf("createSSHKey Error  : no id found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("createSSHKey Error : no id found")
	}
	d.Set("key_id", idres)
	d.SetId(idres)
	return resourceKsyunSSHKeyRead(d, m)
}

func resourceKsyunSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	sksconn := m.(*KsyunClient).sksconn
	readSSHKey := make(map[string]interface{})
	readSSHKey["KeyId.1"] = d.Id()

	action := "DescribeKeys"
	logger.Debug(logger.ReqFormat, action, readSSHKey)
	resp, err := sksconn.DescribeKeys(&readSSHKey)
	logger.Debug(logger.AllFormat, action, readSSHKey, *resp, err)
	if err != nil {
		return fmt.Errorf("Error  : %s", err)
	}
	privateKey, ok := (*resp)["PrivateKey"]
	if ok {
		d.Set("private_key", privateKey)
	}
	itemset, ok := (*resp)["KeySet"]
	if !ok {
		d.SetId("")
		return nil
	}
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], sshKeyKeys, map[string]bool{})
	return nil
}

func resourceKsyunSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {
	sksconn := m.(*KsyunClient).sksconn
	// Enable partial attribute modification
	d.Partial(true)
	// Whether the representative has any modifications
	attributeUpdate := false
	updateReq := make(map[string]interface{})
	updateReq["KeyId"] = d.Id()
	// modify
	if d.HasChange("key_name") && !d.IsNewResource() {
		if v, ok := d.GetOk("key_name"); ok {
			updateReq["KeyName"] = fmt.Sprintf("%v", v)
			attributeUpdate = true
		}
	}
	if attributeUpdate {
		action := "ModifyKey"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := sksconn.ModifyKey(&updateReq)
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		if err != nil {
			return fmt.Errorf("update SSHKey (%v)error:%v", updateReq, err)
		}
		d.SetPartial("key_name")
	}
	d.Partial(false)
	return resourceKsyunSSHKeyRead(d, m)
}

func resourceKsyunSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sksconn
	//delete
	deleteSSHKey := make(map[string]interface{})
	deleteSSHKey["KeyId"] = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteKey"
		logger.Debug(logger.ReqFormat, action, deleteSSHKey)
		resp, err1 := conn.DeleteKey(&deleteSSHKey)
		logger.Debug(logger.AllFormat, action, deleteSSHKey, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		//check
		readSSHKey := make(map[string]interface{})
		readSSHKey["KeyId.1"] = d.Id()
		action = "DescribeKeys"
		logger.Debug(logger.ReqFormat, action, readSSHKey)
		resp, err := conn.DescribeKeys(&readSSHKey)
		logger.Debug(logger.AllFormat, action, readSSHKey, *resp, err)
		if err != nil && notFoundError(err) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading SSHKey when deleting %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["KeySet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting SSHKey %v,%v", d.Id(), err1))
	})

}
