# Redis Complete Learning Guide

> A structured beginner-to-intermediate guide for software engineers to understand, use, and master Redis.  
> **Focus:** Core concepts, internal workings, Spring Boot integration, and FAANG interview readiness.  
> **Audience:** Java/Spring Boot developers with zero prior Redis knowledge.  
> **Related:** For generic cache patterns (cache-aside, stampede, eviction theory), see [Caching — Complete FAANG Interview Guide](../Caching-Complete-Guide.md).

---

## Table of Contents

1. [What is Redis?](#1-what-is-redis)
2. [Why Redis? Real-World Problems It Solves](#2-why-redis-real-world-problems-it-solves)
3. [Redis vs Alternatives](#3-redis-vs-alternatives)
4. [Core Concepts](#4-core-concepts)
5. [Data Structures Deep Dive](#5-data-structures-deep-dive)
6. [Request Flow: End to End](#6-request-flow-end-to-end)
7. [Redis Architecture Diagrams](#7-redis-architecture-diagrams)
8. [Consistency and Durability Semantics](#8-consistency-and-durability-semantics)
9. [Practical Section: Running Redis Locally](#9-practical-section-running-redis-locally)
10. [Spring Boot Integration](#10-spring-boot-integration)
11. [Important Redis Configurations](#11-important-redis-configurations)
12. [Common Interview Questions](#12-common-interview-questions)
13. [Quick Reference Cheat Sheet](#13-quick-reference-cheat-sheet)

---

## 1. What is Redis?

### Simple Definition

**Redis** (Remote Dictionary Server) is an **open-source, in-memory data structure store** used as a database, cache, message broker, and streaming engine. In plain English, it keeps data in RAM so reads and writes complete in **sub-millisecond** time.

- Created by **Salvatore Sanfilippo** in 2009
- Written in **C**
- Used by Twitter, GitHub, Snapchat, Uber, Airbnb, and most FAANG-scale backends

### The Filing Cabinet Analogy

Think of Redis like a **giant filing cabinet in RAM**:

```
YOUR APP (Client)              REDIS (In-Memory Store)           YOUR APP (Another Service)
     │                                  │                                │
     │   "Store user:42 profile"        │                                │
     │ ────────────────────────────────>│                                │
     │                                  │   Drawer labeled "user:42"     │
     │                                  │   holds a Hash of fields       │
     │   "Get user:42"                  │                                │
     │ ────────────────────────────────>│                                │
     │<──────────────────────────────── │                                │
     │   Returns in ~1ms                │   "Get user:42"                │
     │                                  │<───────────────────────────────│
     │                                  │──────────────────────────────> │
```

**Key differences from a normal database:**
- Data lives in **memory** (with optional disk persistence)
- Supports **rich data structures** (not just rows in tables)
- Keys can **expire automatically** (TTL)
- Extremely **fast** but limited by RAM size

### One-Line Summary

> Redis = An in-memory data structure server that gives you microsecond-to-millisecond latency for caching, sessions, counters, leaderboards, pub/sub, and distributed coordination.

### Redis Is Not "Just a Cache"

| Role | Example |
|------|---------|
| **Cache** | Product catalog in front of PostgreSQL |
| **Session store** | JWT blocklist, shopping cart |
| **Primary data store** | Real-time leaderboards, rate limit counters |
| **Message broker** | Pub/Sub notifications, Redis Streams |
| **Coordination** | Distributed locks, idempotency keys |

> **Interview tip:** In interviews, say "Redis as a cache layer in front of the database" for read-heavy designs, but also mention **when Redis holds authoritative ephemeral state** (sessions, rate limits) vs **derived cache** (product details).

---

## 2. Why Redis? Real-World Problems It Solves

### Problems BEFORE Redis

```
Without Redis: Every Read Hits the Database
────────────────────────────────────────────

  10,000 users/sec ──► PostgreSQL
                         │
                         ├── Connection pool exhausted
                         ├── p99 latency spikes to 200ms+
                         └── Replica lag under load
```

```
With Redis: Hot Data Served from Memory
────────────────────────────────────────

  10,000 users/sec ──► Redis (95% hits, ~1ms)
                         │
                         └── 500 misses/sec ──► PostgreSQL
```

### Real-World Use Cases

| Use Case | Example | Why Redis? |
|----------|---------|------------|
| **Caching** | Product pages on Amazon | 100x lower latency than DB |
| **Session storage** | Logged-in user state | Fast TTL-based expiry |
| **Rate limiting** | API gateway 100 req/min | Atomic INCR, sub-ms |
| **Leaderboards** | Gaming scores, trending posts | Sorted Sets with O(log N) rank |
| **Real-time analytics** | Page view counters | Atomic counters, HyperLogLog |
| **Pub/Sub** | Live notifications | Fire-and-forget fan-out |
| **Distributed locks** | Inventory reservation | SET NX EX, Lua scripts |
| **Job queues** | Email retry queue | Lists or Streams |
| **Geospatial** | Nearby drivers (Uber) | GEO commands |
| **Idempotency** | Payment duplicate prevention | SET key NX EX |

### When Redis vs Database vs Kafka

| Need | Use | Why |
|------|-----|-----|
| **Durable business records** | PostgreSQL / MySQL | ACID, disk, complex queries |
| **Hot read path, TTL data** | Redis | Speed, simple structures |
| **Event log, replay, ordering** | Kafka | Append-only log, consumers |
| **Simple KV at extreme scale** | Memcached | Simpler, multi-threaded GET/SET |

```
Typical Microservice Stack
══════════════════════════

  Client ──► API Gateway ──► Service ──► Redis (cache/session)
                              │
                              ├──► PostgreSQL (source of truth)
                              └──► Kafka (async events to other services)
```

> **Interview insight:** Redis and Kafka **complement** each other. Redis = fast state now. Kafka = durable event history for async processing.

---

## 3. Redis vs Alternatives

### Redis vs Memcached

| Feature | Redis | Memcached |
|---------|-------|-----------|
| **Data structures** | Strings, Lists, Sets, Hashes, Sorted Sets, Streams, etc. | Strings only |
| **Persistence** | RDB + AOF | None |
| **Replication** | Built-in primary-replica | None |
| **Clustering** | Redis Cluster (16,384 hash slots) | Client-side consistent hashing |
| **Pub/Sub** | Yes | No |
| **Lua scripting** | Yes | No |
| **Threading** | Single command thread + I/O threads (6.0+) | Multi-threaded |
| **Typical use** | Cache + coordination + structures | Pure simple KV cache |

**When to use Memcached:** Only simple GET/SET at massive scale, no persistence or structures needed.

**When to use Redis:** Default for FAANG interviews — covers 95% of scenarios.

### Redis vs Local Cache (Caffeine / Guava)

| Aspect | Local (JVM) | Redis |
|--------|-------------|-------|
| **Latency** | ~0.001ms | ~0.5–2ms (network) |
| **Shared across instances** | No | Yes |
| **Survives app restart** | No | Yes (if persisted) |
| **Memory limit** | JVM heap | Dedicated Redis RAM |
| **Best for** | Ultra-hot per-instance data | Cross-pod shared state |

> **Interview tip:** Mention **multi-level cache**: Caffeine L1 (per pod) + Redis L2 (shared). See [Part 2 §11](./02-Redis-Complete-Learning-Guide-Part2.md).

### Redis vs Relational Database

| Aspect | Redis | PostgreSQL |
|--------|-------|------------|
| **Storage** | RAM (optional disk backup) | Disk-first |
| **Query model** | Key-based commands | SQL, joins, indexes |
| **Consistency** | Often eventual (replica lag) | Strong ACID transactions |
| **Cost per GB** | Higher (RAM) | Lower (SSD) |
| **Max data size** | Limited by RAM | Terabytes+ |

**Rule:** PostgreSQL owns **truth**; Redis holds **copies or ephemeral state**.

---

## 4. Core Concepts

### 4.1 Key-Value Model

**Simple Definition:**  
Every piece of data in Redis is stored under a **string key**. The value can be a String, Hash, List, Set, Sorted Set, Stream, or more.

**Naming convention (production):**

```
{service}:{entity}:{id}:{field}

Examples:
  product:catalog:12345
  session:auth:abc-def-token
  ratelimit:api:user:42:minute
```

**Rules:**
- Max key size: **512 MB** (keep keys short — under 100 bytes in practice)
- Avoid `KEYS *` in production — use `SCAN`

---

### 4.2 TTL (Time To Live)

**Simple Definition:**  
TTL automatically **deletes a key after a set duration**. Essential for caches and sessions.

```bash
SET session:user:42 "data" EX 3600    # Expires in 3600 seconds
TTL session:user:42                   # Returns remaining seconds
```

| Command | Effect |
|---------|--------|
| `EXPIRE key seconds` | Set expiry on existing key |
| `SET key value EX 60` | Set with expiry atomically |
| `PERSIST key` | Remove expiry |

> **Interview tip:** Always set TTL on cache keys. Unbounded keys = memory leak.

---

### 4.3 Single-Threaded Command Execution

**Simple Definition:**  
Redis processes **one command at a time** per shard (main thread). No locks needed inside the engine.

**Why it is still fast:**
- All data in memory (~100ns per operation)
- I/O multiplexing (epoll/kqueue) handles thousands of connections
- Bottleneck is usually **network**, not CPU

Since Redis 6.0, **I/O threads** handle network read/write; **command execution** stays single-threaded.

```
┌─────────────────────────────────────────┐
│              REDIS SERVER               │
│                                         │
│  Clients ──► I/O Threads (6.0+)         │
│                    │                    │
│                    ▼                    │
│              Event Loop                 │
│         (single command thread)         │
│                    │                    │
│                    ▼                    │
│              In-Memory Data             │
└─────────────────────────────────────────┘
```

---

### 4.4 Memory and Eviction

When `maxmemory` is reached, Redis evicts keys per policy (see [§11](#11-important-redis-configurations)).

**Memory overhead:** Redis uses ~1.5–2x raw data size due to metadata, fragmentation, and data structure encoding.

---

### 4.5 Persistence (Overview)

| Mode | How | Trade-off |
|------|-----|-----------|
| **RDB** | Point-in-time snapshots | Fast recovery; may lose last minutes |
| **AOF** | Append-only log of writes | More durable; larger files |
| **None** | Pure in-memory | Fastest; data lost on restart |

Deep dive in [Part 2 §3](./02-Redis-Complete-Learning-Guide-Part2.md).

---

### 4.6 Replication (Overview)

Primary-replica model: one **primary** handles writes; **replicas** copy data asynchronously.

```
  Writes ──► Primary ──► Replica 1
                │
                └──► Replica 2

  Reads can go to replicas (with lag awareness)
```

Cluster and Sentinel covered in Part 2.

---

### Concept Comparison Table

| Concept | Purpose | Interview keyword |
|---------|---------|-------------------|
| **Key** | Unique identifier | Naming, hot keys |
| **TTL** | Auto-expiry | Cache freshness |
| **Single-threaded** | No lock contention | Hot key bottleneck |
| **RDB/AOF** | Durability | Recovery vs performance |
| **Replica** | Read scaling, HA | Replication lag |

---

## 5. Data Structures Deep Dive

### 5.1 String

**Use cases:** Counters, JSON blobs, simple cache values, distributed locks (`SET NX`).

```bash
SET product:1 "{\"name\":\"Phone\",\"price\":999}"
GET product:1
INCR page:views:homepage
INCRBY wallet:user:42 50
SETNX lock:order:123 "owner-uuid" EX 30   # Lock only if not exists
```

| Command | Description |
|---------|-------------|
| `GET` / `SET` | Read / write |
| `INCR` / `DECR` | Atomic counter |
| `MGET` / `MSET` | Batch read / write |
| `SET key value NX EX ttl` | Distributed lock pattern |

---

### 5.2 Hash

**Use cases:** Object storage (user profile, product fields) without serializing entire JSON.

```bash
HSET user:42 name "Alice" email "alice@example.com" tier "gold"
HGET user:42 name
HGETALL user:42
HMGET user:42 name email
HINCRBY user:42 login_count 1
```

> **Interview tip:** Prefer Hash over giant JSON strings when you update individual fields — saves bandwidth on partial reads.

---

### 5.3 List

**Use cases:** Queues, recent activity feeds, timeline tails.

```bash
LPUSH notifications:user:42 "msg1" "msg2"
RPOP notifications:user:42
LRANGE notifications:user:42 0 9    # Last 10 items
BLPOP task:queue 5                  # Blocking pop (worker pattern)
```

| Pattern | Commands |
|---------|----------|
| Stack | `LPUSH` + `LPOP` |
| Queue | `LPUSH` + `RPOP` |
| Blocking worker | `BLPOP` |

---

### 5.4 Set

**Use cases:** Unique tags, mutual followers, "users who liked this post".

```bash
SADD post:99:likes user:1 user:2 user:3
SISMEMBER post:99:likes user:1
SINTER user:1:following user:2:following    # Mutual following
SCARD post:99:likes
```

---

### 5.5 Sorted Set (ZSet)

**Use cases:** Leaderboards, priority queues, time-window rate limits, delayed jobs.

```bash
ZADD leaderboard 1500 "player:alice" 1200 "player:bob"
ZRANK leaderboard "player:alice"           # 0-based rank
ZREVRANGE leaderboard 0 9 WITHSCORES       # Top 10
ZINCRBY leaderboard 50 "player:alice"
ZREMRANGEBYSCORE events:queue 0 1690000000 # Remove old entries
```

| Command | Complexity | Use |
|---------|------------|-----|
| `ZADD` | O(log N) | Add/update score |
| `ZRANGE` | O(log N + M) | Range by rank |
| `ZRANGEBYSCORE` | O(log N + M) | Range by score (time windows) |

---

### 5.6 Stream

**Use cases:** Event sourcing light, consumer groups, audit logs (see Part 2 vs Kafka).

```bash
XADD orders:stream * orderId 123 status CREATED
XREAD COUNT 10 STREAMS orders:stream 0
XGROUP CREATE orders:stream order-processors $ MKSTREAM
XREADGROUP GROUP order-processors consumer-1 COUNT 1 STREAMS orders:stream >
```

---

### 5.7 HyperLogLog, Bitmap, Geo (Brief)

| Structure | Command example | Use case |
|-----------|-----------------|----------|
| **HyperLogLog** | `PFADD uv:2024-01-01 user:42` | Approximate unique visitors |
| **Bitmap** | `SETBIT online:users 42 1` | Daily active flags |
| **Geo** | `GEOADD drivers -122.4 37.8 "driver:7"` | Nearby search |

---

### Data Structure Selection Guide

| Problem | Structure | Why |
|---------|-----------|-----|
| Cache JSON document | String or Hash | Hash for partial updates |
| API rate limit | String INCR or ZSet | ZSet for sliding window |
| Leaderboard | Sorted Set | Built-in ranking |
| Unique visitors | HyperLogLog | ~12KB fixed memory |
| Task queue | List or Stream | Stream for consumer groups |
| Pub/Sub notifications | Pub/Sub or Stream | Stream for persistence |

---

## 6. Request Flow: End to End

### 6.1 Cache-Aside Read (Most Common)

```
Client          App Service              Redis              Database
  │                 │                      │                    │
  │ GET /product/1  │                      │                    │
  │────────────────>│                      │                    │
  │                 │ GET product:1        │                    │
  │                 │─────────────────────>│                    │
  │                 │                      │                    │
  │                 │ [CACHE HIT]          │                    │
  │                 │<─────────────────────│                    │
  │<────────────────│  Return product      │                    │
  │                 │                      │                    │
  │ GET /product/99 │                      │                    │
  │────────────────>│                      │                    │
  │                 │ GET product:99       │                    │
  │                 │─────────────────────>│                    │
  │                 │ (nil) MISS           │                    │
  │                 │<─────────────────────│                    │
  │                 │ SELECT * ...         │                    │
  │                 │─────────────────────────────────────────>│
  │                 │<─────────────────────────────────────────│
  │                 │ SET product:99 EX 600                    │
  │                 │─────────────────────>│                    │
  │<────────────────│                      │                    │
```

For pattern theory (stampede, penetration), see [Caching Guide §2–6](../Caching-Complete-Guide.md).

### 6.2 Write Path with Cache Invalidation

```
  App writes to DB ──► Commit success ──► DEL cache key (or update cache)
```

**Order matters:** Invalidate **after** DB commit, or use short TTL to tolerate stale reads.

### 6.3 Pipelining

Send multiple commands without waiting for each response — reduces round trips.

```bash
# redis-cli
MULTI
GET product:1
GET product:2
GET product:3
EXEC
```

In Java (Lettuce): use `redisTemplate.executePipelined(...)`.

---

## 7. Redis Architecture Diagrams

### 7.1 Standalone (Development)

```
┌──────────────┐         ┌──────────────┐
│ Spring Boot  │────────>│ Redis        │
│ App          │  TCP    │ :6379        │
└──────────────┘         └──────────────┘
```

### 7.2 Primary-Replica (Production Baseline)

```
                    ┌─────────────┐
         Writes ───>│   Primary   │
                    └──────┬──────┘
                           │ replication
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │ Replica 1│ │ Replica 2│ │ Replica 3│
        └──────────┘ └──────────┘ └──────────┘
              ▲
              └── Reads (optional, eventual)
```

### 7.3 Redis Cluster (Preview)

```
┌─────────────────────────────────────────────────────────┐
│                   Redis Cluster                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │
│  │ Master 1    │  │ Master 2    │  │ Master 3    │      │
│  │ Slots 0-5460│  │ 5461-10922  │  │ 10923-16383 │      │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘      │
│         │ Replica        │ Replica        │ Replica     │
└─────────────────────────────────────────────────────────┘
```

Full cluster details: [Part 2 §1](./02-Redis-Complete-Learning-Guide-Part2.md).

---

## 8. Consistency and Durability Semantics

### 8.1 Cache vs Source of Truth

| Data type | Source of truth | Redis role |
|-----------|-----------------|------------|
| Product catalog | PostgreSQL | Cache (can rebuild) |
| Shopping cart | Redis or DB | Often Redis with persistence |
| Rate limit counter | Redis | Authoritative for window |
| Session token | Redis | Authoritative until expiry |

### 8.2 Consistency Levels

| Setup | Behavior |
|-------|----------|
| Single Redis instance | Strong consistency for commands |
| Primary + replicas | Writes to primary; reads from replica may be **stale** |
| Redis Cluster | Per-key linearizable writes to slot owner |

> **Interview tip:** Say "Redis cache is **eventually consistent** with the database. We accept stale reads within TTL or invalidate on write."

### 8.3 Durability vs Performance

| `appendfsync` | Durability | Speed |
|---------------|------------|-------|
| `always` | Highest | Slowest |
| `everysec` | Good (≤1 sec loss) | Balanced (default) |
| `no` | Lowest | Fastest |

### 8.4 Link to Cache Patterns

| Pattern | Summary | Deep dive |
|---------|---------|-----------|
| Cache-aside | App manages cache | [Caching Guide §2](../Caching-Complete-Guide.md) |
| Write-through | Write to cache + DB together | Caching Guide |
| Write-behind | Write cache, async DB | Caching Guide |
| Read-through | Cache loads on miss | Caching Guide |

---

## 9. Practical Section: Running Redis Locally

### 9.1 Docker Compose (Recommended)

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  redis:
    image: redis:7.2-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data

  redis-insight:
    image: redis/redisinsight:latest
    ports:
      - "5540:5540"
    depends_on:
      - redis

volumes:
  redis-data:
```

Start:

```bash
docker compose up -d
docker compose ps
```

### 9.2 redis-cli Cookbook

```bash
# Connect
docker exec -it <redis-container> redis-cli

# Basic KV
SET greeting "Hello Redis"
GET greeting

# TTL
SET temp:key "value" EX 60
TTL temp:key

# Hash
HSET user:1 name "Bob" age 30
HGETALL user:1

# Counter
INCR visits:today

# Sorted set leaderboard
ZADD game:lb 100 "player1" 200 "player2"
ZREVRANGE game:lb 0 2 WITHSCORES

# Pipeline
echo -e "SET a 1\nINCR a\nGET a" | redis-cli --pipe

# Server info
INFO memory
INFO stats
DBSIZE
```

### 9.3 Spring Boot Local Profile

```yaml
spring:
  data:
    redis:
      host: localhost
      port: 6379
```

---

## 10. Spring Boot Integration

### 10.1 Dependencies

```xml
<dependencies>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-data-redis</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-cache</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-actuator</artifactId>
    </dependency>
</dependencies>
```

### 10.2 application.yml

```yaml
spring:
  application:
    name: redis-demo

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
          max-wait: 2000ms

management:
  endpoints:
    web:
      exposure:
        include: health,info
  health:
    redis:
      enabled: true
```

### 10.3 Redis Configuration Beans

```java
@Configuration
@EnableCaching
public class RedisConfig {

    @Bean
    public RedisTemplate<String, Object> redisTemplate(
            RedisConnectionFactory connectionFactory) {
        RedisTemplate<String, Object> template = new RedisTemplate<>();
        template.setConnectionFactory(connectionFactory);
        template.setKeySerializer(new StringRedisSerializer());

        ObjectMapper mapper = new ObjectMapper();
        mapper.registerModule(new JavaTimeModule());
        GenericJackson2JsonRedisSerializer jsonSerializer =
            new GenericJackson2JsonRedisSerializer(mapper);

        template.setValueSerializer(jsonSerializer);
        template.setHashKeySerializer(new StringRedisSerializer());
        template.setHashValueSerializer(jsonSerializer);
        template.afterPropertiesSet();
        return template;
    }

    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        ObjectMapper mapper = new ObjectMapper();
        mapper.registerModule(new JavaTimeModule());
        GenericJackson2JsonRedisSerializer serializer =
            new GenericJackson2JsonRedisSerializer(mapper);

        RedisCacheConfiguration defaults = RedisCacheConfiguration
            .defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(10))
            .serializeValuesWith(
                RedisSerializationContext.SerializationPair.fromSerializer(serializer))
            .disableCachingNullValues();

        Map<String, RedisCacheConfiguration> perCache = Map.of(
            "products", defaults.entryTtl(Duration.ofMinutes(30)),
            "sessions", defaults.entryTtl(Duration.ofHours(2))
        );

        return RedisCacheManager.builder(factory)
            .cacheDefaults(defaults)
            .withInitialCacheConfigurations(perCache)
            .build();
    }
}
```

### 10.4 Declarative Caching with Annotations

```java
@Service
@Slf4j
public class ProductService {

    private final ProductRepository productRepository;

    public ProductService(ProductRepository productRepository) {
        this.productRepository = productRepository;
    }

    @Cacheable(value = "products", key = "#id", unless = "#result == null")
    public Product getById(Long id) {
        log.info("Cache MISS — loading product {} from database", id);
        return productRepository.findById(id).orElse(null);
    }

    @CachePut(value = "products", key = "#product.id")
    public Product update(Product product) {
        return productRepository.save(product);
    }

    @CacheEvict(value = "products", key = "#id")
    public void delete(Long id) {
        productRepository.deleteById(id);
    }
}
```

### 10.5 Manual Cache-Aside (Interview-Friendly)

```java
@Service
@Slf4j
@RequiredArgsConstructor
public class ProductCacheService {

    private static final String KEY_PREFIX = "product:";
    private static final Duration TTL = Duration.ofMinutes(10);

    private final RedisTemplate<String, Object> redisTemplate;
    private final ProductRepository productRepository;

    public Optional<Product> getProduct(Long id) {
        String key = KEY_PREFIX + id;

        try {
            Product cached = (Product) redisTemplate.opsForValue().get(key);
            if (cached != null) {
                log.info("Cache HIT for key={}", key);
                return Optional.of(cached);
            }
        } catch (DataAccessException ex) {
            log.warn("Redis unavailable, falling back to DB for id={}", id, ex);
            return productRepository.findById(id);
        }

        log.info("Cache MISS for key={}", key);
        Optional<Product> fromDb = productRepository.findById(id);
        fromDb.ifPresent(product -> {
            try {
                redisTemplate.opsForValue().set(key, product, TTL);
            } catch (DataAccessException ex) {
                log.warn("Failed to write cache for key={}", key, ex);
            }
        });
        return fromDb;
    }

    public void evict(Long id) {
        redisTemplate.delete(KEY_PREFIX + id);
    }
}
```

### 10.6 StringRedisTemplate for Counters

```java
@Service
@RequiredArgsConstructor
public class PageViewService {

    private final StringRedisTemplate redisTemplate;

    public long incrementPageView(String pageId) {
        String key = "pageviews:" + pageId;
        Long count = redisTemplate.opsForValue().increment(key);
        redisTemplate.expire(key, Duration.ofDays(1));
        return count != null ? count : 0;
    }
}
```

### 10.7 REST Controller to Test

```java
@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
public class ProductController {

    private final ProductCacheService productCacheService;

    @GetMapping("/{id}")
    public ResponseEntity<Product> getProduct(@PathVariable Long id) {
        return productCacheService.getProduct(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}/cache")
    public ResponseEntity<Void> evictCache(@PathVariable Long id) {
        productCacheService.evict(id);
        return ResponseEntity.noContent().build();
    }
}
```

### 10.8 Health Check

With Actuator enabled, `GET /actuator/health` includes Redis status when the server is reachable.

```json
{
  "status": "UP",
  "components": {
    "redis": { "status": "UP" },
    "db": { "status": "UP" }
  }
}
```

### 10.9 Key Naming and Error Handling Checklist

| Practice | Reason |
|----------|--------|
| Prefix keys by domain | Avoid collisions in shared Redis |
| Set TTL on every cache key | Prevent memory leaks |
| Catch `DataAccessException` | Degrade to DB if Redis is down |
| Do not cache nulls by default | Use `unless` or Bloom filter for misses |
| Use connection pooling (Lettuce) | Avoid connection exhaustion |

---

## 11. Important Redis Configurations

### 11.1 Memory

```conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

| Policy | Evicts | Best for |
|--------|--------|----------|
| `allkeys-lru` | Any key, LRU | General cache |
| `volatile-lru` | Keys with TTL only | Mixed TTL data |
| `allkeys-lfu` | Any key, LFU | Hot key skew (Redis 4.0+) |
| `noeviction` | Nothing — writes fail | When eviction is unacceptable |

### 11.2 Networking

```conf
timeout 300
tcp-keepalive 60
maxclients 10000
```

### 11.3 Persistence (Quick Reference)

```conf
save 900 1
save 300 10
appendonly yes
appendfsync everysec
```

### 11.4 Lettuce Pool (Spring Boot)

| Setting | Guidance |
|---------|----------|
| `max-active` | ≥ expected concurrent Redis users |
| `max-idle` | Keep warm connections |
| `max-wait` | Fail fast if pool exhausted |

> **Interview tip:** If p99 latency spikes under load, check **connection pool exhaustion** and **hot keys** before blaming Redis itself.

---

## 12. Common Interview Questions

### Q1: What is Redis and why is it used?

**Answer:**  
Redis is an in-memory data structure store used as a cache, session store, message broker, and coordination layer. It delivers sub-millisecond latency because data lives in RAM and supports rich structures (Hashes, Sorted Sets, etc.). Teams use it to reduce database load, speed up APIs, and implement rate limiting, leaderboards, and distributed locks.

### Q2: Why is Redis single-threaded and still fast?

**Answer:**  
Commands run on one thread per shard, avoiding lock contention. Operations are in-memory (nanoseconds). The bottleneck is usually network I/O, which Redis 6+ improves with I/O threads. Throughput is often 100K–300K ops/sec per node.

### Q3: What happens when Redis runs out of memory?

**Answer:**  
If `maxmemory` is set with an eviction policy, Redis evicts keys (e.g., LRU). With `noeviction`, write commands fail with OOM errors. Reads may still work. Production always sets `maxmemory` and a policy.

### Q4: RDB vs AOF — which do you choose?

**Answer:**  
RDB: periodic snapshots, compact, fast restarts, may lose recent data. AOF: logs every write, more durable, larger files, rewrite compaction. Production often uses **both**: RDB for fast baseline + AOF for durability. Cache-only instances may disable persistence.

### Q5: How do you handle cache invalidation?

**Answer:**  
(1) TTL for time-bound freshness, (2) explicit `DEL` on DB write, (3) event-driven invalidation via Kafka, (4) versioned keys. Order: commit DB first, then invalidate cache. See [Caching Guide §4](../Caching-Complete-Guide.md).

### Q6: What is a hot key problem?

**Answer:**  
One key (e.g., viral product) gets disproportionate traffic, saturating a single Redis thread/shard. Mitigations: local L1 cache, read replicas, split key into shards (`product:1:partA`), or pre-warm CDN.

### Q7: Redis vs Memcached?

**Answer:**  
Default to Redis: structures, persistence, replication, pub/sub. Memcached only for simple KV at extreme scale with no extra features.

### Q8: Can Redis replace a database?

**Answer:**  
Rarely for core business data. Redis lacks complex queries, joins, and cheap durable storage at TB scale. Use it for cache, sessions, counters, and ephemeral state; PostgreSQL remains source of truth.

### Q9: How does Redis replication work?

**Answer:**  
Primary accepts writes; replicas are async copies. On primary failure, Sentinel or manual promotion elects a new primary. Reads from replicas may lag — not suitable for read-your-writes without routing writes to primary.

### Q10: Explain SET NX EX for distributed locks.

**Answer:**  
`SET lock:resource uuid NX EX 30` sets the key only if absent, with expiry. Prevents deadlocks if the holder crashes. Must verify lock value on unlock (Lua script) and use fencing tokens for external systems. See Part 2 §5.

### Q11: What is cache stampede and how does Redis help?

**Answer:**  
Many requests miss cache simultaneously and hammer the DB. Fixes: mutex with `SETNX`, single-flight pattern, stale-while-revalidate, probabilistic early expiry. See [Caching Guide §6](../Caching-Complete-Guide.md) and Part 2 §11.

### Q12: Redis vs Kafka?

**Answer:**  
Redis: low-latency state, TTL, key-based access, not a durable event log. Kafka: durable ordered log, replay, high-throughput streaming. Use together: Redis for hot state, Kafka for async domain events.

---

## 13. Quick Reference Cheat Sheet

### Terminology

| Term | One-Line Definition |
|------|---------------------|
| **Key** | Unique string identifier for a value |
| **TTL** | Time until key auto-deletes |
| **Primary** | Node that accepts writes |
| **Replica** | Read copy of primary data |
| **RDB** | Snapshot persistence |
| **AOF** | Append-only write log |
| **Eviction** | Removing keys when memory is full |
| **Pipeline** | Batch commands, one round trip |
| **Cluster** | Sharded multi-master deployment |
| **Sentinel** | HA failover for primary-replica |
| **Slot** | Hash range in cluster (16,384 total) |

### CLI Cheat Sheet

```bash
SET key value [EX seconds]
GET key
DEL key
EXISTS key
TTL key
INCR key
HSET key field value
HGET key field
LPUSH key value
RPOP key
SADD key member
ZADD key score member
EXPIRE key seconds
INFO section
SCAN cursor MATCH pattern COUNT 100
```

### Spring Boot Quick Template

```
pom.xml → spring-boot-starter-data-redis, spring-boot-starter-cache
application.yml → host, port, lettuce pool
RedisConfig.java → RedisTemplate + RedisCacheManager
@Service → @Cacheable OR manual cache-aside
Controller → REST endpoints
actuator → /actuator/health
```

---

> **Continue to [Part 2 (Advanced)](./02-Redis-Complete-Learning-Guide-Part2.md):**
> - Redis Cluster and Sentinel
> - Persistence, replication, and failure scenarios
> - Distributed locks, Pub/Sub vs Streams
> - Performance tuning and security
> - Advanced Spring Boot patterns (Redisson, Streams, Testcontainers)
> - Real-world examples: rate limiter, sessions, leaderboards, payments
> - Advanced interview questions and production checklist
