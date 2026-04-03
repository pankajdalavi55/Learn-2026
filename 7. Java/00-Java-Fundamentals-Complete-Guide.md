# Java Fundamentals — Complete Interview Guide

## Table of Contents

- [1. Java Overview](#1-java-overview)
  - [1.1 What is Java?](#11-what-is-java)
  - [1.2 Key Features of Java](#12-key-features-of-java)
  - [1.3 JVM, JRE, JDK](#13-jvm-jre-jdk)
  - [1.4 How Java Code Executes](#14-how-java-code-executes)
- [2. Language Fundamentals](#2-language-fundamentals)
  - [2.1 Data Types](#21-data-types)
  - [2.2 Variables](#22-variables)
  - [2.3 Type Conversion and Casting](#23-type-conversion-and-casting)
  - [2.4 Wrapper Classes and Autoboxing](#24-wrapper-classes-and-autoboxing)
  - [2.5 Naming Conventions](#25-naming-conventions)
- [3. Operators and Assignment](#3-operators-and-assignment)
  - [3.1 Arithmetic Operators](#31-arithmetic-operators)
  - [3.2 Relational Operators](#32-relational-operators)
  - [3.3 Logical Operators](#33-logical-operators)
  - [3.4 Bitwise Operators](#34-bitwise-operators)
  - [3.5 Assignment Operators](#35-assignment-operators)
  - [3.6 Operator Precedence and Associativity](#36-operator-precedence-and-associativity)
  - [3.7 Tricky Interview Questions — Operators](#37-tricky-interview-questions--operators)
- [4. Flow Control](#4-flow-control)
  - [4.1 Conditional Statements](#41-conditional-statements)
  - [4.2 Switch Statement and Expressions](#42-switch-statement-and-expressions)
  - [4.3 Loops](#43-loops)
  - [4.4 Break, Continue, and Labels](#44-break-continue-and-labels)
- [5. Declaration and Access Modifiers](#5-declaration-and-access-modifiers)
  - [5.1 Class, Method, and Variable Declarations](#51-class-method-and-variable-declarations)
  - [5.2 Access Modifiers](#52-access-modifiers)
  - [5.3 Non-Access Modifiers](#53-non-access-modifiers)
  - [5.4 Scope and Visibility Rules](#54-scope-and-visibility-rules)
  - [5.5 Common Interview Traps — Modifiers](#55-common-interview-traps--modifiers)
  - [5.6 Interfaces — Complete Guide](#56-interfaces--complete-guide)
- [6. Quick Revision Cheat Sheet](#6-quick-revision-cheat-sheet)
- [7. Interview Questions with Answers](#7-interview-questions-with-answers)

---

## 1. Java Overview

### 1.1 What is Java?

Java is a **class-based, object-oriented, general-purpose** programming language designed to have as few implementation dependencies as possible. It was created by **James Gosling** at **Sun Microsystems** in 1995 (now owned by Oracle).

Its core philosophy is captured in the phrase: **"Write Once, Run Anywhere" (WORA)** — compiled Java code runs on any platform that has a JVM, without recompilation.

> **Real-World Analogy:** Think of Java like a **universal power adapter**. You write your code once (the adapter), and it works on any machine (any country's socket) as long as the JVM (the adapter converter) is installed.

### 1.2 Key Features of Java

| Feature | Description |
|---|---|
| **Platform Independent** | Bytecode runs on any OS via JVM |
| **Object-Oriented** | Everything is modeled as objects (except primitives) |
| **Strongly Typed** | Every variable must have a declared type; strict type checking at compile time |
| **Automatic Memory Management** | Garbage Collector handles deallocation |
| **Multithreaded** | Built-in support for concurrent programming |
| **Secure** | No explicit pointers, bytecode verification, Security Manager |
| **Robust** | Strong type-checking, exception handling, GC prevent crashes |
| **Distributed** | Built-in support for RMI, sockets, HTTP |
| **Architecture Neutral** | Bytecode is independent of processor architecture |

> **Interview Note:** "Platform independence" applies to **bytecode**, not the JVM itself. The JVM is platform-**dependent** — each OS has its own JVM implementation.

### 1.3 JVM, JRE, JDK

```
┌────────────────────────────────────────────────┐
│  JDK (Java Development Kit)                    │
│  ┌──────────────────────────────────────────┐  │
│  │  JRE (Java Runtime Environment)          │  │
│  │  ┌────────────────────────────────────┐  │  │
│  │  │  JVM (Java Virtual Machine)        │  │  │
│  │  │  - Class Loader                    │  │  │
│  │  │  - Bytecode Verifier               │  │  │
│  │  │  - Execution Engine (JIT + GC)     │  │  │
│  │  └────────────────────────────────────┘  │  │
│  │  + Core Libraries (java.lang, etc.)      │  │
│  └──────────────────────────────────────────┘  │
│  + Compiler (javac), Debugger (jdb),           │
│    Archiver (jar), Javadoc, jshell, etc.       │
└────────────────────────────────────────────────┘
```

| Component | Purpose | Contains |
|---|---|---|
| **JVM** | Executes bytecode | Class Loader, Verifier, Execution Engine (Interpreter + JIT), GC |
| **JRE** | Runtime environment to **run** Java programs | JVM + Core class libraries |
| **JDK** | Full development kit to **develop + run** | JRE + Compiler (`javac`) + Dev tools (`jdb`, `jar`, `javadoc`, `jshell`) |

> **Interview Tip:** Since Java 11, Oracle no longer ships a standalone JRE. You install the JDK, which includes everything.

**Key JVM Internals (high-level):**
- **Class Loader Subsystem** — loads `.class` files (Bootstrap → Extension → Application loaders).
- **Runtime Data Areas** — Method Area, Heap, Stack, PC Register, Native Method Stack.
- **Execution Engine** — Interpreter (line-by-line), JIT Compiler (hot-spot optimization), Garbage Collector.

### 1.4 How Java Code Executes

```
  Source Code (.java)
        │
        ▼
  javac (Compiler)
        │
        ▼
  Bytecode (.class)
        │
        ▼
  JVM ──┬── Class Loader (loads .class)
        ├── Bytecode Verifier (security checks)
        └── Execution Engine
              ├── Interpreter (initial execution)
              └── JIT Compiler (optimizes hot methods → native code)
                      │
                      ▼
              Machine Code (OS-specific)
```

**Step-by-step:**
1. **Write** source code in `.java` files.
2. **Compile** using `javac` → produces platform-independent `.class` bytecode.
3. **Class Loader** loads bytecode into JVM memory.
4. **Bytecode Verifier** ensures code doesn't violate access rules or corrupt memory.
5. **Interpreter** starts executing bytecode line-by-line.
6. **JIT Compiler** identifies frequently executed ("hot") code paths and compiles them directly to **native machine code** for faster execution.

> **Common Mistake:** Saying "Java is interpreted." Java is **both compiled and interpreted**. Source → bytecode is compilation; bytecode → machine code involves interpretation and JIT compilation.

#### Java Platform Evolution (Java 9–21+)

| Version | Type | Key Changes |
|---|---|---|
| **Java 8** (2014) | LTS | Lambdas, Stream API, `default`/`static` interface methods, `Optional`, `java.time` |
| **Java 9** (2017) | — | **Module System (JPMS)**, JShell REPL, `private` interface methods, `List.of()` / `Map.of()` |
| **Java 10** (2018) | — | `var` local variable type inference |
| **Java 11** (2018) | LTS | Single-file execution (`java File.java`), `String.strip()/.isBlank()/.lines()`, HTTP Client API, standalone JRE removed |
| **Java 14** (2020) | — | Switch expressions (finalized), **Helpful NullPointerExceptions**, `record` (preview) |
| **Java 15** (2020) | — | **Text Blocks** (finalized), `sealed` classes (preview) |
| **Java 16** (2021) | — | **Records** (finalized), **`instanceof` pattern matching** (finalized) |
| **Java 17** (2021) | LTS | **Sealed classes** (finalized), `strictfp` made default, Security Manager deprecated |
| **Java 21** (2023) | LTS | **Virtual threads**, **sequenced collections**, pattern matching in switch (finalized), record patterns, unnamed classes (preview) |
| **Java 22** (2024) | — | **Unnamed variables `_`** (finalized), statements before `super()` (preview), Stream Gatherers (preview) |
| **Java 23** (2024) | — | Primitive types in patterns (preview), **Markdown Javadoc comments**, ZGC generational by default, Module Import Declarations (preview) |
| **Java 24** (2025) | — | **Stream Gatherers** (finalized), Flexible Constructor Bodies (3rd preview), Primitive patterns in switch (2nd preview), Class-File API (finalized), Security Manager permanently disabled |

```java
// Java 9+: Module System (module-info.java in module root)
module com.myapp.order {
    requires java.sql;             // depends on java.sql module
    exports com.myapp.order.api;   // only this package is visible to other modules
}

// Java 9+: JShell (interactive REPL)
// $ jshell
// jshell> int x = 42;
// jshell> System.out.println(x * 2);
// 84

// Java 11+: Run single-file programs directly
// $ java HelloWorld.java    (no javac needed for single-file programs)

// Java 23+: Module Import Declarations (preview) — import an entire module
import module java.base;  // imports all packages exported by java.base
// No more: import java.util.*; import java.io.*; import java.math.*; etc.
```

> **Interview Note:** The module system (`JPMS`) adds **strong encapsulation** at the package level. Even `public` classes are invisible to other modules unless explicitly `exports`-ed. This is a layer above access modifiers.

> **Release Cadence:** Since Java 10, Oracle releases a new version every **6 months** (March and September). LTS (Long-Term Support) versions come every **2 years** (11, 17, 21, 25).

---

## 2. Language Fundamentals

### 2.1 Data Types

Java is **statically typed** — every variable must have a declared type at compile time.

#### Primitive Types (8 total)

| Type | Size | Default | Range | Example |
|---|---|---|---|---|
| `byte` | 1 byte | `0` | -128 to 127 | `byte b = 100;` |
| `short` | 2 bytes | `0` | -32,768 to 32,767 | `short s = 30000;` |
| `int` | 4 bytes | `0` | -2³¹ to 2³¹-1 (~±2.1B) | `int i = 100000;` |
| `long` | 8 bytes | `0L` | -2⁶³ to 2⁶³-1 | `long l = 100000L;` |
| `float` | 4 bytes | `0.0f` | ~±3.4×10³⁸ (6-7 sig digits) | `float f = 3.14f;` |
| `double` | 8 bytes | `0.0d` | ~±1.7×10³⁰⁸ (15 sig digits) | `double d = 3.14;` |
| `char` | 2 bytes | `'\u0000'` | 0 to 65,535 (Unicode) | `char c = 'A';` |
| `boolean` | JVM-specific | `false` | `true` or `false` | `boolean flag = true;` |

> **Interview Note:** Java `char` is **2 bytes** (UTF-16), not 1 byte like in C/C++. The size of `boolean` is not precisely defined by the spec — the JVM typically uses 1 byte in arrays, and may use 4 bytes (int) for standalone variables.

#### Non-Primitive (Reference) Types

- **Classes** — `String`, `Scanner`, user-defined classes
- **Interfaces** — `Runnable`, `Comparable`
- **Arrays** — `int[]`, `String[]`
- **Enums** — `enum Day { MON, TUE, ... }`

```java
// Reference types store addresses (references), not actual data
String name = "Java";   // name holds a reference to a String object on the heap
int[] arr = {1, 2, 3};  // arr holds a reference to an array object on the heap
```

#### Text Blocks — Multi-line Strings (Java 15)

Before Java 15, multi-line strings were painful. Text blocks fix this.

```java
// Before (Java 14 and earlier)
String json = "{\n" +
              "  \"name\": \"Alice\",\n" +
              "  \"age\": 30\n" +
              "}";

// After (Java 15+ Text Blocks)
String json = """
        {
          "name": "Alice",
          "age": 30
        }
        """;

// SQL example
String query = """
        SELECT e.name, d.department_name
        FROM employees e
        JOIN departments d ON e.dept_id = d.id
        WHERE e.salary > 50000
        ORDER BY e.name
        """;
```

**Text Block Rules:**
- Opening `"""` must be followed by a newline (no content on same line).
- Indentation is stripped based on the position of the closing `"""`.
- Supports `\s` (preserve trailing space) and `\` (line continuation, no newline).

```java
// Line continuation with \
String singleLine = """
        This is actually \
        one single line\
        """;
// Result: "This is actually one single line"

// Preserve trailing whitespace with \s
String padded = """
        name   \s
        age    \s
        """;
```

> **Interview Note:** Text blocks produce a regular `String` object — no new type. They use the same string pool as literals.

#### String API Enhancements (Java 11–21)

| Method | Version | Example | Result |
|---|---|---|---|
| `isBlank()` | 11 | `"  ".isBlank()` | `true` (whitespace-only) |
| `strip()` | 11 | `"  hi  ".strip()` | `"hi"` (Unicode-aware trim) |
| `stripLeading()` | 11 | `"  hi  ".stripLeading()` | `"hi  "` |
| `stripTrailing()` | 11 | `"  hi  ".stripTrailing()` | `"  hi"` |
| `lines()` | 11 | `"a\nb\nc".lines()` | `Stream<String>` of 3 lines |
| `repeat(n)` | 11 | `"ha".repeat(3)` | `"hahaha"` |
| `indent(n)` | 12 | `"hi".indent(4)` | `"    hi\n"` |
| `transform(fn)` | 12 | `"hello".transform(String::toUpperCase)` | `"HELLO"` |
| `formatted(args)` | 15 | `"Hi %s".formatted("Bob")` | `"Hi Bob"` |

```java
// strip() vs trim() — strip is Unicode-aware
String s = "\u2000 hello \u2000";  // \u2000 = Unicode whitespace
s.trim();   // "\u2000 hello \u2000" — trim only handles ASCII ≤ 32
s.strip();  // "hello" — strip handles all Unicode whitespace

// Practical: process multi-line input
String input = "  Alice  \n  Bob  \n  \n  Charlie  ";
List<String> names = input.lines()
    .map(String::strip)
    .filter(s -> !s.isBlank())
    .toList();
// ["Alice", "Bob", "Charlie"]
```

> **Interview Tip:** Always prefer `strip()` over `trim()` in modern Java. `trim()` only removes ASCII whitespace (char ≤ 32), while `strip()` handles all Unicode whitespace characters.

**Primitive vs Reference — Key Differences:**

| Aspect | Primitive | Reference |
|---|---|---|
| Stored in | Stack (local) or as part of object (heap) | Reference on stack, object on heap |
| Default value | `0`, `false`, `'\u0000'` | `null` |
| Can be `null`? | No | Yes |
| Pass by? | Value (copy) | Value of reference (alias to same object) |
| `==` compares | Actual values | Memory addresses |

> **Common Mistake:** Saying "Java passes objects by reference." Java is **always pass-by-value**. For objects, it passes the **value of the reference** (memory address), not the object itself.

#### Helpful NullPointerExceptions (Java 14+)

Before Java 14, NPE messages were useless for chained calls. Now the JVM tells you exactly what was null.

```java
Employee emp = getEmployee();
String city = emp.getAddress().getCity().toUpperCase();
// If address is null...

// Before (Java 13): "NullPointerException" (which part?!)

// After (Java 14+):
// "Cannot invoke Address.getCity() because the return value of
//  Employee.getAddress() is null"
```

> **Interview Note:** This is enabled by default since Java 15. It works by analyzing the bytecode at the point of the NPE. No performance cost until an NPE actually occurs.

### 2.2 Variables

#### Types of Variables

```java
public class Employee {

    // Instance variable — one copy per object, stored on heap
    private String name;
    private int age;

    // Static (class) variable — one copy per class, stored in Method Area
    private static int employeeCount = 0;

    public void printInfo() {
        // Local variable — exists only within this method, stored on stack
        String info = name + " (" + age + ")";
        System.out.println(info);
    }
}
```

| Variable Type | Declared In | Stored In | Lifetime | Default Value |
|---|---|---|---|---|
| **Local** | Method / block / constructor | Stack | Method execution | **None** (must be initialized) |
| **Instance** | Class (outside methods) | Heap (part of object) | Object lifetime | Type default (`0`, `null`, `false`) |
| **Static** | Class with `static` keyword | Method Area | Class lifetime (until ClassLoader unloads) | Type default |

> **Interview Trap:** Local variables do **not** get default values. Using an uninitialized local variable causes a **compile-time error**.

```java
public void test() {
    int x;
    System.out.println(x); // ❌ Compile error: variable x might not have been initialized
}
```

#### `var` (Local Variable Type Inference — Java 10+)

```java
var list = new ArrayList<String>();  // compiler infers ArrayList<String>
var count = 10;                      // compiler infers int
var name = "Java";                   // compiler infers String

// Restrictions:
// var x;              ❌ Cannot infer type without initializer
// var x = null;       ❌ Cannot infer type from null
// var x = {1, 2, 3}; ❌ Cannot infer array type from initializer
// Cannot be used for method parameters, return types, or fields
```

**Where `var` shines vs where to avoid it:**

```java
// ✅ Good — type is obvious from the right-hand side
var users = new ArrayList<User>();
var response = httpClient.send(request, BodyHandlers.ofString());
var entry = Map.entry("key", "value");

// ❌ Bad — type is unclear, hurts readability
var result = service.process(data);  // What type is result?
var x = calculate();                  // What does calculate return?
```

> **Interview Note:** `var` is NOT a keyword — it's a **reserved type name**. You can still name a variable `var` (don't). No runtime impact — bytecode is identical to explicit types.

#### Unnamed Variables `_` (Java 22)

When you must declare a variable but don't need its value, use `_` to signal intent.

```java
// Before (Java 21 and earlier)
try {
    int num = Integer.parseInt(input);
} catch (NumberFormatException e) {   // 'e' is unused but must be named
    System.out.println("Invalid input");
}

// After (Java 22+)
try {
    int num = Integer.parseInt(input);
} catch (NumberFormatException _) {   // clearly signals: exception object not needed
    System.out.println("Invalid input");
}

// Useful in enhanced for loops
for (var _ : collection) {            // only care about iteration count, not the element
    totalIterations++;
}

// Useful in Map iteration when you only need values
map.forEach((_, value) -> process(value));

// Useful in pattern matching
if (obj instanceof Point(var x, _)) {  // don't care about y
    System.out.println("x = " + x);
}
```

> **Interview Tip:** `_` was a valid identifier until Java 8. Java 9 made it a compile warning. Java 22 repurposes it as a special unnamed variable. Old code using `_` as a variable name won't compile on Java 22+.

### 2.3 Type Conversion and Casting

#### Widening (Implicit) — Automatic, no data loss

```
byte → short → int → long → float → double
               char → int
```

```java
int i = 100;
long l = i;       // implicit widening: int → long
float f = l;      // implicit widening: long → float
double d = f;     // implicit widening: float → double
```

#### Narrowing (Explicit) — Manual cast required, potential data loss

```java
double d = 9.78;
int i = (int) d;         // i = 9 (decimal part truncated, NOT rounded)

long big = 130;
byte b = (byte) big;     // b = -126 (overflow! 130 wraps around in byte range)

int large = 100_000;
short s = (short) large;  // s = -31072 (data loss)
```

> **Interview Gotcha:**

```java
byte a = 10;
byte b = 20;
byte c = a + b;         // ❌ Compile error! a + b is promoted to int
byte c = (byte)(a + b); // ✅ Explicit cast needed

// Why? Java promotes byte/short/char to int in arithmetic expressions.
```

**Promotion Rules in Expressions:**
1. `byte`, `short`, `char` → promoted to `int` in any arithmetic operation.
2. If either operand is `long`, the other is promoted to `long`.
3. If either operand is `float`, the other is promoted to `float`.
4. If either operand is `double`, the other is promoted to `double`.

### 2.4 Wrapper Classes and Autoboxing

Each primitive type has a corresponding **wrapper class** in `java.lang`:

| Primitive | Wrapper | Cache Range |
|---|---|---|
| `byte` | `Byte` | -128 to 127 (all) |
| `short` | `Short` | -128 to 127 |
| `int` | `Integer` | -128 to 127 |
| `long` | `Long` | -128 to 127 |
| `float` | `Float` | None |
| `double` | `Double` | None |
| `char` | `Character` | 0 to 127 |
| `boolean` | `Boolean` | `TRUE` and `FALSE` (both cached) |

#### Autoboxing and Unboxing

```java
// Autoboxing: primitive → Wrapper (compiler inserts Integer.valueOf())
Integer obj = 42;          // int → Integer

// Unboxing: Wrapper → primitive (compiler inserts .intValue())
int val = obj;             // Integer → int

// Works in collections
List<Integer> list = new ArrayList<>();
list.add(10);              // autoboxing
int first = list.get(0);  // unboxing
```

> **Critical Interview Topic — Integer Cache:**

```java
Integer a = 127;
Integer b = 127;
System.out.println(a == b);      // true  (same cached object)

Integer c = 128;
Integer d = 128;
System.out.println(c == d);      // false (different objects, outside cache)
System.out.println(c.equals(d)); // true  (compares values)

Integer e = new Integer(127);    // forces new object (deprecated since Java 9)
Integer f = 127;
System.out.println(e == f);      // false (new object bypasses cache)
```

**Why?** `Integer.valueOf()` caches instances for values -128 to 127. Within this range, autoboxing returns the same object. Outside it, new objects are created.

> **Common Mistake — NullPointerException with Unboxing:**

```java
Integer wrapper = null;
int value = wrapper;  // 💥 NullPointerException at runtime!
// Compiler doesn't warn you. The unboxing calls wrapper.intValue() on null.
```

### 2.5 Naming Conventions

| Element | Convention | Example |
|---|---|---|
| **Package** | All lowercase, reverse domain | `com.amazon.inventory` |
| **Class / Interface** | PascalCase (nouns) | `OrderService`, `Serializable` |
| **Method** | camelCase (verbs) | `calculateTotal()`, `isValid()` |
| **Variable** | camelCase (meaningful nouns) | `itemCount`, `userName` |
| **Constant** (`static final`) | SCREAMING_SNAKE_CASE | `MAX_RETRY_COUNT`, `PI` |
| **Enum Values** | SCREAMING_SNAKE_CASE | `Status.ACTIVE` |
| **Type Parameter** | Single uppercase letter | `T`, `E`, `K`, `V` |

**Best Practices:**
- Prefer descriptive names over abbreviations (`customerAddress` > `custAddr`).
- Boolean variables/methods should read like questions: `isEmpty()`, `hasPermission`, `isActive`.
- Avoid Hungarian notation (`strName`, `intAge`) — the type system handles this.
- Constants should **never** be mutable objects declared as `static final`:

```java
// This is NOT a true constant — the list contents can change
public static final List<String> ROLES = new ArrayList<>(); // ❌ Mutable

// Use Collections.unmodifiableList or List.of
public static final List<String> ROLES = List.of("ADMIN", "USER"); // ✅ Immutable
```

---

## 3. Operators and Assignment

### 3.1 Arithmetic Operators

| Operator | Name | Example | Result |
|---|---|---|---|
| `+` | Addition | `10 + 3` | `13` |
| `-` | Subtraction | `10 - 3` | `7` |
| `*` | Multiplication | `10 * 3` | `30` |
| `/` | Division | `10 / 3` | `3` (integer division) |
| `%` | Modulus | `10 % 3` | `1` |

```java
// Integer division truncates
System.out.println(7 / 2);    // 3 (not 3.5)
System.out.println(7.0 / 2);  // 3.5 (one operand is double → double division)
System.out.println(7 / 2.0);  // 3.5

// Modulus with negative numbers — result takes sign of the dividend (left operand)
System.out.println(10 % 3);   //  1
System.out.println(-10 % 3);  // -1
System.out.println(10 % -3);  //  1
System.out.println(-10 % -3); // -1

// Division by zero
System.out.println(10 / 0);     // 💥 ArithmeticException
System.out.println(10.0 / 0);   // Infinity (IEEE 754)
System.out.println(0.0 / 0.0);  // NaN
```

> **Interview Note:** Integer division by zero throws `ArithmeticException`. Floating-point division by zero returns `Infinity` or `NaN` — no exception.

#### Increment / Decrement

```java
int a = 5;
int b = a++;  // b = 5, a = 6 (post-increment: use, then increment)
int c = ++a;  // c = 7, a = 7 (pre-increment: increment, then use)

// Classic interview question
int x = 5;
int y = x++ + ++x;
// Step 1: x++ → uses 5, x becomes 6
// Step 2: ++x → x becomes 7, uses 7
// y = 5 + 7 = 12
```

### 3.2 Relational Operators

| Operator | Meaning | Example |
|---|---|---|
| `==` | Equal to | `5 == 5` → `true` |
| `!=` | Not equal to | `5 != 3` → `true` |
| `>` | Greater than | `5 > 3` → `true` |
| `<` | Less than | `5 < 3` → `false` |
| `>=` | Greater than or equal | `5 >= 5` → `true` |
| `<=` | Less than or equal | `5 <= 3` → `false` |

> **Critical Distinction: `==` vs `.equals()`**

```java
// For primitives: == compares values
int a = 10, b = 10;
System.out.println(a == b); // true

// For objects: == compares references (memory addresses)
String s1 = new String("hello");
String s2 = new String("hello");
System.out.println(s1 == s2);      // false (different objects)
System.out.println(s1.equals(s2)); // true  (same content)

// String pool complicates things
String s3 = "hello";
String s4 = "hello";
System.out.println(s3 == s4);      // true (both point to same pool entry)
```

#### `instanceof` Pattern Matching (Java 16+)

The `instanceof` operator got a major upgrade — cast and bind in one step.

```java
// Before (Java 15 and earlier)
if (obj instanceof String) {
    String s = (String) obj;           // redundant cast
    System.out.println(s.toUpperCase());
}

// After (Java 16+)
if (obj instanceof String s) {        // test + cast + binding in one step
    System.out.println(s.toUpperCase());
}

// Works with logical operators
if (obj instanceof String s && s.length() > 5) {
    System.out.println("Long string: " + s);
}

// Flow-scoping: pattern variable is available after the check
if (!(obj instanceof String s)) {
    return;
}
System.out.println(s.toUpperCase());  // s is in scope here due to flow analysis
```

> **Interview Note:** Pattern variable scope follows **flow analysis**, not block scoping. The variable is available wherever the compiler can prove the `instanceof` check succeeded. This eliminates `ClassCastException` risks from manual casting.

### 3.3 Logical Operators

| Operator | Name | Short-Circuit? | Behavior |
|---|---|---|---|
| `&&` | Logical AND | Yes | Returns `true` if **both** operands are `true` |
| `\|\|` | Logical OR | Yes | Returns `true` if **at least one** operand is `true` |
| `!` | Logical NOT | N/A | Inverts the boolean value |
| `&` | Bitwise AND (on booleans) | **No** | Evaluates both operands always |
| `\|` | Bitwise OR (on booleans) | **No** | Evaluates both operands always |

```java
// Short-circuit evaluation — second operand skipped if result is determined
int x = 0;
if (x != 0 && (10 / x > 2)) {
    // Safe: 10/x is never evaluated because x != 0 is false
}

// Without short-circuit → ArithmeticException
if (x != 0 & (10 / x > 2)) {
    // 💥 Both sides evaluated: 10 / 0 throws exception
}
```

> **Interview Tip:** Always use `&&` and `||` (short-circuit) in conditional logic. Use `&` and `|` only when you explicitly need both sides evaluated (rare).

### 3.4 Bitwise Operators

| Operator | Name | Example (`a = 5` `0101`, `b = 3` `0011`) | Result |
|---|---|---|---|
| `&` | AND | `5 & 3` | `1` (`0001`) |
| `\|` | OR | `5 \| 3` | `7` (`0111`) |
| `^` | XOR | `5 ^ 3` | `6` (`0110`) |
| `~` | NOT (complement) | `~5` | `-6` (`...11111010`) |
| `<<` | Left shift | `5 << 1` | `10` (`1010`) |
| `>>` | Right shift (signed) | `-8 >> 2` | `-2` (preserves sign bit) |
| `>>>` | Right shift (unsigned) | `-1 >>> 28` | `15` (fills with zeros) |

```java
// Practical use: Check if a number is even or odd
boolean isEven = (n & 1) == 0; // faster than n % 2 == 0

// Swap two numbers without temp variable
a = a ^ b;
b = a ^ b;
a = a ^ b;

// Multiply / Divide by powers of 2
int doubled = n << 1;  // n * 2
int halved  = n >> 1;  // n / 2

// Check if number is a power of 2
boolean isPowerOf2 = (n > 0) && ((n & (n - 1)) == 0);

// Set, clear, toggle specific bit
int setBit3   = n | (1 << 3);   // set bit 3
int clearBit3 = n & ~(1 << 3);  // clear bit 3
int toggleBit = n ^ (1 << 3);   // toggle bit 3
```

> **Interview Note:** Left shift by `n` is equivalent to multiplying by 2ⁿ. Right shift by `n` is equivalent to dividing by 2ⁿ (integer division). `>>>` always fills with zeros (unsigned), while `>>` preserves the sign bit.

### 3.5 Assignment Operators

| Operator | Example | Equivalent To |
|---|---|---|
| `=` | `a = 10` | — |
| `+=` | `a += 5` | `a = a + 5` |
| `-=` | `a -= 5` | `a = a - 5` |
| `*=` | `a *= 5` | `a = a * 5` |
| `/=` | `a /= 5` | `a = a / 5` |
| `%=` | `a %= 5` | `a = a % 5` |
| `<<=` | `a <<= 2` | `a = a << 2` |
| `>>=` | `a >>= 2` | `a = a >> 2` |
| `&=` | `a &= mask` | `a = a & mask` |
| `\|=` | `a \|= mask` | `a = a \| mask` |
| `^=` | `a ^= mask` | `a = a ^ mask` |

> **Interview Gotcha — Compound Assignment Includes Implicit Cast:**

```java
byte b = 50;
b = b + 5;    // ❌ Compile error: b + 5 is int, cannot assign to byte
b += 5;       // ✅ Works! Compiler inserts implicit cast: b = (byte)(b + 5)

// This can silently cause overflow
byte x = 127;
x += 1;       // No compile error, but x is now -128 (overflow)
```

### 3.6 Operator Precedence and Associativity

From **highest** to **lowest** precedence:

| Precedence | Operators | Associativity |
|---|---|---|
| 1 | `()` `[]` `.` | Left to right |
| 2 | `++` `--` (postfix) | Left to right |
| 3 | `++` `--` (prefix) `+` `-` (unary) `~` `!` `(type)` | Right to left |
| 4 | `*` `/` `%` | Left to right |
| 5 | `+` `-` | Left to right |
| 6 | `<<` `>>` `>>>` | Left to right |
| 7 | `<` `<=` `>` `>=` `instanceof` | Left to right |
| 8 | `==` `!=` | Left to right |
| 9 | `&` | Left to right |
| 10 | `^` | Left to right |
| 11 | `\|` | Left to right |
| 12 | `&&` | Left to right |
| 13 | `\|\|` | Left to right |
| 14 | `?:` (ternary) | Right to left |
| 15 | `=` `+=` `-=` `*=` `/=` etc. | Right to left |

> **Interview Tip:** When in doubt, use parentheses for clarity. Most production code should not rely on obscure precedence rules.

### 3.7 Tricky Interview Questions — Operators

```java
// Q1: What's the output?
System.out.println(1 + 2 + "3");   // "33" (1+2=3, then 3+"3"="33")
System.out.println("1" + 2 + 3);   // "123" ("1"+2="12", "12"+3="123")

// Q2: What's the output?
int i = 0;
i = i++;
System.out.println(i);  // 0 (post-increment returns old value, then assignment overwrites)

// Q3: Will this compile?
short s = 10;
s = s + 5;     // ❌ Compile error
s += 5;        // ✅ Compiles (implicit cast)

// Q4: What's the output?
System.out.println(10 == 010);  // false! 010 is octal = 8 in decimal

// Q5: Floating-point precision
System.out.println(0.1 + 0.2 == 0.3);  // false (IEEE 754 representation issues)
System.out.println(0.1 + 0.2);          // 0.30000000000000004

// For precise decimal arithmetic, use BigDecimal
BigDecimal a = new BigDecimal("0.1");
BigDecimal b = new BigDecimal("0.2");
System.out.println(a.add(b).equals(new BigDecimal("0.3"))); // true

// Q6: What happens?
double d = 1 / 0;     // 💥 ArithmeticException — integer division happens first
double d = 1.0 / 0;   // Infinity — floating-point division
```

#### Numeric Literal Enhancements (Java 7+)

```java
// Underscores in numeric literals for readability (Java 7+)
int billion = 1_000_000_000;         // much clearer than 1000000000
long hexBytes = 0xFF_EC_DE_5E;
long creditCard = 1234_5678_9012_3456L;
float pi = 3.14_15F;
int binary = 0b1010_0001_0100_0010;  // binary literal (also Java 7+)

// ❌ Invalid positions for underscore:
// int x = _100;       at the beginning
// int x = 100_;       at the end
// float f = 3._14;    adjacent to decimal point
// long l = 0x_1A;     adjacent to 0x prefix
```

#### `Math` API Enhancements (Java 8–9)

```java
// Exact arithmetic — throws ArithmeticException on overflow instead of wrapping silently
Math.addExact(Integer.MAX_VALUE, 1);       // 💥 ArithmeticException
Math.multiplyExact(Integer.MAX_VALUE, 2);  // 💥 ArithmeticException
Math.toIntExact(3_000_000_000L);           // 💥 ArithmeticException (doesn't fit in int)

// Before: silent overflow was a source of subtle bugs
int result = Integer.MAX_VALUE + 1;        // -2147483648 (wraps silently!)

// floorDiv and floorMod — correct behavior for negative numbers
Math.floorDiv(-7, 2);   // -4 (rounds toward negative infinity)
-7 / 2;                  // -3 (rounds toward zero — Java default)

Math.floorMod(-7, 2);   // 1 (always non-negative for positive divisor)
-7 % 2;                  // -1 (takes sign of dividend — Java default)
```

---

## 4. Flow Control

### 4.1 Conditional Statements

#### `if`, `if-else`, `else-if`, nested `if`

```java
// Basic if-else
int score = 85;
if (score >= 90) {
    System.out.println("A");
} else if (score >= 80) {
    System.out.println("B");
} else if (score >= 70) {
    System.out.println("C");
} else {
    System.out.println("F");
}
```

> **Interview Gotcha — Dangling Else:**

```java
// What does this print?
int x = 5;
if (x > 3)
    if (x > 10)
        System.out.println("Greater than 10");
else
    System.out.println("Not greater");   // Prints "Not greater"

// The else binds to the NEAREST if (x > 10), not (x > 3)
// Always use braces to avoid ambiguity
```

#### Ternary Operator

```java
int max = (a > b) ? a : b;

// Nested ternary (avoid in production — poor readability)
String grade = (score >= 90) ? "A" : (score >= 80) ? "B" : "C";
```

### 4.2 Switch Statement and Expressions

#### Traditional Switch (all versions)

```java
int day = 3;
switch (day) {
    case 1:
        System.out.println("Monday");
        break;
    case 2:
        System.out.println("Tuesday");
        break;
    case 3:
        System.out.println("Wednesday");
        break;
    default:
        System.out.println("Other day");
        break;
}
```

**Switch works with:** `byte`, `short`, `char`, `int`, `String` (Java 7+), `enum`.
**Switch does NOT work with:** `long`, `float`, `double`, `boolean`.

> **Common Mistake — Forgetting `break` (Fall-Through):**

```java
int val = 1;
switch (val) {
    case 1:
        System.out.println("One");
        // no break → falls through!
    case 2:
        System.out.println("Two");
        break;
    case 3:
        System.out.println("Three");
}
// Output: "One" then "Two" (fall-through from case 1 to case 2)
```

#### Modern Switch Expression (Java 14+)

```java
// Arrow syntax — no fall-through, no break needed
String dayName = switch (day) {
    case 1 -> "Monday";
    case 2 -> "Tuesday";
    case 3 -> "Wednesday";
    case 4 -> "Thursday";
    case 5 -> "Friday";
    case 6, 7 -> "Weekend";    // multiple labels
    default -> "Invalid";
};

// With blocks — use yield to return a value
String result = switch (status) {
    case "ACTIVE" -> {
        log("Processing active...");
        yield "Proceed";
    }
    case "INACTIVE" -> {
        log("Skipping...");
        yield "Skip";
    }
    default -> "Unknown";
};
```

#### Pattern Matching in Switch (Java 21+)

```java
static String describe(Object obj) {
    return switch (obj) {
        case Integer i when i > 0 -> "Positive integer: " + i;
        case Integer i            -> "Non-positive integer: " + i;
        case String s             -> "String of length " + s.length();
        case null                 -> "null value";
        default                   -> "Something else";
    };
}
```

#### Primitive Types in Patterns (Java 23+ Preview)

Java 23/24 extends pattern matching to **primitive types** — `int`, `long`, `double`, etc. can now be used directly in `instanceof` and `switch` patterns.

```java
// Before (Java 22): only reference types in patterns
if (obj instanceof Integer i) {
    int val = i;  // still need unboxing
}

// After (Java 23+ preview): primitives directly in patterns
if (obj instanceof int i) {    // unboxing + pattern matching in one step
    System.out.println(i * 2);
}

// Primitive patterns in switch
String classify(int statusCode) {
    return switch (statusCode) {
        case 200       -> "OK";
        case 301       -> "Moved";
        case 404       -> "Not Found";
        case int i when i >= 500 -> "Server Error: " + i;  // guarded primitive pattern
        default        -> "Unknown: " + statusCode;
    };
}

// Safe narrowing conversions without explicit casts
long bigValue = 42L;
if (bigValue instanceof int small) {   // true — 42 fits in int
    System.out.println(small);
}

long overflow = 3_000_000_000L;
if (overflow instanceof int small) {   // false — doesn't fit, no ClassCastException
    System.out.println(small);         // not reached
}
```

> **Interview Note:** This is still in **preview** as of Java 24. But it signals Java's direction — making `switch` a universal dispatch mechanism for any value type. Combined with sealed classes and records, this makes Java's pattern matching rival features in Scala and Kotlin.

### 4.3 Loops

#### `for` Loop

```java
for (int i = 0; i < 10; i++) {
    System.out.println(i);
}

// Multiple variables in for loop
for (int i = 0, j = 10; i < j; i++, j--) {
    System.out.println(i + " " + j);
}

// Infinite loop
for (;;) {
    // runs forever until break
}
```

#### `while` Loop

```java
int count = 0;
while (count < 5) {
    System.out.println(count);
    count++;
}
```

#### `do-while` Loop

```java
// Executes at least once, even if condition is false
int x = 10;
do {
    System.out.println(x);
    x++;
} while (x < 5);
// Prints: 10 (executes once, then condition x < 5 is false)
```

> **Interview Tip:** `do-while` is used when the loop body must execute at least once — e.g., menu-driven programs, input validation loops.

#### Enhanced `for` Loop (for-each)

```java
int[] nums = {1, 2, 3, 4, 5};
for (int n : nums) {
    System.out.println(n);
}

List<String> names = List.of("Alice", "Bob", "Charlie");
for (String name : names) {
    System.out.println(name);
}
```

**Limitations of for-each:**
- Cannot access the index.
- Cannot modify the underlying collection while iterating (throws `ConcurrentModificationException`).
- Cannot iterate in reverse.
- Works only on arrays and `Iterable` types.

#### Sequenced Collections — Reversible Iteration (Java 21)

Java 21 introduced `SequencedCollection`, adding first/last access and **reverse iteration** to ordered collections.

```java
// Before (Java 20 and earlier) — inconsistent access across types
List<String> list = List.of("A", "B", "C");
list.get(0);                          // first
list.get(list.size() - 1);           // last (awkward)

Deque<String> deque = new ArrayDeque<>(list);
deque.getFirst();                     // OK
deque.getLast();                      // OK

SortedSet<String> sorted = new TreeSet<>(list);
sorted.first();                       // different method name!
sorted.last();

// After (Java 21) — unified interface for all ordered collections
SequencedCollection<String> seq = new ArrayList<>(List.of("A", "B", "C"));
seq.getFirst();       // "A"
seq.getLast();        // "C"
seq.addFirst("Z");   // adds at beginning
seq.addLast("D");    // adds at end
seq.reversed();      // returns a reversed VIEW (not a copy)

// Reverse iteration is now trivial
for (String s : seq.reversed()) {
    System.out.println(s);  // D, C, B, A, Z
}
```

| Interface | Extends | Key New Methods |
|---|---|---|
| `SequencedCollection<E>` | `Collection<E>` | `getFirst()`, `getLast()`, `addFirst()`, `addLast()`, `reversed()` |
| `SequencedSet<E>` | `SequencedCollection<E>`, `Set<E>` | Same + no duplicates |
| `SequencedMap<K,V>` | `Map<K,V>` | `firstEntry()`, `lastEntry()`, `reversed()`, `pollFirstEntry()` |

> **Interview Note:** `ArrayList`, `LinkedList`, `TreeSet`, `LinkedHashSet`, `TreeMap`, `LinkedHashMap` all now implement the appropriate sequenced interface. `reversed()` returns a **view** — mutations to the view affect the original collection.

#### Stream Gatherers — Custom Intermediate Operations (Java 24)

Before Java 24, the Stream API had a **fixed set** of intermediate operations (`map`, `filter`, `flatMap`, etc.). If you needed windowing, batching, or custom stateful transformations, you were stuck.

**Stream Gatherers** let you write your own intermediate operations.

```java
// Built-in Gatherers (java.util.stream.Gatherers)

// windowFixed(n) — groups elements into fixed-size lists
List<List<Integer>> windows = Stream.of(1, 2, 3, 4, 5, 6, 7)
    .gather(Gatherers.windowFixed(3))
    .toList();
// [[1, 2, 3], [4, 5, 6], [7]]

// windowSliding(n) — sliding window
List<List<Integer>> sliding = Stream.of(1, 2, 3, 4, 5)
    .gather(Gatherers.windowSliding(3))
    .toList();
// [[1, 2, 3], [2, 3, 4], [3, 4, 5]]

// fold — stateful accumulation (like reduce, but as an intermediate op)
Stream.of(1, 2, 3, 4, 5)
    .gather(Gatherers.fold(() -> 0, Integer::sum))
    .toList();
// [15]

// scan — running accumulation (emits every intermediate result)
Stream.of(1, 2, 3, 4, 5)
    .gather(Gatherers.scan(() -> 0, Integer::sum))
    .toList();
// [1, 3, 6, 10, 15]

// mapConcurrent — process elements concurrently with virtual threads
Stream.of("url1", "url2", "url3")
    .gather(Gatherers.mapConcurrent(10, url -> fetchFromApi(url)))
    .toList();
```

```java
// Custom Gatherer: deduplicate consecutive elements
// [1, 1, 2, 2, 2, 3, 1, 1] → [1, 2, 3, 1]

// Before (Java 23 and earlier) — awkward stateful workaround
List<Integer> result = new ArrayList<>();
Integer prev = null;
for (int n : numbers) {
    if (!Objects.equals(n, prev)) { result.add(n); prev = n; }
}

// After (Java 24) — clean, composable, works with parallel streams
List<Integer> result = numbers.stream()
    .gather(Gatherer.ofSequential(
        () -> new Object() { Integer prev = null; },       // initializer (state)
        (state, element, downstream) -> {                   // integrator
            if (!Objects.equals(element, state.prev)) {
                state.prev = element;
                return downstream.push(element);
            }
            return true;
        }
    ))
    .toList();
```

> **Interview Insight:** Gatherers follow the same philosophy as `Collector` (for terminal ops) but for **intermediate** operations. They have 4 components: `initializer`, `integrator`, `combiner` (for parallel), and `finisher`. Use built-in `Gatherers.*` methods when possible.

### 4.4 Break, Continue, and Labels

```java
// break — exits the nearest enclosing loop
for (int i = 0; i < 10; i++) {
    if (i == 5) break;
    System.out.println(i);  // prints 0, 1, 2, 3, 4
}

// continue — skips current iteration
for (int i = 0; i < 10; i++) {
    if (i % 2 == 0) continue;
    System.out.println(i);  // prints 1, 3, 5, 7, 9
}
```

#### Labeled Break and Continue (for nested loops)

```java
outer:
for (int i = 0; i < 5; i++) {
    for (int j = 0; j < 5; j++) {
        if (j == 3) break outer;     // exits BOTH loops
        System.out.println(i + "," + j);
    }
}
// Output: 0,0  0,1  0,2

search:
for (int i = 0; i < matrix.length; i++) {
    for (int j = 0; j < matrix[i].length; j++) {
        if (matrix[i][j] == target) {
            System.out.println("Found at " + i + "," + j);
            break search;
        }
    }
}
```

> **Real-World Analogy:** A labeled `break` is like an emergency exit in a building. Normal `break` exits the current room (inner loop), but a labeled `break` lets you exit the entire building (outer loop) in one step.

---

## 5. Declaration and Access Modifiers

### 5.1 Class, Method, and Variable Declarations

#### Class Declaration

```java
// Top-level class (one public class per file, filename must match)
public class OrderService {
    // fields, constructors, methods
}

// A file can contain multiple non-public top-level classes
class HelperClass { }

// Nested (inner) class
public class Outer {
    class Inner { }              // non-static inner class
    static class StaticNested { } // static nested class
}
```

#### Records — Data Carrier Classes (Java 16)

Records eliminate boilerplate for classes whose sole purpose is holding data.

```java
// Before (Java 15 and earlier) — 30+ lines of boilerplate
public final class Point {
    private final int x;
    private final int y;
    public Point(int x, int y) { this.x = x; this.y = y; }
    public int x() { return x; }
    public int y() { return y; }
    @Override public boolean equals(Object o) { /* ... */ }
    @Override public int hashCode() { return Objects.hash(x, y); }
    @Override public String toString() { return "Point[x=" + x + ", y=" + y + "]"; }
}

// After (Java 16+) — one line
public record Point(int x, int y) { }
// Compiler auto-generates: constructor, x(), y(), equals(), hashCode(), toString()
```

```java
// Compact constructor for validation
public record Employee(int id, String name, double salary) {
    public Employee {   // no parameter list — compact constructor
        if (salary < 0) throw new IllegalArgumentException("Salary cannot be negative");
        name = name.strip();  // can modify before assignment
    }
}

// Records can implement interfaces, have static fields and extra methods
public record Range(int low, int high) implements Comparable<Range> {
    public int span() { return high - low; }

    @Override
    public int compareTo(Range other) { return Integer.compare(this.low, other.low); }
}
```

**Record Restrictions:**
- Cannot extend a class (implicitly extend `java.lang.Record`).
- Implicitly `final` — cannot be extended.
- All components are `private final` — always immutable.
- Cannot add extra instance fields (only the components declared in the header).

> **Interview Tip:** Records are ideal for DTOs, API request/response objects, value objects, map keys, and event payloads. Use them anywhere you'd write a POJO that just holds data.

#### Sealed Classes — Controlled Inheritance (Java 17)

Sealed classes restrict which classes can extend them.

```java
// Before: any class anywhere could extend Shape — no control
// After (Java 17+): only permitted subclasses
public sealed class Shape
    permits Circle, Rectangle, Triangle { }

public final class Circle extends Shape {               // final: cannot extend further
    private final double radius;
    public Circle(double radius) { this.radius = radius; }
}

public non-sealed class Rectangle extends Shape {       // non-sealed: anyone can extend
    private final double width, height;
    public Rectangle(double w, double h) { this.width = w; this.height = h; }
}

public sealed class Triangle extends Shape              // sealed: controlled further
    permits EquilateralTriangle, RightTriangle { }
```

**Subclass must be one of:** `final`, `sealed`, or `non-sealed`.

**Why it matters:** Combined with switch pattern matching, the compiler can verify **exhaustiveness**.

```java
// Compiler knows ALL subtypes → no default needed
double area(Shape shape) {
    return switch (shape) {
        case Circle c    -> Math.PI * c.radius() * c.radius();
        case Rectangle r -> r.width() * r.height();
        case Triangle t  -> computeTriangleArea(t);
    };
}
```

> **Interview Note:** `sealed` and `non-sealed` are **context-sensitive keywords** — they're only keywords when used as modifiers. You can still have a variable named `sealed`. Adding a new permitted subclass forces all exhaustive switches to update — the compiler enforces it.

#### Method Declaration

```java
[access_modifier] [non-access_modifiers] return_type methodName([parameters]) [throws ExceptionList] {
    // body
}

// Examples
public static void main(String[] args) { }
private final int calculate(int a, int b) throws ArithmeticException { return a / b; }
protected abstract void process();
```

#### Flexible Constructor Bodies (Java 22+ Preview)

Traditionally, `super()` or `this()` **must** be the first statement in a constructor. Java 22+ relaxes this — you can put validation and field initialization **before** the super call.

```java
// Before (Java 21 and earlier): workaround needed for validation before super()
public class PositiveAmount extends Amount {
    public PositiveAmount(double value) {
        super(validate(value));   // forced to use static helper method
    }
    private static double validate(double v) {
        if (v <= 0) throw new IllegalArgumentException("Must be positive: " + v);
        return v;
    }
}

// After (Java 22+ preview): statements allowed before super()
public class PositiveAmount extends Amount {
    public PositiveAmount(double value) {
        if (value <= 0) throw new IllegalArgumentException("Must be positive: " + value);
        super(value);   // super() no longer needs to be first statement
    }
}
```

**Rules:** Statements before `super()` cannot reference `this` (the instance being constructed). They can only validate arguments and compute values to pass to the parent constructor.

> **Interview Note:** This is in **3rd preview** as of Java 24. It addresses a long-standing pain point and makes constructors more readable when input validation is needed.

#### Markdown Documentation Comments (Java 23+)

Java 23 allows writing Javadoc in **Markdown** instead of HTML-based `@tags`.

```java
// Before: HTML-based Javadoc
/**
 * Calculates the sum of two numbers.
 * <p>
 * This method uses <b>exact arithmetic</b> and throws
 * {@link ArithmeticException} on overflow.
 *
 * @param a the first number
 * @param b the second number
 * @return the sum
 * @throws ArithmeticException if the result overflows
 */
public int add(int a, int b) { return Math.addExact(a, b); }

// After: Markdown Javadoc (Java 23+)
/// Calculates the sum of two numbers.
///
/// This method uses **exact arithmetic** and throws
/// [ArithmeticException] on overflow.
///
/// @param a the first number
/// @param b the second number
/// @return the sum
/// @throws ArithmeticException if the result overflows
public int add(int a, int b) { return Math.addExact(a, b); }
```

Uses `///` (triple slash) instead of `/** */`. Supports standard Markdown: `**bold**`, `` `code` ``, `[links]`, lists, headings, code blocks. Traditional `@param`, `@return`, `@throws` tags still work inside Markdown comments.

#### Variable Declaration

```java
[access_modifier] [non-access_modifiers] type variableName [= value];

private static final int MAX_SIZE = 100;
public volatile boolean running = true;
transient String sessionToken;
```

### 5.2 Access Modifiers

Java has **four** access levels (from most restrictive to least):

| Modifier | Class | Package | Subclass (same pkg) | Subclass (diff pkg) | World |
|---|---|---|---|---|---|
| `private` | ✅ | ❌ | ❌ | ❌ | ❌ |
| default (no keyword) | ✅ | ✅ | ✅ | ❌ | ❌ |
| `protected` | ✅ | ✅ | ✅ | ✅* | ❌ |
| `public` | ✅ | ✅ | ✅ | ✅ | ✅ |

> \* `protected` access from a different package is only through **inheritance** — you cannot access a protected member via a reference of the parent's type.

```java
// package com.example.base
public class Animal {
    protected String name = "Animal";
    protected void speak() { System.out.println("..."); }
}

// package com.example.derived (DIFFERENT package)
public class Dog extends Animal {
    public void test() {
        System.out.println(this.name);     // ✅ inherited member
        this.speak();                       // ✅ inherited method

        Animal a = new Animal();
        System.out.println(a.name);         // ❌ Compile error!
        // Cannot access protected member via parent-type reference
        // from a different package
    }
}
```

#### Where Can Each Modifier Be Used?

| Modifier | Top-Level Class | Inner Class | Method | Field | Constructor | Local Variable |
|---|---|---|---|---|---|---|
| `public` | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| `protected` | ❌ | ✅ | ✅ | ✅ | ✅ | ❌ |
| default | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| `private` | ❌ | ✅ | ✅ | ✅ | ✅ | ❌ |

> **Interview Trap:** Top-level classes can only be `public` or default (package-private). They **cannot** be `private` or `protected`.

### 5.3 Non-Access Modifiers

#### `static`

```java
public class MathUtils {
    // Static field — shared across all instances
    public static final double PI = 3.14159;

    // Static method — belongs to class, not instances
    public static int add(int a, int b) {
        return a + b;
    }

    // Static block — executed once when class is loaded
    static {
        System.out.println("MathUtils loaded");
    }

    // Static nested class — does not need outer class instance
    static class Helper { }
}

// Usage: no instance needed
MathUtils.add(3, 5);
double pi = MathUtils.PI;
```

**Key Rules:**
- Static methods **cannot** access instance members directly.
- Static methods **cannot** use `this` or `super`.
- Instance methods **can** access static members.

#### `final`

```java
// Final variable — constant (cannot be reassigned)
final int MAX = 100;
MAX = 200; // ❌ Compile error

// Final reference — reference cannot change, but object contents can
final List<String> list = new ArrayList<>();
list.add("hello");        // ✅ modifying contents is allowed
list = new ArrayList<>(); // ❌ reassigning reference is not

// Final method — cannot be overridden by subclasses
public class Parent {
    public final void critical() { }
}
public class Child extends Parent {
    public void critical() { } // ❌ Compile error: cannot override
}

// Final class — cannot be extended
public final class String { }
public class MyString extends String { } // ❌ Compile error

// Final parameter
public void process(final int value) {
    value = 10; // ❌ Compile error
}
```

> **Interview Note:** `final` does NOT make objects immutable. It only prevents reassignment of the reference. For true immutability, all fields must be private and final with no setters, and the class should be final.

#### `abstract`

```java
// Abstract class — cannot be instantiated
public abstract class Shape {
    private String color;

    public Shape(String color) {
        this.color = color;
    }

    // Abstract method — no body, must be overridden
    public abstract double area();

    // Concrete method — has body, can be inherited as-is
    public String getColor() {
        return color;
    }
}

public class Circle extends Shape {
    private double radius;

    public Circle(String color, double radius) {
        super(color);
        this.radius = radius;
    }

    @Override
    public double area() {
        return Math.PI * radius * radius;
    }
}
```

**Key Rules:**
- Abstract classes **can** have constructors (called via `super()` from subclass).
- Abstract classes **can** have concrete methods.
- If a class has **any** abstract method, the class **must** be abstract.
- Abstract methods **cannot** be `private`, `final`, or `static`.
- First **concrete** subclass must implement **all** abstract methods.

#### `synchronized`

```java
// Synchronized method — only one thread can execute at a time (per object instance)
public synchronized void increment() {
    count++;
}

// Synchronized block — finer-grained control
public void process() {
    synchronized (this) {
        // critical section
    }
}
```

#### `volatile`

```java
// Ensures visibility across threads — reads/writes go directly to main memory
private volatile boolean running = true;
```

#### `transient`

```java
// Excluded from serialization
public class User implements Serializable {
    private String username;
    private transient String password; // won't be serialized
}
```

#### `strictfp`

```java
// Ensures IEEE 754 floating-point consistency across platforms
// Rarely used — removed as a meaningful modifier in Java 17 (all FP is now strict)
public strictfp class Calculator { }
```

#### Summary of Non-Access Modifiers

| Modifier | Class | Method | Variable | Version | Description |
|---|---|---|---|---|---|
| `static` | ✅ (nested) | ✅ | ✅ | 1.0 | Belongs to class, not instance |
| `final` | ✅ | ✅ | ✅ | 1.0 | No extension / override / reassignment |
| `abstract` | ✅ | ✅ | ❌ | 1.0 | Incomplete; must be implemented |
| `synchronized` | ❌ | ✅ | ❌ | 1.0 | Thread-safe access |
| `volatile` | ❌ | ❌ | ✅ | 1.0 | Direct main-memory access |
| `transient` | ❌ | ❌ | ✅ | 1.0 | Skip during serialization |
| `native` | ❌ | ✅ | ❌ | 1.0 | Implemented in C/C++ |
| `strictfp` | ✅ | ✅ | ❌ | 1.2 | Strict floating-point (obsolete Java 17+) |
| `default` | ❌ | ✅ (interface) | ❌ | **8** | Concrete method in interface |
| `sealed` | ✅ | ❌ | ❌ | **17** | Restricts which classes can extend |
| `non-sealed` | ✅ | ❌ | ❌ | **17** | Opts out of sealed restriction |

### 5.4 Scope and Visibility Rules

```java
public class ScopeDemo {
    private int instanceVar = 10;       // visible throughout the class
    private static int classVar = 20;   // visible throughout the class

    public void method() {
        int localVar = 30;              // visible only in this method

        {
            int blockVar = 40;          // visible only in this block
            System.out.println(blockVar); // ✅
        }
        // System.out.println(blockVar); // ❌ out of scope

        for (int i = 0; i < 5; i++) {
            int loopVar = i * 2;        // visible only in this loop iteration
        }
        // System.out.println(i);       // ❌ out of scope
        // System.out.println(loopVar); // ❌ out of scope
    }
}
```

**Variable Shadowing:**

```java
public class Shadow {
    int x = 10; // instance variable

    public void demo() {
        int x = 20; // local variable shadows instance variable
        System.out.println(x);      // 20 (local)
        System.out.println(this.x); // 10 (instance)
    }
}
```

> **Best Practice:** Avoid variable shadowing. It hurts readability and is a common source of bugs.

### 5.5 Common Interview Traps — Modifiers

**Trap 1: `abstract` + `final` conflict**

```java
// ❌ Compile error: abstract and final are contradictory
public abstract final class Widget { }

// abstract requires subclassing; final prevents it
```

**Trap 2: `abstract` + `private` conflict**

```java
public abstract class Base {
    private abstract void doWork(); // ❌ Compile error
    // abstract methods must be visible to subclasses for overriding
}
```

**Trap 3: `static` method cannot be `abstract`**

```java
public abstract class Utility {
    public static abstract void process(); // ❌ Compile error
    // static methods belong to class and cannot be overridden
}
```

**Trap 4: Interface methods**

```java
public interface Service {
    // All methods are implicitly public and abstract (before Java 8)
    void execute();

    // Java 8+: default methods have a body
    default void log() { System.out.println("Logging..."); }

    // Java 8+: static methods
    static void version() { System.out.println("v1.0"); }

    // Java 9+: private methods (helper for default methods)
    private void helper() { }

    // Interface fields are implicitly public, static, and final
    int MAX = 100; // same as: public static final int MAX = 100;
}
```

**Interface Evolution Across Java Versions:**

| Version | What Can an Interface Have? |
|---|---|
| Java 1–7 | Abstract methods + `public static final` constants only |
| **Java 8** | + `default` methods (concrete, with body) + `static` methods |
| **Java 9** | + `private` methods + `private static` methods (helpers for default methods) |
| **Java 17** | Interfaces can be `sealed` (restrict implementations) |

```java
// Full modern interface example (Java 17+)
public sealed interface PaymentProcessor permits CreditCardProcessor, UpiProcessor {

    boolean process(double amount);   // abstract

    default void logPayment(double amount) {    // default
        formatAndLog("Payment", amount);
    }

    static PaymentProcessor create(String type) {    // static factory
        return switch (type) {
            case "CARD" -> new CreditCardProcessor();
            case "UPI"  -> new UpiProcessor();
            default -> throw new IllegalArgumentException("Unknown: " + type);
        };
    }

    private void formatAndLog(String action, double amount) {    // private helper
        System.out.printf("[%s] $%.2f%n", action, amount);
    }
}
```

> **Interview Note:** `default` methods were added to allow interface evolution **without breaking existing implementations**. Before Java 8, adding a method to an interface broke every class that implemented it.

**Trap 5: Constructor cannot be `final`, `abstract`, or `static`**

```java
public class Foo {
    public final Foo() { }    // ❌ Compile error
    public abstract Foo() { } // ❌ Compile error
    public static Foo() { }   // ❌ Compile error
}
```

### 5.6 Interfaces — Complete Guide

An interface is a **contract** — it defines *what* a class must do, without dictating *how*.

> **Real-World Analogy:** A power socket is an interface. It defines the shape (contract) — any appliance (implementing class) that matches the plug shape can connect. The socket doesn't care if it's a lamp, a phone charger, or a blender.

#### Declaring and Implementing Interfaces

```java
public interface Drawable {
    void draw();           // implicitly public and abstract
    double area();         // implicitly public and abstract
}

public class Circle implements Drawable {
    private double radius;

    public Circle(double radius) { this.radius = radius; }

    @Override
    public void draw() {
        System.out.println("Drawing circle with radius " + radius);
    }

    @Override
    public double area() {
        return Math.PI * radius * radius;
    }
}
```

**Key Rules:**
- All methods are implicitly `public abstract` (unless `default`, `static`, or `private`).
- All fields are implicitly `public static final` (constants only).
- A class can implement **multiple** interfaces.
- An interface can extend **multiple** interfaces.

```java
// Multiple interface implementation
public class SmartPhone implements Callable, Browsable, Photographable {
    // must implement ALL abstract methods from all three interfaces
}

// Interface extending multiple interfaces
public interface Serializable extends Readable, Writable {
    // inherits abstract methods from both
}
```

#### Implicit Modifiers — What the Compiler Adds

```java
public interface Example {
    int MAX = 100;                  // compiler sees: public static final int MAX = 100
    void process();                 // compiler sees: public abstract void process()
    default void log() { }         // compiler sees: public default void log()
    static void util() { }         // compiler sees: public static void util()
}
```

> **Interview Trap:** You **cannot** declare interface methods as `protected`, `private` (pre-Java 9), `final`, or `synchronized`. Interface fields cannot be non-static or non-final.

```java
public interface Bad {
    protected void doWork();    // ❌ Compile error
    final void execute();       // ❌ Compile error
    private int count = 0;      // ❌ Compile error (fields are always public static final)
}
```

#### `default` Methods (Java 8+)

Allow adding new methods to interfaces **without breaking** existing implementations.

```java
public interface Collection<E> {
    // Existing method — all implementations already have this
    boolean add(E element);

    // Added in Java 8 — existing classes don't need to change
    default Stream<E> stream() {
        return StreamSupport.stream(spliterator(), false);
    }
}
```

**Default Method Conflict — The Diamond Problem:**

```java
interface A {
    default void hello() { System.out.println("A"); }
}

interface B {
    default void hello() { System.out.println("B"); }
}

// Class implements both — which hello()?
class C implements A, B {
    @Override
    public void hello() {
        A.super.hello();    // explicitly choose A's version
        // or B.super.hello();
        // or provide entirely new implementation
    }
}
```

**Resolution Rules (in order):**
1. **Class always wins** — a concrete method in the class/superclass overrides any default.
2. **Most specific interface wins** — if `B extends A`, then `B`'s default overrides `A`'s.
3. **Ambiguity → must override** — if no rule resolves it, the class must explicitly override.

```java
interface A {
    default void greet() { System.out.println("A"); }
}

interface B extends A {
    default void greet() { System.out.println("B"); }
}

class C implements A, B {
    // No override needed — B is more specific than A
    // C.greet() calls B's version
}

class D implements A, B {
    @Override
    public void greet() {
        A.super.greet();   // explicitly call A's version
    }
}
```

#### `static` Methods in Interfaces (Java 8+)

```java
public interface Validator<T> {
    boolean validate(T value);

    // Static factory methods — called via interface name
    static Validator<String> nonEmpty() {
        return s -> s != null && !s.isEmpty();
    }

    static Validator<Integer> positive() {
        return i -> i != null && i > 0;
    }

    // Cannot be overridden — belongs to the interface, not implementations
}

// Usage
Validator<String> v = Validator.nonEmpty();
v.validate("hello"); // true
```

> **Interview Note:** Interface static methods are NOT inherited by implementing classes or sub-interfaces. You must call them via the interface name: `MyInterface.staticMethod()`.

#### `private` Methods in Interfaces (Java 9+)

```java
public interface Logger {
    default void logInfo(String msg) {
        log("INFO", msg);
    }

    default void logError(String msg) {
        log("ERROR", msg);
    }

    // Shared helper — not part of the public contract
    private void log(String level, String msg) {
        System.out.printf("[%s] %s: %s%n", LocalDateTime.now(), level, msg);
    }

    // Private static helper
    private static String format(String template, Object... args) {
        return String.format(template, args);
    }
}
```

**Why?** Avoids code duplication between `default` methods without exposing helper logic in the public API.

#### Functional Interfaces (Java 8+)

An interface with **exactly one abstract method** — can be used with lambda expressions.

```java
@FunctionalInterface
public interface Predicate<T> {
    boolean test(T t);              // the single abstract method (SAM)

    // default and static methods don't count
    default Predicate<T> negate() {
        return t -> !test(t);
    }

    static <T> Predicate<T> isEqual(Object target) {
        return t -> Objects.equals(t, target);
    }
}

// Lambda assigns the implementation of the single abstract method
Predicate<String> isLong = s -> s.length() > 10;
isLong.test("hello");       // false
isLong.negate().test("hi"); // true
```

**Key Built-in Functional Interfaces (`java.util.function`):**

| Interface | Method | Signature | Example Use |
|---|---|---|---|
| `Predicate<T>` | `test` | `T → boolean` | Filtering |
| `Function<T,R>` | `apply` | `T → R` | Mapping/transforming |
| `Consumer<T>` | `accept` | `T → void` | Performing side effects |
| `Supplier<T>` | `get` | `() → T` | Lazy value generation |
| `UnaryOperator<T>` | `apply` | `T → T` | Transforming same type |
| `BinaryOperator<T>` | `apply` | `(T, T) → T` | Reducing two values |
| `BiFunction<T,U,R>` | `apply` | `(T, U) → R` | Two-arg transformation |
| `BiPredicate<T,U>` | `test` | `(T, U) → boolean` | Two-arg filter |

```java
// Lambda and method reference examples
Function<String, Integer> length = String::length;
Consumer<String> printer = System.out::println;
Supplier<List<String>> listMaker = ArrayList::new;
BinaryOperator<Integer> sum = Integer::sum;

// Composing functional interfaces
Function<String, String> trim = String::strip;
Function<String, String> upper = String::toUpperCase;
Function<String, String> pipeline = trim.andThen(upper);
pipeline.apply("  hello  "); // "HELLO"
```

> **Interview Note:** `@FunctionalInterface` is optional — it just instructs the compiler to enforce the single-abstract-method rule. Any interface with exactly one abstract method is functional, annotated or not. Methods inherited from `Object` (like `equals`) don't count toward the SAM.

#### Marker Interfaces

Interfaces with **no methods** — used to tag/mark a class with a capability.

```java
// Built-in marker interfaces
public interface Serializable { }  // marks objects as serializable
public interface Cloneable { }     // marks objects as safe to clone
public interface RandomAccess { }  // marks list as O(1) index access

// Custom marker interface
public interface Auditable { }

public class Order implements Auditable {
    // no methods to implement — just signals to framework that this entity is auditable
}

// Usage: checked at runtime
if (entity instanceof Auditable) {
    auditService.track(entity);
}
```

> **Interview Question:** "Marker interface vs annotation — which is better?" Annotations (`@Auditable`) are more flexible (can carry metadata), but marker interfaces provide **compile-time type safety** — you can declare `void track(Auditable entity)` to enforce at compile time.

#### `Comparable` vs `Comparator`

| Aspect | `Comparable<T>` | `Comparator<T>` |
|---|---|---|
| Package | `java.lang` | `java.util` |
| Method | `compareTo(T o)` | `compare(T o1, T o2)` |
| Defines | **Natural ordering** (single) | **Custom ordering** (multiple) |
| Modifies class? | Yes (implements in the class) | No (external comparator) |
| Use with | `Collections.sort(list)` | `Collections.sort(list, comparator)` |

```java
// Comparable — natural ordering defined inside the class
public class Employee implements Comparable<Employee> {
    private String name;
    private double salary;

    @Override
    public int compareTo(Employee other) {
        return Double.compare(this.salary, other.salary);  // natural order: by salary
    }
}

Collections.sort(employees);  // uses compareTo

// Comparator — external, multiple strategies
Comparator<Employee> byName = Comparator.comparing(Employee::getName);
Comparator<Employee> bySalaryDesc = Comparator.comparingDouble(Employee::getSalary).reversed();
Comparator<Employee> byNameThenSalary = byName.thenComparingDouble(Employee::getSalary);

employees.sort(byNameThenSalary);
```

**Modern Comparator factory methods (Java 8+):**

```java
// Chained comparators — clean and readable
Comparator<Employee> comparator = Comparator
    .comparing(Employee::getDepartment)
    .thenComparing(Employee::getName)
    .thenComparingDouble(Employee::getSalary)
    .reversed();

// Null-safe comparators
Comparator<Employee> nullSafe = Comparator.nullsLast(
    Comparator.comparing(Employee::getName)
);
```

#### Interface vs Abstract Class — Decision Flowchart

```
Need to define a contract?
├── YES → Will unrelated classes implement it?
│         ├── YES → Use INTERFACE
│         └── NO → Do classes need shared state (fields)?
│                   ├── YES → Use ABSTRACT CLASS
│                   └── NO → Use INTERFACE (prefer)
└── NO → Use a regular CLASS
```

> **Modern Java Best Practice:** Default to **interfaces**. Use abstract classes only when you need shared instance state (fields) or constructors. Since Java 8+ with default methods, interfaces can provide behavior too — making abstract classes less necessary.

#### Sealed Interfaces (Java 17+) — Controlling Implementation

```java
// Only these three can implement Result
public sealed interface Result<T> permits Success, Failure, Pending { }

public record Success<T>(T value)    implements Result<T> { }
public record Failure<T>(String err) implements Result<T> { }
public record Pending<T>()           implements Result<T> { }

// Exhaustive switch — no default needed
String describe(Result<?> result) {
    return switch (result) {
        case Success<?> s -> "OK: " + s.value();
        case Failure<?> f -> "Error: " + f.err();
        case Pending<?> p -> "Waiting...";
    };
}
```

---

## 6. Quick Revision Cheat Sheet

### Data Types
- 8 primitives: `byte`(1), `short`(2), `int`(4), `long`(8), `float`(4), `double`(8), `char`(2), `boolean`
- `char` is 2 bytes (UTF-16), not 1 byte
- `boolean` size is JVM-dependent (typically 1 byte in arrays, 4 bytes standalone)
- Default for references: `null`; for numbers: `0`; for boolean: `false`
- Text blocks `"""..."""` for multi-line strings (Java 15+)
- Prefer `strip()` over `trim()` — Unicode-aware (Java 11+)

### Variables
- Local variables: no default value (must initialize)
- Instance variables: get default values
- Static variables: shared across all instances
- `var` (Java 10+): local-only type inference, same bytecode
- `_` (Java 22+): unnamed variable for intentionally unused values

### Type Casting
- `byte + short` → promoted to `int`
- Compound assignment (`+=`) includes implicit cast
- `(int) 9.78` = `9` (truncation, not rounding)
- `instanceof` pattern matching (Java 16+): test + cast + bind in one step

### Operators
- `==` on objects → compares references
- `.equals()` on objects → compares values (if overridden)
- Integer cache: -128 to 127 for `Integer.valueOf()`
- `&&` short-circuits; `&` does not
- `0.1 + 0.2 != 0.3` (floating-point precision)
- `Math.addExact()` / `multiplyExact()` throw on overflow instead of wrapping (Java 8+)

### Flow Control
- `switch` works with: `byte`, `short`, `char`, `int`, `String`, `enum`
- `switch` does NOT work with: `long`, `float`, `double`, `boolean`
- Missing `break` in switch → fall-through
- Switch expressions `->` (Java 14+): no fall-through, returns values, uses `yield`
- Pattern matching in switch (Java 21+): type patterns + guarded `when` clauses
- `do-while` executes at least once
- Labeled `break` exits outer loops
- `SequencedCollection.reversed()` for reverse iteration (Java 21+)

### Declarations & Modifiers
- `private` → `default` → `protected` → `public` (increasing visibility)
- `final` ≠ immutable (contents of a `final` object can still change)
- `abstract` + `final` → ❌ illegal
- `abstract` + `private` → ❌ illegal
- `static` method → can't access instance members or use `this`
- `volatile` → thread visibility; `transient` → skip serialization
- `record` (Java 16+): immutable data carrier, auto-generates constructor/equals/hashCode/toString
- `sealed`/`non-sealed`/`permits` (Java 17+): control class inheritance hierarchy
- Interfaces: `default` methods (8), `private` methods (9), `sealed` (17)

### JVM / JDK / JRE
- JDK ⊃ JRE ⊃ JVM (standalone JRE removed since Java 11)
- JVM is platform-**dependent**; bytecode is platform-**independent**
- Java is **both** compiled (javac) and interpreted (JVM) + JIT compiled
- Module system JPMS (Java 9+): `module-info.java` controls package visibility
- Helpful NullPointerExceptions (Java 14+): tells exactly which reference was null

---

## 7. Interview Questions with Answers

> Questions 1–15 cover core fundamentals. Questions 16–20 cover modern Java enhancements.

### Q1: Why is Java called platform independent?

**Answer:** Java source code is compiled into **bytecode** (`.class` files) by `javac`, which is platform-independent. This bytecode runs on the **JVM**, which is platform-specific. Since every OS has its own JVM implementation, the same bytecode can run anywhere — "Write Once, Run Anywhere."

The key nuance: **bytecode is platform-independent, the JVM is not.**

---

### Q2: Is Java purely object-oriented? Why or why not?

**Answer:** No. Java is **not** purely object-oriented because:
- It has **8 primitive types** (`int`, `boolean`, etc.) that are not objects.
- It supports `static` methods and variables that belong to the class, not instances.
- Wrapper classes partially bridge this gap, but primitives exist for performance.

Languages like Smalltalk and Ruby are considered purely OO because everything is an object.

---

### Q3: What is the difference between `==` and `.equals()` in Java?

**Answer:**
- `==` compares **references** for objects (memory addresses) and **values** for primitives.
- `.equals()` compares **logical content** (if properly overridden; default in `Object` behaves like `==`).

```java
String a = new String("hello");
String b = new String("hello");
a == b;       // false (different objects)
a.equals(b);  // true  (same content)
```

Gotcha: `String` literals from the pool may return `true` for `==` due to interning.

---

### Q4: Explain the Integer cache. Why does `==` sometimes work for `Integer` objects?

**Answer:** `Integer.valueOf()` caches `Integer` objects for values **-128 to 127**. Autoboxing uses `valueOf()`, so within this range, the same cached object is returned.

```java
Integer a = 100, b = 100;
a == b; // true (same cached object)

Integer c = 200, d = 200;
c == d; // false (different objects, outside cache)
```

**Best practice:** Always use `.equals()` for wrapper comparisons.

---

### Q5: Why does `byte a = 10; byte b = 20; byte c = a + b;` fail to compile?

**Answer:** In Java, **all arithmetic expressions with `byte`, `short`, or `char` are promoted to `int`**. So `a + b` evaluates to an `int`, which cannot be assigned to a `byte` without an explicit cast.

Fix: `byte c = (byte)(a + b);`

However, `byte c = 10 + 20;` **does** compile because `10` and `20` are **compile-time constants** — the compiler evaluates `30` and verifies it fits in a byte.

---

### Q6: What is the difference between `final`, `finally`, and `finalize()`?

**Answer:**

| Keyword | Purpose |
|---|---|
| `final` | Modifier — prevents extension (class), overriding (method), reassignment (variable) |
| `finally` | Block — executes after try/catch regardless of exception |
| `finalize()` | Method — called by GC before object destruction (**deprecated since Java 9**, removed in Java 18) |

```java
final int x = 10;        // constant

try {
    riskyOperation();
} catch (Exception e) {
    handle(e);
} finally {
    cleanup();            // always executes
}
```

---

### Q7: Can we have a `static` method in an `abstract` class? Can a `static` method be `abstract`?

**Answer:**
- **Yes**, an abstract class can have static methods. They belong to the class itself and are called via the class name.
- **No**, a static method cannot be abstract. `static` methods are resolved at compile-time (static binding), but `abstract` methods need runtime polymorphism (dynamic binding). These are contradictory.

---

### Q8: What is the difference between `protected` and default (package-private) access?

**Answer:**
- **Default:** Accessible within the **same package** only.
- **Protected:** Accessible within the **same package** AND in **subclasses** of different packages (via inheritance only).

Key trap: A protected member in a different package is accessible only through `this` or subclass reference, **not** through a parent-type reference.

---

### Q9: Why can `+=` compile when `=` with the same expression cannot?

**Answer:** Compound assignment operators (`+=`, `-=`, `*=`, etc.) include an **implicit narrowing cast**.

```java
byte b = 50;
b = b + 5;  // ❌ b + 5 is int
b += 5;     // ✅ equivalent to b = (byte)(b + 5)
```

This is defined in JLS §15.26.2. It can silently cause **overflow without a compile error**.

---

### Q10: What happens when you unbox a `null` wrapper object?

**Answer:** It throws a **`NullPointerException`** at runtime.

```java
Integer val = null;
int x = val; // NPE! Compiler inserts val.intValue() which is called on null
```

The compiler does not warn about this. It's a common source of bugs in collections and Optional-less code.

---

### Q11: Can a `switch` statement work with `long`? Why or why not?

**Answer:** No. `switch` supports `byte`, `short`, `char`, `int`, `String` (Java 7+), and `enum`. It does **not** support `long`, `float`, `double`, or `boolean`.

Reason: The `switch` statement was originally designed around `tableswitch` and `lookupswitch` JVM instructions, which operate on 32-bit `int` values. `long` is 64-bit and would require fundamentally different bytecode.

---

### Q12: What is the difference between `break` and `continue` in Java?

**Answer:**
- `break` — **exits** the nearest enclosing loop entirely.
- `continue` — **skips** the current iteration and jumps to the next iteration.

Both support labeled forms for nested loops:
- `break outer;` exits the outer loop.
- `continue outer;` skips to the next iteration of the outer loop.

---

### Q13: Is Java pass-by-value or pass-by-reference?

**Answer:** Java is **always pass-by-value**.
- For **primitives**, the actual value is copied.
- For **objects**, the **reference (memory address) is copied** — not the object itself.

This means you can modify the object's state through the copied reference, but you **cannot** make the original reference point to a new object.

```java
void change(StringBuilder sb) {
    sb.append(" World");    // ✅ modifies the same object
    sb = new StringBuilder("New"); // ❌ only changes the local copy of the reference
}

StringBuilder s = new StringBuilder("Hello");
change(s);
System.out.println(s); // "Hello World"
```

---

### Q14: What are the rules for method `main()` in Java?

**Answer:** The signature must be:

```java
public static void main(String[] args)
```

- `public` — JVM needs to access it from outside the class.
- `static` — JVM calls it without creating an instance.
- `void` — returns nothing to the JVM.
- `String[] args` — command-line arguments (can also be `String... args` or `String args[]`).

Variations that compile:
- `public static void main(String... args)` — varargs
- `final public static void main(String[] args)` — `final` and reordered modifiers
- `public static void main(String args[])` — C-style array

**Since Java 21 (preview):** Unnamed classes and instance main methods simplify entry points:

```java
void main() {
    System.out.println("Hello!");
}
```

---

### Q15: What is the output of this code?

```java
public class Test {
    static int x = 10;
    static { x += 5; }
    public static void main(String[] args) {
        System.out.println(x);  // ?
        Test t = new Test();
        System.out.println(t.x); // ?
    }
    static { x *= 2; }
}
```

**Answer:**
1. Static variable `x` initialized to `10`.
2. First static block: `x += 5` → `x = 15`.
3. Second static block: `x *= 2` → `x = 30`.
4. `main()` prints `30`, then `t.x` also prints `30` (static variable, shared).

Static blocks execute **in order of appearance** during class loading, **before** `main()`.

---

### Q16: What are Text Blocks? How are they different from regular Strings?

**Answer:** Text blocks (Java 15) are multi-line string literals delimited by `"""`. They produce a regular `String` object — no new type. Key differences:
- Automatically handle newlines (no `\n` needed).
- Strip common leading whitespace based on closing `"""` position.
- Support `\s` (preserve trailing space) and `\` (line continuation).

```java
String json = """
        {"name": "Alice"}
        """;
// Equivalent to: "{\"name\": \"Alice\"}\n"
```

They use the same string pool as regular literals. The compiler processes indentation at compile time — no runtime overhead.

---

### Q17: What is `instanceof` pattern matching? Why was it introduced?

**Answer:** Java 16 enhanced `instanceof` to combine the type check, cast, and variable binding into one expression, eliminating redundant casts.

```java
// Before: check, then cast (same type written twice)
if (obj instanceof String) {
    String s = (String) obj;
    use(s);
}

// After: test + cast + bind
if (obj instanceof String s) {
    use(s);
}
```

The pattern variable uses **flow-scoping** — it's available wherever the compiler can prove the check succeeded. This eliminates `ClassCastException` risks from manual casts.

---

### Q18: What is a Record in Java? Can it replace all classes?

**Answer:** A `record` (Java 16) is an immutable data carrier. The compiler generates the canonical constructor, accessor methods, `equals()`, `hashCode()`, and `toString()`.

```java
record Point(int x, int y) { }
```

**Cannot** replace all classes because records:
- Cannot have extra instance fields beyond the declared components.
- Are implicitly `final` — cannot be extended.
- All fields are `final` — always immutable (no setters).
- Cannot extend another class (implicitly extend `java.lang.Record`).

Records are ideal for DTOs, value objects, and any class whose identity is defined entirely by its data.

---

### Q19: What is the difference between `sealed` and `final` classes?

**Answer:**
- `final` → **no class** can extend it. The hierarchy is completely closed.
- `sealed` → **only named classes** (via `permits`) can extend it. The hierarchy is controlled but not closed.

```java
public final class String { }           // nobody can extend
public sealed class Shape permits Circle, Rectangle { }  // only Circle, Rectangle can extend
```

Sealed classes enable **exhaustive pattern matching** in switch — the compiler knows all possible subtypes. Each permitted subclass must be `final`, `sealed`, or `non-sealed`.

---

### Q20: What does `SequencedCollection` solve? Name the key methods.

**Answer:** Before Java 21, accessing the first/last element was inconsistent across collection types (`get(0)` for List, `getFirst()` for Deque, `first()` for SortedSet). `SequencedCollection` (Java 21) provides a **unified interface** for all ordered collections.

Key methods: `getFirst()`, `getLast()`, `addFirst()`, `addLast()`, `removeFirst()`, `removeLast()`, `reversed()`.

`reversed()` returns a **view** — not a copy. Changes to the view affect the original. `ArrayList`, `LinkedList`, `TreeSet`, `LinkedHashSet`, `LinkedHashMap` all implement it.

---
