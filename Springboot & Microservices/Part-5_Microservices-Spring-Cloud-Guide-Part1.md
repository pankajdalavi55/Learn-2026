# Microservices & Spring Cloud — Complete Guide (Part 1)

> A comprehensive guide covering Microservices architecture with Spring Boot and Spring Cloud.  
> **Part 1 Scope:** Fundamentals, Communication, Service Discovery, API Gateway, Config Management, Resilience Patterns  
> **Part 2 Scope:** Distributed Tracing, Security, Event-Driven Architecture, Saga Patterns, Kubernetes Deployment

---

## Table of Contents

1. [Microservices Fundamentals](#1-microservices-fundamentals)
2. [Microservices Architecture Patterns](#2-microservices-architecture-patterns)
3. [Communication Patterns](#3-communication-patterns)
4. [Service Discovery](#4-service-discovery)
5. [API Gateway Pattern](#5-api-gateway-pattern)
6. [Configuration Management](#6-configuration-management)
7. [Load Balancing & Resilience Patterns](#7-load-balancing--resilience-patterns)
8. [Interview Questions - Part 1](#8-interview-questions---part-1)

---

## 1. Microservices Fundamentals

### 1.1 What are Microservices?

**Microservices** is an architectural style that structures an application as a collection of **small, autonomous, loosely coupled services**, each running in its own process, communicating through lightweight mechanisms (typically HTTP/REST or messaging).

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Monolith vs Microservices Architecture                        │
│                                                                                  │
│  MONOLITHIC APPLICATION                    MICROSERVICES APPLICATION             │
│  ──────────────────────                    ─────────────────────────             │
│                                                                                  │
│  ┌─────────────────────────────┐          ┌─────────┐  ┌─────────┐              │
│  │                             │          │  User   │  │  Order  │              │
│  │    ┌─────────────────────┐  │          │ Service │  │ Service │              │
│  │    │    User Module      │  │          │  :8081  │  │  :8082  │              │
│  │    ├─────────────────────┤  │          └────┬────┘  └────┬────┘              │
│  │    │    Order Module     │  │               │            │                    │
│  │    ├─────────────────────┤  │          ┌────┴────────────┴────┐              │
│  │    │   Payment Module    │  │          │    API Gateway       │              │
│  │    ├─────────────────────┤  │          │       :8080          │              │
│  │    │  Inventory Module   │  │          └────┬────────────┬────┘              │
│  │    └─────────────────────┘  │               │            │                    │
│  │                             │          ┌────┴────┐  ┌────┴────┐              │
│  │    Single Deployable       │          │ Payment │  │Inventory│              │
│  │    Unit (WAR/EAR)          │          │ Service │  │ Service │              │
│  │                             │          │  :8083  │  │  :8084  │              │
│  └─────────────────────────────┘          └─────────┘  └─────────┘              │
│                                                                                  │
│  Characteristics:                          Characteristics:                      │
│  • Single codebase                        • Independent codebases                │
│  • Single database                        • Database per service                 │
│  • Deploy all or nothing                  • Independent deployments              │
│  • Scale entire app                       • Scale individual services            │
│  • Single tech stack                      • Polyglot (multiple languages)        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Microservices Characteristics

| Characteristic | Description |
|----------------|-------------|
| **Single Responsibility** | Each service handles one business capability |
| **Autonomous** | Services can be developed, deployed, and scaled independently |
| **Decentralized** | Decentralized data management and governance |
| **Failure Isolation** | Failure in one service doesn't cascade to others |
| **Technology Agnostic** | Each service can use different tech stack |
| **Smart Endpoints, Dumb Pipes** | Business logic in services, communication is simple |
| **Design for Failure** | Build with assumption that services will fail |
| **Evolutionary Design** | Services can be replaced or rewritten independently |

### 1.3 When to Use Microservices

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        Microservices Decision Matrix                             │
│                                                                                  │
│  ✅ USE MICROSERVICES WHEN:                ❌ AVOID MICROSERVICES WHEN:         │
│  ─────────────────────────                 ──────────────────────────           │
│                                                                                  │
│  • Large team (10+ developers)             • Small team (< 5 developers)        │
│  • Complex domain with clear boundaries    • Simple CRUD application            │
│  • Different parts scale differently       • Uniform scaling needs              │
│  • Need for independent deployments        • Tight coupling is acceptable       │
│  • Polyglot requirements                   • Single tech stack works            │
│  • High availability requirements          • Monolith downtime acceptable       │
│  • Organization has DevOps maturity        • Limited ops capability             │
│                                                                                  │
│  COMPLEXITY GRAPH:                                                               │
│                                                                                  │
│  Productivity │     Microservices                                                │
│       ▲       │            ╱                                                    │
│       │       │           ╱                                                     │
│       │       │──────────╱────── Crossover Point                                │
│       │       │         ╱                                                       │
│       │       │        ╱   Monolith                                             │
│       │       │       ╱                                                         │
│       └───────┴──────────────────────────▶ System Complexity                    │
│               │       │                                                         │
│               │       └── Start migrating here                                  │
│               └── Start with Monolith                                           │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.4 Microservices Benefits and Challenges

**Benefits:**

```java
// 1. Independent Scaling
// User Service: 10 instances (high traffic)
// Payment Service: 2 instances (low traffic)
// Inventory Service: 5 instances (medium traffic)

// 2. Independent Deployment
// Deploy Order Service without touching User Service
// Zero-downtime deployments per service

// 3. Technology Freedom
// User Service: Java + Spring Boot + PostgreSQL
// Analytics Service: Python + FastAPI + MongoDB
// Real-time Service: Node.js + Socket.IO + Redis

// 4. Fault Isolation
// Payment Service down → Other services continue working
// Circuit breaker prevents cascade failures

// 5. Team Autonomy
// Team A owns User Service (full-stack ownership)
// Team B owns Order Service (independent releases)
```

**Challenges:**

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        Microservices Challenges                                  │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │ DISTRIBUTED SYSTEM COMPLEXITY                                              │  │
│  │ • Network latency and failures                                             │  │
│  │ • Partial failures (some services up, some down)                           │  │
│  │ • Eventual consistency                                                     │  │
│  │ • Distributed transactions                                                 │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │ OPERATIONAL OVERHEAD                                                       │  │
│  │ • More services to deploy and monitor                                      │  │
│  │ • Complex CI/CD pipelines                                                  │  │
│  │ • Distributed logging and tracing                                          │  │
│  │ • Service discovery and load balancing                                     │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │ DATA MANAGEMENT                                                            │  │
│  │ • Database per service → data duplication                                  │  │
│  │ • Cross-service queries are complex                                        │  │
│  │ • Maintaining data consistency                                             │  │
│  │ • Handling distributed transactions (Saga pattern)                         │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │ TESTING COMPLEXITY                                                         │  │
│  │ • Integration tests across services                                        │  │
│  │ • Contract testing                                                         │  │
│  │ • End-to-end testing                                                       │  │
│  │ • Environment management                                                   │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.5 Spring Cloud Overview

**Spring Cloud** provides tools for developers to build some of the common patterns in distributed systems.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         Spring Cloud Ecosystem                                   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                        SPRING BOOT APPLICATION                          │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                      │                                          │
│  ┌───────────────────────────────────┼───────────────────────────────────────┐  │
│  │                          SPRING CLOUD                                     │  │
│  │                                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐   │  │
│  │  │    Service      │  │   API Gateway   │  │   Config Management    │   │  │
│  │  │   Discovery     │  │                 │  │                        │   │  │
│  │  │ ─────────────── │  │ ─────────────── │  │ ───────────────────    │   │  │
│  │  │ • Eureka        │  │ • Spring Cloud  │  │ • Spring Cloud Config │   │  │
│  │  │ • Consul        │  │   Gateway       │  │ • Vault Integration   │   │  │
│  │  │ • Zookeeper     │  │ • Zuul (legacy) │  │ • Consul KV           │   │  │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────────┘   │  │
│  │                                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐   │  │
│  │  │  Load Balancing │  │    Resilience   │  │  Distributed Tracing   │   │  │
│  │  │                 │  │                 │  │                        │   │  │
│  │  │ ─────────────── │  │ ─────────────── │  │ ───────────────────    │   │  │
│  │  │ • Spring Cloud  │  │ • Resilience4j  │  │ • Micrometer Tracing  │   │  │
│  │  │   LoadBalancer  │  │ • Circuit       │  │ • Zipkin              │   │  │
│  │  │ • Ribbon (old)  │  │   Breaker       │  │ • Jaeger              │   │  │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────────┘   │  │
│  │                                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐   │  │
│  │  │    Messaging    │  │    Security     │  │      Kubernetes        │   │  │
│  │  │                 │  │                 │  │     Integration        │   │  │
│  │  │ ─────────────── │  │ ─────────────── │  │ ───────────────────    │   │  │
│  │  │ • Spring Cloud  │  │ • OAuth2        │  │ • Spring Cloud        │   │  │
│  │  │   Stream        │  │ • JWT           │  │   Kubernetes          │   │  │
│  │  │ • Bus           │  │ • Security      │  │ • Service Mesh        │   │  │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────────┘   │  │
│  │                                                                           │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.6 Spring Cloud Version Compatibility

| Spring Cloud | Spring Boot | Notes |
|--------------|-------------|-------|
| 2023.0.x (Leyton) | 3.2.x, 3.3.x | Current (Java 17+) |
| 2022.0.x (Kilburn) | 3.0.x, 3.1.x | Java 17+ required |
| 2021.0.x (Jubilee) | 2.6.x, 2.7.x | Last Java 11 support |
| 2020.0.x (Ilford) | 2.4.x, 2.5.x | Legacy |

---

## 2. Microservices Architecture Patterns

### 2.1 Domain-Driven Design (DDD) and Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    E-Commerce Domain - Bounded Contexts                          │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                           CUSTOMER CONTEXT                               │    │
│  │   ┌────────────┐  ┌────────────┐  ┌────────────┐                        │    │
│  │   │  Customer  │  │  Address   │  │ Preferences│                        │    │
│  │   │   Entity   │  │   Entity   │  │   Entity   │                        │    │
│  │   └────────────┘  └────────────┘  └────────────┘                        │    │
│  │   • Customer Registration          • User Service (Microservice)        │    │
│  │   • Profile Management                                                   │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                            ORDER CONTEXT                                 │    │
│  │   ┌────────────┐  ┌────────────┐  ┌────────────┐                        │    │
│  │   │   Order    │  │ OrderItem  │  │  Customer  │ ← Different definition │    │
│  │   │   Entity   │  │   Entity   │  │  (ID only) │    than Customer ctx   │    │
│  │   └────────────┘  └────────────┘  └────────────┘                        │    │
│  │   • Order Placement                • Order Service (Microservice)       │    │
│  │   • Order Status Tracking                                                │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                          INVENTORY CONTEXT                               │    │
│  │   ┌────────────┐  ┌────────────┐  ┌────────────┐                        │    │
│  │   │  Product   │  │   Stock    │  │ Warehouse  │                        │    │
│  │   │   Entity   │  │   Entity   │  │   Entity   │                        │    │
│  │   └────────────┘  └────────────┘  └────────────┘                        │    │
│  │   • Stock Management               • Inventory Service (Microservice)   │    │
│  │   • Inventory Tracking                                                   │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                          PAYMENT CONTEXT                                 │    │
│  │   ┌────────────┐  ┌────────────┐  ┌────────────┐                        │    │
│  │   │  Payment   │  │Transaction │  │   Refund   │                        │    │
│  │   │   Entity   │  │   Entity   │  │   Entity   │                        │    │
│  │   └────────────┘  └────────────┘  └────────────┘                        │    │
│  │   • Payment Processing             • Payment Service (Microservice)     │    │
│  │   • Refund Handling                                                      │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Database Per Service Pattern

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        Database Per Service Pattern                              │
│                                                                                  │
│     ┌─────────────┐      ┌─────────────┐      ┌─────────────┐                   │
│     │    User     │      │    Order    │      │  Inventory  │                   │
│     │   Service   │      │   Service   │      │   Service   │                   │
│     └──────┬──────┘      └──────┬──────┘      └──────┬──────┘                   │
│            │                    │                    │                          │
│            ▼                    ▼                    ▼                          │
│     ┌─────────────┐      ┌─────────────┐      ┌─────────────┐                   │
│     │ PostgreSQL  │      │   MySQL     │      │  MongoDB    │                   │
│     │  (Users)    │      │  (Orders)   │      │ (Products)  │                   │
│     └─────────────┘      └─────────────┘      └─────────────┘                   │
│                                                                                  │
│  BENEFITS:                              CHALLENGES:                              │
│  ─────────                              ───────────                              │
│  • Loose coupling                       • Cross-service queries                  │
│  • Independent scaling                  • Data consistency                       │
│  • Right database for the job           • Distributed transactions               │
│  • Isolated failures                    • Data duplication                       │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Shared Database Pattern (Anti-Pattern with Exceptions)

```java
// ⚠️ Generally an anti-pattern, but acceptable in specific scenarios:
// 1. Migration phase from monolith
// 2. Read-only shared reference data
// 3. Reporting databases (CQRS read models)

// Bounded Context separation even with shared DB:
@Entity
@Table(name = "products", schema = "inventory")  // Schema separation
public class Product { }

@Entity
@Table(name = "products", schema = "catalog")    // Different view
public class CatalogProduct { }
```

### 2.4 Decomposition Strategies

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                     Service Decomposition Strategies                             │
│                                                                                  │
│  1. DECOMPOSE BY BUSINESS CAPABILITY                                            │
│  ─────────────────────────────────────                                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │   Sales     │ │  Marketing  │ │   Finance   │ │     HR      │               │
│  │  Capability │ │  Capability │ │  Capability │ │  Capability │               │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘               │
│         ▼               ▼               ▼               ▼                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │Order Service│ │Campaign Svc │ │Payment Svc  │ │Employee Svc │               │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘               │
│                                                                                  │
│  2. DECOMPOSE BY SUBDOMAIN (DDD)                                                │
│  ────────────────────────────────                                               │
│  Core Domain → Most business value (build in-house)                             │
│  Supporting Domain → Necessary but not core (could outsource)                   │
│  Generic Domain → Common functionality (use off-the-shelf)                      │
│                                                                                  │
│  3. STRANGLER FIG PATTERN (Monolith Migration)                                  │
│  ─────────────────────────────────────────────                                  │
│                                                                                  │
│  Phase 1:  ┌──────────────────────────────────────┐                             │
│            │           Monolith                   │                             │
│            │  [Users] [Orders] [Payments] [...]  │                             │
│            └──────────────────────────────────────┘                             │
│                                                                                  │
│  Phase 2:  ┌─────────────┐  ┌────────────────────────────┐                      │
│            │User Service │  │        Monolith            │                      │
│            │ (extracted) │  │ [Orders] [Payments] [...]  │                      │
│            └─────────────┘  └────────────────────────────┘                      │
│                   ▲                    │                                         │
│                   └────────────────────┘ Facade routes to new service           │
│                                                                                  │
│  Phase N:  ┌─────────┐ ┌──────────┐ ┌───────────┐ ┌─────────────┐               │
│            │  User   │ │  Order   │ │  Payment  │ │   Other     │               │
│            │ Service │ │  Service │ │  Service  │ │  Services   │               │
│            └─────────┘ └──────────┘ └───────────┘ └─────────────┘               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.5 Sample Project Structure

```
ecommerce-microservices/
├── api-gateway/
│   ├── src/main/java/
│   │   └── com.example.gateway/
│   │       ├── GatewayApplication.java
│   │       └── config/
│   │           └── GatewayConfig.java
│   ├── src/main/resources/
│   │   └── application.yml
│   └── pom.xml
│
├── config-server/
│   ├── src/main/java/
│   │   └── com.example.config/
│   │       └── ConfigServerApplication.java
│   └── pom.xml
│
├── discovery-server/
│   ├── src/main/java/
│   │   └── com.example.discovery/
│   │       └── DiscoveryServerApplication.java
│   └── pom.xml
│
├── user-service/
│   ├── src/main/java/
│   │   └── com.example.user/
│   │       ├── UserServiceApplication.java
│   │       ├── controller/
│   │       ├── service/
│   │       ├── repository/
│   │       ├── model/
│   │       └── dto/
│   └── pom.xml
│
├── order-service/
│   ├── src/main/java/
│   │   └── com.example.order/
│   │       ├── OrderServiceApplication.java
│   │       ├── controller/
│   │       ├── service/
│   │       ├── repository/
│   │       ├── model/
│   │       ├── dto/
│   │       └── client/         # Feign clients
│   │           └── UserClient.java
│   └── pom.xml
│
├── inventory-service/
├── payment-service/
├── notification-service/
│
├── docker-compose.yml
└── pom.xml (parent)
```

---

## 3. Communication Patterns

### 3.1 Synchronous vs Asynchronous Communication

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Communication Patterns Comparison                             │
│                                                                                  │
│  SYNCHRONOUS (Request-Response)          ASYNCHRONOUS (Event-Driven)            │
│  ──────────────────────────────          ─────────────────────────              │
│                                                                                  │
│  ┌─────────┐  Request   ┌─────────┐     ┌─────────┐  Publish   ┌─────────┐     │
│  │ Service │ ─────────▶ │ Service │     │ Service │ ─────────▶ │ Message │     │
│  │    A    │ ◀───────── │    B    │     │    A    │            │  Broker │     │
│  └─────────┘  Response  └─────────┘     └─────────┘            └────┬────┘     │
│                                                                      │          │
│  • A waits for B                         ┌─────────┐   Subscribe    │          │
│  • Tight coupling                        │ Service │ ◀──────────────┘          │
│  • Simple to understand                  │    B    │                            │
│  • Cascading failures                    └─────────┘                            │
│                                                                                  │
│  Use Cases:                              • A doesn't wait                        │
│  • Query data                            • Loose coupling                        │
│  • Need immediate response               • Better fault tolerance                │
│  • CRUD operations                       • Eventual consistency                  │
│                                                                                  │
│                                          Use Cases:                              │
│                                          • Fire and forget                       │
│                                          • Event notifications                   │
│                                          • Long-running processes                │
│                                          • Decoupled workflows                   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 REST Communication with RestTemplate

```java
// ❌ RestTemplate - Legacy (Still works but not recommended)
@Configuration
public class RestTemplateConfig {
    
    @Bean
    @LoadBalanced  // Enable service discovery
    public RestTemplate restTemplate() {
        return new RestTemplateBuilder()
            .setConnectTimeout(Duration.ofSeconds(5))
            .setReadTimeout(Duration.ofSeconds(5))
            .build();
    }
}

@Service
public class OrderService {
    
    private final RestTemplate restTemplate;
    
    public OrderService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }
    
    public User getUserById(Long userId) {
        // Using service name (with @LoadBalanced)
        String url = "http://user-service/api/users/{id}";
        return restTemplate.getForObject(url, User.class, userId);
    }
    
    public User createUser(UserRequest request) {
        String url = "http://user-service/api/users";
        ResponseEntity<User> response = restTemplate.postForEntity(
            url, request, User.class);
        return response.getBody();
    }
}
```

### 3.3 REST Communication with WebClient (Recommended)

```java
// ✅ WebClient - Modern, Reactive, Non-blocking
@Configuration
public class WebClientConfig {
    
    @Bean
    @LoadBalanced
    public WebClient.Builder webClientBuilder() {
        return WebClient.builder()
            .defaultHeader(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON_VALUE);
    }
    
    @Bean
    public WebClient userServiceWebClient(WebClient.Builder builder) {
        return builder
            .baseUrl("http://user-service")
            .build();
    }
}

@Service
@Slf4j
public class OrderService {
    
    private final WebClient userServiceWebClient;
    
    public OrderService(@Qualifier("userServiceWebClient") WebClient webClient) {
        this.userServiceWebClient = webClient;
    }
    
    // Blocking call (for traditional Spring MVC)
    public User getUserById(Long userId) {
        return userServiceWebClient.get()
            .uri("/api/users/{id}", userId)
            .retrieve()
            .onStatus(HttpStatusCode::is4xxClientError, response -> 
                Mono.error(new UserNotFoundException("User not found: " + userId)))
            .onStatus(HttpStatusCode::is5xxServerError, response -> 
                Mono.error(new ServiceException("User service unavailable")))
            .bodyToMono(User.class)
            .timeout(Duration.ofSeconds(5))
            .block();  // Blocks for result
    }
    
    // Non-blocking call (for WebFlux)
    public Mono<User> getUserByIdAsync(Long userId) {
        return userServiceWebClient.get()
            .uri("/api/users/{id}", userId)
            .retrieve()
            .bodyToMono(User.class)
            .timeout(Duration.ofSeconds(5))
            .doOnError(e -> log.error("Error fetching user: {}", e.getMessage()));
    }
    
    // With retry and error handling
    public Mono<User> getUserWithRetry(Long userId) {
        return userServiceWebClient.get()
            .uri("/api/users/{id}", userId)
            .retrieve()
            .bodyToMono(User.class)
            .retryWhen(Retry.backoff(3, Duration.ofMillis(500))
                .filter(ex -> ex instanceof WebClientResponseException.ServiceUnavailable))
            .onErrorResume(ex -> {
                log.error("Failed to get user after retries", ex);
                return Mono.just(User.unknown());  // Fallback
            });
    }
}
```

### 3.4 Declarative REST with OpenFeign

```java
// Add dependency
// spring-cloud-starter-openfeign

// Enable Feign Clients
@SpringBootApplication
@EnableFeignClients
public class OrderServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(OrderServiceApplication.class, args);
    }
}

// Feign Client Interface
@FeignClient(
    name = "user-service",
    fallbackFactory = UserClientFallbackFactory.class,
    configuration = UserClientConfig.class
)
public interface UserClient {
    
    @GetMapping("/api/users/{id}")
    User getUserById(@PathVariable("id") Long id);
    
    @GetMapping("/api/users")
    List<User> getAllUsers();
    
    @GetMapping("/api/users")
    List<User> getUsersByIds(@RequestParam("ids") List<Long> ids);
    
    @PostMapping("/api/users")
    User createUser(@RequestBody UserRequest request);
    
    @PutMapping("/api/users/{id}")
    User updateUser(@PathVariable("id") Long id, @RequestBody UserRequest request);
    
    @DeleteMapping("/api/users/{id}")
    void deleteUser(@PathVariable("id") Long id);
}

// Feign Configuration
@Configuration
public class UserClientConfig {
    
    @Bean
    public Request.Options requestOptions() {
        return new Request.Options(
            5, TimeUnit.SECONDS,   // Connect timeout
            10, TimeUnit.SECONDS,  // Read timeout
            true                    // Follow redirects
        );
    }
    
    @Bean
    public ErrorDecoder errorDecoder() {
        return new UserClientErrorDecoder();
    }
    
    @Bean
    public Logger.Level feignLoggerLevel() {
        return Logger.Level.FULL;
    }
}

// Custom Error Decoder
public class UserClientErrorDecoder implements ErrorDecoder {
    
    @Override
    public Exception decode(String methodKey, Response response) {
        return switch (response.status()) {
            case 404 -> new UserNotFoundException("User not found");
            case 400 -> new BadRequestException("Invalid request");
            case 503 -> new ServiceUnavailableException("User service unavailable");
            default -> new Exception("Error: " + response.status());
        };
    }
}

// Fallback Factory for Resilience
@Component
@Slf4j
public class UserClientFallbackFactory implements FallbackFactory<UserClient> {
    
    @Override
    public UserClient create(Throwable cause) {
        log.error("User service fallback triggered: {}", cause.getMessage());
        
        return new UserClient() {
            @Override
            public User getUserById(Long id) {
                return User.builder()
                    .id(id)
                    .name("Unknown User")
                    .status("UNAVAILABLE")
                    .build();
            }
            
            @Override
            public List<User> getAllUsers() {
                return Collections.emptyList();
            }
            
            // ... other fallback methods
        };
    }
}

// Using Feign Client in Service
@Service
public class OrderService {
    
    private final UserClient userClient;
    
    public OrderService(UserClient userClient) {
        this.userClient = userClient;
    }
    
    @Transactional
    public Order createOrder(CreateOrderRequest request) {
        // Synchronous call to user service
        User user = userClient.getUserById(request.getUserId());
        
        if (user == null || "INACTIVE".equals(user.getStatus())) {
            throw new InvalidOrderException("User not eligible");
        }
        
        // Create order logic...
        Order order = Order.builder()
            .userId(user.getId())
            .userName(user.getName())
            .items(request.getItems())
            .build();
            
        return orderRepository.save(order);
    }
}
```

### 3.5 gRPC Communication (High Performance)

```protobuf
// user-service.proto
syntax = "proto3";

package com.example.user;

option java_multiple_files = true;
option java_package = "com.example.user.grpc";

service UserService {
    rpc GetUser(GetUserRequest) returns (UserResponse);
    rpc GetUsers(GetUsersRequest) returns (stream UserResponse);
    rpc CreateUser(CreateUserRequest) returns (UserResponse);
}

message GetUserRequest {
    int64 id = 1;
}

message GetUsersRequest {
    repeated int64 ids = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message UserResponse {
    int64 id = 1;
    string name = 2;
    string email = 3;
    string status = 4;
}
```

```java
// gRPC Server Implementation (User Service)
@GrpcService
public class UserGrpcService extends UserServiceGrpc.UserServiceImplBase {
    
    private final UserRepository userRepository;
    
    @Override
    public void getUser(GetUserRequest request, 
                        StreamObserver<UserResponse> responseObserver) {
        User user = userRepository.findById(request.getId())
            .orElseThrow(() -> new StatusRuntimeException(Status.NOT_FOUND));
        
        UserResponse response = UserResponse.newBuilder()
            .setId(user.getId())
            .setName(user.getName())
            .setEmail(user.getEmail())
            .setStatus(user.getStatus())
            .build();
        
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
    
    @Override
    public void getUsers(GetUsersRequest request,
                         StreamObserver<UserResponse> responseObserver) {
        // Server-side streaming
        userRepository.findAllById(request.getIdsList())
            .forEach(user -> {
                responseObserver.onNext(UserResponse.newBuilder()
                    .setId(user.getId())
                    .setName(user.getName())
                    .build());
            });
        responseObserver.onCompleted();
    }
}

// gRPC Client (Order Service)
@Service
public class UserGrpcClient {
    
    private final UserServiceGrpc.UserServiceBlockingStub blockingStub;
    
    public UserGrpcClient(@GrpcClient("user-service") Channel channel) {
        this.blockingStub = UserServiceGrpc.newBlockingStub(channel);
    }
    
    public UserResponse getUser(Long userId) {
        GetUserRequest request = GetUserRequest.newBuilder()
            .setId(userId)
            .build();
        return blockingStub.getUser(request);
    }
}
```

### 3.6 Asynchronous Messaging with RabbitMQ

```java
// RabbitMQ Configuration
@Configuration
public class RabbitMQConfig {
    
    public static final String ORDER_EXCHANGE = "order.exchange";
    public static final String ORDER_CREATED_QUEUE = "order.created.queue";
    public static final String ORDER_CREATED_ROUTING_KEY = "order.created";
    
    @Bean
    public TopicExchange orderExchange() {
        return new TopicExchange(ORDER_EXCHANGE);
    }
    
    @Bean
    public Queue orderCreatedQueue() {
        return QueueBuilder.durable(ORDER_CREATED_QUEUE)
            .withArgument("x-dead-letter-exchange", "order.dlx")
            .build();
    }
    
    @Bean
    public Binding orderCreatedBinding() {
        return BindingBuilder
            .bind(orderCreatedQueue())
            .to(orderExchange())
            .with(ORDER_CREATED_ROUTING_KEY);
    }
    
    @Bean
    public Jackson2JsonMessageConverter messageConverter() {
        return new Jackson2JsonMessageConverter();
    }
    
    @Bean
    public RabbitTemplate rabbitTemplate(ConnectionFactory connectionFactory) {
        RabbitTemplate template = new RabbitTemplate(connectionFactory);
        template.setMessageConverter(messageConverter());
        return template;
    }
}

// Event Classes
@Data
@Builder
public class OrderCreatedEvent {
    private Long orderId;
    private Long userId;
    private List<OrderItem> items;
    private BigDecimal totalAmount;
    private LocalDateTime createdAt;
}

// Publisher (Order Service)
@Service
@Slf4j
public class OrderEventPublisher {
    
    private final RabbitTemplate rabbitTemplate;
    
    public void publishOrderCreated(Order order) {
        OrderCreatedEvent event = OrderCreatedEvent.builder()
            .orderId(order.getId())
            .userId(order.getUserId())
            .items(order.getItems())
            .totalAmount(order.getTotalAmount())
            .createdAt(LocalDateTime.now())
            .build();
        
        rabbitTemplate.convertAndSend(
            RabbitMQConfig.ORDER_EXCHANGE,
            RabbitMQConfig.ORDER_CREATED_ROUTING_KEY,
            event
        );
        
        log.info("Published OrderCreatedEvent for orderId: {}", order.getId());
    }
}

// Consumer (Inventory Service)
@Service
@Slf4j
public class OrderEventConsumer {
    
    private final InventoryService inventoryService;
    
    @RabbitListener(queues = RabbitMQConfig.ORDER_CREATED_QUEUE)
    public void handleOrderCreated(OrderCreatedEvent event) {
        log.info("Received OrderCreatedEvent for orderId: {}", event.getOrderId());
        
        try {
            inventoryService.reserveStock(event.getItems());
            log.info("Stock reserved for orderId: {}", event.getOrderId());
        } catch (InsufficientStockException e) {
            log.error("Insufficient stock for orderId: {}", event.getOrderId());
            // Publish compensation event
            throw new AmqpRejectAndDontRequeueException("Insufficient stock");
        }
    }
}
```

### 3.7 Asynchronous Messaging with Apache Kafka

```java
// Kafka Configuration
@Configuration
@EnableKafka
public class KafkaConfig {
    
    @Bean
    public ProducerFactory<String, Object> producerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        config.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, JsonSerializer.class);
        config.put(ProducerConfig.ACKS_CONFIG, "all");
        config.put(ProducerConfig.RETRIES_CONFIG, 3);
        return new DefaultKafkaProducerFactory<>(config);
    }
    
    @Bean
    public KafkaTemplate<String, Object> kafkaTemplate() {
        return new KafkaTemplate<>(producerFactory());
    }
    
    @Bean
    public ConsumerFactory<String, OrderCreatedEvent> consumerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ConsumerConfig.GROUP_ID_CONFIG, "inventory-service");
        config.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
        config.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, JsonDeserializer.class);
        config.put(JsonDeserializer.TRUSTED_PACKAGES, "com.example.*");
        return new DefaultKafkaConsumerFactory<>(config);
    }
    
    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, OrderCreatedEvent> 
            kafkaListenerContainerFactory() {
        ConcurrentKafkaListenerContainerFactory<String, OrderCreatedEvent> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory());
        factory.setConcurrency(3);  // Parallel consumers
        return factory;
    }
}

// Kafka Producer (Order Service)
@Service
@Slf4j
public class OrderKafkaProducer {
    
    private static final String TOPIC = "orders";
    
    private final KafkaTemplate<String, Object> kafkaTemplate;
    
    public void publishOrderCreated(Order order) {
        OrderCreatedEvent event = OrderCreatedEvent.builder()
            .orderId(order.getId())
            .userId(order.getUserId())
            .items(order.getItems())
            .build();
        
        kafkaTemplate.send(TOPIC, order.getId().toString(), event)
            .whenComplete((result, ex) -> {
                if (ex == null) {
                    log.info("Sent message=[{}] with offset=[{}]", 
                        event, result.getRecordMetadata().offset());
                } else {
                    log.error("Unable to send message=[{}]", event, ex);
                }
            });
    }
}

// Kafka Consumer (Inventory Service)
@Service
@Slf4j
public class OrderKafkaConsumer {
    
    private final InventoryService inventoryService;
    
    @KafkaListener(
        topics = "orders",
        groupId = "inventory-service",
        containerFactory = "kafkaListenerContainerFactory"
    )
    public void consume(
            @Payload OrderCreatedEvent event,
            @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
            @Header(KafkaHeaders.OFFSET) long offset) {
        
        log.info("Received event: {} from partition: {} with offset: {}", 
            event, partition, offset);
        
        inventoryService.reserveStock(event.getItems());
    }
    
    // With manual acknowledgment
    @KafkaListener(topics = "orders", groupId = "payment-service")
    public void consumeWithAck(
            OrderCreatedEvent event,
            Acknowledgment acknowledgment) {
        try {
            // Process event
            acknowledgment.acknowledge();
        } catch (Exception e) {
            // Don't acknowledge - message will be redelivered
            log.error("Error processing event", e);
        }
    }
}
```

### 3.8 Communication Pattern Summary

| Pattern | Protocol | When to Use | Trade-offs |
|---------|----------|-------------|------------|
| **REST** | HTTP | CRUD ops, Simple queries | Easy, but blocking |
| **gRPC** | HTTP/2 | High-performance, Streaming | Fast, but complex setup |
| **RabbitMQ** | AMQP | Task queues, RPC | Rich routing, moderate scale |
| **Kafka** | TCP | Event streaming, High volume | High throughput, complex |
| **GraphQL** | HTTP | Flexible queries, Mobile | Flexible, but overhead |

---

## 4. Service Discovery

### 4.1 Why Service Discovery?

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    The Service Discovery Problem                                 │
│                                                                                  │
│  WITHOUT SERVICE DISCOVERY:                                                      │
│  ─────────────────────────────                                                   │
│                                                                                  │
│  ┌──────────────┐         ┌──────────────┐                                      │
│  │Order Service │ ──────▶ │ 192.168.1.10 │  Hardcoded IP!                       │
│  └──────────────┘         │  User Service │                                      │
│                           └──────────────┘                                      │
│                                                                                  │
│  Problems:                                                                       │
│  • IPs change when services restart                                             │
│  • Can't scale dynamically                                                      │
│  • Configuration nightmare                                                       │
│  • No health checking                                                           │
│                                                                                  │
│  WITH SERVICE DISCOVERY:                                                         │
│  ───────────────────────                                                         │
│                                                                                  │
│  ┌──────────────┐  1. Where is    ┌─────────────────┐                           │
│  │Order Service │  user-service?  │ Service Registry │                          │
│  │              │ ──────────────▶ │  ┌───────────┐   │                          │
│  │              │                 │  │user-svc   │   │                          │
│  │              │                 │  │ :8081     │   │                          │
│  │              │                 │  │ :8082     │   │ ◀─── Services register   │
│  │              │                 │  │ :8083     │   │       themselves         │
│  │              │ ◀────────────── │  └───────────┘   │                          │
│  │              │  2. Use         └─────────────────┘                           │
│  │              │  192.168.1.10:8081                                            │
│  └──────────────┘  (healthy instance)                                           │
│                                                                                  │
│  Benefits:                                                                       │
│  • Dynamic service location                                                      │
│  • Automatic load balancing                                                      │
│  • Health checking                                                               │
│  • Self-registration/deregistration                                              │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Netflix Eureka Server

```java
// Discovery Server Application
// pom.xml: spring-cloud-starter-netflix-eureka-server

@SpringBootApplication
@EnableEurekaServer
public class DiscoveryServerApplication {
    public static void main(String[] args) {
        SpringApplication.run(DiscoveryServerApplication.class, args);
    }
}
```

```yaml
# application.yml - Eureka Server
server:
  port: 8761

spring:
  application:
    name: discovery-server

eureka:
  instance:
    hostname: localhost
  client:
    register-with-eureka: false  # Don't register itself
    fetch-registry: false        # Don't fetch registry
    service-url:
      defaultZone: http://${eureka.instance.hostname}:${server.port}/eureka/
  server:
    enable-self-preservation: true
    eviction-interval-timer-in-ms: 5000
    response-cache-update-interval-ms: 3000
```

**Eureka Server Cluster (High Availability):**

```yaml
# eureka-server-1.yml
server:
  port: 8761

spring:
  application:
    name: discovery-server
  profiles:
    active: peer1

eureka:
  instance:
    hostname: eureka-server-1
  client:
    register-with-eureka: true
    fetch-registry: true
    service-url:
      defaultZone: http://eureka-server-2:8762/eureka/,http://eureka-server-3:8763/eureka/

---
# eureka-server-2.yml
server:
  port: 8762

eureka:
  instance:
    hostname: eureka-server-2
  client:
    service-url:
      defaultZone: http://eureka-server-1:8761/eureka/,http://eureka-server-3:8763/eureka/
```

### 4.3 Eureka Client (Service Registration)

```java
// User Service Application
// pom.xml: spring-cloud-starter-netflix-eureka-client

@SpringBootApplication
@EnableDiscoveryClient  // Optional in newer versions (auto-detected)
public class UserServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(UserServiceApplication.class, args);
    }
}
```

```yaml
# application.yml - Eureka Client
server:
  port: 8081

spring:
  application:
    name: user-service  # Service name for registration

eureka:
  client:
    service-url:
      defaultZone: http://localhost:8761/eureka/
    registry-fetch-interval-seconds: 5
    instance-info-replication-interval-seconds: 5
  instance:
    instance-id: ${spring.application.name}:${random.uuid}
    prefer-ip-address: true
    lease-renewal-interval-in-seconds: 10
    lease-expiration-duration-in-seconds: 30
    metadata-map:
      version: 1.0.0
      environment: dev

# Health check for Eureka
management:
  endpoints:
    web:
      exposure:
        include: health,info
  endpoint:
    health:
      show-details: always
```

### 4.4 Service Discovery with DiscoveryClient

```java
@Service
@Slf4j
public class ServiceDiscoveryService {
    
    private final DiscoveryClient discoveryClient;
    
    public ServiceDiscoveryService(DiscoveryClient discoveryClient) {
        this.discoveryClient = discoveryClient;
    }
    
    // Get all registered services
    public List<String> getAllServices() {
        return discoveryClient.getServices();
    }
    
    // Get service instances
    public List<ServiceInstance> getServiceInstances(String serviceName) {
        return discoveryClient.getInstances(serviceName);
    }
    
    // Get single instance URL
    public String getServiceUrl(String serviceName) {
        List<ServiceInstance> instances = discoveryClient.getInstances(serviceName);
        
        if (instances.isEmpty()) {
            throw new ServiceNotFoundException("No instances of " + serviceName);
        }
        
        // Simple round-robin (use Spring Cloud LoadBalancer for production)
        ServiceInstance instance = instances.get(0);
        return instance.getUri().toString();
    }
    
    // Print all services
    public void printRegisteredServices() {
        discoveryClient.getServices().forEach(serviceName -> {
            log.info("Service: {}", serviceName);
            discoveryClient.getInstances(serviceName).forEach(instance -> {
                log.info("  - Instance: {} at {}:{}", 
                    instance.getInstanceId(),
                    instance.getHost(),
                    instance.getPort());
            });
        });
    }
}
```

### 4.5 Consul Service Discovery (Alternative)

```yaml
# application.yml - Consul Client
spring:
  application:
    name: user-service
  cloud:
    consul:
      host: localhost
      port: 8500
      discovery:
        service-name: ${spring.application.name}
        instance-id: ${spring.application.name}:${random.uuid}
        health-check-interval: 10s
        health-check-path: /actuator/health
        prefer-ip-address: true
        tags:
          - version=1.0.0
          - env=dev
```

```java
// Consul-based Service
@SpringBootApplication
@EnableDiscoveryClient
public class UserServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(UserServiceApplication.class, args);
    }
}
```

### 4.6 Kubernetes Service Discovery

In Kubernetes, service discovery is built-in through DNS:

```yaml
# kubernetes/user-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
    - port: 8080
      targetPort: 8080
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
        - name: user-service
          image: user-service:latest
          ports:
            - containerPort: 8080
```

```yaml
# application.yml - Spring Cloud Kubernetes
spring:
  application:
    name: order-service
  cloud:
    kubernetes:
      discovery:
        enabled: true
        all-namespaces: false
      loadbalancer:
        mode: service  # Use K8s Service for load balancing
```

---

## 5. API Gateway Pattern

### 5.1 Why API Gateway?

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        API Gateway Pattern                                       │
│                                                                                  │
│  WITHOUT API GATEWAY:                     WITH API GATEWAY:                      │
│  ─────────────────────                    ──────────────────                     │
│                                                                                  │
│  ┌─────────┐                              ┌─────────┐                           │
│  │ Client  │                              │ Client  │                           │
│  └────┬────┘                              └────┬────┘                           │
│       │                                        │                                 │
│       ├────────────────┬───────────────┐       │  Single Entry Point           │
│       │                │               │       ▼                                 │
│       ▼                ▼               ▼  ┌─────────────────────┐               │
│  ┌─────────┐      ┌─────────┐    ┌─────────┐│                     │              │
│  │  User   │      │  Order  │    │ Payment ││   API GATEWAY       │              │
│  │ Service │      │ Service │    │ Service ││   ─────────────     │              │
│  └─────────┘      └─────────┘    └─────────┘│   • Routing         │              │
│                                            │   • Authentication  │              │
│  Problems:                                 │   • Rate Limiting   │              │
│  • Client knows all services               │   • Load Balancing  │              │
│  • No unified authentication               │   • Monitoring      │              │
│  • Each service manages cross-cutting      │   • Response Cache  │              │
│  • CORS complexity                         │                     │              │
│  • No single point for monitoring          └──────────┬──────────┘              │
│                                                        │                         │
│                                            ┌───────────┼───────────┐            │
│                                            ▼           ▼           ▼            │
│                                       ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│                                       │  User   │ │  Order  │ │ Payment │       │
│                                       │ Service │ │ Service │ │ Service │       │
│                                       └─────────┘ └─────────┘ └─────────┘       │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Spring Cloud Gateway Setup

```java
// API Gateway Application
// pom.xml: spring-cloud-starter-gateway

@SpringBootApplication
public class ApiGatewayApplication {
    public static void main(String[] args) {
        SpringApplication.run(ApiGatewayApplication.class, args);
    }
}
```

```yaml
# application.yml - Spring Cloud Gateway
server:
  port: 8080

spring:
  application:
    name: api-gateway
  cloud:
    gateway:
      discovery:
        locator:
          enabled: true
          lower-case-service-id: true
      routes:
        # User Service Routes
        - id: user-service
          uri: lb://user-service
          predicates:
            - Path=/api/users/**
          filters:
            - StripPrefix=0
            - name: CircuitBreaker
              args:
                name: userServiceCB
                fallbackUri: forward:/fallback/users
        
        # Order Service Routes
        - id: order-service
          uri: lb://order-service
          predicates:
            - Path=/api/orders/**
          filters:
            - StripPrefix=0
            - name: RequestRateLimiter
              args:
                redis-rate-limiter.replenishRate: 10
                redis-rate-limiter.burstCapacity: 20
        
        # Payment Service Routes
        - id: payment-service
          uri: lb://payment-service
          predicates:
            - Path=/api/payments/**
            - Method=POST,GET
          filters:
            - AddRequestHeader=X-Request-Source, API-Gateway
            - AddResponseHeader=X-Response-Time, ${response.time}
        
        # Inventory with Path Rewrite
        - id: inventory-service
          uri: lb://inventory-service
          predicates:
            - Path=/api/v1/inventory/**
          filters:
            - RewritePath=/api/v1/inventory/(?<segment>.*), /inventory/${segment}

eureka:
  client:
    service-url:
      defaultZone: http://localhost:8761/eureka/
```

### 5.3 Gateway Route Predicates

```yaml
spring:
  cloud:
    gateway:
      routes:
        - id: complex-route
          uri: lb://service
          predicates:
            # Path matching
            - Path=/api/**
            
            # HTTP Method
            - Method=GET,POST
            
            # Header presence/value
            - Header=X-Request-Id, \d+
            
            # Query parameter
            - Query=category, electronics
            
            # Cookie
            - Cookie=session, .+
            
            # Time-based (After, Before, Between)
            - After=2025-01-01T00:00:00+00:00
            - Before=2026-12-31T23:59:59+00:00
            
            # Host
            - Host=**.example.com
            
            # Remote Address
            - RemoteAddr=192.168.1.0/24
            
            # Weight (for canary deployments)
            - Weight=group1, 8  # 80% traffic
```

### 5.4 Gateway Filters

```java
// Global Filter - Applied to all routes
@Component
@Slf4j
public class GlobalLoggingFilter implements GlobalFilter, Ordered {
    
    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        String requestId = UUID.randomUUID().toString();
        long startTime = System.currentTimeMillis();
        
        ServerHttpRequest request = exchange.getRequest().mutate()
            .header("X-Request-Id", requestId)
            .build();
        
        log.info("Incoming request: {} {} - RequestId: {}", 
            request.getMethod(), request.getURI(), requestId);
        
        return chain.filter(exchange.mutate().request(request).build())
            .then(Mono.fromRunnable(() -> {
                long duration = System.currentTimeMillis() - startTime;
                log.info("Request completed: {} - Duration: {}ms - Status: {}", 
                    requestId, duration, exchange.getResponse().getStatusCode());
            }));
    }
    
    @Override
    public int getOrder() {
        return -1;  // Execute first
    }
}

// Custom Route Filter
@Component
public class AuthenticationFilter implements GatewayFilter, Ordered {
    
    private final JwtTokenProvider jwtTokenProvider;
    
    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        ServerHttpRequest request = exchange.getRequest();
        
        // Skip auth for public endpoints
        if (isPublicEndpoint(request.getPath().toString())) {
            return chain.filter(exchange);
        }
        
        String authHeader = request.getHeaders().getFirst(HttpHeaders.AUTHORIZATION);
        
        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            return unauthorized(exchange);
        }
        
        String token = authHeader.substring(7);
        
        try {
            Claims claims = jwtTokenProvider.validateToken(token);
            
            ServerHttpRequest modifiedRequest = request.mutate()
                .header("X-User-Id", claims.getSubject())
                .header("X-User-Roles", claims.get("roles", String.class))
                .build();
            
            return chain.filter(exchange.mutate().request(modifiedRequest).build());
            
        } catch (JwtException e) {
            return unauthorized(exchange);
        }
    }
    
    private Mono<Void> unauthorized(ServerWebExchange exchange) {
        exchange.getResponse().setStatusCode(HttpStatus.UNAUTHORIZED);
        return exchange.getResponse().setComplete();
    }
    
    private boolean isPublicEndpoint(String path) {
        return path.startsWith("/api/auth/") || 
               path.startsWith("/api/public/") ||
               path.equals("/actuator/health");
    }
    
    @Override
    public int getOrder() {
        return 0;
    }
}

// Register Custom Filter Factory
@Component
public class AuthenticationGatewayFilterFactory 
        extends AbstractGatewayFilterFactory<AuthenticationGatewayFilterFactory.Config> {
    
    private final JwtTokenProvider jwtTokenProvider;
    
    public AuthenticationGatewayFilterFactory(JwtTokenProvider jwtTokenProvider) {
        super(Config.class);
        this.jwtTokenProvider = jwtTokenProvider;
    }
    
    @Override
    public GatewayFilter apply(Config config) {
        return new AuthenticationFilter(jwtTokenProvider);
    }
    
    public static class Config {
        // Configuration properties
    }
}
```

### 5.5 Rate Limiting with Redis

```yaml
# application.yml
spring:
  cloud:
    gateway:
      routes:
        - id: rate-limited-route
          uri: lb://service
          predicates:
            - Path=/api/**
          filters:
            - name: RequestRateLimiter
              args:
                redis-rate-limiter.replenishRate: 10      # Requests per second
                redis-rate-limiter.burstCapacity: 20      # Max burst
                redis-rate-limiter.requestedTokens: 1     # Tokens per request
                key-resolver: "#{@userKeyResolver}"
  data:
    redis:
      host: localhost
      port: 6379
```

```java
@Configuration
public class RateLimiterConfig {
    
    // Rate limit by user
    @Bean
    public KeyResolver userKeyResolver() {
        return exchange -> {
            String userId = exchange.getRequest().getHeaders()
                .getFirst("X-User-Id");
            return Mono.just(userId != null ? userId : "anonymous");
        };
    }
    
    // Rate limit by IP
    @Bean
    public KeyResolver ipKeyResolver() {
        return exchange -> Mono.just(
            Objects.requireNonNull(exchange.getRequest().getRemoteAddress())
                .getHostString()
        );
    }
    
    // Rate limit by API key
    @Bean
    public KeyResolver apiKeyResolver() {
        return exchange -> Mono.just(
            exchange.getRequest().getHeaders().getFirst("X-API-Key")
        );
    }
}
```

### 5.6 Circuit Breaker in Gateway

```yaml
spring:
  cloud:
    gateway:
      routes:
        - id: user-service
          uri: lb://user-service
          predicates:
            - Path=/api/users/**
          filters:
            - name: CircuitBreaker
              args:
                name: userCircuitBreaker
                fallbackUri: forward:/fallback/users
                statusCodes: 500,502,503

resilience4j:
  circuitbreaker:
    instances:
      userCircuitBreaker:
        slidingWindowSize: 10
        failureRateThreshold: 50
        waitDurationInOpenState: 10000
        permittedNumberOfCallsInHalfOpenState: 5
```

```java
@RestController
@RequestMapping("/fallback")
public class FallbackController {
    
    @GetMapping("/users")
    public ResponseEntity<Map<String, Object>> userServiceFallback() {
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE)
            .body(Map.of(
                "message", "User service is currently unavailable",
                "timestamp", LocalDateTime.now(),
                "status", "FALLBACK"
            ));
    }
    
    @GetMapping("/orders")
    public ResponseEntity<Map<String, Object>> orderServiceFallback() {
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE)
            .body(Map.of(
                "message", "Order service is currently unavailable",
                "timestamp", LocalDateTime.now(),
                "status", "FALLBACK"
            ));
    }
}
```

---

## 6. Configuration Management

### 6.1 Why Centralized Configuration?

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Configuration Management Problem                              │
│                                                                                  │
│  WITHOUT CONFIG SERVER:                    WITH CONFIG SERVER:                   │
│  ───────────────────────                   ────────────────────                  │
│                                                                                  │
│  ┌──────────┐  ┌──────────┐               ┌──────────────────────┐              │
│  │  Service │  │  Service │               │    Config Server     │              │
│  │    A     │  │    B     │               │                      │              │
│  │┌────────┐│  │┌────────┐│               │  ┌────────────────┐  │              │
│  ││config  ││  ││config  ││               │  │   Git Repo     │  │              │
│  ││.yml    ││  ││.yml    ││               │  │ (Single Source │  │              │
│  │└────────┘│  │└────────┘│               │  │   of Truth)    │  │              │
│  └──────────┘  └──────────┘               │  └────────────────┘  │              │
│                                           └──────────┬───────────┘              │
│  Problems:                                           │                          │
│  • Config scattered across services                  │                          │
│  • Rebuild/redeploy for config changes    ┌──────────┼───────────┐             │
│  • No audit trail                         ▼          ▼           ▼             │
│  • Secret management difficult       ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  • Environment consistency           │Service A│ │Service B│ │Service C│       │
│                                      │(fetches │ │(fetches │ │(fetches │       │
│  Benefits:                           │ config) │ │ config) │ │ config) │       │
│  • Single source of truth            └─────────┘ └─────────┘ └─────────┘       │
│  • Runtime config refresh                                                       │
│  • Version control                                                              │
│  • Environment-specific configs                                                 │
│  • Encryption support                                                           │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Spring Cloud Config Server

```java
// Config Server Application
// pom.xml: spring-cloud-config-server

@SpringBootApplication
@EnableConfigServer
public class ConfigServerApplication {
    public static void main(String[] args) {
        SpringApplication.run(ConfigServerApplication.class, args);
    }
}
```

```yaml
# application.yml - Config Server
server:
  port: 8888

spring:
  application:
    name: config-server
  cloud:
    config:
      server:
        # Git Backend (Recommended for production)
        git:
          uri: https://github.com/your-org/config-repo
          default-label: main
          search-paths: '{application}'  # Folder per app
          clone-on-start: true
          force-pull: true
          username: ${GIT_USERNAME}
          password: ${GIT_PASSWORD}
        
        # Native filesystem (for development)
        # native:
        #   search-locations: classpath:/config-repo
  profiles:
    active: git  # or native

# Encryption key for encrypting sensitive values
encrypt:
  key: ${ENCRYPT_KEY:my-secret-key}

eureka:
  client:
    service-url:
      defaultZone: http://localhost:8761/eureka/
```

### 6.3 Config Repository Structure

```
config-repo/
├── application.yml              # Shared by all services
├── application-dev.yml          # Shared dev settings
├── application-prod.yml         # Shared prod settings
│
├── user-service/
│   ├── user-service.yml         # Default
│   ├── user-service-dev.yml     # Dev profile
│   └── user-service-prod.yml    # Prod profile
│
├── order-service/
│   ├── order-service.yml
│   ├── order-service-dev.yml
│   └── order-service-prod.yml
│
└── api-gateway/
    ├── api-gateway.yml
    └── api-gateway-prod.yml
```

```yaml
# config-repo/application.yml (Shared)
spring:
  jpa:
    show-sql: false
    hibernate:
      ddl-auto: none
  
management:
  endpoints:
    web:
      exposure:
        include: health,info,refresh

logging:
  level:
    root: INFO

---
# config-repo/user-service/user-service.yml
server:
  port: 8081

spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/users
    username: ${DB_USERNAME:user}
    password: '{cipher}AQBjF...'  # Encrypted password

app:
  jwt:
    secret: '{cipher}AQCn...'
    expiration: 3600000

---
# config-repo/user-service/user-service-dev.yml
spring:
  jpa:
    show-sql: true
    hibernate:
      ddl-auto: update

logging:
  level:
    com.example: DEBUG
    org.hibernate.SQL: DEBUG

---
# config-repo/user-service/user-service-prod.yml
spring:
  datasource:
    hikari:
      maximum-pool-size: 20
      minimum-idle: 5

logging:
  level:
    root: WARN
```

### 6.4 Config Client Setup

```java
// User Service - Config Client
// pom.xml: spring-cloud-starter-config
```

```yaml
# bootstrap.yml (or application.yml with spring.config.import)
spring:
  application:
    name: user-service
  profiles:
    active: dev
  config:
    import: optional:configserver:http://localhost:8888
  cloud:
    config:
      fail-fast: true
      retry:
        initial-interval: 1000
        max-attempts: 6
        multiplier: 1.5

# Or in application.yml (Spring Boot 2.4+)
spring:
  application:
    name: user-service
  config:
    import: "configserver:"
  cloud:
    config:
      uri: http://localhost:8888
```

### 6.5 Runtime Configuration Refresh

```java
// Enable refresh scope
@Configuration
@RefreshScope
public class AppConfig {
    
    @Value("${app.feature.enabled:false}")
    private boolean featureEnabled;
    
    @Value("${app.message:default}")
    private String message;
    
    // These values will refresh when /actuator/refresh is called
}

// Controller example
@RestController
@RefreshScope
public class ConfigController {
    
    @Value("${app.dynamic.property}")
    private String dynamicProperty;
    
    @GetMapping("/config")
    public String getConfig() {
        return dynamicProperty;
    }
}
```

```bash
# Trigger refresh for single service
curl -X POST http://localhost:8081/actuator/refresh

# Response shows changed properties
["app.dynamic.property"]
```

### 6.6 Spring Cloud Bus (Broadcast Refresh)

```yaml
# Add to all services
# pom.xml: spring-cloud-starter-bus-amqp (or bus-kafka)

spring:
  rabbitmq:
    host: localhost
    port: 5672
    username: guest
    password: guest
  cloud:
    bus:
      enabled: true
      refresh:
        enabled: true
```

```bash
# Refresh ALL connected services at once
curl -X POST http://localhost:8888/actuator/busrefresh

# Refresh specific service
curl -X POST http://localhost:8888/actuator/busrefresh/user-service
```

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                     Spring Cloud Bus Refresh Flow                                │
│                                                                                  │
│  1. POST /actuator/busrefresh                                                   │
│           │                                                                      │
│           ▼                                                                      │
│  ┌─────────────────┐      2. Publish RefreshEvent                               │
│  │  Config Server  │ ─────────────────────────────────┐                         │
│  └─────────────────┘                                  │                         │
│                                                       ▼                         │
│                                              ┌─────────────────┐                │
│                                              │  Message Broker │                │
│                                              │ (RabbitMQ/Kafka)│                │
│                                              └────────┬────────┘                │
│                                                       │                         │
│                      3. All services receive event    │                         │
│           ┌───────────────────┬───────────────────────┼───────────────┐        │
│           ▼                   ▼                       ▼               ▼        │
│     ┌──────────┐        ┌──────────┐           ┌──────────┐    ┌──────────┐   │
│     │ Service  │        │ Service  │           │ Service  │    │ Service  │   │
│     │    A     │        │    B     │           │    C     │    │    D     │   │
│     │ (refresh)│        │ (refresh)│           │ (refresh)│    │ (refresh)│   │
│     └──────────┘        └──────────┘           └──────────┘    └──────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 6.7 Encrypting Sensitive Data

```bash
# Encrypt a value
curl http://localhost:8888/encrypt -d 'my-secret-password'
# Returns: AQBjF7Gq...

# Decrypt a value
curl http://localhost:8888/decrypt -d 'AQBjF7Gq...'
# Returns: my-secret-password
```

```yaml
# Use encrypted value in config
spring:
  datasource:
    password: '{cipher}AQBjF7GqK9H...'
```

```yaml
# Config server encryption setup
encrypt:
  # Symmetric key
  key: ${ENCRYPT_KEY}
  
  # OR Asymmetric (RSA) - more secure
  # key-store:
  #   location: classpath:keystore.jks
  #   password: ${KEYSTORE_PASSWORD}
  #   alias: configkey
```

---

## 7. Load Balancing & Resilience Patterns

### 7.1 Client-Side Load Balancing

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Load Balancing Strategies                                     │
│                                                                                  │
│  SERVER-SIDE LB                           CLIENT-SIDE LB                         │
│  ──────────────                           ──────────────                         │
│                                                                                  │
│  ┌─────────┐                              ┌─────────┐                           │
│  │ Client  │                              │ Client  │                           │
│  └────┬────┘                              │  ┌───┐  │                           │
│       │                                   │  │LB │  │ ◀── LoadBalancer inside   │
│       ▼                                   │  └─┬─┘  │     each client           │
│  ┌─────────┐                              └────┼────┘                           │
│  │   LB    │                                   │                                 │
│  │ (Nginx, │                         ┌────────┼────────┐                        │
│  │  HAProxy│                         ▼        ▼        ▼                        │
│  └─────────┘                    ┌─────────┐┌─────────┐┌─────────┐               │
│       │                         │Instance ││Instance ││Instance │               │
│  ┌────┼────┐                    │   1     ││   2     ││   3     │               │
│  ▼    ▼    ▼                    └─────────┘└─────────┘└─────────┘               │
│ ┌──┐ ┌──┐ ┌──┐                                                                  │
│ │I1│ │I2│ │I3│                  • No single point of failure                    │
│ └──┘ └──┘ └──┘                  • Client has control                            │
│                                 • Works with service discovery                   │
│ • Extra hop                     • Spring Cloud LoadBalancer                      │
│ • SPOF (if LB fails)                                                            │
│ • Extra infra to manage                                                          │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Spring Cloud LoadBalancer

```java
// Enabled by default with spring-cloud-starter-loadbalancer
// Previously Ribbon (deprecated), now Spring Cloud LoadBalancer

@Configuration
public class LoadBalancerConfig {
    
    // Round Robin (default)
    @Bean
    public ServiceInstanceListSupplier serviceInstanceListSupplier(
            ConfigurableApplicationContext context) {
        return ServiceInstanceListSupplier.builder()
            .withDiscoveryClient()
            .withHealthChecks()
            .withCaching()
            .build(context);
    }
}

// Custom load balancing strategy
@Configuration
public class CustomLoadBalancerConfig {
    
    @Bean
    public ReactorLoadBalancer<ServiceInstance> randomLoadBalancer(
            Environment environment,
            LoadBalancerClientFactory loadBalancerClientFactory) {
        String name = environment.getProperty(LoadBalancerClientFactory.PROPERTY_NAME);
        return new RandomLoadBalancer(
            loadBalancerClientFactory.getLazyProvider(name, ServiceInstanceListSupplier.class),
            name
        );
    }
}

// Apply to specific client
@LoadBalancerClient(name = "user-service", configuration = CustomLoadBalancerConfig.class)
public class UserServiceConfig { }
```

### 7.3 Circuit Breaker Pattern with Resilience4j

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        Circuit Breaker States                                    │
│                                                                                  │
│        ┌─────────────────────────────────────────────────────────────────┐      │
│        │                                                                  │      │
│        │    CLOSED                  OPEN                    HALF-OPEN    │      │
│        │    ──────                  ────                    ─────────    │      │
│        │                                                                  │      │
│        │    [Normal]    Failure     [Fast      Wait        [Test]        │      │
│        │    Requests    threshold   Fail]      timeout     Few requests  │      │
│        │    pass ──────reached────▶ All ──────expires────▶ allowed       │      │
│        │    through                 requests                              │      │
│        │                            fail fast                             │      │
│        │        ▲                                              │          │      │
│        │        │                                              │          │      │
│        │        │                  Success                     │          │      │
│        │        └──────────────────rate OK─────────────────────┘          │      │
│        │                           (enough successful calls)              │      │
│        │                                                                  │      │
│        │                           Failure                                │      │
│        │                           ──────▶ Back to OPEN                  │      │
│        │                                                                  │      │
│        └─────────────────────────────────────────────────────────────────┘      │
│                                                                                  │
│  Example Timeline:                                                               │
│                                                                                  │
│  Time:  0──1──2──3──4──5──6──7──8──9──10─11─12─13─14─15─16─17─18─19─20          │
│         │  │  │  │  │  OPEN (fast fail)│  HALF│    CLOSED (recovered)          │
│  State: └──┴──┴──┴──┘                  └OPEN ─┘                                 │
│         CLOSED (failures accumulate)   │                                        │
│                                        │                                        │
│  Requests: ✓ ✓ ✗ ✗ ✗ ⊘ ⊘ ⊘ ⊘ ⊘ ✓ ✓ ✓ ✓ ✓ ✓ ✓ ✓ ✓ ✓                            │
│                                                                                  │
│  ✓ = Success,  ✗ = Failure,  ⊘ = Rejected (circuit open)                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Add dependency: spring-cloud-starter-circuitbreaker-resilience4j

@Service
@Slf4j
public class UserService {
    
    private final UserClient userClient;
    
    // Circuit Breaker annotation
    @CircuitBreaker(name = "userService", fallbackMethod = "getUserFallback")
    @Retry(name = "userService")
    @Bulkhead(name = "userService")
    public User getUserById(Long userId) {
        log.info("Calling user service for userId: {}", userId);
        return userClient.getUserById(userId);
    }
    
    // Fallback method - must have same return type + Exception param
    private User getUserFallback(Long userId, Exception ex) {
        log.error("Fallback triggered for userId: {} - Reason: {}", userId, ex.getMessage());
        return User.builder()
            .id(userId)
            .name("Cached/Default User")
            .status("UNAVAILABLE")
            .build();
    }
    
    // Different fallbacks for different exceptions
    private User getUserFallback(Long userId, CallNotPermittedException ex) {
        log.error("Circuit is OPEN for userId: {}", userId);
        return User.defaultUser(userId);
    }
    
    private User getUserFallback(Long userId, TimeoutException ex) {
        log.error("Timeout for userId: {}", userId);
        return User.timedOutUser(userId);
    }
}
```

```yaml
# application.yml - Resilience4j Configuration
resilience4j:
  circuitbreaker:
    instances:
      userService:
        register-health-indicator: true
        sliding-window-type: COUNT_BASED
        sliding-window-size: 10
        failure-rate-threshold: 50
        wait-duration-in-open-state: 10s
        permitted-number-of-calls-in-half-open-state: 5
        automatic-transition-from-open-to-half-open-enabled: true
        record-exceptions:
          - java.io.IOException
          - java.net.ConnectException
          - feign.FeignException
        ignore-exceptions:
          - com.example.BusinessException
  
  retry:
    instances:
      userService:
        max-attempts: 3
        wait-duration: 500ms
        exponential-backoff-multiplier: 2
        retry-exceptions:
          - java.io.IOException
          - java.net.ConnectException
  
  bulkhead:
    instances:
      userService:
        max-concurrent-calls: 20
        max-wait-duration: 500ms
  
  ratelimiter:
    instances:
      userService:
        limit-for-period: 100
        limit-refresh-period: 1s
        timeout-duration: 0
  
  timelimiter:
    instances:
      userService:
        timeout-duration: 3s
        cancel-running-future: true
```

### 7.4 Resilience4j Programmatic Usage

```java
@Service
public class ResilientUserService {
    
    private final UserClient userClient;
    private final CircuitBreakerRegistry circuitBreakerRegistry;
    private final RetryRegistry retryRegistry;
    private final BulkheadRegistry bulkheadRegistry;
    
    public User getUserById(Long userId) {
        CircuitBreaker circuitBreaker = circuitBreakerRegistry.circuitBreaker("userService");
        Retry retry = retryRegistry.retry("userService");
        Bulkhead bulkhead = bulkheadRegistry.bulkhead("userService");
        
        Supplier<User> supplier = () -> userClient.getUserById(userId);
        
        // Decorate supplier with resilience patterns
        Supplier<User> decoratedSupplier = Decorators.ofSupplier(supplier)
            .withCircuitBreaker(circuitBreaker)
            .withRetry(retry)
            .withBulkhead(bulkhead)
            .withFallback(Arrays.asList(
                CallNotPermittedException.class,
                TimeoutException.class
            ), (ex) -> User.defaultUser(userId))
            .decorate();
        
        return decoratedSupplier.get();
    }
    
    // With metrics
    public void logCircuitBreakerMetrics() {
        CircuitBreaker cb = circuitBreakerRegistry.circuitBreaker("userService");
        CircuitBreaker.Metrics metrics = cb.getMetrics();
        
        log.info("Circuit Breaker State: {}", cb.getState());
        log.info("Failure Rate: {}%", metrics.getFailureRate());
        log.info("Successful Calls: {}", metrics.getNumberOfSuccessfulCalls());
        log.info("Failed Calls: {}", metrics.getNumberOfFailedCalls());
        log.info("Not Permitted Calls: {}", metrics.getNumberOfNotPermittedCalls());
    }
}
```

### 7.5 Retry Pattern

```java
@Service
public class OrderService {
    
    @Retry(name = "orderService", fallbackMethod = "createOrderFallback")
    public Order createOrder(OrderRequest request) {
        // May fail due to transient errors
        return orderClient.createOrder(request);
    }
    
    private Order createOrderFallback(OrderRequest request, Exception ex) {
        // Queue for later processing
        orderQueue.enqueue(request);
        return Order.pending(request);
    }
}
```

```yaml
resilience4j:
  retry:
    instances:
      orderService:
        max-attempts: 3
        wait-duration: 1s
        exponential-backoff-multiplier: 2
        # Wait: 1s, 2s, 4s
        retry-exceptions:
          - java.net.ConnectException
          - java.io.IOException
        ignore-exceptions:
          - com.example.ValidationException
```

### 7.6 Bulkhead Pattern

```java
// Semaphore Bulkhead - limits concurrent calls
@Bulkhead(name = "userService", type = Bulkhead.Type.SEMAPHORE)
public User getUserById(Long userId) {
    return userClient.getUserById(userId);
}

// Thread Pool Bulkhead - isolates in separate thread pool
@Bulkhead(name = "userService", type = Bulkhead.Type.THREADPOOL)
public CompletableFuture<User> getUserByIdAsync(Long userId) {
    return CompletableFuture.supplyAsync(() -> userClient.getUserById(userId));
}
```

```yaml
resilience4j:
  bulkhead:
    instances:
      userService:
        max-concurrent-calls: 10        # Max parallel calls
        max-wait-duration: 500ms        # Wait time when full
  
  thread-pool-bulkhead:
    instances:
      userService:
        max-thread-pool-size: 10
        core-thread-pool-size: 5
        queue-capacity: 100
        keep-alive-duration: 20ms
```

### 7.7 Rate Limiter

```java
@RateLimiter(name = "userService", fallbackMethod = "rateLimitFallback")
public List<User> getAllUsers() {
    return userClient.getAllUsers();
}

private List<User> rateLimitFallback(RequestNotPermitted ex) {
    throw new TooManyRequestsException("Rate limit exceeded. Try again later.");
}
```

---

## 8. Interview Questions - Part 1

### Basic Level

**Q1: What are microservices and how do they differ from monolithic architecture?**
> Microservices is an architectural style where an application is composed of small, independent services that communicate over APIs. Unlike monoliths (single deployable unit), microservices can be developed, deployed, and scaled independently. Each service owns its data and can use different technologies.

**Q2: What is service discovery and why is it needed?**
> Service discovery allows services to find each other dynamically without hardcoded addresses. As services scale up/down or move, their locations change. Discovery servers (like Eureka) maintain a registry of service instances and their locations, enabling dynamic routing.

**Q3: Explain the API Gateway pattern.**
> API Gateway is a single entry point for all clients. It handles cross-cutting concerns like:
> - Routing requests to appropriate services
> - Authentication/Authorization
> - Rate limiting
> - Load balancing
> - Response caching
> - Request/Response transformation

**Q4: What is Spring Cloud Config?**
> Spring Cloud Config provides centralized external configuration management. Services fetch their configuration from a Config Server (backed by Git, filesystem, or vault), enabling:
> - Externalized configuration
> - Environment-specific configs
> - Runtime refresh without restart
> - Encrypted sensitive data

**Q5: What is the difference between synchronous and asynchronous communication?**
> - **Synchronous**: Caller waits for response (REST, gRPC). Simple but creates tight coupling.
> - **Asynchronous**: Caller doesn't wait (messaging queues). Better decoupling and fault tolerance but eventually consistent.

### Intermediate Level

**Q6: Explain the Circuit Breaker pattern.**
> Circuit Breaker prevents cascading failures when a service is unavailable. States:
> - **CLOSED**: Requests pass through normally
> - **OPEN**: After failure threshold, requests fail fast
> - **HALF-OPEN**: Test requests to check recovery
> 
> Tools: Resilience4j, Hystrix (deprecated)

**Q7: What is the difference between Eureka and Consul?**

| Feature | Eureka | Consul |
|---------|--------|--------|
| Protocol | HTTP REST | HTTP + DNS |
| Health Check | Client heartbeat | Active health checks |
| KV Store | No | Yes |
| Service Mesh | No | Yes (Connect) |
| Best For | Spring ecosystem | Multi-platform |

**Q8: How does client-side load balancing work?**
> The client (service making requests) maintains a list of available instances from the service registry and decides which instance to call. Spring Cloud LoadBalancer implements strategies like Round Robin. Unlike server-side LB, there's no single point of failure.

**Q9: Explain the Database per Service pattern.**
> Each microservice has its own private database, ensuring loose coupling and independent scaling. Challenges include cross-service queries (use API composition) and distributed transactions (use Saga pattern).

**Q10: How do you handle distributed transactions?**
> Options:
> - **Saga Pattern**: Sequence of local transactions with compensating transactions
> - **Event Sourcing**: Store events instead of state
> - **Two-Phase Commit**: Coordinator-based (not recommended for microservices)

### Advanced Level

**Q11: Design a resilient inter-service communication strategy.**
```java
// Layered resilience approach:
@CircuitBreaker(name = "userService", fallbackMethod = "fallback")
@Retry(name = "userService")
@Bulkhead(name = "userService")
@TimeLimiter(name = "userService")
public CompletableFuture<User> getUser(Long id) {
    return userClient.getUserById(id);
}

// Configuration: Circuit opens at 50% failure, retries 3x with backoff,
// max 20 concurrent calls, 3s timeout
```

**Q12: How would you handle configuration refresh without downtime?**
> 1. Use Spring Cloud Config with Git backend
> 2. Enable `@RefreshScope` on beans with dynamic properties
> 3. Use Spring Cloud Bus to broadcast refresh events
> 4. POST to `/actuator/busrefresh` triggers all services to reload config
> 5. For critical configs, use feature flags or rolling deployments

**Q13: Explain the strangler fig pattern for monolith migration.**
> Incrementally extract functionality from monolith to microservices:
> 1. Add API Gateway in front of monolith
> 2. Extract one bounded context at a time
> 3. Route specific requests to new service
> 4. Keep facade in monolith for backward compatibility
> 5. Gradually retire monolith features

**Q14: How do you implement idempotency in microservices?**
```java
// Use idempotency keys
@PostMapping("/orders")
public Order createOrder(
        @RequestHeader("Idempotency-Key") String idempotencyKey,
        @RequestBody OrderRequest request) {
    // Check if already processed
    Order existing = cache.get(idempotencyKey);
    if (existing != null) return existing;
    
    // Process and cache result
    Order order = orderService.create(request);
    cache.put(idempotencyKey, order, 24, TimeUnit.HOURS);
    return order;
}
```

**Q15: Describe rate limiting strategies at the gateway level.**
> - **Token Bucket**: Tokens replenish at fixed rate, request takes a token
> - **Leaky Bucket**: Requests processed at constant rate, excess queued
> - **Fixed Window**: Count requests in fixed time windows
> - **Sliding Window**: Rolling time window for smoother limiting
> 
> Implementation: Spring Cloud Gateway + Redis for distributed rate limiting

---

## Quick Reference

### Spring Cloud Components

| Component | Purpose | Dependency |
|-----------|---------|------------|
| Eureka | Service Discovery | `spring-cloud-starter-netflix-eureka-*` |
| Config | Centralized Config | `spring-cloud-config-*` |
| Gateway | API Gateway | `spring-cloud-starter-gateway` |
| OpenFeign | Declarative REST | `spring-cloud-starter-openfeign` |
| LoadBalancer | Client-side LB | `spring-cloud-starter-loadbalancer` |
| Resilience4j | Fault Tolerance | `spring-cloud-starter-circuitbreaker-resilience4j` |
| Bus | Config Broadcast | `spring-cloud-starter-bus-*` |
| Stream | Messaging | `spring-cloud-starter-stream-*` |

### Common application.yml Properties

```yaml
# Eureka Client
eureka:
  client:
    service-url:
      defaultZone: http://localhost:8761/eureka/
  instance:
    prefer-ip-address: true

# Config Client  
spring:
  config:
    import: configserver:http://localhost:8888
  cloud:
    config:
      fail-fast: true

# Gateway Route
spring:
  cloud:
    gateway:
      routes:
        - id: service-route
          uri: lb://service-name
          predicates:
            - Path=/api/**
          filters:
            - StripPrefix=1

# Resilience4j
resilience4j:
  circuitbreaker:
    instances:
      myService:
        sliding-window-size: 10
        failure-rate-threshold: 50
```

---

*Part 2 covers: Distributed Tracing, Security, Event-Driven Architecture, Saga Pattern, Docker/Kubernetes Deployment*

*Last Updated: February 2026*
