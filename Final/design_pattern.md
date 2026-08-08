

## Singleton Design Pattern (Interview Deep Dive)

#### Singleton is a creational design pattern that ensures a class has only one instance throughout the application's lifecycle while providing a global access point to it. It is commonly used for shared resources like configuration managers, loggers, caches, and thread pools where creating multiple instances would waste resources or lead to inconsistent behavior.



Similarly:

- Logger
- Database Connection Pool
- Configuration Manager
- Cache Manager
- Thread Pool
- Spring Bean (Singleton Scope)

### Characteristics

A Singleton class has:

- Private constructor
- Static instance
- Static method to access the instance
- Prevents external object creation

### Complete Eager Singleton

```
public class Logger {

    private static Logger instance = new Logger();

    private Logger(){}

    public static Logger getInstance() {
        return instance;
    }

}
```

Usage

```
Logger logger = Logger.getInstance();
```

Advantages

- Very simple
- Thread-safe
- No synchronization required

Disadvantage

Creates the object even if it is never used.



### Lazy Initialization

Object is created only when needed.

```
Application Starts

instance = null -> First Request -> create Object -> Reuse foreve
```

```
public class Logger {

    private static Logger instance;

    private Logger() {}

    public static Logger getInstance() {

        if(instance == null){
            instance = new Logger();
        }

        return instance;
    }
}
```



Problem : two threads execute simultaneously. in instance == null both return true and create two objects.



### Thread-Safe Singleton

```
public class Logger {

    private static Logger instance;

    private Logger(){}

    public static synchronized Logger getInstance(){

        if(instance == null){
            instance = new Logger();
        }

        return instance;
    }

}
```

Works correctly.

### Downside

Every call acquires a lock.

```
1000 Requests -> 1000 synchronized calls
```

Performance suffers.



## Double-Checked Locking (Best for Interviews)

```
public class Logger {

    private static volatile Logger instance;

    private Logger() {}

    public static Logger getInstance() {

        if (instance == null) {

            synchronized (Logger.class) {

                if (instance == null) {
                    instance = new Logger();
                }

            }

        }

        return instance;
    }
}
```

### Why Two `if` Checks?

First check:

```
Avoid entering synchronized block
```

Second check:

```
Another thread may already have created the object
```

---

# Why `volatile`?

Without `volatile`: Object creation can be reordered.

Normally

```
1 Allocate Memory -> 2 Initialize Object -> 3 Assign Reference
```

CPU optimization may reorder:

```
1 Allocate Memory -> 3 Assign Reference -> 2 Initialize Object
```

Another thread could receive a reference to a **partially initialized object**.

####  `volatile` prevents this reordering and guarantees visibility across threads.



### Bill Pugh Singleton (Recommended)

Uses the JVM's class-loading mechanism.

```
public class Logger {

    private Logger() {}

    private static class Holder {
        private static final Logger INSTANCE = new Logger();
    }

    public static Logger getInstance() {
        return Holder.INSTANCE;
    }
}
```

### Why it's good

-  Lazy initialization 
-  Thread-safe 
-  No synchronization overhead 
-  Clean implementation 

This is one of the best manual implementations.

---



### Enum Singleton (Safest)

Joshua Bloch recommends this approach.

```
public enum Logger {

    INSTANCE;

    public void log() {
        System.out.println("Logging");
    }

}
```

Usage

```
Logger.INSTANCE.log();
```

### Advantages

-  Thread-safe 
-  Prevents reflection attacks 
-  Safe from serialization issues 
-  Very concise



Overall

### Advantages

- Saves memory
- Centralized state
- Easy access
- Thread-safe implementations available
- Ideal for shared resources

---

### Disadvantages

- Introduces global state, increasing coupling.
- Makes unit testing harder because dependencies are hidden.
- Can create concurrency issues if mutable state is stored.
- Violates the Single Responsibility Principle if it accumulates too many responsibilities.
- Poor fit when multiple independent instances may eventually be needed.



### Which Singleton implementation is best?


