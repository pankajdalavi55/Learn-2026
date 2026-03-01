# Senior Golang Interview Question Handbook - Part 2

> **Target Audience:** Experienced Golang Developers (5–12 years)  
> **Purpose:** Prepare for senior backend/product company interviews  
> **Difficulty Legend:** 🟢 Basic | 🟡 Intermediate | 🔴 Advanced | ⚫ Expert

---

## Table of Contents - Part 2

5. [Error Handling & Logging](#5-error-handling--logging)
6. [Go Modules & Dependency Management](#6-go-modules--dependency-management)
7. [Testing in Go](#7-testing-in-go)
8. [Building Production APIs in Go](#8-building-production-apis-in-go)
9. [Database & Caching](#9-database--caching)

---

## 5. Error Handling & Logging

### 5.1 Error Fundamentals

#### 🟢 Q46: What is the idiomatic way to handle errors in Go?

**Answer:**

```go
// Check error immediately after call
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}
// Use result...
```

**Principles:**
- Handle errors at call site, not centrally
- Don't ignore errors (use `_` only when intentional)
- Add context when propagating
- Return early on error (happy path at bottom)

---

#### 🟡 Q47: Explain error wrapping with %w and errors.Is/As.

**Answer:**

**Error wrapping (Go 1.13+):**
```go
func readConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // Wrap with context, preserve original
        return nil, fmt.Errorf("reading config %s: %w", path, err)
    }
    // ...
}
```

**Unwrapping with errors.Is:**
```go
_, err := readConfig("app.yaml")
if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
    return useDefaults()
}
```

**Type assertion with errors.As:**
```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Failed path:", pathErr.Path)
    fmt.Println("Operation:", pathErr.Op)
}
```

**Difference:**
- `errors.Is` — checks error identity (sentinel errors)
- `errors.As` — checks error type (custom error types)

---

#### 🔴 Q48: How do you create custom error types?

**Answer:**

**Method 1: Sentinel errors**
```go
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrValidation   = errors.New("validation failed")
)

func FindUser(id string) (*User, error) {
    user := db.Find(id)
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}

// Usage
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

**Method 2: Custom error type with context**
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Constructor for consistent creation
func NewValidationError(field, msg string) error {
    return &ValidationError{Field: field, Message: msg}
}

// Usage with errors.As
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Invalid field: %s\n", valErr.Field)
}
```

**Method 3: Error with wrapped cause**
```go
type QueryError struct {
    Query string
    Err   error
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query %q failed: %v", e.Query, e.Err)
}

// Implement Unwrap for errors.Is/As to work
func (e *QueryError) Unwrap() error {
    return e.Err
}

// Usage
err := &QueryError{Query: "SELECT...", Err: sql.ErrNoRows}
errors.Is(err, sql.ErrNoRows) // true
```

**Method 4: Multi-error (Go 1.20+)**
```go
func (e *MultiError) Error() string {
    return "multiple errors occurred"
}

func (e *MultiError) Unwrap() []error {
    return e.errs
}

// errors.Is checks all wrapped errors
```

---

### 5.2 Panic vs Error

#### 🟡 Q49: When should you use panic vs returning error?

**Answer:**

| Use `panic` | Use `error` |
|-------------|-------------|
| Unrecoverable programmer error | Expected failure (file not found, network timeout) |
| Nil pointer that shouldn't be nil | User input validation |
| Index out of bounds (bug) | Business logic failures |
| Initialization failures (must succeed) | External service errors |

```go
// panic: Programmer error
func MustCompile(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil {
        panic(fmt.Sprintf("invalid regex %q: %v", pattern, err))
    }
    return re
}

// error: Expected failure
func Compile(pattern string) (*regexp.Regexp, error) {
    return regexp.Compile(pattern)
}

// panic: Invariant violation
func (s *Stack) Pop() int {
    if len(s.items) == 0 {
        panic("pop from empty stack") // Bug in caller
    }
    // ...
}
```

**Rule of thumb:** If the caller can reasonably handle it, return error.

---

#### 🔴 Q50: How does recover work and when should you use it?

**Answer:**

**recover() only works inside deferred function:**
```go
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
            // Optionally log stack trace
            debug.PrintStack()
        }
    }()
    
    fn()
    return nil
}
```

**Valid use cases:**
1. **HTTP handlers** — prevent one request from crashing server
2. **Goroutine boundaries** — prevent panic from taking down all goroutines
3. **Plugin/library boundaries** — isolate third-party code

```go
// HTTP middleware for panic recovery
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

**Antipatterns:**
```go
// WRONG: Using panic/recover for control flow
func find(items []int, target int) int {
    defer func() { recover() }()
    panic("not found") // Don't do this
}

// WRONG: Recovering and ignoring
defer func() { recover() }() // Swallows all panics silently
```

---

### 5.3 Error Handling Patterns

#### 🔴 Q51: What are common error handling patterns in production Go code?

**Answer:**

**Pattern 1: Error types for handling logic**
```go
type ErrorKind int

const (
    ErrKindUnknown ErrorKind = iota
    ErrKindNotFound
    ErrKindValidation
    ErrKindPermission
    ErrKindConflict
)

type AppError struct {
    Kind    ErrorKind
    Message string
    Err     error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

// HTTP handler uses Kind to determine status code
func errorToStatus(err error) int {
    var appErr *AppError
    if errors.As(err, &appErr) {
        switch appErr.Kind {
        case ErrKindNotFound:
            return http.StatusNotFound
        case ErrKindValidation:
            return http.StatusBadRequest
        case ErrKindPermission:
            return http.StatusForbidden
        case ErrKindConflict:
            return http.StatusConflict
        }
    }
    return http.StatusInternalServerError
}
```

**Pattern 2: Functional error handling**
```go
func (s *Service) CreateUser(ctx context.Context, req CreateUserReq) (*User, error) {
    // Validate
    if err := req.Validate(); err != nil {
        return nil, &AppError{Kind: ErrKindValidation, Message: err.Error(), Err: err}
    }
    
    // Check existence
    existing, err := s.repo.FindByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, ErrNotFound) {
        return nil, fmt.Errorf("checking email: %w", err)
    }
    if existing != nil {
        return nil, &AppError{Kind: ErrKindConflict, Message: "email already exists"}
    }
    
    // Create
    user, err := s.repo.Create(ctx, req.ToUser())
    if err != nil {
        return nil, fmt.Errorf("creating user: %w", err)
    }
    
    return user, nil
}
```

**Pattern 3: Error aggregation**
```go
type MultiError struct {
    errors []error
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.errors = append(m.errors, err)
    }
}

func (m *MultiError) Error() string {
    var msgs []string
    for _, e := range m.errors {
        msgs = append(msgs, e.Error())
    }
    return strings.Join(msgs, "; ")
}

func (m *MultiError) ErrorOrNil() error {
    if len(m.errors) == 0 {
        return nil
    }
    return m
}

// Usage
func validateUser(u User) error {
    var errs MultiError
    if u.Name == "" {
        errs.Add(errors.New("name is required"))
    }
    if u.Email == "" {
        errs.Add(errors.New("email is required"))
    }
    return errs.ErrorOrNil()
}
```

---

### 5.4 Observability Best Practices

#### 🔴 Q52: What are logging best practices in production Go services?

**Answer:**

**1. Structured logging:**
```go
import "log/slog"

// Setup
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
slog.SetDefault(logger)

// Usage
slog.Info("user created",
    slog.String("user_id", user.ID),
    slog.String("email", user.Email),
    slog.Duration("latency", time.Since(start)),
)
// Output: {"time":"...","level":"INFO","msg":"user created","user_id":"123","email":"...","latency":"5ms"}
```

**2. Context-aware logging:**
```go
func LoggerFromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

// Middleware adds request context to logger
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        logger := slog.Default().With(
            slog.String("request_id", requestID),
            slog.String("method", r.Method),
            slog.String("path", r.URL.Path),
        )
        ctx := WithLogger(r.Context(), logger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**3. Error logging with context:**
```go
func (s *Service) ProcessOrder(ctx context.Context, orderID string) error {
    logger := LoggerFromContext(ctx).With(slog.String("order_id", orderID))
    
    order, err := s.repo.GetOrder(ctx, orderID)
    if err != nil {
        logger.Error("failed to get order", slog.Any("error", err))
        return fmt.Errorf("getting order: %w", err)
    }
    
    logger.Info("processing order", slog.Float64("amount", order.Amount))
    // ...
}
```

**4. Log levels:**
```go
// DEBUG - Development, verbose tracing
slog.Debug("cache miss", slog.String("key", key))

// INFO - Business events, request lifecycle
slog.Info("order completed", slog.String("order_id", id))

// WARN - Degraded but recoverable
slog.Warn("retry attempt", slog.Int("attempt", 3))

// ERROR - Failures requiring attention
slog.Error("payment failed", slog.Any("error", err))
```

---

## 6. Go Modules & Dependency Management

### 6.1 Module Basics

#### 🟢 Q53: Explain go.mod and go.sum files.

**Answer:**

**go.mod:**
```go
module github.com/myorg/myapp   // Module path

go 1.21                         // Go version

require (
    github.com/gin-gonic/gin v1.9.1        // Direct dependency
    golang.org/x/sync v0.3.0               // Direct dependency
)

require (
    github.com/bytedance/sonic v1.9.1 // indirect  // Transitive dependency
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0  // Replace directive

exclude github.com/broken/pkg v1.2.3  // Exclude version
```

**go.sum:**
- Cryptographic checksums for module versions
- Ensures reproducible builds
- Contains checksums for both `.mod` files and module zip archives

```
github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkP...
github.com/gin-gonic/gin v1.9.1/go.mod h1:RdlIXY9OAW...
```

**Commands:**
```bash
go mod init github.com/myorg/myapp  # Initialize module
go mod tidy                          # Add missing, remove unused
go mod download                      # Download dependencies
go mod verify                        # Verify checksums
go mod why github.com/pkg            # Why is this dependency needed?
go mod graph                         # Print dependency graph
```

---

#### 🟡 Q54: Explain semantic versioning in Go modules.

**Answer:**

**Format:** `vMAJOR.MINOR.PATCH`

| Component | Meaning | Example |
|-----------|---------|---------|
| MAJOR | Breaking changes | v2.0.0 |
| MINOR | New features, backward compatible | v1.1.0 |
| PATCH | Bug fixes, backward compatible | v1.0.1 |

**Go-specific rules:**

**Major version in import path (v2+):**
```go
// v0 and v1
import "github.com/myorg/pkg"

// v2 and above MUST include major version
import "github.com/myorg/pkg/v2"
```

**Module structure for v2+:**
```
mymodule/
├── go.mod          // module github.com/myorg/pkg/v2
├── main.go
└── v2/             // Alternative: subdirectory approach
    └── go.mod
```

**Pre-release and metadata:**
```
v1.0.0-alpha        // Pre-release
v1.0.0-beta.1       // Pre-release with number
v1.0.0+metadata     // Build metadata (ignored in version selection)
```

**Version selection (MVS - Minimal Version Selection):**
```go
// If A requires B v1.2.0 and C requires B v1.3.0
// Go selects B v1.3.0 (minimum version satisfying all)
```

---

#### 🔴 Q55: How do you handle private modules?

**Answer:**

**1. Configure GOPRIVATE:**
```bash
# Single private domain
go env -w GOPRIVATE=github.com/myorg

# Multiple domains
go env -w GOPRIVATE=github.com/myorg,gitlab.internal.com

# Wildcards
go env -w GOPRIVATE=*.internal.com
```

**2. Git authentication:**
```bash
# Option A: .netrc file
# ~/.netrc (Linux/Mac) or %USERPROFILE%\_netrc (Windows)
machine github.com login USERNAME password TOKEN

# Option B: Git credential helper
git config --global url."https://TOKEN:x-oauth-basic@github.com/myorg".insteadOf "https://github.com/myorg"

# Option C: SSH
git config --global url."git@github.com:".insteadOf "https://github.com/myorg/"
```

**3. Related environment variables:**
```bash
GOPRIVATE   # Skip proxy and checksum for these modules
GONOPROXY   # Skip proxy (but still use checksum DB)
GONOSUMDB   # Skip checksum DB
GOPROXY     # Module proxy URL (default: https://proxy.golang.org)
```

**4. Vendoring for CI/CD:**
```bash
go mod vendor            # Copy dependencies to vendor/
go build -mod=vendor     # Build using vendor directory
```

---

### 6.2 Dependency Management

#### 🔴 Q56: How do you handle dependency conflicts and upgrades?

**Answer:**

**View dependencies:**
```bash
go list -m all                    # All dependencies
go list -m -versions github.com/pkg  # Available versions
go mod graph                      # Dependency graph
go mod why github.com/pkg         # Why is it included?
```

**Upgrade dependencies:**
```bash
go get -u ./...                   # Update all direct dependencies
go get -u=patch ./...             # Only patch updates
go get github.com/pkg@v1.2.3      # Specific version
go get github.com/pkg@latest      # Latest version
go get github.com/pkg@master      # Specific branch
```

**Handle conflicts:**
```go
// Force specific version via replace
replace github.com/pkg => github.com/pkg v1.2.3

// Use fork
replace github.com/original/pkg => github.com/myfork/pkg v1.0.0

// Local development
replace github.com/pkg => ../local-pkg
```

**Exclude problematic versions:**
```go
exclude github.com/broken/pkg v1.2.3
```

**Debug version selection:**
```bash
go mod why -m github.com/pkg     # Why is this module here?
go mod graph | grep github.com/pkg  # Who depends on it?
```

---

## 7. Testing in Go

### 7.1 Testing Fundamentals

#### 🟢 Q57: Explain table-driven tests in Go.

**Answer:**

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

**Benefits:**
- Easy to add new cases
- Each case has a descriptive name
- Parallel execution with `t.Parallel()`
- Clear failure messages

**Parallel table tests:**
```go
for _, tt := range tests {
    tt := tt // Capture range variable
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        // Test code...
    })
}
```

---

#### 🟡 Q58: How do you use test fixtures and helpers?

**Answer:**

**Test helpers:**
```go
// Helper function (t.Helper() marks it as helper for better error reporting)
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v, want %v", got, want)
    }
}

// Setup/teardown helper
func setupTestDB(t *testing.T) (*sql.DB, func()) {
    t.Helper()
    db, err := sql.Open("postgres", testDSN)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    
    cleanup := func() {
        db.Exec("DELETE FROM users WHERE email LIKE '%@test.com'")
        db.Close()
    }
    
    return db, cleanup
}

func TestUserRepository(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()
    
    // Tests using db...
}
```

**Test fixtures (testdata directory):**
```
mypackage/
├── parser.go
├── parser_test.go
└── testdata/
    ├── valid_input.json
    ├── invalid_input.json
    └── expected_output.json
```

```go
func TestParser(t *testing.T) {
    input, err := os.ReadFile("testdata/valid_input.json")
    if err != nil {
        t.Fatalf("failed to read fixture: %v", err)
    }
    
    result, err := Parse(input)
    // Assert...
}
```

**Golden files:**
```go
var update = flag.Bool("update", false, "update golden files")

func TestGolden(t *testing.T) {
    result := GenerateOutput()
    golden := filepath.Join("testdata", t.Name()+".golden")
    
    if *update {
        os.WriteFile(golden, result, 0644)
        return
    }
    
    expected, _ := os.ReadFile(golden)
    if !bytes.Equal(result, expected) {
        t.Errorf("output mismatch; run with -update to regenerate")
    }
}
```

---

### 7.2 Mocking

#### 🔴 Q59: What are the mocking approaches in Go?

**Answer:**

**Approach 1: Interface-based mocking (preferred)**
```go
// Define interface for dependency
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// Production implementation
type PostgresUserRepo struct {
    db *sql.DB
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
    // Real DB query
}

// Test mock
type MockUserRepo struct {
    GetByIDFunc func(ctx context.Context, id string) (*User, error)
    SaveFunc    func(ctx context.Context, user *User) error
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
    return m.GetByIDFunc(ctx, id)
}

func (m *MockUserRepo) Save(ctx context.Context, user *User) error {
    return m.SaveFunc(ctx, user)
}

// Test
func TestUserService_GetUser(t *testing.T) {
    mock := &MockUserRepo{
        GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
            return &User{ID: id, Name: "Test"}, nil
        },
    }
    
    service := NewUserService(mock)
    user, err := service.GetUser(context.Background(), "123")
    
    if err != nil || user.Name != "Test" {
        t.Errorf("unexpected result")
    }
}
```

**Approach 2: Using testify/mock**
```go
import "github.com/stretchr/testify/mock"

type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestWithTestify(t *testing.T) {
    mockRepo := new(MockUserRepo)
    mockRepo.On("GetByID", mock.Anything, "123").Return(&User{Name: "Test"}, nil)
    
    service := NewUserService(mockRepo)
    user, _ := service.GetUser(context.Background(), "123")
    
    assert.Equal(t, "Test", user.Name)
    mockRepo.AssertExpectations(t)
}
```

**Approach 3: HTTP mocking**
```go
func TestAPIClient(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/users/123" {
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(`{"id":"123","name":"Test"}`))
            return
        }
        w.WriteHeader(http.StatusNotFound)
    }))
    defer server.Close()
    
    client := NewAPIClient(server.URL)
    user, err := client.GetUser("123")
    
    assert.NoError(t, err)
    assert.Equal(t, "Test", user.Name)
}
```

---

### 7.3 Benchmarks and Race Detection

#### 🟡 Q60: How do you write benchmark tests?

**Answer:**

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// With setup (excluded from timing)
func BenchmarkProcess(b *testing.B) {
    data := generateTestData(1000)
    
    b.ResetTimer() // Reset after setup
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// Sub-benchmarks with different inputs
func BenchmarkSort(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            data := generateData(size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                sort.Ints(data)
            }
        })
    }
}

// Memory allocation tracking
func BenchmarkWithAllocs(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}
```

**Running benchmarks:**
```bash
go test -bench=.                      # Run all benchmarks
go test -bench=BenchmarkFib           # Run specific benchmark
go test -bench=. -benchmem            # Include memory stats
go test -bench=. -benchtime=5s        # Run for 5 seconds
go test -bench=. -count=5             # Run 5 times
go test -bench=. -cpuprofile=cpu.out  # Generate CPU profile
```

**Output:**
```
BenchmarkFibonacci-8    5000000    234 ns/op    0 B/op    0 allocs/op
```

---

#### 🔴 Q61: Explain the race detector and how to use it.

**Answer:**

**Enable race detector:**
```bash
go test -race ./...          # Run tests with race detection
go build -race               # Build with race detection
go run -race main.go         # Run with race detection
```

**What it detects:**
```go
// Data race example
var counter int

func increment() {
    counter++ // Race: concurrent read/write
}

func TestRace(t *testing.T) {
    for i := 0; i < 100; i++ {
        go increment()
    }
    time.Sleep(time.Second)
    fmt.Println(counter)
}
```

**Race detector output:**
```
WARNING: DATA RACE
Write at 0x00c0000a4010 by goroutine 8:
  main.increment()
      main.go:10 +0x38

Previous read at 0x00c0000a4010 by goroutine 7:
  main.increment()
      main.go:10 +0x30

Goroutine 8 (running) created at:
  main.TestRace()
      main.go:15 +0x5c
```

**Fix with sync:**
```go
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// Or use atomic
var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}
```

**Race detector limitations:**
- 2-20x slowdown, 5-10x memory increase
- Only detects races that **actually occur** during execution
- Not suitable for production

---

### 7.4 Test Coverage

#### 🟡 Q62: How do you measure and improve test coverage?

**Answer:**

**Generate coverage:**
```bash
go test -cover ./...                           # Show coverage percentage
go test -coverprofile=coverage.out ./...       # Generate profile
go tool cover -html=coverage.out               # HTML report
go tool cover -func=coverage.out               # Function-level summary
```

**Coverage modes:**
```bash
go test -covermode=set ./...    # Was statement executed? (default)
go test -covermode=count ./...  # How many times executed?
go test -covermode=atomic ./... # Like count, but thread-safe
```

**Cover specific packages:**
```bash
go test -cover -coverpkg=./... ./...  # Include all packages in coverage
```

**Integration with CI:**
```yaml
# GitHub Actions example
- name: Test with coverage
  run: go test -coverprofile=coverage.out ./...

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    file: coverage.out
```

**Coverage targets:**
- 80%+ for critical business logic
- 60-70% for utilities/helpers
- 100% not always practical or necessary
- Focus on critical paths, not metric gaming

---

## 8. Building Production APIs in Go

### 8.1 REST API Best Practices

#### 🟡 Q63: What are best practices for structuring a Go REST API?

**Answer:**

**Project structure:**
```
myapi/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── handler/              # HTTP handlers
│   │   ├── handler.go
│   │   └── user_handler.go
│   ├── service/              # Business logic
│   │   └── user_service.go
│   ├── repository/           # Data access
│   │   └── user_repo.go
│   ├── model/                # Domain models
│   │   └── user.go
│   └── middleware/           # HTTP middleware
│       └── auth.go
├── pkg/                      # Reusable packages
├── api/                      # OpenAPI specs
└── go.mod
```

**Handler pattern:**
```go
type UserHandler struct {
    service UserService
    logger  *slog.Logger
}

func NewUserHandler(s UserService, l *slog.Logger) *UserHandler {
    return &UserHandler{service: s, logger: l}
}

func (h *UserHandler) RegisterRoutes(r *http.ServeMux) {
    r.HandleFunc("GET /users/{id}", h.GetUser)
    r.HandleFunc("POST /users", h.CreateUser)
    r.HandleFunc("PUT /users/{id}", h.UpdateUser)
    r.HandleFunc("DELETE /users/{id}", h.DeleteUser)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // Go 1.22+
    
    user, err := h.service.GetUser(r.Context(), id)
    if err != nil {
        h.handleError(w, r, err)
        return
    }
    
    h.respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        h.respondJSON(w, appErr.StatusCode(), map[string]string{"error": appErr.Message})
        return
    }
    h.logger.Error("unhandled error", slog.Any("error", err))
    h.respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
}
```

---

### 8.2 Middleware

#### 🔴 Q64: How do you implement middleware in Go?

**Answer:**

**Middleware signature:**
```go
type Middleware func(http.Handler) http.Handler
```

**Common middleware implementations:**

**1. Logging:**
```go
func LoggingMiddleware(logger *slog.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer to capture status
            wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
            
            next.ServeHTTP(wrapped, r)
            
            logger.Info("request",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.Int("status", wrapped.status),
                slog.Duration("duration", time.Since(start)),
            )
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (w *responseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}
```

**2. Authentication:**
```go
func AuthMiddleware(validator TokenValidator) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            
            claims, err := validator.Validate(strings.TrimPrefix(token, "Bearer "))
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            
            ctx := context.WithValue(r.Context(), userClaimsKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**3. Rate limiting:**
```go
func RateLimitMiddleware(rps int) Middleware {
    limiter := rate.NewLimiter(rate.Limit(rps), rps)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Chaining middleware:**
```go
func Chain(middlewares ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// Usage
handler := Chain(
    RecoveryMiddleware,
    LoggingMiddleware(logger),
    AuthMiddleware(validator),
)(finalHandler)
```

---

### 8.3 Graceful Shutdown

#### 🔴 Q65: How do you implement graceful shutdown?

**Answer:**

```go
func main() {
    // Setup
    logger := slog.Default()
    handler := setupRoutes()
    
    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Start server in goroutine
    go func() {
        logger.Info("starting server", slog.String("addr", server.Addr))
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server error", slog.Any("error", err))
            os.Exit(1)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    logger.Info("shutting down server...")
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Shutdown HTTP server
    if err := server.Shutdown(ctx); err != nil {
        logger.Error("server shutdown error", slog.Any("error", err))
    }
    
    // Close other resources (DB, cache, etc.)
    closeResources(ctx)
    
    logger.Info("server stopped")
}

func closeResources(ctx context.Context) {
    // Close database connections
    if db != nil {
        db.Close()
    }
    
    // Close Redis connections
    if redisClient != nil {
        redisClient.Close()
    }
    
    // Flush logs, metrics
    // ...
}
```

**Key points:**
- `server.Shutdown()` stops accepting new connections
- Waits for active requests to complete (up to context deadline)
- Close resources in reverse initialization order
- Use buffered channel for signal (prevents missing signal)

---

### 8.4 Configuration Management

#### 🔴 Q66: How do you handle configuration in production Go services?

**Answer:**

**Configuration struct:**
```go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
    Port            int           `yaml:"port" env:"SERVER_PORT" default:"8080"`
    ReadTimeout     time.Duration `yaml:"read_timeout" default:"15s"`
    WriteTimeout    time.Duration `yaml:"write_timeout" default:"15s"`
    ShutdownTimeout time.Duration `yaml:"shutdown_timeout" default:"30s"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST" required:"true"`
    Port     int    `yaml:"port" env:"DB_PORT" default:"5432"`
    User     string `yaml:"user" env:"DB_USER" required:"true"`
    Password string `yaml:"password" env:"DB_PASSWORD" required:"true"`
    Database string `yaml:"database" env:"DB_NAME" required:"true"`
    MaxConns int    `yaml:"max_conns" default:"25"`
}
```

**Loading configuration:**
```go
func LoadConfig(path string) (*Config, error) {
    cfg := &Config{}
    
    // Load from YAML file
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config file: %w", err)
    }
    
    if err := yaml.Unmarshal(data, cfg); err != nil {
        return nil, fmt.Errorf("parsing config: %w", err)
    }
    
    // Override with environment variables
    if err := envconfig.Process("", cfg); err != nil {
        return nil, fmt.Errorf("processing env vars: %w", err)
    }
    
    // Validate
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("validating config: %w", err)
    }
    
    return cfg, nil
}

func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return errors.New("database host is required")
    }
    if c.Database.MaxConns < 1 {
        return errors.New("database max_conns must be positive")
    }
    return nil
}
```

**Configuration precedence:**
1. Default values in struct tags
2. Configuration file (YAML/JSON)
3. Environment variables (highest priority)

**12-Factor App approach:**
- Store config in environment variables
- Never commit secrets to version control
- Use secret managers in production (Vault, AWS Secrets Manager)

---

### 8.5 Dependency Injection

#### 🔴 Q67: How do you implement dependency injection in Go?

**Answer:**

**Manual DI (preferred for most projects):**
```go
// main.go
func main() {
    cfg := loadConfig()
    
    // Build dependencies bottom-up
    db, err := sql.Open("postgres", cfg.Database.DSN())
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Repositories
    userRepo := repository.NewPostgresUserRepo(db)
    orderRepo := repository.NewPostgresOrderRepo(db)
    
    // Services
    userService := service.NewUserService(userRepo)
    orderService := service.NewOrderService(orderRepo, userService)
    
    // Handlers
    userHandler := handler.NewUserHandler(userService)
    orderHandler := handler.NewOrderHandler(orderService)
    
    // Router
    router := http.NewServeMux()
    userHandler.RegisterRoutes(router)
    orderHandler.RegisterRoutes(router)
    
    // Server
    server := &http.Server{
        Addr:    cfg.Server.Addr,
        Handler: router,
    }
    // ...
}
```

**Using wire (compile-time DI):**
```go
// wire.go
//go:build wireinject

