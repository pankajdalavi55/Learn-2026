# Hotel Booking Application — High-Level Design (HLD)

---

## 1. Problem Statement

Design a scalable hotel booking platform (similar to Booking.com, OYO, MakeMyTrip) that allows users to search for hotels, view available rooms, make reservations, process payments, and manage bookings.

---

## 2. Functional Requirements


| #     | Requirement              | Description                                                            |
| ----- | ------------------------ | ---------------------------------------------------------------------- |
| FR-1  | User Registration & Auth | Sign up, login, OAuth, profile management                              |
| FR-2  | Hotel/Room Search        | Search by city, date range, guests, filters (price, rating, amenities) |
| FR-3  | Hotel Listing & Details  | Show hotel info, photos, reviews, room types, pricing                  |
| FR-4  | Room Availability Check  | Real-time availability for selected date range                         |
| FR-5  | Booking / Reservation    | Reserve rooms with date range, guest details, special requests         |
| FR-6  | Payment Processing       | Pay via credit card, UPI, wallet; refund support                       |
| FR-7  | Booking Management       | View, modify, cancel bookings                                          |
| FR-8  | Reviews & Ratings        | Post-stay review and rating system                                     |
| FR-9  | Notifications            | Email/SMS/Push for booking confirmation, reminders, offers             |
| FR-10 | Hotel Partner Portal     | Hotel owners manage inventory, pricing, and view analytics             |


---

## 3. Non-Functional Requirements


| #     | Requirement     | Target                                              |
| ----- | --------------- | --------------------------------------------------- |
| NFR-1 | Availability    | 99.99% uptime                                       |
| NFR-2 | Latency         | Search results < 200ms (p99)                        |
| NFR-3 | Scalability     | Handle 100K+ concurrent users                       |
| NFR-4 | Consistency     | Strong consistency for bookings (no double-booking) |
| NFR-5 | Data Durability | Zero data loss for payment/booking records          |
| NFR-6 | Security        | PCI-DSS compliant, encrypted PII, OWASP Top 10      |
| NFR-7 | Observability   | Centralized logging, metrics, distributed tracing   |


---

## 4. Capacity Estimation


| Metric                   | Value                            |
| ------------------------ | -------------------------------- |
| DAU (Daily Active Users) | ~5 million                       |
| Monthly Bookings         | ~10 million                      |
| Hotels in System         | ~500,000                         |
| Total Rooms              | ~10 million                      |
| Search QPS (peak)        | ~50,000                          |
| Booking QPS (peak)       | ~5,000                           |
| Avg. Payload per Search  | ~5 KB                            |
| Storage (5 years)        | ~50 TB (bookings + media + logs) |


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


| Database          | Use Case                                      | Why                                            |
| ----------------- | --------------------------------------------- | ---------------------------------------------- |
| **PostgreSQL**    | Users, Hotels, Rooms, Bookings, Payments      | ACID transactions, relational integrity        |
| **Elasticsearch** | Hotel search index                            | Full-text search, geo queries, faceted filters |
| **Redis**         | Inventory cache, session store, rate limiting | Sub-ms reads, TTL-based expiry                 |
| **MongoDB**       | Reviews, analytics, audit logs                | Flexible schema, high write throughput         |
| **S3 / Blob**     | Hotel images, invoices, documents             | Cost-effective object storage                  |
| **Kafka**         | Event streaming between services              | Decoupled async communication                  |


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


| Concern              | Solution                                                |
| -------------------- | ------------------------------------------------------- |
| Read-heavy search    | Elasticsearch cluster with replicas; Redis caching      |
| Write-heavy bookings | Partition bookings by hotel_id or region                |
| Horizontal scaling   | Stateless microservices behind load balancer            |
| Database scaling     | Read replicas for PostgreSQL; sharding by region        |
| Global users         | Multi-region deployment with geo-routing (Route53)      |
| Media                | CDN (CloudFront) for images; lazy loading               |
| Spiky traffic        | Auto-scaling groups + queue-based load leveling (Kafka) |


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


| Layer            | What's Cached                   | TTL                 |
| ---------------- | ------------------------------- | ------------------- |
| Client / Browser | Static assets, hotel thumbnails | 1 day               |
| CDN              | Images, JS/CSS bundles          | 7 days              |
| Redis L1         | Session, hot search queries     | 5 min               |
| Redis L2         | Inventory snapshot, pricing     | 1 min               |
| Elasticsearch    | Denormalized hotel+room index   | Near real-time sync |


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


| Pillar   | Tool                                        | Purpose                                     |
| -------- | ------------------------------------------- | ------------------------------------------- |
| Logging  | ELK Stack (Elasticsearch, Logstash, Kibana) | Centralized log aggregation                 |
| Metrics  | Prometheus + Grafana                        | Latency, throughput, error rates, SLAs      |
| Tracing  | Jaeger / Zipkin                             | Distributed request tracing across services |
| Alerting | PagerDuty / OpsGenie                        | On-call rotation, incident management       |


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


| Layer                   | Technology                                                |
| ----------------------- | --------------------------------------------------------- |
| Frontend                | React.js / Next.js (Web), React Native / Flutter (Mobile) |
| API Gateway             | Kong / AWS API Gateway                                    |
| Backend                 | Java (Spring Boot) / Node.js / Go                         |
| Search                  | Elasticsearch                                             |
| Primary DB              | PostgreSQL (with read replicas)                           |
| Document DB             | MongoDB                                                   |
| Cache                   | Redis Cluster                                             |
| Message Queue           | Apache Kafka                                              |
| Object Storage          | AWS S3                                                    |
| CDN                     | CloudFront / Cloudflare                                   |
| Container Orchestration | Kubernetes (EKS)                                          |
| CI/CD                   | GitHub Actions + ArgoCD                                   |
| Monitoring              | Prometheus, Grafana, Jaeger, ELK                          |


---

## 16. Key Design Decisions


| Decision            | Choice                                      | Rationale                                           |
| ------------------- | ------------------------------------------- | --------------------------------------------------- |
| Architecture style  | Microservices                               | Independent scaling, team autonomy, fault isolation |
| Communication       | Sync (REST/gRPC) + Async (Kafka)            | REST for queries, Kafka for event-driven flows      |
| Search engine       | Elasticsearch                               | Battle-tested for geo + full-text search at scale   |
| Booking consistency | Saga pattern with compensating transactions | Distributed transactions across services            |
| Inventory locking   | Redis atomic ops + DB pessimistic lock      | Fast path + correctness guarantee                   |
| Payment integration | Gateway tokenization (Stripe)               | PCI compliance without storing card data            |


---

# Hotel Booking Application — Low-Level Design (LLD)

---

## 1. Database Schema Design

### 1.1 Entity-Relationship Diagram (Conceptual)

