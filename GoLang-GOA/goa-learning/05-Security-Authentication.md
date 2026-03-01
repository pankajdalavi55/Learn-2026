# Part 5: Security & Authentication

> **Goal:** Master Goa's security features - JWT, API keys, Basic Auth, OAuth2, and implementing authentication in your services

---

## 📚 Table of Contents

1. [Security Overview](#security-overview)
2. [Security Schemes in DSL](#security-schemes-in-dsl)
3. [API Key Authentication](#api-key-authentication)
4. [Basic Authentication](#basic-authentication)
5. [JWT Authentication](#jwt-authentication)
6. [OAuth2 Flows](#oauth2-flows)
7. [Applying Security](#applying-security)
8. [Implementing Security Handlers](#implementing-security-handlers)
9. [Security Best Practices](#security-best-practices)
10. [Complete Examples](#complete-examples)
11. [Summary](#summary)
12. [Knowledge Check](#knowledge-check)

---

## 🎯 Security Overview

### What is Authentication vs Authorization?

```
┌─────────────────────────────────────────────────────────────────┐
│              AUTHENTICATION vs AUTHORIZATION                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  AUTHENTICATION (AuthN)                                         │
│  ────────────────────                                           │
│  "Who are you?"                                                 │
│                                                                 │
│  ┌─────────┐        ┌──────────────┐        ┌─────────┐        │
│  │ Client  │───────▶│  Credentials │───────▶│ Identity│        │
│  └─────────┘        └──────────────┘        └─────────┘        │
│                                                                 │
│  Examples:                                                      │
│  • Username/Password  →  User ID                                │
│  • API Key           →  Application ID                          │
│  • JWT Token         →  Claims (user, roles)                    │
│  • OAuth2 Token      →  Scopes, User info                       │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  AUTHORIZATION (AuthZ)                                          │
│  ───────────────────                                            │
│  "What can you do?"                                             │
│                                                                 │
│  ┌─────────┐        ┌──────────────┐        ┌─────────┐        │
│  │Identity │───────▶│ Permissions  │───────▶│ Allow/  │        │
│  │         │        │    Check     │        │  Deny   │        │
│  └─────────┘        └──────────────┘        └─────────┘        │
│                                                                 │
│  Examples:                                                      │
│  • Role-based: admin can delete, user can read                  │
│  • Scope-based: read:users, write:posts                         │
│  • Attribute-based: owner can modify own resource               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Security in Goa Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      REQUEST FLOW                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    1. CLIENT REQUEST                      │  │
│  │                                                           │  │
│  │  GET /api/users                                           │  │
│  │  Authorization: Bearer eyJhbGciOiJIUzI1N...               │  │
│  │  X-API-Key: sk_live_abc123                                │  │
│  └──────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │               2. TRANSPORT LAYER                          │  │
│  │                                                           │  │
│  │  • Extract credentials from headers/query/cookies         │  │
│  │  • Map to Payload attributes                              │  │
│  └──────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │               3. SECURITY HANDLER                         │  │
│  │                                                           │  │
│  │  • Validate credentials                                   │  │
│  │  • Verify signatures/expiry                               │  │
│  │  • Extract identity (user ID, scopes)                     │  │
│  │  • Store in context                                       │  │
│  │                                                           │  │
│  │  ┌────────────────────────────────────────────────────┐  │  │
│  │  │  func(ctx, token) → (context, error)               │  │  │
│  │  │                                                     │  │  │
│  │  │  • SUCCESS: Return context with auth info          │  │  │
│  │  │  • FAILURE: Return error (401/403)                 │  │  │
│  │  └────────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │               4. SERVICE METHOD                           │  │
│  │                                                           │  │
│  │  • Receives validated context                             │  │
│  │  • Can access user info from context                      │  │
│  │  • Performs business logic                                │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Security Scheme Types in Goa

| Scheme | Use Case | Location | Example |
|--------|----------|----------|---------|
| **APIKey** | Server-to-server, simple auth | Header, Query, Cookie | `X-API-Key: abc123` |
| **BasicAuth** | Simple username/password | Header | `Authorization: Basic base64(user:pass)` |
| **JWTSecurity** | Stateless token auth | Header, Query, Cookie | `Authorization: Bearer jwt...` |
| **OAuth2** | Delegated authorization | Header | `Authorization: Bearer access_token` |

---

## 🔑 Security Schemes in DSL

### Defining Security Schemes

Security schemes are defined at the API or Service level and then applied to methods.

#### Basic Structure

```go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// Define security schemes at API level
var _ = API("myapi", func() {
    Title("My Secure API")
    
    // Security schemes are defined here or in Service
})

// API Key scheme
var APIKeyAuth = APIKeySecurity("api_key", func() {
    Description("API key for authentication")
})

// Basic Auth scheme
var BasicAuth = BasicAuthSecurity("basic", func() {
    Description("Username and password authentication")
})

// JWT scheme
var JWTAuth = JWTSecurity("jwt", func() {
    Description("JWT token authentication")
})

// OAuth2 scheme
var OAuth2Auth = OAuth2Security("oauth2", func() {
    Description("OAuth2 authentication")
})
```

### Where to Define Schemes

```
┌─────────────────────────────────────────────────────────────────┐
│                SECURITY SCHEME DEFINITION LEVELS                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PACKAGE LEVEL (Recommended)                                    │
│  ──────────────────────────                                     │
│  var JWTAuth = JWTSecurity("jwt", func() { ... })               │
│                                                                 │
│  • Reusable across services                                     │
│  • Clear, named reference                                       │
│  • Easy to find and maintain                                    │
│                                                                 │
│  SERVICE LEVEL                                                  │
│  ─────────────                                                  │
│  Service("users", func() {                                      │
│      Security(JWTAuth)  // Apply to all methods                 │
│  })                                                             │
│                                                                 │
│  METHOD LEVEL                                                   │
│  ────────────                                                   │
│  Method("get", func() {                                         │
│      Security(JWTAuth)  // Apply to this method only            │
│  })                                                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔐 API Key Authentication

### What is API Key Authentication?

API Key authentication uses a simple token (key) to identify the client. It's commonly used for:
- Server-to-server communication
- Public APIs with rate limiting
- Simple authentication without user context

```
┌─────────────────────────────────────────────────────────────────┐
│                    API KEY FLOW                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Client obtains API key (from dashboard, registration)       │
│                                                                 │
│  2. Client includes key in request                              │
│     ┌─────────────────────────────────────────────────────┐    │
│     │  GET /api/data                                       │    │
│     │  X-API-Key: sk_live_abc123def456                     │    │
│     └─────────────────────────────────────────────────────┘    │
│                                                                 │
│  3. Server validates key                                        │
│     • Key exists in database?                                   │
│     • Key not expired/revoked?                                  │
│     • Rate limit not exceeded?                                  │
│                                                                 │
│  4. Server processes request or returns 401/403                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Defining API Key Security

```go
// design/security.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// API Key in Header (most common)
var APIKeyAuth = APIKeySecurity("api_key", func() {
    Description("API key must be provided in the X-API-Key header")
})

// API Key in Query Parameter
var APIKeyQueryAuth = APIKeySecurity("api_key_query", func() {
    Description("API key in query parameter")
})

// API Key in Cookie
var APIKeyCookieAuth = APIKeySecurity("api_key_cookie", func() {
    Description("API key in cookie")
})
```

### Applying API Key to Service

```go
// design/service.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("data", func() {
    Description("Data access service")
    
    // Apply API key auth to all methods
    Security(APIKeyAuth)
    
    Method("get", func() {
        Description("Get data")
        
        Payload(func() {
            // The token field receives the API key
            // This is automatically mapped by Goa
            TokenField(1, "api_key", String, func() {
                Description("API key")
            })
            
            // Other payload fields
            Attribute("id", String, "Resource ID")
            Required("api_key", "id")
        })
        
        Result(DataResult)
        
        HTTP(func() {
            GET("/data/{id}")
            // Map API key to header
            Header("api_key:X-API-Key")
            Response(StatusOK)
        })
        
        GRPC(func() {
            // Map API key to gRPC metadata
            Metadata(func() {
                Attribute("api_key:x-api-key")
            })
        })
    })
})
```

### API Key Location Options

```go
// Header (Recommended)
Method("endpoint", func() {
    Security(APIKeyAuth)
    Payload(func() {
        TokenField(1, "key", String)
    })
    HTTP(func() {
        GET("/endpoint")
        Header("key:X-API-Key")  // Custom header
        // OR
        Header("key:Authorization")  // Standard header
    })
})

// Query Parameter
Method("endpoint", func() {
    Security(APIKeyQueryAuth)
    Payload(func() {
        TokenField(1, "key", String)
        Attribute("id", String)
    })
    HTTP(func() {
        GET("/endpoint/{id}")
        Param("key:api_key")  // ?api_key=value
    })
})

// Cookie
Method("endpoint", func() {
    Security(APIKeyCookieAuth)
    Payload(func() {
        TokenField(1, "key", String)
    })
    HTTP(func() {
        GET("/endpoint")
        Cookie("key:session")  // Cookie named "session"
    })
})
```

### API Key Implementation

```go
// security.go
package api

import (
    "context"
    "fmt"
    
    data "myproject/gen/data"
)

// APIKeyStore simulates a database of API keys
type APIKeyStore struct {
    keys map[string]*APIKeyInfo
}

type APIKeyInfo struct {
    ClientID    string
    ClientName  string
    Permissions []string
    RateLimit   int
    Active      bool
}

// Context key for storing authenticated client info
type contextKey string

const clientInfoKey contextKey = "client_info"

// APIKeyAuthFunc implements API key authentication
func APIKeyAuthFunc(store *APIKeyStore) func(context.Context, string) (context.Context, error) {
    return func(ctx context.Context, key string) (context.Context, error) {
        // Validate key format
        if key == "" {
            return ctx, data.MakeUnauthorized(fmt.Errorf("API key required"))
        }
        
        // Look up key in store
        info, exists := store.keys[key]
        if !exists {
            return ctx, data.MakeUnauthorized(fmt.Errorf("invalid API key"))
        }
        
        // Check if key is active
        if !info.Active {
            return ctx, data.MakeUnauthorized(fmt.Errorf("API key is disabled"))
        }
        
        // Store client info in context
        ctx = context.WithValue(ctx, clientInfoKey, info)
        
        return ctx, nil
    }
}

// GetClientInfo retrieves client info from context
func GetClientInfo(ctx context.Context) (*APIKeyInfo, bool) {
    info, ok := ctx.Value(clientInfoKey).(*APIKeyInfo)
    return info, ok
}
```

### Wiring API Key Auth to Server

```go
// cmd/server/main.go
package main

import (
    "context"
    "net/http"
    
    goahttp "goa.design/goa/v3/http"
    data "myproject/gen/data"
    datasvr "myproject/gen/http/data/server"
)

func main() {
    // Create API key store
    store := &APIKeyStore{
        keys: map[string]*APIKeyInfo{
            "sk_live_abc123": {
                ClientID:    "client-1",
                ClientName:  "Test Client",
                Permissions: []string{"read", "write"},
                Active:      true,
            },
        },
    }
    
    // Create service with security handler
    svc := NewDataService()
    endpoints := data.NewEndpoints(svc)
    
    // Create HTTP server with security handler
    mux := goahttp.NewMuxer()
    
    // Security handler function
    apiKeyAuth := APIKeyAuthFunc(store)
    
    // Create server with security
    server := datasvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        nil,  // error handler
        apiKeyAuth,  // API key auth handler
    )
    
    datasvr.Mount(mux, server)
    
    http.ListenAndServe(":8080", mux)
}
```

---

## 👤 Basic Authentication

### What is Basic Authentication?

Basic Auth sends username and password encoded in Base64 with each request.

```
┌─────────────────────────────────────────────────────────────────┐
│                    BASIC AUTH FLOW                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Client has username and password                            │
│     username: "john"                                            │
│     password: "secret123"                                       │
│                                                                 │
│  2. Client encodes credentials                                  │
│     base64("john:secret123") = "am9objpzZWNyZXQxMjM="           │
│                                                                 │
│  3. Client sends request                                        │
│     ┌─────────────────────────────────────────────────────┐    │
│     │  GET /api/profile                                    │    │
│     │  Authorization: Basic am9objpzZWNyZXQxMjM=           │    │
│     └─────────────────────────────────────────────────────┘    │
│                                                                 │
│  4. Server decodes and validates                                │
│     • Decode Base64                                             │
│     • Split by ":"                                              │
│     • Verify username/password                                  │
│                                                                 │
│  ⚠️  SECURITY NOTE: Always use HTTPS with Basic Auth!           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Defining Basic Auth Security

```go
// design/security.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// Define Basic Auth scheme
var BasicAuth = BasicAuthSecurity("basic", func() {
    Description("HTTP Basic authentication")
    // Goa handles the Authorization header automatically
})
```

### Applying Basic Auth

```go
// design/users.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("User service with basic auth")
    
    // Apply basic auth to all methods
    Security(BasicAuth)
    
    Error("unauthorized", ErrorResult)
    
    Method("profile", func() {
        Description("Get current user profile")
        
        Payload(func() {
            // Goa provides special fields for Basic Auth
            UsernameField(1, "username", String, func() {
                Description("Username")
            })
            PasswordField(2, "password", String, func() {
                Description("Password")
            })
        })
        
        Result(UserProfile)
        Error("unauthorized")
        
        HTTP(func() {
            GET("/profile")
            // No explicit mapping needed - Goa handles
            // Authorization header automatically
            Response(StatusOK)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    // Public endpoint without auth
    Method("health", func() {
        Description("Health check")
        
        // Override service-level security
        NoSecurity()
        
        Result(String)
        
        HTTP(func() {
            GET("/health")
        })
    })
})
```

### Basic Auth Payload Fields

```go
// UsernameField and PasswordField are special Goa DSL functions
Payload(func() {
    // First argument is field number (for proto)
    // Second is field name
    // Third is type
    // Fourth is optional description/validation
    
    UsernameField(1, "username", String, func() {
        Description("User's username or email")
        MinLength(3)
        MaxLength(100)
    })
    
    PasswordField(2, "password", String, func() {
        Description("User's password")
        MinLength(8)
    })
    
    // Additional fields for the endpoint
    Attribute("include_stats", Boolean, func() {
        Default(false)
    })
})
```

### Basic Auth Implementation

```go
// security.go
package api

import (
    "context"
    "fmt"
    
    "golang.org/x/crypto/bcrypt"
    users "myproject/gen/users"
)

// UserStore simulates user database
type UserStore struct {
    users map[string]*User
}

type User struct {
    ID           string
    Username     string
    PasswordHash string  // bcrypt hash
    Email        string
    Role         string
}

// Context key for user
type userContextKey string

const currentUserKey userContextKey = "current_user"

// BasicAuthFunc implements basic authentication
func BasicAuthFunc(store *UserStore) func(context.Context, string, string) (context.Context, error) {
    return func(ctx context.Context, username, password string) (context.Context, error) {
        // Find user
        user, exists := store.users[username]
        if !exists {
            // Don't reveal if user exists or not
            return ctx, users.MakeUnauthorized(fmt.Errorf("invalid credentials"))
        }
        
        // Verify password with bcrypt
        err := bcrypt.CompareHashAndPassword(
            []byte(user.PasswordHash),
            []byte(password),
        )
        if err != nil {
            return ctx, users.MakeUnauthorized(fmt.Errorf("invalid credentials"))
        }
        
        // Store user in context
        ctx = context.WithValue(ctx, currentUserKey, user)
        
        return ctx, nil
    }
}

// GetCurrentUser retrieves user from context
func GetCurrentUser(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(currentUserKey).(*User)
    return user, ok
}

// HashPassword creates bcrypt hash (for user registration)
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), err
}
```

### Wiring Basic Auth

```go
// cmd/server/main.go
package main

import (
    "context"
    "net/http"
    
    goahttp "goa.design/goa/v3/http"
    users "myproject/gen/users"
    userssvr "myproject/gen/http/users/server"
)

func main() {
    // Create user store with sample user
    passwordHash, _ := HashPassword("secret123")
    store := &UserStore{
        users: map[string]*User{
            "john": {
                ID:           "user-1",
                Username:     "john",
                PasswordHash: passwordHash,
                Email:        "john@example.com",
                Role:         "admin",
            },
        },
    }
    
    // Create service
    svc := NewUsersService()
    endpoints := users.NewEndpoints(svc)
    
    // Create server with basic auth handler
    mux := goahttp.NewMuxer()
    
    server := userssvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        nil,
        nil,  // API key auth (not used)
        BasicAuthFunc(store),  // Basic auth handler
    )
    
    userssvr.Mount(mux, server)
    
    http.ListenAndServe(":8080", mux)
}
```

---

## 🎫 JWT Authentication

### What is JWT?

JWT (JSON Web Token) is a compact, self-contained token format for secure information exchange.

```
┌─────────────────────────────────────────────────────────────────┐
│                      JWT STRUCTURE                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.                          │
│  eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4ifQ.                │
│  SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c                    │
│                                                                 │
│  ├─────────────────┤ ├───────────────────────────┤ ├─────────┤ │
│       HEADER              PAYLOAD (Claims)          SIGNATURE   │
│                                                                 │
│  HEADER (Base64URL encoded JSON)                                │
│  {                                                              │
│    "alg": "HS256",     // Signing algorithm                     │
│    "typ": "JWT"        // Token type                            │
│  }                                                              │
│                                                                 │
│  PAYLOAD (Base64URL encoded JSON)                               │
│  {                                                              │
│    "sub": "1234567890",  // Subject (user ID)                   │
│    "name": "John Doe",   // Custom claim                        │
│    "iat": 1516239022,    // Issued at                           │
│    "exp": 1516242622,    // Expiration                          │
│    "iss": "myapp",       // Issuer                              │
│    "aud": "myapi",       // Audience                            │
│    "roles": ["admin"]    // Custom claim                        │
│  }                                                              │
│                                                                 │
│  SIGNATURE                                                      │
│  HMACSHA256(                                                    │
│    base64UrlEncode(header) + "." + base64UrlEncode(payload),    │
│    secret                                                       │
│  )                                                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### JWT Authentication Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                      JWT AUTH FLOW                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────┐                              ┌──────────┐         │
│  │  Client  │                              │  Server  │         │
│  └────┬─────┘                              └────┬─────┘         │
│       │                                         │               │
│       │  1. POST /login                         │               │
│       │     {username, password}                │               │
│       │─────────────────────────────────────────▶               │
│       │                                         │               │
│       │  2. Validate credentials                │               │
│       │     Generate JWT with claims            │               │
│       │                                         │               │
│       │  3. Return JWT                          │               │
│       │     {token: "eyJhbG..."}                │               │
│       │◀─────────────────────────────────────────               │
│       │                                         │               │
│       │  4. Store token (localStorage, etc)     │               │
│       │                                         │               │
│       │  5. GET /api/protected                  │               │
│       │     Authorization: Bearer eyJhbG...     │               │
│       │─────────────────────────────────────────▶               │
│       │                                         │               │
│       │  6. Verify signature                    │               │
│       │     Check expiration                    │               │
│       │     Extract claims                      │               │
│       │                                         │               │
│       │  7. Return protected data               │               │
│       │◀─────────────────────────────────────────               │
│       │                                         │               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Defining JWT Security

```go
// design/security.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// JWT Security scheme
var JWTAuth = JWTSecurity("jwt", func() {
    Description("JWT token authentication")
    
    // Define the scopes that can be required
    Scope("api:read", "Read access to API")
    Scope("api:write", "Write access to API")
    Scope("admin", "Admin access")
})

// Multiple JWT schemes for different use cases
var UserJWT = JWTSecurity("user_jwt", func() {
    Description("User JWT for web/mobile clients")
    Scope("user:read", "Read user data")
    Scope("user:write", "Modify user data")
})

var ServiceJWT = JWTSecurity("service_jwt", func() {
    Description("Service-to-service JWT")
    Scope("service:internal", "Internal service communication")
})
```

### JWT Token Location Options

```go
// Bearer token in Authorization header (most common)
Method("get", func() {
    Security(JWTAuth)
    Payload(func() {
        TokenField(1, "token", String, func() {
            Description("JWT token")
        })
        Attribute("id", String)
        Required("token", "id")
    })
    HTTP(func() {
        GET("/resource/{id}")
        // Default: looks for "Authorization: Bearer <token>"
    })
})

// Token in custom header
Method("get", func() {
    Security(JWTAuth)
    Payload(func() {
        TokenField(1, "token", String)
        Attribute("id", String)
    })
    HTTP(func() {
        GET("/resource/{id}")
        Header("token:X-Auth-Token")  // Custom header
    })
})

// Token in query parameter (for WebSocket, SSE)
Method("stream", func() {
    Security(JWTAuth)
    Payload(func() {
        TokenField(1, "token", String)
    })
    HTTP(func() {
        GET("/stream")
        Param("token:access_token")  // ?access_token=jwt...
    })
})

// Token in cookie (for web apps)
Method("dashboard", func() {
    Security(JWTAuth)
    Payload(func() {
        TokenField(1, "token", String)
    })
    HTTP(func() {
        GET("/dashboard")
        Cookie("token:auth_token")  // Cookie named "auth_token"
    })
})
```

### JWT with Scopes

```go
var _ = Service("resources", func() {
    Description("Resource management with JWT")
    
    // Require jwt auth for all methods
    Security(JWTAuth)
    
    // Read operation - requires read scope
    Method("list", func() {
        Description("List resources")
        
        // Require specific scope
        Security(JWTAuth, func() {
            Scope("api:read")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
        })
        
        Result(ArrayOf(Resource))
        
        HTTP(func() {
            GET("/resources")
        })
    })
    
    // Write operation - requires write scope
    Method("create", func() {
        Description("Create resource")
        
        Security(JWTAuth, func() {
            Scope("api:write")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("name", String)
            Attribute("data", Any)
            Required("token", "name")
        })
        
        Result(Resource)
        
        HTTP(func() {
            POST("/resources")
        })
    })
    
    // Admin operation - requires admin scope
    Method("delete_all", func() {
        Description("Delete all resources (admin only)")
        
        Security(JWTAuth, func() {
            Scope("admin")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
        })
        
        HTTP(func() {
            DELETE("/resources")
            Response(StatusNoContent)
        })
    })
})
```

### JWT Implementation

```go
// security.go
package api

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    resources "myproject/gen/resources"
)

// JWT configuration
type JWTConfig struct {
    SecretKey     []byte
    Issuer        string
    Audience      string
    TokenDuration time.Duration
}

// Custom claims structure
type Claims struct {
    jwt.RegisteredClaims
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Email    string   `json:"email"`
    Roles    []string `json:"roles"`
    Scopes   []string `json:"scopes"`
}

// Context key for claims
type claimsContextKey string

const jwtClaimsKey claimsContextKey = "jwt_claims"

// JWTAuthFunc creates a JWT authentication handler
func JWTAuthFunc(config *JWTConfig, requiredScopes []string) func(context.Context, string) (context.Context, error) {
    return func(ctx context.Context, tokenString string) (context.Context, error) {
        // Remove "Bearer " prefix if present
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        
        if tokenString == "" {
            return ctx, resources.MakeUnauthorized(fmt.Errorf("token required"))
        }
        
        // Parse and validate token
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            // Validate signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return config.SecretKey, nil
        })
        
        if err != nil {
            return ctx, resources.MakeUnauthorized(fmt.Errorf("invalid token: %v", err))
        }
        
        claims, ok := token.Claims.(*Claims)
        if !ok || !token.Valid {
            return ctx, resources.MakeUnauthorized(fmt.Errorf("invalid token claims"))
        }
        
        // Validate issuer
        if claims.Issuer != config.Issuer {
            return ctx, resources.MakeUnauthorized(fmt.Errorf("invalid issuer"))
        }
        
        // Validate audience
        if !claims.VerifyAudience(config.Audience, true) {
            return ctx, resources.MakeUnauthorized(fmt.Errorf("invalid audience"))
        }
        
        // Validate required scopes
        if !hasRequiredScopes(claims.Scopes, requiredScopes) {
            return ctx, resources.MakeForbidden(fmt.Errorf("insufficient permissions"))
        }
        
        // Store claims in context
        ctx = context.WithValue(ctx, jwtClaimsKey, claims)
        
        return ctx, nil
    }
}

// hasRequiredScopes checks if all required scopes are present
func hasRequiredScopes(tokenScopes, required []string) bool {
    if len(required) == 0 {
        return true
    }
    
    scopeMap := make(map[string]bool)
    for _, s := range tokenScopes {
        scopeMap[s] = true
    }
    
    for _, r := range required {
        if !scopeMap[r] {
            return false
        }
    }
    return true
}

// GetJWTClaims retrieves claims from context
func GetJWTClaims(ctx context.Context) (*Claims, bool) {
    claims, ok := ctx.Value(jwtClaimsKey).(*Claims)
    return claims, ok
}

// GenerateJWT creates a new JWT token
func GenerateJWT(config *JWTConfig, userID, username, email string, roles, scopes []string) (string, error) {
    now := time.Now()
    
    claims := &Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    config.Issuer,
            Subject:   userID,
            Audience:  jwt.ClaimStrings{config.Audience},
            ExpiresAt: jwt.NewNumericDate(now.Add(config.TokenDuration)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            ID:        generateTokenID(),
        },
        UserID:   userID,
        Username: username,
        Email:    email,
        Roles:    roles,
        Scopes:   scopes,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(config.SecretKey)
}

func generateTokenID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

### Login Endpoint for JWT

```go
// design/auth.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("auth", func() {
    Description("Authentication service")
    
    // Login - no security (public endpoint)
    Method("login", func() {
        Description("Authenticate and get JWT token")
        
        NoSecurity()  // Public endpoint
        
        Payload(func() {
            Attribute("username", String, func() {
                MinLength(3)
                MaxLength(50)
            })
            Attribute("password", String, func() {
                MinLength(8)
            })
            Required("username", "password")
        })
        
        Result(func() {
            Attribute("token", String, "JWT access token")
            Attribute("refresh_token", String, "Refresh token")
            Attribute("expires_in", Int, "Token lifetime in seconds")
            Attribute("token_type", String, "Token type", func() {
                Default("Bearer")
            })
            Required("token", "expires_in", "token_type")
        })
        
        Error("unauthorized")
        
        HTTP(func() {
            POST("/auth/login")
            Response(StatusOK)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    // Refresh token
    Method("refresh", func() {
        Description("Refresh access token")
        
        NoSecurity()
        
        Payload(func() {
            Attribute("refresh_token", String, "Refresh token")
            Required("refresh_token")
        })
        
        Result(func() {
            Attribute("token", String)
            Attribute("expires_in", Int)
            Attribute("token_type", String, func() { Default("Bearer") })
            Required("token", "expires_in", "token_type")
        })
        
        Error("unauthorized")
        
        HTTP(func() {
            POST("/auth/refresh")
            Response(StatusOK)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    // Logout (optional - invalidate refresh token)
    Method("logout", func() {
        Description("Invalidate refresh token")
        
        Security(JWTAuth)
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("refresh_token", String)
        })
        
        HTTP(func() {
            POST("/auth/logout")
            Response(StatusNoContent)
        })
    })
})
```

### Login Implementation

```go
// auth.go
package api

import (
    "context"
    "fmt"
    
    "golang.org/x/crypto/bcrypt"
    auth "myproject/gen/auth"
)

type authSvc struct {
    config    *JWTConfig
    userStore *UserStore
}

func NewAuthService(config *JWTConfig, store *UserStore) auth.Service {
    return &authSvc{
        config:    config,
        userStore: store,
    }
}

func (s *authSvc) Login(ctx context.Context, p *auth.LoginPayload) (*auth.LoginResult, error) {
    // Find user
    user, exists := s.userStore.GetByUsername(p.Username)
    if !exists {
        return nil, auth.MakeUnauthorized(fmt.Errorf("invalid credentials"))
    }
    
    // Verify password
    err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p.Password))
    if err != nil {
        return nil, auth.MakeUnauthorized(fmt.Errorf("invalid credentials"))
    }
    
    // Determine scopes based on user role
    scopes := getScopesForRole(user.Role)
    
    // Generate access token
    token, err := GenerateJWT(
        s.config,
        user.ID,
        user.Username,
        user.Email,
        []string{user.Role},
        scopes,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %v", err)
    }
    
    // Generate refresh token (longer lived, stored in DB)
    refreshToken, err := s.generateRefreshToken(user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %v", err)
    }
    
    return &auth.LoginResult{
        Token:        token,
        RefreshToken: &refreshToken,
        ExpiresIn:    int(s.config.TokenDuration.Seconds()),
        TokenType:    "Bearer",
    }, nil
}

func (s *authSvc) Refresh(ctx context.Context, p *auth.RefreshPayload) (*auth.RefreshResult, error) {
    // Validate refresh token and get user
    userID, err := s.validateRefreshToken(p.RefreshToken)
    if err != nil {
        return nil, auth.MakeUnauthorized(err)
    }
    
    user, exists := s.userStore.GetByID(userID)
    if !exists {
        return nil, auth.MakeUnauthorized(fmt.Errorf("user not found"))
    }
    
    // Generate new access token
    scopes := getScopesForRole(user.Role)
    token, err := GenerateJWT(
        s.config,
        user.ID,
        user.Username,
        user.Email,
        []string{user.Role},
        scopes,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %v", err)
    }
    
    return &auth.RefreshResult{
        Token:     token,
        ExpiresIn: int(s.config.TokenDuration.Seconds()),
        TokenType: "Bearer",
    }, nil
}

func (s *authSvc) Logout(ctx context.Context, p *auth.LogoutPayload) error {
    // Invalidate refresh token if provided
    if p.RefreshToken != nil {
        s.invalidateRefreshToken(*p.RefreshToken)
    }
    return nil
}

func getScopesForRole(role string) []string {
    switch role {
    case "admin":
        return []string{"api:read", "api:write", "admin"}
    case "user":
        return []string{"api:read", "api:write"}
    case "viewer":
        return []string{"api:read"}
    default:
        return []string{}
    }
}
```

---

## 🔄 OAuth2 Flows

### What is OAuth2?

OAuth2 is an authorization framework that enables third-party applications to access resources on behalf of a user.

```
┌─────────────────────────────────────────────────────────────────┐
│                    OAUTH2 ACTORS                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐     ┌─────────────────┐                   │
│  │  Resource Owner │     │ Authorization   │                   │
│  │     (User)      │     │    Server       │                   │
│  └────────┬────────┘     └────────┬────────┘                   │
│           │                       │                            │
│           │ Grants permission     │ Issues tokens              │
│           │                       │                            │
│           ▼                       ▼                            │
│  ┌─────────────────┐     ┌─────────────────┐                   │
│  │     Client      │────▶│    Resource     │                   │
│  │  (Your App)     │     │     Server      │                   │
│  └─────────────────┘     │    (Your API)   │                   │
│                          └─────────────────┘                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### OAuth2 Grant Types

```
┌─────────────────────────────────────────────────────────────────┐
│                   OAUTH2 GRANT TYPES                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. AUTHORIZATION CODE (Most Secure)                            │
│     Use: Web apps with server backend                           │
│     Flow: User → Auth Server → Code → Token                     │
│                                                                 │
│  2. AUTHORIZATION CODE WITH PKCE                                │
│     Use: Mobile/SPA apps (no client secret)                     │
│     Flow: Same as above + code verifier                         │
│                                                                 │
│  3. CLIENT CREDENTIALS                                          │
│     Use: Server-to-server (no user)                             │
│     Flow: Client ID + Secret → Token                            │
│                                                                 │
│  4. IMPLICIT (Deprecated)                                       │
│     Use: Legacy SPAs                                            │
│     Flow: User → Auth Server → Token directly                   │
│                                                                 │
│  5. PASSWORD (Deprecated for third-party)                       │
│     Use: First-party apps only                                  │
│     Flow: Username + Password → Token                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Authorization Code Flow

```
┌─────────────────────────────────────────────────────────────────┐
│              AUTHORIZATION CODE FLOW                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌────────┐          ┌────────┐          ┌────────────────┐    │
│  │ Client │          │  User  │          │  Auth Server   │    │
│  └───┬────┘          └───┬────┘          └───────┬────────┘    │
│      │                   │                       │              │
│      │ 1. Redirect to auth server               │              │
│      │──────────────────────────────────────────▶              │
│      │                   │                       │              │
│      │                   │ 2. User logs in       │              │
│      │                   │───────────────────────▶              │
│      │                   │                       │              │
│      │                   │ 3. User consents      │              │
│      │                   │───────────────────────▶              │
│      │                   │                       │              │
│      │ 4. Redirect with authorization code       │              │
│      │◀──────────────────────────────────────────               │
│      │                   │                       │              │
│      │ 5. Exchange code for token                │              │
│      │──────────────────────────────────────────▶              │
│      │                   │                       │              │
│      │ 6. Return access token + refresh token    │              │
│      │◀──────────────────────────────────────────               │
│      │                   │                       │              │
│      │ 7. Call API with token                    │              │
│      │────────────────────────────────────────────────────▶    │
│      │                   │                       │   API        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Defining OAuth2 Security

```go
// design/security.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// OAuth2 with Authorization Code flow
var OAuth2Auth = OAuth2Security("oauth2", func() {
    Description("OAuth2 authentication")
    
    // Authorization code flow (web apps)
    AuthorizationCodeFlow(
        "https://auth.example.com/authorize",  // Authorization URL
        "https://auth.example.com/token",      // Token URL
        "https://auth.example.com/refresh",    // Refresh URL (optional)
    )
    
    // Define scopes
    Scope("read:profile", "Read user profile")
    Scope("write:profile", "Update user profile")
    Scope("read:data", "Read user data")
    Scope("write:data", "Write user data")
    Scope("admin", "Administrator access")
})

// OAuth2 with Implicit flow (deprecated, but supported)
var OAuth2ImplicitAuth = OAuth2Security("oauth2_implicit", func() {
    Description("OAuth2 implicit flow (legacy)")
    
    ImplicitFlow(
        "https://auth.example.com/authorize",
    )
    
    Scope("read", "Read access")
    Scope("write", "Write access")
})

// OAuth2 with Client Credentials flow (server-to-server)
var OAuth2ClientCredentials = OAuth2Security("oauth2_client", func() {
    Description("OAuth2 client credentials for service-to-service")
    
    ClientCredentialsFlow(
        "https://auth.example.com/token",
    )
    
    Scope("service:read", "Service read access")
    Scope("service:write", "Service write access")
})

// OAuth2 with Password flow (first-party only)
var OAuth2Password = OAuth2Security("oauth2_password", func() {
    Description("OAuth2 password grant (first-party apps only)")
    
    PasswordFlow(
        "https://auth.example.com/token",
        "https://auth.example.com/refresh",
    )
    
    Scope("user", "User access")
})
```

### Using OAuth2 in Service

```go
// design/service.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("protected", func() {
    Description("OAuth2 protected service")
    
    // Apply OAuth2 to all methods
    Security(OAuth2Auth)
    
    // Profile endpoint - requires read:profile scope
    Method("get_profile", func() {
        Description("Get user profile")
        
        Security(OAuth2Auth, func() {
            Scope("read:profile")
        })
        
        Payload(func() {
            // AccessToken is special field for OAuth2
            AccessTokenField(1, "token", String, func() {
                Description("OAuth2 access token")
            })
        })
        
        Result(UserProfile)
        
        HTTP(func() {
            GET("/profile")
            // Token extracted from Authorization: Bearer header
        })
    })
    
    // Update profile - requires write:profile scope
    Method("update_profile", func() {
        Description("Update user profile")
        
        Security(OAuth2Auth, func() {
            Scope("write:profile")
        })
        
        Payload(func() {
            AccessTokenField(1, "token", String)
            Attribute("name", String)
            Attribute("bio", String)
            Required("token")
        })
        
        Result(UserProfile)
        
        HTTP(func() {
            PUT("/profile")
        })
    })
    
    // Admin endpoint - requires admin scope
    Method("admin_stats", func() {
        Description("Get admin statistics")
        
        Security(OAuth2Auth, func() {
            Scope("admin")
        })
        
        Payload(func() {
            AccessTokenField(1, "token", String)
        })
        
        Result(AdminStats)
        
        HTTP(func() {
            GET("/admin/stats")
        })
    })
})
```

### OAuth2 Token Validation

```go
// security.go
package api

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"
    
    protected "myproject/gen/protected"
)

// OAuth2Config holds OAuth2 configuration
type OAuth2Config struct {
    IntrospectionURL string  // Token introspection endpoint
    ClientID         string  // Client credentials for introspection
    ClientSecret     string
    
    // Or use local validation with public key
    PublicKey        []byte
    Issuer           string
    Audience         string
}

// TokenInfo represents introspection response
type TokenInfo struct {
    Active    bool     `json:"active"`
    Scope     string   `json:"scope"`
    ClientID  string   `json:"client_id"`
    Username  string   `json:"username"`
    Subject   string   `json:"sub"`
    ExpiresAt int64    `json:"exp"`
    IssuedAt  int64    `json:"iat"`
}

// OAuth2AuthFunc creates OAuth2 authentication handler
func OAuth2AuthFunc(config *OAuth2Config, requiredScopes []string) func(context.Context, string) (context.Context, error) {
    return func(ctx context.Context, tokenString string) (context.Context, error) {
        // Remove "Bearer " prefix
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        
        if tokenString == "" {
            return ctx, protected.MakeUnauthorized(fmt.Errorf("access token required"))
        }
        
        // Validate token via introspection
        tokenInfo, err := introspectToken(config, tokenString)
        if err != nil {
            return ctx, protected.MakeUnauthorized(fmt.Errorf("token validation failed: %v", err))
        }
        
        if !tokenInfo.Active {
            return ctx, protected.MakeUnauthorized(fmt.Errorf("token is not active"))
        }
        
        // Check expiration (in case introspection doesn't)
        if tokenInfo.ExpiresAt > 0 && tokenInfo.ExpiresAt < time.Now().Unix() {
            return ctx, protected.MakeUnauthorized(fmt.Errorf("token has expired"))
        }
        
        // Validate required scopes
        tokenScopes := strings.Split(tokenInfo.Scope, " ")
        if !hasRequiredScopes(tokenScopes, requiredScopes) {
            return ctx, protected.MakeForbidden(fmt.Errorf("insufficient scope"))
        }
        
        // Store token info in context
        ctx = context.WithValue(ctx, oauth2TokenKey, tokenInfo)
        
        return ctx, nil
    }
}

// introspectToken calls the OAuth2 introspection endpoint
func introspectToken(config *OAuth2Config, token string) (*TokenInfo, error) {
    req, err := http.NewRequest("POST", config.IntrospectionURL, 
        strings.NewReader("token="+token))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.SetBasicAuth(config.ClientID, config.ClientSecret)
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("introspection failed: %d", resp.StatusCode)
    }
    
    var info TokenInfo
    if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
        return nil, err
    }
    
    return &info, nil
}

// Context key for OAuth2 token info
type oauth2ContextKey string

const oauth2TokenKey oauth2ContextKey = "oauth2_token"

// GetOAuth2TokenInfo retrieves token info from context
func GetOAuth2TokenInfo(ctx context.Context) (*TokenInfo, bool) {
    info, ok := ctx.Value(oauth2TokenKey).(*TokenInfo)
    return info, ok
}
```

---

## 🛡️ Applying Security

### Service-Level Security

```go
var _ = Service("accounts", func() {
    Description("Account management")
    
    // Apply JWT auth to ALL methods in this service
    Security(JWTAuth)
    
    Method("list", func() {
        // Inherits JWTAuth from service
        Payload(func() {
            TokenField(1, "token", String)
        })
        Result(ArrayOf(Account))
        HTTP(func() {
            GET("/accounts")
        })
    })
    
    Method("get", func() {
        // Inherits JWTAuth from service
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("id", String)
            Required("token", "id")
        })
        Result(Account)
        HTTP(func() {
            GET("/accounts/{id}")
        })
    })
})
```

### Method-Level Security

```go
var _ = Service("mixed", func() {
    Description("Service with mixed security")
    
    // Public endpoint
    Method("health", func() {
        Description("Health check - no auth required")
        
        NoSecurity()  // Explicitly public
        
        Result(HealthStatus)
        HTTP(func() {
            GET("/health")
        })
    })
    
    // API Key protected
    Method("webhook", func() {
        Description("Webhook endpoint - API key auth")
        
        Security(APIKeyAuth)
        
        Payload(func() {
            TokenField(1, "key", String)
            Attribute("event", String)
            Attribute("data", Any)
        })
        
        HTTP(func() {
            POST("/webhook")
            Header("key:X-Webhook-Key")
        })
    })
    
    // JWT protected
    Method("user_data", func() {
        Description("User data - JWT auth")
        
        Security(JWTAuth)
        
        Payload(func() {
            TokenField(1, "token", String)
        })
        
        Result(UserData)
        HTTP(func() {
            GET("/user/data")
        })
    })
    
    // OAuth2 protected with scopes
    Method("sensitive", func() {
        Description("Sensitive data - OAuth2 with scope")
        
        Security(OAuth2Auth, func() {
            Scope("read:sensitive")
        })
        
        Payload(func() {
            AccessTokenField(1, "token", String)
        })
        
        Result(SensitiveData)
        HTTP(func() {
            GET("/sensitive")
        })
    })
})
```

### Multiple Security Schemes (OR relationship)

```go
var _ = Service("flexible", func() {
    Description("Service accepting multiple auth methods")
    
    Method("data", func() {
        Description("Accepts either JWT or API Key")
        
        // Multiple Security() calls = OR relationship
        // User can authenticate with EITHER method
        Security(JWTAuth)
        Security(APIKeyAuth)
        
        Payload(func() {
            // Include fields for both auth types
            TokenField(1, "jwt_token", String, func() {
                Description("JWT token")
            })
            APIKeyField(2, "api_key", String, func() {
                Description("API key")
            })
            
            Attribute("id", String)
            Required("id")
            // Note: Either jwt_token OR api_key is required, not both
        })
        
        Result(Data)
        
        HTTP(func() {
            GET("/data/{id}")
            // JWT from Authorization header
            // API key from X-API-Key header
            Header("api_key:X-API-Key")
        })
    })
})
```

### Security with Different Scopes per Method

```go
var _ = Service("documents", func() {
    Description("Document management with granular permissions")
    
    // Base security for all methods
    Security(JWTAuth)
    
    // Read - minimal scope
    Method("list", func() {
        Security(JWTAuth, func() {
            Scope("documents:read")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("folder_id", String)
        })
        
        Result(ArrayOf(DocumentSummary))
        
        HTTP(func() {
            GET("/documents")
        })
    })
    
    // Read single - same scope
    Method("get", func() {
        Security(JWTAuth, func() {
            Scope("documents:read")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("id", String)
            Required("token", "id")
        })
        
        Result(Document)
        
        HTTP(func() {
            GET("/documents/{id}")
        })
    })
    
    // Create - write scope
    Method("create", func() {
        Security(JWTAuth, func() {
            Scope("documents:write")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("title", String)
            Attribute("content", String)
            Required("token", "title")
        })
        
        Result(Document)
        
        HTTP(func() {
            POST("/documents")
            Response(StatusCreated)
        })
    })
    
    // Delete - delete scope (more privileged)
    Method("delete", func() {
        Security(JWTAuth, func() {
            Scope("documents:delete")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("id", String)
            Required("token", "id")
        })
        
        HTTP(func() {
            DELETE("/documents/{id}")
            Response(StatusNoContent)
        })
    })
    
    // Admin - admin scope
    Method("admin_purge", func() {
        Security(JWTAuth, func() {
            Scope("admin")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("older_than_days", Int)
            Required("token", "older_than_days")
        })
        
        Result(func() {
            Attribute("deleted_count", Int)
            Required("deleted_count")
        })
        
        HTTP(func() {
            DELETE("/documents/admin/purge")
        })
    })
})
```

### NoSecurity for Public Endpoints

```go
var _ = Service("api", func() {
    Description("API with public and protected endpoints")
    
    // Default: JWT auth
    Security(JWTAuth)
    
    // Public: Health check
    Method("health", func() {
        NoSecurity()  // Override service-level security
        
        Result(func() {
            Attribute("status", String)
            Attribute("version", String)
        })
        
        HTTP(func() {
            GET("/health")
        })
    })
    
    // Public: API documentation
    Method("docs", func() {
        NoSecurity()
        
        Result(Bytes)
        
        HTTP(func() {
            GET("/docs")
            Response(StatusOK, func() {
                ContentType("text/html")
            })
        })
    })
    
    // Public: OpenAPI spec
    Method("openapi", func() {
        NoSecurity()
        
        Result(Bytes)
        
        HTTP(func() {
            GET("/openapi.json")
            Response(StatusOK, func() {
                ContentType("application/json")
            })
        })
    })
    
    // Protected: Get user profile
    Method("profile", func() {
        // Uses service-level JWTAuth
        
        Payload(func() {
            TokenField(1, "token", String)
        })
        
        Result(UserProfile)
        
        HTTP(func() {
            GET("/profile")
        })
    })
})
```

---

## 🔧 Implementing Security Handlers

### Security Handler Function Signatures

```go
// API Key handler
func APIKeyAuth(ctx context.Context, key string) (context.Context, error)