package main

import "github.com/google/wire"

func InitializeApp(cfg *Config) (*App, error) {
    wire.Build(
        // Database
        NewDatabase,
        
        // Repositories
        repository.NewUserRepo,
        repository.NewOrderRepo,
        
        // Services
        service.NewUserService,
        service.NewOrderService,
        
        // Handlers
        handler.NewUserHandler,
        handler.NewOrderHandler,
        
        // App
        NewApp,
    )
    return nil, nil
}

// Run: wire gen ./...
```

**Using fx (runtime DI):**
```go
func main() {
    fx.New(
        fx.Provide(
            loadConfig,
            NewDatabase,
            repository.NewUserRepo,
            service.NewUserService,
            handler.NewUserHandler,
        ),
        fx.Invoke(startServer),
    ).Run()
}
```

**Recommendation:** 
- Small projects: Manual DI
- Large projects with many dependencies: Wire (compile-time safety)
- Avoid runtime DI unless you need dynamic behavior

---

## 9. Database & Caching

### 9.1 Database/SQL Package

#### 🟢 Q68: Explain the database/sql package essentials.

**Answer:**

**Opening connection:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq" // Driver (blank import for side effects)
)

db, err := sql.Open("postgres", "postgres://user:pass@localhost/db?sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Verify connection
if err := db.Ping(); err != nil {
    log.Fatal(err)
}
```

