package schemaregistry

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"

	"github.com/Landoop/schema-registry"
)

// Terraform resource ID separator
const IDSeparator = "___"

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
					newJSON, _ := structure.NormalizeJsonString(new)
					oldJSON, _ := structure.NormalizeJsonString(old)
					return newJSON == oldJSON
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

	_, err := client.RegisterNewSchema(subject, schema)
	if err != nil {
		return err
	}

	schemaDefinition, err := client.GetLatestSchema(subject)
	if err != nil {
		return err
	}

	rd.SetId(subject + IDSeparator + strconv.Itoa(schemaDefinition.Version))

	return resourceSchemaRegistrySubjectSchemaRead(rd, meta)
}

func resourceSchemaRegistrySubjectSchemaRead(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	ID := strings.Split(rd.Id(), IDSeparator)
	subject := ID[0]
	schemaID := ID[1]

	log.Printf("[INFO] Reading Schema Registry schema for subject %s with id %s", subject, schemaID)

	subjectSchemaID, err := strconv.Atoi(schemaID)

	if err != nil {
		return err
	}

	schema, err := client.GetSchemaBySubject(subject, subjectSchemaID)
	if err != nil {
		log.Printf("[WARN] Removing %s from Terraform state because it's gone", rd.Id())
		rd.SetId("")
		return nil
	}

	err = rd.Set("schema", schema)

	return nil
}

func resourceSchemaRegistrySubjectSchemaUpdate(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	schema := rd.Get("schema").(string)
	subject := rd.Get("subject").(string)

	log.Printf("[INFO] Updating Schema Registry schema for subject '%s'", subject)

	_, err := client.RegisterNewSchema(subject, schema)
	if err != nil {
		return err
	}

	schemaDefinition, err := client.GetLatestSchema(subject)
	if err != nil {
		return err
	}

	rd.SetId(subject + IDSeparator + strconv.Itoa(schemaDefinition.Version))

	return resourceSchemaRegistrySubjectSchemaRead(rd, meta)
}

func resourceSchemaRegistrySubjectSchemaDelete(rd *schema.ResourceData, meta interface{}) error {
	client := meta.(*schemaregistry.Client)

	ID := strings.Split(rd.Id(), IDSeparator)
	subject := ID[0]

	log.Printf("[INFO] Deleting Schema Registry subject '%s'", subject)

	if _, err := client.DeleteSubject(subject); err != nil {
		return err
	}

	rd.SetId("")

	return nil
}
