---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_sec_group"
sidebar_current: "docs-ksyun-resource-redis-sec-group"
description: |-
  Provides an redis security rule resource.
---

# ksyun_redis_sec_group

Provides an redis security group function.

## Example Usage

```hcl
variable "available_zone" {
  default = "cn-beijing-6a"
}

resource "ksyun_redis_sec_group" "add" {
  available_zone = "${var.available_zone}"
  name = "testAddTerraform"
  description = "testAddTerraform"
}

resource "ksyun_redis_sec_group_rule" "default" {
  available_zone = "${var.available_zone}"
  security_group_id = "${ksyun_redis_sec_group.add.id}"
  rules = ["172.16.0.0/32","192.168.0.0/32"]
}

resource "ksyun_redis_sec_group_allocate" "default" {
  available_zone = "${var.available_zone}"
  security_group_id = "${ksyun_redis_sec_group.add.id}"
  cache_ids = ["122334234"]
}
```

## Argument Reference

The following arguments are supported:

**ksyun_redis_sec_group**

* `available_zone`- (Required) The Zone to launch the security group .
* `name` - (Required) The name of   the security group.
*  `description ` - (Required) The description of   the security group.

**ksyun_redis_sec_group_rule**

- `security_group_id`- (Required) The ID of  the security group .
- `available_zone`- (Required) The Zone to launch the security group .
- `rules` - (Required) The cidr block of source for the instance, multiple cidr separated by comma.

**ksyun_redis_sec_group_allocate**

- `available_zone`- (Required) The Zone to launch the security group .
- `security_group_id`- (Required) The ID of  the security group .
- `cache_ids` - (Required) The ids of   the redis instance .