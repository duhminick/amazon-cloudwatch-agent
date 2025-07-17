# Technology Stack

## Core Technologies

- **Language**: Go 1.24.4
- **Build System**: Make-based build system with cross-platform compilation
- **Dependencies**: Go modules for dependency management
- **Architecture**: Plugin-based system built on OpenTelemetry Collector and Telegraf

## Key Dependencies

- **OpenTelemetry Collector**: Core telemetry collection framework
- **Telegraf**: Metrics collection agent (forked version: `github.com/aws/telegraf`)
- **AWS SDK Go**: AWS service integration (`github.com/aws/aws-sdk-go`)
- **Kubernetes Client**: Container orchestration support (`k8s.io/client-go`)

## Build Commands

### Development
```bash
# Build for all platforms
make build

# Build and package (default target)
make release

# Clean build artifacts
make clean

# Run tests
make test

# Run tests with race detection
make test-data-race

# Format code
make fmt

# Run linter
make lint
```

### Docker
```bash
# Build Docker image from source
make dockerized-build

# Build for specific architectures
make docker-build-amd64
make docker-build-arm64
```

### Platform-Specific Builds
```bash
# Linux builds (AMD64/ARM64)
make amazon-cloudwatch-agent-linux

# Windows builds
make amazon-cloudwatch-agent-windows  

# macOS builds
make amazon-cloudwatch-agent-darwin
```

## Code Quality Tools

- **Linter**: golangci-lint v1.64.2
- **Formatter**: gofmt + goimports
- **License Check**: addlicense tool
- **Import Order**: impi tool
- **Shell Formatting**: shfmt

## Build Modes

- Default: PIE (Position Independent Executable) mode
- Configurable via `CWAGENT_BUILD_MODE` environment variable

## Cross-Platform Support

- Linux: AMD64, ARM64 (RPM, DEB packages)
- Windows: AMD64 (ZIP package)  
- macOS: AMD64, ARM64 (TAR.GZ package)
- Docker: Multi-architecture support