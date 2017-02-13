package main

import (
	"github.com/hashicorp/terraform/builtin/providers/arukas"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: arukas.Provider,
	})
}
