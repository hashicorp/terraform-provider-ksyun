provider "ksyun"{
  region = "cn-beijing-6"
}

data "ksyun_krds" "search-krds"{
  output_file = "output_file"
  db_instance_type = "HRDS,RR,TRDS"
}

