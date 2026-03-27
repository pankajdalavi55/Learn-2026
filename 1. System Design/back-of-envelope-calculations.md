# Back-of-Envelope Calculations Guide

**Navigation:** [← Back to README](./README.md)

---

Back-of-envelope calculations are one of the most underrated skills in system design interviews. They demonstrate that you think quantitatively, validate your architecture with real numbers, and can identify bottlenecks before writing a single line of code. Interviewers at FAANG+ companies explicitly evaluate this — skipping it is a common reason for "Lean No Hire."

---

## 1. Why This Matters

### What Interviewers Are Looking For

| Signal | What It Demonstrates |
|--------|---------------------|
| **Proactively doing math** | Senior-level thinking — you don't guess, you estimate |
| **Knowing key numbers** | You've built real systems and understand their limits |
| **Identifying bottlenecks** | You can spot where a design will break before it breaks |
| **Making reasonable assumptions** | You can work with incomplete data — a core engineering skill |
| **Validating architecture choices** | Your design decisions are grounded in reality, not hand-waving |

### When to Use Back-of-Envelope Math

- **After requirements gathering** — estimate scale (QPS, storage, bandwidth)
- **When choosing infrastructure** — will one database handle this? Do we need sharding?
- **When sizing caches** — how much memory do we need?
- **When discussing trade-offs** — is this cost-effective at our scale?
- **When the interviewer asks** — "How many servers would you need for this?"

---

## 2. Essential Numbers to Memorize

### Powers of 2

You don't need to memorize all of these, but the highlighted rows are critical.

| Power | Exact Value | Approximate | Common Name |
|-------|-------------|-------------|-------------|
| 2^10 | 1,024 | ~1 Thousand | **1 KB** |
| 2^20 | 1,048,576 | ~1 Million | **1 MB** |
| 2^30 | 1,073,741,824 | ~1 Billion | **1 GB** |
| 2^40 | 1,099,511,627,776 | ~1 Trillion | **1 TB** |
| 2^50 | — | ~1 Quadrillion | **1 PB** |

**Mental shortcut:** Every 10 powers of 2 ≈ 3 powers of 10 (i.e., 2^10 ≈ 10^3).

### Data Size Reference

| Data Type | Typical Size |
|-----------|-------------|
| A single character (ASCII) | 1 byte |
| A single character (Unicode/UTF-8 avg) | 2-4 bytes |
| A UUID | 36 bytes (string) / 16 bytes (binary) |
| A tweet (280 chars) | ~560 bytes |
| A typical JSON API response | 1-10 KB |
| A profile photo (compressed) | 200 KB - 1 MB |
| A high-res image | 2-5 MB |
| A 1-minute video (720p, compressed) | ~5-10 MB |
| A 1-hour video (1080p, compressed) | ~1-3 GB |
| A typical log entry | 200-500 bytes |
| A database row (typical) | 200 bytes - 2 KB |

### Latency Numbers Every Engineer Should Know

```
L1 Cache Reference:                     0.5 ns
L2 Cache Reference:                       7 ns
Mutex Lock/Unlock:                       100 ns
Main Memory Reference:                   100 ns
Compress 1KB with Snappy:             3,000 ns  =    3 μs
Send 2KB over 1 Gbps network:        20,000 ns  =   20 μs
SSD Random Read:                     150,000 ns  =  150 μs
Read 1 MB sequentially from memory:  250,000 ns  =  250 μs
Round trip within same datacenter:   500,000 ns  =  500 μs  = 0.5 ms
Read 1 MB sequentially from SSD:   1,000,000 ns  =    1 ms
HDD Seek:                          10,000,000 ns  =   10 ms
Read 1 MB sequentially from HDD:  20,000,000 ns  =   20 ms
Send packet CA → Netherlands → CA: 150,000,000 ns = 150 ms
```

**Key takeaways:**
- Memory is ~100x faster than SSD, SSD is ~10-20x faster than HDD
- Network within a datacenter (~0.5 ms) is fast; cross-region (~150 ms) is slow
- Compression is cheap — always consider it for network transfers

