---
layout: "ksyun"
page_title: "Ksyun: ksyun_lines"
sidebar_current: "docs-ksyun-datasource-lines"
description: |-
  Provides a list of line resources in the current region.
---

# ksyun_lines

This data source provides a list of line resources supported.

## Example Usage

```hcl
# Get  lines
data "ksyun_lines" "default" {
  output_file="output_result"
  line_name="BGP"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lines` - All the lines accourding the argument.
* `total_count` - Total number of lines  that satisfy the condition.

