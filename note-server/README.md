# Note Server

A scalable backend server for the Note application, built with Go and designed for cloud-native deployment.

## Architecture

### Package Structure

The project follows Go best practices with clear separation between internal and public packages:

```
note-server/
├── cmd/
│   └── server/          # Application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── http/            # HTTP handlers and routing
│   ├── service/         # Business logic
│   ├── util/            # Internal utilities (deprecated - use pkg/)
│   └── ws/              # WebSocket handling
├── pkg/                 # Public, reusable packages
│   ├── response/        # HTTP response utilities
│   └── timeutil/        # Time formatting utilities
├── docs/                # Documentation
└── .github/workflows/   # CI/CD pipelines
```

### Design Principles

- **Package Privacy**: Keep `internal/` clean for application-specific code
- **Shared Code**: Put reusable utilities in `pkg/` for potential sharing
- **Future-Ready**: Prepared for `note-shared` module extraction

## Quick Start

### Local Development

1. **Clone and run:**
   ```bash
   git clone <repository-url>
   cd note-server
   go mod download
   go run ./cmd/server
   ```

2. **Test the API:**
   ```bash
   curl http://localhost:8080/healthz
   ```

### Docker

1. **Build and run with Docker:**
   ```bash
   docker build -t note-server .
   docker run -p 8080:8080 note-server
   ```

2. **Using Docker Compose (with frontend):**
   ```bash
   docker-compose up -d
   ```

## Deployment

### Production Deployment Options

- **Docker Compose**: Simple multi-container deployment
- **Cloud Platforms**: Ready for deployment on any container platform

See [Deployment Guide](docs/deployment.md) for detailed instructions.

### CI/CD

The project includes GitHub Actions workflows for:

- **Continuous Integration**: Testing, linting, and security scanning
- **Docker Image Building**: Multi-architecture images pushed to GitHub Container Registry
- **Automated Deployment**: Staging and production deployments to Kubernetes

## API Endpoints

| Endpoint | Method | Description |
|----------|---------|-------------|
| `/healthz` | GET | Health check |
| `/ws` | WebSocket | Real-time communication |
| `/api/notes` | GET/POST | Note operations |
| `/api/transcribe` | POST | Audio transcription |
| `/api/summarize` | POST | Text summarization |

## Configuration

The application uses environment variables for configuration:

```bash
PORT=8080                # Server port
LOG_LEVEL=info          # Logging level
# Add other environment variables as needed
```

## Development

### Running Tests

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...
```

### Code Quality

```bash
# Linting
golangci-lint run

# Security scanning
gosec ./...

# Format code
go fmt ./...
```

## Future Roadmap

### Phase 1: Shared Module (note-shared)
- Extract common packages to separate module
- Enable code sharing between services
- Implement semantic versioning

### Phase 2: Service Expansion
- Note management service
- User authentication service
- File storage service

### Phase 3: Advanced Features
- Event-driven architecture with pub/sub
- Distributed tracing and monitoring
- Multi-region deployment

See [Shared Module Plan](docs/shared-module-plan.md) for detailed roadmap.

## Package Migration

The `internal/util` package is being phased out in favor of `pkg/` packages:

- ✅ `pkg/response` - HTTP response utilities
- ✅ `pkg/timeutil` - Time formatting utilities
- 🔄 `internal/util` - Deprecated wrapper (backward compatibility)

**Migration Guide:**
```go
// Old (deprecated)
import "github.com/your-org/note-server/internal/util"

// New (recommended)
import "github.com/your-org/note-server/pkg/response"
import "github.com/your-org/note-server/pkg/timeutil"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following the existing patterns
4. Add tests for new functionality
5. Submit a pull request

### Guidelines

- Follow Go best practices and idioms
- Use `any` type instead of `interface{}`
- Keep `internal/` packages private and specific
- Put shared utilities in `pkg/` packages
- Write comprehensive tests
- Update documentation

## License

[Add your license information here]

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issues](../../issues)
- 💬 [Discussions](../../discussions)
