provider "ksyun"{
  region = "cn-shanghai-2"
}

data "ksyun_krds_security_groups" "hou_desc" {
  output_file = "output_file"
  security_group_id = 2813
}