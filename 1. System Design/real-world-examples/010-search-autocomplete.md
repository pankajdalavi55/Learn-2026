# Complete System Design: Search Autocomplete / Typeahead (Production-Ready)

> **Complexity Level:** Intermediate to Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** Google Search Suggestions, YouTube Search, Amazon Product Search, Spotify Search

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Trie Data Structure & Algorithms](#7-deep-dive-trie-data-structure--algorithms)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a search autocomplete (typeahead) system like Google's search suggestions that returns relevant query completions as the user types."

**Interviewer's Perspective:**  
This problem assesses:
- Trie data structure and prefix matching
- Ranking algorithms (popularity, freshness, personalization)
- Real-time data pipelines for updating suggestions
- Latency optimization (sub-100ms responses)
- Offline data processing at scale

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I start, let me clarify the requirements."

### 2.1 Functional Requirements

**Candidate:** "For functional requirements:
1. Should suggestions appear as the user types each character?
2. How many suggestions do we show — top 5, 10?
3. Should we rank by popularity, recency, or personalization?
4. Do we need fuzzy matching (handle typos)?
5. Do we need to filter inappropriate suggestions?
6. Should suggestions be personalized based on user history?"

**Interviewer:** "Let's support:
- Return top 5 suggestions after each keystroke
- Rank primarily by popularity, with recency boost
- Basic fuzzy matching as a stretch goal
- Filter inappropriate content
- Personalization optional, discuss approach"

**Candidate:** "Core features:
1. ✅ Prefix-based autocomplete (top 5 suggestions per prefix)
2. ✅ Ranking by popularity with recency boost
3. ✅ Content filtering (block inappropriate queries)
4. ✅ Fast response (<100ms end-to-end)
5. ✅ Optional personalization based on user history"

### 2.2 Non-Functional Requirements

**Candidate:** "For non-functional requirements:
1. What's the expected query volume?
2. How fresh do suggestions need to be — real-time trending?
3. What latency is acceptable?
4. Multi-language support needed?"

**Interviewer:**
- Scale: 5 billion searches/day, ~10K autocomplete QPS
- Freshness: trending queries reflected within 15-30 minutes
- Latency: <100ms for autocomplete response
- English first, multi-language later
- Availability: 99.9%

**Candidate:** "Summary:
- **Scale:** 10K autocomplete QPS, 100M unique queries in corpus
- **Latency:** <100ms end-to-end
- **Freshness:** Trending queries within 15-30 minutes
- **Availability:** 99.9%
- **Corpus:** 100M unique queries, updated from 5B daily searches"

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me estimate the traffic:

**Search Queries:**
- 5 billion searches/day
- Searches/sec: 5B / 86,400 ≈ **58,000 searches/sec**

**Autocomplete Requests:**
- Each search triggers ~4-5 autocomplete requests (after debouncing)
- But many users pick from suggestions early, reducing further keystrokes
- Average: ~3 autocomplete requests per search
- Autocomplete QPS: 58K × 3 ≈ **~175K requests/sec**
- Peak (3x): **~500K requests/sec**

Wait — the interviewer said 10K QPS. Let me recalibrate. Perhaps 10K is average, with peaks at 50K.

**Candidate:** "At 10K QPS average, 50K peak — that's manageable with proper caching."

### 3.2 Storage Estimation

**Candidate:** "For storage:

**Query Corpus:**
- 100M unique queries
- Average query length: 20 characters = 20 bytes
- Metadata per query (frequency, timestamp, score): 50 bytes
- Total: 100M × 70 bytes = **7 GB raw data**

**Trie Storage:**
- Trie with 100M queries: shared prefixes reduce storage
- Estimated trie size: **10-20 GB** (with top-K results at each node)
- Fits comfortably in memory on modern servers (64-128 GB RAM)

**Query Logs (for aggregation):**
- 5B queries/day × 50 bytes = **250 GB/day**
- Stored in HDFS/S3 for batch processing
- Retained for 30 days: **7.5 TB**"

### 3.3 Cache Estimation

**Candidate:** "Following the power law, a small fraction of prefixes covers most requests:

- Top 100K prefixes handle ~80% of autocomplete requests
- Each prefix → 5 suggestions × 50 bytes = 250 bytes
- Cache size: 100K × 250 bytes = **25 MB** — trivially small
- Even caching 10M prefixes: **2.5 GB** — fits in Redis easily"

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌─────────────┐
│   Client    │
│ (Browser/   │
│  Mobile)    │
│             │
│ Debounce    │
│ 200-300ms   │
└──────┬──────┘
       │ HTTPS
       ▼
┌─────────────────────────────────────┐
│      CDN / Edge Cache               │
│   (Cache popular prefixes)          │
└──────────────┬──────────────────────┘
               │ Cache miss
               ▼
┌─────────────────────────────────────┐
│      Load Balancer                  │
└──────────────┬──────────────────────┘
               │
       ┌───────┴────────┐
       ▼                ▼
┌─────────────┐  ┌─────────────┐
│ Autocomplete│  │ Autocomplete│
│  Service 1  │  │  Service N  │
│             │  │             │
│ In-Memory   │  │ In-Memory   │
│ Trie        │  │ Trie        │
└──────┬──────┘  └──────┬──────┘
       │                │
       └───────┬────────┘
               │ (Trie Updates)
               │
┌──────────────┴──────────────────────┐
│         Offline Pipeline            │
│                                     │
│  Query Logs → Kafka → Aggregation   │
│  Service (Spark/Flink) → Trie       │
│  Builder → Deploy to Servers        │
│                                     │
│  ┌──────────┐  ┌─────────────────┐  │
│  │  Query   │  │  Aggregation    │  │
│  │  Logs    │  │  Service        │  │
│  │ (Kafka/  │  │  (Spark/Flink)  │  │
│  │  HDFS)   │  │                 │  │
│  └──────────┘  └─────────────────┘  │
│                                     │
│  ┌──────────┐  ┌─────────────────┐  │
│  │  Trie    │  │  Content        │  │
│  │  Builder │  │  Filter         │  │
│  └──────────┘  └─────────────────┘  │
└─────────────────────────────────────┘
       │
       ▼
┌─────────────────┐  ┌────────────────┐
│ Personalization │  │  Trending      │
│ Service         │  │  Service       │
│ (User History)  │  │  (Real-Time)   │
└─────────────────┘  └────────────────┘
```

### 4.2 API Design

**Candidate:** "Simple API for autocomplete:

**Autocomplete Request:**
```http
GET /api/autocomplete?prefix=how+to+le&userId=user_123&limit=5

Response (200 OK):
{
  "prefix": "how to le",
  "suggestions": [
    { "query": "how to learn python", "score": 98500 },
    { "query": "how to learn guitar", "score": 72300 },
    { "query": "how to learn spanish", "score": 65100 },
    { "query": "how to learn coding", "score": 58200 },
    { "query": "how to learn excel", "score": 45600 }
  ],
  "took_ms": 12
}
```

**Headers:**
```
Cache-Control: public, max-age=300  (5 minute CDN cache)
```
"

### 4.3 Data Flow

**Candidate:** "Two main flows:

**Flow 1: Serving Autocomplete (Real-Time Path)**
1. User types 'how to le' → client debounces (300ms) → sends request
2. CDN checks cache → if hit, return immediately (0ms backend)
3. If CDN miss → Load Balancer → Autocomplete Service
4. Autocomplete Service performs prefix lookup in in-memory trie
5. Trie returns top 5 queries for prefix 'how to le'
6. Optionally merge with personalized results
7. Return JSON response, CDN caches it

**Flow 2: Trie Building (Offline Pipeline)**
1. Query logs stream into Kafka from search service
2. Aggregation service (Spark Streaming / Flink) counts query frequencies
3. Apply time-decay: recent queries weighted higher
4. Content filter removes inappropriate queries
5. Trie Builder constructs optimized trie with top-K at each node
6. New trie deployed to Autocomplete Services (rolling update)
7. Cycle: every 15-30 minutes for trending, daily for full rebuild"

---

## 5. Deep Dive: Core Components

### 5.1 Autocomplete Service

**Candidate:** "The serving layer is optimized for speed:

**Design:**
- Stateless service (trie loaded from shared storage on startup)
- Each instance holds the complete trie in memory (~15 GB)
- Prefix lookup: O(L) where L = prefix length — sub-millisecond
- No database calls on the hot path

**Technology:** Go or C++ for maximum performance
- Go: simple, concurrent, garbage-collected
- C++: zero overhead, manual memory management

**Request Flow:**
```go
func handleAutocomplete(w http.ResponseWriter, r *http.Request) {
    prefix := r.URL.Query().Get("prefix")
    prefix = strings.ToLower(strings.TrimSpace(prefix))

    if len(prefix) < 1 {
        http.Error(w, "prefix too short", 400)
        return
    }

    suggestions := trie.GetTopK(prefix, 5)

    // Optional: merge personalized results
    userId := r.URL.Query().Get("userId")
    if userId != "" {
        personal := personalizationService.GetSuggestions(userId, prefix, 3)
        suggestions = mergeAndRerank(suggestions, personal, 5)
    }

    json.NewEncoder(w).Encode(suggestions)
}
```
"

### 5.2 Aggregation Service

**Candidate:** "Processes raw query logs into frequency counts:

**Batch Pipeline (Daily - Full Rebuild):**
```
HDFS (30 days of logs) → Spark Job → Query Frequencies → Trie Builder
```

**Streaming Pipeline (Real-Time - Trending):**
```
Kafka (live queries) → Flink → Sliding Window Counts (15 min) → Trending Trie Updater
```

**Frequency Calculation with Time Decay:**
```python
import math

def calculate_score(query, occurrences):
    score = 0
    now = time.time()
    for timestamp in occurrences:
        hours_ago = (now - timestamp) / 3600
        # Exponential decay: half-life of 7 days (168 hours)
        decay = math.exp(-0.693 * hours_ago / 168)
        score += decay
    return score
```
"

### 5.3 Content Filter

**Candidate:** "Removes inappropriate, offensive, or harmful queries:

- Blocklist of explicit terms and phrases
- ML-based toxicity classifier
- Manual review queue for borderline cases
- Legal/compliance filtering (DMCA, right to be forgotten)
- Applied during trie building, not at serving time (too slow)
"

### 5.4 Trie Builder

**Candidate:** "Constructs the optimized trie from aggregated data:

1. Read query-frequency pairs from aggregation output
2. Insert each query into trie with its score
3. At each internal node, cache top-K results (avoid DFS at query time)
4. Serialize trie to binary format (memory-mapped file)
5. Upload to shared storage (S3)
6. Autocomplete Services pull new trie and hot-swap
"

---

## 6. Deep Dive: Database Design

### 6.1 Query Logs (Kafka + HDFS/S3)

```
Raw Log Schema:
{
    "query": "how to learn python",
    "userId": "user_123",       // anonymized or null
    "timestamp": 1704067200,
    "resultCount": 1520000,
    "clickPosition": 2,         // which result they clicked
    "sessionId": "sess_abc",
    "locale": "en-US",
    "device": "mobile"
}
```

### 6.2 Query Aggregates (PostgreSQL or Redis)

```sql
CREATE TABLE query_aggregates (
    query_text VARCHAR(200) PRIMARY KEY,
    frequency BIGINT NOT NULL,
    decayed_score DOUBLE PRECISION NOT NULL,
    first_seen TIMESTAMP NOT NULL,
    last_seen TIMESTAMP NOT NULL,
    is_filtered BOOLEAN DEFAULT FALSE,
    category VARCHAR(50)
);

CREATE INDEX idx_score ON query_aggregates(decayed_score DESC);
CREATE INDEX idx_prefix ON query_aggregates USING gin(query_text gin_trgm_ops);
```

### 6.3 User Search History (Cassandra or DynamoDB)

```sql
CREATE TABLE user_search_history (
    user_id UUID,
    query_text TEXT,
    searched_at TIMESTAMP,
    clicked_url TEXT,
    PRIMARY KEY (user_id, searched_at)
) WITH CLUSTERING ORDER BY (searched_at DESC);
```

### 6.4 Trending Queries (Redis)

```
ZSET trending:global → { "query1": score1, "query2": score2, ... }
ZSET trending:en-US → { ... }

Updated every 15 minutes by streaming pipeline
TTL: 1 hour (auto-expire stale trending data)
```

---

## 7. Deep Dive: Trie Data Structure & Algorithms

### 7.1 Basic Trie Implementation

**Candidate:** "The trie (prefix tree) is the heart of the system:

```python
class TrieNode:
    def __init__(self):
        self.children = {}           # char → TrieNode
        self.is_end = False          # marks end of a complete query
        self.frequency = 0           # search frequency
        self.top_k = []              # cached top-K results for this prefix

class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, query, frequency):
        node = self.root
        for char in query:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.is_end = True
        node.frequency = frequency

    def search_prefix(self, prefix):
        node = self.root
        for char in prefix:
            if char not in node.children:
                return []  # no suggestions
            node = node.children[char]
        # Return cached top-K if available
        if node.top_k:
            return node.top_k
        # Otherwise, DFS to find all completions
        return self._dfs_top_k(node, prefix, k=5)

    def _dfs_top_k(self, node, prefix, k):
        results = []
        if node.is_end:
            results.append((prefix, node.frequency))

        for char, child in node.children.items():
            results.extend(self._dfs_top_k(child, prefix + char, k))

        results.sort(key=lambda x: -x[1])
        return results[:k]
