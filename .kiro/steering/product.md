# Amazon CloudWatch Agent

The Amazon CloudWatch Agent is a telemetry collection agent that enables comprehensive monitoring and observability for AWS environments. It collects system-level metrics, custom application metrics, logs, and traces from EC2 instances, on-premises servers, and containerized environments.

## Key Capabilities

- **Metrics Collection**: System-level metrics (CPU, memory, disk, network) and custom application metrics via StatsD and collectd protocols
- **Log Collection**: Application and system logs from Linux and Windows environments
- **Trace Collection**: OpenTelemetry and AWS X-Ray traces for distributed tracing
- **Multi-Platform Support**: Linux, Windows, macOS, and containerized deployments
- **Hybrid Environments**: Works with both AWS-managed and on-premises infrastructure

## Architecture

Built on top of OpenTelemetry Collector and Telegraf, the agent operates as a pipeline processor that can handle both telegraf and OpenTelemetry components alongside custom AWS-specific components. It uses a plugin-based architecture for extensibility and supports various input sources, processors, and output destinations.

## Target Environments

- Amazon EC2 instances
- On-premises servers
- Container environments (Docker, Kubernetes, ECS)
- Hybrid cloud deployments