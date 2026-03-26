# Generic Large-Scale Distributed System Architecture
> Professional system design diagram for interview preparation

---

## 📊 Complete Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌──────────────┐              ┌──────────────┐           ┌──────────────┐        │
│  │ Web Browser  │              │ Mobile App   │           │  Desktop App │        │
│  │              │              │              │           │              │        │
│  │ • Cookies    │              │ • Tokens     │           │ • Tokens     │        │
│  │ • LocalStore │              │ • App Cache  │           │ • App Cache  │        │
│  │ • SessionStore│             │ • Keychain   │           │ • Keychain   │        │
│  └──────┬───────┘              └──────┬───────┘           └──────┬───────┘        │
│         │                             │                          │                │
│         └─────────────────────────────┴──────────────────────────┘                │
│                                       │                                            │
└───────────────────────────────────────┼────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         EDGE & SECURITY LAYER                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌────────────────────────────────────────────────────────────────────┐            │
│  │  DNS Resolution (Route 53 / Cloudflare DNS)                        │            │
│  │  • Geo-routing  • Health checks  • Failover                        │            │
│  └────────────────┬───────────────────────────────────────────────────┘            │
│                   │                                                                 │
│                   ▼                                                                 │
│  ┌────────────────────────────────────────────────────────────────────┐            │
│  │  CDN (CloudFront / Cloudflare / Akamai)                            │            │
│  │  • Static assets (JS, CSS, images, videos)                         │            │
│  │  • Edge caching  • SSL/TLS termination                             │            │
│  │  Cache Hit → Return ↺   Cache Miss → Continue ↓                    │            │
│  └────────────────┬───────────────────────────────────────────────────┘            │
│                   │                                                                 │
│                   ▼                                                                 │
│  ┌────────────────────────────────────────────────────────────────────┐            │
│  │  WAF (Web Application Firewall)                                    │            │
│  │  • SQL injection protection  • XSS protection                      │            │
│  │  • DDoS mitigation  • Bot detection                                │            │
│  └────────────────┬───────────────────────────────────────────────────┘            │
│                   │                                                                 │
│                   ▼                                                                 │
│  ┌────────────────────────────────────────────────────────────────────┐            │
│  │  Rate Limiter (Token Bucket / Sliding Window)                      │            │
│  │  • Per user: 100 req/min  • Per IP: 1000 req/min                   │            │
│  │  Exceeded → 429 Too Many Requests                                  │            │
│  └────────────────┬───────────────────────────────────────────────────┘            │
│                   │                                                                 │
└───────────────────┼─────────────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                    API & TRAFFIC MANAGEMENT                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌────────────────────────────────────────────────────────────────────┐            │
│  │  Load Balancer (ALB / NLB / HAProxy)                               │            │
│  │  • Layer 7 routing  • SSL termination  • Health checks             │            │
│  │  • Sticky sessions  • Cross-zone load balancing                    │            │
│  │  Algorithm: Least connections / Round robin / Weighted             │            │
│  └────────────────┬───────────────────────────────────────────────────┘            │
│                   │                                                                 │
│                   ▼                                                                 │
│  ┌────────────────────────────────────────────────────────────────────┐            │
│  │  API Gateway (Kong / Apigee / AWS API Gateway)                     │            │
│  │                                                                     │            │
│  │  • Authentication/Authorization (JWT, OAuth 2.0, API Keys)         │            │
│  │  • Request routing & transformation                                │            │
│  │  • Rate limiting & throttling                                      │            │
│  │  • API versioning (v1, v2)                                         │            │
│  │  • Request/Response validation                                     │            │
│  │  • Logging & monitoring                                            │            │
│  └────────────────┬───────────────────────────────────────────────────┘            │
│                   │                                                                 │
└───────────────────┼─────────────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         APPLICATION LAYER                                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   User       │  │   Product    │  │   Order      │  │  Payment     │          │
│  │  Service     │  │   Service    │  │  Service     │  │  Service     │  ...     │
│  │              │  │              │  │              │  │              │          │
│  │ • Stateless  │  │ • Stateless  │  │ • Stateless  │  │ • Stateless  │          │
│  │ • Auto-scale │  │ • Auto-scale │  │ • Auto-scale │  │ • Auto-scale │          │
│  │ • Health ✓   │  │ • Health ✓   │  │ • Health ✓   │  │ • Health ✓   │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                 │                    │
│         └─────────────────┴─────────────────┴─────────────────┘                    │
│                           │                                                         │
│         ┌─────────────────┴────────────────────┐                                   │
│         │                                      │                                    │
│         ▼                                      ▼                                    │
│  ┌────────────────────────────┐    ┌────────────────────────────┐                 │
│  │  Service Discovery         │    │  Circuit Breaker           │                 │
│  │  (Consul / Eureka / K8s)   │    │  (Hystrix / Resilience4j)  │                 │
│  │  • Dynamic registration    │    │  • Fail fast               │                 │
│  │  • Service health          │    │  • Fallback mechanisms     │                 │
│  └────────────────────────────┘    └────────────────────────────┘                 │
│                                                                                     │
│  Inter-Service Communication:                                                      │
│  ┌──────────────────────────────────────────────────────────────┐                 │
│  │  Synchronous: REST (HTTP/JSON), gRPC (HTTP/2 + Protobuf)     │                 │
│  │  Asynchronous: Message Queue, Event Streaming                │                 │
│  │  Retry Strategy: Exponential backoff with jitter             │                 │
│  │  Timeout: Connection (5s), Read (30s)                        │                 │
│  └──────────────────────────────────────────────────────────────┘                 │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
                    │                                      │
                    │                                      │
        ┌───────────┴─────────────┐           ┌───────────┴─────────────┐
        │                         │           │                         │
        ▼                         ▼           ▼                         ▼
