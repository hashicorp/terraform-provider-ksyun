---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_instance"
sidebar_current: "docs-ksyun-resource-redis-instance"
description: |-
  Provides an Redis instance resource.
---

# ksyun_redis_instance

Provides an redis instance resource.

## Example Usage

```hcl
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
  available_zone = "${var.available_zone}"
}

resource "ksyun_redis_instance" "default" {
  available_zone        = "${var.available_zone}"
  name                  = "MyRedisInstance1101"
  mode                  = 2
  capacity              = 1
  slaveNum              = 2  
  net_type              = 2
  vnet_id               = "${ksyun_subnet.default.id}"
  vpc_id                = "${ksyun_vpc.default.id}"
  bill_type             = 5
  duration              = ""
  duration_unit         = ""
  pass_word             = "Shiwo1101"
  iam_project_id        = "0"
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "513",
    "zset-max-ziplist-entries"    = "129",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Optional) The Zone to launch the DB instance.
* `name ` - (Optional) The name of DB instance.
* `mode ` - (Optional) The KVStore instance system architecture required by the user. Valid values:  1(cluster),2(single).
* `capacity ` - (Require) The instance capacity required by the user. Valid values :{1, 2, 4, 8, 16,20,24,28, 32, 64}.
* `slaveNum ` - (Optional) The readonly node num required by the user. Valid values ：{0-7}
* `net_type ` - (Require) The network type. Valid values ：2(vpc).
* `vpc_id` - (Require)   Used to retrieve instances belong to specified VPC .
* `vnet_id` - (Require) The ID of subnet. the instance will use the subnet in the current region.
* `bill_type` - (Optional)Valid values are 1 (Monthly), 5(Daily), 87(HourlyInstantSettlement).
* `duration` - (Optional)Only meaningful if bill_type is 1。 Valid values：{1~36}.
* `duration_unit` - (Optional)Only meaningful if bill_type is 1。 Valid values：month.
* `pass_word` - (Optional)The password of the  instance.The password is a string of 8 to 30 characters and must contain uppercase letters, lowercase letters, and numbers.
* `iam_project_id` - (Optional) The project instance belongs to.
* `protocol` - Engine version. Supported values: 2.8, 4.0 and 5.0.
* `parameters` - Set of parameters needs to be set after instance was launched. Available parameters can refer to the  docs https://docs.ksyun.com/documents/1018 .


