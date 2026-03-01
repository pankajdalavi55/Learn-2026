# Spring Framework — Complete Learning Guide

> A comprehensive, interview-ready reference covering the Spring Framework from fundamentals to production-grade patterns.  
> **Scope:** Spring Core only (IoC, DI, AOP, MVC, Data Access, Security, Testing). Spring Boot and Microservices are covered in separate guides.

---

## Table of Contents

1. [Introduction & Architecture](#1-introduction--architecture)
2. [Inversion of Control (IoC) & Dependency Injection (DI)](#2-inversion-of-control-ioc--dependency-injection-di)
3. [Bean Lifecycle & Scopes](#3-bean-lifecycle--scopes)
4. [Configuration Approaches](#4-configuration-approaches)
5. [Aspect-Oriented Programming (AOP)](#5-aspect-oriented-programming-aop)
6. [Spring MVC & REST](#6-spring-mvc--rest)
7. [Data Access — JDBC, ORM & Transactions](#7-data-access--jdbc-orm--transactions)
8. [Spring Security Fundamentals](#8-spring-security-fundamentals)
9. [Testing with Spring](#9-testing-with-spring)
10. [Event System & Async Processing](#10-event-system--async-processing)
11. [Spring Expression Language (SpEL)](#11-spring-expression-language-spel)
12. [Resource Handling & Profiles](#12-resource-handling--profiles)
13. [Best Practices & Anti-Patterns](#13-best-practices--anti-patterns)
14. [Interview Questions by Experience Level](#14-interview-questions-by-experience-level)

---

## 1. Introduction & Architecture

### 1.1 What is the Spring Framework?

Spring is an **open-source, lightweight, enterprise-grade** Java application framework that provides comprehensive infrastructure support. It was created by **Rod Johnson** in 2003 as a response to the complexity of J2EE (now Jakarta EE).

**Historical Context:**
In 2002, Rod Johnson published *"Expert One-on-One J2EE Design and Development"* criticizing the heavyweight nature of EJB 2.x. The book included 30,000 lines of infrastructure code that became the foundation of Spring. The framework officially launched in 2004 with version 1.0.

**Key Milestones:**
| Version | Year | Key Features |
|---|---|---|
| 1.0 | 2004 | IoC, AOP, JDBC abstraction |
| 2.0 | 2006 | XML namespaces, AspectJ integration |
| 2.5 | 2007 | Annotation-based configuration (`@Autowired`, `@Component`) |
| 3.0 | 2009 | Java-based config (`@Configuration`), SpEL, REST support |
| 4.0 | 2013 | Java 8 support, WebSocket, generics-based DI |
| 5.0 | 2017 | Reactive programming (WebFlux), Kotlin support, Java 9+ |
| 6.0 | 2022 | Jakarta EE 9+, native compilation (GraalVM AOT), virtual threads |

**Core Philosophy:**
- Favour **plain Java objects** (POJOs) over heavy enterprise components
- Promote **loose coupling** through IoC/DI
- Provide **declarative** programming via AOP
- Reduce **boilerplate** code
- Embrace **convention over configuration**
- Support **non-invasive** programming (your code doesn't depend on Spring APIs)

### 1.1.1 The POJO Programming Model

Spring advocates that business logic should reside in **Plain Old Java Objects** — classes with no special requirements:

```java
// POJO — No Spring dependency, fully testable
public class OrderProcessor {
    private final PaymentGateway paymentGateway;
    private final InventoryService inventory;

    public OrderProcessor(PaymentGateway paymentGateway, InventoryService inventory) {
        this.paymentGateway = paymentGateway;
        this.inventory = inventory;
    }

    public OrderResult process(Order order) {
        inventory.reserve(order.getItems());
        paymentGateway.charge(order.getPayment());
        return new OrderResult(Status.COMPLETED);
    }
}
```

**Benefits of POJO-centric design:**
- **Testability** — Unit test with simple `new` and mocks
- **Portability** — Not tied to any container
- **Simplicity** — No interfaces to implement, no base classes to extend
- **Maintainability** — Clear separation of concerns

### 1.2 Spring Modules (High-Level)

```
┌─────────────────────────────────────────────────────────────┐
│                      Spring Framework                       │
├────────────┬────────────┬──────────────┬────────────────────┤
│  Core      │  Web       │  Data Access │   Cross-Cutting    │
│  Container │  Layer     │  / ORM       │                    │
├────────────┼────────────┼──────────────┼────────────────────┤
│ spring-core│spring-web  │ spring-jdbc  │ spring-aop         │
│ spring-bean│spring-webmvc│spring-orm   │ spring-aspects     │
│ spring-ctx │spring-websocket│spring-tx │ spring-instrument  │
│ spring-spel│            │ spring-jms   │ spring-test        │
└────────────┴────────────┴──────────────┴────────────────────┘
```

### 1.3 Spring vs Jakarta EE vs Spring Boot

| Aspect | Spring Framework | Jakarta EE (Java EE) | Spring Boot |
|---|---|---|---|
| **Container** | Lightweight IoC container | Application server required | Embedded server (Tomcat/Jetty) |
| **Configuration** | XML / Java / Annotations | XML descriptors + Annotations | Auto-configuration |
| **Startup** | Manual wiring | Container-managed | Opinionated defaults |
| **Testing** | First-class support | CDI/Arquillian | `@SpringBootTest` |
| **Learning Curve** | Moderate | Steep | Low |

### 1.4 Setting Up a Spring Project (Non-Boot)

**Maven `pom.xml`:**
```xml
<properties>
    <spring.version>6.1.4</spring.version>
</properties>

<dependencies>
    <!-- Core -->
    <dependency>
        <groupId>org.springframework</groupId>
        <artifactId>spring-context</artifactId>
        <version>${spring.version}</version>
    </dependency>

    <!-- Web MVC -->
    <dependency>
        <groupId>org.springframework</groupId>
        <artifactId>spring-webmvc</artifactId>
        <version>${spring.version}</version>
    </dependency>

    <!-- JDBC -->
    <dependency>
        <groupId>org.springframework</groupId>
        <artifactId>spring-jdbc</artifactId>
        <version>${spring.version}</version>
    </dependency>
</dependencies>
```

---

## 2. Inversion of Control (IoC) & Dependency Injection (DI)

### 2.0 Theoretical Foundation

#### What is Inversion of Control?

IoC is a **design principle** (not a pattern) where the control of object creation and lifecycle is transferred from the application code to a framework or container. The "inversion" refers to reversing the traditional flow where objects directly instantiate their dependencies.

**Traditional Control Flow:**
```
Application Code → Creates Dependencies → Uses Dependencies
```

**Inverted Control Flow:**
```
Container Creates Dependencies → Injects into Application Code → Code Uses Dependencies
```

#### IoC vs DI — Clearing the Confusion

| Concept | Definition | Relationship |
|---|---|---|
| **IoC** | Design principle | The broader concept |
| **DI** | Implementation technique | One way to achieve IoC |
| **Service Locator** | Alternative technique | Another way to achieve IoC |

**IoC can be achieved through:**
1. **Dependency Injection** (Spring's approach) — Dependencies pushed to the object
2. **Service Locator** — Object pulls dependencies from a registry
3. **Factory Pattern** — Factory creates and returns dependencies
4. **Template Method** — Superclass controls algorithm, subclass provides specifics

> Spring uses DI because it keeps components unaware of the container and promotes testability.

#### Connection to SOLID Principles

DI directly supports **three** of the five SOLID principles:

| Principle | How DI Helps |
|---|---|
| **S**ingle Responsibility | Each class focuses on one job; infrastructure concerns (object creation) handled elsewhere |
| **O**pen/Closed | Add new implementations without modifying existing code (just wire a different bean) |
| **D**ependency Inversion | High-level modules depend on abstractions (interfaces), not concrete implementations |

**Dependency Inversion Principle (DIP) in Practice:**
```
Traditional:                          With DIP:
┌──────────────┐                     ┌──────────────┐
│ OrderService │                     │ OrderService │
│   depends on │                     │   depends on │
│              ▼                     │              ▼
│ JdbcOrderRepo│                     │ <<interface>>│
└──────────────┘                     │ OrderRepo    │
                                     └──────────────┘
                                            ▲
                                     ┌──────┴───────┐
                                     │JdbcOrderRepo │
                                     │MongoOrderRepo│
                                     └──────────────┘
```

#### The Hollywood Principle

DI follows the "Hollywood Principle": **"Don't call us, we'll call you."**

- Your objects don't ask for dependencies
- The container provides them when constructing your objects
- This leads to cleaner, more declarative code

### 2.1 The Problem DI Solves

**Without DI (Tight Coupling):**
```java
public class OrderService {
    // OrderService CREATES its own dependency — tightly coupled
    private final OrderRepository repo = new JdbcOrderRepository();

    public void placeOrder(Order order) {
        repo.save(order);
    }
}
```

**With DI (Loose Coupling):**
```java
public class OrderService {
    private final OrderRepository repo; // interface

    // Dependency INJECTED from outside
    public OrderService(OrderRepository repo) {
        this.repo = repo;
    }

    public void placeOrder(Order order) {
        repo.save(order);
    }
}
```

> **Key Insight:** The object does NOT create or locate its dependencies; an external entity (the Spring container) provides them. This is the *Inversion* in IoC.

### 2.2 The Spring IoC Container

The container is responsible for:
1. **Instantiating** beans
2. **Configuring** them (injecting dependencies)
3. **Managing** their lifecycle (init → use → destroy)

Two main container implementations:

| Container | Interface | Use Case |
|---|---|---|
| **BeanFactory** | `BeanFactory` | Lightweight, lazy init, low memory |
| **ApplicationContext** | `ApplicationContext` | Full-featured: events, i18n, AOP, eager init |

> In practice, always use `ApplicationContext`; `BeanFactory` is legacy.

```java
// Bootstrapping ApplicationContext
// 1. Annotation-based
ApplicationContext ctx = new AnnotationConfigApplicationContext(AppConfig.class);

// 2. XML-based
ApplicationContext ctx = new ClassPathXmlApplicationContext("applicationContext.xml");

// 3. Web-based
// Configured via web.xml ContextLoaderListener or WebApplicationInitializer
```

### 2.3 Types of Dependency Injection

#### Constructor Injection (Recommended)
```java
@Component
public class NotificationService {
    private final EmailSender emailSender;
    private final SmsSender smsSender;

    @Autowired // optional on single constructor since Spring 4.3
    public NotificationService(EmailSender emailSender, SmsSender smsSender) {
        this.emailSender = emailSender;
        this.smsSender = smsSender;
    }
}
```

**Why prefer constructor injection?**
- Ensures **immutability** (`final` fields)
- Makes dependencies **explicit** — visible in constructor signature
- Object is **never in an incomplete state**
- Easier to **unit test** (just call `new` with mocks)

#### Setter Injection
```java
@Component
public class ReportGenerator {
    private DataSource dataSource;

    @Autowired
    public void setDataSource(DataSource dataSource) {
        this.dataSource = dataSource;
    }
}
```

Use setter injection for **optional** dependencies or when you need to allow **reconfiguration** after construction.

#### Field Injection (Avoid in Production)
```java
@Component
public class UserService {
    @Autowired
    private UserRepository userRepository; // no setter, no constructor param
}
```

**Why avoid?**
- Cannot make fields `final`
- Hides dependencies
- Impossible to instantiate without reflection (hard to unit test)
- Violates single-responsibility detection (class can silently accumulate many injected fields)

### 2.4 `@Autowired` Resolution Rules

Spring resolves `@Autowired` in this order:

```
1. Match by TYPE  →  Only one bean of requested type?  ✅ Inject it.
2. Multiple candidates?
   a. @Primary bean present?   ✅ Use it.
   b. @Qualifier specified?    ✅ Match by qualifier name.
   c. Match by FIELD/PARAM NAME as bean name?  ✅ Use it.
   d. None matched?           ❌ Throw NoUniqueBeanDefinitionException.
```

**Example with `@Qualifier`:**
```java
@Component("mysqlRepo")
public class MySqlOrderRepository implements OrderRepository { }

@Component("mongoRepo")
public class MongoOrderRepository implements OrderRepository { }

@Component
public class OrderService {
    private final OrderRepository repo;

    public OrderService(@Qualifier("mongoRepo") OrderRepository repo) {
        this.repo = repo;
    }
}
```

### 2.5 `@Primary` vs `@Qualifier`

| Feature | `@Primary` | `@Qualifier` |
|---|---|---|
| Applied on | Bean definition | Injection point |
| Scope | Global default for a type | Per-injection override |
| Precedence | Lower (fallback) | Higher (explicit) |

```java
@Primary
@Component
public class PostgresRepo implements OrderRepository { }

@Component
public class MongoRepo implements OrderRepository { }

// Without @Qualifier → PostgresRepo wins (it's @Primary)
// With @Qualifier("mongoRepo") → MongoRepo wins
```

### 2.6 Injection of Collections

```java
@Component
public class NotificationRouter {
    private final List<NotificationSender> senders; // ALL beans of this type

    public NotificationRouter(List<NotificationSender> senders) {
        this.senders = senders; // [EmailSender, SmsSender, PushSender, …]
    }

    public void notifyAll(String message) {
        senders.forEach(s -> s.send(message));
    }
}
```

Control order with `@Order(1)`, `@Order(2)`, etc.

### 2.7 `@Lazy` Initialization

```java
@Component
@Lazy // Bean created only when first requested, NOT at context startup
public class HeavyReportEngine { }
```

Can also be applied at the injection point:
```java
public OrderService(@Lazy HeavyReportEngine engine) {
    this.engine = engine; // Spring injects a PROXY; actual bean created on first method call
}
```

### 2.8 `ObjectProvider` — Programmatic, Safe Lookup

```java
@Component
public class PaymentProcessor {
    private final ObjectProvider<FraudChecker> fraudCheckerProvider;

    public PaymentProcessor(ObjectProvider<FraudChecker> fraudCheckerProvider) {
        this.fraudCheckerProvider = fraudCheckerProvider;
    }

    public void process(Payment p) {
        // getIfAvailable — returns null if no bean, no exception
        FraudChecker fc = fraudCheckerProvider.getIfAvailable(NoOpFraudChecker::new);
        fc.check(p);
    }
}
```

---

## 3. Bean Lifecycle & Scopes

### 3.0 Why Bean Lifecycle Matters

Understanding bean lifecycle is crucial for:

1. **Resource Management** — Acquire resources (DB connections, file handles, thread pools) at the right time and release them properly
2. **Initialization Order** — Ensure dependencies are ready before a bean is used
3. **Framework Integration** — Hook into Spring's infrastructure (AOP proxies, transaction management)
4. **Debugging** — Diagnose issues like "bean not found" or "bean not fully initialized"
5. **Performance Optimization** — Defer expensive initialization until necessary

#### The Container's Responsibilities

Spring's IoC container acts as a **sophisticated object factory** with these responsibilities:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Spring IoC Container                         │
├─────────────────────────────────────────────────────────────────┤
│  1. Read Configuration    │  XML, annotations, Java config     │
│  2. Create Bean Graph     │  Resolve dependencies, detect cycles│
│  3. Instantiate Beans     │  Call constructors in correct order │
│  4. Inject Dependencies   │  Wire beans together               │
│  5. Apply Post-Processing │  AOP proxies, validation           │
│  6. Initialize            │  Call @PostConstruct, init methods │
│  7. Make Available        │  Beans ready for use               │
│  8. Destroy               │  Clean shutdown, release resources │
└─────────────────────────────────────────────────────────────────┘
```

#### Eager vs Lazy Initialization

| Mode | When Created | Default For | Trade-off |
|---|---|---|---|
| **Eager** | At container startup | Singletons | Slower startup, faster first access, fail-fast |
| **Lazy** | On first access | Prototypes | Faster startup, slower first access, late failure |

> **Production recommendation:** Prefer eager initialization. Fail-fast at startup catches configuration errors before handling real traffic.

### 3.1 Complete Bean Lifecycle

```
  Container starts
       │
       ▼
  ① Bean Definition Read (XML / @Component scan / @Bean methods)
       │
       ▼
  ② Instantiation (constructor call)
       │
       ▼
  ③ Populate Properties (DI — setter/field injection)
       │
       ▼
  ④ BeanNameAware.setBeanName()
       │
       ▼
  ⑤ BeanFactoryAware.setBeanFactory()
       │
       ▼
  ⑥ ApplicationContextAware.setApplicationContext()
       │
       ▼
  ⑦ BeanPostProcessor.postProcessBeforeInitialization()
       │
       ▼
  ⑧ @PostConstruct / InitializingBean.afterPropertiesSet() / init-method
       │
       ▼
  ⑨ BeanPostProcessor.postProcessAfterInitialization()
       │
       ▼
  ⑩ Bean is READY — available in container
       │
       ▼
  (Application runs…)
       │
       ▼
  ⑪ @PreDestroy / DisposableBean.destroy() / destroy-method
       │
       ▼
  Container shuts down
```

### 3.2 Initialization & Destruction Callbacks

```java
@Component
public class CacheManager {

    // ---------- METHOD 1: JSR-250 annotations (PREFERRED) ----------
    @PostConstruct
    public void init() {
        System.out.println("Cache warmed up");
    }

    @PreDestroy
    public void shutdown() {
        System.out.println("Cache cleared");
    }
}

// ---------- METHOD 2: Spring interfaces ----------
public class CacheManager implements InitializingBean, DisposableBean {
    @Override
    public void afterPropertiesSet() { /* init */ }

    @Override
    public void destroy() { /* cleanup */ }
}

// ---------- METHOD 3: @Bean attributes ----------
@Configuration
public class AppConfig {
    @Bean(initMethod = "init", destroyMethod = "shutdown")
    public CacheManager cacheManager() {
        return new CacheManager();
    }
}
```

**Execution Order (if all three defined):**
1. `@PostConstruct`
2. `InitializingBean.afterPropertiesSet()`
3. Custom `initMethod`

### 3.3 BeanPostProcessor (BPP)

A **hook** to modify or wrap beans **before** and **after** their initialization.

```java
@Component
public class TimingBeanPostProcessor implements BeanPostProcessor {

    @Override
    public Object postProcessBeforeInitialization(Object bean, String beanName) {
        // Called BEFORE @PostConstruct
        return bean;
    }

    @Override
    public Object postProcessAfterInitialization(Object bean, String beanName) {
        // Called AFTER @PostConstruct — great place to create proxies
        if (bean instanceof DataSource ds) {
            return new MonitoredDataSourceProxy(ds); // return wrapped bean
        }
        return bean;
    }
}
```

> **Spring AOP** and **`@Transactional`** work because of BPPs that create proxies in `postProcessAfterInitialization`.

### 3.4 BeanFactoryPostProcessor (BFPP)

Modifies **bean definitions** (metadata) before any bean is instantiated.

```java
@Component
public class CustomPropertyOverrider implements BeanFactoryPostProcessor {
    @Override
    public void postProcessBeanFactory(ConfigurableListableBeanFactory factory) {
        BeanDefinition bd = factory.getBeanDefinition("dataSource");
        bd.getPropertyValues().add("url", "jdbc:mysql://new-host/db");
    }
}
```

> `PropertySourcesPlaceholderConfigurer` (resolves `${…}` placeholders) is itself a BFPP.

### 3.5 Bean Scopes

| Scope | Annotation / XML | Behaviour |
|---|---|---|
| **singleton** (default) | `@Scope("singleton")` | One instance per Spring container |
| **prototype** | `@Scope("prototype")` | New instance every time requested |
| **request** | `@RequestScope` | One per HTTP request |
| **session** | `@SessionScope` | One per HTTP session |
| **application** | `@ApplicationScope` | One per ServletContext |
| **websocket** | `@Scope("websocket")` | One per WebSocket session |

```java
@Component
@Scope("prototype")
public class ShoppingCart { }
```

### 3.6 The Singleton-Prototype Injection Problem

```java
@Component // singleton by default
public class CheckoutService {
    @Autowired
    private ShoppingCart cart; // prototype

    // ⚠ BUG: same cart instance every time — prototype scope "lost"
}
```

**Solutions:**

**1. Method Injection with `@Lookup` (Preferred):**
```java
@Component
public abstract class CheckoutService {

    public void checkout() {
        ShoppingCart cart = getCart(); // fresh prototype each time
        cart.addItem(...);
    }

    @Lookup
    protected abstract ShoppingCart getCart(); // Spring overrides at runtime
}
```

**2. `ObjectProvider`:**
```java
@Component
public class CheckoutService {
    private final ObjectProvider<ShoppingCart> cartProvider;

    public CheckoutService(ObjectProvider<ShoppingCart> cartProvider) {
        this.cartProvider = cartProvider;
    }

    public void checkout() {
        ShoppingCart cart = cartProvider.getObject(); // new instance
    }
}
```

**3. Scoped Proxy:**
```java
@Component
@Scope(value = "prototype", proxyMode = ScopedProxyMode.TARGET_CLASS)
public class ShoppingCart { }
// Spring injects a CGLIB proxy; each method call delegates to a fresh instance
```

---

## 4. Configuration Approaches

### 4.1 Java-Based Configuration (`@Configuration`)

```java
@Configuration
@ComponentScan(basePackages = "com.example")
@PropertySource("classpath:application.properties")
public class AppConfig {

    @Value("${db.url}")
    private String dbUrl;

    @Bean
    public DataSource dataSource() {
        HikariDataSource ds = new HikariDataSource();
        ds.setJdbcUrl(dbUrl);
        ds.setMaximumPoolSize(10);
        return ds;
    }

    @Bean
    public PlatformTransactionManager txManager(DataSource ds) {
        return new DataSourceTransactionManager(ds);
    }
}
```

### 4.2 Why `@Configuration` Classes are Special

`@Configuration` classes are **CGLIB-proxied**. This ensures `@Bean` method calls between methods return the **same singleton instance**:

```java
@Configuration
public class AppConfig {

    @Bean
    public ServiceA serviceA() {
        return new ServiceA(commonDependency()); // ← calls @Bean method
    }

    @Bean
    public ServiceB serviceB() {
        return new ServiceB(commonDependency()); // ← same call
    }

    @Bean
    public CommonDependency commonDependency() {
        return new CommonDependency(); // created ONCE thanks to CGLIB proxy
    }
}
```

> If you use `@Component` instead of `@Configuration`, each call to `commonDependency()` creates a **new** instance (lite mode — no proxying).

### 4.3 XML Configuration (Legacy — Know for Interviews)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:context="http://www.springframework.org/schema/context"
       xsi:schemaLocation="...">

    <!-- Enable annotation processing -->
    <context:annotation-config/>

    <!-- Component scanning -->
    <context:component-scan base-package="com.example"/>

    <!-- Manual bean definition -->
    <bean id="orderService" class="com.example.OrderService">
        <constructor-arg ref="orderRepository"/>
    </bean>

    <bean id="orderRepository" class="com.example.JdbcOrderRepository">
        <property name="dataSource" ref="dataSource"/>
    </bean>
</beans>
```

### 4.4 Stereotype Annotations

| Annotation | Layer | Specialisation |
|---|---|---|
| `@Component` | Generic | Base annotation for any Spring-managed bean |
| `@Service` | Business | Marks service-layer logic (no extra behaviour) |
| `@Repository` | Persistence | Enables **exception translation** (DB exceptions → `DataAccessException`) |
| `@Controller` | Web | Marks MVC controller, returns views |
| `@RestController` | Web | `@Controller` + `@ResponseBody` (every method returns JSON/XML) |
| `@Configuration` | Config | Declares `@Bean` factory methods, CGLIB-proxied |

### 4.5 `@ComponentScan` — Controlling Discovery

```java
@Configuration
@ComponentScan(
    basePackages = "com.example",
    includeFilters = @ComponentScan.Filter(type = FilterType.ANNOTATION, classes = MyCustom.class),
    excludeFilters = @ComponentScan.Filter(type = FilterType.REGEX, pattern = ".*Test.*")
)
public class AppConfig { }
```

### 4.6 `@Conditional` — Conditional Bean Registration

```java
@Bean
@Conditional(OnLinuxCondition.class)
public FileWatcher linuxWatcher() {
    return new InotifyFileWatcher();
}

public class OnLinuxCondition implements Condition {
    @Override
    public boolean matches(ConditionContext ctx, AnnotatedTypeMetadata metadata) {
        return ctx.getEnvironment().getProperty("os.name").contains("Linux");
    }
}
```

> Spring Boot simplifies this with `@ConditionalOnProperty`, `@ConditionalOnClass`, `@ConditionalOnMissingBean`, etc.

### 4.7 `@Import` and Modular Configuration

```java
@Configuration
@Import({DatabaseConfig.class, SecurityConfig.class, CacheConfig.class})
public class RootConfig { }
```

For dynamic imports:
```java
public class MyImportSelector implements ImportSelector {
    @Override
    public String[] selectImports(AnnotationMetadata metadata) {
        return new String[] { "com.example.DynamicConfig" };
    }
}

@Import(MyImportSelector.class)
@Configuration
public class AppConfig { }
```

---

## 5. Aspect-Oriented Programming (AOP)

### 5.0 Theoretical Foundation of AOP

#### The Problem: Crosscutting Concerns

In object-oriented programming, we organize code into classes representing business entities. However, some concerns **cut across** multiple classes:

```
                    Logging    Security    Transactions    Caching
                       │          │            │             │
     ┌─────────────────┼──────────┼────────────┼─────────────┼─────┐
     │ OrderService    ●          ●            ●             ●     │
     ├─────────────────┼──────────┼────────────┼─────────────┼─────┤
     │ PaymentService  ●          ●            ●             ●     │
     ├─────────────────┼──────────┼────────────┼─────────────┼─────┤
     │ InventoryService●          ●            ●             ●     │
     └─────────────────┼──────────┼────────────┼─────────────┼─────┘
                       │          │            │             │
                  Cross-cutting concerns (horizontal slices)
```

This leads to **code tangling** (mixing concerns in one class) and **code scattering** (same concern duplicated across classes).

#### AOP Paradigm

AOP introduces a new dimension of modularity by allowing you to define:
- **What** to do (the advice — logging code, security check, etc.)
- **Where** to do it (the pointcut — which methods/classes)
- **When** to do it (advice type — before, after, around)

#### Weaving Types

**Weaving** is the process of integrating aspects with target code:

| Weaving Type | When | Tool | Pros | Cons |
|---|---|---|---|---|
| **Compile-time** | At compilation | AspectJ compiler (ajc) | Full power, best performance | Requires special compiler |
| **Load-time (LTW)** | At class loading | AspectJ agent | No special compiler, full AspectJ | JVM agent required |
| **Runtime (Proxy)** | At runtime | Spring AOP | Simple, no special tools | Limited to method execution |

> **Spring AOP uses runtime weaving** via dynamic proxies. For field access, constructor interception, or static method advice, use full AspectJ.

#### Spring AOP vs AspectJ

| Feature | Spring AOP | AspectJ |
|---|---|---|
| **Join Points** | Method execution only | Method, constructor, field, static init, etc. |
| **Weaving** | Runtime (proxy) | Compile-time, load-time |
| **Proxy Type** | JDK Dynamic / CGLIB | Direct bytecode modification |
| **Performance** | Slight overhead per call | No runtime overhead |
| **Complexity** | Low (pure Spring) | Higher (ajc compiler or agent) |
| **Use Case** | Enterprise services (tx, security) | Fine-grained interception |

> **Recommendation:** Use Spring AOP for typical enterprise concerns. Use AspectJ only when you need constructor/field interception or maximum performance.

#### The Proxy Pattern in AOP

```
  Client                    Proxy                     Target
    │                         │                          │
    │─── call method() ─────→ │                          │
    │                         │─── run before advice ───→│
    │                         │─── delegate to target ──→│
    │                         │                          │── execute method
    │                         │←── return result ────────│
    │                         │─── run after advice ────→│
    │←── return result ───────│                          │
```

The proxy intercepts all calls, applies advice, and delegates to the actual object.

### 5.1 Why AOP?

**Cross-cutting concerns** — logging, security, transactions, caching — scatter across many classes. AOP lets you modularize them into **aspects**.

```
Without AOP:                          With AOP:
┌───────────────────┐                ┌──────────────┐
│ OrderService      │                │ OrderService  │  ← clean business logic
│  log(...)         │                └──────────────┘
│  checkAuth(...)   │                      ↑
│  startTx(...)     │                ┌─────┴──────┐
│  placeOrder(...)  │                │ AOP Proxies │
│  commitTx(...)    │                │ @Log        │
│  log(...)         │                │ @Secured    │
└───────────────────┘                │ @Transact.  │
                                     └─────────────┘
```

### 5.2 AOP Terminology

| Term | Meaning | Example |
|---|---|---|
| **Aspect** | Module encapsulating a cross-cutting concern | `@Aspect` class `LoggingAspect` |
| **Join Point** | A point during execution | Any method execution |
| **Advice** | Action taken at a join point | The logging code itself |
| **Pointcut** | Expression matching join points | `execution(* com.example.service.*.*(..))` |
| **Target Object** | The object being advised | `OrderService` |
| **Proxy** | The wrapper Spring creates | JDK or CGLIB proxy |
| **Weaving** | Linking aspects with targets | At runtime (Spring default) |

### 5.3 Enabling AOP

```java
@Configuration
@EnableAspectJAutoProxy // enables proxy-based AOP
public class AopConfig { }
```

For class-based proxying (CGLIB — no interface required):
```java
@EnableAspectJAutoProxy(proxyTargetClass = true)
```

### 5.4 Advice Types

```java
@Aspect
@Component
public class LoggingAspect {

    // ────── 1. Before ──────
    @Before("execution(* com.example.service.*.*(..))")
    public void logBefore(JoinPoint jp) {
        log.info("→ Entering: {}", jp.getSignature().getName());
    }

    // ────── 2. After Returning ──────
    @AfterReturning(pointcut = "execution(* com.example.service.*.*(..))", returning = "result")
    public void logAfterReturning(JoinPoint jp, Object result) {
        log.info("← {} returned: {}", jp.getSignature().getName(), result);
    }

    // ────── 3. After Throwing ──────
    @AfterThrowing(pointcut = "execution(* com.example.service.*.*(..))", throwing = "ex")
    public void logException(JoinPoint jp, Exception ex) {
        log.error("✖ {} threw: {}", jp.getSignature().getName(), ex.getMessage());
    }

    // ────── 4. After (Finally) ──────
    @After("execution(* com.example.service.*.*(..))")
    public void logAfter(JoinPoint jp) {
        log.info("◼ {} completed (finally)", jp.getSignature().getName());
    }

    // ────── 5. Around (MOST POWERFUL) ──────
    @Around("execution(* com.example.service.*.*(..))")
    public Object measureTime(ProceedingJoinPoint pjp) throws Throwable {
        long start = System.nanoTime();
        try {
            Object result = pjp.proceed(); // call the actual method
            return result;
        } finally {
            long elapsed = (System.nanoTime() - start) / 1_000_000;
            log.info("⏱ {} took {} ms", pjp.getSignature().getName(), elapsed);
        }
    }
}
```

### 5.5 Pointcut Expressions — Deep Dive

```java
// Match all public methods in service package
execution(public * com.example.service.*.*(..))

// Match methods returning void
execution(void com.example..*Service.delete*(..))

// Match methods with exactly 2 args
execution(* com.example.service.*.*(String, int))

// Combine with && || !
@Pointcut("execution(* com.example.service.*.*(..)) && !execution(* *.get*(..))")
public void nonGetterServiceMethods() {}

// Match by annotation
@annotation(com.example.annotation.Loggable)

// Match all beans ending with "Service"
bean(*Service)

// Match all methods within annotated class
@within(org.springframework.stereotype.Service)

// Match argument annotations
@args(com.example.annotation.Validated)

// Target specific type
target(com.example.service.OrderService)
```

### 5.6 Custom Annotation + AOP

```java
// 1. Define annotation
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface AuditLog {
    String action() default "";
}

// 2. Aspect
@Aspect
@Component
public class AuditAspect {

    @Around("@annotation(auditLog)")
    public Object audit(ProceedingJoinPoint pjp, AuditLog auditLog) throws Throwable {
        String user = SecurityContextHolder.getContext().getAuthentication().getName();
        log.info("AUDIT: user={}, action={}, method={}", user, auditLog.action(),
                 pjp.getSignature().getName());
        return pjp.proceed();
    }
}

// 3. Usage
@Service
public class AccountService {

    @AuditLog(action = "TRANSFER")
    public void transfer(Account from, Account to, BigDecimal amount) { ... }
}
```

### 5.7 Spring AOP — JDK Proxy vs CGLIB Proxy

| Feature | JDK Dynamic Proxy | CGLIB Proxy |
|---|---|---|
| **Requirement** | Target implements an interface | No interface needed |
| **Mechanism** | `java.lang.reflect.Proxy` | Bytecode generation (subclassing) |
| **Speed** | Slightly slower to invoke | Slightly faster invocation |
| **`final` methods** | N/A (interface-based) | Cannot proxy `final` methods/classes |
| **Default in Spring Boot** | ✗ | ✓ (since Spring Boot 2.0) |

> **Self-invocation caveat:** If a bean method calls another method on `this`, the proxy is **bypassed** and AOP advice does **not** run. Solution: inject the bean into itself, or use `AopContext.currentProxy()`.

```java
@Service
public class OrderService {

    @Transactional
    public void placeOrder(Order o) {
        // This WILL have the transaction
    }

    public void bulkPlace(List<Order> orders) {
        orders.forEach(this::placeOrder);
        // ⚠ placeOrder() called via `this` — @Transactional SKIPPED!
    }
}
```

---

## 6. Spring MVC & REST

### 6.0 MVC Design Pattern — Theory

#### The Model-View-Controller Pattern

MVC is an **architectural pattern** that separates an application into three interconnected components:

```
┌─────────────────────────────────────────────────────────────────┐
│                          USER                                   │
│                           │                                     │
│              ┌────────────▼────────────┐                        │
│              │       CONTROLLER        │                        │
│              │    (handles requests,   │                        │
│              │     orchestrates flow)  │                        │
│              └────────────┬────────────┘                        │
│                    │             │                              │
│           updates  │             │  selects                     │
│                    ▼             ▼                              │
│         ┌──────────────┐   ┌──────────────┐                    │
│         │    MODEL     │   │    VIEW      │                    │
│         │  (business   │   │  (renders    │                    │
│         │   data &     │◄──│   UI)        │                    │
│         │   logic)     │   │              │                    │
│         └──────────────┘   └──────────────┘                    │
│                                    │                            │
│                         ┌──────────▼──────────┐                │
│                         │   RESPONSE TO USER  │                │
│                         └─────────────────────┘                │
└─────────────────────────────────────────────────────────────────┘
```

**Component Responsibilities:**

| Component | Responsibility | Spring Equivalent |
|---|---|---|
| **Model** | Business data, state, logic | Domain objects, `@Service`, `Model` attribute |
| **View** | Presentation logic, rendering | JSP, Thymeleaf, JSON/XML serializers |
| **Controller** | Request handling, flow control | `@Controller`, `@RestController` |

#### Front Controller Pattern

Spring MVC implements the **Front Controller** pattern:

> A single controller (DispatcherServlet) handles all incoming requests, then delegates to appropriate handlers.

**Benefits:**
- **Centralized control** — Common processing (security, logging) in one place
- **Consistent handling** — All requests follow the same pipeline
- **Decoupling** — Controllers don't need to know Servlet API details

#### Push vs Pull MVC

| Model | Description | Spring Approach |
|---|---|---|
| **Push (Action-based)** | Controller pushes data to view | Traditional Spring MVC |
| **Pull (Component-based)** | View pulls data from model | Less common in Spring |

#### Request-Response Lifecycle Theory

```
1. Request arrives at server
2. Servlet container routes to DispatcherServlet
3. DispatcherServlet consults HandlerMapping(s)
4. Appropriate handler (controller method) identified
5. HandlerAdapter invokes the handler
6. Handler returns ModelAndView (or @ResponseBody data)
7. ViewResolver resolves logical view name to actual View
8. View renders the response
9. Response sent to client
```

### 6.1 Architecture — DispatcherServlet Request Flow

```
Client Request
     │
     ▼
┌─────────────────────────┐
│    DispatcherServlet     │  ← Front Controller
│    (single entry point)  │
└────────┬────────────────┘
         │
    ① Consult HandlerMapping
         │
         ▼
    ② Find Handler (Controller method)
         │
    ③ Call HandlerAdapter
         │
         ▼
┌──────────────────────────┐
│   @Controller / @Rest    │
│   method executes        │
│   returns ModelAndView   │
│   or @ResponseBody       │
└────────┬─────────────────┘
         │
    ④ ViewResolver (if needed)
         │
         ▼
    ⑤ Render View / Write JSON
         │
         ▼
  HTTP Response → Client
```

### 6.2 `@Controller` vs `@RestController`

```java
// Traditional MVC — returns VIEW name
@Controller
public class PageController {

    @GetMapping("/home")
    public String home(Model model) {
        model.addAttribute("message", "Welcome");
        return "home"; // → ViewResolver → /WEB-INF/views/home.jsp
    }
}

// REST API — returns DATA directly
@RestController // = @Controller + @ResponseBody on every method
@RequestMapping("/api/v1/orders")
public class OrderController {

    @GetMapping
    public List<OrderDto> getAll() {
        return orderService.findAll();
    }

    @GetMapping("/{id}")
    public ResponseEntity<OrderDto> getById(@PathVariable Long id) {
        return orderService.findById(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public OrderDto create(@Valid @RequestBody CreateOrderRequest request) {
        return orderService.create(request);
    }

    @PutMapping("/{id}")
    public OrderDto update(@PathVariable Long id,
                           @Valid @RequestBody UpdateOrderRequest request) {
        return orderService.update(id, request);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void delete(@PathVariable Long id) {
        orderService.delete(id);
    }
}
```

### 6.3 Request Parameter Binding

```java
// Path variable:   /api/users/42
@GetMapping("/users/{id}")
public User getUser(@PathVariable("id") Long userId) { }

// Query params:    /api/users?name=John&age=30
@GetMapping("/users")
public List<User> search(@RequestParam String name,
                         @RequestParam(required = false, defaultValue = "0") int age) { }

// Headers
@GetMapping("/data")
public String getData(@RequestHeader("X-Correlation-Id") String correlationId) { }

// Cookie
@GetMapping("/prefs")
public String prefs(@CookieValue("theme") String theme) { }

// Matrix variables: /cars;color=red;year=2024
@GetMapping("/cars/{specs}")
public List<Car> filter(@MatrixVariable Map<String, String> specs) { }
```

### 6.4 `ResponseEntity` — Full Control

```java
@GetMapping("/{id}")
public ResponseEntity<OrderDto> getOrder(@PathVariable Long id) {
    return orderService.findById(id)
        .map(order -> ResponseEntity
            .ok()
            .header("X-Order-Status", order.getStatus().name())
            .cacheControl(CacheControl.maxAge(30, TimeUnit.SECONDS))
            .body(order))
        .orElse(ResponseEntity.notFound().build());
}
```

### 6.5 Exception Handling

#### Per-Controller: `@ExceptionHandler`
```java
@RestController
public class OrderController {

    @ExceptionHandler(OrderNotFoundException.class)
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ErrorResponse handleNotFound(OrderNotFoundException ex) {
        return new ErrorResponse("ORDER_NOT_FOUND", ex.getMessage());
    }
}
```

#### Global: `@ControllerAdvice`
```java
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorResponse> handleValidation(MethodArgumentNotValidException ex) {
        List<String> errors = ex.getBindingResult().getFieldErrors().stream()
            .map(fe -> fe.getField() + ": " + fe.getDefaultMessage())
            .toList();
        return ResponseEntity.badRequest()
            .body(new ErrorResponse("VALIDATION_FAILED", errors));
    }

    @ExceptionHandler(DataAccessException.class)
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public ErrorResponse handleDbError(DataAccessException ex) {
        return new ErrorResponse("DB_ERROR", "A database error occurred");
    }

    @ExceptionHandler(Exception.class)
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public ErrorResponse handleAll(Exception ex) {
        log.error("Unhandled exception", ex);
        return new ErrorResponse("INTERNAL_ERROR", "Something went wrong");
    }
}
```

### 6.6 Validation with Bean Validation (JSR 380)

```java
public record CreateOrderRequest(
    @NotBlank(message = "Customer name is required")
    String customerName,

    @NotEmpty(message = "At least one item required")
    List<@Valid OrderItemRequest> items,

    @Email
    String contactEmail,

    @Future
    LocalDate deliveryDate
) {}

public record OrderItemRequest(
    @NotNull Long productId,
    @Min(1) @Max(100) int quantity
) {}
```

Use `@Valid` on the controller parameter to trigger validation:
```java
@PostMapping
public OrderDto create(@Valid @RequestBody CreateOrderRequest request) { }
```

### 6.7 Content Negotiation

```java
@Configuration
public class WebConfig implements WebMvcConfigurer {

    @Override
    public void configureContentNegotiation(ContentNegotiationConfigurer configurer) {
        configurer
            .defaultContentType(MediaType.APPLICATION_JSON)
            .favorParameter(true)          // ?format=xml
            .parameterName("format")
            .mediaType("json", MediaType.APPLICATION_JSON)
            .mediaType("xml", MediaType.APPLICATION_XML);
    }
}
```

### 6.8 Interceptors

```java
@Component
public class RequestTimingInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response,
                             Object handler) {
        request.setAttribute("startTime", System.nanoTime());
        return true; // false = abort request
    }

    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response,
                           Object handler, ModelAndView modelAndView) {
        // After controller, before view rendering
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response,
                                Object handler, Exception ex) {
        long start = (long) request.getAttribute("startTime");
        log.info("Request {} took {} ms", request.getRequestURI(),
                 (System.nanoTime() - start) / 1_000_000);
    }
}

@Configuration
public class WebConfig implements WebMvcConfigurer {
    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(new RequestTimingInterceptor())
                .addPathPatterns("/api/**")
                .excludePathPatterns("/api/health");
    }
}
```

### 6.9 Filters vs Interceptors vs AOP

| Feature | Servlet Filter | HandlerInterceptor | AOP Advice |
|---|---|---|---|
| **Level** | Servlet container | Spring MVC | Any Spring bean |
| **Access to** | `HttpServletRequest/Response` | Handler + ModelAndView | Method arguments, return value |
| **Use case** | CORS, encoding, security | Logging, auth, timing | Business cross-cutting |
| **Order** | `@Order` / web.xml | Registry order | `@Order` on aspect |
| **Works with non-Spring?** | Yes | No | No |

---

## 7. Data Access — JDBC, ORM & Transactions

### 7.0 Transaction Theory — ACID Properties

#### What is a Transaction?

A **transaction** is a logical unit of work that must be executed completely or not at all. It groups multiple operations into a single atomic unit.

**Real-world analogy:** Bank transfer — debiting one account and crediting another must both succeed or both fail.

#### ACID Properties

| Property | Meaning | Example |
|---|---|---|
| **Atomicity** | All operations succeed or all fail | Transfer: debit + credit both or neither |
| **Consistency** | Database moves from one valid state to another | Total balance remains unchanged after transfer |
| **Isolation** | Concurrent transactions don't interfere | Two transfers on same account don't corrupt |
| **Durability** | Committed data survives system failure | Once transfer confirmed, it's permanent |

#### Isolation Levels — Deep Theory

**Concurrency Problems:**

| Problem | Description | Example |
|---|---|---|
| **Dirty Read** | Reading uncommitted data from another tx | Tx1 writes $100, Tx2 reads $100, Tx1 rollbacks → Tx2 has phantom data |
| **Non-Repeatable Read** | Same query returns different data within tx | Tx1 reads balance=$100, Tx2 updates to $200, Tx1 reads again → $200 |
| **Phantom Read** | New rows appear in repeated query | Tx1 counts 5 orders, Tx2 inserts 1, Tx1 counts again → 6 orders |

**Isolation Levels Explained:**

```
                            Dirty    Non-Repeatable   Phantom
     Level                  Read     Read             Read        Performance
  ─────────────────────────────────────────────────────────────────────────────
  READ_UNCOMMITTED          Yes      Yes              Yes         Fastest
  READ_COMMITTED            No       Yes              Yes         Fast
  REPEATABLE_READ           No       No               Yes         Moderate
  SERIALIZABLE              No       No               No          Slowest
```

> **Default for most databases:** `READ_COMMITTED` (PostgreSQL, Oracle, SQL Server)  
> **MySQL InnoDB default:** `REPEATABLE_READ`

#### Local vs Distributed Transactions

| Type | Scope | Coordination | Spring Support |
|---|---|---|---|
| **Local** | Single database | None needed | `DataSourceTransactionManager` |
| **Global/XA** | Multiple resources | Two-Phase Commit (2PC) | `JtaTransactionManager` |
| **Saga** | Microservices | Compensating transactions | Manual or frameworks (Axon, Temporal) |

**Two-Phase Commit (2PC):**
```
     Coordinator                 Resource A              Resource B
          │                          │                       │
 Phase 1: │── PREPARE ──────────────→│                       │
          │── PREPARE ────────────────────────────────────→  │
          │←─ VOTE YES ──────────────│                       │
          │←─ VOTE YES ───────────────────────────────────── │
          │                          │                       │
 Phase 2: │── COMMIT ───────────────→│                       │
          │── COMMIT ─────────────────────────────────────→  │
          │←─ ACK ───────────────────│                       │
          │←─ ACK ────────────────────────────────────────── │
```

> **Microservices caution:** 2PC is not recommended for distributed systems (blocking, single point of failure). Use the **Saga pattern** instead.

#### Programmatic vs Declarative Transaction Management

| Approach | How | Pros | Cons |
|---|---|---|---|
| **Programmatic** | `TransactionTemplate`, manual `commit()/rollback()` | Fine-grained control | Boilerplate, error-prone |
| **Declarative** | `@Transactional` annotation | Clean, DRY | Less flexible |

> **Spring recommendation:** Use declarative (`@Transactional`) for 95% of cases. Use programmatic only for complex scenarios requiring dynamic tx behavior.

### 7.1 Spring's Data Access Exception Hierarchy

Spring translates vendor-specific exceptions into a **consistent, unchecked** hierarchy:

```
DataAccessException (root - unchecked)
├── NonTransientDataAccessException
│   ├── DataIntegrityViolationException
│   ├── DuplicateKeyException
│   └── EmptyResultDataAccessException
├── TransientDataAccessException
│   ├── QueryTimeoutException
│   └── ConcurrencyFailureException
└── RecoverableDataAccessException
```

> All are **unchecked** (`RuntimeException`), unlike raw JDBC's `SQLException`.

### 7.2 `JdbcTemplate` — Direct SQL

```java
@Repository
public class JdbcOrderRepository {

    private final JdbcTemplate jdbc;

    public JdbcOrderRepository(DataSource dataSource) {
        this.jdbc = new JdbcTemplate(dataSource);
    }

    // Query for single object
    public Order findById(Long id) {
        return jdbc.queryForObject(
            "SELECT id, customer, total FROM orders WHERE id = ?",
            (rs, rowNum) -> new Order(rs.getLong("id"),
                                       rs.getString("customer"),
                                       rs.getBigDecimal("total")),
            id
        );
    }

    // Query for list
    public List<Order> findByCustomer(String customer) {
        return jdbc.query(
            "SELECT id, customer, total FROM orders WHERE customer = ?",
            (rs, rowNum) -> new Order(rs.getLong("id"),
                                       rs.getString("customer"),
                                       rs.getBigDecimal("total")),
            customer
        );
    }

    // Insert
    public int save(Order order) {
        return jdbc.update(
            "INSERT INTO orders (customer, total) VALUES (?, ?)",
            order.getCustomer(), order.getTotal()
        );
    }

    // Batch insert
    public int[] batchSave(List<Order> orders) {
        return jdbc.batchUpdate(
            "INSERT INTO orders (customer, total) VALUES (?, ?)",
            new BatchPreparedStatementSetter() {
                @Override
                public void setValues(PreparedStatement ps, int i) throws SQLException {
                    ps.setString(1, orders.get(i).getCustomer());
                    ps.setBigDecimal(2, orders.get(i).getTotal());
                }
                @Override
                public int getBatchSize() { return orders.size(); }
            }
        );
    }
}
```

### 7.3 `NamedParameterJdbcTemplate`

```java
@Repository
public class OrderRepository {
    private final NamedParameterJdbcTemplate namedJdbc;

    public List<Order> findByFilters(String customer, BigDecimal minTotal) {
        String sql = "SELECT * FROM orders WHERE customer = :customer AND total >= :minTotal";
        MapSqlParameterSource params = new MapSqlParameterSource()
            .addValue("customer", customer)
            .addValue("minTotal", minTotal);
        return namedJdbc.query(sql, params, orderRowMapper);
    }
}
```

### 7.4 Spring ORM — JPA Integration

```java
@Configuration
@EnableTransactionManagement
public class JpaConfig {

    @Bean
    public LocalContainerEntityManagerFactoryBean entityManagerFactory(DataSource ds) {
        LocalContainerEntityManagerFactoryBean em = new LocalContainerEntityManagerFactoryBean();
        em.setDataSource(ds);
        em.setPackagesToScan("com.example.entity");
        em.setJpaVendorAdapter(new HibernateJpaVendorAdapter());
        em.setJpaProperties(hibernateProperties());
        return em;
    }

    @Bean
    public PlatformTransactionManager transactionManager(EntityManagerFactory emf) {
        return new JpaTransactionManager(emf);
    }

    private Properties hibernateProperties() {
        Properties props = new Properties();
        props.put("hibernate.hbm2ddl.auto", "validate");
        props.put("hibernate.show_sql", "true");
        props.put("hibernate.dialect", "org.hibernate.dialect.PostgreSQLDialect");
        return props;
    }
}
```

### 7.5 `@Transactional` — Deep Dive

```java
@Service
public class OrderService {

    @Transactional // default: propagation=REQUIRED, isolation=DEFAULT, readOnly=false
    public void placeOrder(Order order) {
        orderRepo.save(order);
        inventoryService.reduceStock(order.getItems()); // participates in SAME tx
        notificationService.sendConfirmation(order);     // participates in SAME tx
    }

    @Transactional(readOnly = true)
    public List<Order> getRecentOrders() {
        return orderRepo.findRecent(); // Hibernate flush mode = MANUAL, potential perf gain
    }

    @Transactional(
        propagation = Propagation.REQUIRES_NEW,
        isolation = Isolation.REPEATABLE_READ,
        timeout = 30,
        rollbackFor = BusinessException.class,
        noRollbackFor = NotificationException.class
    )
    public void processPayment(Payment payment) { ... }
}
```

### 7.6 Transaction Propagation Types

| Propagation | Behaviour |
|---|---|
| `REQUIRED` (default) | Join existing tx, or create new if none |
| `REQUIRES_NEW` | Always create a **new** tx (suspend existing) |
| `SUPPORTS` | Use tx if one exists, otherwise run without |
| `NOT_SUPPORTED` | Always run without tx (suspend existing) |
| `MANDATORY` | Must have an existing tx, else throw exception |
| `NEVER` | Must **not** have a tx, else throw exception |
| `NESTED` | Nested tx with savepoints (JDBC only) |

**Visual:**
```
Caller                    placeOrder()              sendNotification()
  │                      REQUIRED                   REQUIRES_NEW
  │─── No Tx ──────────→ Creates Tx A ────────────→ Suspends Tx A
  │                                                  Creates Tx B
  │                                                  Commits Tx B
  │                      ← Resumes Tx A ────────────
  │                      Commits Tx A
```

### 7.7 Isolation Levels

| Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| `READ_UNCOMMITTED` | ✓ | ✓ | ✓ |
| `READ_COMMITTED` | ✗ | ✓ | ✓ |
| `REPEATABLE_READ` | ✗ | ✗ | ✓ |
| `SERIALIZABLE` | ✗ | ✗ | ✗ |

### 7.8 Common `@Transactional` Pitfalls

```java
// ❌ Pitfall 1: Private method — proxy can't intercept
@Transactional
private void updateOrder() { } // AOP proxy ignores private methods

// ❌ Pitfall 2: Self-invocation — bypasses proxy
public void bulkUpdate(List<Order> orders) {
    orders.forEach(this::updateSingle); // `this` — no proxy!
}
@Transactional
public void updateSingle(Order o) { }

// ❌ Pitfall 3: Catching exception silently
@Transactional
public void process() {
    try {
        riskyOperation();
    } catch (Exception e) {
        log.error("Swallowed", e);
        // Transaction will COMMIT because no exception propagated!
    }
}

// ✅ Fix: re-throw or mark for rollback
@Transactional
public void process() {
    try {
        riskyOperation();
    } catch (Exception e) {
        TransactionAspectSupport.currentTransactionStatus().setRollbackOnly();
        throw e;
    }
}
```

---

## 8. Spring Security Fundamentals

### 8.0 Security Concepts — Theory

#### Authentication vs Authorization

| Concept | Question Answered | Example |
|---|---|---|
| **Authentication** | "Who are you?" | Login with username/password |
| **Authorization** | "What can you do?" | Can this user delete records? |

```
  Request → [Authentication] → Identity Established → [Authorization] → Access Granted/Denied
                  │                                          │
             "Is this user                            "Does this user
              who they claim?"                        have permission?"
```

#### Security Principles

| Principle | Description |
|---|---|
| **Defense in Depth** | Multiple security layers; if one fails, others protect |
| **Least Privilege** | Grant minimum permissions necessary |
| **Fail Secure** | Default to deny; explicitly allow |
| **Separation of Duties** | Split critical functions among multiple actors |
| **Security by Design** | Build security in from the start, not bolted on |

#### Common Attack Vectors (and Spring's Protection)

| Attack | Description | Spring Protection |
|---|---|---|
| **CSRF** | Forged requests from authenticated sessions | CSRF tokens (`CsrfFilter`) |
| **XSS** | Malicious scripts injected into pages | Content-Type headers, output encoding |
| **SQL Injection** | Malicious SQL in user input | Parameterized queries (JdbcTemplate) |
| **Session Fixation** | Attacker sets victim's session ID | Session ID regeneration on login |
| **Brute Force** | Repeated login attempts | Rate limiting, account lockout |
| **Clickjacking** | Hidden UI elements trick users | X-Frame-Options header |

#### Authentication Mechanisms

| Mechanism | How It Works | Use Case |
|---|---|---|
| **Form Login** | HTML form → server session | Traditional web apps |
| **HTTP Basic** | Base64 credentials in header | Simple APIs, internal services |
| **HTTP Digest** | Hashed credentials (more secure than Basic) | Legacy systems |
| **OAuth 2.0** | Token-based, delegated authorization | Third-party access, SPAs |
| **JWT** | Self-contained token with claims | Stateless APIs, microservices |
| **SAML** | XML-based SSO standard | Enterprise SSO |
| **LDAP** | Directory-based authentication | Corporate environments |

#### The Principal and GrantedAuthority Model

```
┌────────────────────────────────────────────┐
│           SecurityContext                  │
│  ┌──────────────────────────────────────┐  │
│  │         Authentication                │  │
│  │  ┌────────────────────────────────┐  │  │
│  │  │  Principal (who)               │  │  │
│  │  │  - username                    │  │  │
│  │  │  - user details                │  │  │
│  │  └────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────┐  │  │
│  │  │  Credentials (proof)           │  │  │
│  │  │  - password (cleared after auth)│  │  │
│  │  └────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────┐  │  │
│  │  │  Authorities (permissions)     │  │  │
│  │  │  - ROLE_USER                   │  │  │
│  │  │  - ROLE_ADMIN                  │  │  │
│  │  │  - DELETE_PRIVILEGE            │  │  │
│  │  └────────────────────────────────┘  │  │
│  └──────────────────────────────────────┘  │
└────────────────────────────────────────────┘
```

### 8.1 Architecture Overview

```
HTTP Request
     │
     ▼
┌──────────────────────────┐
│  DelegatingFilterProxy    │  ← Servlet filter, delegates to Spring
│  └─ FilterChainProxy     │
│     └─ SecurityFilterChain│
│        ├─ CorsFilter      │
│        ├─ CsrfFilter      │
│        ├─ Authentication  │
│        │  Filter           │
│        ├─ Authorization   │
│        │  Filter           │
│        └─ ExceptionTransl.│
└──────────────────────────┘
            │
            ▼
     DispatcherServlet → Controller
```

### 8.2 Key Components

| Component | Role |
|---|---|
| `SecurityFilterChain` | Ordered chain of security filters |
| `AuthenticationManager` | Delegates to `AuthenticationProvider`(s) |
| `AuthenticationProvider` | Performs actual authentication (DB lookup, LDAP, etc.) |
| `UserDetailsService` | Loads user data (`UserDetails`) from a store |
| `PasswordEncoder` | Hashes/verifies passwords |
| `SecurityContextHolder` | Stores `Authentication` object (ThreadLocal) |
| `GrantedAuthority` | Represents a permission/role |

### 8.3 Java Configuration (Spring Security 6+)

```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .csrf(csrf -> csrf.disable()) // disable for stateless APIs
            .sessionManagement(sm -> sm.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/public/**").permitAll()
                .requestMatchers("/api/admin/**").hasRole("ADMIN")
                .requestMatchers(HttpMethod.DELETE, "/api/**").hasAuthority("DELETE_PRIVILEGE")
                .anyRequest().authenticated()
            )
            .httpBasic(Customizer.withDefaults())
            .build();
    }

    @Bean
    public UserDetailsService userDetailsService() {
        return username -> userRepository.findByUsername(username)
            .map(u -> User.builder()
                .username(u.getUsername())
                .password(u.getPassword())
                .authorities(u.getRoles().stream()
                    .map(r -> new SimpleGrantedAuthority("ROLE_" + r.getName()))
                    .toList()
                    .toArray(new GrantedAuthority[0]))
                .build())
            .orElseThrow(() -> new UsernameNotFoundException("User not found: " + username));
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(12);
    }
}
```

### 8.4 Method-Level Security

```java
@Configuration
@EnableMethodSecurity // replaces @EnableGlobalMethodSecurity
public class MethodSecurityConfig { }

@Service
public class AccountService {

    @PreAuthorize("hasRole('ADMIN')")
    public void deleteAccount(Long id) { }

    @PreAuthorize("#username == authentication.name or hasRole('ADMIN')")
    public Account getAccount(String username) { }

    @PostAuthorize("returnObject.owner == authentication.name")
    public Document getDocument(Long docId) { }

    @PreFilter("filterObject.owner == authentication.name")
    public void batchDelete(List<Document> docs) { }
}
```

### 8.5 CORS Configuration

```java
@Bean
public CorsConfigurationSource corsConfigurationSource() {
    CorsConfiguration config = new CorsConfiguration();
    config.setAllowedOrigins(List.of("https://frontend.example.com"));
    config.setAllowedMethods(List.of("GET", "POST", "PUT", "DELETE"));
    config.setAllowedHeaders(List.of("*"));
    config.setAllowCredentials(true);
    config.setMaxAge(3600L);

    UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
    source.registerCorsConfiguration("/api/**", config);
    return source;
}
```

---

## 9. Testing with Spring

### 9.0 Testing Philosophy & Theory

#### Why Testing Matters

> "Code without tests is broken by design." — Jacob Kaplan-Moss

Tests provide:
- **Confidence** — Deploy without fear of breaking existing functionality
- **Documentation** — Tests describe how code should behave
- **Design feedback** — Hard-to-test code often indicates poor design
- **Regression prevention** — Catch bugs before they reach production
- **Refactoring safety net** — Change implementation without changing behavior

#### Test Categories

| Category | Scope | Speed | Dependencies | Purpose |
|---|---|---|---|---|
| **Unit** | Single class/method | Milliseconds | Mocked | Verify logic in isolation |
| **Integration** | Multiple components | Seconds | Real (some) | Verify component interaction |
| **End-to-End** | Entire system | Minutes | Real (all) | Verify user scenarios |
| **Contract** | API boundaries | Fast | Stubs | Verify API contracts |
| **Performance** | System under load | Varies | Real | Verify non-functional requirements |

#### The FIRST Principles of Good Tests

| Principle | Meaning |
|---|---|
| **Fast** | Tests should run quickly (feedback loop) |
| **Independent** | Tests shouldn't depend on each other |
| **Repeatable** | Same result every time, any environment |
| **Self-validating** | Pass or fail, no manual interpretation |
| **Timely** | Written before or alongside production code |

#### Arrange-Act-Assert (AAA) Pattern

```java
@Test
void calculateDiscount_premiumCustomer_returns20Percent() {
    // Arrange — Set up test data and dependencies
    Customer customer = new Customer("John", CustomerTier.PREMIUM);
    DiscountCalculator calculator = new DiscountCalculator();
    
    // Act — Execute the behavior being tested
    BigDecimal discount = calculator.calculate(customer, new BigDecimal("100.00"));
    
    // Assert — Verify the expected outcome
    assertThat(discount).isEqualByComparingTo("20.00");
}
```

#### Test Doubles

| Double | Purpose | Example |
|---|---|---|
| **Dummy** | Fill parameter lists (never used) | `new Object()` |
| **Stub** | Return canned answers | `when(repo.find()).thenReturn(entity)` |
| **Spy** | Record calls + delegate to real impl | `spy(realService)` |
| **Mock** | Verify interactions | `verify(mock).method()` |
| **Fake** | Working but simplified impl | In-memory database |

### 9.1 Testing Strategy Pyramid

```
           ┌──────────┐
           │  E2E     │  ← Fewest (expensive, slow, brittle)
           ├──────────┤
           │Integration│  ← @SpringBootTest, @WebMvcTest
           ├──────────┤
           │  Unit     │  ← Most (fast, isolated, reliable)
           └──────────┘
```

**Recommended ratio:** ~70% unit, ~20% integration, ~10% E2E

### 9.2 Unit Testing (No Spring Context)

**Key principle:** Unit tests should NOT load Spring context. Use constructor injection to easily instantiate classes with mock dependencies.

```java
@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock
    private OrderRepository orderRepo;

    @Mock
    private NotificationService notificationService;

    @InjectMocks
    private OrderService orderService;

    @Test
    void placeOrder_savesAndNotifies() {
        Order order = new Order("customer1", BigDecimal.TEN);
        when(orderRepo.save(any())).thenReturn(order);

        orderService.placeOrder(order);

        verify(orderRepo).save(order);
        verify(notificationService).sendConfirmation(order);
    }

    @Test
    void placeOrder_insufficientStock_throwsException() {
        Order order = new Order("customer1", BigDecimal.TEN);
        doThrow(new InsufficientStockException("Item A"))
            .when(orderRepo).save(order);

        assertThrows(InsufficientStockException.class,
                     () -> orderService.placeOrder(order));
        verify(notificationService, never()).sendConfirmation(any());
    }
}
```

### 9.3 Integration Testing with Spring Context

```java
@SpringJUnitConfig(classes = {AppConfig.class, TestDataSourceConfig.class})
@Transactional // auto-rollback after each test
class OrderRepositoryIntegrationTest {

    @Autowired
    private OrderRepository orderRepo;

    @Autowired
    private JdbcTemplate jdbc;

    @Test
    void save_persistsToDatabase() {
        Order order = new Order("customer1", new BigDecimal("99.99"));

        orderRepo.save(order);

        Integer count = jdbc.queryForObject(
            "SELECT COUNT(*) FROM orders WHERE customer = ?",
            Integer.class, "customer1");
        assertThat(count).isEqualTo(1);
    }
}
```

### 9.4 Testing Web Layer — `MockMvc`

```java
@WebMvcTest(OrderController.class)
class OrderControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private OrderService orderService;

    @Test
    void getOrder_returnsOk() throws Exception {
        OrderDto dto = new OrderDto(1L, "customer1", BigDecimal.TEN);
        when(orderService.findById(1L)).thenReturn(Optional.of(dto));

        mockMvc.perform(get("/api/v1/orders/1")
                .accept(MediaType.APPLICATION_JSON))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.customerName").value("customer1"))
            .andExpect(jsonPath("$.total").value(10));
    }

    @Test
    void createOrder_invalidRequest_returns400() throws Exception {
        String invalidJson = """
            { "customerName": "", "items": [] }
            """;

        mockMvc.perform(post("/api/v1/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content(invalidJson))
            .andExpect(status().isBadRequest())
            .andExpect(jsonPath("$.code").value("VALIDATION_FAILED"));
    }
}
```

### 9.5 Testing with Profiles and Properties

```java
@ActiveProfiles("test")
@TestPropertySource(properties = {
    "spring.datasource.url=jdbc:h2:mem:testdb",
    "cache.enabled=false"
})
@SpringJUnitConfig(AppConfig.class)
class ServiceLayerTest { }
```

### 9.6 `@DirtiesContext`

```java
@DirtiesContext(classMode = DirtiesContext.ClassMode.AFTER_EACH_TEST_METHOD)
class StatefulBeanTest {
    // Spring recreates the ApplicationContext after each test
    // Use sparingly — slow!
}
```

---

## 10. Event System & Async Processing

### 10.0 Event-Driven Architecture — Theory

#### Why Events?

Events enable **loose coupling** between components. Instead of direct method calls, components communicate through events:

```
  Direct Coupling:                    Event-Driven:
  
  OrderService                        OrderService
      │                                   │
      ├──→ NotificationService              │──→ publish(OrderPlacedEvent)
      ├──→ InventoryService                       │
      ├──→ AnalyticsService              ┌───────▼───────┐
      └──→ AuditService                  │ Event Bus    │
                                       └───────────────┘
  OrderService knows ALL                    ┌─────┬────┬────┐
  consumers (tight coupling)                ▼    ▼    ▼    ▼
                                         Notif Inv  Audit Analytics
                                         
                                       OrderService knows NOTHING
                                       about consumers (loose coupling)
```

#### Benefits of Event-Driven Design

| Benefit | Description |
|---|---|
| **Decoupling** | Publisher doesn't know (or care) about subscribers |
| **Extensibility** | Add new listeners without modifying publisher |
| **Testability** | Test publisher and listeners independently |
| **Async potential** | Events can be processed asynchronously |
| **Audit trail** | Events naturally create a log of what happened |

#### Event Patterns

| Pattern | Description | Use Case |
|---|---|---|
| **Observer** | Object notifies dependents of state changes | UI updates, caching invalidation |
| **Publish-Subscribe** | Publishers send to topics, subscribers receive | Notification systems, logging |
| **Event Sourcing** | Store state as sequence of events | Audit logs, temporal queries |
| **CQRS** | Separate read and write models | Complex domains, high scalability |

#### Synchronous vs Asynchronous Events

| Mode | Behavior | Trade-offs |
|---|---|---|
| **Synchronous** | Publisher waits for all listeners | Simple, consistent, but slower |
| **Asynchronous** | Publisher continues immediately | Fast, but complex error handling |

> **Spring default:** Synchronous. Add `@Async` to listeners for async processing.

#### Event vs Direct Method Call — When to Choose

| Use Events When... | Use Direct Calls When... |
|---|---|
| Multiple independent reactions | Single, well-defined action |
| Reactions may change/grow | Tight coupling is acceptable |
| Caller shouldn't wait for completion | Response needed immediately |
| Cross-cutting (audit, metrics) | Core business logic flow |

### 10.1 Application Events

```java
// 1. Define event
public class OrderPlacedEvent extends ApplicationEvent {
    private final Order order;

    public OrderPlacedEvent(Object source, Order order) {
        super(source);
        this.order = order;
    }
    public Order getOrder() { return order; }
}

// 2. Publish event
@Service
public class OrderService {
    private final ApplicationEventPublisher publisher;

    public OrderService(ApplicationEventPublisher publisher) {
        this.publisher = publisher;
    }

    @Transactional
    public void placeOrder(Order order) {
        orderRepo.save(order);
        publisher.publishEvent(new OrderPlacedEvent(this, order));
    }
}

// 3. Listen for event
@Component
public class NotificationListener {

    @EventListener
    public void onOrderPlaced(OrderPlacedEvent event) {
        sendEmail(event.getOrder().getCustomerEmail());
    }
}
```

### 10.2 Simplified Event (No Extends Required)

Since Spring 4.2, events don't need to extend `ApplicationEvent`:

```java
// Plain POJO event
public record OrderPlacedEvent(Long orderId, String customer) {}

// Publish
publisher.publishEvent(new OrderPlacedEvent(42L, "John"));

// Listen
@EventListener
public void handle(OrderPlacedEvent event) { }
```

### 10.3 `@TransactionalEventListener`

```java
@TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
public void onOrderCommitted(OrderPlacedEvent event) {
    // Only fires AFTER the transaction that published the event commits
    // Prevents sending email for rolled-back orders
    notificationService.sendConfirmation(event.orderId());
}
```

| Phase | When |
|---|---|
| `AFTER_COMMIT` (default) | After successful commit |
| `AFTER_ROLLBACK` | After rollback |
| `AFTER_COMPLETION` | After commit or rollback |
| `BEFORE_COMMIT` | Before commit |

### 10.4 `@Async` — Asynchronous Method Execution

```java
@Configuration
@EnableAsync
public class AsyncConfig implements AsyncConfigurer {

    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("async-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
    }

    @Override
    public AsyncUncaughtExceptionHandler getAsyncUncaughtExceptionHandler() {
        return (ex, method, params) ->
            log.error("Async error in {}: {}", method.getName(), ex.getMessage());
    }
}

@Service
public class ReportService {

    @Async
    public CompletableFuture<Report> generateReport(ReportRequest request) {
        Report report = heavyComputation(request);
        return CompletableFuture.completedFuture(report);
    }
}
```

### 10.5 `@Scheduled` — Task Scheduling

```java
@Configuration
@EnableScheduling
public class SchedulingConfig { }

@Component
public class CleanupJob {

    @Scheduled(fixedRate = 60_000) // every 60 seconds
    public void cleanTempFiles() { }

    @Scheduled(fixedDelay = 30_000) // 30s after LAST execution ends
    public void syncData() { }

    @Scheduled(cron = "0 0 2 * * MON-FRI") // 2 AM weekdays
    public void generateDailyReport() { }

    @Scheduled(cron = "${cleanup.cron}") // externalized cron
    public void configDrivenCleanup() { }
}
```

---

## 11. Spring Expression Language (SpEL)

### 11.1 Basics

```java
@Value("#{2 + 3}")                     // → 5
private int sum;

@Value("#{T(java.lang.Math).PI}")      // → 3.14159…
private double pi;

@Value("#{systemProperties['user.home']}")
private String userHome;

@Value("#{@orderService.getDefaultCurrency()}")  // call bean method
private String currency;

@Value("#{${max.threads} ?: 10}")      // Elvis operator (default value)
private int maxThreads;
```

### 11.2 SpEL in Annotations

```java
// Conditional bean creation
@Bean
@ConditionalOnExpression("#{${feature.newCheckout} == true}")
public CheckoutService newCheckoutService() { }

// Security expressions
@PreAuthorize("hasRole('ADMIN') and #order.total < 10000")
public void approveOrder(Order order) { }

// Cache key
@Cacheable(value = "orders", key = "#customerId + '_' + #status")
public List<Order> findOrders(String customerId, String status) { }
```

### 11.3 SpEL in XML

```xml
<bean id="appInfo" class="com.example.AppInfo">
    <property name="randomNumber" value="#{T(java.lang.Math).random() * 100}" />
    <property name="dbUrl" value="#{systemProperties['db.url'] ?: 'jdbc:h2:mem:default'}" />
</bean>
```

---

## 12. Resource Handling & Profiles

### 12.1 `@PropertySource` and `Environment`

```java
@Configuration
@PropertySource("classpath:app.properties")
@PropertySource("classpath:db-${spring.profiles.active}.properties") // profile-aware
public class AppConfig {

    @Autowired
    private Environment env;

    @Bean
    public DataSource dataSource() {
        HikariDataSource ds = new HikariDataSource();
        ds.setJdbcUrl(env.getProperty("db.url"));
        ds.setUsername(env.getProperty("db.username"));
        ds.setPassword(env.getProperty("db.password"));
        ds.setMaximumPoolSize(env.getProperty("db.pool.max", Integer.class, 10));
        return ds;
    }
}
```

### 12.2 `@Value` — Property Injection

```java
@Component
public class AppSettings {
    @Value("${app.name}")
    private String appName;

    @Value("${app.timeout:5000}")   // default value after colon
    private int timeout;

    @Value("${app.features}")       // comma-separated → List
    private List<String> features;

    @Value("#{${app.limits}}")      // SpEL to parse as Map from {key:val, ...}
    private Map<String, Integer> limits;
}
```

### 12.3 Profiles

```java
// Activate profiles
// 1. Programmatically
ctx.getEnvironment().setActiveProfiles("dev", "metrics");

// 2. System property
-Dspring.profiles.active=prod

// 3. Annotation
@ActiveProfiles("test") // in tests
```

**Profile-specific beans:**
```java
@Configuration
@Profile("dev")
public class DevDataSourceConfig {
    @Bean
    public DataSource dataSource() {
        return new EmbeddedDatabaseBuilder()
            .setType(EmbeddedDatabaseType.H2)
            .addScript("schema.sql")
            .build();
    }
}

@Configuration
@Profile("prod")
public class ProdDataSourceConfig {
    @Bean
    public DataSource dataSource() {
        HikariDataSource ds = new HikariDataSource();
        ds.setJdbcUrl("jdbc:postgresql://prod-host/mydb");
        return ds;
    }
}
```

**Profile expressions (Spring 5.1+):**
```java
@Profile("prod & us-east")         // AND
@Profile("dev | staging")          // OR
@Profile("!prod")                  // NOT
```

### 12.4 Resource Abstraction

```java
@Component
public class FileLoader {

    @Value("classpath:data/config.json")
    private Resource classpathResource;

    @Value("file:/opt/app/config.json")
    private Resource fileResource;

    @Value("https://api.example.com/config.json")
    private Resource urlResource;

    public String loadContent() throws IOException {
        return new String(classpathResource.getInputStream().readAllBytes(),
                          StandardCharsets.UTF_8);
    }
}
```

| Prefix | Example | Resolves |
|---|---|---|
| `classpath:` | `classpath:config.xml` | Classpath resource |
| `file:` | `file:/opt/data.txt` | File system |
| `https:` | `https://host/data` | URL |
| No prefix | `data/config.xml` | Depends on `ApplicationContext` type |

---

## 12.5 Design Patterns in Spring Framework

Understanding the design patterns Spring uses helps you leverage the framework effectively and answer architecture questions in interviews.

### Patterns Used BY Spring (Internally)

| Pattern | Where Used | Purpose |
|---|---|---|
| **Singleton** | Default bean scope | One instance per container |
| **Factory** | `BeanFactory`, `FactoryBean` | Centralized object creation |
| **Proxy** | AOP, `@Transactional`, `@Async` | Add behavior without modifying code |
| **Template Method** | `JdbcTemplate`, `RestTemplate` | Reusable algorithm with customization hooks |
| **Strategy** | `PlatformTransactionManager`, `Resource` | Swappable implementations |
| **Observer** | Application events | Decouple event publishers from listeners |
| **Decorator** | `BeanPostProcessor` | Add behavior by wrapping |
| **Adapter** | `HandlerAdapter` | Make incompatible interfaces work together |
| **Front Controller** | `DispatcherServlet` | Single entry point for requests |
| **Composite** | `CompositePropertySource` | Treat groups uniformly |

### Deep Dive: Key Patterns

#### 1. Template Method Pattern

The Template Method pattern defines an algorithm's skeleton, deferring some steps to subclasses (or callbacks in Spring's case):

```java
// Spring's JdbcTemplate uses template method pattern
// The template handles: connection, exception, cleanup
// You provide: the actual SQL and row mapping

jdbcTemplate.query(
    "SELECT * FROM users",           // Your SQL
    (rs, rowNum) -> new User(        // Your mapping
        rs.getLong("id"),
        rs.getString("name")
    )
);
// Template handles: getConnection(), createStatement(), 
// executeQuery(), close(), exception translation
```

**Why it matters:** Eliminates boilerplate. You focus on business logic; Spring handles infrastructure.

#### 2. Proxy Pattern

Proxies are **fundamental** to Spring's magic. They enable AOP, transactions, security, and lazy loading:

```
Client                 Proxy                    Target
  │                      │                        │
  │─── call method() ───▶│                        │
  │                      │─── before advice ────▶ │
  │                      │─── call target ──────▶ │
  │                      │◀── return ──────────── │
  │                      │─── after advice ─────▶ │
  │◀── return ───────────│                        │
```

#### 3. Factory Pattern

Spring's entire IoC container is essentially a sophisticated factory:

```java
// You declare what you want
@Component
public class OrderService {
    private final OrderRepository repo;
    
    public OrderService(OrderRepository repo) {
        this.repo = repo;  // Spring (the factory) provides this
    }
}

// Spring (the factory) creates and wires everything
// You never call `new OrderService(new JdbcOrderRepository(...))`
```

#### 4. Strategy Pattern

Different implementations can be swapped without changing client code:

```java
// Strategy interface
public interface PaymentStrategy {
    void pay(BigDecimal amount);
}

// Concrete strategies
@Component("creditCard")
public class CreditCardPayment implements PaymentStrategy { ... }

@Component("paypal")
public class PayPalPayment implements PaymentStrategy { ... }

// Context — uses strategy without knowing concrete type
@Service
public class CheckoutService {
    private final Map<String, PaymentStrategy> strategies;
    
    public CheckoutService(Map<String, PaymentStrategy> strategies) {
        this.strategies = strategies; // Spring injects all implementations!
    }
    
    public void checkout(String method, BigDecimal amount) {
        strategies.get(method).pay(amount);
    }
}
```

### Patterns You Should Use WITH Spring

| Pattern | Use Case | Spring Support |
|---|---|---|
| **Repository** | Data access abstraction | `@Repository`, Spring Data |
| **Service Layer** | Business logic encapsulation | `@Service` |
| **DTO (Data Transfer Object)** | API contracts | Jackson serialization |
| **Builder** | Complex object construction | Lombok `@Builder` |
| **Specification** | Dynamic queries | Spring Data `Specification` |
| **Circuit Breaker** | Resilience | Spring Cloud Circuit Breaker |

---

## 13. Best Practices & Anti-Patterns

### 13.1 Best Practices

| # | Practice | Why |
|---|---|---|
| 1 | **Constructor injection everywhere** | Immutability, testability, explicit deps |
| 2 | **Program to interfaces** | Loose coupling, easy mocking, proxy-friendly |
| 3 | **Use `@Configuration` over `@Component` for config** | CGLIB ensures singleton semantics for `@Bean` methods |
| 4 | **Keep controllers thin** | Delegate to service layer; controller = routing + validation |
| 5 | **`@Transactional` at service layer** | Not on repositories (too fine-grained) or controllers (too coarse) |
| 6 | **Global exception handler** | `@RestControllerAdvice` for consistent error responses |
| 7 | **Externalize config** | `@PropertySource` + profiles, not hardcoded values |
| 8 | **Prefer `@EventListener`** | Decouple side effects from main business logic |
| 9 | **Use `@Async` carefully** | Configure thread pool, handle exceptions, propagate context |
| 10 | **Fail fast at startup** | Eager init (default singleton) catches wiring errors early |

### 13.1.1 Layered Architecture Best Practice

```
┌─────────────────────────────────────────────────────────────┐
│  Presentation Layer                                         │
│  @Controller / @RestController                              │
│  • Handle HTTP requests/responses                           │
│  • Input validation                                         │
│  • NO business logic                                        │
├─────────────────────────────────────────────────────────────┤
│  Service Layer                                              │
│  @Service                                                   │
│  • Business logic and rules                                 │
│  • Transaction boundaries (@Transactional)                  │
│  • Orchestrates repositories                                │
├─────────────────────────────────────────────────────────────┤
│  Repository Layer                                           │
│  @Repository                                                │
│  • Data access logic                                        │
│  • NO business logic                                        │
│  • Exception translation                                    │
├─────────────────────────────────────────────────────────────┤
│  Domain Layer                                               │
│  Plain Java classes                                         │
│  • Entity classes                                           │
│  • Value objects                                            │
│  • Domain events                                            │
└─────────────────────────────────────────────────────────────┘
```

### 13.1.2 Configuration Best Practices

```java
// ✅ DO: Split configuration by concern
@Configuration
public class DataSourceConfig { }

@Configuration  
public class SecurityConfig { }

@Configuration
public class CacheConfig { }

@Configuration
@Import({DataSourceConfig.class, SecurityConfig.class, CacheConfig.class})
public class AppConfig { }

// ❌ DON'T: God configuration class with 100+ beans
@Configuration
public class EverythingConfig { 
    // 100 @Bean methods... 
}
```

### 13.2 Common Anti-Patterns

| Anti-Pattern | Problem | Solution |
|---|---|---|
| **Field injection everywhere** | Hidden deps, untestable | Constructor injection |
| **God `@Configuration`** | One class with 50 `@Bean` methods | Split into modular `@Configuration` classes |
| **`@Autowired` on everything** | Over-reliance on Spring; hard to test | Constructor injection (implicit on single constructor) |
| **Catching exceptions in `@Transactional`** | Tx commits despite failure | Rethrow or call `setRollbackOnly()` |
| **Self-invocation with AOP** | Proxy bypassed | Use `@Lookup`, `ObjectProvider`, or extract to separate bean |
| **`@Transactional` on private methods** | Proxy can't intercept | Make method `public` or `protected` |
| **Circular dependencies** | A → B → A | Redesign; use `@Lazy` as last resort |
| **Ignoring bean scope mismatch** | Singleton holding prototype | Use `ObjectProvider` or `@Lookup` |

---

## 14. Interview Questions by Experience Level

### 14.1 Junior (0–2 Years)

**Q1: What is Dependency Injection?**
> DI is a design pattern where an object receives its dependencies from an external source rather than creating them itself. Spring achieves this through its IoC container. DI promotes loose coupling, making code more testable and maintainable.

**Q2: What are the types of DI in Spring?**
> Constructor injection (preferred — immutable, explicit), setter injection (optional deps), and field injection (avoid — hides dependencies, hard to test).

**Q3: What is the difference between `@Component`, `@Service`, `@Repository`, and `@Controller`?**
> All are stereotype annotations extending `@Component`. `@Service` marks business logic (no extra Spring behavior). `@Repository` enables exception translation. `@Controller` marks MVC controllers. They serve as documentation and enable targeted AOP pointcuts.

**Q4: What is a Spring Bean?**
> An object instantiated, assembled, and managed by the Spring IoC container. By default, beans are singletons.

**Q5: What is `@Autowired`? How does Spring resolve it?**
> It tells Spring to inject a dependency. Resolution: first by type, then `@Primary`, then `@Qualifier`, then by parameter/field name.

**Q6: What is the difference between `BeanFactory` and `ApplicationContext`?**
> Both are IoC containers. `ApplicationContext` (used in practice) extends `BeanFactory` with i18n, event publishing, AOP integration, web support, and eager initialization.

**Q7: Explain singleton vs prototype scope.**
> Singleton: one instance per container (default). Prototype: new instance each time. Prototype beans are NOT managed after creation (no `@PreDestroy`).

---

### 14.2 Mid-Level (2–5 Years)

**Q8: Explain the complete bean lifecycle.**
> Instantiation → DI → Aware callbacks → BPP `postProcessBefore` → `@PostConstruct` → BPP `postProcessAfter` → Ready → `@PreDestroy` → Destroy. BPPs are how Spring implements AOP proxies and `@Transactional`.

**Q9: How does `@Transactional` work internally?**
> Spring creates a proxy (JDK or CGLIB) around the bean. When a `@Transactional` method is called through the proxy, an interceptor starts a transaction, delegates to the actual method, and commits or rolls back based on the outcome. This is why self-invocation bypasses transaction management.

**Q10: What is the self-invocation problem? How do you fix it?**
> When a method in a bean calls another method on `this`, it bypasses the proxy. Fix: inject the bean into itself, use `AopContext.currentProxy()`, extract the called method into a separate bean, or use `@Lookup`.

**Q11: Explain transaction propagation. When would you use `REQUIRES_NEW`?**
> Propagation defines how a transaction relates to an existing one. `REQUIRES_NEW` creates an independent transaction, useful for audit logging that must persist even if the outer transaction rolls back.

**Q12: How does `@ControllerAdvice` work?**
> It's a global interceptor for controllers. `@ExceptionHandler` methods in `@ControllerAdvice` catch exceptions from any controller. Used for centralized error handling. Can be scoped to specific packages or annotations.

**Q13: Explain the difference between JDK Proxy and CGLIB Proxy.**
> JDK: interface-based, uses `java.lang.reflect.Proxy`. CGLIB: subclass-based, no interface needed, generates bytecode. Spring Boot defaults to CGLIB. CGLIB cannot proxy `final` classes/methods.

**Q14: What is the difference between `@Bean` and `@Component`?**
> `@Component` is class-level (detected by classpath scanning). `@Bean` is method-level (in `@Configuration` class) for registering third-party objects you can't annotate. `@Bean` gives you full programmatic control over construction.

**Q15: How would you handle circular dependencies?**
> Best: redesign to eliminate the cycle. Workarounds: `@Lazy` on one injection point, setter injection, `ObjectProvider`. Spring 6+ disallows circular deps by default.

---

### 14.3 Senior / Lead (5+ Years)

**Q16: How does `@Configuration` differ from `@Component` internally (full vs lite mode)?**
> `@Configuration` classes are CGLIB-proxied ensuring inter-`@Bean` method calls return the same singleton. `@Component` with `@Bean` methods runs in "lite mode" — each method call creates a new instance. This affects singleton guarantees.

**Q17: Explain BeanPostProcessor vs BeanFactoryPostProcessor.**
> BFPP modifies **bean definitions** (metadata) before instantiation — e.g., `PropertySourcesPlaceholderConfigurer`. BPP modifies **bean instances** after creation — e.g., `AutowiredAnnotationBeanPostProcessor`. Order: BFPP runs first, then beans are created, then BPP processes them.

**Q18: Design a custom `@Retryable` annotation using Spring AOP.**
> Create `@Retryable(maxAttempts, backoff)` annotation. Write an `@Aspect` with `@Around` advice matching `@annotation(retryable)`. In the advice, loop up to `maxAttempts`, catch exceptions, apply backoff via `Thread.sleep()`, and rethrow if all attempts fail. Register the aspect as a `@Component`.

**Q19: How would you implement a custom scope (e.g., tenant-scoped)?**
> Implement `org.springframework.beans.factory.config.Scope` interface: `get()` looks up/creates bean per tenant, `remove()` cleans up, `registerDestructionCallback()` handles lifecycle. Register via `ConfigurableBeanFactory.registerScope("tenant", new TenantScope())`. Use `@Scope("tenant")` on beans.

**Q20: Explain how Spring's event system works with transactions.**
> `@TransactionalEventListener(phase = AFTER_COMMIT)` defers event processing until the transaction commits. Internally, Spring registers `TransactionSynchronization` callbacks. If the tx rolls back, AFTER_COMMIT listeners never fire — preventing side effects from uncommitted data.

**Q21: How does Spring resolve property values (`${…}`) and SpEL (`#{…}`)?**
> `${…}` placeholders are resolved by `PropertySourcesPlaceholderConfigurer` (a BFPP) which replaces them with values from `Environment`. `#{…}` SpEL expressions are evaluated by `StandardBeanExpressionResolver` at bean definition time. They can reference beans (`@beanName`), call methods, and use operators.

**Q22: What strategies do you use to optimize Spring application startup time?**
> 1) Use `@Lazy` for expensive beans not needed at startup. 2) Limit `@ComponentScan` to specific packages. 3) Avoid classpath scanning overhead with explicit `@Import`. 4) Use Spring AOT (Ahead-of-Time) for GraalVM native images. 5) Profile startup with `ApplicationStartup` metrics. 6) Defer non-critical initialization to `@EventListener(ApplicationReadyEvent.class)`.

**Q23: How do you handle distributed transactions across microservices?**
> Avoid 2PC (two-phase commit) in microservices. Use the **Saga pattern** — either choreography (events) or orchestration (central coordinator). Each service has a local transaction and publishes compensating events on failure. Spring supports this via `@TransactionalEventListener` and frameworks like Axon or Temporal.

**Q24: Design a multi-tenant application with Spring.**
> Use a `TenantContext` (ThreadLocal-based). Implement a `HandlerInterceptor` to extract tenant ID from request headers and set it in context. For data isolation: 1) **Schema-per-tenant**: `AbstractRoutingDataSource` switches DataSource per tenant. 2) **Row-level**: Hibernate `@Filter` adds `WHERE tenant_id = ?` globally. 3) **DB-per-tenant**: Dynamic DataSource registry.

---

### 14.4 Scenario-Based Questions (All Levels)

**Q25: Your application is slow to start. How do you diagnose and fix it?**
> 1) Enable `spring.main.lazy-initialization=true` to identify culprits. 2) Use `ApplicationStartup` with `BufferingApplicationStartup` to profile bean creation. 3) Check for: slow `@PostConstruct` methods, excessive classpath scanning, database connections at startup, HTTP calls during init. 4) Solutions: `@Lazy` on slow beans, narrow `@ComponentScan`, use `@EventListener(ApplicationReadyEvent)` for deferred init.

**Q26: A `@Transactional` method isn't rolling back on exception. What could be wrong?**
> Possible causes: 1) Exception is `checked` and not in `rollbackFor` (default rolls back only for unchecked). 2) Exception is caught inside the method. 3) Method is called via `this` (self-invocation). 4) Method is `private`. 5) Transaction manager not configured. 6) Using wrong `DataSource`. Debug by checking proxy creation and exception propagation.

