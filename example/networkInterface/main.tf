provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_network_interface" "test" {
  instance_id   = "3795bb91-b08d-44ac-b58f-0b5ea697204f"
  network_interface_id = "9018071c-dce4-430b-a2bb-025ca9a07b52"
  security_group_id=["ca4d9b07-05fe-4459-9104-c81bac529aa3","db0e7b26-71b8-4d61-9e4b-4c97923502bf"]
  subnet_id="4a9b6f43-d728-46fa-8c58-a7a13c9ff8d7"
  private_ip_address="10.0.55.5"
}