// Basic Auth handler
func BasicAuth(ctx context.Context, username, password string) (context.Context, error)

// JWT handler
func JWTAuth(ctx context.Context, token string) (context.Context, error)

// OAuth2 handler (same as JWT for token validation)
func OAuth2Auth(ctx context.Context, token string) (context.Context, error)
```

### Complete Security Implementation

```go
// security/security.go
package security

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

// ============ CONTEXT KEYS ============

type contextKey string

const (
    UserKey      contextKey = "user"
    ClientKey    contextKey = "client"
    ScopesKey    contextKey = "scopes"
    RolesKey     contextKey = "roles"
)

// ============ DATA STRUCTURES ============

type User struct {
    ID           string
    Username     string
    Email        string
    PasswordHash string
    Roles        []string
    Active       bool
}

type Client struct {
    ID          string
    Name        string
    APIKey      string
    Permissions []string
    Active      bool
}

type JWTClaims struct {
    jwt.RegisteredClaims
    UserID   string   `json:"uid"`
    Username string   `json:"username"`
    Email    string   `json:"email"`
    Roles    []string `json:"roles"`
    Scopes   []string `json:"scopes"`
}

// ============ STORES ============

type UserStore interface {
    GetByUsername(username string) (*User, error)
    GetByID(id string) (*User, error)
}

