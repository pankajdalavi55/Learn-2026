# Google Search Typeahead (Autocomplete) — High-Level Design

---

## 1. Problem Statement

Design a real-time typeahead / autocomplete system (like Google Search suggestions) that returns the **top-K most relevant search suggestions** as the user types each character, with **sub-50ms latency** at massive global scale.

**Example:**
```
User types: "how to l"

Suggestions:
 1. how to lose weight
 2. how to learn python
 3. how to lower blood pressure
 4. how to learn guitar
 5. how to link aadhar with pan
```

---

## 2. Functional Requirements

| # | Requirement | Description |
|---|-------------|-------------|
| FR-1 | Prefix-based suggestions | Return top suggestions matching the typed prefix |
| FR-2 | Ranked results | Suggestions ranked by popularity / relevance / freshness |
| FR-3 | Real-time response | Suggestions update with every keystroke |
| FR-4 | Personalization (optional) | Factor in user's search history for ranking |
| FR-5 | Trending queries | Boost recently trending / breaking-news queries |
| FR-6 | Multi-language support | Handle queries in multiple languages and scripts |
| FR-7 | Offensive content filtering | Block inappropriate / harmful suggestions |
| FR-8 | Spell correction hints | Suggest corrected versions of misspelled prefixes |

---

## 3. Non-Functional Requirements

| # | Requirement | Target |
|---|-------------|--------|
| NFR-1 | Latency | p99 < 50ms (user-perceived < 100ms including network) |
| NFR-2 | Availability | 99.999% (five nines) |
| NFR-3 | Scalability | Handle 100K+ QPS globally |
| NFR-4 | Freshness | Trending queries reflected within 15–30 minutes |
| NFR-5 | Consistency | Eventual consistency is acceptable |
| NFR-6 | Fault tolerance | Graceful degradation — show stale results if backend is down |

---

## 4. Capacity Estimation

### 4.1 Traffic

| Metric | Value |
|--------|-------|
| Daily Active Users | ~1 billion |
| Searches per user per day | ~5 |
| Avg. characters typed per search | ~10 |
| Keystrokes generating a request | ~5 (with debouncing/throttling) |
| **Total daily suggestion requests** | **~25 billion** |
| **Peak QPS** | **~600,000** |
| Avg. suggestion response size | ~500 bytes (5 suggestions × ~100 bytes) |

### 4.2 Storage

| Data | Estimate |
|------|----------|
| Unique search phrases in corpus | ~5 billion |
| Avg. phrase length | ~25 characters |
| Raw phrase storage | ~125 GB |
| Trie nodes (with metadata) | ~50 GB per shard |
| Total with replication (3×) | ~500 GB – 1 TB |
| Trending data (rolling window) | ~10 GB |

### 4.3 Bandwidth

| Direction | Calculation | Result |
|-----------|-------------|--------|
| Inbound | 600K QPS × ~50 bytes/request | ~30 MB/s |
| Outbound | 600K QPS × ~500 bytes/response | ~300 MB/s |

---

## 5. Core Data Structure — Trie

A **Trie (Prefix Tree)** is the foundational data structure for typeahead. Each node represents a character, and paths from root to nodes represent prefixes.

### 5.1 Basic Trie Structure

```
              (root)
            /   |    \
           h    w     l
          /     |      \
         o      h       e
        /       |        \
       w        a         a
      /         |          \
    (how)      (what)      (learn)
    / | \       |  \
   t   i  d    i    s
  /    |   \    |     \
 o     s    o  (whatis) (whats)
 |          |
(howto)   (howdo)
```

### 5.2 Enhanced Trie Node

```
┌─────────────────────────────────┐
│         TrieNode                │
├─────────────────────────────────┤
│ char: character                  │
│ children: Map<char, TrieNode>   │
│ isEndOfPhrase: boolean          │
│ topSuggestions: List<Suggestion>│  ← Pre-computed top-K at each node
│ frequency: long                  │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│        Suggestion               │
├─────────────────────────────────┤
│ phrase: String                   │
│ score: double                    │  ← Weighted: popularity + freshness + personalization
│ lastUpdated: timestamp           │
└─────────────────────────────────┘
```