```

**Time Complexity:**
- Insert: O(L) where L = query length
- Prefix search (with cached top-K): O(P) where P = prefix length
- Prefix search (without cache, DFS): O(P + N) where N = nodes in subtree

**Space Complexity:**
- Worst case: O(N × K) where N = total characters, K = alphabet size
- With top-K caching: additional O(N × K_results) for cached suggestions"

### 7.2 Compressed Trie (Radix Tree)

**Candidate:** "Optimization: merge single-child chains to reduce nodes:

```
Standard Trie:          Compressed Trie (Radix):
    h                       h
    |                       |
    o                      "ow to"
    |                      /    \
    w                   "learn"  "cook"
    |
    (space)
    |
    t
    |
    o
    |
    (space)
    |
   / \
  l   c
  |   |
  e   o
  |   |
  a   o
  |   |
  r   k
  |
  n
```

**Benefits:**
- 40-60% fewer nodes for natural language queries
- Less memory, faster traversal
- Each edge stores a string instead of single character

```python
class CompressedTrieNode:
    def __init__(self):
        self.children = {}           # prefix_string → CompressedTrieNode
        self.is_end = False
        self.frequency = 0
        self.top_k = []

    def find_child(self, char):
        for prefix, child in self.children.items():
            if prefix[0] == char:
                return prefix, child
        return None, None
