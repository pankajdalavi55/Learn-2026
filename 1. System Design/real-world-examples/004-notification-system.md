# Complete System Design: Notification System (Production-Ready)

> **Complexity Level:** Intermediate to Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** Firebase Cloud Messaging (FCM), AWS SNS, Twilio, SendGrid, Apple APNs

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Delivery Pipeline](#7-deep-dive-delivery-pipeline)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a scalable notification system that can send push notifications, SMS, and emails to millions of users."

**Interviewer's Perspective:**  
This problem is a favorite at FAANG companies because it assesses:
- **Message queue design** — priority handling, ordering, backpressure
- **Delivery guarantees** — at-least-once vs exactly-once semantics
- **Fan-out architecture** — how to broadcast to millions efficiently
- **Multi-channel orchestration** — routing across push, SMS, email with different protocols
- **Reliability engineering** — retry strategies, dead letter queues, failover
- **Third-party integration** — handling unreliable external providers (APNs, FCM, Twilio, SendGrid)

> **Why this matters:** Every major application needs notifications. Designing one that is reliable, fast, and cost-effective at scale is a core infrastructure challenge.

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I dive into the design, I'd like to clarify some requirements. I'll start with functional and then move to non-functional."

### 2.1 Functional Requirements

**Candidate:** "For the functional side:
1. Which notification channels should we support?
2. Should users be able to control their preferences per channel?
3. Do we need scheduling — both immediate and delayed?
4. Should we support templated notifications, or just raw content?
5. Do we need delivery tracking and acknowledgement?
6. Should we batch/group related notifications?"

**Interviewer:** "Good questions. Here's what we need:
- Support push notifications (iOS and Android), SMS, and email
- Users can set preferences per channel (opt-in/out, quiet hours)
- Both immediate and scheduled (delayed) delivery
- Template-based content with variable substitution
- Full delivery tracking: sent, delivered, opened, failed
- Notification grouping/batching — e.g., '5 people liked your post' instead of 5 separate notifications"

**Candidate:** "Got it. Let me summarize the functional requirements:
1. ✅ Multi-channel delivery: Push (iOS/Android), SMS, Email
2. ✅ Template-based content with variable substitution
3. ✅ User notification preferences (per-channel enable/disable, quiet hours)
4. ✅ Scheduling: immediate + delayed (cron-like or one-time future)
5. ✅ Notification grouping/batching (digest mode)
6. ✅ Delivery tracking: sent → delivered → opened → clicked"

### 2.2 Non-Functional Requirements

**Candidate:** "For non-functional requirements:
1. What delivery guarantee do we need? At-most-once, at-least-once, or exactly-once?
2. What's the acceptable latency for real-time notifications?
3. What's the availability target?
4. What scale are we designing for?"

**Interviewer:**  
- "At-least-once delivery — a user might tolerate a duplicate but must not miss a critical notification."
- "Real-time notifications should be delivered within 1 minute of being triggered."
- "99.9% availability — this is critical infrastructure."
- "10 million notifications per day, with 100 million registered users."

**Candidate:** "Let me confirm the non-functional requirements:
1. ✅ At-least-once delivery guarantee (idempotent processing to reduce duplicates)
2. ✅ < 1 minute end-to-end latency for real-time notifications
3. ✅ 99.9% availability (~8.7 hours downtime/year)
4. ✅ 10M notifications/day, 100M registered users, peak ~1000 notifications/sec"

### 2.3 Out of Scope (for now)

**Candidate:** "I'll consider the following out of scope unless we have time:
- In-app notification center (can add later)
- Rich media push notifications (images, action buttons)
- Notification analytics dashboard (open rates, A/B testing)
- Cross-channel escalation (push → SMS → email)
- Multi-language / i18n support"

**Interviewer:** "That's fine. Focus on the core pipeline first; we can discuss advanced features at the end."

---

## 3. Scale Estimation

**Candidate:** "Let me estimate the scale before designing."

### 3.1 Traffic Estimation

| Metric | Value |
|--------|-------|
| Total notifications/day | 10,000,000 |
| Average notifications/sec | ~115/sec |
| Peak notifications/sec | ~1,000/sec (10x average) |
| Registered users | 100,000,000 |

**Per-Channel Breakdown:**

| Channel | % Share | Volume/Day | Avg/Sec | Peak/Sec |
|---------|---------|-----------|---------|----------|
| Push (iOS + Android) | 60% | 6,000,000 | ~70 | ~600 |
| Email | 25% | 2,500,000 | ~29 | ~250 |
| SMS | 15% | 1,500,000 | ~17 | ~150 |

### 3.2 Storage Estimation

**Candidate:** "Each notification record is approximately 500 bytes including metadata."

| Data | Size Per Record | Daily Volume | Daily Storage | Monthly |
|------|----------------|-------------|---------------|---------|
| Notification records | ~500 bytes | 10M | 5 GB | 150 GB |
| Delivery logs | ~200 bytes | 10M | 2 GB | 60 GB |
| User preferences | ~100 bytes | 100M users | 10 GB (total) | — |
| Templates | ~2 KB | ~1,000 | 2 MB (total) | — |

**Total monthly storage:** ~210 GB/month → ~2.5 TB/year

### 3.3 Bandwidth Estimation

| Direction | Calculation | Bandwidth |
|-----------|------------|-----------|
| Inbound (API calls) | 1,000 req/sec × 1 KB avg payload | ~1 MB/s |
| Outbound (to providers) | 1,000 msg/sec × 2 KB avg payload | ~2 MB/s |
| Total | | ~3 MB/s peak |

### 3.4 Cost Estimation (Monthly)

| Channel | Unit Cost | Volume/Month | Monthly Cost |
|---------|----------|-------------|-------------|
| Push (FCM/APNs) | Free | 180M | $0 |
| Email (SendGrid) | $0.00065/email | 75M | ~$48,750 |
| SMS (Twilio) | $0.0075/SMS | 45M | ~$337,500 |
| Infrastructure | — | — | ~$15,000 |
| **Total** | | | **~$401,250** |

**Candidate:** "SMS is by far the most expensive channel. This is why preference management and smart channel routing are critical — we should prefer push when possible, fall back to email, and use SMS only for critical alerts."

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│                              NOTIFICATION CLIENTS                                │
│      Backend Services  |  Admin Dashboard  |  Scheduled Jobs  |  Event Triggers  │
└─────────────────────────────────┬────────────────────────────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │      API Gateway         │  Auth, Rate Limiting, Throttling
                    │    (Kong / NGINX)        │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │  Notification Service    │  Validation, Enrichment,
                    │  (Core Orchestrator)     │  Preference Check, Routing
                    └──┬─────┬─────┬─────┬───┘
                       │     │     │     │
            ┌──────────┘     │     │     └──────────────┐
            ▼                ▼     ▼                    ▼
  ┌──────────────┐  ┌────────────────┐  ┌──────────────────────┐
  │  Template    │  │  User Pref DB  │  │  Scheduler Service   │
  │  Service     │  │  (PostgreSQL)  │  │  (Delayed / Cron)    │
  │  (Redis +    │  └────────────────┘  └──────────┬───────────┘
  │   Postgres)  │                                 │
  └──────────────┘                                 │
                                                   ▼
                    ┌─────────────────────────────────────────────┐
                    │        Priority Message Queue (Kafka)        │
                    │                                             │
                    │  ┌─────────┐ ┌─────────┐ ┌─────────┐      │
                    │  │ P0:     │ │ P1:     │ │ P2:     │      │
                    │  │Critical │ │ High    │ │ Normal  │      │
                    │  └─────────┘ └─────────┘ └─────────┘      │
                    │  ┌─────────┐                               │
                    │  │ P3: Low │                               │
                    │  └─────────┘                               │
                    └────────┬──────────┬──────────┬─────────────┘
                             │          │          │
                    ┌────────┘          │          └────────┐
                    ▼                   ▼                   ▼
          ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
          │ Push Worker  │   │ Email Worker │   │  SMS Worker  │
          │ (Consumer    │   │ (Consumer    │   │  (Consumer   │
          │  Group)      │   │  Group)      │   │   Group)     │
          └──────┬───────┘   └──────┬───────┘   └──────┬───────┘
                 │                  │                   │
        ┌────────┴────┐      ┌─────┴─────┐      ┌─────┴─────┐
        ▼             ▼      ▼           ▼      ▼           ▼
   ┌─────────┐  ┌─────────┐ ┌─────────┐  ┌──────────┐ ┌─────────┐
   │  APNs   │  │   FCM   │ │SendGrid │  │ Mailgun  │ │ Twilio  │
   │ (iOS)   │  │(Android)│ │(Primary)│  │(Fallback)│ │  (SMS)  │
   └─────────┘  └─────────┘ └─────────┘  └──────────┘ └─────────┘
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │   Delivery Tracker      │  Webhooks from providers
                    │   (Status Updates)      │  DLR, bounce, open, click
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    ▼                         ▼
          ┌──────────────┐          ┌──────────────┐
          │ Delivery Log │          │  Analytics   │
          │ (Cassandra)  │          │  (ClickHouse)│
          └──────────────┘          └──────────────┘
```

### 4.2 API Design

**Candidate:** "I'll define three core APIs."

#### Send Notification

```
POST /api/v1/notifications/send
```

```json
{
  "user_ids": ["user_123", "user_456"],
  "template_id": "order_shipped",
  "channels": ["push", "email"],
  "priority": "high",
  "variables": {
    "order_id": "ORD-789",
    "tracking_url": "https://track.example.com/ORD-789",
    "customer_name": "Alice"
  },
  "schedule_at": null,
  "group_key": "order_updates_user_123",
  "idempotency_key": "notif-order-789-shipped"
}
```

**Response:**
```json
{
  "notification_id": "ntf_abc123",
  "status": "queued",
  "recipients": 2,
  "channels": ["push", "email"],
  "created_at": "2026-04-24T10:30:00Z"
}
```

#### Get User Notifications

```
GET /api/v1/notifications/{userId}?channel=push&status=delivered&limit=20&cursor=abc
```

**Response:**
```json
{
  "notifications": [
    {
      "id": "ntf_abc123",
      "channel": "push",
      "title": "Order Shipped!",
      "body": "Your order ORD-789 has been shipped.",
      "status": "delivered",
      "created_at": "2026-04-24T10:30:00Z",
      "delivered_at": "2026-04-24T10:30:12Z"
    }
  ],
  "next_cursor": "def456",
  "has_more": true
}
```

#### Update User Preferences

```
PUT /api/v1/notifications/preferences
```

```json
{
  "user_id": "user_123",
  "preferences": {
    "push": { "enabled": true, "quiet_hours": { "start": "22:00", "end": "08:00", "timezone": "America/New_York" } },
    "email": { "enabled": true, "digest": "daily" },
    "sms": { "enabled": false }
  },
  "category_overrides": {
    "security_alerts": { "sms": true, "push": true },
    "marketing": { "email": false, "push": false }
  }
}
```

### 4.3 End-to-End Data Flow

**Candidate:** "Let me walk through the flow for sending a single notification."

```
1. Client calls POST /api/v1/notifications/send
                    │
2. API Gateway authenticates, rate-limits, forwards to Notification Service
                    │
3. Notification Service:
   a. Validates request (schema, user existence)
   b. Checks idempotency key → if duplicate, return existing notification_id
   c. Fetches user preferences from User Preference DB
   d. Filters out disabled channels (e.g., user disabled SMS)
   e. Checks quiet hours → if active, defers or drops based on priority
   f. Renders template via Template Service (variable substitution)
   g. Creates notification record in database (status: "queued")
   h. Publishes message to Kafka per channel with assigned priority
                    │
4. Kafka routes to priority topic (e.g., notifications.push.p1)
                    │
5. Channel Worker (e.g., Push Worker) consumes the message:
   a. Fetches device token(s) for user
   b. Constructs provider-specific payload (APNs / FCM format)
   c. Sends to third-party provider
   d. Updates notification status → "sent"
                    │
6. Third-party provider delivers to device
                    │
7. Provider sends delivery receipt (webhook / callback)
                    │
8. Delivery Tracker updates status → "delivered"
   Writes to delivery log (Cassandra) and analytics (ClickHouse)
```

---

## 5. Deep Dive: Core Components

### 5.1 Notification Service (Core Orchestrator)

**Candidate:** "This is the brain of the system. It handles validation, enrichment, and routing."

```javascript
// notification-service.js — Core orchestration logic

class NotificationService {
  constructor(kafkaProducer, templateService, preferenceService, deduplicator) {
    this.kafka = kafkaProducer;
    this.templates = templateService;
    this.preferences = preferenceService;
    this.dedup = deduplicator;
  }

  async send(request) {
    // 1. Idempotency check
    const existing = await this.dedup.check(request.idempotency_key);
    if (existing) return existing;

    // 2. Validate and enrich
    const users = await this.validateUsers(request.user_ids);
    const template = await this.templates.get(request.template_id);

    const results = [];

    for (const user of users) {
      // 3. Check preferences per channel
      const prefs = await this.preferences.get(user.id);
      const allowedChannels = this.filterChannels(request.channels, prefs, request.priority);

      for (const channel of allowedChannels) {
        // 4. Render template
        const content = this.templates.render(template, channel, request.variables);

        // 5. Build notification record
        const notification = {
          id: generateId(),
          user_id: user.id,
          channel,
          priority: request.priority,
          content,
          status: 'queued',
          group_key: request.group_key,
          idempotency_key: request.idempotency_key,
          created_at: new Date(),
        };

        // 6. Persist and publish
        await this.saveNotification(notification);
        await this.publishToQueue(notification);
        results.push(notification);
      }
    }

    // 7. Store idempotency key
    await this.dedup.store(request.idempotency_key, results[0]?.id);
    return { notification_id: results[0]?.id, status: 'queued', recipients: results.length };
  }

  filterChannels(requestedChannels, userPrefs, priority) {
    return requestedChannels.filter(channel => {
      const pref = userPrefs[channel];
      if (!pref?.enabled) return false;

      // Critical notifications bypass quiet hours
      if (priority === 'critical') return true;

      // Check quiet hours
      if (pref.quiet_hours && this.isInQuietHours(pref.quiet_hours)) {
        return false;
      }
      return true;
    });
  }

  async publishToQueue(notification) {
    const topic = `notifications.${notification.channel}.${notification.priority}`;
    await this.kafka.send({
      topic,
      messages: [{
        key: notification.user_id,    // partition by user_id for ordering
        value: JSON.stringify(notification),
      }],
    });
  }
}
```

### 5.2 Priority Queue Design

**Candidate:** "I'll use Kafka with separate topics per priority level. Workers consume P0 first, then P1, and so on."

```
Priority Lanes:
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  P0 (Critical)   │ OTP, security alerts, payment failures      │
│  ─────────────   │ SLA: < 10 seconds                           │
│                  │ Partitions: 12, Replication: 3               │
│                                                                 │
│  P1 (High)       │ Order updates, shipping, appointment         │
│  ─────────────   │ SLA: < 30 seconds                           │
│                  │ Partitions: 24, Replication: 3               │
│                                                                 │
│  P2 (Normal)     │ Social interactions, content updates         │
│  ─────────────   │ SLA: < 1 minute                             │
│                  │ Partitions: 48, Replication: 3               │
│                                                                 │
│  P3 (Low)        │ Marketing, digests, recommendations          │
│  ─────────────   │ SLA: < 5 minutes                            │
│                  │ Partitions: 24, Replication: 3               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

```javascript
// priority-consumer.js — Workers poll higher priority topics first

class PriorityConsumer {
  constructor(channel) {
    this.consumers = {
      p0: new KafkaConsumer({ topic: `notifications.${channel}.critical`, groupId: `${channel}-workers` }),
      p1: new KafkaConsumer({ topic: `notifications.${channel}.high`, groupId: `${channel}-workers` }),
      p2: new KafkaConsumer({ topic: `notifications.${channel}.normal`, groupId: `${channel}-workers` }),
      p3: new KafkaConsumer({ topic: `notifications.${channel}.low`, groupId: `${channel}-workers` }),
    };
  }

  async poll() {
    // Drain higher-priority queues first (weighted polling)
    const priorities = ['p0', 'p1', 'p2', 'p3'];
    const weights = [8, 4, 2, 1]; // P0 gets 8x more polling cycles

    for (let i = 0; i < priorities.length; i++) {
      const messages = await this.consumers[priorities[i]].poll({
        maxMessages: weights[i] * 10,
        timeout: 100,
      });
      if (messages.length > 0) return messages;
    }
    return [];
  }
}
```

### 5.3 Channel Workers

**Candidate:** "Each channel has its own worker fleet. Workers handle provider-specific logic and retries."

```python
# push_worker.py — Push notification worker with retry logic

import asyncio
from dataclasses import dataclass

@dataclass
class PushMessage:
    notification_id: str
    user_id: str
    title: str
    body: str
    data: dict
    priority: str

class PushWorker:
    def __init__(self, apns_client, fcm_client, token_store, delivery_tracker):
        self.apns = apns_client
        self.fcm = fcm_client
        self.token_store = token_store
        self.tracker = delivery_tracker

    async def process(self, message: PushMessage):
        tokens = await self.token_store.get_tokens(message.user_id)

        for token in tokens:
            try:
                if token.platform == 'ios':
                    await self.apns.send(token.value, message)
                elif token.platform == 'android':
                    await self.fcm.send(token.value, message)

                await self.tracker.update(message.notification_id, 'sent', token.platform)

            except InvalidTokenError:
                await self.token_store.invalidate(token)
            except ProviderRateLimitError:
                raise RetryableError("Rate limited by provider", backoff=30)
            except ProviderUnavailableError:
                raise RetryableError("Provider down", backoff=60)
```

### 5.4 Template Engine

```javascript
// template-service.js

class TemplateService {
  constructor(db, cache) {
    this.db = db;
    this.cache = cache; // Redis
  }

  async get(templateId) {
    let template = await this.cache.get(`template:${templateId}`);
    if (!template) {
      template = await this.db.query('SELECT * FROM notification_templates WHERE id = $1', [templateId]);
      await this.cache.set(`template:${templateId}`, template, { ttl: 3600 });
    }
    return template;
  }

  render(template, channel, variables) {
    const channelTemplate = template.channels[channel];

    let subject = channelTemplate.subject || '';
    let body = channelTemplate.body || '';

    // Variable substitution: {{variable_name}}
    for (const [key, value] of Object.entries(variables)) {
      const placeholder = new RegExp(`\\{\\{${key}\\}\\}`, 'g');
      subject = subject.replace(placeholder, value);
      body = body.replace(placeholder, value);
    }

    return { subject, body, metadata: channelTemplate.metadata };
  }
}

// Example template in DB:
// {
//   id: "order_shipped",
//   channels: {
//     push: { body: "Your order {{order_id}} has been shipped! Track: {{tracking_url}}" },
//     email: {
//       subject: "Your Order {{order_id}} is on its way!",
//       body: "<h1>Hi {{customer_name}},</h1><p>Your order has shipped...</p>"
//     },
//     sms: { body: "Order {{order_id}} shipped. Track at {{tracking_url}}" }
//   }
// }
```

---

## 6. Deep Dive: Database Design

### 6.1 Schema Design

**Candidate:** "I'll use PostgreSQL for transactional data and Cassandra for the write-heavy delivery log."

#### Notifications Table (PostgreSQL)

```sql
CREATE TABLE notifications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    channel         VARCHAR(10) NOT NULL CHECK (channel IN ('push', 'email', 'sms')),
    template_id     VARCHAR(100),
    priority        VARCHAR(10) NOT NULL DEFAULT 'normal'
                        CHECK (priority IN ('critical', 'high', 'normal', 'low')),
    title           TEXT,
    body            TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'queued'
                        CHECK (status IN ('queued', 'sent', 'delivered', 'opened', 'clicked', 'failed', 'bounced')),
    group_key       VARCHAR(200),
    idempotency_key VARCHAR(200) UNIQUE,
    metadata        JSONB DEFAULT '{}',
    error_message   TEXT,
    retry_count     INT DEFAULT 0,
    scheduled_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at         TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    opened_at       TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ
);

CREATE INDEX idx_notifications_user_id ON notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_status ON notifications (status) WHERE status IN ('queued', 'sent');
CREATE INDEX idx_notifications_scheduled ON notifications (scheduled_at) WHERE scheduled_at IS NOT NULL AND status = 'queued';
CREATE INDEX idx_notifications_group_key ON notifications (group_key, created_at DESC);
CREATE INDEX idx_notifications_idempotency ON notifications (idempotency_key);
```

#### User Preferences Table (PostgreSQL)

```sql
CREATE TABLE user_preferences (
    user_id             UUID NOT NULL,
    channel             VARCHAR(10) NOT NULL CHECK (channel IN ('push', 'email', 'sms')),
    enabled             BOOLEAN NOT NULL DEFAULT TRUE,
    quiet_hours_start   TIME,
    quiet_hours_end     TIME,
    timezone            VARCHAR(50) DEFAULT 'UTC',
    digest_frequency    VARCHAR(20) DEFAULT 'none'
                            CHECK (digest_frequency IN ('none', 'hourly', 'daily', 'weekly')),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, channel)
);

CREATE TABLE user_category_preferences (
    user_id     UUID NOT NULL,
    category    VARCHAR(50) NOT NULL,
    channel     VARCHAR(10) NOT NULL,
    enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (user_id, category, channel)
);
```

#### Notification Templates Table (PostgreSQL)

```sql
CREATE TABLE notification_templates (
    id              VARCHAR(100) PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    category        VARCHAR(50) NOT NULL,
    channel         VARCHAR(10) NOT NULL,
    subject         TEXT,
    body            TEXT NOT NULL,
    variables       JSONB NOT NULL DEFAULT '[]',
    metadata        JSONB DEFAULT '{}',
    is_active       BOOLEAN DEFAULT TRUE,
    version         INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Example: list required variables for validation
-- variables: ["order_id", "customer_name", "tracking_url"]
```

#### Device Tokens Table (PostgreSQL)

```sql
CREATE TABLE device_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    platform    VARCHAR(10) NOT NULL CHECK (platform IN ('ios', 'android', 'web')),
    token       TEXT NOT NULL UNIQUE,
    app_version VARCHAR(20),
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_tokens_user ON device_tokens (user_id) WHERE is_active = TRUE;
```

### 6.2 Delivery Log (Cassandra)

**Candidate:** "The delivery log is write-heavy — every send, delivery receipt, and status change writes here. Cassandra is ideal for this workload."

```sql
-- Cassandra: optimized for high-volume time-series writes
CREATE TABLE delivery_log (
    notification_id UUID,
    event_type      TEXT,        -- 'queued', 'sent', 'delivered', 'opened', 'failed', 'bounced'
    channel         TEXT,
    provider        TEXT,        -- 'apns', 'fcm', 'sendgrid', 'twilio'
    event_time      TIMESTAMP,
    metadata        MAP<TEXT, TEXT>,
    PRIMARY KEY (notification_id, event_time)
) WITH CLUSTERING ORDER BY (event_time DESC)
  AND default_time_to_live = 7776000;  -- 90 days TTL

-- Query pattern: "Get all events for a notification"
-- SELECT * FROM delivery_log WHERE notification_id = ?;

-- Secondary table for user-level queries
CREATE TABLE user_delivery_log (
    user_id         UUID,
    day             DATE,
    event_time      TIMESTAMP,
    notification_id UUID,
    channel         TEXT,
    event_type      TEXT,
    PRIMARY KEY ((user_id, day), event_time)
) WITH CLUSTERING ORDER BY (event_time DESC)
  AND default_time_to_live = 7776000;
```

### 6.3 When to Use Which Database

| Data | Database | Reasoning |
|------|----------|-----------|
| Notifications, Templates, Preferences | PostgreSQL | Transactional, relational queries, moderate write volume |
| Delivery Logs | Cassandra | Write-heavy (10M+ writes/day), time-series, append-only |
| Template Cache | Redis | Low-latency reads, TTL-based invalidation |
| Idempotency Keys | Redis | Fast lookup, auto-expire after 24 hours |
| Analytics | ClickHouse | Column-oriented, fast aggregation over billions of rows |

---

## 7. Deep Dive: Delivery Pipeline

> **This is the KEY SECTION — interviewers spend the most time here.**

### 7.1 Fan-Out Strategies

**Candidate:** "There are two approaches for broadcasting notifications."

```
Strategy 1: Per-User Fan-Out (Write-time fan-out)
─────────────────────────────────────────────────
Trigger: "Send promo to all 100M users"
                    │
                    ▼
    ┌──────────────────────────┐
    │  Fan-Out Service         │
    │  Generates 100M messages │ ← SLOW: minutes to hours
    │  (one per user)          │
    └──────────────────────────┘
                    │
                    ▼
    ┌──────────────────────────┐
    │  Kafka (100M messages)   │ ← HIGH storage cost
    └──────────────────────────┘

Pros: Simple delivery logic, per-user preferences applied at fan-out
Cons: Slow for large audiences, Kafka storage spike

Strategy 2: Per-Topic Fan-Out (Read-time fan-out)
─────────────────────────────────────────────────
Trigger: "Send promo to topic 'all_users'"
                    │
                    ▼
    ┌──────────────────────────┐
    │  1 message published     │ ← FAST: instant
    │  with audience = topic   │
    └──────────────────────────┘
                    │
                    ▼
    ┌──────────────────────────┐
    │  Worker resolves topic   │
    │  → fetches user list     │ ← Paginated, streaming
    │  → sends in batches      │
    └──────────────────────────┘

Pros: Fast publish, low queue storage
Cons: Complex worker logic, harder to track per-user status
```

**Candidate:** "I'd use a hybrid approach:
- **Per-user fan-out** for targeted notifications (< 10,000 recipients)
- **Per-topic fan-out** for broadcasts (> 10,000 recipients), using batch streaming workers"

### 7.2 Push Notification Flow

```
┌────────────────────────────────────────────────────────────┐
│                 Push Notification Pipeline                   │
│                                                            │
│  1. Worker receives message from Kafka                     │
│              │                                             │
│  2. Fetch device tokens for user_id                        │
│              │                                             │
│     ┌────────┴────────┐                                    │
│     ▼                 ▼                                    │
│  iOS tokens       Android tokens                           │
│     │                 │                                    │
│  3. Build APNs     3. Build FCM                            │
│     payload           payload                              │
│     │                 │                                    │
│  4. Send via       4. Send via                             │
│     HTTP/2            HTTP/1.1                              │
│     to APNs           to FCM                               │
│     │                 │                                    │
│  5. Handle response:                                       │
│     ├─ 200 OK → mark "sent"                                │
│     ├─ 410 Gone → invalidate token                         │
│     ├─ 429 Too Many → backoff + retry                      │
│     └─ 503 Unavailable → retry with exponential backoff    │
│                                                            │
│  6. APNs/FCM delivers to device                            │
│  7. Device sends delivery receipt (if supported)           │
└────────────────────────────────────────────────────────────┘
```

#### Token Management

```javascript
// token-manager.js

class TokenManager {
  constructor(db, cache) {
    this.db = db;
    this.cache = cache;
  }

  async getActiveTokens(userId) {
    const cacheKey = `tokens:${userId}`;
    let tokens = await this.cache.get(cacheKey);

    if (!tokens) {
      tokens = await this.db.query(
        'SELECT * FROM device_tokens WHERE user_id = $1 AND is_active = TRUE',
        [userId]
      );
      await this.cache.set(cacheKey, tokens, { ttl: 300 }); // 5 min cache
    }
    return tokens;
  }

  async registerToken(userId, platform, token, appVersion) {
    await this.db.query(`
      INSERT INTO device_tokens (user_id, platform, token, app_version)
      VALUES ($1, $2, $3, $4)
      ON CONFLICT (token)
      DO UPDATE SET user_id = $1, app_version = $4, is_active = TRUE, last_used = NOW()
    `, [userId, platform, token, appVersion]);

    await this.cache.del(`tokens:${userId}`);
  }

  async invalidateToken(token) {
    await this.db.query(
      'UPDATE device_tokens SET is_active = FALSE WHERE token = $1 RETURNING user_id',
      [token]
    );
  }

  // Periodic job: clean up stale tokens (no activity in 90 days)
  async cleanupStaleTokens() {
    const result = await this.db.query(`
      UPDATE device_tokens
      SET is_active = FALSE
      WHERE is_active = TRUE AND last_used < NOW() - INTERVAL '90 days'
      RETURNING user_id, token
    `);
    return result.rows.length;
  }
}
```

### 7.3 Email Delivery Pipeline

```
Email Pipeline:
┌───────────────────────────────────────────────────────────────┐
│                                                               │
│  1. Email Worker receives message                             │
│              │                                                │
│  2. Render HTML template (MJML → HTML)                        │
│     ├─ Inline CSS (email client compatibility)                │
│     ├─ Add tracking pixel (open tracking)                     │
│     └─ Rewrite links (click tracking)                         │
│              │                                                │
│  3. SMTP Pipeline:                                            │
│     ├─ Check sender reputation score                          │
│     ├─ Apply rate limiting per domain                         │
│     │   (Gmail: 2,000/day per sender, Yahoo: 500/hr)          │
│     ├─ Select sending IP from pool (warm IPs)                 │
│     └─ Send via SendGrid API                                  │
│              │                                                │
│  4. Handle webhooks from SendGrid:                            │
│     ├─ "delivered"  → update status                           │
│     ├─ "bounced"    → mark email invalid, suppress future     │
│     ├─ "complained" → add to suppression list                 │
│     ├─ "opened"     → tracking pixel loaded                   │
│     └─ "clicked"    → link redirect captured                  │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

#### Bounce Handling and Suppression

```python
# email_worker.py

class EmailWorker:
    def __init__(self, sendgrid_client, suppression_list, delivery_tracker):
        self.sendgrid = sendgrid_client
        self.suppression = suppression_list
        self.tracker = delivery_tracker

    async def process(self, notification):
        email = notification.user_email

        # Check suppression list (bounced or complained)
        if await self.suppression.is_suppressed(email):
            await self.tracker.update(notification.id, 'skipped', reason='suppressed')
            return

        # Check domain rate limits
        domain = email.split('@')[1]
        if not await self.rate_limiter.acquire(f'email_domain:{domain}'):
            raise RetryableError("Domain rate limit hit", backoff=60)

        response = await self.sendgrid.send(
            to=email,
            subject=notification.content.subject,
            html=notification.content.body,
            headers={
                'X-Notification-ID': notification.id,
                'List-Unsubscribe': f'<https://example.com/unsubscribe/{notification.user_id}>',
            }
        )

        await self.tracker.update(notification.id, 'sent', provider='sendgrid',
                                  provider_id=response.message_id)

    async def handle_bounce_webhook(self, event):
        """Called by webhook endpoint when SendGrid reports a bounce."""
        if event.type == 'hard_bounce':
            await self.suppression.add(event.email, reason='hard_bounce')
        elif event.type == 'spam_complaint':
            await self.suppression.add(event.email, reason='complaint')
            # Complaint rate > 0.1% triggers alert
            await self.check_complaint_rate()
```

#### IP Warm-Up Strategy

```
New IP Warm-Up Schedule:
┌──────────┬────────────────┬──────────────────┐
│   Day    │  Emails/Day    │  Notes           │
├──────────┼────────────────┼──────────────────┤
│  Day 1   │      50        │  Engaged users   │
│  Day 2   │     100        │  only            │
│  Day 3   │     500        │                  │
│  Day 7   │   5,000        │  Monitor bounce  │
│  Day 14  │  25,000        │  rate closely    │
│  Day 21  │ 100,000        │                  │
│  Day 30  │ 500,000        │  Full capacity   │
└──────────┴────────────────┴──────────────────┘
Rule: If bounce rate > 5% or complaint rate > 0.1%, pause and investigate.
```

### 7.4 SMS Delivery Pipeline

```python
# sms_worker.py

class SMSWorker:
    def __init__(self, twilio_client, delivery_tracker):
        self.twilio = twilio_client
        self.tracker = delivery_tracker

    async def process(self, notification):
        phone = notification.user_phone

        # Format to E.164 international standard
        formatted_phone = self.format_e164(phone, notification.user_country)

        # Carrier routing: select optimal sender ID/number
        sender = self.select_sender(formatted_phone)

        # Character limit: SMS = 160 chars (GSM-7) or 70 chars (UCS-2/Unicode)
        body = self.truncate_body(notification.content.body, encoding='gsm7')

        response = await self.twilio.send(
            to=formatted_phone,
            from_=sender,
            body=body,
            status_callback=f'https://api.example.com/webhooks/sms/status/{notification.id}'
        )

        await self.tracker.update(notification.id, 'sent', provider='twilio',
                                  provider_id=response.sid)

    def format_e164(self, phone, country_code):
        """Convert local number to E.164 format."""
        # +1 (US), +91 (India), +44 (UK), etc.
        if phone.startswith('+'):
            return phone
        country_prefixes = {'US': '+1', 'IN': '+91', 'UK': '+44', 'DE': '+49'}
        prefix = country_prefixes.get(country_code, '')
        return f"{prefix}{phone.lstrip('0')}"

    async def handle_dlr(self, notification_id, status):
        """Handle Delivery Receipt (DLR) from Twilio."""
        status_map = {
            'delivered': 'delivered',
            'undelivered': 'failed',
            'failed': 'failed',
        }
        mapped_status = status_map.get(status, 'unknown')
        await self.tracker.update(notification_id, mapped_status, channel='sms')
```

### 7.5 Retry Strategy with Exponential Backoff

**Candidate:** "Retries are critical for at-least-once delivery. I'll use exponential backoff with jitter."

```javascript
// retry-handler.js

class RetryHandler {
  constructor(maxRetries = 5, baseDelay = 1000, maxDelay = 300000) {
    this.maxRetries = maxRetries;
    this.baseDelay = baseDelay;   // 1 second
    this.maxDelay = maxDelay;     // 5 minutes
  }

  async executeWithRetry(fn, context) {
    let lastError;

    for (let attempt = 0; attempt <= this.maxRetries; attempt++) {
      try {
        return await fn();
      } catch (error) {
        lastError = error;

        if (!this.isRetryable(error)) {
          throw error; // Non-retryable: bad request, invalid token, etc.
        }

        if (attempt === this.maxRetries) {
          break; // Exhausted retries
        }

        const delay = this.calculateDelay(attempt);
        console.log(`Retry ${attempt + 1}/${this.maxRetries} for ${context.notificationId}, ` +
                     `waiting ${delay}ms. Error: ${error.message}`);

        await this.sleep(delay);
      }
    }

    // All retries exhausted → send to Dead Letter Queue
    await this.sendToDeadLetterQueue(context, lastError);
    throw lastError;
  }

  calculateDelay(attempt) {
    // Exponential backoff: 1s, 2s, 4s, 8s, 16s ... capped at maxDelay
    const exponentialDelay = this.baseDelay * Math.pow(2, attempt);
    const cappedDelay = Math.min(exponentialDelay, this.maxDelay);

    // Add jitter (±25%) to prevent thundering herd
    const jitter = cappedDelay * 0.25 * (Math.random() * 2 - 1);
    return Math.floor(cappedDelay + jitter);
  }

  isRetryable(error) {
    const retryableCodes = [429, 500, 502, 503, 504];
    return retryableCodes.includes(error.statusCode) ||
           error.code === 'ECONNRESET' ||
           error.code === 'ETIMEDOUT';
  }

  async sendToDeadLetterQueue(context, error) {
    await this.kafka.send({
      topic: 'notifications.dead_letter',
      messages: [{
        key: context.notificationId,
        value: JSON.stringify({
          notification: context,
          error: { message: error.message, code: error.statusCode },
          exhausted_at: new Date().toISOString(),
          total_attempts: this.maxRetries + 1,
        }),
      }],
    });
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

// Retry schedule visualization:
//
// Attempt 0: immediate
// Attempt 1: ~1s    (1000ms ± 250ms)
// Attempt 2: ~2s    (2000ms ± 500ms)
// Attempt 3: ~4s    (4000ms ± 1000ms)
// Attempt 4: ~8s    (8000ms ± 2000ms)
// Attempt 5: → Dead Letter Queue
//
// Total max wait: ~15 seconds before DLQ
```

### 7.6 Dead Letter Queue (DLQ)

```
Dead Letter Queue Flow:
┌──────────────┐     ┌──────────────┐     ┌──────────────────────┐
│ Worker fails │────▶│  DLQ Topic   │────▶│  DLQ Consumer        │
│ all retries  │     │  (Kafka)     │     │  (low-priority job)  │
└──────────────┘     └──────────────┘     └──────────┬───────────┘
                                                     │
                                          ┌──────────┴───────────┐
                                          │  Classify failure:   │
                                          │  ├─ Provider outage  │──▶ Re-queue after 1 hour
                                          │  ├─ Invalid data     │──▶ Alert + log
                                          │  ├─ Rate limited     │──▶ Re-queue with backoff
                                          │  └─ Unknown          │──▶ Manual review queue
                                          └──────────────────────┘
```

### 7.7 Deduplication

**Candidate:** "To prevent duplicate notifications, I use a three-layer deduplication strategy."

```javascript
// deduplication.js

class DeduplicationService {
  constructor(redis) {
    this.redis = redis;
  }

  // Layer 1: API-level — idempotency key from the caller
  async checkIdempotencyKey(key) {
    if (!key) return null;
    const existing = await this.redis.get(`idemp:${key}`);
    return existing ? JSON.parse(existing) : null;
  }

  async storeIdempotencyKey(key, notificationId, ttl = 86400) {
    await this.redis.setex(`idemp:${key}`, ttl, JSON.stringify({ notification_id: notificationId }));
  }

  // Layer 2: Content-based — same user + channel + content hash within time window
  async isDuplicate(userId, channel, contentHash) {
    const dedupKey = `dedup:${userId}:${channel}:${contentHash}`;
    const exists = await this.redis.exists(dedupKey);
    if (exists) return true;

    // Mark as seen, expire in 1 hour (configurable per category)
    await this.redis.setex(dedupKey, 3600, '1');
    return false;
  }

  // Layer 3: Worker-level — Kafka consumer offset tracking
  // Kafka guarantees at-least-once delivery; workers must be idempotent.
  // Before processing, check if notification was already sent:
  async isAlreadySent(notificationId) {
    const status = await this.redis.get(`sent:${notificationId}`);
    return status === 'true';
  }

  async markSent(notificationId, ttl = 86400) {
    await this.redis.setex(`sent:${notificationId}`, ttl, 'true');
  }
}
```

---

## 8. Scaling Strategies

### 8.1 Kafka Partitioning

**Candidate:** "I'll partition Kafka topics by user_id to guarantee per-user ordering."

```
Kafka Topic: notifications.push.high
┌────────────────────────────────────────────────────┐
│                                                    │
│  Partition 0:  user_001, user_049, user_097 ...    │
│  Partition 1:  user_002, user_050, user_098 ...    │
│  Partition 2:  user_003, user_051, user_099 ...    │
│      ...                                           │
│  Partition 47: user_048, user_096, user_144 ...    │
│                                                    │
│  Key: hash(user_id) % num_partitions               │
│  Guarantee: all notifications for a user go to     │
│             the same partition → ordered delivery   │
│                                                    │
└────────────────────────────────────────────────────┘
```

### 8.2 Consumer Groups Per Channel

```
Consumer Group Architecture:
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  Consumer Group: push-workers                                │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │Worker-1 │ │Worker-2 │ │Worker-3 │ │Worker-4 │  ← Auto- │
│  │P0,P1,P2 │ │P3,P4,P5 │ │P6,P7,P8 │ │P9,P10   │  scale  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘          │
│                                                              │
│  Consumer Group: email-workers                               │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                       │
│  │Worker-1 │ │Worker-2 │ │Worker-3 │  ← Fewer workers      │
│  │P0–P3    │ │P4–P7    │ │P8–P11   │    (email is slower)  │
│  └─────────┘ └─────────┘ └─────────┘                       │
│                                                              │
│  Consumer Group: sms-workers                                 │
│  ┌─────────┐ ┌─────────┐                                   │
│  │Worker-1 │ │Worker-2 │  ← Fewest (lowest volume,         │
│  │P0–P5    │ │P6–P11   │    highest cost per message)       │
│  └─────────┘ └─────────┘                                   │
│                                                              │
└──────────────────────────────────────────────────────────────┘

Auto-Scaling Rules:
- Scale up when: Kafka consumer lag > 10,000 messages for > 2 minutes
- Scale down when: Kafka consumer lag < 100 for > 10 minutes
- Max workers: push=20, email=12, sms=6
```

### 8.3 Horizontal Scaling of Workers

```python
# Kubernetes HPA (Horizontal Pod Autoscaler) config concept

# Push workers: scale on Kafka consumer lag
hpa_push_workers = {
    "min_replicas": 4,
    "max_replicas": 20,
    "metrics": [
        {"type": "kafka_consumer_lag", "target": 5000},
        {"type": "cpu", "target_utilization": 70},
    ],
}

# Email workers: scale on queue depth + respect rate limits
hpa_email_workers = {
    "min_replicas": 3,
    "max_replicas": 12,
    "metrics": [
        {"type": "kafka_consumer_lag", "target": 10000},
    ],
    "constraints": {
        "max_send_rate_per_worker": 100,  # emails/sec per worker
    },
}
```

### 8.4 Database Sharding

```
Notification DB Sharding by user_id:
┌─────────────────────────────────────────────────────┐
│                                                     │
│  Shard 0: user_id hash % 16 == 0                    │
│  Shard 1: user_id hash % 16 == 1                    │
│  ...                                                │
│  Shard 15: user_id hash % 16 == 15                  │
│                                                     │
│  Benefits:                                          │
│  - All notifications for a user are co-located      │
│  - "Get my notifications" query hits single shard   │
│  - Write load distributed across 16 instances       │
│                                                     │
│  Cassandra (delivery log): automatic sharding via   │
│  consistent hashing on partition key                 │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 8.5 CDN for Email Assets

```
Email Image Hosting:
┌───────────┐     ┌──────────┐     ┌──────────────┐
│ Email     │────▶│ CloudFront│────▶│ S3 Bucket    │
│ Client    │     │ (CDN)    │     │ (img assets) │
│ loads     │     └──────────┘     └──────────────┘
│ images    │
└───────────┘

- All email images served via CDN (low latency, high availability)
- Tracking pixel also served via CDN → redirect to tracking endpoint
- URL format: https://cdn.example.com/emails/{campaign_id}/{image}.png
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Third-Party Provider Outage

**Scenario:** SendGrid is down; emails are not being delivered.

```
Failover Strategy:
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  Email Worker:                                               │
│  ┌─────────────┐                                             │
│  │ Send via    │──── Success ──▶ Done                        │
│  │ SendGrid    │                                             │
│  │ (primary)   │──── Fail (3x) ──▶ ┌──────────────────┐     │
│  └─────────────┘                   │ Circuit Breaker   │     │
│                                    │ OPENS             │     │
│                                    └────────┬──────────┘     │
│                                             │                │
│                                    ┌────────▼──────────┐     │
│                                    │ Send via Mailgun  │     │
│                                    │ (secondary)       │     │
│                                    └──────────────────┘     │
│                                                              │
│  Circuit Breaker Config:                                     │
│  - Failure threshold: 5 failures in 30 seconds               │
│  - Open duration: 60 seconds (then half-open)                │
│  - Half-open: allow 1 request; if success → close            │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

```javascript
// circuit-breaker.js

class CircuitBreaker {
  constructor(options = {}) {
    this.failureThreshold = options.failureThreshold || 5;
    this.resetTimeout = options.resetTimeout || 60000;
    this.state = 'CLOSED';  // CLOSED → OPEN → HALF_OPEN → CLOSED
    this.failureCount = 0;
    this.lastFailureTime = null;
  }

  async execute(primaryFn, fallbackFn) {
    if (this.state === 'OPEN') {
      if (Date.now() - this.lastFailureTime > this.resetTimeout) {
        this.state = 'HALF_OPEN';
      } else {
        return fallbackFn();
      }
    }

    try {
      const result = await primaryFn();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      if (this.state === 'OPEN') {
        return fallbackFn();
      }
      throw error;
    }
  }

  onSuccess() {
    this.failureCount = 0;
    this.state = 'CLOSED';
  }

  onFailure() {
    this.failureCount++;
    this.lastFailureTime = Date.now();
    if (this.failureCount >= this.failureThreshold) {
      this.state = 'OPEN';
    }
  }
}
```

### 9.2 Kafka Broker Failure

| Scenario | Mitigation |
|----------|------------|
| Single broker down | Replication factor = 3; remaining ISR replicas serve reads/writes |
| Partition leader failure | Controller auto-elects new leader from ISR (~seconds) |
| All brokers in a partition's ISR down | Producers buffer locally; `acks=all` ensures no data loss on commit |
| Full cluster outage | Producers queue locally (bounded buffer); alert triggers manual failover to standby cluster |

### 9.3 Duplicate Notification Delivery

**Cause:** Kafka consumer crashes after sending but before committing offset.

```
Prevention Layers:
1. Idempotency key (API level) — prevents re-submission
2. Content dedup hash (service level) — prevents same content within time window
3. "Already sent" check (worker level) — checks Redis before sending to provider
4. Provider-level dedup — APNs and FCM collapse matching collapse_key notifications
```

### 9.4 User Preference Update Race Condition

**Scenario:** User disables push notifications while a notification is in-flight.

```
Timeline:
  T=0: Notification enters Kafka queue (push channel)
  T=1: User disables push notifications (preference updated)
  T=2: Push worker picks up message, fetches stale preference (cache)

Solution:
  - Preference cache TTL = 60 seconds (not hours)
  - For critical preference changes (unsubscribe), also publish
    invalidation event to Kafka → workers listen for pref changes
  - Workers always re-check preferences for delayed/scheduled notifications
  - Accept eventual consistency: occasional delivery after opt-out is
    tolerable; hard unsubscribes are handled at the provider level
    (APNs token revocation, email unsubscribe header)
```

### 9.5 Quiet Hours Timezone Handling

```python
# quiet_hours.py

from datetime import datetime
import pytz

def is_in_quiet_hours(user_prefs, notification_priority):
    # Critical notifications always bypass quiet hours
    if notification_priority == 'critical':
        return False

    if not user_prefs.quiet_hours_start or not user_prefs.quiet_hours_end:
        return False

    user_tz = pytz.timezone(user_prefs.timezone or 'UTC')
    now_user_local = datetime.now(pytz.utc).astimezone(user_tz).time()

    start = user_prefs.quiet_hours_start  # e.g., time(22, 0)
    end = user_prefs.quiet_hours_end      # e.g., time(8, 0)

    # Handle overnight quiet hours (22:00 → 08:00)
    if start > end:
        return now_user_local >= start or now_user_local < end
    else:
        return start <= now_user_local < end

# Behavior during quiet hours:
# - P0 (critical): deliver immediately (OTP, security alerts)
# - P1 (high): queue and deliver when quiet hours end
# - P2/P3: queue or batch into morning digest
```

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Notification System Metrics                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  DELIVERY METRICS (per channel):                                        │
│  ├─ delivery_rate          — % of sent notifications that were delivered│
│  ├─ delivery_latency_p50   — median time from trigger to delivery       │
│  ├─ delivery_latency_p99   — 99th percentile delivery latency           │
│  ├─ send_rate              — notifications sent per second               │
│  └─ failure_rate           — % of notifications that failed              │
│                                                                         │
│  QUEUE METRICS:                                                         │
│  ├─ queue_depth_p0         — messages waiting in critical queue          │
│  ├─ queue_depth_p1/p2/p3   — messages in other priority queues           │
│  ├─ consumer_lag           — how far behind consumers are                │
│  └─ dlq_size               — dead letter queue message count             │
│                                                                         │
│  EMAIL-SPECIFIC:                                                        │
│  ├─ bounce_rate            — hard bounces / total sent (target < 2%)    │
│  ├─ complaint_rate         — spam complaints / total sent (< 0.1%)     │
│  └─ open_rate              — unique opens / total delivered              │
│                                                                         │
│  PUSH-SPECIFIC:                                                         │
│  ├─ invalid_token_rate     — % of sends hitting invalid tokens           │
│  └─ platform_split         — iOS vs Android delivery ratio               │
│                                                                         │
│  SMS-SPECIFIC:                                                          │
│  ├─ dlr_success_rate       — delivery receipts confirming delivery       │
│  └─ cost_per_message       — tracked per carrier/country                 │
│                                                                         │
│  SYSTEM HEALTH:                                                         │
│  ├─ api_latency_p99        — notification API response time              │
│  ├─ worker_cpu_utilization — per worker group                            │
│  └─ provider_error_rate    — errors from APNs, FCM, SendGrid, Twilio    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Prometheus Metrics (Code)

```python
# metrics.py — Prometheus instrumentation

from prometheus_client import Counter, Histogram, Gauge

# Delivery counters
notifications_sent = Counter(
    'notifications_sent_total',
    'Total notifications sent',
    ['channel', 'priority', 'status']  # status: success, failed, retried
)

notifications_delivered = Counter(
    'notifications_delivered_total',
    'Total notifications confirmed delivered',
    ['channel', 'provider']
)

# Latency histograms
delivery_latency = Histogram(
    'notification_delivery_latency_seconds',
    'Time from creation to delivery',
    ['channel', 'priority'],
    buckets=[0.5, 1, 2, 5, 10, 30, 60, 120, 300]
)

# Queue gauges
queue_depth = Gauge(
    'notification_queue_depth',
    'Current messages in queue',
    ['channel', 'priority']
)

consumer_lag = Gauge(
    'notification_consumer_lag',
    'Kafka consumer lag',
    ['consumer_group', 'topic']
)

# Provider health
provider_errors = Counter(
    'notification_provider_errors_total',
    'Errors from third-party providers',
    ['provider', 'error_type']
)
```

### 10.3 Alerting Rules

```yaml
# Prometheus alert rules

groups:
  - name: notification_alerts
    rules:
      - alert: HighDeliveryFailureRate
        expr: |
          rate(notifications_sent_total{status="failed"}[5m])
          / rate(notifications_sent_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Notification failure rate > 5% for {{ $labels.channel }}"

      - alert: CriticalQueueBacklog
        expr: notification_queue_depth{priority="critical"} > 1000
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Critical priority queue depth > 1000"

      - alert: HighConsumerLag
        expr: notification_consumer_lag > 50000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Consumer group {{ $labels.consumer_group }} lag > 50K"

      - alert: EmailBounceRateHigh
        expr: |
          rate(notifications_delivered_total{channel="email", status="bounced"}[1h])
          / rate(notifications_sent_total{channel="email"}[1h]) > 0.02
        for: 15m
        labels:
          severity: warning
        annotations:
          summary: "Email bounce rate > 2%"

      - alert: ProviderDown
        expr: |
          rate(notification_provider_errors_total[5m]) > 10
        for: 3m
        labels:
          severity: critical
        annotations:
          summary: "Provider {{ $labels.provider }} error rate spike"

      - alert: DLQGrowing
        expr: notification_queue_depth{priority="dead_letter"} > 500
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Dead letter queue has > 500 unprocessed messages"
```

### 10.4 Structured Logging

```javascript
// Structured log format for notification lifecycle

// On send:
{
  "timestamp": "2026-04-24T10:30:00.123Z",
  "level": "info",
  "event": "notification.sent",
  "notification_id": "ntf_abc123",
  "user_id": "user_123",
  "channel": "push",
  "provider": "fcm",
  "priority": "high",
  "template_id": "order_shipped",
  "latency_ms": 145,
  "attempt": 1,
  "trace_id": "trace-xyz-789"
}

// On failure:
{
  "timestamp": "2026-04-24T10:30:02.456Z",
  "level": "error",
  "event": "notification.failed",
  "notification_id": "ntf_abc123",
  "user_id": "user_123",
  "channel": "push",
  "provider": "apns",
  "error_code": "InvalidToken",
  "error_message": "Device token is no longer active",
  "attempt": 3,
  "will_retry": false,
  "trace_id": "trace-xyz-789"
}
```

### 10.5 Grafana Dashboard Layout

```
┌──────────────────────────────────────────────────────────────────────┐
│                  Notification System Dashboard                       │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Row 1: Overview                                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────────┐  │
│  │ Sent Today  │ │ Delivered   │ │ Failed      │ │ DLQ Size     │  │
│  │ 8.2M       │ │ 97.3%      │ │ 0.8%       │ │ 42          │  │
│  └─────────────┘ └─────────────┘ └─────────────┘ └──────────────┘  │
│                                                                      │
│  Row 2: Throughput (time series graph)                               │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  📈 Notifications/sec by channel (push | email | sms)        │   │
│  │  [stacked area chart, 24h window]                            │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  Row 3: Latency + Queue Depth                                        │
│  ┌───────────────────────────┐ ┌────────────────────────────────┐   │
│  │ Delivery Latency (p50/99)│ │ Queue Depth by Priority        │   │
│  │ [line chart per channel] │ │ [bar chart: P0, P1, P2, P3]   │   │
│  └───────────────────────────┘ └────────────────────────────────┘   │
│                                                                      │
│  Row 4: Provider Health                                              │
│  ┌───────────────────────────┐ ┌────────────────────────────────┐   │
│  │ Provider Error Rates     │ │ Email: Bounce + Complaint Rate │   │
│  │ [heatmap: APNs, FCM,    │ │ [line chart, alert thresholds] │   │
│  │  SendGrid, Twilio]       │ │                                │   │
│  └───────────────────────────┘ └────────────────────────────────┘   │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 11. Advanced Features

### 11.1 Notification Aggregation / Digest

**Candidate:** "Instead of sending '5 people liked your post' as 5 separate notifications, I'll aggregate them."

```javascript
// aggregation-service.js

class NotificationAggregator {
  constructor(redis, scheduler) {
    this.redis = redis;
    this.scheduler = scheduler;
    this.WINDOW = 300; // 5-minute aggregation window
  }

  async aggregate(notification) {
    const groupKey = notification.group_key;
    if (!groupKey) return null; // Not eligible for aggregation

    const bucketKey = `agg:${groupKey}`;
    const count = await this.redis.incr(bucketKey);

    if (count === 1) {
      // First notification in this group — start the aggregation window
      await this.redis.expire(bucketKey, this.WINDOW);
      await this.redis.set(`agg:first:${groupKey}`, JSON.stringify(notification));

      // Schedule digest delivery at end of window
      this.scheduler.scheduleOnce(`deliver_digest:${groupKey}`, this.WINDOW * 1000);
      return 'aggregating';
    }

    // Subsequent notifications — append actor to list
    await this.redis.rpush(`agg:actors:${groupKey}`, JSON.stringify({
      actor_name: notification.variables.actor_name,
      timestamp: new Date().toISOString(),
    }));

    return 'aggregated';
  }

  async deliverDigest(groupKey) {
    const count = await this.redis.get(`agg:${groupKey}`);
    const first = JSON.parse(await this.redis.get(`agg:first:${groupKey}`));
    const actors = await this.redis.lrange(`agg:actors:${groupKey}`, 0, -1);

    // Build aggregated notification
    // "Alice, Bob, and 3 others liked your post"
    const actorNames = actors.map(a => JSON.parse(a).actor_name);
    const displayNames = this.formatActorList(actorNames, parseInt(count));

    const digestNotification = {
      ...first,
      body: `${displayNames} liked your post`,
      metadata: { aggregated_count: parseInt(count) },
    };

    // Clean up Redis keys
    await this.redis.del(`agg:${groupKey}`, `agg:first:${groupKey}`, `agg:actors:${groupKey}`);

    return digestNotification;
  }

  formatActorList(names, totalCount) {
    if (totalCount === 1) return names[0];
    if (totalCount === 2) return `${names[0]} and ${names[1]}`;
    const othersCount = totalCount - 2;
    return `${names[0]}, ${names[1]}, and ${othersCount} other${othersCount > 1 ? 's' : ''}`;
  }
}
```

### 11.2 A/B Testing Notification Content

```
A/B Test Flow:
┌──────────────────┐     ┌────────────────┐     ┌─────────────────┐
│ Notification Send │────▶│ A/B Test       │────▶│ Variant A (50%) │
│ (template has     │     │ Assignment     │     │ "Your order     │
│  active test)     │     │ (hash user_id) │     │  has shipped!"  │
│                   │     │                │     ├─────────────────┤
│                   │     │                │────▶│ Variant B (50%) │
│                   │     │                │     │ "Great news!    │
│                   │     └────────────────┘     │  Order on its   │
│                   │                            │  way 🚀"        │
└──────────────────┘                            └─────────────────┘
                                                        │
                                                        ▼
                                                ┌─────────────────┐
                                                │ Track: open     │
                                                │ rate, click     │
                                                │ rate per variant│
                                                └─────────────────┘
```

### 11.3 Rich Push Notifications

```json
{
  "to": "device_token_abc",
  "notification": {
    "title": "Your order has shipped!",
    "body": "Arriving Thursday, April 24",
    "image": "https://cdn.example.com/orders/ORD-789/preview.jpg"
  },
  "data": {
    "deeplink": "myapp://orders/ORD-789",
    "notification_id": "ntf_abc123"
  },
  "android": {
    "notification": {
      "click_action": "OPEN_ORDER_DETAIL",
      "channel_id": "order_updates"
    }
  },
  "apns": {
    "payload": {
      "aps": {
        "mutable-content": 1,
        "category": "ORDER_UPDATE"
      }
    },
    "fcm_options": {
      "image": "https://cdn.example.com/orders/ORD-789/preview.jpg"
    }
  }
}
```

### 11.4 In-App Notification Center

```
┌─────────────────────────────────────────────┐
│         In-App Notification Center           │
│                                             │
│  Storage: Redis Sorted Set per user          │
│  Key: inbox:{user_id}                        │
│  Score: timestamp                            │
│  Value: notification JSON                    │
│                                             │
│  APIs:                                       │
│  GET  /api/v1/inbox/{userId}?limit=20        │
│  POST /api/v1/inbox/{userId}/mark-read       │
│  GET  /api/v1/inbox/{userId}/unread-count    │
│                                             │
│  WebSocket: real-time push to open clients   │
│  ws://api.example.com/ws/notifications       │
│                                             │
│  Retention: 90 days, max 500 per user        │
│  (ZREMRANGEBYRANK to trim oldest)            │
└─────────────────────────────────────────────┘
```

### 11.5 Cross-Channel Escalation

**Candidate:** "If a user doesn't read a push notification within 30 minutes, we escalate to SMS, then email."

```python
# escalation_service.py

class EscalationService:
    ESCALATION_CHAINS = {
        'critical': [
            {'channel': 'push', 'wait': 0},
            {'channel': 'sms',  'wait': 300},    # 5 min after push
            {'channel': 'email', 'wait': 900},   # 15 min after push
        ],
        'high': [
            {'channel': 'push', 'wait': 0},
            {'channel': 'email', 'wait': 1800},  # 30 min after push
        ],
        'normal': [
            {'channel': 'push', 'wait': 0},      # No escalation
        ],
    }

    async def start_escalation(self, notification, priority):
        chain = self.ESCALATION_CHAINS.get(priority, [{'channel': 'push', 'wait': 0}])

        # Send first channel immediately
        first = chain[0]
        await self.send(notification, first['channel'])

        # Schedule subsequent channels
        for step in chain[1:]:
            await self.scheduler.schedule(
                task='check_and_escalate',
                payload={
                    'notification_id': notification.id,
                    'channel': step['channel'],
                    'user_id': notification.user_id,
                },
                delay_seconds=step['wait'],
            )

    async def check_and_escalate(self, notification_id, channel, user_id):
        status = await self.tracker.get_status(notification_id)

        # If already opened/clicked, don't escalate
        if status in ('opened', 'clicked'):
            return

        # Check if user has this channel enabled
        prefs = await self.preferences.get(user_id)
        if not prefs.get(channel, {}).get('enabled'):
            return

        await self.send_via_channel(notification_id, channel)
```

### 11.6 Analytics Dashboard Data Model

```sql
-- ClickHouse: notification analytics (columnar, fast aggregation)

CREATE TABLE notification_analytics (
    notification_id UUID,
    user_id         UUID,
    channel         LowCardinality(String),
    template_id     String,
    category        LowCardinality(String),
    priority        LowCardinality(String),
    status          LowCardinality(String),
    provider        LowCardinality(String),
    created_at      DateTime,
    sent_at         Nullable(DateTime),
    delivered_at    Nullable(DateTime),
    opened_at       Nullable(DateTime),
    clicked_at      Nullable(DateTime),
    country         LowCardinality(String),
    platform        LowCardinality(String)
) ENGINE = MergeTree()
ORDER BY (channel, created_at, template_id);

-- Example queries:
-- Daily delivery rate by channel
SELECT
    channel,
    toDate(created_at) AS day,
    count() AS total_sent,
    countIf(status = 'delivered') AS delivered,
    round(delivered / total_sent * 100, 2) AS delivery_rate
FROM notification_analytics
WHERE created_at >= today() - 7
GROUP BY channel, day
ORDER BY day, channel;

-- Template performance (open rate, click rate)
SELECT
    template_id,
    count() AS total_sent,
    countIf(opened_at IS NOT NULL) AS opened,
    countIf(clicked_at IS NOT NULL) AS clicked,
    round(opened / total_sent * 100, 2) AS open_rate,
    round(clicked / total_sent * 100, 2) AS click_rate
FROM notification_analytics
WHERE created_at >= today() - 30
GROUP BY template_id
ORDER BY total_sent DESC
LIMIT 20;
```

---

## 12. Interview Q&A

### Q1: How do you guarantee exactly-once delivery?

**Candidate:** "True exactly-once delivery is nearly impossible in a distributed system with external providers, but we can get close with at-least-once delivery plus deduplication:

1. **Idempotency keys** at the API layer — callers provide a unique key; repeated calls return the same result.
2. **Content-based deduplication** — hash(user_id + channel + content) checked against Redis with a 1-hour TTL window.
3. **Worker-level 'already sent' check** — before calling the provider, check a Redis flag `sent:{notification_id}`.
4. **Provider-level collapse** — APNs `apns-collapse-id` and FCM `collapse_key` ensure only the latest notification of a type shows on the device.
5. **Kafka consumer offset management** — commit offsets only after successful send and dedup flag set.

The combination gives us effective exactly-once semantics even though individual components offer at-least-once."

---

### Q2: How do you handle notification storms (e.g., system-wide alert to 100M users)?

**Candidate:** "A broadcast to 100M users is a fan-out problem. Here's my approach:

1. **Per-topic fan-out** — publish one message with audience='all_users'. Workers stream the user list in batches of 10,000.
2. **Rate limiting at the worker level** — cap at N sends/second per provider to avoid being throttled or blocked.
3. **Staggered delivery** — for non-critical broadcasts, spread delivery over 1-4 hours to avoid overwhelming downstream systems.
4. **Dedicated broadcast Kafka topic** — separate from real-time notifications so broadcasts don't block time-sensitive P0 messages.
5. **Provider batching** — FCM supports 'topic messaging' (send to a topic, FCM fans out). For email, SendGrid supports batch APIs (1,000 recipients per call).
6. **Pre-computed segments** — maintain materialized user lists for common segments (all users, all iOS users, all US users) so we don't need to query 100M rows at send time.

The goal is to decouple the publish speed from the delivery speed."

---

### Q3: How would you implement notification preferences with quiet hours across timezones?

**Candidate:** "This is tricky because quiet hours are relative to the user's local timezone.

1. **Store timezone per user** — we store the IANA timezone (e.g., 'America/New_York') alongside quiet hours.
2. **Convert to user-local time at send time** — when a notification arrives, we convert UTC now to the user's local time and check if it falls within their quiet hours.
3. **Overnight window handling** — quiet hours like 22:00–08:00 span midnight. The check is: `if start > end: return now >= start OR now < end`.
4. **Priority override** — P0 (critical) notifications like OTPs and security alerts bypass quiet hours entirely.
5. **Deferred delivery** — if a notification is suppressed by quiet hours, we schedule it for delivery at quiet_hours_end in the user's timezone.
6. **Batch timezone processing** — for broadcasts, group users by timezone offset to process the relevant cohort at the right time."

---

### Q4: Push notification token management — how do you handle stale/invalid tokens?

**Candidate:** "Token management is essential for push notification health:

1. **Token registration** — clients register tokens on every app launch (not just first install) via `POST /api/device-tokens`. Use `ON CONFLICT` upsert.
2. **Invalid token detection** — APNs returns HTTP 410 Gone, FCM returns `NOT_REGISTERED`. On these responses, immediately mark the token as inactive.
3. **Stale token cleanup** — a daily job deactivates tokens with no activity in 90 days. This keeps our token database lean.
4. **Token migration** — when a user gets a new device, the old token should become invalid. We handle this by allowing multiple tokens per user but removing old ones when providers report them invalid.
5. **Feedback service** — APNs provides a feedback service listing expired tokens. We poll it daily and clean up.
6. **Metrics** — track invalid_token_rate. If it spikes above 5%, investigate (app update breaking token registration, certificate expiry, etc.)."

---

### Q5: How do you prevent notification fatigue for users?

**Candidate:** "Notification fatigue leads to users disabling all notifications or uninstalling the app. Prevention strategies:

1. **Per-user rate limiting** — cap at N notifications per channel per hour per user. Example: max 5 push notifications/hour, 2 emails/day for marketing.
2. **Frequency capping by category** — separate limits for transactional (high) vs marketing (low) notifications.
3. **Smart aggregation** — batch multiple notifications of the same type into one digest ('5 new messages' instead of 5 separate notifications).
4. **User-level engagement scoring** — track open rates per user. If a user hasn't opened the last 10 notifications, reduce frequency or switch to digest mode.
5. **Preference granularity** — let users control notifications per category (order updates: ON, marketing: OFF, social: digest only).
6. **Time optimization** — ML model predicts the optimal send time per user based on historical engagement patterns."

---

### Q6: How would you implement cross-channel escalation?

**Candidate:** "Cross-channel escalation ensures critical notifications are seen:

1. **Define escalation chains per priority** — Critical: push → SMS (after 5 min) → email (after 15 min). High: push → email (after 30 min).
2. **Scheduler service** — after sending the first channel, schedule a 'check_and_escalate' job with the appropriate delay.
3. **Status check before escalating** — when the job fires, check if the notification has been opened/clicked. If yes, cancel escalation.
4. **Respect user preferences** — if a user has SMS disabled, skip that channel in the chain.
5. **Deduplication across channels** — the escalation message should reference the same notification_id, so users can see it's the same notification in their history.
6. **Cost awareness** — SMS is expensive. Only escalate to SMS for P0/P1 notifications. Include cost tracking per escalation path."

---

### Q7: How do you handle email deliverability and avoid spam filters?

**Candidate:** "Email deliverability is a complex domain:

1. **Authentication** — implement SPF, DKIM, and DMARC records for the sending domain. Without these, emails go straight to spam.
2. **Dedicated sending IPs** — separate IPs for transactional vs marketing email. A marketing spam complaint won't affect OTP delivery.
3. **IP warm-up** — new IPs start with 50 emails/day, gradually increasing over 30 days. Sudden volume spikes trigger spam filters.
4. **Bounce management** — immediately suppress hard-bounced emails. Retry soft bounces up to 3 times. Keep bounce rate < 2%.
5. **Complaint handling** — honor spam complaints instantly. Complaint rate > 0.1% triggers automatic sending pause.
6. **List hygiene** — periodically verify email addresses. Remove unengaged users (no opens in 6 months) from marketing lists.
7. **Content practices** — avoid spam trigger words, maintain good text-to-image ratio, always include unsubscribe link (required by CAN-SPAM/GDPR).
8. **Sender reputation monitoring** — check Google Postmaster Tools, Microsoft SNDS, and third-party tools like SenderScore."

---

### Q8: How would you design notification analytics (open rate, click-through rate)?

**Candidate:** "Analytics requires tracking user interactions:

1. **Open tracking (email)** — embed a 1×1 transparent tracking pixel: `<img src='https://api.example.com/track/open/{notification_id}'>`. When loaded, we record an open event.
2. **Open tracking (push)** — APNs and FCM don't provide open callbacks natively. The mobile app must report opens by calling `POST /api/notifications/{id}/opened` when the user taps or views the notification.
3. **Click tracking** — rewrite all links in emails through a redirect: `https://track.example.com/click/{notification_id}/{link_hash}` → record click → redirect to original URL.
4. **Storage** — write events to Kafka → consume into ClickHouse for fast aggregation.
5. **Metrics computed:**
   - Open rate = unique opens / delivered (email: 20-30% is good, push: 5-15%)
   - Click-through rate = unique clicks / delivered
   - Conversion rate = desired actions / delivered
6. **Privacy considerations** — Apple Mail Privacy Protection pre-fetches images, inflating open rates. Track click-through as the more reliable metric."

---

## 13. Production Checklist

### Pre-Launch

- [ ] Load test all channels at 2x peak (2,000 notifications/sec)
- [ ] Verify Kafka replication factor = 3, min ISR = 2
- [ ] Test provider failover: disable primary, verify secondary picks up
- [ ] Validate email authentication (SPF, DKIM, DMARC) via mail-tester.com
- [ ] Confirm dead letter queue consumer is running and alerting
- [ ] Set up Grafana dashboard with all key metrics
- [ ] Configure PagerDuty alerts for critical failures
- [ ] Test idempotency: send same request twice, verify no duplicate delivery
- [ ] Validate quiet hours logic across DST transitions
- [ ] Security review: API authentication, rate limiting, input sanitization

### Day 1

- [ ] Monitor delivery rates per channel (target > 95%)
- [ ] Verify consumer lag stays below 5,000 messages
- [ ] Check email bounce rate (target < 2%)
- [ ] Confirm push notification invalid token rate (target < 3%)
- [ ] Review DLQ — investigate any messages that landed there
- [ ] Validate structured logs are flowing to centralized logging

### Week 1

- [ ] Analyze delivery latency distributions (p50, p95, p99)
- [ ] Review provider cost against estimates
- [ ] Tune Kafka partition count based on actual throughput
- [ ] Adjust worker auto-scaling thresholds based on observed patterns
- [ ] Run chaos test: kill a worker pod, verify auto-recovery
- [ ] Begin email IP warm-up if using new sending IPs

### Month 1

- [ ] Audit notification volume vs budget (especially SMS costs)
- [ ] Review user opt-out rates — investigate if > 5%
- [ ] Implement notification analytics dashboard (open/click rates)
- [ ] Evaluate need for additional Kafka partitions
- [ ] Plan database archival: move notifications older than 90 days to cold storage
- [ ] Conduct capacity planning for next quarter (project 50% growth)
- [ ] Document runbooks for common operational scenarios

---

## Summary

### Technical Decisions

| Decision | Choice | Reasoning |
|----------|--------|-----------|
| Message Queue | Kafka | High throughput, durable, supports ordering by partition key |
| Primary DB | PostgreSQL | ACID for notifications/preferences, rich querying |
| Delivery Log DB | Cassandra | Write-optimized, handles 10M+ writes/day, TTL support |
| Analytics DB | ClickHouse | Columnar, fast aggregations over billions of events |
| Cache | Redis | Template cache, idempotency keys, dedup, rate limiting |
| Push Provider | APNs + FCM | Native iOS/Android support, free at scale |
| Email Provider | SendGrid (primary) + Mailgun (fallback) | Reliable, good API, webhook support |
| SMS Provider | Twilio | Global coverage, delivery receipts, programmable |
| Delivery Guarantee | At-least-once + deduplication | Practical alternative to true exactly-once |
| Priority System | 4-tier (P0–P3) | Separate SLAs, prevents low-priority floods from blocking critical |

### Scalability Path

```
Phase 1 (10M/day):
  - 3 Kafka brokers, 48 partitions
  - 4 push workers, 3 email workers, 2 SMS workers
  - Single PostgreSQL + Cassandra cluster
  - Estimated infra cost: ~$5,000/month

Phase 2 (100M/day):
  - 6 Kafka brokers, 96 partitions
  - 20 push workers, 12 email workers, 6 SMS workers
  - PostgreSQL sharded (4 shards), Cassandra 6-node cluster
  - Add ClickHouse for analytics
  - Estimated infra cost: ~$25,000/month

Phase 3 (1B/day):
  - 12+ Kafka brokers, 256+ partitions
  - 100+ workers across channels (auto-scaled)
  - PostgreSQL 16+ shards, Cassandra 20+ nodes
  - Multi-region deployment for latency + DR
  - Dedicated email sending infrastructure (own MTAs)
  - Estimated infra cost: ~$150,000/month
```

---

> **Interview Tip:** Start with the high-level architecture, then let the interviewer guide which component to deep-dive into. The delivery pipeline (Section 7) is where most interviewers spend time — be ready to discuss retry strategies, deduplication, and provider failover in detail.
