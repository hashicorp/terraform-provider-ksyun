---
layout: "ksyun"
page_title: "Ksyun: ksyun_security_groups"
sidebar_current: "docs-ksyun-datasource-security-groups"
description: |-
  Provides a list of Security Group resources in the current region.
---

# ksyun_security_groups

This data source provides a list of Security Group resources according to their Security Group ID, name and resource id.

## Example Usage

```hcl
data "ksyun_security_groups" "default" {
  output_file="output_result"
  ids=[]
  vpc_id=[]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Security Group IDs, all the Security Group resources belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - It is a nested type which documented below.
* `total_count` - Total number of Security Group resources that satisfy the condition.

The attribute (`security_groups`) support the following:

* `id` - The ID of Security Group.
* `security_group_name` - The name of Security Group.
* `security_group_entry` - It is a nested type which documented below.
* `security_group_type` - The type of Security Group.
* `vpc_id` - The ID of vpc .
* `create_time` - The time of creation for the security group, formatted in RFC3339 time string.

The attribute (`security_group_entry`) support the following:

* `cidr_block` - The cidr block of source.
* `port_range_from` - The start of port numbers .
* `port_range_to` - The end of port numbers.
* `protocol` - The protocol. Can be `tcp`, `udp`, `icmp`, `ip`.
