# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  instances
data "ksyun_instances" "default" {
  output_file = "output_result"
  ids = []
  name_regex=""
  project_id = []
  network_interface {
    network_interface_id = []
    subnet_id = []
    group_id = []
  }
  instance_state {
    name =  []
  }
  availability_zone {
    name =  []
  }
}

