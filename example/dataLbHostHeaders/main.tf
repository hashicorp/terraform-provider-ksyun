# Specify the provider and access details
provider "ksyun" {
}

# Get  slbs
data "ksyun_lb_host_headers" "default" {
  output_file="output_result"
  ids=[]
  listener_id=[]
}