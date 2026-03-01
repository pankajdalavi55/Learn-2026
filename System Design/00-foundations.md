# Phase 0: Foundations (Mandatory Prerequisites)

**Navigation:** [Next: Core Building Blocks →](01-core-building-blocks.md)

---

These are the fundamental concepts you MUST internalize before any system design interview. They form the vocabulary and mental models that interviewers expect you to apply naturally. If you struggle with any of these, it's a strong signal that you may not be ready for senior-level roles.

---

## 1. Latency vs Throughput

### Concept Overview (What & Why)

**Latency** is the time it takes to complete a single operation (measured in ms or μs).
**Throughput** is the number of operations completed per unit of time (measured in requests/second, transactions/second, etc.).

**Analogy:** Think of a highway. Latency is how long it takes one car to travel from point A to B. Throughput is how many cars pass a checkpoint per hour. You can have a fast highway (low latency) with few cars (low throughput), or a slow highway (high latency) packed with cars (high throughput).

**Why this matters:** Almost every system design involves trade-offs between latency and throughput. Optimizing for one often impacts the other.

**When interviewers expect this:**
- When discussing caching strategies
- When designing real-time vs batch processing systems
- When choosing between synchronous and asynchronous communication
- When sizing infrastructure

### Key Design Principles

| Principle | Explanation |
|-----------|-------------|
| **Batching increases throughput but adds latency** | Grouping requests together is more efficient but each individual request waits longer |
| **Parallelism can improve both** | But only up to a point (Amdahl's Law) |
| **Queuing trades latency for throughput** | Absorbs burst traffic but adds wait time |
| **Caching improves both** | Reduces latency per request AND reduces backend load |

**Rule of Thumb:**
- User-facing APIs: Prioritize latency (< 100ms)
- Background jobs: Prioritize throughput
- Analytics pipelines: Throughput matters more than latency

### Trade-offs & Decision Matrix

| Choice | Latency Impact | Throughput Impact | When to Choose |
|--------|---------------|-------------------|----------------|
| Sync processing | Low latency | Limited throughput | User-facing, small payloads |
| Async processing | Higher latency | Higher throughput | Background jobs, large payloads |
| Batching | Higher latency | Higher throughput | Data pipelines, bulk operations |
| Connection pooling | Lower latency | Higher throughput | Database connections |
| Caching | Lower latency | Higher throughput | Read-heavy workloads |

### Real-World Examples

- **Payment processing:** Visa optimizes for both - p99 latency < 100ms AND 65,000 transactions/second
- **Search engines:** Google prioritizes latency (< 200ms for user-perceived results)
- **Video encoding:** Netflix prioritizes throughput (millions of hours encoded per day, latency can be minutes)
- **Gaming:** Fortnite requires ultra-low latency (< 50ms network round trip)

### Failure Scenarios & Edge Cases

| Failure | Symptom | Mitigation |
|---------|---------|------------|
| Throughput bottleneck | Requests queue up, latency spikes | Scale horizontally, add caching |
| Latency degradation | Users abandon requests | Set timeouts, add monitoring |
| Thundering herd | Both latency and throughput collapse | Rate limiting, request coalescing |
| GC pauses | Periodic latency spikes | Tune GC, use off-heap storage |

### Interview Perspective

**What interviewers look for:**
- Can you articulate the difference clearly?
- Do you consider both when making design decisions?
- Can you identify which matters more for the given problem?

**Common traps:**
- ❌ "We'll just add more servers" (doesn't address latency)
- ❌ "Caching will solve everything" (cache misses still hit backend)
- ❌ Ignoring tail latency (p99 matters more than average)

**Strong signals:**
- ✅ Mentioning specific latency targets (p50, p99)
- ✅ Discussing batching trade-offs
- ✅ Knowing when to prioritize which

**Follow-up questions:**
- "What's an acceptable latency for this system?"
- "How would you handle 10x traffic?"
- "What happens to latency when throughput doubles?"

### One-Page Cheat Sheet

```
LATENCY vs THROUGHPUT

Latency = Time for ONE operation
Throughput = Operations per unit time

Key Rules:
• Batching: ↑ Throughput, ↑ Latency
• Caching: ↓ Latency, ↑ Throughput (for reads)
• Parallelism: Can improve both (up to a point)
• Queuing: ↑ Throughput, ↑ Latency

Targets:
• User-facing APIs: < 100ms (p99)
• Background jobs: Throughput > Latency
• Real-time systems: < 50ms

Always specify:
• Latency: p50, p95, p99 (not average!)
• Throughput: requests/second, with burst capacity

Red Flag: Ignoring tail latency (p99/p999)
```

---

## 2. Availability vs Consistency (CAP Theorem - Real-World Interpretation)

### Concept Overview (What & Why)

**Availability:** Every request receives a response (success or failure), even if some nodes are down.
**Consistency:** Every read returns the most recent write.

The **CAP Theorem** states that in a distributed system experiencing a network partition, you must choose between Consistency and Availability. You cannot have both.

**The real-world interpretation:** Network partitions WILL happen. The question is: when a partition occurs, do you refuse to respond (choosing consistency) or respond with potentially stale data (choosing availability)?

**Why this matters:** This is the most fundamental trade-off in distributed systems. Every database, cache, and message queue makes this trade-off.

**When interviewers expect this:**
- When choosing databases (PostgreSQL vs Cassandra)
- When designing replication strategies
- When handling network failures
- When discussing data consistency requirements

### Key Design Principles

**CAP Choices in Practice:**

| Category | Behavior During Partition | Examples |
|----------|--------------------------|----------|
| **CP (Consistency + Partition Tolerance)** | Refuse requests if uncertain about data state | ZooKeeper, etcd, HBase, MongoDB (default) |
| **AP (Availability + Partition Tolerance)** | Serve requests with potentially stale data | Cassandra, DynamoDB, CouchDB |

**Important Nuances:**
- CAP is about behavior during partitions, not normal operation
- Most systems are "CP most of the time" but can be configured
- PACELC extends CAP: Even without partition, there's a Latency vs Consistency trade-off

### Trade-offs & Decision Matrix

| Requirement | Choose | Why |
|-------------|--------|-----|
| Financial transactions | CP | Wrong balance = lawsuit |
| User session data | AP | Stale session > no session |
| Inventory counts | CP (with caveats) | Overselling = operational nightmare |
| Social media feeds | AP | Stale post > no feed |
| Leader election | CP | Two leaders = disaster |
| Metrics/Analytics | AP | Approximate data is fine |

### Real-World Examples

**CP Systems:**
- **Banking core systems:** Must refuse transaction if unsure about balance
- **Kubernetes etcd:** Leader election cannot have split-brain
- **Booking systems:** Cannot sell the same seat twice

**AP Systems:**
- **Netflix:** Shows continue even if recommendations are stale
- **Twitter timeline:** You might see tweets slightly out of order
- **DNS:** Old IP is better than no IP

**Interesting Hybrid:**
- **Amazon DynamoDB:** Default is eventually consistent, but offers strongly consistent reads (at 2x cost)
- **Cassandra:** Tunable consistency (ONE, QUORUM, ALL)

### Failure Scenarios & Edge Cases

| Scenario | CP Behavior | AP Behavior |
|----------|-------------|-------------|
| Network partition | Return error | Return stale data |
| Node failure | Wait for recovery | Serve from replicas |
| Split brain | One side stops working | Both sides continue (conflicts later) |
| Replication lag | Wait for sync | Serve inconsistent reads |

**How mature engineers handle this:**
- Design for conflict resolution (last-write-wins, vector clocks, CRDTs)
- Implement idempotency for retries
- Use compensation/saga patterns for eventual consistency
- Monitor replication lag as a key metric

### Interview Perspective

**What interviewers look for:**
- Do you understand CAP applies during partitions?
- Can you explain why you'd choose CP vs AP for specific use cases?
- Do you know that "CA" systems don't exist in distributed environments?

**Common traps:**
- ❌ "We'll use a CA database" (partitions will happen)
- ❌ "We'll just sync everything immediately" (physics prevents this)
- ❌ "Eventually consistent is always bad" (it's often the right choice)

**Weak answer:** "We need strong consistency for everything"
**Strong answer:** "For the payment ledger, we need strong consistency - we'll use a CP database with synchronous replication. For the user's purchase history display, eventual consistency is fine - we can use an AP read replica with a few seconds of lag."

**Follow-up questions:**
- "What happens if the primary database is unreachable?"
- "How would you handle conflicts in an AP system?"
- "What's the replication lag tolerance for this use case?"

### One-Page Cheat Sheet

```
CAP THEOREM (PRACTICAL VERSION)

During Network Partition, choose:
• CP: Refuse request if uncertain (banks, bookings)
• AP: Serve potentially stale data (social media, caching)

Real-World Mapping:
• CP: PostgreSQL, MySQL, MongoDB, ZooKeeper, etcd
• AP: Cassandra, DynamoDB (default), CouchDB

Key Insight: CAP only applies DURING PARTITIONS
Normal operation: You can have both C and A

PACELC Extension:
If Partition → choose A or C
Else → choose Latency or Consistency

Tunable Consistency (Cassandra/DynamoDB):
• R + W > N = Strong Consistency
• R = 1, W = 1 = Eventual Consistency

Red Flags:
• "CA system" (doesn't exist distributed)
• "Strong consistency everywhere" (over-engineering)
• No conflict resolution strategy for AP
```

---

## 3. Scalability vs Elasticity

### Concept Overview (What & Why)

**Scalability:** The ability to handle increased load by adding resources (can be manual).
**Elasticity:** The ability to automatically add or remove resources based on demand.

**Analogy:** A scalable restaurant can add more tables when busy. An elastic restaurant has tables that magically appear during rush hour and disappear during slow hours, so you don't pay for empty tables.

**Why this matters:**
- Scalability is about capability; elasticity is about automation and cost optimization
- Cloud-native design assumes elasticity; on-prem often only achieves scalability
- Cost efficiency at scale requires elasticity

**When interviewers expect this:**
- Cloud architecture discussions
- Cost optimization conversations
- Handling variable/bursty traffic patterns
- Comparing cloud vs on-prem

### Key Design Principles

| Principle | Description |
|-----------|-------------|
| **Design for statelessness** | Stateful services are harder to scale elastically |
| **Externalize state** | Store state in managed services (Redis, S3, databases) |
| **Use auto-scaling groups** | Define min, max, and scaling policies |
| **Design for failure** | Instances will be terminated; design accordingly |
| **Monitor the right metrics** | CPU isn't always the right scaling trigger |

**Scaling Triggers:**
- CPU utilization (common but not always best)
- Request queue depth (better for async workloads)
- Custom metrics (business-specific, like orders/minute)
- Scheduled (known traffic patterns like Black Friday)

### Trade-offs & Decision Matrix

| Approach | Cost | Complexity | Best For |
|----------|------|------------|----------|
| Manual scaling | Higher (over-provision) | Low | Predictable, stable workloads |
| Auto-scaling (reactive) | Medium | Medium | Variable traffic with patterns |
| Predictive scaling | Lower | High | Well-understood traffic patterns |
| Serverless | Pay-per-use | Low (app), High (cold starts) | Sporadic, bursty workloads |

### Real-World Examples

- **Netflix:** Elastic scaling for streaming - traffic peaks in evenings, scales down at night
- **Uber:** Scales backend during surge hours; scales down overnight
- **Black Friday:** Retailers pre-scale before the event (predictive), then elastic during
- **Gaming:** Server auto-scaling based on concurrent players per region

### Failure Scenarios & Edge Cases

| Scenario | Problem | Mitigation |
|----------|---------|------------|
| Scaling too slow | Traffic spikes before new instances ready | Over-provision baseline, use predictive scaling |
| Scaling too aggressively | Thrashing (scale up/down repeatedly) | Add cooldown periods |
| Database doesn't scale | App scales, DB becomes bottleneck | Scale DB separately, use read replicas |
| Stateful services | Can't terminate instances freely | Externalize state, use graceful shutdown |
| Cold start latency | New instances not immediately ready | Keep warm pool, use provisioned concurrency |

### Interview Perspective

**What interviewers look for:**
- Understanding that elasticity ≠ scalability
- Awareness of what can and cannot scale elastically
- Knowledge of cloud auto-scaling mechanisms
- Cost consciousness

**Common traps:**
- ❌ "We'll just use auto-scaling" (what about the database?)
- ❌ "Kubernetes handles everything" (you still need to design stateless services)
- ❌ Ignoring cold start / warm-up time

**Strong signals:**
- ✅ Discussing stateless design for elastic scaling
- ✅ Mentioning cooldown periods and scaling policies
- ✅ Considering what to scale independently (web tier, worker tier, cache)

**Follow-up questions:**
- "How would you handle a 10x traffic spike in 5 minutes?"
- "What if auto-scaling can't keep up?"
- "How do you ensure new instances are ready to serve traffic?"

### One-Page Cheat Sheet

```
SCALABILITY vs ELASTICITY

Scalability = CAN handle more load (add resources)
Elasticity = AUTOMATICALLY handles more load + scales down

Requirements for Elasticity:
1. Stateless services
2. Externalized state (Redis, S3, DB)
3. Health checks and readiness probes
4. Graceful shutdown handling
5. Fast startup time

Auto-Scaling Triggers:
• CPU utilization (simple but often wrong)
• Queue depth (better for async)
• Custom metrics (orders/second, etc.)
• Scheduled (known patterns)

Scaling Policies:
• Step scaling: Add 2 instances if CPU > 70%
• Target tracking: Maintain CPU at 50%
• Predictive: ML-based forecasting

Watch Out For:
• Database scaling (harder than app tier)
• Cold start latency
• Thrashing (scale up/down repeatedly)
• Stateful components blocking elastic scaling

Cost Tip: Elastic = pay for what you use
         Manual = pay for peak capacity always
```

---

## 4. Stateful vs Stateless Services

### Concept Overview (What & Why)

**Stateless Service:** Does not store client session data between requests. Each request is independent.
**Stateful Service:** Maintains client state across requests. The server "remembers" previous interactions.

**Why this matters:**
- Stateless services are dramatically easier to scale, deploy, and recover
- Stateful services require careful design for high availability
- Most modern architectures push state to dedicated stateful components (databases, caches)

**When interviewers expect this:**
- Designing web tiers and APIs
- Discussing load balancing strategies
- Scaling services
- Handling failures and deployments

### Key Design Principles

**Stateless Design Principles:**
1. All state lives in databases, caches, or client tokens (JWT)
2. Any instance can handle any request
3. Instances are disposable and replaceable
4. No sticky sessions required

**When Statefulness is Unavoidable:**
- WebSocket connections (ongoing bidirectional communication)
- In-memory caches for performance (but backed by distributed cache)
- Leader-follower patterns (one leader holds state)
- Batch processing with checkpoints

### Trade-offs & Decision Matrix

| Aspect | Stateless | Stateful |
|--------|-----------|----------|
| Scaling | Easy (add instances) | Hard (need state migration) |
| Load balancing | Any algorithm works | Requires sticky sessions |
| Failure recovery | Instant (route to another instance) | Complex (state recovery needed) |
| Deployment | Rolling deploy easy | Need careful coordination |
| Latency | Slightly higher (fetch state each time) | Lower (state already in memory) |
| Cost | Higher (external state storage) | Lower (if you can maintain it) |

### Real-World Examples

**Stateless:**
- REST APIs (each request contains all needed context)
- Serverless functions (Lambda, Cloud Functions)
- Most microservices (state in external stores)

**Stateful:**
- Database servers (obviously)
- WebSocket servers (maintaining connections)
- Kafka brokers (partition assignments)
- Gaming servers (game state in memory)

**Hybrid Approach (Common):**
- Web tier: Stateless API servers
- Session storage: Redis cluster (stateful, but designed for it)
- Database: PostgreSQL with replicas (stateful, managed carefully)

### Failure Scenarios & Edge Cases

| Scenario | Stateless Impact | Stateful Impact |
|----------|-----------------|-----------------|
| Instance crash | Zero impact, LB routes elsewhere | State lost, recovery needed |
| Deployment | Rolling update, no disruption | Need graceful drain, state handoff |
| Network partition | Requests reroute automatically | Clients may lose session |
| Memory pressure | Simple restart | State corruption risk |

**How to handle stateful components:**
- Use replicated storage (Redis Cluster, database replicas)
- Implement graceful shutdown (drain connections)
- Design for state reconstruction (event sourcing)
- Use leader election (ZooKeeper, etcd)

### Interview Perspective

**What interviewers look for:**
- Understanding why stateless is preferred for horizontal scaling
- Knowing when stateful is acceptable/necessary
- Strategies for managing necessary state

**Common traps:**
- ❌ "We'll store session in local memory" (lost on restart/scaling)
- ❌ "Everything must be stateless" (databases exist for a reason)
- ❌ Not considering WebSocket/long-polling implications

**Strong signals:**
- ✅ "Web tier is stateless, state lives in Redis/DB"
- ✅ "For WebSocket, we'll need sticky sessions or a presence service"
- ✅ "Session state in JWT or distributed cache"

**Follow-up questions:**
- "How do you handle user sessions?"
- "What happens if an instance dies mid-request?"
- "How would you deploy this service with zero downtime?"

### One-Page Cheat Sheet

```
STATELESS vs STATEFUL

Stateless: Server doesn't remember you between requests
Stateful: Server maintains your session

MAKE IT STATELESS:
• Store session in JWT (client-side)
• Store session in Redis (distributed)
• Store data in database
• Pass context with each request

WHEN STATEFUL IS OK:
• WebSocket connections (use presence service)
• In-memory cache (backed by Redis)
• Databases (designed for statefulness)
• Leader-based coordination

SCALING IMPACT:
Stateless → Add instances freely
Stateful → Need sticky sessions, state migration

FAILURE IMPACT:
Stateless → No impact, route to another
Stateful → State lost, recovery needed

GOLDEN RULE:
Web/API tier = Stateless
Data tier = Stateful (use managed services)
Cache tier = Stateful but replicated
```

---

## 5. Horizontal vs Vertical Scaling

### Concept Overview (What & Why)

**Vertical Scaling (Scale Up):** Add more resources (CPU, RAM, disk) to existing servers.
**Horizontal Scaling (Scale Out):** Add more servers to distribute the load.

**Why this matters:**
- Vertical scaling has hard limits (biggest machine available)
- Horizontal scaling is theoretically unlimited but adds complexity
- Most large-scale systems require horizontal scaling
- The architecture must be designed for horizontal scaling; it's not free

**When interviewers expect this:**
- Any discussion of handling more traffic
- Database scaling conversations
- When justifying distributed system complexity

### Key Design Principles

**Vertical Scaling:**
- Simple to implement (just get a bigger machine)
- Single point of failure
- Downtime often required for upgrade
- Eventually hits hardware limits
- Good for: Databases in early stage, simple applications

**Horizontal Scaling:**
- Requires stateless design or distributed state
- No single point of failure (if designed correctly)
- Can scale without downtime
- More complex (load balancing, data distribution)
- Good for: Web servers, API servers, worker pools

### Trade-offs & Decision Matrix

| Factor | Vertical | Horizontal |
|--------|----------|------------|
| Complexity | Low | High |
| Cost curve | Exponential (big machines cost disproportionately more) | Linear |
| Limit | Hardware ceiling | Theoretically unlimited |
| Availability | SPOF | Redundancy built-in |
| Data consistency | Simple | Requires careful design |
| Good for | Databases, legacy apps | Web tier, modern microservices |

### Real-World Examples

**Vertical Scaling:**
- Traditional enterprise databases (Oracle on big iron)
- Legacy applications not designed for distribution
- Small-to-medium startups (simpler than distributed)

**Horizontal Scaling:**
- Google (thousands of commodity servers)
- Netflix (auto-scaling across regions)
- Any modern cloud-native application

**Common Pattern:**
- Web tier: Horizontal (easy to scale stateless servers)
- Database: Vertical first, then horizontal when necessary (read replicas, sharding)
- Cache: Horizontal (Redis Cluster, Memcached)

### Failure Scenarios & Edge Cases

| Scenario | Vertical | Horizontal |
|----------|----------|------------|
| Hardware failure | Complete outage | Partial capacity reduction |
| Need more capacity | Downtime for upgrade | Add nodes (no downtime) |
| Database lock contention | Might help (more CPU) | Doesn't help (still one DB) |
| Network-bound workload | Won't help much | Helps (more network interfaces) |

### Interview Perspective

**What interviewers look for:**
- Understanding that horizontal scaling requires architectural support
- Knowing when vertical is actually the right choice
- Awareness of database scaling challenges

**Common traps:**
- ❌ "We'll just add more servers" (need stateless design first)
- ❌ "Vertical scaling is bad" (it's often simpler and sufficient)
- ❌ Not addressing database scaling separately from app scaling

**Strong signals:**
- ✅ "Web tier scales horizontally; stateless design"
- ✅ "Database: vertical initially, read replicas for reads, sharding when necessary"
- ✅ "Horizontal scaling requires handling distributed system concerns"

**Follow-up questions:**
- "How would you scale the database?"
- "What changes are needed to support horizontal scaling?"
- "At what point would you switch from vertical to horizontal?"

### One-Page Cheat Sheet

```
HORIZONTAL vs VERTICAL SCALING

Vertical (Scale Up):
• Bigger machine
• Simple but limited
• Single point of failure
• Good for: databases early on

Horizontal (Scale Out):
• More machines
• Complex but unlimited
• Redundancy built-in
• Good for: web tier, workers

COST COMPARISON:
Vertical: Exponential (2x CPU ≠ 2x cost)
Horizontal: Linear (2x servers = 2x cost)

WHEN TO USE EACH:
Small/Medium traffic → Vertical
High traffic → Horizontal
Database → Vertical first, then replicas/sharding
Web/API → Horizontal from day one (it's easy)

PREREQUISITES FOR HORIZONTAL:
• Stateless services
• Load balancer
• Distributed state management
• Handle partial failures

RULE OF THUMB:
If one big server can handle it → vertical (simpler)
If you need more than biggest available → horizontal
If you need high availability → horizontal (redundancy)
```

---

## 6. Read-Heavy vs Write-Heavy Systems

### Concept Overview (What & Why)

**Read-Heavy:** System receives far more read requests than write requests (e.g., 100:1 ratio).
**Write-Heavy:** System receives a high volume of writes, possibly more than reads.

**Why this matters:**
- The optimization strategies are completely different
- Read-heavy: Caching, read replicas, CDNs are your friends
- Write-heavy: Write-optimized storage, queuing, async processing
- Most web applications are read-heavy; analytics/IoT are often write-heavy

**When interviewers expect this:**
- First question: "What's the read/write ratio?"
- When designing data layer
- When choosing databases and caching strategies

### Key Design Principles

**Read-Heavy Optimizations:**
| Strategy | How It Helps |
|----------|--------------|
| Caching (Redis, Memcached) | Serve reads from memory |
| Read replicas | Distribute read load |
| CDN | Cache static content at edge |
| Denormalization | Trade storage for read speed |
| Materialized views | Pre-compute expensive queries |

**Write-Heavy Optimizations:**
| Strategy | How It Helps |
|----------|--------------|
| Write-behind caching | Batch writes to database |
| Message queues | Buffer and async process |
| LSM-tree storage (Cassandra) | Optimized for sequential writes |
| Sharding | Distribute write load |
| Event sourcing | Append-only, fast writes |

### Trade-offs & Decision Matrix

| System Type | Database Choice | Caching Strategy | Architecture |
|-------------|-----------------|------------------|--------------|
| Read-heavy (blog) | PostgreSQL + replicas | Aggressive caching | CDN + cache + read replicas |
| Write-heavy (IoT) | Cassandra, TimescaleDB | Minimal caching | Queue + batch write |
| Mixed (social) | Both (polyglot) | Write-through for hot data | Hybrid approach |

### Real-World Examples

**Read-Heavy:**
- Wikipedia (millions of reads, few edits)
- Netflix catalog (streaming metadata)
- E-commerce product pages

**Write-Heavy:**
- IoT sensor data (constant ingestion)
- Logging systems (Elasticsearch, Splunk)
- Metrics collection (time-series data)
- Stock trading (order ingestion)

**Mixed but Complex:**
- Social media (reads for feed, writes for posts/likes)
- Gaming (frequent state updates + reads)

### Failure Scenarios & Edge Cases

| Scenario | Read-Heavy Impact | Write-Heavy Impact |
|----------|-------------------|-------------------|
| Cache failure | Thundering herd to DB | Minimal impact |
| Primary DB down | Serve from replica (stale OK) | Writes fail, queue fills |
| Replication lag | Users see stale data | Less relevant |
| Write amplification | N/A | Storage fills fast |

### Interview Perspective

**What interviewers look for:**
- Do you ask about read/write ratio upfront?
- Do you adjust architecture based on the answer?
- Do you understand why strategies differ?

**Common traps:**
- ❌ Designing for read-heavy when it's write-heavy (or vice versa)
- ❌ "We'll cache everything" for write-heavy systems
- ❌ Not considering mixed workloads

**Strong signals:**
- ✅ "First, what's the read/write ratio?"
- ✅ "For 100:1 read/write, we'll use read replicas and aggressive caching"
- ✅ "For write-heavy, we'll use Cassandra and queue incoming writes"

**Follow-up questions:**
- "What if the ratio changes during peak hours?"
- "How do you handle cache invalidation on writes?"
- "What's the acceptable staleness for reads?"

### One-Page Cheat Sheet

```
READ-HEAVY vs WRITE-HEAVY

ALWAYS ASK: "What's the read/write ratio?"

READ-HEAVY (100:1 or more):
• Cache aggressively (Redis)
• Read replicas
• CDN for static content
• Denormalize data
• Database: PostgreSQL, MySQL

WRITE-HEAVY (1:1 or writes > reads):
• Queue writes (Kafka, SQS)
• Batch processing
• LSM-tree storage (Cassandra, RocksDB)
• Append-only logs
• Minimal caching

MIXED WORKLOADS:
• Separate read/write paths (CQRS)
• Different storage for each
• Write to queue → process → update read store

COMMON RATIOS:
• Blog/Wiki: 1000:1 (reads dominate)
• Social feed: 100:1 (read-heavy)
• Chat: 10:1 (more balanced)
• IoT: 1:100 (write-heavy)
• Gaming: 1:10 to 10:1 (varies)

Database Choice:
Read-heavy → SQL + replicas + cache
Write-heavy → NoSQL (Cassandra, DynamoDB)
```

---

## 7. Strong vs Eventual Consistency

### Concept Overview (What & Why)

**Strong Consistency:** After a write, all subsequent reads return that write. The system behaves like a single copy of data.
**Eventual Consistency:** After a write, reads may return stale data for some time. Eventually, all reads return the write.

**Why this matters:**
- Strong consistency is expensive (latency, availability, throughput)
- Eventual consistency enables higher performance and availability
- The right choice depends on business requirements, not technical preference

**When interviewers expect this:**
- Database selection
- Replication strategy design
- Any distributed data discussion

### Key Design Principles

**Strong Consistency Guarantees:**
- Linearizability: Operations appear to happen instantaneously
- Serializability: Transactions appear to execute in some serial order
- Implementation: Synchronous replication, consensus protocols

**Eventual Consistency Variants:**
| Variant | Guarantee |
|---------|-----------|
| Eventual | Data converges "eventually" (undefined time) |
| Bounded staleness | Data is at most N seconds old |
| Session consistency | Your writes are visible to your reads |
| Read-your-writes | You see your own writes immediately |
| Monotonic reads | You never see older data after seeing newer |

### Trade-offs & Decision Matrix

| Factor | Strong Consistency | Eventual Consistency |
|--------|-------------------|---------------------|
| Latency | Higher (wait for sync) | Lower |
| Availability | Lower (fails if partition) | Higher |
| Throughput | Lower | Higher |
| Complexity | Lower (simpler reasoning) | Higher (conflict resolution) |
| Cost | Higher (sync replicas) | Lower |

**Decision Guide:**
| Requirement | Choose |
|-------------|--------|
| Financial transactions | Strong |
| Inventory (can't oversell) | Strong |
| User profile display | Eventual |
| Social feed | Eventual |
| Shopping cart | Session consistency |
| Analytics | Eventual |
| Leader election | Strong |

### Real-World Examples

**Strong Consistency:**
- Bank account balance
- Ticket booking (can't double-book)
- Inventory count for limited items

**Eventual Consistency:**
- Facebook news feed (it's OK if post appears slightly delayed)
- Netflix recommendations
- Product reviews and ratings
- DNS propagation

**Read-Your-Writes (Common Compromise):**
- After posting a tweet, you see it immediately
- Others might see it a few seconds later
- Implemented by routing your reads to the same replica you wrote to

### Failure Scenarios & Edge Cases

| Scenario | Strong Consistency | Eventual Consistency |
|----------|-------------------|---------------------|
| Network partition | System becomes unavailable | System continues, data diverges |
| High latency | Writes slow down | Writes remain fast |
| Conflict | Prevented by design | Must resolve (LWW, vector clocks) |
| Replication lag | Blocks until caught up | Returns stale data |

**Conflict Resolution Strategies (for Eventual):**
- Last-Write-Wins (LWW): Simple but can lose data
- Vector clocks: Track causality, surface conflicts
- CRDTs: Data structures that merge automatically
- Application-level merge: Let business logic decide

### Interview Perspective

**What interviewers look for:**
- Understanding that consistency is a spectrum
- Matching consistency level to business requirements
- Awareness of implementation costs

**Common traps:**
- ❌ "We need strong consistency everywhere" (over-engineering)
- ❌ "Eventual consistency is always fine" (ignoring business requirements)
- ❌ Not knowing how to handle conflicts in EC systems

**Strong signals:**
- ✅ "For payments, strong consistency is non-negotiable"
- ✅ "For the activity feed, eventual consistency with a 5-second bound is acceptable"
- ✅ "We'll implement read-your-writes for user experience"

**Follow-up questions:**
- "How stale can the data be?"
- "What happens during a conflict?"
- "How would you implement read-your-writes?"

### One-Page Cheat Sheet

```
STRONG vs EVENTUAL CONSISTENCY

Strong: All reads see latest write (expensive)
Eventual: Reads may lag behind writes (cheaper)

CONSISTENCY SPECTRUM:
Linearizable (strongest) 
  ↓ Sequential
  ↓ Causal
  ↓ Read-your-writes
  ↓ Monotonic reads
  ↓ Eventual (weakest)

WHEN TO USE STRONG:
• Money/payments
• Inventory that can't oversell
• Booking/reservations
• Access control changes

WHEN TO USE EVENTUAL:
• Social feeds
• Recommendations
• Analytics
• Non-critical notifications

IMPLEMENTATION COST:
Strong: Synchronous replication, consensus
Eventual: Async replication, conflict resolution

CONFLICT RESOLUTION (for Eventual):
• Last-Write-Wins (simple, lossy)
• Vector clocks (track causality)
• CRDTs (auto-merge)
• App-level merge

USEFUL COMPROMISE:
"Read-your-writes" - You see your writes immediately
Others see them eventually
```

---

## 8. Idempotency

### Concept Overview (What & Why)

**Idempotency:** An operation is idempotent if performing it multiple times has the same effect as performing it once.

**Why this matters:**
- Networks fail; retries are inevitable
- Without idempotency, retries can cause duplicate actions (double charges, double posts)
- Critical for payment systems, order processing, any "exactly-once" semantics
- Fundamental to reliable distributed systems

**When interviewers expect this:**
- Payment/order processing
- Any API design
- Message queue consumers
- Database operations

### Key Design Principles

**Naturally Idempotent Operations:**
- GET, HEAD, OPTIONS (HTTP)
- DELETE (deleting an already-deleted resource is OK)
- SET (overwrite, not increment)
- PUT (replace entire resource)

**Non-Idempotent Operations (need special handling):**
- POST (creates new resource each time)
- INCREMENT (each call adds)
- APPEND (each call adds data)
- Payments (each call charges)

**Making Operations Idempotent:**

1. **Idempotency Key:** Client sends unique ID; server deduplicates
```
POST /payments
Idempotency-Key: abc-123
{amount: 100}

Server checks: Have I processed abc-123?
  Yes → Return cached result
  No → Process and store result with key
```

2. **Conditional Updates:** Use version numbers or timestamps
```sql
UPDATE accounts SET balance = 50 WHERE id = 1 AND version = 5
```

3. **Unique Constraints:** Database prevents duplicates
```sql
INSERT INTO orders (order_id, ...) VALUES (uuid, ...)
-- Second insert with same uuid fails
```

### Trade-offs & Decision Matrix

| Approach | Complexity | Storage Cost | Reliability |
|----------|------------|--------------|-------------|
| Client-generated idempotency key | Low | Medium (store keys) | High |
| Server-generated dedup | Medium | Medium | High |
| Optimistic locking (versions) | Medium | Low | Medium |
| Database unique constraints | Low | Low | High |
| Exactly-once queue processing | High | High | Highest |

### Real-World Examples

**Payment Processing (Stripe):**
- Client sends `Idempotency-Key` header
- Stripe stores the key for 24 hours
- Retries with same key return same result
- Different key = new payment

**Order Processing:**
- Order ID generated client-side (or from cart)
- Server checks if order ID already processed
- Retry creates same order, not duplicate

**Message Queues:**
- SQS: At-least-once delivery, consumer must be idempotent
- Kafka: Can achieve exactly-once with transactions
- Consumer tracks processed message IDs

### Failure Scenarios & Edge Cases

| Scenario | Without Idempotency | With Idempotency |
|----------|---------------------|------------------|
| Client timeout, request succeeded | Client retries, action happens twice | Retry returns cached result |
| Network error mid-response | Client retries, uncertain state | Safe to retry |
| Server crash after processing | Duplicate if not persisted | Key stored, detected as duplicate |
| Queue redelivery | Duplicate processing | Consumer checks processed set |

**Edge Cases to Handle:**
- Key collision (if using random keys, use UUID)
- Key expiration (how long to keep keys?)
- Partial failures (operation half-completed)
- Race conditions (two requests with same key simultaneously)

### Interview Perspective

**What interviewers look for:**
- Automatic mention of idempotency for payment/order systems
- Understanding of idempotency key pattern
- Awareness that retry safety requires design

**Common traps:**
- ❌ Assuming the network is reliable
- ❌ "We'll just prevent retries" (can't control client behavior)
- ❌ Not considering partial failure states

**Strong signals:**
- ✅ "For payments, we'll require an idempotency key"
- ✅ "The POST will be idempotent using client-generated order ID"
- ✅ "Queue consumers must be idempotent; we'll track processed message IDs"

**Follow-up questions:**
- "What if the client sends the same payment twice with different keys?"
- "How long do you store idempotency keys?"
- "What if the server crashes between processing and storing the key?"

### One-Page Cheat Sheet

```
IDEMPOTENCY

Definition: Multiple identical requests = same effect as one

NATURALLY IDEMPOTENT:
• GET, PUT, DELETE
• SET operations (x = 5)
• Overwrites

NOT IDEMPOTENT (need handling):
• POST (creates new resource)
• INCREMENT (x = x + 1)
• Payments, order creation

IMPLEMENTATION PATTERNS:

1. Idempotency Key (Best for APIs):
   Client: POST /pay {key: "abc", amount: 100}
   Server: Check if "abc" processed, return cached or process

2. Unique Constraints (Database):
   INSERT INTO orders (id, ...) 
   -- Same ID = constraint violation = duplicate detected

3. Conditional Updates:
   UPDATE ... WHERE version = expected_version

KEY STORAGE:
• Duration: 24 hours typical (Stripe)
• Storage: Redis (fast), Database (durable)
• Response: Store result, return on retry

CRITICAL FOR:
• Payments (double charge = angry customer)
• Order creation (double order = nightmare)
• Message queue consumers (redelivery happens)
• Any non-GET API in production

EDGE CASES:
• Key expiration
• Race conditions (same key, same time)
• Partial failures (processed but key not stored)
```

---

## Phase 0 Summary: The Foundation Mental Model

These eight concepts form the vocabulary of system design interviews. Before you walk into any interview, you should be able to:

1. **Latency vs Throughput:** Know which to optimize for and why
2. **CAP Theorem:** Explain CP vs AP trade-offs with real examples
3. **Scalability vs Elasticity:** Design for automatic scaling
4. **Stateful vs Stateless:** Default to stateless, know when state is unavoidable
5. **Horizontal vs Vertical:** Know the limits and prerequisites of each
6. **Read vs Write Heavy:** Adjust architecture based on access patterns
7. **Consistency Models:** Match consistency to business requirements
8. **Idempotency:** Ensure safe retries for any mutation

**Interviewer's Perspective:**

If a candidate struggles with these foundations, they're not ready for a senior role. These concepts should be second nature—you should reach for them automatically when analyzing a system.

**Red Flags:**
- Confusing latency with throughput
- Saying "CA system" (doesn't exist)
- Ignoring statefulness when designing for scale
- Not asking about read/write ratio
- Not mentioning idempotency for payment systems

**Green Flags:**
- Naturally using these terms correctly
- Immediately asking clarifying questions about these aspects
- Making trade-off decisions based on these fundamentals
- Discussing failure modes with maturity

---

## Common Interview Questions & Model Answers

This section provides realistic interview questions based on the foundations covered above, with ideal answers and follow-up questions you should expect.

---

### Q1: Explain the difference between latency and throughput. How would you optimize for each?

**Ideal Answer:**

"Latency is the time to complete a single operation—for example, how long it takes to process one API request. Throughput is the number of operations completed per unit time—like requests per second.

To optimize for **low latency**, I would:
- Add caching layers (Redis, CDN)
- Use connection pooling to avoid handshake overhead
- Place servers geographically close to users
- Minimize database queries per request
- Use faster data structures (hash maps over arrays for lookups)

To optimize for **high throughput**, I would:
- Implement batching (group multiple operations)
- Scale horizontally (add more servers)
- Use asynchronous processing (message queues)
- Optimize database with indexes and read replicas
- Employ load balancing

The key is that these can conflict—batching increases throughput but adds latency since each request waits for the batch to fill."

**Follow-up Q:** "Your API latency suddenly jumped from 50ms to 500ms. How would you debug this?"

**Ideal Answer:**

"I'd follow a systematic approach:

1. **Check metrics first:** Look at p50, p99, p99.9 latency—if only p99 spiked, it's likely a specific edge case
2. **Look for recent changes:** Deployments, config changes, traffic spikes
3. **Check external dependencies:** Database query time, third-party API calls, cache hit rate
4. **Review resource utilization:** CPU, memory, network—could be GC pauses or disk I/O
5. **Check logs:** Look for error patterns, slow queries, timeout warnings

Common culprits:
- Cache eviction causing backend overload
- Database connection pool exhaustion
- Network partition affecting a replica
- Memory leak causing GC thrashing"

---

### Q2: What is the CAP theorem? Can you give me an example of a CP and an AP system?

**Ideal Answer:**

"CAP theorem states that in a distributed system during a network partition, you must choose between **Consistency** (all nodes see the same data) and **Availability** (all requests get responses).

**Note:** You can't have a CA system in a distributed environment because network partitions are inevitable. CA only exists in single-node systems.

**CP System (Consistency + Partition Tolerance):**
- **Example:** Banking transaction systems, HBase, MongoDB (with majority writes)
- **Behavior:** During a partition, the system rejects writes to minority partitions to ensure consistency
- **Use case:** When correctness is critical—you can't have two users seeing different account balances

**AP System (Availability + Partition Tolerance):**
- **Example:** DNS, Cassandra, DynamoDB
- **Behavior:** During a partition, all nodes accept writes, leading to temporary inconsistency
- **Use case:** When availability matters more—like social media likes/views where eventual consistency is acceptable

Most real systems aren't purely CP or AP—they offer tuneable consistency. For example, Cassandra lets you configure consistency level per query."

**Follow-up Q:** "If I'm designing a shopping cart, should I choose CP or AP?"

**Ideal Answer:**

"I'd lean toward **AP (availability)** for these reasons:

1. **User experience:** A user should always be able to add items to their cart, even during network issues
2. **Business priority:** Losing a potential sale (unavailable cart) is worse than a temporary inconsistency
3. **Conflict resolution:** If two datacenters both accept cart updates, we can merge them (combine items from both)
4. **Eventual consistency is acceptable:** A few seconds of delay syncing the cart across regions won't hurt

However, at **checkout time**, I'd switch to **CP behavior:**
- Inventory checks must be consistent (can't oversell)
- Payment processing must be strongly consistent
- Use distributed transactions or saga pattern

So the answer is: AP for cart operations, CP for checkout. This is a common pattern—tuneable consistency based on operation criticality."

---

### Q3: How does horizontal scaling differ from vertical scaling? When would you choose each?

**Ideal Answer:**

"**Vertical scaling** means adding more resources to a single machine (CPU, RAM, disk). **Horizontal scaling** means adding more machines.

| Aspect | Vertical Scaling | Horizontal Scaling |
|--------|-----------------|-------------------|
| Complexity | Simple (no code changes) | Complex (requires distributed design) |
| Limits | Hardware ceiling (~1TB RAM) | Nearly unlimited |
| Cost | Exponential (high-end hardware is expensive) | Linear |
| Availability | Single point of failure | High availability |
| Consistency | Easy (single DB) | Requires distributed consensus |

**When to choose Vertical:**
- Early stage, validating product-market fit
- Databases that need ACID guarantees (PostgreSQL)
- Applications with high memory requirements (in-memory analytics)
- When you need to move fast and don't have distributed systems expertise

**When to choose Horizontal:**
- Scale beyond single machine limits
- Need high availability (multi-AZ, multi-region)
- Stateless services (web servers, API gateways)
- Read-heavy workloads (add read replicas)

**Real-world approach:** Start vertical, plan for horizontal. Most companies use both—vertical for databases (up to a point), horizontal for application servers."

**Follow-up Q:** "What are the prerequisites for horizontal scaling?"

**Ideal Answer:**

"For horizontal scaling to work, you need:

1. **Statelessness:** Application servers shouldn't store session state locally
   - Use Redis/Memcached for session storage
   - Or use sticky sessions (but this limits scaling)

2. **Load balancing:** Distribute traffic across servers
   - Layer 7 (HTTP) or Layer 4 (TCP)
   - Health checks to route around failures

3. **Shared data layer:** All servers must access the same data source
   - Centralized database or distributed data store
   - Consistent caching strategy

4. **Idempotent operations:** Handle duplicate requests gracefully
   - Use idempotency keys for mutations
   - Design for at-least-once delivery

5. **Service discovery:** New instances must be discoverable
   - DNS, service mesh, or orchestrator (Kubernetes)

Without these, adding more servers won't help—you'll have inconsistent state or routing problems."

---

### Q4: Explain the difference between strong consistency and eventual consistency. When would you use each?

**Ideal Answer:**

"**Strong consistency** guarantees that after a write completes, all subsequent reads will see that write. **Eventual consistency** means there's a delay—reads might return stale data temporarily, but all replicas will converge.

**Strong Consistency:**
- **Examples:** Financial transactions, inventory management, user authentication
- **Implementation:** Quorum writes (W + R > N in Cassandra), distributed transactions (2PC)
- **Trade-off:** Higher latency, lower availability during partitions
- **Pattern:** Write to primary, wait for replication, then acknowledge

**Eventual Consistency:**
- **Examples:** Social media posts, product reviews, view counts
- **Implementation:** Asynchronous replication, gossip protocols
- **Trade-off:** Temporary inconsistency, but lower latency and higher availability
- **Pattern:** Write locally, replicate asynchronously

**Decision Matrix:**

| Use Strong Consistency When | Use Eventual Consistency When |
|------------------------------|-------------------------------|
| Money is involved | User-generated content |
| Overselling is unacceptable | Stale data is tolerable |
| Legal compliance requires it | Low latency is critical |
| CAP partition is rare | Global distribution needed |

**Hybrid approach:** Use both in the same system. For example, Amazon uses strong consistency for payment processing but eventual consistency for product reviews."

**Follow-up Q:** "How would you implement read-your-own-writes consistency in an eventually consistent system?"

**Ideal Answer:**

"Read-your-own-writes ensures a user sees their own updates immediately, even if other users might see stale data. Here's how:

1. **Session stickiness:** Route a user's requests to the same replica where they wrote
   - Use consistent hashing based on user ID
   - Limitation: Doesn't help if that replica fails

2. **Timestamp-based reads:** Include the write timestamp in subsequent reads
   - Client sends: 'Give me data at least as recent as timestamp T'
   - Server waits for replication to catch up to T
   - Works across replicas

3. **Write to cache on mutation:** After writing to database, also update cache
   - User's next read hits cache (which has latest data)
   - Cache entry expires after replication lag window (few seconds)

4. **Read from primary:** For critical operations, explicitly read from the primary database
   - Override eventual consistency for specific queries
   - Trade latency for consistency when needed

I'd use approach #3 (write-through cache) for most scenarios—it's simple and effective for typical replication lag of 1-3 seconds."

---

### Q5: What is idempotency and why is it important in distributed systems?

**Ideal Answer:**

"Idempotency means an operation can be applied multiple times with the same result as applying it once. For example, 'SET user_status = active' is idempotent, but 'increment likes by 1' is not.

**Why it's critical in distributed systems:**

1. **Network failures are common:** A request might succeed but the response is lost
   - Client doesn't know if it succeeded, so it retries
   - Without idempotency, you get duplicate charges, double-increments, etc.

2. **At-least-once delivery:** Message queues (Kafka, SQS) guarantee at-least-once
   - Same message might be delivered multiple times
   - Idempotent handlers prevent duplicate processing

3. **Distributed transactions:** If one step fails and you retry, you need idempotency
   - Otherwise, partial retries corrupt state

**How to implement idempotency:**

| Pattern | Example | Use Case |
|---------|---------|----------|
| Idempotency key | UUID with each request | Payment APIs (Stripe uses this) |
| Natural key | order_id + user_id | Prevent duplicate order creation |
| Version numbers | Optimistic locking with version field | Concurrent updates |
| Unique constraint | Database unique index | Prevent duplicate inserts |
| Upsert operations | INSERT ... ON CONFLICT UPDATE | Safe to retry |

**Example:** Payment processing
```
POST /payments
{
  'idempotency_key': 'uuid-12345',
  'amount': 100,
  'user_id': 456
}
```
- First request: Process payment, store key
- Retry: Check key exists, return original response
- Different key: New payment

Without this, network timeouts could double-charge users."

**Follow-up Q:** "How long should you store idempotency keys?"

**Ideal Answer:**

"The retention period depends on the retry window and compliance requirements:

**Practical considerations:**

1. **Retry window:** Store keys at least as long as clients might reasonably retry
   - For API requests: 24-48 hours (covers user retries, client-side queues)
   - For payment systems: 30 days (compliance, dispute resolution)
   - For background jobs: 7 days (covers retry backoff strategies)

2. **Storage cost vs risk:** Longer is safer but costs more
   - Use TTL in Redis: Automatic expiration
   - Archive to cheaper storage after active period

3. **Compliance requirements:** Some industries mandate retention
   - Financial: Often 7 years for audit trails
   - Healthcare: HIPAA compliance periods

**Implementation strategy:**
- **Hot storage (Redis):** 24 hours for fast lookups
- **Warm storage (Database):** 30 days for recent history
- **Cold storage (S3):** Years for compliance/auditing

**Trade-off:** Stripe stores payment idempotency keys for 24 hours in active memory, then archives for auditing. This balances performance with safety."

---

### Q6: Explain database sharding. What are the challenges?

**Ideal Answer:**

"Sharding is horizontal partitioning—splitting a database into smaller, independent pieces (shards) that can live on different servers. Each shard contains a subset of the data.

**Why shard:**
- Single database can't handle the load (too many queries or too much data)
- Distribute across multiple servers for scale
- Each shard is smaller, so queries are faster

**Sharding strategies:**

1. **Range-based:** Shard by value range (users A-M on shard1, N-Z on shard2)
   - Pro: Easy to add new data
   - Con: Unbalanced (more users with last name starting with S)

2. **Hash-based:** Hash the key, modulo by number of shards
   - Pro: Even distribution
   - Con: Hard to add shards (requires rehashing)

3. **Geographic:** Shard by region (US users → US shard)
   - Pro: Data locality, comply with regulations (GDPR)
   - Con: Uneven load

4. **Directory-based:** Lookup table maps keys to shards
   - Pro: Flexible, can rebalance
   - Con: Lookup table is a bottleneck

**Major challenges:**

1. **Cross-shard queries:** Joining data across shards is slow
   - Solution: Denormalization, or accept expensive scatter-gather

2. **Shard key choice is permanent:** Hard to change later
   - Solution: Choose carefully upfront, plan for consistent hashing

3. **Uneven data distribution (hot spots):**
   - Solution: Composite shard keys, split hot shards

4. **No ACID across shards:** Can't use database transactions
   - Solution: Saga pattern or two-phase commit (slow)

5. **Operational complexity:** Multiple databases to manage
   - Solution: Automation, managed services (AWS Aurora, Vitess)"

**Follow-up Q:** "How would you reshard a live production database?"

**Ideal Answer:**

"Resharding live is one of the hardest operations. Here's a phased approach:

**Phase 1: Preparation (weeks)**
1. **Add shard_id column:** To all tables, populated by application
2. **Dual writes:** Write to current shard AND future shard
3. **Verify dual writes:** Monitor for discrepancies

**Phase 2: Data Migration (days)**
1. **Background job:** Copy historical data to new shards
2. **Use consistent hashing:** Minimize data movement
3. **Validate:** Compare checksums between old and new shards
4. **Catch-up replication:** Sync any changes during migration

**Phase 3: Cutover (hours)**
1. **Stop writes:** Brief maintenance window (or use feature flags)
2. **Final sync:** Ensure all data is copied
3. **Switch routing:** Update application to read from new shards
4. **Monitor:** Watch for errors, ready to rollback
5. **Gradual rollout:** Route 1% traffic, then 10%, then 100%

**Phase 4: Cleanup (days)**
1. **Run both shards in parallel:** For 24-48 hours
2. **Verify new shards work:** No errors, performance acceptable
3. **Drop old shards:** After confirmed success

**Zero-downtime technique:**
- Use a **proxy layer** (Vitess, ProxySQL) to handle routing
- Shift traffic gradually without application changes
- Fallback to old shards if issues arise

**Real example:** Discord resharded from 1 database to 12 using this approach, took 3 months of planning for 2 days of execution."

---

### Q7: What's the difference between stateful and stateless services? How does it impact design?

**Ideal Answer:**

"**Stateless services** don't store any client-specific data between requests. Each request is independent and contains all necessary information.

**Stateful services** maintain client data across requests—like session state, connection state, or in-memory caches.

**Impact on design:**

| Aspect | Stateless | Stateful |
|--------|-----------|----------|
| Scalability | Easy to scale horizontally | Hard to scale (state must be migrated) |
| Load balancing | Any server can handle any request | Needs sticky sessions or state replication |
| Failure handling | Restart without data loss | Requires state recovery/persistence |
| Deployment | Rolling updates are simple | Requires careful state migration |
| Examples | REST APIs, Lambda functions | WebSocket servers, game servers |

**When stateless:**
- Web servers (session in Redis, not in-memory)
- API gateways
- Serverless functions
- Most microservices

**When stateful:**
- WebSocket connections (persistent TCP connection to one server)
- In-memory caches (state is the cache itself)
- Database connections (connection pooling)
- Real-time gaming (player state in server memory)

**Design principle:** Default to stateless, use stateful only when required. When you must be stateful, push state to external storage (Redis, database) rather than in-memory."

**Follow-up Q:** "How would you handle sessions in a stateless architecture?"

**Ideal Answer:**

"In a stateless architecture, session data lives outside the application servers. Here are the approaches:

**1. Client-side sessions (JWT tokens):**
- Store session data in the token itself
- Client sends token with each request
- Server validates signature, extracts data
- **Pros:** Truly stateless, scales infinitely
- **Cons:** Token can't be revoked easily, size limits (4KB in cookies)

**2. External session store (Redis/Memcached):**
- Application generates session ID
- Stores session data in Redis: `session:abc123 -> {user_id, permissions}`
- Client sends session ID with each request
- Application fetches session from Redis
- **Pros:** Can invalidate sessions, no size limits
- **Cons:** Network hop to Redis (but fast: ~1ms)

**3. Database sessions:**
- Similar to Redis, but use database table
- **Pros:** Persistent, survives cache eviction
- **Cons:** Slower, adds DB load

**Best practice (hybrid):**
```
1. Store minimal data in JWT (user_id, role)
2. Store sensitive data in Redis (permissions, preferences)
3. Cache Redis data in application for 1-2 minutes
```

This gives you:
- Fast reads (local cache)
- Security (sensitive data in Redis, not client)
- Revocability (invalidate Redis key)
- Scalability (stateless servers)

**Real example:** Netflix uses this approach—JWT for user identity, Redis for session-specific data like watch history position."

---

### Q8: How would you design a rate limiter?

**Ideal Answer:**

"A rate limiter controls how many requests a user/IP can make in a time window. I'd approach it like this:

**Requirements clarification:**
- What are we limiting? (Per user, per IP, per API key?)
- What's the limit? (100 requests per minute? 1000 per hour?)
- Distributed system or single server?
- Hard limit (reject) or soft limit (throttle)?

**Algorithm choices:**

1. **Token Bucket** (most common)
   - Bucket refills with tokens at fixed rate
   - Each request consumes a token
   - If no tokens, request is rejected
   - **Pros:** Handles bursts, smooth refill
   - **Use:** API rate limiting (AWS, Stripe)

2. **Leaky Bucket**
   - Requests enter a queue, processed at fixed rate
   - Queue overflow → reject
   - **Pros:** Smooth output rate
   - **Use:** Traffic shaping, network packets

3. **Fixed Window**
   - Count requests in fixed time windows (10:00-10:01)
   - Reset counter at window boundary
   - **Pros:** Simple, memory efficient
   - **Cons:** Burst at window boundaries (180 requests in 2 seconds)

4. **Sliding Window Log**
   - Store timestamp of each request
   - Count requests in the last N seconds
   - **Pros:** Accurate, no boundary issues
   - **Cons:** Memory intensive (store all timestamps)

5. **Sliding Window Counter** (best balance)
   - Hybrid of fixed window and sliding window
   - Weighted count from current and previous window
   - **Pros:** Memory efficient, accurate
   - **Use:** My recommendation for most cases

**Implementation (distributed system):**

```
Redis implementation (Token Bucket):

Key: rate_limit:user_123
Value: {tokens: 95, last_refill: 1640000000}

On each request:
1. Calculate tokens to add: (now - last_refill) * refill_rate
2. Add tokens (up to max bucket size)
3. If tokens >= 1: Allow, decrement token
4. Else: Reject with 429 status

Lua script for atomicity (Redis)
```

**Response headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000060
Retry-After: 30 (if rate limited)
```

**Scaling considerations:**
- Use Redis for shared state across servers
- Local in-memory cache with eventual consistency for high throughput
- Different limits for different user tiers (free vs paid)

**Example:** GitHub API uses this—5000 requests/hour for authenticated users."

**Follow-up Q:** "What if the rate limiter itself becomes a bottleneck?"

**Ideal Answer:**

"If the rate limiter is overloaded, it defeats the purpose. Here's how to prevent it:

**1. Local caching with eventual consistency:**
- Cache rate limit state locally in each server
- Sync with Redis every 1-10 seconds
- **Trade-off:** Slight inaccuracy (might allow 105 requests instead of 100)
- **Benefit:** 10-100x less Redis load

**2. Approximate algorithms:**
- Use **HyperLogLog** for counting (probabilistic)
- Use **Bloom filter** to remember recent IPs
- Sacrifice exact accuracy for massive scalability

**3. Multi-layer rate limiting:**
- **Layer 1 (CDN/Edge):** Coarse limits (10,000 req/sec per IP)
- **Layer 2 (Load Balancer):** Medium limits (1,000 req/min per user)
- **Layer 3 (Application):** Fine-grained limits (100 req/min per API key)
- Each layer filters out obvious abuse before it hits the next layer

**4. Partition the problem:**
- Shard rate limit counters by user_id hash
- Multiple Redis instances, each handles a subset
- Consistent hashing for even distribution

**5. Asynchronous updates:**
- Accept request immediately (don't wait for Redis)
- Update rate limit counter in background
- Periodically check if over limit (eventual enforcement)
- **Use case:** Very high throughput where occasional overage is acceptable

**Real example:** Cloudflare uses a multi-layer approach—edge servers do basic IP limiting (no Redis), origin servers do precise limiting. This handles DDoS attacks (millions of req/sec) without overloading Redis."

---

**Navigation:** [\u2190 Previous: Foundations](00-foundations.md) | [Next: Core Building Blocks →](01-core-building-blocks.md)
- Immediately asking clarifying questions about these aspects
- Making trade-off decisions based on these fundamentals
- Discussing failure modes with maturity
