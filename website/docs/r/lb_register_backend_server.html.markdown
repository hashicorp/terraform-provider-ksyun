---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_register_backend_server"
sidebar_current: "docs-ksyun-lb-register-backend-server"
description: |-
  Provides a lb register backend server resource.
---


# ksyun_lb_register_backend_server

Provides a lb register backend server resource.

## Example Usage

```hcl
provider "ksyun" {
}
resource "ksyun_lb_register_backend_server" "default" {
backend_server_group_id="xxxx"
backend_server_ip="192.168.5.xxx"
backend_server_port="8081"
weight=10
}
```

## Argument Reference

The following arguments are supported:

- `backend_server_group_id` - (Required) The ID of backend server group.
- `backend_server_ip` - (Required) The IP of backend server.
- `backend_server_port` - (Required) The port of backend server.Valid Values:1-65535
- `weight` - (Optional) The weight of backend service.Valid Values:0-255Mirror.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the backend service was created.
- `register_id` - The registration ID of binding server group.
- `real_server_ip` - The IP of real server.
- `real_server_port` - The port of real server.Valid Values:1-65535.
- `real_server_type` - The type of real server.Valid Values:'Host'.
- `master_slave_type` - The type of real server.Only MasterSlave listener has this parameter.The Valid Values:'Master','Slave'.
- `instance_id` - The ID of instance.
- `network_interface_id` - The ID of network interface.
- `real_server_state` - The state of real server.Values:'healthy','unhealthy'