type ClientStore interface {
    GetByAPIKey(key string) (*Client, error)
}

// ============ AUTHENTICATORS ============

// APIKeyAuthenticator handles API key validation
type APIKeyAuthenticator struct {
    store ClientStore
}

func NewAPIKeyAuth(store ClientStore) *APIKeyAuthenticator {
    return &APIKeyAuthenticator{store: store}
}

func (a *APIKeyAuthenticator) Authenticate(ctx context.Context, key string) (context.Context, error) {
    if key == "" {
        return ctx, fmt.Errorf("API key is required")
    }
    
    // Remove any prefix (e.g., "ApiKey ")
    key = strings.TrimPrefix(key, "ApiKey ")
    
    client, err := a.store.GetByAPIKey(key)
    if err != nil {
        return ctx, fmt.Errorf("invalid API key")
    }
    
    if !client.Active {
        return ctx, fmt.Errorf("API key is disabled")
    }
    
    // Store client in context
    ctx = context.WithValue(ctx, ClientKey, client)
    ctx = context.WithValue(ctx, ScopesKey, client.Permissions)
    
    return ctx, nil
}

// BasicAuthenticator handles username/password validation
type BasicAuthenticator struct {
    store UserStore
}

func NewBasicAuth(store UserStore) *BasicAuthenticator {
    return &BasicAuthenticator{store: store}
}

