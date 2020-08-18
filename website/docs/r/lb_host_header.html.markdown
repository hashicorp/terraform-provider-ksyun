| layout | page_title                  | sidebar_current           | description                         |
| ------ | --------------------------- | ------------------------- | ----------------------------------- |
| ksyun  | Ksyun: ksyun_lb_host_header | docs-ksyun-lb-host-header | Provides a lb host header resource. |

# ksyun_lb_host_header

Provides a lb host header resource.

## Example Usage

```
resource "ksyun_lb_host_header" "foo" {
listener_id = "xxxx"
host_header = "tf-xuan"
certificate_id = ""
}
```

## Argument Reference

The following arguments are supported:

- `listener_id` - (Required) The ID of the listener.
- `host_header` - (Required) The hostheader.
- `certificate_id` - (Optional) The ID of the certificate, HTTPS type listener creates this parameter which is not default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `create_time` - The time when the hostheader was created.
- `host_header_id` - The ID of hostheader.