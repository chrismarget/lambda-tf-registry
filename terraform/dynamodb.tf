# curl 'https://registry.terraform.io/v1/providers/hashicorp/random/versions'
# curl 'https://registry.terraform.io/v1/providers/hashicorp/random/2.0.0/download/linux/amd64'

#         PK            SK
# namespace/type version/os/arch | protocols url shaUrl sigUrl sigKeys

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

locals {
  provider_details = {
    "4.0.1/linux/arm" = {
      sha     = "1aa2e4c07ddf87f7bda65a4a0f3b45c3edfbe983768d49a105f7ab9f2e4f8320"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_linux_arm.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/linux/amd64" = {
      sha     = "f80791f95f0ea5b332913e533c79ed4820e8c9243c508d8c7d6240b212160aaa"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_linux_amd64.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/windows/386" = {
      sha     = "3874421e4c975e987ade5bdece6d1eacd41065841c82856cc12fde405ea2fe38"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_windows_386.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/freebsd/arm" = {
      sha     = "1b7993daaf659dec421043ccf2dea021972ebacf47e5da3387e1ef35a0ffecbe"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_freebsd_arm.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/linux/arm64" = {
      sha     = "b4eb5438dc4bfbed7223c0044b775a210d52b631a9f37d884d567a3eacc31b92"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_linux_arm64.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/darwin/amd64" = {
      sha     = "b9808ee16fa06b7113a72c8d74f1cb322d0e7364fc34ba4bfdd0424ef7fd93d8"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_darwin_amd64.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/windows/amd64" = {
      sha     = "bdba092ae2939cb7e28380c5fd4a33ee96bead1abadbf9ec95d559cea8c04c3c"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_windows_amd64.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/freebsd/386" = {
      sha     = "4f27e1a90d779ac4bbdbd3db735b4777a90aefc8005905a8ed450bb517c323db"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_freebsd_386.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/darwin/arm64" = {
      sha     = "fe34ecc33c990f045ca5e3828e8aeb8ee86c9072e098e0ac0e4b47cbcb01edc0"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_darwin_arm64.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/freebsd/amd64" = {
      sha     = "1c40b056af93fe792fd468a96f317a6ce918849799906cf619a1b8cf01e79ccb"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_freebsd_amd64.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
    "4.0.1/linux/386" = {
      sha     = "bc5b1913fe841a0d40f28ff70d76e1c22fa3f469ae28011422d12c6001dcb954"
      url     = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_linux_386.zip"
      sha_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS"
      sig_url = "https://releases.hashicorp.com/terraform-provider-tls/4.0.1/terraform-provider-tls_4.0.1_SHA256SUMS.72D7468F.sig"
    }
  }
}

resource "aws_dynamodb_table_item" "providers" {
  for_each = local.provider_details
  table_name = aws_dynamodb_table.registry_providers.name
  hash_key   = aws_dynamodb_table.registry_providers.hash_key
  range_key = aws_dynamodb_table.registry_providers.range_key
  item = jsonencode({
    (aws_dynamodb_table.registry_providers.hash_key) = { S = "hashicorp/tls" }
    (aws_dynamodb_table.registry_providers.range_key) = { S = each.key }
    Protocols = { S = "[\"5.0\"]"}
    SHA = { S = each.value.sha }
    URL = { S = each.value.url }
    SHA_URL = { S = each.value.sha_url }
    Sig_URL = { S = each.value.sig_url }
    Keys = { S = file("signing_key")}
  })
}
