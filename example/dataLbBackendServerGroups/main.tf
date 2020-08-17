# Specify the provider and access details
provider "ksyun" {
}

# Get  availability zones
data "ksyun_lb_backend_server_groups" "default" {
  output_file="out_file"
  ids=[]
  host_header_id=[]
}

