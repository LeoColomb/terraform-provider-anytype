## Related Issue

<!-- If this PR resolves an issue, add "Fixes #<number>" here. -->

## Description

A plain-English summary of what the PR does and why. If it changes the
code-generation pipeline or the OpenAPI inputs, include a brief note about
how you regenerated and tested.

## Checklist

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes (unit tests)
- [ ] If provider schema changed: `make generate` re-run and the resulting
      diff committed
- [ ] Documentation updated under `docs/` and example(s) under `examples/`
- [ ] Entry added to `CHANGELOG.md` under `## [Unreleased]`
