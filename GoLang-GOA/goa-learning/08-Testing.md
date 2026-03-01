# Part 8: Testing

> **Goal:** Master testing in Goa - from Go's built-in testing package to advanced BDD with Ginkgo/Gomega, including unit tests, integration tests, mocking, and using generated clients

---

## 📚 Table of Contents

1. [Testing Overview](#testing-overview)
2. [Go Built-in Testing (Must Know First)](#go-built-in-testing)
3. [Unit Testing Service Logic](#unit-testing-service-logic)
4. [Testing Generated Endpoints](#testing-generated-endpoints)
5. [Mocking Services](#mocking-services)
6. [Integration Testing](#integration-testing)
7. [Using Generated Client for Testing](#using-generated-client-for-testing)
8. [Ginkgo + Gomega (BDD Testing)](#ginkgo--gomega-bdd-testing)
9. [Testing Best Practices](#testing-best-practices)
10. [Complete Examples](#complete-examples)
11. [Summary](#summary)
12. [Knowledge Check](#knowledge-check)

---

## 🎯 Testing Overview

### Why Testing Matters

```
┌─────────────────────────────────────────────────────────────────┐
│                    TESTING PYRAMID                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                        ▲                                        │
│                       /│\      E2E Tests                        │
│                      / │ \     (Few, Slow, Expensive)           │
│                     /  │  \                                     │
│                    /───┼───\                                    │
│                   /    │    \   Integration Tests               │
│                  /     │     \  (Some, Medium Speed)            │
│                 /──────┼──────\                                 │
│                /       │       \                                │
│               /        │        \  Unit Tests                   │
│              /         │         \ (Many, Fast, Cheap)          │
│             ───────────┴───────────                             │
│                                                                 │
│  GOA TESTING STRATEGY:                                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  1. Unit Tests      → Service logic, helpers, utils     │   │
│  │  2. Endpoint Tests  → Generated endpoints, handlers     │   │
│  │  3. Integration     → HTTP/gRPC servers, database       │   │
│  │  4. E2E Tests       → Full system with generated client │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Testing Tools Comparison

| Tool | Type | Use Case | Learning Curve |
|------|------|----------|----------------|
| `testing` | Standard Library | All Go tests | Low |
| `testify` | Assertions/Mocks | Better assertions | Low |
| `gomock` | Mock Generation | Interface mocking | Medium |
| `httptest` | HTTP Testing | HTTP handlers | Low |
| `Ginkgo` | BDD Framework | Behavior tests | Medium |
| `Gomega` | Matcher Library | With Ginkgo | Medium |

---

## 🧪 Go Built-in Testing (Must Know First)

### The `testing` Package - Standard Library

Go has **built-in testing support** — no external framework required!

```
┌─────────────────────────────────────────────────────────────────┐
│                   GO TESTING BASICS                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  File Naming Convention:                                        │
│  ───────────────────────                                        │
│  • Source file:  calculator.go                                  │
│  • Test file:    calculator_test.go     (must end with _test)   │
│                                                                 │
│  Test Function Naming:                                          │
│  ─────────────────────                                          │
│  • Test:      func TestXxx(t *testing.T)                        │
│  • Benchmark: func BenchmarkXxx(b *testing.B)                   │
│  • Example:   func ExampleXxx()                                 │
│                                                                 │
│  Running Tests:                                                 │
│  ──────────────                                                 │
│  • go test              (run tests in current package)          │
│  • go test ./...        (run tests in all packages)             │
│  • go test -v           (verbose output)                        │
│  • go test -run TestAdd (run specific test)                     │
│  • go test -cover       (show coverage)                         │
│  • go test -bench .     (run benchmarks)                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Basic Test Structure

```go
// calculator.go
package calculator

// Add adds two integers
func Add(a, b int) int {
    return a + b
}

// Divide divides a by b, returns error if b is zero
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```go
// calculator_test.go
package calculator

import (
    "testing"
)

// Basic test function
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

// Test with subtests
func TestAddSubtests(t *testing.T) {
    t.Run("positive numbers", func(t *testing.T) {
        result := Add(2, 3)
        if result != 5 {
            t.Errorf("Add(2, 3) = %d; want 5", result)
        }
    })
    
    t.Run("negative numbers", func(t *testing.T) {
        result := Add(-2, -3)
        if result != -5 {
            t.Errorf("Add(-2, -3) = %d; want -5", result)
        }
    })
    
    t.Run("mixed numbers", func(t *testing.T) {
        result := Add(-2, 5)
        if result != 3 {
            t.Errorf("Add(-2, 5) = %d; want 3", result)
        }
    })
}
```

### Table-Driven Tests (Idiomatic Go)

Table-driven tests are the **preferred pattern** in Go — clean, readable, and easy to extend.

```go
// calculator_test.go
package calculator

import "testing"

func TestAdd_TableDriven(t *testing.T) {
    // Define test cases as a slice of structs
    tests := []struct {
        name     string  // Test case name
        a        int     // First input
        b        int     // Second input
        expected int     // Expected output
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed numbers", -2, 5, 3},
        {"zeros", 0, 0, 0},
        {"large numbers", 1000000, 2000000, 3000000},
    }
    
    // Run each test case
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

func TestDivide_TableDriven(t *testing.T) {
    tests := []struct {
        name        string
        a           int
        b           int
        expected    int
        expectError bool
    }{
        {"normal division", 10, 2, 5, false},
        {"division by zero", 10, 0, 0, true},
        {"negative division", -10, 2, -5, false},
        {"integer division", 7, 3, 2, false},  // 7/3 = 2 (integer)
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)
            
            // Check error expectation
            if tt.expectError {
                if err == nil {
                    t.Errorf("Divide(%d, %d) expected error, got nil", 
                        tt.a, tt.b)
                }
                return // Don't check result if error expected
            }
            
            // No error expected
            if err != nil {
                t.Errorf("Divide(%d, %d) unexpected error: %v", 
                    tt.a, tt.b, err)
                return
            }
            
            if result != tt.expected {
                t.Errorf("Divide(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Test Helper Functions

```go
// calculator_test.go
package calculator

import "testing"

// Helper function to compare results
// t.Helper() marks this function as a test helper
func assertEqual(t *testing.T, got, want int) {
    t.Helper() // Report errors from caller's line
    if got != want {
        t.Errorf("got %d; want %d", got, want)
    }
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func assertError(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Error("expected error, got nil")
    }
}

// Using helpers
func TestAddWithHelpers(t *testing.T) {
    assertEqual(t, Add(2, 3), 5)
    assertEqual(t, Add(-2, -3), -5)
    assertEqual(t, Add(0, 0), 0)
}

func TestDivideWithHelpers(t *testing.T) {
    result, err := Divide(10, 2)
    assertNoError(t, err)
    assertEqual(t, result, 5)
    
    _, err = Divide(10, 0)
    assertError(t, err)
}
```

### Setup and Teardown

```go
// database_test.go
package database

import (
    "testing"
    "os"
)

// TestMain runs before/after all tests in the package
func TestMain(m *testing.M) {
    // Setup - runs once before all tests
    setup()
    
    // Run all tests
    code := m.Run()
    
    // Teardown - runs once after all tests
    teardown()
    
    os.Exit(code)
}

func setup() {
    // Initialize database connection
    // Create test tables
    // Seed test data
}

func teardown() {
    // Clean up test data
    // Close connections
}

// Per-test setup/teardown
func TestWithSetup(t *testing.T) {
    // Setup for this specific test
    db := setupTestDB(t)
    defer db.Close() // Teardown after test
    
    // Your test code
}

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    return db
}
```

### Benchmarks

```go
// calculator_test.go
package calculator

import "testing"

// Benchmark function - starts with Benchmark
func BenchmarkAdd(b *testing.B) {
    // b.N is set by the testing framework
    // It runs the loop enough times to get accurate measurement
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}

// Benchmark with different inputs
func BenchmarkAdd_Large(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1000000000, 2000000000)
    }
}

// Benchmark with setup
func BenchmarkDivide(b *testing.B) {
    // Setup code - not measured
    a := 1000000
    divisor := 7
    
    // Reset timer to exclude setup time
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        Divide(a, divisor)
    }
}

// Benchmark with parallel execution
func BenchmarkAdd_Parallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Add(100, 200)
        }
    })
}
```

Run benchmarks:
```bash
go test -bench .                    # Run all benchmarks
go test -bench BenchmarkAdd        # Run specific benchmark
go test -bench . -benchmem         # Include memory allocations
go test -bench . -benchtime 5s     # Run for 5 seconds
```

### Examples (Documentation Tests)

```go
// calculator_test.go
package calculator

import "fmt"

// Example functions appear in godoc
// Output comment is verified by go test
func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

func ExampleDivide() {
    result, _ := Divide(10, 2)
    fmt.Println(result)
    // Output: 5
}

// Example for a type method
func ExampleCalculator_Add() {
    calc := NewCalculator()
    result := calc.Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

// Unordered output (for maps, etc.)
func ExampleGetKeys() {
    keys := GetKeys(map[string]int{"a": 1, "b": 2})
    for _, k := range keys {
        fmt.Println(k)
    }
    // Unordered output:
    // a
    // b
}
```

### Test Coverage

```bash
# Generate coverage report
go test -cover

# Generate detailed coverage file
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out
```

### Using `t.Parallel()` for Concurrent Tests

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"test1", 1, 2},
        {"test2", 2, 4},
        {"test3", 3, 6},
    }
    
    for _, tt := range tests {
        tt := tt // Capture range variable (important!)
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Run this subtest in parallel
            
            result := Double(tt.input)
            if result != tt.expected {
                t.Errorf("got %d; want %d", result, tt.expected)
            }
        })
    }
}
```

### Testing HTTP Handlers with `httptest`

```go
// handler.go
package api

import (
    "encoding/json"
    "net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "id required", http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "id":   id,
        "name": "Test User",
    })
}
```

```go
// handler_test.go
package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthHandler(t *testing.T) {
    // Create a request
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    
    // Create a ResponseRecorder to capture the response
    rec := httptest.NewRecorder()
    
    // Call the handler
    HealthHandler(rec, req)
    
    // Check status code
    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
    }
    
    // Check response body
    var response map[string]string
    json.NewDecoder(rec.Body).Decode(&response)
    
    if response["status"] != "ok" {
        t.Errorf("status = %s; want ok", response["status"])
    }
}

func TestGetUserHandler(t *testing.T) {
    tests := []struct {
        name       string
        queryID    string
        wantStatus int
        wantBody   map[string]string
    }{
        {
            name:       "valid request",
            queryID:    "123",
            wantStatus: http.StatusOK,
            wantBody:   map[string]string{"id": "123", "name": "Test User"},
        },
        {
            name:       "missing id",
            queryID:    "",
            wantStatus: http.StatusBadRequest,
            wantBody:   nil,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            url := "/user"
            if tt.queryID != "" {
                url += "?id=" + tt.queryID
            }
            
            req := httptest.NewRequest(http.MethodGet, url, nil)
            rec := httptest.NewRecorder()
            
            GetUserHandler(rec, req)
            
            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
            
            if tt.wantBody != nil {
                var response map[string]string
                json.NewDecoder(rec.Body).Decode(&response)
                
                for k, v := range tt.wantBody {
                    if response[k] != v {
                        t.Errorf("response[%s] = %s; want %s", k, response[k], v)
                    }
                }
            }
        })
    }
}

// Testing with a test server
func TestWithServer(t *testing.T) {
    // Create a test server
    server := httptest.NewServer(http.HandlerFunc(HealthHandler))
    defer server.Close()
    
    // Make a real HTTP request
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
    }
}
```

### Summary: Built-in Testing

```
┌─────────────────────────────────────────────────────────────────┐
│              GO TESTING QUICK REFERENCE                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  File:     *_test.go                                            │
│  Test:     func TestXxx(t *testing.T)                           │
│  Bench:    func BenchmarkXxx(b *testing.B)                      │
│  Example:  func ExampleXxx()                                    │
│                                                                 │
│  Commands:                                                      │
│  ─────────                                                      │
│  go test              Run tests                                 │
│  go test -v           Verbose                                   │
│  go test -run Name    Run specific test                         │
│  go test -cover       Show coverage                             │
│  go test -bench .     Run benchmarks                            │
│  go test ./...        All packages                              │
│                                                                 │
│  Key Functions:                                                 │
│  ──────────────                                                 │
│  t.Error(msg)         Log error, continue                       │
│  t.Errorf(fmt, ...)   Log formatted error                       │
│  t.Fatal(msg)         Log error, stop test                      │
│  t.Fatalf(fmt, ...)   Log formatted, stop                       │
│  t.Skip(msg)          Skip this test                            │
│  t.Helper()           Mark as helper function                   │
│  t.Parallel()         Run in parallel                           │
│  t.Run(name, fn)      Run subtest                               │
│                                                                 │
│  ✔ Recommended for all Go developers                            │
│  ✔ Fast and simple                                              │
│  ✔ Used in production everywhere                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔬 Unit Testing Service Logic

### Testing Goa Service Methods

```go
// service/users.go
package service

import (
    "context"
    "errors"
    "time"
    
    users "myproject/gen/users"
)

// UsersService implements the users service interface
type UsersService struct {
    repo UserRepository
}

// UserRepository interface for data access
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, offset, limit int) ([]*User, int64, error)
}

// User domain model
type User struct {
    ID        int64
    Name      string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewUsersService creates a new users service
func NewUsersService(repo UserRepository) *UsersService {
    return &UsersService{repo: repo}
}

// Get retrieves a user by ID
func (s *UsersService) Get(ctx context.Context, p *users.GetPayload) (*users.User, error) {
    user, err := s.repo.FindByID(ctx, p.ID)
    if err != nil {
        return nil, users.MakeNotFound(errors.New("user not found"))
    }
    
    return toUserResult(user), nil
}

// Create creates a new user
func (s *UsersService) Create(ctx context.Context, p *users.CreatePayload) (*users.User, error) {
    // Check if email already exists
    existing, _ := s.repo.FindByEmail(ctx, p.Email)
    if existing != nil {
        return nil, users.MakeConflict(errors.New("email already exists"))
    }
    
    // Create user
    user := &User{
        Name:      p.Name,
        Email:     p.Email,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return toUserResult(user), nil
}

// Helper to convert domain model to result type
func toUserResult(u *User) *users.User {
    createdAt := u.CreatedAt.Format(time.RFC3339)
    updatedAt := u.UpdatedAt.Format(time.RFC3339)
    
    return &users.User{
        ID:        &u.ID,
        Name:      &u.Name,
        Email:     &u.Email,
        CreatedAt: &createdAt,
        UpdatedAt: &updatedAt,
    }
}
```

### Unit Tests for Service

```go
// service/users_test.go
package service

import (
    "context"
    "errors"
    "testing"
    "time"
    
    users "myproject/gen/users"
)

// Mock repository for testing
type mockUserRepo struct {
    users    map[int64]*User
    byEmail  map[string]*User
    nextID   int64
    findErr  error
    createErr error
}

func newMockUserRepo() *mockUserRepo {
    return &mockUserRepo{
        users:   make(map[int64]*User),
        byEmail: make(map[string]*User),
        nextID:  1,
    }
}

func (m *mockUserRepo) FindByID(ctx context.Context, id int64) (*User, error) {
    if m.findErr != nil {
        return nil, m.findErr
    }
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
    if m.findErr != nil {
        return nil, m.findErr
    }
    user, ok := m.byEmail[email]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *User) error {
    if m.createErr != nil {
        return m.createErr
    }
    user.ID = m.nextID
    m.nextID++
    m.users[user.ID] = user
    m.byEmail[user.Email] = user
    return nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *User) error {
    m.users[user.ID] = user
    m.byEmail[user.Email] = user
    return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id int64) error {
    delete(m.users, id)
    return nil
}

func (m *mockUserRepo) List(ctx context.Context, offset, limit int) ([]*User, int64, error) {
    var result []*User
    for _, u := range m.users {
        result = append(result, u)
    }
    return result, int64(len(result)), nil
}

// Helper to add test user to mock repo
func (m *mockUserRepo) addUser(user *User) {
    m.users[user.ID] = user
    m.byEmail[user.Email] = user
    if user.ID >= m.nextID {
        m.nextID = user.ID + 1
    }
}

// ===========================================
// TESTS
// ===========================================

func TestUsersService_Get(t *testing.T) {
    tests := []struct {
        name      string
        setupRepo func(*mockUserRepo)
        payload   *users.GetPayload
        wantErr   bool
        wantName  string
    }{
        {
            name: "user exists",
            setupRepo: func(repo *mockUserRepo) {
                repo.addUser(&User{
                    ID:        1,
                    Name:      "John Doe",
                    Email:     "john@example.com",
                    CreatedAt: time.Now(),
                    UpdatedAt: time.Now(),
                })
            },
            payload:  &users.GetPayload{ID: 1},
            wantErr:  false,
            wantName: "John Doe",
        },
        {
            name: "user not found",
            setupRepo: func(repo *mockUserRepo) {
                // Empty repo
            },
            payload: &users.GetPayload{ID: 999},
            wantErr: true,
        },
        {
            name: "repository error",
            setupRepo: func(repo *mockUserRepo) {
                repo.findErr = errors.New("database error")
            },
            payload: &users.GetPayload{ID: 1},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            repo := newMockUserRepo()
            tt.setupRepo(repo)
            svc := NewUsersService(repo)
            
            // Execute
            result, err := svc.Get(context.Background(), tt.payload)
            
            // Verify
            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            
            if *result.Name != tt.wantName {
                t.Errorf("name = %s; want %s", *result.Name, tt.wantName)
            }
        })
    }
}

func TestUsersService_Create(t *testing.T) {
    tests := []struct {
        name      string
        setupRepo func(*mockUserRepo)
        payload   *users.CreatePayload
        wantErr   bool
        errType   string
    }{
        {
            name: "successful creation",
            setupRepo: func(repo *mockUserRepo) {
                // Empty repo
            },
            payload: &users.CreatePayload{
                Name:  "Jane Doe",
                Email: "jane@example.com",
            },
            wantErr: false,
        },
        {
            name: "email already exists",
            setupRepo: func(repo *mockUserRepo) {
                repo.addUser(&User{
                    ID:    1,
                    Name:  "Existing",
                    Email: "jane@example.com",
                })
            },
            payload: &users.CreatePayload{
                Name:  "Jane Doe",
                Email: "jane@example.com",
            },
            wantErr: true,
            errType: "conflict",
        },
        {
            name: "repository create error",
            setupRepo: func(repo *mockUserRepo) {
                repo.createErr = errors.New("database error")
            },
            payload: &users.CreatePayload{
                Name:  "New User",
                Email: "new@example.com",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            repo := newMockUserRepo()
            tt.setupRepo(repo)
            svc := NewUsersService(repo)
            
            // Execute
            result, err := svc.Create(context.Background(), tt.payload)
            
            // Verify
            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            
            // Verify result
            if *result.Name != tt.payload.Name {
                t.Errorf("name = %s; want %s", *result.Name, tt.payload.Name)
            }
            
            if *result.Email != tt.payload.Email {
                t.Errorf("email = %s; want %s", *result.Email, tt.payload.Email)
            }
            
            // Verify ID was assigned
            if result.ID == nil || *result.ID == 0 {
                t.Error("expected ID to be assigned")
            }
        })
    }
}
```

---

## 🎯 Testing Generated Endpoints

### Understanding Goa Endpoints

```
┌─────────────────────────────────────────────────────────────────┐
│                GOA ENDPOINT ARCHITECTURE                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  HTTP Request                                                   │
│       │                                                         │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  HTTP Handler (generated)                                │   │
│  │  • Decodes HTTP request                                  │   │
│  │  • Validates payload                                     │   │
│  │  • Calls endpoint                                        │   │
│  └─────────────────────────────────────────────────────────┘   │
│       │                                                         │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Endpoint (generated)                                    │   │
│  │  • Wraps service method                                  │   │
│  │  • Type: func(ctx, request) (response, error)           │   │
│  └─────────────────────────────────────────────────────────┘   │
│       │                                                         │
│       ▼                                                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Service Method (your code)                              │   │
│  │  • Business logic                                        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Testing Endpoints Directly

```go
// endpoint_test.go
package service

import (
    "context"
    "testing"
    
    users "myproject/gen/users"
)

func TestEndpoint_Get(t *testing.T) {
    // Setup
    repo := newMockUserRepo()
    repo.addUser(&User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    })
    
    svc := NewUsersService(repo)
    
    // Create endpoints from service
    endpoints := users.NewEndpoints(svc)
    
    // Test the endpoint directly
    t.Run("get existing user", func(t *testing.T) {
        // Create payload
        payload := &users.GetPayload{ID: 1}
        
        // Call endpoint
        result, err := endpoints.Get(context.Background(), payload)
        
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        
        user, ok := result.(*users.User)
        if !ok {
            t.Fatalf("unexpected result type: %T", result)
        }
        
        if *user.Name != "Test User" {
            t.Errorf("name = %s; want Test User", *user.Name)
        }
    })
    
    t.Run("get non-existing user", func(t *testing.T) {
        payload := &users.GetPayload{ID: 999}
        
        _, err := endpoints.Get(context.Background(), payload)
        
        if err == nil {
            t.Error("expected error, got nil")
        }
    })
}
```

### Testing HTTP Handlers

```go
// http_test.go
package service

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    goahttp "goa.design/goa/v3/http"
    
    users "myproject/gen/users"
    userssvr "myproject/gen/http/users/server"
)

func TestHTTPHandler_Get(t *testing.T) {
    // Setup
    repo := newMockUserRepo()
    repo.addUser(&User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    })
    
    svc := NewUsersService(repo)
    endpoints := users.NewEndpoints(svc)
    
    // Create HTTP mux
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
    
    // Test cases
    tests := []struct {
        name       string
        path       string
        wantStatus int
        wantBody   map[string]interface{}
    }{
        {
            name:       "get existing user",
            path:       "/users/1",
            wantStatus: http.StatusOK,
            wantBody: map[string]interface{}{
                "name":  "Test User",
                "email": "test@example.com",
            },
        },
        {
            name:       "get non-existing user",
            path:       "/users/999",
            wantStatus: http.StatusNotFound,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tt.path, nil)
            rec := httptest.NewRecorder()
            
            mux.ServeHTTP(rec, req)
            
            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
            
            if tt.wantBody != nil {
                var response map[string]interface{}
                json.NewDecoder(rec.Body).Decode(&response)
                
                for k, v := range tt.wantBody {
                    if response[k] != v {
                        t.Errorf("response[%s] = %v; want %v", k, response[k], v)
                    }
                }
            }
        })
    }
}

func TestHTTPHandler_Create(t *testing.T) {
    repo := newMockUserRepo()
    svc := NewUsersService(repo)
    endpoints := users.NewEndpoints(svc)
    
    mux := goahttp.NewMuxer()
    server := userssvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        nil,
        nil,
    )
    userssvr.Mount(mux, server)
    
    tests := []struct {
        name       string
        body       map[string]string
        wantStatus int
    }{
        {
            name: "valid creation",
            body: map[string]string{
                "name":  "New User",
                "email": "new@example.com",
            },
            wantStatus: http.StatusCreated,
        },
        {
            name: "missing required field",
            body: map[string]string{
                "name": "New User",
                // missing email
            },
            wantStatus: http.StatusBadRequest,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            bodyBytes, _ := json.Marshal(tt.body)
            req := httptest.NewRequest(
                http.MethodPost,
                "/users",
                bytes.NewReader(bodyBytes),
            )
            req.Header.Set("Content-Type", "application/json")
            
            rec := httptest.NewRecorder()
            mux.ServeHTTP(rec, req)
            
            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d\nbody: %s", 
                    rec.Code, tt.wantStatus, rec.Body.String())
            }
        })
    }
}
```

---

## 🎭 Mocking Services

### Using `testify/mock`

```bash
go get github.com/stretchr/testify
```

```go
// mocks/repository_mock.go
package mocks

