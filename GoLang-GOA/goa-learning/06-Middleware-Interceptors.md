# Part 6: Middleware & Interceptors

> **Goal:** Master Goa's middleware system - HTTP middleware, gRPC interceptors, logging, request tracing, and building custom middleware

---

## 📚 Table of Contents

1. [Middleware Overview](#middleware-overview)
2. [Goa Middleware Architecture](#goa-middleware-architecture)
3. [Endpoint Middleware](#endpoint-middleware)
4. [HTTP Middleware](#http-middleware)
5. [Logging Middleware](#logging-middleware)
6. [Request ID Middleware](#request-id-middleware)
7. [Custom Middleware](#custom-middleware)
8. [gRPC Interceptors](#grpc-interceptors)
9. [Middleware Patterns](#middleware-patterns)
10. [Complete Examples](#complete-examples)
11. [Summary](#summary)
12. [Knowledge Check](#knowledge-check)

---

## 🎯 Middleware Overview

### What is Middleware?

Middleware is code that runs between receiving a request and executing your business logic, and between your business logic completing and sending a response.

```
┌─────────────────────────────────────────────────────────────────┐
│                    REQUEST/RESPONSE FLOW                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                        REQUEST                                  │
│                           │                                     │
│                           ▼                                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   MIDDLEWARE 1                           │   │
│  │              (e.g., Request ID)                          │   │
│  │    ┌─────────────────────────────────────────────────┐  │   │
│  │    │                MIDDLEWARE 2                      │  │   │
│  │    │            (e.g., Logging)                       │  │   │
│  │    │    ┌─────────────────────────────────────────┐  │  │   │
│  │    │    │             MIDDLEWARE 3                │  │  │   │
│  │    │    │         (e.g., Auth)                    │  │  │   │
│  │    │    │    ┌─────────────────────────────────┐  │  │  │   │
│  │    │    │    │                                 │  │  │  │   │
│  │    │    │    │      YOUR BUSINESS LOGIC       │  │  │  │   │
│  │    │    │    │                                 │  │  │  │   │
│  │    │    │    └─────────────────────────────────┘  │  │  │   │
│  │    │    │                                         │  │  │   │
│  │    │    └─────────────────────────────────────────┘  │  │   │
│  │    │                                                 │  │   │
│  │    └─────────────────────────────────────────────────┘  │   │
│  │                                                         │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│                           ▼                                     │
│                        RESPONSE                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Why Use Middleware?

| Purpose | Example |
|---------|---------|
| **Cross-cutting concerns** | Logging, metrics, tracing |
| **Request processing** | Authentication, validation, rate limiting |
| **Response processing** | Compression, caching headers |
| **Error handling** | Panic recovery, error transformation |
| **Context enrichment** | Request ID, user info, tenant ID |

### Middleware Types in Goa

```
┌─────────────────────────────────────────────────────────────────┐
│                   GOA MIDDLEWARE LAYERS                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. TRANSPORT MIDDLEWARE (HTTP/gRPC specific)                   │
│     ─────────────────────────────────────────                   │
│     • HTTP: http.Handler wrapper                                │
│     • gRPC: UnaryInterceptor, StreamInterceptor                 │
│     • Access to raw request/response                            │
│     • Protocol-specific operations                              │
│                                                                 │
│  2. ENDPOINT MIDDLEWARE (Transport agnostic)                    │
│     ───────────────────────────────────────                     │
│     • Works with Goa endpoints                                  │
│     • Access to decoded payload and result                      │
│     • Same code for HTTP and gRPC                               │
│     • Business logic wrapping                                   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │           HTTP Request                gRPC Request       │   │
│  │                │                           │             │   │
│  │                ▼                           ▼             │   │
│  │     ┌──────────────────┐       ┌──────────────────┐     │   │
│  │     │  HTTP Middleware │       │ gRPC Interceptor │     │   │
│  │     │  (net/http)      │       │ (grpc)           │     │   │
│  │     └────────┬─────────┘       └────────┬─────────┘     │   │
│  │              │                          │               │   │
│  │              └──────────┬───────────────┘               │   │
│  │                         ▼                               │   │
│  │              ┌──────────────────┐                       │   │
│  │              │ Endpoint         │                       │   │
│  │              │ Middleware       │                       │   │
│  │              │ (goa.Endpoint)   │                       │   │
│  │              └────────┬─────────┘                       │   │
│  │                       ▼                                 │   │
│  │              ┌──────────────────┐                       │   │
│  │              │ Service Method   │                       │   │
│  │              └──────────────────┘                       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🏗️ Goa Middleware Architecture

### Understanding Goa Endpoints

In Goa, an **endpoint** is a function that takes a context and request, and returns a response and error:

```go
// Goa endpoint signature
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
```

### Endpoint Middleware Signature

Middleware wraps endpoints:

```go
// Middleware wraps an endpoint and returns a new endpoint
type Middleware func(Endpoint) Endpoint
```

### How Middleware Works

```go
// Simple middleware example
func LoggingMiddleware(logger *log.Logger) func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // BEFORE: runs before the endpoint
            logger.Printf("Request received: %T", req)
            start := time.Now()
            
            // Call the next middleware/endpoint
            res, err := next(ctx, req)
            
            // AFTER: runs after the endpoint
            logger.Printf("Request completed in %v", time.Since(start))
            
            return res, err
        }
    }
}
```

### Middleware Execution Order

```
┌─────────────────────────────────────────────────────────────────┐
│                  MIDDLEWARE EXECUTION ORDER                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Middleware added:  [M1, M2, M3]                                │
│                                                                 │
│  Execution:                                                     │
│                                                                 │
│  REQUEST  ──▶  M1.Before ──▶  M2.Before ──▶  M3.Before ──▶     │
│                                                                 │
│                              ENDPOINT                           │
│                                                                 │
│  RESPONSE ◀──  M1.After  ◀──  M2.After  ◀──  M3.After  ◀──     │
│                                                                 │
│  Example with Logging, Auth, Metrics:                           │
│                                                                 │
│  REQUEST:                                                       │
│    1. Logging.Before  (log: "request received")                 │
│    2. Auth.Before     (validate token)                          │
│    3. Metrics.Before  (start timer)                             │
│    4. ENDPOINT        (business logic)                          │
│                                                                 │
│  RESPONSE:                                                      │
│    5. Metrics.After   (record latency)                          │
│    6. Auth.After      (nothing)                                 │
│    7. Logging.After   (log: "request completed")                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔌 Endpoint Middleware

### Creating Endpoint Middleware

```go
// middleware/endpoint.go
package middleware

import (
    "context"
    "time"
    
    "goa.design/goa/v3/middleware"
    goa "goa.design/goa/v3/pkg"
)

// TimingMiddleware measures endpoint execution time
func TimingMiddleware() func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            start := time.Now()
            
            res, err := next(ctx, req)
            
            duration := time.Since(start)
            
            // Log or record metrics
            if err != nil {
                log.Printf("Endpoint failed after %v: %v", duration, err)
            } else {
                log.Printf("Endpoint succeeded in %v", duration)
            }
            
            return res, err
        }
    }
}
```

### Applying Endpoint Middleware

```go
// cmd/server/main.go
package main

import (
    "context"
    
    goa "goa.design/goa/v3/pkg"
    users "myproject/gen/users"
    "myproject/middleware"
)

func main() {
    // Create service
    svc := NewUsersService()
    
    // Create endpoints
    endpoints := users.NewEndpoints(svc)
    
    // Apply middleware to all endpoints
    endpoints.Use(middleware.TimingMiddleware())
    endpoints.Use(middleware.LoggingMiddleware(logger))
    
    // Or apply to specific endpoint
    endpoints.Get = middleware.TimingMiddleware()(endpoints.Get)
    
    // ... rest of server setup
}
```

### Built-in Goa Middleware

Goa provides several built-in middleware in `goa.design/goa/v3/middleware`:

```go
import "goa.design/goa/v3/middleware"

// Request ID middleware
endpoints.Use(middleware.RequestID())

// Logging middleware (requires log adapter)
endpoints.Use(middleware.Log(adapter))

// Debug middleware (logs request/response payloads)
endpoints.Use(middleware.Debug(mux, os.Stdout))
```

### Endpoint Middleware with Context

```go
// middleware/context.go
package middleware

import (
    "context"
    
    goa "goa.design/goa/v3/pkg"
)

// Context keys
type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
    TenantIDKey  contextKey = "tenant_id"
)

// ContextEnricher adds values to context
func ContextEnricher(tenantID string) func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // Add tenant ID to context
            ctx = context.WithValue(ctx, TenantIDKey, tenantID)
            
            return next(ctx, req)
        }
    }
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(ctx context.Context) string {
    if v := ctx.Value(TenantIDKey); v != nil {
        return v.(string)
    }
    return ""
}
```

---

## 🌐 HTTP Middleware

### Standard HTTP Middleware Pattern

HTTP middleware in Go follows the `http.Handler` wrapper pattern:

```go
// Standard HTTP middleware signature
type HTTPMiddleware func(http.Handler) http.Handler
```

### Creating HTTP Middleware

```go
// middleware/http.go
package middleware

import (
    "net/http"
    "time"
)

// HTTPLogging logs HTTP requests
func HTTPLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create response wrapper to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        // Call next handler
        next.ServeHTTP(wrapped, r)
        
        // Log after
        log.Printf(
            "%s %s %d %v",
            r.Method,
            r.URL.Path,
            wrapped.statusCode,
            time.Since(start),
        )
    })
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### Applying HTTP Middleware to Goa

```go
// cmd/server/main.go
package main

import (
    "net/http"
    
    goahttp "goa.design/goa/v3/http"
    "myproject/middleware"
)

func main() {
    // Create Goa mux
    mux := goahttp.NewMuxer()
    
    // Mount services...
    
    // Wrap mux with HTTP middleware
    handler := middleware.HTTPLogging(mux)
    handler = middleware.Recovery(handler)
    handler = middleware.CORS(handler)
    
    // Start server
    http.ListenAndServe(":8080", handler)
}
```

### Common HTTP Middleware Examples

#### Recovery Middleware

```go
// middleware/recovery.go
package middleware

import (
    "net/http"
    "runtime/debug"
)

// Recovery recovers from panics and returns 500
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Log the stack trace
                log.Printf("PANIC: %v\n%s", err, debug.Stack())
                
                // Return 500 error
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(`{"error":"internal server error"}`))
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

#### CORS Middleware

```go
// middleware/cors.go
package middleware

import (
    "net/http"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
    AllowOrigins     []string
    AllowMethods     []string
    AllowHeaders     []string
    ExposeHeaders    []string
    AllowCredentials bool
    MaxAge           int
}

// DefaultCORSConfig returns default CORS config
func DefaultCORSConfig() CORSConfig {
    return CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
        MaxAge:       86400,
    }
}

// CORS creates CORS middleware
func CORS(config CORSConfig) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if origin is allowed
            allowed := false
            for _, o := range config.AllowOrigins {
                if o == "*" || o == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", 
                    strings.Join(config.AllowMethods, ", "))
                w.Header().Set("Access-Control-Allow-Headers", 
                    strings.Join(config.AllowHeaders, ", "))
                
                if config.AllowCredentials {
                    w.Header().Set("Access-Control-Allow-Credentials", "true")
                }
                
                if len(config.ExposeHeaders) > 0 {
                    w.Header().Set("Access-Control-Expose-Headers",
                        strings.Join(config.ExposeHeaders, ", "))
                }
                
                if config.MaxAge > 0 {
                    w.Header().Set("Access-Control-Max-Age", 
                        strconv.Itoa(config.MaxAge))
                }
            }
            
            // Handle preflight
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

#### Compression Middleware

```go
// middleware/compress.go
package middleware

import (
    "compress/gzip"
    "io"
    "net/http"
    "strings"
)

// Gzip compresses responses
func Gzip(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if client accepts gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        // Create gzip writer
        gz := gzip.NewWriter(w)
        defer gz.Close()
        
        // Set headers
        w.Header().Set("Content-Encoding", "gzip")
        w.Header().Del("Content-Length")
        
        // Wrap response writer
        gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
        
        next.ServeHTTP(gzw, r)
    })
}

type gzipResponseWriter struct {
    http.ResponseWriter
    io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}
```

#### Rate Limiting Middleware

```go
// middleware/ratelimit.go
package middleware

import (
    "net/http"
    "sync"
    "time"
)

// RateLimiter implements a simple rate limiter
type RateLimiter struct {
    mu       sync.Mutex
    requests map[string][]time.Time
    limit    int
    window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

// RateLimit creates rate limiting middleware
func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get client identifier (IP address)
        clientIP := getClientIP(r)
        
        rl.mu.Lock()
        
        // Clean old requests
        now := time.Now()
        cutoff := now.Add(-rl.window)
        
        var recent []time.Time
        for _, t := range rl.requests[clientIP] {
            if t.After(cutoff) {
                recent = append(recent, t)
            }
        }
        
        // Check if limit exceeded
        if len(recent) >= rl.limit {
            rl.mu.Unlock()
            
            w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
            w.Header().Set("X-RateLimit-Remaining", "0")
            w.Header().Set("Retry-After", "60")
            
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        // Record request
        recent = append(recent, now)
        rl.requests[clientIP] = recent
        
        remaining := rl.limit - len(recent)
        rl.mu.Unlock()
        
        // Set rate limit headers
        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
        w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
        
        next.ServeHTTP(w, r)
    })
}

func getClientIP(r *http.Request) string {
    // Check X-Forwarded-For header
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        parts := strings.Split(xff, ",")
        return strings.TrimSpace(parts[0])
    }
    
    // Check X-Real-IP header
    if xri := r.Header.Get("X-Real-IP"); xri != "" {
        return xri
    }
    
    // Fall back to remote address
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}
```

---

## 📝 Logging Middleware

### Structured Logging Middleware

```go
// middleware/logging.go
package middleware

import (
    "context"
    "log/slog"
    "net/http"
    "time"
    
    goa "goa.design/goa/v3/pkg"
)

// Logger context key
type loggerKey struct{}

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger retrieves logger from context
func GetLogger(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return l
    }
    return slog.Default()
}

// EndpointLogging creates endpoint logging middleware
func EndpointLogging(logger *slog.Logger) func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // Add request ID to logger
            requestID := GetRequestID(ctx)
            reqLogger := logger.With(
                slog.String("request_id", requestID),
            )
            
            // Add logger to context
            ctx = WithLogger(ctx, reqLogger)
            
            // Log request
            reqLogger.Info("endpoint request",
                slog.String("type", fmt.Sprintf("%T", req)),
            )
            
            start := time.Now()
            
            // Call endpoint
            res, err := next(ctx, req)
            
            duration := time.Since(start)
            
            // Log response
            if err != nil {
                reqLogger.Error("endpoint error",
                    slog.Duration("duration", duration),
                    slog.String("error", err.Error()),
                )
            } else {
                reqLogger.Info("endpoint success",
                    slog.Duration("duration", duration),
                    slog.String("result_type", fmt.Sprintf("%T", res)),
                )
            }
            
            return res, err
        }
    }
}

// HTTPLogging creates HTTP logging middleware
func HTTPLogging(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer
            wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: 200}
            
            // Get or generate request ID
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = generateRequestID()
            }
            
            // Add to response header
            w.Header().Set("X-Request-ID", requestID)
            
            // Create request-scoped logger
            reqLogger := logger.With(
                slog.String("request_id", requestID),
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.String("remote_addr", r.RemoteAddr),
            )
            
            // Log request start
            reqLogger.Info("http request started",
                slog.String("user_agent", r.UserAgent()),
            )
            
            // Call next handler
            next.ServeHTTP(wrapped, r)
            
            // Log request completion
            reqLogger.Info("http request completed",
                slog.Int("status", wrapped.statusCode),
                slog.Duration("duration", time.Since(start)),
                slog.Int("bytes", wrapped.bytesWritten),
            )
        })
    }
}

type statusResponseWriter struct {
    http.ResponseWriter
    statusCode   int
    bytesWritten int
}

func (w *statusResponseWriter) WriteHeader(code int) {
    w.statusCode = code
    w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.bytesWritten += n
    return n, err
}
```

### Log Levels and Filtering

```go
// middleware/loglevel.go
package middleware

import (
    "context"
    "log/slog"
    "net/http"
    
    goa "goa.design/goa/v3/pkg"
)

// LogLevel middleware allows dynamic log levels
type LogLevelConfig struct {
    // Paths to log at debug level
    DebugPaths map[string]bool
    
    // Paths to skip logging entirely
    SkipPaths map[string]bool
    
    // Log request body (use with caution)
    LogRequestBody bool
    
    // Log response body (use with caution)
    LogResponseBody bool
}

// ConditionalLogging logs based on configuration
func ConditionalLogging(logger *slog.Logger, config LogLevelConfig) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Skip logging for certain paths
            if config.SkipPaths[r.URL.Path] {
                next.ServeHTTP(w, r)
                return
            }
            
            // Determine log level
            level := slog.LevelInfo
            if config.DebugPaths[r.URL.Path] {
                level = slog.LevelDebug
            }
            
            start := time.Now()
            wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: 200}
            
            // Log at appropriate level
            logger.Log(r.Context(), level, "http request",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
            )
            
            next.ServeHTTP(wrapped, r)
            
            logger.Log(r.Context(), level, "http response",
                slog.Int("status", wrapped.statusCode),
                slog.Duration("duration", time.Since(start)),
            )
        })
    }
}
```

### JSON Logging

```go
// middleware/jsonlog.go
package middleware

import (
    "encoding/json"
    "io"
    "net/http"
    "time"
)

// LogEntry represents a structured log entry
type LogEntry struct {
    Timestamp   string            `json:"timestamp"`
    Level       string            `json:"level"`
    RequestID   string            `json:"request_id"`
    Method      string            `json:"method"`
    Path        string            `json:"path"`
    Status      int               `json:"status,omitempty"`
    Duration    float64           `json:"duration_ms,omitempty"`
    RemoteAddr  string            `json:"remote_addr"`
    UserAgent   string            `json:"user_agent"`
    Error       string            `json:"error,omitempty"`
    Extra       map[string]string `json:"extra,omitempty"`
}

// JSONLogger writes JSON log entries
type JSONLogger struct {
    writer io.Writer
}

// NewJSONLogger creates a new JSON logger
func NewJSONLogger(w io.Writer) *JSONLogger {
    return &JSONLogger{writer: w}
}

// Middleware creates HTTP logging middleware
func (l *JSONLogger) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: 200}
        
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
            w.Header().Set("X-Request-ID", requestID)
        }
        
        // Call next handler
        next.ServeHTTP(wrapped, r)
        
        // Create log entry
        entry := LogEntry{
            Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
            Level:      "INFO",
            RequestID:  requestID,
            Method:     r.Method,
            Path:       r.URL.Path,
            Status:     wrapped.statusCode,
            Duration:   float64(time.Since(start).Microseconds()) / 1000,
            RemoteAddr: getClientIP(r),
            UserAgent:  r.UserAgent(),
        }
        
        // Adjust level based on status
        if wrapped.statusCode >= 500 {
            entry.Level = "ERROR"
        } else if wrapped.statusCode >= 400 {
            entry.Level = "WARN"
        }
        
        // Write JSON log
        json.NewEncoder(l.writer).Encode(entry)
    })
}
```

---

## 🏷️ Request ID Middleware

### What is Request ID?

Request ID (also called Correlation ID or Trace ID) is a unique identifier that follows a request through all services and logs, enabling distributed tracing.

```
┌─────────────────────────────────────────────────────────────────┐
│                    REQUEST ID FLOW                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Client Request                                                 │
│       │                                                         │
│       │  X-Request-ID: (none or client-provided)                │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  API Gateway / Load Balancer                             │   │
│  │  Generate: X-Request-ID: abc-123-def-456                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│       │                                                         │
│       │  X-Request-ID: abc-123-def-456                          │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Service A                                               │   │
│  │  Log: [abc-123-def-456] Processing request               │   │
│  └─────────────────────────────────────────────────────────┘   │
│       │                                                         │
│       │  X-Request-ID: abc-123-def-456                          │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Service B                                               │   │
│  │  Log: [abc-123-def-456] Fetching data                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│       │                                                         │
│       │  X-Request-ID: abc-123-def-456                          │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Database / External Service                             │   │
│  │  Log: [abc-123-def-456] Query executed                   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  All logs can be correlated by: abc-123-def-456                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Request ID Implementation

```go
// middleware/requestid.go
package middleware

import (
    "context"
    "net/http"
    
    "github.com/google/uuid"
    goa "goa.design/goa/v3/pkg"
)

// Request ID header names
const (
    RequestIDHeader     = "X-Request-ID"
    CorrelationIDHeader = "X-Correlation-ID"
)

// Context key for request ID
type requestIDKey struct{}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, id)
}

