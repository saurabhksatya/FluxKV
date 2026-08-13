# FluxKV - ⚠️ **Work in Progress**

A distributed key-value store built in Go with consistent hashing (Murmur3) and multi-leader replication.

## Architecture

FluxKV uses a **decentralized architecture** with no dedicated leader servers:

- **Murmur3 Hash Ring** - Consistent hashing for shard distribution across nodes
- **Multi-leader Replication** - Each server acts as:
  - **Leader** for a subset of shards
  - **Follower** for other shards
- **gRPC** - Inter-node communication and client API
- **YAML Configuration** - Flexible cluster configuration

## Project Structure

```
.
├── commands/       # CLI commands
├── configuration/  # Config parsing (YAML)
├── internal/       # Core internal packages
├── loadtest/       # Load testing utilities
├── replication/    # Replication logic (leader/follower)
├── server/         # gRPC server implementation
├── utils/          # Shared utilities (logger, etc.)
├── main.go         # Entry point
├── config.yaml     # Example configuration
└── docker-compose.yml
```

## Status

⚠️ **Work in Progress** - This project is under active development. APIs, configuration, and behavior may change.

## Getting Started

```bash
# Build
go build -o flux .

# Run with config
./flux --config config.yaml

# Or using Docker
docker-compose up
```

## Configuration

See `config.yaml` for cluster setup:

```yaml
server:
  port: 8000
  node_id: "node-1"

cluster:
  nodes:
    - "node-1:8000"
    - "node-2:8000"
    - "node-3:8000"
  replication_factor: 3
  virtual_nodes: 150
```

## Dependencies

- `github.com/shirou/gopsutil` - System metrics
- `google.golang.org/grpc` - RPC framework
- `google.golang.org/protobuf` - Protocol buffers
- `gopkg.in/yaml.v3` - YAML parsing
