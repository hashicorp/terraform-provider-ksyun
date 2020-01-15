provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.1.0.2/24"
}