// GetRequestID retrieves request ID from context
func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey{}).(string); ok {
        return id
    }
    return ""
}

// generateRequestID creates a new unique request ID
func generateRequestID() string {
    return uuid.New().String()
}

// HTTPRequestID creates HTTP request ID middleware
func HTTPRequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get existing request ID or generate new one
        requestID := r.Header.Get(RequestIDHeader)
        if requestID == "" {
            requestID = r.Header.Get(CorrelationIDHeader)
        }
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // Add to request context
        ctx := WithRequestID(r.Context(), requestID)
        
        // Add to response headers
        w.Header().Set(RequestIDHeader, requestID)
        
        // Call next with updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// EndpointRequestID creates endpoint request ID middleware
// Use this when HTTP middleware already set the ID
func EndpointRequestID() func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // Check if request ID already exists
            if GetRequestID(ctx) == "" {
                // Generate new one if not
                ctx = WithRequestID(ctx, generateRequestID())
            }
            
            return next(ctx, req)
        }
    }
}
```

### Request ID with Goa's Built-in Middleware

```go
// Using Goa's built-in request ID middleware
import (
    "goa.design/goa/v3/middleware"
)

// In server setup
endpoints.Use(middleware.RequestID())

// Access in service
func (s *svc) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    // Get request ID from context
    reqID := middleware.ContextRequestID(ctx)
    
    s.logger.Info("processing request", "request_id", reqID)
    
    // ...
}
```

### Propagating Request ID to Other Services

```go
// middleware/propagate.go
package middleware

