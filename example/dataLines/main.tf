# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_lines" "default" {
  output_file="output_result"
  line_name="BGP"
}

