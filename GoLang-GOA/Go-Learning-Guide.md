# Go Learning Guide

## 1. Foundations and Setup

### 1.1 Installation and Tooling

**Install Go:**
- Download from [golang.org](https://golang.org/dl/)
- Verify: `go version`

**Key Commands:**
```bash
go run main.go       # Compile and run
go build             # Compile to executable
go install           # Install to $GOPATH/bin
go mod init <module> # Initialize module
go mod tidy          # Clean dependencies
```

**Editor Setup:**
- VS Code + Go extension
- Auto-formatting, linting, intellisense enabled

### 1.2 Project Structure

**Module System:**
```go
// go.mod
module github.com/username/myproject

go 1.21

require (
    github.com/pkg/errors v0.9.1
)
```

**Package Organization:**
```
myproject/
├── go.mod
├── main.go
├── internal/          # Private packages
│   └── handler/
├── pkg/              # Public packages
│   └── utils/
└── cmd/              # Multiple binaries
    ├── server/
    └── cli/
```

**Import Paths:**
```go
import (
    "fmt"                                    // Standard library
    "github.com/username/myproject/pkg/utils" // Your package
    "github.com/external/lib"                 // External package
)
```

### 1.3 Basic Syntax

**Variables and Constants:**
```go
// Variable declaration
var name string = "John"
var age int         // Zero value: 0
var isActive bool   // Zero value: false

// Short declaration (inside functions only)
count := 42
message := "Hello"

// Multiple declarations
var (
    host = "localhost"
    port = 8080
)

// Constants
const Pi = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)
```

**Control Flow:**
```go
// If statement
if x > 10 {
    fmt.Println("Greater")
} else if x == 10 {
    fmt.Println("Equal")
} else {
    fmt.Println("Less")
}

// If with initialization
if err := doSomething(); err != nil {
    return err
}

// For loop (only loop in Go)
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// While-style loop
for count < 10 {
    count++
}

// Infinite loop
for {
    // break to exit
}

// Range loop
nums := []int{1, 2, 3}
for index, value := range nums {
    fmt.Println(index, value)
}

// Switch
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("TGIF")
default:
    fmt.Println("Regular day")
}

// Switch without condition (cleaner if-else)
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "C"
}
```

**Functions:**
```go
// Basic function
func add(a, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // naked return
}

// Variadic functions
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

result := sum(1, 2, 3, 4) // 10
```

---

## 2. Core Language Features

### 2.1 Types

**Basic Types:**
```go
// Integers
var i8 int8      // -128 to 127
var u8 uint8     // 0 to 255
var i32 int32
var i64 int64
var i int        // Platform dependent (32 or 64 bit)

// Floats
var f32 float32
var f64 float64

// Others
var b byte       // alias for uint8
var r rune       // alias for int32 (Unicode code point)
var s string
var bl bool
```

**Arrays (Fixed Size):**
```go
var arr [5]int                    // [0 0 0 0 0]
arr2 := [3]string{"a", "b", "c"}
arr3 := [...]int{1, 2, 3}         // Length inferred

// Arrays are values (copied when assigned)
a := [2]int{1, 2}
b := a           // b is a copy
b[0] = 100       // a is unchanged
```

**Slices (Dynamic):**
```go
// Create slices
var s []int                       // nil slice
s = []int{1, 2, 3}
s = make([]int, 5)                // length=5, capacity=5
s = make([]int, 3, 10)            // length=3, capacity=10

// Slice operations
nums := []int{1, 2, 3, 4, 5}
sub := nums[1:4]                  // [2 3 4], shares backing array
sub = nums[:3]                    // [1 2 3]
sub = nums[2:]                    // [3 4 5]

// Append (may reallocate)
nums = append(nums, 6)
nums = append(nums, 7, 8, 9)
nums = append(nums, []int{10, 11}...)

// Length vs Capacity
s := make([]int, 3, 5)
fmt.Println(len(s))               // 3 (current elements)
fmt.Println(cap(s))               // 5 (allocated space)

// Copy slices
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)
```

**Slice Internals:**
```go
// A slice is: {pointer, length, capacity}
// Multiple slices can share the same backing array

original := []int{1, 2, 3, 4, 5}
slice1 := original[0:3]  // [1 2 3]
slice2 := original[2:5]  // [3 4 5]

// Modifying shared backing array
slice1[2] = 999
// original is now [1 2 999 4 5]
// slice2 is now [999 4 5]

// To avoid sharing, copy:
independent := make([]int, len(slice1))
copy(independent, slice1)
```

**Maps:**
```go
// Create maps
var m map[string]int              // nil map (can't add to it)
m = make(map[string]int)
m = map[string]int{
    "Alice": 25,
    "Bob":   30,
}

// Operations
m["Charlie"] = 35                 // Add/update
age := m["Alice"]                 // Get
age, ok := m["Unknown"]           // Check existence
if !ok {
    fmt.Println("Not found")
}
delete(m, "Bob")                  // Delete

// Iterate
for key, value := range m {
    fmt.Println(key, value)
}
// Note: map iteration order is random
```

**Strings:**
```go
s := "Hello, 世界"

// Strings are immutable
// Length in bytes, not runes
fmt.Println(len(s))               // 13 bytes

// Iterate over runes
for i, r := range s {
    fmt.Printf("%d: %c\n", i, r)
}

// String manipulation
import "strings"
strings.Contains(s, "Hello")      // true
strings.Split("a,b,c", ",")       // ["a" "b" "c"]
strings.Join([]string{"a", "b"}, "-") // "a-b"
strings.ToUpper(s)
strings.TrimSpace("  text  ")

// String building (efficient)
import "strings"
var builder strings.Builder
builder.WriteString("Hello")
builder.WriteString(" World")
result := builder.String()
```

### 2.2 Pointers

**Basics:**
```go
x := 42
p := &x          // p is pointer to x (address of x)
fmt.Println(*p)  // 42 (dereference)
*p = 21          // Modify x through pointer
fmt.Println(x)   // 21

// Zero value of pointer is nil
var ptr *int
if ptr == nil {
    fmt.Println("nil pointer")
}
```

**When to Use Pointers:**
```go
// 1. Modify function parameters
func increment(n *int) {
    *n++
}
val := 5
increment(&val)  // val is now 6

// 2. Avoid copying large structs
type LargeStruct struct {
    data [1000000]int
}

func process(ls *LargeStruct) {
    // Works with original, no copy
}

// 3. Indicate optional/nullable values
func find(id int) *User {
    // Return nil if not found
    if notFound {
        return nil
    }
    return &user
}
```

### 2.3 Structs and Methods

**Struct Definition:**
```go
type Person struct {
    Name    string
    Age     int
    Email   string
    address string  // unexported (private)
}

// Create structs
p1 := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}
p2 := Person{"Bob", 25, "bob@example.com", ""}  // Order matters
var p3 Person  // Zero values

// Anonymous structs
config := struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}
```

**Methods:**
```go
// Value receiver (operates on copy)
func (p Person) Greet() string {
    return "Hello, I'm " + p.Name
}

// Pointer receiver (can modify original)
func (p *Person) HaveBirthday() {
    p.Age++
}

// When to use pointer receiver:
// 1. Need to modify receiver
// 2. Receiver is large struct
// 3. Consistency (if one method uses pointer, all should)

p := Person{Name: "Alice", Age: 30}
p.Greet()         // Works
p.HaveBirthday()  // Works (Go auto-references)
```

**Composition:**
```go
// Embedding (composition over inheritance)
type Address struct {
    Street string
    City   string
}

type Employee struct {
    Person          // Embedded
    Address         // Embedded
    EmployeeID int
}

e := Employee{
    Person:     Person{Name: "Alice", Age: 30},
    Address:    Address{Street: "Main St", City: "NYC"},
    EmployeeID: 123,
}

// Access embedded fields directly
fmt.Println(e.Name)    // From Person
fmt.Println(e.City)    // From Address
fmt.Println(e.Greet()) // From Person's method
```

### 2.4 Interfaces

**Interface Definition:**
```go
// Implicit implementation (no "implements" keyword)
type Writer interface {
    Write([]byte) (int, error)
}

type Reader interface {
    Read([]byte) (int, error)
}

// Combining interfaces
type ReadWriter interface {
    Reader
    Writer
}

// Any type that has Write method implements Writer
type FileWriter struct {
    path string
}

func (fw FileWriter) Write(data []byte) (int, error) {
    // Implementation
    return len(data), nil
}

// FileWriter automatically implements Writer interface
var w Writer = FileWriter{path: "file.txt"}
```

**Common Patterns:**
```go
// Empty interface (any type)
var any interface{}
any = 42
any = "hello"
any = Person{Name: "Alice"}

// Type assertion
var i interface{} = "hello"
s := i.(string)        // Panics if wrong type
s, ok := i.(string)    // Safe, ok is false if wrong type

// Type switch
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case Person:
        fmt.Printf("Person: %s\n", v.Name)
    default:
        fmt.Printf("Unknown type\n")
    }
}
```

**Interface Design:**
```go
// Small interfaces are better
type Stringer interface {
    String() string
}

// Accept interfaces, return structs
func Process(r Reader) *Result {  // Good
    // ...
}

// Avoid empty interface pitfalls
func Bad(data interface{}) {      // Too generic
    // Hard to use, requires type assertions
}

func Good(data string) {          // Specific
    // Clear contract
}
```

### 2.5 Packages and Visibility

**Exported vs Unexported:**
```go
package math

// Exported (public) - starts with capital letter
func Add(a, b int) int {
    return a + b
}

type Calculator struct {
    Brand string        // Exported field
    model string        // Unexported field
}

// Unexported (private) - starts with lowercase
func validate(n int) bool {
    return n > 0
}
```

**Package Organization:**
```go
// myproject/pkg/utils/strings.go
package utils

func Reverse(s string) string {
    // Implementation
}

// myproject/main.go
package main

import "github.com/username/myproject/pkg/utils"

func main() {
    result := utils.Reverse("hello")
}
```

**Init Function:**
```go
package database

import "database/sql"

var db *sql.DB

// Runs automatically when package is imported
func init() {
    var err error
    db, err = sql.Open("postgres", "connection_string")
    if err != nil {
        panic(err)
    }
}
```

### 2.6 Error Handling

**Basic Errors:**
```go
import "errors"

// Create errors
err := errors.New("something went wrong")
err := fmt.Errorf("failed to process: %s", filename)

// Return and check errors
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 0)
if err != nil {
    log.Fatal(err)
}
```

**Sentinel Errors:**
```go
package io

var EOF = errors.New("EOF")
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// Usage
data, err := reader.Read()
if err == io.EOF {
    // Handle end of file
}
```

**Custom Errors:**
```go
// Error with context
type ValidationError struct {
    Field string
    Value interface{}
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Msg)
}

// Usage
func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field: "age",
            Value: age,
            Msg:   "must be non-negative",
        }
    }
    return nil
}
```

**Wrapping Errors (Go 1.13+):**
```go
import "fmt"

// Wrap error with context
func readConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("read config: %w", err)  // %w wraps
    }
    // ...
}

// Unwrap errors
err := readConfig("config.json")
if err != nil {
    // Check wrapped error
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("Config file not found")
    }
    
    // Get underlying error
    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Println("Path:", pathErr.Path)
    }
}
```

---

## 3. Concurrency and Synchronization

### 3.1 Goroutines

**Basic Usage:**
```go
// Regular function call (synchronous)
doWork()

// Goroutine (asynchronous)
go doWork()

// Anonymous function
go func() {
    fmt.Println("Running in goroutine")
}()

// With parameters
go func(msg string) {
    fmt.Println(msg)
}("Hello from goroutine")
```

**Lifecycle Considerations:**
```go
func main() {
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Goroutine done")
    }()
    
    // main might exit before goroutine completes
    time.Sleep(2 * time.Second)  // Bad: waiting arbitrarily
}

// Better: use WaitGroup (see below)
```

### 3.2 Channels

**Channel Basics:**
```go
// Unbuffered channel (synchronous)
ch := make(chan int)

// Buffered channel
ch := make(chan int, 3)  // Capacity of 3

// Send and receive
ch <- 42          // Send (blocks if full)
value := <-ch     // Receive (blocks if empty)

// Close channel
close(ch)

// Receive with check
value, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}

// Range over channel (until closed)
for value := range ch {
    fmt.Println(value)
}
```

**Channel Patterns:**
```go
// 1. Signal completion
done := make(chan bool)
go func() {
    doWork()
    done <- true
}()
<-done  // Wait for completion

// 2. Pipeline
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Usage
nums := generate(1, 2, 3, 4)
squares := square(nums)
for sq := range squares {
    fmt.Println(sq)
}

// 3. Fan-out (multiple workers)
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

jobs := make(chan int, 100)
results := make(chan int, 100)

// Start workers
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}

// Send jobs
for j := 1; j <= 5; j++ {
    jobs <- j
}
close(jobs)

// Collect results
for a := 1; a <= 5; a++ {
    <-results
}
```

**Select Statement:**
```go
// Multiplex channels
select {
case msg1 := <-ch1:
    fmt.Println("Received from ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("Received from ch2:", msg2)
case ch3 <- value:
    fmt.Println("Sent to ch3")
default:
    fmt.Println("No channel ready")
}

// Timeout pattern
select {
case result := <-ch:
    fmt.Println("Got result:", result)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
}

// Cancel pattern
quit := make(chan bool)
go func() {
    for {
        select {
        case <-quit:
            return
        default:
            // Do work
        }
    }
}()
// Later: quit <- true
```

### 3.3 Context

**Context Usage:**
```go
import "context"

// Create contexts
ctx := context.Background()              // Root context
ctx := context.TODO()                    // Placeholder

// With cancellation
ctx, cancel := context.WithCancel(ctx)
defer cancel()  // Always call cancel

// With timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// With deadline
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(ctx, deadline)
defer cancel()

// With values (use sparingly)
ctx = context.WithValue(ctx, "requestID", "12345")
```

**Context Patterns:**
```go
// Cancellable goroutine
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Cancelled:", ctx.Err())
            return
        default:
            // Do work
            time.Sleep(500 * time.Millisecond)
        }
    }
}

ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)
time.Sleep(2 * time.Second)
cancel()  // Stop worker

// HTTP request with timeout
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
resp, err := http.DefaultClient.Do(req)
if err != nil {
    // Handle timeout or cancellation
}
```

### 3.4 Synchronization

**WaitGroup:**
```go
import "sync"

var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Println("Worker", id)
    }(i)
}

wg.Wait()  // Wait for all goroutines
fmt.Println("All done")
```

**Mutex:**
```go
import "sync"

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// RWMutex (multiple readers, single writer)
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()         // Read lock
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()          // Write lock
    defer c.mu.Unlock()
    c.data[key] = value
}
```

**Atomic Operations:**
```go
import "sync/atomic"

var counter int64

// Atomic increment
atomic.AddInt64(&counter, 1)

// Atomic load/store
value := atomic.LoadInt64(&counter)
atomic.StoreInt64(&counter, 100)

// Compare and swap
swapped := atomic.CompareAndSwapInt64(&counter, 100, 200)
```

**sync.Once:**
```go
import "sync"

var once sync.Once
var instance *Singleton

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}  // Initialized only once
    })
    return instance
}
```

**sync.Map:**
```go
// Thread-safe map (use when keys are stable)
var cache sync.Map

// Store
cache.Store("key", "value")

// Load
value, ok := cache.Load("key")

// LoadOrStore
actual, loaded := cache.LoadOrStore("key", "default")

// Delete
cache.Delete("key")

// Range
cache.Range(func(key, value interface{}) bool {
    fmt.Println(key, value)
    return true  // Continue iteration
})

// Note: Regular map + mutex is often better
```

### 3.5 Concurrency Patterns

**Worker Pool:**
```go
func workerPool(jobs <-chan int, results chan<- int, workers int) {
    var wg sync.WaitGroup
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- job * 2  // Process job
            }
        }()
    }
    
    wg.Wait()
    close(results)
}

// Usage
jobs := make(chan int, 100)
results := make(chan int, 100)

go workerPool(jobs, results, 3)

// Send jobs
for i := 1; i <= 10; i++ {
    jobs <- i
}
close(jobs)

// Collect results
for result := range results {
    fmt.Println(result)
}
```

**Fan-in Pattern:**
```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

// Usage
ch1 := generateNumbers(1, 3)
ch2 := generateNumbers(4, 6)
merged := fanIn(ch1, ch2)
```

**Avoiding Goroutine Leaks:**
```go
// Bad: goroutine never exits
func leak() <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; ; i++ {
            ch <- i  // Blocks forever if no receiver
        }
    }()
    return ch
}

// Good: use context for cancellation
func noLeak(ctx context.Context) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case ch <- i:
            case <-ctx.Done():
                return  // Exit goroutine
            }
        }
    }()
    return ch
}
```

---

## 4. Tooling, Testing, and Quality

### 4.1 Go Tooling

**Essential Commands:**
```bash
# Format code (auto-fix)
go fmt ./...

# Vet code (detect suspicious constructs)
go vet ./...

# Static analysis
go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
go vet -vettool=$(which shadow) ./...

# Documentation
go doc fmt.Println
go doc -all net/http

# Generate code
go generate ./...

# Module management
go mod download     # Download dependencies
go mod verify       # Verify dependencies
go mod vendor       # Copy deps to vendor/
```

**golangci-lint (Recommended):**
```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run

# Configuration: .golangci.yml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
```

**go generate:**
```go
//go:generate stringer -type=Status
type Status int

const (
    Pending Status = iota
    Active
    Completed
)

// Run: go generate
// Creates status_string.go with String() method
```

### 4.2 Testing

**Basic Tests:**
```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

// Run tests
// go test
// go test -v           # Verbose
// go test ./...        # All packages
// go test -run TestAdd # Specific test
```

**Table-Driven Tests:**
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 5, 5},
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

