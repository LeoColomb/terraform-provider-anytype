# Terraform Provider for Anytype

A [Terraform](https://www.terraform.io) provider for [Anytype](https://anytype.io), built on the
[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) and
generated from the official [Anytype OpenAPI specification](https://github.com/anyproto/anytype-api)
using the Terraform [provider code generation](https://developer.hashicorp.com/terraform/plugin/code-generation)
toolchain.

This provider lets you manage Anytype resources — spaces, types, properties, tags and objects —
declaratively from Terraform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.6
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider)
- A running Anytype API (the desktop application exposes it locally by default on
  `http://127.0.0.1:31009`)
- An Anytype API key — obtained via the `/v1/auth/challenges` + `/v1/auth/api_keys` flow

## Usage

```hcl
terraform {
  required_providers {
    anytype = {
      source  = "LeoColomb/anytype"
      version = "~> 0.1"
    }
  }
}

provider "anytype" {
  endpoint = "http://127.0.0.1:31009" # optional, defaults to local Anytype
  api_key  = var.anytype_api_key       # or set ANYTYPE_API_KEY
}

resource "anytype_space" "wiki" {
  name        = "Engineering Wiki"
  description = "The local-first engineering wiki"
}
```

See [`examples/`](./examples) for more.

## Building

```sh
go install
```

## Code generation

The provider schema is derived from the Anytype OpenAPI document. To regenerate it:

```sh
make generate-spec   # runs tfplugingen-openapi
make generate-code   # runs tfplugingen-framework
make generate        # full pipeline + docs
```

The generator configuration lives in [`codegen/generator_config.yml`](./codegen/generator_config.yml)
and the intermediate provider-code specification is committed at
[`codegen/provider_code_spec.json`](./codegen/provider_code_spec.json).

## Testing

Unit tests:

```sh
make test
```

Acceptance tests require a running Anytype instance and `TF_ACC=1`:

```sh
TF_ACC=1 ANYTYPE_API_KEY=xxx make testacc
```

## License

Mozilla Public License v2.0. See [LICENSE](./LICENSE).