import (
    "context"
    
    "github.com/stretchr/testify/mock"
    
    "myproject/service"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*service.User, error) {
    args := m.Called(ctx, id)
    
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*service.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*service.User, error) {
    args := m.Called(ctx, email)
    
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*service.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *service.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *service.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*service.User, int64, error) {
    args := m.Called(ctx, offset, limit)
    return args.Get(0).([]*service.User), args.Get(1).(int64), args.Error(2)
}
```

### Using Testify Mock in Tests

```go
// service/users_testify_test.go
package service

import (
    "context"
    "errors"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    
    "myproject/mocks"
    users "myproject/gen/users"
)

func TestUsersService_Get_WithTestify(t *testing.T) {
    // Create mock
    mockRepo := new(mocks.MockUserRepository)
    
    // Setup expectations
    testUser := &User{
        ID:        1,
        Name:      "John Doe",
        Email:     "john@example.com",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    mockRepo.On("FindByID", mock.Anything, int64(1)).Return(testUser, nil)
    mockRepo.On("FindByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
    
    // Create service with mock
    svc := NewUsersService(mockRepo)
    
    t.Run("user exists", func(t *testing.T) {
        result, err := svc.Get(context.Background(), &users.GetPayload{ID: 1})
        
        require.NoError(t, err)
        assert.Equal(t, "John Doe", *result.Name)
        assert.Equal(t, "john@example.com", *result.Email)
    })
    
    t.Run("user not found", func(t *testing.T) {
        _, err := svc.Get(context.Background(), &users.GetPayload{ID: 999})
        
        assert.Error(t, err)
    })
    
    // Verify all expectations were met
    mockRepo.AssertExpectations(t)
}

func TestUsersService_Create_WithTestify(t *testing.T) {
    t.Run("successful creation", func(t *testing.T) {
        mockRepo := new(mocks.MockUserRepository)
        
        // Email doesn't exist
        mockRepo.On("FindByEmail", mock.Anything, "new@example.com").
            Return(nil, errors.New("not found"))
        
        // Create succeeds
        mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*service.User")).
            Return(nil).
            Run(func(args mock.Arguments) {
                // Simulate ID assignment
                user := args.Get(1).(*User)
                user.ID = 1
            })
        
        svc := NewUsersService(mockRepo)
        
        result, err := svc.Create(context.Background(), &users.CreatePayload{
            Name:  "New User",
            Email: "new@example.com",
        })
        
        require.NoError(t, err)
        assert.Equal(t, "New User", *result.Name)
        assert.NotNil(t, result.ID)
        
        mockRepo.AssertExpectations(t)
    })
    
    t.Run("email already exists", func(t *testing.T) {
        mockRepo := new(mocks.MockUserRepository)
        
        existingUser := &User{ID: 1, Email: "existing@example.com"}
        mockRepo.On("FindByEmail", mock.Anything, "existing@example.com").
            Return(existingUser, nil)
        
        svc := NewUsersService(mockRepo)
        
        _, err := svc.Create(context.Background(), &users.CreatePayload{
            Name:  "New User",
            Email: "existing@example.com",
        })
        
        assert.Error(t, err)
        mockRepo.AssertExpectations(t)
    })
}
```

### Using gomock

```bash
go install github.com/golang/mock/mockgen@latest
```

Generate mocks:
```bash
mockgen -source=service/repository.go -destination=mocks/mock_repository.go -package=mocks
```

```go
// Using generated gomock mocks
package service

import (
    "context"
    "errors"
    "testing"
    
    "github.com/golang/mock/gomock"
    
    "myproject/mocks"
    users "myproject/gen/users"
)

func TestUsersService_Get_Gomock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    
    testUser := &User{
        ID:    1,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Setup expectations with gomock
    mockRepo.EXPECT().
        FindByID(gomock.Any(), int64(1)).
        Return(testUser, nil).
        Times(1)
    
    svc := NewUsersService(mockRepo)
    
    result, err := svc.Get(context.Background(), &users.GetPayload{ID: 1})
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if *result.Name != "John Doe" {
        t.Errorf("name = %s; want John Doe", *result.Name)
    }
}
```

---

## 🔗 Integration Testing

### Setting Up Integration Tests

```go
// integration/setup_test.go
package integration

import (
    "context"
    "database/sql"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
    
    _ "github.com/lib/pq"
    goahttp "goa.design/goa/v3/http"
    
    users "myproject/gen/users"
    userssvr "myproject/gen/http/users/server"
    "myproject/repository"
    "myproject/service"
)

var (
    testDB     *sql.DB
    testServer *httptest.Server
    testClient *http.Client
)

func TestMain(m *testing.M) {
    // Setup
    var err error
    
    // Connect to test database
    testDB, err = sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        panic(err)
    }
    
    // Run migrations
    if err := runMigrations(testDB); err != nil {
        panic(err)
    }
    
    // Setup test server
    repo := repository.NewPostgresUserRepository(testDB)
    svc := service.NewUsersService(repo)
    endpoints := users.NewEndpoints(svc)
    
    mux := goahttp.NewMuxer()
    server := userssvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        nil,
        nil,
    )
    userssvr.Mount(mux, server)
    
    testServer = httptest.NewServer(mux)
    testClient = testServer.Client()
    
    // Run tests
    code := m.Run()
    
    // Teardown
    testServer.Close()
    testDB.Close()
    
    os.Exit(code)
}

