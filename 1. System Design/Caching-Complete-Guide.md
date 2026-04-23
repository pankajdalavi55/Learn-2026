# Caching — Complete FAANG Interview Guide

**Navigation:** [← Core Building Blocks](01-core-building-blocks.md) | [Distributed Systems →](02-distributed-systems.md)

> A comprehensive, FAANG-level caching guide for backend engineers preparing for senior/staff-level system design interviews.  
> Covers fundamentals through distributed caching, with real-world production scenarios, code, and interview strategies.

---

## Table of Contents

1. [Fundamentals](#1-fundamentals)
2. [Caching Strategies](#2-caching-strategies)
3. [Cache Eviction Policies](#3-cache-eviction-policies)
4. [Cache Invalidation](#4-cache-invalidation)
5. [Distributed Caching](#5-distributed-caching)
6. [Real-World Problems](#6-real-world-problems)
7. [System Design Integration](#7-system-design-integration)
8. [Concurrency & Performance](#8-concurrency--performance)
9. [Monitoring & Metrics](#9-monitoring--metrics)
10. [Tech-Specific: Redis, Memcached, Spring Boot](#10-tech-specific-redis-memcached-spring-boot)
11. [Interview Questions & Answers](#11-interview-questions--answers)
12. [Common Mistakes & How to Answer](#12-common-mistakes--how-to-answer)

---

## 1. Fundamentals

### 1.1 What is Caching?

Caching is storing a **copy of data in a faster storage layer** so future requests for that data are served faster than accessing the primary source.

The core insight: **trade memory (cheap) for latency (expensive).**

```
┌──────────────────────────────────────────────────────────────────┐
│                        REQUEST FLOW                              │
│                                                                  │
│   Client ──► Cache ──[HIT]──► Return immediately (~1ms)         │
│                │                                                 │
│              [MISS]                                              │
│                │                                                 │
│                ▼                                                 │
│           Data Source (DB / API / Disk) ──► Return (~50-500ms)   │
│                │                                                 │
│                ▼                                                 │
│          Populate Cache                                          │
└──────────────────────────────────────────────────────────────────┘
```

### 1.2 Why Caching is Used

| Goal | Without Cache | With Cache |
|------|--------------|------------|
| **Latency** | 50–500ms (DB round-trip) | 0.1–5ms (memory access) |
| **Throughput** | ~1K QPS per DB instance | ~100K+ QPS per Redis node |
| **DB Load** | Every read hits DB | Only cache misses hit DB |
| **Cost** | More DB replicas needed | Fewer DB instances required |
| **User Experience** | Noticeable lag | Near-instant response |

**Key principle:** Caching is most effective for **read-heavy workloads** with a high **temporal locality** (recently accessed data is likely accessed again).

### 1.3 Cache Hit, Miss, and Hit Ratio

**Cache Hit:** Requested data is found in cache. Served directly.  
**Cache Miss:** Data not in cache. Fetched from origin, then optionally stored in cache.

**Hit Ratio Formula:**

```
Hit Ratio = Cache Hits / (Cache Hits + Cache Misses)

Example: 950 hits out of 1000 requests → 950/1000 = 95% hit ratio
```

**Industry Benchmarks:**

| Hit Ratio | Assessment | Action |
|-----------|-----------|--------|
| > 95% | Excellent | Maintain current strategy |
| 85–95% | Good | Monitor and tune |
| 70–85% | Moderate | Investigate miss patterns |
| < 70% | Poor | Re-evaluate caching strategy |

**Miss Ratio** = 1 - Hit Ratio

**Effective Latency Formula:**

```
Effective Latency = (Hit Ratio × Cache Latency) + (Miss Ratio × Origin Latency)

Example: (0.95 × 2ms) + (0.05 × 100ms) = 1.9ms + 5ms = 6.9ms
vs. 100ms without cache → 14x improvement
```

> **Interview Tip:** Always quote latency numbers. Even a 5% miss rate with a 100ms origin penalty drastically improves average latency.

### 1.4 Types of Caching

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        CACHING LAYERS                                   │
│                                                                         │
│   ┌──────────┐   ┌──────────┐   ┌──────────────┐   ┌──────────────┐   │
│   │  Client   │   │   CDN    │   │  Application │   │   Database   │   │
│   │  Cache    │──►│  Cache   │──►│    Cache     │──►│    Cache     │   │
│   └──────────┘   └──────────┘   └──────────────┘   └──────────────┘   │
│                                                                         │
│   Browser        Cloudflare     Redis/Memcached     Query Cache         │
│   LocalStorage   Akamai         In-process cache    Buffer Pool         │
│   HTTP Cache     AWS CloudFront Guava/Caffeine      Materialized Views  │
│                                                                         │
│   ◄── Closest to user                    Closest to data source ──►    │
│   ◄── Lowest latency                     Highest consistency   ──►     │
│   ◄── Hardest to invalidate              Easiest to invalidate ──►     │
└─────────────────────────────────────────────────────────────────────────┘
```

| Layer | Where | Tech | Latency | Best For |
|-------|-------|------|---------|----------|
| **Client** | Browser/App | HTTP headers, LocalStorage, Service Worker | 0ms (no network) | Static assets, user preferences |
| **CDN** | Edge servers | CloudFront, Cloudflare, Akamai | 5–20ms | Images, CSS/JS, API responses |
| **Application** | App server memory or external store | Redis, Memcached, Caffeine, Guava | 1–5ms | Session data, hot queries, computed results |
| **Database** | DB engine | Buffer pool, query cache, materialized views | 5–10ms | Repeated query patterns |

> **Interview Tip:** In system design, always mention at least 2–3 caching layers. Saying "we'll add Redis" is incomplete — talk about CDN for static content, application cache for hot data, and DB-level optimizations.

---

## 2. Caching Strategies

### 2.1 Cache-Aside (Lazy Loading)

The **application** is responsible for reading and writing the cache. The cache does not interact with the data store directly.

```
                    READ PATH                          WRITE PATH
              ┌──────────────────┐              ┌──────────────────┐
              │  1. Check Cache  │              │  1. Write to DB  │
              └────────┬─────────┘              └────────┬─────────┘
                       │                                 │
              ┌────────▼─────────┐              ┌────────▼─────────┐
              │   Cache Hit?     │              │ 2. Invalidate or │
              └───┬──────────┬───┘              │    Update Cache  │
                  │          │                  └──────────────────┘
               [YES]       [NO]
                  │          │
         ┌────────▼───┐  ┌──▼──────────────┐
         │ Return     │  │ 2. Read from DB │
         │ cached data│  └────────┬────────┘
         └────────────┘          │
                        ┌────────▼────────┐
                        │ 3. Store in     │
                        │    Cache        │
                        └────────┬────────┘
                        ┌────────▼────────┐
                        │ 4. Return data  │
                        └─────────────────┘
```

**Pseudocode:**

```java
public User getUser(String userId) {
    // Step 1: Check cache
    User user = cache.get("user:" + userId);
    if (user != null) {
        return user; // Cache hit
    }
    // Step 2: Cache miss — fetch from DB
    user = database.findById(userId);
    // Step 3: Populate cache with TTL
    if (user != null) {
        cache.set("user:" + userId, user, Duration.ofMinutes(30));
    }
    return user;
}
```

**Pros:** Only requested data is cached; cache failure doesn't break reads (falls back to DB).  
**Cons:** Cache miss penalty (three round-trips); data can become stale.

---

### 2.2 Read-Through

The **cache library/provider** is responsible for loading data from the data store on a miss. Application only talks to the cache.

```
  Application ──► Cache ──[MISS]──► Cache loads from DB ──► Returns to App
                   │
                 [HIT]
                   │
                   ▼
              Return data
```

**Pros:** Simpler application code; cache-as-a-data-source abstraction.  
**Cons:** First request is always slow (cold start); cache library must know how to load data.

---

### 2.3 Write-Through

Every write goes to **cache first**, then synchronously to the **data store**.

```
  Application ──► Cache ──► Database
                   │
                   ▼
             Both updated
            before returning
```

**Pseudocode:**

```java
public void updateUser(User user) {
    cache.put("user:" + user.getId(), user);  // Write to cache
    database.save(user);                       // Write to DB synchronously
}
```

**Pros:** Cache is always consistent with DB; read-after-write consistency.  
**Cons:** Write latency increases (two writes); unused data may be cached.

---

### 2.4 Write-Back (Write-Behind)

Writes go to **cache only** initially. Cache asynchronously flushes to the data store.

```
  Application ──► Cache ──► [Async Queue] ──► Database
                   │              │
                   ▼              ▼
             Return              Batch writes
             immediately         periodically
```

**Pros:** Lowest write latency; batch writes reduce DB load; absorbs write spikes.  
**Cons:** Risk of data loss if cache crashes before flush; complex failure handling.

---

### 2.5 Write-Around

Writes go **directly to DB**, bypassing the cache. Cache is populated only on reads (misses).

**Pros:** Cache isn't polluted with data that may not be read.  
**Cons:** Read-after-write will be a cache miss.

---

### 2.6 Comparison Table

| Strategy | Read Latency | Write Latency | Consistency | Data Loss Risk | Complexity | Best Use Case |
|----------|-------------|---------------|-------------|----------------|------------|---------------|
| **Cache-Aside** | Miss penalty | N/A (app writes DB) | Eventual | None | Low | General-purpose, read-heavy |
| **Read-Through** | Miss penalty | N/A | Eventual | None | Medium | Abstracted caching layer |
| **Write-Through** | Low | High (2 writes) | Strong | None | Medium | Read-heavy with consistency needs |
| **Write-Back** | Low | Very Low | Eventual | **Yes** (crash) | High | Write-heavy, burst writes |
| **Write-Around** | Miss after write | Low | Eventual | None | Low | Write-once, read-rarely data |

> **Interview Tip:** The most common answer is cache-aside. But always discuss **why** you'd pick it. For write-heavy workloads (e.g., analytics ingestion), write-back is better. For strong consistency (e.g., financial data), write-through is safer.

**Common Combinations:**

- **Cache-aside + Write-around** — Most common in production. Reads populate cache on miss; writes go directly to DB and invalidate cache.
- **Read-through + Write-through** — Full cache abstraction. Application never talks to DB directly.
- **Read-through + Write-back** — Maximum performance. Used in systems like CPU caches, disk controllers.

---

## 3. Cache Eviction Policies

### 3.1 Overview

When cache is full and a new entry must be stored, the eviction policy determines **which existing entry to remove**.

| Policy | Logic | Time Complexity | Space Overhead | Best For |
|--------|-------|----------------|----------------|----------|
| **LRU** | Remove least recently accessed | O(1) | High (HashMap + DLL) | General-purpose, temporal locality |
| **LFU** | Remove least frequently accessed | O(1) with min-heap trick | Higher (freq counts) | Stable access patterns |
| **FIFO** | Remove oldest inserted | O(1) | Low (queue) | Simple use cases, TTL replacement |
| **Random** | Remove random entry | O(1) | None | When simplicity is paramount |
| **TTL** | Expire after time limit | O(1) amortized | Per-key timer | Session data, API responses |

### 3.2 When to Use Which

```
                         Decision Matrix
 ┌────────────────────────────────────────────────────────┐
 │ Access pattern has temporal locality?                   │
 │   YES ──► LRU                                          │
 │                                                        │
 │ Some items are always popular (power-law)?             │
 │   YES ──► LFU                                          │
 │                                                        │
 │ Data has natural expiry (sessions, tokens)?            │
 │   YES ──► TTL                                          │
 │                                                        │
 │ Need simplicity, okay with suboptimal eviction?        │
 │   YES ──► FIFO or Random                               │
 │                                                        │
 │ Mixed / unsure?                                        │
 │   ──► LRU + TTL (most production systems)              │
 └────────────────────────────────────────────────────────┘
```

### 3.3 Internal Working of LRU

LRU uses a **HashMap + Doubly Linked List** to achieve O(1) for both `get` and `put`.

```
  HashMap (O(1) lookup)
  ┌──────────┬──────────┬──────────┬──────────┐
  │  key: A  │  key: B  │  key: C  │  key: D  │
  │  val: →  │  val: →  │  val: →  │  val: →  │
  └────┬─────┴────┬─────┴────┬─────┴────┬─────┘
       │          │          │          │
       ▼          ▼          ▼          ▼
  Doubly Linked List (O(1) insertion/removal)

  HEAD ◄──► [D] ◄──► [C] ◄──► [B] ◄──► [A] ◄──► TAIL
  (MRU)                                           (LRU)

  GET(C):   Move C to head → HEAD ◄──► [C] ◄──► [D] ◄──► [B] ◄──► [A] ◄──► TAIL
  PUT(E):   Evict A (tail), add E at head → HEAD ◄──► [E] ◄──► [C] ◄──► [D] ◄──► [B] ◄──► TAIL
```

**Operations:**

| Operation | Action | Complexity |
|-----------|--------|-----------|
| `get(key)` | HashMap lookup → move node to head | O(1) |
| `put(key, val)` | If exists: update + move to head. If full: evict tail, then insert at head | O(1) |
| `evict()` | Remove tail node, delete from HashMap | O(1) |

### 3.4 LRU Cache — Java O(1) Implementation

```java
import java.util.HashMap;
import java.util.Map;

public class LRUCache<K, V> {

    private final int capacity;
    private final Map<K, Node<K, V>> map;
    private final Node<K, V> head; // sentinel
    private final Node<K, V> tail; // sentinel

    private static class Node<K, V> {
        K key;
        V value;
        Node<K, V> prev;
        Node<K, V> next;

        Node(K key, V value) {
            this.key = key;
            this.value = value;
        }
    }

    public LRUCache(int capacity) {
        this.capacity = capacity;
        this.map = new HashMap<>();
        this.head = new Node<>(null, null);
        this.tail = new Node<>(null, null);
        head.next = tail;
        tail.prev = head;
    }

    public V get(K key) {
        Node<K, V> node = map.get(key);
        if (node == null) return null;
        moveToHead(node);
        return node.value;
    }

    public void put(K key, V value) {
        Node<K, V> existing = map.get(key);
        if (existing != null) {
            existing.value = value;
            moveToHead(existing);
            return;
        }
        if (map.size() == capacity) {
            Node<K, V> lru = tail.prev;
            removeNode(lru);
            map.remove(lru.key);
        }
        Node<K, V> newNode = new Node<>(key, value);
        addToHead(newNode);
        map.put(key, newNode);
    }

    private void addToHead(Node<K, V> node) {
        node.next = head.next;
        node.prev = head;
        head.next.prev = node;
        head.next = node;
    }

    private void removeNode(Node<K, V> node) {
        node.prev.next = node.next;
        node.next.prev = node.prev;
    }

    private void moveToHead(Node<K, V> node) {
        removeNode(node);
        addToHead(node);
    }
}
```

**Why sentinel nodes?** Eliminates null checks for head/tail edge cases. Every real node always has valid `prev` and `next` pointers.

### 3.5 LFU Cache — Key Insight

LFU can be implemented in O(1) using a **HashMap + frequency-to-DLL map**:

```
  freq=1: [D] ◄──► [E]          ← evict from lowest freq, oldest entry
  freq=2: [B] ◄──► [C]
  freq=5: [A]

  Access D → move D from freq=1 to freq=2
  If freq=1 is now empty and minFreq was 1, update minFreq = 2
```

Each frequency bucket is a doubly linked list. A `minFreq` variable tracks the lowest frequency for O(1) eviction.

---

## 4. Cache Invalidation

> *"There are only two hard things in Computer Science: cache invalidation and naming things."* — Phil Karlton

### 4.1 Why Invalidation is Hard

Cache invalidation is hard because it requires **coordinating state across two independent systems** (cache and data store) that can fail independently, experience network partitions, and process concurrent operations.

**Core challenges:**

1. **No single source of truth** — Cache is a copy; any mutation to the origin must propagate.
2. **Timing** — Between the write to DB and the invalidation of cache, readers see stale data.
3. **Ordering** — Concurrent writes can cause cache to hold an older value than DB.
4. **Distributed coordination** — In a multi-node system, all caches must be invalidated.
5. **Failure modes** — If invalidation fails, stale data persists indefinitely.

### 4.2 TTL-Based Invalidation

Simplest approach: every cache entry expires after a fixed **Time-To-Live**.

```
  SET user:123 "{...}" EX 300    ← Expires in 5 minutes

  Timeline:
  T=0s     ──► Data cached
  T=0-300s ──► Served from cache (possibly stale)
  T=300s   ──► Expired, next read hits DB
```

**Choosing TTL values:**

| Data Type | Suggested TTL | Rationale |
|-----------|--------------|-----------|
| User profile | 5–15 min | Changes infrequently |
| Product catalog | 1–5 min | Updated by merchants |
| Session data | 30 min – 24 hr | Tied to session lifetime |
| Config/feature flags | 30–60 sec | Must propagate quickly |
| Stock prices | 1–5 sec | Near real-time needed |
| Static assets (CDN) | 1 hr – 1 year | Immutable with versioned URLs |

**Pros:** Simple, self-healing (stale data eventually expires).  
**Cons:** Data is stale until TTL expires; too-short TTL negates caching benefit; too-long TTL serves stale data.

### 4.3 Event-Based Invalidation

Cache is invalidated **when the source data changes**, rather than on a timer.

```
  ┌─────────────┐    write    ┌──────────┐   event    ┌──────────┐
  │ Application │───────────►│ Database │───────────►│  Cache   │
  │   Server    │            └──────────┘            │  (evict) │
  └─────────────┘                 │                  └──────────┘
                                  │ CDC / Trigger
                                  ▼
                           ┌──────────────┐
                           │ Message Bus  │
                           │(Kafka/Redis) │
                           └──────────────┘
```

**Implementation patterns:**

1. **Application-level:** After DB write, application explicitly deletes/updates cache.
2. **Database triggers:** DB fires a trigger on mutation that publishes to a message bus.
3. **Change Data Capture (CDC):** Tools like Debezium stream DB changelog to Kafka; consumers invalidate cache.
4. **Pub/Sub:** Application publishes invalidation events; cache nodes subscribe.

**Best practice: Delete, don't update.** Deleting a cache key is idempotent. Updating can cause race conditions if two concurrent writes try to set different values.

```java
// GOOD: Delete and let next read repopulate
public void updateUser(User user) {
    database.save(user);
    cache.delete("user:" + user.getId());
}

// RISKY: Race condition if two concurrent updates
public void updateUserBad(User user) {
    database.save(user);
    cache.set("user:" + user.getId(), user); // may overwrite newer value
}
```

### 4.4 Versioning Strategy

Append a version to cache keys. When data changes, increment the version.

```
  Key: product:456:v3  ──► Current version
  Key: product:456:v2  ──► Old version (will be evicted by LRU/TTL)

  On update:
    1. Write to DB
    2. Increment version in metadata store: product:456:version = 4
    3. New reads use key product:456:v4 → cache miss → fetch from DB
```

**Pros:** No explicit invalidation needed; old keys naturally expire.  
**Cons:** Requires a version counter (single point of coordination); wastes memory with old versions.

**Real-world example:** Immutable cache keys with content hashing for static assets.

```
  /static/app.js         → mutable, hard to invalidate
  /static/app.a3f8b2.js  → immutable, cache forever, deploy new hash
```

### 4.5 Handling Stale Data

**Stale-while-revalidate pattern:**

```
  1. Serve stale data immediately (user gets fast response)
  2. Asynchronously fetch fresh data in background
  3. Update cache for next request

  Timeline:
  ┌────────┐     ┌───────┐     ┌──────────┐
  │Request │────►│ Cache │────►│  Return   │ (stale but fast)
  └────────┘     │(stale)│     │  stale    │
                 └───┬───┘     └──────────┘
                     │
                     ▼ (background)
                ┌──────────┐
                │ Fetch    │
                │ fresh    │──► Update cache
                │ from DB  │
                └──────────┘
```

Used in: HTTP `stale-while-revalidate` header, SWR libraries, CDN edge caching.

**Acceptable staleness depends on the domain:**

| Domain | Acceptable Staleness |
|--------|---------------------|
| Social media feed | Seconds to minutes |
| E-commerce prices | Seconds |
| User profile photo | Minutes to hours |
| Analytics dashboard | Minutes |
| Financial transactions | **Zero** (don't cache) |

### 4.6 Strong vs Eventual Consistency

| Aspect | Strong Consistency | Eventual Consistency |
|--------|-------------------|---------------------|
| **Guarantee** | Read always returns latest write | Read may return stale data temporarily |
| **Implementation** | Write-through + synchronous invalidation | Cache-aside + TTL |
| **Latency** | Higher (must coordinate) | Lower (async) |
| **Availability** | Lower (fails if cache unavailable) | Higher (tolerates failures) |
| **Use Case** | Financial systems, inventory counts | Social feeds, product pages |

**Practical guidance:** Most systems use **eventual consistency** with short TTLs (5–60 sec). Strong consistency for cached data is expensive and often unnecessary.

> **Interview Tip:** When an interviewer asks about consistency, don't default to "strong consistency." Instead say: *"For this use case, eventual consistency with a 30-second TTL is acceptable because showing a slightly stale product description doesn't cause business harm, and it allows us to handle 100x more traffic."*

---

## 5. Distributed Caching

### 5.1 Why Local Cache is Not Enough

```
  Problem: 10 app servers, each with local cache

  Server 1: cache has user:123 = {name: "Alice"}
  Server 2: cache has user:123 = {name: "Alice"}
  ...
  Server 10: cache has user:123 = {name: "Alice"}

  User updates name to "Bob" → hits Server 3
  Server 3 updates DB + its local cache
  Servers 1,2,4-10 still have stale "Alice"
```

**Problems with local-only cache:**

1. **Inconsistency** — Each server has its own copy; updates don't propagate.
2. **Memory waste** — Same data duplicated across N servers.
3. **Cold starts** — New server has empty cache; takes time to warm up.
4. **Cache size** — Limited to single server's memory.
5. **Hit ratio** — With round-robin load balancing, requests for the same key hit different servers.

### 5.2 Distributed Cache Architecture

```
  ┌─────────────────────────────────────────────────────────────────┐
  │                     APPLICATION TIER                            │
  │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐              │
  │  │ App 1  │  │ App 2  │  │ App 3  │  │ App N  │              │
  │  └───┬────┘  └───┬────┘  └───┬────┘  └───┬────┘              │
  │      │           │           │           │                     │
  └──────┼───────────┼───────────┼───────────┼─────────────────────┘
         │           │           │           │
         ▼           ▼           ▼           ▼
  ┌─────────────────────────────────────────────────────────────────┐
  │                   DISTRIBUTED CACHE TIER                        │
  │                                                                 │
  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
  │  │ Cache    │  │ Cache    │  │ Cache    │  │ Cache    │      │
  │  │ Node 1   │  │ Node 2   │  │ Node 3   │  │ Node N   │      │
  │  │ [A-G]   │  │ [H-N]   │  │ [O-T]   │  │ [U-Z]   │      │
  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
  │                                                                 │
  │  Client library determines which node to query (hash(key))     │
  └─────────────────────────────────────────────────────────────────┘
```

### 5.3 Sharding and Partitioning

**Simple modular hashing:**

```
  node = hash(key) % N

  Problem: Adding/removing a node remaps almost ALL keys.
  N=3: hash("user:1") % 3 = 1  →  node 1
  N=4: hash("user:1") % 4 = 2  →  node 2  (REMAPPED)
```

With N nodes and a resize, ~(N-1)/N keys get remapped. For 100 nodes, **99% of keys move** — catastrophic cache miss storm.

### 5.4 Consistent Hashing

Consistent hashing maps both **keys and nodes** onto a circular hash ring. A key is assigned to the **first node clockwise** from its hash position.

```
                        0
                        │
                   ┌────┴────┐
                  /           \
           Node A               Node B
              │                    │
              │    ●key1           │
              │       ●key2       │
              │                    │
           Node D               Node C
                  \           /
                   └────┬────┘
                        │
                       180

  key1 → hash lands between A and B → assigned to Node B
  key2 → hash lands between A and B → assigned to Node B

  If Node B fails:
  key1, key2 → reassigned to Node C (next clockwise)
  Keys on A, C, D are unaffected
```

**With virtual nodes (vnodes):** Each physical node gets multiple positions on the ring to ensure even distribution.

```
  Physical Node A → vnode A1, A2, A3, A4, A5
  Physical Node B → vnode B1, B2, B3, B4, B5

  More vnodes = more uniform distribution
  Typical: 100–200 vnodes per physical node
```

**When a node is added/removed:** Only **K/N** keys are remapped (K = total keys, N = number of nodes). For 100 nodes, only ~1% of keys move.

### 5.5 Replication Strategies

| Strategy | Description | Trade-off |
|----------|-------------|-----------|
| **No replication** | Each key on one node | Fastest writes, data loss on failure |
| **Primary-replica** | Write to primary, async replicate to replicas | Read scaling, slight lag |
| **Multi-primary** | Write to any node, sync to others | High availability, conflict resolution needed |
| **Chain replication** | Write to head, propagate along chain | Strong consistency, higher latency |

**Redis Cluster replication:**

```
  Primary 1 ──► Replica 1a, Replica 1b
  Primary 2 ──► Replica 2a, Replica 2b
  Primary 3 ──► Replica 3a, Replica 3b

  Reads: Can be served by replicas (READONLY mode)
  Writes: Must go to primary
  Failover: Replica promoted to primary on failure
```

### 5.6 CAP Theorem Relevance

```
           Consistency
              /\
             /  \
            /    \
           / CP   \
          /________\
     Availability ── Partition Tolerance
                 AP
```

| Cache System | CAP Trade-off | Behavior During Partition |
|-------------|---------------|--------------------------|
| **Redis Cluster** | CP (configurable) | Rejects writes to minority partition |
| **Memcached** | AP | Each node independent, no replication |
| **Redis Sentinel** | CP | Blocks until new primary elected |
| **Hazelcast** | AP/CP configurable | Split-brain handling configurable |

> **Interview Insight:** In practice, most cache systems choose **AP** — serving stale data is better than being unavailable. Cache is not the source of truth; the database is.

---

## 6. Real-World Problems

### 6.1 Cache Stampede (Thundering Herd)

**Problem:** A popular cache key expires. Hundreds of concurrent requests simultaneously miss cache and hit the database.

```
  T=0: Key "trending_posts" expires (was serving 10K req/sec)

  ┌──────┐ ┌──────┐ ┌──────┐
  │Req 1 │ │Req 2 │ │Req N │    All 10K requests
  └──┬───┘ └──┬───┘ └──┬───┘    simultaneously miss cache
     │        │        │         and hit the database
     ▼        ▼        ▼
  ┌─────────────────────────┐
  │    DATABASE              │ ◄── 10K identical queries
  │    (overwhelmed)        │     DB CPU spikes to 100%
  └─────────────────────────┘
```

**Solutions:**

**1. Locking (Mutex)**

```java
public String getTrendingPosts() {
    String data = cache.get("trending_posts");
    if (data != null) return data;

    // Only one thread fetches; others wait
    if (lock.tryLock("trending_posts", 5, TimeUnit.SECONDS)) {
        try {
            data = cache.get("trending_posts"); // double-check
            if (data != null) return data;

            data = database.fetchTrendingPosts();
            cache.set("trending_posts", data, Duration.ofMinutes(5));
            return data;
        } finally {
            lock.unlock("trending_posts");
        }
    }
    // Failed to acquire lock — return stale data or wait
    return cache.getStale("trending_posts");
}
```

**2. Early/Probabilistic Expiry**

```
  Actual TTL: 300 seconds
  At T=270s (90% of TTL), one random request proactively refreshes the cache.
  Key never actually expires → no stampede.

  probability_of_refresh = (current_time - cache_time) / TTL
  If random() < probability_of_refresh → refresh in background
```

**3. Stale-While-Revalidate**

Return stale data immediately while refreshing asynchronously. Requires storing the actual expiry separately from the cache TTL.

### 6.2 Cache Penetration

**Problem:** Requests for data that **doesn't exist** in DB (e.g., invalid IDs, attack traffic). Every request is a cache miss → hits DB.

```
  Attacker sends: GET /user/99999999 (non-existent)
  Cache: MISS
  DB: SELECT * FROM users WHERE id=99999999 → empty
  Cache: nothing to store → next request hits DB again
  Repeat 100K times/sec → DB overwhelmed
```

**Solutions:**

**1. Cache null/empty results**

```java
User user = database.findById(id);
if (user == null) {
    cache.set("user:" + id, "NULL_MARKER", Duration.ofMinutes(2));
} else {
    cache.set("user:" + id, serialize(user), Duration.ofMinutes(30));
}
```

**2. Bloom filter**

```
  ┌──────────┐    ┌──────────────┐    ┌──────────┐    ┌──────────┐
  │ Request  │──►│ Bloom Filter │──►│  Cache   │──►│ Database │
  └──────────┘   │ "Does key    │   └──────────┘   └──────────┘
                  │  possibly    │
                  │  exist?"     │
                  └──────┬───────┘
                         │
              ┌──────────┴──────────┐
              │                     │
         "Definitely NO"      "Maybe YES"
              │                     │
              ▼                     ▼
         Return 404           Continue to cache/DB
```

Bloom filter uses ~1.2 bytes per element with 1% false positive rate. For 100M users, that's ~120MB — fits in memory.

### 6.3 Cache Avalanche

**Problem:** Large number of keys expire **at the same time**, causing a massive spike in DB load.

```
  T=0:    Populate 1M keys with TTL=3600s
  T=3600: ALL 1M keys expire simultaneously
          → 1M cache misses → DB crushed
```

**Solutions:**

**1. TTL jitter (randomization)**

```java
int baseTTL = 3600;
int jitter = random.nextInt(600); // 0-600 seconds
cache.set(key, value, Duration.ofSeconds(baseTTL + jitter));
// Keys expire between 3600-4200s, spreading the load
```

**2. Staggered warm-up:** Pre-populate cache in batches, not all at once.

**3. Circuit breaker:** If DB error rate exceeds threshold, return cached stale data or degraded response instead of hammering DB.

### 6.4 Hot Key Problem

**Problem:** A single key receives disproportionate traffic (e.g., celebrity post, flash sale item). One cache node becomes a bottleneck.

```
  "product:iphone15" → 500K req/sec → Single Redis node
  Node CPU: 100%, latency spikes, other keys on same node affected
```

**Solutions:**

| Solution | How It Works | Trade-off |
|----------|-------------|-----------|
| **Local cache (L1)** | Cache hot key in app server memory (Caffeine/Guava) | Stale data, memory per server |
| **Key replication** | Store copies as `product:iphone15:1`, `product:iphone15:2`, ... randomly read from one | More memory, invalidation complexity |
| **Read replicas** | Route reads to Redis replicas | Replication lag |
| **Key splitting** | Split value into sub-keys, aggregate on read | Application complexity |

**Production pattern: L1 + L2 approach**

```java
public Product getProduct(String id) {
    // L1: Local in-process cache (Caffeine, 1000 entries, 10s TTL)
    Product product = localCache.getIfPresent("product:" + id);
    if (product != null) return product;

    // L2: Distributed cache (Redis)
    product = redisCache.get("product:" + id);
    if (product != null) {
        localCache.put("product:" + id, product);
        return product;
    }

    // L3: Database
    product = database.findById(id);
    redisCache.set("product:" + id, product, Duration.ofMinutes(5));
    localCache.put("product:" + id, product);
    return product;
}
```

---

## 7. System Design Integration

### 7.1 Where Cache Fits in Architecture

```
┌────────────────────────────────────────────────────────────────────────────┐
│                         TYPICAL WEB ARCHITECTURE                           │
│                                                                            │
│  ┌────────┐   ┌─────────┐   ┌──────────────┐   ┌──────────┐             │
│  │Browser │──►│  CDN    │──►│Load Balancer │──►│ API GW   │             │
│  │ Cache  │   │ (Edge)  │   └──────┬───────┘   └────┬─────┘             │
│  └────────┘   └─────────┘          │                 │                    │
│                                    ▼                 ▼                    │
│                          ┌─────────────────────────────────┐             │
│                          │       APPLICATION SERVERS        │             │
│                          │  ┌──────────────────────────┐   │             │
│                          │  │   L1: In-Process Cache   │   │             │
│                          │  │   (Caffeine / Guava)     │   │             │
│                          │  └────────────┬─────────────┘   │             │
│                          └───────────────┼─────────────────┘             │
│                                          │                               │
│                                          ▼                               │
│                          ┌───────────────────────────────┐               │
│                          │   L2: Distributed Cache       │               │
│                          │   (Redis / Memcached)         │               │
│                          └───────────────┬───────────────┘               │
│                                          │                               │
│                                          ▼                               │
│                          ┌───────────────────────────────┐               │
│                          │   Database (PostgreSQL, etc)  │               │
│                          │   + Internal Buffer/Cache     │               │
│                          └───────────────────────────────┘               │
└────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Multi-Level Caching (L1, L2)

| Property | L1 (In-Process) | L2 (Distributed) |
|----------|-----------------|-------------------|
| **Location** | App server JVM heap | Separate Redis/Memcached cluster |
| **Latency** | ~100ns – 1μs | ~1–5ms (network hop) |
| **Size** | Small (100MB – 1GB) | Large (10GB – 1TB) |
| **Shared** | No (per server) | Yes (all servers) |
| **Consistency** | Weaker (local copy) | Better (single source) |
| **Failure impact** | None (just a miss) | Larger (all servers affected) |
| **Tech** | Caffeine, Guava, ConcurrentHashMap | Redis, Memcached, Hazelcast |

**Invalidation flow for multi-level:**

```
  Write to DB
       │
       ├──► Delete from L2 (Redis)
       │
       └──► Publish event to Pub/Sub
                │
                ▼
         All app servers receive event
                │
                ▼
         Each server deletes from L1 (local cache)
```

### 7.3 Design: Caching for Social Media Feed

**Scenario:** Design the caching layer for a Twitter-like home timeline.

```
  ┌──────────────────────────────────────────────────────────────────┐
  │                    FEED CACHING ARCHITECTURE                      │
  │                                                                   │
  │  User opens app                                                  │
  │       │                                                          │
  │       ▼                                                          │
  │  ┌─────────────────┐                                             │
  │  │ CDN: Cache       │ ← Static assets (profile pics, media)     │
  │  │ static content  │                                             │
  │  └────────┬────────┘                                             │
  │           │                                                      │
  │           ▼                                                      │
  │  ┌─────────────────┐                                             │
  │  │ L1: Pre-computed │ ← Top 20 posts per user (fanout-on-write) │
  │  │ timeline cache  │    Key: timeline:{userId}                   │
  │  │ (Redis sorted   │    Score: timestamp                         │
  │  │  set)           │    TTL: 5 minutes                           │
  │  └────────┬────────┘                                             │
  │           │ [MISS or pagination beyond cached]                   │
  │           ▼                                                      │
  │  ┌─────────────────┐                                             │
  │  │ L2: Individual  │ ← Post details cached separately           │
  │  │ post cache      │    Key: post:{postId}                       │
  │  │ (Redis hash)    │    TTL: 30 minutes                          │
  │  └────────┬────────┘                                             │
  │           │ [MISS]                                               │
  │           ▼                                                      │
  │  ┌─────────────────┐                                             │
  │  │ Database         │ ← Tweets table + joins                     │
  │  └─────────────────┘                                             │
  └──────────────────────────────────────────────────────────────────┘
```

**Key decisions:**

- **Fanout-on-write for non-celebrity users:** When a user posts, push post ID to all followers' timeline caches. Fast reads, slow writes.
- **Fanout-on-read for celebrities (>1M followers):** Don't pre-compute. At read time, merge celebrity posts into the timeline. Avoids writing to millions of cache entries.
- **TTL:** Short TTL (5 min) for timeline, longer for individual posts.
- **Eviction:** LRU + TTL on timeline. Keep only latest 800 posts per user.

### 7.4 Design: Caching for E-Commerce Product Page

**Scenario:** Design caching for an Amazon-like product detail page.

```
  Product Page Components:
  ┌───────────────────────────────────────┐
  │ 1. Product details (name, desc, imgs) │ ← Rarely changes → TTL 10 min
  │ 2. Price                              │ ← Changes on sales → TTL 30 sec
  │ 3. Inventory ("In Stock")             │ ← Real-time → TTL 5 sec / no cache
  │ 4. Reviews summary                    │ ← Aggregated → TTL 5 min
  │ 5. Recommendations                    │ ← ML model output → TTL 1 hour
  │ 6. Seller info                        │ ← Rarely changes → TTL 30 min
  └───────────────────────────────────────┘
```

**Cache strategy per component:**

| Component | Cache Layer | TTL | Invalidation | Strategy |
|-----------|-----------|-----|-------------|----------|
| Product details | CDN + Redis | 10 min | Event (on edit) | Read-through |
| Price | Redis | 30 sec | Event (on price change) | Cache-aside |
| Inventory | No cache or 5s TTL | 5 sec | N/A (read from DB) | Direct read or very short TTL |
| Reviews | Redis | 5 min | Periodic recompute | Cache-aside |
| Recommendations | Redis | 1 hr | Batch job refresh | Write-through |
| Seller info | Redis | 30 min | Event (on update) | Cache-aside |

**Key insight:** Don't cache the entire page as one blob. Decompose into components with different caching characteristics. This lets you cache slow-changing data aggressively while keeping fast-changing data fresh.

### 7.5 Trade-offs: Consistency vs Performance

```
  Performance ◄────────────────────────────────────────► Consistency

  High cache TTL                              No cache / cache-aside
  Write-back                                  Write-through
  Local cache only                            Strong invalidation
  Stale-while-revalidate                      Synchronous refresh
  
  Fast, scalable, stale                       Slow, correct, expensive
```

**Decision framework:**

1. **What is the cost of stale data?** (Financial loss? Bad UX? Compliance violation?)
2. **What is the read:write ratio?** (High reads → cache aggressively. High writes → careful with invalidation.)
3. **What is the access pattern?** (Zipfian/power-law → cache top items. Uniform → less benefit.)
4. **What is the acceptable staleness window?** (Seconds? Minutes? Hours?)

---

## 8. Concurrency & Performance

### 8.1 Thread Safety in Caching

In-process caches are accessed by **multiple threads concurrently**. Without thread safety, you get corrupted state, lost updates, or crashes.

| Approach | Mechanism | Performance |
|----------|----------|-------------|
| `ConcurrentHashMap` | Lock striping (segments) | High read throughput |
| `synchronized` blocks | Exclusive lock per operation | Safe but poor concurrency |
| Caffeine / Guava | Internally thread-safe | Best for production |
| Redis (external) | Single-threaded event loop | Inherently serialized |

### 8.2 Race Conditions in Cache

**Problem: Read-Modify-Write race**

```
  Thread A: reads counter = 10 from cache
  Thread B: reads counter = 10 from cache
  Thread A: sets counter = 11
  Thread B: sets counter = 11   ← Should be 12!
```

**Solution: Atomic operations**

```java
// Redis INCR is atomic
redisTemplate.opsForValue().increment("counter", 1);

// Or use CAS (Compare-And-Swap)
while (true) {
    long current = cache.get("counter");
    if (cache.compareAndSet("counter", current, current + 1)) {
        break;
    }
}
```

**Problem: Delete-then-write race (stale cache after invalidation)**

```
  Timeline:
  T1: Thread A reads from DB → gets value V1 (old)
  T2: Thread B writes V2 to DB
  T3: Thread B deletes cache key
  T4: Thread A writes V1 to cache ← STALE! Cache now has V1 instead of V2
```

**Solution: Delayed double-delete**

```java
public void updateWithDoubleDelete(String key, Object newValue) {
    cache.delete(key);                        // Delete stale cache
    database.update(key, newValue);           // Write to DB
    Thread.sleep(500);                        // Wait for in-flight reads to complete
    cache.delete(key);                        // Delete again in case of race
}
```

### 8.3 Optimistic vs Pessimistic Approaches

| Approach | Mechanism | When to Use |
|----------|----------|-------------|
| **Optimistic** | CAS, versioning, retry on conflict | Low contention, high throughput |
| **Pessimistic** | Locks, mutexes, distributed locks | High contention, must-not-fail operations |

**Optimistic (Redis WATCH/MULTI):**

```
WATCH mykey            ← Watch for changes
val = GET mykey
new_val = val + 1
MULTI                  ← Start transaction
SET mykey new_val
EXEC                   ← Fails if mykey changed since WATCH
```

**Pessimistic (Redis distributed lock via Redlock):**

```java
RLock lock = redisson.getLock("lock:product:123");
lock.lock(10, TimeUnit.SECONDS);
try {
    Product p = getFromCacheOrDB("product:123");
    p.setStock(p.getStock() - 1);
    saveToDBAndCache(p);
} finally {
    lock.unlock();
}
```

### 8.4 Connection Pooling

```java
// Redis connection pool configuration
JedisPoolConfig poolConfig = new JedisPoolConfig();
poolConfig.setMaxTotal(128);       // max connections
poolConfig.setMaxIdle(64);         // max idle connections
poolConfig.setMinIdle(16);         // min idle connections (pre-warmed)
poolConfig.setMaxWaitMillis(2000); // max wait for connection

JedisPool pool = new JedisPool(poolConfig, "redis-host", 6379);
```

**Rule of thumb:** `max_connections = (num_app_servers × threads_per_server) / num_redis_nodes` with some headroom.

---

## 9. Monitoring & Metrics

### 9.1 Key Metrics

| Metric | Formula / Source | Target | Action if Violated |
|--------|-----------------|--------|-------------------|
| **Hit Ratio** | hits / (hits + misses) | > 90% | Increase TTL, check miss patterns |
| **Latency (p50, p99)** | Histogram from client | p50 < 1ms, p99 < 5ms | Check network, key size, slow commands |
| **Eviction Rate** | `evicted_keys` counter | Low & stable | Increase memory or reduce TTL |
| **Memory Usage** | `used_memory` / `maxmemory` | < 80% | Scale out or tune eviction |
| **Connection Count** | `connected_clients` | < maxclients | Fix connection leaks, add pooling |
| **Key Count** | `dbsize` | Predictable growth | Monitor for unexpected growth |
| **Replication Lag** | `master_repl_offset` diff | < 1 sec | Check network, replica load |

### 9.2 Redis Monitoring Commands

```bash
# Real-time stats
redis-cli INFO stats

# Memory analysis
redis-cli INFO memory
redis-cli MEMORY DOCTOR

# Slow queries (commands taking > 10ms)
redis-cli SLOWLOG GET 10

# Key distribution across slots (Redis Cluster)
redis-cli CLUSTER INFO

# Monitor commands in real-time (use briefly — impacts performance)
redis-cli MONITOR

# Big key analysis
redis-cli --bigkeys
```

### 9.3 Observability Stack

```
  ┌────────────────────────────────────────────────────────────┐
  │                 CACHE OBSERVABILITY                         │
  │                                                            │
  │  ┌──────────────┐    ┌───────────────┐    ┌────────────┐  │
  │  │ Redis INFO   │──►│  Prometheus   │──►│  Grafana   │  │
  │  │ + Exporter   │   │  (metrics)    │   │ (dashboard)│  │
  │  └──────────────┘   └───────────────┘   └────────────┘  │
  │                                                            │
  │  ┌──────────────┐    ┌───────────────┐    ┌────────────┐  │
  │  │ Application  │──►│  Datadog /    │──►│  Alerts    │  │
  │  │ Metrics      │   │  New Relic    │   │            │  │
  │  │ (Micrometer) │   │              │   │            │  │
  │  └──────────────┘   └───────────────┘   └────────────┘  │
  └────────────────────────────────────────────────────────────┘
```

**Alerts to configure:**

| Alert | Condition | Severity |
|-------|-----------|----------|
| Hit ratio drop | < 80% for 5 min | Warning |
| Memory > 90% | Used > 90% of max | Critical |
| Latency spike | p99 > 10ms for 2 min | Warning |
| Eviction spike | > 1000 evictions/sec sustained | Warning |
| Connection exhaustion | Clients > 80% of maxclients | Critical |
| Replication broken | Replica disconnected > 30s | Critical |

---

## 10. Tech-Specific: Redis, Memcached, Spring Boot

### 10.1 How Redis Works Internally

**Single-threaded event loop:**

```
  ┌──────────────────────────────────────────────────────────┐
  │                 REDIS ARCHITECTURE                        │
  │                                                          │
  │  ┌─────────────┐                                        │
  │  │ Event Loop  │ ← Single thread processes all commands │
  │  │ (epoll/     │    sequentially. No locks needed.      │
  │  │  kqueue)    │                                        │
  │  └──────┬──────┘                                        │
  │         │                                                │
  │    ┌────┴────┐                                          │
  │    │ I/O     │ ← Since Redis 6.0: I/O threads handle   │
  │    │ Threads │   network read/write. Command execution  │
  │    │ (6.0+)  │   remains single-threaded.               │
  │    └─────────┘                                          │
  │                                                          │
  │  Data Structures (in-memory):                           │
  │  ┌──────────────────────────────────────────────────┐   │
  │  │ String  │ List  │ Hash  │ Set  │ Sorted Set     │   │
  │  │ (SDS)   │(ziplist│(ziplist│     │ (skiplist +    │   │
  │  │         │/quick- │/hash- │     │  hash table)   │   │
  │  │         │ list)  │ table)│     │                │   │
  │  └──────────────────────────────────────────────────┘   │
  │                                                          │
  │  Persistence:                                           │
  │  ┌──────────────────────┐  ┌────────────────────────┐   │
  │  │ RDB (Snapshots)      │  │ AOF (Append-Only File) │   │
  │  │ Fork + copy-on-write │  │ Every write logged     │   │
  │  │ Point-in-time backup │  │ Crash-safe (fsync)     │   │
  │  └──────────────────────┘  └────────────────────────┘   │
  └──────────────────────────────────────────────────────────┘
```

**Why single-threaded is fast:**
- No context switching, no lock contention.
- All operations are in-memory (~100ns per operation).
- I/O multiplexing (epoll) handles 10K+ connections.
- Bottleneck is network, not CPU.

**Redis data structures and their use cases:**

| Structure | Internal Encoding | Use Case |
|-----------|------------------|----------|
| **String** | SDS (Simple Dynamic String) | Counters, session data, simple KV |
| **List** | Quicklist (ziplist + linked list) | Message queues, activity feeds |
| **Hash** | Ziplist (small) / Hashtable (large) | Object storage (user profiles) |
| **Set** | Intset (small) / Hashtable | Tags, unique visitors, intersections |
| **Sorted Set** | Skiplist + Hashtable | Leaderboards, range queries, scheduling |
| **Stream** | Radix tree + listpacks | Event sourcing, log aggregation |
| **HyperLogLog** | Sparse/Dense registers | Cardinality estimation (unique counts) |

### 10.2 Redis vs Memcached

| Feature | Redis | Memcached |
|---------|-------|-----------|
| **Data structures** | Strings, Lists, Sets, Hashes, Sorted Sets, Streams, etc. | Strings only |
| **Persistence** | RDB snapshots + AOF | None (pure cache) |
| **Replication** | Built-in primary-replica | None (use client-side sharding) |
| **Clustering** | Redis Cluster (auto-sharding) | Client-side consistent hashing |
| **Pub/Sub** | Built-in | Not available |
| **Lua scripting** | Yes (server-side scripts) | No |
| **Transactions** | MULTI/EXEC (optimistic) | CAS only |
| **Threading** | Single-threaded + I/O threads (6.0+) | Multi-threaded |
| **Memory efficiency** | Varies by data structure | Better for simple key-value (slab allocator) |
| **Max key size** | 512MB | 250 bytes (key) / 1MB (value) |
| **Eviction policies** | 8 policies (LRU, LFU, etc.) | LRU only |
| **Typical throughput** | ~100K–300K ops/sec | ~100K–700K ops/sec (multi-threaded) |

**When to use Redis:** Need data structures, persistence, pub/sub, or clustering.  
**When to use Memcached:** Pure key-value cache, maximum simplicity, or multi-threaded performance for simple gets/sets.

> **Interview Insight:** Default to Redis in system design interviews. It covers 95% of use cases. Only mention Memcached if specifically asked or if the scenario is purely simple KV caching at extreme scale.

### 10.3 Spring Boot Caching

**Enable caching:**

```java
@SpringBootApplication
@EnableCaching
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

**Core annotations:**

```java
@Service
public class ProductService {

    // Caches the return value. Key derived from method params.
    @Cacheable(value = "products", key = "#id")
    public Product getProduct(Long id) {
        return productRepository.findById(id).orElse(null);
    }

    // Updates the cache after method execution.
    @CachePut(value = "products", key = "#product.id")
    public Product updateProduct(Product product) {
        return productRepository.save(product);
    }

    // Removes entry from cache.
    @CacheEvict(value = "products", key = "#id")
    public void deleteProduct(Long id) {
        productRepository.deleteById(id);
    }

    // Clears entire cache.
    @CacheEvict(value = "products", allEntries = true)
    public void clearProductCache() {
        // Cache is cleared after this method executes
    }

    // Conditional caching: only cache if result is not null
    @Cacheable(value = "products", key = "#id", unless = "#result == null")
    public Product getProductSafe(Long id) {
        return productRepository.findById(id).orElse(null);
    }
}
```

**Redis integration with Spring Boot:**

```yaml
# application.yml
spring:
  cache:
    type: redis
  data:
    redis:
      host: localhost
      port: 6379
      timeout: 2000ms
      lettuce:
        pool:
          max-active: 16
          max-idle: 8
          min-idle: 4
```

```java
@Configuration
@EnableCaching
public class RedisCacheConfig {

    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        RedisCacheConfiguration defaultConfig = RedisCacheConfiguration
            .defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(10))
            .serializeValuesWith(
                RedisSerializationContext.SerializationPair
                    .fromSerializer(new GenericJackson2JsonRedisSerializer())
            )
            .disableCachingNullValues();

        // Per-cache TTL configuration
        Map<String, RedisCacheConfiguration> cacheConfigs = Map.of(
            "products", defaultConfig.entryTtl(Duration.ofMinutes(30)),
            "sessions", defaultConfig.entryTtl(Duration.ofHours(1)),
            "config",   defaultConfig.entryTtl(Duration.ofSeconds(60))
        );

        return RedisCacheManager.builder(factory)
            .cacheDefaults(defaultConfig)
            .withInitialCacheConfigurations(cacheConfigs)
            .build();
    }
}
```

**Multi-level cache (Caffeine L1 + Redis L2):**

```java
@Service
public class MultiLevelCacheService {

    private final Cache<String, Product> localCache = Caffeine.newBuilder()
        .maximumSize(1_000)
        .expireAfterWrite(Duration.ofSeconds(30))
        .build();

    @Autowired
    private RedisTemplate<String, Product> redisTemplate;

    @Autowired
    private ProductRepository repository;

    public Product getProduct(Long id) {
        String key = "product:" + id;

        // L1: Local cache
        Product product = localCache.getIfPresent(key);
        if (product != null) return product;

        // L2: Redis
        product = redisTemplate.opsForValue().get(key);
        if (product != null) {
            localCache.put(key, product);
            return product;
        }

        // L3: Database
        product = repository.findById(id).orElse(null);
        if (product != null) {
            redisTemplate.opsForValue().set(key, product, Duration.ofMinutes(10));
            localCache.put(key, product);
        }
        return product;
    }
}
```

---

## 11. Interview Questions & Answers

### Q1: How would you design a caching layer for a service handling 1M requests/second?

**Answer framework:**

1. **Identify what to cache:** Profile the read patterns. Focus on the top 20% of keys that serve 80% of traffic.
2. **Multi-level cache:** L1 (Caffeine, per-server, 10s TTL, 10K entries) for ultra-hot keys. L2 (Redis Cluster, 6+ shards) for broader dataset.
3. **Consistent hashing** for key distribution across Redis nodes with ~150 virtual nodes per physical node.
4. **Cache-aside pattern** with TTL (5–10 min) + event-based invalidation via Kafka for critical data.
5. **Hot key mitigation:** Detect hot keys via metrics. Auto-replicate to local cache or add key suffixes.
6. **Stampede protection:** Probabilistic early refresh + distributed mutex for expensive computations.
7. **Monitoring:** Hit ratio > 95%, p99 < 5ms, eviction rate alerts.

---

### Q2: You notice cache hit ratio dropped from 95% to 60% overnight. How do you diagnose?

**Structured approach:**

1. **Check if the working set grew** — New features or traffic patterns may have introduced keys that don't fit in cache.
2. **Check eviction rates** — If evictions spiked, the cache is too small for the current working set.
3. **Check TTL configuration** — A deployment may have accidentally shortened TTLs.
4. **Check for cache avalanche** — Many keys expiring simultaneously (look at expiry distribution).
5. **Check for key pattern changes** — A new API or bot traffic hitting cache with low-locality keys.
6. **Check infrastructure** — Redis node went down, partition, or network issues causing failures counted as misses.

---

### Q3: How do you ensure cache consistency with the database?

**Answer:**

No single solution — it depends on the consistency requirements.

- **Eventual consistency (most common):** Cache-aside with short TTL (30s–5min). On write, delete cache key. Next read repopulates from DB. Staleness window = TTL.
- **Stronger consistency:** Write-through cache (write to cache + DB synchronously). Or event-based invalidation via CDC (Debezium) → Kafka → cache invalidation consumer.
- **For critical data (inventory, balance):** Don't cache at all, or cache with very short TTL (1–5s) and read-your-own-writes guarantee (check cache version against DB version).
- **Delayed double-delete** for race conditions in cache-aside: delete cache → write DB → wait 500ms → delete cache again.

---

### Q4: Explain cache stampede and how you'd prevent it in production.

**Answer:**

Cache stampede occurs when a popular key expires and many concurrent requests simultaneously miss cache, all hitting the database with the same expensive query.

**Prevention strategies (pick based on context):**

1. **Distributed lock:** Only one request computes; others wait or get stale data. Adds latency for waiters.
2. **Probabilistic early expiration:** Proactively refresh before actual expiry. `should_refresh = (now - cache_set_time) / TTL > random()`. No coordination needed.
3. **Stale-while-revalidate:** Serve stale value immediately, refresh async in background. Best UX but requires dual TTL tracking.
4. **Never-expire + background refresh:** Cache has no TTL. A background job periodically refreshes. Simplest, but wastes resources if data is rarely read.

In production at scale, we use a combination: probabilistic early refresh as the primary mechanism, with a distributed lock as a fallback for extremely expensive computations.

---

### Q5: Design a caching strategy for a news feed.

**Answer:**

See [Section 7.3](#73-design-caching-for-social-media-feed) for the full architecture. Key points to hit in an interview:

1. Pre-computed timeline stored in Redis Sorted Set (user_id → sorted set of post_ids by timestamp).
2. Fanout-on-write for normal users, fanout-on-read for celebrities.
3. Separate caches for timeline (short TTL) and individual posts (longer TTL).
4. CDN for media assets.
5. Cache warming on login for active users.
6. Pagination: cache first N pages, fetch rest from DB.

---

### Q6: Redis is single-threaded. How does it handle high throughput?

**Answer:**

1. **All operations are in-memory** — RAM access is ~100ns, so a single thread can process hundreds of thousands of operations per second.
2. **I/O multiplexing** — Uses epoll/kqueue to handle thousands of connections on a single thread without blocking.
3. **No locking overhead** — Single-threaded execution means no mutexes, no context switches, no cache coherence protocol overhead.
4. **Since Redis 6.0** — I/O threads handle network read/write in parallel. Command execution remains single-threaded, preserving atomicity.
5. **Efficient data structures** — Purpose-built: SDS strings, skiplists, ziplist encoding for small collections.
6. **Pipeline support** — Clients batch commands, reducing network round-trips.

Typical single-node throughput: 100K–300K ops/sec. For more, use Redis Cluster (horizontal sharding).

---

### Q7: How would you handle a flash sale where one product gets 10M views in 5 minutes?

**Answer:**

1. **CDN** for product images and static assets.
2. **L1 local cache** (Caffeine) on every app server: product details with 5s TTL. Reduces Redis load by ~90%.
3. **Redis replicas** for the product key — route reads across replicas.
4. **Inventory: don't cache**. Use Redis atomic `DECR` on a pre-loaded counter. When counter hits 0, item is sold out.
5. **Rate limiting** at API gateway to prevent abuse.
6. **Pre-warm cache** before the sale starts.
7. **Queue-based writes** — Accept orders into a Kafka queue, process asynchronously. Show "Order Processing" instead of blocking.

---

## 12. Common Mistakes & How to Answer

### Mistakes Candidates Make

| Mistake | Why It's Bad | What to Do Instead |
|---------|-------------|-------------------|
| "Just add Redis" | Ignores cache strategy, invalidation, and failure modes | Discuss strategy, TTL, invalidation, and fallback behavior |
| Caching everything | Wastes memory, complicates invalidation | Cache selectively — identify hot data and read patterns |
| Ignoring invalidation | The hardest part of caching, can't be hand-waved | Explicitly state how and when cache is invalidated |
| Single cache layer | Misses CDN, local cache, and DB cache opportunities | Discuss multi-level caching and what goes where |
| Not mentioning TTL | Implies data lives in cache forever | Always specify TTL and justify the duration |
| Ignoring failure | What if Redis is down? | Discuss graceful degradation (fall through to DB, circuit breaker) |
| Caching mutable data with long TTL | Users see stale data for too long | Match TTL to the data's mutation frequency and staleness tolerance |
| Not considering thundering herd | Cache expiry under load causes DB spikes | Mention stampede protection proactively |

### How to Answer Caching Questions in Interviews

**Step 1: Clarify requirements**

- What is the read:write ratio?
- What is the acceptable staleness?
- What is the expected QPS?
- What data are we caching?

**Step 2: State your caching strategy**

> *"I'll use cache-aside with Redis as L2 and Caffeine as L1. TTL of 5 minutes for product details. Invalidation via event-based deletion on write."*

**Step 3: Address the hard problems proactively**

- Invalidation approach
- Stampede protection
- Hot key handling
- Failure/degradation strategy

**Step 4: Discuss trade-offs**

> *"This gives us eventual consistency with a 5-minute staleness window. If we need stronger consistency, we could reduce TTL to 30 seconds or switch to write-through, but that increases Redis write load by 3x."*

**Step 5: Mention monitoring**

> *"I'd monitor hit ratio, p99 latency, and eviction rate with alerts for anomalies."*

---

## Quick Reference Card

```
┌──────────────────────────────────────────────────────────────────────┐
│                    CACHING CHEAT SHEET                                │
│                                                                      │
│  Strategy:   Cache-aside (default) │ Write-through (consistency)    │
│              Write-back (write-heavy) │ Read-through (abstraction)  │
│                                                                      │
│  Eviction:   LRU (default) │ LFU (stable patterns) │ TTL (expiry) │
│                                                                      │
│  Layers:     CDN → L1 (local) → L2 (Redis) → DB                    │
│                                                                      │
│  Invalidation: TTL + Event-based delete (not update)                │
│                                                                      │
│  Problems:   Stampede → lock/early-refresh                           │
│              Penetration → bloom filter/cache null                   │
│              Avalanche → TTL jitter                                  │
│              Hot key → L1 + replicas                                 │
│                                                                      │
│  Metrics:    Hit ratio > 90% │ p99 < 5ms │ Evictions low           │
│                                                                      │
│  Redis:      Single-threaded │ Rich data structures │ Persistent   │
│  Memcached:  Multi-threaded │ Simple KV │ Pure cache               │
│                                                                      │
│  Rule of thumb: Cache-aside + LRU + TTL + event invalidation       │
│                 covers 80% of production use cases.                  │
└──────────────────────────────────────────────────────────────────────┘
```

---

**End of Guide.**

*This guide covers the complete caching knowledge expected in FAANG system design interviews. Review the quick reference card before your interview, and practice articulating trade-offs out loud — that's what separates good answers from great ones.*