```
┌──────────┐       ┌───────────┐       ┌───────────┐
│   User   │1────M│  Booking   │M────1│   Hotel    │
└──────────┘       └─────┬─────┘       └─────┬─────┘
                         │                    │
                         │1                   │1
                         │                    │
                   ┌─────▼─────┐       ┌─────▼─────┐
                   │  Payment  │       │   Room     │
                   └───────────┘       └─────┬─────┘
                                             │1
                                             │
                                       ┌─────▼─────┐
                                       │ Inventory  │
                                       │  (per day) │
                                       └───────────┘

┌──────────┐       ┌───────────┐       ┌───────────┐
│   User   │1────M│  Review    │M────1│   Hotel    │
└──────────┘       └───────────┘       └───────────┘
```

---

### 1.2 Table Definitions (PostgreSQL)

#### users

```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    phone           VARCHAR(20) UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    full_name       VARCHAR(150) NOT NULL,
    avatar_url      VARCHAR(500),
    role            VARCHAR(20) NOT NULL DEFAULT 'GUEST',  -- GUEST, HOTEL_ADMIN, SYSTEM_ADMIN
    oauth_provider  VARCHAR(50),
    oauth_id        VARCHAR(255),
    is_verified     BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
```

#### hotels

```sql
CREATE TABLE hotels (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id        UUID NOT NULL REFERENCES users(id),
    name            VARCHAR(300) NOT NULL,
    description     TEXT,
    star_rating     SMALLINT CHECK (star_rating BETWEEN 1 AND 5),
    address_line    VARCHAR(500) NOT NULL,
    city            VARCHAR(100) NOT NULL,
    state           VARCHAR(100),
    country         VARCHAR(100) NOT NULL,
    zip_code        VARCHAR(20),
    latitude        DECIMAL(10, 8) NOT NULL,
    longitude       DECIMAL(11, 8) NOT NULL,
    check_in_time   TIME DEFAULT '14:00',
    check_out_time  TIME DEFAULT '11:00',
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_hotels_city ON hotels(city);
CREATE INDEX idx_hotels_location ON hotels USING GIST (
    ll_to_earth(latitude, longitude)
);
CREATE INDEX idx_hotels_owner ON hotels(owner_id);
```

#### hotel_amenities

```sql
CREATE TABLE hotel_amenities (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id    UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    amenity     VARCHAR(100) NOT NULL,  -- WIFI, POOL, GYM, PARKING, SPA, RESTAURANT, etc.
    UNIQUE(hotel_id, amenity)
);
```

#### hotel_images

```sql
CREATE TABLE hotel_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id    UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    image_url   VARCHAR(500) NOT NULL,
    caption     VARCHAR(255),
    sort_order  SMALLINT DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_hotel_images_hotel ON hotel_images(hotel_id);
```

#### room_types

```sql
CREATE TABLE room_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id        UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,   -- Deluxe, Suite, Standard, etc.
    description     TEXT,
    max_guests      SMALLINT NOT NULL DEFAULT 2,
    bed_type        VARCHAR(50),             -- KING, QUEEN, TWIN, DOUBLE
    area_sqft       INT,
    base_price      DECIMAL(10, 2) NOT NULL,
    currency        VARCHAR(3) DEFAULT 'INR',
    total_rooms     INT NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_room_types_hotel ON room_types(hotel_id);
```

#### room_amenities

```sql
CREATE TABLE room_amenities (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_type_id    UUID NOT NULL REFERENCES room_types(id) ON DELETE CASCADE,
    amenity         VARCHAR(100) NOT NULL,  -- AC, TV, MINIBAR, BALCONY, etc.
    UNIQUE(room_type_id, amenity)
);
```

#### room_inventory

```sql
CREATE TABLE room_inventory (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_type_id    UUID NOT NULL REFERENCES room_types(id),
    date            DATE NOT NULL,
    total_rooms     INT NOT NULL,
    booked_rooms    INT NOT NULL DEFAULT 0,
    price_override  DECIMAL(10, 2),          -- NULL means use base_price
    version         INT NOT NULL DEFAULT 0,  -- Optimistic locking
    UNIQUE(room_type_id, date)
);

CREATE INDEX idx_inventory_room_date ON room_inventory(room_type_id, date);
```

#### bookings

```sql
CREATE TABLE bookings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_ref         VARCHAR(20) UNIQUE NOT NULL,  -- Human-readable: HB-2026-XXXX
    user_id             UUID NOT NULL REFERENCES users(id),
    hotel_id            UUID NOT NULL REFERENCES hotels(id),
    room_type_id        UUID NOT NULL REFERENCES room_types(id),
    check_in_date       DATE NOT NULL,
    check_out_date      DATE NOT NULL,
    num_rooms           SMALLINT NOT NULL DEFAULT 1,
    num_guests          SMALLINT NOT NULL,
    total_amount        DECIMAL(12, 2) NOT NULL,
    currency            VARCHAR(3) DEFAULT 'INR',
    status              VARCHAR(20) NOT NULL DEFAULT 'INITIATED',
    -- INITIATED → HELD → CONFIRMED → CHECKED_IN → COMPLETED → CANCELLED
    special_requests    TEXT,
    held_until          TIMESTAMP WITH TIME ZONE,    -- TTL for hold
    cancelled_at        TIMESTAMP WITH TIME ZONE,
    cancellation_reason TEXT,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT chk_dates CHECK (check_out_date > check_in_date)
);

CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_hotel ON bookings(hotel_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_dates ON bookings(check_in_date, check_out_date);
CREATE INDEX idx_bookings_ref ON bookings(booking_ref);
```

#### payments

```sql
CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id          UUID NOT NULL REFERENCES bookings(id),
    user_id             UUID NOT NULL REFERENCES users(id),
    amount              DECIMAL(12, 2) NOT NULL,
    currency            VARCHAR(3) DEFAULT 'INR',
    payment_method      VARCHAR(30) NOT NULL,  -- CREDIT_CARD, UPI, WALLET, NET_BANKING
    gateway             VARCHAR(30) NOT NULL,   -- STRIPE, RAZORPAY, PAYPAL
    gateway_txn_id      VARCHAR(255),
    idempotency_key     VARCHAR(255) UNIQUE NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    -- PENDING → PROCESSING → SUCCESS → FAILED → REFUNDED
    refund_amount       DECIMAL(12, 2),
    refund_reason       TEXT,
    paid_at             TIMESTAMP WITH TIME ZONE,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payments_booking ON payments(booking_id);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_idempotency ON payments(idempotency_key);
```

#### reviews

```sql
CREATE TABLE reviews (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id  UUID UNIQUE NOT NULL REFERENCES bookings(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    hotel_id    UUID NOT NULL REFERENCES hotels(id),
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title       VARCHAR(255),
    body        TEXT,
    is_approved BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_reviews_hotel ON reviews(hotel_id);
CREATE INDEX idx_reviews_user ON reviews(user_id);
```

