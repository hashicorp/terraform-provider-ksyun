| layout | page_title                           | sidebar_current                    | description                                  |
| ------ | ------------------------------------ | ---------------------------------- | -------------------------------------------- |
| ksyun  | Ksyun: ksyun_lb_backend_server_group | docs-ksyun-lb-backend-server-group | Provides a lb backend server group resource. |

# ksyun_lb_backend_server_group

Provides a lb backend server group resource.

## Example Usage

```
provider "ksyun" {
}
resource "ksyun_lb_backend_server_group" "default" {
backend_server_group_name="xuan-tf"
vpc_id=""
backend_server_group_type=""
}
```

## Argument Reference

The following arguments are supported:

- `backend_server_group_name` - (Required) The name of backend server group.
- `vpc_id` - (Required) Virtual private network ID.
- `backend_server_group_type` - (Optional) The type of backend server group.Valid values are Server and Mirror.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the backend server group was created.
- `backend_server_group_id` - The ID of backend server group.
- `backend_server_number` - The number of backend server group.
- `health_check` - Health check information, only the mirror server has this parameter.