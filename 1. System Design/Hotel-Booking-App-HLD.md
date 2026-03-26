# Hotel Booking Application — High-Level Design (HLD)

---

## 1. Problem Statement

Design a scalable hotel booking platform (similar to Booking.com, OYO, MakeMyTrip) that allows users to search for hotels, view available rooms, make reservations, process payments, and manage bookings.

---

## 2. Functional Requirements

| # | Requirement | Description |
|---|-------------|-------------|
| FR-1 | User Registration & Auth | Sign up, login, OAuth, profile management |
| FR-2 | Hotel/Room Search | Search by city, date range, guests, filters (price, rating, amenities) |
| FR-3 | Hotel Listing & Details | Show hotel info, photos, reviews, room types, pricing |
| FR-4 | Room Availability Check | Real-time availability for selected date range |
| FR-5 | Booking / Reservation | Reserve rooms with date range, guest details, special requests |
| FR-6 | Payment Processing | Pay via credit card, UPI, wallet; refund support |
| FR-7 | Booking Management | View, modify, cancel bookings |
| FR-8 | Reviews & Ratings | Post-stay review and rating system |
| FR-9 | Notifications | Email/SMS/Push for booking confirmation, reminders, offers |
| FR-10 | Hotel Partner Portal | Hotel owners manage inventory, pricing, and view analytics |

---

## 3. Non-Functional Requirements

| # | Requirement | Target |
|---|-------------|--------|
| NFR-1 | Availability | 99.99% uptime |
| NFR-2 | Latency | Search results < 200ms (p99) |
| NFR-3 | Scalability | Handle 100K+ concurrent users |
| NFR-4 | Consistency | Strong consistency for bookings (no double-booking) |
| NFR-5 | Data Durability | Zero data loss for payment/booking records |
| NFR-6 | Security | PCI-DSS compliant, encrypted PII, OWASP Top 10 |
| NFR-7 | Observability | Centralized logging, metrics, distributed tracing |

---

## 4. Capacity Estimation

| Metric | Value |
|--------|-------|
| DAU (Daily Active Users) | ~5 million |
| Monthly Bookings | ~10 million |
| Hotels in System | ~500,000 |
| Total Rooms | ~10 million |
| Search QPS (peak) | ~50,000 |
| Booking QPS (peak) | ~5,000 |
| Avg. Payload per Search | ~5 KB |
| Storage (5 years) | ~50 TB (bookings + media + logs) |

---

## 5. High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                           CLIENTS                                   │
│         Mobile App (iOS/Android)  |  Web App  |  Partner Portal     │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                      ┌────▼────┐
                      │   CDN   │  (Static assets, images)
                      └────┬────┘
                           │
                   ┌───────▼────────┐
                   │  API Gateway   │  Rate limiting, Auth, Routing
                   │  (Kong/NGINX)  │
                   └───────┬────────┘
                           │
        ┌──────────────────┼──────────────────────┐
        │                  │                       │
   ┌────▼────┐      ┌─────▼──────┐        ┌──────▼───────┐
   │  User   │      │  Search    │        │  Booking     │
   │ Service │      │  Service   │        │  Service     │
   └────┬────┘      └─────┬──────┘        └──────┬───────┘
        │                  │                       │
        │           ┌──────▼──────┐         ┌─────▼──────┐
        │           │ Hotel/Room  │         │  Payment   │
        │           │  Service    │         │  Service   │
        │           └──────┬──────┘         └─────┬──────┘
        │                  │                       │
   ┌────▼────┐      ┌─────▼──────┐        ┌──────▼───────┐
   │Notifica-│      │  Review    │        │  Inventory   │
   │tion Svc │      │  Service   │        │  Service     │
   └─────────┘      └────────────┘        └──────────────┘
        │                  │                       │
        ▼                  ▼                       ▼
  ┌──────────────────────────────────────────────────────┐
  │                   DATA LAYER                          │
  │  ┌──────────┐  ┌────────────┐  ┌──────────────────┐  │
  │  │PostgreSQL│  │Elasticsearch│  │   Redis Cache    │  │
  │  │ (Primary)│  │  (Search)  │  │  (Sessions/Inv.) │  │
  │  └──────────┘  └────────────┘  └──────────────────┘  │
  │  ┌──────────┐  ┌────────────┐  ┌──────────────────┐  │
  │  │  S3/Blob │  │   Kafka    │  │  MongoDB (Reviews │  │
  │  │ (Media)  │  │(Event Bus) │  │    & Analytics)   │  │
  │  └──────────┘  └────────────┘  └──────────────────┘  │
  └──────────────────────────────────────────────────────┘
