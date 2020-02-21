---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_instance_node"
sidebar_current: "docs-ksyun-resource-redis-instance-node"
description: |-
  Provides an redis instance node resource.
---

# ksyun_redis_instance_node

Provides an redis instance node resource.

## Example Usage

```hcl
resource "ksyun_redis_instance_node" "default" {
  cache_id          = "${ksyun_redis_instance.default.id}"
  available_zone    = "${var.available_zone}"
}

resource "ksyun_redis_instance_node" "node" {
  // creating multiple read-only nodes,
  // not concurrently, requires dependencies to synchronize the execution of creating multiple read-only nodes.
  // if only one read-only node is created, it is not required to fill in.
  pre_node_id       = "${ksyun_redis_instance_node.default.id}"
  cache_id          = "${ksyun_redis_instance.default.id}"
  available_zone    = "${var.available_zone}"
}
```

## Argument Reference

The following arguments are supported:

* `cache_id` - (Optional)  The ID of  the intance .
 * `available_zone` - (Optional) The Zone to launch the DB instance.

