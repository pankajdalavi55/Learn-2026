# Complete System Design: Flash Sale System (Production-Ready)

> **Related:** [Payment Processing System](./007-payment-system.md)

> **Complexity Level:** Advanced
> **Estimated Time:** 60-75 minutes in interview
> **Real-World Examples:** Flipkart Big Billion Days, Amazon Prime Day, Xiaomi Mi Flash Sale, Taobao Singles' Day, Ticketmaster drops

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Inventory Engine (The Heart)](#7-deep-dive-inventory-engine-the-heart)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Follow-Up Questions (Rapid Fire)](#13-follow-up-questions-rapid-fire)
14. [Production Checklist](#14-production-checklist)
15. [Summary](#summary)

---

## 1. Problem Statement

**Initial Question:**
"Design a flash sale system like Flipkart's Big Billion Days, where 10 million users try to buy 10,000 units of a discounted phone the moment the sale opens at 12:00:00 PM."

**Interviewer's Perspective:**
This problem assesses a candidate's ability to reason about:
- **Extreme traffic spikes:** 0 → 1M RPS in under one second, then back to normal
- **Oversell prevention:** never sell the 10,001st unit of a 10,000-unit stock
- **Hot key contention:** millions of concurrent writes to a *single* inventory row
- **Fairness:** genuine users must not lose to bots and scalpers
- **Graceful degradation:** the rest of the marketplace must stay alive while the sale burns
- **Async decoupling:** the checkout path must not block on payment or order persistence

> **Why this is a hard problem:** Normal e-commerce load is spread over time and across millions of SKUs. A flash sale compresses everything into one second, on one SKU, on one database row. The interesting failure is not "the system is slow" — it is "the system oversold 4,000 phones and we must cancel real customer orders."

**Two distinct sub-problems (say this out loud in the interview):**

| Sub-problem | Nature | Solution family |
|---|---|---|
| **Read storm** — 10M users refreshing the product page | 99.9% of traffic, read-only | CDN, static rendering, caching, sold-out short-circuit |
| **Write storm** — 10M users hitting "Buy Now" | 0.1% of traffic, but must be exactly correct | Atomic Redis decrement, queue-based load leveling, async order creation |

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before designing, I want to pin down the scope. Flash sale systems live or die on a few specific decisions — oversell tolerance, fairness policy, and whether the checkout is synchronous — so let me clarify those first."

**Interviewer:** "Go ahead."

### 2.1 Functional Requirements

**Candidate:** "Functional requirements I'd like to confirm:
1. Users browse a sale landing page with a countdown timer before the sale starts.
2. At sale start, users click 'Buy Now' and either win a unit or are told it's sold out.
3. Winners get a reserved unit held for a fixed window (say 10 minutes) to complete payment.
4. If payment is not completed in time, the unit goes back into stock.
5. Per-user purchase limits (e.g., 1 unit per user per sale).
6. Real-time stock display ('Only 240 left!').
7. Order history and status tracking after purchase."

**Interviewer:** "Yes to all. Add one thing — we run multiple flash sales concurrently, up to 500 SKUs during Big Billion Days, and some of them are 'deal of the hour' with staggered start times."

**Candidate:** "Understood. Core functional requirements:
1. ✅ Sale landing page with countdown and live stock indicator
2. ✅ Purchase attempt at sale start — win or sold-out response
3. ✅ Reservation with TTL (10-minute payment window)
4. ✅ Automatic inventory release on reservation expiry
5. ✅ Per-user purchase limits enforced atomically
6. ✅ 500 concurrent sales with staggered start times
7. ✅ Order tracking and status"

**Out of scope (state this explicitly):** payment processing internals (delegated to the payment system), logistics/delivery, product catalog management, pricing engine.

### 2.2 Non-Functional Requirements

**Candidate:** "Non-functional requirements — these are the ones that shape the architecture:
1. What is the oversell tolerance? Zero, or is a small overshoot acceptable?
2. Is undersell acceptable — can we end the sale with 3 unsold units?
3. What's the latency budget for the buy request?
4. What availability do we need, and is degradation acceptable?
5. Do we need strict fairness (first-come-first-served) or is probabilistic fairness OK?"

**Interviewer:** "Zero oversell — that's a hard business constraint, since we'd have to cancel real orders and it becomes a PR problem. Undersell of a handful of units is acceptable. Buy response under 500 ms at p99. The rest of the site must never go down because of a flash sale. Fairness should be 'reasonable' — we don't need microsecond ordering, but bots must not sweep the stock."

**Candidate:** "Summarizing:

| Requirement | Target |
|---|---|
| Consistency | Zero oversell (hard). Undersell of < 0.1% acceptable |
| Latency | p99 < 500 ms for buy request; < 100 ms for page load (CDN) |
| Availability | 99.99% for the sale path; **100% isolation** — flash sale must never take down core commerce |
| Throughput | 1M RPS peak at sale open, sustained for ~10 seconds |
| Fairness | Bot mitigation + per-user limits; approximate FCFS is acceptable |
| Degradation | Sale can fail closed ('try again') — core browsing/checkout must stay up |

> **Key trade-off to state early:** We deliberately choose **CP over AP for inventory** (correctness beats availability for stock), but **AP for everything else** (page rendering, stock display, order status). Displaying a slightly stale stock count is fine; selling a unit twice is not."

### 2.3 Scale

**Candidate:** "What scale should I plan for?"

**Interviewer:** "10 million users on the landing page at T-0. 10,000 units of the headline SKU. 500 concurrent SKUs. Assume peak of 1 million requests per second."

**Candidate:** "Noted:

| Metric | Value |
|---|---|
| Concurrent users at sale start | 10,000,000 |
| Peak RPS | 1,000,000 |
| Headline SKU stock | 10,000 units |
| Concurrent sales (SKUs) | 500 |
| Sale duration | Seconds to minutes (stock-bound) |
| Reservation window | 10 minutes |
| Conversion (winners / attempts) | 10,000 / 10,000,000 = **0.1%** |

> **Critical insight for the interview:** 99.9% of requests are *destined to fail*. The system's primary job is to reject 9,990,000 requests as cheaply and as early as possible — ideally at the CDN or edge, never at the database."

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me model the traffic shape, because a flash sale is not steady-state traffic — it's an impulse."

```
Traffic profile (headline SKU):

  RPS
  1M ┤        ╷
     │        │
 800k┤        │╮
     │        ││
 600k┤        ││╮
     │        │││
 400k┤        ││ ╮
     │        ││  ╰╮
 200k┤        ││    ╰──╮
     │  ──────╯│        ╰────────────────
   0 ┤─────────┴─────────────────────────
      T-60s   T=0   T+5s  T+30s   T+5min

Pre-sale (T-60s to T=0):
  Page refreshes / countdown polls:  ~500,000 RPS  (100% cacheable)
  Stock status polls:                ~200,000 RPS  (cacheable, 1s TTL)

Sale open (T=0 to T+5s):
  Buy requests:                      ~1,000,000 RPS peak
  Total buy attempts:                ~10,000,000 over 5 seconds
  Successful reservations:           10,000  (0.1%)
  Rejected (sold out):               9,990,000 (99.9%)

Post sold-out (T+5s onwards):
  Buy requests still arriving:       ~200,000 RPS decaying
  → All served from the "SOLD OUT" edge flag, zero backend load

Payment phase (T+0 to T+10min):
  Payment initiations:               10,000 spread over 10 min ≈ 17 TPS
  → Trivial load; reuses the existing payment system
```

**Candidate:** "The key number is this: the buy path needs to survive 1M RPS, but only 10,000 requests ever need to reach the database. That's a **100:1 reduction at Redis, and a further 1000:1 reduction before the DB.** Everything in the design follows from that funnel."

### 3.2 The Request Funnel

```
  10,000,000 buy attempts
         │
         ▼  ┌─────────────────────────────────────────────┐
   [CDN / Edge]  Sold-out flag, bot fingerprint, geo block │  drops ~40%
         │  └─────────────────────────────────────────────┘
   6,000,000
         │
         ▼  ┌─────────────────────────────────────────────┐
   [API Gateway]  Rate limit per user/IP/device            │  drops ~50%
         │  └─────────────────────────────────────────────┘
   3,000,000
         │
         ▼  ┌─────────────────────────────────────────────┐
   [Admission Control]  Token/queue, stochastic shedding   │  drops ~97%
         │  └─────────────────────────────────────────────┘
     100,000
         │
         ▼  ┌─────────────────────────────────────────────┐
   [Redis Inventory]  Atomic Lua DECR — the true gate      │  drops ~90%
         │  └─────────────────────────────────────────────┘
      10,000  winners
         │
         ▼  ┌─────────────────────────────────────────────┐
   [Kafka → Order Service → PostgreSQL]                    │  10k writes total
         │  └─────────────────────────────────────────────┘
      10,000 orders created asynchronously

Each layer is 10-100x cheaper than the layer below it.
Design principle: REJECT AS EARLY AS POSSIBLE.
```

### 3.3 Storage Estimation

```
Per reservation record:        ~300 bytes
Per order record:              ~800 bytes
Per audit/event record:        ~400 bytes

Big Billion Days (5 days, 500 SKUs, ~5M total units sold):
  Reservations:   6M × 300 B  ≈ 1.8 GB
  Orders:         5M × 800 B  ≈ 4.0 GB
  Events (Kafka): 50M × 400 B ≈ 20 GB (7-day retention)
  Access logs:    ~2 TB (sampled + compressed for bot analysis)

Redis memory (hot path):
  Inventory counters:  500 SKUs × 20 shards × 64 B      ≈ 640 KB (trivial)
  Per-user dedup set:  10M users × 500 SKUs is too big  → use per-SKU
                       Bitmap keyed by user_id: 10M bits = 1.25 MB per SKU
                       500 SKUs × 1.25 MB                ≈ 625 MB
  Reservation objects: 10k × 300 B per SKU               ≈ 1.5 GB total
  → A modest Redis cluster (6 nodes × 16 GB) is more than enough.

→ Storage is NOT the bottleneck. Single-row write contention is.
```

### 3.4 Bandwidth Estimation

```
Landing page (fully static, CDN):   ~800 KB first load, ~2 KB for polls
  10M users × 800 KB = 8 TB served — but 99.99% from CDN edge, not origin

Buy request:   ~500 bytes  → 1M RPS × 500 B = 500 MB/s inbound at edge
Buy response:  ~200 bytes  → 200 MB/s outbound

Origin bandwidth after CDN offload: < 5% of the above
→ Bandwidth is manageable; CPU spent on connection setup (TLS handshakes)
  is the real edge-layer cost. Use connection reuse / HTTP-2 / keep-alive.
```

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌───────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                       │
│        Web (React)   |   Android/iOS App   |   Partner APIs                │
│   Client-side: server-synced countdown, button disabled until T=0,         │
│   exponential backoff on retry, no auto-refresh loops                      │
└──────────────────────────────┬────────────────────────────────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │   CDN  (Cloudflare / │  Static page, JS, images
                    │   Akamai / CloudFront)│  SOLD_OUT flag @ edge (1s TTL)
                    │   + Edge Workers      │  Bot fingerprint, geo/ASN block
                    └──────────┬───────────┘
                               │  (only ~40-60% gets through)
                    ┌──────────▼──────────┐
                    │    API Gateway       │  TLS termination, JWT auth,
                    │    (Kong / Envoy)    │  per-user + per-IP rate limit,
                    │                      │  request validation
                    └──────────┬───────────┘
                               │
                    ┌──────────▼──────────────────────────┐
                    │      ADMISSION CONTROL LAYER         │
                    │  • Virtual waiting room (queue token)│
                    │  • Stochastic load shedding          │
                    │  • Circuit breaker on downstream     │
                    └──────────┬───────────────────────────┘
                               │  (only ~1-2% gets through)
         ┌─────────────────────┼──────────────────────┐
         │                     │                       │
┌────────▼────────┐  ┌────────▼─────────┐  ┌─────────▼─────────┐
│  Sale Read       │  │  Flash Sale       │  │  Order Query      │
│  Service         │  │  Service (Buy)    │  │  Service          │
│  (stock, config) │  │  ★ HOT PATH ★     │  │  (status, history)│
└────────┬─────────┘  └────────┬──────────┘  └─────────┬─────────┘
         │                     │                        │
         │            ┌────────▼──────────────────────┐ │
         │            │   INVENTORY ENGINE (Redis)    │ │
         │            │   ┌─────────────────────────┐ │ │
         │            │   │ Lua script (atomic):    │ │ │
         │            │   │  1. dedup check (user)  │ │ │
         │            │   │  2. stock > 0 ?         │ │ │
         │            │   │  3. DECR stock          │ │ │
         │            │   │  4. create reservation  │ │ │
         │            │   │  → all-or-nothing       │ │ │
         │            │   └─────────────────────────┘ │ │
         │            │  Sharded counters: sku:{id}:s0│ │
         │            │  ... s19  (hot-key splitting) │ │
         │            └────────┬──────────────────────┘ │
         │                     │ WINNER ONLY (10k total) │
         │            ┌────────▼──────────┐              │
         │            │   Kafka            │  topic: reservation.created
         │            │   (buffer/decouple)│  partitioned by sku_id
         │            └────────┬───────────┘              │
         │                     │                          │
         │      ┌──────────────┼────────────┬─────────────┤
         │      │              │             │             │
┌────────▼──────▼──┐ ┌────────▼──────┐ ┌───▼──────────┐  │
│  Order Service    │ │ Notification  │ │ Analytics /  │  │
│  (creates order,  │ │ Service       │ │ Fraud        │  │
│   calls Payment)  │ │ (push/SMS)    │ │ Scoring      │  │
└────────┬──────────┘ └───────────────┘ └──────────────┘  │
         │                                                 │
┌────────▼─────────────────────────────────────────────────▼──┐
│                        DATA LAYER                            │
│  ┌───────────────┐  ┌───────────────┐  ┌─────────────────┐  │
│  │ PostgreSQL    │  │ Redis Cluster │  │  PostgreSQL     │  │
│  │ (orders,      │  │ (inventory,   │  │  (source-of-    │  │
│  │  reservations)│  │  dedup, cache)│  │   truth stock)  │  │
│  │  + replicas   │  │  + AOF persist│  │   reconciled    │  │
│  └───────────────┘  └───────────────┘  └─────────────────┘  │
└──────────────────────────────────────────────────────────────┘
         │
┌────────▼──────────────────────────────────────────────────┐
│  BACKGROUND JOBS                                           │
│  • Reservation expiry sweeper (releases stock back)        │
│  • Redis ↔ PostgreSQL stock reconciliation (every 30s)     │
│  • Cache pre-warmer (T-30min before each sale)             │
│  • Bot/fraud batch analysis                                │
└────────────────────────────────────────────────────────────┘
```

### 4.2 API Design

**Candidate:** "The APIs are deliberately thin. Note that the buy endpoint returns *immediately* — it does not wait for the order to be created."

#### Get Sale Configuration (heavily cached)
```
GET /api/v1/sales/{sale_id}
Cache-Control: public, max-age=60

Response (200 OK):
{
  "sale_id": "sale_bbd_pixel9",
  "sku_id": "sku_pixel9_128gb",
  "title": "Pixel 9 — Big Billion Days",
  "price_cents": 4999900,
  "starts_at": "2026-08-15T12:00:00Z",
  "server_time": "2026-08-15T11:59:42Z",   // client syncs countdown to this
  "status": "SCHEDULED",                    // SCHEDULED | LIVE | SOLD_OUT | ENDED
  "max_per_user": 1
}
```

#### Get Live Stock (cached at edge, 1s TTL)
```
GET /api/v1/sales/{sale_id}/stock

Response (200 OK):
{
  "sale_id": "sale_bbd_pixel9",
  "status": "LIVE",
  "stock_bucket": "LOW",      // HIGH | MEDIUM | LOW | SOLD_OUT
  "approx_remaining": 240,    // deliberately approximate
  "as_of": "2026-08-15T12:00:03Z"
}
```
> **Design note:** We return a *bucket*, not an exact number, and cache for 1 second. An exact live count would require a Redis read per request — 1M RPS of pure waste. Users cannot act on the difference between 240 and 237.

#### Attempt Purchase (the hot path)
```
POST /api/v1/sales/{sale_id}/purchase
Headers:
  Authorization: Bearer <jwt>
  Idempotency-Key: "user_9931_sale_bbd_pixel9"    // deterministic per user+sale
  X-Queue-Token: "qt_a83f...".                    // from waiting room, if enabled
  X-Device-Fingerprint: "fp_7d2c..."

Request Body:
{ "quantity": 1 }

Response — WINNER (202 Accepted):
{
  "result": "RESERVED",
  "reservation_id": "resv_7x8y9z",
  "expires_at": "2026-08-15T12:10:00Z",
  "payment_url": "/checkout/resv_7x8y9z",
  "order_status_url": "/api/v1/reservations/resv_7x8y9z"
}

Response — SOLD OUT (409 Conflict):
{ "result": "SOLD_OUT", "message": "All units have been claimed." }

Response — ALREADY PURCHASED (409 Conflict):
{ "result": "LIMIT_REACHED", "existing_reservation_id": "resv_7x8y9z" }

Response — SHED / TRY AGAIN (503 + Retry-After):
{ "result": "BUSY", "retry_after_ms": 800 }
```
> **Why 202 and not 201?** The unit is *reserved*, not sold. The order row does not exist yet — it's being created asynchronously off Kafka. Returning 202 is honest and lets us keep the write path async.

#### Confirm Purchase (payment)
```
POST /api/v1/reservations/{reservation_id}/confirm
Headers:
  Idempotency-Key: "resv_7x8y9z_confirm"

Request Body:
{ "payment_method": "pm_card_visa_4242", "address_id": "addr_112" }

Response (200 OK):
{ "order_id": "ord_abc123", "status": "CONFIRMED", "payment_id": "pay_7x8y9z" }
```

#### Waiting Room Token (when enabled)
```
GET /api/v1/sales/{sale_id}/queue

Response (200 OK):
{
  "position": 45231,
  "estimated_wait_seconds": 90,
  "token": null,                 // null until admitted
  "poll_after_ms": 5000
}

Response (when admitted):
{ "position": 0, "token": "qt_a83f...", "token_expires_at": "..." }
```

### 4.3 Purchase Flow — Data Flow

```
                    User taps "Buy Now" at T=0
                              │
                              ▼
                    ┌──────────────────┐
                    │ 1. CDN edge check │  Is SOLD_OUT flag set?
                    └────────┬──────────┘
                   ┌─────────┴─────────┐
             (yes) │                   │ (no)
                   ▼                   ▼
          Return SOLD_OUT      ┌──────────────────┐
          (never hits origin)  │ 2. Gateway:       │
                               │  auth + rate limit│
                               └────────┬──────────┘
                              ┌─────────┴─────────┐
                    (limited) │                   │ (ok)
                              ▼                   ▼
                     429 Too Many        ┌──────────────────┐
                                         │ 3. Admission:     │
                                         │  queue token /    │
                                         │  shed if overload │
                                         └────────┬──────────┘
                                        ┌─────────┴─────────┐
                                (shed)  │                   │ (admit)
                                        ▼                   ▼
                                503 BUSY          ┌──────────────────────┐
                                + Retry-After     │ 4. REDIS LUA (atomic) │
                                                  │  a) SETNX user dedup  │
                                                  │  b) stock > 0 ?       │
                                                  │  c) DECR stock        │
                                                  │  d) HSET reservation  │
                                                  │  e) EXPIRE 10 min     │
                                                  └────────┬──────────────┘
                                        ┌──────────────────┴────────┐
                              (stock=0) │                            │ (success)
                                        ▼                            ▼
                                  409 SOLD_OUT           ┌────────────────────┐
                                  + set edge flag        │ 5. Publish to Kafka│
                                                         │  reservation.created│
                                                         └────────┬───────────┘
                                                                  │
                                                         ┌────────▼───────────┐
                                                         │ 6. Return 202 to    │
                                                         │    user immediately │
                                                         │    (~15 ms total)   │
                                                         └────────┬───────────┘
                                                                  │
                                    ═══════ ASYNC FROM HERE ═══════
                                                                  │
                        ┌─────────────────────────────────────────┼───────────────┐
                        ▼                         ▼               ▼               ▼
              ┌──────────────────┐   ┌──────────────────┐  ┌───────────┐  ┌─────────────┐
              │ 7. Order Service │   │ Notification     │  │ Analytics │  │ Fraud       │
              │  creates PENDING │   │ (push: "You got  │  │ pipeline  │  │ scoring     │
              │  order in PG     │   │  it! Pay now")   │  │           │  │             │
              └────────┬─────────┘   └──────────────────┘  └───────────┘  └─────────────┘
                       │
              ┌────────▼─────────┐
              │ 8. User pays     │  → Payment System (see doc 007)
              │    within 10 min │
              └────────┬─────────┘
              ┌────────┴─────────┐
      (paid)  │                  │ (expired)
              ▼                  ▼
     Order CONFIRMED     Reservation expires
     Stock permanently   → Sweeper INCRs Redis stock
     consumed            → Unit returns to sale
                         → Edge SOLD_OUT flag cleared
```

---

## 5. Deep Dive: Core Components

### 5.1 Flash Sale Service (Hot Path)

**Candidate:** "This service does as little as possible. Its entire job is: validate, run one Lua script against Redis, publish one Kafka message, return. No database access on the hot path — none."

```java
@RestController
@RequestMapping("/api/v1/sales")
public class FlashSaleController {

    private final InventoryEngine inventoryEngine;
    private final KafkaTemplate<String, ReservationEvent> kafka;
    private final SaleConfigCache saleConfigCache;   // local in-process cache
    private final MeterRegistry metrics;

    @PostMapping("/{saleId}/purchase")
    public ResponseEntity<PurchaseResponse> purchase(
            @PathVariable String saleId,
            @RequestHeader("Idempotency-Key") String idempotencyKey,
            @AuthenticationPrincipal UserPrincipal user,
            @RequestBody PurchaseRequest request) {

        // 1. Local cache lookup — zero network hops
        SaleConfig config = saleConfigCache.get(saleId);
        if (config == null) {
            return ResponseEntity.notFound().build();
        }

        // 2. Cheap guards BEFORE touching Redis
        if (Instant.now().isBefore(config.startsAt())) {
            return ResponseEntity.badRequest()
                    .body(PurchaseResponse.notStarted());
        }
        if (config.soldOut()) {              // in-process flag, refreshed via pub/sub
            metrics.counter("flashsale.reject", "reason", "soldout_local").increment();
            return ResponseEntity.status(409).body(PurchaseResponse.soldOut());
        }
        if (request.quantity() > config.maxPerUser()) {
            return ResponseEntity.badRequest()
                    .body(PurchaseResponse.limitExceeded());
        }

        // 3. THE atomic operation — single Redis round trip
        ReservationResult result = inventoryEngine.tryReserve(
                saleId, user.getId(), request.quantity(), config.reservationTtlSeconds());

        switch (result.status()) {
            case SOLD_OUT -> {
                saleConfigCache.markSoldOut(saleId);   // trip local + edge flag
                metrics.counter("flashsale.reject", "reason", "soldout").increment();
                return ResponseEntity.status(409).body(PurchaseResponse.soldOut());
            }
            case ALREADY_RESERVED -> {
                // Idempotent: same user retrying gets their existing reservation
                return ResponseEntity.status(409)
                        .body(PurchaseResponse.limitReached(result.reservationId()));
            }
            case RESERVED -> {
                // 4. Fire-and-forget to Kafka (async, non-blocking send)
                kafka.send("reservation.created", saleId,
                        new ReservationEvent(result.reservationId(), user.getId(),
                                saleId, config.skuId(), request.quantity(),
                                result.expiresAt(), idempotencyKey));

                metrics.counter("flashsale.reserved").increment();
                return ResponseEntity.accepted().body(
                        PurchaseResponse.reserved(result.reservationId(), result.expiresAt()));
            }
            default -> {
                return ResponseEntity.status(503).body(PurchaseResponse.busy(800));
            }
        }
    }
}
```

> **Interview talking point:** Notice there is no `@Transactional`, no JPA, no join, no lock. The hot path is one Redis call. Everything else is pushed to the async side. If someone asks "why not just use a database transaction?", the answer is in §7.1.

### 5.2 Inventory Engine — Atomic Reservation via Lua

**Candidate:** "This is the correctness core of the entire system. Redis is single-threaded for command execution, and a Lua script runs atomically — no other command interleaves. That gives us a compare-and-decrement with dedup in one shot."

```java
@Component
public class InventoryEngine {

    private static final String RESERVE_SCRIPT = """
        -- KEYS[1] = stock counter key      e.g. flashsale:{sale_1}:stock:shard3
        -- KEYS[2] = user dedup set key     e.g. flashsale:{sale_1}:users
        -- KEYS[3] = reservation hash key   e.g. flashsale:{sale_1}:resv:<uuid>
        -- ARGV[1] = user_id
        -- ARGV[2] = quantity
        -- ARGV[3] = reservation_id
        -- ARGV[4] = ttl seconds
        -- ARGV[5] = now (epoch millis)

        -- Step 1: per-user dedup. SADD returns 0 if the user is already there.
        if redis.call('SADD', KEYS[2], ARGV[1]) == 0 then
            return {-2, ''}                       -- ALREADY_RESERVED
        end

        -- Step 2: check stock
        local stock = tonumber(redis.call('GET', KEYS[1]) or '0')
        local qty = tonumber(ARGV[2])
        if stock < qty then
            redis.call('SREM', KEYS[2], ARGV[1])  -- roll back the dedup entry
            return {-1, ''}                       -- SOLD_OUT
        end

        -- Step 3: decrement stock (atomic with everything above)
        redis.call('DECRBY', KEYS[1], qty)

        -- Step 4: write the reservation record with a TTL
        redis.call('HSET', KEYS[3],
            'user_id', ARGV[1],
            'quantity', ARGV[2],
            'status', 'PENDING_PAYMENT',
            'created_at', ARGV[5])
        redis.call('EXPIRE', KEYS[3], ARGV[4])

        -- Step 5: index for the expiry sweeper (score = expiry timestamp)
        redis.call('ZADD', KEYS[2] .. ':expiry',
            tonumber(ARGV[5]) + (tonumber(ARGV[4]) * 1000), ARGV[3])

        return {1, ARGV[3]}                       -- RESERVED
        """;

    private final RedisTemplate<String, String> redis;
    private final DefaultRedisScript<List> script;
    private final int shardCount = 20;

    public ReservationResult tryReserve(String saleId, String userId,
                                        int qty, int ttlSeconds) {
        String reservationId = "resv_" + UUID.randomUUID();

        // Pick a random shard, then fall through to others if empty (see §7.3)
        int start = ThreadLocalRandom.current().nextInt(shardCount);
        for (int i = 0; i < shardCount; i++) {
            int shard = (start + i) % shardCount;
            List<?> out = redis.execute(script,
                List.of(stockKey(saleId, shard), userKey(saleId), resvKey(saleId, reservationId)),
                userId, String.valueOf(qty), reservationId,
                String.valueOf(ttlSeconds), String.valueOf(System.currentTimeMillis()));

            long code = ((Number) out.get(0)).longValue();
            if (code == 1)  return ReservationResult.reserved(reservationId,
                                    Instant.now().plusSeconds(ttlSeconds));
            if (code == -2) return ReservationResult.alreadyReserved(reservationId);
            // code == -1 → this shard is empty; try the next one
        }
        return ReservationResult.soldOut();
    }

    // Hash tags {saleId} keep all keys for one sale on the same Redis Cluster slot,
    // which is REQUIRED for multi-key Lua scripts.
    private String stockKey(String saleId, int shard) {
        return "flashsale:{" + saleId + "}:stock:" + shard;
    }
    private String userKey(String saleId) {
        return "flashsale:{" + saleId + "}:users";
    }
    private String resvKey(String saleId, String resvId) {
        return "flashsale:{" + saleId + "}:resv:" + resvId;
    }
}
```

> **Three details that impress interviewers:**
> 1. **Rollback inside the script** — if stock check fails after `SADD`, we `SREM` so the user isn't wrongly blocked from retrying after a restock.
> 2. **Hash tags `{saleId}`** — Redis Cluster requires all keys in a multi-key operation to hash to the same slot. Forgetting this is the #1 production bug when moving from standalone Redis to a cluster.
> 3. **The dedup structure** — a Set works, but at 10M users a **Bitmap keyed by numeric user id** (`SETBIT`) uses ~1.25 MB instead of ~600 MB. Mention this as the optimization.

### 5.3 Admission Control (Virtual Waiting Room)

**Candidate:** "Rate limiting rejects users. A waiting room *defers* them — much better UX, and it converts a spike into a controlled stream."

```java
@Component
public class AdmissionController {

    private final RedisTemplate<String, String> redis;

    /**
     * Token bucket refilled at a rate we can actually serve.
     * If we can handle 50k RPS at the inventory layer, we admit 50k tokens/sec
     * regardless of whether 1M or 5M people are waiting.
     */
    private static final String ADMIT_SCRIPT = """
        local key       = KEYS[1]
        local now       = tonumber(ARGV[1])
        local rate      = tonumber(ARGV[2])   -- tokens per second
        local capacity  = tonumber(ARGV[3])   -- burst capacity

        local bucket    = redis.call('HMGET', key, 'tokens', 'last_refill')
        local tokens    = tonumber(bucket[1]) or capacity
        local last      = tonumber(bucket[2]) or now

        -- refill based on elapsed time
        local elapsed   = math.max(0, now - last) / 1000.0
        tokens = math.min(capacity, tokens + (elapsed * rate))

        if tokens < 1 then
            redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)
            return 0                          -- not admitted, keep waiting
        end

        redis.call('HMSET', key, 'tokens', tokens - 1, 'last_refill', now)
        redis.call('EXPIRE', key, 3600)
        return 1                              -- admitted
        """;

    public boolean tryAdmit(String saleId, int ratePerSecond, int burst) {
        Long result = redis.execute(admitScript,
            List.of("flashsale:{" + saleId + "}:admission"),
            String.valueOf(System.currentTimeMillis()),
            String.valueOf(ratePerSecond), String.valueOf(burst));
        return result != null && result == 1L;
    }
}
```

**Waiting room mechanics (what the user sees):**

```
User arrives at T-0
   │
   ▼
Assigned a queue position via an atomic Redis INCR on a counter,
signed into a JWT so it can't be forged:
   position = INCR flashsale:{sale}:queue:counter    → 45,231

Client polls GET /queue every 5s with its signed position.
Server publishes an "admitted watermark" (e.g. 60,000) every second.

   if position <= watermark:
       issue a short-lived queue token (HMAC-signed, 60s TTL, single-use)
       → client may now call /purchase
   else:
       show "You are #45,231 in line — about 90 seconds"

Watermark advances at the rate the backend can absorb:
   watermark += admission_rate_per_second

Anti-abuse: the token is bound to (user_id, sale_id, position) and is
consumed atomically on first use, so it cannot be shared or replayed.
```

> **Why a waiting room beats naive rate limiting:** with a 429, ten million clients all retry at once and you get a synchronized retry storm — a self-inflicted DDoS. With a queue, each client has a position and a *time* to come back, so the load is naturally spread.

### 5.4 Order Service (Async Consumer)

**Candidate:** "This runs off Kafka, so it can be slow. It's the first place a relational database appears in the flow."

```java
@Component
public class ReservationConsumer {

    private final OrderRepository orderRepo;
    private final NotificationService notifications;

    @KafkaListener(topics = "reservation.created", groupId = "order-service",
                   concurrency = "12")
    @Transactional
    public void onReservation(ReservationEvent event, Acknowledgment ack) {
        try {
            // Idempotent insert — the consumer may see the same event twice
            // (at-least-once delivery). The unique constraint makes retries safe.
            Order order = Order.builder()
                    .id(UUID.randomUUID())
                    .idempotencyKey(event.idempotencyKey())
                    .reservationId(event.reservationId())
                    .userId(event.userId())
                    .skuId(event.skuId())
                    .quantity(event.quantity())
                    .status(OrderStatus.PENDING_PAYMENT)
                    .expiresAt(event.expiresAt())
                    .build();

            orderRepo.insertIgnoreConflict(order);   // ON CONFLICT DO NOTHING

            notifications.sendPurchaseWindowOpened(event.userId(), order.getId());
            ack.acknowledge();

        } catch (TransientDataAccessException e) {
            // Do NOT ack — Kafka will redeliver. Reservation TTL gives us
            // a 10-minute budget to recover, which is enormous.
            throw e;
        }
    }
}
```

> **Key property:** if the Order Service is completely down for 3 minutes, nothing is lost. Kafka buffers, reservations are still held in Redis with a 10-minute TTL, and the winners are still winners. The hot path never knew.

### 5.5 Reservation Expiry Sweeper

```java
@Component
public class ReservationSweeper {

    private static final String RELEASE_SCRIPT = """
        -- Find reservations whose expiry has passed, release their stock
        local expired = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, 500)
        local released = 0
        for i, resvId in ipairs(expired) do
            local resvKey = KEYS[2] .. resvId
            local status  = redis.call('HGET', resvKey, 'status')
            -- Only release if still unpaid. Paid ones are simply de-indexed.
            if status == 'PENDING_PAYMENT' or status == false then
                local qty = tonumber(redis.call('HGET', resvKey, 'quantity') or '0')
                local uid = redis.call('HGET', resvKey, 'user_id')
                if qty > 0 then
                    redis.call('INCRBY', KEYS[3], qty)   -- stock back to shard 0
                    released = released + qty
                end
                if uid then
                    redis.call('SREM', KEYS[4], uid)     -- let the user try again
                end
                redis.call('DEL', resvKey)
            end
            redis.call('ZREM', KEYS[1], resvId)
        end
        return released
        """;

    @Scheduled(fixedDelay = 5000)
    @SchedulerLock(name = "reservationSweeper")   // ShedLock: only one node runs it
    public void sweep() {
        for (String saleId : activeSaleIds()) {
            Long released = redis.execute(releaseScript, keysFor(saleId),
                    String.valueOf(System.currentTimeMillis()));
            if (released != null && released > 0) {
                log.info("Released {} units back to sale {}", released, saleId);
                saleEvents.publishRestocked(saleId);   // clears edge SOLD_OUT flag
                metrics.counter("flashsale.released", "sale", saleId).increment(released);
            }
        }
    }
}
```

> **Subtle correctness point:** the sweeper must be idempotent and must not double-release. Doing the status check, the `INCRBY`, and the `DEL` inside one Lua script guarantees that two sweeper instances can never release the same reservation twice. `@SchedulerLock` is belt-and-braces on top.

---

## 6. Deep Dive: Database Design

### 6.1 Schema

**Candidate:** "PostgreSQL holds the durable truth. Redis holds the fast truth. They reconcile continuously."

```sql
-- Sale configuration (read-mostly, aggressively cached)
CREATE TABLE flash_sales (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id             UUID NOT NULL REFERENCES products(id),
    title              TEXT NOT NULL,
    price_cents        BIGINT NOT NULL CHECK (price_cents > 0),
    total_stock        INTEGER NOT NULL CHECK (total_stock >= 0),
    reserved_count     INTEGER NOT NULL DEFAULT 0,
    sold_count         INTEGER NOT NULL DEFAULT 0,
    max_per_user       SMALLINT NOT NULL DEFAULT 1,
    starts_at          TIMESTAMPTZ NOT NULL,
    ends_at            TIMESTAMPTZ NOT NULL,
    status             VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED',
    reservation_ttl_s  INTEGER NOT NULL DEFAULT 600,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_status CHECK (
        status IN ('DRAFT','SCHEDULED','LIVE','SOLD_OUT','ENDED','CANCELLED')),
    CONSTRAINT no_oversell CHECK (reserved_count + sold_count <= total_stock)
);

CREATE INDEX idx_sales_starts   ON flash_sales (starts_at) WHERE status = 'SCHEDULED';
CREATE INDEX idx_sales_live     ON flash_sales (status) WHERE status = 'LIVE';

-- Reservations: durable mirror of the Redis reservation
CREATE TABLE reservations (
    id              UUID PRIMARY KEY,
    sale_id         UUID NOT NULL REFERENCES flash_sales(id),
    user_id         UUID NOT NULL,
    quantity        SMALLINT NOT NULL CHECK (quantity > 0),
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING_PAYMENT',
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- The database-level guarantee of "one per user per sale"
    CONSTRAINT uq_user_sale UNIQUE (sale_id, user_id),
    CONSTRAINT valid_resv_status CHECK (
        status IN ('PENDING_PAYMENT','CONFIRMED','EXPIRED','CANCELLED'))
);

CREATE INDEX idx_resv_expiry ON reservations (expires_at)
    WHERE status = 'PENDING_PAYMENT';

-- Orders
CREATE TABLE orders (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key  VARCHAR(255) NOT NULL UNIQUE,
    reservation_id   UUID NOT NULL REFERENCES reservations(id),
    user_id          UUID NOT NULL,
    sku_id           UUID NOT NULL,
    quantity         SMALLINT NOT NULL,
    amount_cents     BIGINT NOT NULL,
    currency         CHAR(3) NOT NULL DEFAULT 'INR',
    status           VARCHAR(20) NOT NULL DEFAULT 'PENDING_PAYMENT',
    payment_id       VARCHAR(255),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_order_status CHECK (
        status IN ('PENDING_PAYMENT','CONFIRMED','PAYMENT_FAILED',
                   'EXPIRED','CANCELLED','FULFILLED'))
);

CREATE INDEX idx_orders_user  ON orders (user_id, created_at DESC);
CREATE INDEX idx_orders_sale  ON orders (sku_id, created_at DESC);

-- Inventory audit log: append-only, explains every unit's journey
CREATE TABLE inventory_events (
    id          BIGSERIAL PRIMARY KEY,
    sale_id     UUID NOT NULL,
    event_type  VARCHAR(30) NOT NULL,  -- RESERVED | RELEASED | SOLD | RESTOCKED | ADJUSTED
    quantity    INTEGER NOT NULL,
    reservation_id UUID,
    user_id     UUID,
    source      VARCHAR(20) NOT NULL,  -- REDIS | SWEEPER | ADMIN | RECONCILER
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inv_events_sale ON inventory_events (sale_id, created_at DESC);
```

### 6.2 The `no_oversell` Constraint — Your Safety Net

**Candidate:** "The `CHECK (reserved_count + sold_count <= total_stock)` constraint is the last line of defence. If a Redis bug, a bad deploy, or a manual restock ever lets an extra reservation slip through, the database write fails loudly instead of silently overselling. We'd rather lose the order and alert than ship a phone we don't have."

```sql
-- Confirming a sale: reserved → sold, atomic and oversell-proof
UPDATE flash_sales
SET reserved_count = reserved_count - 1,
    sold_count     = sold_count + 1,
    updated_at     = NOW()
WHERE id = $1
  AND reserved_count >= 1;
-- 0 rows affected → something is wrong; alert, do not proceed.
```

### 6.3 Why NOT a Pure Database Solution

**Candidate:** "Interviewers usually ask why we can't just do this in Postgres. Let me show the math."

```sql
-- The "obvious" approach — atomic and correct, but it does not scale
UPDATE flash_sales
SET reserved_count = reserved_count + 1
WHERE id = 'sale_1'
  AND reserved_count + sold_count < total_stock;
```

```
This statement IS correct — the WHERE clause makes it a compare-and-swap
and Postgres row locks serialize it. The problem is throughput:

  Every concurrent transaction contends on ONE row.
  Row lock hold time (network + WAL flush) ≈ 1-2 ms.
  Max throughput = 1 / 0.0015 s ≈ 600-1,000 updates/sec.

At 1,000,000 RPS:
  - Connection pool (say 500 connections) is exhausted in microseconds
  - 999,000 requests queue on the lock
  - Lock wait times explode → statement timeouts → connection churn
  - The database becomes unresponsive for EVERY other query
  - The entire marketplace goes down because of one phone

Redis with a Lua script:
  - Single-threaded execution = inherent serialization, no lock manager
  - ~50 µs per script
  - Realistic throughput: 100,000+ ops/sec per node
  - Sharded across 20 counters: 1,000,000+ ops/sec
  → 1000x better, and it protects the database entirely.
```

> **Say this in the interview:** "The database approach isn't *wrong*, it's *slow*. And in a flash sale, slow is a correctness problem too — because a database that's stuck holding locks will start timing out, and timeouts create uncertain states, which is how oversell bugs actually get introduced in production."

---

## 7. Deep Dive: Inventory Engine (The Heart)

### 7.1 Why Redis Is the Source of Truth *During* the Sale

```
             During sale                 After sale
  ┌──────────────────────────┐   ┌──────────────────────────┐
  │  Redis = authoritative    │   │ PostgreSQL = authoritative│
  │  PostgreSQL = eventually  │──▶│ Redis = discarded / TTL'd │
  │    consistent mirror      │   │                           │
  └──────────────────────────┘   └──────────────────────────┘

Justification:
  • The value we protect is a COUNTER, not rich relational data.
  • The window is short (seconds to minutes).
  • Redis AOF with appendfsync=everysec bounds worst-case loss to 1 second.
  • A reconciler compares Redis and PG every 30s and alerts on drift.
  • Undersell (losing a unit) is acceptable; oversell is not — and a
    Redis crash causes UNDERSELL, not oversell, because the recovered
    counter is <= the true remaining stock.
```

> **This asymmetry is the single best point you can make in this interview.** Every failure mode of Redis biases toward *undersell*, which the business accepts. Design your failure modes so they fail in the direction the business tolerates.

### 7.2 The Hot Key Problem

```
Problem: ALL 1M RPS hit ONE Redis key: flashsale:sale_1:stock

  ┌──────────────┐
  │ Redis Node 3 │ ◀── 1,000,000 RPS   💀 CPU pinned at 100%
  │  (owns slot  │                        Other 5 nodes idle
  │   for that   │
  │   one key)   │
  └──────────────┘

Redis Cluster shards by KEY HASH. One key = one node. Adding nodes
does not help. This is the "hot key" / "celebrity key" problem.
```

### 7.3 Solution — Inventory Sharding (Stock Splitting)

```
Split 10,000 units across 20 logical shards of 500 each:

  flashsale:{sale_1}:stock:0   = 500
  flashsale:{sale_1}:stock:1   = 500
  ...
  flashsale:{sale_1}:stock:19  = 500
                                ─────
                        Total  = 10,000

Each request picks a random shard → load spreads across the cluster.
Throughput scales roughly linearly with shard count.

  Request → hash(random) → shard 7 → DECR flashsale:{sale_1}:stock:7

Edge case: SHARD EXHAUSTION
  Shard 7 hits 0 while shard 12 still has 300 units.
  A naive implementation reports SOLD_OUT with 300 units unsold.

Fix — sequential fallback (implemented in §5.2):
  Try shard 7 → empty
  Try shard 8 → empty
  Try shard 9 → success ✅
  Only after all 20 shards are empty do we report SOLD_OUT.

Cost: the *last* few units cost up to 20 Redis calls to find.
      That is fine — by then traffic has collapsed and only a
      handful of requests are still winning.

Optimization: maintain a Redis bitmap of non-empty shards; when a
shard hits zero, clear its bit so later requests skip it entirely.
```

**Trade-off table (be ready for this):**

| Shards | Throughput | Fallback cost | Fairness |
|---|---|---|---|
| 1 | ~100k ops/s | none | perfect FCFS |
| 20 | ~1M ops/s | up to 20 hops at tail | approximate |
| 100 | ~2M ops/s (cluster-bound) | up to 100 hops at tail | approximate |

> **Rule of thumb:** shard count ≈ number of Redis cluster nodes × 3-4. More shards than that buys throughput you can't use and makes the tail worse.

### 7.4 Read Path — Never Read Stock from Redis on the Hot Path

```
NAIVE (do not do this):
  Every page poll → GET flashsale:sale_1:stock → 1M RPS of reads on the hot key
  → You've recreated the hot key problem on the read side.

CORRECT:
  A background publisher reads the summed stock ONCE per second and
  pushes it to:
    1. The CDN edge as a cached JSON (1s TTL)   → serves ~all traffic
    2. Redis pub/sub → each app instance's local in-process cache
    3. WebSocket/SSE fan-out for the "Only 240 left!" ticker

  ┌────────────┐  1/sec   ┌────────────┐  push  ┌─────────────┐
  │ Stock      │─────────▶│ Edge cache  │──────▶│ 10M clients │
  │ Publisher  │          │ (1s TTL)    │        │             │
  └────────────┘          └────────────┘        └─────────────┘

  10,000,000 reads/sec → 1 Redis read/sec.  A 10-million-to-one reduction.

Bucketing: return HIGH / MEDIUM / LOW / SOLD_OUT instead of exact numbers.
  • Cheaper to cache (fewer distinct values → higher hit ratio)
  • Prevents scrapers from inferring exact sale velocity
  • Users cannot act on the difference between 240 and 237 anyway
```

### 7.5 Sold-Out Short-Circuit (The Highest-Leverage Optimization)

```
The moment stock hits 0, propagate a SOLD_OUT flag OUTWARD as fast as possible:

  Redis (t=0ms)
     │  Lua script returns stock == 0
     ▼
  Publish to Redis pub/sub (t≈1ms)
     │
     ├──▶ Every app instance sets a local boolean (t≈5ms)
     │      → subsequent requests rejected in-process, ZERO Redis calls
     │
     ├──▶ Update the edge KV store / cache purge (t≈50-500ms)
     │      → subsequent requests rejected AT THE CDN, ZERO origin calls
     │
     └──▶ WebSocket push to connected clients (t≈100ms)
            → button turns grey, most users never send the request

Impact: after ~500 ms, essentially 100% of the remaining traffic —
potentially hundreds of millions of requests over the next minute —
is absorbed by the CDN at near-zero cost.

Reverse propagation: when the sweeper releases expired reservations,
clear the flag the same way. Publish RESTOCKED → clear local flags →
purge edge cache. Cheap, and it prevents undersell.
```

### 7.6 Preventing Oversell — The Full Chain

```
Layer 1: Redis Lua script
   Atomic check-and-decrement. Single-threaded execution = no race.
   ✅ Primary guarantee

Layer 2: Per-user dedup (SADD / SETBIT inside the same script)
   One reservation per user per sale, enforced atomically with the decrement.
   ✅ Prevents a single user from sweeping stock with parallel requests

Layer 3: PostgreSQL UNIQUE (sale_id, user_id) on reservations
   Durable enforcement of the same rule if Redis state is lost.
   ✅ Backstop

Layer 4: PostgreSQL CHECK (reserved_count + sold_count <= total_stock)
   The write physically cannot succeed if it would oversell.
   ✅ Last line of defence — fails loudly

Layer 5: Reconciler (every 30 s)
   SUM(redis shards) + reserved + sold == total_stock ?
   Drift → P0 alert, and auto-pause the sale if drift > threshold.
   ✅ Detection

Together: a single-layer bug cannot cause oversell. It takes a
coordinated failure of two independent mechanisms.
```

### 7.7 Redis ↔ PostgreSQL Reconciliation

```java
@Component
public class InventoryReconciler {

    @Scheduled(fixedDelay = 30_000)
    @SchedulerLock(name = "inventoryReconciler")
    public void reconcile() {
        for (FlashSale sale : saleRepo.findLive()) {
            long redisRemaining = sumShards(sale.getId());
            long dbReserved     = reservationRepo.countActive(sale.getId());
            long dbSold         = orderRepo.countConfirmed(sale.getId());

            long accounted = redisRemaining + dbReserved + dbSold;
            long drift     = sale.getTotalStock() - accounted;

            if (drift == 0) continue;

            // Drift is EXPECTED to be small and positive during the sale:
            // a reservation exists in Redis but its Kafka event hasn't been
            // consumed into PG yet. That's in-flight, not a bug.
            long inFlight = kafkaLagEstimator.pendingFor(sale.getId());

            if (Math.abs(drift) - inFlight > sale.getTotalStock() * 0.001) {
                alerts.page(Severity.P0,
                    "INVENTORY DRIFT on sale %s: expected %d, accounted %d (drift %d)"
                        .formatted(sale.getId(), sale.getTotalStock(), accounted, drift));

                if (drift < 0) {
                    // Negative drift = we have accounted for MORE than exists.
                    // This is the oversell signature. Halt immediately.
                    saleService.pause(sale.getId(), "RECONCILIATION_DRIFT");
                }
            }
            metrics.gauge("flashsale.drift", drift);
        }
    }
}
```

### 7.8 Bot & Scalper Mitigation

```
Bots are the reason a "fair" flash sale feels unfair. Layered defence:

Layer 1 — Identity friction (before the sale)
  • Verified phone number required to participate
  • Account age > 30 days, or prior purchase history
  • One entry per verified identity, not per account

Layer 2 — Client attestation
  • Device fingerprint (Canvas/WebGL/font entropy) + Play Integrity /
    App Attest on mobile
  • Signed request payload — a raw curl to the API is rejected outright
  • Proof-of-work challenge (~200 ms of hashing) — invisible to a human,
    but multiplies a bot farm's cost by the number of accounts it runs

Layer 3 — Behavioural signals (real time, < 5 ms)
  • Time from page load to click < 50 ms → almost certainly automated
  • No mouse movement / touch events before the click
  • Identical fingerprints across many accounts
  • Datacenter ASN (AWS/GCP/Azure IP ranges) → high suspicion

Layer 4 — Shadow rejection (important!)
  Do NOT tell a bot it was blocked. Return a normal-looking SOLD_OUT.
  If you return "BOT_DETECTED", the operator instantly knows which
  signal fired and iterates. Silent failure is far more expensive
  for them to debug.

Layer 5 — Post-sale enforcement
  • Graph analysis: shared addresses, payment instruments, devices
  • Cancel and restock confirmed scalper orders
  • The restocked units go into a second drop for flagged-clean users
```

> **Honest framing for the interview:** "You cannot win this outright — a determined operator with 10,000 verified SIMs will get through. The goal is to raise cost per successful purchase above the resale margin. That's an economics problem, not a purely technical one." Interviewers like this answer because it shows commercial judgment.

---

## 8. Scaling Strategies

### 8.1 Static-First Page Architecture

```
The landing page is split into two parts:

  STATIC (99.9% of bytes) — pushed to CDN 30 minutes before the sale
    HTML shell, JS bundle, CSS, product images, specifications, reviews
    Cache-Control: public, max-age=3600, immutable
    → 10M users, ~0 origin requests

  DYNAMIC (a few hundred bytes) — fetched separately
    { server_time, sale_status, stock_bucket }
    Cache-Control: public, max-age=1
    → 10M users, ~1 origin request per second per edge POP

Countdown timer MUST sync to server_time, never to the device clock.
Otherwise a user whose phone is 20 seconds fast fires at T-20s and
gets a stream of errors — and every such user retries in a loop.
```

### 8.2 Pre-Warming (T-30 Minutes)

```
A cold start under 1M RPS is a guaranteed outage. Warm everything:

  T-30 min:  Scale app pods 10x (K8s HPA cannot react fast enough to a
             0→1M RPS step function — pre-scale explicitly, or use
             KEDA with a scheduled trigger)
  T-30 min:  Load sale config into every pod's local cache
  T-25 min:  Push static assets to CDN, verify edge hit ratio > 99%
  T-20 min:  Initialize Redis inventory shards, verify sum == total_stock
  T-15 min:  Warm JVM — send synthetic traffic so JIT compiles the hot path
             (a cold JVM interprets bytecode for the first ~10k invocations;
              at 1M RPS that is a disaster in the first second)
  T-10 min:  Warm DB connection pools; pre-create Kafka topic partitions
  T-5  min:  Freeze deploys. No config changes. Read-only change window.
  T-1  min:  On-call bridge open, dashboards up, kill switch tested
  T=0:       Flip status to LIVE via a single atomic flag
```

> **JVM warmup is a genuinely underrated point in Java interviews.** Mention tiered compilation, and that a cold JVM can be 10-20x slower than a warmed one. AOT (`-XX:AOTCache`) or CRaC can help if the platform supports it.

### 8.3 Horizontal Scaling of the Buy Path

```
Flash Sale Service is completely STATELESS:
  • No session affinity required
  • No local locks
  • All shared state is in Redis
  → scales linearly by adding pods

Sizing:
  Per pod: ~5,000 RPS (mostly one Redis round trip + one async Kafka send)
  For 1M RPS admitted: ... but we don't admit 1M RPS.
  Admission control caps the inventory layer at ~100k RPS.
  → 100,000 / 5,000 = 20 pods, plus 2x headroom = 40 pods.

The waiting room absorbs the difference. This is the whole point:
  don't scale to peak demand, scale to peak SERVICEABLE demand
  and buffer the rest.
```

### 8.4 Kafka Partitioning

```
Topic: reservation.created
  Partitions: 50   (must handle a 10k-message burst in ~5 seconds)
  Key: sale_id     → all events for one sale go to one partition,
                     preserving per-sale ordering
  Replication: 3, min.insync.replicas = 2
  acks = 1 on the producer (NOT acks=all)

Why acks=1 on the hot path?
  acks=all adds ~5-10 ms of latency per send. At the buy path's latency
  budget, that matters. The trade-off: a broker failure could lose a
  handful of reservation events.

  Is that acceptable? YES — because the reservation still exists in
  Redis with its TTL. A lost Kafka event means the order row isn't
  created; the reservation expires; the stock is released. We UNDERSELL
  by one unit. Again: every failure biases toward undersell.

  If a hot SKU has one partition and that becomes the bottleneck, key by
  (sale_id, hash(user_id) % 4) to get 4 partitions per sale — ordering
  across users doesn't matter here.
```

### 8.5 Multi-Region

```
Flash sales are usually region-bound (Big Billion Days = India), so
prefer a SINGLE-REGION inventory authority:

  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
  │  Edge: APAC  │    │ Edge: EU    │    │ Edge: US    │
  │  (CDN + WAF) │    │             │    │             │
  └──────┬───────┘    └──────┬──────┘    └──────┬──────┘
         │                   │                  │
         └───────────────────┼──────────────────┘
                             ▼
                  ┌──────────────────────┐
                  │  ap-south-1 (Mumbai) │  ← single inventory authority
                  │  Redis + App + Kafka │
                  └──────────────────────┘

Why not active-active inventory? Because a distributed counter across
regions needs consensus (Raft/Paxos) on every decrement — that's tens of
milliseconds per operation, at which point you're slower than PostgreSQL.

If you genuinely need multi-region inventory, PARTITION the stock:
  ap-south-1: 6,000 units    eu-west-1: 2,000 units    us-east-1: 2,000 units
Each region owns its slice absolutely — no cross-region coordination.
Rebalance between regions every 30 s via a slow async job.
Trade-off: one region may sell out while another has stock (regional
undersell). Acceptable for most businesses; state it explicitly.
```

### 8.6 Protecting the Rest of the Marketplace (Bulkheading)

```
NON-NEGOTIABLE: a flash sale must never take down normal shopping.

  ┌──────────────────────────────────────────────────────┐
  │  Separate everything:                                 │
  │                                                       │
  │  Flash Sale:   own K8s node pool, own Redis cluster,  │
  │                own DB (or at minimum own connection   │
  │                pool + separate replica), own Kafka    │
  │                cluster, own CDN hostname              │
  │                                                       │
  │  Core Commerce: untouched, unaware, unaffected        │
  └──────────────────────────────────────────────────────┘

Bulkhead patterns applied:
  • Separate thread pools (Resilience4j Bulkhead) per downstream
  • Separate connection pools with hard caps
  • Circuit breakers on every cross-service call
  • Rate limits on flash-sale → core-service calls (e.g. user lookup)
  • A kill switch that disables all flash sales in < 5 seconds

Test this: run a game day where you deliberately overload the flash
sale path and verify that core checkout p99 is unchanged.
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Redis Node Failure Mid-Sale

```
Scenario: The Redis primary holding shards 0-6 crashes at T+2s,
          with 3,000 units already sold from those shards.

Timeline:
  T+2.0s: Primary dies
  T+2.1s: Sentinel/Cluster detects failure
  T+2.5s: Replica promoted to primary
  T+2.5s: Replica's state may lag by up to ~1 second of writes

Impact analysis:
  Replica may show MORE stock than it should (it missed some DECRs)
  → risk of OVERSELL on those shards. This is the dangerous direction!

Mitigation stack:
  1. WAIT command: `WAIT 1 50` after the reservation script — blocks
     until at least 1 replica acknowledges, with a 50 ms timeout.
     Adds ~1-5 ms; eliminates most of the loss window.
  2. AOF with appendfsync=everysec on all nodes (bounded loss).
  3. On failover, the reconciler runs IMMEDIATELY (triggered by a
     cluster event, not the 30 s schedule) and hard-corrects the
     counters from PostgreSQL:
         redis_stock = total_stock - db_reserved - db_sold - safety_margin
     The safety margin (a few units) guarantees we correct DOWNWARD.
  4. If drift exceeds the threshold → pause the sale, correct, resume.
     A 10-second pause is vastly better than 400 cancelled orders.

Alternative for zero-tolerance cases: use Redis with WAIT + a
PostgreSQL confirmation write before returning 202 to the user. You
trade ~20 ms of latency for a stronger guarantee. Worth it for
high-value SKUs (a ₹2 lakh laptop), not for a ₹499 accessory.
```

### 9.2 Thundering Herd on Retry

```
Scenario: 1M users get 503 BUSY. Every client retries after 1 second.
          At T+1s, another 1M RPS arrives. The system never recovers.

           RPS
           1M ┤ ╷    ╷    ╷    ╷      ← synchronized retry waves
              │ │    │    │    │
              │ │    │    │    │
            0 ┤─┴────┴────┴────┴──
                T  T+1  T+2  T+3

Mitigations:
  1. Retry-After header with JITTER computed server-side:
       retry_after_ms = base * (1 + random(0, 1))
     Different clients get different values → the wave flattens.
  2. Exponential backoff with full jitter in the client SDK:
       sleep = random(0, min(cap, base * 2^attempt))
  3. The waiting room replaces retries entirely — the client holds a
     position and polls on a schedule the server controls.
  4. Cap retries at 3 in the client; then show a terminal message.
  5. Never auto-retry on 409 SOLD_OUT — that's a terminal state.
```

### 9.3 Cache Stampede at Sale Start

```
Scenario: The sale config cache entry expires at exactly T=0 for all
          40 pods simultaneously. All 40 hit PostgreSQL at once, plus
          every subsequent request until the cache refills.

Mitigations:
  1. NEVER expire the sale config during a sale. Load it at pre-warm
     and invalidate only via explicit pub/sub push.
  2. Probabilistic early expiration (XFetch): refresh when
       now - delta * beta * log(random()) >= expiry
     so refreshes are staggered across instances.
  3. Single-flight / request coalescing: only one thread per pod
     fetches; others wait on the same future.
       Caffeine's LoadingCache does this for you by default.
  4. Serve stale on error: if the DB is unreachable, keep serving the
     last known config rather than failing the request.
```

```java
// Caffeine gives you single-flight + refresh-ahead almost for free
Caffeine.newBuilder()
    .maximumSize(1_000)
    .refreshAfterWrite(Duration.ofSeconds(30))   // async refresh, serve stale
    .expireAfterWrite(Duration.ofMinutes(10))    // hard expiry, well beyond sale
    .build(saleId -> saleRepository.findById(saleId));
```

### 9.4 Kafka Unavailable

```
Scenario: The Kafka cluster is unreachable at T+1s. Reservations
          succeed in Redis, but no order rows are being created.

Impact: Users see "You got it!" but have no order to pay for.

Mitigation:
  1. Producer buffers in memory (buffer.memory=64MB, non-blocking send
     with max.block.ms=0). We ride out a short outage transparently.
  2. If the buffer fills → fall back to the local outbox: write the
     event to a local append-only file / local PostgreSQL outbox table,
     and drain it when Kafka recovers.
  3. The 10-minute reservation TTL is our recovery budget — enormous
     compared to typical Kafka recovery times.
  4. If Kafka is down for > 5 minutes: stop the sale via the kill
     switch. Better to pause than to hand out reservations that can
     never become orders.
  5. Circuit breaker on the producer prevents each request from paying
     a timeout penalty on the hot path.
```

### 9.5 Payment Failure After Reservation

```
Scenario: User wins a unit, starts payment, payment fails or times out.

Flow:
  T+0:     Reservation created, TTL 10 minutes
  T+2min:  User submits payment → payment gateway declines
  T+2min:  Order status = PAYMENT_FAILED

Decision: do we release the unit immediately or hold it?
  → HOLD until the TTL. Give the user a chance to retry with another
    card. Releasing instantly creates a terrible experience and
    encourages hostile "card-testing" behaviour to churn stock.

  T+10min: Reservation expires
           → Sweeper INCRs Redis stock, removes user from dedup set
           → Publishes RESTOCKED → edge SOLD_OUT flag cleared
           → Optionally notifies waitlisted users ("A unit is available!")

Edge case — payment succeeds AFTER the reservation expires:
  The payment webhook arrives at T+10min30s. Two options:
    a) Refund automatically and apologize (safest, default)
    b) Try to re-acquire a unit; if available, honour the order
  Choose (a) by default. It is deterministic and never oversells.
  Make the payment window (8 min) strictly shorter than the
  reservation TTL (10 min) so this race is rare by construction.
```

### 9.6 The Clock Skew Problem

```
Scenario: App servers disagree about when the sale starts. Server A
          thinks it's T=0; Server B thinks it's T-3s and rejects.

Impact: Users randomly get NOT_STARTED errors, retry hard, and
        perceive the sale as broken.

Mitigation:
  1. Never rely on server wall clock for the go-live decision. Use a
     single source of truth: a Redis flag flipped by ONE scheduled job.
       SET flashsale:{sale}:status LIVE
     App instances subscribe to the pub/sub event.
  2. NTP/chrony on all hosts with a monitored max drift alert (< 50 ms).
  3. Serve `server_time` to clients so their countdown is synced.
  4. Grace window: accept requests from T-2s and hold them briefly
     rather than rejecting, so slightly-early clients aren't punished.
```

### 9.7 The "Everyone Wins Nothing" Problem

```
Scenario: The sale opens and closes in 800 ms. 9,990,000 users see
          SOLD_OUT before they could even tap. Rage on social media,
          and accusations that the sale was fake.

This is a PRODUCT failure, not a technical one, but be ready:
  1. Publish honest stock numbers before the sale ("10,000 units").
  2. Show a live counter during the sale so it's visibly real.
  3. Publish a post-sale transparency report: units sold, unique
     winners, bot attempts blocked.
  4. Consider a LOTTERY model instead of FCFS (see §11.1) — everyone
     who registers within a window has an equal chance. Removes the
     latency advantage of bots and of users with better internet.
  5. Waitlist + restock notification for expired reservations.
```

### 9.8 Failure Mode Summary

| Failure | Direction of error | Severity | Recovery |
|---|---|---|---|
| Redis replica lag on failover | **Oversell risk** | 🔴 Critical | `WAIT` + immediate reconcile + downward correction |
| Redis full data loss | Undersell | 🟡 Medium | Rebuild counters from PostgreSQL |
| Kafka down | Undersell | 🟡 Medium | Buffer → local outbox → drain |
| Order Service down | Delayed orders only | 🟢 Low | Kafka buffers; TTL gives 10 min |
| Sweeper down | Undersell | 🟡 Medium | Units stay locked until it recovers |
| App pod crash | Individual request failures | 🟢 Low | Stateless; client retries |
| PostgreSQL down | Sale halts | 🟠 High | Kill switch; Redis holds correct state |
| Clock skew | Confusing errors | 🟡 Medium | Central status flag, not local clocks |

> **The pattern to articulate:** almost every failure mode lands in the undersell column. The one that doesn't — Redis replica lag — gets three independent mitigations because it's the only path to the outcome the business cannot accept.

---

## 10. Monitoring & Observability

### 10.1 Live Dashboard

```
┌──────────────────────────────────────────────────────────────────────┐
│                    FLASH SALE WAR ROOM — sale_bbd_pixel9              │
├──────────────────┬───────────────────┬───────────────────────────────┤
│  Stock Remaining │  Requests/sec     │  Reservations/sec             │
│  ███░░░░░░  2,340│  ████████ 847k    │  ██░░░░░░  1,240              │
├──────────────────┴───────────────────┴───────────────────────────────┤
│                                                                       │
│  Request outcomes (last 10s)          Latency (p50 / p95 / p99)      │
│    Reserved       10,240   0.1%       Buy path:   8ms / 34ms / 96ms  │
│    Sold out    6,120,400  72.1%       Redis Lua:  0.4ms/ 1.2ms/ 4ms  │
│    Rate limited1,890,300  22.3%       Edge:       2ms /  6ms / 18ms  │
│    Shed (503)    420,100   4.9%                                       │
│    Errors (5xx)      840   0.01% ✅                                   │
│                                                                       │
│  Queue depth: 2,340,120 waiting   Admission rate: 48,000/s           │
│  Est. wait for new arrival: 49s                                       │
├───────────────────────────────────────────────────────────────────────┤
│  CDN hit ratio:    99.7% ✅       Redis CPU (max node):  61% ✅       │
│  Kafka lag:        2,100 msgs     Bot blocks:  1.2M (14% of traffic)  │
│  INVENTORY DRIFT:  0 ✅           RECONCILIATION: ✅ CLEAN (12s ago)  │
└───────────────────────────────────────────────────────────────────────┘
```

### 10.2 Critical Alerts

| Alert | Condition | Severity | Action |
|---|---|---|---|
| Inventory drift negative | `redis + reserved + sold > total_stock` | **P0** | Auto-pause sale, page immediately |
| Oversell detected | Any DB `no_oversell` constraint violation | **P0** | Halt sale, freeze fulfilment, exec escalation |
| Redis node down | Cluster reports failed primary | **P0** | Verify failover, force reconcile |
| Buy path error rate | 5xx > 0.5% over 30 s | **P1** | Check Redis/Kafka health, consider shedding harder |
| Buy path p99 latency | > 500 ms over 1 min | **P1** | Reduce admission rate, add pods |
| Redis CPU | Any node > 80% | **P1** | Increase shard count, verify hot key spread |
| Kafka consumer lag | > 50,000 messages | **P2** | Scale consumers; TTL budget still healthy |
| Sweeper stalled | No run in 60 s | **P2** | Restart; risk is undersell only |
| CDN hit ratio | < 95% during sale | **P2** | Check cache headers, purge misconfiguration |
| Bot traffic share | > 40% of requests | **P2** | Tighten attestation, enable proof-of-work |
| Core commerce p99 | Any regression during sale | **P1** | Bulkhead breach — kill switch is on the table |

### 10.3 Structured Event Logging

```java
// One structured event per purchase attempt — sampled at 1% for
// successes, 100% for errors and for all reservations.
{
  "timestamp": "2026-08-15T12:00:00.412Z",
  "event_type": "flashsale.purchase_attempt",
  "sale_id": "sale_bbd_pixel9",
  "user_id_hash": "u_8f3a...",          // hashed — never log raw PII at this volume
  "trace_id": "4bf92f3577b34da6",       // W3C traceparent, propagated from edge
  "outcome": "RESERVED",                 // RESERVED|SOLD_OUT|LIMIT|SHED|RATE_LIMITED
  "reservation_id": "resv_7x8y9z",
  "shard_attempted": 7,
  "shard_hops": 1,                       // how many shards before success
  "latency_ms": { "total": 12, "redis": 0.6, "kafka_enqueue": 0.2 },
  "admission": { "queue_position": 4211, "waited_ms": 3400 },
  "bot_signals": { "score": 0.02, "device_attested": true, "asn_type": "residential" }
}
```

> **Volume warning:** at 1M RPS, logging every request at 500 bytes is 500 MB/s — 1.8 TB/hour. Sample aggressively: 100% of reservations (only 10k), 100% of errors, 0.1% of sold-out rejections. Use metrics (counters/histograms) for the bulk, logs for the exceptions.

### 10.4 Distributed Tracing

```
Trace a single winning request end to end:

  [Edge Worker]────────────────────────── 2ms
     └─[API Gateway: auth + ratelimit]─── 3ms
         └─[Admission: token check]────── 1ms
             └─[FlashSaleService]──────── 12ms
                 ├─[Redis: reserve.lua]── 0.6ms   ← the critical section
                 └─[Kafka: send async]─── 0.2ms
  ══════ async boundary ══════
  [OrderService: consume]───────────────── 45ms  (T+180ms)
     └─[PostgreSQL: insert order]──────── 8ms
  [NotificationService]────────────────── 22ms  (T+210ms)

Sample at 100% for reservations (only 10k traces — cheap and
invaluable for post-mortems) and 0.01% for rejections.
```

---

## 11. Advanced Features

### 11.1 Lottery Mode (Fairness Alternative to FCFS)

```
Instead of first-come-first-served, run a registration window:

  Phase 1 (24h): Registration
    Users register interest. No stock check, no contention.
    Just an INSERT — trivially scalable, no spike at all.

  Phase 2 (instant): Draw
    Deterministic, verifiable, seeded shuffle:
      seed    = SHA256(sale_id || public_random_beacon)
      ranking = sort(registrants, key = HMAC(seed, user_id))
      winners = ranking[0 : stock_count]

    Publish the seed after the draw so anyone can verify the result.

  Phase 3 (48h): Winners claim and pay.
    Unclaimed units cascade to the next names in the ranking.

Advantages:
  ✅ Eliminates the traffic spike entirely — the hardest problem vanishes
  ✅ Perfectly fair: network latency and device speed stop mattering
  ✅ Bots gain nothing from speed (only from having many accounts)
  ✅ Publicly verifiable — huge for trust

Disadvantages:
  ❌ Loses the urgency and buzz that make flash sales a marketing event
  ❌ Longer sale cycle; lower conversion from winners

Xiaomi and several sneaker platforms moved to this model for exactly
these reasons. Bring it up as a design alternative — interviewers
love a candidate who questions the premise.
```

### 11.2 Waitlist & Restock Notification

```
When SOLD_OUT, offer a waitlist instead of a dead end:

  User taps "Notify me" → ZADD flashsale:{sale}:waitlist <timestamp> <user_id>

When the sweeper releases N units:
  1. ZPOPMIN N*3 from the waitlist (oversubscribe 3x — most won't act)
  2. Push notification with a SHORT exclusive window (e.g. 5 minutes)
  3. Issue single-use claim tokens bound to those users
  4. Unclaimed after the window → back to open sale

This converts undersell into sales and turns a frustrating experience
into a second chance. Typically recovers 60-80% of expired units.
```

### 11.3 Tiered / Progressive Release

```
Instead of releasing 10,000 units at once, release in waves:

  12:00:00 → 3,000 units
  12:05:00 → 3,000 units
  12:10:00 → 4,000 units

Benefits:
  • Peak load is 1/3 — cheaper infrastructure
  • Multiple chances for users who missed the first wave
  • Bots must sustain their attack, increasing their cost
  • Extends the marketing window ("more stock at 12:05!")

Implementation: a scheduled job INCRBYs the Redis shards at each
wave time and publishes a RESTOCKED event to clear the edge flag.
```

### 11.4 Graceful Degradation Ladder

```
As load increases, shed capability before shedding availability:

  Level 0 (normal):     Everything on. Exact stock counts. Live ticker.
  Level 1 (elevated):   Stock shown as buckets. Ticker updates every 5s.
  Level 2 (high):       Waiting room enabled. Recommendations disabled.
  Level 3 (critical):   Stochastic shedding (accept 1 in N). Static
                        page only, no personalization.
  Level 4 (emergency):  Sale paused. "Back in a moment" page from CDN.
                        Core commerce completely unaffected.

Each level is a single feature flag, pre-tested, flippable in < 5s.
The on-call runbook says exactly which metric triggers which level.
```

### 11.5 Cross-Reference: Payment Integration

```
The flash sale hands off to the payment system at the reservation step.
Key integration contracts (see doc 007):

  • Idempotency key: "resv_<reservation_id>" — deterministic, so a
    double-submit from the checkout page cannot double-charge.
  • Payment timeout MUST be shorter than the reservation TTL:
        payment window = 8 min  <  reservation TTL = 10 min
    This eliminates the "paid after expiry" race by construction.
  • On payment success → reservation CONFIRMED → order CONFIRMED →
    the Redis reservation record's status flips so the sweeper skips it.
  • On payment failure → hold the reservation until TTL; allow retry.
  • Ledger entries are written by the payment system, not here. The
    flash sale system owns INVENTORY; the payment system owns MONEY.
    Never let one service own both — that's how you get bugs where
    inventory and ledger disagree.
```

---

## 12. Interview Q&A

### Q1: Why not just use a database transaction with `SELECT ... FOR UPDATE`?

**Candidate:** "It's correct but it doesn't scale, and at flash-sale volume 'doesn't scale' becomes a correctness problem too.

`SELECT ... FOR UPDATE` takes an exclusive row lock. Every one of a million concurrent requests queues behind that single row. The lock is held for the duration of the transaction — the update, the WAL flush, the commit — call it 1-2 ms. That caps you at roughly 600-1,000 reservations per second.

Three things then go wrong:
1. The connection pool is exhausted instantly, so requests fail before they even reach the lock.
2. Lock waits exceed statement timeouts, transactions abort, clients retry, and the queue grows.
3. Because the database is saturated, *every other query in the marketplace* slows down — you've taken down browsing and checkout to sell one phone.

And the timeouts themselves create ambiguity: did the reservation commit before the client gave up? That uncertainty is exactly how oversell bugs get introduced.

Redis with a Lua script gives us the same atomicity — Redis executes commands single-threaded, so a script is inherently serialized — at roughly 100,000 ops per second per node, and we shard the counter across 20 keys to get past a million. That's a thousandfold improvement, and it keeps PostgreSQL completely out of the hot path."

---

### Q2: What if Redis crashes and you lose the inventory counter?

**Candidate:** "First, I'd point out the asymmetry that makes this survivable: Redis losing data causes *undersell*, not *oversell*, which is the direction the business tolerates.

Concretely, the recovery path is:
1. **AOF persistence** with `appendfsync=everysec` bounds worst-case loss to about one second of writes.
2. **Replicas plus the `WAIT` command** — I issue `WAIT 1 50` after the reservation script, so it blocks until at least one replica has the write, with a 50 ms cap. That closes most of the failover window for a few milliseconds of latency.
3. **Rebuild from PostgreSQL.** Reservations are mirrored to PostgreSQL asynchronously, so I can always recompute:
   `redis_stock = total_stock − db_reserved − db_sold − safety_margin`
   The safety margin — a handful of units — guarantees I correct downward. If I'm wrong, I undersell a few units rather than overselling.
4. **Immediate reconciliation on a cluster failover event**, not on the 30-second schedule.
5. If drift is beyond threshold, **pause the sale**, correct, and resume. A ten-second pause is far cheaper than cancelling four hundred orders.

There's a stronger variant for very high-value SKUs: write the reservation to PostgreSQL synchronously before returning 202. That costs about 20 ms and gives a durable guarantee. I'd apply it selectively — worth it for a ₹2 lakh laptop, not for a ₹499 accessory."

---

### Q3: How do you prevent one user from buying the entire stock?

**Candidate:** "Four independent layers, and the important thing is that the first one is *inside the same atomic operation* as the stock decrement.

1. **Atomic dedup in the Lua script.** The first thing the script does is `SADD` the user id — or `SETBIT` on a bitmap for memory efficiency. If it returns 0, the user already has a reservation and we return immediately. Because this happens inside the same script as the decrement, two parallel requests from the same user cannot both succeed. If the stock check subsequently fails, I `SREM` to roll back, so the user isn't wrongly blocked from a later restock.

2. **Database `UNIQUE (sale_id, user_id)`** on the reservations table. Durable enforcement of the same rule, and it catches anything that slips through when Redis state is rebuilt.

3. **Identity-level limits, not account-level.** Limits keyed on a verified phone number, payment instrument, and shipping address — because a scalper's real move isn't 100 requests from one account, it's one request each from 100 accounts. Post-sale graph analysis on shared addresses and cards catches the coordinated ones.

4. **Rate limiting** at the gateway per user, per IP, and per device fingerprint, which caps the blast radius of any single actor.

I'd add that layer 3 is where the real fight is. Layers 1 and 2 are easy and solve the naive case completely; the hard problem is adversarial and partly non-technical."

---

### Q4: How do you handle the hot key problem in Redis?

**Candidate:** "Redis Cluster shards by key hash, so a single inventory key lives on exactly one node. Adding nodes doesn't help — that one node takes all million requests and pins its CPU.

The fix is **inventory sharding**: split 10,000 units into 20 counters of 500 each. Each request picks a shard at random, so load spreads across the cluster and throughput scales roughly linearly.

That introduces one edge case worth calling out: **shard exhaustion.** Shard 7 hits zero while shard 12 still has 300 units. A naive implementation reports SOLD_OUT with stock unsold. So on an empty shard I fall through to the next one and only report SOLD_OUT after all twenty are empty. The cost is up to twenty Redis calls for the very last units — which is fine, because by then traffic has collapsed. An optimization is a bitmap of non-empty shards so later requests skip the dead ones.

Two other things matter here. First, **hash tags**: all keys in a multi-key Lua script must hash to the same slot, so I name them `flashsale:{sale_id}:stock:7`. The braces tell Redis Cluster to hash only `sale_id`. Forgetting this is the classic bug when moving from standalone Redis to a cluster — it works in dev and throws CROSSSLOT in production.

Second, the **read side has the same problem** and people forget it. If every page poll does a `GET` on the stock key, you've recreated the hot key on reads. So a background job reads the summed stock once per second and pushes it to the CDN edge and to each pod's local cache. Ten million reads per second become one."

---

### Q5: A user's payment succeeds but the reservation has already expired. What happens?

**Candidate:** "First, I'd design so this is rare rather than handled: I make the payment window strictly shorter than the reservation TTL — eight minutes against ten — so the race only occurs on a genuinely delayed gateway callback.

When it does happen, I default to **automatic refund with an apology and a discount code.** It's deterministic and it can never oversell. The alternative — trying to re-acquire a unit — is tempting but it reintroduces exactly the race condition I've spent the whole design eliminating, and it makes the outcome dependent on timing the user can't see.

That said, there's a middle path worth mentioning: if the sale still has stock at the moment the late callback arrives, attempt one atomic reservation for that user. If it succeeds, honour the order; if not, refund. It's strictly better for the user and still safe, because it goes through the same Lua script as everyone else. I'd only build it if the business asks, because it adds a code path that's hard to test.

Operationally, the key is that the payment system and the inventory system have separate sources of truth, and the reservation record is the join point. The payment webhook looks up the reservation, sees `EXPIRED`, and triggers the refund flow. That state check has to be atomic — a compare-and-set on the reservation status — otherwise you can double-process a duplicated webhook."

---

### Q6: How would you make the sale fair? Bots always win.

**Candidate:** "I'd split this into what technology can do and what it can't.

**What layered defence buys you:** identity friction before the sale — verified phone, account age, purchase history — so an account isn't free. Client attestation with Play Integrity and App Attest, plus signed request payloads, so a raw curl fails. Behavioural signals — a click 50 ms after page load with no mouse movement, or a datacenter ASN — scored in under five milliseconds. And a proof-of-work challenge that costs a human 200 invisible milliseconds but multiplies a bot farm's cost by the number of accounts it runs.

One tactical detail I'd emphasize: **shadow rejection.** Never tell a bot it was blocked. Return an ordinary-looking SOLD_OUT. If you return 'BOT_DETECTED', the operator immediately knows which signal fired and iterates within the hour. Silent failure is far more expensive for them to debug.

**What technology can't fix:** a determined operator with ten thousand verified SIMs and residential proxies will get through. So the honest framing is that this is an economics problem — the goal is to push the cost per successful purchase above the resale margin, not to achieve perfect exclusion.

Which is why I'd genuinely propose the **lottery model** as an alternative. A registration window followed by a verifiable seeded draw removes speed as an advantage entirely — network latency and device performance stop mattering, and bots gain nothing from being fast, only from having many accounts, which the identity layer already attacks. It also eliminates the traffic spike, which is the hardest engineering problem in the whole design. The cost is losing the urgency that makes a flash sale a marketing event, so it's a product decision, not an engineering one. Xiaomi moved this way for exactly these reasons."

---

### Q7: How do you display real-time stock to 10 million users without melting Redis?

**Candidate:** "I don't read Redis on the request path at all.

A background publisher reads the summed stock **once per second** and fans it out three ways: to the CDN edge as a small cached JSON with a one-second TTL, over Redis pub/sub to each pod's in-process cache, and over WebSocket or SSE to connected clients for the live ticker. Ten million reads per second collapse to one.

Two refinements. First, I return a **bucket** — HIGH, MEDIUM, LOW, SOLD_OUT — rather than an exact number. It caches better because there are fewer distinct values, it stops scrapers inferring exact sale velocity, and no user can act on the difference between 240 and 237 units.

Second, and this is the highest-leverage optimization in the entire design: the moment stock hits zero, I propagate a **SOLD_OUT flag outward** as fast as possible. Redis pub/sub sets a local boolean on every pod within about five milliseconds, so subsequent requests are rejected in-process with zero Redis calls. Then a CDN edge KV update within a few hundred milliseconds means requests are rejected at the edge with zero origin calls. And a WebSocket push turns the button grey so most users never send the request at all.

After roughly half a second, effectively all remaining traffic — which could be hundreds of millions of requests over the next minute — is absorbed by the CDN at near-zero cost. And I propagate the flag in reverse when the sweeper releases expired reservations, so we don't undersell."

---

### Q8: How do you guarantee the flash sale doesn't take down the rest of the site?

**Candidate:** "Bulkheading, taken seriously — meaning physical separation, not just a code-level abstraction.

The flash sale gets its own Kubernetes node pool, its own Redis cluster, its own Kafka cluster, its own CDN hostname, and at minimum its own database connection pool against a separate replica. Ideally its own database entirely. The point is that resource exhaustion in the sale path cannot starve core commerce, because they don't share the resource.

On top of that: Resilience4j bulkheads give separate thread pools per downstream call, so a slow dependency can't consume every thread. Circuit breakers on every cross-service call — for instance, the sale service calling user-profile — fail fast rather than queueing. Rate limits on flash-sale-to-core-service calls cap the load the sale can impose even when it's healthy.

And a **kill switch** that disables all flash sales in under five seconds, tested before every event.

The part people skip is verification. I'd run a game day where we deliberately overload the flash sale path and assert that core checkout p99 is statistically unchanged. If you haven't tested the bulkhead under load, you don't have a bulkhead — you have a diagram."

---

### Q9: Walk me through what happens in the first 100 milliseconds of the sale.

**Candidate:** "Let me trace it.

**T=0ms.** A single scheduled job flips a Redis key to LIVE and publishes on pub/sub. I use one central flag rather than each server checking its own clock, because clock skew across forty pods produces random NOT_STARTED errors and a retry storm.

**T=5ms.** Every pod has received the pub/sub message and updated its local cache. Simultaneously, the CDN edge config flips and WebSocket clients get a push that enables the buy button.

**T=10 to 50ms.** The first wave arrives. Requests hit the edge, pass gateway auth and rate limiting, then admission control. Admission is a token bucket refilling at the rate the inventory layer can actually serve — say fifty thousand per second. Everyone above that gets a queue position, not an error.

**T=15 to 100ms.** Admitted requests execute the Lua script. Each one does dedup, stock check, decrement, reservation write — atomically, in about half a millisecond. Winners get a 202 immediately and a Kafka message goes out asynchronously; the response does not wait for it.

**Around T=200ms**, in this scenario, stock is exhausted. The script returns zero, we publish SOLD_OUT, local flags trip within five milliseconds, the edge flag within a few hundred, and from that point virtually all traffic is absorbed at the CDN.

**T=200ms onward**, asynchronously: Kafka consumers create ten thousand PENDING_PAYMENT order rows in PostgreSQL — a trivial write load spread over seconds — and notifications go out.

The thing I'd highlight is that the entire critical section is one Redis script per request. There is no database, no lock, no join, and no synchronous network call to another service on the path that determines who wins."

---

### Q10: How is this different from designing a ticketing system like Ticketmaster?

**Candidate:** "The traffic shape is similar, but three things change the design materially.

**Seats are not fungible.** Flash sale inventory is a counter — any unit is as good as any other, which is exactly why a Redis `DECRBY` works. Seats are distinct items with adjacency constraints, so you need per-seat locking, a seat map, and 'find me three together' queries. That's closer to a spatial reservation problem than a counter problem, and Redis alone won't carry it.

**The hold semantics are richer.** Ticketing users browse a seat map, hold specific seats while they decide, and often change their minds. That means longer holds, hold-then-swap flows, and a much bigger surface for abandoned holds — so the release mechanism and its correctness matter more than they do here.

**Fairness expectations are legal, not just reputational.** Several jurisdictions regulate ticket resale and require certain fairness disclosures, so the waiting room and bot mitigation aren't merely product choices.

What carries over cleanly: the queue-based admission control, the static-first page architecture, the async order creation, the CDN-level rejection of doomed traffic, and the principle of failing toward undersell. If I were designing ticketing, I'd keep this entire outer architecture and replace the Redis counter with a per-seat state machine backed by short-lived distributed locks — probably Redis with fencing tokens, or a purpose-built service with a compacted Kafka log per event."

---

## 13. Follow-Up Questions (Rapid Fire)

> These are the questions interviewers use to probe depth once the main design is on the board. Prepare a 30-60 second answer for each.

### 13.1 Inventory & Consistency

| # | Question | Key points to hit |
|---|---|---|
| 1 | Why is Lua atomic in Redis? | Redis executes commands on a single thread; a script runs to completion with no interleaving. Caveat: it blocks the server, so keep scripts short — a slow script stalls every client. |
| 2 | What if the Lua script itself throws an error halfway? | Redis does **not** roll back partial script effects (unlike a SQL transaction). Design scripts so any partial state is recoverable, do the mutating calls last, and roll back explicitly (the `SREM` in our script). |
| 3 | Could you use `DECR` and check for negative instead of checking first? | Yes — `DECRBY` then `if result < 0 then INCRBY back`. It's one fewer read. But it briefly allows negative stock, so any concurrent reader must clamp at zero. Slightly faster, slightly less obvious. |
| 4 | How do you handle quantity > 1 per user? | Same script, `DECRBY qty`, but shard fallback becomes harder — a shard may have 2 units when the user wants 3. Either allow partial fulfilment, or reserve across shards inside one script (all keys must share a hash tag). |
| 5 | What's your undersell rate and is it acceptable? | Expect < 0.1% — from expired-but-unreleased reservations and shard fallback races. Quantify it, monitor it, and confirm the business tolerance up front. |
| 6 | How do you restock mid-sale? | `INCRBY` on shard 0, write an `inventory_events` row with source=ADMIN, publish RESTOCKED to clear the edge SOLD_OUT flag. Must be idempotent — use a restock id. |
| 7 | What isolation level would you use if you *did* use PostgreSQL? | `READ COMMITTED` is sufficient because `UPDATE ... WHERE stock > 0` takes a row lock and re-evaluates the predicate on the fresh row. `SERIALIZABLE` would add serialization failures and retries for no benefit here. |
| 8 | Optimistic vs pessimistic locking for the DB path? | Optimistic (version column / conditional UPDATE) is better under *low* contention. Under flash-sale contention almost every attempt fails and retries, so pessimistic wins — but both are too slow, which is why we're in Redis. |

### 13.2 Scaling & Performance

| # | Question | Key points to hit |
|---|---|---|
| 9 | How many shards should you pick, and why not 1,000? | ≈ 3-4× the number of Redis nodes. Too many shards makes the tail worse (more fallback hops) and buys throughput the cluster can't deliver anyway. |
| 10 | What is the actual bottleneck at 1M RPS? | Usually TLS handshakes and connection setup at the edge, not application logic. Then Redis single-node CPU. Then Kafka producer buffer. Measure — don't assume it's the database. |
| 11 | How do you load test this? | Distributed load generators (k6/Gatling/Locust) across many regions, ramping 0→1M RPS in under a second. Test the *spike shape*, not just the steady rate — a ramp-up test hides thundering-herd bugs. |
| 12 | How does Kubernetes HPA handle a 0→1M step function? | It doesn't — HPA reacts on a 15-30s metrics window. Pre-scale explicitly via a scheduled KEDA trigger or a manual replica count set at T-30min. |
| 13 | Why does JVM warmup matter here? | Cold JVM interprets bytecode; C2 compiles after ~10k invocations. A cold pod can be 10-20× slower for the first seconds — precisely the seconds that matter. Send synthetic traffic at T-15min, or use AOT/CRaC. |
| 14 | How do you avoid GC pauses on the hot path? | Minimize allocation per request (reuse buffers, avoid boxing), use ZGC or Shenandoah for sub-millisecond pauses, size the heap so young-gen collections are rare. A 200 ms stop-the-world pause at T=0 is thousands of failed requests. |
| 15 | Why Kafka and not RabbitMQ or SQS? | Throughput and replay. We need to absorb a 10k-message burst and, more importantly, replay from an offset if the consumer had a bug. Kafka's log retention is the recovery mechanism. |
| 16 | Why `acks=1` and not `acks=all` on the producer? | Latency on the hot path. The trade-off is safe *only because* a lost event causes undersell (reservation expires, stock returns), never oversell. Always justify a durability downgrade by naming the failure direction. |

### 13.3 Failure & Recovery

| # | Question | Key points to hit |
|---|---|---|
| 17 | Redis primary fails at the exact moment of the last unit — walk through it. | Replica may lag → may show stock that's gone → oversell risk. `WAIT 1 50` shrinks the window; failover triggers immediate reconciliation; correction is always downward with a safety margin. |
| 18 | What if the sweeper double-releases a reservation? | It can't: the status check, `INCRBY`, and `DEL` are in one Lua script, so a second sweeper sees the key gone. `@SchedulerLock` is a second layer. |
| 19 | What if a user's reservation exists in Redis but the order row never got created? | The reservation TTL expires, stock returns, and the user sees "expired" on the status page. We monitor `reserved_in_redis − orders_in_pg` as the Kafka-lag health signal. |
| 20 | How would you roll back a bad deploy mid-sale? | You don't deploy mid-sale — freeze at T-5min. If you must: feature-flag rollback first (instant), pod rollback second. Never a schema change. |
| 21 | What's your kill switch and how fast is it? | A single Redis flag + edge config change; pods react on pub/sub in ~5ms, edge in < 5s. Tested before every event as part of the checklist. |
| 22 | Region-wide outage 30 seconds before the sale? | Postpone. Announce. Do not fail over an inventory authority mid-flight — a split-brain counter is the one scenario that produces guaranteed oversell. |
| 23 | How do you test failure modes before the sale? | Chaos engineering in a staging environment at production scale: kill Redis primaries, add network latency, saturate Kafka, and assert the drift metric stays at zero. |

### 13.4 Product & Fairness

| # | Question | Key points to hit |
|---|---|---|
| 24 | Is FCFS actually fair? | No — it rewards fast internet, fast devices, and geographic proximity to the datacenter. A user in a tier-3 city with 4G loses to a bot on a Mumbai VPS every time. That's the argument for a lottery. |
| 25 | Should the waiting room position be strictly ordered? | Not necessarily. Strict ordering rewards whoever loaded the page first, which is itself a proxy for connection quality. Randomizing within a time bucket is arguably fairer and equally cheap. |
| 26 | How do you handle a user with two devices? | Dedup on user identity, not session or device. But allow a legitimate retry from a second device — that's why the reservation lookup returns the *existing* reservation rather than a hard error. |
| 27 | What do you show a user who lost? | Never a raw error. Show the waitlist option, related in-stock products, and the next drop time. Converting a loss into a waitlist signup is real business value. |
| 28 | How do you prove the sale wasn't rigged? | Publish stock counts beforehand, a live counter during, and a post-sale transparency report. For lottery mode, publish the seed so the draw is independently verifiable. |

### 13.5 Design Alternatives (Be Ready to Defend Your Choice)

| # | Question | Key points to hit |
|---|---|---|
| 29 | Why not a fully queue-based design — every request goes to Kafka, one consumer decrements? | It works and is very simple to reason about (perfect FCFS, single writer). But the user must poll for a result instead of getting a synchronous answer, and the single consumer caps you at its throughput. Good for extreme fairness requirements; worse UX. |
| 30 | Why not put the counter in the CDN edge (Cloudflare Durable Objects / Workers KV)? | Durable Objects give you a genuinely serialized single-threaded counter at the edge, which is elegant. Downsides: vendor lock-in, harder to reconcile with your own DB, and a cold-start penalty. Worth mentioning as a modern alternative. |
| 31 | Could you use a distributed lock (Redlock) instead? | You could, but you don't need to — a Lua script on a single key gives you the same mutual exclusion with none of Redlock's contested safety assumptions. Reaching for a distributed lock when an atomic primitive exists is a red flag. |
| 32 | Would an in-memory counter in the app with a leader work? | Yes — a single leader pod holding the counter in memory is the fastest possible option. But leader election, failover, and state recovery reintroduce every problem Redis already solved for you. Only worth it above ~5M RPS. |
| 33 | How would you adapt this for a 1-unit sale (a single rare item)? | Sharding is meaningless. Go the other way: a single key, strict FCFS, synchronous PostgreSQL write for durability, and heavy admission control. With one unit, latency doesn't matter — correctness and auditability do. |
| 34 | How does this change if inventory is per-warehouse (regional stock)? | The counter becomes per (sku, warehouse), so contention naturally spreads. But now you need to map the user to a warehouse *before* reserving, and handle "sold out near you but available elsewhere," which is a product decision about delivery time vs availability. |

---

## 14. Production Checklist

### Pre-Launch (Week -2 to -1)

| # | Item | Status |
|---|---|---|
| 1 | Load test at 2x expected peak with a true step-function ramp | ☐ |
| 2 | Chaos test: kill Redis primary mid-sale, verify zero oversell | ☐ |
| 3 | Bulkhead verification: core commerce p99 unchanged under sale load | ☐ |
| 4 | Kill switch tested end to end (< 5 s to full stop) | ☐ |
| 5 | Bot mitigation tuned against a red-team script | ☐ |
| 6 | Reconciler running green for 7 days on staging sales | ☐ |
| 7 | Runbooks written for P0s (drift, oversell, Redis failover) | ☐ |
| 8 | CDN cache rules verified — hit ratio > 99% on the static shell | ☐ |
| 9 | Degradation ladder flags tested at each level | ☐ |
| 10 | Legal/PR sign-off on stock disclosure and fairness messaging | ☐ |

### T-30 Minutes

| # | Item | Status |
|---|---|---|
| 1 | Pods pre-scaled to 2x planned capacity | ☐ |
| 2 | JVM warmed with synthetic traffic (verify JIT compilation counters) | ☐ |
| 3 | Redis inventory shards initialized; `SUM(shards) == total_stock` verified | ☐ |
| 4 | Sale config loaded into every pod's local cache (verify via health endpoint) | ☐ |
| 5 | CDN static assets pushed and edge-verified from 3 regions | ☐ |
| 6 | Deploy freeze in effect; change window closed | ☐ |
| 7 | War-room bridge open; dashboards on screen; on-call confirmed | ☐ |
| 8 | Kill switch access verified for at least 2 engineers | ☐ |

### During Sale (T=0 to sold out)

| # | Item | Status |
|---|---|---|
| 1 | Inventory drift metric watched continuously (must stay at 0) | ☐ |
| 2 | Core commerce p99 watched for any regression | ☐ |
| 3 | Redis max-node CPU below 80% | ☐ |
| 4 | Kafka consumer lag draining, not growing | ☐ |
| 5 | Bot block rate within expected range (not 0%, not 90%) | ☐ |
| 6 | Error rate below 0.5%; escalate the degradation ladder if not | ☐ |

### Post-Sale (T+1 hour to T+7 days)

| # | Item | Status |
|---|---|---|
| 1 | Full reconciliation: units sold + reserved + remaining == total_stock | ☐ |
| 2 | All expired reservations released; final stock returned to catalogue | ☐ |
| 3 | Bot/scalper graph analysis run; fraudulent orders cancelled and restocked | ☐ |
| 4 | Post-mortem on every P1/P2 alert that fired | ☐ |
| 5 | Latency and funnel metrics archived as the baseline for the next sale | ☐ |
| 6 | Transparency report published (units, winners, bots blocked) | ☐ |
| 7 | Waitlist users notified of restocked units | ☐ |
| 8 | Cost review: CDN egress, Redis, compute — tune for next time | ☐ |

---

## Summary

| Aspect | Decision | Rationale |
|---|---|---|
| **Inventory store** | Redis + Lua script | Atomic check-and-decrement at 100k+ ops/s per node |
| **Hot key** | 20 sharded counters with fallback | One key = one node; sharding scales linearly |
| **Durable truth** | PostgreSQL with `no_oversell` CHECK | Last line of defence; fails loudly |
| **Order creation** | Async via Kafka | Keeps the hot path to one Redis call |
| **Admission** | Virtual waiting room (token bucket) | Defers users instead of rejecting; prevents retry storms |
| **Read path** | CDN + 1s-TTL bucketed stock | 10M reads/sec → 1 Redis read/sec |
| **Sold-out** | Flag propagated to edge | Absorbs the long tail at near-zero cost |
| **Reservation** | 10-min TTL, swept back to stock | Bounds the damage of abandoned checkouts |
| **Failure bias** | Every mode fails toward **undersell** | The business tolerates undersell, never oversell |
| **Isolation** | Full bulkhead from core commerce | A sale must never take down the marketplace |
| **Fairness** | Layered bot defence + optional lottery | Raise cost above resale margin; lottery removes speed advantage |

### Scalability Path

```
Phase 1 (Current):  10k units, 1M RPS, single region
  → Redis cluster (6 nodes), 20 shards per SKU
  → 40 app pods, pre-scaled
  → Kafka 50 partitions
  → CDN with edge sold-out flag

Phase 2 (500 concurrent SKUs — full Big Billion Days):
  → Per-SKU shard counts tuned by expected demand
  → Dedicated Redis cluster per SKU tier (headline vs long-tail)
  → Staggered start times to flatten the aggregate peak
  → Regional stock partitioning for warehouse-bound SKUs

Phase 3 (10M RPS — Singles' Day scale):
  → Edge-resident counters (Durable Objects / regional Redis at POPs)
  → Hierarchical inventory: edge holds a slice, refills from origin
  → Full lottery mode for headline SKUs to eliminate the spike
  → ML-driven admission rate tuning in real time

Phase 4 (Platform):
  → Flash Sale as a self-serve product for third-party sellers
  → Per-tenant isolation and quota
  → Configurable fairness policy (FCFS / lottery / tiered)
  → SLA-backed oversell guarantee with automated reconciliation reports
```

---

> **Interview Tip:** Flash sale design has one central insight, and everything else follows from it: **99.9% of requests must fail, so reject them as early and as cheaply as possible.** Lead with the request funnel, then the atomic Redis decrement, then the hot key sharding. If you say nothing else memorable, say this: *"I designed every failure mode to bias toward undersell, because the business tolerates undersell and cannot tolerate oversell."* That single sentence demonstrates you understand which correctness property actually matters.
