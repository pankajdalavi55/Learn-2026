# Event-Driven Architecture in Spring Boot — Complete Guide

**Navigation:** [← Caching](Part-8_Caching-Spring-Boot-Complete-Guide.md)

> A comprehensive guide covering Event-Driven Architecture (EDA) fundamentals, patterns, messaging brokers, and production-ready Spring Boot implementations with real-world examples.  
> **Part 9 of the Spring Boot & Microservices Series**  
> **Prerequisites:** Spring Boot basics, Microservices fundamentals (Parts 5–6), Java fundamentals

---

## Table of Contents

1. [Introduction to Event-Driven Architecture](#1-introduction-to-event-driven-architecture)
2. [Core Concepts & Event Design](#2-core-concepts--event-design)
3. [EDA Patterns](#3-eda-patterns)
4. [Messaging Brokers Comparison](#4-messaging-brokers-comparison)
5. [In-Process Events with Spring](#5-in-process-events-with-spring)
6. [Apache Kafka with Spring Boot](#6-apache-kafka-with-spring-boot)
7. [RabbitMQ with Spring Boot](#7-rabbitmq-with-spring-boot)
8. [Spring Cloud Stream](#8-spring-cloud-stream)
9. [Reliable Event Publishing](#9-reliable-event-publishing)
10. [Event Consumption Best Practices](#10-event-consumption-best-practices)
11. [Saga Pattern — Choreography & Orchestration](#11-saga-pattern--choreography--orchestration)
12. [Complete E-Commerce Example](#12-complete-e-commerce-example)
13. [Event Sourcing & CQRS](#13-event-sourcing--cqrs)
14. [Testing Event-Driven Systems](#14-testing-event-driven-systems)
15. [Production Concerns & Observability](#15-production-concerns--observability)
16. [When NOT to Use Event-Driven Architecture](#16-when-not-to-use-event-driven-architecture)
17. [Interview Questions](#17-interview-questions)

---

## 1. Introduction to Event-Driven Architecture

### 1.1 What is Event-Driven Architecture?

**Event-Driven Architecture (EDA)** is an architectural style where services communicate by producing and consuming **events** — records of something that happened in the past — rather than making direct synchronous calls.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SYNCHRONOUS vs EVENT-DRIVEN                                   │
│                                                                                  │
│  SYNCHRONOUS (Request-Response):                                                 │
│  ┌─────────┐  POST /orders  ┌─────────┐  POST /reserve  ┌─────────┐            │
│  │  Order  │ ──────────────▶│Inventory│ ──────────────▶│ Payment │            │
│  │ Service │ ◀──────────────│ Service │ ◀──────────────│ Service │            │
│  └─────────┘   wait...      └─────────┘   wait...       └─────────┘            │
│                                                                                  │
│  • Tight coupling — Order knows Inventory and Payment APIs                      │
│  • Cascading failures — if Payment is down, Order fails                        │
│  • Slow — total latency = sum of all service calls                              │
│                                                                                  │
│  EVENT-DRIVEN (Asynchronous):                                                    │
│  ┌─────────┐  OrderCreated   ┌─────────────┐  StockReserved  ┌─────────┐        │
│  │  Order  │ ───────────────▶│ Event Broker│ ───────────────▶│ Payment │        │
│  │ Service │                 │   (Kafka)   │                 │ Service │        │
│  └─────────┘                 └──────┬──────┘                 └─────────┘        │
│       │                             │                                            │
│       │ returns 201 immediately     │ OrderCreated also consumed by:            │
│       ▼                             ├──▶ Inventory Service                      │
│  "Order accepted"                   ├──▶ Analytics Service                      │
│                                     └──▶ Notification Service                    │
│                                                                                  │
│  • Loose coupling — Order doesn't know who listens                              │
│  • Resilient — broker buffers events if consumers are down                      │
│  • Fast response — Order returns immediately after publishing                   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 What is an Event?

An **event** is an immutable record of a state change that already occurred:

| Property | Description | Example |
|----------|-------------|---------|
| **Past tense** | Describes what happened | `OrderCreated`, not `CreateOrder` |
| **Immutable** | Cannot be changed after creation | Events are append-only |
| **Self-contained** | Carries enough context to act | Order ID, items, amount |
| **Timestamped** | When it happened | `2024-03-15T10:30:00Z` |
| **Identifiable** | Unique event ID for deduplication | `evt-uuid-1234` |

```java
// Event — describes what HAPPENED (past tense)
public record OrderCreatedEvent(
    String eventId,
    String eventType,      // "ORDER_CREATED"
    Instant timestamp,
    Long orderId,
    String userId,
    List<OrderItem> items,
    BigDecimal totalAmount
) {}

// Command — describes what SHOULD happen (imperative) — NOT an event
public record CreateOrderCommand(
    String userId,
    List<OrderItem> items
) {}
```

### 1.3 Why Use Event-Driven Architecture?

| Benefit | Description | Real-World Impact |
|---------|-------------|-------------------|
| **Loose Coupling** | Producers don't know consumers | Add analytics without touching Order service |
| **Scalability** | Services scale independently | Scale notification workers separately |
| **Resilience** | Broker buffers during outages | Payment down? Events queue until recovery |
| **Extensibility** | New consumers subscribe to topics | Add fraud detection as new listener |
| **Audit Trail** | Events are a natural log of changes | Replay events for debugging/compliance |
| **Async Processing** | Non-blocking user experience | Order API returns in 50ms, not 2 seconds |

### 1.4 When EDA Fits Microservices

EDA shines when:

- Multiple services need to react to the same business event
- The caller doesn't need an immediate response from downstream services
- You need eventual consistency across bounded contexts
- You want to decouple teams and deployment cycles
- High throughput event streams are involved (orders, clicks, IoT)

> **Cross-reference:** Part 6 covers EDA at a high level within microservices. This guide goes deep into implementation, patterns, and production concerns.

---

## 2. Core Concepts & Event Design

### 2.1 EDA Building Blocks

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    EVENT-DRIVEN ARCHITECTURE COMPONENTS                            │
│                                                                                  │
│  ┌──────────────┐         ┌──────────────────┐         ┌──────────────┐       │
│  │   PRODUCER   │ publish │   EVENT BROKER   │ deliver │   CONSUMER   │       │
│  │ (Publisher)  │ ──────▶ │  (Kafka/Rabbit)  │ ──────▶ │ (Subscriber) │       │
│  └──────────────┘         └──────────────────┘         └──────────────┘       │
│                                                                                  │
│  Producer:                  Broker:                       Consumer:              │
│  • Order Service            • Kafka Topic                 • Inventory Service  │
│  • User Service             • RabbitMQ Exchange/Queue     • Payment Service    │
│  • Emits events             • Stores & routes events      • Reacts to events   │
│                                                                                  │
│  Optional Components:                                                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ Schema      │  │ Event       │  │ Dead Letter │  │ Stream      │            │
│  │ Registry    │  │ Store       │  │ Queue (DLQ) │  │ Processor   │            │
│  │ (Avro/JSON) │  │ (Event Src) │  │             │  │ (Kafka Str) │            │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Event Envelope Pattern

Wrap business payload in a standard envelope for cross-cutting metadata:

```java
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class EventEnvelope<T> {
    // Metadata
    private String eventId;          // UUID — for idempotency
    private String eventType;        // ORDER_CREATED, PAYMENT_FAILED
    private String aggregateType;    // Order, Payment
    private String aggregateId;      // order-123
    private Instant timestamp;
    private String source;           // order-service
    private String correlationId;    // trace across services
    private int version;             // schema version

    // Payload
    private T data;
}

// Usage
EventEnvelope<OrderCreatedPayload> envelope = EventEnvelope.<OrderCreatedPayload>builder()
    .eventId(UUID.randomUUID().toString())
    .eventType("ORDER_CREATED")
    .aggregateType("Order")
    .aggregateId(order.getId().toString())
    .timestamp(Instant.now())
    .source("order-service")
    .correlationId(MDC.get("traceId"))
    .version(1)
    .data(OrderCreatedPayload.from(order))
    .build();
```

### 2.3 Event Design Principles

| Principle | Good | Bad |
|-----------|------|-----|
| **Past tense naming** | `OrderCreated`, `PaymentProcessed` | `CreateOrder`, `ProcessPayment` |
| **Self-describing** | Include all data consumer needs | Force consumer to call back for details |
| **Small & focused** | One event = one thing happened | `OrderEverythingChanged` |
| **Versioned** | `version: 2` in envelope | Breaking changes without versioning |
| **Idempotent-friendly** | Include `eventId` for dedup | No unique identifier |
| **Partition key** | Use `orderId` as Kafka key | Random key breaks ordering |

### 2.4 Delivery Semantics

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    MESSAGE DELIVERY GUARANTEES                                     │
│                                                                                  │
│  AT-MOST-ONCE:     Message may be lost, never duplicated                         │
│  ─────────────     Fast, but risky for critical business events                 │
│                                                                                  │
│  AT-LEAST-ONCE:    Message always delivered, may be duplicated                  │
│  ──────────────    Most common in production — requires idempotent consumers    │
│                                                                                  │
│  EXACTLY-ONCE:     Message delivered exactly once                               │
│  ─────────────     Hardest to achieve — Kafka transactions + idempotent producer│
│                                                                                  │
│  Production Reality:                                                             │
│  • Design for AT-LEAST-ONCE + idempotent consumers                                │
│  • Use Outbox pattern for reliable publishing                                     │
│  • Store processed eventIds to prevent duplicate processing                     │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. EDA Patterns

### 3.1 Pattern Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    EVENT-DRIVEN ARCHITECTURE PATTERNS                            │
│                                                                                  │
│  1. EVENT NOTIFICATION                                                           │
│  ─────────────────────                                                           │
│  ┌─────────┐   OrderCreated    ┌─────────────┐                                  │
│  │  Order  │ ────────────────▶ │  Inventory  │  "Something happened"            │
│  │ Service │   (minimal data)  │   Service   │  Consumer fetches details        │
│  └─────────┘                   └─────────────┘  via API if needed               │
│                                                                                  │
│  2. EVENT-CARRIED STATE TRANSFER                                                 │
│  ───────────────────────────────                                                 │
│  ┌─────────┐   OrderCreated    ┌─────────────┐                                  │
│  │  Order  │ ────────────────▶ │  Inventory  │  "Here's everything you need"    │
│  │ Service │   (full payload)  │   Service   │  No callback required            │
│  └─────────┘                   └─────────────┘                                  │
│                                                                                  │
│  3. EVENT SOURCING                                                               │
│  ─────────────────                                                               │
│  ┌─────────┐                   ┌─────────────┐                                  │
│  │  Order  │ ────────────────▶ │   Event     │  Store events, not current state │
│  │ Service │   OrderCreated    │   Store     │  Rebuild state by replaying       │
│  └─────────┘   ItemAdded       └──────┬──────┘                                  │
│                OrderCompleted         │                                          │
│                                       ▼                                          │
│                               ┌─────────────┐                                    │
│                               │ Projections │  Materialized read views           │
│                               └─────────────┘                                    │
│                                                                                  │
│  4. CQRS (Command Query Responsibility Segregation)                              │
│  ───────────────────────────────────────────────                                 │
│          Commands                      Queries                                   │
│              │                            │                                      │
│              ▼                            ▼                                      │
│       ┌─────────────┐            ┌─────────────┐                                │
│       │   Write     │   Events   │    Read     │                                │
│       │   Model     │ ─────────▶ │   Model     │                                │
│       │ (normalized)│            │(denormalized)│                                │
│       └─────────────┘            └─────────────┘                                │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Event Notification vs Event-Carried State Transfer

| Aspect | Event Notification | Event-Carried State Transfer |
|--------|-------------------|---------------------------|
| **Payload size** | Minimal (IDs only) | Full entity data |
| **Coupling** | Consumer calls producer API | Fully decoupled |
| **Consistency** | Consumer sees latest on fetch | Point-in-time snapshot |
| **Network** | Extra API calls | Single message |
| **Use when** | Rare consumers, large entities | Multiple consumers, small payloads |

```java
// Event Notification — minimal payload
public record OrderCreatedNotification(Long orderId, String userId) {}

// Consumer must fetch details
@KafkaListener(topics = "order-notifications")
public void handle(OrderCreatedNotification event) {
    OrderDetails order = orderApiClient.getOrder(event.orderId()); // extra call
    inventoryService.reserve(order.items());
}

// Event-Carried State Transfer — full payload
public record OrderCreatedEvent(
    Long orderId,
    String userId,
    List<OrderItem> items,
    BigDecimal totalAmount,
    String shippingAddress
) {}

// Consumer acts immediately — no callback
@KafkaListener(topics = "order-events")
public void handle(OrderCreatedEvent event) {
    inventoryService.reserve(event.items()); // self-contained
}
```

### 3.3 Publish-Subscribe Pattern

```
                    ┌─────────────────────────────────┐
                    │         Kafka Topic: orders      │
                    └───────────────┬─────────────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              │                     │                     │
              ▼                     ▼                     ▼
     ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
     │ Inventory Svc   │  │ Notification Svc│  │ Analytics Svc   │
     │ group: inventory│  │ group: notify   │  │ group: analytics│
     └─────────────────┘  └─────────────────┘  └─────────────────┘

     Each consumer group receives ALL messages independently.
     Adding Analytics doesn't require changing Order Service.
```

---

## 4. Messaging Brokers Comparison

### 4.1 Kafka vs RabbitMQ vs Others

| Aspect | Apache Kafka | RabbitMQ | AWS SQS/SNS |
|--------|-------------|----------|-------------|
| **Model** | Distributed commit log | Message queue (AMQP) | Managed queue/topic |
| **Message retention** | Days/weeks (configurable) | Until consumed | 14 days max |
| **Replay** | Yes — consumers re-read | No — message deleted | Limited |
| **Throughput** | Millions/sec | Thousands/sec | High (managed) |
| **Ordering** | Per partition | Per queue | FIFO queues only |
| **Consumer model** | Pull (consumer polls) | Push (broker delivers) | Pull |
| **Routing** | Topic-based | Exchanges, bindings, routing keys | Topic/queue based |
| **Best for** | Event streaming, analytics, CDC | Task queues, RPC, routing | Cloud-native, serverless |
| **Complexity** | Higher ops overhead | Moderate | Low (fully managed) |

### 4.2 Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CHOOSING A MESSAGE BROKER                                     │
│                                                                                  │
│  Use KAFKA when:                                                                 │
│  ✓ High throughput (100K+ events/sec)                                           │
│  ✓ Event replay / audit trail required                                          │
│  ✓ Multiple consumers need same events                                          │
│  ✓ Stream processing (Kafka Streams, Flink)                                     │
│  ✓ Event sourcing / CDC pipelines                                               │
│                                                                                  │
│  Use RABBITMQ when:                                                              │
│  ✓ Complex routing (topic, fanout, headers exchanges)                           │
│  ✓ Request-reply / RPC patterns                                                 │
│  ✓ Smaller scale (< 50K messages/sec)                                           │
│  ✓ Priority queues, TTL per message                                             │
│  ✓ Simpler operational requirements                                             │
│                                                                                  │
│  Use BOTH when:                                                                  │
│  ✓ Kafka for event streaming backbone                                           │
│  ✓ RabbitMQ for task queues and RPC between services                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 Kafka Core Concepts (Quick Reference)

| Concept | Description |
|---------|-------------|
| **Topic** | Named category for events (e.g., `order-events`) |
| **Partition** | Ordered, immutable log within a topic — unit of parallelism |
| **Offset** | Sequential ID of message within a partition |
| **Producer** | Writes events to topics |
| **Consumer Group** | Group of consumers sharing work — each partition consumed by one member |
| **Broker** | Kafka server storing partitions |
| **Replication Factor** | Copies of each partition across brokers for fault tolerance |

> **Deep dive:** See the [Apache Kafka Complete Learning Guide](../1.%20System%20Design/Kafka/01-Apache-Kafka-Complete-Learning-Guide.md) for internals.

---

## 5. In-Process Events with Spring

Before distributed messaging, Spring provides **in-process event publishing** — useful within a monolith or single service.

### 5.1 Spring Application Events

```java
// 1. Define event (plain POJO — no extends required since Spring 4.2)
public record OrderPlacedEvent(
    Long orderId,
    String customerEmail,
    BigDecimal totalAmount,
    List<OrderItem> items
) {}

// 2. Publish event
@Service
@RequiredArgsConstructor
public class OrderService {

    private final OrderRepository orderRepository;
    private final ApplicationEventPublisher eventPublisher;

    @Transactional
    public Order placeOrder(CreateOrderRequest request) {
        Order order = orderRepository.save(mapToEntity(request));
        eventPublisher.publishEvent(new OrderPlacedEvent(
            order.getId(),
            order.getCustomerEmail(),
            order.getTotalAmount(),
            order.getItems()
        ));
        return order;
    }
}

// 3. Listen synchronously (default)
@Component
@Slf4j
public class NotificationListener {

    @EventListener
    public void onOrderPlaced(OrderPlacedEvent event) {
        log.info("Sending confirmation email for order {}", event.orderId());
        emailService.sendConfirmation(event.customerEmail(), event.orderId());
    }
}

// 4. Listen asynchronously
@Component
@Slf4j
public class AnalyticsListener {

    @Async
    @EventListener
    public void onOrderPlaced(OrderPlacedEvent event) {
        log.info("Recording analytics for order {}", event.orderId());
        analyticsService.recordPurchase(event);
    }
}

// 5. Conditional listener — only for large orders
@EventListener(condition = "#event.totalAmount.compareTo(T(java.math.BigDecimal).valueOf(1000)) > 0")
public void onLargeOrder(OrderPlacedEvent event) {
    fraudService.flagForReview(event.orderId());
}

// 6. Transaction-aware — fire only after commit
@TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
public void onOrderCommitted(OrderPlacedEvent event) {
    // Safe to publish to Kafka — DB transaction already committed
    kafkaTemplate.send("orders", event.orderId().toString(), event);
}
```

### 5.2 In-Process vs Distributed Events

| Aspect | Spring Application Events | Kafka/RabbitMQ Events |
|--------|--------------------------|----------------------|
| **Scope** | Single JVM / application | Cross-service, cross-network |
| **Durability** | Lost if JVM crashes | Persisted in broker |
| **Ordering** | Guaranteed in-process | Per-partition (Kafka) |
| **Use case** | Decouple modules in monolith | Microservice communication |
| **Latency** | Microseconds | Milliseconds |

---

## 6. Apache Kafka with Spring Boot

### 6.1 Project Setup

```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka</artifactId>
</dependency>
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka-test</artifactId>
    <scope>test</scope>
</dependency>
```

```yaml
# application.yml
spring:
  kafka:
    bootstrap-servers: localhost:9092
    producer:
      key-serializer: org.apache.kafka.common.serialization.StringSerializer
      value-serializer: org.springframework.kafka.support.serializer.JsonSerializer
      acks: all
      retries: 3
      properties:
        enable.idempotence: true
    consumer:
      group-id: ${spring.application.name}
      auto-offset-reset: earliest
      key-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      value-deserializer: org.springframework.kafka.support.serializer.JsonDeserializer
      properties:
        spring.json.trusted.packages: com.example.*
      enable-auto-commit: false
    listener:
      ack-mode: manual
```

### 6.2 Kafka Producer — Order Service

```java
@Configuration
@EnableKafka
public class KafkaProducerConfig {

    @Bean
    public ProducerFactory<String, Object> producerFactory(
            @Value("${spring.kafka.bootstrap-servers}") String bootstrapServers) {
        Map<String, Object> config = new HashMap<>();
        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        config.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        config.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, JsonSerializer.class);
        config.put(ProducerConfig.ACKS_CONFIG, "all");
        config.put(ProducerConfig.RETRIES_CONFIG, 3);
        config.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
        return new DefaultKafkaProducerFactory<>(config);
    }

    @Bean
    public KafkaTemplate<String, Object> kafkaTemplate(
            ProducerFactory<String, Object> producerFactory) {
        return new KafkaTemplate<>(producerFactory);
    }
}

@Service
@Slf4j
@RequiredArgsConstructor
public class OrderEventPublisher {

    private static final String TOPIC = "order-events";

    private final KafkaTemplate<String, Object> kafkaTemplate;

    public void publishOrderCreated(Order order) {
        OrderCreatedEvent event = OrderCreatedEvent.from(order);

        // Use orderId as key — ensures all events for same order go to same partition
        kafkaTemplate.send(TOPIC, order.getId().toString(), event)
            .whenComplete((result, ex) -> {
                if (ex == null) {
                    log.info("Published OrderCreatedEvent orderId={} offset={}",
                        order.getId(),
                        result.getRecordMetadata().offset());
                } else {
                    log.error("Failed to publish OrderCreatedEvent orderId={}",
                        order.getId(), ex);
                    throw new EventPublishException("Kafka publish failed", ex);
                }
            });
    }

    public void publishOrderCancelled(Long orderId, String reason) {
        OrderCancelledEvent event = new OrderCancelledEvent(
            UUID.randomUUID().toString(),
            orderId,
            reason,
            Instant.now()
        );
        kafkaTemplate.send(TOPIC, orderId.toString(), event);
    }
}
```

### 6.3 Kafka Consumer — Inventory Service

```java
@Configuration
@EnableKafka
public class KafkaConsumerConfig {

    @Bean
    public ConsumerFactory<String, OrderCreatedEvent> consumerFactory(
            @Value("${spring.kafka.bootstrap-servers}") String bootstrapServers) {
        Map<String, Object> config = new HashMap<>();
        config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        config.put(ConsumerConfig.GROUP_ID_CONFIG, "inventory-service");
        config.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
        config.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, JsonDeserializer.class);
        config.put(JsonDeserializer.TRUSTED_PACKAGES, "com.example.*");
        config.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, false);
        return new DefaultKafkaConsumerFactory<>(config);
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, OrderCreatedEvent>
            kafkaListenerContainerFactory(ConsumerFactory<String, OrderCreatedEvent> consumerFactory) {
        ConcurrentKafkaListenerContainerFactory<String, OrderCreatedEvent> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory);
        factory.setConcurrency(3);
        factory.getContainerProperties().setAckMode(ContainerProperties.AckMode.MANUAL);
        factory.setCommonErrorHandler(new DefaultErrorHandler(
            new FixedBackOff(1000L, 3L) // 3 retries, 1 second apart
        ));
        return factory;
    }
}

@Service
@Slf4j
@RequiredArgsConstructor
public class OrderEventConsumer {

    private final InventoryService inventoryService;
    private final ProcessedEventRepository processedEventRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @KafkaListener(
        topics = "order-events",
        groupId = "inventory-service",
        containerFactory = "kafkaListenerContainerFactory"
    )
    public void handleOrderCreated(
            @Payload OrderCreatedEvent event,
            @Header(KafkaHeaders.RECEIVED_KEY) String key,
            @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
            @Header(KafkaHeaders.OFFSET) long offset,
            Acknowledgment acknowledgment) {

        log.info("Received OrderCreatedEvent orderId={} partition={} offset={}",
            event.orderId(), partition, offset);

        // Idempotency check
        if (processedEventRepository.existsByEventId(event.eventId())) {
            log.warn("Duplicate event skipped: {}", event.eventId());
            acknowledgment.acknowledge();
            return;
        }

        try {
            inventoryService.reserveStock(event.items());

            processedEventRepository.save(new ProcessedEvent(event.eventId(), Instant.now()));

            // Publish success event for saga continuation
            kafkaTemplate.send("inventory-events", key,
                new StockReservedEvent(event.orderId(), event.items()));

            acknowledgment.acknowledge();

        } catch (InsufficientStockException e) {
            log.error("Insufficient stock for orderId={}", event.orderId());

            kafkaTemplate.send("inventory-events", key,
                new StockReserveFailedEvent(event.orderId(), e.getMessage()));

            acknowledgment.acknowledge(); // Don't retry — business failure
        }
    }
}
```

### 6.4 Docker Compose for Local Kafka

```yaml
# docker-compose.yml
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on: [zookeeper]
    ports: ["9092:9092"]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    ports: ["8080:8080"]
    environment:
      KAFKA_CLUSTERS_0_NAME: local
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
    depends_on: [kafka]
```

---

## 7. RabbitMQ with Spring Boot

### 7.1 Project Setup

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-amqp</artifactId>
</dependency>
```

```yaml
spring:
  rabbitmq:
    host: localhost
    port: 5672
    username: guest
    password: guest
    listener:
      simple:
        acknowledge-mode: manual
        retry:
          enabled: true
          max-attempts: 3
          initial-interval: 1000ms
```

### 7.2 Exchange, Queue & Binding Configuration

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    RABBITMQ TOPOLOGY — ORDER FLOW                                │
│                                                                                  │
│  ┌──────────────┐    order.created     ┌──────────────────┐                   │
│  │ order.exchange│ ──────────────────▶ │ order.created.q  │ ──▶ Inventory Svc │
│  │  (Topic)      │                     └──────────────────┘                   │
│  └──────┬───────┘                                                                │
│         │ order.cancelled                                                        │
│         └──────────────────────────▶ ┌──────────────────────┐                   │
│                                      │ order.cancelled.q    │ ──▶ Payment Svc   │
│                                      └──────────────────────┘                   │
│                                                                                  │
│  Dead Letter Exchange (DLX):                                                     │
│  order.created.q ──(fail after retries)──▶ order.dlx ──▶ order.created.dlq       │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
@Configuration
public class RabbitMQConfig {

    public static final String ORDER_EXCHANGE = "order.exchange";
    public static final String ORDER_DLX = "order.dlx";
    public static final String ORDER_CREATED_QUEUE = "order.created.queue";
    public static final String ORDER_CREATED_DLQ = "order.created.dlq";
    public static final String ORDER_CREATED_KEY = "order.created";
    public static final String STOCK_RESERVED_KEY = "stock.reserved";

    @Bean
    public TopicExchange orderExchange() {
        return new TopicExchange(ORDER_EXCHANGE, true, false);
    }

    @Bean
    public DirectExchange deadLetterExchange() {
        return new DirectExchange(ORDER_DLX, true, false);
    }

    @Bean
    public Queue orderCreatedQueue() {
        return QueueBuilder.durable(ORDER_CREATED_QUEUE)
            .withArgument("x-dead-letter-exchange", ORDER_DLX)
            .withArgument("x-dead-letter-routing-key", ORDER_CREATED_DLQ)
            .withArgument("x-message-ttl", 86400000) // 24 hours
            .build();
    }

    @Bean
    public Queue orderCreatedDlq() {
        return QueueBuilder.durable(ORDER_CREATED_DLQ).build();
    }

    @Bean
    public Binding orderCreatedBinding() {
        return BindingBuilder.bind(orderCreatedQueue())
            .to(orderExchange())
            .with(ORDER_CREATED_KEY);
    }

    @Bean
    public Binding orderCreatedDlqBinding() {
        return BindingBuilder.bind(orderCreatedDlq())
            .to(deadLetterExchange())
            .with(ORDER_CREATED_DLQ);
    }

    @Bean
    public Jackson2JsonMessageConverter messageConverter() {
        return new Jackson2JsonMessageConverter();
    }

    @Bean
    public RabbitTemplate rabbitTemplate(ConnectionFactory factory,
            Jackson2JsonMessageConverter converter) {
        RabbitTemplate template = new RabbitTemplate(factory);
        template.setMessageConverter(converter);
        template.setConfirmCallback((correlation, ack, reason) -> {
            if (!ack) {
                log.error("Message not acknowledged: {}", reason);
            }
        });
        return template;
    }
}
```

### 7.3 Publisher & Consumer

```java
@Service
@Slf4j
@RequiredArgsConstructor
public class OrderRabbitPublisher {

    private final RabbitTemplate rabbitTemplate;

    public void publishOrderCreated(Order order) {
        OrderCreatedEvent event = OrderCreatedEvent.from(order);

        rabbitTemplate.convertAndSend(
            RabbitMQConfig.ORDER_EXCHANGE,
            RabbitMQConfig.ORDER_CREATED_KEY,
            event,
            message -> {
                message.getMessageProperties().setMessageId(event.eventId());
                message.getMessageProperties().setCorrelationId(
                    MDC.get("traceId"));
                return message;
            }
        );

        log.info("Published OrderCreatedEvent to RabbitMQ orderId={}", order.getId());
    }
}

@Service
@Slf4j
@RequiredArgsConstructor
public class InventoryRabbitConsumer {

    private final InventoryService inventoryService;
    private final ProcessedEventRepository processedEventRepository;

    @RabbitListener(queues = RabbitMQConfig.ORDER_CREATED_QUEUE)
    public void handleOrderCreated(
            OrderCreatedEvent event,
            Channel channel,
            @Header(AmqpHeaders.DELIVERY_TAG) long deliveryTag,
            @Header(value = AmqpHeaders.MESSAGE_ID, required = false) String messageId) {

        log.info("Received OrderCreatedEvent orderId={}", event.orderId());

        if (processedEventRepository.existsByEventId(event.eventId())) {
            log.warn("Duplicate event skipped: {}", event.eventId());
            acknowledge(channel, deliveryTag);
            return;
        }

        try {
            inventoryService.reserveStock(event.items());
            processedEventRepository.save(new ProcessedEvent(event.eventId(), Instant.now()));
            acknowledge(channel, deliveryTag);

        } catch (InsufficientStockException e) {
            log.error("Insufficient stock — sending to DLQ orderId={}", event.orderId());
            // Reject without requeue → goes to DLQ
            rejectToDlq(channel, deliveryTag);
        } catch (Exception e) {
            log.error("Transient error — will retry orderId={}", event.orderId(), e);
            rejectForRetry(channel, deliveryTag);
        }
    }

    private void acknowledge(Channel channel, long deliveryTag) {
        try { channel.basicAck(deliveryTag, false); }
        catch (IOException e) { throw new AmqpException(e); }
    }

    private void rejectToDlq(Channel channel, long deliveryTag) {
        try { channel.basicReject(deliveryTag, false); }
        catch (IOException e) { throw new AmqpException(e); }
    }

    private void rejectForRetry(Channel channel, long deliveryTag) {
        try { channel.basicNack(deliveryTag, false, true); }
        catch (IOException e) { throw new AmqpException(e); }
    }
}
```

---

## 8. Spring Cloud Stream

Spring Cloud Stream provides a **broker-agnostic abstraction** over Kafka, RabbitMQ, and others using functional programming style.

### 8.1 Dependencies & Configuration

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-stream-kafka</artifactId>
</dependency>
```

```yaml
spring:
  cloud:
    stream:
      function:
        definition: orderProcessor;paymentProcessor
      bindings:
        orderProcessor-in-0:
          destination: order-events
          group: inventory-service
        orderProcessor-out-0:
          destination: inventory-events
        paymentProcessor-in-0:
          destination: inventory-events
          group: payment-service
      kafka:
        binder:
          brokers: localhost:9092
        bindings:
          orderProcessor-in-0:
            consumer:
              enable-dlq: true
              dlq-name: order-events.dlq
```

### 8.2 Functional Style (Spring Cloud Stream 4.x)

```java
@SpringBootApplication
public class InventoryServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(InventoryServiceApplication.class, args);
    }

    // Consumer: process order events
    @Bean
    public Consumer<OrderCreatedEvent> orderProcessor(
            InventoryService inventoryService,
            StreamBridge streamBridge) {
        return event -> {
            log.info("Processing order: {}", event.orderId());
            try {
                inventoryService.reserveStock(event.items());
                streamBridge.send("inventory-events-out-0",
                    StockReservedEvent.from(event));
            } catch (InsufficientStockException e) {
                streamBridge.send("inventory-events-out-0",
                    StockReserveFailedEvent.from(event, e.getMessage()));
            }
        };
    }
}

// StreamBridge — imperative publishing from anywhere
@Service
@RequiredArgsConstructor
public class OrderService {

    private final OrderRepository orderRepository;
    private final StreamBridge streamBridge;

    @Transactional
    public Order createOrder(CreateOrderRequest request) {
        Order order = orderRepository.save(mapToEntity(request));

        boolean sent = streamBridge.send("order-events-out-0",
            OrderCreatedEvent.from(order));

        if (!sent) {
            throw new EventPublishException("Failed to send to order-events");
        }
        return order;
    }
}
```

### 8.3 When to Use Spring Cloud Stream

| Use Spring Cloud Stream | Use Raw Kafka/RabbitMQ Client |
|------------------------|------------------------------|
| Multi-broker portability needed | Kafka-only, need fine-grained control |
| Functional reactive pipelines | Complex consumer offset management |
| Built-in DLQ, partitioning config | Custom retry/backoff logic |
| Rapid prototyping | Maximum performance tuning |

---

## 9. Reliable Event Publishing

### 9.1 The Dual-Write Problem

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    THE DUAL-WRITE PROBLEM                                        │
│                                                                                  │
│  WITHOUT OUTBOX:                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐                │
│  │  @Transactional                                              │                │
│  │  1. Save Order to DB          ✓ SUCCESS                     │                │
│  │  2. Publish to Kafka          ✗ NETWORK FAILURE             │                │
│  └─────────────────────────────────────────────────────────────┘                │
│                                                                                  │
│  Result: Order exists in DB but no event published → inconsistent state         │
│                                                                                  │
│  OR:                                                                             │
│  1. Publish to Kafka            ✓ SUCCESS                                     │
│  2. Save Order to DB            ✗ DB FAILURE                                   │
│  Result: Event published but order doesn't exist → phantom event                │
│                                                                                  │
│  SOLUTION: Transactional Outbox Pattern                                          │
│  Save order AND event in SAME database transaction.                             │
│  Separate process publishes events from outbox table.                           │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Transactional Outbox Pattern

```java
// Outbox entity
@Entity
@Table(name = "outbox_events")
@Data
public class OutboxEvent {
    @Id
    private String id;

    private String aggregateType;   // "Order"
    private String aggregateId;     // "12345"
    private String eventType;       // "ORDER_CREATED"
    private String topic;           // "order-events"

    @Column(columnDefinition = "TEXT")
    private String payload;         // JSON serialized event

    private Instant createdAt;
    private boolean published;
    private Instant publishedAt;
}

// Order service — atomic DB + outbox write
@Service
@RequiredArgsConstructor
public class OrderService {

    private final OrderRepository orderRepository;
    private final OutboxRepository outboxRepository;
    private final ObjectMapper objectMapper;

    @Transactional
    public Order createOrder(CreateOrderRequest request) {
        Order order = orderRepository.save(mapToEntity(request));

        OrderCreatedEvent event = OrderCreatedEvent.from(order);
        OutboxEvent outbox = new OutboxEvent();
        outbox.setId(UUID.randomUUID().toString());
        outbox.setAggregateType("Order");
        outbox.setAggregateId(order.getId().toString());
        outbox.setEventType("ORDER_CREATED");
        outbox.setTopic("order-events");
        outbox.setPayload(objectMapper.writeValueAsString(event));
        outbox.setCreatedAt(Instant.now());
        outbox.setPublished(false);

        outboxRepository.save(outbox);
        return order;
    }
}

// Outbox poller — publishes unpublished events
@Service
@Slf4j
@RequiredArgsConstructor
public class OutboxPoller {

    private final OutboxRepository outboxRepository;
    private final KafkaTemplate<String, String> kafkaTemplate;

    @Scheduled(fixedDelay = 1000)
    @Transactional
    public void pollAndPublish() {
        List<OutboxEvent> pending = outboxRepository
            .findTop100ByPublishedFalseOrderByCreatedAtAsc();

        for (OutboxEvent outbox : pending) {
            try {
                kafkaTemplate.send(
                    outbox.getTopic(),
                    outbox.getAggregateId(),
                    outbox.getPayload()
                ).get(5, TimeUnit.SECONDS);

                outbox.setPublished(true);
                outbox.setPublishedAt(Instant.now());
                outboxRepository.save(outbox);

                log.info("Published outbox event id={} type={}",
                    outbox.getId(), outbox.getEventType());

            } catch (Exception e) {
                log.error("Failed to publish outbox event id={}", outbox.getId(), e);
                // Will retry on next poll
            }
        }
    }
}
```

### 9.3 Debezium CDC — Alternative to Polling

Instead of polling the outbox table, **Debezium** reads database transaction logs (WAL/binlog) and streams changes to Kafka automatically.

```json
// POST http://localhost:8083/connectors
{
  "name": "order-outbox-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "order_service",
    "database.password": "password",
    "database.dbname": "orders",
    "database.server.name": "order-service",
    "table.include.list": "public.outbox_events",
    "transforms": "outbox",
    "transforms.outbox.type": "io.debezium.transforms.outbox.EventRouter",
    "transforms.outbox.table.field.event.key": "aggregate_id",
    "transforms.outbox.table.field.event.type": "event_type",
    "transforms.outbox.table.field.event.payload": "payload",
    "transforms.outbox.route.by.field": "aggregate_type",
    "transforms.outbox.route.topic.replacement": "${routedByValue}-events"
  }
}
```

| Approach | Pros | Cons |
|----------|------|------|
| **Scheduled Poller** | Simple, no extra infra | Polling delay, DB load |
| **Debezium CDC** | Near real-time, no polling | Requires Kafka Connect, ops complexity |
| **Transactional Kafka** | True exactly-once | Kafka 0.11+, complex setup |

---

## 10. Event Consumption Best Practices

### 10.1 Idempotent Consumers

At-least-once delivery means consumers **will** see duplicates. Design every consumer to be idempotent:

```java
@Entity
@Table(name = "processed_events", uniqueConstraints = {
    @UniqueConstraint(columnNames = {"event_id", "consumer_group"})
})
@Data
public class ProcessedEvent {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String eventId;
    private String consumerGroup;
    private Instant processedAt;
}

@Repository
public interface ProcessedEventRepository extends JpaRepository<ProcessedEvent, Long> {
    boolean existsByEventIdAndConsumerGroup(String eventId, String consumerGroup);
}

// Idempotent consumer pattern
public void processEvent(OrderCreatedEvent event, String consumerGroup) {
    if (processedEventRepository.existsByEventIdAndConsumerGroup(
            event.eventId(), consumerGroup)) {
        log.warn("Duplicate event {} — skipping", event.eventId());
        return;
    }

    // Business logic
    inventoryService.reserveStock(event.items());

    // Mark as processed (same transaction as business logic)
    processedEventRepository.save(new ProcessedEvent(
        event.eventId(), consumerGroup, Instant.now()));
}
```

### 10.2 Event Ordering

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    EVENT ORDERING IN KAFKA                                       │
│                                                                                  │
│  Ordering guaranteed ONLY within a partition.                                     │
│                                                                                  │
│  Topic: order-events (3 partitions)                                              │
│                                                                                  │
│  Partition 0: [OrderCreated(1)] [ItemAdded(1)] [OrderPaid(1)]  ← ordered        │
│  Partition 1: [OrderCreated(2)] [OrderPaid(2)]               ← ordered        │
│  Partition 2: [OrderCreated(3)] [ItemAdded(3)] [OrderPaid(3)]  ← ordered        │
│                                                                                  │
│  Across partitions: NO ordering guarantee                                        │
│                                                                                  │
│  RULE: Use aggregate ID as message key                                           │
│  kafkaTemplate.send("order-events", orderId.toString(), event);                 │
│  → All events for order 123 always land in same partition                      │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Sequence number validation for strict ordering
public void processWithVersioning(OrderStatusChangedEvent event) {
    Order order = orderRepository.findById(event.orderId()).orElseThrow();

    if (event.sequenceNumber() <= order.getLastEventSequence()) {
        log.warn("Stale event ignored orderId={} seq={} current={}",
            event.orderId(), event.sequenceNumber(), order.getLastEventSequence());
        return;
    }

    order.applyStatusChange(event);
    order.setLastEventSequence(event.sequenceNumber());
    orderRepository.save(order);
}
```

### 10.3 Dead Letter Queue (DLQ)

```java
// Kafka DLQ with Spring Cloud Stream
spring:
  cloud:
    stream:
      kafka:
        bindings:
          orderProcessor-in-0:
            consumer:
              enable-dlq: true
              dlq-name: order-events.dlq
              dlq-partitions: 3
              auto-commit-on-error: false

// Manual DLQ handler — reprocess or alert
@KafkaListener(topics = "order-events.dlq", groupId = "dlq-monitor")
public void handleDlqMessage(
        @Payload String rawMessage,
        @Header(KafkaHeaders.DLQ_EXCEPTION_MESSAGE) String errorMessage) {

    log.error("Message in DLQ: {} error: {}", rawMessage, errorMessage);
    alertingService.sendAlert("DLQ message requires attention", rawMessage);
    // Optionally: store for manual reprocessing
    dlqRepository.save(new DlqRecord(rawMessage, errorMessage, Instant.now()));
}
```

### 10.4 Retry Strategies

| Strategy | When to Use | Implementation |
|----------|-------------|----------------|
| **Immediate retry** | Transient network blips | Spring Retry, 3 attempts |
| **Exponential backoff** | Downstream service recovery | `FixedBackOff(1s, 2s, 4s)` |
| **DLQ after max retries** | Persistent failures | Route to DLQ, alert ops |
| **Circuit breaker** | Cascading failures | Resilience4j + event buffering |

```java
@Bean
public ConcurrentKafkaListenerContainerFactory<String, OrderCreatedEvent>
        kafkaListenerContainerFactory() {
    ConcurrentKafkaListenerContainerFactory<String, OrderCreatedEvent> factory =
        new ConcurrentKafkaListenerContainerFactory<>();

    // Retry 3 times with exponential backoff, then send to DLQ
    DefaultErrorHandler errorHandler = new DefaultErrorHandler(
        new DeadLetterPublishingRecoverer(kafkaTemplate,
            (record, ex) -> new TopicPartition("order-events.dlq", 0)),
        new ExponentialBackOffWithMaxRetries(3)
    );
    errorHandler.addNotRetryableExceptions(
        InsufficientStockException.class,  // Business errors — no retry
        ValidationException.class
    );
    factory.setCommonErrorHandler(errorHandler);
    return factory;
}
```

---

## 11. Saga Pattern — Choreography & Orchestration

### 11.1 The Problem

Distributed transactions across microservices cannot use traditional ACID. The **Saga pattern** coordinates a sequence of local transactions with compensating actions.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    ORDER SAGA — HAPPY PATH & FAILURE                             │
│                                                                                  │
│  HAPPY PATH:                                                                     │
│  Order ──▶ Inventory ──▶ Payment ──▶ Shipping ──▶ Order Confirmed              │
│  CREATE     RESERVE      CHARGE      SHIP          COMPLETE                   │
│                                                                                  │
│  FAILURE at Payment:                                                             │
│  Order ──▶ Inventory ──▶ Payment ✗                                            │
│  CREATE     RESERVE      FAILED                                                  │
│              │                                                                   │
│              └── COMPENSATE: Release Inventory                                   │
│  Order ── COMPENSATE: Cancel Order                                               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 11.2 Choreography-Based Saga

No central coordinator — services react to events independently.

```java
// === ORDER SERVICE ===
@Service
@RequiredArgsConstructor
public class OrderSagaChoreography {

    private final OrderRepository orderRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Transactional
    public Order createOrder(CreateOrderRequest request) {
        Order order = orderRepository.save(Order.pending(request));
        kafkaTemplate.send("order-events", order.getId().toString(),
            OrderCreatedEvent.from(order));
        return order;
    }

    @KafkaListener(topics = "payment-events", groupId = "order-service")
    public void onPaymentCompleted(PaymentCompletedEvent event) {
        Order order = orderRepository.findById(event.orderId()).orElseThrow();
        order.confirm();
        orderRepository.save(order);
        kafkaTemplate.send("order-events", event.orderId().toString(),
            OrderConfirmedEvent.from(order));
    }

    @KafkaListener(topics = "payment-events", groupId = "order-service")
    public void onPaymentFailed(PaymentFailedEvent event) {
        Order order = orderRepository.findById(event.orderId()).orElseThrow();
        order.cancel(event.reason());
        orderRepository.save(order);
        kafkaTemplate.send("order-events", event.orderId().toString(),
            new OrderCancelledEvent(event.orderId(), event.reason()));
    }
}

// === INVENTORY SERVICE ===
@Service
@RequiredArgsConstructor
public class InventorySagaChoreography {

    private final InventoryRepository inventoryRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @KafkaListener(topics = "order-events", groupId = "inventory-service")
    public void onOrderCreated(OrderCreatedEvent event) {
        try {
            reserveStock(event.items());
            kafkaTemplate.send("inventory-events", event.orderId().toString(),
                new StockReservedEvent(event.orderId(), event.items()));
        } catch (InsufficientStockException e) {
            kafkaTemplate.send("inventory-events", event.orderId().toString(),
                new StockReserveFailedEvent(event.orderId(), e.getMessage()));
        }
    }

    @KafkaListener(topics = "order-events", groupId = "inventory-service")
    public void onOrderCancelled(OrderCancelledEvent event) {
        releaseReservedStock(event.orderId());
    }
}

// === PAYMENT SERVICE ===
@Service
@RequiredArgsConstructor
public class PaymentSagaChoreography {

    private final PaymentGateway paymentGateway;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @KafkaListener(topics = "inventory-events", groupId = "payment-service")
    public void onStockReserved(StockReservedEvent event) {
        try {
            paymentGateway.charge(event.orderId(), event.totalAmount());
            kafkaTemplate.send("payment-events", event.orderId().toString(),
                new PaymentCompletedEvent(event.orderId()));
        } catch (PaymentException e) {
            kafkaTemplate.send("payment-events", event.orderId().toString(),
                new PaymentFailedEvent(event.orderId(), e.getMessage()));
        }
    }
}
```

| Choreography Pros | Choreography Cons |
|-------------------|-------------------|
| Simple, no extra service | Hard to track saga state |
| Truly decoupled | Cyclic dependencies possible |
| Each service owns its logic | Debugging distributed flow is hard |

### 11.3 Orchestration-Based Saga

Central orchestrator coordinates all steps.

```java
// Saga state machine
public enum SagaStep {
    STARTED, INVENTORY_RESERVED, PAYMENT_COMPLETED, SHIPPING_CREATED, COMPLETED,
    COMPENSATING_INVENTORY, COMPENSATING_PAYMENT, FAILED
}

@Entity
@Data
public class OrderSaga {
    @Id
    private String sagaId;
    private Long orderId;
    @Enumerated(EnumType.STRING)
    private SagaStep currentStep;
    private String failureReason;
    private Instant startedAt;
    private Instant completedAt;
}

@Service
@RequiredArgsConstructor
@Slf4j
public class OrderSagaOrchestrator {

    private final OrderSagaRepository sagaRepository;
    private final InventoryClient inventoryClient;
    private final PaymentClient paymentClient;
    private final ShippingClient shippingClient;
    private final OrderRepository orderRepository;

    @Transactional
    public void startSaga(CreateOrderRequest request) {
        Order order = orderRepository.save(Order.pending(request));

        OrderSaga saga = new OrderSaga();
        saga.setSagaId(UUID.randomUUID().toString());
        saga.setOrderId(order.getId());
        saga.setCurrentStep(SagaStep.STARTED);
        saga.setStartedAt(Instant.now());
        sagaRepository.save(saga);

        reserveInventory(saga, order);
    }

    private void reserveInventory(OrderSaga saga, Order order) {
        try {
            inventoryClient.reserve(order.getItems());
            saga.setCurrentStep(SagaStep.INVENTORY_RESERVED);
            sagaRepository.save(saga);
            processPayment(saga, order);
        } catch (InsufficientStockException e) {
            failSaga(saga, e.getMessage());
        }
    }

    private void processPayment(OrderSaga saga, Order order) {
        try {
            paymentClient.charge(order.getId(), order.getTotalAmount());
            saga.setCurrentStep(SagaStep.PAYMENT_COMPLETED);
            sagaRepository.save(saga);
            createShipment(saga, order);
        } catch (PaymentException e) {
            compensateInventory(saga, order);
        }
    }

    private void createShipment(OrderSaga saga, Order order) {
        try {
            shippingClient.createShipment(order);
            saga.setCurrentStep(SagaStep.COMPLETED);
            saga.setCompletedAt(Instant.now());
            sagaRepository.save(saga);

            order.confirm();
            orderRepository.save(order);
        } catch (ShippingException e) {
            compensatePayment(saga, order);
        }
    }

    private void compensateInventory(OrderSaga saga, Order order) {
        saga.setCurrentStep(SagaStep.COMPENSATING_INVENTORY);
        sagaRepository.save(saga);
        inventoryClient.release(order.getItems());
        failSaga(saga, "Payment failed");
    }

    private void compensatePayment(OrderSaga saga, Order order) {
        saga.setCurrentStep(SagaStep.COMPENSATING_PAYMENT);
        sagaRepository.save(saga);
        paymentClient.refund(order.getId());
        compensateInventory(saga, order);
    }

    private void failSaga(OrderSaga saga, String reason) {
        saga.setCurrentStep(SagaStep.FAILED);
        saga.setFailureReason(reason);
        sagaRepository.save(saga);

        Order order = orderRepository.findById(saga.getOrderId()).orElseThrow();
        order.cancel(reason);
        orderRepository.save(order);
    }
}
```

| Orchestration Pros | Orchestration Cons |
|--------------------|-------------------|
| Centralized visibility | Orchestrator is a single point of failure |
| Easy to understand flow | Orchestrator can become complex |
| Clear compensation logic | Tighter coupling to orchestrator |

### 11.4 Choosing Choreography vs Orchestration

| Factor | Choreography | Orchestration |
|--------|-------------|---------------|
| **Saga complexity** | Simple (2-3 steps) | Complex (4+ steps) |
| **Team structure** | Autonomous teams | Central platform team |
| **Visibility needs** | Low | High (need saga dashboard) |
| **Failure handling** | Distributed compensation | Centralized rollback |
| **Example** | Notification on order | Multi-step payment flow |

---

## 12. Complete E-Commerce Example

### 12.1 System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    E-COMMERCE EVENT-DRIVEN ARCHITECTURE                          │
│                                                                                  │
│  ┌─────────────┐                                                                │
│  │   Client    │                                                                │
│  │  (Browser)  │                                                                │
│  └──────┬──────┘                                                                │
│         │ POST /api/orders                                                       │
│         ▼                                                                        │
│  ┌─────────────┐   OrderCreated    ┌─────────────────────────────────────┐    │
│  │   Order     │ ─────────────────▶│           Apache Kafka               │    │
│  │   Service   │                   │                                      │    │
│  │  (PostgreSQL│                   │  Topics:                             │    │
│  │   + Outbox) │                   │  • order-events                      │    │
│  └─────────────┘                   │  • inventory-events                  │    │
│                                    │  • payment-events                    │    │
│                                    │  • notification-events               │    │
│                                    └──────┬──────┬──────┬────────────────┘    │
│                                           │      │      │                       │
│                          ┌────────────────┘      │      └────────────────┐    │
│                          ▼                       ▼                       ▼    │
│                   ┌─────────────┐         ┌─────────────┐         ┌───────────┐ │
│                   │  Inventory  │         │   Payment   │         │Notification│ │
│                   │   Service   │         │   Service   │         │  Service  │ │
│                   │ (PostgreSQL)│         │ (PostgreSQL)│         │  (Email)  │ │
│                   └─────────────┘         └─────────────┘         └───────────┘ │
│                                                                                  │
│  Supporting:                                                                     │
│  • API Gateway (Spring Cloud Gateway)                                           │
│  • Service Discovery (Eureka / K8s DNS)                                         │
│  • Distributed Tracing (Zipkin)                                                   │
│  • Schema Registry (Confluent)                                                    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 12.2 Event Flow — Place Order

```
Step 1: Client → POST /api/orders
        Order Service saves order + outbox event (single transaction)
        Returns 201 Created immediately

Step 2: Outbox Poller → publishes OrderCreatedEvent to order-events topic

Step 3: Inventory Service consumes OrderCreatedEvent
        Reserves stock → publishes StockReservedEvent to inventory-events

Step 4: Payment Service consumes StockReservedEvent
        Charges card → publishes PaymentCompletedEvent to payment-events

Step 5: Order Service consumes PaymentCompletedEvent
        Updates order status to CONFIRMED
        Publishes OrderConfirmedEvent

Step 6: Notification Service consumes OrderConfirmedEvent
        Sends confirmation email to customer

Total client wait time: ~50ms (Step 1 only)
Full order processing: ~2-5 seconds (async)
```

### 12.3 Order Service — Complete Implementation

```java
// Domain model
@Entity
@Table(name = "orders")
@Data
@Builder
public class Order {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String userId;
    private String status;  // PENDING, CONFIRMED, CANCELLED

    @OneToMany(cascade = CascadeType.ALL, orphanRemoval = true)
    private List<OrderItem> items = new ArrayList<>();

    private BigDecimal totalAmount;
    private Instant createdAt;
    private String failureReason;

    public static Order pending(CreateOrderRequest request) {
        BigDecimal total = request.items().stream()
            .map(i -> i.price().multiply(BigDecimal.valueOf(i.quantity())))
            .reduce(BigDecimal.ZERO, BigDecimal::add);

        return Order.builder()
            .userId(request.userId())
            .status("PENDING")
            .items(request.items())
            .totalAmount(total)
            .createdAt(Instant.now())
            .build();
    }

    public void confirm() { this.status = "CONFIRMED"; }
    public void cancel(String reason) {
        this.status = "CANCELLED";
        this.failureReason = reason;
    }
}

// REST Controller
@RestController
@RequestMapping("/api/orders")
@RequiredArgsConstructor
public class OrderController {

    private final OrderService orderService;

    @PostMapping
    public ResponseEntity<OrderResponse> createOrder(
            @Valid @RequestBody CreateOrderRequest request) {
        Order order = orderService.createOrder(request);
        return ResponseEntity.status(HttpStatus.CREATED)
            .body(OrderResponse.from(order));
    }

    @GetMapping("/{id}")
    public OrderResponse getOrder(@PathVariable Long id) {
        return OrderResponse.from(orderService.findById(id));
    }
}

// Service with Outbox
@Service
@RequiredArgsConstructor
@Slf4j
public class OrderService {

    private final OrderRepository orderRepository;
    private final OutboxRepository outboxRepository;
    private final ObjectMapper objectMapper;

    @Transactional
    public Order createOrder(CreateOrderRequest request) {
        Order order = orderRepository.save(Order.pending(request));
        saveToOutbox("ORDER_CREATED", order.getId().toString(),
            "order-events", OrderCreatedEvent.from(order));
        log.info("Order created id={} status=PENDING", order.getId());
        return order;
    }

    @Transactional
    public void confirmOrder(Long orderId) {
        Order order = findById(orderId);
        order.confirm();
        orderRepository.save(order);
        saveToOutbox("ORDER_CONFIRMED", orderId.toString(),
            "order-events", OrderConfirmedEvent.from(order));
        log.info("Order confirmed id={}", orderId);
    }

    @Transactional
    public void cancelOrder(Long orderId, String reason) {
        Order order = findById(orderId);
        order.cancel(reason);
        orderRepository.save(order);
        saveToOutbox("ORDER_CANCELLED", orderId.toString(),
            "order-events", new OrderCancelledEvent(orderId, reason));
        log.info("Order cancelled id={} reason={}", orderId, reason);
    }

    public Order findById(Long id) {
        return orderRepository.findById(id)
            .orElseThrow(() -> new OrderNotFoundException(id));
    }

    private void saveToOutbox(String eventType, String aggregateId,
            String topic, Object event) {
        try {
            OutboxEvent outbox = new OutboxEvent();
            outbox.setId(UUID.randomUUID().toString());
            outbox.setAggregateType("Order");
            outbox.setAggregateId(aggregateId);
            outbox.setEventType(eventType);
            outbox.setTopic(topic);
            outbox.setPayload(objectMapper.writeValueAsString(event));
            outbox.setCreatedAt(Instant.now());
            outbox.setPublished(false);
            outboxRepository.save(outbox);
        } catch (JsonProcessingException e) {
            throw new EventSerializationException(e);
        }
    }
}

// Event consumers in same service
@Component
@RequiredArgsConstructor
@Slf4j
public class OrderEventHandlers {

    private final OrderService orderService;

    @KafkaListener(topics = "payment-events", groupId = "order-service")
    public void onPaymentCompleted(PaymentCompletedEvent event) {
        orderService.confirmOrder(event.orderId());
    }

    @KafkaListener(topics = "payment-events", groupId = "order-service")
    public void onPaymentFailed(PaymentFailedEvent event) {
        orderService.cancelOrder(event.orderId(), event.reason());
    }

    @KafkaListener(topics = "inventory-events", groupId = "order-service")
    public void onStockReserveFailed(StockReserveFailedEvent event) {
        orderService.cancelOrder(event.orderId(), event.reason());
    }
}
```

### 12.4 Full Docker Compose Stack

```yaml
# docker-compose.yml — complete e-commerce stack
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: password
    ports: ["5432:5432"]
    volumes: ["./init-db.sql:/docker-entrypoint-initdb.d/init.sql"]

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on: [zookeeper]
    ports: ["9092:9092"]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

  order-service:
    build: ./order-service
    ports: ["8081:8080"]
    environment:
      SPRING_DATASOURCE_URL: jdbc:postgresql://postgres:5432/orders
      SPRING_KAFKA_BOOTSTRAP_SERVERS: kafka:29092
    depends_on: [postgres, kafka]

  inventory-service:
    build: ./inventory-service
    ports: ["8082:8080"]
    environment:
      SPRING_DATASOURCE_URL: jdbc:postgresql://postgres:5432/inventory
      SPRING_KAFKA_BOOTSTRAP_SERVERS: kafka:29092
    depends_on: [postgres, kafka]

  payment-service:
    build: ./payment-service
    ports: ["8083:8080"]
    environment:
      SPRING_DATASOURCE_URL: jdbc:postgresql://postgres:5432/payments
      SPRING_KAFKA_BOOTSTRAP_SERVERS: kafka:29092
    depends_on: [postgres, kafka]

  notification-service:
    build: ./notification-service
    ports: ["8084:8080"]
    environment:
      SPRING_KAFKA_BOOTSTRAP_SERVERS: kafka:29092
    depends_on: [kafka]
```

---

## 13. Event Sourcing & CQRS

### 13.1 Event Sourcing

Instead of storing current state, store the **sequence of events** that led to that state. Rebuild state by replaying events.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    EVENT SOURCING                                                │
│                                                                                  │
│  Traditional:  orders table → { id: 1, status: CONFIRMED, total: 99.99 }       │
│                                                                                  │
│  Event Sourced: events table →                                                  │
│    [1] OrderCreated     { orderId: 1, items: [...], total: 99.99 }             │
│    [2] StockReserved    { orderId: 1, items: [...] }                            │
│    [3] PaymentCompleted { orderId: 1, amount: 99.99 }                           │
│    [4] OrderConfirmed   { orderId: 1 }                                          │
│                                                                                  │
│  Current state = replay all events for aggregate                                 │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Event store entity
@Entity
@Table(name = "event_store")
@Data
public class StoredEvent {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String aggregateId;
    private String aggregateType;
    private String eventType;
    private int version;           // optimistic locking

    @Column(columnDefinition = "TEXT")
    private String payload;

    private Instant timestamp;
}

// Event-sourced aggregate
public class OrderAggregate {
    private Long id;
    private String status;
    private List<OrderItem> items;
    private BigDecimal totalAmount;
    private final List<Object> uncommittedEvents = new ArrayList<>();

    public static OrderAggregate create(CreateOrderCommand cmd) {
        OrderAggregate order = new OrderAggregate();
        order.apply(new OrderCreatedEvent(cmd.userId(), cmd.items()));
        return order;
    }

    public static OrderAggregate fromHistory(List<StoredEvent> events) {
        OrderAggregate order = new OrderAggregate();
        events.forEach(e -> order.apply(deserialize(e)));
        return order;
    }

    private void apply(Object event) {
        switch (event) {
            case OrderCreatedEvent e -> {
                this.id = e.orderId();
                this.status = "PENDING";
                this.items = e.items();
                this.totalAmount = e.totalAmount();
            }
            case OrderConfirmedEvent e -> this.status = "CONFIRMED";
            case OrderCancelledEvent e -> this.status = "CANCELLED";
            default -> throw new IllegalArgumentException("Unknown event: " + event);
        }
    }

    public void confirm() {
        apply(new OrderConfirmedEvent(this.id));
        uncommittedEvents.add(new OrderConfirmedEvent(this.id));
    }

    public List<Object> getUncommittedEvents() {
        return List.copyOf(uncommittedEvents);
    }
}

// Event store repository
@Repository
@RequiredArgsConstructor
public class EventStore {
    private final StoredEventRepository repository;
    private final ObjectMapper objectMapper;

    public void save(String aggregateId, List<Object> events, int expectedVersion) {
        int version = expectedVersion;
        for (Object event : events) {
            StoredEvent stored = new StoredEvent();
            stored.setAggregateId(aggregateId);
            stored.setEventType(event.getClass().getSimpleName());
            stored.setVersion(++version);
            stored.setPayload(objectMapper.writeValueAsString(event));
            stored.setTimestamp(Instant.now());
            repository.save(stored);
        }
    }

    public List<StoredEvent> load(String aggregateId) {
        return repository.findByAggregateIdOrderByVersionAsc(aggregateId);
    }
}
```

### 13.2 CQRS — Command Query Responsibility Segregation

Separate write model (commands → events) from read model (optimized queries).

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CQRS ARCHITECTURE                                             │
│                                                                                  │
│  WRITE SIDE:                              READ SIDE:                             │
│  ┌─────────────┐                          ┌─────────────┐                        │
│  │ POST /orders│                          │ GET /orders │                        │
│  │  (Command)  │                          │   (Query)   │                        │
│  └──────┬──────┘                          └──────▲──────┘                        │
│         │                                        │                               │
│         ▼                                        │                               │
│  ┌─────────────┐    OrderCreatedEvent    ┌──────┴──────┐                        │
│  │   Command   │ ──────────────────────▶ │  Read Model │                        │
│  │   Handler   │                         │  (denormalized│                        │
│  │             │                         │   view in     │                        │
│  │ Event Store │                         │   PostgreSQL/  │                        │
│  │ (write DB)  │                         │   Elasticsearch)│                       │
│  └─────────────┘                         └─────────────┘                        │
│                                                                                  │
│  Write DB: normalized, event-sourced       Read DB: denormalized, fast queries  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Write side — command handler
@Service
@RequiredArgsConstructor
public class CreateOrderCommandHandler {

    private final EventStore eventStore;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Transactional
    public Long handle(CreateOrderCommand command) {
        OrderAggregate order = OrderAggregate.create(command);
        eventStore.save(
            order.getId().toString(),
            order.getUncommittedEvents(),
            0
        );
        order.getUncommittedEvents().forEach(event ->
            kafkaTemplate.send("order-events", order.getId().toString(), event)
        );
        return order.getId();
    }
}

// Read side — projection updater
@Component
@RequiredArgsConstructor
public class OrderProjectionHandler {

    private final OrderReadRepository readRepository;

    @KafkaListener(topics = "order-events", groupId = "order-projection")
    public void project(OrderCreatedEvent event) {
        OrderReadModel view = OrderReadModel.builder()
            .orderId(event.orderId())
            .userId(event.userId())
            .status("PENDING")
            .totalAmount(event.totalAmount())
            .itemCount(event.items().size())
            .createdAt(event.timestamp())
            .build();
        readRepository.save(view);
    }

    @KafkaListener(topics = "order-events", groupId = "order-projection")
    public void project(OrderConfirmedEvent event) {
        readRepository.findByOrderId(event.orderId())
            .ifPresent(view -> {
                view.setStatus("CONFIRMED");
                readRepository.save(view);
            });
    }
}

// Read side — query
@RestController
@RequestMapping("/api/orders")
@RequiredArgsConstructor
public class OrderQueryController {

    private final OrderReadRepository readRepository;

    @GetMapping("/{id}")
    public OrderReadModel getOrder(@PathVariable Long id) {
        return readRepository.findByOrderId(id)
            .orElseThrow(() -> new OrderNotFoundException(id));
    }

    @GetMapping
    public Page<OrderReadModel> listOrders(
            @RequestParam String userId,
            Pageable pageable) {
        return readRepository.findByUserId(userId, pageable);
    }
}
```

### 13.3 When to Use Event Sourcing & CQRS

| Use When | Avoid When |
|----------|-----------|
| Complete audit trail required | Simple CRUD with few business rules |
| Temporal queries ("state at time T") | Team unfamiliar with pattern |
| Complex domain with rich behavior | Read-after-write consistency needed immediately |
| Multiple read views from same data | Small application, low complexity |

---

## 14. Testing Event-Driven Systems

### 14.1 Unit Testing Event Publishers

```java
@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock private OrderRepository orderRepository;
    @Mock private OutboxRepository outboxRepository;
    @Mock private ObjectMapper objectMapper;

    @InjectMocks private OrderService orderService;

    @Test
    void createOrder_savesOrderAndOutboxEvent() throws Exception {
        CreateOrderRequest request = new CreateOrderRequest(
            "user-1", List.of(new OrderItem("SKU-1", 2, BigDecimal.TEN)));

        Order savedOrder = Order.pending(request);
        savedOrder.setId(1L);

        when(orderRepository.save(any())).thenReturn(savedOrder);
        when(objectMapper.writeValueAsString(any())).thenReturn("{}");

        Order result = orderService.createOrder(request);

        assertThat(result.getId()).isEqualTo(1L);
        assertThat(result.getStatus()).isEqualTo("PENDING");
        verify(outboxRepository).save(argThat(outbox ->
            "ORDER_CREATED".equals(outbox.getEventType()) &&
            "order-events".equals(outbox.getTopic())
        ));
    }
}
```

### 14.2 Integration Testing with Embedded Kafka

```java
@SpringBootTest
@EmbeddedKafka(
    partitions = 1,
    topics = {"order-events", "inventory-events"},
    brokerProperties = {"listeners=PLAINTEXT://localhost:9092", "port=9092"}
)
class OrderEventIntegrationTest {

    @Autowired
    private KafkaTemplate<String, Object> kafkaTemplate;

    @Autowired
    private OrderRepository orderRepository;

    @Value("${spring.embedded.kafka.brokers}")
    private String bootstrapServers;

    @Test
    void publishOrderCreatedEvent_inventoryServiceConsumes() throws Exception {
        OrderCreatedEvent event = new OrderCreatedEvent(
            UUID.randomUUID().toString(),
            1L, "user-1",
            List.of(new OrderItem("SKU-1", 2, BigDecimal.TEN)),
            BigDecimal.valueOf(20)
        );

        kafkaTemplate.send("order-events", "1", event).get(5, TimeUnit.SECONDS);

        // Wait for async processing
        await().atMost(10, TimeUnit.SECONDS).untilAsserted(() -> {
            InventoryReservation reservation =
                inventoryRepository.findByOrderId(1L);
            assertThat(reservation).isNotNull();
            assertThat(reservation.getStatus()).isEqualTo("RESERVED");
        });
    }
}
```

### 14.3 Testing with Testcontainers

```java
@SpringBootTest
@Testcontainers
class OrderServiceKafkaIT {

    @Container
    static KafkaContainer kafka = new KafkaContainer(
        DockerImageName.parse("confluentinc/cp-kafka:7.5.0"));

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.kafka.bootstrap-servers", kafka::getBootstrapServers);
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired
    private OrderService orderService;

    @Autowired
    private KafkaConsumer<String, String> testConsumer;

    @Test
    void createOrder_publishesEventToKafka() {
        CreateOrderRequest request = buildRequest();

        Order order = orderService.createOrder(request);

        ConsumerRecords<String, String> records = KafkaTestUtils.getRecords(
            testConsumer, Duration.ofSeconds(10));

        assertThat(records.count()).isGreaterThan(0);
        String value = records.iterator().next().value();
        assertThat(value).contains("\"orderId\":" + order.getId());
    }
}
```

### 14.4 Testing Event Consumers

```java
@ExtendWith(MockitoExtension.class)
class InventoryEventConsumerTest {

    @Mock private InventoryService inventoryService;
    @Mock private ProcessedEventRepository processedEventRepository;
    @Mock private Acknowledgment acknowledgment;

    @InjectMocks private OrderEventConsumer consumer;

    @Test
    void handleOrderCreated_reservesStockAndAcknowledges() {
        OrderCreatedEvent event = buildEvent();

        when(processedEventRepository.existsByEventId(event.eventId())).thenReturn(false);

        consumer.handleOrderCreated(event, "1", 0, 0L, acknowledgment);

        verify(inventoryService).reserveStock(event.items());
        verify(processedEventRepository).save(any());
        verify(acknowledgment).acknowledge();
    }

    @Test
    void handleOrderCreated_skipsDuplicateEvents() {
        OrderCreatedEvent event = buildEvent();
        when(processedEventRepository.existsByEventId(event.eventId())).thenReturn(true);

        consumer.handleOrderCreated(event, "1", 0, 0L, acknowledgment);

        verify(inventoryService, never()).reserveStock(any());
        verify(acknowledgment).acknowledge();
    }
}
```

---

## 15. Production Concerns & Observability

### 15.1 Schema Evolution

```java
// Use versioned events for backward compatibility
public record OrderCreatedEventV2(
    String eventId,
    Long orderId,
    String userId,
    List<OrderItem> items,
    BigDecimal totalAmount,
    String currency,        // NEW in v2
    String shippingMethod   // NEW in v2
) {
    public static OrderCreatedEventV2 fromV1(OrderCreatedEventV1 v1) {
        return new OrderCreatedEventV2(
            v1.eventId(), v1.orderId(), v1.userId(),
            v1.items(), v1.totalAmount(),
            "USD",              // default for migrated events
            "STANDARD"          // default for migrated events
        );
    }
}

// Consumer handles both versions
@KafkaListener(topics = "order-events")
public void handle(ConsumerRecord<String, String> record) {
    JsonNode node = objectMapper.readTree(record.value());
    int version = node.path("version").asInt(1);

    switch (version) {
        case 1 -> processV1(objectMapper.treeToValue(node, OrderCreatedEventV1.class));
        case 2 -> processV2(objectMapper.treeToValue(node, OrderCreatedEventV2.class));
        default -> log.warn("Unknown event version: {}", version);
    }
}
```

### 15.2 Observability

```yaml
# application.yml — tracing + metrics
management:
  tracing:
    sampling:
      probability: 0.1
  endpoints:
    web:
      exposure:
        include: health,metrics,prometheus

spring:
  kafka:
    producer:
      properties:
        interceptor.classes: io.micrometer.tracing.kafka.TracingProducerInterceptor
    consumer:
      properties:
        interceptor.classes: io.micrometer.tracing.kafka.TracingConsumerInterceptor
```

```java
// Structured logging with correlation IDs
@KafkaListener(topics = "order-events")
public void handle(OrderCreatedEvent event,
        @Header(value = "X-Correlation-Id", required = false) String correlationId) {

    MDC.put("correlationId", correlationId != null ? correlationId : event.eventId());
    MDC.put("orderId", event.orderId().toString());

    try {
        log.info("Processing OrderCreatedEvent");
        inventoryService.reserveStock(event.items());
        log.info("Stock reserved successfully");
    } finally {
        MDC.clear();
    }
}

// Custom metrics
@Service
@RequiredArgsConstructor
public class EventMetrics {

    private final MeterRegistry meterRegistry;

    public void recordEventPublished(String eventType) {
        meterRegistry.counter("events.published", "type", eventType).increment();
    }

    public void recordEventProcessed(String eventType, String consumer, boolean success) {
        meterRegistry.counter("events.processed",
            "type", eventType,
            "consumer", consumer,
            "success", String.valueOf(success)
        ).increment();
    }

    public void recordEventProcessingTime(String eventType, long millis) {
        meterRegistry.timer("events.processing.time", "type", eventType)
            .record(millis, TimeUnit.MILLISECONDS);
    }
}
```

### 15.3 Production Checklist

| Category | Checklist Item |
|----------|---------------|
| **Reliability** | Outbox pattern for all critical events |
| **Reliability** | Idempotent consumers with processed_events table |
| **Reliability** | DLQ configured with alerting |
| **Reliability** | Retry with exponential backoff |
| **Ordering** | Partition key = aggregate ID |
| **Schema** | Event versioning strategy documented |
| **Monitoring** | Consumer lag alerts (Kafka) |
| **Monitoring** | DLQ depth monitoring |
| **Monitoring** | Event publish/process metrics |
| **Security** | SASL/SSL for Kafka in production |
| **Security** | ACLs per service (produce/consume permissions) |
| **Operations** | Topic retention policy defined |
| **Operations** | Consumer group naming convention |
| **Operations** | Runbook for DLQ reprocessing |

---

## 16. When NOT to Use Event-Driven Architecture

| Scenario | Why Not EDA | Better Alternative |
|----------|------------|-------------------|
| **Simple CRUD** | Overhead not justified | Direct REST/gRPC |
| **Immediate response needed** | Async adds latency for reads | Synchronous call |
| **Strong consistency required** | EDA is eventually consistent | Distributed transaction or monolith |
| **Low traffic** | Broker ops overhead | In-process Spring events |
| **Request-reply pattern** | Awkward with pure events | REST, gRPC, RabbitMQ RPC |
| **Small team / startup** | Complexity slows development | Modular monolith first |

```
Decision Framework:

  Need immediate response?  ──YES──▶  Synchronous (REST/gRPC)
         │
         NO
         │
  Multiple consumers?  ──YES──▶  Event-Driven (Kafka)
         │
         NO
         │
  Decouple within monolith?  ──YES──▶  Spring Application Events
         │
         NO
         │
  Simple service call  ──────────▶  Direct REST
```

---

## 17. Interview Questions

### Fundamentals

**Q1: What is Event-Driven Architecture and how does it differ from request-response?**
> EDA is a style where services communicate through events (records of past state changes) via a message broker, rather than direct synchronous API calls. In request-response, the caller waits for a response and is tightly coupled to the callee. In EDA, the producer publishes an event and continues immediately; consumers react independently and asynchronously. EDA provides loose coupling, better scalability, and resilience at the cost of eventual consistency and increased complexity.

**Q2: What is the difference between an event and a command?**
> An event describes something that **already happened** (past tense: `OrderCreated`). It is immutable and broadcast to unknown consumers. A command instructs something **to happen** (imperative: `CreateOrder`). It is sent to a specific handler and implies an expectation of action. Confusing the two leads to coupling — consumers should react to events, not receive commands.

**Q3: Explain eventual consistency in the context of EDA.**
> In EDA, when a service publishes an event, downstream consumers process it asynchronously. There is a window where the system is in an inconsistent state — e.g., order exists but inventory not yet reserved. The system eventually reaches consistency once all consumers process their events. This trade-off enables higher availability and performance. Design compensating actions (sagas) for when consistency cannot be achieved.

### Patterns

**Q4: What is the Transactional Outbox pattern?**
> The Outbox pattern solves the dual-write problem: saving to DB and publishing to a broker atomically. Instead of publishing directly, the service saves the event to an `outbox_events` table in the same database transaction as the business data. A separate poller (or Debezium CDC) reads unpublished events and publishes them to Kafka. This guarantees at-least-once delivery without losing events.

**Q5: Compare choreography vs orchestration sagas.**
> - **Choreography**: No central coordinator. Services publish and react to events independently. Simple and decoupled, but hard to track state and debug.
> - **Orchestration**: A central saga orchestrator directs each step and handles compensation. Easy to visualize and manage, but the orchestrator is a single point of failure and adds coupling.
> Use choreography for simple flows (2-3 steps); orchestration for complex multi-step business processes.

**Q6: What is Event Sourcing? When would you use it?**
> Event Sourcing stores state as a sequence of events rather than current state. To get current state, replay all events for an aggregate. Use when you need a complete audit trail, temporal queries, or complex domain behavior. Avoid for simple CRUD where the operational complexity isn't justified.

**Q7: Explain CQRS and its relationship to Event Sourcing.**
> CQRS separates the write model (commands, normalized) from the read model (queries, denormalized). They can use different databases optimized for their purpose. Event Sourcing is often the write side; events are projected into read-optimized views. They are complementary but independent — you can use CQRS without Event Sourcing.

### Kafka & Messaging

**Q8: Compare Kafka vs RabbitMQ for microservices.**

| Aspect | Kafka | RabbitMQ |
|--------|-------|----------|
| Model | Log-based, pull | Queue-based, push |
| Ordering | Per partition | Per queue |
| Retention | Configurable (days/weeks) | Until consumed |
| Throughput | Millions/sec | Thousands/sec |
| Replay | Yes | No |
| Use case | Event streaming, analytics | Task queues, RPC, routing |

**Q9: How do you ensure exactly-once processing in Kafka?**
> True exactly-once requires: (1) idempotent producer (`enable.idempotence=true`), (2) transactional producer with `transactional.id`, (3) read-process-write in a single Kafka transaction, and (4) idempotent consumer as a safety net. In practice, most teams design for at-least-once delivery with idempotent consumers using a `processed_events` deduplication table.

**Q10: How do you handle event ordering in a distributed system?**
> - Use the aggregate ID as the Kafka message key (same key → same partition → ordered)
> - Include sequence numbers in events; consumers reject stale events
> - Accept that global ordering is expensive; partition-level ordering is sufficient for most cases
> - Never run multiple consumers on the same partition for the same logic

### Production

**Q11: How do you design idempotent event consumers?**
> Store processed `eventId` + `consumerGroup` in a database table with a unique constraint. Before processing, check if the event was already handled. Process the business logic and mark the event as processed in the same transaction. This handles at-least-once delivery safely.

**Q12: What is a Dead Letter Queue and when do messages go there?**
> A DLQ is a separate topic/queue where messages land after exhausting all retry attempts or encountering non-retryable business errors (e.g., validation failure). Operations teams monitor DLQ depth, investigate failures, fix issues, and reprocess messages. Always alert when DLQ depth exceeds a threshold.

**Q13: Design an event-driven order processing system for an e-commerce platform.**
> ```
> 1. Order Service: POST /orders → save order + outbox event → return 201
> 2. Outbox Poller → Kafka topic: order-events
> 3. Inventory Service: consume OrderCreated → reserve stock → publish StockReserved
> 4. Payment Service: consume StockReserved → charge card → publish PaymentCompleted
> 5. Order Service: consume PaymentCompleted → confirm order
> 6. Notification Service: consume OrderConfirmed → send email
>
> Failure: Payment fails → PaymentFailed event → Order Service cancels →
>          Inventory Service releases stock (compensation)
>
> Key design: Outbox pattern, idempotent consumers, partition by orderId,
>             DLQ for unprocessable events, saga for distributed transaction
> ```

### Quick Reference

| Annotation / Config | Purpose |
|-------------------|---------|
| `@KafkaListener` | Consume Kafka messages |
| `@RabbitListener` | Consume RabbitMQ messages |
| `@EventListener` | In-process Spring events |
| `@TransactionalEventListener` | Fire after DB commit |
| `StreamBridge.send()` | Imperative publish (Spring Cloud Stream) |
| `enable.idempotence=true` | Kafka idempotent producer |
| `enable-dlq: true` | Spring Cloud Stream DLQ |
| `x-dead-letter-exchange` | RabbitMQ DLQ routing |

---

## Series Navigation

- [Part 1: Spring Framework](Part-1_Spring-Framework-Complete-Guide.md)
- [Part 2: Spring Boot](Part-2_Spring-Boot-Complete-Guide.md)
- [Part 3: Spring Security](Part-3_Spring-Security-Complete-Guide.md)
- [Part 4: JPA & Hibernate](Part-4_JPA-Hibernate-Complete-Guide.md)
- [Part 5: Microservices (Part 1)](Part-5_Microservices-Spring-Cloud-Guide-Part1.md)
- [Part 6: Microservices (Part 2)](Part-6_Microservices-Spring-Cloud-Guide-Part2.md)
- [Part 7: JUnit & Mockito](Part-7_JUnit-Mockito-Testing-Complete-Guide.md)
- [Part 8: Caching](Part-8_Caching-Spring-Boot-Complete-Guide.md)
- **Part 9: Event-Driven Architecture** (this guide)

---

*Document Version: 1.0*
*Last Updated: 2026*

---

**Navigation:** [← Caching](Part-8_Caching-Spring-Boot-Complete-Guide.md)
