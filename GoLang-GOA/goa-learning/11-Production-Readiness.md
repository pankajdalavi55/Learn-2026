# Part 11: Production Readiness

## Overview

This guide covers everything needed to prepare a Goa service for production deployment. We'll explore structured logging, graceful shutdown, configuration management, containerization, health checks, observability (metrics and tracing), and API versioning.

---

## Table of Contents

1. [Structured Logging](#1-structured-logging)
2. [Graceful Shutdown](#2-graceful-shutdown)
3. [Environment Configuration](#3-environment-configuration)
4. [Dockerizing Goa Services](#4-dockerizing-goa-services)
5. [Health Checks](#5-health-checks)
6. [Metrics with Prometheus](#6-metrics-with-prometheus)
7. [Distributed Tracing with OpenTelemetry](#7-distributed-tracing-with-opentelemetry)
8. [API Versioning](#8-api-versioning)
9. [Complete Production-Ready Example](#9-complete-production-ready-example)
10. [Best Practices & Checklist](#10-best-practices--checklist)

---

## 1. Structured Logging

### Why Structured Logging?

Traditional text logs are difficult to parse and search. Structured logging outputs logs in a consistent format (JSON) that can be easily consumed by log aggregation tools like ELK Stack, Splunk, or CloudWatch.

### Using Goa's Built-in Logger

```go
// main.go
package main

import (
    "context"
    "log"
    "os"

    goahttp "goa.design/goa/v3/http"
    "goa.design/goa/v3/middleware"
)

func main() {
    // Goa provides a simple logger interface
    logger := log.New(os.Stderr, "[api] ", log.Ltime|log.Lmicroseconds)
    
    // Create HTTP server with logging middleware
    handler := goahttp.NewMuxer()
    
    // Add request logging middleware
    handler.Use(middleware.Log(adapter(logger)))
}
```

### Integrating with Zerolog (Recommended for Production)

```go
// internal/logger/zerolog.go
package logger

import (
    "context"
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// Logger wraps zerolog for structured logging
type Logger struct {
    logger zerolog.Logger
}

// Config holds logger configuration
type Config struct {
    Level       string `json:"level" envconfig:"LOG_LEVEL" default:"info"`
    Pretty      bool   `json:"pretty" envconfig:"LOG_PRETTY" default:"false"`
    ServiceName string `json:"service_name" envconfig:"SERVICE_NAME"`
    Environment string `json:"environment" envconfig:"ENVIRONMENT"`
}

// NewLogger creates a new structured logger
func NewLogger(cfg Config) *Logger {
    // Set global log level
    level, err := zerolog.ParseLevel(cfg.Level)
    if err != nil {
        level = zerolog.InfoLevel
    }
    zerolog.SetGlobalLevel(level)

    // Configure output format
    var logger zerolog.Logger
    if cfg.Pretty {
        // Human-readable format for development
        logger = zerolog.New(zerolog.ConsoleWriter{
            Out:        os.Stdout,
            TimeFormat: time.RFC3339,
        })
    } else {
        // JSON format for production
        logger = zerolog.New(os.Stdout)
    }

    // Add common fields
    logger = logger.With().
        Timestamp().
        Str("service", cfg.ServiceName).
        Str("environment", cfg.Environment).
        Caller().
        Logger()

    return &Logger{logger: logger}
}

// Context keys for request-scoped logging
type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
)

// WithRequestID adds request ID to context
func (l *Logger) WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

// FromContext creates a logger with context values
func (l *Logger) FromContext(ctx context.Context) zerolog.Logger {
    logger := l.logger

    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        logger = logger.With().Str("request_id", requestID).Logger()
    }

    if userID, ok := ctx.Value(UserIDKey).(string); ok {
        logger = logger.With().Str("user_id", userID).Logger()
    }

    return logger
}

// Debug logs at debug level
func (l *Logger) Debug(ctx context.Context) *zerolog.Event {
    return l.FromContext(ctx).Debug()
}

// Info logs at info level
func (l *Logger) Info(ctx context.Context) *zerolog.Event {
    return l.FromContext(ctx).Info()
}

// Warn logs at warn level
func (l *Logger) Warn(ctx context.Context) *zerolog.Event {
    return l.FromContext(ctx).Warn()
}

// Error logs at error level
func (l *Logger) Error(ctx context.Context) *zerolog.Event {
    return l.FromContext(ctx).Error()
}

// Fatal logs at fatal level and exits
func (l *Logger) Fatal(ctx context.Context) *zerolog.Event {
    return l.FromContext(ctx).Fatal()
}
```

### Integrating with Zap (High Performance)

```go
// internal/logger/zap.go
package logger

import (
    "context"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// ZapLogger wraps uber's zap for high-performance logging
type ZapLogger struct {
    logger *zap.Logger
    sugar  *zap.SugaredLogger
}

// NewZapLogger creates a production-ready zap logger
func NewZapLogger(cfg Config) (*ZapLogger, error) {
    var zapConfig zap.Config

    if cfg.Environment == "production" {
        zapConfig = zap.NewProductionConfig()
        zapConfig.EncoderConfig.TimeKey = "timestamp"
        zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        zapConfig = zap.NewDevelopmentConfig()
        zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }

    // Parse log level
    level, err := zapcore.ParseLevel(cfg.Level)
    if err != nil {
        level = zapcore.InfoLevel
    }
    zapConfig.Level = zap.NewAtomicLevelAt(level)

    logger, err := zapConfig.Build(
        zap.AddCaller(),
        zap.AddStacktrace(zapcore.ErrorLevel),
        zap.Fields(
            zap.String("service", cfg.ServiceName),
            zap.String("environment", cfg.Environment),
        ),
    )
    if err != nil {
        return nil, err
    }

    return &ZapLogger{
        logger: logger,
        sugar:  logger.Sugar(),
    }, nil
}

// With creates a child logger with additional fields
func (l *ZapLogger) With(fields ...zap.Field) *ZapLogger {
    return &ZapLogger{
        logger: l.logger.With(fields...),
        sugar:  l.logger.With(fields...).Sugar(),
    }
}

// WithContext extracts context values and adds them as fields
func (l *ZapLogger) WithContext(ctx context.Context) *ZapLogger {
    fields := []zap.Field{}

    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        fields = append(fields, zap.String("request_id", requestID))
    }

    if userID, ok := ctx.Value(UserIDKey).(string); ok {
        fields = append(fields, zap.String("user_id", userID))
    }

    return l.With(fields...)
}

// Methods for different log levels
func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
    l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
    l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
    l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
    l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
    l.logger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *ZapLogger) Sync() error {
    return l.logger.Sync()
}
```

### Goa Logging Middleware

```go
// internal/middleware/logging.go
package middleware

import (
    "context"
    "net/http"
    "time"

    "github.com/google/uuid"
    "your-project/internal/logger"
)

// LoggingMiddleware creates a middleware that logs all requests
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Generate or extract request ID
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = uuid.New().String()
            }

            // Add request ID to context
            ctx := log.WithRequestID(r.Context(), requestID)
            r = r.WithContext(ctx)

            // Add request ID to response headers
            w.Header().Set("X-Request-ID", requestID)

            // Create a response writer wrapper to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

            // Log request start
            log.Info(ctx).
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Str("query", r.URL.RawQuery).
                Str("remote_addr", r.RemoteAddr).
                Str("user_agent", r.UserAgent()).
                Msg("request_started")

            // Process request
            next.ServeHTTP(wrapped, r)

            // Log request completion
            duration := time.Since(start)
            logEvent := log.Info(ctx)

            if wrapped.statusCode >= 400 {
                logEvent = log.Warn(ctx)
            }
            if wrapped.statusCode >= 500 {
                logEvent = log.Error(ctx)
            }

            logEvent.
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Int("status", wrapped.statusCode).
                Int64("duration_ms", duration.Milliseconds()).
                Int("response_size", wrapped.bytesWritten).
                Msg("request_completed")
        })
    }
}

// responseWriter wraps http.ResponseWriter to capture response details
type responseWriter struct {
    http.ResponseWriter
    statusCode   int
    bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    n, err := rw.ResponseWriter.Write(b)
    rw.bytesWritten += n
    return n, err
}
```

### Log Output Examples

**Development (Pretty Format):**
```
2024-01-15T10:30:45Z INF request_started method=GET path=/api/users remote_addr=192.168.1.1:54321 request_id=abc123
2024-01-15T10:30:45Z INF request_completed method=GET path=/api/users status=200 duration_ms=45 request_id=abc123
```

**Production (JSON Format):**
```json
{"level":"info","timestamp":"2024-01-15T10:30:45Z","service":"user-api","environment":"production","request_id":"abc123","method":"GET","path":"/api/users","message":"request_started"}
{"level":"info","timestamp":"2024-01-15T10:30:45Z","service":"user-api","environment":"production","request_id":"abc123","method":"GET","path":"/api/users","status":200,"duration_ms":45,"message":"request_completed"}
```

---

## 2. Graceful Shutdown

### Why Graceful Shutdown?

Graceful shutdown ensures:
- Active requests complete before server stops
- Database connections are properly closed
- Background jobs finish processing
- Resources are cleaned up properly

### Basic Graceful Shutdown

```go
// cmd/api/main.go
package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "your-project/internal/logger"
)

func main() {
    log := logger.NewLogger(logger.Config{
        Level:       "info",
        ServiceName: "api",
        Environment: os.Getenv("ENVIRONMENT"),
    })

    // Create HTTP server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      setupRoutes(),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Channel to listen for errors from server
    serverErrors := make(chan error, 1)

    // Start server in goroutine
    go func() {
        log.Info(context.Background()).
            Str("addr", server.Addr).
            Msg("starting server")
        serverErrors <- server.ListenAndServe()
    }()

    // Channel to listen for OS signals
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    // Block until signal or error
    select {
    case err := <-serverErrors:
        log.Fatal(context.Background()).
            Err(err).
            Msg("server error")

    case sig := <-shutdown:
        log.Info(context.Background()).
            Str("signal", sig.String()).
            Msg("shutdown signal received")

        // Create timeout context for shutdown
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        // Attempt graceful shutdown
        if err := server.Shutdown(ctx); err != nil {
            log.Error(context.Background()).
                Err(err).
                Msg("graceful shutdown failed, forcing exit")
            server.Close()
        }

        log.Info(context.Background()).Msg("server stopped")
    }
}
```

### Complete Graceful Shutdown with Resource Cleanup

```go
// internal/server/server.go
package server

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "your-project/internal/logger"
)

// Server encapsulates all server components
type Server struct {
    httpServer *http.Server
    db         *sql.DB
    logger     *logger.Logger
    wg         sync.WaitGroup
    shutdownCh chan struct{}
}

// Config holds server configuration
type Config struct {
    Host            string        `envconfig:"HOST" default:"0.0.0.0"`
    Port            int           `envconfig:"PORT" default:"8080"`
    ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"15s"`
    WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"15s"`
    IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
    ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

// New creates a new server instance
func New(cfg Config, handler http.Handler, db *sql.DB, log *logger.Logger) *Server {
    return &Server{
        httpServer: &http.Server{
            Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
            Handler:      handler,
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            IdleTimeout:  cfg.IdleTimeout,
        },
        db:         db,
        logger:     log,
        shutdownCh: make(chan struct{}),
    }
}

// Start begins serving requests and handles graceful shutdown
func (s *Server) Start(ctx context.Context) error {
    // Channel for server errors
    serverErrors := make(chan error, 1)

    // Start HTTP server
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.logger.Info(ctx).
            Str("addr", s.httpServer.Addr).
            Msg("HTTP server starting")
        
        if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
            serverErrors <- err
        }
    }()

    // Start background workers (if any)
    s.startBackgroundWorkers(ctx)

    // Wait for shutdown signal or error
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

    select {
    case err := <-serverErrors:
        return fmt.Errorf("server error: %w", err)

    case sig := <-shutdown:
        s.logger.Info(ctx).
            Str("signal", sig.String()).
            Msg("shutdown signal received, initiating graceful shutdown")

        return s.shutdown(ctx)
    }
}

// shutdown performs graceful shutdown of all components
func (s *Server) shutdown(ctx context.Context) error {
    // Signal background workers to stop
    close(s.shutdownCh)

    // Create shutdown context with timeout
    shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // Shutdown HTTP server
    s.logger.Info(ctx).Msg("shutting down HTTP server")
    if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
        s.logger.Error(ctx).Err(err).Msg("HTTP server shutdown error")
        s.httpServer.Close()
    }

    // Wait for background workers with timeout
    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        s.logger.Info(ctx).Msg("all background workers stopped")
    case <-shutdownCtx.Done():
        s.logger.Warn(ctx).Msg("timed out waiting for background workers")
    }

    // Close database connection
    if s.db != nil {
        s.logger.Info(ctx).Msg("closing database connection")
        if err := s.db.Close(); err != nil {
            s.logger.Error(ctx).Err(err).Msg("database close error")
        }
    }

    s.logger.Info(ctx).Msg("graceful shutdown complete")
    return nil
}

// startBackgroundWorkers starts any background processing
func (s *Server) startBackgroundWorkers(ctx context.Context) {
    // Example: Start a background job processor
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.runBackgroundJob(ctx)
    }()
}

// runBackgroundJob is an example background worker
func (s *Server) runBackgroundJob(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-s.shutdownCh:
            s.logger.Info(ctx).Msg("background job stopping")
            return
        case <-ticker.C:
            // Do periodic work
            s.logger.Debug(ctx).Msg("background job tick")
        }
    }
}
```

### Connection Draining Pattern

```go
// internal/middleware/draining.go
package middleware

import (
    "net/http"
    "sync"
    "sync/atomic"
)

// ConnectionDrainer tracks active connections for graceful shutdown
type ConnectionDrainer struct {
    activeRequests int64
    draining       int32
    wg             sync.WaitGroup
}

// NewConnectionDrainer creates a new drainer
func NewConnectionDrainer() *ConnectionDrainer {
    return &ConnectionDrainer{}
}

// Middleware returns HTTP middleware that tracks connections
func (d *ConnectionDrainer) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if we're draining
        if atomic.LoadInt32(&d.draining) == 1 {
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
            return
        }

        // Track this request
        d.wg.Add(1)
        atomic.AddInt64(&d.activeRequests, 1)
        defer func() {
            atomic.AddInt64(&d.activeRequests, -1)
            d.wg.Done()
        }()

        next.ServeHTTP(w, r)
    })
}

// StartDraining signals that we're shutting down
func (d *ConnectionDrainer) StartDraining() {
    atomic.StoreInt32(&d.draining, 1)
}

// Wait blocks until all active requests complete
func (d *ConnectionDrainer) Wait() {
    d.wg.Wait()
}

// ActiveRequests returns the current count
func (d *ConnectionDrainer) ActiveRequests() int64 {
    return atomic.LoadInt64(&d.activeRequests)
}
```

---

## 3. Environment Configuration

### Using envconfig

```go
// internal/config/config.go
package config

import (
    "fmt"
    "time"

    "github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration
type Config struct {
    // Server settings
    Server ServerConfig

    // Database settings
    Database DatabaseConfig

    // Redis settings
    Redis RedisConfig

    // Authentication settings
    Auth AuthConfig

    // Observability settings
    Observability ObservabilityConfig
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
    Host            string        `envconfig:"SERVER_HOST" default:"0.0.0.0"`
    Port            int           `envconfig:"SERVER_PORT" default:"8080"`
    ReadTimeout     time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"15s"`
    WriteTimeout    time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s"`
    IdleTimeout     time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
    ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"30s"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
    Host            string        `envconfig:"DB_HOST" required:"true"`
    Port            int           `envconfig:"DB_PORT" default:"5432"`
    User            string        `envconfig:"DB_USER" required:"true"`
    Password        string        `envconfig:"DB_PASSWORD" required:"true"`
    Name            string        `envconfig:"DB_NAME" required:"true"`
    SSLMode         string        `envconfig:"DB_SSL_MODE" default:"disable"`
    MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
    MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
    ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
    ConnMaxIdleTime time.Duration `envconfig:"DB_CONN_MAX_IDLE_TIME" default:"1m"`
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
    Host     string `envconfig:"REDIS_HOST" default:"localhost"`
    Port     int    `envconfig:"REDIS_PORT" default:"6379"`
    Password string `envconfig:"REDIS_PASSWORD"`
    DB       int    `envconfig:"REDIS_DB" default:"0"`
}

// Addr returns the Redis address
func (c *RedisConfig) Addr() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// AuthConfig holds authentication settings
type AuthConfig struct {
    JWTSecret            string        `envconfig:"JWT_SECRET" required:"true"`
    JWTExpiration        time.Duration `envconfig:"JWT_EXPIRATION" default:"24h"`
    RefreshExpiration    time.Duration `envconfig:"REFRESH_EXPIRATION" default:"168h"`
    BCryptCost           int           `envconfig:"BCRYPT_COST" default:"12"`
}

// ObservabilityConfig holds logging, metrics, and tracing settings
type ObservabilityConfig struct {
    // Logging
    LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
    LogFormat string `envconfig:"LOG_FORMAT" default:"json"`

    // Metrics
    MetricsEnabled bool   `envconfig:"METRICS_ENABLED" default:"true"`
    MetricsPort    int    `envconfig:"METRICS_PORT" default:"9090"`
    MetricsPath    string `envconfig:"METRICS_PATH" default:"/metrics"`

    // Tracing
    TracingEnabled     bool    `envconfig:"TRACING_ENABLED" default:"true"`
    TracingEndpoint    string  `envconfig:"TRACING_ENDPOINT" default:"localhost:4317"`
    TracingSampleRate  float64 `envconfig:"TRACING_SAMPLE_RATE" default:"0.1"`
    TracingServiceName string  `envconfig:"TRACING_SERVICE_NAME" default:"api"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
    var cfg Config

    if err := envconfig.Process("", &cfg); err != nil {
        return nil, fmt.Errorf("loading config: %w", err)
    }

    return &cfg, nil
}

// MustLoad loads config or panics
func MustLoad() *Config {
    cfg, err := Load()
    if err != nil {
        panic(fmt.Sprintf("failed to load config: %v", err))
    }
    return cfg
}

// Validate performs validation on configuration
func (c *Config) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }

    if c.Database.MaxOpenConns < c.Database.MaxIdleConns {
        return fmt.Errorf("max open conns must be >= max idle conns")
    }

    return nil
}
```

### Using Viper for Advanced Configuration

```go
// internal/config/viper.go
package config

import (
    "fmt"
    "strings"

    "github.com/spf13/viper"
)

// LoadWithViper loads configuration from multiple sources
func LoadWithViper(configPath string) (*Config, error) {
    v := viper.New()

    // Set defaults
    setDefaults(v)

    // Read from config file
    if configPath != "" {
        v.SetConfigFile(configPath)
        if err := v.ReadInConfig(); err != nil {
            return nil, fmt.Errorf("reading config file: %w", err)
        }
    }

    // Read from environment variables
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    v.AutomaticEnv()

    // Unmarshal to struct
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshaling config: %w", err)
    }

    return &cfg, nil
}

func setDefaults(v *viper.Viper) {
    // Server defaults
    v.SetDefault("server.host", "0.0.0.0")
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.readtimeout", "15s")
    v.SetDefault("server.writetimeout", "30s")

    // Database defaults
    v.SetDefault("database.port", 5432)
    v.SetDefault("database.sslmode", "disable")
    v.SetDefault("database.maxopenconns", 25)
    v.SetDefault("database.maxidleconns", 5)

    // Observability defaults
    v.SetDefault("observability.loglevel", "info")
    v.SetDefault("observability.logformat", "json")
    v.SetDefault("observability.metricsenabled", true)
}
```

### Config File Example (YAML)

```yaml
# config/config.yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "15s"
  write_timeout: "30s"
  idle_timeout: "60s"
  shutdown_timeout: "30s"

database:
  host: "localhost"
  port: 5432
  user: "apiuser"
  password: "${DB_PASSWORD}"  # From environment
  name: "apidb"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5

redis:
  host: "localhost"
  port: 6379
  db: 0

auth:
  jwt_expiration: "24h"
  refresh_expiration: "168h"
  bcrypt_cost: 12

observability:
  log_level: "info"
  log_format: "json"
  metrics_enabled: true
  metrics_port: 9090
  tracing_enabled: true
  tracing_endpoint: "localhost:4317"
  tracing_sample_rate: 0.1
```

### Environment-Specific Configuration

```go
// internal/config/environments.go
package config

import (
    "fmt"
    "os"
    "path/filepath"
)

// Environment represents the runtime environment
type Environment string

const (
    Development Environment = "development"
    Staging     Environment = "staging"
    Production  Environment = "production"
)

// GetEnvironment returns the current environment
func GetEnvironment() Environment {
    env := os.Getenv("ENVIRONMENT")
    switch env {
    case "production", "prod":
        return Production
    case "staging", "stage":
        return Staging
    default:
        return Development
    }
}

// LoadForEnvironment loads environment-specific configuration
func LoadForEnvironment() (*Config, error) {
    env := GetEnvironment()

    // Base config path
    basePath := "config"

    // Load base config
    v := viper.New()
    v.SetConfigFile(filepath.Join(basePath, "config.yaml"))
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("reading base config: %w", err)
    }

    // Load environment-specific overrides
    envConfigPath := filepath.Join(basePath, fmt.Sprintf("config.%s.yaml", env))
    if _, err := os.Stat(envConfigPath); err == nil {
        v.SetConfigFile(envConfigPath)
        if err := v.MergeInConfig(); err != nil {
            return nil, fmt.Errorf("reading env config: %w", err)
        }
    }

    // Environment variables override everything
    v.AutomaticEnv()

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshaling config: %w", err)
    }

    return &cfg, nil
}
```

---

## 4. Dockerizing Goa Services

### Basic Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /app/server ./cmd/api

# Final stage - minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/server /server

# Copy config files if needed
COPY --from=builder /app/config /config

# Expose port
EXPOSE 8080

# Set health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "health"]

# Run as non-root user
USER 1000:1000

# Set entrypoint
ENTRYPOINT ["/server"]
```

### Multi-Stage Dockerfile with Testing

```dockerfile
# Dockerfile.full
# ==========================================
# Stage 1: Dependencies
# ==========================================
FROM golang:1.21-alpine AS deps

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ==========================================
# Stage 2: Testing
# ==========================================
FROM deps AS test

COPY . .

# Run tests with coverage
RUN go test -race -coverprofile=coverage.out -covermode=atomic ./...

# ==========================================
# Stage 3: Build
# ==========================================
FROM deps AS builder

# Build arguments for versioning
ARG VERSION=dev
ARG BUILD_TIME
ARG COMMIT_SHA

WORKDIR /app

COPY . .

# Build with optimizations and version info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s \
        -X main.version=${VERSION} \
        -X main.buildTime=${BUILD_TIME} \
        -X main.commitSHA=${COMMIT_SHA}" \
    -o /server ./cmd/api

# ==========================================  
# Stage 4: Production Runtime
# ==========================================
FROM gcr.io/distroless/static:nonroot AS production

WORKDIR /

# Copy artifacts
COPY --from=builder /server /server
COPY --from=builder /app/config /config

# Use nonroot user
USER nonroot:nonroot

EXPOSE 8080 9090

ENTRYPOINT ["/server"]

# ==========================================
# Stage 5: Development Runtime
# ==========================================
FROM golang:1.21-alpine AS development

WORKDIR /app

# Install development tools
RUN apk add --no-cache git make curl

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

COPY . .

EXPOSE 8080 9090

CMD ["air", "-c", ".air.toml"]
```

### Docker Compose for Local Development

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
      target: development
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      - ENVIRONMENT=development
      - LOG_LEVEL=debug
      - LOG_FORMAT=pretty
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=apiuser
      - DB_PASSWORD=secret
      - DB_NAME=apidb
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=dev-secret-key
      - TRACING_ENDPOINT=jaeger:4317
    volumes:
      - .:/app
      - go-cache:/go/pkg/mod
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: apiuser
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: apidb
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U apiuser -d apidb"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # UI
      - "4317:4317"    # gRPC OTLP
      - "4318:4318"    # HTTP OTLP
    environment:
      - COLLECTOR_OTLP_ENABLED=true

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9091:9090"
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
      - ./config/grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./config/grafana/datasources:/etc/grafana/provisioning/datasources

volumes:
  postgres-data:
  redis-data:
  prometheus-data:
  grafana-data:
  go-cache:
```

### Production Docker Compose

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  api:
    image: ${REGISTRY}/api:${TAG:-latest}
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 128M
      update_config:
        parallelism: 1
        delay: 10s
        failure_action: rollback
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
      - LOG_FORMAT=json
    secrets:
      - db_password
      - jwt_secret
    healthcheck:
      test: ["CMD", "/server", "health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s
    networks:
      - api-network

secrets:
  db_password:
    external: true
  jwt_secret:
    external: true

networks:
  api-network:
    driver: overlay
```

---

## 5. Health Checks

### Design DSL for Health Endpoints

```go
// design/health.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("health", func() {
    Description("Health check service")

    // Liveness probe - is the service running?
    Method("live", func() {
        Description("Liveness check - returns ok if service is running")
        HTTP(func() {
            GET("/health/live")
            Response(StatusOK)
        })
        Result(HealthStatus)
    })

    // Readiness probe - is the service ready to accept traffic?
    Method("ready", func() {
        Description("Readiness check - returns ok if service can handle requests")
        HTTP(func() {
            GET("/health/ready")
            Response(StatusOK)
            Response("not_ready", StatusServiceUnavailable)
        })
        Result(HealthStatus)
        Error("not_ready", HealthStatus, "Service is not ready")
    })

    // Detailed health check for debugging
    Method("status", func() {
        Description("Detailed health status with component checks")
        HTTP(func() {
            GET("/health/status")
            Response(StatusOK)
            Response("degraded", StatusOK)
            Response("unhealthy", StatusServiceUnavailable)
        })
        Result(DetailedHealthStatus)
        Error("degraded", DetailedHealthStatus, "Some components degraded")
        Error("unhealthy", DetailedHealthStatus, "Service unhealthy")
    })
})

// HealthStatus is a simple health response
var HealthStatus = Type("HealthStatus", func() {
    Attribute("status", String, "Health status", func() {
        Enum("ok", "degraded", "unhealthy")
        Example("ok")
    })
    Attribute("timestamp", String, "Check timestamp", func() {
        Format(FormatDateTime)
    })
    Required("status", "timestamp")
})

// DetailedHealthStatus includes component-level health
var DetailedHealthStatus = Type("DetailedHealthStatus", func() {
    Attribute("status", String, "Overall status", func() {
        Enum("ok", "degraded", "unhealthy")
    })
    Attribute("timestamp", String, func() {
        Format(FormatDateTime)
    })
    Attribute("version", String, "Service version")
    Attribute("uptime", String, "Service uptime")
    Attribute("components", MapOf(String, ComponentHealth), "Component health statuses")
    Required("status", "timestamp", "components")
})

// ComponentHealth represents health of a single component
var ComponentHealth = Type("ComponentHealth", func() {
    Attribute("status", String, func() {
        Enum("ok", "degraded", "unhealthy")
    })
    Attribute("message", String, "Status message")
    Attribute("latency_ms", Int64, "Check latency in milliseconds")
    Attribute("last_check", String, func() {
        Format(FormatDateTime)
    })
    Required("status")
})
```

### Health Service Implementation

```go
// internal/service/health.go
package service

import (
    "context"
    "sync"
    "time"

    "your-project/gen/health"
)

// HealthService implements health checks
type HealthService struct {
    startTime time.Time
    version   string
    checkers  map[string]HealthChecker
    mu        sync.RWMutex
}

// HealthChecker interface for component health checks
type HealthChecker interface {
    Name() string
    Check(ctx context.Context) ComponentStatus
}

// ComponentStatus represents the health of a single component
type ComponentStatus struct {
    Status    string
    Message   string
    LatencyMs int64
}

// NewHealthService creates a new health service
func NewHealthService(version string) *HealthService {
    return &HealthService{
        startTime: time.Now(),
        version:   version,
        checkers:  make(map[string]HealthChecker),
    }
}

// RegisterChecker adds a health checker
func (s *HealthService) RegisterChecker(checker HealthChecker) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.checkers[checker.Name()] = checker
}

// Live implements liveness check
func (s *HealthService) Live(ctx context.Context) (*health.HealthStatus, error) {
    return &health.HealthStatus{
        Status:    "ok",
        Timestamp: time.Now().Format(time.RFC3339),
    }, nil
}

// Ready implements readiness check
func (s *HealthService) Ready(ctx context.Context) (*health.HealthStatus, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    // Check all critical components
    for _, checker := range s.checkers {
        status := checker.Check(ctx)
        if status.Status == "unhealthy" {
            return nil, health.MakeNotReady(&health.HealthStatus{
                Status:    "unhealthy",
                Timestamp: time.Now().Format(time.RFC3339),
            })
        }
    }

    return &health.HealthStatus{
        Status:    "ok",
        Timestamp: time.Now().Format(time.RFC3339),
    }, nil
}

// Status implements detailed health check
func (s *HealthService) Status(ctx context.Context) (*health.DetailedHealthStatus, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    components := make(map[string]*health.ComponentHealth)
    overallStatus := "ok"

    // Check all components in parallel
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for name, checker := range s.checkers {
        wg.Add(1)
        go func(n string, c HealthChecker) {
            defer wg.Done()
            
            start := time.Now()
            status := c.Check(ctx)
            latency := time.Since(start).Milliseconds()
            
            mu.Lock()
            defer mu.Unlock()
            
            components[n] = &health.ComponentHealth{
                Status:    status.Status,
                Message:   &status.Message,
                LatencyMs: &latency,
            }
            
            // Determine overall status
            if status.Status == "unhealthy" && overallStatus != "unhealthy" {
                overallStatus = "unhealthy"
            } else if status.Status == "degraded" && overallStatus == "ok" {
                overallStatus = "degraded"
            }
        }(name, checker)
    }
    
    wg.Wait()

    uptime := time.Since(s.startTime).String()
    result := &health.DetailedHealthStatus{
        Status:     overallStatus,
        Timestamp:  time.Now().Format(time.RFC3339),
        Version:    &s.version,
        Uptime:     &uptime,
        Components: components,
    }

    // Return appropriate error for degraded/unhealthy status
    if overallStatus == "unhealthy" {
        return nil, health.MakeUnhealthy(result)
    }
    if overallStatus == "degraded" {
        return nil, health.MakeDegraded(result)
    }

    return result, nil
}
```

### Health Checkers for Common Components

```go
// internal/health/checkers.go
package health

import (
    "context"
    "database/sql"
    "time"

    "github.com/redis/go-redis/v9"
)

// PostgresChecker checks PostgreSQL health
type PostgresChecker struct {
    db      *sql.DB
    timeout time.Duration
}

func NewPostgresChecker(db *sql.DB) *PostgresChecker {
    return &PostgresChecker{
        db:      db,
        timeout: 5 * time.Second,
    }
}

func (c *PostgresChecker) Name() string {
    return "postgres"
}

func (c *PostgresChecker) Check(ctx context.Context) ComponentStatus {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    if err := c.db.PingContext(ctx); err != nil {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    }

    // Check connection pool stats
    stats := c.db.Stats()
    if stats.OpenConnections >= stats.MaxOpenConnections {
        return ComponentStatus{
            Status:  "degraded",
            Message: "connection pool exhausted",
        }
    }

    return ComponentStatus{
        Status:  "ok",
        Message: "connection successful",
    }
}

// RedisChecker checks Redis health
type RedisChecker struct {
    client  *redis.Client
    timeout time.Duration
}

func NewRedisChecker(client *redis.Client) *RedisChecker {
    return &RedisChecker{
        client:  client,
        timeout: 3 * time.Second,
    }
}

func (c *RedisChecker) Name() string {
    return "redis"
}

func (c *RedisChecker) Check(ctx context.Context) ComponentStatus {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    result, err := c.client.Ping(ctx).Result()
    if err != nil {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    }

    if result != "PONG" {
        return ComponentStatus{
            Status:  "degraded",
            Message: "unexpected ping response: " + result,
        }
    }

    return ComponentStatus{
        Status:  "ok",
        Message: "connection successful",
    }
}

// ExternalAPIChecker checks external API health
type ExternalAPIChecker struct {
    name    string
    url     string
    timeout time.Duration
}

func NewExternalAPIChecker(name, url string) *ExternalAPIChecker {
    return &ExternalAPIChecker{
        name:    name,
        url:     url,
        timeout: 10 * time.Second,
    }
}

func (c *ExternalAPIChecker) Name() string {
    return c.name
}

func (c *ExternalAPIChecker) Check(ctx context.Context) ComponentStatus {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
    if err != nil {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 500 {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: fmt.Sprintf("status code: %d", resp.StatusCode),
        }
    }

    if resp.StatusCode >= 400 {
        return ComponentStatus{
            Status:  "degraded",
            Message: fmt.Sprintf("status code: %d", resp.StatusCode),
        }
    }

    return ComponentStatus{
        Status:  "ok",
        Message: "reachable",
    }
}

// DiskSpaceChecker checks available disk space
type DiskSpaceChecker struct {
    path              string
    warningThreshold  uint64 // bytes
    criticalThreshold uint64 // bytes
}

func NewDiskSpaceChecker(path string, warningGB, criticalGB uint64) *DiskSpaceChecker {
    return &DiskSpaceChecker{
        path:              path,
        warningThreshold:  warningGB * 1024 * 1024 * 1024,
        criticalThreshold: criticalGB * 1024 * 1024 * 1024,
    }
}

func (c *DiskSpaceChecker) Name() string {
    return "disk_space"
}

func (c *DiskSpaceChecker) Check(ctx context.Context) ComponentStatus {
    var stat syscall.Statfs_t
    if err := syscall.Statfs(c.path, &stat); err != nil {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    }

    available := stat.Bavail * uint64(stat.Bsize)
    
    if available < c.criticalThreshold {
        return ComponentStatus{
            Status:  "unhealthy",
            Message: fmt.Sprintf("only %d GB available", available/1024/1024/1024),
        }
    }

    if available < c.warningThreshold {
        return ComponentStatus{
            Status:  "degraded",
            Message: fmt.Sprintf("only %d GB available", available/1024/1024/1024),
        }
    }

    return ComponentStatus{
        Status:  "ok",
        Message: fmt.Sprintf("%d GB available", available/1024/1024/1024),
    }
}
```

---

## 6. Metrics with Prometheus

### Prometheus Metrics Setup

```go
// internal/metrics/prometheus.go
package metrics

import (
    "net/http"
    "strconv"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all application metrics
type Metrics struct {
    // HTTP metrics
    HTTPRequestsTotal    *prometheus.CounterVec
    HTTPRequestDuration  *prometheus.HistogramVec
    HTTPRequestsInFlight prometheus.Gauge
    HTTPResponseSize     *prometheus.HistogramVec

    // Business metrics
    UsersRegistered prometheus.Counter
    OrdersCreated   prometheus.Counter
    OrdersTotal     *prometheus.CounterVec
    
    // Database metrics
    DBQueryDuration  *prometheus.HistogramVec
    DBConnectionsOpen prometheus.Gauge
    DBConnectionsIdle prometheus.Gauge

    // Cache metrics
    CacheHits   prometheus.Counter
    CacheMisses prometheus.Counter

    // Custom business metrics
    RevenueTotal    *prometheus.CounterVec
    ActiveSessions  prometheus.Gauge
}

// NewMetrics creates and registers all metrics
func NewMetrics(namespace, subsystem string) *Metrics {
    m := &Metrics{
        // HTTP Request Counter
        HTTPRequestsTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "http_requests_total",
                Help:      "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),

        // HTTP Request Duration Histogram
        HTTPRequestDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "http_request_duration_seconds",
                Help:      "HTTP request duration in seconds",
                Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
            },
            []string{"method", "path", "status"},
        ),

        // In-flight requests
        HTTPRequestsInFlight: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "http_requests_in_flight",
                Help:      "Number of HTTP requests currently being processed",
            },
        ),

        // Response size
        HTTPResponseSize: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "http_response_size_bytes",
                Help:      "HTTP response size in bytes",
                Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
            },
            []string{"method", "path"},
        ),

        // Database metrics
        DBQueryDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Subsystem: "database",
                Name:      "query_duration_seconds",
                Help:      "Database query duration in seconds",
                Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
            },
            []string{"query_type", "table"},
        ),

        DBConnectionsOpen: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: "database",
                Name:      "connections_open",
                Help:      "Number of open database connections",
            },
        ),

        DBConnectionsIdle: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: "database",
                Name:      "connections_idle",
                Help:      "Number of idle database connections",
            },
        ),

        // Cache metrics
        CacheHits: promauto.NewCounter(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "cache",
                Name:      "hits_total",
                Help:      "Total number of cache hits",
            },
        ),

        CacheMisses: promauto.NewCounter(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "cache",
                Name:      "misses_total",
                Help:      "Total number of cache misses",
            },
        ),

        // Business metrics
        UsersRegistered: promauto.NewCounter(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "business",
                Name:      "users_registered_total",
                Help:      "Total number of registered users",
            },
        ),

        OrdersCreated: promauto.NewCounter(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "business",
                Name:      "orders_created_total",
                Help:      "Total number of orders created",
            },
        ),

        OrdersTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "business",
                Name:      "orders_by_status_total",
                Help:      "Total orders by status",
            },
            []string{"status"},
        ),

        RevenueTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "business",
                Name:      "revenue_total",
                Help:      "Total revenue in cents",
            },
            []string{"currency", "payment_method"},
        ),

        ActiveSessions: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: "business",
                Name:      "active_sessions",
                Help:      "Number of active user sessions",
            },
        ),
    }

    return m
}

