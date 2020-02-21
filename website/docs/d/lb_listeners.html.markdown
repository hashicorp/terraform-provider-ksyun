---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_listeners"
sidebar_current: "docs-ksyun-datasource-lb-listeners"
description: |-
  Provides a list of Load Balancer Listener resources in the current region.
---

# ksyun_lb_listeners

This data source provides a list of Load Balancer Listener resources according to their Load Balancer Listener ID.

## Example Usage

```hcl
data "ksyun_listeners" "default" {
  output_file="output_result"
  ids=[""]
  load_balancer_id=["d3fd0421-a35a-4ddb-a939-5c51e8af8e8c","4534d617-9de0-4a4a-9ed5-3561196cacb6"]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of a load balancer.
* `ids` - (Optional) A list of LB Listener IDs, all the LB Listeners belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listeners` - It is a nested type which documented below.
* `total_count` - Total number of LB listeners that satisfy the condition.

The attribute (`lb_listeners`) support the following:

* `id` - The ID of LB Listener.
* `listener_name` - The name of LB Listener.
* `listener_protocol` - LB Listener protocol. Possible values: `http`, `https` , `tcp` and `udp` .
* `listener_port` - Port opened on the LB Listener to receive requests, range: 1-65535.
* `method` - The load balancer method in which the listener is.