---

## 2. Class Diagrams (Core Domain Models)

### 2.1 User Domain

```
┌──────────────────────────────────┐
│            User                  │
├──────────────────────────────────┤
│ - id: UUID                       │
│ - email: String                  │
│ - phone: String                  │
│ - passwordHash: String           │
│ - fullName: String               │
│ - role: UserRole                 │
│ - isVerified: boolean            │
├──────────────────────────────────┤
│ + register(dto): User            │
│ + login(email, pwd): AuthToken   │
│ + updateProfile(dto): User       │
│ + resetPassword(token, pwd): void│
└──────────────────────────────────┘

<<enum>> UserRole
─────────────────
GUEST
HOTEL_ADMIN
SYSTEM_ADMIN
```

### 2.2 Hotel & Room Domain

```
┌──────────────────────────────────┐
│            Hotel                 │
├──────────────────────────────────┤
│ - id: UUID                       │
│ - ownerId: UUID                  │
│ - name: String                   │
│ - starRating: int                │
│ - city: String                   │
│ - latitude: double               │
│ - longitude: double              │
│ - amenities: List<String>        │
│ - images: List<HotelImage>       │
│ - roomTypes: List<RoomType>      │
├──────────────────────────────────┤
│ + addRoomType(dto): RoomType     │
│ + updateDetails(dto): Hotel      │
│ + toggleActive(flag): void       │
└───────────────┬──────────────────┘
                │ 1
                │
                │ *
┌───────────────▼──────────────────┐
│           RoomType               │
├──────────────────────────────────┤
│ - id: UUID                       │
│ - hotelId: UUID                  │
│ - name: String                   │
│ - maxGuests: int                 │
│ - bedType: BedType               │
│ - basePrice: BigDecimal          │
│ - totalRooms: int                │
│ - amenities: List<String>        │
├──────────────────────────────────┤
│ + updatePricing(price): void     │
│ + getAvailability(dates): int    │
└───────────────┬──────────────────┘
                │ 1
                │
                │ *
┌───────────────▼──────────────────┐
│        RoomInventory             │
├──────────────────────────────────┤
│ - id: UUID                       │
│ - roomTypeId: UUID               │
│ - date: LocalDate                │
│ - totalRooms: int                │
│ - bookedRooms: int               │
│ - priceOverride: BigDecimal      │
│ - version: int                   │
├──────────────────────────────────┤
│ + getAvailable(): int            │
│ + holdRooms(count): boolean      │
│ + confirmRooms(count): void      │
│ + releaseRooms(count): void      │
└──────────────────────────────────┘
```

### 2.3 Booking Domain

```
┌──────────────────────────────────────────┐
│              Booking                      │
├──────────────────────────────────────────┤
│ - id: UUID                                │
│ - bookingRef: String                      │
│ - userId: UUID                            │
│ - hotelId: UUID                           │
│ - roomTypeId: UUID                        │
│ - checkInDate: LocalDate                  │
│ - checkOutDate: LocalDate                 │
│ - numRooms: int                           │
│ - numGuests: int                          │
│ - totalAmount: BigDecimal                 │
│ - status: BookingStatus                   │
│ - heldUntil: Instant                      │
├──────────────────────────────────────────┤
│ + initiate(): Booking                     │
│ + hold(ttlMinutes): void                  │
│ + confirm(): void                         │
│ + cancel(reason): void                    │
│ + isHoldExpired(): boolean                │
│ + calculateTotal(nights, price, rooms): $ │
└──────────────────────────────────────────┘

<<enum>> BookingStatus
──────────────────────
INITIATED
HELD
CONFIRMED
CHECKED_IN
COMPLETED
CANCELLED

State Machine:
  INITIATED ──► HELD ──► CONFIRMED ──► CHECKED_IN ──► COMPLETED
       │           │          │
       └───────────┴──────────┴──────► CANCELLED
```

### 2.4 Payment Domain

```
┌──────────────────────────────────────┐
│             Payment                   │
├──────────────────────────────────────┤
│ - id: UUID                            │
│ - bookingId: UUID                     │
│ - userId: UUID                        │
│ - amount: BigDecimal                  │
│ - paymentMethod: PaymentMethod        │
│ - gateway: String                     │
│ - gatewayTxnId: String               │
│ - idempotencyKey: String              │
│ - status: PaymentStatus              │
├──────────────────────────────────────┤
│ + initiatePayment(): Payment          │
│ + processCallback(gatewayResp): void  │
│ + refund(amount, reason): void        │
│ + isRefundable(): boolean             │
└──────────────────────────────────────┘

<<enum>> PaymentStatus
──────────────────────
PENDING
PROCESSING
SUCCESS
FAILED
REFUNDED

<<enum>> PaymentMethod
──────────────────────
CREDIT_CARD
DEBIT_CARD
UPI
WALLET
NET_BANKING
```

---

## 3. API Design (RESTful)

### 3.1 Authentication APIs


| Method | Endpoint                       | Description          | Auth          |
| ------ | ------------------------------ | -------------------- | ------------- |
| POST   | `/api/v1/auth/register`        | Register new user    | Public        |
| POST   | `/api/v1/auth/login`           | Login, get JWT       | Public        |
| POST   | `/api/v1/auth/refresh`         | Refresh access token | Refresh Token |
| POST   | `/api/v1/auth/forgot-password` | Send reset link      | Public        |
| POST   | `/api/v1/auth/reset-password`  | Reset with token     | Public        |


#### POST `/api/v1/auth/register`

**Request:**

```json
{
    "email": "john@example.com",
    "phone": "+919876543210",
    "password": "SecureP@ss123",
    "fullName": "John Doe"
}
```

