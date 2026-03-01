# Part 1: Core Foundations (Must Know First)

> **Prerequisites for Goa Framework** - Master these Go essentials before diving into Goa microservices development.

---

## Table of Contents

1. [Structs & Interfaces](#1-structs--interfaces)
2. [Methods & Receivers](#2-methods--receivers)
3. [Error Handling](#3-error-handling)
4. [Context Package](#4-context-package)
5. [Concurrency (Goroutines & Channels)](#5-concurrency-goroutines--channels)
6. [Modules & Dependency Management](#6-modules--dependency-management)

---

## 1. Structs & Interfaces

### 🔹 Structs

Structs are composite data types that group together variables under a single name. They are Go's primary mechanism for creating custom types and are fundamental to organizing data in your applications.

#### Understanding Structs: The Theory

**What is a Struct?**

A struct (short for "structure") is a user-defined type that contains a collection of named fields. Unlike arrays (which store elements of the same type), structs can hold fields of different types. Think of a struct as a blueprint for creating objects that represent real-world entities.

**Memory Layout**

Structs in Go are stored contiguously in memory. The fields are laid out in the order they are declared, with potential padding for alignment. This makes structs efficient for memory access but means field order can affect memory usage.

```
┌─────────────────────────────────────────────────────┐
│                    User Struct                       │
├─────────┬───────────┬───────────┬─────────┬─────────┤
│   ID    │ FirstName │ LastName  │  Email  │ Active  │
│  (int)  │ (string)  │ (string)  │(string) │ (bool)  │
│ 8 bytes │ 16 bytes  │ 16 bytes  │16 bytes │ 1 byte  │
└─────────┴───────────┴───────────┴─────────┴─────────┘
```

**Value Semantics**

Structs in Go have value semantics by default. When you assign a struct to another variable or pass it to a function, a complete copy is made. This is different from languages like Java where objects are passed by reference.

**Why Structs Matter for Goa?**

In Goa framework:
- **Request/Response Types**: Structs define the shape of API payloads
- **DTOs (Data Transfer Objects)**: Move data between layers
- **Domain Models**: Represent business entities
- **Configuration**: Hold service settings

#### Basic Struct Definition

```go
// Define a struct
type User struct {
    ID        int
    FirstName string
    LastName  string
    Email     string
    Active    bool
}

// Create struct instances
func main() {
    // Method 1: Using field names (recommended)
    user1 := User{
        ID:        1,
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
        Active:    true,
    }

    // Method 2: Positional (not recommended - brittle)
    user2 := User{2, "Jane", "Doe", "jane@example.com", false}

    // Method 3: Zero-value initialization
    var user3 User // All fields get zero values
    user3.ID = 3
    user3.FirstName = "Bob"
}
```

#### Embedded Structs (Composition)

**Theory: Composition over Inheritance**

Go deliberately omits class-based inheritance. Instead, it promotes **composition** - building complex types by combining simpler ones. This follows the principle "favor composition over inheritance" from object-oriented design.

**Benefits of Composition:**
- More flexible than inheritance hierarchies
- Avoids the "fragile base class" problem
- Promotes loose coupling
- Makes code easier to test and maintain

**Field Promotion:**
When you embed a struct, its fields and methods are "promoted" to the outer struct. You can access them directly without referencing the embedded type.

```go
// Base struct
type Address struct {
    Street  string
    City    string
    Country string
    ZipCode string
}

// Embedding Address into Employee
type Employee struct {
    ID      int
    Name    string
    Address // Embedded struct (anonymous field)
    Salary  float64
}

func main() {
    emp := Employee{
        ID:   101,
        Name: "Alice",
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
            ZipCode: "10001",
        },
        Salary: 75000.0,
    }

    // Access embedded fields directly
    fmt.Println(emp.City)        // "New York" - promoted field
    fmt.Println(emp.Address.City) // Also works
}
```

#### Struct Tags

**Theory: Metadata for Structs**

Struct tags are string literals attached to struct fields that provide metadata. They are invisible to normal Go code but can be accessed via reflection. Tags are extensively used in:

- **JSON/XML encoding**: Control field names in serialized output
- **Database ORMs**: Map fields to database columns
- **Validation libraries**: Define validation rules
- **Goa framework**: Define API schemas, validation, and documentation

**Tag Syntax:**
```
`key1:"value1" key2:"value2"`
```

**Common Tag Keys:**
| Tag Key | Purpose | Example |
|---------|---------|--------|
| `json` | JSON field name | `json:"user_id"` |
| `xml` | XML element/attribute name | `xml:"userId,attr"` |
| `validate` | Validation rules | `validate:"required,email"` |
| `gorm` | Database mapping | `gorm:"column:user_id"` |
| `mapstructure` | Config decoding | `mapstructure:"api_key"` |

```go
type APIResponse struct {
    UserID    int    `json:"user_id" xml:"userId" validate:"required"`
    Username  string `json:"username" validate:"required,min=3,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"-"` // Excluded from JSON
    CreatedAt string `json:"created_at,omitempty"` // Omit if empty
}
```

---

### 🔹 Interfaces

Interfaces define behavior through method signatures. They enable polymorphism and loose coupling.

#### Understanding Interfaces: The Theory

**What is an Interface?**

An interface is a type that specifies a set of method signatures. Any type that implements all the methods of an interface automatically satisfies that interface - there's no explicit declaration required (implicit implementation).

**The Duck Typing Philosophy:**
> "If it walks like a duck and quacks like a duck, it's a duck."

Go's interfaces embody this principle. You don't say "this type implements this interface" - if a type has the right methods, it implements the interface automatically.

**Why Interfaces are Powerful:**

1. **Decoupling**: Code depends on behaviors, not concrete types
2. **Testability**: Easy to create mocks and stubs
3. **Flexibility**: New types can satisfy existing interfaces
4. **Composition**: Small interfaces combine into larger ones

**Interface Internals:**

Under the hood, an interface value consists of two components:
```
┌─────────────────────────────────┐
│        Interface Value          │
├────────────────┬────────────────┤
│     Type       │     Data       │
│   (pointer to  │   (pointer to  │
│   type info)   │   actual data) │
└────────────────┴────────────────┘
```

**The Zero Value of Interface:**
A nil interface has both type and value set to nil. This is different from an interface holding a nil pointer!

**Why Interfaces Matter for Goa?**
- Services are defined as interfaces
- Endpoints implement interface contracts
- Middleware uses interfaces for extensibility
- Testing relies heavily on interface mocking

#### Basic Interface

```go
// Define an interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Compose interfaces
type ReadWriter interface {
    Reader
    Writer
}
```

#### Implementing Interfaces (Implicit)

```go
// Interface definition
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle implements Shape implicitly
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle also implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Polymorphic function
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}

    PrintShapeInfo(rect)   // Works!
    PrintShapeInfo(circle) // Works!
}
```

#### Empty Interface

**Theory: The Universal Container**

The empty interface `interface{}` (or `any` in Go 1.18+) has no methods, so every type satisfies it. This makes it a universal container, but use it sparingly as you lose type safety.

**When to Use Empty Interface:**
- Heterogeneous collections (rare)
- Generic utility functions (before generics)
- Interfacing with untyped data (JSON, configs)
- Framework internals

**When NOT to Use:**
- Regular business logic
- When the type is known at compile time
- When generics (Go 1.18+) would work better

```go
// Empty interface can hold any type
func PrintAnything(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

// Modern Go (1.18+) uses 'any' as alias
func PrintAnythingNew(v any) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}
```

#### Type Assertions & Type Switch

**Theory: Recovering Type Information**

When you have an interface value, you may need to access the underlying concrete type. Go provides two mechanisms:

1. **Type Assertion**: Extract a specific type from an interface
2. **Type Switch**: Branch based on the underlying type

**Type Assertion Syntax:**
```go
value, ok := interfaceVar.(ConcreteType)
```
- If successful: `value` contains the concrete value, `ok` is `true`
- If fails: `value` is zero value, `ok` is `false`
- Without `ok`: Panics if assertion fails!

**Best Practice:** Always use the two-value form to avoid panics.

```go
func processValue(v interface{}) {
    // Type assertion
    if str, ok := v.(string); ok {
        fmt.Println("String:", str)
        return
    }

    // Type switch
    switch val := v.(type) {
    case int:
        fmt.Println("Integer:", val)
    case float64:
        fmt.Println("Float:", val)
    case string:
        fmt.Println("String:", val)
    case []int:
        fmt.Println("Int slice:", val)
    default:
        fmt.Println("Unknown type")
    }
}
```

---

## 2. Methods & Receivers

### Understanding Methods: The Theory

**What is a Method?**

A method is a function with a special **receiver** argument that appears between the `func` keyword and the method name. The receiver binds the method to a type, enabling object-oriented programming patterns in Go.

**Methods vs Functions:**
```go
// Function
func PrintUser(u User) { ... }

// Method
func (u User) Print() { ... }
```

Methods provide:
- Better organization (methods are namespaced to types)
- Interface implementation
- Cleaner API (object.Method() vs Function(object))

**The Receiver Concept:**

The receiver is essentially the first parameter of the method, but with special syntax. When you call `user.Print()`, Go passes `user` as the receiver.

**Method Sets:**

Every type has a "method set" - the collection of methods accessible on that type:
- Type `T`: Methods with receiver `T`
- Type `*T`: Methods with receiver `T` OR `*T`

This asymmetry is important for interface satisfaction.

### 🔹 Value Receivers vs Pointer Receivers

**The Core Decision:**

Choosing between value and pointer receivers is one of the most important decisions in Go programming. It affects mutability, performance, and interface satisfaction.

**Value Receiver Mechanics:**
```
Original Struct → Copy Created → Method Operates on Copy
                                         ↓
                              Changes Lost When Method Returns
```

**Pointer Receiver Mechanics:**
```
Original Struct → Pointer Passed → Method Operates on Original
                                         ↓
                              Changes Persist After Method Returns
```

```go
type Counter struct {
    count int
}

// Value receiver - operates on a COPY
func (c Counter) GetCount() int {
    return c.count
}

// Value receiver - modification doesn't affect original
func (c Counter) IncrementWrong() {
    c.count++ // Only modifies the copy!
}

// Pointer receiver - operates on the ORIGINAL
func (c *Counter) Increment() {
    c.count++ // Modifies the original
}

// Pointer receiver - for large structs (avoids copying)
func (c *Counter) Reset() {
    c.count = 0
}

func main() {
    counter := Counter{count: 0}
    
    counter.IncrementWrong()
    fmt.Println(counter.count) // 0 - not changed!
    
    counter.Increment()
    fmt.Println(counter.count) // 1 - changed!
}
```

### 🔹 When to Use Which Receiver

| Use Pointer Receiver `*T` | Use Value Receiver `T` |
|---------------------------|------------------------|
| Method needs to modify the receiver | Method only reads data |
| Struct is large (avoid copying) | Struct is small (few fields) |
| Consistency (if any method uses pointer) | Immutability is desired |
| Struct contains sync.Mutex or similar | Basic types or small structs |

**Rule of Thumb:** When in doubt, use pointer receivers. They're consistent, efficient for larger structs, and allow modification when needed.

**Size Guidelines:**
- Small struct (≤3 fields, basic types): Value receiver is fine
- Medium struct (4-10 fields): Consider pointer
- Large struct (>10 fields or contains slices/maps): Use pointer

**Consistency Principle:**
If ANY method on a type uses a pointer receiver, ALL methods should use pointer receivers for consistency.

### 🔹 Method Sets and Interface Satisfaction

```go
type Printer interface {
    Print()
}

type Document struct {
    Content string
}

// Pointer receiver
func (d *Document) Print() {
    fmt.Println(d.Content)
}

func main() {
    doc := Document{Content: "Hello"}
    
    // This works - Go automatically takes the address
    doc.Print()
    
    var p Printer
    // p = doc       // ERROR: Document doesn't implement Printer
    p = &doc         // OK: *Document implements Printer
    p.Print()
}
```

**Rule:** 
- Value receivers: Both `T` and `*T` satisfy the interface
- Pointer receivers: Only `*T` satisfies the interface

---

## 3. Error Handling

### Understanding Error Handling: The Theory

**Go's Philosophy on Errors**

Go treats errors as values, not exceptions. This is a deliberate design choice that makes error handling explicit and predictable. Unlike try-catch blocks that can hide control flow, Go's error handling is visible in the code structure.

**Why No Exceptions?**

1. **Explicit Control Flow**: You always see where errors can occur
2. **Forcing Handling**: Can't accidentally ignore errors (though you can explicitly with `_`)
3. **Performance**: No stack unwinding overhead
4. **Simplicity**: One mechanism for all error cases

**The Error Handling Pattern:**
```go
result, err := SomeOperation()
if err != nil {
    // Handle error - return, log, retry, etc.
    return err
}
// Continue with result
```

This pattern appears thousands of times in Go code. It's intentionally verbose to make error handling obvious.

**Errors vs Panics:**
| Errors | Panics |
|--------|--------|
| Expected failure conditions | Unexpected/unrecoverable situations |
| Returned as values | Unwind the stack |
| Caller decides how to handle | Usually crash the program |
| Network failures, invalid input | Nil pointer dereference, index out of bounds |

**Error Handling in Goa:**
Goa transforms errors into HTTP status codes and response bodies. Understanding Go's error patterns is essential for proper API error responses.

### 🔹 The `error` Interface

**The Simplest Interface:**

The `error` interface is remarkably simple - just one method. Any type that implements `Error() string` is an error. This simplicity enables powerful patterns.

```go
// Built-in error interface
type error interface {
    Error() string
}
```

### 🔹 Creating Errors

**Three Levels of Error Complexity:**

1. **Simple errors** (`errors.New`): For basic, static error messages
2. **Formatted errors** (`fmt.Errorf`): When you need dynamic context
3. **Custom error types**: When errors carry structured data

**Error Wrapping (Go 1.13+):**

The `%w` verb in `fmt.Errorf` creates a chain of errors, preserving the original error while adding context. This enables `errors.Is` and `errors.As` to inspect the error chain.

```
┌─────────────────────────────────────────────────┐
│              Error Chain                         │
│                                                  │
│  "service error" ──wraps──▶ "db error" ──wraps──▶ "connection refused"
│       (outer)                (middle)              (root cause)
└─────────────────────────────────────────────────┘
```

```go
import (
    "errors"
    "fmt"
)

// Method 1: errors.New
func validateAge(age int) error {
    if age < 0 {
        return errors.New("age cannot be negative")
    }
    return nil
}

// Method 2: fmt.Errorf (with formatting)
func validateUsername(username string) error {
    if len(username) < 3 {
        return fmt.Errorf("username '%s' is too short (min 3 chars)", username)
    }
    return nil
}

// Method 3: fmt.Errorf with %w (error wrapping - Go 1.13+)
func fetchUser(id int) error {
    err := database.Query(id)
    if err != nil {
        return fmt.Errorf("fetchUser failed for id %d: %w", id, err)
    }
    return nil
}
```

### 🔹 Custom Error Types

**When to Use Custom Error Types:**

- Error needs to carry additional data (codes, fields, metadata)
- Callers need to programmatically handle different error conditions
- You want type-safe error handling with `errors.As`
- Building public APIs where error structure is part of the contract

**Design Principles:**
1. Implement the `error` interface with a pointer receiver
2. Include fields that help diagnose/handle the error
3. Provide a clear `Error()` message format
4. Consider adding helper methods for common checks

```go
// Custom error type
type ValidationError struct {
    Field   string
    Message string
    Code    int
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("[%d] %s: %s", e.Code, e.Field, e.Message)
}

// Usage
func validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return &ValidationError{
            Field:   "email",
            Message: "invalid email format",
            Code:    1001,
        }
    }
    return nil
}
```

### 🔹 Error Handling Patterns

```go
// Pattern 1: Check and return early
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    // Process file...
    return nil
}

