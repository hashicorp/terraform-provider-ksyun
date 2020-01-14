# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Get  slbs
data "ksyun_slbs" "default" {
  output_file="output_result"

  ids=["d3fd0421-a35a-4ddb-a939-5c51e8af8e8c"]
  state=""
  vpc_id=[]
}