func (a *BasicAuthenticator) Authenticate(ctx context.Context, username, password string) (context.Context, error) {
    if username == "" || password == "" {
        return ctx, fmt.Errorf("username and password are required")
    }
    
    user, err := a.store.GetByUsername(username)
    if err != nil {
        return ctx, fmt.Errorf("invalid credentials")
    }
    
    if !user.Active {
        return ctx, fmt.Errorf("account is disabled")
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return ctx, fmt.Errorf("invalid credentials")
    }
    
    // Store user in context
    ctx = context.WithValue(ctx, UserKey, user)
    ctx = context.WithValue(ctx, RolesKey, user.Roles)
    
    return ctx, nil
}

// JWTAuthenticator handles JWT validation
type JWTAuthenticator struct {
    secretKey      []byte
    issuer         string
    audience       string
    requiredScopes []string
}

func NewJWTAuth(secretKey []byte, issuer, audience string) *JWTAuthenticator {
    return &JWTAuthenticator{
        secretKey: secretKey,
        issuer:    issuer,
        audience:  audience,
    }
}

func (a *JWTAuthenticator) WithRequiredScopes(scopes ...string) *JWTAuthenticator {
    newAuth := *a
    newAuth.requiredScopes = scopes
    return &newAuth
}

func (a *JWTAuthenticator) Authenticate(ctx context.Context, tokenString string) (context.Context, error) {
    // Remove "Bearer " prefix
    tokenString = strings.TrimPrefix(tokenString, "Bearer ")
    
    if tokenString == "" {
        return ctx, fmt.Errorf("token is required")
    }
    
    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return a.secretKey, nil
    })
    
    if err != nil {
        return ctx, fmt.Errorf("invalid token: %v", err)
    }
    
    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        return ctx, fmt.Errorf("invalid token claims")
    }
    
    // Validate issuer
    if claims.Issuer != a.issuer {
        return ctx, fmt.Errorf("invalid issuer")
    }
    
    // Validate audience
    if !claims.VerifyAudience(a.audience, true) {
        return ctx, fmt.Errorf("invalid audience")
    }
    
    // Validate required scopes
    if len(a.requiredScopes) > 0 {
        scopeMap := make(map[string]bool)
        for _, s := range claims.Scopes {
            scopeMap[s] = true
        }
        for _, required := range a.requiredScopes {
            if !scopeMap[required] {
                return ctx, fmt.Errorf("missing required scope: %s", required)
            }
        }
    }
    
    // Store claims in context
    user := &User{
        ID:       claims.UserID,
        Username: claims.Username,
        Email:    claims.Email,
        Roles:    claims.Roles,
    }
    
    ctx = context.WithValue(ctx, UserKey, user)
    ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)
    ctx = context.WithValue(ctx, RolesKey, claims.Roles)
    
    return ctx, nil
}