// Pattern 2: errors.Is (check specific error)
func handleError(err error) {
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File does not exist")
        return
    }
    if errors.Is(err, os.ErrPermission) {
        fmt.Println("Permission denied")
        return
    }
    fmt.Println("Unknown error:", err)
}

// Pattern 3: errors.As (extract custom error type)
func handleValidationError(err error) {
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation failed - Field: %s, Code: %d\n", 
            validationErr.Field, validationErr.Code)
        return
    }
    fmt.Println("Other error:", err)
}
```

### 🔹 Sentinel Errors

**Theory: Named Error Values**

Sentinel errors are predefined error values that represent specific, well-known error conditions. They're called "sentinel" because they stand guard for specific failure modes.

**When to Use Sentinel Errors:**
- Error condition is well-defined and doesn't need dynamic data
- Multiple call sites need to check for the same error
- Part of your package's public API

**Naming Convention:**
Sentinel errors typically start with `Err` prefix: `ErrNotFound`, `ErrTimeout`, `ErrInvalidInput`.

**Standard Library Examples:**
- `io.EOF` - End of file/stream
- `sql.ErrNoRows` - No rows returned from query
- `context.Canceled` - Context was canceled
- `context.DeadlineExceeded` - Context deadline passed

```go
// Define package-level sentinel errors
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized access")
    ErrInvalidInput = errors.New("invalid input")
)

