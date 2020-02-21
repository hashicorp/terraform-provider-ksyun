---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb"
sidebar_current: "docs-ksyun-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# ksyun_lb

Provides a Load Balancer resource.

## Example Usage

```hcl
resource "ksyun_lb" "default" {
  vpc_id = "74d0a45b-472d-49fc-84ad-221e21ee23aa"
  load_balancer_name = "tf-xun1"
  type = "public"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` - (Optional) The name of the load balancer. 
* `vpc_id` - (Required) The ID of the VPC linked to the Load Balancers.
* `type` - (Optional) The type of load balancer.Valid Values:'public', 'internal'.
* `subnet_id` - (Optional) The id of the subnet.only Internal type is Required.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for load balancer, formatted in RFC3339 time string.
* `private_ip` - The IP address of intranet IP. It is `""` if `internal` is `false`.

## Import

LB can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```