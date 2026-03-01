# Complete Guide to Multithreading in Spring Boot

## Table of Contents
- [1. Introduction](#1-introduction)
  - [1.1 Why Multithreading Matters in Spring Boot](#11-why-multithreading-matters-in-spring-boot)
  - [1.2 How Spring Boot Handles Requests by Default](#12-how-spring-boot-handles-requests-by-default)
  - [1.3 Core Java vs Spring-Managed Concurrency](#13-core-java-vs-spring-managed-concurrency)
- [2. Default Threading Model in Spring Boot](#2-default-threading-model-in-spring-boot)
  - [2.1 Embedded Tomcat Thread Pool Architecture](#21-embedded-tomcat-thread-pool-architecture)
  - [2.2 Request-Per-Thread Model](#22-request-per-thread-model)
  - [2.3 Server Thread Pool Configuration](#23-server-thread-pool-configuration)
  - [2.4 Performance Considerations](#24-performance-considerations)
- [3. Using @Async in Spring Boot](#3-using-async-in-spring-boot)
  - [3.1 Enabling Async Support](#31-enabling-async-support)
  - [3.2 Using @Async Annotation](#32-using-async-annotation)
  - [3.3 Return Types: void vs Future vs CompletableFuture](#33-return-types-void-vs-future-vs-completablefuture)
  - [3.4 Exception Handling in Async Methods](#34-exception-handling-in-async-methods)
  - [3.5 Custom TaskExecutor Configuration](#35-custom-taskexecutor-configuration)
  - [3.6 Complete Working Example](#36-complete-working-example)
- [4. ThreadPoolTaskExecutor Configuration](#4-threadpooltaskexecutor-configuration)
  - [4.1 Understanding Thread Pool Parameters](#41-understanding-thread-pool-parameters)
  - [4.2 Core Pool Size vs Max Pool Size](#42-core-pool-size-vs-max-pool-size)
  - [4.3 Queue Capacity and Behavior](#43-queue-capacity-and-behavior)
  - [4.4 Rejection Policies](#44-rejection-policies)
  - [4.5 Complete Configuration Examples](#45-complete-configuration-examples)
  - [4.6 Production Tuning Guidelines](#46-production-tuning-guidelines)
- [5. Using CompletableFuture in Spring Boot](#5-using-completablefuture-in-spring-boot)
  - [5.1 CompletableFuture Fundamentals](#51-completablefuture-fundamentals)
  - [5.2 Asynchronous Service Calls](#52-asynchronous-service-calls)
  - [5.3 Combining Multiple Async Calls](#53-combining-multiple-async-calls)
  - [5.4 Exception Handling Patterns](#54-exception-handling-patterns)
  - [5.5 Performance Benefits and Patterns](#55-performance-benefits-and-patterns)
- [6. Scheduling Tasks](#6-scheduling-tasks)
  - [6.1 Enabling Scheduling Support](#61-enabling-scheduling-support)
  - [6.2 @Scheduled Annotation](#62-scheduled-annotation)
  - [6.3 Thread Pool for Scheduled Tasks](#63-thread-pool-for-scheduled-tasks)
  - [6.4 Dynamic and Conditional Scheduling](#64-dynamic-and-conditional-scheduling)
  - [6.5 Best Practices and Production Considerations](#65-best-practices-and-production-considerations)
- [7. Spring WebFlux vs Spring MVC](#7-spring-webflux-vs-spring-mvc)
  - [7.1 Thread Model Comparison](#71-thread-model-comparison)
  - [7.2 Event Loop Model](#72-event-loop-model)
  - [7.3 Non-Blocking vs Blocking I/O](#73-non-blocking-vs-blocking-io)
  - [7.4 When to Use WebFlux](#74-when-to-use-webflux)
  - [7.5 Code Comparison and Migration](#75-code-comparison-and-migration)
- [8. Database and Multithreading](#8-database-and-multithreading)
  - [8.1 Connection Pools (HikariCP)](#81-connection-pools-hikaricp)
  - [8.2 Thread Safety in JPA/Hibernate](#82-thread-safety-in-jpahibernate)
  - [8.3 Transaction Boundaries](#83-transaction-boundaries)
  - [8.4 Common Mistakes and Solutions](#84-common-mistakes-and-solutions)
- [9. Concurrency Issues in Spring Boot](#9-concurrency-issues-in-spring-boot)
  - [9.1 Race Conditions in Singleton Beans](#91-race-conditions-in-singleton-beans)
  - [9.2 Thread Safety of Spring Beans](#92-thread-safety-of-spring-beans)
  - [9.3 Using Prototype Scope](#93-using-prototype-scope)
  - [9.4 Using ThreadLocal Safely](#94-using-threadlocal-safely)
- [10. Best Practices in Production](#10-best-practices-in-production)
  - [10.1 Avoid Blocking Calls in Async Threads](#101-avoid-blocking-calls-in-async-threads)
  - [10.2 Proper Executor Sizing](#102-proper-executor-sizing)
  - [10.3 Monitoring Thread Pools](#103-monitoring-thread-pools)
  - [10.4 Handling Backpressure](#104-handling-backpressure)
  - [10.5 Logging and Debugging Async Code](#105-logging-and-debugging-async-code)
  - [10.6 Avoiding Memory Leaks](#106-avoiding-memory-leaks)

---

## 1. Introduction

### 1.1 Why Multithreading Matters in Spring Boot

Multithreading is fundamental to building high-performance, responsive Spring Boot applications. Here's why it matters:

#### **Throughput & Scalability**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    WHY MULTITHREADING MATTERS                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Single-Threaded               vs           Multi-Threaded                  │
│  ─────────────────                          ──────────────────              │
│                                                                             │
│  Request 1 ─────────────────►               Request 1 ─────►                │
│            [    Processing    ]             Request 2 ─────►   Parallel     │
│  Request 2 ─────────────────►               Request 3 ─────►   Execution    │
│            [    Processing    ]             Request 4 ─────►                │
│  Request 3 ─────────────────►                                               │
│            [    Processing    ]                                             │
│                                                                             │
│  Total: 3 × Processing Time                 Total: Max(Processing Times)    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Key Benefits**

| Benefit | Description | Real-World Impact |
|---------|-------------|-------------------|
| **Improved Throughput** | Handle multiple requests simultaneously | Serve 1000s of concurrent users |
| **Better Resource Utilization** | Keep CPU busy while waiting on I/O | Maximize server efficiency |
| **Reduced Latency** | Parallel processing of independent tasks | Faster API response times |
| **Non-Blocking Operations** | I/O operations don't block other requests | Better user experience |
| **Background Processing** | Offload heavy tasks to background threads | Responsive UI/API |

#### **Production Use Cases**

```java
// Use Case 1: E-Commerce Order Processing
@Service
public class OrderService {
    
    @Async  // Process order confirmation email in background
    public void sendOrderConfirmation(Order order) {
        emailService.send(order.getCustomerEmail(), buildEmailContent(order));
    }
    
    // Use Case 2: Parallel Data Aggregation
    public DashboardData getDashboard(String userId) {
        CompletableFuture<UserStats> stats = asyncService.getUserStats(userId);
        CompletableFuture<List<Order>> orders = asyncService.getRecentOrders(userId);
        CompletableFuture<List<Notification>> notifications = asyncService.getNotifications(userId);
        
        // All three calls execute in parallel
        CompletableFuture.allOf(stats, orders, notifications).join();
        
        return new DashboardData(stats.get(), orders.get(), notifications.get());
    }
}
```

---

### 1.2 How Spring Boot Handles Requests by Default

Spring Boot uses an embedded servlet container (Tomcat by default) that manages a thread pool for handling HTTP requests.

#### **Default Request Handling Flow**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SPRING BOOT REQUEST HANDLING                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│    Client                  Tomcat                    Spring Container       │
│    ──────                  ──────                    ─────────────────      │
│                                                                             │
│  ┌─────────┐           ┌─────────────┐            ┌──────────────────┐     │
│  │ Request │──────────►│Thread Pool  │───────────►│DispatcherServlet │     │
│  └─────────┘           │  (200 max)  │            └────────┬─────────┘     │
│                        │             │                     │               │
│                        │ ┌─────────┐ │                     ▼               │
│                        │ │Thread-1 │─┼──►┌────────────────────────────┐    │
│                        │ ├─────────┤ │   │  Filter Chain              │    │
│                        │ │Thread-2 │ │   │    ↓                       │    │
│                        │ ├─────────┤ │   │  Controller → Service      │    │
│                        │ │Thread-3 │ │   │    ↓                       │    │
│                        │ ├─────────┤ │   │  Repository → Database     │    │
│                        │ │  ...    │ │   │    ↓                       │    │
│                        │ ├─────────┤ │   │  Response                  │    │
│                        │ │Thread-N │ │   └────────────────────────────┘    │
│                        │ └─────────┘ │                                      │
│                        └─────────────┘                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Default Tomcat Thread Pool Settings**

```yaml
# Default values (as of Spring Boot 3.x)
server:
  tomcat:
    threads:
      max: 200          # Maximum worker threads
      min-spare: 10     # Minimum idle threads
    max-connections: 8192  # Maximum connections
    accept-count: 100      # Queue size when all threads are busy
```

#### **What Happens Per Request**

```java
// Each HTTP request gets its own thread from the pool
@RestController
public class UserController {
    
    @GetMapping("/users/{id}")
    public User getUser(@PathVariable Long id) {
        // This entire method executes on a single Tomcat thread
        // Thread name: http-nio-8080-exec-{N}
        
        String threadName = Thread.currentThread().getName();
        System.out.println("Handling request on: " + threadName);
        // Output: "Handling request on: http-nio-8080-exec-3"
        
        return userService.findById(id);
    }
}
```

#### **Thread Pool States**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        TOMCAT THREAD POOL STATES                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Low Traffic                    Medium Traffic          High Traffic        │
│  ───────────                    ──────────────          ────────────        │
│                                                                             │
│  ┌──────────┐                  ┌──────────┐            ┌──────────┐         │
│  │ Active:3 │                  │Active:50 │            │Active:200│         │
│  │ Idle: 7  │                  │ Idle: 10 │            │ Idle: 0  │         │
│  │ Queue: 0 │                  │ Queue: 0 │            │Queue: 100│         │
│  └──────────┘                  └──────────┘            └──────────┘         │
│                                                                             │
│  Threads scale up as needed     Operating normally      At capacity!        │
│                                                         New requests wait   │
│                                                         in accept queue     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 1.3 Core Java vs Spring-Managed Concurrency

Understanding the differences helps you choose the right approach for different scenarios.

#### **Comparison Overview**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│              CORE JAVA vs SPRING-MANAGED CONCURRENCY                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Core Java Threading                   Spring-Managed Concurrency           │
│  ────────────────────                  ──────────────────────────           │
│                                                                             │
│  • Manual thread management            • Declarative with @Async            │
│  • Direct ExecutorService usage        • Spring-managed TaskExecutor        │
│  • Manual exception handling           • Built-in exception propagation     │
│  • No transaction support              • Transaction propagation aware      │
│  • No Spring context access            • Full Spring context available      │
│  • Lower-level control                 • Higher-level abstraction           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Code Comparison**

**Core Java Approach:**
```java
@Service
public class ReportService {
    
    private final ExecutorService executor = Executors.newFixedThreadPool(10);
    
    public void generateReports(List<ReportRequest> requests) {
        List<Future<?>> futures = new ArrayList<>();
        
        for (ReportRequest request : requests) {
            Future<?> future = executor.submit(() -> {
                try {
                    // Manual exception handling required
                    generateReport(request);
                } catch (Exception e) {
                    // Must handle manually
                    log.error("Report generation failed", e);
                }
            });
            futures.add(future);
        }
        
        // Manual waiting for completion
        for (Future<?> future : futures) {
            try {
                future.get(30, TimeUnit.SECONDS);
            } catch (TimeoutException | ExecutionException | InterruptedException e) {
                // Handle exceptions
            }
        }
    }
    
    // Must manually shutdown executor
    @PreDestroy
    public void cleanup() {
        executor.shutdown();
        try {
            if (!executor.awaitTermination(60, TimeUnit.SECONDS)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
        }
    }
}
```

**Spring-Managed Approach:**
```java
@Service
@RequiredArgsConstructor
public class ReportService {
    
    private final AsyncReportGenerator asyncReportGenerator;
    
    public void generateReports(List<ReportRequest> requests) {
        List<CompletableFuture<Void>> futures = requests.stream()
            .map(asyncReportGenerator::generateReportAsync)
            .collect(Collectors.toList());
        
        // Wait for all to complete
        CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).join();
    }
}

@Service
public class AsyncReportGenerator {
    
    @Async("reportTaskExecutor")  // Uses Spring-managed executor
    public CompletableFuture<Void> generateReportAsync(ReportRequest request) {
        // Automatic exception handling via AsyncUncaughtExceptionHandler
        generateReport(request);
        return CompletableFuture.completedFuture(null);
    }
}

@Configuration
@EnableAsync
public class AsyncConfig implements AsyncConfigurer {
    
    @Bean("reportTaskExecutor")
    public TaskExecutor reportTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("report-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
        // Spring handles lifecycle automatically!
    }
}
```

#### **Feature Comparison Table**

| Feature | Core Java | Spring-Managed |
|---------|-----------|----------------|
| **Thread Pool Lifecycle** | Manual `shutdown()` required | Auto-managed by Spring |
| **Exception Handling** | Manual try-catch everywhere | `AsyncUncaughtExceptionHandler` |
| **Transaction Support** | Not propagated | Configurable propagation |
| **Spring Beans Access** | Manual injection/lookup | Full DI support |
| **Configuration** | Code-based | Properties/Code-based |
| **Monitoring** | Manual metrics | Actuator integration |
| **Shutdown** | Manual graceful shutdown | Spring graceful shutdown |
| **Testing** | Complex mocking | `@Async` can be disabled |

#### **When to Use Which**

```java
// USE CORE JAVA WHEN:
// 1. Fine-grained control needed
// 2. Non-Spring components
// 3. Complex completion strategies

CompletableFuture.supplyAsync(() -> fetchDataA(), executorA)
    .thenCombine(
        CompletableFuture.supplyAsync(() -> fetchDataB(), executorB),
        (a, b) -> merge(a, b)
    )
    .thenApplyAsync(this::transform, executorC)
    .exceptionally(this::handleError);

// USE SPRING @ASYNC WHEN:
// 1. Simple fire-and-forget
// 2. Background processing
// 3. Need Spring context (transactions, security)

@Async
@Transactional  // Transaction works automatically
public void processOrder(Order order) {
    orderRepository.save(order);
    // If this fails, transaction rolls back
}
```

---

## 2. Default Threading Model in Spring Boot

### 2.1 Embedded Tomcat Thread Pool Architecture

Spring Boot's embedded Tomcat uses a sophisticated thread pool architecture for handling HTTP requests.

#### **Architecture Overview**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TOMCAT THREAD POOL ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                              ┌───────────────────┐                          │
│                              │   NIO Connector   │                          │
│                              │   (Port 8080)     │                          │
│                              └─────────┬─────────┘                          │
│                                        │                                    │
│                                        ▼                                    │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │                         ACCEPTOR THREAD                            │     │
│  │   • Single thread accepting new connections                        │     │
│  │   • Hands off to Poller                                            │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                        │                                    │
│                                        ▼                                    │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │                          POLLER THREAD                             │     │
│  │   • Uses NIO Selector for non-blocking I/O                         │     │
│  │   • Monitors multiple connections efficiently                      │     │
│  │   • Detects when data is ready to read                             │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                        │                                    │
│                                        ▼                                    │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │                    WORKER THREAD POOL                              │     │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │     │
│  │  │Worker-1 │ │Worker-2 │ │Worker-3 │ │  ...    │ │Worker-N │       │     │
│  │  │ (busy)  │ │ (idle)  │ │ (busy)  │ │         │ │ (idle)  │       │     │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘       │     │
│  │                                                                    │     │
│  │  Core Size: 10 (min-spare)     Max Size: 200 (max threads)         │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Thread Pool Components**

```java
// Understanding Tomcat's internal components

/**
 * 1. ACCEPTOR - Accepts incoming connections
 * Thread name: http-nio-8080-Acceptor
 */

/**
 * 2. POLLER - Monitors connection readiness (NIO)
 * Thread name: http-nio-8080-Poller
 * Uses java.nio.channels.Selector
 */

/**
 * 3. WORKER POOL - Processes actual requests
 * Thread name: http-nio-8080-exec-{N}
 * These are the threads your code runs on
 */

@RestController
public class ThreadInfoController {
    
    @GetMapping("/thread-info")
    public Map<String, String> getThreadInfo() {
        Thread current = Thread.currentThread();
        return Map.of(
            "threadName", current.getName(),        // http-nio-8080-exec-1
            "threadGroup", current.getThreadGroup().getName(),
            "threadId", String.valueOf(current.getId())
        );
    }
    
    @GetMapping("/all-tomcat-threads")
    public List<String> getAllTomcatThreads() {
        return Thread.getAllStackTraces().keySet().stream()
            .map(Thread::getName)
            .filter(name -> name.startsWith("http-nio"))
            .sorted()
            .collect(Collectors.toList());
        /*
         * Returns something like:
         * - http-nio-8080-Acceptor
         * - http-nio-8080-Poller
         * - http-nio-8080-exec-1
         * - http-nio-8080-exec-2
         * - ...
         * - http-nio-8080-exec-10
         */
    }
}
```

#### **Connection vs Thread Relationship**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONNECTION vs THREAD MANAGEMENT                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  max-connections: 8192                                                      │
│  ─────────────────────                                                      │
│  • Maximum concurrent TCP connections Tomcat will accept                    │
│  • Connections waiting for worker thread are held by Poller                 │
│                                                                             │
│  threads.max: 200                                                           │
│  ────────────────                                                           │
│  • Maximum worker threads processing requests                               │
│  • Each actively processing request needs one thread                        │
│                                                                             │
│  accept-count: 100                                                          │
│  ─────────────────                                                          │
│  • OS-level TCP accept queue size                                           │
│  • When max-connections reached, new connections wait here                  │
│                                                                             │
│  ┌─────────┐    ┌─────────────────┐    ┌────────────────┐    ┌──────────┐   │
│  │ Accept  │───►│  8192 Active    │───►│  200 Worker    │───►│ Response │   │
│  │ Queue   │    │  Connections    │    │   Threads      │    │          │   │
│  │ (100)   │    │  (Poller mgmt)  │    │  (Processing)  │    │          │   │
│  └─────────┘    └─────────────────┘    └────────────────┘    └──────────┘   │
│                                                                             │
│  If accept-count       If max-connections      If threads.max               │
│  is full, connection   is full, goes to        is full, waits               │
│  refused!              accept-count queue      for thread                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 2.2 Request-Per-Thread Model

The request-per-thread model is the traditional servlet model where each HTTP request is handled by a dedicated thread throughout its lifecycle.

#### **How It Works**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      REQUEST-PER-THREAD MODEL                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Timeline for Request Processing:                                           │
│                                                                             │
│  Thread: http-nio-8080-exec-5                                               │
│  ─────────────────────────────────────────────────────────────────────►    │
│  │                                                                          │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       │
│  ├─►│ Receive │─►│ Filter  │─►│Controller│─►│ Service │─►│ Send    │       │
│  │  │ Request │  │ Chain   │  │         │  │+DB Call │  │ Response│       │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘       │
│  │  [   5ms  ]   [  10ms   ]  [   5ms   ]  [  100ms  ]  [  5ms   ]        │
│  │                                                                          │
│  │  Total: 125ms - Thread is BLOCKED entire time                           │
│  │                                                                          │
│  └─────────────────────────────────────────────────────────────────────►    │
│                                                                             │
│  The Problem with I/O-Bound Operations:                                     │
│  ────────────────────────────────────────                                   │
│  Thread-5:  [Working][████ Waiting for DB ████][Working]                    │
│                       ↑                                                     │
│                       └── Thread is idle but cannot serve other requests   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Code Example: Blocking Operations**

```java
@RestController
@RequiredArgsConstructor
public class OrderController {
    
    private final OrderService orderService;
    private final PaymentService paymentService;
    private final NotificationService notificationService;
    
    /**
     * BLOCKING APPROACH (Request-Per-Thread)
     * Total time: sum of all operations
     * Thread blocked entire time
     */
    @GetMapping("/orders/{id}/details")
    public OrderDetails getOrderDetailsBlocking(@PathVariable Long id) {
        // Thread: http-nio-8080-exec-1 handles everything
        
        Order order = orderService.findById(id);           // 50ms - DB call
        Payment payment = paymentService.getForOrder(id);   // 100ms - External API
        List<Notification> notifs = notificationService.getForOrder(id); // 30ms - DB
        
        // Total: 180ms - Thread blocked entire time
        return new OrderDetails(order, payment, notifs);
    }
    
    /**
     * OPTIMIZED APPROACH (Parallel Execution)
     * Total time: max of all operations
     * Uses additional threads from custom pool
     */
    @GetMapping("/orders/{id}/details-parallel")
    public OrderDetails getOrderDetailsParallel(@PathVariable Long id) {
        // Main thread: http-nio-8080-exec-1
        
        CompletableFuture<Order> orderFuture = CompletableFuture
            .supplyAsync(() -> orderService.findById(id));
        CompletableFuture<Payment> paymentFuture = CompletableFuture
            .supplyAsync(() -> paymentService.getForOrder(id));
        CompletableFuture<List<Notification>> notifsFuture = CompletableFuture
            .supplyAsync(() -> notificationService.getForOrder(id));
        
        CompletableFuture.allOf(orderFuture, paymentFuture, notifsFuture).join();
        
        // Total: ~100ms (slowest operation)
        return new OrderDetails(
            orderFuture.join(),
            paymentFuture.join(),
            notifsFuture.join()
        );
    }
}
```

#### **Thread Utilization Visualization**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THREAD UTILIZATION COMPARISON                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BLOCKING MODEL (Sequential):                                               │
│  ─────────────────────────────                                              │
│  Time:    0ms      50ms     150ms    180ms                                  │
│            │        │         │        │                                    │
│  Thread-1: [═══DB═══|════PAYMENT API════|══DB══]                            │
│            Processing  │         Waiting       │                            │
│                        └─────────100ms─────────┘                            │
│                                                                             │
│  1 Request = 1 Thread blocked for 180ms                                     │
│  200 threads can handle: 200 concurrent requests                            │
│  Throughput: ~1,111 requests/second                                         │
│                                                                             │
│  PARALLEL MODEL (Async):                                                    │
│  ──────────────────────────                                                 │
│  Time:    0ms               100ms                                           │
│            │                  │                                             │
│  Thread-1: [═══Coordinate════]                                              │
│  Thread-2:   [═══DB 50ms═══]                                                │
│  Thread-3:   [═════PAYMENT API 100ms═════]                                  │
│  Thread-4:   [══DB 30ms══]                                                  │
│                                                                             │
│  1 Request = Main thread free after ~100ms                                  │
│  Work distributed across thread pools                                       │
│  Throughput: ~2,000 requests/second (improved)                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **ThreadLocal in Request-Per-Thread Model**

```java
/**
 * ThreadLocal works perfectly in request-per-thread model
 * because same thread handles entire request lifecycle
 */
@Component
public class RequestContext {
    
    private static final ThreadLocal<RequestInfo> context = new ThreadLocal<>();
    
    public void setRequestInfo(RequestInfo info) {
        context.set(info);
    }
    
    public RequestInfo getRequestInfo() {
        return context.get();
    }
    
    public void clear() {
        context.remove();  // Important: prevent memory leaks
    }
}

// Filter to set context
@Component
@Order(1)
public class RequestContextFilter extends OncePerRequestFilter {
    
    @Autowired
    private RequestContext requestContext;
    
    @Override
    protected void doFilterInternal(HttpServletRequest request, 
                                    HttpServletResponse response, 
                                    FilterChain chain) throws ServletException, IOException {
        try {
            RequestInfo info = new RequestInfo(
                UUID.randomUUID().toString(),
                request.getRemoteAddr(),
                Instant.now()
            );
            requestContext.setRequestInfo(info);
            chain.doFilter(request, response);
        } finally {
            requestContext.clear();  // Always cleanup
        }
    }
}

// Works throughout request processing
@Service
public class AuditService {
    
    @Autowired
    private RequestContext requestContext;
    
    public void logAction(String action) {
        RequestInfo info = requestContext.getRequestInfo();
        // info is available - same thread as filter
        log.info("Action: {} by request: {}", action, info.getRequestId());
    }
}
```

---

### 2.3 Server Thread Pool Configuration

#### **Complete Configuration Reference**

```yaml
# application.yml - Comprehensive Thread Pool Configuration
server:
  port: 8080
  
  tomcat:
    # THREAD POOL SETTINGS
    threads:
      max: 200              # Maximum worker threads (default: 200)
      min-spare: 10         # Minimum idle threads (default: 10)
    
    # CONNECTION SETTINGS  
    max-connections: 8192   # Maximum connections (default: 8192)
    accept-count: 100       # Accept queue size (default: 100)
    
    # TIMEOUT SETTINGS
    connection-timeout: 20000   # Connection timeout in ms (default: 20000)
    keep-alive-timeout: 20000   # Keep-alive timeout in ms
    max-keep-alive-requests: 100  # Max requests per keep-alive connection
    
    # ADDITIONAL SETTINGS
    max-http-form-post-size: 2MB  # Max form POST size
    max-swallow-size: 2MB         # Max request body to swallow on error
    
    # URI ENCODING
    uri-encoding: UTF-8
    
    # ACCESS LOG (useful for debugging)
    accesslog:
      enabled: true
      directory: logs
      pattern: "%t %a %r %s %b %D"  # %D = processing time in ms
```

```properties
# application.properties - Same configuration
server.port=8080
server.tomcat.threads.max=200
server.tomcat.threads.min-spare=10
server.tomcat.max-connections=8192
server.tomcat.accept-count=100
server.tomcat.connection-timeout=20000
```

#### **Environment-Specific Configurations**

```yaml
# application-dev.yml
server:
  tomcat:
    threads:
      max: 50        # Lower for development
      min-spare: 5

---
# application-prod.yml
server:
  tomcat:
    threads:
      max: 400       # Higher for production
      min-spare: 50
    max-connections: 10000
    accept-count: 200
```

#### **Programmatic Configuration**

```java
@Configuration
public class TomcatConfig {
    
    /**
     * Programmatic Tomcat customization
     * Use when you need dynamic configuration based on environment
     */
    @Bean
    public WebServerFactoryCustomizer<TomcatServletWebServerFactory> tomcatCustomizer() {
        return factory -> {
            factory.addConnectorCustomizers(connector -> {
                // Access the protocol handler
                if (connector.getProtocolHandler() instanceof AbstractProtocol<?> protocol) {
                    // Thread pool settings
                    protocol.setMaxThreads(getMaxThreads());
                    protocol.setMinSpareThreads(getMinSpareThreads());
                    
                    // Connection settings
                    protocol.setMaxConnections(8192);
                    protocol.setAcceptCount(100);
                    
                    // Timeout settings
                    protocol.setConnectionTimeout(20000);
                    protocol.setKeepAliveTimeout(20000);
                }
            });
        };
    }
    
    private int getMaxThreads() {
        // Dynamic calculation based on CPU cores
        int cpuCores = Runtime.getRuntime().availableProcessors();
        
        // For I/O-bound apps: more threads than CPU cores
        // Typical formula: threads = cores * (1 + waitTime/serviceTime)
        // Assuming waitTime/serviceTime ratio of ~10 for typical web apps
        int calculatedThreads = cpuCores * 10;
        
        // Cap at reasonable maximum
        return Math.min(calculatedThreads, 500);
    }
    
    private int getMinSpareThreads() {
        return getMaxThreads() / 10;  // 10% of max
    }
}
```

#### **Different Embedded Server Configurations**

```yaml
# For Jetty (add spring-boot-starter-jetty dependency)
server:
  jetty:
    threads:
      max: 200
      min: 8
      idle-timeout: 60000ms
    max-http-form-post-size: 2MB

# For Undertow (add spring-boot-starter-undertow dependency)
server:
  undertow:
    threads:
      io: 4              # I/O threads (typically = CPU cores)
      worker: 64         # Worker threads (for blocking operations)
    buffer-size: 1024
    direct-buffers: true
```

#### **Switching Embedded Servers**

```xml
<!-- pom.xml - Switch from Tomcat to Undertow -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
    <exclusions>
        <exclusion>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-tomcat</artifactId>
        </exclusion>
    </exclusions>
</dependency>
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-undertow</artifactId>
</dependency>
```

---

### 2.4 Performance Considerations

#### **Thread Pool Sizing Guidelines**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THREAD POOL SIZING GUIDELINES                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Formula for Optimal Thread Count:                                          │
│  ─────────────────────────────────                                          │
│                                                                             │
│  Threads = Number of CPUs × Target CPU Utilization × (1 + Wait/Service)     │
│                                                                             │
│  Where:                                                                     │
│  • Number of CPUs = Runtime.getRuntime().availableProcessors()              │
│  • Target CPU Utilization = 0.5 to 1.0 (typically 0.8)                      │
│  • Wait Time = Time spent waiting for I/O (DB, HTTP calls)                  │
│  • Service Time = Time spent actually computing                             │
│                                                                             │
│  EXAMPLES:                                                                  │
│  ─────────                                                                  │
│                                                                             │
│  CPU-Bound Application (calculations, in-memory processing):                │
│  • Wait/Service ratio: ~0                                                   │
│  • 8 CPUs × 0.8 × (1 + 0) = 6-8 threads                                    │
│  • More threads = context switching overhead                                │
│                                                                             │
│  I/O-Bound Application (typical web app with DB):                           │
│  • Wait/Service ratio: ~10 (90% waiting, 10% computing)                     │
│  • 8 CPUs × 0.8 × (1 + 10) = ~70 threads                                   │
│                                                                             │
│  Mixed Application (some computation + I/O):                                │
│  • Wait/Service ratio: ~2-5                                                 │
│  • 8 CPUs × 0.8 × (1 + 3) = ~25-35 threads                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Performance Tuning Code Examples**

```java
@Configuration
public class PerformanceTuningConfig {
    
    @Value("${app.performance.profile:balanced}")
    private String performanceProfile;
    
    @Bean
    public WebServerFactoryCustomizer<TomcatServletWebServerFactory> 
           performanceCustomizer() {
        return factory -> {
            factory.addConnectorCustomizers(connector -> {
                if (connector.getProtocolHandler() instanceof AbstractProtocol<?> protocol) {
                    switch (performanceProfile) {
                        case "cpu-bound" -> configureCpuBound(protocol);
                        case "io-bound" -> configureIoBound(protocol);
                        case "high-throughput" -> configureHighThroughput(protocol);
                        default -> configureBalanced(protocol);
                    }
                }
            });
        };
    }
    
    private void configureCpuBound(AbstractProtocol<?> protocol) {
        int cores = Runtime.getRuntime().availableProcessors();
        protocol.setMaxThreads(cores + 1);  // Minimal threads
        protocol.setMinSpareThreads(cores);
        log.info("CPU-bound config: {} threads", cores + 1);
    }
    
    private void configureIoBound(AbstractProtocol<?> protocol) {
        int cores = Runtime.getRuntime().availableProcessors();
        int threads = cores * 10;  // Higher ratio for I/O waiting
        protocol.setMaxThreads(Math.min(threads, 500));
        protocol.setMinSpareThreads(cores * 2);
        log.info("I/O-bound config: {} threads", threads);
    }
    
    private void configureHighThroughput(AbstractProtocol<?> protocol) {
        protocol.setMaxThreads(500);
        protocol.setMinSpareThreads(100);
        protocol.setMaxConnections(20000);
        protocol.setAcceptCount(500);
        log.info("High-throughput config: 500 threads, 20K connections");
    }
    
    private void configureBalanced(AbstractProtocol<?> protocol) {
        protocol.setMaxThreads(200);
        protocol.setMinSpareThreads(20);
        protocol.setMaxConnections(8192);
        protocol.setAcceptCount(100);
        log.info("Balanced config: default settings");
    }
}
```

#### **Monitoring Thread Pool Health**

```java
@Component
@RequiredArgsConstructor
public class ThreadPoolMetrics {
    
    private final MeterRegistry meterRegistry;
    
    @EventListener(ApplicationReadyEvent.class)
    public void registerMetrics() {
        // Get Tomcat's thread pool via JMX
        Gauge.builder("tomcat.threads.active", this::getActiveThreads)
            .description("Active Tomcat threads")
            .register(meterRegistry);
            
        Gauge.builder("tomcat.threads.max", this::getMaxThreads)
            .description("Maximum Tomcat threads")
            .register(meterRegistry);
            
        Gauge.builder("tomcat.connections.active", this::getActiveConnections)
            .description("Active connections")
            .register(meterRegistry);
    }
    
    private double getActiveThreads() {
        try {
            ObjectName name = new ObjectName(
                "Tomcat:type=ThreadPool,name=\"http-nio-8080\"");
            MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();
            return ((Integer) mbs.getAttribute(name, "currentThreadsBusy")).doubleValue();
        } catch (Exception e) {
            return -1;
        }
    }
    
    // Similar methods for getMaxThreads(), getActiveConnections()...
}
```

#### **Common Performance Problems and Solutions**

```java
/**
 * PROBLEM 1: Thread Pool Exhaustion
 * Symptoms: Requests timing out, 503 errors
 */
@RestController
public class ProblematicController {
    
    // BAD: Long-running operation blocking Tomcat thread
    @GetMapping("/bad/report")
    public Report generateReport() {
        return reportService.generateLargeReport();  // Takes 30 seconds
        // Blocks 1 of 200 threads for 30 seconds!
    }
    
    // GOOD: Offload to background thread
    @GetMapping("/good/report")
    public ResponseEntity<Void> generateReportAsync() {
        String reportId = UUID.randomUUID().toString();
        reportService.generateLargeReportAsync(reportId);  // Non-blocking
        return ResponseEntity.accepted()
            .header("Location", "/reports/" + reportId)
            .build();
    }
}

/**
 * PROBLEM 2: Connection Pool Mismatch
 * DB pool too small compared to Tomcat threads
 */
@Configuration
public class DataSourceConfig {
    
    @Bean
    @ConfigurationProperties("spring.datasource.hikari")
    public HikariConfig hikariConfig() {
        HikariConfig config = new HikariConfig();
        
        // Tomcat has 200 threads, but DB pool only has 10
        // 190 threads will wait for DB connections!
        
        // RULE: DB Pool Size ≈ Tomcat Threads × %DB-bound-requests
        // If 50% of requests hit DB: 200 × 0.5 = 100 connections
        config.setMaximumPoolSize(100);
        config.setMinimumIdle(20);
        config.setConnectionTimeout(30000);
        
        return config;
    }
}

/**
 * PROBLEM 3: Slow Consumers
 * Client reads response slowly, holding thread
 */
@Configuration  
public class SlowConsumerConfig {
    
    @Bean
    public WebServerFactoryCustomizer<TomcatServletWebServerFactory> 
           asyncTimeoutCustomizer() {
        return factory -> factory.addConnectorCustomizers(connector -> {
            // Set socket timeout for slow clients
            connector.setProperty("socket.soTimeout", "60000");
            
            // Use async writes with NIO
            connector.setProperty("socket.appWriteTimeout", "30000");
        });
    }
}
```

#### **Quick Reference: Configuration by Use Case**

| Use Case | max-threads | min-spare | max-connections | accept-count |
|----------|------------|-----------|-----------------|--------------|
| **Small App** (low traffic) | 50 | 5 | 2000 | 50 |
| **Standard API** | 200 | 20 | 8192 | 100 |
| **High Traffic API** | 400 | 50 | 15000 | 200 |
| **Microservice** (fast, I/O) | 100 | 20 | 5000 | 100 |
| **File Upload Service** | 50 | 10 | 500 | 50 |
| **WebSocket Server** | 200 | 50 | 20000 | 500 |

#### **Health Check Endpoint for Thread Pool**

```java
@RestController
@RequestMapping("/actuator/health")
public class ThreadPoolHealthController {
    
    @GetMapping("/threadpool")
    public Map<String, Object> getThreadPoolHealth() {
        Map<String, Object> health = new HashMap<>();
        
        try {
            MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();
            ObjectName name = new ObjectName(
                "Tomcat:type=ThreadPool,name=\"http-nio-8080\"");
            
            int maxThreads = (Integer) mbs.getAttribute(name, "maxThreads");
            int currentThreads = (Integer) mbs.getAttribute(name, "currentThreadCount");
            int busyThreads = (Integer) mbs.getAttribute(name, "currentThreadsBusy");
            int connectionCount = (Integer) mbs.getAttribute(name, "connectionCount");
            
            double utilizationPercent = (busyThreads * 100.0) / maxThreads;
            
            health.put("status", utilizationPercent < 80 ? "UP" : "WARN");
            health.put("maxThreads", maxThreads);
            health.put("currentThreads", currentThreads);
            health.put("busyThreads", busyThreads);
            health.put("idleThreads", currentThreads - busyThreads);
            health.put("connectionCount", connectionCount);
            health.put("utilizationPercent", String.format("%.2f%%", utilizationPercent));
            
            // Warnings
            if (utilizationPercent > 90) {
                health.put("warning", "Thread pool near exhaustion!");
            }
            
        } catch (Exception e) {
            health.put("status", "UNKNOWN");
            health.put("error", e.getMessage());
        }
        
        return health;
    }
}
```

---

## Summary

| Topic | Key Points |
|-------|------------|
| **Why Multithreading** | Throughput, resource utilization, responsiveness |
| **Default Model** | Tomcat thread pool, 200 max threads, request-per-thread |
| **Core Java vs Spring** | Manual vs declarative, lifecycle management, Spring context |
| **Thread Pool Config** | `server.tomcat.threads.*` properties, programmatic customization |
| **Performance** | Size based on CPU cores and I/O ratio, monitor utilization |

---

## 3. Using @Async in Spring Boot

The `@Async` annotation is Spring's declarative approach to asynchronous method execution. It allows you to run methods in a separate thread without manually managing thread creation and lifecycle.

### 3.1 Enabling Async Support

#### **Basic Setup**

```java
@Configuration
@EnableAsync  // This annotation enables Spring's async processing
public class AsyncConfig {
    // Configuration beans go here
}
```

#### **Alternative: Enable on Main Application Class**

```java
@SpringBootApplication
@EnableAsync
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

#### **What @EnableAsync Does**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    @EnableAsync PROCESSING                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Without @EnableAsync:                                                      │
│  ─────────────────────                                                      │
│                                                                             │
│  @Async                                                                     │
│  public void sendEmail() { ... }    ──►  Runs on SAME thread (synchronous) │
│                                                                             │
│  With @EnableAsync:                                                         │
│  ──────────────────                                                         │
│                                                                             │
│  Spring creates a PROXY around the bean:                                    │
│                                                                             │
│  ┌─────────────────┐      ┌──────────────────────────┐                     │
│  │  Caller         │─────►│  AsyncProxy              │                     │
│  │  (Main Thread)  │      │  • Intercepts @Async     │                     │
│  └─────────────────┘      │  • Submits to Executor   │                     │
│         │                 │  • Returns immediately   │                     │
│         │                 └──────────────┬───────────┘                     │
│         │                                │                                  │
│         ▼                                ▼                                  │
│  Continues execution           ┌─────────────────┐                         │
│  immediately                   │  TaskExecutor   │                         │
│                                │  Thread Pool    │                         │
│                                │  • Executes     │                         │
│                                │    sendEmail()  │                         │
│                                └─────────────────┘                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **@EnableAsync Attributes**

```java
@Configuration
@EnableAsync(
    annotation = Async.class,              // Which annotation to look for
    proxyTargetClass = false,              // true = CGLIB, false = JDK proxy
    mode = AdviceMode.PROXY                // PROXY or ASPECTJ
)
public class AsyncConfig {
    // PROXY mode (default): AOP proxy-based, method calls from same class won't work
    // ASPECTJ mode: Compile-time weaving, works for internal calls too
}
```

---

### 3.2 Using @Async Annotation

#### **Basic Usage**

```java
@Service
@Slf4j
public class NotificationService {
    
    /**
     * Simple async method - fire and forget
     * Caller doesn't wait for completion
     */
    @Async
    public void sendEmailNotification(String email, String message) {
        log.info("Sending email on thread: {}", Thread.currentThread().getName());
        // Simulate email sending
        try {
            Thread.sleep(3000);  // 3 second delay
            emailClient.send(email, message);
            log.info("Email sent successfully");
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
    
    /**
     * Async with specific executor
     * Uses the bean named "emailTaskExecutor"
     */
    @Async("emailTaskExecutor")
    public void sendBulkEmails(List<String> emails, String message) {
        emails.forEach(email -> emailClient.send(email, message));
    }
}
```

#### **Calling Async Methods**

```java
@RestController
@RequiredArgsConstructor
public class OrderController {
    
    private final OrderService orderService;
    private final NotificationService notificationService;
    
    @PostMapping("/orders")
    public ResponseEntity<Order> createOrder(@RequestBody OrderRequest request) {
        // Main thread: http-nio-8080-exec-1
        log.info("Creating order on thread: {}", Thread.currentThread().getName());
        
        Order order = orderService.create(request);
        
        // This returns IMMEDIATELY - doesn't wait for email
        notificationService.sendEmailNotification(
            request.getCustomerEmail(),
            "Order " + order.getId() + " confirmed!"
        );
        
        // Response sent while email is being processed in background
        return ResponseEntity.ok(order);
    }
}
```

#### **Important: @Async Proxy Limitations**

```java
@Service
public class ProblematicService {
    
    // ❌ WRONG: Internal call - @Async is IGNORED!
    public void processOrder(Order order) {
        saveOrder(order);
        sendNotification(order);  // This runs SYNCHRONOUSLY!
    }
    
    @Async
    public void sendNotification(Order order) {
        // Won't run async when called from same class
        notificationClient.send(order);
    }
}

// ✅ CORRECT: Separate beans for async methods
@Service
@RequiredArgsConstructor
public class OrderService {
    
    private final AsyncNotificationService asyncNotificationService;
    
    public void processOrder(Order order) {
        saveOrder(order);
        asyncNotificationService.sendNotification(order);  // ✅ Works!
    }
}

@Service
public class AsyncNotificationService {
    
    @Async
    public void sendNotification(Order order) {
        // Runs asynchronously because called through proxy
        notificationClient.send(order);
    }
}
```

#### **Self-Injection Workaround (If You Must)**

```java
@Service
public class SelfInjectingService {
    
    @Autowired
    @Lazy  // Important: Breaks circular dependency
    private SelfInjectingService self;
    
    public void processOrder(Order order) {
        saveOrder(order);
        self.sendNotification(order);  // ✅ Goes through proxy
    }
    
    @Async
    public void sendNotification(Order order) {
        notificationClient.send(order);
    }
}
```

---

### 3.3 Return Types: void vs Future vs CompletableFuture

#### **Comparison Overview**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    @ASYNC RETURN TYPES                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  void                          │ Cannot track completion or get result      │
│  ─────                         │ Fire-and-forget pattern                    │
│  @Async                        │ Exceptions silently handled                │
│  public void process() { }     │                                            │
│                                │                                            │
│  Future<T>                     │ Can check completion, get result           │
│  ─────────                     │ Blocking get() method                      │
│  @Async                        │ Basic exception handling                   │
│  public Future<String> get()   │                                            │
│                                │                                            │
│  CompletableFuture<T>          │ Full async composition support             │
│  ───────────────────           │ Non-blocking operations                    │
│  @Async                        │ Chain, combine, transform results          │
│  public CompletableFuture<T>   │ Exception handling with exceptionally()    │
│                                │                                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Return Type: void**

```java
@Service
@Slf4j
public class VoidAsyncService {
    
    /**
     * Fire-and-Forget Pattern
     * Use for: Notifications, logging, audit, cleanup tasks
     */
    @Async
    public void sendWelcomeEmail(User user) {
        log.info("Starting email send for user: {}", user.getId());
        emailService.send(user.getEmail(), buildWelcomeEmail(user));
        log.info("Email sent for user: {}", user.getId());
        // Caller has no way to know if this succeeded or failed
    }
    
    /**
     * Audit logging - fire and forget
     */
    @Async("auditExecutor")
    public void logUserActivity(String userId, String action, Map<String, Object> details) {
        AuditEntry entry = AuditEntry.builder()
            .userId(userId)
            .action(action)
            .details(details)
            .timestamp(Instant.now())
            .build();
        auditRepository.save(entry);
    }
}
```

#### **Return Type: Future<T>**

```java
@Service
@Slf4j
public class FutureAsyncService {
    
    /**
     * Future allows checking completion and getting result
     * But uses blocking get() - less flexible than CompletableFuture
     */
    @Async
    public Future<UserReport> generateUserReport(Long userId) {
        log.info("Generating report for user: {}", userId);
        
        User user = userRepository.findById(userId).orElseThrow();
        List<Order> orders = orderRepository.findByUserId(userId);
        UserStats stats = analyticsService.calculateStats(userId);
        
        UserReport report = new UserReport(user, orders, stats);
        
        // Wrap result in AsyncResult
        return new AsyncResult<>(report);
    }
    
    /**
     * Future with potential exception
     */
    @Async
    public Future<PaymentResult> processPayment(PaymentRequest request) {
        try {
            PaymentResult result = paymentGateway.process(request);
            return new AsyncResult<>(result);
        } catch (PaymentException e) {
            // Exception will be thrown when caller calls future.get()
            throw new AsyncPaymentException("Payment failed", e);
        }
    }
}

// Calling code
@Service
@RequiredArgsConstructor
public class ReportController {
    
    private final FutureAsyncService asyncService;
    
    @GetMapping("/users/{id}/report")
    public UserReport getUserReport(@PathVariable Long userId) throws Exception {
        Future<UserReport> future = asyncService.generateUserReport(userId);
        
        // Option 1: Blocking wait (defeats purpose somewhat)
        UserReport report = future.get();  // Blocks until complete
        
        // Option 2: Wait with timeout
        UserReport report2 = future.get(30, TimeUnit.SECONDS);
        
        // Option 3: Check if done first
        if (future.isDone()) {
            return future.get();
        }
        
        // Option 4: Cancel if taking too long
        if (!future.isDone()) {
            future.cancel(true);
            throw new TimeoutException("Report generation timed out");
        }
        
        return report;
    }
}
```

#### **Return Type: CompletableFuture<T> (Recommended)**

```java
@Service
@Slf4j
public class CompletableFutureAsyncService {
    
    /**
     * CompletableFuture - Most flexible option
     * Supports chaining, combining, non-blocking operations
     */
    @Async
    public CompletableFuture<UserProfile> fetchUserProfile(Long userId) {
        log.info("Fetching profile on thread: {}", Thread.currentThread().getName());
        
        User user = userRepository.findById(userId).orElseThrow();
        UserProfile profile = profileMapper.toProfile(user);
        
        return CompletableFuture.completedFuture(profile);
    }
    
    @Async
    public CompletableFuture<List<Order>> fetchUserOrders(Long userId) {
        log.info("Fetching orders on thread: {}", Thread.currentThread().getName());
        
        List<Order> orders = orderRepository.findByUserIdOrderByCreatedAtDesc(userId);
        return CompletableFuture.completedFuture(orders);
    }
    
    @Async
    public CompletableFuture<UserStats> fetchUserStats(Long userId) {
        log.info("Fetching stats on thread: {}", Thread.currentThread().getName());
        
        UserStats stats = analyticsService.calculateStats(userId);
        return CompletableFuture.completedFuture(stats);
    }
    
    @Async
    public CompletableFuture<List<Recommendation>> fetchRecommendations(Long userId) {
        log.info("Fetching recommendations on thread: {}", Thread.currentThread().getName());
        
        List<Recommendation> recommendations = recommendationEngine.getFor(userId);
        return CompletableFuture.completedFuture(recommendations);
    }
}

// Powerful composition in calling code
@RestController
@RequiredArgsConstructor
public class DashboardController {
    
    private final CompletableFutureAsyncService asyncService;
    
    /**
     * Parallel execution of multiple async operations
     * All 4 calls execute SIMULTANEOUSLY
     */
    @GetMapping("/users/{id}/dashboard")
    public CompletableFuture<DashboardResponse> getDashboard(@PathVariable Long userId) {
        
        CompletableFuture<UserProfile> profileFuture = asyncService.fetchUserProfile(userId);
        CompletableFuture<List<Order>> ordersFuture = asyncService.fetchUserOrders(userId);
        CompletableFuture<UserStats> statsFuture = asyncService.fetchUserStats(userId);
        CompletableFuture<List<Recommendation>> recosFuture = asyncService.fetchRecommendations(userId);
        
        // Combine all results when ALL complete
        return CompletableFuture.allOf(profileFuture, ordersFuture, statsFuture, recosFuture)
            .thenApply(v -> new DashboardResponse(
                profileFuture.join(),
                ordersFuture.join(),
                statsFuture.join(),
                recosFuture.join()
            ));
    }
    
    /**
     * Chaining async operations
     */
    @GetMapping("/users/{id}/enriched-profile")
    public CompletableFuture<EnrichedProfile> getEnrichedProfile(@PathVariable Long userId) {
        
        return asyncService.fetchUserProfile(userId)
            .thenCompose(profile -> 
                // After profile, fetch orders
                asyncService.fetchUserOrders(userId)
                    .thenApply(orders -> new ProfileWithOrders(profile, orders))
            )
            .thenCompose(profileWithOrders ->
                // After orders, fetch recommendations based on orders
                asyncService.fetchRecommendations(userId)
                    .thenApply(recos -> new EnrichedProfile(
                        profileWithOrders.profile(),
                        profileWithOrders.orders(),
                        recos
                    ))
            );
    }
    
    /**
     * First to complete wins
     */
    @GetMapping("/price/{productId}")
    public CompletableFuture<Price> getFastestPrice(@PathVariable Long productId) {
        
        CompletableFuture<Price> supplier1 = priceService.fetchFromSupplier1(productId);
        CompletableFuture<Price> supplier2 = priceService.fetchFromSupplier2(productId);
        CompletableFuture<Price> supplier3 = priceService.fetchFromSupplier3(productId);
        
        // Return first successful result
        return CompletableFuture.anyOf(supplier1, supplier2, supplier3)
            .thenApply(result -> (Price) result);
    }
}
```

#### **Return Type Comparison Table**

| Feature | void | Future<T> | CompletableFuture<T> |
|---------|------|-----------|---------------------|
| **Get Result** | No | Yes (blocking) | Yes (blocking/non-blocking) |
| **Check Completion** | No | Yes | Yes |
| **Cancel Task** | No | Yes | Yes |
| **Chain Operations** | No | No | Yes (thenApply, thenCompose) |
| **Combine Results** | No | Manual | Yes (allOf, anyOf) |
| **Exception Handling** | AsyncUncaughtExceptionHandler | get() throws | exceptionally(), handle() |
| **Timeout Support** | No | get(timeout) | orTimeout(), completeOnTimeout() |
| **Use Case** | Fire-and-forget | Simple async with result | Complex async workflows |

---

### 3.4 Exception Handling in Async Methods

#### **Exception Handling Overview**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ASYNC EXCEPTION HANDLING                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Return Type          Exception Handling                                    │
│  ───────────          ──────────────────                                    │
│                                                                             │
│  void          ─────► AsyncUncaughtExceptionHandler                         │
│                       (Configure globally)                                  │
│                                                                             │
│  Future<T>     ─────► ExecutionException wrapped                            │
│                       (Thrown on .get() call)                               │
│                                                                             │
│  CompletableFuture<T> ► .exceptionally() or .handle()                       │
│                       (Fluent exception handling)                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Handling Exceptions for void Methods**

```java
/**
 * Custom exception handler for void @Async methods
 */
@Slf4j
public class CustomAsyncExceptionHandler implements AsyncUncaughtExceptionHandler {
    
    @Autowired
    private AlertService alertService;
    
    @Autowired
    private MetricRegistry metricRegistry;
    
    @Override
    public void handleUncaughtException(Throwable ex, Method method, Object... params) {
        log.error("Async exception in method: {} with params: {}",
            method.getName(), 
            Arrays.toString(params),
            ex);
        
        // Record metric
        metricRegistry.counter("async.exceptions", 
            "method", method.getName(),
            "exception", ex.getClass().getSimpleName()
        ).increment();
        
        // Alert on critical exceptions
        if (ex instanceof CriticalBusinessException) {
            alertService.sendAlert(
                "Critical async failure",
                String.format("Method %s failed: %s", method.getName(), ex.getMessage())
            );
        }
        
        // Optionally: Store failed task for retry
        if (method.isAnnotationPresent(Retryable.class)) {
            storeForRetry(method, params, ex);
        }
    }
    
    private void storeForRetry(Method method, Object[] params, Throwable ex) {
        // Store in database or queue for later retry
    }
}

// Register the handler
@Configuration
@EnableAsync
public class AsyncConfig implements AsyncConfigurer {
    
    @Override
    public AsyncUncaughtExceptionHandler getAsyncUncaughtExceptionHandler() {
        return new CustomAsyncExceptionHandler();
    }
    
    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(50);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("async-");
        executor.initialize();
        return executor;
    }
}
```

#### **Handling Exceptions for Future Methods**

```java
@Service
public class FutureExceptionService {
    
    @Async
    public Future<OrderResult> processOrderWithFuture(Order order) {
        try {
            // Business logic that may fail
            validateOrder(order);
            OrderResult result = orderProcessor.process(order);
            return new AsyncResult<>(result);
            
        } catch (ValidationException e) {
            // Wrap and rethrow - caller will get this on .get()
            throw new OrderProcessingException("Validation failed", e);
        } catch (Exception e) {
            throw new OrderProcessingException("Processing failed", e);
        }
    }
}

// Calling code
@Service
public class OrderController {
    
    @Autowired
    private FutureExceptionService service;
    
    public OrderResult processOrder(Order order) {
        Future<OrderResult> future = service.processOrderWithFuture(order);
        
        try {
            return future.get(30, TimeUnit.SECONDS);
            
        } catch (ExecutionException e) {
            // Unwrap the actual exception
            Throwable cause = e.getCause();
            if (cause instanceof OrderProcessingException) {
                throw (OrderProcessingException) cause;
            }
            throw new RuntimeException("Unexpected error", cause);
            
        } catch (TimeoutException e) {
            future.cancel(true);
            throw new OrderTimeoutException("Order processing timed out");
            
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("Processing interrupted", e);
        }
    }
}
```

#### **Handling Exceptions for CompletableFuture (Recommended)**

```java
@Service
@Slf4j
public class CompletableFutureExceptionService {
    
    @Async
    public CompletableFuture<PaymentResult> processPayment(PaymentRequest request) {
        log.info("Processing payment: {}", request.getTransactionId());
        
        // This exception will be captured in the CompletableFuture
        if (request.getAmount().compareTo(BigDecimal.ZERO) <= 0) {
            throw new InvalidPaymentException("Amount must be positive");
        }
        
        PaymentResult result = paymentGateway.process(request);
        return CompletableFuture.completedFuture(result);
    }
}

// Calling code with fluent exception handling
@RestController
@RequiredArgsConstructor
public class PaymentController {
    
    private final CompletableFutureExceptionService paymentService;
    
    @PostMapping("/payments")
    public CompletableFuture<ResponseEntity<PaymentResponse>> processPayment(
            @RequestBody PaymentRequest request) {
        
        return paymentService.processPayment(request)
            // Transform successful result
            .thenApply(result -> ResponseEntity.ok(
                new PaymentResponse(result.getTransactionId(), "SUCCESS")
            ))
            // Handle specific exceptions
            .exceptionally(ex -> {
                Throwable cause = ex.getCause() != null ? ex.getCause() : ex;
                
                if (cause instanceof InvalidPaymentException) {
                    return ResponseEntity.badRequest()
                        .body(new PaymentResponse(null, cause.getMessage()));
                }
                if (cause instanceof PaymentDeclinedException) {
                    return ResponseEntity.status(HttpStatus.PAYMENT_REQUIRED)
                        .body(new PaymentResponse(null, "Payment declined"));
                }
                
                log.error("Unexpected payment error", cause);
                return ResponseEntity.internalServerError()
                    .body(new PaymentResponse(null, "Payment processing failed"));
            });
    }
    
    /**
     * Using handle() for both success and failure
     */
    @PostMapping("/payments/v2")
    public CompletableFuture<ResponseEntity<PaymentResponse>> processPaymentV2(
            @RequestBody PaymentRequest request) {
        
        return paymentService.processPayment(request)
            .handle((result, ex) -> {
                if (ex != null) {
                    // Handle exception
                    log.error("Payment failed", ex);
                    return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                        .body(new PaymentResponse(null, "Failed: " + ex.getMessage()));
                }
                // Handle success
                return ResponseEntity.ok(
                    new PaymentResponse(result.getTransactionId(), "SUCCESS")
                );
            });
    }
    
    /**
     * Recover from failure with default value
     */
    @GetMapping("/prices/{productId}")
    public CompletableFuture<Price> getPrice(@PathVariable Long productId) {
        
        return priceService.fetchLivePrice(productId)
            .exceptionally(ex -> {
                log.warn("Failed to fetch live price, using cached", ex);
                return priceCache.getCachedPrice(productId);  // Fallback
            });
    }
    
    /**
     * Timeout handling (Java 9+)
     */
    @GetMapping("/prices/{productId}/with-timeout")
    public CompletableFuture<Price> getPriceWithTimeout(@PathVariable Long productId) {
        
        return priceService.fetchLivePrice(productId)
            .orTimeout(5, TimeUnit.SECONDS)  // Throws TimeoutException after 5s
            .exceptionally(ex -> {
                if (ex instanceof TimeoutException) {
                    return priceCache.getCachedPrice(productId);
                }
                throw new RuntimeException(ex);
            });
    }
    
    /**
     * Complete with default on timeout (Java 9+)
     */
    @GetMapping("/prices/{productId}/with-default")
    public CompletableFuture<Price> getPriceWithDefault(@PathVariable Long productId) {
        
        return priceService.fetchLivePrice(productId)
            .completeOnTimeout(
                priceCache.getCachedPrice(productId),  // Default value
                5,                                      // Timeout
                TimeUnit.SECONDS                        // Unit
            );
    }
}
```

#### **Exception Handling Patterns Diagram**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    EXCEPTION HANDLING PATTERNS                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Pattern 1: exceptionally() - Handle errors, return recovery value          │
│  ────────────────────────────────────────────────────────────────           │
│                                                                             │
│  asyncMethod()                                                              │
│       │                                                                     │
│       ├── Success ──► thenApply(transform) ──► Result                       │
│       │                                                                     │
│       └── Failure ──► exceptionally(ex -> fallback) ──► Fallback Value      │
│                                                                             │
│  Pattern 2: handle() - Single handler for both success and failure          │
│  ────────────────────────────────────────────────────────────────           │
│                                                                             │
│  asyncMethod()                                                              │
│       │                                                                     │
│       └──► handle((result, ex) -> {                                         │
│                 if (ex != null) return handleError(ex);                     │
│                 return handleSuccess(result);                               │
│            })                                                               │
│                                                                             │
│  Pattern 3: whenComplete() - Side effects, doesn't change result            │
│  ────────────────────────────────────────────────────────────────           │
│                                                                             │
│  asyncMethod()                                                              │
│       │                                                                     │
│       └──► whenComplete((result, ex) -> {                                   │
│                 log.info("Completed with: {} / {}", result, ex);            │
│            })  // Original result or exception propagates                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 3.5 Custom TaskExecutor Configuration

```java
@Configuration
@EnableAsync
@Slf4j
public class AsyncExecutorConfig implements AsyncConfigurer {
    
    /**
     * Default executor for @Async methods without qualifier
     */
    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(50);
        executor.setQueueCapacity(500);
        executor.setThreadNamePrefix("default-async-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(60);
        executor.initialize();
        return executor;
    }
    
    /**
     * Dedicated executor for email operations
     * Lower throughput, less critical
     */
    @Bean("emailExecutor")
    public TaskExecutor emailTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(2);
        executor.setMaxPoolSize(5);
        executor.setQueueCapacity(1000);  // Large queue for email backlog
        executor.setThreadNamePrefix("email-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * Dedicated executor for payment processing
     * High priority, needs fast processing
     */
    @Bean("paymentExecutor")
    public TaskExecutor paymentTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(20);
        executor.setMaxPoolSize(100);
        executor.setQueueCapacity(50);  // Small queue - want immediate processing
        executor.setThreadNamePrefix("payment-");
        // AbortPolicy - don't silently drop payments
        executor.setRejectedExecutionHandler(new AbortPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * Dedicated executor for report generation
     * Long-running, resource-intensive
     */
    @Bean("reportExecutor")
    public TaskExecutor reportTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(2);   // Limited concurrency
        executor.setMaxPoolSize(5);    // Don't overwhelm the system
        executor.setQueueCapacity(20); // Limited queue
        executor.setThreadNamePrefix("report-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * Executor with custom thread factory for priority threads
     */
    @Bean("highPriorityExecutor")
    public TaskExecutor highPriorityTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("high-priority-");
        
        // Custom thread factory with higher priority
        executor.setThreadFactory(r -> {
            Thread t = new Thread(r);
            t.setPriority(Thread.MAX_PRIORITY);
            t.setUncaughtExceptionHandler((thread, ex) -> 
                log.error("Uncaught exception in high priority thread", ex));
            return t;
        });
        
        executor.initialize();
        return executor;
    }
    
    @Override
    public AsyncUncaughtExceptionHandler getAsyncUncaughtExceptionHandler() {
        return new CustomAsyncExceptionHandler();
    }
}

// Using specific executors
@Service
public class MultiExecutorService {
    
    @Async  // Uses default executor
    public void defaultAsyncTask() { }
    
    @Async("emailExecutor")
    public void sendEmail(String to, String content) { }
    
    @Async("paymentExecutor")
    public CompletableFuture<PaymentResult> processPayment(Payment payment) {
        return CompletableFuture.completedFuture(paymentGateway.process(payment));
    }
    
    @Async("reportExecutor")
    public CompletableFuture<Report> generateReport(ReportRequest request) {
        return CompletableFuture.completedFuture(reportGenerator.generate(request));
    }
}
```

---

### 3.6 Complete Working Example

#### **Project Structure**

```
src/main/java/com/example/asyncdemo/
├── AsyncDemoApplication.java
├── config/
│   └── AsyncConfig.java
├── controller/
│   └── OrderController.java
├── service/
│   ├── OrderService.java
│   ├── AsyncNotificationService.java
│   ├── AsyncInventoryService.java
│   └── AsyncPaymentService.java
├── model/
│   ├── Order.java
│   ├── OrderRequest.java
│   └── OrderResponse.java
└── exception/
    └── CustomAsyncExceptionHandler.java
```

#### **Application Configuration**

```java
// AsyncDemoApplication.java
@SpringBootApplication
public class AsyncDemoApplication {
    public static void main(String[] args) {
        SpringApplication.run(AsyncDemoApplication.class, args);
    }
}
```

```java
// config/AsyncConfig.java
@Configuration
@EnableAsync
@Slf4j
public class AsyncConfig implements AsyncConfigurer {
    
    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("main-async-");
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(30);
        executor.initialize();
        log.info("Main async executor initialized");
        return executor;
    }
    
    @Bean("notificationExecutor")
    public TaskExecutor notificationExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(2);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(500);
        executor.setThreadNamePrefix("notification-");
        executor.initialize();
        return executor;
    }
    
    @Bean("inventoryExecutor")
    public TaskExecutor inventoryExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(3);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(50);
        executor.setThreadNamePrefix("inventory-");
        executor.initialize();
        return executor;
    }
    
    @Bean("paymentExecutor")
    public TaskExecutor paymentExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(25);
        executor.setThreadNamePrefix("payment-");
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.AbortPolicy());
        executor.initialize();
        return executor;
    }
    
    @Override
    public AsyncUncaughtExceptionHandler getAsyncUncaughtExceptionHandler() {
        return new CustomAsyncExceptionHandler();
    }
}
```

```java
// exception/CustomAsyncExceptionHandler.java
@Slf4j
@Component
public class CustomAsyncExceptionHandler implements AsyncUncaughtExceptionHandler {
    
    @Override
    public void handleUncaughtException(Throwable ex, Method method, Object... params) {
        log.error("Async exception in method [{}] with params [{}]: {}",
            method.getName(),
            Arrays.toString(params),
            ex.getMessage(),
            ex);
        
        // Add monitoring/alerting here
    }
}
```

#### **Model Classes**

```java
// model/OrderRequest.java
@Data
@Builder
public class OrderRequest {
    private String customerId;
    private String customerEmail;
    private List<OrderItem> items;
    private BigDecimal totalAmount;
    private String paymentMethod;
}

// model/Order.java
@Data
@Builder
public class Order {
    private String id;
    private String customerId;
    private OrderStatus status;
    private BigDecimal totalAmount;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}

// model/OrderResponse.java
@Data
@Builder
public class OrderResponse {
    private String orderId;
    private String status;
    private String message;
    private PaymentResult paymentResult;
    private InventoryResult inventoryResult;
}
```

#### **Async Services**

```java
// service/AsyncNotificationService.java
@Service
@Slf4j
public class AsyncNotificationService {
    
    /**
     * Fire-and-forget notification
     * Uses dedicated notification executor
     */
    @Async("notificationExecutor")
    public void sendOrderConfirmationEmail(String email, Order order) {
        log.info("[{}] Sending order confirmation email to: {}",
            Thread.currentThread().getName(), email);
        
        try {
            // Simulate email sending
            Thread.sleep(2000);
            log.info("[{}] Order confirmation email sent successfully to: {}",
                Thread.currentThread().getName(), email);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            log.error("Email sending interrupted", e);
        }
    }
    
    @Async("notificationExecutor")
    public void sendSmsNotification(String phone, String message) {
        log.info("[{}] Sending SMS to: {}", Thread.currentThread().getName(), phone);
        // SMS sending logic
    }
    
    @Async("notificationExecutor")
    public void sendPushNotification(String userId, String title, String body) {
        log.info("[{}] Sending push notification to user: {}",
            Thread.currentThread().getName(), userId);
        // Push notification logic
    }
}
```

```java
// service/AsyncInventoryService.java
@Service
@Slf4j
public class AsyncInventoryService {
    
    /**
     * Async inventory reservation with CompletableFuture
     * Caller can wait for result or chain operations
     */
    @Async("inventoryExecutor")
    public CompletableFuture<InventoryResult> reserveInventory(Order order) {
        log.info("[{}] Reserving inventory for order: {}",
            Thread.currentThread().getName(), order.getId());
        
        try {
            // Simulate inventory check and reservation
            Thread.sleep(1500);
            
            // Simulate business logic
            boolean success = Math.random() > 0.1;  // 90% success rate
            
            if (success) {
                log.info("[{}] Inventory reserved successfully for order: {}",
                    Thread.currentThread().getName(), order.getId());
                return CompletableFuture.completedFuture(
                    InventoryResult.builder()
                        .orderId(order.getId())
                        .status("RESERVED")
                        .reservationId(UUID.randomUUID().toString())
                        .build()
                );
            } else {
                throw new InventoryException("Insufficient inventory");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("Inventory reservation interrupted", e);
        }
    }
    
    @Async("inventoryExecutor")
    public CompletableFuture<Void> releaseInventory(String reservationId) {
        log.info("[{}] Releasing inventory reservation: {}",
            Thread.currentThread().getName(), reservationId);
        // Release logic
        return CompletableFuture.completedFuture(null);
    }
}
```

```java
// service/AsyncPaymentService.java
@Service
@Slf4j
public class AsyncPaymentService {
    
    /**
     * Async payment processing
     * Returns CompletableFuture for result handling
     */
    @Async("paymentExecutor")
    public CompletableFuture<PaymentResult> processPayment(Order order, String paymentMethod) {
        log.info("[{}] Processing payment for order: {}, method: {}",
            Thread.currentThread().getName(), order.getId(), paymentMethod);
        
        try {
            // Simulate payment gateway call
            Thread.sleep(2000);
            
            // Simulate payment result
            boolean success = Math.random() > 0.05;  // 95% success rate
            
            if (success) {
                log.info("[{}] Payment successful for order: {}",
                    Thread.currentThread().getName(), order.getId());
                return CompletableFuture.completedFuture(
                    PaymentResult.builder()
                        .orderId(order.getId())
                        .transactionId(UUID.randomUUID().toString())
                        .status("SUCCESS")
                        .processedAt(LocalDateTime.now())
                        .build()
                );
            } else {
                throw new PaymentException("Payment declined by gateway");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("Payment processing interrupted", e);
        }
    }
    
    @Async("paymentExecutor")
    public CompletableFuture<RefundResult> processRefund(String transactionId, BigDecimal amount) {
        log.info("[{}] Processing refund for transaction: {}",
            Thread.currentThread().getName(), transactionId);
        // Refund logic
        return CompletableFuture.completedFuture(
            RefundResult.builder().status("REFUNDED").build()
        );
    }
}
```

#### **Order Service (Orchestrator)**

```java
// service/OrderService.java
@Service
@RequiredArgsConstructor
@Slf4j
public class OrderService {
    
    private final AsyncNotificationService notificationService;
    private final AsyncInventoryService inventoryService;
    private final AsyncPaymentService paymentService;
    private final OrderRepository orderRepository;
    
    /**
     * Create order with parallel async operations
     */
    public OrderResponse createOrder(OrderRequest request) {
        log.info("[{}] Creating order for customer: {}",
            Thread.currentThread().getName(), request.getCustomerId());
        
        // 1. Create order synchronously (main business logic)
        Order order = Order.builder()
            .id(UUID.randomUUID().toString())
            .customerId(request.getCustomerId())
            .totalAmount(request.getTotalAmount())
            .status(OrderStatus.PENDING)
            .createdAt(LocalDateTime.now())
            .build();
        
        orderRepository.save(order);
        
        // 2. Start parallel async operations
        CompletableFuture<InventoryResult> inventoryFuture = 
            inventoryService.reserveInventory(order);
        CompletableFuture<PaymentResult> paymentFuture = 
            paymentService.processPayment(order, request.getPaymentMethod());
        
        // 3. Fire-and-forget: Send notification (don't wait)
        notificationService.sendOrderConfirmationEmail(
            request.getCustomerEmail(), order);
        
        // 4. Wait for critical operations to complete
        try {
            CompletableFuture.allOf(inventoryFuture, paymentFuture)
                .get(30, TimeUnit.SECONDS);
            
            InventoryResult inventoryResult = inventoryFuture.join();
            PaymentResult paymentResult = paymentFuture.join();
            
            // 5. Update order status
            order.setStatus(OrderStatus.CONFIRMED);
            order.setUpdatedAt(LocalDateTime.now());
            orderRepository.save(order);
            
            log.info("[{}] Order created successfully: {}",
                Thread.currentThread().getName(), order.getId());
            
            return OrderResponse.builder()
                .orderId(order.getId())
                .status("SUCCESS")
                .message("Order created successfully")
                .paymentResult(paymentResult)
                .inventoryResult(inventoryResult)
                .build();
                
        } catch (TimeoutException e) {
            log.error("Order processing timed out", e);
            order.setStatus(OrderStatus.FAILED);
            orderRepository.save(order);
            throw new OrderProcessingException("Order processing timed out");
            
        } catch (ExecutionException e) {
            log.error("Order processing failed", e.getCause());
            order.setStatus(OrderStatus.FAILED);
            orderRepository.save(order);
            
            // Compensating action: release inventory if payment failed
            if (inventoryFuture.isDone() && !inventoryFuture.isCompletedExceptionally()) {
                inventoryService.releaseInventory(
                    inventoryFuture.join().getReservationId());
            }
            
            throw new OrderProcessingException("Order processing failed", e.getCause());
            
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("Order processing interrupted", e);
        }
    }
    
    /**
     * Alternative: Fully async order creation returning CompletableFuture
     */
    public CompletableFuture<OrderResponse> createOrderAsync(OrderRequest request) {
        log.info("[{}] Creating order async for customer: {}",
            Thread.currentThread().getName(), request.getCustomerId());
        
        Order order = Order.builder()
            .id(UUID.randomUUID().toString())
            .customerId(request.getCustomerId())
            .totalAmount(request.getTotalAmount())
            .status(OrderStatus.PENDING)
            .createdAt(LocalDateTime.now())
            .build();
        
        orderRepository.save(order);
        
        // Fire-and-forget notification
        notificationService.sendOrderConfirmationEmail(
            request.getCustomerEmail(), order);
        
        // Parallel execution with combined result
        CompletableFuture<InventoryResult> inventoryFuture = 
            inventoryService.reserveInventory(order);
        CompletableFuture<PaymentResult> paymentFuture = 
            paymentService.processPayment(order, request.getPaymentMethod());
        
        return inventoryFuture
            .thenCombine(paymentFuture, (inventory, payment) -> {
                order.setStatus(OrderStatus.CONFIRMED);
                orderRepository.save(order);
                
                return OrderResponse.builder()
                    .orderId(order.getId())
                    .status("SUCCESS")
                    .paymentResult(payment)
                    .inventoryResult(inventory)
                    .build();
            })
            .exceptionally(ex -> {
                log.error("Order async processing failed", ex);
                order.setStatus(OrderStatus.FAILED);
                orderRepository.save(order);
                
                return OrderResponse.builder()
                    .orderId(order.getId())
                    .status("FAILED")
                    .message(ex.getMessage())
                    .build();
            });
    }
}
```

#### **Controller**

```java
// controller/OrderController.java
@RestController
@RequestMapping("/api/orders")
@RequiredArgsConstructor
@Slf4j
public class OrderController {
    
    private final OrderService orderService;
    
    /**
     * Synchronous endpoint with internal async operations
     */
    @PostMapping
    public ResponseEntity<OrderResponse> createOrder(@RequestBody OrderRequest request) {
        log.info("[{}] Received order request", Thread.currentThread().getName());
        long start = System.currentTimeMillis();
        
        OrderResponse response = orderService.createOrder(request);
        
        log.info("[{}] Order processed in {}ms",
            Thread.currentThread().getName(),
            System.currentTimeMillis() - start);
        
        return ResponseEntity.ok(response);
    }
    
    /**
     * Fully async endpoint returning CompletableFuture
     * Spring handles the async response automatically
     */
    @PostMapping("/async")
    public CompletableFuture<ResponseEntity<OrderResponse>> createOrderAsync(
            @RequestBody OrderRequest request) {
        log.info("[{}] Received async order request", Thread.currentThread().getName());
        
        return orderService.createOrderAsync(request)
            .thenApply(ResponseEntity::ok)
            .exceptionally(ex -> ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(OrderResponse.builder()
                    .status("FAILED")
                    .message(ex.getMessage())
                    .build()));
    }
    
    /**
     * Endpoint to demonstrate thread info
     */
    @GetMapping("/thread-info")
    public Map<String, String> getThreadInfo() {
        return Map.of(
            "currentThread", Thread.currentThread().getName(),
            "threadGroup", Thread.currentThread().getThreadGroup().getName()
        );
    }
}
```

#### **Application Properties**

```yaml
# application.yml
server:
  port: 8080
  tomcat:
    threads:
      max: 200
      min-spare: 10

spring:
  application:
    name: async-demo

logging:
  level:
    com.example.asyncdemo: DEBUG
  pattern:
    console: "%d{HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n"
```

#### **Sample Test Output**

```
14:23:45.123 [http-nio-8080-exec-1] INFO  OrderController - Received order request
14:23:45.125 [http-nio-8080-exec-1] INFO  OrderService - Creating order for customer: C001
14:23:45.130 [notification-1] INFO  AsyncNotificationService - Sending order confirmation email to: user@example.com
14:23:45.132 [inventory-1] INFO  AsyncInventoryService - Reserving inventory for order: abc-123
14:23:45.134 [payment-1] INFO  AsyncPaymentService - Processing payment for order: abc-123, method: CREDIT_CARD
14:23:46.640 [inventory-1] INFO  AsyncInventoryService - Inventory reserved successfully for order: abc-123
14:23:47.135 [notification-1] INFO  AsyncNotificationService - Order confirmation email sent successfully
14:23:47.140 [payment-1] INFO  AsyncPaymentService - Payment successful for order: abc-123
14:23:47.145 [http-nio-8080-exec-1] INFO  OrderService - Order created successfully: abc-123
14:23:47.150 [http-nio-8080-exec-1] INFO  OrderController - Order processed in 2027ms
```

---

## 4. ThreadPoolTaskExecutor Configuration

### 4.1 Understanding Thread Pool Parameters

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THREADPOOLTASKEXECUTOR PARAMETERS                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    ThreadPoolTaskExecutor                            │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  corePoolSize: 10          │ Always-alive threads                   │   │
│  │  ─────────────────         │ Created on demand, never destroyed     │   │
│  │                            │ (unless allowCoreThreadTimeOut=true)   │   │
│  │                                                                      │   │
│  │  maxPoolSize: 50           │ Maximum threads under peak load        │   │
│  │  ─────────────────         │ Additional threads beyond core         │   │
│  │                            │ Destroyed when idle                    │   │
│  │                                                                      │   │
│  │  queueCapacity: 100        │ Tasks waiting when all threads busy   │   │
│  │  ─────────────────         │ New threads created only when full    │   │
│  │                            │ (and below maxPoolSize)                │   │
│  │                                                                      │   │
│  │  keepAliveSeconds: 60      │ Idle time before killing non-core     │   │
│  │  ─────────────────         │ threads                                │   │
│  │                                                                      │   │
│  │  threadNamePrefix: "async-"│ Thread naming for debugging           │   │
│  │  ─────────────────         │                                        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **How Tasks Are Processed**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TASK PROCESSING FLOW                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                              New Task Arrives                               │
│                                     │                                       │
│                                     ▼                                       │
│                    ┌────────────────────────────────┐                       │
│                    │ Current Threads < corePoolSize │                       │
│                    └────────────────┬───────────────┘                       │
│                            Yes      │      No                               │
│                             ▼       │       ▼                               │
│                    ┌──────────┐    │     ┌─────────────────────┐           │
│                    │ Create   │    │     │ Queue has capacity? │           │
│                    │ new core │    │     └──────────┬──────────┘           │
│                    │ thread   │    │        Yes     │     No               │
│                    └──────────┘    │         ▼      │      ▼               │
│                                    │  ┌──────────┐ │  ┌───────────────────┐│
│                                    │  │Add task  │ │  │Threads<maxPoolSize││
│                                    │  │to queue  │ │  └─────────┬─────────┘│
│                                    │  └──────────┘ │    Yes     │    No    │
│                                    │               │     ▼      │     ▼    │
│                                    │               │ ┌────────┐│┌────────┐ │
│                                    │               │ │Create  │││Reject  │ │
│                                    │               │ │non-core│││Policy  │ │
│                                    │               │ │thread  │││Triggered│ │
│                                    │               │ └────────┘│└────────┘ │
│                                    │               │           │           │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 4.2 Core Pool Size vs Max Pool Size

#### **Understanding the Difference**

```java
@Configuration
public class ThreadPoolExplanation {
    
    /**
     * SCENARIO: corePoolSize=5, maxPoolSize=20, queueCapacity=100
     * 
     * State 1: Initial (idle system)
     * - Active threads: 0
     * - Task arrives: Creates thread #1 (core thread)
     * 
     * State 2: Light load (5 concurrent tasks)
     * - Core threads: 5 (all busy)
     * - Queue: 0
     * - New task: Goes to queue
     * 
     * State 3: Moderate load (105 concurrent tasks submitted)
     * - Core threads: 5 (all busy)
     * - Queue: 100 (FULL!)
     * - New task: Creates thread #6 (non-core)
     * 
     * State 4: Heavy load (120 concurrent tasks)
     * - Threads: 20 (5 core + 15 non-core)
     * - Queue: 100 (full)
     * - All threads busy
     * 
     * State 5: Overload (121+ concurrent tasks)
     * - Threads: 20 (at max)
     * - Queue: 100 (full)
     * - New task: REJECTED (policy applied)
     * 
     * State 6: Load decreases, threads idle
     * - After keepAliveSeconds: Non-core threads destroyed
     * - Eventually: Only 5 core threads remain
     */
    
    @Bean
    public ThreadPoolTaskExecutor explainedExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);      // Always keep 5 threads alive
        executor.setMaxPoolSize(20);      // Can grow to 20 under pressure
        executor.setQueueCapacity(100);   // 100 tasks can wait
        executor.setKeepAliveSeconds(60); // Kill idle non-core after 60s
        executor.initialize();
        return executor;
    }
}
```

#### **Visual Timeline**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THREAD POOL LIFECYCLE                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Time ──────────────────────────────────────────────────────────────────►   │
│                                                                             │
│  Load:     Low        Rising        Peak         Declining         Idle    │
│           ─────       ──────       ─────         ─────────        ─────    │
│                                                                             │
│  Threads:   5    →      5     →     20       →      20    →   5             │
│            (core)     (core)    (5+15 non-core) (shrinking)    (core only) │
│                                                                             │
│  Queue:     0    →     50     →    100      →      50     →    0           │
│                      (filling)   (full!)       (draining)     (empty)      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ Threads: ████████████████████████████████████████░░░░░░░░░░░░░░░░░░│   │
│  │ Queue:   ░░░░░░░████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░│   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│             ↑            ↑            ↑             ↑            ↑         │
│          Startup    Core full    Max reached   Threads die    Stable       │
│                     Queue fills  Queue full    after timeout  at core      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 4.3 Queue Capacity and Behavior

#### **Queue Types and Behavior**

```java
@Configuration
public class QueueBehaviorConfig {
    
    /**
     * DEFAULT: LinkedBlockingQueue with capacity
     * Tasks queued until capacity, then new threads up to maxPoolSize
     */
    @Bean("boundedQueueExecutor")
    public ThreadPoolTaskExecutor boundedQueueExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);  // Bounded queue
        executor.initialize();
        return executor;
    }
    
    /**
     * UNBOUNDED QUEUE: Set queueCapacity to Integer.MAX_VALUE
     * WARNING: maxPoolSize is effectively ignored!
     * New threads only created if queue is full (never happens)
     */
    @Bean("unboundedQueueExecutor")
    public ThreadPoolTaskExecutor unboundedQueueExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);  // ⚠️ IGNORED - queue never fills
        executor.setQueueCapacity(Integer.MAX_VALUE);  // Unbounded
        executor.initialize();
        return executor;
    }
    
    /**
     * DIRECT HANDOFF: Set queueCapacity to 0
     * Tasks go directly to threads, no queuing
     * Immediately creates threads up to maxPoolSize
     */
    @Bean("directHandoffExecutor")
    public ThreadPoolTaskExecutor directHandoffExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(0);   // No core threads
        executor.setMaxPoolSize(100);  // Scales up to 100
        executor.setQueueCapacity(0);  // SynchronousQueue - direct handoff
        executor.setKeepAliveSeconds(60);
        executor.initialize();
        return executor;
    }
}
```

#### **Queue Strategy Comparison**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    QUEUE STRATEGY COMPARISON                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Strategy          │ queueCapacity │ Behavior                              │
│  ────────────────────────────────────────────────────────────────────────── │
│                                                                             │
│  BOUNDED QUEUE     │ 100           │ Tasks queue, then threads scale       │
│  (Recommended)     │               │ Predictable resource usage            │
│                    │               │ Good for most applications            │
│                                                                             │
│  UNBOUNDED QUEUE   │ MAX_VALUE     │ Tasks always queue                    │
│  (Careful!)        │               │ maxPoolSize has no effect             │
│                    │               │ Can cause OOM under heavy load        │
│                                                                             │
│  DIRECT HANDOFF    │ 0             │ No queuing, immediate thread creation │
│  (High throughput) │               │ Best for short-lived tasks            │
│                    │               │ Higher rejection rate under load      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 4.4 Rejection Policies

When the thread pool is at max capacity AND the queue is full, a rejection policy determines what happens to new tasks.

#### **Available Policies**

```java
@Configuration
public class RejectionPolicyConfig {
    
    /**
     * ABORT POLICY (Default)
     * Throws RejectedExecutionException
     * Use when: You want to know when system is overloaded
     */
    @Bean("abortExecutor")
    public ThreadPoolTaskExecutor abortPolicyExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(50);
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.AbortPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * CALLER RUNS POLICY (Recommended for most cases)
     * Runs task in the calling thread
     * Provides natural backpressure - slows down task submission
     * Use when: You don't want to lose tasks
     */
    @Bean("callerRunsExecutor")
    public ThreadPoolTaskExecutor callerRunsPolicyExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(50);
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * DISCARD POLICY
     * Silently drops the rejected task
     * Use when: Task loss is acceptable (e.g., metrics, non-critical logging)
     */
    @Bean("discardExecutor")
    public ThreadPoolTaskExecutor discardPolicyExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(50);
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.DiscardPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * DISCARD OLDEST POLICY
     * Drops the oldest unprocessed task, then retries submission
     * Use when: Newer tasks are more important than older ones
     */
    @Bean("discardOldestExecutor")
    public ThreadPoolTaskExecutor discardOldestPolicyExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(50);
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.DiscardOldestPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * CUSTOM POLICY
     * Implement your own handling logic
     */
    @Bean("customPolicyExecutor")
    public ThreadPoolTaskExecutor customPolicyExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(50);
        executor.setRejectedExecutionHandler(new CustomRejectionHandler());
        executor.initialize();
        return executor;
    }
}

/**
 * Custom rejection handler with logging, metrics, and fallback
 */
@Slf4j
public class CustomRejectionHandler implements RejectedExecutionHandler {
    
    private final AtomicLong rejectedCount = new AtomicLong(0);
    private final MeterRegistry meterRegistry;
    
    public CustomRejectionHandler(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
    }
    
    // Simple constructor for when metrics aren't available
    public CustomRejectionHandler() {
        this.meterRegistry = null;
    }
    
    @Override
    public void rejectedExecution(Runnable r, ThreadPoolExecutor executor) {
        long count = rejectedCount.incrementAndGet();
        
        // Log the rejection
        log.warn("Task rejected! Total rejections: {}. Pool status - " +
            "Active: {}, Pool Size: {}, Queue Size: {}",
            count,
            executor.getActiveCount(),
            executor.getPoolSize(),
            executor.getQueue().size());
        
        // Record metric if available
        if (meterRegistry != null) {
            meterRegistry.counter("executor.rejected.tasks").increment();
        }
        
        // Option 1: Run in caller's thread (like CallerRunsPolicy)
        if (!executor.isShutdown()) {
            log.info("Running rejected task in caller thread");
            r.run();
        }
        
        // Option 2: Store for later retry
        // retryQueue.offer(r);
        
        // Option 3: Send alert if too many rejections
        if (count % 100 == 0) {
            sendAlert("High task rejection rate: " + count);
        }
    }
    
    private void sendAlert(String message) {
        // Send to monitoring system
    }
    
    public long getRejectedCount() {
        return rejectedCount.get();
    }
}
```

#### **Rejection Policy Comparison**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    REJECTION POLICY COMPARISON                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Policy          │ Behavior                  │ Best For                    │
│  ────────────────────────────────────────────────────────────────────────── │
│                                                                             │
│  AbortPolicy     │ Throws exception          │ Critical tasks where        │
│  (Default)       │                           │ failure must be known       │
│                                                                             │
│  CallerRunsPolicy│ Caller thread executes    │ Most applications           │
│  (Recommended)   │ Natural backpressure      │ No task loss acceptable     │
│                                                                             │
│  DiscardPolicy   │ Silently drops task       │ Non-critical work           │
│                  │                           │ (metrics, optional logs)    │
│                                                                             │
│  DiscardOldest   │ Drops oldest queued task  │ Time-sensitive data         │
│  Policy          │ Adds new task             │ (latest value matters)      │
│                                                                             │
│  Custom          │ Your logic                │ Complex requirements        │
│                  │ (log, retry, alert)       │ (retry, dead letter queue)  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 4.5 Complete Configuration Examples

#### **Configuration via Java**

```java
@Configuration
@EnableAsync
@Slf4j
public class ProductionAsyncConfig implements AsyncConfigurer {
    
    @Value("${async.core-pool-size:10}")
    private int corePoolSize;
    
    @Value("${async.max-pool-size:50}")
    private int maxPoolSize;
    
    @Value("${async.queue-capacity:500}")
    private int queueCapacity;
    
    /**
     * Production-ready default async executor
     */
    @Override
    public Executor getAsyncExecutor() {
        return createExecutor("default-async-", corePoolSize, maxPoolSize, queueCapacity);
    }
    
    /**
     * High-priority executor for critical operations
     */
    @Bean("criticalExecutor")
    public TaskExecutor criticalTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(20);
        executor.setMaxPoolSize(100);
        executor.setQueueCapacity(100);  // Small queue - want fast processing
        executor.setThreadNamePrefix("critical-");
        executor.setRejectedExecutionHandler(new AbortPolicy());  // Fail fast
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(120);  // Wait longer for critical tasks
        executor.initialize();
        return executor;
    }
    
    /**
     * Bulk operations executor
     */
    @Bean("bulkExecutor")
    public TaskExecutor bulkTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(10000);  // Large queue for bulk ops
        executor.setThreadNamePrefix("bulk-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(300);  // 5 min for bulk completion
        executor.initialize();
        return executor;
    }
    
    /**
     * Schedule/Cron executor
     */
    @Bean("schedulerExecutor")
    public TaskExecutor schedulerTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(2);
        executor.setMaxPoolSize(5);
        executor.setQueueCapacity(50);
        executor.setThreadNamePrefix("scheduler-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * Event processing executor
     */
    @Bean("eventExecutor")
    @ConditionalOnProperty(name = "events.async.enabled", havingValue = "true")
    public TaskExecutor eventTaskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(50);
        executor.setQueueCapacity(1000);
        executor.setThreadNamePrefix("event-");
        executor.setRejectedExecutionHandler(new DiscardOldestPolicy());
        executor.initialize();
        return executor;
    }
    
    private ThreadPoolTaskExecutor createExecutor(
            String prefix, int core, int max, int queue) {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(core);
        executor.setMaxPoolSize(max);
        executor.setQueueCapacity(queue);
        executor.setThreadNamePrefix(prefix);
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(60);
        
        // Task decorator for MDC propagation
        executor.setTaskDecorator(new MdcTaskDecorator());
        
        executor.initialize();
        log.info("Created executor [{}] - core: {}, max: {}, queue: {}",
            prefix, core, max, queue);
        return executor;
    }
    
    @Override
    public AsyncUncaughtExceptionHandler getAsyncUncaughtExceptionHandler() {
        return new ProductionAsyncExceptionHandler();
    }
}

/**
 * MDC Task Decorator - Propagates MDC context to async threads
 */
public class MdcTaskDecorator implements TaskDecorator {
    
    @Override
    public Runnable decorate(Runnable runnable) {
        // Capture MDC from calling thread
        Map<String, String> contextMap = MDC.getCopyOfContextMap();
        
        return () -> {
            try {
                // Set MDC in async thread
                if (contextMap != null) {
                    MDC.setContextMap(contextMap);
                }
                runnable.run();
            } finally {
                MDC.clear();
            }
        };
    }
}
```

#### **Configuration via Properties/YAML**

```yaml
# application.yml - Environment-agnostic base configuration
spring:
  task:
    execution:
      pool:
        core-size: 10
        max-size: 50
        queue-capacity: 500
        keep-alive: 60s
        allow-core-thread-timeout: false
      thread-name-prefix: "spring-async-"
      shutdown:
        await-termination: true
        await-termination-period: 60s

# Custom executor properties
async:
  core-pool-size: ${ASYNC_CORE_POOL_SIZE:10}
  max-pool-size: ${ASYNC_MAX_POOL_SIZE:50}
  queue-capacity: ${ASYNC_QUEUE_CAPACITY:500}

---
# application-dev.yml
spring:
  config:
    activate:
      on-profile: dev

async:
  core-pool-size: 5
  max-pool-size: 10
  queue-capacity: 100

---
# application-prod.yml
spring:
  config:
    activate:
      on-profile: prod

async:
  core-pool-size: 20
  max-pool-size: 100
  queue-capacity: 1000
```

#### **Monitoring with Actuator**

```java
@Component
@RequiredArgsConstructor
public class ExecutorMetricsConfiguration {
    
    private final MeterRegistry meterRegistry;
    
    @EventListener(ApplicationReadyEvent.class)
    public void registerExecutorMetrics(
            @Autowired @Qualifier("taskExecutor") ThreadPoolTaskExecutor executor) {
        
        // Register metrics for the executor
        Gauge.builder("executor.pool.size", executor, ThreadPoolTaskExecutor::getPoolSize)
            .tag("name", "taskExecutor")
            .description("Current pool size")
            .register(meterRegistry);
        
        Gauge.builder("executor.pool.core", executor, ThreadPoolTaskExecutor::getCorePoolSize)
            .tag("name", "taskExecutor")
            .description("Core pool size")
            .register(meterRegistry);
        
        Gauge.builder("executor.pool.max", executor, ThreadPoolTaskExecutor::getMaxPoolSize)
            .tag("name", "taskExecutor")
            .description("Max pool size")
            .register(meterRegistry);
        
        Gauge.builder("executor.active", executor, ThreadPoolTaskExecutor::getActiveCount)
            .tag("name", "taskExecutor")
            .description("Active threads")
            .register(meterRegistry);
        
        Gauge.builder("executor.queue.size", 
                executor, e -> e.getThreadPoolExecutor().getQueue().size())
            .tag("name", "taskExecutor")
            .description("Queue size")
            .register(meterRegistry);
        
        Gauge.builder("executor.queue.remaining",
                executor, e -> e.getThreadPoolExecutor().getQueue().remainingCapacity())
            .tag("name", "taskExecutor")
            .description("Queue remaining capacity")
            .register(meterRegistry);
    }
}

// Actuator endpoint for thread pool status
@RestController
@RequestMapping("/actuator")
@RequiredArgsConstructor
public class ExecutorInfoEndpoint {
    
    private final Map<String, ThreadPoolTaskExecutor> executors;
    
    @Autowired
    public ExecutorInfoEndpoint(ApplicationContext context) {
        this.executors = context.getBeansOfType(ThreadPoolTaskExecutor.class);
    }
    
    @GetMapping("/executors")
    public Map<String, Object> getExecutorInfo() {
        Map<String, Object> info = new HashMap<>();
        
        executors.forEach((name, executor) -> {
            ThreadPoolExecutor pool = executor.getThreadPoolExecutor();
            info.put(name, Map.of(
                "corePoolSize", pool.getCorePoolSize(),
                "maxPoolSize", pool.getMaximumPoolSize(),
                "poolSize", pool.getPoolSize(),
                "activeCount", pool.getActiveCount(),
                "queueSize", pool.getQueue().size(),
                "queueCapacity", pool.getQueue().remainingCapacity() + pool.getQueue().size(),
                "completedTasks", pool.getCompletedTaskCount(),
                "taskCount", pool.getTaskCount()
            ));
        });
        
        return info;
    }
}
```

---

### 4.6 Production Tuning Guidelines

#### **Thread Pool Sizing Formulas**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THREAD POOL SIZING FORMULAS                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  FORMULA 1: General Purpose                                                 │
│  ──────────────────────────────                                             │
│                                                                             │
│  Optimal Threads = CPU Cores × (1 + Wait Time / Service Time)               │
│                                                                             │
│  Example: 8 cores, 100ms wait, 10ms compute                                 │
│  Threads = 8 × (1 + 100/10) = 8 × 11 = 88 threads                          │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  FORMULA 2: Little's Law                                                    │
│  ──────────────────────────────                                             │
│                                                                             │
│  Threads = Target Throughput × Average Latency                              │
│                                                                             │
│  Example: 1000 req/sec target, 200ms average latency                        │
│  Threads = 1000 × 0.2 = 200 threads                                        │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  FORMULA 3: CPU-Bound Tasks                                                 │
│  ──────────────────────────────                                             │
│                                                                             │
│  Threads = CPU Cores + 1                                                    │
│  (Extra thread to utilize CPU during brief I/O)                             │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  FORMULA 4: I/O-Bound Tasks                                                 │
│  ──────────────────────────────                                             │
│                                                                             │
│  Threads = CPU Cores × 2 (minimum)                                          │
│  Up to CPU Cores × 10 for heavy I/O                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Production Configuration by Use Case**

```java
@Configuration
@EnableAsync
public class ProductionTunedConfig {
    
    /**
     * API GATEWAY / HIGH THROUGHPUT SERVICE
     * Characteristics: Many short requests, mostly I/O
     */
    @Bean("apiGatewayExecutor")
    @Profile("gateway")
    public ThreadPoolTaskExecutor apiGatewayExecutor() {
        int cores = Runtime.getRuntime().availableProcessors();
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(cores * 4);        // High core for I/O
        executor.setMaxPoolSize(cores * 8);         // Allow scaling
        executor.setQueueCapacity(200);             // Moderate queue
        executor.setKeepAliveSeconds(30);           // Quick scale down
        executor.setThreadNamePrefix("gateway-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * DATA PROCESSING SERVICE
     * Characteristics: Long-running, CPU-intensive
     */
    @Bean("dataProcessingExecutor")
    @Profile("processing")
    public ThreadPoolTaskExecutor dataProcessingExecutor() {
        int cores = Runtime.getRuntime().availableProcessors();
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(cores);            // Match CPU cores
        executor.setMaxPoolSize(cores + 2);         // Slight headroom
        executor.setQueueCapacity(1000);            // Large queue
        executor.setKeepAliveSeconds(120);          // Longer keepalive
        executor.setThreadNamePrefix("processing-");
        executor.setRejectedExecutionHandler(new CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
    
    /**
     * NOTIFICATION SERVICE
     * Characteristics: I/O heavy, can tolerate delays
     */
    @Bean("notificationExecutor")
    @Profile("notification")
    public ThreadPoolTaskExecutor notificationExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(10000);  // Large queue - notifications can wait
        executor.setKeepAliveSeconds(60);
        executor.setThreadNamePrefix("notify-");
        executor.setRejectedExecutionHandler(new DiscardOldestPolicy());  // Newer matters
        executor.initialize();
        return executor;
    }
    
    /**
     * PAYMENT PROCESSING
     * Characteristics: Critical, can't lose tasks, latency sensitive
     */
    @Bean("paymentExecutor")
    @Profile("payment")
    public ThreadPoolTaskExecutor paymentExecutor() {
        int cores = Runtime.getRuntime().availableProcessors();
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(cores * 2);
        executor.setMaxPoolSize(cores * 4);
        executor.setQueueCapacity(50);      // Small queue - want fast processing
        executor.setKeepAliveSeconds(60);
        executor.setThreadNamePrefix("payment-");
        executor.setRejectedExecutionHandler(new AbortPolicy());  // Fail fast
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(300);  // Wait for in-flight payments
        executor.initialize();
        return executor;
    }
}
```

#### **Tuning Checklist**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PRODUCTION TUNING CHECKLIST                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ☐ SIZING                                                                   │
│    ├─ Profile your workload (CPU vs I/O bound)                             │
│    ├─ Calculate optimal thread count using formulas                        │
│    ├─ Start conservative, scale up based on metrics                        │
│    └─ Consider downstream dependencies (DB pool, external APIs)            │
│                                                                             │
│  ☐ QUEUE CONFIGURATION                                                      │
│    ├─ Bounded queue for predictable memory usage                           │
│    ├─ Size based on acceptable wait time × throughput                      │
│    └─ Monitor queue depth - alert if consistently high                     │
│                                                                             │
│  ☐ REJECTION HANDLING                                                       │
│    ├─ CallerRunsPolicy for most cases (backpressure)                       │
│    ├─ AbortPolicy for critical operations that must not be delayed         │
│    ├─ Custom handler for monitoring and alerting                           │
│    └─ Test rejection scenarios under load                                  │
│                                                                             │
│  ☐ SHUTDOWN BEHAVIOR                                                        │
│    ├─ setWaitForTasksToCompleteOnShutdown(true)                            │
│    ├─ Set appropriate awaitTerminationSeconds                              │
│    └─ Handle interrupted tasks gracefully                                  │
│                                                                             │
│  ☐ MONITORING                                                               │
│    ├─ Export metrics to Prometheus/Grafana                                 │
│    ├─ Alert on: High active count, queue depth, rejections                 │
│    ├─ Log slow tasks and exceptions                                        │
│    └─ Track task completion times                                          │
│                                                                             │
│  ☐ TESTING                                                                  │
│    ├─ Load test with 2-3x expected peak traffic                            │
│    ├─ Test graceful shutdown behavior                                      │
│    ├─ Test rejection policy behavior                                       │
│    └─ Verify MDC/context propagation                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Quick Reference: Pool Configuration by Pattern**

| Pattern | Core | Max | Queue | Rejection Policy | Use Case |
|---------|------|-----|-------|------------------|----------|
| **Fast I/O** | CPUs × 2 | CPUs × 4 | 100-500 | CallerRuns | API calls, DB queries |
| **CPU Intensive** | CPUs | CPUs + 1 | 50-100 | CallerRuns | Calculations, parsing |
| **Background Jobs** | 2-5 | 10-20 | 1000+ | CallerRuns | Emails, reports |
| **Critical Ops** | CPUs × 2 | CPUs × 4 | 25-50 | Abort | Payments, orders |
| **Event Processing** | 5-10 | 50-100 | 500-1000 | DiscardOldest | Events, notifications |
| **Batch Processing** | 2-5 | 10-20 | 10000+ | CallerRuns | Data imports, exports |

---

## Summary - Sections 3 & 4

| Topic | Key Points |
|-------|------------|
| **@EnableAsync** | Enables proxy-based async; method calls from same class won't work |
| **@Async Return Types** | void (fire-forget), Future (basic), CompletableFuture (recommended) |
| **Exception Handling** | AsyncUncaughtExceptionHandler for void; exceptionally()/handle() for CompletableFuture |
| **Core Pool Size** | Always-alive threads; set based on steady-state load |
| **Max Pool Size** | Peak capacity; only reached after queue is full |
| **Queue Capacity** | Buffer for burst traffic; bounded recommended |
| **Rejection Policies** | CallerRunsPolicy for most; AbortPolicy for critical ops |
| **Production Tuning** | Profile workload, monitor metrics, test under load |

---

## 5. Using CompletableFuture in Spring Boot

`CompletableFuture` is Java's powerful abstraction for asynchronous programming, introduced in Java 8. It provides a rich API for composing, combining, and handling asynchronous operations without blocking threads.

### 5.1 CompletableFuture Fundamentals

#### **What is CompletableFuture?**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPLETABLEFUTURE OVERVIEW                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CompletableFuture<T> = Future<T> + CompletionStage<T>                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         CompletableFuture                           │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                     │   │
│  │  From Future:                    From CompletionStage:              │   │
│  │  ────────────                    ─────────────────────              │   │
│  │  • get() - blocking wait         • thenApply() - transform         │   │
│  │  • isDone() - check completion   • thenCompose() - chain           │   │
│  │  • cancel() - cancel task        • thenCombine() - merge           │   │
│  │                                  • exceptionally() - handle errors │   │
│  │                                  • thenAccept() - consume          │   │
│  │                                  • allOf()/anyOf() - combine       │   │
│  │                                                                     │   │
│  │  Additional CompletableFuture methods:                              │   │
│  │  ─────────────────────────────────────                              │   │
│  │  • complete() - manually complete                                   │   │
│  │  • completeExceptionally() - complete with error                    │   │
│  │  • supplyAsync() - start async computation                          │   │
│  │  • runAsync() - start async action                                  │   │
│  │  • orTimeout() - timeout handling (Java 9+)                         │   │
│  │  • completeOnTimeout() - default on timeout (Java 9+)               │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Creating CompletableFutures**

```java
@Service
@Slf4j
public class CompletableFutureBasics {
    
    // 1. Create already completed future
    public CompletableFuture<String> alreadyComplete() {
        return CompletableFuture.completedFuture("immediate result");
    }
    
    // 2. Create failed future
    public CompletableFuture<String> alreadyFailed() {
        return CompletableFuture.failedFuture(new RuntimeException("error"));
    }
    
    // 3. Run async with return value (uses ForkJoinPool.commonPool())
    public CompletableFuture<String> supplyAsyncDefault() {
        return CompletableFuture.supplyAsync(() -> {
            log.info("Running on: {}", Thread.currentThread().getName());
            return "computed result";
        });
    }
    
    // 4. Run async with custom executor (RECOMMENDED for production)
    @Autowired
    private Executor asyncExecutor;
    
    public CompletableFuture<String> supplyAsyncCustomExecutor() {
        return CompletableFuture.supplyAsync(() -> {
            log.info("Running on: {}", Thread.currentThread().getName());
            return "computed result";
        }, asyncExecutor);  // Use Spring-managed executor
    }
    
    // 5. Run async without return value
    public CompletableFuture<Void> runAsyncNoReturn() {
        return CompletableFuture.runAsync(() -> {
            log.info("Fire and forget task");
        }, asyncExecutor);
    }
    
    // 6. Manual completion
    public CompletableFuture<String> manualCompletion() {
        CompletableFuture<String> future = new CompletableFuture<>();
        
        // Complete later (e.g., from callback)
        someAsyncOperation(result -> {
            if (result.isSuccess()) {
                future.complete(result.getData());
            } else {
                future.completeExceptionally(result.getError());
            }
        });
        
        return future;
    }
}
```

#### **Method Categories**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPLETABLEFUTURE METHOD CATEGORIES                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  TRANSFORMATION (Returns new CompletableFuture<R>)                          │
│  ──────────────────────────────────────────────────                         │
│  thenApply(Function<T,R>)       │ Transform result synchronously            │
│  thenApplyAsync(Function<T,R>)  │ Transform result asynchronously           │
│                                                                             │
│  COMPOSITION (Flatten nested futures)                                       │
│  ────────────────────────────────────                                       │
│  thenCompose(Function<T,CompletableFuture<R>>)  │ Chain dependent async     │
│  thenComposeAsync(...)                          │ calls                     │
│                                                                             │
│  CONSUMPTION (Returns CompletableFuture<Void>)                              │
│  ──────────────────────────────────────────────                             │
│  thenAccept(Consumer<T>)        │ Consume result, no return                 │
│  thenRun(Runnable)              │ Run action, ignore result                 │
│                                                                             │
│  COMBINATION (Merge multiple futures)                                       │
│  ────────────────────────────────────                                       │
│  thenCombine(CF<U>, BiFunction<T,U,R>)  │ Combine 2 results                 │
│  allOf(CF<?>...)                         │ Wait for all to complete         │
│  anyOf(CF<?>...)                         │ First to complete wins           │
│                                                                             │
│  ERROR HANDLING                                                             │
│  ──────────────                                                             │
│  exceptionally(Function<Throwable,T>)    │ Recover from exception           │
│  handle(BiFunction<T,Throwable,R>)       │ Handle success or failure        │
│  whenComplete(BiConsumer<T,Throwable>)   │ Side effect, doesn't change      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 5.2 Asynchronous Service Calls

#### **Basic Service Pattern**

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class ProductService {
    
    private final ProductRepository productRepository;
    private final Executor asyncExecutor;
    private final RestTemplate restTemplate;
    
    /**
     * Async database call
     */
    public CompletableFuture<Product> findProductById(Long id) {
        return CompletableFuture.supplyAsync(() -> {
            log.info("Fetching product {} on thread: {}", id, Thread.currentThread().getName());
            return productRepository.findById(id)
                .orElseThrow(() -> new ProductNotFoundException(id));
        }, asyncExecutor);
    }
    
    /**
     * Async external API call
     */
    public CompletableFuture<ProductPrice> fetchPriceFromSupplier(Long productId) {
        return CompletableFuture.supplyAsync(() -> {
            log.info("Calling supplier API on thread: {}", Thread.currentThread().getName());
            String url = "https://supplier-api.com/prices/" + productId;
            return restTemplate.getForObject(url, ProductPrice.class);
        }, asyncExecutor);
    }
    
    /**
     * Async with WebClient (non-blocking) - Preferred for external calls
     */
    @Autowired
    private WebClient webClient;
    
    public CompletableFuture<ProductPrice> fetchPriceReactive(Long productId) {
        return webClient.get()
            .uri("/prices/{id}", productId)
            .retrieve()
            .bodyToMono(ProductPrice.class)
            .toFuture();  // Convert Mono to CompletableFuture
    }
}
```

#### **Sequential vs Parallel Execution**

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class OrderEnrichmentService {
    
    private final UserService userService;
    private final ProductService productService;
    private final ShippingService shippingService;
    private final DiscountService discountService;
    
    /**
     * SEQUENTIAL: Each call waits for previous (SLOW)
     * Total time: sum of all calls
     */
    public EnrichedOrder enrichOrderSequential(Order order) {
        long start = System.currentTimeMillis();
        
        // Each blocking call
        User user = userService.findById(order.getUserId());           // 100ms
        Product product = productService.find(order.getProductId());    // 150ms
        ShippingInfo shipping = shippingService.calculate(order);       // 200ms
        Discount discount = discountService.findApplicable(order);      // 100ms
        
        log.info("Sequential took: {}ms", System.currentTimeMillis() - start);
        // Total: ~550ms
        
        return new EnrichedOrder(order, user, product, shipping, discount);
    }
    
    /**
     * PARALLEL: All calls execute simultaneously (FAST)
     * Total time: max of all calls
     */
    public CompletableFuture<EnrichedOrder> enrichOrderParallel(Order order) {
        long start = System.currentTimeMillis();
        
        // Start all async operations immediately
        CompletableFuture<User> userFuture = userService.findByIdAsync(order.getUserId());
        CompletableFuture<Product> productFuture = productService.findAsync(order.getProductId());
        CompletableFuture<ShippingInfo> shippingFuture = shippingService.calculateAsync(order);
        CompletableFuture<Discount> discountFuture = discountService.findApplicableAsync(order);
        
        // Combine all results
        return CompletableFuture.allOf(userFuture, productFuture, shippingFuture, discountFuture)
            .thenApply(v -> {
                log.info("Parallel took: {}ms", System.currentTimeMillis() - start);
                // Total: ~200ms (slowest call)
                
                return new EnrichedOrder(
                    order,
                    userFuture.join(),
                    productFuture.join(),
                    shippingFuture.join(),
                    discountFuture.join()
                );
            });
    }
}
```

#### **Visualization: Sequential vs Parallel**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SEQUENTIAL vs PARALLEL EXECUTION                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SEQUENTIAL (Blocking):                                                     │
│  ──────────────────────                                                     │
│                                                                             │
│  Time:  0ms      100ms     250ms      450ms      550ms                      │
│          │         │         │          │          │                        │
│  Thread: [==User==][===Product===][====Shipping====][==Discount==]          │
│                                                                             │
│  Total: 100 + 150 + 200 + 100 = 550ms                                       │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  PARALLEL (Non-blocking):                                                   │
│  ────────────────────────                                                   │
│                                                                             │
│  Time:  0ms                              200ms                              │
│          │                                 │                                │
│  T1:     [========User (100ms)========]   │                                │
│  T2:     [==========Product (150ms)==========]                             │
│  T3:     [============Shipping (200ms)============]                        │
│  T4:     [========Discount (100ms)========]                                │
│                                            │                                │
│                                            ▼                                │
│                                     Aggregate Results                       │
│                                                                             │
│  Total: max(100, 150, 200, 100) = 200ms  (63% faster!)                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 5.3 Combining Multiple Async Calls

#### **Pattern 1: thenCombine - Combine Two Results**

```java
@Service
@Slf4j
public class CombinePatternService {
    
    @Autowired
    private Executor asyncExecutor;
    
    /**
     * thenCombine: Combine results of two independent futures
     */
    public CompletableFuture<PriceComparison> comparePrices(Long productId) {
        
        CompletableFuture<BigDecimal> ourPrice = CompletableFuture.supplyAsync(
            () -> priceService.getOurPrice(productId), asyncExecutor);
        
        CompletableFuture<BigDecimal> competitorPrice = CompletableFuture.supplyAsync(
            () -> priceService.getCompetitorPrice(productId), asyncExecutor);
        
        // Combine when BOTH complete
        return ourPrice.thenCombine(competitorPrice, (our, competitor) -> {
            BigDecimal difference = competitor.subtract(our);
            boolean weAreCheaper = difference.compareTo(BigDecimal.ZERO) > 0;
            return new PriceComparison(our, competitor, difference, weAreCheaper);
        });
    }
    
    /**
     * Multiple thenCombine chained
     */
    public CompletableFuture<FullPriceAnalysis> fullPriceAnalysis(Long productId) {
        
        CompletableFuture<BigDecimal> supplierA = getFromSupplierA(productId);
        CompletableFuture<BigDecimal> supplierB = getFromSupplierB(productId);
        CompletableFuture<BigDecimal> supplierC = getFromSupplierC(productId);
        
        return supplierA
            .thenCombine(supplierB, (a, b) -> new PricePair(a, b))
            .thenCombine(supplierC, (pair, c) -> 
                new FullPriceAnalysis(pair.a(), pair.b(), c));
    }
}
```

#### **Pattern 2: allOf - Wait for All**

```java
@Service
@Slf4j
public class AllOfPatternService {
    
    /**
     * allOf: Wait for all futures to complete
     * Note: allOf returns CompletableFuture<Void>
     */
    public CompletableFuture<Dashboard> buildDashboard(String userId) {
        
        CompletableFuture<UserProfile> profileF = userService.getProfile(userId);
        CompletableFuture<List<Order>> ordersF = orderService.getRecentOrders(userId);
        CompletableFuture<List<Notification>> notificationsF = notificationService.getUnread(userId);
        CompletableFuture<UserStats> statsF = analyticsService.getStats(userId);
        CompletableFuture<List<Recommendation>> recommendationsF = recommenderService.get(userId);
        
        // Wait for ALL to complete
        return CompletableFuture.allOf(
                profileF, ordersF, notificationsF, statsF, recommendationsF)
            .thenApply(v -> {
                // All futures are now complete - .join() won't block
                return Dashboard.builder()
                    .profile(profileF.join())
                    .recentOrders(ordersF.join())
                    .notifications(notificationsF.join())
                    .stats(statsF.join())
                    .recommendations(recommendationsF.join())
                    .build();
            });
    }
    
    /**
     * allOf with typed results using utility method
     */
    public CompletableFuture<List<ProductPrice>> getAllPrices(List<Long> productIds) {
        
        List<CompletableFuture<ProductPrice>> futures = productIds.stream()
            .map(id -> priceService.getPriceAsync(id))
            .collect(Collectors.toList());
        
        // Convert List<CompletableFuture<T>> to CompletableFuture<List<T>>
        return sequence(futures);
    }
    
    /**
     * Utility: Convert List<CF<T>> to CF<List<T>>
     */
    public static <T> CompletableFuture<List<T>> sequence(List<CompletableFuture<T>> futures) {
        return CompletableFuture.allOf(futures.toArray(new CompletableFuture[0]))
            .thenApply(v -> futures.stream()
                .map(CompletableFuture::join)
                .collect(Collectors.toList()));
    }
}
```

#### **Pattern 3: anyOf - First to Complete Wins**

```java
@Service
@Slf4j
public class AnyOfPatternService {
    
    /**
     * anyOf: Return first completed result
     * Use case: Redundant calls to multiple services, use fastest response
     */
    public CompletableFuture<Price> getFastestPrice(Long productId) {
        
        CompletableFuture<Price> supplier1 = fetchFromSupplier1(productId);
        CompletableFuture<Price> supplier2 = fetchFromSupplier2(productId);
        CompletableFuture<Price> supplier3 = fetchFromSupplier3(productId);
        CompletableFuture<Price> cache = fetchFromCache(productId);
        
        return CompletableFuture.anyOf(supplier1, supplier2, supplier3, cache)
            .thenApply(result -> (Price) result);  // anyOf returns Object
    }
    
    /**
     * anyOf with timeout fallback
     */
    public CompletableFuture<Price> getPriceWithFallback(Long productId) {
        
        CompletableFuture<Price> livePrice = fetchLivePrice(productId);
        
        // Create timeout future that completes with cached price after 2 seconds
        CompletableFuture<Price> timeoutFallback = CompletableFuture.supplyAsync(() -> {
            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            return getCachedPrice(productId);
        });
        
        // Return whichever completes first
        return CompletableFuture.anyOf(livePrice, timeoutFallback)
            .thenApply(result -> (Price) result);
    }
    
    /**
     * Better timeout handling with orTimeout (Java 9+)
     */
    public CompletableFuture<Price> getPriceWithTimeout(Long productId) {
        return fetchLivePrice(productId)
            .orTimeout(2, TimeUnit.SECONDS)  // Throws TimeoutException
            .exceptionally(ex -> {
                log.warn("Live price timed out, using cached", ex);
                return getCachedPrice(productId);
            });
    }
    
    /**
     * completeOnTimeout - Complete with default value (Java 9+)
     */
    public CompletableFuture<Price> getPriceWithDefault(Long productId) {
        return fetchLivePrice(productId)
            .completeOnTimeout(
                getCachedPrice(productId),  // Default value
                2,                          // Timeout duration
                TimeUnit.SECONDS            // Time unit
            );
    }
}
```

#### **Pattern 4: thenCompose - Chain Dependent Calls**

```java
@Service
@Slf4j
public class ComposePatternService {
    
    /**
     * thenCompose: Chain async calls that depend on previous result
     * Avoids nested CompletableFuture<CompletableFuture<T>>
     */
    public CompletableFuture<OrderConfirmation> processOrder(OrderRequest request) {
        
        return validateOrder(request)                           // CF<ValidatedOrder>
            .thenCompose(validated -> reserveInventory(validated))  // CF<Reservation>
            .thenCompose(reservation -> processPayment(reservation)) // CF<Payment>
            .thenCompose(payment -> createOrder(payment))           // CF<Order>
            .thenCompose(order -> sendConfirmation(order));         // CF<Confirmation>
    }
    
    /**
     * thenCompose vs thenApply
     */
    public void composeVsApply() {
        
        // thenApply: Transform value, returns same type
        CompletableFuture<String> nameF = getUserId()
            .thenApply(id -> "User-" + id);  // String -> String
        
        // thenCompose: Chain futures, flattens result
        CompletableFuture<User> userF = getUserId()
            .thenCompose(id -> getUserById(id));  // Long -> CF<User> -> User
        
        // If you used thenApply with async call:
        CompletableFuture<CompletableFuture<User>> nestedF = getUserId()
            .thenApply(id -> getUserById(id));  // Wrong! Nested futures
    }
    
    /**
     * Complex composition example: E-commerce checkout flow
     */
    public CompletableFuture<CheckoutResult> checkout(Cart cart, PaymentInfo payment) {
        
        return CompletableFuture.supplyAsync(() -> cart, asyncExecutor)
            // Validate cart
            .thenCompose(c -> validateCart(c))
            // Check inventory for all items
            .thenCompose(validCart -> 
                checkInventoryForAll(validCart.getItems())
                    .thenApply(inventory -> new CartWithInventory(validCart, inventory)))
            // Calculate totals with discounts
            .thenCompose(cwi -> calculateTotals(cwi))
            // Process payment
            .thenCompose(totals -> processPayment(totals, payment))
            // Create order
            .thenCompose(paymentResult -> {
                if (paymentResult.isSuccess()) {
                    return createOrder(cart, paymentResult);
                } else {
                    return CompletableFuture.failedFuture(
                        new PaymentFailedException(paymentResult.getReason()));
                }
            })
            // Send notifications (parallel, don't wait)
            .thenApply(order -> {
                sendOrderConfirmationAsync(order);  // Fire and forget
                return new CheckoutResult(order, "SUCCESS");
            });
    }
}
```

#### **Combination Patterns Summary**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMBINATION PATTERNS SUMMARY                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Pattern          │ Use Case                        │ Result                │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  thenCombine()    │ Merge 2 independent results     │ CF<R>                 │
│  CF1 + CF2 → R    │ "Get user AND orders"           │                       │
│                                                                             │
│  allOf()          │ Wait for all, aggregate results │ CF<Void>              │
│  [CF...] → Void   │ "Get 5 things, combine when done"│ (use join to get)    │
│                                                                             │
│  anyOf()          │ First completed wins            │ CF<Object>            │
│  [CF...] → first  │ "Fastest of 3 suppliers"        │ (needs cast)          │
│                                                                             │
│  thenCompose()    │ Chain dependent async calls     │ CF<R>                 │
│  CF → CF → CF     │ "Validate → Reserve → Pay"      │ (flattens nesting)    │
│                                                                             │
│  thenApply()      │ Transform result synchronously  │ CF<R>                 │
│  CF → R           │ "Convert DTO to Entity"         │                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 5.4 Exception Handling Patterns

#### **Exception Flow in CompletableFuture**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    EXCEPTION PROPAGATION                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Stage 1          Stage 2          Stage 3          Result                  │
│  ──────────────────────────────────────────────────────────────────────     │
│                                                                             │
│  SUCCESS PATH:                                                              │
│  ┌───────┐       ┌───────┐       ┌───────┐       ┌───────┐                 │
│  │ Fetch │──────►│ Parse │──────►│ Save  │──────►│ Done  │                 │
│  └───────┘       └───────┘       └───────┘       └───────┘                 │
│                                                                             │
│  EXCEPTION PATH (no handling):                                              │
│  ┌───────┐       ┌───────┐       ┌───────┐       ┌───────┐                 │
│  │ Fetch │───X──►│ Skip  │──────►│ Skip  │──────►│ Error │                 │
│  └───────┘  ↓    └───────┘       └───────┘       └───────┘                 │
│         Exception propagates through all stages─────────────►               │
│                                                                             │
│  EXCEPTION PATH (with exceptionally):                                       │
│  ┌───────┐       ┌────────────┐  ┌───────┐       ┌───────┐                 │
│  │ Fetch │───X──►│exceptionally│─►│ Save  │──────►│ Done  │                 │
│  └───────┘  ↓    │(recover)   │  └───────┘       └───────┘                 │
│         Exception │            │                                            │
│         handled   └────────────┘                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Pattern 1: exceptionally - Recover from Exception**

```java
@Service
@Slf4j
public class ExceptionHandlingService {
    
    /**
     * exceptionally: Recover with fallback value
     */
    public CompletableFuture<Price> getPriceWithFallback(Long productId) {
        return fetchLivePrice(productId)
            .exceptionally(ex -> {
                log.warn("Failed to fetch live price for {}: {}", productId, ex.getMessage());
                return getCachedPrice(productId);  // Fallback to cache
            });
    }
    
    /**
     * exceptionally with specific exception types
     */
    public CompletableFuture<User> getUserWithRecovery(Long userId) {
        return fetchUser(userId)
            .exceptionally(ex -> {
                Throwable cause = ex.getCause() != null ? ex.getCause() : ex;
                
                if (cause instanceof UserNotFoundException) {
                    log.info("User {} not found, returning guest", userId);
                    return User.guest();
                } else if (cause instanceof ServiceUnavailableException) {
                    log.warn("User service down, using cache");
                    return userCache.get(userId);
                } else {
                    log.error("Unexpected error fetching user", cause);
                    throw new RuntimeException("Failed to get user", cause);
                }
            });
    }
    
    /**
     * Re-throw as different exception
     */
    public CompletableFuture<Order> getOrder(Long orderId) {
        return fetchOrder(orderId)
            .exceptionally(ex -> {
                throw new OrderServiceException("Failed to fetch order " + orderId, ex);
            });
    }
}
```

#### **Pattern 2: handle - Handle Success and Failure**

```java
@Service
@Slf4j
public class HandlePatternService {
    
    /**
     * handle: Single handler for both success and failure
     * Always called, regardless of exception
     */
    public CompletableFuture<ApiResponse<User>> getUserSafely(Long userId) {
        return fetchUser(userId)
            .handle((user, ex) -> {
                if (ex != null) {
                    log.error("Failed to fetch user {}", userId, ex);
                    return ApiResponse.error("User fetch failed: " + ex.getMessage());
                }
                return ApiResponse.success(user);
            });
    }
    
    /**
     * handle with metrics recording
     */
    public CompletableFuture<Price> getPriceWithMetrics(Long productId) {
        long startTime = System.currentTimeMillis();
        
        return fetchPrice(productId)
            .handle((price, ex) -> {
                long duration = System.currentTimeMillis() - startTime;
                
                if (ex != null) {
                    metrics.recordFailure("price_fetch", duration, ex.getClass().getSimpleName());
                    return Price.unknown();
                } else {
                    metrics.recordSuccess("price_fetch", duration);
                    return price;
                }
            });
    }
    
    /**
     * handle for transformation with error context
     */
    public CompletableFuture<Result<Order>> processOrderSafely(OrderRequest request) {
        return processOrder(request)
            .handle((order, ex) -> {
                if (ex != null) {
                    Throwable cause = unwrapException(ex);
                    
                    return Result.<Order>failure()
                        .error(cause.getMessage())
                        .errorCode(mapToErrorCode(cause))
                        .timestamp(Instant.now())
                        .build();
                }
                return Result.success(order);
            });
    }
    
    private Throwable unwrapException(Throwable ex) {
        // CompletionException wraps the actual cause
        if (ex instanceof CompletionException && ex.getCause() != null) {
            return ex.getCause();
        }
        return ex;
    }
}
```

#### **Pattern 3: whenComplete - Side Effects**

```java
@Service
@Slf4j
public class WhenCompleteService {
    
    /**
     * whenComplete: Execute side effect, don't change result
     * Original result/exception propagates unchanged
     */
    public CompletableFuture<Order> processOrderWithLogging(OrderRequest request) {
        return processOrder(request)
            .whenComplete((order, ex) -> {
                // Logging side effect - doesn't change result
                if (ex != null) {
                    log.error("Order processing failed for request: {}", request.getId(), ex);
                    auditService.logFailure(request, ex);
                } else {
                    log.info("Order {} processed successfully", order.getId());
                    auditService.logSuccess(order);
                }
            });
        // Original order or exception continues to flow
    }
    
    /**
     * whenComplete for cleanup
     */
    public CompletableFuture<Report> generateReportWithCleanup(ReportRequest request) {
        TempFile tempFile = createTempFile();
        
        return generateReport(request, tempFile)
            .whenComplete((report, ex) -> {
                // Always cleanup temp file
                try {
                    tempFile.delete();
                } catch (IOException e) {
                    log.warn("Failed to delete temp file", e);
                }
            });
    }
}
```

#### **Pattern 4: Combining Exception Handlers**

```java
@Service
@Slf4j
public class CombinedExceptionHandling {
    
    /**
     * Multi-level exception handling
     */
    public CompletableFuture<OrderResult> robustOrderProcessing(OrderRequest request) {
        return validateOrder(request)
            .thenCompose(validated -> {
                return reserveInventory(validated)
                    .exceptionally(ex -> {
                        // Try alternative warehouse if primary fails
                        log.warn("Primary warehouse failed, trying backup");
                        return reserveFromBackupWarehouse(validated).join();
                    });
            })
            .thenCompose(reservation -> {
                return processPayment(reservation)
                    .handle((payment, ex) -> {
                        if (ex != null) {
                            // Record failed payment attempt
                            paymentAudit.recordFailure(reservation, ex);
                            throw new PaymentFailedException(ex);
                        }
                        paymentAudit.recordSuccess(payment);
                        return payment;
                    });
            })
            .thenCompose(payment -> createOrder(payment))
            .whenComplete((result, ex) -> {
                // Final logging regardless of success/failure
                if (ex != null) {
                    log.error("Order failed: {}", request.getId(), ex);
                    notifyOperations(request, ex);
                } else {
                    log.info("Order completed: {}", result.getOrderId());
                }
            })
            .exceptionally(ex -> {
                // Final fallback - return failed result instead of exception
                return OrderResult.failed(request.getId(), ex.getMessage());
            });
    }
}
```

#### **Exception Handling Comparison**

| Method | When Called | Returns | Changes Result? | Use Case |
|--------|-------------|---------|-----------------|----------|
| `exceptionally()` | Only on exception | T | Yes (can recover) | Fallback values |
| `handle()` | Always | R | Yes (transform) | Wrap in Result/Response |
| `whenComplete()` | Always | Void (original propagates) | No | Logging, cleanup |

---

### 5.5 Performance Benefits and Patterns

#### **Performance Comparison: Real Example**

```java
@Service
@Slf4j
public class PerformanceComparisonService {
    
    /**
     * Production scenario: Build product page with data from multiple sources
     */
    
    // BLOCKING APPROACH - Sequential
    public ProductPage buildProductPageBlocking(Long productId) {
        long start = System.currentTimeMillis();
        
        Product product = productService.getProduct(productId);           // 50ms
        List<Review> reviews = reviewService.getReviews(productId);        // 100ms
        List<Product> related = recommendationService.getRelated(productId);// 150ms
        Inventory inventory = inventoryService.getStock(productId);        // 30ms
        PriceInfo price = pricingService.getPrice(productId);              // 80ms
        Seller seller = sellerService.getSellerInfo(product.getSellerId());// 40ms
        
        log.info("Blocking build took: {}ms", System.currentTimeMillis() - start);
        // Total: 50+100+150+30+80+40 = 450ms
        
        return new ProductPage(product, reviews, related, inventory, price, seller);
    }
    
    // NON-BLOCKING APPROACH - Parallel
    public CompletableFuture<ProductPage> buildProductPageAsync(Long productId) {
        long start = System.currentTimeMillis();
        
        // Start all independent calls immediately
        var productF = productService.getProductAsync(productId);
        var reviewsF = reviewService.getReviewsAsync(productId);
        var relatedF = recommendationService.getRelatedAsync(productId);
        var inventoryF = inventoryService.getStockAsync(productId);
        var priceF = pricingService.getPriceAsync(productId);
        
        // Seller depends on product - chain it
        var sellerF = productF.thenCompose(product -> 
            sellerService.getSellerInfoAsync(product.getSellerId()));
        
        return CompletableFuture.allOf(
                productF, reviewsF, relatedF, inventoryF, priceF, sellerF)
            .thenApply(v -> {
                log.info("Async build took: {}ms", System.currentTimeMillis() - start);
                // Total: max(50+40, 100, 150, 30, 80) = 150ms (67% faster!)
                
                return new ProductPage(
                    productF.join(),
                    reviewsF.join(),
                    relatedF.join(),
                    inventoryF.join(),
                    priceF.join(),
                    sellerF.join()
                );
            });
    }
}
```

#### **Performance Visualization**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PERFORMANCE COMPARISON: PRODUCT PAGE                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BLOCKING (450ms total):                                                    │
│  ───────────────────────                                                    │
│                                                                             │
│  Time: 0    50   150      300      330    410    450ms                      │
│        │     │     │        │        │      │      │                        │
│        [Prod][Reviews][Recommendations][Inv][Price][Seller]                 │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  ASYNC (150ms total):                                                       │
│  ────────────────────                                                       │
│                                                                             │
│  Time: 0                                           150ms                    │
│        │                                             │                      │
│  T1:   [====Product 50ms====][Seller 40ms]           │  (depends on prod)   │
│  T2:   [========Reviews 100ms========]               │                      │
│  T3:   [===========Recommendations 150ms===========] │  (slowest)           │
│  T4:   [=Inv 30ms=]                                  │                      │
│  T5:   [======Price 80ms======]                      │                      │
│                                                      ▼                      │
│                                               Aggregate                     │
│                                                                             │
│  Speedup: 450ms → 150ms = 3x faster (67% reduction)                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Best Practices for CompletableFuture**

```java
@Configuration
public class CompletableFutureBestPractices {
    
    /**
     * BEST PRACTICE 1: Always use custom executor in production
     * Default ForkJoinPool.commonPool() is shared across JVM
     */
    @Bean("cfExecutor")
    public Executor completableFutureExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(20);
        executor.setMaxPoolSize(100);
        executor.setQueueCapacity(500);
        executor.setThreadNamePrefix("cf-");
        executor.initialize();
        return executor;
    }
}

@Service
@RequiredArgsConstructor
public class BestPracticesService {
    
    private final Executor cfExecutor;
    
    /**
     * BEST PRACTICE 2: Pass executor to *Async methods
     */
    public CompletableFuture<Result> goodPattern() {
        return CompletableFuture.supplyAsync(() -> fetchData(), cfExecutor)
            .thenApplyAsync(data -> transform(data), cfExecutor)
            .thenComposeAsync(transformed -> save(transformed), cfExecutor);
    }
    
    /**
     * BEST PRACTICE 3: Set timeouts to prevent hanging
     */
    public CompletableFuture<Data> withTimeout() {
        return fetchData()
            .orTimeout(5, TimeUnit.SECONDS)
            .exceptionally(ex -> {
                if (ex instanceof TimeoutException) {
                    return Data.defaultValue();
                }
                throw new RuntimeException(ex);
            });
    }
    
    /**
     * BEST PRACTICE 4: Don't block in async chains
     */
    public CompletableFuture<Result> avoidBlocking() {
        // BAD: Blocking inside async chain
        return CompletableFuture.supplyAsync(() -> {
            Data data = otherFuture.get();  // DON'T DO THIS - blocks thread
            return process(data);
        });
        
        // GOOD: Use composition
        return otherFuture.thenApply(data -> process(data));
    }
    
    /**
     * BEST PRACTICE 5: Handle exceptions at appropriate level
     */
    public CompletableFuture<Order> properExceptionHandling(OrderRequest request) {
        return validateOrder(request)
            .thenCompose(v -> reserveInventory(v))
            .exceptionally(ex -> {
                // Handle inventory-specific failure
                if (ex.getCause() instanceof InventoryException) {
                    return Reservation.backorder();
                }
                throw (RuntimeException) ex;
            })
            .thenCompose(res -> processPayment(res))
            .thenCompose(pay -> createOrder(pay))
            .exceptionally(ex -> {
                // Final catch-all
                log.error("Order failed", ex);
                return Order.failed(ex.getMessage());
            });
    }
    
    /**
     * BEST PRACTICE 6: Use join() over get() when inside composition
     */
    public CompletableFuture<Report> useJoinNotGet() {
        CompletableFuture<DataA> aF = fetchA();
        CompletableFuture<DataB> bF = fetchB();
        
        return CompletableFuture.allOf(aF, bF)
            .thenApply(v -> {
                // join() preferred inside thenApply - won't block (futures complete)
                // join() throws unchecked exception (easier to handle)
                DataA a = aF.join();
                DataB b = bF.join();
                return new Report(a, b);
            });
    }
}
```

#### **CompletableFuture Quick Reference**

| Method | Input | Output | Use When |
|--------|-------|--------|----------|
| `supplyAsync(Supplier)` | - | `CF<T>` | Start async task with result |
| `runAsync(Runnable)` | - | `CF<Void>` | Start async task no result |
| `thenApply(Function)` | `T` | `CF<R>` | Transform result sync |
| `thenApplyAsync(Function)` | `T` | `CF<R>` | Transform result async |
| `thenCompose(Function)` | `T` | `CF<R>` | Chain dependent async calls |
| `thenCombine(CF, BiFunction)` | `T`, `U` | `CF<R>` | Combine 2 futures |
| `allOf(CF...)` | - | `CF<Void>` | Wait for all |
| `anyOf(CF...)` | - | `CF<Object>` | First to complete |
| `exceptionally(Function)` | `Throwable` | `CF<T>` | Recover from error |
| `handle(BiFunction)` | `T`, `Throwable` | `CF<R>` | Handle success/failure |
| `orTimeout(long, TimeUnit)` | - | `CF<T>` | Timeout with exception |
| `completeOnTimeout(T, long, TimeUnit)` | - | `CF<T>` | Timeout with default |

---

## 6. Scheduling Tasks

Spring's scheduling support allows you to execute tasks at fixed intervals, with delays, or using cron expressions.

### 6.1 Enabling Scheduling Support

#### **Basic Setup**

```java
@Configuration
@EnableScheduling  // Enables Spring's scheduling support
public class SchedulingConfig {
    // Basic configuration - uses single-threaded scheduler by default
}

// Alternative: On main application class
@SpringBootApplication
@EnableScheduling
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

#### **What @EnableScheduling Does**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    @EnableScheduling PROCESSING                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Spring Container Startup:                                                  │
│  ─────────────────────────                                                  │
│                                                                             │
│  1. @EnableScheduling detected                                              │
│     │                                                                       │
│     ▼                                                                       │
│  2. ScheduledAnnotationBeanPostProcessor registered                         │
│     │                                                                       │
│     ▼                                                                       │
│  3. Scans all beans for @Scheduled methods                                  │
│     │                                                                       │
│     ▼                                                                       │
│  4. Creates ScheduledTaskRegistrar                                          │
│     │                                                                       │
│     ▼                                                                       │
│  5. Registers tasks with TaskScheduler                                      │
│     │                                                                       │
│     ▼                                                                       │
│  6. Default: Single-threaded scheduler (ConcurrentTaskScheduler)            │
│                                                                             │
│  ⚠️  WARNING: Default is SINGLE THREADED!                                   │
│      If Task A takes 10 seconds, Task B waits!                              │
│      Always configure a proper thread pool for production.                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 6.2 @Scheduled Annotation

#### **Schedule Types Overview**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    @SCHEDULED OPTIONS                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  fixedRate                                                                  │
│  ─────────                                                                  │
│  Time:     0s        5s        10s       15s       20s                      │
│            │         │         │         │         │                        │
│            [Task]    [Task]    [Task]    [Task]    [Task]                   │
│                                                                             │
│  Runs every 5 seconds regardless of task duration                           │
│  If task takes longer than interval, next execution starts immediately      │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  fixedDelay                                                                 │
│  ──────────                                                                 │
│  Time:     0s   3s        8s   11s       16s  19s                          │
│            │    │         │    │         │    │                            │
│            [Task]─5s─►    [Task]─5s─►    [Task]─5s─►                        │
│            └───┘          └───┘          └───┘                              │
│            3s task        3s task        3s task                            │
│                                                                             │
│  Waits 5 seconds AFTER task completes before starting next                  │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  cron                                                                       │
│  ────                                                                       │
│  Executes at specific times based on cron expression                        │
│  "0 0 * * * *" = Every hour at minute 0                                     │
│  "0 0 9 * * MON-FRI" = 9:00 AM on weekdays                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **fixedRate Examples**

```java
@Service
@Slf4j
public class FixedRateScheduler {
    
    /**
     * fixedRate: Execute every 5 seconds
     * Starts immediately when application starts
     */
    @Scheduled(fixedRate = 5000)
    public void runEvery5Seconds() {
        log.info("[fixedRate] Task running at: {}", LocalDateTime.now());
    }
    
    /**
     * fixedRate with initialDelay: Wait before first execution
     */
    @Scheduled(fixedRate = 10000, initialDelay = 30000)
    public void runAfterInitialDelay() {
        // First run: 30 seconds after startup
        // Subsequent runs: every 10 seconds
        log.info("[fixedRate+initialDelay] Running...");
    }
    
    /**
     * Using TimeUnit for clarity
     */
    @Scheduled(fixedRate = 1, timeUnit = TimeUnit.MINUTES)
    public void runEveryMinute() {
        log.info("[fixedRate-minutes] Running every minute");
    }
    
    /**
     * Using properties file values
     */
    @Scheduled(fixedRateString = "${app.scheduler.metrics.rate:60000}")
    public void collectMetrics() {
        // Rate configurable via application.properties
        metricsCollector.collect();
    }
    
    /**
     * WARNING: fixedRate with long-running task
     * If task takes longer than rate, executions pile up
     */
    @Scheduled(fixedRate = 5000)
    public void longRunningTask() {
        log.info("Starting long task...");
        try {
            Thread.sleep(8000);  // Takes 8 seconds, rate is 5 seconds
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        log.info("Long task complete");
        // Next execution starts immediately (was already due)
    }
}
```

#### **fixedDelay Examples**

```java
@Service
@Slf4j
public class FixedDelayScheduler {
    
    /**
     * fixedDelay: Wait 5 seconds AFTER task completes
     */
    @Scheduled(fixedDelay = 5000)
    public void runWithDelay() {
        log.info("[fixedDelay] Starting at: {}", LocalDateTime.now());
        // After this method completes, wait 5 seconds, then run again
    }
    
    /**
     * fixedDelay with initialDelay
     */
    @Scheduled(fixedDelay = 10000, initialDelay = 5000)
    public void delayWithInitial() {
        log.info("[fixedDelay+initialDelay] Running...");
    }
    
    /**
     * Perfect for polling external systems
     * Ensures you don't overwhelm the system with requests
     */
    @Scheduled(fixedDelay = 1000)  // 1 second after completion
    public void pollExternalSystem() {
        log.info("Polling external system...");
        try {
            List<Event> events = externalApi.getNewEvents();
            events.forEach(eventProcessor::process);
        } catch (Exception e) {
            log.warn("Polling failed, will retry after delay", e);
        }
        // Even if this takes 30 seconds, next poll waits 1 second after
    }
    
    /**
     * Using properties
     */
    @Scheduled(fixedDelayString = "${app.scheduler.cleanup.delay:30000}")
    public void cleanup() {
        tempFileService.cleanupOldFiles();
    }
}
```

#### **Cron Expression Examples**

```java
@Service
@Slf4j
public class CronScheduler {
    
    /**
     * Cron format: second minute hour day-of-month month day-of-week
     *              0      0      *    *            *     *
     */
    
    // Every hour at minute 0
    @Scheduled(cron = "0 0 * * * *")
    public void everyHour() {
        log.info("Running hourly task");
    }
    
    // Every day at midnight
    @Scheduled(cron = "0 0 0 * * *")
    public void dailyMidnight() {
        log.info("Running daily midnight task");
        dailyReportService.generate();
    }
    
    // Every day at 9:00 AM
    @Scheduled(cron = "0 0 9 * * *")
    public void dailyMorning() {
        log.info("Good morning! Running 9 AM task");
        notificationService.sendDailyDigest();
    }
    
    // Every weekday (Mon-Fri) at 6:00 PM
    @Scheduled(cron = "0 0 18 * * MON-FRI")
    public void weekdayEvening() {
        log.info("End of business day task");
        dailySummaryService.generateAndSend();
    }
    
    // Every 15 minutes
    @Scheduled(cron = "0 */15 * * * *")
    public void every15Minutes() {
        log.info("15-minute check");
        healthCheckService.check();
    }
    
    // First day of every month at midnight
    @Scheduled(cron = "0 0 0 1 * *")
    public void monthlyTask() {
        log.info("Monthly billing run");
        billingService.generateMonthlyInvoices();
    }
    
    // Every Sunday at 2:00 AM
    @Scheduled(cron = "0 0 2 * * SUN")
    public void weeklyMaintenance() {
        log.info("Weekly maintenance");
        databaseService.optimize();
    }
    
    // With timezone
    @Scheduled(cron = "0 0 9 * * *", zone = "America/New_York")
    public void newYorkMorning() {
        log.info("9 AM in New York");
    }
    
    // Using property placeholder
    @Scheduled(cron = "${app.scheduler.report.cron:0 0 6 * * *}")
    public void configurableCron() {
        log.info("Configurable cron job");
    }
}
```

#### **Cron Expression Reference**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CRON EXPRESSION REFERENCE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Format: second minute hour day-of-month month day-of-week                  │
│                                                                             │
│  Field          Allowed Values      Special Characters                      │
│  ─────────────────────────────────────────────────────────────────────────  │
│  second         0-59                , - * /                                 │
│  minute         0-59                , - * /                                 │
│  hour           0-23                , - * /                                 │
│  day-of-month   1-31                , - * ? L W                             │
│  month          1-12 or JAN-DEC     , - * /                                 │
│  day-of-week    0-7 or SUN-SAT      , - * ? L #                             │
│                 (0 and 7 are Sunday)                                        │
│                                                                             │
│  Special Characters:                                                        │
│  ──────────────────                                                         │
│  *   = All values            "*" in hour = every hour                       │
│  ?   = No specific value     Use in day-of-month or day-of-week             │
│  -   = Range                 "10-12" = 10, 11, 12                           │
│  ,   = List                  "MON,WED,FRI" = those 3 days                   │
│  /   = Increments            "0/15" = 0, 15, 30, 45                         │
│  L   = Last                  "L" in day-of-month = last day of month        │
│  W   = Weekday               "15W" = nearest weekday to 15th                │
│  #   = Nth day               "2#3" = third Monday (day 2, occurrence 3)     │
│                                                                             │
│  COMMON EXAMPLES:                                                           │
│  ────────────────                                                           │
│  "0 0 * * * *"         Every hour                                           │
│  "0 */10 * * * *"      Every 10 minutes                                     │
│  "0 0 8-18 * * *"      Every hour from 8 AM to 6 PM                         │
│  "0 0 0 * * *"         Every day at midnight                                │
│  "0 0 9 * * MON-FRI"   9 AM on weekdays                                     │
│  "0 0 0 1 * *"         First day of every month                             │
│  "0 0 0 L * *"         Last day of every month                              │
│  "0 0 9 ? * 2#1"       First Monday of month at 9 AM                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 6.3 Thread Pool for Scheduled Tasks

#### **The Default Single-Thread Problem**

```java
/**
 * PROBLEM: Default scheduler is SINGLE THREADED!
 * 
 * If you have multiple @Scheduled methods and one blocks,
 * ALL other scheduled tasks are delayed.
 */
@Service
@Slf4j
public class ProblematicScheduler {
    
    @Scheduled(fixedRate = 1000)  // Every second
    public void fastTask() {
        log.info("Fast task - should run every second");
    }
    
    @Scheduled(fixedRate = 5000)  // Every 5 seconds
    public void slowTask() {
        log.info("Slow task starting...");
        try {
            Thread.sleep(30000);  // Takes 30 seconds!
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        log.info("Slow task complete");
    }
    
    // RESULT: fastTask is blocked for 30 seconds while slowTask runs!
}
```

#### **Configuring Thread Pool for Scheduling**

```java
@Configuration
@EnableScheduling
public class SchedulerConfig implements SchedulingConfigurer {
    
    @Override
    public void configureTasks(ScheduledTaskRegistrar taskRegistrar) {
        // Option 1: Simple thread pool
        taskRegistrar.setScheduler(scheduledExecutorService());
    }
    
    @Bean(destroyMethod = "shutdown")
    public ScheduledExecutorService scheduledExecutorService() {
        return Executors.newScheduledThreadPool(10);
    }
}

// Alternative: Using ThreadPoolTaskScheduler (recommended)
@Configuration
@EnableScheduling
public class SchedulerConfig {
    
    @Bean
    public TaskScheduler taskScheduler() {
        ThreadPoolTaskScheduler scheduler = new ThreadPoolTaskScheduler();
        scheduler.setPoolSize(10);  // Number of threads
        scheduler.setThreadNamePrefix("scheduled-");
        scheduler.setAwaitTerminationSeconds(60);
        scheduler.setWaitForTasksToCompleteOnShutdown(true);
        
        // Error handler for scheduled tasks
        scheduler.setErrorHandler(t -> {
            log.error("Scheduled task error", t);
            // Send alert, record metric, etc.
        });
        
        scheduler.initialize();
        return scheduler;
    }
}
```

#### **Configuration via Properties (Spring Boot 2.1+)**

```yaml
# application.yml
spring:
  task:
    scheduling:
      pool:
        size: 10  # Number of threads for @Scheduled tasks
      thread-name-prefix: "scheduling-"
      shutdown:
        await-termination: true
        await-termination-period: 60s
```

#### **Multiple Schedulers for Different Task Types**

```java
@Configuration
@EnableScheduling
public class MultiSchedulerConfig {
    
    /**
     * Primary scheduler for regular tasks
     */
    @Bean
    @Primary
    public TaskScheduler taskScheduler() {
        ThreadPoolTaskScheduler scheduler = new ThreadPoolTaskScheduler();
        scheduler.setPoolSize(5);
        scheduler.setThreadNamePrefix("scheduler-");
        scheduler.initialize();
        return scheduler;
    }
    
    /**
     * Dedicated scheduler for long-running tasks
     */
    @Bean("longRunningScheduler")
    public TaskScheduler longRunningScheduler() {
        ThreadPoolTaskScheduler scheduler = new ThreadPoolTaskScheduler();
        scheduler.setPoolSize(3);
        scheduler.setThreadNamePrefix("long-running-");
        scheduler.initialize();
        return scheduler;
    }
    
    /**
     * High-priority scheduler for critical tasks
     */
    @Bean("criticalScheduler")
    public TaskScheduler criticalScheduler() {
        ThreadPoolTaskScheduler scheduler = new ThreadPoolTaskScheduler();
        scheduler.setPoolSize(5);
        scheduler.setThreadNamePrefix("critical-");
        scheduler.initialize();
        return scheduler;
    }
}

// Using specific scheduler (not directly supported, need programmatic approach)
@Service
@RequiredArgsConstructor
public class ScheduledTasksWithCustomScheduler {
    
    private final TaskScheduler longRunningScheduler;
    
    @PostConstruct
    public void scheduleTasksWithCustomScheduler() {
        longRunningScheduler.scheduleAtFixedRate(
            this::longRunningReport,
            Duration.ofHours(1)
        );
    }
    
    public void longRunningReport() {
        log.info("Running on {}", Thread.currentThread().getName());
        // Long running task...
    }
}
```

---

### 6.4 Dynamic and Conditional Scheduling

#### **Conditional Scheduling with @ConditionalOnProperty**

```java
@Service
@Slf4j
@ConditionalOnProperty(
    name = "app.scheduler.enabled",
    havingValue = "true",
    matchIfMissing = true  // Enable by default
)
public class ConditionalScheduledTasks {
    
    @Scheduled(fixedRate = 60000)
    public void conditionalTask() {
        log.info("This only runs if app.scheduler.enabled=true");
    }
}

// Disable in certain profiles
@Service
@Profile("!test")  // Don't run in test profile
public class ProductionScheduler {
    
    @Scheduled(cron = "0 0 * * * *")
    public void hourlyTask() {
        log.info("Production hourly task");
    }
}
```

#### **Dynamic Scheduling (Programmatic)**

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class DynamicSchedulerService {
    
    private final TaskScheduler taskScheduler;
    private final Map<String, ScheduledFuture<?>> scheduledTasks = new ConcurrentHashMap<>();
    
    /**
     * Schedule a new task dynamically
     */
    public void scheduleTask(String taskId, Runnable task, Duration interval) {
        // Cancel existing task with same ID if present
        cancelTask(taskId);
        
        ScheduledFuture<?> future = taskScheduler.scheduleAtFixedRate(
            task,
            interval
        );
        
        scheduledTasks.put(taskId, future);
        log.info("Scheduled task: {} with interval: {}", taskId, interval);
    }
    
    /**
     * Schedule with cron expression
     */
    public void scheduleTaskWithCron(String taskId, Runnable task, String cronExpression) {
        cancelTask(taskId);
        
        ScheduledFuture<?> future = taskScheduler.schedule(
            task,
            new CronTrigger(cronExpression)
        );
        
        scheduledTasks.put(taskId, future);
        log.info("Scheduled task: {} with cron: {}", taskId, cronExpression);
    }
    
    /**
     * Cancel a scheduled task
     */
    public boolean cancelTask(String taskId) {
        ScheduledFuture<?> future = scheduledTasks.remove(taskId);
        if (future != null) {
            boolean cancelled = future.cancel(false);  // false = don't interrupt if running
            log.info("Cancelled task: {}, success: {}", taskId, cancelled);
            return cancelled;
        }
        return false;
    }
    
    /**
     * Update schedule interval
     */
    public void updateTaskInterval(String taskId, Duration newInterval) {
        ScheduledFuture<?> existingFuture = scheduledTasks.get(taskId);
        if (existingFuture != null) {
            // Need to store the Runnable to reschedule
            // This is a simplified example - in production, store task metadata
            log.info("Rescheduling task {} with new interval: {}", taskId, newInterval);
        }
    }
    
    /**
     * Get all scheduled task IDs
     */
    public Set<String> getScheduledTaskIds() {
        return scheduledTasks.keySet();
    }
    
    /**
     * Check if task is still scheduled
     */
    public boolean isTaskScheduled(String taskId) {
        ScheduledFuture<?> future = scheduledTasks.get(taskId);
        return future != null && !future.isCancelled() && !future.isDone();
    }
}

// Usage example
@RestController
@RequestMapping("/api/scheduler")
@RequiredArgsConstructor
public class SchedulerController {
    
    private final DynamicSchedulerService schedulerService;
    private final ReportService reportService;
    
    @PostMapping("/tasks/report")
    public ResponseEntity<String> scheduleReportTask(@RequestParam long intervalMinutes) {
        schedulerService.scheduleTask(
            "daily-report",
            () -> reportService.generateDailyReport(),
            Duration.ofMinutes(intervalMinutes)
        );
        return ResponseEntity.ok("Report task scheduled");
    }
    
    @DeleteMapping("/tasks/{taskId}")
    public ResponseEntity<String> cancelTask(@PathVariable String taskId) {
        boolean cancelled = schedulerService.cancelTask(taskId);
        return cancelled 
            ? ResponseEntity.ok("Task cancelled")
            : ResponseEntity.notFound().build();
    }
}
```

#### **Scheduling Based on Database Configuration**

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class ConfigDrivenScheduler implements ApplicationRunner {
    
    private final TaskScheduler taskScheduler;
    private final ScheduledTaskRepository taskRepository;
    private final ApplicationContext applicationContext;
    
    private final Map<Long, ScheduledFuture<?>> activeTasks = new ConcurrentHashMap<>();
    
    @Override
    public void run(ApplicationArguments args) {
        // Load and schedule tasks from database on startup
        refreshScheduledTasks();
    }
    
    @Scheduled(fixedRate = 60000)  // Check for config changes every minute
    public void refreshScheduledTasks() {
        List<ScheduledTaskConfig> configs = taskRepository.findAllEnabled();
        
        // Cancel tasks no longer in config
        Set<Long> configIds = configs.stream()
            .map(ScheduledTaskConfig::getId)
            .collect(Collectors.toSet());
        
        activeTasks.keySet().stream()
            .filter(id -> !configIds.contains(id))
            .forEach(this::cancelTask);
        
        // Schedule new or updated tasks
        for (ScheduledTaskConfig config : configs) {
            scheduleFromConfig(config);
        }
    }
    
    private void scheduleFromConfig(ScheduledTaskConfig config) {
        ScheduledFuture<?> existing = activeTasks.get(config.getId());
        
        // Check if already scheduled with same config
        if (existing != null && !configChanged(config)) {
            return;
        }
        
        // Cancel existing if present
        cancelTask(config.getId());
        
        // Get the task bean
        Runnable task = getTaskRunnable(config.getTaskBeanName(), config.getMethodName());
        
        ScheduledFuture<?> future;
        if (config.getCronExpression() != null) {
            future = taskScheduler.schedule(task, new CronTrigger(config.getCronExpression()));
        } else {
            future = taskScheduler.scheduleAtFixedRate(task, Duration.ofMillis(config.getIntervalMs()));
        }
        
        activeTasks.put(config.getId(), future);
        log.info("Scheduled task: {} ({})", config.getName(), config.getId());
    }
    
    private void cancelTask(Long taskId) {
        ScheduledFuture<?> future = activeTasks.remove(taskId);
        if (future != null) {
            future.cancel(false);
            log.info("Cancelled task: {}", taskId);
        }
    }
    
    private Runnable getTaskRunnable(String beanName, String methodName) {
        Object bean = applicationContext.getBean(beanName);
        Method method = ReflectionUtils.findMethod(bean.getClass(), methodName);
        return () -> ReflectionUtils.invokeMethod(method, bean);
    }
}

// Entity for task configuration
@Entity
@Data
public class ScheduledTaskConfig {
    @Id
    @GeneratedValue
    private Long id;
    
    private String name;
    private String taskBeanName;
    private String methodName;
    private String cronExpression;
    private Long intervalMs;
    private boolean enabled;
    private LocalDateTime lastModified;
}
```

---

### 6.5 Best Practices and Production Considerations

#### **Best Practices Checklist**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SCHEDULING BEST PRACTICES                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ☐ THREAD POOL CONFIGURATION                                                │
│    ├─ Always configure multi-threaded scheduler in production               │
│    ├─ Size based on number of concurrent scheduled tasks                    │
│    └─ Consider separate pools for long-running vs quick tasks               │
│                                                                             │
│  ☐ ERROR HANDLING                                                           │
│    ├─ Wrap task body in try-catch                                          │
│    ├─ Log exceptions with context                                          │
│    ├─ Send alerts for critical task failures                               │
│    └─ Consider retry logic for transient failures                          │
│                                                                             │
│  ☐ DISTRIBUTED SYSTEMS                                                      │
│    ├─ Use ShedLock or similar for cluster-safe scheduling                  │
│    ├─ One instance should run the task, not all                            │
│    └─ Consider using external scheduler (Quartz, Kubernetes CronJob)       │
│                                                                             │
│  ☐ MONITORING                                                               │
│    ├─ Log task start/end with duration                                     │
│    ├─ Record metrics for task execution                                    │
│    ├─ Alert on tasks that don't run or take too long                       │
│    └─ Health check endpoint for scheduler status                           │
│                                                                             │
│  ☐ CONFIGURATION                                                            │
│    ├─ Externalize cron expressions and intervals                           │
│    ├─ Allow disabling via property                                         │
│    └─ Different schedules for different environments                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Production-Ready Scheduled Task Template**

```java
@Service
@Slf4j
@RequiredArgsConstructor
public class ProductionScheduledTasks {
    
    private final MeterRegistry meterRegistry;
    private final AlertService alertService;
    
    /**
     * Production-ready scheduled task with:
     * - Proper error handling
     * - Metrics recording
     * - Alerting
     * - Logging
     */
    @Scheduled(cron = "${app.tasks.data-sync.cron:0 0 * * * *}")
    public void dataSyncTask() {
        String taskName = "data-sync";
        long startTime = System.currentTimeMillis();
        log.info("[{}] Starting task", taskName);
        
        try {
            // Actual task logic wrapped in try-catch
            performDataSync();
            
            // Record success metrics
            long duration = System.currentTimeMillis() - startTime;
            recordSuccess(taskName, duration);
            log.info("[{}] Completed successfully in {}ms", taskName, duration);
            
        } catch (Exception e) {
            // Record failure
            long duration = System.currentTimeMillis() - startTime;
            recordFailure(taskName, duration, e);
            log.error("[{}] Failed after {}ms", taskName, duration, e);
            
            // Alert for critical tasks
            alertService.sendAlert(
                AlertLevel.HIGH,
                "Scheduled task failed: " + taskName,
                e.getMessage()
            );
            
            // Re-throw if you want to see in health checks
            // Or swallow if task will retry on next schedule
        }
    }
    
    private void recordSuccess(String taskName, long duration) {
        meterRegistry.counter("scheduled.task.success", "task", taskName).increment();
        meterRegistry.timer("scheduled.task.duration", "task", taskName)
            .record(duration, TimeUnit.MILLISECONDS);
    }
    
    private void recordFailure(String taskName, long duration, Exception e) {
        meterRegistry.counter("scheduled.task.failure", 
            "task", taskName, 
            "exception", e.getClass().getSimpleName()
        ).increment();
        meterRegistry.timer("scheduled.task.duration", "task", taskName)
            .record(duration, TimeUnit.MILLISECONDS);
    }
    
    private void performDataSync() {
        // Task implementation
    }
}
```

#### **Distributed Scheduling with ShedLock**

```java
// Add dependency: net.javacrumbs.shedlock:shedlock-spring

@Configuration
@EnableScheduling
@EnableSchedulerLock(defaultLockAtMostFor = "PT30M")  // 30 minutes max
public class ShedLockConfig {
    
    @Bean
    public LockProvider lockProvider(DataSource dataSource) {
        return new JdbcTemplateLockProvider(
            JdbcTemplateLockProvider.Configuration.builder()
                .withJdbcTemplate(new JdbcTemplate(dataSource))
                .usingDbTime()  // Use DB time for consistency
                .build()
        );
    }
}

@Service
@Slf4j
public class DistributedScheduledTasks {
    
    /**
     * With ShedLock: Only ONE instance in cluster executes this task
     * Other instances skip if lock is held
     */
    @Scheduled(cron = "0 0 * * * *")  // Every hour
    @SchedulerLock(
        name = "hourlyReportTask",        // Lock name (unique identifier)
        lockAtLeastFor = "PT5M",          // Hold lock for minimum 5 minutes
        lockAtMostFor = "PT30M"           // Release lock after max 30 minutes
    )
    public void hourlyReport() {
        log.info("Running hourly report - only on one instance");
        reportService.generateHourlyReport();
    }
    
    /**
     * Critical task - must not run concurrently
     */
    @Scheduled(fixedRate = 60000)
    @SchedulerLock(
        name = "paymentReconciliation",
        lockAtLeastFor = "PT30S",   // Prevent rapid re-execution
        lockAtMostFor = "PT15M"     // Safety timeout
    )
    public void paymentReconciliation() {
        log.info("Running payment reconciliation");
        paymentService.reconcile();
    }
}

/*
 * SQL for ShedLock table:
 * 
 * CREATE TABLE shedlock (
 *   name VARCHAR(64) NOT NULL PRIMARY KEY,
 *   lock_until TIMESTAMP(3) NOT NULL,
 *   locked_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
 *   locked_by VARCHAR(255) NOT NULL
 * );
 */
```

#### **Scheduler Health Check**

```java
@Component
@RequiredArgsConstructor
public class SchedulerHealthIndicator implements HealthIndicator {
    
    private final ThreadPoolTaskScheduler taskScheduler;
    
    @Override
    public Health health() {
        try {
            ThreadPoolExecutor executor = taskScheduler.getScheduledThreadPoolExecutor();
            
            int activeCount = executor.getActiveCount();
            int poolSize = executor.getPoolSize();
            int corePoolSize = executor.getCorePoolSize();
            long completedTasks = executor.getCompletedTaskCount();
            long taskCount = executor.getTaskCount();
            int queueSize = executor.getQueue().size();
            
            // Check for potential issues
            boolean isHealthy = true;
            String reason = "";
            
            // All threads busy and queue growing
            if (activeCount == poolSize && queueSize > poolSize * 2) {
                isHealthy = false;
                reason = "Scheduler may be overloaded";
            }
            
            Health.Builder builder = isHealthy ? Health.up() : Health.down();
            
            return builder
                .withDetail("activeThreads", activeCount)
                .withDetail("poolSize", poolSize)
                .withDetail("corePoolSize", corePoolSize)
                .withDetail("completedTasks", completedTasks)
                .withDetail("totalTasks", taskCount)
                .withDetail("queueSize", queueSize)
                .withDetail("reason", reason)
                .build();
                
        } catch (Exception e) {
            return Health.down()
                .withDetail("error", e.getMessage())
                .build();
        }
    }
}
```

#### **Quick Reference: @Scheduled Attributes**

| Attribute | Type | Description | Example |
|-----------|------|-------------|---------|
| `fixedRate` | long | Execute every N ms | `@Scheduled(fixedRate = 5000)` |
| `fixedDelay` | long | Wait N ms after completion | `@Scheduled(fixedDelay = 5000)` |
| `initialDelay` | long | Wait before first execution | `@Scheduled(fixedRate = 5000, initialDelay = 10000)` |
| `cron` | String | Cron expression | `@Scheduled(cron = "0 0 * * * *")` |
| `zone` | String | Timezone for cron | `@Scheduled(cron = "...", zone = "UTC")` |
| `fixedRateString` | String | Rate from property | `@Scheduled(fixedRateString = "${rate}")` |
| `fixedDelayString` | String | Delay from property | `@Scheduled(fixedDelayString = "${delay}")` |
| `timeUnit` | TimeUnit | Unit for rate/delay | `@Scheduled(fixedRate = 1, timeUnit = MINUTES)` |

---

## Summary - Sections 5 & 6

| Topic | Key Points |
|-------|------------|
| **CompletableFuture** | Non-blocking async composition; always use custom executor |
| **Combining Futures** | thenCombine (2 results), allOf (wait all), anyOf (first wins), thenCompose (chain) |
| **Exception Handling** | exceptionally (recover), handle (transform), whenComplete (side effects) |
| **Performance** | Parallel calls can be 3x+ faster; critical for aggregation endpoints |
| **@EnableScheduling** | Enables task scheduling; default is single-threaded (configure pool!) |
| **@Scheduled** | fixedRate (every N ms), fixedDelay (N ms after completion), cron (calendar-based) |
| **Scheduler Thread Pool** | Always configure for production; use ThreadPoolTaskScheduler |
| **Distributed Scheduling** | Use ShedLock or similar to prevent duplicate execution in clusters |

---

## 7. Spring WebFlux vs Spring MVC

Spring WebFlux is the reactive web framework introduced in Spring 5, offering a non-blocking, event-driven programming model as an alternative to Spring MVC's traditional blocking approach.

### 7.1 Thread Model Comparison

#### **Architecture Overview**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SPRING MVC vs SPRING WEBFLUX                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SPRING MVC (Blocking / Thread-Per-Request)                                 │
│  ──────────────────────────────────────────────                             │
│                                                                             │
│  ┌─────────────┐     ┌──────────────────────────────────────────────┐      │
│  │  Requests   │     │           Tomcat Thread Pool (200)           │      │
│  │             │     │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ │      │
│  │  Req 1 ────────────►│Thread-1│ │Thread-2│ │Thread-3│ │  ...   │ │      │
│  │  Req 2 ────────────►│ [Busy] │ │ [Busy] │ │ [Busy] │ │        │ │      │
│  │  Req 3 ────────────►│        │ │        │ │        │ │        │ │      │
│  │  Req 201 ──── X     │Waiting │ │Waiting │ │Waiting │ │        │ │      │
│  │  (Blocked!)   │     │for I/O │ │for I/O │ │for I/O │ │        │ │      │
│  └─────────────┘     │  └────────┘ └────────┘ └────────┘ └────────┘ │      │
│                       └──────────────────────────────────────────────┘      │
│                                                                             │
│  Each request = 1 dedicated thread (blocked during I/O)                     │
│  Max concurrency = Thread pool size (typically 200)                         │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  SPRING WEBFLUX (Non-Blocking / Event Loop)                                 │
│  ──────────────────────────────────────────────                             │
│                                                                             │
│  ┌─────────────┐     ┌──────────────────────────────────────────────┐      │
│  │  Requests   │     │      Event Loop (Few threads = CPU cores)    │      │
│  │             │     │  ┌────────────────────────────────────────┐  │      │
│  │  Req 1 ────────────►│         Single Thread handles ALL       │  │      │
│  │  Req 2 ────────────►│                                          │  │      │
│  │  Req 3 ────────────►│   ┌─────┐   ┌─────┐   ┌─────┐          │  │      │
│  │    ...      │     │   │Req 1│   │Req 2│   │Req 3│  ...      │  │      │
│  │  Req 10000 ─────────►│   └──┬──┘   └──┬──┘   └──┬──┘          │  │      │
│  │  (OK!)      │     │      │         │         │              │  │      │
│  └─────────────┘     │      ▼         ▼         ▼              │  │      │
│                       │   [Non-blocking I/O callbacks]          │  │      │
│                       └──────────────────────────────────────────────┘      │
│                                                                             │
│  Single thread handles many requests (never blocks)                         │
│  Max concurrency = Limited only by memory                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Resource Comparison**

| Aspect | Spring MVC | Spring WebFlux |
|--------|------------|----------------|
| **Threading Model** | Thread-per-request | Event loop |
| **Threads for 10K concurrent** | 10,000+ threads | ~4-16 threads (CPU cores) |
| **Memory per connection** | ~1MB (thread stack) | ~KB (lightweight) |
| **I/O Model** | Blocking (JDBC, RestTemplate) | Non-blocking (R2DBC, WebClient) |
| **CPU Usage** | Higher (context switching) | Lower (minimal switching) |
| **Latency under load** | Increases (thread exhaustion) | Stable (no blocking) |
| **Code Style** | Imperative | Reactive (functional) |
| **Learning Curve** | Lower | Higher |

#### **Thread Pool Visual Comparison**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THREAD UTILIZATION: 10,000 REQUESTS                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SPRING MVC with I/O-bound requests (200ms DB call each):                   │
│  ────────────────────────────────────────────────────────                   │
│                                                                             │
│  Thread-1:    [████████ Waiting for DB ████████]                            │
│  Thread-2:    [████████ Waiting for DB ████████]                            │
│  Thread-3:    [████████ Waiting for DB ████████]                            │
│    ...        [████████ Waiting for DB ████████]                            │
│  Thread-200:  [████████ Waiting for DB ████████]                            │
│                                                                             │
│  Req 201-10000: ⏳ Waiting in queue (thread pool exhausted!)                │
│                                                                             │
│  Time to complete 10,000 requests: ~10 seconds (batches of 200)             │
│  Memory: 200+ threads × ~1MB = 200+ MB just for threads                     │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  SPRING WEBFLUX with same requests:                                         │
│  ──────────────────────────────────                                         │
│                                                                             │
│  Thread-1: [R1][R2][R3][R4]...[R10000] (processes events, never waits)      │
│  Thread-2: [R5][R6][R7][R8]...(handles completions)                         │
│  Thread-3: (available)                                                      │
│  Thread-4: (available)                                                      │
│                                                                             │
│  All 10,000 initiated immediately → Complete as I/O finishes               │
│                                                                             │
│  Time to complete 10,000 requests: ~200ms (I/O time only!)                  │
│  Memory: 4 threads × ~1MB = ~4 MB for threads                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 7.2 Event Loop Model

#### **How Event Loop Works**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    EVENT LOOP PROCESSING                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌──────────────────┐                                │
│                         │    Event Queue   │                                │
│                         │  ┌────────────┐  │                                │
│                         │  │ Request 1  │  │                                │
│                         │  │ Request 2  │  │                                │
│                         │  │ DB Result  │  │                                │
│                         │  │ HTTP Resp  │  │                                │
│                         │  │ Request 3  │  │                                │
│                         │  └────────────┘  │                                │
│                         └────────┬─────────┘                                │
│                                  │                                          │
│                                  ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐     │
│  │                        EVENT LOOP                                  │     │
│  │                     (Single Thread)                                │     │
│  │                                                                    │     │
│  │   while (true) {                                                   │     │
│  │       Event event = queue.poll();     // Get next event           │     │
│  │       if (event.isNewRequest()) {                                 │     │
│  │           startProcessing(event);      // Non-blocking start      │     │
│  │       } else if (event.isIOComplete()) {                          │     │
│  │           continueProcessing(event);   // Resume with data        │     │
│  │       }                                                            │     │
│  │       // Never blocks - immediately processes next event          │     │
│  │   }                                                                │     │
│  │                                                                    │     │
│  └───────────────────────────────────────────────────────────────────┘     │
│                                  │                                          │
│                                  ▼                                          │
│              ┌───────────────────────────────────────┐                     │
│              │         Non-Blocking I/O              │                     │
│              │  ┌─────────┐  ┌─────────┐  ┌───────┐ │                     │
│              │  │   DB    │  │  HTTP   │  │ File  │ │                     │
│              │  │(R2DBC)  │  │(WebClient) │ (NIO) │ │                     │
│              │  └─────────┘  └─────────┘  └───────┘ │                     │
│              │                                       │                     │
│              │  I/O completes → Event added to queue │                     │
│              └───────────────────────────────────────┘                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Netty Event Loop in WebFlux**

```java
/**
 * WebFlux with Reactor Netty (default)
 * Uses event loop groups for maximum efficiency
 */
@Configuration
public class NettyConfig {
    
    /**
     * Netty creates event loop groups:
     * - Boss group: Accepts connections (1 thread typically)
     * - Worker group: Handles I/O (CPU cores threads)
     */
    @Bean
    public NettyReactiveWebServerFactory nettyFactory() {
        NettyReactiveWebServerFactory factory = new NettyReactiveWebServerFactory();
        factory.addServerCustomizers(httpServer -> httpServer
            .runOn(LoopResources.create(
                "event-loop",
                1,                                    // Select threads (acceptor)
                Runtime.getRuntime().availableProcessors(), // Worker threads
                true                                  // Daemon threads
            ))
        );
        return factory;
    }
}

// Thread naming:
// reactor-http-nio-1 (acceptor)
// reactor-http-nio-2, -3, -4... (workers)
```

#### **Request Lifecycle in Event Loop**

```java
@RestController
@Slf4j
public class WebFluxController {
    
    @Autowired
    private WebClient webClient;
    
    @Autowired
    private ReactiveUserRepository userRepository;  // R2DBC
    
    /**
     * Reactive request processing
     * Thread is NEVER blocked!
     */
    @GetMapping("/users/{id}/profile")
    public Mono<UserProfile> getUserProfile(@PathVariable Long id) {
        log.info("Request received on: {}", Thread.currentThread().getName());
        // Output: reactor-http-nio-1
        
        return userRepository.findById(id)  // Non-blocking DB call
            .doOnNext(user -> log.info("DB returned on: {}", 
                Thread.currentThread().getName()))
            // Output: reactor-http-nio-2 (might be different thread!)
            
            .flatMap(user -> webClient.get()
                .uri("/external/details/{id}", user.getExternalId())
                .retrieve()
                .bodyToMono(ExternalDetails.class)
                .doOnNext(details -> log.info("HTTP returned on: {}", 
                    Thread.currentThread().getName()))
                // Output: reactor-http-nio-3
                
                .map(details -> new UserProfile(user, details))
            );
        // Response goes out on whatever thread completes last
    }
}

/*
 * TIMELINE:
 * 
 * Time 0ms:    Thread-1 receives request, starts DB query, moves to next request
 * Time 50ms:   Thread-2 gets DB result event, starts HTTP call, moves on
 * Time 150ms:  Thread-3 gets HTTP result event, builds response, sends it
 * 
 * Total thread blocking time: 0ms!
 * Threads were always doing useful work.
 */
```

---

### 7.3 Non-Blocking vs Blocking I/O

#### **I/O Comparison**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BLOCKING vs NON-BLOCKING I/O                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BLOCKING I/O (Spring MVC typical)                                          │
│  ─────────────────────────────────                                          │
│                                                                             │
│  Thread-1:                                                                  │
│  ┌──────┐ ┌────────────────────────────────────┐ ┌──────┐                  │
│  │Start │ │   Thread BLOCKED waiting for I/O   │ │Resume│                  │
│  │ I/O  │ │                                    │ │      │                  │
│  └──────┘ └────────────────────────────────────┘ └──────┘                  │
│      │                     ▲                          │                     │
│      │    200ms waiting    │                          │                     │
│      │    (doing nothing)  │                          │                     │
│      └─────────────────────┘                          │                     │
│                                                                             │
│  Problem: Thread sits idle, consuming memory, can't serve other requests   │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  NON-BLOCKING I/O (Spring WebFlux)                                          │
│  ─────────────────────────────────                                          │
│                                                                             │
│  Thread-1:                                                                  │
│  ┌──────┐                                        ┌──────┐                  │
│  │Start │ → Register callback → FREE!            │Resume│                  │
│  │ I/O  │                                        │(callback)               │
│  └──────┘                                        └──────┘                  │
│      │                                                ▲                     │
│      │                                                │                     │
│      │    Thread handles other requests    ┌──────────┘                    │
│      │    during this time!                │                               │
│      │                                      │  I/O complete event          │
│      │                                      │                               │
│      └──────────────────────────────────────┘                               │
│                                                                             │
│  Benefit: Thread immediately available for other work                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Blocking vs Non-Blocking Stacks**

| Layer | Blocking (MVC) | Non-Blocking (WebFlux) |
|-------|----------------|------------------------|
| **Web Server** | Tomcat, Jetty | Netty, Undertow |
| **HTTP Client** | RestTemplate | WebClient |
| **Database** | JDBC, JPA/Hibernate | R2DBC, MongoDB Reactive |
| **Caching** | Jedis, Lettuce (sync) | Lettuce (reactive) |
| **Messaging** | JMS, RabbitMQ (sync) | Reactive RabbitMQ, Kafka Reactor |

#### **Code: Blocking vs Non-Blocking**

```java
// ═══════════════════════════════════════════════════════════════════
// BLOCKING APPROACH (Spring MVC)
// ═══════════════════════════════════════════════════════════════════

@RestController
@RequiredArgsConstructor
public class BlockingController {
    
    private final RestTemplate restTemplate;  // Blocking HTTP client
    private final JdbcTemplate jdbcTemplate;  // Blocking DB access
    
    @GetMapping("/users/{id}/dashboard")
    public Dashboard getDashboard(@PathVariable Long id) {
        // Thread blocked for each call - SEQUENTIAL
        
        // Call 1: Thread waits ~100ms
        User user = jdbcTemplate.queryForObject(
            "SELECT * FROM users WHERE id = ?",
            new BeanPropertyRowMapper<>(User.class),
            id
        );
        
        // Call 2: Thread waits ~200ms  
        Orders orders = restTemplate.getForObject(
            "http://order-service/orders?userId=" + id,
            Orders.class
        );
        
        // Call 3: Thread waits ~150ms
        Notifications notifs = restTemplate.getForObject(
            "http://notification-service/notifications?userId=" + id,
            Notifications.class
        );
        
        // Total: ~450ms of thread blocking
        return new Dashboard(user, orders, notifs);
    }
}

// ═══════════════════════════════════════════════════════════════════
// NON-BLOCKING APPROACH (Spring WebFlux)
// ═══════════════════════════════════════════════════════════════════

@RestController
@RequiredArgsConstructor
public class ReactiveController {
    
    private final WebClient webClient;                    // Non-blocking HTTP
    private final ReactiveCrudRepository<User, Long> userRepository;  // R2DBC
    
    @GetMapping("/users/{id}/dashboard")
    public Mono<Dashboard> getDashboard(@PathVariable Long id) {
        // All calls start IMMEDIATELY - PARALLEL
        // Thread NEVER blocks!
        
        Mono<User> userMono = userRepository.findById(id);  // ~100ms
        
        Mono<Orders> ordersMono = webClient.get()
            .uri("http://order-service/orders?userId={id}", id)
            .retrieve()
            .bodyToMono(Orders.class);  // ~200ms
        
        Mono<Notifications> notifsMono = webClient.get()
            .uri("http://notification-service/notifications?userId={id}", id)
            .retrieve()
            .bodyToMono(Notifications.class);  // ~150ms
        
        // Combine when all complete
        return Mono.zip(userMono, ordersMono, notifsMono)
            .map(tuple -> new Dashboard(
                tuple.getT1(),  // User
                tuple.getT2(),  // Orders
                tuple.getT3()   // Notifications
            ));
        
        // Total time: ~200ms (slowest call) - Thread blocking: 0ms!
    }
}
```

#### **Don't Mix Blocking in Reactive Code!**

```java
@RestController
public class MixedController {
    
    @Autowired
    private JdbcTemplate jdbcTemplate;  // BLOCKING!
    
    // ❌ WRONG: Blocking call in reactive chain
    @GetMapping("/bad")
    public Mono<User> badPattern(@PathVariable Long id) {
        return Mono.fromCallable(() -> {
            // This BLOCKS the event loop thread!
            return jdbcTemplate.queryForObject(
                "SELECT * FROM users WHERE id = ?",
                new BeanPropertyRowMapper<>(User.class),
                id
            );
        }); // Event loop thread blocked - defeats purpose of WebFlux!
    }
    
    // ✅ BETTER: If you MUST use blocking, use bounded elastic scheduler
    @GetMapping("/better")
    public Mono<User> betterPattern(@PathVariable Long id) {
        return Mono.fromCallable(() -> 
            jdbcTemplate.queryForObject(
                "SELECT * FROM users WHERE id = ?",
                new BeanPropertyRowMapper<>(User.class),
                id
            )
        ).subscribeOn(Schedulers.boundedElastic());  // Runs on separate thread pool
        // Event loop thread stays free!
    }
    
    // ✅ BEST: Use reactive database driver
    @Autowired
    private ReactiveUserRepository reactiveRepo;  // R2DBC
    
    @GetMapping("/best")
    public Mono<User> bestPattern(@PathVariable Long id) {
        return reactiveRepo.findById(id);  // Truly non-blocking
    }
}
```

---

### 7.4 When to Use WebFlux

#### **Decision Guide**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    WEBFLUX DECISION GUIDE                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ✅ USE WEBFLUX WHEN:                                                       │
│  ────────────────────                                                       │
│                                                                             │
│  • High concurrency with many I/O-bound operations                          │
│    - 10,000+ concurrent connections                                         │
│    - Microservices gateway/proxy                                            │
│    - Real-time streaming applications                                       │
│                                                                             │
│  • Streaming data scenarios                                                 │
│    - Server-Sent Events (SSE)                                               │
│    - WebSocket heavy applications                                           │
│    - Streaming large files                                                  │
│                                                                             │
│  • You have non-blocking stack available                                    │
│    - MongoDB Reactive / R2DBC for database                                  │
│    - Reactive Redis (Lettuce)                                               │
│    - Reactive messaging (Kafka Reactor)                                     │
│                                                                             │
│  • Microservices with many external calls                                   │
│    - API aggregators                                                        │
│    - BFF (Backend For Frontend)                                             │
│                                                                             │
│  ════════════════════════════════════════════════════════════════════       │
│                                                                             │
│  ❌ STICK WITH MVC WHEN:                                                    │
│  ────────────────────────                                                   │
│                                                                             │
│  • You use blocking dependencies (most common case)                         │
│    - JDBC / JPA / Hibernate                                                 │
│    - Many third-party libraries are blocking                                │
│                                                                             │
│  • Simple CRUD applications with moderate traffic                           │
│    - < 1000 concurrent users                                                │
│    - Traditional business applications                                      │
│                                                                             │
│  • Team unfamiliar with reactive programming                                │
│    - Steep learning curve                                                   │
│    - Debugging is harder                                                    │
│    - Stack traces are complex                                               │
│                                                                             │
│  • CPU-bound operations (WebFlux doesn't help)                              │
│    - Heavy calculations                                                     │
│    - Image processing                                                       │
│    - Complex algorithms                                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Performance Comparison by Scenario**

| Scenario | MVC Performance | WebFlux Performance | Winner |
|----------|-----------------|---------------------|--------|
| **100 concurrent, CPU-bound** | ⭐⭐⭐ | ⭐⭐⭐ | Tie |
| **100 concurrent, I/O-bound** | ⭐⭐⭐ | ⭐⭐⭐ | Tie |
| **1000 concurrent, I/O-bound** | ⭐⭐ | ⭐⭐⭐ | WebFlux |
| **10,000 concurrent, I/O-bound** | ⭐ | ⭐⭐⭐⭐ | WebFlux |
| **Streaming large responses** | ⭐ | ⭐⭐⭐⭐ | WebFlux |
| **Simple CRUD with JDBC** | ⭐⭐⭐ | ⭐⭐ | MVC |
| **Real-time updates (SSE)** | ⭐ | ⭐⭐⭐⭐ | WebFlux |

---

### 7.5 Code Comparison and Migration

#### **Complete Example: REST API**

```java
// ═══════════════════════════════════════════════════════════════════
// SPRING MVC VERSION
// ═══════════════════════════════════════════════════════════════════

// Controller
@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
public class ProductController {
    
    private final ProductService productService;
    
    @GetMapping
    public List<Product> getAllProducts() {
        return productService.findAll();
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<Product> getProduct(@PathVariable Long id) {
        return productService.findById(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public ResponseEntity<Product> createProduct(@RequestBody Product product) {
        Product saved = productService.save(product);
        return ResponseEntity.status(HttpStatus.CREATED).body(saved);
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteProduct(@PathVariable Long id) {
        productService.delete(id);
        return ResponseEntity.noContent().build();
    }
}

// Service
@Service
@RequiredArgsConstructor
public class ProductService {
    
    private final ProductRepository repository;  // JPA Repository
    private final RestTemplate restTemplate;
    
    public List<Product> findAll() {
        return repository.findAll();
    }
    
    public Optional<Product> findById(Long id) {
        return repository.findById(id);
    }
    
    public Product save(Product product) {
        // Enrich with external data
        ProductDetails details = restTemplate.getForObject(
            "http://enrichment-service/details?sku=" + product.getSku(),
            ProductDetails.class
        );
        product.setDetails(details);
        return repository.save(product);
    }
    
    public void delete(Long id) {
        repository.deleteById(id);
    }
}

// ═══════════════════════════════════════════════════════════════════
// SPRING WEBFLUX VERSION
// ═══════════════════════════════════════════════════════════════════

// Controller
@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
public class ReactiveProductController {
    
    private final ReactiveProductService productService;
    
    @GetMapping
    public Flux<Product> getAllProducts() {
        return productService.findAll();
    }
    
    @GetMapping("/{id}")
    public Mono<ResponseEntity<Product>> getProduct(@PathVariable Long id) {
        return productService.findById(id)
            .map(ResponseEntity::ok)
            .defaultIfEmpty(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public Mono<ResponseEntity<Product>> createProduct(@RequestBody Product product) {
        return productService.save(product)
            .map(saved -> ResponseEntity.status(HttpStatus.CREATED).body(saved));
    }
    
    @DeleteMapping("/{id}")
    public Mono<ResponseEntity<Void>> deleteProduct(@PathVariable Long id) {
        return productService.delete(id)
            .then(Mono.just(ResponseEntity.noContent().<Void>build()));
    }
    
    // Streaming endpoint - sends products as they're found
    @GetMapping(value = "/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<Product> streamProducts() {
        return productService.findAll()
            .delayElements(Duration.ofMillis(100));  // Demo: delay between items
    }
}

// Service
@Service
@RequiredArgsConstructor
public class ReactiveProductService {
    
    private final ReactiveProductRepository repository;  // R2DBC Repository
    private final WebClient webClient;
    
    public Flux<Product> findAll() {
        return repository.findAll();
    }
    
    public Mono<Product> findById(Long id) {
        return repository.findById(id);
    }
    
    public Mono<Product> save(Product product) {
        // Enrich with external data - non-blocking!
        return webClient.get()
            .uri("http://enrichment-service/details?sku={sku}", product.getSku())
            .retrieve()
            .bodyToMono(ProductDetails.class)
            .map(details -> {
                product.setDetails(details);
                return product;
            })
            .flatMap(repository::save);
    }
    
    public Mono<Void> delete(Long id) {
        return repository.deleteById(id);
    }
}

// Repository (R2DBC)
public interface ReactiveProductRepository extends ReactiveCrudRepository<Product, Long> {
    
    Flux<Product> findByCategory(String category);
    
    @Query("SELECT * FROM products WHERE price > :minPrice")
    Flux<Product> findByPriceGreaterThan(BigDecimal minPrice);
}
```

#### **Mono vs Flux Quick Reference**

| Type | Description | Analogy | Use Case |
|------|-------------|---------|----------|
| `Mono<T>` | 0 or 1 element | `Optional<T>` or `CompletableFuture<T>` | Single entity, single result |
| `Flux<T>` | 0 to N elements | `Stream<T>` or `List<T>` | Collections, streams |

#### **WebFlux Configuration**

```yaml
# application.yml for WebFlux
spring:
  webflux:
    base-path: /api
  r2dbc:
    url: r2dbc:postgresql://localhost:5432/mydb
    username: user
    password: pass
    pool:
      initial-size: 10
      max-size: 50

server:
  port: 8080
  # Netty configuration
  netty:
    connection-timeout: 30s
    idle-timeout: 60s
```

```java
@Configuration
public class WebFluxConfig {
    
    @Bean
    public WebClient webClient() {
        return WebClient.builder()
            .baseUrl("http://external-service")
            .defaultHeader(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON_VALUE)
            .filter(logRequest())
            .filter(logResponse())
            .build();
    }
    
    private ExchangeFilterFunction logRequest() {
        return ExchangeFilterFunction.ofRequestProcessor(clientRequest -> {
            log.debug("Request: {} {}", clientRequest.method(), clientRequest.url());
            return Mono.just(clientRequest);
        });
    }
    
    private ExchangeFilterFunction logResponse() {
        return ExchangeFilterFunction.ofResponseProcessor(clientResponse -> {
            log.debug("Response status: {}", clientResponse.statusCode());
            return Mono.just(clientResponse);
        });
    }
}
```

---

## 8. Database and Multithreading

Database access in multithreaded applications requires careful handling of connection pools, transaction boundaries, and thread safety.

### 8.1 Connection Pools (HikariCP)

#### **Why Connection Pooling Matters**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE CONNECTION POOLING                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  WITHOUT Pool (Creating connections on demand):                             │
│  ─────────────────────────────────────────────                              │
│                                                                             │
│  Request → Create Connection (~100-500ms) → Query → Close → Response        │
│  Request → Create Connection (~100-500ms) → Query → Close → Response        │
│  Request → Create Connection (~100-500ms) → Query → Close → Response        │
│                                                                             │
│  Problems:                                                                  │
│  • Connection creation is EXPENSIVE (TCP, SSL, auth)                        │
│  • No limit on connections → Can overwhelm database                         │
│  • Can't reuse connections                                                  │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  WITH Pool (Pre-created, reusable connections):                             │
│  ─────────────────────────────────────────────                              │
│                                                                             │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │                    HikariCP Connection Pool                         │    │
│  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ │    │
│  │  │Conn1│ │Conn2│ │Conn3│ │Conn4│ │Conn5│ │Conn6│ │ ... │ │ConnN│ │    │
│  │  │[use]│ │[idle]│ │[use]│ │[idle]│ │[use]│ │[idle]│ │     │ │[idle]│ │    │
│  │  └─────┘ └─────┘ └─────┘ └─────┘ └─────┘ └─────┘ └─────┘ └─────┘ │    │
│  └────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  Request → Borrow Connection (~0.1ms) → Query → Return → Response           │
│                                                                             │
│  Benefits:                                                                  │
│  • Fast connection acquisition (already created)                            │
│  • Controlled number of connections                                         │
│  • Health checks and connection validation                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **HikariCP Configuration**

```yaml
# application.yml - Production HikariCP Configuration
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
    driver-class-name: org.postgresql.Driver
    
    hikari:
      # Pool sizing
      minimum-idle: 10              # Minimum idle connections to maintain
      maximum-pool-size: 50         # Maximum connections in pool
      
      # Timeouts
      connection-timeout: 30000     # Max wait for connection (30 sec)
      idle-timeout: 600000          # Idle connection lifetime (10 min)
      max-lifetime: 1800000         # Max connection lifetime (30 min)
      
      # Validation
      validation-timeout: 5000      # Connection validation timeout
      
      # Performance
      pool-name: HikariPool-Main
      auto-commit: true
      
      # Leak detection (development)
      leak-detection-threshold: 60000  # Warn if connection held > 60 sec
```

```java
@Configuration
public class DataSourceConfig {
    
    /**
     * Programmatic HikariCP configuration
     */
    @Bean
    @ConfigurationProperties(prefix = "spring.datasource.hikari")
    public HikariConfig hikariConfig() {
        return new HikariConfig();
    }
    
    @Bean
    public DataSource dataSource(HikariConfig config) {
        return new HikariDataSource(config);
    }
    
    /**
     * Dynamic pool sizing based on environment
     */
    @Bean
    public HikariDataSource dynamicDataSource(
            @Value("${spring.datasource.url}") String url,
            @Value("${spring.datasource.username}") String username,
            @Value("${spring.datasource.password}") String password) {
        
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl(url);
        config.setUsername(username);
        config.setPassword(password);
        
        // Size based on available CPUs and typical I/O ratio
        int cpuCores = Runtime.getRuntime().availableProcessors();
        int poolSize = cpuCores * 2 + 1;  // Basic formula for I/O-bound
        
        config.setMinimumIdle(Math.max(5, poolSize / 2));
        config.setMaximumPoolSize(Math.min(poolSize * 2, 100));
        
        // Good defaults
        config.setConnectionTimeout(30000);
        config.setIdleTimeout(600000);
        config.setMaxLifetime(1800000);
        config.setValidationTimeout(5000);
        
        return new HikariDataSource(config);
    }
}
```

#### **Pool Sizing Formula**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONNECTION POOL SIZING                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  FORMULA: Pool Size = (core_count * 2) + effective_spindle_count            │
│                                                                             │
│  For SSD (spindle_count = 0): Pool Size ≈ CPU cores * 2 + 1                │
│                                                                             │
│  EXAMPLE: 8 core CPU with SSD                                               │
│  Pool Size = (8 * 2) + 1 = 17 connections                                   │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  CRITICAL CONSIDERATION: Tomcat Thread Pool vs DB Pool                      │
│                                                                             │
│  If Tomcat has 200 threads but DB pool has only 10 connections:             │
│  → 190 threads will WAIT for DB connections!                                │
│  → Creates bottleneck and increases latency                                 │
│                                                                             │
│  RULE OF THUMB:                                                             │
│  DB Pool Size ≥ (Tomcat max threads) × (% of requests hitting DB)           │
│                                                                             │
│  Example: 200 Tomcat threads, 50% requests use DB                           │
│  → DB Pool Size ≥ 200 × 0.5 = 100 connections                              │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  WARNING: More connections ≠ Better performance!                            │
│                                                                             │
│  Database has finite resources. Too many connections:                       │
│  • Increases memory usage on DB server                                      │
│  • More contention for locks                                                │
│  • Context switching overhead                                               │
│                                                                             │
│  Sweet spot usually: 20-50 connections per application instance             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Monitoring Connection Pool**

```java
@Component
@RequiredArgsConstructor
@Slf4j
public class ConnectionPoolMonitor {
    
    private final DataSource dataSource;
    private final MeterRegistry meterRegistry;
    
    @PostConstruct
    public void registerMetrics() {
        if (dataSource instanceof HikariDataSource hikari) {
            HikariPoolMXBean pool = hikari.getHikariPoolMXBean();
            
            Gauge.builder("hikari.connections.active", pool, HikariPoolMXBean::getActiveConnections)
                .tag("pool", hikari.getPoolName())
                .register(meterRegistry);
            
            Gauge.builder("hikari.connections.idle", pool, HikariPoolMXBean::getIdleConnections)
                .tag("pool", hikari.getPoolName())
                .register(meterRegistry);
            
            Gauge.builder("hikari.connections.total", pool, HikariPoolMXBean::getTotalConnections)
                .tag("pool", hikari.getPoolName())
                .register(meterRegistry);
            
            Gauge.builder("hikari.connections.waiting", pool, HikariPoolMXBean::getThreadsAwaitingConnection)
                .tag("pool", hikari.getPoolName())
                .register(meterRegistry);
        }
    }
    
    @Scheduled(fixedRate = 60000)
    public void logPoolStats() {
        if (dataSource instanceof HikariDataSource hikari) {
            HikariPoolMXBean pool = hikari.getHikariPoolMXBean();
            log.info("Connection Pool [{}] - Active: {}, Idle: {}, Total: {}, Waiting: {}",
                hikari.getPoolName(),
                pool.getActiveConnections(),
                pool.getIdleConnections(),
                pool.getTotalConnections(),
                pool.getThreadsAwaitingConnection());
            
            // Alert if pool is under pressure
            if (pool.getThreadsAwaitingConnection() > 10) {
                log.warn("High connection wait count - consider increasing pool size");
            }
        }
    }
}
```

---

### 8.2 Thread Safety in JPA/Hibernate

#### **EntityManager Thread Safety**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ENTITYMANAGER THREAD SAFETY                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ⚠️  ENTITYMANAGER IS NOT THREAD-SAFE!                                      │
│                                                                             │
│  EntityManager maintains:                                                   │
│  • First-level cache (Persistence Context)                                  │
│  • Tracks entity states (NEW, MANAGED, DETACHED, REMOVED)                   │
│  • Database transaction                                                     │
│                                                                             │
│  If shared between threads:                                                 │
│  • Cache corruption                                                         │
│  • Wrong data returned                                                      │
│  • Transaction issues                                                       │
│  • ConcurrentModificationException                                          │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  Spring's Solution: Transaction-scoped EntityManager                        │
│                                                                             │
│  Thread-1 ──► Transaction ──► EntityManager-1 (scoped to this thread)       │
│  Thread-2 ──► Transaction ──► EntityManager-2 (separate instance)           │
│  Thread-3 ──► Transaction ──► EntityManager-3 (separate instance)           │
│                                                                             │
│  Spring injects a PROXY that delegates to thread-specific EntityManager     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

```java
@Repository
@RequiredArgsConstructor
public class ProductRepository {
    
    // Spring injects a PROXY, not the actual EntityManager
    // The proxy routes to the correct instance for current thread/transaction
    @PersistenceContext
    private EntityManager entityManager;  // ✅ SAFE - Spring manages this
    
    // ❌ WRONG: Don't store EntityManager in instance variable from method
    private EntityManager cachedEntityManager;  // DANGEROUS!
    
    public void badPattern() {
        // ❌ Don't cache the actual EntityManager
        this.cachedEntityManager = entityManager;  // Could cause thread issues
    }
    
    // ✅ CORRECT: Use @PersistenceContext injected EntityManager
    public Product findById(Long id) {
        return entityManager.find(Product.class, id);
    }
    
    public List<Product> findByCategory(String category) {
        return entityManager
            .createQuery("SELECT p FROM Product p WHERE p.category = :category", Product.class)
            .setParameter("category", category)
            .getResultList();
    }
}
```

#### **Entity and Collection Thread Safety**

```java
@Entity
public class Order {
    
    @Id
    @GeneratedValue
    private Long id;
    
    // ❌ PROBLEM: Lazy collections accessed from multiple threads
    @OneToMany(mappedBy = "order", fetch = FetchType.LAZY)
    private List<OrderItem> items;  // Not thread-safe if accessed outside transaction
    
    // Potential issues:
    // - LazyInitializationException if accessed after session closes
    // - ConcurrentModificationException if multiple threads modify
}

@Service
@RequiredArgsConstructor
public class OrderService {
    
    private final OrderRepository orderRepository;
    
    // ❌ WRONG: Returning entities with lazy collections
    @Transactional(readOnly = true)
    public Order getOrder(Long id) {
        Order order = orderRepository.findById(id).orElseThrow();
        return order;  // Items collection is lazy - will fail if accessed later
    }
    
    // ✅ CORRECT: Initialize collections before returning
    @Transactional(readOnly = true)
    public Order getOrderWithItems(Long id) {
        Order order = orderRepository.findById(id).orElseThrow();
        Hibernate.initialize(order.getItems());  // Force initialization
        return order;
    }
    
    // ✅ BETTER: Use DTO to avoid lazy loading issues
    @Transactional(readOnly = true)
    public OrderDTO getOrderDTO(Long id) {
        Order order = orderRepository.findById(id).orElseThrow();
        return OrderDTO.builder()
            .id(order.getId())
            .status(order.getStatus())
            .items(order.getItems().stream()
                .map(item -> new OrderItemDTO(item.getId(), item.getProductName()))
                .collect(Collectors.toList()))
            .build();
    }
    
    // ✅ BEST: Use fetch join in query
    @Transactional(readOnly = true)
    public Order getOrderWithItemsJoin(Long id) {
        return orderRepository.findByIdWithItems(id)  // Uses JOIN FETCH
            .orElseThrow();
    }
}

// Repository with fetch join
public interface OrderRepository extends JpaRepository<Order, Long> {
    
    @Query("SELECT o FROM Order o LEFT JOIN FETCH o.items WHERE o.id = :id")
    Optional<Order> findByIdWithItems(@Param("id") Long id);
}
```

#### **Async and JPA**

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class AsyncOrderService {
    
    private final OrderRepository orderRepository;
    
    // ❌ WRONG: Transaction doesn't propagate to async method
    @Transactional
    public void processOrderBad(Long orderId) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        
        asyncProcessItems(order);  // Transaction NOT available here!
    }
    
    @Async
    public void asyncProcessItems(Order order) {
        // order is DETACHED here! 
        // Items collection will throw LazyInitializationException
        for (OrderItem item : order.getItems()) {  // ❌ FAILS!
            processItem(item);
        }
    }
    
    // ✅ CORRECT: Start new transaction in async method
    @Transactional
    public void processOrderGood(Long orderId) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        
        // Pass ID, not entity
        asyncProcessItems(orderId);
    }
    
    @Async
    @Transactional(propagation = Propagation.REQUIRES_NEW)  // New transaction!
    public void asyncProcessItems(Long orderId) {
        // Fresh lookup with new transaction
        Order order = orderRepository.findByIdWithItems(orderId).orElseThrow();
        
        for (OrderItem item : order.getItems()) {  // ✅ Works!
            processItem(item);
        }
    }
    
    // ✅ ALSO CORRECT: Initialize before going async
    @Transactional
    public void processOrderAlternative(Long orderId) {
        Order order = orderRepository.findByIdWithItems(orderId).orElseThrow();
        
        // Extract data before leaving transaction
        List<Long> itemIds = order.getItems().stream()
            .map(OrderItem::getId)
            .collect(Collectors.toList());
        
        asyncProcessItemIds(itemIds);  // Pass primitive data, not entities
    }
    
    @Async
    public void asyncProcessItemIds(List<Long> itemIds) {
        // Process IDs - no lazy loading issues
        itemIds.forEach(this::processItemById);
    }
}
```

---

### 8.3 Transaction Boundaries

#### **Transaction Flow**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TRANSACTION BOUNDARIES                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  @Transactional Method Execution:                                           │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ Method Call                                                          │   │
│  │      │                                                               │   │
│  │      ▼                                                               │   │
│  │ ┌──────────────────────────────────┐                                │   │
│  │ │ Spring AOP Proxy                  │                                │   │
│  │ │ 1. Get DB connection from pool    │                                │   │
│  │ │ 2. Begin transaction              │                                │   │
│  │ │ 3. Bind connection to thread      │                                │   │
│  │ └─────────────┬────────────────────┘                                │   │
│  │               │                                                      │   │
│  │               ▼                                                      │   │
│  │ ┌──────────────────────────────────┐                                │   │
│  │ │ Actual Method Execution           │ ◄── Same thread, same TX      │   │
│  │ │ • All DB operations use same      │                                │   │
│  │ │   connection                      │                                │   │
│  │ │ • EntityManager bound to TX       │                                │   │
│  │ └─────────────┬────────────────────┘                                │   │
│  │               │                                                      │   │
│  │      ┌────────┴────────┐                                            │   │
│  │      │                 │                                            │   │
│  │      ▼                 ▼                                            │   │
│  │  [Success]         [Exception]                                      │   │
│  │      │                 │                                            │   │
│  │      ▼                 ▼                                            │   │
│  │ ┌─────────┐      ┌───────────┐                                      │   │
│  │ │ COMMIT  │      │ ROLLBACK  │                                      │   │
│  │ └─────────┘      └───────────┘                                      │   │
│  │                                                                      │   │
│  │ Return connection to pool                                           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### **Transaction Propagation with Async**

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class TransactionBoundaryService {
    
    private final OrderRepository orderRepository;
    private final PaymentService paymentService;
    private final NotificationService notificationService;
    
    // ═══════════════════════════════════════════════════════════════════
    // Scenario 1: Transaction DOES NOT propagate to @Async
    // ═══════════════════════════════════════════════════════════════════
    
    @Transactional
    public void createOrderWithAsyncNotification(OrderRequest request) {
        // This runs in Transaction-A on Thread-1
        Order order = new Order(request);
        orderRepository.save(order);  // Uses TX-A
        
        // Async call runs on different thread!
        notificationService.sendAsync(order.getId());  // NO TRANSACTION!
        
        // If this fails, order is rolled back but notification was already sent!
        paymentService.charge(order);  // Uses TX-A
    }
    
    // ═══════════════════════════════════════════════════════════════════
    // Scenario 2: Proper handling with CompletableFuture
    // ═══════════════════════════════════════════════════════════════════
    
    @Transactional
    public Order createOrderProperAsync(OrderRequest request) {
        // All DB operations in same transaction
        Order order = new Order(request);
        orderRepository.save(order);
        paymentService.charge(order);
        
        // Save the ID before transaction ends
        Long orderId = order.getId();
        
        // Schedule notification AFTER transaction commits
        TransactionSynchronizationManager.registerSynchronization(
            new TransactionSynchronization() {
                @Override
                public void afterCommit() {
                    // This runs AFTER transaction commits successfully
                    notificationService.sendAsync(orderId);
                }
            }
        );
        
        return order;
    }
    
    // ═══════════════════════════════════════════════════════════════════
    // Scenario 3: Using @TransactionalEventListener
    // ═══════════════════════════════════════════════════════════════════
    
    @Transactional
    public Order createOrderWithEvent(OrderRequest request) {
        Order order = new Order(request);
        orderRepository.save(order);
        paymentService.charge(order);
        
        // Publish event - listeners will be notified after commit
        applicationEventPublisher.publishEvent(new OrderCreatedEvent(order.getId()));
        
        return order;
    }
}

// Event listener that runs AFTER transaction commits
@Component
@RequiredArgsConstructor
public class OrderEventHandler {
    
    private final NotificationService notificationService;
    
    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    @Async  // Run asynchronously after commit
    public void handleOrderCreated(OrderCreatedEvent event) {
        // Safe to send notification - order is definitely committed
        notificationService.send(event.getOrderId());
    }
    
    @TransactionalEventListener(phase = TransactionPhase.AFTER_ROLLBACK)
    public void handleOrderFailed(OrderCreatedEvent event) {
        // Cleanup if transaction failed
        log.warn("Order creation failed for: {}", event.getOrderId());
    }
}
```

#### **Thread-Local Transaction Binding**

```java
@Service
@Slf4j
public class TransactionDebugService {
    
    @Transactional
    public void demonstrateTransactionBinding() {
        // Show current transaction state
        boolean isActive = TransactionSynchronizationManager.isActualTransactionActive();
        String txName = TransactionSynchronizationManager.getCurrentTransactionName();
        boolean readOnly = TransactionSynchronizationManager.isCurrentTransactionReadOnly();
        
        log.info("Transaction active: {}, name: {}, readOnly: {}", 
            isActive, txName, readOnly);
        // Output: Transaction active: true, name: com.example.Service.method, readOnly: false
        
        // The connection is bound to THIS thread
        // Any DB operation on this thread uses the same connection
    }
    
    @Transactional
    public void parallelDbOperations() {
        // Start parallel operations
        CompletableFuture<User> userFuture = CompletableFuture.supplyAsync(() -> {
            // ❌ This runs on ForkJoinPool thread - NO TRANSACTION!
            log.info("Running on: {}, TX active: {}",
                Thread.currentThread().getName(),
                TransactionSynchronizationManager.isActualTransactionActive());
            // Output: Running on: ForkJoinPool.commonPool-worker-1, TX active: false
            
            return userRepository.findById(1L).orElse(null);
            // This works but is NOT part of the transaction!
        });
        
        // Wait for result
        User user = userFuture.join();
    }
}
```

---

### 8.4 Common Mistakes and Solutions

#### **Mistake 1: Connection Pool Exhaustion**

```java
// ❌ MISTAKE: Long-running transaction holds connection
@Service
public class BadConnectionUsage {
    
    @Transactional  // Holds DB connection for entire method!
    public void processWithExternalCall(Long orderId) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        
        // External HTTP call - can take 30+ seconds!
        ExternalResult result = restTemplate.getForObject(
            "http://slow-external-service/process", 
            ExternalResult.class
        );
        // DB connection held during entire HTTP call!
        
        order.setStatus(result.getStatus());
        orderRepository.save(order);
    }
}

// ✅ SOLUTION: Minimize transaction scope
@Service
@RequiredArgsConstructor
public class GoodConnectionUsage {
    
    private final OrderRepository orderRepository;
    private final RestTemplate restTemplate;
    
    public void processWithExternalCall(Long orderId) {
        // Transaction 1: Quick read
        Order order = readOrder(orderId);
        
        // No transaction: External call (can take long)
        ExternalResult result = restTemplate.getForObject(
            "http://slow-external-service/process", 
            ExternalResult.class
        );
        // No DB connection held!
        
        // Transaction 2: Quick update
        updateOrderStatus(orderId, result.getStatus());
    }
    
    @Transactional(readOnly = true)
    public Order readOrder(Long orderId) {
        return orderRepository.findById(orderId).orElseThrow();
    }
    
    @Transactional
    public void updateOrderStatus(Long orderId, String status) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        order.setStatus(status);
        // Auto-saved on transaction commit
    }
}
```

#### **Mistake 2: N+1 Query Problem with Async**

```java
// ❌ MISTAKE: Accessing lazy collection in parallel streams
@Service
public class NPlus1Problem {
    
    @Transactional(readOnly = true)
    public List<OrderSummary> getOrderSummaries(List<Long> orderIds) {
        List<Order> orders = orderRepository.findAllById(orderIds);  // 1 query
        
        return orders.parallelStream()  // BAD! Parallel stream with lazy loading
            .map(order -> {
                // Each access to items triggers a query!
                // And it's on different threads - transaction issues!
                int itemCount = order.getItems().size();  // N queries
                return new OrderSummary(order.getId(), itemCount);
            })
            .collect(Collectors.toList());
    }
}

// ✅ SOLUTION: Eager fetch or separate queries
@Service
public class NPlus1Solution {
    
    // Option 1: Fetch join
    @Transactional(readOnly = true)
    public List<OrderSummary> getOrderSummariesJoin(List<Long> orderIds) {
        // Single query with JOIN FETCH
        List<Order> orders = orderRepository.findAllByIdWithItems(orderIds);
        
        return orders.stream()  // Regular stream - all data already loaded
            .map(order -> new OrderSummary(
                order.getId(), 
                order.getItems().size()  // No additional query
            ))
            .collect(Collectors.toList());
    }
    
    // Option 2: Projection query
    @Transactional(readOnly = true)
    public List<OrderSummary> getOrderSummariesProjection(List<Long> orderIds) {
        // Single query returning exactly what we need
        return orderRepository.findOrderSummaries(orderIds);
    }
}

// Repository
public interface OrderRepository extends JpaRepository<Order, Long> {
    
    @Query("SELECT o FROM Order o LEFT JOIN FETCH o.items WHERE o.id IN :ids")
    List<Order> findAllByIdWithItems(@Param("ids") List<Long> ids);
    
    @Query("SELECT new com.example.OrderSummary(o.id, COUNT(i)) " +
           "FROM Order o LEFT JOIN o.items i " +
           "WHERE o.id IN :ids GROUP BY o.id")
    List<OrderSummary> findOrderSummaries(@Param("ids") List<Long> ids);
}
```

#### **Mistake 3: Modifying Detached Entities**

```java
// ❌ MISTAKE: Modifying entity after transaction ends
@Service
public class DetachedEntityProblem {
    
    @Transactional(readOnly = true)
    public Order getOrder(Long id) {
        return orderRepository.findById(id).orElseThrow();
    }  // Transaction ends - entity becomes DETACHED
    
    public void processOrder(Long orderId) {
        Order order = getOrder(orderId);  // Returns detached entity
        
        order.setStatus("PROCESSED");  // Modifying detached entity!
        // This change is NOT persisted!
    }
}

// ✅ SOLUTION: Proper patterns for modifying entities
@Service
@RequiredArgsConstructor
public class DetachedEntitySolution {
    
    private final OrderRepository orderRepository;
    
    // Solution 1: Do everything in one transaction
    @Transactional
    public Order processOrder(Long orderId) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        order.setStatus("PROCESSED");  // Entity is MANAGED - auto-persisted
        return order;
    }
    
    // Solution 2: Merge detached entity
    @Transactional
    public Order updateDetachedOrder(Order detachedOrder) {
        return orderRepository.save(detachedOrder);  // Merge -> Managed
    }
    
    // Solution 3: Load and update
    @Transactional
    public void updateOrderStatus(Long orderId, String newStatus) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        order.setStatus(newStatus);
        // No explicit save needed - dirty checking handles it
    }
    
    // Solution 4: Direct update query (no entity loading)
    @Modifying
    @Query("UPDATE Order o SET o.status = :status WHERE o.id = :id")
    void updateStatusDirectly(@Param("id") Long id, @Param("status") String status);
}
```

#### **Mistake 4: Forgetting @Transactional on @Async Callee**

```java
// ❌ MISTAKE: Async method without its own transaction
@Service
public class AsyncTransactionMistake {
    
    @Transactional
    public void createOrder(Order order) {
        orderRepository.save(order);
        processAsync(order.getId());  // Async - new thread, no TX!
    }
    
    @Async
    // Missing @Transactional!
    public void processAsync(Long orderId) {
        // No transaction here!
        Order order = orderRepository.findById(orderId).orElseThrow();
        order.setProcessed(true);
        orderRepository.save(order);  // Works but not transactional
        
        // If exception here, partial state might be saved
        externalService.process(order);  // Can fail
        order.setExternalId("...");
        orderRepository.save(order);  // Separate save
    }
}

// ✅ SOLUTION: Add transaction to async method
@Service
public class AsyncTransactionSolution {
    
    @Transactional
    public void createOrder(Order order) {
        orderRepository.save(order);
        processAsync(order.getId());
    }
    
    @Async
    @Transactional(propagation = Propagation.REQUIRES_NEW)  // New transaction!
    public void processAsync(Long orderId) {
        Order order = orderRepository.findById(orderId).orElseThrow();
        order.setProcessed(true);
        
        externalService.process(order);
        order.setExternalId("...");
        
        // Both changes committed together or rolled back together
    }
}
```

#### **Common Mistakes Quick Reference**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMMON DATABASE MULTITHREADING MISTAKES                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  MISTAKE                          │ SOLUTION                                │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  Long transaction with HTTP calls │ Split: read TX → HTTP → update TX      │
│                                                                             │
│  TX not propagating to @Async     │ Add @Transactional(REQUIRES_NEW)        │
│                                                                             │
│  LazyInitializationException      │ Use JOIN FETCH or DTO projection        │
│                                                                             │
│  Modifying detached entities      │ Load fresh in new TX or use merge       │
│                                                                             │
│  N+1 with parallel streams        │ Fetch all data upfront with JOIN        │
│                                                                             │
│  Connection pool exhaustion       │ Minimize TX scope, tune pool size       │
│                                                                             │
│  Sharing EntityManager            │ Always use @PersistenceContext          │
│                                                                             │
│  Blocking in reactive code        │ Use R2DBC or Schedulers.boundedElastic  │
│                                                                             │
│  DB pool < Tomcat threads         │ Size DB pool based on concurrent usage  │
│                                                                             │
│  No TX on repository method       │ JpaRepository methods auto-TX           │
│                                   │ But custom queries need @Transactional  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Summary - Sections 7 & 8

| Topic | Key Points |
|-------|------------|
| **MVC vs WebFlux** | Thread-per-request vs Event loop; MVC for blocking I/O, WebFlux for high concurrency |
| **Event Loop** | Single thread handles many requests; never blocks; uses callbacks |
| **Non-Blocking I/O** | WebClient, R2DBC, Reactive MongoDB; thread immediately available |
| **When WebFlux** | High concurrency (10K+), streaming, non-blocking stack available |
| **When MVC** | JDBC/JPA apps, moderate traffic, team experience, blocking libraries |
| **Connection Pool** | Pre-created connections; size = CPU cores * 2 + 1; tune based on load |
| **EntityManager** | NOT thread-safe; Spring provides thread-bound proxies |
| **Transaction Boundaries** | Bound to thread; don't leak across @Async; use @TransactionalEventListener |
| **Common Mistakes** | Long TX, lazy loading in async, detached entities, pool exhaustion |

---

## 9. Concurrency Issues in Spring Boot

Concurrency issues are among the most difficult bugs to detect and fix. Spring Boot's default singleton scope combined with multithreaded request handling creates unique challenges.

### 9.1 Race Conditions in Singleton Beans

#### **What is a Race Condition?**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    RACE CONDITION EXPLAINED                                   │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  Two threads accessing shared mutable state without proper synchronization    │
│                                                                               │
│  Thread-1                    Shared Counter               Thread-2            │
│  ─────────                   ──────────────               ─────────           │
│                                                                               │
│  read counter (= 100)              100                                        │
│                                                          read counter (= 100) │
│  add 1 (= 101)                                           add 1 (= 101)        │
│  write counter                     101                                        │
│                                    101 ◄───────────────  write counter        │
│                                                                               │
│  EXPECTED: 102                                                                │
│  ACTUAL:   101  ◄── Lost update!                                              │
│                                                                               │
│  This is a CHECK-THEN-ACT race condition                                      │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

#### **Race Condition in Singleton Bean**

```java
// ❌ DANGEROUS: Mutable state in singleton bean
@Service
public class CounterService {
    
    // Instance variable in singleton = SHARED across all threads!
    private int requestCount = 0;  // Mutable state
    private Map<String, UserSession> sessions = new HashMap<>();  // Not thread-safe!
    
    public void processRequest() {
        // Race condition! Multiple threads can:
        // 1. Read same value
        // 2. Increment
        // 3. Write (overwrites other thread's increment)
        requestCount++;  // NOT ATOMIC!
    }
    
    public void addSession(String id, UserSession session) {
        // HashMap is NOT thread-safe!
        // Can cause infinite loop, data corruption, or lost entries
        sessions.put(id, session);  // ❌ DANGEROUS!
    }
    
    public int getRequestCount() {
        return requestCount;
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// TIMELINE showing the race:
// ══════════════════════════════════════════════════════════════════════════════
//
// Time   │ Thread-1 (http-nio-8080-exec-1) │ Thread-2 (http-nio-8080-exec-2)
// ───────┼─────────────────────────────────┼─────────────────────────────────
// T1     │ read requestCount = 0           │
// T2     │                                 │ read requestCount = 0
// T3     │ compute 0 + 1 = 1               │
// T4     │                                 │ compute 0 + 1 = 1
// T5     │ write requestCount = 1          │
// T6     │                                 │ write requestCount = 1  ← OVERWRITES!
// ───────┼─────────────────────────────────┼─────────────────────────────────
// Result │ Expected: 2, Actual: 1          │ ONE INCREMENT LOST!
```

#### **Solutions for Race Conditions**

```java
// ✅ SOLUTION 1: Use AtomicInteger for counters
@Service
public class ThreadSafeCounterService {
    
    private final AtomicInteger requestCount = new AtomicInteger(0);
    private final AtomicLong totalProcessingTime = new AtomicLong(0);
    
    public void processRequest() {
        requestCount.incrementAndGet();  // Atomic operation!
    }
    
    public int getRequestCount() {
        return requestCount.get();
    }
}

// ✅ SOLUTION 2: Use ConcurrentHashMap for maps
@Service
public class SessionService {
    
    // ConcurrentHashMap is thread-safe for concurrent access
    private final Map<String, UserSession> sessions = new ConcurrentHashMap<>();
    
    public void addSession(String id, UserSession session) {
        sessions.put(id, session);  // Thread-safe!
    }
    
    public UserSession getSession(String id) {
        return sessions.get(id);  // Thread-safe!
    }
    
    // For compound operations, use compute methods
    public UserSession getOrCreateSession(String id) {
        return sessions.computeIfAbsent(id, key -> new UserSession(key));
        // computeIfAbsent is ATOMIC - no race condition!
    }
}

// ✅ SOLUTION 3: Use synchronized for complex operations
@Service
public class ComplexStateService {
    
    private final Object lock = new Object();
    private List<Transaction> pendingTransactions = new ArrayList<>();
    private BigDecimal totalAmount = BigDecimal.ZERO;
    
    public void addTransaction(Transaction tx) {
        synchronized (lock) {
            // Both operations happen atomically
            pendingTransactions.add(tx);
            totalAmount = totalAmount.add(tx.getAmount());
        }
    }
    
    public TransactionSummary getSummary() {
        synchronized (lock) {
            return new TransactionSummary(
                new ArrayList<>(pendingTransactions),  // Defensive copy
                totalAmount
            );
        }
    }
}

// ✅ SOLUTION 4: Use ReentrantLock for more control
@Service
public class AdvancedLockingService {
    
    private final ReentrantLock lock = new ReentrantLock();
    private final Condition notEmpty = lock.newCondition();
    private final Queue<Task> taskQueue = new LinkedList<>();
    
    public void addTask(Task task) {
        lock.lock();
        try {
            taskQueue.offer(task);
            notEmpty.signal();  // Wake up waiting consumers
        } finally {
            lock.unlock();  // ALWAYS unlock in finally!
        }
    }
    
    public Task takeTask() throws InterruptedException {
        lock.lock();
        try {
            while (taskQueue.isEmpty()) {
                notEmpty.await();  // Wait for task
            }
            return taskQueue.poll();
        } finally {
            lock.unlock();
        }
    }
    
    // Try-lock pattern for avoiding deadlocks
    public boolean tryAddTask(Task task, long timeout, TimeUnit unit) {
        try {
            if (lock.tryLock(timeout, unit)) {
                try {
                    taskQueue.offer(task);
                    notEmpty.signal();
                    return true;
                } finally {
                    lock.unlock();
                }
            }
            return false;  // Couldn't acquire lock in time
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            return false;
        }
    }
}
```

---

### 9.2 Thread Safety of Spring Beans

#### **Spring Bean Scopes Overview**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    SPRING BEAN SCOPES AND THREAD SAFETY                       │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  SINGLETON (default)                                                          │
│  ════════════════════                                                         │
│  • ONE instance for entire application                                        │
│  • ALL requests share the same instance                                       │
│  • Instance variables = SHARED STATE (dangerous!)                             │
│                                                                               │
│  ┌──────────────────────────────────────────────────────────────┐             │
│  │              SINGLETON BEAN: UserService                     │             │
│  │  ┌─────────────────────────────────────────────────────────┐ │             │
│  │  │  Instance created once at startup                       │ │             │
│  │  │  All threads access THIS SAME object                    │ │             │
│  │  └─────────────────────────────────────────────────────────┘ │             │
│  └──────────────────────────────────────────────────────────────┘             │
│         ▲              ▲              ▲              ▲                        │
│         │              │              │              │                        │
│    Thread-1       Thread-2       Thread-3       Thread-N                      │
│    (Request 1)   (Request 2)   (Request 3)   (Request N)                      │
│                                                                               │
│  ─────────────────────────────────────────────────────────────────────────    │
│                                                                               │
│  PROTOTYPE                                                                    │
│  ═════════                                                                    │
│  • NEW instance for each injection point                                      │
│  • Not automatically thread-safe (same thread can use for multiple requests)  │
│  • Spring doesn't manage lifecycle after creation                             │
│                                                                               │
│  REQUEST SCOPE (web apps)                                                     │
│  ═══════════════════════                                                      │
│  • ONE instance per HTTP request                                              │
│  • Different threads get different instances                                  │
│  • Inherently thread-safe for request data                                    │
│                                                                               │
│  SESSION SCOPE (web apps)                                                     │
│  ═══════════════════════                                                      │
│  • ONE instance per HTTP session                                              │
│  • Can be accessed by different requests (threads) in same session            │
│  • Needs synchronization for concurrent AJAX calls                            │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

#### **Making Singleton Beans Thread-Safe**

```java
// ══════════════════════════════════════════════════════════════════════════════
// PATTERN 1: Stateless Service (BEST APPROACH)
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class StatelessOrderService {
    
    // ✅ Dependencies are injected (final, immutable references)
    private final OrderRepository orderRepository;
    private final PaymentGateway paymentGateway;
    
    public StatelessOrderService(OrderRepository orderRepository, 
                                  PaymentGateway paymentGateway) {
        this.orderRepository = orderRepository;
        this.paymentGateway = paymentGateway;
    }
    
    // ✅ All data comes from method parameters (local to thread)
    public Order processOrder(OrderRequest request) {
        // Local variables are thread-safe (each thread has its own stack)
        Order order = new Order(request);  // Local variable
        
        PaymentResult result = paymentGateway.charge(order);  // Stateless call
        
        if (result.isSuccessful()) {
            order.setStatus(OrderStatus.PAID);
            return orderRepository.save(order);  // Returns new managed entity
        }
        
        throw new PaymentFailedException(result.getError());
    }
    
    // ✅ No instance variables that hold state
    // ✅ All methods work only with parameters and return values
}

// ══════════════════════════════════════════════════════════════════════════════
// PATTERN 2: Immutable Values + Atomic References
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class ConfigurationService {
    
    // Immutable configuration - safe to share
    private record AppConfig(
        String apiUrl,
        int timeout,
        List<String> allowedOrigins
    ) {
        AppConfig {
            // Defensive copy for list
            allowedOrigins = List.copyOf(allowedOrigins);
        }
    }
    
    // AtomicReference for safe updates
    private final AtomicReference<AppConfig> currentConfig = 
        new AtomicReference<>();
    
    @PostConstruct
    public void loadConfig() {
        currentConfig.set(loadFromSource());
    }
    
    // Thread-safe read
    public AppConfig getConfig() {
        return currentConfig.get();  // Returns immutable object
    }
    
    // Thread-safe update (atomic swap)
    public void updateConfig(AppConfig newConfig) {
        currentConfig.set(newConfig);  // Atomic
    }
    
    // Compare-and-swap for conditional updates
    public boolean updateIfCurrent(AppConfig expected, AppConfig newConfig) {
        return currentConfig.compareAndSet(expected, newConfig);
    }
    
    private AppConfig loadFromSource() {
        // Load from DB/file/remote
        return new AppConfig("https://api.example.com", 30, List.of("localhost"));
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// PATTERN 3: Thread-Safe Collections
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class CacheService {
    
    // Thread-safe map implementations
    private final ConcurrentHashMap<String, CachedItem> cache = 
        new ConcurrentHashMap<>();
    
    // Thread-safe queue for eviction
    private final ConcurrentLinkedQueue<String> evictionQueue = 
        new ConcurrentLinkedQueue<>();
    
    // Copy-on-write for read-heavy, write-rare scenarios
    private final CopyOnWriteArrayList<CacheEventListener> listeners = 
        new CopyOnWriteArrayList<>();
    
    public CachedItem get(String key) {
        return cache.get(key);
    }
    
    public void put(String key, CachedItem item) {
        cache.put(key, item);
        evictionQueue.offer(key);
        
        // Safe iteration (snapshot)
        listeners.forEach(l -> l.onCacheUpdate(key, item));
    }
    
    // Atomic compute operations
    public CachedItem getOrCompute(String key, Function<String, CachedItem> loader) {
        return cache.computeIfAbsent(key, loader);
    }
    
    // Atomic update
    public CachedItem updateExpiry(String key, Duration newExpiry) {
        return cache.computeIfPresent(key, (k, existing) -> 
            existing.withExpiry(newExpiry));
    }
}
```

#### **Thread-Safety Decision Guide**

| Data Type | Thread-Safe Alternative | Use Case |
|-----------|------------------------|----------|
| `int`, `long` | `AtomicInteger`, `AtomicLong` | Counters, metrics |
| `boolean` | `AtomicBoolean` | Flags, toggles |
| `Object reference` | `AtomicReference<T>` | Immutable config swaps |
| `HashMap` | `ConcurrentHashMap` | Caches, lookups |
| `ArrayList` | `CopyOnWriteArrayList` | Listeners (read-heavy) |
| `HashSet` | `ConcurrentHashMap.newKeySet()` | Thread-safe sets |
| `LinkedList` (queue) | `ConcurrentLinkedQueue` | Work queues |
| Multiple related fields | `synchronized` block | Complex state |

---

### 9.3 Using Prototype Scope

#### **When to Use Prototype Scope**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    PROTOTYPE SCOPE USE CASES                                  │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  ✅ USE PROTOTYPE WHEN:                                                       │
│  ─────────────────────────                                                    │
│  • Bean holds request-specific state                                          │
│  • Bean is expensive to make thread-safe                                      │
│  • Different configuration per use case                                       │
│  • Builder pattern beans                                                      │
│                                                                               │
│  ❌ DON'T USE PROTOTYPE FOR:                                                  │
│  ───────────────────────────                                                  │
│  • Stateless services (use singleton)                                         │
│  • Just to "avoid" thread safety (wrong approach)                             │
│  • Heavy initialization (defeats performance)                                 │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

```java
// ══════════════════════════════════════════════════════════════════════════════
// PROTOTYPE BEAN DEFINITION
// ══════════════════════════════════════════════════════════════════════════════

@Component
@Scope("prototype")  // New instance each time
public class ReportBuilder {
    
    // Instance state is SAFE - each caller gets own instance
    private String title;
    private List<ReportSection> sections = new ArrayList<>();
    private ReportFormat format;
    
    public ReportBuilder withTitle(String title) {
        this.title = title;
        return this;
    }
    
    public ReportBuilder addSection(ReportSection section) {
        this.sections.add(section);
        return this;
    }
    
    public ReportBuilder format(ReportFormat format) {
        this.format = format;
        return this;
    }
    
    public Report build() {
        return new Report(title, List.copyOf(sections), format);
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// ⚠️ THE PROTOTYPE-SINGLETON INJECTION PROBLEM
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class WrongReportService {
    
    // ❌ WRONG: Prototype injected into singleton
    // This creates ONE instance at singleton creation time!
    @Autowired
    private ReportBuilder reportBuilder;  // SAME instance for all calls!
    
    public Report generateReport(ReportRequest request) {
        // ❌ BUG: All requests share the same builder!
        return reportBuilder
            .withTitle(request.getTitle())  // Overwrites previous!
            .addSection(request.getData())   // Accumulates!
            .build();
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// ✅ CORRECT WAYS TO USE PROTOTYPE FROM SINGLETON
// ══════════════════════════════════════════════════════════════════════════════

// SOLUTION 1: ObjectProvider (recommended)
@Service
@RequiredArgsConstructor
public class CorrectReportService1 {
    
    // ObjectProvider creates new instance on each get()
    private final ObjectProvider<ReportBuilder> reportBuilderProvider;
    
    public Report generateReport(ReportRequest request) {
        // ✅ Fresh prototype instance for each call
        ReportBuilder builder = reportBuilderProvider.getObject();
        return builder
            .withTitle(request.getTitle())
            .addSection(request.getData())
            .build();
    }
}

// SOLUTION 2: ApplicationContext (works but verbose)
@Service
@RequiredArgsConstructor
public class CorrectReportService2 {
    
    private final ApplicationContext context;
    
    public Report generateReport(ReportRequest request) {
        // ✅ Fresh instance each time
        ReportBuilder builder = context.getBean(ReportBuilder.class);
        return builder
            .withTitle(request.getTitle())
            .addSection(request.getData())
            .build();
    }
}

// SOLUTION 3: @Lookup method injection
@Service
public abstract class CorrectReportService3 {
    
    // Spring overrides this method to return prototype
    @Lookup
    protected abstract ReportBuilder createReportBuilder();
    
    public Report generateReport(ReportRequest request) {
        // ✅ Fresh instance each time
        return createReportBuilder()
            .withTitle(request.getTitle())
            .addSection(request.getData())
            .build();
    }
}

// SOLUTION 4: Scoped proxy (automatic)
@Component
@Scope(value = "prototype", proxyMode = ScopedProxyMode.TARGET_CLASS)
public class ProxiedReportBuilder {
    // Same implementation as ReportBuilder
}

@Service
@RequiredArgsConstructor
public class CorrectReportService4 {
    
    // ✅ Proxy automatically creates new instance per method call
    private final ProxiedReportBuilder reportBuilder;
    
    public Report generateReport(ReportRequest request) {
        // Each method call goes through proxy → new instance
        return reportBuilder
            .withTitle(request.getTitle())
            .addSection(request.getData())
            .build();
    }
}
```

#### **Request Scope Alternative**

```java
// Request scope: One instance per HTTP request
@Component
@Scope(value = WebApplicationContext.SCOPE_REQUEST, 
       proxyMode = ScopedProxyMode.TARGET_CLASS)
public class RequestContext {
    
    private String requestId;
    private Instant startTime;
    private List<String> processingSteps = new ArrayList<>();
    
    @PostConstruct
    public void init() {
        this.requestId = UUID.randomUUID().toString();
        this.startTime = Instant.now();
    }
    
    public void addStep(String step) {
        processingSteps.add(step);
    }
    
    public String getRequestId() {
        return requestId;
    }
    
    public Duration getElapsedTime() {
        return Duration.between(startTime, Instant.now());
    }
}

@Service
@RequiredArgsConstructor
public class OrderService {
    
    // Injected once, but proxy returns request-specific instance
    private final RequestContext requestContext;
    private final OrderRepository orderRepository;
    
    @Transactional
    public Order createOrder(OrderRequest request) {
        requestContext.addStep("Validating order");
        validateOrder(request);
        
        requestContext.addStep("Saving order");
        Order order = orderRepository.save(new Order(request));
        
        log.info("Request {} completed in {} - Steps: {}",
            requestContext.getRequestId(),
            requestContext.getElapsedTime(),
            requestContext.getProcessingSteps());
        
        return order;
    }
}
```

---

### 9.4 Using ThreadLocal Safely

#### **What is ThreadLocal?**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    THREADLOCAL EXPLAINED                                      │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  ThreadLocal provides thread-isolated storage                                 │
│  Each thread sees its own independent value                                   │
│                                                                               │
│  ┌────────────────────────────────────────────────────────────────────────┐   │
│  │                    ThreadLocal<UserContext>                            │   │
│  │                                                                        │   │
│  │   Thread-1               Thread-2               Thread-3               │   │
│  │   ┌─────────────┐        ┌─────────────┐        ┌─────────────┐        │   │
│  │   │ User: Alice │        │ User: Bob   │        │ User: Carol │        │   │
│  │   │ Role: Admin │        │ Role: User  │        │ Role: User  │        │   │
│  │   └─────────────┘        └─────────────┘        └─────────────┘        │   │
│  │        ▲                      ▲                       ▲                │   │
│  │        │                      │                       │                │   │
│  │   threadLocal.get()      threadLocal.get()       threadLocal.get()     │   │
│  │   returns Alice          returns Bob             returns Carol         │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                               │
│  COMMON USE CASES:                                                            │
│  • User context (authentication info)                                         │
│  • Request metadata (correlation ID, tenant ID)                               │
│  • Transaction context                                                        │
│  • Database connection (already used by Spring)                               │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

#### **Basic ThreadLocal Usage**

```java
// Thread-safe user context using ThreadLocal
public class UserContextHolder {
    
    private static final ThreadLocal<UserContext> contextHolder = 
        new ThreadLocal<>();
    
    public static void setContext(UserContext context) {
        contextHolder.set(context);
    }
    
    public static UserContext getContext() {
        return contextHolder.get();
    }
    
    public static void clear() {
        contextHolder.remove();
    }
    
    // InheritableThreadLocal - child threads inherit parent's value
    private static final InheritableThreadLocal<String> correlationId = 
        new InheritableThreadLocal<>();
}

// Usage in filter
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class UserContextFilter implements Filter {
    
    @Override
    public void doFilter(ServletRequest request, ServletResponse response, 
                         FilterChain chain) throws IOException, ServletException {
        HttpServletRequest httpRequest = (HttpServletRequest) request;
        
        try {
            // Extract and set user context
            UserContext context = extractUserContext(httpRequest);
            UserContextHolder.setContext(context);
            
            chain.doFilter(request, response);
            
        } finally {
            // ⚠️ CRITICAL: Always clean up!
            UserContextHolder.clear();
        }
    }
    
    private UserContext extractUserContext(HttpServletRequest request) {
        String userId = request.getHeader("X-User-Id");
        String role = request.getHeader("X-User-Role");
        return new UserContext(userId, role);
    }
}

// Access anywhere in the call chain
@Service
public class AuditService {
    
    public void logAction(String action) {
        UserContext context = UserContextHolder.getContext();
        log.info("User {} performed action: {}", context.getUserId(), action);
    }
}
```

#### **⚠️ ThreadLocal Pitfalls and Solutions**

```java
// ══════════════════════════════════════════════════════════════════════════════
// PITFALL 1: Memory Leak - Not cleaning up ThreadLocal
// ══════════════════════════════════════════════════════════════════════════════

// ❌ WRONG: ThreadLocal not cleaned up
@Service
public class LeakyService {
    
    private static final ThreadLocal<LargeObject> cache = new ThreadLocal<>();
    
    public void process() {
        cache.set(new LargeObject());  // Set value
        doWork();
        // ❌ Never cleaned! Thread returns to pool with value attached
    }
}
// In thread pools (like Tomcat), threads are REUSED!
// LargeObject stays in memory as long as thread exists
// → Memory grows with each new request

// ✅ CORRECT: Always clean up in finally block
@Service
public class SafeService {
    
    private static final ThreadLocal<LargeObject> cache = new ThreadLocal<>();
    
    public void process() {
        try {
            cache.set(new LargeObject());
            doWork();
        } finally {
            cache.remove();  // ✅ Always clean up!
        }
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// PITFALL 2: ThreadLocal not propagating to @Async threads
// ══════════════════════════════════════════════════════════════════════════════

// ❌ WRONG: Context lost in async call
@Service
public class AsyncContextProblem {
    
    public void processRequest() {
        UserContextHolder.setContext(new UserContext("user123", "admin"));
        asyncProcess();  // Context NOT available!
    }
    
    @Async
    public void asyncProcess() {
        // ❌ Returns null! Different thread, different ThreadLocal
        UserContext context = UserContextHolder.getContext();
        log.info("User: {}", context);  // NullPointerException!
    }
}

// ✅ SOLUTION 1: Pass context explicitly
@Service
public class ExplicitContextPassing {
    
    public void processRequest() {
        UserContext context = UserContextHolder.getContext();
        asyncProcess(context);  // Pass explicitly
    }
    
    @Async
    public void asyncProcess(UserContext context) {
        // ✅ Context available as parameter
        log.info("User: {}", context.getUserId());
    }
}

// ✅ SOLUTION 2: TaskDecorator to propagate context
@Configuration
@EnableAsync
public class AsyncConfig implements AsyncConfigurer {
    
    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(100);
        
        // Decorate tasks to copy ThreadLocal values
        executor.setTaskDecorator(new ContextCopyingDecorator());
        
        executor.initialize();
        return executor;
    }
}

public class ContextCopyingDecorator implements TaskDecorator {
    
    @Override
    public Runnable decorate(Runnable runnable) {
        // Capture context from calling thread
        UserContext context = UserContextHolder.getContext();
        String correlationId = MDC.get("correlationId");
        
        return () -> {
            try {
                // Set context in worker thread
                UserContextHolder.setContext(context);
                MDC.put("correlationId", correlationId);
                
                runnable.run();
                
            } finally {
                // Clean up worker thread
                UserContextHolder.clear();
                MDC.clear();
            }
        };
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// PITFALL 3: ThreadLocal with parallel streams
// ══════════════════════════════════════════════════════════════════════════════

// ❌ WRONG: Parallel stream uses ForkJoinPool threads
@Service
public class ParallelStreamProblem {
    
    public void processBatch(List<Order> orders) {
        UserContextHolder.setContext(new UserContext("admin", "ADMIN"));
        
        orders.parallelStream()  // ❌ Different threads!
            .forEach(order -> {
                // Context might be null on ForkJoinPool threads
                UserContext ctx = UserContextHolder.getContext();  // ❌
                processOrder(order, ctx);
            });
    }
}

// ✅ CORRECT: Capture before stream, pass to lambdas
@Service
public class ParallelStreamSolution {
    
    public void processBatch(List<Order> orders) {
        // Capture first
        final UserContext context = UserContextHolder.getContext();
        
        orders.parallelStream()
            .forEach(order -> {
                // Use captured context (effectively final)
                processOrder(order, context);  // ✅
            });
    }
}
```

#### **Spring's Built-in ThreadLocal Usage**

```java
// Spring uses ThreadLocal extensively. Be aware of these:

// 1. RequestContextHolder - holds current request attributes
ServletRequestAttributes attrs = (ServletRequestAttributes) 
    RequestContextHolder.getRequestAttributes();
HttpServletRequest request = attrs.getRequest();

// 2. SecurityContextHolder - holds authentication
Authentication auth = SecurityContextHolder.getContext().getAuthentication();

// 3. TransactionSynchronizationManager - holds transaction state
boolean inTransaction = TransactionSynchronizationManager.isActualTransactionActive();

// 4. LocaleContextHolder - holds current locale
Locale locale = LocaleContextHolder.getLocale();

// These all have the same @Async propagation issue!
// Use TaskDecorator to propagate them:

public class SpringContextDecorator implements TaskDecorator {
    
    @Override
    public Runnable decorate(Runnable runnable) {
        // Capture Spring contexts
        RequestAttributes requestAttributes = RequestContextHolder.getRequestAttributes();
        SecurityContext securityContext = SecurityContextHolder.getContext();
        Map<String, String> mdcContext = MDC.getCopyOfContextMap();
        
        return () -> {
            try {
                // Restore in worker thread
                RequestContextHolder.setRequestAttributes(requestAttributes);
                SecurityContextHolder.setContext(securityContext);
                if (mdcContext != null) {
                    MDC.setContextMap(mdcContext);
                }
                
                runnable.run();
                
            } finally {
                // Clean up
                RequestContextHolder.resetRequestAttributes();
                SecurityContextHolder.clearContext();
                MDC.clear();
            }
        };
    }
}
```

---

## 10. Best Practices in Production

Production-ready multithreaded applications require careful attention to performance, monitoring, and reliability.

### 10.1 Avoid Blocking Calls in Async Threads

#### **The Problem**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    BLOCKING CALLS IN ASYNC THREADS                            │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  Your async thread pool has LIMITED threads (e.g., 10 threads)                │
│  If threads block on I/O, you run out of workers!                             │
│                                                                               │
│  SCENARIO: 100 async tasks, 10 threads, each task blocks 5 seconds on I/O     │
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │                    Thread Pool (10 threads)                             │  │
│  │                                                                         │  │
│  │  Thread-1: [████████ BLOCKING ON HTTP ████████] → Released after 5s     │  │
│  │  Thread-2: [████████ BLOCKING ON HTTP ████████]                         │  │
│  │  Thread-3: [████████ BLOCKING ON HTTP ████████]                         │  │
│  │  ...                                                                    │  │
│  │  Thread-10: [███████ BLOCKING ON HTTP ████████]                         │  │
│  │                                                                         │  │
│  │  Tasks 11-100: ⏳ WAITING IN QUEUE (no threads available!)              │  │
│  │                                                                         │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                               │
│  Time to complete 100 tasks: (100/10) × 5s = 50 SECONDS!                      │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

#### **Solutions**

```java
// ══════════════════════════════════════════════════════════════════════════════
// SOLUTION 1: Use WebClient instead of RestTemplate for HTTP calls
// ══════════════════════════════════════════════════════════════════════════════

@Service
@RequiredArgsConstructor
public class NonBlockingHttpService {
    
    private final WebClient webClient;
    
    // ❌ BLOCKING: RestTemplate blocks thread
    public Data fetchDataBlocking(String id) {
        return restTemplate.getForObject("/api/data/" + id, Data.class);
    }
    
    // ✅ NON-BLOCKING: WebClient returns immediately
    public CompletableFuture<Data> fetchDataNonBlocking(String id) {
        return webClient.get()
            .uri("/api/data/{id}", id)
            .retrieve()
            .bodyToMono(Data.class)
            .toFuture();  // Convert to CompletableFuture
    }
    
    // Batch processing - all calls parallel!
    public CompletableFuture<List<Data>> fetchAllData(List<String> ids) {
        List<CompletableFuture<Data>> futures = ids.stream()
            .map(this::fetchDataNonBlocking)
            .collect(Collectors.toList());
        
        return CompletableFuture.allOf(futures.toArray(new CompletableFuture[0]))
            .thenApply(v -> futures.stream()
                .map(CompletableFuture::join)
                .collect(Collectors.toList()));
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// SOLUTION 2: Separate thread pools for different workloads
// ══════════════════════════════════════════════════════════════════════════════

@Configuration
@EnableAsync
public class MultiPoolConfig {
    
    // CPU-bound tasks: threads = CPU cores
    @Bean("cpuBoundExecutor")
    public TaskExecutor cpuBoundExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        int cores = Runtime.getRuntime().availableProcessors();
        executor.setCorePoolSize(cores);
        executor.setMaxPoolSize(cores);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("cpu-");
        executor.initialize();
        return executor;
    }
    
    // I/O-bound tasks: more threads to handle blocking
    @Bean("ioBoundExecutor")
    public TaskExecutor ioBoundExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(50);   // Many threads for blocking I/O
        executor.setMaxPoolSize(200);
        executor.setQueueCapacity(500);
        executor.setThreadNamePrefix("io-");
        executor.initialize();
        return executor;
    }
}

@Service
public class TaskService {
    
    // Route CPU work to CPU pool
    @Async("cpuBoundExecutor")
    public CompletableFuture<Result> computeHeavyCalculation(Data input) {
        // No blocking, just CPU work
        return CompletableFuture.completedFuture(calculate(input));
    }
    
    // Route I/O work to I/O pool
    @Async("ioBoundExecutor")
    public CompletableFuture<Data> fetchExternalData(String id) {
        // This will block - that's OK in I/O pool
        Data data = restTemplate.getForObject("/api/" + id, Data.class);
        return CompletableFuture.completedFuture(data);
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// SOLUTION 3: Use timeouts to prevent indefinite blocking
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class TimeoutAwareService {
    
    @Autowired
    private RestTemplate restTemplate;
    
    public Data fetchWithTimeout(String id) {
        // Configure RestTemplate with timeouts
        HttpComponentsClientHttpRequestFactory factory = 
            new HttpComponentsClientHttpRequestFactory();
        factory.setConnectTimeout(5000);   // 5 second connection timeout
        factory.setReadTimeout(10000);     // 10 second read timeout
        
        RestTemplate timeoutTemplate = new RestTemplate(factory);
        
        try {
            return timeoutTemplate.getForObject("/api/" + id, Data.class);
        } catch (ResourceAccessException e) {
            log.warn("Request timed out for id: {}", id);
            return getDefaultData();  // Fallback
        }
    }
    
    // CompletableFuture with timeout
    public CompletableFuture<Data> fetchAsyncWithTimeout(String id) {
        return CompletableFuture.supplyAsync(() -> fetchData(id))
            .orTimeout(10, TimeUnit.SECONDS)
            .exceptionally(ex -> {
                if (ex instanceof TimeoutException) {
                    log.warn("Async operation timed out for id: {}", id);
                    return getDefaultData();
                }
                throw new RuntimeException(ex);
            });
    }
}
```

---

### 10.2 Proper Executor Sizing

#### **Sizing Formulas**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    THREAD POOL SIZING GUIDE                                   │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  CPU-BOUND TASKS (calculations, transformations)                              │
│  ════════════════════════════════════════════════                             │
│                                                                               │
│  Formula: threads = CPU cores                                                 │
│                                                                               │
│  • More threads = more context switching overhead                             │
│  • CPU can only run (cores) threads simultaneously                            │
│  • Extra threads just wait and waste memory                                   │
│                                                                               │
│  Example: 8-core CPU → 8 threads for CPU-bound work                           │
│                                                                               │
│  ─────────────────────────────────────────────────────────────────────────    │
│                                                                               │
│  I/O-BOUND TASKS (HTTP calls, database, file I/O)                             │
│  ════════════════════════════════════════════════                             │
│                                                                               │
│  Formula: threads = CPU cores × (1 + wait_time / compute_time)                │
│                                                                               │
│  • Threads spend most time waiting for I/O                                    │
│  • While one thread waits, others can use CPU                                 │
│  • Need enough threads to keep CPU busy                                       │
│                                                                               │
│  Example: 8 cores, tasks wait 200ms, compute 20ms                             │
│  threads = 8 × (1 + 200/20) = 8 × 11 = 88 threads                             │
│                                                                               │
│  ─────────────────────────────────────────────────────────────────────────    │
│                                                                               │
│  MIXED WORKLOADS                                                              │
│  ═══════════════                                                              │
│                                                                               │
│  • Create separate pools for CPU and I/O tasks                                │
│  • Size each pool appropriately                                               │
│  • Monitor and adjust based on actual metrics                                 │
│                                                                               │
│  ─────────────────────────────────────────────────────────────────────────    │
│                                                                               │
│  PRACTICAL DEFAULTS                                                           │
│  ═════════════════                                                            │
│                                                                               │
│  • Start with: cores × 2 for general async work                               │
│  • Queue: 100-500 depending on burst expectations                             │
│  • Monitor metrics and adjust                                                 │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

```java
@Configuration
@EnableAsync
public class OptimizedExecutorConfig {
    
    @Bean("optimizedExecutor")
    public TaskExecutor optimizedExecutor(
            @Value("${async.wait-time-ms:200}") int waitTimeMs,
            @Value("${async.compute-time-ms:20}") int computeTimeMs) {
        
        int cpuCores = Runtime.getRuntime().availableProcessors();
        
        // Apply Little's Law: L = λW
        // threads = cores × (1 + wait/compute)
        int optimalThreads = (int) (cpuCores * (1 + (double) waitTimeMs / computeTimeMs));
        
        // Cap at reasonable maximum
        int maxThreads = Math.min(optimalThreads * 2, 200);
        
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(optimalThreads);
        executor.setMaxPoolSize(maxThreads);
        executor.setQueueCapacity(optimalThreads * 10);  // ~10 tasks per thread queued
        executor.setKeepAliveSeconds(60);
        executor.setThreadNamePrefix("opt-async-");
        executor.setRejectionHandler(new CallerRunsPolicy());
        executor.initialize();
        
        log.info("Initialized executor: cores={}, optimalThreads={}, maxThreads={}",
            cpuCores, optimalThreads, maxThreads);
        
        return executor;
    }
}
```

---

### 10.3 Monitoring Thread Pools

#### **Spring Boot Actuator Integration**

```yaml
# application.yml
management:
  endpoints:
    web:
      exposure:
        include: health,metrics,prometheus,threaddump
  endpoint:
    health:
      show-details: always
  metrics:
    tags:
      application: ${spring.application.name}
```

```java
@Configuration
@EnableAsync
public class MonitoredExecutorConfig {
    
    @Bean("monitoredExecutor")
    public TaskExecutor monitoredExecutor(MeterRegistry meterRegistry) {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(50);
        executor.setQueueCapacity(200);
        executor.setThreadNamePrefix("async-");
        executor.initialize();
        
        // Register metrics with Micrometer
        ThreadPoolExecutor threadPool = executor.getThreadPoolExecutor();
        
        Gauge.builder("async.pool.size", threadPool, ThreadPoolExecutor::getPoolSize)
            .tag("executor", "async")
            .description("Current pool size")
            .register(meterRegistry);
        
        Gauge.builder("async.pool.active", threadPool, ThreadPoolExecutor::getActiveCount)
            .tag("executor", "async")
            .description("Active threads")
            .register(meterRegistry);
        
        Gauge.builder("async.queue.size", threadPool, e -> e.getQueue().size())
            .tag("executor", "async")
            .description("Queue size")
            .register(meterRegistry);
        
        Gauge.builder("async.queue.remaining", threadPool, e -> e.getQueue().remainingCapacity())
            .tag("executor", "async")
            .description("Queue remaining capacity")
            .register(meterRegistry);
        
        FunctionCounter.builder("async.tasks.completed", threadPool, 
                ThreadPoolExecutor::getCompletedTaskCount)
            .tag("executor", "async")
            .description("Completed tasks")
            .register(meterRegistry);
        
        return executor;
    }
}

// Custom health indicator
@Component
public class ThreadPoolHealthIndicator implements HealthIndicator {
    
    @Autowired
    @Qualifier("monitoredExecutor")
    private TaskExecutor taskExecutor;
    
    @Override
    public Health health() {
        if (taskExecutor instanceof ThreadPoolTaskExecutor executor) {
            ThreadPoolExecutor pool = executor.getThreadPoolExecutor();
            
            int queueSize = pool.getQueue().size();
            int queueCapacity = queueSize + pool.getQueue().remainingCapacity();
            double queueUsage = (double) queueSize / queueCapacity;
            
            Map<String, Object> details = new HashMap<>();
            details.put("activeThreads", pool.getActiveCount());
            details.put("poolSize", pool.getPoolSize());
            details.put("corePoolSize", pool.getCorePoolSize());
            details.put("maxPoolSize", pool.getMaximumPoolSize());
            details.put("queueSize", queueSize);
            details.put("queueCapacity", queueCapacity);
            details.put("queueUsagePercent", String.format("%.1f%%", queueUsage * 100));
            details.put("completedTasks", pool.getCompletedTaskCount());
            
            // Health status based on queue usage
            if (queueUsage > 0.9) {
                return Health.down()
                    .withDetails(details)
                    .withDetail("reason", "Queue nearly full")
                    .build();
            } else if (queueUsage > 0.7) {
                return Health.status("WARNING")
                    .withDetails(details)
                    .withDetail("reason", "Queue usage high")
                    .build();
            }
            
            return Health.up().withDetails(details).build();
        }
        
        return Health.unknown().build();
    }
}
```

#### **Custom Metrics Dashboard Queries (Prometheus/Grafana)**

```promql
# Alert when queue is filling up
async_queue_size{executor="async"} / 
  (async_queue_size{executor="async"} + async_queue_remaining{executor="async"}) > 0.8

# Thread pool utilization
async_pool_active{executor="async"} / async_pool_size{executor="async"}

# Task throughput (tasks/minute)
rate(async_tasks_completed{executor="async"}[1m]) * 60

# Queue wait time estimate
async_queue_size{executor="async"} / 
  (rate(async_tasks_completed{executor="async"}[5m]) + 0.001)
```

---

### 10.4 Handling Backpressure

#### **What is Backpressure?**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    BACKPRESSURE EXPLAINED                                     │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  Producer generates work faster than consumer can process                     │
│                                                                               │
│  Without backpressure handling:                                               │
│                                                                               │
│  Producer ──100 req/s──► [█████████████████ QUEUE OVERFLOW] ──► Consumer      │
│  (fast)                        ▲                              (slow: 50 req/s)│
│                                │                                              │
│                       OutOfMemoryError!                                       │
│                                                                               │
│  ─────────────────────────────────────────────────────────────────────────    │
│                                                                               │
│  With backpressure handling:                                                  │
│                                                                               │
│  Producer ──50 req/s──► [███████░░░░░░░░░░ BOUNDED QUEUE] ──► Consumer        │
│  (throttled)                     ▲                            (50 req/s)      │
│                                  │                                            │
│                        Signal to slow down                                    │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

#### **Backpressure Strategies**

```java
@Configuration
public class BackpressureConfig {
    
    // Strategy 1: Bounded queue + CallerRunsPolicy
    // When queue is full, producer thread executes the task itself (slows down)
    @Bean("callerRunsExecutor")
    public TaskExecutor callerRunsExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);  // Bounded queue
        executor.setRejectionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        // When queue full: calling thread runs the task → natural backpressure
        executor.initialize();
        return executor;
    }
    
    // Strategy 2: Blocking offer with timeout
    @Bean("blockingExecutor")
    public TaskExecutor blockingExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor() {
            @Override
            protected BlockingQueue<Runnable> createQueue(int queueCapacity) {
                // ArrayBlockingQueue for blocking behavior
                return new ArrayBlockingQueue<>(queueCapacity);
            }
        };
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);
        executor.initialize();
        return executor;
    }
}

@Service
@RequiredArgsConstructor
@Slf4j
public class BackpressureAwareService {
    
    @Qualifier("callerRunsExecutor")
    private final TaskExecutor executor;
    
    // Strategy 3: Rate limiting with Semaphore
    private final Semaphore permits = new Semaphore(100);  // Max 100 concurrent
    
    public CompletableFuture<Result> processWithBackpressure(Request request) {
        // Try to acquire permit
        if (!permits.tryAcquire()) {
            log.warn("Rate limit exceeded, rejecting request");
            return CompletableFuture.failedFuture(
                new RateLimitExceededException("Too many concurrent requests"));
        }
        
        return CompletableFuture.supplyAsync(() -> {
            try {
                return doProcess(request);
            } finally {
                permits.release();  // Release permit when done
            }
        }, executor);
    }
    
    // Strategy 4: Reactive backpressure with Project Reactor
    public Flux<Result> processStreamWithBackpressure(Flux<Request> requests) {
        return requests
            .limitRate(100)              // Request only 100 at a time from upstream
            .onBackpressureBuffer(500)   // Buffer up to 500
            .onBackpressureDrop(req -> log.warn("Dropped request: {}", req.getId()))
            .flatMap(
                this::processReactive,
                10                        // Max 10 concurrent
            );
    }
    
    private Mono<Result> processReactive(Request request) {
        return Mono.fromCallable(() -> doProcess(request))
            .subscribeOn(Schedulers.boundedElastic());
    }
}

// Strategy 5: Circuit breaker pattern (Resilience4j)
@Service
public class CircuitBreakerService {
    
    @CircuitBreaker(name = "externalService", fallbackMethod = "fallback")
    @Bulkhead(name = "externalService", type = Bulkhead.Type.THREADPOOL)
    @RateLimiter(name = "externalService")
    public CompletableFuture<Response> callExternalService(Request request) {
        return webClient.post()
            .uri("/api/process")
            .bodyValue(request)
            .retrieve()
            .bodyToMono(Response.class)
            .toFuture();
    }
    
    public CompletableFuture<Response> fallback(Request request, Exception e) {
        log.warn("Circuit breaker fallback for request: {}", request.getId());
        return CompletableFuture.completedFuture(Response.defaultResponse());
    }
}
```

```yaml
# Resilience4j configuration
resilience4j:
  circuitbreaker:
    instances:
      externalService:
        slidingWindowSize: 10
        failureRateThreshold: 50
        waitDurationInOpenState: 10s
  
  bulkhead:
    instances:
      externalService:
        maxConcurrentCalls: 25
        maxWaitDuration: 500ms
  
  ratelimiter:
    instances:
      externalService:
        limitForPeriod: 100
        limitRefreshPeriod: 1s
        timeoutDuration: 0s  # Fail fast
```

---

### 10.5 Logging and Debugging Async Code

#### **Challenge: Lost Context in Async**

```java
// Problem: Log correlation is lost in async execution
// Request-1: [correlationId=abc123] Starting process
// Request-2: [correlationId=def456] Starting process
// Async:     [correlationId=???] Processing completed   ← Which request?
```

#### **Solution: MDC Propagation**

```java
@Configuration
@EnableAsync
public class LoggingAsyncConfig implements AsyncConfigurer {
    
    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(50);
        executor.setQueueCapacity(200);
        executor.setThreadNamePrefix("async-");
        executor.setTaskDecorator(new MdcTaskDecorator());
        executor.initialize();
        return executor;
    }
}

public class MdcTaskDecorator implements TaskDecorator {
    
    @Override
    public Runnable decorate(Runnable runnable) {
        // Capture MDC from calling thread
        Map<String, String> contextMap = MDC.getCopyOfContextMap();
        
        return () -> {
            try {
                // Set MDC in worker thread
                if (contextMap != null) {
                    MDC.setContextMap(contextMap);
                }
                runnable.run();
            } finally {
                MDC.clear();
            }
        };
    }
}

// Filter to set MDC at request entry
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class MdcFilter implements Filter {
    
    @Override
    public void doFilter(ServletRequest request, ServletResponse response, 
                         FilterChain chain) throws IOException, ServletException {
        try {
            String correlationId = ((HttpServletRequest) request)
                .getHeader("X-Correlation-Id");
            if (correlationId == null) {
                correlationId = UUID.randomUUID().toString().substring(0, 8);
            }
            
            MDC.put("correlationId", correlationId);
            MDC.put("requestPath", ((HttpServletRequest) request).getRequestURI());
            
            chain.doFilter(request, response);
            
        } finally {
            MDC.clear();
        }
    }
}

// logback-spring.xml
// <pattern>%d{HH:mm:ss.SSS} [%thread] [%X{correlationId}] %-5level %logger{36} - %msg%n</pattern>
```

#### **Debugging Async Issues**

```java
@Service
@Slf4j
public class DebuggableAsyncService {
    
    @Async
    public CompletableFuture<Result> processWithDebugging(Request request) {
        String taskId = UUID.randomUUID().toString().substring(0, 8);
        
        log.debug("[Task-{}] Started on thread: {} - Request: {}",
            taskId, Thread.currentThread().getName(), request.getId());
        
        Instant start = Instant.now();
        
        try {
            Result result = doProcess(request);
            
            log.debug("[Task-{}] Completed in {}ms on thread: {}",
                taskId, 
                Duration.between(start, Instant.now()).toMillis(),
                Thread.currentThread().getName());
            
            return CompletableFuture.completedFuture(result);
            
        } catch (Exception e) {
            log.error("[Task-{}] Failed after {}ms on thread: {} - Error: {}",
                taskId,
                Duration.between(start, Instant.now()).toMillis(),
                Thread.currentThread().getName(),
                e.getMessage(),
                e);  // Include stack trace for errors
            
            return CompletableFuture.failedFuture(e);
        }
    }
    
    // Debug utility: Log thread pool state
    public void logPoolState(ThreadPoolTaskExecutor executor) {
        ThreadPoolExecutor pool = executor.getThreadPoolExecutor();
        
        log.info("Thread Pool State: active={}/{}, pool={}/{}, queue={}/{}, completed={}",
            pool.getActiveCount(),
            pool.getMaximumPoolSize(),
            pool.getPoolSize(),
            pool.getMaximumPoolSize(),
            pool.getQueue().size(),
            pool.getQueue().size() + pool.getQueue().remainingCapacity(),
            pool.getCompletedTaskCount());
    }
}

// Actuator endpoint for thread dump
// GET /actuator/threaddump

// Custom endpoint for pool details
@RestController
@RequestMapping("/debug")
@Profile("debug")  // Only in debug profile
public class DebugController {
    
    @Autowired
    private Map<String, ThreadPoolTaskExecutor> executors;
    
    @GetMapping("/thread-pools")
    public Map<String, Map<String, Object>> getThreadPoolStats() {
        return executors.entrySet().stream()
            .collect(Collectors.toMap(
                Map.Entry::getKey,
                e -> getPoolStats(e.getValue())
            ));
    }
    
    private Map<String, Object> getPoolStats(ThreadPoolTaskExecutor executor) {
        ThreadPoolExecutor pool = executor.getThreadPoolExecutor();
        return Map.of(
            "activeCount", pool.getActiveCount(),
            "poolSize", pool.getPoolSize(),
            "corePoolSize", pool.getCorePoolSize(),
            "maxPoolSize", pool.getMaximumPoolSize(),
            "queueSize", pool.getQueue().size(),
            "queueCapacity", pool.getQueue().remainingCapacity(),
            "completedTasks", pool.getCompletedTaskCount(),
            "rejectedTasks", getRejectedCount(pool)
        );
    }
}
```

---

### 10.6 Avoiding Memory Leaks

#### **Common Memory Leak Sources**

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                    MEMORY LEAKS IN MULTITHREADED CODE                         │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  1. ThreadLocal not cleaned up                                                │
│  ─────────────────────────────                                                │
│  • Threads in pool are reused                                                 │
│  • ThreadLocal values persist across requests                                 │
│  • Memory grows with each new value                                           │
│                                                                               │
│  2. Listeners/Callbacks not deregistered                                      │
│  ─────────────────────────────────────────                                    │
│  • Object registers as listener                                               │
│  • Object goes out of scope but listener holds reference                      │
│  • GC cannot collect the object                                               │
│                                                                               │
│  3. Unbounded caches/collections                                              │
│  ────────────────────────────────                                             │
│  • Adding to cache without eviction                                           │
│  • Map keys that keep growing                                                 │
│  • Collections that are never cleaned                                         │
│                                                                               │
│  4. Futures that are never consumed                                           │
│  ──────────────────────────────────                                           │
│  • Creating CompletableFuture but never joining                               │
│  • Results accumulate in memory                                               │
│                                                                               │
│  5. ExecutorService not shut down                                             │
│  ────────────────────────────────                                             │
│  • Threads keep running after app should stop                                 │
│  • In tests, creates thread leaks                                             │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

#### **Prevention Patterns**

```java
// ══════════════════════════════════════════════════════════════════════════════
// LEAK 1: ThreadLocal - Always clean up!
// ══════════════════════════════════════════════════════════════════════════════

public class SafeThreadLocalUsage {
    
    private static final ThreadLocal<ExpensiveObject> context = new ThreadLocal<>();
    
    // ❌ LEAK: No cleanup
    public void leakyMethod() {
        context.set(new ExpensiveObject());
        doWork();
        // ThreadLocal value stays forever!
    }
    
    // ✅ SAFE: Always cleanup in finally
    public void safeMethod() {
        try {
            context.set(new ExpensiveObject());
            doWork();
        } finally {
            context.remove();  // ALWAYS clean up!
        }
    }
    
    // ✅ BETTER: Use try-with-resources pattern
    public void betterMethod() {
        try (var ctx = new ThreadLocalContext<>(context, new ExpensiveObject())) {
            doWork();
        }  // Auto-cleaned
    }
}

// Auto-closeable ThreadLocal wrapper
public class ThreadLocalContext<T> implements AutoCloseable {
    private final ThreadLocal<T> threadLocal;
    
    public ThreadLocalContext(ThreadLocal<T> threadLocal, T value) {
        this.threadLocal = threadLocal;
        threadLocal.set(value);
    }
    
    @Override
    public void close() {
        threadLocal.remove();
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// LEAK 2: Unbounded cache - Use bounded cache with eviction
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class SafeCacheService {
    
    // ❌ LEAK: Unbounded HashMap
    private final Map<String, Data> leakyCache = new HashMap<>();
    
    // ✅ SAFE: Bounded with size limit (using Caffeine)
    private final Cache<String, Data> safeCache = Caffeine.newBuilder()
        .maximumSize(10_000)
        .expireAfterWrite(Duration.ofMinutes(10))
        .recordStats()
        .build();
    
    // ✅ ALTERNATIVE: LRU LinkedHashMap
    private final Map<String, Data> lruCache = Collections.synchronizedMap(
        new LinkedHashMap<>(100, 0.75f, true) {
            @Override
            protected boolean removeEldestEntry(Map.Entry<String, Data> eldest) {
                return size() > 1000;  // Max 1000 entries
            }
        }
    );
    
    // ✅ WeakHashMap for caches that can be GC'd
    private final Map<Key, Data> weakCache = 
        Collections.synchronizedMap(new WeakHashMap<>());
}

// ══════════════════════════════════════════════════════════════════════════════
// LEAK 3: Executor shutdown - Always shutdown!
// ══════════════════════════════════════════════════════════════════════════════

@Configuration
public class ExecutorLifecycleConfig {
    
    @Bean(destroyMethod = "shutdown")  // Automatic shutdown on context close
    public ExecutorService executorService() {
        return Executors.newFixedThreadPool(10);
    }
    
    // Spring's ThreadPoolTaskExecutor handles this automatically
    @Bean
    public ThreadPoolTaskExecutor taskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setWaitForTasksToCompleteOnShutdown(true);
        executor.setAwaitTerminationSeconds(30);
        executor.initialize();
        return executor;
    }
}

// In tests - always shutdown!
@Test
void testWithExecutor() {
    ExecutorService executor = Executors.newFixedThreadPool(5);
    try {
        // Test code
    } finally {
        executor.shutdown();
        if (!executor.awaitTermination(5, TimeUnit.SECONDS)) {
            executor.shutdownNow();
        }
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// LEAK 4: Futures - Always consume or cancel
// ══════════════════════════════════════════════════════════════════════════════

@Service
@Slf4j
public class FutureLeakPrevention {
    
    // ❌ LEAK: Fire and forget without handling
    public void leakyFireAndForget() {
        CompletableFuture.supplyAsync(() -> heavyComputation());
        // Future never consumed - result sits in memory
    }
    
    // ✅ SAFE: If you don't need result, use runAsync
    public void safeFireAndForget() {
        CompletableFuture.runAsync(() -> {
            heavyComputation();  // Result not stored
        }).exceptionally(e -> {
            log.error("Background task failed", e);
            return null;
        });
    }
    
    // ✅ SAFE: Track futures and clean up periodically
    private final Set<CompletableFuture<?>> pendingFutures = 
        ConcurrentHashMap.newKeySet();
    
    public <T> CompletableFuture<T> trackFuture(CompletableFuture<T> future) {
        pendingFutures.add(future);
        return future.whenComplete((r, e) -> pendingFutures.remove(future));
    }
    
    @Scheduled(fixedRate = 60000)
    public void cleanupStaleFutures() {
        pendingFutures.removeIf(f -> f.isDone() || f.isCancelled());
        log.debug("Pending futures count: {}", pendingFutures.size());
    }
}

// ══════════════════════════════════════════════════════════════════════════════
// LEAK 5: Listeners - Use weak references or explicit cleanup
// ══════════════════════════════════════════════════════════════════════════════

@Service
public class SafeEventService {
    
    // ❌ LEAK: Strong references to listeners
    private final List<EventListener> leakyListeners = new ArrayList<>();
    
    // ✅ SAFE: Weak references allow GC
    private final List<WeakReference<EventListener>> safeListeners = 
        new CopyOnWriteArrayList<>();
    
    public void addListener(EventListener listener) {
        safeListeners.add(new WeakReference<>(listener));
    }
    
    public void fireEvent(Event event) {
        Iterator<WeakReference<EventListener>> iter = safeListeners.iterator();
        while (iter.hasNext()) {
            EventListener listener = iter.next().get();
            if (listener == null) {
                // Listener was GC'd - remove weak reference
                iter.remove();
            } else {
                listener.onEvent(event);
            }
        }
    }
}
```

#### **Memory Leak Detection**

```java
@Component
@Slf4j
public class MemoryMonitor {
    
    @Scheduled(fixedRate = 60000)
    public void logMemoryUsage() {
        Runtime runtime = Runtime.getRuntime();
        
        long totalMemory = runtime.totalMemory();
        long freeMemory = runtime.freeMemory();
        long usedMemory = totalMemory - freeMemory;
        long maxMemory = runtime.maxMemory();
        
        double usagePercent = (double) usedMemory / maxMemory * 100;
        
        log.info("Memory: used={}MB, total={}MB, max={}MB, usage={:.1f}%",
            usedMemory / 1024 / 1024,
            totalMemory / 1024 / 1024,
            maxMemory / 1024 / 1024,
            usagePercent);
        
        if (usagePercent > 80) {
            log.warn("High memory usage detected: {:.1f}%", usagePercent);
        }
    }
    
    // Force GC and check memory (for debugging only!)
    @PostMapping("/debug/gc")
    @Profile("debug")
    public Map<String, Long> forceGC() {
        long beforeGC = Runtime.getRuntime().totalMemory() - 
                        Runtime.getRuntime().freeMemory();
        
        System.gc();
        
        long afterGC = Runtime.getRuntime().totalMemory() - 
                       Runtime.getRuntime().freeMemory();
        
        return Map.of(
            "beforeGC_MB", beforeGC / 1024 / 1024,
            "afterGC_MB", afterGC / 1024 / 1024,
            "freed_MB", (beforeGC - afterGC) / 1024 / 1024
        );
    }
}
```

---

## Summary - Sections 9 & 10

| Topic | Key Points |
|-------|------------|
| **Race Conditions** | Use AtomicInteger, ConcurrentHashMap, synchronized blocks |
| **Singleton Thread Safety** | Keep beans stateless; use thread-safe collections |
| **Prototype Scope** | Use ObjectProvider or @Lookup for proper prototype injection |
| **ThreadLocal** | Always clean up in finally; use TaskDecorator for @Async |
| **Blocking in Async** | Use WebClient; separate pools for CPU/IO; set timeouts |
| **Executor Sizing** | CPU: cores; I/O: cores × (1 + wait/compute) |
| **Monitoring** | Actuator metrics; custom health indicators; Prometheus/Grafana |
| **Backpressure** | CallerRunsPolicy; Semaphore; Circuit breaker; Reactive limitRate |
| **Logging Async** | Propagate MDC with TaskDecorator; correlation IDs |
| **Memory Leaks** | Clean ThreadLocal; bounded caches; shutdown executors; consume futures |

---

*Next sections will cover:*
- *Virtual Threads (Java 21+)*
- *Testing Multithreaded Code*
- *Interview Questions & Answers*
