# Specify the provider and access details
provider "ksyun" {
}

# Get  eips
data "ksyun_subnets" "default" {
  output_file="output_result"

  ids=[]
  vpc_id=[]
  nat_id=[]
  network_acl_id=[]
  subnet_type=[]
  availability_zone_name=[]

}

