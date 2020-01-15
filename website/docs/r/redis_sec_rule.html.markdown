---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_sec_rule"
sidebar_current: "docs-ksyun-resource-redis-sec-rule"
description: |-
  Provides an redis security rule resource.
---

# ksyun_redis_sec_rule

Provides an redis security rule resource.

## Example Usage

```hcl
resource "ksyun_redis_sec_rule" "default" {
  cache_id          = "${ksyun_redis_instance.default.id}"
  available_zone    = "${var.available_zone}"
  rules=[
    {
      cidr="10.0.3.4/32"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The id of instance, .
* `cidrs` - (Required) The cidr block of source for the instance, multiple cidr separated by comma.

