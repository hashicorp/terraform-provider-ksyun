# Specify the provider and access details
provider "ksyun" {
  region = "eu-east-1"
}

resource "ksyun_eip_associate" "slb" {
  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
  instance_type="Slb"
  instance_id="7fae85e4-ab1a-415c-aef9-03a402c79d97"
  network_interface_id=""
}
resource "ksyun_eip_associate" "server" {
  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
  instance_type="Ipfwd"
  instance_id="566567677-6766-4743-afb7-7c7081214092"
  network_interface_id="87945980-59659-04548-759045803"
}
