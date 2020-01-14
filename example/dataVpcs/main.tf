# Specify the provider and access details
provider "ksyun" {
}

data "ksyun_vpcs" "default" {
  output_file="output_result"
  ids=["cf08e947-6577-44eb-80bc-2902d54818c3"]
}

