# Redis Complete Learning Guide — Part 2 (Advanced)

> A structured intermediate-to-advanced guide continuing from Part 1.  
> **Focus:** Redis Cluster, Sentinel, persistence, distributed patterns, production tuning, advanced Spring Boot, real-world examples, and FAANG interview depth.  
> **Prerequisite:** Complete [Part 1](./01-Redis-Complete-Learning-Guide.md) (core concepts, data structures, Spring Boot basics).

---

## Table of Contents

1. [Redis Cluster](#1-redis-cluster)
2. [Redis Sentinel](#2-redis-sentinel)
3. [Persistence Deep Dive](#3-persistence-deep-dive)
4. [Replication and Read Scaling](#4-replication-and-read-scaling)
5. [Distributed Locking](#5-distributed-locking)
6. [Pub/Sub vs Streams](#6-pubsub-vs-streams)
7. [Performance Tuning and Best Practices](#7-performance-tuning-and-best-practices)
8. [Security](#8-security)
9. [Monitoring and Observability](#9-monitoring-and-observability)
10. [Redis Design Patterns in Microservices](#10-redis-design-patterns-in-microservices)
11. [Advanced Spring Boot Patterns](#11-advanced-spring-boot-patterns)
12. [Real-World Examples](#12-real-world-examples)
13. [Advanced Interview Questions](#13-advanced-interview-questions)
14. [Quick Reference Cheat Sheet — Part 2](#14-quick-reference-cheat-sheet--part-2)

---

## 1. Redis Cluster

### 1.1 What is Redis Cluster?

**Simple Definition:**  
Redis Cluster is a **distributed deployment** that automatically shards data across multiple masters using **hash slots**, with optional replicas per master for high availability.

**Real-World Analogy:**  
A library with **16,384 shelves** (slots). Each book's ISBN is hashed to a shelf number. Multiple buildings (nodes) each hold a range of shelves. If one building closes, a backup building has copies of those shelves.

### 1.2 Hash Slots

- Total slots: **16,384** (0–16383)
- Key routing: `CRC16(key) mod 16384` → slot → responsible master
- Hash tags: `user:{42}:profile` and `user:{42}:orders` share slot `{42}`

```
Key: "product:12345"
        │
        ▼
   CRC16 mod 16384 = slot 8912
        │
        ▼
   Master node owning slots 5461–10922
```

### 1.3 Cluster Topology

```
┌────────────────────────────────────────────────────────────────┐
│                     Redis Cluster (6 nodes)                     │
│                                                                │
│  Master A          Master B          Master C                  │
│  Slots 0-5460      5461-10922        10923-16383               │
│      │                 │                  │                    │
│  Replica A'        Replica B'         Replica C'               │
└────────────────────────────────────────────────────────────────┘
```

### 1.4 Failover

1. Master fails → replicas detect via gossip protocol
2. Replica promoted to master for that slot range
3. Cluster reconfigures; clients fetch updated topology via `CLUSTER SLOTS`

### 1.5 CAP in Cluster

| Aspect | Behavior |
|--------|----------|
| **Partition tolerance** | Yes — cluster continues with majority |
| **Availability** | Minority partition cannot serve writes |
| **Consistency** | Per-key linearizable writes to slot owner |

> **Interview tip:** Redis Cluster is **CP** for writes during partitions — minority side rejects writes. For cache workloads, that is usually acceptable.

### 1.6 Useful Commands

```bash
CLUSTER INFO
CLUSTER NODES
CLUSTER SLOTS
CLUSTER KEYSLOT mykey
redis-cli -c -p 7000   # -c enables cluster redirect
```

### 1.7 When NOT to Use Cluster

- Dataset fits one master with replicas (simpler Sentinel setup)
- Need multi-key transactions across unrelated keys (use hash tags or redesign)
- Very small deployments where ops complexity outweighs benefit

For full "design Redis-like system" interviews, see [009-distributed-cache-redis-like.md](../real-world-examples/009-distributed-cache-redis-like.md).

---

## 2. Redis Sentinel

### 2.1 What is Sentinel?

**Simple Definition:**  
Sentinel is a **high-availability solution** for Redis primary-replica setups. It monitors nodes, detects failures, and **automatically promotes** a replica to primary.

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ Sentinel 1  │  │ Sentinel 2  │  │ Sentinel 3  │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │
       └────────────────┼────────────────┘
                        │ monitor
                        ▼
                 ┌─────────────┐
                 │   Primary   │
                 └──────┬──────┘
                        │
            ┌───────────┴───────────┐
            ▼                       ▼
     ┌─────────────┐         ┌─────────────┐
     │  Replica 1  │         │  Replica 2  │
     └─────────────┘         └─────────────┘
```

### 2.2 Failover Steps

1. Sentinels agree primary is down (quorum)
2. Select best replica (priority, offset, runid)
3. `REPLICAOF NO ONE` on chosen replica
4. Reconfigure other replicas and notify clients

### 2.3 Sentinel vs Cluster

| Feature | Sentinel | Cluster |
|---------|----------|---------|
| **Sharding** | No (single primary) | Yes (16K slots) |
| **Failover** | Yes | Yes |
| **Max dataset** | One node RAM | Sum of all masters |
| **Complexity** | Lower | Higher |
| **Use when** | < RAM of one machine | Need horizontal scale |

### 2.4 Spring Boot with Sentinel

```yaml
spring:
  data:
    redis:
      sentinel:
        master: mymaster
        nodes:
          - sentinel1:26379
          - sentinel2:26379
          - sentinel3:26379
```

---

## 3. Persistence Deep Dive

### 3.1 RDB (Redis Database)

**How it works:**
1. `fork()` creates child process
2. Child writes snapshot to disk
3. Parent continues serving (copy-on-write for changed pages)

```
Timeline:
  t0 ── fork ──► Child writes dump.rdb
  t1 ── Parent serves writes (COW pages)
  t2 ── Child completes snapshot
```

| Pros | Cons |
|------|------|
| Compact file | Point-in-time only |
| Fast restart | fork() latency spike on large RAM |
| Good for backups | May lose data since last snapshot |

```conf
save 900 1      # save if 1 key changed in 900 sec
save 300 10
save 60 10000
dbfilename dump.rdb
```

### 3.2 AOF (Append Only File)

Every write appended to log. Replay on startup reconstructs dataset.

| `appendfsync` | Behavior |
|---------------|----------|
| `always` | fsync every write — safest, slowest |
| `everysec` | fsync every second — default, ≤1s loss |
| `no` | OS decides — fastest, least safe |

**AOF rewrite:** Compacts log by summarizing current key state (like garbage collection).

```conf
appendonly yes
appendfilename "appendonly.aof"
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
```

### 3.3 Production Recommendation

| Workload | Recommendation |
|----------|----------------|
| Pure cache | No persistence or RDB only |
| Sessions | AOF everysec + RDB backup |
| Critical counters | AOF always (accept latency cost) |

### 3.4 Disaster Recovery

1. Regular RDB snapshots to S3
2. AOF rewrite monitoring
3. Test restore procedure quarterly
4. Document RTO/RPO: "We accept up to 1 second of writes"

---

## 4. Replication and Read Scaling

### 4.1 Async Replication

Primary accepts write → returns to client → replicates to replicas **asynchronously**.

```
Write ──► Primary ──► ACK to client
              │
              └── (async) ──► Replica (may lag 1-50ms+)
```

### 4.2 Replication Lag Problems

| Problem | Scenario | Fix |
|---------|----------|-----|
| **Stale read** | Read replica after write | Read from primary for critical reads |
| **Lost write** | Primary dies before replicate | Use `WAIT` or min replicas (trade latency) |
| **Split brain** | Network partition | Sentinel quorum, odd number of sentinels |

### 4.3 WAIT Command

```bash
SET key value
WAIT 1 1000   # Wait until 1 replica acknowledges, timeout 1000ms
```

Use when you need stronger durability without full sync replication.

### 4.4 Read Replica Routing

```
                    ┌── Read (stale OK) ──► Replica
  App ── Write ──► Primary
                    └── Read (must be fresh) ──► Primary
```

> **Interview tip:** Never assume read-your-writes from replicas. Route session validation and payment status to primary.

---

## 5. Distributed Locking

### 5.1 Basic Lock with SET NX EX

```bash
SET lock:order:123 "uuid-owner-abc" NX EX 30
```

| Flag | Meaning |
|------|---------|
| `NX` | Set only if not exists |
| `EX 30` | Auto-expire in 30 seconds |

### 5.2 Safe Unlock (Lua Script)

Never `DEL` without verifying owner — another process may have acquired the lock.

```lua
-- unlock.lua
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end
```

### 5.3 Redlock (Conceptual)

Acquire lock on **majority** of N independent Redis nodes with same TTL. Debated in industry (Martin Kleppmann critique) — understand trade-offs for interviews.

### 5.4 Fencing Tokens

```
1. Client acquires lock, gets token=34
2. Client slow, lock expires
3. Client B acquires lock, token=35
4. Client A writes to DB with token=34
5. Storage rejects token < 35
```

### 5.5 When NOT to Use Redis Locks

- Long-held locks (minutes) — use DB or dedicated coordinator
- Cross-datacenter strong consistency — etcd/ZooKeeper may fit better
- Locks without fencing on external systems

---

## 6. Pub/Sub vs Streams

### 6.1 Pub/Sub

**Fire-and-forget:** Subscribers offline = messages lost.

```
Publisher ──► Channel "notifications" ──► Subscriber A
                                      └──► Subscriber B
```

```bash
PUBLISH notifications "Hello"
SUBSCRIBE notifications
```

### 6.2 Redis Streams

**Persistent log** with consumer groups, ACKs, and replay.

```
Producer ──► Stream "orders" ──► Consumer Group "fulfillment"
                           └──► Consumer Group "analytics"
```

| Feature | Pub/Sub | Streams | Kafka |
|---------|---------|---------|-------|
| Persistence | No | Yes (bounded) | Yes |
| Consumer groups | No | Yes | Yes |
| Replay | No | Yes | Yes |
| Throughput | High | High | Very high |
| Ops complexity | Low | Low | High |

> **Interview tip:** Use Pub/Sub for live ephemeral events (typing indicators). Use Streams for lightweight queuing. Use Kafka for durable cross-service event backbone.

---

## 7. Performance Tuning and Best Practices

### 7.1 Pipelining and Batching

```java
List<Object> results = redisTemplate.executePipelined((RedisCallback<Object>) connection -> {
    for (Long id : productIds) {
        connection.stringCommands().get(("product:" + id).getBytes());
    }
    return null;
});
```

### 7.2 Avoid KEYS *

Use `SCAN` for iteration:

```bash
SCAN 0 MATCH product:* COUNT 100
```

### 7.3 Hot Key Mitigation

| Technique | How |
|-----------|-----|
| **Local L1 cache** | Caffeine in app for top 100 keys |
| **Key splitting** | `product:1:copy0` … `product:1:copy9`, random read |
| **Read replicas** | Spread read load |
| **Pre-warm** | Load before traffic spike |

### 7.4 Big Key Problem

Large Hash/List/Set blocks the single thread during deletion.

| Prevention | Fix |
|------------|-----|
| Chunk data | `UNLINK` (async delete, Redis 4.0+) |
| Monitor `MEMORY USAGE key` | Split into smaller keys |

### 7.5 Connection Pool Sizing

```
connections ≈ (cores × 2) + effective_spindle_count
```

For Spring Boot: monitor pool wait time; increase `max-active` if threads block.

### 7.6 Memory Optimization

- Use Hash instead of JSON string for objects
- Set appropriate TTL
- Enable `activerehashing yes`
- Monitor `mem_fragmentation_ratio` (ideal 1.0–1.5)

---

## 8. Security

### 8.1 Network

- Redis **never** exposed to public internet
- VPC private subnets, security groups
- TLS in transit (Redis 6+ ACL + TLS)

### 8.2 ACL (Redis 6+)

```bash
ACL SETUSER appuser on >strongpassword ~product:* +get +set +del
ACL LIST
```

### 8.3 Spring Boot TLS

```yaml
spring:
  data:
    redis:
      ssl:
        enabled: true
      password: ${REDIS_PASSWORD}
```

### 8.4 Checklist

| Control | Purpose |
|---------|---------|
| `requirepass` / ACL | Authentication |
| TLS | Encryption in transit |
| Rename dangerous commands | Disable `FLUSHALL` in prod |
| Separate instances per env | No dev → prod sharing |

---

## 9. Monitoring and Observability

### 9.1 INFO Commands

```bash
INFO memory          # used_memory, fragmentation
INFO stats           # ops/sec, hits/misses
INFO replication     # lag, role
INFO clients         # connected_clients
LATENCY DOCTOR
SLOWLOG GET 10
```

### 9.2 Key Metrics

| Metric | Alert threshold |
|--------|-----------------|
| `used_memory` / `maxmemory` | > 80% |
| `connected_clients` | > 80% of maxclients |
| `replication lag` | > 10s |
| `evicted_keys` rate | Sudden spike |
| `instantaneous_ops_per_sec` | Baseline anomaly |
| Hit ratio | < 85% for cache |

### 9.3 Prometheus + Grafana

Use `redis_exporter` sidecar:

```
Redis ──► redis_exporter:9121 ──► Prometheus ──► Grafana dashboards
```

### 9.4 Spring Boot Metrics

```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-registry-prometheus</artifactId>
</dependency>
```

Micrometer exposes Lettuce connection pool and command latency when configured.

---

## 10. Redis Design Patterns in Microservices

| Pattern | Redis feature | Example |
|---------|---------------|---------|
| **Rate limiter** | INCR + EXPIRE or ZSet | API gateway |
| **Session store** | Hash + TTL | Auth service |
| **Idempotency** | SET NX EX | Payment API |
| **Leaderboard** | Sorted Set | Gaming |
| **Delayed jobs** | ZSet score = run timestamp | Retry queue |
| **Bloom filter** | RedisBloom module | Cache penetration |
| **Feature flags** | String/Hash | Toggle per user |
| **Distributed lock** | SET NX + Lua | Inventory |

```
Microservice A ──┐
Microservice B ──┼──► Redis Cluster ◄── shared patterns
Microservice C ──┘         │
                             ├── Cache
                             ├── Sessions
                             ├── Rate limits
                             └── Locks
```

---

## 11. Advanced Spring Boot Patterns

### 11.1 Redisson Distributed Lock

```xml
<dependency>
    <groupId>org.redisson</groupId>
    <artifactId>redisson-spring-boot-starter</artifactId>
    <version>3.27.0</version>
</dependency>
```

```java
@Service
@RequiredArgsConstructor
public class InventoryService {

    private final RedissonClient redisson;

    public void reserveStock(String productId, int quantity) {
        RLock lock = redisson.getLock("lock:inventory:" + productId);
        try {
            if (lock.tryLock(5, 30, TimeUnit.SECONDS)) {
                // decrement stock in DB
            } else {
                throw new IllegalStateException("Could not acquire lock");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
            }
        }
    }
}
```

### 11.2 Redisson Rate Limiter

```java
RRateLimiter limiter = redisson.getRateLimiter("ratelimit:user:" + userId);
limiter.trySetRate(RateType.OVERALL, 100, 1, RateIntervalUnit.MINUTES);
if (!limiter.tryAcquire()) {
    throw new TooManyRequestsException();
}
```

### 11.3 Pub/Sub Listener

```java
@Configuration
public class RedisPubSubConfig {

    @Bean
    RedisMessageListenerContainer container(
            RedisConnectionFactory factory,
            MessageListenerAdapter adapter) {
        RedisMessageListenerContainer container = new RedisMessageListenerContainer();
        container.setConnectionFactory(factory);
        container.addMessageListener(adapter, new ChannelTopic("notifications"));
        return container;
    }

    @Bean
    MessageListenerAdapter listenerAdapter(NotificationSubscriber subscriber) {
        return new MessageListenerAdapter(subscriber, "onMessage");
    }
}

@Component
@Slf4j
public class NotificationSubscriber {
    public void onMessage(String message) {
        log.info("Received notification: {}", message);
    }
}
```

### 11.4 Redis Streams Consumer

```java
@Service
@Slf4j
@RequiredArgsConstructor
public class OrderStreamConsumer {

    private final StringRedisTemplate redisTemplate;

    @Scheduled(fixedDelay = 1000)
    public void poll() {
        List<MapRecord<String, Object, Object>> records = redisTemplate.opsForStream().read(
            Consumer.from("order-processors", "instance-1"),
            StreamReadOptions.empty().count(10).block(Duration.ofSeconds(2)),
            StreamOffset.create("orders:stream", ReadOffset.lastConsumed())
        );
        if (records == null) return;
        for (MapRecord<String, Object, Object> record : records) {
            log.info("Processing order event: {}", record.getValue());
            redisTemplate.opsForStream().acknowledge("orders:stream", "order-processors", record.getId());
        }
    }
}
```

### 11.5 Cache Stampede — Single Flight

```java
public Optional<Product> getWithSingleFlight(Long id) {
    String key = "product:" + id;
    String lockKey = "lock:" + key;

    Product cached = (Product) redisTemplate.opsForValue().get(key);
    if (cached != null) return Optional.of(cached);

    Boolean acquired = redisTemplate.opsForValue()
        .setIfAbsent(lockKey, "1", Duration.ofSeconds(5));

    if (Boolean.TRUE.equals(acquired)) {
        try {
            Optional<Product> fromDb = productRepository.findById(id);
            fromDb.ifPresent(p -> redisTemplate.opsForValue().set(key, p, Duration.ofMinutes(10)));
            return fromDb;
        } finally {
            redisTemplate.delete(lockKey);
        }
    } else {
        Thread.sleep(50);
        return getWithSingleFlight(id);
    }
}
```

### 11.6 Multi-Level Cache (Caffeine L1 + Redis L2)

```
Request ──► L1 Caffeine (per pod, ~30s TTL)
                │ miss
                ▼
            L2 Redis (shared, ~10min TTL)
                │ miss
                ▼
            Database
```

```java
@Service
@RequiredArgsConstructor
public class MultiLevelProductService {

    private final Cache<Long, Product> localCache = Caffeine.newBuilder()
        .maximumSize(1_000)
        .expireAfterWrite(Duration.ofSeconds(30))
        .build();

    private final RedisTemplate<String, Object> redisTemplate;
    private final ProductRepository productRepository;

    public Product getProduct(Long id) {
        return localCache.get(id, key -> {
            String redisKey = "product:" + id;
            Product fromRedis = (Product) redisTemplate.opsForValue().get(redisKey);
            if (fromRedis != null) return fromRedis;

            Product fromDb = productRepository.findById(id)
                .orElseThrow(() -> new NotFoundException(id));
            redisTemplate.opsForValue().set(redisKey, fromDb, Duration.ofMinutes(10));
            return fromDb;
        });
    }
}
```

### 11.7 Testcontainers Integration Test

```java
@SpringBootTest
@Testcontainers
class ProductCacheServiceIT {

    @Container
    static GenericContainer<?> redis = new GenericContainer<>("redis:7.2-alpine")
        .withExposedPorts(6379);

    @DynamicPropertySource
    static void redisProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.data.redis.host", redis::getHost);
        registry.add("spring.data.redis.port", () -> redis.getMappedPort(6379));
    }

    @Autowired
    ProductCacheService productCacheService;

    @Test
    void cachesProductOnSecondRead() {
        // test cache hit/miss behavior
    }
}
```

### 11.8 Cloud Profiles

| Environment | Config |
|-------------|--------|
| **Local** | `localhost:6379`, Docker Compose |
| **AWS** | ElastiCache cluster endpoint, TLS, auth token |
| **GCP** | Memorystore, VPC peering |
| **Azure** | Azure Cache for Redis |

```yaml
# application-prod.yml
spring:
  data:
    redis:
      cluster:
        nodes:
          - redis-cluster.example.com:6379
      ssl:
        enabled: true
      password: ${REDIS_AUTH_TOKEN}
```

---

## 12. Real-World Examples

Each example follows: **Problem → Why Redis → Data model → Commands → Spring Boot sketch → Pitfalls → Interview one-liner**.

---

### 12.1 API Rate Limiter (Sliding Window)

**Problem:** Limit each user to 100 requests per minute across all API pods.

**Why Redis:** Shared counter state, atomic operations, sub-ms latency.

**Architecture:**

```
Client ──► API Gateway / Service ──► Redis (rate limit check)
                         │
                         └── Allow (200) or Reject (429)
```

**Data model:** Sorted Set — score = request timestamp (ms), member = unique request id.

```bash
ZADD ratelimit:user:42 1699000000123 "req-uuid-1"
ZREMRANGEBYSCORE ratelimit:user:42 0 1698999940123   # remove older than 60s
ZCARD ratelimit:user:42                              # count in window
EXPIRE ratelimit:user:42 120
```

**Spring Boot sketch:**

```java
@Service
@RequiredArgsConstructor
public class SlidingWindowRateLimiter {

    private final StringRedisTemplate redis;
    private static final int LIMIT = 100;
    private static final long WINDOW_MS = 60_000;

    public boolean allow(String userId) {
        String key = "ratelimit:" + userId;
        long now = System.currentTimeMillis();
        String member = UUID.randomUUID().toString();

        redis.opsForZSet().add(key, member, now);
        redis.opsForZSet().removeRangeByScore(key, 0, now - WINDOW_MS);
        Long count = redis.opsForZSet().zCard(key);
        redis.expire(key, Duration.ofMinutes(2));

        return count != null && count <= LIMIT;
    }
}
```

**Pitfalls:** Clock skew across servers (use Redis `TIME`); hot key for celebrity user (shard or local token bucket).

**Interview one-liner:** "I use a Redis Sorted Set sliding window for distributed rate limiting — O(log N) per request with automatic expiry."

**See also:** [003-rate-limiter.md](../real-world-examples/003-rate-limiter.md)

---

### 12.2 Session Store and JWT Blocklist

**Problem:** Store session data; invalidate tokens on logout.

**Why Redis:** TTL-based expiry, shared across pods, faster than DB session table.

```
Login ──► Auth Service ──► SET session:{tokenId} userPayload EX 3600
Logout ──► DEL session:{tokenId}  OR  SET blocklist:{jti} 1 EX ttl-remaining
```

**Data model:**

```bash
HSET session:abc123 userId 42 role admin loginAt 1699000000
EXPIRE session:abc123 3600

SET blocklist:jti-xyz 1 EX 1800
```

**Spring Boot sketch:**

```java
@Service
@RequiredArgsConstructor
public class SessionService {

    private final StringRedisTemplate redis;

    public void createSession(String sessionId, String userId) {
        String key = "session:" + sessionId;
        redis.opsForHash().put(key, "userId", userId);
        redis.expire(key, Duration.ofHours(2));
    }

    public boolean isTokenBlocked(String jti) {
        return Boolean.TRUE.equals(redis.hasKey("blocklist:" + jti));
    }
}
```

**Pitfalls:** Sticky sessions not needed if all state in Redis; encrypt sensitive fields.

**Interview one-liner:** "Sessions live in Redis Hash with TTL; logout adds JTI to a blocklist with remaining token lifetime."

---

### 12.3 E-Commerce Product Cache

**Problem:** Product page reads overwhelm PostgreSQL during sales.

**Why Redis:** 95%+ cache hit ratio, Hash for partial field updates.

```
GET /product/1 ──► Cache-aside ──► Redis Hash product:1
                                      │ miss
                                      └──► PostgreSQL
```

```bash
HSET product:1 name "Laptop" price 999 stock 50 category "electronics"
EXPIRE product:1 600
```

**On price update:** Update DB → `DEL product:1` or `HSET` + refresh TTL.

**Pitfalls:** Cache stampede on hot SKU — use single-flight (§11.5). Stale price — shorten TTL or event-driven invalidation from Kafka.

**See also:** [Caching Guide — Amazon example](../Caching-Complete-Guide.md)

**Interview one-liner:** "Product cache-aside with Hash structure, 10-minute TTL, and Kafka-driven invalidation on inventory/price change."

---

### 12.4 Gaming Leaderboard

**Problem:** Real-time top-100 players globally.

**Why Redis:** Sorted Set provides O(log N) rank updates and range queries.

```bash
ZADD leaderboard:global 2500 "player:alice"
ZINCRBY leaderboard:global 100 "player:alice"
ZREVRANGE leaderboard:global 0 99 WITHSCORES
ZRANK leaderboard:global "player:alice"
```

**Spring Boot:**

```java
public List<LeaderboardEntry> topPlayers(String board, int n) {
    Set<ZSetOperations.TypedTuple<String>> tuples =
        redis.opsForZSet().reverseRangeWithScores(board, 0, n - 1);
    // map to DTOs
}
```

**Pitfalls:** Massive board (millions) — shard by region `leaderboard:us-east`.

**Interview one-liner:** "Sorted Sets give O(log N) score updates and O(log N + M) top-N queries — ideal for leaderboards."

---

### 12.5 Inventory Reservation

**Problem:** Prevent overselling when 1000 users buy the last 10 items.

**Why Redis:** Atomic `DECR` + distributed lock for checkout window.

```
Reserve flow:
  1. SET lock:sku:99 NX EX 10
  2. DECR stock:99  (or check >= 0 in Lua)
  3. If fail, release lock
  4. On payment timeout, INCR stock:99
```

**Lua atomic decrement:**

```lua
local stock = tonumber(redis.call('GET', KEYS[1]) or '0')
if stock >= tonumber(ARGV[1]) then
  return redis.call('DECRBY', KEYS[1], ARGV[1])
else
  return -1
end
```

**Pitfalls:** Redis is not source of truth for inventory long-term — sync to DB. Use Redisson lock for multi-step checkout.

**Interview one-liner:** "Redis holds tentative stock counters with Lua atomicity; PostgreSQL is authoritative and reconciled async."

---

### 12.6 Idempotent Payment API

**Problem:** Client retries POST /pay — must not double-charge.

**Why Redis:** `SET idempotency:{key} NX EX` is atomic and fast.

```
POST /pay  Idempotency-Key: abc-123
  │
  ├── SET idempotency:abc-123 "processing" NX EX 86400
  │     └── fails → return cached response
  ├── Process payment
  └── SET idempotency:abc-123 "{response json}" XX EX 86400
```

**Spring Boot:**

```java
public ResponseEntity<PaymentResult> pay(String idempotencyKey, PaymentRequest req) {
    String key = "idempotency:" + idempotencyKey;
    Boolean first = redis.opsForValue().setIfAbsent(key, "processing", Duration.ofHours(24));
    if (Boolean.FALSE.equals(first)) {
        String cached = redis.opsForValue().get(key);
        return ResponseEntity.ok(deserialize(cached));
    }
    PaymentResult result = paymentProcessor.charge(req);
    redis.opsForValue().set(key, serialize(result), Duration.ofHours(24));
    return ResponseEntity.ok(result);
}
```

**Pitfalls:** Store **response**, not just flag. Key must be client-supplied UUID.

**See also:** [007-payment-system.md](../real-world-examples/007-payment-system.md)

**Interview one-liner:** "Idempotency keys map to SET NX with cached response body — safe retries without duplicate charges."

---

### 12.7 Notification Fan-Out

**Problem:** Push real-time notification to online users.

**Why Redis Pub/Sub:** Low-latency broadcast; each service instance subscribes.

```
Order Service ──PUBLISH──► channel:user:42 ──► Notification Service (subscriber)
                                                    │
                                                    └── WebSocket to client
```

```bash
PUBLISH channel:user:42 '{"type":"ORDER_SHIPPED","orderId":123}'
```

**Pitfalls:** No persistence — offline users miss message. Use Streams or Kafka for guaranteed delivery.

**See also:** [004-notification-system.md](../real-world-examples/004-notification-system.md)

**Interview one-liner:** "Pub/Sub for online real-time push; Kafka for durable notification pipeline to email/SMS workers."

---

### 12.8 URL Shortener Hot Path

**Problem:** Redirect `abc.xyz/x` must be < 10ms at billions of hits.

**Why Redis:** O(1) GET for short code → long URL mapping.

```
GET redirect:x  →  "https://example.com/very/long/url"
```

```bash
SET redirect:x "https://example.com/long" EX 86400
```

**Architecture:**

```
User ──► CDN (optional) ──► Redirect Service ──► Redis
                                    │ miss
                                    └──► DB + populate cache
```

**Pitfalls:** Viral link = hot key — local cache + read replicas.

**See also:** [001-url-shortener-complete.md](../real-world-examples/001-url-shortener-complete.md)

**Interview one-liner:** "Redirect path is cache-aside String lookup in Redis with DB fallback; CDN handles static edge caching."

---

### 12.9 Real-World Examples Summary

| Example | Redis structures | Primary pattern |
|---------|------------------|-----------------|
| Rate limiter | ZSet / String | Sliding window / token bucket |
| Session | Hash + TTL | Authoritative ephemeral state |
| Product cache | Hash / String | Cache-aside |
| Leaderboard | Sorted Set | In-memory ranking |
| Inventory | String + Lua / Lock | Coordination |
| Idempotency | String SET NX | Dedup |
| Notifications | Pub/Sub or Stream | Fan-out |
| URL shortener | String | Hot path cache |

---

## 13. Advanced Interview Questions

### Q1: Explain Redis Cluster hash slots.

**Answer:** 16,384 slots partition the keyspace. `CRC16(key) mod 16384` maps each key to a slot owned by one master. Replicas provide HA. Hash tags `{userId}` co-locate related keys. Clients follow `MOVED` redirects when topology changes.

### Q2: What is replication lag and how do you handle it?

**Answer:** Replicas copy primary asynchronously. Lag causes stale reads. Fix: read from primary after writes, use `WAIT` for critical paths, or design for eventual consistency with short TTL.

### Q3: Debate: Is Redlock safe?

**Answer:** Redlock acquires majority of independent Redis nodes. Critics (Kleppmann) note clock drift and lack of fencing. In interviews: acknowledge controversy; prefer fencing tokens, short TTL, or etcd for strict correctness.

### Q4: How do you design cache invalidation in microservices?

**Answer:** Options: (1) TTL only for low-risk data, (2) domain events via Kafka → all services `DEL` keys, (3) versioned keys `product:1:v5`, (4) write-through for critical data. Never invalidate before DB commit.

### Q5: Redis down — what happens to your app?

**Answer:** Cache-aside should **degrade to DB** (catch `DataAccessException`). Circuit breaker prevents hammering dead Redis. For sessions/rate limits, outage is user-visible — need Redis HA (Sentinel/Cluster) and multi-AZ.

### Q6: How to find and fix a hot key?

**Answer:** Monitor `redis-cli --hotkeys` or latency spikes per command. Fix: local L1 cache, key splitting, read replicas, or pre-compute (CDN). Prevent: avoid monotonic keys on one shard.

### Q7: Big key deletion blocks Redis — what do you do?

**Answer:** Use `UNLINK` for async deletion. Prevent by chunking large collections. Monitor with `MEMORY USAGE key`. Schedule deletions off-peak.

### Q8: When would you choose Kafka over Redis Streams?

**Answer:** Need long retention (days/weeks), replay across many teams, high throughput event backbone, schema registry. Redis Streams for simple in-cluster queuing with lower ops overhead.

### Q9: How do you size Redis memory?

**Answer:** Estimate: `(avg_key_size + overhead) × key_count × replication_factor`. Add 30% for fragmentation. Set `maxmemory` below physical RAM. Load test with production-like key distribution.

### Q10: Design a distributed cache — where do you start?

**Answer:** Clarify requirements (latency, consistency, data structures). For **using** Redis: cache-aside, TTL, eviction policy, HA mode. For **building** Redis: point to consistent hashing, replication, failure detection — see [009-distributed-cache-redis-like.md](../real-world-examples/009-distributed-cache-redis-like.md).

### Q11: Explain cache penetration, breakdown, avalanche with Redis fixes.

| Problem | Fix in Redis |
|---------|--------------|
| **Penetration** | Cache null with short TTL; Bloom filter |
| **Breakdown** | Mutex SETNX; never expire hot keys without refresh |
| **Avalanche** | TTL jitter; random stagger; multi-level cache |

### Q12: How does Redis achieve high performance?

**Answer:** In-memory data, single-threaded command loop (no locks), I/O multiplexing, efficient encodings (ziplist → hashtable upgrade), pipelining, and optional I/O threads for network.

### Q13: Compare ElastiCache vs self-hosted Redis.

| Aspect | Managed (ElastiCache) | Self-hosted |
|--------|----------------------|-------------|
| Ops | AWS handles patching, failover | You manage Sentinel/Cluster |
| Cost | Higher $/GB | Lower infra, higher labor |
| Features | Version limits | Full control |
| Best for | Most production teams | Custom modules, cost at scale |

### Q14: What metrics prove your cache is healthy?

**Answer:** Hit ratio > 90%, p99 latency < 2ms, zero OOM evictions, low `evicted_keys`, replication lag < 1s, connection pool utilization < 80%.

### Q15: Spring @Cacheable vs manual cache-aside?

**Answer:** `@Cacheable` is cleaner for simple CRUD. Manual cache-aside gives control over key design, stampede protection, conditional caching, and graceful Redis failure handling — preferred in high-scale interviews.

---

## 14. Quick Reference Cheat Sheet — Part 2

### Cluster / Sentinel Commands

```bash
CLUSTER INFO
CLUSTER NODES
CLUSTER KEYSLOT mykey
redis-cli -c -h host -p 7000

SENTINEL masters
SENTINEL replicas mymaster
SENTINEL failover mymaster
```

### Lock Pattern

```bash
SET lock:resource uuid NX EX 30
# unlock via Lua comparing uuid
```

### Stream Pattern

```bash
XADD stream * field value
XGROUP CREATE stream group $ MKSTREAM
XREADGROUP GROUP group consumer STREAMS stream >
XACK stream group id
```

### Production Checklist

| Area | Check |
|------|-------|
| Memory | `maxmemory` + eviction policy set |
| HA | Sentinel or Cluster in prod |
| Persistence | RDB/AOF aligned with RPO |
| Security | ACL, TLS, private network |
| Monitoring | Exporter, hit ratio, lag alerts |
| App | TTL on all keys, degrade path if Redis down |
| Keys | No `KEYS *`, naming convention documented |

### Spring Boot Patterns Index

| Pattern | Technology |
|---------|------------|
| Declarative cache | `@Cacheable` + `RedisCacheManager` |
| Manual cache-aside | `RedisTemplate` |
| Distributed lock | Redisson / SET NX + Lua |
| Rate limit | ZSet sliding window / Redisson |
| Pub/Sub | `RedisMessageListenerContainer` |
| Streams | `opsForStream()` |
| Multi-level | Caffeine + Redis |
| Integration test | Testcontainers |

---

> **Back to [Part 1 (Core)](./01-Redis-Complete-Learning-Guide.md)** for fundamentals, data structures, and basic Spring Boot setup.