// ============ CONTEXT HELPERS ============

func GetUser(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(UserKey).(*User)
    return user, ok
}

func GetClient(ctx context.Context) (*Client, bool) {
    client, ok := ctx.Value(ClientKey).(*Client)
    return client, ok
}

func GetScopes(ctx context.Context) []string {
    scopes, ok := ctx.Value(ScopesKey).([]string)
    if !ok {
        return nil
    }
    return scopes
}

func GetRoles(ctx context.Context) []string {
    roles, ok := ctx.Value(RolesKey).([]string)
    if !ok {
        return nil
    }
    return roles
}

func HasScope(ctx context.Context, scope string) bool {
    scopes := GetScopes(ctx)
    for _, s := range scopes {
        if s == scope {
            return true
        }
    }
    return false
}

func HasRole(ctx context.Context, role string) bool {
    roles := GetRoles(ctx)
    for _, r := range roles {
        if r == role {
            return true
        }
    }
    return false
}
```

### Wiring Security to Server

```go
// cmd/server/main.go
package main

import (
    "context"
    "log"
    "net/http"
    
    goahttp "goa.design/goa/v3/http"
    api "myproject/gen/api"
    apisvr "myproject/gen/http/api/server"
    "myproject/security"
)

func main() {
    // Create stores
    userStore := NewInMemoryUserStore()
    clientStore := NewInMemoryClientStore()
    
    // Create authenticators
    apiKeyAuth := security.NewAPIKeyAuth(clientStore)
    basicAuth := security.NewBasicAuth(userStore)
    jwtAuth := security.NewJWTAuth(
        []byte("your-secret-key"),
        "myapp",
        "myapi",
    )
    
    // Create service
    svc := NewAPIService()
    endpoints := api.NewEndpoints(svc)
    
    // Create server with all security handlers
    mux := goahttp.NewMuxer()
    
    server := apisvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        errorHandler,
        // Security handlers in order matching DSL definition
        apiKeyAuth.Authenticate,    // For APIKeyAuth
        basicAuth.Authenticate,     // For BasicAuth
        jwtAuth.Authenticate,       // For JWTAuth
    )
    
    apisvr.Mount(mux, server)
    
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", mux)
}

