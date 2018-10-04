package main

import (
	"github.com/Brightspace/terraform-provider-cloudability/cloudability"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: cloudability.Provider})
}
