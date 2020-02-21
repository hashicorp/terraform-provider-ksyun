---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_listener_associate_acl"
sidebar_current: "docs-ksyun-resource-lb-listener"
description: |-
  Associate a Load Balancer Listener resource with acl.
---

# ksyun_lb_listener_associate_acl

Associate a Load Balancer Listener resource with acl.

## Example Usage

```hcl
resource "ksyun_lb_listener_associate_acl_associate_acl" "default" {
  listener_id = "b330eae5-11a3-4e9e-bf7d-a7a1117a5878"
  load_balancer_acl_id = "7e94fa82-05c7-496c-ae5e-35fd32ff3cf2"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of load balancer instance.
* `listener_name` - (Optional) The id of the listener. 

