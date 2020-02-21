---
layout: "ksyun"
page_title: "Ksyun: ksyun_health_checks"
sidebar_current: "docs-ksyun-datasource-health_checks"
description: |-
  Provides a list of healthcheck resources in the current region.
---

# ksyun_health_checks

This data source provides a list of healthcheck resources  according to their healthcheck ID or listener ID.

## Example Usage

```hcl
data "ksyun_health_checks" "default" {
  output_file="output_result"
  ids=[]
  listener_id=["8d1dac22-6c6c-42ea-93e2-2702d44ddb93","70467f7e-23dc-465a-a609-fb1525fc6b16"]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of health check IDs, all the healthcheck belong to this region will be retrieved if the ID is `""`.
* `listener_id` - (Optional) A list of listener IDs, all the healthcheck belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `health_checks` - It is a nested type which documented below.
* `total_count` - Total number of Elastic IPs that satisfy the condition.