**Subtests:**
```go
func TestMath(t *testing.T) {
    t.Run("Add", func(t *testing.T) {
        if Add(1, 2) != 3 {
            t.Error("failed")
        }
    })
    
    t.Run("Subtract", func(t *testing.T) {
        if Subtract(5, 3) != 2 {
            t.Error("failed")
        }
    })
}

// Run: go test -run TestMath/Add
```

**Test Coverage:**
```bash
# Generate coverage
go test -cover
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Coverage by function
go tool cover -func=coverage.out
```

**Benchmarks:**
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

func BenchmarkConcatenate(b *testing.B) {
    b.ResetTimer()  // Reset after setup
    for i := 0; i < b.N; i++ {
        result := "Hello" + " " + "World"
        _ = result
    }
}

// Run benchmarks
// go test -bench=.
// go test -bench=Add -benchmem  # Show memory allocations
```

**Test Helpers:**
```go
func TestSomething(t *testing.T) {
    t.Helper()  // Mark as helper (better error reporting)
    // ...
}

// Setup and teardown
func TestMain(m *testing.M) {
    // Setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Teardown
    teardown()
    
    os.Exit(code)
}
```

### 4.3 Mocks and Testability

**Interface-Driven Design:**
```go
// Define interface
type UserRepository interface {
    GetUser(id int) (*User, error)
    SaveUser(user *User) error
}