### 5.3 Why Trie? (vs alternatives)

| Approach | Lookup Time | Space | Pros | Cons |
|----------|-------------|-------|------|------|
| **Trie** | O(prefix_len) | Moderate | Fastest prefix lookup; pre-computed top-K | Memory-intensive if naive |
| HashMap (prefix → results) | O(1) | Very high | Simple | Explodes with all possible prefixes |
| SQL `LIKE 'prefix%'` | O(N) | Low | Easy to implement | Far too slow at scale |
| Elasticsearch prefix | O(log N) | Moderate | Feature-rich | Higher latency, overkill |
| Sorted array + binary search | O(log N) | Low | Compact | Slow for top-K aggregation |

**Decision**: Trie with pre-computed top-K suggestions at every node — O(1) lookup at query time.

---

## 6. High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                     │
│             Browser  |  Mobile App  |  Search Widget                     │
│     (debounce 100-200ms, cancel in-flight on new keystroke)              │
└─────────────────────────────┬────────────────────────────────────────────┘
                              │
                         ┌────▼─────┐
                         │   CDN    │  Cache popular prefix responses
                         └────┬─────┘
                              │
                      ┌───────▼────────┐
                      │  Load Balancer │  (L7, geo-routed)
                      └───────┬────────┘
                              │
               ┌──────────────┼──────────────────┐
               │              │                   │
        ┌──────▼──────┐ ┌────▼─────┐    ┌───────▼────────┐
        │ Suggestion  │ │Suggestion│    │  Suggestion    │
        │  Server 1   │ │ Server 2 │    │  Server N      │
        │ (Trie Shard)│ │(Trie Shd)│    │ (Trie Shard)   │
        └──────┬──────┘ └────┬─────┘    └───────┬────────┘
               │              │                   │
               └──────────────┼───────────────────┘
                              │
                    ┌─────────▼──────────┐
                    │  Aggregator Layer   │  Merge + rank results from shards
                    └─────────┬──────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                     │
  ┌──────▼──────┐    ┌───────▼───────┐    ┌───────▼───────┐
  │  Trie Build │    │  Data         │    │   Analytics   │
  │  Service    │    │  Collection   │    │   Pipeline    │
  │ (offline)   │    │  Service      │    │   (Spark/     │
  │             │    │ (query logs)  │    │    Flink)     │
  └──────┬──────┘    └───────┬───────┘    └───────┬───────┘
         │                    │                     │
         └────────────────────┼─────────────────────┘
                              │
                   ┌──────────▼──────────┐
                   │     DATA STORES     │
                   │  ┌───────────────┐  │
                   │  │  Query Logs   │  │  (Kafka → HDFS/S3)
                   │  │  (raw)        │  │
                   │  └───────────────┘  │
                   │  ┌───────────────┐  │
                   │  │  Phrase DB    │  │  (phrase → frequency, metadata)
                   │  │  (aggregated) │  │
                   │  └───────────────┘  │
                   │  ┌───────────────┐  │
                   │  │  Trie Snapshot│  │  (serialized trie → blob storage)
                   │  │  (S3/HDFS)   │  │
                   │  └───────────────┘  │
                   └─────────────────────┘
```

---

## 7. Component Deep Dive

### 7.1 Client-Side Optimizations

These are **critical** to keep server load manageable and UX smooth:

| Technique | Description |
|-----------|-------------|
| **Debouncing** | Wait 100–200ms after last keystroke before firing request |
| **Throttling** | Max 1 request per 150ms regardless of typing speed |
| **Cancel in-flight** | Abort previous AJAX/fetch call when new character is typed |
| **Client-side caching** | Cache prefix → results in memory (LRU). If user typed "how t" and got results, "how to" can reuse cached "how t" results filtered client-side |
| **Min prefix length** | Don't query server for prefixes < 2-3 characters |
| **Prefetch** | On focus, prefetch popular/personalized suggestions |

```
Timeline of keystrokes: h - o - w - (space) - t - o

