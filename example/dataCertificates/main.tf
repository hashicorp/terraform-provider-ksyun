# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_certificates" "default" {
  output_file="output_result"
  ids=[]
}