| Implementation         | Thread Safe | Lazy | Recommended                     |
| ---------------------- | ----------- | ---- | ------------------------------- |
| Eager                  | ✅           | ❌    | Good for always-used objects    |
| Lazy                   | ❌           | ✅    | No                              |
| Synchronized           | ✅           | ✅    | Acceptable but slower           |
| Double-Checked Locking | ✅           | ✅    | Very Good                       |
| Bill Pugh              | ✅           | ✅    | Excellent                       |
| Enum                   | ✅           | N/A  | Best for many Java applications |


---

#### The Singleton pattern ensures that only one instance of a class exists and provides a global access point to it. It's useful for shared resources such as loggers, configuration managers, cache managers, and thread pools. The implementation typically uses a private constructor to prevent direct instantiation and a static method like `getInstance()` to return the shared object.

#### In Spring Boot, beans are singleton-scoped by default, with the Spring container ensuring that all injections receive the same instance. The main trade-offs are global state and tighter coupling, so Singleton should be reserved for resources that truly need a single shared instance.



---

---

## Factory Design Pattern

#### Factory Pattern is a creational design pattern that centralizes object creation. Instead of instantiating concrete classes directly, clients request objects from a factory, which decides which implementation to create. This reduces coupling, makes the code extensible, and aligns with the Open/Closed Principle.

The key idea is:

> **"Create objects without exposing the object creation logic to the client."**



# Problem Without Factory

```
class OrderService {

    public void checkout(String paymentType) {

        Payment payment;

        if (paymentType.equals("CARD")) {
            payment = new CreditCardPayment();
        } else if (paymentType.equals("UPI")) {
            payment = new UpiPayment();
        } else {
            payment = new WalletPayment();
        }

        payment.pay();
    }
}
```

Problems:

-  Huge `if-else` 
-  Tight coupling 
-  Hard to extend 
-  Violates Open/Closed Principle 

# Solution: Factory Pattern

Move object creation into a dedicated class.

```
Client -- Factory -- Creates Correct Object -- Returns Interface
```

### Step 1: Create Interface

```
interface Payment {

    void pay();

}
```

---

### Step 2: Implement Different Payments

```
class CreditCardPayment implements Payment {

    @Override
    public void pay() {
        System.out.println("Paid using Credit Card");
    }
}
```

```
class UpiPayment implements Payment {

    @Override
    public void pay() {
        System.out.println("Paid using UPI");
    }
}
```

```
class WalletPayment implements Payment {

    @Override
    public void pay() {
        System.out.println("Paid using Wallet");
    }
}
```

---

### Step 3: Create Factory

```
class PaymentFactory {

    public static Payment getPayment(String type) {

        switch (type.toUpperCase()) {

            case "CARD":
                return new CreditCardPayment();

            case "UPI":
                return new UpiPayment();

            case "WALLET":
                return new WalletPayment();

            default:
                throw new IllegalArgumentException("Invalid Payment Type");
        }
    }
}
```

---

### Step 4: Client

```
public class Main {

    public static void main(String[] args) {

        Payment payment = PaymentFactory.getPayment("UPI");

        payment.pay();

    }
}
```

Output

```
Paid using UPI
```

Notice the client doesn't know which concrete class is created.

---

### Class Diagram

```
                 Payment
                    ▲
      ┌─────────────┼─────────────┐
      │             │             │
CreditCard      UPIPayment    WalletPayment
      ▲             ▲             ▲
      └─────────────┼─────────────┘
                    │
            PaymentFactory
                    │
             getPayment(type)
                    │
                  Client
```



# Advantages

### 1. Loose Coupling

Without Factory

```
new CreditCardPayment();
```

Client depends on implementation.

With Factory

```
Payment payment = PaymentFactory.getPayment(type);
```

Client depends only on the interface.

---

### 2. Centralized Object Creation

All creation logic stays in one place.

```
PaymentFactory -- Create -- Initialize --> Return
```

If constructors change, only the factory changes.

---

### 3. Easy to Extend

Need Apple Pay?

Just add:

```
class ApplePayPayment implements Payment
```

and update the factory.

---

### 4. Easier Testing

Instead of creating concrete classes, tests can inject mock implementations.