func errorHandler(ctx context.Context, w http.ResponseWriter, err error) {
    log.Printf("Error: %v", err)
}
```

---

## 🏆 Security Best Practices

### Do's and Don'ts

```
┌─────────────────────────────────────────────────────────────────┐
│                   SECURITY BEST PRACTICES                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ✅ DO                                                          │
│  ────                                                           │
│  • Always use HTTPS in production                               │
│  • Hash passwords with bcrypt (cost >= 10)                      │
│  • Use short-lived access tokens (15-60 min)                    │
│  • Use refresh tokens for long sessions                         │
│  • Validate all token claims (iss, aud, exp)                    │
│  • Use constant-time comparison for secrets                     │
│  • Log authentication failures                                  │
│  • Rate limit authentication endpoints                          │
│  • Rotate secrets/keys periodically                             │
│  • Use scopes for fine-grained access control                   │
│                                                                 │
│  ❌ DON'T                                                       │
│  ──────                                                         │
│  • Don't store plain-text passwords                             │
│  • Don't expose sensitive data in error messages                │
│  • Don't use symmetric JWT keys in distributed systems          │
│  • Don't trust client-provided data without validation          │
│  • Don't disable certificate validation                         │
│  • Don't use weak secrets or hardcoded credentials              │
│  • Don't log sensitive information (tokens, passwords)          │
│  • Don't use deprecated OAuth2 flows (implicit)                 │
│  • Don't skip signature verification                            │
│  • Don't use MD5/SHA1 for password hashing                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Token Security Guidelines

