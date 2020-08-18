---
layout: "ksyun"
page_title: "Ksyun: ksyun_instances"
sidebar_current: "docs-ksyun-datasource-instances"
description: |-
  Provides a list of instance resources in the current region.
---

# ksyun_instances

This data source providers a list of instance resources according to their availability zone, instance ID.

## Example Usage

```h
# Get  instances
data "ksyun_instances" "default" {
  output_file = "output_result"
  ids = []
  search = ""
  project_id = []
  network_interface {
  network_interface_id = []
  subnet_id = []
  group_id = []
  }
  instance_state {
  name =  []
  }
  availability_zone {
  name =  []
  }
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `research` - (Optional) A regex string to filter results by instance name or privateIpAddress.
* `image_id` - (Optional) The image ID of some instance used.
* `subnet_id` - (Optional) The ID of subnet. the instance will use the subnet in the current region.
* `security_group_id` - (Optional) Security Group to associate with.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - (Optional) A list of instance IDs.
* `instances` - instances documented below.

The attribute (`instances`) support the following:

* `instance_id` - The ID of instance.
* `instance_state` - The state of instance.
* `vpc_id` - The ID of VPC linked to the instance.
* `subnet_id` - The ID of subnet linked to the instance.
* `image_id` - The ID for the image to use for the instance.
* `instance_type` -  The type of instance.
* `security_group_id` - Security Group to associate with.
* `instance_name` - The name of instance.
* `project_id` -  The project instance belongs to.
* `user_data` - The user data to be specified into this instance. 
* `creation_date` - Time of creation.
* `charge_type` - Instance charge type.
* `availability_zone_name` - Name of the instance belongs to.
* `private_ip_address` - Instance private IP address.
* `project_id` - The ID of project_id the instance belongs to.
* `disk_size` - The size of systemdisk.
* `disk_type` - The type of systemdisk.