---

## Drawback

The simple factory shown above still needs modification whenever a new payment type is added.

```
switch(type)
```

This means it still partially violates the Open/Closed Principle.



## Factory Pattern in Spring Boot

BeanFactory, FactoryBean, ObjectMapper



### When Should You Use Factory?

Use Factory when:

- Object creation is complex.
- The client should not know the concrete implementation.
- Multiple implementations exist.
- You want to reduce coupling.
- You expect new implementations in the future.

Avoid using it when:

- There is only one implementation.
- Object creation is trivial.
- Adding a factory would introduce unnecessary abstraction.



## Abstract Factory

#### Abstract Factory is a creational design pattern that provides an interface to create families of related or dependent objects without exposing their concrete implementations. It is useful when a system needs to work with multiple product families while ensuring compatibility among the created objects.

The **Abstract Factory Pattern** is a **Creational Design Pattern** that provides an interface for creating **families of related objects** without specifying their concrete classes.

> **Factory creates one object.**
>
> **Abstract Factory creates a family of related objects.**



## Builder Design Pattern

The **Builder Pattern** is a **Creational Design Pattern** used to **construct complex objects step by step**. It is especially useful when an object has **many optional parameters** or when creating the object requires multiple configuration steps.

> **"Builder separates object construction from its representation."**

---

> **The Builder Pattern is a creational design pattern that constructs complex objects step by step using a fluent API. It avoids telescoping constructors, improves readability, supports immutable objects, and allows optional parameters without requiring multiple overloaded constructors.**





# How Builder Works

```
Client -> Builder -> Set Properties -> build() -> Final Object
```



# Immutable Objects

Builder is commonly used with immutable classes.

```
private final String name;
```

No setters.

Once created,

```
Employee employee = builder.build();
```

the object cannot be modified.

Advantages:

-  Thread-safe 
-  Predictable 
-  Easier to cache 
-  Easier to reason about 

---

# Validation During Build

Builder can validate before creating the object.

```
public Employee build() {

    if (name == null) {
        throw new IllegalArgumentException("Name is mandatory");
    }

    if (salary < 0) {
        throw new IllegalArgumentException("Invalid salary");
    }

    return new Employee(this);
}
```

This ensures invalid objects are never created.



# Fluent Builder ⭐⭐⭐⭐⭐ (Most Common)

This is what almost every Java project uses.

```
Employee employee = Employee.builder()
        .name("John")
        .age(25)
        .salary(100000)
        .build();
```

Implementation

```
public class Employee {

    private String name;
    private int age;

    private Employee(Builder builder) {
        this.name = builder.name;
        this.age = builder.age;
    }

    public static class Builder {

        private String name;
        private int age;

        public Builder name(String name) {
            this.name = name;
            return this;
        }

        public Builder age(int age) {
            this.age = age;
            return this;
        }

        public Employee build() {
            return new Employee(this);
        }
    }
}
```

### Pros

-  Readable 
-  Immutable objects 
-  Optional parameters 
-  Most popular



# Advantages

### 1. Readable

Instead of

```
new Employee(1, "Pankaj", ...)
```

Use

```
Employee.builder()
        .name("Pankaj")
```

---

### 2. No Constructor Explosion :Avoids dozens of overloaded constructors.

### 3. Optional Parameters : Only specify what you need.

### 4. Immutable Objects : Very common in enterprise applications.

### 5. Validation :Validate before object creation.

### 6. Easy Maintenance

Adding a new field:

```
private String city;
```

doesn't require creating multiple new constructors.

---

# Disadvantages

-  More code if written manually (although Lombok removes much of this). 
-  Slightly more objects are created (the Builder plus the target object), though this overhead is negligible for most applications. 
-  Overkill for simple objects with only a few fields.



# When Should You Use Builder?

Use Builder when:

- A class has many constructor parameters.
- Many fields are optional.
- You want immutable objects.
- Object creation requires validation.
- Readability is important.

Avoid Builder when:

- The object has only two or three simple required fields.
- Construction is trivial.

---

# Common Interview Questions

### Why not use setters?

