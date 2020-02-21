---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_acls"
sidebar_current: "docs-ksyun-datasource-lb-acls"
description: |-
  Provides a list of Load Balancer Rule resources belong to the Load Balancer listener.
---

# ksyun_lb_acls

This data source provides a list of Load Balancer Rule resources according to their Load Balancer Rule ID.

## Example Usage

```hcl
data "ksyun_lb_acls" "default" {
  output_file="output_result"
  ids=[]

}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of LB Rule IDs, all the LB Rules belong to the Load Balancer listener will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lb_acls` - It is a nested type which documented below.
* `total_count` - Total number of LB Rules that satisfy the condition.

