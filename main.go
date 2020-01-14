package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ksyun.Provider,
	})
}