import (
    "context"
    "net/http"
)

// PropagatingHTTPClient wraps http.Client to propagate request ID
type PropagatingHTTPClient struct {
    client *http.Client
}

// NewPropagatingClient creates a new propagating client
func NewPropagatingClient(client *http.Client) *PropagatingHTTPClient {
    if client == nil {
        client = http.DefaultClient
    }
    return &PropagatingHTTPClient{client: client}
}

// Do executes request with request ID propagation
func (c *PropagatingHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
    // Get request ID from context
    requestID := GetRequestID(ctx)
    if requestID != "" {
        req.Header.Set(RequestIDHeader, requestID)
    }
    
    // Forward context
    req = req.WithContext(ctx)
    
    return c.client.Do(req)
}

// Get is a convenience method for GET requests
func (c *PropagatingHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    return c.Do(ctx, req)
}
```

---

## 🛠️ Custom Middleware

### Authentication Middleware

```go
// middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    
    goa "goa.design/goa/v3/pkg"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
    // Header name for the token
    HeaderName string
    
    // Token prefix (e.g., "Bearer ")
    TokenPrefix string
    
    // Validator function
    Validator func(ctx context.Context, token string) (context.Context, error)
    
    // Paths to skip authentication
    SkipPaths map[string]bool
}

// HTTPAuth creates HTTP authentication middleware
func HTTPAuth(config AuthConfig) func(http.Handler) http.Handler {
    if config.HeaderName == "" {
        config.HeaderName = "Authorization"
    }
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Skip authentication for certain paths
            if config.SkipPaths[r.URL.Path] {
                next.ServeHTTP(w, r)
                return
            }
            
            // Get token from header
            authHeader := r.Header.Get(config.HeaderName)
            if authHeader == "" {
                http.Error(w, "Authorization required", http.StatusUnauthorized)
                return
            }
            
            // Remove prefix
            token := strings.TrimPrefix(authHeader, config.TokenPrefix)
            
            // Validate token
            ctx, err := config.Validator(r.Context(), token)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // Continue with enriched context
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// EndpointAuth creates endpoint authentication middleware
func EndpointAuth(validator func(ctx context.Context, token string) (context.Context, error)) func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // Extract token from payload
            token := extractToken(req)
            if token == "" {
                return nil, fmt.Errorf("authorization required")
            }
            
            // Validate
            ctx, err := validator(ctx, token)
            if err != nil {
                return nil, err
            }
            
            return next(ctx, req)
        }
    }
}

