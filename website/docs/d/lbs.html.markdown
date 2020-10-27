---
layout: "ksyun"
page_title: "Ksyun: ksyun_lbs"
sidebar_current: "docs-ksyun-datasource-lbs"
description: |-
  Provides a list of Load Balancer resources in the current region.
---

# ksyun_lbs

This data source provides a list of Load Balancer resources according to their Load Balancer ID, VPC ID and Subnet ID.

## Example Usage

```hcl
data "ksyun_lbs" "default" {
  output_file="output_result"
  name_regex=""
  ids=[]
  state=""
  vpc_id=[]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Load Balancer IDs, all the LBs belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter resulting lbs by name.
* `vpc_id` - (Optional) The ID of the VPC linked to the Load Balancers.
* `subnet_id` - (Optional) The ID of subnet that intrant load balancer belongs to.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lbs` - It is a nested type which documented below.
* `total_count` - Total number of Load Balancers that satisfy the condition.

The attribute (`lbs`) support the following:

* `id` - The ID of Load Balancer.
* `name` - The name of Load Balancer.
* `vpc_id` - The ID of the VPC linked to the Load Balancers.
* `subnet_id` - (Optional) The ID of subnet that intrant load balancer belongs to. 
* `private_ip` - The IP address of intranet IP.