### Time Conversions (The Cheat Sheet)

| Time Unit | Seconds | Useful For |
|-----------|---------|------------|
| 1 day | 86,400 ≈ **10^5** | Daily traffic → QPS |
| 1 month | 2,592,000 ≈ **2.5 × 10^6** | Monthly active users → QPS |
| 1 year | 31,536,000 ≈ **3 × 10^7** | Storage growth projections |

**Simplification trick:** Use **100,000 seconds/day** (it's 86,400, but 10^5 is close enough and makes division trivial).

### Throughput References

| Component | Typical Throughput |
|-----------|--------------------|
| A single web server (e.g., Nginx) | 10K-100K requests/second |
| A single application server (Java/Spring Boot) | 500-5,000 requests/second |
| A single MySQL/PostgreSQL instance | 5K-20K queries/second (simple reads) |
| A single Redis instance | 100K-500K operations/second |
| A single Kafka broker | 100K-200K messages/second |
| A single Elasticsearch node | 5K-20K queries/second |
| Typical SSD IOPS | 10K-100K |
| Typical HDD IOPS | 100-200 |

---

## 3. The Estimation Framework (Step-by-Step)

### The 4-Step Method

```
Step 1: Clarify the scale  →  How many users? How many requests?
Step 2: Estimate traffic    →  QPS (Queries Per Second)
Step 3: Estimate storage    →  How much data over time?
Step 4: Estimate bandwidth  →  Data transfer rate
```

### Step 1: Clarify the Scale

Always start with these questions:
- **Total users?** (e.g., 500 million)
- **Daily active users (DAU)?** (typically 20-30% of total for social apps)
- **Read:write ratio?** (e.g., 10:1 for most social platforms)
- **Retention period?** (how long do we keep data?)
- **Peak vs average?** (peak is typically 2-5x average)

### Step 2: Estimate QPS

**Formula:**
```
Average QPS = (DAU × actions_per_user_per_day) / seconds_in_a_day

Peak QPS   = Average QPS × peak_multiplier (typically 2x-5x)
```

### Step 3: Estimate Storage

**Formula:**
```
Daily Storage   = daily_writes × size_per_write

Yearly Storage  = Daily Storage × 365

5-Year Storage  = Yearly Storage × 5 × replication_factor
```

### Step 4: Estimate Bandwidth

**Formula:**
```
Incoming bandwidth = Write QPS × size_per_write

Outgoing bandwidth = Read QPS × size_per_read
```

---

## 4. Worked Examples

### Example 1: Twitter-like Service

**Requirements:** Design a service where users post short messages (tweets) and read a timeline feed.

#### Given Assumptions
- 500 million total users
- 200 million DAU
- Each user posts 2 tweets/day on average
- Each user reads timeline 10 times/day, fetching 20 tweets each time
- Average tweet size: 400 bytes (text) + 200 bytes (metadata) = 600 bytes
- 20% of tweets have a media attachment (~500 KB average)
- Data retention: 5 years
- Replication factor: 3

#### Traffic Estimation

```
Write QPS (tweets):
  200M users × 2 tweets/day = 400M tweets/day
  400M / 100,000 seconds/day = 4,000 writes/second

Read QPS (timeline):
  200M users × 10 reads/day = 2B reads/day
  2B / 100,000 = 20,000 reads/second

Peak QPS (3x multiplier):
  Write Peak: ~12,000/sec
  Read Peak:  ~60,000/sec
```

**Insight:** Read-heavy system (5:1 ratio). Caching the timeline will be critical.

#### Storage Estimation

```
Text storage per day:
  400M tweets × 600 bytes = 240 GB/day

Media storage per day:
  400M × 20% = 80M media attachments
  80M × 500 KB = 40 TB/day

Total daily:  ~40 TB/day (media dominates)

Per year:     40 TB × 365 = 14.6 PB/year

5-year total: 14.6 PB × 5 = 73 PB

With 3x replication: ~220 PB
```

**Insight:** Media storage dwarfs text storage. You need object storage (S3-like), not a traditional database, for media. Text can go in a database with sharding.

#### Bandwidth Estimation

```
Incoming (writes):
  Text:  4,000 writes/sec × 600 bytes = 2.4 MB/sec
  Media: 4,000 × 0.2 × 500 KB = 400 MB/sec
  Total incoming: ~400 MB/sec ≈ 3.2 Gbps

Outgoing (reads):
  Each timeline read = 20 tweets × 600 bytes = 12 KB (text only)
  20,000 reads/sec × 12 KB = 240 MB/sec (text)
  
  If 20% of tweets have media thumbnails (~50 KB each):
  20,000 × 20 × 0.2 × 50 KB = 4 GB/sec
  Total outgoing: ~4 GB/sec ≈ 32 Gbps
```

**Insight:** Outgoing bandwidth is huge. A CDN is mandatory for serving media.

#### Summary Table

| Metric | Value |
|--------|-------|
| Write QPS | ~4,000 (peak: ~12,000) |
| Read QPS | ~20,000 (peak: ~60,000) |
| Text Storage (5yr) | ~440 TB (with replication) |
| Media Storage (5yr) | ~220 PB (with replication) |
| Incoming Bandwidth | ~3.2 Gbps |
| Outgoing Bandwidth | ~32 Gbps |

---

### Example 2: URL Shortener (like bit.ly)

**Requirements:** Shorten URLs and redirect users to the original URL.

#### Given Assumptions
- 500 million new URLs shortened per month
- Read:write ratio = 100:1
- Each URL mapping: short URL (7 chars) + long URL (average 200 chars) + metadata = ~500 bytes
- Data retention: 10 years

#### Traffic Estimation

```
Write QPS:
  500M / month = 500M / (30 × 86,400) ≈ 500M / 2.5M ≈ 200 writes/sec

Read QPS:
  200 × 100 (read:write ratio) = 20,000 reads/sec

Peak QPS (3x):
  Write Peak: ~600/sec
  Read Peak:  ~60,000/sec
```

#### Storage Estimation

```
Per month: 500M × 500 bytes = 250 GB

Per year:  250 GB × 12 = 3 TB

10 years:  30 TB

With replication (3x): 90 TB
```

**Insight:** Storage is modest. A sharded relational database or key-value store handles this easily.

#### Unique Key Space

```
If we use base62 encoding (a-z, A-Z, 0-9):
  7 characters → 62^7 = 3.5 trillion possible keys

We need: 500M/month × 12 × 10 = 60 billion keys over 10 years

3.5 trillion >> 60 billion → 7 characters is more than sufficient
```

#### Cache Estimation

```
Following the 80/20 rule (80% of traffic hits 20% of URLs):

Daily read requests: 20,000 QPS × 86,400 = ~1.7 billion/day

Cache 20% of daily URLs:
  1.7B × 0.2 = 340 million entries
  340M × 500 bytes = 170 GB

A few Redis instances (128 GB RAM each) handle this comfortably.
```

#### Summary Table

| Metric | Value |
|--------|-------|
| Write QPS | ~200 (peak: ~600) |
| Read QPS | ~20,000 (peak: ~60,000) |
| Storage (10yr, replicated) | ~90 TB |
| Cache Needed | ~170 GB |
| Key Space (7 chars, base62) | 3.5 trillion (sufficient) |

---

### Example 3: YouTube-like Video Platform

**Requirements:** Users upload and stream videos.

#### Given Assumptions
- 2 billion total users, 800 million DAU
- 5 million video uploads per day
- Average video: 300 MB (after compression, multiple resolutions)
- Average user watches 5 videos/day, each ~10 minutes
- Average video stream bitrate: 5 Mbps (1080p)
- Data retention: indefinite

#### Traffic Estimation

```
Upload QPS:
  5M videos/day / 100,000 sec/day = 50 uploads/sec

Video watch QPS:
  800M users × 5 videos/day = 4 billion views/day
  4B / 100,000 = 40,000 video-start requests/sec
```

#### Storage Estimation

```
Daily upload storage:
  5M videos × 300 MB = 1.5 PB/day

Per year:    1.5 PB × 365 = 547 PB/year ≈ 0.55 EB/year

With replication (3x): ~1.6 EB/year
```

**Insight:** At this scale, you need a custom distributed storage system (like Google's Colossus or similar). Traditional storage solutions don't work.

#### Bandwidth Estimation

```
Incoming (uploads):
  50 uploads/sec × 300 MB = 15 GB/sec = 120 Gbps

Outgoing (streaming):
  Concurrent viewers at any moment ≈ 
    800M DAU × 50 min watching/day / 1,440 min/day ≈ 28 million concurrent
  28M × 5 Mbps = 140 Tbps (Terabits per second!)
```

**Insight:** 140 Tbps outgoing bandwidth is astronomical. This is why YouTube uses a global CDN with edge servers in thousands of locations. The CDN serves >95% of traffic — origin servers would collapse otherwise.

#### Summary Table

| Metric | Value |
|--------|-------|
| Upload QPS | ~50/sec |
| Watch QPS | ~40,000/sec |
| Concurrent Viewers | ~28 million |
| Daily Storage | ~1.5 PB |
| Yearly Storage (replicated) | ~1.6 EB |
| Outgoing Bandwidth | ~140 Tbps |

---

### Example 4: Chat/Messaging Service (like WhatsApp)

**Requirements:** Real-time 1:1 and group messaging.

#### Given Assumptions
- 2 billion total users, 500 million DAU
- Average user sends 40 messages/day
- Average message size: 100 bytes (text) + 100 bytes (metadata) = 200 bytes
- 10% of messages have media (average 200 KB)
- Messages stored for 30 days on server, indefinitely in cloud backup
- Max group size: 256 members

#### Traffic Estimation

```
Message QPS:
  500M DAU × 40 messages/day = 20 billion messages/day
  20B / 100,000 = 200,000 messages/sec

Peak QPS (3x): 600,000 messages/sec
```

**Insight:** 600K messages/sec at peak is extremely high. This requires a distributed message broker and connection-based architecture (WebSockets), not HTTP polling.

#### Storage Estimation

```
Text messages per day:
  20B × 200 bytes = 4 TB/day

Media messages per day:
  20B × 10% × 200 KB = 400 TB/day

Total daily: ~400 TB/day

30-day server retention:
  400 TB × 30 = 12 PB (hot storage)
```

#### Connection Estimation

```
Concurrent connections:
  At any moment, ~30-50% of DAU may be online
  500M × 0.4 = 200 million concurrent WebSocket connections

If each server handles 50,000 connections:
  200M / 50,000 = 4,000 connection servers needed
```

#### Summary Table

| Metric | Value |
|--------|-------|
| Message QPS | ~200,000 (peak: ~600,000) |
| Daily Storage (text + media) | ~400 TB |
| 30-Day Hot Storage | ~12 PB |
| Concurrent Connections | ~200 million |
| Connection Servers Needed | ~4,000 |

---

### Example 5: Rate Limiter

**Requirements:** Limit API requests to 100 requests/user/minute.

#### Given Assumptions
- 10 million active API users
- Using a sliding window counter in Redis
- Each rate limit entry: user_id (8 bytes) + counter (4 bytes) + timestamp (8 bytes) = ~20 bytes
- Key with Redis overhead: ~100 bytes per entry

#### Memory Estimation

```
Memory needed:
  10M users × 100 bytes = 1 GB

With some overhead for Redis data structures: ~2-3 GB
```

**Insight:** A single Redis instance (typical 64-128 GB RAM) handles this trivially. Even at 100M users, it's only ~10-30 GB — still one or a few Redis nodes.

#### Throughput Check

```
If each API call requires 1 Redis operation (GET + conditional SET):
  Assume average 50 API calls/user/minute at peak
  10M × 50 / 60 seconds = ~8.3 million Redis ops/sec

A single Redis instance handles ~500K ops/sec
  → Need ~17 Redis instances (or a Redis Cluster)
```

---

## 5. Server Estimation

### How Many Application Servers?

**General formula:**
```
Servers needed = Peak QPS / QPS_per_server
```

**Typical QPS per server (Java/Spring Boot):**
- CPU-bound (computation): 500-2,000 QPS
- I/O-bound (DB calls, API calls): 2,000-10,000 QPS
- Simple proxy/routing: 10,000-50,000 QPS

**Example:** Twitter-like service with 60K peak read QPS
```
If each server handles 5,000 QPS:
  60,000 / 5,000 = 12 servers

With 30% headroom for failover: 12 × 1.3 ≈ 16 servers
```

### How Many Database Servers?

**General formula:**
```
DB servers = Peak QPS / QPS_per_DB_instance

Also consider:
  Storage-based sharding = Total data / storage_per_server
  
  Take the LARGER of the two.
```

**Example:** URL Shortener — 90 TB data, 60K read QPS
```
QPS-based:   60,000 / 15,000 = 4 DB instances
Storage-based: 90 TB / 5 TB per server = 18 DB instances

Answer: 18 shards (storage is the bottleneck, not QPS)
```

---

## 6. Quick Estimation Cheat Sheet

### Traffic Shortcuts

| Daily Requests | ≈ QPS |
|---------------|-------|
| 1 million | ~12 |
| 10 million | ~120 |
| 100 million | ~1,200 |
| 1 billion | ~12,000 |
| 10 billion | ~120,000 |

**Memorize this:** 1 million/day ≈ 12/sec. Then just scale by 10x.

### Storage Shortcuts

| Calculation | Result |
|-------------|--------|
| 1 KB × 1 billion = | 1 TB |
| 1 MB × 1 million = | 1 TB |
| 1 MB × 1 billion = | 1 PB |
| 100 bytes × 1 billion = | 100 GB |

### Byte Unit Conversions

```
1 KB  = 10^3  bytes = 1,000 bytes
1 MB  = 10^6  bytes = 1,000 KB
1 GB  = 10^9  bytes = 1,000 MB
1 TB  = 10^12 bytes = 1,000 GB
1 PB  = 10^15 bytes = 1,000 TB
```

---

## 7. Common Mistakes

### Mistake 1: Being Too Precise

**Wrong:** "We need exactly 11,574 QPS..."
**Right:** "Roughly 12,000 QPS. Let's round to ~12K for simplicity."

Back-of-envelope means **order of magnitude** — whether you need 10 servers or 100 matters; whether you need 10 or 12 does not.

### Mistake 2: Forgetting Replication

Storage estimates should always account for:
- **Data replication** (typically 3x for distributed databases)
- **Backups** (add 1-2x more)
- **Temporary data** (logs, caches, intermediate processing)

### Mistake 3: Ignoring Peak vs Average

Systems must handle **peak load**, not average load. Design for:
- **2-3x average** for most systems
- **5-10x average** for systems with viral/bursty traffic (social media, e-commerce flash sales)
- **Auto-scaling buffer** — even with auto-scaling, you need baseline capacity

### Mistake 4: Ignoring the Read/Write Ratio

Most systems are read-heavy (10:1 to 1000:1). This fundamentally shapes architecture:
- **Read-heavy** → invest in caching, read replicas, CDN
- **Write-heavy** → invest in write-optimized databases, message queues, async processing

### Mistake 5: Not Stating Assumptions

Always state your assumptions explicitly:
```
✗ "We'll have about 1000 QPS"
✓ "Assuming 100M DAU, each making ~10 requests/day, 
   that's 1B requests/day ÷ 100K seconds/day ≈ 10,000 QPS"
```

This shows your reasoning. Even if the interviewer disagrees with your assumptions, they can see your process is sound.

---

## 8. Practice Problems

Try these on your own before checking your work:

### Problem 1: Instagram-like Photo Sharing
- 1 billion total users, 500 million DAU
- Each user uploads 1 photo every 3 days (avg)
- Average photo size: 2 MB
- Each user views 50 photos/day
- Data retention: indefinite

**Estimate:** QPS (read/write), daily storage, 5-year storage, bandwidth, cache size for hot photos.

### Problem 2: Uber-like Ride Service
- 50 million DAU (riders + drivers)
- 20 million rides per day
- Each active driver sends GPS location every 3 seconds
- 5 million active drivers at peak
- Each GPS ping: 50 bytes

**Estimate:** GPS update QPS, ride request QPS, daily GPS storage, bandwidth for location updates.

### Problem 3: Notification Service
- 500 million registered devices
- Average 5 notifications per device per day
- Average notification payload: 1 KB
- 40% of notifications are sent during a 4-hour peak window

**Estimate:** Average and peak notification throughput, daily bandwidth, storage for 7-day notification history.

### Problem 4: Search Autocomplete (Typeahead)
- 1 billion search queries per day
- Average user types 4 characters before selecting a suggestion
- Each keystroke triggers an autocomplete request
- Top 10 suggestions returned per request, each ~50 bytes
- Cache the top 20% most frequent prefixes

**Estimate:** QPS for autocomplete, response payload size, cache storage needed.

---

## 9. Interview Script: How to Present Estimates

Here's how to walk through estimates during an interview:

```
"Before diving into the architecture, let me estimate the scale 
we're dealing with.

We have 200 million daily active users. If each user makes about 
10 requests per day, that's 2 billion requests per day.

Dividing by roughly 100,000 seconds in a day, we get about 20,000 
QPS on average. At peak, maybe 3x that — around 60,000 QPS.

For storage, if each request generates a 1 KB record, that's 
2 billion × 1 KB = 2 TB per day, or about 730 TB per year.

With 3x replication, we're looking at roughly 2 PB per year.

This tells me we'll need:
- A sharded database (single instance can't hold 2 PB)
- A caching layer (60K QPS is manageable but we want sub-100ms latency)
- A CDN if we're serving media

Let me design with these numbers in mind..."
```

**Why this works:**
- Shows structured thinking
- Makes assumptions explicit
- Derives architecture decisions from numbers (not the other way around)
- Takes only 1-2 minutes
- Signals senior-level quantitative reasoning

---

## 10. Quick Reference Card

```
┌─────────────────────────────────────────────────────┐
│          BACK-OF-ENVELOPE CHEAT SHEET               │
├─────────────────────────────────────────────────────┤
│                                                     │
│  TIME:  1 day ≈ 10^5 seconds                        │
│         1 month ≈ 2.5 × 10^6 seconds                │
│         1 year ≈ 3 × 10^7 seconds                   │
│                                                     │
│  QPS:   1M/day ≈ 12/sec                             │
│         1B/day ≈ 12K/sec                             │
│         Peak = Average × 2~5x                       │
│                                                     │
│  STORAGE: 1 KB × 1B = 1 TB                          │
│           1 MB × 1M = 1 TB                           │
│           Always multiply by 3 for replication       │
│                                                     │
│  SERVERS: App server ≈ 2K-10K QPS                   │
│           DB server  ≈ 5K-20K QPS (simple reads)    │
│           Redis      ≈ 100K-500K ops/sec            │
│           Kafka      ≈ 100K-200K msgs/sec           │
│                                                     │
│  NETWORK: Same DC roundtrip    ≈ 0.5 ms             │
│           Cross-region         ≈ 150 ms              │
│                                                     │
│  RULE OF THUMB:                                     │
│    - 80/20 rule for caching (cache top 20%)         │
│    - Peak = 2-5x average                            │
│    - Always add 30% headroom                        │
│    - Round aggressively — order of magnitude matters │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

*Back-of-envelope calculations are not about getting the exact right answer. They're about demonstrating that you think in systems — that you can translate vague requirements into concrete infrastructure decisions. Practice these until they feel natural, because in an interview, confidence with numbers is what separates "Lean Hire" from "Strong Hire."*
