package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunEpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEpcCreate,
		Update: resourceKsyunEpcUpdate,
		Read:   resourceKsyunEpcRead,
		Delete: resourceKsyunEpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateName,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"raid": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Raid1",
					"Raid10",
					"Raid5",
					"Raid50",
					"SRaid0",
				}, false),
			},

			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"network_interface_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bond4",
					"single",
					"dual",
				}, false),
			},

			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"private_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},

			"security_group_id": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dns1": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},

			"dns2": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},

			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Monthly",
					"Daily",
					"PostPaidByDay",
					"PrePaidByMonth",
				}, false),
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"purchase_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 36),
			},

			"sn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstancePassword,
			},

			"security_agent": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "classic",
			},

			"cloud_monitor_agent": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "classic",
			},

			"extension_network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"extension_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"extension_private_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},

			"extension_security_group_id": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"extension_dns1": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},

			"extension_dns2": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},

			"address_band_width": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"line_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"address_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PrePaidByMonth",
					"PostPaidByPeak",
					"PostPaidByDay",
					"PostPaidByTransfer",
					"PostPaidByHour",
				}, false),
			},

			"address_purchase_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"address_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"system_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EXT4",
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
			},

			"data_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "XFS",
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
			},

			"data_disk_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/DATA/disk",
				ValidateFunc: validation.StringInSlice([]string{
					"/DATA/disk",
					"/data",
				}, false),
			},

			"data_disk_catalogue_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NaturalNumber",
				ValidateFunc: validation.StringInSlice([]string{
					"NoSuffix",
					"NaturalNumber",
					"NaturalNumberFromZero",
				}, false),
			},

			"hyper_threading": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NoChange",
				ValidateFunc: validation.StringInSlice([]string{
					"Open",
					"Close",
					"NoChange",
				}, false),
			},

			"hyper_threading_status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"allow_modify_hyper_threading": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"releasable_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cabinet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"rack_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cabinet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_bond": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"product_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"memory": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_container": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"cpu": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequence": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"gpu": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequence": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"disk_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_space": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_attribute": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network_interface_attribute_set": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"d_n_s1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"d_n_s2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceKsyunEpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.epcconn

	var resp *map[string]interface{}
	var err error
	createReq := make(map[string]interface{})
	createParam := []string{
		"availability_zone",
		"host_type",
		"raid",
		"image_id",
		"network_interface_mode",
		"subnet_id",
		"private_ip_address",
		"dns1",
		"dns2",
		"key_id",
		"host_name",
		"charge_type",
		"project_id",
		"purchase_time",
		"password",
		"security_agent",
		"cloudmonitor_agent",
		"extension_subnet_id",
		"extension_private_ip_address",
		"extension_dns1",
		"extension_dns2",
		"address_band_width",
		"line_id",
		"address_charge_type",
		"address_purchase_time",
		"address_project_id",
		"system_file_type",
		"data_file_type",
		"data_disk_catalogue",
		"data_disk_catalogue_suffix",
		"hyper_threading",
	}

	for _, v := range createParam {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createReq[vv] = fmt.Sprintf("%v", v1)
		}
	}

	createStruct := []string{
		"security_group_id",
		"extension_security_group_id",
	}

	for _, v := range createStruct {
		if v1, ok := d.GetOk(v); ok {
			keys := SchemaSetToStringSlice(v1)
			prefix := Downline2Hump(v)
			for k, v := range keys {
				key := fmt.Sprintf("%s.%d", prefix, k+1)
				createReq[key] = fmt.Sprintf("%v", v)
			}
		}
	}

	action := "CreateEpc"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.CreateEpc(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating Epc, %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)

	if resp != nil {
		epc := (*resp)["Host"].(map[string]interface{})
		d.SetId(epc["HostId"].(string))
	}
	// after create instance, we need to wait it initialized
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Initializing"},
		Target:     []string{"Running"},
		Refresh:    epcInstanceStateRefreshFunc(conn, d.Id(), "Running"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on waiting for epc %s complete creating, %s", d.Id(), err)
	}
	return resourceKsyunEpcRead(d, meta)
}

func resourceKsyunEpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.epcconn

	readReq := make(map[string]interface{})
	readReq["HostId.1"] = d.Id()

	action := "DescribeEpcs"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeEpcs(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading Epc %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)

	if resp != nil {
		items := (*resp)["HostSet"].([]interface{})
		if len(items) == 0 {
			d.SetId("")
			return nil
		}
		item := items[0].(map[string]interface{})

		excludesKeys := map[string]bool{
			"Cpu":                          true,
			"Gpu":                          true,
			"DiskSet":                      true,
			"NetworkInterfaceAttributeSet": true,
		}
		excludes := SetDByResp(d, item, epcInstanceKeys, excludesKeys)

		if excludes["Cpu"] != nil {
			itemSet := GetSubDByRep(excludes["Cpu"], epcCpuKeys, map[string]bool{})
			d.Set(Hump2Downline("Cpu"), itemSet)
		}
		if excludes["Gpu"] != nil {
			itemSet := GetSubDByRep(excludes["Gpu"], epcGpuKeys, map[string]bool{})
			d.Set(Hump2Downline("Gpu"), itemSet)
		}
		if excludes["DiskSet"] != nil {
			itemSet := GetSubSliceDByRep(excludes["DiskSet"].([]interface{}), epcDiskSetKeys)
			d.Set(Hump2Downline("DiskSet"), itemSet)
		}
		if excludes["NetworkInterfaceAttributeSet"] != nil {
			itemSet := GetSubSliceDByRep(excludes["NetworkInterfaceAttributeSet"].([]interface{}), epcNetworkInterfaceKeys)
			for k, v := range itemSet {
				if sg, ok := v["security_group_set"]; ok {
					itemSetSub := GetSubSliceDByRep(sg.([]interface{}), epcSecurityGroupKeys)
					itemSet[k]["security_group_set"] = itemSetSub
				}
			}
			d.Set(Hump2Downline("NetworkInterfaceAttributeSet"), itemSet)

			//为输入网卡属性赋值 方便使用
			for _, v := range itemSet {
				if v1, ok := v["network_interface_type"]; ok {
					if "primary" == v1.(string) {
						d.Set("network_interface_id", v["network_interface_id"])
						d.Set("subnet_id", v["subnet_id"])
						d.Set("private_ip_address", v["private_ip_address"])
						d.Set("security_group_id", v["security_group_id"])
						d.Set("dns1", v["d_n_s1"])
						d.Set("dns2", v["d_n_s2"])
						if sg, ok := v["security_group_set"]; ok {
							sgids := []string{}
							for _, v2 := range sg.([]map[string]interface{}) {
								sgids = append(sgids, v2["security_group_id"].(string))
							}
							d.Set("security_group_id", sgids)
						}
					} else {
						d.Set("extension_network_interface_id", v["network_interface_id"])
						d.Set("extension_subnet_id", v["subnet_id"])
						d.Set("extension_private_ip_address", v["private_ip_address"])
						d.Set("extension_security_group_id", v["security_group_id"])
						d.Set("extension_dns1", v["d_n_s1"])
						d.Set("extension_dns2", v["d_n_s2"])
						if sg, ok := v["security_group_set"]; ok {
							sgids := []string{}
							for _, v2 := range sg.([]map[string]interface{}) {
								sgids = append(sgids, v2["security_group_id"].(string))
							}
							d.Set("extension_security_group_id", sgids)
						}
					}
				}
			}
		}
	}
	return nil
}

func resourceKsyunEpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.epcconn
	//attributeUpdate := false
	d.Partial(true)
	updateReq := make(map[string]interface{})
	updateReq["HostId"] = d.Id()

	action := "ModifyEpc"
	//修改 物理机名称
	if d.HasChange("host_name") && !d.IsNewResource() {
		updateReq["HostName"] = d.Get("host_name").(string)
		logger.Debug(logger.ReqFormat, action, updateReq)
		_, err := conn.ModifyEpc(&updateReq)
		if err != nil {
			return fmt.Errorf("error on modify Epc, %s", err)
		}
		d.SetPartial("host_name")
		//清除参数
		clearUpdateReq(&updateReq)
	}
	//修改DNS  目前支持主网卡
	if d.HasChange("dns1") && !d.IsNewResource() {
		dns1, ok := d.GetOk("dns1")
		networkInterfaceId, ok1 := d.GetOk("network_interface_id")
		if ok && ok1 {
			updateReq["DNS1"] = dns1
			updateReq["NetworkInterfaceId"] = networkInterfaceId
			if d.HasChange("dns2") {
				updateReq["DNS2"] = d.Get("dns2")
			}
			action := "ModifyDns"
			logger.Debug(logger.ReqFormat, action, updateReq)
			_, err := conn.ModifyDns(&updateReq)
			if err != nil {
				return fmt.Errorf("error on modifydns Epc, %s", err)
			}
			d.SetPartial("dns1")
			d.SetPartial("dns2")
			//清除参数
			clearUpdateReq(&updateReq)
		}

	}

	//修改子网 目前支持主网卡
	if d.HasChange("subnet_id") && !d.IsNewResource() {
		subnetId, ok := d.GetOk("subnet_id")
		networkInterfaceId, ok1 := d.GetOk("network_interface_id")
		if ok && ok1 {
			updateReq["SubnetId"] = subnetId
			updateReq["NetworkInterfaceId"] = networkInterfaceId
			if d.HasChange("ip_address") {
				updateReq["IpAddress"] = d.Get("ip_address")
			}
			if d.HasChange("security_group_id") {
				keys := SchemaSetToStringSlice(d.Get("security_group_id"))
				prefix := Downline2Hump("security_group_id")
				for k, v := range keys {
					key := fmt.Sprintf("%s.%d", prefix, k+1)
					updateReq[key] = fmt.Sprintf("%v", v)
				}
			}
			action := "ModifyNetworkInterfaceAttribute"
			logger.Debug(logger.ReqFormat, action, updateReq)
			_, err := conn.ModifyNetworkInterfaceAttribute(&updateReq)
			if err != nil {
				return fmt.Errorf("error on modifydns Epc, %s", err)
			}
			d.SetPartial("subnet_id")
			d.SetPartial("ip_address")
			d.SetPartial("security_group_id")
			//清除参数
			clearUpdateReq(&updateReq)
		}
	} else {
		//修改安全组
		if d.HasChange("security_group_id") && !d.IsNewResource() {
			_, ok := d.GetOk("security_group_id")
			networkInterfaceId, ok1 := d.GetOk("network_interface_id")
			if ok && ok1 {
				updateReq["NetworkInterfaceId"] = networkInterfaceId
				keys := SchemaSetToStringSlice(d.Get("security_group_id"))
				prefix := Downline2Hump("security_group_id")
				for k, v := range keys {
					key := fmt.Sprintf("%s.%d", prefix, k+1)
					updateReq[key] = fmt.Sprintf("%v", v)
				}
				action := "ModifySecurityGroup"
				logger.Debug(logger.ReqFormat, action, updateReq)
				_, err := conn.ModifySecurityGroup(&updateReq)
				if err != nil {
					return fmt.Errorf("error on modifydns Epc, %s", err)
				}
				d.SetPartial("security_group_id")
				d.SetPartial("network_interface_id")
				//清除参数
				clearUpdateReq(&updateReq)
			}
		}
	}

	//重装系统
	if d.HasChange("image_id") && !d.IsNewResource() {
		err := operateEpcStatusForWait(d, "Reinstalling", updateReq, conn)
		if err != nil {
			return err
		}
	}
	d.Partial(false)
	return resourceKsyunEpcRead(d, meta)
}

func operateEpcStatusForWait(d *schema.ResourceData, hostStatus string, updateReq map[string]interface{}, conn *epc.Epc) error {
	_, err := operateEpcStatus(d, hostStatus, updateReq, conn)
	if err != nil {
		return fmt.Errorf("error on operateEpcStatus Epc, %s", err)
	}
	if hostStatus != "Stopped" {
		hostStatus = "Running"
	}
	// after create instance, we need to wait it initialized
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Initializing"},
		Target:     []string{hostStatus},
		Refresh:    epcInstanceStateRefreshFunc(conn, d.Id(), hostStatus),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on waiting for epc %s complete creating, %s", d.Id(), err)
	}
	return nil
}