```
┌─────────────────────────────────────────────────────────────────┐
│                   TOKEN SECURITY                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ACCESS TOKEN                                                   │
│  ────────────                                                   │
│  Lifetime: 15-60 minutes                                        │
│  Storage: Memory only (not localStorage)                        │
│  Purpose: API authorization                                     │
│  Contents: User ID, scopes, minimal claims                      │
│                                                                 │
│  REFRESH TOKEN                                                  │
│  ─────────────                                                  │
│  Lifetime: Days to weeks                                        │
│  Storage: HTTP-only secure cookie or secure storage             │
│  Purpose: Get new access tokens                                 │
│  Should be: Rotated on use, revocable                           │
│                                                                 │
│  API KEY                                                        │
│  ───────                                                        │
│  Lifetime: Long-lived, manual rotation                          │
│  Storage: Secure vault, environment variables                   │
│  Purpose: Server-to-server auth                                 │
│  Should be: Prefixed (sk_live_), easily revocable               │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Token Type     │  Lifetime    │  Revocable │  Storage  │   │
│  ├─────────────────┼──────────────┼────────────┼───────────┤   │
│  │  Access Token   │  15-60 min   │  No*       │  Memory   │   │
│  │  Refresh Token  │  7-30 days   │  Yes       │  Cookie   │   │
│  │  API Key        │  Long-lived  │  Yes       │  Server   │   │
│  │  ID Token       │  5-15 min    │  No        │  Memory   │   │
│  └─────────────────┴──────────────┴────────────┴───────────┘   │
│                                                                 │
│  * Access tokens are typically not revocable (stateless),       │
│    but can be if using token blacklist/introspection            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Error Handling in Authentication

```go
// GOOD: Generic error messages that don't leak information
func (a *BasicAuthenticator) Authenticate(ctx context.Context, username, password string) (context.Context, error) {
    user, err := a.store.GetByUsername(username)
    if err != nil {
        // Don't reveal if user exists
        return ctx, fmt.Errorf("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword(...); err != nil {
        // Same error for wrong password
        return ctx, fmt.Errorf("invalid credentials")
    }
    
    // ...
}

// BAD: Reveals information about system
func (a *BasicAuthenticator) AuthenticateBad(ctx context.Context, username, password string) (context.Context, error) {
    user, err := a.store.GetByUsername(username)
    if err != nil {
        return ctx, fmt.Errorf("user '%s' not found", username)  // BAD!
    }
    
    if err := bcrypt.CompareHashAndPassword(...); err != nil {
        return ctx, fmt.Errorf("incorrect password for '%s'", username)  // BAD!
    }
    
    // ...
}
```

---

## 📦 Complete Examples

### Complete Secure API Design

```go
// design/security.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// Security schemes
var APIKeyAuth = APIKeySecurity("api_key", func() {
    Description("API key authentication for server-to-server")
})

var BasicAuth = BasicAuthSecurity("basic", func() {
    Description("Basic authentication for simple clients")
})

var JWTAuth = JWTSecurity("jwt", func() {
    Description("JWT authentication for web/mobile clients")
    Scope("read", "Read access")
    Scope("write", "Write access")
    Scope("admin", "Administration access")
})
```

```go
// design/types.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var User = Type("User", func() {
    Attribute("id", String, "User ID")
    Attribute("username", String, "Username")
    Attribute("email", String, "Email")
    Attribute("roles", ArrayOf(String), "User roles")
    Attribute("created_at", String, "Creation timestamp")
    Required("id", "username", "email")
})

var LoginPayload = Type("LoginPayload", func() {
    Attribute("username", String, func() {
        MinLength(3)
        MaxLength(50)
    })
    Attribute("password", String, func() {
        MinLength(8)
    })
    Required("username", "password")
})

var TokenResult = Type("TokenResult", func() {
    Attribute("access_token", String, "JWT access token")
    Attribute("refresh_token", String, "Refresh token")
    Attribute("token_type", String, "Token type", func() {
        Default("Bearer")
    })
    Attribute("expires_in", Int, "Token lifetime in seconds")
    Required("access_token", "token_type", "expires_in")
})

