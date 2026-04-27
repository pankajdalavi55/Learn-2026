# Modern Java Features: Java 8 to 21 — Comprehensive Guide

> For senior engineers and FAANG-level interview preparation. Covers every significant language and API change from Java 8 through Java 21, with production-grade code examples and interview tips.

[← Previous: Memory Model & JVM Internals](07-Java-Memory-Model-and-JVM-Internals.md) | [Home](README.md) | [Next: Design Patterns →](09-Java-Design-Patterns-Guide.md)

---

## Table of Contents

1. [Java Version History Overview](#1-java-version-history-overview)
2. [Java 8 (LTS - March 2014)](#2-java-8-lts---march-2014)
3. [Java 9 (September 2017)](#3-java-9-september-2017)
4. [Java 10 (March 2018)](#4-java-10-march-2018)
5. [Java 11 (LTS - September 2018)](#5-java-11-lts---september-2018)
6. [Java 12-13](#6-java-12-13)
7. [Java 14](#7-java-14)
8. [Java 15](#8-java-15)
9. [Java 16](#9-java-16)
10. [Java 17 (LTS - September 2021)](#10-java-17-lts---september-2021)
11. [Java 18-20 (Non-LTS Highlights)](#11-java-18-20-non-lts-highlights)
12. [Java 21 (LTS - September 2023)](#12-java-21-lts---september-2023)
13. [Migration Guide](#13-migration-guide)
14. [Interview-Focused Summary](#14-interview-focused-summary)

---

## 1. Java Version History Overview

| Version | Release     | LTS? | Key Features (One-Liner)                                      |
|---------|-------------|------|---------------------------------------------------------------|
| 8       | Mar 2014    | Yes  | Lambdas, Streams, Optional, `java.time`, default methods      |
| 9       | Sep 2017    | No   | JPMS (modules), JShell, factory collection methods             |
| 10      | Mar 2018    | No   | `var` (local variable type inference)                          |
| 11      | Sep 2018    | Yes  | HTTP Client API, new String methods, single-file execution     |
| 12      | Mar 2019    | No   | Switch expressions (preview), `Collectors.teeing()`            |
| 13      | Sep 2019    | No   | Text blocks (preview), `yield` in switch                       |
| 14      | Mar 2020    | No   | Records (preview), `instanceof` pattern matching (preview)     |
| 15      | Sep 2020    | No   | Sealed classes (preview), text blocks (standard)               |
| 16      | Mar 2021    | No   | Records (standard), `Stream.toList()`                          |
| 17      | Sep 2021    | Yes  | Sealed classes (standard), pattern matching switch (preview)   |
| 18      | Mar 2022    | No   | UTF-8 default, simple web server, `@snippet` in Javadoc        |
| 19      | Sep 2022    | No   | Virtual threads (preview), record patterns (preview)           |
| 20      | Mar 2023    | No   | Scoped values (incubator), virtual threads second preview      |
| 21      | Sep 2023    | Yes  | Virtual threads (standard), sequenced collections, record patterns (standard) |

> **Interview Tip:** Memorize the four LTS versions — 8, 11, 17, 21 — and three headline features for each.

---

## 2. Java 8 (LTS - March 2014)

Java 8 was the most transformative release in Java's history, introducing functional programming constructs and a modern date/time API.

### 2.1 Lambda Expressions & Stream API

Lambdas and Streams are covered extensively in the **Streams & Functional Programming Guide**. Brief recap:

```java
List<String> names = List.of("Alice", "Bob", "Charlie", "Diana");

List<String> filtered = names.stream()
    .filter(n -> n.length() > 3)
    .map(String::toUpperCase)
    .sorted()
    .collect(Collectors.toList());
// [ALICE, CHARLIE, DIANA]
```

### 2.2 Optional\<T\>

`Optional` eliminates explicit null checks and communicates intent — "this value may be absent."

```java
public Optional<User> findUserByEmail(String email) {
    User user = userRepository.query(email);
    return Optional.ofNullable(user);
}

// Usage
String displayName = findUserByEmail("alice@example.com")
    .map(User::getDisplayName)
    .orElse("Anonymous");

// Throw if absent
User user = findUserByEmail("bob@example.com")
    .orElseThrow(() -> new UserNotFoundException("bob@example.com"));
```

**Rules of thumb:**
- Never use `Optional` for fields or method parameters — it's for return types.
- Never call `get()` without `isPresent()` — prefer `orElse()`, `orElseGet()`, `map()`.
- `Optional.of(null)` throws `NullPointerException` — use `ofNullable()` when uncertain.

### 2.3 java.time API

The `java.time` package replaced the deeply flawed `java.util.Date` and `Calendar`. All types are **immutable** and **thread-safe**.

```java
// LocalDate — date without time or timezone
LocalDate today = LocalDate.now();
LocalDate birthday = LocalDate.of(1990, Month.JUNE, 15);
LocalDate parsed = LocalDate.parse("2024-03-15");

long age = ChronoUnit.YEARS.between(birthday, today);

// LocalTime — time without date or timezone
LocalTime now = LocalTime.now();
LocalTime meeting = LocalTime.of(14, 30);
LocalTime afterMeeting = meeting.plusHours(1).plusMinutes(30);

// LocalDateTime — date + time, no timezone
LocalDateTime appointmentAt = LocalDateTime.of(2024, 6, 15, 14, 30);
LocalDateTime twoWeeksLater = appointmentAt.plusWeeks(2);

// ZonedDateTime — full date-time with timezone
ZonedDateTime nyTime = ZonedDateTime.now(ZoneId.of("America/New_York"));
ZonedDateTime tokyoTime = nyTime.withZoneSameInstant(ZoneId.of("Asia/Tokyo"));

// Instant — machine timestamp (epoch-based)
Instant timestamp = Instant.now();
Instant later = timestamp.plus(Duration.ofHours(2));

// Duration (time-based) vs Period (date-based)
Duration twoHours = Duration.ofHours(2);
Duration between = Duration.between(LocalTime.of(9, 0), LocalTime.of(17, 30)); // 8h30m

Period sixMonths = Period.ofMonths(6);
Period agePeriod = Period.between(birthday, today);

// DateTimeFormatter
DateTimeFormatter fmt = DateTimeFormatter.ofPattern("dd-MMM-yyyy HH:mm");
String formatted = LocalDateTime.now().format(fmt);       // "15-Jun-2024 14:30"
LocalDateTime back = LocalDateTime.parse("15-Jun-2024 14:30", fmt);
```

**Why `java.time` replaced `Date`/`Calendar`:**

| Problem with `Date`/`Calendar`        | Solution in `java.time`             |
|---------------------------------------|-------------------------------------|
| Mutable, not thread-safe              | All types are immutable             |
| Months are 0-indexed (Jan = 0)        | `Month.JANUARY` or 1               |
| No separation of date/time/timezone   | Distinct types for each concern     |
| Formatting tied to `SimpleDateFormat` | Thread-safe `DateTimeFormatter`     |
| No duration/period concept            | `Duration` and `Period` classes     |

### 2.4 Default and Static Methods in Interfaces

```java
@FunctionalInterface
public interface PaymentProcessor {
    boolean process(Payment payment); // single abstract method

    default boolean processWithRetry(Payment payment, int maxRetries) {
        for (int attempt = 1; attempt <= maxRetries; attempt++) {
            if (process(payment)) return true;
            log("Retry " + attempt + " for payment " + payment.getId());
        }
        return false;
    }

    static PaymentProcessor noOp() {
        return payment -> true;
    }
}
```

**Diamond Problem Resolution:** If two interfaces provide the same default method, the implementing class **must** override it:

```java
interface A { default String greet() { return "Hello from A"; } }
interface B { default String greet() { return "Hello from B"; } }

class C implements A, B {
    @Override
    public String greet() {
        return A.super.greet(); // explicit choice
    }
}
```

### 2.5 StringJoiner and String.join()

```java
StringJoiner sj = new StringJoiner(", ", "[", "]");
sj.add("alpha").add("beta").add("gamma");
System.out.println(sj); // [alpha, beta, gamma]

String csv = String.join(",", "one", "two", "three"); // "one,two,three"
```

### 2.6 Nashorn JavaScript Engine

```java
ScriptEngine engine = new ScriptEngineManager().getEngineByName("nashorn");
engine.eval("print('Hello from JS')");
```

> **Note:** Nashorn was deprecated in Java 11 and removed in Java 15. Use GraalVM's JavaScript engine as a replacement.

### 2.7 CompletableFuture

Covered in the **Concurrency & Multithreading Guide**. It enables non-blocking asynchronous pipelines with `thenApply()`, `thenCompose()`, `thenCombine()`, and more.

---

## 3. Java 9 (September 2017)

### 3.1 Java Platform Module System (JPMS)

The module system provides **strong encapsulation** and **reliable configuration** — packages are only accessible if explicitly exported, and dependencies are declared up front.

```java
// module-info.java for com.app.order module
module com.app.order {
    requires java.sql;                        // compile + runtime dependency
    requires transitive com.app.common;       // consumers also get this
    exports com.app.order.api;                // public API package
    exports com.app.order.model to com.app.web; // qualified export
    opens com.app.order.internal to com.google.gson; // reflection access
    provides com.app.order.spi.OrderValidator
        with com.app.order.internal.DefaultOrderValidator;
    uses com.app.order.spi.PaymentGateway;    // service consumer
}
```

**Key directives:**

| Directive            | Purpose                                                      |
|----------------------|--------------------------------------------------------------|
| `requires`           | Declare a dependency on another module                       |
| `requires transitive`| Dependency is also visible to modules that depend on this one|
| `exports`            | Make a package accessible to other modules                   |
| `opens`              | Allow reflective access (for frameworks like Spring, Gson)   |
| `provides ... with`  | Register a service implementation                            |
| `uses`               | Declare that this module consumes a service                  |

**Classpath vs Modulepath:** Code on the classpath lives in the **unnamed module** (no encapsulation). Code on the modulepath gets full module encapsulation. Libraries can be used on either — the module system is opt-in.

### 3.2 JShell (REPL)

```text
$ jshell
jshell> var list = List.of(1, 2, 3)
list ==> [1, 2, 3]

jshell> list.stream().mapToInt(Integer::intValue).sum()
$2 ==> 6

jshell> /exit
```

### 3.3 Private Interface Methods

```java
public interface DataExporter {
    default void exportAsJson(List<Record> records) {
        String data = serialize(records);
        writeToFile(data, "export.json");
    }

    default void exportAsCsv(List<Record> records) {
        String data = serialize(records);
        writeToFile(data, "export.csv");
    }

    private String serialize(List<Record> records) {
        return records.stream()
            .map(Record::toString)
            .collect(Collectors.joining("\n"));
    }

    private void writeToFile(String data, String filename) { /* ... */ }
}
```

### 3.4 Immutable Collection Factory Methods

```java
List<String> colors = List.of("red", "green", "blue");     // immutable
Set<Integer> primes = Set.of(2, 3, 5, 7, 11);              // immutable, no dupes
Map<String, Integer> scores = Map.of("Alice", 95, "Bob", 87); // up to 10 entries

// For more than 10 entries:
Map<String, Integer> large = Map.ofEntries(
    Map.entry("Alice", 95),
    Map.entry("Bob", 87),
    Map.entry("Charlie", 92)
);
```

**Constraints:** No `null` keys or values. `Set.of()` / `Map.of()` throw on duplicates. Returned collections are truly unmodifiable — `add()` / `put()` throw `UnsupportedOperationException`.

### 3.5 Stream Additions

```java
// takeWhile — takes elements while predicate is true (ordered streams)
List<Integer> taken = Stream.of(2, 4, 6, 7, 8, 10)
    .takeWhile(n -> n % 2 == 0)
    .toList(); // [2, 4, 6]

// dropWhile — drops elements while predicate is true
List<Integer> dropped = Stream.of(2, 4, 6, 7, 8, 10)
    .dropWhile(n -> n % 2 == 0)
    .toList(); // [7, 8, 10]

// ofNullable — 0 or 1 element stream
Stream<String> s = Stream.ofNullable(getNullableValue()); // empty if null

// iterate with predicate (finite stream without limit())
Stream.iterate(1, n -> n <= 100, n -> n * 2)
    .forEach(System.out::println); // 1, 2, 4, 8, 16, 32, 64
```

### 3.6 Optional Additions

```java
Optional<String> opt = findConfigValue("db.url");

// ifPresentOrElse
opt.ifPresentOrElse(
    url -> connectTo(url),
    () -> log.warn("No DB URL configured, using default")
);

// or — lazy fallback Optional
String dbUrl = findConfigValue("db.url")
    .or(() -> findConfigValue("db.fallback.url"))
    .orElse("jdbc:h2:mem:default");

// stream — integrate Optional into Stream pipelines
List<String> allUrls = configKeys.stream()
    .map(this::findConfigValue)     // Stream<Optional<String>>
    .flatMap(Optional::stream)      // Stream<String>, empties removed
    .toList();
```

### 3.7 ProcessHandle API

```java
ProcessHandle current = ProcessHandle.current();
System.out.println("PID: " + current.pid());
current.info().command().ifPresent(cmd -> System.out.println("Command: " + cmd));
current.info().totalCpuDuration().ifPresent(d -> System.out.println("CPU: " + d));

// List all child processes
current.children().forEach(ph ->
    System.out.println("Child PID: " + ph.pid()));
```

### 3.8 Multi-Release JAR Files & try-with-resources Enhancement

```java
// Java 9: effectively final variables can be used in try-with-resources
Connection conn = DriverManager.getConnection(url);
try (conn) { // no need for re-declaration
    // use conn
}
```

---

## 4. Java 10 (March 2018)

### 4.1 Local Variable Type Inference (`var`)

The `var` keyword lets the compiler infer the type from the initializer. It is **not** dynamic typing — the type is fixed at compile time.

```java
// Good uses of var
var users = new ArrayList<User>();               // ArrayList<User>
var stream = users.stream();                     // Stream<User>
var entry = Map.entry("key", 42);                // Map.Entry<String, Integer>
var response = httpClient.send(request, ofString()); // HttpResponse<String>

// Especially helpful with complex generics
var grouped = users.stream()
    .collect(Collectors.groupingBy(
        User::getDepartment,
        Collectors.mapping(User::getName, Collectors.toList())
    )); // Map<String, List<String>> — inferred
```

**Where `var` CANNOT be used:**

```java
// Fields
// var count = 0;              // COMPILE ERROR

// Method parameters
// void process(var item) {}   // COMPILE ERROR

// Return types
// var getUser() { ... }       // COMPILE ERROR

// Uninitialized variables
// var x;                      // COMPILE ERROR

// null initializer
// var nothing = null;         // COMPILE ERROR

// Array initializer without new
// var arr = {1, 2, 3};        // COMPILE ERROR
```

**Style guidelines:**
- Use `var` when the type is obvious from the right-hand side: `var list = new ArrayList<String>()`
- Avoid `var` when the type is not clear: `var result = compute()` — what type is `result`?
- Prefer descriptive variable names when using `var`: `var userCount` over `var x`

### 4.2 Unmodifiable Collection Copies

```java
List<String> original = new ArrayList<>(List.of("a", "b", "c"));
List<String> copy = List.copyOf(original);  // unmodifiable snapshot

original.add("d");
System.out.println(copy); // [a, b, c] — not affected

// In stream pipelines
List<String> unmodifiable = original.stream()
    .filter(s -> s.length() > 1)
    .collect(Collectors.toUnmodifiableList());
```

### 4.3 G1 GC Improvements

Java 10 introduced **parallel full GC for G1**, reducing worst-case pause times. G1 became the default GC in Java 9 — this made its fallback full collection significantly faster.

---

## 5. Java 11 (LTS - September 2018)

### 5.1 HTTP Client API (`java.net.http`)

Replaces the legacy `HttpURLConnection`. Supports HTTP/1.1, HTTP/2, WebSocket, and async operations.

```java
HttpClient client = HttpClient.newBuilder()
    .version(HttpClient.Version.HTTP_2)
    .connectTimeout(Duration.ofSeconds(10))
    .followRedirects(HttpClient.Redirect.NORMAL)
    .build();

// Synchronous GET
HttpRequest getRequest = HttpRequest.newBuilder()
    .uri(URI.create("https://api.example.com/users/42"))
    .header("Accept", "application/json")
    .GET()
    .build();

HttpResponse<String> response = client.send(getRequest, HttpResponse.BodyHandlers.ofString());
System.out.println(response.statusCode());  // 200
System.out.println(response.body());        // {"id": 42, "name": "Alice"}

// Asynchronous POST
String jsonBody = """
    {"name": "Bob", "email": "bob@example.com"}
    """;

HttpRequest postRequest = HttpRequest.newBuilder()
    .uri(URI.create("https://api.example.com/users"))
    .header("Content-Type", "application/json")
    .POST(HttpRequest.BodyPublishers.ofString(jsonBody))
    .build();

client.sendAsync(postRequest, HttpResponse.BodyHandlers.ofString())
    .thenApply(HttpResponse::body)
    .thenAccept(body -> System.out.println("Created: " + body))
    .join();
```

### 5.2 New String Methods

```java
"   ".isBlank();           // true (whitespace-only or empty)
"  hello  ".strip();       // "hello"  (Unicode-aware, unlike trim())
"  hello  ".stripLeading();// "hello  "
"  hello  ".stripTrailing();// "  hello"
"ha".repeat(3);            // "hahaha"

"line1\nline2\nline3".lines()  // Stream<String> of lines
    .filter(l -> !l.isBlank())
    .forEach(System.out::println);
```

### 5.3 File Utility Methods

```java
Path path = Path.of("/tmp/config.txt");

// One-liner read/write
String content = Files.readString(path);
Files.writeString(path, "key=value\n", StandardOpenOption.APPEND);
```

### 5.4 Optional.isEmpty()

```java
Optional<String> opt = Optional.empty();
if (opt.isEmpty()) { // clearer than !opt.isPresent()
    System.out.println("No value present");
}
```

### 5.5 var in Lambda Parameters

Enables adding annotations to lambda parameters:

```java
list.sort((@NotNull var a, @NotNull var b) -> a.compareToIgnoreCase(b));
```

### 5.6 Single-File Source Code Execution

```text
$ java HelloWorld.java
```

No need for `javac` first — the JVM compiles and runs in one step. Useful for scripts and quick experiments.

### 5.7 Removed Java EE and CORBA Modules

`java.xml.ws`, `java.xml.bind` (JAXB), `java.activation`, `java.corba` were removed. Add these as Maven/Gradle dependencies if needed:

```xml
<dependency>
    <groupId>jakarta.xml.bind</groupId>
    <artifactId>jakarta.xml.bind-api</artifactId>
    <version>4.0.0</version>
</dependency>
```

---

## 6. Java 12-13

### 6.1 Switch Expressions (Preview in 12, Enhanced in 13)

Traditional switch is statement-based and error-prone (fall-through). Switch expressions produce a value and use arrow syntax.

```java
// Java 12: arrow syntax, no fall-through
String dayType = switch (day) {
    case MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY -> "Weekday";
    case SATURDAY, SUNDAY -> "Weekend";
};

// Java 13: yield for multi-line blocks
int numLetters = switch (day) {
    case MONDAY, FRIDAY, SUNDAY -> 6;
    case TUESDAY -> 7;
    case WEDNESDAY -> 9;
    case THURSDAY, SATURDAY -> 8;
    default -> {
        String trimmed = day.toString().trim();
        yield trimmed.length(); // yield returns value from block
    }
};
```

### 6.2 Text Blocks (Preview in 13)

```java
String json = """
        {
            "name": "Alice",
            "roles": ["admin", "user"],
            "active": true
        }
        """;

String html = """
        <html>
            <body>
                <h1>Welcome</h1>
            </body>
        </html>
        """;
```

**Whitespace rules:** The compiler strips **incidental** whitespace (the common leading indentation) but preserves **essential** whitespace (indentation relative to the closing `"""`).

### 6.3 String.indent() and String.transform()

```java
String indented = "hello\nworld".indent(4);
// "    hello\n    world\n"

String result = "  hello  "
    .transform(String::strip)
    .transform(String::toUpperCase);
// "HELLO"
```

### 6.4 Collectors.teeing()

Applies two downstream collectors and merges their results:

```java
record Stats(long count, double average) {}

Stats stats = List.of(10, 20, 30, 40, 50).stream()
    .collect(Collectors.teeing(
        Collectors.counting(),
        Collectors.averagingInt(Integer::intValue),
        Stats::new
    ));
// Stats[count=5, average=30.0]
```

---

## 7. Java 14

### 7.1 Records (Preview)

Records are transparent carriers for immutable data. The compiler generates the constructor, accessors, `equals()`, `hashCode()`, and `toString()`.

```java
public record Employee(String name, String department, BigDecimal salary) {

    // Compact canonical constructor — for validation
    public Employee {
        Objects.requireNonNull(name, "name must not be null");
        Objects.requireNonNull(department, "department must not be null");
        if (salary.compareTo(BigDecimal.ZERO) < 0) {
            throw new IllegalArgumentException("Salary cannot be negative");
        }
    }

    // Custom constructor must delegate to canonical
    public Employee(String name, String department) {
        this(name, department, BigDecimal.ZERO);
    }

    // Custom accessor (computed property)
    public String displayName() {
        return name + " (" + department + ")";
    }
}
```

**What records generate:**

| Generated Member     | Description                                       |
|----------------------|---------------------------------------------------|
| Constructor          | All-args canonical constructor                    |
| Accessors            | `name()`, `department()`, `salary()` (no `get`)   |
| `equals()`           | Component-wise equality                           |
| `hashCode()`         | Based on all components                           |
| `toString()`         | `Employee[name=Alice, department=Eng, salary=100]`|

**Restrictions:** Records are implicitly `final`, cannot extend other classes (implicitly extend `Record`), fields are `final`, no instance fields beyond components. Records **can** implement interfaces.

```java
public sealed interface Shape permits Circle, Rectangle {}
public record Circle(double radius) implements Shape {}
public record Rectangle(double width, double height) implements Shape {}
```

### 7.2 Helpful NullPointerExceptions

```text
// Before Java 14
Exception in thread "main" java.lang.NullPointerException

// Java 14+ with -XX:+ShowCodeDetailsInExceptionMessages (default from Java 17)
Exception in thread "main" java.lang.NullPointerException:
  Cannot invoke "String.length()" because "user.getAddress().getCity()" is null
```

Pinpoints exactly which part of a chained expression was null — invaluable for debugging.

### 7.3 Pattern Matching for instanceof (Preview)

```java
// Before: cast after check
if (obj instanceof String) {
    String s = (String) obj;
    System.out.println(s.toUpperCase());
}

// Java 14+: pattern variable in scope
if (obj instanceof String s) {
    System.out.println(s.toUpperCase());
}

// Works with logical operators (flow scoping)
if (obj instanceof String s && s.length() > 5) {
    process(s);
}
```

### 7.4 Switch Expressions (Standard)

Switch expressions became a permanent feature in Java 14 (no longer preview).

---

## 8. Java 15

### 8.1 Sealed Classes (Preview)

Sealed classes restrict which classes can extend them, enabling **exhaustive** pattern matching.

```java
public sealed interface PaymentMethod
    permits CreditCard, BankTransfer, DigitalWallet {}

public record CreditCard(String number, YearMonth expiry) implements PaymentMethod {}
public record BankTransfer(String iban, String bic) implements PaymentMethod {}
public non-sealed class DigitalWallet implements PaymentMethod {
    private final String provider;
    private final String accountId;
    // constructor, getters...
}
```

**Permitted subtypes must be one of:**
- `final` — no further extension
- `sealed` — restricts its own subtypes
- `non-sealed` — opens up extension (opt-out)

### 8.2 Text Blocks (Standard)

Text blocks became standard in Java 15, with two new escape sequences:

```java
String singleLine = """
        This is a very long line that we want to keep as a \
        single line in the resulting string.""";
// No newline where \ appears

String withTrailingSpaces = """
        Name:   Alice\s
        Role:   Admin\s
        """;
// \s is a space that prevents trailing whitespace stripping

String formatted = """
        Dear %s,
        Your order #%d has shipped.
        """.formatted("Alice", 12345);
```

### 8.3 Hidden Classes

Hidden classes cannot be discovered by other classes and are meant for framework use — JVM generates them at runtime for dynamic proxies, lambda metafactories, etc. Not typically used directly by application developers.

---

## 9. Java 16

### 9.1 Records and Pattern Matching for instanceof (Standard)

Both features graduated from preview to **permanent** standard features in Java 16.

### 9.2 Stream.toList()

```java
List<String> names = employees.stream()
    .map(Employee::name)
    .toList(); // returns unmodifiable List — more concise than collect(Collectors.toList())
```

> **Subtle difference:** `Collectors.toList()` returns a mutable `ArrayList`; `Stream.toList()` returns an unmodifiable list. Choose accordingly.

### 9.3 Stream.mapMulti()

An imperative alternative to `flatMap`, useful when the mapping produces zero or few elements:

```java
// Flatten Optional results without flatMap
List<String> activeNames = users.stream()
    .<String>mapMulti((user, consumer) -> {
        if (user.isActive()) {
            consumer.accept(user.getName());
        }
    })
    .toList();

// One-to-many mapping
List<Integer> flattened = Stream.of(1, 2, 3)
    .<Integer>mapMulti((num, consumer) -> {
        consumer.accept(num);
        consumer.accept(num * 10);
    })
    .toList(); // [1, 10, 2, 20, 3, 30]
```

### 9.4 Day Period Support in DateTimeFormatter

```java
DateTimeFormatter formatter = DateTimeFormatter.ofPattern("h:mm B");
System.out.println(LocalTime.of(8, 30).format(formatter));  // 8:30 in the morning
System.out.println(LocalTime.of(14, 30).format(formatter)); // 2:30 in the afternoon
System.out.println(LocalTime.of(22, 0).format(formatter));  // 10:00 at night
```

---

## 10. Java 17 (LTS - September 2021)

Java 17 is the most widely adopted LTS after Java 8/11. It finalizes sealed classes and introduces pattern matching for switch as a preview.

### 10.1 Sealed Classes (Standard, Final)

```java
public sealed interface Shape permits Circle, Rectangle, Triangle {}

public record Circle(double radius) implements Shape {
    public double area() { return Math.PI * radius * radius; }
}

public record Rectangle(double width, double height) implements Shape {
    public double area() { return width * height; }
}

public record Triangle(double base, double height) implements Shape {
    public double area() { return 0.5 * base * height; }
}

// Exhaustive switch — compiler knows all subtypes
public static double calculateArea(Shape shape) {
    return switch (shape) {
        case Circle c    -> c.area();
        case Rectangle r -> r.area();
        case Triangle t  -> t.area();
        // no default needed — sealed + all permits covered
    };
}
```

### 10.2 Pattern Matching for Switch (Preview)

```java
public static String describe(Object obj) {
    return switch (obj) {
        case Integer i when i > 0  -> "Positive integer: " + i;
        case Integer i             -> "Non-positive integer: " + i;
        case String s when s.isBlank() -> "Blank string";
        case String s              -> "String of length " + s.length();
        case int[] arr             -> "Int array of length " + arr.length;
        case null                  -> "null value";
        default                    -> "Unknown: " + obj.getClass().getSimpleName();
    };
}
```

### 10.3 Stronger Encapsulation of JDK Internals

Internal APIs like `sun.misc.Unsafe` are no longer accessible by default. Use `--add-opens` on the command line or migrate to supported APIs:

```text
java --add-opens java.base/sun.nio.ch=ALL-UNNAMED -jar app.jar
```

### 10.4 New Random Generator API

```java
RandomGenerator rng = RandomGeneratorFactory.of("L128X256MixRandom").create();

// Stream of random integers in range
rng.ints(10, 1, 100).forEach(System.out::println);

// Pick a random algorithm
RandomGeneratorFactory.all()
    .map(RandomGeneratorFactory::name)
    .sorted()
    .forEach(System.out::println);
```

### 10.5 Deprecations and Removals

- **Applet API** deprecated for removal
- **Security Manager** deprecated for removal
- **RMI Activation** removed
- Nashorn JavaScript engine removed (already deprecated in 11)

---

## 11. Java 18-20 (Non-LTS Highlights)

### 11.1 Java 18

**UTF-8 by Default:** `Charset.defaultCharset()` now returns UTF-8 on all platforms. No more platform-dependent encoding surprises.

**Simple Web Server:**

```text
$ jwebserver --port 8080 --directory /var/www
Binding to loopback by default. For all interfaces use "-b 0.0.0.0" or "-b ::".
Serving /var/www and target directories on 127.0.0.1 port 8080
```

Serves static files — useful for quick prototyping. Not for production.

**Code Snippets in Javadoc:**

```java
/**
 * Calculates compound interest.
 *
 * {@snippet :
 *   double result = Finance.compoundInterest(1000.0, 0.05, 12);
 *   // result = 1795.86
 * }
 */
public static double compoundInterest(double principal, double rate, int years) {
    return principal * Math.pow(1 + rate, years);
}
```

### 11.2 Java 19-20: Record Patterns (Preview)

Record patterns allow deconstructing record values inside `instanceof` and `switch`:

```java
record Point(int x, int y) {}
record Line(Point start, Point end) {}

// Nested deconstruction
static String describeLine(Object obj) {
    if (obj instanceof Line(Point(var x1, var y1), Point(var x2, var y2))) {
        return "Line from (%d,%d) to (%d,%d)".formatted(x1, y1, x2, y2);
    }
    return "Not a line";
}
```

### 11.3 Java 19-20: Virtual Threads (Preview)

Virtual threads are lightweight threads managed by the JVM, not the OS. They enable a "thread per request" model without exhausting system threads.

```java
// Preview syntax — same concepts, finalized in Java 21
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    IntStream.range(0, 10_000).forEach(i ->
        executor.submit(() -> {
            Thread.sleep(Duration.ofSeconds(1));
            return i;
        })
    );
}
```

### 11.4 Structured Concurrency (Incubator)

```java
// Incubator API — manages a group of concurrent subtasks as a unit
try (var scope = new StructuredTaskScope.ShutdownOnFailure()) {
    Subtask<User> userTask  = scope.fork(() -> fetchUser(userId));
    Subtask<Order> orderTask = scope.fork(() -> fetchOrder(orderId));

    scope.join().throwIfFailed();

    return new UserOrder(userTask.get(), orderTask.get());
}
```

### 11.5 Foreign Function & Memory API (Preview)

Provides a safe, pure-Java alternative to JNI for calling native code:

```java
// Call strlen from C standard library
Linker linker = Linker.nativeLinker();
SymbolLookup stdlib = linker.defaultLookup();

MethodHandle strlen = linker.downcallHandle(
    stdlib.find("strlen").orElseThrow(),
    FunctionDescriptor.of(JAVA_LONG, ADDRESS)
);

try (Arena arena = Arena.ofConfined()) {
    MemorySegment cString = arena.allocateFrom("Hello");
    long len = (long) strlen.invoke(cString);
    System.out.println("Length: " + len); // 5
}
```

---

## 12. Java 21 (LTS - September 2023)

Java 21 is the latest LTS and brings several major features to standard status.

### 12.1 Virtual Threads (Standard, Final)

Virtual threads are cheap to create (a few hundred bytes of stack) and are scheduled by the JVM onto a small pool of platform (OS) threads.

```java
// Create a virtual thread directly
Thread vThread = Thread.ofVirtual()
    .name("worker-", 0)
    .start(() -> {
        System.out.println(Thread.currentThread());
        // Thread[#42,worker-0,5,VirtualThreads]
    });
vThread.join();

// Virtual thread per task executor — the recommended pattern
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    List<Future<String>> futures = new ArrayList<>();
    for (int i = 0; i < 100_000; i++) {
        final int taskId = i;
        futures.add(executor.submit(() -> {
            Thread.sleep(Duration.ofMillis(500));
            return "Result-" + taskId;
        }));
    }
    // 100K concurrent tasks — each on its own virtual thread
    for (var future : futures) {
        future.get(); // blocks the virtual thread, not the OS thread
    }
}
```

> **Interview Tip:** Virtual threads shine for I/O-bound workloads (HTTP handlers, DB queries). They don't help CPU-bound work. Avoid pinning: don't hold `synchronized` blocks during blocking calls — use `ReentrantLock` instead. See the **Concurrency Guide** for deep coverage.

### 12.2 Pattern Matching for Switch (Standard, Final)

```java
sealed interface Notification permits Email, Sms, Push {}
record Email(String to, String subject, String body) implements Notification {}
record Sms(String phone, String message) implements Notification {}
record Push(String deviceId, String title) implements Notification {}

public String formatNotification(Notification n) {
    return switch (n) {
        case Email e when e.body().length() > 1000 ->
            "Long email to %s: %s...".formatted(e.to(), e.subject());
        case Email e ->
            "Email to %s: %s".formatted(e.to(), e.subject());
        case Sms s ->
            "SMS to %s: %s".formatted(s.phone(), s.message());
        case Push p ->
            "Push to device %s: %s".formatted(p.deviceId(), p.title());
        // exhaustive — no default needed for sealed types
    };
}

// Null handling in switch
public String safeDescribe(Object obj) {
    return switch (obj) {
        case null         -> "null";
        case String s     -> "String: " + s;
        case Integer i    -> "Integer: " + i;
        default           -> "Other: " + obj;
    };
}
```

### 12.3 Record Patterns (Standard, Final)

```java
record Address(String city, String zip) {}
record Customer(String name, Address address) {}

// Nested deconstruction in switch
static String describeCustomer(Object obj) {
    return switch (obj) {
        case Customer(var name, Address(var city, var zip))
            when city.equals("Seattle") ->
                name + " is a local customer (ZIP: " + zip + ")";
        case Customer(var name, Address(var city, _)) ->
                name + " is from " + city;
        default -> "Not a customer";
    };
}
```

### 12.4 Sequenced Collections

Before Java 21, getting the first/last element or iterating in reverse required different code for `List`, `Deque`, `SortedSet`, and `LinkedHashMap`. Sequenced collections unify this.

**Interface hierarchy:**

```text
          SequencedCollection
           /              \
  SequencedSet       (List already had)
       |
  SequencedMap
```

```java
// SequencedCollection methods
List<String> list = new ArrayList<>(List.of("a", "b", "c"));
list.getFirst();   // "a"
list.getLast();     // "c"
list.addFirst("z"); // ["z", "a", "b", "c"]
list.addLast("d");  // ["z", "a", "b", "c", "d"]
list.reversed();    // reversed view: ["d", "c", "b", "a", "z"]

// SequencedMap methods
LinkedHashMap<String, Integer> map = new LinkedHashMap<>();
map.put("one", 1);
map.put("two", 2);
map.put("three", 3);

map.firstEntry();    // one=1
map.lastEntry();     // three=3
map.pollFirstEntry();// removes and returns one=1

SequencedMap<String, Integer> reversed = map.reversed();
reversed.forEach((k, v) -> System.out.println(k + "=" + v));
// three=3, two=2

// SequencedSet
LinkedHashSet<String> set = new LinkedHashSet<>(List.of("x", "y", "z"));
set.getFirst();  // "x"
set.getLast();   // "z"
set.reversed();  // [z, y, x]
```

**Full interface hierarchy (Mermaid):**

```text
                    Collection
                        |
              SequencedCollection
               /        |        \
          List      Deque    SequencedSet
                                  |
                              SortedSet
                                  |
                            NavigableSet

                       Map
                        |
                  SequencedMap
                        |
                    SortedMap
                        |
                  NavigableMap
```

### 12.5 String Templates (Preview)

```java
// STR template processor (preview)
String name = "Alice";
int age = 30;

String greeting = STR."Hello \{name}, you are \{age} years old.";
// "Hello Alice, you are 30 years old."

// Expressions in templates
String info = STR."Next year you'll be \{age + 1}. Name length: \{name.length()}.";

// Multi-line
String json = STR."""
    {
        "name": "\{name}",
        "age": \{age},
        "active": \{age < 65}
    }
    """;
```

> **Note:** String templates were removed from preview in later JDK releases due to design reconsideration. They may return in a different form.

### 12.6 Unnamed Patterns and Variables

Use `_` for intentionally unused variables — no compiler warnings:

```java
// Unnamed variable in enhanced for
for (var _ : collection) {
    totalIterations++;
}

// Unnamed pattern in switch
switch (shape) {
    case Circle c   -> handleCircle(c);
    case Rectangle _ -> handleGenericShape();  // don't need the variable
    case Triangle _  -> handleGenericShape();
}

// Unnamed variable in try-with-resources
try (var _ = ScopedContext.open()) {
    performWork();
}

// Unnamed in catch
try {
    riskyOperation();
} catch (IllegalArgumentException _) {
    log.warn("Invalid argument — using default");
}
```

### 12.7 Unnamed Classes and Instance Main Methods (Preview)

Simplifies the entry point for beginners — no class declaration needed:

```java
// Old way
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}

// Java 21 preview — unnamed class, instance main
void main() {
    System.out.println("Hello, World!");
}
```

---

## 13. Migration Guide

### 13.1 Java 8 → 17: Key Breaking Changes

| What Changed                          | Impact                                    | Action Needed                                           |
|---------------------------------------|-------------------------------------------|---------------------------------------------------------|
| Java EE modules removed (Java 11)    | `javax.xml.bind`, `javax.annotation` gone | Add Jakarta dependencies via Maven/Gradle               |
| `javax.*` → `jakarta.*` namespace    | Package names changed                     | Update imports when using Jakarta EE 9+                 |
| Nashorn removed (Java 15)            | JS engine no longer bundled               | Use GraalVM JS or standalone Nashorn                    |
| Strong encapsulation (Java 16+)      | Internal APIs inaccessible                | Use `--add-opens` or migrate to public APIs             |
| `SecurityManager` deprecated         | Will be removed in future                 | Migrate to modern security approaches                   |
| Applet API deprecated                | No browser plugin support                 | Use web technologies instead                            |
| Default GC changed to G1 (Java 9)    | Different performance characteristics     | Benchmark; tune G1 or switch to ZGC/Shenandoah          |
| String.split() behavior change       | Edge-case behavior fixed                  | Review regex-based splits for trailing empty strings     |
| `var` is a reserved type name        | Can't use `var` as a type name            | Rename classes/interfaces named `var`                   |
| Records, sealed classes              | New restricted keywords                   | Rename if using `record`, `sealed`, `permits` as idents |

**Module system impact:** Most applications can run on the classpath unchanged (unnamed module). Migrate to modules incrementally if desired. Libraries should provide `module-info.java` for proper modularity.

**Reflection restrictions:** Add `--add-opens` flags for frameworks that rely on deep reflection:

```text
--add-opens java.base/java.lang=ALL-UNNAMED
--add-opens java.base/java.util=ALL-UNNAMED
```

### 13.2 Java 17 → 21: Features to Adopt

| Feature                        | Status in 21 | Adoption Guidance                                     |
|--------------------------------|-------------|-------------------------------------------------------|
| Virtual threads                | Standard    | Use for I/O-bound workloads; replace thread pools      |
| Record patterns                | Standard    | Use in switch/instanceof for cleaner deconstruction    |
| Pattern matching for switch    | Standard    | Replace if-else chains with exhaustive switches        |
| Sequenced collections          | Standard    | Use `getFirst()`/`getLast()` instead of index tricks   |
| String templates               | Preview     | Wait for finalization before production use            |
| Structured concurrency         | Preview     | Experiment; not yet production-ready                   |

### 13.3 Migration Checklist

1. **Update build tools:** Maven 3.8+, Gradle 7.3+ for Java 17 support
2. **Check dependencies:** Ensure all libraries support target Java version
3. **Run with `--illegal-access=warn`** first (Java 16) to identify reflection issues
4. **Replace removed APIs:** JAXB, JAX-WS, CORBA — use Jakarta or third-party
5. **Test thoroughly:** Especially around date/time, string handling, serialization
6. **Update CI/CD:** Docker base images, GraalVM native-image configs
7. **Enable helpful NPE messages:** Default from Java 17
8. **Adopt new features gradually:** records → sealed → pattern matching

---

## 14. Interview-Focused Summary

### Quick-Reference Q&A

| # | Question | Key Answer |
|---|----------|------------|
| 1 | What are records in Java? | Immutable data carriers; compiler generates constructor, accessors, `equals`, `hashCode`, `toString`. Declared with `record Name(Type field) {}`. Cannot extend classes (implicitly extend `Record`), can implement interfaces. |
| 2 | Explain sealed classes. | Restrict which classes can extend/implement using `sealed ... permits`. Subtypes must be `final`, `sealed`, or `non-sealed`. Enables exhaustive pattern matching in switch. |
| 3 | What is `var` and where can it be used? | Local variable type inference. Usable only for local variables with initializers. Cannot be used for fields, method parameters, return types, or `null` initializers. Type is resolved at compile time — Java remains statically typed. |
| 4 | What are text blocks? | Multi-line string literals using `"""`. Support `\` (line continuation), `\s` (space), and `formatted()`. Incidental whitespace is stripped. Standard since Java 15. |
| 5 | Key differences between Java 8, 11, 17, 21? | **8:** Lambdas, streams, `java.time`. **11:** HTTP Client, new String methods, single-file exec. **17:** Sealed classes, stronger encapsulation, new random API. **21:** Virtual threads, sequenced collections, record patterns. |
| 6 | What is JPMS? | Java Platform Module System (Java 9). Uses `module-info.java` to declare dependencies (`requires`), public API (`exports`), and reflection access (`opens`). Provides strong encapsulation and reliable configuration. |
| 7 | What are virtual threads? | Lightweight threads managed by the JVM, not the OS. Created via `Thread.ofVirtual()` or `Executors.newVirtualThreadPerTaskExecutor()`. Ideal for I/O-bound tasks. Standard in Java 21. |
| 8 | What is pattern matching for instanceof? | `if (obj instanceof String s)` — combines type check and cast. The pattern variable `s` is in scope where the match is guaranteed. Standard since Java 16. |
| 9 | What are sequenced collections? | New interfaces (`SequencedCollection`, `SequencedSet`, `SequencedMap`) in Java 21. Provide `getFirst()`, `getLast()`, `addFirst()`, `addLast()`, `reversed()`. Unify encounter-order operations. |
| 10 | How does pattern matching for switch work? | Switch can match types: `case String s ->`. Supports guarded patterns: `case String s when s.length() > 5 ->`. Handles null: `case null ->`. Exhaustive for sealed types. Standard in Java 21. |
| 11 | What is `Optional` and how should it be used? | Container for a possibly-absent value. Use as return type, never for fields or parameters. Prefer `map()`, `orElse()`, `orElseThrow()` over `get()`. |
| 12 | What did `java.time` replace and why? | Replaced `java.util.Date` and `Calendar`. Old API was mutable, not thread-safe, had 0-indexed months. `java.time` is immutable, thread-safe, has clear types for date/time/zone/duration. |
| 13 | What are default methods in interfaces? | Methods with implementation in interfaces using `default` keyword. Added in Java 8 to evolve APIs without breaking implementors. Diamond problem resolved by requiring override. |
| 14 | Explain switch expressions vs switch statements. | Expressions return a value, use `->` syntax, no fall-through. Use `yield` for multi-line blocks. Standard since Java 14. |
| 15 | What is `Stream.toList()` and how does it differ from `collect(toList())`? | `Stream.toList()` (Java 16) returns an unmodifiable list. `Collectors.toList()` returns a mutable `ArrayList`. |
| 16 | What is `Collectors.teeing()`? | Java 12 collector that applies two downstream collectors to the same stream and merges results with a BiFunction. |
| 17 | What are record patterns? | Deconstruct records in `instanceof` and `switch`: `case Point(var x, var y)`. Support nesting. Standard in Java 21. |
| 18 | What is the HTTP Client API? | `java.net.http` package (Java 11). Supports HTTP/1.1, HTTP/2, async via `CompletableFuture`, builder pattern. Replaces `HttpURLConnection`. |
| 19 | What are the String methods added in Java 11? | `isBlank()`, `strip()`, `stripLeading()`, `stripTrailing()`, `repeat(int)`, `lines()`. `strip()` is Unicode-aware (unlike `trim()`). |
| 20 | What is structured concurrency? | Treats groups of concurrent tasks as a single unit of work. If one fails, others are cancelled. Uses `StructuredTaskScope`. Still preview/incubator as of Java 21. |
| 21 | What are unnamed variables (`_`)? | Java 21 feature. Use `_` for intentionally unused variables in `for`, `catch`, `try-with-resources`, and pattern matching. Eliminates "unused variable" warnings. |
| 22 | How do you migrate from Java 8 to 17? | Key steps: replace removed Java EE modules with Jakarta deps, add `--add-opens` for reflection, update build tools, rename any `var`/`record` identifiers, run `jdeps` for internal API usage. |
| 23 | What is the Foreign Function & Memory API? | Safe replacement for JNI. Call native code and manage off-heap memory from pure Java. Uses `Linker`, `SymbolLookup`, `MemorySegment`, `Arena`. Preview through Java 21. |
| 24 | What is `mapMulti()` in streams? | Java 16 imperative alternative to `flatMap()`. Uses a `BiConsumer<T, Consumer<R>>` — you push zero or more elements to the consumer. More efficient for small expansions. |

### Feature Maturity Cheat Sheet

| Feature                     | Preview | Standard |
|-----------------------------|---------|----------|
| Switch expressions          | 12      | 14       |
| Text blocks                 | 13      | 15       |
| Records                     | 14      | 16       |
| Pattern matching instanceof | 14      | 16       |
| Sealed classes              | 15      | 17       |
| Pattern matching switch     | 17      | 21       |
| Record patterns             | 19      | 21       |
| Virtual threads             | 19      | 21       |
| Structured concurrency      | 19      | —        |
| String templates            | 21      | —        |

> **Final Interview Tip:** When discussing Java versions, focus on the **four LTS releases** (8, 11, 17, 21). Know the headline feature of each, understand the progression from preview to standard for major features, and be prepared to discuss trade-offs (e.g., when to use records vs classes, virtual threads vs platform threads, `var` vs explicit types). Demonstrating awareness of the migration path from 8 → 17 → 21 signals production experience.

---

*End of Modern Java Features Guide (Java 8–21)*

---

[← Previous: Memory Model & JVM Internals](07-Java-Memory-Model-and-JVM-Internals.md) | [Home](README.md) | [Next: Design Patterns →](09-Java-Design-Patterns-Guide.md)