// Real implementation
type DBUserRepository struct {
    db *sql.DB
}

func (r *DBUserRepository) GetUser(id int) (*User, error) {
    // Database logic
}

// Mock implementation for testing
type MockUserRepository struct {
    Users map[int]*User
}

func (r *MockUserRepository) GetUser(id int) (*User, error) {
    user, ok := r.Users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

// Service depends on interface
type UserService struct {
    repo UserRepository
}

// Test with mock
func TestUserService(t *testing.T) {
    mockRepo := &MockUserRepository{
        Users: map[int]*User{
            1: {ID: 1, Name: "Alice"},
        },
    }
    
    service := &UserService{repo: mockRepo}
    user, err := service.GetUserByID(1)
    // Assertions...
}
```

**httptest:**
```go
import "net/http/httptest"

func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/users", nil)
    w := httptest.NewRecorder()
    
    handler := UserHandler()
    handler.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
    
    // Check response body
    body := w.Body.String()
    if !strings.Contains(body, "users") {
        t.Error("Response missing users")
    }
}

// Test server
func TestClient(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    }))
    defer server.Close()
    
    resp, err := http.Get(server.URL)
    // Test client logic...
}
```

### 4.4 Modules and Versioning

**Semantic Versioning:**
```bash
# Tag releases
git tag v1.0.0
git push origin v1.0.0