var ErrorResult = Type("ErrorResult", func() {
    Attribute("code", String, "Error code")
    Attribute("message", String, "Error message")
    Required("code", "message")
})
```

```go
// design/auth_service.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("auth", func() {
    Description("Authentication service")
    
    Error("unauthorized", ErrorResult)
    Error("bad_request", ErrorResult)
    
    // Login - public
    Method("login", func() {
        Description("Authenticate with username/password")
        NoSecurity()
        
        Payload(LoginPayload)
        Result(TokenResult)
        Error("unauthorized")
        
        HTTP(func() {
            POST("/auth/login")
            Response(StatusOK)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    // Refresh - public (but needs valid refresh token)
    Method("refresh", func() {
        Description("Refresh access token")
        NoSecurity()
        
        Payload(func() {
            Attribute("refresh_token", String)
            Required("refresh_token")
        })
        
        Result(TokenResult)
        Error("unauthorized")
        
        HTTP(func() {
            POST("/auth/refresh")
            Response(StatusOK)
            Response("unauthorized", StatusUnauthorized)
        })
    })
    
    // Logout - authenticated
    Method("logout", func() {
        Description("Invalidate refresh token")
        Security(JWTAuth)
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("refresh_token", String)
        })
        
        HTTP(func() {
            POST("/auth/logout")
            Response(StatusNoContent)
        })
    })
    
    // Get current user - authenticated
    Method("me", func() {
        Description("Get current user info")
        Security(JWTAuth)
        
        Payload(func() {
            TokenField(1, "token", String)
        })
        
        Result(User)
        
        HTTP(func() {
            GET("/auth/me")
        })
    })
})
```

```go
// design/users_service.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("User management service")
    
    // Default security for all methods
    Security(JWTAuth)
    
    Error("not_found", ErrorResult)
    Error("unauthorized", ErrorResult)
    Error("forbidden", ErrorResult)
    Error("bad_request", ErrorResult)
    
    // List users - requires read scope
    Method("list", func() {
        Description("List all users")
        
        Security(JWTAuth, func() {
            Scope("read")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("page", Int, func() { Default(1) })
            Attribute("limit", Int, func() { Default(20) })
        })
        
        Result(func() {
            Attribute("users", ArrayOf(User))
            Attribute("total", Int)
            Required("users", "total")
        })
        
        HTTP(func() {
            GET("/users")
            Param("page")
            Param("limit")
        })
    })
    
    // Get user - requires read scope
    Method("get", func() {
        Description("Get user by ID")
        
        Security(JWTAuth, func() {
            Scope("read")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("id", String, "User ID")
            Required("token", "id")
        })
        
        Result(User)
        Error("not_found")
        
        HTTP(func() {
            GET("/users/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Create user - requires write scope
    Method("create", func() {
        Description("Create new user")
        
        Security(JWTAuth, func() {
            Scope("write")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("username", String, func() {
                MinLength(3)
                MaxLength(50)
            })
            Attribute("email", String, func() {
                Format(FormatEmail)
            })
            Attribute("password", String, func() {
                MinLength(8)
            })
            Attribute("roles", ArrayOf(String))
            Required("token", "username", "email", "password")
        })
        
        Result(User)
        Error("bad_request")
        
        HTTP(func() {
            POST("/users")
            Response(StatusCreated)
            Response("bad_request", StatusBadRequest)
        })
    })
    
    // Update user - requires write scope
    Method("update", func() {
        Description("Update user")
        
        Security(JWTAuth, func() {
            Scope("write")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("id", String)
            Attribute("email", String, func() {
                Format(FormatEmail)
            })
            Attribute("roles", ArrayOf(String))
            Required("token", "id")
        })
        
        Result(User)
        Error("not_found")
        
        HTTP(func() {
            PUT("/users/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    // Delete user - requires admin scope
    Method("delete", func() {
        Description("Delete user")
        
        Security(JWTAuth, func() {
            Scope("admin")
        })
        
        Payload(func() {
            TokenField(1, "token", String)
            Attribute("id", String)
            Required("token", "id")
        })
        
        Error("not_found")
        Error("forbidden")
        
        HTTP(func() {
            DELETE("/users/{id}")
            Response(StatusNoContent)
            Response("not_found", StatusNotFound)
            Response("forbidden", StatusForbidden)
        })
    })
})
```

```go
// design/webhook_service.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("webhooks", func() {
    Description("Webhook service for integrations")
    
    // API Key auth for webhooks
    Security(APIKeyAuth)
    
    Error("unauthorized", ErrorResult)
    Error("bad_request", ErrorResult)
    
    Method("receive", func() {
        Description("Receive webhook event")
        
        Payload(func() {
            APIKeyField(1, "api_key", String, func() {
                Description("Webhook API key")
            })
            Attribute("event_type", String, func() {
                Enum("user.created", "user.updated", "user.deleted",
                     "order.created", "order.completed")
            })
            Attribute("data", Any, "Event data")
            Attribute("timestamp", Int64, "Event timestamp")
            Required("api_key", "event_type", "data")
        })
        
        Result(func() {
            Attribute("received", Boolean)
            Attribute("event_id", String)
            Required("received", "event_id")
        })
        
        HTTP(func() {
            POST("/webhooks")
            Header("api_key:X-Webhook-Key")
            Response(StatusOK)
            Response("unauthorized", StatusUnauthorized)
            Response("bad_request", StatusBadRequest)
        })
    })
})
```

### Complete Server Main

```go
// cmd/server/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    goahttp "goa.design/goa/v3/http"
    
    auth "myproject/gen/auth"
    users "myproject/gen/users"
    webhooks "myproject/gen/webhooks"
    
    authsvr "myproject/gen/http/auth/server"
    userssvr "myproject/gen/http/users/server"
    webhookssvr "myproject/gen/http/webhooks/server"
    
    "myproject/security"
)

func main() {
    // Configuration
    jwtConfig := &security.JWTConfig{
        SecretKey:     []byte(getEnv("JWT_SECRET", "your-256-bit-secret")),
        Issuer:        getEnv("JWT_ISSUER", "myapp"),
        Audience:      getEnv("JWT_AUDIENCE", "myapi"),
        TokenDuration: 15 * time.Minute,
    }
    
    // Create stores
    userStore := security.NewInMemoryUserStore()
    clientStore := security.NewInMemoryClientStore()
    
    // Seed test data
    seedTestData(userStore, clientStore)
    
    // Create authenticators
    apiKeyAuth := security.NewAPIKeyAuth(clientStore)
    jwtAuth := security.NewJWTAuth(jwtConfig.SecretKey, jwtConfig.Issuer, jwtConfig.Audience)
    
    // Create services
    authSvc := NewAuthService(jwtConfig, userStore)
    usersSvc := NewUsersService(userStore)
    webhooksSvc := NewWebhooksService()
    
    // Create endpoints
    authEndpoints := auth.NewEndpoints(authSvc)
    usersEndpoints := users.NewEndpoints(usersSvc)
    webhooksEndpoints := webhooks.NewEndpoints(webhooksSvc)
    
    // Create HTTP mux
    mux := goahttp.NewMuxer()
    
    // Error handler
    errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
        log.Printf("ERROR: %v", err)
    }
    
    // Mount auth service (no security handlers needed - all public or JWT)
    authServer := authsvr.New(
        authEndpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        errHandler,
        jwtAuth.Authenticate,  // For /auth/me and /auth/logout
    )
    authsvr.Mount(mux, authServer)
    
    // Mount users service (JWT auth)
    usersServer := userssvr.New(
        usersEndpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        errHandler,
        jwtAuth.Authenticate,
    )
    userssvr.Mount(mux, usersServer)
    
    // Mount webhooks service (API key auth)
    webhooksServer := webhookssvr.New(
        webhooksEndpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        errHandler,
        apiKeyAuth.Authenticate,
    )
    webhookssvr.Mount(mux, webhooksServer)
    
    // Create HTTP server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    
    // Start server
    go func() {
        log.Printf("HTTP server listening on %s", server.Addr)
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()
    
    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // Graceful shutdown
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    server.Shutdown(ctx)
    log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func seedTestData(users *security.InMemoryUserStore, clients *security.InMemoryClientStore) {
    // Add test user
    users.Add(&security.User{
        ID:           "user-1",
        Username:     "admin",
        Email:        "admin@example.com",
        PasswordHash: security.MustHashPassword("admin123"),
        Roles:        []string{"admin"},
        Active:       true,
    })
    
    users.Add(&security.User{
        ID:           "user-2",
        Username:     "john",
        Email:        "john@example.com",
        PasswordHash: security.MustHashPassword("password123"),
        Roles:        []string{"user"},
        Active:       true,
    })
    
    // Add test API client
    clients.Add(&security.Client{
        ID:          "client-1",
        Name:        "Webhook Client",
        APIKey:      "whk_test_abc123def456",
        Permissions: []string{"webhook:receive"},
        Active:      true,
    })
}
```

---

## 📝 Summary

### Security Schemes
- **APIKeySecurity**: Simple key-based auth, best for server-to-server
- **BasicAuthSecurity**: Username/password, use only with HTTPS
- **JWTSecurity**: Stateless tokens with claims and scopes
- **OAuth2Security**: Delegated authorization with multiple flows

### Applying Security
- **Service-level**: `Security(Scheme)` applies to all methods
- **Method-level**: Override or add specific security
- **NoSecurity()**: Make endpoints public
- **Multiple schemes**: OR relationship (any scheme works)
- **Scopes**: Fine-grained permissions within a scheme

### Token Fields
- **TokenField**: General token (API Key, JWT)
- **APIKeyField**: Explicitly for API keys
- **UsernameField / PasswordField**: For Basic Auth
- **AccessTokenField**: For OAuth2

### Implementation
- Security handlers receive credentials, return context
- Store authenticated entity in context
- Return errors for auth failures (401, 403)
- Use constant-time comparison for secrets

---

## 📋 Knowledge Check

Before proceeding, ensure you can:

- [ ] Define APIKeySecurity, BasicAuthSecurity, JWTSecurity, and OAuth2Security schemes
- [ ] Apply security at service and method levels
- [ ] Use NoSecurity() for public endpoints
- [ ] Map tokens to headers, query params, and cookies
- [ ] Define and require scopes for methods
- [ ] Implement API key validation handler
- [ ] Implement Basic Auth with bcrypt password verification
- [ ] Implement JWT validation with claims extraction
- [ ] Implement OAuth2 token introspection
- [ ] Store and retrieve auth info from context
- [ ] Wire security handlers to generated server
- [ ] Apply security best practices

---

## 🔗 Quick Reference Links

- [Goa Security DSL](https://pkg.go.dev/goa.design/goa/v3/dsl#Security)
- [Goa Security Examples](https://github.com/goadesign/examples/tree/master/security)
- [JWT.io](https://jwt.io/) - JWT debugger
- [golang-jwt Library](https://github.com/golang-jwt/jwt)
- [OAuth2 Specification](https://oauth.net/2/)
- [OWASP Authentication Cheatsheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)

---

> **Next Up:** Part 6 - Middleware & Plugins (Logging, Tracing, Rate Limiting, Custom Middleware)
