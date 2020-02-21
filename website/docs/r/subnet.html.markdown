---
layout: "ksyun"
page_title: "Ksyun: ksyun_subnet"
sidebar_current: "docs-ksyun-resource-subnet"
description: |-
  Provides a Subnet resource under VPC resource.
---

# ksyun_subnet

Provides a Subnet resource under VPC resource.

## Example Usage

```hcl
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "example" {
  subnet_name      = "tf-acc-subnet1"
  	cidr_block = "10.0.5.0/24"
      subnet_type = "Normal"
      dhcp_ip_from = "10.0.5.2"
      dhcp_ip_to = "10.0.5.253"
      vpc_id  = "${ksyun_vpc.test.id}"
      gateway_ip = "10.0.5.1"
      dns1 = "198.18.254.41"
      dns2 = "198.18.254.40"
      availability_zone = "cn-shanghai-2a"
}
```

## Argument Reference

The following arguments are supported:

* `subnet_name` - (Required) The name of the subnet.
* `cidr_block` - (Required) The CIDR block assigned to the subnet.
* `subnet_type ` - (Required) The type of subnet.Valid Values:'Reserve', 'Normal', 'Physical'.
* `dhcp_ip_from` - (Required) DHCP start IP.
* `dhcp_ip_to` - (Required) DHCP end IP.
* `gateway_ip` - (Required) The IP of gateway.
* `vpc_id` - (Required) The id of the vpc.
* `dns1` - (Optional) The dns of the subnet.
* `dns2` - (Optional) The dns of the subnet.
* `availability_zone` - (Optional) The name of the availability zone. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of subnet, formatted in RFC3339 time string.

## Import

Subnet can be imported using the `id`, e.g.

```
$ terraform import ksyun_subnet.example subnet-abc123456
```