# Specify the provider and access details
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region = "cn-beijing-6"
}

data "ksyun_availability_zones" "default" {
  output_file = ""
  ids = []
}
resource "ksyun_vpc" "default" {
  vpc_name = "ksyun-vpc-tf"
  cidr_block = "10.1.0.0/23"
}
resource "ksyun_subnet" "default" {
  subnet_name = "ksyun-subnet-tf"
  cidr_block = "10.1.0.0/23"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_mongodb_instance" "repset" {
  name = "mongodb_repset"
  instance_account = "root"
  instance_password = "admin"
  instance_class = "1C2G"
  storage = 5
  node_num = 3
  vpc_id = "${ksyun_vpc.default.id}"
  vnet_id = "${ksyun_subnet.default.id}"
  db_version = "3.6"
  pay_type = "byDay"
  iam_project_id = "0"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_mongodb_shard_instance" "shard" {
  name = "mongodb_shard"
  instance_account = "root"
  instance_password = "admin"
  mongos_class = "1C2G"
  mongos_num = 2
  shard_class = "1C2G"
  shard_num = 2
  storage = 5
  vpc_id = "${ksyun_vpc.default.id}"
  vnet_id = "${ksyun_subnet.default.id}"
  db_version = "3.6"
  pay_type = "hourlyInstantSettlement"
  iam_project_id = "0"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_mongodb_security_rule" "repset" {
  instance_id = "${ksyun_mongodb_instance.repset.id}"
  cidrs = "192.168.10.1/32"
}

resource "ksyun_mongodb_security_rule" "shard" {
  instance_id = "${ksyun_mongodb_shard_instance.shard.id}"
  cidrs = "192.168.10.1/32,192.168.20.1/32"
}