func runMigrations(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    return err
}

// Helper to clean database between tests
func cleanDB(t *testing.T) {
    t.Helper()
    _, err := testDB.Exec("DELETE FROM users")
    if err != nil {
        t.Fatalf("failed to clean db: %v", err)
    }
}
```

### Integration Test Examples

```go
// integration/users_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"
)

func TestUsersCRUD_Integration(t *testing.T) {
    cleanDB(t)
    
    var userID int64
    
    t.Run("create user", func(t *testing.T) {
        body := map[string]string{
            "name":  "Integration Test User",
            "email": "integration@test.com",
        }
        bodyBytes, _ := json.Marshal(body)
        
        resp, err := testClient.Post(
            testServer.URL+"/users",
            "application/json",
            bytes.NewReader(bodyBytes),
        )
        if err != nil {
            t.Fatalf("request failed: %v", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusCreated {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusCreated)
        }
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        
        userID = int64(result["id"].(float64))
        if userID == 0 {
            t.Error("expected user ID to be assigned")
        }
    })
    
    t.Run("get user", func(t *testing.T) {
        resp, err := testClient.Get(
            fmt.Sprintf("%s/users/%d", testServer.URL, userID),
        )
        if err != nil {
            t.Fatalf("request failed: %v", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
        }
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        
        if result["name"] != "Integration Test User" {
            t.Errorf("name = %v; want Integration Test User", result["name"])
        }
    })
    
    t.Run("update user", func(t *testing.T) {
        body := map[string]string{
            "name": "Updated Name",
        }
        bodyBytes, _ := json.Marshal(body)
        
        req, _ := http.NewRequest(
            http.MethodPut,
            fmt.Sprintf("%s/users/%d", testServer.URL, userID),
            bytes.NewReader(bodyBytes),
        )
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := testClient.Do(req)
        if err != nil {
            t.Fatalf("request failed: %v", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
        }
    })
    
    t.Run("delete user", func(t *testing.T) {
        req, _ := http.NewRequest(
            http.MethodDelete,
            fmt.Sprintf("%s/users/%d", testServer.URL, userID),
            nil,
        )
        
        resp, err := testClient.Do(req)
        if err != nil {
            t.Fatalf("request failed: %v", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusNoContent {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusNoContent)
        }
        
        // Verify deletion
        resp, _ = testClient.Get(
            fmt.Sprintf("%s/users/%d", testServer.URL, userID),
        )
        if resp.StatusCode != http.StatusNotFound {
            t.Errorf("expected 404 after deletion, got %d", resp.StatusCode)
        }
    })
}

func TestUsersList_Integration(t *testing.T) {
    cleanDB(t)
    
    // Create multiple users
    for i := 1; i <= 5; i++ {
        body := map[string]string{
            "name":  fmt.Sprintf("User %d", i),
            "email": fmt.Sprintf("user%d@test.com", i),
        }
        bodyBytes, _ := json.Marshal(body)
        
        resp, _ := testClient.Post(
            testServer.URL+"/users",
            "application/json",
            bytes.NewReader(bodyBytes),
        )
        resp.Body.Close()
    }
    
    t.Run("list all users", func(t *testing.T) {
        resp, err := testClient.Get(testServer.URL + "/users")
        if err != nil {
            t.Fatalf("request failed: %v", err)
        }
        defer resp.Body.Close()
        
        var result struct {
            Users      []map[string]interface{} `json:"users"`
            Total      int                      `json:"total"`
            Page       int                      `json:"page"`
            PerPage    int                      `json:"per_page"`
        }
        json.NewDecoder(resp.Body).Decode(&result)
        
        if result.Total != 5 {
            t.Errorf("total = %d; want 5", result.Total)
        }
    })
    
    t.Run("list with pagination", func(t *testing.T) {
        resp, _ := testClient.Get(testServer.URL + "/users?page=1&per_page=2")
        defer resp.Body.Close()
        
        var result struct {
            Users []map[string]interface{} `json:"users"`
        }
        json.NewDecoder(resp.Body).Decode(&result)
        
        if len(result.Users) != 2 {
            t.Errorf("users count = %d; want 2", len(result.Users))
        }
    })
}
```

### Using Docker for Integration Tests

```go
// integration/docker_test.go
package integration

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "testing"
    "time"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/lib/pq"
)

var (
    pgContainer testcontainers.Container
    dbURL       string
)

func TestMain(m *testing.M) {
    ctx := context.Background()
    
    // Start PostgreSQL container
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections").
            WithOccurrence(2).
            WithStartupTimeout(60 * time.Second),
    }
    
    var err error
    pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        panic(err)
    }
    
    // Get connection details
    host, _ := pgContainer.Host(ctx)
    port, _ := pgContainer.MappedPort(ctx, "5432")
    dbURL = fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", 
        host, port.Port())
    
    // Setup database
    testDB, err = sql.Open("postgres", dbURL)
    if err != nil {
        panic(err)
    }
    
    // Wait for connection
    for i := 0; i < 10; i++ {
        if err := testDB.Ping(); err == nil {
            break
        }
        time.Sleep(time.Second)
    }
    
    // Run migrations
    runMigrations(testDB)
    
    // Setup server...
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    pgContainer.Terminate(ctx)
    
    os.Exit(code)
}
```

---

## 🤖 Using Generated Client for Testing

### Goa Generated Client

Goa generates type-safe HTTP and gRPC clients that you can use for testing.

```go
// client_test.go
package integration

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    
    goahttp "goa.design/goa/v3/http"
    
    users "myproject/gen/users"
    userssvr "myproject/gen/http/users/server"
    userscli "myproject/gen/http/users/client"
)

