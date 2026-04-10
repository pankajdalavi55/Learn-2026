# System Design — FAANG Interview Questions & Answers

> Mostly asked questions in FAANG / MAANG interviews. Short, precise answers in bullet points. Organized topic-wise for quick revision.

---

## Topic 1: Fundamentals

---

### Q: What is the difference between latency and throughput?

- **Latency** = time for ONE operation (ms). **Throughput** = operations per unit time (req/s)
- Optimizing one often impacts the other — batching increases throughput but adds latency
- User-facing APIs: prioritize latency (< 100ms p99). Background jobs: prioritize throughput
- Always talk in percentiles (p50, p95, p99) — never averages

**Follow-up:** "How would you handle 10x traffic?"

- Horizontal scaling + caching + async processing + rate limiting + graceful degradation

---

### Q: Explain CAP theorem. Can you have all three?

- In a distributed system during a **network partition**, you must choose between **Consistency** (every read sees latest write) and **Availability** (every request gets a response)
- Network partitions WILL happen → you're really choosing CP or AP
- **CP examples:** ZooKeeper, etcd, PostgreSQL — refuse requests if uncertain
- **AP examples:** Cassandra, DynamoDB (default), CouchDB — serve potentially stale data
- "CA system" does not exist in distributed environments
- **PACELC** extends CAP: even without partition, there's a Latency vs Consistency trade-off

**Follow-up:** "Shopping cart — CP or AP?"

- AP — stale cart is better than no cart; merge conflicts on checkout

---

### Q: Horizontal vs vertical scaling — when to use each?

- **Vertical (scale up):** bigger machine — simple, SPOF, hits hardware ceiling, exponential cost
- **Horizontal (scale out):** more machines — complex, redundant, theoretically unlimited, linear cost
- Web/API tier → horizontal from day one (stateless, easy)
- Database → vertical first, then read replicas, then sharding when necessary
- Horizontal requires: stateless services, load balancer, distributed state management

**Follow-up:** "Prerequisites for horizontal scaling?"

- Stateless design, externalized state (Redis/DB), health checks, graceful shutdown, fast startup

---

### Q: What is idempotency and why does it matter?

- An operation is idempotent if performing it multiple times has the same effect as once
- **Critical because:** networks fail, retries are inevitable — without idempotency you get double charges, duplicate orders
- Naturally idempotent: GET, PUT, DELETE, SET operations
- Not idempotent (need handling): POST, INCREMENT, payments
- **Implementation:** client sends unique idempotency key → server deduplicates → returns cached result on retry
- Stripe stores idempotency keys for 24 hours

**Follow-up:** "What if server crashes between processing and storing the key?"

- Use database transaction: process + store key atomically; on crash, neither is committed

---

### Q: Strong consistency vs eventual consistency — when to use each?

