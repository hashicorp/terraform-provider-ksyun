---
layout: "ksyun"
page_title: "Ksyun: ksyun_mongodbs"
sidebar_current: "docs-ksyun-datasource-mongodbs"
description: |-
  Provides a list of MongoDB resources in the current region.
---

# ksyun_mongodbs

This data source provides a list of MongoDB resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to .

## Example Usage

```hcl
# Get  mongodbs
data "ksyun_mongodbs" "default" {
  output_file = "output_result"
  iam_project_id = ""
  instance_id = ""
  vnet_id = ""
  vpc_id = ""
  name = ""
  vip = ""
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional)  The name of MongoDB, all the MongoDBs belong to this region will be retrieved if the name is `""`.
* `instance_id` - (Optional)  The id of MongoDB, all the MongoDBs belong to this region will be retrieved if the instance_id is `""`.
* `iam_project_id` - (Optional)  The project instance belongs to.
* `vpc_id` - (Optional)   Used to retrieve instances belong to specified VPC .
* `vnet_id` - (Optional) The ID of subnet. the instance will use the subnet in the current region.
* `vip` - (Optional) The vip of instances. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type which documented below.
* `total_count` - Total number of MongoDBs that satisfy the condition.

