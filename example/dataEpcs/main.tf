# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
  access_key = "your ak"
  secret_key = "your sk"
}

# Get  epcs

data "ksyun_epcs" "default" {
  ids = ["4c3656de-9907-4c0d-868e-969cee396f37"]
  //  project_id = ["1"]
  //  host_name = ["jeeep-test_1"]
  //  subnet_id = ["25f2947a-afdc-40ae-b4c9-cab16c38a63a"]
  //  vpc_id = ["5f103e58-c3c2-4af7-92fa-98b87abb4232"]
  //  cabinet_id = ["f174dceb-22ba-45be-941e-3f11eabd69f1"]
  //  host_type =["CAL","SSD"]
  //  epc_host_status = ["Stopping"]
  //  os_name = ["CentOS-7.2 64位1"]
  //  product_type = ["customer"]
  //  cluster_id = ["f174dceb-22ba-45be-941e-3f11eabd69f1"]
  //  enable_container = ["false"]
  output_file = "output_result"
}

