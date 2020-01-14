---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_listener_servers"
sidebar_current: "docs-ksyun-datasource-lb-listener-servers"
description: |-
  Provides a list of Load Balancer Listener Server resources in the current region.
---

# ksyun_lb_listener_servers

This data source provides a list of Load Balancer Listener  Server resources according to their Load Balancer Listener Server ID.

## Example Usage

```hcl
data "ksyun_listener_servers" "default" {
  output_file="output_result"
  ids=[]
  listener_id=[]
  real_server_ip=["10.72.20.126","172.31.16.20"]
}
```

## Argument Reference

The following arguments are supported:


* `ids` - (Optional) A list of LB Listener Server IDs, all the LB Listener Servers belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listener_id` - The ID of LB Listener.
* `real_server_ip` - The name of LB Listener Server.