**Q27: How would you implement rate limiting in a Spring application?**
> Options: 1) Custom `HandlerInterceptor` with in-memory counter or Redis. 2) Servlet `Filter` with bucket4j library. 3) Spring Cloud Gateway with built-in rate limiter. 4) AOP aspect with custom `@RateLimited` annotation. Key decisions: scope (per user, IP, API key), storage (memory vs Redis for distributed), algorithm (token bucket, sliding window).

**Q28: You have a memory leak in production. How do you identify if it's Spring-related?**
> 1) Take heap dumps at intervals, analyze with MAT or VisualVM. 2) Check for: prototype beans never garbage collected, caches without eviction, event listeners holding references, thread-local not cleaned. 3) Spring-specific culprits: `@Async` without bounded queue, `ApplicationContext` not closed in tests, circular references preventing GC.

**Q29: How do you ensure your REST API is backward compatible?**
> 1) API versioning (URL path `/v1/`, header, or media type). 2) Never remove fields, only add (use `@JsonIgnoreProperties(ignoreUnknown = true)` on clients). 3) Deprecate before removing. 4) Use consumer-driven contract testing (Spring Cloud Contract, Pact). 5) Document changes. 6) Maintain multiple versions in parallel during transition.

**Q30: Your service needs to call another service. How do you handle failures?**
> 1) **Timeouts**: Configure connection and read timeouts on `RestTemplate`/`WebClient`. 2) **Retries**: Use `spring-retry` with exponential backoff. 3) **Circuit breaker**: Resilience4j to fail fast when service is down. 4) **Fallback**: Graceful degradation with cached or default data. 5) **Bulkhead**: Isolate thread pools per downstream service.

