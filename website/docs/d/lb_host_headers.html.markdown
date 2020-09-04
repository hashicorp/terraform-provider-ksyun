---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_host_headers"
sidebar_current: "docs-ksyun-datasource-lb-host-headers"
description: |-
  Provides a list of lb host headers in the current region.
---

# ksyun_lb_host_headers

  Provides a list of lb host headers in the current region.

## Example Usage

```hcl
provider "ksyun" {
}
# Get slbs
data "ksyun_lb_host_headers" "default" {
output_file="output_result"
ids=[]
listener_id=[]
}
```

### Argument Reference

- `ids` - (Optional) A list of hostheader IDs.
- `listener_id` - (Optional) The ID of listener.
- `output_file` - (Optional) File name where to save data source results (after running terraform plan).

### Attributes Reference

- `create_time` - The time when the hostheader was created.
- `host_header_id` - The ID of hostheader.
- `host_header` - The hostheader.
- `certificate_id` - The ID of certificate, HTTPS type listener creates this parameter which is not default.