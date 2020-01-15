provider "ksyun" {
  region = "cn-beijing-6"
}
resource "ksyun_vpc" "test" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}
resource "ksyun_subnet" "test" {
  subnet_name      = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id  = "${ksyun_vpc.test.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "cn-beijing-6a"
}

resource "ksyun_security_group" "test" {
  vpc_id = "${ksyun_vpc.test.id}"
  security_group_name="${var.security_group_name}"
}
resource "ksyun_security_group_entry" "test1" {
  description = "26231a41-4c6b-4a10-94ed-27088d5679df"
  security_group_id="${ksyun_security_group.test.id}"
  cidr_block="10.0.1.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
resource "ksyun_security_group_entry" "test2" {
  description = "26231a41-4c6b-4a10-94ed-27088d5679df"
  security_group_id="${ksyun_security_group.test.id}"
  cidr_block="10.0.1.1/32"
  direction="out"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
