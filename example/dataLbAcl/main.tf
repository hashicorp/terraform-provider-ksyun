# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Get  lb acls
data "ksyun_lb_acls" "default" {
  output_file="output_result"
  ids=[]

}