┌─────────────────┐      ┌─────────────────────────────────┐    ┌──────────────────┐
│  CACHING LAYER  │      │     DATA LAYER                  │    │  ASYNC LAYER     │
├─────────────────┤      ├─────────────────────────────────┤    ├──────────────────┤
│                 │      │                                 │    │                  │
│ ┌─────────────┐ │      │  ┌────────────────────────┐    │    │ ┌──────────────┐ │
│ │   Redis     │ │      │  │ Relational DB          │    │    │ │ Message      │ │
│ │  Cluster    │ │      │  │ (PostgreSQL / MySQL)   │    │    │ │ Queue        │ │
│ │             │ │      │  │                        │    │    │ │ (Kafka/      │ │
│ │ • Master-   │ │      │  │  ┌──────────┐          │    │    │ │  RabbitMQ/   │ │
│ │   Replica   │ │      │  │  │ Primary  │          │    │    │ │  SQS)        │ │
│ │ • Sentinal  │ │◄─────┼──┼──│   DB     │          │    │    │ │              │ │
│ │ • TTL-based │ │  ▲   │  │  └────┬─────┘          │    │    │ │ Topics:      │ │
│ │   eviction  │ │  │   │  │       │ Replication    │    │    │ │ • order.     │ │
│ └─────────────┘ │  │   │  │       ▼                │    │    │ │   created    │ │
│                 │  │   │  │  ┌──────────┐          │    │    │ │ • payment.   │ │
│ ┌─────────────┐ │  │   │  │  │  Read    │          │    │    │ │   processed  │ │
│ │ Memcached   │ │  │   │  │  │ Replica  │          │    │    │ │ • email.     │ │
│ │             │ │  │   │  │  │    1     │          │    │    │ │   send       │ │
│ │ • Session   │ │  │   │  │  └──────────┘          │    │    │ └──────┬───────┘ │
│ │   data      │ │  │   │  │                        │    │    │        │         │
│ │ • Hot data  │ │  │   │  │  ┌──────────┐          │    │    │        ▼         │
│ └─────────────┘ │  │   │  │  │  Read    │          │    │    │ ┌──────────────┐ │
│                 │  │   │  │  │ Replica  │          │    │    │ │ Background   │ │
│ Cache Strategy: │  │   │  │  │    N     │          │    │    │ │ Workers/     │ │
│ ┌─────────────┐ │  │   │  │  └──────────┘          │    │    │ │ Consumers    │ │
│ │Cache-Aside: │ │  │   │  │                        │    │    │ │              │ │
│ │             │ │  │   │  │ • ACID transactions    │    │    │ │ • Email      │ │
│ │1.Check cache│─┼──┘   │  │ • Indexing             │    │    │ │   sender     │ │
│ │2.Cache miss │ │      │  │ • Connection pooling   │    │    │ │ • Image      │ │
│ │3.Query DB   │─┼──────┼─▶│                        │    │    │ │   processor  │ │
│ │4.Write cache│◄┼──────┘  │                        │    │    │ │ • Analytics  │ │
│ │5.Return data│ │         │                        │    │    │ │ • Reporting  │ │
│ └─────────────┘ │         └────────────────────────┘    │    │ └──────────────┘ │
│                 │                                       │    │                  │
│                 │         ┌────────────────────────┐    │    │ Retry Logic:     │
│                 │         │  NoSQL DB              │    │    │ • Max retries: 3 │
│                 │         │  (MongoDB / Cassandra/ │    │    │ • Exponential    │
│                 │         │   DynamoDB)            │    │    │   backoff        │
│                 │         │                        │    │    │ • Dead letter    │
│                 │         │  ┌──────────┐          │    │    │   queue          │
│                 │         │  │ Shard 1  │          │    │    └──────────────────┘
│                 │         │  │ (A-F)    │          │    │
│                 │         │  └──────────┘          │    │
│                 │         │                        │    │
│                 │         │  ┌──────────┐          │    │
│                 │         │  │ Shard 2  │          │    │
│                 │         │  │ (G-M)    │          │    │
│                 │         │  └──────────┘          │    │
│                 │         │                        │    │
│                 │         │  ┌──────────┐          │    │
│                 │         │  │ Shard N  │          │    │
│                 │         │  │ (N-Z)    │          │    │
│                 │         │  └──────────┘          │    │
│                 │         │                        │    │
│                 │         │ • Horizontal scaling   │    │
│                 │         │ • Eventual consistency │    │
│                 │         │ • High throughput      │    │
│                 │         └────────────────────────┘    │
└─────────────────┘                                       └──────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         OBSERVABILITY & MONITORING                                  │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐                 │
│  │ Centralized      │  │ Metrics &        │  │ Distributed      │                 │
│  │ Logging          │  │ Monitoring       │  │ Tracing          │                 │
│  │                  │  │                  │  │                  │                 │
│  │ • ELK Stack      │  │ • Prometheus     │  │ • OpenTelemetry  │                 │
│  │ • OpenSearch     │  │ • Grafana        │  │ • Jaeger         │                 │
│  │ • Splunk         │  │ • Datadog        │  │ • Zipkin         │                 │
│  │                  │  │                  │  │                  │                 │
│  │ Logs:            │  │ Metrics:         │  │ Traces:          │                 │
│  │ • Application    │  │ • CPU, Memory    │  │ • Request ID     │                 │
│  │ • Access logs    │  │ • Request rate   │  │ • Span tracking  │                 │
│  │ • Error logs     │  │ • Error rate     │  │ • Latency        │                 │
│  │ • Audit logs     │  │ • Latency p95/99 │  │   breakdown      │                 │
│  │                  │  │ • DB connections │  │ • Bottleneck ID  │                 │
│  └──────────────────┘  └────────┬─────────┘  └──────────────────┘                 │
│                                 │                                                   │
│                                 ▼                                                   │
│                    ┌────────────────────────┐                                      │
│                    │  Alerting System       │                                      │
│                    │  (PagerDuty / Opsgenie)│                                      │
│                    │                        │                                      │
│                    │  • Threshold alerts    │                                      │
│                    │  • Anomaly detection   │                                      │
│                    │  • On-call rotation    │                                      │
│                    │  • Incident management │                                      │
│                    └────────────────────────┘                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                            SECURITY & CONFIGURATION                                 │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐                 │
│  │ Secrets Manager  │  │ Configuration    │  │ IAM & Access     │                 │
│  │                  │  │ Management       │  │ Control          │                 │
│  │ • Vault          │  │                  │  │                  │                 │
│  │ • AWS Secrets    │  │ • Consul KV      │  │ • IAM Roles      │                 │
│  │   Manager        │  │ • etcd           │  │ • Policies       │                 │
│  │ • Azure KeyVault │  │ • ConfigMaps     │  │ • MFA            │                 │
│  │                  │  │                  │  │ • Least privilege│                 │
│  │ Stores:          │  │ Stores:          │  │                  │                 │
│  │ • DB credentials │  │ • Feature flags  │  │ Service-to-      │                 │
│  │ • API keys       │  │ • Environment    │  │ service auth:    │                 │
│  │ • Certificates   │  │   variables      │  │ • mTLS           │                 │
│  │ • Encryption keys│  │ • App settings   │  │ • Service mesh   │                 │
│  │                  │  │                  │  │   (Istio)        │                 │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘                 │
│                                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐                 │
│  │  Encryption                                                   │                 │
│  │  • In Transit: TLS 1.3, mTLS                                 │                 │
│  │  • At Rest: AES-256, Database encryption                     │                 │
│  │  • End-to-End: Client-side encryption for sensitive data     │                 │
│  └──────────────────────────────────────────────────────────────┘                 │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Complete Request Flow (Detailed)

