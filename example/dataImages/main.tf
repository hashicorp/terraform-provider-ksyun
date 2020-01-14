# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  ksyun_images
data "ksyun_images" "default" {
  output_file="output_result"
  ids=[]
  name_regex="centos-7.0-20180927115228"
  is_public=true
  image_source="system"
}

data "ksyun_images" "centos-7_5" {
  output_file="output_result1"
  platform= "centos-7.5"
  is_public=true
}
data "ksyun_images" "windows-2016-64-zh" {
  output_file="output_result1"
  platform= "windows-server_2016_datacenter_64_zh"
  is_public=true
}
data "ksyun_images" "windows-2016-64-en" {
  output_file="output_result1"
  platform= "windows-server_2016_datacenter_64_en"
  is_public=true
}

data "ksyun_images" "debian-8_2" {
  output_file="output_result1"
  platform= "debian-8.2"
  is_public=true
}

data "ksyun_images" "ubuntu-18_04" {
  output_file="output_result1"
  platform= "ubuntu-18.04"
  is_public=true
}


