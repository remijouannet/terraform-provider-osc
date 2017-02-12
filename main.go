package main

import (
	"github.com/remijouannet/terraform-provider-osc/osc"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: osc.Provider,
	})
}