### READ Operation Flow (e.g., GET /api/products/123)

```
1. Client Request
   │
   ├─ Browser sends: GET /api/products/123
   │  Headers: Authorization: Bearer <JWT>
   │
2. DNS Resolution
   │
   ├─ DNS resolves api.example.com → CDN IP
   │
3. CDN Check
   │
   ├─ Is request cacheable? (GET with cache headers)
   ├─ Cache HIT? → Return cached response (END)
   └─ Cache MISS? → Continue to origin
   │
4. WAF Security Check
   │
   ├─ Check for malicious patterns
   ├─ Validate request size/headers
   └─ Pass → Continue, Fail → 403 Forbidden
   │
5. Rate Limiter
   │
   ├─ Check user rate limit (100 req/min)
   ├─ Within limit? → Continue
   └─ Exceeded? → 429 Too Many Requests
   │
6. Load Balancer
   │
   ├─ Select healthy backend (health check ✓)
   ├─ Algorithm: Least connections
   └─ Route to API Gateway instance
   │
7. API Gateway
   │
   ├─ Verify JWT token signature
   ├─ Check user permissions (ACL)
   ├─ Request validation
   ├─ Log request (correlation ID)
   └─ Route to Product Service
   │
8. Product Service (Application Layer)
   │
   ├─ Receive request with trace ID
   ├─ Check Redis cache first
   │  │
   │  ├─ Cache HIT?
   │  │  └─ Return cached data (skip DB) → Step 12
   │  │
   │  └─ Cache MISS?
   │     └─ Continue to DB
   │
9. Database Query
   │
   ├─ Query: SELECT * FROM products WHERE id = 123
   ├─ Read from Read Replica (reduce load on primary)
   ├─ Connection from pool (avoid connection overhead)
   └─ Return result
   │
10. Cache Write (Cache-Aside Pattern)
   │
   ├─ Write data to Redis
   ├─ Set TTL: 1 hour
   └─ Continue
   │
11. Service Response Processing
   │
   ├─ Transform data (DTO)
   ├─ Log response time
   └─ Return JSON response
   │
12. API Gateway Response
   │
   ├─ Add response headers (CORS, Cache-Control)
   ├─ Log response (status, latency)
   └─ Return to client
   │
13. Load Balancer → Client
   │
   └─ Response: 200 OK + JSON data
   
⏱️ Total Time: ~50-150ms (with cache: 5-20ms)
```

