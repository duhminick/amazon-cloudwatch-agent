# Project Structure

## Root Directory Organization

### Core Application
- **`cmd/`** - Main application entry points and executables
  - `amazon-cloudwatch-agent/` - Main agent executable
  - `amazon-cloudwatch-agent-config-wizard/` - Configuration wizard
  - `config-downloader/` - Configuration download utility
  - `config-translator/` - Configuration translation utility
  - `start-amazon-cloudwatch-agent/` - Agent startup utility

### Source Code
- **`internal/`** - Private application packages (not importable by external projects)
  - `constants/` - Application constants
  - `ec2metadataprovider/` - EC2 metadata integration
  - `ecsservicediscovery/` - ECS service discovery
  - `k8sCommon/` - Kubernetes utilities
  - `util/` - Common utilities
- **`cfg/`** - Configuration management
  - `aws/` - AWS-specific configuration
  - `commonconfig/` - Shared configuration
  - `envconfig/` - Environment-based configuration

### Plugin Architecture
- **`plugins/`** - CloudWatch Agent specific plugins
  - `inputs/` - Input plugins (logfile, nvidia_smi, prometheus, statsd, etc.)
  - `outputs/` - Output plugins (cloudwatch, cloudwatchlogs)
  - `processors/` - Processing plugins (awsapplicationsignals, ec2tagger, etc.)
- **`receiver/`** - OpenTelemetry receivers
- **`processor/`** - OpenTelemetry processors
- **`extension/`** - OpenTelemetry extensions

### Translation & Configuration
- **`translator/`** - Configuration translation logic
  - `translate/` - Translation implementations
  - `tocwconfig/` - CloudWatch config generation
  - `util/` - Translation utilities
- **`tool/`** - Command-line tools and utilities

### Infrastructure
- **`service/`** - Service layer components
- **`sdk/`** - AWS SDK customizations
- **`packaging/`** - Platform-specific packaging files
  - `debian/` - Debian package configuration
  - `linux/` - RPM package configuration  
  - `windows/` - Windows package configuration
  - `darwin/` - macOS package configuration

### Container Support
- **`amazon-cloudwatch-container-insights/`** - Container-specific implementations
  - `cloudwatch-agent-dockerfile/` - Docker build configurations
  - `k8s-yaml-templates/` - Kubernetes deployment templates

### Build & Development
- **`build/`** - Build artifacts (generated)
- **`Tools/`** - Build and packaging scripts
- **`licensing/`** - License files

## Naming Conventions

### Go Packages
- Use lowercase package names
- Avoid underscores in package names (except for specific cases like `windows_events`)
- Use descriptive names that reflect functionality

### File Organization
- Group related functionality in packages
- Separate platform-specific code with build tags or separate files
- Use `_test.go` suffix for test files
- Use `_unix.go`, `_windows.go`, `_darwin.go` for platform-specific implementations

### Import Organization
- Standard library imports first
- Third-party imports second  
- Local imports last (prefixed with `github.com/aws/amazon-cloudwatch-agent`)
- Use `goimports` tool to maintain proper import order

## Configuration Files
- **`go.mod`** - Go module definition with extensive replace directives for AWS forks
- **`Makefile`** - Build system configuration
- **`.golangci.yml`** - Linter configuration
- **`codecov.yml`** - Code coverage configuration