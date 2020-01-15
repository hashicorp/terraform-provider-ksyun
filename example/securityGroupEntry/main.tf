# Specify the provider and access details
provider "ksyun" {
  region = "eu-east-1"
}
resource "ksyun_security_group_entry" "default1" {
  description = ""
  security_group_id="7385c8ea-79f7-4e9c-b99f-517fc3726256"
  cidr_block="10.0.0.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
resource "ksyun_security_group_entry" "default2" {
  description = ""
  security_group_id="7385c8ea-79f7-4e9c-b99f-517fc3726256"
  cidr_block="10.0.0.1/32"
  direction="out"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
