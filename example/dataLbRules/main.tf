# Specify the provider and access details
provider "ksyun" {
}

# Get  slb rule
data "ksyun_lb_rules" "default" {
  output_file="output_result"
  ids=[]
  host_header_id=[]
}

