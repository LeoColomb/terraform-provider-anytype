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
