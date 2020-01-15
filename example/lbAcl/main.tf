# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Create Load Balancer Listener Acl
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "tf-xun2"
}
