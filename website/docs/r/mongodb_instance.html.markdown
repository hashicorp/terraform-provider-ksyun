---
layout: "ksyun"
page_title: "Ksyun: ksyun_mongodb_instance"
sidebar_current: "docs-ksyun-resource-mongodb-instance"
description: |-
  Provides an replica set MongoDB resource.
---

# ksyun_mongodb_instance

Provides an replica set MongoDB resource.

## Example Usage

```hcl
resource "ksyun_mongodb_instance" "default" {
  name = "InstanceName"
  instance_account = "root"
  instance_password = "admin"
  instance_class = "1C2G"
  storage = 5
  node_num = 3
  vpc_id = "VpcId"
  vnet_id = "VnetId"
  db_version = "3.6"
  pay_type = "byDay"
  iam_project_id = "0"
  availability_zone = "cn-shanghai-3b"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of instance, which contains 6-64 characters and only support Chinese, English, numbers, '-', '_'.
* `instance_account` - (Optional) The administrator name of instance, if not defined `instance_account`, the instance will use `root`.
* `instance_password` - (Required) The administrator password of instance.
* `instance_class` - (Required) The class of instance cpu and memory.
* `storage` - (Required) The size of instance disk, measured in GB (GigaByte).
* `node_num` - (Required) The num of instance node.
* `vpc_id` - (Required) The id of VPC linked to the instance.
* `vnet_id` - (Required) The id of subnet linked to the instance.
* `db_version` - (Required) The version of instance engine, and support `3.2` and `3.6`
* `pay_type` - (Optional) Instance charge type, if not defined `pay_type`, the instance will use `byMonth`.
* `duration` - (Optional) The duration of instance use, if `pay_type` is `byMonth`, the duration is required.
* `iam_project_id` - (Optional) The project id of instance belong, if not defined `iam_project_id`, the instance will use `0`.
* `availability_zone` - (Required) Availability zone where instance is located.


