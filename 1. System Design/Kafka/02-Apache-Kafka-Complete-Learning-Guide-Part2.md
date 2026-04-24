# Apache Kafka Complete Learning Guide — Part 2 (Advanced)

> A structured intermediate-to-advanced guide continuing from Part 1.  
> **Focus:** Kafka Streams, Kafka Connect, Schema Registry, Security, Performance Tuning, Design Patterns, and Production Readiness.  
> **Prerequisite:** Complete understanding of [Part 1](./Apache-Kafka-Complete-Learning-Guide.md) (Core concepts, architecture, Spring Boot basics).

---

## Table of Contents

1. [Kafka Streams — Stream Processing](#1-kafka-streams--stream-processing)
2. [Kafka Connect — Data Integration](#2-kafka-connect--data-integration)
3. [Schema Registry & Schema Evolution](#3-schema-registry--schema-evolution)
4. [Log Compaction Deep Dive](#4-log-compaction-deep-dive)
5. [Multi-Broker Cluster Setup & Failure Scenarios](#5-multi-broker-cluster-setup--failure-scenarios)
6. [Monitoring & Observability](#6-monitoring--observability)
7. [Kafka Security](#7-kafka-security)
8. [Performance Tuning & Best Practices](#8-performance-tuning--best-practices)
9. [Kafka Design Patterns](#9-kafka-design-patterns)
10. [Advanced Spring Boot Kafka Patterns](#10-advanced-spring-boot-kafka-patterns)
11. [Troubleshooting Common Issues](#11-troubleshooting-common-issues)
12. [Advanced Interview Questions](#12-advanced-interview-questions)
13. [Quick Reference Cheat Sheet — Part 2](#13-quick-reference-cheat-sheet--part-2)

---

## 1. Kafka Streams — Stream Processing

### 1.1 What is Kafka Streams?

**Simple Definition:**  
Kafka Streams is a **client library** (not a separate cluster) for building real-time stream processing applications that read from and write to Kafka topics.

**Real-World Analogy:**  
Think of a **water treatment plant**. Water (data) flows in continuously from a river (Kafka topic). The plant filters, transforms, and routes it to different destinations (output topics). Kafka Streams is the treatment plant — it processes data as it flows, not in batches.

**Key Properties:**
- **No separate cluster needed** — it runs inside your Java/Spring Boot application
- **Exactly-once processing** built in
- **Fault-tolerant** — automatically recovers from failures
- **Scalable** — just run more application instances
- **Stateful processing** — can maintain state (counts, aggregations, joins)

```
KAFKA STREAMS vs TRADITIONAL BATCH PROCESSING
═══════════════════════════════════════════════

Batch Processing (Spark, Hadoop):
  ┌──────────┐     ┌──────────┐     ┌──────────┐
  │ Collect  │ ──> │ Process  │ ──> │ Output   │
  │ all data │     │ in bulk  │     │ results  │
  │ (hours)  │     │ (minutes)│     │          │
  └──────────┘     └──────────┘     └──────────┘
  Latency: Minutes to Hours

Stream Processing (Kafka Streams):
  Data ─────> Process ─────> Output
  flows       each record    immediately
  in real     as it          (ms latency)
  time        arrives
  Latency: Milliseconds
```

### 1.2 Kafka Streams vs Other Frameworks

| Feature | Kafka Streams | Apache Flink | Apache Spark Streaming |
|---------|--------------|--------------|----------------------|
| **Deployment** | Library (runs in your app) | Separate cluster | Separate cluster |
| **Infra needed** | Just Kafka | Flink cluster + Kafka | Spark cluster + Kafka |
| **Latency** | Milliseconds | Milliseconds | Seconds (micro-batch) |
| **Exactly-once** | Built-in | Built-in | Possible but complex |
| **State management** | RocksDB (local) | Managed state | Checkpointing |
| **Learning curve** | Low (Java library) | High | Medium |
| **Best for** | Kafka-centric apps | Complex event processing | Big data / ML pipelines |
| **Scaling** | Add app instances | Add TaskManagers | Add executors |

**When to Use Kafka Streams:**
- Your input AND output are Kafka topics
- You want low-latency processing (ms, not seconds)
- You don't want to manage another cluster (Flink/Spark)
- Your team already knows Java/Spring Boot

### 1.3 Core Concepts

#### KStream vs KTable

```
KStream — A stream of ALL events (like an event log)
══════════════════════════════════════════════════════

  Time ──────────────────────────────────────────────>

  Key: user-1    Key: user-2    Key: user-1    Key: user-1
  Val: login     Val: login     Val: purchase  Val: logout

  KStream sees ALL 4 records.
  Every event is independent.
  Like: "What happened?"


KTable — A changelog of LATEST state per key (like a database table)
════════════════════════════════════════════════════════════════════

  Time ──────────────────────────────────────────────>

  Key: user-1    Key: user-2    Key: user-1    Key: user-1
  Val: login     Val: login     Val: purchase  Val: logout

  KTable at the END:
    ┌──────────┬──────────┐
    │   Key    │  Value   │
    ├──────────┼──────────┤
    │ user-1   │ logout   │  ← latest value for user-1
    │ user-2   │ login    │  ← latest value for user-2
    └──────────┴──────────┘

  KTable keeps only the LATEST value per key.
  Like: "What is the current state?"
```

| Concept | KStream | KTable |
|---------|---------|--------|
| **Represents** | Event stream (append-only) | Changelog (latest per key) |
| **Analogy** | Bank transaction history | Current account balance |
| **Record meaning** | "Something happened" | "Current state of this key" |
| **Duplicates** | All kept | Latest per key overwrites |
| **Source** | Any topic | Compacted topic (usually) |

#### GlobalKTable

A `GlobalKTable` is like a `KTable`, but replicated to **every application instance**. Useful for small, slowly-changing reference data (e.g., country codes, product catalog) that every instance needs for joins.

```
GlobalKTable — Full Copy on Every Instance
═══════════════════════════════════════════

  Instance 1                Instance 2                Instance 3
  ┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
  │ GlobalKTable:   │      │ GlobalKTable:   │      │ GlobalKTable:   │
  │ "countries"     │      │ "countries"     │      │ "countries"     │
  │ ┌─────┬───────┐ │      │ ┌─────┬───────┐ │      │ ┌─────┬───────┐ │
  │ │ US  │ USA   │ │      │ │ US  │ USA   │ │      │ │ US  │ USA   │ │
  │ │ IN  │ India │ │      │ │ IN  │ India │ │      │ │ IN  │ India │ │
  │ │ UK  │ UK    │ │      │ │ UK  │ UK    │ │      │ │ UK  │ UK    │ │
  │ └─────┴───────┘ │      │ └─────┴───────┘ │      │ └─────┴───────┘ │
  └─────────────────┘      └─────────────────┘      └─────────────────┘

  Every instance has the COMPLETE table.
  No need for co-partitioning for joins.
```

### 1.4 Stream Processing Operations

#### Stateless Operations

```
STATELESS OPERATIONS (no memory of previous records)
════════════════════════════════════════════════════

  filter()     — Keep only records matching a condition
  ─────────────────────────────────────────────────
  Input:   [("u1", 10), ("u2", 5), ("u3", 15)]
  .filter((k, v) -> v > 8)
  Output:  [("u1", 10), ("u3", 15)]


  map() / mapValues()  — Transform each record
  ─────────────────────────────────────────────
  Input:   [("u1", "hello"), ("u2", "world")]
  .mapValues(v -> v.toUpperCase())
  Output:  [("u1", "HELLO"), ("u2", "WORLD")]


  flatMap()    — One record → zero or more records
  ─────────────────────────────────────────────────
  Input:   [("u1", "hello world")]
  .flatMapValues(v -> Arrays.asList(v.split(" ")))
  Output:  [("u1", "hello"), ("u1", "world")]


  branch()     — Split stream into multiple streams
  ─────────────────────────────────────────────────
  Input stream:  [orders with amount]
  .split()
    .branch((k, v) -> v.amount > 1000, Branched.as("high"))
    .branch((k, v) -> v.amount > 100,  Branched.as("medium"))
    .defaultBranch(Branched.as("low"));

  → highStream, mediumStream, lowStream


  selectKey()  — Change the key of records
  ─────────────────────────────────────────
  Input:   [("orderId", {userId: "u1", ...})]
  .selectKey((k, v) -> v.getUserId())
  Output:  [("u1", {userId: "u1", ...})]
```

#### Stateful Operations

```
STATEFUL OPERATIONS (maintain state across records)
═══════════════════════════════════════════════════

  count()      — Count records per key
  ────────────────────────────────────
  Input:   [("A", 1), ("B", 1), ("A", 1), ("A", 1)]
  .groupByKey().count()
  Output:  KTable: { "A": 3, "B": 1 }


  aggregate()  — Custom aggregation per key
  ──────────────────────────────────────────
  Input:   [("u1", 100), ("u1", 200), ("u2", 50)]
  .groupByKey().aggregate(
      () -> 0,                             // initializer
      (key, value, agg) -> agg + value     // aggregator
  )
  Output:  KTable: { "u1": 300, "u2": 50 }


  reduce()     — Combine records per key
  ──────────────────────────────────────
  Input:   [("u1", 10), ("u1", 20), ("u1", 5)]
  .groupByKey().reduce((a, b) -> Math.max(a, b))
  Output:  KTable: { "u1": 20 }


  join()       — Combine two streams/tables
  ──────────────────────────────────────────
  (Explained in detail below)
```

#### Windowed Operations

```
WINDOWED OPERATIONS — Aggregation over time windows
════════════════════════════════════════════════════

Tumbling Window (fixed, non-overlapping)
────────────────────────────────────────
  |── 5 min ──|── 5 min ──|── 5 min ──|
  | e1 e2 e3  | e4 e5     | e6 e7 e8  |
  Each event belongs to exactly ONE window.

  .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofMinutes(5)))
  Use case: "Count orders per 5-minute window"


Hopping Window (fixed, overlapping)
───────────────────────────────────
  |── 10 min ──────────|
       |── 10 min ──────────|
            |── 10 min ──────────|
  Advance: 5 min
  Each event can belong to MULTIPLE windows.

  .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofMinutes(10))
                          .advanceBy(Duration.ofMinutes(5)))
  Use case: "Rolling 10-minute average, updated every 5 min"


Sliding Window (event-triggered, for joins)
───────────────────────────────────────────
  Defined by a time difference between events.
  Window only exists when events from both sides occur within the window.

  .join(otherStream, ..., JoinWindows.ofTimeDifferenceWithNoGrace(Duration.ofMinutes(5)))
  Use case: "Match order-placed and payment-received within 5 min"


Session Window (inactivity-based)
─────────────────────────────────
  |─ events ─|  gap  |─ events ─|  gap  |─ events ─|
  Sessions end after a period of inactivity.
  Variable-length windows.

  .windowedBy(SessionWindows.ofInactivityGapWithNoGrace(Duration.ofMinutes(30)))
  Use case: "User session activity (session ends after 30 min of inactivity)"
```

### 1.5 Joins in Kafka Streams

```
JOIN TYPES
══════════

KStream–KStream Join (Windowed)
───────────────────────────────
  Orders Stream     Payments Stream
  ┌───────────┐     ┌──────────────┐
  │ orderId:1 │     │ orderId:1    │
  │ amount:100│     │ status:PAID  │
  └─────┬─────┘     └──────┬───────┘
        │                  │
        └──── JOIN ────────┘
              │
              ▼
  ┌───────────────────────────────┐
  │ orderId:1, amount:100, PAID   │
  └───────────────────────────────┘

  Both sides are infinite event streams.
  REQUIRES a time window (events must match within a time range).


KStream–KTable Join (Enrichment)
────────────────────────────────
  Orders Stream     Users KTable
  ┌───────────┐     ┌──────────────┐
  │ userId:42 │     │ userId:42    │
  │ item:Book │     │ name:Alice   │
  └─────┬─────┘     └──────┬───────┘
        │                  │
        └──── JOIN ────────┘
              │
              ▼
  ┌───────────────────────────────┐
  │ userId:42, item:Book,         │
  │ name:Alice                    │
  └───────────────────────────────┘

  Enrich stream events with lookup data.
  NO window needed (table always has latest state).


KTable–KTable Join
──────────────────
  Like a SQL table join.
  Both sides are tables (latest value per key).
  Result is also a KTable.
```

| Join Type | Left | Right | Window Required | Use Case |
|-----------|------|-------|-----------------|----------|
| KStream–KStream | Stream | Stream | Yes | Correlating events within time |
| KStream–KTable | Stream | Table | No | Enriching events with reference data |
| KStream–GlobalKTable | Stream | GlobalKTable | No | Enriching without co-partitioning |
| KTable–KTable | Table | Table | No | Joining two changing datasets |

### 1.6 Kafka Streams with Spring Boot

**Step 1: Add Dependency**

```xml
<dependency>
    <groupId>org.apache.kafka</groupId>
    <artifactId>kafka-streams</artifactId>
</dependency>
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka</artifactId>
</dependency>
```

**Step 2: Configuration**

```yaml
spring:
  kafka:
    bootstrap-servers: localhost:9092
    streams:
      application-id: order-analytics-stream
      properties:
        default.key.serde: org.apache.kafka.common.serialization.Serdes$StringSerde
        default.value.serde: org.apache.kafka.common.serialization.Serdes$StringSerde
        state.dir: /tmp/kafka-streams
```

**Step 3: Stream Topology — Word Count Example**

```java
@Configuration
@EnableKafkaStreams
public class WordCountStreamConfig {

    @Bean
    public KStream<String, String> wordCountStream(StreamsBuilder builder) {
        KStream<String, String> textLines = builder.stream("text-input");

        KTable<String, Long> wordCounts = textLines
            .flatMapValues(value -> Arrays.asList(value.toLowerCase().split("\\W+")))
            .groupBy((key, word) -> word)
            .count(Materialized.as("word-counts-store"));

        wordCounts.toStream().to("word-count-output",
            Produced.with(Serdes.String(), Serdes.Long()));

        return textLines;
    }
}
```

**Step 4: Stream Topology — Order Analytics**

```java
@Configuration
@EnableKafkaStreams
public class OrderAnalyticsStreamConfig {

    @Bean
    public KStream<String, OrderEvent> orderAnalyticsStream(StreamsBuilder builder) {

        JsonSerde<OrderEvent> orderSerde = new JsonSerde<>(OrderEvent.class);

        KStream<String, OrderEvent> orders = builder.stream(
            "order-events",
            Consumed.with(Serdes.String(), orderSerde)
        );

        KStream<String, OrderEvent>[] branches = orders
            .split(Named.as("order-"))
            .branch((key, order) -> order.getAmount() > 1000,
                    Branched.as("high-value"))
            .branch((key, order) -> order.getAmount() > 100,
                    Branched.as("medium-value"))
            .defaultBranch(Branched.as("low-value"))
            .noDefaultBranch();

        KTable<Windowed<String>, Long> orderCountsPer5Min = orders
            .groupByKey()
            .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofMinutes(5)))
            .count(Materialized.as("order-counts-5min"));

        KTable<String, Double> totalAmountPerUser = orders
            .groupBy((key, order) -> order.getUserId())
            .aggregate(
                () -> 0.0,
                (userId, order, totalAmount) -> totalAmount + order.getAmount(),
                Materialized.<String, Double, KeyValueStore<Bytes, byte[]>>as(
                    "total-amount-per-user")
                    .withValueSerde(Serdes.Double())
            );

        return orders;
    }
}
```

**Step 5: Interactive Queries (Query State from REST API)**

```java
@RestController
@RequestMapping("/api/analytics")
@RequiredArgsConstructor
public class AnalyticsController {

    private final StreamsBuilderFactoryBean factoryBean;

    @GetMapping("/word-count/{word}")
    public ResponseEntity<Long> getWordCount(@PathVariable String word) {
        KafkaStreams streams = factoryBean.getKafkaStreams();
        ReadOnlyKeyValueStore<String, Long> store = streams.store(
            StoreQueryParameters.fromNameAndType(
                "word-counts-store", QueryableStoreTypes.keyValueStore()));

        Long count = store.get(word);
        return ResponseEntity.ok(count != null ? count : 0L);
    }

    @GetMapping("/user-spending/{userId}")
    public ResponseEntity<Double> getUserSpending(@PathVariable String userId) {
        KafkaStreams streams = factoryBean.getKafkaStreams();
        ReadOnlyKeyValueStore<String, Double> store = streams.store(
            StoreQueryParameters.fromNameAndType(
                "total-amount-per-user", QueryableStoreTypes.keyValueStore()));

        Double total = store.get(userId);
        return ResponseEntity.ok(total != null ? total : 0.0);
    }
}
```

### 1.7 Kafka Streams Architecture

```
KAFKA STREAMS INTERNAL ARCHITECTURE
════════════════════════════════════

  Your Application (JVM Process)
  ┌──────────────────────────────────────────────────────────┐
  │                                                          │
  │  ┌────────────────────┐  ┌────────────────────┐          │
  │  │  Stream Thread 1   │  │  Stream Thread 2   │   ...    │
  │  │  ┌──────────────┐  │  │  ┌──────────────┐  │          │
  │  │  │  Task 0_0    │  │  │  │  Task 0_2    │  │          │
  │  │  │  (P0 of T)   │  │  │  │  (P2 of T)   │  │          │
  │  │  │  ┌─────────┐ │  │  │  │  ┌─────────┐ │  │          │
  │  │  │  │ RocksDB │ │  │  │  │  │ RocksDB │ │  │          │
  │  │  │  │ (state) │ │  │  │  │  │ (state) │ │  │          │
  │  │  │  └─────────┘ │  │  │  │  └─────────┘ │  │          │
  │  │  └──────────────┘  │  │  └──────────────┘  │          │
  │  │  ┌──────────────┐  │  │                    │          │
  │  │  │  Task 0_1    │  │  │                    │          │
  │  │  │  (P1 of T)   │  │  │                    │          │
  │  │  └──────────────┘  │  │                    │          │
  │  └────────────────────┘  └────────────────────┘          │
  │                                                          │
  │  Each Task processes ONE partition of the input topic.   │
  │  State is stored locally in RocksDB.                     │
  │  State is backed up to a changelog topic in Kafka.       │
  └──────────────────────────────────────────────────────────┘

  Scaling: Run MORE instances of the same app.
  ─────────────────────────────────────────────
  Instance 1 (Tasks: 0_0, 0_1)  ← handles partitions 0, 1
  Instance 2 (Tasks: 0_2)        ← handles partition 2
  Instance 3 (added) → rebalance → some tasks migrate here
```

---

## 2. Kafka Connect — Data Integration

### 2.1 What is Kafka Connect?

**Simple Definition:**  
Kafka Connect is a **framework** for streaming data between Kafka and external systems (databases, file systems, search engines, cloud storage) without writing any code.

**Real-World Analogy:**  
Kafka Connect is like **USB adapters**. You don't build a custom cable for every device — you use standard adapters (connectors). Need to move data from MySQL to Kafka? Plug in the "MySQL connector." Need data from Kafka to Elasticsearch? Plug in the "Elasticsearch connector."

```
KAFKA CONNECT — THE BIG PICTURE
════════════════════════════════

  External Systems                   KAFKA                    External Systems
  (Sources)                                                   (Sinks)

  ┌──────────┐                                               ┌──────────────┐
  │  MySQL   │──┐                                        ┌──>│ Elasticsearch│
  └──────────┘  │     ┌────────────────────────┐         │   └──────────────┘
  ┌──────────┐  │     │     KAFKA CLUSTER      │         │   ┌──────────────┐
  │ Postgres │──┼────>│                        │──────── ┼──>│   Amazon S3  │
  └──────────┘  │     │  topic-1  topic-2  ... │         │   └──────────────┘
  ┌──────────┐  │     │                        │         │   ┌──────────────┐
  │  MongoDB │──┘     └────────────────────────┘         └──>│   Snowflake  │
  └──────────┘                                               └──────────────┘

                SOURCE                                    SINK
              CONNECTORS                               CONNECTORS
           (External → Kafka)                       (Kafka → External)
```

### 2.2 Source vs Sink Connectors

| Aspect | Source Connector | Sink Connector |
|--------|-----------------|----------------|
| **Direction** | External system → Kafka | Kafka → External system |
| **Examples** | JDBC Source, Debezium (CDC), File Source | JDBC Sink, Elasticsearch Sink, S3 Sink |
| **Use case** | Ingest data into Kafka | Push data out of Kafka |
| **Reads from** | Database, file, API | Kafka topic |
| **Writes to** | Kafka topic | Database, file, API |

### 2.3 Popular Connectors

| Connector | Type | Use Case |
|-----------|------|----------|
| **JDBC Source** | Source | Poll database tables and send rows to Kafka |
| **Debezium MySQL/Postgres** | Source | Change Data Capture — capture every INSERT/UPDATE/DELETE |
| **File Source** | Source | Stream file contents to Kafka |
| **JDBC Sink** | Sink | Write Kafka messages to a database table |
| **Elasticsearch Sink** | Sink | Index Kafka data into Elasticsearch |
| **S3 Sink** | Sink | Store Kafka data in S3 (Parquet, JSON, Avro) |
| **HDFS Sink** | Sink | Write to Hadoop HDFS |
| **MongoDB Source/Sink** | Both | Bi-directional sync with MongoDB |

### 2.4 Kafka Connect Architecture

```
KAFKA CONNECT ARCHITECTURE
═══════════════════════════

  ┌─────────────────────────────────────────────────────────┐
  │               Connect Cluster (Workers)                 │
  │                                                         │
  │  ┌───────────────────┐   ┌───────────────────┐          │
  │  │  Worker Node 1    │   │  Worker Node 2    │   ...    │
  │  │  ┌─────────────┐  │   │  ┌─────────────┐  │          │
  │  │  │ Connector A │  │   │  │ Connector A │  │          │
  │  │  │  Task 0     │  │   │  │  Task 1     │  │          │
  │  │  └─────────────┘  │   │  └─────────────┘  │          │
  │  │  ┌─────────────┐  │   │  ┌─────────────┐  │          │
  │  │  │ Connector B │  │   │  │ Connector B │  │          │
  │  │  │  Task 0     │  │   │  │  Task 1     │  │          │
  │  │  └─────────────┘  │   │  └─────────────┘  │          │
  │  └───────────────────┘   └───────────────────┘          │
  └─────────────────────────────────────────────────────────┘

  Connector = Configuration defining what to move and how
  Task      = Unit of work (parallelism)
  Worker    = JVM process running tasks
```

**Key Terminology:**

| Term | Description |
|------|-------------|
| **Worker** | A JVM process that runs connector tasks |
| **Connector** | A logical job definition (source or sink) |
| **Task** | A unit of parallelism within a connector |
| **Converter** | Translates data between Kafka Connect format and serialization format |
| **Transform (SMT)** | Single Message Transform — lightweight per-record transformation |

### 2.5 Setting Up Kafka Connect with Docker

```yaml
version: '3.8'
services:
  kafka-connect:
    image: confluentinc/cp-kafka-connect:7.5.0
    depends_on:
      - kafka
    ports:
      - "8083:8083"
    environment:
      CONNECT_BOOTSTRAP_SERVERS: kafka:29092
      CONNECT_REST_PORT: 8083
      CONNECT_GROUP_ID: connect-cluster
      CONNECT_CONFIG_STORAGE_TOPIC: connect-configs
      CONNECT_OFFSET_STORAGE_TOPIC: connect-offsets
      CONNECT_STATUS_STORAGE_TOPIC: connect-status
      CONNECT_KEY_CONVERTER: org.apache.kafka.connect.json.JsonConverter
      CONNECT_VALUE_CONVERTER: org.apache.kafka.connect.json.JsonConverter
      CONNECT_CONFIG_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_OFFSET_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_STATUS_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_PLUGIN_PATH: /usr/share/java,/usr/share/confluent-hub-components
```

### 2.6 Example: JDBC Source Connector (MySQL → Kafka)

```json
{
  "name": "mysql-source-connector",
  "config": {
    "connector.class": "io.confluent.connect.jdbc.JdbcSourceConnector",
    "connection.url": "jdbc:mysql://mysql:3306/mydb",
    "connection.user": "root",
    "connection.password": "password",
    "table.whitelist": "orders,users",
    "mode": "incrementing",
    "incrementing.column.name": "id",
    "topic.prefix": "mysql-",
    "poll.interval.ms": "1000",
    "tasks.max": "2"
  }
}
```

**Deploy via REST API:**

```bash
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d @mysql-source-connector.json
```

**Result:** Every row in `orders` table → message in `mysql-orders` topic.

### 2.7 Example: Elasticsearch Sink Connector (Kafka → Elasticsearch)

```json
{
  "name": "elasticsearch-sink-connector",
  "config": {
    "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "connection.url": "http://elasticsearch:9200",
    "topics": "mysql-orders",
    "type.name": "_doc",
    "key.ignore": "true",
    "schema.ignore": "true",
    "tasks.max": "2"
  }
}
```

### 2.8 Debezium — Change Data Capture (CDC)

```
DEBEZIUM CDC — Capture Every Database Change
═════════════════════════════════════════════

  Traditional Polling (JDBC Source):
  ─────────────────────────────────
  Poll every N seconds → "Any new rows?"
  ⚠ Misses UPDATEs and DELETEs (unless carefully configured)
  ⚠ Adds load to the database

  Debezium CDC:
  ────────────
  Reads the database's TRANSACTION LOG (binlog for MySQL, WAL for Postgres)
  ✓ Captures INSERT, UPDATE, DELETE — everything
  ✓ Near real-time (milliseconds)
  ✓ Minimal database impact (reads the log, not the tables)

  ┌─────────────────┐    binlog     ┌──────────────┐         ┌───────┐
  │     MySQL       │ ────────────> │   Debezium   │ ──────> │ Kafka │
  │  (any change)   │   change      │  Connector   │  topic  │       │
  └─────────────────┘   stream      └──────────────┘         └───────┘
```

**Debezium MySQL Source Connector Config:**

```json
{
  "name": "mysql-debezium-connector",
  "config": {
    "connector.class": "io.debezium.connector.mysql.MySqlConnector",
    "database.hostname": "mysql",
    "database.port": "3306",
    "database.user": "debezium",
    "database.password": "dbz",
    "database.server.id": "1",
    "topic.prefix": "dbserver1",
    "database.include.list": "inventory",
    "schema.history.internal.kafka.bootstrap.servers": "kafka:29092",
    "schema.history.internal.kafka.topic": "schema-changes.inventory"
  }
}
```

### 2.9 Kafka Connect REST API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/connectors` | List all connectors |
| `POST` | `/connectors` | Create a new connector |
| `GET` | `/connectors/{name}` | Get connector info |
| `GET` | `/connectors/{name}/status` | Get connector status |
| `PUT` | `/connectors/{name}/config` | Update connector config |
| `POST` | `/connectors/{name}/restart` | Restart a connector |
| `PUT` | `/connectors/{name}/pause` | Pause a connector |
| `PUT` | `/connectors/{name}/resume` | Resume a connector |
| `DELETE` | `/connectors/{name}` | Delete a connector |
| `GET` | `/connectors/{name}/tasks` | List tasks for a connector |

---

## 3. Schema Registry & Schema Evolution

### 3.1 The Problem Without Schema Registry

```
THE SCHEMA PROBLEM
══════════════════

Day 1: Producer sends
  { "orderId": "123", "amount": 100.0 }

Day 30: Developer changes the schema
  { "orderId": "123", "amount": 100.0, "currency": "USD" }

Day 31: Old consumer CRASHES 💥
  → doesn't understand "currency" field
  → or worse: silently ignores it and loses data

Without Schema Registry:
  • No validation of what's being sent
  • No contract between producer and consumer
  • Breaking changes go undetected until RUNTIME
  • Debugging is a nightmare
```

### 3.2 What is Schema Registry?

**Simple Definition:**  
Schema Registry is a **centralized service** that stores and validates message schemas (the structure/format of data). It ensures producers and consumers agree on the data format.

**Real-World Analogy:**  
Schema Registry is like a **contract office**. Before a producer sends data, it registers the contract (schema). Before a consumer reads data, it checks the contract. If someone tries to break the contract (incompatible change), the registry rejects it.

```
SCHEMA REGISTRY IN THE ARCHITECTURE
════════════════════════════════════

  ┌──────────┐          ┌──────────────────┐          ┌──────────┐
  │ Producer │          │  Schema Registry │          │ Consumer │
  │          │   1.     │                  │   4.     │          │
  │ Create   │ ──────── │  Stores schemas: │ ──────── │ Fetch    │
  │ schema   │ Register │  • order-v1      │  Fetch   │ schema   │
  │          │ schema   │  • order-v2      │  schema  │          │
  │ 2. Send  │          │  • user-v1       │          │ 5. Use   │
  │ data     │          │                  │          │ schema   │
  │ with     │          │  Validates       │          │ to       │
  │ schema   │          │  compatibility   │          │ decode   │
  │ ID       │          │                  │          │ data     │
  └────┬─────┘          └──────────────────┘          └────┬─────┘
       │                                                   │
       │  3. Message = [schema-id (4 bytes)] + [data]      │
       └───────────────── KAFKA ───────────────────────────┘
```

### 3.3 Serialization Formats

| Format | Schema Support | Size | Speed | Human Readable | Best For |
|--------|---------------|------|-------|----------------|----------|
| **JSON** | Optional (JSON Schema) | Large | Slow | Yes | Debugging, small scale |
| **Avro** | Required | Small | Fast | No (binary) | Production (most popular) |
| **Protobuf** | Required | Smallest | Fastest | No (binary) | High-performance, gRPC |
| **JSON Schema** | Required | Large | Slow | Yes | When readability matters |

### 3.4 Apache Avro Deep Dive

**Avro Schema Example:**

```json
{
  "type": "record",
  "name": "OrderEvent",
  "namespace": "com.example.events",
  "fields": [
    { "name": "orderId", "type": "string" },
    { "name": "userId", "type": "string" },
    { "name": "amount", "type": "double" },
    { "name": "status", "type": {
        "type": "enum",
        "name": "OrderStatus",
        "symbols": ["CREATED", "CONFIRMED", "SHIPPED", "DELIVERED", "CANCELLED"]
      }
    },
    { "name": "timestamp", "type": { "type": "long", "logicalType": "timestamp-millis" } },
    { "name": "notes", "type": ["null", "string"], "default": null }
  ]
}
```

**Key Avro Features:**
- **Schema is stored with the data** (actually, only the schema ID — 4 bytes)
- **Compact binary format** (no field names in the data, unlike JSON)
- **Schema evolution** — reader and writer can use different schema versions
- **Null handling** — explicit via union types (`["null", "string"]`)

### 3.5 Schema Compatibility Types

```
COMPATIBILITY TYPES
═══════════════════

BACKWARD (Default) — New schema can READ data written with old schema
──────────────────────────────────────────────────────────────────────
  v1: { orderId, amount }              Writer (old)
  v2: { orderId, amount, currency }    Reader (new)

  Reader v2 reads v1 data → "currency" is missing → uses DEFAULT value
  ✓ Safe to deploy NEW consumers FIRST

  Allowed changes:
    ✓ Add field WITH default value
    ✓ Remove field that had default value
    ✗ Add field WITHOUT default value
    ✗ Change field type


FORWARD — Old schema can READ data written with new schema
──────────────────────────────────────────────────────────
  v1: { orderId, amount }              Reader (old)
  v2: { orderId, amount, currency }    Writer (new)

  Reader v1 reads v2 data → ignores "currency" field
  ✓ Safe to deploy NEW producers FIRST

  Allowed changes:
    ✓ Add field (old reader ignores it)
    ✓ Remove field that had default value
    ✗ Remove required field


FULL — Both backward AND forward compatible
───────────────────────────────────────────
  Both old and new schemas can read each other's data.
  Most restrictive but safest.

  Allowed changes:
    ✓ Add field WITH default value
    ✓ Remove field WITH default value


NONE — No compatibility checking
────────────────────────────────
  ⚠ Any change allowed. No safety net.
  Not recommended for production.
```

| Compatibility | Deploy Order | Allowed Changes | Safety |
|--------------|-------------|-----------------|--------|
| **BACKWARD** | Consumers first | Add optional fields, remove fields with defaults | Medium |
| **FORWARD** | Producers first | Add fields, remove optional fields | Medium |
| **FULL** | Any order | Add/remove only optional fields | High |
| **NONE** | N/A | Anything | None |

### 3.6 Spring Boot with Schema Registry (Avro)

**Step 1: Add Dependencies**

```xml
<dependency>
    <groupId>io.confluent</groupId>
    <artifactId>kafka-avro-serializer</artifactId>
    <version>7.5.0</version>
</dependency>
<dependency>
    <groupId>io.confluent</groupId>
    <artifactId>kafka-schema-registry-client</artifactId>
    <version>7.5.0</version>
</dependency>
<dependency>
    <groupId>org.apache.avro</groupId>
    <artifactId>avro</artifactId>
    <version>1.11.3</version>
</dependency>
```

**Step 2: Add Confluent Maven Repository**

```xml
<repositories>
    <repository>
        <id>confluent</id>
        <url>https://packages.confluent.io/maven/</url>
    </repository>
</repositories>
```

**Step 3: Configuration**

```yaml
spring:
  kafka:
    bootstrap-servers: localhost:9092
    properties:
      schema.registry.url: http://localhost:8081
    producer:
      key-serializer: org.apache.kafka.common.serialization.StringSerializer
      value-serializer: io.confluent.kafka.serializers.KafkaAvroSerializer
    consumer:
      key-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      value-deserializer: io.confluent.kafka.serializers.KafkaAvroDeserializer
      properties:
        specific.avro.reader: true
```

### 3.7 Schema Registry Docker Setup

```yaml
services:
  schema-registry:
    image: confluentinc/cp-schema-registry:7.5.0
    depends_on:
      - kafka
    ports:
      - "8081:8081"
    environment:
      SCHEMA_REGISTRY_HOST_NAME: schema-registry
      SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS: kafka:29092
      SCHEMA_REGISTRY_LISTENERS: http://0.0.0.0:8081
```

### 3.8 Schema Registry REST API

```bash
# List all subjects
curl http://localhost:8081/subjects

# Get latest schema for a subject
curl http://localhost:8081/subjects/order-events-value/versions/latest

# Register a new schema
curl -X POST http://localhost:8081/subjects/order-events-value/versions \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"schema": "{\"type\":\"record\",\"name\":\"Order\",\"fields\":[{\"name\":\"orderId\",\"type\":\"string\"}]}"}'

# Check compatibility
curl -X POST http://localhost:8081/compatibility/subjects/order-events-value/versions/latest \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"schema": "<new-schema-json>"}'

# Set compatibility level
curl -X PUT http://localhost:8081/config/order-events-value \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"compatibility": "FULL"}'
```

---

## 4. Log Compaction Deep Dive

### 4.1 What is Log Compaction?

**Simple Definition:**  
Log compaction ensures that Kafka retains **at least the last known value for each message key** within a topic partition. Instead of deleting old messages by time, it keeps the latest value per key.

**Real-World Analogy:**  
Think of a **phone book**. When someone changes their phone number, you don't add a second entry — you update the existing one. The phone book (compacted topic) always shows the latest number for each person.

```
LOG COMPACTION — BEFORE vs AFTER
═════════════════════════════════

BEFORE COMPACTION:
  Offset:  0     1     2     3     4     5     6     7
  Key:    [K1]  [K2]  [K1]  [K3]  [K1]  [K2]  [K3]  [K1]
  Value:  [v1]  [v1]  [v2]  [v1]  [v3]  [v2]  [v2]  [v4]

  Multiple entries for K1, K2, K3

AFTER COMPACTION:
  Offset:  5     6     7
  Key:    [K2]  [K3]  [K1]
  Value:  [v2]  [v2]  [v4]

  Only the LATEST value for each key survives.
  Offsets are preserved (not reassigned).
  Old offsets (0-4) are gone.
```

### 4.2 Delete vs Compact Retention Policies

| Aspect | `cleanup.policy=delete` | `cleanup.policy=compact` | `cleanup.policy=delete,compact` |
|--------|------------------------|-------------------------|-------------------------------|
| **Retains** | Messages within time/size limit | Latest per key (forever) | Both: compact + delete old |
| **Deletion trigger** | `retention.ms` or `retention.bytes` | Log cleaner thread | Both |
| **Use case** | Event logs, metrics | State stores, snapshots | Bounded state stores |
| **Key required** | No | Yes (null key = delete marker) | Yes |

### 4.3 Tombstone Records (Deletes in Compacted Topics)

```
TOMBSTONE — Deleting a key from a compacted topic
══════════════════════════════════════════════════

  To "delete" a key, produce a message with:
    Key = "user-42"
    Value = null       ← This is a TOMBSTONE

  Before compaction:
    [("user-42", "Alice"), ("user-42", null)]

  After compaction:
    [("user-42", null)]    ← tombstone retained briefly

  After tombstone retention period (delete.retention.ms):
    (record completely removed)

  Use case: GDPR deletion — set value to null for a user's key
```

### 4.4 When to Use Log Compaction

| Use Case | Why Compaction? |
|----------|----------------|
| **Kafka Streams state stores** | Changelog topics need latest state per key |
| **Database snapshots** | Keep current state of each DB row |
| **Configuration distribution** | Always have latest config per service |
| **User profile cache** | Latest profile per user ID |
| **`__consumer_offsets`** | Internal topic: latest offset per consumer group |

### 4.5 Configuration

```bash
# Create a compacted topic
kafka-topics --create \
  --topic user-profiles \
  --partitions 3 \
  --replication-factor 2 \
  --config cleanup.policy=compact \
  --config min.cleanable.dirty.ratio=0.5 \
  --config delete.retention.ms=86400000 \
  --bootstrap-server localhost:9092
```

| Config | Default | Description |
|--------|---------|-------------|
| `cleanup.policy` | `delete` | Set to `compact` for log compaction |
| `min.cleanable.dirty.ratio` | `0.5` | Ratio of dirty (uncompacted) log to trigger compaction |
| `delete.retention.ms` | `86400000` (24h) | How long tombstones are retained after compaction |
| `min.compaction.lag.ms` | `0` | Minimum time before a message can be compacted |
| `max.compaction.lag.ms` | `Long.MAX` | Maximum time before compaction is forced |
| `segment.ms` | `604800000` (7d) | Time before active segment is rolled (compaction only works on closed segments) |

---

## 5. Multi-Broker Cluster Setup & Failure Scenarios

### 5.1 Production Multi-Broker Docker Compose

```yaml
version: '3.8'
services:
  kafka-1:
    image: apache/kafka:3.7.0
    ports:
      - "9092:9092"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:9092
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@kafka-1:9093,2@kafka-2:9093,3@kafka-3:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
      KAFKA_MIN_INSYNC_REPLICAS: 2
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk

  kafka-2:
    image: apache/kafka:3.7.0
    ports:
      - "9093:9092"
    environment:
      KAFKA_NODE_ID: 2
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:9092
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@kafka-1:9093,2@kafka-2:9093,3@kafka-3:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
      KAFKA_MIN_INSYNC_REPLICAS: 2
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk

  kafka-3:
    image: apache/kafka:3.7.0
    ports:
      - "9094:9092"
    environment:
      KAFKA_NODE_ID: 3
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-3:9092
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@kafka-1:9093,2@kafka-2:9093,3@kafka-3:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
      KAFKA_MIN_INSYNC_REPLICAS: 2
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk
```

### 5.2 Failure Scenarios and Recovery

```
FAILURE SCENARIO 1: Single Broker Failure
══════════════════════════════════════════

  Cluster: 3 brokers, replication-factor=3, min.insync.replicas=2

  BEFORE:
  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
  │  Broker 1   │  │  Broker 2   │  │  Broker 3   │
  │  P0 (L)     │  │  P0 (F)     │  │  P0 (F)     │
  │  P1 (F)     │  │  P1 (L)     │  │  P1 (F)     │
  │  P2 (F)     │  │  P2 (F)     │  │  P2 (L)     │
  └─────────────┘  └─────────────┘  └─────────────┘
  ISR(P0) = {1,2,3}  ISR(P1) = {1,2,3}  ISR(P2) = {1,2,3}

  Broker 2 CRASHES!

  AFTER:
  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
  │  Broker 1   │  │  Broker 2   │  │  Broker 3   │
  │  P0 (L) ✓   │  │    DOWN     │  │  P0 (F) ✓   │
  │  P1 (NEW L) │  │             │  │  P1 (F) ✓   │
  │  P2 (F) ✓   │  │             │  │  P2 (L) ✓   │
  └─────────────┘  └─────────────┘  └─────────────┘
  ISR(P0) = {1,3}   ISR(P1) = {1,3}   ISR(P2) = {1,3}

  ✓ P1 leadership transfers to Broker 1
  ✓ Cluster is still fully operational
  ✓ min.insync.replicas=2 still satisfied (2 brokers alive)
  ✓ ZERO data loss, ZERO downtime


FAILURE SCENARIO 2: Two Brokers Down (Under-replicated)
═══════════════════════════════════════════════════════

  Broker 2 and Broker 3 both DOWN:

  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
  │  Broker 1   │  │  Broker 2   │  │  Broker 3   │
  │  P0 (L)     │  │    DOWN     │  │    DOWN      │
  │  P1 (L)     │  │             │  │              │
  │  P2 (L)     │  │             │  │              │
  └─────────────┘  └─────────────┘  └─────────────┘
  ISR(P0) = {1}    ← only 1 in-sync replica!

  With min.insync.replicas=2:
  ⚠ Producers with acks=all get NOT_ENOUGH_REPLICAS error
  ⚠ Writes are REJECTED (data safety preserved)
  ✓ Reads still work (consumers can read existing data)
  → Wait for brokers to recover, or reduce min.insync.replicas


FAILURE SCENARIO 3: Unclean Leader Election
════════════════════════════════════════════

  Problem: ALL ISR replicas are down. Only an out-of-sync replica is alive.

  Option A: unclean.leader.election.enable=false (default)
  → Partition becomes UNAVAILABLE until an ISR replica recovers
  → No data loss, but downtime

  Option B: unclean.leader.election.enable=true
  → Out-of-sync replica becomes leader
  → AVAILABILITY maintained but SOME DATA MAY BE LOST
  → Messages that weren't replicated to this broker are gone

  Decision:
    Data is more important → Keep false (default)
    Availability is more important → Set true
```

### 5.3 KRaft vs ZooKeeper

```
KRAFT MODE (Kafka Raft) — The Future
═════════════════════════════════════

  ZooKeeper Mode (Legacy):
  ┌──────────┐     ┌──────────┐     ┌──────────┐
  │   ZK 1   │ ──  │   ZK 2   │ ──  │   ZK 3   │  ← Separate ZK ensemble
  └──────────┘     └──────────┘     └──────────┘
       │                │                │
  ┌──────────┐     ┌──────────┐     ┌──────────┐
  │ Broker 1 │     │ Broker 2 │     │ Broker 3 │  ← Kafka brokers
  └──────────┘     └──────────┘     └──────────┘

  Problems:
  • Two separate systems to operate and monitor
  • ZooKeeper is a bottleneck for metadata operations
  • Limits cluster to ~200K partitions
  • Complex deployment


  KRaft Mode (Kafka 3.3+, production-ready in 3.5+):
  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
  │ Broker 1         │  │ Broker 2         │  │ Broker 3         │
  │ + Controller     │  │ + Controller     │  │ + Controller     │
  │ (Raft consensus) │  │ (Raft consensus) │  │ (Raft consensus) │
  └──────────────────┘  └──────────────────┘  └──────────────────┘

  Benefits:
  ✓ Single system to operate
  ✓ Faster metadata operations
  ✓ Supports millions of partitions
  ✓ Faster controller failover (seconds vs minutes)
  ✓ Simpler deployment
```

| Aspect | ZooKeeper | KRaft |
|--------|-----------|-------|
| **Architecture** | Separate ZK cluster | Built into Kafka |
| **Max partitions** | ~200K | Millions |
| **Controller failover** | Minutes | Seconds |
| **Ops complexity** | High (two systems) | Low (one system) |
| **Status (2024+)** | Deprecated | Recommended |
| **When to use** | Legacy migrations only | All new deployments |

---

## 6. Monitoring & Observability

### 6.1 Key Metrics to Monitor

```
CRITICAL KAFKA METRICS (Organized by Priority)
═══════════════════════════════════════════════

🔴 CRITICAL (Page immediately):
  • UnderReplicatedPartitions > 0
    → Data at risk! Replicas falling behind.
  • OfflinePartitionsCount > 0
    → Partitions with NO leader. Data unavailable.
  • ActiveControllerCount != 1
    → No controller = no leader elections = cluster paralyzed.

🟡 WARNING (Investigate soon):
  • Consumer Lag (records-lag-max)
    → Consumer falling behind producer. May cause data staleness.
  • RequestHandlerAvgIdlePercent < 0.3
    → Brokers overloaded. Need more brokers or better hardware.
  • NetworkProcessorAvgIdlePercent < 0.3
    → Network threads saturated.

🟢 INFORMATIONAL (Dashboard):
  • MessagesInPerSec → throughput
  • BytesInPerSec / BytesOutPerSec → bandwidth
  • PartitionCount → total partitions
  • ISRShrinkRate / ISRExpandRate → cluster stability
```

### 6.2 Monitoring Architecture

```
KAFKA MONITORING STACK
══════════════════════

  ┌──────────────┐   JMX    ┌──────────────┐  scrape  ┌────────────┐
  │ Kafka Broker │ ────────>│ JMX Exporter │ ────────>│ Prometheus │
  │              │          │ (Agent JAR)  │          │            │
  └──────────────┘          └──────────────┘          └─────┬──────┘
                                                            │
  ┌──────────────┐   JMX    ┌──────────────┐  scrape       │
  │ Kafka Broker │ ────────>│ JMX Exporter │ ─────────────┤
  └──────────────┘          └──────────────┘               │
                                                            │
                                                            ▼
                                                     ┌────────────┐
                                                     │  Grafana   │
                                                     │ Dashboards │
                                                     └────────────┘

  Alternative: Kafka UI / AKHQ / Confluent Control Center
```

### 6.3 Consumer Lag Monitoring

**Consumer lag** = the difference between the latest message produced and the latest message consumed.

```
CONSUMER LAG EXPLAINED
══════════════════════

  Partition 0:
  Produced:  [0] [1] [2] [3] [4] [5] [6] [7] [8] [9]
                                          ▲              ▲
                                          │              │
                                   Consumer offset    Log end
                                      (offset 6)     (offset 9)

  LAG = Log End Offset - Consumer Offset = 9 - 6 = 3 messages behind

  Healthy:  lag = 0 (or close to 0)
  Warning:  lag growing steadily → consumer too slow
  Critical: lag growing unbounded → consumer stuck or too many messages
```

**Monitor with CLI:**

```bash
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --describe --group my-consumer-group

# Output:
# TOPIC     PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG
# orders    0          1000            1005            5
# orders    1          2000            2000            0
# orders    2          1500            1520            20
```

### 6.4 Kafka UI Tools

| Tool | Type | Features | Cost |
|------|------|----------|------|
| **Kafka UI** (provectus) | Open Source | Topics, consumers, messages, schema registry | Free |
| **AKHQ** | Open Source | Topics, consumer groups, schema registry, connect | Free |
| **Confluent Control Center** | Commercial | Full management, monitoring, alerting | Paid |
| **Kafdrop** | Open Source | Lightweight, topic browser | Free |
| **Redpanda Console** | Open Source | Modern UI, compatible with Kafka | Free |

**Kafka UI Docker Setup:**

```yaml
services:
  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    ports:
      - "8080:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: local
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
      KAFKA_CLUSTERS_0_SCHEMAREGISTRY: http://schema-registry:8081
      KAFKA_CLUSTERS_0_KAFKACONNECT_0_NAME: connect
      KAFKA_CLUSTERS_0_KAFKACONNECT_0_ADDRESS: http://kafka-connect:8083
```

### 6.5 JMX Metrics with Spring Boot

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health,metrics,prometheus
  metrics:
    tags:
      application: order-service
    export:
      prometheus:
        enabled: true

spring:
  kafka:
    producer:
      properties:
        metric.reporters: org.apache.kafka.common.metrics.JmxReporter
    consumer:
      properties:
        metric.reporters: org.apache.kafka.common.metrics.JmxReporter
```

**Key Spring Kafka Micrometer Metrics:**

| Metric | Description |
|--------|-------------|
| `kafka.producer.record.send.total` | Total records sent |
| `kafka.producer.record.error.total` | Total send errors |
| `kafka.producer.request.latency.avg` | Average request latency |
| `kafka.consumer.records.consumed.total` | Total records consumed |
| `kafka.consumer.records.lag.max` | Maximum consumer lag |
| `kafka.consumer.fetch.latency.avg` | Average fetch latency |

---

## 7. Kafka Security

### 7.1 Security Overview

```
KAFKA SECURITY LAYERS
═════════════════════

  Client (Producer/Consumer)
     │
     │  Layer 1: ENCRYPTION (SSL/TLS)
     │  ─────────────────────────────
     │  Data encrypted in transit.
     │  Prevents eavesdropping.
     │
     │  Layer 2: AUTHENTICATION
     │  ───────────────────────
     │  "Who are you?"
     │  Options: SSL certificates, SASL/PLAIN, SASL/SCRAM, OAuth
     │
     │  Layer 3: AUTHORIZATION (ACLs)
     │  ─────────────────────────────
     │  "What are you allowed to do?"
     │  Read topic X? Write topic Y? Create topics?
     │
     ▼
  Kafka Broker
```

### 7.2 Authentication Mechanisms

| Mechanism | How It Works | Use Case |
|-----------|-------------|----------|
| **SSL/TLS (mTLS)** | Client and broker exchange certificates | Strong identity, PKI infra available |
| **SASL/PLAIN** | Username + password (plain text over TLS) | Simple setup, dev/test |
| **SASL/SCRAM** | Salted challenge-response (password never sent) | Production, no PKI needed |
| **SASL/GSSAPI** | Kerberos authentication | Enterprise with existing Kerberos |
| **SASL/OAUTHBEARER** | OAuth 2.0 tokens | Modern cloud-native apps |

### 7.3 SSL/TLS Configuration

**Broker Configuration:**

```properties
# Inter-broker communication
listeners=PLAINTEXT://:9092,SSL://:9093
advertised.listeners=PLAINTEXT://broker1:9092,SSL://broker1:9093
security.inter.broker.protocol=SSL

ssl.keystore.location=/var/ssl/private/kafka.broker.keystore.jks
ssl.keystore.password=keystore-password
ssl.key.password=key-password
ssl.truststore.location=/var/ssl/private/kafka.broker.truststore.jks
ssl.truststore.password=truststore-password
ssl.client.auth=required
```

**Spring Boot Client Configuration:**

```yaml
spring:
  kafka:
    bootstrap-servers: broker1:9093
    security:
      protocol: SSL
    ssl:
      trust-store-location: classpath:kafka.client.truststore.jks
      trust-store-password: truststore-password
      key-store-location: classpath:kafka.client.keystore.jks
      key-store-password: keystore-password
      key-password: key-password
```

### 7.4 SASL/SCRAM Configuration

**Broker Configuration:**

```properties
listeners=SASL_SSL://:9093
security.inter.broker.protocol=SASL_SSL
sasl.mechanism.inter.broker.protocol=SCRAM-SHA-256
sasl.enabled.mechanisms=SCRAM-SHA-256
```

**Spring Boot Client Configuration:**

```yaml
spring:
  kafka:
    bootstrap-servers: broker1:9093
    properties:
      security.protocol: SASL_SSL
      sasl.mechanism: SCRAM-SHA-256
      sasl.jaas.config: >
        org.apache.kafka.common.security.scram.ScramLoginModule required
        username="my-app"
        password="secret";
    ssl:
      trust-store-location: classpath:kafka.client.truststore.jks
      trust-store-password: truststore-password
```

### 7.5 Authorization with ACLs

```bash
# Grant producer permission to user "order-service"
kafka-acls --bootstrap-server localhost:9092 \
  --add --allow-principal User:order-service \
  --operation Write --topic order-events

# Grant consumer permission
kafka-acls --bootstrap-server localhost:9092 \
  --add --allow-principal User:analytics-service \
  --operation Read --topic order-events \
  --group analytics-group

# List ACLs
kafka-acls --bootstrap-server localhost:9092 --list

# Remove ACLs
kafka-acls --bootstrap-server localhost:9092 \
  --remove --allow-principal User:order-service \
  --operation Write --topic order-events
```

**Common ACL Operations:**

| Operation | Description |
|-----------|-------------|
| `Read` | Consume from a topic |
| `Write` | Produce to a topic |
| `Create` | Create topics |
| `Delete` | Delete topics |
| `Alter` | Change topic configurations |
| `Describe` | View topic metadata |
| `ClusterAction` | Cluster-level operations |
| `All` | All operations |

---

## 8. Performance Tuning & Best Practices

### 8.1 Producer Performance Tuning

```
PRODUCER PERFORMANCE LEVERS
════════════════════════════

  ┌──────────────────────────────────────────────────────────┐
  │                     Producer                              │
  │                                                          │
  │  Record ──> Serializer ──> Partitioner ──> Buffer        │
  │                                             │            │
  │                                    batch.size (32KB)     │
  │                                    linger.ms (5-20ms)    │
  │                                             │            │
  │                                    buffer.memory (32MB)  │
  │                                             │            │
  │                                    compression.type      │
  │                                    (snappy/lz4/zstd)     │
  │                                             │            │
  │                               Network Send ──> Broker    │
  └──────────────────────────────────────────────────────────┘

  THROUGHPUT KNOBS:
  ────────────────
  ↑ batch.size          → Larger batches = fewer network calls
  ↑ linger.ms           → Wait longer to fill batches
  ✓ compression.type    → Compress batches (snappy is good balance)
  ↑ buffer.memory       → More buffering capacity
  ↑ max.in.flight=5     → More parallel requests (with idempotence)

  LATENCY KNOBS:
  ──────────────
  ↓ linger.ms = 0       → Send immediately (no batching wait)
  ↓ batch.size          → Smaller batches, send sooner
  acks=1 instead of all → Don't wait for replicas
```

| Tuning Goal | Config | Recommended Value | Trade-off |
|-------------|--------|-------------------|-----------|
| **Higher throughput** | `batch.size` | `65536` (64KB) | More memory, slight latency |
| **Higher throughput** | `linger.ms` | `20` | +20ms latency |
| **Higher throughput** | `compression.type` | `lz4` or `snappy` | CPU usage |
| **Lower latency** | `linger.ms` | `0` | Lower throughput |
| **Lower latency** | `acks` | `1` | Risk of data loss |
| **Durability** | `acks` | `all` | Higher latency |
| **No duplicates** | `enable.idempotence` | `true` | Slight overhead |

### 8.2 Consumer Performance Tuning

```
CONSUMER PERFORMANCE LEVERS
════════════════════════════

  Broker ──> Network Fetch ──> Deserializer ──> Processing ──> Commit
                  │                                              │
          fetch.min.bytes                              enable.auto.commit
          fetch.max.wait.ms                            auto.commit.interval.ms
          max.poll.records                             (or manual commit)
          max.partition.fetch.bytes
```

| Tuning Goal | Config | Value | Effect |
|-------------|--------|-------|--------|
| **Higher throughput** | `fetch.min.bytes` | `1048576` (1MB) | Fetch larger chunks |
| **Higher throughput** | `max.poll.records` | `1000` | Process more per poll |
| **Higher throughput** | Increase partitions | N/A | More consumers in parallel |
| **Lower latency** | `fetch.min.bytes` | `1` | Fetch immediately |
| **Lower latency** | `fetch.max.wait.ms` | `100` | Don't wait to accumulate |
| **Avoid rebalance** | `max.poll.interval.ms` | Match processing time | Prevent consumer timeout |
| **Avoid rebalance** | `session.timeout.ms` | `30000` | Reasonable heartbeat window |

### 8.3 Broker Performance Tuning

| Config | Default | Recommendation | Description |
|--------|---------|----------------|-------------|
| `num.network.threads` | `3` | Number of CPU cores | Threads for network I/O |
| `num.io.threads` | `8` | 2x number of disks | Threads for disk I/O |
| `socket.send.buffer.bytes` | `102400` | `1048576` | Socket send buffer |
| `socket.receive.buffer.bytes` | `102400` | `1048576` | Socket receive buffer |
| `log.flush.interval.messages` | `Long.MAX` | Keep default | Let OS handle flushing |
| `log.dirs` | `/tmp/kafka-logs` | Multiple disks | Spread I/O across disks |
| `num.partitions` | `1` | `3-6` default | Default for auto-created topics |
| `default.replication.factor` | `1` | `3` | Default replication |

### 8.4 Compression Comparison

| Algorithm | Compression Ratio | CPU Cost | Speed | Best For |
|-----------|------------------|----------|-------|----------|
| **None** | 1:1 | None | Fastest | CPU-constrained, small messages |
| **Snappy** | ~1.5-2x | Low | Fast | General purpose (recommended) |
| **LZ4** | ~2-3x | Low | Very fast | High throughput, low latency |
| **ZSTD** | ~3-4x | Medium | Medium | Best compression, bandwidth-limited |
| **GZIP** | ~3-4x | High | Slow | Archival, not recommended for real-time |

### 8.5 Partition Count Guidelines

```
HOW MANY PARTITIONS?
════════════════════

  Too Few Partitions:
  • Limited consumer parallelism
  • Potential throughput bottleneck
  • Uneven load distribution

  Too Many Partitions:
  • More memory used on brokers (each partition uses ~1MB+)
  • More file handles
  • Longer leader election during failures
  • Slower rebalancing
  • Higher end-to-end latency

  FORMULA:
  ────────
  Partitions = max(T/P, T/C)

  Where:
    T = Target throughput (messages/sec)
    P = Max throughput per producer partition (~100K msg/s)
    C = Max throughput per consumer partition (~50K msg/s)

  GUIDELINES:
  ──────────
  Small topic (< 10K msg/s):  3 partitions
  Medium topic (10K-100K):    6-12 partitions
  Large topic (100K-1M):      12-30 partitions
  Very large topic (> 1M):    30-100 partitions

  Always:
  • Number of partitions >= number of consumers in the group
  • You can INCREASE partitions, but NEVER decrease
  • Start conservative, increase when needed
```

---

## 9. Kafka Design Patterns

### 9.1 Event Sourcing

```
EVENT SOURCING PATTERN
══════════════════════

  Traditional CRUD:
  ┌─────────────────────────────────┐
  │  Account: user-42               │
  │  Balance: $500                  │  ← Only current state
  │  Last Updated: 2024-01-15      │
  └─────────────────────────────────┘
  History? GONE. Can't reconstruct how we got here.


  Event Sourcing (with Kafka):
  ┌──────────────────────────────────────────────┐
  │  Topic: "account-events"  Key: "user-42"    │
  │                                              │
  │  offset 0: { type: "CREATED",   balance: 0 } │
  │  offset 1: { type: "DEPOSITED", amount: 1000}│
  │  offset 2: { type: "WITHDREW",  amount: 200} │
  │  offset 3: { type: "DEPOSITED", amount: 300} │
  │  offset 4: { type: "WITHDREW",  amount: 600} │
  └──────────────────────────────────────────────┘

  Current balance? Replay events: 0 + 1000 - 200 + 300 - 600 = $500
  Balance at offset 2? Replay up to 2: 0 + 1000 - 200 = $800
  Full audit trail preserved!

  Kafka Topics in Event Sourcing:
  • Event topic (append-only): ALL events ever
  • Snapshot topic (compacted): Current state per key
  • Materialized views: Query-optimized projections
```

### 9.2 CQRS (Command Query Responsibility Segregation)

```
CQRS WITH KAFKA
════════════════

  ┌──────────────┐                              ┌──────────────────┐
  │   Command    │                              │     Query        │
  │   Service    │                              │     Service      │
  │              │                              │                  │
  │  POST /order │                              │  GET /orders     │
  │  PUT /order  │                              │  GET /order/{id} │
  └──────┬───────┘                              └────────┬─────────┘
         │                                               │
         │ write                                    read │
         ▼                                               ▼
  ┌──────────────┐     ┌──────────┐     ┌──────────────────────┐
  │  Write DB    │ ──> │  KAFKA   │ ──> │  Read DB             │
  │  (Postgres)  │     │  (Event  │     │  (Elasticsearch /    │
  │  Normalized  │     │   Bus)   │     │   Redis / MongoDB)   │
  └──────────────┘     └──────────┘     │  Denormalized,       │
                                        │  query-optimized     │
                                        └──────────────────────┘

  Benefits:
  • Write model optimized for writes (normalized, consistent)
  • Read model optimized for reads (denormalized, fast queries)
  • Kafka decouples the two (eventual consistency)
  • Can have MULTIPLE read models for different query patterns
```

### 9.3 Saga Pattern (Distributed Transactions)

```
SAGA PATTERN WITH KAFKA — Choreography
═══════════════════════════════════════

  Scenario: E-commerce order → payment → inventory → shipping

  ┌─────────┐   order    ┌──────────┐   payment   ┌───────────┐   ship
  │  Order  │──created──>│ Payment  │──confirmed──>│ Inventory │──reserved──>
  │ Service │            │ Service  │              │ Service   │
  └─────────┘            └──────────┘              └───────────┘
       ▲                      │                         │
       │                      │ payment                 │ inventory
       │                      │ failed                  │ failed
       │                      ▼                         ▼
       │               ┌──────────┐              ┌───────────┐
       └───cancel──────│  Cancel  │──compensate──│  Restore  │
          order        │  Order   │  inventory   │  Stock    │
                       └──────────┘              └───────────┘

  Each service:
  1. Listens to a Kafka topic for events
  2. Performs its local transaction
  3. Publishes a success or failure event
  4. On failure: publishes compensating events (rollback)

  Topics:
  • order-events:     { CREATED, CONFIRMED, CANCELLED }
  • payment-events:   { CONFIRMED, FAILED }
  • inventory-events: { RESERVED, FAILED, RESTORED }
  • shipping-events:  { SHIPPED, FAILED }
```

### 9.4 Transactional Outbox Pattern

```
TRANSACTIONAL OUTBOX — Reliable Event Publishing
═════════════════════════════════════════════════

  Problem: How to atomically update a DB AND publish to Kafka?

  BAD Approach:
  ─────────────
  1. Save order to DB       ✓
  2. Publish to Kafka       ✗ (if this fails, DB and Kafka are inconsistent)


  OUTBOX Pattern:
  ───────────────
  1. Save order to DB            }
  2. Save event to OUTBOX table  }  Single DB transaction (atomic)

  3. Debezium (CDC) reads outbox table changes → publishes to Kafka
     OR
  3. Background job polls outbox table → publishes to Kafka → marks as sent

  ┌───────────────────────────────────────────────────┐
  │  Database                                          │
  │                                                    │
  │  ┌──────────────┐    ┌─────────────────────────┐  │
  │  │ orders table │    │  outbox table           │  │
  │  │              │    │                         │  │
  │  │ id: 123      │    │ id: 1                   │  │
  │  │ user: alice   │    │ aggregate_type: Order   │  │
  │  │ amount: 100  │    │ aggregate_id: 123       │  │
  │  │              │    │ event_type: OrderCreated │  │
  │  └──────────────┘    │ payload: {...}          │  │
  │                      │ sent: false             │  │
  │  Both written in     └─────────┬───────────────┘  │
  │  ONE transaction                │                  │
  └─────────────────────────────────┼──────────────────┘
                                    │
                          Debezium CDC / Polling
                                    │
                                    ▼
                             ┌──────────┐
                             │  KAFKA   │
                             │  topic   │
                             └──────────┘
```

### 9.5 Dead Letter Queue (DLQ) Pattern

```
DEAD LETTER QUEUE PATTERN
═════════════════════════

  ┌──────────┐     ┌──────────────┐     ┌──────────────┐
  │  Source   │     │   Consumer   │     │  Processing  │
  │  Topic   │────>│              │────>│  Logic       │
  └──────────┘     └──────┬───────┘     └──────┬───────┘
                          │                     │
                        Success               Failure
                          │                     │
                          ▼                     ▼
                    ┌──────────┐          ┌──────────────┐
                    │  Commit  │          │  Retry (3x)  │
                    │  Offset  │          └──────┬───────┘
                    └──────────┘                 │
                                             Still fails
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │  DLT Topic   │
                                          │ "topic.DLT"  │
                                          └──────┬───────┘
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │  Alert +     │
                                          │  Manual      │
                                          │  Investigation│
                                          └──────────────┘

  DLT messages contain:
  • Original message
  • Exception details
  • Retry count
  • Timestamp of failure
  • Source topic and partition
```

### 9.6 Event-Driven Microservices Architecture

```
COMPLETE EVENT-DRIVEN ARCHITECTURE
═══════════════════════════════════

  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │   API    │  │  Order   │  │ Payment  │  │ Shipping │
  │ Gateway  │  │ Service  │  │ Service  │  │ Service  │
  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
       │              │              │              │
       │              ▼              ▼              ▼
       │         ┌─────────────────────────────────────┐
       │         │          KAFKA CLUSTER              │
       │         │                                     │
       │         │  ┌─────────────────────────────┐    │
       └────────>│  │  order-events               │    │
                 │  │  payment-events             │    │
                 │  │  shipping-events            │    │
                 │  │  notification-events        │    │
                 │  │  audit-events               │    │
                 │  └─────────────────────────────┘    │
                 │                                     │
                 └──────────┬──────────────────────────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
       ┌────────────┐ ┌──────────┐ ┌──────────────┐
       │Notification│ │ Analytics│ │   Audit      │
       │  Service   │ │ Service  │ │   Service    │
       │ (email,sms)│ │ (reports)│ │ (compliance) │
       └────────────┘ └──────────┘ └──────────────┘

  Each service:
  • Owns its data (separate database)
  • Publishes events to Kafka (what happened)
  • Subscribes to events it cares about
  • Can be deployed and scaled independently
```

---

## 10. Advanced Spring Boot Kafka Patterns

### 10.1 Custom Serializer/Deserializer

```java
public class OrderEventSerializer implements Serializer<OrderEvent> {

    private final ObjectMapper objectMapper = new ObjectMapper()
        .registerModule(new JavaTimeModule());

    @Override
    public byte[] serialize(String topic, OrderEvent data) {
        try {
            return objectMapper.writeValueAsBytes(data);
        } catch (JsonProcessingException e) {
            throw new SerializationException("Error serializing OrderEvent", e);
        }
    }
}

public class OrderEventDeserializer implements Deserializer<OrderEvent> {

    private final ObjectMapper objectMapper = new ObjectMapper()
        .registerModule(new JavaTimeModule());

    @Override
    public OrderEvent deserialize(String topic, byte[] data) {
        try {
            return objectMapper.readValue(data, OrderEvent.class);
        } catch (IOException e) {
            throw new SerializationException("Error deserializing OrderEvent", e);
        }
    }
}
```

### 10.2 Kafka Headers (Metadata with Messages)

```java
@Service
public class TracingProducer {

    private final KafkaTemplate<String, OrderEvent> kafkaTemplate;

    public void sendWithHeaders(OrderEvent event) {
        ProducerRecord<String, OrderEvent> record = new ProducerRecord<>(
            "order-events", event.getOrderId(), event);

        record.headers()
            .add("correlation-id", UUID.randomUUID().toString().getBytes())
            .add("source-service", "order-service".getBytes())
            .add("event-type", "ORDER_CREATED".getBytes())
            .add("timestamp", String.valueOf(System.currentTimeMillis()).getBytes());

        kafkaTemplate.send(record);
    }
}

@Service
public class TracingConsumer {

    @KafkaListener(topics = "order-events")
    public void consume(ConsumerRecord<String, OrderEvent> record) {
        String correlationId = new String(
            record.headers().lastHeader("correlation-id").value());
        String sourceService = new String(
            record.headers().lastHeader("source-service").value());

        MDC.put("correlationId", correlationId);
        log.info("Processing order from {} with correlationId {}",
            sourceService, correlationId);
    }
}
```

### 10.3 Batch Consumer

```java
@Configuration
public class BatchConsumerConfig {

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, OrderEvent>
            batchFactory(ConsumerFactory<String, OrderEvent> cf) {

        ConcurrentKafkaListenerContainerFactory<String, OrderEvent> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(cf);
        factory.setBatchListener(true);
        factory.setConcurrency(3);
        factory.getContainerProperties().setAckMode(
            ContainerProperties.AckMode.BATCH);
        return factory;
    }
}

@Service
@Slf4j
public class BatchOrderConsumer {

    @KafkaListener(
        topics = "order-events",
        groupId = "batch-processor",
        containerFactory = "batchFactory"
    )
    public void consumeBatch(List<ConsumerRecord<String, OrderEvent>> records,
                             Acknowledgment ack) {
        log.info("Received batch of {} records", records.size());

        List<OrderEvent> orders = records.stream()
            .map(ConsumerRecord::value)
            .toList();

        batchInsertToDatabase(orders);
        ack.acknowledge();
    }
}
```

### 10.4 Conditional Consumer (Kafka Filter)

```java
@Configuration
public class FilteredConsumerConfig {

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, OrderEvent>
            filteredFactory(ConsumerFactory<String, OrderEvent> cf) {

        ConcurrentKafkaListenerContainerFactory<String, OrderEvent> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(cf);

        factory.setRecordFilterStrategy(record ->
            record.value().getAmount() < 100.0
        );

        return factory;
    }
}

@KafkaListener(
    topics = "order-events",
    containerFactory = "filteredFactory"
)
public void consumeHighValueOrders(OrderEvent event) {
    log.info("Processing high-value order: {}", event.getOrderId());
}
```

### 10.5 Request-Reply Pattern with Kafka

```java
@Configuration
public class ReplyingKafkaConfig {

    @Bean
    public ReplyingKafkaTemplate<String, String, String> replyingTemplate(
            ProducerFactory<String, String> pf,
            ConcurrentMessageListenerContainer<String, String> repliesContainer) {

        return new ReplyingKafkaTemplate<>(pf, repliesContainer);
    }

    @Bean
    public ConcurrentMessageListenerContainer<String, String> repliesContainer(
            ConcurrentKafkaListenerContainerFactory<String, String> factory) {

        ConcurrentMessageListenerContainer<String, String> container =
            factory.createContainer("replies-topic");
        container.getContainerProperties().setGroupId("replies-group");
        container.setAutoStartup(false);
        return container;
    }
}

@RestController
public class SyncKafkaController {

    private final ReplyingKafkaTemplate<String, String, String> replyingTemplate;

    @GetMapping("/api/validate-order/{orderId}")
    public ResponseEntity<String> validateOrder(@PathVariable String orderId)
            throws Exception {

        ProducerRecord<String, String> record =
            new ProducerRecord<>("validation-requests", orderId);
        record.headers().add(
            KafkaHeaders.REPLY_TOPIC, "replies-topic".getBytes());

        RequestReplyFuture<String, String, String> future =
            replyingTemplate.sendAndReceive(record, Duration.ofSeconds(10));

        ConsumerRecord<String, String> reply = future.get(10, TimeUnit.SECONDS);
        return ResponseEntity.ok(reply.value());
    }
}

@Service
public class ValidationService {

    @KafkaListener(topics = "validation-requests")
    @SendTo
    public String validateOrder(String orderId) {
        boolean isValid = performValidation(orderId);
        return isValid ? "VALID" : "INVALID: " + getReason(orderId);
    }
}
```

### 10.6 Retry with Exponential Backoff and DLT

```java
@Configuration
public class RetryConfig {

    @Bean
    public DefaultErrorHandler errorHandler(
            KafkaTemplate<String, Object> kafkaTemplate) {

        DeadLetterPublishingRecoverer recoverer =
            new DeadLetterPublishingRecoverer(kafkaTemplate,
                (record, ex) -> new TopicPartition(
                    record.topic() + ".DLT", record.partition()));

        ExponentialBackOff backOff = new ExponentialBackOff(1000L, 2.0);
        backOff.setMaxElapsedTime(60000L);

        DefaultErrorHandler errorHandler =
            new DefaultErrorHandler(recoverer, backOff);

        errorHandler.addNotRetryableExceptions(
            DeserializationException.class,
            ClassCastException.class
        );

        return errorHandler;
    }
}

@Service
@Slf4j
public class DltConsumer {

    @KafkaListener(topics = "order-events.DLT", groupId = "dlt-processor")
    public void processDlt(ConsumerRecord<String, OrderEvent> record) {
        log.error("DLT received: key={}, value={}, exception={}",
            record.key(), record.value(),
            new String(record.headers().lastHeader(
                "kafka_dlt-exception-message").value()));

        alertOpsTeam(record);
        saveToFailedEventsTable(record);
    }
}
```

### 10.7 Kafka Transaction with Spring Boot

```java
@Configuration
public class KafkaTransactionConfig {

    @Bean
    public ProducerFactory<String, Object> producerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ProducerConfig.TRANSACTIONAL_ID_CONFIG, "tx-");
        config.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);

        DefaultKafkaProducerFactory<String, Object> factory =
            new DefaultKafkaProducerFactory<>(config);
        factory.setTransactionIdPrefix("tx-");
        return factory;
    }

    @Bean
    public KafkaTransactionManager<String, Object> kafkaTransactionManager(
            ProducerFactory<String, Object> pf) {
        return new KafkaTransactionManager<>(pf);
    }
}

@Service
@RequiredArgsConstructor
public class TransactionalOrderService {

    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Transactional
    public void processOrder(OrderEvent order) {
        kafkaTemplate.send("order-events", order.getOrderId(), order);
        kafkaTemplate.send("audit-events", order.getOrderId(),
            new AuditEvent("ORDER_PROCESSED", order.getOrderId()));
        kafkaTemplate.send("notification-events", order.getUserId(),
            new NotificationEvent("Your order " + order.getOrderId() + " is confirmed"));
    }
}
```

---

## 11. Troubleshooting Common Issues

### 11.1 Consumer Not Receiving Messages

```
TROUBLESHOOTING: Consumer Not Receiving Messages
═════════════════════════════════════════════════

Check 1: Is the consumer in the right group?
  → kafka-consumer-groups --describe --group <group-id>
  → Verify topic and partition assignments

Check 2: What is auto.offset.reset?
  → If "latest" and no committed offset → consumer only sees NEW messages
  → Solution: Set to "earliest" or produce new messages

Check 3: Is the consumer assigned to any partitions?
  → If consumers > partitions → some are IDLE
  → Solution: Increase partitions or reduce consumers

Check 4: Is there a deserialization error?
  → Check logs for SerializationException
  → Producer and consumer must use same serializer/deserializer

Check 5: Is the consumer lagging?
  → kafka-consumer-groups --describe --group <group-id>
  → Check LAG column

Check 6: Is the consumer repeatedly rebalancing?
  → Check for: max.poll.interval.ms too low
  → Processing taking too long between poll() calls
  → Solution: Increase max.poll.interval.ms or reduce max.poll.records
```

### 11.2 Consumer Rebalancing Too Often

```
TROUBLESHOOTING: Frequent Rebalances
═════════════════════════════════════

Symptoms:
  • Consumer keeps printing "Revoke partitions" / "Assign partitions"
  • Messages are processed multiple times
  • Throughput drops

Causes & Solutions:
  ┌─────────────────────────────────────┬──────────────────────────────┐
  │  Cause                              │  Solution                    │
  ├─────────────────────────────────────┼──────────────────────────────┤
  │ Processing time > max.poll.interval │ ↑ max.poll.interval.ms       │
  │                                     │ ↓ max.poll.records           │
  ├─────────────────────────────────────┼──────────────────────────────┤
  │ Session timeout too aggressive      │ ↑ session.timeout.ms         │
  │                                     │ (set heartbeat = timeout/3)  │
  ├─────────────────────────────────────┼──────────────────────────────┤
  │ Consumer instances flapping         │ Fix deployment/health issues │
  │ (crash-loop, OOM, etc.)             │                              │
  ├─────────────────────────────────────┼──────────────────────────────┤
  │ Too many consumers joining/leaving  │ Use static group membership  │
  │                                     │ (group.instance.id)          │
  └─────────────────────────────────────┴──────────────────────────────┘

Static Group Membership (reduce rebalances):
  spring.kafka.consumer.properties.group.instance.id=consumer-1
  → Consumer gets a stable identity
  → Rejoining after restart doesn't trigger full rebalance
  → Partitions are preserved if consumer rejoins within session.timeout
```

### 11.3 Producer Message Loss

```
TROUBLESHOOTING: Message Loss
══════════════════════════════

Check 1: What is your acks setting?
  acks=0  → Messages may be lost (fire and forget)
  acks=1  → Message may be lost if leader crashes before replication
  acks=all → Safest (combine with min.insync.replicas ≥ 2)

Check 2: Are retries enabled?
  retries=0 → Transient failures cause permanent loss
  Solution: retries=3 (or higher), with enable.idempotence=true

Check 3: Is buffer.memory full?
  If buffer is full → block.on.buffer.full or max.block.ms applies
  After timeout → BufferExhaustedException → message lost
  Solution: Increase buffer.memory or handle the exception

Check 4: Are you handling send() errors?
  kafkaTemplate.send().whenComplete((result, ex) -> {
      if (ex != null) {
          // THIS MUST NOT BE IGNORED
          log.error("Send failed!", ex);
          retryOrSaveToFallback(message);
      }
  });

SAFE PRODUCER CHECKLIST:
  ✓ acks=all
  ✓ enable.idempotence=true
  ✓ retries=3+
  ✓ min.insync.replicas=2 (broker/topic config)
  ✓ Handle send() errors in callback
  ✓ Monitor producer metrics
```

### 11.4 High Consumer Lag

```
TROUBLESHOOTING: Growing Consumer Lag
══════════════════════════════════════

  Diagnosis:
  kafka-consumer-groups --describe --group my-group
  → If LAG is growing → consumers can't keep up

  Solutions (in order of effort):
  ┌────────────────────────────────────────────────────────────────┐
  │  1. Increase consumers (up to partition count)                │
  │     → Add more instances of your consumer app                 │
  │                                                               │
  │  2. Increase max.poll.records                                 │
  │     → Process more records per poll cycle                     │
  │                                                               │
  │  3. Optimize processing logic                                 │
  │     → Batch DB writes, use async I/O, reduce processing time  │
  │                                                               │
  │  4. Increase partitions (requires producer-side changes)       │
  │     → More partitions = more consumer parallelism             │
  │                                                               │
  │  5. Use batch consumption                                     │
  │     → Process records in batches instead of one-by-one        │
  │                                                               │
  │  6. Check for processing bottlenecks                          │
  │     → Database writes slow? External API slow? Memory issues? │
  └────────────────────────────────────────────────────────────────┘
```

### 11.5 Common Exceptions and Fixes

| Exception | Cause | Fix |
|-----------|-------|-----|
| `TimeoutException` | Broker unreachable | Check network, bootstrap-servers config |
| `SerializationException` | Incompatible serializer/deserializer | Match producer and consumer serializers |
| `RecordTooLargeException` | Message > `max.message.bytes` | Increase broker limit or reduce message size |
| `NotEnoughReplicasException` | ISR count < `min.insync.replicas` | Check broker health, wait for recovery |
| `OffsetOutOfRangeException` | Requested offset doesn't exist | Set `auto.offset.reset=earliest` |
| `CommitFailedException` | Consumer was kicked out during commit | Increase `max.poll.interval.ms` |
| `RebalanceInProgressException` | Consumer group is rebalancing | Retry or increase `session.timeout.ms` |
| `InvalidGroupIdException` | Missing or invalid `group.id` | Set `spring.kafka.consumer.group-id` |
| `TopicAuthorizationException` | No ACL permission | Grant proper ACLs |

---

## 12. Advanced Interview Questions

### Q1: How does Kafka achieve high throughput despite writing to disk?

**Answer:**  
Kafka uses several techniques:
1. **Sequential I/O** — Kafka writes to the end of log files sequentially, which is extremely fast on modern disks (500MB/s+). Random I/O is slow (~100 IOPS).
2. **Zero-copy transfer** — Data goes from disk to network socket without copying through the application (uses `sendfile()` system call).
3. **Batching** — Messages are batched both on producer (record accumulator) and broker (segment files), reducing I/O operations.
4. **Page cache** — Kafka relies on the OS page cache rather than maintaining its own cache in JVM heap. Hot data is served from memory without explicit caching code.
5. **Compression** — Messages can be compressed (snappy/lz4) in batches, reducing both disk and network I/O.

### Q2: Explain how Kafka ensures exactly-once semantics internally.

**Answer:**  
Kafka's exactly-once relies on two mechanisms:

**Idempotent Producer:**
- Each producer gets a unique **Producer ID (PID)** from the broker
- Each message gets a **sequence number** per partition
- Broker maintains the latest sequence number per PID per partition
- If a duplicate arrives (same PID + sequence), broker silently drops it
- Handles retries without duplicates

**Transactions:**
- Producer wraps multiple sends in a transaction
- Uses a **Transaction Coordinator** (a broker) to manage transaction state
- Writes transaction markers (COMMIT/ABORT) to each partition
- Consumers with `isolation.level=read_committed` skip uncommitted/aborted messages
- Atomic: all messages in a transaction are visible, or none are

### Q3: What is a Consumer Group Coordinator and how does it work?

**Answer:**  
The Group Coordinator is a specific broker responsible for managing a consumer group:

1. **Selection:** The coordinator for a group is determined by hashing the `group.id` and mapping it to a partition of `__consumer_offsets`. The broker leading that partition becomes the coordinator.
2. **Responsibilities:**
   - Receiving `JoinGroup` and `SyncGroup` requests
   - Managing the rebalance protocol
   - Storing committed offsets in `__consumer_offsets`
   - Monitoring consumer heartbeats
   - Declaring consumers dead after `session.timeout.ms`
3. **Rebalance Protocol:**
   - Consumer sends `JoinGroup` → coordinator creates a new "generation"
   - One consumer is elected as the **Group Leader** (client-side)
   - Leader computes partition assignments using the configured `PartitionAssignor`
   - Leader sends assignments back via `SyncGroup`
   - Coordinator distributes assignments to all consumers

### Q4: Explain the difference between Log Compaction and Log Deletion.

**Answer:**

| Aspect | Log Deletion | Log Compaction |
|--------|-------------|----------------|
| **Trigger** | Time (`retention.ms`) or size (`retention.bytes`) | Dirty ratio (`min.cleanable.dirty.ratio`) |
| **What's removed** | Entire old segments | Only older duplicate keys |
| **Guarantee** | "Data younger than X" | "At least latest per key" |
| **Key required** | No | Yes |
| **Use case** | Event logs, metrics | State stores, snapshots |
| **Config** | `cleanup.policy=delete` | `cleanup.policy=compact` |

### Q5: How would you design a Kafka-based system for real-time fraud detection?

**Answer:**

```
Design:
1. Transaction events → "transactions" topic (partitioned by card-id)
2. Kafka Streams application:
   - Windowed aggregation: count transactions per card per 5-minute window
   - State store: maintains spending patterns per card
   - KStream-KTable join: enrich with customer risk profile
   - Rules engine: flag if > 5 transactions in 1 minute,
     or amount > 3x average, or unusual geography
3. Flagged events → "fraud-alerts" topic
4. Alert service consumes fraud-alerts → notifies fraud team

Key decisions:
- Partition by card-id → all transactions for a card in one partition
  → ordering guaranteed, stateful processing works
- Use exactly-once semantics (can't double-flag or miss flags)
- Low latency: Kafka Streams (ms latency)
- State backed by changelog topic (fault-tolerant)
```

### Q6: What is the Cooperative Sticky Assignor and why is it better?

**Answer:**

**Eager Rebalancing (Default before Kafka 2.4):**
- ALL consumers revoke ALL partitions at start of rebalance
- ALL partitions are reassigned from scratch
- Brief period where NO consumer is processing anything
- High impact even for minor changes (one consumer joining/leaving)

**Cooperative Sticky Assignor (Kafka 2.4+):**
- Only partitions that need to MOVE are revoked
- Other consumers keep processing their existing partitions
- Rebalancing happens in multiple phases (incremental)
- Minimal disruption

```
Eager: Consumer A has P0,P1 | Consumer B has P2
       Consumer C joins → ALL stop → reassign → A:P0, B:P1, C:P2
       Downtime: ~seconds (all stopped)

Cooperative: Consumer A has P0,P1 | Consumer B has P2
             Consumer C joins → only P1 is revoked from A → assigned to C
             A keeps processing P0, B keeps processing P2
             Downtime: near-zero
```

**Configuration:**

```yaml
spring:
  kafka:
    consumer:
      properties:
        partition.assignment.strategy: org.apache.kafka.clients.consumer.CooperativeStickyAssignor
```

### Q7: How do you handle schema evolution in a Kafka-based system?

**Answer:**
1. Use **Schema Registry** with **Avro** or **Protobuf**
2. Set compatibility mode (BACKWARD recommended as default)
3. Rules:
   - Always add new fields with **default values**
   - Never remove required fields without a default
   - Never change field types
   - Never rename fields
4. Deploy consumers first (BACKWARD), then producers
5. Schema Registry validates compatibility on registration — incompatible schemas are rejected
6. Use specific Avro reader (`specific.avro.reader=true`) for type-safe POJOs, or GenericRecord for flexible handling

### Q8: What are ISR shrink and expand events, and what do they indicate?

**Answer:**

- **ISR Shrink:** A follower is removed from the ISR because it fell behind the leader beyond `replica.lag.time.max.ms` (default 30s). Indicates:
  - Network issues between brokers
  - Disk I/O bottleneck on the follower
  - GC pauses on the follower
  - Follower broker overloaded

- **ISR Expand:** A previously removed follower caught up and rejoined the ISR. Indicates the issue was resolved.

- **Monitoring:** Frequent shrink/expand cycles indicate an unstable cluster. Alert on `UnderReplicatedPartitions > 0` and investigate broker health.

### Q9: How does Kafka handle backpressure?

**Answer:**  
Kafka handles backpressure naturally through its **pull-based** consumer model:

1. **Consumer side:** Consumers call `poll()` at their own pace. If processing is slow, they simply poll less frequently. Messages stay in Kafka until consumed (within retention period).
2. **Producer side:** If brokers are slow, the producer's internal buffer fills up. When `buffer.memory` is exhausted, `send()` blocks for up to `max.block.ms` and then throws `BufferExhaustedException`.
3. **Broker side:** Quotas can be set per client to limit produce/consume bandwidth:
   ```bash
   kafka-configs --alter --add-config 'producer_byte_rate=1048576,consumer_byte_rate=2097152' \
     --entity-type clients --entity-name my-app --bootstrap-server localhost:9092
   ```

### Q10: Explain Kafka's internal storage format.

**Answer:**

```
Topic "orders" on disk:
/kafka-logs/
  └── orders-0/                  ← Partition 0
      ├── 00000000000000000000.log       ← Segment file (messages)
      ├── 00000000000000000000.index     ← Offset index
      ├── 00000000000000000000.timeindex ← Timestamp index
      ├── 00000000000049152.log          ← Next segment (starts at offset 49152)
      ├── 00000000000049152.index
      ├── 00000000000049152.timeindex
      └── leader-epoch-checkpoint
  └── orders-1/                  ← Partition 1
      ├── ...
```

- **Segment files (.log):** Actual message data. Split by size (`log.segment.bytes`, default 1GB) or time (`log.roll.ms`).
- **Offset index (.index):** Maps offset → physical position in the .log file. Sparse index (not every offset).
- **Time index (.timeindex):** Maps timestamp → offset. Used for time-based lookups.
- **Active segment:** The latest segment being written to. Only closed segments can be deleted or compacted.

---

## 13. Quick Reference Cheat Sheet — Part 2

### Kafka Streams Terminology

| Term | One-Line Definition |
|------|---------------------|
| **KStream** | An unbounded stream of events (all records, append-only) |
| **KTable** | A changelog stream (latest value per key, like a table) |
| **GlobalKTable** | A KTable fully replicated on every stream instance |
| **Topology** | The processing graph (sources → processors → sinks) |
| **State Store** | Local key-value storage (RocksDB) for stateful operations |
| **Windowed Operation** | Aggregation over time-bounded groups of records |
| **Punctuator** | Timer-based callback within a processor |
| **Interactive Query** | REST API to query local state stores |

### Kafka Connect Terminology

| Term | One-Line Definition |
|------|---------------------|
| **Source Connector** | Reads from external system, writes to Kafka |
| **Sink Connector** | Reads from Kafka, writes to external system |
| **Worker** | JVM process that runs connector tasks |
| **Task** | Unit of parallelism within a connector |
| **SMT** | Single Message Transform — lightweight per-record transformation |
| **Converter** | Translates between internal format and serialized format |
| **CDC** | Change Data Capture — capture database changes in real-time |
| **Debezium** | Most popular open-source CDC platform for Kafka Connect |

### Schema Registry Terminology

| Term | One-Line Definition |
|------|---------------------|
| **Subject** | A named scope for schema evolution (usually `<topic>-value`) |
| **Schema ID** | Unique numeric ID for a registered schema |
| **Compatibility** | Rules governing allowed schema changes (BACKWARD, FORWARD, FULL) |
| **Avro** | Compact binary serialization format with schema evolution support |
| **Protobuf** | Google's binary serialization format (smallest, fastest) |
| **Schema Evolution** | Changing a schema over time while maintaining compatibility |

### Advanced CLI Commands

```bash
# Kafka Streams
kafka-streams-application-reset --application-id my-stream-app \
  --input-topics input-topic --bootstrap-server localhost:9092

# Reassign partitions (JSON file based)
kafka-reassign-partitions --bootstrap-server localhost:9092 \
  --reassignment-json-file reassignment.json --execute

# Change topic replication factor
kafka-reassign-partitions --bootstrap-server localhost:9092 \
  --reassignment-json-file increase-replication.json --execute

# View log segments
kafka-dump-log --files /kafka-logs/orders-0/00000000000000000000.log \
  --print-data-log

# Check broker configs
kafka-configs --describe --entity-type brokers --entity-name 0 \
  --bootstrap-server localhost:9092

# Producer/Consumer performance test
kafka-producer-perf-test --topic perf-test --num-records 1000000 \
  --record-size 1024 --throughput -1 --producer-props bootstrap.servers=localhost:9092

kafka-consumer-perf-test --topic perf-test --messages 1000000 \
  --bootstrap-server localhost:9092
```

### Docker Compose — Complete Kafka Ecosystem

```yaml
version: '3.8'
services:
  kafka:
    image: apache/kafka:3.7.0
    ports:
      - "9092:9092"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093,INTERNAL://0.0.0.0:29092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092,INTERNAL://kafka:29092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,INTERNAL:PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@localhost:9093
      KAFKA_INTER_BROKER_LISTENER_NAME: INTERNAL
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk

  schema-registry:
    image: confluentinc/cp-schema-registry:7.5.0
    depends_on:
      - kafka
    ports:
      - "8081:8081"
    environment:
      SCHEMA_REGISTRY_HOST_NAME: schema-registry
      SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS: kafka:29092

  kafka-connect:
    image: confluentinc/cp-kafka-connect:7.5.0
    depends_on:
      - kafka
      - schema-registry
    ports:
      - "8083:8083"
    environment:
      CONNECT_BOOTSTRAP_SERVERS: kafka:29092
      CONNECT_REST_PORT: 8083
      CONNECT_GROUP_ID: connect-cluster
      CONNECT_CONFIG_STORAGE_TOPIC: connect-configs
      CONNECT_OFFSET_STORAGE_TOPIC: connect-offsets
      CONNECT_STATUS_STORAGE_TOPIC: connect-status
      CONNECT_KEY_CONVERTER: org.apache.kafka.connect.json.JsonConverter
      CONNECT_VALUE_CONVERTER: io.confluent.connect.avro.AvroConverter
      CONNECT_VALUE_CONVERTER_SCHEMA_REGISTRY_URL: http://schema-registry:8081
      CONNECT_CONFIG_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_OFFSET_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_STATUS_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_PLUGIN_PATH: /usr/share/java,/usr/share/confluent-hub-components

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    depends_on:
      - kafka
      - schema-registry
      - kafka-connect
    ports:
      - "8080:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: local
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
      KAFKA_CLUSTERS_0_SCHEMAREGISTRY: http://schema-registry:8081
      KAFKA_CLUSTERS_0_KAFKACONNECT_0_NAME: connect
      KAFKA_CLUSTERS_0_KAFKACONNECT_0_ADDRESS: http://kafka-connect:8083
```

### Production Checklist

```
KAFKA PRODUCTION READINESS CHECKLIST
═════════════════════════════════════

Cluster:
  ☐ Minimum 3 brokers
  ☐ KRaft mode (no ZooKeeper for new deployments)
  ☐ Dedicated disks for Kafka log directories
  ☐ Separate controller and broker roles (for large clusters)

Topics:
  ☐ replication-factor >= 3
  ☐ min.insync.replicas = 2
  ☐ Appropriate partition count (not too few, not too many)
  ☐ retention.ms set based on use case

Producers:
  ☐ acks=all
  ☐ enable.idempotence=true
  ☐ retries > 0
  ☐ Error handling in send callbacks
  ☐ Compression enabled (lz4 or snappy)

Consumers:
  ☐ enable.auto.commit=false (manual commit)
  ☐ Proper auto.offset.reset (earliest for new groups)
  ☐ max.poll.interval.ms matches processing time
  ☐ Idempotent processing (handle duplicates)
  ☐ Dead Letter Topic configured

Security:
  ☐ SSL/TLS for encryption in transit
  ☐ SASL authentication enabled
  ☐ ACLs configured per service
  ☐ Separate credentials per application

Monitoring:
  ☐ Consumer lag monitoring with alerting
  ☐ UnderReplicatedPartitions alert
  ☐ OfflinePartitionsCount alert
  ☐ Broker disk usage monitoring
  ☐ JMX metrics exported to Prometheus/Grafana

Operations:
  ☐ Kafka UI deployed for visibility
  ☐ Schema Registry for schema management
  ☐ Backup strategy for critical topics
  ☐ Runbooks for common failure scenarios
  ☐ Regular partition rebalancing
```

---

> **Further Learning:**
> - [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
> - [Confluent Developer](https://developer.confluent.io/)
> - [Kafka Streams Documentation](https://kafka.apache.org/documentation/streams/)
> - [Debezium Documentation](https://debezium.io/documentation/)
> - Practice: Build a complete event-driven microservices project with Kafka Streams, Connect, and Schema Registry
