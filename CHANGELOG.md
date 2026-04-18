# Changelog

## [Unreleased]

### Added

- Initial provider scaffold built from the HashiCorp
  `terraform-provider-scaffolding-framework` template.
- OpenAPI-driven code generation pipeline (`tfplugingen-openapi` +
  `tfplugingen-framework`) wired up via `make generate`.
- `anytype` provider with `endpoint`, `api_key`, and `api_version`
  configuration (also readable from `ANYTYPE_ENDPOINT`, `ANYTYPE_API_KEY`,
  and `ANYTYPE_API_VERSION` environment variables).
- `anytype_space` resource (create, read, update, import).
- `anytype_space` and `anytype_spaces` data sources.
- `anytype_type` resource (create, read, update, delete, import) with support
  for linking properties by `key`/`name`/`format`, plus matching
  `anytype_type` and `anytype_types` data sources.
- `anytype_property` resource (create, read, update, delete, import) with
  optional create-time `tags` seeding for `select` / `multi_select`
  formats, plus matching `anytype_property` and `anytype_properties` data
  sources.
- End-to-end examples under `examples/resources/anytype_type/` showing how
  to declare a space, properties, and a type with linked properties
  together.
- Thin typed HTTP client for the Anytype API with unit tests.
- GitHub Actions workflows for build + lint + unit tests on every push/PR,
  a reproducibility check that `make generate` produces no diff, and an
  acceptance test matrix gated on `TF_ACC` + `ANYTYPE_API_KEY`.
- GoReleaser-driven release workflow that produces Terraform Registry
  compatible artefacts (zipped multi-arch binaries, SHA256SUMS, detached
  GPG signatures) when a `v*` tag is pushed.
- `.golangci.yml` wired up to the same linter preset as the HashiCorp
  scaffolding template, plus a `make lint` / `make check` target.
- Dependabot config covering Go modules and GitHub Actions, and a
  repository pull request template.
- Daily `Sync Anytype API` GitHub Actions workflow that watches the upstream
  `versions.json`, runs `scripts/sync-anytype-api.sh` to bump the pinned
  spec version and regenerate the provider code, and opens a pull request
  when a new version is available.

### Changed

- Generated provider artefacts (`codegen/openapi.yaml`,
  `codegen/provider_code_spec.json`, everything under `internal/generated/`)
  are now gitignored. CI regenerates them on every build, test, release,
  and Anytype API sync job by calling `make generate` directly. The
  codegen CLIs (`tfplugingen-openapi`, `tfplugingen-framework`) are pinned
  as indirect dependencies in `tools/tools.go` and resolved via `go run`,
  so no separate install step is required.
- Dropped `.github/pull_request_template.md`.
