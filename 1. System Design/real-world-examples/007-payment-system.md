# Complete System Design: Payment Processing System (Production-Ready)

> **Complexity Level:** Advanced  
> **Estimated Time:** 60-75 minutes in interview  
> **Real-World Examples:** Stripe, PayPal, Square, Razorpay, Adyen

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Payment Processing Engine](#7-deep-dive-payment-processing-engine)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a payment processing system like Stripe that handles credit card payments, refunds, and payouts for an e-commerce platform."

**Interviewer's Perspective:**  
This problem assesses a candidate's ability to reason about:
- **Reliability:** Payments must never be lost or duplicated
- **Idempotency:** Retrying a request must not charge a customer twice
- **Financial data consistency:** Every cent must be accounted for (double-entry bookkeeping)
- **Security:** PCI DSS compliance, tokenization, encryption
- **Fault tolerance:** Graceful handling of PSP timeouts and network partitions
- **Regulatory awareness:** Audit trails, data retention, cross-border compliance

> **Why this is a hard problem:** Unlike most systems where occasional inconsistency is tolerable, payment systems demand *exactly-once* semantics. A single bug can cause real financial loss or regulatory violations.

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before diving into the design, I'd like to clarify the scope and constraints. Payment systems are highly regulated, so understanding the exact requirements is critical."

**Interviewer:** "Go ahead."

### 2.1 Functional Requirements

**Candidate:** "Let me walk through the functional requirements:
1. Can customers pay via credit card, debit card, and bank transfer?
2. Do we need to support full and partial refunds?
3. Are we handling merchant payouts (settling funds to merchant bank accounts)?
4. Do we need real-time payment status tracking?
5. Should we support multiple currencies?
6. Is idempotent payment processing a hard requirement?"

**Interviewer:** "Yes to all. We're building a platform like Stripe — merchants integrate with our API, customers pay through various methods, and we settle funds to merchants on a schedule."

**Candidate:** "Got it. So the core functional requirements are:
1. ✅ Process payments (credit card, debit card, bank transfer)
2. ✅ Handle refunds (full and partial)
3. ✅ Manage merchant payouts on configurable schedules
4. ✅ Real-time payment status tracking via API and webhooks
5. ✅ Multi-currency support with conversion
6. ✅ Idempotent payment processing (retry-safe)"

### 2.2 Non-Functional Requirements

**Candidate:** "For non-functional requirements, I want to confirm:
1. What's the consistency model? I assume strong consistency for financial operations.
2. What's the latency budget for a payment API call?
3. What availability target are we aiming for?
4. Do we need PCI DSS compliance?
5. What are the audit and data retention requirements?"

**Interviewer:** "Strong consistency is non-negotiable. Payment latency under 2 seconds. We need five-nines availability for the payment API. Full PCI DSS compliance and a complete audit trail for at least 7 years."

**Candidate:** "Summarizing non-functional requirements:

| Requirement | Target |
|---|---|
| Consistency | Exactly-once payment processing (no double charges) |
| Latency | < 2 seconds end-to-end payment processing |
| Availability | 99.999% for payment API (~5 min downtime/year) |
| Security | PCI DSS Level 1 compliance |
| Audit | Full audit trail, 7-year retention |
| Encryption | Data encrypted at rest (AES-256) and in transit (TLS 1.3) |

### 2.3 Scale

**Candidate:** "What's the expected scale?"

**Interviewer:** "Plan for 1 million transactions per day, $10 billion annual GMV, supporting 100,000 merchants."

**Candidate:** "Noted:

| Metric | Value |
|---|---|
| Daily transactions | 1,000,000 |
| Annual GMV | $10,000,000,000 |
| Active merchants | 100,000 |
| Payment methods | Credit card, debit card, bank transfer |
| Supported currencies | 50+ |

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me estimate the traffic patterns."

```
Transactions per day:     1,000,000
Transactions per second:  1,000,000 / 86,400 ≈ 12 TPS (average)
Peak TPS (10x average):   ~120 TPS
Black Friday peak (20x):  ~240 TPS

Read operations (status checks, history):
  Assume 10:1 read-to-write ratio
  Read QPS (average):     ~120 QPS
  Read QPS (peak):        ~1,200 QPS

API calls (total):
  Payment creation:       12 TPS
  Status checks:          60 TPS
  Webhooks delivery:      24 TPS (2 per payment avg)
  Refund requests:        ~1 TPS (5-8% refund rate)
```

### 3.2 Storage Estimation

**Candidate:** "For storage, financial data is never deleted — only soft-deleted or archived."

```
Single transaction record:  ~2 KB
  - Payment record:         ~500 bytes
  - Ledger entries (2+):    ~400 bytes
  - Audit log entries:      ~600 bytes
  - Metadata/indexes:       ~500 bytes

Daily storage:              1M × 2 KB = 2 GB/day
Monthly storage:            60 GB/month
Annual storage:             730 GB/year
7-year retention:           ~5.1 TB

Monetary precision:
  CRITICAL — use integer cents (bigint), NEVER floating point
  $99.99 → stored as 9999 cents
  Avoids IEEE 754 rounding errors (e.g., 0.1 + 0.2 ≠ 0.3)
```

### 3.3 Bandwidth Estimation

```
Average payment request:    ~1 KB
Average payment response:   ~500 bytes

Inbound bandwidth:  12 TPS × 1 KB = 12 KB/s (trivial)
Outbound bandwidth: 12 TPS × 500 B = 6 KB/s

Webhook delivery:   24/s × 500 B = 12 KB/s
Total bandwidth:    ~30 KB/s average, ~600 KB/s peak

→ Bandwidth is NOT the bottleneck; database consistency is.
```

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌───────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                      │
│      Merchant Frontend  |  Mobile SDK  |  Server-to-Server API            │
└──────────────────────────────┬────────────────────────────────────────────┘
                               │
                        ┌──────▼──────┐
                        │ API Gateway │  Rate Limiting, Auth, TLS Termination
                        │  (Kong)     │  IP Whitelisting, Request Validation
                        └──────┬──────┘
                               │
              ┌────────────────┼────────────────────┐
              │                │                     │
     ┌────────▼────────┐ ┌────▼──────┐  ┌──────────▼──────────┐
     │  Payment API    │ │ Refund API│  │  Payout API          │
     │  Service        │ │ Service   │  │  Service             │
     └────────┬────────┘ └────┬──────┘  └──────────┬──────────┘
              │                │                     │
              │         ┌──────▼──────────────────────▼──────┐
              │         │         Idempotency Store          │
              │         │         (Redis Cluster)            │
              │         └──────────────┬─────────────────────┘
              │                        │
     ┌────────▼────────────────────────▼────────┐
     │            Payment Orchestrator           │
     │         (Payment State Machine)           │
     └────┬──────────┬──────────┬───────────────┘
          │          │          │
   ┌──────▼───┐ ┌───▼────┐ ┌──▼──────────┐
   │  PSP     │ │ Ledger │ │   Wallet    │
   │  Gateway │ │ Service│ │   Service   │
   │ (Stripe, │ │(Double │ │ (Merchant   │
   │  Adyen)  │ │ Entry) │ │  Balances)  │
   └──────┬───┘ └───┬────┘ └──┬──────────┘
          │          │          │
          ▼          ▼          ▼
   ┌─────────────────────────────────────────────┐
   │              DATA LAYER                      │
   │  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
   │  │PostgreSQL│  │PostgreSQL│  │  Redis     │  │
   │  │(Payments)│  │ (Ledger) │  │ (Cache/   │  │
   │  │  Primary │  │ Primary  │  │  Idempot.)│  │
   │  └────┬─────┘  └────┬─────┘  └───────────┘  │
   │       │              │                        │
   │  ┌────▼─────┐  ┌────▼─────┐                  │
   │  │ Read     │  │ Read     │                  │
   │  │ Replicas │  │ Replicas │                  │
   │  └──────────┘  └──────────┘                  │
   └─────────────────────────────────────────────┘
          │
   ┌──────▼──────────────────────────────────────┐
   │              EVENT BUS (Kafka)               │
   │  Topics: payment.created, payment.completed, │
   │  payment.failed, refund.initiated,           │
   │  payout.scheduled, ledger.entry.created      │
   └──────┬──────────┬───────────────┬───────────┘
          │          │               │
   ┌──────▼───┐ ┌───▼────────┐ ┌───▼───────────┐
   │ Webhook  │ │Reconcilia- │ │  Analytics /  │
   │ Delivery │ │tion Service│ │  Reporting    │
   │ Service  │ │            │ │  (OLAP)       │
   └──────────┘ └────────────┘ └───────────────┘
```

### 4.2 API Design

**Candidate:** "Here are the core API endpoints."

#### Create Payment
```
POST /api/v1/payments
Headers:
  Authorization: Bearer <api_key>
  Idempotency-Key: "txn_abc123_attempt1"
  Content-Type: application/json

Request Body:
{
  "amount": 9999,              // in cents ($99.99)
  "currency": "USD",
  "payment_method": "pm_card_visa_4242",
  "merchant_id": "merch_001",
  "description": "Order #12345",
  "metadata": {
    "order_id": "order_12345",
    "customer_email": "user@example.com"
  },
  "capture": true              // auth + capture in one step
}

Response (201 Created):
{
  "id": "pay_7x8y9z",
  "status": "CAPTURED",
  "amount": 9999,
  "currency": "USD",
  "psp_reference": "ch_3abc123",
  "created_at": "2026-04-24T10:30:00Z"
}
```

#### Create Refund
```
POST /api/v1/refunds
Headers:
  Authorization: Bearer <api_key>
  Idempotency-Key: "refund_pay7x8y9z_1"

Request Body:
{
  "payment_id": "pay_7x8y9z",
  "amount": 5000,              // partial refund: $50.00
  "reason": "customer_request"
}

Response (201 Created):
{
  "id": "ref_abc123",
  "payment_id": "pay_7x8y9z",
  "amount": 5000,
  "status": "PROCESSING",
  "created_at": "2026-04-24T11:00:00Z"
}
```

#### Get Payment Status
```
GET /api/v1/payments/pay_7x8y9z
Headers:
  Authorization: Bearer <api_key>

Response (200 OK):
{
  "id": "pay_7x8y9z",
  "status": "CAPTURED",
  "amount": 9999,
  "amount_refunded": 5000,
  "currency": "USD",
  "payment_method": "pm_card_visa_4242",
  "psp_reference": "ch_3abc123",
  "refunds": [
    { "id": "ref_abc123", "amount": 5000, "status": "COMPLETED" }
  ],
  "created_at": "2026-04-24T10:30:00Z",
  "updated_at": "2026-04-24T11:05:00Z"
}
```

#### Schedule Payout
```
POST /api/v1/payouts
Headers:
  Authorization: Bearer <api_key>

Request Body:
{
  "merchant_id": "merch_001",
  "amount": 950000,            // $9,500.00
  "currency": "USD",
  "destination": "ba_merchant_bank_001"
}
```

### 4.3 Payment Lifecycle Data Flow

```
Customer clicks "Pay"
       │
       ▼
  ┌─────────────┐     ┌──────────────┐
  │ 1. API call  │────▶│ 2. Check     │
  │    received  │     │  idempotency │
  └─────────────┘     └──────┬───────┘
                              │
              ┌───────────────┼──────────────────┐
              │ (key exists)  │ (new key)         │
              ▼               ▼                   │
     Return cached    ┌──────────────┐            │
     response         │ 3. Create    │            │
                      │   payment    │            │
                      │   record     │            │
                      │   (CREATED)  │            │
                      └──────┬───────┘            │
                             │                    │
                      ┌──────▼───────┐            │
                      │ 4. Call PSP  │            │
                      │   (Stripe/   │            │
                      │    Adyen)    │            │
                      └──────┬───────┘            │
                             │                    │
              ┌──────────────┼──────────────┐     │
              │ (success)    │ (failure)     │     │
              ▼              ▼               │     │
     ┌──────────────┐ ┌───────────┐         │     │
     │ 5. Update    │ │ Mark      │         │     │
     │   status =   │ │ FAILED,   │         │     │
     │   CAPTURED   │ │ return    │         │     │
     └──────┬───────┘ │ error     │         │     │
            │          └───────────┘         │     │
     ┌──────▼───────┐                        │     │
     │ 6. Write     │                        │     │
     │   ledger     │                        │     │
     │   entries    │                        │     │
     └──────┬───────┘                        │     │
            │                                │     │
     ┌──────▼───────┐                        │     │
     │ 7. Publish   │                        │     │
     │   event to   │                        │     │
     │   Kafka      │                        │     │
     └──────┬───────┘                        │     │
            │                                │     │
     ┌──────▼───────┐                        │     │
     │ 8. Return    │                        │     │
     │   response   │                        │     │
     │   to client  │                        │     │
     └──────────────┘                        │     │
```

---

## 5. Deep Dive: Core Components

### 5.1 Payment Service (Orchestrator)

**Candidate:** "The Payment Service is the central orchestrator. It coordinates the entire payment flow using a state machine pattern."

```javascript
// Payment Orchestrator - Node.js
class PaymentOrchestrator {
  constructor(idempotencyStore, paymentRepo, pspGateway, ledgerService, eventBus) {
    this.idempotencyStore = idempotencyStore;
    this.paymentRepo = paymentRepo;
    this.pspGateway = pspGateway;
    this.ledgerService = ledgerService;
    this.eventBus = eventBus;
  }

  async processPayment(request, idempotencyKey) {
    // Step 1: Check idempotency
    const cached = await this.idempotencyStore.get(idempotencyKey);
    if (cached) return cached;

    // Step 2: Create payment record in DB (status = CREATED)
    const payment = await this.paymentRepo.create({
      idempotency_key: idempotencyKey,
      merchant_id: request.merchant_id,
      amount_cents: request.amount,
      currency: request.currency,
      status: 'CREATED',
    });

    try {
      // Step 3: Transition to PROCESSING
      await this.paymentRepo.updateStatus(payment.id, 'PROCESSING');

      // Step 4: Call PSP
      const pspResult = await this.pspGateway.charge({
        amount: request.amount,
        currency: request.currency,
        payment_method: request.payment_method,
        idempotency_key: idempotencyKey,
      });

      // Step 5: Transition to CAPTURED
      await this.paymentRepo.updateWithPSP(payment.id, {
        status: 'CAPTURED',
        psp_reference: pspResult.reference,
      });

      // Step 6: Create ledger entries (double-entry)
      await this.ledgerService.recordPayment(payment.id, request);

      // Step 7: Publish event
      await this.eventBus.publish('payment.completed', { payment_id: payment.id });

      // Step 8: Cache response for idempotency
      const response = { id: payment.id, status: 'CAPTURED', ...pspResult };
      await this.idempotencyStore.set(idempotencyKey, response, TTL_24H);

      return response;

    } catch (error) {
      await this.paymentRepo.updateStatus(payment.id, 'FAILED', error.message);
      await this.eventBus.publish('payment.failed', { payment_id: payment.id, error });
      throw error;
    }
  }
}
```

### 5.2 Payment Gateway — PSP Abstraction Layer

**Candidate:** "We abstract over multiple PSPs so we can failover between them."

```python
# PSP abstraction layer - Python
from abc import ABC, abstractmethod
from dataclasses import dataclass
from enum import Enum

class PSPProvider(Enum):
    STRIPE = "stripe"
    ADYEN = "adyen"
    SQUARE = "square"

@dataclass
class PSPChargeRequest:
    amount_cents: int
    currency: str
    payment_method_token: str
    idempotency_key: str
    merchant_reference: str

@dataclass
class PSPChargeResponse:
    success: bool
    reference: str
    status: str
    raw_response: dict

class PSPGateway(ABC):
    @abstractmethod
    def charge(self, request: PSPChargeRequest) -> PSPChargeResponse:
        pass

    @abstractmethod
    def refund(self, psp_reference: str, amount_cents: int) -> PSPChargeResponse:
        pass

    @abstractmethod
    def get_status(self, psp_reference: str) -> str:
        pass


class StripeGateway(PSPGateway):
    def charge(self, request: PSPChargeRequest) -> PSPChargeResponse:
        response = stripe.PaymentIntent.create(
            amount=request.amount_cents,
            currency=request.currency,
            payment_method=request.payment_method_token,
            confirm=True,
            idempotency_key=request.idempotency_key,
        )
        return PSPChargeResponse(
            success=response.status == "succeeded",
            reference=response.id,
            status=response.status,
            raw_response=response.to_dict(),
        )

    def refund(self, psp_reference: str, amount_cents: int) -> PSPChargeResponse:
        response = stripe.Refund.create(
            payment_intent=psp_reference,
            amount=amount_cents,
        )
        return PSPChargeResponse(
            success=response.status == "succeeded",
            reference=response.id,
            status=response.status,
            raw_response=response.to_dict(),
        )


class PSPRouter:
    """Routes payment requests to the appropriate PSP with failover."""

    def __init__(self, providers: dict[PSPProvider, PSPGateway]):
        self.providers = providers
        self.primary = PSPProvider.STRIPE
        self.fallback = PSPProvider.ADYEN

    def charge(self, request: PSPChargeRequest) -> PSPChargeResponse:
        try:
            return self.providers[self.primary].charge(request)
        except PSPUnavailableError:
            # Failover to secondary PSP
            return self.providers[self.fallback].charge(request)
```

### 5.3 Ledger Service (Double-Entry Accounting)

**Candidate:** "Every financial operation creates balanced ledger entries. This is the source of truth for money movement."

```python
class LedgerService:
    """
    Double-entry bookkeeping: every transaction creates at least two entries
    that sum to zero. This invariant is enforced in the database.
    """

    def record_payment(self, payment_id: str, merchant_id: str,
                       amount_cents: int, currency: str):
        entries = [
            LedgerEntry(
                payment_id=payment_id,
                account_id=f"customer_receivable",
                entry_type="DEBIT",
                amount_cents=amount_cents,
                currency=currency,
            ),
            LedgerEntry(
                payment_id=payment_id,
                account_id=f"merchant_{merchant_id}_payable",
                entry_type="CREDIT",
                amount_cents=amount_cents - self.calculate_fee(amount_cents),
                currency=currency,
            ),
            LedgerEntry(
                payment_id=payment_id,
                account_id="platform_revenue",
                entry_type="CREDIT",
                amount_cents=self.calculate_fee(amount_cents),
                currency=currency,
            ),
        ]
        # All entries written atomically — total DEBIT == total CREDIT
        self.write_entries_atomic(entries)
```

### 5.4 Idempotency Layer

```javascript
// Redis-based idempotency store
class IdempotencyStore {
  constructor(redis) {
    this.redis = redis;
    this.TTL = 24 * 60 * 60; // 24 hours
  }

  async get(key) {
    const data = await this.redis.get(`idempotency:${key}`);
    return data ? JSON.parse(data) : null;
  }

  async set(key, response, ttl = this.TTL) {
    await this.redis.set(
      `idempotency:${key}`,
      JSON.stringify(response),
      'EX', ttl
    );
  }

  async acquireLock(key) {
    // Prevent concurrent processing of the same idempotency key
    const acquired = await this.redis.set(
      `lock:idempotency:${key}`,
      'processing',
      'NX', 'EX', 30  // 30-second lock TTL
    );
    return acquired === 'OK';
  }

  async releaseLock(key) {
    await this.redis.del(`lock:idempotency:${key}`);
  }
}
```

---

## 6. Deep Dive: Database Design

### 6.1 Schema Design

**Candidate:** "I'll use PostgreSQL with serializable isolation for all financial writes. Here's the schema."

```sql
-- Payments: core transaction record
CREATE TABLE payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    merchant_id     UUID NOT NULL REFERENCES merchant_accounts(id),
    amount_cents    BIGINT NOT NULL CHECK (amount_cents > 0),
    currency        CHAR(3) NOT NULL,               -- ISO 4217
    status          VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    payment_method  VARCHAR(255) NOT NULL,
    psp_provider    VARCHAR(50),
    psp_reference   VARCHAR(255),
    description     TEXT,
    metadata        JSONB DEFAULT '{}',
    failure_reason  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_status CHECK (
        status IN ('CREATED','PROCESSING','AUTHORIZED','CAPTURED',
                   'SETTLED','FAILED','CANCELLED','REFUNDED','PARTIALLY_REFUNDED')
    )
);

