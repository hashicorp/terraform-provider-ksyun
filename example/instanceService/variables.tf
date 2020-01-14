variable "availability_zone" {
  default = "cn-beijing-6a"
}
variable "security_group_name" {
  default = "ksyun_security_group_tf"
}
variable "instance_type" {
  default = ""
}
variable "instance_name" {
  default = "ksyun_instance_tf"
}
variable "subnet_name" {
  default = "ksyun_subnet_tf"
}
variable "vpc_name" {
  default = "ksyun_vpc_tf"
}

variable "vpc_cidr" {
  default = "10.1.0.0/21"
}
