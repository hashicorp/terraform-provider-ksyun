# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region = "cn-shanghai-3"
}

data "ksyun_volumes" "default" {
  output_file="output_result"
  ids=[]
  volume_category=""
  volume_status=""
  volume_type=""
  availability_zone=""
}