# Version format: vMAJOR.MINOR.PATCH
# v1.2.3
#  ^ ^ ^
#  | | +-- Patch (bug fixes)
#  | +---- Minor (new features, backward compatible)
#  +------ Major (breaking changes)
```

**go.mod Directives:**
```go
module github.com/user/project

go 1.21

require (
    github.com/pkg/errors v0.9.1
    github.com/gorilla/mux v1.8.0
)

// Replace directive (for local development)
replace github.com/user/dep => ../dep

// Exclude directive
exclude github.com/old/dep v1.0.0

// Indirect dependencies
require github.com/some/dep v1.0.0 // indirect
```

**Private Modules:**
```bash
# Configure Git for private repos
git config --global url."git@github.com:".insteadOf "https://github.com/"

# Set GOPRIVATE
export GOPRIVATE=github.com/mycompany/*

# Or in go.env
go env -w GOPRIVATE=github.com/mycompany/*
```

**Vendoring:**
```bash
# Copy dependencies to vendor/
go mod vendor

# Use vendored dependencies
go build -mod=vendor
```

### 4.5 Profiling and Debugging

**pprof - CPU Profiling:**
```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Your code here
}

// Analyze
// go tool pprof cpu.prof
// (pprof) top
// (pprof) list functionName
// (pprof) web  # Generate graph (requires graphviz)
```

**pprof - Memory Profiling:**
```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    // Your code here
    
    f, _ := os.Create("mem.prof")
    pprof.WriteHeapProfile(f)
    f.Close()
}

// Analyze
// go tool pprof mem.prof
```

**HTTP pprof:**
```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Your application code
}

// Access profiles:
// http://localhost:6060/debug/pprof/
// go tool pprof http://localhost:6060/debug/pprof/heap
// go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

**Race Detector:**
```bash
# Detect data races
go test -race
go run -race main.go
go build -race

# Example race condition
var counter int

func increment() {
    counter++  // Race!
}

go increment()
go increment()
```

**Debugging with Delve:**
```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug main.go

# Common commands
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) step
(dlv) print variable
(dlv) quit
```

---

## 5. Standard Library and Common Tasks

### 5.1 I/O and Files

**Reading Files:**
```go
import (
    "os"
    "io"
    "bufio"
)

// Read entire file
data, err := os.ReadFile("file.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))

// Read with os.Open
file, err := os.Open("file.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Read all
content, err := io.ReadAll(file)

// Buffered reading (efficient for large files)
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}
if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
```

**Writing Files:**
```go
// Write entire file
data := []byte("Hello, World!")
err := os.WriteFile("file.txt", data, 0644)

// Write with os.Create
file, err := os.Create("file.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

file.WriteString("Hello\n")
file.Write([]byte("World\n"))

// Buffered writing
writer := bufio.NewWriter(file)
writer.WriteString("Buffered write\n")
writer.Flush()  // Important!
```

**File Operations:**
```go
// Check if exists
if _, err := os.Stat("file.txt"); os.IsNotExist(err) {
    fmt.Println("File does not exist")
}

// Delete
os.Remove("file.txt")

// Rename/Move
os.Rename("old.txt", "new.txt")

// Create directory
os.Mkdir("mydir", 0755)
os.MkdirAll("path/to/dir", 0755)  // Create parents too

// Remove directory
os.Remove("mydir")
os.RemoveAll("mydir")  // Recursive

// List directory
entries, err := os.ReadDir(".")
for _, entry := range entries {
    fmt.Println(entry.Name(), entry.IsDir())
}
```

**io.Reader and io.Writer:**
```go
// Copy from reader to writer
src, _ := os.Open("source.txt")
dst, _ := os.Create("dest.txt")
defer src.Close()
defer dst.Close()

io.Copy(dst, src)

// Chain readers/writers
file, _ := os.Open("file.gz")
gzipReader, _ := gzip.NewReader(file)
scanner := bufio.NewScanner(gzipReader)
```

### 5.2 JSON and Encoding

**JSON Marshaling:**
```go
import "encoding/json"

type Person struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email,omitempty"`  // Omit if empty
    password string  // Not exported, won't be marshaled
}

// Struct to JSON
person := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}
jsonData, err := json.Marshal(person)
// {"name":"Alice","age":30,"email":"alice@example.com"}

// Pretty print
jsonData, err := json.MarshalIndent(person, "", "  ")

// JSON to struct
jsonStr := `{"name":"Bob","age":25}`
var p Person
err := json.Unmarshal([]byte(jsonStr), &p)
```

**JSON Streaming:**
```go
// Encode to writer
file, _ := os.Create("data.json")
encoder := json.NewEncoder(file)
encoder.Encode(person)

// Decode from reader
file, _ := os.Open("data.json")
decoder := json.NewDecoder(file)
var p Person
decoder.Decode(&p)

// Decode array
decoder := json.NewDecoder(strings.NewReader(`[{"name":"Alice"},{"name":"Bob"}]`))
var people []Person
decoder.Decode(&people)
```

**Custom JSON Marshaling:**
```go
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
    return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
    str := string(data[1 : len(data)-1])  // Remove quotes
    parsed, err := time.Parse(time.RFC3339, str)
    if err != nil {
        return err
    }
    *t = Timestamp(parsed)
    return nil
}
```

**Other Encodings:**
```go
import (
    "encoding/xml"
    "encoding/csv"
    "encoding/gob"
)

// XML
type Book struct {
    XMLName xml.Name `xml:"book"`
    Title   string   `xml:"title"`
    Author  string   `xml:"author"`
}

xmlData, _ := xml.MarshalIndent(book, "", "  ")

// CSV
file, _ := os.Create("data.csv")
writer := csv.NewWriter(file)
writer.Write([]string{"Name", "Age"})
writer.Write([]string{"Alice", "30"})
writer.Flush()

// Gob (Go binary format, efficient)
var buf bytes.Buffer
encoder := gob.NewEncoder(&buf)
encoder.Encode(person)

decoder := gob.NewDecoder(&buf)
var p Person
decoder.Decode(&p)
```

### 5.3 HTTP

**HTTP Server:**
```go
import "net/http"

// Simple handler function
func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello, World!"))
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/api/users", usersHandler)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler with method check
func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // List users
        json.NewEncoder(w).Encode(users)
    case http.MethodPost:
        // Create user
        var user User
        json.NewDecoder(r.Body).Decode(&user)
        // Save user...
        w.WriteHeader(http.StatusCreated)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}
```

**http.Handler Interface:**
```go
type MyHandler struct {
    db *sql.DB
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Handle request
    w.Write([]byte("Response"))
}

// Usage
handler := &MyHandler{db: db}
http.ListenAndServe(":8080", handler)
```

**Middleware Pattern:**
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Chain middlewares
handler := loggingMiddleware(authMiddleware(http.HandlerFunc(helloHandler)))
http.ListenAndServe(":8080", handler)
```

**Request Parsing:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    // URL parameters
    id := r.URL.Query().Get("id")
    
    // Form data
    r.ParseForm()
    name := r.FormValue("name")
    
    // JSON body
    var data map[string]interface{}
    json.NewDecoder(r.Body).Decode(&data)
    defer r.Body.Close()
    
    // Headers
    contentType := r.Header.Get("Content-Type")
    
    // Cookies
    cookie, err := r.Cookie("session")
}
```

**HTTP Client:**
```go
// Simple GET
resp, err := http.Get("https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)

// POST with JSON
user := User{Name: "Alice"}
jsonData, _ := json.Marshal(user)

resp, err := http.Post(
    "https://api.example.com/users",
    "application/json",
    bytes.NewBuffer(jsonData),
)

// Custom request
req, _ := http.NewRequest("PUT", "https://api.example.com/users/1", nil)
req.Header.Set("Authorization", "Bearer token")
req.Header.Set("Content-Type", "application/json")

client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)

// With context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
resp, err := http.DefaultClient.Do(req)
```

**Server Configuration:**
```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}

log.Fatal(server.ListenAndServe())

// Graceful shutdown
go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}()

// Wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Shutdown(ctx)
```

### 5.4 Time and Context

**Time Basics:**
```go
import "time"

// Current time
now := time.Now()
fmt.Println(now)

// Create time
t := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

// Parsing
layout := "2006-01-02 15:04:05"  // Reference time
t, err := time.Parse(layout, "2024-01-01 12:00:00")

// Common layouts
time.RFC3339    // "2006-01-02T15:04:05Z07:00"
time.RFC822     // "02 Jan 06 15:04 MST"
time.Kitchen    // "3:04PM"

// Formatting
formatted := now.Format("2006-01-02 15:04:05")

// Durations
duration := 5 * time.Second
duration := time.Hour + 30*time.Minute

// Time arithmetic
future := now.Add(24 * time.Hour)
past := now.Add(-1 * time.Hour)
diff := future.Sub(now)

// Comparison
if t1.Before(t2) { }
if t1.After(t2) { }
if t1.Equal(t2) { }
```

**Timers and Tickers:**
```go
// Timer (one-time)
timer := time.NewTimer(2 * time.Second)
<-timer.C
fmt.Println("Timer expired")

// Stop timer
if !timer.Stop() {
    <-timer.C
}

// Ticker (repeating)
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for {
    select {
    case t := <-ticker.C:
        fmt.Println("Tick at", t)
    }
}

// After (simpler one-time)
select {
case <-time.After(5 * time.Second):
    fmt.Println("Timeout")
}
```

**Context with Time:**
```go
// Timeout
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

select {
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())  // context deadline exceeded
}

// Deadline
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

### 5.5 CLI Applications

**Command-Line Flags:**
```go
import "flag"

func main() {
    // Define flags
    host := flag.String("host", "localhost", "Server host")
    port := flag.Int("port", 8080, "Server port")
    verbose := flag.Bool("verbose", false, "Enable verbose logging")
    
    // Parse flags
    flag.Parse()
    
    // Use flags
    fmt.Printf("Starting server on %s:%d\n", *host, *port)
    if *verbose {
        log.SetLevel(log.DebugLevel)
    }
    
    // Positional arguments
    args := flag.Args()
    fmt.Println("Arguments:", args)
}

// Usage:
// go run main.go -host=0.0.0.0 -port=9000 -verbose
// go run main.go -h  # Show help
```

**Environment Variables:**
```go
import "os"

// Get environment variable
host := os.Getenv("HOST")
if host == "" {
    host = "localhost"  // Default
}

// Set environment variable
os.Setenv("DEBUG", "true")

// Check if set
value, exists := os.LookupEnv("API_KEY")
if !exists {
    log.Fatal("API_KEY not set")
}

// Get all environment variables
for _, env := range os.Environ() {
    fmt.Println(env)
}
```

**Logging:**
```go
import "log"

// Basic logging
log.Println("Info message")
log.Printf("User %s logged in", username)
log.Fatal("Critical error")  // Exits with status 1
log.Panic("Panic!")          // Calls panic()

// Configure logger
log.SetPrefix("[APP] ")
log.SetFlags(log.LstdFlags | log.Lshortfile)
// Output: [APP] 2024/01/01 12:00:00 main.go:10: Info message

// Custom logger
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
logger := log.New(file, "[CUSTOM] ", log.LstdFlags)
logger.Println("Custom log")

// Structured logging (use external library)
// Popular: github.com/sirupsen/logrus, go.uber.org/zap
```

**Complete CLI Example:**
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
)

func main() {
    // Flags
    var (
        configPath = flag.String("config", "", "Config file path")
        verbose    = flag.Bool("v", false, "Verbose output")
        port       = flag.Int("port", 8080, "Server port")
    )
    
    flag.Parse()
    
    // Environment variables with fallback
    host := os.Getenv("HOST")
    if host == "" {
        host = "localhost"
    }
    
    // Logging setup
    if *verbose {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    }
    
    // Validate required flags
    if *configPath == "" {
        fmt.Fprintf(os.Stderr, "Error: -config flag is required\n")
        flag.Usage()
        os.Exit(1)
    }
    
    // Run application
    log.Printf("Starting server on %s:%d", host, *port)
    log.Printf("Using config: %s", *configPath)
    
    // Your application logic...
}
```

### 5.6 Database Access

**database/sql Basics:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // PostgreSQL driver
)

// Open connection
db, err := sql.Open("postgres", "host=localhost user=postgres password=secret dbname=mydb sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Test connection
if err := db.Ping(); err != nil {
    log.Fatal(err)
}

// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

**Queries:**
```go
// Query single row
var name string
var age int
err := db.QueryRow("SELECT name, age FROM users WHERE id = $1", 1).Scan(&name, &age)
if err == sql.ErrNoRows {
    fmt.Println("No user found")
} else if err != nil {
    log.Fatal(err)
}

// Query multiple rows
rows, err := db.Query("SELECT id, name, age FROM users WHERE age > $1", 18)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    var age int
    if err := rows.Scan(&id, &name, &age); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%d: %s (%d)\n", id, name, age)
}

if err := rows.Err(); err != nil {
    log.Fatal(err)
}
```

**Inserts and Updates:**
```go
// Insert
result, err := db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "Alice", 30)
if err != nil {
    log.Fatal(err)
}

id, _ := result.LastInsertId()
rowsAffected, _ := result.RowsAffected()

// Update
result, err := db.Exec("UPDATE users SET age = $1 WHERE id = $2", 31, 1)

// Delete
result, err := db.Exec("DELETE FROM users WHERE id = $1", 1)
```

**Transactions:**
```go
// Begin transaction
tx, err := db.Begin()
if err != nil {
    log.Fatal(err)
}

// Execute queries
_, err = tx.Exec("INSERT INTO accounts (user_id, balance) VALUES ($1, $2)", 1, 100)
if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

_, err = tx.Exec("UPDATE users SET account_created = true WHERE id = $1", 1)
if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

// Commit
if err := tx.Commit(); err != nil {
    log.Fatal(err)
}
```

**Prepared Statements:**
```go
// Prepare statement
stmt, err := db.Prepare("SELECT name FROM users WHERE id = $1")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

// Use multiple times
var name string
stmt.QueryRow(1).Scan(&name)
stmt.QueryRow(2).Scan(&name)
```

**Context with Database:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Query with context
rows, err := db.QueryContext(ctx, "SELECT * FROM users")

// Execute with context
_, err := db.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")

// Transaction with context
tx, err := db.BeginTx(ctx, nil)
```

**Example Repository Pattern:**
```go
type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
    user := &User{}
    err := r.db.QueryRowContext(ctx,
        "SELECT id, name, email, age FROM users WHERE id = $1", id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return user, err
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    err := r.db.QueryRowContext(ctx,
        "INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id",
        user.Name, user.Email, user.Age,
    ).Scan(&user.ID)
    return err
}

func (r *UserRepository) List(ctx context.Context) ([]*User, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, age FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*User
    for rows.Next() {
        user := &User{}
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, rows.Err()
}
```

---

## Quick Reference

### Go Commands
```bash
go run main.go          # Run
go build                # Compile
go install              # Install binary
go test                 # Run tests
go test -v              # Verbose tests
go test -cover          # Coverage
go test -race           # Race detection
go test -bench=.        # Benchmarks
go mod init             # Initialize module
go mod tidy             # Clean dependencies
go fmt ./...            # Format code
go vet ./...            # Static analysis
```

### Common Patterns
```go
// Error handling
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// Defer cleanup
f, err := os.Open("file")
if err != nil {
    return err
}
defer f.Close()

// Check interface implementation at compile time
var _ io.Reader = (*MyType)(nil)

// Zero values
var i int       // 0
var s string    // ""
var p *int      // nil
var sl []int    // nil
var m map[string]int  // nil
```

### Best Practices
- Use `go fmt` always
- Handle all errors explicitly
- Use interfaces for abstraction
- Prefer composition over inheritance
- Keep functions small and focused
- Use table-driven tests
- Accept interfaces, return structs
- Don't ignore errors (use `_` intentionally)
- Use context for cancellation
- Avoid global state
- Make zero value useful
- Use meaningful variable names

---

This guide covers the essential Go concepts with practical examples. Practice by building real projects: CLI tools, REST APIs, concurrent processors, and database applications.