// Handler returns the Prometheus HTTP handler
func Handler() http.Handler {
    return promhttp.Handler()
}
```

### Prometheus Metrics Middleware

```go
// internal/middleware/metrics.go
package middleware

import (
    "net/http"
    "strconv"
    "time"

    "your-project/internal/metrics"
)

// MetricsMiddleware creates HTTP metrics middleware
func MetricsMiddleware(m *metrics.Metrics) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Track in-flight requests
            m.HTTPRequestsInFlight.Inc()
            defer m.HTTPRequestsInFlight.Dec()

            // Create response wrapper
            wrapped := &metricsResponseWriter{
                ResponseWriter: w,
                statusCode:     http.StatusOK,
            }

            // Record start time
            start := time.Now()

            // Process request
            next.ServeHTTP(wrapped, r)

            // Record metrics
            duration := time.Since(start).Seconds()
            status := strconv.Itoa(wrapped.statusCode)
            path := normalizePath(r.URL.Path) // Prevent high cardinality

            m.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
            m.HTTPRequestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
            m.HTTPResponseSize.WithLabelValues(r.Method, path).Observe(float64(wrapped.bytesWritten))
        })
    }
}

type metricsResponseWriter struct {
    http.ResponseWriter
    statusCode   int
    bytesWritten int
}