```
"

### 7.3 Pre-Computed Top-K at Each Node

**Candidate:** "The most important optimization — avoid DFS at query time:

```python
def build_trie_with_top_k(queries_with_scores, k=5):
    trie = Trie()

    # Phase 1: Insert all queries
    for query, score in queries_with_scores:
        trie.insert(query, score)

    # Phase 2: Bottom-up propagation of top-K
    def propagate_top_k(node, prefix):
        candidates = []

        # If this node is an end of a query, include it
        if node.is_end:
            candidates.append((prefix, node.frequency))

        # Collect top-K from all children
        for char, child in node.children.items():
            propagate_top_k(child, prefix + char)
            candidates.extend(child.top_k)

        # Keep only top-K candidates
        candidates.sort(key=lambda x: -x[1])
        node.top_k = candidates[:k]

    propagate_top_k(trie.root, '')
    return trie
```

**With this optimization:**
- Query time: O(P) — just traverse to the prefix node and return cached top-K
- No DFS needed at serving time
- Trade-off: more memory (store top-K at every node), longer build time
- **This is what makes <10ms serving latency possible**"

### 7.4 Ranking Strategy

**Candidate:** "Multiple signals feed into the ranking score:

**Signal 1: Raw Frequency**
```
score = total_search_count
```
Problem: biased toward old, historically popular queries.

**Signal 2: Time-Decayed Frequency**
```python
def time_decayed_score(query_occurrences, half_life_hours=168):
    score = 0
    now = time.time()
    for ts in query_occurrences:
        hours_ago = (now - ts) / 3600
        decay = 2 ** (-hours_ago / half_life_hours)
        score += decay
    return score
