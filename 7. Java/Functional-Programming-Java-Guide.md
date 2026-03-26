# Functional Programming in Java – Internals, Performance & Senior-Level Interview Guide

> **Target Audience:** Staff-level and Senior-level Java Engineers (8+ years experience)  
> **Purpose:** Deep-dive reference for advanced interviews and architectural discussions covering lambda internals, streams mechanics, FP patterns, performance tuning, and production considerations.

---

## Table of Contents

1. [Evolution of Functional Programming in Java](#1-evolution-of-functional-programming-in-java)
    - 1.1 Historical Timeline
    - 1.2 Pre-Java 8 Limitations
    - 1.3 Why Java is NOT Purely Functional
    - 1.4 Java vs Scala vs Kotlin — FP Comparison
    - 1.5 Impact on Enterprise Systems
    - **1.6 Core FP Concepts & Terminology — Theory** *(NEW)*
2. [Lambda Expressions – Deep Dive](#2-lambda-expressions--deep-dive)
    - 2.1 Syntax, Type Inference, and Target Typing
    - **2.2 SAM Conversion Rules & Type Theory** *(NEW)*
    - 2.3 Effectively Final Variables & Variable Capture
    - 2.4 Closure Behavior — How Capture Works Internally
    - 2.5 Anonymous Class vs Lambda — Critical Differences
    - 2.6 Internal Implementation — invokedynamic & LambdaMetafactory
    - 2.7 Real Production Use Cases
3. [Functional Interfaces – Internals & Usage](#3-functional-interfaces--internals--usage)
4. [Streams API – Internal Mechanics](#4-streams-api--internal-mechanics)
    - **4.0 Theoretical Foundation — Stream Algebra & Lazy Evaluation** *(NEW)*
    - 4.1 Stream Pipeline Architecture
    - 4.2 Stateless vs Stateful Operations
    - 4.3 Internal Implementation — Sink Chaining
    - 4.4 Spliterator — Source Decomposition
    - 4.5 Parallel Streams — ForkJoinPool
    - 4.6 Memory & GC Impact
5. [Optional – Beyond Null Handling](#5-optional--beyond-null-handling)
    - 5.1 Design Intent
    - **5.2 Optional as a Monad — Formal Laws & Theory** *(NEW)*
    - 5.3 Correct Usage Patterns
    - 5.4 Anti-Patterns
    - 5.5 Performance Concerns
    - 5.6 Serialization Implications
6. [Functional Programming Patterns in Java](#6-functional-programming-patterns-in-java)
    - 6.1 Immutability
    - 6.2 Pure Functions & Side-Effect Minimization
    - 6.3 Function Composition & Higher-Order Functions
    - 6.4 Monadic Thinking
    - **6.5 Either / Try / Result Pattern — Railway-Oriented Programming** *(NEW)*
    - 6.6 Currying & Partial Application (Theory & Limitations)
    - 6.7 Lazy Evaluation Patterns
    - 6.8 Custom Collector Design
7. [Concurrency & Functional Programming](#7-concurrency--functional-programming)
    - **7.0 Theoretical Foundation — Why FP Enables Safe Concurrency** *(NEW)*
    - 7.1 Statelessness & Thread Safety
    - 7.2 Non-Interference Rule
    - 7.3 Reduction — Associativity & Identity
    - 7.4 Side Effects in Parallel Streams
    - 7.5 Memory Visibility Concerns
    - 7.6 ForkJoinPool Tuning
    - 7.7 Production Debugging Scenario
8. [Performance & Optimization](#8-performance--optimization)
9. [Advanced Topics](#9-advanced-topics)
    - 9.1 CompletableFuture — Functional Async Chaining
    - 9.2 Functional Style in Reactive Programming
    - 9.3 Functional Configuration in Spring
    - 9.4 Records & Pattern Matching — FP Alignment
    - 9.5 Virtual Threads & FP Considerations
    - 9.6 Gatherers (Java 22+ Preview)
    - **9.7 Type-Driven Development & Functional Architecture** *(NEW)*
10. [Staff-Level / Senior-Level Interview Questions](#10-staff-level--senior-level-interview-questions)
    - 10A. Internal & JVM Questions (Q1–Q5)
    - 10B. Performance & Debugging Scenarios (Q6–Q10)
    - 10C. Architectural & Design Questions (Q11–Q15)

---

## 1. Evolution of Functional Programming in Java

### 1.1 Historical Timeline

```
1996 ─── Java 1.0  ── Anonymous inner classes (verbose callback pattern)
  │
1998 ─── Java 1.2  ── Collections Framework (Iterator, Comparator)
  │                    └─ All based on interface + anonymous class boilerplate
  │
2004 ─── Java 5    ── Generics, Enhanced for-loop, Iterable
  │                    └─ Type-safe collections, but still imperative iteration
  │
2006 ─── Java 6    ── No FP improvements (focus: JVM performance, scripting)
  │
  │      Meanwhile: Scala (2004), Clojure (2007), Groovy lambdas
  │                 → JVM ecosystem demanded first-class functions
  │
2011 ─── Project Lambda (JSR 335) begins
  │
2014 ─── Java 8    ── PARADIGM SHIFT
  │      • Lambda expressions
  │      • Functional interfaces (@FunctionalInterface)
  │      • Streams API
  │      • Optional
  │      • Default methods in interfaces
  │      • Method references
  │      • java.util.function package
  │      • CompletableFuture
  │
2017 ─── Java 9    ── Reactive Streams (Flow API)
  │      • Stream: takeWhile, dropWhile, ofNullable, iterate(seed, hasNext, next)
  │      • Optional: ifPresentOrElse, or, stream()
  │      • Collectors: flatMapping, filtering
  │      • Process API with CompletableFuture
  │
2018 ─── Java 10   ── var (local variable type inference — aids lambda readability)
  │
2018 ─── Java 11   ── var in lambda parameters: (var x, var y) -> x + y
  │      • String/Predicate utilities (isBlank, lines, strip)
  │      • HttpClient (Builder pattern + CompletableFuture)
  │
2021 ─── Java 16   ── Records (immutable data carriers — aligns with FP values)
  │      • Pattern matching for instanceof
  │
2022 ─── Java 17   ── Sealed classes (algebraic data types)
  │
2023 ─── Java 21   ── Pattern matching for switch (finalized)
  │      • Record patterns (deconstruction)
  │      • Virtual threads (impact on async/FP patterns)
  │      • Sequenced collections (better stream sources)
  │
2024+ ── Future    ── Gatherers (Stream intermediate operation customization)
                      String templates, Structured concurrency
```

### 1.2 Pre-Java 8 Limitations

```java
// Pre-Java 8: Sorting required verbose anonymous class
Collections.sort(employees, new Comparator<Employee>() {
    @Override
    public int compare(Employee a, Employee b) {
        return a.getSalary().compareTo(b.getSalary());
    }
});

// Thread creation: anonymous Runnable
new Thread(new Runnable() {
    @Override
    public void run() {
        processData();
    }
}).start();

// Event handling: layers of anonymous classes
button.addActionListener(new ActionListener() {
    @Override
    public void actionPerformed(ActionEvent e) {
        handleClick(e);
    }
});
```

**Core problems:**
- **Verbosity**: Single-method interfaces required ~5 lines of boilerplate per callback
- **No first-class functions**: Functions couldn't be passed as values — only objects wrapping a single method
- **Imperative-only iteration**: External iteration via `for`/`while` — no declarative data pipeline
- **Mutable-by-default**: No language pressure toward immutability
- **No lazy evaluation**: Every collection operation eager and materialized

### 1.3 Why Java is NOT Purely Functional

| Pure FP Requirement | Java's Reality |
|---------------------|---------------|
| **Immutable by default** | Mutable by default; `final` is opt-in, records (Java 16+) help but aren't universal |
| **No side effects** | Side effects everywhere: I/O, state mutation, exceptions |
| **Functions are values** | Lambdas are syntactic sugar over interfaces, not true function types |
| **Algebraic Data Types** | Sealed classes + records (Java 17+) approximate this, but pattern matching is still evolving |
| **Tail-call optimization** | JVM does NOT support TCO — deep recursion → `StackOverflowError` |
| **Higher-kinded types** | Not supported — can't abstract over `Optional<T>`, `Stream<T>`, `List<T>` generically |
| **Pattern matching** | Progressing (switch patterns in Java 21+) but not as powerful as Scala/Haskell |
| **Lazy evaluation** | Streams are lazy, but the language is eager by default |

### 1.4 Java vs Scala vs Kotlin — FP Comparison (Conceptual)

| Feature | Java | Kotlin | Scala |
|---------|------|--------|-------|
| Lambda syntax | `(x) -> x + 1` | `{ x -> x + 1 }` | `x => x + 1` |
| Function types | No real function types (uses interfaces) | `(Int) -> Int` | `Int => Int` |
| Immutable collections | `List.of()` (Java 9+), unmodifiable wrappers | Default (`listOf`) | Default (`List`) |
| Pattern matching | Switch expressions (Java 21+) | `when` expression | Full match/case |
| Tail recursion | Not optimized | `tailrec` keyword | `@tailrec` |
| Extension functions | Not supported | Native support | Implicit classes |
| Null handling | `Optional` (wrapper object) | Nullable types (`String?`) | `Option[T]` |
| Monadic types | `Optional`, `Stream`, `CompletableFuture` | Coroutines, Flow | For-comprehensions |
| Type inference | Limited (`var`, lambda params) | Extensive | Very extensive |

**Architectural Trade-offs of FP in OOP Java:**

1. **Readability**: FP + OOP mix can confuse teams without FP experience — stream chains with 10 intermediate operations become unreadable
2. **Debugging**: Stack traces through lambda/stream pipelines are cryptic (synthetic method names like `lambda$0`)
3. **Performance**: Streams add overhead (object allocation, virtual dispatch) vs simple loops — measurable in hot paths
4. **Testability**: Pure functions are trivially testable; side-effect-free transformations reduce mock complexity
5. **Parallelism**: Stateless operations can be parallelized safely — but Java's `parallel()` is often misused

### 1.5 Impact on Enterprise Systems

- **Spring Framework**: Adopted `Function<>`, `Supplier<>`, `Consumer<>` extensively — Spring Cloud Function, WebFlux (reactive), functional bean registration
- **Data processing**: Streams replaced verbose Guava/Apache Commons collection utilities
- **API design**: Fluent builder APIs + functional callbacks became the standard (e.g., `HttpClient.newBuilder()...`)
- **Microservices**: `CompletableFuture` chains + reactive streams enabled non-blocking service composition
- **Testing**: Lambdas simplified assertion libraries (AssertJ, Mockito argument matchers, custom validators)

### 1.6 Core FP Concepts & Terminology — Theory

Understanding functional programming theory is essential for senior engineers. These principles underpin every FP feature in Java.

#### First-Class Functions

A language has **first-class functions** if functions can be:
1. Assigned to variables
2. Passed as arguments to other functions
3. Returned as results from other functions
4. Stored in data structures

```java
// Java achieves this via functional interfaces (not true function types):
Function<String, Integer> parse = Integer::parseInt;    // 1. Assign to variable
List<String> result = transform(data, String::trim);    // 2. Pass as argument
Predicate<Order> filter = buildFilter(minAmount);       // 3. Return from function
Map<String, Function<Event, Result>> handlers = Map.of( // 4. Store in structure
    "ORDER", this::handleOrder,
    "PAYMENT", this::handlePayment
);

// Limitation: Java uses nominal typing (interface name matters), not structural typing.
// Two interfaces with identical signatures are NOT interchangeable:
//   Supplier<String> s = ...;
//   Callable<String> c = s;   // ❌ Compilation error — different types
```

#### Referential Transparency

An expression is **referentially transparent** if it can be replaced by its value without changing the program's behavior. This is the formal definition of a **pure function**.

```java
// Referentially transparent (pure):
int square(int x) { return x * x; }
// square(5) can ALWAYS be replaced with 25 — anywhere, any time.

// NOT referentially transparent (impure):
int count = 0;
int increment() { return ++count; }
// increment() returns different values each time — cannot substitute.

// NOT referentially transparent (hidden dependency):
String format(LocalDate date) { 
    return date.format(DateTimeFormatter.ofPattern("dd/MM", Locale.getDefault()));
}
// Result depends on JVM Locale setting — same input, different output on different machines.
```

**Why this matters in practice:**
- Referentially transparent functions can be **memoized** (cached) safely
- They can be **parallelized** without locks — no shared mutable state
- They enable **equational reasoning** — reason about code by substitution
- They make **testing trivial** — no mocks, no setup, just input→output assertions

#### Higher-Order Functions (HOF)

A function that takes a function as a parameter or returns a function as a result.

```
Higher-Order Function Classification:

  ┌──────────────────────────────────────────────────────────────────┐
  │  HOF Type           │ Example                │ Java Equivalent   │
  ├─────────────────────┼────────────────────────┼───────────────────┤
  │ Takes function arg  │ map(f, list)           │ stream.map(f)     │
  │ Returns function    │ compose(f, g) → h      │ f.andThen(g)      │
  │ Both                │ decorator(f) → f'      │ Middleware pattern │
  │ Callback pattern    │ onEvent(handler)       │ button.onClick(h) │
  │ Strategy injection  │ sort(comparator)       │ Collections.sort() │
  └─────────────────────┴────────────────────────┴───────────────────┘

  Classic HOFs in Java streams:
    map     : (Stream<T>, T→R)      → Stream<R>      // transform each
    filter  : (Stream<T>, T→bool)   → Stream<T>      // keep matching
    reduce  : (Stream<T>, T×T→T)    → T              // fold into one value
    flatMap : (Stream<T>, T→Stream<R>) → Stream<R>   // transform and flatten
    sorted  : (Stream<T>, T×T→int)  → Stream<T>      // order by comparator
```

#### Purity vs Side Effects

```
Side Effect Spectrum:

  PURE (no side effects)                    IMPURE (side effects)
  ◄──────────────────────────────────────────────────────────────────►
  │                                                                   │
  Math.max(a, b)     String::trim     System.out.println    DB write
  List.of(1,2,3)     record.name()    log.info(...)         HTTP call
  Optional.map(f)    Predicate.and()  Random.nextInt()      File I/O
  
  Categories of Side Effects:
  ┌──────────────────┬──────────────────────────────────────────────┐
  │ Observable        │ I/O (print, file, network), throwing         │
  │                   │ exceptions, modifying external state          │
  ├──────────────────┼──────────────────────────────────────────────┤
  │ Hidden            │ Reading mutable global state, system clock,  │
  │                   │ random numbers, thread-local values           │
  ├──────────────────┼──────────────────────────────────────────────┤
  │ Benign / Internal │ Lazy initialization, caching/memoization,    │
  │                   │ logging (sometimes considered acceptable)     │
  └──────────────────┴──────────────────────────────────────────────┘
```

**The "Functional Core, Imperative Shell" architecture:**
```
  ┌─────────────────────────────────────────────────────────┐
  │                   Imperative Shell                       │
  │   (I/O, DB, HTTP, logging, config reading)              │
  │                                                          │
  │  ┌───────────────────────────────────────────────────┐  │
  │  │              Functional Core                       │  │
  │  │   (pure transformations, business logic,           │  │
  │  │    validation, calculation — NO side effects)      │  │
  │  │                                                     │  │
  │  │   Input Data ──► Pure Functions ──► Output Data    │  │
  │  └───────────────────────────────────────────────────┘  │
  │                                                          │
  │   Read input ──► Pass to core ──► Write output           │
  └─────────────────────────────────────────────────────────┘

  Benefits:
  • Core is trivially testable (no mocks)
  • Core is parallelizable (no shared state)
  • Side effects are isolated and explicit
  • Easier to reason about and refactor
```

#### Type Theory Basics for Java FP

| FP Type Concept | Java Equivalent | What It Means |
|----------------|-----------------|---------------|
| **Product Type** | `record Point(int x, int y)` | Combines multiple values — AND (x AND y) |
| **Sum Type** | `sealed interface Shape permits Circle, Rect` | One of several alternatives — OR (Circle OR Rect) |
| **Unit Type** | `Void` / `void` | A type with exactly one value (no information) |
| **Bottom Type** | `Nothing` (not in Java) | A type with no values — indicates non-termination or error |
| **Parametric Polymorphism** | `<T> List<T>` | Works for ANY type T without knowing what T is |
| **Ad-hoc Polymorphism** | Method overloading / interfaces | Different behavior for different types |
| **Higher-Kinded Types** | NOT supported in Java | Abstracting over type constructors like `M<_>` |
| **Phantom Types** | `class Id<T> { String value; }` | Type parameter used for compile-time safety, not runtime |

```java
// Product types (AND): Both fields present
record Employee(String name, Department dept) {}  // name AND dept

// Sum types (OR): Exactly one variant active (Java 17+ sealed)
sealed interface PaymentResult permits Success, Failure, Pending {}
record Success(String txnId, BigDecimal amount) implements PaymentResult {}
record Failure(String errorCode, String message) implements PaymentResult {}
record Pending(String txnId, Instant estimatedCompletion) implements PaymentResult {}

// Exhaustive handling guaranteed by compiler:
String describe(PaymentResult result) {
    return switch (result) {
        case Success s -> "Paid: " + s.txnId();
        case Failure f -> "Failed: " + f.errorCode();
        case Pending p -> "Pending until: " + p.estimatedCompletion();
        // No default needed — compiler knows all cases covered
    };
}

// Phantom types for compile-time safety:
class UserId extends TypedId<User> { }   // Can't accidentally pass OrderId where UserId expected
class OrderId extends TypedId<Order> { }

void processUser(UserId id) { ... }
// processUser(new OrderId("123"));  // ❌ Compile error — type safety without runtime cost
```

#### Declarative vs Imperative Paradigms

```
Imperative (HOW):                          Declarative / Functional (WHAT):
─────────────────────                      ────────────────────────────────
List<String> result = new ArrayList<>();   List<String> result = employees.stream()
for (Employee e : employees) {                 .filter(e -> e.getSalary() > 50000)
    if (e.getSalary() > 50000) {               .map(Employee::getName)
        result.add(e.getName());               .sorted()
    }                                          .toList();
}
Collections.sort(result);

│                                           │
├─ Explicit iteration (for loop)            ├─ Describes WHAT, not HOW
├─ Explicit state mutation (list.add)       ├─ No explicit state mutation
├─ Control flow is manual                   ├─ Pipeline of transformations
├─ Easy to introduce bugs (off-by-one)      ├─ Harder to have iteration bugs
├─ Harder to parallelize (shared list)      ├─ Trivially parallelizable (.parallel())
└─ Closer to machine execution model        └─ Closer to mathematical specification
```

**When to prefer imperative in Java:**
- Performance-critical tight loops (nanosecond-sensitive hot paths)
- Complex control flow with multiple early exits / `break` / `continue`
- Algorithms that require mutable counters or indices (e.g., two-pointer technique)
- Code that primarily performs side effects (file I/O, DB writes)

**When to prefer functional/declarative:**
- Data transformation pipelines (filter → map → reduce)
- Composition of independent operations
- When parallelism may be needed in the future
- When the business logic maps naturally to "describe the result"

---

## 2. Lambda Expressions – Deep Dive

### 2.1 Syntax, Type Inference, and Target Typing

```java
// Full syntax
(String s) -> { return s.toUpperCase(); }

// Inferred parameter type (target typing from context)
(s) -> { return s.toUpperCase(); }

// Single parameter — parentheses optional
s -> s.toUpperCase()

// Multiple parameters
(a, b) -> a.compareTo(b)

// No parameters
() -> System.currentTimeMillis()

// Block body (multiple statements)
(x, y) -> {
    long sum = x + y;
    log.debug("Sum: {}", sum);
    return sum;
}
```

**Target Typing:** Lambda's type is determined by the **context** where it's used, not the lambda itself.

```java
// Same lambda, two different target types:
Callable<String> c = () -> "hello";     // Target: Callable<String>
Supplier<String> s = () -> "hello";     // Target: Supplier<String>
// Both are valid — the lambda itself has no intrinsic type

// Ambiguity example — compiler error:
// void process(Callable<String> c) { }
// void process(Supplier<String> s) { }
// process(() -> "hello");  // ERROR: ambiguous — both overloads match
```

### 2.2 SAM Conversion Rules & Type Theory

**SAM (Single Abstract Method) conversion** is the mechanism by which the compiler determines if a lambda can be assigned to a target type.

```
SAM Conversion Rules:

  1. Target type MUST be a functional interface (exactly 1 abstract method)
  2. Lambda parameter types must be compatible with the SAM's parameter types
  3. Lambda return type must be compatible with the SAM's return type
  4. Lambda checked exceptions must be declared by the SAM method
  
  Resolution order when compiler encounters a lambda:
  ┌─────────────────────────────────────────────────────────────┐
  │ 1. Determine TARGET TYPE from context:                       │
  │    - Assignment: Function<String,Integer> f = ...           │
  │    - Method argument: list.stream().map(...)                │
  │    - Return statement: return x -> x + 1;                   │
  │    - Cast: (Predicate<String>) s -> s.isEmpty()             │
  │                                                              │
  │ 2. Check SAM compatibility:                                  │
  │    - Target is functional interface? ✓                       │
  │    - Parameter arity matches? ✓                              │
  │    - Types assignable (covariant return, contravariant args)?│
  │    - Checked exceptions subset of SAM's throws clause? ✓    │
  │                                                              │
  │ 3. Generate invokedynamic call site for the SAM method       │
  └─────────────────────────────────────────────────────────────┘
```

**Intersection types with lambdas — serializable lambdas:**

```java
// Cast to intersection type — lambda implements MULTIPLE interfaces:
Comparator<String> comp = (Comparator<String> & Serializable) 
    (a, b) -> a.compareToIgnoreCase(b);
// This lambda is both Comparator AND Serializable

// Practical use: lambdas in distributed systems (Spark, Hazelcast, Ignite)
// Serializable lambdas use LambdaMetafactory.altMetafactory (more complex)
// ⚠ Fragile: serialized form depends on synthetic method names

// Type inference with var (Java 11+):
var func = (Function<String, Integer> & Serializable) s -> s.length();
// 'var' infers the intersection type — useful for local declarations
```

**Type erasure impact on lambdas:**

```java
// Generics are erased at runtime — affects lambda behavior:
Function<String, Integer> f1 = s -> s.length();
Function<Object, Object> f2 = (Function) f1;  // Unchecked cast — compiles with warning
Object result = f2.apply(42);  // ClassCastException at runtime — 42 is not String

// This is why Stream<T>.toArray() returns Object[], not T[]:
// The runtime cannot create a generic array — type T is erased

// Erasure and method references:
List<String> strings = List.of("a", "b");
List<Object> objects = (List) strings;  // Heap pollution
objects.forEach(System.out::println);   // Works (toString on any Object)
objects.stream().map(String::toUpperCase);  // ClassCastException if list has non-Strings!

// Reification (values know their types): Java primitives, arrays
// Erasure (types lost at runtime): generics, lambdas, functional interfaces
// This is why: list.stream().toArray(String[]::new) needs the constructor reference —
// it's the only way to convey the component type at runtime
```

### 2.3 Effectively Final Variables & Variable Capture

Lambdas can capture variables from the enclosing scope, but those variables must be **effectively final** (assigned exactly once).

```java
// Effectively final — no reassignment after initialization
int threshold = 10;  // effectively final (never reassigned)
Predicate<Integer> isAbove = n -> n > threshold;  // ✅ captures threshold

// NOT effectively final — compilation error
int counter = 0;
// Runnable r = () -> counter++;  // ❌ ERROR: counter is modified

// Workaround: use AtomicInteger or single-element array (mutable container)
AtomicInteger counter = new AtomicInteger(0);
Runnable r = () -> counter.incrementAndGet();  // ✅ reference is effectively final
```

**Why effectively final?**
- Lambda captures the **value** (for primitives) or **reference** (for objects) at capture time
- If the variable could change after capture, the lambda would hold a stale copy → confusing semantics
- Unlike closures in JavaScript/Python, Java copies the value into a synthetic field on the lambda instance

### 2.4 Closure Behavior — How Capture Works Internally

```java
public Function<Integer, Integer> createAdder(int base) {
    // 'base' is captured — copied into the lambda's synthetic class
    return x -> x + base;
}

// After compilation, the lambda internally looks like:
// static int lambda$createAdder$0(int captured_base, int x) {
//     return x + captured_base;
// }
// The captured 'base' is passed as an additional argument at the call site
```

```
Variable Capture Model:

  Enclosing Method Stack Frame:
  ┌──────────────────────┐
  │ base = 5             │──── value copied ────┐
  │ threshold = 10       │──── value copied ──┐ │
  │ list = @ref_0x1234   │──── ref copied ──┐ │ │
  └──────────────────────┘                  │ │ │
                                            ▼ ▼ ▼
  Lambda Instance (or static method args):
  ┌──────────────────────┐
  │ captured_list = @ref_0x1234  │ ← SAME object (shared mutable state!)
  │ captured_threshold = 10      │ ← copy (primitive value)
  │ captured_base = 5            │ ← copy (primitive value)
  └──────────────────────┘

  ⚠ Object references are shared — mutations to list inside lambda
    ARE visible outside, and vice versa. Only the reference is final,
    not the object's contents.
```

### 2.5 Anonymous Class vs Lambda — Critical Differences

| Aspect | Anonymous Class | Lambda Expression |
|--------|----------------|-------------------|
| **Bytecode** | Generates a separate `.class` file (`Outer$1.class`) | No separate class file; uses `invokedynamic` |
| **`this` keyword** | Refers to the anonymous class instance | Refers to the **enclosing class** instance |
| **Object creation** | New object allocated on every instantiation | May be cached/reused (JVM optimization) |
| **State** | Can have fields, multiple methods | Stateless (only captured variables) |
| **Type** | Creates a new type in the type hierarchy | Desugared to a method; type generated at runtime |
| **Overhead** | Class loading + object allocation each time | `invokedynamic` bootstrap once; minimal allocation |
| **Serialization** | Straightforward (implements Serializable) | Possible but fragile (depends on synthetic method names) |

```java
// `this` difference:
class Outer {
    String name = "Outer";
    
    void demo() {
        // Anonymous class: `this` = anonymous instance
        Runnable anon = new Runnable() {
            String name = "Anon";
            public void run() {
                System.out.println(this.name);  // prints "Anon"
            }
        };
        
        // Lambda: `this` = Outer instance
        Runnable lambda = () -> {
            System.out.println(this.name);  // prints "Outer"
        };
    }
}
```

### 2.6 Internal Implementation — invokedynamic & LambdaMetafactory

This is the most critical section for staff-level interviews.

```
Lambda Compilation & Runtime Flow:

  Source Code:                Compiled Bytecode:           Runtime (First Call):
  ┌──────────────┐          ┌───────────────────┐        ┌─────────────────────────┐
  │ list.forEach │          │ invokedynamic     │        │ Bootstrap Method:       │
  │ (s -> print  │  javac   │   #accept         │  JVM   │ LambdaMetafactory       │
  │     (s))     │────────►│                   │───────►│  .metafactory()         │
  └──────────────┘          │ + desugared       │        │                         │
                            │   static method:  │        │ Generates a class at    │
                            │   lambda$0(String) │       │ runtime implementing    │
                            └───────────────────┘        │ Consumer<String>        │
                                                         │                         │
                                                         │ Returns a CallSite      │
                                                         │ (cached for future use) │
                                                         └─────────────────────────┘
```

#### Step-by-Step Internal Process:

**1. Compilation (javac):**
- Lambda body is **desugared** into a private static (or instance) method in the enclosing class
- The lambda call site is replaced with an `invokedynamic` instruction
- No anonymous class file is generated at compile time

```java
// Source:
list.forEach(s -> System.out.println(s));

// javac desugars to:
// 1. A private static method in the same class:
private static void lambda$main$0(String s) {
    System.out.println(s);
}

// 2. An invokedynamic instruction at the call site:
// invokedynamic #accept:()*Consumer  [bootstrap: LambdaMetafactory.metafactory]
```

**2. First Invocation (Bootstrap):**

```java
// JVM calls the bootstrap method: LambdaMetafactory.metafactory()
// Parameters:
//   - MethodHandles.Lookup  caller       (access context)
//   - String                interfaceName ("accept")  
//   - MethodType            factoryType  (()Consumer — how to create the lambda)
//   - MethodType            interfaceMethodType (generic: (Object)void)
//   - MethodHandle          implementation (→ lambda$main$0)
//   - MethodType            dynamicMethodType (specialized: (String)void)
//
// LambdaMetafactory uses ASM to generate a class at runtime:
//
//   final class EnclosingClass$$Lambda$1 implements Consumer<String> {
//       public void accept(String s) {
//           EnclosingClass.lambda$main$0(s);  // delegates to desugared method
//       }
//   }
//
// Returns a ConstantCallSite → cached → subsequent calls skip bootstrap entirely
```

**3. Subsequent Invocations:**
- `invokedynamic` resolves via the cached `CallSite` — effectively zero overhead
- Non-capturing lambdas: reuse a **singleton instance** (no allocation per call)
- Capturing lambdas: allocate a new instance (carries captured values)

#### Bytecode Analysis

```bash
# Compile and examine bytecode
javac LambdaDemo.java
javap -v -p LambdaDemo.class
```

```
// Relevant bytecode output:

// 1. The invokedynamic instruction:
  invokedynamic #2,  0   // InvokeDynamic #0:accept:()Ljava/util/function/Consumer;

// 2. Bootstrap methods table:
BootstrapMethods:
  0: #27 REF_invokeStatic java/lang/invoke/LambdaMetafactory.metafactory
    Method arguments:
      #28 (Ljava/lang/Object;)V                           // erased type
      #29 REF_invokeStatic LambdaDemo.lambda$main$0       // implementation
      #30 (Ljava/lang/String;)V                           // specialized type

// 3. The desugared lambda method:
  private static void lambda$main$0(java.lang.String);
    Code:
       0: getstatic     #3    // System.out
       3: aload_0
       4: invokevirtual #4    // println(String)
       7: return
```

#### Why invokedynamic Instead of Anonymous Classes?

| Concern | Anonymous Class Approach | invokedynamic Approach |
|---------|--------------------------|----------------------|
| **Class files** | One `.class` per lambda → classloading overhead, larger JARs | No compile-time class files |
| **Runtime flexibility** | Locked into bytecode at compile time | JVM can choose optimal strategy at runtime |
| **Memory** | New object instance every time | Non-capturing: singleton; JVM can even intrinsify |
| **Future optimizations** | Requires recompilation | JVM can evolve without changing bytecode (e.g., value types) |
| **Startup** | Load N classes for N lambdas | Generate only when first invoked (lazy) |
| **Metaspace** | N class metadata entries at compile time | Generated classes are lightweight, can be GC'd with classloader |

#### Memory Impact

```
Non-Capturing Lambda (no free variables):
  → Singleton instance — ZERO per-invocation allocation
  → Example: s -> s.toUpperCase()

Capturing Lambda (captures local variables):
  → New instance per invocation (carries captured values)
  → Example:
      int threshold = 10;
      filter(n -> n > threshold)  // captures 'threshold'
  → Generated class has a field: final int arg$1;
  → Each call: new LambdaClass$$Lambda$N(threshold)
  → GC pressure proportional to capture frequency

Instance-Capturing Lambda (captures 'this'):
  → Captures the enclosing instance reference
  → Prevents enclosing object from being GC'd while lambda is alive
  → Common source of memory leaks in long-lived callbacks
```

### 2.7 Real Production Use Cases

```java
// 1. Stream pipeline in data processing service
List<OrderDTO> activeOrders = orders.stream()
    .filter(o -> o.getStatus() == Status.ACTIVE)
    .filter(o -> o.getAmount().compareTo(minAmount) > 0)
    .sorted(Comparator.comparing(Order::getCreatedAt).reversed())
    .map(orderMapper::toDTO)
    .collect(Collectors.toList());

// 2. CompletableFuture composition in async service
CompletableFuture<EnrichedUser> enrichUser(String userId) {
    return userService.fetchAsync(userId)
        .thenCompose(user -> addressService.fetchAsync(user.getAddressId())
            .thenCombine(preferenceService.fetchAsync(userId),
                (address, prefs) -> new EnrichedUser(user, address, prefs)))
        .orTimeout(3, TimeUnit.SECONDS)
        .exceptionally(ex -> EnrichedUser.fallback(userId));
}

// 3. Spring WebFlux reactive handler
@GetMapping("/users/{id}")
public Mono<ResponseEntity<UserDTO>> getUser(@PathVariable String id) {
    return userService.findById(id)
        .map(userMapper::toDTO)
        .map(ResponseEntity::ok)
        .defaultIfEmpty(ResponseEntity.notFound().build());
}

// 4. Custom validation using Predicate composition
Predicate<Transaction> isValid = ((Predicate<Transaction>) Transaction::isNonNull)
    .and(t -> t.getAmount().compareTo(BigDecimal.ZERO) > 0)
    .and(t -> t.getCurrency() != null)
    .and(t -> !t.isFlagged());
    
transactions.stream().filter(isValid).forEach(processor::process);
```

---

## 3. Functional Interfaces – Internals & Usage

### 3.1 @FunctionalInterface — What It Actually Does

```java
@FunctionalInterface
public interface Transformer<T, R> {
    R transform(T input);
    
    // Allowed: default methods
    default <V> Transformer<T, V> andThen(Transformer<R, V> after) {
        return t -> after.transform(this.transform(t));
    }
    
    // Allowed: static methods
    static <T> Transformer<T, T> identity() {
        return t -> t;
    }
    
    // Allowed: java.lang.Object methods (toString, equals, hashCode)
    @Override
    String toString();
    
    // NOT allowed: second abstract method → compilation error
    // R transform2(T input);  // ❌
}
```

**Key facts:**
- `@FunctionalInterface` is **informational** — the compiler enforces the single-abstract-method (SAM) rule only when the annotation is present
- An interface with exactly one abstract method is a functional interface **even without** the annotation
- `default` and `static` methods don't count
- Methods inherited from `Object` (equals, hashCode, toString) don't count

### 3.2 Built-in Functional Interfaces — Complete Reference

```
java.util.function package — 43 interfaces organized by pattern:

  ┌───────────────────────────────────────────────────────────────────┐
  │                    Core Functional Interfaces                      │
  ├──────────────┬──────────────┬──────────────┬──────────────────────┤
  │  Function    │  Predicate   │  Consumer    │  Supplier            │
  │  T → R      │  T → boolean │  T → void    │  () → T              │
  ├──────────────┼──────────────┼──────────────┼──────────────────────┤
  │              │              │              │                      │
  │ andThen()   │ and()        │ andThen()    │ (no composition)     │
  │ compose()   │ or()         │              │                      │
  │             │ negate()     │              │                      │
  ├──────────────┴──────────────┴──────────────┴──────────────────────┤
  │                    Operator Specializations                        │
  ├──────────────┬──────────────────────────────────────────────────── │
  │UnaryOperator │  T → T  (extends Function<T,T>)                   │
  │BinaryOperator│  (T, T) → T  (extends BiFunction<T,T,T>)          │
  ├──────────────┴────────────────────────────────────────────────────┤
  │                    Bi-Arity Variants                               │
  ├──────────────┬──────────────┬─────────────────────────────────────┤
  │ BiFunction   │ BiPredicate  │ BiConsumer                          │
  │ (T,U) → R   │ (T,U) → bool │ (T,U) → void                       │
  ├──────────────┴──────────────┴─────────────────────────────────────┤
  │              Primitive Specializations (avoid autoboxing)          │
  ├───────────────────────────────────────────────────────────────────┤
  │ IntFunction<R>    IntPredicate      IntConsumer    IntSupplier    │
  │ LongFunction<R>   LongPredicate     LongConsumer   LongSupplier  │
  │ DoubleFunction<R> DoublePredicate   DoubleConsumer DoubleSupplier│
  │ IntToLongFunction    IntToDoubleFunction                          │
  │ LongToIntFunction    LongToDoubleFunction                         │
  │ DoubleToIntFunction  DoubleToLongFunction                         │
  │ ToIntFunction<T>     ToLongFunction<T>    ToDoubleFunction<T>     │
  │ ObjIntConsumer<T>    ObjLongConsumer<T>   ObjDoubleConsumer<T>    │
  │ IntUnaryOperator     IntBinaryOperator                            │
  │ LongUnaryOperator    LongBinaryOperator                           │
  │ DoubleUnaryOperator  DoubleBinaryOperator                         │
  └───────────────────────────────────────────────────────────────────┘
```

### 3.3 Method References — Four Types

```java
// Type 1: Static method reference
Function<String, Integer> parser = Integer::parseInt;
// Equivalent: s -> Integer.parseInt(s)

// Type 2: Instance method on a particular object
String prefix = "Hello";
Function<String, String> concat = prefix::concat;
// Equivalent: s -> prefix.concat(s)
// ⚠ Captures 'prefix' reference — instance-capturing

// Type 3: Instance method on an ARBITRARY object of a type
Function<String, String> upper = String::toUpperCase;
// Equivalent: s -> s.toUpperCase()
// First parameter becomes the receiver — no capture

// Type 4: Constructor reference
Supplier<ArrayList<String>> factory = ArrayList::new;
// Equivalent: () -> new ArrayList<>()

Function<Integer, ArrayList<String>> sizedFactory = ArrayList::new;
// Equivalent: size -> new ArrayList<>(size)
// JVM infers which constructor based on target type's method signature
```

**Method Reference Internals:**
- Compiled identically to lambdas — **same `invokedynamic` mechanism**
- Type 1 and Type 3: typically non-capturing (singleton instance)
- Type 2: capturing (holds reference to the bound instance)
- Type 4: non-capturing

### 3.4 Composition & Chaining — Higher-Order Functions

```java
// Function composition (right-to-left: compose, left-to-right: andThen)
Function<String, String> trim = String::trim;
Function<String, String> lower = String::toLowerCase;
Function<String, Integer> length = String::length;

// Pipeline: trim → toLowerCase → length
Function<String, Integer> pipeline = trim.andThen(lower).andThen(length);
int result = pipeline.apply("  HELLO  ");  // 5

// Predicate composition
Predicate<Employee> senior = e -> e.getYearsExp() >= 8;
Predicate<Employee> highPerf = e -> e.getRating() >= 4.5;
Predicate<Employee> engineering = e -> e.getDept() == Dept.ENGINEERING;

Predicate<Employee> promotionEligible = senior.and(highPerf).and(engineering);
Predicate<Employee> needsReview = senior.and(highPerf.negate());

// Consumer chaining
Consumer<Order> validate = orderValidator::validate;
Consumer<Order> enrich   = orderEnricher::enrich;
Consumer<Order> persist  = orderRepository::save;
Consumer<Order> notify   = notificationService::sendConfirmation;

Consumer<Order> processOrder = validate.andThen(enrich).andThen(persist).andThen(notify);
orders.forEach(processOrder);
```

### 3.5 Designing Custom Functional Interfaces

```java
// When to create custom vs use built-in:
// ✅ Custom: domain-specific name improves readability
// ✅ Custom: need checked exception support
// ✅ Custom: need more than 2 parameters
// ❌ Custom: if Function/Predicate/Consumer/Supplier fits naturally

// Custom interface with checked exception support
@FunctionalInterface
public interface ThrowingFunction<T, R, E extends Exception> {
    R apply(T t) throws E;
    
    // Adapter: convert to unchecked Function
    static <T, R, E extends Exception> Function<T, R> unchecked(
            ThrowingFunction<T, R, E> f) {
        return t -> {
            try {
                return f.apply(t);
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        };
    }
}

// Usage in streams (checked exceptions + streams is a common pain point):
List<Config> configs = paths.stream()
    .map(ThrowingFunction.unchecked(Files::readString))  // IOException handled
    .map(ConfigParser::parse)
    .collect(toList());

// Tri-function (not in JDK):
@FunctionalInterface
public interface TriFunction<A, B, C, R> {
    R apply(A a, B b, C c);
}
```

### 3.6 Anti-Patterns & Best Practices

**Anti-Patterns:**

```java
// ❌ Lambda that's too complex — should be a named method
items.stream()
    .filter(item -> {
        if (item.getType() == Type.A) {
            return item.getScore() > 50 && item.isActive();
        } else if (item.getType() == Type.B) {
            return item.getScore() > 30;
        }
        return false;
    })
    .collect(toList());

// ✅ Extract to named method for readability and testability
items.stream()
    .filter(this::isEligible)  // Clear intent, testable independently
    .collect(toList());

// ❌ Side effects in Predicate/Function (breaks FP contract)
Predicate<Order> isValid = order -> {
    auditLog.record(order);  // SIDE EFFECT in a predicate!
    return order.isValid();
};

// ❌ Overusing functional chaining when imperative is clearer
// 3+ levels of nested flatMap/compose typically needs refactoring

// ❌ Ignoring primitive specializations
stream.map(i -> i * 2)     // autoboxes int → Integer at each step
stream.mapToInt(i -> i * 2) // ✅ stays in primitive land
```

**Best Practices for Large Systems:**
1. **Name complex lambdas**: Extract to private methods or Predicate/Function constants
2. **Use primitive specializations** in hot paths to avoid autoboxing
3. **Limit chain depth**: If a stream pipeline exceeds 6-7 operations, consider splitting
4. **Document functional parameter contracts**: What are the invariants? Can it return null?
5. **Prefer method references** when they're equally readable — less syntactic noise


---

## 4. Streams API – Internal Mechanics

### 4.0 Theoretical Foundation — Stream Algebra & Lazy Evaluation

Before diving into implementation, understanding the **mathematical foundation** of streams helps reason about correctness and optimization.

#### Streams as Lazy Lists (Conceptual Model)

In FP theory, a stream is a **lazily evaluated sequence** — elements are computed only when demanded by a terminal operation. This connects Java Streams to fundamental CS concepts:

```
Theoretical Model:

  Eager List:   [1, 2, 3, 4, 5]  — all elements in memory at construction
  Lazy Stream:  1 → ? → ? → ?    — each element computed on demand
  
  Category Theory view:
    Stream<T> is a FUNCTOR — supports map: (T→R) → (Stream<T> → Stream<R>)
    Stream<T> is a MONAD  — supports flatMap: (T→Stream<R>) → (Stream<T> → Stream<R>)
    
  This gives us mathematical LAWS that Java streams obey:
  
  Functor Laws:
    1. Identity:     stream.map(x -> x)  ≡  stream
    2. Composition:  stream.map(f).map(g) ≡  stream.map(f.andThen(g))
    
  Monad Laws:
    1. Left identity:   Stream.of(x).flatMap(f)        ≡  f.apply(x)
    2. Right identity:   stream.flatMap(Stream::of)     ≡  stream
    3. Associativity:   stream.flatMap(f).flatMap(g)    ≡  stream.flatMap(x -> f(x).flatMap(g))
```

**Why these laws matter practically:**
- **Functor composition law** → JVM CAN fuse `map(f).map(g)` into `map(f.andThen(g))` — one pass instead of two
- **Monad associativity** → Nested `flatMap` calls can be reordered/flattened without changing results (enables optimizer freedom)
- If these laws were violated, stream pipelines would behave unpredictably when the JVM optimizes

#### Operation Fusion Theory

**Operation fusion** (also called **loop fusion** or **stream fusion**) is the core optimization that makes streams practical:

```
Without Fusion (hypothetical eager evaluation):

  list.stream()
      .filter(predicate)    →  [intermediate List<T>]       // N elements → allocate
      .map(function)        →  [intermediate List<R>]       // M elements → allocate
      .collect(toList())    →  [final List<R>]              // Final → allocate
  
  Cost: 3 iterations + 2 intermediate collections + 3 allocations

With Fusion (Java's actual behavior):

  for (T element : source) {
      if (predicate.test(element)) {            // filter (inline)
          R result = function.apply(element);   // map (inline)
          finalList.add(result);                // collect (inline)
      }
  }
  
  Cost: 1 iteration + 0 intermediate collections + 1 allocation
  
  This is possible because:
  1. ALL stateless operations process one element at a time (Sink interface)
  2. Operations are composed via Sink DELEGATION (not buffering)
  3. JIT compiler can inline the Sink chain into a tight loop
```

**Fusion boundaries — operations that BREAK fusion:**

| Operation | Breaks Fusion? | Why |
|-----------|---------------|-----|
| `filter`, `map`, `peek` | No | Stateless — element-at-a-time processing |
| `flatMap` | Partial | Creates sub-stream per element — local fusion within sub-stream |
| `sorted` | **Yes** | Must buffer ALL elements before emitting ANY |
| `distinct` | **Yes** | Must track all seen elements (HashSet internally) |
| `limit`, `skip` | Partial | Stateful counter, but doesn't buffer elements |
| `collect`, `reduce` | Terminal | Drives the fused pipeline |

#### Short-Circuit Mechanics

Short-circuiting is a form of **lazy evaluation** where the pipeline stops processing once the result is determined:

```
Short-Circuit Flow (findFirst):

  Source: [A, B, C, D, E, F, G, ...]    (could be infinite)
  
  Pipeline: .filter(expensive_predicate).findFirst()
  
  Execution:
    A → filter → REJECT
    B → filter → REJECT
    C → filter → PASS → findFirst receives C → CANCEL PIPELINE
    D, E, F, G ... → NEVER PROCESSED
    
  Implemented via Sink.cancellationRequested():
    For each element:
      if (downstream.cancellationRequested()) → stop iterating source
      
  This enables streams over INFINITE sources:
    Stream.iterate(0, i -> i + 1)   // Infinite!
        .filter(i -> i % 7 == 0)
        .limit(10)                   // Short-circuits after 10 matches
        .toList();                   // [0, 7, 14, 21, ..., 63]
```

#### Push vs Pull Model

```
Pull Model (Iterator / traditional):
  Consumer PULLS elements from source on demand.
  
  while (iterator.hasNext()) {     // Consumer controls pace
      T element = iterator.next(); // Consumer pulls
      process(element);
  }

Push Model (Stream Sink chain / Reactive):
  Source PUSHES elements into the Sink chain.
  
  spliterator.forEachRemaining(element -> {  // Source drives iteration
      sink.accept(element);                  // Source pushes
  });

Java Streams use INTERNAL ITERATION (push):
  ✓ Enables operation fusion (sink chain)
  ✓ Enables parallel splitting (source controls decomposition)
  ✓ Enables short-circuit optimization (source checks cancellation flag)
  ✗ Consumer cannot control pace (no backpressure)
  
  Reactive Streams (Reactor/RxJava) add BACKPRESSURE:
  ✓ Subscriber requests N elements → Publisher pushes at most N
  ✓ Prevents overwhelm when producer is faster than consumer
  ✗ More complex API (Publisher/Subscriber/Subscription)
```

### 4.1 Stream Pipeline Architecture

```
Stream Pipeline Structure:

  ┌──────────┐     ┌──────────────────────────────┐     ┌──────────────┐
  │  Source   │────►│  Intermediate Operations     │────►│   Terminal   │
  │          │     │  (lazy — nothing happens yet) │     │  Operation   │
  └──────────┘     └──────────────────────────────┘     │  (triggers)  │
                                                         └──────────────┘
  Sources:            Intermediate:                  Terminal:
  .stream()           .filter()   (stateless)       .collect()
  .parallelStream()   .map()      (stateless)       .forEach()
  Stream.of()         .flatMap()  (stateless)       .reduce()
  Stream.generate()   .peek()     (stateless)       .count()
  Stream.iterate()    .sorted()   (stateful)        .findFirst()
  Arrays.stream()     .distinct() (stateful)        .toArray()
  Files.lines()       .limit()    (short-circuit)   .min() / .max()
  BufferedReader       .skip()     (stateful)        .anyMatch() (short-circuit)
    .lines()          .takeWhile()(short-circuit)    .allMatch() (short-circuit)
                      .dropWhile()(stateful)         .noneMatch()(short-circuit)
                                                     .findAny()  (short-circuit)
                                                     .toList()   (Java 16+)
```

**Lazy Evaluation — Nothing Happens Until Terminal:**

```java
// This does NOTHING — no filter, no map, no print:
Stream<String> lazy = names.stream()
    .filter(n -> {
        System.out.println("filtering: " + n);  // NEVER printed
        return n.length() > 3;
    })
    .map(String::toUpperCase);
// 'lazy' is just a pipeline description, not an executed computation

// Only when a terminal operation is called:
List<String> result = lazy.collect(toList());  // NOW filter+map execute
```

### 4.2 Stateless vs Stateful Operations

| Category | Operations | Internal Behavior | Parallel Impact |
|----------|-----------|-------------------|-----------------|
| **Stateless** | `filter`, `map`, `flatMap`, `peek`, `mapToInt` | Each element processed independently; no buffers | Perfectly parallelizable |
| **Stateful** | `sorted`, `distinct`, `limit`, `skip`, `dropWhile` | Must see some/all elements before emitting | Require synchronization barriers; limit parallelism |
| **Short-Circuit** | `limit`, `findFirst`, `anyMatch`, `takeWhile` | Can finish without consuming all elements | Cancels remaining splits |

```java
// Performance trap: stateful operation in parallel stream
List<String> result = hugeList.parallelStream()
    .filter(s -> s.startsWith("A"))     // stateless — perfect
    .sorted()                            // stateful — MUST synchronize
    .limit(100)                          // short-circuit — but sorted sees ALL first
    .collect(toList());
// sorted() destroys parallelism benefit: must gather everything to sort
// → Entire dataset filtered into a single buffer, sorted, then limited
```

### 4.3 Internal Implementation — Pipeline Stages & Sink Chaining

The Streams API (java.util.stream) uses an internal pipeline model that fuses operations together:

```
Internal Class Hierarchy:

  BaseStream
    └── AbstractPipeline (linked list of stages)
        ├── ReferencePipeline.Head    (source stage)
        ├── ReferencePipeline.StatelessOp   (filter, map)
        ├── ReferencePipeline.StatefulOp    (sorted, distinct)
        └── IntPipeline, LongPipeline, DoublePipeline (primitive variants)
        
  Pipeline is a linked list: Head → Op → Op → Op → TerminalOp
```

```
Pipeline Linking (at construction):

  list.stream()   .filter(p)    .map(f)      .collect(toList())
       │              │            │               │
       ▼              ▼            ▼               ▼
  ┌─────────┐   ┌──────────┐ ┌──────────┐   ┌──────────────┐
  │  Head   │──►│FilterOp  │►│  MapOp   │   │ TerminalOp   │
  │(source) │   │(stores   │ │(stores   │   │ (triggers     │
  │         │   │ predicate│ │ function)│   │  execution)   │
  └─────────┘   └──────────┘ └──────────┘   └──────────────┘
  
  ← previousStage links     → nextStage links
```

#### Sink Chaining — The Execution Model

When the terminal operation triggers execution, the pipeline builds a chain of **Sinks**:

```
Sink Chain (at execution):

  Terminal calls: pipeline.wrapSink(terminalSink)
  Each stage wraps the downstream sink:

  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
  │ FilterSink   │────►│  MapSink     │────►│ CollectorSink│
  │              │     │              │     │  (terminal)  │
  │ accept(e):   │     │ accept(e):   │     │ accept(e):   │
  │  if pred(e)  │     │  e2=f(e)     │     │  list.add(e) │
  │    downstream│     │  downstream  │     │              │
  │    .accept(e)│     │  .accept(e2) │     │              │
  └──────────────┘     └──────────────┘     └──────────────┘
  
  Source calls filterSink.accept(element) for each element
  → filter tests predicate → if true, calls mapSink.accept(element)
  → map transforms → calls collectorSink.accept(transformed)
  
  ★ Key insight: NO intermediate collections created!
    Each element flows through the entire chain before the next element.
    This is "operation fusion" — same as loop fusion in compilers.
```

```java
// What the JVM effectively executes (after fusion):
// INSTEAD of:  filter all → collect → map all → collect → terminal
// IT DOES:     for each element: filter → map → collect (one pass)

for (T element : source) {
    if (predicate.test(element)) {          // FilterSink
        R mapped = function.apply(element); // MapSink
        resultList.add(mapped);             // CollectorSink
    }
}
// Single pass, no intermediate allocations — this is the fusion benefit
```

### 4.4 Spliterator — The Source Decomposition Engine

`Spliterator` (Splitable Iterator) is the internal abstraction for sourcing elements and splitting for parallelism.

```java
public interface Spliterator<T> {
    boolean tryAdvance(Consumer<? super T> action); // Process one element
    Spliterator<T> trySplit();  // Split into two halves for parallel processing
    long estimateSize();        // Approximate remaining elements
    int characteristics();      // Bitfield of properties (ORDERED, SORTED, SIZED, etc.)
}
```

```
Spliterator Characteristics:

  ORDERED    (0x00000010) — Elements have a defined encounter order
  DISTINCT   (0x00000001) — No duplicates
  SORTED     (0x00000004) — Elements are sorted by natural/comparator order
  SIZED      (0x00000040) — estimateSize() returns exact count
  NONNULL    (0x00000100) — Elements are guaranteed non-null
  IMMUTABLE  (0x00000400) — Source cannot be modified
  CONCURRENT (0x00001000) — Source safely modified concurrently
  SUBSIZED   (0x00004000) — trySplit() produces equally-sized halves
  
  These characteristics enable optimizations:
  - SIZED → count() returns immediately (no traversal)
  - SIZED + SUBSIZED → parallel splits are balanced
  - SORTED → sorted() is a no-op
  - DISTINCT → distinct() is a no-op
```

```java
// Custom Spliterator for a binary tree (enables efficient parallel tree processing):
class TreeSpliterator<T> implements Spliterator<T> {
    private Deque<TreeNode<T>> stack = new ArrayDeque<>();
    
    TreeSpliterator(TreeNode<T> root) {
        if (root != null) stack.push(root);
    }
    
    @Override
    public boolean tryAdvance(Consumer<? super T> action) {
        if (stack.isEmpty()) return false;
        TreeNode<T> node = stack.pop();
        action.accept(node.value);
        if (node.right != null) stack.push(node.right);
        if (node.left != null) stack.push(node.left);
        return true;
    }
    
    @Override
    public Spliterator<T> trySplit() {
        // Split: give the right subtree to another worker
        if (stack.isEmpty()) return null;
        TreeNode<T> node = stack.peek();
        if (node.right == null) return null;
        TreeNode<T> rightSubtree = node.right;
        node.right = null; // prevent double-processing
        return new TreeSpliterator<>(rightSubtree);
    }
    
    @Override public long estimateSize() { return Long.MAX_VALUE; }
    @Override public int characteristics() { return NONNULL; }
}
```

### 4.5 Parallel Streams — ForkJoinPool & Splitting Strategy

```
Parallel Stream Execution:

  Source Spliterator
         │
         ├── trySplit() ──────────────┐
         │                             │
    Left Split                   Right Split
         │                             │
    ├── trySplit()──┐          ├── trySplit()──┐
    │               │          │               │
  Chunk1         Chunk2      Chunk3         Chunk4
    │               │          │               │
  Thread 1       Thread 2   Thread 3       Thread 4
  (ForkJoin      (ForkJoin   (ForkJoin     (ForkJoin
   Worker)        Worker)     Worker)       Worker)
    │               │          │               │
    ▼               ▼          ▼               ▼
  Result1        Result2     Result3        Result4
    │               │          │               │
    └──── Combine ──┘          └── Combine ────┘
           │                          │
           └──────── Combine ─────────┘
                       │
                  Final Result
                  
  Default pool: ForkJoinPool.commonPool()
  Threads: Runtime.getRuntime().availableProcessors() - 1  (+ calling thread)
```

**When Parallel Streams HURT Performance:**

| Scenario | Why It Hurts | What to Do Instead |
|----------|-------------|-------------------|
| Small data set (< 10K elements) | Fork/join overhead > computation savings | Use sequential stream or loop |
| Blocking I/O in pipeline | Blocks ForkJoinPool threads; starves other parallel streams in the JVM | Use custom executor + CompletableFuture |
| `LinkedList` source | Poor splitting (no random access; `trySplit()` uses slow sequential scans) | Convert to ArrayList first, or use arrays |
| Stateful operations (`sorted`, `distinct`) | Requires synchronization barriers | Pre-sort data or use sequential for that phase |
| Shared mutable state | Data race — wrong results, not just slow | Eliminate mutation; use reduction/collection |
| Ordered operations (`forEachOrdered`, `limit`) | Forces sequential processing within parallel framework | Drop ordering if possible (`unordered()`) |
| Global `ForkJoinPool.commonPool()` | One slow parallel stream blocks all others in the JVM | Custom ForkJoinPool (see below) |

**Custom ForkJoinPool for isolation:**
```java
// ⚠ Undocumented but widely used pattern — submit parallel stream to custom pool:
ForkJoinPool customPool = new ForkJoinPool(4);  // 4 threads, isolated from commonPool
try {
    List<Result> results = customPool.submit(() ->
        data.parallelStream()
            .map(this::expensiveComputation)
            .collect(toList())
    ).get();  // .get() blocks until parallel stream completes
} finally {
    customPool.shutdown();
}
// Caveat: This works because ForkJoinTask.fork() uses the current pool's thread
// Not officially guaranteed by the spec — but used extensively in production
```

### 4.6 Memory & GC Impact

```
Stream Pipeline Memory Model:

  Sequential Stream:
  ┌─────────────────────────────────────────────────────────┐
  │ Pipeline Object Graph:                                   │
  │   Head → FilterOp → MapOp → (terminal triggers GC-able)│
  │   ~3 objects + lambda instances per intermediate op     │
  │   Fused Sink chain: ~1 object per op (wraps downstream) │
  │   NO intermediate collections (fusion)                  │
  │                                                          │
  │   Memory: O(pipeline_depth) — NOT O(data_size)          │
  └─────────────────────────────────────────────────────────┘
  
  Parallel Stream:
  ┌─────────────────────────────────────────────────────────┐
  │ Additional allocation:                                   │
  │   ForkJoinTask per split (recursive decomposition tree) │
  │   Intermediate results per split (for combining)        │
  │   Spliterator instances per split                       │
  │                                                          │
  │   Memory: O(pipeline_depth × splits × intermediate)     │
  │   GC pressure: higher due to task + partial result alloc│
  └─────────────────────────────────────────────────────────┘
  
  ⚠ flatMap: Creates a NEW stream per element — significant GC pressure
    for large datasets. Consider alternatives if allocation profiling shows
    flatMap as a hotspot.
```

```java
// flatMap GC pressure example:
orders.stream()
    .flatMap(order -> order.getLineItems().stream())  // new Stream per order!
    // For 1M orders × small lists = 1M Stream + Spliterator objects
    .filter(item -> item.getPrice() > 100)
    .collect(toList());

// Alternative with less GC pressure (Java 16+ mapMulti):
orders.stream()
    .<LineItem>mapMulti((order, consumer) -> {
        for (LineItem item : order.getLineItems()) {
            consumer.accept(item);  // No intermediate Stream objects
        }
    })
    .filter(item -> item.getPrice() > 100)
    .collect(toList());
```

---

## 5. Optional – Beyond Null Handling

### 5.1 Design Intent

`Optional<T>` was designed with a single, narrow purpose: **as a return type for methods that may not have a result**. It was NOT designed as a general-purpose null-replacement or a field type.

> *"Optional was intended to provide a limited mechanism for library method return types where there is a clear need to represent 'no result,' and where using null for that was overwhelmingly likely to cause errors."*  
> — Brian Goetz, Java Language Architect

### 5.2 Optional as a Monad — Formal Laws & Theory

In category theory, a **monad** is a type constructor `M` with two operations:
- **unit** (also called `return` or `of`): wraps a value → `T → M<T>`
- **bind** (also called `flatMap`): chains computations → `M<T> → (T → M<R>) → M<R>`

Java's `Optional` satisfies the three monad laws:

```java
// Let f: T → Optional<R>  and  g: R → Optional<S>

// LAW 1: Left Identity — wrapping a value and flatMapping is same as applying directly
Optional.of(x).flatMap(f)  ≡  f.apply(x)
// Wrapping then unwrapping is a no-op

// LAW 2: Right Identity — flatMapping with the wrapper is identity
optional.flatMap(Optional::of)  ≡  optional  
// Wrapping the contents produces the same Optional

// LAW 3: Associativity — nesting order doesn't matter
optional.flatMap(f).flatMap(g)  ≡  optional.flatMap(x -> f.apply(x).flatMap(g))
// (m >>= f) >>= g  ≡  m >>= (λx → f x >>= g)
```

**Why this matters in Java:**
```java
// Because Optional obeys monad laws, chains are PREDICTABLE:
Optional<String> city = findUser(userId)          // Optional<User>
    .flatMap(User::getAddress)                     // Optional<Address>
    .flatMap(Address::getCity)                     // Optional<String>
    .filter(c -> !c.isBlank());                    // Optional<String>

// We KNOW this is equivalent to:
Optional<String> city = findUser(userId)
    .flatMap(u -> u.getAddress().flatMap(Address::getCity))
    .filter(c -> !c.isBlank());

// We can refactor safely because the monad laws GUARANTEE equivalence.
// If Optional violated these laws, refactoring could silently change behavior.
```

```
Optional Operations Mapped to Monad/Functor Terminology:

  FP Concept    │ Optional Method      │ Type Signature
  ──────────────┼──────────────────────┼───────────────────────────
  unit / return │ Optional.of(v)       │ T → Optional<T>
  map (functor) │ optional.map(f)      │ Optional<T> → (T→R) → Optional<R>
  bind/flatMap  │ optional.flatMap(f)  │ Optional<T> → (T→Optional<R>) → Optional<R>
  filter        │ optional.filter(p)   │ Optional<T> → (T→boolean) → Optional<T>
  fold          │ map + orElse         │ Optional<T> → (T→R) → R → R
  orElse        │ optional.or(alt)     │ Optional<T> → (()→Optional<T>) → Optional<T>
  empty         │ Optional.empty()     │ () → Optional<T>  (the "zero" / mzero)
  
  MonadPlus laws (Optional.empty as zero):
    Optional.empty().flatMap(f)  ≡  Optional.empty()   // Left zero
    optional.flatMap(x -> Optional.empty())  ≡  Optional.empty()  // Right zero
```

**Comparison with other languages' Option types:**

| Feature | Java `Optional` | Scala `Option` | Kotlin `?` types | Haskell `Maybe` |
|---------|-----------------|----------------|-------------------|-----------------|
| Monad? | Yes (informally) | Yes (with for-comprehension) | No (nullable types, not wrapped) | Yes (formal) |
| Pattern match | No (use `map`/`orElse`) | `match { case Some(v) => ... }` | Smart casts (`if (x != null)`) | `case Just x -> ... Nothing -> ...` |
| Can be null itself | Yes (design flaw) | No | N/A (not wrapped) | No (enforced by type system) |
| Collection-like | `.stream()` (Java 9+) | Extends `Iterable` | N/A | `Foldable`, `Traversable` |
| Field usage | Discouraged | Common | Native nullable | Common |

### 5.3 Correct Usage Patterns

```java
// ✅ Method return type — expressive API contract
public Optional<User> findByEmail(String email) {
    return Optional.ofNullable(userRepository.findOne(email));
}

// ✅ Functional chaining — map / flatMap / filter
Optional<String> city = findByEmail("a@b.com")
    .map(User::getAddress)
    .map(Address::getCity)
    .filter(c -> !c.isBlank());

// ✅ map vs flatMap — critical distinction
Optional<Address> address = user.map(User::getAddress);       // User::getAddress returns Address
Optional<Address> address = user.flatMap(User::getAddress);   // User::getAddress returns Optional<Address>

// map:     (T → R)             → Optional<R>               (wraps result)
// flatMap: (T → Optional<R>)   → Optional<R>               (doesn't double-wrap)

// ✅ Providing defaults
String name = findUser(id)
    .map(User::getName)
    .orElse("Unknown");                       // Eager default (always evaluated)

String name = findUser(id)
    .map(User::getName)
    .orElseGet(() -> generateDefault(id));    // Lazy default (computed only if empty)

// ✅ Throwing on absence
User user = findUser(id)
    .orElseThrow(() -> new UserNotFoundException(id));  // Java 10+: orElseThrow() (no args)

// ✅ Java 9+: ifPresentOrElse
findUser(id).ifPresentOrElse(
    user -> processUser(user),
    () -> log.warn("User {} not found", id)
);

// ✅ Java 9+: or() — lazy alternative Optional
Optional<User> user = findInCache(id)
    .or(() -> findInDatabase(id))     // Only queries DB if cache is empty
    .or(() -> findInArchive(id));     // Only queries archive if DB is empty

// ✅ Java 9+: stream() — integrate with stream pipelines
List<String> names = userIds.stream()
    .map(this::findUser)                   // Stream<Optional<User>>
    .flatMap(Optional::stream)             // Filter out empty, unwrap present
    .map(User::getName)
    .collect(toList());
```

### 5.4 Anti-Patterns — What NOT To Do

```java
// ❌ NEVER: Optional as method parameter
public void processUser(Optional<User> user) { }
// Why: Callers must wrap with Optional.of()/ofNullable() — adds verbosity with no benefit
// Fix: Use @Nullable annotation or overloaded methods

// ❌ NEVER: Optional as field type
public class Order {
    private Optional<Discount> discount;  // ❌
    // Why: Adds 16 bytes overhead per Optional wrapper; not serializable by default
    //      Every access requires unwrap; field can itself be null (Optional that is null)
    // Fix: Use nullable field + getter returns Optional
    private Discount discount;  // ✅
    public Optional<Discount> getDiscount() { return Optional.ofNullable(discount); }
}

// ❌ NEVER: Optional in collections
List<Optional<String>> items;  // ❌ 
// Fix: Filter out nulls at the collection boundary
List<String> items;  // ✅ — contract: no nulls

// ❌ NEVER: Optional.get() without isPresent() check
String name = findUser(id).get();  // NoSuchElementException if empty
// Fix: Use orElse, orElseThrow, or functional methods

// ❌ Wrapping and immediately unwrapping
Optional.ofNullable(value).orElse(default);  // Pointless overhead
// Fix: Just use ternary
value != null ? value : default;

// ❌ Using Optional for conditional logic (isPresent + get)
Optional<User> opt = findUser(id);
if (opt.isPresent()) {          // Treating Optional as an if-null check
    process(opt.get());         // This is NOT the purpose of Optional
}
// Fix: opt.ifPresent(this::process);

// ❌ Returning Optional.of(null) — NullPointerException
Optional.of(null);  // NPE!
// Fix: Optional.ofNullable(potentiallyNullValue);
```

### 5.5 Performance Concerns

```
Optional Allocation Cost:

  Optional.of(value):
    → Allocates a new Optional object wrapping the value
    → 16 bytes (object header) + 8 bytes (reference field) = ~24 bytes
    → Minor GC eligible (short-lived) — but millions per second adds GC pressure

  Optional.empty():
    → Returns a cached singleton — NO allocation
    → Comparing with == is safe: Optional.empty() == Optional.empty()

  Escape Analysis Opportunity:
    → If the Optional doesn't escape the method, JIT may eliminate the allocation
    → Example: createOptional().map(f).orElse(default) — may be scalar-replaced
    → But: if the Optional crosses method boundaries, escape analysis often fails

  Primitive Alternatives (avoid autoboxing):
    OptionalInt, OptionalLong, OptionalDouble
    → Do NOT extend Optional<T>
    → Cannot be used with map/flatMap chains as smoothly
    → getAsInt(), orElse(int), ifPresent(IntConsumer)
```

### 5.6 Serialization Implications

- `Optional` does **NOT** implement `Serializable`
- Jackson can handle it with `jackson-datatype-jdk8` module (or included in Spring Boot by default)
- JPA/Hibernate: Do **not** use Optional as an entity field type — use nullable fields with Optional getters
- Protocol Buffers / Avro: `Optional` doesn't map directly; use `has_field()` patterns

```java
// Jackson serialization (with jackson-datatype-jdk8):
@JsonProperty
private String name;               // Always present

@JsonInclude(Include.NON_ABSENT)   // Omit from JSON if empty
public Optional<String> getNickname() { return Optional.ofNullable(nickname); }
// Serializes: {"name":"John"} or {"name":"John","nickname":"Johnny"}

// JPA entity — correct pattern:
@Entity
public class Order {
    @Column(nullable = true)
    private LocalDate shippedDate;  // nullable in DB
    
    public Optional<LocalDate> getShippedDate() {
        return Optional.ofNullable(shippedDate);
    }
}
```

---

## 6. Functional Programming Patterns in Java

### 6.1 Immutability

```java
// Pre-records: manual immutable class (verbose, error-prone)
public final class Money {
    private final BigDecimal amount;
    private final Currency currency;
    
    public Money(BigDecimal amount, Currency currency) {
        this.amount = Objects.requireNonNull(amount);
        this.currency = Objects.requireNonNull(currency);
    }
    // getters, equals, hashCode, toString — boilerplate
    
    // Functional transformation: return NEW instance (never mutate)
    public Money add(Money other) {
        if (!this.currency.equals(other.currency)) throw new IllegalArgumentException();
        return new Money(this.amount.add(other.amount), this.currency);
    }
}

// Java 16+: Records — immutable by default
public record Money(BigDecimal amount, Currency currency) {
    // Compact constructor for validation
    public Money {
        Objects.requireNonNull(amount);
        Objects.requireNonNull(currency);
    }
    
    public Money add(Money other) {
        if (!this.currency.equals(other.currency)) throw new IllegalArgumentException();
        return new Money(this.amount.add(other.amount), this.currency);
    }
}

// Immutable collections (Java 9+):
List<String> names = List.of("Alice", "Bob");           // Truly immutable
Map<String, Integer> scores = Map.of("A", 100, "B", 90); // Throws on mutation
// vs Collections.unmodifiableList() — just a VIEW, source can still change
```

### 6.2 Pure Functions & Side-Effect Minimization

```java
// PURE function: same input → same output, no side effects
Function<List<Integer>, OptionalInt> max = list ->
    list.stream().mapToInt(Integer::intValue).max();

// IMPURE: depends on external state
Function<String, String> greet = name ->
    "Hello " + name + " at " + LocalTime.now();  // Time = external state

// Strategy: Push side effects to the boundary
public class OrderService {
    // PURE: Business logic transformation (testable, parallelizable)
    public OrderResult calculateOrder(Order order, PricingRules rules) {
        BigDecimal subtotal = order.getItems().stream()
            .map(item -> rules.getPrice(item).multiply(BigDecimal.valueOf(item.getQty())))
            .reduce(BigDecimal.ZERO, BigDecimal::add);
        
        BigDecimal tax = rules.calculateTax(subtotal, order.getRegion());
        BigDecimal discount = rules.applyDiscounts(subtotal, order.getCoupons());
        
        return new OrderResult(subtotal, tax, discount); // No I/O, no mutation
    }
    
    // IMPURE boundary: Side effects (I/O) isolated here
    public void processOrder(String orderId) {
        Order order = orderRepo.findById(orderId);       // I/O (side effect)
        PricingRules rules = pricingService.getRules();   // I/O (side effect)
        OrderResult result = calculateOrder(order, rules); // PURE
        orderRepo.save(result);                           // I/O (side effect)
        notificationService.notify(result);               // I/O (side effect)
    }
}
```

### 6.3 Function Composition & Higher-Order Functions

```java
// Higher-order function: takes/returns functions
public static <T> Predicate<T> not(Predicate<T> p) {
    return p.negate();  // Java 11+: Predicate.not()
}

// Function composition for data transformation pipeline
Function<String, String> sanitize = s -> s.trim().toLowerCase();
Function<String, String> removeSpecial = s -> s.replaceAll("[^a-z0-9 ]", "");
Function<String, List<String>> tokenize = s -> List.of(s.split("\\s+"));
Function<List<String>, List<String>> stopWordFilter = 
    words -> words.stream().filter(w -> !STOP_WORDS.contains(w)).collect(toList());

// Compose into a reusable pipeline
Function<String, List<String>> textProcessor = 
    sanitize.andThen(removeSpecial).andThen(tokenize).andThen(stopWordFilter);

List<String> tokens = textProcessor.apply("  Hello, World! This is a TEST.  ");
// ["hello", "world", "test"]

// Strategy pattern via Function (replacing class hierarchy):
Map<PaymentType, Function<Payment, PaymentResult>> processors = Map.of(
    PaymentType.CREDIT_CARD, this::processCreditCard,
    PaymentType.BANK_TRANSFER, this::processBankTransfer,
    PaymentType.CRYPTO, this::processCrypto
);

PaymentResult result = processors.get(payment.getType()).apply(payment);
```

### 6.4 Monadic Thinking — Optional & Stream as Monads

In FP, a monad is a type that supports `flatMap` (bind) and a unit operation (constructor). Java's `Optional` and `Stream` exhibit monadic behavior:

```java
// Optional monad — chaining computations that may fail
Optional<String> userCity = findUser(userId)          // Optional<User>
    .flatMap(User::getAddress)                         // Optional<Address>
    .flatMap(Address::getCity)                         // Optional<String>
    .filter(city -> !city.isBlank());                  // Optional<String>
// Each flatMap: if empty → short-circuits; if present → applies function
// Equivalent to nested null checks, but declarative

// Stream monad — chaining transformations over collections
List<String> allSkills = employees.stream()            // Stream<Employee>
    .flatMap(e -> e.getSkills().stream())               // Stream<String>
    .distinct()
    .sorted()
    .collect(toList());

// CompletableFuture monad — chaining async computations
CompletableFuture<Order> enrichedOrder = fetchOrder(orderId)  // CF<Order>
    .thenCompose(order -> fetchPricing(order))                 // CF<PricedOrder>
    .thenCompose(priced -> applyDiscounts(priced))             // CF<DiscountedOrder>
    .thenApply(discounted -> finalize(discounted));            // CF<Order>
// thenCompose = flatMap; thenApply = map
```

### 6.5 Either / Try / Result Pattern — Error Handling the FP Way

Java's `Optional` only models **presence/absence**. It cannot carry error information. In functional languages, this is solved with **sum types** like `Either<L, R>` or `Try<T>`. While Java doesn't include these in the standard library, the pattern is essential for senior-level FP understanding.

#### The Problem with Exceptions in FP Pipelines

```java
// Exceptions break functional composition:
List<Config> configs = paths.stream()
    .map(path -> Files.readString(path))   // ❌ Throws IOException — UNHANDLED
    .map(ConfigParser::parse)
    .toList();
// Won't compile! Checked exceptions can't escape lambda boundaries
// in standard functional interfaces (Function doesn't declare throws)

// Ugly workaround:
.map(path -> {
    try { return Files.readString(path); }
    catch (IOException e) { throw new UncheckedIOException(e); }    // Wrap and rethrow
})
// This defeats the purpose of functional composition — we're back to try/catch
```

#### Either<L, R> — The Functional Alternative

`Either<L, R>` represents a value of one of two types: **Left** (conventionally the error) or **Right** (conventionally the success).

```java
// Custom Either implementation (or use vavr/cyclops library):
public sealed interface Either<L, R> permits Either.Left, Either.Right {
    record Left<L, R>(L value) implements Either<L, R> {}
    record Right<L, R>(R value) implements Either<L, R> {}
    
    static <L, R> Either<L, R> right(R value) { return new Right<>(value); }
    static <L, R> Either<L, R> left(L error) { return new Left<>(error); }
    
    default <T> Either<L, T> map(Function<R, T> f) {
        return switch (this) {
            case Right<L, R> r -> Either.right(f.apply(r.value()));
            case Left<L, R> l -> (Either<L, T>) l;  // Error propagates unchanged
        };
    }
    
    default <T> Either<L, T> flatMap(Function<R, Either<L, T>> f) {
        return switch (this) {
            case Right<L, R> r -> f.apply(r.value());
            case Left<L, R> l -> (Either<L, T>) l;
        };
    }
    
    default R getOrElse(R defaultValue) {
        return switch (this) { case Right<L,R> r -> r.value(); case Left<L,R> l -> defaultValue; };
    }
}
```

#### Railway-Oriented Programming

This pattern (coined by Scott Wlaschin) models a processing pipeline as a **railway track**: the happy path continues on the "right track," while any error diverts to the "left track" — and subsequent operations are skipped automatically.

```
Railway-Oriented Programming Visualization:

  Input ──► validate ──► enrich ──► calculate ──► persist ──► Result
            │             │          │             │
            ▼             ▼          ▼             ▼
          Error         Error      Error         Error
          (Left)        (Left)     (Left)        (Left)
          
  Once on the error track, all subsequent map/flatMap operations are SKIPPED.
  The error propagates to the end automatically — no if/else chains needed.

  Compare to imperative:
    if (validationResult.hasError()) return error;
    if (enrichResult.hasError()) return error;
    if (calcResult.hasError()) return error;
    // ... verbose and error-prone error checking
```

```java
// Railway-oriented order processing:
Either<OrderError, OrderResult> processOrder(OrderRequest request) {
    return validateRequest(request)                         // Either<Error, ValidRequest>
        .flatMap(this::checkInventory)                      // Either<Error, InventoryCheck>
        .flatMap(this::calculatePricing)                    // Either<Error, PricedOrder>
        .flatMap(this::applyDiscounts)                      // Either<Error, DiscountedOrder>
        .map(this::buildResult);                            // Either<Error, OrderResult>
}
// If validateRequest fails → returns Left(error) → all subsequent flatMaps skip → Left propagates

// In a Stream — collecting errors instead of throwing:
List<Either<String, Config>> results = paths.stream()
    .map(path -> {
        try { return Either.<String, Config>right(loadConfig(path)); }
        catch (Exception e) { return Either.<String, Config>left("Failed: " + path + " - " + e.getMessage()); }
    })
    .toList();

List<Config> successes = results.stream().filter(e -> e instanceof Either.Right).map(e -> e.getOrElse(null)).toList();
List<String> errors = results.stream().filter(e -> e instanceof Either.Left).map(e -> ((Either.Left<String, Config>) e).value()).toList();
```

#### Try<T> — Exception-Safe Computation

```java
// Try monad — wraps computation that may throw
public sealed interface Try<T> permits Try.Success, Try.Failure {
    record Success<T>(T value) implements Try<T> {}
    record Failure<T>(Throwable error) implements Try<T> {}
    
    static <T> Try<T> of(Callable<T> computation) {
        try { return new Success<>(computation.call()); }
        catch (Exception e) { return new Failure<>(e); }
    }
    
    default <R> Try<R> map(Function<T, R> f) {
        return switch (this) {
            case Success<T> s -> Try.of(() -> f.apply(s.value()));
            case Failure<T> fail -> (Try<R>) fail;
        };
    }
    
    default <R> Try<R> flatMap(Function<T, Try<R>> f) {
        return switch (this) {
            case Success<T> s -> f.apply(s.value());
            case Failure<T> fail -> (Try<R>) fail;
        };
    }
    
    default T getOrElse(T fallback) {
        return switch (this) { case Success<T> s -> s.value(); case Failure<T> f -> fallback; };
    }
    
    default Try<T> recover(Function<Throwable, T> recovery) {
        return switch (this) {
            case Failure<T> f -> Try.of(() -> recovery.apply(f.error()));
            case Success<T> s -> s;
        };
    }
}

// Usage — exceptions never leak into stream pipeline:
List<Config> configs = paths.stream()
    .map(path -> Try.of(() -> loadConfig(path)))           // Try<Config>
    .map(t -> t.recover(ex -> Config.defaultConfig()))     // Try<Config> (recovered)
    .filter(t -> t instanceof Try.Success)
    .map(t -> t.getOrElse(null))
    .toList();
```

**Available libraries for production use:**
- **Vavr** (`io.vavr`): Full FP library — `Either`, `Try`, `Tuple`, `Lazy`, persistent collections
- **Cyclops** (`com.oath.cyclops`): Reactive FP — `Either`, `Try`, integration with Reactor/RxJava
- **Result** (custom): Lightweight, often project-specific

### 6.6 Currying & Partial Application (Theory & Limitations)

```
Currying Theory:

  In lambda calculus, EVERY function takes exactly ONE argument.
  Multi-argument functions are achieved through currying:
  
  f(a, b, c) = a + b + c
  
  Curried form:
  f = λa. (λb. (λc. a + b + c))
  
  Evaluation:
  f(1)     → λb. (λc. 1 + b + c)    // Partially applied — waiting for b and c
  f(1)(2)  → λc. 1 + 2 + c           // Partially applied — waiting for c
  f(1)(2)(3) → 1 + 2 + 3 = 6         // Fully applied
  
  Partial Application ≠ Currying:
  • Currying: transform f(a,b,c) into f(a)(b)(c) — always one arg at a time
  • Partial Application: fix some arguments, return function for remaining
    e.g., g = f(1, _, 3) → g(b) = f(1, b, 3)
  
  Java supports PARTIAL APPLICATION naturally (via lambdas capturing arguments)
  but makes CURRYING verbose due to nested Function<> types.
```

```java
// Currying: transforming f(a, b) into f(a)(b)
// Java doesn't support this natively, but can be emulated:

// Standard BiFunction:
BiFunction<Double, Double, Double> add = (a, b) -> a + b;
double result = add.apply(3.0, 4.0);  // 7.0

// Curried version:
Function<Double, Function<Double, Double>> curriedAdd = a -> b -> a + b;
Function<Double, Double> addThree = curriedAdd.apply(3.0);  // Partial application
double result = addThree.apply(4.0);  // 7.0

// Practical use: creating specialized functions from general ones
Function<String, Function<String, String>> formatter = 
    prefix -> value -> prefix + ": " + value;

Function<String, String> errorFormatter = formatter.apply("ERROR");
Function<String, String> infoFormatter = formatter.apply("INFO");

String msg = errorFormatter.apply("Connection timeout");  // "ERROR: Connection timeout"

// Generic curry/uncurry utilities:
static <A, B, R> Function<A, Function<B, R>> curry(BiFunction<A, B, R> f) {
    return a -> b -> f.apply(a, b);
}
static <A, B, R> BiFunction<A, B, R> uncurry(Function<A, Function<B, R>> f) {
    return (a, b) -> f.apply(a).apply(b);
}

// Limitation: Java's type system makes currying verbose beyond 2-3 levels
// Scala: def add(a: Int)(b: Int) = a + b  — native currying syntax
// Haskell: ALL functions are curried by default — add a b = a + b
```

### 6.7 Lazy Evaluation Patterns

```java
// Lazy initialization with Supplier
class ExpensiveResource {
    private Supplier<Connection> connectionSupplier = () -> {
        Connection conn = createConnection();  // Expensive
        this.connectionSupplier = () -> conn;  // Replace with memoized version
        return conn;
    };
    
    public Connection getConnection() {
        return connectionSupplier.get();  // First call: creates; subsequent: returns cached
    }
}

// Thread-safe lazy with double-checked locking via Supplier
public static <T> Supplier<T> memoize(Supplier<T> delegate) {
    AtomicReference<T> ref = new AtomicReference<>();
    return () -> {
        T val = ref.get();
        if (val == null) {
            synchronized (ref) {
                val = ref.get();
                if (val == null) {
                    val = Objects.requireNonNull(delegate.get());
                    ref.set(val);
                }
            }
        }
        return val;
    };
}

// Lazy stream evaluation — infinite sequences
Stream<BigInteger> fibonacci = Stream.iterate(
    new BigInteger[]{BigInteger.ZERO, BigInteger.ONE},
    pair -> new BigInteger[]{pair[1], pair[0].add(pair[1])}
).map(pair -> pair[0]);

// Only computes as many as consumed:
fibonacci.limit(100).forEach(System.out::println);
```

### 6.8 Custom Collector Design

```java
// Collector<T, A, R>:
//   T = input element type
//   A = mutable accumulator type
//   R = result type

// Custom Collector: Group into batches of size N
public static <T> Collector<T, ?, List<List<T>>> toBatches(int batchSize) {
    return Collector.of(
        // Supplier: create accumulator
        ArrayList::new,
        
        // Accumulator: add element to accumulator
        (List<List<T>> batches, T element) -> {
            if (batches.isEmpty() || batches.get(batches.size() - 1).size() >= batchSize) {
                batches.add(new ArrayList<>());
            }
            batches.get(batches.size() - 1).add(element);
        },
        
        // Combiner: merge two accumulators (for parallel streams)
        (left, right) -> {
            // Merge last batch of left with first batch of right if possible
            if (!left.isEmpty() && !right.isEmpty()) {
                List<T> lastLeft = left.get(left.size() - 1);
                if (lastLeft.size() < batchSize) {
                    List<T> firstRight = right.get(0);
                    int space = batchSize - lastLeft.size();
                    lastLeft.addAll(firstRight.subList(0, Math.min(space, firstRight.size())));
                    if (space < firstRight.size()) {
                        right.set(0, new ArrayList<>(firstRight.subList(space, firstRight.size())));
                    } else {
                        right.remove(0);
                    }
                }
            }
            left.addAll(right);
            return left;
        },
        
        // Characteristics
        Collector.Characteristics.IDENTITY_FINISH  // A == R, no finisher needed
    );
}

// Usage:
List<List<Order>> batches = orders.stream().collect(toBatches(500));
batches.forEach(batch -> orderRepository.saveAll(batch)); // Batch DB inserts
```

**Collector Characteristics:**

| Characteristic | Meaning | Impact |
|---------------|---------|--------|
| `CONCURRENT` | Accumulator supports concurrent access; combiner not needed | Enables parallel accumulation into single container |
| `UNORDERED` | Collection doesn't depend on encounter order | Parallel: no ordering constraints for combining |
| `IDENTITY_FINISH` | Finisher is identity; accumulator IS the result | Skip finisher step; cast A directly to R |

```java
// Reduction patterns:
// reduce(identity, accumulator, combiner) — for parallel correctness

// Sum (with identity):
int total = orders.stream()
    .mapToInt(Order::getQuantity)
    .reduce(0, Integer::sum);  // identity=0, accumulator=sum

// Custom reduction — most expensive order:
Optional<Order> mostExpensive = orders.stream()
    .reduce((a, b) -> a.getTotal().compareTo(b.getTotal()) >= 0 ? a : b);

// Mutable reduction using collect (preferred for collections/strings):
String csv = items.stream()
    .map(Item::getName)
    .collect(Collectors.joining(", "));
```

---

## 7. Concurrency & Functional Programming

### 7.0 Theoretical Foundation — Why FP Enables Safe Concurrency

#### The Shared Mutable State Problem

Concurrency bugs arise from the intersection of three properties. Remove any ONE and concurrency is safe:

```
The Concurrency Hazard Triangle:

         Shared State
            /\
           /  \
          /    \
         / BUG  \
        /  ZONE  \
       /          \
      /____________\
  Mutation      Concurrency
  
  Remove SHARING   → thread-local, message passing (Actor model)
  Remove MUTATION   → immutable data, pure functions (FP approach)
  Remove CONCURRENCY → single-threaded execution
  
  FP attacks the MUTATION vertex:
  • Pure functions produce new values instead of modifying existing ones
  • Immutable data can be safely shared across any number of threads
  • No locks needed when nothing changes
```

#### Formal Safety Guarantees of Pure Functions

```
Theorem: A pure function f: A → B is INHERENTLY THREAD-SAFE.

Proof:
  1. f has no side effects → no writes to shared state
  2. f depends only on its arguments → no reads from shared mutable state
  3. Multiple threads calling f(x) simultaneously:
     - Each thread has its own stack frame (arguments, locals)
     - No shared memory is read or written
     - Result is determined entirely by input
  4. Therefore: no data races, no need for synchronization
  
  Corollary: A stream pipeline composed entirely of pure functions
  can be safely parallelized with .parallel() — GUARANTEED correct.
  
  This is WHY the Streams API requires:
  • Non-interference (don't modify source)
  • Statelessness (don't depend on external mutable state)
  • No side effects in behavioral parameters
  — These are the conditions that make the above theorem applicable.
```

#### Java Memory Model (JMM) Guarantees for FP

The JMM defines **happens-before** relationships that determine when one thread's writes are visible to another thread. FP leverages these implicitly:

```
JMM Happens-Before Rules Relevant to Streams & Lambdas:

  1. PROGRAM ORDER RULE:
     Each action in a thread happens-before every subsequent action in that thread.
     → Sequential stream operations see each other's results (trivially)

  2. FORK-JOIN RULE (ForkJoinPool):
     fork() of a ForkJoinTask HB the task's execution
     Task completion HB join() that observes it
     → Parallel stream splits and merges have correct visibility
     → combiner() in Collector sees all accumulated results from both sides

  3. VOLATILE VARIABLE RULE:
     Write to volatile HB subsequent read of that volatile
     → Stream's internal cancellation flag (for short-circuiting) uses this

  4. FINAL FIELD SEMANTICS:
     Construction of an object with final fields HB any thread reading those fields
     → Records (all fields final) and immutable objects are safely publishable
     → Lambda captures of effectively final variables are safely visible

  5. THREAD START RULE:
     Thread.start() HB any action in the started thread
     → ForkJoinPool worker threads see all state established before parallel stream began
```

```java
// Practical implications:

// ✅ SAFE: Immutable data shared before parallel stream
List<Config> configs = List.of(config1, config2, config3);  // Immutable
long count = configs.parallelStream()                        // Safe — List.of is unmodifiable
    .filter(c -> c.isEnabled())                              // Safe — Config is immutable record
    .count();

// ✅ SAFE: Accumulator in Collector — framework provides HB guarantees
Map<String, List<Order>> grouped = orders.parallelStream()
    .collect(Collectors.groupingBy(Order::getRegion));
// Internally: each thread accumulates into its OWN map → combined via combiner
// combine() is called after split processing completes → HB guarantee by ForkJoin framework

// ❌ UNSAFE: Shared mutable state without HB relationship
HashMap<String, Integer> counts = new HashMap<>();           // Mutable, shared
orders.parallelStream().forEach(o -> 
    counts.merge(o.getRegion(), 1, Integer::sum));            // DATA RACE
// Multiple threads write to same HashMap — no HB between different ForkJoin tasks
// Fix: Use Collectors.groupingBy(... counting()) — lets framework handle concurrency
```

### 7.1 Statelessness & Thread Safety

The core promise of FP for concurrency: **stateless transformations don't need synchronization**.

```java
// Thread-safe by design — no shared mutable state:
Function<Order, OrderDTO> mapper = order -> new OrderDTO(
    order.getId(), 
    order.getTotal(), 
    order.getStatus().name()
);
// Can be safely shared across threads, used in parallel streams,
// called from multiple CompletableFuture chains — no locks needed

// UNSAFE — shared mutable state in a lambda:
List<String> results = new ArrayList<>();  // NOT thread-safe
items.parallelStream()
    .filter(item -> item.isActive())
    .forEach(item -> results.add(item.getName()));  // ❌ Race condition!
// ArrayList.add is not atomic → lost updates, ArrayIndexOutOfBoundsException

// FIX: Use collect() — thread-safe accumulation:
List<String> results = items.parallelStream()
    .filter(Item::isActive)
    .map(Item::getName)
    .collect(Collectors.toList());  // ✅ Internally uses combiner for parallel merge
```

### 7.2 Non-Interference Rule

The **non-interference** contract: the stream source must NOT be modified during stream pipeline execution.

```java
// ❌ VIOLATION: Modifying source during stream execution
List<String> names = new ArrayList<>(List.of("Alice", "Bob", "Charlie"));
names.stream()
    .filter(n -> {
        if (n.equals("Bob")) names.add("Dave");  // ConcurrentModificationException!
        return n.length() > 3;
    })
    .collect(toList());

// ❌ SUBTLE VIOLATION: Stateful predicate
Set<String> seen = new HashSet<>();
list.parallelStream()
    .filter(s -> seen.add(s))  // Uses mutable state in predicate — race condition
    .collect(toList());
// Fix: Use .distinct() instead — stream framework handles concurrency internally

// ✅ SAFE: Behavioral parameters must be:
//   - Non-interfering (don't modify the source)
//   - Stateless (or use thread-safe state if absolutely necessary)
//   - Without side effects (for most intermediate operations)
```

### 7.3 Reduction Operations — Associativity & Identity Requirements

```java
// reduce(identity, accumulator, combiner)
// Requirements for correctness in parallel:
//   1. accumulator must be ASSOCIATIVE:  (a op b) op c == a op (b op c)
//   2. identity must be a TRUE identity: identity op x == x
//   3. combiner must be compatible with accumulator

// ✅ Correct parallel reduction:
int sum = numbers.parallelStream().reduce(0, Integer::sum);
// 0 is identity for addition; addition is associative

// ❌ BROKEN parallel reduction — non-associative:
double average = numbers.parallelStream()
    .reduce(0.0, (a, b) -> (a + b) / 2.0);  // NOT associative!
// Sequential: ((1+2)/2 + 3)/2 = 2.25
// Parallel:   Thread1: (1+2)/2=1.5   Thread2: (3+4)/2=3.5  Combined: (1.5+3.5)/2=2.5
// Different results! Use: stream.mapToDouble(...).average()

// ❌ BROKEN identity:
int product = numbers.parallelStream().reduce(1, Integer::sum);
// Identity 1 is WRONG for sum — should be 0
// In parallel: each chunk starts with 1 → adds extra 1 per chunk

// Mutable reduction — Collector guarantees:
// supplier() → creates fresh accumulator PER CHUNK in parallel
// combiner() → merges two accumulators (must produce correct union)
// These guarantees make collect() inherently safe for parallel streams
```

### 7.4 Side Effects in Parallel Streams — Common Pitfalls

```java
// ❌ Pitfall 1: Logging order assumption
orders.parallelStream()
    .map(this::processOrder)
    .forEach(result -> log.info("Processed: {}", result));
// Log order is UNDEFINED in parallel — use forEachOrdered() if order matters
// But forEachOrdered() negates parallelism benefit for the terminal operation

// ❌ Pitfall 2: Writing to shared file/DB without synchronization
orders.parallelStream()
    .forEach(order -> {
        String line = order.toCsv();
        writer.write(line);  // BufferedWriter is NOT thread-safe → corrupt output
    });
// Fix: Collect first, then write sequentially
String csv = orders.parallelStream()
    .map(Order::toCsv)
    .collect(Collectors.joining("\n"));
writer.write(csv);

// ❌ Pitfall 3: ThreadLocal in parallel streams
private static final ThreadLocal<SimpleDateFormat> df = 
    ThreadLocal.withInitial(() -> new SimpleDateFormat("yyyy-MM-dd"));

dates.parallelStream()
    .map(d -> df.get().format(d))  // Different threads, different ThreadLocal values
    // Leaks ThreadLocal values into ForkJoinPool worker threads (long-lived)
    .collect(toList());
// Fix: Use DateTimeFormatter (thread-safe, immutable)
```

### 7.5 Memory Visibility Concerns

```java
// ForkJoinPool tasks use work-stealing — happens-before is guaranteed by the framework:
// 1. The fork of a task HB its execution
// 2. The completion of a task HB the join that observes it
// 3. Stream framework handles this internally — no manual synchronization needed

// BUT: if you bring your own shared state, YOU must ensure visibility:
class SharedCounter {
    private int count = 0;  // No volatile, no synchronization
    
    void process(Stream<Item> items) {
        items.parallelStream().forEach(item -> {
            count++;  // ❌ NOT atomic, NOT visible across threads
        });
        System.out.println(count);  // May print stale/wrong value
    }
}
// Fix: Use AtomicInteger, or better — reduce/collect instead of mutating shared state
```

### 7.6 ForkJoinPool Tuning for Parallel Streams

```java
// Default parallelism = Runtime.availableProcessors()
// Affects ALL parallel streams in the JVM (shared commonPool)

// Change default parallelism (system property, must be set at startup):
// -Djava.util.concurrent.ForkJoinPool.common.parallelism=16

// Isolation: use custom pool for stream execution
ForkJoinPool ioPool = new ForkJoinPool(32);  // More threads for I/O-bound work
try {
    List<Response> responses = ioPool.submit(() ->
        urls.parallelStream()
            .map(this::httpGet)  // Blocking I/O
            .collect(toList())
    ).get();
} finally {
    ioPool.shutdown();
}
// ⚠ CAUTION: Blocking I/O in ForkJoinPool is generally an anti-pattern
// Better: Use CompletableFuture + custom ThreadPoolExecutor for I/O
// Or: Virtual threads (Java 21+) for blocking I/O at scale
```

### 7.7 Real Production Debugging Scenario

```
Scenario: CPU spike + increased p99 latency in a Spring Boot service

Investigation:
  1. Thread dump (jstack) shows multiple threads in:
     ForkJoinPool.commonPool-worker-N → BLOCKED on synchronized block
     
  2. Root cause: A developer used parallelStream() and called a method
     that internally uses synchronized (legacy DAO with connection pooling)
     
  3. Effect: ForkJoinPool workers (shared across the JVM) are blocked
     waiting for DB connections → OTHER parallel streams in the same JVM
     are starved → cascade of timeouts
     
  4. Fix:
     a. Replaced parallelStream() with sequential stream for DB-accessing code
     b. For CPU-bound parallel work: continued using parallel streams
     c. For I/O-bound parallel work: migrated to CompletableFuture + bounded executor
     d. Added custom ForkJoinPool where parallel streams touch external resources
     
  5. Prevention: Code review rule —
     "No parallelStream() unless profiling proves CPU-bound benefit with >10K elements"
```

---

## 8. Performance & Optimization

### 8.1 Lambda Allocation Cost

```
Lambda Instantiation Cost:

  Non-Capturing Lambda (no free variables):
  ┌─────────────────────────────────────────────────┐
  │ s -> s.toUpperCase()                             │
  │                                                   │
  │ COST: ZERO per call (after first invocation)     │
  │ The JVM caches a singleton instance              │
  │ → No heap allocation, no GC pressure             │
  │ → Equivalent to a static method call             │
  └─────────────────────────────────────────────────┘

  Capturing Lambda (captures local variables):
  ┌─────────────────────────────────────────────────┐
  │ int threshold = 10;                              │
  │ filter(n -> n > threshold)                       │
  │                                                   │
  │ COST: One object allocation per invocation       │
  │ → New instance with field: final int arg$1 = 10  │
  │ → ~24 bytes on heap, short-lived → Minor GC      │
  │ → JIT escape analysis may eliminate it            │
  └─────────────────────────────────────────────────┘

  Instance-Capturing Lambda (captures `this`):
  ┌─────────────────────────────────────────────────┐
  │ items.forEach(item -> this.process(item))        │
  │                                                   │
  │ COST: One object allocation per invocation       │
  │ → Holds reference to enclosing instance          │
  │ → Prevents enclosing object from GC              │
  │ → Potential memory leak in long-lived callbacks  │
  └─────────────────────────────────────────────────┘
  
  Method Reference:
  ┌─────────────────────────────────────────────────┐
  │ String::toUpperCase (unbound) → singleton        │
  │ this::process (bound) → new instance per call    │
  │ System.out::println (bound) → new instance*      │
  │   *JIT may optimize if target is effectively     │
  │    constant (System.out rarely changes)           │
  └─────────────────────────────────────────────────┘
```

### 8.2 Autoboxing Overhead — The Silent Killer

```java
// ❌ Autoboxing in stream pipeline (Stream<Integer> instead of IntStream):
long sum = numbers.stream()                 // Stream<Integer>
    .filter(n -> n > 0)                      // unbox Integer → int for comparison
    .map(n -> n * 2)                         // unbox, multiply, REBOX to Integer
    .reduce(0, Integer::sum);               // unbox both, add, REBOX result
// Per element: up to 3 unbox + 2 box operations = 5 boxing operations
// Bottleneck: Integer.valueOf() creates objects (cached only for -128 to 127)

// ✅ Primitive stream — ZERO boxing:
long sum = numbers.stream()
    .mapToInt(Integer::intValue)            // Unbox once to IntStream
    .filter(n -> n > 0)                      // int comparison (no boxing)
    .map(n -> n * 2)                         // int multiply (no boxing)
    .sum();                                  // native sum (no boxing)
```

**Autoboxing performance impact:**

| Operation | Boxed (`Stream<Integer>`) | Primitive (`IntStream`) | Overhead |
|-----------|--------------------------|------------------------|----------|
| Element storage | 16 bytes (Integer object) | 4 bytes (int value) | 4x memory |
| Comparison | Unbox + compare | Direct compare | ~2x slower |
| Arithmetic | Unbox + compute + rebox | Direct compute | ~3-5x slower |
| GC pressure | N Integer objects per pipeline | Zero objects | Significant for >10K elements |
| Cache efficiency | Poor (object references scattered) | Excellent (contiguous array) | 5-10x throughput difference |

### 8.3 Stream vs Loop Performance — Comprehensive Comparison

```
                    Performance Comparison (typical measurements via JMH):
                    
  Operation          │ for-loop │ Stream  │ Parallel │ Notes
  ───────────────────┼──────────┼─────────┼──────────┼──────────────────────
  Sum 1M ints        │ ~0.5ms   │ ~2.5ms  │ ~1.2ms   │ Stream: overhead from
                     │          │         │          │ pipeline objects
  ───────────────────┼──────────┼─────────┼──────────┼──────────────────────
  Filter + map       │ ~1.0ms   │ ~3.0ms  │ ~1.5ms   │ Loop wins for
  (1M elements)      │          │         │          │ simple operations
  ───────────────────┼──────────┼─────────┼──────────┼──────────────────────
  Complex transform  │ ~15ms    │ ~16ms   │ ~5ms     │ Streams competitive
  (1M, heavy map)    │          │         │          │ when ops are expensive
  ───────────────────┼──────────┼─────────┼──────────┼──────────────────────
  GroupBy + Count    │ ~25ms    │ ~27ms   │ ~12ms    │ Collectors very
  (1M elements)      │          │         │          │ efficient
  ───────────────────┼──────────┼─────────┼──────────┼──────────────────────
  String concat      │ ~8ms     │ ~45ms   │ ~60ms    │ Streams + String
  (100K strings)     │ (SB)    │ joining │ joining  │ have more overhead
  ───────────────────┼──────────┼─────────┼──────────┼──────────────────────
  Find first match   │ ~0.01ms  │ ~0.02ms │ ~0.5ms   │ Parallel overhead
  (early termination)│          │         │          │ dominates for small work
  ───────────────────┴──────────┴─────────┴──────────┴──────────────────────
  
  Key takeaway:
  • Loops: 2-5x faster for simple operations on small-medium data
  • Streams: competitive for complex transformations; ~10-20% overhead
  • Parallel: only beneficial for CPU-intensive ops on >50K elements
  • Readability usually trumps micro-performance differences
```

### 8.4 GC Pressure from Streams

```java
// Stream pipeline allocations:
List<R> result = list.stream()          // 1 Head, 1 Spliterator
    .filter(predicate)                   // 1 FilterOp stage
    .map(function)                       // 1 MapOp stage
    .collect(toList());                  // 1 TerminalOp, 1 ArrayList, Sink chain

// Total: ~8-12 objects per pipeline execution (negligible for most apps)
// BUT: if called 1M times/sec in a hot loop → 8-12M objects/sec → GC pressure

// flatMap is the most allocation-heavy operation:
orders.stream()
    .flatMap(o -> o.getItems().stream())  // New stream + spliterator PER element
    .collect(toList());
// For 100K orders: 100K extra Stream objects + 100K Spliterators

// Java 16+ alternative: mapMulti (no intermediate streams)
orders.stream()
    .<Item>mapMulti((order, consumer) -> {
        order.getItems().forEach(consumer);
    })
    .collect(toList());
// ZERO intermediate stream allocations
```

### 8.5 Escape Analysis & Stream Optimization

```
JIT Escape Analysis for Streams:

  Ideal case (may be scalar-replaced):
    Optional<String> result = Optional.of("hello")
        .map(String::toUpperCase)
        .filter(s -> s.length() > 3);
    // Optional intermediaries may be eliminated by escape analysis
    // if the entire chain is inlined and Optional doesn't escape

  Typical case (escape analysis FAILS):
    List<String> result = list.stream()     // Pipeline head stored in local
        .filter(s -> s.length() > 3)         // New StatelessOp stage (escapes to pipeline)
        .map(String::toUpperCase)            // Another stage (escapes)
        .collect(toList());                  // Terminal processes pipeline
    // Pipeline stages reference each other → escape to linked list → heap allocated
    // JIT generally CANNOT eliminate stream pipeline objects

  Practical impact:
  - Stream pipeline objects: small, short-lived → collected in Minor GC
  - For hot paths (>100K calls/sec): consider loops to eliminate this overhead
  - For business logic (10-1000 calls/sec): stream overhead is irrelevant
```

### 8.6 Microbenchmarking with JMH

```java
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.MICROSECONDS)
@Warmup(iterations = 5, time = 1)
@Measurement(iterations = 10, time = 1)
@Fork(2)
@State(Scope.Benchmark)
public class StreamVsLoopBenchmark {

    @Param({"100", "10000", "1000000"})
    private int size;
    
    private List<Integer> data;
    
    @Setup
    public void setup() {
        data = new Random(42).ints(size, 0, 1000).boxed().collect(toList());
    }
    
    @Benchmark
    public long forLoop() {
        long sum = 0;
        for (int n : data) {
            if (n > 500) sum += n * 2L;
        }
        return sum;
    }
    
    @Benchmark
    public long streamSequential() {
        return data.stream()
            .filter(n -> n > 500)
            .mapToLong(n -> n * 2L)
            .sum();
    }
    
    @Benchmark
    public long streamParallel() {
        return data.parallelStream()
            .filter(n -> n > 500)
            .mapToLong(n -> n * 2L)
            .sum();
    }
    
    @Benchmark
    public long intStreamOptimized() {
        return data.stream()
            .mapToInt(Integer::intValue)
            .filter(n -> n > 500)
            .mapToLong(n -> n * 2L)
            .sum();
    }
}

// Run: java -jar target/benchmarks.jar StreamVsLoop -f 2 -wi 5 -i 10
```

**JMH Common Mistakes:**
- Benchmarking without warmup → measuring interpreter, not JIT-compiled code
- Not using `@Fork` → JIT optimizations from previous benchmarks bleed over
- Dead code elimination → JIT removes computation whose result isn't used (use `Blackhole.consume()`)
- Loop hoisting → JIT moves invariant computation outside the loop
- Not testing at realistic data sizes

### 8.7 Guidelines for High-Throughput Systems

| Guideline | Rationale |
|-----------|-----------|
| Use primitive streams (`IntStream`, `LongStream`) for numeric processing | Eliminates autoboxing overhead entirely |
| Prefer `collect()` over `forEach()` with side effects | Thread-safe, composable, GC-friendly |
| Avoid `parallelStream()` for I/O-bound operations | Blocks ForkJoinPool threads; use async I/O |
| Use `toArray()` over `collect(toList())` when array is sufficient | Direct array allocation, less overhead |
| Replace `flatMap` with `mapMulti` (Java 16+) for high-volume flattening | Eliminates per-element Stream allocation |
| Pre-size result collections: `Collectors.toCollection(() -> new ArrayList<>(expectedSize))` | Avoids internal array resizing |
| Extract hot-path lambdas to `static final` fields | Ensures non-capturing singleton optimization |
| Profile before optimizing — use JFR allocation profiling | Avoid premature optimization; measure actual hotspots |

```java
// Hot-path optimization: extract non-capturing lambda
private static final Predicate<Transaction> IS_VALID = 
    t -> t.getAmount().compareTo(BigDecimal.ZERO) > 0 && t.getStatus() == Status.ACTIVE;

private static final Function<Transaction, TransactionDTO> TO_DTO = 
    t -> new TransactionDTO(t.getId(), t.getAmount(), t.getTimestamp());

// These are guaranteed singletons — zero allocation per use
List<TransactionDTO> dtos = transactions.stream()
    .filter(IS_VALID)
    .map(TO_DTO)
    .collect(toList());
```

---

## 9. Advanced Topics

### 9.1 CompletableFuture — Functional Async Chaining

`CompletableFuture` is Java's monadic abstraction for asynchronous programming, supporting `map` (thenApply), `flatMap` (thenCompose), and error handling.

```java
// Full async pipeline with error handling and composition:
public CompletableFuture<OrderSummary> processOrderAsync(String orderId) {
    
    return fetchOrder(orderId)                              // CF<Order>
        // map: synchronous transformation
        .thenApply(order -> validateOrder(order))            // CF<ValidatedOrder>
        
        // flatMap: async transformation (returns another CF)
        .thenCompose(valid -> enrichWithPricing(valid))      // CF<PricedOrder>
        
        // parallel composition: combine two async results
        .thenCombine(
            fetchInventory(orderId),                         // CF<Inventory>
            (priced, inventory) -> reserveStock(priced, inventory)  // CF<ReservedOrder>
        )
        
        // apply async: run transformation on different executor
        .thenApplyAsync(reserved -> calculateTaxes(reserved), cpuExecutor)
        
        // error recovery
        .exceptionally(ex -> {
            log.error("Order processing failed: {}", orderId, ex);
            metricsService.incrementFailure("order.processing");
            return OrderSummary.failed(orderId, ex.getMessage());
        })
        
        // timeout (Java 9+)
        .orTimeout(5, TimeUnit.SECONDS)
        
        // post-processing
        .whenComplete((result, error) -> {
            if (error != null) alertService.notifyFailure(orderId, error);
        });
}
```

```
CompletableFuture Method Mapping to FP Concepts:

  FP Concept       │ CompletableFuture Method  │ Signature
  ─────────────────┼───────────────────────────┼──────────────────────────
  map              │ thenApply                 │ CF<T> → (T→U) → CF<U>
  flatMap          │ thenCompose               │ CF<T> → (T→CF<U>) → CF<U>
  zip / combine    │ thenCombine               │ CF<T> × CF<U> → ((T,U)→V) → CF<V>
  onError          │ exceptionally             │ CF<T> → (Throwable→T) → CF<T>
  recover          │ handle                    │ CF<T> → ((T,Throwable)→U) → CF<U>
  tap / peek       │ whenComplete              │ CF<T> → ((T,Throwable)→void) → CF<T>
  race (first)     │ anyOf                     │ CF<?>... → CF<Object>
  all              │ allOf                     │ CF<?>... → CF<Void>
```

**Production Patterns:**

```java
// Fan-out / Fan-in: multiple parallel calls with timeout
public CompletableFuture<DashboardData> fetchDashboard(String userId) {
    CompletableFuture<UserProfile> profile = userService.getProfile(userId);
    CompletableFuture<List<Order>> orders = orderService.getRecent(userId);
    CompletableFuture<WalletBalance> wallet = walletService.getBalance(userId);
    
    return CompletableFuture.allOf(profile, orders, wallet)
        .thenApply(v -> new DashboardData(
            profile.join(),   // Safe — allOf guarantees completion
            orders.join(),
            wallet.join()
        ))
        .orTimeout(3, TimeUnit.SECONDS)
        .exceptionally(ex -> DashboardData.partial(userId));
}

// ⚠ Always specify executor for I/O operations:
.thenApplyAsync(this::blockingCall, ioExecutor)
// Default: ForkJoinPool.commonPool() — DO NOT use for blocking I/O
```

### 9.2 Functional Style in Reactive Programming

```java
// Reactor (Spring WebFlux) — functional stream processing over async data
@Service
public class OrderReactiveService {
    
    public Flux<OrderEvent> processOrders(Flux<Order> orderStream) {
        return orderStream
            .filter(order -> order.getTotal().compareTo(MIN_AMOUNT) > 0)
            .flatMap(order -> enrichOrder(order)            // Mono<EnrichedOrder>
                .timeout(Duration.ofSeconds(2))
                .onErrorResume(ex -> Mono.just(EnrichedOrder.fallback(order))))
            .groupBy(EnrichedOrder::getRegion)
            .flatMap(regionGroup -> regionGroup
                .buffer(Duration.ofSeconds(1))              // Batch by time window
                .flatMap(batch -> processBatch(batch)))     // Mono<List<OrderEvent>>
            .doOnNext(event -> metricsService.record(event))
            .subscribeOn(Schedulers.boundedElastic());
    }
}

// Key Reactor operators mapped to FP concepts:
// Mono  = Optional  (0 or 1 element, async)
// Flux  = Stream    (0 to N elements, async, backpressure-aware)
// map, flatMap, filter, reduce — same semantics, async execution
// Schedulers = Executor abstraction for thread control
```

### 9.3 Functional Configuration in Spring

```java
// Spring functional bean registration (alternative to @Bean / @Component)
public class AppConfig implements ApplicationContextInitializer<GenericApplicationContext> {
    @Override
    public void initialize(GenericApplicationContext ctx) {
        ctx.registerBean(UserRepository.class, () -> new JdbcUserRepository(ctx.getBean(DataSource.class)));
        ctx.registerBean(UserService.class, () -> new UserService(ctx.getBean(UserRepository.class)));
        ctx.registerBean(UserController.class, () -> new UserController(ctx.getBean(UserService.class)));
    }
}
// Benefits: no reflection, no component scanning, faster startup (GraalVM-friendly)

// Spring WebFlux functional endpoints (alternative to @Controller):
@Bean
public RouterFunction<ServerResponse> routes(UserHandler handler) {
    return RouterFunctions.route()
        .GET("/users/{id}", handler::getUser)
        .GET("/users", handler::listUsers)
        .POST("/users", handler::createUser)
        .filter((request, next) -> {
            log.info("Request: {} {}", request.method(), request.path());
            return next.handle(request);
        })
        .build();
}
```

### 9.4 Records & Pattern Matching — FP Alignment

```java
// Records (Java 16+) — immutable data carriers aligned with FP principles
public record Point(double x, double y) {
    // Compact canonical constructor for validation
    public Point {
        if (Double.isNaN(x) || Double.isNaN(y)) throw new IllegalArgumentException("NaN");
    }
    
    // Derived computation (pure function)
    public double distanceTo(Point other) {
        return Math.sqrt(Math.pow(this.x - other.x, 2) + Math.pow(this.y - other.y, 2));
    }
}

// Pattern matching for switch (Java 21+) — algebraic data type decomposition
sealed interface Shape permits Circle, Rectangle, Triangle {}
record Circle(double radius) implements Shape {}
record Rectangle(double width, double height) implements Shape {}
record Triangle(double base, double height) implements Shape {}

double area(Shape shape) {
    return switch (shape) {
        case Circle c    -> Math.PI * c.radius() * c.radius();
        case Rectangle r -> r.width() * r.height();
        case Triangle t  -> 0.5 * t.base() * t.height();
        // Exhaustive: compiler ensures all cases covered (sealed)
    };
}

// Record patterns (deconstruction) — Java 21+:
String describe(Object obj) {
    return switch (obj) {
        case Point(var x, var y) when x == 0 && y == 0 -> "Origin";
        case Point(var x, var y) -> "Point(%s, %s)".formatted(x, y);
        case Circle(var r) -> "Circle with radius " + r;
        default -> obj.toString();
    };
}
```

### 9.5 Virtual Threads & FP Considerations (Java 21+)

```java
// Virtual threads change the calculus for async FP:
// Before: CompletableFuture chains needed to avoid blocking ForkJoinPool threads
// After: Virtual threads make blocking I/O cheap — simple sequential code is OK

// Pre-virtual threads (async boilerplate):
CompletableFuture<User> future = CompletableFuture
    .supplyAsync(() -> httpClient.get("/users/" + id), ioExecutor)
    .thenApply(response -> mapper.readValue(response, User.class));

// With virtual threads (simple blocking code, each on its own virtual thread):
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    List<Future<User>> futures = userIds.stream()
        .map(id -> executor.submit(() -> {
            String response = httpClient.get("/users/" + id);  // Blocking is FINE
            return mapper.readValue(response, User.class);
        }))
        .toList();
    
    List<User> users = futures.stream()
        .map(f -> { try { return f.get(); } catch (Exception e) { throw new RuntimeException(e); } })
        .toList();
}

// Structured concurrency (Preview, Java 21+):
try (var scope = new StructuredTaskScope.ShutdownOnFailure()) {
    Subtask<User> userTask = scope.fork(() -> userService.fetch(userId));
    Subtask<List<Order>> ordersTask = scope.fork(() -> orderService.fetch(userId));
    
    scope.join();           // Wait for both
    scope.throwIfFailed();  // Propagate first failure
    
    return new Dashboard(userTask.get(), ordersTask.get());
}
// Benefits: clear parent-child task relationship, automatic cancellation,
// structured error propagation — cleaner than CompletableFuture chains
```

**Impact on FP patterns:**
- `CompletableFuture` chains become less necessary for I/O coordination — simpler sequential code on virtual threads is often clearer
- `parallelStream()` still uses `ForkJoinPool` (platform threads) — not affected by virtual threads
- Reactive frameworks (WebFlux, RxJava) still relevant for backpressure scenarios, but virtual threads reduce the need for reactive in pure I/O concurrency

### 9.6 Gatherers (Java 22+ Preview) — Custom Intermediate Operations

```java
// Gatherers let you write custom intermediate stream operations
// (currently equivalent to what Collector is for terminal operations)

// Built-in Gatherers:
Stream.of(1, 2, 3, 4, 5)
    .gather(Gatherers.windowFixed(3))   // [[1,2,3], [4,5]]
    .toList();

Stream.of(1, 2, 3, 4, 5)
    .gather(Gatherers.windowSliding(3)) // [[1,2,3], [2,3,4], [3,4,5]]
    .toList();

// Custom Gatherer: rate-limited emission
Gatherer<T, ?, T> rateLimiter(Duration minInterval) {
    return Gatherer.ofSequential(
        () -> new AtomicLong(0),  // state: last emission time
        (state, element, downstream) -> {
            long now = System.nanoTime();
            long last = state.get();
            if (now - last >= minInterval.toNanos()) {
                state.set(now);
                return downstream.push(element);
            }
            return true;  // continue but don't emit
        }
    );
}
```

### 9.7 Type-Driven Development & Functional Architecture

Type-driven development uses the **type system as documentation and enforcement** — making invalid states unrepresentable at compile time.

#### Making Illegal States Unrepresentable

```java
// ❌ Stringly-typed / primitive-obsessed API:
void processPayment(String orderId, double amount, String currency, String status) { }
// Any String can go in any parameter — bugs caught only at runtime

// ✅ Type-driven API: compiler prevents misuse
record OrderId(String value) {}
record Money(BigDecimal amount, Currency currency) {
    public Money { 
        if (amount.compareTo(BigDecimal.ZERO) < 0) throw new IllegalArgumentException("Negative");
    }
}
sealed interface PaymentStatus permits Pending, Completed, Failed {}

void processPayment(OrderId id, Money amount, PaymentStatus status) { }
// Wrong argument order → COMPILE ERROR
// Negative amount → PREVENTED by record constructor
// Invalid status → PREVENTED by sealed type
```

#### State Machine via Sealed Types

```java
// Model a workflow where each state has different valid operations:
sealed interface OrderState permits Draft, Submitted, Approved, Shipped {}

record Draft(OrderId id, List<LineItem> items) implements OrderState {
    Submitted submit() { 
        if (items.isEmpty()) throw new IllegalStateException("Empty order");
        return new Submitted(id, items, Instant.now()); 
    }
    // Can ONLY submit a Draft — can't ship or approve it
}

record Submitted(OrderId id, List<LineItem> items, Instant submittedAt) implements OrderState {
    Approved approve(UserId approver) { return new Approved(id, items, approver, Instant.now()); }
    Draft reject(String reason) { return new Draft(id, items); }  // Back to draft
    // Can't ship a Submitted order — must approve first
}

record Approved(OrderId id, List<LineItem> items, UserId approver, Instant approvedAt) implements OrderState {
    Shipped ship(String trackingNumber) { return new Shipped(id, trackingNumber, Instant.now()); }
    // Only Approved orders can be shipped — enforced by type system!
}

record Shipped(OrderId id, String trackingNumber, Instant shippedAt) implements OrderState {
    // Terminal state — no further transitions. Compiler enforces this.
}

// Usage:
OrderState result = switch (currentState) {
    case Draft d -> d.submit();
    case Submitted s -> s.approve(currentUser);
    case Approved a -> a.ship(generateTracking());
    case Shipped s -> s;  // Already shipped — no-op
};
// Compiler guarantees EXHAUSTIVE handling of all states
```

#### Functional Architecture Patterns for Enterprise Java

```
Hexagonal / Ports and Adapters with FP:

  ┌──────────────────────────────────────────────────────────────┐
  │                    Application Layer                          │
  │  ┌──────────────────────────────────────────────────────┐    │
  │  │              FUNCTIONAL CORE                          │    │
  │  │  (Pure domain logic — no framework dependencies)      │    │
  │  │                                                        │    │
  │  │  Function<OrderRequest, Either<Error, Order>> create  │    │
  │  │  Function<Order, Either<Error, Invoice>> invoice       │    │
  │  │  Predicate<Order> validate                             │    │
  │  │  BiFunction<Order, Discount, Order> applyDiscount      │    │
  │  │                                                        │    │
  │  │  • Receives immutable inputs, returns immutable outputs │    │
  │  │  • ALL dependencies injected as function parameters    │    │
  │  │  • Zero side effects — testable with just assertions   │    │
  │  └──────────────────────────────────────────────────────┘    │
  │          ▲                                    ▲               │
  │          │  Ports (interfaces)                │               │
  │  ┌───────┴──────┐                  ┌──────────┴──────────┐   │
  │  │   Input       │                  │   Output             │   │
  │  │   Adapters    │                  │   Adapters           │   │
  │  │ (Controllers, │                  │ (Repositories,       │   │
  │  │  CLI, Events) │                  │  Messaging, HTTP)    │   │
  │  └──────────────┘                  └─────────────────────┘   │
  └──────────────────────────────────────────────────────────────┘

  Key principle: "Dependency injection via function parameters"
  Instead of:  @Autowired OrderRepository repo;   // field injection
  Use:         Function<OrderId, Optional<Order>> findOrder  // function parameter
  
  Benefits:
  • Core logic is portable — no Spring/Jakarta dependency
  • Functions compose — buildPipeline(validate, enrich, price, persist)
  • Testing needs only lambdas, not mock frameworks
  • GraalVM native image compatible (no reflection in core)
```

```java
// Practical example — functional service composition:
public class OrderProcessor {
    private final Function<OrderId, Optional<Order>>    findOrder;
    private final Function<Order, Either<Error, Order>> validate;
    private final Function<Order, Order>                enrich;
    private final Function<Order, Either<Error, Order>> persist;
    
    // All dependencies are FUNCTIONS — injected via constructor
    public OrderProcessor(
            Function<OrderId, Optional<Order>> findOrder,
            Function<Order, Either<Error, Order>> validate,
            Function<Order, Order> enrich,
            Function<Order, Either<Error, Order>> persist) {
        this.findOrder = findOrder;
        this.validate = validate;
        this.enrich = enrich;
        this.persist = persist;
    }
    
    public Either<Error, Order> process(OrderId id) {
        return findOrder.apply(id)
            .map(Either::<Error, Order>right)
            .orElse(Either.left(new Error("Not found")))
            .map(enrich)
            .flatMap(validate)
            .flatMap(persist);
    }
}

// In tests — NO mocks needed, just lambdas:
var processor = new OrderProcessor(
    id -> Optional.of(testOrder),              // stub findOrder
    order -> Either.right(order),              // stub validate (always pass)
    order -> order,                            // stub enrich (identity)
    order -> Either.right(order)               // stub persist (always succeed)
);
assertEquals(Either.right(testOrder), processor.process(testId));
```

#### Effect Systems (Conceptual — Beyond Java)

In pure FP languages, **effect systems** track what side effects a function can perform in its **type signature**. Java doesn't have this, but understanding the concept helps reason about FP design:

```
Haskell example (for comparison):
  
  pure computation:    add :: Int -> Int -> Int
  with I/O effect:     readFile :: FilePath -> IO String
  with state effect:   counter :: State Int Int
  with failure effect: lookup :: Key -> Maybe Value
  
  The TYPE tells you what effects the function can have.
  A function without IO in its type CANNOT do I/O — compiler enforced.

Java approximation (conventions, not compiler-enforced):
  
  Function<A, B>                    → Pure computation (by convention)
  Function<A, CompletableFuture<B>> → Async computation
  Function<A, Optional<B>>          → Computation that may fail silently
  Function<A, Either<Err, B>>       → Computation that may fail with reason
  Function<A, Stream<B>>            → Computation producing multiple values
  Function<A, Mono<B>>              → Async computation with backpressure
  
  Pattern: The RETURN TYPE communicates what kind of "effect" the function has.
  This is why returning Either<Error, Result> is more informative than throws Exception.
  The caller's code visibly handles the error case (map/flatMap) vs invisible try/catch.
```

---

## 10. Senior-Level Interview Questions & Answers

### Part A: Internal & JVM-Level Questions

---

#### Q1: How are lambdas implemented at the bytecode level? Why doesn't Java create an anonymous inner class for each lambda?

**Expected Answer:**

Java compiles lambdas using `invokedynamic` (introduced in Java 7, leveraged for lambdas in Java 8). The process:

1. **Compile time:** The compiler generates an `invokedynamic` instruction pointing to a bootstrap method (`LambdaMetafactory.metafactory`). The lambda body is compiled into a **private static method** (`lambda$methodName$0`) inside the enclosing class.

2. **First invocation (linkage):** The bootstrap method runs once. `LambdaMetafactory` uses `ASM` to dynamically generate a class at runtime that implements the target functional interface and delegates to the desugared static method.

3. **Subsequent invocations:** The generated `CallSite` (linked via `MethodHandle`) is reused — no further class generation.

**Why not anonymous inner classes?**
- Anonymous classes create a `.class` file per lambda at compile time → metaspace bloat
- Anonymous classes always allocate a new object on the heap
- `invokedynamic` allows the JVM to choose the best strategy at runtime:
  - **Stateless lambdas** (no captures): singleton instance (zero allocation)
  - **Capturing lambdas**: new instance per invocation (lightweight — no outer class `this$0` reference unless needed)
  - Future JVMs can change strategy without recompiling bytecode

**Key bytecode to verify:**
```bash
javap -c -p MyClass.class
# Look for: invokedynamic #0:run  (bootstrap: LambdaMetafactory.metafactory)
# Look for: private static lambda$main$0()V  (desugared body)
```

**Common Mistake:** Saying "lambdas are syntactic sugar for anonymous classes." They are NOT — the implementation mechanism is fundamentally different.

**Follow-up Probes:**
- What happens if a lambda captures a local variable vs. an instance field?
- Can you explain how `MethodHandle` differs from reflection?
- How does `LambdaMetafactory.altMetafactory` differ from `metafactory`?

---

#### Q2: What is the difference between a stateless and a capturing lambda? What are the memory implications?

**Expected Answer:**

| Aspect | Stateless Lambda | Capturing Lambda |
|--------|-----------------|-----------------|
| Captures variables? | No | Yes (local vars or `this`) |
| Instance allocation | Singleton (reused) | New object per evaluation |
| GC pressure | Zero | Per-invocation allocation |
| Desugared method | `private static` | `private static` (captured values passed as args) |
| Generated class fields | None | Fields for each captured variable |

```java
// Stateless — singleton, zero allocation:
Predicate<String> notEmpty = s -> !s.isEmpty();
// Same instance reused every time this line executes

// Capturing local variable — new instance per evaluation:
String prefix = "Hello";
Function<String, String> greeter = name -> prefix + " " + name;
// New instance each time, holds reference to "Hello" string

// Capturing 'this' — new instance, holds enclosing object reference:
public Predicate<Order> isExpensive() {
    return order -> order.getTotal() > this.threshold;
    // Captures 'this' → prevents GC of enclosing object if lambda is stored
}
```

**Senior-level concern:** Capturing `this` in a lambda stored in a long-lived collection (cache, listener registry) creates a **memory leak** — the entire enclosing object is retained.

**Follow-up Probes:**
- How would you detect this kind of leak in a heap dump?
- What does "effectively final" mean and why is it required for captures?

---

#### Q3: Explain how a Stream pipeline executes internally. What is the Sink chain?

**Expected Answer:**

A Stream pipeline has three parts: **Source → Intermediate ops → Terminal op**. Nothing executes until the terminal operation is invoked (lazy evaluation).

**Internal execution model:**

```
Source → op1 → op2 → op3 → terminal
                                   
Internally becomes a Sink chain (reverse wiring):

terminal.begin(size)
  └→ op3.Sink.begin()
       └→ op2.Sink.begin()
            └→ op1.Sink.begin()

For each element from source:
  op1.Sink.accept(element)
    └→ op2.Sink.accept(transformed)
       └→ op3.Sink.accept(filtered)
          └→ terminal.Sink.accept(result)

terminal.end()
  └→ cascading end() calls
```

**Key implementation classes:**
- `AbstractPipeline` — linked list of stages (head → intermediate → intermediate)
- `Sink<T>` — the element-processing callback with `begin(long size)`, `accept(T)`, `end()`, `cancellationRequested()`
- `Spliterator` — provides elements, supports splitting for parallelism
- `TerminalOp` — triggers pipeline wiring and execution

**Operation fusion:** The JVM does NOT create intermediate collections. Each element flows through the entire Sink chain before the next element is processed (depth-first, not breadth-first).

**Stateful vs. stateless:**
- Stateless (`filter`, `map`): process one element at a time, no buffering
- Stateful (`sorted`, `distinct`): must see all elements before emitting → creates internal buffer → breaks pure streaming

**Common Mistake:** Thinking `filter().map().collect()` creates three intermediate lists. It doesn't — elements flow one-at-a-time through fused Sink chain.

**Follow-up Probes:**
- What happens when you call `peek()` — does it add a pipeline stage?
- How does short-circuiting (like `findFirst()`) work through the Sink chain?
- What is `StreamShape` and why does it matter for primitive specialization?

---

#### Q4: What role does Spliterator play in parallel streams? How does splitting affect performance?

**Expected Answer:**

`Spliterator` (Splittable Iterator) is the source abstraction for streams. It provides:

1. **`tryAdvance(Consumer)`** — process one element
2. **`forEachRemaining(Consumer)`** — bulk traversal (optimized)
3. **`trySplit()`** — split into two halves for parallel processing
4. **`estimateSize()`** — hint for splitting strategy
5. **`characteristics()`** — bitfield (ORDERED, DISTINCT, SORTED, SIZED, NONNULL, IMMUTABLE, CONCURRENT, SUBSIZED)

**Splitting for parallel streams:**

```
Original Spliterator [1, 2, 3, 4, 5, 6, 7, 8]
                    ↓ trySplit()
 Left [1, 2, 3, 4]          Right [5, 6, 7, 8]
       ↓ trySplit()                ↓ trySplit()
 [1, 2]    [3, 4]           [5, 6]    [7, 8]
   ↓          ↓               ↓          ↓
 Thread1   Thread2         Thread3   Thread4
   ↓          ↓               ↓          ↓
 result1   result2         result3   result4
       ↓                         ↓
     combine               combine
              ↓
          final result
```

**Splitting efficiency by data source:**

| Source | Splitting Quality | Reason |
|--------|------------------|--------|
| `ArrayList`, arrays | Excellent | Index-based split at midpoint, O(1) |
| `HashSet`, `HashMap` | Good | Bucket-based splitting |
| `TreeSet`, `TreeMap` | Good | Structural splitting |
| `LinkedList` | Terrible | Must traverse to find midpoint, O(n) |
| `Stream.iterate()` | Terrible | Unknown size, sequential dependency |
| `BufferedReader.lines()` | Poor | I/O bound, sequential reading |

**Characteristics matter:**
- `SIZED + SUBSIZED`: enables perfect splitting — each sub-spliterator knows its exact size
- `ORDERED`: forces ordered merge in parallel → constrains parallelism
- `IMMUTABLE`/`CONCURRENT`: safe for parallel without `CopyOnWriteArrayList`

**Follow-up Probes:**
- How would you write a custom Spliterator for a database cursor?
- What happens when characteristics are wrong (e.g., claiming SORTED when not)?
- How does `SUBSIZED` differ from `SIZED`?

---

#### Q5: What is the `invokedynamic` call site caching mechanism for lambdas? Can it fail?

**Expected Answer:**

When the JVM encounters `invokedynamic` for a lambda:

1. **Bootstrap (first call):** `LambdaMetafactory.metafactory()` is invoked. It generates a hidden class implementing the functional interface. A `ConstantCallSite` (for non-capturing) or `MutableCallSite` (theoretically) is returned.

2. **Caching:** The `CallSite` is linked permanently. On subsequent executions, the JVM jumps directly to the target `MethodHandle` — no bootstrap re-invocation.

3. **Non-capturing lambdas:** The `CallSite` returns the **same singleton instance** every time. The target `MethodHandle` is a constant that returns a cached instance.

4. **Capturing lambdas:** The `CallSite` target is a factory `MethodHandle` that creates a new instance with the captured values as constructor arguments.

**Can it fail or deoptimize?**
- The call site is a `ConstantCallSite` — it cannot be relinked. It will not "fail" under normal circumstances.
- However, `LambdaMetafactory` can throw `LambdaConversionException` at link time if there's a type mismatch.
- Under extreme Metaspace pressure, the generated hidden class could trigger `OutOfMemoryError: Metaspace`.
- Serializable lambdas use `altMetafactory` which has more complex (and slightly slower) linking.

**Common Mistake:** Confusing `invokedynamic` with reflection. `invokedynamic` has NO runtime reflection overhead — it's resolved to a direct method handle at link time.

**Follow-up Probes:**
- What's the difference between `ConstantCallSite` and `MutableCallSite`?
- Where are the generated lambda classes stored in the JVM memory model?
- How does `-Djdk.internal.lambda.dumpProxyClasses` help in debugging?

---

### Part B: Performance & Debugging Scenario Questions

---

#### Q6: A production service shows high GC pressure after migrating from for-loops to Streams. How do you diagnose and fix it?

**Expected Answer (structured troubleshooting):**

**Step 1: Confirm the hypothesis**
```bash
# Enable GC logging:
-Xlog:gc*:file=gc.log:time,uptime,level,tags:filecount=5,filesize=50m

# Check allocation rate with JFR:
jcmd <pid> JFR.start duration=60s filename=alloc.jfr
# Analyze: jfr print --events jdk.ObjectAllocationInNewTLAB alloc.jfr
```

**Step 2: Identify common Stream GC culprits:**

| Root Cause | What Happens | Fix |
|-----------|-------------|-----|
| **Autoboxing** | `IntStream` → `Stream<Integer>` via `.boxed()` or `.mapToObj()` creates wrapper objects | Use primitive specializations: `IntStream`, `LongStream`, `DoubleStream` |
| **Lambda captures** | Capturing lambdas create new object per evaluation | Extract to static method references or stateless lambdas |
| **Unnecessary Stream creation** | `collection.stream().forEach()` instead of `collection.forEach()` | Use `Iterable.forEach()` when no pipeline needed |
| **`flatMap` explosion** | `flatMap(x -> Stream.of(...))` creates Stream object per element | Use `mapMulti()` (Java 16+) — reuses downstream consumer |
| **`Collectors.toList()`** | Creates `ArrayList` with resizing | Use `.toList()` (Java 16+) — returns compact unmodifiable list |
| **String concatenation in map** | `map(x -> x.name + ":" + x.value)` creates intermediate Strings | Use `StringBuilder` in collect, or `String.formatted()` |

**Step 3: JMH benchmark before/after:**
```java
@Benchmark
public int streamSum(Blackhole bh) {
    return list.stream().mapToInt(Integer::intValue).sum(); // Primitive stream
}

@Benchmark  
public int loopSum() {
    int sum = 0;
    for (int val : list) sum += val;
    return sum;
}
// Compare throughput and gc.alloc.rate metrics
```

**Step 4: Production mitigation if refactoring is costly:**
- Increase Young Gen (`-Xmn`) to accommodate higher allocation rate
- Switch to ZGC/Shenandoah for lower pause times
- Profile with `async-profiler` in allocation mode: `./asprof -e alloc -d 30 <pid>`

**Common Mistake:** Blindly converting all streams back to loops. Most stream pipelines have negligible overhead — the issue is usually one or two hot paths with autoboxing or excessive `flatMap`.

**Follow-up Probes:**
- How does escape analysis interact with lambda allocations?
- What is TLAB and how does it mitigate short-lived object allocation?
- Can you show a `mapMulti` replacement for `flatMap`?

---

#### Q7: A parallel stream operation intermittently produces incorrect results. How do you debug it?

**Expected Answer:**

**Root causes (in order of likelihood):**

1. **Shared mutable state (most common):**
```java
// BUG: ArrayList is not thread-safe
List<Result> results = new ArrayList<>();
data.parallelStream().forEach(item -> results.add(process(item)));  // Race condition!

// FIX: Use Collector (thread-safe reduction)
List<Result> results = data.parallelStream()
    .map(this::process)
    .collect(Collectors.toList());
```

2. **Non-associative reduction:**
```java
// BUG: subtraction is not associative
int result = numbers.parallelStream().reduce(0, (a, b) -> a - b);
// Sequential: ((0-1)-2)-3 = -6
// Parallel: might compute (0-1) and (0-2) then combine: (-1) - (-2) = 1  ← WRONG

// FIX: Use only associative operators (add, multiply, min, max, concat)
```

3. **Identity violation in reduce:**
```java
// BUG: 10 is not the identity for addition
int result = numbers.parallelStream().reduce(10, Integer::sum);
// Each parallel partition starts with 10 → result is 10 * numPartitions + actual sum

// FIX: Use 0 as identity, add the constant separately
int result = numbers.parallelStream().reduce(0, Integer::sum) + 10;
```

4. **Stateful lambda in intermediate op:**
```java
// BUG: shared AtomicInteger for indexing — order undefined in parallel
AtomicInteger counter = new AtomicInteger();
list.parallelStream()
    .map(item -> counter.getAndIncrement() + ":" + item)  // Indices are random!
    .collect(Collectors.toList());
```

**Debugging approach:**
```java
// 1. Add thread visibility:
.peek(item -> System.out.println(Thread.currentThread().getName() + ": " + item))

// 2. Compare sequential vs parallel results:
List<Result> sequential = data.stream().map(this::process).toList();
List<Result> parallel = data.parallelStream().map(this::process).toList();
assert sequential.equals(parallel); // If this fails → side-effect or ordering issue

// 3. Use -Djava.util.concurrent.ForkJoinPool.common.parallelism=1
// Forces single thread — if bug disappears, it's a concurrency issue
```

**Follow-up Probes:**
- What is the non-interference contract for stream sources?
- How would you write a custom thread-safe Collector for this scenario?
- Can `ConcurrentHashMap` be safely used as a stream source during parallel processing?

---

#### Q8: You're seeing `OutOfMemoryError` in a service that processes large CSV files using Streams. What's happening?

**Expected Answer:**

**Likely causes:**

1. **Holding reference to entire Stream in a `sorted()` or `distinct()` operation:**
```java
// BUG: sorted() must buffer ALL elements before emitting
Files.lines(path)          // Lazy — good
    .map(this::parseLine)
    .sorted(comparing(Record::timestamp))  // Buffers ALL lines in memory — OOM!
    .forEach(this::process);

// FIX: Sort the file externally or use a bounded buffer with external merge sort
// Or: If file fits in memory, load it intentionally:
List<Record> records = Files.lines(path).map(this::parseLine).toList();
records.sort(comparing(Record::timestamp));
records.forEach(this::process);
```

2. **`collect(Collectors.toList())` on unbounded stream:**
```java
// BUG: collecting millions of records into a list
List<Record> all = Files.lines(hugePath).map(this::parse).collect(Collectors.toList());

// FIX: Process in streaming fashion:
Files.lines(hugePath)
    .map(this::parse)
    .filter(this::isValid)
    .forEach(this::writeToDb);  // Process-and-forget each element
```

3. **`flatMap` creating massive intermediate streams:**
```java
// BUG: each line expands to many records, all held in memory
lines.flatMap(line -> Arrays.stream(line.split(",")).map(Record::new))

// FIX: Use mapMulti for controlled expansion
lines.<Record>mapMulti((line, consumer) -> {
    for (String part : line.split(",")) consumer.accept(new Record(part));
});
```

4. **Stream source not being closed (file handle leak → eventual OOM):**
```java
// BUG: Files.lines() returns a Stream that must be closed
Stream<String> lines = Files.lines(path);  // File handle stays open
lines.forEach(this::process);
// If called in a loop: file handles accumulate → OOM or "Too many open files"

// FIX: Always use try-with-resources
try (Stream<String> lines = Files.lines(path)) {
    lines.forEach(this::process);
}
```

**Diagnostic steps:**
```bash
# 1. Heap dump on OOM:
-XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/tmp/heapdump.hprof

# 2. In Eclipse MAT or VisualVM: look for large Object[] or ArrayList 
#    retained by AbstractPipeline or Collectors internal accumulators

# 3. Allocation profiling:
async-profiler -e alloc -d 60 --flamegraph alloc.html <pid>
```

**Follow-up Probes:**
- How does `Files.lines()` differ from `BufferedReader.readLine()` in terms of memory?
- What is back-pressure and how does it relate to Streams vs. Reactive Streams?
- How would you process a 50GB file in Java with constant memory?

---

#### Q9: Autoboxing in Streams — identify the performance problem and fix it:

```java
List<Integer> numbers = IntStream.rangeClosed(1, 1_000_000)
    .boxed()
    .collect(Collectors.toList());

double average = numbers.stream()
    .reduce(0, Integer::sum) / (double) numbers.size();
```

**Expected Answer:**

**Problems identified:**

1. **Unnecessary boxing:** `IntStream.rangeClosed()` produces primitives, but `.boxed()` converts to `Integer` objects — 1M wrapper objects (~16 bytes each = ~16MB heap allocation).

2. **Reduction with boxing:** `numbers.stream().reduce(0, Integer::sum)` operates on `Stream<Integer>` — each `sum` operation unboxes two `Integer`s and boxes the result back.

3. **Missing primitive specialization:** `Stream<Integer>` has no `average()` method — it's only on `IntStream`.

**Fixed version:**
```java
// Zero boxing, zero wrapper objects:
double average = IntStream.rangeClosed(1, 1_000_000)
    .average()
    .orElse(0.0);

// If List<Integer> is required elsewhere, compute stats from primitive stream:
IntSummaryStatistics stats = numbers.stream()
    .mapToInt(Integer::intValue)  // Unbox once, then stay in int domain
    .summaryStatistics();

double avg = stats.getAverage();
long sum = stats.getSum();
int max = stats.getMax();
```

**Performance impact:**
- Original: ~1M `Integer` allocations + ~1M unbox/box cycles in reduce
- Fixed: Zero allocations (primitive stack operations throughout)
- Benchmark difference: typically **3-5x faster** for large datasets

**Follow-up Probes:**
- What is `IntSummaryStatistics` and when would you use it?
- How does `mapToInt` differ from `map(Integer::intValue)` internally?
- What happens if you use `sum()` instead of `reduce(0, Integer::sum)` on `Stream<Integer>`?

---

#### Q10: A developer reports that `parallelStream()` is slower than sequential. What are the likely causes?

**Expected Answer:**

**Common causes (in order of impact):**

| Cause | Why It's Slower | Solution |
|-------|----------------|----------|
| **Small dataset** | Fork/join overhead > processing time | Threshold: >10K elements OR >1ms per element |
| **Shared ForkJoinPool saturation** | Other parallel streams or `CompletableFuture` tasks compete for same pool | Custom `ForkJoinPool` or virtual threads |
| **LinkedList source** | O(n) splitting — traverses list to find midpoint | Convert to array or ArrayList first |
| **I/O-bound operations** | Parallel doesn't help when bottleneck is disk/network | Use async I/O or virtual threads instead |
| **Ordered operations** | `forEachOrdered`, ordered `collect` forces sequential merge | Use `unordered()` hint or `forEach` |
| **Heavy stateful ops** | `sorted()` or `distinct()` require synchronization | Move sorting before/after parallel section |
| **Lock contention in processing** | Synchronized methods called from lambda | Redesign for lock-free processing |
| **False sharing** | Adjacent array elements processed by different threads | Chunk-based processing |

**Diagnostic approach:**
```java
// Quick test: time both approaches
long start = System.nanoTime();
result = data.stream().map(this::process).collect(toList());       // sequential
long seqTime = System.nanoTime() - start;

start = System.nanoTime();
result = data.parallelStream().map(this::process).collect(toList()); // parallel
long parTime = System.nanoTime() - start;

System.out.printf("Sequential: %dms, Parallel: %dms, Speedup: %.2fx%n",
    seqTime/1_000_000, parTime/1_000_000, (double)seqTime/parTime);

// If speedup < 1.0x → parallel is slower → don't use it
```

**Rule of thumb for parallel streams:**
```
NQ Model (Brian Goetz):
  N = number of elements
  Q = cost per element (computation time)

  If N * Q > 10,000 microseconds → parallel likely beneficial
  If N * Q < 1,000 microseconds → sequential is faster
```

**Follow-up Probes:**
- How do you submit a parallel stream to a custom ForkJoinPool?
- What does `Stream.unordered()` actually do to the pipeline?
- Can you have nested parallel streams and what happens to the ForkJoinPool?

---

### Part C: Architectural & Design Questions

---

#### Q11: When should you NOT use Streams? Give concrete enterprise scenarios.

**Expected Answer:**

**Avoid Streams when:**

| Scenario | Why Not Streams | Better Alternative |
|----------|----------------|-------------------|
| **Simple iteration with side effects** | Streams are designed for transformations, not imperative side effects | Enhanced for-loop or `Iterable.forEach()` |
| **Index-dependent processing** | No built-in index access — workarounds are ugly | Traditional for-loop |
| **Multi-step mutation of complex state** | Stream pipelines should be stateless; tracking state across elements is anti-pattern | Imperative loop with local variables |
| **Early exit with complex conditions** | Short-circuit ops (`findFirst`, `anyMatch`) exist but can't break mid-pipeline arbitrarily | `for` + `break`/`return` |
| **Exception-heavy processing** | Checked exceptions in lambdas require ugly wrappers | `for` loop with try-catch |
| **Performance-critical tight loops** | Stream overhead (pipeline setup, Sink chain, virtual dispatch) matters in nanosecond-sensitive code | Primitive arrays with index loops |
| **Debugging complex transformations** | Long stream chains are hard to step through in debugger | Break into named methods or use loops |

**Code example — when imperative is clearer:**
```java
// Stream version — hard to read and maintain:
Map<String, List<ValidationError>> errors = orders.stream()
    .flatMap(order -> validate(order).stream()
        .map(error -> Map.entry(order.getId(), error)))
    .collect(Collectors.groupingBy(
        Map.Entry::getKey,
        Collectors.mapping(Map.Entry::getValue, Collectors.toList())));

// Imperative version — clearer intent:
Map<String, List<ValidationError>> errors = new LinkedHashMap<>();
for (Order order : orders) {
    List<ValidationError> orderErrors = validate(order);
    if (!orderErrors.isEmpty()) {
        errors.put(order.getId(), orderErrors);
    }
}
```

**Senior-level principle:** "Use streams for data transformation pipelines. Use loops for procedural logic with side effects. The goal is readability and maintainability, not functional purity."

**Follow-up Probes:**
- How do you handle checked exceptions in stream lambdas?
- What's the team readability cost of deep stream chains?
- When does stream debugging tooling (IntelliJ Stream Trace) help?

---

#### Q12: How would you design a functional-style API for a domain service? Show the principles.

**Expected Answer:**

**Principles for functional API design:**

1. **Immutable inputs and outputs** — methods should not modify parameters
2. **Return new objects, don't mutate** — enable chaining and thread safety
3. **Use Optional instead of null** — make absence explicit in the type system
4. **Accept functional interfaces as parameters** — enable composition and customization
5. **Make operations chainable** — return `this` or new instances for builder-like flow

```java
// Functional-style Order Processing API:
public sealed interface OrderResult 
    permits OrderResult.Success, OrderResult.Failure {
    
    record Success(Order order, Receipt receipt) implements OrderResult {}
    record Failure(Order order, List<String> reasons) implements OrderResult {}
    
    // Pattern-match friendly (Java 21+):
    default <T> T fold(
        Function<Success, T> onSuccess,
        Function<Failure, T> onFailure
    ) {
        return switch (this) {
            case Success s -> onSuccess.apply(s);
            case Failure f -> onFailure.apply(f);
        };
    }
}

// Service API using functional composition:
public class OrderService {
    
    // Accept strategy as functional interface:
    public OrderResult process(
        Order order,
        UnaryOperator<Order> enrichment,        // Customizable enrichment step
        Predicate<Order> validationRule,         // Pluggable validation
        Function<Order, Receipt> billing         // Swappable billing strategy
    ) {
        return Optional.of(order)
            .map(enrichment)                     // Apply caller's enrichment
            .filter(validationRule)              // Apply caller's validation
            .map(validated -> {
                Receipt receipt = billing.apply(validated);
                return (OrderResult) new OrderResult.Success(validated, receipt);
            })
            .orElseGet(() -> new OrderResult.Failure(
                order, List.of("Validation failed")));
    }
}

// Usage — caller composes behavior:
OrderResult result = orderService.process(
    order,
    o -> o.withDiscount(loyaltyDiscount),        // Enrichment
    o -> o.total().compareTo(MIN) > 0,           // Validation
    billingService::charge                        // Billing strategy
);

String message = result.fold(
    success -> "Order confirmed: " + success.receipt().id(),
    failure -> "Rejected: " + String.join(", ", failure.reasons())
);
```

**Why this design works:**
- **Testable:** Each function can be unit tested independently
- **Composable:** Caller controls enrichment/validation/billing logic
- **Type-safe:** `sealed` + records enable exhaustive pattern matching
- **Immutable:** All data objects are records (immutable by default)

**Follow-up Probes:**
- How does this compare to the Strategy pattern in OOP?
- How would you handle multiple validation rules that each return different errors?
- What are the limitations of this approach for complex workflow orchestration?

---

#### Q13: Compare FP trade-offs in enterprise Java — when does FP help vs. hurt at scale?

**Expected Answer:**

```
FP Trade-off Matrix for Enterprise Java:

                        Helps                          Hurts
  ─────────────────────────────────────────────────────────────
  Readability      │ Simple transformations,         │ Deep nested chains,
                   │ data pipelines, mapping         │ complex business logic
  ─────────────────┼─────────────────────────────────┼─────────────────────
  Testability      │ Pure functions, immutability     │ Mocking functional
                   │ = easy unit tests               │ interfaces in legacy code
  ─────────────────┼─────────────────────────────────┼─────────────────────
  Performance      │ Lazy evaluation, parallel        │ Autoboxing, allocation,
                   │ streams for CPU-bound work       │ overhead in hot loops
  ─────────────────┼─────────────────────────────────┼─────────────────────
  Concurrency      │ Immutability eliminates races,   │ Parallel streams have
                   │ stateless lambdas are safe       │ hidden FJP dependencies
  ─────────────────┼─────────────────────────────────┼─────────────────────
  Team Adoption    │ Familiar pattern for data        │ Steep curve for monadic
                   │ processing                       │ patterns, custom collectors
  ─────────────────┼─────────────────────────────────┼─────────────────────
  Debugging        │ Named method references           │ Anonymous lambdas in
                   │ = readable stack traces           │ stack traces unreadable
  ─────────────────┼─────────────────────────────────┼─────────────────────
  Maintenance      │ Declarative code ages well        │ Overuse of Optional,
                   │ for simple cases                  │ streams for everything
```

**Architectural guidelines for staff engineers:**

1. **Data transformation layer** → Embrace FP (Streams, map/filter/reduce, immutable DTOs)
2. **Business logic / orchestration** → Mix FP and OOP (use loops for complex conditionals, pure functions for calculations)
3. **Infrastructure / framework code** → FP for composition (middleware chains, plugin systems, functional config)
4. **Performance-critical paths** → Benchmark, prefer primitives and loops where measured
5. **API design** → Accept `Function`/`Predicate` parameters for extensibility, but provide reasonable defaults

**Anti-patterns observed in enterprise FP adoption:**
- "Stream everything" — using streams where a simple if-else or loop is clearer
- "Optional everywhere" — Optional as method parameter, field, or collection element
- "Lambda spaghetti" — 10+ line lambdas instead of named methods
- "Parallel by default" — adding `.parallel()` hoping for speed without measurement
- "Collector hell" — deeply nested `Collectors.groupingBy(... mapping(... reducing(...)))` instead of a simple loop

**Follow-up Probes:**
- How do you establish FP coding guidelines for a team of 20 developers?
- When would you choose Kotlin over Java for functional-style enterprise code?
- How does Java's FP story compare to C#'s LINQ approach?

---

#### Q14: How would you handle a large dataset (100M+ records) functionally in Java without running out of memory?

**Expected Answer:**

**Principle: Streaming processing with constant memory footprint.**

```java
// Pattern 1: File-based streaming with try-with-resources
try (Stream<String> lines = Files.lines(Path.of("data.csv"), StandardCharsets.UTF_8)) {
    Map<String, DoubleSummaryStatistics> regionStats = lines
        .skip(1)                        // Skip header
        .map(Record::fromCsv)           // Parse lazily — one line at a time
        .filter(Record::isValid)        // Drop invalid immediately
        .collect(Collectors.groupingBy(
            Record::region,
            Collectors.summarizingDouble(Record::amount)
        ));
    // Only the Map + stats objects in memory, NOT all 100M records
}

// Pattern 2: Database cursor streaming (JPA/JDBC)
@Transactional(readOnly = true)
public void processAllOrders(Consumer<Order> processor) {
    try (Stream<Order> orders = orderRepository.streamAll()) {
        orders
            .filter(order -> order.getStatus() == PENDING)
            .peek(processor)
            .forEach(entityManager::detach);  // Prevent persistence context bloat
    }
}

// Pattern 3: Chunked parallel processing with bounded memory
try (Stream<String> lines = Files.lines(hugePath)) {
    AtomicLong count = new AtomicLong();
    lines
        .collect(Collectors.groupingBy(
            line -> count.getAndIncrement() / CHUNK_SIZE))  // Group into chunks
        .values()
        .parallelStream()
        .forEach(chunk -> {
            // Process chunk — bounded memory per chunk
            List<Record> records = chunk.stream().map(Record::fromCsv).toList();
            batchInsert(records);
        });
}

// Pattern 4: Custom Spliterator for database pagination
public class PaginatedSpliterator<T> extends Spliterators.AbstractSpliterator<T> {
    private final Function<Integer, List<T>> pageFetcher;
    private int currentPage = 0;
    private Iterator<T> currentBatch = Collections.emptyIterator();
    
    public PaginatedSpliterator(Function<Integer, List<T>> pageFetcher, int pageSize) {
        super(Long.MAX_VALUE, ORDERED | NONNULL);
        this.pageFetcher = pageFetcher;
    }
    
    @Override
    public boolean tryAdvance(Consumer<? super T> action) {
        if (!currentBatch.hasNext()) {
            List<T> page = pageFetcher.apply(currentPage++);
            if (page.isEmpty()) return false;
            currentBatch = page.iterator();
        }
        action.accept(currentBatch.next());
        return true;
    }
}

// Usage:
StreamSupport.stream(
    new PaginatedSpliterator<>(page -> orderDao.findPage(page, 1000), 1000), false)
    .filter(Order::isPending)
    .forEach(this::process);
```

**Key principles for large-scale functional processing:**
1. **Never collect everything** — use terminal operations that aggregate (`count`, `sum`, `reduce`, `forEach`)
2. **Close stream resources** — `Files.lines()`, JDBC streams MUST be closed
3. **Avoid stateful intermediate ops** — `sorted()`, `distinct()` buffer entire dataset
4. **Detach/evict entities** — JPA persistence context grows unbounded otherwise
5. **Use primitive specializations** — avoid 100M autoboxed wrapper objects
6. **Consider `mapMulti` over `flatMap`** — avoids intermediate Stream object creation

**Follow-up Probes:**
- How does Hibernate's `ScrollableResults` compare to JPA stream?
- When would you switch from Streams to a batch processing framework like Spring Batch?
- How would you implement backpressure-like behavior in plain Java Streams?

---

#### Q15: Design a functional middleware/interceptor chain in Java (like Express.js middleware or Spring interceptors).

**Expected Answer:**

```java
// Functional middleware pattern:
@FunctionalInterface
public interface Middleware<T, R> {
    R handle(T request, Function<T, R> next);
    
    // Compose middlewares into a chain:
    default Middleware<T, R> andThen(Middleware<T, R> after) {
        return (request, next) -> 
            this.handle(request, req -> after.handle(req, next));
    }
    
    // Convert chain to a simple Function:
    static <T, R> Function<T, R> chain(
        Function<T, R> finalHandler,
        List<Middleware<T, R>> middlewares
    ) {
        Function<T, R> chain = finalHandler;
        // Build from inside out (last middleware wraps closest to handler):
        for (int i = middlewares.size() - 1; i >= 0; i--) {
            final Middleware<T, R> mw = middlewares.get(i);
            final Function<T, R> current = chain;
            chain = request -> mw.handle(request, current);
        }
        return chain;
    }
}

// Define middlewares as lambdas:
Middleware<HttpRequest, HttpResponse> logging = (req, next) -> {
    long start = System.nanoTime();
    log.info("→ {} {}", req.method(), req.path());
    HttpResponse response = next.apply(req);
    log.info("← {} {}ms", response.status(), (System.nanoTime() - start) / 1_000_000);
    return response;
};

Middleware<HttpRequest, HttpResponse> auth = (req, next) -> {
    if (!isAuthenticated(req)) return HttpResponse.unauthorized();
    return next.apply(req);
};

Middleware<HttpRequest, HttpResponse> rateLimit = (req, next) -> {
    if (rateLimiter.isExceeded(req.clientIp())) return HttpResponse.tooManyRequests();
    return next.apply(req);
};

Middleware<HttpRequest, HttpResponse> errorHandler = (req, next) -> {
    try {
        return next.apply(req);
    } catch (Exception e) {
        log.error("Unhandled error", e);
        return HttpResponse.internalError(e.getMessage());
    }
};

// Build the chain:
Function<HttpRequest, HttpResponse> pipeline = Middleware.chain(
    controller::handleRequest,              // Final handler
    List.of(errorHandler, logging, auth, rateLimit)  // Middlewares (outside → inside)
);

// Execute:
HttpResponse response = pipeline.apply(incomingRequest);
```

**Execution flow:**
```
Request → errorHandler → logging → auth → rateLimit → controller
                                                          ↓
Response ← errorHandler ← logging ← auth ← rateLimit ← controller
```

**Why this is powerful:**
- Each middleware is a single `Middleware` lambda — independently testable
- Composition via `andThen` or `chain` — no class hierarchy needed
- Order is explicit and configurable at runtime
- Same pattern works for: HTTP middleware, message processors, validation chains, retry/circuit-breaker wrappers

**Follow-up Probes:**
- How does this compare to the Chain of Responsibility pattern?
- How would you add async support (returning `CompletableFuture<R>`)?
- In Spring, where do you see this pattern used internally?

---

> **End of Guide**
> 
> This guide covers the depth expected at **Staff / Senior Engineer** level (8+ years).
> Topics range from lambda bytecode internals to production performance debugging to architectural design patterns.
> All code examples target **Java 17+** with forward references to Java 21+ features.
>
> **Recommended study path:**
> 1. Sections 1-3: Foundation — understand lambda/interface mechanics
> 2. Section 4: Streams internals — know the Sink chain and Spliterator
> 3. Sections 5-6: Patterns — apply FP idiomatically
> 4. Sections 7-8: Production reality — concurrency pitfalls and performance
> 5. Section 9: Modern Java — records, virtual threads, reactive
> 6. Section 10: Interview practice — work through all 15 questions with a timer