Request sent:      ✗   ✗   ✓("how")  ✗    ✓("how t")  ✓("how to")
                          debounce         debounce     debounce
```

### 7.2 Suggestion Servers (Query Path)

Each server holds a **shard of the trie** in memory.

**Query flow:**
1. Receive prefix from load balancer
2. Traverse trie to the node matching the prefix — O(prefix length)
3. Return **pre-computed top-K suggestions** stored at that node — O(1)
4. Apply filters (offensive content, personalization boost)
5. Return JSON response

**Why in-memory?** The entire trie (or a shard) fits in RAM (50 GB per shard). Disk-based lookup would add unacceptable latency.

### 7.3 Trie Sharding Strategy

With billions of phrases, a single trie won't fit on one server. Sharding options:

| Strategy | How it Works | Pros | Cons |
|----------|-------------|------|------|
| **Range-based (by first char)** | Shard A: a-f, Shard B: g-m, etc. | Simple routing | Uneven load ("s" has far more queries than "z") |
| **Hash-based** | hash(prefix[0:2]) % N shards | Even distribution | Need to query multiple shards for aggregation |
| **Popularity-based** | Hot prefixes on dedicated shards | Optimized for hot paths | Complex rebalancing |

**Recommended: Two-level approach**
- **Level 1**: Route by first 2 characters (676 buckets for a-z × a-z)
- **Level 2**: Consistent hashing to map buckets to physical servers

```
Prefix: "how to learn"
  → Route key: "ho"
  → Consistent hash("ho") → Server 7
  → Server 7's in-memory trie handles the full traversal
```

### 7.4 Aggregator Layer

When a single query might span multiple shards (rare, but for very short prefixes like "a"), an aggregator merges results:

```
Client → LB → Aggregator
                  │
          ┌───────┼───────┐
          ▼       ▼       ▼
       Shard 1  Shard 2  Shard 3
          │       │       │
          └───────┼───────┘
                  ▼
            Merge top-K from each shard
            Re-rank by global score
            Return final top-5
```

For prefixes ≥ 2 characters, the query typically hits a **single shard** (no aggregation needed).

### 7.5 Data Collection Service

Captures raw search queries from the search engine:

```
User Search → Search Engine → Kafka Topic: "raw-queries"
                                    │
                              ┌─────▼──────┐
                              │ Flink/Spark │
                              │ Streaming   │
                              └─────┬──────┘
                                    │
                              ┌─────▼──────┐
                              │ Aggregated  │
                              │ Phrase DB   │
                              │ (phrase →   │
                              │  frequency) │
                              └────────────┘
```

**Sampling**: Don't log every single query (too expensive). Sample 1-in-100 or 1-in-1000, then extrapolate frequencies. This is sufficient for ranking.

### 7.6 Trie Build Service (Offline Pipeline)

Periodically rebuilds the trie from the aggregated phrase database.

```
┌─────────────────────────────────────────────────────────┐
│                  TRIE BUILD PIPELINE                     │
│                                                          │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────┐   │
│  │ Phrase DB │───►│ Build    │───►│ Serialized Trie  │   │
│  │ (phrases +│    │ Trie     │    │ Snapshot (S3)    │   │
│  │  scores)  │    │ In-Memory│    │                  │   │
│  └──────────┘    └──────────┘    └────────┬─────────┘   │
│                                            │             │
│                                   ┌────────▼─────────┐   │
│                                   │ Distribute to    │   │
│                                   │ Suggestion       │   │
│                                   │ Servers           │   │
│                                   └──────────────────┘   │
│                                                          │
│  Frequency: Every 15-30 minutes (for trending)           │
│  Full rebuild: Every few hours                           │
└─────────────────────────────────────────────────────────┘
```

**Build steps:**
1. Query Phrase DB for all phrases with `score > threshold`
2. Insert each phrase into trie
3. At every trie node, compute and store top-K suggestions (using a min-heap of size K)
4. Serialize trie to binary format
5. Upload to S3/HDFS
6. Suggestion servers pull the new snapshot and hot-swap

---

## 8. Ranking / Scoring Algorithm

Suggestions aren't just sorted by raw frequency. A **weighted score** balances multiple signals:

```
score(phrase) = w1 × popularity(phrase)
             + w2 × freshness(phrase)
             + w3 × personalization(phrase, user)
             + w4 × trending_boost(phrase)
             - w5 × penalty(phrase)    // spam, offensive
