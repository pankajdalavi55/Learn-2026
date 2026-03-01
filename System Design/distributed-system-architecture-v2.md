# Large-Scale Distributed System Architecture
> Clean whiteboard-style diagram for system design interviews

---

## 🏗️ High-Level Architecture Overview

```
╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                                    CLIENTS                                               ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║     🌐 Web Browser              📱 Mobile App               🖥️ Desktop App              ║
║     ┌─────────────┐             ┌─────────────┐             ┌─────────────┐             ║
║     │ • Cookies   │             │ • Keychain  │             │ • Token     │             ║
║     │ • LocalStore│             │ • SQLite    │             │   Storage   │             ║
║     │ • SessionStr│             │ • Shared    │             │ • File      │             ║
║     │ • Browser   │             │   Prefs     │             │   Cache     │             ║
║     │   Cache     │             │ • NSCache   │             │             │             ║
║     └──────┬──────┘             └──────┬──────┘             └──────┬──────┘             ║
║            │                           │                           │                    ║
║            └───────────────────────────┼───────────────────────────┘                    ║
║                                        │                                                 ║
║                                        ▼                                                 ║
║                          ┌─────────────────────────┐                                    ║
║                          │    HTTPS Request        │                                    ║
║                          │  GET /api/v1/products   │                                    ║
║                          │  Authorization: Bearer  │                                    ║
║                          └─────────────────────────┘                                    ║
║                                        │                                                 ║
╚════════════════════════════════════════╪════════════════════════════════════════════════╝
                                         │
                                         ▼
╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                              EDGE & SECURITY LAYER                                       ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║  ┌──────────────────────────────────────────────────────────────────────────────────┐   ║
║  │ 🌍 DNS (Route 53 / Cloudflare DNS)                                               │   ║
║  │    ├─ Geo-routing: Route users to nearest datacenter                             │   ║
║  │    ├─ Failover: Automatic switch on health check failure                         │   ║
║  │    └─ TTL: 300s (5 min) for quick failover                                       │   ║
║  └────────────────────────────────────┬─────────────────────────────────────────────┘   ║
║                                       │                                                  ║
║                                       ▼                                                  ║
║  ┌──────────────────────────────────────────────────────────────────────────────────┐   ║
║  │ 🚀 CDN (CloudFront / Cloudflare / Akamai)                                        │   ║
║  │    ├─ Edge Locations: 200+ PoPs globally                                         │   ║
║  │    ├─ Static Assets: JS, CSS, Images, Fonts, Videos                              │   ║
║  │    ├─ Cache-Control: max-age=31536000, immutable                                 │   ║
║  │    │                                                                              │   ║
║  │    │  ┌───────────┐  Cache HIT   ┌─────────────────────┐                         │   ║
║  │    └──│ Request   │─────────────►│ Return cached asset │ ──► Client              │   ║
║  │       └───────────┘              └─────────────────────┘                         │   ║
║  │             │                                                                     │   ║
║  │             │ Cache MISS                                                          │   ║
║  │             ▼                                                                     │   ║
║  └─────────────┼────────────────────────────────────────────────────────────────────┘   ║
║                │                                                                         ║
║                ▼                                                                         ║
║  ┌──────────────────────────────────────────────────────────────────────────────────┐   ║
║  │ 🛡️ WAF (AWS WAF / Cloudflare WAF)                                                │   ║
║  │    ├─ OWASP Top 10 protection                                                    │   ║
║  │    ├─ SQL injection blocking                                                     │   ║
║  │    ├─ XSS prevention                                                             │   ║
║  │    ├─ Bot detection & CAPTCHA                                                    │   ║
║  │    └─ Geo-blocking (if required)                                                 │   ║
║  │                                                                                   │   ║
║  │    BLOCK ──► 403 Forbidden                                                       │   ║
║  │    PASS  ──► Continue ▼                                                          │   ║
║  └─────────────┬────────────────────────────────────────────────────────────────────┘   ║
║                │                                                                         ║
║                ▼                                                                         ║
║  ┌──────────────────────────────────────────────────────────────────────────────────┐   ║
║  │ ⏱️ RATE LIMITER                                                                   │   ║
║  │                                                                                   │   ║
║  │    Algorithm: Token Bucket / Sliding Window                                      │   ║
║  │                                                                                   │   ║
║  │    ┌─────────────────┬──────────────────┬──────────────────┐                     │   ║
║  │    │  Tier           │  Limit           │  Window          │                     │   ║
║  │    ├─────────────────┼──────────────────┼──────────────────┤                     │   ║
║  │    │  Anonymous      │  20 req          │  per minute      │                     │   ║
║  │    │  Free User      │  100 req         │  per minute      │                     │   ║
║  │    │  Premium User   │  1000 req        │  per minute      │                     │   ║
║  │    │  Per IP         │  500 req         │  per minute      │                     │   ║
║  │    └─────────────────┴──────────────────┴──────────────────┘                     │   ║
║  │                                                                                   │   ║
║  │    EXCEEDED ──► 429 Too Many Requests + Retry-After header                       │   ║
║  │    OK       ──► Continue ▼                                                       │   ║
║  └─────────────┬────────────────────────────────────────────────────────────────────┘   ║
║                │                                                                         ║
╚════════════════╪════════════════════════════════════════════════════════════════════════╝
                 │
                 ▼
╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                           API & TRAFFIC MANAGEMENT                                       ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║  ┌──────────────────────────────────────────────────────────────────────────────────┐   ║
║  │ ⚖️ LOAD BALANCER (ALB / NLB / nginx / HAProxy)                                   │   ║
║  │                                                                                   │   ║
║  │    Layer 7 (Application):                                                        │   ║
║  │    ├─ SSL/TLS termination                                                        │   ║
║  │    ├─ Path-based routing (/api/*, /admin/*, /ws/*)                               │   ║
║  │    ├─ Host-based routing (api.example.com, admin.example.com)                    │   ║
║  │    └─ WebSocket support                                                          │   ║
║  │                                                                                   │   ║
║  │    Health Checks:                                                                │   ║
║  │    ├─ Endpoint: GET /health                                                      │   ║
║  │    ├─ Interval: 10 seconds                                                       │   ║
║  │    ├─ Unhealthy threshold: 2 failures                                            │   ║
║  │    └─ Healthy threshold: 3 successes                                             │   ║
║  │                                                                                   │   ║
║  │    Algorithm: Least Connections / Round Robin / Weighted                         │   ║
║  │                                                                                   │   ║
║  │         ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐                   │   ║
║  │         │ Server 1 │  │ Server 2 │  │ Server 3 │  │ Server N │                   │   ║
║  │         │    ✓     │  │    ✓     │  │    ✗     │  │    ✓     │                   │   ║
║  │         └────┬─────┘  └────┬─────┘  └──────────┘  └────┬─────┘                   │   ║
║  │              │             │         (unhealthy)       │                          │   ║
║  └──────────────┼─────────────┼───────────────────────────┼─────────────────────────┘   ║
║                 │             │                           │                              ║
║                 └─────────────┴─────────────┬─────────────┘                              ║
║                                             │                                            ║
║                                             ▼                                            ║
║  ┌──────────────────────────────────────────────────────────────────────────────────┐   ║
║  │ 🚪 API GATEWAY (Kong / AWS API Gateway / Apigee)                                 │   ║
║  │                                                                                   │   ║
║  │    ┌─────────────────────────────────────────────────────────────────────────┐   │   ║
║  │    │  1. AUTHENTICATION                                                      │   │   ║
║  │    │     ├─ JWT validation (verify signature, expiry)                        │   │   ║
║  │    │     ├─ OAuth 2.0 token introspection                                    │   │   ║
║  │    │     ├─ API Key validation                                               │   │   ║
║  │    │     └─ INVALID ──► 401 Unauthorized                                     │   │   ║
║  │    └─────────────────────────────────────────────────────────────────────────┘   │   ║
║  │                                             │                                     │   ║
║  │    ┌─────────────────────────────────────────────────────────────────────────┐   │   ║
║  │    │  2. AUTHORIZATION                                                       │   │   ║
║  │    │     ├─ Role-based access control (RBAC)                                 │   │   ║
║  │    │     ├─ Scope validation (read, write, admin)                            │   │   ║
║  │    │     └─ FORBIDDEN ──► 403 Forbidden                                      │   │   ║
║  │    └─────────────────────────────────────────────────────────────────────────┘   │   ║
║  │                                             │                                     │   ║
║  │    ┌─────────────────────────────────────────────────────────────────────────┐   │   ║
║  │    │  3. REQUEST TRANSFORMATION                                              │   │   ║
║  │    │     ├─ Add correlation ID (X-Request-ID)                                │   │   ║
║  │    │     ├─ Version routing (v1, v2, v3)                                     │   │   ║
║  │    │     ├─ Request validation (JSON schema)                                 │   │   ║
║  │    │     └─ Header enrichment                                                │   │   ║
║  │    └─────────────────────────────────────────────────────────────────────────┘   │   ║
║  │                                             │                                     │   ║
║  │    ┌─────────────────────────────────────────────────────────────────────────┐   │   ║
║  │    │  4. SERVICE ROUTING                                                     │   │   ║
║  │    │     ├─ /api/v1/users/*    ──► User Service                              │   │   ║
║  │    │     ├─ /api/v1/products/* ──► Product Service                           │   │   ║
║  │    │     ├─ /api/v1/orders/*   ──► Order Service                             │   │   ║
║  │    │     └─ /api/v1/payments/* ──► Payment Service                           │   │   ║
║  │    └─────────────────────────────────────────────────────────────────────────┘   │   ║
║  │                                                                                   │   ║
║  └──────────────────────────────────────────┬───────────────────────────────────────┘   ║
║                                             │                                            ║
╚═════════════════════════════════════════════╪════════════════════════════════════════════╝
                                              │
                                              ▼
╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                              APPLICATION LAYER                                           ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║    ┌───────────────────────────────────────────────────────────────────────────────┐    ║
║    │                        MICROSERVICES CLUSTER                                  │    ║
║    │                                                                               │    ║
║    │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐      │    ║
║    │   │   👤 USER   │   │  📦 PRODUCT │   │  🛒 ORDER   │   │  💳 PAYMENT │      │    ║
║    │   │   SERVICE   │   │   SERVICE   │   │   SERVICE   │   │   SERVICE   │      │    ║
║    │   │             │   │             │   │             │   │             │      │    ║
║    │   │ Instances:  │   │ Instances:  │   │ Instances:  │   │ Instances:  │      │    ║
║    │   │   3-10      │   │   5-20      │   │   3-15      │   │   2-5       │      │    ║
║    │   │             │   │             │   │             │   │             │      │    ║
║    │   │ • Stateless │   │ • Stateless │   │ • Stateless │   │ • Stateless │      │    ║
║    │   │ • Auto-scale│   │ • Auto-scale│   │ • Auto-scale│   │ • Auto-scale│      │    ║
║    │   │ • Health ✓  │   │ • Health ✓  │   │ • Health ✓  │   │ • Health ✓  │      │    ║
║    │   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘      │    ║
║    │          │                 │                 │                 │              │    ║
║    │          └─────────────────┴─────────────────┴─────────────────┘              │    ║
║    │                                    │                                          │    ║
║    └────────────────────────────────────┼──────────────────────────────────────────┘    ║
║                                         │                                               ║
║    ┌────────────────────────────────────┴──────────────────────────────────────────┐    ║
║    │                     INTER-SERVICE COMMUNICATION                               │    ║
║    │                                                                               │    ║
║    │   ┌─────────────────────────┐    ┌─────────────────────────┐                 │    ║
║    │   │   SYNCHRONOUS           │    │   ASYNCHRONOUS          │                 │    ║
║    │   │                         │    │                         │                 │    ║
║    │   │  ┌───────────────────┐  │    │  ┌───────────────────┐  │                 │    ║
║    │   │  │ REST (HTTP/JSON)  │  │    │  │ Message Queue     │  │                 │    ║
║    │   │  │ • Simple          │  │    │  │ • Kafka           │  │                 │    ║
║    │   │  │ • Widely adopted  │  │    │  │ • RabbitMQ        │  │                 │    ║
║    │   │  │ • Cacheable       │  │    │  │ • Amazon SQS      │  │                 │    ║
║    │   │  └───────────────────┘  │    │  └───────────────────┘  │                 │    ║
║    │   │                         │    │                         │                 │    ║
║    │   │  ┌───────────────────┐  │    │  ┌───────────────────┐  │                 │    ║
║    │   │  │ gRPC (HTTP/2)     │  │    │  │ Event Streaming   │  │                 │    ║
║    │   │  │ • Fast (binary)   │  │    │  │ • Event sourcing  │  │                 │    ║
║    │   │  │ • Strongly typed  │  │    │  │ • CQRS pattern    │  │                 │    ║
║    │   │  │ • Bi-directional  │  │    │  │ • Replay capable  │  │                 │    ║
║    │   │  └───────────────────┘  │    │  └───────────────────┘  │                 │    ║
║    │   │                         │    │                         │                 │    ║
║    │   └─────────────────────────┘    └─────────────────────────┘                 │    ║
║    │                                                                               │    ║
║    └───────────────────────────────────────────────────────────────────────────────┘    ║
║                                                                                          ║
║    ┌─────────────────────────────────────────────────────────────────────────────────┐  ║
║    │  🔍 SERVICE DISCOVERY (Consul / Eureka / Kubernetes DNS)                        │  ║
║    │                                                                                 │  ║
║    │     Service Registry:                                                           │  ║
║    │     ┌────────────────┬─────────────────────────┬──────────────┐                │  ║
║    │     │ Service        │ Instances               │ Health       │                │  ║
║    │     ├────────────────┼─────────────────────────┼──────────────┤                │  ║
║    │     │ user-service   │ 10.0.1.10, 10.0.1.11   │ ✓ ✓          │                │  ║
║    │     │ product-service│ 10.0.2.10, 10.0.2.11   │ ✓ ✓          │                │  ║
║    │     │ order-service  │ 10.0.3.10, 10.0.3.11   │ ✓ ✓          │                │  ║
║    │     └────────────────┴─────────────────────────┴──────────────┘                │  ║
║    │                                                                                 │  ║
║    │     • Automatic registration on startup                                        │  ║
║    │     • Health check every 10 seconds                                            │  ║
║    │     • Auto-deregistration on failure                                           │  ║
║    └─────────────────────────────────────────────────────────────────────────────────┘  ║
║                                                                                          ║
╚══════════════════════════════════════════════════════════════════════════════════════════╝
                    │                              │
                    ▼                              ▼
    ┌───────────────────────────┐    ┌───────────────────────────┐
    │                           │    │                           │
    ▼                           ▼    ▼                           ▼
╔═══════════════════════╗  ╔═══════════════════════╗  ╔═══════════════════════╗
║    CACHING LAYER      ║  ║     DATA LAYER        ║  ║   ASYNC & EVENTING    ║
╠═══════════════════════╣  ╠═══════════════════════╣  ╠═══════════════════════╣
║                       ║  ║                       ║  ║                       ║
║  ┌─────────────────┐  ║  ║  ┌─────────────────┐  ║  ║  ┌─────────────────┐  ║
║  │   REDIS CLUSTER │  ║  ║  │  RELATIONAL DB  │  ║  ║  │  MESSAGE BROKER │  ║
║  │                 │  ║  ║  │  (PostgreSQL)   │  ║  ║  │  (Apache Kafka) │  ║
║  │  ┌───────────┐  │  ║  ║  │                 │  ║  ║  │                 │  ║
║  │  │  Master   │  │  ║  ║  │   ┌─────────┐   │  ║  ║  │  Topics:        │  ║
║  │  │   (R/W)   │  │  ║  ║  │   │ PRIMARY │   │  ║  ║  │  • orders       │  ║
║  │  └─────┬─────┘  │  ║  ║  │   │  (R/W)  │   │  ║  ║  │  • payments     │  ║
║  │        │        │  ║  ║  │   └────┬────┘   │  ║  ║  │  • notifications║  ║
║  │        ▼        │  ║  ║  │        │        │  ║  ║  │  • analytics    │  ║
║  │  ┌───────────┐  │  ║  ║  │  Replication    │  ║  ║  │                 │  ║
║  │  │ Replica 1 │  │  ║  ║  │        │        │  ║  ║  │  Partitions: 12 │  ║
║  │  │  (Read)   │  │  ║  ║  │        ▼        │  ║  ║  │  Replication: 3 │  ║
║  │  └───────────┘  │  ║  ║  │  ┌─────────┐    │  ║  ║  │                 │  ║
║  │  ┌───────────┐  │  ║  ║  │  │ REPLICA │    │  ║  ║  └────────┬────────┘  ║
║  │  │ Replica 2 │  │  ║  ║  │  │ (Read)  │    │  ║  ║           │          ║
║  │  │  (Read)   │  │  ║  ║  │  └─────────┘    │  ║  ║           ▼          ║
║  │  └───────────┘  │  ║  ║  │  ┌─────────┐    │  ║  ║  ┌─────────────────┐  ║
║  │                 │  ║  ║  │  │ REPLICA │    │  ║  ║  │    CONSUMERS    │  ║
║  │  Sentinel for   │  ║  ║  │  │ (Read)  │    │  ║  ║  │                 │  ║
║  │  auto-failover  │  ║  ║  │  └─────────┘    │  ║  ║  │  ┌───────────┐  │  ║
║  └─────────────────┘  ║  ║  └─────────────────┘  ║  ║  │  │   Email   │  │  ║
║                       ║  ║                       ║  ║  │  │  Worker   │  │  ║
║  ┌─────────────────┐  ║  ║  ┌─────────────────┐  ║  ║  │  └───────────┘  │  ║
║  │    MEMCACHED    │  ║  ║  │   NoSQL DB      │  ║  ║  │  ┌───────────┐  │  ║
║  │                 │  ║  ║  │   (MongoDB)     │  ║  ║  │  │  Analytics│  │  ║
║  │  • Session data │  ║  ║  │                 │  ║  ║  │  │  Worker   │  │  ║
║  │  • Temp data    │  ║  ║  │  ┌───────────┐  │  ║  ║  │  └───────────┘  │  ║
║  │  • Short TTL    │  ║  ║  │  │  Shard 1  │  │  ║  ║  │  ┌───────────┐  │  ║
║  │                 │  ║  ║  │  │  (A-M)    │  │  ║  ║  │  │  Image    │  │  ║
║  └─────────────────┘  ║  ║  │  └───────────┘  │  ║  ║  │  │ Processor │  │  ║
║                       ║  ║  │  ┌───────────┐  │  ║  ║  │  └───────────┘  │  ║
║ ┌───────────────────┐ ║  ║  │  │  Shard 2  │  │  ║  ║  │  ┌───────────┐  │  ║
║ │ CACHE-ASIDE FLOW  │ ║  ║  │  │  (N-Z)    │  │  ║  ║  │  │  Report   │  │  ║
║ │                   │ ║  ║  │  └───────────┘  │  ║  ║  │  │ Generator │  │  ║
║ │  1. Check Cache   │ ║  ║  │                 │  ║  ║  │  └───────────┘  │  ║
║ │       │           │ ║  ║  │  • Document DB  │  ║  ║  │                 │  ║
║ │       ▼           │ ║  ║  │  • Schema-less  │  ║  ║  └─────────────────┘  ║
║ │  HIT? ─► Return   │ ║  ║  │  • Horizontal   │  ║  ║                       ║
║ │       │           │ ║  ║  │    scaling      │  ║  ║  ┌─────────────────┐  ║
║ │       ▼ MISS      │ ║  ║  └─────────────────┘  ║  ║  │  DEAD LETTER Q  │  ║
║ │  2. Query DB      │ ║  ║                       ║  ║  │                 │  ║
║ │       │           │ ║  ║  ┌─────────────────┐  ║  ║  │  Failed msgs    │  ║
║ │       ▼           │ ║  ║  │  OBJECT STORAGE │  ║  ║  │  for manual     │  ║
║ │  3. Write Cache   │ ║  ║  │     (S3)        │  ║  ║  │  review         │  ║
║ │       │           │ ║  ║  │                 │  ║  ║  │                 │  ║
║ │       ▼           │ ║  ║  │  • Images       │  ║  ║  │  Max retries: 3 │  ║
║ │  4. Return Data   │ ║  ║  │  • Videos       │  ║  ║  │  Backoff: exp   │  ║
║ │                   │ ║  ║  │  • Documents    │  ║  ║  └─────────────────┘  ║
║ └───────────────────┘ ║  ║  │  • Backups      │  ║  ║                       ║
║                       ║  ║  └─────────────────┘  ║  ║                       ║
╚═══════════════════════╝  ╚═══════════════════════╝  ╚═══════════════════════╝


╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                               OBSERVABILITY LAYER                                        ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║    ┌─────────────────────┐   ┌─────────────────────┐   ┌─────────────────────┐          ║
║    │ 📋 LOGGING          │   │ 📊 METRICS          │   │ 🔎 TRACING          │          ║
║    │                     │   │                     │   │                     │          ║
║    │  Stack:             │   │  Stack:             │   │  Stack:             │          ║
║    │  • Elasticsearch    │   │  • Prometheus       │   │  • OpenTelemetry    │          ║
║    │  • Logstash         │   │  • Grafana          │   │  • Jaeger           │          ║
║    │  • Kibana           │   │  • Alertmanager     │   │  • Zipkin           │          ║
║    │                     │   │                     │   │                     │          ║
║    │  Log Types:         │   │  Key Metrics:       │   │  Features:          │          ║
║    │  • Application      │   │  • Request rate     │   │  • Distributed      │          ║
║    │  • Access           │   │  • Error rate       │   │    trace context    │          ║
║    │  • Error            │   │  • Latency p50/95/99│   │  • Span timing      │          ║
║    │  • Audit            │   │  • CPU/Memory       │   │  • Service map      │          ║
║    │  • Security         │   │  • DB connections   │   │  • Bottleneck ID    │          ║
║    │                     │   │  • Cache hit rate   │   │                     │          ║
║    │  Retention: 30 days │   │  Retention: 15 days │   │  Retention: 7 days  │          ║
║    └──────────┬──────────┘   └──────────┬──────────┘   └──────────┬──────────┘          ║
║               │                         │                         │                      ║
║               └─────────────────────────┼─────────────────────────┘                      ║
║                                         ▼                                                ║
║                          ┌────────────────────────────┐                                 ║
║                          │   🚨 ALERTING SYSTEM       │                                 ║
║                          │   (PagerDuty / Opsgenie)   │                                 ║
║                          │                            │                                 ║
║                          │  Alert Rules:              │                                 ║
║                          │  • Error rate > 5%    → P1 │                                 ║
║                          │  • Latency p99 > 1s   → P2 │                                 ║
║                          │  • Service down       → P1 │                                 ║
║                          │  • DB conn pool > 80% → P3 │                                 ║
║                          │                            │                                 ║
║                          │  ──► Slack / Email / SMS   │                                 ║
║                          │  ──► On-call rotation      │                                 ║
║                          └────────────────────────────┘                                 ║
║                                                                                          ║
╚══════════════════════════════════════════════════════════════════════════════════════════╝


╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                         RELIABILITY & MAINTAINABILITY                                    ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║  ┌────────────────────────┐  ┌────────────────────────┐  ┌────────────────────────┐     ║
║  │ 🔌 CIRCUIT BREAKER     │  │ 🔄 RETRY & TIMEOUT     │  │ ❤️ HEALTH CHECKS       │     ║
║  │                        │  │                        │  │                        │     ║
║  │  States:               │  │  Retry Strategy:       │  │  Endpoints:            │     ║
║  │  • CLOSED (normal)     │  │  • Max attempts: 3     │  │  • /health             │     ║
║  │  • OPEN (fail fast)    │  │  • Backoff: exponential│  │  • /ready              │     ║
║  │  • HALF-OPEN (test)    │  │  • Jitter: random      │  │  • /live               │     ║
║  │                        │  │                        │  │                        │     ║
║  │  Config:               │  │  Timeouts:             │  │  Checks:               │     ║
║  │  • Failure %: 50       │  │  • Connection: 5s      │  │  • Database conn       │     ║
║  │  • Timeout: 30s        │  │  • Read: 30s           │  │  • Redis conn          │     ║
║  │  • Half-open req: 3    │  │  • Total: 60s          │  │  • External API        │     ║
║  │                        │  │                        │  │  • Disk space          │     ║
║  │  Fallback:             │  │  Wait Times:           │  │                        │     ║
║  │  • Return cached data  │  │  • 1s → 2s → 4s → fail │  │  Interval: 10s         │     ║
║  │  • Return default      │  │                        │  │                        │     ║
║  │  • Graceful degrade    │  │                        │  │                        │     ║
║  └────────────────────────┘  └────────────────────────┘  └────────────────────────┘     ║
║                                                                                          ║
╚══════════════════════════════════════════════════════════════════════════════════════════╝


╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                             SECURITY & CONFIGURATION                                     ║
╠══════════════════════════════════════════════════════════════════════════════════════════╣
║                                                                                          ║
║  ┌────────────────────────┐  ┌────────────────────────┐  ┌────────────────────────┐     ║
║  │ 🔐 SECRETS MANAGER     │  │ ⚙️ CONFIG MANAGEMENT   │  │ 🛡️ IAM & ACCESS        │     ║
║  │                        │  │                        │  │                        │     ║
║  │  Tools:                │  │  Tools:                │  │  Principles:           │     ║
║  │  • HashiCorp Vault     │  │  • Consul KV           │  │  • Least privilege     │     ║
║  │  • AWS Secrets Manager │  │  • etcd                │  │  • Role-based access   │     ║
║  │  • Azure Key Vault     │  │  • Spring Cloud Config │  │  • Service accounts    │     ║
║  │                        │  │  • K8s ConfigMaps      │  │                        │     ║
║  │  Stores:               │  │                        │  │  Service-to-Service:   │     ║
║  │  • DB passwords        │  │  Stores:               │  │  • mTLS               │     ║
║  │  • API keys            │  │  • Feature flags       │  │  • Service mesh       │     ║
║  │  • TLS certificates    │  │  • Environment config  │  │  • JWT validation     │     ║
║  │  • Encryption keys     │  │  • App settings        │  │                        │     ║
║  │                        │  │                        │  │  Audit:                │     ║
║  │  Rotation: 30 days     │  │  Hot reload: Yes       │  │  • All access logged   │     ║
║  │                        │  │                        │  │  • Anomaly detection   │     ║
║  └────────────────────────┘  └────────────────────────┘  └────────────────────────┘     ║
║                                                                                          ║
║  ┌───────────────────────────────────────────────────────────────────────────────────┐  ║
║  │  🔒 ENCRYPTION                                                                     │  ║
║  │                                                                                    │  ║
║  │  In Transit:                      At Rest:                                        │  ║
║  │  • TLS 1.3 for all connections    • AES-256 for database                          │  ║
║  │  • mTLS for service-to-service    • Encryption for S3 objects                     │  ║
║  │  • Certificate rotation           • Key rotation every 90 days                    │  ║
║  │                                                                                    │  ║
║  └───────────────────────────────────────────────────────────────────────────────────┘  ║
║                                                                                          ║
╚══════════════════════════════════════════════════════════════════════════════════════════╝
```