CREATE INDEX idx_payments_merchant    ON payments (merchant_id, created_at DESC);
CREATE INDEX idx_payments_status      ON payments (status);
CREATE INDEX idx_payments_psp_ref     ON payments (psp_reference);
CREATE INDEX idx_payments_created     ON payments (created_at DESC);

-- Ledger Entries: double-entry bookkeeping (immutable, append-only)
CREATE TABLE ledger_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID NOT NULL REFERENCES payments(id),
    account_id      VARCHAR(255) NOT NULL,
    entry_type      VARCHAR(6) NOT NULL,             -- DEBIT or CREDIT
    amount_cents    BIGINT NOT NULL CHECK (amount_cents > 0),
    currency        CHAR(3) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_entry_type CHECK (entry_type IN ('DEBIT', 'CREDIT'))
);

CREATE INDEX idx_ledger_payment  ON ledger_entries (payment_id);
CREATE INDEX idx_ledger_account  ON ledger_entries (account_id, created_at DESC);

-- Refunds
CREATE TABLE refunds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID NOT NULL REFERENCES payments(id),
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    amount_cents    BIGINT NOT NULL CHECK (amount_cents > 0),
    currency        CHAR(3) NOT NULL,
    reason          VARCHAR(100),
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    psp_reference   VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_refund_status CHECK (
        status IN ('PENDING','PROCESSING','COMPLETED','FAILED')
    )
);