// extractToken extracts token from various payload types
func extractToken(req interface{}) string {
    // Use reflection or type assertions
    // This is an example - actual implementation depends on your payload structure
    if p, ok := req.(interface{ GetToken() string }); ok {
        return p.GetToken()
    }
    return ""
}
```

### Metrics Middleware

```go
// middleware/metrics.go
package middleware

import (
    "context"
    "net/http"
    "strconv"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    goa "goa.design/goa/v3/pkg"
)

// Prometheus metrics
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    
    endpointRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "endpoint_requests_total",
            Help: "Total number of endpoint requests",
        },
        []string{"service", "method", "status"},
    )
    
    endpointRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "endpoint_request_duration_seconds",
            Help:    "Endpoint request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"service", "method"},
    )
)

// HTTPMetrics creates HTTP metrics middleware
func HTTPMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start).Seconds()
        
        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            strconv.Itoa(wrapped.statusCode),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration)
    })
}

// EndpointMetrics creates endpoint metrics middleware
func EndpointMetrics(serviceName, methodName string) func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            start := time.Now()
            
            res, err := next(ctx, req)
            
            duration := time.Since(start).Seconds()
            status := "success"
            if err != nil {
                status = "error"
            }
            
            endpointRequestsTotal.WithLabelValues(
                serviceName,
                methodName,
                status,
            ).Inc()
            
            endpointRequestDuration.WithLabelValues(
                serviceName,
                methodName,
            ).Observe(duration)
            
            return res, err
        }
    }
}
```

### Caching Middleware

```go
// middleware/cache.go
package middleware

