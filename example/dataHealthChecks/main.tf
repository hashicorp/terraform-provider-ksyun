# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Get  healthchecks
data "ksyun_health_checks" "default" {
  output_file="output_result"

  ids=[]
  listener_id=["8d1dac22-6c6c-42ea-93e2-2702d44ddb93","70467f7e-23dc-465a-a609-fb1525fc6b16"]

}

