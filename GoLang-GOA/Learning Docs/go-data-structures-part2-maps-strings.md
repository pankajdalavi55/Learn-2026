# Go Core Data Structures - Part 2: Maps and Strings

## Table of Contents
1. [Maps](#maps)
   - [Introduction and Basics](#maps-introduction)
   - [Syntax and Declaration](#maps-syntax)
   - [Basic Operations](#maps-operations)
   - [Code Examples](#maps-examples)
   - [Reference Behavior](#maps-reference)
   - [Zero Value of Maps](#maps-zero-value)
   - [Checking Key Existence](#maps-existence)
   - [Iteration Order](#maps-iteration)
   - [Concurrency Considerations](#maps-concurrency)
   - [Memory Behavior](#maps-memory)
   - [Common Mistakes](#maps-mistakes)
   - [Interview Tips](#maps-interview)
2. [Strings and Runes](#strings)
   - [Introduction and Basics](#strings-introduction)
   - [UTF-8 Encoding](#strings-utf8)
   - [Bytes vs Runes](#strings-bytes-runes)
   - [String Declaration and Operations](#strings-operations)
   - [Iterating Over Strings](#strings-iteration)
   - [String Immutability](#strings-immutability)
   - [Code Examples](#strings-examples)
   - [Common String Operations](#strings-common-ops)
   - [Performance Considerations](#strings-performance)
   - [Common Mistakes](#strings-mistakes)
   - [Interview Tips](#strings-interview)

---

<a name="maps"></a>
## 1. Maps

<a name="maps-introduction"></a>
### Introduction and Basics

A **map** is an unordered collection of key-value pairs. Maps are Go's built-in hash table implementation.

**Key Characteristics:**
- **Unordered** - No guaranteed iteration order
- **Reference type** - Maps are passed by reference
- **Dynamic** - Can grow and shrink
- **Fast lookups** - O(1) average time complexity
- **Unique keys** - Each key appears at most once
- **Type-safe** - Keys and values must match declared types

**Real-world analogy:** Like a dictionary - you look up a word (key) to find its definition (value).

---

<a name="maps-syntax"></a>
### Syntax and Declaration

```go
// Syntax
var mapName map[KeyType]ValueType

// Examples
var ages map[string]int           // Map from string to int
var scores map[int]float64         // Map from int to float64
var cache map[string]interface{}   // Map from string to any type
```

**Key Requirements:**
- Key type must be **comparable** (can use `==` and `!=`)
- Valid key types: int, string, bool, pointers, structs (with comparable fields), arrays
- Invalid key types: slices, maps, functions

---

<a name="maps-operations"></a>
### Basic Operations

```go
// Create map
m := make(map[string]int)

// Add/Update
m["key"] = value

// Read
value := m["key"]

// Delete
delete(m, "key")

// Check existence
value, exists := m["key"]

// Length
len(m)
```

---

<a name="maps-examples"></a>
### Code Examples

#### Example 1: Creating Maps (Different Methods)

```go
package main

import "fmt"

func main() {
    // Method 1: Using make
    ages := make(map[string]int)
    ages["Alice"] = 25
    ages["Bob"] = 30
    fmt.Println("Method 1 (make):", ages)
    
    // Method 2: Map literal (empty)
    scores := map[string]float64{}
    scores["Math"] = 95.5
    scores["Physics"] = 88.0
    fmt.Println("Method 2 (literal):", scores)
    
    // Method 3: Map literal with initialization
    capitals := map[string]string{
        "USA":    "Washington D.C.",
        "France": "Paris",
        "Japan":  "Tokyo",
    }
    fmt.Println("Method 3 (initialized):", capitals)
    
    // Method 4: Nil map (cannot add elements!)
    var nilMap map[string]int
    fmt.Printf("Method 4 (nil map): %v, isNil=%v\n", nilMap, nilMap == nil)
}
```

**Output:**
```
Method 1 (make): map[Alice:25 Bob:30]
Method 2 (literal): map[Math:95.5 Physics:88]
Method 3 (initialized): map[France:Paris Japan:Tokyo USA:Washington D.C.]
Method 4 (nil map): map[], isNil=true
```

---

#### Example 2: Basic Map Operations

```go
package main

import "fmt"

func main() {
    // Create and populate map
    prices := make(map[string]float64)
    
    // Add elements
    prices["Apple"] = 1.50
    prices["Banana"] = 0.75
    prices["Orange"] = 2.00
    
    fmt.Println("Initial prices:", prices)
    
    // Access element
    applePrice := prices["Apple"]
    fmt.Printf("Apple price: $%.2f\n", applePrice)
    
    // Update element
    prices["Apple"] = 1.75
    fmt.Printf("Updated Apple price: $%.2f\n", prices["Apple"])
    
    // Delete element
    delete(prices, "Banana")
    fmt.Println("After deleting Banana:", prices)
    
    // Get map length
    fmt.Println("Number of items:", len(prices))
}
```

**Output:**
```
Initial prices: map[Apple:1.5 Banana:0.75 Orange:2]
Apple price: $1.50
Updated Apple price: $1.75
After deleting Banana: map[Apple:1.75 Orange:2]
Number of items: 2
```

---

#### Example 3: Iterating Over Maps

```go
package main

import "fmt"

func main() {
    grades := map[string]int{
        "Alice":   90,
        "Bob":     85,
        "Charlie": 92,
        "Diana":   88,
    }
    
    // Method 1: Iterate with key and value
    fmt.Println("Method 1: Key and Value")
    for name, grade := range grades {
        fmt.Printf("%s: %d\n", name, grade)
    }
    
    // Method 2: Iterate with keys only
    fmt.Println("\nMethod 2: Keys only")
    for name := range grades {
        fmt.Println(name)
    }
    
    // Method 3: Iterate with values only (using blank identifier)
    fmt.Println("\nMethod 3: Values only")
    for _, grade := range grades {
        fmt.Println(grade)
    }
}
```

**Output (order may vary):**
```
Method 1: Key and Value
Alice: 90
Bob: 85
Charlie: 92
Diana: 88

Method 2: Keys only
Alice
Bob
Charlie
Diana

Method 3: Values only
90
85
92
88
```

**Important:** The iteration order is **random** and can change between runs!

---

<a name="maps-reference"></a>
### Reference Behavior

Maps are reference types - they're passed by reference, not copied.

```go
package main

import "fmt"

func modifyMap(m map[string]int) {
    m["New"] = 100
    m["Alice"] = 999
}

func main() {
    ages := map[string]int{
        "Alice": 25,
        "Bob":   30,
    }
    
    fmt.Println("Before function:", ages)
    
    // Pass map to function
    modifyMap(ages)
    
    fmt.Println("After function:", ages)
    // Original map IS modified!
}
```

**Output:**
```
Before function: map[Alice:25 Bob:30]
After function: map[Alice:999 Bob:30 New:100]
```

**Key Point:** Changes made inside the function affect the original map.

---

#### Creating Independent Copies

```go
package main

import "fmt"

func copyMap(original map[string]int) map[string]int {
    copied := make(map[string]int)
    for key, value := range original {
        copied[key] = value
    }
    return copied
}

func main() {
    original := map[string]int{"A": 1, "B": 2, "C": 3}
    
    // Create independent copy
    copied := copyMap(original)
    
    // Modify copy
    copied["A"] = 999
    copied["D"] = 4
    
    fmt.Println("Original:", original)  // Unchanged
    fmt.Println("Copied:", copied)      // Modified
}
```

**Output:**
```
Original: map[A:1 B:2 C:3]
Copied: map[A:999 B:2 C:3 D:4]
```

---

<a name="maps-zero-value"></a>
### Zero Value of Maps

The zero value of a map is `nil`. A `nil` map behaves like an empty map for reading, but **cannot be written to**.

```go
package main

import "fmt"

func main() {
    var nilMap map[string]int
    
    // Check if nil
    fmt.Printf("nilMap == nil: %v\n", nilMap == nil)
    fmt.Printf("Length: %d\n", len(nilMap))
    
    // Reading from nil map is safe (returns zero value)
    value := nilMap["key"]
    fmt.Printf("Reading from nil map: %d\n", value)
    
    // Check existence is safe
    _, exists := nilMap["key"]
    fmt.Printf("Key exists: %v\n", exists)
    
    // Iterating over nil map is safe (no iterations)
    for k, v := range nilMap {
        fmt.Printf("This won't print: %s=%d\n", k, v)
    }
    
    // Writing to nil map causes PANIC!
    // nilMap["key"] = 10  // panic: assignment to entry in nil map
    
    // Must initialize before writing
    nilMap = make(map[string]int)
    nilMap["key"] = 10
    fmt.Println("After initialization:", nilMap)
}
```

**Output:**
```
nilMap == nil: true
Length: 0
Reading from nil map: 0
Key exists: false
After initialization: map[key:10]
```

**Critical Rule:** Always initialize maps with `make()` or a map literal before adding elements.

---

#### Nil Map vs Empty Map

```go
package main

import "fmt"

func main() {
    // Nil map
    var nilMap map[string]int
    
    // Empty map (not nil)
    emptyMap := make(map[string]int)
    
    // Another empty map
    emptyMap2 := map[string]int{}
    
    fmt.Printf("nilMap == nil: %v, len=%d\n", nilMap == nil, len(nilMap))
    fmt.Printf("emptyMap == nil: %v, len=%d\n", emptyMap == nil, len(emptyMap))
    fmt.Printf("emptyMap2 == nil: %v, len=%d\n", emptyMap2 == nil, len(emptyMap2))
    
    // Can write to empty maps
    emptyMap["key"] = 42
    emptyMap2["key"] = 42
    
    // Cannot write to nil map
    // nilMap["key"] = 42  // PANIC!
}
```

**Output:**
```
nilMap == nil: true, len=0
emptyMap == nil: false, len=0
emptyMap2 == nil: false, len=0
```

---

<a name="maps-existence"></a>
### Checking Key Existence

The "comma ok" idiom is the standard way to check if a key exists in a map.

```go
package main

import "fmt"

func main() {
    ages := map[string]int{
        "Alice": 25,
        "Bob":   30,
        "Carol": 0,  // Note: value is zero
    }
    
    // Method 1: Just get the value (returns zero value if missing)
    age1 := ages["Alice"]
    age2 := ages["Unknown"]
    fmt.Printf("Alice: %d, Unknown: %d\n", age1, age2)
    
    // Method 2: Check existence with "comma ok" idiom
    if age, exists := ages["Alice"]; exists {
        fmt.Printf("Alice exists, age: %d\n", age)
    } else {
        fmt.Println("Alice not found")
    }
    
    if age, exists := ages["Unknown"]; exists {
        fmt.Printf("Unknown exists, age: %d\n", age)
    } else {
        fmt.Println("Unknown not found")
    }
    
    // Important: Distinguish between zero value and missing key
    if age, exists := ages["Carol"]; exists {
        fmt.Printf("Carol exists with age: %d (zero value)\n", age)
    } else {
        fmt.Println("Carol not found")
    }
}
```

**Output:**
```
Alice: 25, Unknown: 0
Alice exists, age: 25
Unknown not found
Carol exists with age: 0 (zero value)
```

**Best Practice:** Always use the "comma ok" idiom when you need to distinguish between a missing key and a key with a zero value.

---

#### Practical Example: Word Counter

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "go is great go is simple go go go"
    words := strings.Fields(text)
    
    // Count word occurrences
    wordCount := make(map[string]int)
    
    for _, word := range words {
        // Increment count (zero value is 0, so this works even for new words)
        wordCount[word]++
    }
    
    fmt.Println("Word counts:")
    for word, count := range wordCount {
        fmt.Printf("%s: %d\n", word, count)
    }
}
```

**Output:**
```
Word counts:
go: 5
is: 2
great: 1
simple: 1
```

---

<a name="maps-iteration"></a>
### Iteration Order

**Critical Fact:** Map iteration order is **not guaranteed** and is **intentionally randomized** in Go.

```go
package main

import "fmt"

func main() {
    m := map[int]string{
        1: "one",
        2: "two",
        3: "three",
        4: "four",
        5: "five",
    }
    
    fmt.Println("Run 1:")
    for k, v := range m {
        fmt.Printf("%d: %s\n", k, v)
    }
    
    fmt.Println("\nRun 2:")
    for k, v := range m {
        fmt.Printf("%d: %s\n", k, v)
    }
    
    // Order may differ between runs!
}
```

**Output (example - will vary):**
```
Run 1:
1: one
2: two
3: three
4: four
5: five

Run 2:
4: four
1: one
5: five
2: two
3: three
```

---

#### Sorted Iteration Over Maps

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    ages := map[string]int{
        "Diana":   28,
        "Alice":   25,
        "Bob":     30,
        "Charlie": 22,
    }
    
    // Get keys and sort them
    keys := make([]string, 0, len(ages))
    for key := range ages {
        keys = append(keys, key)
    }
    sort.Strings(keys)
    
    // Iterate in sorted order
    fmt.Println("Sorted by name:")
    for _, key := range keys {
        fmt.Printf("%s: %d\n", key, ages[key])
    }
}
```

**Output:**
```
Sorted by name:
Alice: 25
Bob: 30
Charlie: 22
Diana: 28
```

---

<a name="maps-concurrency"></a>
### Concurrency Considerations

**⚠️ WARNING:** Maps are **NOT safe for concurrent use**. Concurrent reads and writes will cause a runtime panic.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    m := make(map[int]int)
    
    // UNSAFE: This will panic
    // for i := 0; i < 100; i++ {
    //     go func(n int) {
    //         m[n] = n  // Concurrent writes - PANIC!
    //     }(i)
    // }
    
    // SAFE: Using mutex
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            mu.Lock()
            m[n] = n
            mu.Unlock()
        }(i)
    }
    
    wg.Wait()
    fmt.Printf("Map has %d elements\n", len(m))
}
```

**Output:**
```
Map has 100 elements
```

---

#### Using sync.Map for Concurrent Access

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var m sync.Map
    var wg sync.WaitGroup
    
    // Concurrent writes (safe with sync.Map)
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            m.Store(n, n*n)  // Store key-value
        }(i)
    }
    
    wg.Wait()
    
    // Read values
    if value, ok := m.Load(10); ok {
        fmt.Printf("Key 10: %v\n", value)
    }
    
    // Count elements
    count := 0
    m.Range(func(key, value interface{}) bool {
        count++
        return true  // Continue iteration
    })
    fmt.Printf("Total elements: %d\n", count)
}
```

**Output:**
```
Key 10: 100
Total elements: 100
```

**When to use sync.Map:**
- Multiple goroutines reading and writing
- Mostly reads with occasional writes
- Keys are written once and read many times

**When to use regular map + mutex:**
- Simple concurrent access patterns
- More type safety (sync.Map uses `interface{}`)

---

<a name="maps-memory"></a>
### Memory Behavior and Internal Representation

**Map Internal Structure:**

Maps in Go are implemented as hash tables with separate chaining for collision resolution.

```
Conceptual Structure:

map[string]int
     ↓
Hash Table
┌─────┬─────┬─────┬─────┬─────┐
│  0  │  1  │  2  │  3  │  4  │ ... buckets
└─────┴─────┴─────┴─────┴─────┘
   ↓
Bucket (holds up to 8 key-value pairs)
```

**Memory Characteristics:**

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    // Empty map
    m1 := make(map[string]int)
    fmt.Printf("Empty map size: %d bytes\n", unsafe.Sizeof(m1))
    
    // Map with 1000 elements
    m2 := make(map[string]int, 1000)
    for i := 0; i < 1000; i++ {
        m2[fmt.Sprintf("key%d", i)] = i
    }
    fmt.Printf("Map with 1000 elements size: %d bytes\n", unsafe.Sizeof(m2))
    
    // Note: unsafe.Sizeof returns the size of the map header, not the data
}
```

**Output:**
```
Empty map size: 8 bytes
Map with 1000 elements size: 8 bytes
```

**Key Points:**
- Map variable is a pointer to the map header (8 bytes on 64-bit)
- Actual data is stored on the heap
- Maps automatically grow when needed
- Deleted entries leave "tombstones" (memory not immediately freed)

---

#### Map Pre-allocation

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    n := 1000000
    
    // Without pre-allocation
    start := time.Now()
    m1 := make(map[int]int)
    for i := 0; i < n; i++ {
        m1[i] = i
    }
    fmt.Printf("Without pre-allocation: %v\n", time.Since(start))
    
    // With pre-allocation
    start = time.Now()
    m2 := make(map[int]int, n)
    for i := 0; i < n; i++ {
        m2[i] = i
    }
    fmt.Printf("With pre-allocation: %v\n", time.Since(start))
}
```

**Output (example):**
```
Without pre-allocation: 125ms
With pre-allocation: 85ms
```

**Best Practice:** If you know the approximate size, pre-allocate with `make(map[K]V, size)` for better performance.

---

<a name="maps-mistakes"></a>
### Common Mistakes and Pitfalls

#### Mistake 1: Writing to Nil Map

```go
package main

import "fmt"

func main() {
    var m map[string]int
    
    // This will PANIC!
    // m["key"] = 42  // panic: assignment to entry in nil map
    
    // CORRECT: Initialize first
    m = make(map[string]int)
    m["key"] = 42
    fmt.Println(m)
}
```

**Output:**
```
map[key:42]
```

---

#### Mistake 2: Assuming Map Iteration Order

```go
package main

import "fmt"

func main() {
    m := map[int]string{1: "a", 2: "b", 3: "c"}
    
    // WRONG: Assuming order
    // for k, v := range m {
    //     // Don't assume k will be 1, 2, 3 in that order!
    // }
    
    // CORRECT: Sort keys if order matters
    fmt.Println("Order is not guaranteed!")
}
```

---

#### Mistake 3: Not Checking Key Existence

```go
package main

import "fmt"

func main() {
    scores := map[string]int{
        "Alice": 90,
        "Bob":   0,  // Actual score is 0
    }
    
    // WRONG: Can't distinguish between zero value and missing key
    score1 := scores["Bob"]      // Returns 0 (exists with value 0)
    score2 := scores["Charlie"]  // Returns 0 (doesn't exist)
    fmt.Printf("Bob: %d, Charlie: %d - Can't tell the difference!\n", score1, score2)
    
    // CORRECT: Use comma ok idiom
    if score, exists := scores["Bob"]; exists {
        fmt.Printf("Bob's score: %d (exists)\n", score)
    }
    
    if score, exists := scores["Charlie"]; exists {
        fmt.Printf("Charlie's score: %d\n", score)
    } else {
        fmt.Println("Charlie not found")
    }
}
```

**Output:**
```
Bob: 0, Charlie: 0 - Can't tell the difference!
Bob's score: 0 (exists)
Charlie not found
```

---

#### Mistake 4: Concurrent Map Access Without Synchronization

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    m := make(map[int]int)
    
    // WRONG: No synchronization (will panic or data race)
    // for i := 0; i < 1000; i++ {
    //     go func(n int) {
    //         m[n] = n  // UNSAFE!
    //     }(i)
    // }
    
    // CORRECT: Use mutex
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            mu.Lock()
            m[n] = n
            mu.Unlock()
        }(i)
    }
    
    wg.Wait()
    fmt.Printf("Safe: map has %d elements\n", len(m))
}
```

**Output:**
```
Safe: map has 1000 elements
```

---

#### Mistake 5: Trying to Get Address of Map Element

```go
package main

import "fmt"

func main() {
    m := map[string]int{"Alice": 25}
    
    // This does NOT compile
    // ptr := &m["Alice"]  // cannot take the address of m["Alice"]
    
    // Why? Map may relocate elements during growth
    
    // CORRECT: If you need pointers, store pointers as values
    m2 := make(map[string]*int)
    age := 25
    m2["Alice"] = &age
    
    fmt.Println("Value:", *m2["Alice"])
}
```

**Output:**
```
Value: 25
```

---

<a name="maps-interview"></a>
### Interview Tips and Key Takeaways

#### 🎯 Frequently Asked Interview Questions

**Q1: What is the time complexity of map operations in Go?**

**Answer:**
- **Average case**: O(1) for insert, lookup, and delete
- **Worst case**: O(n) when there are hash collisions
- In practice, Go's hash function and collision resolution make worst case rare

**Q2: Can you compare two maps using ==?**

**Answer:** No, maps can only be compared to `nil`. To compare two maps, you must iterate and compare each key-value pair manually or use `reflect.DeepEqual()`.

**Q3: What happens if you access a key that doesn't exist?**

**Answer:** The map returns the **zero value** of the value type. Use the "comma ok" idiom to distinguish between a zero value and a missing key:
```go
value, exists := m["key"]
```

**Q4: Are maps safe for concurrent use?**

**Answer:** No, maps are **not safe** for concurrent reads and writes. Use `sync.Mutex` or `sync.Map` for concurrent access.

**Q5: What is the zero value of a map?**

**Answer:** The zero value is `nil`. You can read from a nil map (returns zero values), but writing to it causes a panic.

**Q6: What types can be used as map keys?**

**Answer:** Any **comparable** type:
- ✅ Valid: int, float, string, bool, pointers, structs (with comparable fields), arrays
- ❌ Invalid: slices, maps, functions (not comparable)

---

#### 🔥 Critical Points for Interviews

1. **Maps are Reference Types**
   ```go
   m1 := map[string]int{"a": 1}
   m2 := m1  // m2 points to same map
   m2["a"] = 2
   // m1["a"] is now 2
   ```

2. **Iteration Order is Random**
   ```go
   // NEVER assume order!
   for k, v := range m { }
   ```

3. **Comma OK Idiom**
   ```go
   if value, ok := m[key]; ok {
       // key exists
   }
   ```

4. **Nil Map vs Empty Map**
   ```go
   var nilMap map[string]int      // nil, can't write
   emptyMap := make(map[string]int)  // not nil, can write
   ```

5. **Cannot Take Address of Map Elements**
   ```go
   // ptr := &m["key"]  // Compile error!
   ```

---

#### ⚠️ Common Interview Pitfalls

| Pitfall | Problem | Solution |
|---------|---------|----------|
| **Writing to nil map** | Runtime panic | Initialize with `make()` or literal |
| **Concurrent access** | Data race/panic | Use `sync.Mutex` or `sync.Map` |
| **Assuming order** | Non-deterministic bugs | Sort keys if order matters |
| **Comparing maps** | Compile error | Use manual comparison or `reflect.DeepEqual()` |
| **Taking addresses** | Compile error | Store pointers as values if needed |

---

#### 💡 Best Practices for Production Code

✅ **DO:**
- Initialize maps before use: `m := make(map[K]V)`
- Use "comma ok" to check key existence
- Pre-allocate if size is known: `make(map[K]V, capacity)`
- Use `sync.Mutex` for concurrent access (or `sync.Map`)
- Sort keys if deterministic order is needed
- Return empty map instead of nil from functions (when appropriate)

❌ **DON'T:**
- Don't write to nil maps
- Don't assume iteration order
- Don't use maps concurrently without synchronization
- Don't compare maps with `==`
- Don't try to take addresses of map elements

---

<a name="strings"></a>
## 2. Strings and Runes

<a name="strings-introduction"></a>
### Introduction and Basics

A **string** in Go is a read-only slice of bytes. Strings are encoded in **UTF-8** by default.

**Key Characteristics:**
- **Immutable** - Cannot be modified after creation
- **UTF-8 encoded** - Supports international characters
- **Byte sequence** - Internally, a string is `[]byte`
- **Efficient** - String operations are optimized
- **Value type** - But with copy-on-write optimization

**String vs Rune:**
- **string**: Sequence of bytes
- **rune**: Alias for `int32`, represents a Unicode code point

---

<a name="strings-utf8"></a>
### UTF-8 Encoding

Go strings are UTF-8 encoded, meaning characters can be 1-4 bytes long.

**ASCII characters:** 1 byte  
**Special characters:** 2-4 bytes

```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    // ASCII string
    s1 := "Hello"
    fmt.Printf("'%s': %d bytes, %d runes\n", 
        s1, len(s1), utf8.RuneCountInString(s1))
    
    // String with multi-byte characters
    s2 := "Hello, 世界"  // "世界" = "World" in Chinese
    fmt.Printf("'%s': %d bytes, %d runes\n", 
        s2, len(s2), utf8.RuneCountInString(s2))
    
    // Emoji (4 bytes)
    s3 := "Go 🚀"
    fmt.Printf("'%s': %d bytes, %d runes\n", 
        s3, len(s3), utf8.RuneCountInString(s3))
}
```

**Output:**
```
'Hello': 5 bytes, 5 runes
'Hello, 世界': 13 bytes, 9 runes
'Go 🚀': 6 bytes, 4 runes
```

**Key Insight:**
- `len(s)` returns **number of bytes**
- `utf8.RuneCountInString(s)` returns **number of characters (runes)**

---

<a name="strings-bytes-runes"></a>
### Bytes vs Runes

**Byte (`uint8`):** 8-bit unsigned integer  
**Rune (`int32`):** 32-bit integer representing a Unicode code point

```go
package main

import "fmt"

func main() {
    s := "Hello, 世界"
    
    // Iterate over bytes
    fmt.Println("Bytes:")
    for i := 0; i < len(s); i++ {
        fmt.Printf("%d: %x (%c)\n", i, s[i], s[i])
    }
    
    fmt.Println("\nRunes:")
    // Iterate over runes
    for i, r := range s {
        fmt.Printf("%d: %U (%c)\n", i, r, r)
    }
}
```

**Output:**
```
Bytes:
0: 48 (H)
1: 65 (e)
2: 6c (l)
3: 6c (l)
4: 6f (o)
5: 2c (,)
6: 20 ( )
7: e4 (ä)
8: b8 (¸)
9: 96 ()
10: e7 (ç)
11: 95 ()
12: 8c ()

Runes:
0: U+0048 (H)
1: U+0065 (e)
2: U+006C (l)
3: U+006C (l)
4: U+006F (o)
5: U+002C (,)
6: U+0020 ( )
7: U+4E16 (世)
10: U+754C (界)
```

**Notice:** When iterating over bytes, multi-byte characters are split. When iterating over runes, each character is correctly recognized.

---

<a name="strings-operations"></a>
### String Declaration and Operations

#### Declaration

```go
package main

import "fmt"

func main() {
    // Method 1: Double quotes (interpreted string)
    s1 := "Hello, World!"
    
    // Method 2: Backticks (raw string literal)
    s2 := `This is a
multi-line
string`
    
    // Method 3: String concatenation
    s3 := "Hello" + ", " + "Go!"
    
    // Method 4: From bytes
    bytes := []byte{72, 101, 108, 108, 111}
    s4 := string(bytes)
    
    // Method 5: From runes
    runes := []rune{'H', 'e', 'l', 'l', 'o'}
    s5 := string(runes)
    
    fmt.Println(s1)
    fmt.Println(s2)
    fmt.Println(s3)
    fmt.Println(s4)
    fmt.Println(s5)
}
```

**Output:**
```
Hello, World!
This is a
multi-line
string
Hello, Go!
Hello
Hello
```

---

#### Basic String Operations

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, Go!"
    
    // Length (bytes)
    fmt.Println("Length:", len(s))
    
    // Access individual byte
    fmt.Printf("First byte: %c\n", s[0])
    
    // Substring (slicing)
    fmt.Println("Substring [0:5]:", s[0:5])
    fmt.Println("Substring [7:]:", s[7:])
    
    // Concatenation
    s2 := s + " Programming"
    fmt.Println("Concatenated:", s2)
    
    // Common operations with strings package
    fmt.Println("Contains 'Go':", strings.Contains(s, "Go"))
    fmt.Println("Starts with 'Hello':", strings.HasPrefix(s, "Hello"))
    fmt.Println("Ends with '!':", strings.HasSuffix(s, "!"))
    fmt.Println("Index of 'Go':", strings.Index(s, "Go"))
    fmt.Println("Uppercase:", strings.ToUpper(s))
    fmt.Println("Lowercase:", strings.ToLower(s))
    fmt.Println("Replace:", strings.Replace(s, "Go", "Python", 1))
}
```

**Output:**
```
Length: 10
First byte: H
Substring [0:5]: Hello
Substring [7:]: Go!
Concatenated: Hello, Go! Programming
Contains 'Go': true
Starts with 'Hello': true
Ends with '!': true
Index of 'Go': 7
Uppercase: HELLO, GO!
Lowercase: hello, go!
Replace: Hello, Python!
```

---

<a name="strings-iteration"></a>
### Iterating Over Strings

There are two main ways to iterate over strings, with very different results.

#### Method 1: Byte Iteration (Using Index)

```go
package main

import "fmt"

func main() {
    s := "Go语言"  // "语言" = "language" in Chinese
    
    fmt.Printf("String: %s (%d bytes)\n", s, len(s))
    
    // Iterate by index (bytes)
    fmt.Println("\nByte-by-byte iteration:")
    for i := 0; i < len(s); i++ {
        fmt.Printf("Index %d: %x\n", i, s[i])
    }
}
```

**Output:**
```
String: Go语言 (8 bytes)

Byte-by-byte iteration:
Index 0: 47
Index 1: 6f
Index 2: e8
Index 3: af
Index 4: ad
Index 5: e8
Index 6: a8
Index 7: 80
```

**Problem:** Multi-byte characters are split into individual bytes!

---

#### Method 2: Rune Iteration (Using range)

```go
package main

import "fmt"

func main() {
    s := "Go语言"
    
    fmt.Printf("String: %s\n", s)
    
    // Iterate using range (runes)
    fmt.Println("\nRune-by-rune iteration:")
    for index, runeValue := range s {
        fmt.Printf("Index %d: %U (%c)\n", index, runeValue, runeValue)
    }
}
```

**Output:**
```
String: Go语言

Rune-by-rune iteration:
Index 0: U+0047 (G)
Index 1: U+006F (o)
Index 2: U+8BED (语)
Index 5: U+8A00 (言)
```

**Notice:** 
- Index jumps from 2 to 5 (multi-byte character)
- Each character is correctly recognized

**Best Practice:** Use `range` when iterating over strings to handle multi-byte characters correctly.

---

#### Converting String to Runes and Bytes

```go
package main

import "fmt"

func main() {
    s := "Hello, 世界"
    
    // Convert to bytes
    bytes := []byte(s)
    fmt.Printf("Bytes: %v (length: %d)\n", bytes, len(bytes))
    
    // Convert to runes
    runes := []rune(s)
    fmt.Printf("Runes: %v (length: %d)\n", runes, len(runes))
    
    // Iterate over runes slice
    fmt.Println("\nRunes as characters:")
    for i, r := range runes {
        fmt.Printf("%d: %c\n", i, r)
    }
}
```

**Output:**
```
Bytes: [72 101 108 108 111 44 32 228 184 150 231 149 140] (length: 13)
Runes: [72 101 108 108 111 44 32 19990 30028] (length: 9)

Runes as characters:
0: H
1: e
2: l
3: l
4: o
5: ,
6:  
7: 世
8: 界
```

---

<a name="strings-immutability"></a>
### String Immutability

Strings in Go are **immutable** - you cannot change individual characters.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello"
    
    // This does NOT compile
    // s[0] = 'h'  // cannot assign to s[0]
    
    // Method 1: Create new string (concatenation)
    s2 := "h" + s[1:]
    fmt.Println("Method 1:", s2)
    
    // Method 2: Convert to []byte, modify, convert back
    bytes := []byte(s)
    bytes[0] = 'h'
    s3 := string(bytes)
    fmt.Println("Method 2:", s3)
    
    // Method 3: Convert to []rune (for multi-byte chars)
    runes := []rune(s)
    runes[0] = 'h'
    s4 := string(runes)
    fmt.Println("Method 3:", s4)
    
    // Method 4: Use strings.Replace
    s5 := strings.Replace(s, "H", "h", 1)
    fmt.Println("Method 4:", s5)
    
    // Original string is unchanged
    fmt.Println("Original:", s)
}
```

**Output:**
```
Method 1: hello
Method 2: hello
Method 3: hello
Method 4: hello
Original: Hello
```

**Why Immutability?**
- **Safety**: Strings can be shared without worry
- **Performance**: Compiler optimizations
- **Concurrency**: Safe to read from multiple goroutines

---

<a name="strings-examples"></a>
### More Code Examples

#### Example: String Builder for Efficient Concatenation

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

func concatWithPlus(n int) string {
    s := ""
    for i := 0; i < n; i++ {
        s += "a"
    }
    return s
}

func concatWithBuilder(n int) string {
    var builder strings.Builder
    for i := 0; i < n; i++ {
        builder.WriteString("a")
    }
    return builder.String()
}

func main() {
    n := 10000
    
    // Method 1: Using + (slow, creates many intermediate strings)
    start := time.Now()
    s1 := concatWithPlus(n)
    fmt.Printf("Using +: %v (length: %d)\n", time.Since(start), len(s1))
    
    // Method 2: Using strings.Builder (fast)
    start = time.Now()
    s2 := concatWithBuilder(n)
    fmt.Printf("Using Builder: %v (length: %d)\n", time.Since(start), len(s2))
}
```

**Output (example):**
```
Using +: 15ms (length: 10000)
Using Builder: 0.5ms (length: 10000)
```

**Best Practice:** Use `strings.Builder` for building strings in loops.

---

<a name="strings-common-ops"></a>
### Common String Operations

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // Split
    text := "apple,banana,cherry"
    fruits := strings.Split(text, ",")
    fmt.Println("Split:", fruits)
    
    // Join
    joined := strings.Join(fruits, " | ")
    fmt.Println("Join:", joined)
    
    // Trim
    s := "  Hello, World!  "
    fmt.Printf("Trim: '%s'\n", strings.TrimSpace(s))
    
    // Count occurrences
    count := strings.Count("hello", "l")
    fmt.Println("Count 'l' in 'hello':", count)
    
    // Repeat
    repeated := strings.Repeat("Go", 3)
    fmt.Println("Repeat:", repeated)
    
    // Replace all
    replaced := strings.ReplaceAll("foo foo foo", "foo", "bar")
    fmt.Println("Replace all:", replaced)
    
    // Fields (split by whitespace)
    sentence := "Go   is    great"
    words := strings.Fields(sentence)
    fmt.Println("Fields:", words)
}
```

**Output:**
```
Split: [apple banana cherry]
Join: apple | banana | cherry
Trim: 'Hello, World!'
Count 'l' in 'hello': 2
Repeat: GoGoGo
Replace all: bar bar bar
Fields: [Go is great]
```

---

<a name="strings-performance"></a>
### Performance Considerations

#### String Concatenation Performance

```go
package main

import (
    "bytes"
    "fmt"
    "strings"
    "time"
)

func benchmarkConcat(n int) {
    // Method 1: Using + operator
    start := time.Now()
    s := ""
    for i := 0; i < n; i++ {
        s += "x"
    }
    fmt.Printf("+ operator: %v\n", time.Since(start))
    
    // Method 2: Using strings.Builder
    start = time.Now()
    var builder strings.Builder
    for i := 0; i < n; i++ {
        builder.WriteString("x")
    }
    _ = builder.String()
    fmt.Printf("strings.Builder: %v\n", time.Since(start))
    
    // Method 3: Using bytes.Buffer
    start = time.Now()
    var buffer bytes.Buffer
    for i := 0; i < n; i++ {
        buffer.WriteString("x")
    }
    _ = buffer.String()
    fmt.Printf("bytes.Buffer: %v\n", time.Since(start))
}

func main() {
    fmt.Println("Concatenating 10,000 strings:")
    benchmarkConcat(10000)
}
```

**Output (example):**
```
Concatenating 10,000 strings:
+ operator: 12ms
strings.Builder: 0.3ms
bytes.Buffer: 0.4ms
```

**Recommendation:**
- Use `+` for simple, few concatenations
- Use `strings.Builder` for loops and many concatenations
- `strings.Builder` is faster than `bytes.Buffer` for string building

---

<a name="strings-mistakes"></a>
### Common Mistakes and Pitfalls

#### Mistake 1: Confusing Bytes and Runes

```go
package main

import "fmt"

func main() {
    s := "Hello, 世界"
    
    // WRONG: Assuming len() gives character count
    fmt.Printf("Wrong assumption: %d characters\n", len(s))
    
    // CORRECT: Use utf8.RuneCountInString or convert to []rune
    fmt.Printf("Correct: %d characters\n", len([]rune(s)))
}
```

**Output:**
```
Wrong assumption: 13 characters
Correct: 9 characters
```

---

#### Mistake 2: Trying to Modify Strings

```go
package main

func main() {
    s := "Hello"
    
    // This does NOT compile
    // s[0] = 'h'  // cannot assign to s[0]
    
    // CORRECT: Create a new string
    s = "h" + s[1:]
}
```

---

#### Mistake 3: Incorrect String Slicing with Multi-byte Characters

```go
package main

import "fmt"

func main() {
    s := "世界Hello"
    
    // WRONG: Slicing by bytes can break multi-byte characters
    // First Chinese character is 3 bytes (indices 0, 1, 2)
    broken := s[0:2]  // Only part of first character!
    fmt.Printf("Broken: %s (shows incorrectly)\n", broken)
    
    // CORRECT: Convert to runes first
    runes := []rune(s)
    correct := string(runes[0:2])
    fmt.Printf("Correct: %s\n", correct)
}
```

**Output:**
```
Broken: � (shows incorrectly)
Correct: 世界
```

---

#### Mistake 4: Inefficient String Building

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // WRONG: Inefficient in loops
    result := ""
    for i := 0; i < 1000; i++ {
        result += "a"  // Creates new string each iteration!
    }
    
    // CORRECT: Use strings.Builder
    var builder strings.Builder
    for i := 0; i < 1000; i++ {
        builder.WriteString("a")
    }
    result = builder.String()
    fmt.Printf("Built string of length: %d\n", len(result))
}
```

**Output:**
```
Built string of length: 1000
```

---

<a name="strings-interview"></a>
### Interview Tips and Key Takeaways

#### 🎯 Frequently Asked Interview Questions

**Q1: Are strings mutable or immutable in Go?**

**Answer:** Strings are **immutable**. Once created, they cannot be changed. Any operation that appears to modify a string actually creates a new string.

**Q2: What's the difference between `len(s)` and `utf8.RuneCountInString(s)`?**

**Answer:**
- `len(s)` returns the number of **bytes**
- `utf8.RuneCountInString(s)` returns the number of **characters (runes)**

For ASCII strings, they're the same. For strings with multi-byte characters, they differ.

**Q3: How do you iterate over a string correctly?**

**Answer:** Use `range`:
```go
for index, runeValue := range s {
    // runeValue is a rune, handles multi-byte characters correctly
}
```

**Q4: What's the difference between a byte and a rune?**

**Answer:**
- **byte** (`uint8`): 8-bit value, represents a single byte
- **rune** (`int32`): 32-bit value, represents a Unicode code point (character)

**Q5: How do you efficiently build strings in a loop?**

**Answer:** Use `strings.Builder`:
```go
var builder strings.Builder
for i := 0; i < n; i++ {
    builder.WriteString("text")
}
result := builder.String()
```

**Q6: Can you compare strings using ==?**

**Answer:** Yes, strings can be compared using `==`, `!=`, `<`, `>`, `<=`, `>=`. Comparison is lexicographic (byte-by-byte).

---

#### 🔥 Critical Points for Interviews

1. **UTF-8 Encoding**
   ```go
   s := "Hello, 世界"
   len(s)                      // 13 bytes
   utf8.RuneCountInString(s)   // 9 runes
   ```

2. **String Immutability**
   ```go
   s := "Hello"
   // s[0] = 'h'  // Compile error
   s = "h" + s[1:]  // Creates new string
   ```

3. **Byte vs Rune Iteration**
   ```go
   for i := 0; i < len(s); i++ { }  // Bytes
   for i, r := range s { }           // Runes (correct for Unicode)
   ```

4. **Efficient String Building**
   ```go
   var builder strings.Builder
   builder.WriteString("text")
   result := builder.String()
   ```

5. **String Slicing**
   ```go
   s := "Hello"
   s[1:4]  // "ell" - creates new string
   ```

---

#### ⚠️ Common Interview Pitfalls

| Pitfall | Problem | Solution |
|---------|---------|----------|
| **len() for characters** | Returns bytes, not characters | Use `utf8.RuneCountInString()` or `len([]rune(s))` |
| **Byte iteration on Unicode** | Breaks multi-byte characters | Use `range` for rune iteration |
| **String modification** | Strings are immutable | Create new string or convert to []byte/[]rune |
| **+ in loops** | Inefficient, many allocations | Use `strings.Builder` |
| **Byte slicing** | Can break multi-byte characters | Convert to []rune first |

---

#### 💡 Best Practices for Production Code

✅ **DO:**
- Use `range` to iterate over strings (handles Unicode correctly)
- Use `strings.Builder` for building strings in loops
- Use `strings` package for common operations
- Check string length with `utf8.RuneCountInString()` for multi-byte strings
- Use raw string literals `` `...` `` for multi-line or literal strings

❌ **DON'T:**
- Don't assume `len(s)` gives character count
- Don't iterate by index for Unicode strings
- Don't use `+` for concatenation in loops
- Don't try to modify strings directly
- Don't slice strings by bytes if they contain multi-byte characters

---

## Summary - Part 2

### Key Takeaways

**Maps:**
- ✓ Unordered key-value collections
- ✓ Reference types (passed by reference)
- ✓ Keys must be comparable types
- ✓ Zero value is `nil` (can't write to nil map)
- ✓ Use "comma ok" idiom to check existence
- ✓ Iteration order is random
- ✓ Not safe for concurrent use (use `sync.Mutex` or `sync.Map`)
- ✓ Pre-allocate with `make(map[K]V, size)` for better performance

**Strings and Runes:**
- ✓ Strings are immutable byte sequences
- ✓ UTF-8 encoded by default
- ✓ `len(s)` returns bytes, not characters
- ✓ Use `range` for correct Unicode iteration
- ✓ `byte` (uint8) vs `rune` (int32, Unicode code point)
- ✓ Use `strings.Builder` for efficient concatenation
- ✓ Convert to `[]rune` for character-based operations
- ✓ Strings can be compared with `==`, `<`, etc.

### Interview Readiness Checklist

**Maps:**
- [ ] Understand reference behavior
- [ ] Know when to use "comma ok" idiom
- [ ] Understand nil vs empty map
- [ ] Know iteration order is random
- [ ] Understand concurrency issues
- [ ] Know which types can be keys
- [ ] Can explain time complexity

**Strings:**
- [ ] Understand UTF-8 encoding
- [ ] Know difference between bytes and runes
- [ ] Can iterate correctly over Unicode strings
- [ ] Understand string immutability
- [ ] Know efficient string building techniques
- [ ] Can explain `len()` vs `RuneCountInString()`
- [ ] Understand string slicing pitfalls

---

**Next:** Part 3 will cover **Structs** and **Zero Values** 🚀