```

---

## 6. Core Components

### 6.1 API Gateway
- Single entry point for all clients
- Handles authentication (JWT validation), rate limiting, request routing
- SSL termination, request/response transformation
- Tech: Kong, AWS API Gateway, or NGINX

### 6.2 User Service
- Registration, login (email/phone + OAuth)
- Profile management, password reset
- JWT token issuance and refresh
- Stores data in PostgreSQL (users table)

### 6.3 Search Service
- Full-text and geo-based hotel search
- Filters: price range, star rating, amenities, distance
- Backed by Elasticsearch for fast, ranked search results
- Caches hot queries in Redis (TTL: 5 min)

### 6.4 Hotel / Room Service
- CRUD operations for hotel and room data (partner portal)
- Room types, amenities, photos, pricing rules
- Stores structured data in PostgreSQL; images in S3/Blob

### 6.5 Inventory Service
- Real-time room availability tracking
- Date-range based inventory matrix (hotel × room_type × date → available_count)
- Uses Redis for fast reads; PostgreSQL as source of truth
- Pessimistic locking or optimistic concurrency for booking conflicts

### 6.6 Booking Service
- Orchestrates the reservation flow (check availability → hold → confirm)
- Temporary hold (TTL ~10 min) to prevent double-booking during payment
- State machine: `INITIATED → HELD → CONFIRMED → CHECKED_IN → COMPLETED → CANCELLED`
- Emits events to Kafka on state transitions

### 6.7 Payment Service
- Integrates with payment gateways (Stripe, Razorpay, PayPal)
- Handles charge, refund, and payout to hotel partners
- Idempotency keys to prevent duplicate charges
- PCI-DSS compliant; never stores raw card data

### 6.8 Notification Service
- Consumes Kafka events (booking confirmed, payment received, etc.)
- Sends Email (SES/SendGrid), SMS (Twilio), Push (FCM/APNs)
- Template-based message rendering
- Retry with exponential backoff for failures

### 6.9 Review Service
- Post-checkout review and 1–5 star rating
- Stored in MongoDB (flexible schema for text + media)
- Moderation pipeline (spam/abuse detection)
- Aggregated ratings synced to hotel search index

---

## 7. Data Flow — Booking Journey

```
User                    System
 │                        │
 │  1. Search hotels      │
 │───────────────────────►│  Search Service → Elasticsearch
 │  ◄─── Hotel list ──────│
 │                        │
 │  2. Select hotel/room  │
 │───────────────────────►│  Hotel Service → Room details + pricing
 │  ◄─── Room details ────│
 │                        │
 │  3. Check availability │
 │───────────────────────►│  Inventory Service → Redis/PostgreSQL
 │  ◄─── Available ───────│
 │                        │
 │  4. Initiate booking   │
 │───────────────────────►│  Booking Service → HOLD inventory (10 min TTL)
 │  ◄─── Booking ID ──────│
 │                        │
 │  5. Make payment       │
 │───────────────────────►│  Payment Service → Gateway charge
 │  ◄─── Payment OK ──────│
 │                        │
 │  6. Confirm booking    │
 │       (automatic)      │  Booking Service → CONFIRM, deduct inventory
 │  ◄─── Confirmation ────│  Notification Service → email/SMS
 │                        │
```

---

## 8. Database Strategy

| Database | Use Case | Why |
|----------|----------|-----|
| **PostgreSQL** | Users, Hotels, Rooms, Bookings, Payments | ACID transactions, relational integrity |
| **Elasticsearch** | Hotel search index | Full-text search, geo queries, faceted filters |
| **Redis** | Inventory cache, session store, rate limiting | Sub-ms reads, TTL-based expiry |
| **MongoDB** | Reviews, analytics, audit logs | Flexible schema, high write throughput |
| **S3 / Blob** | Hotel images, invoices, documents | Cost-effective object storage |
| **Kafka** | Event streaming between services | Decoupled async communication |

---

## 9. Handling Double-Booking (Concurrency)

Double-booking is the **most critical correctness problem**. Strategy:

1. **Pessimistic Locking** — `SELECT ... FOR UPDATE` on inventory row during booking. Simple but limits throughput.
2. **Optimistic Locking** — Version column on inventory row; retry on conflict. Better throughput.
3. **Redis Atomic Decrement** — `DECR available_count` is atomic. If result < 0, rollback and reject.
4. **Two-Phase Hold** — Booking creates a temporary hold (TTL ~10 min). Payment must complete within TTL. On expiry, hold auto-releases.

**Recommended**: Combine Redis atomic decrement (fast path) + PostgreSQL transaction (source of truth) with saga pattern for distributed consistency.

---

## 10. Scalability Strategy

| Concern | Solution |
|---------|----------|
| Read-heavy search | Elasticsearch cluster with replicas; Redis caching |
| Write-heavy bookings | Partition bookings by hotel_id or region |
| Horizontal scaling | Stateless microservices behind load balancer |
| Database scaling | Read replicas for PostgreSQL; sharding by region |
| Global users | Multi-region deployment with geo-routing (Route53) |
| Media | CDN (CloudFront) for images; lazy loading |
| Spiky traffic | Auto-scaling groups + queue-based load leveling (Kafka) |

---

## 11. Caching Strategy

```
┌──────────┐      ┌──────────┐      ┌──────────────┐
│  Client   │─────│   CDN    │─────│  App Cache    │
│  Cache    │     │ (Images) │     │  (Redis)      │
└──────────┘      └──────────┘      └──────┬───────┘
                                           │
                                    ┌──────▼───────┐
                                    │   Database   │
                                    └──────────────┘
