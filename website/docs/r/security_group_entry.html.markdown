---
layout: "ksyun"
page_title: "Ksyun: ksyun_security_group_entry"
sidebar_current: "docs-ksyun-resource-security-group"
description: |-
  Provides a Security Group resource.
---

# ksyun_security_group_entry

Provides a Security Group resource.

## Example Usage

```hcl
resource "ksyun_security_group_entry" "default" {
  security_group_id="7385c8ea-79f7-4e9c-b99f-517fc3726256"
  cidr_block="10.0.0.1/32"
  direction="in"
  protocol="ip"
}

```

## Argument Reference

The following arguments are supported:


* `description` - (Optional) The description of the security group .
* `security_group_id` - (Required) The ID of the security group.
* `cidr_block` - (Required) The cidr block of security group rules.
* `direction` - (Required) .Valid Values:'in', 'out'.
* `protocol` - (Required) protocol.Valid Values:'ip', 'tcp', 'udp', 'icmp'.
* `icmp_type` - (Optional) ICMP protocol.The required if protocol type is 'icmp'.
* `icmp_code` - (Optional) ICMP protocol.The required if protocol type is 'icmp'.
* `port_range_from` - (Optional) Port rule start port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.
* `port_range_to` - (Optional) Port rule start port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.
 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of security group, formatted in RFC3339 time string.

## Import

Security Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group_entry.example firewall-abc123456
```