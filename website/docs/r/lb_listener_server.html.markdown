---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_listener_server"
sidebar_current: "docs-ksyun-resource-lb-listener-server"
description: |-
  Provides a Load Balancer Listener server resource.
---

# ksyun_lb_listener_server

Provides a Load Balancer Listener server resource.

## Example Usage

```hcl
resource "ksyun_lb_listener_server" "default" {
  listener_id = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
  real_server_ip = "10.0.77.20"
  real_server_port = 8000
  real_server_type = "host"
  instance_id = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
  weight = 10
}

```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required) The id of the listener.
* `real_server_type` - (Required) The type of real server.Valid Values:'Host', 'DirectConnectGateway', 'VpnTunnel'.
* `instance_id` - (Optional) The ID of instance.
* `real_server_ip` - (Required) The IP of real server.
* `real_server_port` - (Required) The port of real server.Valid Values:1-65535
* `weight` - (Optional) The weight of backend service.Valid Values:1-255

## Import

LB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_listener.example vserver-abcdefg
```