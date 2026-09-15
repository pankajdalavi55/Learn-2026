Since you're preparing for **Senior Backend (6 YOE)** and targeting **FAANG + Product Companies + Germany**, don't study System Design randomly. Learn it layer by layer so every topic builds on the previous one.

---

# Phase 1: Foundation (Must Know)

These are the building blocks used in almost every system.

## 1. Scalability

* Vertical Scaling
* Horizontal Scaling
* Stateless vs Stateful services
* Bottlenecks
* Capacity Planning

**Interview Questions**

* How do you scale an application?
* When do you choose vertical over horizontal scaling?

---

## 2. Back-of-the-envelope Estimation

Learn to calculate:

* QPS
* Storage
* Bandwidth
* Memory
* Cache Size
* Number of Servers

Know important numbers:

* 1 KB = 10³ Bytes
* 1 MB = 10⁶ Bytes
* 1 GB = 10⁹ Bytes

Latency cheat sheet

| Component       | Approx    |
| --------------- | --------- |
| L1 Cache        | 1 ns      |
| L2 Cache        | 4 ns      |
| RAM             | 100 ns    |
| SSD             | 100 μs    |
| HDD             | 10 ms     |
| Network Same DC | 0.5 ms    |
| Cross Region    | 50-150 ms |

---

## 3. Networking Basics

Learn

* HTTP/HTTPS
* TCP vs UDP
* TLS
* DNS
* CDN
* Load Balancer
* Reverse Proxy
* API Gateway

Interview

* Difference between Load Balancer and API Gateway?
* Why HTTPS?
* Why CDN?

---

## 4. Database Fundamentals

Relational Database

* Primary Key
* Foreign Key
* Normalization
* ACID
* Transactions

NoSQL

* Key Value
* Document
* Column
* Graph

Comparison

* SQL vs NoSQL
* PostgreSQL vs MongoDB

---

# Phase 2: Data Layer

## 5. Indexing

Learn

* B Tree
* Hash Index
* Composite Index
* Covering Index
* Clustered Index
* Secondary Index

Questions

* Why B Tree?
* How index works internally?

---

## 6. Partitioning

* Horizontal
* Vertical
* Range
* Hash
* List

---

## 7. Sharding

Learn

* Why Sharding
* Shard Key
* Hot Shards
* Resharding
* Consistent Hashing

Questions

* Why Consistent Hashing?
* What happens when new shard added?

---

## 8. Replication

* Leader-Follower
* Multi Leader
* Leaderless

Concepts

* Read Replicas
* Failover
* Replication Lag

---

# Phase 3: Caching

## 9. Caching

Types

* Browser Cache
* CDN
* API Cache
* Application Cache
* Database Cache
* Redis

Policies

* LRU
* LFU
* FIFO

Strategies

* Cache Aside
* Read Through
* Write Through
* Write Around
* Write Back

Problems

* Cache Stampede
* Cache Penetration
* Cache Avalanche

---

# Phase 4: Distributed Systems

## 10. CAP Theorem

Understand

* Consistency
* Availability
* Partition Tolerance

Questions

* Why can't we have all three?

---

## 11. Consistency Models

* Strong
* Eventual
* Causal
* Read Your Writes

---

## 12. Distributed Transactions

* Two Phase Commit
* Saga Pattern

Types

* Choreography
* Orchestration

---

## 13. Idempotency

Learn

* Idempotency Keys
* Duplicate Requests
* Payment Systems

---

## 14. Distributed Locking

* Redis Lock
* ZooKeeper
* etcd

---

# Phase 5: Messaging

## 15. Message Queues

Learn

* Kafka
* RabbitMQ
* SQS

Concepts

* Producer
* Consumer
* Topic
* Partition
* Offset
* Consumer Group

Questions

* Why Kafka?
* Ordering?
* Retry?
* DLQ?

---

## 16. Event Driven Architecture

Learn

* Event
* Event Bus
* Event Sourcing
* CQRS

---

# Phase 6: Reliability

## 17. Load Balancing

Algorithms

* Round Robin
* Least Connections
* Weighted
* Consistent Hashing

---

## 18. Fault Tolerance

Learn

* Retry
* Timeout
* Circuit Breaker
* Bulkhead
* Rate Limiter

---

## 19. High Availability

* Active Active
* Active Passive
* Multi Region
* Disaster Recovery

---

# Phase 7: Security

## 20. Authentication

* JWT
* OAuth2
* OpenID Connect
* Session
* API Keys

---

## 21. Authorization

* RBAC
* ABAC
* ACL
* Policy Based Access

Since your project is based on authorization, know this deeply.

---

# Phase 8: Observability

## 22. Logging

* Centralized Logging
* ELK

---

## 23. Metrics

* Prometheus
* Grafana

Golden Signals

* Latency
* Traffic
* Errors
* Saturation

---

## 24. Tracing

* OpenTelemetry
* Jaeger
* Zipkin

---

# Phase 9: Cloud & Deployment

## 25. Containers

* Docker
* Images
* Layers
* Volumes

---

## 26. Kubernetes

* Pod
* Deployment
* ReplicaSet
* Service
* Ingress
* ConfigMap
* Secret

---

## 27. CI/CD

* Build
* Test
* Deploy
* Rollback

Deployment

* Blue Green
* Canary
* Rolling

---

# Phase 10: System Design Patterns

These are the patterns interviewers expect.

* Caching
* Queue
* Pub/Sub
* Sharding
* Replication
* CQRS
* Event Sourcing
* Saga
* Outbox Pattern
* API Gateway
* Service Discovery
* Sidecar
* Leader Election
* Rate Limiter
* Bloom Filter
* Consistent Hashing

---

# Phase 11: Design Popular Systems

Practice these in order of increasing difficulty:

### Beginner

* URL Shortener
* TinyURL
* Pastebin
* Parking Lot
* Library Management
* Rate Limiter

### Intermediate

* Chat Application
* Notification Service
* News Feed
* Distributed Cache
* Search Autocomplete
* File Storage (Google Drive)
* Logging System

### Advanced

* YouTube
* Uber
* WhatsApp
* Instagram
* Netflix
* Google Docs
* Payment Gateway
* Ride Sharing
* Food Delivery
* Online Judge (LeetCode)
* Ticket Booking (BookMyShow)
* E-commerce (Amazon)

---

# Phase 12: Interview Framework

For every design question, follow this sequence:

1. Clarify requirements (functional & non-functional)
2. Estimate scale (QPS, storage, bandwidth)
3. Identify core entities
4. Design APIs
5. Design database schema
6. Draw the high-level architecture
7. Explain request flow
8. Identify bottlenecks
9. Add scalability components (cache, queues, sharding, replication)
10. Discuss consistency and trade-offs
11. Address failure scenarios and reliability
12. Cover monitoring, security, and deployment

---

# Recommended Learning Roadmap (12 Modules)

1. Foundations & Estimation
2. Networking & APIs
3. Databases & Indexing
4. Partitioning, Sharding & Replication
5. Caching
6. Distributed Systems Concepts
7. Messaging & Event-Driven Architecture
8. Reliability & Fault Tolerance
9. Security & Authorization
10. Observability & Operations
11. Cloud, Containers & Deployment
12. End-to-End System Design Case Studies

Mastering these modules will cover the vast majority of topics expected in senior backend and system design interviews.