- **Strong:** all reads see latest write (expensive — sync replication, consensus protocols)
- **Eventual:** reads may lag behind writes temporarily (cheaper — async replication)
- Strong for: payments, inventory (can't oversell), booking, access control
- Eventual for: social feeds, recommendations, analytics, notifications
- **Read-your-writes** is a useful middle ground — you see your own writes immediately, others see them eventually
- Conflict resolution for eventual: Last-Write-Wins, vector clocks, CRDTs, app-level merge

---

### Q: Stateful vs stateless services — impact on scaling?

- **Stateless:** no client data stored between requests — any instance handles any request, easy to scale/replace
- **Stateful:** maintains session across requests — requires sticky sessions, state migration on failure
- Golden rule: web/API tier = stateless, data tier = stateful (managed services)
- Move state to: JWT (client-side), Redis (distributed cache), database
- Stateful is OK for: WebSocket servers (use presence service), databases, Kafka brokers

---

## Topic 2: Networking & API Design

---

### Q: REST vs gRPC — when to use which?

- **REST:** HTTP + JSON, universal browser support, human-readable, mature caching → public APIs, CRUD, third-party integrations
- **gRPC:** HTTP/2 + Protobuf (binary), 7-10x smaller payloads, native streaming, strongly typed → internal microservices, streaming, performance-critical paths
- Common pattern: Public → REST → API Gateway → gRPC → Internal services
- gRPC limitation: no native browser support (needs grpc-web proxy)

**Follow-up:** "What about GraphQL?"

- Client specifies exactly what data it needs — solves over/under-fetching
- Best for: complex frontend data needs, mobile apps (bandwidth savings)
- Trade-off: harder to cache, complex query optimization, N+1 risk on server

---

### Q: How would you version an API?

- **URL versioning** (`/v1/users`) — most common, explicit, cacheable (used by Stripe, GitHub, Twitter)
- **Header versioning** — cleaner URLs, harder to test
- Version from day one; support N-1 for 6-12 months
- Add optional fields (safe); removing/renaming fields is breaking
- Use sunset headers to announce deprecation

---

### Q: Offset pagination vs cursor pagination?

- **Offset (`OFFSET 100 LIMIT 20`):** simple, supports jump-to-page, but O(n) for large offsets and inconsistent with inserts
- **Cursor (`WHERE id > last_id LIMIT 20`):** fast at any depth, consistent with inserts, no jump-to-page
- Use cursor for: infinite scroll, large datasets (millions+), real-time feeds
- Use offset for: small datasets, admin dashboards needing page numbers

---

### Q: How would you design a rate limiter?

- **Token Bucket** (most common): allows bursts up to bucket capacity, steady rate = refill rate (used by AWS, Stripe)
- **Sliding Window Counter:** good accuracy + low memory (best balance)
- **Fixed Window:** simple but 2x burst at window boundaries
- Distributed: Redis `INCR` + `EXPIRE` (atomic); accept approximate counts across nodes
- Respond with: HTTP 429 + `Retry-After` + `X-RateLimit-Remaining` headers
- Rate limit by user ID (not just IP) — prevents multi-account bypass

**Follow-up:** "How to handle different API tiers (free vs premium)?"

- Lookup tier from user context → apply different bucket configs per tier

---

### Q: Explain OAuth2 and JWT in microservices.

- **Authentication (AuthN):** "Who are you?" — **Authorization (AuthZ):** "What can you do?"
- JWT: stateless token (Header.Payload.Signature) — any service can validate without DB lookup
- Trade-off: can't revoke JWT until expiry → use short-lived access tokens (15 min) + refresh tokens (7 days)
- Microservices pattern: Client → Auth Service (JWT) → API Gateway (validates JWT) → downstream services trust gateway context
- Store in httpOnly cookies (not localStorage — XSS vulnerable)
- OAuth2 flows: Authorization Code + PKCE for mobile/SPA, Client Credentials for machine-to-machine

---

## Topic 3: Load Balancing

---

### Q: L4 vs L7 load balancing?

- **L4 (Transport):** operates at TCP/UDP — routes by IP:port, fast, low CPU, no content inspection
- **L7 (Application):** operates at HTTP — inspects headers/URLs/cookies, SSL termination, path-based routing
- L4 for: database connections, maximum throughput, pass-through
- L7 for: HTTP APIs, URL routing, canary deployments, advanced health checks
- Common architecture: Internet → L4 LB → L7 LB → App servers
- AWS: NLB (L4), ALB (L7)

---

### Q: Which LB algorithm would you choose?

- **Round Robin:** homogeneous stateless services (simple, fair)
- **Least Connections:** long-lived connections like WebSocket, DB pools
- **Consistent Hashing:** cache servers (minimizes cache misses on topology change)
- **Weighted Round Robin:** mix of server capacities
- **IP Hash:** crude session affinity (breaks with NAT)

**Follow-up:** "LB is a single point of failure — how to handle?"

- Active-passive pair (VRRP), active-active with DNS, managed cloud LBs (multi-AZ by default)

---

### Q: What are sticky sessions and why avoid them?

- All requests from a client routed to same server — needed when state stored in server memory
- Problems: uneven load, scaling issues, SPOF per user session, complicates deployments
- Better alternatives: externalize session to Redis, use JWT (client-side state), design stateless services
- Acceptable for: WebSocket connections (inherently sticky — design for reconnection)

---

## Topic 4: Databases

---

### Q: SQL vs NoSQL — when to choose each?

- **SQL (PostgreSQL, MySQL):** ACID, complex joins, well-defined schema, ad-hoc queries → orders, payments, inventory, CRM
- **NoSQL Key-Value (Redis, DynamoDB):** sub-ms lookups by key → caching, sessions
- **NoSQL Document (MongoDB):** flexible schema, nested data → product catalogs, user profiles with varying fields
- **NoSQL Column (Cassandra, HBase):** write-heavy, time-series, linear scalability → IoT, logging, activity feeds
- **Graph (Neo4j):** relationship-heavy → social networks, recommendations
- Real-world: polyglot persistence — PostgreSQL for transactions + Redis for cache + Cassandra for logs + Elasticsearch for search

---

### Q: How does database indexing work?

- Index = B+ tree data structure → O(log n) lookups vs O(n) full scan
- **Composite index:** leftmost prefix rule — `(A, B, C)` serves queries on `(A)`, `(A,B)`, or `(A,B,C)` only
- **Covering index:** includes all queried columns → index-only scan, no table lookup
- Index WHERE, JOIN, ORDER BY columns
- **Don't index:** small tables, low cardinality (boolean), frequently updated columns, write-heavy tables
- Every index slows down writes (INSERT/UPDATE/DELETE must update index too)

**Follow-up:** "Query uses index but still slow?"

- Leading wildcard, function on column, type mismatch, large OFFSET, stale statistics → use `EXPLAIN ANALYZE`

---

### Q: Explain database sharding.

- Splitting data across multiple DB instances by a shard key
- **Hash-based** (`hash(user_id) % N`): even distribution, no range queries
- **Range-based** (`user_id 1-1M → shard 1`): range queries easy, hot spots possible
- **Good shard key:** high cardinality, even distribution, matches query patterns (usually user_id or tenant_id)
- **Bad shard key:** timestamp (today's shard always hot), country (uneven), status (low cardinality)
- Cross-shard joins: nearly impossible → application-level joins
- Cross-shard transactions: use saga pattern
- Shard as last resort — after: vertical scaling, read replicas, caching

**Follow-up:** "How to reshard without downtime?"

- Dual-write → background migration (chunks + CDC) → gradual read switchover (feature flags) → write switchover → cleanup

---

### Q: What is the hot partition problem?

- One partition/shard gets disproportionate traffic → bottleneck
- Causes: celebrity user, viral product, all writes to today's date partition
- **Solutions:** salt/suffix key (`id + random(0-99)`) with scatter-gather reads, dedicated shards for hot entities, caching hot data in Redis
- Prevention: high-cardinality keys, consistent hashing, monitor per-shard QPS

---

### Q: ACID vs BASE?

- **ACID:** Atomicity, Consistency, Isolation, Durability — SQL databases, strong guarantees
- **BASE:** Basically Available, Soft state, Eventually consistent — NoSQL, high availability
- ACID for: financial transactions, inventory, bookings
- BASE for: social feeds, analytics, logging, recommendations
- Most systems use both: ACID for critical path, BASE for everything else

---

### Q: Explain isolation levels.

- **Read Uncommitted:** see uncommitted data (dirty reads) — almost never used
- **Read Committed:** only see committed data — PostgreSQL default, sufficient for most apps
- **Repeatable Read:** same query returns same results within transaction — MySQL InnoDB default
- **Serializable:** transactions execute as if serial — slowest, safest, for financial operations
- **Optimistic locking:** check version at commit (`WHERE version = 5`) — low contention
- **Pessimistic locking:** lock rows during transaction — high contention

---

## Topic 5: Caching

---

### Q: What caching strategy would you use?

- **Cache-Aside (most common):** app checks cache → miss → DB → store in cache; write → DB → invalidate cache
- **Write-Through:** write → cache → DB synchronously — strong consistency, higher write latency
- **Write-Behind:** write → cache → async DB — fastest writes, data loss risk
- **Read-Through:** cache itself fetches from DB on miss — cleaner app logic
- Use cache-aside for most read-heavy systems; write-through for critical data; write-behind for analytics/logs

---

### Q: How do you handle cache invalidation?

- "Two hard problems in CS: cache invalidation and naming things"
- **Write-invalidate + TTL** (most common): delete cache on DB write, TTL as safety net
- **Event-driven** (microservices): DB update → Kafka event → cache delete
- **Versioned keys** (`product:123:v2`): new version = new key, old key expires
- Pitfalls: race condition (cache updated before DB commit), lost invalidation (cache delete fails → stale forever), thundering herd after invalidation
- Best practice: combine write-invalidate + TTL + retry on delete failure + monitoring cache hit ratio

---

### Q: What is cache stampede and how to prevent it?

- Popular key expires → thousands of requests hit DB simultaneously → overload
- **Mutex/locking:** one request fetches, others wait → prevents DB overload, adds wait latency
- **Probabilistic early expiry:** refresh before TTL with jitter → spreads load
- **Stale-while-revalidate:** serve stale data, refresh in background → best UX
- **Background refresh:** proactively refresh hot keys before expiry → no miss for hot data
- Best practice: combine stale-while-revalidate + TTL jitter + locking for expensive computations

---

### Q: Cache eviction — LRU vs LFU vs TTL?

- **LRU (Least Recently Used):** evict oldest-accessed — best general-purpose default
- **LFU (Least Frequently Used):** evict least-accessed — better for stable popularity patterns, risk of stale popular items
- **TTL (Time To Live):** expire after fixed time — essential for data freshness
- Redis uses approximate LRU/LFU (sampling-based) for performance
- Production default: LRU + TTL combined

---

### Q: Where would you place caches in the architecture?

- **Browser/Client:** static assets, API responses (Cache-Control headers) — fastest, hardest to invalidate
- **CDN:** static assets at edge locations — reduces origin load, global reach
- **API Gateway:** API response caching — seconds-to-minutes TTL
- **Application (Redis/Memcached):** DB query results, computed data — most control
- **Database:** query cache, buffer pool — automatic, transparent
- Request flow: Browser → CDN → Gateway → App Cache → Redis → DB

---

## Topic 6: Distributed Systems

---

### Q: Explain consistent hashing.

- **Problem:** `hash(key) % N` — if N changes, almost all keys remap → cache invalidation storm
- **Solution:** arrange servers on a virtual ring (0 to 2^32), key routes to first server clockwise
- Adding/removing server → only ~1/N keys move (vs ~100% with modulo)
- **Virtual nodes:** each server has multiple positions on ring → even distribution, gradual failover
- Used by: Cassandra, DynamoDB, Memcached, CDN request routing

---

### Q: Leader-based vs leaderless replication?

- **Leader-based:** one primary accepts writes, replicates to followers — strong consistency, simpler, SPOF risk
  - PostgreSQL, MySQL, MongoDB — best for most OLTP
- **Leaderless:** all nodes accept writes, quorum-based — no SPOF, higher availability, conflict resolution needed
  - Cassandra, DynamoDB, Riak — best when availability > consistency
- Default to leader-based; choose leaderless only when availability is paramount

**Follow-up:** "What is split-brain?"

- Network partition → two nodes think they're leader → divergent data
- Prevention: quorum election (majority vote), fencing tokens, STONITH, consensus layer (etcd)

---

### Q: Explain quorum reads/writes.

- N = replicas, W = write ack count, R = read count
- **W + R > N** → at least one node in read overlaps with write → strong consistency
- N=3, W=2, R=2 → strong consistency with tolerance for 1 node failure
- Tunable: W=1, R=1 (fast, eventual); W=N, R=1 (slow writes, fast consistent reads)
- Cassandra/DynamoDB support per-query tunable consistency

---

### Q: Kafka vs RabbitMQ vs SQS?

- **Kafka:** log-based, retains messages (replay), ordering per partition, 100K+ msg/s → event streaming, audit logs, data pipelines
- **RabbitMQ:** traditional queue, deletes after consumption, complex routing (fanout/topic), low latency → task queues, RPC, notifications
- **SQS:** fully managed, auto-scales, standard (at-least-once) or FIFO (exactly-once) → AWS-native, simple decoupling
- Kafka when: replay needed, high throughput, event sourcing
- RabbitMQ when: complex routing, low per-message latency
- SQS when: AWS ecosystem, managed simplicity

---

### Q: At-least-once vs exactly-once delivery?

- **At-most-once:** fire and forget — may lose messages (OK for metrics/logs)
- **At-least-once:** retry until ack — duplicates possible (default for most systems)
- **Exactly-once:** no loss, no duplicates — expensive, rare (Kafka transactions)
- **Best practice:** use at-least-once + idempotent consumers (deduplicate by message ID, DB unique constraints)
- Exactly-once is effectively "at-least-once delivery + idempotent processing"

---

### Q: What is event-driven architecture?

- Services communicate by publishing/subscribing to events instead of direct calls
- **Benefits:** loose coupling, independent scaling, resilience (events queue up if subscriber down)
- **Challenges:** eventual consistency, harder debugging (no linear call stack), schema evolution
- Event design: eventId, eventType, timestamp, aggregateId, data
- Use events when: fan-out, loose coupling, eventual consistency OK
- Use sync calls when: immediate response, simple flow, strong consistency needed

---

### Q: How do you handle backpressure?

- Producers generate data faster than consumers can process → queue grows unbounded
- **Solutions:** bounded queues (reject overflow), throttle producers, auto-scale consumers, sample/drop non-critical data
- Kafka: consumers control pace via offset; monitor lag metric
- RabbitMQ: memory/disk alarms block producers; prefetch limits consumers
- Alert on queue depth > threshold

---

## Topic 7: Scalability Patterns

---

### Q: Monolith vs microservices — when to choose?

- **Start with monolith/modular monolith** — simpler, faster iteration, less operational overhead
- **Extract to microservices when:** team > 50 engineers, components need independent scaling, different release cadences, clear bounded contexts
- Microservices add: distributed tracing, eventual consistency, saga complexity, infrastructure cost
- 80% of systems should stay as modular monolith
- Decompose by: business capability, DDD bounded context, team ownership, data ownership
- Anti-pattern: "distributed monolith" — services that must deploy together

**Follow-up:** "How do you decide service boundaries?"

- Each service: one team, one reason to change, owns its data, independently deployable
- When in doubt, keep together — can always split later

---

### Q: Explain the Saga pattern.

- Distributed transaction as sequence of local transactions + compensating actions on failure
- **Choreography:** services publish/subscribe events — loose coupling, harder to trace (good for 2-3 steps)
- **Orchestration:** central coordinator directs steps — clear flow, easier debugging (good for complex flows)
- Example: CreateOrder → ReserveInventory → ChargePayment → ConfirmOrder; if payment fails → ReleaseInventory → CancelOrder
- Every step and compensation must be idempotent
- Use outbox pattern for reliable event publishing

**Follow-up:** "What if compensation fails?"

- Retry with backoff, dead letter queue, alert for manual intervention; compensations must be idempotent

---

### Q: Explain the Circuit Breaker pattern.

- Prevents cascading failures by stopping calls to a failing service
- **Closed:** normal operation → **Open:** fail fast (after threshold failures) → **Half-Open:** test with limited requests
- Configuration: 5-10 failure threshold, 30-60s open timeout, 1-3 half-open test requests
- Count 5xx and timeouts as failures; NOT 4xx (client errors)
- Combine with fallback: cached data, default values, degraded experience
- Use libraries (Resilience4j, Polly) — don't build from scratch

**Follow-up:** "Circuit Breaker vs Retry?"

- Retry = transient failures (try again). Circuit Breaker = sustained failures (stop trying)
- Use both: retry inside circuit breaker — retries handle blips, circuit breaks on prolonged failure

---

### Q: Explain retry with exponential backoff.

- `delay = base_delay × 2^attempt + random_jitter`
- **Jitter is critical** — without it all clients retry simultaneously → thundering herd
- Only retry: transient errors (503, timeouts) — NOT permanent errors (400, 404)
- Only retry idempotent operations (or with idempotency key)
- Max 2-3 retries with cap on delay

---

### Q: What is the Bulkhead pattern?

- Isolate resources per dependency — one failure doesn't affect others
- Analogy: ship compartments — one floods, others stay dry
- Types: separate thread pools, semaphore limits, connection pool isolation
- Without bulkhead: slow Service C uses all threads → Service A and B also fail
- With bulkhead: slow Service C only affects its own pool

---

### Q: Explain timeouts in distributed systems.

- Set timeouts on ALL external calls — without them, slow services hold resources indefinitely
- **Connection timeout:** 1-5s. **Read timeout:** varies by operation. **Total timeout:** end-to-end budget
- **Timeout propagation:** each hop's timeout shorter than caller's — prevents caller timeout while downstream continues
- gRPC deadline propagation: client sets deadline, propagated through entire call chain
- Combine with circuit breaker for resilience

---

## Topic 8: Observability & Reliability

---

### Q: What are the three pillars of observability?

- **Metrics:** numerical measurements over time (Prometheus, Datadog) — dashboards, alerts, trends
- **Logs:** discrete events with details (ELK, Splunk) — debugging, auditing
- **Traces:** request flow across services (Jaeger, Zipkin) — understanding distributed call chains
- Correlation IDs propagated across all services for traceability
- Structured logs (JSON), no PII, appropriate levels

---

### Q: Explain RED and USE metrics.

- **RED (for services):** Rate (req/s), Errors (error rate), Duration (p50/p95/p99)
- **USE (for resources):** Utilization (% used), Saturation (queue depth), Errors
- Monitor RED for every service; USE for every infrastructure component
- Alert on symptoms (RED), investigate causes (USE)

---

### Q: What are SLIs, SLOs, and SLAs?

- **SLI (Indicator):** what you measure — p99 latency, error rate, availability %
- **SLO (Objective):** internal target — "p99 < 200ms", "99.95% success rate"
- **SLA (Agreement):** external contract with penalties — always less strict than SLO (buffer)
- **Error budget** = 100% − SLO → for 99.9%: 43 min downtime/month allowed
- Alert on SLO burn rate (not individual errors) to avoid alert fatigue
- 99.9% (43 min/mo) is a good starting point; 99.99% (4.3 min/mo) requires significant investment

---

## Topic 9: System Design Building Blocks (Quick Reference)

---

### Q: How would you design a URL shortener?

- Write: generate short code (base62 of auto-increment ID or hash), store mapping in DB
- Read (100:1 ratio): check cache (Redis) → DB fallback → 301/302 redirect
- Scale: read replicas, cache-aside, CDN for popular links
- Analytics: async write to Kafka → analytics pipeline
- Shard by short code hash

---

### Q: How would you design a notification system?

- Ingest: events from services → Kafka topic
- Fan-out: notification service reads events → determines recipients + channels (push/email/SMS)
- Delivery: per-channel workers with retry + exponential backoff
- Deduplication: idempotency key per (user, event, channel)
- User preferences: read from cache, respect Do Not Disturb
- Scale: partition Kafka by user_id, auto-scale workers by queue depth

---

### Q: How would you design a chat system?

- Real-time: WebSocket connections (stateful) with presence service
- Message storage: write to DB + push to recipient's WebSocket
- Offline delivery: store in inbox, deliver on reconnect
- Group chat: fan-out on write (small groups) or fan-out on read (large groups)
- Scale: shard by conversation_id, connection servers behind L4 LB
- Read receipts and typing indicators: ephemeral pub/sub (no persistence)

---

### Q: How would you design a news feed?

- **Fan-out on write (push):** on post, write to all followers' feeds → fast reads, expensive writes (celebrity problem)
- **Fan-out on read (pull):** on read, query all followees → slow reads, simple writes
- **Hybrid:** push for normal users, pull for celebrities (> 1M followers)
- Feed storage: Redis sorted set per user (score = timestamp)
- Scale: shard feed by user_id, cache hot feeds, async fan-out via Kafka

---

### Q: How would you design a rate limiter service?

- Token Bucket algorithm per (client_id, endpoint)
- Distributed: Redis `INCR` + `EXPIRE` with Lua script for atomicity
- Deploy at API Gateway (first line of defense) + application layer (business limits)
- Respond: 429 + `Retry-After` header
- Sliding window counter for smoother enforcement
- Shard Redis by client_id for scale

---

## Topic 10: Estimation & Numbers to Know

---

### Q: Key latency numbers every engineer should know?


| Operation                            | Latency    |
| ------------------------------------ | ---------- |
| L1 cache reference                   | 0.5 ns     |
| L2 cache reference                   | 7 ns       |
| Main memory reference                | 100 ns     |
| SSD random read                      | 150 μs     |
| HDD seek                             | 10 ms      |
| Network round trip (same DC)         | 0.5 ms     |
| Network round trip (cross-continent) | 150 ms     |
| Redis GET                            | 0.1-0.5 ms |
| Simple DB query (indexed)            | 1-5 ms     |
| Complex DB query                     | 10-100 ms  |


---

### Q: Common capacity estimates?


| Metric                                 | Approximate            |
| -------------------------------------- | ---------------------- |
| Requests/day (100M users, 10 req/user) | 1B req/day ≈ 12K req/s |
| Storage (1B records × 1KB each)        | 1 TB                   |
| Bandwidth (12K req/s × 10KB)           | 120 MB/s               |
| Redis throughput                       | 100K+ ops/s per node   |
| Kafka throughput                       | 1M+ msg/s per cluster  |
| Single PostgreSQL                      | 10K-50K queries/s      |


---

### Q: Availability in 9s?


| Availability      | Downtime/year | Downtime/month |
| ----------------- | ------------- | -------------- |
| 99% (two 9s)      | 3.65 days     | 7.2 hours      |
| 99.9% (three 9s)  | 8.76 hours    | 43 minutes     |
| 99.99% (four 9s)  | 52.6 minutes  | 4.3 minutes    |
| 99.999% (five 9s) | 5.26 minutes  | 26 seconds     |


---

> **Revision tip:** Read this file the morning of your interview. Focus on trade-offs and "it depends" answers — interviewers value reasoning over memorized facts.





Got it — you want a **clean, structured segregation of System Design topics** (not just a list, but grouped like an architect thinks).

Here is the **industry-grade categorization** 👇

---

# 🧠 1. FOUNDATION (Must Know Before Anything)

These are **decision-driving concepts**.

### Core Concepts

- Functional vs Non-Functional Requirements
- Latency vs Throughput
- Availability vs Consistency
- CAP Theorem
- ACID vs BASE

### Estimation

- QPS (Queries per second)
- Storage estimation
- Bandwidth calculation

---

# 🏗️ 2. ARCHITECTURE LAYER

### Architecture Patterns

- Monolith
- Microservices
- Event-driven architecture
- Serverless architecture

### Communication Styles

- Sync (REST, gRPC)
- Async (Queues, Events)

---

# 🌐 3. API & INTERFACE DESIGN

### API Design

- REST principles
- GraphQL basics
- Idempotency
- Pagination (cursor-based)

### Gateway Layer

- API Gateway
- Rate limiting
- Request validation

---

# ⚙️ 4. COMPUTE LAYER (Backend Processing)

### Application Layer

- Stateless services
- Business logic separation

### Concurrency

- Multi-threading
- Async programming
- Non-blocking I/O
- Virtual Threads (Java 21)

---

# 🗄️ 5. DATA LAYER

### Database Types

- SQL → MySQL, PostgreSQL
- NoSQL → MongoDB, Cassandra, DynamoDB

### Data Engineering Concepts

- Indexing
- Partitioning (Sharding)
- Replication
- Schema design

---

# ⚡ 6. PERFORMANCE LAYER

### Caching

- Redis, Memcached
- Cache strategies:
  - Cache Aside
  - Write Through
  - Write Back

### CDN

- Static content delivery
- Geo-distribution

---

# 📩 7. ASYNCHRONOUS PROCESSING

### Messaging Systems

- Kafka
- RabbitMQ
- AWS SQS

### Patterns

- Pub/Sub
- Event streaming
- Event sourcing (advanced)

---

# 📂 8. STORAGE SYSTEMS

### Types

- Object storage (S3)
- Blob storage
- File systems

### Use Cases

- Media (images/videos)
- Logs
- Backups

---

# 🔍 9. SEARCH & INDEXING

### Tools

- Elasticsearch
- Solr

### Concepts

- Inverted index
- Ranking algorithms
- Full-text search

---

# 📈 10. SCALABILITY & DISTRIBUTION

### Scaling Types

- Vertical scaling
- Horizontal scaling

### Techniques

- Load balancing
- Auto scaling
- Data partitioning

---

# ⚖️ 11. CONSISTENCY & DISTRIBUTED SYSTEMS

### Models

- Strong consistency
- Eventual consistency
- Causal consistency

### Distributed Concepts

- Leader election
- Consensus (Raft, Paxos)
- Distributed locks

---

# 🔐 12. SECURITY LAYER

### Authentication & Authorization

- JWT
- OAuth
- RBAC

### Protection

- HTTPS/TLS
- Rate limiting
- DDoS protection

---

# 🔁 13. RELIABILITY & FAULT TOLERANCE

### Patterns

- Retry
- Circuit Breaker
- Bulkhead
- Graceful degradation

---

# 📊 14. OBSERVABILITY

### Monitoring

- Metrics (Prometheus)

### Logging

- ELK Stack

### Tracing

- Jaeger, Zipkin

---

# 🚀 15. DEVOPS & DEPLOYMENT

### Tools & Practices

- CI/CD pipelines
- Docker
- Kubernetes

### Deployment Strategies

- Blue-Green
- Canary releases

---

# 🧩 16. ADVANCED TOPICS (Senior+)

- Event Sourcing
- CQRS
- Distributed Transactions (2PC, Saga)
- Data Lake vs Data Warehouse
- Stream Processing

---

# 🧠 How to Use This (Important)

For interviews, think in **this exact order**:

```
1. Foundation (requirements + estimation)
2. Architecture
3. Data layer
4. Scaling + caching
5. Async systems
6. Reliability
7. Trade-offs

```

---

# 🔥 If You Want Next Level

Since you’re targeting **high-paying backend / global roles**, next logical step:

👉 I can convert this into:

- **30-day System Design roadmap (aligned to your Java + AWS stack)**
- OR
- **Hands-on designs (WhatsApp, YouTube, Uber)** with deep backend focus

Just say:  
👉 “Roadmap” or “Start design WhatsApp”