---

## 📡 Request/Response Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────────────────────┐
│                              END-TO-END REQUEST FLOW                                     │
└──────────────────────────────────────────────────────────────────────────────────────────┘

     CLIENT                                                                 BACKEND
        │                                                                      │
        │  ①  HTTPS Request                                                   │
        │  ─────────────────────────────────────────────────────────────────► │
        │     GET /api/v1/products/123                                        │
        │     Headers: Authorization: Bearer eyJhbG...                        │
        │                                                                      │
        │                          ┌─────────────────────┐                    │
        │                          │     DNS LOOKUP      │                    │
        │                          │  api.example.com    │                    │
        │                          │  → 52.10.20.30      │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │       CDN           │                    │
        │                          │  Cache: MISS        │                    │
        │                          │  → Forward to origin│                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │       WAF           │                    │
        │                          │  Security: PASS     │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │   RATE LIMITER      │                    │
        │                          │  Limit: OK          │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │   LOAD BALANCER     │                    │
        │                          │  Select: Server 2   │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │    API GATEWAY      │                    │
        │                          │  ② Auth: VALID      │                    │
        │                          │  ③ Route: Product   │                    │
        │                          │     Service         │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │  PRODUCT SERVICE    │                    │
        │                          │  ④ Check Redis      │                    │
        │                          │     Cache: MISS     │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │    DATABASE         │                    │
        │                          │  ⑤ Query: SELECT    │                    │
        │                          │     Read Replica    │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │                                    ▼                                 │
        │                          ┌─────────────────────┐                    │
        │                          │  PRODUCT SERVICE    │                    │
        │                          │  ⑥ Write to Redis   │                    │
        │                          │     TTL: 1 hour     │                    │
        │                          └─────────────────────┘                    │
        │                                    │                                 │
        │  ⑦  HTTPS Response                │                                 │
        │  ◄───────────────────────────────────────────────────────────────── │
        │     Status: 200 OK                                                  │
        │     Headers: Cache-Control: max-age=300                             │
        │     Body: { "id": 123, "name": "Product", "price": 99.99 }         │
        │                                                                      │
        │                                                                      │
        ▼                                                                      ▼

     ┌──────────────────────────────────────────────────────────────────────────┐
     │  TIMELINE: Total ~120ms                                                  │
     │                                                                          │
     │  DNS Lookup      │████│ 20ms                                            │
     │  CDN             │██│ 5ms                                               │
     │  WAF             │██│ 5ms                                               │
     │  Rate Limiter    │██│ 2ms                                               │
     │  Load Balancer   │██│ 3ms                                               │
     │  API Gateway     │████│ 15ms                                            │
     │  Redis Check     │██│ 5ms (cache miss)                                  │
     │  Database Query  │████████████│ 50ms                                    │
     │  Redis Write     │██│ 5ms                                               │
     │  Response        │████│ 10ms                                            │
     │                  └──────────────────────────────────────────────────────│
     │                  0ms                 50ms                100ms     120ms │
     └──────────────────────────────────────────────────────────────────────────┘