**Response (201):**

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "fullName": "John Doe",
    "role": "GUEST",
    "createdAt": "2026-03-25T10:30:00Z"
}
```

#### POST `/api/v1/auth/login`

**Request:**

```json
{
    "email": "john@example.com",
    "password": "SecureP@ss123"
}
```

**Response (200):**

```json
{
    "accessToken": "eyJhbGciOiJSUzI1NiIs...",
    "refreshToken": "dGhpcyBpcyBhIHJlZnJl...",
    "expiresIn": 3600,
    "tokenType": "Bearer"
}
```

---

### 3.2 Hotel Search APIs


| Method | Endpoint                         | Description          | Auth   |
| ------ | -------------------------------- | -------------------- | ------ |
| GET    | `/api/v1/hotels/search`          | Search hotels        | Public |
| GET    | `/api/v1/hotels/{hotelId}`       | Hotel details        | Public |
| GET    | `/api/v1/hotels/{hotelId}/rooms` | Room types for hotel | Public |


#### GET `/api/v1/hotels/search`

**Query Parameters:**

```
city=Mumbai
checkIn=2026-04-10
checkOut=2026-04-12
guests=2
minPrice=2000
maxPrice=10000
starRating=4,5
amenities=WIFI,POOL
sortBy=price_asc
page=1
size=20
```

**Response (200):**

```json
{
    "hotels": [
        {
            "id": "hotel-uuid-1",
            "name": "Grand Hyatt Mumbai",
            "starRating": 5,
            "city": "Mumbai",
            "thumbnailUrl": "https://cdn.example.com/hotels/gh-mumbai.jpg",
            "avgRating": 4.6,
            "reviewCount": 2340,
            "startingPrice": 8500.00,
            "currency": "INR",
            "amenities": ["WIFI", "POOL", "GYM", "SPA"],
            "distanceKm": 3.2
        }
    ],
    "pagination": {
        "page": 1,
        "size": 20,
        "totalElements": 145,
        "totalPages": 8
    }
}
```

#### GET `/api/v1/hotels/{hotelId}`

**Response (200):**

```json
{
    "id": "hotel-uuid-1",
    "name": "Grand Hyatt Mumbai",
    "description": "Luxury 5-star hotel in the heart of Mumbai...",
    "starRating": 5,
    "address": {
        "line": "Off Western Express Highway, Santacruz East",
        "city": "Mumbai",
        "state": "Maharashtra",
        "country": "India",
        "zipCode": "400055"
    },
    "location": {
        "latitude": 19.0956,
        "longitude": 72.8617
    },
    "checkInTime": "14:00",
    "checkOutTime": "11:00",
    "amenities": ["WIFI", "POOL", "GYM", "SPA", "RESTAURANT", "PARKING"],
    "images": [
        {
            "url": "https://cdn.example.com/hotels/gh-1.jpg",
            "caption": "Lobby"
        }
    ],
    "avgRating": 4.6,
    "reviewCount": 2340
}
```

---

### 3.3 Availability & Pricing APIs


| Method | Endpoint                                | Description             | Auth   |
| ------ | --------------------------------------- | ----------------------- | ------ |
| GET    | `/api/v1/hotels/{hotelId}/availability` | Check room availability | Public |


#### GET `/api/v1/hotels/{hotelId}/availability?checkIn=2026-04-10&checkOut=2026-04-12&guests=2`

**Response (200):**

```json
{
    "hotelId": "hotel-uuid-1",
    "checkIn": "2026-04-10",
    "checkOut": "2026-04-12",
    "nights": 2,
    "rooms": [
        {
            "roomTypeId": "rt-uuid-1",
            "name": "Deluxe Room",
            "maxGuests": 2,
            "bedType": "KING",
            "areaSqft": 450,
            "amenities": ["AC", "TV", "MINIBAR", "WIFI"],
            "available": 5,
            "pricePerNight": 8500.00,
            "totalPrice": 17000.00,
            "currency": "INR"
        },
        {
            "roomTypeId": "rt-uuid-2",
            "name": "Premium Suite",
            "maxGuests": 3,
            "bedType": "KING",
            "areaSqft": 800,
            "amenities": ["AC", "TV", "MINIBAR", "WIFI", "BALCONY", "JACUZZI"],
            "available": 2,
            "pricePerNight": 15000.00,
            "totalPrice": 30000.00,
            "currency": "INR"
        }
    ]
}
```

---

### 3.4 Booking APIs


| Method | Endpoint                              | Description           | Auth |
| ------ | ------------------------------------- | --------------------- | ---- |
| POST   | `/api/v1/bookings`                    | Create booking (hold) | JWT  |
| GET    | `/api/v1/bookings/{bookingId}`        | Get booking details   | JWT  |
| GET    | `/api/v1/bookings/me`                 | List user's bookings  | JWT  |
| PUT    | `/api/v1/bookings/{bookingId}/cancel` | Cancel booking        | JWT  |


#### POST `/api/v1/bookings`

**Request:**

```json
{
    "hotelId": "hotel-uuid-1",
    "roomTypeId": "rt-uuid-1",
    "checkInDate": "2026-04-10",
    "checkOutDate": "2026-04-12",
    "numRooms": 1,
    "numGuests": 2,
    "specialRequests": "Late check-in around 10 PM"
}
```

**Response (201):**

```json
{
    "id": "booking-uuid-1",
    "bookingRef": "HB-2026-7A3F",
    "status": "HELD",
    "hotelName": "Grand Hyatt Mumbai",
    "roomType": "Deluxe Room",
    "checkInDate": "2026-04-10",
    "checkOutDate": "2026-04-12",
    "nights": 2,
    "numRooms": 1,
    "totalAmount": 17000.00,
    "currency": "INR",
    "heldUntil": "2026-03-25T10:40:00Z",
    "paymentDeadlineMinutes": 10,
    "createdAt": "2026-03-25T10:30:00Z"
}
```

---

### 3.5 Payment APIs


| Method | Endpoint                              | Description        | Auth           |
| ------ | ------------------------------------- | ------------------ | -------------- |
| POST   | `/api/v1/payments`                    | Initiate payment   | JWT            |
| GET    | `/api/v1/payments/{paymentId}`        | Get payment status | JWT            |
| POST   | `/api/v1/payments/webhook`            | Gateway callback   | Gateway Secret |
| POST   | `/api/v1/payments/{paymentId}/refund` | Request refund     | JWT            |


#### POST `/api/v1/payments`

**Request:**

```json
{
    "bookingId": "booking-uuid-1",
    "paymentMethod": "UPI",
    "idempotencyKey": "idem-key-unique-12345"
}
```

**Response (200):**

```json
{
    "id": "payment-uuid-1",
    "bookingId": "booking-uuid-1",
    "amount": 17000.00,
    "currency": "INR",
    "status": "PROCESSING",
    "gatewayRedirectUrl": "https://razorpay.com/pay/xyz123",
    "createdAt": "2026-03-25T10:31:00Z"
}
```

#### POST `/api/v1/payments/webhook` (Gateway Callback)

**Request (from Razorpay):**

```json
{
    "event": "payment.captured",
    "payload": {
        "payment": {
            "entity": {
                "id": "pay_xyz123",
                "amount": 1700000,
                "currency": "INR",
                "status": "captured",
                "notes": {
                    "bookingId": "booking-uuid-1",
                    "paymentId": "payment-uuid-1"
                }
            }
        }
    }
}
```

---

### 3.6 Review APIs


| Method | Endpoint                           | Description            | Auth   |
| ------ | ---------------------------------- | ---------------------- | ------ |
| POST   | `/api/v1/reviews`                  | Submit review          | JWT    |
| GET    | `/api/v1/hotels/{hotelId}/reviews` | List reviews for hotel | Public |


#### POST `/api/v1/reviews`

**Request:**

```json
{
    "bookingId": "booking-uuid-1",
    "hotelId": "hotel-uuid-1",
    "rating": 5,
    "title": "Exceptional stay!",
    "body": "The room was spotless, staff very courteous..."
}
```

---

### 3.7 Hotel Partner APIs


| Method | Endpoint                                       | Description           | Auth        |
| ------ | ---------------------------------------------- | --------------------- | ----------- |
| POST   | `/api/v1/partner/hotels`                       | Register hotel        | HOTEL_ADMIN |
| PUT    | `/api/v1/partner/hotels/{hotelId}`             | Update hotel          | HOTEL_ADMIN |
| POST   | `/api/v1/partner/hotels/{hotelId}/rooms`       | Add room type         | HOTEL_ADMIN |
| PUT    | `/api/v1/partner/rooms/{roomTypeId}/pricing`   | Update pricing        | HOTEL_ADMIN |
| PUT    | `/api/v1/partner/rooms/{roomTypeId}/inventory` | Bulk update inventory | HOTEL_ADMIN |
| GET    | `/api/v1/partner/hotels/{hotelId}/analytics`   | Booking analytics     | HOTEL_ADMIN |


---

## 4. Service Layer Design

### 4.1 BookingService — Core Orchestration Logic

```java
@Service
@Transactional
public class BookingService {

