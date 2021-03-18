package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/trafi/terraform-provider-schemaregistry/schemaregistry"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: schemaregistry.Provider})
}
