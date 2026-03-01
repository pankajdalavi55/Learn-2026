# Part 10: Project Structure & Architecture

## 📦 Building Maintainable Goa Applications

> **Goal**: Master clean architecture patterns, proper package organization, and dependency injection to build scalable, testable Goa microservices.

---

## 📋 Table of Contents

1. [Clean Architecture Overview](#1-clean-architecture-overview)
2. [Goa Project Structure](#2-goa-project-structure)
3. [Separating Transport from Business Logic](#3-separating-transport-from-business-logic)
4. [Internal Packages](#4-internal-packages)
5. [Dependency Injection Patterns](#5-dependency-injection-patterns)
6. [Real-World Example: E-Commerce Service](#6-real-world-example-e-commerce-service)
7. [Best Practices & Anti-Patterns](#7-best-practices--anti-patterns)

---

## 1. Clean Architecture Overview

### 1.1 What is Clean Architecture?

Clean Architecture (by Robert C. Martin) separates concerns into layers, making code:
- **Testable** - Business logic independent of frameworks
- **Maintainable** - Changes in one layer don't affect others
- **Flexible** - Easy to swap implementations

```
┌─────────────────────────────────────────────────────────────────┐
│                     CLEAN ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│    ┌─────────────────────────────────────────────────────┐     │
│    │              External Layer (Frameworks)            │     │
│    │   ┌─────────────────────────────────────────────┐   │     │
│    │   │         Interface Adapters                  │   │     │
│    │   │   ┌─────────────────────────────────────┐   │   │     │
│    │   │   │        Application Layer            │   │   │     │
│    │   │   │   ┌─────────────────────────────┐   │   │   │     │
│    │   │   │   │     Domain/Entities         │   │   │   │     │
│    │   │   │   │      (Core Business)        │   │   │   │     │
│    │   │   │   └─────────────────────────────┘   │   │   │     │
│    │   │   │                                     │   │   │     │
│    │   │   │   Use Cases / Application Logic    │   │   │     │
│    │   │   └─────────────────────────────────────┘   │   │     │
│    │   │                                             │   │     │
│    │   │   Controllers, Gateways, Presenters        │   │     │
│    │   └─────────────────────────────────────────────┘   │     │
│    │                                                     │     │
│    │   Goa, HTTP, gRPC, Database, External APIs         │     │
│    └─────────────────────────────────────────────────────┘     │
│                                                                 │
│              Dependencies point INWARD only →                   │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Clean Architecture Layers

| Layer | Purpose | Goa Mapping |
|-------|---------|-------------|
| **Entities** | Core business objects & rules | Domain models |
| **Use Cases** | Application-specific business rules | Service interface implementations |
| **Interface Adapters** | Convert data between layers | Goa endpoints, controllers |
| **Frameworks** | External tools & delivery | Goa transport, HTTP/gRPC |

### 1.3 The Dependency Rule

```
┌───────────────────────────────────────────────────────────────┐
│                    DEPENDENCY DIRECTION                       │
├───────────────────────────────────────────────────────────────┤
│                                                               │
│   Outer Layers          →          Inner Layers               │
│                                                               │
│   ┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐│
│   │Framework│ ──► │Adapter  │ ──► │Use Case │ ──► │ Entity  ││
│   │  (Goa)  │     │(Endpoint│     │(Service)│     │(Domain) ││
│   └─────────┘     └─────────┘     └─────────┘     └─────────┘│
│                                                               │
│   Dependencies ALWAYS point inward (toward business logic)    │
│   Inner layers know NOTHING about outer layers                │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

**Key Rules:**
1. Inner layers cannot import outer layers
2. Data crosses boundaries as simple structs
3. Interfaces are defined by inner layers, implemented by outer layers

---

## 2. Goa Project Structure

### 2.1 Standard Goa Project Layout

```
┌─────────────────────────────────────────────────────────────────┐
│                    GOA PROJECT STRUCTURE                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   my-service/                                                   │
│   ├── cmd/                        # Application entry points    │
│   │   └── my-service/                                          │
│   │       └── main.go             # Main application           │
│   │                                                            │
│   ├── design/                     # Goa DSL definitions        │
│   │   ├── design.go               # API design                 │
│   │   ├── types.go                # Type definitions           │
│   │   └── security.go             # Security schemes           │
│   │                                                            │
│   ├── gen/                        # Generated code (DO NOT EDIT)│
│   │   ├── myservice/              # Service interfaces         │
│   │   ├── http/                   # HTTP transport             │
│   │   └── grpc/                   # gRPC transport             │
│   │                                                            │
│   ├── internal/                   # Private application code   │
│   │   ├── domain/                 # Domain models & logic      │
│   │   ├── service/                # Service implementations    │
│   │   ├── repository/             # Data access layer          │
│   │   └── infrastructure/         # External services          │
│   │                                                            │
│   ├── pkg/                        # Public shared packages     │
│   │                                                            │
│   ├── go.mod                                                   │
│   └── go.sum                                                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Recommended Directory Structure

```
my-service/
├── cmd/
│   └── my-service/
│       └── main.go                 # Entry point, wiring
│
├── design/
│   ├── api.go                      # API definition
│   ├── types.go                    # Request/Response types
│   ├── security.go                 # Auth schemes
│   └── errors.go                   # Error definitions
│
├── gen/                            # 🔒 GENERATED - DO NOT EDIT
│   ├── my_service/
│   │   ├── service.go              # Service interface
│   │   ├── endpoints.go            # Endpoint definitions
│   │   └── client.go               # Service client
│   ├── http/
│   │   └── my_service/
│   │       ├── server/
│   │       │   ├── encode_decode.go
│   │       │   ├── paths.go
│   │       │   ├── server.go
│   │       │   └── types.go
│   │       └── client/
│   └── grpc/
│       └── my_service/
│
├── internal/
│   ├── domain/                     # 💎 Core Business Logic
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   ├── order.go
│   │   │   └── product.go
│   │   ├── valueobject/
│   │   │   ├── email.go
│   │   │   ├── money.go
│   │   │   └── address.go
│   │   ├── repository/             # Repository interfaces
│   │   │   ├── user_repository.go
│   │   │   └── order_repository.go
│   │   └── service/                # Domain services
│   │       └── pricing_service.go
│   │
│   ├── application/                # 📋 Use Cases
│   │   ├── usecase/
│   │   │   ├── create_user.go
│   │   │   ├── place_order.go
│   │   │   └── process_payment.go
│   │   └── dto/                    # Data Transfer Objects
│   │       ├── user_dto.go
│   │       └── order_dto.go
│   │
│   ├── service/                    # 🔌 Goa Service Implementation
│   │   ├── user_service.go
│   │   ├── order_service.go
│   │   └── mapper/                 # Type mappers
│   │       ├── user_mapper.go
│   │       └── order_mapper.go
│   │
│   ├── repository/                 # 💾 Data Access Implementation
│   │   ├── postgres/
│   │   │   ├── user_repository.go
│   │   │   └── order_repository.go
│   │   └── memory/
│   │       └── user_repository.go
│   │
│   ├── infrastructure/             # 🌐 External Services
│   │   ├── email/
│   │   │   └── smtp_client.go
│   │   ├── payment/
│   │   │   └── stripe_client.go
│   │   └── cache/
│   │       └── redis_client.go
│   │
│   └── config/                     # ⚙️ Configuration
│       └── config.go
│
├── pkg/                            # 📦 Public Packages
│   ├── logger/
│   │   └── logger.go
│   └── errors/
│       └── errors.go
│
├── migrations/                     # Database migrations
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
│
├── scripts/                        # Build & deploy scripts
│   └── generate.sh
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── go.sum
```

### 2.3 Layer Mapping in Goa

```
┌─────────────────────────────────────────────────────────────────┐
│                 LAYER MAPPING IN GOA PROJECT                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Clean Architecture         Goa Project Directory             │
│   ─────────────────         ────────────────────────           │
│                                                                 │
│   ┌─────────────────┐       ┌─────────────────────┐            │
│   │   Frameworks &  │  ═══► │ gen/                │            │
│   │   Drivers       │       │ cmd/                │            │
│   └─────────────────┘       └─────────────────────┘            │
│           │                          │                          │
│           ▼                          ▼                          │
│   ┌─────────────────┐       ┌─────────────────────┐            │
│   │   Interface     │  ═══► │ internal/service/   │            │
│   │   Adapters      │       │ internal/repository/│            │
│   └─────────────────┘       └─────────────────────┘            │
│           │                          │                          │
│           ▼                          ▼                          │
│   ┌─────────────────┐       ┌─────────────────────┐            │
│   │   Application   │  ═══► │ internal/application│            │
│   │   (Use Cases)   │       │ internal/usecase/   │            │
│   └─────────────────┘       └─────────────────────┘            │
│           │                          │                          │
│           ▼                          ▼                          │
│   ┌─────────────────┐       ┌─────────────────────┐            │
│   │   Entities      │  ═══► │ internal/domain/    │            │
│   │   (Domain)      │       │   entity/           │            │
│   └─────────────────┘       └─────────────────────┘            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. Separating Transport from Business Logic

### 3.1 The Problem: Mixed Concerns

❌ **Bad: Business logic in Goa service implementation**

```go
// ❌ DON'T DO THIS - Business logic mixed with transport concerns
type userServiceImpl struct {
    db *sql.DB
}

func (s *userServiceImpl) CreateUser(ctx context.Context, p *userservice.CreateUserPayload) (*userservice.User, error) {
    // ❌ Direct database access
    // ❌ Business rules mixed with transport
    // ❌ Hard to test
    // ❌ Hard to reuse
    
    // Validation (should be in domain)
    if len(p.Password) < 8 {
        return nil, userservice.MakeBadRequest(errors.New("password too short"))
    }
    
    // Hashing (should be in domain/service)
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(p.Password), 10)
    
    // Database (should be in repository)
    result, err := s.db.ExecContext(ctx, 
        "INSERT INTO users (email, password) VALUES ($1, $2)",
        p.Email, hashedPassword)
    
    // Error mapping (okay here, but cluttered)
    if err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            return nil, userservice.MakeConflict(errors.New("email exists"))
        }
        return nil, err
    }
    
    id, _ := result.LastInsertId()
    
    return &userservice.User{
        ID:    int(id),
        Email: p.Email,
    }, nil
}
```

### 3.2 The Solution: Layered Architecture

✅ **Good: Clean separation of concerns**

```
┌─────────────────────────────────────────────────────────────────┐
│              TRANSPORT / BUSINESS LOGIC SEPARATION              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   HTTP Request                                                  │
│        │                                                        │
│        ▼                                                        │
│   ┌──────────────────────────────────────────────────────┐     │
│   │              gen/http/ (Generated)                    │     │
│   │   • Decode HTTP request to Goa payload               │     │
│   │   • Encode response to HTTP response                 │     │
│   └──────────────────────────────────────────────────────┘     │
│        │                                                        │
│        ▼                                                        │
│   ┌──────────────────────────────────────────────────────┐     │
│   │           internal/service/ (Adapter Layer)          │     │
│   │   • Implement Goa service interface                  │     │
│   │   • Map Goa types ↔ Domain types                     │     │
│   │   • Delegate to use cases                            │     │
│   │   • Map errors to Goa errors                         │     │
│   └──────────────────────────────────────────────────────┘     │
│        │                                                        │
│        ▼                                                        │
│   ┌──────────────────────────────────────────────────────┐     │
│   │         internal/application/ (Use Case Layer)       │     │
│   │   • Orchestrate business operations                  │     │
│   │   • Transaction management                           │     │
│   │   • Call domain services & repositories              │     │
│   └──────────────────────────────────────────────────────┘     │
│        │                                                        │
│        ▼                                                        │
│   ┌──────────────────────────────────────────────────────┐     │
│   │           internal/domain/ (Domain Layer)            │     │
│   │   • Entities with business rules                     │     │
│   │   • Value objects                                    │     │
│   │   • Repository interfaces (ports)                    │     │
│   │   • Domain services                                  │     │
│   └──────────────────────────────────────────────────────┘     │
│        │                                                        │
│        ▼                                                        │
│   ┌──────────────────────────────────────────────────────┐     │
│   │       internal/repository/ (Infrastructure)          │     │
│   │   • Repository implementations                       │     │
│   │   • Database access                                  │     │
│   │   • External service clients                         │     │
│   └──────────────────────────────────────────────────────┘     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.3 Implementation Example

#### Step 1: Domain Layer (internal/domain/)

```go
// internal/domain/entity/user.go
package entity

import (
    "errors"
    "time"
)

// User is a domain entity with business rules
type User struct {
    ID           string
    Email        Email      // Value object
    PasswordHash string
    FirstName    string
    LastName     string
    Status       UserStatus
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// UserStatus is an enum value object
type UserStatus string

const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
    UserStatusBanned   UserStatus = "banned"
)

// NewUser creates a new user with validation
func NewUser(email Email, firstName, lastName string, hashedPassword string) (*User, error) {
    if firstName == "" || lastName == "" {
        return nil, errors.New("first and last name required")
    }
    
    return &User{
        Email:        email,
        FirstName:    firstName,
        LastName:     lastName,
        PasswordHash: hashedPassword,
        Status:       UserStatusActive,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }, nil
}

// FullName returns the user's full name (business logic)
func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// CanLogin checks if user can authenticate (business rule)
func (u *User) CanLogin() bool {
    return u.Status == UserStatusActive
}

// Ban bans the user (business operation)
func (u *User) Ban() error {
    if u.Status == UserStatusBanned {
        return errors.New("user already banned")
    }
    u.Status = UserStatusBanned
    u.UpdatedAt = time.Now()
    return nil
}
```

```go
// internal/domain/valueobject/email.go
package valueobject

import (
    "errors"
    "regexp"
    "strings"
)

// Email is a value object with validation
type Email struct {
    value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a validated email
func NewEmail(email string) (Email, error) {
    normalized := strings.ToLower(strings.TrimSpace(email))
    
    if !emailRegex.MatchString(normalized) {
        return Email{}, errors.New("invalid email format")
    }
    
    return Email{value: normalized}, nil
}

// String returns the email string
func (e Email) String() string {
    return e.value
}

// Domain returns the email domain
func (e Email) Domain() string {
    parts := strings.Split(e.value, "@")
    if len(parts) == 2 {
        return parts[1]
    }
    return ""
}

// Equals compares two emails
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

```go
// internal/domain/repository/user_repository.go
package repository

import (
    "context"
    
    "myservice/internal/domain/entity"
    "myservice/internal/domain/valueobject"
)

// UserRepository defines the port for user persistence
// This is an interface defined in domain, implemented in infrastructure
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id string) (*entity.User, error)
    GetByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*entity.User, int, error)
}
```

#### Step 2: Application Layer (internal/application/)

```go
// internal/application/usecase/create_user.go
package usecase

import (
    "context"
    "errors"
    
    "myservice/internal/domain/entity"
    "myservice/internal/domain/repository"
    "myservice/internal/domain/valueobject"
)

// CreateUserInput is the use case input DTO
type CreateUserInput struct {
    Email     string
    Password  string
    FirstName string
    LastName  string
}

// CreateUserOutput is the use case output DTO
type CreateUserOutput struct {
    ID        string
    Email     string
    FirstName string
    LastName  string
}

// CreateUserUseCase handles user creation
type CreateUserUseCase struct {
    userRepo       repository.UserRepository
    passwordHasher PasswordHasher
}

// PasswordHasher is a port for password hashing
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(password, hash string) error
}

// NewCreateUserUseCase creates the use case
func NewCreateUserUseCase(
    userRepo repository.UserRepository,
    hasher PasswordHasher,
) *CreateUserUseCase {
    return &CreateUserUseCase{
        userRepo:       userRepo,
        passwordHasher: hasher,
    }
}

// Execute runs the use case
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // 1. Validate email (creates value object)
    email, err := valueobject.NewEmail(input.Email)
    if err != nil {
        return nil, &ValidationError{Field: "email", Err: err}
    }
    
    // 2. Check if email exists
    existing, err := uc.userRepo.GetByEmail(ctx, email)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        return nil, err
    }
    if existing != nil {
        return nil, &ConflictError{Resource: "user", Field: "email"}
    }
    
    // 3. Hash password
    hashedPassword, err := uc.passwordHasher.Hash(input.Password)
    if err != nil {
        return nil, err
    }
    
    // 4. Create domain entity
    user, err := entity.NewUser(email, input.FirstName, input.LastName, hashedPassword)
    if err != nil {
        return nil, &ValidationError{Err: err}
    }
    
    // 5. Persist
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 6. Return output DTO
    return &CreateUserOutput{
        ID:        user.ID,
        Email:     user.Email.String(),
        FirstName: user.FirstName,
        LastName:  user.LastName,
    }, nil
}
```

```go
// internal/application/usecase/errors.go
package usecase

import "fmt"

// ValidationError represents a validation failure
type ValidationError struct {
    Field string
    Err   error
}

func (e *ValidationError) Error() string {
    if e.Field != "" {
        return fmt.Sprintf("validation error on field %s: %v", e.Field, e.Err)
    }
    return fmt.Sprintf("validation error: %v", e.Err)
}

// ConflictError represents a conflict (duplicate)
type ConflictError struct {
    Resource string
    Field    string
}

func (e *ConflictError) Error() string {
    return fmt.Sprintf("%s with this %s already exists", e.Resource, e.Field)
}

// NotFoundError represents a missing resource
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}
```

#### Step 3: Service Layer - Goa Adapter (internal/service/)

```go
// internal/service/user_service.go
package service

import (
    "context"
    "errors"
    
    "myservice/gen/userservice"
    "myservice/internal/application/usecase"
)

// UserService implements the Goa userservice.Service interface
type UserService struct {
    createUser *usecase.CreateUserUseCase
    getUser    *usecase.GetUserUseCase
    listUsers  *usecase.ListUsersUseCase
}

// NewUserService creates a new user service
func NewUserService(
    createUser *usecase.CreateUserUseCase,
    getUser *usecase.GetUserUseCase,
    listUsers *usecase.ListUsersUseCase,
) *UserService {
    return &UserService{
        createUser: createUser,
        getUser:    getUser,
        listUsers:  listUsers,
    }
}

// CreateUser implements userservice.Service.CreateUser
func (s *UserService) CreateUser(ctx context.Context, p *userservice.CreateUserPayload) (*userservice.User, error) {
    // 1. Map Goa payload to use case input
    input := usecase.CreateUserInput{
        Email:     p.Email,
        Password:  p.Password,
        FirstName: p.FirstName,
        LastName:  p.LastName,
    }
    
    // 2. Execute use case
    output, err := s.createUser.Execute(ctx, input)
    if err != nil {
        // 3. Map domain errors to Goa errors
        return nil, mapToGoaError(err)
    }
    
    // 4. Map output to Goa response
    return &userservice.User{
        ID:        output.ID,
        Email:     output.Email,
        FirstName: output.FirstName,
        LastName:  output.LastName,
    }, nil
}

// GetUser implements userservice.Service.GetUser
func (s *UserService) GetUser(ctx context.Context, p *userservice.GetUserPayload) (*userservice.User, error) {
    output, err := s.getUser.Execute(ctx, p.ID)
    if err != nil {
        return nil, mapToGoaError(err)
    }
    
    return mapUserToGoa(output), nil
}

// ListUsers implements userservice.Service.ListUsers
func (s *UserService) ListUsers(ctx context.Context, p *userservice.ListUsersPayload) (*userservice.UserList, error) {
    limit := 20
    if p.Limit != nil {
        limit = *p.Limit
    }
    
    offset := 0
    if p.Offset != nil {
        offset = *p.Offset
    }
    
    output, err := s.listUsers.Execute(ctx, limit, offset)
    if err != nil {
        return nil, mapToGoaError(err)
    }
    
    users := make([]*userservice.User, len(output.Users))
    for i, u := range output.Users {
        users[i] = mapUserToGoa(u)
    }
    
    return &userservice.UserList{
        Users: users,
        Total: output.Total,
    }, nil
}

// mapToGoaError maps domain/application errors to Goa errors
func mapToGoaError(err error) error {
    var validationErr *usecase.ValidationError
    if errors.As(err, &validationErr) {
        return userservice.MakeBadRequest(err)
    }
    
    var conflictErr *usecase.ConflictError
    if errors.As(err, &conflictErr) {
        return userservice.MakeConflict(err)
    }
    
    var notFoundErr *usecase.NotFoundError
    if errors.As(err, &notFoundErr) {
        return userservice.MakeNotFound(err)
    }
    
    // Unknown error - return as internal error
    return userservice.MakeInternalError(err)
}

// mapUserToGoa maps use case output to Goa type
func mapUserToGoa(u *usecase.UserOutput) *userservice.User {
    return &userservice.User{
        ID:        u.ID,
        Email:     u.Email,
        FirstName: u.FirstName,
        LastName:  u.LastName,
    }
}
```

#### Step 4: Repository Implementation (internal/repository/)

```go
// internal/repository/postgres/user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "errors"
    "time"
    
    "github.com/google/uuid"
    
    "myservice/internal/domain/entity"
    "myservice/internal/domain/repository"
    "myservice/internal/domain/valueobject"
)

type userRepository struct {
    db *sql.DB
}

// NewUserRepository creates a PostgreSQL user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    user.ID = uuid.New().String()
    
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO users (id, email, password_hash, first_name, last_name, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, user.ID, user.Email.String(), user.PasswordHash, user.FirstName, user.LastName, 
       user.Status, user.CreatedAt, user.UpdatedAt)
    
    return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
    var user entity.User
    var emailStr string
    
    err := r.db.QueryRowContext(ctx, `
        SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
        FROM users WHERE id = $1
    `, id).Scan(
        &user.ID, &emailStr, &user.PasswordHash, 
        &user.FirstName, &user.LastName, &user.Status,
        &user.CreatedAt, &user.UpdatedAt,
    )
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, repository.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    
    // Reconstruct value object
    email, _ := valueobject.NewEmail(emailStr)
    user.Email = email
    
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
    var user entity.User
    var emailStr string
    
    err := r.db.QueryRowContext(ctx, `
        SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
        FROM users WHERE email = $1
    `, email.String()).Scan(
        &user.ID, &emailStr, &user.PasswordHash,
        &user.FirstName, &user.LastName, &user.Status,
        &user.CreatedAt, &user.UpdatedAt,
    )
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, repository.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    
    user.Email = email
    return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
    user.UpdatedAt = time.Now()
    
    result, err := r.db.ExecContext(ctx, `
        UPDATE users SET 
            email = $2, first_name = $3, last_name = $4, 
            status = $5, updated_at = $6
        WHERE id = $1
    `, user.ID, user.Email.String(), user.FirstName, user.LastName,
       user.Status, user.UpdatedAt)
    
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return repository.ErrNotFound
    }
    
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
    result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return repository.ErrNotFound
    }
    
    return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, int, error) {
    // Get total count
    var total int
    err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    // Get users
    rows, err := r.db.QueryContext(ctx, `
        SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    var users []*entity.User
    for rows.Next() {
        var user entity.User
        var emailStr string
        
        if err := rows.Scan(
            &user.ID, &emailStr, &user.PasswordHash,
            &user.FirstName, &user.LastName, &user.Status,
            &user.CreatedAt, &user.UpdatedAt,
        ); err != nil {
            return nil, 0, err
        }
        
        email, _ := valueobject.NewEmail(emailStr)
        user.Email = email
        users = append(users, &user)
    }
    
    return users, total, nil
}
```

---

## 4. Internal Packages

### 4.1 Understanding `internal/` Directory

Go's `internal/` directory provides **compiler-enforced encapsulation**:

```
┌─────────────────────────────────────────────────────────────────┐
│                    INTERNAL PACKAGE RULES                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   my-service/                                                   │
│   ├── internal/                                                 │
│   │   └── domain/           ◄── Can be imported by:            │
│   │       └── entity/           • my-service/cmd/...           │
│   │                             • my-service/internal/...      │
│   │                             • my-service/pkg/...           │
│   │                                                            │
│   │                         ✗ CANNOT be imported by:           │
│   │                             • other-service/...            │
│   │                             • any external package         │
│   │                                                            │
│   └── pkg/                  ◄── Can be imported by ANYONE      │
│       └── logger/                                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Internal Package Organization

```go
// internal/domain/entity/user.go
// ✅ Private to this module - contains core business entities
package entity

// internal/domain/repository/interfaces.go
// ✅ Private interfaces - implemented by infrastructure layer
package repository

// internal/application/usecase/create_user.go
// ✅ Private use cases - business operations
package usecase

// internal/service/user_service.go
// ✅ Private - Goa service implementation
package service

// internal/repository/postgres/user.go
// ✅ Private - database implementation
package postgres

// internal/infrastructure/email/smtp.go
// ✅ Private - external service implementations
package email
```

### 4.3 What Goes in `internal/` vs `pkg/`

| Directory | Purpose | Examples |
|-----------|---------|----------|
| `internal/domain/` | Core business logic | Entities, Value Objects, Domain Services |
| `internal/application/` | Use cases | CreateUser, ProcessOrder |
| `internal/service/` | Goa implementations | Service adapters |
| `internal/repository/` | Data access | PostgreSQL, MongoDB impl |
| `internal/infrastructure/` | External services | Email, Payment clients |
| `internal/config/` | Configuration | App config struct |
| `pkg/logger/` | Reusable logging | Structured logger |
| `pkg/errors/` | Shared error types | Common error types |
| `pkg/middleware/` | Reusable middleware | Auth, CORS, etc. |

### 4.4 Nested Internal Packages

```
my-service/
├── internal/
│   ├── domain/
│   │   ├── internal/          ◄── Only domain/ can import
│   │   │   └── validation/    # Domain-specific validators
│   │   ├── entity/
│   │   └── repository/
│   │
│   └── repository/
│       ├── internal/          ◄── Only repository/ can import
│       │   └── queries/       # SQL query builders
│       └── postgres/
```

---

## 5. Dependency Injection Patterns

### 5.1 Why Dependency Injection?

```
┌─────────────────────────────────────────────────────────────────┐
│                    DEPENDENCY INJECTION BENEFITS                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Without DI (hard-coded dependencies):                         │
│   ┌──────────────────┐                                         │
│   │   UserService    │──────────► PostgresUserRepo (concrete)  │
│   │   (hard to test) │                                         │
│   └──────────────────┘                                         │
│                                                                 │
│   With DI (injected dependencies):                              │
│   ┌──────────────────┐                                         │
│   │   UserService    │──────────► UserRepository (interface)   │
│   │   (easy to test) │                   ▲                     │
│   └──────────────────┘                   │                     │
│                              ┌───────────┴───────────┐         │
│                              │                       │         │
│                        PostgresRepo             MockRepo       │
│                        (production)              (test)        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Constructor Injection (Recommended)

```go
// internal/service/user_service.go
package service

import (
    "myservice/internal/application/usecase"
    "myservice/internal/domain/repository"
)

// UserService implements Goa service interface
type UserService struct {
    userRepo       repository.UserRepository
    passwordHasher usecase.PasswordHasher
    logger         Logger
}

// Logger interface for dependency injection
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, err error, fields ...interface{})
}

// NewUserService creates a UserService with injected dependencies
func NewUserService(
    userRepo repository.UserRepository,
    passwordHasher usecase.PasswordHasher,
    logger Logger,
) *UserService {
    return &UserService{
        userRepo:       userRepo,
        passwordHasher: passwordHasher,
        logger:         logger,
    }
}
```

### 5.3 Manual Wiring (No Framework)

```go
// cmd/my-service/main.go
package main

import (
    "context"
    "database/sql"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    
    _ "github.com/lib/pq"
    
    "myservice/gen/userservice"
    goahttp "myservice/gen/http/userservice/server"
    "myservice/internal/application/usecase"
    "myservice/internal/config"
    "myservice/internal/infrastructure/hasher"
    "myservice/internal/repository/postgres"
    "myservice/internal/service"
    "myservice/pkg/logger"
)

func main() {
    ctx := context.Background()
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize logger
    log := logger.New(cfg.LogLevel)
    
    // Initialize database
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database", err)
    }
    defer db.Close()
    
    // Build dependency graph manually
    
    // Layer 1: Infrastructure
    userRepo := postgres.NewUserRepository(db)
    passwordHasher := hasher.NewBcryptHasher()
    
    // Layer 2: Use Cases
    createUserUC := usecase.NewCreateUserUseCase(userRepo, passwordHasher)
    getUserUC := usecase.NewGetUserUseCase(userRepo)
    listUsersUC := usecase.NewListUsersUseCase(userRepo)
    
    // Layer 3: Service (Goa implementation)
    userSvc := service.NewUserService(createUserUC, getUserUC, listUsersUC)
    
    // Layer 4: Endpoints
    endpoints := userservice.NewEndpoints(userSvc)
    
    // Layer 5: HTTP Transport
    mux := goahttp.New(nil, nil)
    server := goahttp.NewServer(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
    goahttp.Mount(mux, server)
    
    // Start server
    httpServer := &http.Server{
        Addr:    cfg.HTTPAddr,
        Handler: mux,
    }
    
    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        
        log.Info("Shutting down server...")
        httpServer.Shutdown(ctx)
    }()
    
    log.Info("Starting server", "addr", cfg.HTTPAddr)
    if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatal("Server failed", err)
    }
}
```

### 5.4 Using Wire (Google's DI Framework)

Wire is a compile-time dependency injection tool that generates code.

```go
// cmd/my-service/wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "database/sql"
    
    "github.com/google/wire"
    
    "myservice/gen/userservice"
    "myservice/internal/application/usecase"
    "myservice/internal/config"
    "myservice/internal/infrastructure/hasher"
    "myservice/internal/repository/postgres"
    "myservice/internal/service"
    "myservice/pkg/logger"
)

// ProviderSet groups providers by layer
var infrastructureSet = wire.NewSet(
    postgres.NewUserRepository,
    hasher.NewBcryptHasher,
    wire.Bind(new(usecase.PasswordHasher), new(*hasher.BcryptHasher)),
)

var useCaseSet = wire.NewSet(
    usecase.NewCreateUserUseCase,
    usecase.NewGetUserUseCase,
    usecase.NewListUsersUseCase,
)

var serviceSet = wire.NewSet(
    service.NewUserService,
    // Bind to Goa interface
    wire.Bind(new(userservice.Service), new(*service.UserService)),
)

// InitializeApp creates the application with all dependencies
func InitializeApp(cfg *config.Config, db *sql.DB, log *logger.Logger) (*App, error) {
    wire.Build(
        infrastructureSet,
        useCaseSet,
        serviceSet,
        NewApp,
    )
    return nil, nil
}
```

```go
// cmd/my-service/app.go
package main

import (
    "context"
    "net/http"
    
    "myservice/gen/userservice"
    goahttp "myservice/gen/http/userservice/server"
    "myservice/internal/config"
    "myservice/pkg/logger"
)

// App holds the application state
type App struct {
    cfg       *config.Config
    log       *logger.Logger
    endpoints *userservice.Endpoints
}

// NewApp creates a new application
func NewApp(
    cfg *config.Config,
    log *logger.Logger,
    svc userservice.Service,
) *App {
    return &App{
        cfg:       cfg,
        log:       log,
        endpoints: userservice.NewEndpoints(svc),
    }
}

// Run starts the HTTP server
func (a *App) Run(ctx context.Context) error {
    mux := goahttp.New(nil, nil)
    server := goahttp.NewServer(a.endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
    goahttp.Mount(mux, server)
    
    httpServer := &http.Server{
        Addr:    a.cfg.HTTPAddr,
        Handler: mux,
    }
    
    return httpServer.ListenAndServe()
}
```

Run Wire to generate the wiring code:

```bash
# Install wire
go install github.com/google/wire/cmd/wire@latest

# Generate wire_gen.go
cd cmd/my-service
wire
```

### 5.5 Using fx (Uber's DI Framework)

fx provides runtime dependency injection with lifecycle management.

```go
// cmd/my-service/main.go
package main

import (
    "context"
    "database/sql"
    "net/http"
    
    "go.uber.org/fx"
    
    "myservice/gen/userservice"
    goahttp "myservice/gen/http/userservice/server"
    "myservice/internal/application/usecase"
    "myservice/internal/config"
    "myservice/internal/infrastructure/hasher"
    "myservice/internal/repository/postgres"
    "myservice/internal/service"
    "myservice/pkg/logger"
)

func main() {
    fx.New(
        // Configuration
        fx.Provide(config.Load),
        
        // Infrastructure
        fx.Provide(
            NewDatabase,
            logger.New,
            postgres.NewUserRepository,
            hasher.NewBcryptHasher,
        ),
        
        // Bind interfaces
        fx.Provide(
            fx.Annotate(
                hasher.NewBcryptHasher,
                fx.As(new(usecase.PasswordHasher)),
            ),
        ),
        
        // Use Cases
        fx.Provide(
            usecase.NewCreateUserUseCase,
            usecase.NewGetUserUseCase,
            usecase.NewListUsersUseCase,
        ),
        
        // Service
        fx.Provide(
            service.NewUserService,
            fx.Annotate(
                service.NewUserService,
                fx.As(new(userservice.Service)),
            ),
        ),
        
        // Transport
        fx.Provide(
            userservice.NewEndpoints,
            NewHTTPServer,
        ),
        
        // Start server
        fx.Invoke(RegisterHooks),
    ).Run()
}

// NewDatabase creates a database connection
func NewDatabase(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }
    
    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return db.Close()
        },
    })
    
    return db, nil
}

