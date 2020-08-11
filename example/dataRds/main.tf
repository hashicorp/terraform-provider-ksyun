provider "ksyun"{
  region = "cn-shanghai-2"
}

data "ksyun_krds" "search-krds"{
  output_file = "output_file"
  db_instance_type = "HRDS,RR,TRDS"
}

output "instance_id" {
  value = "${data.ksyun_krds.search-krds.krds.1.vip}"
}
output "krds" {
  value = [for instance in data.ksyun_krds.search-krds.krds:
  instance.vip
  ]
  #value = "${data.ksyun_krds.krds.krds.1.vip}"
}