### WRITE Operation Flow (e.g., POST /api/orders)

```
1-7. [Same as READ: Client → DNS → CDN → WAF → Rate Limiter → LB → API Gateway]
   │
8. Order Service (Application Layer)
   │
   ├─ Validate request payload
   ├─ Check inventory (call Product Service via gRPC)
   │  └─ Circuit Breaker: If Product Service down → Fail fast
   │
9. Database Write (Primary DB)
   │
   ├─ BEGIN TRANSACTION
   ├─ INSERT INTO orders (...) VALUES (...)
   ├─ UPDATE inventory SET quantity = quantity - 1
   ├─ COMMIT
   └─ Return order ID
   │
10. Cache Invalidation
   │
   ├─ Delete related cache keys
   │  └─ DEL product:123
   │  └─ DEL user:orders:456
   │
11. Publish Event (Async Communication)
   │
   ├─ Publish to Kafka: "order.created"
   │  Payload: { orderId, userId, items, amount }
   │
12. Background Processing (Async Workers)
   │
   ├─ Email Service: Consumes "order.created" → Send confirmation email
   ├─ Payment Service: Process payment → Publish "payment.completed"
   ├─ Analytics Service: Update metrics
   └─ Notification Service: Send push notification
   │
13. Response to Client
   │
   └─ 201 Created + { orderId: 789, status: "pending" }
   
⏱️ Total Time: ~200-500ms (user gets immediate response, async tasks continue)
```