**Query patterns:**
```go
// Query single row
var user User
err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).
    Scan(&user.ID, &user.Name, &user.Email)
if err == sql.ErrNoRows {
    return nil, ErrNotFound
}
if err != nil {
    return nil, err
}

// Query multiple rows
rows, err := db.Query("SELECT id, name FROM users WHERE status = $1", status)
if err != nil {
    return nil, err
}
defer rows.Close()

var users []User
for rows.Next() {
    var u User
    if err := rows.Scan(&u.ID, &u.Name); err != nil {
        return nil, err
    }
    users = append(users, u)
}
if err := rows.Err(); err != nil {
    return nil, err
}

// Execute (INSERT, UPDATE, DELETE)
result, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
if err != nil {
    return err
}
id, _ := result.LastInsertId()      // Not supported by all drivers
affected, _ := result.RowsAffected()
```

**With context (always use in production):**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := db.QueryRowContext(ctx, "SELECT ...", args...).Scan(&result)
```

---

### 9.2 Connection Pooling

#### 🔴 Q69: How do you configure connection pooling in Go?

**Answer:**

```go
db, _ := sql.Open("postgres", dsn)

// Connection pool settings
db.SetMaxOpenConns(25)           // Max open connections to database
db.SetMaxIdleConns(10)           // Max idle connections in pool
db.SetConnMaxLifetime(5 * time.Minute)  // Max time a connection can be reused
db.SetConnMaxIdleTime(1 * time.Minute)  // Max time a connection can be idle
```

**Guidelines:**

| Setting | Guideline | Reason |
|---------|-----------|--------|
| `MaxOpenConns` | Start with 25, tune based on load | Prevents overwhelming DB |
| `MaxIdleConns` | ~50% of MaxOpenConns | Balance between ready connections and resources |
| `ConnMaxLifetime` | < DB's wait_timeout | Prevent using stale connections |
| `ConnMaxIdleTime` | 1-5 minutes | Release unused connections |

**Monitoring pool stats:**
```go
stats := db.Stats()
log.Printf("Open: %d, InUse: %d, Idle: %d, WaitCount: %d, WaitDuration: %v",
    stats.OpenConnections,
    stats.InUse,
    stats.Idle,
    stats.WaitCount,
    stats.WaitDuration,
)