    private final BookingRepository bookingRepo;
    private final InventoryService inventoryService;
    private final PaymentService paymentService;
    private final NotificationService notificationService;
    private final EventPublisher eventPublisher;

    public BookingResponse createBooking(BookingRequest request, UUID userId) {
        // 1. Validate dates and room type
        RoomType roomType = validateRoomType(request.getRoomTypeId());

        // 2. Check and hold inventory atomically
        boolean held = inventoryService.holdRooms(
            request.getRoomTypeId(),
            request.getCheckInDate(),
            request.getCheckOutDate(),
            request.getNumRooms()
        );
        if (!held) {
            throw new RoomNotAvailableException("Rooms no longer available");
        }

        // 3. Calculate total
        BigDecimal total = calculateTotal(
            roomType, request.getCheckInDate(),
            request.getCheckOutDate(), request.getNumRooms()
        );

        // 4. Create booking in HELD status
        Booking booking = Booking.builder()
            .bookingRef(generateBookingRef())
            .userId(userId)
            .hotelId(roomType.getHotelId())
            .roomTypeId(request.getRoomTypeId())
            .checkInDate(request.getCheckInDate())
            .checkOutDate(request.getCheckOutDate())
            .numRooms(request.getNumRooms())
            .numGuests(request.getNumGuests())
            .totalAmount(total)
            .status(BookingStatus.HELD)
            .heldUntil(Instant.now().plusSeconds(600))  // 10 min TTL
            .build();

        bookingRepo.save(booking);

        // 5. Publish event
        eventPublisher.publish(new BookingHeldEvent(booking));

        return BookingResponse.from(booking);
    }

    public void confirmBooking(UUID bookingId) {
        Booking booking = bookingRepo.findById(bookingId)
            .orElseThrow(() -> new BookingNotFoundException(bookingId));

        if (booking.isHoldExpired()) {
            cancelBooking(bookingId, "Hold expired");
            throw new BookingExpiredException("Booking hold has expired");
        }

        booking.setStatus(BookingStatus.CONFIRMED);
        bookingRepo.save(booking);

        inventoryService.confirmRooms(
            booking.getRoomTypeId(),
            booking.getCheckInDate(),
            booking.getCheckOutDate(),
            booking.getNumRooms()
        );

        eventPublisher.publish(new BookingConfirmedEvent(booking));
    }

    public void cancelBooking(UUID bookingId, String reason) {
        Booking booking = bookingRepo.findById(bookingId)
            .orElseThrow(() -> new BookingNotFoundException(bookingId));

        booking.setStatus(BookingStatus.CANCELLED);
        booking.setCancelledAt(Instant.now());
        booking.setCancellationReason(reason);
        bookingRepo.save(booking);

        // Release held inventory
        inventoryService.releaseRooms(
            booking.getRoomTypeId(),
            booking.getCheckInDate(),
            booking.getCheckOutDate(),
            booking.getNumRooms()
        );

        // Trigger refund if payment was made
        paymentService.refundIfApplicable(bookingId);

        eventPublisher.publish(new BookingCancelledEvent(booking));
    }
}
```

### 4.2 InventoryService — Concurrency-Safe Room Management

```java
@Service
public class InventoryService {

    private final RoomInventoryRepository inventoryRepo;
    private final RedisTemplate<String, Integer> redisTemplate;

    public boolean holdRooms(UUID roomTypeId, LocalDate checkIn,
                             LocalDate checkOut, int numRooms) {
        List<LocalDate> dates = checkIn.datesUntil(checkOut).toList();

        // Fast path: atomic Redis check
        for (LocalDate date : dates) {
            String key = inventoryKey(roomTypeId, date);
            Long remaining = redisTemplate.opsForValue().decrement(key, numRooms);
            if (remaining != null && remaining < 0) {
                // Rollback Redis decrements
                rollbackRedis(roomTypeId, checkIn, date, numRooms);
                return false;
            }
        }

        // Slow path: persist to DB with optimistic lock
        try {
            for (LocalDate date : dates) {
                RoomInventory inv = inventoryRepo
                    .findByRoomTypeIdAndDate(roomTypeId, date)
                    .orElseThrow();
                if (inv.getAvailable() < numRooms) {
                    throw new InsufficientInventoryException();
                }
                inv.setBookedRooms(inv.getBookedRooms() + numRooms);
                inventoryRepo.save(inv); // version check on save
            }
            return true;
        } catch (OptimisticLockingFailureException | InsufficientInventoryException e) {
            // Rollback Redis
            rollbackRedis(roomTypeId, checkIn, checkOut, numRooms);
            return false;
        }
    }

    public void releaseRooms(UUID roomTypeId, LocalDate checkIn,
                             LocalDate checkOut, int numRooms) {
        List<LocalDate> dates = checkIn.datesUntil(checkOut).toList();
        for (LocalDate date : dates) {
            // Redis increment
            String key = inventoryKey(roomTypeId, date);
            redisTemplate.opsForValue().increment(key, numRooms);

            // DB update
            RoomInventory inv = inventoryRepo
                .findByRoomTypeIdAndDate(roomTypeId, date).orElseThrow();
            inv.setBookedRooms(inv.getBookedRooms() - numRooms);
            inventoryRepo.save(inv);
        }
    }

    private String inventoryKey(UUID roomTypeId, LocalDate date) {
        return "inv:" + roomTypeId + ":" + date;
    }
}
```

### 4.3 PaymentService — Idempotent Payment Processing

```java
@Service
public class PaymentService {