func TestWithGeneratedClient(t *testing.T) {
    // Setup server
    repo := newMockUserRepo()
    repo.addUser(&User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    })
    
    svc := NewUsersService(repo)
    endpoints := users.NewEndpoints(svc)
    
    mux := goahttp.NewMuxer()
    server := userssvr.New(
        endpoints,
        mux,
        goahttp.RequestDecoder,
        goahttp.ResponseEncoder,
        nil,
        nil,
    )
    userssvr.Mount(mux, server)
    
    // Create test server
    ts := httptest.NewServer(mux)
    defer ts.Close()
    
    // Create generated client
    httpClient := &http.Client{}
    client := userscli.NewClient(
        ts.URL,
        httpClient,
        goahttp.RequestEncoder,
        goahttp.ResponseDecoder,
        false,
    )
    
    // Create endpoints from client
    clientEndpoints := users.Endpoints{
        Get:    userscli.BuildGetEndpoint(client),
        Create: userscli.BuildCreateEndpoint(client),
        List:   userscli.BuildListEndpoint(client),
    }
    
    t.Run("get user via client", func(t *testing.T) {
        result, err := clientEndpoints.Get(context.Background(), &users.GetPayload{ID: 1})
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        
        user := result.(*users.User)
        if *user.Name != "Test User" {
            t.Errorf("name = %s; want Test User", *user.Name)
        }
    })
    
    t.Run("create user via client", func(t *testing.T) {
        result, err := clientEndpoints.Create(context.Background(), &users.CreatePayload{
            Name:  "New User",
            Email: "new@example.com",
        })
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        
        user := result.(*users.User)
        if *user.Name != "New User" {
            t.Errorf("name = %s; want New User", *user.Name)
        }
    })
}
```

### Testing gRPC with Generated Client

```go
// grpc_client_test.go
package integration

