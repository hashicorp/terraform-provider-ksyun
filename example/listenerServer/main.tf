# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}


# Attach instances to Load Balancer
resource "ksyun_lb_listener_server" "default" {
  listener_id = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
  real_server_ip = "10.0.77.20"
  real_server_port = 8000
  real_server_type = "host"
  instance_id = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
  weight = 10
}