    private final PaymentRepository paymentRepo;
    private final PaymentGatewayClient gatewayClient;
    private final BookingService bookingService;
    private final EventPublisher eventPublisher;

    public PaymentResponse initiatePayment(PaymentRequest request, UUID userId) {
        // Idempotency check
        Optional<Payment> existing = paymentRepo
            .findByIdempotencyKey(request.getIdempotencyKey());
        if (existing.isPresent()) {
            return PaymentResponse.from(existing.get());
        }

        Booking booking = bookingService.getBooking(request.getBookingId());
        if (booking.isHoldExpired()) {
            throw new BookingExpiredException("Payment window expired");
        }

        Payment payment = Payment.builder()
            .bookingId(request.getBookingId())
            .userId(userId)
            .amount(booking.getTotalAmount())
            .currency(booking.getCurrency())
            .paymentMethod(request.getPaymentMethod())
            .idempotencyKey(request.getIdempotencyKey())
            .status(PaymentStatus.PROCESSING)
            .build();

        paymentRepo.save(payment);

        // Call external gateway
        GatewayResponse gatewayResp = gatewayClient.createOrder(
            payment.getAmount(),
            payment.getCurrency(),
            payment.getId().toString()
        );

        payment.setGateway(gatewayResp.getGateway());
        payment.setGatewayTxnId(gatewayResp.getTransactionId());
        paymentRepo.save(payment);

        return PaymentResponse.from(payment, gatewayResp.getRedirectUrl());
    }

    @Transactional
    public void handleGatewayWebhook(WebhookPayload payload) {
        Payment payment = paymentRepo.findByGatewayTxnId(payload.getTxnId())
            .orElseThrow();

        if ("captured".equals(payload.getStatus())) {
            payment.setStatus(PaymentStatus.SUCCESS);
            payment.setPaidAt(Instant.now());
            paymentRepo.save(payment);

            bookingService.confirmBooking(payment.getBookingId());
            eventPublisher.publish(new PaymentSuccessEvent(payment));
        } else {
            payment.setStatus(PaymentStatus.FAILED);
            paymentRepo.save(payment);

            bookingService.cancelBooking(
                payment.getBookingId(), "Payment failed"
            );
            eventPublisher.publish(new PaymentFailedEvent(payment));
        }
    }
}
```

---

## 5. Search Service — Elasticsearch Integration

### 5.1 Hotel Search Index Mapping

```json
{
    "mappings": {
        "properties": {
            "hotelId":      { "type": "keyword" },
            "name":         { "type": "text", "analyzer": "standard" },
            "city":         { "type": "keyword" },
            "country":      { "type": "keyword" },
            "starRating":   { "type": "integer" },
            "avgRating":    { "type": "float" },
            "reviewCount":  { "type": "integer" },
            "startingPrice":{ "type": "float" },
            "location":     { "type": "geo_point" },
            "amenities":    { "type": "keyword" },
            "isActive":     { "type": "boolean" },
            "updatedAt":    { "type": "date" }
        }
    }
}
```

### 5.2 Search Query Builder

```java
@Service
public class HotelSearchService {

    private final ElasticsearchClient esClient;
    private final RedisTemplate<String, String> redisTemplate;

    public SearchResponse searchHotels(SearchRequest request) {
        // Check cache first
        String cacheKey = buildCacheKey(request);
        String cached = redisTemplate.opsForValue().get(cacheKey);
        if (cached != null) {
            return deserialize(cached);
        }

        BoolQuery.Builder boolQuery = new BoolQuery.Builder();

        // City filter
        boolQuery.filter(q -> q.term(t ->
            t.field("city").value(request.getCity())
        ));

        // Star rating filter
        if (request.getStarRatings() != null) {
            boolQuery.filter(q -> q.terms(t ->
                t.field("starRating")
                 .terms(tv -> tv.value(
                     request.getStarRatings().stream()
                         .map(FieldValue::of).toList()
                 ))
            ));
        }

        // Price range filter
        if (request.getMinPrice() != null || request.getMaxPrice() != null) {
            boolQuery.filter(q -> q.range(r -> {
                var range = r.field("startingPrice");
                if (request.getMinPrice() != null) range.gte(JsonData.of(request.getMinPrice()));
                if (request.getMaxPrice() != null) range.lte(JsonData.of(request.getMaxPrice()));
                return range;
            }));
        }

        // Amenities filter
        if (request.getAmenities() != null) {
            for (String amenity : request.getAmenities()) {
                boolQuery.filter(q -> q.term(t ->
                    t.field("amenities").value(amenity)
                ));
            }
        }

        // Active only
        boolQuery.filter(q -> q.term(t ->
            t.field("isActive").value(true)
        ));

        // Execute search
        var searchResp = esClient.search(s -> s
            .index("hotels")
            .query(q -> q.bool(boolQuery.build()))
            .sort(buildSort(request.getSortBy()))
            .from((request.getPage() - 1) * request.getSize())
            .size(request.getSize()),
            HotelDocument.class
        );

        SearchResponse response = mapResponse(searchResp);

        // Cache result
        redisTemplate.opsForValue().set(
            cacheKey, serialize(response),
            Duration.ofMinutes(5)
        );

        return response;
    }
}
```

---

## 6. Event-Driven Architecture (Kafka Topics)

### 6.1 Topic Design


| Topic              | Producer          | Consumers               | Purpose                        |
| ------------------ | ----------------- | ----------------------- | ------------------------------ |
| `booking.events`   | Booking Service   | Notification, Analytics | Booking state changes          |
| `payment.events`   | Payment Service   | Booking, Notification   | Payment state changes          |
| `inventory.events` | Inventory Service | Search (ES sync)        | Availability changes           |
| `review.events`    | Review Service    | Search (rating sync)    | New reviews                    |
| `user.events`      | User Service      | Notification            | Registrations, profile updates |


### 6.2 Event Schemas

```json
// BookingConfirmedEvent
{
    "eventId": "evt-uuid",
    "eventType": "BOOKING_CONFIRMED",
    "timestamp": "2026-03-25T10:35:00Z",
    "payload": {
        "bookingId": "booking-uuid-1",
        "bookingRef": "HB-2026-7A3F",
        "userId": "user-uuid-1",
        "hotelId": "hotel-uuid-1",
        "hotelName": "Grand Hyatt Mumbai",
        "roomType": "Deluxe Room",
        "checkInDate": "2026-04-10",
        "checkOutDate": "2026-04-12",
        "totalAmount": 17000.00
    }
}
```

### 6.3 Notification Event Consumer

```java
@Component
public class BookingEventConsumer {

    private final NotificationService notificationService;