import (
    "context"
    "net"
    "testing"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/test/bufconn"
    
    users "myproject/gen/users"
    usersgrpc "myproject/gen/grpc/users/server"
    userspb "myproject/gen/grpc/users/pb"
    usersclient "myproject/gen/grpc/users/client"
)

const bufSize = 1024 * 1024

func TestGRPCClient(t *testing.T) {
    // Setup in-memory connection
    lis := bufconn.Listen(bufSize)
    
    // Setup gRPC server
    repo := newMockUserRepo()
    repo.addUser(&User{ID: 1, Name: "Test User", Email: "test@example.com"})
    
    svc := NewUsersService(repo)
    endpoints := users.NewEndpoints(svc)
    
    grpcServer := grpc.NewServer()
    usersServer := usersgrpc.New(endpoints, nil)
    userspb.RegisterUsersServer(grpcServer, usersServer)
    
    go func() {
        grpcServer.Serve(lis)
    }()
    defer grpcServer.Stop()
    
    // Create client connection
    conn, err := grpc.DialContext(
        context.Background(),
        "bufnet",
        grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
            return lis.Dial()
        }),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        t.Fatalf("failed to dial: %v", err)
    }
    defer conn.Close()
    
    // Create client endpoints
    client := usersclient.NewClient(conn, nil)
    
    t.Run("get user via gRPC client", func(t *testing.T) {
        result, err := client.Get()(context.Background(), &users.GetPayload{ID: 1})
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        
        user := result.(*users.User)
        if *user.Name != "Test User" {
            t.Errorf("name = %s; want Test User", *user.Name)
        }
    })
}
```

---

## 🎭 Ginkgo + Gomega (BDD Testing)

### What is BDD Testing?

**Behavior-Driven Development (BDD)** focuses on describing the behavior of a system from the user's perspective.

```
┌─────────────────────────────────────────────────────────────────┐
│                   BDD vs TDD                                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Traditional TDD:                                               │
│  ─────────────────                                              │
│  func TestAdd(t *testing.T) {                                   │
│      result := Add(2, 3)                                        │
│      if result != 5 {                                           │
│          t.Errorf("Add(2, 3) = %d; want 5", result)             │
│      }                                                          │
│  }                                                              │
│                                                                 │
│  BDD with Ginkgo:                                               │
│  ─────────────────                                              │
│  Describe("Add", func() {                                       │
│      When("adding 2 and 3", func() {                            │
│          It("should return 5", func() {                         │
│              Expect(Add(2, 3)).To(Equal(5))                     │
│          })                                                     │
│      })                                                         │
│  })                                                             │
│                                                                 │
│  Advantages of BDD:                                             │
│  • Human-readable specifications                                │
│  • Self-documenting tests                                       │
│  • Focus on behavior, not implementation                        │
│  • Better for complex scenarios                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Installing Ginkgo and Gomega

```bash
# Install Ginkgo CLI
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Add to your project
go get github.com/onsi/ginkgo/v2
go get github.com/onsi/gomega
```

### Bootstrap Ginkgo Test Suite

```bash
# Generate test suite in your package
cd service
ginkgo bootstrap

# Generate test file for a specific file
ginkgo generate users_service
```

### Basic Ginkgo Structure

```go
// service/service_suite_test.go (generated by ginkgo bootstrap)
package service_test

import (
    "testing"
    
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Service Suite")
}
```

### Ginkgo Test Structure