// NewHTTPServer creates the HTTP server
func NewHTTPServer(endpoints *userservice.Endpoints) *http.Server {
    mux := goahttp.New(nil, nil)
    server := goahttp.NewServer(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
    goahttp.Mount(mux, server)
    
    return &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
}

// RegisterHooks starts and stops the HTTP server
func RegisterHooks(lc fx.Lifecycle, srv *http.Server, log *logger.Logger) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go func() {
                log.Info("Starting HTTP server", "addr", srv.Addr)
                if err := srv.ListenAndServe(); err != http.ErrServerClosed {
                    log.Error("HTTP server failed", err)
                }
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            log.Info("Stopping HTTP server")
            return srv.Shutdown(ctx)
        },
    })
}
```

### 5.6 DI Comparison

| Approach | Pros | Cons |
|----------|------|------|
| **Manual** | Simple, no deps, explicit | Verbose, error-prone |
| **Wire** | Compile-time safe, fast | Build step, learning curve |
| **fx** | Lifecycle mgmt, flexible | Runtime overhead, harder to debug |

**Recommendation**: Start with manual wiring. Move to Wire for larger projects.

---

## 6. Real-World Example: E-Commerce Service

### 6.1 Complete Project Structure

```
ecommerce-service/
├── cmd/
│   └── ecommerce/
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
│
├── design/
│   ├── api.go
│   ├── types.go
│   ├── products.go
│   ├── orders.go
│   ├── users.go
│   └── errors.go
│
├── gen/                        # 🔒 GENERATED
│   ├── ecommerce/
│   ├── http/
│   └── grpc/
│
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   ├── order_item.go
│   │   │   ├── user.go
│   │   │   └── cart.go
│   │   ├── valueobject/
│   │   │   ├── money.go
│   │   │   ├── email.go
│   │   │   ├── address.go
│   │   │   └── quantity.go
│   │   ├── repository/
│   │   │   ├── product_repository.go
│   │   │   ├── order_repository.go
│   │   │   ├── user_repository.go
│   │   │   └── errors.go
│   │   └── service/
│   │       ├── pricing_service.go
│   │       └── inventory_service.go
│   │
│   ├── application/
│   │   ├── usecase/
│   │   │   ├── product/
│   │   │   │   ├── create_product.go
│   │   │   │   ├── get_product.go
│   │   │   │   └── list_products.go
│   │   │   ├── order/
│   │   │   │   ├── create_order.go
│   │   │   │   ├── get_order.go
│   │   │   │   └── cancel_order.go
│   │   │   └── user/
│   │   │       ├── register_user.go
│   │   │       └── authenticate_user.go
│   │   ├── dto/
│   │   │   ├── product_dto.go
│   │   │   └── order_dto.go
│   │   └── port/
│   │       ├── payment_gateway.go
│   │       └── notification_service.go
│   │
│   ├── service/
│   │   ├── product_service.go
│   │   ├── order_service.go
│   │   ├── user_service.go
│   │   └── mapper/
│   │       ├── product_mapper.go
│   │       ├── order_mapper.go
│   │       └── user_mapper.go
│   │
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── product_repository.go
│   │   │   ├── order_repository.go
│   │   │   ├── user_repository.go
│   │   │   └── connection.go
│   │   └── memory/
│   │       └── product_repository.go
│   │
│   ├── infrastructure/
│   │   ├── payment/
│   │   │   ├── stripe_gateway.go
│   │   │   └── mock_gateway.go
│   │   ├── notification/
│   │   │   ├── email_service.go
│   │   │   └── sms_service.go
│   │   └── cache/
│   │       └── redis_cache.go
│   │
│   └── config/
│       ├── config.go
│       └── env.go
│
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── errors/
│   │   └── errors.go
│   └── middleware/
│       ├── logging.go
│       ├── recovery.go
│       └── auth.go
│
├── migrations/
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_products.up.sql
│   ├── 002_create_products.down.sql
│   ├── 003_create_orders.up.sql
│   └── 003_create_orders.down.sql
│
├── scripts/
│   ├── generate.sh
│   └── migrate.sh
│
├── api/                        # API documentation
│   └── openapi.yaml
│
├── deployments/
│   ├── docker/
│   │   └── Dockerfile
│   └── kubernetes/
│       ├── deployment.yaml
│       └── service.yaml
│
├── Makefile
├── docker-compose.yml
├── go.mod
└── go.sum
```

### 6.2 Domain Entity Example

```go
// internal/domain/entity/order.go
package entity