```

---

## ⚡ Quick Reference Card

### Layer Summary

| Layer | Components | Purpose |
|-------|------------|---------|
| **Client** | Web/Mobile, Cache, Storage | User interface, local state |
| **Edge** | DNS, CDN, WAF, Rate Limiter | Security, performance, protection |
| **API** | Load Balancer, API Gateway | Traffic management, routing |
| **Application** | Microservices, Discovery | Business logic, processing |
| **Cache** | Redis, Memcached | Performance, reduce DB load |
| **Data** | SQL, NoSQL, Replicas | Persistence, durability |
| **Async** | Kafka, Workers | Decoupling, background jobs |
| **Observability** | Logs, Metrics, Traces | Monitoring, debugging |
| **Security** | Secrets, IAM, Encryption | Protection, compliance |

### Key Design Decisions

```
┌─────────────────────────────────────────────────────────────────────────┐
│  DECISION                          │  WHEN TO USE                       │
├─────────────────────────────────────────────────────────────────────────┤
│  SQL (PostgreSQL)                  │  ACID, relations, complex queries  │
│  NoSQL (MongoDB)                   │  Flexibility, horizontal scale     │
│  Redis Cache                       │  Hot data, sessions, rate limiting │
│  Kafka                             │  Event streaming, high throughput  │
│  REST API                          │  Simple, cacheable, widely adopted │
│  gRPC                              │  Low latency, internal services    │
│  Read Replicas                     │  Read-heavy workloads              │
│  Sharding                          │  Write scaling, data isolation     │
│  CDN                               │  Static assets, global users       │
│  Circuit Breaker                   │  Prevent cascading failures        │
└─────────────────────────────────────────────────────────────────────────┘
```

### Scaling Cheat Sheet

```
1K → 10K Users:
  └── Add caching (Redis)

