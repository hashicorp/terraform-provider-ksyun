---
layout: "ksyun"
page_title: "Ksyun: ksyun_healthcheck"
sidebar_current: "docs-ksyun-resource-healthcheck"
description: |-
  Provides an Health Check resource.
---

# ksyun_healthcheck

Provides an Health Check resource.

## Example Usage

```hcl
resource "ksyun_healthcheck" "default" {
  listener_id = "537e2e7b-0007-4a75-9749-882167dbc93d"
  health_check_state = "stop"
  healthy_threshold = 2
  interval = 20
  timeout = 200
  unhealthy_threshold = 2
  url_path = "/monitor"
  is_default_host_name = true
  host_name = "www.ksyun.com"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required) The id of the listener.
* `health_check_state` - (Required) Status maintained by health examination.Valid Values:'start', 'stop'.
* `healthy_threshold` - (Required) Health threshold.Valid Values:1-10.
* `interval` - (Required) Interval of health examination.Valid Values:1-3600.
* `timeout` - (Required) Health check timeout.Valid Values:1-3600.
* `unhealthy_threshold` - (Required) Unhealthy threshold.Valid Values:1-10.
* `url_path ` - (Optional) Link to HTTP type listener health check.
* `host_name` - (Optional) Domain name of HTTP type health check.