    @KafkaListener(topics = "booking.events", groupId = "notification-service")
    public void handleBookingEvent(BookingEvent event) {
        switch (event.getEventType()) {
            case "BOOKING_CONFIRMED" -> {
                notificationService.sendEmail(
                    event.getUserEmail(),
                    "booking-confirmation",
                    Map.of(
                        "bookingRef", event.getBookingRef(),
                        "hotelName", event.getHotelName(),
                        "checkIn", event.getCheckInDate(),
                        "checkOut", event.getCheckOutDate(),
                        "total", event.getTotalAmount()
                    )
                );
                notificationService.sendSms(event.getUserPhone(),
                    "Booking confirmed! Ref: " + event.getBookingRef());
            }
            case "BOOKING_CANCELLED" ->
                notificationService.sendEmail(
                    event.getUserEmail(),
                    "booking-cancellation",
                    Map.of("bookingRef", event.getBookingRef())
                );
        }
    }
}
```

---

## 7. Scheduled Jobs


| Job                        | Schedule         | Description                                                            |
| -------------------------- | ---------------- | ---------------------------------------------------------------------- |
| `HoldExpiryJob`            | Every 1 minute   | Finds HELD bookings past `heldUntil`, cancels them, releases inventory |
| `InventorySyncJob`         | Every 5 minutes  | Syncs Redis inventory cache from PostgreSQL                            |
| `SearchIndexSyncJob`       | Every 10 minutes | Full re-index of hotel data into Elasticsearch                         |
| `PaymentReconciliationJob` | Daily 2 AM       | Reconciles payment records with gateway reports                        |
| `ReviewAggregationJob`     | Every 30 minutes | Recalculates average rating per hotel                                  |


### HoldExpiryJob Implementation

```java
@Component
public class HoldExpiryJob {

    private final BookingRepository bookingRepo;
    private final BookingService bookingService;

    @Scheduled(fixedRate = 60_000) // every 1 min
    public void expireHeldBookings() {
        List<Booking> expired = bookingRepo
            .findByStatusAndHeldUntilBefore(
                BookingStatus.HELD, Instant.now()
            );

        for (Booking booking : expired) {
            try {
                bookingService.cancelBooking(
                    booking.getId(), "Hold expired automatically"
                );
            } catch (Exception e) {
                log.error("Failed to expire booking {}", booking.getId(), e);
            }
        }
    }
}
```

---

## 8. Error Handling & Response Format

### 8.1 Standard Error Response

```json
{
    "error": {
        "code": "ROOM_NOT_AVAILABLE",
        "message": "The requested room type is no longer available for the selected dates.",
        "details": {
            "roomTypeId": "rt-uuid-1",
            "checkIn": "2026-04-10",
            "checkOut": "2026-04-12"
        },
        "timestamp": "2026-03-25T10:30:00Z",
        "traceId": "abc-123-def-456"
    }
}
```

### 8.2 Error Code Catalog


| HTTP Status | Error Code                | Description                        |
| ----------- | ------------------------- | ---------------------------------- |
| 400         | INVALID_REQUEST           | Missing or invalid request fields  |
| 400         | INVALID_DATE_RANGE        | check-out must be after check-in   |
| 401         | UNAUTHORIZED              | Missing or invalid JWT             |
| 403         | FORBIDDEN                 | Insufficient role permissions      |
| 404         | HOTEL_NOT_FOUND           | Hotel ID does not exist            |
| 404         | BOOKING_NOT_FOUND         | Booking ID does not exist          |
| 409         | ROOM_NOT_AVAILABLE        | Rooms sold out for selected dates  |
| 409         | BOOKING_ALREADY_CANCELLED | Booking is already cancelled       |
| 410         | BOOKING_HOLD_EXPIRED      | Payment window has passed          |
| 422         | DUPLICATE_REVIEW          | User already reviewed this booking |
| 429         | RATE_LIMIT_EXCEEDED       | Too many requests                  |
| 502         | PAYMENT_GATEWAY_ERROR     | External payment gateway failure   |


---

## 9. Sequence Diagrams

### 9.1 Complete Booking Flow

```
Client          API GW         Booking Svc     Inventory Svc    Payment Svc     Notification
  │                │                │                │                │                │
  │ POST /bookings │                │                │                │                │
  │───────────────►│                │                │                │                │
  │                │ route + auth   │                │                │                │
  │                │───────────────►│                │                │                │
  │                │                │ holdRooms()    │                │                │
  │                │                │───────────────►│                │                │
  │                │                │    ◄─── OK ────│                │                │
  │                │                │                │                │                │
  │                │  ◄── 201 ──────│ (status=HELD)  │                │                │
  │  ◄── 201 ──────│                │                │                │                │
  │                │                │                │                │                │
  │ POST /payments │                │                │                │                │
  │───────────────►│                │                │                │                │
  │                │───────────────────────────────────────────────►│                │
  │                │                │                │   createOrder() │                │
  │  ◄── redirect ─│                │                │ ◄── gwUrl ─────│                │
  │                │                │                │                │                │
  │ (user pays on  │                │                │                │                │
  │  gateway page) │                │                │                │                │
  │                │                │                │                │                │
  │                │   webhook      │                │                │                │
  │                │───────────────────────────────────────────────►│                │
  │                │                │                │  handleWebhook()│                │
  │                │                │ confirmBooking()│               │                │
  │                │                │◄───────────────────────────────│                │
  │                │                │ confirmRooms() │                │                │
  │                │                │───────────────►│                │                │
  │                │                │                │                │                │
  │                │                │──── Kafka: BookingConfirmedEvent ──────────────►│
  │                │                │                │                │  sendEmail()   │
  │                │                │                │                │  sendSMS()     │
  │                │                │                │                │                │
```

### 9.2 Cancellation Flow

```
Client          Booking Svc     Inventory Svc    Payment Svc     Notification
  │                │                │                │                │
  │ PUT /cancel    │                │                │                │
  │───────────────►│                │                │                │
  │                │ releaseRooms() │                │                │
  │                │───────────────►│                │                │
  │                │   ◄─── OK ────│                │                │
  │                │                │                │                │
  │                │ refundIfApplicable()            │                │
  │                │───────────────────────────────►│                │
  │                │                │  ◄── refund ──│                │
  │                │                │                │                │
  │                │───── Kafka: BookingCancelledEvent ─────────────►│
  │  ◄── 200 ─────│                │                │  sendEmail()   │
  │                │                │                │                │