import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "sync"
    "time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
    Body       []byte
    Headers    http.Header
    StatusCode int
    ExpiresAt  time.Time
}

// Cache is a simple in-memory cache
type Cache struct {
    mu      sync.RWMutex
    entries map[string]*CacheEntry
    ttl     time.Duration
}

// NewCache creates a new cache
func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        entries: make(map[string]*CacheEntry),
        ttl:     ttl,
    }
    
    // Start cleanup goroutine
    go c.cleanup()
    
    return c
}

// Middleware creates caching middleware
func (c *Cache) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only cache GET requests
        if r.Method != http.MethodGet {
            next.ServeHTTP(w, r)
            return
        }
        
        // Generate cache key
        key := c.generateKey(r)
        
        // Check cache
        c.mu.RLock()
        entry, found := c.entries[key]
        c.mu.RUnlock()
        
        if found && time.Now().Before(entry.ExpiresAt) {
            // Return cached response
            for k, v := range entry.Headers {
                w.Header()[k] = v
            }
            w.Header().Set("X-Cache", "HIT")
            w.WriteHeader(entry.StatusCode)
            w.Write(entry.Body)
            return
        }
        
        // Cache miss - call handler and cache response
        rec := &responseRecorder{
            ResponseWriter: w,
            statusCode:     200,
            body:           &bytes.Buffer{},
        }
        
        next.ServeHTTP(rec, r)
        
        // Only cache successful responses
        if rec.statusCode >= 200 && rec.statusCode < 300 {
            c.mu.Lock()
            c.entries[key] = &CacheEntry{
                Body:       rec.body.Bytes(),
                Headers:    rec.Header().Clone(),
                StatusCode: rec.statusCode,
                ExpiresAt:  time.Now().Add(c.ttl),
            }
            c.mu.Unlock()
        }
        
        w.Header().Set("X-Cache", "MISS")
    })
}

func (c *Cache) generateKey(r *http.Request) string {
    // Create key from method, path, and query
    h := sha256.New()
    h.Write([]byte(r.Method))
    h.Write([]byte(r.URL.Path))
    h.Write([]byte(r.URL.RawQuery))
    return hex.EncodeToString(h.Sum(nil))
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, entry := range c.entries {
            if now.After(entry.ExpiresAt) {
                delete(c.entries, key)
            }
        }
        c.mu.Unlock()
    }
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
    body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
    r.statusCode = code
    r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
    r.body.Write(b)
    return r.ResponseWriter.Write(b)
}
```

### Timeout Middleware

```go
// middleware/timeout.go
package middleware

import (
    "context"
    "net/http"
    "time"
    
    goa "goa.design/goa/v3/pkg"
)

// HTTPTimeout creates HTTP timeout middleware
func HTTPTimeout(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            
            // Create channel for completion
            done := make(chan struct{})
            
            go func() {
                next.ServeHTTP(w, r.WithContext(ctx))
                close(done)
            }()
            
            select {
            case <-done:
                // Request completed
            case <-ctx.Done():
                // Timeout
                if ctx.Err() == context.DeadlineExceeded {
                    http.Error(w, "Request timeout", http.StatusGatewayTimeout)
                }
            }
        })
    }
}

// EndpointTimeout creates endpoint timeout middleware
func EndpointTimeout(timeout time.Duration) func(goa.Endpoint) goa.Endpoint {
    return func(next goa.Endpoint) goa.Endpoint {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            ctx, cancel := context.WithTimeout(ctx, timeout)
            defer cancel()
            
            // Channel for result
            type result struct {
                res interface{}
                err error
            }
            ch := make(chan result, 1)
            
            go func() {
                res, err := next(ctx, req)
                ch <- result{res, err}
            }()
            
            select {
            case r := <-ch:
                return r.res, r.err
            case <-ctx.Done():
                return nil, fmt.Errorf("request timeout: %w", ctx.Err())
            }
        }
    }
}
```

---

## 📡 gRPC Interceptors

### Understanding gRPC Interceptors

gRPC interceptors are the gRPC equivalent of HTTP middleware. There are two types:

```
┌─────────────────────────────────────────────────────────────────┐
│                    gRPC INTERCEPTOR TYPES                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  UNARY INTERCEPTOR                                              │
│  ─────────────────                                              │
│  For single request/response calls                              │
│                                                                 │
│  Client ──── Request ────▶ Server                               │
│  Client ◀─── Response ──── Server                               │
│                                                                 │
│  STREAM INTERCEPTOR                                             │
│  ──────────────────                                             │
│  For streaming calls (server, client, or bidirectional)         │
│                                                                 │
│  Client ◀────── Stream ──────▶ Server                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Unary Interceptor Signature