func GetUser(id int) (*User, error) {
    user := db.Find(id)
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}

// Caller can check:
user, err := GetUser(123)
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}
```

---

## 4. Context Package

### Understanding Context: The Theory

**What is Context?**

Context is Go's solution for managing request-scoped data, deadlines, and cancellation signals across API boundaries and goroutines. It's essential for building robust, resource-aware services.

**The Problem Context Solves:**

In a microservices architecture:
```
Client Request
     │
     ▼
┌──────────────┐
│   Service A    │ ──▶ Database Query
└───────┬──────┘
        │
        ▼
┌──────────────┐
│   Service B    │ ──▶ External API Call
└───────┬──────┘
        │
        ▼
┌──────────────┐
│   Service C    │ ──▶ Cache Lookup
└──────────────┘
```

If the client disconnects after making a request, all those operations should stop. Without context, you'd waste resources completing work nobody cares about.

**Context Hierarchy:**

Contexts form a tree structure. When a parent context is canceled, all children are automatically canceled.

```
          Background (root)
              │
     ┌───────┴───────┐
     ▼               ▼
  Request 1      Request 2
     │               │
  ┌──┴──┐         ┌──┴──┐
  ▼     ▼         ▼     ▼
 DB   Cache    DB   External
```

**The Four Context Functions:**

| Function | Purpose | Cancellation |
|----------|---------|-------------|
| `Background()` | Root context, never cancels | Never |
| `TODO()` | Placeholder during refactoring | Never |
| `WithCancel()` | Manual cancellation control | Call cancel() |
| `WithTimeout()` | Auto-cancel after duration | After timeout |
| `WithDeadline()` | Auto-cancel at specific time | At deadline |
| `WithValue()` | Carry request-scoped data | Inherits parent |

**Context is Immutable:**

You never modify a context - you derive new contexts from existing ones. Each `With*` function returns a new context.

The `context` package is crucial for managing request lifecycles, cancellation, and passing request-scoped values.

### 🔹 Why Context Matters for Goa

- **Request cancellation**: Client disconnects → cancel ongoing work
- **Timeouts**: Prevent operations from running forever
- **Request tracing**: Pass trace IDs through call chain
- **Deadline propagation**: Cascade timeouts through services

### 🔹 Creating Contexts

**Understanding Context Creation:**

Always start with a root context (`Background` or `TODO`) and derive child contexts as needed. Never pass `nil` as context.

**The Cancellation Contract:**

When you call `WithCancel`, `WithTimeout`, or `WithDeadline`, you receive a `cancel` function. This MUST be called to release resources, even if the context completes normally.

**Why `defer cancel()` is Critical:**

```go
// Without defer cancel() - RESOURCE LEAK!
ctx, _ := context.WithTimeout(parent, 5*time.Second)
// Goroutines and timers continue running!

