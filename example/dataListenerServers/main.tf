# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "cn-beijing-6"
}

# Get  listener_servers
data "ksyun_listener_servers" "default" {
  output_file="output_result"

  ids=[]
  listener_id=[]
  real_server_ip=["10.72.20.126","172.31.16.20"]
}

