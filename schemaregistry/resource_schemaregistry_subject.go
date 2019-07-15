package schemaregistry

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"

	"github.com/Landoop/schema-registry"
)

func resourceSchemaRegistrySubjectSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchemaRegistrySubjectSchemaCreate,
		Read:   resourceSchemaRegistrySubjectSchemaRead,
		Update: resourceSchemaRegistrySubjectSchemaUpdate,
		Delete: resourceSchemaRegistrySubjectSchemaDelete,
		Schema: map[string]*schema.Schema{
			"subject": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
			},
		},
	}
}

func resourceSchemaRegistrySubjectSchemaCreate(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	subject := rd.Get("subject").(string)
	schema := rd.Get("schema").(string)

	log.Printf("[INFO] Creating Schema Registry schema for subject %s", subject)

	schemaID, err := client.RegisterNewSchema(subject, schema)
	if err != nil {
		return err
	}

	rd.SetId(strconv.Itoa(schemaID))

	return resourceSchemaRegistrySubjectSchemaRead(rd, meta)
}

func resourceSchemaRegistrySubjectSchemaRead(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	schemaID := rd.Id()

	log.Printf("[INFO] Reading Schema Registry schema for id %s", schemaID)

	schemaId, err := strconv.Atoi(schemaID)

	if err != nil {
		return err
	}

	schema, err := client.GetSchemaByID(schemaId)
	if err != nil {
		return handleNotFoundError(err, rd)
	}

	err = rd.Set("schema", schema)

	return nil
}

func resourceSchemaRegistrySubjectSchemaUpdate(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	schema := rd.Get("schema").(string)
	subject := rd.Get("subject").(string)

	log.Printf("[INFO] Updating Schema Registry schema for subject '%s'", subject)

	schemaID, err := client.RegisterNewSchema(subject, schema)
	if err != nil {
		return err
	}

	rd.SetId(strconv.Itoa(schemaID))

	return resourceSchemaRegistrySubjectSchemaRead(rd, meta)
}

func resourceSchemaRegistrySubjectSchemaDelete(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	subject := rd.Get("subject").(string)

	log.Printf("[INFO] Deleting Schema Registry subject '%s'", subject)

	if _, err := client.DeleteSubject(subject); err != nil {
		return err
	}

	rd.SetId("")

	return nil
}
