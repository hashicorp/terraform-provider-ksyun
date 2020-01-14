# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  eips
data "ksyun_eips" "default" {
  output_file="output_result"

  ids=[]
  project_id=["0"]
  instance_type=["Ipfwd"]
  network_interface_id=[]
  internet_gateway_id=[]
  band_width_share_id=[]
  line_id=[]
  public_ip=[]
}