func operateEpcStatus(d *schema.ResourceData, hostStatus string, updateReq map[string]interface{}, conn *epc.Epc) (*map[string]interface{}, error) {
	switch hostStatus {
	case "Stopped":
		action := "StopEpc"
		logger.Debug(logger.ReqFormat, action, updateReq)
		return conn.StopEpc(&updateReq)
	case "Running":
		action := "StartEpc"
		logger.Debug(logger.ReqFormat, action, updateReq)
		return conn.StartEpc(&updateReq)
	case "Rebooting":
		action := "RebootEpc"
		logger.Debug(logger.ReqFormat, action, updateReq)
		return conn.RebootEpc(&updateReq)
	case "StartHyperThreading":
		updateReq["HyperThreadingStatus"] = hostStatus
		action := "ModifyHyperThreading"
		logger.Debug(logger.ReqFormat, action, updateReq)
		return conn.ModifyHyperThreading(&updateReq)
	case "StopHyperThreading":
		updateReq["HyperThreadingStatus"] = hostStatus
		action := "ModifyHyperThreading"
		logger.Debug(logger.ReqFormat, action, updateReq)
		return conn.ModifyHyperThreading(&updateReq)
	case "ResetPassword":
		if v, ok := d.GetOk("password"); ok {
			if d.HasChange("password") {
				updateReq["Password"] = v
				action := "ResetPassword"
				logger.Debug(logger.ReqFormat, action, updateReq)
				resp, err := conn.ResetPassword(&updateReq)
				if err == nil {
					d.SetPartial("password")
				}
				return resp, err
			}
		}
	case "Reinstalling":
		updateReq["Password"] = d.Get("password")
		updateReq["ImageId"] = d.Get("image_id")
		updateReq["KeyId"] = d.Get("key_id")
		updateReq["NetworkInterfaceMode"] = d.Get("network_interface_mode")
		if v, ok := d.GetOk("security_agent"); ok {
			updateReq["SecurityAgent"] = v
		}
		if v, ok := d.GetOk("cloudmonitor_agent"); ok {
			updateReq["CloudMonitorAgent"] = v
		}
		updateReq["Raid"] = d.Get("raid")
		updateReq["SystemFileType"] = d.Get("system_file_type")
		updateReq["DataFileType"] = d.Get("data_file_type")
		updateReq["DataDiskCatalogue"] = d.Get("data_disk_catalogue")
		updateReq["DataDiskCatalogueSuffix"] = d.Get("data_disk_catalogue_suffix")
		updateReq["HyperThreading"] = d.Get("hyper_threading")
		action := "ReinstallEpc"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := conn.ReinstallEpc(&updateReq)
		if err == nil {
			d.SetPartial("password")
			d.SetPartial("image_id")
			d.SetPartial("key_id")
			d.SetPartial("network_interface_mode")
			d.SetPartial("security_agent")
			d.SetPartial("cloudmonitor_agent")
			d.SetPartial("raid")
			d.SetPartial("system_file_type")
			d.SetPartial("data_file_type")
			d.SetPartial("data_disk_catalogue")
			d.SetPartial("data_disk_catalogue_suffix")
			d.SetPartial("hyper_threading")
		}
		return resp, err
	}
	return nil, nil
}

//清除除HostId其他的属性
func clearUpdateReq(updateReq *map[string]interface{}) {
	if *updateReq != nil {
		for k := range *updateReq {
			if k != "HostId" {
				delete(*updateReq, k)
			}
		}
	}

}

func resourceKsyunEpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.epcconn
	deleteEpc := make(map[string]interface{})
	deleteEpc["HostId"] = d.Id()

	_, err := conn.DeleteEpc(&deleteEpc)
	if err != nil {
		return fmt.Errorf("error on deleting Epc, %s", err)
	}
	return nil
}

func epcInstanceStateRefreshFunc(conn *epc.Epc, instanceId, target string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		readEpc := make(map[string]interface{})
		readEpc["HostId.1"] = instanceId
		resp, err := conn.DescribeEpcs(&readEpc)
		if err != nil {
			if isNotFoundError(err) {
				return nil, "Initializing", nil
			}
			return nil, "Initializing", err
		}

		if resp != nil {
			l := (*resp)["HostSet"].([]interface{})
			if len(l) == 0 {
				return nil, "Initializing", nil
			}
			item := l[0].(map[string]interface{})
			state := item["HostStatus"].(string)
			if state != target {
				state = "Initializing"
			}
			return l[0], state, nil
		}
		return nil, "Initializing", err
	}
}