```go
// Server-side unary interceptor
type UnaryServerInterceptor func(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error)

// Client-side unary interceptor
type UnaryClientInterceptor func(
    ctx context.Context,
    method string,
    req, reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption,
) error
```

### Creating gRPC Interceptors

```go
// middleware/grpc.go
package middleware

import (
    "context"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
)

// UnaryLogging creates a logging interceptor
func UnaryLogging(logger *slog.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        start := time.Now()
        
        // Get request ID from metadata
        requestID := ""
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if ids := md.Get("x-request-id"); len(ids) > 0 {
                requestID = ids[0]
            }
        }
        
        // Generate if not present
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // Add to context
        ctx = WithRequestID(ctx, requestID)
        
        // Log request
        logger.Info("grpc request",
            slog.String("method", info.FullMethod),
            slog.String("request_id", requestID),
        )
        
        // Call handler
        resp, err := handler(ctx, req)
        
        // Log response
        duration := time.Since(start)
        code := codes.OK
        if err != nil {
            code = status.Code(err)
        }
        
        logger.Info("grpc response",
            slog.String("method", info.FullMethod),
            slog.String("request_id", requestID),
            slog.String("code", code.String()),
            slog.Duration("duration", duration),
        )
        
        return resp, err
    }
}

// UnaryRecovery creates a panic recovery interceptor
func UnaryRecovery(logger *slog.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (resp interface{}, err error) {
        defer func() {
            if r := recover(); r != nil {
                logger.Error("grpc panic",
                    slog.String("method", info.FullMethod),
                    slog.Any("panic", r),
                    slog.String("stack", string(debug.Stack())),
                )
                
                err = status.Errorf(codes.Internal, "internal error")
            }
        }()
        
        return handler(ctx, req)
    }
}

// UnaryRequestID creates a request ID interceptor
func UnaryRequestID() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // Get or generate request ID
        requestID := ""
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if ids := md.Get("x-request-id"); len(ids) > 0 {
                requestID = ids[0]
            }
        }
        
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // Add to context
        ctx = WithRequestID(ctx, requestID)
        
        // Add to outgoing metadata (for response)
        grpc.SetHeader(ctx, metadata.Pairs("x-request-id", requestID))
        
        return handler(ctx, req)
    }
}

// UnaryMetrics creates a metrics interceptor
func UnaryMetrics() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        start := time.Now()
        
        resp, err := handler(ctx, req)
        
        duration := time.Since(start).Seconds()
        code := codes.OK
        if err != nil {
            code = status.Code(err)
        }
        
        // Record metrics (using Prometheus)
        grpcRequestsTotal.WithLabelValues(
            info.FullMethod,
            code.String(),
        ).Inc()
        
        grpcRequestDuration.WithLabelValues(
            info.FullMethod,
        ).Observe(duration)
        
        return resp, err
    }
}

// UnaryAuth creates an authentication interceptor
func UnaryAuth(validator func(ctx context.Context, token string) (context.Context, error)) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // Extract token from metadata
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }
        
        tokens := md.Get("authorization")
        if len(tokens) == 0 {
            return nil, status.Error(codes.Unauthenticated, "missing authorization")
        }
        
        // Validate token
        token := strings.TrimPrefix(tokens[0], "Bearer ")
        ctx, err := validator(ctx, token)
        if err != nil {
            return nil, status.Error(codes.Unauthenticated, "invalid token")
        }
        
        return handler(ctx, req)
    }
}
```

### Stream Interceptors

```go
// middleware/grpc_stream.go
package middleware

import (
    "context"
    
    "google.golang.org/grpc"
)

// StreamLogging creates a stream logging interceptor
func StreamLogging(logger *slog.Logger) grpc.StreamServerInterceptor {
    return func(
        srv interface{},
        ss grpc.ServerStream,
        info *grpc.StreamServerInfo,
        handler grpc.StreamHandler,
    ) error {
        start := time.Now()
        
        // Get request ID
        ctx := ss.Context()
        requestID := GetRequestID(ctx)
        if requestID == "" {
            if md, ok := metadata.FromIncomingContext(ctx); ok {
                if ids := md.Get("x-request-id"); len(ids) > 0 {
                    requestID = ids[0]
                }
            }
        }
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        logger.Info("grpc stream started",
            slog.String("method", info.FullMethod),
            slog.String("request_id", requestID),
            slog.Bool("client_stream", info.IsClientStream),
            slog.Bool("server_stream", info.IsServerStream),
        )
        
        // Wrap stream to add context
        wrapped := &wrappedServerStream{
            ServerStream: ss,
            ctx:          WithRequestID(ctx, requestID),
        }
        
        err := handler(srv, wrapped)
        
        logger.Info("grpc stream ended",
            slog.String("method", info.FullMethod),
            slog.String("request_id", requestID),
            slog.Duration("duration", time.Since(start)),
            slog.Any("error", err),
        )
        
        return err
    }
}

// StreamRecovery creates a stream recovery interceptor
func StreamRecovery(logger *slog.Logger) grpc.StreamServerInterceptor {
    return func(
        srv interface{},
        ss grpc.ServerStream,
        info *grpc.StreamServerInfo,
        handler grpc.StreamHandler,
    ) (err error) {
        defer func() {
            if r := recover(); r != nil {
                logger.Error("grpc stream panic",
                    slog.String("method", info.FullMethod),
                    slog.Any("panic", r),
                )
                err = status.Errorf(codes.Internal, "internal error")
            }
        }()
        
        return handler(srv, ss)
    }
}

// wrappedServerStream wraps grpc.ServerStream with custom context
type wrappedServerStream struct {
    grpc.ServerStream
    ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
    return w.ctx
}
```

### Client Interceptors

