# OOP, Exception Handling & Modern Java — Complete Interview Guide
## Table of Contents
- [Part 1: Object-Oriented Programming](#part-1-object-oriented-programming)
  - [1.1 What is OOP? Why OOP?](#11-what-is-oop-why-oop)
  - [1.2 Class and Object](#12-class-and-object)
  - [1.3 Encapsulation](#13-encapsulation)
  - [1.4 Inheritance](#14-inheritance)
  - [1.5 Polymorphism](#15-polymorphism)
  - [1.6 Abstraction](#16-abstraction)
  - [1.7 Association, Aggregation, Composition](#17-association-aggregation-composition)
  - [1.8 SOLID Principles](#18-solid-principles)
  - [1.9 Real-World Design Examples](#19-real-world-design-examples)
- [Part 2: Java OOP Internals](#part-2-java-oop-internals)
  - [2.1 Object Creation in Memory](#21-object-creation-in-memory)
  - [2.2 Method Dispatch — Static vs Dynamic Binding](#22-method-dispatch--static-vs-dynamic-binding)
  - [2.3 equals() and hashCode()](#23-equals-and-hashcode)
  - [2.4 toString()](#24-tostring)
  - [2.5 final Keyword Deep Dive](#25-final-keyword-deep-dive)
  - [2.6 static Keyword Deep Dive](#26-static-keyword-deep-dive)
  - [2.7 this vs super](#27-this-vs-super)
  - [2.8 Object Class Methods Overview](#28-object-class-methods-overview)
- [Part 3: Exception Handling](#part-3-exception-handling)
  - [3.1 What is an Exception?](#31-what-is-an-exception)
  - [3.2 Exception Hierarchy](#32-exception-hierarchy)
  - [3.3 Checked vs Unchecked Exceptions](#33-checked-vs-unchecked-exceptions)
  - [3.4 Error vs Exception](#34-error-vs-exception)
  - [3.5 try, catch, finally](#35-try-catch-finally)
  - [3.6 Multiple Catch Blocks and Order](#36-multiple-catch-blocks-and-order)
  - [3.7 throw vs throws](#37-throw-vs-throws)
  - [3.8 Custom Exceptions](#38-custom-exceptions)
  - [3.9 try-with-resources](#39-try-with-resources)
  - [3.10 Exception Propagation](#310-exception-propagation)
- [Part 4: Advanced Exception Handling](#part-4-advanced-exception-handling)
  - [4.1 Best Practices — DOs and DON'Ts](#41-best-practices--dos-and-donts)
  - [4.2 When NOT to Catch Exceptions](#42-when-not-to-catch-exceptions)
  - [4.3 Logging vs Rethrowing](#43-logging-vs-rethrowing)
  - [4.4 Fail-Fast vs Fail-Safe](#44-fail-fast-vs-fail-safe)
  - [4.5 Real-World Exception Design](#45-real-world-exception-design)
- [Part 5: Modern Java Features (14+ / 17+ / 21)](#part-5-modern-java-features-14--17--21)
  - [5.1 var — Local Variable Type Inference](#51-var--local-variable-type-inference)
  - [5.2 Records](#52-records)
  - [5.3 Sealed Classes](#53-sealed-classes)
  - [5.4 Pattern Matching](#54-pattern-matching)
  - [5.5 Updated Switch Expressions](#55-updated-switch-expressions)
  - [5.6 Virtual Threads (Project Loom — Java 21)](#56-virtual-threads-project-loom--java-21)
- [6. Quick Revision Cheat Sheet](#6-quick-revision-cheat-sheet)
- [7. Interview Questions with Answers (20)](#7-interview-questions-with-answers-20)
---
# Part 1: Object-Oriented Programming
## 1.1 What is OOP? Why OOP?
Object-Oriented Programming is a paradigm that organizes software design around **objects** (data + behavior) rather than functions and logic.
**Why OOP over procedural programming?**
| Aspect | Procedural | OOP |
|---|---|---|
| Focus | Functions/procedures | Objects and interactions |
| Data security | Global data, open access | Encapsulated, controlled access |
| Code reuse | Copy-paste or function libraries | Inheritance, composition |
| Scalability | Hard to scale | Modular, extensible |
| Real-world modeling | Poor | Natural mapping |
| Maintenance | Change ripples everywhere | Localized changes |
> **Real-World Analogy:** Procedural programming is like writing a single recipe where all steps are listed sequentially. OOP is like organizing a kitchen — the oven, refrigerator, and blender are each **objects** with their own internal state and operations. You interact with them through defined interfaces (buttons, dials) without knowing their internals.
**Four Pillars of OOP:**
1. **Encapsulation** — Bundle data + methods; hide internals
2. **Inheritance** — Acquire properties from parent
3. **Polymorphism** — Same interface, different behavior
4. **Abstraction** — Expose what, hide how
## 1.2 Class and Object
A **class** is a blueprint/template. An **object** is an instance of that class created at runtime.
```java
public class Employee {
    private String name;
    private double salary;
    public Employee(String name, double salary) {
        this.name = name;
        this.salary = salary;
    }
    public double annualSalary() {
        return salary * 12;
    }
}
Employee emp = new Employee("Alice", 80000); // object creation
```
### Memory-Level Explanation
```
Stack (Thread-specific)                  Heap (Shared)
┌─────────────────────┐                ┌──────────────────────────┐
│ main()              │                │  Employee Object         │
│  emp ──────────────────────────────►│  name ──► "Alice"        │
│  (reference: 0xA1)  │                │  salary: 80000.0         │
└─────────────────────┘                │  (Object header: class   │
                                       │   pointer, hashCode,     │
                                       │   lock info, GC age)     │
                                       └──────────────────────────┘
```
**What happens with `new Employee("Alice", 80000)`?**
1. **Class loading** — JVM loads `Employee.class` if not already loaded.
2. **Memory allocation** — JVM allocates memory on the **heap** for the object (fields + object header).
3. **Default initialization** — Fields set to defaults (`null`, `0.0`).
4. **Constructor execution** — Fields assigned actual values.
5. **Reference returned** — The memory address is stored in the stack variable `emp`.
> **Interview Note:** The reference variable `emp` lives on the **stack**. The actual object lives on the **heap**. When `emp` goes out of scope, the object becomes eligible for GC (if no other references exist).
## 1.3 Encapsulation
Encapsulation = **data hiding** + **controlled access** through methods.
> **Real-World Analogy:** A bank account. You can't directly modify your balance — you use deposit() and withdraw() methods that enforce rules (minimum balance, overdraft limits).
```java
public class BankAccount {
    private double balance;  // hidden from outside
    public BankAccount(double initialBalance) {
        if (initialBalance < 0) throw new IllegalArgumentException("Balance cannot be negative");
        this.balance = initialBalance;
    }
    public double getBalance() {
        return balance;
    }
    public void deposit(double amount) {
        if (amount <= 0) throw new IllegalArgumentException("Deposit must be positive");
        balance += amount;
    }
    public void withdraw(double amount) {
        if (amount > balance) throw new InsufficientFundsException("Not enough balance");
        balance -= amount;
    }
}
```
**Benefits:**
- **Validation** — prevent invalid state (`negative balance`).
- **Flexibility** — internal representation can change without affecting clients.
- **Security** — only expose what's necessary.
### Creating Immutable Classes
```java
public final class Money {
    private final String currency;
    private final BigDecimal amount;
    public Money(String currency, BigDecimal amount) {
        this.currency = currency;
        this.amount = amount;
    }
    public String getCurrency() { return currency; }
    public BigDecimal getAmount() { return amount; }
    public Money add(Money other) {
        if (!this.currency.equals(other.currency))
            throw new IllegalArgumentException("Currency mismatch");
        return new Money(currency, this.amount.add(other.amount));
    }
}
```
**Rules for Immutability:**
1. Class declared `final` (no subclass can break contract).
2. All fields `private final`.
3. No setters.
4. **Defensive copies** for mutable fields (e.g., `Date`, `List`).
5. Return new objects instead of modifying existing ones.
> **Interview Trap:** `final` reference ≠ immutable object. A `final List<String>` can still have items added/removed. You need `Collections.unmodifiableList()` or `List.of()` for true immutability.
## 1.4 Inheritance
Inheritance allows a class to **acquire** the fields and methods of another class.
> **Real-World Analogy:** A `SavingsAccount` **is-a** `BankAccount`. It inherits all base behavior (deposit, withdraw) and adds interest calculation.
```java
public class Animal {
    protected String name;
    public Animal(String name) { this.name = name; }
    public void eat() { System.out.println(name + " is eating"); }
    public void sound() { System.out.println("Some generic sound"); }
}
public class Dog extends Animal {
    private String breed;
    public Dog(String name, String breed) {
        super(name);         // MUST call parent constructor
        this.breed = breed;
    }
    @Override
    public void sound() {   // method overriding
        System.out.println(name + " says: Woof!");
    }
    public void fetch() {   // new behavior
        System.out.println(name + " fetches the ball");
    }
}
```
### Types of Inheritance in Java
| Type | Java Support | Example |
|---|---|---|
| **Single** | ✅ | `class Dog extends Animal` |
| **Multilevel** | ✅ | `class Puppy extends Dog extends Animal` |
| **Hierarchical** | ✅ | `Dog extends Animal`, `Cat extends Animal` |
| **Multiple (classes)** | ❌ | Not supported (diamond problem) |
| **Multiple (interfaces)** | ✅ | `class X implements A, B` |
### Diamond Problem and Java's Solution
```
       Animal
      /      \
   Dog        Cat
      \      /
       Pet        ← Which Animal constructor/method to use?
```
Java **prevents** multiple class inheritance. For interfaces, Java resolves conflicts with explicit rules:
```java
interface Flyable {
    default void move() { System.out.println("Flying"); }
}
interface Swimmable {
    default void move() { System.out.println("Swimming"); }
}
class Duck implements Flyable, Swimmable {
    @Override
    public void move() {
        Flyable.super.move();    // explicitly choose, or provide own implementation
    }
}
```
### Method Overriding Rules
| Rule | Requirement |
|---|---|
| Method name | Must be identical |
| Parameters | Must be identical (count, type, order) |
| Return type | Same or **covariant** (subtype) |
| Access modifier | Same or **wider** (not narrower) |
| Exceptions | Same or **narrower** checked exceptions |
| `static` methods | Cannot override (hidden, not overridden) |
| `final` methods | Cannot override |
| `private` methods | Cannot override (not inherited) |
| Constructor | Cannot override (not inherited) |
```java
class Parent {
    protected Number calculate() throws IOException {
        return 42;
    }
}
class Child extends Parent {
    @Override
    public Integer calculate() throws FileNotFoundException { // ✅ All valid:
        return 42;                                            // wider access (protected → public)
    }                                                          // covariant return (Number → Integer)
}                                                              // narrower exception (IOException → FileNotFoundException)
```
> **Interview Gotcha:** `static` methods are **hidden**, not overridden. The method called depends on the **reference type**, not the object type.
```java
class Parent {
    static void greet() { System.out.println("Parent"); }
}
class Child extends Parent {
    static void greet() { System.out.println("Child"); }
}
Parent p = new Child();
p.greet();  // "Parent" — resolved at compile time based on reference type
```
## 1.5 Polymorphism
**Polymorphism** = "many forms." The same method call behaves differently depending on the object.
### Compile-Time Polymorphism (Method Overloading)
Resolved at **compile time** by the compiler based on method signature.
```java
public class Calculator {
    public int add(int a, int b) { return a + b; }
    public double add(double a, double b) { return a + b; }
    public int add(int a, int b, int c) { return a + b + c; }
    public String add(String a, String b) { return a + b; }
}
```
**Overloading Rules:**
- Must differ in **parameter list** (number, type, or order).
- Return type alone is **not** sufficient to overload.
- Access modifier and exceptions can differ freely.
> **Tricky Case — Method Resolution with Widening, Autoboxing, Varargs:**
```java
public class Overload {
    void test(int a)       { System.out.println("int"); }
    void test(long a)      { System.out.println("long"); }
    void test(Integer a)   { System.out.println("Integer"); }
    void test(int... a)    { System.out.println("varargs"); }
}
Overload o = new Overload();
o.test(5);
// Priority: exact match (int) → widening (long) → autoboxing (Integer) → varargs (int...)
// Output: "int"
```
**Resolution priority:** Exact match → Widening → Autoboxing → Varargs.
### Runtime Polymorphism (Method Overriding)
Resolved at **runtime** by the JVM based on the actual object type.
```java
Animal animal = new Dog("Rex", "Labrador");
animal.sound();  // "Rex says: Woof!" — Dog's overridden method called at runtime
animal.eat();    // "Rex is eating" — inherited from Animal
// animal.fetch(); // ❌ Compile error: Animal reference doesn't know about fetch()
```
```
              Compile Time                        Runtime
              ──────────                         ─────────
Reference:    Animal                             Object: Dog
Checks:       Does Animal have sound()? ✅       Calls: Dog.sound() ✅
              Does Animal have fetch()? ❌        (vtable lookup)
```
> **Real-World Analogy:** A TV remote (reference) can control any TV (object) through the same buttons (methods). The behavior depends on which TV is connected — Samsung shows Samsung UI, Sony shows Sony UI.
## 1.6 Abstraction
Abstraction = exposing **what** an object does, hiding **how** it does it.
### Abstract Classes vs Interfaces
```java
// Abstract class — partial implementation
public abstract class Payment {
    protected double amount;
    public Payment(double amount) {
        this.amount = amount;
    }
    public abstract boolean process();    // subclass MUST implement
    public void logTransaction() {        // concrete method — shared logic
        System.out.println("Processed: $" + amount);
    }
}
public class CreditCardPayment extends Payment {
    private String cardNumber;
    public CreditCardPayment(double amount, String cardNumber) {
        super(amount);
        this.cardNumber = cardNumber;
    }
    @Override
    public boolean process() {
        // call payment gateway, charge card
        return true;
    }
}
```
```java
// Interface — pure contract (before Java 8: only abstract methods)
public interface Notifiable {
    void send(String message, String recipient);
    default void sendBulk(List<String> messages, String recipient) {   // Java 8+
        messages.forEach(msg -> send(msg, recipient));
    }
    static Notifiable email() {    // Java 8+ static factory
        return (msg, to) -> System.out.println("Email to " + to + ": " + msg);
    }
}
```
### Comprehensive Comparison
| Feature | Abstract Class | Interface |
|---|---|---|
| **Instantiation** | ❌ Cannot | ❌ Cannot |
| **Constructors** | ✅ Yes | ❌ No |
| **Fields** | Any type (instance, static) | Only `public static final` |
| **Method types** | Abstract + concrete | Abstract + `default` + `static` + `private` (9+) |
| **Access modifiers** | Any | Methods: `public` (default) or `private` (9+); Fields: `public static final` |
| **Multiple inheritance** | ❌ Single class only | ✅ Implement multiple interfaces |
| **Inheritance keyword** | `extends` | `implements` |
| **State** | ✅ Can hold instance state | ❌ No instance state |
| **When to use** | Shared state/code among related classes | Define a contract for unrelated classes |
> **Decision Guide:**
> - Use **interface** when multiple unrelated classes need the same capability (`Serializable`, `Comparable`).
> - Use **abstract class** when related classes share state and behavior (`InputStream` → `FileInputStream`, `ByteArrayInputStream`).
> - Prefer interfaces (design to an interface, not implementation).
## 1.7 Association, Aggregation, Composition
These describe **relationships** between objects — critical for system design interviews.
### Association (Uses-A)
A general relationship where one object **uses** or **interacts with** another. Neither owns the other.
```java
class Teacher {
    private String name;
    void teach(Student student) {    // association: Teacher uses Student
        System.out.println(name + " teaches " + student.getName());
    }
}
class Student {
    private String name;
    String getName() { return name; }
}
```
### Aggregation (Has-A, Weak Ownership)
A whole-part relationship where the part can **exist independently** of the whole.
> **Real-World Analogy:** A **Department** has **Professors**. If the department is dissolved, the professors still exist — they can join another department.
```java
class Department {
    private String name;
    private List<Professor> professors;  // aggregation: Department HAS professors
    public Department(String name, List<Professor> professors) {
        this.professors = professors;     // professors created outside, passed in
    }
}
class Professor {
    private String name;
    // exists independently of Department
}
// Usage
Professor p1 = new Professor("Dr. Smith");
Professor p2 = new Professor("Dr. Jones");
Department cs = new Department("CS", List.of(p1, p2));
// If cs is garbage collected, p1 and p2 still exist
```
### Composition (Has-A, Strong Ownership)
A whole-part relationship where the part **cannot exist** without the whole. The whole **creates and owns** the part.
> **Real-World Analogy:** A **House** has **Rooms**. If the house is demolished, the rooms cease to exist.
```java
class Engine {
    private String type;
    Engine(String type) { this.type = type; }
}
class Car {
    private final Engine engine;
    public Car(String engineType) {
        this.engine = new Engine(engineType);  // Car CREATES the engine
    }
}
// Engine lifecycle is tied to Car
// When Car is GC'd, its Engine is also GC'd (no external reference)
```
### Summary Table
| Relationship | Ownership | Lifecycle | Example |
|---|---|---|---|
| **Association** | None | Independent | Teacher – Student |
| **Aggregation** | Weak (has-a) | Independent | Department – Professor |
| **Composition** | Strong (has-a) | Dependent | Car – Engine, House – Room |
| **Inheritance** | N/A (is-a) | Coupled | Dog – Animal |
> **Interview Tip:** Composition is preferred over inheritance ("Favor composition over inheritance" — Effective Java, Item 18). Inheritance breaks encapsulation; composition provides flexibility.
## 1.8 SOLID Principles
### S — Single Responsibility Principle
> A class should have **one, and only one, reason to change**.
```java
// ❌ Violates SRP: UserService handles business logic AND email sending
class UserService {
    void registerUser(User user) { /* save to DB */ }
    void sendWelcomeEmail(User user) { /* send email */ }
}
// ✅ Follows SRP: separate concerns
class UserService {
    void registerUser(User user) { /* save to DB */ }
}
class EmailService {
    void sendWelcomeEmail(User user) { /* send email */ }
}
```
### O — Open/Closed Principle
> Classes should be **open for extension**, **closed for modification**.
```java
// ❌ Adding a new shape requires modifying AreaCalculator
class AreaCalculator {
    double calculate(Object shape) {
        if (shape instanceof Circle c) return Math.PI * c.radius * c.radius;
        if (shape instanceof Rectangle r) return r.width * r.height;
        // adding Triangle means modifying this class
        return 0;
    }
}
// ✅ New shapes extend without modifying existing code
interface Shape {
    double area();
}
class Circle implements Shape {
    double radius;
    public double area() { return Math.PI * radius * radius; }
}
class Triangle implements Shape {
    double base, height;
    public double area() { return 0.5 * base * height; }
}
```
### L — Liskov Substitution Principle
> Subtypes must be substitutable for their base types **without altering program correctness**.
```java
// ❌ Violates LSP: Square changes the behavior contract of Rectangle
class Rectangle {
    protected int width, height;
    void setWidth(int w) { width = w; }
    void setHeight(int h) { height = h; }
    int area() { return width * height; }
}
class Square extends Rectangle {
    @Override
    void setWidth(int w) { width = w; height = w; }  // breaks parent's contract
    @Override
    void setHeight(int h) { width = h; height = h; }
}
Rectangle r = new Square();
r.setWidth(5);
r.setHeight(3);
r.area(); // Expected: 15. Actual: 9. ← LSP violation
```
### I — Interface Segregation Principle
> Clients should not be forced to depend on methods they don't use.
```java
// ❌ Fat interface: forces Robot to implement eat()
interface Worker {
    void work();
    void eat();
}
// ✅ Segregated interfaces
interface Workable { void work(); }
interface Feedable { void eat(); }
class Human implements Workable, Feedable { /* both */ }
class Robot implements Workable { /* only work */ }
```
### D — Dependency Inversion Principle
> High-level modules should not depend on low-level modules. Both should depend on **abstractions**.
```java
// ❌ High-level OrderService depends on concrete MySQLRepository
class OrderService {
    private MySQLRepository repo = new MySQLRepository();
}
// ✅ Both depend on abstraction
interface OrderRepository { void save(Order order); }
class MySQLRepository implements OrderRepository { /* ... */ }
class MongoRepository implements OrderRepository { /* ... */ }
class OrderService {
    private final OrderRepository repo;  // depends on abstraction
    OrderService(OrderRepository repo) { this.repo = repo; }
}
```
## 1.9 Real-World Design Examples
### E-Commerce Order System
```java
// Demonstrates: Encapsulation, Composition, Polymorphism, SRP
public class Order {
    private final String orderId;
    private final List<OrderItem> items;         // composition
    private final PaymentMethod paymentMethod;   // polymorphism via interface
    private OrderStatus status;
    public Order(String orderId, List<OrderItem> items, PaymentMethod paymentMethod) {
        this.orderId = orderId;
        this.items = List.copyOf(items);          // defensive copy
        this.paymentMethod = paymentMethod;
        this.status = OrderStatus.CREATED;
    }
    public double totalAmount() {
        return items.stream().mapToDouble(OrderItem::subtotal).sum();
    }
    public boolean checkout() {
        boolean paid = paymentMethod.pay(totalAmount());
        this.status = paid ? OrderStatus.PAID : OrderStatus.PAYMENT_FAILED;
        return paid;
    }
}
interface PaymentMethod {
    boolean pay(double amount);
}
class CreditCard implements PaymentMethod {
    public boolean pay(double amount) { /* charge card */ return true; }
}
class UPI implements PaymentMethod {
    public boolean pay(double amount) { /* UPI flow */ return true; }
}
```
---
# Part 2: Java OOP Internals
## 2.1 Object Creation in Memory
### What `new` Does Internally
```java
Employee emp = new Employee("Alice", 80000);
```
```
Step 1: Class Loading
  ┌─────────────────────────────────────────────────────────────┐
  │  Method Area (Metaspace since Java 8)                       │
  │  ┌───────────────────────────────────────────────────────┐  │
  │  │  Employee.class metadata                              │  │
  │  │  - Field descriptors (name: String, salary: double)   │  │
  │  │  - Method bytecode (annualSalary, constructor)        │  │
  │  │  - Constant pool (string literals, class refs)        │  │
  │  │  - vtable (virtual method table for polymorphism)     │  │
  │  └───────────────────────────────────────────────────────┘  │
  └─────────────────────────────────────────────────────────────┘
Step 2-5: Object Allocation and Initialization
  Stack                              Heap
  ┌──────────────┐                  ┌──────────────────────────────┐
  │ main()       │                  │  Employee Object @ 0xA1      │
  │  emp = 0xA1 ─┼─────────────────►│  [Object Header]             │
  │              │                  │    Mark Word: hashCode, lock  │
  │              │                  │    Klass Pointer → Employee   │
  │              │                  │  [Instance Data]              │
  │              │                  │    name → "Alice" (String obj)│
  │              │                  │    salary = 80000.0           │
  │              │                  │  [Padding] (alignment to 8B)  │
  └──────────────┘                  └──────────────────────────────┘
```
**Object Header (typically 12 bytes on 64-bit JVM with compressed oops):**
- **Mark Word** (8 bytes): hash code, GC age, lock state, biased locking info.
- **Klass Pointer** (4 bytes compressed): pointer to class metadata in Metaspace.
- **Array Length** (4 bytes, only for arrays).
> **Interview Note:** Even an empty `new Object()` consumes **16 bytes** on the heap (12-byte header + 4-byte padding for 8-byte alignment).
### Constructor Chaining
```java
class A {
    A() { System.out.println("A constructor"); }
}
class B extends A {
    B() { System.out.println("B constructor"); }  // implicit super() as first line
}
class C extends B {
    C() { System.out.println("C constructor"); }  // implicit super() as first line
}
new C();
// Output: A constructor → B constructor → C constructor
// Constructors chain from top of hierarchy downward
```
## 2.2 Method Dispatch — Static vs Dynamic Binding
### Static Binding (Early Binding) — Compile Time
Used for: `private`, `static`, `final` methods, and method **overloading**.
```java
class Calculator {
    static int add(int a, int b) { return a + b; }
    static double add(double a, double b) { return a + b; }
}
Calculator.add(3, 5);     // resolved at compile time → add(int, int)
Calculator.add(3.0, 5.0); // resolved at compile time → add(double, double)
```
### Dynamic Binding (Late Binding) — Runtime
Used for: **overridden** instance methods. JVM uses the **vtable** (virtual method table).
```java
Animal a = new Dog("Rex", "Lab");
a.sound(); // Compile time: "Does Animal have sound()? Yes." → compiles
           // Runtime: object is Dog → Dog's vtable → Dog.sound() executed
```
```
vtable for Dog:
┌──────────────┬────────────────────────┐
│ Method       │ Points To              │
├──────────────┼────────────────────────┤
│ toString()   │ Object.toString()      │
│ eat()        │ Animal.eat()           │
│ sound()      │ Dog.sound() ← override │
│ fetch()      │ Dog.fetch()            │
└──────────────┴────────────────────────┘
```
> **Interview Note:** Fields are **never** polymorphic. Field access is resolved at compile time based on reference type.
```java
class Parent { int x = 10; }
class Child extends Parent { int x = 20; }
Parent p = new Child();
System.out.println(p.x); // 10 (Parent's field, compile-time resolution)
```
## 2.3 equals() and hashCode()
### The Contract
1. If `a.equals(b)` is `true`, then `a.hashCode() == b.hashCode()` **must** be `true`.
2. If `a.hashCode() == b.hashCode()`, `a.equals(b)` is **not necessarily** `true` (hash collisions).
3. If you override `equals()`, you **must** override `hashCode()`.
### Default Behavior (Object class)
```java
// Object.equals() — compares references (same as ==)
public boolean equals(Object obj) {
    return (this == obj);
}
// Object.hashCode() — typically derived from memory address (native method)
```
### Correct Implementation
```java
public class Employee {
    private final int id;
    private final String name;
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;                     // same reference
        if (o == null || getClass() != o.getClass()) return false;  // null or different class
        Employee emp = (Employee) o;
        return id == emp.id && Objects.equals(name, emp.name);
    }
    @Override
    public int hashCode() {
        return Objects.hash(id, name);  // must use same fields as equals()
    }
}
```
> **What breaks if you don't override hashCode()?**
```java
Set<Employee> set = new HashSet<>();
Employee e1 = new Employee(1, "Alice");
Employee e2 = new Employee(1, "Alice");
set.add(e1);
set.contains(e2); // false! Different hashCode → different bucket → never finds e1
// Even though e1.equals(e2) is true
```
**equals() Rules (from the spec):**
- **Reflexive:** `x.equals(x)` → `true`
- **Symmetric:** `x.equals(y)` ↔ `y.equals(x)`
- **Transitive:** `x.equals(y)` && `y.equals(z)` → `x.equals(z)`
- **Consistent:** Multiple calls return the same result
- **Null:** `x.equals(null)` → `false`
> **Interview Trap — `equals()` with inheritance:**
```java
// Using instanceof is dangerous across inheritance hierarchies
class Point {
    int x, y;
    @Override
    public boolean equals(Object o) {
        if (!(o instanceof Point p)) return false;
        return x == p.x && y == p.y;
    }
}
class ColorPoint extends Point {
    String color;
    @Override
    public boolean equals(Object o) {
        if (!(o instanceof ColorPoint cp)) return false;
        return super.equals(cp) && Objects.equals(color, cp.color);
    }
}
Point p = new Point(1, 2);
ColorPoint cp = new ColorPoint(1, 2, "RED");
p.equals(cp);  // true  (Point.equals sees a Point)
cp.equals(p);  // false (ColorPoint.equals needs a ColorPoint)
// Symmetry violated! Use getClass() instead of instanceof for strict equality.
```
## 2.4 toString()
```java
// Default Object.toString()
// Returns: "com.example.Employee@1a2b3c4d" (className@hexHashCode)
@Override
public String toString() {
    return "Employee{id=" + id + ", name='" + name + "'}";
}
```
**Why it matters:**
- Used by `System.out.println()`, loggers, debuggers, string concatenation.
- Without overriding, you get useless `ClassName@hash` output.
- Records (Java 16+) auto-generate `toString()`.
## 2.5 final Keyword Deep Dive
| Applied To | Effect | Example |
|---|---|---|
| **Variable** | Cannot be reassigned | `final int MAX = 100;` |
| **Reference** | Reference can't change; object contents can | `final List<String> list = new ArrayList<>();` |
| **Blank final** | Must be assigned exactly once in constructor | `final int x; // assigned in constructor` |
| **Method** | Cannot be overridden by subclasses | `public final void critical() {}` |
| **Class** | Cannot be extended | `public final class String {}` |
| **Parameter** | Cannot be reassigned inside method | `void foo(final int x) {}` |
```java
// Blank final — useful for dependency injection pattern
class Service {
    private final Repository repo; // blank final
    Service(Repository repo) {
        this.repo = repo;  // assigned once in constructor ✅
    }
}
```
> **Interview Note:** All variables used inside a **lambda** or **anonymous class** must be effectively final (not reassigned after initialization), even without the `final` keyword.
```java
int count = 0;
Runnable r = () -> System.out.println(count); // ✅ effectively final
count = 1; // ❌ now it's NOT effectively final → compile error in lambda
```
## 2.6 static Keyword Deep Dive
### Static Context Rules
```
┌──────────────────────────────────────────────────────┐
│  static members → loaded with class                  │
│  instance members → created with each new object     │
│                                                      │
│  static CAN access:     static members only          │
│  static CANNOT access:  instance members, this/super │
│  instance CAN access:   both static and instance     │
└──────────────────────────────────────────────────────┘
```
### Static Initialization Order
```java
class Demo {
    static int x = initX();      // 1st: static field initializer
    static { x += 10; }          // 2nd: static block
    static int y = initY();      // 3rd: another static field
    int a = initA();             // 4th (per instance): instance field initializer
    { a += 100; }                // 5th (per instance): instance initializer block
    Demo() { a += 1000; }       // 6th (per instance): constructor
    static int initX() { System.out.println("initX"); return 1; }
    static int initY() { System.out.println("initY"); return 2; }
    int initA() { System.out.println("initA"); return 3; }
}
```
**Execution order:**
1. Static fields and static blocks — **in order of appearance** (once, at class loading).
2. Instance fields and instance blocks — **in order of appearance** (every `new`).
3. Constructor body.
### Static Method Cannot Be Overridden
```java
class Parent {
    static void hello() { System.out.println("Parent"); }
}
class Child extends Parent {
    static void hello() { System.out.println("Child"); }  // method HIDING, not overriding
}
Parent p = new Child();
p.hello(); // "Parent" — resolved by reference type at compile time
```
## 2.7 this vs super
| Keyword | Refers To | Usage |
|---|---|---|
| `this` | Current object instance | Access instance members, call other constructors |
| `super` | Parent class | Access parent members, call parent constructor |
```java
class Animal {
    String name;
    Animal(String name) { this.name = name; }
    void info() { System.out.println("Animal: " + name); }
}
class Dog extends Animal {
    String breed;
    Dog(String name, String breed) {
        super(name);          // MUST be first statement if calling super()
        this.breed = breed;
    }
    Dog() {
        this("Unknown", "Mixed"); // this() calls another constructor; MUST be first statement
    }
    @Override
    void info() {
        super.info();        // call parent's info()
        System.out.println("Breed: " + breed);
    }
}
```
**Rules:**
- `this()` and `super()` **cannot** both appear in the same constructor.
- Both must be the **first statement** in the constructor.
- `this` and `super` cannot be used in a `static` context.
## 2.8 Object Class Methods Overview
Every class in Java implicitly extends `java.lang.Object`. Key methods:
| Method | Purpose | When to Override |
|---|---|---|
| `equals(Object)` | Logical equality | When objects should be compared by value |
| `hashCode()` | Hash for hash-based collections | Always with `equals()` |
| `toString()` | String representation | Almost always (debugging, logging) |
| `getClass()` | Runtime class info | Never (final method) |
| `clone()` | Create a copy | Rarely (use copy constructors instead) |
| `finalize()` | Pre-GC cleanup | Never (deprecated since Java 9) |
| `wait()` / `notify()` / `notifyAll()` | Thread synchronization | Advanced concurrency |
> **Interview Note:** `clone()` performs **shallow copy** by default. For deep copy, override and manually clone nested objects, or use serialization/copy constructors.
```java
// Prefer copy constructor over clone()
public class Address {
    private String city;
    private String street;
    // Copy constructor
    public Address(Address other) {
        this.city = other.city;
        this.street = other.street;
    }
}
```
---
# Part 3: Exception Handling
## 3.1 What is an Exception?
An exception is an **event that disrupts the normal flow** of program execution. Java uses an object-oriented approach — exceptions are objects that encapsulate error information.
> **Real-World Analogy:** An exception is like a fire alarm in a building. Normal operations (people working) are interrupted. The alarm (exception) is caught by the fire response team (catch block), who handle it (log, evacuate). If nobody handles it, the building shuts down (program crashes).
## 3.2 Exception Hierarchy
```
java.lang.Object
  └── java.lang.Throwable
        ├── java.lang.Error                    (Unrecoverable — DON'T catch)
        │     ├── OutOfMemoryError
        │     ├── StackOverflowError
        │     ├── NoClassDefFoundError
        │     └── VirtualMachineError
        │
        └── java.lang.Exception                (Recoverable)
              ├── IOException                  ┐
              ├── SQLException                 │ Checked Exceptions
              ├── ClassNotFoundException       │ (Must handle or declare)
              ├── InterruptedException         ┘
              │
              └── RuntimeException             ┐
                    ├── NullPointerException    │
                    ├── ArrayIndexOutOfBounds   │ Unchecked Exceptions
                    ├── ClassCastException      │ (No obligation to handle)
                    ├── IllegalArgumentException│
                    ├── ArithmeticException     │
                    └── ConcurrentModification  ┘
```
## 3.3 Checked vs Unchecked Exceptions
| Aspect | Checked | Unchecked |
|---|---|---|
| **Superclass** | `Exception` (not `RuntimeException`) | `RuntimeException` |
| **Compile-time enforcement** | Must handle (`try-catch`) or declare (`throws`) | No obligation |
| **When** | External failures you can anticipate | Programming errors (bugs) |
| **Examples** | `IOException`, `SQLException` | `NullPointerException`, `ClassCastException` |
| **Philosophy** | Recoverable situations | Indicates defective code |
```java
// Checked — compiler forces you to handle it
public void readFile(String path) throws IOException {  // must declare
    BufferedReader reader = new BufferedReader(new FileReader(path));
    String line = reader.readLine();
}
// Unchecked — compiler doesn't enforce handling
public int divide(int a, int b) {
    return a / b;  // may throw ArithmeticException, but no compile-time requirement
}
```
> **Interview Insight:** The debate: Many modern frameworks (Spring) wrap checked exceptions into unchecked ones. Effective Java recommends using checked exceptions only for **recoverable** conditions where the caller can **take meaningful action**.
## 3.4 Error vs Exception
| Aspect | Error | Exception |
|---|---|---|
| **Recoverability** | Generally unrecoverable | Generally recoverable |
| **Caused by** | JVM/system problems | Application logic or external systems |
| **Handle?** | Almost never | Yes, handle where recovery is possible |
| **Examples** | `OutOfMemoryError`, `StackOverflowError` | `IOException`, `NullPointerException` |
| **Hierarchy** | Extends `Throwable` directly | Extends `Throwable` directly |
```java
// Catching Error is usually wrong, but sometimes necessary
try {
    riskyOperation();
} catch (OutOfMemoryError e) {
    // Dangerous: JVM may be in an unstable state
    // Only valid in specific cases (e.g., graceful shutdown, alert)
    logger.fatal("OOM detected, initiating shutdown", e);
    System.exit(1);
}
```
## 3.5 try, catch, finally
### Execution Flow
```java
try {
    // Code that may throw an exception
    int result = 10 / 0;
    System.out.println("This won't print");
} catch (ArithmeticException e) {
    // Handles the exception
    System.out.println("Caught: " + e.getMessage());
} finally {
    // ALWAYS executes (cleanup)
    System.out.println("Finally block");
}
// Output:
// Caught: / by zero
// Finally block
```
### When Does `finally` NOT Execute?
1. `System.exit()` is called.
2. JVM crashes or is killed (`kill -9`).
3. The thread running the try block is killed/interrupted in a fatal way.
4. Infinite loop or deadlock in try/catch block.
> **Classic Interview Question — What does this return?**
```java
public static int test() {
    try {
        return 1;
    } catch (Exception e) {
        return 2;
    } finally {
        return 3;  // ⚠️ This overrides the try/catch return!
    }
}
// Returns: 3
// finally's return overwrites the try block's return
// NEVER return from finally — it suppresses exceptions too
```
```java
// Even worse: finally swallows exceptions
public static int dangerous() {
    try {
        throw new RuntimeException("Error!");
    } finally {
        return 42; // exception is silently swallowed — VERY BAD
    }
}
// Returns 42, exception is lost!
```
## 3.6 Multiple Catch Blocks and Order
```java
try {
    riskyOperation();
} catch (FileNotFoundException e) {     // more specific FIRST
    System.out.println("File not found");
} catch (IOException e) {               // broader parent SECOND
    System.out.println("IO error");
} catch (Exception e) {                 // broadest LAST
    System.out.println("General error");
}
// Order: most specific → most general
// Reversing the order causes compile error (unreachable catch block)
```
### Multi-Catch (Java 7+)
```java
try {
    parse(input);
} catch (NumberFormatException | ParseException e) {
    // handles both types in single block
    // e is implicitly final — cannot be reassigned
    System.out.println("Parsing failed: " + e.getMessage());
}
```
**Rules:**
- Exceptions in multi-catch must not have a parent-child relationship.
- `catch (IOException | FileNotFoundException e)` → ❌ compile error (FileNotFoundException is a subtype of IOException).
## 3.7 throw vs throws
| Keyword | Purpose | Location | Usage |
|---|---|---|---|
| `throw` | **Throws** an exception object | Inside method body | `throw new IllegalArgumentException("Invalid");` |
| `throws` | **Declares** that method may throw | Method signature | `void read() throws IOException` |
```java
public class UserService {
    // throws: declares checked exceptions that may propagate
    public User findUser(int id) throws UserNotFoundException {
        User user = repository.findById(id);
        if (user == null) {
            // throw: creates and throws the exception
            throw new UserNotFoundException("User not found: " + id);
        }
        return user;
    }
}
```
> **Interview Note:** You can `throw` unchecked exceptions without `throws` in the signature. You **must** declare checked exceptions with `throws` (or catch them).
## 3.8 Custom Exceptions
### Best Practices for Custom Exceptions
```java
// Checked exception — caller can recover
public class InsufficientFundsException extends Exception {
    private final double currentBalance;
    private final double requestedAmount;
    public InsufficientFundsException(String message, double currentBalance, double requestedAmount) {
        super(message);
        this.currentBalance = currentBalance;
        this.requestedAmount = requestedAmount;
    }
    public double getDeficit() {
        return requestedAmount - currentBalance;
    }
}
// Unchecked exception — programming error or unrecoverable
public class InvalidOrderStateException extends RuntimeException {
    private final String orderId;
    private final String currentState;
    public InvalidOrderStateException(String orderId, String currentState, String attemptedAction) {
        super("Cannot " + attemptedAction + " order " + orderId + " in state " + currentState);
        this.orderId = orderId;
        this.currentState = currentState;
    }
}
```
**Guidelines:**
- Extend `Exception` if the caller can and should **recover**.
- Extend `RuntimeException` if it represents a **bug** or **unrecoverable** situation.
- Include **contextual information** (IDs, states, amounts) — not just a message.
- Provide constructors that accept a `Throwable cause` for chaining.
- Follow naming convention: `XxxException`.
## 3.9 try-with-resources
Automatically closes resources that implement `AutoCloseable`.
### Before (Java 6 style)
```java
BufferedReader reader = null;
try {
    reader = new BufferedReader(new FileReader("data.txt"));
    String line = reader.readLine();
} catch (IOException e) {
    e.printStackTrace();
} finally {
    if (reader != null) {
        try {
            reader.close();  // close can itself throw!
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
```
### After (Java 7+ try-with-resources)
```java
try (BufferedReader reader = new BufferedReader(new FileReader("data.txt"))) {
    String line = reader.readLine();
} catch (IOException e) {
    e.printStackTrace();
}
// reader.close() called automatically, even if an exception occurs
```
### Multiple Resources
```java
try (
    Connection conn = dataSource.getConnection();
    PreparedStatement stmt = conn.prepareStatement(sql);
    ResultSet rs = stmt.executeQuery()
) {
    while (rs.next()) {
        // process results
    }
}
// All three are closed in REVERSE order: rs → stmt → conn
```
### Suppressed Exceptions
```java
class MyResource implements AutoCloseable {
    @Override
    public void close() throws Exception {
        throw new Exception("Close failed");
    }
}
try (MyResource res = new MyResource()) {
    throw new Exception("Primary exception");
}
// Primary exception is thrown
// Close exception is "suppressed" — accessible via getSuppressed()
// catch (Exception e) {
//     e.getSuppressed(); // returns array containing "Close failed" exception
// }
```
## 3.10 Exception Propagation
Exceptions propagate **up the call stack** until caught or until they reach `main()` (causing program termination).
```
main()
  └── processOrder()
        └── validatePayment()
              └── chargeCard()   ← exception thrown here
Propagation path: chargeCard → validatePayment → processOrder → main → JVM (crash)
```
```java
void chargeCard(String cardNumber, double amount) throws PaymentException {
    // throws PaymentException — propagates to caller
    throw new PaymentException("Card declined");
}
void validatePayment(Order order) throws PaymentException {
    // doesn't catch — exception propagates further
    chargeCard(order.getCard(), order.getTotal());
}
void processOrder(Order order) {
    try {
        validatePayment(order);
    } catch (PaymentException e) {
        // finally caught here — fallback logic
        order.setStatus(OrderStatus.PAYMENT_FAILED);
        notifyUser(order, e.getMessage());
    }
}
```
> **Rule for checked exceptions:** Each method in the call chain must either `catch` the exception or declare it with `throws`. Unchecked exceptions propagate freely without `throws`.
---
# Part 4: Advanced Exception Handling
## 4.1 Best Practices — DOs and DON'Ts
### DO
```java
// 1. Catch specific exceptions
try {
    parse(json);
} catch (JsonParseException e) {     // ✅ specific
    handleParseError(e);
}
// 2. Include context in exception messages
throw new OrderNotFoundException(
    "Order not found: orderId=" + orderId + ", userId=" + userId
);
// 3. Use exception chaining (preserve root cause)
try {
    repository.save(entity);
} catch (SQLException e) {
    throw new DataAccessException("Failed to save entity: " + entity.getId(), e); // wraps cause
}
// 4. Clean up resources in finally or try-with-resources
// 5. Log at the point where you handle, not where you rethrow
// 6. Fail fast — validate inputs at method entry
public void transfer(Account from, Account to, double amount) {
    Objects.requireNonNull(from, "Source account must not be null");
    Objects.requireNonNull(to, "Target account must not be null");
    if (amount <= 0) throw new IllegalArgumentException("Amount must be positive: " + amount);
    // proceed...
}
```
### DON'T
```java
// 1. ❌ Don't catch Exception/Throwable broadly
try { ... }
catch (Exception e) { /* swallows everything, including NPE, which is a bug */ }
// 2. ❌ Don't swallow exceptions silently
try { ... }
catch (IOException e) { /* empty catch block — the worst antipattern */ }
// 3. ❌ Don't use exceptions for flow control
try {
    int i = 0;
    while (true) {
        array[i++] = process();  // uses ArrayIndexOutOfBoundsException to exit loop
    }
} catch (ArrayIndexOutOfBoundsException e) { /* loop done */ }
// ✅ Use: for (int i = 0; i < array.length; i++) { ... }
// 4. ❌ Don't log AND rethrow (leads to duplicate logs)
try { ... }
catch (IOException e) {
    logger.error("Error", e);  // logged here
    throw e;                    // logged again by the caller's handler → duplicate
}
// 5. ❌ Don't return from finally block (swallows exceptions)
// 6. ❌ Don't throw exceptions in finally block (overrides original exception)
```
## 4.2 When NOT to Catch Exceptions
- **NullPointerException** — fix the code (add null checks, use Optional).
- **ClassCastException** — fix the type logic (use generics, `instanceof`).
- **StackOverflowError** — fix the recursion (add base case, use iteration).
- **OutOfMemoryError** — fix memory usage, increase heap, find leaks.
- Any exception that indicates a **bug** — fix the bug, don't mask it.
> **Principle:** Catch exceptions when you can **meaningfully recover**. Let them propagate when the caller is better positioned to handle them.
## 4.3 Logging vs Rethrowing
```java
// Pattern 1: Handle and log (terminal handler)
public Response handleRequest(Request req) {
    try {
        return processRequest(req);
    } catch (BusinessException e) {
        logger.warn("Business rule violation: {}", e.getMessage());
        return Response.badRequest(e.getMessage());
    } catch (Exception e) {
        logger.error("Unexpected error processing request: {}", req.getId(), e);
        return Response.internalError();
    }
}
// Pattern 2: Wrap and rethrow (intermediate layer — DON'T log here)
public User getUser(int id) {
    try {
        return jdbcTemplate.queryForObject(sql, mapper, id);
    } catch (EmptyResultDataAccessException e) {
        throw new UserNotFoundException("User not found: " + id, e); // wrap, don't log
    }
}
// Pattern 3: Translate exception (layer boundary)
// Repository throws SQLException → Service throws DataAccessException → Controller handles
```
**Rule of thumb:** Log at the **outermost handler** (controller, message listener, scheduler). Wrap and rethrow at **layer boundaries**.
## 4.4 Fail-Fast vs Fail-Safe
| Approach | Behavior | Example |
|---|---|---|
| **Fail-Fast** | Immediately reports errors at the earliest point | `Objects.requireNonNull()`, `HashMap` iterator |
| **Fail-Safe** | Tolerates errors, continues operating | `ConcurrentHashMap` iterator, default values |
```java
// Fail-fast: validate early, crash loud
public void processOrder(Order order) {
    Objects.requireNonNull(order, "Order must not be null");
    if (order.getItems().isEmpty()) {
        throw new IllegalArgumentException("Order must have at least one item");
    }
    // proceed with validated data
}
// Fail-fast iterators (ArrayList, HashMap)
List<String> list = new ArrayList<>(List.of("a", "b", "c"));
for (String s : list) {
    list.remove(s);  // 💥 ConcurrentModificationException (fail-fast)
}
// Fail-safe iterators (ConcurrentHashMap, CopyOnWriteArrayList)
ConcurrentHashMap<String, Integer> map = new ConcurrentHashMap<>();
map.put("a", 1); map.put("b", 2);
for (Map.Entry<String, Integer> entry : map.entrySet()) {
    map.remove(entry.getKey());  // ✅ no exception (fail-safe, works on a snapshot/segment)
}
```
## 4.5 Real-World Exception Design
### Service Layer Pattern (Spring Boot)
```java
// Custom exception hierarchy for an e-commerce system
public abstract class ApplicationException extends RuntimeException {
    private final String errorCode;
    protected ApplicationException(String errorCode, String message) {
        super(message);
        this.errorCode = errorCode;
    }
    protected ApplicationException(String errorCode, String message, Throwable cause) {
        super(message, cause);
        this.errorCode = errorCode;
    }
    public String getErrorCode() { return errorCode; }
}
public class ResourceNotFoundException extends ApplicationException {
    public ResourceNotFoundException(String resource, Object id) {
        super("NOT_FOUND", resource + " not found with id: " + id);
    }
}
public class BusinessRuleException extends ApplicationException {
    public BusinessRuleException(String message) {
        super("BUSINESS_RULE_VIOLATION", message);
    }
}
```
### Global Exception Handler (Spring Boot)
```java
@RestControllerAdvice
public class GlobalExceptionHandler {
    @ExceptionHandler(ResourceNotFoundException.class)
    public ResponseEntity<ErrorResponse> handleNotFound(ResourceNotFoundException ex) {
        return ResponseEntity.status(404)
            .body(new ErrorResponse(ex.getErrorCode(), ex.getMessage()));
    }
    @ExceptionHandler(BusinessRuleException.class)
    public ResponseEntity<ErrorResponse> handleBusinessRule(BusinessRuleException ex) {
        return ResponseEntity.status(422)
            .body(new ErrorResponse(ex.getErrorCode(), ex.getMessage()));
    }
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorResponse> handleUnexpected(Exception ex) {
        log.error("Unexpected error", ex);
        return ResponseEntity.status(500)
            .body(new ErrorResponse("INTERNAL_ERROR", "An unexpected error occurred"));
    }
}
```
---
# Part 5: Modern Java Features (14+ / 17+ / 21)
## 5.1 var — Local Variable Type Inference
**Introduced:** Java 10 | **Why:** Reduce boilerplate for local variables where the type is obvious.
### Before vs After
```java
// Before (Java 9 and earlier)
Map<String, List<Employee>> departmentMap = new HashMap<String, List<Employee>>();
BufferedReader reader = new BufferedReader(new FileReader("data.txt"));
// After (Java 10+)
var departmentMap = new HashMap<String, List<Employee>>();
var reader = new BufferedReader(new FileReader("data.txt"));
```
### Where `var` CAN Be Used
```java
var list = new ArrayList<String>();       // local variable with initializer
var stream = list.stream();               // inferred from RHS
var entry = Map.entry("key", "value");    // inferred as Map.Entry<String, String>
for (var item : list) { }                 // enhanced for loop
for (var i = 0; i < 10; i++) { }         // traditional for loop
try (var conn = dataSource.getConnection()) { }  // try-with-resources
```
### Where `var` CANNOT Be Used
```java
var x;                    // ❌ no initializer → can't infer type
var x = null;             // ❌ null has no type
var x = {1, 2, 3};       // ❌ array initializer needs explicit type
var lambda = (x) -> x;   // ❌ lambda needs target type
var methodRef = System.out::println;  // ❌ method reference needs target type
// ❌ Not allowed for:
class Foo {
    var field = 10;              // fields
    var method(var param) { }    // parameters and return types
}
```
### Readability Trade-Offs
```java
// ✅ Good — type is obvious from RHS
var users = new ArrayList<User>();
var response = httpClient.send(request, BodyHandlers.ofString());
// ❌ Bad — type is unclear
var result = service.process(data);  // What type is result?
var x = calculate();                  // What does calculate return?
// Rule of thumb: use var when the type is visible on the same line
```
> **Interview Insight:** `var` is NOT a keyword — it's a **reserved type name**. You can still have a variable named `var` (though you shouldn't). `var` is compile-time sugar — no runtime impact. The bytecode is identical.
## 5.2 Records
**Introduced:** Java 14 (preview), Java 16 (stable) | **Why:** Eliminate boilerplate for simple data carrier classes.
### Before (Traditional POJO)
```java
public final class Point {
    private final int x;
    private final int y;
    public Point(int x, int y) {
        this.x = x;
        this.y = y;
    }
    public int x() { return x; }
    public int y() { return y; }
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Point)) return false;
        Point p = (Point) o;
        return x == p.x && y == p.y;
    }
    @Override
    public int hashCode() { return Objects.hash(x, y); }
    @Override
    public String toString() { return "Point[x=" + x + ", y=" + y + "]"; }
}
```
### After (Record)
```java
public record Point(int x, int y) { }
// That's it! Compiler generates:
// - private final fields (x, y)
// - canonical constructor
// - accessor methods: x(), y() (NOT getX(), getY())
// - equals(), hashCode(), toString()
```
### Customizing Records
```java
public record Employee(int id, String name, double salary) {
    // Compact constructor — validation without parameter assignment
    public Employee {
        if (salary < 0) throw new IllegalArgumentException("Salary cannot be negative");
        name = name.trim();  // can modify parameters before assignment
    }
    // Additional constructor
    public Employee(int id, String name) {
        this(id, name, 0.0);  // must delegate to canonical constructor
    }
    // Additional methods
    public double annualSalary() {
        return salary * 12;
    }
    // Can implement interfaces
    // public record Employee(int id, String name) implements Serializable { }
}
```
### What Records CANNOT Do
- **Cannot extend** a class (implicitly extend `java.lang.Record`).
- **Cannot be extended** (implicitly `final`).
- **Cannot have mutable instance fields** (all components are `final`).
- **Cannot declare additional instance fields** beyond the record components.
### Records vs POJOs vs Lombok
| Feature | POJO | Lombok `@Data` | Record |
|---|---|---|---|
| Boilerplate | Extensive | Minimal (annotation) | Minimal (built-in) |
| Immutability | Manual | `@Value` for immutable | Built-in (always immutable) |
| Inheritance | ✅ | ✅ | ❌ (final) |
| Extra fields | ✅ | ✅ | ❌ (only components) |
| Serialization | Manual | Manual | Built-in friendly |
| No dependency | ✅ | ❌ (Lombok lib) | ✅ |
> **When to use Records:** DTOs, API responses/requests, value objects, configuration, event payloads — any class whose sole purpose is carrying data.
## 5.3 Sealed Classes
**Introduced:** Java 15 (preview), Java 17 (stable) | **Why:** Control which classes can extend/implement a type — enabling exhaustive pattern matching.
### The Problem
```java
// Without sealed: anyone anywhere can extend Shape
public abstract class Shape { }
// In another module, someone does:
public class WeirdShape extends Shape { } // ← No way to prevent this
```
### The Solution
```java
public sealed class Shape
    permits Circle, Rectangle, Triangle {
    // Only Circle, Rectangle, Triangle can extend Shape
}
public final class Circle extends Shape {       // final: no further extension
    private final double radius;
    public Circle(double radius) { this.radius = radius; }
}
public sealed class Rectangle extends Shape     // sealed: controlled further extension
    permits Square {
    private final double width, height;
    public Rectangle(double width, double height) { this.width = width; this.height = height; }
}
public final class Square extends Rectangle {   // final: no further extension
    public Square(double side) { super(side, side); }
}
public non-sealed class Triangle extends Shape { // non-sealed: open for extension
    // anyone can extend Triangle
}
```
**Subclass modifier rules:**
- `final` — no further subclassing.
- `sealed` — restricted further subclassing.
- `non-sealed` — open for any subclassing (escape hatch).
### Sealed Interfaces
```java
public sealed interface Result<T>
    permits Success, Failure {
}
public record Success<T>(T value) implements Result<T> { }
public record Failure<T>(String error) implements Result<T> { }
```
### Exhaustive Switch with Sealed Types
```java
// Compiler knows ALL subtypes → can verify exhaustiveness (no default needed)
double area(Shape shape) {
    return switch (shape) {
        case Circle c    -> Math.PI * c.radius() * c.radius();
        case Rectangle r -> r.width() * r.height();
        case Triangle t  -> t.base() * t.height() / 2;
        // No default needed — compiler knows these are ALL possibilities
    };
}
```
> **Interview Insight:** Sealed classes enable the compiler to enforce **exhaustive pattern matching** — if you add a new permitted subclass, all switch expressions over that type will fail to compile until updated. This is a powerful domain modeling tool.
## 5.4 Pattern Matching
### instanceof Pattern Matching (Java 16)
**Why:** Eliminate redundant casting after `instanceof` checks.
```java
// Before (Java 15 and earlier)
if (obj instanceof String) {
    String s = (String) obj;    // redundant cast
    System.out.println(s.length());
}
// After (Java 16+)
if (obj instanceof String s) {  // test + cast + binding in one step
    System.out.println(s.length());
}
// Works with logical operators
if (obj instanceof String s && s.length() > 5) {
    System.out.println("Long string: " + s);
}
// Pattern variable scope
if (!(obj instanceof String s)) {
    return;  // early return
}
// s is in scope here (flow-scoping)
System.out.println(s.toUpperCase());
```
### Pattern Matching in Switch (Java 21)
```java
// Combines instanceof patterns with switch expressions
static String format(Object obj) {
    return switch (obj) {
        case Integer i    -> "Integer: %d".formatted(i);
        case Long l       -> "Long: %d".formatted(l);
        case Double d     -> "Double: %.2f".formatted(d);
        case String s     -> "String: \"%s\"".formatted(s);
        case int[] arr    -> "Array of length %d".formatted(arr.length);
        case null         -> "null";
        default           -> "Unknown: " + obj.getClass().getSimpleName();
    };
}
```
### Guarded Patterns (Java 21)
```java
static String categorize(Object obj) {
    return switch (obj) {
        case Integer i when i < 0  -> "Negative integer";
        case Integer i when i == 0 -> "Zero";
        case Integer i             -> "Positive integer";
        case String s when s.isEmpty() -> "Empty string";
        case String s              -> "String: " + s;
        case null                  -> "null";
        default                    -> "Other";
    };
}
```
> **Interview Insight:** Pattern matching in switch is **order-sensitive** — more specific patterns must come before general ones. `case Integer i when i < 0` must come before `case Integer i`. The compiler enforces this.
### Record Patterns (Java 21)
```java
record Point(int x, int y) { }
record Line(Point start, Point end) { }
// Destructuring records in pattern matching
static String describe(Object obj) {
    return switch (obj) {
        case Point(int x, int y) -> "Point at (%d, %d)".formatted(x, y);
        case Line(Point(var x1, var y1), Point(var x2, var y2)) ->
            "Line from (%d,%d) to (%d,%d)".formatted(x1, y1, x2, y2);
        default -> "Unknown";
    };
}
```
## 5.5 Updated Switch Expressions
**Introduced:** Java 14 (stable) | **Why:** Make switch less error-prone (no fall-through) and usable as an expression that returns a value.
### Traditional vs Modern
```java
// Traditional (statement — no return value, fall-through risk)
String dayType;
switch (day) {
    case MONDAY:
    case TUESDAY:
    case WEDNESDAY:
    case THURSDAY:
    case FRIDAY:
        dayType = "Weekday";
        break;
    case SATURDAY:
    case SUNDAY:
        dayType = "Weekend";
        break;
    default:
        dayType = "Unknown";
        break;
}
// Modern (expression — returns value, no fall-through)
String dayType = switch (day) {
    case MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY -> "Weekday";
    case SATURDAY, SUNDAY -> "Weekend";
};
// No default needed if all enum values are covered
```
### Arrow Syntax (`->`)
```java
// Arrow form: no fall-through, no break needed
switch (status) {
    case ACTIVE   -> System.out.println("Active");
    case INACTIVE -> System.out.println("Inactive");
    case PENDING  -> {
        log("Pending state detected");
        System.out.println("Pending");
    }
}
```
### `yield` Keyword
```java
// yield returns a value from a block in a switch expression
int numLetters = switch (day) {
    case MONDAY, FRIDAY, SUNDAY -> 6;
    case TUESDAY -> 7;
    case WEDNESDAY -> 9;
    case THURSDAY, SATURDAY -> 8;
};
// When you need a block with logic:
String result = switch (code) {
    case 200 -> "OK";
    case 404 -> "Not Found";
    default -> {
        logger.warn("Unexpected code: {}", code);
        yield "Unknown (" + code + ")";  // yield returns the value
    }
};
```
### Key Differences Summary
| Aspect | Traditional Switch | Switch Expression |
|---|---|---|
| Fall-through | Yes (without break) | No (arrow syntax) |
| Returns value | No (statement) | Yes (expression) |
| Multiple labels | Stacked cases | `case A, B, C ->` |
| Block return | `break` | `yield` |
| Exhaustiveness | Not enforced | Enforced (for enums/sealed) |
## 5.6 Virtual Threads (Project Loom — Java 21)
**Why:** Traditional platform threads are expensive (1-2 MB stack each, limited by OS). Modern microservices handle thousands of concurrent I/O-bound requests. Virtual threads make concurrency **cheap and scalable**.
### The Problem
```
Traditional Threading Model:
┌──────────────────────────────────────────────────────────────┐
│  1 request = 1 platform thread = 1 OS thread                │
│  Platform thread stack: ~1 MB                                │
│  Max threads ≈ available RAM / stack size ≈ few thousand     │
│                                                              │
│  10,000 concurrent requests → 10 GB RAM just for stacks!    │
│  Most of that time: thread is BLOCKED waiting for I/O        │
└──────────────────────────────────────────────────────────────┘
```
### Platform Threads vs Virtual Threads
| Aspect | Platform Thread | Virtual Thread |
|---|---|---|
| **Managed by** | OS kernel | JVM |
| **Stack size** | ~1 MB (fixed) | Starts small, grows as needed (~few KB) |
| **Creation cost** | Expensive (kernel call) | Cheap (Java object) |
| **Max count** | Thousands | Millions |
| **Blocking I/O** | Blocks the OS thread | Unmounts from carrier thread; carrier freed for other work |
| **Best for** | CPU-bound tasks | I/O-bound tasks |
| **Scheduler** | OS scheduler | JVM `ForkJoinPool` (work-stealing) |
### Before vs After
```java
// Before: Platform threads with thread pool (limited scalability)
ExecutorService executor = Executors.newFixedThreadPool(200); // capped at 200
List<Future<String>> futures = new ArrayList<>();
for (int i = 0; i < 10_000; i++) {
    futures.add(executor.submit(() -> {
        // many tasks wait in queue because pool is limited to 200
        return callExternalApi();
    }));
}
// After: Virtual threads (Java 21) — one thread per task
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    List<Future<String>> futures = new ArrayList<>();
    for (int i = 0; i < 10_000; i++) {
        futures.add(executor.submit(() -> {
            return callExternalApi();  // each gets its own virtual thread
        }));
    }
}
// 10,000 concurrent virtual threads, each using only KBs of memory
```
### How Virtual Threads Work Internally
```
Virtual Threads (managed by JVM)          Carrier Threads (platform, ForkJoinPool)
┌──────────────────────────┐              ┌───────────────────────┐
│ VThread-1 (running)    ──┼──mounted──►  │ Carrier-1 (executing) │
│ VThread-2 (blocked/IO)   │              │                       │
│ VThread-3 (running)    ──┼──mounted──►  │ Carrier-2 (executing) │
│ VThread-4 (blocked/IO)   │              │                       │
│ VThread-5 (waiting)      │              │ Carrier-3 (idle)      │
│ ...                      │              │                       │
│ VThread-10000 (blocked)  │              │ Carrier-N (# of CPUs) │
└──────────────────────────┘              └───────────────────────┘
When VThread-1 hits a blocking I/O call:
  1. VThread-1 is UNMOUNTED from Carrier-1
  2. VThread-1's stack is saved to heap
  3. Carrier-1 picks up VThread-4 (or any other runnable virtual thread)
  4. When I/O completes, VThread-1 is re-scheduled on any available carrier
```
### Creating Virtual Threads
```java
// 1. Thread.startVirtualThread()
Thread vt = Thread.startVirtualThread(() -> {
    System.out.println("Running on: " + Thread.currentThread());
});
// 2. Thread.ofVirtual().start()
Thread vt = Thread.ofVirtual()
    .name("my-vthread")
    .start(() -> System.out.println("Hello from virtual thread"));
// 3. ExecutorService (recommended for production)
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    IntStream.range(0, 100_000).forEach(i ->
        executor.submit(() -> {
            Thread.sleep(Duration.ofSeconds(1));
            return "Task " + i;
        })
    );
}
// 4. Spring Boot 3.2+ (application.properties)
// spring.threads.virtual.enabled=true
// All request-handling threads become virtual threads automatically
```
### When to Use and When NOT to Use
| Use Virtual Threads | Don't Use Virtual Threads |
|---|---|
| HTTP client calls | CPU-intensive computation (sorting, hashing) |
| Database queries | `synchronized` blocks holding during I/O (pins carrier) |
| File I/O | Tasks requiring thread-local caching (pooled threads) |
| Message queue consumers | Very short-lived CPU tasks |
| Microservice-to-microservice calls | When you need thread pooling semantics |
> **Interview Insight — Pinning:** Virtual threads **pin** to their carrier thread inside `synchronized` blocks or native methods. While pinned, the carrier cannot serve other virtual threads. Use `ReentrantLock` instead of `synchronized` for virtual-thread-friendly code.
```java
// ❌ Pins virtual thread to carrier
synchronized (lock) {
    callExternalApi(); // carrier thread is blocked while waiting for I/O
}
// ✅ Virtual-thread friendly
private final ReentrantLock lock = new ReentrantLock();
lock.lock();
try {
    callExternalApi(); // virtual thread unmounts; carrier is freed
} finally {
    lock.unlock();
}
```
### Impact on Spring Boot / Microservices
```
Traditional (Tomcat with 200 threads):
  200 concurrent requests → thread pool exhausted
  201st request waits in queue
With Virtual Threads (Spring Boot 3.2+):
  10,000+ concurrent requests → each gets a virtual thread
  No pooling bottleneck for I/O-bound workloads
  Throughput scales with I/O parallelism, not thread count
```
Enabling in Spring Boot 3.2+:
```properties
# application.properties
spring.threads.virtual.enabled=true
```
This makes all Tomcat request-handler threads, `@Async` methods, and `@Scheduled` methods use virtual threads.
---
# 6. Quick Revision Cheat Sheet
### OOP Pillars
| Pillar | One-Liner | Key Mechanism |
|---|---|---|
| Encapsulation | Hide data, expose methods | `private` fields + `public` getters/setters |
| Inheritance | Reuse via "is-a" relationship | `extends` / `implements` |
| Polymorphism | Same call, different behavior | Overloading (compile) / Overriding (runtime) |
| Abstraction | Show what, hide how | Abstract classes, interfaces |
### Overloading vs Overriding
| Aspect | Overloading | Overriding |
|---|---|---|
| Binding | Compile-time | Runtime |
| Methods | Same name, different params | Same name, same params |
| Return type | Can differ | Same or covariant |
| Access | Can differ | Same or wider |
| `static` | Yes | No (hiding, not overriding) |
### Association Types
- **Association** — uses-a (Teacher uses Student)
- **Aggregation** — has-a, independent lifecycle (Department has Professor)
- **Composition** — has-a, dependent lifecycle (Car has Engine)
### Exception Handling
- **Checked**: `Exception` subclass (not `RuntimeException`) → must handle or declare
- **Unchecked**: `RuntimeException` subclass → no obligation
- **Error**: JVM/system level → don't catch
- `finally` always runs (except `System.exit()`, JVM crash)
- Never return from `finally` (swallows exceptions)
- Multi-catch: `catch (A | B e)` — types must not be related
- try-with-resources: auto-closes `AutoCloseable` in reverse order
### equals() / hashCode() Contract
- `equals` → `true` implies `hashCode` must be equal
- Equal `hashCode` does NOT imply `equals` → `true`
- Override both or neither
### Modern Java Quick Reference
| Feature | Version | One-Liner |
|---|---|---|
| `var` | 10 | Local type inference |
| Switch expressions | 14 | Switch returns values, arrow syntax |
| Records | 16 | Immutable data carriers |
| Sealed classes | 17 | Controlled inheritance |
| Pattern matching (`instanceof`) | 16 | Cast-free type checks |
| Pattern matching (`switch`) | 21 | Type patterns in switch |
| Virtual threads | 21 | Lightweight JVM-managed threads for I/O |
### SOLID in One Line Each
- **S**: One class, one responsibility
- **O**: Open for extension, closed for modification
- **L**: Subtypes must be substitutable for base types
- **I**: Don't force clients to depend on unused methods
- **D**: Depend on abstractions, not concretions
---
# 7. Interview Questions with Answers (20)
### Q1: What are the four pillars of OOP? Explain with a real example.
**Answer:** Using an e-commerce system:
- **Encapsulation**: `Order` class hides its `items` list and `total` calculation behind methods. External code can't set `total` directly.
- **Inheritance**: `ElectronicsProduct extends Product` — inherits `name`, `price`, adds `warrantyPeriod`.
- **Polymorphism**: `paymentMethod.process()` calls different implementations depending on whether the runtime object is `CreditCard`, `UPI`, or `Wallet`.
- **Abstraction**: The `PaymentGateway` interface defines `charge(amount)` — the caller doesn't know if it's Stripe, Razorpay, or PayPal underneath.
---
### Q2: What is the difference between method overloading and method overriding?
**Answer:**
- **Overloading**: Same method name, **different** parameters (compile-time, static binding). Used for convenience — `add(int, int)` vs `add(double, double)`.
- **Overriding**: Same method name + same parameters in a subclass (runtime, dynamic binding). Used for polymorphism — `Dog.sound()` overrides `Animal.sound()`.
Key: Overloading is resolved by the **compiler** based on reference type + arguments. Overriding is resolved by the **JVM** based on the actual object type at runtime.
---
### Q3: Can we override a static method in Java?
**Answer:** No. Static methods are **hidden**, not overridden. When a subclass defines a static method with the same signature, it's called **method hiding**. The method called depends on the **reference type** (compile-time), not the object type (runtime). This is because static methods use **static binding**.
```java
Parent p = new Child();
p.staticMethod(); // calls Parent.staticMethod() — resolved by reference type
```
---
### Q4: What is the diamond problem? How does Java solve it?
**Answer:** The diamond problem occurs when a class inherits from two classes that have a common ancestor, creating ambiguity about which method implementation to use.
Java **prevents** multiple class inheritance entirely. For interfaces with conflicting `default` methods, the implementing class **must** override the method and explicitly choose (`InterfaceA.super.method()`) or provide its own implementation.
---
### Q5: Abstract class vs Interface — when to use which?
**Answer:**
- **Interface** when: defining a **capability** that unrelated classes share (`Serializable`, `Comparable`, `Closeable`). Use when you need multiple inheritance of type.
- **Abstract class** when: related classes share **state and behavior** (`InputStream` family). Use when subclasses need common fields or constructor logic.
Since Java 8+ (default methods), the line has blurred. Prefer interfaces for new designs; use abstract classes when shared state is needed.
---
### Q6: Explain the `equals()` and `hashCode()` contract. What happens if you break it?
**Answer:** Contract: if `a.equals(b)` is `true`, then `a.hashCode()` must equal `b.hashCode()`.
If broken (e.g., override `equals()` but not `hashCode()`): objects that are "equal" may land in **different hash buckets** in `HashMap`/`HashSet`. This means `set.contains(equalObject)` returns `false`, `map.get(equalKey)` returns `null`, and duplicate "equal" entries appear in sets.
---
### Q7: What is the difference between Aggregation and Composition?
**Answer:**
- **Aggregation**: Weak "has-a" — the part can **exist independently**. Example: `Department` has `Professor`. If the department is deleted, professors still exist.
- **Composition**: Strong "has-a" — the part's lifecycle is **tied to the whole**. Example: `House` has `Room`. If the house is demolished, rooms are destroyed.
Implementation difference: In composition, the whole **creates** the part internally. In aggregation, the part is **passed in** from outside.
---
### Q8: Explain exception propagation in Java.
**Answer:** When an exception is thrown, the JVM searches the current method for a matching `catch` block. If not found, the method terminates and the exception propagates to the **caller**. This continues up the call stack. If no handler is found all the way to `main()`, the JVM prints the stack trace and terminates.
For checked exceptions, every method in the chain must either `catch` or declare with `throws`. Unchecked exceptions propagate freely.
---
### Q9: Can a `finally` block prevent an exception from being thrown?
**Answer:** Yes — if the `finally` block has a `return` statement, it **swallows** the exception silently. This is a dangerous antipattern.
```java
try {
    throw new RuntimeException("Error!");
} finally {
    return 42; // Exception is silently lost. Method returns 42.
}
```
Similarly, if `finally` throws a new exception, the original exception is lost (unless you use `addSuppressed()`). Never return from `finally` or throw from `finally`.
---
### Q10: What are suppressed exceptions in try-with-resources?
**Answer:** If both the `try` block and the `close()` method throw exceptions, the `try` block's exception is the primary one that's thrown. The `close()` exception is added as a **suppressed exception**, accessible via `Throwable.getSuppressed()`.
This prevents the close exception from hiding the original error — a major improvement over manual try/finally patterns.
---
### Q11: Checked vs Unchecked exceptions — which should custom exceptions extend?
**Answer:**
- Extend `Exception` (checked) when the caller can **meaningfully recover** — e.g., `InsufficientFundsException` (retry with lower amount, switch payment method).
- Extend `RuntimeException` (unchecked) when it represents a **programming error** or **unrecoverable** condition — e.g., `InvalidConfigurationException`, `IllegalArgumentException`.
Modern Java frameworks (Spring) lean toward unchecked exceptions to reduce boilerplate.
---
### Q12: What is the purpose of Records in Java? Can they replace all POJOs?
**Answer:** Records are immutable data carriers that auto-generate constructor, accessors, `equals()`, `hashCode()`, and `toString()`. They eliminate boilerplate for DTOs, value objects, and event payloads.
They **cannot** replace all POJOs because: they cannot have mutable fields, cannot be extended, cannot declare extra instance fields, and accessors are named `field()` not `getField()` (may conflict with frameworks expecting JavaBean conventions, though most modern frameworks support records).
---
### Q13: What are sealed classes and why are they useful?
**Answer:** Sealed classes restrict which classes can extend them using the `permits` keyword. This enables:
1. **Exhaustive pattern matching** — the compiler knows all subtypes, so `switch` doesn't need a `default`.
2. **Domain modeling** — express "a Shape is either a Circle, Rectangle, or Triangle" in the type system.
3. **Library design** — control the inheritance hierarchy while still allowing multiple implementations.
Combined with records and pattern matching, sealed classes enable algebraic data types in Java.
---
### Q14: What problem do virtual threads solve? When should you NOT use them?
**Answer:** Virtual threads solve **thread scalability** for I/O-bound workloads. Traditional platform threads cost ~1MB each, limiting servers to thousands. Virtual threads cost ~KBs, enabling millions.
**Don't use** for: CPU-bound tasks (no I/O blocking to optimize), code with `synchronized` blocks around I/O (causes pinning), tasks that rely on thread-local caching in thread pools, or very short CPU tasks where scheduling overhead matters.
---
### Q15: What is "pinning" in virtual threads?
**Answer:** Pinning occurs when a virtual thread cannot unmount from its carrier (platform) thread — specifically inside `synchronized` blocks or native method calls. While pinned, the carrier thread is blocked and cannot serve other virtual threads, negating the benefit.
**Fix:** Replace `synchronized` with `ReentrantLock` for code that performs I/O while holding a lock.
---
### Q16: Why is `final` not the same as immutability?
**Answer:** `final` prevents **reassignment** of a reference. It does NOT prevent **modification** of the object the reference points to.
```java
final List<String> list = new ArrayList<>();
list.add("item"); // ✅ contents can change
list = new ArrayList<>(); // ❌ reference can't change
```
For true immutability: use `final` class + `private final` fields + no setters + defensive copies + `List.of()` / `Collections.unmodifiableList()`.
---
### Q17: Explain covariant return types with an example.
**Answer:** A covariant return type allows an overriding method to return a **subtype** of the return type declared in the parent class.
```java
class Animal {
    Animal create() { return new Animal(); }
}
class Dog extends Animal {
    @Override
    Dog create() { return new Dog(); } // covariant: Dog is a subtype of Animal
}
```
This avoids unnecessary casting at the call site. It works because a `Dog` **is-a** `Animal`, so the Liskov Substitution Principle holds.
---
### Q18: What happens when you call a method on a `null` reference? What about static methods?
**Answer:**
- **Instance method on null**: `NullPointerException` at runtime.
- **Static method on null**: **Works!** Static methods are resolved at compile time based on the reference type, not the object.
```java
String s = null;
s.length();       // 💥 NullPointerException
Integer i = null;
Integer.valueOf(5); // ✅ works fine (static method)
i.valueOf(5);       // ✅ also works! Compiles to Integer.valueOf(5) — null is irrelevant
```
---
### Q19: Can constructors be inherited in Java?
**Answer:** No. Constructors are **never** inherited. Each class must define its own constructors. However, a subclass constructor **must** call a superclass constructor (explicitly with `super()` or implicitly — the compiler inserts `super()` if omitted).
If the parent has no no-arg constructor and the child doesn't explicitly call a parameterized `super(...)`, it's a **compile error**.
---
### Q20: What is the difference between `fail-fast` and `fail-safe` iterators?
**Answer:**
- **Fail-fast** (`ArrayList`, `HashMap`): Detects structural modification during iteration via a `modCount` check. Throws `ConcurrentModificationException` immediately. Not guaranteed in all cases (best-effort).
- **Fail-safe** (`ConcurrentHashMap`, `CopyOnWriteArrayList`): Iterates over a snapshot or uses segment-level locking. Never throws `ConcurrentModificationException`. May not reflect recent modifications.
Choose based on requirements: fail-fast for detecting bugs early; fail-safe for concurrent access patterns.
---