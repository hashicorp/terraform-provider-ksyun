# Specify the provider and access details
provider "ksyun" {
  region="cn-beijing-6"
}

# Get  network_interfaces
data "ksyun_network_interfaces" "default" {
  output_file="output_result"
  ids=[]
  vpc_id=[]
  subnet_id=[]
  securitygroup_id=[]
  instance_type=[]
  instance_id=[]
  private_ip_address=[]

}

