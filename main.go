package main

import (
	"github.com/hashicorp/terraform/plugin"
	schemaregistry "github.com/trafi/terraform-provider-schema-registry/schema-registry"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: schemaregistry.Provider})
}
