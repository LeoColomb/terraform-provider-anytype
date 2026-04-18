SHELL := /usr/bin/env bash

CODEGEN_DIR       ?= codegen
OPENAPI_SPEC      ?= $(CODEGEN_DIR)/openapi.yaml
OPENAPI_URL       ?= https://raw.githubusercontent.com/anyproto/anytype-api/main/docs/reference/openapi-2025-11-08.yaml
GENERATOR_CONFIG  ?= $(CODEGEN_DIR)/generator_config.yml
PROVIDER_SPEC     ?= $(CODEGEN_DIR)/provider_code_spec.json
GEN_RESOURCES_OUT    ?= internal/generated/resource_schemas
GEN_DATASOURCES_OUT  ?= internal/generated/datasource_schemas
GEN_PROVIDER_OUT     ?= internal/generated/provider_schema

default: fmt lint build test

build:
	go build -v ./...

install: build
	go install -v ./...

tidy:
	go mod tidy

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

check: fmt lint test

## --- Code generation ------------------------------------------------------

$(OPENAPI_SPEC):
	mkdir -p $(CODEGEN_DIR)
	curl -sSL $(OPENAPI_URL) -o $(OPENAPI_SPEC)

fetch-spec: $(OPENAPI_SPEC)

# Generate the Provider Code Specification (IR JSON) from the OpenAPI document.
generate-spec: fetch-spec
	tfplugingen-openapi generate \
		--config $(GENERATOR_CONFIG) \
		--output $(PROVIDER_SPEC) \
		$(OPENAPI_SPEC)

# Generate Terraform Plugin Framework Go code (schemas + models) from the IR.
# Resource, data source, and provider schemas are emitted into separate Go
# packages so that type names do not collide between resource and data source
# flavours of the same schema.
generate-code:
	tfplugingen-framework generate resources \
		--input $(PROVIDER_SPEC) \
		--output $(GEN_RESOURCES_OUT) \
		--package resource_schemas
	tfplugingen-framework generate data-sources \
		--input $(PROVIDER_SPEC) \
		--output $(GEN_DATASOURCES_OUT) \
		--package datasource_schemas
	tfplugingen-framework generate provider \
		--input $(PROVIDER_SPEC) \
		--output $(GEN_PROVIDER_OUT) \
		--package provider_schema

generate: generate-spec generate-code fmt

docs:
	cd tools && go generate ./...

.PHONY: default build install tidy fmt lint check test testacc fetch-spec generate-spec generate-code generate docs
