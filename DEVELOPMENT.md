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

## Development Commands

### Setup & Installation

| Command | Description |
|---------|-------------|
| `make setup-dev` | Setup complete development environment |
| `make install` | Install dependencies for development |
| `make create-dev-compose` | Create development docker-compose override file |
| `make create-air-config` | Create Air config for Go hot reloading |
| `make create-dev-dockerfiles` | Create development Dockerfiles |

### Development Environment

| Command | Description |
|---------|-------------|
| `make dev` | Start development environment with hot reloading |
| `make dev-build` | Build and start development environment |
| `make dev-down` | Stop development environment |
| `make dev-server` | Start only server in development mode |
| `make dev-web` | Start only web in development mode |
| `make quick-dev` | Quick dev: build and run in development mode |

### Building & Deployment

| Command | Description |
|---------|-------------|
| `make build` | Build all Docker images |
| `make build-server` | Build only the server Docker image |
| `make build-web` | Build only the web Docker image |
| `make up` | Start all services in production mode |
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

1. Use `make dev` for daily development
2. Use `make health` to check if services are running
3. Use `make logs` to debug issues
4. Use `make clean` when you have Docker issues
5. Run `make help` anytime to see available commands
