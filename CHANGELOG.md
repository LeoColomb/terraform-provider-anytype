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
