# Specify the provider and access details
provider "ksyun" {
}

# Get  ksc_subnet_allocated_ip_addresses
data "ksyun_subnet_allocated_ip_addresses" "default" {
  output_file="output_result"
  ids=["494c3a64-eff9-4438-aa7c-694b7ba472d5"]
  subnet_id=["494c3a64-eff9-4438-aa7c-694b7ba472d5"]
}