Setters allow partially initialized or mutable objects. Builders can enforce required fields and produce immutable, fully initialized instances.

---

### Does Builder follow SOLID?

Yes:

- **Single Responsibility Principle:** Construction logic is separated from the object itself.
- **Open/Closed Principle:** New optional fields can often be added without changing existing client code.

---

### Why is the constructor private?

To ensure objects are created only through the Builder, allowing validation and preserving immutability.



# 3-Minute Interview Answer

> "The Builder Pattern is a creational design pattern used to construct complex objects step by step. It solves the problem of telescoping constructors and improves readability by using a fluent API. Instead of passing many constructor parameters, the client sets only the required properties through chained builder methods and finally calls `build()` to create the object. Builders are commonly used with immutable classes because they allow validation before object creation and prevent partially initialized objects. In Java, Lombok's `@Builder` annotation, `HttpRequest.newBuilder()`, Spring Security's `User.withUsername()`, and many AWS SDK request classes are examples of this pattern. I typically choose the Builder Pattern when an object has many optional fields or when readability and immutability are important."



## Strategy Design Pattern

The **Strategy Pattern** is a **Behavioral Design Pattern** that defines a **family of algorithms**, encapsulates each algorithm in a separate class, and makes them interchangeable at runtime.

> **"The Strategy Pattern lets you change an object's behavior without modifying its code."**

It is one of the **most frequently asked design patterns** in Java and Spring Boot interviews because it demonstrates **polymorphism**, **composition over inheritance**, and the **Open/Closed Principle**.

---

# Interview Definition (1-minute answer)

> "The Strategy Pattern defines a family of algorithms, encapsulates each algorithm into a separate class, and allows the client to choose the appropriate algorithm at runtime. It removes large if-else or switch statements, promotes loose coupling, and makes the system easy to extend."



# How Strategy Works

```
Client -- PaymentService -- PaymentService -- Concrete Strategy

The client decides which strategy to use.

Step 1: Strategy Interface

interface PaymentStrategy {
    void pay(double amount);
}

Step 2: Concrete Strategies

Credit Card

class CreditCardPayment implements PaymentStrategy {

    @Override
    public void pay(double amount) {
        System.out.println("Paid " + amount + " using Credit Card");
    }
}

UPI
class UpiPayment implements PaymentStrategy {

    @Override
    public void pay(double amount) {
        System.out.println("Paid " + amount + " using UPI");
    }
}

Wallet
class WalletPayment implements PaymentStrategy {

    @Override
    public void pay(double amount) {
        System.out.println("Paid " + amount + " using Wallet");
    }
}

Step 3: Context

class PaymentService {

    private final PaymentStrategy paymentStrategy;

    PaymentService(PaymentStrategy paymentStrategy) {
        this.paymentStrategy = paymentStrategy;
    }

    public void checkout(double amount) {
        paymentStrategy.pay(amount);
    }
}

Notice:

PaymentService doesn't know which payment method is used.

It only knows the interface.

Step 4: Client

public class Main {

    public static void main(String[] args) {

        PaymentStrategy strategy = new UpiPayment();

        PaymentService service = new PaymentService(strategy);

        service.checkout(1000);

    }
}

Output

Paid 1000 using UPI
```



# Advantages

## 1. Removes if-else

2. Open/Closed Principle

## 3. Easy Testing

Mock the interface.

```

```

```
PaymentStrategy mock =
    Mockito.mock(PaymentStrategy.class);
```

Testing becomes simple.

---

## 4. Loose Coupling

Context depends on the interface.

Not on implementations.

---

## 5. Runtime Selection

User chooses:

```
UPI - Card - Wallet
```

No code changes required.

---

# Disadvantages

-  More classes. 
-  Slight increase in complexity. 
-  Overkill if only one algorithm exists. 
-  Client (or factory) must decide which strategy to use.



# Where Strategy is Used in Java

### Comparator

```
Collections.sort(list, comparator);
```

Different comparators implement different sorting strategies.

---

### ExecutorService : Different thread pool implementations execute tasks differently.

---

### Compression

```
ZIP, RAR, GZIP
```

