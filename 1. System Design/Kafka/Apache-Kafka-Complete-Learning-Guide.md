# Apache Kafka Complete Learning Guide

> A structured beginner-to-intermediate guide for software engineers to understand, use, and master Apache Kafka.  
> **Focus:** Core concepts, internal workings, Spring Boot integration, and interview readiness.  
> **Audience:** Java/Spring Boot developers with zero prior Kafka knowledge.

---

## Table of Contents

1. [What is Apache Kafka?](#1-what-is-apache-kafka)
2. [Why Kafka? Real-World Problems It Solves](#2-why-kafka-real-world-problems-it-solves)
3. [Kafka vs Traditional Messaging Systems](#3-kafka-vs-traditional-messaging-systems)
4. [Core Concepts](#4-core-concepts)
5. [Message Flow: End to End](#5-message-flow-end-to-end)
6. [Kafka Architecture Diagrams](#6-kafka-architecture-diagrams)
7. [Deep Understanding: Why Things Work This Way](#7-deep-understanding-why-things-work-this-way)
8. [Delivery Semantics](#8-delivery-semantics)
9. [Practical Section: Running Kafka Locally](#9-practical-section-running-kafka-locally)
10. [Spring Boot Integration](#10-spring-boot-integration)
11. [Important Kafka Configurations](#11-important-kafka-configurations)
12. [Common Interview Questions](#12-common-interview-questions)
13. [Quick Reference Cheat Sheet](#13-quick-reference-cheat-sheet)

---

## 1. What is Apache Kafka?

### Simple Definition

Apache Kafka is an **open-source distributed event streaming platform**. In plain English, it is a system that lets applications **send**, **store**, and **read** streams of data (messages/events) in real-time, at massive scale.

- Originally developed at **LinkedIn** in 2010
- Open-sourced through the **Apache Software Foundation**
- Written in **Scala** and **Java**
- Used by 80%+ of Fortune 100 companies

### The Post Office Analogy

Think of Kafka like a **super-efficient post office**:

```
YOU (Producer)                   POST OFFICE (Kafka)                FRIEND (Consumer)
     │                                  │                                │
     │   Write a letter                 │                                │
     │   Put it in an envelope          │                                │
     │   Address it to a TOPIC          │                                │
     │ ────────────────────────────────>│                                │
     │                                  │   Sorts into the right         │
     │                                  │   mailbox (Partition)          │
     │                                  │   Keeps a copy (Durability)    │
     │                                  │ ──────────────────────────────>│
     │                                  │                   Reads letter │
     │                                  │                   at own pace  │
```

**Key differences from a real post office:**
- The post office **keeps every letter forever** (or for a configured time)
- Multiple friends can **read the same letter** independently
- Letters are **numbered** (offsets), so your friend knows exactly where they left off
- The post office can handle **millions of letters per second**

### One-Line Summary

> Kafka = A distributed commit log that lets producers write events and consumers read them, independently, at their own pace, with guaranteed ordering and durability.

---

## 2. Why Kafka? Real-World Problems It Solves

### Problems BEFORE Kafka

```
Without Kafka: Point-to-Point Chaos
────────────────────────────────────

  Service A ──────> Service D
  Service A ──────> Service E
  Service B ──────> Service D
  Service B ──────> Service F
  Service C ──────> Service E
  Service C ──────> Service F

  6 direct connections! (grows as N×M)
  Each connection = different protocol, different format, tight coupling
```

```
With Kafka: Clean Decoupled Architecture
─────────────────────────────────────────

  Service A ──┐                    ┌──> Service D
  Service B ──┼──> [ KAFKA ] ─────┼──> Service E
  Service C ──┘                    └──> Service F

  All services talk to ONE system
  Loose coupling, easy to add/remove services
```

### Real-World Use Cases

| Use Case | Example | Why Kafka? |
|----------|---------|------------|
| **Event Streaming** | User activity tracking at Netflix | Millions of events/sec, real-time |
| **Log Aggregation** | Collecting logs from 1000+ servers | Central pipeline, no data loss |
| **Metrics Collection** | Application performance monitoring | Time-series data, high throughput |
| **Order Processing** | E-commerce order pipeline | Guaranteed delivery, ordering |
| **Data Integration** | Syncing data between MySQL → Elasticsearch | Decoupled, reliable data flow |
| **Microservice Communication** | Async communication between services | Loose coupling, fault tolerance |
| **Fraud Detection** | Real-time credit card transaction analysis | Low latency stream processing |
| **IoT Data Ingestion** | Sensor data from millions of devices | Massive scale, buffering |

### Why Not Just Use a Database or REST API?

| Approach | Limitation Kafka Solves |
|----------|------------------------|
| **REST API (sync)** | Caller waits for response; if receiver is down, data is lost |
| **Database polling** | Expensive; not real-time; tight coupling to DB schema |
| **Traditional MQ (RabbitMQ)** | Message deleted after consumption; no replay; lower throughput |
| **Kafka** | Async, durable, replayable, high throughput, decoupled |

---

## 3. Kafka vs Traditional Messaging Systems

| Feature | Traditional MQ (RabbitMQ, ActiveMQ) | Apache Kafka |
|---------|-------------------------------------|--------------|
| **Model** | Message Queue (point-to-point) | Distributed Commit Log |
| **Message Retention** | Deleted after consumption | Retained for configured period (even after consumption) |
| **Replay** | Not possible (message gone) | Possible (consumers can re-read) |
| **Throughput** | Thousands/sec | **Millions/sec** |
| **Ordering** | No guaranteed ordering | **Guaranteed within partition** |
| **Consumer Model** | Push (broker pushes to consumer) | **Pull** (consumer pulls from broker) |
| **Scaling** | Vertical (limited) | **Horizontal** (add more brokers) |
| **Use Case** | Task queues, request-reply | Event streaming, data pipelines |
| **Message Routing** | Complex routing rules (exchanges) | Simple topic-based |
| **Backpressure** | Can overwhelm slow consumers | Consumers read at their own pace |
| **Storage** | In-memory (mostly) | **Disk-based** (sequential I/O) |
| **Cluster** | Complex setup | **Built-in** distributed design |

### When to Use What?

- **Use RabbitMQ** when: You need complex routing, request-reply pattern, small scale
- **Use Kafka** when: You need high throughput, event replay, data pipelines, stream processing

---

## 4. Core Concepts

### 4.1 Topic

**Simple Definition:**  
A topic is a **named category or feed** to which messages are published. Think of it as a **folder** or a **channel name**.

**Real-World Analogy:**  
A topic is like a **TV channel**. "Sports" channel has sports content, "News" channel has news. Producers broadcast to a specific channel, and consumers tune into the channels they care about.

**Internal Working:**
- A topic is a logical name (e.g., `order-events`, `user-signups`)
- Each topic is split into one or more **partitions**
- Topics are identified by their name (must be unique within a cluster)
- Topics can be created manually or auto-created on first use

**Example:**
```
Topic: "order-events"
  ├── Partition 0: [msg0, msg1, msg2, msg3]
  ├── Partition 1: [msg0, msg1, msg2]
  └── Partition 2: [msg0, msg1, msg2, msg3, msg4]
```

---

### 4.2 Partition

**Simple Definition:**  
A partition is a **numbered, ordered, immutable sequence of messages** within a topic. It is the unit of parallelism in Kafka.

**Real-World Analogy:**  
Imagine a topic "Complaints" at a bank. Instead of one long queue, they open **3 counters** (partitions). Customers are distributed across counters. Each counter maintains its own order. More counters = faster processing.

**Internal Working:**
- Each partition is an **append-only log** stored on disk
- Messages within a partition get a sequential ID called **offset**
- Messages are **ordered within a partition**, but NOT across partitions
- Partitions are distributed across different brokers for scalability
- You choose the number of partitions when creating a topic

**Why Partitions Matter:**
- **Parallelism:** More partitions = more consumers can read in parallel
- **Scalability:** Partitions spread across brokers = distributed load
- **Throughput:** Each partition can be read/written independently

```
Topic: "order-events" (3 partitions)

Partition 0:  [0] [1] [2] [3] [4] [5]  ───>  (append direction)
Partition 1:  [0] [1] [2] [3]          ───>
Partition 2:  [0] [1] [2] [3] [4]      ───>

Each [] is a message. The number is the offset.
New messages are ALWAYS appended to the end.
Old messages are NEVER modified.
```

---

### 4.3 Broker

**Simple Definition:**  
A broker is a **single Kafka server** that stores data and serves client requests. A Kafka cluster is a group of brokers working together.

**Real-World Analogy:**  
A broker is like a **warehouse in a logistics network**. Amazon doesn't have one giant warehouse — they have many across the country. Each warehouse (broker) stores a portion of the inventory (partitions). Together, they form the delivery network (cluster).

**Internal Working:**
- Each broker is identified by a unique **ID** (integer)
- A broker stores one or more **partitions** of one or more topics
- One broker in the cluster acts as the **Controller** (manages partition assignments)
- Brokers communicate with each other for replication
- Any broker can handle a client's request (it knows which broker has the data)

```
Kafka Cluster (3 Brokers)
═════════════════════════

┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Broker 0      │  │   Broker 1      │  │   Broker 2      │
│                 │  │                 │  │                 │
│  Topic-A P0 (L) │  │  Topic-A P1 (L) │  │  Topic-A P2 (L) │
│  Topic-A P1 (F) │  │  Topic-A P2 (F) │  │  Topic-A P0 (F) │
│  Topic-B P0 (L) │  │  Topic-B P1 (L) │  │  Topic-B P0 (F) │
│                 │  │                 │  │  Topic-B P1 (F) │
└─────────────────┘  └─────────────────┘  └─────────────────┘

(L) = Leader partition   (F) = Follower/Replica partition
```

---

### 4.4 Producer

**Simple Definition:**  
A producer is a **client application** that sends (publishes/writes) messages to Kafka topics.

**Real-World Analogy:**  
A producer is like a **news reporter**. They write a news story (message) and submit it to a specific section (topic) of the newspaper. They don't care who reads it or when.

**Internal Working:**
1. Producer connects to any broker (called **bootstrap server**)
2. Asks for **metadata** (which broker leads which partition)
3. Sends messages directly to the **leader broker** of the target partition
4. Can choose partition via:
   - **Round Robin** (default, no key) — distributes evenly
   - **Key-based hashing** — same key always goes to same partition
   - **Custom partitioner** — your own logic

**Key Decision: With or Without a Key?**

| Scenario | Key | Partition Strategy | Use When |
|----------|-----|--------------------|----------|
| Even distribution | No key (null) | Round Robin | Order doesn't matter |
| Related messages together | Key = `userId` | Hash(key) % partitions | Need ordering per user |
| Custom logic | Custom | Custom Partitioner class | Special routing rules |

**Example:**
```java
// Without key — round robin across partitions
producer.send(new ProducerRecord<>("order-events", "Order #123 placed"));

// With key — all orders for user-42 go to the same partition
producer.send(new ProducerRecord<>("order-events", "user-42", "Order #123 placed"));
```

---

### 4.5 Consumer

**Simple Definition:**  
A consumer is a **client application** that reads (subscribes to/pulls) messages from Kafka topics.

**Real-World Analogy:**  
A consumer is like a **newspaper subscriber**. They subscribe to sections (topics) they're interested in. They read at their own pace. If they miss a day, they can go back and read old issues (replay). Multiple subscribers can read the same newspaper independently.

**Internal Working:**
1. Consumer subscribes to one or more topics
2. Consumer **pulls** messages from the broker (Kafka does NOT push)
3. Consumer tracks its position using **offsets**
4. After processing, consumer **commits** the offset ("I've read up to message #5")
5. If consumer crashes and restarts, it resumes from the last committed offset

```
Consumer Reading from Partition 0:

  [0] [1] [2] [3] [4] [5] [6] [7] [8]
                    ▲              ▲
                    │              │
            Last committed     Latest message
              offset (3)        (offset 8)

  Consumer will read messages 4, 5, 6, 7, 8 next.
```

---

### 4.6 Consumer Group

**Simple Definition:**  
A consumer group is a **set of consumers** that cooperate to consume messages from a topic. Each partition is assigned to **exactly one consumer** within the group.

**Real-World Analogy:**  
Think of a **pizza delivery team**. A neighborhood (topic) is divided into zones (partitions). Each delivery driver (consumer) handles specific zones. No two drivers deliver to the same zone. If a driver calls in sick, their zone is reassigned to another driver (rebalancing).

**Internal Working:**
- Each consumer group has a unique **group ID**
- Kafka ensures each partition is consumed by **only one consumer** in a group
- If consumers > partitions → some consumers sit idle
- If consumers < partitions → some consumers handle multiple partitions
- If a consumer joins/leaves → **rebalancing** happens automatically

**Consumer Group Rules:**

| Partitions | Consumers in Group | Result |
|------------|-------------------|--------|
| 3 | 1 | 1 consumer reads all 3 partitions |
| 3 | 2 | 1 consumer reads 2 partitions, 1 reads 1 |
| 3 | 3 | Each consumer reads exactly 1 partition (ideal) |
| 3 | 4 | 3 consumers read 1 each, **1 consumer sits IDLE** |

```
Topic "orders" (3 partitions) — Consumer Group "order-service"
══════════════════════════════════════════════════════════════

           ┌───────────────────────────────────┐
           │        Topic: "orders"            │
           │                                   │
           │  Partition 0   P1       P2        │
           └──────┬────────┬────────┬──────────┘
                  │        │        │
                  ▼        ▼        ▼
           ┌──────────┐ ┌────┐ ┌──────────┐
           │Consumer A│ │ C-B│ │Consumer C│
           └──────────┘ └────┘ └──────────┘
           └──────── Consumer Group: "order-service" ────────┘

Multiple Consumer Groups can read the SAME topic independently:

           ┌───────────────────────────────────┐
           │        Topic: "orders"            │
           └──────┬────────┬────────┬──────────┘
                  │        │        │
          ┌───────┴────────┴────────┴────────┐
          │  Group "order-service"            │
          │  Consumer A ← P0, P1             │
          │  Consumer B ← P2                 │
          └──────────────────────────────────┘
          ┌───────┴────────┴────────┴────────┐
          │  Group "analytics-service"       │
          │  Consumer X ← P0, P1, P2         │
          └──────────────────────────────────┘

  Both groups read ALL messages, independently!
```

---

### 4.7 Offset

**Simple Definition:**  
An offset is a **unique sequential ID** assigned to each message within a partition. It represents the position of a message.

**Real-World Analogy:**  
An offset is like a **page number** in a book. If you stop reading at page 42, you bookmark it (commit offset). Next time, you start from page 43. You can also go back and re-read page 10 (replay).

**Internal Working:**
- Offsets start at **0** and increment by 1 for each new message
- Offsets are **per partition** (Partition 0 and Partition 1 each have their own offset sequence)
- Consumers track 3 types of offsets:

| Offset Type | Meaning |
|-------------|---------|
| **Current offset** | The offset of the next message to be fetched |
| **Committed offset** | The offset the consumer has confirmed processing |
| **Log-end offset** | The offset of the latest message in the partition |

- Committed offsets are stored in a special internal topic: `__consumer_offsets`
- Offsets are **never reused** (even if old messages are deleted)

```
Partition 0 Timeline:

  Offset:   0    1    2    3    4    5    6    7    8    9
          ┌────┬────┬────┬────┬────┬────┬────┬────┬────┬────┐
          │ m0 │ m1 │ m2 │ m3 │ m4 │ m5 │ m6 │ m7 │ m8 │ m9 │
          └────┴────┴────┴────┴────┴────┴────┴────┴────┴────┘
                          ▲                        ▲         ▲
                          │                        │         │
                   Committed Offset          Current    Log-End
                   (processed till here)     Position   Offset
```

---

### 4.8 Replication

**Simple Definition:**  
Replication means keeping **copies of each partition on multiple brokers** so that data is not lost if a broker crashes.

**Real-World Analogy:**  
Think of **backup copies** of important documents. You keep the original at home, a copy in a bank locker, and another with a trusted friend. If your house burns down, you still have the document.

**Internal Working:**
- Configured per topic via `replication-factor` (typically 3)
- One copy is the **Leader** — handles all reads/writes
- Other copies are **Followers** — continuously replicate from the leader
- **ISR (In-Sync Replicas):** The set of replicas that are up to date with the leader
- If the leader fails, one of the ISR followers becomes the new leader

```
Topic "orders", Partition 0, Replication Factor = 3

  ┌─────────────────────────────────────────────────────────────┐
  │                                                             │
  │   Broker 0                Broker 1              Broker 2    │
  │   ┌───────────────┐      ┌──────────────┐      ┌─────────────┐
  │   │ P0 (LEADER)   │ ───> │ P0 (FOLLOWER)│      │ P0 (FOLLOWER)│
  │   │               │      │              │      │              │
  │   │ [0][1][2][3]  │      │ [0][1][2][3] │      │ [0][1][2][3] │
  │   └───────────────┘      └──────────────┘      └─────────────┘
  │         ▲                       │                     │       │
  │         │                       │                     │       │
  │    All reads &           Replicates from        Replicates    │
  │    writes go here        the Leader             from Leader   │
  │                                                               │
  │   ISR = {Broker 0, Broker 1, Broker 2}                       │
  └─────────────────────────────────────────────────────────────┘

  If Broker 0 crashes → Broker 1 or 2 becomes the new LEADER
```

---

### 4.9 Leader & Follower

**Simple Definition:**
- **Leader:** The primary replica that handles all read and write requests for a partition
- **Follower:** A backup replica that copies data from the leader and can take over if the leader fails

**Real-World Analogy:**  
Think of a **team lead and backup team leads**. The team lead (leader) handles all client meetings and decisions. Backup leads (followers) sit in on every meeting and take notes. If the team lead gets sick, a backup lead can immediately step in with full context.

**Internal Working:**
- Every partition has **exactly one leader** and zero or more followers
- Producers and consumers **only talk to the leader**
- Followers continuously fetch data from the leader (like a consumer)
- If a follower falls too far behind, it's removed from the ISR
- Leader election happens automatically via the **Controller** broker

| Aspect | Leader | Follower |
|--------|--------|----------|
| **Handles reads** | Yes | No (by default) |
| **Handles writes** | Yes | No |
| **Count per partition** | Exactly 1 | 0 to (replication-factor - 1) |
| **Purpose** | Serve client requests | Provide fault tolerance |
| **Data flow** | Receives from producer | Replicates from leader |
| **Failure behavior** | Triggers leader election | Can be promoted to leader |
| **In ISR** | Always | Only if caught up |

---

### Concept Comparison Tables

**Partition vs Topic:**

| Aspect | Topic | Partition |
|--------|-------|-----------|
| **What** | Logical category/name | Physical subdivision of a topic |
| **Analogy** | A book | A chapter in the book |
| **Ordering** | No ordering guarantee | Strict ordering within partition |
| **Count** | Many topics per cluster | Many partitions per topic |
| **Stored on** | Distributed across brokers | A single broker (per replica) |
| **Consumer mapping** | A group subscribes to a topic | A partition maps to one consumer in a group |

**Producer vs Consumer:**

| Aspect | Producer | Consumer |
|--------|----------|----------|
| **Role** | Writes/sends messages | Reads/receives messages |
| **Direction** | Data INTO Kafka | Data OUT OF Kafka |
| **Talks to** | Leader of target partition | Leader of assigned partition |
| **Key concern** | Acknowledgments (acks) | Offset management |
| **Scaling** | Add more producer instances | Add more consumers to group |
| **Delivery mode** | Fire-and-forget / sync / async | Pull-based (consumer controls pace) |

---

## 5. Message Flow: End to End

### 5.1 Producer → Kafka → Consumer Flow

```
STEP-BY-STEP MESSAGE FLOW
══════════════════════════

Step 1: Producer sends message
──────────────────────────────
  Producer App
      │
      │  send("order-events", key="user-42", value="Order #123")
      │
      ▼
  ┌──────────────────────────────────┐
  │        Serializer                │
  │   Key → bytes, Value → bytes    │
  └──────────────┬───────────────────┘
                 │
                 ▼
  ┌──────────────────────────────────┐
  │        Partitioner               │
  │   hash("user-42") % 3 = 1       │
  │   → Partition 1                  │
  └──────────────┬───────────────────┘
                 │
                 ▼
  ┌──────────────────────────────────┐
  │     Record Accumulator (Buffer)  │
  │   Batches messages per partition  │
  │   Sends when batch full or       │
  │   linger.ms expires              │
  └──────────────┬───────────────────┘
                 │
                 ▼
Step 2: Network send to Leader Broker
──────────────────────────────────────
  ┌──────────────────────────────────┐
  │   Broker 1 (Leader of P1)       │
  │                                  │
  │   1. Writes message to disk log  │
  │   2. Assigns offset (e.g., 47)   │
  │   3. Replicates to followers     │
  │   4. Returns acknowledgment      │
  └──────────────┬───────────────────┘
                 │
                 ▼
Step 3: Follower Replication
────────────────────────────
  Broker 1 (Leader)  ──replicate──>  Broker 0 (Follower)
                     ──replicate──>  Broker 2 (Follower)

Step 4: Acknowledgment to Producer
───────────────────────────────────
  acks=0  → No waiting (fire and forget)
  acks=1  → Leader wrote to disk
  acks=all → Leader + all ISR followers wrote

Step 5: Consumer reads the message
───────────────────────────────────
  Consumer (group="order-service")
      │
      │  poll()  →  Fetch from Partition 1, offset ≥ last committed
      │
      ▼
  ┌──────────────────────────────────┐
  │   Broker 1 (Leader of P1)       │
  │   Returns messages from offset   │
  │   48 onward (batch of messages)  │
  └──────────────┬───────────────────┘
                 │
                 ▼
  ┌──────────────────────────────────┐
  │        Deserializer              │
  │   bytes → Key, Value objects     │
  └──────────────┬───────────────────┘
                 │
                 ▼
  Consumer processes message → commits offset
```

### 5.2 Consumer Group Rebalancing Flow

Rebalancing happens when a consumer **joins**, **leaves**, or **crashes** within a group.

```
CONSUMER GROUP REBALANCING
══════════════════════════

BEFORE: 2 consumers, 3 partitions
─────────────────────────────────
  Consumer A ← Partition 0, Partition 1
  Consumer B ← Partition 2

EVENT: Consumer C joins the group
─────────────────────────────────
  1. Consumer C sends JoinGroup request to Group Coordinator (a broker)
  2. Group Coordinator triggers rebalance
  3. ALL consumers in the group:
     a. Stop processing
     b. Commit their current offsets
     c. Revoke current partition assignments
  4. Group Coordinator reassigns partitions:

AFTER: 3 consumers, 3 partitions
─────────────────────────────────
  Consumer A ← Partition 0
  Consumer B ← Partition 1
  Consumer C ← Partition 2

  ✓ Perfectly balanced!

FAILURE SCENARIO: Consumer B crashes
─────────────────────────────────────
  1. Consumer B stops sending heartbeats
  2. After session.timeout.ms, Coordinator declares B dead
  3. Rebalance triggered:

  Consumer A ← Partition 0, Partition 1
  Consumer C ← Partition 2

  Consumer A picks up where B left off (from B's last committed offset)
```

### 5.3 Replication Flow

```
REPLICATION FLOW (Detail)
═════════════════════════

                    ┌──────────────────────────┐
                    │        Producer           │
                    └────────────┬─────────────┘
                                │
                       write(msg)
                                │
                                ▼
                    ┌──────────────────────────┐
                    │  Leader (Broker 0)        │
                    │                          │
                    │  Log: [0][1][2][3][4]    │
                    │            HW=3 LEO=5    │
                    └─────┬──────────┬─────────┘
                          │          │
                    fetch │          │ fetch
                          │          │
                          ▼          ▼
            ┌──────────────┐  ┌──────────────┐
            │ Follower      │  │ Follower      │
            │ (Broker 1)    │  │ (Broker 2)    │
            │               │  │               │
            │ [0][1][2][3]  │  │ [0][1][2]     │
            │        LEO=4  │  │       LEO=3   │
            └──────────────┘  └──────────────┘

  HW  = High Watermark (offset up to which consumers can read)
  LEO = Log End Offset (latest offset on that broker)

  HW = min(LEO of all ISR replicas) = min(5, 4, 3) = 3
  Consumers can only read up to offset 3 (committed data)

  After Broker 2 catches up:
    LEO(Broker 0)=5, LEO(Broker 1)=5, LEO(Broker 2)=5
    HW moves to 5 → consumers can now read all messages
```

---

## 6. Kafka Architecture Diagrams

### 6.1 Complete Kafka Architecture

```
┌───────────────────────────────────────────────────────────────────────────────────┐
│                            KAFKA ECOSYSTEM                                        │
│                                                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                              │
│  │ Producer 1  │  │ Producer 2  │  │ Producer 3  │     ← Applications that       │
│  │ (Order Svc) │  │ (User Svc)  │  │ (Payment)   │       write events            │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                              │
│         │                │                │                                       │
│         └────────────────┼────────────────┘                                       │
│                          │                                                        │
│                          ▼                                                        │
│  ┌────────────────────────────────────────────────────────────────────────────┐   │
│  │                        KAFKA CLUSTER                                       │   │
│  │                                                                            │   │
│  │   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐                 │   │
│  │   │  Broker 0    │   │  Broker 1    │   │  Broker 2    │                 │   │
│  │   │              │   │              │   │              │                 │   │
│  │   │  orders-P0(L)│   │  orders-P1(L)│   │  orders-P2(L)│                 │   │
│  │   │  orders-P1(F)│   │  orders-P2(F)│   │  orders-P0(F)│                 │   │
│  │   │  users-P0(L) │   │  users-P1(L) │   │  users-P0(F) │                 │   │
│  │   │              │   │              │   │  users-P1(F) │                 │   │
│  │   └──────────────┘   └──────────────┘   └──────────────┘                 │   │
│  │                                                                            │   │
│  │   ┌────────────────────────────────────────────────────────────────┐       │   │
│  │   │  ZooKeeper / KRaft (Cluster Metadata Management)              │       │   │
│  │   │  • Broker registration         • Leader election              │       │   │
│  │   │  • Topic configuration          • Consumer group coordination │       │   │
│  │   └────────────────────────────────────────────────────────────────┘       │   │
│  └────────────────────────────────────────────────────────────────────────────┘   │
│                          │                                                        │
│         ┌────────────────┼────────────────┐                                       │
│         │                │                │                                       │
│         ▼                ▼                ▼                                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                              │
│  │ Consumer 1  │  │ Consumer 2  │  │ Consumer 3  │     ← Applications that       │
│  │ (Analytics) │  │ (Search)    │  │ (Notif Svc) │       read events             │
│  └─────────────┘  └─────────────┘  └─────────────┘                              │
│                                                                                   │
└───────────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Partition Distribution Across Brokers

```
PARTITION DISTRIBUTION EXAMPLE
══════════════════════════════

Topic: "orders" — 6 partitions, replication factor 3
Topic: "users"  — 3 partitions, replication factor 2

Broker 0                    Broker 1                    Broker 2
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│                     │    │                     │    │                     │
│ orders-P0 [LEADER]  │    │ orders-P0 [FOLLOWER]│    │ orders-P0 [FOLLOWER]│
│ orders-P1 [FOLLOWER]│    │ orders-P1 [LEADER]  │    │ orders-P1 [FOLLOWER]│
│ orders-P2 [FOLLOWER]│    │ orders-P2 [FOLLOWER]│    │ orders-P2 [LEADER]  │
│ orders-P3 [LEADER]  │    │ orders-P3 [FOLLOWER]│    │ orders-P3 [FOLLOWER]│
│ orders-P4 [FOLLOWER]│    │ orders-P4 [LEADER]  │    │ orders-P4 [FOLLOWER]│
│ orders-P5 [FOLLOWER]│    │ orders-P5 [FOLLOWER]│    │ orders-P5 [LEADER]  │
│                     │    │                     │    │                     │
│ users-P0 [LEADER]   │    │ users-P0 [FOLLOWER] │    │                     │
│ users-P1 [FOLLOWER] │    │ users-P1 [LEADER]   │    │ users-P1 [FOLLOWER] │ ← not present
│                     │    │                     │    │ users-P2 [LEADER]   │
│ users-P2 [FOLLOWER] │    │                     │    │                     │
│                     │    │                     │    │                     │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘

Key Points:
  • Leaders are DISTRIBUTED evenly across brokers (load balancing)
  • Each partition's replicas are on DIFFERENT brokers (fault tolerance)
  • If Broker 1 crashes → Broker 0 or 2 take over its leader partitions
```

### 6.3 Consumer Group Working (Detailed)

```
CONSUMER GROUP SCENARIOS
════════════════════════

Scenario 1: Consumers < Partitions (underutilized)
───────────────────────────────────────────────────
  Topic "events" has 4 partitions, Consumer Group has 2 consumers

  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐
  │  P0  │  │  P1  │  │  P2  │  │  P3  │
  └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘
     │         │         │         │
     └────┬────┘         └────┬────┘
          │                   │
          ▼                   ▼
     ┌─────────┐         ┌─────────┐
     │  C-1    │         │  C-2    │
     │ (2 part)│         │ (2 part)│
     └─────────┘         └─────────┘

Scenario 2: Consumers = Partitions (ideal)
──────────────────────────────────────────
  Topic "events" has 4 partitions, Consumer Group has 4 consumers

  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐
  │  P0  │  │  P1  │  │  P2  │  │  P3  │
  └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘
     │         │         │         │
     ▼         ▼         ▼         ▼
  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐
  │ C-1 │  │ C-2 │  │ C-3 │  │ C-4 │
  └─────┘  └─────┘  └─────┘  └─────┘

     ✓ Maximum parallelism achieved!

Scenario 3: Consumers > Partitions (waste)
──────────────────────────────────────────
  Topic "events" has 3 partitions, Consumer Group has 5 consumers

  ┌──────┐  ┌──────┐  ┌──────┐
  │  P0  │  │  P1  │  │  P2  │
  └──┬───┘  └──┬───┘  └──┬───┘
     │         │         │
     ▼         ▼         ▼
  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐
  │ C-1 │  │ C-2 │  │ C-3 │  │ C-4 │  │ C-5 │
  └─────┘  └─────┘  └─────┘  └─────┘  └─────┘
                               IDLE      IDLE

     ⚠ C-4 and C-5 do nothing. Wasted resources!

RULE: Max useful consumers = Number of partitions
```

### 6.4 Replication Leader/Follower Diagram

```
LEADER ELECTION AFTER FAILURE
═════════════════════════════

BEFORE FAILURE:
───────────────
  Broker 0 (LEADER)     Broker 1 (FOLLOWER)    Broker 2 (FOLLOWER)
  ┌────────────────┐    ┌────────────────┐     ┌────────────────┐
  │ P0: [0..99]    │───>│ P0: [0..99]    │     │ P0: [0..97]    │
  │ (ISR member)   │    │ (ISR member)   │     │ (ISR member*)  │
  └────────────────┘    └────────────────┘     └────────────────┘
        ▲                                       * slightly behind
    All reads                                     but within
    & writes                                      replica.lag threshold

BROKER 0 CRASHES!
─────────────────
  Broker 0 (DOWN)       Broker 1 (FOLLOWER)    Broker 2 (FOLLOWER)
  ┌────────────────┐    ┌────────────────┐     ┌────────────────┐
  │    ╳ DEAD ╳    │    │ P0: [0..99]    │     │ P0: [0..97]    │
  │                │    │ (ISR member)   │     │ (ISR member)   │
  └────────────────┘    └────────────────┘     └────────────────┘

  Controller detects Broker 0 failure
  → Selects Broker 1 as new leader (it was fully in sync)

AFTER ELECTION:
───────────────
  Broker 0 (DOWN)       Broker 1 (NEW LEADER)  Broker 2 (FOLLOWER)
  ┌────────────────┐    ┌────────────────┐     ┌────────────────┐
  │    ╳ DEAD ╳    │    │ P0: [0..99]    │     │ P0: [0..99]    │
  │                │    │ ★ NOW LEADER   │     │ (catches up)   │
  └────────────────┘    └────────────────┘     └────────────────┘
                              ▲
                          All reads
                          & writes now

  ISR = {Broker 1, Broker 2}    (Broker 0 removed)
  Zero data loss. Zero downtime (for consumers).
```

---

## 7. Deep Understanding: Why Things Work This Way

### 7.1 Why Partitions Are Important (Parallelism)

**The Core Problem:** One consumer can only read so fast. If your topic gets 100,000 messages/second, one consumer might only handle 10,000/sec.

**The Solution:** Split the topic into partitions. Now 10 consumers can each handle 10,000/sec from their own partition.

```
Without Partitions (Bottleneck):
  100,000 msg/sec → [Single Queue] → [Single Consumer] → 10,000 msg/sec
                                                          ⚠ 90,000 msg/sec BACKLOG!

With 10 Partitions:
  100,000 msg/sec → [P0] → [Consumer 0] → 10,000 msg/sec  ✓
                    [P1] → [Consumer 1] → 10,000 msg/sec  ✓
                    [P2] → [Consumer 2] → 10,000 msg/sec  ✓
                    ...
                    [P9] → [Consumer 9] → 10,000 msg/sec  ✓
                                                           Total: 100,000 msg/sec ✓
```

**Choosing the Right Number of Partitions:**

| Factor | More Partitions | Fewer Partitions |
|--------|----------------|-----------------|
| **Throughput** | Higher (more parallelism) | Lower |
| **Consumer count** | More consumers possible | Fewer needed |
| **Memory on broker** | More (each partition uses memory) | Less |
| **Rebalancing time** | Slower (more to reassign) | Faster |
| **End-to-end latency** | Slightly higher | Lower |
| **File handles** | More open files | Fewer |

**Rule of Thumb:**
- Start with **number of partitions = expected peak throughput / throughput per consumer**
- Or simply: **number of partitions >= number of consumers you plan to run**
- Common default: **3 to 12** for most use cases
- You can **increase** partitions later but **cannot decrease** them

---

### 7.2 Why Offsets Matter

Offsets are the reason Kafka consumers are so **resilient** and **flexible**.

**Without offsets (traditional MQ):**
```
Consumer reads message → MQ deletes message → message is GONE
If consumer crashes mid-processing → message is LOST
No way to re-read old messages
```

**With offsets (Kafka):**
```
Consumer reads message at offset 5 → message STAYS in Kafka
Consumer processes it → commits offset 5 → "I'm done with 5"
Consumer crashes → restarts → asks "what was my last committed offset?" → 5
Consumer resumes from offset 6 → NO MESSAGE LOST
Consumer wants to re-process? → reset offset to 0 → replay everything!
```

**Offset Commit Strategies:**

| Strategy | How | Pros | Cons |
|----------|-----|------|------|
| **Auto commit** | `enable.auto.commit=true` | Simple, no code needed | May lose messages (commits before processing) |
| **Manual sync** | `consumer.commitSync()` | Guaranteed commit before proceeding | Blocks the consumer thread |
| **Manual async** | `consumer.commitAsync()` | Non-blocking | May fail silently; harder to reason about |
| **Manual per-record** | Commit after each record | Precise control | Lower throughput due to frequent commits |

---

### 7.3 How Kafka Ensures Durability and Fault Tolerance

Kafka uses **multiple layers** of protection:

```
DURABILITY & FAULT TOLERANCE LAYERS
════════════════════════════════════

Layer 1: DISK PERSISTENCE
─────────────────────────
  Every message is written to DISK immediately.
  Kafka uses sequential I/O (very fast on modern disks).
  Messages are retained for a configurable period (default: 7 days).

Layer 2: REPLICATION
────────────────────
  Each partition is replicated across multiple brokers.
  Default: replication-factor=3 (1 leader + 2 followers).
  Even if 2 brokers die, data survives.

Layer 3: ISR (In-Sync Replicas)
───────────────────────────────
  Only followers that are "caught up" are in the ISR.
  Producer with acks=all waits for ALL ISR to acknowledge.
  Guarantees: if leader dies, new leader has ALL data.

Layer 4: ACKNOWLEDGMENTS (acks)
───────────────────────────────
  acks=0:   No guarantee (fastest)
  acks=1:   Leader wrote to disk (balanced)
  acks=all: Leader + all ISR wrote (strongest)

Layer 5: MIN.INSYNC.REPLICAS
─────────────────────────────
  min.insync.replicas=2 means:
  → At least 2 replicas (including leader) must acknowledge
  → If only 1 replica is alive, producer gets an error
  → Prevents writing to an under-replicated partition
```

**Durability Configuration Matrix:**

| Setting | `acks=0` | `acks=1` | `acks=all` |
|---------|----------|----------|------------|
| **Speed** | Fastest | Fast | Slowest |
| **Data Safety** | May lose data | May lose if leader crashes | No data loss (with ISR) |
| **Use Case** | Metrics, logs | General purpose | Financial transactions |
| **Throughput** | Highest | High | Lower |

---

### 7.4 Delivery Semantics

## 8. Delivery Semantics

This is one of the **most important** concepts for interviews and real-world usage.

| Semantic | Meaning | How It Happens | Use Case |
|----------|---------|----------------|----------|
| **At-most-once** | Message may be lost, never duplicated | Consumer commits offset BEFORE processing | Metrics, logs where losing some data is OK |
| **At-least-once** | Message never lost, may be duplicated | Consumer commits offset AFTER processing | Most common; requires idempotent consumers |
| **Exactly-once** | Message delivered exactly once | Kafka transactions + idempotent producer | Financial systems, billing, critical data |

### At-Most-Once (May Lose Messages)

```
1. Consumer fetches message at offset 5
2. Consumer COMMITS offset 5  ← (committed before processing!)
3. Consumer starts processing
4. Consumer CRASHES mid-processing
5. Consumer restarts → resumes from offset 6
   ⚠ Message at offset 5 was NEVER fully processed = LOST
```

### At-Least-Once (May Duplicate Messages)

```
1. Consumer fetches message at offset 5
2. Consumer PROCESSES message successfully
3. Consumer tries to COMMIT offset 5
4. Consumer CRASHES before commit completes
5. Consumer restarts → resumes from offset 5 (last committed was 4)
   ⚠ Message at offset 5 is PROCESSED AGAIN = DUPLICATE

   Solution: Make consumers IDEMPOTENT
   (processing the same message twice has the same effect as once)
   Example: Use a unique order-id as database primary key
            → INSERT fails on duplicate = no harm done
```

### Exactly-Once (Kafka Transactions)

```
Producer Side:
  enable.idempotence=true
  transactional.id="my-transactional-producer"

  producer.beginTransaction();
  producer.send(record1);
  producer.send(record2);
  producer.commitTransaction();   // atomic: all or nothing

Consumer Side:
  isolation.level=read_committed
  → Consumer only sees messages from COMMITTED transactions
  → Never sees partial writes

Result: End-to-end exactly-once semantics
```

**Decision Guide:**

```
Do you need exactly-once? ────── YES ───> Use Kafka Transactions
        │                                 (higher complexity, lower throughput)
        NO
        │
        ▼
Can you tolerate duplicates? ─── YES ───> At-least-once (recommended default)
        │                                 + Idempotent consumers
        NO
        │
        ▼
Can you tolerate data loss? ──── YES ───> At-most-once
        │                                 (simplest, fastest)
        NO
        │
        ▼
Use Kafka Transactions (exactly-once)
```

---

## 9. Practical Section: Running Kafka Locally

### 9.1 Option A: Using Docker Compose (Recommended)

**Step 1: Create `docker-compose.yml`**

```yaml
version: '3.8'
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - "2181:2181"

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

**Step 2: Start Kafka**

```bash
docker-compose up -d
```

**Step 3: Verify it's running**

```bash
docker-compose ps
```

### 9.2 Option B: Using KRaft (No Zookeeper, Kafka 3.3+)

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
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@localhost:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk
```

### 9.3 Basic Operations (Command Line)

**Create a Topic:**

```bash
# Create a topic named "my-first-topic" with 3 partitions
docker exec -it <kafka-container> kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic my-first-topic \
  --partitions 3 \
  --replication-factor 1
```

**List Topics:**

```bash
docker exec -it <kafka-container> kafka-topics \
  --bootstrap-server localhost:9092 \
  --list
```

**Describe a Topic:**

```bash
docker exec -it <kafka-container> kafka-topics \
  --bootstrap-server localhost:9092 \
  --describe \
  --topic my-first-topic
```

**Produce Messages (Terminal 1):**

```bash
docker exec -it <kafka-container> kafka-console-producer \
  --bootstrap-server localhost:9092 \
  --topic my-first-topic

> Hello Kafka!
> This is my first message
> Order #123 placed by user-42
```

**Consume Messages (Terminal 2):**

```bash
# Read from beginning
docker exec -it <kafka-container> kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic my-first-topic \
  --from-beginning

# Read only new messages (no --from-beginning)
docker exec -it <kafka-container> kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic my-first-topic \
  --group my-consumer-group
```

**Check Consumer Group Offsets:**

```bash
docker exec -it <kafka-container> kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group my-consumer-group
```

---

## 10. Spring Boot Integration

### 10.1 Project Setup

**Step 1: Add Dependencies (`pom.xml`)**

```xml
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka</artifactId>
</dependency>
```

Spring Boot auto-configures Kafka. No extra version needed if using Spring Boot's dependency management.

**Step 2: Configure (`application.yml`)**

```yaml
spring:
  kafka:
    bootstrap-servers: localhost:9092

    producer:
      key-serializer: org.apache.kafka.common.serialization.StringSerializer
      value-serializer: org.apache.kafka.common.serialization.StringSerializer
      acks: all
      retries: 3

    consumer:
      group-id: my-application-group
      key-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      value-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      auto-offset-reset: earliest
      enable-auto-commit: false

    listener:
      ack-mode: manual_immediate
```

### 10.2 Producer Example

```java
@Service
@RequiredArgsConstructor
public class OrderEventProducer {

    private final KafkaTemplate<String, String> kafkaTemplate;

    private static final String TOPIC = "order-events";

    public void sendOrderEvent(String orderId, String eventPayload) {
        kafkaTemplate.send(TOPIC, orderId, eventPayload)
            .whenComplete((result, ex) -> {
                if (ex == null) {
                    log.info("Sent message=[{}] with offset=[{}]",
                        eventPayload,
                        result.getRecordMetadata().offset());
                } else {
                    log.error("Failed to send message=[{}]", eventPayload, ex);
                }
            });
    }
}
```

### 10.3 Consumer Example

```java
@Service
@Slf4j
public class OrderEventConsumer {

    @KafkaListener(
        topics = "order-events",
        groupId = "order-processing-group",
        concurrency = "3"
    )
    public void consume(
            @Payload String message,
            @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
            @Header(KafkaHeaders.OFFSET) long offset,
            Acknowledgment acknowledgment) {

        log.info("Received message=[{}] from partition=[{}] offset=[{}]",
            message, partition, offset);

        try {
            processOrder(message);
            acknowledgment.acknowledge();
        } catch (Exception e) {
            log.error("Error processing message: {}", message, e);
            // Don't acknowledge → message will be re-delivered
        }
    }

    private void processOrder(String message) {
        // Business logic here
    }
}
```

### 10.4 Sending JSON Objects

**Step 1: Configure JSON Serializer**

```yaml
spring:
  kafka:
    producer:
      value-serializer: org.springframework.kafka.support.serializer.JsonSerializer
    consumer:
      value-deserializer: org.springframework.kafka.support.serializer.JsonDeserializer
      properties:
        spring.json.trusted.packages: "com.example.dto"
```

**Step 2: Create a DTO**

```java
@Data
@AllArgsConstructor
@NoArgsConstructor
public class OrderEvent {
    private String orderId;
    private String userId;
    private String product;
    private double amount;
    private String status;
    private LocalDateTime timestamp;
}
```

**Step 3: Producer with JSON**

```java
@Service
@RequiredArgsConstructor
public class OrderEventProducer {

    private final KafkaTemplate<String, OrderEvent> kafkaTemplate;

    public void sendOrderEvent(OrderEvent event) {
        kafkaTemplate.send("order-events", event.getOrderId(), event);
    }
}
```

**Step 4: Consumer with JSON**

```java
@Service
@Slf4j
public class OrderEventConsumer {

    @KafkaListener(topics = "order-events", groupId = "order-processing-group")
    public void consume(OrderEvent event, Acknowledgment ack) {
        log.info("Received order: {} for user: {}", event.getOrderId(), event.getUserId());
        processOrder(event);
        ack.acknowledge();
    }
}
```

### 10.5 Creating Topics Programmatically

```java
@Configuration
public class KafkaTopicConfig {

    @Bean
    public NewTopic orderEventsTopic() {
        return TopicBuilder.name("order-events")
            .partitions(3)
            .replicas(1)
            .config(TopicConfig.RETENTION_MS_CONFIG, "604800000") // 7 days
            .build();
    }

    @Bean
    public NewTopic userEventsTopic() {
        return TopicBuilder.name("user-events")
            .partitions(3)
            .replicas(1)
            .build();
    }
}
```

### 10.6 REST Controller to Test

```java
@RestController
@RequestMapping("/api/orders")
@RequiredArgsConstructor
public class OrderController {

    private final OrderEventProducer producer;

    @PostMapping
    public ResponseEntity<String> createOrder(@RequestBody OrderEvent event) {
        event.setTimestamp(LocalDateTime.now());
        event.setStatus("CREATED");
        producer.sendOrderEvent(event);
        return ResponseEntity.ok("Order event published: " + event.getOrderId());
    }
}
```

### 10.7 Error Handling and Retry

```java
@Configuration
public class KafkaConsumerConfig {

    @Bean
    public DefaultErrorHandler errorHandler() {
        // Retry 3 times with 1 second backoff, then send to DLT
        BackOff backOff = new FixedBackOff(1000L, 3);
        return new DefaultErrorHandler(
            new DeadLetterPublishingRecoverer(kafkaTemplate()),
            backOff
        );
    }
}
```

```
ERROR HANDLING FLOW
═══════════════════

  Message arrives
       │
       ▼
  Consumer processes ── Success ──> Acknowledge ✓
       │
     Failure
       │
       ▼
  Retry 1 (after 1s) ── Success ──> Acknowledge ✓
       │
     Failure
       │
       ▼
  Retry 2 (after 1s) ── Success ──> Acknowledge ✓
       │
     Failure
       │
       ▼
  Retry 3 (after 1s) ── Success ──> Acknowledge ✓
       │
     Failure
       │
       ▼
  Send to Dead Letter Topic (DLT): "order-events.DLT"
  → Investigate and fix later
```

---

## 11. Important Kafka Configurations

### 11.1 Producer Configurations

| Property | Default | Recommended | Description |
|----------|---------|-------------|-------------|
| `acks` | `1` | `all` | How many replicas must acknowledge. `all` = strongest durability |
| `retries` | `2147483647` | `3` | Number of retries on failure |
| `batch.size` | `16384` (16KB) | `32768` (32KB) | Batch size for network send |
| `linger.ms` | `0` | `5-20` | Time to wait before sending a batch (allows more batching) |
| `buffer.memory` | `33554432` (32MB) | Keep default | Total memory for buffering unsent messages |
| `max.in.flight.requests.per.connection` | `5` | `1` (if ordering critical) | Max unacknowledged requests |
| `enable.idempotence` | `true` (Kafka 3.0+) | `true` | Prevent duplicate messages |
| `compression.type` | `none` | `snappy` or `lz4` | Compress messages (saves network & disk) |

### 11.2 Consumer Configurations

| Property | Default | Recommended | Description |
|----------|---------|-------------|-------------|
| `group.id` | — | Always set | Consumer group name |
| `auto.offset.reset` | `latest` | `earliest` (initial) | Where to start if no committed offset |
| `enable.auto.commit` | `true` | `false` | Auto-commit offsets? (prefer manual) |
| `auto.commit.interval.ms` | `5000` | — | Auto-commit interval (if auto-commit on) |
| `max.poll.records` | `500` | `100-500` | Max records per poll() call |
| `max.poll.interval.ms` | `300000` (5min) | Adjust to processing time | Max time between poll() calls |
| `session.timeout.ms` | `45000` | `10000-30000` | Time before consumer is considered dead |
| `heartbeat.interval.ms` | `3000` | `session.timeout / 3` | Heartbeat frequency |
| `fetch.min.bytes` | `1` | `1-1048576` | Min data per fetch (higher = more batching) |

### 11.3 Topic Configurations

| Property | Default | Description |
|----------|---------|-------------|
| `num.partitions` | `1` | Default partitions for auto-created topics |
| `replication.factor` | `1` | Default replication for auto-created topics |
| `retention.ms` | `604800000` (7 days) | How long to keep messages |
| `retention.bytes` | `-1` (unlimited) | Max size per partition before deletion |
| `cleanup.policy` | `delete` | `delete` = remove old messages, `compact` = keep latest per key |
| `min.insync.replicas` | `1` | Min replicas that must acknowledge (with acks=all) |
| `segment.bytes` | `1073741824` (1GB) | Log segment file size |

### 11.4 The "Safe Producer" Configuration

For production systems that cannot afford data loss:

```yaml
spring:
  kafka:
    producer:
      acks: all
      retries: 3
      properties:
        enable.idempotence: true
        max.in.flight.requests.per.connection: 5
        min.insync.replicas: 2
```

```
SAFE PRODUCER EXPLAINED
═══════════════════════

  acks=all
    → Leader waits for ALL ISR replicas to write before acknowledging

  enable.idempotence=true
    → Kafka deduplicates messages (even if producer retries)
    → Uses producer-id + sequence-number internally

  min.insync.replicas=2
    → At least 2 replicas (including leader) must be alive
    → If only 1 broker is up → producer gets NotEnoughReplicas error
    → Prevents writing data that exists on only 1 broker

  retries=3
    → Retries transient failures (network blip, leader election)

  Result: Messages are NEVER lost (if cluster has 2+ live brokers)
```

---

## 12. Common Interview Questions

### Q1: What is Kafka and why is it used?

**Answer:**  
Kafka is a distributed event streaming platform used for building real-time data pipelines and streaming applications. It acts as a high-throughput, fault-tolerant message broker where producers publish events and consumers read them asynchronously. Unlike traditional MQs, Kafka retains messages on disk, supports replay, and scales horizontally.

### Q2: How does Kafka guarantee message ordering?

**Answer:**  
Kafka guarantees ordering **within a single partition only**. All messages with the same key go to the same partition (via hashing). Within that partition, messages are strictly ordered by offset. There is no ordering guarantee across partitions.

### Q3: What happens when a broker goes down?

**Answer:**  
If a broker holding leader partitions goes down, the Controller broker detects this (via ZooKeeper or KRaft) and promotes a follower from the ISR to be the new leader. Consumers and producers are redirected to the new leader with minimal disruption. Once the broker recovers, it rejoins as a follower and catches up.

### Q4: Explain Consumer Groups.

**Answer:**  
A consumer group is a set of consumers sharing a group ID. Kafka distributes partitions among consumers in the group so that each partition is consumed by exactly one consumer. This enables parallel processing. If a consumer dies, Kafka rebalances and assigns its partitions to remaining consumers. Different consumer groups read the same topic independently.

### Q5: What is the difference between acks=0, acks=1, and acks=all?

| Setting | Behavior | Data Safety | Performance |
|---------|----------|-------------|-------------|
| `acks=0` | Producer doesn't wait for any acknowledgment | May lose data | Fastest |
| `acks=1` | Producer waits for leader to write | May lose if leader crashes before replication | Balanced |
| `acks=all` | Producer waits for leader + all ISR to write | No data loss (with min.insync.replicas >= 2) | Slowest |

### Q6: How does Kafka achieve exactly-once semantics?

**Answer:**  
Through two mechanisms:
1. **Idempotent Producer** (`enable.idempotence=true`): Kafka assigns a producer-id and sequence number to each message. Duplicates are detected and dropped.
2. **Transactions**: Producer wraps multiple sends in a transaction (`beginTransaction()`, `commitTransaction()`). Consumer reads with `isolation.level=read_committed` to see only committed data.

### Q7: What is the role of ZooKeeper in Kafka?

**Answer:**  
ZooKeeper manages cluster metadata: broker registration, leader election, topic configuration, and consumer group coordination. Starting with Kafka 3.3+ (KRaft mode), ZooKeeper is being replaced by Kafka's own built-in consensus protocol for metadata management.

### Q8: How would you choose the number of partitions?

**Answer:**  
Consider: (1) desired throughput / throughput per consumer, (2) number of consumers in the group, (3) memory and file handle overhead per partition, (4) rebalancing latency. A good starting point is 3-12 partitions. You can increase later but never decrease.

### Q9: What is a Dead Letter Topic (DLT)?

**Answer:**  
A DLT is a separate Kafka topic where messages that fail processing after all retries are sent. This prevents a bad message from blocking the consumer indefinitely. The DLT can be monitored and failed messages can be investigated and reprocessed.

### Q10: Kafka vs RabbitMQ — when to use which?

| Factor | Choose Kafka | Choose RabbitMQ |
|--------|-------------|-----------------|
| Throughput | Millions/sec | Thousands/sec |
| Message replay | Needed | Not needed |
| Ordering | Required | Not critical |
| Routing complexity | Simple (topics) | Complex (exchanges, bindings) |
| Use case | Event streaming, data pipelines | Task queues, request-reply |
| Consumer model | Pull-based | Push-based |
| Message retention | Long-term (days/weeks) | Immediate consumption |

---

## 13. Quick Reference Cheat Sheet

### Kafka Terminology in 1 Line Each

| Term | One-Line Definition |
|------|---------------------|
| **Topic** | A named stream of messages (like a category or channel) |
| **Partition** | An ordered, immutable sub-stream within a topic (unit of parallelism) |
| **Offset** | A sequential ID for each message within a partition (like a page number) |
| **Broker** | A single Kafka server that stores partitions and serves clients |
| **Cluster** | A group of brokers working together |
| **Producer** | An application that writes messages to Kafka topics |
| **Consumer** | An application that reads messages from Kafka topics |
| **Consumer Group** | A set of consumers sharing work on a topic (each partition → 1 consumer) |
| **Leader** | The primary replica handling all reads/writes for a partition |
| **Follower** | A backup replica that copies data from the leader |
| **ISR** | In-Sync Replicas — followers that are caught up with the leader |
| **Replication Factor** | Number of copies of each partition across brokers |
| **ZooKeeper/KRaft** | Manages cluster metadata, leader election, configuration |
| **Commit Log** | Kafka's internal append-only data structure on disk |
| **Segment** | A chunk of the commit log (rotated by size or time) |
| **High Watermark** | The offset up to which consumers can safely read |
| **DLT** | Dead Letter Topic — where failed messages go after retries exhausted |
| **Idempotent Producer** | A producer that prevents duplicates using sequence numbers |
| **Compacted Topic** | A topic that keeps only the latest value for each key |

### Command Cheat Sheet

```bash
# Topic Management
kafka-topics --create --topic <name> --partitions <n> --replication-factor <r> --bootstrap-server localhost:9092
kafka-topics --list --bootstrap-server localhost:9092
kafka-topics --describe --topic <name> --bootstrap-server localhost:9092
kafka-topics --alter --topic <name> --partitions <n> --bootstrap-server localhost:9092

# Producer
kafka-console-producer --topic <name> --bootstrap-server localhost:9092
kafka-console-producer --topic <name> --bootstrap-server localhost:9092 --property "key.separator=:" --property "parse.key=true"

# Consumer
kafka-console-consumer --topic <name> --bootstrap-server localhost:9092 --from-beginning
kafka-console-consumer --topic <name> --bootstrap-server localhost:9092 --group <group-id>

# Consumer Groups
kafka-consumer-groups --list --bootstrap-server localhost:9092
kafka-consumer-groups --describe --group <group-id> --bootstrap-server localhost:9092
kafka-consumer-groups --reset-offsets --group <group-id> --topic <name> --to-earliest --execute --bootstrap-server localhost:9092
```

### Spring Boot Quick Template

```
pom.xml → spring-kafka dependency
application.yml → bootstrap-servers, serializer/deserializer, group-id
KafkaTopicConfig.java → NewTopic beans
Producer.java → KafkaTemplate.send()
Consumer.java → @KafkaListener
```

---

> **Continue to [Part 2 (Advanced)](./Apache-Kafka-Complete-Learning-Guide-Part2.md):**
> - Kafka Streams for stream processing
> - Kafka Connect for data integration (JDBC, Elasticsearch, S3, Debezium CDC)
> - Schema Registry (Avro/Protobuf) for schema evolution
> - Log Compaction, Multi-broker clusters, and failure scenarios
> - Monitoring with Kafka UI, Prometheus, and Grafana
> - Security (SSL/TLS, SASL, ACLs)
> - Performance tuning and production best practices
> - Design patterns: Event Sourcing, CQRS, Saga, Transactional Outbox
> - Advanced Spring Boot Kafka patterns
