# Java 8 to 17 - Modern Features Guide

## Why Learn Modern Java Features?

Modern Java features make code:
- ✅ **Shorter** - Less boilerplate
- ✅ **Cleaner** - More readable
- ✅ **Faster** - Better performance
- ✅ **Safer** - Fewer null pointer exceptions

**For Interviews:** Companies expect you to know Java 8+ features!

---

## Java 8 (Released 2014) - The Game Changer

### 1. Lambda Expressions - Anonymous Functions

**Before Java 8 (Verbose):**
```java
List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5);

// Old way - anonymous class
Collections.sort(numbers, new Comparator<Integer>() {
    @Override
    public int compare(Integer a, Integer b) {
        return a - b;
    }
});
```

**After Java 8 (Concise):**
```java
List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5);

// Lambda expression
Collections.sort(numbers, (a, b) -> a - b);
```

**Lambda Syntax:**
```java
// No parameters
() -> System.out.println("Hello")

// One parameter (parentheses optional)
x -> x * x
(x) -> x * x

// Multiple parameters
(a, b) -> a + b

// Multiple statements (need braces)
(a, b) -> {
    int sum = a + b;
    return sum;
}
```

**Practical Examples:**

```java
import java.util.*;

public class LambdaExamples {
    public static void main(String[] args) {
        List<String> names = Arrays.asList("Alice", "Bob", "Charlie");
        
        // forEach with lambda
        names.forEach(name -> System.out.println(name));
        
        // Filter and process
        names.stream()
             .filter(name -> name.startsWith("A"))
             .forEach(System.out::println);
        
        // Sorting with lambda
        names.sort((a, b) -> b.compareTo(a));  // Descending
        
        // Thread with lambda
        Thread thread = new Thread(() -> {
            System.out.println("Running in thread");
        });
        thread.start();
    }
}
```

---

### 2. Functional Interfaces

**Functional Interface:** Interface with exactly one abstract method

**Built-in Functional Interfaces:**

```java
import java.util.function.*;

public class FunctionalInterfacesDemo {
    public static void main(String[] args) {
        
        // 1. Predicate<T> - Takes T, returns boolean
        Predicate<Integer> isEven = n -> n % 2 == 0;
        System.out.println(isEven.test(4));  // true
        
        // 2. Function<T, R> - Takes T, returns R
        Function<String, Integer> length = s -> s.length();
        System.out.println(length.apply("Hello"));  // 5
        
        // 3. Consumer<T> - Takes T, returns nothing
        Consumer<String> printer = s -> System.out.println(s);
        printer.accept("Hello");
        
        // 4. Supplier<T> - Takes nothing, returns T
        Supplier<Double> randomSupplier = () -> Math.random();
        System.out.println(randomSupplier.get());
        
        // 5. BiFunction<T, U, R> - Takes T and U, returns R
        BiFunction<Integer, Integer, Integer> add = (a, b) -> a + b;
        System.out.println(add.apply(5, 3));  // 8
        
        // 6. UnaryOperator<T> - Takes T, returns T
        UnaryOperator<Integer> square = x -> x * x;
        System.out.println(square.apply(5));  // 25
        
        // 7. BinaryOperator<T> - Takes two T, returns T
        BinaryOperator<Integer> max = (a, b) -> a > b ? a : b;
        System.out.println(max.apply(10, 20));  // 20
    }
}
```

**Create Custom Functional Interface:**

```java
@FunctionalInterface
interface Calculator {
    int calculate(int a, int b);
}

public class CustomFunctionalInterface {
    public static void main(String[] args) {
        Calculator add = (a, b) -> a + b;
        Calculator multiply = (a, b) -> a * b;
        
        System.out.println(add.calculate(5, 3));       // 8
        System.out.println(multiply.calculate(5, 3));  // 15
    }
}
```

---

### 3. Stream API - Process Collections Declaratively

