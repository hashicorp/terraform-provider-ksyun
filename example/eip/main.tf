# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}

# Create an eip
resource "ksyun_eip" "default1" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =1
  charge_type = "PostPaidByDay"
  purchase_time =1
  project_id=0
}
