# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Git Repositories

**Working repositories:**
- https://github.com/jniedergang/terraform-provider-harvester
- https://gitea.home.zypp.fr/jniedergang/terraform-provider-harvester

**Reference (upstream):**
- https://github.com/harvester/terraform-provider-harvester

## Build Commands

```bash
# Build provider binaries (amd64 and arm64)
make build

# Run tests with coverage (generates coverage.html)
make test

# Run linters and formatters (golangci-lint, go fmt, go generate)
make validate

# Generate documentation (terraform-plugin-docs)
make generate

# Interactive dev helper
./dev-env.sh
```

### Acceptance Tests

Acceptance tests require a `kubeconfig_test.yaml` file in the repository root:
```bash
TF_ACC=1 go test ./internal/tests/... -v
```

## Architecture

### Provider Overview

Terraform provider for Harvester HCI (Hyper-Converged Infrastructure). Uses HashiCorp Terraform Plugin SDK v2 with Kubernetes client-go for API interactions.

- Minimum Terraform: >= 0.13.x
- Minimum Harvester: v1.1.0 (v1.0.x not supported)

### Project Structure

```
internal/
├── config/          # Provider configuration and K8s client initialization
├── provider/        # Provider definition and all resources/data sources
│   └── <resource>/  # Each resource has its own directory
├── tests/           # Acceptance tests (resource_*_test.go)
└── util/            # Constructor pattern, schema utilities, state management

pkg/
├── client/          # Multi-client aggregator (KubeVirt, Harvester CRDs, etc.)
├── constants/       # Field and resource type constants
├── helper/          # ID/naming utilities (namespace/name format)
└── importer/        # Resource importers
```

### Resource Pattern

Each resource follows this structure in `internal/provider/<resource>/`:

| File | Purpose |
|------|---------|
| `schema_<name>.go` | Schema definition with field types and validation |
| `resource_<name>.go` | CRUD operations (Create, Read, Update, Delete) |
| `resource_<name>_constructor.go` | Builds Kubernetes objects from Terraform state |
| `datasource_<name>.go` | Read-only data source implementation |
| `schema_<name>_*.go` | Nested schema definitions for complex types |

### Key Patterns

**Constructor Pattern** (`internal/util/constructor.go`):
- `Constructor` interface with `Setup()`, `Result()`, `Validate()` methods
- `Processor` pattern for field parsing with required/optional validation

**Schema Wrapping** (`internal/util/schema.go`):
- `NamespacedSchemaWrap()` - adds common fields (name, namespace, tags, labels, description, state, message)
- `NonNamespacedSchemaWrap()` - for cluster-level resources
- `DataSourceSchemaWrap()` - converts resource schema to read-only

**ID Format** (`pkg/helper/id.go`):
- Namespaced resources: `"namespace/name"`
- Non-namespaced: `"name"`
- Use `helper.BuildID()` and `helper.IDParts()`

**Constants** (`pkg/constants/`):
- Resource types: `constants.ResourceTypeVirtualMachine`
- Field names: `constants.FieldCommonName`, `constants.FieldVirtualMachineCPU`
- Define new constants in `constants_<resource>.go`

### Client Architecture

`pkg/client/client.go` aggregates multiple clients:
- `KubeClient` - standard Kubernetes
- `HarvesterClient` - Harvester CRDs
- `HarvesterNetworkClient` - network controller CRDs
- `HarvesterLoadbalancerClient` - load balancer CRDs
- `KubeVirtSubresourceClient` - KubeVirt subresources

### Testing

Tests use HashiCorp's testing framework in `internal/tests/`:
- `VMResourceBuilder` - fluent builder for VM test configs
- `testAccPreCheck()` - validates kubeconfig before tests
- Use `getStateChangeConf()` for waiting on resource deletion

## Adding a New Resource

1. Create directory `internal/provider/<resource>/`
2. Add constants in `pkg/constants/constants_<resource>.go`
3. Implement schema, resource operations, and constructor
4. Register in `internal/provider/provider.go` (ResourcesMap and DataSourcesMap)
5. Add acceptance tests in `internal/tests/resource_<name>_test.go`
6. Run `go generate` to update documentation