func (w *metricsResponseWriter) WriteHeader(code int) {
    w.statusCode = code
    w.ResponseWriter.WriteHeader(code)
}

func (w *metricsResponseWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.bytesWritten += n
    return n, err
}

// normalizePath prevents high cardinality by replacing dynamic segments
func normalizePath(path string) string {
    // Replace UUIDs with placeholder
    uuidPattern := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
    path = uuidPattern.ReplaceAllString(path, ":id")

    // Replace numeric IDs with placeholder
    numericPattern := regexp.MustCompile(`/\d+`)
    path = numericPattern.ReplaceAllString(path, "/:id")

    return path
}
```

### Database Metrics Collector

```go
// internal/metrics/db_collector.go
package metrics

import (
    "context"
    "database/sql"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

// DBCollector collects database connection pool metrics
type DBCollector struct {
    db *sql.DB

    maxOpenDesc     *prometheus.Desc
    openDesc        *prometheus.Desc
    inUseDesc       *prometheus.Desc
    idleDesc        *prometheus.Desc
    waitCountDesc   *prometheus.Desc
    waitDurationDesc *prometheus.Desc
}

// NewDBCollector creates a new database metrics collector
func NewDBCollector(db *sql.DB, namespace string) *DBCollector {
    return &DBCollector{
        db: db,
        maxOpenDesc: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "db", "max_open_connections"),
            "Maximum number of open connections to the database",
            nil, nil,
        ),
        openDesc: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "db", "open_connections"),
            "Number of established connections both in use and idle",
            nil, nil,
        ),
        inUseDesc: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "db", "in_use_connections"),
            "Number of connections currently in use",
            nil, nil,
        ),
        idleDesc: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "db", "idle_connections"),
            "Number of idle connections",
            nil, nil,
        ),
        waitCountDesc: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "db", "wait_count_total"),
            "Total number of connections waited for",
            nil, nil,
        ),
        waitDurationDesc: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "db", "wait_duration_seconds_total"),
            "Total time blocked waiting for a new connection",
            nil, nil,
        ),
    }
}

