# Docker Setup for Note Web Application

This directory contains Docker configurations for running the Next.js web application in both production and development modes.

## Files Overview

- `Dockerfile` - Multi-stage production build
- `Dockerfile.dev` - Development environment with hot reloading
- `docker-compose.yml` - Orchestration for both production and development
- `.dockerignore` - Excludes unnecessary files from Docker build context

## Production Build

### Build and Run with Docker

```bash
# Build the production image
docker build -t note-web .

# Run the container
docker run -p 3000:3000 note-web
```

### Using Docker Compose

```bash
# Build and start the production service
docker-compose up --build

# Run in detached mode
docker-compose up -d

# Stop the service
docker-compose down
```

The application will be available at `http://localhost:3000`

## Development Mode

### Using Docker Compose for Development

```bash
# Start development service with hot reloading
docker-compose --profile dev up web-dev

# Build and start development service
docker-compose --profile dev up --build web-dev
```

The development server will be available at `http://localhost:3001` with hot reloading enabled.

## Multi-Stage Build Details

The production Dockerfile uses a 3-stage build process:

1. **Dependencies Stage**: Installs npm dependencies in an optimized way
2. **Builder Stage**: Builds the Next.js application with all optimizations
3. **Runner Stage**: Creates a minimal runtime image with only production assets

### Key Optimizations

- Uses Alpine Linux for smaller image size
- Leverages Next.js standalone output for minimal runtime
- Runs as non-root user for security
- Excludes dev dependencies from final image
- Implements proper layer caching for faster rebuilds

## Environment Variables

The following environment variables are set by default:

- `NODE_ENV=production` (for production builds)
- `NEXT_TELEMETRY_DISABLED=1` (disables Next.js telemetry)
- `PORT=3000` (application port)
- `HOSTNAME=0.0.0.0` (bind to all interfaces)

## Health Checks

The production service includes health checks that:
- Test the application endpoint every 30 seconds
- Allow 40 seconds for startup
- Retry 3 times before marking as unhealthy

## Security Features

- Non-root user execution
- Minimal attack surface with Alpine Linux
- No unnecessary packages in production image
- Proper file permissions

## Troubleshooting

### Build Issues

If you encounter build issues:

```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker build --no-cache -t note-web .
```

### Development Issues

For development mode issues:

```bash
# Check if ports are available
lsof -i :3001

# View container logs
docker-compose --profile dev logs web-dev

# Restart development container
docker-compose --profile dev restart web-dev
```

## Image Sizes

Expected image sizes:
- Development image: ~400MB (includes dev dependencies and tools)
- Production image: ~150MB (optimized for runtime only)

The production image is significantly smaller due to the multi-stage build process that excludes development dependencies and build tools.