```

| Signal | Description | Weight |
|--------|-------------|--------|
| **Popularity** | Historical search frequency (log-normalized) | 0.5 |
| **Freshness** | Recency-weighted frequency (exponential decay) | 0.2 |
| **Trending** | Spike detection vs. baseline (z-score) | 0.15 |
| **Personalization** | User's past search history overlap | 0.1 |
| **Penalty** | Negative score for blocked/offensive content | 0.05 |

### 8.1 Freshness with Exponential Decay

```
freshness_score = Σ (count_in_window × e^(-λ × age_in_hours))
```

Where `λ` controls decay rate. Recent searches contribute more to the score.

### 8.2 Trending Detection

```
z_score = (current_hour_count - rolling_avg) / rolling_stddev

if z_score > 3.0 → mark as trending → apply trending_boost multiplier
```

This ensures breaking news (e.g., "earthquake in...") appears quickly in suggestions.

---

## 9. Handling Trending / Real-Time Queries

A pure offline trie rebuild (every few hours) is too slow for trending topics. Solution: **two-tier architecture**.

```
┌─────────────────────────────────────────────────────┐
│              SUGGESTION SERVER                       │
│                                                      │
│   ┌────────────────────┐   ┌──────────────────────┐ │
│   │  STATIC TRIE       │   │  TRENDING TRIE       │ │
│   │  (rebuilt hourly)  │   │  (updated every       │ │
│   │                    │   │   15-30 min via       │ │
│   │  Billions of       │   │   streaming)          │ │
│   │  phrases           │   │                       │ │
│   │                    │   │  Thousands of         │ │
│   │                    │   │  trending phrases     │ │
│   └────────┬───────────┘   └──────────┬───────────┘ │
│            │                          │              │
│            └──────────┬───────────────┘              │
│                       │                              │
│              ┌────────▼────────┐                     │
│              │  MERGE & RANK   │                     │
│              │  (union top-K   │                     │
│              │   from both)    │                     │
│              └─────────────────┘                     │
└─────────────────────────────────────────────────────┘
```

| Tier | Data Source | Update Frequency | Size |
|------|------------|-----------------|------|
| Static Trie | Historical query logs (batch) | Every 1-4 hours | Billions of phrases |
| Trending Trie | Real-time streaming (Flink/Kafka) | Every 15-30 min | Top ~10K trending phrases |

At query time, both tries are traversed and results are merged with trending results boosted.

---

## 10. Caching Strategy

### 10.1 Multi-Layer Caching

```
┌────────────┐     ┌────────────┐     ┌────────────┐     ┌──────────┐
│  Browser   │────►│    CDN     │────►│   Redis    │────►│  Trie    │
│  Cache     │     │  (Edge)    │     │  (Server)  │     │ (Memory) │
│ (LRU,      │     │            │     │            │     │          │
│  in-memory)│     │ popular    │     │ warm       │     │ cold     │
│            │     │ prefixes   │     │ prefixes   │     │ prefixes │
└────────────┘     └────────────┘     └────────────┘     └──────────┘
   ~80% hit          ~10% hit          ~8% hit           ~2% miss