CREATE INDEX idx_refunds_payment ON refunds (payment_id);

-- Merchant Accounts
CREATE TABLE merchant_accounts (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(255) NOT NULL,
    email             VARCHAR(255) NOT NULL,
    balance_cents     BIGINT NOT NULL DEFAULT 0,
    currency          CHAR(3) NOT NULL DEFAULT 'USD',
    payout_schedule   VARCHAR(20) NOT NULL DEFAULT 'DAILY',
    bank_account_id   VARCHAR(255),
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_payout_schedule CHECK (
        payout_schedule IN ('DAILY','WEEKLY','MONTHLY','MANUAL')
    )
);

-- Payment Status Audit Log (immutable, append-only)
CREATE TABLE payment_status_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID NOT NULL REFERENCES payments(id),
    from_status     VARCHAR(20),
    to_status       VARCHAR(20) NOT NULL,
    reason          TEXT,
    actor           VARCHAR(100) NOT NULL,           -- 'system', 'merchant', 'admin'
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_status_log_payment ON payment_status_log (payment_id, created_at);
```

### 6.2 Ledger Balance Invariant

```sql
-- Verify ledger balance: DEBITS must always equal CREDITS per payment
-- This query should always return zero rows in a healthy system
SELECT payment_id,
       SUM(CASE WHEN entry_type = 'DEBIT' THEN amount_cents ELSE 0 END) AS total_debits,
       SUM(CASE WHEN entry_type = 'CREDIT' THEN amount_cents ELSE 0 END) AS total_credits
FROM ledger_entries
GROUP BY payment_id
HAVING SUM(CASE WHEN entry_type = 'DEBIT' THEN amount_cents ELSE 0 END) !=
       SUM(CASE WHEN entry_type = 'CREDIT' THEN amount_cents ELSE 0 END);
```

### 6.3 Why PostgreSQL with Serializable Isolation?

**Interviewer:** "Why not eventual consistency here?"

**Candidate:** "In payment systems, eventual consistency is unacceptable for write paths:
- **Double-charge risk:** Two concurrent retries could both succeed without serializable isolation
- **Ledger imbalance:** Partial writes to the ledger violate double-entry invariants
- **Regulatory requirement:** Auditors require a single source of truth with provable consistency

We use `SERIALIZABLE` isolation for payment writes and ledger entries. Read-heavy workloads (dashboards, reports) can use read replicas with `READ COMMITTED` isolation."

```sql
-- Payment creation uses serializable isolation
BEGIN ISOLATION LEVEL SERIALIZABLE;

INSERT INTO payments (idempotency_key, merchant_id, amount_cents, currency, status)
VALUES ('txn_abc123', 'merch_001', 9999, 'USD', 'CREATED');

-- ... PSP call happens here ...

UPDATE payments SET status = 'CAPTURED', psp_reference = 'ch_3abc123'
WHERE id = 'pay_7x8y9z' AND status = 'PROCESSING';

INSERT INTO ledger_entries (payment_id, account_id, entry_type, amount_cents, currency)
VALUES
  ('pay_7x8y9z', 'customer_receivable',       'DEBIT',  9999, 'USD'),
  ('pay_7x8y9z', 'merchant_merch001_payable', 'CREDIT', 9710, 'USD'),
  ('pay_7x8y9z', 'platform_revenue',          'CREDIT',  289, 'USD');

COMMIT;
```

---

## 7. Deep Dive: Payment Processing Engine

> **This is the KEY section — interviewers spend the most time here.**

### 7.1 Payment State Machine

**Candidate:** "Every payment follows a strict state machine. Illegal transitions are rejected by the system."

```
                    ┌──────────┐
                    │ CREATED  │
                    └────┬─────┘
                         │ initiate()
                    ┌────▼──────┐
            ┌───────│PROCESSING │───────┐
            │       └────┬──────┘       │
            │ fail()     │ authorize()  │ timeout()
            │            │              │
    ┌───────▼──┐   ┌─────▼──────┐  ┌───▼────────┐
    │  FAILED  │   │ AUTHORIZED │  │ UNCERTAIN  │
    └──────────┘   └─────┬──────┘  └───┬────────┘
                         │ capture()    │ resolve()
                         │              │
                    ┌────▼──────┐       │
                    │ CAPTURED  │◄──────┘
                    └────┬──────┘
                         │ settle()
                    ┌────▼──────┐
                    │  SETTLED  │
                    └────┬──────┘
                         │ refund()
              ┌──────────┼──────────┐
              │                     │
    ┌─────────▼──────────┐  ┌──────▼─────────────┐
    │  PARTIALLY_REFUNDED│  │     REFUNDED        │
    └────────────────────┘  └────────────────────┘
