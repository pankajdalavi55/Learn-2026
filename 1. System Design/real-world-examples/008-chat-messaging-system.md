# Complete System Design: Chat / Messaging System (Production-Ready)

> **Complexity Level:** Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** WhatsApp, Facebook Messenger, Slack, Telegram, Signal

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Real-Time Messaging Engine](#7-deep-dive-real-time-messaging-engine)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a real-time chat/messaging system like WhatsApp that supports 1-on-1 and group messaging with delivery guarantees."

**Interviewer's Perspective:**  
This problem assesses:
- Real-time communication (WebSocket management at scale)
- Message ordering and delivery guarantees
- Presence/status tracking
- Fan-out for group messages
- Offline message delivery
- End-to-end encryption concepts

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I dive in, let me clarify both functional and non-functional requirements."

### 2.1 Functional Requirements

**Candidate:** "For functional requirements:
1. Should we support 1-on-1 messaging and group chats?
2. What is the maximum group size?
3. What message types do we need — text only, or also images, video, files?
4. Do we need read receipts and delivery indicators?
5. Do we need online/offline presence indicators?
6. Should we support message history and search?
7. Do we need push notifications for offline users?"

**Interviewer:** "Yes, let's support:
- 1-on-1 and group chats (up to 500 members)
- Text, image, video, file messages
- Read receipts and delivery status
- Online/offline presence
- Message history
- Push notifications for offline users"

**Candidate:** "Got it. Core features:
1. ✅ 1-on-1 messaging (real-time)
2. ✅ Group messaging (up to 500 members)
3. ✅ Multi-media messages (text, image, video, file)
4. ✅ Read receipts (sent → delivered → read)
5. ✅ Presence (online/offline/last seen)
6. ✅ Message history and search
7. ✅ Push notifications for offline users"

### 2.2 Non-Functional Requirements

**Candidate:** "For non-functional requirements:
1. What's the expected scale — daily active users?
2. What's the acceptable message delivery latency?
3. What guarantees do we need — at-least-once, exactly-once?
4. Do we need end-to-end encryption?
5. What's the availability target?"

**Interviewer:**
- Scale: 500M daily active users
- Latency: <100ms for online recipients
- Delivery: at-least-once (no message loss)
- E2E encryption: optional, discuss approach
- Availability: 99.99%

**Candidate:** "Summary:
- **Scale:** 500M DAU, ~50M concurrent connections
- **Latency:** <100ms message delivery for online users
- **Delivery:** At-least-once, no message loss
- **Availability:** 99.99% (52 minutes downtime/year)
- **Ordering:** Messages ordered within a conversation
- **Storage:** Persistent message history"

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me estimate the traffic:

**Messages:**
- 500M DAU × 40 messages/user/day = 20 billion messages/day
- Messages/sec: 20B / 86,400 ≈ **230,000 messages/sec**
- Peak (3x): **~700K messages/sec**

**Connections:**
- 10% of DAU online simultaneously = 50M concurrent WebSocket connections
- Each connection requires a persistent TCP socket

**Read receipts and presence:**
- At least 2x message rate for delivery/read confirmations
- Presence heartbeats: 50M users × every 30 sec = **1.7M heartbeats/sec**"

### 3.2 Storage Estimation

**Candidate:** "For storage:

**Per Message:**
- message_id: 16 bytes (UUID)
- conversation_id: 16 bytes
- sender_id: 8 bytes
- content: average 100 bytes (text)
- metadata: 60 bytes (timestamps, status, type)
- Total: ~200 bytes per text message

**Daily Storage:**
- Text: 20B × 200 bytes = **4 TB/day**
- Media: stored in object storage (S3), only URLs in message DB
- Media volume: assume 5% of messages have media, avg 500KB = 1B × 500KB = **500 TB/day**

**Yearly Storage:**
- Text messages: 4 TB/day × 365 = **1.5 PB/year**
- Media: separate object storage with lifecycle policies"

### 3.3 Bandwidth Estimation

**Candidate:** "For bandwidth:

**Text Messages:**
- 230K messages/sec × 200 bytes = **46 MB/sec** (negligible)

**Media Messages:**
- 230K × 5% × 500KB = **5.75 GB/sec** (significant)
- CDN absorbs most media delivery bandwidth

**WebSocket Overhead:**
- 50M connections × 200 bytes heartbeat × every 30 sec = **333 MB/sec**"

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌──────────────┐
│   Clients    │
│(iOS/Android/ │
│    Web)      │
└──────┬───────┘
       │ WebSocket + HTTPS
       ▼
┌──────────────────────────────────────────┐
│         Load Balancer (L4/L7)            │
└──────────────────┬───────────────────────┘
                   │
       ┌───────────┼───────────┐
       │           │           │
       ▼           ▼           ▼
┌────────────┐┌────────────┐┌────────────┐
│ WebSocket  ││ WebSocket  ││ WebSocket  │
│ Gateway 1  ││ Gateway 2  ││ Gateway N  │
│            ││            ││            │
│ (Manages   ││ (Manages   ││ (Manages   │
│  connections││ connections││ connections│
│  per user)  ││ per user)  ││ per user)  │
└─────┬──────┘└─────┬──────┘└─────┬──────┘
      │             │             │
      └─────────────┼─────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
        ▼           ▼           ▼
┌─────────────┐┌──────────┐┌──────────────┐
│   Chat      ││ Presence ││   Group      │
│  Service    ││ Service  ││  Service     │
│             ││ (Redis)  ││              │
│ (Message    ││          ││ (Membership, │
│  routing)   ││ Online/  ││  fan-out)    │
│             ││ Offline  ││              │
└──────┬──────┘└──────────┘└──────────────┘
       │
       ▼
┌──────────────────┐    ┌───────────────┐
│   Kafka          │    │   Media       │
│  (Message Queue) │    │   Service     │
│                  │    │   (S3 + CDN)  │
└────────┬─────────┘    └───────────────┘
         │
    ┌────┼────┐
    │         │
    ▼         ▼
┌─────────┐ ┌──────────────┐
│Cassandra│ │    Push       │
│(Message │ │ Notification  │
│ Store)  │ │   Service     │
└─────────┘ └──────────────┘
```

### 4.2 API Design

**Candidate:** "The system uses both WebSocket events and REST APIs:

**WebSocket Events (Real-Time):**

```javascript
// Client → Server
{
  "type": "send_message",
  "data": {
    "conversationId": "conv_123",
    "content": "Hello!",
    "messageType": "text",
    "clientMessageId": "client_msg_456"  // for deduplication
  }
}

// Server → Client
{
  "type": "new_message",
  "data": {
    "messageId": "msg_789",
    "conversationId": "conv_123",
    "senderId": "user_001",
    "content": "Hello!",
    "messageType": "text",
    "timestamp": "2026-01-15T10:30:00Z"
  }
}

// Typing indicator
{ "type": "typing", "data": { "conversationId": "conv_123", "userId": "user_001" } }

// Read receipt
{ "type": "read_receipt", "data": { "conversationId": "conv_123", "messageId": "msg_789" } }
```

**REST APIs (Non-Real-Time):**

```http
GET /api/conversations
GET /api/conversations/{id}/messages?before={messageId}&limit=50
POST /api/groups
PUT /api/groups/{id}/members
POST /api/media/upload
GET /api/users/{id}/status
```
"

### 4.3 Data Flow

**Candidate:** "Let me walk through the main flows:

**Flow 1: Send Message (Recipient Online)**
1. Sender's client sends message via WebSocket to Gateway A
2. Gateway A forwards to Chat Service
3. Chat Service persists message to Cassandra
4. Chat Service looks up recipient's gateway in connection registry (Redis)
5. Routes message to Gateway B (where recipient is connected)
6. Gateway B pushes message to recipient's WebSocket
7. Recipient's client sends delivery ACK
8. Chat Service updates message status to 'delivered'

**Flow 2: Send Message (Recipient Offline)**
1. Steps 1-3 same as above
2. Chat Service checks connection registry — recipient not found
3. Message stays persisted in Cassandra
4. Chat Service triggers Push Notification Service (APNs/FCM)
5. When recipient comes online, client fetches undelivered messages via sync API

**Flow 3: Group Message**
1. Sender sends message to group conversation
2. Chat Service persists message once
3. Group Service retrieves member list (e.g., 200 members)
4. For each online member: route to their gateway
5. For each offline member: queue push notification
6. Each member reads from the same message in Cassandra (no duplication)"

---

## 5. Deep Dive: Core Components

### 5.1 WebSocket Gateway

**Candidate:** "The WebSocket Gateway is the most critical component:

**Responsibilities:**
- Manage persistent WebSocket connections
- Authenticate connections (JWT token on handshake)
- Route messages between clients and backend services
- Handle heartbeats (detect stale connections)
- Register/deregister connections in the connection registry

**Technology:** Node.js or Go (both excellent for concurrent connections)

**Connection Registry (Redis):**
```
Key: user:{userId}:connection
Value: { gatewayId: "gw-002", connectedAt: timestamp, deviceId: "..." }
TTL: 5 minutes (refreshed on heartbeat)
```

**Connection Capacity:**
- Each gateway server: ~500K concurrent WebSocket connections
- 50M concurrent users / 500K per server = **100 gateway servers**
- Each connection uses ~10KB memory = 500K × 10KB = **5GB RAM per server**"

### 5.2 Chat Service

**Candidate:** "The Chat Service handles message routing logic:

```javascript
async function handleSendMessage(senderId, conversationId, content, clientMsgId) {
    // Deduplication check
    const existing = await redis.get(`dedup:${clientMsgId}`);
    if (existing) return JSON.parse(existing);

    const messageId = generateSnowflakeId();
    const message = {
        messageId,
        conversationId,
        senderId,
        content,
        type: 'text',
        status: 'sent',
        createdAt: Date.now()
    };

    // Persist to Cassandra
    await cassandra.execute(
        'INSERT INTO messages (conversation_id, message_id, sender_id, content, type, created_at) VALUES (?, ?, ?, ?, ?, ?)',
        [conversationId, messageId, senderId, content, 'text', message.createdAt]
    );

    // Set dedup key (TTL 24 hours)
    await redis.setex(`dedup:${clientMsgId}`, 86400, JSON.stringify(message));

    // Get recipients
    const recipients = await getConversationMembers(conversationId);

    for (const recipientId of recipients) {
        if (recipientId === senderId) continue;
        await routeMessageToRecipient(recipientId, message);
    }

    return message;
}

async function routeMessageToRecipient(recipientId, message) {
    const connectionInfo = await redis.get(`user:${recipientId}:connection`);

    if (connectionInfo) {
        const { gatewayId } = JSON.parse(connectionInfo);
        await publishToGateway(gatewayId, recipientId, message);
    } else {
        await pushNotificationService.send(recipientId, message);
    }
}
```
"

### 5.3 Presence Service

**Candidate:** "Presence tracks who is online:

**Approach: Heartbeat-based with Redis**

```javascript
async function handleHeartbeat(userId) {
    const key = `presence:${userId}`;
    const data = JSON.stringify({
        status: 'online',
        lastSeen: Date.now(),
        deviceId: deviceId
    });
    await redis.setex(key, 60, data);  // TTL 60 seconds
}

async function getUserStatus(userId) {
    const data = await redis.get(`presence:${userId}`);
    if (data) {
        return { status: 'online', ...JSON.parse(data) };
    }
    // Fallback: check last_seen in DB
    const user = await db.query('SELECT last_seen FROM users WHERE id = ?', [userId]);
    return { status: 'offline', lastSeen: user.last_seen };
}
```

**Presence fan-out:** When a user comes online, notify their contacts. For a user with 500 contacts, that is 500 presence updates. At 50M online users, this is a lot of traffic — so we only fan out to users who are currently viewing the contact list."

### 5.4 Group Service

**Candidate:** "Manages group metadata and membership:

- Stores group info: name, avatar, admin, member list
- Handles join/leave/add/remove operations
- For message fan-out: provides member list for a given group
- Caches member lists in Redis (invalidate on membership change)

**Optimization for large groups:**
- Groups > 100 members: store member list in Cassandra, cache in Redis
- Fan-out via Kafka: publish message once to a group topic, each member's gateway consumes
"

---

## 6. Deep Dive: Database Design

### 6.1 Message Store (Cassandra)

**Candidate:** "Cassandra is ideal for messages because of write-heavy workload and time-series access patterns.

```sql
CREATE TABLE messages (
    conversation_id UUID,
    message_id BIGINT,        -- Snowflake ID (time-ordered)
    sender_id UUID,
    content TEXT,
    message_type TEXT,         -- 'text', 'image', 'video', 'file'
    media_url TEXT,
    status TEXT,               -- 'sent', 'delivered', 'read'
    created_at TIMESTAMP,
    PRIMARY KEY (conversation_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);
```

**Access Patterns:**
- Get latest messages: `WHERE conversation_id = ? ORDER BY message_id DESC LIMIT 50`
- Get messages before cursor: `WHERE conversation_id = ? AND message_id < ? LIMIT 50`
- Partition by conversation_id: all messages for a chat are co-located

**Why Cassandra:**
1. ✅ High write throughput (230K writes/sec)
2. ✅ Time-ordered data within partition (clustering key)
3. ✅ Linear scalability (add nodes for more capacity)
4. ✅ Tunable consistency (use LOCAL_QUORUM for messages)
5. ✅ Automatic data distribution and replication"

### 6.2 Conversations & Groups (PostgreSQL)

```sql
CREATE TABLE conversations (
    id UUID PRIMARY KEY,
    type VARCHAR(10) NOT NULL,     -- 'direct' or 'group'
    name VARCHAR(255),             -- NULL for direct chats
    avatar_url TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    last_message_at TIMESTAMP,
    last_message_preview TEXT
);

CREATE TABLE conversation_members (
    conversation_id UUID REFERENCES conversations(id),
    user_id UUID REFERENCES users(id),
    role VARCHAR(20) DEFAULT 'member',  -- 'admin', 'member'
    joined_at TIMESTAMP DEFAULT NOW(),
    last_read_message_id BIGINT,
    muted_until TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_user_conversations ON conversation_members(user_id, conversation_id);
```

### 6.3 User & Device Management (PostgreSQL)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    display_name VARCHAR(100),
    avatar_url TEXT,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_devices (
    user_id UUID REFERENCES users(id),
    device_id VARCHAR(100),
    platform VARCHAR(20),          -- 'ios', 'android', 'web'
    push_token TEXT,
    public_key TEXT,               -- for E2E encryption
    last_active TIMESTAMP,
    PRIMARY KEY (user_id, device_id)
);
```

---

## 7. Deep Dive: Real-Time Messaging Engine

### 7.1 WebSocket Connection Lifecycle

**Candidate:** "The lifecycle of a WebSocket connection:

```
1. CONNECT
   Client → wss://chat.example.com/ws?token=JWT_TOKEN
   Server validates JWT, extracts userId
   Server registers in connection registry

2. SUBSCRIBE
   Server auto-subscribes client to all their conversations
   Loads undelivered messages since last_seen_message_id

3. HEARTBEAT (every 30 seconds)
   Client → { type: "ping" }
   Server → { type: "pong" }
   Server refreshes presence TTL in Redis

4. MESSAGE EXCHANGE
   Bidirectional: send/receive messages, typing, receipts

5. DISCONNECT
   Client closes connection OR heartbeat timeout
   Server deregisters from connection registry
   Server updates last_seen in user DB
```
"

### 7.2 Connection Management at Scale

**Candidate:** "Managing 50M concurrent connections is a significant challenge:

**Connection Registry:**
```
Redis Cluster:
  user:user_001:connection → { gateway: "gw-007", device: "iphone_14", ts: 1704067200 }
  user:user_002:connection → { gateway: "gw-023", device: "pixel_8", ts: 1704067180 }
  ...
```

**Multi-Device Support:**
A user may be connected from multiple devices simultaneously:
```
user:user_001:connections → {
  "iphone_14": { gateway: "gw-007", ts: 1704067200 },
  "macbook":   { gateway: "gw-015", ts: 1704067100 }
}
```
Messages are delivered to ALL connected devices.

**Gateway-to-Gateway Communication:**
When a message needs to go from Gateway A to a user on Gateway B:

```
Option 1: Internal message bus (Redis Pub/Sub)
  - Each gateway subscribes to its own channel: gateway:gw-007
  - Chat Service publishes to the target gateway's channel

Option 2: Direct gRPC between gateways
  - Lower latency, but requires service discovery
  - Each gateway exposes a gRPC endpoint for message delivery

My recommendation: Redis Pub/Sub for simplicity, switch to gRPC at higher scale.
```
"

### 7.3 Message Ordering

**Candidate:** "Message ordering is critical for chat. Here's how we guarantee it:

**Per-Conversation Ordering:**
- Use Snowflake IDs as message IDs (time-based, monotonically increasing)
- Cassandra clustering key orders messages within a partition
- Client displays messages sorted by message_id

**Handling Out-of-Order Delivery:**
```javascript
// Client-side: maintain a local message buffer
class MessageBuffer {
    constructor() {
        this.messages = new Map(); // conversationId → sorted messages
    }

    addMessage(conversationId, message) {
        if (!this.messages.has(conversationId)) {
            this.messages.set(conversationId, []);
        }
        const msgs = this.messages.get(conversationId);

        // Insert in sorted position by messageId
        const insertIdx = this.binarySearch(msgs, message.messageId);
        msgs.splice(insertIdx, 0, message);

        // Detect gaps: check if there is a gap in sequence
        this.detectAndFillGaps(conversationId, msgs);
    }

    async detectAndFillGaps(conversationId, msgs) {
        // If we receive msg 105 but last was 102, fetch 103-104 from server
        // ... gap detection logic
    }
}
```

**Concurrent Senders:**
- Two users sending simultaneously to the same conversation
- Each gets a unique Snowflake ID (different worker bits)
- Cassandra INSERT is last-write-wins at the partition level, but since each message has a unique ID, there is no conflict
- Total ordering: all participants see the same message_id order"

### 7.4 Group Message Fan-Out

**Candidate:** "For group messages, we fan out to all members:

**Small Groups (≤100 members):**
```javascript
async function fanOutGroupMessage(groupId, message) {
    const members = await redis.smembers(`group:${groupId}:members`);

    const deliveryPromises = members.map(async (memberId) => {
        if (memberId === message.senderId) return;

        const conn = await redis.get(`user:${memberId}:connection`);
        if (conn) {
            const { gateway } = JSON.parse(conn);
            await redis.publish(`gateway:${gateway}`, JSON.stringify({
                type: 'new_message',
                recipientId: memberId,
                message
            }));
        } else {
            await pushNotificationQueue.add({ userId: memberId, message });
        }
    });

    await Promise.allSettled(deliveryPromises);
}
```

**Large Groups (100-500 members):**
- Use Kafka for fan-out instead of direct routing
- Publish message once to topic `group-messages-{groupId}`
- Each gateway has a consumer that filters relevant messages for its connected users
- Reduces redundant network hops

**Message Storage:**
- Store message ONCE in Cassandra (partitioned by conversation_id)
- All group members read from the same partition
- No duplication — very storage efficient"

### 7.5 Read Receipts and Delivery Status

**Candidate:** "Message status follows this state machine:

```
  SENT ──→ DELIVERED ──→ READ
   │
   └──→ FAILED (retry or give up)

Transitions:
- SENT: Server has persisted the message
- DELIVERED: Recipient's device received the message (ACK)
- READ: Recipient opened the conversation (client sends read_receipt)
```

**Implementation:**
```javascript
// When recipient's client receives message
client.on('new_message', async (message) => {
    displayMessage(message);
    // Send delivery ACK
    ws.send(JSON.stringify({
        type: 'delivery_ack',
        messageId: message.messageId,
        conversationId: message.conversationId
    }));
});

// When user opens conversation
function onConversationOpened(conversationId) {
    const latestMessageId = getLatestMessageId(conversationId);
    ws.send(JSON.stringify({
        type: 'read_receipt',
        conversationId,
        messageId: latestMessageId  // "read up to this message"
    }));
}
```

**Server handles receipts:**
```javascript
async function handleDeliveryAck(userId, messageId, conversationId) {
    await cassandra.execute(
        'UPDATE messages SET status = ? WHERE conversation_id = ? AND message_id = ?',
        ['delivered', conversationId, messageId]
    );
    // Notify sender
    const message = await getMessage(conversationId, messageId);
    await routeToUser(message.senderId, {
        type: 'status_update',
        messageId, status: 'delivered'
    });
}
```

**Group read receipts:**
- Track per-member read position: `conversation_members.last_read_message_id`
- Show 'read by N' instead of individual checkmarks (WhatsApp-style)
- Update asynchronously to avoid write amplification"

### 7.6 Typing Indicators

**Candidate:** "Typing indicators are ephemeral — never persisted:

```javascript
// Client sends typing event (throttled to once per 3 seconds)
function onUserTyping(conversationId) {
    if (Date.now() - lastTypingEvent < 3000) return;
    lastTypingEvent = Date.now();

    ws.send(JSON.stringify({
        type: 'typing',
        conversationId
    }));
}

// Server broadcasts to other participants (NO persistence)
async function handleTyping(senderId, conversationId) {
    const members = await getOnlineMembers(conversationId);
    for (const memberId of members) {
        if (memberId === senderId) continue;
        await routeToUser(memberId, {
            type: 'typing',
            conversationId,
            userId: senderId
        });
    }
}
```

**Client-side display:**
- Show 'User is typing...' for 5 seconds after last typing event
- Auto-clear if no new typing event received"

### 7.7 Offline Message Sync

**Candidate:** "When a user comes back online, they need all missed messages:

```javascript
// Client reconnection flow
async function onReconnect(userId, lastSyncedMessageId) {
    // For each conversation, fetch messages after lastSyncedMessageId
    const conversations = await getActiveConversations(userId);

    for (const conv of conversations) {
        const lastRead = conv.lastReadMessageId;
        const newMessages = await cassandra.execute(
            'SELECT * FROM messages WHERE conversation_id = ? AND message_id > ? LIMIT 100',
            [conv.id, lastRead]
        );

        // Send all missed messages to client in order
        for (const msg of newMessages.rows) {
            ws.send(JSON.stringify({ type: 'new_message', data: msg }));
        }
    }
}
```

**Optimization for heavy offline periods:**
- If user was offline for days, paginate sync (100 messages per batch)
- Priority: sync most recent conversations first
- Background sync older conversations
- Client shows a 'syncing...' indicator"

### 7.8 End-to-End Encryption (Overview)

**Candidate:** "For E2E encryption, I'd follow the Signal Protocol approach:

**Key Concepts:**
1. Each device has a long-term identity key pair and medium-term signed pre-keys
2. Pre-keys are uploaded to the server during registration
3. When Alice wants to message Bob, she fetches Bob's pre-key bundle from the server
4. Alice performs X3DH (Extended Triple Diffie-Hellman) key agreement to establish a shared secret
5. Messages are encrypted using Double Ratchet Algorithm (forward secrecy)

**Server's role in E2E:**
- Store encrypted messages (server cannot decrypt)
- Distribute public keys and pre-key bundles
- Relay encrypted messages between users

```javascript
// Simplified: server stores and relays ciphertext
async function handleEncryptedMessage(senderId, recipientId, encryptedPayload) {
    const message = {
        messageId: generateId(),
        senderId,
        recipientId,
        encryptedContent: encryptedPayload,  // opaque to server
        timestamp: Date.now()
    };

    await persistMessage(message);
    await routeToRecipient(recipientId, message);
}
```

**Trade-offs:**
- Pro: Maximum privacy, server can't read messages
- Con: Server-side search not possible, device-bound keys, backup complexity"

### 7.9 Media Message Handling

**Candidate:** "For images/videos/files:

```
1. Client uploads media to Media Service
   POST /api/media/upload → returns mediaUrl

2. Client sends message with mediaUrl
   { type: "image", mediaUrl: "https://cdn.example.com/abc123.jpg", thumbnail: "base64..." }

3. Recipient receives message with mediaUrl
   Client downloads media from CDN

Optimization:
- Generate thumbnail on upload (server-side)
- Send thumbnail inline with message (small base64)
- Lazy-load full media on tap
- Media stored in S3, served via CloudFront CDN
- Encryption: media encrypted before upload (E2E), decryption key in message
```
"

---

## 8. Scaling Strategies

### 8.1 Current Bottlenecks

**Candidate:** "At 230K messages/sec:

1. **WebSocket Gateways:** 50M connections across 100 servers — manageable
2. **Cassandra:** 230K writes/sec — well within Cassandra cluster capacity
3. **Redis (Connection Registry):** 230K lookups/sec — trivial for Redis
4. **Fan-out for groups:** 500-member groups × 230K/sec — Kafka handles this

Primary bottleneck: Gateway-to-gateway message routing at peak."

### 8.2 Scaling to 10x (2.3M messages/sec)

**Candidate:** "

**Step 1: Scale WebSocket Gateways**
- 500M concurrent connections / 500K per server = 1,000 servers
- Use consistent hashing to distribute users across gateways

**Step 2: Scale Cassandra**
- Add nodes to the ring (Cassandra scales linearly)
- Ensure partition key (conversation_id) distributes evenly

**Step 3: Kafka Partitioning**
- Partition by conversation_id for ordering
- Increase partition count as throughput grows

**Step 4: Connection Registry Sharding**
- Redis Cluster with hash slot distribution
- User ID hash determines which Redis shard"

### 8.3 Scaling to 100x (23M messages/sec)

**Candidate:** "At 100x, we need:

**Multi-Region Deployment:**
```
US Region:
  - WebSocket Gateways (users in US connect here)
  - Regional Cassandra cluster
  - Regional Redis cluster

EU Region:
  - Same setup for European users

Asia Region:
  - Same setup

Cross-Region Messages:
  - User in US messages user in EU
  - US Chat Service → Kafka bridge → EU Chat Service → EU Gateway
  - Cassandra cross-datacenter replication for message history
```

**Protocol Optimization:**
- Switch from JSON to Protocol Buffers (50% bandwidth reduction)
- Message batching: aggregate multiple messages into one frame
- Connection multiplexing for multi-device users"

---

## 9. Failure Scenarios & Mitigation

### 9.1 WebSocket Gateway Crash

**Scenario:** A gateway server dies, dropping 500K connections.

**Impact:**
- 500K users temporarily disconnected
- Messages to those users can't be delivered via WebSocket

**Mitigation:**
```javascript
// Client-side: automatic reconnection with exponential backoff
class WebSocketClient {
    connect() {
        this.ws = new WebSocket(this.url);
        this.ws.onclose = () => this.reconnect();
        this.ws.onerror = () => this.reconnect();
    }

    reconnect() {
        const delay = Math.min(1000 * Math.pow(2, this.retryCount), 30000);
        setTimeout(() => {
            this.retryCount++;
            this.connect();
            // On reconnect, sync missed messages
            this.syncMissedMessages();
        }, delay);
    }
}
```
- Load balancer detects failure in 10 seconds, routes new connections to healthy servers
- Connection registry entries expire via TTL (60 seconds)
- Messages sent during disconnect are persisted in Cassandra and delivered on reconnect
- **RTO: 10-30 seconds** (time for client reconnection)

### 9.2 Cassandra Node Failure

**Scenario:** A Cassandra node goes down.

**Impact:**
- Partition ranges owned by that node temporarily unavailable
- If replication factor = 3, other replicas serve reads/writes

**Mitigation:**
- Replication factor 3 across racks
- Consistency level LOCAL_QUORUM (2 out of 3 must ACK)
- Hinted handoff: surviving nodes store writes for the failed node
- Node replacement: Cassandra rebalances automatically
- **Impact: Zero downtime** if RF=3 and only 1 node fails

### 9.3 Message Ordering Violation

**Scenario:** Due to network delays, messages arrive out of order.

**Impact:**
- Conversation appears jumbled

**Mitigation:**
- Snowflake IDs guarantee monotonic ordering from server's perspective
- Client sorts by message_id before display
- If client detects gap (message 105 arrives but 103, 104 missing), request sync from server
- Cassandra clustering key ensures storage-level ordering

### 9.4 Split Brain in Presence Service

**Scenario:** Redis cluster partition causes inconsistent presence data.

**Impact:**
- User appears online on one partition, offline on another
- Messages may be routed incorrectly

**Mitigation:**
- Use Redis Cluster with majority quorum
- Presence data has short TTL (60 seconds) — self-healing
- If message delivery fails (user not actually on that gateway), fall back to push notification
- Accept eventual consistency for presence (it's not critical)

### 9.5 Duplicate Message Delivery

**Scenario:** Network hiccup causes client to retry, or server processes message twice.

**Impact:**
- User sees the same message twice

**Mitigation:**
```javascript
// Client-generated deduplication ID
const clientMsgId = `${deviceId}-${localSeqNumber}`;

// Server deduplication
async function handleMessage(clientMsgId, message) {
    const dedupKey = `dedup:${clientMsgId}`;
    const exists = await redis.get(dedupKey);
    if (exists) return JSON.parse(exists); // return cached response

    const result = await processMessage(message);
    await redis.setex(dedupKey, 86400, JSON.stringify(result));
    return result;
}
```

### 9.6 Push Notification Failure

**Scenario:** APNs or FCM service is down; offline users don't get notified.

**Mitigation:**
- Retry with exponential backoff (up to 3 attempts)
- Dead letter queue for permanently failed notifications
- Monitor push delivery rate
- When user comes online, sync catches up all missed messages regardless of push status
- Push notification is a best-effort enhancement, not the primary delivery mechanism

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

**Candidate:** "For a chat system, I'd track:

**Application Metrics (RED):**
1. **Rate:** Messages/sec (sent, delivered, read)
2. **Errors:** Failed deliveries, WebSocket errors
3. **Duration:** Message delivery latency (send to receive)

**Business Metrics:**
- Active conversations per hour
- Messages per user per day
- Group creation rate
- Media message percentage

**Infrastructure Metrics (USE):**
1. **Utilization:** WebSocket connections per gateway, Cassandra disk usage
2. **Saturation:** Kafka consumer lag, connection queue depth
3. **Errors:** Cassandra timeouts, Redis connection failures

**Example Dashboard (Grafana):**
```
Row 1: Message Delivery
- [Graph] Messages sent/delivered/read per second
- [Gauge] Delivery success rate (%)
- [Heatmap] End-to-end delivery latency

Row 2: Connections
- [Graph] Active WebSocket connections per gateway
- [Graph] Connection/disconnection rate
- [Gauge] Total concurrent connections

Row 3: Infrastructure
- [Graph] Cassandra write/read latency
- [Graph] Kafka consumer lag per partition
- [Graph] Redis memory usage and ops/sec

Row 4: Push Notifications
- [Graph] Push delivery rate per platform (iOS/Android)
- [Graph] Push failure rate
- [Graph] Time to deliver push notification
```
"

### 10.2 Alerting Rules

```yaml
alert: HighMessageDeliveryLatency
expr: histogram_quantile(0.99, message_delivery_latency_seconds) > 1
for: 5m
severity: critical
message: "p99 message delivery latency above 1 second"

alert: WebSocketConnectionDrop
expr: rate(websocket_disconnections_total[5m]) > 1000
for: 2m
severity: critical
message: "Unusual WebSocket disconnection rate - possible gateway failure"

alert: KafkaConsumerLag
expr: kafka_consumer_lag > 100000
for: 5m
severity: warning
message: "Kafka consumer lag exceeding 100K messages"

alert: CassandraWriteLatency
expr: cassandra_write_latency_p99 > 0.1
for: 5m
severity: warning
message: "Cassandra p99 write latency above 100ms"
```

### 10.3 Distributed Tracing

```
Trace: msg_delivery_abc123

Span 1: WebSocket Gateway receive (2ms)
Span 2: Chat Service process (5ms)
  Span 2a: Deduplication check - Redis (1ms)
  Span 2b: Cassandra write (8ms)
  Span 2c: Connection registry lookup (1ms)
Span 3: Inter-gateway routing (3ms)
Span 4: WebSocket Gateway deliver (1ms)
Total: 21ms
```

---

## 11. Advanced Features

### 11.1 Message Reactions

```javascript
// WebSocket event
{
    "type": "reaction",
    "data": {
        "messageId": "msg_789",
        "conversationId": "conv_123",
        "emoji": "👍",
        "action": "add"  // or "remove"
    }
}

// Store reactions in Cassandra
// CREATE TABLE message_reactions (
//     conversation_id UUID,
//     message_id BIGINT,
//     user_id UUID,
//     emoji TEXT,
//     created_at TIMESTAMP,
//     PRIMARY KEY ((conversation_id, message_id), user_id)
// );
```

### 11.2 Message Editing and Deletion

```javascript
async function editMessage(userId, conversationId, messageId, newContent) {
    // Verify sender
    const original = await getMessage(conversationId, messageId);
    if (original.senderId !== userId) throw new Error('Unauthorized');

    await cassandra.execute(
        'UPDATE messages SET content = ?, edited_at = ? WHERE conversation_id = ? AND message_id = ?',
        [newContent, Date.now(), conversationId, messageId]
    );

    // Notify all participants
    await broadcastToConversation(conversationId, {
        type: 'message_edited',
        messageId, newContent, editedAt: Date.now()
    });
}

async function deleteMessage(userId, conversationId, messageId) {
    // Soft delete: replace content
    await cassandra.execute(
        'UPDATE messages SET content = ?, deleted_at = ? WHERE conversation_id = ? AND message_id = ?',
        ['This message was deleted', Date.now(), conversationId, messageId]
    );

    await broadcastToConversation(conversationId, {
        type: 'message_deleted',
        messageId
    });
}
```

### 11.3 Voice/Video Calls (WebRTC Signaling)

**Candidate:** "For calls, the chat infrastructure handles signaling only:

```
1. Caller sends 'call_offer' via WebSocket (contains SDP offer)
2. Server routes to callee's gateway
3. Callee receives offer, sends 'call_answer' (SDP answer)
4. Both exchange ICE candidates via WebSocket
5. Direct peer-to-peer connection established via WebRTC
6. Media flows directly between devices (not through server)

For group calls: use an SFU (Selective Forwarding Unit) server
```
"

### 11.4 Disappearing Messages

```javascript
// When creating a conversation, set message TTL
await db.query(
    'UPDATE conversations SET message_ttl_seconds = ? WHERE id = ?',
    [86400, conversationId]  // 24 hours
);

// Background job: delete expired messages
async function cleanupExpiredMessages() {
    const conversations = await db.query(
        'SELECT id, message_ttl_seconds FROM conversations WHERE message_ttl_seconds IS NOT NULL'
    );

    for (const conv of conversations) {
        const cutoff = Date.now() - (conv.message_ttl_seconds * 1000);
        await cassandra.execute(
            'DELETE FROM messages WHERE conversation_id = ? AND created_at < ?',
            [conv.id, cutoff]
        );
    }
}
// Run every hour
```

### 11.5 Message Search

**Candidate:** "Full-text search across billions of messages:

- Index messages in Elasticsearch (async via Kafka)
- Partition Elasticsearch index by user_id (privacy: each user only searches their own messages)
- For E2E encrypted chats: search only on client-side (server can't index encrypted content)

```http
GET /api/messages/search?q=meeting+tomorrow&conversationId=conv_123

Response:
{
    "results": [
        {
            "messageId": "msg_456",
            "conversationId": "conv_123",
            "content": "Let's have a meeting tomorrow at 3pm",
            "sender": "Alice",
            "timestamp": "2026-01-15T09:00:00Z",
            "highlights": ["<em>meeting</em> <em>tomorrow</em>"]
        }
    ]
}
```
"

---

## 12. Interview Q&A

### Q1: How do you ensure message ordering in a distributed system?

**Answer:**
We use Snowflake IDs (time-based, 64-bit) as message IDs. Since the timestamp is embedded in the ID, messages are naturally time-ordered. Within a Cassandra partition (conversation_id), the clustering key (message_id) maintains strict order. The client also sorts by message_id before display and detects gaps for re-sync.

For concurrent senders, each Chat Service instance generates unique Snowflake IDs (different machine bits), so there are no collisions. The total order is determined by the ID value, which reflects approximate send time. This is sufficient for chat — users don't need nanosecond precision in ordering.

### Q2: How do you handle 50M concurrent WebSocket connections?

**Answer:**
We horizontally scale WebSocket gateways. Each server handles ~500K connections (Go/Node.js with epoll/kqueue). At 50M connections, we need ~100 gateway servers. Each connection uses ~10KB memory, so each server needs ~5GB RAM for connections.

The key challenge is routing: when User A (on Gateway 7) messages User B (on Gateway 23), we use a Redis-based connection registry to look up which gateway serves each user. Inter-gateway communication uses Redis Pub/Sub or gRPC.

For sticky connections, we use the load balancer's connection-level (L4) routing — once a WebSocket is established, all frames stay on the same gateway.

### Q3: How would you implement end-to-end encryption?

**Answer:**
We'd implement the Signal Protocol:

1. **Key generation:** Each device generates identity key pair + signed pre-keys + one-time pre-keys
2. **Key distribution:** Server stores public pre-key bundles (cannot decrypt messages)
3. **Session establishment:** X3DH (Extended Triple Diffie-Hellman) key agreement
4. **Message encryption:** Double Ratchet Algorithm provides forward secrecy and break-in recovery
5. **Multi-device:** Each device has its own keys; messages are encrypted separately for each device

The server only stores and relays ciphertext. Trade-off: server-side search becomes impossible, and lost devices mean lost message history (unless encrypted backups are implemented).

### Q4: How do you handle message delivery to offline users?

**Answer:**
Messages for offline users are persisted in Cassandra regardless of online status. Additionally:

1. Push notification sent via APNs (iOS) or FCM (Android) — shows preview
2. When user comes online, client sends last_synced_message_id
3. Server returns all messages after that ID for each conversation
4. Client processes and displays them in order

For users offline for extended periods, we paginate the sync (100 messages per batch, newest conversations first) and sync older data in the background.

### Q5: How do you scale group chat with 500 members?

**Answer:**
For group messages, we store the message ONCE in Cassandra (partition by conversation_id) and fan out delivery notifications:

- **Small groups (≤50):** Direct fan-out — look up each member's gateway and route the message
- **Medium groups (50-200):** Kafka-based fan-out — publish to a group topic, each gateway's consumer delivers to local members
- **Large groups (200-500):** Same as medium, but with batched delivery and rate limiting to prevent thundering herd

Read receipts in large groups are aggregated ("read by 45 members") instead of showing individual indicators, reducing write amplification.

### Q6: How would you implement read receipts at scale?

**Answer:**
Read receipts use a "high-water mark" approach: instead of tracking each message individually, we track the latest message_id each user has read per conversation.

```sql
UPDATE conversation_members
SET last_read_message_id = 12345
WHERE conversation_id = 'conv_123' AND user_id = 'user_456';
```

This means one update per conversation-open event, not per message. For display:
- 1-on-1: compare partner's last_read_message_id with each message_id (blue checkmarks)
- Groups: count members whose last_read_message_id >= message_id ("read by N")

Read receipt updates are sent to the conversation sender in real-time via WebSocket, but batched (every 2 seconds) to avoid excessive traffic.

### Q7: How do you handle message search across billions of messages?

**Answer:**
We use Elasticsearch for full-text search, indexed asynchronously:

1. Messages flow through Kafka → Elasticsearch indexer
2. Index partitioned by user_id (privacy boundary)
3. Each user can only search conversations they belong to
4. Elasticsearch handles fuzzy matching, phrase queries, and relevance ranking

For E2E encrypted conversations, server-side search is impossible. Options:
- Client-side search (load messages to device, search locally)
- Encrypted search indexes (homomorphic encryption — experimental, high overhead)
- Store searchable metadata separately (sender, timestamp) while content remains encrypted

### Q8: How would you implement disappearing messages?

**Answer:**
Disappearing messages have a configurable TTL per conversation:

1. When enabled, new messages include an `expires_at` timestamp
2. Client-side timer deletes messages from local storage after TTL
3. Server-side background job periodically scans and deletes expired messages from Cassandra
4. Using Cassandra TTL feature: `INSERT INTO messages ... USING TTL 86400` (auto-deleted after 24 hours)

Challenges:
- Client clock manipulation: server-enforced TTL is the source of truth
- Screenshots: cannot prevent (display warning to users)
- Cached messages in notifications: clear notification content after TTL
- Backup sync: exclude disappearing messages from backup

---

## 13. Production Checklist

### 13.1 Pre-Launch

- [ ] **Load Testing:** Simulate 1M concurrent WebSocket connections using k6 or Artillery
- [ ] **Message Delivery Test:** Verify delivery under all conditions (online, offline, multi-device)
- [ ] **Failover Testing:** Kill gateway servers, verify seamless reconnection
- [ ] **Security Audit:**
  - [ ] WebSocket authentication (JWT validation)
  - [ ] Rate limiting on message sending
  - [ ] Input validation (message size limits, content type validation)
  - [ ] TLS for all connections
- [ ] **Data Migration:** If migrating from existing system, plan message history migration
- [ ] **Push Notification Setup:** Configure APNs/FCM certificates and tokens
- [ ] **Monitoring:** Dashboards, alerts, on-call rotation

### 13.2 Day-1 Operations

- [ ] Monitor WebSocket connection count and stability
- [ ] Monitor message delivery latency (p50, p95, p99)
- [ ] Check Cassandra write/read latency
- [ ] Verify push notification delivery rate
- [ ] Monitor Kafka consumer lag

### 13.3 Week-1 Optimization

- [ ] Analyze message delivery patterns (peak hours, geographic distribution)
- [ ] Tune WebSocket heartbeat intervals
- [ ] Optimize Cassandra compaction strategy
- [ ] Review and tune Kafka partition counts
- [ ] Analyze and resolve any message ordering issues

### 13.4 Month-1 Scaling

- [ ] Review gateway capacity (plan for 3x growth)
- [ ] Optimize media pipeline (thumbnail generation, CDN hit rate)
- [ ] Implement message search if not in MVP
- [ ] Plan multi-region deployment based on user distribution
- [ ] Cost optimization (reserved instances, storage tiering)

---

## Summary: Key Takeaways

### Technical Decisions

| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Real-Time Protocol** | WebSocket | Bidirectional, low-latency, persistent connection |
| **Message Store** | Cassandra | Write-heavy, time-series, linear scalability |
| **Connection Registry** | Redis | Sub-ms lookups, TTL for presence, Pub/Sub for routing |
| **Message Queue** | Kafka | Fan-out, ordering by partition, replay capability |
| **User/Group Metadata** | PostgreSQL | ACID, complex queries, relational data |
| **Media Storage** | S3 + CDN | Scalable object storage, global delivery |
| **Push Notifications** | APNs + FCM | Platform-native, reliable background delivery |

### Scalability Path

1. **Current (230K msg/sec):** 100 gateways, single Cassandra cluster, Redis cluster
2. **10x (2.3M msg/sec):** 1K gateways, Cassandra scaling, Kafka partitioning
3. **100x (23M msg/sec):** Multi-region, protocol optimization (protobuf), connection multiplexing

### Interview Performance Tips

1. ✅ Start with WebSocket vs HTTP long-polling vs SSE trade-off
2. ✅ Address message ordering early (Snowflake IDs)
3. ✅ Discuss online and offline delivery paths separately
4. ✅ Deep dive into group message fan-out strategies
5. ✅ Mention E2E encryption at a conceptual level
6. ✅ Discuss read receipts as a write amplification problem
7. ✅ Address failure scenarios (gateway crash, partition)

---

**End of Chat / Messaging System Design**  
[← Back to Main Index](../README.md)
