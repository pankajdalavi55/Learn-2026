# Part 2: Goa Core Concepts (Very Important)

> **The Heart of Goa** - Master these concepts to effectively design and build microservices with Goa framework.

---

## Table of Contents

1. [Design-First Philosophy](#1-design-first-philosophy)
2. [Services](#2-services)
3. [Methods](#3-methods)
4. [Data Modeling](#4-data-modeling)
5. [Security Basics](#5-security-basics)
6. [Service Implementation](#6-service-implementation)
7. [Complete Working Example](#7-complete-working-example)

---

## 1. Design-First Philosophy

### Understanding Design-First: The Theory

**What is Design-First?**

Design-First (also called API-First or Contract-First) is an approach where you define your API contract before implementing business logic. You describe "what" your API does before "how" it does it.

**Design-First vs Code-First:**

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Code-First Approach                              │
│                                                                      │
│   Write Code ──▶ Generate Docs ──▶ Docs often out of sync!          │
│   (handlers)      (swagger/openapi)                                  │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                    Design-First Approach (Goa)                       │
│                                                                      │
│   Design DSL ──▶ Generate Code ──▶ Implement Logic ──▶ Always in sync│
│   (design/)      (goa gen)        (business logic)                   │
└─────────────────────────────────────────────────────────────────────┘
```

**Benefits of Design-First:**

| Benefit | Description |
|---------|-------------|
| **Contract Stability** | API contract is defined upfront, reducing breaking changes |
| **Parallel Development** | Frontend/Backend teams can work simultaneously |
| **Auto-Generated Code** | Less boilerplate, fewer bugs |
| **Documentation Sync** | OpenAPI spec always matches implementation |
| **Type Safety** | Generated types ensure consistency |
| **Validation Built-in** | Request validation generated from design |

**Why Goa Uses Design-First:**

Goa embraces design-first because:
1. APIs are your product's contract - they deserve careful design
2. Generated code eliminates hand-written boilerplate errors
3. Changes to the design automatically propagate everywhere
4. Multiple transports (HTTP, gRPC) from single design

---

### 🔹 How Goa Generates Code

**The Generation Pipeline:**

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│   design/*.go    │────▶│     goa gen      │────▶│  Generated Code  │
│                  │     │                  │     │                  │
│  - Services      │     │  - Parses DSL    │     │  - Types         │
│  - Methods       │     │  - Validates     │     │  - Endpoints     │
│  - Types         │     │  - Generates     │     │  - Transport     │
│  - Transports    │     │                  │     │  - OpenAPI       │
└──────────────────┘     └──────────────────┘     └──────────────────┘
         │                                                  │
         │              Your Code Never Touches             │
         └──────────────────────────────────────────────────┘
                              gen/ folder
```

**What Gets Generated:**

| Generated File | Purpose | Location |
|----------------|---------|----------|
| `types.go` | Request/Response types | `gen/<service>/` |
| `service.go` | Service interface | `gen/<service>/` |
| `endpoints.go` | Endpoint definitions | `gen/<service>/` |
| `client.go` | Service client | `gen/<service>/` |
| `server.go` | HTTP server handlers | `gen/http/<service>/server/` |
| `client.go` | HTTP client | `gen/http/<service>/client/` |
| `encode_decode.go` | HTTP encoding/decoding | `gen/http/<service>/` |
| `paths.go` | URL path constructors | `gen/http/<service>/` |
| `openapi.json` | OpenAPI 3.0 spec | `gen/http/` |
| `openapi.yaml` | OpenAPI 3.0 spec (YAML) | `gen/http/` |

**Generation is Idempotent:**

Running `goa gen` multiple times produces the same output. You can safely regenerate after design changes without losing consistency.

**Never Edit Generated Code:**

Generated files contain this header:
```go
// Code generated by goa v3.x.x, DO NOT EDIT.
```

If you need to customize behavior, use:
- Middleware for cross-cutting concerns
- Custom encoders/decoders for special serialization
- Interceptors for request/response modification

---

### 🔹 Project Structure

**Standard Goa Project Layout:**

```
myservice/
├── cmd/                          # Entry points
│   └── myservice/
│       ├── main.go               # Main function
│       └── http.go               # HTTP server setup (optional)
│
├── design/                       # 📐 API Design (you write this)
│   ├── design.go                 # Main design file
│   ├── types.go                  # Shared type definitions
│   └── api.go                    # API-level settings (optional)
│
├── gen/                          # 🤖 Generated code (never edit!)
│   ├── myservice/                # Service layer
│   │   ├── client.go
│   │   ├── endpoints.go
│   │   ├── service.go
│   │   └── types.go
│   │
│   └── http/                     # HTTP transport layer
│       ├── myservice/
│       │   ├── server/
│       │   │   ├── encode_decode.go
│       │   │   ├── paths.go
│       │   │   ├── server.go
│       │   │   └── types.go
│       │   └── client/
│       │       └── ...
│       ├── openapi.json
│       ├── openapi.yaml
│       └── openapi3.json
│
├── myservice.go                  # 🔧 Service implementation (you write this)
├── go.mod
└── go.sum
```

**Understanding the Layers:**

```
┌─────────────────────────────────────────────────────────────────┐
│                        Your Code                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  design/        │  Implementation  │  cmd/              │   │
│  │  (API Design)   │  (Business Logic)│  (Entry Point)     │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Generated Code (gen/)                        │
│  ┌──────────────────────┐  ┌──────────────────────────────┐    │
│  │  Service Layer       │  │  Transport Layer             │    │
│  │  - Interfaces        │  │  - HTTP handlers             │    │
│  │  - Types             │  │  - Encoding/Decoding         │    │
│  │  - Endpoints         │  │  - Client/Server             │    │
│  └──────────────────────┘  └──────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

**File Ownership:**

| File/Folder | Owner | Edit? |
|-------------|-------|-------|
| `design/` | You | ✅ Yes |
| `gen/` | Goa | ❌ Never |
| `cmd/` | You | ✅ Yes |
| `*_service.go` | You | ✅ Yes |
| `go.mod` | You | ✅ Yes |

---

### 🔹 Design Package Concept

**What is the Design Package?**

The design package is a special Go package that uses Goa's DSL (Domain Specific Language) to describe your API. It's located in the `design/` folder and is only used at generation time - it's not compiled into your final binary.

**The DSL is Just Go Code:**

Goa's DSL functions are regular Go functions. The design package is valid Go code that gets executed by `goa gen` to produce the API description.

```go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// This IS Go code - it uses Goa's DSL functions
var _ = API("myapi", func() {
    Title("My API")
    Description("A demonstration API")
    Version("1.0")
})
```

**API() DSL - Complete Configuration:**

The `API()` function is the top-level entry point for your entire API definition. It sets global configuration:

```go
var _ = API("ecommerce", func() {
    // Basic Info
    Title("E-Commerce Platform API")
    Description("RESTful API for the e-commerce platform")
    Version("2.0")
    
    // Terms of Service
    TermsOfService("https://example.com/terms")
    
    // Contact Information
    Contact(func() {
        Name("API Support")
        Email("api@example.com")
        URL("https://example.com/support")
    })
    
    // License
    License(func() {
        Name("Apache 2.0")
        URL("https://www.apache.org/licenses/LICENSE-2.0")
    })
    
    // Documentation
    Docs(func() {
        Description("Additional documentation")
        URL("https://example.com/docs")
    })
    
    // Server definitions
    Server("production", func() {
        Description("Production server")
        Host("production", func() {
            Description("Production host")
            URI("https://api.example.com")
        })
    })
    
    Server("development", func() {
        Description("Development server")
        Host("localhost", func() {
            URI("http://localhost:8080")
        })
    })
    
    // HTTP configuration (applies to all services)
    HTTP(func() {
        // Global path prefix
        Path("/api")
        
        // Response content type
        Consumes("application/json")
        Produces("application/json")
    })
})
```

**Server and Host Configuration:**

```go
// Multiple servers for different environments
Server("api", func() {
    Description("Main API server")
    
    // Multiple hosts per server
    Host("production", func() {
        Description("Production")
        URI("https://api.example.com")
        URI("grpcs://grpc.example.com:443")  // gRPC endpoint
    })
    
    Host("staging", func() {
        Description("Staging")
        URI("https://staging-api.example.com")
    })
    
    Host("development", func() {
        Description("Development")
        // Variables in host definition
        URI("http://localhost:{port}")
        Variable("port", String, "HTTP port", func() {
            Default("8080")
            Enum("8080", "8081", "8082")
        })
    })
})
```

**DSL Building Blocks:**

```
┌─────────────────────────────────────────────────────────────────┐
│                        Goa DSL Hierarchy                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  API()              ← Top-level API definition                   │
│    │                                                             │
│    ├── Service()    ← Group of related methods                   │
│    │     │                                                       │
│    │     ├── Method()   ← Single operation                       │
│    │     │     │                                                 │
│    │     │     ├── Payload()  ← Input                            │
│    │     │     ├── Result()   ← Output                           │
│    │     │     └── Error()    ← Error definitions                │
│    │     │                                                       │
│    │     └── HTTP()     ← HTTP transport mapping                 │
│    │           │                                                 │
│    │           ├── GET/POST/PUT/DELETE                           │
│    │           └── Response()                                    │
│    │                                                             │
│    └── Type()       ← Reusable type definitions                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**The Dot Import:**

Notice the dot import:
```go
import . "goa.design/goa/v3/dsl"
```

This allows using DSL functions without the package prefix, making designs more readable:
```go
// With dot import
Service("users", func() { ... })

// Without dot import (verbose)
dsl.Service("users", func() { ... })
```

**Design Package Rules:**

1. Package name must be `design` (by convention)
2. Files must import Goa DSL
3. Use `var _ = ...` to execute DSL at package init time
4. Design is executed, not compiled into binary
5. Keep designs organized (split into multiple files for large APIs)

---

### 🔹 Code Generation Lifecycle (goa gen)

**Installation:**

```bash
# Install Goa CLI
go install goa.design/goa/v3/cmd/goa@latest

# Verify installation
goa version
```

**Basic Generation Command:**

```bash
# Generate code from design
goa gen <module-path>/design

# Example
goa gen github.com/username/myservice/design
```

**The Generation Process:**

```
┌─────────────────────────────────────────────────────────────────┐
│                    goa gen Execution Flow                        │
└─────────────────────────────────────────────────────────────────┘

Step 1: Load Design
┌──────────────────────┐
│  Parse design/*.go   │
│  Execute DSL funcs   │
│  Build API model     │
└──────────┬───────────┘
           │
           ▼
Step 2: Validate
┌──────────────────────┐
│  Check for errors    │
│  Validate types      │
│  Verify consistency  │
└──────────┬───────────┘
           │
           ▼
Step 3: Generate
┌──────────────────────┐
│  Service layer       │
│  Transport layer     │
│  OpenAPI specs       │
└──────────┬───────────┘
           │
           ▼
Step 4: Write Files
┌──────────────────────┐
│  Write to gen/       │
│  Format with gofmt   │
│  Report results      │
└──────────────────────┘
```

**Generation Options:**

```bash
# Generate with specific output directory
goa gen <design-path> -o ./output

# Generate only example implementation
goa example <design-path>

# Common workflow
goa gen mymodule/design && goa example mymodule/design
```

**The `goa example` Command:**

This generates skeleton implementation files that you can customize:

```bash
goa example github.com/username/myservice/design
```

Generated example files:
- `cmd/<service>/main.go` - Main entry point
- `cmd/<service>/http.go` - HTTP server setup
- `<service>.go` - Service implementation skeleton

**Regeneration Workflow:**

```
Edit design/*.go
      │
      ▼
Run: goa gen <design-path>
      │
      ▼
gen/ folder updated
      │
      ▼
Your implementation automatically
uses new types/interfaces
```

**Common Generation Errors:**

| Error | Cause | Solution |
|-------|-------|----------|
| `undefined: Service` | Missing DSL import | Add `. "goa.design/goa/v3/dsl"` |
| `duplicate method` | Same method name twice | Rename one method |
| `invalid type` | Undefined type reference | Define the type first |
| `missing payload` | Method requires input | Add `Payload()` |

**Best Practices:**

1. **Add to Makefile:**
   ```makefile
   generate:
       goa gen github.com/username/myservice/design
   
   example:
       goa example github.com/username/myservice/design
   ```

2. **Git Ignore Strategy:**
   ```gitignore
   # Some teams ignore gen/ (regenerate in CI)
   # gen/
   
   # Others commit gen/ for visibility
   # (recommended for smaller teams)
   ```

3. **Pre-commit Hook:**
   ```bash
   #!/bin/sh
   goa gen ./design && git add gen/
   ```

---

## 2. Services

### Understanding Services: The Theory

**What is a Service in Goa?**

A Service is a logical grouping of related operations (Methods) that work on a common domain. Think of it as a module or component of your API.

**Service Granularity:**

```
Too Coarse (Monolithic)           Just Right                    Too Fine (Fragmented)
┌───────────────────────┐    ┌─────────┐ ┌─────────┐ ┌─────────┐    ┌─────┐ ┌─────┐
│                       │    │  Users  │ │ Orders  │ │Products │    │GetU │ │SetU │
│    EverythingService  │    │ Service │ │ Service │ │ Service │    │     │ │     │
│                       │    │         │ │         │ │         │    └─────┘ └─────┘
│  - CreateUser         │    │-Create  │ │-Create  │ │-List    │    ┌─────┐ ┌─────┐
│  - GetUser            │    │-Get     │ │-Get     │ │-Get     │    │DelU │ │LstU │
│  - CreateOrder        │    │-Update  │ │-List    │ │-Update  │    │     │ │     │
│  - GetOrder           │    │-Delete  │ │-Cancel  │ │         │    └─────┘ └─────┘
│  - ListProducts       │    │         │ │         │ │         │
│  - ...100 more        │    └─────────┘ └─────────┘ └─────────┘
└───────────────────────┘
      ❌ Hard to             ✅ Cohesive                  ❌ Too many
      maintain               and focused                  services
```

**Service Design Principles:**

| Principle | Description |
|-----------|-------------|
| **Single Responsibility** | Each service handles one domain |
| **Cohesion** | Related methods stay together |
| **Clear Boundaries** | Well-defined interface and scope |
| **Independent Deployment** | Can be deployed separately (microservices) |

---

### 🔹 Defining a Service

**Basic Service Definition:**

```go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("The users service manages user accounts")
    
    // Methods go here...
})
```

**Service with Full Configuration:**

```go
var _ = Service("users", func() {
    // Documentation
    Description("User management service for the platform")
    
    // Error definitions (service-wide)
    Error("not_found", ErrorResult, "User not found")
    Error("unauthorized", ErrorResult, "Authentication required")
    Error("validation_error", ErrorResult, "Invalid input data")
    
    // HTTP transport configuration
    HTTP(func() {
        // Base path for all methods in this service
        Path("/api/v1/users")
    })
    
    // gRPC transport configuration (if using gRPC)
    GRPC(func() {
        // gRPC-specific settings
    })
    
    // Methods defined here
    Method("list", func() { /* ... */ })
    Method("get", func() { /* ... */ })
    Method("create", func() { /* ... */ })
    Method("update", func() { /* ... */ })
    Method("delete", func() { /* ... */ })
})
```

**Service-Level Errors:**

Define errors at service level for reuse across methods:

```go
var _ = Service("orders", func() {
    Description("Order management service")
    
    // These errors can be used by any method in this service
    Error("not_found", ErrorResult, func() {
        Description("Order not found")
    })
    Error("already_shipped", ErrorResult, func() {
        Description("Cannot modify shipped order")
    })
    Error("insufficient_stock", ErrorResult, func() {
        Description("Not enough inventory")
    })
    
    Method("cancel", func() {
        // Can reference service-level errors
        Error("not_found")
        Error("already_shipped")
        // ... rest of method
    })
})
```

**Error Types in Detail:**

```go
// Standard error type
var ErrorResult = Type("ErrorResult", func() {
    Description("Standard error response")
    Attribute("code", String, "Error code", func() {
        Example("VALIDATION_ERROR")
    })
    Attribute("message", String, "Human-readable message", func() {
        Example("Invalid email format")
    })
    Attribute("field", String, "Field that caused the error", func() {
        Example("email")
    })
    Attribute("details", MapOf(String, Any), "Additional error details")
    Required("code", "message")
})

// Using errors in service
var _ = Service("users", func() {
    // Define all possible errors
    Error("not_found", ErrorResult, "Resource not found")
    Error("unauthorized", ErrorResult, "Authentication required")
    Error("forbidden", ErrorResult, "Insufficient permissions")
    Error("validation", ErrorResult, "Validation failed")
    Error("conflict", ErrorResult, "Resource already exists")
    Error("internal", ErrorResult, "Internal server error")
    
    HTTP(func() {
        Path("/users")
        // Map errors to HTTP status codes
        Response("not_found", StatusNotFound)
        Response("unauthorized", StatusUnauthorized)
        Response("forbidden", StatusForbidden)
        Response("validation", StatusBadRequest)
        Response("conflict", StatusConflict)
        Response("internal", StatusInternalServerError)
    })
    
    Method("create", func() {
        Payload(CreateUserPayload)
        Result(User)
        
        // List errors this method can return
        Error("validation")
        Error("conflict")
        Error("internal")
        
        HTTP(func() {
            POST("/")
            Response(StatusCreated)
            // Method-specific error mappings (optional, inherits from service)
        })
    })
})
```

**Fault vs Timeout vs Temporary Errors:**

```go
// Temporary error - client can retry
Error("service_unavailable", ErrorResult, func() {
    Description("Service temporarily unavailable")
    Temporary()  // Marks error as temporary/retryable
})

// Timeout error
Error("timeout", ErrorResult, func() {
    Description("Request timed out")
    Timeout()  // Marks error as timeout
})

// Fault - server-side issue
Error("internal", ErrorResult, func() {
    Description("Internal server error")
    Fault()  // Marks error as server fault
})
```

**Generated Service Interface:**

From the design, Goa generates an interface:

```go
// gen/users/service.go (generated)
type Service interface {
    // List returns all users
    List(context.Context, *ListPayload) (*ListResult, error)
    // Get returns a user by ID
    Get(context.Context, *GetPayload) (*User, error)
    // Create adds a new user
    Create(context.Context, *CreatePayload) (*User, error)
    // Update modifies an existing user
    Update(context.Context, *UpdatePayload) (*User, error)
    // Delete removes a user
    Delete(context.Context, *DeletePayload) error
}
```

You implement this interface in your service implementation file.

---

### 🔹 Methods Inside a Service

**Method Structure:**

```go
Method("methodName", func() {
    // Documentation
    Description("What this method does")
    
    // Input
    Payload(PayloadType)
    
    // Output
    Result(ResultType)
    
    // Errors this method can return
    Error("error_name")
    
    // HTTP mapping
    HTTP(func() {
        GET("/path/{id}")
        Response(StatusOK)
    })
})
```

**Complete Method Example:**

```go
var _ = Service("users", func() {
    Description("User management service")
    
    Method("get", func() {
        Description("Retrieve a user by their ID")
        
        // What the client sends
        Payload(func() {
            Attribute("id", String, "User ID", func() {
                Format(FormatUUID)
            })
            Required("id")
        })
        
        // What the server returns
        Result(User)
        
        // Possible errors
        Error("not_found")
        Error("unauthorized")
        
        // HTTP specifics
        HTTP(func() {
            GET("/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    Method("create", func() {
        Description("Create a new user account")
        
        Payload(func() {
            Attribute("email", String, func() {
                Format(FormatEmail)
            })
            Attribute("password", String, func() {
                MinLength(8)
                MaxLength(100)
            })
            Attribute("name", String, func() {
                MinLength(1)
                MaxLength(255)
            })
            Required("email", "password", "name")
        })
        
        Result(User)
        
        Error("validation_error")
        Error("email_taken")
        
        HTTP(func() {
            POST("/")
            Body(func() {
                Attribute("email")
                Attribute("password")
                Attribute("name")
            })
            Response(StatusCreated)
            Response("validation_error", StatusBadRequest)
            Response("email_taken", StatusConflict)
        })
    })
    
    Method("list", func() {
        Description("List users with pagination")
        
        Payload(func() {
            Attribute("page", Int, "Page number", func() {
                Minimum(1)
                Default(1)
            })
            Attribute("limit", Int, "Items per page", func() {
                Minimum(1)
                Maximum(100)
                Default(20)
            })
            Attribute("sort", String, "Sort field", func() {
                Enum("created_at", "updated_at", "name")
                Default("created_at")
            })
        })
        
        Result(func() {
            Attribute("users", ArrayOf(User))
            Attribute("total", Int, "Total number of users")
            Attribute("page", Int, "Current page")
            Attribute("pages", Int, "Total pages")
            Required("users", "total", "page", "pages")
        })
        
        HTTP(func() {
            GET("/")
            Param("page")
            Param("limit")
            Param("sort")
            Response(StatusOK)
        })
    })
})
```

---

### 🔹 Multiple Services in One Project

**When to Use Multiple Services:**

- Different domains (users, orders, products)
- Different access patterns (public API vs admin API)
- Different scaling requirements
- Team ownership boundaries

**Multi-Service Design:**

```go
// design/design.go - API definition
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = API("ecommerce", func() {
    Title("E-Commerce Platform API")
    Description("API for the e-commerce platform")
    Version("1.0")
    
    Server("api", func() {
        Description("API server")
        Host("localhost", func() {
            URI("http://localhost:8080")
        })
    })
})
```

```go
// design/users.go - Users service
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("User account management")
    
    HTTP(func() {
        Path("/api/v1/users")
    })
    
    Method("get", func() { /* ... */ })
    Method("create", func() { /* ... */ })
    Method("update", func() { /* ... */ })
    Method("delete", func() { /* ... */ })
})
```

```go
// design/orders.go - Orders service
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("orders", func() {
    Description("Order processing and management")
    
    HTTP(func() {
        Path("/api/v1/orders")
    })
    
    Method("create", func() { /* ... */ })
    Method("get", func() { /* ... */ })
    Method("list", func() { /* ... */ })
    Method("cancel", func() { /* ... */ })
})
```

```go
// design/products.go - Products service
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("products", func() {
    Description("Product catalog management")
    
    HTTP(func() {
        Path("/api/v1/products")
    })
    
    Method("list", func() { /* ... */ })
    Method("get", func() { /* ... */ })
    Method("search", func() { /* ... */ })
})
```

**Generated Structure for Multiple Services:**

```
gen/
├── users/
│   ├── client.go
│   ├── endpoints.go
│   ├── service.go
│   └── types.go
├── orders/
│   ├── client.go
│   ├── endpoints.go
│   ├── service.go
│   └── types.go
├── products/
│   ├── client.go
│   ├── endpoints.go
│   ├── service.go
│   └── types.go
└── http/
    ├── users/
    ├── orders/
    ├── products/
    ├── openapi.json
    └── openapi3.json
```

**Shared Types Across Services:**

```go
// design/types.go - Shared types
package design

import (
    . "goa.design/goa/v3/dsl"
)

// Shared across all services
var User = Type("User", func() {
    Attribute("id", String, func() {
        Format(FormatUUID)
    })
    Attribute("email", String)
    Attribute("name", String)
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Required("id", "email", "name", "created_at")
})

// Used by orders service
var OrderItem = Type("OrderItem", func() {
    Attribute("product_id", String)
    Attribute("quantity", Int)
    Attribute("price", Float64)
    Required("product_id", "quantity", "price")
})

// Shared error result
var ErrorResult = Type("ErrorResult", func() {
    Attribute("code", String, "Error code")
    Attribute("message", String, "Human-readable message")
    Attribute("details", MapOf(String, String), "Additional details")
    Required("code", "message")
})
```

---

## 3. Methods

### Understanding Methods: The Theory

**What is a Method?**

A Method represents a single operation in your API. It defines:
- **What input it accepts** (Payload)
- **What output it returns** (Result)
- **What can go wrong** (Errors)
- **How it maps to transport** (HTTP/gRPC)

**Method Anatomy:**

```
┌─────────────────────────────────────────────────────────────────┐
│                        Method Definition                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Method("create_user", func() {                                  │
│                                                                  │
│      ┌────────────────────────────────────────────────────┐     │
│      │ METADATA                                           │     │
│      │ Description("Creates a new user account")          │     │
│      │ Security(JWTAuth)                                  │     │
│      └────────────────────────────────────────────────────┘     │
│                                                                  │
│      ┌────────────────────────────────────────────────────┐     │
│      │ INPUT                                              │     │
│      │ Payload(CreateUserPayload)                         │     │
│      └────────────────────────────────────────────────────┘     │
│                                                                  │
│      ┌────────────────────────────────────────────────────┐     │
│      │ OUTPUT                                             │     │
│      │ Result(User)                                       │     │
│      └────────────────────────────────────────────────────┘     │
│                                                                  │
│      ┌────────────────────────────────────────────────────┐     │
│      │ ERRORS                                             │     │
│      │ Error("validation_error")                          │     │
│      │ Error("email_taken")                               │     │
│      └────────────────────────────────────────────────────┘     │
│                                                                  │
│      ┌────────────────────────────────────────────────────┐     │
│      │ TRANSPORT                                          │     │
│      │ HTTP(func() { POST("/") ... })                     │     │
│      │ GRPC(func() { ... })                               │     │
│      └────────────────────────────────────────────────────┘     │
│                                                                  │
│  })                                                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

### 🔹 Payload

**What is Payload?**

Payload defines the input data a method receives. It can come from various sources in HTTP:
- Request body
- URL path parameters
- Query parameters
- Headers

**Payload Definition Options:**

```go
// Option 1: Inline definition
Method("create", func() {
    Payload(func() {
        Attribute("name", String)
        Attribute("email", String)
        Required("name", "email")
    })
})

// Option 2: Reference a Type
Method("create", func() {
    Payload(CreateUserPayload)
})

// Option 3: Simple type for single value
Method("get", func() {
    Payload(String)  // Just a single string
})

// Option 4: No payload (no input needed)
Method("list_all", func() {
    // No Payload() - method takes no input
    Result(ArrayOf(User))
})
```

**Payload Sources in HTTP:**

```go
Method("update", func() {
    Payload(func() {
        // From URL path
        Attribute("id", String, "User ID from path")
        
        // From request body
        Attribute("name", String, "New name")
        Attribute("email", String, "New email")
        
        // From query parameter
        Attribute("notify", Boolean, "Send notification")
        
        // From header
        Attribute("authorization", String, "Auth token")
        
        Required("id")
    })
    
    HTTP(func() {
        PUT("/{id}")               // 'id' from path
        Param("notify")            // 'notify' from query
        Header("authorization:Authorization")  // From header
        Body(func() {              // Rest from body
            Attribute("name")
            Attribute("email")
        })
    })
})
```

**Payload Mapping Diagram:**

```
HTTP Request
┌─────────────────────────────────────────────────────────────┐
│ PUT /users/abc-123?notify=true HTTP/1.1                     │
│ Authorization: Bearer xyz789                                 │
│ Content-Type: application/json                              │
│                                                             │
│ {                                                           │
│   "name": "John Doe",                                       │
│   "email": "john@example.com"                               │
│ }                                                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Payload Struct                            │
│                                                              │
│  type UpdatePayload struct {                                 │
│      ID            string  // from path                      │
│      Name          string  // from body                      │
│      Email         string  // from body                      │
│      Notify        bool    // from query                     │
│      Authorization string  // from header                    │
│  }                                                           │
└─────────────────────────────────────────────────────────────┘
```

---

### 🔹 Result

**What is Result?**

Result defines what the method returns on success. Like Payload, it can be defined inline, reference a type, or be omitted.

**Result Definition Options:**

```go
// Option 1: Inline definition
Method("get", func() {
    Result(func() {
        Attribute("id", String)
        Attribute("name", String)
        Attribute("email", String)
        Required("id", "name", "email")
    })
})

// Option 2: Reference a Type
Method("get", func() {
    Result(User)
})

// Option 3: Collection result
Method("list", func() {
    Result(ArrayOf(User))
})

// Option 4: No result (void method)
Method("delete", func() {
    Payload(func() {
        Attribute("id", String)
        Required("id")
    })
    // No Result() - returns nothing on success
})

// Option 5: ResultType for streaming
Method("download", func() {
    Result(func() {
        ContentType("application/octet-stream")
    })
})
```

**Complex Result with Metadata:**

```go
Method("list", func() {
    Description("List users with pagination")
    
    Payload(ListPayload)
    
    Result(func() {
        Attribute("data", ArrayOf(User), "List of users")
        Attribute("meta", func() {
            Attribute("total", Int, "Total count")
            Attribute("page", Int, "Current page")
            Attribute("per_page", Int, "Items per page")
            Attribute("total_pages", Int, "Total pages")
            Required("total", "page", "per_page", "total_pages")
        })
        Required("data", "meta")
    })
    
    HTTP(func() {
        GET("/")
        Response(StatusOK)
    })
})
```

**Result HTTP Mapping:**

```go
Method("create", func() {
    Payload(CreateUserPayload)
    Result(User)
    
    HTTP(func() {
        POST("/")
        Response(StatusCreated, func() {
            // Custom header in response
            Header("location:Location")
        })
    })
})
```

---

### 🔹 Required Fields

**Theory:**

`Required()` specifies which attributes must be present. Missing required fields cause validation errors before your code even runs.

**Required Placement:**

```go
// In Type definition
var User = Type("User", func() {
    Attribute("id", String)
    Attribute("email", String)
    Attribute("name", String)
    Attribute("bio", String)  // Optional
    
    Required("id", "email", "name")  // bio is optional
})

// In Payload definition
Method("create", func() {
    Payload(func() {
        Attribute("email", String)
        Attribute("password", String)
        Attribute("name", String)
        Attribute("referral_code", String)  // Optional
        
        Required("email", "password", "name")
    })
})

// In Result definition
Method("get", func() {
    Result(func() {
        Attribute("user", User)
        Attribute("permissions", ArrayOf(String))
        
        Required("user", "permissions")
    })
})
```

**Generated Validation:**

When a required field is missing, Goa automatically returns a validation error:

```json
{
  "name": "bad_request",
  "message": "missing required field \"email\""
}
```

**Required vs Optional in HTTP:**

| Attribute | In Body | In Path | In Query |
|-----------|---------|---------|----------|
| Required | Must be in JSON | Always required | Error if missing |
| Optional | Can be omitted | N/A (path always required) | Uses default or nil |

---

### 🔹 Default Values

**Theory:**

Default values are used when an optional attribute is not provided. They reduce boilerplate in client code and make APIs more forgiving.

**Setting Defaults:**

```go
Method("list", func() {
    Payload(func() {
        // Numeric defaults
        Attribute("page", Int, "Page number", func() {
            Default(1)
            Minimum(1)
        })
        
        Attribute("limit", Int, "Items per page", func() {
            Default(20)
            Minimum(1)
            Maximum(100)
        })
        
        // String defaults
        Attribute("sort_by", String, "Sort field", func() {
            Default("created_at")
            Enum("created_at", "updated_at", "name")
        })
        
        Attribute("order", String, "Sort order", func() {
            Default("desc")
            Enum("asc", "desc")
        })
        
        // Boolean defaults
        Attribute("include_deleted", Boolean, func() {
            Default(false)
        })
    })
})
```

**Default Value Types:**

```go
// String default
Attribute("status", String, func() {
    Default("pending")
})

// Integer default
Attribute("retries", Int, func() {
    Default(3)
})

// Float default
Attribute("tax_rate", Float64, func() {
    Default(0.0)
})

// Boolean default
Attribute("active", Boolean, func() {
    Default(true)
})

// Array default (empty array)
Attribute("tags", ArrayOf(String), func() {
    Default([]string{})
})
```

**How Defaults Work:**

```
Client Request                    After Default Application
┌───────────────────────┐         ┌───────────────────────┐
│ GET /users            │   ──▶   │ page: 1               │
│ (no query params)     │         │ limit: 20             │
└───────────────────────┘         │ sort_by: "created_at" │
                                  │ order: "desc"         │
┌───────────────────────┐         └───────────────────────┘
│ GET /users?page=5     │   ──▶   ┌───────────────────────┐
│                       │         │ page: 5               │
└───────────────────────┘         │ limit: 20 (default)   │
                                  │ sort_by: "created_at" │
                                  │ order: "desc"         │
                                  └───────────────────────┘
```

---

### 🔹 Validation Rules

**Built-in Validation Functions:**

Goa provides extensive validation that runs before your code:

**String Validations:**

```go
Attribute("username", String, func() {
    MinLength(3)           // At least 3 characters
    MaxLength(50)          // At most 50 characters
    Pattern("^[a-z0-9_]+$") // Must match regex
})

Attribute("email", String, func() {
    Format(FormatEmail)    // Must be valid email
})

Attribute("website", String, func() {
    Format(FormatURI)      // Must be valid URI
})

Attribute("id", String, func() {
    Format(FormatUUID)     // Must be valid UUID
})

Attribute("role", String, func() {
    Enum("admin", "user", "guest")  // Must be one of these
})
```

**Numeric Validations:**

```go
Attribute("age", Int, func() {
    Minimum(0)             // Must be >= 0
    Maximum(150)           // Must be <= 150
})

Attribute("quantity", Int, func() {
    Minimum(1)             // At least 1
    Maximum(1000)          // At most 1000
})

Attribute("price", Float64, func() {
    Minimum(0.01)          // Positive price
})

Attribute("discount", Float64, func() {
    Minimum(0)             // 0% minimum
    Maximum(100)           // 100% maximum
})
```

**Array Validations:**

```go
Attribute("tags", ArrayOf(String), func() {
    MinLength(1)           // At least 1 tag
    MaxLength(10)          // At most 10 tags
})

Attribute("items", ArrayOf(OrderItem), func() {
    MinLength(1)           // Order must have items
})
```

**Format Constants:**

| Format | Description | Example |
|--------|-------------|---------|
| `FormatDate` | RFC3339 date | `2024-01-15` |
| `FormatDateTime` | RFC3339 datetime | `2024-01-15T10:30:00Z` |
| `FormatUUID` | UUID v4 | `550e8400-e29b-41d4-a716-446655440000` |
| `FormatEmail` | Email address | `user@example.com` |
| `FormatHostname` | Hostname | `example.com` |
| `FormatIPv4` | IPv4 address | `192.168.1.1` |
| `FormatIPv6` | IPv6 address | `::1` |
| `FormatIP` | IPv4 or IPv6 | Either format |
| `FormatURI` | URI | `https://example.com/path` |

**Complete Validation Example:**

```go
var CreateUserPayload = Type("CreateUserPayload", func() {
    Description("Payload for creating a new user")
    
    Attribute("email", String, "User's email address", func() {
        Format(FormatEmail)
        MaxLength(255)
        Example("user@example.com")
    })
    
    Attribute("password", String, "Account password", func() {
        MinLength(8)
        MaxLength(100)
        Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)`)  // Requires lowercase, uppercase, digit
        Example("SecurePass123")
    })
    
    Attribute("username", String, "Unique username", func() {
        MinLength(3)
        MaxLength(30)
        Pattern("^[a-z][a-z0-9_]*$")  // Starts with letter, alphanumeric + underscore
        Example("john_doe")
    })
    
    Attribute("age", Int, "User's age", func() {
        Minimum(13)   // Must be 13 or older
        Maximum(120)
    })
    
    Attribute("website", String, "Personal website", func() {
        Format(FormatURI)
        Example("https://johndoe.com")
    })
    
    Attribute("role", String, "User role", func() {
        Enum("user", "moderator", "admin")
        Default("user")
    })
    
    Required("email", "password", "username")
})
```

**Validation Error Response:**

When validation fails, Goa returns detailed errors:

```json
{
  "name": "bad_request",
  "id": "validation_error",
  "message": "password must be at least 8 characters long",
  "temporary": false,
  "timeout": false,
  "fault": false
}
```

---

## 4. Data Modeling

### Understanding Data Modeling: The Theory

**What is Data Modeling in Goa?**

Data modeling in Goa defines the structure of data flowing through your API. Good data models ensure:
- Clear contracts between client and server
- Automatic validation and serialization
- Generated documentation accuracy
- Type safety in generated code

**Model Hierarchy:**

```
┌─────────────────────────────────────────────────────────────────┐
│                        Type Definitions                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Primitive Types                                                 │
│  ├── Boolean, Int, Int32, Int64                                 │
│  ├── UInt, UInt32, UInt64                                       │
│  ├── Float32, Float64                                           │
│  ├── String, Bytes                                              │
│  └── Any                                                         │
│                                                                  │
│  Composite Types                                                 │
│  ├── Type()         → Named struct-like type                    │
│  ├── ArrayOf()      → Slice/List                                │
│  ├── MapOf()        → Key-value map                             │
│  └── ResultType()   → Response with views                       │
│                                                                  │
│  Special Types                                                   │
│  ├── ErrorResult    → Standardized error                        │
│  └── Any            → Dynamic type (avoid)                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

### 🔹 Attributes

**What are Attributes?**

Attributes are the fields of a type. Each attribute has:
- A name
- A type
- Optional description
- Optional validations and metadata

**Attribute Syntax:**

```go
// Basic syntax
Attribute("name", Type)

// With description
Attribute("name", Type, "Description text")

// With configuration function
Attribute("name", Type, func() {
    Description("Description text")
    // Validations, defaults, examples...
})

// Full form
Attribute("name", Type, "Description", func() {
    // Configuration
})
```

**Attribute Examples:**

```go
var Product = Type("Product", func() {
    Description("Product in the catalog")
    
    // Simple attributes
    Attribute("id", String)
    Attribute("name", String)
    Attribute("price", Float64)
    Attribute("quantity", Int)
    Attribute("active", Boolean)
    
    // With descriptions
    Attribute("sku", String, "Stock Keeping Unit")
    
    // With full configuration
    Attribute("category", String, "Product category", func() {
        Enum("electronics", "clothing", "food", "other")
        Default("other")
    })
    
    // Nested type reference
    Attribute("dimensions", Dimensions)
    
    // Array attribute
    Attribute("tags", ArrayOf(String))
    
    // Map attribute  
    Attribute("metadata", MapOf(String, String))
    
    Required("id", "name", "price")
})
```

**Attribute Modifiers:**

```go
Attribute("field", String, func() {
    // Documentation
    Description("Field description")
    
    // Validations (covered in Methods section)
    MinLength(1)
    MaxLength(100)
    Pattern("^[a-z]+$")
    Format(FormatEmail)
    Enum("value1", "value2")
    Minimum(0)
    Maximum(100)
    
    // Default value
    Default("default_value")
    
    // Example for documentation
    Example("example_value")
    
    // Metadata (custom key-value)
    Meta("key", "value")
})
```

---

### 🔹 Types

**Creating Named Types:**

```go
// Basic Type definition
var User = Type("User", func() {
    Description("A user of the system")
    
    Attribute("id", String, "Unique identifier", func() {
        Format(FormatUUID)
    })
    Attribute("email", String, "Email address", func() {
        Format(FormatEmail)
    })
    Attribute("name", String, "Display name")
    Attribute("created_at", String, "Creation timestamp", func() {
        Format(FormatDateTime)
    })
    
    Required("id", "email", "name", "created_at")
})
```

**Type Composition:**

```go
// Base type
var Address = Type("Address", func() {
    Attribute("street", String)
    Attribute("city", String)
    Attribute("country", String)
    Attribute("postal_code", String)
    Required("street", "city", "country")
})

// Type using another type
var Customer = Type("Customer", func() {
    Attribute("id", String)
    Attribute("name", String)
    Attribute("billing_address", Address)   // References Address type
    Attribute("shipping_address", Address)  // Can use multiple times
    Required("id", "name")
})
```

**Extending Types:**

```go
// Define base attributes in a function
var CommonFields = func() {
    Attribute("id", String, func() {
        Format(FormatUUID)
    })
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("updated_at", String, func() {
        Format(FormatDateTime)
    })
}

// Use in multiple types
var User = Type("User", func() {
    CommonFields()  // Include common fields
    
    Attribute("email", String)
    Attribute("name", String)
    
    Required("id", "email", "name", "created_at", "updated_at")
})

var Product = Type("Product", func() {
    CommonFields()  // Same common fields
    
    Attribute("name", String)
    Attribute("price", Float64)
    
    Required("id", "name", "price", "created_at", "updated_at")
})
```

**ResultType for Views:**

ResultType allows different views of the same data:

```go
var User = ResultType("User", func() {
    Description("User account")
    
    Attributes(func() {
        Attribute("id", String)
        Attribute("email", String)
        Attribute("name", String)
        Attribute("password_hash", String)  // Sensitive!
        Attribute("created_at", String)
        Attribute("updated_at", String)
        Attribute("last_login", String)
        
        Required("id", "email", "name")
    })
    
    // Default view - what most responses return
    View("default", func() {
        Attribute("id")
        Attribute("email")
        Attribute("name")
        Attribute("created_at")
    })
    
    // Full view - for admin or detailed responses
    View("full", func() {
        Attribute("id")
        Attribute("email")
        Attribute("name")
        Attribute("created_at")
        Attribute("updated_at")
        Attribute("last_login")
        // Note: password_hash excluded from all views!
    })
    
    // Minimal view - for lists
    View("tiny", func() {
        Attribute("id")
        Attribute("name")
    })
})
```

Using views in methods:

```go
Method("get", func() {
    Payload(func() {
        Attribute("id", String)
        Attribute("view", String, func() {
            Enum("default", "full", "tiny")
            Default("default")
        })
        Required("id")
    })
    
    Result(User, func() {
        View("default")  // Can be overridden based on request
    })
})
```

---

### 🔹 Custom Types

**When to Create Custom Types:**

- Reused across multiple services/methods
- Complex validation logic
- Domain concepts that deserve a name
- Improving API documentation clarity

**Custom Scalar Types:**

```go
// Custom string type with validation
var Email = Type("Email", String, func() {
    Description("Valid email address")
    Format(FormatEmail)
    MaxLength(255)
    Example("user@example.com")
})

var UUID = Type("UUID", String, func() {
    Description("Universally Unique Identifier")
    Format(FormatUUID)
    Example("550e8400-e29b-41d4-a716-446655440000")
})

var PhoneNumber = Type("PhoneNumber", String, func() {
    Description("International phone number")
    Pattern(`^\+[1-9]\d{1,14}$`)
    Example("+14155551234")
})

var Currency = Type("Currency", String, func() {
    Description("ISO 4217 currency code")
    Enum("USD", "EUR", "GBP", "JPY", "CNY")
    Example("USD")
})
```

**Using Custom Types:**

```go
var User = Type("User", func() {
    Attribute("id", UUID)           // Uses custom UUID type
    Attribute("email", Email)       // Uses custom Email type
    Attribute("phone", PhoneNumber) // Uses custom PhoneNumber type
    Required("id", "email")
})

var Transaction = Type("Transaction", func() {
    Attribute("id", UUID)
    Attribute("amount", Float64)
    Attribute("currency", Currency)
    Required("id", "amount", "currency")
})
```

**Custom Complex Types:**

```go
// Money type (amount + currency together)
var Money = Type("Money", func() {
    Description("Monetary amount with currency")
    
    Attribute("amount", Float64, "Amount in decimal", func() {
        Minimum(0)
        Example(99.99)
    })
    Attribute("currency", Currency, "Currency code")
    
    Required("amount", "currency")
})

// Date range type
var DateRange = Type("DateRange", func() {
    Description("A range of dates")
    
    Attribute("start", String, "Start date", func() {
        Format(FormatDate)
        Example("2024-01-01")
    })
    Attribute("end", String, "End date", func() {
        Format(FormatDate)
        Example("2024-12-31")
    })
    
    Required("start", "end")
})

// Pagination parameters
var PaginationParams = Type("PaginationParams", func() {
    Description("Standard pagination parameters")
    
    Attribute("page", Int, func() {
        Minimum(1)
        Default(1)
    })
    Attribute("per_page", Int, func() {
        Minimum(1)
        Maximum(100)
        Default(20)
    })
    Attribute("sort_by", String)
    Attribute("sort_order", String, func() {
        Enum("asc", "desc")
        Default("desc")
    })
})
```

---

### 🔹 Arrays & Maps

**Arrays (ArrayOf):**

```go
// Simple array
Attribute("tags", ArrayOf(String))

// Array with validation
Attribute("tags", ArrayOf(String), func() {
    MinLength(1)   // At least 1 item
    MaxLength(10)  // At most 10 items
})

// Array of custom type
Attribute("items", ArrayOf(OrderItem))

// Array of complex inline type
Attribute("addresses", ArrayOf(func() {
    Attribute("label", String)
    Attribute("street", String)
    Attribute("city", String)
    Required("street", "city")
}))
```

**Array Examples in Types:**

```go
var Order = Type("Order", func() {
    Attribute("id", UUID)
    Attribute("customer_id", UUID)
    
    // Array of order items
    Attribute("items", ArrayOf(OrderItem), func() {
        MinLength(1)  // Must have at least one item
    })
    
    // Array of strings
    Attribute("tags", ArrayOf(String))
    
    // Array of UUIDs
    Attribute("related_orders", ArrayOf(UUID))
    
    Required("id", "customer_id", "items")
})

var OrderItem = Type("OrderItem", func() {
    Attribute("product_id", UUID)
    Attribute("quantity", Int, func() {
        Minimum(1)
    })
    Attribute("price", Money)
    Required("product_id", "quantity", "price")
})
```

**Maps (MapOf):**

```go
// String to String map
Attribute("metadata", MapOf(String, String))

// String to Any (flexible values)
Attribute("properties", MapOf(String, Any))

// String to typed value
Attribute("settings", MapOf(String, Setting))

// Enum keys
Attribute("translations", MapOf(
    func() {
        Enum("en", "es", "fr", "de")
    }, 
    String,
))
```

**Map Examples in Types:**

```go
var Product = Type("Product", func() {
    Attribute("id", UUID)
    Attribute("name", String)
    
    // Key-value metadata
    Attribute("metadata", MapOf(String, String), func() {
        Description("Additional product metadata")
        Example(map[string]string{
            "color": "blue",
            "size":  "large",
        })
    })
    
    // Translations
    Attribute("translations", MapOf(String, ProductTranslation), func() {
        Description("Product name/description in different languages")
    })
    
    // Feature flags
    Attribute("features", MapOf(String, Boolean), func() {
        Description("Product feature toggles")
        Example(map[string]bool{
            "returnable":  true,
            "gift_wrap":   true,
            "fragile":     false,
        })
    })
    
    Required("id", "name")
})

var ProductTranslation = Type("ProductTranslation", func() {
    Attribute("name", String)
    Attribute("description", String)
    Required("name")
})
```

**Nested Arrays and Maps:**

```go
var Report = Type("Report", func() {
    // Array of arrays (matrix)
    Attribute("data", ArrayOf(ArrayOf(Float64)))
    
    // Map of arrays
    Attribute("categories", MapOf(String, ArrayOf(Product)))
    
    // Array of maps
    Attribute("records", ArrayOf(MapOf(String, Any)))
})
```

---

### 🔹 Metadata

**What is Metadata?**

Metadata adds arbitrary key-value pairs to types and attributes. It's used for:
- Documentation customization
- Code generation hints
- External tool integration
- Custom processing

**Meta Function:**

```go
var User = Type("User", func() {
    // Type-level metadata
    Meta("struct:tag:json", "user")
    Meta("openapi:tag", "Users")
    Meta("database:table", "users")
    
    Attribute("id", String, func() {
        // Attribute-level metadata
        Meta("struct:tag:json", "id,omitempty")
        Meta("struct:tag:gorm", "primaryKey")
        Meta("database:column", "user_id")
    })
    
    Attribute("email", String, func() {
        Meta("struct:tag:json", "email")
        Meta("struct:tag:gorm", "uniqueIndex")
    })
})
```

**Common Meta Keys:**

| Meta Key | Purpose | Example |
|----------|---------|---------|
| `struct:tag:json` | Custom JSON tag | `Meta("struct:tag:json", "user_id,omitempty")` |
| `struct:tag:*` | Any struct tag | `Meta("struct:tag:xml", "userId")` |
| `openapi:tag` | OpenAPI grouping | `Meta("openapi:tag", "Users")` |
| `openapi:summary` | OpenAPI summary | `Meta("openapi:summary", "Get user")` |
| `swagger:*` | Swagger extensions | `Meta("swagger:extension:x-custom", "value")` |

**Practical Meta Usage:**

```go
var _ = Service("users", func() {
    // Group in OpenAPI documentation
    Meta("openapi:tag:Users")
    
    Method("get", func() {
        Meta("openapi:summary", "Retrieve a user")
        Meta("openapi:operationId", "getUser")
        
        Payload(func() {
            Attribute("id", String, func() {
                Meta("openapi:example", "abc-123")
            })
            Required("id")
        })
        
        Result(User)
    })
})
```

---

### 🔹 Example Values

**Why Examples Matter:**

Examples appear in:
- OpenAPI/Swagger documentation
- Generated API clients
- Testing and mocking
- Developer onboarding

**Adding Examples:**

```go
// Simple example
Attribute("email", String, func() {
    Example("john@example.com")
})

// Multiple examples
Attribute("status", String, func() {
    Enum("pending", "active", "suspended")
    Example("pending")
})

// Complex type example
var User = Type("User", func() {
    Attribute("id", String, func() {
        Example("user_abc123")
    })
    Attribute("email", String, func() {
        Example("john.doe@example.com")
    })
    Attribute("name", String, func() {
        Example("John Doe")
    })
    Attribute("role", String, func() {
        Example("admin")
    })
    Attribute("created_at", String, func() {
        Example("2024-01-15T10:30:00Z")
    })
    
    Required("id", "email", "name", "created_at")
})
```

**Type-Level Examples:**

```go
var Address = Type("Address", func() {
    Description("Physical address")
    
    Attribute("street", String, func() {
        Example("123 Main Street")
    })
    Attribute("city", String, func() {
        Example("San Francisco")
    })
    Attribute("state", String, func() {
        Example("CA")
    })
    Attribute("zip", String, func() {
        Example("94102")
    })
    Attribute("country", String, func() {
        Example("USA")
    })
    
    Required("street", "city", "country")
    
    // Full example for the entire type
    Example(map[string]interface{}{
        "street":  "456 Oak Avenue, Suite 100",
        "city":    "New York",
        "state":   "NY",
        "zip":     "10001",
        "country": "USA",
    })
})
```

**Examples in Arrays and Maps:**

```go
Attribute("tags", ArrayOf(String), func() {
    Example([]string{"electronics", "featured", "sale"})
})

Attribute("metadata", MapOf(String, String), func() {
    Example(map[string]string{
        "color":    "blue",
        "size":     "medium",
        "material": "cotton",
    })
})
```

---

## 5. Security Basics

### Understanding Security in Goa

**Why Security Matters:**

Goa provides first-class support for API security. Security schemes are defined once and applied to methods, ensuring consistent authentication across your API.

**Security Scheme Types:**

| Scheme Type | Use Case | Example |
|-------------|----------|---------|
| `BasicAuthSecurity` | Username/password | Internal APIs |
| `APIKeySecurity` | API keys | Third-party integrations |
| `JWTSecurity` | JWT tokens | Modern web/mobile apps |
| `OAuth2Security` | OAuth 2.0 flows | Social logins, delegated auth |

---

### 🔹 Defining Security Schemes

**API Key Security:**

```go
// design/security.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// API Key in header
var APIKeyAuth = APIKeySecurity("api_key", func() {
    Description("API key authentication")
})

// API Key in query parameter
var APIKeyQueryAuth = APIKeySecurity("api_key_query", func() {
    Description("API key in query parameter")
    Query("api_key")  // ?api_key=xxx
})
```

**JWT Security:**

```go
// JWT Bearer token
var JWTAuth = JWTSecurity("jwt", func() {
    Description("JWT Bearer token authentication")
    Scope("api:read", "Read access")
    Scope("api:write", "Write access")
    Scope("admin", "Admin access")
})
```

**Basic Auth:**

```go
var BasicAuth = BasicAuthSecurity("basic", func() {
    Description("Basic authentication")
})
```

**OAuth2 Security:**

```go
var OAuth2Auth = OAuth2Security("oauth2", func() {
    Description("OAuth2 authentication")
    
    // Authorization code flow
    AuthorizationCodeFlow(
        "https://auth.example.com/authorize",
        "https://auth.example.com/token",
        "https://auth.example.com/refresh",
    )
    
    Scope("read:users", "Read user data")
    Scope("write:users", "Modify user data")
})
```

---

### 🔹 Applying Security to Methods

**Service-Level Security:**

```go
var _ = Service("users", func() {
    Description("User management service")
    
    // All methods in this service require JWT
    Security(JWTAuth)
    
    Method("list", func() {
        // Inherits JWTAuth from service
        Result(ArrayOf(User))
    })
    
    Method("get", func() {
        // Inherits JWTAuth from service
        Payload(func() {
            Attribute("id", String)
            Required("id")
        })
        Result(User)
    })
    
    // Public method - no auth required
    Method("health", func() {
        NoSecurity()  // Override service security
        Result(String)
    })
})
```

**Method-Level Security:**

```go
var _ = Service("admin", func() {
    Description("Admin operations")
    
    Method("dashboard", func() {
        // Requires JWT with admin scope
        Security(JWTAuth, func() {
            Scope("admin")
        })
        Result(DashboardData)
    })
    
    Method("stats", func() {
        // Multiple security options (OR logic)
        Security(JWTAuth)
        Security(APIKeyAuth)  // Either JWT OR API key works
        Result(StatsData)
    })
})
```

---

### 🔹 Security in Payload

**Token in Payload:**

```go
Method("protected", func() {
    Security(JWTAuth)
    
    Payload(func() {
        // Token attribute for security
        Token("token", String, "JWT token")
        
        // Other attributes
        Attribute("data", String)
        Required("token")
    })
    
    Result(String)
    
    HTTP(func() {
        POST("/protected")
        // Map token to Authorization header
        Header("token:Authorization")
    })
})
```

**API Key in Payload:**

```go
Method("external", func() {
    Security(APIKeyAuth)
    
    Payload(func() {
        APIKey("api_key", "key", String, "API Key")
        Attribute("query", String)
        Required("key", "query")
    })
    
    HTTP(func() {
        GET("/search")
        Header("key:X-API-Key")  // Key in custom header
        Param("query")
    })
})
```

---

### 🔹 Implementing Security Handler

When you implement your service, you need to provide security handlers:

```go
// In your main.go or http.go
package main

import (
    "context"
    "errors"
    
    "github.com/golang-jwt/jwt/v5"
    usersvc "mymodule/gen/users"
)

// JWTAuth handler
func JWTAuthHandler(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
    // Parse and validate JWT
    claims := &jwt.RegisteredClaims{}
    parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })
    
    if err != nil || !parsedToken.Valid {
        return ctx, usersvc.MakeUnauthorized(errors.New("invalid token"))
    }
    
    // Add claims to context
    ctx = context.WithValue(ctx, "user_id", claims.Subject)
    
    return ctx, nil
}

// APIKey handler
func APIKeyAuthHandler(ctx context.Context, key string, scheme *security.APIKeyScheme) (context.Context, error) {
    // Validate API key (check against database, etc.)
    if key != "valid-api-key" {
        return ctx, errors.New("invalid API key")
    }
    return ctx, nil
}
```

---

## 6. Service Implementation

### Understanding Service Implementation

After running `goa gen`, you get generated interfaces that you must implement. This is where your business logic lives.

---

### 🔹 Generated Interface

Goa generates a service interface based on your design:

```go
// gen/users/service.go (generated - DO NOT EDIT)
package users

import (
    "context"
)

// Service describes the users service interface.
type Service interface {
    // List returns all users.
    List(context.Context, *ListPayload) (*ListResult, error)
    // Get returns a user by ID.
    Get(context.Context, *GetPayload) (*User, error)
    // Create creates a new user.
    Create(context.Context, *CreatePayload) (*User, error)
    // Update modifies an existing user.
    Update(context.Context, *UpdatePayload) (*User, error)
    // Delete removes a user.
    Delete(context.Context, *DeletePayload) error
}
```

---

### 🔹 Implementing the Service

Create a file (e.g., `users.go`) to implement the interface:

```go
// users.go - YOUR code
package main

import (
    "context"
    "log"
    "sync"
    "time"
    
    "github.com/google/uuid"
    users "mymodule/gen/users"
)

// usersService implements the users.Service interface
type usersService struct {
    logger *log.Logger
    mu     sync.RWMutex
    users  map[string]*users.User  // In-memory storage (use DB in production)
}

// NewUsersService creates a new users service
func NewUsersService(logger *log.Logger) users.Service {
    return &usersService{
        logger: logger,
        users:  make(map[string]*users.User),
    }
}

// List returns all users
func (s *usersService) List(ctx context.Context, p *users.ListPayload) (*users.ListResult, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    s.logger.Printf("Listing users, page: %d, limit: %d", p.Page, p.Limit)
    
    // Convert map to slice
    userList := make([]*users.User, 0, len(s.users))
    for _, u := range s.users {
        userList = append(userList, u)
    }
    
    // Calculate pagination
    total := len(userList)
    pages := (total + *p.Limit - 1) / *p.Limit
    
    // Apply pagination
    start := (*p.Page - 1) * *p.Limit
    end := start + *p.Limit
    if start > len(userList) {
        userList = []*users.User{}
    } else if end > len(userList) {
        userList = userList[start:]
    } else {
        userList = userList[start:end]
    }
    
    return &users.ListResult{
        Users: userList,
        Total: total,
        Page:  *p.Page,
        Pages: pages,
    }, nil
}

// Get returns a user by ID
func (s *usersService) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    s.logger.Printf("Getting user: %s", p.ID)
    
    user, ok := s.users[p.ID]
    if !ok {
        return nil, users.MakeNotFound(fmt.Errorf("user %s not found", p.ID))
    }
    
    return user, nil
}

// Create creates a new user
func (s *usersService) Create(ctx context.Context, p *users.CreatePayload) (*users.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Printf("Creating user: %s", p.Email)
    
    // Check if email already exists
    for _, u := range s.users {
        if u.Email == p.Email {
            return nil, users.MakeConflict(fmt.Errorf("email %s already taken", p.Email))
        }
    }
    
    // Create new user
    now := time.Now().UTC().Format(time.RFC3339)
    user := &users.User{
        ID:        uuid.New().String(),
        Email:     p.Email,
        Name:      p.Name,
        CreatedAt: now,
    }
    
    s.users[user.ID] = user
    
    return user, nil
}

// Update modifies an existing user
func (s *usersService) Update(ctx context.Context, p *users.UpdatePayload) (*users.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Printf("Updating user: %s", p.ID)
    
    user, ok := s.users[p.ID]
    if !ok {
        return nil, users.MakeNotFound(fmt.Errorf("user %s not found", p.ID))
    }
    
    // Update fields
    if p.Name != nil {
        user.Name = *p.Name
    }
    if p.Email != nil {
        user.Email = *p.Email
    }
    
    return user, nil
}

// Delete removes a user
func (s *usersService) Delete(ctx context.Context, p *users.DeletePayload) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Printf("Deleting user: %s", p.ID)
    
    if _, ok := s.users[p.ID]; !ok {
        return users.MakeNotFound(fmt.Errorf("user %s not found", p.ID))
    }
    
    delete(s.users, p.ID)
    
    return nil
}
```

---

### 🔹 Returning Errors

Goa generates error constructors for each error defined in your design:

```go
// Generated error constructors (in gen/users/service.go)
// MakeNotFound builds a users service not_found error.
func MakeNotFound(err error) *goa.ServiceError { ... }

// MakeUnauthorized builds a users service unauthorized error.
func MakeUnauthorized(err error) *goa.ServiceError { ... }

// MakeValidation builds a users service validation error.
func MakeValidation(err error) *goa.ServiceError { ... }
```

**Using Error Constructors:**

```go
func (s *usersService) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    user, err := s.db.FindUser(p.ID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // Return "not_found" error defined in design
            return nil, users.MakeNotFound(fmt.Errorf("user %s not found", p.ID))
        }
        // Return "internal" error for unexpected errors
        return nil, users.MakeInternal(err)
    }
    return user, nil
}
```

---

### 🔹 Accessing Security Context

When using security, access authenticated user info from context:

```go
func (s *usersService) GetProfile(ctx context.Context, p *users.GetProfilePayload) (*users.User, error) {
    // Get user ID from JWT claims (set in security handler)
    userID, ok := ctx.Value("user_id").(string)
    if !ok {
        return nil, users.MakeUnauthorized(errors.New("user ID not in context"))
    }
    
    // Use userID to fetch profile
    return s.getUserByID(userID)
}
```

---

## 7. Complete Working Example

### Full Example: User Service

Let's create a complete, working Goa service from design to implementation.

---

### 🔹 Step 1: Project Setup

```bash
# Create project directory
mkdir goa-users-api && cd goa-users-api

# Initialize Go module
go mod init github.com/username/goa-users-api

# Install Goa
go install goa.design/goa/v3/cmd/goa@latest

# Create design directory
mkdir design
```

---

### 🔹 Step 2: Design Files

**design/design.go:**

```go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = API("usersapi", func() {
    Title("Users API")
    Description("A simple user management API built with Goa")
    Version("1.0")
    
    Server("usersapi", func() {
        Host("localhost", func() {
            URI("http://localhost:8080")
        })
    })
})
```

**design/types.go:**

```go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// User represents a user in the system
var User = Type("User", func() {
    Description("A user of the system")
    
    Attribute("id", String, "Unique user ID", func() {
        Example("550e8400-e29b-41d4-a716-446655440000")
    })
    Attribute("email", String, "Email address", func() {
        Format(FormatEmail)
        Example("john@example.com")
    })
    Attribute("name", String, "Full name", func() {
        MinLength(1)
        MaxLength(100)
        Example("John Doe")
    })
    Attribute("created_at", String, "Creation timestamp", func() {
        Format(FormatDateTime)
        Example("2024-01-15T10:30:00Z")
    })
    
    Required("id", "email", "name", "created_at")
})

// ErrorResponse for all errors
var ErrorResponse = Type("ErrorResponse", func() {
    Attribute("code", String, "Error code", func() {
        Example("USER_NOT_FOUND")
    })
    Attribute("message", String, "Error message", func() {
        Example("The requested user was not found")
    })
    Required("code", "message")
})
```

**design/users.go:**

```go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("User management service")
    
    // Service-level errors
    Error("not_found", ErrorResponse, "User not found")
    Error("bad_request", ErrorResponse, "Invalid request")
    Error("conflict", ErrorResponse, "User already exists")
    
    HTTP(func() {
        Path("/users")
    })
    
    // List all users
    Method("list", func() {
        Description("List all users with pagination")
        
        Payload(func() {
            Attribute("page", Int, "Page number", func() {
                Minimum(1)
                Default(1)
            })
            Attribute("limit", Int, "Items per page", func() {
                Minimum(1)
                Maximum(100)
                Default(20)
            })
        })
        
        Result(func() {
            Attribute("users", ArrayOf(User), "List of users")
            Attribute("total", Int, "Total count")
            Attribute("page", Int, "Current page")
            Attribute("pages", Int, "Total pages")
            Required("users", "total", "page", "pages")
        })
        
        HTTP(func() {
            GET("/")
            Param("page")
            Param("limit")
            Response(StatusOK)
        })
    })
    
    // Get a single user
    Method("get", func() {
        Description("Get a user by ID")
        
        Payload(func() {
            Attribute("id", String, "User ID", func() {
                Example("550e8400-e29b-41d4-a716-446655440000")
            })
            Required("id")
        })
        
        Result(User)
        
        Error("not_found")
        
        HTTP(func() {
            GET("/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Create a new user
    Method("create", func() {
        Description("Create a new user")
        
        Payload(func() {
            Attribute("email", String, "Email address", func() {
                Format(FormatEmail)
            })
            Attribute("name", String, "Full name", func() {
                MinLength(1)
                MaxLength(100)
            })
            Required("email", "name")
        })
        
        Result(User)
        
        Error("bad_request")
        Error("conflict")
        
        HTTP(func() {
            POST("/")
            Body(func() {
                Attribute("email")
                Attribute("name")
            })
            Response(StatusCreated)
            Response("bad_request", StatusBadRequest)
            Response("conflict", StatusConflict)
        })
    })
    
    // Update a user
    Method("update", func() {
        Description("Update an existing user")
        
        Payload(func() {
            Attribute("id", String, "User ID")
            Attribute("email", String, "New email", func() {
                Format(FormatEmail)
            })
            Attribute("name", String, "New name", func() {
                MinLength(1)
                MaxLength(100)
            })
            Required("id")
        })
        
        Result(User)
        
        Error("not_found")
        Error("bad_request")
        
        HTTP(func() {
            PUT("/{id}")
            Body(func() {
                Attribute("email")
                Attribute("name")
            })
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("bad_request", StatusBadRequest)
        })
    })
    
    // Delete a user
    Method("delete", func() {
        Description("Delete a user")
        
        Payload(func() {
            Attribute("id", String, "User ID")
            Required("id")
        })
        
        Error("not_found")
        
        HTTP(func() {
            DELETE("/{id}")
            Response(StatusNoContent)
            Response("not_found", StatusNotFound)
        })
    })
})
```

---

### 🔹 Step 3: Generate Code

```bash
# Generate code from design
goa gen github.com/username/goa-users-api/design

# Generate example implementation
goa example github.com/username/goa-users-api/design
```

**Generated structure:**

```
goa-users-api/
├── design/
│   ├── design.go
│   ├── types.go
│   └── users.go
├── gen/
│   ├── users/
│   │   ├── client.go
│   │   ├── endpoints.go
│   │   ├── service.go
│   │   └── types.go
│   └── http/
│       ├── users/
│       │   ├── client/
│       │   ├── server/
│       │   └── ...
│       ├── openapi.json
│       └── openapi3.json
├── cmd/
│   └── usersapi/
│       ├── main.go
│       └── http.go
└── users.go  (skeleton - edit this!)
```

---

### 🔹 Step 4: Implement Service

Edit `users.go` (generated skeleton):

```go
package usersapi

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "github.com/google/uuid"
    users "github.com/username/goa-users-api/gen/users"
)

// userssrvc implements users.Service
type userssrvc struct {
    logger *log.Logger
    mu     sync.RWMutex
    store  map[string]*users.User
}

// NewUsers returns the users service implementation.
func NewUsers(logger *log.Logger) users.Service {
    return &userssrvc{
        logger: logger,
        store:  make(map[string]*users.User),
    }
}

// List returns all users with pagination
func (s *userssrvc) List(ctx context.Context, p *users.ListPayload) (*users.ListResult, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    s.logger.Printf("users.list called with page=%d, limit=%d", *p.Page, *p.Limit)
    
    // Collect all users
    all := make([]*users.User, 0, len(s.store))
    for _, u := range s.store {
        all = append(all, u)
    }
    
    total := len(all)
    limit := *p.Limit
    page := *p.Page
    pages := (total + limit - 1) / limit
    if pages == 0 {
        pages = 1
    }
    
    // Paginate
    start := (page - 1) * limit
    end := start + limit
    if start >= len(all) {
        return &users.ListResult{
            Users: []*users.User{},
            Total: total,
            Page:  page,
            Pages: pages,
        }, nil
    }
    if end > len(all) {
        end = len(all)
    }
    
    return &users.ListResult{
        Users: all[start:end],
        Total: total,
        Page:  page,
        Pages: pages,
    }, nil
}

// Get returns a user by ID
func (s *userssrvc) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    s.logger.Printf("users.get called for ID=%s", p.ID)
    
    user, exists := s.store[p.ID]
    if !exists {
        return nil, users.MakeNotFound(fmt.Errorf("user %q not found", p.ID))
    }
    
    return user, nil
}

// Create creates a new user
func (s *userssrvc) Create(ctx context.Context, p *users.CreatePayload) (*users.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Printf("users.create called for email=%s", p.Email)
    
    // Check for duplicate email
    for _, u := range s.store {
        if u.Email == p.Email {
            return nil, users.MakeConflict(fmt.Errorf("user with email %q already exists", p.Email))
        }
    }
    
    // Create user
    user := &users.User{
        ID:        uuid.New().String(),
        Email:     p.Email,
        Name:      p.Name,
        CreatedAt: time.Now().UTC().Format(time.RFC3339),
    }
    
    s.store[user.ID] = user
    
    return user, nil
}

// Update updates an existing user
func (s *userssrvc) Update(ctx context.Context, p *users.UpdatePayload) (*users.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Printf("users.update called for ID=%s", p.ID)
    
    user, exists := s.store[p.ID]
    if !exists {
        return nil, users.MakeNotFound(fmt.Errorf("user %q not found", p.ID))
    }
    
    if p.Email != nil {
        user.Email = *p.Email
    }
    if p.Name != nil {
        user.Name = *p.Name
    }
    
    return user, nil
}

// Delete removes a user
func (s *userssrvc) Delete(ctx context.Context, p *users.DeletePayload) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.logger.Printf("users.delete called for ID=%s", p.ID)
    
    if _, exists := s.store[p.ID]; !exists {
        return users.MakeNotFound(fmt.Errorf("user %q not found", p.ID))
    }
    
    delete(s.store, p.ID)
    
    return nil
}
```

---

### 🔹 Step 5: Run the Service

```bash
# Install dependencies
go mod tidy

# Run the service
go run ./cmd/usersapi
```

---

### 🔹 Step 6: Test the API

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "name": "John Doe"}'

# List users
curl http://localhost:8080/users

# Get a specific user
curl http://localhost:8080/users/{id}

# Update a user
curl -X PUT http://localhost:8080/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'

# Delete a user
curl -X DELETE http://localhost:8080/users/{id}
```

---

### 🔹 Step 7: View OpenAPI Spec

The generated OpenAPI spec is at `gen/http/openapi3.json`. You can:

1. **Serve it with Swagger UI:**
   ```bash
   # Add static file serving or use online Swagger Editor
   ```

2. **Import into Postman:**
   - Import `openapi3.json` directly

3. **Generate client SDKs:**
   - Use OpenAPI Generator for other languages

---

## 📝 Summary: Key Takeaways

### Design-First Philosophy
- **Design before code**: Define API contract in DSL, then implement
- **Generated code**: Never edit `gen/` folder, regenerate after design changes
- **Project structure**: `design/` (yours), `gen/` (generated), `cmd/` (entry points)
- **DSL is Go**: Uses regular Go functions with closures

### Services
- **Logical grouping**: Related methods under one service
- **Service-level errors**: Define once, reuse in methods
- **Multiple services**: Separate domains into different services
- **Shared types**: Define in `types.go`, use across services

### Methods
- **Payload**: Input from body, path, query, headers
- **Result**: Output on success, can have views
- **Required**: Must be present, validated automatically
- **Default**: Used when optional field is missing
- **Validation**: Built-in rules for strings, numbers, arrays

### Data Modeling
- **Types**: Named structures with attributes
- **Custom types**: Reusable validated scalars and composites
- **Arrays**: `ArrayOf(Type)` with length validation
- **Maps**: `MapOf(KeyType, ValueType)` for dictionaries
- **Metadata**: Custom key-value for generation hints
- **Examples**: Documentation and testing values

### Security
- **Security schemes**: API Key, JWT, Basic Auth, OAuth2
- **Service-level security**: Apply to all methods
- **NoSecurity()**: Exempt specific methods from auth
- **Token in Payload**: Map to headers for HTTP transport
- **Security handlers**: Implement authentication logic

### Service Implementation
- **Generated interface**: Implement the `Service` interface
- **Error constructors**: Use `MakeNotFound()`, `MakeUnauthorized()`, etc.
- **Context for auth**: Access user info from security handlers
- **Business logic**: Keep implementation separate from design

---

## 📋 Knowledge Check

Before proceeding to Part 3 (HTTP Transport), ensure you can:

- [ ] Explain the difference between design-first and code-first
- [ ] Configure API() with servers, hosts, and metadata
- [ ] Run `goa gen` and understand generated file structure
- [ ] Define a service with multiple methods
- [ ] Create Payload and Result definitions
- [ ] Use Required, Default, and validation rules
- [ ] Create reusable Types and custom types
- [ ] Use ArrayOf and MapOf for collections
- [ ] Add metadata and examples for documentation
- [ ] Define and apply security schemes
- [ ] Implement the generated service interface
- [ ] Return appropriate errors using generated constructors
- [ ] Build a complete API from design to running service

---

## 🔗 Quick Reference Links

- [Goa Documentation](https://goa.design/)
- [Goa DSL Reference](https://pkg.go.dev/goa.design/goa/v3/dsl)
- [Goa Examples](https://github.com/goadesign/examples)
- [Goa Security Examples](https://github.com/goadesign/examples/tree/master/security)
- [OpenAPI Specification](https://swagger.io/specification/)

---

> **Next Up:** Part 3 - HTTP Transport (Routes, Request/Response Mapping, Status Codes)
