---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_backend_server_groups"
sidebar_current: "docs-ksyun-datasource-ksyun-lb-backend-server-groups"
description: |-
  Provides a list of lb backend server groups in the current region.
---

# ksyun_lb_backend_server_groups

  Provides a list of lb backend server groups in the current region.

## Example Usage

```hcl
provider "ksyun" {
}
# Get availability zones
data "ksyun_lb_backend_server_groups" "default" {
output_file="out_file"
ids=[]
}
```

## Argument Reference

The following arguments are supported:

- `ids` - (Optional) A list of backend server group IDs.
- `output_file` - (Optional) File name where to save data source results (after running terraform plan).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the backend server group was created.
- `vpc_id` - Virtual private network ID.

- `backend_server_group_id` - The id of backend server group.
- `backend_server_group_name` - The name of backend server group.
- `backend_server_number` - The number of backend server number.
- `backend_server_group_type` - The type of backend server group.Valid values are Server and Mirror.
- `health_check` - Health check information, only the mirror server has this parameter.