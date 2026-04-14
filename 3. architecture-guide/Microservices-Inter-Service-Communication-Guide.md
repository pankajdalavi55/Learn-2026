# Microservices Inter-Service Communication — Complete Production Guide

## Table of Contents
1. [Core Fundamentals](#1-core-fundamentals)
2. [REST-Based Communication (Spring Boot)](#2-rest-based-communication-spring-boot)
3. [Design & Architecture](#3-design--architecture)
4. [Synchronous vs Asynchronous](#4-synchronous-vs-asynchronous)
5. [Fault Tolerance & Resilience](#5-fault-tolerance--resilience)
6. [Real Production Problems](#6-real-production-problems)
7. [Data Consistency & Transactions](#7-data-consistency--transactions)
8. [Security in Inter-Service Communication](#8-security-in-inter-service-communication)
9. [gRPC vs REST](#9-grpc-vs-rest)
10. [Performance & Optimization](#10-performance--optimization)
11. [Observability & Debugging](#11-observability--debugging)
12. [Testing Inter-Service Communication](#12-testing-inter-service-communication)
13. [Scenario-Based Questions](#13-scenario-based-questions-interview-gold)

---

## 1. Core Fundamentals

### 1.1 What Is Inter-Service Communication in Microservices?

Inter-service communication is **how independent microservices talk to each other** to fulfill business operations that span multiple service boundaries. In a monolith, this is a simple method call. In microservices, it's a network call — and that changes everything.

```
Monolith:                            Microservices:
┌──────────────────────┐             ┌─────────┐   HTTP/gRPC   ┌─────────┐
│  OrderModule          │             │  Order  │──────────────→│ Payment │
│    .placeOrder()      │             │ Service │               │ Service │
│       │               │             └────┬────┘               └─────────┘
│       ▼               │                  │
│  PaymentModule        │                  │  Kafka Event
│    .charge()          │                  ▼
│       │               │             ┌─────────┐
│       ▼               │             │Inventory│
│  InventoryModule      │             │ Service │
│    .reserve()         │             └─────────┘
└──────────────────────┘
  In-process method call              Network calls (latency, failures, serialization)
```

**Why this is fundamentally different:**

| Aspect | Monolith (Method Call) | Microservices (Network Call) |
|--------|----------------------|------------------------------|
| **Latency** | Nanoseconds | Milliseconds to seconds |
| **Failure** | Process crash = all fail | Partial failure possible |
| **Serialization** | None (shared memory) | JSON/Protobuf/Avro encoding |
| **Transaction** | Single DB transaction | Distributed — no ACID guarantee |
| **Discovery** | Compile-time reference | Runtime discovery needed |
| **Versioning** | Deploy together | Independent deployments |

---

### 1.2 Types of Communication

```
Inter-Service Communication
        │
        ├── By Interaction Style
        │     ├── Request/Response (ask and wait for answer)
        │     └── Event-Driven (fire and forget / pub-sub)
        │
        ├── By Timing
        │     ├── Synchronous (caller blocks until response)
        │     └── Asynchronous (caller doesn't block)
        │
        └── By Participants
              ├── One-to-One (point-to-point)
              └── One-to-Many (pub-sub broadcast)
```

| Pattern | Style | Example | When to Use |
|---------|-------|---------|-------------|
| **Sync Request/Response** | 1:1, blocking | REST, gRPC | Need immediate answer (GET user, validate payment) |
| **Async Request/Response** | 1:1, non-blocking | Message queue with reply | Long operations (PDF generation, batch processing) |
| **Event Notification** | 1:many, fire-and-forget | Kafka topic publish | State changes others might care about |
| **Event-Carried State Transfer** | 1:many, data in event | Kafka with full payload | Consumers need data without callback |
| **Command** | 1:1, async | RabbitMQ queue | Trigger action in another service |

---

### 1.3 Synchronous vs Asynchronous

```
Synchronous:
┌─────────┐  request   ┌─────────┐
│ Service │────────────→│ Service │
│    A    │   blocked   │    B    │
│         │◄────────────│         │
└─────────┘  response   └─────────┘
    A waits for B to respond.
    If B is slow → A is slow.
    If B is down → A fails.

Asynchronous:
┌─────────┐  message   ┌──────────┐           ┌─────────┐
│ Service │────────────→│  Message │──────────→│ Service │
│    A    │  continues  │  Broker  │           │    B    │
│         │  working    │ (Kafka)  │           │         │
└─────────┘             └──────────┘           └─────────┘
    A publishes and moves on.
    B processes when ready.
    A and B are temporally decoupled.
```

| Aspect | Synchronous | Asynchronous |
|--------|------------|--------------|
| **Coupling** | Temporal + spatial | Only spatial (or neither with events) |
| **Latency** | Additive (A + B + C) | Non-blocking, parallel possible |
| **Failure handling** | Cascading failures | Isolated failures |
| **Complexity** | Simpler to implement | More infrastructure (broker, DLQ, idempotency) |
| **Data consistency** | Easier (request/response) | Eventual consistency |
| **Debugging** | Straightforward stack traces | Distributed tracing needed |
| **Best for** | Queries, validations | Commands, notifications, data sync |

---

### 1.4 When to Prefer REST over Messaging?

```
Use REST/Sync when:                    Use Messaging/Async when:
─────────────────────                  ──────────────────────────
✓ Need immediate response              ✓ Fire-and-forget is acceptable
✓ Simple CRUD operations               ✓ Long-running operations
✓ User is waiting for result            ✓ Multiple consumers need the data
✓ Low volume, low latency required      ✓ Need to absorb traffic spikes
✓ Read/query operations                 ✓ Services have different throughput
✓ Simple request/reply pattern          ✓ Event-driven workflow
✓ Few service hops (2-3 max)            ✓ Decoupling is priority
```

**Production decision matrix:**

| Scenario | Recommendation | Reason |
|----------|---------------|--------|
| Get user profile for UI | REST | User waiting, low latency needed |
| Validate payment card | REST/gRPC | Sync validation, immediate result |
| Send order confirmation email | Messaging | User doesn't wait for email delivery |
| Update search index after product change | Messaging | Eventual consistency acceptable |
| Real-time inventory check at checkout | REST | Must be current, user waiting |
| Process monthly billing for 1M users | Messaging | Long-running, needs queue/batching |
| Propagate user profile change to 10 services | Messaging (pub-sub) | Fan-out to many consumers |

---

### 1.5 Common Communication Protocols

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Communication Protocols                           │
├──────────────┬──────────────┬──────────────┬───────────────────────┤
│  HTTP/REST   │    gRPC      │  WebSocket   │  Messaging            │
│              │              │              │  (Kafka/RabbitMQ)     │
├──────────────┼──────────────┼──────────────┼───────────────────────┤
│ Text/JSON    │ Binary/Proto │ Full-duplex  │ Async/Durable         │
│ HTTP/1.1     │ HTTP/2       │ Persistent   │ Pub-Sub / Queue       │
│ Req/Response │ Streaming    │ Bi-directional│ Guaranteed delivery  │
│ Stateless    │ Typed (IDL)  │ Real-time    │ Ordered (partitioned) │
├──────────────┼──────────────┼──────────────┼───────────────────────┤
│ External APIs│ Internal svc │ Live updates │ Event-driven arch     │
│ CRUD ops     │ High perf    │ Chat/Gaming  │ Decoupled workflows   │
│ Public APIs  │ Internal comms│ Notifications│ Data pipelines       │
└──────────────┴──────────────┴──────────────┴───────────────────────┘
```

| Protocol | Latency | Throughput | Learning Curve | Browser Support | Use Case |
|----------|---------|------------|----------------|-----------------|----------|
| **HTTP/REST** | Medium | Medium | Low | Full | External/public APIs |
| **gRPC** | Low | High | Medium | Limited (gRPC-Web) | Internal high-perf service-to-service |
| **WebSocket** | Very Low | Medium | Medium | Full | Real-time bidirectional (chat, dashboards) |
| **Kafka** | Medium | Very High | High | None (backend only) | Event streaming, data pipelines |
| **RabbitMQ** | Low | High | Medium | None (backend only) | Task queues, request routing |

---

### 1.6 Tight Coupling vs Loose Coupling

```
Tight Coupling:                           Loose Coupling:
┌─────────┐                              ┌─────────┐
│ Order   │──── knows IP ────→           │ Order   │── publish event ──→┌────────┐
│ Service │──── knows API ────→Payment   │ Service │                    │ Kafka  │
│         │──── waits for ────→Service   │         │                    └───┬────┘
└─────────┘     response                 └─────────┘                        │
                                                              ┌─────────────┼──────────┐
    If Payment changes API → Order breaks                     ▼             ▼          ▼
    If Payment is down → Order fails                     Payment       Inventory   Analytics
    If Payment is slow → Order is slow                   Service       Service     Service
                                                         (each processes independently)
```

**Coupling dimensions:**

| Dimension | Tight | Loose | How to Achieve Loose |
|-----------|-------|-------|---------------------|
| **Temporal** | Must be online simultaneously | Can process at different times | Messaging/events |
| **Spatial** | Knows exact address/location | Doesn't know who consumes | Service registry, event bus |
| **Schema** | Shares data models / DTOs | Own internal models | API contracts, anti-corruption layer |
| **Technology** | Same language/framework | Any tech stack | Protocol-based integration (HTTP, gRPC) |
| **Deployment** | Deploy together | Deploy independently | Versioned APIs, backward compatibility |

**Production rule of thumb:** Prefer loose coupling. Accept the complexity of eventual consistency and messaging infrastructure. Tight coupling erases the benefits of microservices.

---

## 2. REST-Based Communication (Spring Boot)

### 2.1 RestTemplate vs WebClient vs OpenFeign

```
Evolution of HTTP Clients in Spring:

RestTemplate (Spring 3.0, 2009)
    │  Synchronous, blocking, thread-per-request
    │  Maintenance mode since Spring 5
    ▼
WebClient (Spring 5.0, 2017)
    │  Non-blocking, reactive, supports sync & async
    │  Recommended replacement for RestTemplate
    ▼
OpenFeign (Spring Cloud)
    │  Declarative HTTP client (interface + annotations)
    │  Integrates with Spring Cloud ecosystem
    ▼
HTTP Interface Client (Spring 6.1+, 2023)
       Declarative (like Feign) but native Spring
       Uses WebClient or RestClient under the hood
```

#### Detailed Comparison

| Feature | RestTemplate | WebClient | OpenFeign |
|---------|-------------|-----------|-----------|
| **Style** | Imperative | Reactive / Imperative | Declarative |
| **Blocking** | Always blocking | Non-blocking (can block with `.block()`) | Blocking by default |
| **Thread model** | Thread-per-request | Event loop (Netty) | Thread-per-request |
| **Spring Cloud integration** | Manual | Manual | Native (service discovery, LB, CB) |
| **Learning curve** | Low | Medium-High | Very Low |
| **Error handling** | RestClientException | Reactive operators | Feign ErrorDecoder |
| **Interceptors** | ClientHttpRequestInterceptor | ExchangeFilterFunction | RequestInterceptor |
| **Streaming** | No | Yes (SSE, WebSocket) | No |
| **Status** | Maintenance mode | Active development | Active (Spring Cloud) |
| **Best for** | Legacy apps, simple cases | High-throughput reactive apps | Microservice-to-microservice |

#### Code Comparison

**RestTemplate:**

```java
@Service
public class PaymentClient {

    private final RestTemplate restTemplate;

    public PaymentResponse charge(PaymentRequest request) {
        ResponseEntity<PaymentResponse> response = restTemplate.postForEntity(
            "http://payment-service/api/payments",
            request,
            PaymentResponse.class
        );
        return response.getBody();
    }
}
```

**WebClient:**

```java
@Service
public class PaymentClient {

    private final WebClient webClient;

    public Mono<PaymentResponse> charge(PaymentRequest request) {
        return webClient.post()
            .uri("/api/payments")
            .bodyValue(request)
            .retrieve()
            .bodyToMono(PaymentResponse.class);
    }

    // Blocking usage (when not fully reactive)
    public PaymentResponse chargeBlocking(PaymentRequest request) {
        return charge(request).block(Duration.ofSeconds(5));
    }
}
```

**OpenFeign:**

```java
@FeignClient(
    name = "payment-service",
    fallbackFactory = PaymentClientFallbackFactory.class
)
public interface PaymentClient {

    @PostMapping("/api/payments")
    PaymentResponse charge(@RequestBody PaymentRequest request);

    @GetMapping("/api/payments/{id}")
    PaymentResponse getPayment(@PathVariable("id") Long id);
}
```

---

### 2.2 Why Is RestTemplate Deprecated (Maintenance Mode)?

RestTemplate isn't technically deprecated — it's in **maintenance mode**. No new features, only critical bug fixes.

```
Problems with RestTemplate:
─────────────────────────────────────────────────────────

1. Blocking I/O — one thread per request
   ┌─────────────────────────────────────────┐
   │  Thread Pool (200 threads)              │
   │  ┌───┐┌───┐┌───┐┌───┐...┌───┐         │
   │  │ T1 ││ T2 ││ T3 ││ T4 │   │T200│    │
   │  │busy││busy││busy││busy│   │busy│    │
   │  └───┘└───┘└───┘└───┘   └───┘         │
   │  All 200 threads BLOCKED waiting for    │
   │  downstream HTTP responses              │
   │  Thread pool exhausted → 503 errors     │
   └─────────────────────────────────────────┘

2. WebClient alternative — non-blocking
   ┌─────────────────────────────────────────┐
   │  Event Loop (4 threads for 10k+ reqs)   │
   │  ┌───┐┌───┐┌───┐┌───┐                  │
   │  │ T1 ││ T2 ││ T3 ││ T4 │              │
   │  │ ok ││ ok ││ ok ││ ok │              │
   │  └───┘└───┘└───┘└───┘                  │
   │  Threads never block — handle callbacks  │
   │  when responses arrive                   │
   └─────────────────────────────────────────┘
```

**When RestTemplate is still acceptable:**
- Low-traffic internal tools
- Batch jobs with limited concurrency
- Migrating legacy code (not worth rewriting just for this)
- Spring MVC apps that don't need reactive capabilities

**When to use WebClient instead:**
- High concurrency (hundreds/thousands of concurrent outbound calls)
- Reactive stack (Spring WebFlux)
- Streaming responses (SSE, chunked)
- Need non-blocking behavior for performance

---

### 2.3 How Does OpenFeign Work Internally?

```
1. At startup:
   ┌────────────────────────────────────────────────────────┐
   │ @EnableFeignClients scans for @FeignClient interfaces   │
   │       │                                                │
   │       ▼                                                │
   │ FeignClientFactoryBean creates a dynamic proxy          │
   │ for each @FeignClient interface                         │
   │       │                                                │
   │       ▼                                                │
   │ Proxy wires in:                                         │
   │   • Encoder (object → HTTP body)                       │
   │   • Decoder (HTTP body → object)                       │
   │   • Contract (annotation parsing → HTTP metadata)      │
   │   • Client (actual HTTP execution)                     │
   │   • Interceptors (auth headers, logging)               │
   │   • ErrorDecoder (HTTP error → exception)              │
   │   • LoadBalancer (if service discovery enabled)        │
   └────────────────────────────────────────────────────────┘

2. At runtime (when you call paymentClient.charge(request)):
   ┌────────────────────────────────────────────────────────┐
   │ Your Code: paymentClient.charge(request)               │
   │       │                                                │
   │       ▼                                                │
   │ InvocationHandler intercepts method call               │
   │       │                                                │
   │       ▼                                                │
   │ Build RequestTemplate from annotations:                │
   │   POST /api/payments                                   │
   │   Content-Type: application/json                       │
   │   Body: {serialized PaymentRequest}                    │
   │       │                                                │
   │       ▼                                                │
   │ Apply RequestInterceptors (add auth header, trace ID)  │
   │       │                                                │
   │       ▼                                                │
   │ LoadBalancer resolves "payment-service" → 10.0.1.5:8080│
   │       │                                                │
   │       ▼                                                │
   │ HTTP Client executes request                           │
   │       │                                                │
   │       ▼                                                │
   │ Response Decoder converts response → PaymentResponse   │
   │       │                                                │
   │       ▼                                                │
   │ Return to your code                                    │
   └────────────────────────────────────────────────────────┘
```

---

### 2.4 How to Define Feign Clients in Spring Boot

**Step 1: Dependencies**

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-openfeign</artifactId>
</dependency>
```

**Step 2: Enable Feign**

```java
@SpringBootApplication
@EnableFeignClients
public class OrderServiceApplication { }
```

**Step 3: Define Client Interface**

```java
@FeignClient(
    name = "payment-service",                    // Service name (for discovery)
    url = "${payment.service.url:}",             // Direct URL (overrides discovery)
    configuration = PaymentClientConfig.class,   // Custom config
    fallbackFactory = PaymentFallbackFactory.class
)
public interface PaymentClient {

    @PostMapping("/api/v1/payments")
    PaymentResponse charge(@RequestBody PaymentRequest request);

    @GetMapping("/api/v1/payments/{id}")
    PaymentResponse getPayment(@PathVariable("id") Long paymentId);

    @GetMapping("/api/v1/payments")
    List<PaymentResponse> getPayments(
        @RequestParam("userId") Long userId,
        @RequestParam("status") String status
    );
}
```

**Step 4: Custom Configuration**

```java
public class PaymentClientConfig {

    @Bean
    public RequestInterceptor authInterceptor() {
        return template -> {
            String token = SecurityContextHolder.getContext()
                .getAuthentication().getCredentials().toString();
            template.header("Authorization", "Bearer " + token);
        };
    }

    @Bean
    public ErrorDecoder errorDecoder() {
        return (methodKey, response) -> {
            if (response.status() == 404) {
                return new PaymentNotFoundException("Payment not found");
            }
            if (response.status() == 429) {
                return new RateLimitExceededException("Too many requests");
            }
            return new FeignException.InternalServerError(
                "Payment service error", response.request(), null, null);
        };
    }

    @Bean
    public Retryer retryer() {
        return new Retryer.Default(100, 1000, 3);
    }
}
```

---

### 2.5 Handling Timeouts, Retries, and Fallbacks

#### Timeouts

```yaml
# application.yml — Feign timeout configuration
spring:
  cloud:
    openfeign:
      client:
        config:
          default:                           # Global defaults
            connect-timeout: 3000            # 3 seconds to establish connection
            read-timeout: 5000               # 5 seconds to read response
          payment-service:                   # Per-client overrides
            connect-timeout: 2000
            read-timeout: 10000              # Payment takes longer
```

```
Timeout Flow:
─────────────────────────────────────────────────────────────
Client                        Server
  │── TCP SYN ─────────────→│
  │                          │  ← connect-timeout applies here
  │←─ TCP SYN-ACK ──────────│
  │── HTTP Request ─────────→│
  │                          │  ← read-timeout applies here
  │                          │    (server processing)
  │←─ HTTP Response ─────────│
```

**Production values:**

```
Service Type          Connect Timeout    Read Timeout
──────────────────    ───────────────    ────────────
Internal microservice    1-3 seconds       3-5 seconds
Payment gateway          3-5 seconds       10-30 seconds
Batch/Report service     3-5 seconds       30-60 seconds
Health check             1 second           2 seconds
```

#### Retries

```java
// Custom retry configuration
public class PaymentClientConfig {

    @Bean
    public Retryer retryer() {
        // period=100ms, maxPeriod=1000ms, maxAttempts=3
        return new Retryer.Default(100, 1000, 3);
    }
}
```

```
Retry with backoff:
Attempt 1: [REQUEST]──[500 ERROR]── wait 100ms
Attempt 2: [REQUEST]──[500 ERROR]── wait 200ms
Attempt 3: [REQUEST]──[200 OK] ✓

Only retry on:
  ✓ Connection timeouts
  ✓ 5xx server errors
  ✗ 4xx client errors (data problem, won't fix on retry)
  ✗ Non-idempotent operations (POST without idempotency key)
```

#### Fallbacks (with Resilience4j Circuit Breaker)

```java
@Component
public class PaymentFallbackFactory implements FallbackFactory<PaymentClient> {

    @Override
    public PaymentClient create(Throwable cause) {
        return new PaymentClient() {
            @Override
            public PaymentResponse charge(PaymentRequest request) {
                if (cause instanceof CircuitBreakerOpenException) {
                    return PaymentResponse.queued("Circuit open, payment queued for retry");
                }
                return PaymentResponse.failed("Payment service unavailable: " + cause.getMessage());
            }

            @Override
            public PaymentResponse getPayment(Long paymentId) {
                return PaymentResponse.cached(paymentId);
            }

            @Override
            public List<PaymentResponse> getPayments(Long userId, String status) {
                return Collections.emptyList();
            }
        };
    }
}
```

```yaml
# Enable circuit breaker with Feign
spring:
  cloud:
    openfeign:
      circuitbreaker:
        enabled: true
```

---

## 3. Design & Architecture

### 3.1 How Do Services Discover Each Other?

Without discovery, every service needs hardcoded addresses. When instances scale up/down or IPs change, everything breaks.

```
Without Service Discovery:
┌─────────┐  http://10.0.1.5:8080   ┌─────────┐
│  Order  │─────────────────────────→│ Payment │
│ Service │  HARDCODED IP!           │ Service │
└─────────┘                          └─────────┘
    If Payment moves to 10.0.2.8 → Order breaks!

With Service Discovery:
┌─────────┐  "payment-service"  ┌──────────────┐  10.0.1.5:8080  ┌─────────┐
│  Order  │────────────────────→│   Service    │────────────────→│ Payment │
│ Service │  logical name       │   Registry   │  resolved IP    │ Service │
└─────────┘                     │ (Eureka/     │                 │ (inst 1)│
                                │  Consul)     │────────────────→├─────────┤
                                └──────────────┘                 │ Payment │
                                      ▲                          │ (inst 2)│
                                      │  register                └─────────┘
                                      │  + heartbeat
                                 ┌─────────┐
                                 │ Payment │
                                 │ Service │
                                 └─────────┘
```

#### Eureka (Netflix OSS — Spring Cloud Netflix)

```
┌─────────────────────────────────────────────────────┐
│                   Eureka Server                      │
│                                                     │
│  Registry:                                          │
│  ┌─────────────────────────────────────────────┐   │
│  │ payment-service:                             │   │
│  │   - 10.0.1.5:8080 (UP, last heartbeat: 2s)  │   │
│  │   - 10.0.1.6:8080 (UP, last heartbeat: 5s)  │   │
│  │                                              │   │
│  │ inventory-service:                           │   │
│  │   - 10.0.2.3:8081 (UP, last heartbeat: 1s)  │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
│  Self-preservation mode:                            │
│  If > 85% instances stop sending heartbeats,        │
│  assume NETWORK issue, not mass failure.            │
│  Keep all registrations alive.                      │
└─────────────────────────────────────────────────────┘
```

#### Consul (HashiCorp)

```
Consul provides more than discovery:
  ├── Service Discovery (like Eureka)
  ├── Health Checking (HTTP, TCP, script-based)
  ├── KV Store (dynamic configuration)
  ├── Multi-datacenter support (built-in)
  └── Service Mesh (Consul Connect with mTLS)
```

#### Kubernetes-Native Discovery

```
In Kubernetes, discovery is built-in:
┌─────────────────────────────────────────────┐
│  Kubernetes DNS                              │
│                                             │
│  payment-service.default.svc.cluster.local  │
│       │                                     │
│       └──→ ClusterIP → Pod 1 (10.0.1.5)    │
│                      → Pod 2 (10.0.1.6)    │
│                                             │
│  No Eureka/Consul needed!                   │
│  K8s Service does discovery + load balancing│
└─────────────────────────────────────────────┘
```

**Production recommendation:** If running on Kubernetes, use K8s-native service discovery. Eureka/Consul add unnecessary complexity. If running on VMs or bare metal, Consul is preferred over Eureka for its health checking and multi-DC support.

---

### 3.2 Client-Side vs Server-Side Load Balancing

```
Server-Side Load Balancing:
┌─────────┐         ┌──────────────┐         ┌─────────┐
│  Order  │────────→│   Load       │────────→│Payment-1│
│ Service │         │  Balancer    │────────→│Payment-2│
└─────────┘         │ (Nginx/ALB)  │────────→│Payment-3│
                    └──────────────┘
    Client sends to ONE address.
    LB decides which instance.
    Extra network hop.

Client-Side Load Balancing:
┌─────────────────────────────┐
│  Order Service              │         ┌─────────┐
│  ┌────────────────────────┐ │────────→│Payment-1│
│  │ Spring Cloud           │ │         ├─────────┤
│  │ LoadBalancer           │ │────────→│Payment-2│
│  │                        │ │         ├─────────┤
│  │ Knows all instances    │ │────────→│Payment-3│
│  │ from Service Registry  │ │         └─────────┘
│  └────────────────────────┘ │
└─────────────────────────────┘
    Client knows ALL instances.
    Client picks one (round-robin, random, weighted).
    No extra hop — direct connection.
```

| Aspect | Server-Side (Nginx/ALB) | Client-Side (Spring Cloud LB) |
|--------|------------------------|-------------------------------|
| **Extra hop** | Yes — LB is in the path | No — direct to instance |
| **Single point of failure** | LB itself (mitigated with HA) | None — logic in client |
| **Instance awareness** | LB knows health | Client fetches from registry |
| **Customizable** | Limited (L4/L7 rules) | Fully programmable |
| **Scalability** | LB can bottleneck | Scales with clients |
| **Use for** | External traffic, edge | Internal service-to-service |

**Production pattern:** Use **server-side LB** (ALB/Nginx) at the edge (external traffic into your cluster) and **client-side LB** (Spring Cloud LoadBalancer) for internal service-to-service calls.

---

### 3.3 How Does Spring Cloud LoadBalancer Work?

```java
// Spring Cloud LoadBalancer replaces Netflix Ribbon (deprecated)
// It integrates transparently with Feign and WebClient

@FeignClient(name = "payment-service") // Name used for service lookup
public interface PaymentClient {
    @PostMapping("/api/payments")
    PaymentResponse charge(@RequestBody PaymentRequest request);
}
```

```
Request Flow:
─────────────────────────────────────────────────────
1. Feign sees name = "payment-service"
2. Spring Cloud LoadBalancer intercepts
3. Queries ServiceInstanceListSupplier (Eureka/Consul/K8s)
4. Gets list: [10.0.1.5:8080, 10.0.1.6:8080, 10.0.2.3:8080]
5. Applies algorithm (default: Round Robin)
6. Selects: 10.0.1.6:8080
7. Feign sends HTTP to http://10.0.1.6:8080/api/payments
```

**Load balancing algorithms:**

```java
// Custom load balancer — prefer same-zone instances
@Bean
public ServiceInstanceListSupplier discoveryClientServiceInstanceListSupplier(
        ConfigurableApplicationContext context) {
    return ServiceInstanceListSupplier.builder()
        .withDiscoveryClient()
        .withZonePreference()         // Prefer same availability zone
        .withHealthChecks()           // Filter unhealthy instances
        .withCaching()                // Cache instance list
        .build(context);
}
```

---

### 3.4 What Is API Gateway and Why Is It Important?

```
Without API Gateway:
┌─────────┐     ┌──────────────┐
│ Mobile  │────→│ Order Svc    │  Client knows all service addresses
│  App    │────→│ Payment Svc  │  Each service handles auth, rate limiting
│         │────→│ User Svc     │  Cross-cutting concerns duplicated
└─────────┘     └──────────────┘

With API Gateway:
┌─────────┐     ┌──────────────┐     ┌──────────────┐
│ Mobile  │────→│  API Gateway │────→│ Order Svc    │
│  App    │     │              │────→│ Payment Svc  │
│         │     │ (Spring Cloud│────→│ User Svc     │
└─────────┘     │  Gateway)    │     └──────────────┘
                └──────────────┘
                Single entry point.
                All cross-cutting concerns centralized.
```

**What API Gateway handles:**

```
┌──────────────────────────────────────────────────────┐
│                  API Gateway                          │
├──────────────────────────────────────────────────────┤
│                                                      │
│  ┌──────────────────┐  ┌──────────────────────────┐ │
│  │ Authentication   │  │ Rate Limiting             │ │
│  │ (JWT validation) │  │ (100 req/min per user)    │ │
│  └──────────────────┘  └──────────────────────────┘ │
│                                                      │
│  ┌──────────────────┐  ┌──────────────────────────┐ │
│  │ Request Routing  │  │ Response Caching          │ │
│  │ /api/orders → Svc│  │ (GET responses cached)    │ │
│  └──────────────────┘  └──────────────────────────┘ │
│                                                      │
│  ┌──────────────────┐  ┌──────────────────────────┐ │
│  │ Load Balancing   │  │ Circuit Breaking          │ │
│  └──────────────────┘  └──────────────────────────┘ │
│                                                      │
│  ┌──────────────────┐  ┌──────────────────────────┐ │
│  │ Protocol Transl. │  │ Request/Response Transform│ │
│  │ (gRPC → REST)    │  │ (API versioning, mapping) │ │
│  └──────────────────┘  └──────────────────────────┘ │
│                                                      │
│  ┌──────────────────┐  ┌──────────────────────────┐ │
│  │ Logging/Metrics  │  │ CORS Handling             │ │
│  └──────────────────┘  └──────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

**Spring Cloud Gateway example:**

```yaml
spring:
  cloud:
    gateway:
      routes:
        - id: order-service
          uri: lb://order-service
          predicates:
            - Path=/api/orders/**
          filters:
            - StripPrefix=1
            - name: CircuitBreaker
              args:
                name: orderServiceCB
                fallbackUri: forward:/fallback/orders

        - id: payment-service
          uri: lb://payment-service
          predicates:
            - Path=/api/payments/**
          filters:
            - StripPrefix=1
            - name: RequestRateLimiter
              args:
                redis-rate-limiter.replenishRate: 10
                redis-rate-limiter.burstCapacity: 20
```

---

### 3.5 How Does Communication Change with API Gateway?

```
Without Gateway (mesh of direct calls):
┌────────┐   ┌────────┐   ┌────────┐
│ Client │──→│ Svc A  │──→│ Svc B  │
│        │──→│        │   │        │
│        │──→│ Svc C  │──→│ Svc D  │
└────────┘   └────────┘   └────────┘
  Client makes N different calls
  Knows about N services

With Gateway:
┌────────┐      ┌─────────┐      ┌────────┐   ┌────────┐
│ Client │─────→│ Gateway │─────→│ Svc A  │──→│ Svc B  │
│        │      │         │─────→│ Svc C  │──→│ Svc D  │
└────────┘      └─────────┘      └────────┘   └────────┘
  Client makes 1-2 calls to Gateway
  Gateway fans out internally
  Internal services still talk directly to each other
```

**Key insight:** The API Gateway sits between **external clients** and your services. **Internal** service-to-service communication (Service A → Service B) still happens directly via Feign/gRPC — it does NOT go through the gateway. Gateway is for north-south traffic, not east-west.

---

## 4. Synchronous vs Asynchronous

### 4.1 Deep Dive: Sync vs Async Communication

```
Synchronous Chain (worst case):
─────────────────────────────────────────────────────────────
User Request → Order(50ms) → Payment(200ms) → Inventory(100ms) → Notification(80ms)
                                                                     │
Total latency = 50 + 200 + 100 + 80 = 430ms (ADDITIVE)              │
If any fails → entire chain fails                                    │

Asynchronous (event-driven):
─────────────────────────────────────────────────────────────
User Request → Order(50ms) → publish event → Response to User (50ms!)
                                    │
                    ┌───────────────┼────────────────┐
                    ▼               ▼                ▼
               Payment         Inventory        Notification
              (processes       (processes        (processes
               async)           async)            async)

Total latency = 50ms (just the order creation)
Each consumer processes independently
```

---

### 4.2 When NOT to Use Synchronous Calls

```
Avoid sync when:
─────────────────────────────────────────────────────────────

1. Chain depth > 3 services
   A → B → C → D → E
   Latency: sum of all
   Failure probability: multiplied (0.99^5 = 0.95 = 5% failure rate!)

2. Downstream can be slow/unreliable
   Your service: 50ms p99
   Payment gateway: 2-30 seconds
   Don't block your thread for 30 seconds

3. Result is not needed immediately
   "Send confirmation email"
   "Update analytics dashboard"
   "Sync to search index"

4. Fan-out to many services
   Order created → notify 8 services
   8 sync calls = 8x latency risk

5. Spiky traffic patterns
   Black Friday: 100x normal traffic
   Sync calls → all downstream services must scale simultaneously
   Async → message queue absorbs the spike
```

---

### 4.3 Problems with Synchronous Chaining

#### Latency Amplification

```
Individual service SLAs:
  Order:     p99 = 50ms
  Payment:   p99 = 200ms
  Inventory: p99 = 100ms
  Shipping:  p99 = 150ms

Chain latency:
  p99 = 50 + 200 + 100 + 150 = 500ms

But it's WORSE than that:
  p99 of chain ≈ p99.75 of individual services
  (because ANY slow response makes the chain slow)

  Real p99 of chain ≈ 800-1200ms (tail latency amplification)
```

#### Cascading Failures

```
Normal:
  Order → Payment → Inventory → Shipping  (all healthy)

Shipping goes slow (30s timeout):
  ┌───────┐     ┌─────────┐     ┌──────────┐     ┌──────────┐
  │ Order │────→│ Payment │────→│Inventory │────→│ Shipping │
  │ 200   │     │ 200     │     │ 200      │     │ SLOW     │
  │threads│     │ threads │     │ threads  │     │ 30s resp │
  │BLOCKED│     │ BLOCKED │     │ BLOCKED  │     │          │
  └───────┘     └─────────┘     └──────────┘     └──────────┘
  All upstream services run out of threads.
  One slow service takes down the entire system.
  This is the #1 production outage pattern in microservices.
```

---

### 4.4 What Is Event-Driven Architecture?

```
┌──────────────────────────────────────────────────────────────┐
│                Event-Driven Architecture                      │
│                                                              │
│  Producers                  Event Backbone              Consumers
│  ┌─────────┐               ┌───────────┐              ┌──────────┐
│  │ Order   │──OrderCreated─→│           │─────────────→│ Payment  │
│  │ Service │               │           │              │ Service  │
│  └─────────┘               │   Kafka   │              └──────────┘
│                            │   /       │              ┌──────────┐
│  ┌─────────┐               │  RabbitMQ │─────────────→│Inventory │
│  │ Payment │──PaymentDone──→│           │              │ Service  │
│  │ Service │               │           │              └──────────┘
│  └─────────┘               │           │              ┌──────────┐
│                            │           │─────────────→│Analytics │
│                            └───────────┘              │ Service  │
│                                                       └──────────┘
│                                                              │
│  Key properties:                                             │
│  • Producers don't know consumers (loose coupling)           │
│  • Events are facts (immutable records of what happened)     │
│  • Consumers process at their own pace                       │
│  • New consumers can be added without changing producers     │
│  • Events can be replayed for recovery or new consumers      │
└──────────────────────────────────────────────────────────────┘
```

**Event types:**

```
1. Event Notification (thin event)
   { "type": "ORDER_CREATED", "orderId": "123" }
   Consumer calls back to get details.

2. Event-Carried State Transfer (fat event)
   { "type": "ORDER_CREATED", "orderId": "123",
     "items": [...], "totalAmount": 500, "userId": "42" }
   Consumer has all data. No callback needed.

3. Domain Event (DDD)
   Represents something meaningful that happened in the domain.
   Bounded context publishes, other contexts subscribe.
```

---

### 4.5 Apache Kafka vs RabbitMQ

```
Kafka Architecture:
┌─────────────────────────────────────────────────────────────┐
│  Topic: order-events                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Partition 0   │  │ Partition 1   │  │ Partition 2   │     │
│  │ [msg1][msg3]  │  │ [msg2][msg5]  │  │ [msg4][msg6]  │    │
│  │ [msg7]...     │  │ [msg8]...     │  │ [msg9]...     │    │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │                │
│   Consumer A-0      Consumer A-1      Consumer A-2           │
│   (Consumer Group A)                                         │
│                                                              │
│   Consumer B-0 reads ALL partitions                          │
│   (Consumer Group B — independent)                           │
│                                                              │
│  Messages are RETAINED (days/weeks), not deleted on consume  │
│  Consumers track their own offset                            │
│  REPLAY possible: reset offset to re-process                 │
└─────────────────────────────────────────────────────────────┘

RabbitMQ Architecture:
┌─────────────────────────────────────────────────────────────┐
│  Exchange: order-exchange (topic/direct/fanout)              │
│       │                                                     │
│       ├── routing key: order.created ──→ Queue: payment-q   │
│       │                                   │                 │
│       │                              Consumer 1 (Payment)   │
│       │                                                     │
│       ├── routing key: order.* ──────→ Queue: inventory-q   │
│       │                                   │                 │
│       │                              Consumer 2 (Inventory) │
│       │                                                     │
│       └── routing key: # ────────────→ Queue: audit-q       │
│                                           │                 │
│                                      Consumer 3 (Audit)     │
│                                                              │
│  Messages DELETED once acknowledged by consumer              │
│  Smart routing via exchanges + bindings                      │
│  No replay capability (once consumed, gone)                  │
└─────────────────────────────────────────────────────────────┘
```

| Feature | Apache Kafka | RabbitMQ |
|---------|-------------|----------|
| **Model** | Distributed log (append-only) | Message broker (queue-based) |
| **Ordering** | Per partition | Per queue |
| **Retention** | Time/size-based (days/weeks) | Until consumed + acknowledged |
| **Replay** | Yes (reset consumer offset) | No (message deleted after ack) |
| **Throughput** | Very high (millions/sec) | High (tens of thousands/sec) |
| **Latency** | Medium (batching) | Low (per-message delivery) |
| **Routing** | Topic + partition key | Flexible (direct, topic, fanout, headers) |
| **Consumer model** | Pull (consumer polls) | Push (broker pushes to consumer) |
| **Scaling** | Horizontal (add partitions) | Vertical + clustering |
| **Delivery** | At-least-once / exactly-once | At-least-once / at-most-once |
| **Protocol** | Custom binary protocol | AMQP (standard) |
| **Built-in features** | Stream processing (Kafka Streams) | Priority queues, DLX, TTL |

---

### 4.6 When to Use Kafka vs RabbitMQ

```
Choose KAFKA when:                        Choose RABBITMQ when:
──────────────────                        ────────────────────
✓ High-throughput event streaming          ✓ Complex routing logic needed
✓ Event sourcing / audit log               ✓ Low-latency per-message delivery
✓ Need to replay events                    ✓ Priority queues required
✓ Multiple consumer groups for same data   ✓ RPC / request-reply pattern
✓ Data pipeline (ETL, analytics)           ✓ Simpler operational model
✓ Ordering within a partition key          ✓ Message TTL and dead-lettering
✓ Long-term event storage                  ✓ Existing AMQP ecosystem
✓ Stream processing (Kafka Streams)        ✓ Smaller scale (< 100k msg/sec)
```

**Production decision:**

| Scenario | Recommendation |
|----------|---------------|
| Order events consumed by 5 different services | Kafka (fan-out to consumer groups) |
| Task queue: process image uploads | RabbitMQ (work queue with ack) |
| Real-time analytics pipeline | Kafka (high throughput + stream processing) |
| Delayed message delivery (send email in 30 min) | RabbitMQ (delayed message plugin) |
| CDC (database change data capture) | Kafka (Debezium + Kafka Connect) |
| Chat message routing | RabbitMQ (low latency, flexible routing) |
| Event sourcing (replay full history) | Kafka (immutable log, retention) |

---

## 5. Fault Tolerance & Resilience

### 5.1 Circuit Breaker Pattern

Prevents cascading failures by stopping calls to a failing service.

```
States:
─────────────────────────────────────────────────────────────

CLOSED (normal operation):
  All requests go through.
  Track failure rate.
  [req✓][req✓][req✗][req✓][req✗][req✓]
  Failure rate: 33% (below threshold 50%)

  If failure rate >= threshold:
      │
      ▼

OPEN (circuit tripped):
  All requests FAIL IMMEDIATELY (no network call).
  Return fallback response.
  [req→fallback][req→fallback][req→fallback]
  Wait for cooldown period (e.g., 30 seconds).

  After cooldown:
      │
      ▼

HALF-OPEN (testing recovery):
  Allow LIMITED requests through.
  [req✓][req✓][req✗]
  If success rate good → CLOSED
  If failures continue → OPEN (reset cooldown)
```

```
                   failure rate >= threshold
         ┌──────────────────────────────────────┐
         │                                      ▼
    ┌────────┐                            ┌──────────┐
    │ CLOSED │                            │   OPEN   │
    │(normal)│                            │(fail fast│
    └────────┘                            │ fallback)│
         ▲                                └────┬─────┘
         │                                     │
         │   success rate good                 │ wait timeout
         │                                     │
         └──────────────┐                      │
                        │                      ▼
                   ┌─────────┐
                   │HALF-OPEN│
                   │ (probe) │
                   └─────────┘
                        │
                        │ failures continue
                        └──────────────────→ OPEN
```

---

### 5.2 Implementing Circuit Breaker with Resilience4j

```xml
<dependency>
    <groupId>io.github.resilience4j</groupId>
    <artifactId>resilience4j-spring-boot3</artifactId>
</dependency>
```

```java
@Service
public class PaymentService {

    @CircuitBreaker(name = "paymentService", fallbackMethod = "chargeFallback")
    @Retry(name = "paymentService")
    @Bulkhead(name = "paymentService")
    @TimeLimiter(name = "paymentService")
    public PaymentResponse charge(PaymentRequest request) {
        return paymentClient.charge(request);
    }

    private PaymentResponse chargeFallback(PaymentRequest request, Throwable t) {
        if (t instanceof CallNotPermittedException) {
            // Circuit is OPEN — don't even try
            return PaymentResponse.circuitOpen("Payment service temporarily unavailable");
        }
        return PaymentResponse.failed("Payment failed: " + t.getMessage());
    }
}
```

```yaml
# application.yml
resilience4j:
  circuitbreaker:
    instances:
      paymentService:
        sliding-window-size: 10           # Evaluate last 10 calls
        sliding-window-type: COUNT_BASED  # or TIME_BASED
        failure-rate-threshold: 50        # Open if 50% fail
        slow-call-rate-threshold: 80      # Open if 80% are slow
        slow-call-duration-threshold: 3s  # "Slow" = > 3 seconds
        wait-duration-in-open-state: 30s  # Stay open for 30s
        permitted-number-of-calls-in-half-open-state: 5
        minimum-number-of-calls: 5        # Need 5 calls before evaluating
        record-exceptions:
          - java.io.IOException
          - java.util.concurrent.TimeoutException
          - feign.FeignException.InternalServerError
        ignore-exceptions:
          - com.example.BusinessValidationException
```

---

### 5.3 Retry Pattern

```
Retry with Exponential Backoff + Jitter:
─────────────────────────────────────────────────────────────
Attempt 1: [REQUEST]──[FAIL]── wait 100ms (+ jitter 0-50ms)
Attempt 2: [REQUEST]──[FAIL]── wait 200ms (+ jitter 0-100ms)
Attempt 3: [REQUEST]──[SUCCESS] ✓

Without jitter (Thundering Herd):
  All 1000 clients retry at exactly 100ms → server hit with 1000 simultaneous requests

With jitter (spread retries):
  Clients retry between 100-150ms → gradual ramp, server recovers
```

```yaml
resilience4j:
  retry:
    instances:
      paymentService:
        max-attempts: 3
        wait-duration: 500ms
        enable-exponential-backoff: true
        exponential-backoff-multiplier: 2
        retry-exceptions:
          - java.io.IOException
          - java.util.concurrent.TimeoutException
        ignore-exceptions:
          - com.example.BusinessValidationException
```

**Critical rule: Retry only idempotent operations.** A GET is safe to retry. A POST that creates a payment is dangerous unless you have an idempotency key.

---

### 5.4 Bulkhead Pattern

Isolates failures by limiting concurrent calls to a downstream service. Prevents one slow service from consuming all threads.

```
Without Bulkhead:
┌──────────────────────────────────────────────┐
│  Shared Thread Pool (200 threads)            │
│                                              │
│  Payment calls: 180 threads (slow!)          │
│  Inventory calls: 15 threads (blocked)       │
│  User calls: 5 threads (starved!)            │
│                                              │
│  Slow payment service consumes ALL threads   │
│  Even healthy services can't get threads      │
└──────────────────────────────────────────────┘

With Bulkhead:
┌──────────────────────────────────────────────┐
│  Payment Pool:    [50 threads max]           │
│  Inventory Pool:  [30 threads max]           │
│  User Pool:       [20 threads max]           │
│  Shared Pool:     [100 threads]              │
│                                              │
│  Payment is slow? Only 50 threads affected.  │
│  Inventory and User continue normally.        │
└──────────────────────────────────────────────┘
```

```yaml
resilience4j:
  bulkhead:
    instances:
      paymentService:
        max-concurrent-calls: 25            # Max 25 concurrent calls
        max-wait-duration: 500ms            # Wait up to 500ms for a slot
  thread-pool-bulkhead:
    instances:
      paymentService:
        max-thread-pool-size: 10
        core-thread-pool-size: 5
        queue-capacity: 20
```

---

### 5.5 Timeout Handling

```yaml
resilience4j:
  timelimiter:
    instances:
      paymentService:
        timeout-duration: 3s               # Cancel if no response in 3s
        cancel-running-future: true        # Actually cancel the thread
```

**Layered timeout strategy (defense in depth):**

```
┌───────────────────────────────────────────────────────────┐
│ Layer 1: Resilience4j TimeLimiter      (3 seconds)        │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ Layer 2: Feign read-timeout          (5 seconds)    │  │
│  │  ┌───────────────────────────────────────────────┐  │  │
│  │  │ Layer 3: HikariCP connection-timeout (10 sec) │  │  │
│  │  │  ┌─────────────────────────────────────────┐  │  │  │
│  │  │  │ Layer 4: DB query timeout   (30 sec)    │  │  │  │
│  │  │  └─────────────────────────────────────────┘  │  │  │
│  │  └───────────────────────────────────────────────┘  │  │
│  └─────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────┘

Innermost timeout should be LONGEST.
Outermost should be SHORTEST.
This way, outer layers fail fast before inner ones.
```

---

### 5.6 Fallback Mechanisms

```java
@Service
public class OrderService {

    @CircuitBreaker(name = "inventoryService", fallbackMethod = "inventoryFallback")
    public InventoryStatus checkInventory(Long productId) {
        return inventoryClient.checkStock(productId);
    }

    // Fallback strategies (in order of preference):
    private InventoryStatus inventoryFallback(Long productId, Throwable t) {

        // Strategy 1: Return cached data
        InventoryStatus cached = cacheService.getInventory(productId);
        if (cached != null && !cached.isStale(Duration.ofMinutes(5))) {
            return cached;
        }

        // Strategy 2: Return degraded response
        return InventoryStatus.unknown(productId, "Inventory check unavailable, assuming in stock");

        // Strategy 3: Queue for later processing (for writes)
        // pendingQueue.add(new InventoryCheckRequest(productId));
        // return InventoryStatus.pending();

        // Strategy 4: Fail fast with meaningful error (last resort)
        // throw new ServiceUnavailableException("Inventory service down");
    }
}
```

**Fallback strategy selection:**

| Scenario | Fallback Strategy | Example |
|----------|------------------|---------|
| **Read operation** | Return cached data | Product catalog, user profile |
| **Non-critical write** | Queue for retry | Analytics events, audit logs |
| **Critical write** | Return pending status + reconcile | Payment (queue and confirm later) |
| **Validation** | Default to permissive (or restrictive depending on risk) | Fraud check: block if unsure |
| **No fallback possible** | Fail fast with clear error | Authentication service down |

---

### 5.7 Resilience4j Decoration Order

The order in which patterns are applied matters.

```
Incoming Request
      │
      ▼
┌─────────────────┐
│    Retry         │  ← Outermost: retries the entire decorated call
│  ┌─────────────┐│
│  │Circuit Breaker│  ← Checks if circuit is open before proceeding
│  │┌───────────┐││
│  ││ Rate Limit │││  ← Controls call rate
│  ││┌─────────┐│││
│  │││Bulkhead  ││││  ← Limits concurrent calls
│  │││┌───────┐ ││││
│  ││││Timeout │ ││││  ← Innermost: times out the actual call
│  ││││  Call  │ ││││
│  │││└───────┘ ││││
│  ││└─────────┘│││
│  │└───────────┘││
│  └─────────────┘│
└─────────────────┘

Retry → CircuitBreaker → RateLimiter → Bulkhead → TimeLimiter → Call
```

---

## 6. Real Production Problems

### 6.1 What Happens If Service B Is Down When Service A Calls It?

```
Scenario: Order Service calls Payment Service, Payment is DOWN
─────────────────────────────────────────────────────────────

Without resilience:
  Order → Payment (DOWN) → ConnectionRefusedException → 500 to user
  Every request fails for every user.
  Thread blocked for connect-timeout duration.

With full resilience stack:
  Order → CircuitBreaker check
     │
     ├── Circuit CLOSED: try Payment
     │     └── Timeout (3s) → Retry (2 more times)
     │           └── All retries fail → Circuit records failure
     │                 └── Failure rate > 50% → Circuit OPENS
     │
     ├── Circuit OPEN: skip Payment entirely
     │     └── Return fallback immediately (< 1ms)
     │           "Payment queued for processing"
     │
     └── Circuit HALF-OPEN (after 30s): probe Payment
           └── Success → Circuit CLOSES, resume normal
           └── Fail → Circuit stays OPEN, wait another 30s
```

**Production response plan:**

```java
@CircuitBreaker(name = "payment", fallbackMethod = "paymentDown")
public PaymentResult charge(PaymentRequest req) {
    return paymentClient.charge(req);
}

private PaymentResult paymentDown(PaymentRequest req, Throwable t) {
    // 1. Save to pending_payments table
    pendingPaymentRepo.save(new PendingPayment(req, Instant.now()));

    // 2. Publish event for async retry
    eventPublisher.publish(new PaymentPendingEvent(req));

    // 3. Return graceful response to user
    return PaymentResult.pending("Payment will be processed shortly");
}

// Background job retries pending payments when service recovers
@Scheduled(fixedDelay = 60000)
public void retryPendingPayments() {
    List<PendingPayment> pending = pendingPaymentRepo.findUnprocessed();
    for (PendingPayment p : pending) {
        try {
            paymentClient.charge(p.toRequest());
            p.markProcessed();
        } catch (Exception e) {
            p.incrementRetryCount();
            if (p.getRetryCount() > MAX_RETRIES) {
                p.markFailed();
                alertService.notifyOps("Payment permanently failed: " + p.getId());
            }
        }
    }
}
```

---

### 6.2 How to Prevent Cascading Failures

```
The Cascade:
─────────────────────────────────────────────────────────────
  Shipping (slow)
     ↑ blocks
  Inventory (threads exhausted)
     ↑ blocks
  Payment (threads exhausted)
     ↑ blocks
  Order (threads exhausted) → ALL users get 503!

Prevention stack (defense in depth):
─────────────────────────────────────────────────────────────

Layer 1: TIMEOUTS — Don't wait forever
  ├── Connection timeout: 2s
  ├── Read timeout: 5s
  └── Kill slow calls fast

Layer 2: CIRCUIT BREAKER — Stop calling broken services
  ├── Fail fast when downstream is unhealthy
  └── Auto-recover when downstream recovers

Layer 3: BULKHEAD — Isolate blast radius
  ├── Limit concurrent calls per downstream service
  └── Slow service can't steal threads from others

Layer 4: RETRY with BACKOFF + JITTER — Recover from transient failures
  ├── Don't overwhelm recovering service
  └── Jitter prevents thundering herd

Layer 5: FALLBACK — Graceful degradation
  ├── Serve stale data from cache
  ├── Return default/degraded response
  └── Queue for async processing

Layer 6: ASYNC DECOUPLING — Remove temporal coupling
  ├── Use events instead of sync calls where possible
  └── Message queue absorbs spikes
```

---

### 6.3 What Is Backpressure and How to Handle It?

**Backpressure** is when a consumer cannot keep up with the rate of incoming data from the producer.

```
Without backpressure handling:
Producer: [msg][msg][msg][msg][msg][msg][msg][msg][msg][msg]...
                                    │
                                    ▼
Consumer: [msg]...[processing]...[msg]...[processing]...
           Slow consumer.
           Messages pile up.
           Memory exhaustion → OOM crash!

With backpressure handling:
Producer: [msg][msg][SLOW DOWN][msg][msg][SLOW DOWN]...
                        │                    │
                   Consumer signals          Consumer signals
                   "I'm behind!"             "Still catching up"
```

**Backpressure strategies:**

| Strategy | How It Works | When to Use |
|----------|-------------|-------------|
| **Drop** | Discard excess messages | Metrics, telemetry (latest value matters) |
| **Buffer** | Queue messages in memory/disk | Short bursts, predictable load |
| **Throttle producer** | Producer slows down | When producer can be controlled |
| **Scale consumer** | Add more consumer instances | Sustained high volume |
| **Sample** | Process every Nth message | Analytics, monitoring |
| **Batch** | Process messages in groups | Database writes, bulk APIs |

**Kafka-specific backpressure:**

```yaml
# Consumer configuration for backpressure
spring:
  kafka:
    consumer:
      max-poll-records: 100          # Process 100 records per poll (not 500 default)
      max-poll-interval-ms: 300000   # 5 min to process a batch before rebalance
      fetch-min-bytes: 1024          # Wait for at least 1KB of data
      fetch-max-wait-ms: 500         # Or max 500ms

    # If consumer is truly too slow:
    # 1. Increase partitions + add consumer instances
    # 2. Batch database writes
    # 3. Use async processing within consumer
```

---

### 6.4 Retry Storm Problem

```
Normal state: 1000 requests/sec to Payment Service
Payment goes DOWN for 30 seconds.

Without protection:
  1000 req/sec × 3 retries = 3000 req/sec hitting Payment
  Payment comes back UP
  Backlog: 30 × 1000 = 30,000 queued retries
  All retries fire simultaneously = 30,000 req SPIKE
  Payment goes down AGAIN under the load
  This is a RETRY STORM → infinite loop of failure

┌─────────────────────────────────────────────────────────┐
│  Payment Service                                        │
│  Capacity: 1000 req/sec                                 │
│                                                         │
│  Normal:  ████████████ 1000/sec (OK)                    │
│  Outage:  starts at T=0                                 │
│  T=30s:   Recovery attempt                              │
│  Retries: ████████████████████████████ 30,000 req burst │
│           (Payment dies again!)                          │
│  T=60s:   Recovery attempt (even more retries queued)   │
│           Exponential amplification!                     │
└─────────────────────────────────────────────────────────┘
```

**Solutions:**

```
1. Circuit Breaker → stops retries when service is confirmed down
2. Exponential backoff → spaces out retries
3. Jitter → prevents synchronized retry waves
4. Max retry count → finite retries, then fail
5. Retry budget → "only 10% of requests can be retries"
6. Queue-based retry → buffer retries in a queue, process gradually
```

```java
// Retry budget example
public class RetryBudgetFilter {
    private final AtomicInteger totalRequests = new AtomicInteger(0);
    private final AtomicInteger retryRequests = new AtomicInteger(0);

    private static final double MAX_RETRY_RATIO = 0.1; // 10% retry budget

    public boolean shouldRetry() {
        double retryRatio = (double) retryRequests.get() / totalRequests.get();
        return retryRatio < MAX_RETRY_RATIO;
    }
}
```

---

### 6.5 Avoiding Duplicate Requests in Retries

When a request times out but actually succeeded on the server, the retry creates a duplicate.

```
Client: POST /payments (idempotency-key: pay-123)
   │── Request sent ──→ Server processes ✓ (but response lost)
   │── TIMEOUT!
   │── Retry: POST /payments (idempotency-key: pay-123)
   │── Server sees pay-123 already processed → returns original result
   │── Client gets response ✓ (no duplicate payment!)
```

```java
@RestController
public class PaymentController {

    @PostMapping("/api/payments")
    public ResponseEntity<PaymentResponse> createPayment(
            @RequestHeader("Idempotency-Key") String idempotencyKey,
            @RequestBody PaymentRequest request) {

        // Check if already processed
        Optional<PaymentResponse> existing = idempotencyStore.get(idempotencyKey);
        if (existing.isPresent()) {
            return ResponseEntity.ok(existing.get());
        }

        // Process and store result atomically
        PaymentResponse response = paymentService.charge(request);
        idempotencyStore.put(idempotencyKey, response, Duration.ofHours(24));

        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }
}
```

---

## 7. Data Consistency & Transactions

### 7.1 Why Are Distributed Transactions Hard?

```
Monolith: Single database, single transaction
┌──────────────────────────────────────┐
│  @Transactional                      │
│  order.save() + payment.save()       │
│  + inventory.update()                │
│  → ONE commit or ONE rollback        │
└──────────────────────────────────────┘

Microservices: Three databases, three transactions
┌──────────┐    ┌──────────┐    ┌──────────┐
│ Order DB │    │Payment DB│    │Invent DB │
│  TX-1 ✓  │    │  TX-2 ✓  │    │  TX-3 ✗  │
└──────────┘    └──────────┘    └──────────┘
  Committed       Committed      FAILED!
  Order exists.   Payment charged. Inventory NOT reserved.
  INCONSISTENT STATE!
  No way to "undo" TX-1 and TX-2 atomically.
```

**The fundamental problems:**

| Problem | Description |
|---------|-------------|
| **No shared transaction** | Each DB has its own transaction manager |
| **Network unreliability** | Commit message can be lost between services |
| **Partial failure** | Some services succeed, others fail |
| **No global rollback** | Can't undo a committed transaction in another service |
| **Timing** | Services process at different speeds |

---

### 7.2 Saga Pattern — Choreography vs Orchestration

Already introduced in the Transaction Guide. Here's the complete comparison with code.

#### Choreography (Event-Driven)

```
┌───────────┐  OrderCreated   ┌────────────┐  PaymentDone  ┌────────────┐
│   Order   │ ──────────────→ │  Payment   │ ────────────→ │ Inventory  │
│  Service  │                 │  Service   │               │  Service   │
└─────┬─────┘                 └─────┬──────┘               └─────┬──────┘
      │                             │                            │
      │  InventoryFailed            │  InventoryFailed           │
      │◄────────────────────────────│◄───────────────────────────│
      │                             │                            │
  Set order                    Refund                       (already failed)
  CANCELLED                    payment
```

```java
// Order Service — publishes event
@Transactional
public Order createOrder(OrderRequest request) {
    Order order = orderRepo.save(new Order(request, OrderStatus.PENDING));
    outboxRepo.save(new OutboxEvent("order.created", order.toEvent()));
    return order;
}

// Payment Service — listens and reacts
@KafkaListener(topics = "order-events")
@Transactional
public void handleOrderCreated(OrderCreatedEvent event) {
    Payment payment = paymentService.charge(event);
    outboxRepo.save(new OutboxEvent("payment.completed", payment.toEvent()));
}

// Payment Service — compensates on inventory failure
@KafkaListener(topics = "inventory-events")
@Transactional
public void handleInventoryFailed(InventoryFailedEvent event) {
    paymentService.refund(event.getPaymentId());
    outboxRepo.save(new OutboxEvent("payment.refunded", event.getOrderId()));
}
```

#### Orchestration (Central Coordinator)

```java
@Service
public class OrderSagaOrchestrator {

    @Transactional
    public OrderResult executeOrderSaga(OrderRequest request) {
        SagaExecution saga = sagaRepo.save(new SagaExecution(request));

        try {
            // Step 1: Create order
            Order order = orderService.createOrder(request);
            saga.completeStep("CREATE_ORDER", order.getId());

            // Step 2: Process payment
            PaymentResult payment = paymentService.charge(request.getPaymentInfo());
            saga.completeStep("CHARGE_PAYMENT", payment.getId());

            // Step 3: Reserve inventory
            inventoryService.reserve(request.getItems());
            saga.completeStep("RESERVE_INVENTORY");

            saga.markCompleted();
            return OrderResult.success(order);

        } catch (Exception e) {
            compensate(saga);
            saga.markFailed(e.getMessage());
            return OrderResult.failed(e.getMessage());
        }
    }

    private void compensate(SagaExecution saga) {
        List<SagaStep> completedSteps = saga.getCompletedSteps();
        Collections.reverse(completedSteps);

        for (SagaStep step : completedSteps) {
            try {
                switch (step.getName()) {
                    case "RESERVE_INVENTORY" -> inventoryService.release(step.getResourceId());
                    case "CHARGE_PAYMENT" -> paymentService.refund(step.getResourceId());
                    case "CREATE_ORDER" -> orderService.cancelOrder(step.getResourceId());
                }
                step.markCompensated();
            } catch (Exception e) {
                step.markCompensationFailed(e.getMessage());
                alertService.notifyOps("Compensation failed for saga: " + saga.getId());
            }
        }
    }
}
```

| Aspect | Choreography | Orchestration |
|--------|-------------|---------------|
| **Coupling** | Very loose | Tighter (orchestrator knows steps) |
| **Visibility** | Hard to see full flow | Full flow visible in orchestrator |
| **Debugging** | Distributed logs, hard | Central saga state, easier |
| **Adding steps** | Add consumer, no code change | Modify orchestrator |
| **Failure handling** | Each service handles its own | Centralized compensation |
| **Best for** | Simple flows (2-3 services) | Complex flows (4+ services) |

---

### 7.3 How Do Services Maintain Data Consistency?

```
┌──────────────────────────────────────────────────────────┐
│           Data Consistency Patterns                       │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  1. Outbox Pattern (guaranteed event publishing)         │
│     @Transactional                                       │
│     save(order) + save(outboxEvent) → single local TX    │
│     CDC/Poller reads outbox → publishes to Kafka         │
│                                                          │
│  2. Event Sourcing (events as source of truth)           │
│     Don't store state, store events that produce state   │
│     OrderCreated → ItemAdded → PaymentReceived → ...     │
│     Replay events to rebuild current state               │
│                                                          │
│  3. CQRS (separate read/write models)                    │
│     Write model: normalized, event sourced               │
│     Read model: denormalized, optimized for queries      │
│     Sync via events (eventually consistent)              │
│                                                          │
│  4. Idempotent Consumer (handle duplicates)              │
│     Store processed event IDs                            │
│     If already processed → skip                          │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

---

## 8. Security in Inter-Service Communication

### 8.1 How Do Services Authenticate Each Other?

```
┌────────────────────────────────────────────────────────────┐
│           Service-to-Service Authentication                 │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Pattern 1: JWT Token Propagation                          │
│  ┌──────┐  JWT  ┌─────────┐  JWT  ┌─────────┐            │
│  │Client│──────→│ Gateway │──────→│ Order   │             │
│  └──────┘       └─────────┘       │ Service │             │
│                                   └────┬────┘             │
│                                        │ Same JWT          │
│                                        ▼ (or new internal) │
│                                   ┌─────────┐             │
│                                   │ Payment │             │
│                                   │ Service │             │
│                                   └─────────┘             │
│                                                            │
│  Pattern 2: Service Account (Client Credentials)           │
│  ┌─────────┐  client_id/secret  ┌──────┐  JWT  ┌───────┐ │
│  │ Order   │──────────────────→│ Auth │──────→│ Order  │ │
│  │ Service │                    │Server│       │Service │ │
│  └─────────┘                    └──────┘       └───┬───┘ │
│                                                     │      │
│      service-to-service JWT (machine identity)      ▼      │
│                                               ┌─────────┐  │
│                                               │ Payment │  │
│                                               │ Service │  │
│                                               └─────────┘  │
│                                                            │
│  Pattern 3: mTLS (Mutual TLS)                              │
│  Both client and server present certificates               │
│  Zero-trust network — every call authenticated             │
│  Managed by service mesh (Istio, Linkerd)                  │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

**JWT propagation in Feign:**

```java
public class JwtRequestInterceptor implements RequestInterceptor {

    @Override
    public void apply(RequestTemplate template) {
        ServletRequestAttributes attrs =
            (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();

        if (attrs != null) {
            String token = attrs.getRequest().getHeader("Authorization");
            if (token != null) {
                template.header("Authorization", token);
            }
        }
    }
}
```

---

### 8.2 What Is mTLS (Mutual TLS)?

```
Standard TLS (one-way):
  Client ──────→ Server
  Client verifies server certificate.
  Server doesn't verify client.

Mutual TLS (two-way):
  Client ←─────→ Server
  Client verifies server certificate. ✓
  Server verifies client certificate. ✓
  Both identities confirmed.
```

```
mTLS Handshake:
─────────────────────────────────────────────────────────────
1. Client → Server: "Hello, supported ciphers"
2. Server → Client: Server certificate + "send YOUR cert"
3. Client → Server: Client certificate
4. Server validates client cert against trusted CA
5. Client validates server cert against trusted CA
6. Both verified → encrypted channel established
```

**In practice with Istio (service mesh):**

```
┌─────────────────────┐          ┌─────────────────────┐
│  Order Service Pod   │          │  Payment Service Pod │
│  ┌─────┐  ┌──────┐  │  mTLS   │  ┌──────┐  ┌─────┐ │
│  │ App │→ │Envoy │──┼─────────┼─→│Envoy │→ │ App │ │
│  │     │  │Proxy │  │encrypted│  │Proxy │  │     │ │
│  └─────┘  └──────┘  │         │  └──────┘  └─────┘ │
└─────────────────────┘          └─────────────────────┘
  Your app code doesn't handle TLS.
  Envoy sidecar handles mTLS transparently.
  Certificates auto-rotated by Istio CA.
```

---

### 8.3 Securing Internal APIs

| Layer | Mechanism | Purpose |
|-------|-----------|---------|
| **Transport** | mTLS / TLS | Encrypt data in transit |
| **Authentication** | JWT / Service accounts | Verify caller identity |
| **Authorization** | RBAC / Scopes | Control what caller can do |
| **Network** | Network policies (K8s) | Restrict which pods can communicate |
| **API Gateway** | Rate limiting, WAF | Protect against abuse |
| **Secrets** | Vault / K8s Secrets | Secure credential storage |

```yaml
# Kubernetes NetworkPolicy — only allow Order → Payment
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-order-to-payment
spec:
  podSelector:
    matchLabels:
      app: payment-service
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: order-service
      ports:
        - port: 8080
```

---

## 9. gRPC vs REST

### 9.1 What Is gRPC and How Is It Different from REST?

```
REST:
  Client: POST /api/payments  {"amount": 500}  (JSON text)
     │
     │  HTTP/1.1 (text-based headers, one request per connection)
     │
  Server: 200 OK  {"id": 1, "status": "SUCCESS"}  (JSON text)

gRPC:
  Client: paymentService.charge(PaymentRequest{amount: 500})  (binary Protobuf)
     │
     │  HTTP/2 (binary, multiplexed, header compression)
     │
  Server: PaymentResponse{id: 1, status: SUCCESS}  (binary Protobuf)
```

```protobuf
// payment.proto — Interface Definition Language
syntax = "proto3";

service PaymentService {
    rpc Charge(PaymentRequest) returns (PaymentResponse);
    rpc GetPayment(PaymentId) returns (PaymentResponse);
    rpc StreamPayments(PaymentQuery) returns (stream PaymentResponse);  // Server streaming
}

message PaymentRequest {
    string order_id = 1;
    double amount = 2;
    string currency = 3;
}

message PaymentResponse {
    string payment_id = 1;
    PaymentStatus status = 2;
}
```

| Feature | REST (HTTP/JSON) | gRPC (HTTP/2 + Protobuf) |
|---------|-----------------|--------------------------|
| **Protocol** | HTTP/1.1 or HTTP/2 | HTTP/2 always |
| **Payload** | JSON (text, ~2-10x larger) | Protobuf (binary, compact) |
| **Contract** | OpenAPI/Swagger (optional) | `.proto` file (mandatory) |
| **Code generation** | Optional (openapi-generator) | Built-in (protoc compiler) |
| **Streaming** | Limited (SSE, WebSocket) | Native bidirectional streaming |
| **Latency** | Higher (text parsing, larger payload) | Lower (binary, header compression) |
| **Browser support** | Full native support | Limited (gRPC-Web proxy needed) |
| **Caching** | HTTP caching (CDN, browser) | No built-in caching |
| **Tooling** | curl, Postman, browser | grpcurl, BloomRPC, Postman (newer) |
| **Human readability** | Easy (JSON is readable) | Hard (binary on wire) |
| **Load balancing** | Standard L7 LB (Nginx, ALB) | Needs gRPC-aware LB (Envoy, gRPC-LB) |

---

### 9.2 When to Prefer gRPC?

```
Prefer gRPC:                           Prefer REST:
─────────────                          ──────────────
✓ Internal service-to-service           ✓ External/public APIs
✓ High-throughput, low-latency          ✓ Browser clients
✓ Streaming data (real-time feeds)      ✓ Simple CRUD
✓ Polyglot services (Go ↔ Java ↔ Py)   ✓ Quick prototyping
✓ Strong typing needed                  ✓ Cacheable responses
✓ Bandwidth-sensitive (mobile backend)  ✓ Human-debuggable
✓ Bidirectional communication           ✓ Existing REST ecosystem
```

**Production pattern:** Many companies use **REST for external APIs** (client-facing) and **gRPC for internal service-to-service** communication. The API Gateway translates between them.

```
┌────────┐  REST   ┌──────────┐  gRPC   ┌─────────┐  gRPC  ┌─────────┐
│ Mobile │────────→│   API    │────────→│ Order   │───────→│ Payment │
│  App   │  JSON   │ Gateway  │ Protobuf│ Service │        │ Service │
└────────┘         └──────────┘         └─────────┘        └─────────┘
                    Translates
                    REST → gRPC
```

---

### 9.3 Drawbacks of gRPC

| Drawback | Impact | Mitigation |
|----------|--------|------------|
| **No browser support** | Can't call from JavaScript directly | gRPC-Web proxy (Envoy) |
| **Not human-readable** | Can't debug with curl easily | Use grpcurl, Postman, server reflection |
| **Load balancing complexity** | L7 LB must understand HTTP/2 framing | Use Envoy or gRPC-native LB |
| **No HTTP caching** | Can't use CDN or browser cache | Implement application-level caching |
| **Breaking changes** | Proto field renumbering breaks clients | Follow Protobuf evolution rules |
| **Learning curve** | Teams unfamiliar with Protobuf, HTTP/2 | Training, gradual adoption |
| **Debugging in production** | Binary payload hard to inspect | Enable gRPC reflection, structured logging |

---

## 10. Performance & Optimization

### 10.1 How to Reduce Latency in Service Calls

```
┌──────────────────────────────────────────────────────────────┐
│              Latency Reduction Strategies                      │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Connection Pooling — reuse TCP connections               │
│     Before: [DNS][TCP][TLS][Request][Response] = 200ms      │
│     After:  [Request][Response] = 20ms (existing connection) │
│                                                              │
│  2. HTTP/2 — multiplexing over single connection             │
│     Before: 10 requests = 10 sequential connections          │
│     After:  10 requests = 1 connection, multiplexed          │
│                                                              │
│  3. Caching — avoid the call entirely                        │
│     Local cache (Caffeine): < 1ms                            │
│     Distributed cache (Redis): 1-5ms                         │
│     Service call: 10-200ms                                   │
│                                                              │
│  4. Parallel calls — fan-out, don't chain                    │
│     Before: A→B(100ms)→C(100ms)→D(100ms) = 300ms           │
│     After:  A→[B,C,D in parallel] = 100ms (max of all)     │
│                                                              │
│  5. Data locality — keep data close                          │
│     Replicate needed data via events                         │
│     Avoid cross-service joins                                │
│                                                              │
│  6. Compression — smaller payloads, faster transfer          │
│     gzip for REST, Protobuf for gRPC (already compact)      │
│                                                              │
│  7. Connection keep-alive — avoid TCP/TLS setup per request  │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

---

### 10.2 Connection Pooling

```
Without pooling:
  Request 1: [DNS][TCP SYN/ACK][TLS handshake][HTTP Request][Response][TCP FIN]
  Request 2: [DNS][TCP SYN/ACK][TLS handshake][HTTP Request][Response][TCP FIN]
  Request 3: [DNS][TCP SYN/ACK][TLS handshake][HTTP Request][Response][TCP FIN]
  Each request: ~200ms overhead

With pooling:
  Request 1: [DNS][TCP][TLS][HTTP Request][Response]  ← first creates connection
  Request 2:                  [HTTP Request][Response]  ← reuses connection
  Request 3:                  [HTTP Request][Response]  ← reuses connection
  Subsequent requests: ~20ms (just the HTTP round-trip)
```

```java
// Apache HttpClient connection pool (used by Feign)
@Bean
public CloseableHttpClient httpClient() {
    PoolingHttpClientConnectionManager connManager =
        new PoolingHttpClientConnectionManager();
    connManager.setMaxTotal(200);              // Total connections across all hosts
    connManager.setDefaultMaxPerRoute(50);     // Max connections per host

    return HttpClients.custom()
        .setConnectionManager(connManager)
        .setKeepAliveStrategy((response, context) -> 30_000) // 30s keep-alive
        .evictIdleConnections(60, TimeUnit.SECONDS)
        .build();
}

// WebClient connection pool (Reactor Netty)
@Bean
public WebClient webClient() {
    ConnectionProvider provider = ConnectionProvider.builder("custom")
        .maxConnections(200)
        .maxIdleTime(Duration.ofSeconds(30))
        .maxLifeTime(Duration.ofMinutes(5))
        .pendingAcquireTimeout(Duration.ofSeconds(10))
        .build();

    HttpClient httpClient = HttpClient.create(provider)
        .option(ChannelOption.CONNECT_TIMEOUT_MILLIS, 3000);

    return WebClient.builder()
        .clientConnector(new ReactorClientHttpConnector(httpClient))
        .build();
}
```

---

### 10.3 How HTTP/2 Helps

```
HTTP/1.1:
┌────────────────────────────────────────────────┐
│  Connection 1: [Request 1]──[Response 1]       │
│  Connection 2: [Request 2]──[Response 2]       │
│  Connection 3: [Request 3]──[Response 3]       │
│  6 connections max per host (browser limit)    │
│  Head-of-line blocking within each connection  │
└────────────────────────────────────────────────┘

HTTP/2:
┌────────────────────────────────────────────────┐
│  Single Connection (multiplexed):              │
│  Stream 1: [Req1]─────────[Resp1]              │
│  Stream 2: [Req2]──[Resp2]                     │
│  Stream 3: [Req3]──────────────[Resp3]         │
│                                                │
│  All requests on ONE TCP connection            │
│  No head-of-line blocking between streams      │
│  Binary framing (not text)                     │
│  Header compression (HPACK)                    │
└────────────────────────────────────────────────┘
```

| Feature | HTTP/1.1 | HTTP/2 |
|---------|---------|--------|
| **Multiplexing** | No (one request per connection) | Yes (many streams per connection) |
| **Header compression** | None | HPACK (reduces header size ~90%) |
| **Server push** | No | Yes (proactively send resources) |
| **Binary framing** | Text-based | Binary (more efficient parsing) |
| **Connection count** | Many (6 per host typical) | One (with many streams) |
| **gRPC compatible** | No | Yes (gRPC requires HTTP/2) |

---

### 10.4 Batching Requests

```java
// Instead of N individual calls:
for (Long userId : userIds) {
    User user = userClient.getUser(userId);    // N network round-trips!
}

// Batch into single call:
List<User> users = userClient.getUsersByIds(userIds);  // 1 network round-trip

// Or use parallel calls for different services:
CompletableFuture<User> userFuture = CompletableFuture.supplyAsync(
    () -> userClient.getUser(userId));
CompletableFuture<List<Order>> ordersFuture = CompletableFuture.supplyAsync(
    () -> orderClient.getOrdersByUser(userId));
CompletableFuture<PaymentSummary> paymentFuture = CompletableFuture.supplyAsync(
    () -> paymentClient.getSummary(userId));

// Wait for all (parallel, not sequential)
CompletableFuture.allOf(userFuture, ordersFuture, paymentFuture).join();

UserDashboard dashboard = new UserDashboard(
    userFuture.get(),
    ordersFuture.get(),
    paymentFuture.get()
);
// Total latency = max(user, orders, payment) instead of sum
```

---

## 11. Observability & Debugging

### 11.1 Distributed Tracing (Micrometer Tracing / Sleuth)

```
Without tracing:
  User reports: "My order took 5 seconds"
  You: Check Order logs → 50ms processing time
       Check Payment logs → 100ms processing time
       Where did the other 4.8 seconds go?!

With distributed tracing:
  Trace ID: abc-123
  ┌──────────────────────────────────────────────────────────────┐
  │ Span 1: API Gateway (total: 5000ms)                          │
  │   ├── Span 2: Order Service (processing: 50ms)               │
  │   │     ├── Span 3: Inventory Service (processing: 80ms)     │
  │   │     └── Span 4: Payment Service (processing: 100ms)      │
  │   │           └── Span 5: Payment Gateway (HTTP: 4500ms!) ← │
  │   └── Span 6: Notification (async: 30ms)                     │
  └──────────────────────────────────────────────────────────────┘
  Found it! External payment gateway is the bottleneck.
```

```
Trace propagation:
─────────────────────────────────────────────────────────────
                  traceparent: 00-{traceId}-{spanId}-01
Client ────────────────────────────→ API Gateway
  traceId=abc123, spanId=001              │
                                          │ Creates child span
                                          ▼
                                    Order Service
  traceId=abc123, spanId=002              │
  parentId=001                            │ Propagates via HTTP header
                                          ▼
                                    Payment Service
  traceId=abc123, spanId=003
  parentId=002

All services share the SAME traceId → correlate all logs
```

**Spring Boot 3 setup (Micrometer Tracing + Zipkin):**

```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-tracing-bridge-brave</artifactId>
</dependency>
<dependency>
    <groupId>io.zipkin.reporter2</groupId>
    <artifactId>zipkin-reporter-brave</artifactId>
</dependency>
```

```yaml
management:
  tracing:
    sampling:
      probability: 1.0    # Sample 100% in dev, 10-20% in production
  zipkin:
    tracing:
      endpoint: http://zipkin:9411/api/v2/spans
```

---

### 11.2 Correlation IDs

```
┌──────────────────────────────────────────────────────────────┐
│  Correlation ID = a unique identifier that travels with      │
│  the request across ALL services and log entries             │
│                                                              │
│  Request from client:                                        │
│  X-Correlation-ID: req-550e8400-e29b                        │
│                                                              │
│  Order Service log:                                          │
│  [req-550e8400] INFO Creating order for user 42              │
│                                                              │
│  Payment Service log:                                        │
│  [req-550e8400] INFO Charging $50 for order 123              │
│                                                              │
│  Inventory Service log:                                      │
│  [req-550e8400] INFO Reserving 3 items for order 123         │
│                                                              │
│  One grep: all logs for this request across all services     │
└──────────────────────────────────────────────────────────────┘
```

```java
// MDC-based correlation ID filter
@Component
public class CorrelationIdFilter extends OncePerRequestFilter {

    private static final String CORRELATION_HEADER = "X-Correlation-ID";

    @Override
    protected void doFilterInternal(HttpServletRequest request,
            HttpServletResponse response, FilterChain chain)
            throws ServletException, IOException {

        String correlationId = request.getHeader(CORRELATION_HEADER);
        if (correlationId == null) {
            correlationId = UUID.randomUUID().toString();
        }

        MDC.put("correlationId", correlationId);
        response.setHeader(CORRELATION_HEADER, correlationId);

        try {
            chain.doFilter(request, response);
        } finally {
            MDC.remove("correlationId");
        }
    }
}
```

```xml
<!-- logback-spring.xml pattern -->
<pattern>%d{yyyy-MM-dd HH:mm:ss} [%X{correlationId}] [%thread] %-5level %logger - %msg%n</pattern>
```

---

### 11.3 Debugging Slow Service Calls

```
Debugging Playbook:
─────────────────────────────────────────────────────────────

Step 1: Identify the slow span (Zipkin/Jaeger)
  └── Which service/call is taking the most time?

Step 2: Check metrics for that service
  ├── CPU/Memory (is it resource-constrained?)
  ├── Connection pool utilization (is it exhausted?)
  ├── Thread pool utilization (are threads blocked?)
  └── GC pauses (is JVM garbage collecting?)

Step 3: Check downstream dependencies
  ├── Database query performance (slow query log)
  ├── External API latency (metrics, traces)
  └── Message broker lag (Kafka consumer lag)

Step 4: Check network
  ├── DNS resolution time
  ├── TCP connection time
  └── TLS handshake time (missing connection pooling?)

Step 5: Profile the application
  ├── Thread dump (jstack) — where are threads stuck?
  ├── Heap dump — memory issues?
  └── Async-profiler — CPU flame graph

Common findings:
  • Missing connection pool → TCP/TLS overhead on every call
  • N+1 query → 100 DB calls instead of 1
  • Missing index → full table scan
  • Large payload serialization → JSON processing time
  • Thread pool exhaustion → requests queuing
  • GC storm → stop-the-world pauses
```

---

## 12. Testing Inter-Service Communication

### 12.1 Testing Feign Clients

```java
// Integration test with WireMock
@SpringBootTest
@AutoConfigureWireMock(port = 0) // Random port
class PaymentClientTest {

    @Autowired
    private PaymentClient paymentClient;

    @Test
    void shouldChargePayment() {
        stubFor(post(urlEqualTo("/api/v1/payments"))
            .withRequestBody(matchingJsonPath("$.amount", equalTo("500")))
            .willReturn(aResponse()
                .withStatus(200)
                .withHeader("Content-Type", "application/json")
                .withBody("""
                    {
                      "paymentId": "pay-123",
                      "status": "SUCCESS"
                    }
                    """)));

        PaymentResponse response = paymentClient.charge(
            new PaymentRequest("order-1", BigDecimal.valueOf(500)));

        assertThat(response.getStatus()).isEqualTo("SUCCESS");
    }

    @Test
    void shouldHandleTimeout() {
        stubFor(post(urlEqualTo("/api/v1/payments"))
            .willReturn(aResponse()
                .withStatus(200)
                .withFixedDelay(10000))); // 10 second delay

        assertThrows(FeignException.class, () ->
            paymentClient.charge(new PaymentRequest("order-1", BigDecimal.valueOf(500))));
    }

    @Test
    void shouldHandleServerError() {
        stubFor(post(urlEqualTo("/api/v1/payments"))
            .willReturn(aResponse().withStatus(500)));

        // Verify fallback is triggered
        PaymentResponse response = paymentClient.charge(
            new PaymentRequest("order-1", BigDecimal.valueOf(500)));

        assertThat(response.getStatus()).isEqualTo("FAILED");
    }
}
```

---

### 12.2 Contract Testing with Pact

Contract testing ensures the **consumer** and **provider** agree on the API shape without deploying both together.

```
Problem without contract testing:
  Team A (Order) develops against Payment API v1
  Team B (Payment) changes API to v2
  Integration breaks in production!

With Pact:
  ┌────────────────┐          ┌────────────────┐
  │  Order Service  │          │ Payment Service │
  │  (Consumer)     │          │  (Provider)     │
  │                 │          │                 │
  │ 1. Define      │          │                 │
  │    expectations │          │                 │
  │    (pact file)  │─────────→│ 2. Verify pact │
  │                 │  shared  │    against real │
  │                 │   pact   │    API          │
  └────────────────┘          └────────────────┘
                  │
                  ▼
         ┌──────────────┐
         │  Pact Broker  │  (stores and shares contracts)
         └──────────────┘
```

**Consumer side (Order Service):**

```java
@ExtendWith(PactConsumerTestExt.class)
class PaymentClientPactTest {

    @Pact(provider = "payment-service", consumer = "order-service")
    public V4Pact createPact(PactDslWithProvider builder) {
        return builder
            .given("payment service is available")
            .uponReceiving("a charge request")
                .path("/api/v1/payments")
                .method("POST")
                .body(new PactDslJsonBody()
                    .stringType("orderId", "order-123")
                    .decimalType("amount", 500.0))
            .willRespondWith()
                .status(200)
                .body(new PactDslJsonBody()
                    .stringType("paymentId")
                    .stringValue("status", "SUCCESS"))
            .toPact(V4Pact.class);
    }

    @Test
    @PactTestFor(pactMethod = "createPact")
    void testCharge(MockServer mockServer) {
        PaymentClient client = createClient(mockServer.getUrl());
        PaymentResponse response = client.charge(new PaymentRequest("order-123", 500.0));
        assertThat(response.getStatus()).isEqualTo("SUCCESS");
    }
}
```

**Provider side (Payment Service):**

```java
@Provider("payment-service")
@PactBroker(url = "https://pact-broker.company.com")
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class PaymentProviderPactTest {

    @TestTemplate
    @ExtendWith(PactVerificationInvocationContextProvider.class)
    void verifyPact(PactVerificationContext context) {
        context.verifyInteraction();
    }

    @State("payment service is available")
    void setupState() {
        // Set up test data for this state
    }
}
```

---

### 12.3 Mocking Downstream Services

| Tool | Type | Best For |
|------|------|----------|
| **WireMock** | HTTP mock server | Feign/REST client testing |
| **MockServer** | HTTP mock + proxy | Complex scenarios, record/replay |
| **Testcontainers** | Real service in Docker | Integration tests with real Kafka/RabbitMQ |
| **@MockBean** | Spring mock | Unit tests, mock the client bean |
| **Pact** | Contract mock | Consumer-driven contract testing |

```java
// Strategy 1: @MockBean (unit test — fast, no network)
@SpringBootTest
class OrderServiceTest {

    @MockBean
    private PaymentClient paymentClient;

    @Test
    void shouldCreateOrder() {
        when(paymentClient.charge(any()))
            .thenReturn(new PaymentResponse("pay-1", "SUCCESS"));

        Order order = orderService.createOrder(request);
        assertThat(order.getStatus()).isEqualTo(OrderStatus.CONFIRMED);
    }
}

// Strategy 2: Testcontainers (integration — slower, real broker)
@SpringBootTest
@Testcontainers
class OrderEventTest {

    @Container
    static KafkaContainer kafka = new KafkaContainer(
        DockerImageName.parse("confluentinc/cp-kafka:7.5.0"));

    @DynamicPropertySource
    static void kafkaProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.kafka.bootstrap-servers", kafka::getBootstrapServers);
    }

    @Test
    void shouldPublishOrderEvent() {
        orderService.createOrder(request);
        // Verify event was published to Kafka
        ConsumerRecord<String, String> record = KafkaTestUtils
            .getSingleRecord(consumer, "order-events");
        assertThat(record.value()).contains("ORDER_CREATED");
    }
}
```

---

## 13. Scenario-Based Questions (Interview Gold)

### 13.1 Design Communication: Order → Payment → Inventory

```
┌────────────────────────────────────────────────────────────────┐
│                Production Architecture                          │
│                                                                │
│  ┌────────┐   REST    ┌──────────────┐                        │
│  │ Client │──────────→│ API Gateway  │                        │
│  └────────┘           │ (rate limit, │                        │
│                       │  auth, route)│                        │
│                       └──────┬───────┘                        │
│                              │ REST/gRPC                      │
│                              ▼                                │
│                       ┌─────────────┐                         │
│                       │   Order     │                         │
│                       │   Service   │                         │
│                       └──────┬──────┘                         │
│                              │                                │
│                    ┌─────────┼─────────┐                      │
│                    │         │         │                       │
│              Sync (REST)    Outbox   Async (Kafka)            │
│                    │      Pattern      │                      │
│                    ▼         │         ▼                       │
│             ┌──────────┐    │   ┌──────────────┐             │
│             │ Payment  │    │   │  Inventory   │             │
│             │ Service  │    │   │  Service     │             │
│             └──────────┘    │   └──────────────┘             │
│                             │                                 │
│                             ▼                                 │
│                    ┌──────────────┐                           │
│                    │ Notification │                           │
│                    │   Service    │                           │
│                    └──────────────┘                           │
│                                                                │
│  Decisions:                                                    │
│  • Payment: SYNC (need immediate result for user response)    │
│  • Inventory: ASYNC (can reserve after order confirmed)       │
│  • Notification: ASYNC (fire-and-forget)                      │
│  • Use Saga for rollback on failure                           │
└────────────────────────────────────────────────────────────────┘
```

**Implementation flow:**

```java
@Service
public class OrderOrchestrator {

    public OrderResult placeOrder(OrderRequest request) {
        // Step 1: Create order in PENDING state (local TX)
        Order order = orderService.createOrder(request);

        // Step 2: Charge payment SYNC (user needs to know)
        try {
            PaymentResult payment = paymentClient.charge(
                new PaymentRequest(order.getId(), request.getAmount()));

            if (!payment.isSuccess()) {
                orderService.cancelOrder(order.getId());
                return OrderResult.paymentFailed(payment.getReason());
            }
        } catch (Exception e) {
            orderService.cancelOrder(order.getId());
            return OrderResult.paymentError("Payment service unavailable");
        }

        // Step 3: Confirm order + publish event (local TX + outbox)
        orderService.confirmOrder(order.getId());
        // Outbox event triggers:
        //   - Inventory reservation (async)
        //   - Notification (async)
        //   - Analytics (async)

        return OrderResult.success(order);
    }
}
```

---

### 13.2 If Latency Is High Between Services, What Will You Do?

```
Investigation & Resolution Playbook:
─────────────────────────────────────────────────────────────

1. MEASURE (where is the latency?)
   └── Distributed tracing → identify slowest span
       ├── Network latency? (DNS, TCP, TLS)
       ├── Processing latency? (CPU, DB query)
       └── Queue latency? (message broker backlog)

2. OPTIMIZE THE CALL
   ├── Connection pooling (eliminate TCP/TLS overhead)
   ├── HTTP/2 (multiplexing)
   ├── gRPC (binary protocol, smaller payload)
   ├── Compression (gzip for large payloads)
   └── Request batching (1 batch call vs N individual calls)

3. AVOID THE CALL
   ├── Caching (local: Caffeine, distributed: Redis)
   ├── Data replication via events (keep local copy)
   ├── API Gateway aggregation (BFF pattern)
   └── Precomputed views (CQRS read model)

4. PARALLELIZE
   ├── CompletableFuture.allOf() for independent calls
   ├── Reactor Flux.merge() for reactive
   └── Don't chain what can run in parallel

5. MAKE IT ASYNC
   ├── If result not needed immediately → use events
   ├── Return "accepted" → process in background
   └── WebSocket/SSE to push result when ready

6. SCALE
   ├── More instances of slow service
   ├── Auto-scaling based on latency metrics
   └── Geo-distribute (if network latency)
```

---

### 13.3 If Kafka Consumer Is Slow, How Will You Handle It?

```
Diagnosis:
─────────────────────────────────────────────────────────────
Consumer lag = (latest offset - consumer offset) per partition

  Partition 0: Latest=1000, Consumer=400  → Lag: 600 ⚠️
  Partition 1: Latest=1000, Consumer=950  → Lag: 50  ✓
  Partition 2: Latest=1000, Consumer=200  → Lag: 800 ⚠️

Solutions (in order of effort):
─────────────────────────────────────────────────────────────

1. Scale consumers (match partition count)
   Current:  3 partitions, 1 consumer  → 1 consumer handles all
   Improved: 3 partitions, 3 consumers → 1 consumer per partition

2. Increase partitions + consumers
   Current:  3 partitions  → max 3 consumers in group
   Improved: 12 partitions → up to 12 consumers

3. Optimize processing per message
   ├── Batch database writes (save 100 records at once, not 1 by 1)
   ├── Async processing within consumer
   ├── Remove unnecessary processing
   └── Use connection pooling for DB/HTTP calls

4. Increase max.poll.records
   Process more records per poll cycle (batch processing)

5. Check for processing bottleneck
   ├── Slow database query? → add index
   ├── Slow HTTP call? → add caching/circuit breaker
   ├── GC pauses? → tune JVM
   └── Single slow partition? → check for hot key

6. Dead letter queue for poison pills
   If one message causes repeated failures,
   send to DLQ after N retries → don't block the partition
```

```yaml
spring:
  kafka:
    consumer:
      max-poll-records: 500          # Process more per poll
      max-poll-interval-ms: 600000   # 10 min to process batch
      fetch-min-bytes: 10240         # Wait for 10KB before returning
    listener:
      concurrency: 3                 # 3 consumer threads
      type: batch                    # Batch listener mode
```

---

### 13.4 How to Migrate from REST to Async Messaging

```
Phase 1: Identify candidates
─────────────────────────────────────────────────────────────
  Audit all REST calls:
  ┌──────────────────────────┬───────────┬──────────────────┐
  │ Call                     │ Sync Need │ Migration?        │
  ├──────────────────────────┼───────────┼──────────────────┤
  │ GET /users/{id}          │ Yes       │ No (read query)   │
  │ POST /payments           │ Yes       │ No (need result)  │
  │ POST /notifications/send │ No        │ Yes (fire & forget│
  │ PUT /inventory/reserve   │ Maybe     │ Yes (async OK)    │
  │ POST /analytics/event    │ No        │ Yes (async OK)    │
  └──────────────────────────┴───────────┴──────────────────┘

Phase 2: Strangler Fig pattern (gradual migration)
─────────────────────────────────────────────────────────────

  Step 1: Add event publishing alongside REST call
  ┌──────────┐  REST (existing)  ┌──────────────┐
  │  Order   │──────────────────→│ Notification │
  │  Service │                   │   Service    │
  │          │──event (new)─────→│              │
  └──────────┘    Kafka          └──────────────┘
       Both paths active. Verify events match REST calls.

  Step 2: Consumer reads from events, REST still available
  ┌──────────┐  REST (deprecated)┌──────────────┐
  │  Order   │──────────────────→│ Notification │
  │  Service │                   │   Service    │
  │          │──event (primary)─→│ (reads Kafka)│
  └──────────┘                   └──────────────┘

  Step 3: Remove REST call
  ┌──────────┐                   ┌──────────────┐
  │  Order   │──event (only)────→│ Notification │
  │  Service │    Kafka          │   Service    │
  └──────────┘                   └──────────────┘

Phase 3: Verify and cleanup
  ✓ Monitor consumer lag
  ✓ Verify no data loss (compare counts)
  ✓ Remove REST endpoint after grace period
  ✓ Update documentation and API contracts
```

---

### 13.5 How to Handle API Versioning

```
┌────────────────────────────────────────────────────────────────┐
│                   API Versioning Strategies                      │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. URL Path Versioning (MOST COMMON)                          │
│     GET /api/v1/orders                                         │
│     GET /api/v2/orders                                         │
│     Pros: Simple, clear, cacheable                             │
│     Cons: URL clutter, many routes                             │
│                                                                │
│  2. Header Versioning                                          │
│     GET /api/orders                                            │
│     Accept: application/vnd.company.v2+json                    │
│     Pros: Clean URLs                                           │
│     Cons: Hidden, harder to test in browser                    │
│                                                                │
│  3. Query Parameter                                            │
│     GET /api/orders?version=2                                  │
│     Pros: Easy to switch versions                              │
│     Cons: Optional param, easy to forget                       │
│                                                                │
│  4. No Versioning (Backward Compatible Evolution)              │
│     Only add fields, never remove                              │
│     Consumers ignore unknown fields                            │
│     Pros: Simplest, no version management                      │
│     Cons: Can't make breaking changes                          │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

**Production versioning strategy:**

```
For REST APIs:
  → URL path versioning (/api/v1/, /api/v2/)
  → Support N-1 versions (current + previous)
  → Deprecation headers: Sunset: Sat, 01 Jan 2027 00:00:00 GMT
  → Minimum 6-month deprecation period

For Event Schemas (Kafka):
  → Schema Registry (Avro/Protobuf) with compatibility checks
  → BACKWARD compatible by default (new schema reads old data)
  → Never remove fields, only add with defaults

For gRPC:
  → Protobuf built-in evolution (never reuse field numbers)
  → Package versioning: package company.payment.v2;
  → Add fields freely, deprecate with reserved keyword
```

```java
// Supporting multiple versions in Spring Boot
@RestController
public class OrderController {

    @GetMapping("/api/v1/orders/{id}")
    public OrderV1Response getOrderV1(@PathVariable Long id) {
        Order order = orderService.getOrder(id);
        return orderMapper.toV1Response(order);
    }

    @GetMapping("/api/v2/orders/{id}")
    public OrderV2Response getOrderV2(@PathVariable Long id) {
        Order order = orderService.getOrder(id);
        return orderMapper.toV2Response(order); // Includes new fields
    }
}
```

---

## Quick Reference Cheat Sheet

### Communication Pattern Selection

```
Need immediate result?
  ├── YES → Sync (REST/gRPC)
  │         ├── External client? → REST
  │         └── Internal, high-perf? → gRPC
  │
  └── NO  → Async (Messaging)
            ├── One consumer? → RabbitMQ queue
            ├── Many consumers? → Kafka topic
            └── Complex routing? → RabbitMQ exchange
```

### Resilience Patterns

```
Timeout      → Don't wait forever (always set timeouts)
Retry        → Recover from transient failures (with backoff + jitter)
Circuit Break → Stop calling broken services (fail fast)
Bulkhead     → Isolate failures (limit concurrent calls)
Fallback     → Degrade gracefully (cache, default, queue)
```

### Golden Rules for Production

```
1.  Set timeouts on EVERY external call (connect + read)
2.  Implement circuit breaker for all sync dependencies
3.  Use async messaging wherever immediate result isn't needed
4.  Always have a fallback strategy (cache, default, queue)
5.  Add idempotency keys to all mutating operations
6.  Use connection pooling for all HTTP clients
7.  Implement distributed tracing from day one
8.  Test failure scenarios (WireMock latency/errors, Chaos Engineering)
9.  Monitor consumer lag for all Kafka consumers
10. Version your APIs and event schemas from the start
```

---
