# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region = "region"
}

data "ksyun_redis_instances" "default" {
  output_file       = "output_result1"
  fuzzy_search      = ""
  iam_project_id    = ""
  cache_id          = ""
  vnet_id           = ""
  vpc_id            = ""
  name              = ""
  vip               = ""
}