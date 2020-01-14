---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_acl"
sidebar_current: "docs-ksyun-resource-lb-acl"
description: |-
  Provides a Load Balancer acl resource to add content forwarding policies for Load Balancer backend resource.
---

# ksyun_lb_acl

Provides a Load Balancer acl resource to add content forwarding policies for Load Balancer backend resource.

## Example Usage

```hcl
# Create Load Balancer Listener Acl
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "tf-xun2"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_acl_name` - (Required) The name of a load balancer acl.