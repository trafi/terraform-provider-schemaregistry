package schemaregistry

import (
	"fmt"
	"github.com/Landoop/schema-registry"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
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

func handleNotFoundError(err error, rd *schema.ResourceData) error {
	if srerr, ok := err.(*schemaregistry.ResourceError); ok && srerr.ErrorCode == 404 {
		log.Printf("[WARN] Removing %s from Terraform state because it's gone", rd.Id())
		rd.SetId("")
		return nil
	}

	return fmt.Errorf("Error reading: %s: %s", rd.Id(), err)
}