// Expose as metrics
poolOpenConnections.Set(float64(stats.OpenConnections))
poolInUse.Set(float64(stats.InUse))
poolWaitCount.Add(float64(stats.WaitCount))
```

**Tuning steps:**
1. Monitor `WaitCount` and `WaitDuration`
2. If high wait times, increase `MaxOpenConns`
3. If connections frequently closed/opened, increase `MaxIdleConns`
4. Watch database server connection limits

---

### 9.3 Transactions

#### 🔴 Q70: How do you handle transactions properly in Go?

**Answer:**

**Basic transaction:**
```go
func TransferFunds(ctx context.Context, db *sql.DB, from, to string, amount float64) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback() // No-op if committed
    
    // Debit
    _, err = tx.ExecContext(ctx, 
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2", 
        amount, from)
    if err != nil {
        return fmt.Errorf("debit: %w", err)
    }
    
    // Credit
    _, err = tx.ExecContext(ctx, 
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2", 
        amount, to)
    if err != nil {
        return fmt.Errorf("credit: %w", err)
    }
    
    // Commit
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    
    return nil
}
```

**Transaction helper:**
```go
func WithTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p) // Re-throw panic after rollback
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// Usage
err := WithTransaction(ctx, db, func(tx *sql.Tx) error {
    if _, err := tx.ExecContext(ctx, "INSERT ...", args...); err != nil {
        return err
    }
    if _, err := tx.ExecContext(ctx, "UPDATE ...", args...); err != nil {
        return err
    }
    return nil
})
```

**Transaction isolation levels:**
```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable, // Strictest
    ReadOnly:  true,                  // Read-only transaction
})