```go
// middleware/grpc_client.go
package middleware

import (
    "context"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

// UnaryClientRequestID propagates request ID to server
func UnaryClientRequestID() grpc.UnaryClientInterceptor {
    return func(
        ctx context.Context,
        method string,
        req, reply interface{},
        cc *grpc.ClientConn,
        invoker grpc.UnaryInvoker,
        opts ...grpc.CallOption,
    ) error {
        // Get request ID from context
        requestID := GetRequestID(ctx)
        if requestID != "" {
            // Add to outgoing metadata
            ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", requestID)
        }
        
        return invoker(ctx, method, req, reply, cc, opts...)
    }
}

// UnaryClientLogging logs client requests
func UnaryClientLogging(logger *slog.Logger) grpc.UnaryClientInterceptor {
    return func(
        ctx context.Context,
        method string,
        req, reply interface{},
        cc *grpc.ClientConn,
        invoker grpc.UnaryInvoker,
        opts ...grpc.CallOption,
    ) error {
        start := time.Now()
        
        logger.Debug("grpc client request",
            slog.String("method", method),
            slog.String("target", cc.Target()),
        )
        
        err := invoker(ctx, method, req, reply, cc, opts...)
        
        logger.Debug("grpc client response",
            slog.String("method", method),
            slog.Duration("duration", time.Since(start)),
            slog.Any("error", err),
        )
        
        return err
    }
}

// UnaryClientRetry creates a retry interceptor
func UnaryClientRetry(maxRetries int, backoff time.Duration) grpc.UnaryClientInterceptor {
    return func(
        ctx context.Context,
        method string,
        req, reply interface{},
        cc *grpc.ClientConn,
        invoker grpc.UnaryInvoker,
        opts ...grpc.CallOption,
    ) error {
        var lastErr error
        
        for attempt := 0; attempt <= maxRetries; attempt++ {
            err := invoker(ctx, method, req, reply, cc, opts...)
            if err == nil {
                return nil
            }
            
            lastErr = err
            
            // Check if error is retryable
            code := status.Code(err)
            if !isRetryable(code) {
                return err
            }
            
            // Check context
            if ctx.Err() != nil {
                return ctx.Err()
            }
            
            // Wait before retry
            if attempt < maxRetries {
                time.Sleep(backoff * time.Duration(attempt+1))
            }
        }
        
        return lastErr
    }
}

func isRetryable(code codes.Code) bool {
    switch code {
    case codes.Unavailable, codes.ResourceExhausted, codes.Aborted:
        return true
    default:
        return false
    }
}
```

### Applying gRPC Interceptors

```go
// cmd/server/main.go
package main

import (
    "google.golang.org/grpc"
    
    "myproject/middleware"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    // Create gRPC server with interceptors
    grpcServer := grpc.NewServer(
        // Chain multiple unary interceptors
        grpc.ChainUnaryInterceptor(
            middleware.UnaryRecovery(logger),
            middleware.UnaryRequestID(),
            middleware.UnaryLogging(logger),
            middleware.UnaryMetrics(),
            middleware.UnaryAuth(tokenValidator),
        ),
        // Chain multiple stream interceptors
        grpc.ChainStreamInterceptor(
            middleware.StreamRecovery(logger),
            middleware.StreamLogging(logger),
        ),
    )
    
    // Register services...
    
    // Start server...
}

// For client
func createClient() *grpc.ClientConn {
    conn, _ := grpc.Dial(
        "localhost:8081",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithChainUnaryInterceptor(
            middleware.UnaryClientRequestID(),
            middleware.UnaryClientLogging(logger),
            middleware.UnaryClientRetry(3, 100*time.Millisecond),
        ),
    )
    return conn
}
```

---

## 🎨 Middleware Patterns

### Chain Pattern

```go
// middleware/chain.go
package middleware

import (
    "net/http"
    
    goa "goa.design/goa/v3/pkg"
)

// HTTPChain chains multiple HTTP middleware
func HTTPChain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// Usage:
// chain := HTTPChain(Recovery, RequestID, Logging, CORS)
// handler := chain(mux)

// EndpointChain chains multiple endpoint middleware
func EndpointChain(middlewares ...func(goa.Endpoint) goa.Endpoint) func(goa.Endpoint) goa.Endpoint {
    return func(endpoint goa.Endpoint) goa.Endpoint {
        for i := len(middlewares) - 1; i >= 0; i-- {
            endpoint = middlewares[i](endpoint)
        }
        return endpoint
    }
}

// Usage:
// chain := EndpointChain(Logging, Timing, Auth)
// endpoint := chain(originalEndpoint)
```

### Conditional Middleware

```go
// middleware/conditional.go
package middleware

import (
    "net/http"
    
    goa "goa.design/goa/v3/pkg"
)

// HTTPConditional applies middleware based on condition
func HTTPConditional(
    condition func(*http.Request) bool,
    middleware func(http.Handler) http.Handler,
) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        withMiddleware := middleware(next)
        
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if condition(r) {
                withMiddleware.ServeHTTP(w, r)
            } else {
                next.ServeHTTP(w, r)
            }
        })
    }
}

// Usage examples:

// Only apply auth to /api paths
func isAPIPath(r *http.Request) bool {
    return strings.HasPrefix(r.URL.Path, "/api")
}
authMiddleware := HTTPConditional(isAPIPath, AuthMiddleware)

// Only log non-health check requests
func isNotHealthCheck(r *http.Request) bool {
    return r.URL.Path != "/health"
}
loggingMiddleware := HTTPConditional(isNotHealthCheck, LoggingMiddleware)
```

### Middleware Registry

```go
// middleware/registry.go
package middleware

import (
    "net/http"
    "sync"
)

// Registry manages middleware registration
type Registry struct {
    mu          sync.RWMutex
    middlewares map[string]func(http.Handler) http.Handler
    order       []string
}

// NewRegistry creates a new middleware registry
func NewRegistry() *Registry {
    return &Registry{
        middlewares: make(map[string]func(http.Handler) http.Handler),
    }
}

// Register adds middleware to the registry
func (r *Registry) Register(name string, mw func(http.Handler) http.Handler) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.middlewares[name] = mw
    r.order = append(r.order, name)
}

// Get retrieves middleware by name
func (r *Registry) Get(name string) (func(http.Handler) http.Handler, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    mw, ok := r.middlewares[name]
    return mw, ok
}

// Chain returns all middleware chained in registration order
func (r *Registry) Chain() func(http.Handler) http.Handler {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var mws []func(http.Handler) http.Handler
    for _, name := range r.order {
        mws = append(mws, r.middlewares[name])
    }
    
    return HTTPChain(mws...)
}

// Usage:
// registry := NewRegistry()
// registry.Register("recovery", Recovery)
// registry.Register("requestid", RequestID)
// registry.Register("logging", Logging)
// 
// handler := registry.Chain()(mux)
```

