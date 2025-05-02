# OSP_backend

A RESTful API service for managing online surveys and responses, built with Go and MongoDB following hexagonal architecture and domain-driven design principles.

## Quick Start

### Prerequisites
- Go 1.24.2 or later
- MongoDB 8.0 or later

### Installation
```bash
git clone https://github.com/lccheungperry/OSP_backend.git
cd OSP_backend
go mod download
```

### Configuration
1. Copy `.env.example` to `.env`
2. Update environment variables as needed

### Running
```bash
go run cmd/main.go
```

## Documentation

- [System Overview](docs/overview.md) - Detailed system architecture, features, and domain organization
- [API Documentation](docs/apiDocumentation.md) - Complete API reference
- [Project Structure](docs/projectStructure.md) - Code organization and design decisions