```

```python
# State machine implementation
class PaymentStateMachine:
    VALID_TRANSITIONS = {
        'CREATED':             ['PROCESSING'],
        'PROCESSING':          ['AUTHORIZED', 'CAPTURED', 'FAILED', 'UNCERTAIN'],
        'AUTHORIZED':          ['CAPTURED', 'CANCELLED', 'FAILED'],
        'CAPTURED':            ['SETTLED', 'PARTIALLY_REFUNDED', 'REFUNDED'],
        'SETTLED':             ['PARTIALLY_REFUNDED', 'REFUNDED'],
        'UNCERTAIN':           ['CAPTURED', 'FAILED'],
        'FAILED':              [],                # terminal state
        'CANCELLED':           [],                # terminal state
        'REFUNDED':            [],                # terminal state
        'PARTIALLY_REFUNDED':  ['REFUNDED'],
    }

    @staticmethod
    def transition(payment, new_status: str, reason: str = None):
        current = payment.status
        if new_status not in PaymentStateMachine.VALID_TRANSITIONS.get(current, []):
            raise InvalidStateTransition(
                f"Cannot transition from {current} to {new_status}"
            )
        payment.status = new_status
        payment.updated_at = datetime.utcnow()
        # Append to audit log
        PaymentStatusLog.create(
            payment_id=payment.id,
            from_status=current,
            to_status=new_status,
            reason=reason,
            actor='system',
        )
        return payment
```

### 7.2 Idempotency Implementation

**Candidate:** "Idempotency is implemented at three levels: API layer, database constraint, and PSP layer."

```javascript
// Idempotency middleware
async function idempotencyMiddleware(req, res, next) {
  const key = req.headers['idempotency-key'];

  if (!key) {
    return res.status(400).json({ error: 'Idempotency-Key header is required' });
  }

  // Level 1: Check Redis for cached response
  const cached = await idempotencyStore.get(key);
  if (cached) {
    res.set('Idempotent-Replayed', 'true');
    return res.status(cached.statusCode).json(cached.body);
  }

  // Level 2: Acquire distributed lock to prevent concurrent processing
  const lockAcquired = await idempotencyStore.acquireLock(key);
  if (!lockAcquired) {
    return res.status(409).json({ error: 'Request already in progress' });
  }

  // Level 3: Database UNIQUE constraint on idempotency_key prevents duplicates
  // even if Redis fails

  // Wrap response to capture and cache it
  const originalJson = res.json.bind(res);
  res.json = async (body) => {
    await idempotencyStore.set(key, {
      statusCode: res.statusCode,
      body: body,
    });
    await idempotencyStore.releaseLock(key);
    return originalJson(body);
  };

  next();
}
```

**Key design decisions:**
- **TTL:** 24 hours — long enough for retry storms, short enough to not bloat Redis
- **Key format:** Provided by the client (e.g., `order_12345_payment_1`) — the client controls retry semantics
- **Lock TTL:** 30 seconds — prevents stuck locks from blocking retries indefinitely

### 7.3 Double-Entry Ledger

**Candidate:** "The ledger is the financial source of truth. Every money movement produces balanced entries."

#### Example: Full Payment Flow Ledger Entries

```
Payment of $99.99 (9999 cents) with 2.9% + 30¢ fee:

Fee = round(9999 × 0.029) + 30 = 290 + 30 = 320 cents ($3.20)
Merchant receives = 9999 - 320 = 9679 cents ($96.79)

Step 1: Payment Captured
┌─────────────────────────────────┬────────┬───────┬──────────┐
│ Account                         │ Type   │ Cents │ Balance  │
├─────────────────────────────────┼────────┼───────┼──────────┤
│ customer_receivable             │ DEBIT  │ 9999  │          │
│ merchant_merch001_payable       │ CREDIT │ 9679  │          │
│ platform_revenue                │ CREDIT │  320  │          │
├─────────────────────────────────┼────────┼───────┼──────────┤
│ TOTAL                           │        │       │ 0 ✅      │
└─────────────────────────────────┴────────┴───────┴──────────┘

Step 2: Partial Refund of $50.00 (5000 cents)
Refund fee returned proportionally: round(5000/9999 × 320) = 160 cents

┌─────────────────────────────────┬────────┬───────┬──────────┐
│ Account                         │ Type   │ Cents │ Balance  │
├─────────────────────────────────┼────────┼───────┼──────────┤
│ customer_receivable             │ CREDIT │ 5000  │          │
│ merchant_merch001_payable       │ DEBIT  │ 4840  │          │
│ platform_revenue                │ DEBIT  │  160  │          │
├─────────────────────────────────┼────────┼───────┼──────────┤
│ TOTAL                           │        │       │ 0 ✅      │
└─────────────────────────────────┴────────┴───────┴──────────┘

Step 3: Merchant Payout of $4,839 (accumulated balance)
┌─────────────────────────────────┬────────┬─────────┬────────┐
│ Account                         │ Type   │ Cents   │Balance │
├─────────────────────────────────┼────────┼─────────┼────────┤
│ merchant_merch001_payable       │ DEBIT  │ 483900  │        │
│ merchant_merch001_bank          │ CREDIT │ 483900  │        │
├─────────────────────────────────┼────────┼─────────┼────────┤
│ TOTAL                           │        │         │ 0 ✅    │
└─────────────────────────────────┴────────┴─────────┴────────┘
```

### 7.4 Exactly-Once Semantics

**Candidate:** "Exactly-once processing is achieved through three reinforcing mechanisms."

```
Exactly-Once = Idempotency + State Machine + DB Transactions

Layer 1: Idempotency Key (client-provided)
  → Same key always returns the same response
  → Prevents duplicate payment creation

Layer 2: State Machine (server-enforced)
  → PROCESSING → CAPTURED is a one-way transition
  → Cannot charge twice because re-entering PROCESSING is illegal

Layer 3: Database Transaction (SERIALIZABLE isolation)
  → status update + ledger write are atomic
  → If any step fails, everything rolls back
```

```python
def process_payment_exactly_once(idempotency_key: str, request: PaymentRequest):
    # Layer 1: Idempotency — return cached result if key exists
    cached = idempotency_store.get(idempotency_key)
    if cached:
        return cached

    with db.transaction(isolation_level='SERIALIZABLE'):
        # Layer 2: State Machine — create or find payment
        payment = Payment.get_or_create(idempotency_key=idempotency_key)

        if payment.status == 'CAPTURED':
            return payment  # already processed

        if payment.status not in ('CREATED', 'PROCESSING'):
            raise InvalidStateError(f"Payment in terminal state: {payment.status}")

        # Transition: CREATED → PROCESSING
        PaymentStateMachine.transition(payment, 'PROCESSING')
        db.save(payment)

        # Call PSP (with PSP-level idempotency key)
        psp_result = psp_gateway.charge(request, idempotency_key=idempotency_key)

        if psp_result.success:
            # Layer 3: Atomic state update + ledger entries
            PaymentStateMachine.transition(payment, 'CAPTURED')
            payment.psp_reference = psp_result.reference
            db.save(payment)

            ledger_service.record_payment(payment)  # within same transaction
        else:
            PaymentStateMachine.transition(payment, 'FAILED', reason=psp_result.error)
            db.save(payment)

    # Cache the result outside the transaction
    idempotency_store.set(idempotency_key, payment)
    return payment
```

### 7.5 PSP Integration and Retry Logic

```python
class PSPRetryHandler:
    MAX_RETRIES = 3
    RETRY_DELAYS = [1, 2, 4]  # exponential backoff in seconds

    def charge_with_retry(self, request: PSPChargeRequest) -> PSPChargeResponse:
        last_error = None
        for attempt in range(self.MAX_RETRIES):
            try:
                response = self.psp_gateway.charge(request)
                if response.success:
                    return response
                if not self._is_retryable(response):
                    return response  # permanent failure (e.g., insufficient funds)
                last_error = response
            except PSPTimeoutError:
                # CRITICAL: PSP timeout = uncertain state
                # Do NOT retry blindly — check payment status first
                status = self.psp_gateway.get_status(request.idempotency_key)
                if status == 'succeeded':
                    return PSPChargeResponse(success=True, reference=status.reference)
                elif status == 'failed':
                    return PSPChargeResponse(success=False)
                # status is still pending — retry with same idempotency key
                last_error = PSPTimeoutError()
            except PSPConnectionError as e:
                last_error = e

            time.sleep(self.RETRY_DELAYS[attempt])

        # After all retries exhausted, mark payment as UNCERTAIN
        raise PaymentUncertainError(f"PSP unreachable after {self.MAX_RETRIES} attempts")

    def _is_retryable(self, response):
        non_retryable = {'insufficient_funds', 'card_declined', 'invalid_card'}
        return response.error_code not in non_retryable
