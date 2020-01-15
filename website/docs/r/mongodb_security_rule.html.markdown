---
layout: "ksyun"
page_title: "Ksyun: ksyun_mongodb_security_rule"
sidebar_current: "docs-ksyun-resource-mongodb-security-rule"
description: |-
  Provides an MongoDB Security Rule resource.
---

# ksyun_mongodb_security_rule

Provides an MongoDB Security Rule resource.

## Example Usage

```hcl
resource "ksyun_mongodb_security_rule" "default" {
  instance_id = "InstanceId"
  cidrs = "192.168.10.1/32"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The id of instance, .
* `cidrs` - (Required) The cidr block of source for the instance, multiple cidr separated by comma.

