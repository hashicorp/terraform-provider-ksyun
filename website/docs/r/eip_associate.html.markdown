---
layout: "ksyun"
page_title: "Ksyun: ksyun_eip_association"
sidebar_current: "docs-ksyun-resource-eip-association"
description: |-
  Provides an EIP Association resource for associating Elastic IP to UHost Instance, Load Balancer, etc..
---

# ksyun_eip_association

Provides an EIP Association resource for associating Elastic IP to UHost Instance, Load Balancer, etc.

## Example Usage

```hcl
resource "ksyun_eip_associate" "slb" {
  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
  instance_type="Slb"
  instance_id="7fae85e4-ab1a-415c-aef9-03a402c79d97"
}
resource "ksyun_eip_associate" "server" {
  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
  instance_type="Ipfwd"
  instance_id="566567677-6766-4743-afb7-7c7081214092"
  network_interface_id="87945980-59659-04548-759045803"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required) The ID of EIP.
* `instance_type` - (Required) The type of the instance.Valid Values:'Ipfwd', 'Slb'.
* `instance_id` - (Required) The id of the instance.
* `network_interface_id` - (Optional) The id of the network interface.