// Levels: LevelDefault, LevelReadUncommitted, LevelReadCommitted, 
//         LevelRepeatableRead, LevelSerializable
```

---

### 9.4 N+1 Query Problem

#### 🔴 Q71: How do you prevent N+1 queries in Go?

**Answer:**

**N+1 Problem:**
```go
// BAD: N+1 queries
func GetUsersWithOrders(ctx context.Context) ([]UserWithOrders, error) {
    users, _ := db.Query("SELECT id, name FROM users") // 1 query
    
    var result []UserWithOrders
    for users.Next() {
        var u User
        users.Scan(&u.ID, &u.Name)
        
        // N queries - one per user!
        orders, _ := db.Query("SELECT * FROM orders WHERE user_id = $1", u.ID)
        u.Orders = scanOrders(orders)
        result = append(result, u)
    }
    return result, nil
}
```

**Solution 1: JOIN**
```go
func GetUsersWithOrders(ctx context.Context) ([]UserWithOrders, error) {
    rows, err := db.Query(`
        SELECT u.id, u.name, o.id, o.product, o.amount
        FROM users u
        LEFT JOIN orders o ON u.id = o.user_id
        ORDER BY u.id
    `)
    // Process and group results
}
```

**Solution 2: Batch loading (DataLoader pattern)**
```go
func GetUsersWithOrders(ctx context.Context, userIDs []string) ([]UserWithOrders, error) {
    // Query 1: Get all users
    users, _ := getUsersByIDs(ctx, userIDs)
    
    // Query 2: Get all orders for these users in one query
    orders, _ := getOrdersByUserIDs(ctx, userIDs)
    
    // Map orders to users in memory
    ordersByUser := make(map[string][]Order)
    for _, o := range orders {
        ordersByUser[o.UserID] = append(ordersByUser[o.UserID], o)
    }
    
    // Combine
    for i := range users {
        users[i].Orders = ordersByUser[users[i].ID]
    }
    
    return users, nil
}

