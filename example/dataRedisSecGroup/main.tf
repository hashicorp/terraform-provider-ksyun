# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region = "region"
}

data "ksyun_redis_security_groups" "default" {
  output_file       = "output_result1"
}