import (
    "errors"
    "time"
    
    "github.com/google/uuid"
    
    "ecommerce/internal/domain/valueobject"
)

// OrderStatus represents the order state
type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCancelled OrderStatus = "cancelled"
)

// Order is the aggregate root for order domain
type Order struct {
    ID              string
    UserID          string
    Items           []OrderItem
    ShippingAddress valueobject.Address
    BillingAddress  valueobject.Address
    Status          OrderStatus
    Subtotal        valueobject.Money
    Tax             valueobject.Money
    ShippingCost    valueobject.Money
    Total           valueobject.Money
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// OrderItem represents an item in an order
type OrderItem struct {
    ProductID   string
    ProductName string
    Quantity    valueobject.Quantity
    UnitPrice   valueobject.Money
    Total       valueobject.Money
}

// NewOrder creates a new order
func NewOrder(userID string, shippingAddr, billingAddr valueobject.Address) *Order {
    return &Order{
        ID:              uuid.New().String(),
        UserID:          userID,
        Items:           []OrderItem{},
        ShippingAddress: shippingAddr,
        BillingAddress:  billingAddr,
        Status:          OrderStatusPending,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }
}

// AddItem adds an item to the order (business logic)
func (o *Order) AddItem(productID, productName string, qty valueobject.Quantity, unitPrice valueobject.Money) error {
    if o.Status != OrderStatusPending {
        return errors.New("cannot modify non-pending order")
    }
    
    // Check if item already exists
    for i, item := range o.Items {
        if item.ProductID == productID {
            newQty, err := item.Quantity.Add(qty)
            if err != nil {
                return err
            }
            o.Items[i].Quantity = newQty
            o.Items[i].Total = unitPrice.Multiply(newQty.Value())
            o.recalculateTotals()
            return nil
        }
    }
    
    // Add new item
    o.Items = append(o.Items, OrderItem{
        ProductID:   productID,
        ProductName: productName,
        Quantity:    qty,
        UnitPrice:   unitPrice,
        Total:       unitPrice.Multiply(qty.Value()),
    })
    
    o.recalculateTotals()
    return nil
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(productID string) error {
    if o.Status != OrderStatusPending {
        return errors.New("cannot modify non-pending order")
    }
    
    for i, item := range o.Items {
        if item.ProductID == productID {
            o.Items = append(o.Items[:i], o.Items[i+1:]...)
            o.recalculateTotals()
            return nil
        }
    }
    
    return errors.New("item not found in order")
}

// Confirm confirms the order (state transition)
func (o *Order) Confirm() error {
    if o.Status != OrderStatusPending {
        return errors.New("only pending orders can be confirmed")
    }
    if len(o.Items) == 0 {
        return errors.New("cannot confirm empty order")
    }
    
    o.Status = OrderStatusConfirmed
    o.UpdatedAt = time.Now()
    return nil
}

// Ship marks the order as shipped
func (o *Order) Ship() error {
    if o.Status != OrderStatusConfirmed {
        return errors.New("only confirmed orders can be shipped")
    }
    
    o.Status = OrderStatusShipped
    o.UpdatedAt = time.Now()
    return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
    if o.Status == OrderStatusDelivered {
        return errors.New("cannot cancel delivered order")
    }
    if o.Status == OrderStatusCancelled {
        return errors.New("order already cancelled")
    }
    
    o.Status = OrderStatusCancelled
    o.UpdatedAt = time.Now()
    return nil
}

// CanBeCancelled checks if order can be cancelled
func (o *Order) CanBeCancelled() bool {
    return o.Status != OrderStatusDelivered && o.Status != OrderStatusCancelled
}

// recalculateTotals updates all monetary totals
func (o *Order) recalculateTotals() {
    subtotal := valueobject.ZeroMoney(o.Subtotal.Currency())
    
    for _, item := range o.Items {
        subtotal = subtotal.Add(item.Total)
    }
    
    o.Subtotal = subtotal
    o.Tax = subtotal.Multiply(0.1) // 10% tax
    o.Total = subtotal.Add(o.Tax).Add(o.ShippingCost)
    o.UpdatedAt = time.Now()
}
```

### 6.3 Use Case Example

```go
// internal/application/usecase/order/create_order.go
package order

import (
    "context"
    "errors"
    
    "ecommerce/internal/application/port"
    "ecommerce/internal/domain/entity"
    "ecommerce/internal/domain/repository"
    "ecommerce/internal/domain/valueobject"
)

// CreateOrderInput is the use case input
type CreateOrderInput struct {
    UserID          string
    Items           []OrderItemInput
    ShippingAddress AddressInput
    BillingAddress  AddressInput
    PaymentMethodID string
}

type OrderItemInput struct {
    ProductID string
    Quantity  int
}

type AddressInput struct {
    Street     string
    City       string
    State      string
    Country    string
    PostalCode string
}

// CreateOrderOutput is the use case output
type CreateOrderOutput struct {
    OrderID string
    Status  string
    Total   float64
}

// CreateOrderUseCase handles order creation
type CreateOrderUseCase struct {
    orderRepo      repository.OrderRepository
    productRepo    repository.ProductRepository
    userRepo       repository.UserRepository
    paymentGateway port.PaymentGateway
    txManager      repository.TransactionManager
}

// NewCreateOrderUseCase creates the use case
func NewCreateOrderUseCase(
    orderRepo repository.OrderRepository,
    productRepo repository.ProductRepository,
    userRepo repository.UserRepository,
    paymentGateway port.PaymentGateway,
    txManager repository.TransactionManager,
) *CreateOrderUseCase {
    return &CreateOrderUseCase{
        orderRepo:      orderRepo,
        productRepo:    productRepo,
        userRepo:       userRepo,
        paymentGateway: paymentGateway,
        txManager:      txManager,
    }
}

// Execute runs the use case
func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
    // 1. Validate user exists
    user, err := uc.userRepo.GetByID(ctx, input.UserID)
    if err != nil {
        return nil, errors.New("user not found")
    }
    if !user.CanPlaceOrders() {
        return nil, errors.New("user cannot place orders")
    }
    
    // 2. Create value objects
    shippingAddr, err := valueobject.NewAddress(
        input.ShippingAddress.Street,
        input.ShippingAddress.City,
        input.ShippingAddress.State,
        input.ShippingAddress.Country,
        input.ShippingAddress.PostalCode,
    )
    if err != nil {
        return nil, err
    }
    
    billingAddr, err := valueobject.NewAddress(
        input.BillingAddress.Street,
        input.BillingAddress.City,
        input.BillingAddress.State,
        input.BillingAddress.Country,
        input.BillingAddress.PostalCode,
    )
    if err != nil {
        return nil, err
    }
    
    // 3. Create order entity
    order := entity.NewOrder(input.UserID, shippingAddr, billingAddr)
    
    // 4. Add items (within transaction)
    err = uc.txManager.Execute(ctx, func(ctx context.Context) error {
        for _, item := range input.Items {
            // Get product
            product, err := uc.productRepo.GetByID(ctx, item.ProductID)
            if err != nil {
                return err
            }
            
            // Check inventory
            if !product.HasStock(item.Quantity) {
                return errors.New("insufficient stock for " + product.Name)
            }
            
            // Create quantity value object
            qty, err := valueobject.NewQuantity(item.Quantity)
            if err != nil {
                return err
            }
            
            // Add to order
            if err := order.AddItem(product.ID, product.Name, qty, product.Price); err != nil {
                return err
            }
            
            // Reserve inventory
            if err := product.ReserveStock(item.Quantity); err != nil {
                return err
            }
            
            // Save product
            if err := uc.productRepo.Update(ctx, product); err != nil {
                return err
            }
        }
        
        // 5. Process payment
        paymentResult, err := uc.paymentGateway.Charge(ctx, port.ChargeInput{
            Amount:          order.Total.Amount(),
            Currency:        order.Total.Currency(),
            PaymentMethodID: input.PaymentMethodID,
            Description:     "Order " + order.ID,
        })
        if err != nil {
            return err
        }
        
        if !paymentResult.Success {
            return errors.New("payment failed: " + paymentResult.ErrorMessage)
        }
        
        // 6. Confirm order
        if err := order.Confirm(); err != nil {
            return err
        }
        
        // 7. Save order
        return uc.orderRepo.Create(ctx, order)
    })
    
    if err != nil {
        return nil, err
    }
    
    return &CreateOrderOutput{
        OrderID: order.ID,
        Status:  string(order.Status),
        Total:   order.Total.Amount(),
    }, nil
}
```

### 6.4 Goa Service Implementation

```go
// internal/service/order_service.go
package service

import (
    "context"
    
    "ecommerce/gen/ecommerce"
    orderuc "ecommerce/internal/application/usecase/order"
    "ecommerce/internal/service/mapper"
)

// OrderService implements Goa ecommerce.Service for orders
type OrderService struct {
    createOrder *orderuc.CreateOrderUseCase
    getOrder    *orderuc.GetOrderUseCase
    cancelOrder *orderuc.CancelOrderUseCase
}

// NewOrderService creates the order service
func NewOrderService(
    createOrder *orderuc.CreateOrderUseCase,
    getOrder *orderuc.GetOrderUseCase,
    cancelOrder *orderuc.CancelOrderUseCase,
) *OrderService {
    return &OrderService{
        createOrder: createOrder,
        getOrder:    getOrder,
        cancelOrder: cancelOrder,
    }
}

// CreateOrder implements ecommerce.Service.CreateOrder
func (s *OrderService) CreateOrder(ctx context.Context, p *ecommerce.CreateOrderPayload) (*ecommerce.Order, error) {
    // Map Goa payload to use case input
    input := mapper.ToCreateOrderInput(p)
    
    // Execute use case
    output, err := s.createOrder.Execute(ctx, input)
    if err != nil {
        return nil, mapOrderError(err)
    }
    
    // Map output to Goa response
    return mapper.ToGoaOrder(output), nil
}

// GetOrder implements ecommerce.Service.GetOrder
func (s *OrderService) GetOrder(ctx context.Context, p *ecommerce.GetOrderPayload) (*ecommerce.Order, error) {
    output, err := s.getOrder.Execute(ctx, p.OrderID)
    if err != nil {
        return nil, mapOrderError(err)
    }
    
    return mapper.ToGoaOrderDetail(output), nil
}

// CancelOrder implements ecommerce.Service.CancelOrder
func (s *OrderService) CancelOrder(ctx context.Context, p *ecommerce.CancelOrderPayload) error {
    err := s.cancelOrder.Execute(ctx, p.OrderID)
    if err != nil {
        return mapOrderError(err)
    }
    
    return nil
}
```

---

## 7. Best Practices & Anti-Patterns

### 7.1 Best Practices

```
┌─────────────────────────────────────────────────────────────────┐
│                      BEST PRACTICES                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ✅ DO                                                          │
│  ────                                                           │
│  • Define interfaces in the consumer package                    │
│  • Keep domain layer free of framework imports                  │
│  • Use value objects for validation and type safety             │
│  • Map errors at service layer boundaries                       │
│  • Use constructor injection for dependencies                   │
│  • Write business logic in domain entities                      │
│  • Keep use cases focused on single operations                  │
│  • Use meaningful package names (not util, common, etc.)        │
│  • Document public APIs and complex business rules              │
│  • Use internal/ for implementation details                     │
│                                                                 │
│  ❌ DON'T                                                       │
│  ──────                                                         │
│  • Import Goa types in domain layer                             │
│  • Put business logic in Goa service implementations            │
│  • Use global variables for dependencies                        │
│  • Create circular dependencies between packages                │
│  • Edit generated code in gen/ directory                        │
│  • Mix transport concerns with business logic                   │
│  • Over-engineer for problems you don't have                    │
│  • Create deep package hierarchies                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 7.2 Common Anti-Patterns

#### ❌ Anti-Pattern 1: Anemic Domain Model

```go
// ❌ BAD: Entity is just a data container
type User struct {
    ID       string
    Email    string
    Password string
    Status   string
}

// Business logic scattered in services
type UserService struct {}

func (s *UserService) CanLogin(u *User) bool {
    return u.Status == "active"
}

func (s *UserService) Ban(u *User) {
    u.Status = "banned"
}
```

```go
// ✅ GOOD: Rich domain model
type User struct {
    id           string
    email        Email // Value object
    passwordHash string
    status       UserStatus
}

// Business logic in entity
func (u *User) CanLogin() bool {
    return u.status == UserStatusActive
}

func (u *User) Ban() error {
    if u.status == UserStatusBanned {
        return ErrAlreadyBanned
    }
    u.status = UserStatusBanned
    return nil
}
```

#### ❌ Anti-Pattern 2: Framework Coupling

```go
// ❌ BAD: Domain depends on Goa
package domain

import (
    "ecommerce/gen/ecommerce" // Framework dependency!
)

type OrderService struct {}

func (s *OrderService) CreateOrder(p *ecommerce.CreateOrderPayload) (*ecommerce.Order, error) {
    // Domain logic mixed with Goa types
}
```

```go
// ✅ GOOD: Domain is framework-agnostic
package domain

type OrderService struct {}

type CreateOrderInput struct {
    UserID string
    Items  []OrderItem
}

type CreateOrderOutput struct {
    OrderID string
    Total   Money
}

func (s *OrderService) CreateOrder(input CreateOrderInput) (*CreateOrderOutput, error) {
    // Pure domain logic - no framework types
}
```

#### ❌ Anti-Pattern 3: Package Cycles

```go
// ❌ BAD: Circular dependency
// package user imports package order
// package order imports package user

package user

import "myapp/order"

func (u *User) GetOrders() []*order.Order {
    // ...
}

package order

import "myapp/user"

func (o *Order) GetUser() *user.User {
    // ...
}
```

```go
// ✅ GOOD: Use interfaces to break cycles
package user

type Order interface {
    ID() string
    Total() float64
}

func (u *User) GetOrders() []Order {
    // Returns interface, not concrete type
}

package order

type User interface {
    ID() string
    CanPlaceOrders() bool
}

func (o *Order) ValidateUser(u User) error {
    // Accepts interface
}
```

### 7.3 Layer Dependency Rules

```
┌─────────────────────────────────────────────────────────────────┐
│                    DEPENDENCY RULES                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Layer              Can Import                Cannot Import    │
│   ─────              ──────────                ─────────────    │
│                                                                 │
│   domain/entity      • Standard library        • application/   │
│                      • domain/valueobject      • service/       │
│                      • domain/repository       • repository/    │
│                        (interfaces only)       • infrastructure/│
│                                                • gen/           │
│                                                • cmd/           │
│                                                                 │
│   domain/repository  • domain/entity           • application/   │
│   (interfaces)       • Standard library        • service/       │
│                                                • infrastructure/│
│                                                                 │
│   application/       • domain/*                • service/       │
│                      • application/port        • repository/    │
│                        (interfaces)            • infrastructure/│
│                      • Standard library        • gen/           │
│                                                                 │
│   service/           • gen/ (Goa types)        • repository/    │
│                      • application/*           • infrastructure/│
│                      • pkg/*                                    │
│                                                                 │
│   repository/        • domain/*                • application/   │
│   (implementations)  • DB drivers              • service/       │
│                      • Standard library                         │
│                                                                 │
│   infrastructure/    • domain/                 • service/       │
│                      • application/port        • gen/           │
│                      • External libraries                       │
│                                                                 │
│   cmd/               • Everything              -                │
│                      (wiring layer)                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📝 Quick Reference

### Project Structure Checklist

```
□ cmd/               - Entry points only
□ design/            - Goa DSL definitions
□ gen/               - Generated code (never edit)
□ internal/domain/   - Business entities & rules
□ internal/application/ - Use cases
□ internal/service/  - Goa service implementations
□ internal/repository/ - Data access implementations
□ internal/infrastructure/ - External services
□ internal/config/   - Configuration
□ pkg/               - Shared public packages
□ migrations/        - Database migrations
□ scripts/           - Build/deploy scripts
```

### Dependency Injection Guide

| Size | Recommended Approach |
|------|---------------------|
| Small (< 10 deps) | Manual wiring |
| Medium (10-50 deps) | Wire (compile-time) |
| Large (50+ deps) | fx (runtime) |

### Layer Responsibilities

| Layer | Responsibility | Knows About |
|-------|---------------|-------------|
| Domain | Business rules | Nothing external |
| Application | Use case orchestration | Domain |
| Service | Transport adaptation | Application, Goa |
| Repository | Data persistence | Domain, Database |
| Infrastructure | External services | Domain, External APIs |

---

## 🎯 Summary

Key takeaways for Goa project architecture:

1. **Separate Concerns**: Transport (Goa) → Service → Use Case → Domain
2. **Depend Inward**: Outer layers depend on inner, never reverse
3. **Use `internal/`**: Enforce encapsulation at compile time
4. **Inject Dependencies**: Use constructor injection for testability
5. **Rich Domain Model**: Put business logic in entities, not services
6. **Map at Boundaries**: Convert types between layers

---

**Next Up**: Part 11 - Production Readiness (Logging, Graceful Shutdown, Docker, Health Checks, Metrics, Tracing)