```

| Layer | What's Cached | TTL |
|-------|---------------|-----|
| Client / Browser | Static assets, hotel thumbnails | 1 day |
| CDN | Images, JS/CSS bundles | 7 days |
| Redis L1 | Session, hot search queries | 5 min |
| Redis L2 | Inventory snapshot, pricing | 1 min |
| Elasticsearch | Denormalized hotel+room index | Near real-time sync |

---

## 12. Security Architecture

- **Authentication**: JWT (access + refresh tokens), OAuth 2.0 (Google, Facebook)
- **Authorization**: RBAC (User, Hotel Admin, System Admin)
- **API Security**: Rate limiting, IP whitelisting for partner APIs, CORS
- **Data Encryption**: TLS 1.3 in transit, AES-256 at rest
- **PCI Compliance**: Tokenized card storage via payment gateway; no raw card data
- **Input Validation**: Server-side validation, parameterized queries (SQL injection prevention)
- **Audit Trail**: Immutable audit logs for all booking and payment events

---

## 13. Observability

| Pillar | Tool | Purpose |
|--------|------|---------|
| Logging | ELK Stack (Elasticsearch, Logstash, Kibana) | Centralized log aggregation |
| Metrics | Prometheus + Grafana | Latency, throughput, error rates, SLAs |
| Tracing | Jaeger / Zipkin | Distributed request tracing across services |
| Alerting | PagerDuty / OpsGenie | On-call rotation, incident management |

---

## 14. Deployment Architecture

```
┌─────────────────────────────────────────────┐
│              Kubernetes Cluster              │
│                                             │
│  ┌─────────┐  ┌──────────┐  ┌───────────┐  │
│  │ User Svc│  │Search Svc│  │Booking Svc│  │
│  │ (3 pods)│  │ (5 pods) │  │ (5 pods)  │  │
│  └─────────┘  └──────────┘  └───────────┘  │
│  ┌─────────┐  ┌──────────┐  ┌───────────┐  │
│  │Hotel Svc│  │Payment   │  │Inventory  │  │
│  │ (3 pods)│  │Svc(3pods)│  │Svc(5 pods)│  │
│  └─────────┘  └──────────┘  └───────────┘  │
│  ┌──────────────────┐  ┌─────────────────┐  │
│  │ Notification Svc │  │  Review Svc     │  │
│  │    (2 pods)      │  │   (2 pods)      │  │
│  └──────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────┘

CI/CD: GitHub Actions → Docker Build → ECR → ArgoCD → K8s
Environments: Dev → Staging → Production (Blue-Green)
```

---

## 15. Tech Stack Summary

| Layer | Technology |
|-------|------------|
| Frontend | React.js / Next.js (Web), React Native / Flutter (Mobile) |
| API Gateway | Kong / AWS API Gateway |
| Backend | Java (Spring Boot) / Node.js / Go |
| Search | Elasticsearch |
| Primary DB | PostgreSQL (with read replicas) |
| Document DB | MongoDB |
| Cache | Redis Cluster |
| Message Queue | Apache Kafka |
| Object Storage | AWS S3 |
| CDN | CloudFront / Cloudflare |
| Container Orchestration | Kubernetes (EKS) |
| CI/CD | GitHub Actions + ArgoCD |
| Monitoring | Prometheus, Grafana, Jaeger, ELK |

---

## 16. Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Architecture style | Microservices | Independent scaling, team autonomy, fault isolation |
| Communication | Sync (REST/gRPC) + Async (Kafka) | REST for queries, Kafka for event-driven flows |
| Search engine | Elasticsearch | Battle-tested for geo + full-text search at scale |
| Booking consistency | Saga pattern with compensating transactions | Distributed transactions across services |
| Inventory locking | Redis atomic ops + DB pessimistic lock | Fast path + correctness guarantee |
| Payment integration | Gateway tokenization (Stripe) | PCI compliance without storing card data |

---

*Next: See [Hotel-Booking-App-LLD.md](./Hotel-Booking-App-LLD.md) for the Low-Level Design with class diagrams, API contracts, database schemas, and detailed component design.*
