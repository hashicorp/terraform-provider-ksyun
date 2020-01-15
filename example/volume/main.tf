# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region = "cn-shanghai-3"
}

resource "ksyun_volume" "default" {
  volume_name="test"
  volume_type="SSD3.0"
  size=15
  charge_type="Daily"
  availability_zone="cn-shanghai-3a"
  volume_desc="test"
}