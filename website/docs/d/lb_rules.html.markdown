| layout | page_title            | sidebar_current     | description                                                  |
| ------ | --------------------- | ------------------- | ------------------------------------------------------------ |
| ksyun  | Ksyun：ksyun_lb_rules | docs-ksyun-lb-rules | Provides a list of ksyun lb rules resources in the current region. |

# ksyun_lb_rules

  Provides a list of ksyun lb rules resources in the current region.

## Example Usage

```
provider "ksyun" {
}
# Get slb rule
data "ksyun_lb_rules" "default" {
output_file="output_result"
ids=[]
host_header_id=[]
}
```

## Argument Reference

The following arguments are supported:

- `ids` - （Optional）A list of rule IDs.
- `output_file` - (Optional) File name where to save data source results (after running terraform plan).
- `host_header_id` - (Optional）The id of host header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the rule was created.
- `rule_id` - The ID of rule.
- `path` - The path of rule
- `backend_server_group_id` - The id of backend server group
- `listener_sync` - Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'.
- `method` - Forwarding mode of listener.Valid Values:'RoundRobin', 'LeastConnections'.
- `session_state` - The state of session.Valid Values:'start', 'stop'.
- `session_persistence_period` - Session hold timeout.Valid Values:1-86400
- `cookie_type` - The type of the cookie.Valid Values:'ImplantCookie', 'RewriteCookie'.
- `cookie_name` - The name of cookie.The CookieType is valid and required when it is 'RewriteCookie'; otherwise, this value is ignored.
- `timeout` - Health check timeout.Valid Values:1-3600.
- `interval` - Interval of health examination.Valid Values:1-3600.
- `health_check_state` - Status maintained by health examination.The health check state is valid and selected when the ListenerSync is 'off ',otherwise, this value is ignored.Valid Values:'start', 'stop'.
- `healthy_threshold` - Health threshold.Valid and required when HealthCheckState is 'start', this value is ignored in other cases.Valid Values:1-10.
- `unhealthy_threshold` - Unhealthy threshold.Valid Values:1-10.
- `url_path` - Link to HTTP type listener health check.
- `host_name` - Domain name of HTTP type health check.