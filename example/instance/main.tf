# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_instance" "default" {
  image_id="6e37ed46-61a2-4f0a-9f4a-dcdd0817917c"
  instance_type="S4.1A"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  #only support part type
  data_disk =[
   {
      type="SSD3.0"
      size=20
      delete_with_instance=true
   }
 ]
  subnet_id="2ea4195d-8111-4cd8-91da-6a34bb06663b"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["6e3dee9c-291c-4647-bfc2-4c1eaa93fb80"]
  private_ip_address=""
  instance_name="xuan-tf"
  instance_name_suffix=""
  sriov_net_support="false"
  project_id=0
  data_guard_id=""
  key_id=[]
  d_n_s1 =""
  d_n_s2 =""
  force_delete =true
  user_data=""
  host_name="xuan-tf"
}
