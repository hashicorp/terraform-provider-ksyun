package ksyun

//import "github.com/hashicorp/terraform/helper/schema"

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"log"
	"strings"
	"time"
)

func resourceKsyunInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstanceCreate,
		Update: resourceKsyunInstanceUpdate,
		Read:   resourceKsyunInstanceRead,
		Delete: resourceKsyunInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"data_disk_gb": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			/*
			   "max_count": {
			      Type:     schema.TypeInt,
			      Required: true,
			   },
			   "min_count": {
			      Type:     schema.TypeInt,
			      Required: true,
			   },
			*/
			"instance_password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"purchase_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sriov_net_support": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"data_guard_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"address_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_configure": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"v_c_p_u": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"g_p_u": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_disk_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						/*
						   "root_disk_gb": {
						      Type:     schema.TypeInt,
						      Computed: true,
						   },
						   "data_disk_type": {
						      Type:     schema.TypeString,
						      Computed: true,
						   },
						*/
					},
				},
			},
			"instance_state": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"monitoring": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"group_set": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"d_n_s1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"d_n_s2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stopped_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_what": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_scaling_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_show_sriov_net_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceKsyunInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	var err error
	creates := []string{
		"image_id",
		"instance_type",
		// "system_disk",
		"data_disk_gb",
		// "data_disk",
		//"max_count",=1
		//"min_count",=1
		"subnet_id",
		//"instance_password",于keyid冲突
		"keep_image_login",
		"charge_type",
		"purchase_time",
		//"security_group_id",
		"private_ip_address",
		"instance_name",
		"instance_name_suffix",
		"sriov_net_support",
		"project_id",
		"data_guard_id",
		"address_band_width",
		"line_id",
		"address_charge_type",
		"address_purchase_time",
		"address_project_id",
		"host_name",
		"user_data",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createReq[vv] = fmt.Sprintf("%v", v1)
		}
	}
	keyIds, ok := d.GetOk("key_id")
	var keyset bool
	if ok {
		keys := SchemaSetToStringSlice(keyIds)
		for k, v := range keys {
			key := fmt.Sprintf("KeyId.%d", k+1)
			createReq[key] = fmt.Sprintf("%v", v)
		}
		if len(keys) > 0 {
			keyset = true
		}
	}
	keepImageLogin := d.Get("keep_image_login")
	var keepImage bool
	if ok {
		keepImage = keepImageLogin.(bool)
	}
	if (!keepImage) && (!keyset) {
		createReq["InstancePassword"] = d.Get("instance_password")
	} else {
		if err := d.Set("instance_password", ""); err != nil {
			return err
		}
	}
	securityGroupIds, ok := d.GetOk("security_group_id")
	if !ok {
		return fmt.Errorf("no SecurityGroupId get")
	}
	securityGroups := SchemaSetToStringSlice(securityGroupIds)
	if len(securityGroups) == 0 {
		return fmt.Errorf("no SecurityGroupId get")
	}
	createReq["SecurityGroupId"] = securityGroups[0]
	createReq["MaxCount"] = "1"
	createReq["MinCount"] = "1"
	createStructs := []string{
		"system_disk",
	}
	for _, v := range createStructs {
		if v1, ok := d.GetOk(v); ok {
			FlatternStructPrefix(v1, &createReq, "SystemDisk")
		}
	}
	if v1, ok := d.GetOk("data_disk"); ok {
		FlatternStructSlicePrefix(v1, &createReq, "DataDisk")
	}
	action := "RunInstances"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.RunInstances(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating Instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		instances := (*resp)["InstancesSet"].([]interface{})
		if len(instances) == 0 {
			return fmt.Errorf("error on creating Instance")
		}
		Instance := instances[0].(map[string]interface{})
		InstanceId := Instance["InstanceId"].(string)
		d.SetId(InstanceId)
	}
	// after create instance, we need to wait it initialized
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"active"},
		Refresh:    instanceStateRefreshForCreateFunc(conn, d.Id(), []string{"active"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err := resourceKsyunInstanceRead(d, meta); err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("error on waiting for instance %q complete creating, %s", d.Id(), err)
	}
	return nil
}

func resourceKsyunInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	readReq := make(map[string]interface{})
	readReq["InstanceId.1"] = d.Id()
	if pd, ok := d.GetOk("project_id"); ok {
		readReq["project_id"] = fmt.Sprintf("%v", pd)
	}
	action := "DescribeInstances"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeInstances(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.AllFormat, action, readReq, *resp, err)
	itemset := (*resp)["InstancesSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	excludesKeys := map[string]bool{
		"InstanceConfigure":   true,
		"InstanceState":       true,
		"Monitoring":          true,
		"NetworkInterfaceSet": true,
		"SystemDisk":          true,
		"KeySet":              true,
	}
	excludes := SetDByResp(d, items[0], instanceKeys, excludesKeys)
	// if excludes["KeySet"] != nil {
	if err := d.Set("key_id", excludes["KeySet"]); err != nil {
		return err
	}
	//log.Println("key_id:%v", excludes["KeySet"])
	// }
	if excludes["InstanceConfigure"] != nil {
		itemSet := GetSubDByRep(excludes["InstanceConfigure"], instanceConfigureKeys, map[string]bool{})
		if len(itemSet) > 0 {
			if instanceConfigure, ok := itemSet[0].(map[string]interface{}); ok {
				if err := d.Set("data_disk_gb", instanceConfigure["data_disk_gb"]); err != nil {
					return err
				}
			}
		}
		if err := d.Set(Hump2Downline("InstanceConfigure"), itemSet); err != nil {
			return err
		}
	} else {
		if err := d.Set(Hump2Downline("InstanceConfigure"), nil); err != nil {
			return err
		}
	}
	if excludes["InstanceState"] != nil {
		itemSet := GetSubDByRep(excludes["InstanceState"], instanceStateKeys, map[string]bool{})
		if err := d.Set(Hump2Downline("InstanceState"), itemSet); err != nil {
			return err
		}
	} else {
		if err := d.Set(Hump2Downline("InstanceState"), nil); err != nil {
			return err
		}
	}
	if excludes["Monitoring"] != nil {
		itemSet := GetSubDByRep(excludes["Monitoring"], monitoringKeys, map[string]bool{})
		if err := d.Set(Hump2Downline("Monitoring"), itemSet); err != nil {
			return err
		}
	} else {
		if err := d.Set(Hump2Downline("Monitoring"), nil); err != nil {
			return err
		}
	}
	if excludes["NetworkInterfaceSet"] != nil {
		networkSet := excludes["NetworkInterfaceSet"]
		networks, ok := networkSet.([]interface{})
		if !ok {
			return fmt.Errorf("no network interfaces on reading Instance %q", d.Id())
		}
		var networkInterfaceItem map[string]interface{}
		log.Printf("networks:%v", networks)
		for _, v := range networks {
			value, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("no network interface on reading Instance %q", d.Id())
			}
			nit, ok := value["NetworkInterfaceType"]
			if !ok || nit != "primary" { //only get master network interface
				continue
			}
			networkInterfaceItem = value
		}
		log.Printf("networkInterfaceItem:%v", networkInterfaceItem)
		excludeKeys := map[string]bool{
			"SecurityGroupSet": true,
			"GroupSet":         true,
		}
		SetDByResp(d, networkInterfaceItem, kecNetworkInterfaceKeys, excludeKeys)
		if gs, ok := networkInterfaceItem["GroupSet"]; ok {
			itemSetSub := GetSubSliceDByRep(gs.([]interface{}), groupSetKeys)
			if err := d.Set("group_set", itemSetSub); err != nil {
				return err
			}
		}
		//SecurityGroupSet instance interface only support one ,so should get from network interface.

		/*if sg, ok := v["security_group_set"]; ok {
		     itemSetSub := GetSubSliceDByRep(sg.([]interface{}), kecSecurityGroupKeys)
		     itemSet[k]["security_group_set"] = itemSetSub
		     if len(itemSetSub) == 1 {
		        if sgIdnew, ok := itemSetSub[0]["security_group_id"]; ok {
		           d.Set("security_group_id", sgIdnew)
		        }
		     }
		  }
		*/
		vpcConn := meta.(*KsyunClient).vpcconn
		readReq := make(map[string]interface{})
		netId := d.Get("network_interface_id")
		readReq["NetworkInterfaceId.1"] = netId
		action := "DescribeNetworkInterfaces"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err := vpcConn.DescribeNetworkInterfaces(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.AllFormat, action, readReq, *resp, err)
		itemset := (*resp)["NetworkInterfaceSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return fmt.Errorf("no data on reading Instance (%q) networkInterfaceSet(%q)", d.Id(), netId)
		}

		excludesKeys := map[string]bool{
			"SecurityGroupSet": true,
			"InstanceType":     true,
		}
		excludes := SetDByResp(d, items[0], networkInterfaceKeys, excludesKeys)
		if sg, ok := excludes["SecurityGroupSet"]; ok {
			itemSetSub := GetSubSliceDByRep(sg.([]interface{}), kecSecurityGroupKeys)
			if len(itemSetSub) != 0 {
				var itemSetSlice []string
				for _, v := range itemSetSub {
					for k1, v1 := range v {
						if k1 == "security_group_id" {
							itemSetSlice = append(itemSetSlice, fmt.Sprintf("%v", v1))
						}
					}
				}
				if err := d.Set("security_group_id", itemSetSlice); err != nil {
					return err
				}
			}
		}
	} else {
		if err := d.Set(Hump2Downline("NetworkInterfaceSet"), nil); err != nil {
			return err
		}
	}
	if excludes["SystemDisk"] != nil {
		itemSet := GetSubDByRep(excludes["SystemDisk"], systemDiskKeys, map[string]bool{})
		if err := d.Set(Hump2Downline("SystemDisk"), itemSet); err != nil {
			return err
		}
	} else {
		if err := d.Set(Hump2Downline("SystemDisk"), nil); err != nil {
			return err
		}
	}
	return nil
}

func resourceKsyunInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	d.Partial(true)
	imageUpdate := false
	updateReq := make(map[string]interface{})
	updateReq["InstanceId"] = d.Id()

	//ModifyInstanceAttribute instancename
	attributeNameUpdate := false
	updateInstanceAttributes := []string{
		"instance_name",
	}
	for _, v := range updateInstanceAttributes {
		if d.HasChange(v) && !d.IsNewResource() {
			updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			attributeNameUpdate = true
		}
	}
	if attributeNameUpdate {
		action := "ModifyInstanceAttribute"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := conn.ModifyInstanceAttribute(&updateReq)
		if err != nil {
			return fmt.Errorf("error on updating  instance name, %s", err)
		}
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		for _, v := range updateInstanceAttributes {
			d.SetPartial(v)
		}
	}

	//ModifyInstanceImage
	updateImages := []string{
		"image_id",
		// "system_disk",
		//"instance_password",
		"keep_image_login",
	}
	var imageUpdated []string
	for _, v := range updateImages {
		if d.HasChange(v) && !d.IsNewResource() {
			updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			imageUpdated = append(imageUpdated, v)
			imageUpdate = true
		}
	}
	if d.HasChange("system_disk") && !d.IsNewResource() {
		FlatternStructPrefix(d.Get("system_disk"), &updateReq, "SystemDisk")
		imageUpdate = true
		imageUpdated = append(imageUpdated, "system_disk")
	}
	var initState string
	var needStart bool
	//stop instance
	if imageUpdate || d.HasChange("key_id") || d.HasChange("instance_password") || d.HasChange("host_name") {
		var err error
		initState, err = instanceStop(d, meta)
		if err != nil {
			return err
		}
	}
	if imageUpdate {
		var keyset bool
		keyIds, ok := d.GetOk("key_id")
		if ok {
			keys := SchemaSetToStringSlice(keyIds)
			for k, v := range keys {
				key := fmt.Sprintf("KeyId.%d", k+1)
				updateReq[key] = fmt.Sprintf("%v", v)
			}
			if len(keys) > 0 {
				keyset = true
			}
		}
		keepImageLogin := d.Get("keep_image_login")
		var keepImage bool
		if ok {
			keepImage = keepImageLogin.(bool)
		}
		if (!keepImage) && (!keyset) {
			updateReq["InstancePassword"] = fmt.Sprintf("%v", d.Get("instance_password"))
		}
		action := "ModifyInstanceImage"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := conn.ModifyInstanceImage(&updateReq)
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		if err != nil {
			return fmt.Errorf("error on updating instance image, %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"rebuilding", "overriding"},
			Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"rebuilding", "overriding"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error on waiting for reinstalling instance when ModifyInstanceImage %q, %s", d.Id(), err)
		}
		stateConf = &resource.StateChangeConf{
			Pending: []string{statusPending},
			Target:  []string{"active"},
			//final state may be "stopped" ,need to return error
			Refresh:    instanceStateRefreshForReinstallFunc(conn, d.Id(), []string{"active"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error on waiting for starting instance when ModifyInstanceImage %q, %s", d.Id(), err)
		}
		for _, v := range imageUpdated {
			d.SetPartial(v)
		}
		d.SetPartial("key_id")
	}
	if !imageUpdate && d.HasChange("key_id") && !d.IsNewResource() {
		old, new := d.GetChange("key_id")
		olds := old.(*schema.Set).List()
		if len(olds) > 0 {
			err := instanceDetachKey(d.Id(), olds, conn)
			if err != nil {
				return fmt.Errorf("error instnceDetachKey when ModifyInstanceKey %q, %s", d.Id(), err)
			}
			time.Sleep(time.Second * 10)
			stateConf := &resource.StateChangeConf{
				Pending:    []string{statusPending},
				Target:     []string{"stopped"},
				Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      3 * time.Second,
				MinTimeout: 2 * time.Second,
			}
			if _, err = stateConf.WaitForState(); err != nil {
				return fmt.Errorf("error on waiting  when updateing %q, %s", d.Id(), err)
			}
		}
		news := new.(*schema.Set).List()
		if len(news) > 0 {
			err := instanceAttachKey(d.Id(), news, conn)
			if err != nil {
				return fmt.Errorf("error instnceAttachKey when ModifyInstanceKey %q, %s", d.Id(), err)
			}
			time.Sleep(time.Second * 10)
			stateConf := &resource.StateChangeConf{
				Pending:    []string{statusPending},
				Target:     []string{"stopped"},
				Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      3 * time.Second,
				MinTimeout: 2 * time.Second,
			}
			if _, err = stateConf.WaitForState(); err != nil {
				return fmt.Errorf("error on waiting  when updateing %q, %s", d.Id(), err)
			}
			d.SetPartial("key_id")
		}
		if initState == "active" {
			needStart = true
		}
	}

	if !imageUpdate && d.HasChange("instance_password") && !d.IsNewResource() {
		updatedAttributePassword := []string{
			"instance_password",
		}
		updateReq3 := make(map[string]interface{})
		updateReq3["InstanceId"] = d.Id()
		for _, v := range updatedAttributePassword {
			if d.HasChange(v) && !d.IsNewResource() {
				updateReq3[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			}
		}
		action := "ModifyInstanceAttribute"
		logger.Debug(logger.ReqFormat, action, updateReq3)
		resp, err := conn.ModifyInstanceAttribute(&updateReq3)
		if err != nil {
			return fmt.Errorf("error on updating  instance password, %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq3, *resp)
		//主机updating_password状态变化过快，有可能检测不到，所以暂时去掉状态检测，直接sleep(10s)
	/*	stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"updating_password"},
			Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"updating_password"}),
			Timeout:    *schema.DefaultTimeout(3 * time.Minute),
			Delay:      1 * time.Second,
			MinTimeout: 1 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error on waiting for starting instance when update password %q, %s", d.Id(), err)
		}*/
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"stopped"},
			Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
			Timeout:    *schema.DefaultTimeout(5 * time.Minute),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error on waiting for starting instance when update password%q, %s", d.Id(), err)
		}
		for _, v := range updatedAttributePassword {
			d.SetPartial(v)
		}
		if initState == "active" {
			needStart = true
		}
	}
	if d.HasChange("host_name") && !d.IsNewResource() {
		updatedAttributeHostName := []string{
			"host_name",
		}
		updateReq4 := make(map[string]interface{})
		updateReq4["InstanceId"] = d.Id()
		for _, v := range updatedAttributeHostName {
			if d.HasChange(v) && !d.IsNewResource() {
				updateReq4[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			}
		}
		action := "ModifyInstanceAttribute"
		logger.Debug(logger.ReqFormat, action, updateReq4)
		resp, err := conn.ModifyInstanceAttribute(&updateReq4)
		if err != nil {
			return fmt.Errorf("error on updating  instance host_name, %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq4, *resp)
		for _, v := range updatedAttributeHostName {
			d.SetPartial(v)
		}
		if initState == "active" {
			needStart = true
		}
	}
	if needStart {
		updateReq1 := make(map[string]interface{})
		updateReq1["InstanceId.1"] = d.Id()
		action := "StartInstances"
		logger.Debug(logger.ReqFormat, action, updateReq1)
		resp, err := conn.StartInstances(&updateReq1) //sync
		if err != nil {
			return fmt.Errorf("error on RebootInstances instance(%v) %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, updateReq1, *resp)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"active"},
			Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"active"}),
			Timeout:    *schema.DefaultTimeout(5 * time.Minute),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error on waiting for starting instance when update %q, %s", d.Id(), err)
		}
	}
	//ModifyInstanceType //need reboot
	typeUpdate := false
	updateInstanceTypes := []string{
		"instance_type",
		"data_disk_gb",
	}
	var typeUpdated []string
	for _, v := range updateInstanceTypes {
		if d.HasChange(v) && !d.IsNewResource() {
			updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			typeUpdate = true
			typeUpdated = append(typeUpdated, v)
		}
	}
	if typeUpdate {
		_, err := instanceStop(d, meta)
		if err != nil {
			return err
		}
		action := "ModifyInstanceType"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := conn.ModifyInstanceType(&updateReq)
		if err != nil {
			return fmt.Errorf("error on updating  instance type, %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq, *resp)
		for _, v := range typeUpdated {
			d.SetPartial(v)
		}
	}

	/*
	   //need shutdown later
	   //StopInstances
	     if d.HasChange("force_stop") && !d.IsNewResource() {
	      if d.Get("force_stop").(bool) {
	         _, err := conn.StopInstances(&updateReq1)
	         if err != nil {
	            return fmt.Errorf("error on stop instance , %s", err)
	         }
	      }
	     }


	         if passwordUpdate {
	            action := "ModifyInstanceType"
	            logger.Debug(logger.ReqFormat, action, updateReq3)
	            resp, err := conn.ModifyInstanceType(&updateReq3)
	            if err != nil {
	               return fmt.Errorf("error on updating  instance type, %s", err)
	            }
	            logger.Debug(logger.RespFormat, action, updateReq3, *resp)
	            for _, v := range updatedAttributePassword {
	               d.SetPartial(v)
	            }
	         }

	      if initState == "active" {
	         //judge init state,and need start if active
	         action := "RebootInstances"
	         logger.Debug(logger.ReqFormat, action, updateReq1)
	         resp, err := conn.RebootInstances(&updateReq1) //sync
	         if err != nil {
	            return fmt.Errorf("error on updating instance type, %s", err)
	         }
	         logger.Debug(logger.RespFormat, action, updateReq1, *resp)
	      }
	*/
	//modify network interface information
	var networkUpdate bool
	updateNetworkReq := make(map[string]interface{})
	updateNetworkReq["InstanceId"] = d.Id()
	updateNetworkReq["NetworkInterfaceId"] = fmt.Sprintf("%v", d.Get("network_interface_id"))
	updateNetworkReq["SubnetId"] = fmt.Sprintf("%v", d.Get("subnet_id"))
	allAttributes := []string{
		"private_ip_address",
		"d_n_s1",
		"d_n_s2",
	}
	var updates []string
	for _, v := range allAttributes {
		if d.HasChange(v) {
			networkUpdate = true
			updates = append(updates, v)
		}
	}
	if d.HasChange("subnet_id") && !d.IsNewResource() {
		networkUpdate = true
	}

	if d.HasChange("security_group_id") && !d.IsNewResource() {
		networkUpdate = true
	}
	if networkUpdate {
		if v, ok := d.GetOk("security_group_id"); ok {
			securityGroupIds := SchemaSetToStringSlice(v)
			for k, v := range securityGroupIds {
				updateNetworkReq[fmt.Sprintf("SecurityGroupId.%v", k+1)] = v
			}
		}
		for _, v := range updates {
			if v1, ok := d.GetOk(v); ok {
				updateNetworkReq[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
			}
		}
		action := "ModifyNetworkInterfaceAttribute"
		logger.Debug(logger.ReqFormat, action, updateNetworkReq)
		resp, err := conn.ModifyNetworkInterfaceAttribute(&updateNetworkReq)
		logger.Debug(logger.AllFormat, action, updateNetworkReq, *resp, err)
		if err != nil {
			return fmt.Errorf("update NetworkInterface (%v)error:%v", updateNetworkReq, err)
		}
		result, ok := (*resp)["Return"]
		if !ok || fmt.Sprintf("%v", result) != "true" {
			return fmt.Errorf("update NetworkInterface (%v)error:%v", updateNetworkReq, result)
		}
		d.SetPartial("subnet_id")
		d.SetPartial("security_group_id")
		for _, v := range updates {
			d.SetPartial(v)
		}
	}
	d.Partial(false)
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	//delete
	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId.1"] = d.Id()
	return resource.Retry(30*time.Minute, func() *resource.RetryError {
		readReq := make(map[string]interface{})
		readReq["InstanceId.1"] = d.Id()
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err1 := conn.DescribeInstances(&readReq)
		logger.Debug(logger.AllFormat, action, readReq, *resp, err1)
		/*
		   {
		       "Marker": 0,
		       "InstanceCount": 0,
		       "RequestId": "7c34a18b-b562-44f3-8ea9-a350e3afe649",
		       "InstancesSet": []
		   }
		*/
		if err1 != nil && notFoundError(err1) {
			return nil
		}
		if err1 != nil {
			return resource.NonRetryableError(err1)
		}
		log.Printf("force_delete:%v", d.Get("force_delete"))
		forceDelete, ok := d.GetOk("force_delete")
		if !ok {
			forceDelete = false
		}
		//recycling will cant't delete security group
		/*else { //recycling will cant't delete security group
		     if err2 == nil || notFoundError(err2) {
		        return nil
		     }
		  }
		*/
		itemset, ok := (*resp)["InstancesSet"]
		if !ok {
			return nil
		}
		items, _ := itemset.([]interface{})
		if len(items) == 0 && fmt.Sprintf("%v", forceDelete) != "true" {
			return nil
		}
		if len(items) > 0 { //not in recycling
			action = "TerminateInstances"
			logger.Debug(logger.ReqFormat, action, deleteReq)
			resp, err2 := conn.TerminateInstances(&deleteReq)
			logger.Debug(logger.AllFormat, action, deleteReq, *resp, err2)
			if err2 != nil && inUseError(err2) {
				return resource.RetryableError(err2)
			}
		}

		/*
		   if strings.ToLower(initState) != "stopped" &&
		      strings.ToLower(initState) != "error" &&
		      strings.ToLower(initState) != "recycling" {
		      action = "StopInstances"
		      logger.Debug(logger.ReqFormat, action, readReq)
		      resp, err := conn.StopInstances(&readReq) //sync
		      logger.Debug(logger.AllFormat, action, readReq, *resp, err)
		      if err1 != nil && notFoundError(err1) {
		         return nil
		      }
		      if err != nil {
		         return resource.RetryableError(err)
		      }
		      stateConf := &resource.StateChangeConf{
		         Pending:    []string{statusPending},
		         Target:     []string{"stopped"},
		         Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
		         Timeout:    d.Timeout(schema.TimeoutUpdate),
		         Delay:      3 * time.Second,
		         MinTimeout: 2 * time.Second,
		      }
		      if _, err = stateConf.WaitForState(); err != nil {
		         return resource.RetryableError(err)
		      }
		   }

		*/

		//check
		readReq = make(map[string]interface{})
		readReq["InstanceId.1"] = d.Id()
		readReq["Filter.1.Name"] = "instance-state.name"
		readReq["Filter.1.Value.1"] = "recycling"
		action = "DescribeInstances"
		index := 1
		for {
			if index > 32 {
				return resource.NonRetryableError(fmt.Errorf("no data on geting recycling kec when delete %q", d.Id()))
			}
			time.Sleep(time.Second * time.Duration(index))
			index = index * 2

			logger.Debug(logger.ReqFormat, action, readReq)
			resp, err := conn.DescribeInstances(&readReq)
			logger.Debug(logger.AllFormat, action, readReq, *resp)
			if err != nil && notFoundError(err) {
				return nil
			}
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error on  reading kec when delete %q, %s", d.Id(), err))
			}
			itemset, ok = (*resp)["InstancesSet"]
			if !ok {
				return nil
			}
			items, _ = itemset.([]interface{})
			if len(items) == 0 {
				return nil
			}
			state, ok := items[0].(map[string]interface{})["InstanceState"]
			if !ok {
				return nil
			}
			initState, ok := state.(map[string]interface{})["Name"].(string)
			if !ok {
				return nil
			}
			if initState == "recycling" {
				break
			}
		}

		//deleteReq["IsRefundResouce"] = "false"
		deleteReq["ForceDelete"] = fmt.Sprintf("%v", forceDelete)
		action = "TerminateInstances"
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err3 := conn.TerminateInstances(&deleteReq)
		logger.Debug(logger.AllFormat, action, deleteReq, *resp, err3)
		if err3 == nil || notFoundError(err3) {
			return nil
		}
		if err3 != nil && inUseError(err3) {
			return resource.RetryableError(err3)
		}
		return resource.NonRetryableError(err3)
	})
}
func instanceStateRefreshFunc(client *kec.Kec, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"InstanceId.1": instanceId}
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeInstances(&req)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		itemset := (*resp)["InstancesSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		item, ok1 := items[0].(map[string]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		instanceState, ok2 := item["InstanceState"]
		if !ok2 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		instancestate, ok3 := instanceState.(map[string]interface{})
		if !ok3 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		state := strings.ToLower(instancestate["Name"].(string))
		for k, v := range target {
			if v == state {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func instanceStateRefreshForReinstallFunc(client *kec.Kec, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"InstanceId.1": instanceId}
		resp, err := client.DescribeInstances(&req)
		if err != nil {
			return nil, "", err
		}
		itemset := (*resp)["InstancesSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		item, ok1 := items[0].(map[string]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		instanceState, ok2 := item["InstanceState"]
		if !ok2 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		instancestate, ok3 := instanceState.(map[string]interface{})
		if !ok3 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		state := strings.ToLower(instancestate["Name"].(string))
		if state == "stopped" {
			return nil, "", fmt.Errorf("instance restart error")
		}
		for k, v := range target {
			if v == state {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func instanceStateRefreshForCreateFunc(client *kec.Kec, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"InstanceId.1": instanceId}
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeInstances(&req)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		itemset := (*resp)["InstancesSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		item, ok1 := items[0].(map[string]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		instanceState, ok2 := item["InstanceState"]
		if !ok2 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		instancestate, ok3 := instanceState.(map[string]interface{})
		if !ok3 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		state := strings.ToLower(instancestate["Name"].(string))
		if state == "error" {
			return nil, "", fmt.Errorf("instance create error")
		}
		for k, v := range target {
			if v == state {
				//resp cant't be null else will try
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func instanceDetachKey(instanceId string, keyIds []interface{}, conn *kec.Kec) error {
	req := make(map[string]interface{})
	req["InstanceId.1"] = fmt.Sprintf("%v", instanceId)
	for k, v := range keyIds {
		req[fmt.Sprintf("KeyId.%v", k+1)] = fmt.Sprintf("%v", v)
	}
	action := "DetachKey"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DetachKey(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("Error DetachKey : %s", err)
	}
	instanceSet, ok := (*resp)["InstancesSet"]
	if !ok {
		return fmt.Errorf("Error DetachKey1 ")
	}
	instances, ok := instanceSet.([]interface{})
	if !ok {
		return fmt.Errorf("Error DetachKey2 ")
	}
	instance, ok := instances[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error DetachKey3 ")
	}
	status, ok := instance["Return"]
	if !ok {
		return fmt.Errorf("Error DetachKey4 ")
	}
	status1, ok := status.(bool)
	if !ok || !status1 {
		return fmt.Errorf("Error DetachKey:fail ")
	}
	return nil
}
func instanceAttachKey(instanceId string, keyIds []interface{}, conn *kec.Kec) error {
	req := make(map[string]interface{})
	req["InstanceId.1"] = fmt.Sprintf("%v", instanceId)
	for k, v := range keyIds {
		req[fmt.Sprintf("KeyId.%v", k+1)] = fmt.Sprintf("%v", v)
	}
	action := "AttachKey"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.AttachKey(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("Error AttachKey : %s", err)
	}
	instanceSet, ok := (*resp)["InstancesSet"]
	if !ok {
		return fmt.Errorf("Error AttachKey1 ")
	}
	instances, ok := instanceSet.([]interface{})
	if !ok {
		return fmt.Errorf("Error AttachKey2 ")
	}
	instance, ok := instances[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error AttachKey3 ")
	}
	status, ok := instance["Return"]
	if !ok {
		return fmt.Errorf("Error AttachKey4 ")
	}
	status1, ok := status.(bool)
	if !ok || !status1 {
		return fmt.Errorf("Error AttachKey:fail ")
	}
	return nil
}
func instanceStop(d *schema.ResourceData, meta interface{}) (string, error) {
	conn := meta.(*KsyunClient).kecconn
	readReq := make(map[string]interface{})
	readReq["InstanceId.1"] = d.Id()
	action := "DescribeInstances"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeInstances(&readReq)
	logger.Debug(logger.AllFormat, action, readReq, *resp, err)
	if err != nil {
		return "", fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}

	itemset := (*resp)["InstancesSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		return "", fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}
	state := items[0].(map[string]interface{})["InstanceState"]
	initState := state.(map[string]interface{})["Name"].(string)
	if initState == "error" {
		return initState, fmt.Errorf("instance with error state")
	}
	if initState != "stopped" {
		action = "StopInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err = conn.StopInstances(&readReq) //同步
		logger.Debug(logger.AllFormat, action, readReq, *resp, err)
		if err != nil {
			return initState, fmt.Errorf("error on stop  instance %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"stopped"},
			Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return initState, fmt.Errorf("error on waiting for starting instance when stopping %q, %s", d.Id(), err)
		}
	}
	return initState, nil
}