Each compression algorithm is a strategy.

---

# Where Strategy is Used in Spring

-  Payment gateway selection 
-  Authentication providers 
-  Message converters 
-  Cache implementations 
-  Retry policies 
-  Validation strategies 
-  Notification services 

---

# When Should You Use Strategy?

Use Strategy when:

-  Multiple algorithms solve the same problem. 
-  Algorithms should be interchangeable. 
-  You want to eliminate `if-else` chains. 
-  New behaviors are expected in the future. 
-  Runtime selection is required. 

Avoid Strategy when:

-  Only one algorithm exists. 
-  Algorithms are unlikely to change. 
-  Simplicity is more important than flexibility. 

---

# Common Interview Questions

### Why not use `if-else`?

Because `if-else` makes the code difficult to maintain and extend. Strategy encapsulates each algorithm in its own class, allowing new behaviors without modifying existing logic.

---

### Does Strategy follow SOLID?

Yes:

- **Open/Closed Principle:** Add new strategies without changing existing code. 
- **Dependency Inversion Principle:** The context depends on the `PaymentStrategy` abstraction, not concrete implementations. 

---

### Why composition instead of inheritance?

With composition, behavior can be changed at runtime by providing a different strategy object. Inheritance fixes behavior at compile time.

---

# 3-Minute Interview Answer

> "The Strategy Pattern is a behavioral design pattern that encapsulates a family of algorithms into separate classes and makes them interchangeable at runtime. It is useful when an application supports multiple ways of performing the same task, such as different payment methods, notification channels, or shipping providers. Instead of using large if-else or switch statements, the context depends on a common strategy interface, while each implementation contains its own algorithm. This follows the Open/Closed Principle because new strategies can be added without modifying existing business logic. In Spring Boot, a common approach is to inject all strategy implementations into a `Map<String, PaymentStrategy>`, allowing the appropriate strategy to be selected dynamically. This pattern improves maintainability, testability, and extensibility and is widely used in enterprise applications."



# Chain of Responsibility Design Pattern (Interview Deep Dive)

The **Chain of Responsibility (CoR)** is a **Behavioral Design Pattern** that passes a request through a **chain of handlers**. Each handler decides whether to:

1. Handle the request.
2. Pass it to the next handler.
3. Stop the chain.

> **"Instead of one object handling a request, multiple objects get a chance to process it."**

This pattern is heavily used in **Spring Boot**, **Servlet Filters**, **Spring Security**, **API Gateways**, and **middleware pipelines**.



Chain of Responsibility is a behavioral design pattern where a request passes through a sequence of handlers. Each handler can process the request or forward it to the next handler. It decouples the sender from the receiver and makes request-processing pipelines easy to extend.



### Imagine an online banking API.

Every request must go through:

```
Authentication - Authorization - Rate Limiting - Logging - Validation - Business Logic
```

Without Chain of Responsibility:

```
public void transferMoney(Request request) {
    authenticate(request);
    authorize(request);
    validate(request);
    rateLimit(request);
    log(request);
    transfer(request);
}
```

Problems:

- One large method
- Tight coupling
- Difficult to add/remove steps
- Violates Open/Closed Principle

Tomorrow:

Need Fraud Detection?

Need IP Whitelisting?

Need Audit?

Need Encryption?

You keep modifying the same class.

# Solution

Separate each responsibility into a handler.

```
Request - Authentication Handler - Authorization Handler - Validation Handler - Logging Handler - Transfer Handler

Each handler focuses on one responsibility.
```

# Components

The pattern consists of:

```
Client - Handler - Concrete Handler 1 - Concrete Handler 2
```

# Step 1: Handler

```
abstract class Handler {

    protected Handler next;

    public Handler setNext(Handler next) {
        this.next = next;
        return next;
    }

    public abstract void handle(Request request);

}
```

Notice

Every handler knows only the next handler.

---

# Step 2: Request

```
class Request {

    boolean authenticated;

    boolean authorized;

}
```

---

# Step 3: Authentication Handler

