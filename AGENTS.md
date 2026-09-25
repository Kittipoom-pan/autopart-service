# Autopart-Service - AI Agent Guidelines

## Tech Stack

- **Language**: Go 1.24.0
- **Web Framework**: Fiber v2 (github.com/gofiber/fiber/v2)
- **Database**: MySQL with go-sql-driver/mysql
- **Query Builder**: sqlc (SQL code generator)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Configuration**: godotenv for environment variables
- **Logging**: zerolog (github.com/rs/zerolog)
- **Hot Reload**: Air (.air.toml)
- **Migration**: golang-migrate
- **Containerization**: Docker (MySQL + phpMyAdmin)

## Project Structure

```
autopart-service/
├── cmd/
│   └── main.go                    # Application entry point
├── config/
│   └── config.go                  # Configuration loading (Server, DB, JWT)
├── internal/
│   ├── auth/                      # Authentication utilities
│   │   ├── jwt.go                 # JWT token generation/validation
│   │   ├── password.go            # Password hashing
│   │   ├── role.go                # Role constants
│   │   └── user.go                # User context helpers
│   ├── common/                    # Common types and enums
│   │   ├── base_response.go       # Standard response structures
│   │   └── enum.go                # Common enums (HTTP status codes)
│   ├── helper/                    # HTTP response helpers
│   │   └── http_response.go       # RespondSuccess/RespondError functions
│   ├── infrastructure/
│   │   └── database/
│   │       ├── mysql.go           # MySQL connection setup
│   │       ├── schema.sql         # Database schema
│   │       ├── query/             # SQL queries for sqlc
│   │       ├── sqlc.yaml         # sqlc configuration
│   │       └── sqlc/              # Generated code (DO NOT EDIT)
│   ├── logger/                    # Logging setup
│   │   └── logger.go              # zerolog initialization
│   ├── middleware/                # HTTP middleware
│   │   ├── jwt_middleware.go      # JWT authentication
│   │   └── timeout_middleware.go  # Request timeout handling
│   ├── module/                    # Business modules
│   │   ├── admin/                 # Admin module
│   │   │   ├── controller/        # HTTP handlers
│   │   │   ├── entity/            # Domain entities (requests/responses)
│   │   │   ├── repository/        # Data access layer
│   │   │   ├── usecase/           # Business logic layer
│   │   │   │   └── validation/    # Request validation
│   │   │   └── route.go           # Route setup
│   │   ├── customer/              # Customer module (same structure)
│   │   └── part/                  # Part module (same structure)
│   └── server/                    # Server setup
│       ├── server.go              # Fiber app initialization
│       └── route_handler.go      # Route mapping
├── migrations/                    # Database migrations
├── pkg/                           # Reusable packages
│   ├── error/                     # Custom error types
│   │   └── custom_error.go
│   └── utils/                     # Utility functions
│       ├── convert_util.go
│       ├── env_util.go
│       └── time_util.go
├── .air.toml                      # Air hot reload configuration
├── docker-compose.yml             # MySQL and phpMyAdmin setup
├── go.mod                         # Go module dependencies
└── go.sum                         # Dependency checksums
```

## How to Run the Project

### Prerequisites
- Go 1.24.0 or higher
- MySQL database
- .env file with required environment variables

### Environment Variables (.env)
Required environment variables:
```env
SERVER_HOST=localhost
SERVER_PORT=3000
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=autopart_db
MYSQL_USER=your_user
MYSQL_PASSWORD=your_password
JWT_SECRET_KEY=your_secret_key
JWT_EXPIRY=3600
MYSQL_ROOT_PASSWORD=root_password
PHPMYADMIN_PASSWORD=phpmyadmin_password
```

### Development (Hot Reload)
```bash
# Install air if not already installed
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Production Build
```bash
# Build the application
go build -o bin/main.exe ./cmd

# Run the binary
./bin/main.exe
```

### Docker Setup
```bash
# Start MySQL and phpMyAdmin
docker-compose up -d

# Stop containers
docker-compose down
```

### Database Migration
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq <name>

# Run migrations up
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" up

# Rollback migrations
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" down

# Check migration version
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" version

# Force migration version (use carefully for dirty state)
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" force <version>
```

### Generate SQL Code with sqlc
```bash
# Generate Go code from SQL queries
sqlc generate -f internal/infrastructure/database/sqlc.yaml
```

## Code Style and Architecture

### Clean Architecture Pattern
The project follows a 3-layer Clean Architecture:
1. **Controller Layer** (`internal/module/*/controller/`): HTTP request handling, response formatting
2. **Usecase Layer** (`internal/module/*/usecase/`): Business logic, validation
3. **Repository Layer** (`internal/module/*/repository/`): Data access, database operations

### Module Structure
Each business module (admin, customer, part) follows this structure:
```
module/
├── controller/          # HTTP handlers
├── entity/             # Domain entities and DTOs
│   ├── mapper.go       # Entity mapping functions
│   ├── *_request.go    # Request structures
│   └── *_response.go   # Response structures
├── repository/         # Data access interfaces and implementations
├── usecase/            # Business logic interfaces and implementations
│   └── validation/     # Request validation logic
└── route.go            # Route setup (public/private)
```