```

### 10.2 Cache Details

| Layer | What's Cached | TTL | Hit Rate |
|-------|--------------|-----|----------|
| **Browser** | Recent prefix → results (JS LRU map) | Session-long | ~80% |
| **CDN Edge** | Top 1000 most popular prefixes | 5 min | ~10% |
| **Redis (Server-side)** | Warm prefixes (top 100K) | 2 min | ~8% |
| **In-Memory Trie** | Everything else | Until next rebuild | ~2% |

**Why 80% browser cache hit?** Users type incrementally: "h" → "ho" → "how". After "how" fetches results, the client can locally filter for "how " and "how t" without a server call.

---

## 11. Data Flow — End to End

### 11.1 Query Path (Online — Read)

```
User types "how t"
    │
    ▼
Browser debounce (150ms) → fire request
    │
    ▼
CDN Edge Cache: MISS
    │
    ▼
Load Balancer → route to nearest data center
    │
    ▼
Consistent hash("ho") → Suggestion Server #7
    │
    ▼
Server #7:
  1. Check Redis cache for "how t" → MISS
  2. Traverse in-memory trie: root → h → o → w → ' ' → t
  3. Read pre-computed top-5 at node 't'
  4. Merge with trending trie results
  5. Apply offensive filter
  6. Store in Redis (TTL: 2 min)
  7. Return response
    │
    ▼
Response (< 10ms server-side):
[
  "how to lose weight",
  "how to learn python",
  "how to tie a tie",
  "how to take screenshot",
  "how to transfer money"
]
```

### 11.2 Data Path (Offline — Write)

```
User completes search: "how to learn python"
    │
    ▼
Search Engine logs query → Kafka topic: raw-queries
    │
    ▼
Flink streaming job:
  - Sample (1:1000)
  - Aggregate: phrase → count per time window
  - Detect trending spikes (z-score)
    │
    ├──► Trending Trie (hot path, updated every 15 min)
    │
    └──► Phrase DB (HDFS/Cassandra)
             │
             ▼
         Trie Build Service (batch, every 1-4 hours)
             │
             ▼
         New trie snapshot → S3
             │
             ▼
         Suggestion Servers pull & hot-swap
```

---

## 12. Trie Update Strategy — Zero Downtime

Rebuilding the trie must not disrupt live serving.

### Approach: Blue-Green Trie Swap

```
Time T0:  Server serving Trie v1 (active)
          Background: building Trie v2

Time T1:  Trie v2 build complete, loaded into memory
          ┌──────────────────────────┐
          │  Trie v1 (serving)       │
          │  Trie v2 (standby, warm) │
          └──────────────────────────┘

Time T2:  Atomic pointer swap → Trie v2 now active
          Trie v1 scheduled for GC

Timeline: ──v1 serving────────┤──v2 serving────►
                              ^ swap (atomic, <1ms)
```

No downtime. No lock contention. Old trie is garbage-collected after swap.

---

## 13. Filtering & Safety

### 13.1 Offensive Content Pipeline

```
Phrase Candidate
    │
    ├──► Blocklist Check (exact match against known bad phrases)
    │
    ├──► ML Classifier (toxicity score, trained on flagged queries)
    │
    ├──► Regex Patterns (known offensive patterns)
    │
    └──► Human Review Queue (edge cases, appeals)
          │
          ▼
    ALLOW / BLOCK decision stored in Phrase DB
    Blocked phrases excluded during trie build
