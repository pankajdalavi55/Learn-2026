# Functional Programming in Java — Problem → Solution Map & Practical Examples

> **Purpose:** Quick-reference guide mapping real development problems to the right FP feature.  
> **Use this when:** You're coding and think *"What FP tool solves this?"* — look up the problem category below.

---

## Table of Contents

1. [Master Problem → FP Feature Map](#1-master-problem--fp-feature-map)
2. [Decision Flowcharts — Which FP Feature to Use](#2-decision-flowcharts--which-fp-feature-to-use)
3. [Feature-by-Feature: Problems It Solves + Examples](#3-feature-by-feature-problems-it-solves--examples)
    - 3.1 Lambda Expressions
    - 3.2 Functional Interfaces (Predicate, Function, Consumer, Supplier)
    - 3.3 Method References
    - 3.4 Streams API
    - 3.5 Optional
    - 3.6 Function Composition & Chaining
    - 3.7 Collectors
    - 3.8 CompletableFuture
    - 3.9 Records & Sealed Classes
    - 3.10 Predicate Composition
4. [Problem Category Deep Dives with Examples](#4-problem-category-deep-dives-with-examples)
    - 4.1 Data Transformation Problems
    - 4.2 Filtering & Searching Problems
    - 4.3 Aggregation & Statistical Problems
    - 4.4 Grouping & Partitioning Problems
    - 4.5 Validation Problems
    - 4.6 Error Handling Problems
    - 4.7 Configuration & Strategy Problems
    - 4.8 Event & Callback Problems
    - 4.9 Concurrency & Async Problems
    - 4.10 Builder & Factory Problems
5. [Real-World Scenario Examples](#5-real-world-scenario-examples)
6. [Anti-Patterns — When NOT to Use FP](#6-anti-patterns--when-not-to-use-fp)
7. [Quick Reference Cheat Sheet](#7-quick-reference-cheat-sheet)

---

## 1. Master Problem → FP Feature Map

```
╔══════════════════════════════════════╦═══════════════════════════════════════════╗
║  PROBLEM CATEGORY                    ║  FP FEATURES TO USE                       ║
╠══════════════════════════════════════╬═══════════════════════════════════════════╣
║                                      ║                                           ║
║  Transform data (A → B)             ║  Stream.map(), Function<T,R>              ║
║                                      ║  Method references, Records               ║
║                                      ║                                           ║
║  Filter / Search data               ║  Stream.filter(), Predicate<T>            ║
║                                      ║  Predicate.and() / .or() / .negate()     ║
║                                      ║  findFirst(), findAny(), anyMatch()       ║
║                                      ║                                           ║
║  Aggregate / Summarize              ║  Stream.reduce(), Collectors              ║
║                                      ║  summaryStatistics(), counting()          ║
║                                      ║                                           ║
║  Group / Partition data             ║  Collectors.groupingBy()                  ║
║                                      ║  Collectors.partitioningBy()              ║
║                                      ║  Collectors.toMap()                       ║
║                                      ║                                           ║
║  Handle null / absent values        ║  Optional<T>, Optional.map()             ║
║                                      ║  .flatMap(), .orElse(), .or()            ║
║                                      ║                                           ║
║  Handle errors functionally         ║  Either<L,R>, Try<T>                     ║
║                                      ║  Optional, CompletableFuture              ║
║                                      ║  .exceptionally(), .handle()             ║
║                                      ║                                           ║
║  Validate input                     ║  Predicate<T> composition                ║
║                                      ║  Function chaining with Either            ║
║                                      ║                                           ║
║  Strategy / Behavior injection      ║  Function<T,R>, Consumer<T>              ║
║                                      ║  Supplier<T>, Comparator<T>              ║
║                                      ║  Map<Key, Function<>>                    ║
║                                      ║                                           ║
║  Event handling / Callbacks         ║  Consumer<T>, Runnable                   ║
║                                      ║  Lambda expressions                       ║
║                                      ║  BiConsumer<T,U>                          ║
║                                      ║                                           ║
║  Async / Parallel processing        ║  CompletableFuture, parallelStream()     ║
║                                      ║  thenApply(), thenCompose()              ║
║                                      ║  thenCombine(), allOf()                  ║
║                                      ║                                           ║
║  Build objects / configurations     ║  Supplier<T>, Builder pattern + lambdas  ║
║                                      ║  Function composition                     ║
║                                      ║  UnaryOperator<T> chaining               ║
║                                      ║                                           ║
║  Control flow / Conditional logic   ║  Optional.map()/filter()                 ║
║                                      ║  Pattern matching (sealed + switch)       ║
║                                      ║  Map<Condition, Action>                   ║
║                                      ║                                           ║
║  Lazy computation / Deferred work   ║  Supplier<T>, Stream (lazy)              ║
║                                      ║  Stream.generate(), Stream.iterate()     ║
║                                      ║                                           ║
║  Pipeline / Chain processing        ║  Function.andThen() / .compose()         ║
║                                      ║  Consumer.andThen()                       ║
║                                      ║  Stream pipeline                          ║
║                                      ║                                           ║
║  Flatten nested structures          ║  Stream.flatMap()                         ║
║                                      ║  Optional.flatMap()                       ║
║                                      ║  mapMulti() (Java 16+)                   ║
║                                      ║                                           ║
║  Immutable data modeling            ║  Records (Java 16+)                      ║
║                                      ║  List.of(), Map.of(), Set.of()           ║
║                                      ║  Sealed interfaces (Java 17+)            ║
║                                      ║                                           ║
║  Replace if/else chains             ║  Map<Key, Function<>>                    ║
║                                      ║  sealed interface + switch expression     ║
║                                      ║  Predicate composition                    ║
║                                      ║                                           ║
║  Sorting / Ordering                 ║  Comparator.comparing()                  ║
║                                      ║  .thenComparing(), .reversed()           ║
║                                      ║  Stream.sorted()                          ║
║                                      ║                                           ║
║  String processing                  ║  Stream + Collectors.joining()            ║
║                                      ║  String.chars(), Pattern.splitAsStream() ║
║                                      ║  mapToInt(), map()                       ║
║                                      ║                                           ║
║  Deduplication                      ║  Stream.distinct(), Collectors.toSet()   ║
║                                      ║  Collectors.toMap() with merge           ║
╚══════════════════════════════════════╩═══════════════════════════════════════════╝
```

---

## 2. Decision Flowcharts — Which FP Feature to Use

### Flowchart 1: "I have a collection and want to..."

```
START: I have a Collection/Array
  │
  ├─── Transform each element?
  │     └─► stream().map(function)
  │           └─ Nested collections? → flatMap()
  │           └─ To primitive type? → mapToInt() / mapToLong() / mapToDouble()
  │
  ├─── Filter elements?
  │     └─► stream().filter(predicate)
  │           └─ Multiple conditions? → predicate1.and(predicate2).or(predicate3)
  │           └─ Find one match? → filter().findFirst() → returns Optional
  │           └─ Check if any match? → anyMatch(predicate)
  │           └─ Check if ALL match? → allMatch(predicate)
  │
  ├─── Sort elements?
  │     └─► stream().sorted(Comparator.comparing(keyExtractor))
  │           └─ Multiple fields? → .thenComparing(field2)
  │           └─ Descending? → .reversed()
  │           └─ Nulls? → Comparator.nullsFirst() / nullsLast()
  │
  ├─── Group elements?
  │     └─► stream().collect(Collectors.groupingBy(classifier))
  │           └─ Two groups only? → partitioningBy(predicate)
  │           └─ Group + count? → groupingBy(key, counting())
  │           └─ Group + sum? → groupingBy(key, summingInt(f))
  │
  ├─── Aggregate to single value?
  │     ├─ Sum/average/min/max? → mapToInt().sum() / .average() / .max()
  │     ├─ Custom fold? → reduce(identity, accumulator)
  │     ├─ String join? → map(::toString).collect(Collectors.joining(", "))
  │     └─ Collect to new type? → collect(Collectors.toList/toSet/toMap)
  │
  ├─── Remove duplicates?
  │     └─► stream().distinct()
  │           └─ By specific field? → collect(toMap(keyFn, valueFn, mergeFunction))
  │
  ├─── Limit/Skip elements?
  │     └─► stream().skip(n).limit(m)
  │           └─ Take while condition true? → takeWhile(predicate) (Java 9+)
  │           └─ Drop while condition true? → dropWhile(predicate) (Java 9+)
  │
  └─── Process each element (side effect)?
        └─► forEach(consumer)
              └─ Need ordered processing? → forEachOrdered(consumer)
              └─ Need index? → IntStream.range(0, size).forEach(i -> ...)
```

### Flowchart 2: "I might have a null / absent value..."

```
START: A method might return null or no result
  │
  ├─── Returning from my method?
  │     └─► Return Optional<T> instead of null
  │           └─ Optional.of(value)        — when value is NEVER null
  │           └─ Optional.ofNullable(val)  — when value MIGHT be null
  │           └─ Optional.empty()          — when no result
  │
  ├─── Received an Optional, need to transform?
  │     └─► optional.map(function)
  │           └─ Function returns Optional? → flatMap() to avoid Optional<Optional<>>
  │           └─ Need to filter? → .filter(predicate)
  │
  ├─── Need a default value?
  │     ├─ Default is constant/cheap? → .orElse(defaultValue)
  │     ├─ Default is expensive to compute? → .orElseGet(() -> compute())
  │     └─ Should throw if absent? → .orElseThrow(() -> new NotFoundException())
  │
  ├─── Need to chain multiple Optional sources?
  │     └─► optional.or(() -> alternative1).or(() -> alternative2)  (Java 9+)
  │
  └─── Working with a Stream of Optionals?
        └─► stream.flatMap(Optional::stream)  (Java 9+)
              Filters out empties and unwraps presents in one step
```

### Flowchart 3: "I need to handle different behaviors..."

```
START: Different behavior depending on input/condition
  │
  ├─── A few known strategies?
  │     └─► Map<Key, Function<Input, Output>> strategies
  │         strategies.get(type).apply(input)
  │
  ├─── Processing pipeline with steps?
  │     └─► Function<A,B> step1 = ...;
  │         Function<B,C> step2 = ...;
  │         step1.andThen(step2).andThen(step3).apply(input)
  │
  ├─── Multiple validations?
  │     └─► Predicate<T> rule1 = ...; rule2 = ...; rule3 = ...;
  │         rule1.and(rule2).and(rule3).test(input)
  │
  ├─── Multiple side-effect actions?
  │     └─► Consumer<T> action1.andThen(action2).andThen(action3)
  │
  ├─── Branching on type (sealed hierarchy)?
  │     └─► switch (value) {
  │             case TypeA a -> handleA(a);
  │             case TypeB b -> handleB(b);
  │         }
  │
  └─── Lazy initialization / deferred computation?
        └─► Supplier<T> lazy = () -> expensiveComputation();
            lazy.get(); // computed only when needed
```

### Flowchart 4: "I need async / parallel processing..."

```
START: I need concurrent or async work
  │
  ├─── CPU-bound computation on large dataset (>10K elements)?
  │     └─► collection.parallelStream().map(cpuIntensiveWork).collect(...)
  │
  ├─── Multiple independent I/O calls?
  │     └─► CompletableFuture.allOf(
  │             cf1 = supplyAsync(() -> httpCall1()),
  │             cf2 = supplyAsync(() -> httpCall2()),
  │             cf3 = supplyAsync(() -> httpCall3())
  │         ).thenApply(v -> combine(cf1.join(), cf2.join(), cf3.join()))
  │
  ├─── Sequential async pipeline (A → B → C)?
  │     └─► cf.thenApply(transformSync)      // map
  │           .thenCompose(transformAsync)     // flatMap
  │           .thenApply(anotherTransform)
  │           .exceptionally(fallback)
  │
  ├─── First-to-complete from multiple sources?
  │     └─► CompletableFuture.anyOf(cf1, cf2, cf3)
  │
  └─── Need timeout on async?
        └─► cf.orTimeout(3, TimeUnit.SECONDS)                    // throws on timeout
            cf.completeOnTimeout(defaultVal, 3, TimeUnit.SECONDS) // fallback on timeout
```

---

## 3. Feature-by-Feature: Problems It Solves + Examples

### 3.1 Lambda Expressions

**Problems solved:** Eliminating boilerplate for single-method interfaces, enabling inline behavior.

```java
// PROBLEM: Verbose anonymous class for simple behavior
// ❌ Before (Java 7):
Collections.sort(employees, new Comparator<Employee>() {
    @Override
    public int compare(Employee a, Employee b) {
        return a.getName().compareTo(b.getName());
    }
});

// ✅ After (Lambda):
employees.sort((a, b) -> a.getName().compareTo(b.getName()));

// ✅ Even better (Method reference + Comparator API):
employees.sort(Comparator.comparing(Employee::getName));
```

```java
// PROBLEM: Runnable boilerplate for thread/executor
// ❌ Before:
executor.submit(new Runnable() {
    public void run() { sendEmail(user); }
});

// ✅ After:
executor.submit(() -> sendEmail(user));
```

```java
// PROBLEM: Event handlers in UI / frameworks
// ❌ Before:
button.setOnAction(new EventHandler<ActionEvent>() {
    public void handle(ActionEvent event) { processClick(event); }
});

// ✅ After:
button.setOnAction(event -> processClick(event));
// ✅ Or: button.setOnAction(this::processClick);
```

---

### 3.2 Functional Interfaces — Predicate, Function, Consumer, Supplier

```
When to use WHICH functional interface:

  ╔═══════════════════╦══════════════════════╦══════════════════════════════════╗
  ║ Interface          ║ Signature            ║ USE WHEN YOU NEED TO...          ║
  ╠═══════════════════╬══════════════════════╬══════════════════════════════════╣
  ║ Predicate<T>      ║ T → boolean          ║ Test / validate / filter         ║
  ║ Function<T,R>     ║ T → R                ║ Transform / convert / map        ║
  ║ Consumer<T>       ║ T → void             ║ Perform action / side effect     ║
  ║ Supplier<T>       ║ () → T               ║ Generate / provide / lazy init   ║
  ║ UnaryOperator<T>  ║ T → T (same type)    ║ Modify-like transform            ║
  ║ BiFunction<T,U,R> ║ (T,U) → R            ║ Combine two inputs               ║
  ║ BiPredicate<T,U>  ║ (T,U) → boolean      ║ Test with two inputs             ║
  ║ BiConsumer<T,U>   ║ (T,U) → void         ║ Action on two inputs             ║
  ║ Comparator<T>     ║ (T,T) → int          ║ Ordering / sorting               ║
  ║ Runnable          ║ () → void             ║ Execute action (no input/output) ║
  ║ Callable<T>       ║ () → T (throws)       ║ Async task with result           ║
  ╚═══════════════════╩══════════════════════╩══════════════════════════════════╝
```

```java
// PREDICATE — Validation rules
Predicate<String> isNotBlank = s -> s != null && !s.isBlank();
Predicate<String> isValidEmail = s -> s.matches("^[\\w.-]+@[\\w.-]+\\.[a-zA-Z]{2,}$");
Predicate<String> isValidInput = isNotBlank.and(isValidEmail);

if (isValidInput.test(email)) { /* proceed */ }
users.stream().filter(u -> isValidInput.test(u.getEmail())).toList();

// FUNCTION — Data transformation
Function<Employee, EmployeeDTO> toDTO = emp -> new EmployeeDTO(
    emp.getId(), emp.getName(), emp.getDepartment()
);
List<EmployeeDTO> dtos = employees.stream().map(toDTO).toList();

// CONSUMER — Side-effect actions
Consumer<Order> logOrder = order -> log.info("Processing: {}", order.getId());
Consumer<Order> validateOrder = order -> { if (!order.isValid()) throw new ValidationException(); };
Consumer<Order> saveOrder = orderRepository::save;

Consumer<Order> processOrder = logOrder.andThen(validateOrder).andThen(saveOrder);
orders.forEach(processOrder);

// SUPPLIER — Lazy initialization & factories
Supplier<Connection> connFactory = () -> DriverManager.getConnection(url);
Supplier<List<String>> lazyConfig = () -> loadConfigFromFile();  // Only loaded when .get() called
Supplier<String> idGenerator = () -> UUID.randomUUID().toString();
```

---

### 3.3 Method References

```
Four types and WHEN to use each:

  ╔═══════════════════════════════════╦══════════════════════════════════════════════╗
  ║ Type                               ║ When to Use                                  ║
  ╠═══════════════════════════════════╬══════════════════════════════════════════════╣
  ║ Static: Integer::parseInt          ║ Calling a utility/static method on each       ║
  ║ Bound:  System.out::println        ║ Calling method on a SPECIFIC object           ║
  ║ Unbound: String::toLowerCase       ║ Calling method ON each stream element         ║
  ║ Constructor: ArrayList::new        ║ Creating new instances                        ║
  ╚═══════════════════════════════════╩══════════════════════════════════════════════╝
```

```java
// Static method reference — parsing, converting
List<Integer> numbers = strings.stream().map(Integer::parseInt).toList();
List<Path> paths = filenames.stream().map(Path::of).toList();

// Bound instance method — using a specific object's method
Logger log = LoggerFactory.getLogger(MyClass.class);
errors.forEach(log::error);    // equivalent to: errors.forEach(e -> log.error(e))

ObjectMapper mapper = new ObjectMapper();
List<Order> orders = jsonStrings.stream().map(mapper::readValueAsOrder).toList();

// Unbound instance method — calling a method on each element
List<String> upperNames = names.stream().map(String::toUpperCase).toList();
List<String> trimmed = inputs.stream().map(String::trim).toList();
boolean allValid = orders.stream().allMatch(Order::isValid);

// Constructor reference — creating objects
List<Thread> threads = tasks.stream().map(Thread::new).toList();
String[] array = names.stream().toArray(String[]::new);
Map<String, List<String>> map = keys.stream()
    .collect(Collectors.toMap(Function.identity(), k -> new ArrayList<>()));
```

---

### 3.4 Streams API

**Problems Solved:**

| Problem | Stream Solution |
|---------|----------------|
| Transform every element | `.map(function)` |
| Transform + flatten nested | `.flatMap(x -> x.getItems().stream())` |
| Keep only matching | `.filter(predicate)` |
| Remove duplicates | `.distinct()` |
| Sort | `.sorted(comparator)` |
| Take first N | `.limit(n)` |
| Skip first N | `.skip(n)` |
| Find one element | `.findFirst()` / `.findAny()` |
| Check existence | `.anyMatch()` / `.allMatch()` / `.noneMatch()` |
| Count | `.count()` |
| Collect to List | `.toList()` (Java 16+) or `.collect(toList())` |
| Collect to Set | `.collect(toSet())` |
| Collect to Map | `.collect(toMap(keyFn, valueFn))` |
| Join strings | `.collect(Collectors.joining(", "))` |
| Sum/Avg/Min/Max | `.mapToInt(fn).sum()` / `.average()` / `.max()` |
| Custom reduction | `.reduce(identity, accumulator)` |
| Group by category | `.collect(groupingBy(classifier))` |
| Partition (true/false) | `.collect(partitioningBy(predicate))` |
| Process lazily from file | `Files.lines(path).map(...).filter(...)` |
| Generate sequence | `Stream.iterate(seed, hasNext, next)` |

```java
// Example: Complete order processing pipeline
record Order(String id, String customer, List<LineItem> items, Status status) {}
record LineItem(String product, int quantity, BigDecimal price) {}

// Revenue by customer for active orders
Map<String, BigDecimal> revenueByCustomer = orders.stream()
    .filter(o -> o.status() == Status.ACTIVE)
    .collect(Collectors.groupingBy(
        Order::customer,
        Collectors.reducing(
            BigDecimal.ZERO,
            o -> o.items().stream()
                .map(li -> li.price().multiply(BigDecimal.valueOf(li.quantity())))
                .reduce(BigDecimal.ZERO, BigDecimal::add),
            BigDecimal::add
        )
    ));

// Top 5 most ordered products across all orders
List<String> top5Products = orders.stream()
    .flatMap(o -> o.items().stream())
    .collect(Collectors.groupingBy(LineItem::product, Collectors.summingInt(LineItem::quantity)))
    .entrySet().stream()
    .sorted(Map.Entry.<String, Integer>comparingByValue().reversed())
    .limit(5)
    .map(Map.Entry::getKey)
    .toList();
```

---

### 3.5 Optional

**Problems Solved:**

| Problem | Optional Solution |
|---------|------------------|
| Method might return null | Return `Optional<T>` |
| Chain null-safe access | `.map(getter).map(getter2)` |
| Nested optional access | `.flatMap(x -> x.getOptionalField())` |
| Provide default | `.orElse(default)` / `.orElseGet(supplier)` |
| Throw if missing | `.orElseThrow(exceptionSupplier)` |
| Conditional processing | `.ifPresent(consumer)` / `.ifPresentOrElse(action, emptyAction)` |
| Fall back to alternative | `.or(() -> alternativeOptional)` |
| Filter valid values | `.filter(predicate)` |
| Integrate with Streams | `.stream()` (Java 9+) → turns Optional into 0-or-1 element stream |

```java
// PROBLEM: Deep null checking (the "null pyramid")
// ❌ Before:
String city = null;
if (user != null) {
    Address addr = user.getAddress();
    if (addr != null) {
        City c = addr.getCity();
        if (c != null) {
            city = c.getName();
        }
    }
}
if (city == null) city = "Unknown";

// ✅ After:
String city = Optional.ofNullable(user)
    .map(User::getAddress)
    .map(Address::getCity)
    .map(City::getName)
    .orElse("Unknown");
```

```java
// PROBLEM: Try cache first, then DB, then default
// ❌ Before:
User user = cache.get(userId);
if (user == null) {
    user = database.find(userId);
    if (user == null) {
        user = User.anonymous();
    }
}

// ✅ After:
User user = cache.findById(userId)             // Optional<User>
    .or(() -> database.findById(userId))        // Try DB if cache miss
    .orElseGet(() -> User.anonymous());          // Default if both miss
```

```java
// PROBLEM: Converting Stream<Optional<T>> to Stream<T>
List<Optional<User>> optionalUsers = ids.stream().map(userService::findById).toList();

// ✅ Unwrap and filter in one step:
List<User> users = ids.stream()
    .map(userService::findById)      // Stream<Optional<User>>
    .flatMap(Optional::stream)        // Stream<User> (empties removed)
    .toList();
```

---

### 3.6 Function Composition & Chaining

**Problems Solved:** Building reusable processing pipelines from small, testable pieces.

```java
// PROBLEM: Multi-step data transformation pipeline
// Each step is independently testable:

Function<String, String> sanitize     = s -> s.trim().toLowerCase();
Function<String, String> removeEmoji  = s -> s.replaceAll("[^\\p{ASCII}]", "");
Function<String, String> normalize    = s -> s.replaceAll("\\s+", " ");
Function<String, String> truncate     = s -> s.length() > 100 ? s.substring(0, 100) : s;

// Compose into a single reusable pipeline:
Function<String, String> cleanInput = sanitize
    .andThen(removeEmoji)
    .andThen(normalize)
    .andThen(truncate);

String clean = cleanInput.apply("  Hello 🌍  World!   Extra   Spaces  ");
// "hello world! extra spaces"

// Can add/remove/reorder steps without touching other code
// Each function is independently unit-testable
```

```java
// PROBLEM: Building configurable processing pipelines
// Different processing based on configuration:

public Function<Order, Order> buildPipeline(Config config) {
    Function<Order, Order> pipeline = Function.identity(); // start with no-op
    
    if (config.isValidationEnabled()) 
        pipeline = pipeline.andThen(this::validate);
    if (config.isEnrichmentEnabled()) 
        pipeline = pipeline.andThen(this::enrich);
    if (config.isPricingEnabled())
        pipeline = pipeline.andThen(this::applyPricing);
    if (config.isDiscountEnabled())
        pipeline = pipeline.andThen(order -> applyDiscount(order, config.getDiscountRate()));
    
    return pipeline;
}

// Usage:
Function<Order, Order> processor = buildPipeline(loadConfig());
Order result = processor.apply(incomingOrder);
```

```java
// PROBLEM: Decorator / middleware pattern (request processing, logging, timing)
@FunctionalInterface
interface Middleware<T, R> {
    R handle(T request, Function<T, R> next);
}

// Build middleware chain:
Function<HttpRequest, HttpResponse> handler = this::processRequest;

// Wrap with logging:
Function<HttpRequest, HttpResponse> withLogging = req -> {
    log.info("→ {} {}", req.method(), req.path());
    HttpResponse resp = handler.apply(req);
    log.info("← {} ({}ms)", resp.status(), resp.duration());
    return resp;
};

// Wrap with auth:
Function<HttpRequest, HttpResponse> withAuth = req -> {
    if (!isAuthenticated(req)) return HttpResponse.unauthorized();
    return withLogging.apply(req);
};

// Wrap with rate limiting:
Function<HttpRequest, HttpResponse> withRateLimit = req -> {
    if (rateLimiter.isExceeded(req.clientIp())) return HttpResponse.tooManyRequests();
    return withAuth.apply(req);
};
```

---

### 3.7 Collectors

**Problems Solved:** Complex aggregations, grouping, and output formatting.

```java
// PROBLEM: Group employees by department
Map<Department, List<Employee>> byDept = employees.stream()
    .collect(Collectors.groupingBy(Employee::getDepartment));

// PROBLEM: Group + count
Map<Department, Long> headcount = employees.stream()
    .collect(Collectors.groupingBy(Employee::getDepartment, Collectors.counting()));

// PROBLEM: Group + average salary
Map<Department, Double> avgSalary = employees.stream()
    .collect(Collectors.groupingBy(
        Employee::getDepartment,
        Collectors.averagingDouble(Employee::getSalary)));

// PROBLEM: Group + collect names as comma-separated string
Map<Department, String> namesByDept = employees.stream()
    .collect(Collectors.groupingBy(
        Employee::getDepartment,
        Collectors.mapping(Employee::getName, Collectors.joining(", "))));

// PROBLEM: Partition into two groups (pass/fail, active/inactive)
Map<Boolean, List<Employee>> partitioned = employees.stream()
    .collect(Collectors.partitioningBy(e -> e.getSalary() > 100000));
List<Employee> highEarners = partitioned.get(true);
List<Employee> others = partitioned.get(false);

// PROBLEM: Create lookup Map (id → entity)
Map<String, Employee> employeeById = employees.stream()
    .collect(Collectors.toMap(Employee::getId, Function.identity()));

// PROBLEM: Create Map with duplicate key handling
Map<String, Employee> latestByName = employees.stream()
    .collect(Collectors.toMap(
        Employee::getName,
        Function.identity(),
        (existing, replacement) -> replacement  // keep latest on duplicate
    ));

// PROBLEM: Multi-level grouping
Map<Department, Map<String, List<Employee>>> byDeptAndCity = employees.stream()
    .collect(Collectors.groupingBy(
        Employee::getDepartment,
        Collectors.groupingBy(Employee::getCity)));

// PROBLEM: Downstream transformation after grouping
Map<Department, Optional<Employee>> highestPaidByDept = employees.stream()
    .collect(Collectors.groupingBy(
        Employee::getDepartment,
        Collectors.maxBy(Comparator.comparing(Employee::getSalary))));

// PROBLEM: Collecting to unmodifiable collections
List<String> immutableNames = employees.stream()
    .map(Employee::getName)
    .collect(Collectors.toUnmodifiableList()); // Java 10+
```

---

### 3.8 CompletableFuture

**Problems Solved:** Async composition, parallel I/O, timeout handling.

```java
// PROBLEM: Multiple independent API calls in parallel
CompletableFuture<UserProfile> profileCF = CompletableFuture
    .supplyAsync(() -> userService.getProfile(userId));
CompletableFuture<List<Order>> ordersCF = CompletableFuture
    .supplyAsync(() -> orderService.getOrders(userId));
CompletableFuture<WalletBalance> walletCF = CompletableFuture
    .supplyAsync(() -> walletService.getBalance(userId));

// Wait for ALL, then combine:
DashboardData dashboard = CompletableFuture.allOf(profileCF, ordersCF, walletCF)
    .thenApply(v -> new DashboardData(profileCF.join(), ordersCF.join(), walletCF.join()))
    .orTimeout(5, TimeUnit.SECONDS)
    .join();
```

```java
// PROBLEM: Sequential async pipeline (fetch → transform → save)
CompletableFuture<Void> result = fetchOrder(orderId)          // CF<Order>
    .thenApply(this::validate)                                 // CF<ValidatedOrder>
    .thenCompose(this::enrichWithPricing)                      // CF<PricedOrder> (async)
    .thenApply(this::calculateTax)                             // CF<TaxedOrder>
    .thenAcceptAsync(this::saveToDatabase, dbExecutor)         // CF<Void>
    .exceptionally(ex -> { log.error("Failed", ex); return null; });
```

```java
// PROBLEM: First successful result from multiple sources (race pattern)
CompletableFuture<Price> bestPrice = CompletableFuture.anyOf(
    fetchPriceFrom(supplierA),
    fetchPriceFrom(supplierB),
    fetchPriceFrom(supplierC)
).thenApply(result -> (Price) result);
```

```java
// PROBLEM: Retry with fallback
public <T> CompletableFuture<T> withRetry(Supplier<CompletableFuture<T>> task, int maxRetries) {
    CompletableFuture<T> cf = task.get();
    for (int i = 0; i < maxRetries; i++) {
        cf = cf.thenApply(CompletableFuture::completedFuture)
               .exceptionally(ex -> task.get())
               .thenCompose(Function.identity());
    }
    return cf;
}
```

---

### 3.9 Records & Sealed Classes

**Problems Solved:** Immutable data modeling, exhaustive type handling.

```java
// PROBLEM: Boilerplate POJO for data transfer
// ❌ Before: 50+ lines for fields, constructor, getters, equals, hashCode, toString
// ✅ After:
record EmployeeDTO(String id, String name, String department, BigDecimal salary) {}

// PROBLEM: Representing a fixed set of outcomes / states (algebraic data types)
sealed interface ApiResponse<T> permits Success, ClientError, ServerError {}
record Success<T>(T data, int statusCode) implements ApiResponse<T> {}
record ClientError<T>(String message, int statusCode) implements ApiResponse<T> {}
record ServerError<T>(String message, Throwable cause) implements ApiResponse<T> {}

// Exhaustive handling:
String handle(ApiResponse<Order> response) {
    return switch (response) {
        case Success<Order> s    -> "Order: " + s.data().id();
        case ClientError<Order> e -> "Client error: " + e.message();
        case ServerError<Order> e -> "Server error: " + e.message();
        // No default needed — compiler checks exhaustiveness!
    };
}
```

```java
// PROBLEM: Value object with validation
record Email(String value) {
    public Email {  // Compact constructor
        if (!value.matches("^[\\w.-]+@[\\w.-]+\\.[a-zA-Z]{2,}$"))
            throw new IllegalArgumentException("Invalid email: " + value);
        value = value.toLowerCase().trim();
    }
}

record Money(BigDecimal amount, String currency) {
    public Money {
        if (amount.compareTo(BigDecimal.ZERO) < 0) 
            throw new IllegalArgumentException("Negative amount");
        Objects.requireNonNull(currency);
    }
    
    public Money add(Money other) {
        if (!this.currency.equals(other.currency)) throw new CurrencyMismatchException();
        return new Money(this.amount.add(other.amount), this.currency);
    }
}
```

---

### 3.10 Predicate Composition

**Problems Solved:** Building complex validation/filter rules from simple reusable parts.

```java
// PROBLEM: Complex filtering logic that changes per business context
// Define atomic rules:
Predicate<Employee> isActive     = e -> e.getStatus() == Status.ACTIVE;
Predicate<Employee> isSenior     = e -> e.getYearsExperience() >= 8;
Predicate<Employee> isHighPerf   = e -> e.getRating() >= 4.5;
Predicate<Employee> isEngineering = e -> e.getDept() == Dept.ENGINEERING;
Predicate<Employee> isRemote     = Employee::isRemote;

// Compose for different business rules:
Predicate<Employee> promotionEligible = isActive.and(isSenior).and(isHighPerf);
Predicate<Employee> layoffProtected   = isHighPerf.or(isSenior.and(isActive));
Predicate<Employee> remoteEngineers   = isActive.and(isEngineering).and(isRemote);
Predicate<Employee> needsImprovement  = isActive.and(isHighPerf.negate());

// Use in different contexts:
List<Employee> toPromote = employees.stream().filter(promotionEligible).toList();
List<Employee> toReview  = employees.stream().filter(needsImprovement).toList();
long remoteEngCount      = employees.stream().filter(remoteEngineers).count();
```

```java
// PROBLEM: Dynamic filter building from user search criteria
public Predicate<Product> buildSearchFilter(SearchCriteria criteria) {
    Predicate<Product> filter = product -> true; // accept all initially
    
    if (criteria.getCategory() != null)
        filter = filter.and(p -> p.getCategory().equals(criteria.getCategory()));
    if (criteria.getMinPrice() != null)
        filter = filter.and(p -> p.getPrice().compareTo(criteria.getMinPrice()) >= 0);
    if (criteria.getMaxPrice() != null)
        filter = filter.and(p -> p.getPrice().compareTo(criteria.getMaxPrice()) <= 0);
    if (criteria.isInStockOnly())
        filter = filter.and(Product::isInStock);
    if (criteria.getKeyword() != null)
        filter = filter.and(p -> p.getName().toLowerCase().contains(criteria.getKeyword().toLowerCase()));
    
    return filter;
}

// Usage:
List<Product> results = products.stream()
    .filter(buildSearchFilter(userCriteria))
    .sorted(Comparator.comparing(Product::getPrice))
    .limit(50)
    .toList();
```

---

## 4. Problem Category Deep Dives with Examples

### 4.1 Data Transformation Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Convert entity to DTO (most common FP use case)
// ─────────────────────────────────────────────────────────
record Employee(String id, String name, String email, Dept dept, BigDecimal salary) {}
record EmployeeDTO(String id, String name, String dept) {}

// Reusable mapper:
Function<Employee, EmployeeDTO> toDTO = e -> new EmployeeDTO(e.id(), e.name(), e.dept().name());

List<EmployeeDTO> dtos = employees.stream().map(toDTO).toList();

// ─────────────────────────────────────────────────────────
// PROBLEM: Transform + flatten nested collection
// ─────────────────────────────────────────────────────────
record Department(String name, List<Employee> employees) {}

// Get ALL employee names across ALL departments:
List<String> allNames = departments.stream()
    .flatMap(dept -> dept.employees().stream())
    .map(Employee::name)
    .distinct()
    .sorted()
    .toList();

// ─────────────────────────────────────────────────────────
// PROBLEM: Convert Map<K,V> to different Map<K2,V2>
// ─────────────────────────────────────────────────────────
Map<String, Integer> nameToAge = Map.of("Alice", 30, "Bob", 25);

Map<String, String> nameToAgeLabel = nameToAge.entrySet().stream()
    .collect(Collectors.toMap(
        Map.Entry::getKey, 
        e -> e.getValue() + " years old"
    ));

// ─────────────────────────────────────────────────────────
// PROBLEM: Pivot data — rows to columns
// ─────────────────────────────────────────────────────────
record Sale(String product, String month, int quantity) {}

Map<String, Map<String, Integer>> pivot = sales.stream()
    .collect(Collectors.groupingBy(
        Sale::product,
        Collectors.toMap(Sale::month, Sale::quantity, Integer::sum)
    ));
// { "Widget": {"Jan": 100, "Feb": 150}, "Gadget": {"Jan": 200, "Feb": 80} }

// ─────────────────────────────────────────────────────────
// PROBLEM: CSV line → Object
// ─────────────────────────────────────────────────────────
Function<String, Employee> csvToEmployee = line -> {
    String[] parts = line.split(",");
    return new Employee(parts[0], parts[1], parts[2], Dept.valueOf(parts[3]), new BigDecimal(parts[4]));
};

List<Employee> employees = Files.lines(Path.of("employees.csv"))
    .skip(1)  // skip header
    .map(csvToEmployee)
    .toList();
```

---

### 4.2 Filtering & Searching Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Find first match with default
// ─────────────────────────────────────────────────────────
Employee fallback = Employee.unknown();

Employee found = employees.stream()
    .filter(e -> e.email().equals("target@company.com"))
    .findFirst()
    .orElse(fallback);

// ─────────────────────────────────────────────────────────
// PROBLEM: Check if ANY item meets a condition
// ─────────────────────────────────────────────────────────
boolean hasOverdue = invoices.stream()
    .anyMatch(inv -> inv.getDueDate().isBefore(LocalDate.now()) && !inv.isPaid());

// ─────────────────────────────────────────────────────────
// PROBLEM: Filter by multiple dynamic conditions
// ─────────────────────────────────────────────────────────
List<Predicate<Transaction>> rules = List.of(
    t -> t.getAmount().compareTo(BigDecimal.ZERO) > 0,
    t -> t.getStatus() != Status.CANCELLED,
    t -> t.getDate().isAfter(startDate),
    t -> allowedCurrencies.contains(t.getCurrency())
);

// Combine all rules with AND:
Predicate<Transaction> allRules = rules.stream()
    .reduce(t -> true, Predicate::and);

List<Transaction> valid = transactions.stream().filter(allRules).toList();

// ─────────────────────────────────────────────────────────
// PROBLEM: Find max / min by a property
// ─────────────────────────────────────────────────────────
Optional<Employee> highestPaid = employees.stream()
    .max(Comparator.comparing(Employee::salary));

Optional<Order> oldestOrder = orders.stream()
    .min(Comparator.comparing(Order::getCreatedAt));

// ─────────────────────────────────────────────────────────
// PROBLEM: Get distinct values of a specific field
// ─────────────────────────────────────────────────────────
Set<String> uniqueDepartments = employees.stream()
    .map(Employee::dept)
    .map(Dept::name)
    .collect(Collectors.toSet());
```

---

### 4.3 Aggregation & Statistical Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Sum, average, min, max, count in one pass
// ─────────────────────────────────────────────────────────
IntSummaryStatistics stats = employees.stream()
    .mapToInt(e -> e.salary().intValue())
    .summaryStatistics();

System.out.println("Count: " + stats.getCount());
System.out.println("Sum:   " + stats.getSum());
System.out.println("Avg:   " + stats.getAverage());
System.out.println("Min:   " + stats.getMin());
System.out.println("Max:   " + stats.getMax());

// ─────────────────────────────────────────────────────────
// PROBLEM: Calculate total with BigDecimal (money)
// ─────────────────────────────────────────────────────────
BigDecimal totalRevenue = orders.stream()
    .map(Order::getTotal)
    .reduce(BigDecimal.ZERO, BigDecimal::add);

// ─────────────────────────────────────────────────────────
// PROBLEM: Weighted average
// ─────────────────────────────────────────────────────────
record Exam(String subject, double score, double weight) {}

double weightedAvg = exams.stream()
    .mapToDouble(e -> e.score() * e.weight())
    .sum() 
    / exams.stream().mapToDouble(Exam::weight).sum();

// ─────────────────────────────────────────────────────────
// PROBLEM: Running total / cumulative sum
// ─────────────────────────────────────────────────────────
List<Integer> values = List.of(1, 2, 3, 4, 5);
List<Integer> cumulative = new ArrayList<>();
values.stream().reduce(0, (running, val) -> {
    int sum = running + val;
    cumulative.add(sum);
    return sum;
});
// cumulative = [1, 3, 6, 10, 15]
// Note: This uses a side-effect — for pure FP, use scan (not in Java stdlib)

// ─────────────────────────────────────────────────────────
// PROBLEM: Frequency map (count occurrences of each item)
// ─────────────────────────────────────────────────────────
List<String> words = List.of("apple", "banana", "apple", "cherry", "banana", "apple");

Map<String, Long> frequency = words.stream()
    .collect(Collectors.groupingBy(Function.identity(), Collectors.counting()));
// {apple=3, banana=2, cherry=1}

// Top N most frequent:
List<Map.Entry<String, Long>> top3 = frequency.entrySet().stream()
    .sorted(Map.Entry.<String, Long>comparingByValue().reversed())
    .limit(3)
    .toList();
```

---

### 4.4 Grouping & Partitioning Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Group orders by status, then get total value per status
// ─────────────────────────────────────────────────────────
Map<OrderStatus, BigDecimal> valueByStatus = orders.stream()
    .collect(Collectors.groupingBy(
        Order::getStatus,
        Collectors.reducing(BigDecimal.ZERO, Order::getTotal, BigDecimal::add)
    ));

// ─────────────────────────────────────────────────────────
// PROBLEM: Group by range/bucket (age groups, price tiers)
// ─────────────────────────────────────────────────────────
Function<Employee, String> ageGroup = e -> {
    int age = e.getAge();
    if (age < 25) return "Junior (< 25)";
    if (age < 35) return "Mid (25-34)";
    if (age < 50) return "Senior (35-49)";
    return "Staff (50+)";
};

Map<String, List<Employee>> byAgeGroup = employees.stream()
    .collect(Collectors.groupingBy(ageGroup));

// ─────────────────────────────────────────────────────────
// PROBLEM: Group by composite key
// ─────────────────────────────────────────────────────────
record DeptCity(String dept, String city) {}

Map<DeptCity, List<Employee>> grouped = employees.stream()
    .collect(Collectors.groupingBy(e -> new DeptCity(e.getDept(), e.getCity())));

// ─────────────────────────────────────────────────────────
// PROBLEM: Partition with downstream collector
// ─────────────────────────────────────────────────────────
Map<Boolean, Long> passFailCount = students.stream()
    .collect(Collectors.partitioningBy(
        s -> s.getScore() >= 60,    // true = pass, false = fail
        Collectors.counting()
    ));
// {true=45, false=5}
```

---

### 4.5 Validation Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Validate an object with multiple rules, collect ALL errors
// ─────────────────────────────────────────────────────────
record ValidationRule<T>(Predicate<T> rule, String errorMessage) {}

List<ValidationRule<Order>> orderRules = List.of(
    new ValidationRule<>(o -> o.getItems() != null && !o.getItems().isEmpty(), "Order must have items"),
    new ValidationRule<>(o -> o.getTotal().compareTo(BigDecimal.ZERO) > 0, "Total must be positive"),
    new ValidationRule<>(o -> o.getCustomer() != null, "Customer is required"),
    new ValidationRule<>(o -> o.getShippingAddress() != null, "Shipping address is required"),
    new ValidationRule<>(o -> o.getItems().size() <= 100, "Max 100 items per order")
);

List<String> errors = orderRules.stream()
    .filter(rule -> !rule.rule().test(order))       // find FAILING rules
    .map(ValidationRule::errorMessage)               // extract messages
    .toList();

if (!errors.isEmpty()) {
    throw new ValidationException(errors);
}

// ─────────────────────────────────────────────────────────
// PROBLEM: Validate and transform (parse with error collection)
// ─────────────────────────────────────────────────────────
record ParseResult<T>(Optional<T> value, List<String> errors) {
    static <T> ParseResult<T> success(T value) { return new ParseResult<>(Optional.of(value), List.of()); }
    static <T> ParseResult<T> failure(String error) { return new ParseResult<>(Optional.empty(), List.of(error)); }
    boolean isSuccess() { return value.isPresent(); }
}

Function<String, ParseResult<Integer>> safeParseInt = s -> {
    try { return ParseResult.success(Integer.parseInt(s)); }
    catch (NumberFormatException e) { return ParseResult.failure("Not a number: " + s); }
};

List<ParseResult<Integer>> results = inputs.stream().map(safeParseInt).toList();
List<Integer> values = results.stream().filter(ParseResult::isSuccess).map(r -> r.value().get()).toList();
List<String> parseErrors = results.stream().flatMap(r -> r.errors().stream()).toList();
```

---

### 4.6 Error Handling Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Checked exception in lambda (Stream + IOException)
// ─────────────────────────────────────────────────────────

// Reusable wrapper for checked exceptions:
@FunctionalInterface
interface ThrowingFunction<T, R> {
    R apply(T t) throws Exception;
    
    static <T, R> Function<T, R> unchecked(ThrowingFunction<T, R> f) {
        return t -> {
            try { return f.apply(t); }
            catch (Exception e) { throw new RuntimeException(e); }
        };
    }
}

// Usage — clean stream with checked exceptions:
List<String> contents = paths.stream()
    .map(ThrowingFunction.unchecked(Files::readString))
    .toList();

// ─────────────────────────────────────────────────────────
// PROBLEM: Process items, skip failures, collect both results and errors
// ─────────────────────────────────────────────────────────
record Result<T>(T value, Throwable error) {
    boolean isSuccess() { return error == null; }
    static <T> Result<T> success(T value) { return new Result<>(value, null); }
    static <T> Result<T> failure(Throwable error) { return new Result<>(null, error); }
}

List<Result<ProcessedOrder>> results = rawOrders.stream()
    .map(order -> {
        try { return Result.success(processOrder(order)); }
        catch (Exception e) { return Result.<ProcessedOrder>failure(e); }
    })
    .toList();

List<ProcessedOrder> successes = results.stream()
    .filter(Result::isSuccess).map(Result::value).toList();
List<Throwable> failures = results.stream()
    .filter(r -> !r.isSuccess()).map(Result::error).toList();

log.info("Processed: {} success, {} failed", successes.size(), failures.size());
```

---

### 4.7 Configuration & Strategy Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Replace if/else or switch for strategy selection
// ─────────────────────────────────────────────────────────

// ❌ Before — growing switch/if chain:
BigDecimal calculateDiscount(Order order) {
    switch (order.getCustomerType()) {
        case REGULAR: return order.getTotal().multiply(new BigDecimal("0.05"));
        case PREMIUM: return order.getTotal().multiply(new BigDecimal("0.10"));
        case VIP:     return order.getTotal().multiply(new BigDecimal("0.20"));
        case EMPLOYEE: return order.getTotal().multiply(new BigDecimal("0.30"));
        default: return BigDecimal.ZERO;
    }
}

// ✅ After — Map + Function (open for extension, closed for modification):
private static final Map<CustomerType, UnaryOperator<BigDecimal>> DISCOUNT_STRATEGIES = Map.of(
    CustomerType.REGULAR,  total -> total.multiply(new BigDecimal("0.05")),
    CustomerType.PREMIUM,  total -> total.multiply(new BigDecimal("0.10")),
    CustomerType.VIP,      total -> total.multiply(new BigDecimal("0.20")),
    CustomerType.EMPLOYEE, total -> total.multiply(new BigDecimal("0.30"))
);

BigDecimal calculateDiscount(Order order) {
    return DISCOUNT_STRATEGIES
        .getOrDefault(order.getCustomerType(), total -> BigDecimal.ZERO)
        .apply(order.getTotal());
}
// Adding new customer type = ONE line in the map. No if/else changes.

// ─────────────────────────────────────────────────────────
// PROBLEM: Dynamic routing / dispatch based on message type
// ─────────────────────────────────────────────────────────
Map<String, Consumer<JsonNode>> handlers = Map.of(
    "ORDER_CREATED",   this::handleOrderCreated,
    "ORDER_CANCELLED", this::handleOrderCancelled,
    "PAYMENT_RECEIVED", this::handlePayment,
    "SHIPMENT_SENT",   this::handleShipment
);

public void dispatch(Event event) {
    handlers.getOrDefault(event.getType(), this::handleUnknown)
            .accept(event.getPayload());
}

// ─────────────────────────────────────────────────────────
// PROBLEM: Feature flags with functional behavior
// ─────────────────────────────────────────────────────────
Map<String, Supplier<Function<Order, Order>>> featureProcessors = Map.of(
    "NEW_PRICING",    () -> this::applyNewPricing,
    "LOYALTY_POINTS", () -> this::applyLoyaltyPoints,
    "FRAUD_CHECK",    () -> this::performFraudCheck
);

Function<Order, Order> pipeline = enabledFeatures.stream()
    .map(featureProcessors::get)
    .filter(Objects::nonNull)
    .map(Supplier::get)
    .reduce(Function.identity(), Function::andThen);

Order processed = pipeline.apply(order);
```

---

### 4.8 Event & Callback Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Observer pattern (event listeners)
// ─────────────────────────────────────────────────────────
class EventBus<T> {
    private final List<Consumer<T>> listeners = new CopyOnWriteArrayList<>();
    
    public void subscribe(Consumer<T> listener) { listeners.add(listener); }
    public void unsubscribe(Consumer<T> listener) { listeners.remove(listener); }
    public void publish(T event) { listeners.forEach(l -> l.accept(event)); }
}

// Usage:
EventBus<OrderEvent> orderEvents = new EventBus<>();
orderEvents.subscribe(event -> log.info("Order: {}", event));
orderEvents.subscribe(event -> metricsService.record(event));
orderEvents.subscribe(event -> notificationService.notify(event));

orderEvents.publish(new OrderEvent(orderId, "CREATED"));

// ─────────────────────────────────────────────────────────
// PROBLEM: Callback / hook pattern
// ─────────────────────────────────────────────────────────
class DataLoader<T> {
    private Consumer<T> onSuccess = t -> {};
    private Consumer<Throwable> onError = e -> {};
    private Runnable onComplete = () -> {};
    
    public DataLoader<T> onSuccess(Consumer<T> handler) { this.onSuccess = handler; return this; }
    public DataLoader<T> onError(Consumer<Throwable> handler) { this.onError = handler; return this; }
    public DataLoader<T> onComplete(Runnable handler) { this.onComplete = handler; return this; }
    
    public void load(Supplier<T> source) {
        try {
            T result = source.get();
            onSuccess.accept(result);
        } catch (Exception e) {
            onError.accept(e);
        } finally {
            onComplete.run();
        }
    }
}

// Usage:
new DataLoader<List<User>>()
    .onSuccess(users -> display(users))
    .onError(ex -> showAlert("Failed: " + ex.getMessage()))
    .onComplete(() -> hideSpinner())
    .load(() -> userService.fetchAll());
```

---

### 4.9 Concurrency & Async Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Process batch in parallel with limited concurrency
// ─────────────────────────────────────────────────────────
ExecutorService executor = Executors.newFixedThreadPool(10);

List<CompletableFuture<Result>> futures = items.stream()
    .map(item -> CompletableFuture.supplyAsync(() -> process(item), executor))
    .toList();

List<Result> results = futures.stream()
    .map(CompletableFuture::join)
    .toList();

executor.shutdown();

// ─────────────────────────────────────────────────────────
// PROBLEM: Timeout with fallback for each async call
// ─────────────────────────────────────────────────────────
Function<String, CompletableFuture<Price>> fetchWithFallback = productId ->
    CompletableFuture
        .supplyAsync(() -> pricingService.getPrice(productId))
        .completeOnTimeout(Price.defaultPrice(productId), 2, TimeUnit.SECONDS)
        .exceptionally(ex -> Price.defaultPrice(productId));

// ─────────────────────────────────────────────────────────
// PROBLEM: Collect all results, even if some fail
// ─────────────────────────────────────────────────────────
List<CompletableFuture<Result<Order>>> futures = orderIds.stream()
    .map(id -> CompletableFuture
        .supplyAsync(() -> orderService.fetch(id))
        .thenApply(Result::success)
        .exceptionally(Result::failure))
    .toList();

List<Result<Order>> allResults = futures.stream().map(CompletableFuture::join).toList();
long successCount = allResults.stream().filter(Result::isSuccess).count();
```

---

### 4.10 Builder & Factory Problems

```java
// ─────────────────────────────────────────────────────────
// PROBLEM: Builder pattern with functional configuration
// ─────────────────────────────────────────────────────────
record HttpRequest(String url, String method, Map<String, String> headers, String body) {}

class RequestBuilder {
    private String url;
    private String method = "GET";
    private Map<String, String> headers = new HashMap<>();
    private String body;
    
    public RequestBuilder url(String url) { this.url = url; return this; }
    public RequestBuilder method(String method) { this.method = method; return this; }
    public RequestBuilder header(String key, String value) { headers.put(key, value); return this; }
    public RequestBuilder body(String body) { this.body = body; return this; }
    
    public HttpRequest build() { return new HttpRequest(url, method, Map.copyOf(headers), body); }
    
    // Accept functional customization:
    public RequestBuilder customize(Consumer<RequestBuilder> customizer) {
        customizer.accept(this);
        return this;
    }
}

// Use with lambda customization:
HttpRequest request = new RequestBuilder()
    .url("https://api.example.com/orders")
    .method("POST")
    .customize(b -> {
        b.header("Authorization", "Bearer " + getToken());
        b.header("Content-Type", "application/json");
    })
    .body(jsonPayload)
    .build();

// ─────────────────────────────────────────────────────────
// PROBLEM: Object configuration via UnaryOperator chaining
// ─────────────────────────────────────────────────────────
record ServerConfig(int port, String host, boolean ssl, int maxConnections, Duration timeout) {}

static ServerConfig defaultConfig() {
    return new ServerConfig(8080, "localhost", false, 100, Duration.ofSeconds(30));
}

// Each customization is a UnaryOperator<ServerConfig> = ServerConfig → ServerConfig
UnaryOperator<ServerConfig> withSsl = c -> new ServerConfig(443, c.host(), true, c.maxConnections(), c.timeout());
UnaryOperator<ServerConfig> withHighCapacity = c -> new ServerConfig(c.port(), c.host(), c.ssl(), 1000, c.timeout());
UnaryOperator<ServerConfig> withProdHost = c -> new ServerConfig(c.port(), "api.prod.com", c.ssl(), c.maxConnections(), c.timeout());

// Compose:
ServerConfig prodConfig = Stream.of(withSsl, withHighCapacity, withProdHost)
    .reduce(UnaryOperator.identity(), (a, b) -> a.andThen(b)::apply)
    .apply(defaultConfig());
```

---

## 5. Real-World Scenario Examples

### Scenario 1: E-Commerce Order Report

```java
/*
 * BUSINESS REQUIREMENT:
 * Generate a monthly sales report showing:
 * - Revenue by product category
 * - Top 10 selling products
 * - Average order value by customer tier
 * - Orders with potential fraud signals
 */

record Order(String orderId, String customerId, CustomerTier tier, 
             LocalDate date, List<OrderItem> items, BigDecimal total) {}
record OrderItem(String productId, String productName, String category, 
                 int quantity, BigDecimal unitPrice) {}
enum CustomerTier { STANDARD, PREMIUM, VIP }

public SalesReport generateMonthlyReport(List<Order> orders, YearMonth month) {
    
    // Filter to target month:
    List<Order> monthlyOrders = orders.stream()
        .filter(o -> YearMonth.from(o.date()).equals(month))
        .toList();
    
    // 1. Revenue by category:
    Map<String, BigDecimal> revenueByCategory = monthlyOrders.stream()
        .flatMap(o -> o.items().stream())
        .collect(Collectors.groupingBy(
            OrderItem::category,
            Collectors.reducing(BigDecimal.ZERO,
                item -> item.unitPrice().multiply(BigDecimal.valueOf(item.quantity())),
                BigDecimal::add)));
    
    // 2. Top 10 selling products:
    List<ProductSales> top10 = monthlyOrders.stream()
        .flatMap(o -> o.items().stream())
        .collect(Collectors.groupingBy(
            OrderItem::productName,
            Collectors.summingInt(OrderItem::quantity)))
        .entrySet().stream()
        .sorted(Map.Entry.<String, Integer>comparingByValue().reversed())
        .limit(10)
        .map(e -> new ProductSales(e.getKey(), e.getValue()))
        .toList();
    
    // 3. Average order value by tier:
    Map<CustomerTier, Double> avgByTier = monthlyOrders.stream()
        .collect(Collectors.groupingBy(
            Order::tier,
            Collectors.averagingDouble(o -> o.total().doubleValue())));
    
    // 4. Fraud signals (multiple orders > $1000 in same day from same customer):
    List<String> suspiciousCustomers = monthlyOrders.stream()
        .filter(o -> o.total().compareTo(new BigDecimal("1000")) > 0)
        .collect(Collectors.groupingBy(
            o -> o.customerId() + "|" + o.date(),
            Collectors.counting()))
        .entrySet().stream()
        .filter(e -> e.getValue() > 2)    // More than 2 high-value orders in a day
        .map(e -> e.getKey().split("\\|")[0])
        .distinct()
        .toList();
    
    return new SalesReport(revenueByCategory, top10, avgByTier, suspiciousCustomers);
}
```

---

### Scenario 2: REST API Response Processing

```java
/*
 * BUSINESS REQUIREMENT:
 * Fetch user data from multiple microservices, combine, handle failures gracefully
 */

public CompletableFuture<UserDashboard> getUserDashboard(String userId) {
    
    // Fan-out: 4 independent service calls in parallel
    var profileCF = userService.getProfile(userId)
        .thenApply(Optional::of)
        .exceptionally(ex -> { log.warn("Profile fetch failed", ex); return Optional.empty(); });
    
    var ordersCF = orderService.getRecentOrders(userId, 10)
        .exceptionally(ex -> { log.warn("Orders fetch failed", ex); return List.of(); });
    
    var recommendationsCF = recommendationService.getForUser(userId)
        .completeOnTimeout(List.of(), 2, TimeUnit.SECONDS);  // Timeout → empty list
    
    var notificationsCF = notificationService.getUnread(userId)
        .exceptionally(ex -> List.of());
    
    // Fan-in: combine all results
    return CompletableFuture.allOf(profileCF, ordersCF, recommendationsCF, notificationsCF)
        .thenApply(v -> {
            UserProfile profile = profileCF.join().orElse(UserProfile.anonymous(userId));
            
            List<OrderSummary> orders = ordersCF.join().stream()
                .map(order -> new OrderSummary(
                    order.getId(),
                    order.getStatus().name(),
                    order.getTotal(),
                    order.getItems().stream().map(OrderItem::productName).toList()))
                .toList();
            
            return new UserDashboard(profile, orders, 
                recommendationsCF.join(), notificationsCF.join());
        })
        .orTimeout(5, TimeUnit.SECONDS);
}
```

---

### Scenario 3: Log File Analysis Pipeline

```java
/*
 * BUSINESS REQUIREMENT:
 * Parse access logs, find top endpoints, error rates, slow requests
 */

record LogEntry(Instant timestamp, String method, String path, 
                int statusCode, long responseTimeMs, String userAgent) {}

public LogAnalysis analyzeAccessLog(Path logFile) throws IOException {
    
    try (Stream<String> lines = Files.lines(logFile)) {
        List<LogEntry> entries = lines
            .filter(line -> !line.isBlank())
            .map(this::parseLogLine)
            .filter(Optional::isPresent)
            .map(Optional::get)
            .toList();
        
        // Top 10 most hit endpoints:
        List<Map.Entry<String, Long>> topEndpoints = entries.stream()
            .collect(Collectors.groupingBy(
                e -> e.method() + " " + e.path(), Collectors.counting()))
            .entrySet().stream()
            .sorted(Map.Entry.<String, Long>comparingByValue().reversed())
            .limit(10)
            .toList();
        
        // Error rate by endpoint:
        Map<String, Double> errorRateByEndpoint = entries.stream()
            .collect(Collectors.groupingBy(
                LogEntry::path,
                Collectors.collectingAndThen(
                    Collectors.partitioningBy(e -> e.statusCode() >= 400, Collectors.counting()),
                    partition -> {
                        long errors = partition.get(true);
                        long total = errors + partition.get(false);
                        return total > 0 ? (double) errors / total * 100 : 0.0;
                    })));
        
        // Slow requests (p95):
        long[] sortedTimes = entries.stream()
            .mapToLong(LogEntry::responseTimeMs)
            .sorted()
            .toArray();
        long p95 = sortedTimes[(int)(sortedTimes.length * 0.95)];
        
        List<LogEntry> slowRequests = entries.stream()
            .filter(e -> e.responseTimeMs() > p95)
            .sorted(Comparator.comparingLong(LogEntry::responseTimeMs).reversed())
            .limit(20)
            .toList();
        
        return new LogAnalysis(topEndpoints, errorRateByEndpoint, p95, slowRequests);
    }
}
```

---

### Scenario 4: Data Migration / ETL Pipeline

```java
/*
 * BUSINESS REQUIREMENT:
 * Migrate users from legacy system — validate, transform, deduplicate, batch insert
 */

public MigrationReport migrateUsers(Path legacyCsv) throws IOException {
    
    AtomicInteger lineNumber = new AtomicInteger(0);
    
    // Step 1: Parse + Validate (collect errors, don't stop)
    List<Result<NewUser>> results;
    try (Stream<String> lines = Files.lines(legacyCsv)) {
        results = lines
            .skip(1)  // header
            .map(line -> {
                int num = lineNumber.incrementAndGet();
                try {
                    LegacyUser legacy = parseLegacyCsv(line);
                    List<String> errors = validateLegacyUser(legacy);
                    if (!errors.isEmpty()) 
                        return Result.<NewUser>failure(new ValidationException("Line " + num + ": " + errors));
                    return Result.success(transformToNewUser(legacy));
                } catch (Exception e) {
                    return Result.<NewUser>failure(new ParseException("Line " + num + ": " + e.getMessage()));
                }
            })
            .toList();
    }
    
    // Step 2: Separate successes and failures
    Map<Boolean, List<Result<NewUser>>> partitioned = results.stream()
        .collect(Collectors.partitioningBy(Result::isSuccess));
    
    List<NewUser> validUsers = partitioned.get(true).stream()
        .map(Result::value)
        .toList();
    List<String> errors = partitioned.get(false).stream()
        .map(r -> r.error().getMessage())
        .toList();
    
    // Step 3: Deduplicate by email (keep latest)
    List<NewUser> deduplicated = validUsers.stream()
        .collect(Collectors.toMap(
            NewUser::email,
            Function.identity(),
            (existing, replacement) -> replacement))  // keep last
        .values().stream().toList();
    
    // Step 4: Batch insert (groups of 500)
    List<List<NewUser>> batches = IntStream.range(0, (deduplicated.size() + 499) / 500)
        .mapToObj(i -> deduplicated.subList(
            i * 500, Math.min((i + 1) * 500, deduplicated.size())))
        .toList();
    
    int inserted = batches.stream()
        .mapToInt(batch -> userRepository.batchInsert(batch))
        .sum();
    
    return new MigrationReport(results.size(), inserted, deduplicated.size() - validUsers.size(), errors);
}
```

---

### Scenario 5: Notification Rule Engine

```java
/*
 * BUSINESS REQUIREMENT:
 * Evaluate notification rules against events, send matching notifications.
 * Rules are configurable and composable.
 */

record NotificationRule(
    String name,
    Predicate<Event> condition,
    Function<Event, Notification> builder,
    Consumer<Notification> sender
) {}

class NotificationEngine {
    private final List<NotificationRule> rules;
    
    public NotificationEngine(List<NotificationRule> rules) {
        this.rules = List.copyOf(rules);
    }
    
    public List<String> process(Event event) {
        return rules.stream()
            .filter(rule -> rule.condition().test(event))
            .map(rule -> {
                Notification notification = rule.builder().apply(event);
                rule.sender().accept(notification);
                return rule.name();
            })
            .toList();  // Return names of triggered rules
    }
}

// Configuration — each rule is composed from FP building blocks:
List<NotificationRule> rules = List.of(
    new NotificationRule(
        "High-Value Order Alert",
        event -> event.type().equals("ORDER") 
              && event.getAmount() > 10000,
        event -> new Notification("SLACK", "High-value order: $" + event.getAmount()),
        slackService::send
    ),
    new NotificationRule(
        "Error Spike Alert",
        event -> event.type().equals("ERROR") 
              && errorRateTracker.isAboveThreshold(),
        event -> new Notification("PAGERDUTY", "Error spike detected: " + event.getMessage()),
        pagerDutyService::alert
    ),
    new NotificationRule(
        "New User Welcome",
        event -> event.type().equals("USER_REGISTERED"),
        event -> new Notification("EMAIL", buildWelcomeEmail(event)),
        emailService::send
    )
);

NotificationEngine engine = new NotificationEngine(rules);
List<String> triggered = engine.process(incomingEvent);
```

---

## 6. Anti-Patterns — When NOT to Use FP

```
╔══════════════════════════════════════╦═══════════════════════════════════════╦══════════════════════════════╗
║ SITUATION                            ║ WHY NOT FP                            ║ BETTER APPROACH              ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Simple loop (<5 lines)              ║ Stream adds overhead + verbosity     ║ Plain for-loop               ║
║                                      ║ for trivial operations               ║                              ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Complex control flow (break,        ║ Streams can't break/continue/return  ║ Traditional loop with        ║
║ continue, early return)             ║ from enclosing method                ║ explicit control             ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Mutating objects in-place           ║ FP encourages immutability;          ║ for-loop or                  ║
║ (performance-critical batch update) ║ creating copies adds GC pressure     ║ forEach with mutation        ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Index-dependent logic               ║ Streams don't provide element index  ║ IntStream.range(0, n) or    ║
║ (i-th element, neighbors)           ║ easily                               ║ traditional for-loop         ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Exception-heavy operations          ║ Checked exceptions can't escape      ║ for-loop with try/catch      ║
║ (each item may throw differently)   ║ lambda cleanly                       ║ or ThrowingFunction wrapper  ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Very small collections (1-3 items)  ║ Stream setup overhead exceeds        ║ Direct access: list.get(0)   ║
║                                      ║ any benefit                          ║                              ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Algorithm requires multiple passes  ║ Streams are one-shot;                ║ Collect to list first, then  ║
║ over same data                      ║ can't reuse a consumed stream        ║ iterate multiple times       ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ Stream chain > 7-8 operations       ║ Becomes hard to read and debug       ║ Break into named variables   ║
║                                      ║                                      ║ or extract methods           ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ parallelStream() for I/O-bound      ║ Blocks ForkJoinPool common threads   ║ CompletableFuture +          ║
║ operations (HTTP, DB calls)         ║ → starves entire JVM                 ║ custom executor              ║
╠══════════════════════════════════════╬═══════════════════════════════════════╬══════════════════════════════╣
║ parallelStream() on small data      ║ Fork/join overhead > benefit         ║ Sequential stream or loop    ║
║ or fast operations (<10K elements)  ║                                      ║                              ║
╚══════════════════════════════════════╩═══════════════════════════════════════╩══════════════════════════════╝
```

```java
// ❌ DON'T: Stream for simple iteration with side-effects
list.stream().forEach(item -> save(item));
// ✅ DO: Direct forEach (no stream overhead)
list.forEach(this::save);

// ❌ DON'T: Stream chain with 10+ operations (unreadable)
result = data.stream().filter(a).map(b).flatMap(c).filter(d).map(e)
    .sorted(f).distinct().limit(g).map(h).collect(i);
// ✅ DO: Break into named steps
var filtered = data.stream().filter(a).filter(d).toList();
var transformed = filtered.stream().map(b).flatMap(c).toList();
var result = transformed.stream().map(e).sorted(f).distinct().limit(g).map(h).collect(i);

// ❌ DON'T: Optional for conditional logic (replacing if/else)
Optional.ofNullable(input).map(this::transform).orElse(fallback);
// ✅ DO: Simple ternary when you already know it might be null
input != null ? transform(input) : fallback;

// ❌ DON'T: Stream.of(single).map(f).findFirst().get()
// ✅ DO: f.apply(single)
```

---

## 7. Quick Reference Cheat Sheet

### Stream Operations At a Glance

```
Source              → Intermediate (lazy)     → Terminal (triggers)
─────────────────    ────────────────────       ──────────────────
.stream()            .filter(Predicate)         .collect(Collector)
.parallelStream()    .map(Function)             .toList()
Stream.of(...)       .flatMap(→ Stream)         .forEach(Consumer)
Stream.iterate()     .mapToInt/Long/Double()    .reduce(identity, op)
Stream.generate()    .sorted(Comparator)        .count()
Arrays.stream()      .distinct()                .findFirst() / findAny()
Files.lines()        .limit(n) / .skip(n)       .anyMatch() / .allMatch()
IntStream.range()    .peek(Consumer)            .noneMatch()
map.entrySet()       .takeWhile(Predicate)      .min() / .max()
  .stream()          .dropWhile(Predicate)      .toArray()
                     .mapMulti()                .sum() / .average()
                     .gather(Gatherer)          .summaryStatistics()
```

### Collector Recipes

```
Collector                                        │  Result Type
─────────────────────────────────────────────────┼──────────────────────
toList()                                         │  List<T>
toSet()                                          │  Set<T>
toUnmodifiableList()                             │  List<T> (immutable)
toMap(keyFn, valueFn)                            │  Map<K,V>
toMap(keyFn, valueFn, mergeFunction)             │  Map<K,V>
joining(", ")                                    │  String
groupingBy(classifier)                           │  Map<K, List<T>>
groupingBy(classifier, downstream)               │  Map<K, D>
partitioningBy(predicate)                        │  Map<Boolean, List<T>>
counting()                                       │  Long
summingInt(fn)                                   │  Integer
averagingDouble(fn)                              │  Double
maxBy(comparator)                                │  Optional<T>
minBy(comparator)                                │  Optional<T>
mapping(fn, downstream)                          │  applies fn then collects
filtering(predicate, downstream)                 │  filters then collects
collectingAndThen(collector, finisher)           │  applies finisher to result
reducing(identity, op)                           │  T
teeing(collector1, collector2, merger)           │  R  (Java 12+)
```

### Comparator Recipes

```java
// Single field:
Comparator.comparing(Employee::getName)

// Multiple fields:
Comparator.comparing(Employee::getDepartment)
          .thenComparing(Employee::getName)
          .thenComparing(Employee::getSalary)

// Descending:
Comparator.comparing(Employee::getSalary).reversed()

// Null-safe:
Comparator.comparing(Employee::getName, Comparator.nullsLast(Comparator.naturalOrder()))

// Custom extraction:
Comparator.comparingInt(s -> s.length())          // avoid autoboxing
Comparator.comparingDouble(Employee::getSalary)    // primitive comparator
```

### Function Composition Recipes

```java
// Sequential pipeline:     f.andThen(g).andThen(h).apply(x)     → h(g(f(x)))
// Reverse composition:     f.compose(g).apply(x)                 → f(g(x))

// Predicate logic:         p1.and(p2).or(p3).negate()
// Consumer sequence:       c1.andThen(c2).andThen(c3)

// Dynamic pipeline from list:
Function<T,T> pipeline = steps.stream()
    .reduce(Function.identity(), Function::andThen);

// Dynamic predicate from list:
Predicate<T> combined = rules.stream()
    .reduce(t -> true, Predicate::and);
```

---

> **End of Guide**
> 
> Use this as a lookup reference during development:
> 1. **Section 1** — Find the right FP feature for your problem type
> 2. **Section 2** — Follow the decision flowchart
> 3. **Section 3** — See feature-specific examples
> 4. **Section 4** — Deep-dive by problem category
> 5. **Section 5** — Full real-world scenarios
> 6. **Section 6** — Know when NOT to use FP
> 7. **Section 7** — Quick copy-paste recipes
