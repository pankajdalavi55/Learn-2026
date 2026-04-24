# Complete System Design: Distributed Rate Limiter (Production-Ready)

> **Complexity Level:** Intermediate to Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** Cloudflare Rate Limiting, AWS API Gateway, Stripe API, GitHub API

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Rate Limiting Algorithms](#7-deep-dive-rate-limiting-algorithms)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a distributed rate limiter that can throttle requests based on configurable rules."

**Interviewer's Perspective:**  
This is a frequently asked intermediate-to-advanced problem that assesses:
- Understanding of distributed systems and coordination
- Knowledge of rate limiting algorithms and their trade-offs
- Consistency vs availability decisions under partitions
- Ability to handle race conditions in concurrent environments
- Real-world awareness of how APIs protect themselves (Stripe: 100 req/sec, GitHub: 5,000 req/hour)

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I dive into the design, I'd like to clarify both the functional and non-functional requirements."

### 2.1 Functional Requirements

**Candidate:** "For functional requirements, I'd like to confirm:
1. Should the rate limiter throttle by user ID, IP address, API key, or all of these?
2. Should the limiting rules be configurable at runtime, or are they baked into code?
3. What response should a throttled client receive?
4. Do we need to support multiple strategies — e.g., per-endpoint limits, global limits, tiered limits?
5. Should clients be informed of their remaining quota?"

**Interviewer:** "Good questions. Here's what we need:
- Throttle by user ID, IP address, and API key — all three
- Rules must be dynamically configurable without redeployment
- Return HTTP 429 Too Many Requests with appropriate headers
- Support multiple strategies (per-endpoint, per-user, global)
- Yes, include rate limit headers in responses"

**Candidate:** "Got it. So the core functional requirements are:
1. ✅ Limit requests per user/IP/API key with configurable rules
2. ✅ Configurable rate (e.g., 100 requests per 60-second window)
3. ✅ Return HTTP 429 with `X-RateLimit-*` headers when throttled
4. ✅ Support multiple limiting strategies (fixed window, sliding window, token bucket)
5. ✅ CRUD API for managing rate limit rules at runtime
6. ✅ Support hierarchical limits (per-endpoint + per-user + global)"

### 2.2 Non-Functional Requirements

**Candidate:** "For non-functional requirements:
1. What's the maximum latency overhead the rate limiter can introduce?
2. What's the availability target?
3. Should it work across multiple data centers or a single region?
4. How strict do we need to be — can we tolerate a small percentage of over-limit requests?
5. What's the expected request throughput?"

**Interviewer:**
- Latency overhead: must be under 5ms at p99
- Availability: 99.99% — the limiter must never be a single point of failure
- Multi-region: yes, requests come from globally distributed API gateways
- Tolerance: small over-count is acceptable (< 1% false negatives), but avoid false positives
- Scale: 1 million requests/second across the entire platform

**Candidate:** "Understood. Summarizing non-functional requirements:
1. ✅ Low latency: < 5ms overhead per request at p99
2. ✅ High availability: 99.99% uptime — must not block legitimate traffic
3. ✅ Distributed: works across multiple nodes and data centers
4. ✅ Minimal false positives (blocking legitimate traffic)
5. ✅ Acceptable false negatives < 1% (letting a few extra requests through)
6. ✅ Scale: 1 million requests/second globally"

---

## 3. Scale Estimation

**Candidate:** "Let me estimate the scale to drive our infrastructure decisions."

### 3.1 Traffic Estimation

| Metric | Value |
|--------|-------|
| Peak throughput | 1,000,000 requests/sec |
| Daily requests | 1M × 86,400 = **86.4 billion/day** |
| Monthly requests | ~2.6 trillion/month |
| Average request size (for rate check) | ~200 bytes (key + metadata) |

### 3.2 Storage Estimation

**Rule Configuration Storage:**
| Item | Estimate |
|------|----------|
| Average rule size | ~10 KB (JSON with conditions, thresholds) |
| Total rules | 100,000 (across all tenants/endpoints) |
| Total rule storage | 100K × 10 KB = **1 GB** |
| Storage type | PostgreSQL / DynamoDB (durable, infrequently written) |

**Counter Storage (Hot Path):**
| Item | Estimate |
|------|----------|
| Bytes per counter | ~100 bytes (key + count + timestamp + TTL) |
| Active counters | 10 million (unique user-endpoint-window combos) |
| Total counter storage | 10M × 100 B = **1 GB** |
| Storage type | Redis Cluster (in-memory, fast reads/writes) |

### 3.3 Bandwidth Estimation

| Direction | Calculation | Total |
|-----------|------------|-------|
| Incoming (rate check request) | 1M req/s × 200 B | **200 MB/s** |
| Outgoing (rate check response) | 1M req/s × 50 B | **50 MB/s** |
| Redis traffic (read + write) | 1M req/s × 150 B | **150 MB/s** |

### 3.4 Infrastructure Estimation

**Candidate:** "Based on these numbers:
- **Redis Cluster:** 6-node cluster (3 masters + 3 replicas) — each master handles ~333K ops/sec, well within Redis's ~1M ops/sec per node capacity
- **Rate Limiter Service:** 20-30 stateless instances behind a load balancer
- **Rules DB:** Single PostgreSQL instance with read replicas (low write volume)
- **Network:** ~400 MB/s total bandwidth — manageable within a modern data center"

---

## 4. High-Level Design

### 4.1 Architecture Diagram

**Candidate:** "Here's the high-level architecture:"

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                     │
│           Mobile App  |  Web App  |  Third-Party API Consumer            │
└──────────────────────────────┬───────────────────────────────────────────┘
                               │
                        ┌──────▼──────┐
                        │ Load        │
                        │ Balancer    │
                        └──────┬──────┘
                               │
                    ┌──────────▼──────────┐
                    │    API Gateway      │  (Authentication, Routing)
                    │  (Kong / Envoy)     │
                    └──────────┬──────────┘
                               │
              ┌────────────────▼────────────────┐
              │     Rate Limiter Middleware      │
              │  (Embedded in gateway or sidecar)│
              │                                  │
              │  1. Extract key (user/IP/API key)│
              │  2. Lookup matching rules         │
              │  3. Check counter against limit   │
              │  4. Allow or reject (HTTP 429)    │
              └───────┬──────────────┬───────────┘
                      │              │
           ┌──────────▼──┐    ┌─────▼──────────┐
           │ Redis        │    │  Rules Cache   │
           │ Cluster      │    │  (Local/Redis) │
           │              │    │                │
           │ - Counters   │    │ - Rule configs │
           │ - Sliding    │    │ - TTL: 30s     │
           │   window logs│    │ - Invalidation │
           │ - Token      │    │   via pub/sub  │
           │   buckets    │    └───────┬────────┘
           └──────────────┘            │
                                ┌──────▼──────┐
                                │ Rules DB    │
                                │ (PostgreSQL)│
                                │             │
                                │ - Rule CRUD │
                                │ - Audit log │
                                └─────────────┘
                               │
              ┌────────────────▼────────────────┐
              │        Backend Services          │
              │  (Order API, Payment API, etc.)   │
              └──────────────────────────────────┘
```

### 4.2 API Design

**Candidate:** "We need two categories of APIs: the rate limit check (internal) and the rule management CRUD (admin)."

**1. Rate Limit Check (Internal — called by gateway middleware)**
```http
POST /internal/ratelimit/check
Content-Type: application/json

Request:
{
  "entity_type": "api_key",
  "entity_id": "sk_live_abc123",
  "endpoint": "/v1/charges",
  "method": "POST"
}

Response (200 OK — Allowed):
{
  "allowed": true,
  "limit": 100,
  "remaining": 73,
  "reset_at": 1714003260
}

Response (429 Too Many Requests — Throttled):
{
  "allowed": false,
  "limit": 100,
  "remaining": 0,
  "reset_at": 1714003260,
  "retry_after": 12
}
```

**2. Create Rate Limit Rule (Admin)**
```http
POST /api/v1/rules
Authorization: Bearer <admin_token>
Content-Type: application/json

Request:
{
  "name": "Stripe-style API key limit",
  "entity_type": "api_key",
  "endpoint_pattern": "/v1/*",
  "max_requests": 100,
  "window_seconds": 60,
  "algorithm": "sliding_window_counter",
  "action": "reject",
  "priority": 10
}

Response (201 Created):
{
  "rule_id": "rl_rule_7f3a2b",
  "name": "Stripe-style API key limit",
  "created_at": "2026-04-24T10:00:00Z",
  "status": "active"
}
```

**3. List Rules**
```http
GET /api/v1/rules?entity_type=api_key&status=active

Response:
{
  "rules": [ ... ],
  "total": 42,
  "page": 1
}
```

**4. Update / Delete Rule**
```http
PATCH /api/v1/rules/{rule_id}
DELETE /api/v1/rules/{rule_id}
```

### 4.3 Data Flow

**Candidate:** "Let me trace through both the allowed and throttled request paths."

**Flow 1: Request Allowed**
```
1. Client → POST /v1/charges (with API key header)
2. API Gateway authenticates the request, resolves entity (api_key: sk_live_abc123)
3. Rate Limiter Middleware:
   a. Builds key: "rl:api_key:sk_live_abc123:/v1/charges:window_1714003200"
   b. Looks up matching rule from local cache (or Redis → PostgreSQL on miss)
   c. Executes atomic Redis INCR + EXPIRE (or Lua script) against the counter
   d. Counter = 74, limit = 100 → ALLOWED
4. Middleware adds headers:
      X-RateLimit-Limit: 100
      X-RateLimit-Remaining: 26
      X-RateLimit-Reset: 1714003260
5. Request forwarded to backend service → processes charge
6. Response returned to client with rate limit headers
```

**Flow 2: Request Throttled**
```
1. Client → POST /v1/charges (with API key header)
2. API Gateway authenticates, resolves entity
3. Rate Limiter Middleware:
   a. Builds key, looks up rule
   b. Executes atomic Redis INCR
   c. Counter = 101, limit = 100 → REJECTED
   d. Decrements counter back (or uses Lua to check-before-increment)
4. Returns immediately:
      HTTP 429 Too Many Requests
      X-RateLimit-Limit: 100
      X-RateLimit-Remaining: 0
      X-RateLimit-Reset: 1714003260
      Retry-After: 12
      Body: {"error": "rate_limit_exceeded", "message": "Too many requests"}
5. Request is NOT forwarded to backend — saving resources
6. Throttle event logged asynchronously for monitoring
```

---

## 5. Deep Dive: Core Components

### 5.1 Rate Limiter Middleware

**Candidate:** "The rate limiter is best implemented as middleware in the API gateway, not as a separate network hop, to minimize latency."

**Deployment Options:**

| Option | Latency | Complexity | When to Use |
|--------|---------|-----------|-------------|
| Gateway plugin (Kong, Envoy filter) | ~1-2ms | Low | Single gateway platform |
| Sidecar proxy (Istio/Envoy) | ~2-3ms | Medium | Service mesh environments |
| Standalone service | ~3-10ms | High | Multi-platform / legacy systems |
| In-process library | ~0.5ms | Low | Monolith or single-language stack |

**Recommended: Gateway plugin** — keeps it on the hot path with minimal latency.

```javascript
// Rate limiter middleware (Express.js example)
const Redis = require('ioredis');
const redis = new Redis.Cluster([
  { host: 'redis-1', port: 6379 },
  { host: 'redis-2', port: 6379 },
  { host: 'redis-3', port: 6379 },
]);

async function rateLimiterMiddleware(req, res, next) {
  const entityKey = extractEntityKey(req); // e.g., "api_key:sk_live_abc123"
  const rule = await getRuleForRequest(req, entityKey);

  if (!rule) {
    return next(); // no rule matched — allow
  }

  const result = await checkRateLimit(entityKey, rule);

  res.set({
    'X-RateLimit-Limit': rule.max_requests,
    'X-RateLimit-Remaining': Math.max(0, result.remaining),
    'X-RateLimit-Reset': result.resetAt,
  });

  if (!result.allowed) {
    return res.status(429).json({
      error: 'rate_limit_exceeded',
      message: `Rate limit of ${rule.max_requests} requests per ${rule.window_seconds}s exceeded`,
      retry_after: result.retryAfter,
    });
  }

  next();
}

function extractEntityKey(req) {
  if (req.headers['x-api-key']) return `api_key:${req.headers['x-api-key']}`;
  if (req.user?.id)             return `user:${req.user.id}`;
  return `ip:${req.ip}`;
}
```

### 5.2 Rules Engine

**Candidate:** "The rules engine matches incoming requests against configured rules. Rules have priorities, conditions, and an associated algorithm."

```javascript
class RulesEngine {
  constructor(rulesCache) {
    this.rulesCache = rulesCache;
  }

  async getRuleForRequest(req, entityKey) {
    const rules = await this.rulesCache.getActiveRules();
    const entityType = entityKey.split(':')[0]; // "api_key", "user", "ip"

    const matchingRules = rules
      .filter(rule => rule.entity_type === entityType || rule.entity_type === '*')
      .filter(rule => this.matchEndpoint(rule.endpoint_pattern, req.path))
      .filter(rule => !rule.method || rule.method === req.method)
      .sort((a, b) => b.priority - a.priority); // higher priority first

    return matchingRules[0] || null; // return highest priority match
  }

  matchEndpoint(pattern, path) {
    if (pattern === '*') return true;
    const regex = new RegExp('^' + pattern.replace(/\*/g, '.*') + '$');
    return regex.test(path);
  }
}
```

**Rule Caching Strategy:**
- Rules are cached locally in each gateway node with a 30-second TTL
- Redis Pub/Sub broadcasts invalidation on rule changes
- Fallback to stale cache if Redis is unreachable (availability over freshness)

```javascript
class RulesCache {
  constructor(redis, db) {
    this.redis = redis;
    this.db = db;
    this.localCache = new Map();
    this.lastRefresh = 0;
    this.TTL_MS = 30_000;

    this.redis.subscribe('rule_updates', () => {
      this.invalidateLocal();
    });
  }

  async getActiveRules() {
    if (Date.now() - this.lastRefresh < this.TTL_MS && this.localCache.size > 0) {
      return Array.from(this.localCache.values());
    }

    try {
      const cached = await this.redis.get('rl:rules:active');
      if (cached) {
        const rules = JSON.parse(cached);
        this.refreshLocal(rules);
        return rules;
      }
    } catch (err) {
      if (this.localCache.size > 0) return Array.from(this.localCache.values());
    }

    const rules = await this.db.query('SELECT * FROM rate_limit_rules WHERE status = $1', ['active']);
    await this.redis.setex('rl:rules:active', 60, JSON.stringify(rules));
    this.refreshLocal(rules);
    return rules;
  }

  refreshLocal(rules) {
    this.localCache.clear();
    rules.forEach(r => this.localCache.set(r.rule_id, r));
    this.lastRefresh = Date.now();
  }

  invalidateLocal() {
    this.lastRefresh = 0;
  }
}
```

### 5.3 Technology Choices: Redis vs Memcached

| Feature | Redis | Memcached |
|---------|-------|-----------|
| Atomic operations | ✅ INCR, Lua scripts | ✅ incr/decr only |
| Data structures | Sorted sets, hashes, streams | Key-value only |
| Persistence | Optional RDB/AOF | None |
| Cluster mode | ✅ Native | ✅ Client-side sharding |
| Pub/Sub | ✅ Built-in | ❌ Not available |
| Lua scripting | ✅ Atomic multi-step ops | ❌ Not available |
| TTL granularity | Milliseconds | Seconds |
| **Verdict** | **✅ Preferred** — Lua scripts enable atomic check-and-increment | ⚠️ Only for simple fixed-window counters |

**Candidate:** "Redis is the clear choice here. The ability to run Lua scripts atomically is critical for avoiding race conditions in distributed counters. Sorted sets are also essential for sliding window log implementations."

---

## 6. Deep Dive: Database Design

### 6.1 Rate Limit Rules Schema

```sql
CREATE TABLE rate_limit_rules (
    rule_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    entity_type     VARCHAR(50) NOT NULL,  -- 'api_key', 'user', 'ip', '*'
    entity_id       VARCHAR(255),          -- specific entity or NULL for all
    endpoint_pattern VARCHAR(500) DEFAULT '*',
    method          VARCHAR(10),           -- 'GET', 'POST', or NULL for all
    max_requests    INTEGER NOT NULL,
    window_seconds  INTEGER NOT NULL,
    algorithm       VARCHAR(50) NOT NULL DEFAULT 'sliding_window_counter',
        -- 'fixed_window', 'sliding_window_log', 'sliding_window_counter',
        -- 'token_bucket', 'leaky_bucket'
    burst_size      INTEGER,               -- for token bucket: max burst
    refill_rate     DECIMAL(10,2),         -- for token bucket: tokens/sec
    action          VARCHAR(20) NOT NULL DEFAULT 'reject',
        -- 'reject', 'throttle', 'log_only'
    priority        INTEGER NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
        -- 'active', 'inactive', 'testing'
    tier            VARCHAR(50),           -- 'free', 'pro', 'enterprise'
    created_by      VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT chk_algorithm CHECK (algorithm IN (
        'fixed_window', 'sliding_window_log', 'sliding_window_counter',
        'token_bucket', 'leaky_bucket'
    )),
    CONSTRAINT chk_action CHECK (action IN ('reject', 'throttle', 'log_only')),
    CONSTRAINT chk_status CHECK (status IN ('active', 'inactive', 'testing'))
);

CREATE INDEX idx_rules_entity_type ON rate_limit_rules(entity_type, status);
CREATE INDEX idx_rules_priority ON rate_limit_rules(priority DESC);
CREATE INDEX idx_rules_tier ON rate_limit_rules(tier);
```

### 6.2 Rate Limit Overrides Schema

```sql
CREATE TABLE rate_limit_overrides (
    override_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id         UUID NOT NULL REFERENCES rate_limit_rules(rule_id),
    entity_id       VARCHAR(255) NOT NULL,  -- specific user/API key
    max_requests    INTEGER NOT NULL,
    window_seconds  INTEGER NOT NULL,
    reason          TEXT,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(rule_id, entity_id)
);

CREATE INDEX idx_overrides_entity ON rate_limit_overrides(entity_id);
```

### 6.3 Rate Limit Counters (Redis Schemas)

**Fixed Window Counter:**
```
Key:    rl:{entity_type}:{entity_id}:{endpoint}:{window_start}
Value:  integer counter
TTL:    window_seconds + buffer
Example: rl:api_key:sk_live_abc:charges:1714003200 → 73
```

**Sliding Window Log:**
```
Key:    rl:log:{entity_type}:{entity_id}:{endpoint}
Type:   Sorted Set (ZSET)
Members: unique request IDs or timestamps
Scores:  request timestamps (epoch ms)
Example: ZRANGEBYSCORE rl:log:api_key:sk_live_abc:charges 1714003200000 +inf
```

**Token Bucket:**
```
Key:    rl:bucket:{entity_type}:{entity_id}:{endpoint}
Type:   Hash
Fields:
  - tokens:       current tokens remaining (float)
  - last_refill:  epoch timestamp of last refill (ms)
Example: HGETALL rl:bucket:api_key:sk_live_abc:charges
         → {"tokens": "42.5", "last_refill": "1714003245000"}
```

### 6.4 Audit & Analytics Tables

```sql
CREATE TABLE rate_limit_events (
    event_id        BIGSERIAL PRIMARY KEY,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       VARCHAR(255) NOT NULL,
    endpoint        VARCHAR(500),
    rule_id         UUID REFERENCES rate_limit_rules(rule_id),
    action_taken    VARCHAR(20) NOT NULL,  -- 'allowed', 'rejected', 'logged'
    current_count   INTEGER,
    limit_value     INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Partition per day for efficient cleanup
CREATE TABLE rate_limit_events_20260424 PARTITION OF rate_limit_events
    FOR VALUES FROM ('2026-04-24') TO ('2026-04-25');

-- Aggregate stats (materialized by a cron job)
CREATE TABLE rate_limit_stats_hourly (
    stat_id         BIGSERIAL PRIMARY KEY,
    rule_id         UUID NOT NULL,
    entity_type     VARCHAR(50),
    hour_bucket     TIMESTAMPTZ NOT NULL,
    total_requests  BIGINT NOT NULL DEFAULT 0,
    rejected_count  BIGINT NOT NULL DEFAULT 0,
    avg_counter     DECIMAL(10,2),
    p99_latency_us  INTEGER,
    
    UNIQUE(rule_id, entity_type, hour_bucket)
);

CREATE INDEX idx_stats_hour ON rate_limit_stats_hourly(hour_bucket);
```

---

## 7. Deep Dive: Rate Limiting Algorithms

**Candidate:** "This is the heart of the system. There are five major algorithms, each with distinct trade-offs. Let me walk through all of them."

---

### 7.1 Fixed Window Counter

**How it works:** Divide time into fixed windows (e.g., 60-second intervals). Maintain a counter per window. Increment on each request; reject if counter exceeds the limit.

```
Timeline:    |--- Window 1 ---|--- Window 2 ---|--- Window 3 ---|
             00:00       01:00  01:00      02:00  02:00      03:00
Requests:      |||||||||||        ||||||||          ||||||
Counter:          45                34                 22
Limit:           100               100                100
Result:        ALLOW             ALLOW              ALLOW
```

**Implementation:**

```javascript
async function fixedWindowCheck(redis, key, rule) {
  const windowStart = Math.floor(Date.now() / 1000 / rule.window_seconds) * rule.window_seconds;
  const redisKey = `${key}:${windowStart}`;

  const count = await redis.incr(redisKey);

  if (count === 1) {
    await redis.expire(redisKey, rule.window_seconds + 1);
  }

  return {
    allowed: count <= rule.max_requests,
    remaining: Math.max(0, rule.max_requests - count),
    resetAt: windowStart + rule.window_seconds,
  };
}
```

```python
import time, redis

def fixed_window_check(r: redis.Redis, key: str, max_requests: int, window_sec: int) -> dict:
    window_start = int(time.time()) // window_sec * window_sec
    redis_key = f"{key}:{window_start}"

    count = r.incr(redis_key)
    if count == 1:
        r.expire(redis_key, window_sec + 1)

    return {
        "allowed": count <= max_requests,
        "remaining": max(0, max_requests - count),
        "reset_at": window_start + window_sec,
    }
```

**Pros:**
- Simple to implement and understand
- Memory efficient — one counter per entity per window
- Fast: single Redis `INCR` command

**Cons: The Boundary Burst Problem**

```
         Window 1                    Window 2
    |---- 60 sec ----|          |---- 60 sec ----|
                     ▼          ▼
    .......[100 reqs]|[100 reqs].................
           ↑ last 5s   first 5s ↑
           of window 1  of window 2

    Result: 200 requests in a 10-second span!
    The limit is 100/min, but a burst at the boundary
    allows 2× the intended rate.
```

This is the most commonly asked follow-up question. The boundary burst can allow up to **2× the configured limit** in the worst case.

---

### 7.2 Sliding Window Log

**How it works:** Store the timestamp of every request in a sorted set. For each new request, count how many timestamps fall within the last `window_seconds`. Reject if the count exceeds the limit.

```
Timeline (window = 60s, limit = 5):

Sorted Set (ZSET):
  Score (timestamp)    Member
  ─────────────────    ─────────
  1714003201.000       req_001
  1714003215.000       req_002
  1714003228.000       req_003
  1714003241.000       req_004
  1714003252.000       req_005
  1714003258.000       req_006  ← New request

Step 1: Remove entries older than (now - 60) = 1714003198
Step 2: Count remaining = 5
Step 3: 5 >= limit(5) → REJECT req_006
```

**Implementation:**

```javascript
async function slidingWindowLogCheck(redis, key, rule) {
  const now = Date.now();
  const windowStart = now - (rule.window_seconds * 1000);
  const requestId = `${now}:${Math.random().toString(36).slice(2, 8)}`;

  const result = await redis.multi()
    .zremrangebyscore(key, 0, windowStart)           // remove expired entries
    .zcard(key)                                       // count current entries
    .zadd(key, now, requestId)                        // add this request
    .expire(key, rule.window_seconds + 1)             // set TTL
    .exec();

  const currentCount = result[1][1]; // count BEFORE adding current request

  if (currentCount >= rule.max_requests) {
    await redis.zrem(key, requestId); // rollback
    const oldestTimestamp = await redis.zrange(key, 0, 0, 'WITHSCORES');
    const resetAt = oldestTimestamp.length > 1
      ? Math.ceil((parseFloat(oldestTimestamp[1]) + rule.window_seconds * 1000) / 1000)
      : Math.ceil(now / 1000) + rule.window_seconds;

    return { allowed: false, remaining: 0, resetAt };
  }

  return {
    allowed: true,
    remaining: rule.max_requests - currentCount - 1,
    resetAt: Math.ceil(now / 1000) + rule.window_seconds,
  };
}
```

```python
import time, redis

def sliding_window_log_check(r: redis.Redis, key: str, max_requests: int, window_sec: int) -> dict:
    now_ms = time.time() * 1000
    window_start_ms = now_ms - (window_sec * 1000)
    request_id = f"{now_ms}:{id(object())}"

    pipe = r.pipeline()
    pipe.zremrangebyscore(key, 0, window_start_ms)
    pipe.zcard(key)
    pipe.zadd(key, {request_id: now_ms})
    pipe.expire(key, window_sec + 1)
    results = pipe.execute()

    current_count = results[1]

    if current_count >= max_requests:
        r.zrem(key, request_id)  # rollback
        return {"allowed": False, "remaining": 0, "reset_at": int(time.time()) + window_sec}

    return {
        "allowed": True,
        "remaining": max_requests - current_count - 1,
        "reset_at": int(time.time()) + window_sec,
    }
```

**Pros:**
- Perfectly accurate — no boundary burst problem
- Exact count of requests in any rolling window

**Cons:**
- Memory intensive: stores every request timestamp (at scale, this is significant)
- At 100 req/min per user with 10M users: 10M × 100 × ~20 bytes = **~20 GB** in Redis
- Multiple Redis commands per check (though pipelined)

---

### 7.3 Sliding Window Counter (Hybrid)

**How it works:** Combines fixed window counters with a weighted average to approximate a sliding window. Uses two counters (current window + previous window) and calculates a weighted count.

```
                 Previous Window          Current Window
                |---- 60 sec ----|      |---- 60 sec ----|
                                  ←─ 45s into current window ─→

Previous count: 84               Current count: 36
Weight for previous = 1 - (elapsed / window) = 1 - (45/60) = 0.25

Weighted count = (previous × weight) + current
               = (84 × 0.25) + 36
               = 21 + 36
               = 57

Limit = 100 → 57 < 100 → ALLOWED
```

**Implementation:**

```javascript
async function slidingWindowCounterCheck(redis, key, rule) {
  const now = Math.floor(Date.now() / 1000);
  const windowSize = rule.window_seconds;
  const currentWindow = Math.floor(now / windowSize) * windowSize;
  const previousWindow = currentWindow - windowSize;
  const elapsedInCurrent = now - currentWindow;
  const weight = 1 - (elapsedInCurrent / windowSize);

  const currentKey = `${key}:${currentWindow}`;
  const previousKey = `${key}:${previousWindow}`;

  const [prevCount, currCount] = await redis.mget(previousKey, currentKey);
  const prev = parseInt(prevCount) || 0;
  const curr = parseInt(currCount) || 0;
  const estimatedCount = Math.floor(prev * weight) + curr;

  if (estimatedCount >= rule.max_requests) {
    return {
      allowed: false,
      remaining: 0,
      resetAt: currentWindow + windowSize,
    };
  }

  await redis.multi()
    .incr(currentKey)
    .expire(currentKey, windowSize * 2 + 1)
    .exec();

  return {
    allowed: true,
    remaining: Math.max(0, rule.max_requests - estimatedCount - 1),
    resetAt: currentWindow + windowSize,
  };
}
```

```python
import time, redis

def sliding_window_counter_check(r: redis.Redis, key: str, max_requests: int, window_sec: int) -> dict:
    now = int(time.time())
    current_window = (now // window_sec) * window_sec
    previous_window = current_window - window_sec
    elapsed = now - current_window
    weight = 1 - (elapsed / window_sec)

    current_key = f"{key}:{current_window}"
    previous_key = f"{key}:{previous_window}"

    prev_count = int(r.get(previous_key) or 0)
    curr_count = int(r.get(current_key) or 0)
    estimated = int(prev_count * weight) + curr_count

    if estimated >= max_requests:
        return {"allowed": False, "remaining": 0, "reset_at": current_window + window_sec}

    pipe = r.pipeline()
    pipe.incr(current_key)
    pipe.expire(current_key, window_sec * 2 + 1)
    pipe.execute()

    return {
        "allowed": True,
        "remaining": max(0, max_requests - estimated - 1),
        "reset_at": current_window + window_sec,
    }
```

**Pros:**
- Good balance of accuracy and memory efficiency
- Only 2 counters per entity (vs. N timestamps in sliding log)
- Cloudflare uses a variant of this approach

**Cons:**
- Approximation, not exact — Cloudflare reports < 0.003% error rate in practice
- Slightly more complex than fixed window

---

### 7.4 Token Bucket

**How it works:** Each entity has a "bucket" with a maximum capacity. Tokens are added at a constant refill rate. Each request consumes one token. If the bucket is empty, the request is rejected. The bucket cannot exceed its maximum capacity.

```
Bucket Config: capacity = 10, refill_rate = 2 tokens/sec

Time 0s:   [##########]  10/10 tokens  → Request uses 1 → 9 tokens
Time 0.5s: [##########]  10/10 tokens  (refilled 1, capped at 10)
Time 1s:   [#########-]   9/10 tokens  → 3 requests → 6 tokens
Time 2s:   [########--]   8/10 tokens  (refilled 2)
Time 5s:   [##########]  10/10 tokens  (refilled 6, capped at 10)
Time 5s:   [----------]   0/10 tokens  → Burst: 10 requests at once!
Time 6s:   [##--------]   2/10 tokens  (refilled 2)

Key insight: Token bucket allows BURSTS up to bucket capacity.
```

**Implementation:**

```javascript
async function tokenBucketCheck(redis, key, rule) {
  const luaScript = `
    local key = KEYS[1]
    local capacity = tonumber(ARGV[1])
    local refill_rate = tonumber(ARGV[2])
    local now = tonumber(ARGV[3])
    local requested = tonumber(ARGV[4])

    local bucket = redis.call('hmget', key, 'tokens', 'last_refill')
    local tokens = tonumber(bucket[1])
    local last_refill = tonumber(bucket[2])

    -- Initialize bucket if it doesn't exist
    if tokens == nil then
      tokens = capacity
      last_refill = now
    end

    -- Calculate tokens to add since last refill
    local elapsed = (now - last_refill) / 1000
    tokens = math.min(capacity, tokens + (elapsed * refill_rate))

    -- Check if enough tokens
    local allowed = tokens >= requested
    if allowed then
      tokens = tokens - requested
    end

    -- Save state
    redis.call('hmset', key, 'tokens', tostring(tokens), 'last_refill', tostring(now))
    redis.call('expire', key, math.ceil(capacity / refill_rate) + 60)

    return {allowed and 1 or 0, math.floor(tokens), math.ceil((requested - tokens) / refill_rate * 1000)}
  `;

  const now = Date.now();
  const result = await redis.eval(
    luaScript, 1, key,
    rule.burst_size,        // capacity
    rule.refill_rate,       // tokens per second
    now,                    // current timestamp
    1                       // tokens requested
  );

  return {
    allowed: result[0] === 1,
    remaining: result[1],
    retryAfter: result[0] === 1 ? 0 : Math.ceil(result[2] / 1000),
    resetAt: Math.ceil(now / 1000) + Math.ceil(rule.burst_size / rule.refill_rate),
  };
}
```

```python
import time, redis

TOKEN_BUCKET_LUA = """
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local bucket = redis.call('hmget', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1]) or capacity
local last_refill = tonumber(bucket[2]) or now

local elapsed = (now - last_refill) / 1000
tokens = math.min(capacity, tokens + elapsed * refill_rate)

local allowed = 0
if tokens >= 1 then
    tokens = tokens - 1
    allowed = 1
end

redis.call('hmset', key, 'tokens', tostring(tokens), 'last_refill', tostring(now))
redis.call('expire', key, math.ceil(capacity / refill_rate) + 60)
return {allowed, math.floor(tokens)}
"""

def token_bucket_check(r: redis.Redis, key: str, capacity: int, refill_rate: float) -> dict:
    now_ms = int(time.time() * 1000)
    result = r.eval(TOKEN_BUCKET_LUA, 1, key, capacity, refill_rate, now_ms)

    return {
        "allowed": result[0] == 1,
        "remaining": result[1],
        "reset_at": int(time.time()) + int(capacity / refill_rate),
    }
```

**Pros:**
- Allows controlled bursts (important for real-world APIs like Stripe)
- Smooth average rate over time
- Used by AWS and Stripe

**Cons:**
- Two parameters to tune (capacity and refill rate)
- Bursts can overwhelm downstream services if capacity is set too high

---

### 7.5 Leaky Bucket

**How it works:** Requests enter a FIFO queue (the "bucket") and are processed at a constant rate. If the queue is full, incoming requests are dropped. The output rate is always smooth regardless of input burstiness.

```
Input (bursty):    ████  ██     ████████   ██
                    ↓     ↓        ↓        ↓
               ┌──────────────────────────────┐
               │          QUEUE (Bucket)       │
               │  Capacity: 10 requests        │
               │  [req][req][req][req][ ][ ]   │
               └──────────────┬────────────────┘
                              │ drain rate: 2 req/sec
                              ▼
Output (smooth):   ██  ██  ██  ██  ██  ██  ██

If queue is FULL → new requests are DROPPED (HTTP 429)
```

**Implementation:**

```javascript
async function leakyBucketCheck(redis, key, rule) {
  const luaScript = `
    local key = KEYS[1]
    local capacity = tonumber(ARGV[1])
    local leak_rate = tonumber(ARGV[2])  -- requests drained per second
    local now = tonumber(ARGV[3])

    local bucket = redis.call('hmget', key, 'queue_size', 'last_leak')
    local queue_size = tonumber(bucket[1]) or 0
    local last_leak = tonumber(bucket[2]) or now

    -- Drain leaked requests
    local elapsed = (now - last_leak) / 1000
    local leaked = math.floor(elapsed * leak_rate)
    queue_size = math.max(0, queue_size - leaked)

    if leaked > 0 then
      last_leak = now
    end

    -- Try to add to queue
    local allowed = queue_size < capacity
    if allowed then
      queue_size = queue_size + 1
    end

    redis.call('hmset', key, 'queue_size', tostring(queue_size), 'last_leak', tostring(now))
    redis.call('expire', key, math.ceil(capacity / leak_rate) + 60)

    return {allowed and 1 or 0, capacity - queue_size}
  `;

  const now = Date.now();
  const result = await redis.eval(
    luaScript, 1, key,
    rule.burst_size,        // queue capacity
    rule.refill_rate,       // leak rate (reusing field)
    now
  );

  return {
    allowed: result[0] === 1,
    remaining: result[1],
    resetAt: Math.ceil(now / 1000) + Math.ceil(1 / rule.refill_rate),
  };
}
```

```python
LEAKY_BUCKET_LUA = """
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local leak_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local bucket = redis.call('hmget', key, 'queue_size', 'last_leak')
local queue_size = tonumber(bucket[1]) or 0
local last_leak = tonumber(bucket[2]) or now

local elapsed = (now - last_leak) / 1000
local leaked = math.floor(elapsed * leak_rate)
queue_size = math.max(0, queue_size - leaked)
if leaked > 0 then last_leak = now end

local allowed = 0
if queue_size < capacity then
    queue_size = queue_size + 1
    allowed = 1
end

redis.call('hmset', key, 'queue_size', tostring(queue_size), 'last_leak', tostring(now))
redis.call('expire', key, math.ceil(capacity / leak_rate) + 60)
return {allowed, capacity - queue_size}
"""

def leaky_bucket_check(r: redis.Redis, key: str, capacity: int, leak_rate: float) -> dict:
    now_ms = int(time.time() * 1000)
    result = r.eval(LEAKY_BUCKET_LUA, 1, key, capacity, leak_rate, now_ms)

    return {
        "allowed": result[0] == 1,
        "remaining": result[1],
        "reset_at": int(time.time()) + int(1 / leak_rate),
    }
```

**Pros:**
- Guarantees smooth, constant output rate — ideal for protecting downstream
- Simple mental model (FIFO queue with constant drain)

**Cons:**
- Does not allow bursts — even legitimate spikes are queued or dropped
- Queue delay can introduce latency for requests that are waiting

---

### 7.6 Algorithm Comparison

| Algorithm | Accuracy | Memory | Burst Handling | Complexity | Best For |
|-----------|----------|--------|---------------|------------|----------|
| **Fixed Window** | ⚠️ Boundary burst | 🟢 Very low | ⚠️ 2× burst at edge | 🟢 Simple | Internal APIs, non-critical limits |
| **Sliding Window Log** | 🟢 Exact | 🔴 High | 🟢 No bursts | 🟡 Medium | Financial APIs, strict compliance |
| **Sliding Window Counter** | 🟡 ~99.97% accurate | 🟢 Low | 🟡 Minimal burst | 🟡 Medium | **General purpose (recommended)** |
| **Token Bucket** | 🟢 Exact | 🟢 Low | 🟢 Controlled bursts | 🟡 Medium | Public APIs, user-facing rate limits |
| **Leaky Bucket** | 🟢 Exact | 🟢 Low | 🔴 No bursts allowed | 🟡 Medium | Protecting fragile downstream services |

### 7.7 When to Use Which

**Candidate:** "My recommendations based on the use case:

- **Cloudflare-style CDN/WAF:** Sliding Window Counter — high throughput, low memory, accurate enough
- **Stripe-style API:** Token Bucket — allows legitimate bursts (e.g., batch API calls), configurable
- **Payment processing / financial:** Sliding Window Log — exactness matters more than memory
- **Message queues / worker protection:** Leaky Bucket — smooth the load on downstream consumers
- **Quick prototype / internal tools:** Fixed Window — simplicity wins when boundary bursts are acceptable"

---

## 8. Scaling Strategies

### 8.1 Race Conditions in Distributed Counters

**Candidate:** "The biggest challenge in distributed rate limiting is race conditions. If two requests from the same user hit different nodes simultaneously, they might both read counter=99 (limit=100), both increment to 100, and both get allowed — exceeding the limit."

**The Problem:**
```
Node A reads counter = 99       Node B reads counter = 99
Node A: 99 < 100 → ALLOW       Node B: 99 < 100 → ALLOW
Node A writes counter = 100     Node B writes counter = 100
                                
Result: Both allowed, but only 1 slot was available!
```

**Solution: Redis Lua Script for Atomic Check-and-Increment**

```lua
-- Atomic rate limit check: no race conditions possible
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local current = tonumber(redis.call('GET', key) or "0")

if current >= limit then
    return {0, current, tonumber(redis.call('TTL', key))}  -- rejected
end

current = redis.call('INCR', key)
if current == 1 then
    redis.call('EXPIRE', key, window)
end

return {1, limit - current, tonumber(redis.call('TTL', key))}  -- allowed
```

**Why this works:** Lua scripts in Redis execute atomically — no other command can interleave. The read-check-write is a single atomic operation.

### 8.2 Synchronization Across Data Centers

```
               ┌────────────────────────────┐
               │       Global Aggregator     │
               │    (Async sync every 5s)    │
               └───────┬────────────┬────────┘
                       │            │
          ┌────────────▼──┐    ┌───▼─────────────┐
          │  US-East DC   │    │   EU-West DC     │
          │               │    │                  │
          │ Redis Cluster │◄──►│  Redis Cluster   │
          │ Local counter │    │  Local counter   │
          │ Limit: 80/min │    │  Limit: 80/min   │
          │ (80% of 100)  │    │  (80% of 100)    │
          └───────────────┘    └──────────────────┘

Strategy: Split global limit across DCs with a buffer.
- Global limit: 100 req/min
- Each DC gets 80% of the limit (allows overlap headroom)
- Background sync reconciles every 5 seconds
- Trade-off: may allow up to ~160% of limit briefly (acceptable per NFR)
```

**Approaches:**

| Strategy | Consistency | Latency | Complexity |
|----------|------------|---------|------------|
| Centralized Redis | Strong | High (cross-DC RTT ~50-100ms) | Low |
| Local counters + async sync | Eventual | Low (~1ms) | Medium |
| Split limit per DC | Approximate | Low | Medium |
| CRDT-based counters | Eventual | Low | High |

**Candidate:** "For our 5ms latency requirement, centralized Redis across data centers won't work. I'd use the **split limit approach** with periodic sync. Each DC independently enforces a fraction of the global limit and periodically shares its counts with other DCs via an async aggregator."

### 8.3 Local + Global Rate Limiting Hybrid

```javascript
async function hybridRateLimit(req, rule) {
  const entityKey = extractEntityKey(req);

  // Layer 1: In-process local counter (zero network latency)
  const localResult = localRateLimiter.check(entityKey, {
    max: Math.floor(rule.max_requests / CLUSTER_SIZE),
    window: rule.window_seconds,
  });

  if (!localResult.allowed) {
    return { allowed: false, source: 'local' };
  }

  // Layer 2: Distributed Redis counter (shared state)
  const globalResult = await redisRateLimiter.check(entityKey, rule);
  return { allowed: globalResult.allowed, source: 'global' };
}
```

**Benefits:**
- Local limiter catches obvious over-limit cases instantly (no Redis call)
- Reduces Redis load by 60-70% in practice
- Redis is only consulted when local limit isn't exceeded

### 8.4 Consistent Hashing for Counter Distribution

```
                      Hash Ring (Redis Cluster)

                         Node A
                        ╱      ╲
                   Node F        Node B
                   │                  │
                   Node E        Node C
                        ╲      ╱
                         Node D

    Key: "rl:api_key:sk_abc:charges"
      → hash("rl:api_key:sk_abc:charges") → lands on Node B
    
    All counters for sk_abc's /charges endpoint go to Node B.
    Redis Cluster handles this via hash slots (16,384 slots).
```

**Candidate:** "Redis Cluster already uses hash slots for key distribution. We can use hash tags like `{sk_abc}` to ensure all counters for the same entity land on the same shard, enabling atomic multi-key operations."

```
Key format: rl:{entity_id}:{endpoint}:{window}
Redis hash tag ensures rl:{sk_abc}:charges:172400 goes to same slot as rl:{sk_abc}:charges:172340
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Redis Cluster Failure

**Candidate:** "This is the most critical failure scenario. The rate limiter depends on Redis for state. If Redis is down, we have two options:"

| Strategy | Behavior | Risk | When to Use |
|----------|----------|------|-------------|
| **Fail-Open** | Allow all traffic | Over-limit abuse | User-facing APIs (protect availability) |
| **Fail-Closed** | Block all traffic | Legitimate users denied | Security-critical (DDoS, auth) |
| **Degraded** | Fall back to local in-memory limits | Approximate enforcement | Balanced approach (recommended) |

```javascript
async function resilientRateLimit(entityKey, rule) {
  try {
    return await redisRateLimiter.check(entityKey, rule);
  } catch (error) {
    metrics.increment('rate_limiter.redis_failure');

    if (rule.action === 'security') {
      // Fail-closed for security-critical rules (DDoS, brute force)
      return { allowed: false, remaining: 0, source: 'fail_closed' };
    }

    // Fail-open with degraded local enforcement for everything else
    return localFallback.check(entityKey, {
      max: Math.floor(rule.max_requests / CLUSTER_SIZE),
      window: rule.window_seconds,
    });
  }
}
```

### 9.2 Clock Drift Across Nodes

**Problem:** If two nodes have clocks that differ by 2 seconds, they'll compute different window boundaries, leading to inconsistent counting.

**Mitigations:**
1. **NTP synchronization** — all nodes sync clocks via NTP (typical drift < 10ms)
2. **Server-side timestamps only** — never trust client timestamps
3. **Redis as the time authority** — use `redis.call('TIME')` in Lua scripts instead of local clock

```lua
-- Use Redis server time instead of client time
local time = redis.call('TIME')
local now_ms = tonumber(time[1]) * 1000 + math.floor(tonumber(time[2]) / 1000)
```

### 9.3 Network Partition Handling

```
Scenario: Network partition splits DC into two halves

   [Gateway A] ─── [Redis Master]      PARTITION       [Gateway B]
                                    ═══════════════
                                                        [Redis Replica]

Gateway B can't reach Redis Master.
```

**Solution:** Redis Sentinel / Cluster auto-failover promotes the replica. During failover (~10-30s), gateways in partition B use local fallback.

### 9.4 Thundering Herd After Limit Reset

**Problem:** If 10,000 clients are all rate-limited and waiting for the window to reset at exactly `T=60s`, they all retry simultaneously at `T=60.001s`, causing a spike.

**Mitigations:**

```javascript
function calculateRetryAfter(resetAt) {
  const baseDelay = resetAt - Math.floor(Date.now() / 1000);
  const jitter = Math.random() * 5; // 0-5 seconds of random jitter
  return Math.max(1, baseDelay + jitter);
}
```

- Add random jitter (0-5s) to `Retry-After` headers
- Stagger window boundaries per entity (hash-based offset)
- Use token bucket (gradual refill) instead of fixed window (cliff-edge reset)

### 9.5 Counter Overflow

**Problem:** Integer overflow if a counter accumulates beyond `2^53` (JavaScript) or `2^63` (Redis).

**Mitigations:**
- Redis INCR supports 64-bit signed integers (max: 9.2 × 10^18) — practically unreachable
- TTL-based auto-expiry ensures counters never accumulate indefinitely
- Monitor for anomalous counter values as a canary for bugs

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|----------------|
| `rate_limiter.check.latency_ms` | Time to perform rate check | p99 > 5ms |
| `rate_limiter.requests.allowed` | Count of allowed requests | Sudden drop > 50% |
| `rate_limiter.requests.rejected` | Count of rejected requests | Spike > 10% of traffic |
| `rate_limiter.redis.errors` | Redis connectivity failures | Any sustained errors |
| `rate_limiter.fallback.active` | Using local fallback | True for > 30s |
| `rate_limiter.false_positive.rate` | Legitimate requests wrongly blocked | > 0.1% |
| `rate_limiter.rules.active_count` | Number of active rules | Change > 10% |

### 10.2 Grafana Dashboard Layout

```
┌──────────────────────────────────────────────────────────────────────┐
│  RATE LIMITER DASHBOARD                                   [24h ▼]   │
├──────────────────────┬───────────────────────────────────────────────┤
│  Allowed Requests/s  │  Rejected Requests/s                         │
│  ████████████████    │  █                                           │
│  ████████████████    │  ██                                          │
│  ████████████████    │  █                                           │
│  956,234 /s          │  43,766 /s  (4.4% reject rate)              │
├──────────────────────┼───────────────────────────────────────────────┤
│  Check Latency (ms)  │  Redis Operations/s                          │
│  p50: 0.8ms          │  ██████████████████████                      │
│  p95: 2.1ms          │  ██████████████████████                      │
│  p99: 3.4ms          │  1,842,000 ops/s                             │
├──────────────────────┼───────────────────────────────────────────────┤
│  Top Throttled Keys  │  Rule Hit Distribution                       │
│  1. ip:203.0.113.42  │  api_key rules: 72%                         │
│  2. user:u_malicious │  ip rules:      18%                         │
│  3. api_key:sk_free1 │  user rules:    10%                         │
├──────────────────────┴───────────────────────────────────────────────┤
│  Throttle Rate by Rule (Top 10)                                      │
│  ┌─────────────────────────────────────────────────────────────────┐ │
│  │  free_tier_limit    ████████████████████████  12,340/min       │ │
│  │  ip_global_limit    ██████████                 5,120/min       │ │
│  │  auth_brute_force   ████                       2,100/min       │ │
│  └─────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────┘
```

### 10.3 Prometheus Alerting Rules

```yaml
groups:
  - name: rate_limiter_alerts
    rules:
      - alert: RateLimiterHighLatency
        expr: histogram_quantile(0.99, rate(rate_limiter_check_duration_seconds_bucket[5m])) > 0.005
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Rate limiter p99 latency exceeds 5ms"
          description: "p99 latency is {{ $value }}s over the last 5 minutes"

      - alert: RateLimiterHighRejectRate
        expr: |
          rate(rate_limiter_requests_total{result="rejected"}[5m])
          /
          rate(rate_limiter_requests_total[5m]) > 0.15
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Rate limiter rejecting >15% of traffic"

      - alert: RateLimiterRedisDown
        expr: rate_limiter_redis_errors_total > 0
        for: 30s
        labels:
          severity: critical
        annotations:
          summary: "Rate limiter cannot reach Redis cluster"

      - alert: RateLimiterFallbackActive
        expr: rate_limiter_fallback_active == 1
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Rate limiter using local fallback (Redis unavailable)"
```

### 10.4 Structured Logging

```javascript
function logThrottledRequest(req, rule, result) {
  logger.warn({
    event: 'rate_limit_exceeded',
    entity_type: rule.entity_type,
    entity_id: extractEntityId(req),
    endpoint: req.path,
    method: req.method,
    rule_id: rule.rule_id,
    rule_name: rule.name,
    current_count: result.currentCount,
    limit: rule.max_requests,
    window_seconds: rule.window_seconds,
    source_ip: req.ip,
    user_agent: req.headers['user-agent'],
    timestamp: new Date().toISOString(),
  });
}
```

---

## 11. Advanced Features

### 11.1 Rate Limit Response Headers

**Standard headers** (used by GitHub, Stripe, Twitter):

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100          # Maximum requests allowed in window
X-RateLimit-Remaining: 0        # Requests remaining in current window
X-RateLimit-Reset: 1714003260   # Unix timestamp when the window resets
Retry-After: 12                 # Seconds until the client should retry
Content-Type: application/json

{
  "error": {
    "type": "rate_limit_error",
    "code": "rate_limit_exceeded",
    "message": "Too many requests. Please retry after 12 seconds.",
    "documentation_url": "https://docs.api.com/rate-limiting"
  }
}
```

**Implementation:**

```javascript
function setRateLimitHeaders(res, rule, result) {
  res.set('X-RateLimit-Limit', String(rule.max_requests));
  res.set('X-RateLimit-Remaining', String(Math.max(0, result.remaining)));
  res.set('X-RateLimit-Reset', String(result.resetAt));

  if (!result.allowed) {
    const retryAfter = Math.max(1, result.resetAt - Math.floor(Date.now() / 1000));
    res.set('Retry-After', String(retryAfter));
  }
}
```

### 11.2 Adaptive Rate Limiting

**Concept:** Dynamically adjust rate limits based on system health. When backend services are under stress, tighten limits; when healthy, relax them.

```javascript
class AdaptiveRateLimiter {
  constructor(baseLimit, healthChecker) {
    this.baseLimit = baseLimit;
    this.healthChecker = healthChecker;
    this.multiplier = 1.0;

    setInterval(() => this.adjustLimits(), 10_000);
  }

  adjustLimits() {
    const health = this.healthChecker.getScore(); // 0.0 (unhealthy) to 1.0 (healthy)

    if (health < 0.3) {
      this.multiplier = 0.25; // severely reduce limits
    } else if (health < 0.6) {
      this.multiplier = 0.5;
    } else if (health < 0.8) {
      this.multiplier = 0.75;
    } else {
      this.multiplier = 1.0;
    }
  }

  getEffectiveLimit() {
    return Math.max(1, Math.floor(this.baseLimit * this.multiplier));
  }
}

class HealthChecker {
  getScore() {
    const cpuUsage = os.loadavg()[0] / os.cpus().length;
    const errorRate = metrics.get('http_5xx_rate_1m');
    const p99Latency = metrics.get('http_latency_p99_ms');

    let score = 1.0;
    if (cpuUsage > 0.8) score -= 0.3;
    if (errorRate > 0.05) score -= 0.3;
    if (p99Latency > 500) score -= 0.2;

    return Math.max(0, Math.min(1, score));
  }
}
```

### 11.3 Rate Limiting by Cost / Weight

**Concept:** Not all API calls are equal. A `GET /users` is cheap, but a `POST /reports/generate` is expensive. Assign weights to endpoints.

```javascript
const ENDPOINT_COSTS = {
  'GET /v1/users':           1,
  'GET /v1/users/:id':       1,
  'POST /v1/charges':        5,
  'POST /v1/reports':       20,
  'POST /v1/bulk-import':   50,
};

async function weightedRateLimit(req, rule) {
  const cost = ENDPOINT_COSTS[`${req.method} ${req.route.path}`] || 1;
  const entityKey = extractEntityKey(req);

  // Token bucket: consume `cost` tokens instead of 1
  return tokenBucketCheck(redis, entityKey, {
    ...rule,
    tokensRequested: cost,
  });
}
```

**Example:** With a limit of 100 tokens/minute:
- User can make 100 `GET /users` calls (1 token each), or
- 20 `POST /charges` calls (5 tokens each), or
- 5 `POST /reports` calls (20 tokens each)

### 11.4 Distributed Rate Limiting with Gossip Protocol

**Concept:** Instead of centralized Redis, nodes share their local counts via a gossip protocol. Each node maintains approximate global state.

```
Node A: local_count = 30            Node B: local_count = 25
    │                                    │
    └─── gossip(A=30) ──────────────────►│
    │◄── gossip(B=25) ──────────────────┘│
    │                                    │
    │ global_estimate = 30 + 25 = 55     │ global_estimate = 30 + 25 = 55
    │ limit = 100 → ALLOW               │ limit = 100 → ALLOW

Gossip interval: every 1-2 seconds
Trade-off: eventual consistency (up to gossip_interval × node_count overcount)
```

**When to use:** Ultra-low-latency requirements where even a Redis round-trip is too slow, or in environments without centralized state stores (edge computing, CDN nodes).

### 11.5 Client-Side Rate Limiting

**Concept:** Well-behaved clients implement rate limiting locally to avoid wasting requests and receiving 429 errors.

```javascript
class ClientRateLimiter {
  constructor() {
    this.tokens = 100;
    this.maxTokens = 100;
    this.refillRate = 100 / 60; // 100 per minute
    this.lastRefill = Date.now();
    this.queue = [];
  }

  async request(url, options) {
    this.refill();

    if (this.tokens < 1) {
      const waitMs = ((1 - this.tokens) / this.refillRate) * 1000;
      await new Promise(resolve => setTimeout(resolve, waitMs));
      this.refill();
    }

    this.tokens -= 1;
    const response = await fetch(url, options);

    // Sync with server's view
    if (response.headers.has('X-RateLimit-Remaining')) {
      const serverRemaining = parseInt(response.headers.get('X-RateLimit-Remaining'));
      this.tokens = Math.min(this.tokens, serverRemaining);
    }

    if (response.status === 429) {
      const retryAfter = parseInt(response.headers.get('Retry-After') || '1');
      await new Promise(resolve => setTimeout(resolve, retryAfter * 1000));
      return this.request(url, options); // retry
    }

    return response;
  }

  refill() {
    const now = Date.now();
    const elapsed = (now - this.lastRefill) / 1000;
    this.tokens = Math.min(this.maxTokens, this.tokens + elapsed * this.refillRate);
    this.lastRefill = now;
  }
}
```

---

## 12. Interview Q&A

### Q1: How do you handle rate limiting in a multi-region setup?

**Candidate:** "There are three main strategies, each trading off accuracy for latency:

1. **Centralized counter (strong consistency):** All regions query a single Redis cluster. Accurate but adds cross-region latency (~50-100ms). Only works if latency budget allows it.

2. **Split limits per region (partitioned):** Divide the global limit proportionally across regions (e.g., 60% US, 30% EU, 10% APAC). Each region enforces independently. Simple but static — doesn't adapt to shifting traffic patterns.

3. **Local counters + async sync (recommended):** Each region maintains local Redis counters and enforces a local limit. A background job syncs counts across regions every 1-5 seconds. This gives sub-millisecond local latency with approximate global enforcement. The error margin is bounded by `sync_interval × max_rate_per_region`.

I'd choose option 3 for our use case since we need sub-5ms latency and accept < 1% false negatives."

---

### Q2: Fixed window vs sliding window — when would you choose each?

**Candidate:** "Fixed window is appropriate when:
- The rate limit is not mission-critical (internal APIs, monitoring endpoints)
- Simplicity and performance are the top priorities
- A 2× boundary burst is acceptable

Sliding window (counter variant) is better when:
- You're building a public-facing API with SLAs (like Stripe or GitHub)
- You need predictable, consistent enforcement
- Users are sophisticated enough to exploit boundary bursts

The sliding window counter is my default recommendation because it's nearly as simple as fixed window but eliminates 99.97% of the boundary burst issue, at the cost of tracking one additional counter."

---

### Q3: How do you prevent race conditions in distributed counters?

**Candidate:** "The core technique is **atomic operations**. In Redis, I use Lua scripts that bundle the read-check-increment into a single atomic execution:

```lua
local current = tonumber(redis.call('GET', KEYS[1]) or '0')
if current >= tonumber(ARGV[1]) then
    return 0  -- rejected
end
redis.call('INCR', KEYS[1])
return 1  -- allowed
```

This runs atomically on the Redis server — no interleaving is possible. For additional defense:
- Use Redis Cluster to partition keys, reducing contention per shard
- Batch local counts and flush to Redis periodically (trades accuracy for throughput)
- Use `MULTI/EXEC` transactions for simpler operations without Lua"

---

### Q4: What happens when Redis goes down? Fail-open or fail-closed?

**Candidate:** "This is a critical design decision that depends on the rule's purpose:

**Fail-open (allow all traffic):**
- Use for general API rate limits (e.g., 'free tier gets 1000 req/hour')
- Protects user experience — legitimate traffic is never blocked
- Risk: temporary abuse during outage

**Fail-closed (block all traffic):**
- Use for security-critical limits (brute-force login, DDoS mitigation)
- Protects system integrity
- Risk: blocks legitimate users during outage

**My approach:** Tag each rule with a `failure_mode` field. Security rules fail-closed; business rules fail-open. Additionally, fall back to in-memory local counters as a degraded middle ground. This local fallback enforces `limit / node_count` per node — imprecise, but better than nothing."

---

### Q5: How would you implement tiered rate limits (free vs premium)?

**Candidate:** "I'd use the rules engine with a `tier` field:

```sql
-- Free tier: 60 req/min for all endpoints
INSERT INTO rate_limit_rules (entity_type, tier, max_requests, window_seconds)
VALUES ('api_key', 'free', 60, 60);

-- Pro tier: 1000 req/min
INSERT INTO rate_limit_rules (entity_type, tier, max_requests, window_seconds)
VALUES ('api_key', 'pro', 1000, 60);

-- Enterprise tier: 10000 req/min
INSERT INTO rate_limit_rules (entity_type, tier, max_requests, window_seconds)
VALUES ('api_key', 'enterprise', 10000, 60);
```

At request time, the middleware resolves the user's tier from the API key metadata (cached) and selects the matching rule. This is how Stripe and GitHub implement it — your API key carries your tier, and different tiers map to different limits.

For per-entity overrides (e.g., a free-tier user who got a special exception), I'd use the `rate_limit_overrides` table which takes priority over the tier-based rule."

---

### Q6: How do you handle clock synchronization across nodes?

**Candidate:** "Clock drift is a real issue in distributed systems. Three complementary strategies:

1. **NTP synchronization:** All servers run NTP (or chrony) to synchronize with time servers. Typical drift is < 10ms, which is acceptable for second-granularity windows.

2. **Single time authority:** Use the Redis server's clock (`redis.call('TIME')`) in Lua scripts instead of the application server's clock. Since the counter and the time check both happen on the same Redis node, drift is irrelevant.

3. **Window overlap tolerance:** Design windows with small overlap buffers. For a 60-second window, a 100ms clock drift means at most 0.17% error — well within our < 1% tolerance.

The practical answer: NTP + Redis server time solves this for 99.99% of cases."

---

### Q7: How would you rate limit WebSocket connections?

**Candidate:** "WebSocket rate limiting is different from HTTP because it's a long-lived connection. I'd apply limits at three layers:

1. **Connection rate:** Limit new WebSocket handshakes per user (e.g., 10 connections/min) — prevents connection flooding. Use a standard token bucket.

2. **Message rate:** Limit messages per connection (e.g., 100 messages/sec). This is enforced in the WebSocket handler itself using an in-memory token bucket per connection.

3. **Bandwidth rate:** Limit total bytes per second. Track cumulative payload size and throttle when the byte budget is exceeded.

```javascript
class WebSocketRateLimiter {
  constructor(messagesPerSec, bytesPerSec) {
    this.msgBucket = new TokenBucket(messagesPerSec, messagesPerSec);
    this.byteBucket = new TokenBucket(bytesPerSec, bytesPerSec);
  }

  allowMessage(messageBytes) {
    if (!this.msgBucket.consume(1)) return false;
    if (!this.byteBucket.consume(messageBytes)) return false;
    return true;
  }
}
```

Key difference: since WebSocket connections are stateful and pinned to a server, local in-memory rate limiting is efficient and sufficient."

---

### Q8: How do you test a rate limiter in production?

**Candidate:** "Testing rate limiters requires care because incorrect limits can either break legitimate users or fail to protect the system. My approach:

1. **Shadow mode / log-only:** Deploy new rules with `action: 'log_only'`. They evaluate but don't enforce — just emit metrics. Compare predicted rejections against expectations for 24-48 hours.

2. **Canary deployment:** Roll out enforcement to 5% of traffic first. Monitor reject rate, false positives, and user complaints. Gradually increase to 100%.

3. **Synthetic load testing:** Use a dedicated test API key with a known limit. Continuously send requests at exactly the limit boundary and verify the limiter behaves correctly:

```python
import asyncio, aiohttp, time

async def test_rate_limit(url, api_key, expected_limit, window_sec):
    async with aiohttp.ClientSession() as session:
        results = {'allowed': 0, 'rejected': 0}
        start = time.time()
        
        for i in range(expected_limit + 20):
            async with session.get(url, headers={'X-API-Key': api_key}) as resp:
                if resp.status == 200:
                    results['allowed'] += 1
                elif resp.status == 429:
                    results['rejected'] += 1

        elapsed = time.time() - start
        assert results['allowed'] == expected_limit, f"Expected {expected_limit}, got {results['allowed']}"
        assert results['rejected'] == 20, f"Expected 20 rejections, got {results['rejected']}"
        print(f"PASS: {results['allowed']} allowed, {results['rejected']} rejected in {elapsed:.2f}s")

asyncio.run(test_rate_limit('https://api.example.com/v1/test', 'sk_test_abc', 100, 60))
```

4. **Chaos engineering:** Intentionally kill Redis nodes and verify graceful degradation to local fallback. Introduce artificial clock drift and confirm the system remains within tolerance.

5. **A/B testing limits:** For business rate limits (not security), run A/B tests with different thresholds to find the optimal balance between protection and user experience."

---

## 13. Production Checklist

### Pre-Launch

- [ ] Redis Cluster deployed with 3+ masters and replicas
- [ ] Lua scripts tested for atomicity under concurrent load
- [ ] Graceful degradation (fail-open/fail-closed) tested per rule category
- [ ] Rate limit headers (`X-RateLimit-*`) returning correctly
- [ ] Rule CRUD API secured (admin-only auth)
- [ ] Load test at 2× expected peak (2M req/sec) — verify < 5ms p99 latency
- [ ] Client documentation published (limits, headers, retry behavior)

### Day-1

- [ ] Deploy all rules in `log_only` mode — no enforcement
- [ ] Monitor throttle rate vs expectations from shadow data
- [ ] Verify Grafana dashboard and alerts are firing correctly
- [ ] Confirm Redis memory usage is within provisioned capacity
- [ ] Test 429 response body and headers match API specification

### Week-1

- [ ] Enable enforcement for security-critical rules first (brute-force, DDoS)
- [ ] Canary enforce business rules on 5% → 25% → 50% → 100% of traffic
- [ ] Investigate any false positive reports from users
- [ ] Tune adaptive rate limiting thresholds based on real load patterns
- [ ] Verify counter TTLs are expiring — no memory leak in Redis

### Month-1

- [ ] Review throttle rate analytics — identify rules that need tuning
- [ ] Implement multi-region sync if expanding to additional data centers
- [ ] Set up automated regression tests for rate limit accuracy
- [ ] Create runbook for Redis cluster failure scenarios
- [ ] Benchmark counter growth and project Redis capacity for next 6 months
- [ ] Review audit logs for rule changes — ensure change management process works

---

## Summary

### Technical Decision Table

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Counter store | Redis Cluster | Sub-ms latency, atomic Lua scripts, TTL, sorted sets |
| Default algorithm | Sliding Window Counter | Best accuracy-to-memory ratio, Cloudflare-proven |
| Rule storage | PostgreSQL | Durable, relational queries, low write volume |
| Rule caching | Local + Redis (30s TTL) | Minimizes DB reads, pub/sub invalidation |
| Deployment model | API Gateway plugin | Lowest latency, no extra network hop |
| Failure mode | Per-rule (fail-open / fail-closed) | Security rules strict, business rules lenient |
| Multi-region | Local counters + async sync | Sub-5ms latency with approximate global enforcement |
| Clock strategy | Redis server time via Lua | Eliminates cross-node drift |
| Race condition prevention | Lua atomic scripts | Read-check-write in single atomic op |

### Scalability Path

```
Phase 1 (< 100K req/s):
  Single Redis instance, gateway middleware, fixed window
  ↓
Phase 2 (100K - 1M req/s):
  Redis Cluster (6 nodes), sliding window counter, rule caching
  ↓
Phase 3 (1M - 10M req/s):
  Multi-region Redis, local + global hybrid, adaptive limits
  ↓
Phase 4 (10M+ req/s):
  Edge-based rate limiting (Cloudflare Workers), gossip protocol,
  CRDT counters, per-PoP enforcement with global aggregation
```

---

> **Interview Tip:** When designing a rate limiter, always start by asking: "What entity are we limiting? What's the limit? What's the window?" Then discuss the algorithm trade-offs. Interviewers care more about your reasoning through trade-offs (consistency vs latency, memory vs accuracy) than memorizing a single "correct" answer. The best candidates explicitly state: "Here's what I'm trading off, and here's why."
