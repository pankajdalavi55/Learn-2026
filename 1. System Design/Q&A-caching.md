1. What is a cache?

---

1. cache Level?

## 🎯 3-Min Interview Answer: Cache Levels

> In distributed systems, caching is implemented at multiple levels to improve performance, reduce latency, and minimize load on backend systems. These levels are typically organized from closest to the user to closest to the database.

The first level is **client-side caching**, which happens in the browser or mobile app. It uses mechanisms like HTTP headers such as Cache-Control and ETag. This is the fastest layer since it avoids network calls entirely, but it’s harder to invalidate and may serve stale data.

The second level is the **CDN or edge cache**, where content is stored on geographically distributed servers. This helps deliver static assets like images, videos, and APIs with low latency by serving data from the nearest location, reducing load on the origin servers.

The third level is the **application-level cache**, typically implemented using in-memory systems like Redis or Memcached. This is the most commonly used layer in backend systems and helps reduce database load by caching frequently accessed data such as user profiles or product details.

The final level is the **database cache**, where the database itself caches query results or data pages, such as buffer pools. This is transparent to the application but provides performance improvements for repeated queries.

In a typical request flow, the system checks cache hierarchically: client → CDN → application cache → database. As we move closer to the database, consistency improves but latency increases.

> Overall, a well-designed system combines multiple cache levels to balance performance, scalability, and consistency based on use case requirements.

---

1. Cache Eviction policies?

## 🎯 3-Min Interview Answer: Cache Eviction Policies (with Problems & Solutions)

> Cache eviction policies determine how a system removes data when the cache becomes full. Since cache memory is limited, the goal is to retain the most useful data and evict less valuable data to maintain a high cache hit ratio and system performance.

The most commonly used policy is **Least Recently Used (LRU)**, where the system evicts data that hasn’t been accessed for the longest time. It works well for typical workloads but can fail in cyclic access patterns where frequently needed data gets evicted. This is usually solved by combining LRU with TTL or using approximations like segmented LRU.

Another policy is **Least Frequently Used (LFU)**, which removes data with the lowest access frequency. While it is effective for stable access patterns, it can lead to cache pollution where old popular data never gets evicted. To solve this, systems use LFU with aging or decay so that old frequency counts reduce over time.

**First In First Out (FIFO)** removes the oldest inserted data regardless of usage. The main issue is that it ignores access patterns, so even frequently used data might get evicted. Because of this, FIFO is rarely used alone in critical systems and is often combined with other strategies.

There is also **Random Eviction**, where entries are removed randomly. This is simple and has low overhead but can evict important data unpredictably. In practice, it is sometimes combined with sampling techniques to approximate better policies like LRU.

Another important approach is **TTL-based eviction**, where data automatically expires after a fixed time. While simple, it can lead to problems like cache avalanche if many keys expire simultaneously or removal of still-useful data. This is typically mitigated by adding TTL randomization or combining it with LRU or LFU.

Modern systems like Redis use hybrid approaches that combine multiple strategies to balance performance and memory efficiency.

> Overall, each eviction policy has trade-offs, so real-world systems use a combination of strategies based on access patterns, ensuring optimal memory usage and consistent performance.

---

1. Explain Cache patterns?

## 🎯 3-Min Interview Answer: Cache Patterns

> Cache patterns define how an application interacts with the cache and database for reading and writing data. Choosing the right pattern is important to balance performance, consistency, and scalability.

The most commonly used pattern is **Cache-Aside (Lazy Loading)**. In this approach, the application first checks the cache, and on a cache miss, it fetches data from the database and updates the cache. It is simple and widely used, but it can lead to issues like cache miss latency and stale data. These are usually handled using TTL, proper invalidation, or techniques like stale-while-revalidate.

Another pattern is **Read-Through Cache**, where the cache itself is responsible for fetching data from the database on a miss. This simplifies application logic, but it reduces control over caching behavior and adds complexity to the cache layer.

For write operations, **Write-Through Cache** ensures that data is written to both the cache and the database simultaneously. This provides strong consistency, but it increases write latency and may store unnecessary data in the cache. It is typically used in systems where consistency is critical, such as financial applications.

In contrast, **Write-Back (Write-Behind)** writes data only to the cache and updates the database asynchronously. This improves write performance and throughput but introduces the risk of data loss if the cache fails before persisting to the database. It is suitable for high-write systems like analytics or logging.

