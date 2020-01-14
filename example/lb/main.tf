# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Create Load Balancer
resource "ksyun_lb" "default" {
  vpc_id = "74d0a45b-472d-49fc-84ad-221e21ee23aa"
  load_balancer_name = "tf-xun1"
  type = "public"
  subnet_id = "609d1736-d8d7-492d-abd3-1183bb60329e"
  load_balancer_state = "stop"
  private_ip_address = "10.0.77.11"
}
