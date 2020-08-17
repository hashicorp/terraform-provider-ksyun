# Specify the provider and access details
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