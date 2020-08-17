# Specify the provider and access details
provider "ksyun" {
}

resource "ksyun_lb_register_backend_server" "default" {
  backend_server_group_id="xxxx"
  backend_server_ip="192.168.5.xxx"
  backend_server_port="8081"
  weight=10
}