// Correct - always defer cancel
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel() // Resources cleaned up
```

```go
import (
    "context"
    "time"
)

func main() {
    // Background context (root context)
    ctx := context.Background()

    // TODO context (placeholder during development)
    ctx = context.TODO()

    // With cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Always defer cancel to prevent leaks

    // With timeout
    ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // With deadline
    deadline := time.Now().Add(10 * time.Second)
    ctx, cancel = context.WithDeadline(context.Background(), deadline)
    defer cancel()

    // With value (use sparingly!)
    ctx = context.WithValue(ctx, "requestID", "abc-123")
}
```

### 🔹 Using Context in Functions

```go
// HTTP handler with context
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    result, err := processWithContext(ctx)
    if err != nil {
        if errors.Is(err, context.Canceled) {
            // Client disconnected
            return
        }
        if errors.Is(err, context.DeadlineExceeded) {
            http.Error(w, "Request timeout", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(result)
}

// Function respecting context cancellation
func processWithContext(ctx context.Context) (string, error) {
    // Check context before expensive operation
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    default:
    }

    // Simulate work with context awareness
    for i := 0; i < 100; i++ {
        select {
        case <-ctx.Done():
            return "", ctx.Err()
        default:
            // Do work
            time.Sleep(10 * time.Millisecond)
        }
    }
    
    return "completed", nil
}
```

### 🔹 Context with External Services

```go
func fetchFromAPI(ctx context.Context, url string) ([]byte, error) {
    // Create request with context
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

func queryDatabase(ctx context.Context, query string) (*sql.Rows, error) {
    // Database query with context
    return db.QueryContext(ctx, query)
}
```

### 🔹 Context Best Practices

```go
// ✅ DO: Pass context as first parameter
func DoSomething(ctx context.Context, arg string) error { ... }

// ❌ DON'T: Store context in struct
type BadService struct {
    ctx context.Context // Never do this!
}

// ✅ DO: Use custom type for context keys
type contextKey string
const requestIDKey contextKey = "requestID"

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
```

---

## 5. Concurrency (Goroutines & Channels)

### Understanding Go Concurrency: The Theory

**Concurrency vs Parallelism:**

- **Concurrency**: Dealing with many things at once (structure)
- **Parallelism**: Doing many things at once (execution)

Go is designed for concurrency. Parallelism is an optional execution optimization.

```
           Concurrency                    Parallelism
    ┌───────┬───────┬───────┐    ┌───────┐ ┌───────┐
    │ Task1 │ Task2 │ Task3 │    │ Task1 │ │ Task2 │
    └───┬───┴───┬───┴───┬───┘    └───┬───┘ └───┬───┘
        │       │       │            │         │
    ─────────────────────▶    ─────▶   ─────▶
         Single Core               Core 1    Core 2
        (interleaved)              (parallel)
```

**Go's Concurrency Philosophy:**

> "Don't communicate by sharing memory; share memory by communicating."

This means prefer channels over shared variables protected by mutexes.

**CSP (Communicating Sequential Processes):**

Go's concurrency model is based on CSP, a formal language from 1978. The key idea: independent processes communicate by sending messages through channels.

**Why Concurrency Matters for Goa:**
- Handle multiple HTTP requests simultaneously
- Make concurrent calls to multiple services
- Process background tasks
- Implement timeouts and cancellations

### 🔹 Goroutines

**What is a Goroutine?**

A goroutine is a lightweight thread of execution managed by the Go runtime. Unlike OS threads (which typically use 1-2MB stack), goroutines start with a tiny stack (~2KB) that grows as needed.

**Goroutine Scheduling:**

```
┌────────────────────────────────────────┐
│            Go Runtime Scheduler           │
└────────────────────┬───────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   ┌────────┐   ┌────────┐   ┌────────┐
   │ OS      │   │ OS      │   │ OS      │
   │ Thread 1│   │ Thread 2│   │ Thread 3│
   └────┬───┘   └────┬───┘   └────┬───┘
       │            │            │
   ┌───┴───┐   ┌───┴───┐   ┌───┴───┐
   │G1  G4 │   │G2  G5 │   │G3  G6 │
   └───────┘   └───────┘   └───────┘
   (Goroutines multiplexed onto OS threads)
```

**Goroutine Characteristics:**
- Extremely cheap to create (thousands are normal)
- Automatically multiplexed onto OS threads
- Cooperative scheduling at safe points
- Stack grows and shrinks dynamically

**Common Pitfall - Loop Variable Capture:**

```go
// WRONG - All goroutines see the same value!
for _, task := range tasks {
    go func() {
        process(task)  // Captures loop variable
    }()
}

// CORRECT - Pass as parameter
for _, task := range tasks {
    go func(t string) {
        process(t)  // Fresh copy each iteration
    }(task)
}
```

Goroutines are lightweight threads managed by Go runtime.

```go
func main() {
    // Start a goroutine
    go func() {
        fmt.Println("Hello from goroutine!")
    }()

    // Named function as goroutine
    go processTask("task-1")

    // Wait for goroutines (simple approach)
    time.Sleep(time.Second) // Don't use in production!
}

func processTask(id string) {
    fmt.Printf("Processing %s\n", id)
}
```

### 🔹 WaitGroup

**Theory: Synchronization Primitive**

`sync.WaitGroup` is a counter-based synchronization mechanism. It waits for a collection of goroutines to finish.

**How It Works:**
1. `Add(n)` - Increment counter by n
2. `Done()` - Decrement counter by 1
3. `Wait()` - Block until counter reaches 0

**Important Rules:**
- Call `Add()` BEFORE launching goroutines (not inside them)
- `Done()` must be called exactly once per goroutine
- `Wait()` blocks until all `Done()` calls complete

```go
import "sync"

func main() {
    var wg sync.WaitGroup
    
    tasks := []string{"task-1", "task-2", "task-3"}
    
    for _, task := range tasks {
        wg.Add(1)
        go func(t string) {
            defer wg.Done()
            processTask(t)
        }(task)
    }
    
    wg.Wait() // Block until all goroutines complete
    fmt.Println("All tasks completed")
}
```

### 🔹 Channels

**Theory: The Communication Primitive**

Channels are typed conduits through which goroutines communicate. They're Go's primary synchronization mechanism.

**Unbuffered vs Buffered Channels:**

```
Unbuffered Channel            Buffered Channel (capacity 3)
┌──────────────────┐          ┌──────────────────┐
│                  │          │ [1] [2] [3]      │
│  Sender blocks   │          │ Send doesn't     │
│  until receiver  │          │ block until full │
│  is ready        │          │                  │
└──────────────────┘          └──────────────────┘
     Synchronous                   Asynchronous
```

**Channel Operations:**

| Operation | Unbuffered | Buffered |
|-----------|------------|----------|
| Send to full | Blocks forever | Blocks until space |
| Send to closed | Panic! | Panic! |
| Receive from empty | Blocks forever | Blocks until data |
| Receive from closed | Returns zero value + false | Drains remaining, then zero |
| Close twice | Panic! | Panic! |

**Directional Channels:**
```go
func worker(jobs <-chan int, results chan<- int) {
    // jobs: receive-only (can't send or close)
    // results: send-only (can't receive)
}
```

Channels are typed conduits for communication between goroutines.

```go
// Unbuffered channel
ch := make(chan int)

// Buffered channel
ch := make(chan int, 10)

// Send to channel
ch <- 42

// Receive from channel
value := <-ch

// Close channel
close(ch)
```

### 🔹 Channel Patterns

**Common Concurrency Patterns:**

Go's channels enable several powerful patterns used extensively in production systems:

| Pattern | Use Case |
|---------|---------|
| Worker Pool | Process jobs with fixed number of workers |
| Fan-Out/Fan-In | Distribute work, merge results |
| Pipeline | Chain processing stages |
| Semaphore | Limit concurrent operations |
| Pub/Sub | Broadcast to multiple consumers |

#### Worker Pool Pattern

**Theory:**

The worker pool pattern limits concurrency to prevent resource exhaustion. Instead of spawning unlimited goroutines, you create a fixed pool of workers that process jobs from a shared channel.

**When to Use:**
- Rate-limiting API calls
- Database connection limits
- CPU-bound tasks (workers = num CPUs)
- Memory-constrained environments

```go
func workerPool() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start workers
    numWorkers := 3
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for a := 1; a <= 9; a++ {
        <-results
    }
}

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}
```

#### Fan-Out, Fan-In Pattern

```go
func fanOutFanIn() {
    input := make(chan int)
    
    // Fan-out: distribute work to multiple goroutines
    c1 := process(input)
    c2 := process(input)
    c3 := process(input)
    
    // Fan-in: merge results
    output := merge(c1, c2, c3)
    
    // Send input
    go func() {
        for i := 0; i < 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // Read merged output
    for result := range output {
        fmt.Println(result)
    }
}

func merge(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }
    
    wg.Add(len(channels))
    for _, c := range channels {
        go output(c)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### 🔹 Select Statement

**Theory: Multiplexing Channel Operations**

The `select` statement is Go's way of waiting on multiple channel operations simultaneously. It's like a switch statement, but for channels.

**How Select Works:**
1. All cases are evaluated simultaneously
2. If multiple cases are ready, one is chosen randomly
3. If no case is ready, it blocks (unless there's a `default`)
4. `default` makes select non-blocking

**Select Patterns:**

| Pattern | Code | Purpose |
|---------|------|---------|
| Timeout | `case <-time.After(d):` | Prevent infinite waits |
| Cancellation | `case <-ctx.Done():` | Respond to context cancel |
| Non-blocking | `default:` | Try without blocking |
| Multiple sources | Multiple `case` | First-ready wins |

**Why Select is Essential for Goa:**
- Implement request timeouts
- Handle graceful shutdown
- Coordinate multiple service calls
- Implement circuit breakers

```go
func selectExample(ctx context.Context) {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from channel 1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from channel 2"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        case <-ctx.Done():
            fmt.Println("Context cancelled")
            return
        case <-time.After(3 * time.Second):
            fmt.Println("Timeout!")
            return
        }
    }
}
```

### 🔹 Mutex for Shared State

**Theory: When Channels Aren't Enough**

While Go prefers "share memory by communicating," sometimes shared state is the right choice. Mutexes protect shared resources from concurrent access.

**When to Use Mutex vs Channels:**

| Use Mutex | Use Channels |
|-----------|-------------|
| Protecting shared state | Passing data ownership |
| Simple counter/flag | Coordinating goroutines |
| Cache with many readers | Pipeline processing |
| Performance-critical hot paths | Signal completion |

**Types of Locks:**

| Lock Type | Use Case |
|-----------|----------|
| `sync.Mutex` | Exclusive access (read OR write) |
| `sync.RWMutex` | Multiple readers OR single writer |

**RWMutex Performance:**
```
Workload: 90% reads, 10% writes

Mutex:    [===R===][===R===][===W===][===R===]  (Serial)
RWMutex:  [R][R][R][R][===W===][R][R][R][R]     (Parallel reads)
```

**Deadlock Prevention Rules:**
1. Always `Lock()` before `Unlock()`
2. Use `defer Unlock()` to ensure unlock
3. Lock in consistent order across goroutines
4. Avoid nested locks when possible

```go
import "sync"

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.RLock() // Use RWMutex for read-heavy workloads
    defer c.mu.RUnlock()
    return c.count
}

// Using sync.RWMutex for read-heavy scenarios
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}
```

---

## 6. Modules & Dependency Management

### Understanding Go Modules: The Theory

**What is a Module?**

A module is a collection of Go packages versioned together as a single unit. It's Go's solution for:
- **Dependency management**: Track and version external code
- **Reproducible builds**: Ensure consistent builds across environments
- **Semantic versioning**: Communicate compatibility guarantees

**Evolution of Go Dependency Management:**

```
2009-2017: $GOPATH (no versioning, global workspace)
    ↓
2017-2018: dep, glide, godep (community tools)
    ↓
2018-now:  Go Modules (official, integrated)
```

**Module Path:**

The module path is the unique identifier for your module, typically matching your repository URL:
```
module github.com/username/project
```

**Semantic Versioning (SemVer):**

```
vMAJOR.MINOR.PATCH
  │     │     │
  │     │     └─── Bug fixes (backward compatible)
  │     └───────── New features (backward compatible)
  └─────────────── Breaking changes
```

**Major Version in Import Path:**

Go enforces semantic import versioning. v2+ must include version in path:
```go
import "github.com/user/project/v2"  // v2.x.x
import "github.com/user/project/v3"  // v3.x.x
```

This allows different major versions to coexist in the same build.

**MVS (Minimal Version Selection):**

Go uses MVS algorithm - it selects the minimum version that satisfies all requirements. This provides:
- Reproducible builds
- Simple mental model
- No SAT solver complexity

### 🔹 Go Modules Basics

```bash
# Initialize a new module
go mod init github.com/username/project

# Add dependencies (automatically updates go.mod)
go get github.com/gin-gonic/gin

# Add specific version
go get github.com/gin-gonic/gin@v1.9.1

# Update all dependencies
go get -u ./...

# Tidy up (remove unused, add missing)
go mod tidy

# Verify dependencies
go mod verify

# Download dependencies
go mod download
```

### 🔹 go.mod File Structure

**Understanding go.mod:**

The `go.mod` file is the manifest for your module. It declares:
- Module identity (path)
- Go version requirement
- Direct dependencies
- Indirect dependencies (transitive)
- Replacements and exclusions

**Directive Types:**

| Directive | Purpose |
|-----------|--------|
| `module` | Declares module path |
| `go` | Minimum Go version |
| `require` | Dependencies needed |
| `replace` | Override dependency location |
| `exclude` | Block specific versions |
| `retract` | Mark own versions as broken |

```go
module github.com/username/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    goa.design/goa/v3 v3.14.0
    google.golang.org/grpc v1.59.0
)

require (
    // Indirect dependencies (auto-managed)
    github.com/go-playground/validator/v10 v10.15.5 // indirect
    golang.org/x/net v0.17.0 // indirect
)

// Replace for local development
replace github.com/myorg/shared => ../shared

// Exclude problematic version
exclude github.com/some/package v1.2.3
```

### 🔹 go.sum File

**Theory: Cryptographic Verification**

The `go.sum` file contains cryptographic checksums of module content. It ensures:
- **Integrity**: Downloaded modules match expected content
- **Authenticity**: No tampering during download
- **Reproducibility**: Same content across all builds

**Format:**
```
module/path version hash
module/path version/go.mod hash
```

**Should you commit go.sum?**
Yes! It's essential for reproducible builds. It should be version controlled alongside go.mod.

```
github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
```

### 🔹 Workspace Mode (Go 1.18+)

**Theory: Multi-Module Development**

Workspaces solve a common problem: developing multiple related modules simultaneously. Before workspaces, you needed `replace` directives in each go.mod, which were error-prone and shouldn't be committed.

**When to Use Workspaces:**
- Developing a library and a service that uses it
- Microservices that share common packages
- Contributing to open-source while testing locally
- Monorepo with multiple Go modules

**Workspace vs Replace:**

| `replace` in go.mod | `go.work` |
|--------------------|----------|
| Committed (bad) or gitignored | Never committed |
| Per-module configuration | Workspace-wide |
| Affects builds anywhere | Local development only |

```bash
# Create workspace file
go work init

# Add modules to workspace
go work use ./service1 ./service2 ./shared

# go.work file
go 1.21

use (
    ./service1
    ./service2
    ./shared
)
```

### 🔹 Vendoring

**Theory: Self-Contained Builds**

Vendoring copies all dependencies into a `vendor/` directory in your project. This makes your project self-contained, not dependent on external sources.

**When to Vendor:**
- CI/CD without external network access
- Regulatory requirements for auditing dependencies
- Protecting against dependency disappearance
- Offline development environments

**Trade-offs:**

| Pros | Cons |
|------|------|
| Reproducible builds | Repository bloat |
| No network needed | Harder to update deps |
| Full code audit possible | Merge conflicts |
| Protection from supply chain attacks | Large diffs |

```bash
# Create vendor directory
go mod vendor

# Build using vendor
go build -mod=vendor ./...

# Verify vendor matches go.mod
go mod verify
```

### 🔹 Private Modules

**Theory: Handling Private Dependencies**

By default, Go tries to fetch modules from a public proxy (proxy.golang.org). Private modules need special configuration to bypass the proxy and authenticate directly.

**Key Environment Variables:**

| Variable | Purpose |
|----------|--------|
| `GOPRIVATE` | Skip proxy and checksum DB |
| `GOPROXY` | Configure module proxy |
| `GONOPROXY` | Bypass proxy for patterns |
| `GONOSUMDB` | Skip checksum verification |

**Authentication:**

Go uses Git for fetching modules. Configure Git for authentication:
```bash
# HTTPS with credential helper
git config --global credential.helper store

# SSH (recommended for CI)
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

```bash
# Set GOPRIVATE for private repos
go env -w GOPRIVATE=github.com/mycompany/*

# Or in .bashrc/.zshrc
export GOPRIVATE=github.com/mycompany/*,gitlab.internal.com/*
```

### 🔹 Common Module Commands

| Command | Description |
|---------|-------------|
| `go mod init <module>` | Initialize new module |
| `go mod tidy` | Add missing, remove unused deps |
| `go mod download` | Download deps to cache |
| `go mod verify` | Verify deps haven't been modified |
| `go mod graph` | Print module dependency graph |
| `go mod why <module>` | Explain why module is needed |
| `go mod edit -require=<mod>` | Edit go.mod programmatically |
| `go list -m all` | List all dependencies |
| `go list -m -versions <mod>` | List available versions |

---

## � Summary: Key Takeaways

### Structs & Interfaces
- **Structs** group related data; use tags for metadata (JSON, validation)
- **Composition over inheritance**: Embed structs instead of extending classes
- **Interfaces** are implicit - implement methods, automatically satisfy interface
- **Small interfaces** are better: io.Reader has just one method

### Methods & Receivers
- **Value receiver**: Operates on copy, use for read-only or small structs
- **Pointer receiver**: Modifies original, use for mutations or large structs
- **Consistency**: If one method uses pointer, all should
- **Method sets**: Pointer receivers only satisfy interface with pointer type

### Error Handling
- **Errors are values**: Check `err != nil` explicitly
- **Wrap errors**: Use `fmt.Errorf("context: %w", err)` to add context
- **`errors.Is`**: Check for specific errors in chain
- **`errors.As`**: Extract custom error types from chain
- **Sentinel errors**: Package-level errors like `ErrNotFound`

### Context Package
- **Always pass context**: First parameter, named `ctx`
- **Never store in structs**: Pass through function calls
- **Cancel to release resources**: `defer cancel()` after WithCancel/Timeout
- **Values sparingly**: Only request-scoped data (trace IDs, auth)
- **Propagate through call chain**: All functions should accept and respect context

### Concurrency
- **Goroutines are cheap**: Use thousands, they're ~2KB each
- **Communicate via channels**: Don't share memory
- **WaitGroup for coordination**: Add before spawning, Done in defer
- **Select for multiplexing**: Wait on multiple channels
- **Mutex when needed**: For shared state access patterns

### Modules
- **Semantic versioning**: MAJOR.MINOR.PATCH
- **go mod tidy**: Your most used command
- **go.sum must be committed**: For reproducible builds
- **Workspaces for multi-module**: Never commit go.work

---

## �📋 Knowledge Check

Before proceeding to Part 2 (Goa Framework Basics), ensure you can:

- [ ] Define structs with tags and embedded fields
- [ ] Create interfaces and implement them implicitly
- [ ] Choose between value and pointer receivers appropriately
- [ ] Handle errors using `errors.Is`, `errors.As`, and wrapping
- [ ] Use context for cancellation, timeouts, and value propagation
- [ ] Create goroutines and coordinate them with WaitGroup
- [ ] Use channels for communication and select for multiplexing
- [ ] Initialize and manage Go modules

---

## 🔗 Quick Reference Links

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go Modules Reference](https://go.dev/ref/mod)

---

> **Next Up:** Part 2 - Goa Framework Basics (DSL, Code Generation, Service Design)
