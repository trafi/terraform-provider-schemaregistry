package schemaregistry

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCHEMAREGISTRY_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"schemaregistry_subject_schema": resourceSchemaRegistrySubjectSchema(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	config := Config{
		URL: data.Get("url").(string),
	}

	log.Println("[INFO] Initializing Schema Registry client")

	return config.Client()
}