Another important pattern is **Refresh-Ahead**, where the system proactively refreshes cache entries before they expire. This avoids cache misses for frequently accessed data but can increase load on the database if not managed carefully.

> Overall, each cache pattern has trade-offs between consistency, latency, and complexity. Cache-Aside is most commonly used, while write-through and write-back are chosen based on consistency and performance needs. In real-world systems, a combination of these patterns is often used to achieve optimal results.

---

1. cache problems?

Here’s a **concise 3-minute interview answer** covering all major cache problems clearly:

---

## 🎯 3-Min Answer: Cache Problems

> In distributed systems, caching improves performance but introduces several challenges that need careful handling.

The most common problem is **cache stampede**, where a popular key expires and multiple concurrent requests hit the database simultaneously, causing overload. This is typically solved using **mutex locking, request coalescing, or stale-while-revalidate**.

Another issue is **cache penetration**, where requests are made for non-existent data. Since the data is never cached, every request hits the database. This can be mitigated by **caching null values, using Bloom filters, or validating inputs**.

We also have **cache breakdown (hot key problem)**, where a highly popular key expires and leads to a sudden spike in database traffic. Solutions include **never expiring hot keys, using refresh-ahead, or locking mechanisms**.

A related issue is **cache avalanche**, where many keys expire at the same time, causing a massive surge in database requests. This can be prevented using **TTL randomization, rate limiting, and multi-level caching**.

Another critical challenge is **stale data**, where cache returns outdated information due to improper invalidation. This is handled using **write-invalidate, write-through, or event-driven invalidation mechanisms**.

Additionally, **cache pollution** occurs when rarely used data occupies cache space, reducing efficiency. This is addressed using **LRU/LFU eviction policies and selective caching**.

Finally, **cache failures**, such as Redis outages, can push all traffic to the database. To handle this, systems use **circuit breakers, fallback strategies, and replication**.

> Overall, designing a robust caching layer requires balancing performance, consistency, and fault tolerance while planning for these failure scenarios.

---

1. cache stampede and how to avoid?

## 🎯 3-Min Interview Answer: Cache Stampede

> Cache stampede, also known as the thundering herd problem, occurs when a cached item expires and multiple concurrent requests try to fetch the same data from the database at the same time. This leads to a sudden spike in database load and can degrade system performance or even cause failures.

This typically happens in high-traffic systems when a popular or “hot” key expires. Since all requests see a cache miss, they simultaneously hit the database, resulting in excessive load, increased latency, and potential cascading failures.

One common solution is **mutex locking**, where only one request is allowed to fetch data from the database while others wait. This ensures that only a single database call is made and the cache is updated once.

Another approach is **stale-while-revalidate**, where the system serves stale data temporarily while refreshing the cache in the background. This reduces latency and prevents sudden load spikes on the database.

**TTL randomization (jitter)** is also used to avoid multiple keys expiring at the same time, thereby spreading out the load more evenly. Additionally, **request coalescing** can be used to combine multiple identical requests into a single database call.

However, these solutions also come with trade-offs. Locking can increase latency for waiting requests, and stale-while-revalidate may serve slightly outdated data. Therefore, the choice depends on the system’s consistency and performance requirements.

> Overall, cache stampede is a critical issue in high-scale systems, and it is typically mitigated using a combination of locking, background refresh, and smart expiration strategies to protect the database and maintain performance.

---

1. cache Invalidation?

## 🎯 3-Min Interview Answer: Cache Invalidation

> Cache invalidation is the process of ensuring that cached data remains consistent with the source of truth, typically the database. It is one of the hardest problems in distributed systems because it requires balancing performance with data correctness.

The most common approach is **write-invalidate**, where after updating the database, the corresponding cache entry is deleted. This ensures that the next read fetches fresh data from the database. However, it introduces a cache miss on the next request and can increase database load temporarily.

Another approach is **TTL-based invalidation**, where cached data automatically expires after a fixed time. This is simple to implement but can result in stale data being served until the TTL expires. It can also lead to issues like cache avalanche if many keys expire simultaneously, which is usually mitigated using TTL randomization.

**Write-through caching** is another strategy where data is updated in both the cache and database at the same time. This ensures strong consistency but increases write latency and may lead to unnecessary cache updates.