```
class AuthenticationHandler extends Handler {

    @Override
    public void handle(Request request) {

        if (!request.authenticated) {

            System.out.println("Authentication Failed");

            return;
        }

        System.out.println("Authenticated");

        if (next != null) {
            next.handle(request);
        }
    }
}
```

---

# Step 4: Authorization Handler

```
class AuthorizationHandler extends Handler {

    @Override
    public void handle(Request request) {

        if (!request.authorized) {

            System.out.println("Authorization Failed");

            return;
        }

        System.out.println("Authorized");

        if (next != null) {
            next.handle(request);
        }
    }
}
```

---

# Step 5: Business Handler

```
class TransferHandler extends Handler {

    @Override
    public void handle(Request request) {

        System.out.println("Money Transferred");

    }
}
```

---

# Step 6: Client

```
public class Main {

    public static void main(String[] args) {

        Handler auth = new AuthenticationHandler();

        Handler authorization = new AuthorizationHandler();

        Handler transfer = new TransferHandler();

        auth.setNext(authorization)
                .setNext(transfer);

        Request request = new Request();

        request.authenticated = true;
        request.authorized = true;

        auth.handle(request);

    }
}
```

Output

```
Authenticated

Authorized

Money Transferred
```



# Spring Boot Implementation

Spring makes this elegant.

```
@Component
@Order(1)
class CustomerValidator implements InvoiceValidator {}
```

```
@Component
@Order(2)
class GstValidator implements InvoiceValidator {}
```

```
@Component
@Order(3)
class DuplicateValidator implements InvoiceValidator {}
```

Spring injects

```
List<InvoiceValidator> validators;
```

already sorted by `@Order`.

Adding a new validator is easy:

```
@Component
@Order(4)
class CreditLimitValidator implements InvoiceValidator {}
```

No changes to `InvoiceValidationService`.

# API Gateway Example

```
Incoming Request - Authentication -> Authorization -> Rate Limiter -> Request Transformation -> Routing -> Microservice
```

Every API Gateway works like this.

---

# Banking Example

Money Transfer

```
Transfer Request -> Authenticate User -> Check Account -> Check Balance -> AML/Fraud Check -> Daily Limit Check -> Transfer Money -> Audit Log
```

If balance is insufficient:

```
Check Balance - Reject
```

The chain stops before transfer.

---

# Advantages

## 1. Loose Coupling : The client knows only the first handler.

---

## 2. Single Responsibility : Each handler has one responsibility.

---

## 3. Open/Closed Principle :

#### **Need Fraud Detection?**

Create

```
class FraudHandler extends Handler
```

Insert into the chain.

No existing handler changes.

---

## 4. Flexible Order

Today

```
Authentication - Validation
```

Tomorrow

```
Validation - Authentication
```

Just reorder the chain.

---

## 5. Reusable Handlers : Authentication can be reused across APIs.

---

# Disadvantages

-  Debugging can be harder because requests flow through multiple handlers. 
-  The order of handlers is critical; an incorrect order can introduce bugs or security issues. 
-  If a handler forgets to call the next handler, processing stops unexpectedly. 
-  Very long chains may add latency.

# Where Chain of Responsibility is Used in Java

- Servlet Filters
- Spring Security Filter Chain
- Netty Channel Pipeline
- Java Logging Frameworks
- Validation Frameworks

---

# Where Chain is Used in Spring

- `FilterChain`
- `OncePerRequestFilter`
- Spring Security Filters
- `HandlerInterceptor`
- Spring Cloud Gateway Filters
- WebFlux `WebFilter`

# When Should You Use Chain of Responsibility?

Use it when:

- A request must pass through multiple processing steps.
- The sender should not know which component handles the request.
- Processing steps may be reordered, added, or removed.
- Each step has a single responsibility.

Avoid it when:

- Only one component handles the request.
- The processing sequence is fixed and unlikely to change.
- A simple method call is sufficient.

# Difference Between Strategy and Chain of Responsibility


| Strategy                          | Chain of Responsibility               |
| --------------------------------- | ------------------------------------- |
| One algorithm is chosen           | Multiple handlers execute in sequence |
| Context delegates to one strategy | Request flows through many handlers   |
| Runtime algorithm selection       | Pipeline/workflow processing          |
| Example: Google Maps travel mode  | Example: Invoice validation           |