func getOrdersByUserIDs(ctx context.Context, userIDs []string) ([]Order, error) {
    query := "SELECT * FROM orders WHERE user_id = ANY($1)"
    return db.Query(query, pq.Array(userIDs))
}
```

**Solution 3: Use ORM with preloading**
```go
// GORM example
var users []User
db.Preload("Orders").Find(&users)
```

---

### 9.5 Redis Usage Patterns

#### 🔴 Q72: What are common Redis patterns in Go?

**Answer:**

**Setup:**
```go
import "github.com/redis/go-redis/v9"

rdb := redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    Password:     "",
    DB:           0,
    PoolSize:     100,
    MinIdleConns: 10,
})
```

**Pattern 1: Cache-aside**
```go
func GetUser(ctx context.Context, id string) (*User, error) {
    // Try cache first
    cached, err := rdb.Get(ctx, "user:"+id).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    if err != redis.Nil {
        log.Printf("redis error: %v", err)
    }
    
    // Cache miss - fetch from DB
    user, err := db.GetUser(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    data, _ := json.Marshal(user)
    rdb.Set(ctx, "user:"+id, data, 15*time.Minute)
    
    return user, nil
}

func InvalidateUser(ctx context.Context, id string) {
    rdb.Del(ctx, "user:"+id)
}
```

**Pattern 2: Distributed lock**
```go
func AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    // SET key value NX EX seconds
    result, err := rdb.SetNX(ctx, "lock:"+key, "1", ttl).Result()
    return result, err
}

