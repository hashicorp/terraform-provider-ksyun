| layout | page_title                            | sidebar_current                                | description                                                  |
| ------ | ------------------------------------- | ---------------------------------------------- | ------------------------------------------------------------ |
| ksyun  | Ksyun: ksyun_lb_register_backend_servers | docs-ksyun-datasource-lb-register-backend-servers | Provides a list of register backend servers in the current region. |

# ksyun_register_backend_servers

Provides a list of register backend servers in the current region.

## Example Usage

```
provider "ksyun" {
region="cn-beijing-6"
}
data "ksyun_lb_register_backend_servers" "foo" {
output_file="output_result"
ids=[]
backend_server_group_id=[]
}
```

## Argument Reference

The following arguments are supported:

-  `ids` - (Optional) A list of backend service IDs.
- `backend_server_group_id` - (Optional) The ID of backend server group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the backend server was created.
- `backend_server_group_id` - The id of backend server group.
- `register_id` - The registration ID of the binding server group. 
- `real_server_ip` - The IP of real server.
- `real_server_port` - The port of real server.Valid Values:1-65535.
- `real_server_type` - The type of real server.Valid Values:'Host'.
- `master_slave_type` - The type of real server.Only MasterSlave listener has this parameter.The Valid Values:'Master','Slave'.
- `instance_id` - The ID of instance.
- `network_interface_id` - The ID of network interface.
- `real_server_state` - The state of real server.Values:'healthy','unhealthy'
- `weight` - The weight of backend service.Valid Values:1-255