```go
// service/users_service_test.go
package service_test

import (
    "context"
    "errors"
    "time"
    
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    
    "myproject/service"
    users "myproject/gen/users"
)

var _ = Describe("UsersService", func() {
    var (
        svc      *service.UsersService
        mockRepo *MockUserRepository
        ctx      context.Context
    )
    
    // Setup before each test
    BeforeEach(func() {
        ctx = context.Background()
        mockRepo = NewMockUserRepository()
        svc = service.NewUsersService(mockRepo)
    })
    
    // Cleanup after each test (if needed)
    AfterEach(func() {
        // Cleanup code
    })
    
    // Describe a method
    Describe("Get", func() {
        // Context: specific scenario
        Context("when the user exists", func() {
            BeforeEach(func() {
                mockRepo.AddUser(&service.User{
                    ID:        1,
                    Name:      "John Doe",
                    Email:     "john@example.com",
                    CreatedAt: time.Now(),
                    UpdatedAt: time.Now(),
                })
            })
            
            It("should return the user", func() {
                result, err := svc.Get(ctx, &users.GetPayload{ID: 1})
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeNil())
                Expect(*result.Name).To(Equal("John Doe"))
                Expect(*result.Email).To(Equal("john@example.com"))
            })
        })
        
        Context("when the user does not exist", func() {
            It("should return a not found error", func() {
                result, err := svc.Get(ctx, &users.GetPayload{ID: 999})
                
                Expect(err).To(HaveOccurred())
                Expect(result).To(BeNil())
            })
        })
        
        Context("when the repository returns an error", func() {
            BeforeEach(func() {
                mockRepo.SetFindError(errors.New("database error"))
            })
            
            It("should propagate the error", func() {
                _, err := svc.Get(ctx, &users.GetPayload{ID: 1})
                
                Expect(err).To(HaveOccurred())
            })
        })
    })
    
    Describe("Create", func() {
        Context("when creating a new user", func() {
            It("should create the user successfully", func() {
                payload := &users.CreatePayload{
                    Name:  "Jane Doe",
                    Email: "jane@example.com",
                }
                
                result, err := svc.Create(ctx, payload)
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeNil())
                Expect(*result.Name).To(Equal("Jane Doe"))
                Expect(*result.Email).To(Equal("jane@example.com"))
                Expect(result.ID).NotTo(BeNil())
            })
        })
        
        Context("when the email already exists", func() {
            BeforeEach(func() {
                mockRepo.AddUser(&service.User{
                    ID:    1,
                    Name:  "Existing User",
                    Email: "existing@example.com",
                })
            })
            
            It("should return a conflict error", func() {
                payload := &users.CreatePayload{
                    Name:  "New User",
                    Email: "existing@example.com",
                }
                
                result, err := svc.Create(ctx, payload)
                
                Expect(err).To(HaveOccurred())
                Expect(result).To(BeNil())
            })
        })
    })
    
    Describe("List", func() {
        Context("when there are multiple users", func() {
            BeforeEach(func() {
                for i := 1; i <= 10; i++ {
                    mockRepo.AddUser(&service.User{
                        ID:        int64(i),
                        Name:      fmt.Sprintf("User %d", i),
                        Email:     fmt.Sprintf("user%d@example.com", i),
                        CreatedAt: time.Now(),
                        UpdatedAt: time.Now(),
                    })
                }
            })
            
            It("should return paginated results", func() {
                payload := &users.ListPayload{
                    Page:    1,
                    PerPage: 5,
                }
                
                result, err := svc.List(ctx, payload)
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result.Users).To(HaveLen(5))
                Expect(*result.Total).To(Equal(int64(10)))
            })
            
            It("should return second page", func() {
                payload := &users.ListPayload{
                    Page:    2,
                    PerPage: 5,
                }
                
                result, err := svc.List(ctx, payload)
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result.Users).To(HaveLen(5))
            })
        })
        
        Context("when there are no users", func() {
            It("should return an empty list", func() {
                payload := &users.ListPayload{
                    Page:    1,
                    PerPage: 10,
                }
                
                result, err := svc.List(ctx, payload)
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result.Users).To(BeEmpty())
                Expect(*result.Total).To(Equal(int64(0)))
            })
        })
    })
})
```

### Gomega Matchers - Complete Reference

```go
// gomega_matchers_test.go
package service_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Gomega Matchers Reference", func() {
    
    // ==========================================
    // EQUALITY MATCHERS
    // ==========================================
    
    Describe("Equality Matchers", func() {
        It("Equal - deep equality", func() {
            Expect(5).To(Equal(5))
            Expect("hello").To(Equal("hello"))
            Expect([]int{1, 2}).To(Equal([]int{1, 2}))
        })
        
        It("BeEquivalentTo - type-converting equality", func() {
            Expect(int32(5)).To(BeEquivalentTo(int64(5)))
            Expect(5.0).To(BeEquivalentTo(5))
        })
        
        It("BeIdenticalTo - pointer identity", func() {
            x := "hello"
            Expect(&x).To(BeIdenticalTo(&x))
        })
        
        It("BeNumerically - numeric comparisons", func() {
            Expect(5).To(BeNumerically("==", 5))
            Expect(5).To(BeNumerically("<", 10))
            Expect(5).To(BeNumerically(">", 3))
            Expect(5).To(BeNumerically("<=", 5))
            Expect(5).To(BeNumerically(">=", 5))
            Expect(5).To(BeNumerically("~", 5.5, 0.6)) // approximately
        })
    })
    
    // ==========================================
    // NIL/ZERO MATCHERS
    // ==========================================
    
    Describe("Nil/Zero Matchers", func() {
        It("BeNil - nil check", func() {
            var p *string
            Expect(p).To(BeNil())
            
            str := "hello"
            Expect(&str).NotTo(BeNil())
        })
        
        It("BeZero - zero value", func() {
            Expect(0).To(BeZero())
            Expect("").To(BeZero())
            Expect([]int(nil)).To(BeZero())
            
            var s struct{ Name string }
            Expect(s).To(BeZero())
        })
    })
    
    // ==========================================
    // BOOLEAN MATCHERS
    // ==========================================
    
    Describe("Boolean Matchers", func() {
        It("BeTrue / BeFalse", func() {
            Expect(true).To(BeTrue())
            Expect(false).To(BeFalse())
        })
    })
    
    // ==========================================
    // ERROR MATCHERS
    // ==========================================
    
    Describe("Error Matchers", func() {
        It("HaveOccurred - error exists", func() {
            err := errors.New("something went wrong")
            Expect(err).To(HaveOccurred())
            
            var noErr error
            Expect(noErr).NotTo(HaveOccurred())
        })
        
        It("Succeed - no error", func() {
            Expect(nil).To(Succeed())
        })
        
        It("MatchError - error message matching", func() {
            err := errors.New("not found")
            
            Expect(err).To(MatchError("not found"))
            Expect(err).To(MatchError(ContainSubstring("found")))
        })
        
        It("BeAssignableToTypeOf - error type", func() {
            err := &CustomError{Code: 404}
            Expect(err).To(BeAssignableToTypeOf(&CustomError{}))
        })
    })
    
    // ==========================================
    // STRING MATCHERS
    // ==========================================
    
    Describe("String Matchers", func() {
        It("ContainSubstring", func() {
            Expect("hello world").To(ContainSubstring("world"))
        })
        
        It("HavePrefix / HaveSuffix", func() {
            Expect("hello world").To(HavePrefix("hello"))
            Expect("hello world").To(HaveSuffix("world"))
        })
        
        It("MatchRegexp", func() {
            Expect("hello123").To(MatchRegexp(`hello\d+`))
        })
    })
    
    // ==========================================
    // COLLECTION MATCHERS
    // ==========================================
    
    Describe("Collection Matchers", func() {
        It("BeEmpty", func() {
            Expect([]int{}).To(BeEmpty())
            Expect(map[string]int{}).To(BeEmpty())
            Expect("").To(BeEmpty())
        })
        
        It("HaveLen - exact length", func() {
            Expect([]int{1, 2, 3}).To(HaveLen(3))
            Expect("hello").To(HaveLen(5))
        })
        
        It("HaveCap - capacity", func() {
            s := make([]int, 0, 10)
            Expect(s).To(HaveCap(10))
        })
        
        It("ContainElement - element exists", func() {
            Expect([]int{1, 2, 3}).To(ContainElement(2))
            Expect([]string{"a", "b"}).To(ContainElement("a"))
        })
        
        It("ContainElements - multiple elements (any order)", func() {
            Expect([]int{1, 2, 3, 4}).To(ContainElements(3, 1))
        })
        
        It("ConsistOf - exact elements (any order)", func() {
            Expect([]int{1, 2, 3}).To(ConsistOf(3, 2, 1))
        })
        
        It("HaveEach - all elements match", func() {
            Expect([]int{2, 4, 6}).To(HaveEach(BeNumerically("<=", 6)))
        })
        
        It("HaveKey / HaveKeyWithValue - map matchers", func() {
            m := map[string]int{"a": 1, "b": 2}
            
            Expect(m).To(HaveKey("a"))
            Expect(m).To(HaveKeyWithValue("a", 1))
        })
    })
    
    // ==========================================
    // POINTER MATCHERS
    // ==========================================
    
    Describe("Pointer Matchers", func() {
        It("PointTo - pointer value", func() {
            x := 5
            Expect(&x).To(PointTo(Equal(5)))
            
            str := "hello"
            Expect(&str).To(PointTo(ContainSubstring("ell")))
        })
    })
    
    // ==========================================
    // STRUCT MATCHERS
    // ==========================================
    
    Describe("Struct Matchers", func() {
        It("HaveField - struct field", func() {
            user := User{Name: "John", Email: "john@example.com"}
            
            Expect(user).To(HaveField("Name", "John"))
            Expect(user).To(HaveField("Email", ContainSubstring("@")))
        })
        
        It("MatchFields - multiple fields", func() {
            user := User{Name: "John", Email: "john@example.com"}
            
            Expect(user).To(MatchFields(IgnoreExtras, Fields{
                "Name":  Equal("John"),
                "Email": ContainSubstring("example.com"),
            }))
        })
    })
    
    // ==========================================
    // TYPE MATCHERS
    // ==========================================
    
    Describe("Type Matchers", func() {
        It("BeAssignableToTypeOf", func() {
            var i interface{} = "hello"
            Expect(i).To(BeAssignableToTypeOf(""))
        })
    })
    
    // ==========================================
    // ASYNC MATCHERS
    // ==========================================
    
    Describe("Async Matchers", func() {
        It("Eventually - async assertions", func() {
            counter := 0
            go func() {
                time.Sleep(100 * time.Millisecond)
                counter = 5
            }()
            
            Eventually(func() int {
                return counter
            }).Should(Equal(5))
        })
        
        It("Eventually with timeout", func() {
            ch := make(chan int, 1)
            go func() {
                time.Sleep(50 * time.Millisecond)
                ch <- 42
            }()
            
            Eventually(ch).WithTimeout(1 * time.Second).Should(Receive(Equal(42)))
        })
        
        It("Consistently - always true", func() {
            value := 5
            
            Consistently(func() int {
                return value
            }).Should(Equal(5))
        })
    })
    
    // ==========================================
    // COMBINING MATCHERS
    // ==========================================
    
    Describe("Combining Matchers", func() {
        It("And - all must pass", func() {
            Expect(5).To(And(
                BeNumerically(">", 0),
                BeNumerically("<", 10),
            ))
            
            // Or use SatisfyAll
            Expect(5).To(SatisfyAll(
                BeNumerically(">", 0),
                BeNumerically("<", 10),
            ))
        })
        
        It("Or - at least one must pass", func() {
            Expect(5).To(Or(
                Equal(3),
                Equal(5),
            ))
            
            // Or use SatisfyAny
            Expect(5).To(SatisfyAny(
                Equal(3),
                Equal(5),
            ))
        })
        
        It("Not - negation", func() {
            Expect(5).To(Not(Equal(3)))
            // Same as:
            Expect(5).NotTo(Equal(3))
        })
    })
    
    // ==========================================
    // CHANNEL MATCHERS
    // ==========================================
    
    Describe("Channel Matchers", func() {
        It("Receive - receive from channel", func() {
            ch := make(chan int, 1)
            ch <- 42
            
            Expect(ch).To(Receive(Equal(42)))
        })
        
        It("BeSent - send to channel", func() {
            ch := make(chan int, 1)
            
            Expect(ch).To(BeSent(42))
        })
        
        It("BeClosed", func() {
            ch := make(chan int)
            close(ch)
            
            Expect(ch).To(BeClosed())
        })
    })
})
```

