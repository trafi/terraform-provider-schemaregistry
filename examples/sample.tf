provider "schemaregistry" {
  version = "1.0.0"
  url = "https://username:password@hostname"
}

resource "schemaregistry_subject_schema" "kafka_schemas" {
  for_each = {
    "rafal-test-4": "\"string\""
  }
  subject = each.key
  schema = each.value

  //compatibility_level = var.schema_compatibility
}