```
Better: recent queries get higher weight; old queries naturally fade.

**Signal 3: Combined Score**
```python
def combined_score(query):
    freq_score = time_decayed_frequency(query)
    ctr_score = click_through_rate(query)        # how often users pick this suggestion
    freshness = recency_boost(query)              # boost for very recent queries

    return (0.6 * freq_score) + (0.3 * ctr_score) + (0.1 * freshness)
```

**Signal 4: Personalized Ranking**
```python
def personalized_score(query, user_id, global_score):
    user_freq = get_user_query_frequency(user_id, query)
    personal_weight = 0.3

    if user_freq > 0:
        personal_score = user_freq * 10  # boost personal history
        return personal_weight * personal_score + (1 - personal_weight) * global_score
    return global_score
```
"

### 7.5 Trie Update Strategy

**Candidate:** "Keeping the trie fresh is critical for trending queries:

**Strategy 1: Full Rebuild (Simple)**
```
Schedule: Every 4-6 hours
Process:
  1. Spark job aggregates all query logs (last 30 days)
  2. Apply time decay scoring
  3. Build complete new trie
  4. Serialize to binary file, upload to S3
  5. Autocomplete servers pull new trie and hot-swap

Pros: Simple, consistent, no incremental complexity
Cons: 4-6 hour delay for new trending queries
```

**Strategy 2: Incremental Updates (Real-Time Trending)**
```
For trending queries (appearing in last 15 min):
  1. Flink streaming job counts queries in sliding window
  2. If query frequency exceeds threshold, it's "trending"
  3. Push trending queries to autocomplete servers via Kafka
  4. Servers merge trending queries into their trie in-memory

Code:
async function handleTrendingUpdate(trendingQueries) {
    for (const { query, score } of trendingQueries) {
        trie.insert(query, score);
        // Update top-K along the path from root to this query
        trie.updateTopKPath(query);
    }
}
```

**Strategy 3: Hybrid (Recommended)**
```
- Full rebuild: every 6 hours (comprehensive, all data)
- Incremental trending: every 15 minutes (hot queries only)
- Result: trending queries appear within 15 min, full corpus refreshed every 6 hours
```
"

### 7.6 Fuzzy Matching (Spell Correction)

**Candidate:** "Handle typos with edit distance:

```python
def fuzzy_search(trie, query, max_distance=1):
    results = []

    def search_recursive(node, prefix, remaining, distance):
        if distance > max_distance:
            return
        if not remaining:
            if node.is_end:
                results.append((prefix, node.frequency, distance))
            # Also check top_k at this node
            for suggestion, score in node.top_k:
                results.append((suggestion, score, distance))
            return

        char = remaining[0]
        rest = remaining[1:]

        for child_char, child_node in node.children.items():
            if child_char == char:
                # Exact match, no edit cost
                search_recursive(child_node, prefix + child_char, rest, distance)
            else:
                # Substitution (cost 1)
                search_recursive(child_node, prefix + child_char, rest, distance + 1)

            # Insertion (skip a character in the trie)
            search_recursive(child_node, prefix + child_char, remaining, distance + 1)

        # Deletion (skip a character in the query)
        search_recursive(node, prefix, rest, distance + 1)

    search_recursive(trie.root, '', query, 0)

    # Sort by score (prefer exact matches, then by frequency)
    results.sort(key=lambda x: (x[2], -x[1]))
    return results[:5]
```

**Alternative: BK-Tree for spell correction**
- Pre-build BK-tree of all queries
- Query with edit distance threshold
- O(log N) average lookup time
- Use for 'did you mean' suggestions when trie returns no results
"

### 7.7 Client-Side Optimizations

**Candidate:** "The client is equally important for perceived performance:

```javascript
class AutocompleteClient {
    constructor() {
        this.cache = new Map();     // prefix → suggestions
        this.debounceTimer = null;
        this.abortController = null;
    }