---

## 🎯 Scalability Patterns Applied

### Horizontal Scaling
```
Application Layer:
✅ Stateless services (any instance can handle any request)
✅ Auto-scaling groups (scale based on CPU/memory/request count)
✅ No session affinity required

Data Layer:
✅ Database sharding (partition by user_id, region, etc.)
✅ Read replicas (scale reads independently)
✅ NoSQL for massive scale
```

### Caching Strategy
```
Multi-Level Caching:
1. Browser/App cache (client-side)
2. CDN (edge cache)
3. Redis/Memcached (server-side)
4. Database query cache

Cache Invalidation:
• Time-based (TTL)
• Event-based (on write operations)
• Manual purge (admin action)
```

### Async Processing
```
Why Async?
✅ Decouple services (loose coupling)
✅ Handle traffic spikes (queue buffering)
✅ Retry failed operations
✅ Scale consumers independently

Examples:
• Send email (not urgent)
• Generate reports (heavy CPU)
• Image processing (time-consuming)
• Analytics tracking (fire-and-forget)
```

---

## 🛡️ Reliability Patterns Applied

### Circuit Breaker
```
Problem: Service A calls Service B, B is down → A keeps trying → Cascading failure

Solution:
┌─────────────────────────┐
│ Circuit States:         │
│                         │
│ 1. CLOSED (normal)      │
│    ↓ failures > threshold│
│ 2. OPEN (fail fast)     │
│    ↓ timeout period     │
│ 3. HALF-OPEN (test)     │
│    ↓ success → CLOSED   │
│    ↓ failure → OPEN     │
└─────────────────────────┘

Configuration:
• Failure threshold: 50%
• Timeout: 60 seconds
• Half-open attempts: 3
```

### Retry with Backoff
```
Request failed? Don't retry immediately!

Exponential Backoff:
Attempt 1: Immediate
Attempt 2: Wait 1s
Attempt 3: Wait 2s
Attempt 4: Wait 4s
Attempt 5: Wait 8s

Add jitter (randomness) to avoid thundering herd:
Wait time = base_delay * (2 ^ attempt) + random(0, 1s)

Max retries: 3-5 (then fail)
```

