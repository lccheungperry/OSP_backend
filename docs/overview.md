# OSP_backend

## Overview
OSP_backend is a RESTful API service for managing online surveys and responses. Built with Go and MongoDB, it follows hexagonal architecture and domain-driven design principles with a clear domain organization.

## Domain Organization
The system is organized around a main "Survey Platform" domain with three subdomains:

#### Survey Subdomain
- Manages survey creation and lifecycle
- Handles token generation for public access
- Manages survey-question assignments

#### Question Subdomain
- Manages reusable question templates
- Handles question format validation
- Supports multiple question types
- Maintains question specifications

#### Response Subdomain
- Manages survey response collection
- Handles response validation
- Maintains response history

## Features

### Survey Management
- CRUD operations for surveys
- Token generation for public survey access
- Question management with reusability

### Response Management
- Submit responses to surveys
- View survey responses
- Response validation

## Technology Stack
- **Language**: Go
- **Database**: MongoDB
- **Architecture**: Hexagonal Architecture
- **Design**: Domain-Driven Design
- **API**: RESTful

## Key Design Principles
1. Domain-Driven Design with clear subdomains
2. Clean separation of concerns through hexagonal architecture
3. Repository pattern for data access abstraction
4. Comprehensive error handling strategy
5. Question reusability through survey-question assignments

## Getting Started
1. Prerequisites
- Go 1.24.2 or later
- MongoDB 8.0 or later

2. Installation
```bash
git clone https://github.com/lccheungperry/OSP_backend.git
cd OSP_backend
go mod download
```

3. Configuration
- Copy `.env.example` to `.env`
- Update environment variables as needed

4. Running the Application
```bash
go run cmd/main.go
```

## API Documentation
See [API Documentation](apiDocumentation.md) for detailed endpoint specifications.

## Project Structure
See [Project Structure](projectStructure.md) for detailed component organization.
