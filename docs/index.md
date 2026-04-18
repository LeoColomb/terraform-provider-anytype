# Anytype Provider

The Anytype provider lets you manage [Anytype](https://anytype.io) resources
declaratively from Terraform. It is built on the
[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
and generated from the
[official Anytype OpenAPI specification](https://github.com/anyproto/anytype-api)
using the HashiCorp
[code generation toolchain](https://developer.hashicorp.com/terraform/plugin/code-generation).

## Example Usage

```terraform
terraform {
  required_providers {
    anytype = {
      source  = "LeoColomb/anytype"
      version = "~> 0.1"
    }
  }
}

provider "anytype" {
  endpoint = "http://127.0.0.1:31009"
  api_key  = var.anytype_api_key
}

resource "anytype_space" "wiki" {
  name        = "Engineering Wiki"
  description = "The local-first engineering wiki"
}
```

## Schema

### Optional

- `endpoint` (String) — Anytype API endpoint. Defaults to the local desktop API
  (`http://127.0.0.1:31009`). Also read from `ANYTYPE_ENDPOINT`.
- `api_key` (String, Sensitive) — Anytype API key obtained via the
  `/v1/auth/api_keys` flow. Also read from `ANYTYPE_API_KEY`.
- `api_version` (String) — Value sent in the `Anytype-Version` header.
  Defaults to the version the provider was generated against. Also read
  from `ANYTYPE_API_VERSION`.

## Authentication

1. Start an authentication challenge:

   ```sh
   curl -X POST http://127.0.0.1:31009/v1/auth/challenges \
     -H 'Anytype-Version: 2025-11-08' \
     -d '{"app_name": "terraform"}'
   ```

2. Approve the challenge in the Anytype desktop app and exchange the 4-digit
   code for an API key:

   ```sh
   curl -X POST http://127.0.0.1:31009/v1/auth/api_keys \
     -H 'Anytype-Version: 2025-11-08' \
     -d '{"challenge_id": "...", "code": "1234"}'
   ```

3. Export the returned key:

   ```sh
   export ANYTYPE_API_KEY=...
   ```