    onInput(prefix) {
        // 1. Check local cache first
        const cached = this.cache.get(prefix);
        if (cached) {
            this.renderSuggestions(cached);
            return;
        }

        // 2. Optimistic: check if a parent prefix's results contain matches
        for (let i = prefix.length - 1; i >= 1; i--) {
            const parentPrefix = prefix.substring(0, i);
            const parentResults = this.cache.get(parentPrefix);
            if (parentResults) {
                const filtered = parentResults.filter(s =>
                    s.query.startsWith(prefix)
                );
                if (filtered.length > 0) {
                    this.renderSuggestions(filtered);
                    // Still fetch from server for better results
                    break;
                }
            }
        }

        // 3. Debounce: wait 200ms before making API call
        clearTimeout(this.debounceTimer);
        if (this.abortController) this.abortController.abort();

        this.debounceTimer = setTimeout(async () => {
            this.abortController = new AbortController();
            try {
                const response = await fetch(
                    `/api/autocomplete?prefix=${encodeURIComponent(prefix)}&limit=5`,
                    { signal: this.abortController.signal }
                );
                const data = await response.json();
                this.cache.set(prefix, data.suggestions);
                this.renderSuggestions(data.suggestions);
            } catch (err) {
                if (err.name !== 'AbortError') console.error(err);
            }
        }, 200);
    }

    renderSuggestions(suggestions) {
        // Highlight matching prefix in bold
        // Render dropdown below search input
    }
}
```

**Key client optimizations:**
1. **Debouncing:** 200-300ms delay prevents request per keystroke
2. **Local caching:** Cache recent prefix results (reduces server calls by ~50%)
3. **Request cancellation:** Abort in-flight request when new keystroke arrives
4. **Optimistic rendering:** Use parent prefix results while waiting for server
5. **Prefetching:** When user pauses, prefetch next likely characters"

---

## 8. Scaling Strategies

### 8.1 Current Bottlenecks

**Candidate:** "At 10K QPS:

1. **Autocomplete Service:** Each server handles ~10K requests/sec easily (in-memory trie lookup is sub-ms). 2-3 servers with replication is sufficient.
2. **CDN:** Popular prefixes cached at edge, absorbing 80% of traffic.
3. **Trie Size:** 15 GB fits in memory on a single server.
4. **Aggregation Pipeline:** Spark job runs in batch, not a bottleneck."

### 8.2 Scaling to 100K QPS

**Candidate:** "At 10x traffic:

**Step 1: Add more Autocomplete Service replicas**
- Each server handles the full trie
- Load balancer distributes requests round-robin
- 10-20 servers handle 100K QPS easily

**Step 2: Improve CDN caching**
- Cache top 1M prefixes at CDN edge
- Expected cache hit rate: 90%+ (power law distribution)
- Only 10K QPS reaches origin servers

**Step 3: Regional deployment**
- Deploy autocomplete servers in each region (US, EU, Asia)
- Reduces latency for global users"

### 8.3 Scaling to 1M+ QPS (Google-Scale)

**Candidate:** "At Google scale, the trie becomes too large for a single machine:

**Trie Sharding by Prefix Range:**
```
Shard 1: prefixes starting with a-f
Shard 2: prefixes starting with g-m
Shard 3: prefixes starting with n-s
Shard 4: prefixes starting with t-z, 0-9
```

**Routing:**
```javascript
function getShardForPrefix(prefix) {
    const firstChar = prefix[0].toLowerCase();
    if (firstChar >= 'a' && firstChar <= 'f') return 'shard-1';
    if (firstChar >= 'g' && firstChar <= 'm') return 'shard-2';
    if (firstChar >= 'n' && firstChar <= 's') return 'shard-3';
    return 'shard-4';
}
```

**Each shard:**
- Holds a sub-trie for its prefix range
- Replicated 3x for redundancy and read scaling
- Independent scaling based on traffic distribution

**Challenge: Uneven distribution**
- Prefix 's' has way more queries than prefix 'x'
- Solution: dynamic sharding based on traffic, not just alphabetical"

---

## 9. Failure Scenarios & Mitigation

### 9.1 Trie Server Failure

**Scenario:** An autocomplete server crashes or becomes unresponsive.

**Impact:** Slightly reduced capacity.

**Mitigation:**
- Load balancer health checks detect failure in 10 seconds
- Traffic rerouted to healthy servers (each holds complete trie)
- Auto-scaling launches replacement
- No data loss (trie is rebuilt from shared storage)

### 9.2 Stale Trie (Aggregation Pipeline Failure)

**Scenario:** Spark job fails, trie isn't updated for 24+ hours.

**Impact:** Missing recent trending queries, stale suggestions.

**Mitigation:**
- Monitor trie age (alert if >12 hours since last build)
- Fall back to last known good trie
- Trending overlay (Flink) provides partial freshness independently
- Stale suggestions are still useful (most queries are evergreen)

### 9.3 Cache Stampede on Trie Deployment

**Scenario:** New trie deployed, CDN cache expires, all requests hit origin.

**Impact:** Thundering herd on autocomplete servers.

**Mitigation:**
```javascript
// Staggered CDN cache invalidation
async function deployNewTrie() {
    for (const region of regions) {
        await deployToRegion(region);
        // Stagger by 5 minutes per region
        await sleep(5 * 60 * 1000);
    }
}

