output "instance_id" {
  value = "${ksyun_instance.default.id}"
}
output "eip_id" {
  value = "${ksyun_eip.default.id}"
}
output "subnet_id" {
  value = "${ksyun_subnet.default.id}"
}
output "vpc_id" {
  value = "${ksyun_vpc.default.id}"
}