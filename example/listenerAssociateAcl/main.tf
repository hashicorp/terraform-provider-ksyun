# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

resource "ksyun_lb_listener_associate_acl" "default" {
  listener_id = "b330eae5-11a3-4e9e-bf7d-a7a1117a5878"
  load_balancer_acl_id = "7e94fa82-05c7-496c-ae5e-35fd32ff3cf2"
}