---

## 📦 Complete Examples

### Complete Server with All Middleware

```go
// cmd/server/main.go
package main

import (
    "context"
    "log/slog"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "google.golang.org/grpc"
    goahttp "goa.design/goa/v3/http"
    
    users "myproject/gen/users"
    userssvr "myproject/gen/http/users/server"
    usersgrpc "myproject/gen/grpc/users/server"
    userspb "myproject/gen/grpc/users/pb"
    
    "myproject/middleware"
)

func main() {
    // Create logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    
    // Create service
    svc := NewUsersService(logger)
    endpoints := users.NewEndpoints(svc)
    
    // Apply endpoint middleware
    endpoints.Use(middleware.EndpointRequestID())
    endpoints.Use(middleware.EndpointLogging(logger))
    endpoints.Use(middleware.TimingMiddleware())
    
    // Setup HTTP server
    httpServer := setupHTTPServer(endpoints, logger)
    
    // Setup gRPC server
    grpcServer := setupGRPCServer(endpoints, logger)
    
    // Start servers
    errChan := make(chan error, 2)
    
    go func() {
        logger.Info("starting HTTP server", slog.String("addr", ":8080"))
        errChan <- httpServer.ListenAndServe()
    }()
    
    go func() {
        lis, err := net.Listen("tcp", ":8081")
        if err != nil {
            errChan <- err
            return
        }
        logger.Info("starting gRPC server", slog.String("addr", ":8081"))
        errChan <- grpcServer.Serve(lis)
    }()
    
    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case err := <-errChan:
        logger.Error("server error", slog.Any("error", err))
    case sig := <-sigChan:
        logger.Info("received signal", slog.String("signal", sig.String()))
    }
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    logger.Info("shutting down servers")
    httpServer.Shutdown(ctx)
    grpcServer.GracefulStop()
    
    logger.Info("servers stopped")
}

func setupHTTPServer(endpoints *users.Endpoints, logger *slog.Logger) *http.Server {
    // Create Goa mux
    mux := goahttp.NewMuxer()
    
    // Create server
    server := userssvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        nil,
        nil,
    )
    userssvr.Mount(mux, server)
    
    // Apply HTTP middleware chain
    var handler http.Handler = mux
    
    // Apply in reverse order (first in chain = outermost)
    handler = middleware.HTTPMetrics(handler)
    handler = middleware.HTTPLogging(logger)(handler)
    handler = middleware.HTTPRequestID(handler)
    handler = middleware.CORS(middleware.DefaultCORSConfig())(handler)
    handler = middleware.Recovery(handler)
    
    return &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
}

func setupGRPCServer(endpoints *users.Endpoints, logger *slog.Logger) *grpc.Server {
    // Create gRPC server with interceptors
    grpcServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            middleware.UnaryRecovery(logger),
            middleware.UnaryRequestID(),
            middleware.UnaryLogging(logger),
            middleware.UnaryMetrics(),
        ),
        grpc.ChainStreamInterceptor(
            middleware.StreamRecovery(logger),
            middleware.StreamLogging(logger),
        ),
    )
    
    // Create and register service
    usersServer := usersgrpc.New(endpoints, nil)
    userspb.RegisterUsersServer(grpcServer, usersServer)
    
    return grpcServer
}
```

### Middleware Package Structure

```
middleware/
├── endpoint.go         # Endpoint middleware
├── http.go             # HTTP middleware base
├── http_logging.go     # HTTP logging
├── http_recovery.go    # Panic recovery
├── http_cors.go        # CORS handling
├── http_compress.go    # Response compression
├── http_ratelimit.go   # Rate limiting
├── http_cache.go       # Response caching
├── http_timeout.go     # Request timeout
├── requestid.go        # Request ID handling
├── grpc.go             # gRPC server interceptors
├── grpc_client.go      # gRPC client interceptors
├── grpc_stream.go      # Stream interceptors
├── metrics.go          # Metrics collection
├── auth.go             # Authentication
├── context.go          # Context utilities
└── chain.go            # Middleware chaining
```

---

## 📝 Summary

### Middleware Types
- **Endpoint Middleware**: Transport-agnostic, works with Goa endpoints
- **HTTP Middleware**: Standard `http.Handler` wrapper pattern
- **gRPC Interceptors**: Unary and Stream interceptors for gRPC

### Key Middleware Components
- **Request ID**: Unique identifier for distributed tracing
- **Logging**: Structured logging with context
- **Recovery**: Panic recovery and error handling
- **Metrics**: Prometheus/OpenTelemetry integration
- **Authentication**: Token validation and context enrichment
- **Rate Limiting**: Request throttling
- **CORS**: Cross-origin resource sharing

### Execution Order
- Middleware executes in order: first added = outermost
- "Before" logic runs top-to-bottom
- "After" logic runs bottom-to-top

### Best Practices
- Keep middleware focused on single responsibility
- Use context for passing data between middleware
- Chain middleware for clean composition
- Apply security middleware early in chain
- Apply logging/metrics middleware to capture all requests

---

## 📋 Knowledge Check

Before proceeding, ensure you can:

- [ ] Explain the difference between endpoint and transport middleware
- [ ] Create custom endpoint middleware
- [ ] Create custom HTTP middleware
- [ ] Implement request ID middleware and propagation
- [ ] Create structured logging middleware
- [ ] Implement panic recovery middleware
- [ ] Create rate limiting middleware
- [ ] Create gRPC unary interceptors
- [ ] Create gRPC stream interceptors
- [ ] Chain multiple middleware together
- [ ] Apply middleware conditionally
- [ ] Wire middleware to Goa servers

---

## 🔗 Quick Reference Links

- [Goa Middleware Package](https://pkg.go.dev/goa.design/goa/v3/middleware)
- [Go HTTP Middleware](https://pkg.go.dev/net/http#Handler)
- [gRPC Interceptors](https://grpc.io/docs/languages/go/interceptors/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [slog - Structured Logging](https://pkg.go.dev/log/slog)
- [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)

---

> **Next Up:** Part 7 - Testing & Deployment (Unit Testing, Integration Testing, Mocking, Docker, Kubernetes)
