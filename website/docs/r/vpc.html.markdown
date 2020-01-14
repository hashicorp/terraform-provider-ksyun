---
layout: "ksyun"
page_title: "Ksyun: ksyun_vpc"
sidebar_current: "docs-ksyun-resource-vpc"
description: |-
  Provides a VPC resource.
---

# ksyun_vpc

Provides a VPC resource.

~> **Note**  The network segment can only be created or deleted, can not perform both of them at the same time.
## Example Usage

```hcl
resource "ksyun_vpc" "example" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.1.0.2/24"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required) The CIDR blocks of VPC.
* `vpc_name` - (Optional) The name of the vpc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for VPC, formatted in RFC3339 time string.
* `cidr_block` - The CIDR block of the VPC.

## Import

VPC can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpc.example uvnet-abc123456
```