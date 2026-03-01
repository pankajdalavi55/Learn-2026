# Go Core Data Structures - Part 1: Arrays and Slices

## Table of Contents
1. [Arrays](#arrays)
   - [Introduction and Basics](#arrays-introduction)
   - [Syntax and Declaration](#arrays-syntax)
   - [Code Examples](#arrays-examples)
   - [Memory Behavior](#arrays-memory)
   - [Common Mistakes](#arrays-mistakes)
   - [Interview Tips](#arrays-interview)
2. [Slices](#slices)
   - [Introduction and Basics](#slices-introduction)
   - [Syntax and Declaration](#slices-syntax)
   - [Length vs Capacity](#slices-length-capacity)
   - [Slice Expressions](#slice-expressions)
   - [The append() Function](#slice-append)
   - [Underlying Array Sharing](#slice-sharing)
   - [Copying Slices Safely](#slice-copying)
   - [Code Examples](#slices-examples)
   - [Memory Behavior](#slices-memory)
   - [Common Mistakes](#slices-mistakes)
   - [Interview Tips](#slices-interview)
3. [Arrays vs Slices: Complete Comparison](#arrays-vs-slices)

---

<a name="arrays"></a>
## 1. Arrays

<a name="arrays-introduction"></a>
### Introduction and Basics

An **array** in Go is a fixed-size, numbered sequence of elements of a single type. Once declared, the size of an array cannot be changed.

**Key Characteristics:**
- **Fixed size** - Length is part of the array's type
- **Value type** - Arrays are passed by value (copied when assigned or passed to functions)
- **Zero-indexed** - First element is at index 0
- **Homogeneous** - All elements must be the same type
- **Contiguous memory** - Elements stored sequentially in memory

<a name="arrays-syntax"></a>
### Syntax and Declaration

```go
// Syntax
var arrayName [size]Type

// Examples
var numbers [5]int                    // Array of 5 integers
var names [3]string                   // Array of 3 strings
var flags [4]bool                     // Array of 4 booleans
```

<a name="arrays-examples"></a>
### Code Examples

#### Example 1: Basic Array Declaration and Initialization

```go
package main

import "fmt"

func main() {
    // Declaration without initialization (zero values)
    var numbers [5]int
    fmt.Println("Zero-valued array:", numbers)
    
    // Declaration with initialization
    var fruits [3]string = [3]string{"Apple", "Banana", "Cherry"}
    fmt.Println("Fruits:", fruits)
    
    // Short declaration with initialization
    colors := [4]string{"Red", "Green", "Blue", "Yellow"}
    fmt.Println("Colors:", colors)
    
    // Array literal with size inference
    days := [...]string{"Mon", "Tue", "Wed", "Thu", "Fri"}
    fmt.Println("Days:", days)
    fmt.Println("Length:", len(days))
}
```

**Output:**
```
Zero-valued array: [0 0 0 0 0]
Fruits: [Apple Banana Cherry]
Colors: [Red Green Blue Yellow]
Days: [Mon Tue Wed Thu Fri]
Length: 5
```

---

#### Example 2: Accessing and Modifying Array Elements

```go
package main

import "fmt"

func main() {
    scores := [5]int{85, 90, 78, 92, 88}
    
    // Accessing elements
    fmt.Println("First score:", scores[0])
    fmt.Println("Last score:", scores[4])
    
    // Modifying elements
    scores[2] = 95
    fmt.Println("Updated scores:", scores)
    
    // Array length
    fmt.Println("Number of scores:", len(scores))
    
    // Accessing out of bounds causes compile error
    // fmt.Println(scores[10]) // This will NOT compile
}
```

**Output:**
```
First score: 85
Last score: 88
Updated scores: [85 90 95 92 88]
Number of scores: 5
```

---

#### Example 3: Iterating Over Arrays

```go
package main

import "fmt"

func main() {
    temperatures := [7]float64{22.5, 23.0, 21.5, 24.0, 25.5, 23.5, 22.0}
    
    // Method 1: Using traditional for loop
    fmt.Println("Method 1: Traditional for loop")
    for i := 0; i < len(temperatures); i++ {
        fmt.Printf("Day %d: %.1f°C\n", i+1, temperatures[i])
    }
    
    fmt.Println("\nMethod 2: Using range (index and value)")
    // Method 2: Using range
    for index, temp := range temperatures {
        fmt.Printf("Day %d: %.1f°C\n", index+1, temp)
    }
    
    fmt.Println("\nMethod 3: Using range (value only)")
    // Method 3: Using range (ignoring index)
    dayNum := 1
    for _, temp := range temperatures {
        fmt.Printf("Day %d: %.1f°C\n", dayNum, temp)
        dayNum++
    }
}
```

**Output:**
```
Method 1: Traditional for loop
Day 1: 22.5°C
Day 2: 23.0°C
Day 3: 21.5°C
Day 4: 24.0°C
Day 5: 25.5°C
Day 6: 23.5°C
Day 7: 22.0°C

Method 2: Using range (index and value)
Day 1: 22.5°C
Day 2: 23.0°C
Day 3: 21.5°C
Day 4: 24.0°C
Day 5: 25.5°C
Day 6: 23.5°C
Day 7: 22.0°C

Method 3: Using range (value only)
Day 1: 22.5°C
Day 2: 23.0°C
Day 3: 21.5°C
Day 4: 24.0°C
Day 5: 25.5°C
Day 6: 23.5°C
Day 7: 22.0°C
```

---

#### Example 4: Arrays are Value Types (Passed by Copy)

```go
package main

import "fmt"

func modifyArray(arr [3]int) {
    arr[0] = 999
    fmt.Println("Inside function:", arr)
}

func main() {
    original := [3]int{1, 2, 3}
    fmt.Println("Before function call:", original)
    
    // Array is passed by value (copied)
    modifyArray(original)
    
    fmt.Println("After function call:", original)
    // Original array is unchanged!
}
```

**Output:**
```
Before function call: [1 2 3]
Inside function: [999 2 3]
After function call: [1 2 3]
```

**Key Point:** The original array is NOT modified because arrays are value types and are copied when passed to functions.

---

#### Example 5: Array Comparison

```go
package main

import "fmt"

func main() {
    arr1 := [3]int{1, 2, 3}
    arr2 := [3]int{1, 2, 3}
    arr3 := [3]int{1, 2, 4}
    
    // Arrays can be compared using == and !=
    fmt.Println("arr1 == arr2:", arr1 == arr2) // true
    fmt.Println("arr1 == arr3:", arr1 == arr3) // false
    fmt.Println("arr1 != arr3:", arr1 != arr3) // true
    
    // Arrays of different sizes cannot be compared
    // arr4 := [4]int{1, 2, 3, 4}
    // fmt.Println(arr1 == arr4) // Compile error: mismatched types
}
```

**Output:**
```
arr1 == arr2: true
arr1 == arr3: false
arr1 != arr3: true
```

---

#### Example 6: Multidimensional Arrays

```go
package main

import "fmt"

func main() {
    // 2D array (3 rows, 4 columns)
    var matrix [3][4]int
    
    // Initialize with values
    matrix[0] = [4]int{1, 2, 3, 4}
    matrix[1] = [4]int{5, 6, 7, 8}
    matrix[2] = [4]int{9, 10, 11, 12}
    
    fmt.Println("2D Array:")
    for i := 0; i < len(matrix); i++ {
        for j := 0; j < len(matrix[i]); j++ {
            fmt.Printf("%3d ", matrix[i][j])
        }
        fmt.Println()
    }
    
    // 2D array with initialization
    grid := [2][3]string{
        {"A1", "A2", "A3"},
        {"B1", "B2", "B3"},
    }
    
    fmt.Println("\nGrid:")
    for _, row := range grid {
        fmt.Println(row)
    }
}
```

**Output:**
```
2D Array:
  1   2   3   4 
  5   6   7   8 
  9  10  11  12 

Grid:
[A1 A2 A3]
[B1 B2 B3]
```

---

<a name="arrays-memory"></a>
### Memory Behavior and Internal Representation

**Memory Layout:**
```
Array: [5]int{10, 20, 30, 40, 50}

Memory: [10][20][30][40][50]
         ↑   ↑   ↑   ↑   ↑
      index 0 1  2  3  4
      
All elements stored contiguously in memory
```

**Key Memory Characteristics:**

1. **Contiguous Allocation**: All elements are stored sequentially in memory
2. **Stack Allocation**: Small arrays are typically allocated on the stack
3. **Value Semantics**: Copying an array copies all elements
4. **Fixed Size**: Memory size is known at compile time

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    arr := [5]int{1, 2, 3, 4, 5}
    
    // Size of the entire array
    fmt.Printf("Size of array: %d bytes\n", unsafe.Sizeof(arr))
    
    // Size of one element
    fmt.Printf("Size of one int: %d bytes\n", unsafe.Sizeof(arr[0]))
    
    // Total size = 5 elements × 8 bytes = 40 bytes (on 64-bit system)
}
```

**Output:**
```
Size of array: 40 bytes
Size of one int: 8 bytes
```

---

<a name="arrays-mistakes"></a>
### Common Mistakes and Pitfalls

#### Mistake 1: Arrays of Different Sizes are Different Types

```go
package main

import "fmt"

func printArray(arr [3]int) {
    fmt.Println(arr)
}

func main() {
    arr1 := [3]int{1, 2, 3}
    arr2 := [5]int{1, 2, 3, 4, 5}
    
    printArray(arr1) // OK
    // printArray(arr2) // Compile error: cannot use arr2 (type [5]int) as type [3]int
    
    // Even same values, different sizes = different types
    // var a [3]int = arr2 // Compile error
}
```

**Output:**
```
[1 2 3]
```

**Lesson:** `[3]int` and `[5]int` are completely different types in Go.

---

#### Mistake 2: Expecting Pass-by-Reference Behavior

```go
package main

import "fmt"

func doubleValues(arr [3]int) {
    for i := range arr {
        arr[i] *= 2
    }
    // Changes are lost when function returns!
}

func doubleValuesCorrect(arr *[3]int) {
    for i := range arr {
        arr[i] *= 2
    }
}

func main() {
    numbers := [3]int{1, 2, 3}
    
    // Wrong approach
    doubleValues(numbers)
    fmt.Println("After doubleValues:", numbers) // Unchanged!
    
    // Correct approach: pass pointer
    doubleValuesCorrect(&numbers)
    fmt.Println("After doubleValuesCorrect:", numbers) // Modified!
}
```

**Output:**
```
After doubleValues: [1 2 3]
After doubleValuesCorrect: [2 4 6]
```

**Lesson:** Pass a pointer to the array if you need to modify it in a function.

---

#### Mistake 3: Large Arrays Performance Issues

```go
package main

import (
    "fmt"
    "time"
)

func processLargeArray(arr [1000000]int) int {
    sum := 0
    for _, v := range arr {
        sum += v
    }
    return sum
}

func processArrayPointer(arr *[1000000]int) int {
    sum := 0
    for _, v := range arr {
        sum += v
    }
    return sum
}

func main() {
    var bigArray [1000000]int
    for i := range bigArray {
        bigArray[i] = i
    }
    
    // Passing by value (slow - copies 8MB of data)
    start := time.Now()
    _ = processLargeArray(bigArray)
    fmt.Printf("By value: %v\n", time.Since(start))
    
    // Passing by pointer (fast - copies only 8 bytes)
    start = time.Now()
    _ = processArrayPointer(&bigArray)
    fmt.Printf("By pointer: %v\n", time.Since(start))
}
```

**Output (example):**
```
By value: 3.5ms
By pointer: 1.2ms
```

**Lesson:** For large arrays, pass pointers to avoid expensive copying, or use slices instead.

---

<a name="arrays-interview"></a>
### Interview Tips and Key Takeaways

#### 🎯 Frequently Asked Interview Questions

**Q1: What is the difference between an array and a slice in Go?**

**Answer:**
- **Array**: Fixed size, value type, size is part of the type, copied when assigned
- **Slice**: Dynamic size, reference type (backed by array), flexible, efficient to pass

**Q2: Are arrays passed by value or reference in Go?**

**Answer:** Arrays are **passed by value** (copied). To avoid copying, pass a pointer `*[N]Type` or use slices instead.

**Q3: Can you compare arrays in Go?**

**Answer:** Yes, arrays of the **same type** (same size and element type) can be compared using `==` and `!=`. Arrays of different sizes cannot be compared.

**Q4: What is the zero value of an array?**

**Answer:** An array's zero value is an array with all elements set to the zero value of the element type (e.g., `0` for int, `""` for string, `false` for bool).

---

#### 🔥 Important Points to Remember

1. **Size is Part of Type**
   ```go
   var a [3]int  // Type: [3]int
   var b [5]int  // Type: [5]int - DIFFERENT TYPE!
   ```

2. **Arrays are Value Types**
   ```go
   a := [3]int{1, 2, 3}
   b := a           // b is a COPY of a
   b[0] = 999
   fmt.Println(a)   // [1 2 3] - unchanged
   ```

3. **Array Length is Compile-Time Constant**
   ```go
   // Valid
   const size = 5
   var arr [size]int
   
   // Invalid
   n := 5
   // var arr2 [n]int  // Compile error: non-constant array bound
   ```

4. **Bounds Checking**
   ```go
   arr := [3]int{1, 2, 3}
   // arr[5] = 10  // Compile error or runtime panic
   ```

---

#### ⚠️ Common Pitfalls for Interviews

| Pitfall | Issue | Solution |
|---------|-------|----------|
| **Different sizes** | `[3]int` and `[5]int` are incompatible | Understand type system |
| **Large arrays** | Inefficient copying | Use pointers or slices |
| **Expecting mutations** | Functions can't modify array | Pass pointer `*[N]Type` |
| **Runtime size** | Can't use variable for size | Use slices for dynamic sizing |

---

#### 💡 Best Practices

✅ **Use arrays when:**
- Size is known and fixed at compile time
- You need value semantics (defensive copying)
- Working with small, fixed-size collections
- Performance-critical code with known bounds

❌ **Avoid arrays when:**
- Size is dynamic or unknown at compile time
- Passing to functions frequently (expensive copying)
- Need to append/remove elements
- Working with large data sets

**Pro Tip:** In most Go code, **slices are preferred over arrays** due to their flexibility and efficiency.

---

<a name="slices"></a>
## 2. Slices

<a name="slices-introduction"></a>
### Introduction and Basics

A **slice** is a dynamically-sized, flexible view into the elements of an array. Slices are the most commonly used collection type in Go.

**Key Characteristics:**
- **Dynamic size** - Can grow and shrink
- **Reference type** - Points to an underlying array
- **Three components** - Pointer, length, and capacity
- **Efficient** - Passed by reference (no copying of elements)
- **Flexible** - Can be appended, sliced, and manipulated

**Slice Structure:**
```
Slice internally has 3 fields:
┌─────────┬────────┬──────────┐
│ Pointer │ Length │ Capacity │
└─────────┴────────┴──────────┘
    ↓
Underlying array
```

---

<a name="slices-syntax"></a>
### Syntax and Declaration

```go
// Syntax
var sliceName []Type

// Examples
var numbers []int              // Slice of integers
var names []string             // Slice of strings
var flags []bool               // Slice of booleans
```

**Note:** Unlike arrays, slices do **not** have a size in the type declaration.

---

<a name="slices-length-capacity"></a>
### Length vs Capacity

**Length (`len`)**: Number of elements in the slice

**Capacity (`cap`)**: Number of elements in the underlying array from the slice's first element

```go
package main

import "fmt"

func main() {
    // Create a slice with make
    numbers := make([]int, 3, 5)
    // length = 3, capacity = 5
    
    fmt.Printf("Slice: %v\n", numbers)
    fmt.Printf("Length: %d\n", len(numbers))
    fmt.Printf("Capacity: %d\n", cap(numbers))
    
    // Visual representation
    fmt.Println("\nVisual:")
    fmt.Println("Slice view: [0 0 0]")
    fmt.Println("Underlying array capacity: [0 0 0 _ _]")
}
```

**Output:**
```
Slice: [0 0 0]
Length: 3
Capacity: 5

Visual:
Slice view: [0 0 0]
Underlying array capacity: [0 0 0 _ _]
```

**Understanding Length and Capacity:**

```go
package main

import "fmt"

func main() {
    s := make([]int, 3, 5)
    s[0], s[1], s[2] = 1, 2, 3
    
    fmt.Printf("Initial - len: %d, cap: %d, slice: %v\n", len(s), cap(s), s)
    
    // Append uses capacity before allocating new array
    s = append(s, 4)
    fmt.Printf("After append(4) - len: %d, cap: %d, slice: %v\n", len(s), cap(s), s)
    
    s = append(s, 5)
    fmt.Printf("After append(5) - len: %d, cap: %d, slice: %v\n", len(s), cap(s), s)
    
    // Capacity exceeded, new array allocated (capacity typically doubles)
    s = append(s, 6)
    fmt.Printf("After append(6) - len: %d, cap: %d, slice: %v\n", len(s), cap(s), s)
}
```

**Output:**
```
Initial - len: 3, cap: 5, slice: [1 2 3]
After append(4) - len: 4, cap: 5, slice: [1 2 3 4]
After append(5) - len: 5, cap: 5, slice: [1 2 3 4 5]
After append(6) - len: 6, cap: 10, slice: [1 2 3 4 5 6]
```

**Key Insight:** When capacity is reached, `append()` allocates a new, larger array (usually double the size) and copies all elements.

---

<a name="slice-expressions"></a>
### Slice Expressions

Slice expressions create a new slice from an existing array or slice.

**Syntax:** `slice[low:high]` (includes `low`, excludes `high`)

```go
package main

import "fmt"

func main() {
    numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    
    // Basic slicing
    fmt.Println("Full slice:", numbers)        // [0 1 2 3 4 5 6 7 8 9]
    fmt.Println("numbers[2:5]:", numbers[2:5]) // [2 3 4]
    fmt.Println("numbers[:4]:", numbers[:4])   // [0 1 2 3]
    fmt.Println("numbers[6:]:", numbers[6:])   // [6 7 8 9]
    fmt.Println("numbers[:]:", numbers[:])     // [0 1 2 3 4 5 6 7 8 9]
    
    // Slicing a slice
    subset := numbers[3:7]
    fmt.Println("subset:", subset)              // [3 4 5 6]
    
    // Further slicing
    fmt.Println("subset[1:3]:", subset[1:3])   // [4 5]
}
```

**Output:**
```
Full slice: [0 1 2 3 4 5 6 7 8 9]
numbers[2:5]: [2 3 4]
numbers[:4]: [0 1 2 3]
numbers[6:]: [6 7 8 9]
numbers[:]: [0 1 2 3 4 5 6 7 8 9]
subset: [3 4 5 6]
subset[1:3]: [4 5]
```

**Full Slice Expression (3-index):**

Syntax: `slice[low:high:max]` - sets capacity

```go
package main

import "fmt"

func main() {
    original := []int{0, 1, 2, 3, 4, 5}
    
    // Normal slice: s = original[1:4]
    s1 := original[1:4]
    fmt.Printf("s1 = original[1:4] - len: %d, cap: %d, %v\n", len(s1), cap(s1), s1)
    
    // Full slice expression: s = original[1:4:4]
    s2 := original[1:4:4]
    fmt.Printf("s2 = original[1:4:4] - len: %d, cap: %d, %v\n", len(s2), cap(s2), s2)
}
```

**Output:**
```
s1 = original[1:4] - len: 3, cap: 5, [1 2 3]
s2 = original[1:4:4] - len: 3, cap: 3, [1 2 3]
```

**Why use 3-index slicing?** To prevent a slice from accessing parts of the underlying array beyond what's intended, useful for defensive programming.

---

<a name="slice-append"></a>
### The append() Function

`append()` adds elements to a slice and returns a new slice.

#### Example 1: Basic Append

```go
package main

import "fmt"

func main() {
    var numbers []int
    fmt.Printf("Initial: len=%d, cap=%d, %v\n", len(numbers), cap(numbers), numbers)
    
    // Append single elements
    numbers = append(numbers, 1)
    fmt.Printf("After append(1): len=%d, cap=%d, %v\n", len(numbers), cap(numbers), numbers)
    
    numbers = append(numbers, 2)
    fmt.Printf("After append(2): len=%d, cap=%d, %v\n", len(numbers), cap(numbers), numbers)
    
    numbers = append(numbers, 3)
    fmt.Printf("After append(3): len=%d, cap=%d, %v\n", len(numbers), cap(numbers), numbers)
    
    // Append multiple elements
    numbers = append(numbers, 4, 5, 6)
    fmt.Printf("After append(4,5,6): len=%d, cap=%d, %v\n", len(numbers), cap(numbers), numbers)
    
    // Append another slice
    moreNumbers := []int{7, 8, 9}
    numbers = append(numbers, moreNumbers...)
    fmt.Printf("After append slice: len=%d, cap=%d, %v\n", len(numbers), cap(numbers), numbers)
}
```

**Output:**
```
Initial: len=0, cap=0, []
After append(1): len=1, cap=1, [1]
After append(2): len=2, cap=2, [1 2]
After append(3): len=3, cap=4, [1 2 3]
After append(4,5,6): len=6, cap=8, [1 2 3 4 5 6]
After append slice: len=9, cap=16, [1 2 3 4 5 6 7 8 9]
```

**Notice:** Capacity grows automatically (typically doubles when exceeded).

---

#### Example 2: Append Behavior and Reallocation

```go
package main

import "fmt"

func main() {
    // Create slice with specific capacity
    s := make([]int, 0, 3)
    fmt.Printf("Initial: len=%d, cap=%d\n", len(s), cap(s))
    
    for i := 1; i <= 6; i++ {
        s = append(s, i)
        fmt.Printf("Append %d: len=%d, cap=%d, slice=%v\n", i, len(s), cap(s), s)
    }
}
```

**Output:**
```
Initial: len=0, cap=3
Append 1: len=1, cap=3, slice=[1]
Append 2: len=2, cap=3, slice=[1 2]
Append 3: len=3, cap=3, slice=[1 2 3]
Append 4: len=4, cap=6, slice=[1 2 3 4]
Append 5: len=5, cap=6, slice=[1 2 3 4 5]
Append 6: len=6, cap=6, slice=[1 2 3 4 5 6]
```

**Key Observation:** When capacity is exceeded (at element 4), a new array with larger capacity is allocated.

---

<a name="slice-sharing"></a>
### Underlying Array Sharing

**Critical Concept:** Multiple slices can share the same underlying array!

```go
package main

import "fmt"

func main() {
    original := []int{0, 1, 2, 3, 4, 5}
    fmt.Println("Original:", original)
    
    // Create slices from original
    slice1 := original[1:4]  // [1 2 3]
    slice2 := original[2:5]  // [2 3 4]
    
    fmt.Println("Slice1:", slice1)
    fmt.Println("Slice2:", slice2)
    
    // Modify slice1
    slice1[1] = 999
    
    // All slices sharing the array are affected!
    fmt.Println("\nAfter slice1[1] = 999:")
    fmt.Println("Original:", original)  // [0 1 999 3 4 5]
    fmt.Println("Slice1:", slice1)      // [1 999 3]
    fmt.Println("Slice2:", slice2)      // [999 3 4]
}
```

**Output:**
```
Original: [0 1 2 3 4 5]
Slice1: [1 2 3]
Slice2: [2 3 4]

After slice1[1] = 999:
Original: [0 1 999 3 4 5]
Slice1: [1 999 3]
Slice2: [999 3 4]
```

**Explanation:** All three slices point to the same underlying array, so changes to one affect all.

---

#### Append Can Break Sharing

```go
package main

import "fmt"

func main() {
    original := []int{1, 2, 3}
    slice1 := original        // Share array
    slice2 := original        // Share array
    
    fmt.Println("Initial:")
    fmt.Printf("original: %v (cap: %d)\n", original, cap(original))
    fmt.Printf("slice1: %v (cap: %d)\n", slice1, cap(slice1))
    fmt.Printf("slice2: %v (cap: %d)\n", slice2, cap(slice2))
    
    // Append to slice1 (triggers reallocation)
    slice1 = append(slice1, 4, 5, 6)
    
    fmt.Println("\nAfter slice1 = append(slice1, 4, 5, 6):")
    fmt.Printf("original: %v (cap: %d)\n", original, cap(original))
    fmt.Printf("slice1: %v (cap: %d)\n", slice1, cap(slice1))
    fmt.Printf("slice2: %v (cap: %d)\n", slice2, cap(slice2))
    
    // Modify slice1
    slice1[0] = 999
    
    fmt.Println("\nAfter slice1[0] = 999:")
    fmt.Println("original:", original)  // Unchanged
    fmt.Println("slice1:", slice1)      // Changed
    fmt.Println("slice2:", slice2)      // Unchanged
}
```

**Output:**
```
Initial:
original: [1 2 3] (cap: 3)
slice1: [1 2 3] (cap: 3)
slice2: [1 2 3] (cap: 3)

After slice1 = append(slice1, 4, 5, 6):
original: [1 2 3] (cap: 3)
slice1: [1 2 3 4 5 6] (cap: 6)
slice2: [1 2 3] (cap: 3)

After slice1[0] = 999:
original: [1 2 3]
slice1: [999 2 3 4 5 6]
slice2: [1 2 3]
```

**Explanation:** `append()` caused `slice1` to get a new underlying array, breaking the sharing.

---

<a name="slice-copying"></a>
### Copying Slices Safely

To avoid unintended sharing, use the `copy()` function.

```go
package main

import "fmt"

func main() {
    original := []int{1, 2, 3, 4, 5}
    
    // Method 1: Using copy()
    copied := make([]int, len(original))
    n := copy(copied, original)
    fmt.Printf("Copied %d elements\n", n)
    
    // Modify copy
    copied[0] = 999
    
    fmt.Println("Original:", original)  // [1 2 3 4 5]
    fmt.Println("Copied:", copied)      // [999 2 3 4 5]
    
    // Method 2: Using append with empty slice
    cloned := append([]int{}, original...)
    cloned[1] = 888
    
    fmt.Println("Original:", original)  // [1 2 3 4 5]
    fmt.Println("Cloned:", cloned)      // [1 888 3 4 5]
}
```

**Output:**
```
Copied 5 elements
Original: [1 2 3 4 5]
Copied: [999 2 3 4 5]
Original: [1 2 3 4 5]
Cloned: [1 888 3 4 5]
```

**copy() Function Details:**

```go
package main

import "fmt"

func main() {
    src := []int{1, 2, 3, 4, 5}
    
    // Destination smaller than source
    dst1 := make([]int, 3)
    copy(dst1, src)
    fmt.Println("dst1 (size 3):", dst1)  // [1 2 3] - only first 3 copied
    
    // Destination larger than source
    dst2 := make([]int, 10)
    n := copy(dst2, src)
    fmt.Printf("dst2 (size 10): %v, copied %d elements\n", dst2, n)
}
```

**Output:**
```
dst1 (size 3): [1 2 3]
dst2 (size 10): [1 2 3 4 5 0 0 0 0 0], copied 5 elements
```

**Rule:** `copy()` copies `min(len(dst), len(src))` elements.

---

<a name="slices-examples"></a>
### More Code Examples

#### Example: Creating Slices (Different Methods)

```go
package main

import "fmt"

func main() {
    // Method 1: Slice literal
    s1 := []int{1, 2, 3, 4, 5}
    fmt.Printf("s1: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
    
    // Method 2: Using make (length only)
    s2 := make([]int, 5)
    fmt.Printf("s2: %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
    
    // Method 3: Using make (length and capacity)
    s3 := make([]int, 3, 10)
    fmt.Printf("s3: %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))
    
    // Method 4: Slice from array
    arr := [5]int{10, 20, 30, 40, 50}
    s4 := arr[1:4]
    fmt.Printf("s4: %v, len=%d, cap=%d\n", s4, len(s4), cap(s4))
    
    // Method 5: Nil slice
    var s5 []int
    fmt.Printf("s5: %v, len=%d, cap=%d, isNil=%v\n", s5, len(s5), cap(s5), s5 == nil)
}
```

**Output:**
```
s1: [1 2 3 4 5], len=5, cap=5
s2: [0 0 0 0 0], len=5, cap=5
s3: [0 0 0], len=3, cap=10
s4: [20 30 40], len=3, cap=4
s5: [], len=0, cap=0, isNil=true
```

---

#### Example: Common Slice Operations

```go
package main

import "fmt"

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // 1. Append to end
    numbers = append(numbers, 6)
    fmt.Println("After append:", numbers)
    
    // 2. Prepend to beginning
    numbers = append([]int{0}, numbers...)
    fmt.Println("After prepend:", numbers)
    
    // 3. Insert in middle (at index 3)
    index := 3
    numbers = append(numbers[:index], append([]int{99}, numbers[index:]...)...)
    fmt.Println("After insert at 3:", numbers)
    
    // 4. Delete element at index 4
    index = 4
    numbers = append(numbers[:index], numbers[index+1:]...)
    fmt.Println("After delete at 4:", numbers)
    
    // 5. Remove duplicates (simple version)
    unique := []int{}
    seen := make(map[int]bool)
    for _, num := range numbers {
        if !seen[num] {
            seen[num] = true
            unique = append(unique, num)
        }
    }
    fmt.Println("Unique:", unique)
}
```

**Output:**
```
After append: [1 2 3 4 5 6]
After prepend: [0 1 2 3 4 5 6]
After insert at 3: [0 1 2 99 3 4 5 6]
After delete at 4: [0 1 2 99 4 5 6]
Unique: [0 1 2 99 4 5 6]
```

---

<a name="slices-memory"></a>
### Memory Behavior and Internal Representation

**Slice Internal Structure:**

```
type slice struct {
    ptr *ElementType  // Pointer to underlying array
    len int           // Number of elements in slice
    cap int           // Capacity of underlying array
}
```

**Memory Diagram:**

```
Slice: s := []int{1, 2, 3, 4, 5}

Slice Header (24 bytes on 64-bit):
┌────────────┬─────┬─────┐
│   ptr      │ len │ cap │
│  (8 bytes) │  5  │  5  │
└──────┬─────┴─────┴─────┘
       │
       ↓
Underlying Array (heap):
[1][2][3][4][5]
```

**Memory Example:**

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    s := []int{1, 2, 3, 4, 5}
    
    // Size of slice header (not the data)
    fmt.Printf("Size of slice header: %d bytes\n", unsafe.Sizeof(s))
    
    // Actual memory of elements
    elementSize := unsafe.Sizeof(s[0])
    totalDataSize := int(elementSize) * len(s)
    fmt.Printf("Size of elements: %d bytes (%d elements × %d bytes)\n", 
        totalDataSize, len(s), elementSize)
}
```

**Output:**
```
Size of slice header: 24 bytes
Size of elements: 40 bytes (5 elements × 8 bytes)
```

---

<a name="slices-mistakes"></a>
### Common Mistakes and Pitfalls

#### Mistake 1: Ignoring append() Return Value

```go
package main

import "fmt"

func main() {
    s := []int{1, 2, 3}
    
    // WRONG: Ignoring return value
    append(s, 4)
    fmt.Println("s:", s)  // [1 2 3] - NOT modified!
    
    // CORRECT: Assign return value
    s = append(s, 4)
    fmt.Println("s:", s)  // [1 2 3 4] - modified!
}
```

**Output:**
```
s: [1 2 3]
s: [1 2 3 4]
```

**Lesson:** Always assign the result of `append()` back to the slice variable.

---

#### Mistake 2: Slice Sharing Leads to Unexpected Changes

```go
package main

import "fmt"

func modifySlice(s []int) {
    s[0] = 999
}

func appendToSlice(s []int) []int {
    return append(s, 100)
}

func main() {
    original := []int{1, 2, 3, 4, 5}
    
    // Passing slice to function
    modifySlice(original)
    fmt.Println("After modifySlice:", original)  // [999 2 3 4 5] - MODIFIED!
    
    // append in function doesn't affect original if capacity exceeded
    original2 := []int{1, 2, 3}
    result := appendToSlice(original2)
    fmt.Println("original2:", original2)  // [1 2 3] - unchanged
    fmt.Println("result:", result)        // [1 2 3 100]
}
```

**Output:**
```
After modifySlice: [999 2 3 4 5]
original2: [1 2 3]
result: [1 2 3 100]
```

**Lesson:** Slices are passed by reference, so modifications inside functions affect the original.

---

#### Mistake 3: Nil vs Empty Slice

```go
package main

import "fmt"

func main() {
    // Nil slice
    var nilSlice []int
    
    // Empty slice (non-nil)
    emptySlice := []int{}
    
    // Both have len=0 and cap=0
    fmt.Printf("nilSlice: %v, len=%d, cap=%d, isNil=%v\n", 
        nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)
    fmt.Printf("emptySlice: %v, len=%d, cap=%d, isNil=%v\n", 
        emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)
    
    // Both can be used with append
    nilSlice = append(nilSlice, 1)
    emptySlice = append(emptySlice, 1)
    
    fmt.Println("After append:")
    fmt.Printf("nilSlice: %v\n", nilSlice)
    fmt.Printf("emptySlice: %v\n", emptySlice)
}
```

**Output:**
```
nilSlice: [], len=0, cap=0, isNil=true
emptySlice: [], len=0, cap=0, isNil=false
After append:
nilSlice: [1]
emptySlice: [1]
```

**Lesson:** Nil and empty slices behave the same in most cases, but `nil` is preferred for zero-value slices.

---

#### Mistake 4: Slicing Beyond Capacity

```go
package main

import "fmt"

func main() {
    s := make([]int, 3, 5)
    s[0], s[1], s[2] = 1, 2, 3
    
    fmt.Printf("s: %v, len=%d, cap=%d\n", s, len(s), cap(s))
    
    // This is OK: slicing within capacity
    s2 := s[:cap(s)]
    fmt.Printf("s2: %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
    
    // This will PANIC: slicing beyond capacity
    // s3 := s[:10]  // panic: runtime error: slice bounds out of range
}
```

**Output:**
```
s: [1 2 3], len=3, cap=5
s2: [1 2 3 0 0], len=5, cap=5
```

---

<a name="slices-interview"></a>
### Interview Tips and Key Takeaways

#### 🎯 Frequently Asked Interview Questions

**Q1: What are the three components of a slice?**

**Answer:** 
1. **Pointer** to the underlying array
2. **Length** - number of elements accessible
3. **Capacity** - total size of underlying array from the slice's start

**Q2: What happens when you append to a slice and capacity is exceeded?**

**Answer:** Go allocates a new underlying array (typically double the capacity), copies all existing elements, and then appends the new element(s).

**Q3: Can you compare slices using ==?**

**Answer:** No, slices can only be compared to `nil`. Use `reflect.DeepEqual()` or write a custom comparison function.

**Q4: What's the difference between nil slice and empty slice?**

**Answer:**
- **Nil slice**: `var s []int` - `s == nil` is `true`, no underlying array
- **Empty slice**: `s := []int{}` - `s == nil` is `false`, has an underlying array (even if empty)
- Both have `len=0` and `cap=0`

**Q5: How do you remove an element from a slice?**

**Answer:**
```go
// Remove element at index i
s = append(s[:i], s[i+1:]...)
```

---

#### 🔥 Critical Points for Interviews

1. **Slice is NOT an Array**
   ```go
   var arr [5]int    // Array: fixed size, value type
   var sli []int     // Slice: dynamic size, reference type
   ```

2. **Always Capture append() Result**
   ```go
   s := []int{1, 2, 3}
   s = append(s, 4)  // MUST assign back
   ```

3. **Capacity Growth Strategy**
   - When capacity < 1024: doubles
   - When capacity >= 1024: grows by ~25%

4. **Slices Share Underlying Arrays**
   ```go
   s1 := []int{1, 2, 3, 4, 5}
   s2 := s1[1:4]
   s2[0] = 999
   // s1 is now [1 999 3 4 5]
   ```

5. **Use copy() for Independence**
   ```go
   dst := make([]int, len(src))
   copy(dst, src)
   ```

---

#### ⚠️ Common Interview Pitfalls

| Pitfall | Problem | Solution |
|---------|---------|----------|
| **Not using append() result** | Original slice unchanged | Always: `s = append(s, x)` |
| **Assuming independence** | Slices share arrays | Use `copy()` for separate data |
| **Comparing slices with ==** | Compile error | Use loop or `reflect.DeepEqual()` |
| **Slice beyond bounds** | Runtime panic | Check bounds before slicing |
| **Modifying while iterating** | Unexpected behavior | Copy slice or iterate backwards |

---

#### 💡 Best Practices for Production Code

✅ **DO:**
- Pre-allocate slices when size is known: `make([]int, 0, expectedSize)`
- Use `copy()` when you need independent data
- Always assign result of `append()` back to the variable
- Use `range` for iteration (safe and idiomatic)
- Check for `nil` when appropriate

❌ **DON'T:**
- Don't assume slices are independent after slicing
- Don't ignore the difference between `len` and `cap`
- Don't use arrays when slices are more appropriate
- Don't compare slices with `==`
- Don't modify slices while iterating without care

---

<a name="arrays-vs-slices"></a>
## 3. Arrays vs Slices: Complete Comparison

### Side-by-Side Comparison

| Feature | Array | Slice |
|---------|-------|-------|
| **Size** | Fixed at compile time | Dynamic, can grow/shrink |
| **Type** | Size is part of type: `[3]int` | Size not in type: `[]int` |
| **Passing** | Passed by value (copied) | Passed by reference (no copy) |
| **Zero Value** | Array with zero-valued elements | `nil` |
| **Memory** | Value type | Header + underlying array |
| **Comparison** | Can use `==`, `!=` | Only compare with `nil` |
| **When to use** | Fixed-size, value semantics | Most use cases, flexible |

---

### Detailed Comparison Examples

#### Example 1: Declaration and Initialization

```go
package main

import "fmt"

func main() {
    // ARRAY
    var arr [3]int = [3]int{1, 2, 3}
    arr2 := [...]int{10, 20, 30}  // Size inferred
    
    // SLICE
    var sli []int = []int{1, 2, 3}
    sli2 := []int{10, 20, 30}
    sli3 := make([]int, 3)  // Slice-specific
    
    fmt.Println("Arrays:", arr, arr2)
    fmt.Println("Slices:", sli, sli2, sli3)
}
```

**Output:**
```
Arrays: [1 2 3] [10 20 30]
Slices: [1 2 3] [10 20 30] [0 0 0]
```

---

#### Example 2: Type System

```go
package main

import "fmt"

func main() {
    // Arrays of different sizes are DIFFERENT TYPES
    var a1 [3]int
    var a2 [5]int
    // a1 = a2  // Compile error: cannot assign [5]int to [3]int
    
    fmt.Printf("Type of a1: %T\n", a1)
    fmt.Printf("Type of a2: %T\n", a2)
    
    // Slices of same element type are SAME TYPE
    var s1 []int = make([]int, 3)
    var s2 []int = make([]int, 5)
    s1 = s2  // OK: both are []int
    
    fmt.Printf("Type of s1: %T\n", s1)
    fmt.Printf("Type of s2: %T\n", s2)
}
```

**Output:**
```
Type of a1: [3]int
Type of a2: [5]int
Type of s1: []int
Type of s2: []int
```

---

#### Example 3: Passing to Functions

```go
package main

import "fmt"

func modifyArray(arr [3]int) {
    arr[0] = 999
    fmt.Println("Inside modifyArray:", arr)
}

func modifySlice(sli []int) {
    sli[0] = 999
    fmt.Println("Inside modifySlice:", sli)
}

func main() {
    // ARRAY: passed by value (copied)
    arr := [3]int{1, 2, 3}
    fmt.Println("Before modifyArray:", arr)
    modifyArray(arr)
    fmt.Println("After modifyArray:", arr)  // UNCHANGED
    
    fmt.Println()
    
    // SLICE: passed by reference (no copy)
    sli := []int{1, 2, 3}
    fmt.Println("Before modifySlice:", sli)
    modifySlice(sli)
    fmt.Println("After modifySlice:", sli)  // CHANGED
}
```

**Output:**
```
Before modifyArray: [1 2 3]
Inside modifyArray: [999 2 3]
After modifyArray: [1 2 3]

Before modifySlice: [1 2 3]
Inside modifySlice: [999 2 3]
After modifySlice: [999 2 3]
```

---

#### Example 4: Flexibility

```go
package main

import "fmt"

func main() {
    // ARRAY: size is fixed
    arr := [3]int{1, 2, 3}
    // arr = append(arr, 4)  // Compile error: first argument to append must be slice
    fmt.Println("Array:", arr)
    
    // SLICE: can grow dynamically
    sli := []int{1, 2, 3}
    sli = append(sli, 4)
    sli = append(sli, 5, 6, 7)
    fmt.Println("Slice:", sli)
    
    // Slice can be created from array
    sliceFromArray := arr[:]
    fmt.Println("Slice from array:", sliceFromArray)
}
```

**Output:**
```
Array: [1 2 3]
Slice: [1 2 3 4 5 6 7]
Slice from array: [1 2 3]
```

---

### When to Use Arrays vs Slices

#### ✅ Use Arrays When:

1. **Fixed, known size** at compile time
   ```go
   var rgb [3]uint8  // Always 3 components
   ```

2. **Value semantics needed** (defensive copying)
   ```go
   func process(data [100]byte) {
       // data is a copy, safe to modify
   }
   ```

3. **Performance critical** with small, fixed data
   ```go
   type Vector3D [3]float64  // Math operations
   ```

4. **Avoiding heap allocations**
   ```go
   var buffer [1024]byte  // Stack allocated
   ```

#### ✅ Use Slices When:

1. **Dynamic size** or size unknown at compile time
   ```go
   userInput := make([]string, 0)
   ```

2. **Passing to/from functions** frequently
   ```go
   func processData(items []int) { }
   ```

3. **Need append/delete operations**
   ```go
   items = append(items, newItem)
   ```

4. **Working with collections** in general
   ```go
   names := []string{"Alice", "Bob", "Charlie"}
   ```

---

### Quick Decision Matrix

```
Need fixed size known at compile time? 
├─ Yes → Consider Array
└─ No → Use Slice

Need to pass data to functions?
├─ Frequently → Use Slice (no copy overhead)
└─ Rarely → Either

Need value semantics (defensive copy)?
├─ Yes → Use Array
└─ No → Use Slice

Need to append/remove elements?
├─ Yes → Use Slice
└─ No → Either

In doubt?
└─ Use Slice (99% of Go code uses slices)
```

---

## Summary - Part 1

### Key Takeaways

**Arrays:**
- ✓ Fixed-size value types
- ✓ Size is part of the type
- ✓ Copied when passed or assigned
- ✓ Can be compared with `==`
- ✓ Use for fixed-size, value-semantic data

**Slices:**
- ✓ Dynamic-size reference types
- ✓ Three components: pointer, length, capacity
- ✓ Efficient to pass (no copying)
- ✓ Can share underlying arrays
- ✓ Use `append()` to grow
- ✓ Use `copy()` for independence
- ✓ **Preferred in 99% of cases**

### Interview Readiness Checklist

- [ ] Understand array vs slice differences
- [ ] Know when capacity triggers reallocation
- [ ] Understand underlying array sharing
- [ ] Can explain len vs cap
- [ ] Know how to safely copy slices
- [ ] Understand append() behavior
- [ ] Know slice expressions and 3-index form
- [ ] Can identify common pitfalls
- [ ] Know best practices for production code

---

**Next:** Part 2 will cover **Maps** and **Strings & Runes** 🚀