### Health Checks
```
Application Health:
GET /health → 200 OK
{
  "status": "healthy",
  "database": "connected",
  "cache": "connected",
  "uptime": 3600
}

Load Balancer Health Checks:
• Every 10 seconds
• 2 consecutive failures → Mark unhealthy
• 3 consecutive successes → Mark healthy
• Unhealthy instance → Remove from pool
```

### Timeout Configuration
```
Service Chain: Client → API Gateway → Service A → Service B

Timeouts:
Client:          60s
API Gateway:     55s  (less than client)
Service A:       50s  (less than gateway)
Service B:       45s  (less than Service A)

Why? Prevent timeout cascading
Each layer should timeout before parent layer
```

---

## 🔐 Security Best Practices

### Defense in Depth
```
Layer 1: WAF (block common attacks)
Layer 2: Rate limiting (prevent abuse)
Layer 3: Authentication (verify identity)
Layer 4: Authorization (check permissions)
Layer 5: Input validation (prevent injection)
Layer 6: Encryption (protect data)
Layer 7: Audit logging (detect breaches)
```

### Zero Trust Architecture
```
Never trust, always verify

Principles:
✅ Verify every request (even internal)
✅ Least privilege access
✅ Assume breach (segment network)
✅ Log everything (audit trail)

Implementation:
• mTLS for service-to-service
• Service mesh (Istio, Linkerd)
• API keys with expiration
• Regular credential rotation
```

---

## 📊 Monitoring & Observability Strategy

### The Three Pillars

**1. Logs (What happened?)**
```
Structured logging:
{
  "timestamp": "2026-01-11T10:30:45Z",
  "level": "ERROR",
  "service": "order-service",
  "trace_id": "abc-123",
  "user_id": "user-456",
  "message": "Payment failed",
  "error": "Insufficient funds",
  "latency_ms": 250
}

Storage: Elasticsearch (7-30 days retention)
```

**2. Metrics (How much/many?)**
```
Key Metrics:
• Request rate (req/sec)
• Error rate (%)
• Latency (p50, p95, p99)
• Saturation (CPU, memory, disk, network)

RED Method:
R - Rate (requests per second)
E - Errors (error rate)
D - Duration (latency distribution)

USE Method:
U - Utilization (% busy)
S - Saturation (queue depth)
E - Errors (error count)
```