In distributed systems, **event-driven invalidation** is commonly used. When data is updated, an event is published (for example using Apache Kafka), and other services listening to the event invalidate or update their caches. While this improves scalability and decoupling, it introduces complexity and possible delays, leading to temporary inconsistency.

These approaches also have challenges. For example, race conditions can occur if cache updates happen before database commits, or cache deletion might fail, leaving stale data. To address this, systems use retries, idempotent operations, and proper ordering of updates.

> Overall, cache invalidation strategies like write-invalidate, TTL, write-through, and event-driven approaches are chosen based on consistency and scalability needs. In practice, systems combine multiple strategies to maintain performance while ensuring data correctness.


### Which eviction policy would you choose and why?
#### 🎯 3-Min Interview Answer: Which Eviction Policy Would You Choose and Why?

> The choice of eviction policy depends on the system’s access pattern, but in most real-world applications, I would choose **LRU (Least Recently Used)** or a **hybrid approach like LRU + TTL**, because it provides a good balance between performance, simplicity, and effectiveness.

LRU works well because it is based on the principle of **temporal locality**, meaning recently accessed data is more likely to be accessed again. This makes it highly effective for common use cases like user sessions, product pages, or API responses.

However, LRU alone has limitations. For example, in cyclic access patterns, it may evict useful data. To address this, I would combine it with **TTL (Time-To-Live)** to ensure stale data does not stay indefinitely and to provide time-based control over cache entries.

In cases where access patterns are highly skewed, such as trending content or hot keys, I might prefer **LFU (Least Frequently Used)** or **LFU with aging**, because it prioritizes frequently accessed data over time. But LFU can cause cache pollution if not tuned properly, so aging or decay mechanisms are important.

In practice, systems like Redis use **approximations of LRU or LFU along with TTL**, which is often the best approach in production systems.

> So overall, I would choose **LRU with TTL as a default**, and adapt to LFU or hybrid strategies depending on the workload, ensuring a balance between cache efficiency, memory usage, and system performance.


### how Redis implements LRU/LFU internally
#### 🎯 3-Min Interview Answer: How Redis Implements LRU / LFU Internally

> Redis does not implement *true* LRU or LFU because maintaining exact ordering or counters for every key would be expensive in terms of memory and CPU. Instead, it uses **approximated algorithms** that are efficient and scalable.

---

### 🔹 LRU in Redis (Approximate LRU)

Redis uses a **sampling-based LRU** approach:

* Each key stores a small **LRU field (timestamp / idle time info)**
* When eviction is needed, Redis:

  1. Randomly samples a few keys (default ~5)
  2. Chooses the key that was **least recently used among the sample**
  3. Evicts that key

### ⚠️ Problem

* Not perfectly accurate (may not evict the absolute least recently used key)

### ✅ Why it works

* Very low overhead
* Good enough approximation for large-scale systems

👉 Config:

```bash
maxmemory-policy allkeys-lru
```

---

### 🔹 LFU in Redis (Approximate LFU)

Redis uses a clever **probabilistic counter with decay**:

* Each key stores:

  * **Access frequency counter (8-bit)**
  * **Last decay time**

* On every access:

  * Counter increases **probabilistically** (not linearly)

* Over time:

  * Counter **decays** (reduces) to forget old popularity

### ⚠️ Problem

* Not exact frequency count

### ✅ Why it works

* Prevents cache pollution
* Adapts to changing access patterns
* Very memory efficient

👉 Config:

```bash
maxmemory-policy allkeys-lfu
```

---

### 🔹 Why Redis Uses Approximation (Key Insight)

* True LRU → needs doubly linked list updates on every access (high cost)
* True LFU → needs exact counters (high memory overhead)

👉 Redis optimizes for:

* ⚡ Performance
* 📉 Low memory overhead
* 📈 Scalability

---

### 🔹 Additional Optimizations

* **Configurable sample size** → better accuracy vs performance
* **TTL support** → works alongside eviction
* **Lazy eviction** → eviction happens only when needed

---

### 🔹 Final Summary

> Redis implements LRU and LFU using approximate algorithms. LRU is implemented using random sampling to pick the least recently used key among a subset, while LFU uses probabilistic counters with decay to track frequency. These approaches trade perfect accuracy for high performance and low memory usage, making them suitable for large-scale distributed systems.

