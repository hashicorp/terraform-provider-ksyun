---
layout: "ksyun"
page_title: "Ksyun: ksyun_zones"
sidebar_current: "docs-ksyun-datasource-zones"
description: |-
  Provides a list of available zones in the current region.
---

# ksyun_zones

This data source provides a list of available zones in the current region.

## Example Usage

```h
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone_name` - Name of the zone.