// Describe implements prometheus.Collector
func (c *DBCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.maxOpenDesc
    ch <- c.openDesc
    ch <- c.inUseDesc
    ch <- c.idleDesc
    ch <- c.waitCountDesc
    ch <- c.waitDurationDesc
}

// Collect implements prometheus.Collector
func (c *DBCollector) Collect(ch chan<- prometheus.Metric) {
    stats := c.db.Stats()

    ch <- prometheus.MustNewConstMetric(c.maxOpenDesc, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
    ch <- prometheus.MustNewConstMetric(c.openDesc, prometheus.GaugeValue, float64(stats.OpenConnections))
    ch <- prometheus.MustNewConstMetric(c.inUseDesc, prometheus.GaugeValue, float64(stats.InUse))
    ch <- prometheus.MustNewConstMetric(c.idleDesc, prometheus.GaugeValue, float64(stats.Idle))
    ch <- prometheus.MustNewConstMetric(c.waitCountDesc, prometheus.CounterValue, float64(stats.WaitCount))
    ch <- prometheus.MustNewConstMetric(c.waitDurationDesc, prometheus.CounterValue, stats.WaitDuration.Seconds())
}

// Register registers the collector with Prometheus
func (c *DBCollector) Register() error {
    return prometheus.Register(c)
}
```

### Prometheus Configuration

```yaml
# config/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

rule_files:
  # - "alerts.yml"

scrape_configs:
  - job_name: 'api'
    static_configs:
      - targets: ['api:9090']
    metrics_path: /metrics
    scrape_interval: 10s

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### Sample Grafana Dashboard JSON

```json
{
  "dashboard": {
    "title": "API Service Dashboard",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(api_http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}} {{status}}"
          }
        ]
      },
      {
        "title": "Request Latency (p95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(api_http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p95"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "sum(rate(api_http_requests_total{status=~\"5..\"}[5m])) / sum(rate(api_http_requests_total[5m])) * 100",
            "legendFormat": "Error %"
          }
        ]
      },
      {
        "title": "Database Connections",
        "type": "graph",
        "targets": [
          {
            "expr": "api_db_open_connections",
            "legendFormat": "Open"
          },
          {
            "expr": "api_db_in_use_connections",
            "legendFormat": "In Use"
          },
          {
            "expr": "api_db_idle_connections",
            "legendFormat": "Idle"
          }
        ]
      }
    ]
  }
}
```

---

## 7. Distributed Tracing with OpenTelemetry

### OpenTelemetry Setup

```go
// internal/telemetry/tracing.go
package telemetry

import (
    "context"
    "fmt"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
    "go.opentelemetry.io/otel/trace"
)

// TracerConfig holds tracing configuration
type TracerConfig struct {
    ServiceName    string
    ServiceVersion string
    Environment    string
    Endpoint       string
    SampleRate     float64
    Enabled        bool
}

// TracerProvider wraps the OpenTelemetry tracer provider
type TracerProvider struct {
    provider *sdktrace.TracerProvider
    tracer   trace.Tracer
}

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(ctx context.Context, cfg TracerConfig) (*TracerProvider, error) {
    if !cfg.Enabled {
        // Return a no-op tracer
        return &TracerProvider{
            tracer: otel.Tracer(cfg.ServiceName),
        }, nil
    }

    // Create OTLP exporter
    client := otlptracegrpc.NewClient(
        otlptracegrpc.WithEndpoint(cfg.Endpoint),
        otlptracegrpc.WithInsecure(), // Use TLS in production
    )

    exporter, err := otlptrace.New(ctx, client)
    if err != nil {
        return nil, fmt.Errorf("creating OTLP exporter: %w", err)
    }

    // Create resource with service information
    res, err := resource.Merge(
        resource.Default(),
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceVersion(cfg.ServiceVersion),
            semconv.DeploymentEnvironment(cfg.Environment),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("creating resource: %w", err)
    }

    // Create sampler based on configuration
    var sampler sdktrace.Sampler
    if cfg.SampleRate >= 1.0 {
        sampler = sdktrace.AlwaysSample()
    } else if cfg.SampleRate <= 0 {
        sampler = sdktrace.NeverSample()
    } else {
        sampler = sdktrace.TraceIDRatioBased(cfg.SampleRate)
    }

    // Create tracer provider
    provider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.ParentBased(sampler)),
    )

    // Set global provider and propagator
    otel.SetTracerProvider(provider)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return &TracerProvider{
        provider: provider,
        tracer:   provider.Tracer(cfg.ServiceName),
    }, nil
}

// Shutdown gracefully shuts down the tracer
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
    if tp.provider != nil {
        return tp.provider.Shutdown(ctx)
    }
    return nil
}

// Tracer returns the tracer instance
func (tp *TracerProvider) Tracer() trace.Tracer {
    return tp.tracer
}

// StartSpan starts a new span
func (tp *TracerProvider) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
    return tp.tracer.Start(ctx, name, opts...)
}
```

### Tracing Middleware for HTTP

```go
// internal/middleware/tracing.go
package middleware

import (
    "net/http"

    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates tracing middleware using otelhttp
func TracingMiddleware(serviceName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return otelhttp.NewHandler(next, serviceName,
            otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
            otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
                return r.Method + " " + r.URL.Path
            }),
        )
    }
}

// CustomTracingMiddleware provides more control over tracing
func CustomTracingMiddleware(tracer trace.Tracer) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract context from incoming request
            ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

            // Start span
            ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
                trace.WithSpanKind(trace.SpanKindServer),
                trace.WithAttributes(
                    attribute.String("http.method", r.Method),
                    attribute.String("http.url", r.URL.String()),
                    attribute.String("http.host", r.Host),
                    attribute.String("http.user_agent", r.UserAgent()),
                    attribute.String("http.remote_addr", r.RemoteAddr),
                ),
            )
            defer span.End()

            // Create wrapped response writer to capture status code
            wrapped := &tracingResponseWriter{ResponseWriter: w, statusCode: 200}

            // Update request with new context
            r = r.WithContext(ctx)

            // Call next handler
            next.ServeHTTP(wrapped, r)

            // Add response attributes
            span.SetAttributes(
                attribute.Int("http.status_code", wrapped.statusCode),
                attribute.Int("http.response_size", wrapped.bytesWritten),
            )

            // Mark span as error if status >= 400
            if wrapped.statusCode >= 400 {
                span.SetAttributes(attribute.Bool("error", true))
            }
        })
    }
}

type tracingResponseWriter struct {
    http.ResponseWriter
    statusCode   int
    bytesWritten int
}

func (w *tracingResponseWriter) WriteHeader(code int) {
    w.statusCode = code
    w.ResponseWriter.WriteHeader(code)
}

func (w *tracingResponseWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.bytesWritten += n
    return n, err
}
```

### Tracing in Service Layer

```go
// internal/service/user.go
package service

import (
    "context"
    "fmt"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

type UserService struct {
    repo   UserRepository
    tracer trace.Tracer
}

func NewUserService(repo UserRepository, tracer trace.Tracer) *UserService {
    return &UserService{
        repo:   repo,
        tracer: tracer,
    }
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    ctx, span := s.tracer.Start(ctx, "UserService.GetUser",
        trace.WithAttributes(attribute.String("user.id", id)),
    )
    defer span.End()

    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, fmt.Errorf("finding user: %w", err)
    }

    if user == nil {
        span.SetAttributes(attribute.Bool("user.found", false))
        return nil, ErrUserNotFound
    }

    span.SetAttributes(
        attribute.Bool("user.found", true),
        attribute.String("user.email", user.Email),
    )

    return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    ctx, span := s.tracer.Start(ctx, "UserService.CreateUser",
        trace.WithAttributes(
            attribute.String("user.email", input.Email),
        ),
    )
    defer span.End()

    // Validate input
    if err := s.validateCreateInput(ctx, input); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "validation failed")
        return nil, err
    }

    // Check for existing user
    existing, err := s.checkExistingUser(ctx, input.Email)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    if existing != nil {
        span.SetStatus(codes.Error, "user already exists")
        return nil, ErrUserAlreadyExists
    }

    // Create user
    user, err := s.repo.Create(ctx, input)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, fmt.Errorf("creating user: %w", err)
    }

    span.SetAttributes(attribute.String("user.id", user.ID))
    span.SetStatus(codes.Ok, "user created")

    return user, nil
}

func (s *UserService) validateCreateInput(ctx context.Context, input CreateUserInput) error {
    _, span := s.tracer.Start(ctx, "UserService.validateCreateInput")
    defer span.End()

    if input.Email == "" {
        return ErrEmailRequired
    }
    if input.Password == "" {
        return ErrPasswordRequired
    }

    span.SetStatus(codes.Ok, "validation passed")
    return nil
}

func (s *UserService) checkExistingUser(ctx context.Context, email string) (*User, error) {
    ctx, span := s.tracer.Start(ctx, "UserService.checkExistingUser",
        trace.WithAttributes(attribute.String("email", email)),
    )
    defer span.End()

    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }

    span.SetAttributes(attribute.Bool("exists", user != nil))
    return user, nil
}
```

### Database Tracing

```go
// internal/repository/traced_repo.go
package repository

import (
    "context"
    "database/sql"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

// TracedDB wraps sql.DB with tracing
type TracedDB struct {
    db     *sql.DB
    tracer trace.Tracer
}

func NewTracedDB(db *sql.DB, tracer trace.Tracer) *TracedDB {
    return &TracedDB{
        db:     db,
        tracer: tracer,
    }
}

func (t *TracedDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    ctx, span := t.tracer.Start(ctx, "db.Query",
        trace.WithSpanKind(trace.SpanKindClient),
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", query),
        ),
    )
    defer span.End()

    rows, err := t.db.QueryContext(ctx, query, args...)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    span.SetStatus(codes.Ok, "query successful")
    return rows, nil
}

func (t *TracedDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
    _, span := t.tracer.Start(ctx, "db.QueryRow",
        trace.WithSpanKind(trace.SpanKindClient),
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", query),
        ),
    )
    defer span.End()

    return t.db.QueryRowContext(ctx, query, args...)
}

func (t *TracedDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
    ctx, span := t.tracer.Start(ctx, "db.Exec",
        trace.WithSpanKind(trace.SpanKindClient),
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", query),
        ),
    )
    defer span.End()

    result, err := t.db.ExecContext(ctx, query, args...)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    if rowsAffected, err := result.RowsAffected(); err == nil {
        span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
    }

    span.SetStatus(codes.Ok, "execution successful")
    return result, nil
}
```

---

## 8. API Versioning

### URL Path Versioning

```go
// design/api_versioned.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

// API with version in URL path
var _ = API("myapi", func() {
    Title("My API")
    Description("API with versioning support")
    Version("1.0.0")

    Server("api", func() {
        Host("localhost", func() {
            URI("http://localhost:8080")
        })
    })
})

// V1 Users Service
var _ = Service("users_v1", func() {
    Description("Users service v1")

    HTTP(func() {
        Path("/api/v1/users")
    })

    Method("list", func() {
        Result(ArrayOf(UserV1))
        HTTP(func() {
            GET("/")
            Response(StatusOK)
        })
    })

    Method("get", func() {
        Payload(func() {
            Attribute("id", String, "User ID")
            Required("id")
        })
        Result(UserV1)
        Error("not_found")
        HTTP(func() {
            GET("/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
})

// V2 Users Service with enhanced features
var _ = Service("users_v2", func() {
    Description("Users service v2 with enhanced features")

    HTTP(func() {
        Path("/api/v2/users")
    })

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
                Enum("created_at", "name", "email")
                Default("created_at")
            })
        })
        Result(UserListV2)
        HTTP(func() {
            GET("/")
            Param("page")
            Param("limit")
            Param("sort")
            Response(StatusOK)
        })
    })

    Method("get", func() {
        Payload(func() {
            Attribute("id", String, "User ID")
            Attribute("include", ArrayOf(String), "Related resources to include")
            Required("id")
        })
        Result(UserV2)
        Error("not_found")
        HTTP(func() {
            GET("/{id}")
            Param("include")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
})

// User types for different versions
var UserV1 = Type("UserV1", func() {
    Attribute("id", String)
    Attribute("name", String)
    Attribute("email", String)
    Required("id", "name", "email")
})

var UserV2 = Type("UserV2", func() {
    Attribute("id", String)
    Attribute("name", String)
    Attribute("email", String)
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("updated_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("profile", UserProfile)
    Attribute("roles", ArrayOf(String))
    Required("id", "name", "email")
})

var UserProfile = Type("UserProfile", func() {
    Attribute("avatar_url", String)
    Attribute("bio", String)
    Attribute("location", String)
})

var UserListV2 = Type("UserListV2", func() {
    Attribute("users", ArrayOf(UserV2))
    Attribute("pagination", Pagination)
    Required("users", "pagination")
})

var Pagination = Type("Pagination", func() {
    Attribute("page", Int)
    Attribute("limit", Int)
    Attribute("total", Int)
    Attribute("total_pages", Int)
    Required("page", "limit", "total", "total_pages")
})
```

### Header-based Versioning

```go
// design/header_versioning.go
package design

import (
    . "goa.design/goa/v3/dsl"
)

var _ = Service("users", func() {
    Description("Users service with header-based versioning")

    HTTP(func() {
        Path("/api/users")
    })

    // List users - version determined by header
    Method("list", func() {
        Payload(func() {
            Attribute("api_version", String, "API version", func() {
                Enum("v1", "v2")
                Default("v1")
            })
        })
        Result(Any) // Returns different types based on version
        HTTP(func() {
            GET("/")
            Header("api_version:X-API-Version")
            Response(StatusOK)
        })
    })
})

// Middleware to handle version routing
```

### Version Router Implementation

```go
// internal/router/versioned.go
package router

import (
    "net/http"
    "strings"
)

// VersionedRouter routes requests based on API version
type VersionedRouter struct {
    v1Handler http.Handler
    v2Handler http.Handler
    v3Handler http.Handler
}

func NewVersionedRouter(v1, v2, v3 http.Handler) *VersionedRouter {
    return &VersionedRouter{
        v1Handler: v1,
        v2Handler: v2,
        v3Handler: v3,
    }
}

func (r *VersionedRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Check header first
    version := req.Header.Get("X-API-Version")

    // Fall back to URL path
    if version == "" {
        version = extractVersionFromPath(req.URL.Path)
    }

    // Fall back to query parameter
    if version == "" {
        version = req.URL.Query().Get("version")
    }

    // Default to v1
    if version == "" {
        version = "v1"
    }

    // Route to appropriate handler
    switch version {
    case "v3":
        r.v3Handler.ServeHTTP(w, req)
    case "v2":
        r.v2Handler.ServeHTTP(w, req)
    default:
        r.v1Handler.ServeHTTP(w, req)
    }
}

func extractVersionFromPath(path string) string {
    parts := strings.Split(path, "/")
    for _, part := range parts {
        if strings.HasPrefix(part, "v") && len(part) == 2 {
            return part
        }
    }
    return ""
}
```

### Version-Aware Service

```go
// internal/service/versioned_user.go
package service

import (
    "context"

    usersv1 "your-project/gen/users_v1"
    usersv2 "your-project/gen/users_v2"
)

// UserServiceV1 implements v1 API
type UserServiceV1 struct {
    repo UserRepository
}

func NewUserServiceV1(repo UserRepository) *UserServiceV1 {
    return &UserServiceV1{repo: repo}
}

func (s *UserServiceV1) List(ctx context.Context) ([]*usersv1.UserV1, error) {
    users, err := s.repo.FindAll(ctx)
    if err != nil {
        return nil, err
    }

    result := make([]*usersv1.UserV1, len(users))
    for i, u := range users {
        result[i] = &usersv1.UserV1{
            ID:    u.ID,
            Name:  u.Name,
            Email: u.Email,
        }
    }
    return result, nil
}

func (s *UserServiceV1) Get(ctx context.Context, p *usersv1.GetPayload) (*usersv1.UserV1, error) {
    user, err := s.repo.FindByID(ctx, p.ID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, usersv1.MakeNotFound(nil)
    }

    return &usersv1.UserV1{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

// UserServiceV2 implements v2 API with enhanced features
type UserServiceV2 struct {
    repo UserRepository
}

func NewUserServiceV2(repo UserRepository) *UserServiceV2 {
    return &UserServiceV2{repo: repo}
}

func (s *UserServiceV2) List(ctx context.Context, p *usersv2.ListPayload) (*usersv2.UserListV2, error) {
    // Apply pagination
    page := 1
    if p.Page != nil {
        page = *p.Page
    }
    limit := 20
    if p.Limit != nil {
        limit = *p.Limit
    }
    sort := "created_at"
    if p.Sort != nil {
        sort = *p.Sort
    }

    users, total, err := s.repo.FindAllPaginated(ctx, page, limit, sort)
    if err != nil {
        return nil, err
    }

    result := make([]*usersv2.UserV2, len(users))
    for i, u := range users {
        result[i] = s.toV2(u)
    }

    totalPages := (total + limit - 1) / limit

    return &usersv2.UserListV2{
        Users: result,
        Pagination: &usersv2.Pagination{
            Page:       page,
            Limit:      limit,
            Total:      total,
            TotalPages: totalPages,
        },
    }, nil
}

func (s *UserServiceV2) Get(ctx context.Context, p *usersv2.GetPayload) (*usersv2.UserV2, error) {
    user, err := s.repo.FindByID(ctx, p.ID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, usersv2.MakeNotFound(nil)
    }

    result := s.toV2(user)

    // Handle includes
    if p.Include != nil {
        for _, inc := range p.Include {
            switch inc {
            case "profile":
                // Profile is already included
            case "roles":
                // Roles are already included
            }
        }
    }

    return result, nil
}

func (s *UserServiceV2) toV2(u *User) *usersv2.UserV2 {
    result := &usersv2.UserV2{
        ID:        u.ID,
        Name:      u.Name,
        Email:     u.Email,
        CreatedAt: ptr(u.CreatedAt.Format(time.RFC3339)),
        UpdatedAt: ptr(u.UpdatedAt.Format(time.RFC3339)),
        Roles:     u.Roles,
    }

    if u.Profile != nil {
        result.Profile = &usersv2.UserProfile{
            AvatarURL: u.Profile.AvatarURL,
            Bio:       u.Profile.Bio,
            Location:  u.Profile.Location,
        }
    }

    return result
}

func ptr(s string) *string {
    return &s
}
```

---

## 9. Complete Production-Ready Example

### Main Entry Point

```go
// cmd/api/main.go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "your-project/gen/health"
    healthsrv "your-project/gen/http/health/server"
    usersv1 "your-project/gen/http/users_v1/server"
    "your-project/internal/config"
    "your-project/internal/database"
    "your-project/internal/logger"
    "your-project/internal/metrics"
    "your-project/internal/middleware"
    "your-project/internal/service"
    "your-project/internal/telemetry"

    goahttp "goa.design/goa/v3/http"
    httpmdlwr "goa.design/goa/v3/http/middleware"
)

var (
    version   = "dev"
    buildTime = "unknown"
    commitSHA = "unknown"
)

func main() {
    ctx := context.Background()

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
    }

    // Initialize logger
    log := logger.NewLogger(logger.Config{
        Level:       cfg.Observability.LogLevel,
        Pretty:      cfg.Observability.LogFormat == "pretty",
        ServiceName: "api",
        Environment: cfg.Environment,
    })

    log.Info(ctx).
        Str("version", version).
        Str("build_time", buildTime).
        Str("commit", commitSHA).
        Msg("starting service")

    // Initialize tracing
    tracer, err := telemetry.InitTracer(ctx, telemetry.TracerConfig{
        ServiceName:    cfg.Observability.TracingServiceName,
        ServiceVersion: version,
        Environment:    cfg.Environment,
        Endpoint:       cfg.Observability.TracingEndpoint,
        SampleRate:     cfg.Observability.TracingSampleRate,
        Enabled:        cfg.Observability.TracingEnabled,
    })
    if err != nil {
        log.Fatal(ctx).Err(err).Msg("failed to initialize tracer")
    }
    defer tracer.Shutdown(ctx)

    // Initialize metrics
    m := metrics.NewMetrics("api", "http")

    // Initialize database
    db, err := database.Connect(ctx, cfg.Database)
    if err != nil {
        log.Fatal(ctx).Err(err).Msg("failed to connect to database")
    }
    defer db.Close()

    // Register database metrics collector
    dbCollector := metrics.NewDBCollector(db, "api")
    if err := dbCollector.Register(); err != nil {
        log.Warn(ctx).Err(err).Msg("failed to register db collector")
    }

    // Initialize services
    healthSvc := service.NewHealthService(version)
    healthSvc.RegisterChecker(health.NewPostgresChecker(db))

    userSvc := service.NewUserServiceV1(
        repository.NewUserRepository(db),
    )

    // Create Goa endpoints
    healthEndpoints := health.NewEndpoints(healthSvc)
    userEndpoints := usersv1.NewEndpoints(userSvc)

    // Create HTTP mux
    mux := goahttp.NewMuxer()

    // Mount health endpoints
    healthServer := healthsrv.New(healthEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
    healthsrv.Mount(mux, healthServer)

    // Mount user endpoints
    userServer := usersv1.New(userEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
    usersv1.Mount(mux, userServer)

    // Build middleware chain
    var handler http.Handler = mux

    // Add middlewares (order matters - first added is outermost)
    handler = middleware.MetricsMiddleware(m)(handler)
    handler = middleware.TracingMiddleware(tracer.Tracer())(handler)
    handler = middleware.LoggingMiddleware(log)(handler)
    handler = middleware.RecoveryMiddleware(log)(handler)
    handler = httpmdlwr.RequestID()(handler)

    // Create HTTP server
    server := &http.Server{
        Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
        Handler:      handler,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
    }

    // Start metrics server
    metricsServer := &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.Observability.MetricsPort),
        Handler: metrics.Handler(),
    }
    go func() {
        log.Info(ctx).
            Int("port", cfg.Observability.MetricsPort).
            Msg("metrics server starting")
        if err := metricsServer.ListenAndServe(); err != http.ErrServerClosed {
            log.Error(ctx).Err(err).Msg("metrics server error")
        }
    }()

    // Start main server
    serverErrors := make(chan error, 1)
    go func() {
        log.Info(ctx).
            Str("addr", server.Addr).
            Msg("HTTP server starting")
        serverErrors <- server.ListenAndServe()
    }()

    // Wait for shutdown signal
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    select {
    case err := <-serverErrors:
        log.Fatal(ctx).Err(err).Msg("server error")

    case sig := <-shutdown:
        log.Info(ctx).
            Str("signal", sig.String()).
            Msg("shutdown signal received")

        // Graceful shutdown
        shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
        defer cancel()

        // Shutdown metrics server
        if err := metricsServer.Shutdown(shutdownCtx); err != nil {
            log.Error(ctx).Err(err).Msg("metrics server shutdown error")
        }

        // Shutdown main server
        if err := server.Shutdown(shutdownCtx); err != nil {
            log.Error(ctx).Err(err).Msg("server shutdown error")
            server.Close()
        }

        log.Info(ctx).Msg("server stopped gracefully")
    }
}
```

### Recovery Middleware

```go
// internal/middleware/recovery.go
package middleware

import (
    "context"
    "net/http"
    "runtime/debug"

    "your-project/internal/logger"
)

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    ctx := r.Context()
                    stack := debug.Stack()

                    log.Error(ctx).
                        Interface("panic", rec).
                        Str("stack", string(stack)).
                        Str("method", r.Method).
                        Str("path", r.URL.Path).
                        Msg("panic recovered")

                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 10. Best Practices & Checklist

### Production Readiness Checklist

```markdown
## Pre-Production Checklist

### Configuration
- [ ] All secrets stored in environment variables or secret manager
- [ ] Configuration validation on startup
- [ ] Environment-specific configuration files
- [ ] Default values for non-critical settings

### Logging
- [ ] Structured JSON logging in production
- [ ] Log levels configurable via environment
- [ ] Request ID propagation
- [ ] Sensitive data redaction
- [ ] Log aggregation configured (ELK, CloudWatch, etc.)

### Health Checks
- [ ] Liveness endpoint (/health/live)
- [ ] Readiness endpoint (/health/ready)
- [ ] Detailed status endpoint for debugging
- [ ] All critical dependencies checked

### Graceful Shutdown
- [ ] Signal handling (SIGTERM, SIGINT)
- [ ] In-flight request completion
- [ ] Database connection cleanup
- [ ] Background worker shutdown
- [ ] Configurable shutdown timeout

### Observability
- [ ] Prometheus metrics exposed
- [ ] Request rate, latency, error rate metrics
- [ ] Business metrics (orders, revenue, etc.)
- [ ] Database connection pool metrics
- [ ] Distributed tracing enabled
- [ ] Trace context propagation
- [ ] Grafana dashboards configured

### Security
- [ ] TLS/HTTPS enabled
- [ ] Authentication middleware
- [ ] Authorization checks
- [ ] Rate limiting
- [ ] CORS configuration
- [ ] Security headers

### Docker
- [ ] Multi-stage build for minimal image size
- [ ] Non-root user
- [ ] Health check in Dockerfile
- [ ] No secrets in image
- [ ] Proper signal handling (PID 1)

### API
- [ ] Versioning strategy defined
- [ ] Deprecation headers for old versions
- [ ] Clear error responses
- [ ] Request validation
- [ ] API documentation (OpenAPI)

### Database
- [ ] Connection pooling configured
- [ ] Migrations automated
- [ ] Query timeouts set
- [ ] Retry logic for transient errors

### Testing
- [ ] Unit test coverage > 80%
- [ ] Integration tests
- [ ] Load testing performed
- [ ] Chaos testing (optional)
```

### Common Anti-Patterns to Avoid

```go
// ❌ BAD: Logging sensitive data
log.Info(ctx).
    Str("password", user.Password).  // Never log passwords!
    Str("credit_card", card.Number). // Never log PII!
    Msg("user created")

// ✅ GOOD: Redact sensitive data
log.Info(ctx).
    Str("user_id", user.ID).
    Str("email_domain", extractDomain(user.Email)).
    Msg("user created")

// ❌ BAD: No timeout on shutdown
server.Shutdown(context.Background()) // Could hang forever

// ✅ GOOD: Shutdown with timeout
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
server.Shutdown(ctx)

// ❌ BAD: High cardinality metrics labels
metrics.HTTPRequestDuration.WithLabelValues(
    r.Method,
    r.URL.Path,     // Includes IDs, causes explosion
    userID,         // Unique per user!
).Observe(duration)

// ✅ GOOD: Normalized labels
metrics.HTTPRequestDuration.WithLabelValues(
    r.Method,
    normalizePath(r.URL.Path), // /users/:id
    status,
).Observe(duration)

// ❌ BAD: Ignoring errors in defer
defer func() {
    db.Close() // Error ignored!
}()

// ✅ GOOD: Handle errors in defer
defer func() {
    if err := db.Close(); err != nil {
        log.Error(ctx).Err(err).Msg("error closing database")
    }
}()

// ❌ BAD: Panic in production code
if user == nil {
    panic("user not found") // Don't panic!
}

// ✅ GOOD: Return proper errors
if user == nil {
    return nil, ErrUserNotFound
}
```

### Performance Optimization Tips

```go
// 1. Use connection pooling
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)