func ReleaseLock(ctx context.Context, key string) {
    rdb.Del(ctx, "lock:"+key)
}

// Usage
if acquired, _ := AcquireLock(ctx, "order:123", 30*time.Second); acquired {
    defer ReleaseLock(ctx, "order:123")
    // Process order...
}
```

**Pattern 3: Rate limiting**
```go
func CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    pipe := rdb.Pipeline()
    
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    return incr.Val() <= int64(limit), nil
}
```

**Pattern 4: Pub/Sub**
```go
// Publisher
func PublishEvent(ctx context.Context, channel string, event Event) error {
    data, _ := json.Marshal(event)
    return rdb.Publish(ctx, channel, data).Err()
}

// Subscriber
func Subscribe(ctx context.Context, channel string, handler func(Event)) {
    sub := rdb.Subscribe(ctx, channel)
    defer sub.Close()
    
    for msg := range sub.Channel() {
        var event Event
        json.Unmarshal([]byte(msg.Payload), &event)
        handler(event)
    }
}
```

**Pattern 5: Session storage**
```go
type SessionStore struct {
    rdb *redis.Client
    ttl time.Duration
}

func (s *SessionStore) Set(ctx context.Context, sessionID string, data map[string]interface{}) error {
    json, _ := json.Marshal(data)
    return s.rdb.Set(ctx, "session:"+sessionID, json, s.ttl).Err()
}