```

### 7.6 Currency Handling

**Candidate:** "Currencies are always stored in minor units. We never use floating point for money."

```python
CURRENCY_MINOR_UNITS = {
    'USD': 2,  # 1 dollar = 100 cents
    'EUR': 2,  # 1 euro = 100 cents
    'GBP': 2,  # 1 pound = 100 pence
    'JPY': 0,  # yen has no minor unit
    'INR': 2,  # 1 rupee = 100 paise
    'BHD': 3,  # 1 dinar = 1000 fils
    'KWD': 3,  # 1 dinar = 1000 fils
}

def to_minor_units(amount_decimal: str, currency: str) -> int:
    """Convert human-readable amount to minor units (cents/paise/etc.)."""
    exponent = CURRENCY_MINOR_UNITS.get(currency, 2)
    from decimal import Decimal, ROUND_HALF_UP
    d = Decimal(amount_decimal)
    return int(d * (10 ** exponent))

def from_minor_units(amount_cents: int, currency: str) -> str:
    """Convert minor units back to human-readable string."""
    exponent = CURRENCY_MINOR_UNITS.get(currency, 2)
    from decimal import Decimal
    return str(Decimal(amount_cents) / (10 ** exponent))

# Multi-currency conversion
class CurrencyConverter:
    def __init__(self, rate_provider):
        self.rate_provider = rate_provider  # e.g., Open Exchange Rates API

    def convert(self, amount_cents: int, from_currency: str, to_currency: str) -> int:
        if from_currency == to_currency:
            return amount_cents
        rate = self.rate_provider.get_rate(from_currency, to_currency)
        from decimal import Decimal, ROUND_HALF_UP
        converted = Decimal(amount_cents) * Decimal(str(rate))
        return int(converted.quantize(Decimal('1'), rounding=ROUND_HALF_UP))
```

### 7.7 Reconciliation Service

**Candidate:** "Reconciliation compares our internal ledger with PSP settlement reports to catch discrepancies."

```python
class ReconciliationService:
    """
    Runs daily: compares internal payment records against PSP settlement files.
    Flags mismatches for manual review.
    """

    def reconcile_daily(self, date: str):
        # 1. Fetch our records for the date
        internal_payments = self.payment_repo.get_captured_by_date(date)

        # 2. Fetch PSP settlement report
        psp_settlements = self.psp_gateway.get_settlement_report(date)

        # 3. Build lookup maps
        internal_map = {p.psp_reference: p for p in internal_payments}
        psp_map = {s.reference: s for s in psp_settlements}

        discrepancies = []

        # 4. Check: every internal record exists in PSP
        for ref, payment in internal_map.items():
            if ref not in psp_map:
                discrepancies.append({
                    'type': 'MISSING_IN_PSP',
                    'reference': ref,
                    'internal_amount': payment.amount_cents,
                })
            elif payment.amount_cents != psp_map[ref].amount_cents:
                discrepancies.append({
                    'type': 'AMOUNT_MISMATCH',
                    'reference': ref,
                    'internal_amount': payment.amount_cents,
                    'psp_amount': psp_map[ref].amount_cents,
                })

        # 5. Check: every PSP record exists internally
        for ref, settlement in psp_map.items():
            if ref not in internal_map:
                discrepancies.append({
                    'type': 'MISSING_INTERNALLY',
                    'reference': ref,
                    'psp_amount': settlement.amount_cents,
                })

        # 6. Report
        if discrepancies:
            self.alert_service.send_alert(
                severity='HIGH',
                message=f"Reconciliation found {len(discrepancies)} discrepancies for {date}",
                details=discrepancies,
            )
            self.discrepancy_repo.save_all(discrepancies)

        return {
            'date': date,
            'internal_count': len(internal_payments),
            'psp_count': len(psp_settlements),
            'discrepancies': len(discrepancies),
            'status': 'CLEAN' if not discrepancies else 'DISCREPANCIES_FOUND',
        }
```

### 7.8 PCI DSS Compliance

**Candidate:** "We never store raw card numbers. Here's how we handle PCI compliance."

```
PCI DSS Compliance Architecture:

┌──────────────┐     ┌──────────────────┐     ┌──────────────┐
│   Customer   │────▶│  Payment Form    │────▶│  Tokenization│
│   Browser    │     │  (Stripe.js /    │     │  Service     │
│              │     │   iframe)        │     │  (PSP-hosted)│
└──────────────┘     └──────────────────┘     └──────┬───────┘
                                                      │
                                                      │ Returns token
                                                      │ (e.g., pm_card_visa_4242)
                                                      ▼
                     ┌──────────────────┐     ┌──────────────┐
                     │  Our Backend     │◄────│  Token only  │
                     │  (never sees     │     │  (no PAN,    │
                     │   raw card data) │     │   no CVV)    │
                     └──────────────────┘     └──────────────┘

Key Compliance Measures:
┌────────────────────────────────────────────────────────────┐
│ 1. TOKENIZATION: Card data goes directly to PSP via       │
│    client-side SDK — our servers never see raw PAN/CVV    │
│                                                            │
│ 2. ENCRYPTION: TLS 1.3 for all API communication          │
│    AES-256 for data at rest                                │
│                                                            │
│ 3. ACCESS CONTROL: Role-based access, MFA for admin       │
│    PCI-scoped network segmentation                         │
│                                                            │
│ 4. AUDIT LOGGING: Every access to payment data is logged  │
│    Immutable audit trail, 7-year retention                 │
│                                                            │
│ 5. KEY MANAGEMENT: HSM-backed encryption keys             │
│    Regular key rotation (90-day cycle)                     │
│                                                            │
│ 6. VULNERABILITY MANAGEMENT: Quarterly ASV scans          │
│    Annual penetration testing                              │
└────────────────────────────────────────────────────────────┘
```

---

## 8. Scaling Strategies

### 8.1 Database Sharding

**Candidate:** "We shard by `merchant_id` using consistent hashing."

```
Shard Strategy: merchant_id-based consistent hashing

  merchant_id ──▶ hash(merchant_id) % num_shards ──▶ shard_N

  Shard 0: merchants A-G     (PostgreSQL cluster 0)
  Shard 1: merchants H-N     (PostgreSQL cluster 1)
  Shard 2: merchants O-T     (PostgreSQL cluster 2)
  Shard 3: merchants U-Z     (PostgreSQL cluster 3)

Advantages:
  ✅ All data for one merchant lives on one shard (no cross-shard transactions)
  ✅ Merchant dashboard queries are single-shard
  ✅ Payout calculations are single-shard

Cross-shard queries (e.g., platform-wide analytics):
  → Handled by OLAP database (see 8.3)
```

### 8.2 Read Replicas

```
Write Path (financial operations):
  Client → Primary DB (SERIALIZABLE isolation)

Read Path (dashboards, status checks):
  Client → Read Replica (READ COMMITTED, slight lag OK)

  Primary ──sync──▶ Replica 1 (same AZ)
           └─async─▶ Replica 2 (different AZ, for HA)
           └─async─▶ Replica 3 (analytics queries)
```

### 8.3 OLTP vs OLAP Separation

```
┌──────────────────┐          ┌─────────────────────┐
│  OLTP (Postgres) │──CDC────▶│  OLAP (ClickHouse/  │
│  Payment writes  │  stream  │  BigQuery)           │
│  Ledger entries  │          │  Analytics queries   │
│  Real-time state │          │  Revenue reports     │
└──────────────────┘          │  Fraud ML features   │
                              └─────────────────────┘

CDC = Change Data Capture (Debezium on Kafka)
```

### 8.4 Event Sourcing for Audit Trail

```javascript
// Every state change is an immutable event
const paymentEvents = [
  { type: 'PaymentCreated',    timestamp: '...', data: { amount: 9999, ... } },
  { type: 'PaymentProcessing', timestamp: '...', data: { psp: 'stripe' } },
  { type: 'PaymentAuthorized', timestamp: '...', data: { auth_code: 'abc' } },
  { type: 'PaymentCaptured',   timestamp: '...', data: { psp_ref: 'ch_123' } },
  { type: 'PaymentSettled',    timestamp: '...', data: { settlement_id: 's_1' } },
];

