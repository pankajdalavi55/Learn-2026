# Part 3: Transport Layer (HTTP & gRPC)

> **Goal:** Master Goa's transport layer - HTTP routing, request/response mapping, and gRPC integration

---

## 📚 Table of Contents

1. [Transport Layer Overview](#transport-layer-overview)
2. [HTTP Transport](#http-transport)
   - [HTTP DSL Basics](#http-dsl-basics)
   - [Routing](#routing)
   - [Path Parameters](#path-parameters)
   - [Query Parameters](#query-parameters)
   - [Headers](#headers)
   - [Request/Response Mapping](#requestresponse-mapping)
   - [Status Codes](#status-codes)
   - [Content Types](#content-types)
   - [File Uploads](#file-uploads)
3. [gRPC Transport](#grpc-transport)
   - [Proto Generation](#proto-generation)
   - [gRPC Mapping](#grpc-mapping)
   - [Streaming](#streaming)
   - [Using Generated Clients](#using-generated-clients)
4. [Multi-Transport Services](#multi-transport-services)
5. [Complete Examples](#complete-examples)
6. [Summary](#summary)
7. [Knowledge Check](#knowledge-check)

---

## 🎯 Transport Layer Overview

### What is the Transport Layer?

The transport layer in Goa separates **what** your API does (service layer) from **how** it communicates (HTTP, gRPC). This separation provides:

```
┌─────────────────────────────────────────────────────────────────┐
│                       CLIENT REQUEST                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    TRANSPORT LAYER                              │
│  ┌─────────────────────┐      ┌─────────────────────┐          │
│  │    HTTP Handler     │      │    gRPC Handler     │          │
│  │  - Parse request    │      │  - Decode protobuf  │          │
│  │  - Validate input   │      │  - Validate input   │          │
│  │  - Map to payload   │      │  - Map to payload   │          │
│  └─────────────────────┘      └─────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    SERVICE LAYER                                │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Your Business Logic                         │   │
│  │         (Same code for HTTP and gRPC)                    │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    TRANSPORT LAYER                              │
│  ┌─────────────────────┐      ┌─────────────────────┐          │
│  │    HTTP Encoder     │      │    gRPC Encoder     │          │
│  │  - Map from result  │      │  - Encode protobuf  │          │
│  │  - Set status code  │      │  - Set status code  │          │
│  │  - Write JSON/XML   │      │  - Stream response  │          │
│  └─────────────────────┘      └─────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       CLIENT RESPONSE                           │
└─────────────────────────────────────────────────────────────────┘
```

### Why Separate Transport from Service?

| Benefit | Explanation |
|---------|-------------|
| **Protocol Independence** | Same business logic serves HTTP and gRPC clients |
| **Clean Architecture** | Service layer knows nothing about HTTP codes or headers |
| **Easy Testing** | Test business logic without HTTP setup |
| **Flexible Deployment** | Add/remove transports without changing core code |
| **Consistent Validation** | Same validation rules regardless of transport |

### Transport Layer in Goa DSL

```go
// design/design.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// Service definition (transport-agnostic)
var _ = Service("users", func() {
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int, "User ID")
            Required("id")
        })
        Result(User)
        Error("not_found")
        
        // HTTP transport mapping
        HTTP(func() {
            GET("/users/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
        
        // gRPC transport mapping
        GRPC(func() {
            Response(CodeOK)
            Response("not_found", CodeNotFound)
        })
    })
})
```

---

## 🌐 HTTP Transport

### HTTP DSL Basics

The `HTTP` function defines how a service or method maps to HTTP. It can appear at both service and method levels.

#### Service-Level HTTP

```go
var _ = Service("users", func() {
    Description("User management service")
    
    // Service-level HTTP settings
    HTTP(func() {
        // Base path for all methods in this service
        Path("/api/v1/users")
        
        // Common headers for all methods
        Header("X-Request-ID:request_id", String, "Request tracking ID")
        
        // Common query parameters
        Param("version:api_version", String, "API version")
    })
    
    Method("list", func() {
        // This method will be at /api/v1/users
        HTTP(func() {
            GET("")  // Empty = use service base path
        })
    })
    
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int)
            Required("id")
        })
        Result(User)
        
        // This method will be at /api/v1/users/{id}
        HTTP(func() {
            GET("/{id}")
        })
    })
})
```

#### Generated HTTP Code Structure

```
gen/
├── http/
│   └── users/
│       ├── client/          # HTTP client code
│       │   ├── client.go    # Client struct with methods
│       │   ├── encode_decode.go  # Request/response encoding
│       │   ├── paths.go     # URL path construction
│       │   └── types.go     # HTTP-specific types
│       └── server/          # HTTP server code
│           ├── server.go    # Server struct and mounts
│           ├── encode_decode.go  # Request/response handling
│           ├── paths.go     # Path patterns
│           └── types.go     # HTTP-specific types
└── users/
    ├── service.go           # Service interface
    └── endpoints.go         # Transport-agnostic endpoints
```

### Routing

#### HTTP Methods

```go
Method("create", func() {
    Payload(CreatePayload)
    Result(User)
    
    HTTP(func() {
        POST("/users")      // Create resource
    })
})

Method("list", func() {
    Result(ArrayOf(User))
    
    HTTP(func() {
        GET("/users")       // List resources
    })
})

Method("get", func() {
    Payload(func() {
        Attribute("id", Int)
    })
    Result(User)
    
    HTTP(func() {
        GET("/users/{id}")  // Get single resource
    })
})

Method("update", func() {
    Payload(UpdatePayload)
    Result(User)
    
    HTTP(func() {
        PUT("/users/{id}")  // Full update
    })
})

Method("patch", func() {
    Payload(PatchPayload)
    Result(User)
    
    HTTP(func() {
        PATCH("/users/{id}")  // Partial update
    })
})

Method("delete", func() {
    Payload(func() {
        Attribute("id", Int)
    })
    
    HTTP(func() {
        DELETE("/users/{id}")  // Delete resource
    })
})

Method("options", func() {
    HTTP(func() {
        OPTIONS("/users")  // CORS preflight
    })
})

Method("head", func() {
    HTTP(func() {
        HEAD("/users/{id}")  // Check existence
    })
})
```

#### Route Patterns

```
┌─────────────────────────────────────────────────────────────────┐
│                    ROUTE PATTERN EXAMPLES                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  /users                    → Collection endpoint                │
│  /users/{id}               → Single resource                    │
│  /users/{id}/posts         → Nested resource                    │
│  /users/{id}/posts/{pid}   → Deeply nested                      │
│  /api/v1/users             → Versioned API                      │
│  /search                   → Action endpoint                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### RESTful Routing Patterns

```go
// Complete RESTful resource definition
var _ = Service("posts", func() {
    HTTP(func() {
        Path("/api/v1")  // Base path
    })
    
    // POST /api/v1/posts
    Method("create", func() {
        Payload(CreatePostPayload)
        Result(Post)
        HTTP(func() {
            POST("/posts")
            Response(StatusCreated)
        })
    })
    
    // GET /api/v1/posts
    Method("list", func() {
        Payload(ListPayload)
        Result(PostList)
        HTTP(func() {
            GET("/posts")
        })
    })
    
    // GET /api/v1/posts/{id}
    Method("show", func() {
        Payload(func() {
            Attribute("id", String, "Post ID")
            Required("id")
        })
        Result(Post)
        Error("not_found")
        HTTP(func() {
            GET("/posts/{id}")
            Response("not_found", StatusNotFound)
        })
    })
    
    // PUT /api/v1/posts/{id}
    Method("update", func() {
        Payload(UpdatePostPayload)
        Result(Post)
        Error("not_found")
        HTTP(func() {
            PUT("/posts/{id}")
            Response("not_found", StatusNotFound)
        })
    })
    
    // DELETE /api/v1/posts/{id}
    Method("delete", func() {
        Payload(func() {
            Attribute("id", String)
            Required("id")
        })
        Error("not_found")
        HTTP(func() {
            DELETE("/posts/{id}")
            Response(StatusNoContent)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Nested resources: GET /api/v1/posts/{id}/comments
    Method("list_comments", func() {
        Payload(func() {
            Attribute("id", String, "Post ID")
            Required("id")
        })
        Result(ArrayOf(Comment))
        HTTP(func() {
            GET("/posts/{id}/comments")
        })
    })
})
```

### Path Parameters

Path parameters capture values from the URL path itself.

#### Basic Path Parameters

```go
Method("get", func() {
    Payload(func() {
        // Define the attribute
        Attribute("user_id", Int, "User identifier")
        Required("user_id")
    })
    Result(User)
    
    HTTP(func() {
        // {user_id} references the payload attribute
        GET("/users/{user_id}")
    })
})
```

#### Path Parameter Mapping

```go
Method("get", func() {
    Payload(func() {
        // Payload uses Go-style naming
        Attribute("userID", Int, "User identifier")
        Required("userID")
    })
    Result(User)
    
    HTTP(func() {
        // Map URL param name to payload attribute
        // Format: {url_name:payload_attribute}
        GET("/users/{user_id:userID}")
    })
})
```

#### Multiple Path Parameters

```go
Method("get_comment", func() {
    Payload(func() {
        Attribute("post_id", String, "Post identifier")
        Attribute("comment_id", String, "Comment identifier")
        Required("post_id", "comment_id")
    })
    Result(Comment)
    
    HTTP(func() {
        GET("/posts/{post_id}/comments/{comment_id}")
    })
})

// Example: GET /posts/abc123/comments/xyz789
// post_id = "abc123", comment_id = "xyz789"
```

#### Path Parameter Types

```go
Method("examples", func() {
    Payload(func() {
        // String (default)
        Attribute("slug", String)
        
        // Integer - automatically parsed
        Attribute("id", Int)
        
        // UUID with validation
        Attribute("uuid", String, func() {
            Format(FormatUUID)
        })
        
        // Enum in path
        Attribute("status", String, func() {
            Enum("active", "inactive", "pending")
        })
    })
    
    HTTP(func() {
        GET("/items/{status}/{id}/{slug}/{uuid}")
    })
})
```

#### Path Parameter Validation Flow

```
┌─────────────────────────────────────────────────────────────────┐
│              PATH PARAMETER VALIDATION FLOW                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Request: GET /users/abc                                        │
│                     │                                           │
│                     ▼                                           │
│  ┌─────────────────────────────────────────────┐               │
│  │  1. Extract from path: "abc"                │               │
│  └─────────────────────────────────────────────┘               │
│                     │                                           │
│                     ▼                                           │
│  ┌─────────────────────────────────────────────┐               │
│  │  2. Type conversion: String → Int           │               │
│  │     ERROR: "abc" is not a valid integer     │───────┐       │
│  └─────────────────────────────────────────────┘       │       │
│                     │ (if success)                     │       │
│                     ▼                                  │       │
│  ┌─────────────────────────────────────────────┐       │       │
│  │  3. Validation: Minimum(1)                  │       │       │
│  │     ERROR: 0 is less than minimum           │───────┤       │
│  └─────────────────────────────────────────────┘       │       │
│                     │ (if success)                     │       │
│                     ▼                                  │       │
│  ┌─────────────────────────────────────────────┐       │       │
│  │  4. Assign to Payload struct                │       ▼       │
│  └─────────────────────────────────────────────┘  400 Bad      │
│                     │                             Request       │
│                     ▼                                           │
│           Call Service Method                                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Query Parameters

Query parameters capture values from the URL query string (`?key=value`).

#### Basic Query Parameters

```go
Method("list", func() {
    Payload(func() {
        Attribute("page", Int, "Page number", func() {
            Default(1)
            Minimum(1)
        })
        Attribute("limit", Int, "Items per page", func() {
            Default(20)
            Minimum(1)
            Maximum(100)
        })
        Attribute("sort", String, "Sort field", func() {
            Default("created_at")
            Enum("created_at", "updated_at", "name")
        })
        Attribute("order", String, "Sort order", func() {
            Default("desc")
            Enum("asc", "desc")
        })
    })
    Result(UserList)
    
    HTTP(func() {
        GET("/users")
        // All payload attributes become query params
        // GET /users?page=2&limit=50&sort=name&order=asc
    })
})
```

#### Explicit Query Parameter Mapping

```go
Method("search", func() {
    Payload(func() {
        Attribute("query", String, "Search query")
        Attribute("category", String, "Filter by category")
        Attribute("min_price", Float64, "Minimum price")
        Attribute("max_price", Float64, "Maximum price")
        Required("query")
    })
    Result(SearchResults)
    
    HTTP(func() {
        GET("/search")
        
        // Explicit mapping with Param()
        // Format: "url_name:payload_attribute"
        Param("q:query")           // ?q=value maps to query
        Param("cat:category")      // ?cat=value maps to category
        Param("min:min_price")     // ?min=value maps to min_price
        Param("max:max_price")     // ?max=value maps to max_price
    })
})

// Example: GET /search?q=shoes&cat=footwear&min=50&max=200
```

#### Array Query Parameters

```go
Method("filter", func() {
    Payload(func() {
        // Array of strings
        Attribute("tags", ArrayOf(String), "Filter by tags")
        
        // Array of integers
        Attribute("ids", ArrayOf(Int), "Filter by IDs")
        
        // Array with validation
        Attribute("statuses", ArrayOf(String), func() {
            Elem(func() {
                Enum("active", "pending", "archived")
            })
        })
    })
    Result(ArrayOf(Item))
    
    HTTP(func() {
        GET("/items")
        // Arrays are passed as repeated params:
        // GET /items?tags=go&tags=api&ids=1&ids=2&ids=3
    })
})
```

#### Mixing Path and Query Parameters

```go
Method("list_user_posts", func() {
    Payload(func() {
        // Path parameter
        Attribute("user_id", Int, "User ID")
        
        // Query parameters
        Attribute("page", Int, func() { Default(1) })
        Attribute("limit", Int, func() { Default(10) })
        Attribute("status", String, func() {
            Enum("draft", "published", "archived")
        })
        
        Required("user_id")
    })
    Result(PostList)
    
    HTTP(func() {
        GET("/users/{user_id}/posts")
        // user_id comes from path
        // page, limit, status come from query
        
        // Example: GET /users/123/posts?page=2&status=published
    })
})
```

#### Boolean Query Parameters

```go
Method("list", func() {
    Payload(func() {
        Attribute("active", Boolean, "Filter active only")
        Attribute("verified", Boolean, "Filter verified only")
    })
    Result(ArrayOf(User))
    
    HTTP(func() {
        GET("/users")
        // Boolean params accept: true, false, 1, 0
        // GET /users?active=true&verified=1
    })
})
```

### Headers

Headers allow passing metadata with HTTP requests and responses.

#### Request Headers

```go
Method("create", func() {
    Payload(func() {
        // Header attributes
        Attribute("auth_token", String, "Authentication token")
        Attribute("request_id", String, "Request tracking ID")
        Attribute("content_lang", String, "Content language")
        
        // Body attributes
        Attribute("name", String)
        Attribute("email", String)
        
        Required("auth_token", "name", "email")
    })
    Result(User)
    
    HTTP(func() {
        POST("/users")
        
        // Map payload attributes to headers
        // Format: "Header-Name:payload_attribute"
        Header("Authorization:auth_token")
        Header("X-Request-ID:request_id")
        Header("Accept-Language:content_lang")
        
        // name and email go to body by default
    })
})
```

#### Response Headers

```go
Method("create", func() {
    Payload(CreatePayload)
    Result(func() {
        Attribute("user", User)
        Attribute("location", String, "Resource URL")
        Attribute("request_id", String, "Request tracking ID")
        Required("user", "location")
    })
    
    HTTP(func() {
        POST("/users")
        Response(StatusCreated, func() {
            // Map result attributes to response headers
            Header("Location:location")
            Header("X-Request-ID:request_id")
            // "user" goes to response body
        })
    })
})
```

#### Common Header Patterns

```go
// Service-level headers (apply to all methods)
var _ = Service("api", func() {
    HTTP(func() {
        Path("/api/v1")
        
        // Common request headers
        Header("X-API-Key:api_key", String, "API key for authentication")
        Header("X-Request-ID:request_id", String, "Request tracking")
    })
    
    Method("any_method", func() {
        Payload(func() {
            // These are inherited from service
            Attribute("api_key", String)
            Attribute("request_id", String)
            
            // Method-specific attributes
            Attribute("data", String)
            
            Required("api_key")
        })
        // ...
    })
})
```

#### Header Validation

```go
Method("upload", func() {
    Payload(func() {
        Attribute("content_type", String, func() {
            Enum("image/jpeg", "image/png", "image/gif")
        })
        Attribute("content_length", Int64, func() {
            Minimum(1)
            Maximum(10485760)  // 10MB
        })
        Required("content_type", "content_length")
    })
    
    HTTP(func() {
        POST("/upload")
        Header("Content-Type:content_type")
        Header("Content-Length:content_length")
    })
})
```

#### Header Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    REQUEST HEADER FLOW                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  HTTP Request:                                                  │
│  POST /users                                                    │
│  Authorization: Bearer xyz123                                   │
│  X-Request-ID: req-456                                          │
│  Content-Type: application/json                                 │
│                                                                 │
│  {"name": "John", "email": "john@example.com"}                  │
│                                                                 │
│                         │                                       │
│                         ▼                                       │
│  ┌─────────────────────────────────────────────┐               │
│  │           HTTP Decoder                       │               │
│  │  - Extract "Authorization" → auth_token      │               │
│  │  - Extract "X-Request-ID" → request_id       │               │
│  │  - Parse JSON body → name, email             │               │
│  └─────────────────────────────────────────────┘               │
│                         │                                       │
│                         ▼                                       │
│  ┌─────────────────────────────────────────────┐               │
│  │           Payload Struct                     │               │
│  │  {                                           │               │
│  │    AuthToken:  "Bearer xyz123",              │               │
│  │    RequestID:  "req-456",                    │               │
│  │    Name:       "John",                       │               │
│  │    Email:      "john@example.com"            │               │
│  │  }                                           │               │
│  └─────────────────────────────────────────────┘               │
│                         │                                       │
│                         ▼                                       │
│                Service Method                                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Request/Response Mapping

#### Complete Mapping Example

```go
Method("create_order", func() {
    Payload(func() {
        // Path parameters
        Attribute("customer_id", String, "Customer identifier")
        
        // Query parameters
        Attribute("notify", Boolean, "Send notification", func() {
            Default(true)
        })
        
        // Header parameters
        Attribute("idempotency_key", String, "Idempotency key")
        Attribute("auth_token", String, "Bearer token")
        
        // Body parameters
        Attribute("items", ArrayOf(OrderItem), "Order items")
        Attribute("shipping_address", Address, "Shipping address")
        Attribute("notes", String, "Order notes")
        
        Required("customer_id", "auth_token", "items", "shipping_address")
    })
    
    Result(func() {
        // Response body
        Attribute("order", Order, "Created order")
        
        // Response headers
        Attribute("order_id", String, "Order identifier")
        Attribute("location", String, "Resource URL")
        
        Required("order", "order_id", "location")
    })
    
    Error("validation_error")
    Error("unauthorized")
    
    HTTP(func() {
        POST("/customers/{customer_id}/orders")
        
        // Query parameter mapping
        Param("notify")
        
        // Request header mapping
        Header("Idempotency-Key:idempotency_key")
        Header("Authorization:auth_token")
        
        // Body gets: items, shipping_address, notes (by default)
        
        Response(StatusCreated, func() {
            // Response header mapping
            Header("X-Order-ID:order_id")
            Header("Location:location")
            // Body gets: order (remaining attribute)
        })
        
        Response("validation_error", StatusBadRequest)
        Response("unauthorized", StatusUnauthorized)
    })
})
```

#### Explicit Body Mapping

```go
Method("update", func() {
    Payload(func() {
        Attribute("id", Int)
        Attribute("name", String)
        Attribute("email", String)
        Attribute("auth", String)
        Required("id", "name", "email", "auth")
    })
    Result(User)
    
    HTTP(func() {
        PUT("/users/{id}")
        Header("Authorization:auth")
        
        // Explicitly specify which attributes go to body
        Body(func() {
            Attribute("name")
            Attribute("email")
        })
    })
})
```

#### Body as Single Attribute

```go
// When you want the entire body to be a single field
Method("raw_upload", func() {
    Payload(func() {
        Attribute("id", String)
        Attribute("data", Bytes, "Raw file data")
        Required("id", "data")
    })
    
    HTTP(func() {
        POST("/upload/{id}")
        // The entire request body is the "data" field
        Body("data")
    })
})
```

#### Mapping Summary Table

| DSL Element | Payload Source | Example |
|-------------|---------------|---------|
| `GET("/path/{id}")` | URL path | `/users/123` → id=123 |
| `Param("name")` | Query string | `?name=john` → name="john" |
| `Header("X-Key:key")` | HTTP header | `X-Key: abc` → key="abc" |
| `Body(func(){...})` | Request body | JSON fields |
| `Body("field")` | Entire body | Raw bytes/string |

### Status Codes

Goa provides constants for all standard HTTP status codes.

#### Success Status Codes

```go
Method("example", func() {
    HTTP(func() {
        POST("/resource")
        
        // 200 OK - Default for successful responses
        Response(StatusOK)
        
        // 201 Created - Resource created
        Response(StatusCreated)
        
        // 202 Accepted - Request accepted for processing
        Response(StatusAccepted)
        
        // 204 No Content - Success with no body
        Response(StatusNoContent)
    })
})
```

#### Error Status Codes

```go
var _ = Service("users", func() {
    Error("not_found", ErrorResult, "User not found")
    Error("bad_request", ErrorResult, "Invalid input")
    Error("unauthorized", ErrorResult, "Authentication required")
    Error("forbidden", ErrorResult, "Access denied")
    Error("conflict", ErrorResult, "Resource conflict")
    Error("gone", ErrorResult, "Resource no longer available")
    Error("unprocessable", ErrorResult, "Validation failed")
    Error("too_many", ErrorResult, "Rate limit exceeded")
    Error("internal", ErrorResult, "Internal server error")
    Error("unavailable", ErrorResult, "Service unavailable")
    
    Method("create", func() {
        Payload(CreatePayload)
        Result(User)
        
        // Reference service-level errors
        Error("not_found")
        Error("bad_request")
        Error("unauthorized")
        Error("forbidden")
        Error("conflict")
        Error("unprocessable")
        Error("too_many")
        Error("internal")
        Error("unavailable")
        
        HTTP(func() {
            POST("/users")
            
            // Map errors to HTTP status codes
            Response(StatusCreated)
            Response("bad_request", StatusBadRequest)           // 400
            Response("unauthorized", StatusUnauthorized)        // 401
            Response("forbidden", StatusForbidden)              // 403
            Response("not_found", StatusNotFound)               // 404
            Response("conflict", StatusConflict)                // 409
            Response("gone", StatusGone)                        // 410
            Response("unprocessable", StatusUnprocessableEntity)// 422
            Response("too_many", StatusTooManyRequests)         // 429
            Response("internal", StatusInternalServerError)     // 500
            Response("unavailable", StatusServiceUnavailable)   // 503
        })
    })
})
```

#### HTTP Status Code Reference

```
┌─────────────────────────────────────────────────────────────────┐
│                    HTTP STATUS CODE GUIDE                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  2xx SUCCESS                                                    │
│  ───────────                                                    │
│  200 OK              - Standard success response                │
│  201 Created         - Resource created (POST)                  │
│  202 Accepted        - Request accepted, processing async       │
│  204 No Content      - Success, no body (DELETE)                │
│                                                                 │
│  3xx REDIRECT                                                   │
│  ────────────                                                   │
│  301 Moved Permanently  - Resource moved permanently            │
│  302 Found              - Temporary redirect                    │
│  304 Not Modified       - Cached version is still valid         │
│                                                                 │
│  4xx CLIENT ERROR                                               │
│  ────────────────                                               │
│  400 Bad Request        - Invalid request syntax/data           │
│  401 Unauthorized       - Authentication required               │
│  403 Forbidden          - Authenticated but not allowed         │
│  404 Not Found          - Resource doesn't exist                │
│  405 Method Not Allowed - HTTP method not supported             │
│  409 Conflict           - Request conflicts with state          │
│  410 Gone               - Resource permanently deleted          │
│  422 Unprocessable      - Validation errors                     │
│  429 Too Many Requests  - Rate limit exceeded                   │
│                                                                 │
│  5xx SERVER ERROR                                               │
│  ────────────────                                               │
│  500 Internal Server Error  - Unexpected server error           │
│  502 Bad Gateway            - Invalid upstream response         │
│  503 Service Unavailable    - Server temporarily down           │
│  504 Gateway Timeout        - Upstream timeout                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Conditional Response Bodies

```go
Method("get", func() {
    Payload(func() {
        Attribute("id", Int)
        Required("id")
    })
    
    // Different result types for different scenarios
    Result(User)
    Error("not_found", NotFoundError)
    Error("deleted", DeletedError)
    
    HTTP(func() {
        GET("/users/{id}")
        
        Response(StatusOK)
        
        Response("not_found", StatusNotFound, func() {
            // Custom body for 404
            Body(func() {
                Attribute("message")
                Attribute("resource_type")
            })
        })
        
        Response("deleted", StatusGone, func() {
            // Custom body for 410
            Body(func() {
                Attribute("deleted_at")
                Attribute("reason")
            })
        })
    })
})
```

### Content Types

#### JSON (Default)

```go
// JSON is the default content type
Method("create", func() {
    Payload(CreatePayload)
    Result(User)
    
    HTTP(func() {
        POST("/users")
        // Automatically uses application/json
    })
})
```

#### Custom Content Types

```go
Method("upload_image", func() {
    Payload(func() {
        Attribute("id", String)
        Attribute("image", Bytes)
        Required("id", "image")
    })
    
    HTTP(func() {
        POST("/images/{id}")
        Body("image")
        
        // Accept specific content types
        ContentType("image/png")
        ContentType("image/jpeg")
    })
})
```

#### Multiple Content Types

```go
Method("export", func() {
    Payload(func() {
        Attribute("format", String, func() {
            Enum("json", "xml", "csv")
            Default("json")
        })
    })
    Result(Bytes)
    
    HTTP(func() {
        GET("/export")
        
        Response(StatusOK, func() {
            // Content type set dynamically based on format
            ContentType("application/json")
            ContentType("application/xml")
            ContentType("text/csv")
        })
    })
})
```

### File Uploads

#### Simple File Upload

```go
Method("upload", func() {
    Payload(func() {
        Attribute("file", Bytes, "File content")
        Attribute("filename", String, "Original filename")
        Required("file")
    })
    Result(UploadResult)
    
    HTTP(func() {
        POST("/upload")
        Header("X-Filename:filename")
        Body("file")
        ContentType("application/octet-stream")
    })
})
```

#### Multipart Form Upload

```go
Method("upload_multipart", func() {
    Payload(func() {
        Attribute("file", Bytes, "File content")
        Attribute("description", String, "File description")
        Attribute("tags", ArrayOf(String), "File tags")
    })
    Result(UploadResult)
    
    HTTP(func() {
        POST("/upload")
        MultipartRequest()
    })
})
```

---

## 📡 gRPC Transport

### Proto Generation

Goa automatically generates Protocol Buffer definitions from your DSL.

#### Enabling gRPC

```go
var _ = Service("calculator", func() {
    Method("add", func() {
        Payload(func() {
            Attribute("a", Int)
            Attribute("b", Int)
            Required("a", "b")
        })
        Result(func() {
            Attribute("result", Int)
            Required("result")
        })
        
        // Add gRPC transport
        GRPC(func() {
            Response(CodeOK)
        })
    })
})
```

#### Generated Proto File

```protobuf
// gen/grpc/calculator/pb/calculator.proto
syntax = "proto3";

package calculator;

option go_package = "/calculator";

service Calculator {
    rpc Add(AddRequest) returns (AddResponse);
}

message AddRequest {
    int64 a = 1;
    int64 b = 2;
}

message AddResponse {
    int64 result = 1;
}
```

#### Proto Generation Command

```bash
# Generate code including proto files
goa gen myproject/design

# Generated structure
gen/
├── grpc/
│   └── calculator/
│       ├── pb/
│       │   ├── calculator.proto    # Protocol Buffer definition
│       │   ├── calculator.pb.go    # Generated Go code
│       │   └── calculator_grpc.pb.go  # gRPC service code
│       ├── client/
│       │   ├── client.go          # gRPC client
│       │   ├── encode_decode.go   # Message encoding
│       │   └── types.go           # Client types
│       └── server/
│           ├── server.go          # gRPC server
│           ├── encode_decode.go   # Message encoding
│           └── types.go           # Server types
└── calculator/
    └── service.go                 # Service interface (same for HTTP & gRPC)
```

#### Proto Field Numbers

```go
// Goa assigns field numbers automatically, but you can control them
Method("create", func() {
    Payload(func() {
        Attribute("id", String, func() {
            // Control proto field number with Meta
            Meta("rpc:tag", "1")
        })
        Attribute("name", String, func() {
            Meta("rpc:tag", "2")
        })
        Attribute("email", String, func() {
            Meta("rpc:tag", "3")
        })
    })
})
```

### gRPC Mapping

#### gRPC Status Codes

```go
Method("get", func() {
    Payload(GetPayload)
    Result(User)
    
    Error("not_found")
    Error("invalid")
    Error("internal")
    
    GRPC(func() {
        Response(CodeOK)
        Response("not_found", CodeNotFound)
        Response("invalid", CodeInvalidArgument)
        Response("internal", CodeInternal)
    })
})
```

#### gRPC Status Code Reference

| Goa Code | gRPC Code | HTTP Equivalent | Usage |
|----------|-----------|-----------------|-------|
| `CodeOK` | 0 OK | 200 | Success |
| `CodeCanceled` | 1 CANCELED | 499 | Client cancelled |
| `CodeUnknown` | 2 UNKNOWN | 500 | Unknown error |
| `CodeInvalidArgument` | 3 INVALID_ARGUMENT | 400 | Invalid input |
| `CodeDeadlineExceeded` | 4 DEADLINE_EXCEEDED | 504 | Timeout |
| `CodeNotFound` | 5 NOT_FOUND | 404 | Resource not found |
| `CodeAlreadyExists` | 6 ALREADY_EXISTS | 409 | Resource exists |
| `CodePermissionDenied` | 7 PERMISSION_DENIED | 403 | Access denied |
| `CodeResourceExhausted` | 8 RESOURCE_EXHAUSTED | 429 | Rate limited |
| `CodeFailedPrecondition` | 9 FAILED_PRECONDITION | 400 | State error |
| `CodeAborted` | 10 ABORTED | 409 | Concurrent conflict |
| `CodeOutOfRange` | 11 OUT_OF_RANGE | 400 | Value out of range |
| `CodeUnimplemented` | 12 UNIMPLEMENTED | 501 | Not implemented |
| `CodeInternal` | 13 INTERNAL | 500 | Internal error |
| `CodeUnavailable` | 14 UNAVAILABLE | 503 | Service unavailable |
| `CodeDataLoss` | 15 DATA_LOSS | 500 | Data corruption |
| `CodeUnauthenticated` | 16 UNAUTHENTICATED | 401 | Not authenticated |

#### gRPC Metadata (Headers)

```go
Method("create", func() {
    Payload(func() {
        Attribute("auth_token", String, "Authentication token")
        Attribute("request_id", String, "Request tracking ID")
        Attribute("data", CreateData)
        Required("auth_token", "data")
    })
    Result(User)
    
    GRPC(func() {
        // Map attributes to gRPC metadata (headers)
        Metadata(func() {
            Attribute("auth_token:authorization")  // metadata key: authorization
            Attribute("request_id:x-request-id")
        })
        
        // Message gets "data" field
        Message(func() {
            Attribute("data")
        })
        
        Response(CodeOK)
    })
})
```

#### Complete gRPC Mapping Example

```go
var _ = Service("users", func() {
    Error("not_found", ErrorResult, "User not found")
    Error("unauthorized", ErrorResult, "Not authenticated")
    Error("forbidden", ErrorResult, "Access denied")
    Error("invalid", ErrorResult, "Invalid input")
    
    Method("create", func() {
        Payload(func() {
            // Metadata (headers)
            Attribute("auth_token", String)
            Attribute("idempotency_key", String)
            
            // Message (body)
            Attribute("name", String)
            Attribute("email", String)
            
            Required("auth_token", "name", "email")
        })
        
        Result(func() {
            Attribute("user", User)
            Attribute("request_id", String)
            Required("user")
        })
        
        Error("not_found")
        Error("unauthorized")
        Error("forbidden")
        Error("invalid")
        
        GRPC(func() {
            Metadata(func() {
                // Request metadata
                Attribute("auth_token:authorization")
                Attribute("idempotency_key:idempotency-key")
            })
            
            Message(func() {
                // Request message fields
                Attribute("name")
                Attribute("email")
            })
            
            Response(CodeOK, func() {
                // Response metadata
                Metadata(func() {
                    Attribute("request_id:x-request-id")
                })
                // Response message gets "user"
            })
            
            Response("not_found", CodeNotFound)
            Response("unauthorized", CodeUnauthenticated)
            Response("forbidden", CodePermissionDenied)
            Response("invalid", CodeInvalidArgument)
        })
    })
})
```

### Streaming

gRPC supports four types of streaming:

```
┌─────────────────────────────────────────────────────────────────┐
│                    gRPC STREAMING TYPES                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  UNARY (default)                                                │
│  Client ──── Single Request ────▶ Server                        │
│  Client ◀─── Single Response ─── Server                         │
│                                                                 │
│  SERVER STREAMING                                               │
│  Client ──── Single Request ────▶ Server                        │
│  Client ◀─── Stream of Responses ─ Server                       │
│                                                                 │
│  CLIENT STREAMING                                               │
│  Client ──── Stream of Requests ─▶ Server                       │
│  Client ◀─── Single Response ──── Server                        │
│                                                                 │
│  BIDIRECTIONAL STREAMING                                        │
│  Client ◀─── Stream ────▶ Server                                │
│  (both send multiple messages)                                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Server Streaming

```go
Method("watch", func() {
    Description("Watch for user changes")
    
    Payload(func() {
        Attribute("user_id", Int, "User to watch")
        Required("user_id")
    })
    
    // StreamingResult indicates server streaming
    StreamingResult(UserEvent)
    
    GRPC(func() {
        Response(CodeOK)
    })
})
```

Generated service interface:

```go
// gen/users/service.go
type Service interface {
    // Watch watches for user changes
    Watch(ctx context.Context, p *WatchPayload, stream WatchServerStream) error
}

// WatchServerStream is the interface for server streaming
type WatchServerStream interface {
    // Send streams an instance of UserEvent
    Send(*UserEvent) error
    // Close closes the stream
    Close() error
}
```

Implementation:

```go
// users.go
func (s *usersSvc) Watch(ctx context.Context, p *users.WatchPayload, stream users.WatchServerStream) error {
    defer stream.Close()
    
    // Watch for changes and stream them
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            event := &users.UserEvent{
                UserID:    p.UserID,
                EventType: "heartbeat",
                Timestamp: time.Now().Unix(),
            }
            if err := stream.Send(event); err != nil {
                return err
            }
        }
    }
}
```

#### Client Streaming

```go
Method("upload_chunks", func() {
    Description("Upload file in chunks")
    
    // StreamingPayload indicates client streaming
    StreamingPayload(FileChunk)
    
    Result(UploadResult)
    
    GRPC(func() {
        Response(CodeOK)
    })
})
```

Generated service interface:

```go
type Service interface {
    // UploadChunks uploads file in chunks
    UploadChunks(ctx context.Context, stream UploadChunksServerStream) error
}

type UploadChunksServerStream interface {
    // Recv receives a FileChunk
    Recv() (*FileChunk, error)
    // SendAndClose sends the result and closes the stream
    SendAndClose(*UploadResult) error
}
```

Implementation:

```go
func (s *uploadSvc) UploadChunks(ctx context.Context, stream upload.UploadChunksServerStream) error {
    var totalSize int64
    var chunks [][]byte
    
    for {
        chunk, err := stream.Recv()
        if err == io.EOF {
            // All chunks received
            break
        }
        if err != nil {
            return err
        }
        
        chunks = append(chunks, chunk.Data)
        totalSize += int64(len(chunk.Data))
    }
    
    // Process complete file...
    
    return stream.SendAndClose(&upload.UploadResult{
        FileID:    uuid.New().String(),
        TotalSize: totalSize,
        Chunks:    int32(len(chunks)),
    })
}
```

#### Bidirectional Streaming

```go
Method("chat", func() {
    Description("Real-time chat")
    
    // Both streaming indicates bidirectional
    StreamingPayload(ChatMessage)
    StreamingResult(ChatMessage)
    
    GRPC(func() {
        Response(CodeOK)
    })
})
```

Generated service interface:

```go
type Service interface {
    // Chat handles bidirectional chat messages
    Chat(ctx context.Context, stream ChatServerStream) error
}

type ChatServerStream interface {
    // Send streams a ChatMessage to client
    Send(*ChatMessage) error
    // Recv receives a ChatMessage from client
    Recv() (*ChatMessage, error)
    // Close closes the stream
    Close() error
}
```

Implementation:

```go
func (s *chatSvc) Chat(ctx context.Context, stream chat.ChatServerStream) error {
    defer stream.Close()
    
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        
        // Echo back with modification
        response := &chat.ChatMessage{
            From:      "server",
            To:        msg.From,
            Content:   "Echo: " + msg.Content,
            Timestamp: time.Now().Unix(),
        }
        
        if err := stream.Send(response); err != nil {
            return err
        }
    }
}
```

### Using Generated Clients

#### HTTP Client Usage

```go
// main.go or client code
package main

import (
    "context"
    "net/http"
    
    goahttp "goa.design/goa/v3/http"
    usersc "myproject/gen/http/users/client"
    users "myproject/gen/users"
)

func main() {
    // Create HTTP client
    httpClient := &http.Client{}
    
    // Create Goa HTTP client
    endpoint := "http://localhost:8080"
    client := usersc.NewClient(
        goahttp.NewClient(httpClient),
        endpoint,
        goahttp.RequestEncoder,
        goahttp.ResponseDecoder,
        false, // debug
    )
    
    // Create endpoints
    endpoints := users.NewClient(
        client.Get(),
        client.List(),
        client.Create(),
        client.Update(),
        client.Delete(),
    )
    
    // Use the client
    ctx := context.Background()
    
    // Get user
    result, err := endpoints.Get(ctx, &users.GetPayload{ID: 1})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", result)
    
    // Create user
    newUser, err := endpoints.Create(ctx, &users.CreatePayload{
        Name:  "John Doe",
        Email: "john@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created: %+v\n", newUser)
}
```

#### gRPC Client Usage

```go
// main.go or client code
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    userspb "myproject/gen/grpc/users/pb"
    usersc "myproject/gen/grpc/users/client"
    users "myproject/gen/users"
)

func main() {
    // Create gRPC connection
    conn, err := grpc.Dial(
        "localhost:8081",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // Create gRPC client
    grpcClient := userspb.NewUsersClient(conn)
    client := usersc.NewClient(grpcClient)
    
    // Create endpoints
    endpoints := users.NewClient(
        client.Get(),
        client.List(),
        client.Create(),
        client.Update(),
        client.Delete(),
    )
    
    // Use the client
    ctx := context.Background()
    
    // Get user
    result, err := endpoints.Get(ctx, &users.GetPayload{ID: 1})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", result)
}
```

#### Streaming Client Usage

```go
// Server streaming client
func watchUser(ctx context.Context, client users.Client, userID int) error {
    // Start streaming
    stream, err := client.Watch(ctx, &users.WatchPayload{UserID: userID})
    if err != nil {
        return err
    }
    
    // Receive events
    for {
        event, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        fmt.Printf("Event: %+v\n", event)
    }
}

// Client streaming client
func uploadFile(ctx context.Context, client upload.Client, data []byte) (*upload.UploadResult, error) {
    stream, err := client.UploadChunks(ctx)
    if err != nil {
        return nil, err
    }
    
    // Send chunks
    chunkSize := 1024 * 1024 // 1MB chunks
    for i := 0; i < len(data); i += chunkSize {
        end := i + chunkSize
        if end > len(data) {
            end = len(data)
        }
        
        if err := stream.Send(&upload.FileChunk{
            Data:     data[i:end],
            Sequence: int32(i / chunkSize),
        }); err != nil {
            return nil, err
        }
    }
    
    // Close and get result
    return stream.CloseAndRecv()
}

// Bidirectional streaming client
func chatSession(ctx context.Context, client chat.Client) error {
    stream, err := client.Chat(ctx)
    if err != nil {
        return err
    }
    
    // Send and receive concurrently
    errCh := make(chan error, 2)
    
    // Sender goroutine
    go func() {
        for i := 0; i < 10; i++ {
            if err := stream.Send(&chat.ChatMessage{
                From:    "client",
                Content: fmt.Sprintf("Message %d", i),
            }); err != nil {
                errCh <- err
                return
            }
            time.Sleep(time.Second)
        }
        stream.CloseSend()
        errCh <- nil
    }()
    
    // Receiver goroutine
    go func() {
        for {
            msg, err := stream.Recv()
            if err == io.EOF {
                errCh <- nil
                return
            }
            if err != nil {
                errCh <- err
                return
            }
            fmt.Printf("Received: %s\n", msg.Content)
        }
    }()
    
    // Wait for both to complete
    for i := 0; i < 2; i++ {
        if err := <-errCh; err != nil {
            return err
        }
    }
    return nil
}
```

---

## 🔄 Multi-Transport Services

### Dual Transport Design

```go
// design/design.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = API("multiservice", func() {
    Title("Multi-Transport Service")
    Description("Service supporting both HTTP and gRPC")
    Version("1.0")
    
    Server("main", func() {
        Host("local", func() {
            URI("http://localhost:8080")
            URI("grpc://localhost:8081")
        })
    })
})

var User = Type("User", func() {
    Attribute("id", Int, "User ID")
    Attribute("name", String, "User name")
    Attribute("email", String, "Email address")
    Required("id", "name", "email")
})

var _ = Service("users", func() {
    Description("User management service")
    
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int, "User ID")
            Required("id")
        })
        Result(User)
        Error("not_found")
        
        // HTTP transport
        HTTP(func() {
            GET("/users/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
        
        // gRPC transport
        GRPC(func() {
            Response(CodeOK)
            Response("not_found", CodeNotFound)
        })
    })
    
    Method("create", func() {
        Payload(func() {
            Attribute("name", String)
            Attribute("email", String)
            Required("name", "email")
        })
        Result(User)
        Error("invalid")
        
        HTTP(func() {
            POST("/users")
            Response(StatusCreated)
            Response("invalid", StatusBadRequest)
        })
        
        GRPC(func() {
            Response(CodeOK)
            Response("invalid", CodeInvalidArgument)
        })
    })
    
    Method("list", func() {
        Payload(func() {
            Attribute("page", Int, func() { Default(1) })
            Attribute("limit", Int, func() { Default(10) })
        })
        Result(ArrayOf(User))
        
        HTTP(func() {
            GET("/users")
            Param("page")
            Param("limit")
        })
        
        GRPC(func() {
            Response(CodeOK)
        })
    })
    
    // Streaming only available in gRPC
    Method("watch", func() {
        Payload(func() {
            Attribute("id", Int)
            Required("id")
        })
        StreamingResult(UserEvent)
        
        // No HTTP - streaming not supported
        GRPC(func() {
            Response(CodeOK)
        })
    })
})

var UserEvent = Type("UserEvent", func() {
    Attribute("user_id", Int)
    Attribute("event_type", String)
    Attribute("timestamp", Int64)
    Required("user_id", "event_type", "timestamp")
})
```

### Multi-Transport Server Setup

```go
// cmd/server/main.go
package main

import (
    "context"
    "net"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"

    "google.golang.org/grpc"
    
    goahttp "goa.design/goa/v3/http"
    goagrpc "goa.design/goa/v3/grpc"
    
    users "myproject/gen/users"
    userssvr "myproject/gen/http/users/server"
    usersgrpc "myproject/gen/grpc/users/server"
    userspb "myproject/gen/grpc/users/pb"
)

func main() {
    // Create service implementation
    usersSvc := NewUsersService()
    
    // Create endpoints
    usersEndpoints := users.NewEndpoints(usersSvc)
    
    // Setup HTTP server
    mux := goahttp.NewMuxer()
    usersServer := userssvr.New(usersEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
    userssvr.Mount(mux, usersServer)
    
    httpServer := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    // Setup gRPC server
    grpcServer := grpc.NewServer()
    usersGRPCServer := usersgrpc.New(usersEndpoints, nil)
    userspb.RegisterUsersServer(grpcServer, usersGRPCServer)
    
    // Start servers
    var wg sync.WaitGroup
    ctx, cancel := context.WithCancel(context.Background())
    
    // HTTP server goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Printf("HTTP server listening on %s", httpServer.Addr)
        if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()
    
    // gRPC server goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        lis, err := net.Listen("tcp", ":8081")
        if err != nil {
            log.Fatalf("gRPC listen error: %v", err)
        }
        log.Printf("gRPC server listening on :8081")
        if err := grpcServer.Serve(lis); err != nil {
            log.Printf("gRPC server error: %v", err)
        }
    }()
    
    // Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    log.Println("Shutting down servers...")
    cancel()
    
    // Shutdown HTTP
    httpServer.Shutdown(ctx)
    
    // Shutdown gRPC
    grpcServer.GracefulStop()
    
    wg.Wait()
    log.Println("Servers stopped")
}
```

---

## 📦 Complete Examples

### Complete HTTP Service

```go
// design/api.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = API("bookstore", func() {
    Title("Bookstore API")
    Description("A complete bookstore REST API example")
    Version("1.0")
    
    Server("bookstore", func() {
        Host("local", func() {
            URI("http://localhost:8080")
        })
        Host("production", func() {
            URI("https://api.bookstore.com")
        })
    })
    
    HTTP(func() {
        Path("/api/v1")
    })
})
```

```go
// design/types.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var Book = Type("Book", func() {
    Description("A book in the store")
    
    Attribute("id", String, "Unique identifier", func() {
        Example("book-123")
    })
    Attribute("isbn", String, "ISBN-13", func() {
        Pattern(`^\d{13}$`)
        Example("9780134190440")
    })
    Attribute("title", String, "Book title", func() {
        MinLength(1)
        MaxLength(200)
        Example("The Go Programming Language")
    })
    Attribute("author", String, "Author name", func() {
        MinLength(1)
        MaxLength(100)
        Example("Alan A. A. Donovan")
    })
    Attribute("price", Float64, "Price in USD", func() {
        Minimum(0)
        Example(49.99)
    })
    Attribute("category", String, "Book category", func() {
        Enum("fiction", "non-fiction", "technical", "children")
    })
    Attribute("stock", Int, "Available copies", func() {
        Minimum(0)
        Default(0)
    })
    Attribute("created_at", String, "Creation timestamp", func() {
        Format(FormatDateTime)
    })
    
    Required("id", "isbn", "title", "author", "price", "category")
})

var BookList = Type("BookList", func() {
    Attribute("books", ArrayOf(Book), "List of books")
    Attribute("total", Int, "Total count")
    Attribute("page", Int, "Current page")
    Attribute("limit", Int, "Page size")
    Required("books", "total", "page", "limit")
})

var CreateBookPayload = Type("CreateBookPayload", func() {
    Attribute("isbn", String, func() {
        Pattern(`^\d{13}$`)
    })
    Attribute("title", String, func() {
        MinLength(1)
        MaxLength(200)
    })
    Attribute("author", String, func() {
        MinLength(1)
        MaxLength(100)
    })
    Attribute("price", Float64, func() {
        Minimum(0)
    })
    Attribute("category", String, func() {
        Enum("fiction", "non-fiction", "technical", "children")
    })
    Attribute("stock", Int, func() {
        Minimum(0)
        Default(0)
    })
    Required("isbn", "title", "author", "price", "category")
})

var UpdateBookPayload = Type("UpdateBookPayload", func() {
    Attribute("title", String, func() {
        MinLength(1)
        MaxLength(200)
    })
    Attribute("author", String, func() {
        MinLength(1)
        MaxLength(100)
    })
    Attribute("price", Float64, func() {
        Minimum(0)
    })
    Attribute("category", String, func() {
        Enum("fiction", "non-fiction", "technical", "children")
    })
    Attribute("stock", Int, func() {
        Minimum(0)
    })
})

var ErrorResult = Type("ErrorResult", func() {
    Attribute("code", String, "Error code")
    Attribute("message", String, "Error message")
    Attribute("details", MapOf(String, String), "Additional details")
    Required("code", "message")
})
```

```go
// design/books.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("books", func() {
    Description("Book management service")
    
    // Service-level errors
    Error("not_found", ErrorResult, "Book not found")
    Error("bad_request", ErrorResult, "Invalid request")
    Error("conflict", ErrorResult, "ISBN already exists")
    Error("internal", ErrorResult, "Internal server error")
    
    // List books
    Method("list", func() {
        Description("List all books with pagination and filtering")
        
        Payload(func() {
            Attribute("page", Int, "Page number", func() {
                Default(1)
                Minimum(1)
            })
            Attribute("limit", Int, "Items per page", func() {
                Default(20)
                Minimum(1)
                Maximum(100)
            })
            Attribute("category", String, "Filter by category", func() {
                Enum("fiction", "non-fiction", "technical", "children")
            })
            Attribute("search", String, "Search in title/author")
            Attribute("sort", String, "Sort field", func() {
                Default("created_at")
                Enum("title", "author", "price", "created_at")
            })
            Attribute("order", String, "Sort order", func() {
                Default("desc")
                Enum("asc", "desc")
            })
        })
        
        Result(BookList)
        
        HTTP(func() {
            GET("/books")
            Param("page")
            Param("limit")
            Param("category")
            Param("search:q")  // ?q=value
            Param("sort")
            Param("order")
            Response(StatusOK)
        })
    })
    
    // Get single book
    Method("get", func() {
        Description("Get a book by ID")
        
        Payload(func() {
            Attribute("id", String, "Book ID")
            Required("id")
        })
        
        Result(Book)
        Error("not_found")
        
        HTTP(func() {
            GET("/books/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Get by ISBN
    Method("get_by_isbn", func() {
        Description("Get a book by ISBN")
        
        Payload(func() {
            Attribute("isbn", String, "ISBN-13", func() {
                Pattern(`^\d{13}$`)
            })
            Required("isbn")
        })
        
        Result(Book)
        Error("not_found")
        
        HTTP(func() {
            GET("/books/isbn/{isbn}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Create book
    Method("create", func() {
        Description("Add a new book")
        
        Payload(func() {
            Attribute("auth_token", String, "Authorization token")
            Extend(CreateBookPayload)
            Required("auth_token")
        })
        
        Result(Book)
        Error("bad_request")
        Error("conflict")
        
        HTTP(func() {
            POST("/books")
            Header("Authorization:auth_token")
            Response(StatusCreated, func() {
                Header("Location")
            })
            Response("bad_request", StatusBadRequest)
            Response("conflict", StatusConflict)
        })
    })
    
    // Update book
    Method("update", func() {
        Description("Update book details")
        
        Payload(func() {
            Attribute("id", String, "Book ID")
            Attribute("auth_token", String, "Authorization token")
            Extend(UpdateBookPayload)
            Required("id", "auth_token")
        })
        
        Result(Book)
        Error("not_found")
        Error("bad_request")
        
        HTTP(func() {
            PUT("/books/{id}")
            Header("Authorization:auth_token")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("bad_request", StatusBadRequest)
        })
    })
    
    // Partial update
    Method("patch", func() {
        Description("Partially update book")
        
        Payload(func() {
            Attribute("id", String, "Book ID")
            Attribute("auth_token", String)
            Attribute("stock", Int, "Update stock only", func() {
                Minimum(0)
            })
            Required("id", "auth_token", "stock")
        })
        
        Result(Book)
        Error("not_found")
        
        HTTP(func() {
            PATCH("/books/{id}")
            Header("Authorization:auth_token")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Delete book
    Method("delete", func() {
        Description("Remove a book")
        
        Payload(func() {
            Attribute("id", String, "Book ID")
            Attribute("auth_token", String)
            Required("id", "auth_token")
        })
        
        Error("not_found")
        
        HTTP(func() {
            DELETE("/books/{id}")
            Header("Authorization:auth_token")
            Response(StatusNoContent)
            Response("not_found", StatusNotFound)
        })
    })
})
```

### Implementation

```go
// books.go
package api

import (
    "context"
    "sync"
    "time"
    
    "github.com/google/uuid"
    books "myproject/gen/books"
)

type booksSvc struct {
    mu     sync.RWMutex
    books  map[string]*books.Book
    byISBN map[string]string  // ISBN -> ID
}

func NewBooksService() books.Service {
    return &booksSvc{
        books:  make(map[string]*books.Book),
        byISBN: make(map[string]string),
    }
}

func (s *booksSvc) List(ctx context.Context, p *books.ListPayload) (*books.BookList, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var result []*books.Book
    
    for _, book := range s.books {
        // Apply category filter
        if p.Category != nil && book.Category != *p.Category {
            continue
        }
        
        // Apply search filter
        if p.Search != nil {
            search := strings.ToLower(*p.Search)
            if !strings.Contains(strings.ToLower(book.Title), search) &&
               !strings.Contains(strings.ToLower(book.Author), search) {
                continue
            }
        }
        
        result = append(result, book)
    }
    
    // Sort results
    sortField := "created_at"
    if p.Sort != nil {
        sortField = *p.Sort
    }
    sortOrder := "desc"
    if p.Order != nil {
        sortOrder = *p.Order
    }
    
    sort.Slice(result, func(i, j int) bool {
        var less bool
        switch sortField {
        case "title":
            less = result[i].Title < result[j].Title
        case "author":
            less = result[i].Author < result[j].Author
        case "price":
            less = result[i].Price < result[j].Price
        default:
            less = *result[i].CreatedAt < *result[j].CreatedAt
        }
        if sortOrder == "desc" {
            return !less
        }
        return less
    })
    
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
    
    return &books.BookList{
        Books: result[start:end],
        Total: total,
        Page:  *p.Page,
        Limit: *p.Limit,
    }, nil
}

func (s *booksSvc) Get(ctx context.Context, p *books.GetPayload) (*books.Book, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    book, ok := s.books[p.ID]
    if !ok {
        return nil, books.MakeNotFound(fmt.Errorf("book %q not found", p.ID))
    }
    return book, nil
}

func (s *booksSvc) GetByIsbn(ctx context.Context, p *books.GetByIsbnPayload) (*books.Book, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    id, ok := s.byISBN[p.Isbn]
    if !ok {
        return nil, books.MakeNotFound(fmt.Errorf("book with ISBN %q not found", p.Isbn))
    }
    return s.books[id], nil
}

func (s *booksSvc) Create(ctx context.Context, p *books.CreatePayload) (*books.Book, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Check for duplicate ISBN
    if _, exists := s.byISBN[p.Isbn]; exists {
        return nil, books.MakeConflict(fmt.Errorf("book with ISBN %q already exists", p.Isbn))
    }
    
    id := uuid.New().String()
    now := time.Now().Format(time.RFC3339)
    
    stock := 0
    if p.Stock != nil {
        stock = *p.Stock
    }
    
    book := &books.Book{
        ID:        id,
        Isbn:      p.Isbn,
        Title:     p.Title,
        Author:    p.Author,
        Price:     p.Price,
        Category:  p.Category,
        Stock:     &stock,
        CreatedAt: &now,
    }
    
    s.books[id] = book
    s.byISBN[p.Isbn] = id
    
    return book, nil
}

func (s *booksSvc) Update(ctx context.Context, p *books.UpdatePayload) (*books.Book, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    book, ok := s.books[p.ID]
    if !ok {
        return nil, books.MakeNotFound(fmt.Errorf("book %q not found", p.ID))
    }
    
    if p.Title != nil {
        book.Title = *p.Title
    }
    if p.Author != nil {
        book.Author = *p.Author
    }
    if p.Price != nil {
        book.Price = *p.Price
    }
    if p.Category != nil {
        book.Category = *p.Category
    }
    if p.Stock != nil {
        book.Stock = p.Stock
    }
    
    return book, nil
}

func (s *booksSvc) Patch(ctx context.Context, p *books.PatchPayload) (*books.Book, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    book, ok := s.books[p.ID]
    if !ok {
        return nil, books.MakeNotFound(fmt.Errorf("book %q not found", p.ID))
    }
    
    book.Stock = &p.Stock
    return book, nil
}

func (s *booksSvc) Delete(ctx context.Context, p *books.DeletePayload) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    book, ok := s.books[p.ID]
    if !ok {
        return books.MakeNotFound(fmt.Errorf("book %q not found", p.ID))
    }
    
    delete(s.byISBN, book.Isbn)
    delete(s.books, p.ID)
    
    return nil
}
```

### Complete gRPC Streaming Service

```go
// design/events.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("events", func() {
    Description("Real-time event streaming service")
    
    // Server streaming: Subscribe to events
    Method("subscribe", func() {
        Description("Subscribe to event stream")
        
        Payload(func() {
            Attribute("topics", ArrayOf(String), "Topics to subscribe")
            Attribute("from_sequence", Int64, "Start from sequence number")
            Required("topics")
        })
        
        StreamingResult(Event)
        
        GRPC(func() {
            Response(CodeOK)
        })
    })
    
    // Client streaming: Batch publish events
    Method("publish_batch", func() {
        Description("Publish multiple events")
        
        StreamingPayload(PublishEvent)
        
        Result(func() {
            Attribute("published", Int, "Number published")
            Attribute("failed", Int, "Number failed")
            Required("published", "failed")
        })
        
        GRPC(func() {
            Response(CodeOK)
        })
    })
    
    // Bidirectional: Real-time sync
    Method("sync", func() {
        Description("Bidirectional event synchronization")
        
        StreamingPayload(SyncRequest)
        StreamingResult(SyncResponse)
        
        GRPC(func() {
            Response(CodeOK)
        })
    })
})

var Event = Type("Event", func() {
    Attribute("id", String, "Event ID")
    Attribute("topic", String, "Event topic")
    Attribute("data", Bytes, "Event data")
    Attribute("sequence", Int64, "Sequence number")
    Attribute("timestamp", Int64, "Unix timestamp")
    Required("id", "topic", "data", "sequence", "timestamp")
})

var PublishEvent = Type("PublishEvent", func() {
    Attribute("topic", String, "Event topic")
    Attribute("data", Bytes, "Event data")
    Attribute("idempotency_key", String, "Idempotency key")
    Required("topic", "data")
})

var SyncRequest = Type("SyncRequest", func() {
    Attribute("client_id", String, "Client identifier")
    Attribute("events", ArrayOf(Event), "Events to sync")
    Attribute("ack_sequence", Int64, "Acknowledge up to sequence")
    Required("client_id")
})

var SyncResponse = Type("SyncResponse", func() {
    Attribute("events", ArrayOf(Event), "Events from server")
    Attribute("server_sequence", Int64, "Server's latest sequence")
    Required("server_sequence")
})
```

---

## 📝 Summary

### HTTP Transport
- **Routing**: Use `GET`, `POST`, `PUT`, `PATCH`, `DELETE` with path patterns
- **Path parameters**: `{param}` or `{url_param:payload_attr}` for mapping
- **Query parameters**: `Param("name")` or implicit from payload
- **Headers**: `Header("Header-Name:attribute")` for request/response
- **Body**: Automatic JSON marshaling, or explicit `Body()` mapping
- **Status codes**: `Response(StatusCode)` for success and errors

### gRPC Transport
- **Proto generation**: Automatic from DSL via `goa gen`
- **Status codes**: Use `Code*` constants matching gRPC semantics
- **Metadata**: Map to headers with `Metadata()` block
- **Message**: Map to proto fields with `Message()` block
- **Field numbers**: Control with `Meta("rpc:tag", "N")`

### Streaming
- **Server streaming**: `StreamingResult(Type)` - server sends multiple responses
- **Client streaming**: `StreamingPayload(Type)` - client sends multiple requests
- **Bidirectional**: Both `StreamingPayload` and `StreamingResult`

### Generated Clients
- **HTTP client**: Uses net/http, full request/response handling
- **gRPC client**: Uses grpc.ClientConn, protobuf encoding
- **Both share**: Same service interface and endpoint abstraction

---

## 📋 Knowledge Check

Before proceeding to Part 4 (Advanced Topics), ensure you can:

- [ ] Configure HTTP routing with path patterns
- [ ] Map payload attributes to path, query, headers, and body
- [ ] Use explicit `Param()`, `Header()`, and `Body()` mappings
- [ ] Set appropriate HTTP status codes for success and errors
- [ ] Enable gRPC transport with `GRPC()` block
- [ ] Understand generated proto file structure
- [ ] Map gRPC metadata and message fields
- [ ] Implement server streaming methods
- [ ] Implement client streaming methods
- [ ] Implement bidirectional streaming
- [ ] Use generated HTTP and gRPC clients
- [ ] Set up a multi-transport server

---

## 🔗 Quick Reference Links

- [Goa HTTP DSL](https://pkg.go.dev/goa.design/goa/v3/dsl#HTTP)
- [Goa gRPC DSL](https://pkg.go.dev/goa.design/goa/v3/dsl#GRPC)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [gRPC Status Codes](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)
- [Protocol Buffers](https://protobuf.dev/)
- [gRPC Streaming](https://grpc.io/docs/what-is-grpc/core-concepts/#server-streaming-rpc)

---

> **Next Up:** Part 4 - Advanced Topics (Middleware, Logging, Error Handling, Testing)
