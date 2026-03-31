# Phase 1: Core Building Blocks

**Navigation:** [← Previous: Foundations](00-foundations.md) | [Next: Distributed Systems →](02-distributed-systems.md)

---

These are the essential components you'll use in every system design. Master these before moving to distributed systems concepts. Each component represents a tool in your architect's toolkit.

---

# Section 1: Networking & API Design

## 1.1 REST vs gRPC

### Concept Overview (What & Why)

**REST (Representational State Transfer):**

- Architectural style using HTTP methods (GET, POST, PUT, DELETE)
- Text-based (usually JSON), human-readable
- Stateless, resource-oriented
- Universal browser/client support

**gRPC (Google Remote Procedure Call):**

- Framework using Protocol Buffers (binary serialization)
- HTTP/2 based (multiplexing, streaming)
- Strongly typed via .proto definitions
- Better performance, smaller payloads

**Why this matters:**

- API choice affects latency, developer experience, and system capabilities
- Wrong choice can lead to painful migrations
- Most systems use both: REST for public APIs, gRPC for internal services

**When interviewers expect this:**

- Designing microservices communication
- Public API design
- Performance-sensitive systems
- Mobile applications

### Key Design Principles


| Aspect          | REST                   | gRPC                      |
| --------------- | ---------------------- | ------------------------- |
| Protocol        | HTTP/1.1 (usually)     | HTTP/2                    |
| Payload         | JSON/XML (text)        | Protocol Buffers (binary) |
| Typing          | Loose, schema optional | Strong, .proto required   |
| Streaming       | Polling/WebSocket      | Native bidirectional      |
| Browser support | Native                 | Requires proxy (grpc-web) |
| Tooling         | Ubiquitous             | Growing but less mature   |
| Human readable  | Yes                    | No (binary)               |


**REST Best Practices:**

- Use nouns for resources (`/users`, `/orders`)
- Use HTTP verbs correctly (GET = read, POST = create, PUT = update, DELETE = delete)
- Version your API (`/v1/users`)
- Return appropriate status codes (200, 201, 400, 404, 500)

**gRPC Best Practices:**

- Define clear .proto contracts
- Use deadlines (not timeouts)
- Implement proper error handling with status codes
- Consider backward compatibility in proto evolution

### Trade-offs & Decision Matrix


| Criterion              | Choose REST           | Choose gRPC                |
| ---------------------- | --------------------- | -------------------------- |
| Public API             | ✅ (universal support) | ❌ (limited browser)        |
| Mobile clients         | ✅ (simpler)           | ⚠️ (smaller payloads help) |
| Internal microservices | ⚠️ (works fine)       | ✅ (performance)            |
| Streaming needed       | ❌ (awkward)           | ✅ (native support)         |
| Developer familiarity  | ✅ (everyone knows)    | ⚠️ (learning curve)        |
| Performance critical   | ❌ (JSON overhead)     | ✅ (binary, HTTP/2)         |
| Polyglot services      | ✅ (ubiquitous)        | ✅ (code generation)        |


### Real-World Examples

**REST Usage:**

- Public APIs (Twitter, GitHub, Stripe)
- Web applications
- Third-party integrations
- CRUD applications

**gRPC Usage:**

- Google internal services
- Netflix microservices
- Kubernetes components
- Real-time gaming backends

**Hybrid Approach (Common):**

```
[Mobile/Web] → REST → [API Gateway] → gRPC → [Internal Services]
```

- Public surface: REST for compatibility
- Internal: gRPC for performance
- API Gateway handles translation

### Failure Scenarios & Edge Cases


| Scenario         | REST                       | gRPC                    |
| ---------------- | -------------------------- | ----------------------- |
| Client timeout   | Retry with idempotency key | Deadline propagation    |
| Version mismatch | 404 or parse error         | Proto evolution handles |
| Large payload    | May timeout                | Streaming handles       |
| Connection drop  | New TCP connection         | Connection kept alive   |


**Common Issues:**

*REST:*

- Over-fetching (getting more data than needed)
- Under-fetching (N+1 queries)
- No standard error format

*gRPC:*

- Debugging binary traffic is harder
- Browser support requires grpc-web proxy
- Load balancer must support HTTP/2

### Interview Perspective

**What interviewers look for:**

- Appropriate choice for the context
- Understanding of trade-offs
- Knowledge of when to use each

**Common traps:**

