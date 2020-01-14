# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_images" "centos-7_5" {
  output_file=""
  platform= "centos-7.5"
  is_public=true
}
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
data "ksyun_lines" "default" {
  output_file=""
  line_name="BGP"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}
resource "ksyun_subnet" "default" {
  subnet_name      = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="${var.security_group_name}"
}
resource "ksyun_security_group" "default2" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="${var.security_group_name}"
}
resource "ksyun_security_group_entry" "test1" {
  description = "test1"
  security_group_id="${ksyun_security_group.default.id}"
  cidr_block="10.0.1.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
resource "ksyun_security_group_entry" "test2" {
  description = "test2"
  security_group_id="${ksyun_security_group.default.id}"
  cidr_block="10.0.1.6/32"
  direction="out"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
resource "ksyun_security_group_entry" "test3" {
  description = "test3"
  security_group_id="${ksyun_security_group.default2.id}"
  cidr_block="10.0.1.6/32"
  direction="out"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
resource "ksyun_ssh_key" "default" {
  key_name="ssh_key_tf"
  public_key=""
}
resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  data_disk =[
    {
      type="SSD3.0"
      size=20
      delete_with_instance=true
    }
  ]
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}","${ksyun_security_group.default2.id}"]
  private_ip_address=""
  instance_name="xuan-tf-combine"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=["${ksyun_ssh_key.default.id}"]
  force_delete=true
}

resource "ksyun_eip" "default" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =1
  charge_type = "PostPaidByDay"
  purchase_time =1
  project_id=0
}
resource "ksyun_eip_associate" "default" {
  allocation_id="${ksyun_eip.default.id}"
  instance_type="Ipfwd"
  instance_id="${ksyun_instance.default.id}"
  network_interface_id="${ksyun_instance.default.network_interface_id}"
}
