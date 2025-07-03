# Development Guide

This project consists of two main components: a server and a web application.

## Architecture Overview

The application is built using a microservices architecture with the following components:

- **Server**: Backend API service built with Go
- **Web**: Frontend application built with Next.js

## Quick Start

```bash
# Build and start all services
make build up

# Or use quick start
make quick-start
```

## Available Commands

Run `make help` to see all available commands:

```bash
make help
```

## Production Testing
```bash
# Build and start production containers
make build up

# Quick start (build + up)
make quick-start
```

## Debugging
```bash
# View logs
make logs
make logs-server  # Server only
make logs-web     # Web only

# Check service status
make status

# Health check
make health
```

## Individual Services
```bash
# Build individual services
make build-server
make build-web

# Rebuild individual services
make rebuild-server
make rebuild-web
```

## Cleanup
```bash
# Clean up Docker resources
make clean

# Clean everything (careful!)
make clean-all
```

## Port Mapping

- **Web App**: http://localhost:3000
- **Server API**: http://localhost:8080

## Command Reference

### Setup & Installation

| Command | Description |
|---------|-------------|
| `make install` | Install dependencies |

### Building & Deployment

| Command | Description |
|---------|-------------|
| `make build` | Build all Docker images |
| `make build-server` | Build only the server Docker image |
| `make build-web` | Build only the web Docker image |
| `make up` | Start all services |
| `make down` | Stop all services |
| `make restart` | Restart all services |
| `make quick-start` | Quick start: build and run |

### Rebuilding & Maintenance

| Command | Description |
|---------|-------------|
| `make rebuild` | Rebuild and restart all services |
| `make rebuild-server` | Rebuild and restart only server |
| `make rebuild-web` | Rebuild and restart only web |

### Monitoring & Debugging

| Command | Description |
|---------|-------------|
| `make logs` | Show logs from all services |
| `make logs-server` | Show logs from server only |
| `make logs-web` | Show logs from web only |
| `make status` | Show status of all services |
| `make health` | Check if services are healthy |
| `make shell-server` | Open shell in server container |
| `make shell-web` | Open shell in web container |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run tests |
| `make test-web` | Run web tests only |
| `make test-server` | Run server tests only |

### Cleanup

| Command | Description |
|---------|-------------|
| `make clean` | Clean up Docker resources (containers, images, volumes) |
| `make clean-all` | Clean up everything including images |

### Help

| Command | Description |
|---------|-------------|
| `make help` | Show this help message |

**Note:** Running `make help` prints the same information with descriptions for all available commands.

## Tips

1. Use `make quick-start` to build and start all services
2. Use `make health` to check if services are running
3. Use `make logs` to debug issues
4. Use `make clean` when you have Docker issues
5. Run `make help` anytime to see available commands