### Dependency Injection
- Use constructor functions for dependency injection
- Example: `NewAdminController(usecase usecase.AdminUsecase) *AdminController`
- Pass dependencies through function parameters, not global variables

### Context Propagation
- Always pass `context.Context` as the first parameter to usecase and repository methods
- Use `c.UserContext()` from Fiber to get request context
- Example: `func (u *adminUsecase) GetAdminByID(ctx context.Context, id int)`

### Error Handling
- Use custom error types from `pkg/error/custom_error.go`
- Available error types:
  - `APIError` - General API errors with status code and message
  - `NotFoundError` - Resource not found errors
  - `UnauthorizedError` - Authentication errors
  - `ForbiddenError` - Authorization errors
- Use helper functions: `RespondError(c, err)` and `RespondSuccess(c, status, data, message)`
- Always return structured error responses to clients

### Logging
- Use zerolog for structured logging
- Initialize logger with component name: `log.With().Str("component", "admin_controller").Logger()`
- Log levels: Debug, Info, Warn, Error
- Include relevant context (ids, usernames, etc.) in log entries

### Validation
- Validation logic belongs in `usecase/validation/` package
- Use custom error types for validation errors
- Return validation errors before business logic execution
- Example: `validation.ValidateAdminRequest(req, isUpdate)`

### Entity Mapping
- Use mapper functions in `entity/mapper.go` to convert between database models and API entities
- Keep mapping logic centralized within each module

### Route Organization
- Separate public routes (no auth) from private routes (JWT protected)
- Use middleware consistently:
  - `TimeoutMiddleware` for all routes
  - `JWTMiddleware` for private routes
  - `RequireRoles` for role-based access control
- Route grouping: `/v1/customer`, `/v1/admin`, `/v1/part`

## Iron Rules (NEVER DO)

### Security
- ❌ **NEVER commit secrets** (API keys, passwords, JWT secrets) to the repository
- ❌ **NEVER log sensitive data** (passwords, tokens, personal information)
- ❌ **NEVER bypass authentication/authorization** for development shortcuts
- ❌ **NEVER use weak password hashing** (always use auth.HashPassword)

### Code Generation
- ❌ **NEVER edit sqlc generated code** in `internal/infrastructure/database/sqlc/`
- ❌ **NEVER write raw SQL queries** in business logic - use sqlc queries
- ❌ **NEVER manually modify generated files** - regenerate with sqlc

### Architecture
- ❌ **NEVER skip layers** (e.g., controller calling repository directly)
- ❌ **NEVER put business logic in controllers** - belongs in usecase
- ❌ **NEVER put validation in controllers** - belongs in usecase/validation
- ❌ **NEVER use global variables** for dependencies - use dependency injection
- ❌ **NEVER ignore context propagation** - always pass context through layers

### Error Handling
- ❌ **NEVER return generic errors** without proper error types
- ❌ **NEVER expose internal error details** to clients
- ❌ **NEVER ignore errors** - always handle or log them
- ❌ **NEVER use panic for error handling** - return errors instead

### Database
- ❌ **NEVER run migrations without backup** in production
- ❌ **NEVER use force migration** without understanding the consequences
- ❌ **NEVER modify migration files** after they've been applied
- ❌ **NEVER commit .env files** to the repository

### Code Style
- ❌ **NEVER add unnecessary comments** - code should be self-documenting
- ❌ **NEVER copy-paste code** - extract to reusable functions
- ❌ **NEVER ignore existing patterns** - follow the established architecture
- ❌ **NEVER add dependencies** without checking if they're already in use

## Development Workflow

1. **Add new feature/module**:
   - Create module structure following existing patterns
   - Add SQL queries to `internal/infrastructure/database/query/`
   - Run `sqlc generate` to generate code
   - Implement repository layer
   - Implement usecase layer with validation
   - Implement controller layer
   - Add routes in `route.go`
   - Register routes in `internal/server/route_handler.go`

2. **Modify database schema**:
   - Create new migration file
   - Update `internal/infrastructure/database/schema.sql`
   - Add/update queries in `query/` directory
   - Run `sqlc generate`
   - Update entity structures if needed
   - Test migration up/down

3. **Debug issues**:
   - Check logs with zerolog structured output
   - Verify database connection and schema
   - Ensure context is properly propagated
   - Check error handling chain

## Testing and Verification

Before committing changes:
- Run `go build ./cmd` to ensure code compiles
- Run `sqlc generate` after modifying SQL queries
- Test migrations up/down
- Verify hot reload works with `air`
- Check for any linting issues
- Ensure no secrets are in the code

## Additional Notes

- The project uses MySQL with sqlc for type-safe SQL queries
- JWT authentication is required for private routes
- Role-based access control: `super_admin`, `staff`, `customer`
- All private routes have a 3-second timeout
- The project supports both admin and customer user types
- phpMyAdmin is available at http://localhost:8080 when using docker-compose