### Ginkgo Lifecycle Hooks

```go
var _ = Describe("Lifecycle Hooks", func() {
    
    // Run once before all specs in this Describe block
    BeforeAll(func() {
        fmt.Println("BeforeAll: Setup expensive resources")
    })
    
    // Run once after all specs in this Describe block
    AfterAll(func() {
        fmt.Println("AfterAll: Cleanup expensive resources")
    })
    
    // Run before each spec
    BeforeEach(func() {
        fmt.Println("BeforeEach: Setup for each test")
    })
    
    // Run after each spec
    AfterEach(func() {
        fmt.Println("AfterEach: Cleanup after each test")
    })
    
    // Cleanup that runs even on failure (like defer)
    JustBeforeEach(func() {
        fmt.Println("JustBeforeEach: Final setup before It")
    })
    
    JustAfterEach(func() {
        fmt.Println("JustAfterEach: First cleanup after It")
    })
    
    // DeferCleanup - like defer but better for tests
    It("uses DeferCleanup", func() {
        resource := acquireResource()
        DeferCleanup(func() {
            resource.Close()
        })
        
        // Use resource...
    })
})
```

### Table-Driven Tests in Ginkgo

```go
var _ = Describe("Calculator with Table Tests", func() {
    DescribeTable("Add function",
        func(a, b, expected int) {
            result := Add(a, b)
            Expect(result).To(Equal(expected))
        },
        Entry("positive numbers", 2, 3, 5),
        Entry("negative numbers", -2, -3, -5),
        Entry("mixed numbers", -2, 5, 3),
        Entry("zeros", 0, 0, 0),
        Entry("large numbers", 1000000, 2000000, 3000000),
    )
    
    DescribeTable("Divide function",
        func(a, b int, expected int, shouldErr bool) {
            result, err := Divide(a, b)
            
            if shouldErr {
                Expect(err).To(HaveOccurred())
            } else {
                Expect(err).NotTo(HaveOccurred())
                Expect(result).To(Equal(expected))
            }
        },
        Entry("normal division", 10, 2, 5, false),
        Entry("division by zero", 10, 0, 0, true),
        Entry("negative division", -10, 2, -5, false),
    )
})
```

### Running Ginkgo Tests

```bash
# Run all tests
ginkgo

# Run with verbose output
ginkgo -v

# Run specific tests by label
ginkgo --label-filter="integration"

# Run tests in parallel
ginkgo -p

# Run with coverage
ginkgo --cover

# Run only focused tests
ginkgo --focus="Create"

# Skip certain tests
ginkgo --skip="slow"

# Watch mode - rerun on file changes
ginkgo watch

# Generate JUnit report
ginkgo --junit-report=report.xml
```

### Focusing and Pending Tests

```go
var _ = Describe("Focus and Pending", func() {
    // FDescribe, FContext, FIt - Focus: only run these
    FIt("focused test - only this will run", func() {
        Expect(true).To(BeTrue())
    })
    
    // PDescribe, PContext, PIt - Pending: skip these
    PIt("pending test - will be skipped", func() {
        Expect(true).To(BeTrue())
    })
    
    // XDescribe, XContext, XIt - Disabled: skip these
    XIt("disabled test", func() {
        Expect(true).To(BeTrue())
    })
    
    // Skip at runtime
    It("conditionally skipped", func() {
        if os.Getenv("CI") == "" {
            Skip("only runs in CI")
        }
        // Test code...
    })
})
```

---

## ✅ Testing Best Practices

### Test Organization

```
┌─────────────────────────────────────────────────────────────────┐
│                TEST ORGANIZATION                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  project/                                                       │
│  ├── service/                                                   │
│  │   ├── users.go              # Production code                │
│  │   ├── users_test.go         # Unit tests (same package)      │
│  │   └── users_suite_test.go   # Ginkgo suite                   │
│  ├── mocks/                                                     │
│  │   └── mock_repository.go    # Generated/manual mocks         │
│  └── integration/                                               │
│      ├── setup_test.go         # Integration test setup         │
│      ├── users_test.go         # Integration tests              │
│      └── docker-compose.yml    # Test infrastructure            │
│                                                                 │
│  Best Practices:                                                │
│  • Unit tests: same package (access unexported)                 │
│  • Integration tests: separate package                          │
│  • Mocks: dedicated package                                     │
│  • Build tags for slow tests                                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Test Naming

```go
// Good test names describe WHAT is being tested and EXPECTED behavior

