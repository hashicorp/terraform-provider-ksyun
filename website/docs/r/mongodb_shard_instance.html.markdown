---
layout: "ksyun"
page_title: "Ksyun: ksyun_mongodb_shard_instance"
sidebar_current: "docs-ksyun-resource-mongodb-shard-instance"
description: |-
  Provides an shard MongoDB resource.
---

# ksyun_mongodb_shard_instance

Provides an shard MongoDB resource.

## Example Usage

```hcl
resource "ksyun_mongodb_shard_instance" "default" {
  name = "InstanceName"
  instance_account = "root"
  instance_password = "admin"
  mongos_class = "1C2G"
  mongos_num = 2
  shard_class = "1C2G"
  shard_num = 2
  storage = 5
  vpc_id = "VpcId"
  vnet_id = "VnetId"
  db_version = "3.6"
  pay_type = "hourlyInstantSettlement"
  iam_project_id = "0"
  availability_zone = "cn-shanghai-3b"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of instance, which contains 6-64 characters and only support Chinese, English, numbers, '-', '_'.
* `instance_account` - (Optional) The administrator name of instance, if not defined `instance_account`, the instance will use `root`.
* `instance_password` - (Required) The administrator password of instance.
* `mongos_class` - (Required) The class of instance mongo node cpu and memory.
* `mongos_num` - (Required) The num of instance mongo node.
* `shard_class` - (Required) The class of instance shard node cpu and memory.
* `shard_num` - (Required) The num of instance shard node.
* `storage` - (Required) The size of instance disk, measured in GB (GigaByte).
* `vpc_id` - (Required) The id of VPC linked to the instance.
* `vnet_id` - (Required) The id of subnet linked to the instance.
* `db_version` - (Required) The version of instance engine, and support `3.2` and `3.6`
* `pay_type` - (Optional) Instance charge type, if not defined `pay_type`, the instance will use `byMonth`.
* `duration` - (Optional) The duration of instance use, if `pay_type` is `byMonth`, the duration is required.
* `iam_project_id` - (Optional) The project id of instance belong, if not defined `iam_project_id`, the instance will use `0`.
* `availability_zone` - (Required) Availability zone where instance is located.