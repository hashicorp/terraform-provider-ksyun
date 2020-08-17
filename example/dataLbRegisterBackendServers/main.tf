# Specify the provider and access details
provider "ksyun" {
  region="cn-beijing-6"
}

data "ksyun_lb_register_backend_servers" "foo" {
  output_file="output_result"
  ids=[]
  backend_server_group_id=[]
}

