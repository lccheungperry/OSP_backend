# Project Structure

## Directory Layout
```
OSP_backend/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── domain/                 # Business logic and models
│   │   └── survey_platform/    # Main domain
│   │       ├── survey/         # Survey subdomain
│   │       │   ├── model/      # Survey data models
│   │       │   │   └── types.go      # Core domain types
│   │       │   ├── command/    # Survey commands
│   │       │   │   ├── types.go      # Command types
│   │       │   │   └── handler.go    # Command handlers
│   │       │   ├── query/      # Survey queries
│   │       │   │   ├── types.go      # Query types
│   │       │   │   └── handler.go    # Query handlers
│   │       │   ├── repository/ # Survey repository
│   │       │   │   └── interface.go  # Repository interface
│   │       │   └── service/    # Survey service coordination
│   │       ├── question/       # Question subdomain
│   │       │   ├── model/      # Question data models
│   │       │   │   └── types.go      # Core domain types
│   │       │   ├── command/    # Question commands
│   │       │   │   ├── types.go      # Command types
│   │       │   │   └── handler.go    # Command handlers
│   │       │   ├── query/      # Question queries
│   │       │   │   ├── types.go      # Query types
│   │       │   │   └── handler.go    # Query handlers
│   │       │   ├── repository/ # Question repository
│   │       │   │   └── interface.go  # Repository interface
│   │       │   └── service/    # Question service coordination
│   │       └── response/       # Response subdomain
│   │           ├── model/      # Response data models
│   │           │   └── types.go      # Core domain types
│   │           ├── command/    # Response commands
│   │           │   ├── types.go      # Command types
│   │           │   └── handler.go    # Command handlers
│   │           ├── query/      # Response queries
│   │           │   ├── types.go      # Query types
│   │           │   └── handler.go    # Query handlers
│   │           ├── repository/ # Response repository
│   │           │   └── interface.go  # Repository interface
│   │           └── service/    # Response service coordination
│   ├── infrastructure/         # External services integration
│   │   ├── mongodb/           # MongoDB implementation
│   │   │   ├── connection.go  # Database connection
│   │   │   └── repository/    # MongoDB repository implementations
│   │   └── middleware/        # Cross-cutting concerns
│   │       ├── logging.go     # Logging middleware
│   │       ├── validation.go  # Validation middleware
│   │       └── tracing.go     # Tracing middleware
│   └── interfaces/            # API and repository implementations
│       ├── http/              # HTTP handlers
│       │   ├── handlers/      # Request handlers
│       │   └── middleware/    # HTTP middleware
│       └── repository/        # Repository implementations
├── pkg/                       # Public packages
│   ├── utils/                 # Utility functions
│   └── config/                # Configuration management
├── docs/                      # Documentation
├── tests/                     # Test files
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── README.md                  # Project documentation
```

## Component Responsibilities

### Domain Layer (`internal/domain/survey_platform/`)
The main domain representing the survey platform, containing three subdomains:

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

### Model Organization
Each subdomain's model directory contains:

1. **types.go**
  - Core domain types and value objects
  - Validation rules and constraints
  - Domain-specific errors

### Command Organization
Each subdomain's command directory contains:

1. **types.go**
  - Command types for mutating operations
  - Command validation rules
  - Command-specific errors

2. **handler.go**
  - Command handler implementations
  - Command processing logic
  - Error handling

### Query Organization
Each subdomain's query directory contains:

1. **types.go**
  - Query types for read operations
  - Query validation rules
  - Query-specific errors

2. **handler.go**
  - Query handler implementations
  - Query processing logic
  - Error handling

### Repository Organization
Each subdomain's repository directory contains:

1. **interface.go**
  - Repository interface definition
  - Data access contracts
  - Error definitions

### Service Organization
Each subdomain's service directory contains:

1. **service.go**
  - Command and query coordination
  - Handler registration
  - Cross-cutting concerns

## Key Design Decisions

1. **Command Query Responsibility Segregation (CQRS)**
  - Clear separation between commands and queries
  - Explicit handler registration
  - Easy to add cross-cutting concerns

2. **Domain-Driven Design**
  - Clear bounded contexts
  - Ubiquitous language
  - Aggregate roots

3. **Hexagonal Architecture**
  - Clean separation of concerns
  - Domain layer at the core

4. **Repository Pattern**
  - Abstract data access
  - Consistent interface
  - Easy to switch implementations
  - Clear separation of concerns