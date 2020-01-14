provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}

data "ksyun_krds_security_groups" "hou_desc" {
  output_file = "output_file"
  security_group_id = 27937
}