- ❌ "gRPC is always better" (not for public APIs)
- ❌ "REST is outdated" (it's still the right choice for many cases)
- ❌ Ignoring browser support constraints

**Strong signals:**

- ✅ "For public API, REST for compatibility; internally, gRPC for performance"
- ✅ "We'll use gRPC streaming for real-time updates"
- ✅ "REST for CRUD operations, gRPC for high-throughput service-to-service"

**Follow-up questions:**

- "How would you handle API evolution?"
- "What about clients that don't support HTTP/2?"
- "How do you debug gRPC traffic?"

### One-Page Cheat Sheet

```
REST vs gRPC

REST: Representational State Transfer
• HTTP + JSON
• Universal support
• Human readable
• Browser native
• Best for: Public APIs, CRUD, simplicity

gRPC: Google Remote Procedure Call
• HTTP/2 + Protobuf (Protocol Buffers) - binary serialization
• Binary (smaller, faster)
• Strongly typed via .proto definitions
• Native streaming
• Best for: Internal services, streaming, performance

COMMON PATTERN:
Public → REST → API Gateway → gRPC → Internal

PERFORMANCE COMPARISON:
REST JSON: ~10x larger than Protobuf
gRPC: 2-10x faster for serialization
HTTP/2: Multiplexing reduces latency

WHEN TO USE EACH:
Public API → REST
Internal microservices → gRPC
Streaming → gRPC
Browser clients → REST (or gRPC-web)
Mobile → Either (gRPC saves bandwidth)
```

---

## 1.2 API Versioning

### Concept Overview (What & Why)

APIs evolve. Versioning ensures old clients continue working while new features are added.

**Why this matters:**

- Breaking changes without versioning = broken clients
- Mobile apps can't force-update; old versions exist forever
- Enterprise clients need stability guarantees

**Common Approaches:**


| Method              | Example                                      | Pros             | Cons                   |
| ------------------- | -------------------------------------------- | ---------------- | ---------------------- |
| URL path            | `/v1/users`                                  | Clear, cacheable | URL pollution          |
| Query param         | `/users?version=1`                           | Easy to add      | Can be forgotten       |
| Header              | `Accept: application/vnd.api+json;version=1` | Clean URLs       | Hidden, harder to test |
| Content negotiation | `Accept: application/vnd.company.v1+json`    | REST-purist      | Complex                |


### Key Design Principles

**Best Practices:**

1. **Version from day one** - Retrofitting is painful
2. **Support at least N-1** - Give clients time to migrate
3. **Deprecate gracefully** - Announce, warn, then remove
4. **Avoid breaking changes** - Add fields, don't remove
5. **Document breaking changes** - Changelog for each version

**What's a Breaking Change?**

- Removing a field ❌
- Changing a field type ❌
- Changing response structure ❌
- Adding a required field ❌
- Adding optional field ✅
- Adding new endpoint ✅

### Trade-offs & Decision Matrix


| Approach                | Use When                                         |
| ----------------------- | ------------------------------------------------ |
| URL versioning (`/v1/`) | Public APIs, clear separation needed             |
| Header versioning       | Internal APIs, cleaner URLs preferred            |
| No versioning           | Simple internal tools, full control over clients |


### Interview Perspective

**Strong signals:**

- ✅ "We'll use URL versioning for clarity: `/v1/orders`"
- ✅ "Breaking changes go to v2; v1 deprecated with 6-month warning"
- ✅ "Add fields optionally to avoid breaking existing clients"

**Follow-up questions:**

- "How do you deprecate an old version?"
- "What if a bug fix needs to go to all versions?"

---

## 1.3 Pagination Strategies

### Concept Overview (What & Why)

When returning large datasets, pagination prevents:

- Memory exhaustion
- Network timeouts
- Poor user experience

**Three Main Approaches:**


| Method       | Mechanism                                 | Pros                  | Cons                                              |
| ------------ | ----------------------------------------- | --------------------- | ------------------------------------------------- |
| Offset/Limit | `OFFSET 100 LIMIT 20`                     | Simple, random access | Slow for large offsets, inconsistent with inserts |
| Cursor-based | `WHERE id > cursor`                       | Consistent, fast      | No random access                                  |
| Keyset       | `WHERE (date, id) > (last_date, last_id)` | Fast, sorted          | Complex for multi-column sort                     |


### Key Design Principles

**Offset Pagination:**

```
GET /users?offset=100&limit=20
Response: {users: [...], total: 5000}
```

- Problem: `OFFSET 1000000` scans 1M rows
- Problem: If new item inserted, you miss or duplicate items

**Cursor Pagination (Recommended):**

```
GET /users?limit=20&after=cursor_abc123
Response: {users: [...], next_cursor: "cursor_def456", has_more: true}
```

- Cursor = opaque token (often base64 encoded ID)
- Fast: `WHERE id > decoded_cursor LIMIT 20`
- Consistent even with inserts

**Keyset Pagination:**

```sql
SELECT * FROM orders 
WHERE (created_at, id) > ('2024-01-15', 12345)
ORDER BY created_at, id
LIMIT 20
```

### Trade-offs & Decision Matrix


| Requirement                  | Recommended Approach |
| ---------------------------- | -------------------- |
| Simple UI, small dataset     | Offset               |
| Infinite scroll              | Cursor               |
| Consistency during iteration | Cursor or Keyset     |
| Jump to page N               | Offset (or hybrid)   |
| Large dataset (millions)     | Cursor/Keyset only   |


### Interview Perspective

**Strong signals:**

- ✅ "Cursor-based for large datasets; offset is O(n)"
- ✅ "Cursor is opaque; encodes the last ID"
- ✅ "For infinite scroll, cursor is standard"

**Common trap:**

- ❌ "Just use OFFSET" for millions of records

---

## 1.4 Rate Limiting

### Concept Overview (What & Why)

Rate limiting controls how many requests a client can make in a time window.

**Why this matters:**

- Prevents abuse and DDoS
- Ensures fair usage among clients
- Protects backend from overload
- Enforces business limits (API tiers)

### Key Algorithms

**1. Token Bucket (Most Common)**

```
Bucket capacity: 100 tokens
Refill rate: 10 tokens/second
Request costs 1 token
Bucket full? Extra tokens discarded
Bucket empty? Request rejected (429)
```

- Allows bursts up to bucket size
- Smooth average rate equals refill rate
- Used by: AWS, Stripe

**2. Leaky Bucket**

```
Requests enter bucket (queue)
Processed at fixed rate
Bucket full? Request rejected
```

- Smooths out bursts
- Strict rate enforcement
- Used when consistent rate is needed

**3. Fixed Window**

```
Count requests per time window (e.g., per minute)
Reset count at window boundary
Over limit? Reject
```

- Simple to implement
- Problem: Burst at window boundaries (double rate)

**4. Sliding Window Log**

```
Store timestamp of each request
Count requests in last N seconds
Over limit? Reject
```

- Accurate but memory-intensive
- Stores every request timestamp

**5. Sliding Window Counter (Best Balance)**

```
Combine current and previous window
Weight by how far into current window
```

- Smooths window boundary issue
- Memory efficient

### Trade-offs & Decision Matrix


| Algorithm       | Burst Handling        | Memory | Accuracy  | Complexity |
| --------------- | --------------------- | ------ | --------- | ---------- |
| Token Bucket    | Allows bursts         | Low    | Good      | Medium     |
| Leaky Bucket    | Smooths bursts        | Low    | Good      | Medium     |
| Fixed Window    | Poor (2x at boundary) | Low    | Poor      | Low        |
| Sliding Log     | Good                  | High   | Excellent | Medium     |
| Sliding Counter | Good                  | Low    | Good      | Medium     |


### Real-World Implementation

**Where to Rate Limit:**

1. **Client-side** - Prevent accidental floods
2. **Load balancer/API Gateway** - First line of defense
3. **Application layer** - Business logic limits
4. **Database** - Connection limits

**Response When Limited:**

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 30
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640000000
```

**Distributed Rate Limiting:**

- Challenge: Multiple servers need shared state
- Solution: Redis with atomic operations (INCR, EXPIRE)
- Alternative: Approximate with local counters + periodic sync

### Interview Perspective

**What interviewers look for:**

- Understanding of different algorithms
- Trade-off awareness (bursts vs smoothing)
- Distributed implementation awareness

**Strong signals:**

- ✅ "Token bucket for API rate limiting; allows bursts"
- ✅ "Use Redis for distributed rate limiting"
- ✅ "Return 429 with Retry-After header"

**Follow-up questions:**

- "How do you rate limit in a distributed system?"
- "What if someone creates many accounts to bypass limits?"
- "How do you handle different limits for different API tiers?"

### One-Page Cheat Sheet

```
RATE LIMITING

ALGORITHMS:
Token Bucket: Allows bursts, common choice
Leaky Bucket: Smooths to constant rate
Fixed Window: Simple but edge case at boundaries
Sliding Window: Best balance of accuracy/memory

IMPLEMENTATION:
Client ID: API key, IP, user ID
Storage: Redis (distributed), local memory (single server)
Response: 429 + Retry-After header

RATE LIMIT HEADERS:
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000

DISTRIBUTED CHALLENGES:
• Race conditions → Use Redis INCR (atomic)
• Synchronization → Accept approximate counts
• Hot keys → Shard by client ID

BYPASS PROTECTION:
• Rate limit by user ID, not just IP
• Limit account creation rate
• Use CAPTCHAs for suspicious patterns
```

---

## 1.5 Authentication vs Authorization (OAuth2, JWT)

### Concept Overview (What & Why)

**Authentication (AuthN):** "Who are you?" - Verifying identity
**Authorization (AuthZ):** "What can you do?" - Verifying permissions

**Why this matters:**

- Security is non-negotiable
- Wrong implementation = data breaches
- Complex topic that interviewers love to explore

### Key Concepts

**Authentication Methods:**


| Method            | Description                                  | Use Case             |
| ----------------- | -------------------------------------------- | -------------------- |
| Session-based     | Server stores session, client has session ID | Traditional web apps |
| Token-based (JWT) | Stateless token with claims                  | APIs, microservices  |
| OAuth2            | Delegated authorization                      | Third-party access   |
| API Keys          | Simple secret                                | Service-to-service   |


**JWT (JSON Web Token):**

```
Header.Payload.Signature

Header: {"alg": "HS256", "typ": "JWT"}
Payload: {"sub": "1234", "name": "John", "exp": 1640000000}
Signature: HMAC(header + payload, secret)
```

**Pros of JWT:**

- Stateless (no session lookup)
- Scalable (any server can validate)
- Contains claims (roles, permissions)

**Cons of JWT:**

- Can't revoke until expiry
- Token size larger than session ID
- Sensitive data in payload (base64 encoded, not encrypted)

**OAuth2 Flows:**


| Flow                      | Use Case                      |
| ------------------------- | ----------------------------- |
| Authorization Code        | Web apps (most secure)        |
| Authorization Code + PKCE | Mobile/SPA (no client secret) |
| Client Credentials        | Machine-to-machine            |
| Implicit (deprecated)     | Legacy SPAs                   |


### Trade-offs & Decision Matrix


| Approach            | Scalability | Revocation          | Complexity |
| ------------------- | ----------- | ------------------- | ---------- |
| Session + Redis     | High        | Instant             | Medium     |
| JWT (short-lived)   | Highest     | At expiry           | Low        |
| JWT + Refresh Token | High        | Refresh rotation    | Medium     |
| OAuth2              | High        | Token introspection | High       |


### Real-World Implementation

**Microservices Auth Pattern:**

```
1. Client authenticates with Auth Service → Gets JWT
2. Client sends JWT to API Gateway
3. Gateway validates JWT signature
4. Gateway forwards request with user context
5. Downstream services trust gateway-provided context
```

**Token Refresh Strategy:**

- Access token: Short-lived (15 min)
- Refresh token: Longer-lived (7 days), stored securely
- On access token expiry, use refresh token to get new access token
- Refresh token rotation: Issue new refresh token with each refresh

### Interview Perspective

**What interviewers look for:**

- Clear distinction between AuthN and AuthZ
- Understanding of JWT trade-offs
- OAuth2 flow selection

**Common traps:**

- ❌ "JWT can be revoked instantly" (need token blacklist)
- ❌ "Store JWT in localStorage" (XSS vulnerable)
- ❌ Confusing OAuth2 flows

**Strong signals:**

- ✅ "Short-lived JWT + refresh tokens for revocation"
- ✅ "OAuth2 Authorization Code + PKCE for mobile"
- ✅ "API Gateway validates JWT, downstream trusts context"

### One-Page Cheat Sheet

```
AUTHENTICATION vs AUTHORIZATION

AuthN: Who are you? (Login)
AuthZ: What can you do? (Permissions)

JWT STRUCTURE:
Header.Payload.Signature
Payload has claims: sub, exp, roles

JWT TRADE-OFFS:
✅ Stateless, scalable
❌ Can't revoke until expiry
Solution: Short expiry + refresh tokens

OAUTH2 FLOWS:
Authorization Code → Web apps
+ PKCE → Mobile/SPA
Client Credentials → Machine-to-machine

TOKEN STORAGE (Browser):
❌ localStorage (XSS vulnerable)
✅ httpOnly cookie (CSRF protection needed)

MICROSERVICES PATTERN:
Client → Auth Service → JWT
Client → API Gateway (validates JWT) → Services

REFRESH STRATEGY:
Access token: 15 min
Refresh token: 7 days, rotate on use
Store refresh token securely (httpOnly cookie)
```

---

# Section 2: Load Balancing

## 2.1 L4 vs L7 Load Balancing

### Concept Overview (What & Why)

**L4 (Transport Layer):**

- Operates at TCP/UDP level
- Routes based on IP and port
- No knowledge of HTTP content
- Fast, simple, less CPU intensive

**L7 (Application Layer):**

- Operates at HTTP level
- Can inspect headers, URLs, cookies
- Content-based routing
- SSL termination
- More features, more overhead

**Why this matters:**

- L4: Raw throughput, simple routing
- L7: Intelligent routing, required for modern apps

### Key Design Principles


| Feature           | L4                | L7            |
| ----------------- | ----------------- | ------------- |
| Speed             | Faster            | Slower        |
| CPU usage         | Lower             | Higher        |
| SSL termination   | No (pass-through) | Yes           |
| Content routing   | No                | Yes           |
| WebSocket support | Pass-through      | Full support  |
| Health checks     | TCP connect       | HTTP endpoint |
| Sticky sessions   | IP hash           | Cookie-based  |


**When to Use Each:**

**L4:**

- Database connections
- Raw TCP services
- Maximum throughput needed
- Don't need content inspection

**L7:**

- HTTP APIs
- Need URL-based routing
- Need SSL termination
- Need request/response modification
- Need advanced health checks

### Real-World Examples

**Common Architecture:**

```
Internet → L4 LB (fast, handles TCP) → L7 LB (routing) → App Servers
```

**AWS Services:**

- L4: Network Load Balancer (NLB)
- L7: Application Load Balancer (ALB)

**Popular L7 Load Balancers:**

- NGINX, HAProxy, Envoy, AWS ALB, GCP Cloud Load Balancer

### Interview Perspective

**Strong signals:**

- ✅ "L4 for database pooling; L7 for HTTP with path-based routing"
- ✅ "SSL termination at L7 LB reduces backend load"
- ✅ "L7 for canary deployments based on headers"

---

## 2.2 Load Balancing Algorithms

### Concept Overview


| Algorithm                  | Description                               | Best For                         |
| -------------------------- | ----------------------------------------- | -------------------------------- |
| Round Robin                | Rotate through servers                    | Homogeneous servers, equal load  |
| Weighted Round Robin       | Rotate with weights                       | Heterogeneous servers            |
| Least Connections          | Route to server with fewest connections   | Long-lived connections           |
| Weighted Least Connections | Least connections with weights            | Mixed workloads                  |
| IP Hash                    | Hash client IP to server                  | Session affinity without cookies |
| Consistent Hashing         | Minimize redistribution on server changes | Caching, stateful services       |
| Random                     | Random selection                          | Simple, good enough often        |


### Key Design Principles

**Round Robin:**

```
Servers: A, B, C
Requests: 1→A, 2→B, 3→C, 4→A, 5→B, ...
```

- Simple and fair
- Problem: Doesn't account for request complexity

**Least Connections:**

```
Server A: 10 connections
Server B: 5 connections
Server C: 8 connections
New request → B (least)
```

- Better for varying request durations
- Problem: New server gets flooded

**Consistent Hashing:**

```
Hash ring with servers at positions
Request hashed to position, routed to next server
Server removed: Only 1/n requests move
```

- Essential for distributed caches
- Minimizes cache misses on topology changes

### Trade-offs & Decision Matrix


| Scenario                            | Recommended Algorithm   |
| ----------------------------------- | ----------------------- |
| Homogeneous servers, short requests | Round Robin             |
| Heterogeneous servers               | Weighted Round Robin    |
| Long-lived connections (WebSocket)  | Least Connections       |
| Stateful services (cache)           | Consistent Hashing      |
| Need session affinity               | IP Hash or Cookie-based |


### Interview Perspective

**Strong signals:**

- ✅ "Least connections for WebSocket servers"
- ✅ "Consistent hashing for cache layer"
- ✅ "Round robin is fine for stateless, homogeneous services"

---

## 2.3 Sticky Sessions

### Concept Overview (What & Why)

Sticky sessions (session affinity) route all requests from a client to the same server.

**When Needed:**

- Stateful applications with local session storage
- WebSocket connections
- In-memory caches that benefit from repeat access

**Why to Avoid:**

- Creates uneven load distribution
- Complicates scaling and deployments
- Single point of failure for that user's session

### Key Design Principles

**Implementation Methods:**

1. **Cookie-based:** LB sets cookie with server ID
2. **IP-based:** Hash client IP to server (less reliable with NAT)

**Better Alternatives:**

1. **Externalize session:** Store in Redis/Memcached
2. **JWT tokens:** Client holds session data
3. **Stateless design:** No server-side session needed

### Trade-offs & Decision Matrix


| Approach               | Scaling | Failure Handling        | Complexity |
| ---------------------- | ------- | ----------------------- | ---------- |
| Sticky sessions        | Limited | User loses session      | Low        |
| External session store | Easy    | Survives server failure | Medium     |
| JWT/Stateless          | Easy    | No session to lose      | Medium     |


### Interview Perspective

**What interviewers look for:**

- Understanding that sticky sessions are a crutch
- Knowledge of better alternatives
- Knowing when they're actually necessary

**Strong signals:**

- ✅ "Avoid sticky sessions; externalize state to Redis"
- ✅ "If unavoidable, use cookie-based affinity"
- ✅ "WebSocket connections inherently sticky; design for reconnection"

**Common trap:**

- ❌ Defaulting to sticky sessions without considering alternatives

### One-Page Cheat Sheet

```
LOAD BALANCING

L4 vs L7:
L4: TCP level, fast, simple
L7: HTTP level, smart routing, SSL termination

ALGORITHMS:
Round Robin → Equal servers, short requests
Least Connections → Long-lived connections
Consistent Hashing → Cache servers
IP Hash → Session affinity (crude)

STICKY SESSIONS:
Definition: All requests from client → same server
Problem: Uneven load, scaling issues, SPOF

AVOID STICKY SESSIONS:
• Externalize session to Redis
• Use JWT (client-side state)
• Design stateless services

HEALTH CHECKS:
L4: TCP connect succeeds
L7: HTTP 200 from /health endpoint
Interval: 5-30 seconds
Threshold: 2-3 failures before removal
```

---

# Section 3: Databases

## 3.1 SQL Databases

### 3.1.1 Indexing (B-tree, Composite Indexes)

### Concept Overview (What & Why)

**Index:** A data structure that improves query speed at the cost of write overhead and storage.

**B-tree Index (Most Common):**

- Balanced tree structure
- O(log n) lookups
- Supports range queries
- Default for most SQL databases

**Why this matters:**

- Without indexes: Full table scan (O(n))
- With index: O(log n) lookup
- Wrong indexes: Slow writes, wasted storage

### Key Design Principles

**Index Types:**


| Type     | Description           | Use Case                       |
| -------- | --------------------- | ------------------------------ |
| B-tree   | Balanced tree, sorted | General purpose, range queries |
| Hash     | Hash table            | Exact match only               |
| GiST/GIN | Generalized           | Full-text search, arrays, JSON |
| Bitmap   | Bit arrays            | Low cardinality columns        |


**Composite Index:**

```sql
CREATE INDEX idx_user_status_date ON orders(user_id, status, created_at);
```

- Leftmost prefix rule: Index used for queries on (user_id), (user_id, status), or all three
- NOT used for queries on just (status) or (created_at)

**Index Selection Rules:**

1. Index columns in WHERE clauses
2. Index columns in JOIN conditions
3. Index columns in ORDER BY (avoid filesort)
4. Consider covering indexes (include all columns needed)

### Trade-offs & Decision Matrix


| Consideration | More Indexes | Fewer Indexes |
| ------------- | ------------ | ------------- |
| Read speed    | Faster       | Slower        |
| Write speed   | Slower       | Faster        |
| Storage       | More         | Less          |
| Maintenance   | Complex      | Simple        |


**When NOT to Index:**

- Low cardinality columns (e.g., boolean)
- Frequently updated columns
- Small tables (full scan is fast)
- Columns rarely used in queries

### Real-World Examples

```sql
-- Slow: Full table scan
SELECT * FROM orders WHERE user_id = 123;

-- Fast: Index on user_id
CREATE INDEX idx_orders_user ON orders(user_id);

-- Covering index (avoids table lookup)
CREATE INDEX idx_orders_user_status ON orders(user_id, status) INCLUDE (total);
SELECT user_id, status, total FROM orders WHERE user_id = 123;
```

### Interview Perspective

**Strong signals:**

- ✅ "Index on high-cardinality columns used in WHERE"
- ✅ "Composite index with leftmost prefix consideration"
- ✅ "Covering index to avoid table lookup"

**Common traps:**

- ❌ "Index everything" (destroys write performance)
- ❌ Ignoring composite index order

---

### 3.1.2 Transactions and Isolation Levels

### Concept Overview (What & Why)

**ACID Properties:**

- **Atomicity:** All or nothing
- **Consistency:** Valid state to valid state
- **Isolation:** Concurrent transactions don't interfere
- **Durability:** Committed = permanent

**Isolation Levels (from weak to strong):**


| Level            | Dirty Read | Non-repeatable Read | Phantom Read |
| ---------------- | ---------- | ------------------- | ------------ |
| Read Uncommitted | Possible   | Possible            | Possible     |
| Read Committed   | Prevented  | Possible            | Possible     |
| Repeatable Read  | Prevented  | Prevented           | Possible     |
| Serializable     | Prevented  | Prevented           | Prevented    |


**Terminology:**

- **Dirty Read:** See uncommitted changes from another transaction
- **Non-repeatable Read:** Same query returns different results (row modified)
- **Phantom Read:** Same query returns different results (rows added/removed)

### Key Design Principles

**Practical Guidance:**

- Most applications: Read Committed (PostgreSQL default)
- Financial/critical: Serializable (slower but safest)
- Performance critical: Understand trade-offs, accept some anomalies

**Locking Strategies:**


| Strategy    | Description                  | Use Case        |
| ----------- | ---------------------------- | --------------- |
| Optimistic  | Check version at commit time | Low contention  |
| Pessimistic | Lock rows during transaction | High contention |


**Optimistic Locking Example:**

```sql
UPDATE accounts 
SET balance = 100, version = version + 1 
WHERE id = 1 AND version = 5;
-- If affected rows = 0, someone else modified it
```

### Interview Perspective

**Strong signals:**

- ✅ "Read Committed is usually sufficient"
- ✅ "Optimistic locking for low-contention scenarios"
- ✅ "Serializable for financial transactions if needed"

---

### 3.1.3 Sharding Strategies

### Concept Overview (What & Why)

**Sharding:** Splitting data across multiple database instances.

**Why Sharding:**

- Single database can't handle the load
- Data size exceeds single machine capacity
- Geographic distribution needed

**When to Shard:**

- Last resort after: Read replicas, vertical scaling, caching
- Significant complexity cost
- Cross-shard queries are painful

### Key Strategies


| Strategy        | Method                 | Pros                       | Cons                       |
| --------------- | ---------------------- | -------------------------- | -------------------------- |
| Range-based     | user_id 1-1M → shard 1 | Simple, range queries easy | Hot spots possible         |
| Hash-based      | hash(user_id) % N      | Even distribution          | Range queries cross shards |
| Directory-based | Lookup table           | Flexible                   | Lookup service = SPOF      |
| Geographic      | Region → shard         | Data locality              | Cross-region queries hard  |


**Shard Key Selection (Critical):**

- Even distribution
- Query locality (most queries hit one shard)
- Growth pattern consideration
- Avoid hot spots

### Trade-offs & Decision Matrix


| Challenge                | Impact        | Mitigation                         |
| ------------------------ | ------------- | ---------------------------------- |
| Cross-shard queries      | Slow, complex | Denormalize, choose key wisely     |
| Cross-shard transactions | Very hard     | Avoid, use saga pattern            |
| Rebalancing              | Painful       | Consistent hashing, virtual shards |
| Operational complexity   | High          | Automation, monitoring             |


### Real-World Examples

- **Instagram:** Shard by user_id (photos stored with owner)
- **Discord:** Shard by guild_id (messages in guild together)
- **Slack:** Shard by workspace

### Interview Perspective

**Strong signals:**

- ✅ "Shard key selection is critical; query pattern analysis first"
- ✅ "Avoid sharding until necessary; significant complexity"
- ✅ "Consistent hashing for rebalancing"

**Common trap:**

- ❌ "Shard from day one" (premature optimization)

---

## 3.2 NoSQL Databases

### 3.2.1 Key-Value vs Document vs Column Stores

### Concept Overview


| Type          | Model             | Examples                   | Best For                          |
| ------------- | ----------------- | -------------------------- | --------------------------------- |
| Key-Value     | Simple key→value  | Redis, DynamoDB, Memcached | Caching, sessions, simple lookups |
| Document      | Key→JSON document | MongoDB, Couchbase         | Flexible schema, nested data      |
| Column Family | Row key→columns   | Cassandra, HBase           | Time-series, write-heavy          |
| Graph         | Nodes + edges     | Neo4j, Neptune             | Relationships, social networks    |


### Key Design Principles

**Key-Value:**

```
key: "user:123"
value: {name: "John", email: "john@example.com"}
```

- Simplest model
- No query by value (need to know key)
- Ultra-fast lookups

**Document:**

```json
{
  "_id": "order_123",
  "user": {"id": 1, "name": "John"},
  "items": [{"sku": "A", "qty": 2}],
  "total": 99.99
}
```

- Flexible schema
- Nested documents reduce joins
- Can query by any field (with indexes)

**Column Family:**

```
Row key: "user:123"
Columns: {
  "profile:name": "John",
  "profile:email": "john@x.com",
  "activity:last_login": "2024-01-15"
}
```

- Wide rows (millions of columns possible)
- Columns stored together (efficient for sparse data)
- Great for time-series (timestamp as column)

### Trade-offs & Decision Matrix


| Need                  | Choose                              |
| --------------------- | ----------------------------------- |
| Simple cache          | Key-Value (Redis)                   |
| Flexible documents    | Document (MongoDB)                  |
| High write throughput | Column (Cassandra)                  |
| Time-series data      | Column or specialized (TimescaleDB) |
| Complex relationships | Graph (Neo4j)                       |


---

### 3.2.2 DynamoDB / Cassandra Partitioning

### Concept Overview

Both use partition key to distribute data:

```
Partition Key → Hash → Partition/Node
```

**DynamoDB:**

- Partition key (hash key)
- Optional sort key
- Partition key determines storage location
- Sort key enables range queries within partition

**Cassandra:**

- Partition key
- Clustering columns (for sorting within partition)
- Consistent hashing with virtual nodes

### Key Design Principles

**Good Partition Key:**

- High cardinality (many unique values)
- Even access distribution
- Query pattern aligned (most queries specify partition key)

**Bad Partition Key Examples:**

- Date (all today's writes hit one partition)
- Country (US partition overwhelmed)
- Status (few values, uneven distribution)

### 3.2.3 Hot Partitions Problem

**What It Is:**
One partition receives disproportionate traffic, becoming a bottleneck.

**Causes:**

- Partition key with low cardinality
- Celebrity problem (one user has millions of followers)
- Temporal patterns (all writes to today's partition)

**Solutions:**


| Solution             | Description                                     |
| -------------------- | ----------------------------------------------- |
| Add suffix/salt      | `partition_key = user_id + random(0-9)`         |
| Scatter-gather       | Write to multiple partitions, aggregate on read |
| Time-based bucketing | `partition_key = user_id + hour`                |
| Write sharding       | Shard hot keys, aggregate in background         |


**Example - Celebrity Problem:**

```
Instead of: partition_key = celebrity_id
Use: partition_key = celebrity_id + random(0-99)
Read: Query all 100 partitions, aggregate
```

### Interview Perspective

**Strong signals:**

- ✅ "Partition key selection based on access patterns"
- ✅ "Hot partition mitigation with salting"
- ✅ "High cardinality partition keys"

**Common trap:**

- ❌ Using timestamp as partition key (all writes to one partition)

### One-Page Cheat Sheet

```
DATABASE SELECTION

SQL (PostgreSQL, MySQL):
• ACID transactions
• Complex queries, joins
• Structured data
• When: Most applications, consistency critical

NoSQL - Key-Value (Redis, DynamoDB):
• Simple lookups by key
• Caching, sessions
• When: Need speed, simple access patterns

NoSQL - Document (MongoDB):
• Flexible schema
• Nested documents
• When: Evolving schema, document-oriented data

NoSQL - Column (Cassandra):
• Write-heavy, time-series
• Linear scalability
• When: High write throughput, append-heavy

INDEXING:
• B-tree: Range queries, general purpose
• Composite: Leftmost prefix rule
• Covering: Include all needed columns

PARTITIONING/SHARDING:
• Key selection is critical
• High cardinality
• Aligned with query patterns
• Watch for hot partitions

HOT PARTITION SOLUTIONS:
• Salt/suffix the key
• Scatter-gather pattern
• Time-based bucketing
```

---

# Section 4: Caching

### What is caching?

Caching is a performance optimization technique where frequently accessed data is stored in a fast, temporary storage layer so future requests can be served quicker without recomputing or fetching from slow systems.

**🔹 Simple Definition**

Caching = storing frequently used data in memory (or fast storage) to reduce latency, load, and cost.

---

## 4.1 Cache Levels: Client vs CDN vs Server

### Concept Overview (What & Why)


| Level                     | Location         | What's Cached                | TTL                |
| ------------------------- | ---------------- | ---------------------------- | ------------------ |
| Browser/Client            | User's device    | Static assets, API responses | Minutes to days    |
| CDN                       | Edge locations   | Static assets, some dynamic  | Hours to days      |
| API Gateway/Reverse Proxy | Application edge | API responses                | Seconds to minutes |
| Application               | Service memory   | Computed data, DB results    | Seconds to minutes |
| Distributed Cache         | Redis/Memcached  | Shared data                  | Seconds to hours   |
| Database                  | Query cache      | Query results                | Automatic          |


### Key Design Principles

**Multi-tier Caching:**

```
Request → Browser Cache → CDN → API Gateway Cache → App Cache → Redis → Database
```

Each tier reduces load on the next.

**What to Cache Where:**


| Data Type                       | Cache Location      |
| ------------------------------- | ------------------- |
| Static assets (JS, CSS, images) | CDN, browser        |
| User session                    | Redis               |
| Database query results          | Application + Redis |
| Computed aggregations           | Application + Redis |
| Full page responses             | CDN (for anonymous) |


### Interview Perspective

**Strong signals:**

- ✅ "CDN for static assets, Redis for session/dynamic data"
- ✅ "Cache-Control headers for browser caching"
- ✅ "Multi-tier caching reduces database load"

---

## 4.2 Cache Eviction: LRU, LFU, TTL

### Concept Overview


| Policy                      | Evicts          | Best For                              |
| --------------------------- | --------------- | ------------------------------------- |
| LRU (Least Recently Used)   | Oldest accessed | General purpose, recency matters      |
| LFU (Least Frequently Used) | Least accessed  | Frequency matters, long-term patterns |
| TTL (Time To Live)          | Expired items   | Data with known freshness             |
| FIFO (First In First Out)   | Oldest inserted | Simple, predictable                   |
| Random                      | Random item     | Simple, surprisingly effective        |


### Key Design Principles

**LRU (Most Common):**

- Good default choice
- Implemented with hash map + doubly linked list
- O(1) for get and put
- Problem: Scan resistance (full scan pollutes cache)

**LFU:**

- Better for long-term popularity patterns
- Problem: New items hard to compete
- Solution: LFU with aging

**TTL:**

- Essential for data freshness
- Set based on consistency requirements
- Short TTL = more DB hits, fresher data
- Long TTL = fewer DB hits, staler data

### Trade-offs & Decision Matrix


| Scenario                 | Recommended Policy             |
| ------------------------ | ------------------------------ |
| General web caching      | LRU + TTL                      |
| Trending/popular content | LFU                            |
| Session data             | TTL only                       |
| CDN                      | TTL with cache-control headers |
| Limited memory           | LRU (simple, effective)        |


---

## 4.3 Cache Patterns: Cache-Aside vs Write-Through

Caching patterns define how your application interacts with cache and database—especially for reads and writes. These are critical in system design interviews because they directly impact consistency, latency, and scalability.

### Concept Overview

**Cache-Aside (Lazy Loading):**

```
Read:
1. App Checks cache
2. Cache miss → Query database
3. Store in cache
4. Return data/response

Write:
1. Write to database
2. Invalidate cache (or let it expire)

✅ Pros:
Simple, flexible
Only cache what is needed

❌ Cons:
Cache miss latency (first request slow)
Stale data possible

📌 Use Case:
Product catalog, user profiles
```

**Write-Through:**

```
🔧 Flow
Write → Cache → DB (synchronously)

Write:
1. Write to cache
2. Cache synchronously writes to database
3. Return success

Read:
1. Always from cache (cache is source of truth)

✅ Pros:
Strong consistency
Cache always up-to-date

❌ Cons:
Higher write latency
Unnecessary cache writes

📌 Use Case:
Banking systems, critical data
```

**Write-Behind (Write-Back):**

```
🔧 Flow
Write → Cache
DB update happens later (async)
 
Write:
1. Write to cache
2. Return success immediately
3. Async: Cache writes to database in batches

✅ Pros:
Very fast writes
High throughput

❌ Cons:
Risk of data loss (if cache crash)
Eventual consistency

📌 Use Case:
Analytics, logging, high-write systems
```

**Read-through (loader-owned miss path):**

```
🔧 Flow
App → Cache
Cache automatically fetches from DB on miss

Read:
1. App asks cache
2. On miss, cache (or cache library) calls a loader that reads DB and populates cache
3. Return value

✅ Pros:
Cleaner app logic
Cache handles data loading

❌ Cons:
More complex cache system
Less control

📌 Use Case:
Standardized platforms, middleware-heavy systems

```

Refresh-Ahead

```
🔧 Flow
Cache refreshes data before expiry (TTL)

✅ Pros
Low latency always
No cache miss

❌ Cons:
Extra load on DB
Might refresh unused data

📌 Use Case:
Hot data (trending products, dashboards)
```

The app does not manually “if miss then DB then set”; the cache layer triggers the load. Compared to cache-aside, logic is centralized—good for libraries— but the loader contract and timeouts must be designed carefully.

### Theory tie-in

**Who owns the truth:** In cache-aside, the *application* orchestrates DB + cache and usually invalidates on writes. In write-through/write-behind, the *cache path* is coupled to persistence—simpler reads, more machinery on writes.

**Invalidation vs update:** On a write you can *delete* cache keys (next read refills) or *update* them in place (fewer misses, riskier if the update is partial or wrong). Delete + refill is the common default because it is easier to reason about.

### Key Design Principles


| Pattern       | Consistency | Read Performance  | Write Performance | Complexity |
| ------------- | ----------- | ----------------- | ----------------- | ---------- |
| Cache-Aside   | Eventual    | Fast (after warm) | Fast              | Low        |
| Write-Through | Strong      | Fast              | Slower (sync)     | Medium     |
| Write-Behind  | Eventual    | Fast              | Fastest           | High       |
| Read-Through  | Eventual    | Fast              | N/A               | Medium     |


### Trade-offs & Decision Matrix


| Requirement                    | Recommended Pattern |
| ------------------------------ | ------------------- |
| Read-heavy, can tolerate stale | Cache-Aside         |
| Strong consistency             | Write-Through       |
| Write-heavy, can tolerate loss | Write-Behind        |
| Simple implementation          | Cache-Aside         |


---

**🔹 Real System Mapping**

```
Instagram Feed → Cache-Aside + TTL
Bank Transactions → Write-Through
Analytics Systems → Write-Back
Netflix/CDN → Refresh-Ahead
```

---

## 4.4 Cache Stampede Problem

### Concept Overview (What & Why)

Cache Stampede (Thundering Herd):

- When a popular cache entry expires, many concurrent requests hit the database simultaneously.
- Cache stampede occurs when many requests hit the system at the same time after a cache entry expires, causing all of them to fetch data from the database simultaneously—overloading it.

---

```
scenario:
A popular key (e.g., “trending products”) is cached
TTL expires
Suddenly 1000 requests arrive
All see cache miss ❌
All hit DB at once 💥

👉 Result:

DB overload
Increased latency
Possible system crash
```

**Why It’s Dangerous**

- 📉 Database gets overwhelmed
- ⏳ Latency spikes for users
- ❌ Cascading failures possible
- 📈 Traffic amplification (1 request → 1000 DB calls)

### Solutions


| Solution                        | Description                           | Complexity |
| ------------------------------- | ------------------------------------- | ---------- |
| Mutex / Locking (Single Flight) | Only one request fetches, others wait | Medium     |
| Probabilistic early expiry      | Refresh before expiry randomly        | Low        |
| Background refresh              | Async refresh before expiry           | Medium     |
| Request coalescing              | Combine concurrent requests           | High       |


**Locking Example:**

```
if (cache miss) {
acquire lock
if still miss → fetch from DB
update cache
release lock
}

✅ Pros
Prevents DB overload

❌ Cons
Adds latency for waiting requests
```

**Probabilistic Early Expiry/Cache Expiry Randomization (Jitter) :**

```
Idea: Don’t let all keys expire at same time

TTL = baseTTL + random(0, 5 min)

✅ Pros
Avoids sudden traffic spikes
```

**Refresh-Ahead (Proactive Refresh)/Background refresh**

```
Idea: Refresh cache before it expires
   - Background job updates popular keys

✅ Pros
   - No cache miss for hot data
```

**Stale-While-Revalidate**

```
Idea: Serve old (stale) data while refreshing in background

if expired:
   return stale data
   trigger async refresh

✅ Pros
- Low latency
- Great UX

❌ Cons
- Slightly outdated data
```

**Request Coalescing**

```
Idea: Combine multiple identical requests into one

Only one DB call per key
```

**⚡ Rate Limiting / Circuit Breaker**

```
Protect DB during spikes
Fail fast or degrade gracefully
```

---

**Real-World Example**

```
Flash sale on Flipkart/Amazon
Cache for product pricing expires
Millions of users hit at same time → DB crash if not handled
```

**🔹 Best Practice Combination ⭐**
In real systems, you don’t use just one:

Cache-Aside + Mutex Lock
TTL + Jitter
Stale-While-Revalidate for hot data
Background refresh for critical keys


## 🔥 Cache Invalidation (Most Difficult Problem in Caching)

### 🔹 Interview Definition

**Cache invalidation is the process of ensuring cached data stays consistent with the source of truth (database) by updating or removing stale entries when underlying data changes.**

### 🔹 Why It’s Hard

- Cache is **fast but not always correct**
- DB is **correct but slow**
- Challenge = **balance consistency vs performance**

👉 Classic quote:

> “There are only two hard things in computer science: cache invalidation and naming things.”

---

### 🔹 When Invalidation Happens

- Data update (price change, profile update)
- Delete operation
- TTL expiry
- Bulk updates (e.g., flash sale)

---

### 🔹 Core Invalidation Strategies

### 1. ⏱️ TTL (Time-To-Live)

### 🔧 How it works

- Cache expires automatically after time

```text
product:123 → expires in 10 min

```

### ✅ Pros

- Simple
- No extra logic

### ❌ Cons

- Can serve stale data until expiry

### 📌 Use Case

- Non-critical data (feeds, analytics)

---

## 2. ❌ Write Invalidate (Most Common)

### 🔧 Flow

1. Update DB
2. Delete cache

```text
UPDATE product SET price=500 WHERE id=123
DELETE cache:product:123

```

### ✅ Pros

- Ensures fresh reads

### ❌ Cons

- Next read → cache miss (DB hit)

### 📌 Use Case

- E-commerce product updates

---

## 3. ✏️ Write Through (Update Cache + DB)

### 🔧 Flow

- Write to DB and cache together

### ✅ Pros

- Strong consistency

### ❌ Cons

- Slower writes
- Cache pollution

---

## 4. 🔄 Event-Based Invalidation (Advanced)

### 🔧 Flow

- DB update → event → invalidate cache

Example using Apache Kafka:

```text
Service A updates DB
→ publishes event
→ Service B invalidates cache

```

### ✅ Pros

- Works in microservices
- Decoupled

### ❌ Cons

- Complexity
- Event delay → temporary inconsistency

---

## 5. 🧊 Versioning (Advanced)

### 🔧 Idea

- Add version to cache key

```text
product:123:v2

```

- New update → increment version

### ✅ Pros

- No stale reads

### ❌ Cons

- Old keys remain (memory overhead)

---

## 6. 🔁 Stale-While-Revalidate

### 🔧 Flow

- Serve stale data
- Refresh in background

### ✅ Pros

- Low latency
- Good UX

### ❌ Cons

- Slightly outdated data

---

# 🔹 Choosing Strategy (Interview Insight)


| Requirement         | Strategy               |
| ------------------- | ---------------------- |
| Simple system       | TTL                    |
| Strong consistency  | Write Through          |
| Balanced approach ⭐ | Write Invalidate       |
| Microservices       | Event-Based            |
| High read traffic   | Stale-While-Revalidate |


---

# 🔹 Real-World Mapping

- **Amazon / Flipkart** → Write Invalidate + TTL
- **Banking Systems** → Write Through
- **Netflix / YouTube** → TTL + Refresh + CDN
- **Microservices systems** → Event-driven (Kafka)

---

# 🔹 Common Pitfalls (Interview Gold)

### ❌ 1. Race Condition

- Cache updated before DB commit

### ❌ 2. Lost Invalidation

- Cache not deleted → stale forever

### ❌ 3. Thundering Herd

- After invalidation → DB spike

### ❌ 4. Partial Failures

- DB success, cache delete failed

---

# 🔹 Best Practice (Production Setup)

👉 Combine strategies:

- Write Invalidate + TTL
- Add retry on cache delete
- Use event-based invalidation in microservices
- Add monitoring (cache hit ratio)

---

# 🔹 3-Min Interview Answer

> Cache invalidation ensures that cached data remains consistent with the database. Common strategies include TTL-based expiration, write-invalidate where cache is deleted after DB updates, and write-through where both cache and DB are updated together. In distributed systems, event-based invalidation using systems like Kafka is used. The main challenge is balancing consistency and performance while handling race conditions and stale data.

---

### Interview Perspective

**What interviewers look for:**

- Awareness of the problem
- Multiple solution approaches
- Trade-off understanding

**Strong signals:**

- ✅ "Locking to prevent stampede"
- ✅ "Background refresh before expiry"
- ✅ "This is critical for popular cache entries"

### One-Page Cheat Sheet

```
CACHING
storing frequently used data in memory (or fast storage) to reduce latency, load, and cost.

CACHE LEVELS:
Browser → CDN → Gateway → App → Redis → Database

EVICTION POLICIES:
LRU: General purpose (default choice)
LFU: Popularity-based
TTL: Time-based freshness

CACHE PATTERNS:
Cache-Aside: App checks cache, fills on miss
Read-Through: Cache triggers loader on miss (app does not hand-fill)
Write-Through: Cache writes to DB synchronously
Write-Behind: Cache writes to DB async

CACHE-ASIDE (Most Common):
Read: Cache miss → DB → Store in cache
Write: DB → Invalidate cache

CACHE STAMPEDE:
Problem: Hot key expires, all requests hit DB
Solutions:
• Locking (one fetches, others wait)
• Early probabilistic refresh
• Background refresh before expiry

CACHE INVALIDATION:
"Two hard problems: cache invalidation and naming things"
Strategies:
• TTL (simple, eventual consistency)
• Event-driven invalidation (complex, immediate)
• Version keys (user:123:v2)

NUMBERS:
Redis: 100k+ ops/sec
Memcached: 200k+ ops/sec
Network round trip: 0.5ms same DC
```

---

# Phase 1 Summary: Building Blocks Integration

These core building blocks form the foundation of any system design. In an interview, you should:

1. **API Design:** Choose REST for public, gRPC for internal. Version from day one.
2. **Load Balancing:** L7 for HTTP services, proper algorithm selection, avoid sticky sessions.
3. **Database:** SQL for transactions, NoSQL for scale/flexibility. Indexing and sharding are advanced topics.
4. **Caching:** Multi-tier strategy, appropriate eviction, handle stampede.

**Integration Example (E-commerce):**

```
Mobile App ← REST API
         ↓
    API Gateway (L7 LB, Rate Limiting, Auth)
         ↓
    Product Service (gRPC) ←→ Redis Cache
         ↓                       ↓
    PostgreSQL (Primary) ←→ Read Replicas
         ↓
    CDN (Product Images)
```

**Interviewer Expectations by Level:**


| Level          | Expectation                                             |
| -------------- | ------------------------------------------------------- |
| Senior (L5)    | Know all building blocks, make reasonable choices       |
| Staff (L6)     | Deep understanding of trade-offs, can justify decisions |
| Principal (L7) | Can challenge conventional wisdom, knows edge cases     |



---

## Common Interview Questions & Model Answers

---

### Q1: When would you choose REST over gRPC?

**Expected Answer:**

- REST for public-facing APIs — universal browser/client support, human-readable JSON, mature caching (ETags, 304)
- gRPC for internal microservices — binary Protobuf (7-10x smaller payloads), HTTP/2 multiplexing, native bidirectional streaming
- Hybrid is the standard pattern: Public → REST → API Gateway → gRPC → Internal services
- REST when developer onboarding, third-party integrations, or caching is the priority
- gRPC when latency, bandwidth, or streaming is critical (e.g., real-time stock quotes)

**Follow-up Q:** "How would you handle API versioning?"

- URL versioning (`/v1/users`) for public APIs — explicit, cacheable, easy deprecation (used by Stripe, GitHub)
- Header versioning for internal APIs — cleaner URLs
- Support N-1 versions; deprecate with 6-12 month sunset period
- Add optional fields (safe), never remove fields without version bump

**Follow-up Q:** "How do you debug gRPC traffic?"

- gRPC reflection + grpcurl or Postman gRPC support
- Log requests/responses at interceptor level, convert Protobuf to JSON for readability
- Load balancer must support HTTP/2; confirm with health-check endpoints

---

### Q2: What load balancing algorithm would you use and why?

**Expected Answer:**

- Round Robin for homogeneous stateless services (simple, fair)
- Least Connections for long-lived connections (WebSocket, database pools)
- Consistent Hashing for cache layers (minimizes cache misses on topology changes)
- Weighted Round Robin for heterogeneous servers (bigger servers get more traffic)
- L4 for raw TCP/UDP throughput; L7 for HTTP routing, SSL termination, path-based routing

**Follow-up Q:** "How do you avoid the load balancer becoming a single point of failure?"

- Active-passive pair with VRRP (1-3s failover)
- Active-active with DNS round-robin across multiple LBs
- Managed cloud LBs (AWS ALB/NLB) — multi-AZ by default, no SPOF
- L4 + L7 layered approach for redundancy

**Follow-up Q:** "When would you use sticky sessions?"

- Prefer to avoid — creates uneven load, complicates scaling, SPOF per user session
- If unavoidable (WebSocket, in-memory cache): use cookie-based affinity
- Better alternative: externalize state to Redis, use JWT for client-side sessions

---

### Q3: How does database indexing work? When would you NOT add an index?

**Expected Answer:**

- Index = B+ tree data structure enabling O(log n) lookups vs O(n) full table scan
- Composite index follows leftmost prefix rule: `(A, B, C)` works for queries on `(A)`, `(A, B)`, or `(A, B, C)`
- Covering index includes all queried columns — avoids table lookup entirely (index-only scan)
- **Do NOT index:** small tables (<1K rows), low cardinality columns (boolean/status), frequently updated columns, write-heavy log tables

**Follow-up Q:** "Query is using an index but still slow. What could be wrong?"

- Leading wildcard (`LIKE '%smith'`) — index can't help
- Function on indexed column (`WHERE LOWER(email) = ...`) — use functional index
- Implicit type conversion (string vs int) — mismatched types bypass index
- Large OFFSET pagination — switch to cursor-based
- Stale optimizer statistics — run `ANALYZE` to refresh

---

### Q4: When would you choose SQL vs NoSQL? Give a concrete example of using both.

**Expected Answer:**

- **SQL:** ACID required, complex joins, well-defined schema, ad-hoc queries (orders, payments, inventory)
- **NoSQL Key-Value (Redis):** caching, sessions, simple lookups — sub-ms latency
- **NoSQL Document (MongoDB):** flexible schema, nested data, evolving models
- **NoSQL Column (Cassandra):** write-heavy, time-series, append-only workloads — 1M+ writes/sec
- Polyglot persistence is the norm: PostgreSQL for transactions + Redis for cache + Cassandra for logs + Elasticsearch for search

**Follow-up Q:** "How would you handle a sharding key that creates hot partitions?"

- Salt/suffix the key (`celebrity_id + random(0-99)`) to spread writes across partitions
- Scatter-gather on reads (query all sub-partitions, merge)
- Use consistent hashing for easier rebalancing
- Dedicated shards for known hot entities (celebrity accounts)

---

### Q5: Explain cache-aside vs write-through. When do you use each?

**Expected Answer:**

- **Cache-aside:** App checks cache → miss → query DB → store in cache; write goes to DB, cache invalidated
  - Simple, only caches accessed data, risk of stale reads
  - Best for: read-heavy general-purpose (product catalogs, user profiles)
- **Write-through:** Write → cache → DB synchronously; reads always from cache
  - Strong consistency, higher write latency
  - Best for: critical data (banking, configuration)
- **Write-behind:** Write → cache → async DB; fastest writes, risk of data loss on cache crash
  - Best for: analytics, logging, high-write systems

**Follow-up Q:** "How do you prevent cache stampede?"

- Mutex/locking: one request fetches, others wait
- Probabilistic early expiry: refresh before TTL with randomized jitter
- Stale-while-revalidate: serve stale data, refresh in background
- Background refresh for known hot keys (trending items)

**Follow-up Q:** "What is your cache invalidation strategy?"

- Write-invalidate + TTL (most common): delete cache on DB write, TTL as safety net
- Event-based invalidation in microservices (DB update → Kafka event → cache delete)
- Add retry on cache delete failure; monitor cache hit ratio as key metric

---

### Q6: How would you design a rate limiter?

**Expected Answer:**

- Token Bucket (most common): allows bursts up to bucket size, steady average at refill rate
- Sliding Window Counter for balanced accuracy and memory efficiency
- Fixed Window is simple but has boundary burst problem (2x rate at window edges)
- Distributed: Redis with atomic `INCR` + `EXPIRE`; accept approximate counts
- Response: HTTP 429 + `Retry-After` header + rate limit headers (`X-RateLimit-Remaining`)
- Rate limit by user ID (not just IP) to prevent multi-account bypass

**Follow-up Q:** "What if the rate limiter itself becomes a bottleneck?"

- Shard rate limiting state by client ID across multiple Redis instances
- Use local in-memory counters with periodic sync (accept approximation)
- Cache rate limit decisions locally for short TTL (e.g., 1 second)

---

### Q7: Explain JWT-based authentication in a microservices architecture.

**Expected Answer:**

- Client authenticates with Auth Service → receives JWT (short-lived, ~15 min) + refresh token (~7 days)
- JWT is stateless: contains claims (sub, roles, exp), verified by signature — any service can validate
- API Gateway validates JWT signature, forwards user context to downstream services
- Trade-off: cannot revoke JWT until expiry — mitigate with short TTL + refresh token rotation
- Store tokens in httpOnly cookies (not localStorage — XSS vulnerable)

**Follow-up Q:** "How do you revoke a JWT before it expires?"

- Maintain a token blacklist in Redis (check on each request) — adds statefulness
- Keep access tokens very short-lived (5-15 min) — reduces revocation window
- Rotate refresh tokens on each use — detect theft if old token reused

---

**Navigation:** [← Previous: Foundations](00-foundations.md) | [Next: Distributed Systems →](02-distributed-systems.md)
