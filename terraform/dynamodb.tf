resource "aws_dynamodb_table" "registry_providers" {
  name           = "registry-providers"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1

  hash_key  = "NamespaceType"
  range_key = "VersionOsArch"

  attribute {
    name = "NamespaceType"
    type = "S"
  }

  attribute {
    name = "VersionOsArch"
    type = "S"
  }
}
