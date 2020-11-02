---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_security_groups"
sidebar_current: "docs-ksyun-datasource-redis-security-groups"
description: |-
  Provides a list of Redis security group resources in the current region.
---

# ksyun_redis_security_groups

This data source provides a list of redis security group resources according to their security Group Id, name, description they belong to .

## Example Usage

```hcl
# Get redis security groups
data "ksyun_redis_security_groups" "default" {
  output_file       = "output_result1"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type which documented below.
* `total_count` - Total number of Redis security groups  that satisfy the condition.