// 2. Use context timeouts
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// 3. Batch database operations
func (r *Repository) CreateMany(ctx context.Context, users []*User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx, `INSERT INTO users (name, email) VALUES ($1, $2)`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, u := range users {
        _, err := stmt.ExecContext(ctx, u.Name, u.Email)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

// 4. Use sync.Pool for frequently allocated objects
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processRequest(data []byte) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    // Use buffer...
}

// 5. Parallel health checks
func (s *HealthService) checkAll(ctx context.Context) map[string]Status {
    results := make(map[string]Status)
    var mu sync.Mutex
    var wg sync.WaitGroup

    for name, checker := range s.checkers {
        wg.Add(1)
        go func(n string, c Checker) {
            defer wg.Done()
            status := c.Check(ctx)
            mu.Lock()
            results[n] = status
            mu.Unlock()
        }(name, checker)
    }

    wg.Wait()
    return results
}
```

---

## Summary

This guide covered essential production-readiness topics for Goa services:

1. **Structured Logging** - JSON logging with zerolog/zap, request correlation
2. **Graceful Shutdown** - Signal handling, connection draining, resource cleanup
3. **Configuration** - Environment variables, config files, validation
4. **Docker** - Multi-stage builds, security best practices
5. **Health Checks** - Liveness, readiness, component health
6. **Metrics** - Prometheus integration, custom metrics, dashboards
7. **Tracing** - OpenTelemetry setup, context propagation
8. **API Versioning** - URL path and header-based strategies

### Key Takeaways

- Always use structured logging in production
- Implement graceful shutdown to avoid data loss
- Expose health endpoints for orchestration systems
- Collect metrics for observability and alerting
- Use distributed tracing for debugging complex flows
- Plan your versioning strategy from the start
- Follow the production checklist before deployment

### Next Steps

- **Part 12**: Code Generation Deep Dive
- **Part 13**: Goa CLI & Tooling
- **Part 14**: Microservices Patterns with Goa
- **Part 15**: Scaling & Enterprise Topics
