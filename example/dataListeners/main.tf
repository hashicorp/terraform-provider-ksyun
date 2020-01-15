# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Get  listeners
data "ksyun_listeners" "default" {
  output_file="output_result"
  ids=[""]
  load_balancer_id=["d3fd0421-a35a-4ddb-a939-5c51e8af8e8c","4534d617-9de0-4a4a-9ed5-3561196cacb6"]
}