```

### 13.2 Legal & Privacy

- Remove queries that contain PII (email addresses, phone numbers, SSN patterns)
- Comply with "right to be forgotten" (GDPR) — remove specific suggestions on request
- Country-specific legal filtering (e.g., censorship laws)

---

## 14. Multi-Region Deployment

```
                    ┌──────────────────┐
                    │   Global DNS     │
                    │  (Geo-routing)   │
                    └────────┬─────────┘
                             │
           ┌─────────────────┼─────────────────┐
           │                 │                  │
    ┌──────▼──────┐  ┌──────▼──────┐  ┌───────▼──────┐
    │  US-East    │  │  EU-West    │  │  AP-South    │
    │  Region     │  │  Region     │  │  Region      │
    │             │  │             │  │              │
    │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌──────────┐ │
    │ │Sugg.Svrs│ │  │ │Sugg.Svrs│ │  │ │Sugg.Svrs │ │
    │ │(trie    │ │  │ │(trie    │ │  │ │(trie     │ │
    │ │shards)  │ │  │ │shards)  │ │  │ │shards)   │ │
    │ └─────────┘ │  │ └─────────┘ │  │ └──────────┘ │
    │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌──────────┐ │
    │ │  Redis  │ │  │ │  Redis  │ │  │ │  Redis   │ │
    │ └─────────┘ │  │ └─────────┘ │  │ └──────────┘ │
    └─────────────┘  └─────────────┘  └──────────────┘
           │                 │                  │
           └─────────────────┼──────────────────┘
                             │
                   ┌─────────▼──────────┐
                   │  Central Trie      │
                   │  Build Pipeline    │
                   │  (distributes to   │
                   │   all regions)     │
                   └────────────────────┘
```

- Each region has its own set of suggestion servers with localized trie
- Region-specific trending queries are handled locally
- Global trie snapshots are built centrally and distributed
- User is routed to nearest region via DNS geo-routing

---

## 15. Personalization (Optional Layer)

```
┌─────────────┐     ┌──────────────┐     ┌────────────────┐
│ User types   │────►│ Suggestion   │────►│ Personalization│
│ "py"         │     │ Server       │     │ Re-ranker      │
│              │     │ (top-10 from │     │                │
│              │     │  trie)       │     │ Boost if phrase│
│ User context:│     │              │     │ matches user's │
│ - past       │     └──────────────┘     │ history/profile│
│   searches   │                          └───────┬────────┘
│ - location   │                                  │
│ - language   │                          ┌───────▼────────┐
└─────────────┘                           │  Final top-5   │
                                          │  (personalized)│
                                          └────────────────┘
```

| Signal | Example |
|--------|---------|
| Search history | User often searches "python tutorial" → boost "python" suggestions |
| Location | User in Mumbai → boost "python jobs in mumbai" |
| Language | User's browser set to Hindi → include Hindi suggestions |
| Device | Mobile users → shorter suggestion phrases |

Personalized data stored per-user in a fast KV store (Redis / DynamoDB). Fetched at query time and used only for re-ranking (not for trie lookup).

---

## 16. Fault Tolerance & Graceful Degradation

| Failure Scenario | Mitigation |
|------------------|------------|
| Suggestion server goes down | Load balancer routes to replica; each shard has 3+ replicas |
| Trie build fails | Continue serving previous trie version |
| Redis cache down | Fall through to in-memory trie (slightly higher latency) |
| Kafka pipeline lag | Trending trie may be stale; static trie still serves |
| Entire region down | DNS failover to next-closest region |
| Trending trie corrupted | Circuit breaker: disable trending merge, serve static only |

---

## 17. Tech Stack Summary

| Layer | Technology |
|-------|------------|
| Client | Browser (JS fetch + LRU cache), Mobile (native HTTP) |
| CDN | Cloudflare / CloudFront |
| Load Balancer | Envoy / NGINX (L7, geo-aware) |
| Suggestion Servers | Java / Go (custom trie in-memory) |
| Cache | Redis Cluster |
| Streaming | Apache Kafka + Apache Flink |
| Batch Processing | Apache Spark |
| Phrase Storage | Cassandra / HDFS + Hive |
| Trie Snapshots | S3 / GCS (binary serialized) |
| Offensive Filtering | ML model (TensorFlow Serving) + blocklist |
| Monitoring | Prometheus + Grafana + PagerDuty |
| Deployment | Kubernetes (multi-region), ArgoCD |

---

## 18. API Design

### GET `/api/v1/suggestions`

**Request:**
```
GET /api/v1/suggestions?q=how+to+l&lang=en&limit=5
Headers:
  Authorization: Bearer <optional, for personalization>
  X-Request-ID: uuid