The Chain of Responsibility Pattern is a behavioral design pattern that passes a request through a sequence of handlers. Each handler performs its responsibility and then decides whether to forward the request to the next handler or stop processing. This decouples the sender from individual processing steps and follows the Single Responsibility and Open/Closed Principles. A classic example is the Spring Security filter chain, where every HTTP request passes through authentication, authorization, CSRF, logging, and other filters before reaching the controller. The pattern is also widely used in API gateways, servlet filters, middleware pipelines, and validation frameworks. It makes the request-processing pipeline modular, reusable, and easy to extend by adding or reordering handlers without changing existing code.



The Chain of Responsibility pattern passes a request through a sequence of handlers, where each handler performs one responsibility and decides whether to continue processing or stop the chain. It is useful when multiple independent processing steps must be executed in order without tightly coupling them together. In Spring Boot, it's commonly used for security filters, servlet filters, interceptors, and validation pipelines. For example, in an invoice validation system, each validator implements a common `InvoiceValidator` interface. Spring injects them as an ordered list, and the validation service iterates through them. Depending on the requirement, the chain can either fail fast on the first error or continue and collect all validation errors. This design follows the Single Responsibility and Open/Closed principles because new handlers can be added without modifying existing ones."



## Observer Design Pattern

The **Observer Pattern** is a **Behavioral Design Pattern** that defines a **one-to-many dependency** between objects. When one object (called the **Subject** or **Publisher**) changes its state, all dependent objects (called **Observers** or **Subscribers**) are automatically notified.

> **"One object changes → Automatically notify all interested objects."**

It is the foundation of **Event-Driven Architecture**, messaging systems, and many Spring Framework features.

---

# Interview Definition (1-minute answer)

> "The Observer Pattern is a behavioral design pattern where a Subject maintains a list of Observers and automatically notifies them whenever its state changes. It promotes loose coupling between the publisher and subscribers and is widely used in event-driven systems."



# Problem Without Observer

Suppose you have an Order Service.

```
class OrderService {

    public void placeOrder(Order order) {

        saveOrder(order);

        sendEmail(order);

        sendSMS(order);

        updateInventory(order);

        generateInvoice(order);

        sendWhatsApp(order);

    }
}
```

Problems:

-  Tight coupling 
-  Huge method 
-  Hard to extend 
-  Every new feature requires modifying `OrderService` 

Tomorrow

Need Analytics?

Need Loyalty Points?

Need Kafka Event?

Need Audit Logging?

You keep modifying the same class.

This violates the **Open/Closed Principle**.

---

# Solution

Instead of calling everyone directly:

```
Order Placed

↓

Publish Event

↓

Observer 1

Observer 2

Observer 3

Observer N
```

OrderService knows nothing about subscribers.



# Components

The Observer Pattern consists of four parts:

```
            Subject
               │
    -------------------------
    │          │           │
 Observer1  Observer2  Observer3
```


| Component           | Responsibility                              |
| ------------------- | ------------------------------------------- |
| Subject (Publisher) | Maintains observers and sends notifications |
| Observer            | Interface for receiving updates             |
| Concrete Observer   | Implements the response to notifications    |
| Client              | Registers observers with the subject        |


# Step 1: Observer Interface

```
interface Observer {
    void update(String message);
}
```

---

# Step 2: Concrete Observer

```
class MobileUser implements Observer {

    private final String name;

    MobileUser(String name) {
        this.name = name;
    }

    @Override
    public void update(String message) {

        System.out.println(name + " received : " + message);

    }
}
```

---

# Step 3: Subject Interface

```
interface Subject {

    void subscribe(Observer observer);

    void unsubscribe(Observer observer);

    void notifyObservers();

}
```

---

# Step 4: Concrete Subject