```

---

## 10. Caching Patterns Used

### 10.1 Cache-Aside (Hotel Details)

```java
public Hotel getHotel(UUID hotelId) {
    String key = "hotel:" + hotelId;
    Hotel cached = redisTemplate.opsForValue().get(key);
    if (cached != null) return cached;

    Hotel hotel = hotelRepo.findById(hotelId).orElseThrow();
    redisTemplate.opsForValue().set(key, hotel, Duration.ofMinutes(30));
    return hotel;
}
```

### 10.2 Write-Through (Inventory)

```java
public void updateInventory(UUID roomTypeId, LocalDate date, int booked) {
    // Write to DB first
    RoomInventory inv = inventoryRepo
        .findByRoomTypeIdAndDate(roomTypeId, date).orElseThrow();
    inv.setBookedRooms(booked);
    inventoryRepo.save(inv);

    // Then update cache
    String key = "inv:" + roomTypeId + ":" + date;
    redisTemplate.opsForValue().set(key, inv.getAvailable(), Duration.ofMinutes(10));
}
```

---

## 11. Design Patterns Used


| Pattern              | Where Used                       | Purpose                                           |
| -------------------- | -------------------------------- | ------------------------------------------------- |
| **Builder**          | Booking, Payment entity creation | Clean object construction with many fields        |
| **Strategy**         | Payment gateway selection        | Swap between Stripe, Razorpay, PayPal             |
| **Observer / Event** | Kafka event publishing           | Decouple booking from notification                |
| **Saga**             | Booking → Payment → Confirm flow | Distributed transaction with compensating actions |
| **Repository**       | All DB access layers             | Abstract persistence from business logic          |
| **Factory**          | Notification channel selection   | Create Email, SMS, or Push sender                 |
| **Circuit Breaker**  | External gateway calls           | Resilience against downstream failures            |
| **Idempotency Key**  | Payment processing               | Prevent duplicate charges on retry                |


---

## 12. Folder / Package Structure (Spring Boot)

```
hotel-booking-service/
├── src/main/java/com/hotelbooking/
│   ├── HotelBookingApplication.java
│   ├── config/
│   │   ├── SecurityConfig.java
│   │   ├── RedisConfig.java
│   │   ├── KafkaConfig.java
│   │   └── ElasticsearchConfig.java
│   ├── controller/
│   │   ├── AuthController.java
│   │   ├── HotelSearchController.java
│   │   ├── BookingController.java
│   │   ├── PaymentController.java
│   │   ├── ReviewController.java
│   │   └── PartnerController.java
│   ├── service/
│   │   ├── UserService.java
│   │   ├── HotelService.java
│   │   ├── HotelSearchService.java
│   │   ├── InventoryService.java
│   │   ├── BookingService.java
│   │   ├── PaymentService.java
│   │   ├── NotificationService.java
│   │   └── ReviewService.java
│   ├── repository/
│   │   ├── UserRepository.java
│   │   ├── HotelRepository.java
│   │   ├── RoomTypeRepository.java
│   │   ├── RoomInventoryRepository.java
│   │   ├── BookingRepository.java
│   │   ├── PaymentRepository.java
│   │   └── ReviewRepository.java
│   ├── model/
│   │   ├── entity/
│   │   │   ├── User.java
│   │   │   ├── Hotel.java
│   │   │   ├── RoomType.java
│   │   │   ├── RoomInventory.java
│   │   │   ├── Booking.java
│   │   │   ├── Payment.java
│   │   │   └── Review.java
│   │   ├── enums/
│   │   │   ├── UserRole.java
│   │   │   ├── BookingStatus.java
│   │   │   ├── PaymentStatus.java
│   │   │   ├── PaymentMethod.java
│   │   │   └── BedType.java
│   │   ├── dto/
│   │   │   ├── request/
│   │   │   │   ├── BookingRequest.java
│   │   │   │   ├── PaymentRequest.java
│   │   │   │   ├── SearchRequest.java
│   │   │   │   └── ReviewRequest.java
│   │   │   └── response/
│   │   │       ├── BookingResponse.java
│   │   │       ├── PaymentResponse.java
│   │   │       ├── SearchResponse.java
│   │   │       └── HotelDetailResponse.java
│   │   └── event/
│   │       ├── BookingHeldEvent.java
│   │       ├── BookingConfirmedEvent.java
│   │       ├── BookingCancelledEvent.java
│   │       ├── PaymentSuccessEvent.java
│   │       └── PaymentFailedEvent.java
│   ├── kafka/
│   │   ├── EventPublisher.java
│   │   ├── BookingEventConsumer.java
│   │   └── PaymentEventConsumer.java
│   ├── scheduler/
│   │   ├── HoldExpiryJob.java
│   │   ├── InventorySyncJob.java
│   │   └── PaymentReconciliationJob.java
│   ├── gateway/
│   │   ├── PaymentGatewayClient.java
│   │   ├── StripeGateway.java
│   │   └── RazorpayGateway.java
│   ├── exception/
│   │   ├── GlobalExceptionHandler.java
│   │   ├── BookingNotFoundException.java
│   │   ├── BookingExpiredException.java
│   │   ├── RoomNotAvailableException.java
│   │   └── PaymentFailedException.java
│   └── util/
│       ├── BookingRefGenerator.java
│       └── JwtUtil.java
├── src/main/resources/
│   ├── application.yml
│   ├── application-dev.yml
│   ├── application-prod.yml
│   └── db/migration/
│       ├── V1__create_users.sql
│       ├── V2__create_hotels.sql
│       ├── V3__create_rooms.sql
│       ├── V4__create_inventory.sql
│       ├── V5__create_bookings.sql
│       ├── V6__create_payments.sql
│       └── V7__create_reviews.sql
└── src/test/java/com/hotelbooking/
    ├── service/
    │   ├── BookingServiceTest.java
    │   ├── InventoryServiceTest.java
    │   └── PaymentServiceTest.java
    └── controller/
        ├── BookingControllerTest.java
        └── HotelSearchControllerTest.java
```

---

## 13. Key Interview Discussion Points


| Topic                                      | Talking Point                                                                                 |
| ------------------------------------------ | --------------------------------------------------------------------------------------------- |
| **Why microservices?**                     | Independent scaling (search scales differently than bookings), team autonomy, fault isolation |
| **Why not a single DB?**                   | Polyglot persistence — relational for bookings (ACID), ES for search, Redis for caching       |
| **How to avoid double-booking?**           | Redis atomic decrement + DB pessimistic/optimistic lock + hold TTL                            |
| **Why Kafka over REST for notifications?** | Async, retry-safe, decoupled; notification failure shouldn't block booking                    |
| **How to handle payment failures?**        | Saga pattern — compensating action (cancel booking, release inventory)                        |
| **How to scale search?**                   | Elasticsearch horizontal scaling + Redis caching + CDN for static data                        |
| **How to handle hotel pricing changes?**   | Event-driven: pricing update → Kafka → Elasticsearch re-index                                 |
| **CAP theorem trade-off?**                 | CP for bookings (consistency > availability); AP for search (eventual consistency OK)         |


---

*See also: [Hotel-Booking-App-HLD.md](./Hotel-Booking-App-HLD.md) for the High-Level Architecture overview.*