```

**Response (200 OK):**
```json
{
    "query": "how to l",
    "suggestions": [
        {
            "text": "how to lose weight",
            "score": 0.98,
            "trending": false
        },
        {
            "text": "how to learn python",
            "score": 0.95,
            "trending": true
        },
        {
            "text": "how to lower blood pressure",
            "score": 0.91,
            "trending": false
        },
        {
            "text": "how to learn guitar",
            "score": 0.88,
            "trending": false
        },
        {
            "text": "how to link aadhar with pan",
            "score": 0.85,
            "trending": false
        }
    ],
    "latencyMs": 8
}
```

**Headers:**
```
Cache-Control: public, max-age=120
X-Served-By: suggestion-server-7
X-Trie-Version: v2026032510
```

---

## 19. Key Optimizations Summary

| Optimization | Impact |
|-------------|--------|
| Pre-computed top-K at each trie node | Query time O(prefix_len) instead of O(N) tree traversal |
| Client-side debounce + cancellation | ~70% reduction in server QPS |
| Browser LRU cache | ~80% of keystrokes don't hit server |
| CDN caching of popular prefixes | Offloads ~10% of remaining traffic |
| Trie sharding by prefix | Single server handles each query (no scatter-gather) |
| Blue-green trie swap | Zero-downtime updates |
| Two-tier trie (static + trending) | Freshness without full rebuild cost |
| Sampling query logs (1:1000) | Manageable data pipeline; statistically accurate |

---

## 20. Key Interview Discussion Points

| Question | Answer |
|----------|--------|
| **Why trie over a database?** | O(prefix_len) lookup, pre-computed top-K → O(1) result fetch. DB can't match sub-10ms at 600K QPS. |
| **How do you handle a new trending topic?** | Two-tier: streaming pipeline updates a small trending trie every 15 min; merged at query time with static trie. |
| **How do you prevent offensive suggestions?** | Blocklist + ML toxicity classifier + human review. Filtered during trie build. |
| **How do you scale to 600K QPS?** | Shard trie by prefix range, replicate each shard 3×, multi-region deployment, aggressive caching. |
| **What happens if the trie build fails?** | Serve previous trie version. Trie snapshots are versioned in S3. Rollback is instant. |
| **How is the trie updated without downtime?** | Blue-green swap: build new trie in background, atomic pointer swap when ready. |
| **CAP trade-off?** | AP — eventual consistency is fine. Stale suggestions for a few minutes are acceptable. |
| **Why not Elasticsearch?** | ES adds network hop + serialization overhead. In-memory trie is 10× faster for pure prefix lookup. ES is better for complex search (fuzzy, faceted). |
| **How to support multiple languages?** | Separate tries per language. Route by `Accept-Language` header or user preference. Unicode-aware trie nodes. |
| **Memory estimate for trie?** | ~50 GB per shard for billions of phrases. Fits in modern server RAM (128+ GB). |

---

## 21. ASCII Summary — System at a Glance

```
         ╔═══════════════════════════════════════════════════╗
         ║          GOOGLE SEARCH TYPEAHEAD                  ║
         ╠═══════════════════════════════════════════════════╣
         ║                                                   ║
         ║   READ PATH (online, <50ms):                      ║
         ║   User → CDN → LB → Trie Server → Top-K results  ║
         ║                                                   ║
         ║   WRITE PATH (offline, batch + streaming):        ║
         ║   Query Logs → Kafka → Flink → Phrase DB          ║
         ║                         → Spark → Trie Build      ║
         ║                                   → S3 Snapshot   ║
         ║                                   → Server Swap   ║
         ║                                                   ║
         ║   KEY INSIGHT:                                    ║
         ║   Pre-compute top-K at every trie node            ║
         ║   → Query = traverse + return (no aggregation)    ║
         ║   → Rebuild offline, serve in-memory              ║
         ║                                                   ║
         ╚═══════════════════════════════════════════════════╝
```

---
