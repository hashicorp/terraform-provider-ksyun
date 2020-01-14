# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_ssh_key" "default1" {
  key_name="ssh_key_tf"
  public_key="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCMJ0oADMm0dgfdfgfdghfhRFEREHHnaEye9NZJvj02gYQcz6dRSivaibOk6XU/tbuzCpKr+eyuxmOYnwfuN2efat83WEphZHT74RozKaroSI6XF8jHSbKxPdhKt5L+E9cS5pGCxp/zUhAWChBaDT8GYQ== root\n"
}
