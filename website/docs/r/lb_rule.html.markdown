---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_rule"
sidebar_current: "docs-ksyun-lb-rule"
description: |-
  Provides a lb rule resource.
---


# ksyun_lb_rule

  Provides a lb rule resource.

## Example Usage

```hcl
provider "ksyun" {
}
resource "ksyun_lb_rule" "default" {
path = "/tfxun/update",
host_header_id = "",
backend_server_group_id=""
listener_sync="on"
method="RoundRobin"
session {
session_state = "start"
session_persistence_period = 1000
cookie_type = "ImplantCookie"
cookie_name = "cookiexunqq"
}
health_check{
health_check_state = "start"
healthy_threshold = 2
interval = 200
timeout = 2000
unhealthy_threshold = 2
url_path = "/monitor"
host_name = "www.ksyun.com"
}
}
```

## Argument Reference

The following arguments are supported:

- `path` - (Required) The path of rule
- `host_header_id` - (Required）The id of host header id
- `backend_server_group_id` - （Required）The id of backend server group
- `listener_sync` - (Required）Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'.
- `method` - (Optional) Forwarding mode of listener.Valid Values:'RoundRobin', 'LeastConnections'.
- `session_state` - (Optional) The state of session.Valid Values:'start', 'stop'.
- `session_persistence_period` - (Optional) Session hold timeout.Valid Values:1-86400
- `cookie_type` - (Optional) The type of the cookie.Valid Values:'ImplantCookie', 'RewriteCookie'.
- `cookie_name` - (Optional) The name of cookie.The CookieType is valid and required when it is 'RewriteCookie'; otherwise, this value is ignored.
- `timeout` - (Optional) Health check timeout.Valid Values:1-3600.
- `interval` - (Optional) Interval of health examination.Valid Values:1-3600.
- `health_check_state` - (Optional) Status maintained by health examination.The health check state is valid and selected when the ListenerSync is 'off ',otherwise, this value is ignored.Valid Values:'start', 'stop'.
- `healthy_threshold` - (Optional) Health threshold.Valid and required when HealthCheckState is 'start', this value is ignored in other cases.Valid Values:1-10.
- `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10.
- `url_path` - (Optional) Link to HTTP type listener health check.
- `host_name` - (Optional) Domain name of HTTP type health check.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the rule was created.
- `rule_id` - The ID of rule.
- `cookie_expiration_period` - Session holds timeout time.Valid Values: 0-86400.
- `backend_server_group_id` - The id of backend server group
- `listener_sync` - Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'.
- `method` - Forwarding mode of listener.Valid Values:'RoundRobin', 'LeastConnections'.
- `session_state` - The state of session.Valid Values:'start', 'stop'.
- `session_persistence_period` - Session hold timeout.Valid Values:1-86400
- `cookie_type` - The type of the cookie.Valid Values:'ImplantCookie', 'RewriteCookie'.
- `cookie_name` - The name of cookie.The CookieType is valid and required when it is 'RewriteCookie'; otherwise, this value is ignored.
- `timeout` - Health check timeout.Valid Values:1-3600.
- `interval` -Interval of health examination.Valid Values:1-3600.
- `health_check_state` -Status maintained by health examination.The health check state is valid and selected when the ListenerSync is 'off ',otherwise, this value is ignored.Valid Values:'start', 'stop'.
- `healthy_threshold` - Health threshold.Valid and required when HealthCheckState is 'start', this value is ignored in other cases.Valid Values:1-10.
- `unhealthy_threshold` - Unhealthy threshold.Valid Values:1-10.
- `url_path` - Link to HTTP type listener health check.
- `host_name` - Domain name of HTTP type health check.