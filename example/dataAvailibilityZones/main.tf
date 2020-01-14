# Specify the provider and access details
provider "ksyun" {
  region="cn-beijing-6"
}

# Get  availability zones
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}

