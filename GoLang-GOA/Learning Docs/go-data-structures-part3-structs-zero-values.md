# Go Core Data Structures - Part 3: Structs and Zero Values

## Table of Contents
1. [Structs](#structs)
   - [Introduction and Basics](#structs-introduction)
   - [Syntax and Declaration](#structs-syntax)
   - [Field Visibility](#structs-visibility)
   - [Struct Literals](#structs-literals)
   - [Anonymous Structs](#structs-anonymous)
   - [Embedded Structs (Composition)](#structs-embedding)
   - [Methods on Structs](#structs-methods)
   - [Pointer vs Value Receivers](#structs-receivers)
   - [Struct Comparison](#structs-comparison)
   - [Struct Tags](#structs-tags)
   - [Code Examples](#structs-examples)
   - [Memory Behavior](#structs-memory)
   - [Common Mistakes](#structs-mistakes)
   - [Interview Tips](#structs-interview)
2. [Zero Values](#zero-values)
   - [Introduction](#zero-introduction)
   - [Zero Values by Type](#zero-by-type)
   - [Zero Values in Practice](#zero-practice)
   - [Designing with Zero Values](#zero-design)
   - [Code Examples](#zero-examples)
   - [Common Patterns](#zero-patterns)
   - [Interview Tips](#zero-interview)

---

<a name="structs"></a>
## 1. Structs

<a name="structs-introduction"></a>
### Introduction and Basics

A **struct** is a composite data type that groups together variables (fields) under a single name. Structs are the primary way to create custom types in Go.

**Key Characteristics:**
- **Value type** - Structs are passed by value (copied when assigned or passed to functions)
- **Type-safe** - Each field has a specific type
- **Flexible** - Can contain any types including other structs
- **No inheritance** - Go uses composition instead of inheritance
- **Efficient** - Stored contiguously in memory

**Real-world analogy:** Like a form or record - a person has name, age, email, etc.

---

<a name="structs-syntax"></a>
### Syntax and Declaration

```go
// Define a struct type
type StructName struct {
    Field1 Type1
    Field2 Type2
    Field3 Type3
}

// Create an instance
var instance StructName
```

#### Example: Basic Struct Definition

```go
package main

import "fmt"

// Define a Person struct
type Person struct {
    Name    string
    Age     int
    Email   string
    IsAdmin bool
}

func main() {
    // Method 1: Zero value struct
    var p1 Person
    fmt.Printf("p1: %+v\n", p1)
    
    // Method 2: Struct literal with field names
    p2 := Person{
        Name:    "Alice",
        Age:     30,
        Email:   "alice@example.com",
        IsAdmin: false,
    }
    fmt.Printf("p2: %+v\n", p2)
    
    // Method 3: Struct literal without field names (order matters)
    p3 := Person{"Bob", 25, "bob@example.com", true}
    fmt.Printf("p3: %+v\n", p3)
    
    // Method 4: Partial initialization
    p4 := Person{Name: "Charlie", Age: 35}
    fmt.Printf("p4: %+v\n", p4)
}
```

**Output:**
```
p1: {Name: Age:0 Email: IsAdmin:false}
p2: {Name:Alice Age:30 Email:alice@example.com IsAdmin:false}
p3: {Name:Bob Age:25 Email:bob@example.com IsAdmin:true}
p4: {Name:Charlie Age:35 Email: IsAdmin:false}
```

---

#### Accessing and Modifying Fields

```go
package main

import "fmt"

type Book struct {
    Title  string
    Author string
    Pages  int
    Price  float64
}

func main() {
    book := Book{
        Title:  "The Go Programming Language",
        Author: "Alan Donovan",
        Pages:  380,
        Price:  44.99,
    }
    
    // Access fields
    fmt.Println("Title:", book.Title)
    fmt.Println("Author:", book.Author)
    
    // Modify fields
    book.Price = 39.99
    book.Pages = 400
    
    fmt.Printf("Updated: %+v\n", book)
}
```

**Output:**
```
Title: The Go Programming Language
Author: Alan Donovan
Updated: {Title:The Go Programming Language Author:Alan Donovan Pages:400 Price:39.99}
```

---

<a name="structs-visibility"></a>
### Field Visibility (Exported vs Unexported)

In Go, visibility is controlled by capitalization:
- **Capitalized** fields are **exported** (public, accessible from other packages)
- **Lowercase** fields are **unexported** (private, accessible only within the same package)

```go
package main

import "fmt"

type User struct {
    // Exported fields (public)
    ID       int
    Username string
    Email    string
    
    // Unexported fields (private)
    password       string
    sessionToken   string
    failedAttempts int
}

func NewUser(id int, username, email, password string) *User {
    return &User{
        ID:             id,
        Username:       username,
        Email:          email,
        password:       password,
        sessionToken:   "",
        failedAttempts: 0,
    }
}

// Method to access private field (getter)
func (u *User) CheckPassword(pwd string) bool {
    if pwd == u.password {
        return true
    }
    u.failedAttempts++
    return false
}

// Method to update private field (setter)
func (u *User) ChangePassword(oldPwd, newPwd string) bool {
    if u.password == oldPwd {
        u.password = newPwd
        return true
    }
    return false
}

func main() {
    user := NewUser(1, "alice", "alice@example.com", "secret123")
    
    // Access exported fields
    fmt.Println("Username:", user.Username)
    fmt.Println("Email:", user.Email)
    
    // Cannot access unexported fields from outside the package
    // fmt.Println(user.password)  // Compile error if in different package
    
    // Use methods to interact with private fields
    fmt.Println("Password correct:", user.CheckPassword("secret123"))
    fmt.Println("Password correct:", user.CheckPassword("wrong"))
    
    // Change password
    if user.ChangePassword("secret123", "newSecret456") {
        fmt.Println("Password changed successfully")
    }
}
```

**Output:**
```
Username: alice
Email: alice@example.com
Password correct: true
Password correct: false
Password changed successfully
```

**Best Practice:** 
- Export only what needs to be public
- Use methods (getters/setters) for controlled access to private fields
- This provides encapsulation and data protection

---

<a name="structs-literals"></a>
### Struct Literals

There are multiple ways to initialize structs.

```go
package main

import "fmt"

type Point struct {
    X, Y int
}

type Rectangle struct {
    Width, Height float64
    Color         string
}

func main() {
    // 1. Zero value (all fields have zero values)
    var p1 Point
    fmt.Printf("p1: %+v\n", p1)
    
    // 2. Named fields (recommended, order doesn't matter)
    p2 := Point{X: 10, Y: 20}
    fmt.Printf("p2: %+v\n", p2)
    
    // 3. Positional (order matters, must include all fields)
    p3 := Point{30, 40}
    fmt.Printf("p3: %+v\n", p3)
    
    // 4. Partial initialization (unspecified fields get zero values)
    r1 := Rectangle{Width: 10.5}
    fmt.Printf("r1: %+v\n", r1)
    
    // 5. Pointer to struct using &
    p4 := &Point{X: 50, Y: 60}
    fmt.Printf("p4: %+v (type: %T)\n", p4, p4)
    
    // 6. Using new (returns pointer, all fields zero)
    p5 := new(Point)
    fmt.Printf("p5: %+v (type: %T)\n", p5, p5)
}
```

**Output:**
```
p1: {X:0 Y:0}
p2: {X:10 Y:20}
p3: {X:30 Y:40}
r1: {Width:10.5 Height:0 Color:}
p4: &{X:50 Y:60} (type: *main.Point)
p5: &{X:0 Y:0} (type: *main.Point)
```

**Best Practice:** Use named field initialization for clarity and maintainability, especially when structs have many fields.

---

<a name="structs-anonymous"></a>
### Anonymous Structs

Anonymous structs are structs without a named type. Useful for one-off data structures.

```go
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    // Anonymous struct variable
    person := struct {
        Name string
        Age  int
    }{
        Name: "Alice",
        Age:  30,
    }
    
    fmt.Printf("person: %+v\n", person)
    
    // Anonymous struct in slice
    people := []struct {
        Name string
        Age  int
    }{
        {"Bob", 25},
        {"Charlie", 35},
        {"Diana", 28},
    }
    
    fmt.Println("\nPeople:")
    for _, p := range people {
        fmt.Printf("  %s is %d years old\n", p.Name, p.Age)
    }
    
    // Common use case: JSON marshaling
    response := struct {
        Status  int    `json:"status"`
        Message string `json:"message"`
        Data    struct {
            UserID   int    `json:"user_id"`
            Username string `json:"username"`
        } `json:"data"`
    }{
        Status:  200,
        Message: "Success",
        Data: struct {
            UserID   int    `json:"user_id"`
            Username string `json:"username"`
        }{
            UserID:   123,
            Username: "alice",
        },
    }
    
    jsonData, _ := json.MarshalIndent(response, "", "  ")
    fmt.Printf("\nJSON Response:\n%s\n", jsonData)
}
```

**Output:**
```
person: {Name:Alice Age:30}

People:
  Bob is 25 years old
  Charlie is 35 years old
  Diana is 28 years old

JSON Response:
{
  "status": 200,
  "message": "Success",
  "data": {
    "user_id": 123,
    "username": "alice"
  }
}
```

**Use Cases:**
- Table-driven tests
- JSON response structures
- Temporary data grouping
- One-time use data structures

---

<a name="structs-embedding"></a>
### Embedded Structs (Composition Over Inheritance)

Go doesn't have inheritance, but it has **composition** through struct embedding.

#### Basic Embedding

```go
package main

import "fmt"

type Address struct {
    Street  string
    City    string
    Country string
}

type Person struct {
    Name    string
    Age     int
    Address Address  // Nested struct
}

func main() {
    person := Person{
        Name: "Alice",
        Age:  30,
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
        },
    }
    
    // Access nested fields
    fmt.Println("Name:", person.Name)
    fmt.Println("City:", person.Address.City)
    fmt.Println("Street:", person.Address.Street)
}
```

**Output:**
```
Name: Alice
City: New York
Street: 123 Main St
```

---

#### Anonymous Embedding (Field Promotion)

```go
package main

import "fmt"

type Address struct {
    Street  string
    City    string
    Country string
}

type Person struct {
    Name    string
    Age     int
    Address  // Anonymous field (embedded)
}

func main() {
    person := Person{
        Name: "Bob",
        Age:  25,
        Address: Address{
            Street:  "456 Oak Ave",
            City:    "Boston",
            Country: "USA",
        },
    }
    
    // Fields are "promoted" - can access directly
    fmt.Println("Name:", person.Name)
    fmt.Println("City:", person.City)      // Promoted from Address
    fmt.Println("Street:", person.Street)  // Promoted from Address
    
    // Can still access via Address
    fmt.Println("Country:", person.Address.Country)
}
```

**Output:**
```
Name: Bob
City: Boston
Street: 456 Oak Ave
Country: USA
```

**Key Point:** Embedded struct fields are "promoted" to the outer struct, allowing direct access.

---

#### Embedding with Methods (Composition Pattern)

```go
package main

import "fmt"

// Base type with methods
type Engine struct {
    Horsepower int
    Type       string
}

func (e *Engine) Start() {
    fmt.Printf("Engine started: %d HP %s\n", e.Horsepower, e.Type)
}

func (e *Engine) Stop() {
    fmt.Println("Engine stopped")
}

// Car embeds Engine
type Car struct {
    Brand  string
    Model  string
    Engine  // Embedded - inherits Engine's methods
}

// Car can have its own methods too
func (c *Car) Drive() {
    fmt.Printf("Driving %s %s\n", c.Brand, c.Model)
}

func main() {
    car := Car{
        Brand: "Toyota",
        Model: "Camry",
        Engine: Engine{
            Horsepower: 200,
            Type:       "V6",
        },
    }
    
    // Call embedded Engine methods directly
    car.Start()
    
    // Call Car's own method
    car.Drive()
    
    // Stop engine
    car.Stop()
    
    // Access embedded fields
    fmt.Printf("Horsepower: %d\n", car.Horsepower)
}
```

**Output:**
```
Engine started: 200 HP V6
Driving Toyota Camry
Engine stopped
Horsepower: 200
```

**This is Go's approach to "inheritance":**
- Composition over inheritance
- Embedded struct's methods are promoted
- More flexible than traditional inheritance

---

<a name="structs-methods"></a>
### Methods on Structs

Methods are functions with a receiver argument.

```go
package main

import (
    "fmt"
    "math"
)

type Circle struct {
    Radius float64
}

// Method with value receiver
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Method with value receiver
func (c Circle) Circumference() float64 {
    return 2 * math.Pi * c.Radius
}

// Method with pointer receiver (can modify struct)
func (c *Circle) Scale(factor float64) {
    c.Radius *= factor
}

// Method with value receiver (cannot modify original)
func (c Circle) ScaleValue(factor float64) {
    c.Radius *= factor  // Modifies copy, not original
}

func main() {
    circle := Circle{Radius: 5.0}
    
    fmt.Printf("Circle radius: %.2f\n", circle.Radius)
    fmt.Printf("Area: %.2f\n", circle.Area())
    fmt.Printf("Circumference: %.2f\n", circle.Circumference())
    
    // Modify using pointer receiver
    circle.Scale(2.0)
    fmt.Printf("After Scale(2.0), radius: %.2f\n", circle.Radius)
    
    // Try to modify using value receiver
    circle.ScaleValue(2.0)
    fmt.Printf("After ScaleValue(2.0), radius: %.2f (unchanged!)\n", circle.Radius)
}
```

**Output:**
```
Circle radius: 5.00
Area: 78.54
Circumference: 31.42
After Scale(2.0), radius: 10.00
After ScaleValue(2.0), radius: 10.00 (unchanged!)
```

---

<a name="structs-receivers"></a>
### Pointer vs Value Receivers

**Critical Decision:** When to use pointer receivers vs value receivers?

```go
package main

import "fmt"

type Counter struct {
    Count int
}

// Value receiver - receives a COPY
func (c Counter) IncrementValue() {
    c.Count++
    fmt.Printf("  Inside IncrementValue: %d\n", c.Count)
}

// Pointer receiver - receives a REFERENCE
func (c *Counter) IncrementPointer() {
    c.Count++
    fmt.Printf("  Inside IncrementPointer: %d\n", c.Count)
}

func main() {
    counter := Counter{Count: 0}
    
    fmt.Println("Initial count:", counter.Count)
    
    // Value receiver - doesn't modify original
    fmt.Println("\nCalling IncrementValue:")
    counter.IncrementValue()
    fmt.Println("After IncrementValue, count:", counter.Count)
    
    // Pointer receiver - modifies original
    fmt.Println("\nCalling IncrementPointer:")
    counter.IncrementPointer()
    fmt.Println("After IncrementPointer, count:", counter.Count)
}
```

**Output:**
```
Initial count: 0

Calling IncrementValue:
  Inside IncrementValue: 1
After IncrementValue, count: 0

Calling IncrementPointer:
  Inside IncrementPointer: 1
After IncrementPointer, count: 1
```

---

#### When to Use Pointer vs Value Receivers

**Use Pointer Receivers When:**

✅ Method needs to modify the receiver  
✅ Struct is large (avoid copying overhead)  
✅ Receiver contains fields that shouldn't be copied (sync.Mutex, etc.)  
✅ Consistency - if some methods use pointer receivers, use for all  

**Use Value Receivers When:**

✅ Method doesn't modify the receiver  
✅ Struct is small (< 32 bytes as a rule of thumb)  
✅ Receiver is a value type (int, string, etc.) or small immutable struct  
✅ You want to work with copies for safety  

```go
package main

import "fmt"

// Small struct - value receiver is fine
type Point struct {
    X, Y int
}

func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// Large struct - pointer receiver is better
type Image struct {
    Pixels [1000000]byte
    Width  int
    Height int
}

func (img *Image) Resize(newWidth, newHeight int) {
    img.Width = newWidth
    img.Height = newHeight
    // ... resize logic
}

func main() {
    p := Point{X: 10, Y: 20}
    fmt.Println(p.String())
    
    img := &Image{Width: 1920, Height: 1080}
    img.Resize(1280, 720)
    fmt.Printf("Image: %dx%d\n", img.Width, img.Height)
}
```

**Output:**
```
(10, 20)
Image: 1280x720
```

---

<a name="structs-comparison"></a>
### Struct Comparison

Structs can be compared using `==` and `!=` if all their fields are comparable.

```go
package main

import "fmt"

type Point struct {
    X, Y int
}

type Person struct {
    Name string
    Age  int
}

type Container struct {
    Data []int  // Slice is not comparable
}

func main() {
    // Comparing structs with comparable fields
    p1 := Point{X: 10, Y: 20}
    p2 := Point{X: 10, Y: 20}
    p3 := Point{X: 15, Y: 25}
    
    fmt.Println("p1 == p2:", p1 == p2)  // true
    fmt.Println("p1 == p3:", p1 == p3)  // false
    
    // Comparing more complex structs
    person1 := Person{Name: "Alice", Age: 30}
    person2 := Person{Name: "Alice", Age: 30}
    person3 := Person{Name: "Bob", Age: 25}
    
    fmt.Println("person1 == person2:", person1 == person2)  // true
    fmt.Println("person1 == person3:", person1 == person3)  // false
    
    // Structs with non-comparable fields cannot be compared
    c1 := Container{Data: []int{1, 2, 3}}
    c2 := Container{Data: []int{1, 2, 3}}
    
    // This would NOT compile:
    // fmt.Println(c1 == c2)  // invalid operation: c1 == c2 (struct containing []int cannot be compared)
    
    fmt.Printf("c1: %+v\n", c1)
    fmt.Printf("c2: %+v\n", c2)
}
```

**Output:**
```
p1 == p2: true
p1 == p3: false
person1 == person2: true
person1 == person3: false
c1: {Data:[1 2 3]}
c2: {Data:[1 2 3]}
```

**Comparable Field Types:**
- int, float, string, bool, pointers, arrays (with comparable elements), structs (with comparable fields)

**Non-Comparable Field Types:**
- slices, maps, functions

---

<a name="structs-tags"></a>
### Struct Tags

Struct tags are metadata attached to struct fields, commonly used for JSON, XML, database mappings, etc.

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"-"`                    // Never include in JSON
    FullName  string `json:"full_name,omitempty"`  // Omit if empty
    IsActive  bool   `json:"is_active"`
}

func main() {
    user := User{
        ID:       123,
        Username: "alice",
        Email:    "alice@example.com",
        Password: "secret123",
        FullName: "",
        IsActive: true,
    }
    
    // Marshal to JSON
    jsonData, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("JSON output:")
    fmt.Println(string(jsonData))
    
    // Unmarshal from JSON
    jsonInput := `{
        "id": 456,
        "username": "bob",
        "email": "bob@example.com",
        "is_active": false,
        "full_name": "Bob Smith"
    }`
    
    var newUser User
    json.Unmarshal([]byte(jsonInput), &newUser)
    fmt.Printf("\nUnmarshaled: %+v\n", newUser)
}
```

**Output:**
```
JSON output:
{
  "id": 123,
  "username": "alice",
  "email": "alice@example.com",
  "is_active": true
}

Unmarshaled: {ID:456 Username:bob Email:bob@example.com Password: FullName:Bob Smith IsActive:false}
```

**Common Tag Options:**
- `json:"field_name"` - Custom JSON field name
- `json:"-"` - Skip field in JSON
- `json:",omitempty"` - Omit if zero value
- `xml:"field_name"` - XML marshaling
- `db:"column_name"` - Database column mapping
- `validate:"required"` - Validation rules

---

<a name="structs-examples"></a>
### More Code Examples

#### Example: Struct with Multiple Methods

```go
package main

import (
    "fmt"
    "strings"
)

type BankAccount struct {
    AccountNumber string
    HolderName    string
    Balance       float64
}

func NewBankAccount(number, name string, initialBalance float64) *BankAccount {
    return &BankAccount{
        AccountNumber: number,
        HolderName:    name,
        Balance:       initialBalance,
    }
}

func (ba *BankAccount) Deposit(amount float64) {
    if amount > 0 {
        ba.Balance += amount
        fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, ba.Balance)
    }
}

func (ba *BankAccount) Withdraw(amount float64) bool {
    if amount > 0 && amount <= ba.Balance {
        ba.Balance -= amount
        fmt.Printf("Withdrew $%.2f. New balance: $%.2f\n", amount, ba.Balance)
        return true
    }
    fmt.Println("Insufficient funds")
    return false
}

func (ba *BankAccount) GetBalance() float64 {
    return ba.Balance
}

func (ba *BankAccount) String() string {
    return fmt.Sprintf("Account: %s, Holder: %s, Balance: $%.2f",
        ba.AccountNumber, ba.HolderName, ba.Balance)
}

func main() {
    account := NewBankAccount("ACC001", "Alice Johnson", 1000.00)
    
    fmt.Println(account)
    
    account.Deposit(500.00)
    account.Withdraw(200.00)
    account.Withdraw(2000.00)  // Insufficient funds
    
    fmt.Printf("\nFinal: %s\n", account)
}
```

**Output:**
```
Account: ACC001, Holder: Alice Johnson, Balance: $1000.00
Deposited $500.00. New balance: $1500.00
Withdrew $200.00. New balance: $1300.00
Insufficient funds

Final: Account: ACC001, Holder: Alice Johnson, Balance: $1300.00
```

---

#### Example: Struct Composition in Real-World Scenario

```go
package main

import (
    "fmt"
    "time"
)

type Timestamp struct {
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Author struct {
    Name  string
    Email string
}

type Article struct {
    ID      int
    Title   string
    Content string
    Author  Author
    Timestamp
}

func NewArticle(id int, title, content string, author Author) *Article {
    now := time.Now()
    return &Article{
        ID:      id,
        Title:   title,
        Content: content,
        Author:  author,
        Timestamp: Timestamp{
            CreatedAt: now,
            UpdatedAt: now,
        },
    }
}

func (a *Article) Update(newContent string) {
    a.Content = newContent
    a.UpdatedAt = time.Now()
}

func (a *Article) Display() {
    fmt.Printf("Article #%d: %s\n", a.ID, a.Title)
    fmt.Printf("By: %s (%s)\n", a.Author.Name, a.Author.Email)
    fmt.Printf("Content: %s\n", a.Content)
    fmt.Printf("Created: %s\n", a.CreatedAt.Format("2006-01-02 15:04:05"))
    fmt.Printf("Updated: %s\n", a.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func main() {
    author := Author{
        Name:  "Alice Johnson",
        Email: "alice@example.com",
    }
    
    article := NewArticle(1, "Introduction to Go", "Go is a great language...", author)
    
    article.Display()
    
    fmt.Println("\n--- Updating article ---")
    time.Sleep(2 * time.Second)
    article.Update("Go is an amazing language for building scalable systems!")
    
    fmt.Println()
    article.Display()
}
```

**Output (example):**
```
Article #1: Introduction to Go
By: Alice Johnson (alice@example.com)
Content: Go is a great language...
Created: 2026-01-11 10:30:45
Updated: 2026-01-11 10:30:45

--- Updating article ---

Article #1: Introduction to Go
By: Alice Johnson (alice@example.com)
Content: Go is an amazing language for building scalable systems!
Created: 2026-01-11 10:30:45
Updated: 2026-01-11 10:30:47
```

---

<a name="structs-memory"></a>
### Memory Behavior and Internal Representation

**Struct Memory Layout:**

Structs are stored contiguously in memory, with fields laid out in order (subject to alignment/padding).

```go
package main

import (
    "fmt"
    "unsafe"
)

type Small struct {
    A int8   // 1 byte
    B int8   // 1 byte
    C int8   // 1 byte
}

type WithPadding struct {
    A int8   // 1 byte
    B int64  // 8 bytes (needs alignment)
    C int8   // 1 byte
}

type Optimized struct {
    B int64  // 8 bytes
    A int8   // 1 byte
    C int8   // 1 byte
}

func main() {
    var s Small
    var wp WithPadding
    var opt Optimized
    
    fmt.Printf("Small struct size: %d bytes\n", unsafe.Sizeof(s))
    fmt.Printf("WithPadding struct size: %d bytes (has padding)\n", unsafe.Sizeof(wp))
    fmt.Printf("Optimized struct size: %d bytes (optimized layout)\n", unsafe.Sizeof(opt))
    
    // Field offsets
    fmt.Printf("\nWithPadding field offsets:\n")
    fmt.Printf("  A: %d\n", unsafe.Offsetof(wp.A))
    fmt.Printf("  B: %d\n", unsafe.Offsetof(wp.B))
    fmt.Printf("  C: %d\n", unsafe.Offsetof(wp.C))
}
```

**Output:**
```
Small struct size: 3 bytes
WithPadding struct size: 24 bytes (has padding)
Optimized struct size: 16 bytes (optimized layout)

WithPadding field offsets:
  A: 0
  B: 8
  C: 16
```

**Key Points:**
- Fields are aligned based on their type (int64 needs 8-byte alignment)
- Padding is added between fields for proper alignment
- Order of fields affects total size
- Larger fields first can reduce total size

---

#### Stack vs Heap Allocation

```go
package main

import "fmt"

type Point struct {
    X, Y int
}

func createOnStack() Point {
    // Allocated on stack
    p := Point{X: 10, Y: 20}
    return p
}

func createOnHeap() *Point {
    // Allocated on heap (escapes to heap)
    p := Point{X: 10, Y: 20}
    return &p
}

func main() {
    // Stack allocation
    p1 := createOnStack()
    fmt.Printf("Stack: %+v\n", p1)
    
    // Heap allocation
    p2 := createOnHeap()
    fmt.Printf("Heap: %+v\n", p2)
}
```

**Output:**
```
Stack: {X:10 Y:20}
Heap: &{X:10 Y:20}
```

**Escape Analysis:**
- If struct doesn't escape function, allocated on stack (faster)
- If struct escapes (returned as pointer, stored in global, etc.), allocated on heap
- Compiler automatically determines this

---

<a name="structs-mistakes"></a>
### Common Mistakes and Pitfalls

#### Mistake 1: Forgetting Structs are Value Types

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func modifyPerson(p Person) {
    p.Age = 100
    fmt.Printf("Inside function: %+v\n", p)
}

func modifyPersonPointer(p *Person) {
    p.Age = 100
    fmt.Printf("Inside function: %+v\n", p)
}

func main() {
    person := Person{Name: "Alice", Age: 30}
    
    // Pass by value - doesn't modify original
    modifyPerson(person)
    fmt.Printf("After modifyPerson: %+v\n\n", person)
    
    // Pass by pointer - modifies original
    modifyPersonPointer(&person)
    fmt.Printf("After modifyPersonPointer: %+v\n", person)
}
```

**Output:**
```
Inside function: {Name:Alice Age:100}
After modifyPerson: {Name:Alice Age:30}
Inside function: &{Name:Alice Age:100}
After modifyPersonPointer: {Name:Alice Age:100}
```

---

#### Mistake 2: Comparing Structs with Non-Comparable Fields

```go
package main

import "fmt"

type User struct {
    Name  string
    Roles []string  // Slice is not comparable
}

func main() {
    u1 := User{Name: "Alice", Roles: []string{"admin"}}
    u2 := User{Name: "Alice", Roles: []string{"admin"}}
    
    // This does NOT compile:
    // if u1 == u2 { }  // invalid operation: u1 == u2
    
    // CORRECT: Compare fields manually
    if u1.Name == u2.Name {
        fmt.Println("Same name")
    }
    
    // Or use reflect.DeepEqual
    // if reflect.DeepEqual(u1, u2) { }
}
```

**Output:**
```
Same name
```

---

#### Mistake 3: Not Using Pointer Receivers When Needed

```go
package main

import "fmt"

type Counter struct {
    Count int
}

// WRONG: Value receiver doesn't modify original
func (c Counter) IncrementWrong() {
    c.Count++
}

// CORRECT: Pointer receiver modifies original
func (c *Counter) IncrementCorrect() {
    c.Count++
}

func main() {
    counter := Counter{Count: 0}
    
    counter.IncrementWrong()
    fmt.Println("After IncrementWrong:", counter.Count)  // Still 0
    
    counter.IncrementCorrect()
    fmt.Println("After IncrementCorrect:", counter.Count)  // Now 1
}
```

**Output:**
```
After IncrementWrong: 0
After IncrementCorrect: 1
```

---

#### Mistake 4: Mixing Pointer and Value Receivers

```go
package main

import "fmt"

type Document struct {
    Title   string
    Content string
}

// Some methods with value receivers
func (d Document) GetTitle() string {
    return d.Title
}

// Some methods with pointer receivers
func (d *Document) SetTitle(title string) {
    d.Title = title
}

func main() {
    // This works fine
    doc := Document{Title: "Original", Content: "..."}
    fmt.Println(doc.GetTitle())
    doc.SetTitle("Updated")
    fmt.Println(doc.GetTitle())
}
```

**Output:**
```
Original
Updated
```

**Best Practice:** Be consistent - if you use pointer receivers for some methods, use them for all methods on that type.

---

<a name="structs-interview"></a>
### Interview Tips and Key Takeaways

#### 🎯 Frequently Asked Interview Questions

**Q1: Are structs passed by value or reference in Go?**

**Answer:** Structs are **passed by value** (copied). To pass by reference, use a pointer: `*StructType`.

**Q2: What's the difference between a value receiver and a pointer receiver?**

**Answer:**
- **Value receiver**: `func (s StructType) Method()` - receives a copy, can't modify original
- **Pointer receiver**: `func (s *StructType) Method()` - receives a reference, can modify original

**Q3: Can you compare structs using ==?**

**Answer:** Yes, if **all fields are comparable**. Structs containing slices, maps, or functions cannot be compared with `==`.

**Q4: How does Go implement inheritance?**

**Answer:** Go **doesn't have inheritance**. Instead, it uses **composition** through struct embedding. Embedded structs promote their fields and methods to the outer struct.

**Q5: What are struct tags used for?**

**Answer:** Struct tags provide metadata for fields, commonly used for:
- JSON/XML marshaling (`json:"field_name"`)
- Database mappings (`db:"column_name"`)
- Validation (`validate:"required"`)

**Q6: What's the zero value of a struct?**

**Answer:** A struct with all fields set to their zero values. For example:
```go
type Person struct {
    Name string  // ""
    Age  int     // 0
}
var p Person  // {Name:"", Age:0}
```

---

#### 🔥 Critical Points for Interviews

1. **Structs are Value Types**
   ```go
   s1 := MyStruct{}
   s2 := s1  // s2 is a COPY, not a reference
   ```

2. **Pointer vs Value Receivers**
   ```go
   func (s StructType) Method()   // Value receiver (copy)
   func (s *StructType) Method()  // Pointer receiver (reference)
   ```

3. **Composition Over Inheritance**
   ```go
   type Car struct {
       Engine  // Embedded - promotes Engine's fields/methods
   }
   ```

4. **Field Visibility**
   ```go
   type User struct {
       Name     string  // Exported (public)
       password string  // Unexported (private)
   }
   ```

5. **Struct Comparison**
   ```go
   // OK if all fields are comparable
   p1 := Point{X: 1, Y: 2}
   p2 := Point{X: 1, Y: 2}
   if p1 == p2 { }  // Works
   ```

---

#### ⚠️ Common Interview Pitfalls

| Pitfall | Problem | Solution |
|---------|---------|----------|
| **Value semantics** | Expecting modifications to propagate | Use pointers when needed |
| **Receiver type** | Using value receiver when modification needed | Use pointer receiver |
| **Comparing structs** | Comparing structs with slices/maps | Check all fields are comparable |
| **Embedding confusion** | Thinking it's inheritance | Understand it's composition |
| **Field order** | Not considering memory alignment | Order large fields first |

---

#### 💡 Best Practices for Production Code

✅ **DO:**
- Use named field initialization for clarity: `Person{Name: "Alice", Age: 30}`
- Use pointer receivers when modifying state or struct is large
- Embed structs for composition, not inheritance
- Capitalize fields that should be exported
- Use struct tags for JSON, database mappings
- Group related fields together
- Order fields by size (largest first) for memory efficiency

❌ **DON'T:**
- Don't mix pointer and value receivers unnecessarily
- Don't compare structs with non-comparable fields using `==`
- Don't export fields that should be private
- Don't use positional initialization for structs with many fields
- Don't forget structs are copied when passed to functions

---

<a name="zero-values"></a>
## 2. Zero Values

<a name="zero-introduction"></a>
### Introduction

In Go, variables declared without an explicit initial value are given their **zero value**. This is a key feature that makes Go code safer and more predictable.

**Philosophy:** Zero values should be useful and safe by default.

---

<a name="zero-by-type"></a>
### Zero Values by Type

| Type | Zero Value | Description |
|------|------------|-------------|
| `int`, `int8`, `int16`, `int32`, `int64` | `0` | Zero |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` | Zero |
| `float32`, `float64` | `0.0` | Zero |
| `bool` | `false` | False |
| `string` | `""` | Empty string |
| `pointer` | `nil` | Nil pointer |
| `slice` | `nil` | Nil slice (len=0, cap=0) |
| `map` | `nil` | Nil map |
| `channel` | `nil` | Nil channel |
| `function` | `nil` | Nil function |
| `interface` | `nil` | Nil interface |
| `struct` | Struct with all fields zero | All fields have zero values |
| `array` | Array with all elements zero | All elements have zero values |

---

#### Examples of Zero Values

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    // Numeric types
    var i int
    var f float64
    fmt.Printf("int: %d, float64: %.1f\n", i, f)
    
    // Boolean
    var b bool
    fmt.Printf("bool: %t\n", b)
    
    // String
    var s string
    fmt.Printf("string: '%s' (empty)\n", s)
    
    // Pointer
    var p *int
    fmt.Printf("pointer: %v (nil: %t)\n", p, p == nil)
    
    // Slice
    var slice []int
    fmt.Printf("slice: %v (nil: %t, len: %d, cap: %d)\n", 
        slice, slice == nil, len(slice), cap(slice))
    
    // Map
    var m map[string]int
    fmt.Printf("map: %v (nil: %t, len: %d)\n", 
        m, m == nil, len(m))
    
    // Struct
    var person Person
    fmt.Printf("struct: %+v\n", person)
    
    // Array
    var arr [3]int
    fmt.Printf("array: %v\n", arr)
}
```

**Output:**
```
int: 0, float64: 0.0
bool: false
string: '' (empty)
pointer: <nil> (nil: true)
slice: [] (nil: true, len: 0, cap: 0)
map: map[] (nil: true, len: 0)
struct: {Name: Age:0}
array: [0 0 0]
```

---

<a name="zero-practice"></a>
### Zero Values in Practice

#### Example 1: Zero Values Are Usable

Go's zero values are designed to be immediately usable.

```go
package main

import "fmt"

func main() {
    // Slice: nil slice is usable with append
    var numbers []int
    numbers = append(numbers, 1, 2, 3)
    fmt.Println("Slice:", numbers)
    
    // String: empty string is usable
    var message string
    message += "Hello"
    message += ", "
    message += "World!"
    fmt.Println("String:", message)
    
    // Numeric: zero is a valid number
    var counter int
    counter++
    counter++
    fmt.Println("Counter:", counter)
    
    // Bool: false is a valid boolean
    var flag bool
    if !flag {
        fmt.Println("Flag is false (zero value)")
    }
}
```

**Output:**
```
Slice: [1 2 3]
String: Hello, World!
Counter: 2
Flag is false (zero value)
```

---

#### Example 2: sync.Mutex Zero Value is Ready to Use

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.Mutex  // Zero value is ready to use!
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func main() {
    // Zero value struct with zero-value mutex
    var counter SafeCounter
    
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    fmt.Println("Final count:", counter.Value())
}
```

**Output:**
```
Final count: 1000
```

**Key Point:** `sync.Mutex` zero value is unlocked and ready to use - no initialization needed!

---

#### Example 3: bytes.Buffer Zero Value is Ready to Use

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    // Zero value bytes.Buffer is ready to use
    var buf bytes.Buffer
    
    buf.WriteString("Hello")
    buf.WriteString(", ")
    buf.WriteString("World!")
    
    fmt.Println(buf.String())
}
```

**Output:**
```
Hello, World!
```

---

<a name="zero-design"></a>
### Designing with Zero Values

**Make the zero value useful** - this is a Go idiom.

#### Example: Well-Designed Type with Useful Zero Value

```go
package main

import (
    "fmt"
    "strings"
)

type StringBuilder struct {
    parts []string
}

// Methods work with zero value
func (sb *StringBuilder) WriteString(s string) {
    sb.parts = append(sb.parts, s)
}

func (sb *StringBuilder) String() string {
    return strings.Join(sb.parts, "")
}

func main() {
    // Zero value is immediately usable
    var sb StringBuilder
    
    sb.WriteString("Go")
    sb.WriteString(" is")
    sb.WriteString(" awesome!")
    
    fmt.Println(sb.String())
}
```

**Output:**
```
Go is awesome!
```

**Design principle:** The zero value `StringBuilder` (with nil slice) works correctly because `append` works on nil slices.

---

#### Example: Configuration with Sensible Zero Values

```go
package main

import (
    "fmt"
    "time"
)

type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    MaxConns     int
}

func (c *ServerConfig) ApplyDefaults() {
    // Apply defaults only if zero
    if c.Host == "" {
        c.Host = "localhost"
    }
    if c.Port == 0 {
        c.Port = 8080
    }
    if c.ReadTimeout == 0 {
        c.ReadTimeout = 5 * time.Second
    }
    if c.WriteTimeout == 0 {
        c.WriteTimeout = 10 * time.Second
    }
    if c.MaxConns == 0 {
        c.MaxConns = 100
    }
}

func main() {
    // Zero value config
    var config ServerConfig
    fmt.Printf("Before defaults: %+v\n", config)
    
    config.ApplyDefaults()
    fmt.Printf("After defaults: %+v\n", config)
    
    // Partial config
    config2 := ServerConfig{Host: "example.com", Port: 9000}
    config2.ApplyDefaults()
    fmt.Printf("Partial config: %+v\n", config2)
}
```

**Output:**
```
Before defaults: {Host: Port:0 ReadTimeout:0s WriteTimeout:0s MaxConns:0}
After defaults: {Host:localhost Port:8080 ReadTimeout:5s WriteTimeout:10s MaxConns:100}
Partial config: {Host:example.com Port:9000 ReadTimeout:5s WriteTimeout:10s MaxConns:100}
```

---

<a name="zero-examples"></a>
### More Zero Value Examples

#### Example: Detecting Zero Value

```go
package main

import "fmt"

type User struct {
    ID       int
    Username string
    Email    string
}

func (u User) IsZero() bool {
    return u == User{}
}

func main() {
    var u1 User
    u2 := User{ID: 1, Username: "alice"}
    var u3 User
    
    fmt.Printf("u1.IsZero(): %t\n", u1.IsZero())
    fmt.Printf("u2.IsZero(): %t\n", u2.IsZero())
    fmt.Printf("u3.IsZero(): %t\n", u3.IsZero())
}
```

**Output:**
```
u1.IsZero(): true
u2.IsZero(): false
u3.IsZero(): true
```

---

#### Example: Pointer vs Value Zero Semantics

```go
package main

import "fmt"

type Config struct {
    Enabled bool
    Timeout int
}

func main() {
    // Value type: zero value is usable
    var cfg1 Config
    fmt.Printf("Value type: %+v (Enabled: %t)\n", cfg1, cfg1.Enabled)
    
    // Pointer type: zero value is nil
    var cfg2 *Config
    fmt.Printf("Pointer type: %v (nil: %t)\n", cfg2, cfg2 == nil)
    
    // Must initialize pointer before use
    cfg2 = &Config{}
    fmt.Printf("After init: %+v (Enabled: %t)\n", cfg2, cfg2.Enabled)
}
```

**Output:**
```
Value type: {Enabled:false Timeout:0} (Enabled: false)
Pointer type: <nil> (nil: true)
After init: &{Enabled:false Timeout:0} (Enabled: false)
```

---

<a name="zero-patterns"></a>
### Common Patterns with Zero Values

#### Pattern 1: Optional Fields with Pointers

```go
package main

import "fmt"

type Article struct {
    Title   string
    Content string
    Views   *int  // nil means "not set", 0 means "zero views"
}

func main() {
    // Article without view count
    a1 := Article{Title: "Go Basics", Content: "..."}
    fmt.Printf("a1.Views: %v (not set: %t)\n", a1.Views, a1.Views == nil)
    
    // Article with zero views
    zeroViews := 0
    a2 := Article{Title: "Advanced Go", Content: "...", Views: &zeroViews}
    fmt.Printf("a2.Views: %v (not set: %t, value: %d)\n", 
        a2.Views, a2.Views == nil, *a2.Views)
}
```

**Output:**
```
a1.Views: <nil> (not set: true)
a2.Views: 0xc000012098 (not set: false, value: 0)
```

**Pattern:** Use pointers to distinguish "not set" (nil) from "set to zero value".

---

#### Pattern 2: Builder Pattern with Zero Value

```go
package main

import "fmt"

type HTTPClient struct {
    baseURL string
    timeout int
    retries int
}

type HTTPClientBuilder struct {
    client HTTPClient
}

func (b *HTTPClientBuilder) BaseURL(url string) *HTTPClientBuilder {
    b.client.baseURL = url
    return b
}

func (b *HTTPClientBuilder) Timeout(t int) *HTTPClientBuilder {
    b.client.timeout = t
    return b
}

func (b *HTTPClientBuilder) Retries(r int) *HTTPClientBuilder {
    b.client.retries = r
    return b
}

func (b *HTTPClientBuilder) Build() HTTPClient {
    // Apply defaults for zero values
    if b.client.timeout == 0 {
        b.client.timeout = 30
    }
    if b.client.retries == 0 {
        b.client.retries = 3
    }
    return b.client
}

func main() {
    // Zero value builder is ready to use
    var builder HTTPClientBuilder
    
    client := builder.
        BaseURL("https://api.example.com").
        Timeout(60).
        Build()
    
    fmt.Printf("Client: %+v\n", client)
}
```

**Output:**
```
Client: {baseURL:https://api.example.com timeout:60 retries:3}
```

---

<a name="zero-interview"></a>
### Interview Tips and Key Takeaways

#### 🎯 Frequently Asked Interview Questions

**Q1: What is the zero value in Go?**

**Answer:** The zero value is the default value assigned to a variable when it's declared without explicit initialization. Every type has a well-defined zero value:
- Numbers: `0`
- Booleans: `false`
- Strings: `""`
- Pointers, slices, maps, channels, functions, interfaces: `nil`
- Structs: struct with all fields set to their zero values

**Q2: Are zero values safe to use in Go?**

**Answer:** Yes! Go's zero values are designed to be **immediately usable and safe**. For example:
- Nil slices can be appended to
- Empty strings can be concatenated
- Zero-value `sync.Mutex` is ready to use
- Zero numeric values are valid numbers

**Q3: What's the difference between nil slice and empty slice?**

**Answer:**
```go
var nilSlice []int        // nil slice (nil, len=0, cap=0)
emptySlice := []int{}     // empty slice (not nil, len=0, cap=0)
```
Both behave the same in most cases, but `nilSlice == nil` is `true`, while `emptySlice == nil` is `false`.

**Q4: Can you write to a nil map?**

**Answer:** **No!** Writing to a nil map causes a runtime panic. You must initialize with `make()` or a map literal first:
```go
var m map[string]int    // nil map
// m["key"] = 1         // PANIC!
m = make(map[string]int)
m["key"] = 1            // OK
```

**Q5: How do you check if a struct is zero value?**

**Answer:** Compare with an empty struct literal:
```go
if myStruct == (MyStructType{}) {
    // It's zero value
}
```

---

#### 🔥 Critical Points for Interviews

1. **Zero Values Are Usable**
   ```go
   var s []int     // nil slice, but can append
   s = append(s, 1)  // Works!
   ```

2. **Nil Map vs Nil Slice**
   ```go
   var slice []int  // Can append, read (safe)
   var m map[string]int  // Can read, but NOT write (panic)
   ```

3. **Struct Zero Value**
   ```go
   var p Person  // All fields have zero values
   // Name: "", Age: 0, etc.
   ```

4. **Pointer Zero Value**
   ```go
   var p *Person  // nil
   // Must initialize before dereferencing
   ```

5. **Design with Zero Value in Mind**
   ```go
   type Buffer struct {
       data []byte  // nil slice is fine
   }
   // Methods should work with zero value
   ```

---

#### 💡 Best Practices

✅ **DO:**
- Design types so their zero value is useful
- Use nil slices instead of empty slices for return values
- Check for nil before dereferencing pointers
- Use zero values for sensible defaults
- Document when zero value is not usable

❌ **DON'T:**
- Don't write to nil maps
- Don't dereference nil pointers
- Don't assume zero value is always safe for custom types
- Don't compare non-comparable types

---

## Summary - Part 3

### Key Takeaways

**Structs:**
- ✓ Value types (copied when passed/assigned)
- ✓ Use pointer receivers to modify state
- ✓ Composition over inheritance (embedding)
- ✓ Exported (capitalized) vs unexported (lowercase) fields
- ✓ Can be compared if all fields are comparable
- ✓ Struct tags for metadata (JSON, validation, etc.)
- ✓ Zero value is struct with all fields zero-valued
- ✓ Memory layout considers field alignment

**Zero Values:**
- ✓ Every type has a well-defined zero value
- ✓ Zero values are designed to be usable and safe
- ✓ Nil slices can be appended to (safe)
- ✓ Nil maps cannot be written to (panic)
- ✓ Design APIs to work with zero values
- ✓ Use pointers to distinguish "not set" from "zero"
- ✓ Zero value struct has all fields at zero values

### Interview Readiness Checklist

**Structs:**
- [ ] Understand value vs reference semantics
- [ ] Know when to use pointer vs value receivers
- [ ] Understand composition through embedding
- [ ] Know field visibility rules
- [ ] Can explain struct comparison rules
- [ ] Understand struct tags and common uses
- [ ] Know memory layout and alignment basics
- [ ] Can design types with useful zero values

**Zero Values:**
- [ ] Know zero value for each type
- [ ] Understand which zero values are safe to use
- [ ] Know nil slice vs empty slice vs nil map
- [ ] Can design APIs using zero values effectively
- [ ] Understand pointer zero value semantics
- [ ] Know how to check for zero value structs

---

## Complete Series Summary

You've now covered all three parts of Go's core data structures:

**Part 1: Arrays and Slices**
- Fixed-size vs dynamic collections
- Value types vs reference-like behavior
- Capacity management and append behavior

**Part 2: Maps and Strings**
- Unordered key-value collections
- UTF-8 encoding and runes
- Reference types and concurrency

**Part 3: Structs and Zero Values**
- Custom composite types
- Composition patterns
- Zero value design philosophy

**You're now ready for Go interviews! 🚀**