func (s *SessionStore) Get(ctx context.Context, sessionID string) (map[string]interface{}, error) {
    data, err := s.rdb.Get(ctx, "session:"+sessionID).Result()
    if err == redis.Nil {
        return nil, ErrSessionNotFound
    }
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    json.Unmarshal([]byte(data), &result)
    return result, nil
}

func (s *SessionStore) Refresh(ctx context.Context, sessionID string) error {
    return s.rdb.Expire(ctx, "session:"+sessionID, s.ttl).Err()
}
```

---

### 9.6 Database Best Practices

#### ⚫ Q73: What are production database best practices in Go?

**Answer:**

**1. Always use context with timeout:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

row := db.QueryRowContext(ctx, "SELECT ...")
```

**2. Prepared statements for repeated queries:**
```go
type UserRepo struct {
    db         *sql.DB
    stmtGetByID *sql.Stmt
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {
    stmt, err := db.Prepare("SELECT id, name, email FROM users WHERE id = $1")
    if err != nil {
        return nil, err
    }
    return &UserRepo{db: db, stmtGetByID: stmt}, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
    var u User
    err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(&u.ID, &u.Name, &u.Email)
    return &u, err
}

func (r *UserRepo) Close() {
    r.stmtGetByID.Close()
}
```

**3. Handle NULL values:**
```go
type User struct {
    ID        string
    Name      string
    DeletedAt sql.NullTime  // Nullable column
}

// Or use pointers
type User struct {
    ID        string
    Name      string
    DeletedAt *time.Time
}
```

**4. Scan into struct with sqlx:**
```go
import "github.com/jmoiron/sqlx"

type User struct {
    ID    string `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

func (r *UserRepo) GetAll(ctx context.Context) ([]User, error) {
    var users []User
    err := r.db.SelectContext(ctx, &users, "SELECT id, name, email FROM users")
    return users, err
}
```

**5. Health checks:**
```go
func (r *Repo) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    return r.db.PingContext(ctx)
}
```

**6. Migrations:**
```go
import "github.com/golang-migrate/migrate/v4"

func RunMigrations(dbURL string) error {
    m, err := migrate.New("file://migrations", dbURL)
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}
```

---

## Part 2 Summary Checklist

### Error Handling
- [ ] Error wrapping with `%w`, `errors.Is`, `errors.As`
- [ ] Custom error types (sentinel, typed with context)
- [ ] `Unwrap()` method for error chains
- [ ] Panic vs error usage (unrecoverable vs expected)
- [ ] `recover()` in middleware/boundaries
- [ ] Production error patterns (error kinds, HTTP status mapping)
- [ ] Structured logging (`log/slog`)

### Go Modules
- [ ] `go.mod` and `go.sum` structure
- [ ] Semantic versioning (v2+ import path)
- [ ] Private modules (GOPRIVATE, authentication)
- [ ] Dependency management (`go mod tidy`, `replace`, `exclude`)

### Testing
- [ ] Table-driven tests
- [ ] Test fixtures and helpers (`t.Helper()`, testdata/)
- [ ] Interface-based mocking
- [ ] HTTP test server (`httptest.NewServer`)
- [ ] Benchmarks (`b.N`, `b.ResetTimer()`, `b.ReportAllocs()`)
- [ ] Race detector (`go test -race`)
- [ ] Coverage (`go test -cover`, `-coverprofile`)

### Production APIs
- [ ] Project structure (cmd, internal, pkg)
- [ ] Handler patterns (dependency injection)
- [ ] Middleware (logging, auth, rate limiting, recovery)
- [ ] Middleware chaining
- [ ] Graceful shutdown (`server.Shutdown()`, signal handling)
- [ ] Configuration (YAML, env vars, validation)
- [ ] Dependency injection (manual, wire, fx)

### Database & Caching
- [ ] `database/sql` patterns (Query, QueryRow, Exec)
- [ ] Connection pooling configuration
- [ ] Transaction handling (`BeginTx`, `Commit`, `Rollback`)
- [ ] N+1 query prevention (JOIN, batch loading)
- [ ] Redis patterns (cache-aside, locks, rate limiting, pub/sub)
- [ ] Production best practices (context, prepared statements, health checks)

---

> **Continue to Part 3:** Performance & Profiling, Distributed Systems, System Design Questions, Rapid Revision
