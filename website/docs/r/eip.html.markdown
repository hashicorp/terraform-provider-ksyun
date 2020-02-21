---
layout: "ksyun"
page_title: "Ksyun: ksyun_eip"
sidebar_current: "docs-ksyun-resource-eip"
description: |-
  Provides an Elastic IP resource.
---

# ksyun_eip

Provides an Elastic IP resource.

## Example Usage

```hcl
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}
resource "ksyun_eip" "default1" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =1
  charge_type = "PrePaidByMonth"
  purchase_time =1
  project_id=0
}
```

## Argument Reference

The following arguments are supported:

* `line_id` - (Required) The id of the line.
* `band_width` - (Required) The band width of the public address.
* `charge_type` - (Required) The charge type of the Elastic IP address.Valid Values:'PrePaidByMonth', 'PostPaidByPeak', 'PostPaidByDay', 'PostPaidByTransfer', 'PostPaidByHour', 'HourlyInstantSettlement'.
* `purchase_time` - (Required) Purchase time.
* `project_id` - (Optional) The id of the project.

 
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_ip` -  The Elastic IP address.
* `id` -  The ID of the Elastic IP .
