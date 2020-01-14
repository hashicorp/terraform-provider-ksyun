# Specify the provider and access details
provider "ksyun" {
}

data "ksyun_security_groups" "default" {
  output_file="output_result"
  ids=[]
  vpc_id=[]
}

