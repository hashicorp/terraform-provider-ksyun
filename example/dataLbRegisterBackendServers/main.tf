# Specify the provider and access details
provider "ksyun" {
  region="cn-beijing-6"
}

data "ksyun_lb_register_backend_servers" "default" {
  output_file="out_file"
  ids=[]
  backend_server_group_id=[]
}

