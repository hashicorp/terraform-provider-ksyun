provider "ksyun" {
  region = "cn-beijing-6"
  access_key = "your ak"
  secret_key = "your sk"
}

resource "ksyun_epc" "default" {
  host_name = "eeeee-test_4"
  host_type = "SSD"
  image_id = "2c9d8f29-6eb9-4bc7-90e5-b0bd7a9e2d3a"
  key_id = "0617edba-18fe-4a04-aed1-3b0874157b75"
  network_interface_mode = "bond4"
  raid = "Raid5"
  availability_zone = "cn-shanghai-3b"
  charge_type = "PostPaidByDay"
  security_agent = "classic"
  cloud_monitor_agent = "classic"
  subnet_id = "ecc2aeb2-c933-4fe9-a58d-9dd507be551c"
  security_group_id = ["9493e7f0-cbcf-4809-a79c-decc5db0bd8e"]
  network_interface_id = "07e83de8-c458-43ef-a2dd-7dbf3e52e31e"
  dns1 = "198.18.224.10"
  dns2 = "198.18.224.11"
//  private_ip_address = "10.0.80.14"
//  dns1 = "198.18.224.10"
//  dns2 = "198.18.224.11"
  password = "eeeade33ddAAA@ee4"
}