// CDN: jittered cache TTL
// Instead of all entries expiring at once:
// TTL = 300 + random(0, 60)  (5 min + jitter)
```

### 9.4 Inappropriate Query Surfacing

**Scenario:** An offensive or harmful query appears in suggestions.

**Impact:** PR disaster, user trust damage.

**Mitigation:**
- Pre-filter during trie building (blocklist + ML classifier)
- Real-time blocklist: can push emergency removals to servers within minutes
- Human review queue for borderline cases
- Incident response: ability to wipe a specific suggestion within 5 minutes
```javascript
// Emergency removal endpoint (admin only)
app.post('/admin/block-suggestion', auth, async (req, res) => {
    const { query } = req.body;
    await redis.sadd('blocked_queries', query);
    // Servers check blocked list before returning results
    // Full removal in next trie build
});
```

### 9.5 Personalization Service Timeout

**Scenario:** Personalization service is slow or down.

**Impact:** Personalized suggestions unavailable.

**Mitigation:**
- Circuit breaker: if personalization fails, return global-only results
- Personalization has a strict 20ms timeout
- Global suggestions are always the fallback and are always correct
```javascript
async function getSuggestions(prefix, userId) {
    const globalResults = trie.getTopK(prefix, 5);  // always fast

    try {
        const personalResults = await withTimeout(
            personalizationService.get(userId, prefix),
            20  // 20ms timeout
        );
        return mergeResults(globalResults, personalResults);
    } catch (err) {
        return globalResults;  // graceful degradation
    }
}
```

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

**Candidate:** "For autocomplete:

**Application Metrics (RED):**
1. **Rate:** Autocomplete requests/sec, CDN hit rate
2. **Errors:** Error rate (timeouts, trie not loaded)
3. **Duration:** End-to-end latency (p50, p95, p99)

**Business Metrics:**
- Suggestion click-through rate (CTR)
- Position of clicked suggestion (1st, 2nd, 3rd...)
- Zero-results rate (prefix with no suggestions)
- Query coverage (% of searches that used autocomplete)

**Infrastructure Metrics:**
- Trie size (memory usage per server)
- Trie age (time since last build)
- Aggregation pipeline lag
- CDN cache hit ratio

**Example Dashboard (Grafana):**
```
Row 1: Traffic & Latency
- [Graph] Autocomplete QPS
- [Heatmap] Response latency distribution
- [Gauge] CDN cache hit rate (%)

Row 2: Quality
- [Graph] Suggestion CTR over time
- [Graph] Zero-result rate (%)
- [Graph] Avg suggestion position clicked

Row 3: Data Pipeline
- [Gauge] Trie age (hours since last build)
- [Graph] Aggregation pipeline throughput
- [Graph] Trending queries detected per hour

Row 4: Infrastructure
- [Graph] Memory usage per server
- [Graph] Active servers / replicas
- [Graph] Query log ingestion rate
```
"

### 10.2 Alerting Rules

```yaml
alert: HighAutocompleteLatency
expr: histogram_quantile(0.99, autocomplete_latency_seconds) > 0.1
for: 5m
severity: critical
message: "Autocomplete p99 latency above 100ms"

alert: StaleTrie
expr: trie_age_hours > 12
for: 10m
severity: warning
message: "Trie not updated in 12+ hours"

alert: LowCTR
expr: autocomplete_ctr < 0.3
for: 30m
severity: warning
message: "Autocomplete CTR dropped below 30%"

alert: HighZeroResultRate
expr: autocomplete_zero_results_rate > 0.1
for: 15m
severity: warning
message: "More than 10% of prefixes returning zero results"
```

### 10.3 Logging

```javascript
logger.info('autocomplete_served', {
    prefix: 'how to le',
    prefixLength: 9,
    numResults: 5,
    latencyMs: 3,
    cacheHit: false,
    userId: 'user_123',   // anonymized
    selectedSuggestion: null  // logged when user clicks
});

logger.info('suggestion_clicked', {
    prefix: 'how to le',
    selectedQuery: 'how to learn python',
    position: 1,
    userId: 'user_123'
});
```

---

## 11. Advanced Features

### 11.1 Trending Searches

**Candidate:** "Real-time trending detection with Flink:

```python
# Flink windowed aggregation
class TrendingDetector:
    def process(self, query_stream):
        # Sliding window: last 15 minutes, slide every 1 minute
        windowed = query_stream \
            .key_by(lambda q: q.query_text) \
            .window(SlidingEventTimeWindows.of(
                Time.minutes(15), Time.minutes(1)
            )) \
            .aggregate(CountAggregate())

        # Compare current count vs baseline (last 24h avg)
        trending = windowed.filter(
            lambda count, baseline:
                count.current > baseline.avg * 3  # 3x spike = trending
        )

        return trending
