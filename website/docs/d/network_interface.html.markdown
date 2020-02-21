---
layout: "ksyun"
page_title: "Ksyun: ksyun_network_interface"
sidebar_current: "docs-ksyun-datasource-network_interface"
description: |-
  Provides a list of Network Interface in the current region.
---

# ksyun_network_interface

This data source provides a list of Network Interface resources according to their Network Interface ID.

## Example Usage

```hcl
data "ksyun_network_interfaces" "default" {
  output_file="output_result"
  ids=[]
  securitygroup_id=[]
  instance_type=[]
  instance_id=[]
  private_ip_address=[]
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Network Interface IDs, all the Network Interfaces belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_interfaces` - It is a nested type which documented below.

The attribute (`network_interfaces`) support the following:

* `id` - The id of the network interface.
* `securitygroup_id` - The ID of the security group.
* `instance_type` - The type of network interface.
* `instance_id` - The ID of instance.
* `private_ip_address` - The private IP address assigned to the instance.