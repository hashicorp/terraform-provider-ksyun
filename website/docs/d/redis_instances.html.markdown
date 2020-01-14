---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_instances"
sidebar_current: "docs-ksyun-datasource-redis-instances"
description: |-
  Provides a list of Redis resources in the current region.
---

# ksyun_redis_instances

This data source provides a list of redis resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to .

## Example Usage

```hcl
# Get  redis instances
data "ksyun_redis_instances" "default" {
  output_file       = "output_result1"
  fuzzy_search      = ""
  iam_project_id    = ""
  cache_id          = ""
  vnet_id           = ""
  vpc_id            = ""
  name              = ""
  vip               = ""
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional)  The name of redis instance, all the Redis instances belong to this region will be retrieved if the name is `""`.
* `iam_project_id` - (Optional)  The project instance belongs to.
* `cache_id` - (Optional)  The ID of  the intance .
* `vpc_id` - (Optional)   Used to retrieve instances belong to specified VPC .
* `vnet_id` - (Optional) The ID of subnet. the instance will use the subnet in the current region.
* `vip` - (Optional) Private IP address of the instance. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type which documented below.
* `total_count` - Total number of Redis instances that satisfy the condition.

