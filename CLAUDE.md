# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Octopipe is a Go HTTP service that provides namespace allocation for branch-based deployments, integrating with Zadig CI/CD platform. It serves as a middleware for managing Kubernetes namespace resources based on Git branch names.

## Core Architecture

### Main Components

- **HTTP Server** (`pkg/api/server`): Gin-based HTTP server with health checks and graceful shutdown
- **Handler Layer** (`internal/handler`): Business logic for namespace allocation and Zadig integration
- **Cache Layer** (`internal/cache`): Thread-safe in-memory cache using sync.Map for branch-to-namespace mappings
- **Router** (`internal/router`): HTTP route definitions using Gin framework

### Key Workflows

1. **Namespace Allocation**: POST `/namespace_allocator` endpoint processes requests based on branch type:
   - `main/master`: Main branch handling
   - `release`: Release branch handling  
   - `dev/develop`: Development branch handling
   - `feat/*/feature/*`: Feature branch handling with dynamic namespace selection

2. **Feature Branch Logic**: 
   - Checks cache for existing branch-to-namespace mapping
   - If cache miss, queries Zadig for used namespaces via `/openapi/environments`
   - Selects available namespace from pool (test1-test50)
   - Caches the mapping for future requests

3. **Zadig Integration**: Service communicates with Zadig platform at `https://zadigx.shub.us/openapi/environments` to fetch environment status

## Development Commands

### Build and Run
```bash
go build -o octopipe .                    # Build binary
go run main.go                           # Run with default config
go run main.go --help                    # Show all available flags
```

### Configuration
```bash
# Run with custom profile and port
go run main.go --profile=prod --port=8080

# Run with config file
go run main.go --config-path=/path/to/config --config=config.yaml
```

### Testing
```bash
go test ./...                            # Run all tests
go test ./internal/cache                 # Run specific package tests
go test -v ./internal/cache              # Run with verbose output
go test -race ./...                      # Run with race detector
```

### Code Quality
```bash
go fmt ./...                             # Format code
go vet ./...                             # Run static analysis
go mod tidy                              # Clean up dependencies
go mod download                          # Download dependencies
```

## Application Configuration

The service accepts configuration via:
- Command line flags (see `--help`)
- Environment variables (automatically mapped)
- YAML config file (specified via `--config-path` and `--config`)

Key configuration options:
- `--host`: HTTP bind address (default: 127.0.0.1)
- `--port`: HTTP port (default: 6652)  
- `--profile`: Environment profile (dev/test/uat/prod)
- `--level`: Log level (debug/info/warn/error/fatal/panic)

## Kubernetes Deployment

Helm charts are available in `charts/octopipe/`:
```bash
helm install octopipe ./charts/octopipe
helm upgrade octopipe ./charts/octopipe
```

## Important Implementation Details

- Uses singleton pattern for branch-to-namespace cache (`internal/cache/branch_namespace.go:18`)
- Hard-coded JWT token in Zadig integration (`internal/handler/zadig.go:25`) - needs rotation management
- Namespace pool limited to test1-test50 range (`internal/handler/namespace_allocator.go:82`)
- HTTP client configured with 3 retries and exponential backoff (`internal/handler/handler.go:16`)
- Graceful shutdown with configurable timeout (`main.go:89`)
- User list injected via USERS_LIST environment variable (`main.go:50`)