// Current state = fold over all events
function currentState(events) {
  return events.reduce((state, event) => {
    switch (event.type) {
      case 'PaymentCreated':    return { ...state, status: 'CREATED', ...event.data };
      case 'PaymentProcessing': return { ...state, status: 'PROCESSING' };
      case 'PaymentCaptured':   return { ...state, status: 'CAPTURED', ...event.data };
      case 'PaymentSettled':    return { ...state, status: 'SETTLED', ...event.data };
      default: return state;
    }
  }, {});
}
```

### 8.5 Async Processing for Non-Critical Paths

```
Synchronous (in payment request path):
  ✅ Idempotency check
  ✅ Payment creation in DB
  ✅ PSP charge call
  ✅ Ledger entry creation
  ✅ Response to client

Asynchronous (via Kafka events):
  📨 Webhook delivery to merchant
  📨 Email receipt to customer
  📨 Analytics event ingestion
  📨 Fraud scoring (post-transaction)
  📨 Reconciliation data preparation
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Double Payment (Idempotency Failure)

```
Scenario: Client sends payment, network drops before response arrives,
          client retries with same idempotency key.

Timeline:
  T0: Client → POST /payments (key="txn_001") → Server receives
  T1: Server processes, charges PSP successfully
  T2: Response lost in transit ❌
  T3: Client retries → POST /payments (key="txn_001")
  T4: Server finds key="txn_001" in Redis → returns cached response ✅

Mitigation Stack:
  Layer 1: Redis idempotency cache (fast path)
  Layer 2: DB UNIQUE constraint on idempotency_key (backup)
  Layer 3: PSP-level idempotency key (prevents double charge at PSP)
```

### 9.2 PSP Timeout (Uncertain State)

```
Scenario: We send a charge request to Stripe, but the response times out.
          Did the charge succeed or not? We don't know.

                    ┌──────────┐      ┌──────────┐
                    │   Our    │─────▶│  Stripe  │
                    │  Server  │  ??  │          │   Charge may or may
                    │          │◄─ ✗ ─│          │   not have succeeded
                    └──────────┘      └──────────┘

Resolution Strategy:
  1. Mark payment as UNCERTAIN (not FAILED)
  2. Start async resolution job:
     a. Query PSP for payment status using idempotency key
     b. If PSP says "succeeded" → transition to CAPTURED
     c. If PSP says "failed" → transition to FAILED
     d. If PSP says "not found" → safe to retry or mark FAILED
  3. Retry the status check with exponential backoff (up to 1 hour)
  4. If still unresolved after 1 hour → alert operations team
  5. NEVER auto-retry the charge without checking status first
```

### 9.3 Partial Refund Failure

```
Scenario: Merchant requests $50 refund on a $100 payment.
          PSP processes the refund, but our DB update fails.

Mitigation:
  1. Create refund record (status=PENDING) BEFORE calling PSP
  2. Call PSP refund API with idempotency key
  3. Update refund record (status=COMPLETED)
  4. If step 3 fails: reconciliation service catches the discrepancy
  5. Manual review resolves: update internal records to match PSP state
```

### 9.4 Ledger Imbalance Detection

```python
# Scheduled job: runs every 15 minutes
def detect_ledger_imbalances():
    query = """
        SELECT payment_id,
               SUM(CASE WHEN entry_type='DEBIT' THEN amount_cents ELSE 0 END) as debits,
               SUM(CASE WHEN entry_type='CREDIT' THEN amount_cents ELSE 0 END) as credits
        FROM ledger_entries
        WHERE created_at > NOW() - INTERVAL '1 hour'
        GROUP BY payment_id
        HAVING SUM(CASE WHEN entry_type='DEBIT' THEN amount_cents ELSE 0 END) !=
               SUM(CASE WHEN entry_type='CREDIT' THEN amount_cents ELSE 0 END)
    """
    imbalances = db.execute(query)
    if imbalances:
        alert(severity='CRITICAL', message=f"LEDGER IMBALANCE: {len(imbalances)} payments")
        # Immediately halt new transactions for affected merchants
        for row in imbalances:
            quarantine_payment(row.payment_id)
```

### 9.5 Network Partition During Payment

```
Scenario: Network partition between our service and the database
          DURING a payment that has already been sent to the PSP.

Timeline:
  T0: Payment CREATED in DB
  T1: PSP charge succeeds (Stripe debits card)
  T2: DB write fails (network partition) ❌
  T3: Client receives error, retries
  T4: Network restored, retry finds payment in PROCESSING state
  T5: Query PSP → confirms charge succeeded → update to CAPTURED

Key insight: The PSP is the authoritative source during partition.
  We can always recover by querying the PSP with the idempotency key.
```

### 9.6 Reconciliation Discrepancy

```
Common discrepancy types and resolutions:

| Type              | Cause                          | Resolution                 |
|-------------------|--------------------------------|----------------------------|
| MISSING_IN_PSP    | Our record exists, PSP doesn't | Charge failed silently     |
|                   |                                | → mark as FAILED           |
| MISSING_INTERNAL  | PSP charged, our record lost   | DB failure during write    |
|                   |                                | → create internal record   |
| AMOUNT_MISMATCH   | Currency conversion difference | Rounding discrepancy       |
|                   |                                | → adjust ledger entry      |
| STATUS_MISMATCH   | PSP settled, we show CAPTURED  | Settlement webhook missed  |
|                   |                                | → update to SETTLED        |
```

---

## 10. Monitoring & Observability

### 10.1 Key Metrics Dashboard

```
┌─────────────────────────────────────────────────────────────────┐
│                    PAYMENT SYSTEM DASHBOARD                      │
├─────────────────┬───────────────────┬───────────────────────────┤
│  Success Rate   │  Avg Latency      │  Active Transactions      │
│  ██████████ 98.7%│  ████░░ 850ms    │  ███░░░ 47 in-flight     │
├─────────────────┴───────────────────┴───────────────────────────┤
│                                                                  │
│  Transactions/min (last 1h)          PSP Latency by Provider    │
│  800 ┤                               Stripe:  ████░░ 620ms     │
│  600 ┤    ╭──╮                        Adyen:   █████░ 780ms     │
│  400 ┤───╯  ╰───╮                    Square:  ███░░░ 490ms     │
│  200 ┤           ╰───                                           │
│    0 ┤─────────────────              Refund Rate: 5.2%          │
│      10:00  10:30  11:00             Chargeback Rate: 0.3%      │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│  LEDGER STATUS: ✅ BALANCED          RECONCILIATION: ✅ CLEAN    │
│  Last check: 2 min ago              Last run: 6h ago            │
└──────────────────────────────────────────────────────────────────┘
```

### 10.2 Critical Alerts

| Alert | Condition | Severity | Action |
|---|---|---|---|
| Payment failure spike | Failure rate > 5% (5-min window) | P1 CRITICAL | Page on-call, check PSP status |
| Ledger imbalance | Any payment with DEBIT ≠ CREDIT | P0 CRITICAL | Halt processing, investigate |
| PSP latency high | p99 > 5 seconds | P2 HIGH | Consider failover to backup PSP |
| Reconciliation mismatch | > 10 discrepancies in daily run | P1 CRITICAL | Finance team + engineering review |
| Idempotency store down | Redis cluster unreachable | P1 CRITICAL | Fallback to DB-only idempotency |
| Payment stuck in PROCESSING | > 5 min in PROCESSING state | P2 HIGH | Trigger async resolution job |
| Unusual transaction volume | > 3x normal volume in 5 min | P2 HIGH | Potential fraud or bot attack |

### 10.3 Structured Audit Logging

```javascript
// Every payment action produces a structured audit log entry
const auditLog = {
  timestamp: '2026-04-24T10:30:00.123Z',
  event_type: 'payment.captured',
  payment_id: 'pay_7x8y9z',
  merchant_id: 'merch_001',
  actor: 'system',
  ip_address: '203.0.113.42',
  request_id: 'req_abc123',
  idempotency_key: 'txn_abc123_attempt1',
  details: {
    amount_cents: 9999,
    currency: 'USD',
    psp_provider: 'stripe',
    psp_reference: 'ch_3abc123',
    processing_time_ms: 847,
  },
  compliance: {
    pci_scope: false,        // no card data in this log
    data_classification: 'FINANCIAL',
    retention_years: 7,
  },
};

// Shipped to immutable log store (e.g., S3 + Athena, or Elasticsearch)
```

---

## 11. Advanced Features

### 11.1 Subscription / Recurring Payments

