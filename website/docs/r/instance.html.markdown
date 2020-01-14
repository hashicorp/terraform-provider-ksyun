---
layout: "ksyun"
page_title: "Ksyun: ksyun_instance"
sidebar_current: "docs-ksyun-resource-instance"
description: |-
  Provides an Host Instance resource.
---


# ksyun_instance

Provides a KEC instance resource.

**Note**  At present, 'Monthly' instance cannot be deleted and must wait it to be outdated and released automatically.

## Example Usage

```h
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}

data "ksyun_lines" "default" {
  output_file=""
  line_name="BGP"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "ksyun_subnet" "default" {
  subnet_name      = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="${var.security_group_name}"
}

resource "ksyun_security_group_entry" "default" {
  description = "test1"
  security_group_id="${ksyun_security_group.default.id}"
  cidr_block="10.0.1.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}

resource "ksyun_ssh_key" "default" {
  key_name="ssh_key_tf"
  public_key=""
}

resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  data_disk =[
    {
      type="SSD3.0"
      size=20
      delete_with_instance=true
    }
  ]
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  private_ip_address=""
  instance_name="xuan-tf-combine"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=["${ksyun_ssh_key.default.id}"]
  force_delete=true
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The ID for the image to use for the instance.
* `instance_type` -  (Required) The type of instance to start.
* `system_disk` - (Required) System disk parameters.
    - `disk_type` - System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk.
    - `disk_size` - The size of the data disk.
* `data_disk_gb` - (Optional) The local SSD disk.
* `data_disk` - (Optional) The list of data disks created with instance.
    - `type` - Data disk type. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk.
    - `size` - Data disk type size.
    - `delete_with_instance` -  Delete this data disk when the instance is destroyed. It only works on SSD3.0, EHDD, disk.
* `subnet_id` - (Required) The ID of subnet. the instance will use the subnet in the current region.
* `security_group_id` - (Required) Security Group to associate with.
* `instance_password` - (Optional) Password to an instance is a string of 8 to 32 characters. 
* `instance_name` - (Optional) The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.
* `keep_image_login` - (Optional) Keep the initial settings of the custom image.
* `charge_type` - (Required) Valid values are Monthly, Daily, HourlyInstantSettlement.
* `purchase_time` - (Optional) The duration that you will buy the resource.
* `private_ip_address` - (Optional) Instance private IP address can be specified when you creating new instance.
* `sriov_net_support` (Optional) Network enhancement.
* `project_id` - (Optional) The project instance belongs to.
* `user_data` - (Optional) The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `creation_date` - The time of creation for instance, formatted in ISO8601 time string.
* `instance_state` - Instance current status. Possible values are `active`, `building`, `stopped`, `deleting`.


## Import

Instance can be imported using the `id`, e.g.

```
$ terraform import 
```