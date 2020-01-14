# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

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