```python
class SubscriptionService:
    def process_recurring_payment(self, subscription_id: str):
        sub = self.subscription_repo.get(subscription_id)

        # Generate deterministic idempotency key for this billing cycle
        idempotency_key = f"sub_{subscription_id}_cycle_{sub.current_period_start}"

        try:
            payment = self.payment_service.create_payment(
                amount_cents=sub.amount_cents,
                currency=sub.currency,
                merchant_id=sub.merchant_id,
                payment_method=sub.default_payment_method,
                idempotency_key=idempotency_key,
            )
            sub.last_payment_id = payment.id
            sub.current_period_start = sub.next_billing_date
            sub.next_billing_date = self.calculate_next_billing(sub)
            self.subscription_repo.save(sub)
        except PaymentFailedError:
            self.handle_dunning(sub)  # retry logic with escalation

    def handle_dunning(self, subscription):
        """Smart retry schedule for failed recurring payments."""
        retry_schedule = [
            (1, 'retry_same_method'),       # Day 1: retry
            (3, 'notify_customer'),          # Day 3: email customer
            (5, 'retry_same_method'),        # Day 5: retry again
            (7, 'try_backup_method'),        # Day 7: try backup card
            (14, 'cancel_subscription'),     # Day 14: cancel
        ]
```

### 11.2 Split Payments (Marketplace Model)

```
Marketplace Payment Split:
Customer pays $100 for an order with items from 2 sellers

┌──────────────────────────────────────────────────────────┐
│  Customer pays: $100.00                                   │
│                                                           │
│  Split:                                                   │
│    Seller A (item $60):  $60 - $1.74 fee = $58.26        │
│    Seller B (item $40):  $40 - $1.16 fee = $38.84        │
│    Platform fee:         $1.74 + $1.16   = $2.90         │
│                                                           │
│  Ledger:                                                  │
│    DEBIT   customer_receivable       $100.00              │
│    CREDIT  seller_a_payable           $58.26              │
│    CREDIT  seller_b_payable           $38.84              │
│    CREDIT  platform_revenue            $2.90              │
│    Balance:                            $0.00 ✅            │
└──────────────────────────────────────────────────────────┘
```

### 11.3 3D Secure Authentication

```
3DS Flow (for high-risk transactions):

  Client → Payment API → PSP returns "requires_action"
                              │
                              ▼
                   ┌─────────────────────┐
                   │  Redirect customer  │
                   │  to bank's 3DS page │
                   └──────────┬──────────┘
                              │
                   Customer enters OTP / biometric
                              │
                   ┌──────────▼──────────┐
                   │  Bank confirms auth │
                   │  Callback to our    │
                   │  return URL         │
                   └──────────┬──────────┘
                              │
                   ┌──────────▼──────────┐
                   │  Confirm payment    │
                   │  with PSP           │
                   │  (status → CAPTURED)│
                   └─────────────────────┘

Payment status during 3DS: AUTHORIZED (awaiting customer action)
```

### 11.4 Webhook Notifications to Merchants

```javascript
class WebhookDeliveryService {
  async deliverWebhook(merchantId, event) {
    const merchant = await this.merchantRepo.get(merchantId);
    const payload = {
      id: `evt_${uuid()}`,
      type: event.type,
      data: event.data,
      created_at: new Date().toISOString(),
    };

    const signature = this.signPayload(payload, merchant.webhook_secret);

    // Retry with exponential backoff: 1s, 5s, 30s, 2m, 10m, 1h, 6h
    const retrySchedule = [1000, 5000, 30000, 120000, 600000, 3600000, 21600000];

    for (let attempt = 0; attempt <= retrySchedule.length; attempt++) {
      try {
        const response = await fetch(merchant.webhook_url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Webhook-Signature': signature,
            'X-Webhook-Id': payload.id,
          },
          body: JSON.stringify(payload),
          timeout: 5000,
        });

        if (response.status >= 200 && response.status < 300) {
          await this.logDelivery(payload.id, 'DELIVERED', attempt);
          return;
        }
      } catch (err) { /* retry */ }

      if (attempt < retrySchedule.length) {
        await this.scheduleRetry(payload, retrySchedule[attempt]);
      }
    }
    await this.logDelivery(payload.id, 'FAILED_PERMANENTLY');
    await this.alertMerchantWebhookFailure(merchantId);
  }

  signPayload(payload, secret) {
    const hmac = crypto.createHmac('sha256', secret);
    hmac.update(JSON.stringify(payload));
    return hmac.digest('hex');
  }
}
```

### 11.5 Dispute / Chargeback Management

```
Chargeback Flow:
  1. Customer disputes charge with their bank
  2. Bank notifies PSP → PSP sends webhook to us
  3. We create a Dispute record, debit merchant account
  4. Merchant can submit evidence (receipts, delivery proof)
  5. We forward evidence to PSP → bank makes final decision
  6. If merchant wins: re-credit merchant account
  7. If merchant loses: dispute amount + fee permanently deducted

Ledger entries for chargeback:
  DEBIT   merchant_merch001_payable     $100.00  (hold funds)
  CREDIT  chargeback_reserve            $100.00
```

### 11.6 Fraud Detection

```python
class FraudDetectionService:
    def score_transaction(self, payment: PaymentRequest) -> FraudScore:
        signals = []

        # Velocity checks
        recent_count = self.get_recent_txn_count(
            payment.payment_method, window_minutes=60
        )
        if recent_count > 10:
            signals.append(('high_velocity', 0.8))

        # Amount anomaly
        avg_amount = self.get_avg_amount(payment.merchant_id)
        if payment.amount_cents > avg_amount * 5:
            signals.append(('amount_anomaly', 0.6))

        # Geographic mismatch
        if self.ip_country(payment.ip) != self.card_country(payment.payment_method):
            signals.append(('geo_mismatch', 0.7))

        # New card on high-risk merchant
        if self.is_first_use(payment.payment_method) and \
           self.is_high_risk_merchant(payment.merchant_id):
            signals.append(('new_card_high_risk', 0.5))

        # Aggregate risk score
        risk_score = min(1.0, sum(s[1] for s in signals))

        if risk_score > 0.9:
            return FraudScore(action='BLOCK', score=risk_score, signals=signals)
        elif risk_score > 0.6:
            return FraudScore(action='REVIEW', score=risk_score, signals=signals)
        else:
            return FraudScore(action='ALLOW', score=risk_score, signals=signals)
```

---

## 12. Interview Q&A

### Q1: How do you ensure exactly-once payment processing?

**Candidate:** "Exactly-once is achieved through three reinforcing layers:

1. **Client-provided idempotency key** — included in every payment request. The same key always returns the same result. Cached in Redis with a 24-hour TTL and backed by a UNIQUE database constraint.

2. **Payment state machine** — enforces legal transitions. Once a payment reaches CAPTURED, it cannot be re-processed. The transition from CREATED → PROCESSING → CAPTURED is one-directional.

3. **Database transaction with SERIALIZABLE isolation** — the status update and ledger entry creation happen atomically. If the ledger write fails, the entire transaction rolls back, including the status change.

Additionally, the PSP receives our idempotency key, so even if we accidentally send two charge requests, the PSP only processes one."

---

### Q2: What happens when the PSP times out — how do you handle uncertain payment state?

**Candidate:** "This is one of the hardest problems in payment systems. When a PSP times out, we genuinely don't know if the charge succeeded.

1. We transition the payment to an **UNCERTAIN** state — not FAILED.
2. We **never** auto-retry the charge, because the original might have succeeded.
3. We start an async resolution job that queries the PSP's status API using the idempotency key.
4. Based on the PSP response:
   - `succeeded` → transition to CAPTURED, create ledger entries
   - `failed` → transition to FAILED
   - `not_found` → safe to consider it FAILED or retry
5. We retry the status query with exponential backoff (1s, 5s, 30s, 2min, 10min) up to 1 hour.
6. If still unresolved after 1 hour, we page the on-call engineer for manual resolution.

The idempotency key is critical here — it's our correlation ID across the uncertain state."

---

### Q3: Explain double-entry bookkeeping in the context of payments.

**Candidate:** "Double-entry bookkeeping means every financial transaction produces at least two ledger entries that sum to zero: a DEBIT and a CREDIT of equal total value.