```

Display trending searches when the search box is empty (before user types anything)."

### 11.2 Entity-Aware Suggestions

**Candidate:** "Instead of just query strings, suggest entities:

```json
{
  "suggestions": [
    { "type": "query", "text": "how to learn python", "score": 98500 },
    { "type": "entity", "text": "Python (programming language)", "entityId": "Q123", "image": "..." },
    { "type": "recent", "text": "python tutorial for beginners", "source": "history" }
  ]
}
```

Merge entity index (knowledge graph) with query trie. Separate entity index for people, places, products, etc."

### 11.3 Query Completion vs. Query Correction

```
User types: "pythn tutrial"

Query Completion: "python tutorial for beginners" (complete the query)
Query Correction: "Did you mean: python tutorial" (fix the typo)

Implementation:
1. First, try exact prefix match in trie
2. If no results or few results, run fuzzy match
3. If fuzzy match finds better alternatives, show "Did you mean..."
4. Return both completions and corrections, labeled differently
```

### 11.4 Contextual Suggestions

**Candidate:** "Enhance suggestions based on context:

- **Location:** User in NYC → 'restaurants' suggests 'restaurants in new york'
- **Time:** Evening → 'dinner' ranks higher than 'lunch'
- **Device:** Mobile → shorter suggestions, touch-friendly
- **Session context:** User searched 'python', next query suggests Python-related topics
- **Language:** Detect input language, return suggestions in that language"

### 11.5 Voice Search Autocomplete

**Candidate:** "For voice input:
- Convert speech to text (STT) in real-time
- Feed partial transcript to autocomplete API
- Show suggestions as user speaks
- Handle homophones (their/there/they're) by returning all variants"

---

## 12. Interview Q&A

### Q1: How do you handle billions of queries to build the trie?

**Answer:**
We use a two-phase approach:

1. **Batch processing (Spark/MapReduce):** Process 30 days of query logs (stored in HDFS/S3). Map phase: emit (query, 1) for each log entry. Reduce phase: sum counts per query. Apply time decay. Output: query → score pairs.

2. **Trie construction:** Read the aggregated query-score pairs (100M entries), insert into trie, propagate top-K to each internal node. This runs on a single large-memory machine (128 GB RAM) since the trie fits in memory.

Total build time: ~2 hours for the full pipeline. Runs every 6 hours. Trending queries are handled separately via streaming (Flink) for real-time updates every 15 minutes.

### Q2: How do you update the trie in real-time for trending queries?

**Answer:**
We use a hybrid approach:
- **Base trie:** Full rebuild every 6 hours from batch pipeline
- **Trending overlay:** Flink streaming job detects query frequency spikes in 15-minute sliding windows. If a query's frequency is >3x its baseline, it's marked as trending. Trending queries are pushed to autocomplete servers via Kafka, and merged into the in-memory trie.

The merge is lightweight: insert the trending query and update top-K lists along the path. This avoids full rebuilds while keeping suggestions fresh.

### Q3: Trie vs inverted index — when would you use each for autocomplete?

**Answer:**
- **Trie:** Best for pure prefix matching. O(P) lookup where P is prefix length. Memory-efficient for prefix queries. Ideal when suggestions are complete query strings.

- **Inverted index (Elasticsearch):** Best for substring matching and full-text search. Supports fuzzy matching natively. Better when suggestions come from a document corpus (product names, article titles).

- **Hybrid:** Use trie for query autocomplete (fast prefix matching), inverted index for entity search (product names that match anywhere in the string, not just prefix).

For Google-style search autocomplete, trie is the right choice. For Amazon product search, an inverted index with edge n-gram tokenization may be better.

### Q4: How do you handle multi-language autocomplete?

**Answer:**
Multi-language introduces several challenges:

1. **Character encoding:** Trie must be Unicode-aware. CJK languages (Chinese, Japanese, Korean) have thousands of characters — trie nodes need hash maps instead of fixed arrays.

2. **Tokenization:** English is space-delimited, but CJK languages are not. Use language-specific tokenizers (e.g., ICU segmenter, MeCab for Japanese).

3. **Script detection:** Determine input language from the first few characters (Latin vs Cyrillic vs CJK) and route to language-specific trie.

4. **Separate tries per language:** Build separate tries for each language with language-specific aggregation. Route based on detected script.

5. **Transliteration:** Support Romanized input for non-Latin scripts (e.g., "nihongo" → "日本語" for Japanese).

### Q5: How do you personalize suggestions without compromising latency?

**Answer:**
Personalization must be fast (<20ms). Approach:

1. **Pre-computed personal suggestions:** For active users, maintain a small personal trie (top 100 recent queries) in Redis, keyed by user_id. At query time, merge personal results with global results.

2. **Lightweight merge:** Global trie returns top 5 results. Personal cache returns top 3 matches. Merge with weighted scoring: `final_score = 0.7 * global + 0.3 * personal`.

3. **Timeout protection:** Personalization has a strict 20ms timeout. If it exceeds that, return global-only results.

4. **Privacy:** Personal history is never shared in global suggestions. Only the user sees their own personalized results.

### Q6: How do you prevent inappropriate/offensive suggestions?

**Answer:**
Multi-layered content filtering:

1. **Blocklist:** Explicit list of banned terms and phrases, applied during trie building
2. **ML classifier:** Toxicity model (e.g., Perspective API) scores each query during aggregation. Queries above threshold are filtered.
3. **Manual review:** Edge cases flagged for human review
4. **Emergency removal:** Admin API to instantly block a suggestion, with server-side blocklist checked at serving time
5. **Frequency threshold:** Only surface queries searched by enough unique users (prevents one user from injecting offensive suggestions)
6. **Legal compliance:** GDPR right to be forgotten, DMCA takedowns

### Q7: How do you handle prefix ambiguity (same prefix, very different intents)?

**Answer:**
Example: prefix "ap" could mean "apple", "apartment", "application", "apex legends"

Solutions:
1. **Diversified ranking:** Don't just rank by frequency. Ensure category diversity in top 5 results — pick top result from each category, then fill remaining by score.
2. **Context signals:** Use session context (previous queries), location, and user history to disambiguate.
3. **Query clustering:** Group queries by semantic similarity during aggregation. Ensure top-K includes representatives from different clusters.

### Q8: How would you implement "did you mean" spell correction?

**Answer:**
Two-phase approach:

1. **Detection:** When prefix has no trie matches or very few suggestions, trigger spell correction.

2. **Correction:** Use a BK-tree (Burkhard-Keller tree) built from the query corpus. Query with edit distance ≤ 2. BK-tree provides O(log N) average lookup.

```python
def did_you_mean(query, bk_tree, max_distance=2):
    candidates = bk_tree.search(query, max_distance)
    if not candidates:
        return None
    # Sort by: edit distance (ascending), then frequency (descending)
    candidates.sort(key=lambda c: (c.distance, -c.frequency))
    return candidates[0].query