```
class YouTubeChannel implements Subject {

    private final List<Observer> observers = new ArrayList<>();

    private String latestVideo;

    @Override
    public void subscribe(Observer observer) {
        observers.add(observer);
    }

    @Override
    public void unsubscribe(Observer observer) {
        observers.remove(observer);
    }

    @Override
    public void notifyObservers() {

        for (Observer observer : observers) {
            observer.update(latestVideo);
        }
    }

    public void uploadVideo(String title) {

        this.latestVideo = title;

        notifyObservers();

    }
}
```

---

# Step 5: Client

```
public class Main {

    public static void main(String[] args) {

        YouTubeChannel channel = new YouTubeChannel();

        Observer user1 = new MobileUser("Pankaj");
        Observer user2 = new MobileUser("Rahul");

        channel.subscribe(user1);
        channel.subscribe(user2);

        channel.uploadVideo("Spring Boot Tutorial");
    }
}
```

Output

```
Pankaj received : Spring Boot Tutorial

Rahul received : Spring Boot Tutorial
```



# Advantages

## 1. Loose Coupling

Publisher only knows the Observer interface.

---

## 2. Open/Closed Principle

Need WhatsApp notifications?

Create

```

```

```
class WhatsAppListener
```

No changes elsewhere.

---

## 3. Easy Extension

Tomorrow add

-  Audit 
-  Rewards 
-  Metrics 
-  Fraud Detection 

Simply register more observers.

---

## 4. Supports Event-Driven Architecture

Perfect for

-  Kafka 
-  RabbitMQ 
-  Spring Events 
-  Domain Events 

---

## 5. Reusability

Observers are independent.

Each performs one responsibility.

---

# Disadvantages

-  Too many observers can make execution harder to trace. 
-  Notification order is not always guaranteed unless explicitly controlled. 
-  Synchronous notifications can slow the publisher if an observer is slow. 
-  Failed observers must be handled carefully to avoid affecting others. 
-  Debugging event-driven flows is often more difficult than direct method calls.

# Observer vs Pub-Sub

They're related but not identical.


| Observer                                       | Publish-Subscribe             |
| ---------------------------------------------- | ----------------------------- |
| Direct reference between Subject and Observers | Broker/message bus in between |
| Usually same process                           | Often distributed systems     |
| In-memory                                      | Kafka, RabbitMQ, Pulsar       |
| Synchronous by default                         | Often asynchronous            |


# Where Observer is Used in Java

- `PropertyChangeListener`
- Swing/AWT event listeners
- JavaBeans events

---

# Where Observer is Used in Spring

- `ApplicationEventPublisher`
- `@EventListener`
- Spring Integration
- Spring Cloud Stream
- Domain Events
- Spring Security authentication events

# When Should You Use Observer?

Use it when:

- One event should trigger multiple independent actions.
- The publisher should not know who the subscribers are.
- You expect subscribers to change over time.
- You're building event-driven or extensible systems.

Avoid it when:

- Only one component needs the result.
- Immediate return values are required.
- The workflow is simple and doesn't benefit from decoupling.

### Does Observer follow SOLID?

Yes:

- **Open/Closed Principle:** Add new observers without changing the publisher.
- **Dependency Inversion Principle:** The publisher depends on the `Observer` abstraction rather than concrete listeners.
- **Single Responsibility Principle:** Each observer handles one concern (email, inventory, analytics, etc.).

---

### Is Kafka an Observer Pattern?

Conceptually, yes. Kafka implements a **publish-subscribe architecture**, which generalizes the Observer idea for distributed systems. Producers publish events, and multiple consumers independently react to them through a broker.

---

# 3-Minute Interview Answer

> "The Observer Pattern is a behavioral design pattern that establishes a one-to-many relationship between a subject and its observers. When the subject's state changes, all registered observers are automatically notified. It helps decouple the publisher from subscribers and is commonly used in event-driven systems. A classic example is order processing: when an order is placed, the system publishes an `OrderPlaced` event, and independent listeners send emails, update inventory, generate invoices, and record analytics. In Spring Boot, this is implemented using `ApplicationEventPublisher` and `@EventListener`. In distributed microservices, Kafka and RabbitMQ extend the same concept through publish-subscribe messaging. The Observer Pattern improves extensibility and maintainability because new listeners can be added without changing the publisher."

