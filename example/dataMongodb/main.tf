# Specify the provider and access details
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region = "cn-beijing-6"
}

data "ksyun_mongodbs" "default" {
  output_file = "output_result"
  iam_project_id = ""
  instance_id = ""
  vnet_id = ""
  vpc_id = ""
  name = ""
  vip = ""
}
