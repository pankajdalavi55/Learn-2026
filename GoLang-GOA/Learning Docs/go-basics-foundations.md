# Go Basics (Foundations) - A Beginner's Guide

## Table of Contents
1. [Introduction to Go](#introduction-to-go)
2. [Installing Go and Setting Up Workspace](#installing-go-and-setting-up-workspace)
3. [Go Command-Line Tools](#go-command-line-tools)
4. [Packages and Imports](#packages-and-imports)
5. [The main Package and main() Function](#the-main-package-and-main-function)
6. [Variables](#variables)
7. [Basic Data Types](#basic-data-types)
8. [Constants](#constants)
9. [Control Structures](#control-structures)
10. [Loops in Go](#loops-in-go)
11. [Comments and Documentation](#comments-and-documentation)

---

## Introduction to Go

### What is Go?

Go (also called Golang) is an open-source programming language created by Google in 2007 and released in 2009. It was designed by Robert Griesemer, Rob Pike, and Ken Thompson to address shortcomings in other languages while keeping their strengths.

### Key Features of Go

- **Simple and Easy to Learn**: Go has a clean, minimalist syntax that's easy for beginners
- **Fast Compilation**: Programs compile quickly into machine code
- **Concurrent Programming**: Built-in support for concurrent programming with goroutines
- **Garbage Collection**: Automatic memory management
- **Static Typing**: Type safety with compile-time type checking
- **Cross-Platform**: Write once, compile for multiple platforms
- **Standard Library**: Rich standard library for common tasks

### Why Learn Go?

- Used by major companies: Google, Uber, Docker, Kubernetes, Netflix, Dropbox
- Great for building web servers, microservices, CLI tools, and cloud applications
- Growing job market with competitive salaries
- Active and supportive community

---

## Installing Go and Setting Up Workspace

### Installation Steps

#### Windows
1. Download the installer from [https://go.dev/dl/](https://go.dev/dl/)
2. Run the `.msi` file and follow the installation wizard
3. The installer will add Go to your PATH automatically

#### macOS
```bash
# Using Homebrew
brew install go
```

#### Linux
```bash
# Download and extract
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

### Verifying Installation

```bash
go version
```

**Output:**
```
go version go1.21.0 windows/amd64
```

### Setting Up Your Workspace

Go uses a workspace structure. By default, your workspace is in `$HOME/go` (or `%USERPROFILE%\go` on Windows).

**Workspace Structure:**
```
$HOME/go/
├── bin/      # Compiled executables
├── pkg/      # Package objects
└── src/      # Source files
```

**Setting GOPATH (optional for modern Go):**
```bash
# Add to your shell profile
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

**Note:** Starting with Go 1.11+, you can use Go Modules, which don't require GOPATH for most projects.

---

## Go Command-Line Tools

Go comes with powerful command-line tools to build, run, and test your programs.

### go run

The `go run` command compiles and runs your Go program in one step. It's great for quick testing during development.

**Syntax:**
```bash
go run filename.go
```

**Example:**

Create a file `hello.go`:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Run it:
```bash
go run hello.go
```

**Output:**
```
Hello, World!
```

**Best Practice:** Use `go run` during development for quick iteration and testing.

---

### go build

The `go build` command compiles your Go program into an executable binary file but doesn't run it.

**Syntax:**
```bash
go build filename.go
```

**Example:**

```bash
# Build the program
go build hello.go

# This creates an executable file named 'hello' (or 'hello.exe' on Windows)
# Run the executable
./hello      # Linux/macOS
hello.exe    # Windows
```

**Output:**
```
Hello, World!
```

**Build with Custom Name:**
```bash
go build -o myapp hello.go
```

**Best Practice:** Use `go build` when you want to create a distributable executable or test the build process.

---

### go test

The `go test` command runs tests written for your Go packages. Test files must end with `_test.go`.

**Example:**

Create `math_utils.go`:
```go
package utils

func Add(a, b int) int {
    return a + b
}
```

Create `math_utils_test.go`:
```go
package utils

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
    }
}
```

Run the test:
```bash
go test
```

**Output:**
```
PASS
ok      utils   0.123s
```

**Run tests with verbose output:**
```bash
go test -v
```

**Best Practice:** Write tests for all your packages to ensure code reliability.

---

## Packages and Imports

### What is a Package?

A package is a collection of Go source files in the same directory that are compiled together. Packages help organize and reuse code.

### Package Declaration

Every Go file must start with a package declaration:

```go
package packagename
```

### Importing Packages

Use the `import` keyword to use code from other packages.

**Single Import:**
```go
import "fmt"
```

**Multiple Imports (Recommended):**
```go
import (
    "fmt"
    "math"
    "strings"
)
```

### Example with Multiple Imports

```go
package main

import (
    "fmt"
    "math"
    "strings"
)

func main() {
    // Using fmt package
    fmt.Println("Hello, Go!")
    
    // Using math package
    sqrt := math.Sqrt(16)
    fmt.Println("Square root of 16:", sqrt)
    
    // Using strings package
    upper := strings.ToUpper("hello")
    fmt.Println("Uppercase:", upper)
}
```

**Output:**
```
Hello, Go!
Square root of 16: 4
Uppercase: HELLO
```

### Common Standard Library Packages

- `fmt` - Formatted I/O (printing, scanning)
- `math` - Mathematical functions
- `strings` - String manipulation
- `time` - Time and date functions
- `os` - Operating system functionality
- `io` - Input/output primitives
- `net/http` - HTTP client and server

**Best Practice:** Import only the packages you use. The Go compiler will give an error for unused imports.

---

## The main Package and main() Function

### The main Package

The `main` package is special in Go. It defines an executable program rather than a library.

### The main() Function

The `main()` function is the entry point of your program. When you run a Go program, execution starts from `main()`.

**Basic Structure:**
```go
package main

import "fmt"

func main() {
    // Your code here
    fmt.Println("Program starts here!")
}
```

**Output:**
```
Program starts here!
```

### Complete Example

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Welcome to Go!")
    fmt.Println("Current time:", time.Now())
    
    // Call other functions
    greet("Alice")
    greet("Bob")
}

func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}
```

**Output:**
```
Welcome to Go!
Current time: 2026-01-11 10:30:45.123456 +0000 UTC
Hello, Alice!
Hello, Bob!
```

**Key Points:**
- Every executable Go program must have exactly one `main` package
- The `main()` function takes no arguments and returns no value
- You can define other functions and call them from `main()`

---

## Variables

In Go, variables are explicitly declared and can be declared in multiple ways.

### Using the var Keyword

The `var` keyword declares one or more variables with a type.

**Syntax:**
```go
var name type
var name type = value
```

**Examples:**

```go
package main

import "fmt"

func main() {
    // Declare without initialization (default value)
    var age int
    fmt.Println("Age:", age) // Output: Age: 0
    
    // Declare with initialization
    var name string = "Alice"
    fmt.Println("Name:", name) // Output: Name: Alice
    
    // Type inference
    var country = "USA" // Type inferred as string
    fmt.Println("Country:", country)
    
    // Multiple variable declaration
    var x, y, z int = 1, 2, 3
    fmt.Println("x, y, z:", x, y, z)
    
    // Multiple variables with different types
    var (
        firstName string = "John"
        lastName  string = "Doe"
        userAge   int    = 30
        isActive  bool   = true
    )
    fmt.Println(firstName, lastName, userAge, isActive)
}
```

**Output:**
```
Age: 0
Name: Alice
Country: USA
x, y, z: 1 2 3
John Doe 30 true
```

### Short Declaration Operator (:=)

The `:=` operator provides a shorthand for declaring and initializing variables. It can only be used inside functions.

**Syntax:**
```go
name := value
```

**Examples:**

```go
package main

import "fmt"

func main() {
    // Short declaration
    message := "Hello, Go!"
    fmt.Println(message)
    
    // Type is inferred
    count := 42        // int
    price := 19.99     // float64
    isValid := true    // bool
    
    fmt.Printf("count: %d, price: %.2f, isValid: %t\n", count, price, isValid)
    
    // Multiple variable declaration
    firstName, lastName := "Jane", "Smith"
    fmt.Println("Name:", firstName, lastName)
    
    // Mix with existing variables (at least one must be new)
    age := 25
    age, city := 26, "New York" // age is reassigned, city is new
    fmt.Println("Age:", age, "City:", city)
}
```

**Output:**
```
Hello, Go!
count: 42, price: 19.99, isValid: true
Name: Jane Smith
Age: 26 City: New York
```

### Default Values (Zero Values)

Variables declared without an explicit initial value are given their zero value:

```go
package main

import "fmt"

func main() {
    var i int       // 0
    var f float64   // 0.0
    var b bool      // false
    var s string    // "" (empty string)
    
    fmt.Printf("int: %d, float: %.1f, bool: %t, string: '%s'\n", i, f, b, s)
}
```

**Output:**
```
int: 0, float: 0.0, bool: false, string: ''
```

**Best Practices:**
- Use `:=` for concise code inside functions
- Use `var` when you need to declare a variable without initializing it
- Use `var` for package-level variables (outside functions)
- Choose meaningful variable names

---

## Basic Data Types

Go has several built-in basic data types.

### Integer Types (int)

Integers are whole numbers without decimal points.

```go
package main

import "fmt"

func main() {
    var age int = 25
    year := 2026
    
    fmt.Println("Age:", age)
    fmt.Println("Year:", year)
    
    // Integer arithmetic
    sum := 10 + 5
    difference := 10 - 5
    product := 10 * 5
    quotient := 10 / 3    // Integer division
    remainder := 10 % 3   // Modulus
    
    fmt.Println("Sum:", sum)
    fmt.Println("Difference:", difference)
    fmt.Println("Product:", product)
    fmt.Println("Quotient:", quotient)
    fmt.Println("Remainder:", remainder)
}
```

**Output:**
```
Age: 25
Year: 2026
Sum: 15
Difference: 5
Product: 50
Quotient: 3
Remainder: 1
```

**Integer Types:**
- `int`, `int8`, `int16`, `int32`, `int64` (signed)
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` (unsigned)

---

### Floating-Point Types (float)

Floating-point numbers represent decimal values.

```go
package main

import "fmt"

func main() {
    var price float64 = 19.99
    temperature := 72.5
    
    fmt.Println("Price:", price)
    fmt.Println("Temperature:", temperature)
    
    // Floating-point arithmetic
    total := 10.5 + 5.3
    division := 10.0 / 3.0
    
    fmt.Println("Total:", total)
    fmt.Printf("Division: %.2f\n", division) // 2 decimal places
}
```

**Output:**
```
Price: 19.99
Temperature: 72.5
Total: 15.8
Division: 3.33
```

**Float Types:**
- `float32` - 32-bit floating point
- `float64` - 64-bit floating point (default and recommended)

---

### String Type (string)

Strings represent sequences of characters.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    var greeting string = "Hello"
    name := "Alice"
    
    // String concatenation
    message := greeting + ", " + name + "!"
    fmt.Println(message)
    
    // String length
    fmt.Println("Length:", len(message))
    
    // String manipulation with strings package
    upper := strings.ToUpper(message)
    lower := strings.ToLower(message)
    
    fmt.Println("Upper:", upper)
    fmt.Println("Lower:", lower)
    
    // Multi-line strings
    poem := `Roses are red,
Violets are blue,
Go is awesome,
And so are you!`
    
    fmt.Println(poem)
}
```

**Output:**
```
Hello, Alice!
Length: 13
Upper: HELLO, ALICE!
Lower: hello, alice!
Roses are red,
Violets are blue,
Go is awesome,
And so are you!
```

**String Features:**
- Strings are immutable (cannot be changed after creation)
- Use double quotes `"` for regular strings
- Use backticks `` ` `` for raw/multi-line strings

---

### Boolean Type (bool)

Boolean values represent true or false.

```go
package main

import "fmt"

func main() {
    var isActive bool = true
    isComplete := false
    
    fmt.Println("Active:", isActive)
    fmt.Println("Complete:", isComplete)
    
    // Boolean operations
    age := 25
    hasLicense := true
    
    canDrive := age >= 18 && hasLicense
    fmt.Println("Can drive:", canDrive)
    
    // Comparison operators
    x, y := 10, 20
    fmt.Println("x == y:", x == y)
    fmt.Println("x != y:", x != y)
    fmt.Println("x < y:", x < y)
    fmt.Println("x > y:", x > y)
    
    // Logical operators
    fmt.Println("true && false:", true && false)
    fmt.Println("true || false:", true || false)
    fmt.Println("!true:", !true)
}
```

**Output:**
```
Active: true
Complete: false
Can drive: true
x == y: false
x != y: true
x < y: true
x > y: false
true && false: false
true || false: true
!true: false
```

---

## Constants

Constants are fixed values that cannot be changed during program execution. They are declared using the `const` keyword.

### Declaring Constants

```go
package main

import "fmt"

func main() {
    // Single constant
    const pi = 3.14159
    
    // Typed constant
    const greeting string = "Hello, World!"
    
    // Multiple constants
    const (
        daysInWeek   = 7
        hoursInDay   = 24
        minutesInDay = 24 * 60
    )
    
    fmt.Println("Pi:", pi)
    fmt.Println("Greeting:", greeting)
    fmt.Println("Minutes in a day:", minutesInDay)
}
```

**Output:**
```
Pi: 3.14159
Greeting: Hello, World!
Minutes in a day: 1440
```

### Constants with iota

`iota` is a special keyword for creating enumerated constants.

```go
package main

import "fmt"

func main() {
    const (
        Sunday    = iota // 0
        Monday           // 1
        Tuesday          // 2
        Wednesday        // 3
        Thursday         // 4
        Friday           // 5
        Saturday         // 6
    )
    
    fmt.Println("Sunday:", Sunday)
    fmt.Println("Wednesday:", Wednesday)
    fmt.Println("Saturday:", Saturday)
}
```

**Output:**
```
Sunday: 0
Wednesday: 3
Saturday: 6
```

### Real-World Example

```go
package main

import "fmt"

func main() {
    // API configuration constants
    const (
        apiVersion  = "v1"
        maxRetries  = 3
        timeout     = 30
        apiBaseURL  = "https://api.example.com"
    )
    
    fmt.Printf("API: %s/%s\n", apiBaseURL, apiVersion)
    fmt.Printf("Max Retries: %d, Timeout: %d seconds\n", maxRetries, timeout)
}
```

**Output:**
```
API: https://api.example.com/v1
Max Retries: 3, Timeout: 30 seconds
```

**Best Practices:**
- Use constants for values that should never change
- Use ALL_CAPS or camelCase for constant names
- Group related constants together
- Constants can only be numbers, strings, or booleans

---

## Control Structures

Control structures determine the flow of program execution.

### if Statement

The `if` statement executes code based on a condition.

**Basic if:**
```go
package main

import "fmt"

func main() {
    age := 20
    
    if age >= 18 {
        fmt.Println("You are an adult")
    }
}
```

**Output:**
```
You are an adult
```

**if-else:**
```go
package main

import "fmt"

func main() {
    temperature := 25
    
    if temperature > 30 {
        fmt.Println("It's hot!")
    } else {
        fmt.Println("It's pleasant")
    }
}
```

**Output:**
```
It's pleasant
```

**if-else if-else:**
```go
package main

import "fmt"

func main() {
    score := 85
    
    if score >= 90 {
        fmt.Println("Grade: A")
    } else if score >= 80 {
        fmt.Println("Grade: B")
    } else if score >= 70 {
        fmt.Println("Grade: C")
    } else {
        fmt.Println("Grade: F")
    }
}
```

**Output:**
```
Grade: B
```

**if with Short Statement:**
```go
package main

import "fmt"

func main() {
    // Variable num is scoped to the if block
    if num := 10; num > 5 {
        fmt.Println("Number is greater than 5:", num)
    }
    // num is not accessible here
}
```

**Output:**
```
Number is greater than 5: 10
```

---

### switch Statement

The `switch` statement is a cleaner way to write multiple `if-else` conditions.

**Basic switch:**
```go
package main

import "fmt"

func main() {
    day := 3
    
    switch day {
    case 1:
        fmt.Println("Monday")
    case 2:
        fmt.Println("Tuesday")
    case 3:
        fmt.Println("Wednesday")
    case 4:
        fmt.Println("Thursday")
    case 5:
        fmt.Println("Friday")
    default:
        fmt.Println("Weekend")
    }
}
```

**Output:**
```
Wednesday
```

**Multiple values in case:**
```go
package main

import "fmt"

func main() {
    day := "Saturday"
    
    switch day {
    case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
        fmt.Println("It's a weekday")
    case "Saturday", "Sunday":
        fmt.Println("It's the weekend!")
    default:
        fmt.Println("Invalid day")
    }
}
```

**Output:**
```
It's the weekend!
```

**switch without expression (like if-else):**
```go
package main

import "fmt"

func main() {
    temperature := 28
    
    switch {
    case temperature < 0:
        fmt.Println("Freezing")
    case temperature < 15:
        fmt.Println("Cold")
    case temperature < 25:
        fmt.Println("Pleasant")
    default:
        fmt.Println("Warm")
    }
}
```

**Output:**
```
Warm
```

**switch with short statement:**
```go
package main

import "fmt"

func main() {
    switch hour := 14; {
    case hour < 12:
        fmt.Println("Good morning!")
    case hour < 18:
        fmt.Println("Good afternoon!")
    default:
        fmt.Println("Good evening!")
    }
}
```

**Output:**
```
Good afternoon!
```

**Key Points:**
- Go's `switch` doesn't require `break` statements (no fall-through by default)
- Use `fallthrough` keyword if you want fall-through behavior
- Can switch on any type, not just integers

---

## Loops in Go

Unlike most languages, Go has only one looping construct: the `for` loop. However, it's versatile enough to handle all looping scenarios.

### Basic for Loop

```go
package main

import "fmt"

func main() {
    // Classic for loop with init, condition, post
    for i := 1; i <= 5; i++ {
        fmt.Println("Iteration:", i)
    }
}
```

**Output:**
```
Iteration: 1
Iteration: 2
Iteration: 3
Iteration: 4
Iteration: 5
```

### for as while Loop

```go
package main

import "fmt"

func main() {
    count := 1
    
    // Only condition (like while loop)
    for count <= 5 {
        fmt.Println("Count:", count)
        count++
    }
}
```

**Output:**
```
Count: 1
Count: 2
Count: 3
Count: 4
Count: 5
```

### Infinite Loop

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    counter := 0
    
    for {
        counter++
        fmt.Println("Loop iteration:", counter)
        
        if counter >= 3 {
            break // Exit the loop
        }
        
        time.Sleep(1 * time.Second)
    }
    
    fmt.Println("Loop ended")
}
```

**Output:**
```
Loop iteration: 1
Loop iteration: 2
Loop iteration: 3
Loop ended
```

### for with range

The `range` keyword iterates over arrays, slices, maps, and strings.

**Iterate over string:**
```go
package main

import "fmt"

func main() {
    word := "Go"
    
    for index, char := range word {
        fmt.Printf("Index: %d, Character: %c\n", index, char)
    }
}
```

**Output:**
```
Index: 0, Character: G
Index: 1, Character: o
```

**Iterate over array/slice:**
```go
package main

import "fmt"

func main() {
    numbers := []int{10, 20, 30, 40, 50}
    
    for index, value := range numbers {
        fmt.Printf("Index: %d, Value: %d\n", index, value)
    }
    
    fmt.Println("\nValues only:")
    for _, value := range numbers {
        fmt.Println(value)
    }
}
```

**Output:**
```
Index: 0, Value: 10
Index: 1, Value: 20
Index: 2, Value: 30
Index: 3, Value: 40
Index: 4, Value: 50

Values only:
10
20
30
40
50
```

### continue and break

**continue** - Skip to next iteration:
```go
package main

import "fmt"

func main() {
    // Print only odd numbers
    for i := 1; i <= 10; i++ {
        if i%2 == 0 {
            continue // Skip even numbers
        }
        fmt.Println(i)
    }
}
```

**Output:**
```
1
3
5
7
9
```

**break** - Exit loop:
```go
package main

import "fmt"

func main() {
    // Find first number divisible by 7
    for i := 1; i <= 100; i++ {
        if i%7 == 0 {
            fmt.Println("First number divisible by 7:", i)
            break
        }
    }
}
```

**Output:**
```
First number divisible by 7: 7
```

### Nested Loops

```go
package main

import "fmt"

func main() {
    // Multiplication table
    for i := 1; i <= 3; i++ {
        for j := 1; j <= 3; j++ {
            fmt.Printf("%d x %d = %d\n", i, j, i*j)
        }
        fmt.Println("---")
    }
}
```

**Output:**

```
1 x 1 = 1
1 x 2 = 2
1 x 3 = 3
---
2 x 1 = 2
2 x 2 = 4
2 x 3 = 6
---
3 x 1 = 3
3 x 2 = 6
3 x 3 = 9
---
```

**Best Practices:**
- Use `for i := 0; i < n; i++` for counter-based loops
- Use `for condition` when you need while-like behavior
- Use `range` when iterating over collections
- Use `_` (blank identifier) to ignore values you don't need

---

## Comments and Documentation

Comments help explain code and make it more readable. Good documentation is essential for maintainable code.

### Single-Line Comments

Use `//` for single-line comments.

```go
package main

import "fmt"

func main() {
    // This is a single-line comment
    fmt.Println("Hello, World!")
    
    age := 25 // Comment after code
    
    // Comments can explain complex logic
    // Multiple single-line comments
    fmt.Println("Age:", age)
}
```

**Output:**
```
Hello, World!
Age: 25
```

### Multi-Line Comments

Use `/* */` for multi-line comments.

```go
package main

import "fmt"

/*
This is a multi-line comment.
It can span multiple lines.
Useful for longer explanations.
*/

func main() {
    /*
        Calculate total price
        with tax and discount
    */
    price := 100.0
    tax := price * 0.08
    total := price + tax
    
    fmt.Printf("Total: $%.2f\n", total)
}
```

**Output:**
```
Total: $108.00
```

### Documentation Comments

Documentation comments (also called doc comments) describe packages, functions, types, and constants. They appear immediately before the declaration, with no blank line.

**Package Documentation:**
```go
// Package calculator provides basic arithmetic operations
// for performing mathematical calculations.
package calculator

import "fmt"

// Add returns the sum of two integers.
// It takes two parameters a and b and returns their sum.
func Add(a, b int) int {
    return a + b
}

// Multiply returns the product of two integers.
func Multiply(a, b int) int {
    return a * b
}

func main() {
    result := Add(5, 3)
    fmt.Println("5 + 3 =", result)
    
    product := Multiply(4, 7)
    fmt.Println("4 * 7 =", product)
}
```

**Output:**
```
5 + 3 = 8
4 * 7 = 28
```

### Using godoc

`godoc` is a tool that extracts and displays documentation from your Go source code.

**Example with proper documentation:**

Create a file `mathutil.go`:
```go
// Package mathutil provides utility functions for mathematical operations.
package mathutil

// Pi represents the mathematical constant π
const Pi = 3.14159

// Square returns the square of a number.
// It multiplies the number by itself.
//
// Example:
//   result := Square(5)  // result = 25
func Square(x float64) float64 {
    return x * x
}

// CircleArea calculates the area of a circle given its radius.
// The formula used is: A = π * r²
//
// Parameters:
//   radius - the radius of the circle
//
// Returns:
//   The area of the circle
func CircleArea(radius float64) float64 {
    return Pi * Square(radius)
}
```

**View documentation:**
```bash
# Install godoc (if not already installed)
go install golang.org/x/tools/cmd/godoc@latest

# Run godoc server
godoc -http=:6060

# Open in browser: http://localhost:6060/pkg/yourpackage/
```

Alternatively, use `go doc` command:
```bash
go doc mathutil
go doc mathutil.Square
go doc mathutil.CircleArea
```

### Documentation Best Practices

```go
package main

import "fmt"

// User represents a user in the system.
// Each user has a unique ID, name, and email address.
type User struct {
    ID    int    // Unique identifier
    Name  string // Full name of the user
    Email string // Email address
}

// NewUser creates and returns a new User instance.
// It initializes the user with the provided values.
//
// Example:
//   user := NewUser(1, "John Doe", "john@example.com")
func NewUser(id int, name, email string) *User {
    return &User{
        ID:    id,
        Name:  name,
        Email: email,
    }
}

// Greet prints a greeting message for the user.
func (u *User) Greet() {
    fmt.Printf("Hello, %s!\n", u.Name)
}

func main() {
    user := NewUser(1, "Alice", "alice@example.com")
    user.Greet()
}
```

**Output:**
```
Hello, Alice!
```

### Comment Guidelines

**DO:**
- Write comments that explain *why*, not *what*
- Keep comments up-to-date with code changes
- Use complete sentences with proper punctuation
- Document all exported (public) functions, types, and constants
- Start doc comments with the name of the thing being described

**DON'T:**
- Don't state the obvious
- Don't leave outdated comments
- Don't use comments to explain bad code - refactor instead
- Don't comment out code for long periods - use version control

**Example of Good vs. Bad Comments:**

```go
package main

import "fmt"

func main() {
    // BAD: Increment i by 1
    i := 0
    i++
    
    // GOOD: Adjust for zero-based indexing
    index := userInput - 1
    
    // BAD: Loop from 0 to 9
    for i := 0; i < 10; i++ {
        fmt.Println(i)
    }
    
    // GOOD: Process only the first 10 records to limit API calls
    for i := 0; i < 10; i++ {
        fmt.Println(i)
    }
}
```

---

## Summary

Congratulations! You've learned the fundamental building blocks of Go programming:

✓ **Introduction to Go** - Understanding what Go is and why it's popular  
✓ **Installation** - Setting up Go on your system  
✓ **Command-line tools** - Using `go run`, `go build`, and `go test`  
✓ **Packages** - Organizing and importing code  
✓ **Variables** - Declaring and using variables with `var` and `:=`  
✓ **Data types** - Working with `int`, `float64`, `string`, and `bool`  
✓ **Constants** - Defining immutable values  
✓ **Control structures** - Making decisions with `if` and `switch`  
✓ **Loops** - Iterating with the versatile `for` loop  
✓ **Comments** - Documenting your code effectively  

### Next Steps

Now that you have a solid foundation, consider exploring:

- **Functions** - Creating reusable code blocks
- **Arrays and Slices** - Working with collections
- **Maps** - Key-value data structures
- **Structs** - Custom data types
- **Pointers** - Memory addresses and references
- **Methods** - Functions with receivers
- **Interfaces** - Defining behavior contracts
- **Error Handling** - Managing errors gracefully
- **Goroutines and Channels** - Concurrent programming
- **Testing** - Writing unit tests

### Practice Projects

Build these projects to reinforce your learning:

1. **Calculator** - Command-line calculator with basic operations
2. **Todo List** - Simple task management program
3. **Temperature Converter** - Convert between Celsius and Fahrenheit
4. **Number Guessing Game** - Interactive game with user input
5. **Grade Calculator** - Calculate average grades and assign letter grades

### Resources

- **Official Documentation**: [https://go.dev/doc/](https://go.dev/doc/)
- **Go Tour**: [https://go.dev/tour/](https://go.dev/tour/)
- **Go by Example**: [https://gobyexample.com/](https://gobyexample.com/)
- **Effective Go**: [https://go.dev/doc/effective_go](https://go.dev/doc/effective_go)

Happy coding! 🚀