// For table-driven tests
tests := []struct {
    name string // Descriptive name
}{
    {"returns user when exists"},
    {"returns error when not found"},
    {"returns error when database fails"},
}

// For Ginkgo
Describe("UsersService", func() {
    Describe("Get", func() {
        Context("when user exists", func() {
            It("should return the user with all fields populated", func() {
                // ...
            })
        })
    })
})
```

### Testing Checklist

```
┌─────────────────────────────────────────────────────────────────┐
│                 TESTING CHECKLIST                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Unit Tests:                                                    │
│  ☐ Test happy path                                              │
│  ☐ Test error cases                                             │
│  ☐ Test edge cases (nil, empty, zero)                           │
│  ☐ Test validation errors                                       │
│  ☐ Use table-driven tests for multiple cases                    │
│  ☐ Keep tests fast (<100ms each)                                │
│                                                                 │
│  Integration Tests:                                             │
│  ☐ Test full request/response cycle                             │
│  ☐ Test database operations                                     │
│  ☐ Test external service calls                                  │
│  ☐ Clean database between tests                                 │
│  ☐ Use containers for dependencies                              │
│                                                                 │
│  Coverage:                                                      │
│  ☐ Aim for 80%+ code coverage                                   │
│  ☐ Focus on critical business logic                             │
│  ☐ Don't test generated code                                    │
│                                                                 │
│  General:                                                       │
│  ☐ Tests are deterministic (no random failures)                 │
│  ☐ Tests are independent (can run in any order)                 │
│  ☐ Tests are maintainable (clear and readable)                  │
│  ☐ Tests run in CI/CD pipeline                                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📦 Complete Examples

### Complete Test Suite for Goa Service

```go
// service/users_service_test.go
package service_test

import (
    "context"
    "testing"
    "time"
    
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    
    "myproject/service"
    "myproject/mocks"
    users "myproject/gen/users"
)

func TestService(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Service Suite")
}

var _ = Describe("UsersService", func() {
    var (
        svc      *service.UsersService
        mockRepo *mocks.MockUserRepository
        ctx      context.Context
    )
    
    BeforeEach(func() {
        ctx = context.Background()
        mockRepo = mocks.NewMockUserRepository()
        svc = service.NewUsersService(mockRepo)
    })
    
    Describe("Get", func() {
        Context("when the user exists", func() {
            var testUser *service.User
            
            BeforeEach(func() {
                testUser = &service.User{
                    ID:        1,
                    Name:      "John Doe",
                    Email:     "john@example.com",
                    CreatedAt: time.Now(),
                    UpdatedAt: time.Now(),
                }
                mockRepo.AddUser(testUser)
            })
            
            It("should return the user", func() {
                result, err := svc.Get(ctx, &users.GetPayload{ID: 1})
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeNil())
                Expect(*result.ID).To(Equal(int64(1)))
                Expect(*result.Name).To(Equal("John Doe"))
                Expect(*result.Email).To(Equal("john@example.com"))
            })
            
            It("should return timestamps", func() {
                result, err := svc.Get(ctx, &users.GetPayload{ID: 1})
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result.CreatedAt).NotTo(BeNil())
                Expect(result.UpdatedAt).NotTo(BeNil())
            })
        })
        
        Context("when the user does not exist", func() {
            It("should return a not found error", func() {
                result, err := svc.Get(ctx, &users.GetPayload{ID: 999})
                
                Expect(err).To(HaveOccurred())
                Expect(result).To(BeNil())
                
                // Check error type
                var serviceErr *users.NotFound
                Expect(errors.As(err, &serviceErr)).To(BeTrue())
            })
        })
    })
    
    Describe("Create", func() {
        DescribeTable("with various inputs",
            func(name, email string, shouldSucceed bool) {
                payload := &users.CreatePayload{
                    Name:  name,
                    Email: email,
                }
                
                result, err := svc.Create(ctx, payload)
                
                if shouldSucceed {
                    Expect(err).NotTo(HaveOccurred())
                    Expect(result).NotTo(BeNil())
                    Expect(*result.Name).To(Equal(name))
                } else {
                    Expect(err).To(HaveOccurred())
                }
            },
            Entry("valid user", "John Doe", "john@example.com", true),
            Entry("another valid user", "Jane Smith", "jane@example.com", true),
        )
        
        Context("when email already exists", func() {
            BeforeEach(func() {
                mockRepo.AddUser(&service.User{
                    ID:    1,
                    Email: "existing@example.com",
                })
            })
            
            It("should return conflict error", func() {
                payload := &users.CreatePayload{
                    Name:  "New User",
                    Email: "existing@example.com",
                }
                
                _, err := svc.Create(ctx, payload)
                
                Expect(err).To(HaveOccurred())
            })
        })
    })
    
    Describe("List", func() {
        BeforeEach(func() {
            for i := 1; i <= 25; i++ {
                mockRepo.AddUser(&service.User{
                    ID:    int64(i),
                    Name:  fmt.Sprintf("User %d", i),
                    Email: fmt.Sprintf("user%d@example.com", i),
                })
            }
        })
        
        DescribeTable("pagination",
            func(page, perPage, expectedCount int, expectedTotal int64) {
                p := int32(page)
                pp := int32(perPage)
                payload := &users.ListPayload{
                    Page:    &p,
                    PerPage: &pp,
                }
                
                result, err := svc.List(ctx, payload)
                
                Expect(err).NotTo(HaveOccurred())
                Expect(result.Users).To(HaveLen(expectedCount))
                Expect(*result.Total).To(Equal(expectedTotal))
            },
            Entry("first page", 1, 10, 10, int64(25)),
            Entry("second page", 2, 10, 10, int64(25)),
            Entry("last page", 3, 10, 5, int64(25)),
            Entry("custom page size", 1, 5, 5, int64(25)),
        )
    })
})
```

---

## 📝 Summary

### Testing Tools Comparison

| Tool | When to Use |
|------|-------------|
| `testing` | Always - standard library, required |
| `testify` | Need better assertions or mocking |
| `httptest` | Testing HTTP handlers |
| `gomock` | Interface mocking with code generation |
| `Ginkgo/Gomega` | BDD style, complex test scenarios |
| `testcontainers` | Integration tests with Docker |

### Key Testing Patterns

1. **Table-Driven Tests** - Multiple cases, one test function
2. **Mock Repositories** - Isolate service logic
3. **HTTP Test Server** - Integration testing HTTP endpoints
4. **Generated Client Tests** - End-to-end API testing
5. **BDD with Ginkgo** - Human-readable test specifications

### Test Commands Reference

```bash
# Standard Go testing
go test ./...                 # All packages
go test -v                    # Verbose
go test -run TestName         # Specific test
go test -cover                # Coverage
go test -bench .              # Benchmarks

# Ginkgo
ginkgo ./...                  # All packages
ginkgo -v                     # Verbose
ginkgo -p                     # Parallel
ginkgo watch                  # Watch mode
ginkgo --cover                # Coverage
```

---

## 📋 Knowledge Check

Before proceeding, ensure you can:

- [ ] Write table-driven tests using `testing` package
- [ ] Create subtests with `t.Run()`
- [ ] Use `httptest` for HTTP handler testing
- [ ] Create mock implementations for interfaces
- [ ] Use `testify/mock` for mocking
- [ ] Write integration tests with test databases
- [ ] Use Goa generated clients for testing
- [ ] Write BDD tests with Ginkgo and Gomega
- [ ] Use Gomega matchers effectively
- [ ] Run tests with coverage reporting

---

## 🔗 Quick Reference Links

- [Go Testing Package](https://pkg.go.dev/testing)
- [httptest Package](https://pkg.go.dev/net/http/httptest)
- [testify](https://github.com/stretchr/testify)
- [gomock](https://github.com/golang/mock)
- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [Gomega Documentation](https://onsi.github.io/gomega/)
- [testcontainers-go](https://github.com/testcontainers/testcontainers-go)

---

> **Next Up:** Part 9 - Deployment (Docker, Kubernetes, CI/CD, Monitoring)