**3. Traces (Where's the bottleneck?)**
```
Distributed Trace Example:
Trace ID: abc-123

Span 1: API Gateway       [0-5ms]
Span 2: Order Service     [5-50ms]
  ├─ Span 3: Redis cache  [10-12ms]
  ├─ Span 4: Product API  [15-35ms]  ← Slow!
  └─ Span 5: DB write     [40-48ms]
Span 6: Kafka publish     [50-52ms]

Total: 52ms (Product API is bottleneck)
```

### Alerting Rules
```
Critical (Page immediately):
• Error rate > 5% for 5 minutes
• p99 latency > 1s for 5 minutes
• Service down (health check fails)
• Database connection pool exhausted

Warning (Notify on Slack):
• Error rate > 1% for 10 minutes
• p95 latency > 500ms for 10 minutes
• Cache hit rate < 70%
• Disk usage > 80%
```

---

## 💡 Interview Discussion Points

### Scalability Trade-offs

**Q: How do you scale from 1K to 10M users?**
```
1K users:
• Single server (monolith)
• Single database
• No caching needed

10K users:
• Split web + app + database
• Add load balancer
• Add Redis cache

100K users:
• Microservices
• Database read replicas
• CDN for static assets
• Async processing (message queue)

1M users:
• Auto-scaling groups
• Database sharding
• Distributed cache cluster
• Multiple regions

10M users:
• Global CDN
• Multi-region deployment
• Event-driven architecture
• Data partitioning strategies
• Dedicated teams per service
```

### Consistency vs Availability (CAP Theorem)

**Q: How do you handle database replication lag?**
```
Problem: Write to primary → Read from replica (lag: 100ms)

Solutions:
1. Read from primary after write (consistency > availability)
2. Use cache for recent writes
3. Add version/timestamp to detect stale reads
4. Eventually consistent (acceptable for non-critical data)

Example:
• Bank balance → Read from primary (strong consistency)
• Social media likes → Read from replica (eventual consistency OK)
```

### Failure Scenarios

**Q: What happens if Redis goes down?**
```
Without cache:
❌ All requests hit database → Database overload → Cascading failure

With proper design:
1. Circuit breaker trips (after 50% errors)
2. Requests bypass cache → Go directly to DB
3. Rate limiter protects DB (limit concurrent connections)
4. Database connection pooling prevents exhaustion
5. Read replicas distribute load
6. Auto-scaling adds DB read replicas
7. Alerts fired → Engineers notified
8. Redis recovered → Circuit breaker closes → Normal flow

Result: Degraded performance, but system survives
```

**Q: How do you handle traffic spikes (10x normal load)?**
```
1. Auto-scaling (horizontal scaling)
   • Trigger: CPU > 70% for 2 minutes
   • Action: Add instances (2 → 4 → 8)
   
2. CDN absorbs static content traffic
   
3. Rate limiting protects backend
   • Tiered limits: Free (10/min), Premium (100/min)
   
4. Queue buffering (message queue)
   • Accept requests fast
   • Process asynchronously
   
5. Database connection pooling
   • Limit max connections
   • Queue requests in application
   
6. Circuit breaker
   • Fail fast if services overwhelmed
   
7. Graceful degradation
   • Disable non-essential features
   • Show cached/static content
```

---

## 🎓 Key Takeaways for Interviews

### When Presenting This Architecture

**1. Start with Requirements**
```
"Let me first clarify the requirements:
• Expected users? (1M daily active)
• Read/write ratio? (90% read, 10% write)
• Latency requirements? (<200ms p95)
• Consistency needs? (eventual OK? or strong?)
• Global or regional? (single region to start)"
```

**2. Draw Incrementally**
```
Don't draw everything at once!

Step 1: Client → Server → Database
Step 2: Add load balancer (scale)
Step 3: Add cache (reduce DB load)
Step 4: Add microservices (separation of concerns)
Step 5: Add async processing (decouple)
Step 6: Add monitoring (observability)
```

**3. Discuss Trade-offs**
```
Every decision has trade-offs:

Microservices:
✅ Independent scaling
✅ Team autonomy
❌ Complexity (debugging, deployment)
❌ Network overhead

Caching:
✅ Reduced latency
✅ Reduced DB load
❌ Stale data risk
❌ Cache invalidation complexity

Async Processing:
✅ Better throughput
✅ Fault tolerance
❌ Eventual consistency
❌ Ordering challenges
```

**4. Numbers Matter**
```
Have rough numbers ready:
• Redis: ~10k-50k ops/sec per instance
• Database: ~1k-5k QPS per instance
• Load balancer: ~10k concurrent connections
• Kafka: ~100k-1M messages/sec
• CDN cache hit rate: 80-95%
```

**5. Address Failure Cases**
```
Interviewers love this:
• "What if the database goes down?"
• "What if there's a network partition?"
• "What if Redis cache is stale?"

Show you think about reliability!
```

---

## 🚀 Next Steps

To use this diagram effectively in interviews:

1. **Practice drawing** simplified versions on whiteboard
2. **Memorize key components** and their purposes
3. **Understand trade-offs** for each choice
4. **Prepare numbers** (throughput, latency, capacity)
5. **Study real systems** (Netflix, Uber, Twitter architectures)
6. **Ask clarifying questions** before jumping into design

**Remember:** There's no single "correct" architecture. Focus on demonstrating your thought process, trade-off analysis, and ability to adapt based on requirements.

Good luck with your interviews! 🎯
