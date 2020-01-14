# Specify the provider and access details
provider "ksyun" {
  region = "eu-east-1"
}
# Create an eip
resource "ksyun_security_group" "default" {
  vpc_id = "26231a41-4c6b-4a10-94ed-27088d5679df"
  security_group_name="xuan-tf--s"
}