For example, when a customer pays $100 with a 2.9% + $0.30 fee:
- **DEBIT** customer_receivable: $100.00
- **CREDIT** merchant_payable: $96.80 (merchant's share)
- **CREDIT** platform_revenue: $3.20 (our fee)
- Total: $100.00 DEBIT = $100.00 CREDIT ✅

This invariant lets us detect any accounting error immediately. We run a ledger balance check every 15 minutes, and any imbalance triggers a P0 critical alert.

The ledger is append-only — we never update or delete entries. Corrections are made by adding new counterbalancing entries (reversals), which maintains a complete audit trail."

---

### Q4: How do you handle multi-currency transactions?

**Candidate:** "Three key principles:

1. **Store in minor units (integers)** — $99.99 is stored as 9999 cents, ¥1000 as 1000 (JPY has no minor unit), 1.234 KWD as 1234 fils (3 decimal places). This avoids floating-point errors.

2. **Record the original currency** — the payment record stores both the amount and the currency the customer paid in. Conversion happens at settlement time, not at payment time.

3. **Lock the exchange rate at transaction time** — we snapshot the exchange rate when the payment is created and store it in the payment metadata. The merchant sees the exact amount they'll receive in their payout currency.

For FX conversion, we use a rate provider (e.g., Open Exchange Rates) with a markup, and the conversion is done using Python's `Decimal` type with explicit rounding rules (ROUND_HALF_UP) to avoid accumulating rounding errors."

---

### Q5: How would you implement reconciliation between your system and the PSP?

**Candidate:** "Reconciliation runs as a daily batch job comparing our internal records with the PSP's settlement reports:

1. **Fetch internal data** — all payments with status CAPTURED or SETTLED for the date
2. **Fetch PSP settlement file** — CSV/API export from Stripe/Adyen
3. **Three-way comparison:**
   - Every internal record should exist in PSP → if not: `MISSING_IN_PSP`
   - Every PSP record should exist internally → if not: `MISSING_INTERNALLY`
   - Amounts must match → if not: `AMOUNT_MISMATCH`
4. **Discrepancy resolution** — auto-resolve where possible (e.g., status updates), flag complex cases for manual review
5. **Reporting** — daily reconciliation report to finance team

We also run continuous micro-reconciliation: after every payment, we compare our recorded PSP reference and amount with a callback/webhook from the PSP. This catches issues in near-real-time instead of waiting for the daily batch."

---

### Q6: How do you handle PCI DSS compliance?

**Candidate:** "The cardinal rule is: **never let raw card data touch our servers.**

1. **Client-side tokenization** — the customer enters card details into a PSP-hosted iframe (Stripe Elements, Adyen Drop-in). Card data goes directly to the PSP; we only receive a token (e.g., `pm_card_visa_4242`).

2. **SAQ A compliance** — since we never process, store, or transmit cardholder data, we qualify for the simplest PCI DSS self-assessment questionnaire.

3. **Defense in depth:**
   - TLS 1.3 for all API traffic
   - AES-256 encryption at rest for all financial data
   - Network segmentation — payment services in isolated VPC
   - HSM-backed key management with 90-day rotation
   - Role-based access control with MFA

4. **Audit requirements** — immutable logs, quarterly ASV scans, annual penetration testing, regular access reviews."

---

### Q7: How would you design a fraud detection system for payments?

**Candidate:** "I'd implement a multi-layered fraud detection system:

**Layer 1 — Rule-based (real-time, synchronous):**
- Velocity checks: block if > 10 transactions in 1 hour from same card
- Amount thresholds: flag transactions > 5x the merchant's average
- Geographic mismatch: IP country ≠ card issuing country
- BIN checks: block known high-risk card BINs

**Layer 2 — ML scoring (real-time, synchronous):**
- Feature vector: amount, time of day, device fingerprint, IP reputation, card age, merchant category
- Model: gradient-boosted tree trained on historical fraud data
- Output: risk score 0.0–1.0
- Actions: ALLOW (< 0.6), REVIEW (0.6–0.9), BLOCK (> 0.9)

**Layer 3 — Post-transaction analysis (async):**
- Network analysis: graph of shared devices, IPs, addresses across accounts
- Pattern detection: unusual purchasing patterns, testing card behavior (many small transactions)
- Chargeback feedback loop: confirmed fraud feeds back into ML training

Decisions are made in < 100ms to stay within the 2-second payment SLA."

---

### Q8: How do you handle payment retries without double-charging?

**Candidate:** "Payment retries are safe because of our idempotency architecture:

1. **Client retries with the same idempotency key** — the system returns the cached result instead of creating a new charge.

2. **Three-layer protection:**
   - **Redis:** fast lookup of previously processed idempotency keys (O(1))
   - **PostgreSQL:** UNIQUE constraint on `idempotency_key` column (backup if Redis is down)
   - **PSP:** we pass our idempotency key to the PSP, so even if we accidentally send two requests, only one charge is created

3. **Distributed lock:** when a request arrives with a new idempotency key, we acquire a Redis lock (NX with 30-second TTL) to prevent concurrent processing of the same key. If a second request arrives while the first is in-flight, it receives a 409 Conflict.

4. **Key ownership:** the idempotency key is generated by the client (e.g., `order_12345_payment_attempt_1`). This means the client controls retry semantics. If they want a new charge, they use a new key."

---

## 13. Production Checklist

### Pre-Launch (Week -2 to -1)

| # | Item | Status |
|---|---|---|
| 1 | PCI DSS SAQ-A completed and filed | ☐ |
| 2 | Penetration test conducted, findings remediated | ☐ |
| 3 | PSP sandbox integration tested (all payment methods) | ☐ |
| 4 | PSP production credentials provisioned and secured in vault | ☐ |
| 5 | Idempotency layer load-tested (10x expected peak) | ☐ |
| 6 | Ledger balance verification job running green for 7 days | ☐ |
| 7 | Database backups tested (restore verified) | ☐ |
| 8 | Encryption at rest enabled and key rotation tested | ☐ |
| 9 | Runbooks for P0 scenarios (ledger imbalance, PSP outage) written | ☐ |
| 10 | Legal review of merchant agreements completed | ☐ |

### Day 1

| # | Item | Status |
|---|---|---|
| 1 | Gradual rollout: 1% of traffic to new payment system | ☐ |
| 2 | Real-time monitoring dashboard active (Grafana) | ☐ |
| 3 | On-call engineer assigned with access to PSP dashboard | ☐ |
| 4 | Manual reconciliation of all Day 1 transactions | ☐ |
| 5 | Webhook delivery success rate monitored (target: > 99%) | ☐ |

### Week 1

| # | Item | Status |
|---|---|---|
| 1 | Ramp to 25% traffic, monitor success rates | ☐ |
| 2 | First automated reconciliation run completed | ☐ |
| 3 | Review all FAILED and UNCERTAIN payments manually | ☐ |
| 4 | Validate ledger balance across all merchants | ☐ |
| 5 | Performance baseline established (p50, p95, p99 latencies) | ☐ |
| 6 | First merchant payout executed and verified | ☐ |

### Month 1

| # | Item | Status |
|---|---|---|
| 1 | Ramp to 100% traffic | ☐ |
| 2 | Automated reconciliation running daily, zero discrepancies | ☐ |
| 3 | Chargeback handling process tested end-to-end | ☐ |
| 4 | Fraud detection rules tuned based on live data | ☐ |
| 5 | Capacity planning reviewed based on actual traffic patterns | ☐ |
| 6 | PSP failover tested (simulate primary PSP outage) | ☐ |
| 7 | Disaster recovery drill completed (database failover) | ☐ |
| 8 | Quarterly ASV scan scheduled | ☐ |

---

## Summary

| Aspect | Decision | Rationale |
|---|---|---|
| **Primary DB** | PostgreSQL (SERIALIZABLE) | Financial data requires strong consistency |
| **Idempotency** | Redis + DB UNIQUE constraint | Fast dedup with durable backup |
| **Event Bus** | Kafka | Reliable async processing, audit trail |
| **PSP Strategy** | Multi-PSP with failover | No single point of failure for payments |
| **Ledger** | Double-entry bookkeeping | Financial industry standard, self-verifying |
| **Currency** | Integer minor units (cents) | Eliminates floating-point errors |
| **Sharding** | By merchant_id | Keeps merchant data co-located |
| **Compliance** | PCI DSS SAQ-A + tokenization | Never handle raw card data |
| **Monitoring** | Grafana + PagerDuty | Real-time alerting on financial anomalies |
| **Reconciliation** | Daily batch + real-time micro | Catch discrepancies early |

### Scalability Path

```
Phase 1 (Current):  1M txns/day, single region
  → Single PostgreSQL primary + 2 read replicas
  → Single PSP (Stripe) + backup (Adyen)
  → Redis cluster for idempotency

Phase 2 (10M txns/day):
  → Database sharding by merchant_id (4 shards)
  → Dedicated OLAP database (ClickHouse) for analytics
  → Multi-PSP intelligent routing (cost optimization)
  → Horizontal scaling of payment service (K8s)

Phase 3 (100M txns/day):
  → Multi-region active-active deployment
  → Per-region PSP routing (latency optimization)
  → Event sourcing for complete audit trail
  → Real-time ML fraud detection pipeline
  → Dedicated compliance team + automated PCI scanning

Phase 4 (1B txns/day - Stripe scale):
  → Custom payment processing infrastructure
  → Direct bank integrations (bypass PSP for top merchants)
  → In-house card network partnerships
  → Real-time cross-border settlement
  → Advanced treasury management
```

---

> **Interview Tip:** Payment system design is one of the hardest system design questions because correctness matters more than performance. Always lead with **idempotency**, **double-entry ledger**, and **state machine** — these three concepts demonstrate you understand the unique challenges of financial systems.
