---
layout: "ksyun"
page_title: "Ksyun: ksyun_security_group"
sidebar_current: "docs-ksyun-resource-security-group"
description: |-
  Provides a Security Group resource.
---

# ksyun_security_group

Provides a Security Group resource.

## Example Usage

```hcl
resource "ksyun_security_group" "default" {
  vpc_id = "26231a41-4c6b-4a10-94ed-27088d5679df"
  security_group_name="xuan-tf--s"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_name` - (Optional) The name of the security group which contains 1-63 characters and only support Chinese, English, numbers, '-', '_' and '.'. 
* `vpc_id` - (Optional) The Id of the vpc.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of security group, formatted in RFC3339 time string.

## Import

Security Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group.example firewall-abc123456
```