**What is Stream?**
- Sequence of elements supporting sequential and parallel operations
- **NOT** a data structure (doesn't store data)
- Allows functional-style operations

**Creating Streams:**

```java
import java.util.*;
import java.util.stream.*;

public class CreateStreams {
    public static void main(String[] args) {
        // From collection
        List<Integer> list = Arrays.asList(1, 2, 3, 4, 5);
        Stream<Integer> stream1 = list.stream();
        
        // From array
        String[] arr = {"a", "b", "c"};
        Stream<String> stream2 = Arrays.stream(arr);
        
        // Using Stream.of()
        Stream<Integer> stream3 = Stream.of(1, 2, 3, 4, 5);
        
        // Infinite stream with limit
        Stream<Integer> stream4 = Stream.iterate(0, n -> n + 2).limit(10);
        
        // Generate stream
        Stream<Double> stream5 = Stream.generate(Math::random).limit(5);
        
        // Range
        IntStream range = IntStream.range(1, 6);  // 1 to 5
        IntStream rangeClosed = IntStream.rangeClosed(1, 5);  // 1 to 5 (inclusive)
    }
}
```

---

**Stream Operations:**

#### Intermediate Operations (return Stream)

```java
import java.util.*;
import java.util.stream.*;

public class IntermediateOps {
    public static void main(String[] args) {
        List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
        
        // 1. filter - Select elements matching condition
        List<Integer> evens = numbers.stream()
                                    .filter(n -> n % 2 == 0)
                                    .collect(Collectors.toList());
        // [2, 4, 6, 8, 10]
        
        // 2. map - Transform each element
        List<Integer> squares = numbers.stream()
                                      .map(n -> n * n)
                                      .collect(Collectors.toList());
        // [1, 4, 9, 16, 25, ...]
        
        // 3. sorted - Sort elements
        List<Integer> sorted = numbers.stream()
                                     .sorted(Comparator.reverseOrder())
                                     .collect(Collectors.toList());
        
        // 4. distinct - Remove duplicates
        List<Integer> distinct = Arrays.asList(1, 2, 2, 3, 3, 3).stream()
                                       .distinct()
                                       .collect(Collectors.toList());
        // [1, 2, 3]
        
        // 5. limit - Take first N elements
        List<Integer> first3 = numbers.stream()
                                     .limit(3)
                                     .collect(Collectors.toList());
        // [1, 2, 3]
        
        // 6. skip - Skip first N elements
        List<Integer> skip3 = numbers.stream()
                                    .skip(3)
                                    .collect(Collectors.toList());
        // [4, 5, 6, 7, 8, 9, 10]
        
        // 7. peek - Debug/perform action without consuming
        numbers.stream()
               .peek(n -> System.out.println("Processing: " + n))
               .filter(n -> n > 5)
               .collect(Collectors.toList());
    }
}
```

---

#### Terminal Operations (produce result)

```java
import java.util.*;
import java.util.stream.*;

public class TerminalOps {
    public static void main(String[] args) {
        List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5);
        
        // 1. collect - Accumulate into collection
        List<Integer> list = numbers.stream().collect(Collectors.toList());
        Set<Integer> set = numbers.stream().collect(Collectors.toSet());
        
        // 2. forEach - Perform action on each element
        numbers.stream().forEach(System.out::println);
        
        // 3. count - Count elements
        long count = numbers.stream().count();  // 5
        
        // 4. reduce - Combine elements
        int sum = numbers.stream().reduce(0, (a, b) -> a + b);  // 15
        int product = numbers.stream().reduce(1, (a, b) -> a * b);  // 120
        
        // 5. min/max
        Optional<Integer> min = numbers.stream().min(Integer::compareTo);
        Optional<Integer> max = numbers.stream().max(Integer::compareTo);
        
        // 6. anyMatch/allMatch/noneMatch
        boolean hasEven = numbers.stream().anyMatch(n -> n % 2 == 0);    // true
        boolean allPositive = numbers.stream().allMatch(n -> n > 0);     // true
        boolean noneNegative = numbers.stream().noneMatch(n -> n < 0);   // true
        
        // 7. findFirst/findAny
        Optional<Integer> first = numbers.stream().findFirst();
        Optional<Integer> any = numbers.stream().findAny();
        
        // 8. toArray
        Integer[] arr = numbers.stream().toArray(Integer[]::new);
    }
}
```

---

**Stream API - Real Examples:**

```java
import java.util.*;
import java.util.stream.*;

public class StreamRealExamples {
    public static void main(String[] args) {
        
        // Example 1: Sum of even numbers
        List<Integer> nums = Arrays.asList(1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
        int sumEven = nums.stream()
                         .filter(n -> n % 2 == 0)
                         .mapToInt(Integer::intValue)
                         .sum();
        System.out.println("Sum of evens: " + sumEven);  // 30
        
        // Example 2: Find names starting with 'A'
        List<String> names = Arrays.asList("Alice", "Bob", "Anna", "Charlie");
        List<String> filteredNames = names.stream()
                                          .filter(name -> name.startsWith("A"))
                                          .collect(Collectors.toList());
        // [Alice, Anna]
        
        // Example 3: Convert to uppercase and sort
        List<String> upperSorted = names.stream()
                                       .map(String::toUpperCase)
                                       .sorted()
                                       .collect(Collectors.toList());
        // [ALICE, ANNA, BOB, CHARLIE]
        
        // Example 4: Group by length
        Map<Integer, List<String>> groupedByLength = names.stream()
                .collect(Collectors.groupingBy(String::length));
        // {3=[Bob], 4=[Anna], 5=[Alice], 7=[Charlie]}
        
        // Example 5: Partition by condition
        Map<Boolean, List<Integer>> partitioned = nums.stream()
                .collect(Collectors.partitioningBy(n -> n % 2 == 0));
        // {false=[1,3,5,7,9], true=[2,4,6,8,10]}
        
        // Example 6: Get statistics
        IntSummaryStatistics stats = nums.stream()
                .mapToInt(Integer::intValue)
                .summaryStatistics();
        System.out.println("Average: " + stats.getAverage());  // 5.5
        System.out.println("Max: " + stats.getMax());          // 10
        
        // Example 7: Joining strings
        String joined = names.stream()
                            .collect(Collectors.joining(", "));
        // "Alice, Bob, Anna, Charlie"
        
        // Example 8: flatMap - Flatten nested lists
        List<List<Integer>> nested = Arrays.asList(
            Arrays.asList(1, 2),
            Arrays.asList(3, 4),
            Arrays.asList(5, 6)
        );
        List<Integer> flattened = nested.stream()
                                       .flatMap(List::stream)
                                       .collect(Collectors.toList());
        // [1, 2, 3, 4, 5, 6]
    }
}
```

---

### 4. Method References - Shorthand for Lambdas

**Syntax:** `ClassName::methodName`

```java
import java.util.*;

public class MethodReferences {
    public static void main(String[] args) {
        List<String> names = Arrays.asList("Alice", "Bob", "Charlie");
        
        // 1. Reference to static method
        // Lambda: n -> Math.abs(n)
        // Method reference: Math::abs
        List<Integer> nums = Arrays.asList(-1, -2, -3);
        nums.stream().map(Math::abs).forEach(System.out::println);
        
        // 2. Reference to instance method of particular object
        // Lambda: s -> System.out.println(s)
        // Method reference: System.out::println
        names.forEach(System.out::println);
        
        // 3. Reference to instance method of arbitrary object
        // Lambda: (s1, s2) -> s1.compareTo(s2)
        // Method reference: String::compareTo
        names.sort(String::compareTo);
        
        // 4. Reference to constructor
        // Lambda: () -> new ArrayList<>()
        // Method reference: ArrayList::new
        List<String> newList = names.stream()
                                   .collect(ArrayList::new, 
                                           ArrayList::add, 
                                           ArrayList::addAll);
    }
}
```

---

### 5. Optional - Avoid Null Pointer Exceptions

**Problem with null:**
```java
String name = getName();  // Might return null
int length = name.length();  // NullPointerException!
```

**Solution with Optional:**

```java
import java.util.Optional;

public class OptionalDemo {
    public static void main(String[] args) {
        
        // Creating Optional
        Optional<String> empty = Optional.empty();
        Optional<String> name = Optional.of("Alice");  // Throws if null
        Optional<String> nullable = Optional.ofNullable(null);  // Safe for null
        
        // Check if value present
        if (name.isPresent()) {
            System.out.println(name.get());
        }
        
        // Better way - ifPresent
        name.ifPresent(System.out::println);
        
        // orElse - Provide default value
        String value = nullable.orElse("Default");
        
        // orElseGet - Provide supplier for default
        String value2 = nullable.orElseGet(() -> "Generated Default");
        
        // orElseThrow - Throw exception if empty
        String value3 = name.orElseThrow(() -> new RuntimeException("No value"));
        
        // map - Transform value if present
        Optional<Integer> length = name.map(String::length);
        
        // filter - Keep only if matches condition
        Optional<String> longName = name.filter(n -> n.length() > 5);
        
        // flatMap - Avoid Optional<Optional<T>>
        Optional<String> upper = name.flatMap(n -> Optional.of(n.toUpperCase()));
    }
    
    // Example: Optional in return type
    public static Optional<String> findNameById(int id) {
        if (id == 1) {
            return Optional.of("Alice");
        }
        return Optional.empty();
    }
    
    // Usage
    public static void demo() {
        findNameById(1)
            .map(String::toUpperCase)
            .ifPresent(System.out::println);  // ALICE
        
        String name = findNameById(2)
                        .orElse("Unknown");   // Unknown
    }
}
```

**Best Practices:**
- ✅ Use Optional as return type
- ❌ Don't use Optional as parameter
- ❌ Don't use Optional for fields
- ✅ Never call `get()` without checking `isPresent()`

---

### 6. Default Methods in Interfaces

**Before Java 8:** Can't add methods to interface without breaking implementations

**After Java 8:** Can add default implementations

```java
interface Vehicle {
    // Abstract method
    void start();
    
    // Default method (has implementation)
    default void stop() {
        System.out.println("Vehicle stopped");
    }
    
    // Static method
    static void checkVehicle() {
        System.out.println("Checking vehicle");
    }
}

class Car implements Vehicle {
    @Override
    public void start() {
        System.out.println("Car started");
    }
    
    // Can override default method
    @Override
    public void stop() {
        System.out.println("Car stopped");
    }
}

public class DefaultMethodDemo {
    public static void main(String[] args) {
        Car car = new Car();
        car.start();  // Car started
        car.stop();   // Car stopped
        
        Vehicle.checkVehicle();  // Checking vehicle
    }
}
```

---

### 7. Date and Time API (java.time)

**Old API Problems:** Date/Calendar classes are mutable and not thread-safe

**New API:** Immutable and thread-safe

```java
import java.time.*;
import java.time.format.DateTimeFormatter;

public class DateTimeDemo {
    public static void main(String[] args) {
        
        // LocalDate - Date without time
        LocalDate today = LocalDate.now();
        LocalDate specificDate = LocalDate.of(2024, 12, 25);
        LocalDate parsedDate = LocalDate.parse("2024-12-25");
        
        System.out.println("Today: " + today);
        System.out.println("Year: " + today.getYear());
        System.out.println("Month: " + today.getMonthValue());
        System.out.println("Day: " + today.getDayOfMonth());
        
        // LocalTime - Time without date
        LocalTime now = LocalTime.now();
        LocalTime specificTime = LocalTime.of(14, 30, 0);
        
        // LocalDateTime - Date and Time
        LocalDateTime dateTime = LocalDateTime.now();
        LocalDateTime specific = LocalDateTime.of(2024, 12, 25, 14, 30);
        
        // Period - Difference between dates
        LocalDate birth = LocalDate.of(1990, 1, 15);
        Period age = Period.between(birth, today);
        System.out.println("Age: " + age.getYears() + " years");
        
        // Duration - Difference between times
        LocalTime start = LocalTime.of(9, 0);
        LocalTime end = LocalTime.of(17, 0);
        Duration workHours = Duration.between(start, end);
        System.out.println("Work hours: " + workHours.toHours());
        
        // Adding/Subtracting
        LocalDate tomorrow = today.plusDays(1);
        LocalDate nextWeek = today.plusWeeks(1);
        LocalDate lastMonth = today.minusMonths(1);
        
        // Formatting
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("dd-MM-yyyy");
        String formatted = today.format(formatter);
        
        // Parsing
        LocalDate parsed = LocalDate.parse("25-12-2024", formatter);
    }
}
```

---

## Java 9 (Released 2017)

### 1. Module System (Project Jigsaw)

**Creates:** Better encapsulation and dependency management

```java
// module-info.java
module com.example.myapp {
    requires java.sql;
    exports com.example.myapp.api;
}
```

---

### 2. Factory Methods for Collections

**Before Java 9:**
```java
List<String> list = new ArrayList<>();
list.add("a");
list.add("b");
list.add("c");
List<String> immutable = Collections.unmodifiableList(list);
```

**After Java 9:**
```java
import java.util.*;

public class FactoryMethods {
    public static void main(String[] args) {
        
        // Immutable List
        List<String> list = List.of("a", "b", "c");
        
        // Immutable Set
        Set<Integer> set = Set.of(1, 2, 3);
        
        // Immutable Map
        Map<String, Integer> map = Map.of(
            "Alice", 25,
            "Bob", 30,
            "Charlie", 35
        );
        
        // Map.ofEntries for more entries
        Map<String, Integer> bigMap = Map.ofEntries(
            Map.entry("a", 1),
            Map.entry("b", 2),
            Map.entry("c", 3)
        );
        
        // These are IMMUTABLE - throws UnsupportedOperationException
        // list.add("d");  // Error!
    }
}
```

---

### 3. Stream API Improvements

```java
import java.util.stream.*;

public class StreamJava9 {
    public static void main(String[] args) {
        
        // takeWhile - Take elements while condition is true
        Stream.of(1, 2, 3, 4, 5, 1, 2)
              .takeWhile(n -> n < 4)
              .forEach(System.out::println);  // 1, 2, 3
        
        // dropWhile - Drop elements while condition is true
        Stream.of(1, 2, 3, 4, 5, 1, 2)
              .dropWhile(n -> n < 4)
              .forEach(System.out::println);  // 4, 5, 1, 2
        
        // iterate with predicate
        Stream.iterate(1, n -> n <= 10, n -> n + 1)
              .forEach(System.out::println);  // 1 to 10
        
        // ofNullable - Create stream from nullable
        Stream<String> stream = Stream.ofNullable(null);  // Empty stream
    }
}
```

---

### 4. Optional Improvements

```java
import java.util.Optional;
import java.util.stream.Stream;

public class OptionalJava9 {
    public static void main(String[] args) {
        
        // ifPresentOrElse - Execute action or else
        Optional<String> name = Optional.of("Alice");
        name.ifPresentOrElse(
            System.out::println,
            () -> System.out.println("No name")
        );
        
        // or - Provide alternative Optional
        Optional<String> empty = Optional.empty();
        Optional<String> result = empty.or(() -> Optional.of("Default"));
        
        // stream - Convert Optional to Stream
        Stream<String> stream = name.stream();
    }
}
```

---

## Java 10 (Released 2018)

### 1. Local Variable Type Inference (var)

**Type inference:** Compiler determines type automatically

```java
import java.util.*;

public class VarKeyword {
    public static void main(String[] args) {
        
        // Instead of: String name = "Alice";
        var name = "Alice";  // Compiler infers String
        
        // Instead of: List<String> list = new ArrayList<>();
        var list = new ArrayList<String>();
        
        // Instead of: Map<String, List<Integer>> map = new HashMap<>();
        var map = new HashMap<String, List<Integer>>();
        
        // Works with loops
        var numbers = List.of(1, 2, 3, 4, 5);
        for (var num : numbers) {
            System.out.println(num);
        }
        
        // Works with streams
        var result = numbers.stream()
                           .filter(n -> n > 2)
                           .collect(Collectors.toList());
    }
}
```

**When NOT to use var:**
```java
// ❌ BAD - Not clear what type
var data = getData();

// ✅ GOOD - Clear from right side
var name = "Alice";
var count = 10;
var list = new ArrayList<String>();

// ❌ Can't use for fields
// var field = 10;  // Error!

// ❌ Can't use for method parameters
// public void method(var param) { }  // Error!

// ❌ Can't use without initializer
// var x;  // Error!
```

---

### 2. Collection.copyOf()

```java
import java.util.*;

public class CopyOf {
    public static void main(String[] args) {
        List<String> list = new ArrayList<>(Arrays.asList("a", "b", "c"));
        
        // Create immutable copy
        List<String> copy = List.copyOf(list);
        
        // copy is immutable
        // copy.add("d");  // UnsupportedOperationException
        
        // Original can still be modified
        list.add("d");
        
        System.out.println(list);  // [a, b, c, d]
        System.out.println(copy);  // [a, b, c]
    }
}
```

---

## Java 11 (Released 2018) - LTS Version

### 1. String Methods

```java
public class StringJava11 {
    public static void main(String[] args) {
        
        // isBlank - Check if string is empty or whitespace
        System.out.println("".isBlank());       // true
        System.out.println("  ".isBlank());     // true
        System.out.println("a".isBlank());      // false
        
        // lines - Split by line terminators
        String multiline = "Line 1\nLine 2\nLine 3";
        multiline.lines().forEach(System.out::println);
        
        // strip - Remove leading/trailing whitespace (Unicode-aware)
        String text = "  Hello  ";
        System.out.println(text.strip());       // "Hello"
        System.out.println(text.stripLeading()); // "Hello  "
        System.out.println(text.stripTrailing());// "  Hello"
        
        // repeat - Repeat string N times
        System.out.println("*".repeat(10));     // **********
        System.out.println("abc".repeat(3));    // abcabcabc
    }
}
```

---

### 2. Files Methods

```java
import java.io.IOException;
import java.nio.file.*;

public class FilesJava11 {
    public static void main(String[] args) throws IOException {
        Path path = Path.of("test.txt");
        
        // writeString - Write string to file
        Files.writeString(path, "Hello World");
        
        // readString - Read entire file as string
        String content = Files.readString(path);
        System.out.println(content);  // Hello World
    }
}
```

---

### 3. Collection.toArray() Enhancement

```java
import java.util.*;

public class ToArrayJava11 {
    public static void main(String[] args) {
        List<String> list = Arrays.asList("a", "b", "c");
        
        // Java 10 and before
        String[] arr1 = list.toArray(new String[0]);
        
        // Java 11 - Cleaner!
        String[] arr2 = list.toArray(String[]::new);
    }
}
```

---

## Java 12-13 (Released 2019)

### 1. Switch Expressions (Preview in 12, Standard in 14)

**Old Switch (Statement):**
```java
String day = "MONDAY";
String result;
switch (day) {
    case "MONDAY":
    case "FRIDAY":
        result = "Working day";
        break;
    case "SATURDAY":
    case "SUNDAY":
        result = "Weekend";
        break;
    default:
        result = "Unknown";
}
```

**New Switch (Expression):**
```java
String day = "MONDAY";

// Returns value directly
String result = switch (day) {
    case "MONDAY", "FRIDAY" -> "Working day";
    case "SATURDAY", "SUNDAY" -> "Weekend";
    default -> "Unknown";
};

// With blocks
int numLetters = switch (day) {
    case "MONDAY", "FRIDAY", "SUNDAY" -> 6;
    case "TUESDAY" -> 7;
    case "THURSDAY", "SATURDAY" -> {
        System.out.println("Calculating...");
        yield 8;  // Use yield for blocks
    }
    default -> throw new IllegalArgumentException("Invalid day");
};
```

---

### 2. Text Blocks (Preview in 13, Standard in 15)

**Before Text Blocks:**
```java
String json = "{\n" +
              "  \"name\": \"Alice\",\n" +
              "  \"age\": 25\n" +
              "}";
```

**After Text Blocks:**
```java
public class TextBlocks {
    public static void main(String[] args) {
        
        // JSON
        String json = """
                {
                  "name": "Alice",
                  "age": 25
                }
                """;
        
        // SQL
        String sql = """
                SELECT id, name, email
                FROM users
                WHERE age > 18
                ORDER BY name
                """;
        
        // HTML
        String html = """
                <html>
                    <body>
                        <h1>Hello World</h1>
                    </body>
                </html>
                """;
        
        System.out.println(json);
    }
}
```

**Features:**
- No need for `\n`
- No need for `+` concatenation
- Maintains indentation
- Much more readable!

---

## Java 14 (Released 2020)

### 1. Pattern Matching for instanceof (Preview in 14, Standard in 16)

**Before:**
```java
Object obj = "Hello";

if (obj instanceof String) {
    String str = (String) obj;  // Manual casting
    System.out.println(str.length());
}
```

**After:**
```java
Object obj = "Hello";

if (obj instanceof String str) {  // Pattern variable
    System.out.println(str.length());  // No casting needed!
}

// Can use in same expression
if (obj instanceof String str && str.length() > 5) {
    System.out.println("Long string: " + str);
}
```

---

### 2. Helpful NullPointerExceptions

**Before Java 14:**
```
Exception in thread "main" java.lang.NullPointerException
    at Main.main(Main.java:5)
```

**After Java 14:**
```
Exception in thread "main" java.lang.NullPointerException: 
    Cannot invoke "String.length()" because "str" is null
    at Main.main(Main.java:5)
```

Shows exactly WHAT is null!

---

## Java 15 (Released 2020)

### 1. Sealed Classes (Preview in 15, Standard in 17)

**Control which classes can extend/implement**

```java
// Only Circle, Rectangle, Triangle can extend Shape
public sealed class Shape permits Circle, Rectangle, Triangle {
    // Common methods
}

// Must use: final, sealed, or non-sealed
final class Circle extends Shape {
    // Cannot be extended further
}

sealed class Rectangle extends Shape permits Square {
    // Only Square can extend Rectangle
}

non-sealed class Triangle extends Shape {
    // Anyone can extend Triangle
}

final class Square extends Rectangle {
}
```

**Why use sealed classes?**
- Better API design
- Exhaustive pattern matching
- Domain modeling

```java
// Compiler knows all possible types
double area = switch (shape) {
    case Circle c -> Math.PI * c.radius * c.radius;
    case Rectangle r -> r.length * r.width;
    case Triangle t -> 0.5 * t.base * t.height;
    // No default needed - compiler knows these are all possibilities!
};
```

---

## Java 16 (Released 2021)

### 1. Records - Data Classes

**Before Records:**
```java
public class Person {
    private final String name;
    private final int age;
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
    
    public String getName() { return name; }
    public int getAge() { return age; }
    
    @Override
    public boolean equals(Object o) {
        // ... lots of boilerplate
    }
    
    @Override
    public int hashCode() {
        // ... lots of boilerplate
    }
    
    @Override
    public String toString() {
        return "Person[name=" + name + ", age=" + age + "]";
    }
}
```

**After Records:**
```java
public record Person(String name, int age) {
    // That's it! Auto-generates:
    // - Constructor
    // - Getters (name(), age())
    // - equals(), hashCode(), toString()
}

public class RecordDemo {
    public static void main(String[] args) {
        Person person = new Person("Alice", 25);
        
        System.out.println(person.name());  // Alice
        System.out.println(person.age());   // 25
        System.out.println(person);  // Person[name=Alice, age=25]
        
        // Records are immutable
        // person.name = "Bob";  // Error! No setters
        
        // Custom methods allowed
    }
}

// Can add custom methods
record Point(int x, int y) {
    // Custom constructor
    public Point {
        if (x < 0 || y < 0) {
            throw new IllegalArgumentException("Negative coordinates");
        }
    }
    
    // Custom method
    public double distanceFromOrigin() {
        return Math.sqrt(x * x + y * y);
    }
}
```

---

### 2. Stream.toList()

```java
import java.util.*;
import java.util.stream.*;

public class StreamToList {
    public static void main(String[] args) {
        List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5);
        
        // Before Java 16
        List<Integer> evens1 = numbers.stream()
                                     .filter(n -> n % 2 == 0)
                                     .collect(Collectors.toList());
        
        // After Java 16 - Shorter!
        List<Integer> evens2 = numbers.stream()
                                     .filter(n -> n % 2 == 0)
                                     .toList();  // Returns immutable list
    }
}
```

---

## Java 17 (Released 2021) - LTS Version

### 1. Pattern Matching for switch (Preview)

```java
public class PatternMatchingSwitch {
    
    public static String formatValue(Object obj) {
        return switch (obj) {
            case Integer i -> "Integer: " + i;
            case String s -> "String: " + s;
            case Double d -> "Double: " + d;
            case null -> "null value";
            default -> "Unknown: " + obj.getClass().getName();
        };
    }
    
    public static void main(String[] args) {
        System.out.println(formatValue(42));        // Integer: 42
        System.out.println(formatValue("Hello"));   // String: Hello
        System.out.println(formatValue(3.14));      // Double: 3.14
        System.out.println(formatValue(null));      // null value
    }
}
```

---

### 2. Enhanced Pseudo-Random Number Generators

```java
import java.util.random.*;

public class RandomJava17 {
    public static void main(String[] args) {
        // New random generator interface
        RandomGenerator random = RandomGenerator.of("L64X128MixRandom");
        
        int num = random.nextInt(100);  // 0 to 99
        double d = random.nextDouble();
        
        // Stream of random numbers
        random.ints(10, 0, 100)
              .forEach(System.out::println);
    }
}
```

---

## Quick Comparison Table

| Feature | Java Version | Status | Use Case |
|---------|--------------|--------|----------|
| **Lambda Expressions** | 8 | Standard | Functional programming |
| **Stream API** | 8 | Standard | Processing collections |
| **Optional** | 8 | Standard | Avoiding null |
| **Method References** | 8 | Standard | Cleaner lambdas |
| **Date/Time API** | 8 | Standard | Working with dates |
| **var keyword** | 10 | Standard | Type inference |
| **String methods** | 11 | Standard | String processing |
| **Switch expressions** | 14 | Standard | Returning values from switch |
| **Text blocks** | 15 | Standard | Multi-line strings |
| **Pattern matching instanceof** | 16 | Standard | Avoiding casting |
| **Records** | 16 | Standard | Data classes |
| **Sealed classes** | 17 | Standard | Controlled inheritance |

---

## Top Interview Questions

### Q1: What are the major features of Java 8?
**Answer:** Lambda expressions, Stream API, Optional, Method references, Default methods in interfaces, New Date/Time API

---

### Q2: Difference between map() and flatMap()?
**Answer:**
- **map():** Transforms each element (one-to-one)
- **flatMap():** Flattens nested structures (one-to-many)

```java
// map - returns Stream<Integer[]>
Stream<Integer[]> stream1 = Stream.of(new int[]{1,2}, new int[]{3,4})
                                  .map(arr -> Arrays.stream(arr).boxed().toArray(Integer[]::new));

// flatMap - returns Stream<Integer>
Stream<Integer> stream2 = Stream.of(new int[]{1,2}, new int[]{3,4})
                                .flatMap(arr -> Arrays.stream(arr).boxed());
```

---

### Q3: What is a functional interface?
**Answer:** Interface with exactly one abstract method. Can have default and static methods. Used for lambda expressions.

```java
@FunctionalInterface
interface MyFunction {
    int apply(int x);  // Single abstract method
    
    default void print() { }  // Default methods allowed
    static void log() { }      // Static methods allowed
}
```

---

### Q4: Difference between Predicate and Function?
**Answer:**
- **Predicate<T>:** Takes T, returns boolean (for filtering)
- **Function<T, R>:** Takes T, returns R (for transforming)

---

### Q5: What's the difference between var in Java 10 vs var in JavaScript?
**Answer:**
- **Java var:** Compile-time type inference, strongly typed
- **JavaScript var:** Runtime, dynamic typing

---

### Q6: Are streams lazy or eager?
**Answer:** **Lazy!** Intermediate operations (filter, map) don't execute until terminal operation (collect, forEach) is called.

---

### Q7: Can we reuse a stream?
**Answer:** **No!** Streams can only be used once. Throws IllegalStateException if reused.

```java
Stream<Integer> stream = Stream.of(1, 2, 3);
stream.forEach(System.out::println);  // OK
stream.forEach(System.out::println);  // IllegalStateException!
```

---

### Q8: What's the difference between Collection and Stream?
**Answer:**
- **Collection:** Data structure, stores elements, can be iterated multiple times
- **Stream:** Processing pipeline, doesn't store data, one-time use

---

### Q9: When to use Optional?
**Answer:**
- ✅ As return type when value might be absent
- ❌ Never as method parameter
- ❌ Never as class field

---

### Q10: Difference between record and class?
**Answer:**
- **Record:** Immutable, auto-generates constructor/getters/equals/hashCode/toString
- **Class:** Mutable by default, manual boilerplate needed

---

## Key Takeaways for Competitive Programming

✅ **Stream API** - Master for processing collections efficiently
✅ **Lambda & Method References** - Write cleaner, shorter code
✅ **var keyword** - Reduce verbosity in local variables
✅ **Text blocks** - Handle multi-line strings easily
✅ **Records** - Quick data classes for test cases
✅ **Switch expressions** - Cleaner conditional logic
✅ **Optional** - Safer null handling

---

## What to Focus On

**For Competitive Programming:**
1. Stream API (most used!)
2. Lambda expressions
3. Method references
4. var keyword
5. String new methods

**For Interviews:**
1. All Java 8 features (mandatory!)
2. Stream API internals
3. Functional interfaces
4. var keyword (Java 10)
5. Records (Java 16)
6. Pattern matching

---

[← Back: Collections Framework](./Collections-Framework.md) | [Next: Comparators →](./Comparators.md)
