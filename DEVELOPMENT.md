# Development Guide

This project includes a comprehensive Makefile for easy development and deployment.

## Quick Start

```bash
# Setup development environment (first time only)
make setup-dev

# Start development with hot reloading
make dev

# Or start production mode
make quick-start
```

## Available Commands

Run `make help` to see all available commands:

```bash
make help
```

## Development Workflow

### First Time Setup
```bash
# Install dependencies and setup development files
make setup-dev
```

### Daily Development
```bash
# Start development with hot reloading
make dev

# Or build and start development
make dev-build
```

### Production Testing
```bash
# Build and start production containers
make build up

# Quick start (build + up)
make quick-start
```

### Debugging
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

### Individual Services
```bash
# Build individual services
make build-server
make build-web

# Rebuild individual services
make rebuild-server
make rebuild-web
```

### Cleanup
```bash
# Clean up Docker resources
make clean

# Clean everything (careful!)
make clean-all
```

## Hot Reloading

The development environment includes:

- **Server**: Uses Air for Go hot reloading
- **Web**: Uses Next.js built-in hot reloading
- **Volumes**: Source code is mounted for instant updates

## Development vs Production

### Development Mode (`make dev`)
- Hot reloading enabled
- Source code mounted as volumes
- Debug logging enabled
- Development Dockerfiles used

### Production Mode (`make up`)
- Optimized builds
- No hot reloading
- Production Dockerfiles used
- Minimal container sizes

## Port Mapping

- **Web App**: http://localhost:3000
- **Server API**: http://localhost:8080

## Tips

1. Use `make dev` for daily development
2. Use `make health` to check if services are running
3. Use `make logs` to debug issues
4. Use `make clean` when you have Docker issues
5. Run `make help` anytime to see available commands
