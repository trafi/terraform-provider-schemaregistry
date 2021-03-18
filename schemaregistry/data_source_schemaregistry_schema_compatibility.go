package schemaregistry

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"

	"github.com/Landoop/schema-registry"
)

func dataSourceSchemaRegistrySchemaCompatibility() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSchemaRegistrySchemaCompatibilityRead,

		Schema: map[string]*schema.Schema{
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Schema Registry subject to validate new schema against",
			},
			"schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "New schema version to validate",
			},
			"compatible": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if new schema version is compatible with subject's latest schema version",
			},
		},
	}
}

func dataSourceSchemaRegistrySchemaCompatibilityRead(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	subject := rd.Get("subject").(string)
	schema := rd.Get("schema").(string)

	log.Printf("[INFO] Checking compatibility for '%s' with latest schema", subject)
	compatibility, err := client.IsLatestSchemaCompatible(subject, schema)

	log.Printf("[INFO] Schema for '%s' is compatible: '%t' (%v)", subject, compatibility, err)

	if srerr, ok := err.(schemaregistry.ResourceError); ok {
		// Missing schema means this is an initial schema submission therefore
		// it's compatible by default
		if srerr.ErrorCode == 40401 {
			compatibility = true
		} else {
			return err
		}
	}

	rd.SetId(subject)
	rd.Set("compatible", compatibility)

	return nil
}
