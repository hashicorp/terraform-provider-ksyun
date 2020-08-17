# Specify the provider and access details
provider "ksyun" {
}

resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name="xuan-tf"
  vpc_id=""
  backend_server_group_type=""
}