10K → 100K Users:
  ├── Add load balancer
  ├── Add read replicas
  └── Add CDN

100K → 1M Users:
  ├── Microservices split
  ├── Message queues (async)
  ├── Auto-scaling
  └── Multiple regions

1M → 10M Users:
  ├── Database sharding
  ├── Global CDN
  ├── Event-driven architecture
  └── Dedicated services per domain
```

---

## 🎤 Interview Talking Points

### Opening Statement
> "Let me start with the high-level architecture and then dive into specific components based on the requirements."

### Scalability Question
> "To handle 10x traffic, I'd implement horizontal scaling with auto-scaling groups, add Redis caching to reduce database load, use CDN for static content, and introduce message queues to decouple services and handle traffic spikes."

### Reliability Question  
> "For high availability, I'd deploy across multiple availability zones, implement circuit breakers to prevent cascading failures, use health checks with automatic failover, and ensure all services are stateless so any instance can handle any request."

### Trade-off Discussion
> "There's always a trade-off. For example, with caching, we gain performance but risk serving stale data. With microservices, we gain scalability and team independence but add complexity in debugging and deployment. I'd start simple and add complexity only when needed."

### Failure Scenario
> "If the database goes down, the circuit breaker opens, we serve cached data where possible, and gracefully degrade non-critical features. Meanwhile, alerting notifies the on-call engineer within 30 seconds."

---

*This architecture is designed to be drawn incrementally on a whiteboard. Start with Client → Server → Database, then add layers based on requirements and constraints discussed with the interviewer.*