```

Alternative: Norvig's spell corrector approach — generate all possible edits (insertions, deletions, substitutions, transpositions) and look them up in the trie. For edit distance 1, this generates ~50L candidates (manageable). For distance 2, use pruning.

---

## 13. Production Checklist

### 13.1 Pre-Launch

- [ ] **Load Testing:** Simulate 50K QPS with realistic prefix distribution
- [ ] **Latency Testing:** Verify p99 < 100ms end-to-end
- [ ] **Content Review:** Audit top 10K suggestions for appropriateness
- [ ] **CDN Configuration:** Edge caching rules, cache TTL, purge mechanism
- [ ] **Trie Build Pipeline:** Verify end-to-end pipeline runs successfully
- [ ] **Failover Testing:** Kill autocomplete servers, verify graceful degradation
- [ ] **Client Integration:** Debouncing, caching, abort controller working correctly
- [ ] **A/B Testing Framework:** Ability to test ranking changes

### 13.2 Day-1 Operations

- [ ] Monitor autocomplete latency (p50, p95, p99)
- [ ] Monitor CDN cache hit rate (target >80%)
- [ ] Verify trie build pipeline completed successfully
- [ ] Check suggestion quality (manual spot-check)
- [ ] Monitor zero-result rate

### 13.3 Week-1 Optimization

- [ ] Analyze suggestion CTR and adjust ranking weights
- [ ] Tune debounce interval based on user behavior data
- [ ] Optimize CDN cache TTL based on hit/miss ratios
- [ ] Identify and cache additional popular prefixes
- [ ] Review and expand content blocklist

### 13.4 Month-1 Scaling

- [ ] Review capacity (plan for 3x growth)
- [ ] Implement personalization if not in MVP
- [ ] Add multi-language support based on user distribution
- [ ] Optimize trie memory usage (radix tree compression)
- [ ] Implement trending queries if not in MVP

---

## Summary: Key Takeaways

### Technical Decisions

| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Data Structure** | Trie with pre-computed top-K | O(P) lookup, sub-ms latency |
| **Serving** | In-memory trie on each server | No network calls, maximum speed |
| **Batch Processing** | Spark/MapReduce | Handle billions of query logs |
| **Real-Time** | Flink streaming | Trending query detection in 15 min |
| **Cache** | CDN + client-side | 80%+ cache hit, reduce server load |
| **Storage** | HDFS/S3 for logs, Redis for trending | Write-heavy logs, fast trending lookup |
| **Ranking** | Time-decayed frequency + CTR | Balance popularity, freshness, quality |

### Scalability Path

1. **Current (10K QPS):** 3 servers with full trie, CDN caching
2. **10x (100K QPS):** 20 servers, aggressive CDN, regional deployment
3. **100x (1M QPS):** Trie sharding by prefix, multi-region, personalization

### Interview Performance Tips

1. ✅ Start with data flow: how do queries become suggestions (offline pipeline)
2. ✅ Explain trie with pre-computed top-K (the key optimization)
3. ✅ Discuss ranking: frequency alone is not enough (need freshness, CTR)
4. ✅ Address client-side: debouncing, caching, abort (often overlooked)
5. ✅ Deep dive into trie update strategy (batch vs streaming vs hybrid)
6. ✅ Discuss content safety (filtering inappropriate suggestions)
7. ✅ Mention CDN caching as a primary scaling lever

---

**End of Search Autocomplete System Design**  
[← Back to Main Index](../README.md)
