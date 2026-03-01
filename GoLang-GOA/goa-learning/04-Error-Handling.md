# Part 4: Error Handling (Very Important)

> **Goal:** Master Goa's error handling - defining errors, custom types, transport mapping, and response structures

---

## 📚 Table of Contents

1. [Error Handling Philosophy](#error-handling-philosophy)
2. [Defining Errors in Design](#defining-errors-in-design)
   - [Service-Level Errors](#service-level-errors)
   - [Method-Level Errors](#method-level-errors)
   - [Error Inheritance](#error-inheritance)
3. [Custom Error Types](#custom-error-types)
   - [Default ErrorResult](#default-errorresult)
   - [Custom Error Structures](#custom-error-structures)
   - [Error with Multiple Fields](#error-with-multiple-fields)
   - [Nested Error Types](#nested-error-types)
4. [HTTP Status Mapping](#http-status-mapping)
   - [Standard Mappings](#standard-mappings)
   - [Custom Response Bodies](#custom-response-bodies)
   - [Response Headers on Error](#response-headers-on-error)
5. [gRPC Error Mapping](#grpc-error-mapping)
   - [gRPC Status Codes](#grpc-status-codes)
   - [Error Details](#error-details)
   - [Rich Error Model](#rich-error-model)
6. [Error Response Structure](#error-response-structure)
   - [Generated Error Types](#generated-error-types)
   - [Error Constructors](#error-constructors)
   - [Returning Errors](#returning-errors)
7. [Error Handling Patterns](#error-handling-patterns)
   - [Validation Errors](#validation-errors)
   - [Business Logic Errors](#business-logic-errors)
   - [External Service Errors](#external-service-errors)
   - [Retry Hints](#retry-hints)
8. [Complete Examples](#complete-examples)
9. [Best Practices](#best-practices)
10. [Summary](#summary)

---

## 🎯 Error Handling Philosophy

### Why Errors Matter in API Design

Errors are as important as successful responses. Well-designed errors help clients:
- Understand what went wrong
- Know if they can retry
- Fix their request
- Present meaningful messages to users

```
┌─────────────────────────────────────────────────────────────────┐
│                    ERROR HANDLING FLOW                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Client Request                                                 │
│       │                                                         │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────┐               │
│  │       Transport Layer (HTTP/gRPC)           │               │
│  │  - Decode request                           │               │
│  │  - Validate format                          │───────┐       │
│  └─────────────────────────────────────────────┘       │       │
│       │                                                │       │
│       ▼                                                │       │
│  ┌─────────────────────────────────────────────┐       │       │
│  │          Validation Layer                   │       │       │
│  │  - Check required fields                    │       │       │
│  │  - Validate patterns, ranges                │───────┤       │
│  └─────────────────────────────────────────────┘       │       │
│       │                                                │       │
│       ▼                                                │       │
│  ┌─────────────────────────────────────────────┐       │       │
│  │         Business Logic Layer                │       │       │
│  │  - Authorization checks                     │       │       │
│  │  - Domain rules                             │───────┤       │
│  │  - Resource operations                      │       │       │
│  └─────────────────────────────────────────────┘       │       │
│       │                                                │       │
│       ▼                                                ▼       │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Error Response                         │   │
│  │  - Appropriate status code                               │   │
│  │  - Structured error body                                 │   │
│  │  - Actionable error message                              │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Error Categories

| Category | Description | HTTP Status | Retryable |
|----------|-------------|-------------|-----------|
| **Client Error** | Invalid request from client | 4xx | No (fix request) |
| **Validation Error** | Input doesn't meet requirements | 400, 422 | No (fix input) |
| **Authentication Error** | Missing/invalid credentials | 401 | No (authenticate) |
| **Authorization Error** | Insufficient permissions | 403 | No (get permission) |
| **Not Found Error** | Resource doesn't exist | 404 | Maybe (create first) |
| **Conflict Error** | State conflict | 409 | Maybe (resolve conflict) |
| **Rate Limit Error** | Too many requests | 429 | Yes (wait) |
| **Server Error** | Internal failure | 5xx | Yes (retry) |
| **Unavailable Error** | Service temporarily down | 503 | Yes (retry later) |

---

## 📝 Defining Errors in Design

### Service-Level Errors

Service-level errors are defined once and shared across all methods in the service.

```go
// design/users.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("User management service")
    
    // Service-level errors - available to all methods
    Error("not_found", ErrorResult, "Resource not found")
    Error("unauthorized", ErrorResult, "Authentication required")
    Error("forbidden", ErrorResult, "Insufficient permissions")
    Error("bad_request", ErrorResult, "Invalid request data")
    Error("conflict", ErrorResult, "Resource already exists")
    Error("internal", ErrorResult, "Internal server error")
    
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int)
            Required("id")
        })
        Result(User)
        
        // Reference service-level errors
        Error("not_found")
        Error("unauthorized")
        
        HTTP(func() {
            GET("/users/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    Method("create", func() {
        Payload(CreateUserPayload)
        Result(User)
        
        // Use different subset of errors
        Error("bad_request")
        Error("conflict")
        Error("unauthorized")
        
        HTTP(func() {
            POST("/users")
            Response(StatusCreated)
            Response("bad_request", StatusBadRequest)
            Response("conflict", StatusConflict)
            Response("unauthorized", StatusUnauthorized)
        })
    })
})
```

### Method-Level Errors

Define errors specific to individual methods.

```go
Method("transfer", func() {
    Description("Transfer funds between accounts")
    
    Payload(TransferPayload)
    Result(TransferResult)
    
    // Method-specific errors
    Error("insufficient_funds", ErrorResult, "Not enough balance")
    Error("account_locked", ErrorResult, "Account is locked")
    Error("daily_limit_exceeded", ErrorResult, "Daily transfer limit reached")
    Error("same_account", ErrorResult, "Cannot transfer to same account")
    
    // Also can use service-level errors
    Error("not_found")      // Source or destination account not found
    Error("unauthorized")   // Not authorized for this account
    
    HTTP(func() {
        POST("/accounts/{from_account_id}/transfer")
        
        Response(StatusOK)
        Response("not_found", StatusNotFound)
        Response("unauthorized", StatusUnauthorized)
        Response("insufficient_funds", StatusPaymentRequired)       // 402
        Response("account_locked", StatusForbidden)                 // 403
        Response("daily_limit_exceeded", StatusTooManyRequests)     // 429
        Response("same_account", StatusBadRequest)                  // 400
    })
})
```

### Error Inheritance

```
┌─────────────────────────────────────────────────────────────────┐
│                    ERROR INHERITANCE                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  API Level                                                      │
│  ───────────                                                    │
│  var _ = API("myapi", func() {                                  │
│      // Errors defined here are NOT inherited                   │
│      // API-level settings are for metadata only                │
│  })                                                             │
│                                                                 │
│  Service Level                                                  │
│  ─────────────                                                  │
│  var _ = Service("users", func() {                              │
│      Error("not_found")      ──────────────────────────┐        │
│      Error("unauthorized")   ──────────────────────────┤        │
│      Error("internal")       ──────────────────────────┤        │
│                                                        │        │
│      Method("get", func() {                            │        │
│          Error("not_found")     ◄──── References ──────┤        │
│          Error("unauthorized")  ◄──── service errors ──┘        │
│          // Can also define method-specific errors              │
│          Error("deleted", DeletedError)                         │
│      })                                                         │
│                                                                 │
│      Method("create", func() {                                  │
│          Error("unauthorized")  ◄──── Can use subset            │
│          Error("conflict")      ◄──── New error for method      │
│      })                                                         │
│  })                                                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Shared Error Definitions

Create a separate file for reusable error types.

```go
// design/errors.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// Standard error result used across services
var ErrorResult = Type("ErrorResult", func() {
    Description("Standard error response")
    
    Attribute("code", String, "Error code", func() {
        Example("USER_NOT_FOUND")
    })
    Attribute("message", String, "Human-readable error message", func() {
        Example("The requested user does not exist")
    })
    Attribute("id", String, "Unique error instance ID for tracking", func() {
        Example("err-abc123")
    })
    
    Required("code", "message")
})

// Validation error with field details
var ValidationError = Type("ValidationError", func() {
    Description("Validation error with field details")
    
    Attribute("code", String, "Error code", func() {
        Default("VALIDATION_ERROR")
    })
    Attribute("message", String, "General error message")
    Attribute("fields", ArrayOf(FieldError), "Field-specific errors")
    
    Required("code", "message", "fields")
})

var FieldError = Type("FieldError", func() {
    Attribute("field", String, "Field name", func() {
        Example("email")
    })
    Attribute("error", String, "Error description", func() {
        Example("must be a valid email address")
    })
    Attribute("value", Any, "The invalid value provided")
    
    Required("field", "error")
})

// Rate limit error with retry info
var RateLimitError = Type("RateLimitError", func() {
    Description("Rate limit exceeded error")
    
    Attribute("code", String, func() {
        Default("RATE_LIMIT_EXCEEDED")
    })
    Attribute("message", String)
    Attribute("retry_after", Int, "Seconds until retry is allowed")
    Attribute("limit", Int, "Request limit")
    Attribute("remaining", Int, "Remaining requests")
    Attribute("reset_at", String, "When the limit resets (RFC3339)")
    
    Required("code", "message", "retry_after")
})
```

---

## 🔧 Custom Error Types

### Default ErrorResult

Goa provides a default `ErrorResult` type if you don't specify one.

```go
// If you write:
Error("not_found")

// Goa uses its built-in error type equivalent to:
var GoaErrorResult = Type("Error", func() {
    Attribute("name", String, "Name of error")
    Attribute("id", String, "Unique error ID")
    Attribute("message", String, "Error message")
    Attribute("temporary", Boolean, "Is error temporary")
    Attribute("timeout", Boolean, "Is error a timeout")
    Attribute("fault", Boolean, "Is error a server fault")
    Required("name", "id", "message")
})
```

### Custom Error Structures

Define your own error types for better API documentation.

```go
// Simple custom error
var NotFoundError = Type("NotFoundError", func() {
    Description("Resource not found error")
    
    Attribute("resource_type", String, "Type of resource", func() {
        Enum("user", "post", "comment", "file")
        Example("user")
    })
    Attribute("resource_id", String, "ID of missing resource", func() {
        Example("usr_123")
    })
    Attribute("message", String, "Error message", func() {
        Example("User not found")
    })
    
    Required("resource_type", "resource_id", "message")
})

// Use in service
Method("get", func() {
    Payload(func() {
        Attribute("id", String)
        Required("id")
    })
    Result(User)
    
    // Use custom error type
    Error("not_found", NotFoundError, "User not found")
    
    HTTP(func() {
        GET("/users/{id}")
        Response(StatusOK)
        Response("not_found", StatusNotFound)
    })
})
```

### Error with Multiple Fields

```go
var DetailedError = Type("DetailedError", func() {
    Description("Detailed error response with context")
    
    // Core error information
    Attribute("code", String, "Machine-readable error code", func() {
        Pattern("^[A-Z][A-Z0-9_]*$")
        Example("INVALID_CREDENTIALS")
    })
    
    Attribute("message", String, "Human-readable message", func() {
        Example("The provided credentials are invalid")
    })
    
    // Error classification
    Attribute("type", String, "Error category", func() {
        Enum("validation", "authentication", "authorization", 
             "not_found", "conflict", "rate_limit", "internal")
    })
    
    // Debugging information
    Attribute("request_id", String, "Request tracking ID", func() {
        Example("req_abc123xyz")
    })
    
    Attribute("timestamp", String, "When error occurred", func() {
        Format(FormatDateTime)
    })
    
    // Help for developers
    Attribute("documentation_url", String, "Link to error documentation", func() {
        Format(FormatURI)
        Example("https://api.example.com/docs/errors/INVALID_CREDENTIALS")
    })
    
    // Additional context
    Attribute("details", MapOf(String, Any), "Additional error context")
    
    Required("code", "message", "type", "request_id", "timestamp")
})
```

### Nested Error Types

```go
// Error with nested validation details
var ValidationErrorResponse = Type("ValidationErrorResponse", func() {
    Attribute("code", String, func() {
        Default("VALIDATION_FAILED")
    })
    Attribute("message", String, func() {
        Default("Request validation failed")
    })
    Attribute("errors", ArrayOf(ValidationDetail), "List of validation errors")
    
    Required("code", "message", "errors")
})

var ValidationDetail = Type("ValidationDetail", func() {
    Attribute("path", String, "JSON path to invalid field", func() {
        Example("$.user.email")
    })
    Attribute("field", String, "Field name", func() {
        Example("email")
    })
    Attribute("constraint", String, "Violated constraint", func() {
        Enum("required", "format", "min_length", "max_length", 
             "minimum", "maximum", "pattern", "enum")
        Example("format")
    })
    Attribute("message", String, "Error description", func() {
        Example("must be a valid email address")
    })
    Attribute("expected", Any, "Expected value/format")
    Attribute("actual", Any, "Actual value provided")
    
    Required("field", "constraint", "message")
})

// Error with nested resource reference
var ConflictError = Type("ConflictError", func() {
    Attribute("code", String, func() {
        Default("RESOURCE_CONFLICT")
    })
    Attribute("message", String)
    Attribute("existing_resource", ResourceReference, "Reference to conflicting resource")
    
    Required("code", "message", "existing_resource")
})

var ResourceReference = Type("ResourceReference", func() {
    Attribute("type", String, "Resource type", func() {
        Example("user")
    })
    Attribute("id", String, "Resource ID", func() {
        Example("usr_123")
    })
    Attribute("url", String, "Resource URL", func() {
        Format(FormatURI)
        Example("https://api.example.com/users/usr_123")
    })
    
    Required("type", "id")
})
```

---

## 🌐 HTTP Status Mapping

### Standard Mappings

```go
var _ = Service("products", func() {
    // Define all possible errors at service level
    Error("bad_request", ValidationErrorResponse)
    Error("unauthorized", ErrorResult)
    Error("forbidden", ErrorResult)
    Error("not_found", NotFoundError)
    Error("conflict", ConflictError)
    Error("gone", ErrorResult)
    Error("unprocessable", ValidationErrorResponse)
    Error("rate_limited", RateLimitError)
    Error("internal", ErrorResult)
    Error("unavailable", ErrorResult)
    
    Method("create", func() {
        Payload(CreateProductPayload)
        Result(Product)
        
        Error("bad_request")
        Error("unauthorized")
        Error("forbidden")
        Error("conflict")
        Error("unprocessable")
        Error("rate_limited")
        Error("internal")
        
        HTTP(func() {
            POST("/products")
            
            // Success response
            Response(StatusCreated)
            
            // Client errors (4xx)
            Response("bad_request", StatusBadRequest)           // 400
            Response("unauthorized", StatusUnauthorized)        // 401
            Response("forbidden", StatusForbidden)              // 403
            Response("conflict", StatusConflict)                // 409
            Response("unprocessable", StatusUnprocessableEntity)// 422
            Response("rate_limited", StatusTooManyRequests)     // 429
            
            // Server errors (5xx)
            Response("internal", StatusInternalServerError)     // 500
        })
    })
})
```

### HTTP Status Code Reference

```
┌─────────────────────────────────────────────────────────────────┐
│              HTTP STATUS CODE DECISION TREE                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Was the request successful?                                    │
│  │                                                              │
│  ├── YES → 2xx Success                                          │
│  │   ├── Resource returned? → 200 OK                            │
│  │   ├── Resource created? → 201 Created                        │
│  │   ├── Accepted for processing? → 202 Accepted                │
│  │   └── No content to return? → 204 No Content                 │
│  │                                                              │
│  └── NO → Was it client's fault?                                │
│      │                                                          │
│      ├── YES → 4xx Client Error                                 │
│      │   ├── Malformed request? → 400 Bad Request               │
│      │   ├── Not authenticated? → 401 Unauthorized              │
│      │   ├── Not authorized? → 403 Forbidden                    │
│      │   ├── Resource not found? → 404 Not Found                │
│      │   ├── Method not allowed? → 405 Method Not Allowed       │
│      │   ├── State conflict? → 409 Conflict                     │
│      │   ├── Resource deleted? → 410 Gone                       │
│      │   ├── Validation failed? → 422 Unprocessable Entity      │
│      │   └── Too many requests? → 429 Too Many Requests         │
│      │                                                          │
│      └── NO → 5xx Server Error                                  │
│          ├── Unexpected error? → 500 Internal Server Error      │
│          ├── Not implemented? → 501 Not Implemented             │
│          ├── Bad gateway? → 502 Bad Gateway                     │
│          ├── Service down? → 503 Service Unavailable            │
│          └── Timeout upstream? → 504 Gateway Timeout            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Custom Response Bodies

```go
Method("get", func() {
    Payload(func() {
        Attribute("id", String)
        Required("id")
    })
    Result(Product)
    
    Error("not_found", NotFoundError)
    Error("deleted", DeletedError)
    
    HTTP(func() {
        GET("/products/{id}")
        
        Response(StatusOK)
        
        // Custom body for 404
        Response("not_found", StatusNotFound, func() {
            // Specify which fields go in the body
            Body(func() {
                Attribute("resource_type")
                Attribute("resource_id")
                Attribute("message")
            })
        })
        
        // Custom body for 410 (Gone)
        Response("deleted", StatusGone, func() {
            Body(func() {
                Attribute("deleted_at")
                Attribute("deleted_by")
                Attribute("message")
            })
        })
    })
})

var DeletedError = Type("DeletedError", func() {
    Attribute("message", String)
    Attribute("deleted_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("deleted_by", String, "User who deleted")
    Required("message", "deleted_at")
})
```

### Response Headers on Error

```go
Method("create", func() {
    Payload(CreatePayload)
    Result(Resource)
    
    Error("rate_limited", RateLimitError)
    Error("conflict", ConflictError)
    
    HTTP(func() {
        POST("/resources")
        
        Response(StatusCreated, func() {
            Header("Location")
        })
        
        // Rate limit error with headers
        Response("rate_limited", StatusTooManyRequests, func() {
            // Standard rate limit headers
            Header("X-RateLimit-Limit:limit")
            Header("X-RateLimit-Remaining:remaining")
            Header("X-RateLimit-Reset:reset_at")
            Header("Retry-After:retry_after")
            
            Body(func() {
                Attribute("code")
                Attribute("message")
            })
        })
        
        // Conflict with location of existing resource
        Response("conflict", StatusConflict, func() {
            Header("Location:existing_url")
            
            Body(func() {
                Attribute("code")
                Attribute("message")
                Attribute("existing_resource")
            })
        })
    })
})

var RateLimitError = Type("RateLimitError", func() {
    Attribute("code", String)
    Attribute("message", String)
    Attribute("limit", Int)
    Attribute("remaining", Int)
    Attribute("reset_at", String)
    Attribute("retry_after", Int)
    Required("code", "message", "retry_after")
})

var ConflictError = Type("ConflictError", func() {
    Attribute("code", String)
    Attribute("message", String)
    Attribute("existing_resource", ResourceRef)
    Attribute("existing_url", String)
    Required("code", "message")
})
```

---

## 📡 gRPC Error Mapping

### gRPC Status Codes

```go
var _ = Service("users", func() {
    Error("not_found", ErrorResult)
    Error("invalid", ValidationError)
    Error("unauthorized", ErrorResult)
    Error("forbidden", ErrorResult)
    Error("conflict", ErrorResult)
    Error("internal", ErrorResult)
    Error("unavailable", ErrorResult)
    
    Method("get", func() {
        Payload(GetPayload)
        Result(User)
        
        Error("not_found")
        Error("unauthorized")
        Error("forbidden")
        
        HTTP(func() {
            GET("/users/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("unauthorized", StatusUnauthorized)
            Response("forbidden", StatusForbidden)
        })
        
        GRPC(func() {
            Response(CodeOK)
            Response("not_found", CodeNotFound)
            Response("unauthorized", CodeUnauthenticated)
            Response("forbidden", CodePermissionDenied)
        })
    })
    
    Method("create", func() {
        Payload(CreatePayload)
        Result(User)
        
        Error("invalid")
        Error("conflict")
        Error("internal")
        
        HTTP(func() {
            POST("/users")
            Response(StatusCreated)
            Response("invalid", StatusBadRequest)
            Response("conflict", StatusConflict)
            Response("internal", StatusInternalServerError)
        })
        
        GRPC(func() {
            Response(CodeOK)
            Response("invalid", CodeInvalidArgument)
            Response("conflict", CodeAlreadyExists)
            Response("internal", CodeInternal)
        })
    })
})
```

### gRPC to HTTP Status Code Mapping

| gRPC Code | HTTP Status | Use Case |
|-----------|-------------|----------|
| `CodeOK` | 200 | Success |
| `CodeCanceled` | 499 | Client cancelled request |
| `CodeUnknown` | 500 | Unknown error |
| `CodeInvalidArgument` | 400 | Invalid request data |
| `CodeDeadlineExceeded` | 504 | Request timeout |
| `CodeNotFound` | 404 | Resource not found |
| `CodeAlreadyExists` | 409 | Resource conflict |
| `CodePermissionDenied` | 403 | Authorization failure |
| `CodeResourceExhausted` | 429 | Rate limit exceeded |
| `CodeFailedPrecondition` | 400 | Invalid state for operation |
| `CodeAborted` | 409 | Concurrency conflict |
| `CodeOutOfRange` | 400 | Value out of valid range |
| `CodeUnimplemented` | 501 | Method not implemented |
| `CodeInternal` | 500 | Internal server error |
| `CodeUnavailable` | 503 | Service unavailable |
| `CodeDataLoss` | 500 | Data corruption |
| `CodeUnauthenticated` | 401 | Missing authentication |

### Error Details

```go
// Rich error details for gRPC
var GRPCError = Type("GRPCError", func() {
    Attribute("code", String, "Error code")
    Attribute("message", String, "Error message")
    
    // Standard gRPC error detail types
    Attribute("error_info", ErrorInfo, "Structured error information")
    Attribute("retry_info", RetryInfo, "Retry guidance")
    Attribute("debug_info", DebugInfo, "Debugging information")
    Attribute("bad_request", BadRequestDetails, "Field violations")
    
    Required("code", "message")
})

var ErrorInfo = Type("ErrorInfo", func() {
    Attribute("reason", String, "Error reason code")
    Attribute("domain", String, "Error domain")
    Attribute("metadata", MapOf(String, String), "Additional metadata")
    Required("reason", "domain")
})

var RetryInfo = Type("RetryInfo", func() {
    Attribute("retry_delay_seconds", Int, "Recommended retry delay")
    Required("retry_delay_seconds")
})

var DebugInfo = Type("DebugInfo", func() {
    Attribute("stack_entries", ArrayOf(String), "Stack trace")
    Attribute("detail", String, "Debug details")
})

var BadRequestDetails = Type("BadRequestDetails", func() {
    Attribute("field_violations", ArrayOf(FieldViolation))
    Required("field_violations")
})

var FieldViolation = Type("FieldViolation", func() {
    Attribute("field", String, "Field path")
    Attribute("description", String, "Violation description")
    Required("field", "description")
})
```

### Rich Error Model

```go
Method("create", func() {
    Payload(CreatePayload)
    Result(User)
    
    // Use rich error type
    Error("validation_failed", ValidationErrorResponse)
    
    GRPC(func() {
        Response(CodeOK)
        
        // Rich error with details
        Response("validation_failed", CodeInvalidArgument, func() {
            // gRPC supports detailed error messages
            Message(func() {
                Attribute("code")
                Attribute("message")
                Attribute("errors")  // Array of validation details
            })
        })
    })
})
```

---

## 📋 Error Response Structure

### Generated Error Types

When you define errors in DSL, Goa generates:

```go
// gen/users/service.go

// Error constructors
func MakeNotFound(err error) *goa.ServiceError {
    return goa.NewServiceError(err, "not_found", false, false, false)
}

func MakeUnauthorized(err error) *goa.ServiceError {
    return goa.NewServiceError(err, "unauthorized", false, false, false)
}

func MakeBadRequest(err error) *goa.ServiceError {
    return goa.NewServiceError(err, "bad_request", false, false, false)
}

func MakeInternal(err error) *goa.ServiceError {
    return goa.NewServiceError(err, "internal", false, false, true)  // fault=true
}

func MakeUnavailable(err error) *goa.ServiceError {
    return goa.NewServiceError(err, "unavailable", true, false, false)  // temporary=true
}

func MakeTimeout(err error) *goa.ServiceError {
    return goa.NewServiceError(err, "timeout", true, true, false)  // temporary=true, timeout=true
}
```

### Error Constructors

The `Make<ErrorName>` functions create properly typed errors:

```go
// Basic usage - just an error message
return nil, users.MakeNotFound(fmt.Errorf("user %d not found", userID))

// With custom error type
return nil, users.MakeNotFound(&users.NotFoundError{
    ResourceType: "user",
    ResourceID:   fmt.Sprintf("%d", userID),
    Message:      "User not found",
})
```

### Returning Errors

#### Simple Error Return

```go
// users.go - Implementation
func (s *usersSvc) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    user, err := s.store.GetUser(ctx, p.ID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, users.MakeNotFound(
                fmt.Errorf("user with ID %d not found", p.ID),
            )
        }
        // Unexpected error - mark as internal/fault
        return nil, users.MakeInternal(
            fmt.Errorf("failed to fetch user: %w", err),
        )
    }
    return user, nil
}
```

#### Custom Error Type Return

```go
// With custom error type defined in DSL
func (s *usersSvc) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    user, err := s.store.GetUser(ctx, p.ID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // Return custom error type
            return nil, &users.NotFoundError{
                ResourceType: "user",
                ResourceID:   strconv.Itoa(p.ID),
                Message:      fmt.Sprintf("User %d does not exist", p.ID),
            }
        }
        return nil, users.MakeInternal(err)
    }
    return user, nil
}
```

#### Validation Error Return

```go
func (s *usersSvc) Create(ctx context.Context, p *users.CreatePayload) (*users.User, error) {
    var fieldErrors []*users.FieldError
    
    // Custom validation beyond DSL
    if !isValidDomain(p.Email) {
        fieldErrors = append(fieldErrors, &users.FieldError{
            Field: "email",
            Error: "email domain is not allowed",
            Value: p.Email,
        })
    }
    
    if existingUser, _ := s.store.GetByEmail(ctx, p.Email); existingUser != nil {
        fieldErrors = append(fieldErrors, &users.FieldError{
            Field: "email",
            Error: "email already in use",
            Value: p.Email,
        })
    }
    
    if len(fieldErrors) > 0 {
        return nil, &users.ValidationError{
            Code:    "VALIDATION_FAILED",
            Message: "Request validation failed",
            Fields:  fieldErrors,
        }
    }
    
    // Continue with creation...
    return s.store.CreateUser(ctx, p)
}
```

---

## 🎨 Error Handling Patterns

### Validation Errors

```go
// design/errors.go
var ValidationError = Type("ValidationError", func() {
    Attribute("code", String, func() {
        Default("VALIDATION_FAILED")
    })
    Attribute("message", String, func() {
        Default("One or more fields failed validation")
    })
    Attribute("errors", ArrayOf(FieldValidationError))
    Required("code", "message", "errors")
})

var FieldValidationError = Type("FieldValidationError", func() {
    Attribute("field", String, "Field name or path")
    Attribute("rule", String, "Validation rule that failed", func() {
        Enum("required", "format", "min", "max", "pattern", "enum", "unique")
    })
    Attribute("message", String, "Human-readable error message")
    Attribute("params", MapOf(String, Any), "Rule parameters")
    Required("field", "rule", "message")
})

// Implementation
func (s *svc) Create(ctx context.Context, p *service.CreatePayload) (*service.Resource, error) {
    var errors []*service.FieldValidationError
    
    // Validate email format (beyond DSL)
    if p.Email != "" && !isValidEmail(p.Email) {
        errors = append(errors, &service.FieldValidationError{
            Field:   "email",
            Rule:    "format",
            Message: "Invalid email format",
            Params:  map[string]any{"expected": "email"},
        })
    }
    
    // Check uniqueness
    if exists, _ := s.store.EmailExists(ctx, p.Email); exists {
        errors = append(errors, &service.FieldValidationError{
            Field:   "email",
            Rule:    "unique",
            Message: "Email already registered",
        })
    }
    
    // Validate password strength
    if len(p.Password) < 8 {
        errors = append(errors, &service.FieldValidationError{
            Field:   "password",
            Rule:    "min",
            Message: "Password must be at least 8 characters",
            Params:  map[string]any{"min": 8, "actual": len(p.Password)},
        })
    }
    
    if len(errors) > 0 {
        return nil, &service.ValidationError{
            Code:    "VALIDATION_FAILED",
            Message: "Request validation failed",
            Errors:  errors,
        }
    }
    
    return s.store.Create(ctx, p)
}
```

### Business Logic Errors

```go
// domain-specific errors
var InsufficientFundsError = Type("InsufficientFundsError", func() {
    Attribute("code", String, func() {
        Default("INSUFFICIENT_FUNDS")
    })
    Attribute("message", String)
    Attribute("available_balance", Float64, "Current balance")
    Attribute("required_amount", Float64, "Amount needed")
    Attribute("currency", String, "Currency code")
    Required("code", "message", "available_balance", "required_amount", "currency")
})

var AccountLockedError = Type("AccountLockedError", func() {
    Attribute("code", String, func() {
        Default("ACCOUNT_LOCKED")
    })
    Attribute("message", String)
    Attribute("locked_at", String, func() { Format(FormatDateTime) })
    Attribute("reason", String, "Lock reason", func() {
        Enum("fraud_suspected", "too_many_failed_attempts", "manual_lock", "compliance")
    })
    Attribute("unlock_at", String, "When account unlocks (if time-based)", func() {
        Format(FormatDateTime)
    })
    Attribute("support_ticket", String, "Support ticket for manual unlock")
    Required("code", "message", "locked_at", "reason")
})

// Implementation
func (s *bankingSvc) Transfer(ctx context.Context, p *banking.TransferPayload) (*banking.TransferResult, error) {
    // Check account status
    account, err := s.store.GetAccount(ctx, p.FromAccountID)
    if err != nil {
        return nil, banking.MakeNotFound(err)
    }
    
    if account.IsLocked {
        return nil, &banking.AccountLockedError{
            Code:      "ACCOUNT_LOCKED",
            Message:   "Source account is locked",
            LockedAt:  account.LockedAt.Format(time.RFC3339),
            Reason:    account.LockReason,
            UnlockAt:  account.UnlockAt.Format(time.RFC3339),
        }
    }
    
    // Check balance
    if account.Balance < p.Amount {
        return nil, &banking.InsufficientFundsError{
            Code:             "INSUFFICIENT_FUNDS",
            Message:          "Not enough balance for this transfer",
            AvailableBalance: account.Balance,
            RequiredAmount:   p.Amount,
            Currency:         account.Currency,
        }
    }
    
    // Process transfer...
    return s.processTransfer(ctx, account, p)
}
```

### External Service Errors

```go
// Wrapping external service failures
var ExternalServiceError = Type("ExternalServiceError", func() {
    Attribute("code", String)
    Attribute("message", String)
    Attribute("service", String, "Name of failed service")
    Attribute("retry_after", Int, "Seconds to wait before retry")
    Attribute("fallback_used", Boolean, "Whether fallback was used")
    Required("code", "message", "service")
})

// Implementation
func (s *orderSvc) Create(ctx context.Context, p *orders.CreatePayload) (*orders.Order, error) {
    // Call payment service
    paymentResult, err := s.paymentClient.Charge(ctx, p.PaymentMethod, p.Amount)
    if err != nil {
        var netErr net.Error
        if errors.As(err, &netErr) && netErr.Timeout() {
            return nil, &orders.ExternalServiceError{
                Code:       "PAYMENT_SERVICE_TIMEOUT",
                Message:    "Payment service timed out",
                Service:    "payment-service",
                RetryAfter: 5,
            }
        }
        
        // Check if it's a known payment error
        var paymentErr *payment.Error
        if errors.As(err, &paymentErr) {
            return nil, &orders.PaymentFailedError{
                Code:         paymentErr.Code,
                Message:      paymentErr.Message,
                DeclineCode:  paymentErr.DeclineCode,
            }
        }
        
        // Unknown error - wrap as internal
        return nil, orders.MakeInternal(fmt.Errorf("payment failed: %w", err))
    }
    
    return s.store.CreateOrder(ctx, p, paymentResult)
}
```

### Retry Hints

```go
// Error types with retry information
Error("rate_limited", RateLimitError, "Too many requests")
Error("unavailable", UnavailableError, "Service temporarily unavailable")
Error("timeout", TimeoutError, "Request timed out")

// Use Goa's built-in error flags
func (s *svc) Process(ctx context.Context, p *service.ProcessPayload) error {
    result, err := s.externalService.Call(ctx, p.Data)
    if err != nil {
        var netErr net.Error
        if errors.As(err, &netErr) {
            if netErr.Timeout() {
                // Timeout error - retryable
                return service.MakeTimeout(fmt.Errorf("external service timeout: %w", err))
            }
            if netErr.Temporary() {
                // Temporary error - retryable
                return service.MakeUnavailable(fmt.Errorf("external service unavailable: %w", err))
            }
        }
        // Permanent failure
        return service.MakeInternal(err)
    }
    return nil
}

// Client-side retry logic based on error flags
func callWithRetry(ctx context.Context, client service.Client, payload *service.Payload) (*service.Result, error) {
    var lastErr error
    
    for attempt := 0; attempt < 3; attempt++ {
        result, err := client.Process(ctx, payload)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Check if error is retryable
        var serviceErr *goa.ServiceError
        if errors.As(err, &serviceErr) {
            if !serviceErr.Temporary && !serviceErr.Timeout {
                // Not retryable - return immediately
                return nil, err
            }
            
            // Calculate backoff
            backoff := time.Duration(attempt+1) * time.Second
            
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(backoff):
                continue
            }
        }
        
        // Unknown error type - don't retry
        return nil, err
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

---

## 📦 Complete Examples

### Complete Error Design

```go
// design/errors.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// ============================================================
// STANDARD ERROR TYPES
// ============================================================

// Generic error for simple cases
var StandardError = Type("StandardError", func() {
    Description("Standard error response")
    Attribute("code", String, "Machine-readable error code")
    Attribute("message", String, "Human-readable error message")
    Attribute("request_id", String, "Request tracking ID")
    Required("code", "message")
})

// Validation error with field details
var ValidationError = Type("ValidationError", func() {
    Description("Validation error with field-level details")
    Attribute("code", String, func() { Default("VALIDATION_FAILED") })
    Attribute("message", String, func() { Default("Request validation failed") })
    Attribute("request_id", String)
    Attribute("errors", ArrayOf(FieldError), "Field-specific errors")
    Required("code", "message", "errors")
})

var FieldError = Type("FieldError", func() {
    Attribute("field", String, "Field name or path")
    Attribute("code", String, "Error code for this field")
    Attribute("message", String, "Error description")
    Attribute("value", Any, "Invalid value (optional)")
    Required("field", "code", "message")
})

// Not found error
var NotFoundError = Type("NotFoundError", func() {
    Attribute("code", String, func() { Default("NOT_FOUND") })
    Attribute("message", String)
    Attribute("resource_type", String, "Type of resource")
    Attribute("resource_id", String, "ID of missing resource")
    Attribute("request_id", String)
    Required("code", "message", "resource_type", "resource_id")
})

// Conflict error
var ConflictError = Type("ConflictError", func() {
    Attribute("code", String, func() { Default("CONFLICT") })
    Attribute("message", String)
    Attribute("conflicting_resource", String, "ID of conflicting resource")
    Attribute("conflict_field", String, "Field causing conflict")
    Attribute("request_id", String)
    Required("code", "message")
})

// Rate limit error
var RateLimitError = Type("RateLimitError", func() {
    Attribute("code", String, func() { Default("RATE_LIMIT_EXCEEDED") })
    Attribute("message", String)
    Attribute("limit", Int, "Request limit")
    Attribute("remaining", Int, "Remaining requests")
    Attribute("reset_at", Int64, "Unix timestamp when limit resets")
    Attribute("retry_after", Int, "Seconds until retry is allowed")
    Attribute("request_id", String)
    Required("code", "message", "retry_after")
})

// Internal error (minimal info for security)
var InternalError = Type("InternalError", func() {
    Attribute("code", String, func() { Default("INTERNAL_ERROR") })
    Attribute("message", String, func() { Default("An internal error occurred") })
    Attribute("request_id", String, "Use this ID when contacting support")
    Required("code", "message", "request_id")
})
```

```go
// design/products.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("products", func() {
    Description("Product catalog service")
    
    // ============================================================
    // SERVICE-LEVEL ERRORS
    // ============================================================
    Error("not_found", NotFoundError, "Product not found")
    Error("validation_error", ValidationError, "Validation failed")
    Error("conflict", ConflictError, "Product SKU already exists")
    Error("unauthorized", StandardError, "Authentication required")
    Error("forbidden", StandardError, "Insufficient permissions")
    Error("rate_limited", RateLimitError, "Too many requests")
    Error("internal", InternalError, "Internal server error")
    
    // ============================================================
    // METHODS
    // ============================================================
    
    Method("list", func() {
        Description("List products with pagination")
        
        Payload(func() {
            Attribute("page", Int, func() { Default(1); Minimum(1) })
            Attribute("limit", Int, func() { Default(20); Minimum(1); Maximum(100) })
            Attribute("category", String)
        })
        
        Result(ProductList)
        
        Error("validation_error")
        Error("rate_limited")
        Error("internal")
        
        HTTP(func() {
            GET("/products")
            Param("page")
            Param("limit")
            Param("category")
            
            Response(StatusOK)
            Response("validation_error", StatusBadRequest)
            Response("rate_limited", StatusTooManyRequests, func() {
                Header("X-RateLimit-Limit:limit")
                Header("X-RateLimit-Remaining:remaining")
                Header("X-RateLimit-Reset:reset_at")
                Header("Retry-After:retry_after")
            })
            Response("internal", StatusInternalServerError)
        })
        
        GRPC(func() {
            Response(CodeOK)
            Response("validation_error", CodeInvalidArgument)
            Response("rate_limited", CodeResourceExhausted)
            Response("internal", CodeInternal)
        })
    })
    
    Method("get", func() {
        Description("Get a product by ID")
        
        Payload(func() {
            Attribute("id", String, "Product ID")
            Required("id")
        })
        
        Result(Product)
        
        Error("not_found")
        Error("internal")
        
        HTTP(func() {
            GET("/products/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("internal", StatusInternalServerError)
        })
        
        GRPC(func() {
            Response(CodeOK)
            Response("not_found", CodeNotFound)
            Response("internal", CodeInternal)
        })
    })
    
    Method("create", func() {
        Description("Create a new product")
        
        Payload(func() {
            Attribute("auth_token", String, "Bearer token")
            Attribute("sku", String, "Product SKU", func() {
                Pattern("^[A-Z0-9-]+$")
                MinLength(3)
                MaxLength(50)
            })
            Attribute("name", String, func() {
                MinLength(1)
                MaxLength(200)
            })
            Attribute("price", Float64, func() {
                Minimum(0)
            })
            Attribute("category", String)
            Required("auth_token", "sku", "name", "price")
        })
        
        Result(Product)
        
        Error("validation_error")
        Error("conflict")
        Error("unauthorized")
        Error("forbidden")
        Error("internal")
        
        HTTP(func() {
            POST("/products")
            Header("Authorization:auth_token")
            
            Response(StatusCreated, func() {
                Header("Location")
            })
            Response("validation_error", StatusBadRequest)
            Response("conflict", StatusConflict)
            Response("unauthorized", StatusUnauthorized)
            Response("forbidden", StatusForbidden)
            Response("internal", StatusInternalServerError)
        })
        
        GRPC(func() {
            Metadata(func() {
                Attribute("auth_token:authorization")
            })
            
            Response(CodeOK)
            Response("validation_error", CodeInvalidArgument)
            Response("conflict", CodeAlreadyExists)
            Response("unauthorized", CodeUnauthenticated)
            Response("forbidden", CodePermissionDenied)
            Response("internal", CodeInternal)
        })
    })
    
    Method("delete", func() {
        Description("Delete a product")
        
        Payload(func() {
            Attribute("id", String, "Product ID")
            Attribute("auth_token", String)
            Required("id", "auth_token")
        })
        
        Error("not_found")
        Error("unauthorized")
        Error("forbidden")
        Error("internal")
        
        HTTP(func() {
            DELETE("/products/{id}")
            Header("Authorization:auth_token")
            
            Response(StatusNoContent)
            Response("not_found", StatusNotFound)
            Response("unauthorized", StatusUnauthorized)
            Response("forbidden", StatusForbidden)
            Response("internal", StatusInternalServerError)
        })
        
        GRPC(func() {
            Metadata(func() {
                Attribute("auth_token:authorization")
            })
            
            Response(CodeOK)
            Response("not_found", CodeNotFound)
            Response("unauthorized", CodeUnauthenticated)
            Response("forbidden", CodePermissionDenied)
            Response("internal", CodeInternal)
        })
    })
})
```

### Complete Implementation

```go
// products.go
package api

import (
    "context"
    "fmt"
    "strings"
    "sync"
    "time"
    
    "github.com/google/uuid"
    products "myproject/gen/products"
)

type productsSvc struct {
    mu       sync.RWMutex
    products map[string]*products.Product
    bySKU    map[string]string
    
    // Rate limiting (simple in-memory)
    requests    map[string]int
    requestTime map[string]time.Time
}

func NewProductsService() products.Service {
    return &productsSvc{
        products:    make(map[string]*products.Product),
        bySKU:       make(map[string]string),
        requests:    make(map[string]int),
        requestTime: make(map[string]time.Time),
    }
}

// requestID extracts or generates request ID from context
func requestID(ctx context.Context) string {
    if id, ok := ctx.Value("request_id").(string); ok {
        return id
    }
    return uuid.New().String()
}

func (s *productsSvc) List(ctx context.Context, p *products.ListPayload) (*products.ProductList, error) {
    reqID := requestID(ctx)
    
    // Check rate limit
    if err := s.checkRateLimit(ctx, reqID); err != nil {
        return nil, err
    }
    
    // Validate (beyond DSL validation)
    if p.Category != nil && !isValidCategory(*p.Category) {
        return nil, &products.ValidationError{
            Code:      "VALIDATION_FAILED",
            Message:   "Invalid category",
            RequestID: &reqID,
            Errors: []*products.FieldError{
                {
                    Field:   "category",
                    Code:    "INVALID_VALUE",
                    Message: fmt.Sprintf("'%s' is not a valid category", *p.Category),
                },
            },
        }
    }
    
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var result []*products.Product
    for _, product := range s.products {
        if p.Category != nil && product.Category != *p.Category {
            continue
        }
        result = append(result, product)
    }
    
    // Pagination
    total := len(result)
    start := (*p.Page - 1) * *p.Limit
    end := start + *p.Limit
    if start > total {
        start = total
    }
    if end > total {
        end = total
    }
    
    return &products.ProductList{
        Products: result[start:end],
        Total:    total,
        Page:     *p.Page,
        Limit:    *p.Limit,
    }, nil
}

func (s *productsSvc) Get(ctx context.Context, p *products.GetPayload) (*products.Product, error) {
    reqID := requestID(ctx)
    
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    product, ok := s.products[p.ID]
    if !ok {
        return nil, &products.NotFoundError{
            Code:         "NOT_FOUND",
            Message:      fmt.Sprintf("Product '%s' not found", p.ID),
            ResourceType: "product",
            ResourceID:   p.ID,
            RequestID:    &reqID,
        }
    }
    
    return product, nil
}

func (s *productsSvc) Create(ctx context.Context, p *products.CreatePayload) (*products.Product, error) {
    reqID := requestID(ctx)
    
    // Validate auth (simplified - real implementation would verify JWT)
    if p.AuthToken == "" || !strings.HasPrefix(p.AuthToken, "Bearer ") {
        return nil, &products.StandardError{
            Code:      "UNAUTHORIZED",
            Message:   "Valid bearer token required",
            RequestID: &reqID,
        }
    }
    
    // Business validation
    var fieldErrors []*products.FieldError
    
    // Validate SKU format
    if !isValidSKU(p.Sku) {
        fieldErrors = append(fieldErrors, &products.FieldError{
            Field:   "sku",
            Code:    "INVALID_FORMAT",
            Message: "SKU must contain only uppercase letters, numbers, and hyphens",
        })
    }
    
    // Validate price
    if p.Price <= 0 {
        fieldErrors = append(fieldErrors, &products.FieldError{
            Field:   "price",
            Code:    "INVALID_VALUE",
            Message: "Price must be greater than 0",
        })
    }
    
    if len(fieldErrors) > 0 {
        return nil, &products.ValidationError{
            Code:      "VALIDATION_FAILED",
            Message:   "Request validation failed",
            RequestID: &reqID,
            Errors:    fieldErrors,
        }
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Check for duplicate SKU
    if existingID, exists := s.bySKU[p.Sku]; exists {
        return nil, &products.ConflictError{
            Code:               "CONFLICT",
            Message:            fmt.Sprintf("Product with SKU '%s' already exists", p.Sku),
            ConflictingResource: &existingID,
            ConflictField:      strPtr("sku"),
            RequestID:          &reqID,
        }
    }
    
    // Create product
    id := uuid.New().String()
    now := time.Now().Format(time.RFC3339)
    
    product := &products.Product{
        ID:        id,
        Sku:       p.Sku,
        Name:      p.Name,
        Price:     p.Price,
        Category:  p.Category,
        CreatedAt: now,
        UpdatedAt: now,
    }
    
    s.products[id] = product
    s.bySKU[p.Sku] = id
    
    return product, nil
}

func (s *productsSvc) Delete(ctx context.Context, p *products.DeletePayload) error {
    reqID := requestID(ctx)
    
    // Validate auth
    if p.AuthToken == "" {
        return &products.StandardError{
            Code:      "UNAUTHORIZED",
            Message:   "Authentication required",
            RequestID: &reqID,
        }
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    product, ok := s.products[p.ID]
    if !ok {
        return &products.NotFoundError{
            Code:         "NOT_FOUND",
            Message:      fmt.Sprintf("Product '%s' not found", p.ID),
            ResourceType: "product",
            ResourceID:   p.ID,
            RequestID:    &reqID,
        }
    }
    
    delete(s.bySKU, product.Sku)
    delete(s.products, p.ID)
    
    return nil
}

// Rate limiting helper
func (s *productsSvc) checkRateLimit(ctx context.Context, reqID string) error {
    // Simplified rate limiting - real implementation would use Redis, etc.
    const limit = 100
    const window = time.Minute
    
    clientIP := "default" // Would extract from context in real implementation
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now()
    lastRequest, exists := s.requestTime[clientIP]
    
    if !exists || now.Sub(lastRequest) > window {
        s.requests[clientIP] = 1
        s.requestTime[clientIP] = now
        return nil
    }
    
    s.requests[clientIP]++
    
    if s.requests[clientIP] > limit {
        resetAt := lastRequest.Add(window)
        retryAfter := int(resetAt.Sub(now).Seconds())
        if retryAfter < 1 {
            retryAfter = 1
        }
        
        return &products.RateLimitError{
            Code:       "RATE_LIMIT_EXCEEDED",
            Message:    "Too many requests, please slow down",
            Limit:      intPtr(limit),
            Remaining:  intPtr(0),
            ResetAt:    int64Ptr(resetAt.Unix()),
            RetryAfter: retryAfter,
            RequestID:  &reqID,
        }
    }
    
    return nil
}

// Helper functions
func isValidSKU(sku string) bool {
    for _, c := range sku {
        if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
            return false
        }
    }
    return true
}

func isValidCategory(cat string) bool {
    validCategories := map[string]bool{
        "electronics": true,
        "clothing":    true,
        "books":       true,
        "home":        true,
    }
    return validCategories[cat]
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int { return &i }
func int64Ptr(i int64) *int64 { return &i }
```

---

## ✅ Best Practices

### 1. Be Specific with Error Types

```go
// ❌ Bad - Generic error
Error("error", ErrorResult)

// ✅ Good - Specific errors
Error("not_found", NotFoundError)
Error("validation_error", ValidationError)
Error("conflict", ConflictError)
```

### 2. Include Actionable Information

```go
// ❌ Bad - No context
return nil, products.MakeNotFound(fmt.Errorf("not found"))

// ✅ Good - Actionable information
return nil, &products.NotFoundError{
    Code:         "NOT_FOUND",
    Message:      "Product with SKU 'ABC-123' does not exist",
    ResourceType: "product",
    ResourceID:   "ABC-123",
    RequestID:    requestID,
}
```

### 3. Use Appropriate Status Codes

```go
// ❌ Bad - Wrong status code
Response("not_found", StatusBadRequest)  // 400 for not found?

// ✅ Good - Correct status codes
Response("not_found", StatusNotFound)           // 404
Response("validation_error", StatusBadRequest)  // 400
Response("conflict", StatusConflict)            // 409
```

### 4. Mark Retryable Errors

```go
// ❌ Bad - No retry information
Error("timeout")  // Client doesn't know if they can retry

// ✅ Good - Clear retry hints
func (s *svc) Call(...) error {
    // Timeout - retryable
    return MakeTimeout(err)  // temporary=true, timeout=true
    
    // Rate limit - retryable with delay
    return &RateLimitError{RetryAfter: 60}
    
    // Validation - not retryable
    return &ValidationError{...}  // no temporary flag
}
```

### 5. Don't Expose Internal Details

```go
// ❌ Bad - Exposes internal details
return nil, products.MakeInternal(
    fmt.Errorf("postgres: connection refused to db-primary:5432"),
)

// ✅ Good - Safe error message
return nil, &products.InternalError{
    Code:      "INTERNAL_ERROR",
    Message:   "An internal error occurred",
    RequestID: reqID,  // For support/debugging
}
```

### 6. Log Internal Errors Server-Side

```go
func (s *svc) Create(ctx context.Context, p *Payload) (*Result, error) {
    result, err := s.store.Insert(ctx, p)
    if err != nil {
        // Log full error server-side
        s.logger.Error("database insert failed",
            "error", err,
            "request_id", requestID(ctx),
            "payload", p,
        )
        
        // Return sanitized error to client
        return nil, &InternalError{
            Code:      "INTERNAL_ERROR",
            Message:   "Failed to create resource",
            RequestID: requestID(ctx),
        }
    }
    return result, nil
}
```

---

## 📝 Summary

### Error Definition
- **Service-level**: Define once, use across all methods
- **Method-level**: Specific to individual methods
- **Custom types**: Use Type() for structured error responses

### Custom Error Types
- Define clear, specific error types for each error category
- Include relevant context (resource type, ID, timestamps)
- Use consistent structure across your API

### HTTP Mapping
- Use appropriate status codes (4xx for client, 5xx for server)
- Include response headers for rate limits (Retry-After)
- Custom body structure per error type

### gRPC Mapping
- Use gRPC-specific status codes
- Map to equivalent HTTP codes where applicable
- Include error details in messages

### Implementation Patterns
- Use generated `Make<Error>` functions for simple errors
- Return custom error type instances for rich errors
- Include request ID for traceability
- Mark retryable errors with Temporary/Timeout flags

---

## 📋 Knowledge Check

Before proceeding to Part 5 (Security), ensure you can:

- [ ] Define service-level and method-level errors
- [ ] Create custom error types with relevant fields
- [ ] Map errors to appropriate HTTP status codes
- [ ] Map errors to gRPC status codes
- [ ] Return structured error responses
- [ ] Implement validation error handling with field details
- [ ] Handle business logic errors with domain context
- [ ] Include retry hints for transient errors
- [ ] Log internal errors safely without exposing details

---

## 🔗 Quick Reference Links

- [Goa Error DSL](https://pkg.go.dev/goa.design/goa/v3/dsl#Error)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [gRPC Status Codes](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)
- [Google Cloud Error Model](https://cloud.google.com/apis/design/errors)
- [RFC 7807 Problem Details](https://tools.ietf.org/html/rfc7807)

---

> **Next Up:** Part 5 - Security & Authentication (JWT, API Keys, OAuth2)