---

## Quick Reference — Annotation Cheat Sheet

| Annotation | Module | Purpose |
|---|---|---|
| `@Component` | Core | Generic managed bean |
| `@Service` | Core | Business-layer bean |
| `@Repository` | Core | Persistence-layer + exception translation |
| `@Controller` | Web | MVC controller (returns views) |
| `@RestController` | Web | REST controller (returns data) |
| `@Configuration` | Core | Declares `@Bean` methods, CGLIB-proxied |
| `@Bean` | Core | Factory method for bean registration |
| `@Autowired` | Core | Dependency injection |
| `@Qualifier` | Core | Disambiguate bean by name |
| `@Primary` | Core | Default bean for a type |
| `@Lazy` | Core | Defer initialization |
| `@Scope` | Core | Bean scope (singleton, prototype, …) |
| `@PostConstruct` | Core | Init callback |
| `@PreDestroy` | Core | Destroy callback |
| `@Value` | Core | Inject property/SpEL value |
| `@PropertySource` | Core | Load properties file |
| `@Profile` | Core | Conditional on active profile |
| `@Conditional` | Core | Conditional bean registration |
| `@Import` | Core | Import configuration class |
| `@ComponentScan` | Core | Configure component scanning |
| `@Transactional` | TX | Transaction demarcation |
| `@EnableTransactionManagement` | TX | Enable annotation-driven TX |
| `@Aspect` | AOP | Declare aspect class |
| `@Before` / `@After` / `@Around` | AOP | Advice types |
| `@Pointcut` | AOP | Reusable pointcut expression |
| `@EnableAspectJAutoProxy` | AOP | Enable AOP proxying |
| `@GetMapping` / `@PostMapping` | Web | HTTP method mapping shortcuts |
| `@RequestBody` | Web | Deserialize request body |
| `@ResponseBody` | Web | Serialize return value |
| `@PathVariable` | Web | Extract URI template variable |
| `@RequestParam` | Web | Extract query parameter |
| `@ExceptionHandler` | Web | Handle specific exceptions |
| `@ControllerAdvice` | Web | Global exception/advice handler |
| `@Valid` | Web | Trigger bean validation |
| `@EnableWebSecurity` | Security | Enable Spring Security |
| `@PreAuthorize` | Security | Method-level authorization |
| `@EnableAsync` | Async | Enable `@Async` support |
| `@Async` | Async | Execute method asynchronously |
| `@EnableScheduling` | Scheduling | Enable `@Scheduled` support |
| `@Scheduled` | Scheduling | Schedule recurring task |
| `@EventListener` | Events | Listen to application events |
| `@TransactionalEventListener` | Events | Transaction-aware listener |

---

> **Next Steps:** See the companion guides (coming soon) for **Spring Boot** and **Microservices** patterns.
