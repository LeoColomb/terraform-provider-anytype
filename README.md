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
make generate   # fetch OpenAPI + regenerate schemas (required once)
go install
```

## Code generation

The provider schema is derived from the Anytype OpenAPI document. The
generated artefacts — `codegen/openapi.yaml`,
`codegen/provider_code_spec.json`, and everything under
`internal/generated/` — are **not** committed to git; they are produced
locally and in CI from [`codegen/generator_config.yml`](./codegen/generator_config.yml)
by running:

```sh
make generate-spec   # runs tfplugingen-openapi via `go run`
make generate-code   # runs tfplugingen-framework via `go run`
make generate        # full pipeline
```

Both codegen CLIs are pinned as indirect dependencies in
[`tools/tools.go`](./tools/tools.go), so `go run` resolves them at the
version recorded in `go.mod` — no separate `go install` step is needed.

CI runs `make generate` before every build, test, and release job, so
generated code is always produced from the currently pinned spec version.

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
