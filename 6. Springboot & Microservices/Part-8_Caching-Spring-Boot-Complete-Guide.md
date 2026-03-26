# Caching in Spring Boot — Complete Guide

> A comprehensive guide covering caching fundamentals, Spring Boot cache abstraction, providers, and production best practices.  
> **Part 8 of the Spring Boot & Microservices Series**  
> **Prerequisites:** Spring Boot basics, Java fundamentals

---

## Table of Contents

1. [Introduction to Caching](#1-introduction-to-caching)
2. [Core Caching Concepts](#2-core-caching-concepts)
3. [How Caching Works in Spring Boot](#3-how-caching-works-in-spring-boot)
4. [Caching Providers in Spring Boot](#4-caching-providers-in-spring-boot)
5. [When to Choose Which Cache Provider](#5-when-to-choose-which-cache-provider)
6. [Advanced Topics](#6-advanced-topics)
7. [Conclusion](#7-conclusion)

---

## 1. Introduction to Caching

### 1.1 What is Caching?

**Caching** is a technique for storing frequently accessed data in a high-speed storage layer (cache) to reduce the time and resources needed to fetch that data from its original source (database, API, file system).

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         CACHING CONCEPT                                          │
│                                                                                  │
│     ┌──────────┐         ┌──────────┐         ┌──────────────┐                  │
│     │  Client  │ ──────► │  Cache   │ ──────► │  Data Source │                  │
│     │ Request  │         │  Layer   │         │  (Database)  │                  │
│     └──────────┘         └──────────┘         └──────────────┘                  │
│                                │                                                 │
│                                │                                                 │
│                         ┌──────▼──────┐                                         │
│                         │ Cache Hit?  │                                         │
│                         └──────┬──────┘                                         │
│                                │                                                 │
│                    ┌───────────┴───────────┐                                    │
│                    │                       │                                    │
│               ┌────▼────┐            ┌─────▼─────┐                              │
│               │   YES   │            │    NO     │                              │
│               │ Return  │            │ Fetch from│                              │
│               │ cached  │            │ source &  │                              │
│               │  data   │            │ cache it  │                              │
│               └─────────┘            └───────────┘                              │
│                                                                                  │
│  WITHOUT CACHE:  Client → Database → Client  (100ms - 500ms)                    │
│  WITH CACHE:     Client → Cache → Client     (1ms - 10ms) ✓                     │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Why is Caching Important?

| Benefit | Description | Impact |
|---------|-------------|--------|
| **Performance Improvement** | Data served from memory is 10-100x faster than disk | Response time: 200ms → 5ms |
| **Reduced Latency** | Eliminates network round-trips to database/APIs | Better user experience |
| **Scalability** | Handle more requests with same resources | 10x throughput increase |
| **Reduced Database Load** | Fewer queries hit the database | Database costs down 60-80% |
| **Cost Efficiency** | Less compute and database resources needed | Lower infrastructure costs |
| **High Availability** | Serve stale data when source is unavailable | Improved resilience |

### 1.3 Common Caching Patterns

#### 1.3.1 Cache-Aside (Lazy Loading)

The application is responsible for reading from and writing to the cache.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CACHE-ASIDE PATTERN                               │
│                                                                      │
│   READ FLOW:                                                         │
│   ┌─────────┐      1. Check      ┌─────────┐                        │
│   │   App   │ ─────────────────► │  Cache  │                        │
│   └────┬────┘                    └────┬────┘                        │
│        │                              │                              │
│        │  3. Return data         2. Cache Miss                      │
│        │◄─────────────────────────────┘                              │
│        │                                                             │
│        │  2. Query DB            ┌──────────┐                       │
│        └─────────────────────────►│ Database │                       │
│                                  └──────────┘                       │
│                                                                      │
│   WRITE FLOW:                                                        │
│   App writes to DB first, then invalidates/updates cache            │
│                                                                      │
│   PROS:                          CONS:                               │
│   • Simple to implement          • Cache miss penalty                │
│   • Cache only what's needed     • Potential stale data              │
│   • Resilient to cache failures  • Application complexity            │
└─────────────────────────────────────────────────────────────────────┘
```

```java
// Cache-Aside Pattern Example
public class ProductService {
    
    private final Cache cache;
    private final ProductRepository repository;
    
    public Product getProduct(Long id) {
        // 1. Check cache first
        String cacheKey = "product:" + id;
        Product cached = cache.get(cacheKey, Product.class);
        
        if (cached != null) {
            return cached;  // Cache hit
        }
        
        // 2. Cache miss - fetch from database
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        
        // 3. Store in cache for future requests
        cache.put(cacheKey, product);
        
        return product;
    }
    
    public Product updateProduct(Long id, ProductRequest request) {
        // 1. Update database
        Product updated = repository.save(mapToEntity(request));
        
        // 2. Invalidate cache
        cache.evict("product:" + id);
        
        return updated;
    }
}
```

#### 1.3.2 Write-Through Cache

Data is written to the cache and database simultaneously.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WRITE-THROUGH PATTERN                             │
│                                                                      │
│   ┌─────────┐      1. Write     ┌─────────┐      2. Write           │
│   │   App   │ ─────────────────►│  Cache  │─────────────────►       │
│   └─────────┘                   └─────────┘               │         │
│                                                           │         │
│                                                    ┌──────▼─────┐   │
│                                                    │  Database  │   │
│                                                    └────────────┘   │
│                                                                      │
│   CHARACTERISTICS:                                                   │
│   • Synchronous write to both cache and DB                          │
│   • Data consistency guaranteed                                      │
│   • Higher write latency                                            │
│                                                                      │
│   PROS:                          CONS:                               │
│   • Data always consistent       • Slower writes                     │
│   • Simple read operations       • Cache may have unused data        │
│   • No stale data                • Higher write latency              │
└─────────────────────────────────────────────────────────────────────┘
```

#### 1.3.3 Write-Behind (Write-Back) Cache

Data is written to cache immediately, then asynchronously written to database.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WRITE-BEHIND PATTERN                              │
│                                                                      │
│   ┌─────────┐      1. Write     ┌─────────┐                         │
│   │   App   │ ─────────────────►│  Cache  │                         │
│   └─────────┘      (fast)       └────┬────┘                         │
│        │                             │                               │
│        │  Return immediately         │  2. Async batch write         │
│        ◄─────────────────────────────┤                               │
│                                      │                               │
│                                      │     ┌────────────┐            │
│                                      └────►│  Database  │            │
│                                            └────────────┘            │
│                                                                      │
│   CHARACTERISTICS:                                                   │
│   • Fast write response                                             │
│   • Batch writes improve DB efficiency                              │
│   • Risk of data loss if cache fails before DB write                │
│                                                                      │
│   PROS:                          CONS:                               │
│   • Very fast writes             • Risk of data loss                 │
│   • Reduced DB load              • Complex implementation            │
│   • Batch optimization           • Eventual consistency              │
└─────────────────────────────────────────────────────────────────────┘
```

#### 1.3.4 Read-Through Cache

Cache sits between application and database; cache fetches data automatically on miss.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    READ-THROUGH PATTERN                              │
│                                                                      │
│   ┌─────────┐      1. Request   ┌─────────┐      2. Auto-fetch      │
│   │   App   │ ─────────────────►│  Cache  │─────────────────►       │
│   └─────────┘                   └────┬────┘               │         │
│        │                             │                    │         │
│        │  3. Return data             │             ┌──────▼─────┐   │
│        ◄─────────────────────────────┘             │  Database  │   │
│                                                    └────────────┘   │
│                                                                      │
│   CHARACTERISTICS:                                                   │
│   • Cache handles data fetching transparently                       │
│   • Application only interacts with cache                           │
│   • Simpler application code                                        │
│                                                                      │
│   PROS:                          CONS:                               │
│   • Simpler application code     • Cache library dependency          │
│   • Centralized caching logic    • Less control over fetching        │
│   • Automatic cache population   • First request always slow         │
└─────────────────────────────────────────────────────────────────────┘
```

#### 1.3.5 Pattern Comparison

| Pattern | Write Latency | Read Latency | Consistency | Data Loss Risk | Use Case |
|---------|---------------|--------------|-------------|----------------|----------|
| **Cache-Aside** | Medium | Low (hit) / High (miss) | Eventual | Low | Most applications |
| **Write-Through** | High | Low | Strong | Very Low | Financial systems |
| **Write-Behind** | Very Low | Low | Eventual | Medium | High-write systems |
| **Read-Through** | N/A | Low (hit) / High (miss) | Eventual | Low | Read-heavy systems |

### 1.4 Real-World Caching Examples

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REAL-WORLD CACHING USE CASES                                  │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  E-COMMERCE PLATFORM                                                     │   │
│  │  • Product catalog (rarely changes, frequently accessed)                 │   │
│  │  • User sessions                                                         │   │
│  │  • Shopping cart                                                         │   │
│  │  • Search results                                                        │   │
│  │  • Price calculations                                                    │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  SOCIAL MEDIA                                                            │   │
│  │  • User profiles                                                         │   │
│  │  • News feed                                                             │   │
│  │  • Friend lists                                                          │   │
│  │  • Notification counts                                                   │   │
│  │  • Trending topics                                                       │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  BANKING/FINANCIAL                                                       │   │
│  │  • Exchange rates                                                        │   │
│  │  • Account balances (with care)                                          │   │
│  │  • Transaction history                                                   │   │
│  │  • Reference data (branch info, currency codes)                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  CONTENT DELIVERY                                                        │   │
│  │  • Static assets (images, CSS, JS)                                       │   │
│  │  • API responses                                                         │   │
│  │  • Database query results                                                │   │
│  │  • Computed/aggregated data                                              │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.5 When to Use Caching

**Good Candidates for Caching:**
- Data that is read frequently but modified rarely
- Expensive computations or database queries
- Data that can tolerate some staleness
- Reference data (country codes, product categories)
- Session data
- API responses from external services

**Poor Candidates for Caching:**
- Frequently changing data
- Data requiring real-time accuracy (stock prices, inventory counts)
- Unique/one-time requests
- Large datasets that would exhaust cache memory
- Sensitive data with strict security requirements

---

## 2. Core Caching Concepts

### 2.1 Cache Key

The **cache key** is a unique identifier used to store and retrieve cached data.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         CACHE KEY STRUCTURE                                      │
│                                                                                  │
│   GOOD KEY DESIGN:                                                              │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  [namespace]:[entity]:[identifier]:[version]                            │  │
│   │                                                                         │  │
│   │  Examples:                                                              │  │
│   │  • products:catalog:12345                                               │  │
│   │  • users:profile:user_abc123                                            │  │
│   │  • orders:summary:ORD-2024-001                                          │  │
│   │  • api:weather:NYC:v2                                                   │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   KEY GENERATION STRATEGIES:                                                     │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  Method signature based:  getUserById_123                               │  │
│   │  Entity based:            User#123                                      │  │
│   │  Hash based:              SHA256(params)                                │  │
│   │  Composite:               users:region:US:state:CA:123                  │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   BEST PRACTICES:                                                                │
│   ✓ Keep keys short but descriptive                                             │
│   ✓ Use consistent naming convention                                            │
│   ✓ Include version for schema changes                                          │
│   ✓ Avoid special characters                                                    │
│   ✗ Don't use sensitive data in keys                                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Cache Key Examples in Spring Boot
@Service
public class ProductService {
    
    // Simple key - uses method parameter
    @Cacheable(value = "products", key = "#id")
    public Product getProductById(Long id) {
        return repository.findById(id).orElseThrow();
    }
    
    // Composite key - multiple parameters
    @Cacheable(value = "products", key = "#category + ':' + #page + ':' + #size")
    public Page<Product> getProductsByCategory(String category, int page, int size) {
        return repository.findByCategory(category, PageRequest.of(page, size));
    }
    
    // Custom key with SpEL
    @Cacheable(value = "users", key = "T(java.lang.String).valueOf(#user.id) + ':' + #user.region")
    public UserProfile getUserProfile(User user) {
        return profileService.loadProfile(user);
    }
}
```

### 2.2 Cache Eviction

**Cache eviction** is the process of removing entries from the cache to free up memory or invalidate stale data.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         CACHE EVICTION POLICIES                                  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  LRU (Least Recently Used)                                              │  │
│   │  • Evicts entries not accessed for longest time                         │  │
│   │  • Good for: General-purpose caching                                    │  │
│   │  • Memory: [A B C D E] → Access A → [B C D E A] → Add F → [C D E A F]  │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  LFU (Least Frequently Used)                                            │  │
│   │  • Evicts entries accessed least number of times                        │  │
│   │  • Good for: Data with varying access patterns                          │  │
│   │  • Tracks: Access count per entry                                       │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  FIFO (First In First Out)                                              │  │
│   │  • Evicts oldest entries first                                          │  │
│   │  • Good for: Time-sensitive data                                        │  │
│   │  • Simple but may evict frequently used items                           │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  TTL (Time To Live)                                                     │  │
│   │  • Evicts entries after specified time                                  │  │
│   │  • Good for: Data that becomes stale over time                          │  │
│   │  • Can combine with other policies                                      │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  Size-Based                                                             │  │
│   │  • Evicts when cache reaches max size                                   │  │
│   │  • Usually combined with LRU or LFU                                     │  │
│   │  • Good for: Memory-constrained environments                            │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 TTL (Time To Live)

**TTL** defines how long a cached entry remains valid before automatic expiration.

```java
// TTL Configuration Examples

// 1. Caffeine Cache with TTL
@Configuration
public class CacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        CaffeineCacheManager cacheManager = new CaffeineCacheManager();
        cacheManager.setCaffeine(Caffeine.newBuilder()
            .expireAfterWrite(Duration.ofMinutes(10))    // TTL: 10 minutes
            .expireAfterAccess(Duration.ofMinutes(5))    // Idle timeout: 5 minutes
            .maximumSize(1000));
        return cacheManager;
    }
}

// 2. Redis Cache with TTL
@Configuration
public class RedisCacheConfig {
    
    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        RedisCacheConfiguration config = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofHours(1))  // Default TTL: 1 hour
            .disableCachingNullValues();
        
        // Different TTLs for different caches
        Map<String, RedisCacheConfiguration> cacheConfigs = new HashMap<>();
        cacheConfigs.put("products", config.entryTtl(Duration.ofHours(24)));
        cacheConfigs.put("sessions", config.entryTtl(Duration.ofMinutes(30)));
        cacheConfigs.put("shortLived", config.entryTtl(Duration.ofMinutes(5)));
        
        return RedisCacheManager.builder(factory)
            .cacheDefaults(config)
            .withInitialCacheConfigurations(cacheConfigs)
            .build();
    }
}
```

### 2.4 Cache Invalidation Strategies

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE INVALIDATION STRATEGIES                                 │
│                                                                                  │
│   "There are only two hard things in Computer Science:                          │
│    cache invalidation and naming things." — Phil Karlton                        │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  1. TIME-BASED INVALIDATION (TTL)                                       │  │
│   │     • Simplest approach                                                 │  │
│   │     • Data expires after fixed time                                     │  │
│   │     • May serve stale data until expiry                                 │  │
│   │     • Best for: Data that changes predictably                           │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  2. EVENT-BASED INVALIDATION                                            │  │
│   │     • Invalidate on specific events (update, delete)                    │  │
│   │     • Most accurate but complex                                         │  │
│   │     • Best for: Data requiring consistency                              │  │
│   │                                                                         │  │
│   │     [Update Event] ──► [Message Queue] ──► [Cache Eviction]            │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  3. VERSION-BASED INVALIDATION                                          │  │
│   │     • Include version in cache key                                      │  │
│   │     • Increment version on data change                                  │  │
│   │     • Old versions naturally expire                                     │  │
│   │                                                                         │  │
│   │     Key: products:v1:123 → products:v2:123                             │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  4. WRITE-THROUGH INVALIDATION                                          │  │
│   │     • Update cache when database is updated                             │  │
│   │     • Keeps cache always fresh                                          │  │
│   │     • Higher write latency                                              │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Event-Based Cache Invalidation Example
@Service
public class ProductService {
    
    private final ApplicationEventPublisher eventPublisher;
    
    @CacheEvict(value = "products", key = "#id")
    @Transactional
    public Product updateProduct(Long id, ProductRequest request) {
        Product updated = repository.save(mapToEntity(id, request));
        
        // Publish event for distributed cache invalidation
        eventPublisher.publishEvent(new ProductUpdatedEvent(id));
        
        return updated;
    }
}

// Event Listener for distributed invalidation
@Component
public class CacheInvalidationListener {
    
    private final CacheManager cacheManager;
    
    @EventListener
    @Async
    public void handleProductUpdate(ProductUpdatedEvent event) {
        Cache cache = cacheManager.getCache("products");
        if (cache != null) {
            cache.evict(event.getProductId());
        }
    }
}
```

### 2.5 Strong vs Eventual Consistency

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CONSISTENCY MODELS IN CACHING                                 │
│                                                                                  │
│   STRONG CONSISTENCY:                                                            │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  • Cache and database always in sync                                    │  │
│   │  • Every read returns most recent write                                 │  │
│   │  • Higher latency due to synchronization                                │  │
│   │  • Use: Financial transactions, inventory                               │  │
│   │                                                                         │  │
│   │  [Write] ──► [Update DB] ──► [Update Cache] ──► [Acknowledge]          │  │
│   │              (synchronous - wait for both)                              │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   EVENTUAL CONSISTENCY:                                                          │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  • Cache may temporarily have stale data                                │  │
│   │  • System converges to consistent state over time                       │  │
│   │  • Lower latency, higher availability                                   │  │
│   │  • Use: Product catalog, user profiles, social feeds                    │  │
│   │                                                                         │  │
│   │  [Write] ──► [Update DB] ──► [Acknowledge] ──► [Async Cache Update]    │  │
│   │              (asynchronous - don't wait for cache)                      │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   COMPARISON:                                                                    │
│   ─────────────────────────────────────────────────────────────────────────     │
│   Aspect              │ Strong           │ Eventual                             │
│   ─────────────────────────────────────────────────────────────────────────     │
│   Read accuracy       │ Always current   │ May be stale                         │
│   Write latency       │ Higher           │ Lower                                │
│   Availability        │ Lower            │ Higher                               │
│   Implementation      │ Complex          │ Simpler                              │
│   Use case            │ Critical data    │ Most applications                    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.6 Local vs Distributed Cache

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    LOCAL VS DISTRIBUTED CACHE                                    │
│                                                                                  │
│   LOCAL CACHE (In-Process):                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                         │  │
│   │    ┌────────────────────────┐    ┌────────────────────────┐            │  │
│   │    │      Instance 1        │    │      Instance 2        │            │  │
│   │    │  ┌──────────────────┐  │    │  ┌──────────────────┐  │            │  │
│   │    │  │  Application     │  │    │  │  Application     │  │            │  │
│   │    │  │  ┌────────────┐  │  │    │  │  ┌────────────┐  │  │            │  │
│   │    │  │  │ Local Cache│  │  │    │  │  │ Local Cache│  │  │            │  │
│   │    │  │  └────────────┘  │  │    │  │  └────────────┘  │  │            │  │
│   │    │  └──────────────────┘  │    │  └──────────────────┘  │            │  │
│   │    └────────────────────────┘    └────────────────────────┘            │  │
│   │                                                                         │  │
│   │    • Each instance has own cache (not shared)                          │  │
│   │    • Fastest access (no network)                                       │  │
│   │    • Cache inconsistency between instances                             │  │
│   │    • Examples: Caffeine, Ehcache (local mode), Guava                   │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   DISTRIBUTED CACHE (Out-of-Process):                                           │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                         │  │
│   │    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │  │
│   │    │  Instance 1  │    │  Instance 2  │    │  Instance 3  │            │  │
│   │    └──────┬───────┘    └──────┬───────┘    └──────┬───────┘            │  │
│   │           │                   │                   │                     │  │
│   │           └───────────────────┼───────────────────┘                     │  │
│   │                               │                                         │  │
│   │                     ┌─────────▼─────────┐                               │  │
│   │                     │  DISTRIBUTED      │                               │  │
│   │                     │  CACHE CLUSTER    │                               │  │
│   │                     │  (Redis/Hazelcast)│                               │  │
│   │                     └───────────────────┘                               │  │
│   │                                                                         │  │
│   │    • Shared cache across all instances                                 │  │
│   │    • Network latency overhead                                          │  │
│   │    • Consistent view of data                                           │  │
│   │    • Examples: Redis, Hazelcast, Memcached                             │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   HYBRID (L1 + L2 Cache):                                                        │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │    Instance with L1 (local) + L2 (distributed) cache                   │  │
│   │                                                                         │  │
│   │    [Request] → [L1 Local Cache] → miss → [L2 Distributed] → miss → [DB]│  │
│   │                       ↓ hit                      ↓ hit                  │  │
│   │                   [Return]                   [Return + populate L1]     │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.7 L1 vs L2 Cache

| Aspect | L1 Cache (Local) | L2 Cache (Distributed) |
|--------|------------------|------------------------|
| **Location** | In-process memory | External cache server |
| **Speed** | ~1ms | ~5-20ms |
| **Size** | Limited (MB) | Large (GB-TB) |
| **Consistency** | Instance-local | Shared across instances |
| **Survivability** | Lost on restart | Persists across restarts |
| **Use Case** | Hot data, frequent access | Shared state, large datasets |

```java
// L1 + L2 Hybrid Cache Implementation
@Configuration
public class HybridCacheConfig {
    
    @Bean
    public CacheManager cacheManager(RedisConnectionFactory redisFactory) {
        // L1: Caffeine (local, fast)
        CaffeineCacheManager l1CacheManager = new CaffeineCacheManager();
        l1CacheManager.setCaffeine(Caffeine.newBuilder()
            .maximumSize(1000)
            .expireAfterWrite(Duration.ofMinutes(5)));  // Short TTL for L1
        
        // L2: Redis (distributed, shared)
        RedisCacheManager l2CacheManager = RedisCacheManager.builder(redisFactory)
            .cacheDefaults(RedisCacheConfiguration.defaultCacheConfig()
                .entryTtl(Duration.ofHours(1)))  // Longer TTL for L2
            .build();
        
        // Composite cache manager
        return new CompositeCacheManager(l1CacheManager, l2CacheManager);
    }
}
```

### 2.8 Cache Hit vs Cache Miss

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE HIT VS CACHE MISS                                       │
│                                                                                  │
│   CACHE HIT:                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  Data found in cache → Return immediately                               │  │
│   │                                                                         │  │
│   │  [Request] ──► [Cache] ──► [Data Found] ──► [Return]                   │  │
│   │                                                                         │  │
│   │  Latency: 1-10ms | Cost: Very Low                                      │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   CACHE MISS:                                                                    │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  Data NOT found in cache → Fetch from source → Store in cache          │  │
│   │                                                                         │  │
│   │  [Request] ──► [Cache] ──► [Not Found] ──► [DB] ──► [Store] ──► [Return]│  │
│   │                                                                         │  │
│   │  Latency: 50-500ms | Cost: Higher (DB query + cache write)             │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   KEY METRICS:                                                                   │
│   ─────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│   Cache Hit Ratio = (Cache Hits) / (Cache Hits + Cache Misses) × 100           │
│                                                                                  │
│   Target Hit Ratios:                                                            │
│   • Good:      > 80%                                                            │
│   • Excellent: > 95%                                                            │
│   • Needs work: < 70%                                                           │
│                                                                                  │
│   Low hit ratio causes:                                                         │
│   • TTL too short                                                               │
│   • Cache size too small                                                        │
│   • Poor key design                                                             │
│   • Data not suitable for caching                                               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.9 Cold Cache vs Warm Cache

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    COLD CACHE VS WARM CACHE                                      │
│                                                                                  │
│   COLD CACHE:                                                                    │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  • Empty or recently cleared cache                                      │  │
│   │  • Occurs after: Application restart, cache flush, new deployment       │  │
│   │  • All requests result in cache misses                                  │  │
│   │  • High database load initially                                         │  │
│   │                                                                         │  │
│   │  [App Start] → [Cache Empty] → [100% Miss Rate] → [DB Overload Risk]   │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   WARM CACHE:                                                                    │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  • Cache populated with frequently accessed data                        │  │
│   │  • High hit ratio (>80%)                                                │  │
│   │  • Normal operating state                                               │  │
│   │  • Low database load                                                    │  │
│   │                                                                         │  │
│   │  [Steady State] → [Cache Populated] → [High Hit Rate] → [Fast Response]│  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   CACHE WARMING STRATEGIES:                                                      │
│   ─────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│   1. Lazy Loading: Let cache populate naturally with traffic                    │
│      - Simple but slow to warm up                                               │
│                                                                                  │
│   2. Pre-Loading: Load expected data at startup                                 │
│      - Fast warm-up but may load unnecessary data                               │
│                                                                                  │
│   3. Scheduled Refresh: Periodically refresh hot data                           │
│      - Maintains warm cache, prevents mass expiration                           │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Cache Warming at Application Startup
@Component
public class CacheWarmer {
    
    private final ProductService productService;
    private final CategoryService categoryService;
    
    @EventListener(ApplicationReadyEvent.class)
    public void warmUpCache() {
        log.info("Starting cache warm-up...");
        
        // Pre-load frequently accessed data
        CompletableFuture.runAsync(() -> {
            // Warm up popular products
            List<Long> popularProductIds = getPopularProductIds();
            popularProductIds.forEach(productService::getProductById);
            
            // Warm up all categories
            categoryService.getAllCategories();
            
            log.info("Cache warm-up completed. Loaded {} products", 
                popularProductIds.size());
        });
    }
    
    private List<Long> getPopularProductIds() {
        // Return top 100 most accessed product IDs
        return analyticsService.getTopProductIds(100);
    }
}
```

---

## 3. How Caching Works in Spring Boot

### 3.1 Spring Cache Abstraction

Spring provides a **cache abstraction** layer that allows you to add caching to your application without being tied to a specific cache provider. This abstraction is similar to Spring's transaction abstraction.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SPRING CACHE ABSTRACTION ARCHITECTURE                         │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                        APPLICATION CODE                                  │  │
│   │            @Cacheable  @CachePut  @CacheEvict  @Caching                 │  │
│   └────────────────────────────────┬────────────────────────────────────────┘  │
│                                    │                                            │
│                                    ▼                                            │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                     SPRING CACHE ABSTRACTION                             │  │
│   │                                                                          │  │
│   │  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐     │  │
│   │  │  CacheManager   │    │     Cache       │    │  KeyGenerator   │     │  │
│   │  │   Interface     │    │   Interface     │    │    Interface    │     │  │
│   │  └─────────────────┘    └─────────────────┘    └─────────────────┘     │  │
│   │                                                                          │  │
│   └────────────────────────────────┬────────────────────────────────────────┘  │
│                                    │                                            │
│              ┌─────────────────────┼─────────────────────┐                     │
│              │                     │                     │                     │
│              ▼                     ▼                     ▼                     │
│   ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐            │
│   │  SimpleCacheMan  │  │ CaffeineCacheMgr │  │  RedisCacheMan   │            │
│   │  (ConcurrentMap) │  │                  │  │                  │            │
│   └──────────────────┘  └──────────────────┘  └──────────────────┘            │
│              │                     │                     │                     │
│              ▼                     ▼                     ▼                     │
│   ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐            │
│   │   ConcurrentMap  │  │    Caffeine      │  │      Redis       │            │
│   │   (JVM Memory)   │  │   (JVM Memory)   │  │    (External)    │            │
│   └──────────────────┘  └──────────────────┘  └──────────────────┘            │
│                                                                                  │
│   KEY COMPONENTS:                                                                │
│   • CacheManager: Creates and manages Cache instances                           │
│   • Cache: Represents a single cache (collection of key-value pairs)            │
│   • KeyGenerator: Generates cache keys from method parameters                   │
│   • CacheResolver: Determines which caches to use at runtime                    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Enabling Caching

```java
// Step 1: Add dependency (Spring Boot Starter Cache)
// pom.xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-cache</artifactId>
</dependency>

// Step 2: Enable caching in your application
@SpringBootApplication
@EnableCaching  // Enables Spring's annotation-driven cache management
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

### 3.3 Key Annotations

#### 3.3.1 @EnableCaching

Enables Spring's annotation-driven cache management capability.

```java
@Configuration
@EnableCaching
public class CacheConfig {
    
    // Optional: Customize cache behavior
    @Bean
    public CacheManager cacheManager() {
        SimpleCacheManager cacheManager = new SimpleCacheManager();
        cacheManager.setCaches(List.of(
            new ConcurrentMapCache("products"),
            new ConcurrentMapCache("users"),
            new ConcurrentMapCache("categories")
        ));
        return cacheManager;
    }
}
```

#### 3.3.2 @Cacheable

Caches the result of a method call. If the cache contains the value, the method is NOT executed.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         @CACHEABLE FLOW                                          │
│                                                                                  │
│   [Method Call] ──► [Check Cache] ──┬──► [Cache Hit] ──► [Return Cached Value]  │
│                                     │                                            │
│                                     └──► [Cache Miss] ──► [Execute Method] ─────│
│                                                               │                  │
│                                                               ▼                  │
│                                                         [Store Result in Cache] │
│                                                               │                  │
│                                                               ▼                  │
│                                                         [Return Result]         │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
@Service
@Slf4j
public class ProductService {
    
    private final ProductRepository repository;
    
    // Basic usage - caches result using 'id' as key
    @Cacheable("products")
    public Product getProductById(Long id) {
        log.info("Fetching product from database: {}", id);
        return repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
    }
    
    // With explicit key
    @Cacheable(value = "products", key = "#id")
    public Product findProduct(Long id) {
        return repository.findById(id).orElseThrow();
    }
    
    // Multiple parameters - composite key
    @Cacheable(value = "productSearch", key = "#category + ':' + #page")
    public Page<Product> searchProducts(String category, int page, int size) {
        log.info("Searching products in category: {}, page: {}", category, page);
        return repository.findByCategory(category, PageRequest.of(page, size));
    }
    
    // Conditional caching - only cache if result is not null
    @Cacheable(value = "products", key = "#sku", unless = "#result == null")
    public Product getProductBySku(String sku) {
        return repository.findBySku(sku).orElse(null);
    }
    
    // Conditional caching - only cache for active products
    @Cacheable(value = "products", condition = "#active == true")
    public List<Product> getProducts(boolean active) {
        return repository.findByActive(active);
    }
    
    // Cache with sync - prevents cache stampede
    @Cacheable(value = "products", key = "#id", sync = true)
    public Product getProductSync(Long id) {
        return repository.findById(id).orElseThrow();
    }
}
```

#### 3.3.3 @CachePut

Always executes the method and updates the cache with the result. Unlike @Cacheable, it doesn't skip method execution.

```java
@Service
public class ProductService {
    
    // Always executes and updates cache
    @CachePut(value = "products", key = "#result.id")
    public Product createProduct(ProductRequest request) {
        Product product = mapToEntity(request);
        return repository.save(product);
    }
    
    // Update cache after modification
    @CachePut(value = "products", key = "#id")
    public Product updateProduct(Long id, ProductRequest request) {
        Product existing = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        updateEntity(existing, request);
        return repository.save(existing);
    }
    
    // Conditional update - only cache if product is active
    @CachePut(value = "products", key = "#result.id", condition = "#result.active")
    public Product updateProductStatus(Long id, boolean active) {
        Product product = repository.findById(id).orElseThrow();
        product.setActive(active);
        return repository.save(product);
    }
}
```

#### 3.3.4 @CacheEvict

Removes entries from the cache. Use this when data is updated or deleted.

```java
@Service
public class ProductService {
    
    // Evict single entry
    @CacheEvict(value = "products", key = "#id")
    public void deleteProduct(Long id) {
        repository.deleteById(id);
    }
    
    // Evict all entries in a cache
    @CacheEvict(value = "products", allEntries = true)
    public void clearProductCache() {
        log.info("Product cache cleared");
    }
    
    // Evict before method execution (default is after)
    @CacheEvict(value = "products", key = "#id", beforeInvocation = true)
    public void deleteProductSafe(Long id) {
        // Even if this throws exception, cache will be evicted
        repository.deleteById(id);
    }
    
    // Evict multiple caches
    @CacheEvict(value = {"products", "productSearch", "categories"}, allEntries = true)
    @Scheduled(cron = "0 0 2 * * ?")  // Daily at 2 AM
    public void refreshAllCaches() {
        log.info("All caches refreshed");
    }
    
    // Evict on update
    @CacheEvict(value = "products", key = "#id")
    @Transactional
    public Product updateProduct(Long id, ProductRequest request) {
        Product product = repository.findById(id).orElseThrow();
        updateEntity(product, request);
        return repository.save(product);
    }
}
```

#### 3.3.5 @Caching

Combines multiple cache operations in a single method.

```java
@Service
public class ProductService {
    
    // Multiple cache operations
    @Caching(
        put = {
            @CachePut(value = "products", key = "#result.id"),
            @CachePut(value = "productsBySku", key = "#result.sku")
        },
        evict = {
            @CacheEvict(value = "productSearch", allEntries = true),
            @CacheEvict(value = "productsByCategory", allEntries = true)
        }
    )
    public Product createProduct(ProductRequest request) {
        Product product = mapToEntity(request);
        return repository.save(product);
    }
    
    // Evict from multiple caches on delete
    @Caching(evict = {
        @CacheEvict(value = "products", key = "#id"),
        @CacheEvict(value = "productsBySku", key = "#sku"),
        @CacheEvict(value = "productSearch", allEntries = true)
    })
    public void deleteProduct(Long id, String sku) {
        repository.deleteById(id);
    }
    
    // Cacheable with fallback
    @Caching(cacheable = {
        @Cacheable(value = "products", key = "#id"),
        @Cacheable(value = "productsBackup", key = "#id")
    })
    public Product getProductWithFallback(Long id) {
        return repository.findById(id).orElseThrow();
    }
}
```

#### 3.3.6 @CacheConfig

Class-level annotation to share common cache-related settings.

```java
@Service
@CacheConfig(cacheNames = "products", keyGenerator = "customKeyGenerator")
public class ProductService {
    
    // Uses cache name "products" from class-level config
    @Cacheable
    public Product getProductById(Long id) {
        return repository.findById(id).orElseThrow();
    }
    
    // Override class-level cache name
    @Cacheable(cacheNames = "productSearch")
    public List<Product> searchProducts(String keyword) {
        return repository.searchByKeyword(keyword);
    }
    
    @CacheEvict
    public void deleteProduct(Long id) {
        repository.deleteById(id);
    }
}
```

### 3.4 Cache Key Generation

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE KEY GENERATION                                          │
│                                                                                  │
│   DEFAULT KEY GENERATION:                                                        │
│   ─────────────────────────────────────────────────────────────────────────     │
│   • No params:     Uses SimpleKey.EMPTY                                         │
│   • One param:     Uses that parameter directly                                 │
│   • Multiple:      Creates SimpleKey with all params                            │
│                                                                                  │
│   EXAMPLES:                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │  getProduct()                    → SimpleKey.EMPTY                      │  │
│   │  getProduct(123L)                → 123L                                 │  │
│   │  getProduct("SKU", 123L)         → SimpleKey["SKU", 123L]              │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   SpEL KEY EXPRESSIONS:                                                          │
│   ─────────────────────────────────────────────────────────────────────────     │
│   #paramName          → Parameter by name                                        │
│   #p0, #p1            → Parameter by index                                       │
│   #a0, #a1            → Alias for parameter index                                │
│   #root.method        → Method being invoked                                     │
│   #root.target        → Target object                                            │
│   #root.caches        → Caches being used                                        │
│   #root.methodName    → Method name                                              │
│   #result             → Method result (only in @CachePut, unless)               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
@Service
public class ProductService {
    
    // Default key - uses parameter directly
    @Cacheable("products")
    public Product getById(Long id) { ... }  // Key: id value
    
    // Explicit key using SpEL
    @Cacheable(value = "products", key = "#id")
    public Product findById(Long id) { ... }
    
    // Multiple parameters - composite key
    @Cacheable(value = "products", key = "#category + ':' + #status")
    public List<Product> find(String category, String status) { ... }
    
    // Using object properties
    @Cacheable(value = "products", key = "#product.sku")
    public Product save(Product product) { ... }
    
    // Complex key expression
    @Cacheable(value = "searchResults", 
               key = "T(java.util.Objects).hash(#criteria.category, #criteria.minPrice, #criteria.maxPrice)")
    public List<Product> search(SearchCriteria criteria) { ... }
    
    // Using result in key (CachePut)
    @CachePut(value = "products", key = "#result.id")
    public Product create(ProductRequest request) { ... }
    
    // Using method name in key
    @Cacheable(value = "analytics", key = "#root.methodName + ':' + #id")
    public Analytics getAnalytics(Long id) { ... }
}
```

### 3.5 Custom Key Generator

```java
// Custom Key Generator Implementation
@Component("customKeyGenerator")
public class CustomKeyGenerator implements KeyGenerator {
    
    @Override
    public Object generate(Object target, Method method, Object... params) {
        StringBuilder keyBuilder = new StringBuilder();
        
        // Include class name
        keyBuilder.append(target.getClass().getSimpleName());
        keyBuilder.append(":");
        
        // Include method name
        keyBuilder.append(method.getName());
        keyBuilder.append(":");
        
        // Include parameters
        for (Object param : params) {
            if (param != null) {
                keyBuilder.append(param.toString());
                keyBuilder.append(":");
            }
        }
        
        return keyBuilder.toString();
    }
}

// Usage
@Service
public class ProductService {
    
    @Cacheable(value = "products", keyGenerator = "customKeyGenerator")
    public Product getProduct(Long id, String region) {
        // Key will be: ProductService:getProduct:123:US:
        return repository.findByIdAndRegion(id, region);
    }
}
```

```java
// Advanced Key Generator with Hashing
@Component("hashKeyGenerator")
public class HashKeyGenerator implements KeyGenerator {
    
    private final ObjectMapper objectMapper = new ObjectMapper();
    
    @Override
    public Object generate(Object target, Method method, Object... params) {
        try {
            String methodSignature = target.getClass().getName() + "." + method.getName();
            String paramsJson = objectMapper.writeValueAsString(params);
            String combined = methodSignature + paramsJson;
            
            // Generate SHA-256 hash for consistent, fixed-length keys
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(combined.getBytes(StandardCharsets.UTF_8));
            
            return Base64.getEncoder().encodeToString(hash);
        } catch (Exception e) {
            throw new RuntimeException("Key generation failed", e);
        }
    }
}
```

### 3.6 Conditional Caching

```java
@Service
public class ProductService {
    
    // condition: Evaluated BEFORE method execution
    // Unless result is determined, condition must be true for caching
    @Cacheable(value = "products", 
               condition = "#id > 0")  // Only cache if id is positive
    public Product getProduct(Long id) {
        return repository.findById(id).orElseThrow();
    }
    
    // unless: Evaluated AFTER method execution
    // Cache unless the condition is true
    @Cacheable(value = "products", 
               unless = "#result == null")  // Don't cache null results
    public Product findProduct(Long id) {
        return repository.findById(id).orElse(null);
    }
    
    // Combined conditions
    @Cacheable(value = "products",
               condition = "#active == true",     // Only consider caching if active=true
               unless = "#result.price < 10.0")   // Don't cache cheap products
    public Product getCacheableProduct(boolean active) {
        return repository.findRandomProduct(active);
    }
    
    // Conditional eviction
    @CacheEvict(value = "products", 
                key = "#product.id",
                condition = "#product.active == false")  // Only evict inactive products
    public void updateProduct(Product product) {
        repository.save(product);
    }
    
    // Complex condition with SpEL
    @Cacheable(value = "reports",
               condition = "#dateRange.daysBetween() <= 30",  // Only cache short ranges
               unless = "#result.isEmpty()")
    public List<Report> getReports(DateRange dateRange) {
        return reportRepository.findByDateRange(dateRange);
    }
}
```

### 3.7 Cache Manager Configuration

```java
@Configuration
@EnableCaching
public class CacheConfig {
    
    // Simple CacheManager with ConcurrentMap
    @Bean
    public CacheManager simpleCacheManager() {
        SimpleCacheManager cacheManager = new SimpleCacheManager();
        cacheManager.setCaches(List.of(
            new ConcurrentMapCache("products"),
            new ConcurrentMapCache("users"),
            new ConcurrentMapCache("categories")
        ));
        return cacheManager;
    }
}
```

```java
// Caffeine Cache Manager with custom configuration
@Configuration
@EnableCaching
public class CaffeineCacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        CaffeineCacheManager cacheManager = new CaffeineCacheManager();
        
        // Global settings
        cacheManager.setCaffeine(Caffeine.newBuilder()
            .initialCapacity(100)
            .maximumSize(500)
            .expireAfterWrite(Duration.ofMinutes(10))
            .expireAfterAccess(Duration.ofMinutes(5))
            .recordStats());  // Enable statistics
        
        // Pre-define cache names
        cacheManager.setCacheNames(List.of("products", "users", "categories"));
        
        return cacheManager;
    }
    
    // Multiple caches with different settings
    @Bean
    public CacheManager multiConfigCacheManager() {
        SimpleCacheManager cacheManager = new SimpleCacheManager();
        
        List<CaffeineCache> caches = List.of(
            buildCache("products", 1000, Duration.ofHours(1)),
            buildCache("users", 500, Duration.ofMinutes(30)),
            buildCache("sessions", 10000, Duration.ofMinutes(15)),
            buildCache("shortLived", 100, Duration.ofMinutes(1))
        );
        
        cacheManager.setCaches(caches);
        return cacheManager;
    }
    
    private CaffeineCache buildCache(String name, int maxSize, Duration ttl) {
        return new CaffeineCache(name, Caffeine.newBuilder()
            .maximumSize(maxSize)
            .expireAfterWrite(ttl)
            .recordStats()
            .build());
    }
}
```

### 3.8 Configuration via application.yml

```yaml
# application.yml
spring:
  cache:
    type: caffeine  # or redis, ehcache, simple, none
    cache-names:
      - products
      - users
      - categories
      - sessions
    
    # Caffeine specific settings
    caffeine:
      spec: maximumSize=500,expireAfterWrite=600s,recordStats
    
    # Redis specific settings (when type: redis)
    redis:
      time-to-live: 3600000  # 1 hour in milliseconds
      cache-null-values: false
      use-key-prefix: true
      key-prefix: "myapp:"

---
# Profile-specific configuration
spring:
  profiles: development
  cache:
    type: simple  # Use simple cache in dev

---
spring:
  profiles: production
  cache:
    type: redis  # Use Redis in production
    redis:
      time-to-live: 7200000  # 2 hours
```

### 3.9 Complete Practical Examples

#### 3.9.1 Simple Service-Level Caching

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class ProductService {
    
    private final ProductRepository repository;
    
    // Cache product by ID
    @Cacheable(value = "products", key = "#id")
    public ProductDTO getProductById(Long id) {
        log.info("Cache MISS - Fetching product from DB: {}", id);
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        return mapToDTO(product);
    }
    
    // Cache product by SKU
    @Cacheable(value = "productsBySku", key = "#sku", unless = "#result == null")
    public ProductDTO getProductBySku(String sku) {
        log.info("Cache MISS - Fetching product by SKU: {}", sku);
        return repository.findBySku(sku)
            .map(this::mapToDTO)
            .orElse(null);
    }
    
    // Cache list of products by category
    @Cacheable(value = "productsByCategory", key = "#category")
    public List<ProductDTO> getProductsByCategory(String category) {
        log.info("Cache MISS - Fetching products for category: {}", category);
        return repository.findByCategory(category)
            .stream()
            .map(this::mapToDTO)
            .collect(Collectors.toList());
    }
    
    private ProductDTO mapToDTO(Product product) {
        return ProductDTO.builder()
            .id(product.getId())
            .name(product.getName())
            .price(product.getPrice())
            .category(product.getCategory())
            .sku(product.getSku())
            .build();
    }
}
```

#### 3.9.2 Cache Eviction Example

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class ProductCacheService {
    
    private final ProductRepository repository;
    
    // Evict single product from multiple caches
    @Caching(evict = {
        @CacheEvict(value = "products", key = "#id"),
        @CacheEvict(value = "productsByCategory", allEntries = true)
    })
    public void deleteProduct(Long id) {
        log.info("Deleting product and evicting from cache: {}", id);
        repository.deleteById(id);
    }
    
    // Update product - evict old cache and update with new value
    @Caching(
        evict = {
            @CacheEvict(value = "productsByCategory", allEntries = true),
            @CacheEvict(value = "productsBySku", key = "#request.oldSku", 
                        condition = "#request.oldSku != #request.newSku")
        },
        put = {
            @CachePut(value = "products", key = "#id"),
            @CachePut(value = "productsBySku", key = "#result.sku")
        }
    )
    public ProductDTO updateProduct(Long id, ProductUpdateRequest request) {
        log.info("Updating product: {}", id);
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        
        product.setName(request.getName());
        product.setPrice(request.getPrice());
        product.setCategory(request.getCategory());
        product.setSku(request.getNewSku());
        
        Product saved = repository.save(product);
        return mapToDTO(saved);
    }
    
    // Scheduled cache refresh
    @CacheEvict(value = {"products", "productsBySku", "productsByCategory"}, 
                allEntries = true)
    @Scheduled(fixedRate = 3600000)  // Every hour
    public void refreshCaches() {
        log.info("Scheduled cache refresh completed");
    }
}
```

#### 3.9.3 Custom Key Example

```java
@Service
@Slf4j
public class SearchService {
    
    private final SearchRepository searchRepository;
    
    // Custom key using SpEL
    @Cacheable(value = "searchResults", 
               key = "T(com.example.util.CacheKeyUtil).generateSearchKey(#criteria)")
    public SearchResult search(SearchCriteria criteria) {
        log.info("Cache MISS - Executing search: {}", criteria);
        return searchRepository.search(criteria);
    }
    
    // Pagination-aware caching
    @Cacheable(value = "pagedProducts",
               key = "'page:' + #pageable.pageNumber + ':size:' + #pageable.pageSize + ':sort:' + #pageable.sort.toString()")
    public Page<ProductDTO> getProductsPage(Pageable pageable) {
        log.info("Cache MISS - Fetching page: {}", pageable);
        return repository.findAll(pageable).map(this::mapToDTO);
    }
    
    // User-specific caching
    @Cacheable(value = "userRecommendations",
               key = "#userId + ':' + #category + ':' + #limit")
    public List<ProductDTO> getRecommendations(Long userId, String category, int limit) {
        log.info("Cache MISS - Generating recommendations for user: {}", userId);
        return recommendationEngine.generate(userId, category, limit);
    }
}

// Cache Key Utility
public class CacheKeyUtil {
    
    public static String generateSearchKey(SearchCriteria criteria) {
        return String.format("search:%s:%s:%s:%d:%d",
            criteria.getKeyword(),
            criteria.getCategory(),
            criteria.getSortBy(),
            criteria.getMinPrice(),
            criteria.getMaxPrice()
        );
    }
}
```

#### 3.9.4 Global Cache Configuration Example

```java
@Configuration
@EnableCaching
@Slf4j
public class GlobalCacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        SimpleCacheManager cacheManager = new SimpleCacheManager();
        
        List<Cache> caches = new ArrayList<>();
        
        // High-traffic, long TTL cache
        caches.add(createCaffeineCache("products", 
            CacheSpec.builder()
                .maxSize(10000)
                .expireAfterWrite(Duration.ofHours(2))
                .expireAfterAccess(Duration.ofMinutes(30))
                .build()));
        
        // Medium-traffic cache
        caches.add(createCaffeineCache("users",
            CacheSpec.builder()
                .maxSize(5000)
                .expireAfterWrite(Duration.ofMinutes(30))
                .build()));
        
        // Short-lived session cache
        caches.add(createCaffeineCache("sessions",
            CacheSpec.builder()
                .maxSize(50000)
                .expireAfterAccess(Duration.ofMinutes(15))
                .build()));
        
        // Search results - shorter TTL
        caches.add(createCaffeineCache("searchResults",
            CacheSpec.builder()
                .maxSize(1000)
                .expireAfterWrite(Duration.ofMinutes(5))
                .build()));
        
        cacheManager.setCaches(caches);
        return cacheManager;
    }
    
    private CaffeineCache createCaffeineCache(String name, CacheSpec spec) {
        log.info("Creating cache '{}' with spec: {}", name, spec);
        
        Caffeine<Object, Object> builder = Caffeine.newBuilder()
            .maximumSize(spec.getMaxSize())
            .recordStats();
        
        if (spec.getExpireAfterWrite() != null) {
            builder.expireAfterWrite(spec.getExpireAfterWrite());
        }
        if (spec.getExpireAfterAccess() != null) {
            builder.expireAfterAccess(spec.getExpireAfterAccess());
        }
        
        // Add removal listener for monitoring
        builder.removalListener((key, value, cause) -> 
            log.debug("Cache '{}' evicted key '{}' due to: {}", name, key, cause));
        
        return new CaffeineCache(name, builder.build());
    }
    
    @Bean
    public KeyGenerator customKeyGenerator() {
        return new CustomKeyGenerator();
    }
    
    // Cache error handler
    @Bean
    public CacheErrorHandler cacheErrorHandler() {
        return new CacheErrorHandler() {
            @Override
            public void handleCacheGetError(RuntimeException e, Cache cache, Object key) {
                log.error("Cache GET error - cache: {}, key: {}", cache.getName(), key, e);
            }
            
            @Override
            public void handleCachePutError(RuntimeException e, Cache cache, Object key, Object value) {
                log.error("Cache PUT error - cache: {}, key: {}", cache.getName(), key, e);
            }
            
            @Override
            public void handleCacheEvictError(RuntimeException e, Cache cache, Object key) {
                log.error("Cache EVICT error - cache: {}, key: {}", cache.getName(), key, e);
            }
            
            @Override
            public void handleCacheClearError(RuntimeException e, Cache cache) {
                log.error("Cache CLEAR error - cache: {}", cache.getName(), e);
            }
        };
    }
}

@Data
@Builder
class CacheSpec {
    private int maxSize;
    private Duration expireAfterWrite;
    private Duration expireAfterAccess;
}
```

---

## 4. Caching Providers in Spring Boot

Spring Boot supports multiple cache providers out of the box. The choice of provider depends on your application's requirements for performance, scalability, and features.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SPRING BOOT CACHE PROVIDERS OVERVIEW                          │
│                                                                                  │
│   IN-MEMORY (LOCAL):                    DISTRIBUTED:                            │
│   ┌───────────────────────────────┐    ┌───────────────────────────────┐       │
│   │  • ConcurrentMap (Simple)     │    │  • Redis                      │       │
│   │  • Caffeine                   │    │  • Hazelcast                  │       │
│   │  • Ehcache (local mode)       │    │  • Infinispan                 │       │
│   │  • Guava (deprecated)         │    │  • Memcached                  │       │
│   └───────────────────────────────┘    └───────────────────────────────┘       │
│                                                                                  │
│   PROVIDER DETECTION ORDER (Auto-configuration):                                │
│   1. Generic (if CacheManager bean exists)                                      │
│   2. JCache (JSR-107)                                                           │
│   3. Hazelcast                                                                  │
│   4. Infinispan                                                                 │
│   5. Couchbase                                                                  │
│   6. Redis                                                                      │
│   7. Caffeine                                                                   │
│   8. Simple (ConcurrentMap) - fallback                                         │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 4.1 ConcurrentMap (SimpleCacheManager)

#### Overview
The simplest cache provider using Java's `ConcurrentHashMap`. Provided by Spring Framework out of the box with no additional dependencies.

#### Architecture Type
**In-Memory (Local)** - Data stored in JVM heap memory.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CONCURRENT MAP CACHE ARCHITECTURE                             │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                           JVM Process                                    │  │
│   │                                                                          │  │
│   │   ┌─────────────────────┐                                               │  │
│   │   │    Application      │                                               │  │
│   │   │        │            │                                               │  │
│   │   │        ▼            │                                               │  │
│   │   │  ┌────────────┐     │     ┌─────────────────────────────────┐      │  │
│   │   │  │ Cache      │     │     │      Heap Memory                │      │  │
│   │   │  │ Manager    │─────┼────►│  ┌─────────────────────────┐    │      │  │
│   │   │  └────────────┘     │     │  │ ConcurrentHashMap       │    │      │  │
│   │   │                     │     │  │  Key1 → Value1          │    │      │  │
│   │   │                     │     │  │  Key2 → Value2          │    │      │  │
│   │   │                     │     │  │  Key3 → Value3          │    │      │  │
│   │   │                     │     │  └─────────────────────────┘    │      │  │
│   │   └─────────────────────┘     └─────────────────────────────────┘      │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Configuration

```java
@Configuration
@EnableCaching
public class SimpleCacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        SimpleCacheManager cacheManager = new SimpleCacheManager();
        cacheManager.setCaches(List.of(
            new ConcurrentMapCache("products"),
            new ConcurrentMapCache("users"),
            new ConcurrentMapCache("categories")
        ));
        return cacheManager;
    }
}
```

```yaml
# application.yml
spring:
  cache:
    type: simple
    cache-names: products, users, categories
```

#### Best Use Cases
- Development and testing environments
- Small applications with single instance
- Prototyping and POC
- Unit tests

#### Pros & Cons

| Pros | Cons |
|------|------|
| Zero dependencies | No TTL support |
| Simple configuration | No size limits |
| Fast (same JVM) | Not distributed |
| No serialization overhead | Memory grows unbounded |
| Thread-safe | No eviction policies |
| Good for testing | Data lost on restart |

#### When to Choose
✅ Development/testing environments  
✅ Single-instance applications with small datasets  
✅ When you need quick caching without external dependencies  
❌ Production multi-instance deployments  
❌ When you need TTL or eviction policies  

---

### 4.2 Caffeine

#### Overview
High-performance, near-optimal caching library for Java. It's the successor to Guava Cache and is the recommended local cache for Spring Boot applications.

#### Architecture Type
**In-Memory (Local)** - Uses sophisticated algorithms (TinyLFU) for optimal memory usage.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CAFFEINE CACHE ARCHITECTURE                                   │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                         Caffeine Cache                                   │  │
│   │                                                                          │  │
│   │   ┌─────────────────────────────────────────────────────────────────┐   │  │
│   │   │                    Window TinyLFU                                │   │  │
│   │   │                                                                  │   │  │
│   │   │  ┌─────────────┐   ┌──────────────┐   ┌───────────────────┐    │   │  │
│   │   │  │   Window    │   │   Probation  │   │    Protected      │    │   │  │
│   │   │  │   (1%)      │   │   (20%)      │   │    (80%)          │    │   │  │
│   │   │  │             │   │              │   │                   │    │   │  │
│   │   │  │  New items  │──►│  Candidates  │──►│  Frequently used  │    │   │  │
│   │   │  │  enter here │   │  for promote │   │  items stay here  │    │   │  │
│   │   │  └─────────────┘   └──────────────┘   └───────────────────┘    │   │  │
│   │   │                                                                  │   │  │
│   │   │  EVICTION: TinyLFU approximates optimal by tracking frequency   │   │  │
│   │   └─────────────────────────────────────────────────────────────────┘   │  │
│   │                                                                          │  │
│   │   FEATURES:                                                              │  │
│   │   • Size-based eviction          • Statistics collection                │  │
│   │   • Time-based expiration        • Async loading                        │  │
│   │   • Reference-based eviction     • Refresh after write                  │  │
│   │   • Removal listeners            • Weak/Soft references                 │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Dependency

```xml
<dependency>
    <groupId>com.github.ben-manes.caffeine</groupId>
    <artifactId>caffeine</artifactId>
</dependency>
```

#### Configuration

```java
@Configuration
@EnableCaching
public class CaffeineCacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        CaffeineCacheManager cacheManager = new CaffeineCacheManager();
        cacheManager.setCaffeine(caffeineCacheBuilder());
        cacheManager.setCacheNames(List.of("products", "users", "categories"));
        return cacheManager;
    }
    
    Caffeine<Object, Object> caffeineCacheBuilder() {
        return Caffeine.newBuilder()
            .initialCapacity(100)
            .maximumSize(1000)
            .expireAfterWrite(Duration.ofMinutes(10))
            .expireAfterAccess(Duration.ofMinutes(5))
            .weakKeys()
            .recordStats();
    }
}

// Multiple caches with different configurations
@Configuration
@EnableCaching
public class MultiCaffeineCacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        SimpleCacheManager manager = new SimpleCacheManager();
        manager.setCaches(List.of(
            buildCache("products", 10000, Duration.ofHours(1)),
            buildCache("users", 5000, Duration.ofMinutes(30)),
            buildCache("sessions", 50000, Duration.ofMinutes(15)),
            buildCache("shortLived", 500, Duration.ofMinutes(1))
        ));
        return manager;
    }
    
    private CaffeineCache buildCache(String name, int maxSize, Duration ttl) {
        return new CaffeineCache(name, 
            Caffeine.newBuilder()
                .maximumSize(maxSize)
                .expireAfterWrite(ttl)
                .recordStats()
                .build());
    }
}
```

```yaml
# application.yml
spring:
  cache:
    type: caffeine
    caffeine:
      spec: maximumSize=10000,expireAfterWrite=3600s,recordStats
    cache-names:
      - products
      - users
      - categories
```

#### Best Use Cases
- High-throughput single-instance applications
- L1 (local) cache in hybrid caching setup
- Session caching (single node)
- Frequently accessed reference data
- When optimal hit ratio is critical

#### Pros & Cons

| Pros | Cons |
|------|------|
| High performance (near-optimal) | Local only (not distributed) |
| Sophisticated eviction (TinyLFU) | Data lost on restart |
| Rich feature set (TTL, size limits) | Memory bounded by JVM |
| Built-in statistics | Not shared across instances |
| Async loading support | |
| Low memory footprint | |

#### When to Choose
✅ High-performance local caching needs  
✅ Single-instance applications  
✅ As L1 cache in multi-tier caching  
✅ When memory efficiency is important  
❌ Multi-instance deployments needing shared cache  
❌ When data persistence is required  

---

### 4.3 Ehcache

#### Overview
Robust, feature-rich, scalable cache that can operate as both in-memory (local) and distributed cache. Supports JSR-107 (JCache) specification.

#### Architecture Type
**In-Memory (Local)** or **Distributed** (with Terracotta)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    EHCACHE ARCHITECTURE                                          │
│                                                                                  │
│   STANDALONE MODE:                                                               │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │    ┌──────────────────────────────────────────────────────────────┐    │  │
│   │    │                    Ehcache (Local)                            │    │  │
│   │    │                                                               │    │  │
│   │    │  ┌─────────────┐   ┌─────────────┐   ┌─────────────────┐    │    │  │
│   │    │  │   Heap      │   │  Off-Heap   │   │     Disk        │    │    │  │
│   │    │  │   Tier      │──►│   Tier      │──►│     Tier        │    │    │  │
│   │    │  │ (fastest)   │   │ (large)     │   │  (persistent)   │    │    │  │
│   │    │  └─────────────┘   └─────────────┘   └─────────────────┘    │    │  │
│   │    │                                                               │    │  │
│   │    └──────────────────────────────────────────────────────────────┘    │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   CLUSTERED MODE (with Terracotta):                                             │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │  ┌────────────┐    ┌────────────┐    ┌────────────┐                    │  │
│   │  │ Instance 1 │    │ Instance 2 │    │ Instance 3 │                    │  │
│   │  │  Ehcache   │    │  Ehcache   │    │  Ehcache   │                    │  │
│   │  └─────┬──────┘    └─────┬──────┘    └─────┬──────┘                    │  │
│   │        │                 │                 │                            │  │
│   │        └─────────────────┼─────────────────┘                            │  │
│   │                          │                                              │  │
│   │                ┌─────────▼─────────┐                                    │  │
│   │                │  Terracotta       │                                    │  │
│   │                │  Server Array     │                                    │  │
│   │                │  (Distributed)    │                                    │  │
│   │                └───────────────────┘                                    │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Dependency

```xml
<dependency>
    <groupId>org.ehcache</groupId>
    <artifactId>ehcache</artifactId>
</dependency>
<!-- For JCache support -->
<dependency>
    <groupId>javax.cache</groupId>
    <artifactId>cache-api</artifactId>
</dependency>
```

#### Configuration

```xml
<!-- ehcache.xml -->
<config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xmlns="http://www.ehcache.org/v3"
        xsi:schemaLocation="http://www.ehcache.org/v3 
                            http://www.ehcache.org/schema/ehcache-core-3.0.xsd">

    <!-- Cache templates -->
    <cache-template name="defaultTemplate">
        <expiry>
            <ttl unit="minutes">10</ttl>
        </expiry>
        <heap unit="entries">1000</heap>
    </cache-template>

    <!-- Products cache with tiered storage -->
    <cache alias="products" uses-template="defaultTemplate">
        <heap unit="entries">10000</heap>
        <offheap unit="MB">100</offheap>
        <disk unit="GB" persistent="true">1</disk>
    </cache>

    <!-- Users cache -->
    <cache alias="users">
        <expiry>
            <ttl unit="minutes">30</ttl>
        </expiry>
        <heap unit="entries">5000</heap>
        <offheap unit="MB">50</offheap>
    </cache>

    <!-- Short-lived cache -->
    <cache alias="sessions">
        <expiry>
            <tti unit="minutes">15</tti>
        </expiry>
        <heap unit="entries">50000</heap>
    </cache>
</config>
```

```java
@Configuration
@EnableCaching
public class EhcacheConfig {
    
    @Bean
    public CacheManager cacheManager() {
        return new JCacheCacheManager(
            Caching.getCachingProvider().getCacheManager(
                getClass().getResource("/ehcache.xml").toURI(),
                getClass().getClassLoader()
            )
        );
    }
}
```

```yaml
# application.yml
spring:
  cache:
    type: jcache
    jcache:
      config: classpath:ehcache.xml
```

#### Best Use Cases
- Enterprise applications requiring tiered caching
- Applications needing disk persistence
- Hibernate second-level cache
- When off-heap memory storage is needed
- Large datasets that don't fit in heap

#### Pros & Cons

| Pros | Cons |
|------|------|
| Tiered storage (heap, off-heap, disk) | Complex configuration |
| Disk persistence | Heavier than Caffeine |
| JCache (JSR-107) compliant | Terracotta needed for clustering |
| Mature and battle-tested | Learning curve |
| Clustering support (with Terracotta) | |
| Rich management/monitoring | |

#### When to Choose
✅ Need for tiered caching (heap + off-heap + disk)  
✅ Data persistence requirement  
✅ Hibernate L2 cache  
✅ Enterprise applications with complex caching needs  
❌ Simple caching needs (use Caffeine)  
❌ Cloud-native distributed caching (use Redis)  

---

### 4.4 Redis

#### Overview
Open-source, in-memory data structure store used as a database, cache, and message broker. The most popular choice for distributed caching in microservices.

#### Architecture Type
**Distributed** - External cache server(s)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS CACHE ARCHITECTURE                                      │
│                                                                                  │
│   STANDALONE:                                                                    │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐              │  │
│   │  │ App Instance │    │ App Instance │    │ App Instance │              │  │
│   │  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘              │  │
│   │         │                   │                   │                       │  │
│   │         └───────────────────┼───────────────────┘                       │  │
│   │                             │                                           │  │
│   │                    ┌────────▼────────┐                                  │  │
│   │                    │  Redis Server   │                                  │  │
│   │                    │   (Primary)     │                                  │  │
│   │                    └─────────────────┘                                  │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   CLUSTER MODE:                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │  ┌────────────────────────────────────────────────────────────────────┐ │  │
│   │  │                      Redis Cluster                                  │ │  │
│   │  │                                                                     │ │  │
│   │  │  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐        │ │  │
│   │  │  │ Master 1│    │ Master 2│    │ Master 3│    │ Master N│        │ │  │
│   │  │  │ Slot    │    │ Slot    │    │ Slot    │    │ Slot    │        │ │  │
│   │  │  │ 0-5460  │    │ 5461-   │    │ 10923-  │    │  ...    │        │ │  │
│   │  │  │         │    │ 10922   │    │ 16383   │    │         │        │ │  │
│   │  │  └────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘        │ │  │
│   │  │       │              │              │              │              │ │  │
│   │  │  ┌────▼────┐    ┌────▼────┐    ┌────▼────┐    ┌────▼────┐        │ │  │
│   │  │  │Replica 1│    │Replica 2│    │Replica 3│    │Replica N│        │ │  │
│   │  │  └─────────┘    └─────────┘    └─────────┘    └─────────┘        │ │  │
│   │  │                                                                     │ │  │
│   │  └────────────────────────────────────────────────────────────────────┘ │  │
│   │                                                                          │  │
│   │  FEATURES:                                                               │  │
│   │  • Data sharding across nodes       • Pub/Sub messaging                 │  │
│   │  • Automatic failover               • Lua scripting                     │  │
│   │  • Hash slot distribution           • Transactions                      │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Dependency

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-redis</artifactId>
</dependency>
```

#### Configuration

```java
@Configuration
@EnableCaching
public class RedisCacheConfig {
    
    @Bean
    public RedisConnectionFactory redisConnectionFactory() {
        RedisStandaloneConfiguration config = new RedisStandaloneConfiguration();
        config.setHostName("localhost");
        config.setPort(6379);
        config.setPassword(RedisPassword.of("password"));
        config.setDatabase(0);
        return new LettuceConnectionFactory(config);
    }
    
    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        RedisCacheConfiguration defaultConfig = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofHours(1))
            .disableCachingNullValues()
            .serializeKeysWith(RedisSerializationContext.SerializationPair
                .fromSerializer(new StringRedisSerializer()))
            .serializeValuesWith(RedisSerializationContext.SerializationPair
                .fromSerializer(new GenericJackson2JsonRedisSerializer()));
        
        // Custom configurations per cache
        Map<String, RedisCacheConfiguration> cacheConfigs = new HashMap<>();
        cacheConfigs.put("products", defaultConfig.entryTtl(Duration.ofHours(24)));
        cacheConfigs.put("users", defaultConfig.entryTtl(Duration.ofMinutes(30)));
        cacheConfigs.put("sessions", defaultConfig.entryTtl(Duration.ofMinutes(15)));
        cacheConfigs.put("shortLived", defaultConfig.entryTtl(Duration.ofMinutes(5)));
        
        return RedisCacheManager.builder(factory)
            .cacheDefaults(defaultConfig)
            .withInitialCacheConfigurations(cacheConfigs)
            .transactionAware()
            .build();
    }
    
    @Bean
    public RedisTemplate<String, Object> redisTemplate(RedisConnectionFactory factory) {
        RedisTemplate<String, Object> template = new RedisTemplate<>();
        template.setConnectionFactory(factory);
        template.setKeySerializer(new StringRedisSerializer());
        template.setValueSerializer(new GenericJackson2JsonRedisSerializer());
        template.setHashKeySerializer(new StringRedisSerializer());
        template.setHashValueSerializer(new GenericJackson2JsonRedisSerializer());
        return template;
    }
}
```

```yaml
# application.yml
spring:
  cache:
    type: redis
    redis:
      time-to-live: 3600000  # 1 hour
      cache-null-values: false
      use-key-prefix: true
      key-prefix: "myapp:"
  
  redis:
    host: localhost
    port: 6379
    password: ${REDIS_PASSWORD:}
    timeout: 2000ms
    lettuce:
      pool:
        max-active: 8
        max-idle: 8
        min-idle: 2
        max-wait: -1ms
```

```yaml
# Redis Cluster Configuration
spring:
  redis:
    cluster:
      nodes:
        - redis-node-1:6379
        - redis-node-2:6379
        - redis-node-3:6379
      max-redirects: 3
    lettuce:
      cluster:
        refresh:
          adaptive: true
          period: 30s
```

#### Best Use Cases
- Distributed caching across microservices
- Session management
- Rate limiting
- Real-time leaderboards
- Pub/Sub messaging
- High-availability caching

#### Pros & Cons

| Pros | Cons |
|------|------|
| Distributed and shared | Network latency overhead |
| High availability (Redis Sentinel/Cluster) | Additional infrastructure |
| Rich data structures | Serialization overhead |
| Persistence options (RDB, AOF) | Cost (managed services) |
| Pub/Sub support | Memory-bound |
| Lua scripting | |
| Active community | |

#### When to Choose
✅ Microservices architecture  
✅ Need for distributed, shared cache  
✅ High availability requirements  
✅ Session sharing across instances  
✅ Cloud-native applications  
❌ Simple single-instance applications  
❌ Ultra-low latency requirements (use local cache + Redis L2)  

---

### 4.5 Hazelcast

#### Overview
In-memory data grid that provides distributed caching, computing, and messaging. Known for its ease of clustering and auto-discovery features.

#### Architecture Type
**Distributed** - Embedded or Client-Server mode

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    HAZELCAST ARCHITECTURE                                        │
│                                                                                  │
│   EMBEDDED MODE (Peer-to-Peer):                                                 │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│   │  │                    Hazelcast Cluster                             │   │  │
│   │  │                                                                  │   │  │
│   │  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐   │   │  │
│   │  │  │   App + HZ     │  │   App + HZ     │  │   App + HZ     │   │   │  │
│   │  │  │   Node 1       │──│   Node 2       │──│   Node 3       │   │   │  │
│   │  │  │                │  │                │  │                │   │   │  │
│   │  │  │  [Partition 1] │  │  [Partition 2] │  │  [Partition 3] │   │   │  │
│   │  │  │  [Backup of 3] │  │  [Backup of 1] │  │  [Backup of 2] │   │   │  │
│   │  │  └────────────────┘  └────────────────┘  └────────────────┘   │   │  │
│   │  │                                                                  │   │  │
│   │  │  • Data partitioned across nodes                                │   │  │
│   │  │  • Automatic backup distribution                                │   │  │
│   │  │  • Auto-discovery (multicast/TCP)                              │   │  │
│   │  └─────────────────────────────────────────────────────────────────┘   │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│   CLIENT-SERVER MODE:                                                            │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │  ┌────────────┐    ┌────────────┐    ┌────────────┐                    │  │
│   │  │ App Client │    │ App Client │    │ App Client │                    │  │
│   │  └─────┬──────┘    └─────┬──────┘    └─────┬──────┘                    │  │
│   │        │                 │                 │                            │  │
│   │        └─────────────────┼─────────────────┘                            │  │
│   │                          │                                              │  │
│   │  ┌───────────────────────▼─────────────────────────┐                   │  │
│   │  │              Hazelcast Cluster                   │                   │  │
│   │  │   ┌─────────┐    ┌─────────┐    ┌─────────┐    │                   │  │
│   │  │   │ Member 1│    │ Member 2│    │ Member 3│    │                   │  │
│   │  │   └─────────┘    └─────────┘    └─────────┘    │                   │  │
│   │  └─────────────────────────────────────────────────┘                   │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Dependency

```xml
<dependency>
    <groupId>com.hazelcast</groupId>
    <artifactId>hazelcast-spring</artifactId>
</dependency>
```

#### Configuration

```java
@Configuration
@EnableCaching
public class HazelcastCacheConfig {
    
    @Bean
    public Config hazelcastConfig() {
        Config config = new Config();
        config.setInstanceName("hazelcast-cache");
        
        // Network configuration
        config.getNetworkConfig()
            .setPort(5701)
            .setPortAutoIncrement(true)
            .getJoin()
            .getMulticastConfig()
            .setEnabled(true);
        
        // Or use TCP/IP discovery
        // config.getNetworkConfig()
        //     .getJoin()
        //     .getTcpIpConfig()
        //     .setEnabled(true)
        //     .addMember("192.168.1.1")
        //     .addMember("192.168.1.2");
        
        // Cache configurations
        config.addMapConfig(new MapConfig()
            .setName("products")
            .setMaxIdleSeconds(3600)
            .setTimeToLiveSeconds(86400)
            .setEvictionConfig(new EvictionConfig()
                .setMaxSizePolicy(MaxSizePolicy.PER_NODE)
                .setSize(10000)
                .setEvictionPolicy(EvictionPolicy.LRU))
            .setBackupCount(1));
        
        config.addMapConfig(new MapConfig()
            .setName("users")
            .setMaxIdleSeconds(1800)
            .setTimeToLiveSeconds(3600)
            .setEvictionConfig(new EvictionConfig()
                .setMaxSizePolicy(MaxSizePolicy.PER_NODE)
                .setSize(5000)
                .setEvictionPolicy(EvictionPolicy.LFU)));
        
        return config;
    }
    
    @Bean
    public HazelcastInstance hazelcastInstance(Config config) {
        return Hazelcast.newHazelcastInstance(config);
    }
    
    @Bean
    public CacheManager cacheManager(HazelcastInstance hazelcastInstance) {
        return new HazelcastCacheManager(hazelcastInstance);
    }
}
```

```yaml
# application.yml
spring:
  cache:
    type: hazelcast
  hazelcast:
    config: classpath:hazelcast.xml
```

```xml
<!-- hazelcast.xml -->
<hazelcast xmlns="http://www.hazelcast.com/schema/config">
    <instance-name>hazelcast-cache</instance-name>
    
    <network>
        <port auto-increment="true">5701</port>
        <join>
            <multicast enabled="false"/>
            <tcp-ip enabled="true">
                <member>192.168.1.1</member>
                <member>192.168.1.2</member>
            </tcp-ip>
        </join>
    </network>
    
    <map name="products">
        <time-to-live-seconds>86400</time-to-live-seconds>
        <max-idle-seconds>3600</max-idle-seconds>
        <eviction size="10000" max-size-policy="PER_NODE" eviction-policy="LRU"/>
        <backup-count>1</backup-count>
    </map>
    
    <map name="users">
        <time-to-live-seconds>3600</time-to-live-seconds>
        <eviction size="5000" max-size-policy="PER_NODE" eviction-policy="LFU"/>
    </map>
</hazelcast>
```

#### Best Use Cases
- Distributed caching with auto-clustering
- In-memory data grid scenarios
- Distributed computing
- Event streaming
- When Kubernetes auto-discovery is needed

#### Pros & Cons

| Pros | Cons |
|------|------|
| Easy clustering (auto-discovery) | Higher memory usage |
| Embedded or client-server modes | Complex for simple use cases |
| Rich distributed data structures | JVM-based (no native persistence) |
| Near-cache support | |
| Kubernetes discovery built-in | |
| No single point of failure | |

#### When to Choose
✅ Need for auto-clustering  
✅ Distributed computing beyond caching  
✅ Kubernetes deployments  
✅ Peer-to-peer data sharing  
❌ Simple caching needs  
❌ When Redis ecosystem features needed  

---

### 4.6 Infinispan

#### Overview
Highly scalable, distributed key-value data store and cache. Originally developed by Red Hat, it's the default cache for WildFly/JBoss.

#### Architecture Type
**Distributed** - Embedded or Client-Server mode

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    INFINISPAN ARCHITECTURE                                       │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                    Infinispan Cluster                                    │  │
│   │                                                                          │  │
│   │    CACHE MODES:                                                          │  │
│   │    ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐       │  │
│   │    │    LOCAL        │  │  REPLICATED     │  │  DISTRIBUTED    │       │  │
│   │    │                 │  │                 │  │                 │       │  │
│   │    │  [Node 1]       │  │ [N1]==[N2]==[N3]│  │ [N1]──[N2]──[N3]│       │  │
│   │    │  Single node    │  │ Full copy each  │  │ Partitioned     │       │  │
│   │    │  No clustering  │  │ Small datasets  │  │ Large datasets  │       │  │
│   │    └─────────────────┘  └─────────────────┘  └─────────────────┘       │  │
│   │                                                                          │  │
│   │    FEATURES:                                                             │  │
│   │    • JCache (JSR-107) compliant                                         │  │
│   │    • Hibernate L2 cache support                                         │  │
│   │    • Cross-site replication                                              │  │
│   │    • REST/Hot Rod/Memcached protocols                                   │  │
│   │    • Kubernetes operator                                                │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Dependency

```xml
<dependency>
    <groupId>org.infinispan</groupId>
    <artifactId>infinispan-spring-boot-starter-embedded</artifactId>
    <version>14.0.0.Final</version>
</dependency>
```

#### Configuration

```java
@Configuration
@EnableCaching
public class InfinispanCacheConfig {
    
    @Bean
    public EmbeddedCacheManager cacheManager() {
        GlobalConfigurationBuilder global = GlobalConfigurationBuilder.defaultClusteredBuilder();
        global.transport().clusterName("my-cluster");
        
        ConfigurationBuilder config = new ConfigurationBuilder();
        config.clustering()
            .cacheMode(CacheMode.DIST_SYNC)
            .hash().numOwners(2);
        
        config.expiration()
            .lifespan(1, TimeUnit.HOURS)
            .maxIdle(30, TimeUnit.MINUTES);
        
        config.memory()
            .maxCount(10000);
        
        EmbeddedCacheManager manager = new DefaultCacheManager(
            global.build(), 
            config.build()
        );
        
        // Define specific caches
        manager.defineConfiguration("products", config.build());
        manager.defineConfiguration("users", config.build());
        
        return manager;
    }
    
    @Bean
    public SpringEmbeddedCacheManager springCacheManager(EmbeddedCacheManager cacheManager) {
        return new SpringEmbeddedCacheManager(cacheManager);
    }
}
```

```yaml
# application.yml
infinispan:
  embedded:
    cluster-name: my-cluster
    config-xml: infinispan.xml
```

#### Best Use Cases
- JBoss/WildFly applications
- Cross-data-center replication
- When JCache compliance is required
- Hibernate L2 cache in clustered environments

#### Pros & Cons

| Pros | Cons |
|------|------|
| JCache compliant | Smaller community than Redis |
| Multiple cache modes | Complex configuration |
| Cross-site replication | JVM-centric |
| Hibernate integration | |
| Kubernetes operator | |
| Transactions support | |

#### When to Choose
✅ JBoss/WildFly applications  
✅ Need for cross-data-center replication  
✅ Hibernate L2 cache (clustered)  
✅ JCache (JSR-107) requirement  
❌ Simple caching needs  
❌ Polyglot environments (use Redis)  

---

## 5. When to Choose Which Cache Provider

### 5.1 Comparison Table

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                              CACHE PROVIDER COMPARISON                                                       │
├────────────────┬────────────────┬────────────────┬────────────────┬────────────────┬────────────────────────┤
│   Feature      │  ConcurrentMap │   Caffeine     │    Ehcache     │     Redis      │    Hazelcast           │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Architecture   │ Local          │ Local          │ Local/Dist     │ Distributed    │ Distributed            │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Performance    │ ★★★★★          │ ★★★★★          │ ★★★★☆          │ ★★★★☆          │ ★★★★☆                  │
│ (Local ops)    │ Fastest        │ Near-optimal   │ Very fast      │ Network latency│ Network latency        │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Distributed    │ ✗ No           │ ✗ No           │ ✓ (Terracotta) │ ✓ Yes          │ ✓ Yes                  │
│ Support        │                │                │                │                │                        │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Persistence    │ ✗ No           │ ✗ No           │ ✓ Yes (Disk)   │ ✓ Yes          │ ✓ Yes                  │
│                │                │                │                │ (RDB/AOF)      │ (External stores)      │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ TTL Support    │ ✗ No           │ ✓ Yes          │ ✓ Yes          │ ✓ Yes          │ ✓ Yes                  │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Eviction       │ ✗ None         │ ✓ LRU/LFU/Size │ ✓ LRU/LFU/FIFO │ ✓ LRU/LFU/     │ ✓ LRU/LFU/Random       │
│ Policies       │                │   Window TinyLFU│   TTL/TTI     │   Volatile     │   TTL/TTI              │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Scalability    │ Single JVM     │ Single JVM     │ Multi-tier/    │ Cluster/       │ Cluster/               │
│                │                │                │ Cluster        │ Sentinel       │ Auto-partition         │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Configuration  │ ★★★★★          │ ★★★★☆          │ ★★★☆☆          │ ★★★★☆          │ ★★★☆☆                  │
│ Simplicity     │ Simplest       │ Simple         │ Moderate       │ Moderate       │ Moderate               │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Memory         │ Heap only      │ Heap           │ Heap/Off-heap/ │ External       │ Heap/Off-heap          │
│                │                │                │ Disk           │                │                        │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Dependencies   │ None           │ Single JAR     │ Multiple JARs  │ External       │ Multiple JARs          │
│                │ (built-in)     │                │                │ server         │                        │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Monitoring     │ ✗ None         │ ✓ Stats API    │ ✓ JMX/Stats    │ ✓ Redis CLI/   │ ✓ Management Center    │
│                │                │                │                │ Monitoring     │                        │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Best For       │ Dev/Test       │ Single-node    │ Enterprise/    │ Microservices/ │ Auto-clustering/       │
│                │ Simple apps    │ High-perf      │ Tiered storage │ Distributed    │ Data grid              │
├────────────────┼────────────────┼────────────────┼────────────────┼────────────────┼────────────────────────┤
│ Production     │ ★★☆☆☆          │ ★★★★★          │ ★★★★★          │ ★★★★★          │ ★★★★☆                  │
│ Readiness      │ Not recommended│ Excellent      │ Excellent      │ Excellent      │ Good                   │
└────────────────┴────────────────┴────────────────┴────────────────┴────────────────┴────────────────────────┘
```

### 5.2 Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE PROVIDER DECISION MATRIX                                │
│                                                                                  │
│   QUESTION FLOW:                                                                 │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                                                                          │  │
│   │                     Need distributed cache?                              │  │
│   │                              │                                           │  │
│   │              ┌───────────────┼───────────────┐                          │  │
│   │              │               │               │                          │  │
│   │             NO              YES             BOTH                        │  │
│   │              │               │            (Hybrid)                      │  │
│   │              ▼               ▼               │                          │  │
│   │         ┌────────┐     ┌────────────┐       ▼                          │  │
│   │         │ Single │     │ Multi-node │   L1: Caffeine                   │  │
│   │         │ Node   │     │ deployment │   L2: Redis                      │  │
│   │         └───┬────┘     └─────┬──────┘                                  │  │
│   │             │                │                                          │  │
│   │             ▼                ▼                                          │  │
│   │      ┌──────────────┐  ┌─────────────────┐                             │  │
│   │      │ CAFFEINE     │  │ Cloud-native?   │                             │  │
│   │      │ (Best choice)│  └────────┬────────┘                             │  │
│   │      └──────────────┘           │                                       │  │
│   │                     ┌───────────┼───────────┐                          │  │
│   │                    YES          │          NO                          │  │
│   │                     │           │           │                          │  │
│   │                     ▼           │           ▼                          │  │
│   │               ┌─────────┐      │      ┌──────────┐                    │  │
│   │               │  REDIS  │      │      │ Auto-    │                    │  │
│   │               │         │      │      │ cluster? │                    │  │
│   │               └─────────┘      │      └────┬─────┘                    │  │
│   │                                │           │                          │  │
│   │                                │     ┌─────┼─────┐                    │  │
│   │                                │    YES    │    NO                    │  │
│   │                                │     │     │     │                    │  │
│   │                                │     ▼     │     ▼                    │  │
│   │                                │ HAZELCAST │  REDIS                   │  │
│   │                                │           │                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Decision Guidelines by Scenario

#### Scenario 1: Small Application with Single Instance

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  SCENARIO: Small single-instance application                                     │
│                                                                                  │
│  CHARACTERISTICS:                                                                │
│  • Single server deployment                                                      │
│  • Low to moderate traffic                                                       │
│  • Simple caching needs                                                          │
│  • Budget constraints                                                            │
│                                                                                  │
│  RECOMMENDATION: CAFFEINE                                                        │
│                                                                                  │
│  WHY:                                                                            │
│  ✓ No external dependencies or infrastructure                                   │
│  ✓ High performance with near-optimal hit ratio                                 │
│  ✓ Simple configuration                                                         │
│  ✓ Built-in TTL and size-based eviction                                         │
│  ✓ Statistics for monitoring                                                    │
│                                                                                  │
│  CONFIGURATION:                                                                  │
│  spring:                                                                         │
│    cache:                                                                        │
│      type: caffeine                                                              │
│      caffeine:                                                                   │
│        spec: maximumSize=10000,expireAfterWrite=3600s                           │
│                                                                                  │
│  ALTERNATIVE: ConcurrentMap (for development/testing only)                      │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Scenario 2: Microservices Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  SCENARIO: Microservices with multiple instances                                 │
│                                                                                  │
│  CHARACTERISTICS:                                                                │
│  • Multiple service instances                                                    │
│  • Horizontal scaling                                                            │
│  • Shared state requirements                                                     │
│  • Cloud deployment (AWS, GCP, Azure)                                           │
│                                                                                  │
│  RECOMMENDATION: REDIS                                                           │
│                                                                                  │
│  WHY:                                                                            │
│  ✓ Distributed cache shared across all instances                                │
│  ✓ Managed services available (AWS ElastiCache, Azure Cache)                   │
│  ✓ Session sharing                                                              │
│  ✓ Rich feature set (pub/sub, data structures)                                 │
│  ✓ High availability with Sentinel/Cluster                                     │
│                                                                                  │
│  CONFIGURATION:                                                                  │
│  spring:                                                                         │
│    cache:                                                                        │
│      type: redis                                                                 │
│    redis:                                                                        │
│      cluster:                                                                    │
│        nodes: redis-1:6379,redis-2:6379,redis-3:6379                           │
│                                                                                  │
│  ALTERNATIVE: Hazelcast (if auto-discovery is preferred)                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Scenario 3: High Throughput System

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  SCENARIO: High throughput, low latency requirements                             │
│                                                                                  │
│  CHARACTERISTICS:                                                                │
│  • Millions of requests per second                                              │
│  • Sub-millisecond latency requirements                                         │
│  • Heavy read operations                                                         │
│  • Cost of network round-trip is significant                                    │
│                                                                                  │
│  RECOMMENDATION: HYBRID (L1 + L2)                                               │
│                    Caffeine (L1) + Redis (L2)                                   │
│                                                                                  │
│  WHY:                                                                            │
│  ✓ L1 (Caffeine): Ultra-fast local cache, sub-ms latency                       │
│  ✓ L2 (Redis): Distributed cache for consistency                               │
│  ✓ Reduces Redis network calls by 80-90%                                       │
│  ✓ Best of both worlds                                                          │
│                                                                                  │
│  ARCHITECTURE:                                                                   │
│  [Request] → [L1 Caffeine] → miss → [L2 Redis] → miss → [Database]            │
│                   ↓ hit                 ↓ hit                                   │
│               [Return]            [Return + Cache in L1]                        │
│                                                                                  │
│  L1 SETTINGS:                                                                    │
│  • Small size (1000-5000 entries)                                               │
│  • Short TTL (1-5 minutes)                                                       │
│                                                                                  │
│  L2 SETTINGS:                                                                    │
│  • Larger size                                                                   │
│  • Longer TTL (30 min - 24 hours)                                               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Hybrid L1 + L2 Cache Implementation
@Configuration
@EnableCaching
public class HybridCacheConfig {
    
    @Bean
    @Primary
    public CacheManager cacheManager(RedisConnectionFactory redisFactory) {
        // L1: Caffeine
        CaffeineCacheManager l1CacheManager = new CaffeineCacheManager();
        l1CacheManager.setCaffeine(Caffeine.newBuilder()
            .maximumSize(5000)
            .expireAfterWrite(Duration.ofMinutes(5))
            .recordStats());
        l1CacheManager.setCacheNames(List.of("products", "users"));
        
        // L2: Redis
        RedisCacheManager l2CacheManager = RedisCacheManager.builder(redisFactory)
            .cacheDefaults(RedisCacheConfiguration.defaultCacheConfig()
                .entryTtl(Duration.ofHours(1)))
            .build();
        
        return new TieredCacheManager(l1CacheManager, l2CacheManager);
    }
}

// Custom Tiered Cache Manager
public class TieredCacheManager implements CacheManager {
    
    private final CacheManager l1CacheManager;
    private final CacheManager l2CacheManager;
    
    public TieredCacheManager(CacheManager l1, CacheManager l2) {
        this.l1CacheManager = l1;
        this.l2CacheManager = l2;
    }
    
    @Override
    public Cache getCache(String name) {
        Cache l1Cache = l1CacheManager.getCache(name);
        Cache l2Cache = l2CacheManager.getCache(name);
        return new TieredCache(name, l1Cache, l2Cache);
    }
    
    @Override
    public Collection<String> getCacheNames() {
        Set<String> names = new HashSet<>();
        names.addAll(l1CacheManager.getCacheNames());
        names.addAll(l2CacheManager.getCacheNames());
        return names;
    }
}

// Tiered Cache Implementation
public class TieredCache implements Cache {
    
    private final String name;
    private final Cache l1Cache;
    private final Cache l2Cache;
    
    @Override
    public ValueWrapper get(Object key) {
        // Try L1 first
        ValueWrapper l1Value = l1Cache.get(key);
        if (l1Value != null) {
            return l1Value;  // L1 hit
        }
        
        // Try L2
        ValueWrapper l2Value = l2Cache.get(key);
        if (l2Value != null) {
            // Populate L1 from L2
            l1Cache.put(key, l2Value.get());
            return l2Value;  // L2 hit
        }
        
        return null;  // Miss
    }
    
    @Override
    public void put(Object key, Object value) {
        l1Cache.put(key, value);
        l2Cache.put(key, value);
    }
    
    @Override
    public void evict(Object key) {
        l1Cache.evict(key);
        l2Cache.evict(key);
    }
    
    @Override
    public void clear() {
        l1Cache.clear();
        l2Cache.clear();
    }
    
    // ... other methods
}
```

#### Scenario 4: Need for Data Persistence

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  SCENARIO: Cache data must survive restarts                                      │
│                                                                                  │
│  CHARACTERISTICS:                                                                │
│  • Expensive to rebuild cache                                                    │
│  • Long cache warm-up time is unacceptable                                      │
│  • Regulatory/compliance requirements                                            │
│                                                                                  │
│  RECOMMENDATIONS (in order of preference):                                       │
│                                                                                  │
│  1. REDIS with persistence                                                       │
│     • RDB snapshots for periodic backup                                         │
│     • AOF for durability (every write logged)                                   │
│     • Best for distributed environments                                          │
│                                                                                  │
│  2. EHCACHE with disk tier                                                       │
│     • Local disk persistence                                                     │
│     • Survives JVM restarts                                                      │
│     • Good for single-node applications                                          │
│                                                                                  │
│  REDIS PERSISTENCE CONFIG:                                                       │
│  redis.conf:                                                                     │
│    # RDB (point-in-time snapshots)                                              │
│    save 900 1      # Save after 900 sec if at least 1 key changed              │
│    save 300 10     # Save after 300 sec if at least 10 keys changed            │
│                                                                                  │
│    # AOF (append-only file for durability)                                      │
│    appendonly yes                                                                │
│    appendfsync everysec                                                          │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Scenario 5: Cloud-Native Deployments

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  SCENARIO: Kubernetes/Cloud deployment                                           │
│                                                                                  │
│  CHARACTERISTICS:                                                                │
│  • Dynamic scaling (pods come and go)                                           │
│  • Service discovery                                                             │
│  • Ephemeral containers                                                          │
│  • Infrastructure as code                                                        │
│                                                                                  │
│  RECOMMENDATIONS:                                                                │
│                                                                                  │
│  1. MANAGED REDIS (Preferred)                                                    │
│     • AWS ElastiCache                                                            │
│     • Azure Cache for Redis                                                      │
│     • Google Cloud Memorystore                                                   │
│     • Redis Cloud                                                                │
│                                                                                  │
│  2. HAZELCAST (If auto-clustering needed)                                        │
│     • Kubernetes discovery plugin                                                │
│     • Automatic cluster formation                                                │
│     • No external dependencies                                                   │
│                                                                                  │
│  HAZELCAST K8s CONFIG:                                                           │
│  hazelcast:                                                                      │
│    network:                                                                      │
│      join:                                                                       │
│        kubernetes:                                                               │
│          enabled: true                                                           │
│          namespace: default                                                      │
│          service-name: hazelcast-service                                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 5.4 Quick Decision Guide

| Your Situation | Recommended Provider | Reason |
|----------------|---------------------|--------|
| Development/Testing | ConcurrentMap or Caffeine | Zero setup, fast iteration |
| Single-instance production | **Caffeine** | Best performance, no infrastructure |
| Multi-instance, cloud | **Redis** | Distributed, managed options |
| Auto-clustering needed | Hazelcast | Built-in discovery |
| Hibernate L2 cache | Ehcache or Infinispan | Native support |
| Very large datasets | Ehcache (tiered) | Heap + off-heap + disk |
| Maximum throughput | Caffeine + Redis | Hybrid L1/L2 approach |
| JBoss/WildFly | Infinispan | Default, well-integrated |

---

## 6. Advanced Topics

### 6.1 Cache Invalidation Challenges

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE INVALIDATION CHALLENGES                                 │
│                                                                                  │
│  "There are only two hard things in Computer Science:                           │
│   cache invalidation and naming things." — Phil Karlton                         │
│                                                                                  │
│  CHALLENGE 1: Stale Data                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Problem: Cache contains outdated data after DB update                   │   │
│  │                                                                          │   │
│  │  Timeline:                                                               │   │
│  │  T1: [Cache: Product A = $100] [DB: Product A = $100]                   │   │
│  │  T2: [DB Update: Product A = $150]                                      │   │
│  │  T3: [Cache: Product A = $100] ← STALE! [DB: Product A = $150]         │   │
│  │                                                                          │   │
│  │  Solutions:                                                              │   │
│  │  • Short TTL (accept some staleness)                                    │   │
│  │  • Event-driven invalidation (publish on update)                        │   │
│  │  • Write-through pattern                                                │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  CHALLENGE 2: Distributed Cache Invalidation                                     │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Problem: Multiple cache instances holding different versions            │   │
│  │                                                                          │   │
│  │  [Instance 1 Cache: v1] [Instance 2 Cache: v2] [Instance 3 Cache: v1]  │   │
│  │                                                                          │   │
│  │  Solutions:                                                              │   │
│  │  • Centralized cache (Redis)                                            │   │
│  │  • Pub/Sub invalidation notifications                                   │   │
│  │  • Cache versioning with namespace                                      │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  CHALLENGE 3: Cascading Invalidation                                             │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Problem: Updating one entity requires invalidating related caches       │   │
│  │                                                                          │   │
│  │  Update Product → Invalidate:                                           │   │
│  │    • Product cache                                                       │   │
│  │    • Category listing cache                                             │   │
│  │    • Search results cache                                               │   │
│  │    • Recommendation cache                                               │   │
│  │    • Related products cache                                             │   │
│  │                                                                          │   │
│  │  Solutions:                                                              │   │
│  │  • Careful cache dependency mapping                                     │   │
│  │  • Event-driven architecture                                            │   │
│  │  • Accept eventual consistency                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Solution: Event-Driven Cache Invalidation
@Service
@RequiredArgsConstructor
public class ProductService {
    
    private final ProductRepository repository;
    private final ApplicationEventPublisher eventPublisher;
    
    @Transactional
    @CacheEvict(value = "products", key = "#id")
    public Product updateProduct(Long id, ProductRequest request) {
        Product product = repository.findById(id).orElseThrow();
        updateEntity(product, request);
        Product saved = repository.save(product);
        
        // Publish event for related cache invalidation
        eventPublisher.publishEvent(new ProductUpdatedEvent(
            saved.getId(),
            saved.getCategory(),
            saved.getSku()
        ));
        
        return saved;
    }
}

@Component
@Slf4j
public class CacheInvalidationHandler {
    
    private final CacheManager cacheManager;
    
    @EventListener
    @Async
    public void handleProductUpdate(ProductUpdatedEvent event) {
        log.info("Invalidating related caches for product: {}", event.getProductId());
        
        // Invalidate category cache
        Cache categoryCache = cacheManager.getCache("productsByCategory");
        if (categoryCache != null) {
            categoryCache.evict(event.getCategory());
        }
        
        // Invalidate search cache (all entries - expensive but ensures consistency)
        Cache searchCache = cacheManager.getCache("searchResults");
        if (searchCache != null) {
            searchCache.clear();
        }
        
        // Invalidate SKU cache
        Cache skuCache = cacheManager.getCache("productsBySku");
        if (skuCache != null) {
            skuCache.evict(event.getSku());
        }
    }
}
```

### 6.2 Cache Stampede Problem

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE STAMPEDE (Thundering Herd)                              │
│                                                                                  │
│  PROBLEM:                                                                        │
│  When a popular cache entry expires, many concurrent requests hit the database  │
│  simultaneously, potentially causing database overload.                          │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  Before expiry:      [Request] → [Cache HIT] → [Return]                 │   │
│  │                      [Request] → [Cache HIT] → [Return]                 │   │
│  │                      [Request] → [Cache HIT] → [Return]                 │   │
│  │                                                                          │   │
│  │  At expiry (T=0):    [Request 1] → [Cache MISS] → [DB Query]           │   │
│  │                      [Request 2] → [Cache MISS] → [DB Query]            │   │
│  │                      [Request 3] → [Cache MISS] → [DB Query]            │   │
│  │                      ...                                                 │   │
│  │                      [Request N] → [Cache MISS] → [DB Query]            │   │
│  │                                                                          │   │
│  │                      ↓                                                   │   │
│  │                      DATABASE OVERLOAD! 💥                              │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  SOLUTIONS:                                                                      │
│                                                                                  │
│  1. MUTEX/LOCK: Only one thread fetches from DB, others wait                    │
│  2. EARLY REFRESH: Refresh before expiry (probabilistic or scheduled)          │
│  3. STALE-WHILE-REVALIDATE: Return stale data while fetching fresh             │
│  4. RANDOMIZED TTL: Add jitter to prevent synchronized expiry                   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Solution 1: Using @Cacheable with sync=true (Mutex)
@Service
public class ProductService {
    
    // sync=true ensures only one thread computes the value
    // Others wait for the result
    @Cacheable(value = "products", key = "#id", sync = true)
    public Product getProduct(Long id) {
        log.info("Cache miss - fetching from DB: {}", id);
        return repository.findById(id).orElseThrow();
    }
}

// Solution 2: Manual Locking with Redis
@Service
@RequiredArgsConstructor
public class ProductServiceWithLock {
    
    private final ProductRepository repository;
    private final RedisTemplate<String, Object> redisTemplate;
    private final CacheManager cacheManager;
    
    public Product getProductWithLock(Long id) {
        String cacheKey = "products:" + id;
        String lockKey = "lock:" + cacheKey;
        
        // Try cache first
        Cache cache = cacheManager.getCache("products");
        Product cached = cache.get(id, Product.class);
        if (cached != null) {
            return cached;
        }
        
        // Acquire lock
        Boolean acquired = redisTemplate.opsForValue()
            .setIfAbsent(lockKey, "locked", Duration.ofSeconds(10));
        
        if (Boolean.TRUE.equals(acquired)) {
            try {
                // Double-check cache after acquiring lock
                cached = cache.get(id, Product.class);
                if (cached != null) {
                    return cached;
                }
                
                // Fetch from DB and cache
                Product product = repository.findById(id).orElseThrow();
                cache.put(id, product);
                return product;
            } finally {
                redisTemplate.delete(lockKey);
            }
        } else {
            // Wait and retry
            try {
                Thread.sleep(100);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            return getProductWithLock(id);  // Retry
        }
    }
}

// Solution 3: Probabilistic Early Refresh
@Service
public class ProductServiceWithEarlyRefresh {
    
    private final ProductRepository repository;
    private final CacheManager cacheManager;
    private static final double EARLY_REFRESH_PROBABILITY = 0.1;  // 10%
    
    @Cacheable(value = "products", key = "#id")
    public Product getProduct(Long id) {
        return repository.findById(id).orElseThrow();
    }
    
    // Call this from a wrapper method
    public Product getProductWithEarlyRefresh(Long id) {
        Cache cache = cacheManager.getCache("products");
        Cache.ValueWrapper wrapper = cache.get(id);
        
        if (wrapper != null) {
            // Probabilistically refresh before expiry
            if (Math.random() < EARLY_REFRESH_PROBABILITY) {
                CompletableFuture.runAsync(() -> refreshCache(id));
            }
            return (Product) wrapper.get();
        }
        
        return getProduct(id);
    }
    
    @CachePut(value = "products", key = "#id")
    public Product refreshCache(Long id) {
        return repository.findById(id).orElseThrow();
    }
}

// Solution 4: Randomized TTL
@Configuration
public class RandomizedTtlCacheConfig {
    
    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        // Base TTL with randomization to prevent stampede
        RedisCacheConfiguration config = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(60));  // Base TTL
        
        return new RandomizedTtlRedisCacheManager(
            RedisCacheWriter.nonLockingRedisCacheWriter(factory),
            config,
            0.1  // 10% jitter
        );
    }
}

// Custom CacheManager with randomized TTL
public class RandomizedTtlRedisCacheManager extends RedisCacheManager {
    
    private final double jitterFactor;
    
    @Override
    protected RedisCache createRedisCache(String name, RedisCacheConfiguration config) {
        Duration baseTtl = config.getTtl();
        if (!baseTtl.isZero()) {
            // Add random jitter: TTL ± (jitterFactor * TTL)
            long jitter = (long) (baseTtl.toMillis() * jitterFactor * (Math.random() - 0.5) * 2);
            Duration randomizedTtl = baseTtl.plusMillis(jitter);
            config = config.entryTtl(randomizedTtl);
        }
        return super.createRedisCache(name, config);
    }
}
```

### 6.3 Cache Penetration

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE PENETRATION                                             │
│                                                                                  │
│  PROBLEM:                                                                        │
│  Requests for non-existent data always miss the cache and hit the database.    │
│  Malicious users can exploit this to attack the database.                        │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  [Request: id=-1] → [Cache MISS] → [DB: NULL] → Return null             │   │
│  │  [Request: id=-2] → [Cache MISS] → [DB: NULL] → Return null             │   │
│  │  [Request: id=-3] → [Cache MISS] → [DB: NULL] → Return null             │   │
│  │  ...                                                                     │   │
│  │  [Attacker floods with non-existent IDs]                                │   │
│  │                                                                          │   │
│  │  Result: Every request hits the database!                               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  SOLUTIONS:                                                                      │
│                                                                                  │
│  1. CACHE NULL VALUES: Store null/empty results with short TTL                  │
│  2. BLOOM FILTER: Check existence before querying                               │
│  3. INPUT VALIDATION: Reject obviously invalid requests                         │
│  4. RATE LIMITING: Limit requests per client                                    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Solution 1: Cache Null Values
@Service
public class ProductService {
    
    // Cache null values with short TTL
    @Cacheable(value = "products", key = "#id", unless = "false")  // Always cache
    public ProductDTO getProduct(Long id) {
        return repository.findById(id)
            .map(this::mapToDTO)
            .orElse(null);  // null will be cached
    }
}

// Redis configuration to handle null caching
@Bean
public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
    RedisCacheConfiguration defaultConfig = RedisCacheConfiguration.defaultCacheConfig()
        .entryTtl(Duration.ofHours(1))
        .serializeValuesWith(RedisSerializationContext.SerializationPair
            .fromSerializer(new GenericJackson2JsonRedisSerializer()));
    
    // Separate config for caches that should store nulls
    RedisCacheConfiguration nullableCacheConfig = RedisCacheConfiguration.defaultCacheConfig()
        .entryTtl(Duration.ofMinutes(5))  // Short TTL for null values
        // Don't call .disableCachingNullValues()
        .serializeValuesWith(RedisSerializationContext.SerializationPair
            .fromSerializer(new GenericJackson2JsonRedisSerializer()));
    
    Map<String, RedisCacheConfiguration> configs = new HashMap<>();
    configs.put("products", nullableCacheConfig);  // Allow null caching
    configs.put("users", defaultConfig.disableCachingNullValues());  // No nulls
    
    return RedisCacheManager.builder(factory)
        .cacheDefaults(defaultConfig)
        .withInitialCacheConfigurations(configs)
        .build();
}

// Solution 2: Bloom Filter
@Service
@RequiredArgsConstructor
public class ProductServiceWithBloomFilter {
    
    private final ProductRepository repository;
    private final BloomFilter<Long> productIdBloomFilter;  // Google Guava
    
    @PostConstruct
    public void initBloomFilter() {
        // Initialize bloom filter with existing product IDs
        repository.findAllIds().forEach(productIdBloomFilter::put);
    }
    
    @Cacheable(value = "products", key = "#id")
    public ProductDTO getProduct(Long id) {
        // Check bloom filter first
        if (!productIdBloomFilter.mightContain(id)) {
            // Definitely doesn't exist - return immediately
            return null;
        }
        
        // Might exist - check database
        return repository.findById(id)
            .map(this::mapToDTO)
            .orElse(null);
    }
    
    @CachePut(value = "products", key = "#result.id")
    public ProductDTO createProduct(ProductRequest request) {
        Product saved = repository.save(mapToEntity(request));
        productIdBloomFilter.put(saved.getId());  // Add to bloom filter
        return mapToDTO(saved);
    }
}

// Bloom Filter Configuration
@Configuration
public class BloomFilterConfig {
    
    @Bean
    public BloomFilter<Long> productIdBloomFilter() {
        return BloomFilter.create(
            Funnels.longFunnel(),
            1_000_000,  // Expected insertions
            0.01        // False positive probability (1%)
        );
    }
}

// Solution 3: Input Validation + Rate Limiting
@Service
@RequiredArgsConstructor
public class ProductServiceWithValidation {
    
    private final ProductRepository repository;
    private final RateLimiter rateLimiter;  // Resilience4j or custom
    
    @Cacheable(value = "products", key = "#id")
    public ProductDTO getProduct(Long id) {
        // Input validation
        if (id == null || id <= 0) {
            throw new IllegalArgumentException("Invalid product ID: " + id);
        }
        
        // Rate limiting
        if (!rateLimiter.acquirePermission()) {
            throw new RateLimitExceededException("Too many requests");
        }
        
        return repository.findById(id)
            .map(this::mapToDTO)
            .orElse(null);
    }
}
```

### 6.4 Cache Breakdown

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE BREAKDOWN                                               │
│                                                                                  │
│  PROBLEM:                                                                        │
│  A single "hot" key expires and causes massive concurrent database requests.    │
│  Different from stampede - this affects one specific popular key.               │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  Hot Key: "homepage_banner" (accessed 10,000 times/second)              │   │
│  │                                                                          │   │
│  │  Normal:    [10,000 requests] → [Cache HIT] → [Return]                  │   │
│  │                                                                          │   │
│  │  At expiry: [10,000 requests] → [Cache MISS] → [DB] 💥                  │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  SOLUTIONS:                                                                      │
│                                                                                  │
│  1. NEVER EXPIRE: Set no TTL for critical hot keys                             │
│  2. MUTEX: Only one thread refreshes, others get stale data                    │
│  3. EARLY REFRESH: Refresh hot keys before expiry                               │
│  4. LOGICAL EXPIRY: Store expiry time in value, refresh asynchronously         │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Solution: Logical Expiry with Async Refresh
@Service
@RequiredArgsConstructor
@Slf4j
public class HotKeyService {
    
    private final BannerRepository repository;
    private final RedisTemplate<String, CacheWrapper<Banner>> redisTemplate;
    private final ExecutorService refreshExecutor = Executors.newFixedThreadPool(5);
    
    public Banner getHomepageBanner() {
        String key = "homepage_banner";
        
        CacheWrapper<Banner> wrapper = redisTemplate.opsForValue().get(key);
        
        if (wrapper == null) {
            // Cold start - fetch and cache
            return fetchAndCache(key);
        }
        
        // Check logical expiry
        if (wrapper.isExpired()) {
            // Return stale data immediately
            Banner staleData = wrapper.getData();
            
            // Refresh asynchronously
            if (wrapper.tryLock()) {  // Prevent multiple refreshes
                refreshExecutor.submit(() -> {
                    try {
                        fetchAndCache(key);
                    } finally {
                        wrapper.unlock();
                    }
                });
            }
            
            return staleData;  // Return stale while refreshing
        }
        
        return wrapper.getData();
    }
    
    private Banner fetchAndCache(String key) {
        Banner banner = repository.findHomepageBanner();
        
        CacheWrapper<Banner> wrapper = new CacheWrapper<>(
            banner,
            Instant.now().plusSeconds(300),  // Logical expiry: 5 minutes
            new AtomicBoolean(false)
        );
        
        // Physical TTL much longer than logical expiry
        redisTemplate.opsForValue().set(key, wrapper, Duration.ofHours(24));
        
        return banner;
    }
}

@Data
@AllArgsConstructor
public class CacheWrapper<T> implements Serializable {
    private T data;
    private Instant logicalExpiry;
    private AtomicBoolean refreshing;
    
    public boolean isExpired() {
        return Instant.now().isAfter(logicalExpiry);
    }
    
    public boolean tryLock() {
        return refreshing.compareAndSet(false, true);
    }
    
    public void unlock() {
        refreshing.set(false);
    }
}
```

### 6.5 Monitoring and Metrics (Micrometer Integration)

```java
// Cache metrics configuration
@Configuration
@EnableCaching
public class CacheMetricsConfig {
    
    @Bean
    public CacheManager cacheManager(MeterRegistry meterRegistry) {
        CaffeineCacheManager cacheManager = new CaffeineCacheManager();
        cacheManager.setCaffeine(Caffeine.newBuilder()
            .maximumSize(10000)
            .expireAfterWrite(Duration.ofMinutes(10))
            .recordStats());  // Enable statistics
        
        cacheManager.setCacheNames(List.of("products", "users", "categories"));
        
        return cacheManager;
    }
    
    // Register cache metrics with Micrometer
    @Bean
    public CacheMetricsRegistrar cacheMetricsRegistrar(
            CacheManager cacheManager, 
            MeterRegistry meterRegistry) {
        return new CacheMetricsRegistrar(cacheManager, meterRegistry);
    }
}

@Component
@RequiredArgsConstructor
@Slf4j
public class CacheMetricsRegistrar implements ApplicationRunner {
    
    private final CacheManager cacheManager;
    private final MeterRegistry meterRegistry;
    
    @Override
    public void run(ApplicationArguments args) {
        for (String cacheName : cacheManager.getCacheNames()) {
            Cache cache = cacheManager.getCache(cacheName);
            
            if (cache instanceof CaffeineCache) {
                CaffeineCache caffeineCache = (CaffeineCache) cache;
                com.github.benmanes.caffeine.cache.Cache<Object, Object> nativeCache = 
                    caffeineCache.getNativeCache();
                
                // Register Caffeine metrics
                CaffeineCacheMetrics.monitor(meterRegistry, nativeCache, cacheName);
                
                log.info("Registered metrics for cache: {}", cacheName);
            }
        }
    }
}

// Custom cache metrics endpoint
@RestController
@RequestMapping("/api/cache")
@RequiredArgsConstructor
public class CacheMetricsController {
    
    private final CacheManager cacheManager;
    
    @GetMapping("/stats")
    public Map<String, CacheStats> getCacheStats() {
        Map<String, CacheStats> stats = new HashMap<>();
        
        for (String cacheName : cacheManager.getCacheNames()) {
            Cache cache = cacheManager.getCache(cacheName);
            
            if (cache instanceof CaffeineCache) {
                CaffeineCache caffeineCache = (CaffeineCache) cache;
                com.github.benmanes.caffeine.cache.stats.CacheStats nativeStats = 
                    caffeineCache.getNativeCache().stats();
                
                stats.put(cacheName, CacheStats.builder()
                    .hitCount(nativeStats.hitCount())
                    .missCount(nativeStats.missCount())
                    .hitRate(nativeStats.hitRate())
                    .evictionCount(nativeStats.evictionCount())
                    .averageLoadPenalty(nativeStats.averageLoadPenalty())
                    .build());
            }
        }
        
        return stats;
    }
    
    @DeleteMapping("/{cacheName}")
    public ResponseEntity<Void> clearCache(@PathVariable String cacheName) {
        Cache cache = cacheManager.getCache(cacheName);
        if (cache != null) {
            cache.clear();
            return ResponseEntity.ok().build();
        }
        return ResponseEntity.notFound().build();
    }
}

@Data
@Builder
class CacheStats {
    private long hitCount;
    private long missCount;
    private double hitRate;
    private long evictionCount;
    private double averageLoadPenalty;
}
```

```yaml
# application.yml - Actuator endpoints for cache metrics
management:
  endpoints:
    web:
      exposure:
        include: health, metrics, caches
  metrics:
    tags:
      application: ${spring.application.name}
    export:
      prometheus:
        enabled: true
```

### 6.6 Best Practices

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHING BEST PRACTICES                                        │
│                                                                                  │
│  1. CACHE DESIGN                                                                 │
│  ──────────────────────────────────────────────────────────────────────────     │
│  ✓ Cache at the right layer (service layer, not repository)                    │
│  ✓ Use meaningful, consistent cache names                                       │
│  ✓ Design cache keys to avoid collisions                                        │
│  ✓ Consider cache granularity (fine vs coarse)                                 │
│  ✓ Decide on consistency model (strong vs eventual)                            │
│                                                                                  │
│  2. TTL STRATEGY                                                                 │
│  ──────────────────────────────────────────────────────────────────────────     │
│  ✓ Set appropriate TTL based on data volatility                                │
│  ✓ Use different TTLs for different data types                                 │
│  ✓ Consider randomized TTL to prevent stampede                                 │
│  ✓ Balance freshness vs performance                                            │
│                                                                                  │
│  3. MEMORY MANAGEMENT                                                            │
│  ──────────────────────────────────────────────────────────────────────────     │
│  ✓ Set maximum cache size                                                       │
│  ✓ Monitor memory usage                                                         │
│  ✓ Use appropriate eviction policy (LRU/LFU)                                   │
│  ✓ Consider off-heap or distributed cache for large datasets                   │
│                                                                                  │
│  4. MONITORING & OBSERVABILITY                                                   │
│  ──────────────────────────────────────────────────────────────────────────     │
│  ✓ Track cache hit/miss rates                                                   │
│  ✓ Monitor eviction counts                                                      │
│  ✓ Set up alerts for abnormal patterns                                         │
│  ✓ Log cache operations at debug level                                         │
│                                                                                  │
│  5. ERROR HANDLING                                                               │
│  ──────────────────────────────────────────────────────────────────────────     │
│  ✓ Implement graceful degradation (fail to DB on cache error)                  │
│  ✓ Use circuit breaker for distributed cache                                   │
│  ✓ Don't let cache failures bring down the application                         │
│                                                                                  │
│  6. TESTING                                                                      │
│  ──────────────────────────────────────────────────────────────────────────     │
│  ✓ Test cache behavior (hit, miss, eviction)                                   │
│  ✓ Use ConcurrentMap cache for unit tests                                      │
│  ✓ Integration test with real cache provider                                   │
│  ✓ Test cache invalidation scenarios                                           │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Best Practice: Graceful Degradation
@Service
@RequiredArgsConstructor
@Slf4j
public class ResilientProductService {
    
    private final ProductRepository repository;
    private final CacheManager cacheManager;
    
    public Product getProduct(Long id) {
        try {
            return getProductFromCache(id);
        } catch (Exception e) {
            log.warn("Cache error, falling back to database: {}", e.getMessage());
            return repository.findById(id).orElseThrow();
        }
    }
    
    private Product getProductFromCache(Long id) {
        Cache cache = cacheManager.getCache("products");
        if (cache == null) {
            return repository.findById(id).orElseThrow();
        }
        
        Product cached = cache.get(id, Product.class);
        if (cached != null) {
            return cached;
        }
        
        Product product = repository.findById(id).orElseThrow();
        cache.put(id, product);
        return product;
    }
}

// Best Practice: Cache Error Handler
@Configuration
public class CacheErrorHandlerConfig extends CachingConfigurerSupport {
    
    @Override
    public CacheErrorHandler errorHandler() {
        return new CacheErrorHandler() {
            @Override
            public void handleCacheGetError(RuntimeException exception, Cache cache, Object key) {
                log.warn("Cache get error - cache: {}, key: {}, error: {}", 
                    cache.getName(), key, exception.getMessage());
                // Don't rethrow - allow fallback to DB
            }
            
            @Override
            public void handleCachePutError(RuntimeException exception, Cache cache, Object key, Object value) {
                log.warn("Cache put error - cache: {}, key: {}, error: {}", 
                    cache.getName(), key, exception.getMessage());
            }
            
            @Override
            public void handleCacheEvictError(RuntimeException exception, Cache cache, Object key) {
                log.warn("Cache evict error - cache: {}, key: {}, error: {}", 
                    cache.getName(), key, exception.getMessage());
            }
            
            @Override
            public void handleCacheClearError(RuntimeException exception, Cache cache) {
                log.error("Cache clear error - cache: {}, error: {}", 
                    cache.getName(), exception.getMessage());
            }
        };
    }
}
```

### 6.7 Common Pitfalls

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    COMMON CACHING PITFALLS                                       │
│                                                                                  │
│  ❌ PITFALL 1: Caching Mutable Objects                                          │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Cached object modified, affecting all consumers                        │
│  Solution: Cache immutable objects or deep copies                               │
│                                                                                  │
│  // BAD                                                                          │
│  @Cacheable("products")                                                          │
│  public Product getProduct(Long id) { ... }  // Returns mutable entity         │
│                                                                                  │
│  // Elsewhere: product.setPrice(0); // Modifies cached object!                  │
│                                                                                  │
│  // GOOD                                                                         │
│  @Cacheable("products")                                                          │
│  public ProductDTO getProduct(Long id) { ... }  // Returns immutable DTO        │
│                                                                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│  ❌ PITFALL 2: Caching in @Transactional Methods                                │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Cache updated before transaction commits                              │
│  Solution: Use @CachePut after transaction or transaction-aware cache manager   │
│                                                                                  │
│  // BAD                                                                          │
│  @Transactional                                                                  │
│  @CachePut(value = "products", key = "#result.id")                              │
│  public Product create(ProductRequest req) {                                    │
│      Product saved = repository.save(entity);                                   │
│      // Cache updated, but transaction might rollback!                          │
│      return saved;                                                              │
│  }                                                                               │
│                                                                                  │
│  // GOOD: Use transactionAware cache manager                                    │
│  RedisCacheManager.builder(factory).transactionAware().build();                 │
│                                                                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│  ❌ PITFALL 3: No Size Limits                                                   │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Cache grows unbounded, causing OOM                                    │
│  Solution: Always set maximum size or entry count                               │
│                                                                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│  ❌ PITFALL 4: Caching Large Objects                                            │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Serialization overhead, memory pressure                               │
│  Solution: Cache references/IDs, or use off-heap storage                        │
│                                                                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│  ❌ PITFALL 5: Self-Invocation Doesn't Trigger Cache                            │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Calling @Cacheable method from same class bypasses proxy             │
│  Solution: Extract to separate service or use self-injection                    │
│                                                                                  │
│  // BAD                                                                          │
│  public Product processProduct(Long id) {                                       │
│      return getProduct(id);  // Cache bypassed - direct method call!           │
│  }                                                                               │
│                                                                                  │
│  @Cacheable("products")                                                          │
│  public Product getProduct(Long id) { ... }                                     │
│                                                                                  │
│  // GOOD: Self-injection                                                         │
│  @Autowired @Lazy private ProductService self;                                  │
│                                                                                  │
│  public Product processProduct(Long id) {                                       │
│      return self.getProduct(id);  // Goes through proxy - cache works!         │
│  }                                                                               │
│                                                                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│  ❌ PITFALL 6: Ignoring Cache Warm-up                                           │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Cold cache after restart causes degraded performance                  │
│  Solution: Implement cache warming at startup                                   │
│                                                                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│                                                                                  │
│  ❌ PITFALL 7: Serialization Issues with Redis                                  │
│  ──────────────────────────────────────────────────────────────────────────     │
│  Problem: Class changes break deserialization                                   │
│  Solution: Use JSON serialization, version your DTOs                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Fix for self-invocation issue
@Service
public class ProductService {
    
    @Autowired
    @Lazy  // Prevent circular dependency
    private ProductService self;
    
    public Product processProduct(Long id) {
        // Use self reference to go through proxy
        return self.getProduct(id);  // Cache will work!
    }
    
    @Cacheable("products")
    public Product getProduct(Long id) {
        return repository.findById(id).orElseThrow();
    }
}

// Fix for serialization issues
@Configuration
public class RedisCacheConfig {
    
    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        // Use JSON serialization (more resilient to class changes)
        ObjectMapper objectMapper = new ObjectMapper()
            .registerModule(new JavaTimeModule())
            .disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS)
            .activateDefaultTyping(
                LaissezFaireSubTypeValidator.instance,
                ObjectMapper.DefaultTyping.NON_FINAL,
                JsonTypeInfo.As.PROPERTY
            );
        
        GenericJackson2JsonRedisSerializer serializer = 
            new GenericJackson2JsonRedisSerializer(objectMapper);
        
        RedisCacheConfiguration config = RedisCacheConfiguration.defaultCacheConfig()
            .serializeValuesWith(RedisSerializationContext.SerializationPair
                .fromSerializer(serializer));
        
        return RedisCacheManager.builder(factory)
            .cacheDefaults(config)
            .build();
    }
}
```

---

## 7. Conclusion

### 7.1 Why Caching Matters

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    THE IMPACT OF CACHING                                         │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  WITHOUT CACHING                    WITH CACHING                        │   │
│  │  ────────────────                   ─────────────                        │   │
│  │                                                                          │   │
│  │  Response Time: 500ms               Response Time: 5ms                  │   │
│  │  DB Load: 10,000 QPS                DB Load: 1,000 QPS                  │   │
│  │  Server Cost: $5,000/mo             Server Cost: $1,000/mo              │   │
│  │  User Experience: Poor              User Experience: Excellent          │   │
│  │  Scalability: Limited               Scalability: High                   │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  KEY BENEFITS:                                                                   │
│                                                                                  │
│  ⚡ PERFORMANCE                                                                  │
│     • Sub-millisecond response times for cached data                           │
│     • Reduced latency improves user experience                                 │
│     • Faster page loads increase conversion rates                              │
│                                                                                  │
│  💰 COST REDUCTION                                                              │
│     • Fewer database read operations                                           │
│     • Reduced compute requirements                                             │
│     • Lower infrastructure costs                                               │
│                                                                                  │
│  📈 SCALABILITY                                                                 │
│     • Handle traffic spikes without database bottlenecks                       │
│     • Scale horizontally with distributed caching                              │
│     • Support millions of concurrent users                                     │
│                                                                                  │
│  🛡️ RESILIENCE                                                                  │
│     • Serve data even when database is slow/unavailable                        │
│     • Circuit breaker patterns with cache fallback                             │
│     • Graceful degradation under load                                          │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 How Spring Boot Simplifies Caching

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SPRING BOOT CACHING ADVANTAGES                                │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  SIMPLE TO START                                                         │   │
│  │  ─────────────────                                                       │   │
│  │                                                                          │   │
│  │  1. Add @EnableCaching                                                  │   │
│  │  2. Add @Cacheable to methods                                           │   │
│  │  3. Done! (Uses ConcurrentMap by default)                               │   │
│  │                                                                          │   │
│  │  @EnableCaching                                                          │   │
│  │  @SpringBootApplication                                                  │   │
│  │  public class App { }                                                    │   │
│  │                                                                          │   │
│  │  @Cacheable("products")                                                  │   │
│  │  public Product getProduct(Long id) { return repo.findById(id); }       │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  PROVIDER AGNOSTIC                                                       │   │
│  │  ─────────────────                                                       │   │
│  │                                                                          │   │
│  │  Same annotations work with any provider:                               │   │
│  │                                                                          │   │
│  │  [Your Code] → [Spring Cache Abstraction] → [ConcurrentMap]             │   │
│  │                                            → [Caffeine]                  │   │
│  │                                            → [Redis]                     │   │
│  │                                            → [Ehcache]                   │   │
│  │                                            → [Hazelcast]                 │   │
│  │                                                                          │   │
│  │  Change provider without changing business code!                         │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  AUTO-CONFIGURATION                                                      │   │
│  │  ──────────────────                                                      │   │
│  │                                                                          │   │
│  │  Spring Boot auto-detects cache providers:                              │   │
│  │                                                                          │   │
│  │  • Add dependency → Provider auto-configured                            │   │
│  │  • Sensible defaults out of the box                                     │   │
│  │  • Customize via application.yml                                        │   │
│  │  • Production-ready metrics integration                                 │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Selecting the Right Caching Provider

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE PROVIDER SELECTION FLOWCHART                            │
│                                                                                  │
│                              START                                               │
│                                │                                                 │
│                                ▼                                                 │
│                    ┌─────────────────────┐                                      │
│                    │ Multiple instances  │                                      │
│                    │  or microservices?  │                                      │
│                    └─────────────────────┘                                      │
│                      │               │                                           │
│                     YES              NO                                          │
│                      │               │                                           │
│                      ▼               ▼                                           │
│              ┌───────────────┐  ┌───────────────┐                               │
│              │ Need session  │  │ High perf     │                               │
│              │ or data sync? │  │ requirements? │                               │
│              └───────────────┘  └───────────────┘                               │
│                │          │        │          │                                  │
│               YES         NO      YES         NO                                 │
│                │          │        │          │                                  │
│                ▼          ▼        ▼          ▼                                  │
│          ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────┐                │
│          │Hazelcast│  │  Redis  │  │Caffeine │  │ConcurrentMap│                │
│          └─────────┘  └─────────┘  └─────────┘  └─────────────┘                │
│                                                                                  │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│  QUICK DECISION GUIDE:                                                          │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐    │
│  │  SITUATION                              RECOMMENDED PROVIDER           │    │
│  ├────────────────────────────────────────────────────────────────────────┤    │
│  │  Single instance, simple needs          ConcurrentMap (default)       │    │
│  │  Single instance, high performance      Caffeine                      │    │
│  │  Multi-instance, shared cache           Redis                         │    │
│  │  Multi-instance, data grid needs        Hazelcast                     │    │
│  │  Tiered caching (disk + memory)         Ehcache                       │    │
│  │  JCache compliance required             Any JCache provider           │    │
│  │  Best of both worlds                    Caffeine (L1) + Redis (L2)   │    │
│  └────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 7.4 Final Recommendations

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    KEY TAKEAWAYS                                                 │
│                                                                                  │
│  1. START SIMPLE                                                                 │
│     • Begin with ConcurrentMap or Caffeine                                      │
│     • Evolve to distributed cache when needed                                   │
│     • Don't over-engineer from the start                                        │
│                                                                                  │
│  2. CACHE THE RIGHT DATA                                                         │
│     ✓ Frequently accessed, rarely changed                                       │
│     ✓ Expensive to compute or fetch                                             │
│     ✓ Can tolerate some staleness                                               │
│     ✗ Highly dynamic, real-time data                                            │
│     ✗ Sensitive data (without encryption)                                       │
│                                                                                  │
│  3. DESIGN FOR FAILURE                                                           │
│     • Implement graceful degradation                                            │
│     • Application should work without cache                                     │
│     • Use timeouts and circuit breakers                                         │
│                                                                                  │
│  4. MONITOR AND TUNE                                                             │
│     • Track hit/miss ratios                                                     │
│     • Monitor memory usage                                                      │
│     • Adjust TTL and size based on patterns                                     │
│                                                                                  │
│  5. HANDLE COMMON PROBLEMS                                                       │
│     • Cache stampede → Use sync=true or mutex                                  │
│     • Cache penetration → Bloom filter or cache nulls                          │
│     • Cache breakdown → Never expire hot keys                                  │
│     • Invalidation → Event-driven architecture                                 │
│                                                                                  │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│  REMEMBER:                                                                       │
│  ─────────                                                                       │
│  "Caching is not a silver bullet. It's a trade-off between                     │
│   consistency and performance. Understand your requirements                     │
│   and choose the right strategy for your use case."                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Document Summary

| Section | Topics Covered |
|---------|----------------|
| 1. Introduction | What is caching, why it matters, caching patterns, real-world examples |
| 2. Core Concepts | Keys, eviction, TTL, invalidation, consistency, local vs distributed |
| 3. Spring Boot Caching | Abstraction, annotations, key generation, configuration |
| 4. Providers | ConcurrentMap, Caffeine, Ehcache, Redis, Hazelcast, Infinispan |
| 5. Provider Selection | Comparison, decision matrix, scenario-based guidance |
| 6. Advanced Topics | Invalidation, stampede, penetration, breakdown, monitoring, pitfalls |
| 7. Conclusion | Why caching matters, Spring Boot advantages, selection guide |
| 8. Redis Deep Dive | Complete Redis caching implementation with Spring Boot |

---

# Part B: Redis Caching Deep Dive

---

## 8. Introduction to Redis for Caching

### 8.1 What is Redis?

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              WHAT IS REDIS?                                      │
│                                                                                  │
│  Redis = REmote DIctionary Server                                               │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  • Open-source, in-memory data structure store                          │   │
│  │  • Can be used as database, cache, message broker, and streaming engine │   │
│  │  • Supports various data structures:                                    │   │
│  │    - Strings, Hashes, Lists, Sets, Sorted Sets                         │   │
│  │    - Bitmaps, HyperLogLogs, Geospatial indexes, Streams               │   │
│  │  • Single-threaded event loop (extremely fast)                         │   │
│  │  • Sub-millisecond latency                                              │   │
│  │  • Written in C, runs on Linux, macOS, Windows (via WSL)               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  KEY CHARACTERISTICS:                                                            │
│  ────────────────────                                                            │
│                                                                                  │
│  ⚡ SPEED           - 100,000+ operations per second                           │
│  💾 IN-MEMORY       - All data stored in RAM                                    │
│  🔄 PERSISTENCE     - Optional disk persistence (RDB, AOF)                      │
│  📡 REPLICATION     - Master-replica architecture                               │
│  🔒 ATOMIC          - All operations are atomic                                  │
│  📊 RICH DATA TYPES - Beyond simple key-value                                    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Why Redis is Popular for Caching

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    WHY REDIS FOR CACHING?                                        │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  1. BLAZING FAST PERFORMANCE                                            │   │
│  │     • In-memory storage = microsecond latency                           │   │
│  │     • Single-threaded = no lock contention                              │   │
│  │     • Optimized data structures                                         │   │
│  │                                                                          │   │
│  │  2. DISTRIBUTED BY DESIGN                                               │   │
│  │     • Share cache across multiple application instances                 │   │
│  │     • Built-in clustering and replication                               │   │
│  │     • Horizontal scalability                                            │   │
│  │                                                                          │   │
│  │  3. RICH FEATURE SET                                                    │   │
│  │     • TTL support (automatic expiration)                                │   │
│  │     • Multiple eviction policies                                        │   │
│  │     • Pub/Sub for cache invalidation                                    │   │
│  │     • Transactions and Lua scripting                                    │   │
│  │                                                                          │   │
│  │  4. BATTLE-TESTED                                                       │   │
│  │     • Used by Twitter, GitHub, Pinterest, Snapchat                      │   │
│  │     • Mature ecosystem                                                   │   │
│  │     • Excellent Spring Boot integration                                 │   │
│  │                                                                          │   │
│  │  5. VERSATILE                                                           │   │
│  │     • Not just caching: sessions, rate limiting, queues                 │   │
│  │     • Multiple data types for different use cases                       │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 8.3 In-Memory Data Store Concept

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    IN-MEMORY DATA STORE                                          │
│                                                                                  │
│  TRADITIONAL DATABASE                      REDIS (IN-MEMORY)                    │
│  ────────────────────                      ─────────────────                    │
│                                                                                  │
│  ┌─────────────────┐                      ┌─────────────────┐                  │
│  │   Application   │                      │   Application   │                  │
│  └────────┬────────┘                      └────────┬────────┘                  │
│           │                                        │                            │
│           ▼                                        ▼                            │
│  ┌─────────────────┐                      ┌─────────────────┐                  │
│  │    Database     │                      │      Redis      │                  │
│  │    Process      │                      │    (in RAM)     │                  │
│  └────────┬────────┘                      └────────┬────────┘                  │
│           │                                        │                            │
│           ▼                                        ▼ (optional)                 │
│  ┌─────────────────┐                      ┌─────────────────┐                  │
│  │   Disk (HDD/    │                      │   Disk (RDB/    │                  │
│  │      SSD)       │                      │     AOF)        │                  │
│  └─────────────────┘                      └─────────────────┘                  │
│                                                                                  │
│  Latency: 1-10ms                          Latency: <1ms                        │
│  (disk I/O bottleneck)                    (RAM access only)                    │
│                                                                                  │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│  WHY RAM IS FASTER:                                                             │
│  ──────────────────                                                             │
│                                                                                  │
│  │ Storage Type │ Access Time    │ Relative Speed │                            │
│  │──────────────│────────────────│────────────────│                            │
│  │ CPU Cache    │ ~1 ns          │ 1x             │                            │
│  │ RAM          │ ~100 ns        │ 100x           │                            │
│  │ SSD          │ ~100,000 ns    │ 100,000x       │                            │
│  │ HDD          │ ~10,000,000 ns │ 10,000,000x    │                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 8.4 Redis as Cache vs Database

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS AS CACHE VS DATABASE                                    │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  ASPECT              │ REDIS AS CACHE      │ REDIS AS DATABASE         │   │
│  │  ────────────────────┼─────────────────────┼────────────────────────── │   │
│  │  Data Persistence    │ Optional/Disabled   │ Required (RDB + AOF)     │   │
│  │  Data Loss Tolerance │ Acceptable          │ Not acceptable           │   │
│  │  Primary Data Store  │ No (DB is primary)  │ Yes                      │   │
│  │  Eviction Policy     │ Enabled (LRU/LFU)   │ noeviction               │   │
│  │  Memory Management   │ Bounded, auto-evict │ Scale memory as needed   │   │
│  │  TTL Usage           │ Common              │ Rare                      │   │
│  │  Replication Focus   │ Read scaling        │ High availability        │   │
│  │  Recovery on Failure │ Rebuild from DB     │ Restore from disk        │   │
│  │  Cost                │ Lower (temporary)   │ Higher (persistent)      │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  CACHE MODE CONFIGURATION:                                                       │
│  ──────────────────────────                                                      │
│                                                                                  │
│  # redis.conf for CACHE mode                                                    │
│  maxmemory 2gb                     # Limit memory                               │
│  maxmemory-policy allkeys-lru      # Evict when full                           │
│  appendonly no                     # No persistence                             │
│  save ""                           # Disable RDB snapshots                      │
│                                                                                  │
│  # redis.conf for DATABASE mode                                                 │
│  maxmemory 0                       # No limit (use all available)              │
│  maxmemory-policy noeviction       # Return error when full                    │
│  appendonly yes                    # Enable AOF                                 │
│  appendfsync everysec              # Sync to disk every second                 │
│  save 900 1                        # RDB snapshot every 15 min if 1 key changed│
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 8.5 When Redis Should Be Used for Caching

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    WHEN TO USE REDIS FOR CACHING                                 │
│                                                                                  │
│  ✅ USE REDIS CACHING WHEN:                                                     │
│  ──────────────────────────                                                      │
│                                                                                  │
│  1. DISTRIBUTED APPLICATIONS                                                     │
│     • Multiple application instances need shared cache                          │
│     • Microservices architecture                                                │
│     • Kubernetes/containerized deployments                                      │
│                                                                                  │
│  2. HIGH-TRAFFIC SCENARIOS                                                       │
│     • Read-heavy workloads (>80% reads)                                        │
│     • API response caching                                                      │
│     • Database query result caching                                             │
│                                                                                  │
│  3. SESSION MANAGEMENT                                                           │
│     • Stateless application instances                                           │
│     • Sticky sessions not desired                                               │
│     • Session sharing across services                                           │
│                                                                                  │
│  4. REAL-TIME FEATURES                                                           │
│     • Leaderboards, counters                                                    │
│     • Rate limiting                                                              │
│     • Real-time analytics                                                        │
│                                                                                  │
│  5. EXPENSIVE OPERATIONS                                                         │
│     • Complex database queries                                                  │
│     • Third-party API responses                                                 │
│     • Computed/aggregated data                                                  │
│                                                                                  │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│  ❌ DON'T USE REDIS CACHING WHEN:                                               │
│  ────────────────────────────────                                               │
│                                                                                  │
│  • Single instance application (use Caffeine instead - less overhead)          │
│  • Data changes very frequently (cache invalidation overhead)                   │
│  • Strong consistency required (eventual consistency trade-off)                 │
│  • Large objects (>10MB per key - serialization overhead)                      │
│  • Budget constraints and simple needs (adds infrastructure complexity)        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 9. How Redis Works as a Cache

### 9.1 Key-Value Storage Model

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS KEY-VALUE MODEL                                         │
│                                                                                  │
│  Redis stores data as key-value pairs where:                                    │
│  • KEY: Always a string (binary-safe, max 512MB)                               │
│  • VALUE: Can be various data types                                             │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  KEY                          VALUE (Type: String)                      │   │
│  │  ═══════════════════════════════════════════════════════════════════   │   │
│  │  "user:1001"                  "{\"id\":1001,\"name\":\"John\"}"        │   │
│  │  "product:SKU-ABC"            "{\"sku\":\"ABC\",\"price\":99.99}"      │   │
│  │  "session:abc123"             "user_data_serialized"                   │   │
│  │  "counter:page_views"         "1542367"                                 │   │
│  │                                                                          │   │
│  │  KEY                          VALUE (Type: Hash)                        │   │
│  │  ═══════════════════════════════════════════════════════════════════   │   │
│  │  "user:1001:profile"          {name: "John", email: "j@mail.com"}      │   │
│  │                                                                          │   │
│  │  KEY                          VALUE (Type: List)                        │   │
│  │  ═══════════════════════════════════════════════════════════════════   │   │
│  │  "user:1001:recent_orders"    [order1, order2, order3]                 │   │
│  │                                                                          │   │
│  │  KEY                          VALUE (Type: Set)                         │   │
│  │  ═══════════════════════════════════════════════════════════════════   │   │
│  │  "product:123:tags"           {electronics, sale, featured}            │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  KEY NAMING CONVENTIONS:                                                         │
│  ───────────────────────                                                         │
│  • Use colons as separators: "object:id:field"                                 │
│  • Keep keys short but descriptive                                              │
│  • Examples: "user:1001", "cache:products:category:electronics"                │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 TTL (Time To Live)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    TTL (TIME TO LIVE)                                            │
│                                                                                  │
│  TTL defines how long a key remains in Redis before automatic deletion.         │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  TIMELINE:                                                               │   │
│  │                                                                          │   │
│  │  T=0          SET key "value" EX 300 (5 minutes)                        │   │
│  │               ┌───────────────────────┐                                 │   │
│  │               │ KEY EXISTS, TTL=300s  │                                 │   │
│  │               └───────────────────────┘                                 │   │
│  │                                                                          │   │
│  │  T=150s       GET key → "value" ✓                                       │   │
│  │               ┌───────────────────────┐                                 │   │
│  │               │ KEY EXISTS, TTL=150s  │                                 │   │
│  │               └───────────────────────┘                                 │   │
│  │                                                                          │   │
│  │  T=300s       KEY EXPIRES (auto-deleted)                                │   │
│  │               ┌───────────────────────┐                                 │   │
│  │               │ KEY DOES NOT EXIST    │                                 │   │
│  │               └───────────────────────┘                                 │   │
│  │                                                                          │   │
│  │  T=301s       GET key → nil ✗                                           │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  REDIS TTL COMMANDS:                                                             │
│  ────────────────────                                                            │
│                                                                                  │
│  SET key value EX 300        # Set with 300 seconds TTL                         │
│  SET key value PX 300000     # Set with 300000 milliseconds TTL                 │
│  SETEX key 300 value         # Set with TTL (seconds)                           │
│  EXPIRE key 300              # Set TTL on existing key (seconds)                │
│  PEXPIRE key 300000          # Set TTL on existing key (milliseconds)           │
│  TTL key                     # Get remaining TTL (seconds)                       │
│  PTTL key                    # Get remaining TTL (milliseconds)                  │
│  PERSIST key                 # Remove TTL (key becomes persistent)              │
│  EXPIREAT key timestamp      # Set expiry at Unix timestamp                     │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 9.3 Expiration Policies

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS EXPIRATION POLICIES                                     │
│                                                                                  │
│  Redis uses TWO strategies to handle key expiration:                            │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  1. PASSIVE (LAZY) EXPIRATION                                           │   │
│  │  ────────────────────────────────                                        │   │
│  │  Key is checked and deleted only when accessed                          │   │
│  │                                                                          │   │
│  │  GET expired_key                                                         │   │
│  │    → Redis checks: Is TTL expired?                                      │   │
│  │    → YES: Delete key, return nil                                        │   │
│  │    → NO: Return value                                                   │   │
│  │                                                                          │   │
│  │  Pros: Zero CPU overhead when key not accessed                          │   │
│  │  Cons: Expired keys linger in memory until accessed                     │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  2. ACTIVE EXPIRATION (Background)                                      │   │
│  │  ─────────────────────────────────                                       │   │
│  │  Redis periodically samples keys and deletes expired ones               │   │
│  │                                                                          │   │
│  │  Every 100ms (10 times/second):                                         │   │
│  │    1. Sample 20 random keys with TTL                                    │   │
│  │    2. Delete all expired keys found                                     │   │
│  │    3. If >25% expired, repeat immediately                               │   │
│  │                                                                          │   │
│  │  This ensures expired keys are eventually removed                       │   │
│  │  without blocking the main thread                                       │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  COMBINED APPROACH:                                                              │
│  ──────────────────                                                              │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────┐     │
│  │ Passive (on access)  +  Active (background)  =  Efficient Expiration │     │
│  │                                                                       │     │
│  │ • No expired data returned to clients (passive)                      │     │
│  │ • Memory reclaimed even for unaccessed keys (active)                 │     │
│  │ • Minimal CPU impact                                                 │     │
│  └───────────────────────────────────────────────────────────────────────┘     │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 9.4 Eviction Policies

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS EVICTION POLICIES                                       │
│                                                                                  │
│  When Redis reaches maxmemory, it uses eviction policy to free space.           │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  POLICY           │ DESCRIPTION                                         │   │
│  │  ─────────────────┼────────────────────────────────────────────────────│   │
│  │  noeviction       │ Return error on write when memory full             │   │
│  │  allkeys-lru      │ Evict LEAST RECENTLY USED key from ALL keys        │   │
│  │  volatile-lru     │ Evict LRU key from keys WITH TTL set               │   │
│  │  allkeys-lfu      │ Evict LEAST FREQUENTLY USED key from ALL keys      │   │
│  │  volatile-lfu     │ Evict LFU key from keys WITH TTL set               │   │
│  │  allkeys-random   │ Evict RANDOM key from ALL keys                     │   │
│  │  volatile-random  │ Evict RANDOM key from keys WITH TTL set            │   │
│  │  volatile-ttl     │ Evict key with SHORTEST TTL remaining              │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ══════════════════════════════════════════════════════════════════════════    │
│                                                                                  │
│  DETAILED EXPLANATION:                                                           │
│  ─────────────────────                                                           │
│                                                                                  │
│  ❶ NOEVICTION                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ • Returns OOM error when maxmemory reached                              │   │
│  │ • Read operations still work                                            │   │
│  │ • Use case: Redis as database (don't want to lose data)                │   │
│  │ • NOT recommended for caching                                           │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❷ ALLKEYS-LRU (Recommended for Caching)                                        │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ • Evicts least recently used keys from ALL keys                        │   │
│  │ • Best general-purpose cache eviction                                  │   │
│  │ • Keeps "hot" (frequently accessed) data                               │   │
│  │ • Use case: Generic caching, unknown access patterns                   │   │
│  │                                                                          │   │
│  │ Example:                                                                 │   │
│  │   key1: accessed 1 minute ago  ← Kept                                  │   │
│  │   key2: accessed 10 minutes ago ← EVICTED                              │   │
│  │   key3: accessed 2 minutes ago  ← Kept                                 │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❸ VOLATILE-LRU                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ • Only evicts keys that HAVE a TTL set                                 │   │
│  │ • Keys without TTL are never evicted                                   │   │
│  │ • Use case: Mix of cache data (TTL) and persistent data (no TTL)      │   │
│  │                                                                          │   │
│  │ Example:                                                                 │   │
│  │   key1 (TTL=300): accessed 10 min ago ← EVICTED                        │   │
│  │   key2 (no TTL):  accessed 1 hour ago ← KEPT (no TTL)                  │   │
│  │   key3 (TTL=600): accessed 1 min ago  ← Kept                           │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❹ ALLKEYS-LFU (Best for Frequency-Based Access)                                │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ • Evicts least FREQUENTLY used keys                                    │   │
│  │ • Tracks access count over time (decays)                               │   │
│  │ • Better than LRU for "warm" data that's accessed periodically        │   │
│  │ • Use case: When some items are accessed in bursts                     │   │
│  │                                                                          │   │
│  │ Example:                                                                 │   │
│  │   key1: accessed 100 times today    ← Kept                             │   │
│  │   key2: accessed 2 times today      ← EVICTED                          │   │
│  │   key3: accessed 50 times today     ← Kept                             │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❺ VOLATILE-LFU                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ • LFU eviction but only on keys with TTL                               │   │
│  │ • Combines frequency tracking with TTL requirement                     │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❻ VOLATILE-TTL                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ • Evicts keys with shortest remaining TTL first                        │   │
│  │ • Use case: When TTL correlates with importance                        │   │
│  │                                                                          │   │
│  │ Example:                                                                 │   │
│  │   key1 (TTL remaining: 10s)  ← EVICTED first                           │   │
│  │   key2 (TTL remaining: 300s) ← Kept                                    │   │
│  │   key3 (TTL remaining: 60s)  ← Evicted second                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ══════════════════════════════════════════════════════════════════════════    │
│                                                                                  │
│  CONFIGURATION:                                                                  │
│  ──────────────                                                                  │
│                                                                                  │
│  # redis.conf                                                                   │
│  maxmemory 2gb                                                                  │
│  maxmemory-policy allkeys-lru                                                   │
│                                                                                  │
│  # Or via Redis CLI                                                             │
│  CONFIG SET maxmemory 2gb                                                       │
│  CONFIG SET maxmemory-policy allkeys-lru                                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 9.5 Cache Hit vs Cache Miss Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE HIT VS CACHE MISS                                       │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         CACHE HIT FLOW                                   │   │
│  │                                                                          │   │
│  │    ┌──────────┐    ┌─────────────┐                                      │   │
│  │    │  Client  │───▶│ Application │                                      │   │
│  │    └──────────┘    └──────┬──────┘                                      │   │
│  │         ▲                 │ 1. Request data                             │   │
│  │         │                 ▼                                              │   │
│  │         │          ┌─────────────┐                                      │   │
│  │         │          │    Redis    │                                      │   │
│  │         │          │   (Cache)   │                                      │   │
│  │         │          └──────┬──────┘                                      │   │
│  │         │                 │ 2. Data FOUND ✓                             │   │
│  │         │                 ▼                                              │   │
│  │    5. Return        3. Return cached data                               │   │
│  │    Response         (No DB query needed)                                │   │
│  │                                                                          │   │
│  │    Latency: ~1-5ms                                                      │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         CACHE MISS FLOW                                  │   │
│  │                                                                          │   │
│  │    ┌──────────┐    ┌─────────────┐                                      │   │
│  │    │  Client  │───▶│ Application │                                      │   │
│  │    └──────────┘    └──────┬──────┘                                      │   │
│  │         ▲                 │ 1. Request data                             │   │
│  │         │                 ▼                                              │   │
│  │         │          ┌─────────────┐                                      │   │
│  │         │          │    Redis    │                                      │   │
│  │         │          │   (Cache)   │                                      │   │
│  │         │          └──────┬──────┘                                      │   │
│  │         │                 │ 2. Data NOT found ✗                         │   │
│  │         │                 ▼                                              │   │
│  │         │          ┌─────────────┐                                      │   │
│  │         │          │  Database   │                                      │   │
│  │         │          └──────┬──────┘                                      │   │
│  │         │                 │ 3. Query database                           │   │
│  │         │                 ▼                                              │   │
│  │         │          4. Store in Redis (with TTL)                         │   │
│  │         │                 │                                              │   │
│  │    6. Return              ▼                                              │   │
│  │    Response         5. Return data                                      │   │
│  │                                                                          │   │
│  │    Latency: ~50-500ms (DB query included)                               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Cache Hit/Miss Implementation
@Service
@RequiredArgsConstructor
@Slf4j
public class ProductService {
    
    private final ProductRepository repository;
    private final RedisTemplate<String, Object> redisTemplate;
    
    private static final String CACHE_PREFIX = "product:";
    private static final Duration CACHE_TTL = Duration.ofMinutes(30);
    
    public Product getProduct(Long id) {
        String cacheKey = CACHE_PREFIX + id;
        
        // Step 1: Try to get from cache
        Product cached = (Product) redisTemplate.opsForValue().get(cacheKey);
        
        if (cached != null) {
            // CACHE HIT
            log.debug("Cache HIT for key: {}", cacheKey);
            return cached;
        }
        
        // CACHE MISS
        log.debug("Cache MISS for key: {}", cacheKey);
        
        // Step 2: Fetch from database
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        
        // Step 3: Store in cache with TTL
        redisTemplate.opsForValue().set(cacheKey, product, CACHE_TTL);
        log.debug("Cached product with key: {}, TTL: {}", cacheKey, CACHE_TTL);
        
        return product;
    }
}
```

### 9.6 Cache-Aside Pattern with Redis

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CACHE-ASIDE PATTERN (Lazy Loading)                            │
│                                                                                  │
│  Application manages cache explicitly - reads from cache first,                 │
│  falls back to database on miss, then populates cache.                          │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │                        READ OPERATION                                    │   │
│  │                                                                          │   │
│  │    ┌──────────────┐                                                     │   │
│  │    │  Application │                                                     │   │
│  │    └───────┬──────┘                                                     │   │
│  │            │                                                             │   │
│  │      ┌─────┴─────┐                                                      │   │
│  │      │           │                                                       │   │
│  │      ▼           │                                                       │   │
│  │  ┌───────┐       │                                                       │   │
│  │  │ Redis │       │ 1. Check cache                                       │   │
│  │  └───┬───┘       │                                                       │   │
│  │      │           │                                                       │   │
│  │   HIT│    MISS   │                                                       │   │
│  │      │     ┌─────┘                                                       │   │
│  │      │     │                                                             │   │
│  │      │     ▼                                                             │   │
│  │      │  ┌──────────┐                                                    │   │
│  │      │  │ Database │  2. Query DB                                       │   │
│  │      │  └────┬─────┘                                                    │   │
│  │      │       │                                                           │   │
│  │      │       ▼                                                           │   │
│  │      │  ┌───────┐                                                        │   │
│  │      │  │ Redis │ 3. Populate cache                                     │   │
│  │      │  └───┬───┘                                                        │   │
│  │      │      │                                                            │   │
│  │      ▼      ▼                                                            │   │
│  │    Return data                                                           │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │                        WRITE OPERATION                                   │   │
│  │                                                                          │   │
│  │    ┌──────────────┐                                                     │   │
│  │    │  Application │                                                     │   │
│  │    └───────┬──────┘                                                     │   │
│  │            │                                                             │   │
│  │      ┌─────┴─────┐                                                      │   │
│  │      ▼           ▼                                                       │   │
│  │  ┌──────────┐  ┌───────┐                                                │   │
│  │  │ Database │  │ Redis │                                                │   │
│  │  │  UPDATE  │  │EVICT/ │                                                │   │
│  │  │          │  │UPDATE │                                                │   │
│  │  └──────────┘  └───────┘                                                │   │
│  │                                                                          │   │
│  │  Option A: Invalidate cache (safer)                                     │   │
│  │  Option B: Update cache (faster if consistent)                          │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Cache-Aside Pattern Implementation
@Service
@RequiredArgsConstructor
@Slf4j
public class ProductCacheAsideService {
    
    private final ProductRepository repository;
    private final StringRedisTemplate redisTemplate;
    private final ObjectMapper objectMapper;
    
    private static final String CACHE_PREFIX = "product:";
    private static final Duration CACHE_TTL = Duration.ofMinutes(30);
    
    // READ: Cache-Aside
    public Product getProduct(Long id) {
        String cacheKey = CACHE_PREFIX + id;
        
        // 1. Try cache first
        String cached = redisTemplate.opsForValue().get(cacheKey);
        if (cached != null) {
            log.debug("Cache HIT: {}", cacheKey);
            return deserialize(cached);
        }
        
        // 2. Cache miss - query database
        log.debug("Cache MISS: {}", cacheKey);
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        
        // 3. Populate cache
        redisTemplate.opsForValue().set(cacheKey, serialize(product), CACHE_TTL);
        
        return product;
    }
    
    // WRITE: Update DB, then invalidate cache
    @Transactional
    public Product updateProduct(Long id, ProductUpdateRequest request) {
        String cacheKey = CACHE_PREFIX + id;
        
        // 1. Update database
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        product.setName(request.getName());
        product.setPrice(request.getPrice());
        Product saved = repository.save(product);
        
        // 2. Invalidate cache (safer than updating)
        redisTemplate.delete(cacheKey);
        log.debug("Invalidated cache: {}", cacheKey);
        
        return saved;
    }
    
    // DELETE: Remove from DB and cache
    @Transactional
    public void deleteProduct(Long id) {
        String cacheKey = CACHE_PREFIX + id;
        
        repository.deleteById(id);
        redisTemplate.delete(cacheKey);
        log.debug("Deleted product and cache: {}", id);
    }
    
    private String serialize(Product product) {
        try {
            return objectMapper.writeValueAsString(product);
        } catch (JsonProcessingException e) {
            throw new CacheSerializationException("Failed to serialize product", e);
        }
    }
    
    private Product deserialize(String json) {
        try {
            return objectMapper.readValue(json, Product.class);
        } catch (JsonProcessingException e) {
            throw new CacheSerializationException("Failed to deserialize product", e);
        }
    }
}
```

### 9.7 Write-Through vs Write-Behind with Redis

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    WRITE-THROUGH VS WRITE-BEHIND                                 │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  WRITE-THROUGH                                                           │   │
│  │  ═════════════                                                           │   │
│  │                                                                          │   │
│  │  Data written to cache AND database SYNCHRONOUSLY                       │   │
│  │                                                                          │   │
│  │    ┌──────────────┐                                                     │   │
│  │    │  Application │                                                     │   │
│  │    └───────┬──────┘                                                     │   │
│  │            │  WRITE                                                      │   │
│  │            ▼                                                             │   │
│  │       ┌────────┐                                                         │   │
│  │       │ Redis  │ ← 1. Write to cache (sync)                             │   │
│  │       └────┬───┘                                                         │   │
│  │            │                                                             │   │
│  │            ▼                                                             │   │
│  │       ┌──────────┐                                                       │   │
│  │       │ Database │ ← 2. Write to DB (sync, waits for completion)        │   │
│  │       └──────────┘                                                       │   │
│  │            │                                                             │   │
│  │            ▼                                                             │   │
│  │       Response returned only AFTER both writes complete                 │   │
│  │                                                                          │   │
│  │  ✓ Strong consistency                                                   │   │
│  │  ✓ Simple to understand                                                 │   │
│  │  ✗ Higher write latency (both operations are synchronous)               │   │
│  │  ✗ Failure in either store needs handling                               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  WRITE-BEHIND (Write-Back)                                              │   │
│  │  ═════════════════════════                                               │   │
│  │                                                                          │   │
│  │  Data written to cache, then ASYNCHRONOUSLY to database                 │   │
│  │                                                                          │   │
│  │    ┌──────────────┐                                                     │   │
│  │    │  Application │                                                     │   │
│  │    └───────┬──────┘                                                     │   │
│  │            │  WRITE                                                      │   │
│  │            ▼                                                             │   │
│  │       ┌────────┐                                                         │   │
│  │       │ Redis  │ ← 1. Write to cache (sync)                             │   │
│  │       └────┬───┘                                                         │   │
│  │            │                                                             │   │
│  │            ▼                                                             │   │
│  │       Response returned IMMEDIATELY                                      │   │
│  │                                                                          │   │
│  │            │ (async, background)                                        │   │
│  │            ▼                                                             │   │
│  │       ┌──────────┐                                                       │   │
│  │       │ Database │ ← 2. Write to DB (async, batched)                    │   │
│  │       └──────────┘                                                       │   │
│  │                                                                          │   │
│  │  ✓ Very low write latency                                               │   │
│  │  ✓ Can batch writes to database                                         │   │
│  │  ✗ Eventual consistency (DB lags behind cache)                          │   │
│  │  ✗ Data loss risk if cache fails before DB write                        │   │
│  │  ✗ More complex implementation                                           │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ══════════════════════════════════════════════════════════════════════════    │
│                                                                                  │
│  COMPARISON:                                                                     │
│                                                                                  │
│  │ Aspect          │ Write-Through     │ Write-Behind        │                 │
│  │─────────────────┼───────────────────┼─────────────────────│                 │
│  │ Consistency     │ Strong            │ Eventual            │                 │
│  │ Write Latency   │ Higher            │ Lower               │                 │
│  │ Data Loss Risk  │ Lower             │ Higher              │                 │
│  │ Complexity      │ Simple            │ Complex             │                 │
│  │ Use Case        │ Critical data     │ High write volume   │                 │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Write-Through Implementation
@Service
@RequiredArgsConstructor
@Slf4j
public class WriteThroughCacheService {
    
    private final ProductRepository repository;
    private final RedisTemplate<String, Product> redisTemplate;
    
    @Transactional
    public Product createProduct(ProductRequest request) {
        // 1. Save to database (synchronous)
        Product product = repository.save(mapToEntity(request));
        
        // 2. Write to cache (synchronous)
        String cacheKey = "product:" + product.getId();
        redisTemplate.opsForValue().set(cacheKey, product, Duration.ofMinutes(30));
        
        log.info("Write-through: Saved product {} to DB and cache", product.getId());
        
        // Response returned only after both complete
        return product;
    }
}

// Write-Behind Implementation
@Service
@RequiredArgsConstructor
@Slf4j
public class WriteBehindCacheService {
    
    private final ProductRepository repository;
    private final RedisTemplate<String, Product> redisTemplate;
    private final BlockingQueue<Product> writeQueue = new LinkedBlockingQueue<>();
    
    @PostConstruct
    public void startBackgroundWriter() {
        // Background thread to batch-write to database
        Thread writer = new Thread(() -> {
            List<Product> batch = new ArrayList<>();
            while (true) {
                try {
                    // Collect batch
                    Product item = writeQueue.poll(1, TimeUnit.SECONDS);
                    if (item != null) {
                        batch.add(item);
                    }
                    
                    // Flush batch when size reached or timeout
                    if (batch.size() >= 100 || (!batch.isEmpty() && item == null)) {
                        repository.saveAll(batch);
                        log.info("Write-behind: Flushed {} items to database", batch.size());
                        batch.clear();
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    break;
                }
            }
        });
        writer.setDaemon(true);
        writer.start();
    }
    
    public Product createProduct(ProductRequest request) {
        Product product = mapToEntity(request);
        product.setId(generateId());  // Pre-generate ID
        
        // 1. Write to cache immediately (synchronous)
        String cacheKey = "product:" + product.getId();
        redisTemplate.opsForValue().set(cacheKey, product, Duration.ofMinutes(30));
        
        // 2. Queue for database write (asynchronous)
        writeQueue.add(product);
        
        log.info("Write-behind: Cached product {}, queued for DB", product.getId());
        
        // Response returned immediately after cache write
        return product;
    }
}
```

---

## 10. Implementing Redis Cache in Spring Boot

### Step 1: Add Dependencies

```xml
<!-- Maven: pom.xml -->
<dependencies>
    <!-- Spring Boot Starter for Redis -->
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-data-redis</artifactId>
    </dependency>
    
    <!-- Spring Boot Starter for Caching -->
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-cache</artifactId>
    </dependency>
    
    <!-- Optional: Connection pooling (Lettuce uses Netty by default) -->
    <dependency>
        <groupId>org.apache.commons</groupId>
        <artifactId>commons-pool2</artifactId>
    </dependency>
    
    <!-- Optional: Jackson for JSON serialization -->
    <dependency>
        <groupId>com.fasterxml.jackson.core</groupId>
        <artifactId>jackson-databind</artifactId>
    </dependency>
    
    <!-- Optional: Jackson Java 8 date/time support -->
    <dependency>
        <groupId>com.fasterxml.jackson.datatype</groupId>
        <artifactId>jackson-datatype-jsr310</artifactId>
    </dependency>
</dependencies>
```

```groovy
// Gradle: build.gradle
dependencies {
    // Spring Boot Starter for Redis
    implementation 'org.springframework.boot:spring-boot-starter-data-redis'
    
    // Spring Boot Starter for Caching
    implementation 'org.springframework.boot:spring-boot-starter-cache'
    
    // Optional: Connection pooling
    implementation 'org.apache.commons:commons-pool2'
    
    // Optional: Jackson for JSON serialization
    implementation 'com.fasterxml.jackson.core:jackson-databind'
    implementation 'com.fasterxml.jackson.datatype:jackson-datatype-jsr310'
}
```

### Step 2: Configuration (application.yml)

```yaml
# application.yml - Complete Redis Configuration
spring:
  # Cache configuration
  cache:
    type: redis                          # Use Redis as cache provider
    redis:
      time-to-live: 3600000             # Default TTL: 1 hour (in milliseconds)
      cache-null-values: false           # Don't cache null values
      use-key-prefix: true               # Add prefix to keys
      key-prefix: "myapp:cache:"         # Custom prefix
      enable-statistics: true            # Enable cache statistics

  # Redis connection configuration
  redis:
    host: localhost                      # Redis server host
    port: 6379                           # Redis server port
    password: ${REDIS_PASSWORD:}         # Password (empty if not set)
    database: 0                          # Database index (0-15)
    timeout: 2000ms                      # Connection timeout
    connect-timeout: 2000ms              # Connect timeout
    
    # Lettuce client configuration (default client)
    lettuce:
      pool:
        enabled: true                    # Enable connection pooling
        max-active: 16                   # Max connections in pool
        max-idle: 8                      # Max idle connections
        min-idle: 4                      # Min idle connections
        max-wait: 1000ms                 # Max wait for connection
        time-between-eviction-runs: 60s  # Eviction check interval
      shutdown-timeout: 100ms            # Graceful shutdown timeout

# Logging for debugging
logging:
  level:
    org.springframework.cache: DEBUG
    org.springframework.data.redis: DEBUG
    io.lettuce.core: INFO
```

```yaml
# application-production.yml - Production Configuration
spring:
  redis:
    host: ${REDIS_HOST:redis-cluster.internal}
    port: ${REDIS_PORT:6379}
    password: ${REDIS_PASSWORD}
    ssl: true                            # Enable SSL for production
    
    lettuce:
      pool:
        max-active: 50                   # Higher for production load
        max-idle: 20
        min-idle: 10
        max-wait: 2000ms

  cache:
    redis:
      time-to-live: 1800000              # 30 minutes default TTL
      cache-null-values: true            # Cache nulls to prevent penetration
```

### Step 3: Enable Caching and Configure Beans

```java
// RedisCacheConfig.java - Complete Configuration Class
package com.example.config;

import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.databind.jsontype.impl.LaissezFaireSubTypeValidator;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.CacheManager;
import org.springframework.cache.annotation.CachingConfigurerSupport;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.cache.interceptor.CacheErrorHandler;
import org.springframework.cache.interceptor.KeyGenerator;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.cache.RedisCacheConfiguration;
import org.springframework.data.redis.cache.RedisCacheManager;
import org.springframework.data.redis.connection.RedisConnectionFactory;
import org.springframework.data.redis.connection.RedisStandaloneConfiguration;
import org.springframework.data.redis.connection.lettuce.LettuceConnectionFactory;
import org.springframework.data.redis.connection.lettuce.LettucePoolingClientConfiguration;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.serializer.GenericJackson2JsonRedisSerializer;
import org.springframework.data.redis.serializer.RedisSerializationContext;
import org.springframework.data.redis.serializer.StringRedisSerializer;

import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

@Configuration
@EnableCaching
public class RedisCacheConfig extends CachingConfigurerSupport {

    @Value("${spring.redis.host:localhost}")
    private String redisHost;
    
    @Value("${spring.redis.port:6379}")
    private int redisPort;
    
    @Value("${spring.redis.password:}")
    private String redisPassword;

    // ============================================
    // Redis Connection Factory
    // ============================================
    
    @Bean
    public LettuceConnectionFactory redisConnectionFactory() {
        RedisStandaloneConfiguration config = new RedisStandaloneConfiguration();
        config.setHostName(redisHost);
        config.setPort(redisPort);
        
        if (redisPassword != null && !redisPassword.isEmpty()) {
            config.setPassword(redisPassword);
        }
        
        // Lettuce client configuration with pooling
        LettucePoolingClientConfiguration clientConfig = LettucePoolingClientConfiguration.builder()
            .commandTimeout(Duration.ofSeconds(2))
            .build();
        
        return new LettuceConnectionFactory(config, clientConfig);
    }

    // ============================================
    // Object Mapper for JSON Serialization
    // ============================================
    
    @Bean
    public ObjectMapper redisObjectMapper() {
        ObjectMapper mapper = new ObjectMapper();
        
        // Support for Java 8 date/time types
        mapper.registerModule(new JavaTimeModule());
        mapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
        
        // Enable type information for polymorphic deserialization
        mapper.activateDefaultTyping(
            LaissezFaireSubTypeValidator.instance,
            ObjectMapper.DefaultTyping.NON_FINAL,
            JsonTypeInfo.As.PROPERTY
        );
        
        return mapper;
    }

    // ============================================
    // Redis Cache Manager with Custom TTL per Cache
    // ============================================
    
    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory connectionFactory) {
        // JSON serializer for cache values
        GenericJackson2JsonRedisSerializer jsonSerializer = 
            new GenericJackson2JsonRedisSerializer(redisObjectMapper());
        
        // Default cache configuration
        RedisCacheConfiguration defaultConfig = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofHours(1))                    // Default TTL: 1 hour
            .serializeKeysWith(RedisSerializationContext.SerializationPair
                .fromSerializer(new StringRedisSerializer()))
            .serializeValuesWith(RedisSerializationContext.SerializationPair
                .fromSerializer(jsonSerializer))
            .prefixCacheNameWith("myapp:cache:")              // Key prefix
            .disableCachingNullValues();                      // Don't cache nulls
        
        // Per-cache TTL configuration
        Map<String, RedisCacheConfiguration> cacheConfigs = new HashMap<>();
        
        // Products cache: 30 minutes TTL
        cacheConfigs.put("products", defaultConfig.entryTtl(Duration.ofMinutes(30)));
        
        // Users cache: 1 hour TTL
        cacheConfigs.put("users", defaultConfig.entryTtl(Duration.ofHours(1)));
        
        // Categories cache: 24 hours TTL (rarely changes)
        cacheConfigs.put("categories", defaultConfig.entryTtl(Duration.ofHours(24)));
        
        // Search results cache: 5 minutes TTL (more dynamic)
        cacheConfigs.put("searchResults", defaultConfig.entryTtl(Duration.ofMinutes(5)));
        
        // Sessions cache: 30 minutes TTL, allow null values
        cacheConfigs.put("sessions", defaultConfig
            .entryTtl(Duration.ofMinutes(30))
            // Uncomment below if you need to cache null values
            // Don't call .disableCachingNullValues()
        );
        
        return RedisCacheManager.builder(connectionFactory)
            .cacheDefaults(defaultConfig)
            .withInitialCacheConfigurations(cacheConfigs)
            .transactionAware()     // Sync cache operations with transactions
            .build();
    }

    // ============================================
    // RedisTemplate for Direct Redis Operations
    // ============================================
    
    @Bean
    public RedisTemplate<String, Object> redisTemplate(RedisConnectionFactory connectionFactory) {
        RedisTemplate<String, Object> template = new RedisTemplate<>();
        template.setConnectionFactory(connectionFactory);
        
        // Key serializer: String
        template.setKeySerializer(new StringRedisSerializer());
        template.setHashKeySerializer(new StringRedisSerializer());
        
        // Value serializer: JSON
        GenericJackson2JsonRedisSerializer jsonSerializer = 
            new GenericJackson2JsonRedisSerializer(redisObjectMapper());
        template.setValueSerializer(jsonSerializer);
        template.setHashValueSerializer(jsonSerializer);
        
        template.afterPropertiesSet();
        return template;
    }

    // ============================================
    // StringRedisTemplate for String Operations
    // ============================================
    
    @Bean
    public StringRedisTemplate stringRedisTemplate(RedisConnectionFactory connectionFactory) {
        return new StringRedisTemplate(connectionFactory);
    }

    // ============================================
    // Custom Key Generator
    // ============================================
    
    @Bean("customKeyGenerator")
    public KeyGenerator customKeyGenerator() {
        return (target, method, params) -> {
            StringBuilder sb = new StringBuilder();
            sb.append(target.getClass().getSimpleName());
            sb.append(":");
            sb.append(method.getName());
            for (Object param : params) {
                sb.append(":");
                sb.append(param != null ? param.toString() : "null");
            }
            return sb.toString();
        };
    }

    // ============================================
    // Cache Error Handler (Graceful Degradation)
    // ============================================
    
    @Override
    public CacheErrorHandler errorHandler() {
        return new RedisCacheErrorHandler();
    }
}

// Separate class for error handling
@Slf4j
class RedisCacheErrorHandler implements CacheErrorHandler {
    
    @Override
    public void handleCacheGetError(RuntimeException exception, Cache cache, Object key) {
        log.warn("Redis GET failed - cache: {}, key: {}, error: {}", 
            cache.getName(), key, exception.getMessage());
        // Don't rethrow - allow fallback to database
    }
    
    @Override
    public void handleCachePutError(RuntimeException exception, Cache cache, Object key, Object value) {
        log.warn("Redis PUT failed - cache: {}, key: {}, error: {}", 
            cache.getName(), key, exception.getMessage());
    }
    
    @Override
    public void handleCacheEvictError(RuntimeException exception, Cache cache, Object key) {
        log.warn("Redis EVICT failed - cache: {}, key: {}, error: {}", 
            cache.getName(), key, exception.getMessage());
    }
    
    @Override
    public void handleCacheClearError(RuntimeException exception, Cache cache) {
        log.error("Redis CLEAR failed - cache: {}, error: {}", 
            cache.getName(), exception.getMessage());
    }
}
```

### Serialization: JSON vs JDK

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SERIALIZATION COMPARISON                                      │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  ASPECT              │ JDK Serialization    │ JSON Serialization        │   │
│  │  ────────────────────┼──────────────────────┼────────────────────────── │   │
│  │  Performance         │ Faster               │ Slightly slower           │   │
│  │  Size                │ Larger               │ Smaller, readable        │   │
│  │  Compatibility       │ Java only            │ Language agnostic        │   │
│  │  Debugging           │ Hard (binary)        │ Easy (human readable)    │   │
│  │  Class Changes       │ Breaks easily        │ More resilient           │   │
│  │  Setup               │ Simple               │ Requires ObjectMapper    │   │
│  │  Recommendation      │ Avoid                │ Use this ✓               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  JSON IN REDIS (Readable):                                                       │
│  ─────────────────────────                                                       │
│  myapp:cache:products::1001                                                     │
│  {                                                                               │
│    "@class": "com.example.dto.ProductDTO",                                      │
│    "id": 1001,                                                                   │
│    "name": "Laptop",                                                            │
│    "price": 999.99,                                                             │
│    "createdAt": "2024-01-15T10:30:00"                                          │
│  }                                                                               │
│                                                                                  │
│  JDK IN REDIS (Binary blob - not readable):                                     │
│  ───────────────────────────────────────────                                     │
│  myapp:cache:products::1001                                                     │
│  \xac\xed\x00\x05sr\x00\x1ecom.example.dto.ProductDTO...                       │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Step 4: Using Caching Annotations

#### Service Layer with @Cacheable, @CachePut, @CacheEvict

```java
// ProductService.java - Complete Example with All Annotations
package com.example.service;

import com.example.dto.ProductDTO;
import com.example.dto.ProductRequest;
import com.example.entity.Product;
import com.example.exception.ProductNotFoundException;
import com.example.repository.ProductRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.cache.annotation.*;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
@CacheConfig(cacheNames = "products")  // Default cache name for this class
public class ProductService {

    private final ProductRepository repository;

    // ============================================
    // @Cacheable - Cache the result
    // ============================================
    
    /**
     * Get single product - cached by ID
     * Key: myapp:cache:products::1001
     */
    @Cacheable(key = "#id")
    public ProductDTO getProduct(Long id) {
        log.info("Fetching product from database: {}", id);
        return repository.findById(id)
            .map(this::toDTO)
            .orElseThrow(() -> new ProductNotFoundException(id));
    }
    
    /**
     * Get product by SKU - custom key
     * Key: myapp:cache:products::sku:ABC-123
     */
    @Cacheable(key = "'sku:' + #sku")
    public ProductDTO getProductBySku(String sku) {
        log.info("Fetching product by SKU from database: {}", sku);
        return repository.findBySku(sku)
            .map(this::toDTO)
            .orElseThrow(() -> new ProductNotFoundException("SKU: " + sku));
    }
    
    /**
     * Get products by category - different cache
     * Key: myapp:cache:productsByCategory::electronics
     */
    @Cacheable(cacheNames = "productsByCategory", key = "#category")
    public List<ProductDTO> getProductsByCategory(String category) {
        log.info("Fetching products by category from database: {}", category);
        return repository.findByCategory(category).stream()
            .map(this::toDTO)
            .collect(Collectors.toList());
    }
    
    /**
     * Conditional caching - only cache if price > 100
     */
    @Cacheable(
        key = "#id",
        condition = "#id > 0",                    // Only cache if ID is positive
        unless = "#result.price < 100"            // Don't cache cheap products
    )
    public ProductDTO getExpensiveProduct(Long id) {
        log.info("Fetching expensive product: {}", id);
        return repository.findById(id)
            .map(this::toDTO)
            .orElseThrow(() -> new ProductNotFoundException(id));
    }
    
    /**
     * Using custom key generator
     */
    @Cacheable(keyGenerator = "customKeyGenerator")
    public ProductDTO getProductWithCustomKey(Long id, String locale) {
        log.info("Fetching product {} for locale {}", id, locale);
        return repository.findById(id)
            .map(p -> toDTO(p, locale))
            .orElseThrow(() -> new ProductNotFoundException(id));
    }
    
    /**
     * Sync = true to prevent cache stampede
     * Only one thread computes value, others wait
     */
    @Cacheable(key = "#id", sync = true)
    public ProductDTO getHotProduct(Long id) {
        log.info("Fetching hot product (sync): {}", id);
        return repository.findById(id)
            .map(this::toDTO)
            .orElseThrow(() -> new ProductNotFoundException(id));
    }

    // ============================================
    // @CachePut - Update cache after operation
    // ============================================
    
    /**
     * Create product and add to cache
     * Key derived from result.id
     */
    @CachePut(key = "#result.id")
    @Transactional
    public ProductDTO createProduct(ProductRequest request) {
        log.info("Creating product: {}", request.getName());
        
        Product product = Product.builder()
            .name(request.getName())
            .sku(request.getSku())
            .price(request.getPrice())
            .category(request.getCategory())
            .build();
        
        Product saved = repository.save(product);
        return toDTO(saved);
    }
    
    /**
     * Update product and refresh cache
     */
    @CachePut(key = "#id")
    @Transactional
    public ProductDTO updateProduct(Long id, ProductRequest request) {
        log.info("Updating product: {}", id);
        
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        
        product.setName(request.getName());
        product.setSku(request.getSku());
        product.setPrice(request.getPrice());
        product.setCategory(request.getCategory());
        
        Product saved = repository.save(product);
        return toDTO(saved);
    }

    // ============================================
    // @CacheEvict - Remove from cache
    // ============================================
    
    /**
     * Delete product and evict from cache
     */
    @CacheEvict(key = "#id")
    @Transactional
    public void deleteProduct(Long id) {
        log.info("Deleting product: {}", id);
        
        if (!repository.existsById(id)) {
            throw new ProductNotFoundException(id);
        }
        
        repository.deleteById(id);
    }
    
    /**
     * Evict multiple caches on update
     */
    @Caching(evict = {
        @CacheEvict(value = "products", key = "#id"),
        @CacheEvict(value = "productsByCategory", key = "#category"),
        @CacheEvict(value = "searchResults", allEntries = true)
    })
    @Transactional
    public ProductDTO updateProductCategory(Long id, String category) {
        log.info("Updating product {} category to {}", id, category);
        
        Product product = repository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        
        product.setCategory(category);
        Product saved = repository.save(product);
        return toDTO(saved);
    }
    
    /**
     * Clear all entries in a cache
     */
    @CacheEvict(allEntries = true)
    public void clearProductCache() {
        log.info("Cleared all product cache entries");
    }
    
    /**
     * Clear cache before method execution
     */
    @CacheEvict(allEntries = true, beforeInvocation = true)
    public void refreshAllProducts() {
        log.info("Cache cleared before refresh");
        // Trigger reload logic
    }

    // ============================================
    // @Caching - Multiple cache operations
    // ============================================
    
    /**
     * Complex caching: cache result + evict related caches
     */
    @Caching(
        put = {
            @CachePut(value = "products", key = "#result.id"),
            @CachePut(value = "products", key = "'sku:' + #result.sku")
        },
        evict = {
            @CacheEvict(value = "productsByCategory", key = "#request.category")
        }
    )
    @Transactional
    public ProductDTO createAndCacheProduct(ProductRequest request) {
        log.info("Creating product with multi-cache: {}", request.getName());
        
        Product product = Product.builder()
            .name(request.getName())
            .sku(request.getSku())
            .price(request.getPrice())
            .category(request.getCategory())
            .build();
        
        Product saved = repository.save(product);
        return toDTO(saved);
    }

    // ============================================
    // Helper methods
    // ============================================
    
    private ProductDTO toDTO(Product product) {
        return ProductDTO.builder()
            .id(product.getId())
            .name(product.getName())
            .sku(product.getSku())
            .price(product.getPrice())
            .category(product.getCategory())
            .createdAt(product.getCreatedAt())
            .build();
    }
    
    private ProductDTO toDTO(Product product, String locale) {
        // Localized DTO conversion
        return toDTO(product);
    }
}
```

#### Repository Layer

```java
// ProductRepository.java
package com.example.repository;

import com.example.entity.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {
    
    Optional<Product> findBySku(String sku);
    
    List<Product> findByCategory(String category);
    
    @Query("SELECT p FROM Product p WHERE p.price BETWEEN :minPrice AND :maxPrice")
    List<Product> findByPriceRange(double minPrice, double maxPrice);
    
    List<Product> findByNameContainingIgnoreCase(String name);
}
```

#### REST Controller

```java
// ProductController.java
package com.example.controller;

import com.example.dto.ProductDTO;
import com.example.dto.ProductRequest;
import com.example.service.ProductService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import jakarta.validation.Valid;
import java.util.List;

@RestController
@RequestMapping("/api/v1/products")
@RequiredArgsConstructor
public class ProductController {

    private final ProductService productService;

    // GET /api/v1/products/1001
    // First call: Cache MISS → fetches from DB, stores in Redis
    // Subsequent calls: Cache HIT → returns from Redis
    @GetMapping("/{id}")
    public ResponseEntity<ProductDTO> getProduct(@PathVariable Long id) {
        return ResponseEntity.ok(productService.getProduct(id));
    }

    // GET /api/v1/products/sku/ABC-123
    @GetMapping("/sku/{sku}")
    public ResponseEntity<ProductDTO> getProductBySku(@PathVariable String sku) {
        return ResponseEntity.ok(productService.getProductBySku(sku));
    }

    // GET /api/v1/products/category/electronics
    @GetMapping("/category/{category}")
    public ResponseEntity<List<ProductDTO>> getProductsByCategory(
            @PathVariable String category) {
        return ResponseEntity.ok(productService.getProductsByCategory(category));
    }

    // POST /api/v1/products
    // Creates product in DB and caches the result
    @PostMapping
    public ResponseEntity<ProductDTO> createProduct(
            @Valid @RequestBody ProductRequest request) {
        ProductDTO created = productService.createProduct(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(created);
    }

    // PUT /api/v1/products/1001
    // Updates product in DB and refreshes cache
    @PutMapping("/{id}")
    public ResponseEntity<ProductDTO> updateProduct(
            @PathVariable Long id,
            @Valid @RequestBody ProductRequest request) {
        return ResponseEntity.ok(productService.updateProduct(id, request));
    }

    // DELETE /api/v1/products/1001
    // Deletes from DB and evicts from cache
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteProduct(@PathVariable Long id) {
        productService.deleteProduct(id);
        return ResponseEntity.noContent().build();
    }

    // POST /api/v1/products/cache/clear
    // Admin endpoint to clear product cache
    @PostMapping("/cache/clear")
    public ResponseEntity<Void> clearCache() {
        productService.clearProductCache();
        return ResponseEntity.ok().build();
    }
}
```

#### DTOs and Entity

```java
// ProductDTO.java
package com.example.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProductDTO implements Serializable {
    private static final long serialVersionUID = 1L;
    
    private Long id;
    private String name;
    private String sku;
    private BigDecimal price;
    private String category;
    private LocalDateTime createdAt;
}

// ProductRequest.java
package com.example.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.Data;

import java.math.BigDecimal;

@Data
public class ProductRequest {
    @NotBlank(message = "Name is required")
    private String name;
    
    @NotBlank(message = "SKU is required")
    private String sku;
    
    @NotNull(message = "Price is required")
    @Positive(message = "Price must be positive")
    private BigDecimal price;
    
    @NotBlank(message = "Category is required")
    private String category;
}

// Product.java (Entity)
package com.example.entity;

import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "products")
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Product {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(nullable = false)
    private String name;
    
    @Column(unique = true, nullable = false)
    private String sku;
    
    @Column(precision = 10, scale = 2)
    private BigDecimal price;
    
    @Column(nullable = false)
    private String category;
    
    @CreationTimestamp
    private LocalDateTime createdAt;
    
    @UpdateTimestamp
    private LocalDateTime updatedAt;
}
```

---

## 11. Advanced Redis Caching Techniques

### 11.1 Distributed Caching Across Multiple Instances

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    DISTRIBUTED CACHING ARCHITECTURE                              │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  Load Balancer                                                           │   │
│  │       │                                                                   │   │
│  │       ├────────────┬────────────┬────────────┐                          │   │
│  │       ▼            ▼            ▼            ▼                           │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐                    │   │
│  │  │ App 1   │  │ App 2   │  │ App 3   │  │ App N   │                    │   │
│  │  │Instance │  │Instance │  │Instance │  │Instance │                    │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘                    │   │
│  │       │            │            │            │                           │   │
│  │       └────────────┴─────┬──────┴────────────┘                          │   │
│  │                          │                                               │   │
│  │                          ▼                                               │   │
│  │              ┌───────────────────────┐                                  │   │
│  │              │    Redis (Shared)     │                                  │   │
│  │              │    Distributed Cache  │                                  │   │
│  │              └───────────────────────┘                                  │   │
│  │                                                                          │   │
│  │  All instances share the same cache:                                    │   │
│  │  • Consistent data across all instances                                 │   │
│  │  • Single source of truth for cached data                               │   │
│  │  • Cache operations visible to all instances immediately               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 11.2 Redis Cluster

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS CLUSTER ARCHITECTURE                                    │
│                                                                                  │
│  Redis Cluster provides:                                                        │
│  • Automatic data sharding across multiple nodes                                │
│  • High availability with automatic failover                                    │
│  • Linear scalability up to 1000 nodes                                         │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  Hash Slots: 0-16383 (divided among masters)                            │   │
│  │                                                                          │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐               │   │
│  │  │   Master 1    │  │   Master 2    │  │   Master 3    │               │   │
│  │  │ Slots: 0-5460 │  │Slots: 5461-   │  │Slots: 10923-  │               │   │
│  │  │               │  │      10922    │  │      16383    │               │   │
│  │  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘               │   │
│  │          │                  │                  │                        │   │
│  │          ▼                  ▼                  ▼                        │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐               │   │
│  │  │   Replica 1   │  │   Replica 2   │  │   Replica 3   │               │   │
│  │  │  (Standby)    │  │  (Standby)    │  │  (Standby)    │               │   │
│  │  └───────────────┘  └───────────────┘  └───────────────┘               │   │
│  │                                                                          │   │
│  │  Key routing: HASH_SLOT = CRC16(key) mod 16384                         │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```yaml
# application.yml - Redis Cluster Configuration
spring:
  redis:
    cluster:
      nodes:
        - redis-node-1:6379
        - redis-node-2:6379
        - redis-node-3:6379
        - redis-node-4:6379
        - redis-node-5:6379
        - redis-node-6:6379
      max-redirects: 3          # Follow redirects for MOVED/ASK responses
    password: ${REDIS_PASSWORD}
    
    lettuce:
      cluster:
        refresh:
          adaptive: true        # Enable adaptive refresh
          period: 30s           # Refresh cluster topology every 30s
      pool:
        max-active: 50
        max-idle: 20
```

```java
// Redis Cluster Configuration
@Configuration
public class RedisClusterConfig {
    
    @Value("${spring.redis.cluster.nodes}")
    private List<String> clusterNodes;
    
    @Value("${spring.redis.password:}")
    private String password;
    
    @Bean
    public LettuceConnectionFactory redisConnectionFactory() {
        RedisClusterConfiguration clusterConfig = new RedisClusterConfiguration(clusterNodes);
        
        if (password != null && !password.isEmpty()) {
            clusterConfig.setPassword(RedisPassword.of(password));
        }
        
        // Cluster-specific client options
        ClusterClientOptions clientOptions = ClusterClientOptions.builder()
            .autoReconnect(true)
            .validateClusterNodeMembership(false)
            .topologyRefreshOptions(ClusterTopologyRefreshOptions.builder()
                .enablePeriodicRefresh(Duration.ofSeconds(30))
                .enableAllAdaptiveRefreshTriggers()
                .build())
            .build();
        
        LettuceClientConfiguration clientConfig = LettucePoolingClientConfiguration.builder()
            .commandTimeout(Duration.ofSeconds(2))
            .clientOptions(clientOptions)
            .build();
        
        return new LettuceConnectionFactory(clusterConfig, clientConfig);
    }
}
```

### 11.3 Redis Sentinel (High Availability)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS SENTINEL ARCHITECTURE                                   │
│                                                                                  │
│  Sentinel provides:                                                             │
│  • Monitoring: Constantly checks master and replicas                            │
│  • Notification: Alerts when something goes wrong                               │
│  • Automatic failover: Promotes replica to master if master fails              │
│  • Configuration provider: Clients connect via Sentinel                        │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │        ┌────────────┐  ┌────────────┐  ┌────────────┐                  │   │
│  │        │ Sentinel 1 │  │ Sentinel 2 │  │ Sentinel 3 │                  │   │
│  │        └──────┬─────┘  └──────┬─────┘  └──────┬─────┘                  │   │
│  │               │               │               │                         │   │
│  │               └───────────────┼───────────────┘                         │   │
│  │                               │ Monitoring                              │   │
│  │                               ▼                                          │   │
│  │                      ┌───────────────┐                                  │   │
│  │                      │    Master     │                                  │   │
│  │                      │   (Primary)   │                                  │   │
│  │                      └───────┬───────┘                                  │   │
│  │                              │ Replication                              │   │
│  │                    ┌─────────┴─────────┐                               │   │
│  │                    ▼                   ▼                                │   │
│  │            ┌───────────────┐   ┌───────────────┐                       │   │
│  │            │   Replica 1   │   │   Replica 2   │                       │   │
│  │            │  (Standby)    │   │  (Standby)    │                       │   │
│  │            └───────────────┘   └───────────────┘                       │   │
│  │                                                                          │   │
│  │  Failover: If master fails, Sentinels elect new master from replicas   │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```yaml
# application.yml - Redis Sentinel Configuration
spring:
  redis:
    sentinel:
      master: mymaster                    # Sentinel master name
      nodes:
        - sentinel-1:26379
        - sentinel-2:26379
        - sentinel-3:26379
    password: ${REDIS_PASSWORD}
    database: 0
    
    lettuce:
      pool:
        max-active: 20
        max-idle: 10
```

```java
// Redis Sentinel Configuration
@Configuration
public class RedisSentinelConfig {
    
    @Value("${spring.redis.sentinel.master}")
    private String master;
    
    @Value("${spring.redis.sentinel.nodes}")
    private List<String> sentinelNodes;
    
    @Value("${spring.redis.password:}")
    private String password;
    
    @Bean
    public LettuceConnectionFactory redisConnectionFactory() {
        RedisSentinelConfiguration sentinelConfig = new RedisSentinelConfiguration()
            .master(master);
        
        for (String node : sentinelNodes) {
            String[] parts = node.split(":");
            sentinelConfig.sentinel(parts[0], Integer.parseInt(parts[1]));
        }
        
        if (password != null && !password.isEmpty()) {
            sentinelConfig.setPassword(RedisPassword.of(password));
        }
        
        LettuceClientConfiguration clientConfig = LettucePoolingClientConfiguration.builder()
            .commandTimeout(Duration.ofSeconds(2))
            .build();
        
        return new LettuceConnectionFactory(sentinelConfig, clientConfig);
    }
}
```

### 11.4 Handling Cache Stampede with Locks

```java
// Distributed lock to prevent cache stampede
@Service
@RequiredArgsConstructor
@Slf4j
public class CacheStampedePreventionService {
    
    private final ProductRepository repository;
    private final StringRedisTemplate redisTemplate;
    private final ObjectMapper objectMapper;
    
    private static final String LOCK_PREFIX = "lock:";
    private static final Duration LOCK_TIMEOUT = Duration.ofSeconds(10);
    private static final Duration CACHE_TTL = Duration.ofMinutes(30);
    
    /**
     * Get product with distributed lock to prevent stampede
     */
    public Product getProductWithLock(Long id) {
        String cacheKey = "product:" + id;
        String lockKey = LOCK_PREFIX + cacheKey;
        
        // 1. Try to get from cache
        String cached = redisTemplate.opsForValue().get(cacheKey);
        if (cached != null) {
            return deserialize(cached);
        }
        
        // 2. Cache miss - try to acquire lock
        String lockValue = UUID.randomUUID().toString();
        Boolean acquired = redisTemplate.opsForValue()
            .setIfAbsent(lockKey, lockValue, LOCK_TIMEOUT);
        
        if (Boolean.TRUE.equals(acquired)) {
            try {
                // 3. Double-check cache (another thread may have populated it)
                cached = redisTemplate.opsForValue().get(cacheKey);
                if (cached != null) {
                    return deserialize(cached);
                }
                
                // 4. Fetch from database
                log.info("Lock acquired, fetching from database: {}", id);
                Product product = repository.findById(id)
                    .orElseThrow(() -> new ProductNotFoundException(id));
                
                // 5. Populate cache
                redisTemplate.opsForValue().set(cacheKey, serialize(product), CACHE_TTL);
                
                return product;
            } finally {
                // 6. Release lock (only if we own it)
                releaseLock(lockKey, lockValue);
            }
        } else {
            // 7. Wait and retry
            log.debug("Lock not acquired, waiting...");
            try {
                Thread.sleep(100);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            return getProductWithLock(id);  // Retry
        }
    }
    
    /**
     * Release lock safely using Lua script (atomic operation)
     */
    private void releaseLock(String lockKey, String lockValue) {
        String script = 
            "if redis.call('get', KEYS[1]) == ARGV[1] then " +
            "    return redis.call('del', KEYS[1]) " +
            "else " +
            "    return 0 " +
            "end";
        
        redisTemplate.execute(
            new DefaultRedisScript<>(script, Long.class),
            Collections.singletonList(lockKey),
            lockValue
        );
    }
    
    private String serialize(Product product) {
        try {
            return objectMapper.writeValueAsString(product);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Serialization failed", e);
        }
    }
    
    private Product deserialize(String json) {
        try {
            return objectMapper.readValue(json, Product.class);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Deserialization failed", e);
        }
    }
}
```

### 11.5 Preventing Cache Penetration

```java
// Cache penetration prevention with Bloom Filter
@Service
@RequiredArgsConstructor
@Slf4j
public class CachePenetrationService {
    
    private final ProductRepository repository;
    private final StringRedisTemplate redisTemplate;
    private final ObjectMapper objectMapper;
    
    private static final String BLOOM_FILTER_KEY = "bloom:products";
    private static final String CACHE_PREFIX = "product:";
    private static final String NULL_VALUE = "NULL";
    private static final Duration NULL_TTL = Duration.ofMinutes(5);
    private static final Duration CACHE_TTL = Duration.ofMinutes(30);
    
    @PostConstruct
    public void initBloomFilter() {
        // Initialize Bloom filter with existing product IDs
        // Using Redis Bloom Module (if available) or application-level Bloom filter
        log.info("Initializing bloom filter with existing product IDs");
        repository.findAllIds().forEach(this::addToBloomFilter);
    }
    
    /**
     * Get product with null caching to prevent penetration
     */
    public Product getProductWithNullCaching(Long id) {
        String cacheKey = CACHE_PREFIX + id;
        
        String cached = redisTemplate.opsForValue().get(cacheKey);
        
        // Check for cached null
        if (NULL_VALUE.equals(cached)) {
            log.debug("Returning cached null for: {}", id);
            return null;  // Return null without hitting DB
        }
        
        if (cached != null) {
            return deserialize(cached);
        }
        
        // Cache miss - check bloom filter first
        if (!mightExist(id)) {
            log.debug("Bloom filter says {} doesn't exist, caching null", id);
            redisTemplate.opsForValue().set(cacheKey, NULL_VALUE, NULL_TTL);
            return null;
        }
        
        // Query database
        Product product = repository.findById(id).orElse(null);
        
        if (product == null) {
            // Cache the null to prevent repeated DB hits
            log.debug("Product {} not found, caching null", id);
            redisTemplate.opsForValue().set(cacheKey, NULL_VALUE, NULL_TTL);
            return null;
        }
        
        // Cache the product
        redisTemplate.opsForValue().set(cacheKey, serialize(product), CACHE_TTL);
        return product;
    }
    
    /**
     * Add to bloom filter when creating product
     */
    @CachePut(value = "products", key = "#result.id")
    @Transactional
    public Product createProduct(ProductRequest request) {
        Product product = repository.save(mapToEntity(request));
        addToBloomFilter(product.getId());
        return product;
    }
    
    // Bloom filter operations (using Redis BITSET as simple bloom filter)
    private void addToBloomFilter(Long id) {
        // Simple hash - in production, use multiple hash functions
        long hash1 = id % 10000;
        long hash2 = (id * 31) % 10000;
        long hash3 = (id * 37) % 10000;
        
        redisTemplate.opsForValue().setBit(BLOOM_FILTER_KEY, hash1, true);
        redisTemplate.opsForValue().setBit(BLOOM_FILTER_KEY, hash2, true);
        redisTemplate.opsForValue().setBit(BLOOM_FILTER_KEY, hash3, true);
    }
    
    private boolean mightExist(Long id) {
        long hash1 = id % 10000;
        long hash2 = (id * 31) % 10000;
        long hash3 = (id * 37) % 10000;
        
        Boolean bit1 = redisTemplate.opsForValue().getBit(BLOOM_FILTER_KEY, hash1);
        Boolean bit2 = redisTemplate.opsForValue().getBit(BLOOM_FILTER_KEY, hash2);
        Boolean bit3 = redisTemplate.opsForValue().getBit(BLOOM_FILTER_KEY, hash3);
        
        return Boolean.TRUE.equals(bit1) && Boolean.TRUE.equals(bit2) && Boolean.TRUE.equals(bit3);
    }
}
```

### 11.6 Handling Null Values

```java
// Null value handling configuration
@Configuration
public class NullValueCacheConfig {
    
    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory connectionFactory) {
        GenericJackson2JsonRedisSerializer serializer = 
            new GenericJackson2JsonRedisSerializer();
        
        // Configuration that ALLOWS null values
        RedisCacheConfiguration configWithNulls = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(30))
            .serializeValuesWith(RedisSerializationContext.SerializationPair
                .fromSerializer(serializer));
            // Note: not calling .disableCachingNullValues()
        
        // Configuration that DISALLOWS null values
        RedisCacheConfiguration configNoNulls = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(30))
            .serializeValuesWith(RedisSerializationContext.SerializationPair
                .fromSerializer(serializer))
            .disableCachingNullValues();
        
        Map<String, RedisCacheConfiguration> cacheConfigs = new HashMap<>();
        cacheConfigs.put("productsAllowNull", configWithNulls);  // Allow null
        cacheConfigs.put("usersNoNull", configNoNulls);           // Disallow null
        
        return RedisCacheManager.builder(connectionFactory)
            .cacheDefaults(configNoNulls)
            .withInitialCacheConfigurations(cacheConfigs)
            .build();
    }
}

// Service using null-aware caching
@Service
@Slf4j
public class NullAwareCacheService {
    
    private final UserRepository userRepository;
    
    // Using unless to handle null returns
    @Cacheable(
        value = "users",
        key = "#id",
        unless = "#result == null"  // Don't cache if result is null
    )
    public UserDTO getUser(Long id) {
        log.info("Fetching user: {}", id);
        return userRepository.findById(id)
            .map(this::toDTO)
            .orElse(null);
    }
    
    // Using cache that allows null values
    @Cacheable(
        value = "productsAllowNull",
        key = "#id"
    )
    public ProductDTO getProductAllowNull(Long id) {
        log.info("Fetching product (null allowed): {}", id);
        return productRepository.findById(id)
            .map(this::toDTO)
            .orElse(null);  // null WILL be cached
    }
}
```

### 11.7 Using Lua Scripts for Atomic Operations

```java
// Lua scripts for atomic cache operations
@Service
@RequiredArgsConstructor
@Slf4j
public class AtomicCacheOperationsService {
    
    private final StringRedisTemplate redisTemplate;
    
    /**
     * Increment and get with expiry (atomic)
     * Use case: Rate limiting, counters
     */
    public Long incrementWithExpiry(String key, long ttlSeconds) {
        String script = 
            "local current = redis.call('INCR', KEYS[1]) " +
            "if current == 1 then " +
            "    redis.call('EXPIRE', KEYS[1], ARGV[1]) " +
            "end " +
            "return current";
        
        DefaultRedisScript<Long> redisScript = new DefaultRedisScript<>();
        redisScript.setScriptText(script);
        redisScript.setResultType(Long.class);
        
        return redisTemplate.execute(
            redisScript,
            Collections.singletonList(key),
            String.valueOf(ttlSeconds)
        );
    }
    
    /**
     * Get and delete (atomic pop)
     * Use case: One-time tokens, queue processing
     */
    public String getAndDelete(String key) {
        String script = 
            "local value = redis.call('GET', KEYS[1]) " +
            "redis.call('DEL', KEYS[1]) " +
            "return value";
        
        DefaultRedisScript<String> redisScript = new DefaultRedisScript<>();
        redisScript.setScriptText(script);
        redisScript.setResultType(String.class);
        
        return redisTemplate.execute(redisScript, Collections.singletonList(key));
    }
    
    /**
     * Compare and set (CAS operation)
     * Use case: Optimistic locking
     */
    public Boolean compareAndSet(String key, String expectedValue, String newValue) {
        String script = 
            "if redis.call('GET', KEYS[1]) == ARGV[1] then " +
            "    redis.call('SET', KEYS[1], ARGV[2]) " +
            "    return 1 " +
            "else " +
            "    return 0 " +
            "end";
        
        DefaultRedisScript<Long> redisScript = new DefaultRedisScript<>();
        redisScript.setScriptText(script);
        redisScript.setResultType(Long.class);
        
        Long result = redisTemplate.execute(
            redisScript,
            Collections.singletonList(key),
            expectedValue,
            newValue
        );
        
        return result != null && result == 1L;
    }
    
    /**
     * Sliding window rate limiter
     * Use case: API rate limiting
     */
    public boolean isAllowed(String key, int maxRequests, int windowSeconds) {
        String script = 
            "local current_time = tonumber(ARGV[1]) " +
            "local window = tonumber(ARGV[2]) " +
            "local max_requests = tonumber(ARGV[3]) " +
            "local window_start = current_time - window " +
            "" +
            "redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', window_start) " +
            "local current_count = redis.call('ZCARD', KEYS[1]) " +
            "" +
            "if current_count < max_requests then " +
            "    redis.call('ZADD', KEYS[1], current_time, current_time .. '-' .. math.random()) " +
            "    redis.call('EXPIRE', KEYS[1], window) " +
            "    return 1 " +
            "else " +
            "    return 0 " +
            "end";
        
        DefaultRedisScript<Long> redisScript = new DefaultRedisScript<>();
        redisScript.setScriptText(script);
        redisScript.setResultType(Long.class);
        
        Long result = redisTemplate.execute(
            redisScript,
            Collections.singletonList(key),
            String.valueOf(System.currentTimeMillis()),
            String.valueOf(windowSeconds * 1000),
            String.valueOf(maxRequests)
        );
        
        return result != null && result == 1L;
    }
}

// Usage example - Rate Limiting
@RestController
@RequiredArgsConstructor
public class RateLimitedController {
    
    private final AtomicCacheOperationsService cacheOps;
    private final ProductService productService;
    
    @GetMapping("/api/products/{id}")
    public ResponseEntity<?> getProduct(
            @PathVariable Long id,
            HttpServletRequest request) {
        
        String clientIp = request.getRemoteAddr();
        String rateLimitKey = "ratelimit:" + clientIp + ":products";
        
        // Allow 100 requests per minute
        if (!cacheOps.isAllowed(rateLimitKey, 100, 60)) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                .body("Rate limit exceeded. Try again later.");
        }
        
        return ResponseEntity.ok(productService.getProduct(id));
    }
}
```

---

## 12. Performance Considerations

### 12.1 Memory Optimization

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS MEMORY OPTIMIZATION                                     │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  STRATEGY                  │ DESCRIPTION                                │   │
│  │  ──────────────────────────┼────────────────────────────────────────── │   │
│  │  Use appropriate data types│ Hash for objects (more memory efficient)  │   │
│  │  Short key names          │ "p:1001" vs "product:1001" (saves bytes)  │   │
│  │  Compress large values    │ gzip JSON before storing                   │   │
│  │  Set maxmemory            │ Prevent Redis from using all RAM          │   │
│  │  Use eviction policy      │ Auto-remove old data when full            │   │
│  │  Avoid large keys         │ Split into smaller chunks if > 1MB        │   │
│  │  Use EXPIRE/TTL           │ Auto-cleanup unused data                   │   │
│  │  Monitor memory usage     │ INFO memory, MEMORY USAGE key             │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  MEMORY-EFFICIENT DATA STRUCTURES:                                               │
│  ──────────────────────────────────                                              │
│                                                                                  │
│  Instead of:                                                                     │
│    SET product:1001:name "Laptop"                                               │
│    SET product:1001:price "999.99"                                              │
│    SET product:1001:category "electronics"                                      │
│                                                                                  │
│  Use Hash:                                                                       │
│    HSET product:1001 name "Laptop" price "999.99" category "electronics"       │
│                                                                                  │
│  Memory savings: Up to 10x for small objects!                                   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Memory-efficient caching with Hash
@Service
@RequiredArgsConstructor
public class MemoryEfficientCacheService {
    
    private final RedisTemplate<String, Object> redisTemplate;
    private final HashOperations<String, String, String> hashOps;
    
    @PostConstruct
    public void init() {
        hashOps = redisTemplate.opsForHash();
    }
    
    // Store object as Hash (memory efficient)
    public void cacheProductAsHash(Product product) {
        String key = "product:" + product.getId();
        
        Map<String, String> fields = new HashMap<>();
        fields.put("name", product.getName());
        fields.put("price", product.getPrice().toString());
        fields.put("category", product.getCategory());
        fields.put("sku", product.getSku());
        
        hashOps.putAll(key, fields);
        redisTemplate.expire(key, Duration.ofMinutes(30));
    }
    
    // Retrieve object from Hash
    public Product getProductFromHash(Long id) {
        String key = "product:" + id;
        
        Map<String, String> fields = hashOps.entries(key);
        
        if (fields.isEmpty()) {
            return null;
        }
        
        return Product.builder()
            .id(id)
            .name(fields.get("name"))
            .price(new BigDecimal(fields.get("price")))
            .category(fields.get("category"))
            .sku(fields.get("sku"))
            .build();
    }
    
    // Get specific field (partial read)
    public BigDecimal getProductPrice(Long id) {
        String key = "product:" + id;
        String price = hashOps.get(key, "price");
        return price != null ? new BigDecimal(price) : null;
    }
    
    // Update specific field (partial update)
    public void updateProductPrice(Long id, BigDecimal newPrice) {
        String key = "product:" + id;
        hashOps.put(key, "price", newPrice.toString());
    }
}
```

### 12.2 Data Serialization Strategies

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SERIALIZATION COMPARISON                                      │
│                                                                                  │
│  │ Serializer           │ Size   │ Speed  │ Readability │ Cross-Lang │         │
│  │──────────────────────┼────────┼────────┼─────────────┼────────────│         │
│  │ JDK Serialization    │ Large  │ Fast   │ None        │ No         │         │
│  │ Jackson JSON         │ Medium │ Medium │ High        │ Yes        │ ✓       │
│  │ Kryo                 │ Small  │ Fast   │ None        │ No         │         │
│  │ Protocol Buffers     │ Small  │ Fast   │ Medium      │ Yes        │         │
│  │ MessagePack          │ Small  │ Fast   │ Medium      │ Yes        │         │
│                                                                                  │
│  RECOMMENDATION: Use Jackson JSON for most use cases (balance of features)      │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Custom serializers configuration
@Configuration
public class SerializerConfig {
    
    // Jackson JSON serializer (recommended)
    @Bean
    public RedisSerializer<Object> jsonSerializer() {
        ObjectMapper mapper = new ObjectMapper();
        mapper.registerModule(new JavaTimeModule());
        mapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
        mapper.activateDefaultTyping(
            LaissezFaireSubTypeValidator.instance,
            ObjectMapper.DefaultTyping.NON_FINAL
        );
        return new GenericJackson2JsonRedisSerializer(mapper);
    }
    
    // Kryo serializer (for performance-critical scenarios)
    @Bean
    public RedisSerializer<Object> kryoSerializer() {
        return new RedisSerializer<Object>() {
            private final ThreadLocal<Kryo> kryoLocal = ThreadLocal.withInitial(() -> {
                Kryo kryo = new Kryo();
                kryo.register(Product.class);
                kryo.register(ProductDTO.class);
                kryo.register(ArrayList.class);
                return kryo;
            });
            
            @Override
            public byte[] serialize(Object o) throws SerializationException {
                if (o == null) return new byte[0];
                ByteArrayOutputStream baos = new ByteArrayOutputStream();
                Output output = new Output(baos);
                kryoLocal.get().writeClassAndObject(output, o);
                output.close();
                return baos.toByteArray();
            }
            
            @Override
            public Object deserialize(byte[] bytes) throws SerializationException {
                if (bytes == null || bytes.length == 0) return null;
                Input input = new Input(bytes);
                return kryoLocal.get().readClassAndObject(input);
            }
        };
    }
    
    // Compression for large objects
    @Bean
    public RedisSerializer<Object> compressingSerializer() {
        return new RedisSerializer<Object>() {
            private final ObjectMapper mapper = new ObjectMapper();
            
            @Override
            public byte[] serialize(Object o) throws SerializationException {
                if (o == null) return new byte[0];
                try {
                    byte[] json = mapper.writeValueAsBytes(o);
                    // Compress if larger than 1KB
                    if (json.length > 1024) {
                        ByteArrayOutputStream baos = new ByteArrayOutputStream();
                        GZIPOutputStream gzip = new GZIPOutputStream(baos);
                        gzip.write(json);
                        gzip.close();
                        return baos.toByteArray();
                    }
                    return json;
                } catch (IOException e) {
                    throw new SerializationException("Compression failed", e);
                }
            }
            
            @Override
            public Object deserialize(byte[] bytes) throws SerializationException {
                if (bytes == null || bytes.length == 0) return null;
                try {
                    // Check if compressed (gzip magic number)
                    if (bytes.length > 2 && bytes[0] == (byte) 0x1f && bytes[1] == (byte) 0x8b) {
                        GZIPInputStream gzip = new GZIPInputStream(new ByteArrayInputStream(bytes));
                        bytes = gzip.readAllBytes();
                        gzip.close();
                    }
                    return mapper.readValue(bytes, Object.class);
                } catch (IOException e) {
                    throw new SerializationException("Decompression failed", e);
                }
            }
        };
    }
}
```

### 12.3 Choosing Eviction Policy

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    EVICTION POLICY SELECTION GUIDE                               │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  USE CASE                        │ RECOMMENDED POLICY                   │   │
│  │  ────────────────────────────────┼─────────────────────────────────────│   │
│  │  General caching                 │ allkeys-lru ✓                       │   │
│  │  Frequency-based access pattern  │ allkeys-lfu                         │   │
│  │  Mix of cache + persistent data  │ volatile-lru                        │   │
│  │  Cache with time-based priority  │ volatile-ttl                        │   │
│  │  Redis as database (no eviction) │ noeviction                          │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  DECISION FLOWCHART:                                                             │
│                                                                                  │
│                     ┌─────────────────────────┐                                 │
│                     │ Does all data have TTL? │                                 │
│                     └───────────┬─────────────┘                                 │
│                           │          │                                           │
│                          YES         NO                                          │
│                           │          │                                           │
│                           ▼          ▼                                           │
│         ┌─────────────────────┐  ┌─────────────────────┐                       │
│         │Priority by TTL      │  │Some data is         │                       │
│         │remaining?           │  │persistent?          │                       │
│         └──────────┬──────────┘  └──────────┬──────────┘                       │
│              │         │              │          │                              │
│             YES        NO            YES         NO                             │
│              │         │              │          │                              │
│              ▼         ▼              ▼          ▼                              │
│        ┌─────────┐┌─────────┐  ┌─────────┐ ┌─────────┐                        │
│        │volatile-││volatile-│  │volatile-│ │allkeys- │                        │
│        │ttl      ││lru/lfu  │  │lru/lfu  │ │lru/lfu  │                        │
│        └─────────┘└─────────┘  └─────────┘ └─────────┘                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 12.4 Monitoring Redis

```java
// Redis monitoring endpoint
@RestController
@RequestMapping("/api/admin/redis")
@RequiredArgsConstructor
public class RedisMonitoringController {
    
    private final StringRedisTemplate redisTemplate;
    
    @GetMapping("/info")
    public Map<String, Object> getRedisInfo() {
        Properties info = redisTemplate.getConnectionFactory()
            .getConnection()
            .info();
        
        Map<String, Object> result = new HashMap<>();
        
        // Memory info
        result.put("used_memory_human", info.getProperty("used_memory_human"));
        result.put("used_memory_peak_human", info.getProperty("used_memory_peak_human"));
        result.put("maxmemory_human", info.getProperty("maxmemory_human"));
        
        // Stats
        result.put("total_connections_received", info.getProperty("total_connections_received"));
        result.put("total_commands_processed", info.getProperty("total_commands_processed"));
        result.put("keyspace_hits", info.getProperty("keyspace_hits"));
        result.put("keyspace_misses", info.getProperty("keyspace_misses"));
        
        // Calculate hit rate
        long hits = Long.parseLong(info.getProperty("keyspace_hits", "0"));
        long misses = Long.parseLong(info.getProperty("keyspace_misses", "0"));
        double hitRate = hits + misses > 0 ? (double) hits / (hits + misses) * 100 : 0;
        result.put("hit_rate_percent", String.format("%.2f%%", hitRate));
        
        // Connected clients
        result.put("connected_clients", info.getProperty("connected_clients"));
        
        // Keys
        result.put("db0_keys", info.getProperty("db0"));
        
        return result;
    }
    
    @GetMapping("/memory/{key}")
    public Map<String, Object> getKeyMemory(@PathVariable String key) {
        Long memoryUsage = redisTemplate.execute((RedisCallback<Long>) connection -> 
            connection.serverCommands().objectEncoding(key.getBytes()));
        
        Long ttl = redisTemplate.getExpire(key);
        String type = redisTemplate.type(key).code();
        
        Map<String, Object> result = new HashMap<>();
        result.put("key", key);
        result.put("type", type);
        result.put("ttl_seconds", ttl);
        result.put("exists", redisTemplate.hasKey(key));
        
        return result;
    }
    
    @GetMapping("/slowlog")
    public List<Map<String, Object>> getSlowLog() {
        return redisTemplate.execute((RedisCallback<List<Map<String, Object>>>) connection -> {
            List<Object> slowlog = connection.serverCommands().slowLogGet(10);
            // Parse slowlog entries
            return Collections.emptyList(); // Simplified
        });
    }
}
```

```yaml
# application.yml - Redis monitoring config
management:
  endpoints:
    web:
      exposure:
        include: health, metrics, redis
  health:
    redis:
      enabled: true

# Micrometer metrics for Redis
spring:
  redis:
    lettuce:
      pool:
        enabled: true
```

### 12.5 Connection Pooling: Lettuce vs Jedis

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    LETTUCE VS JEDIS                                              │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  ASPECT               │ LETTUCE (Default)    │ JEDIS                   │   │
│  │  ─────────────────────┼──────────────────────┼───────────────────────  │   │
│  │  Threading            │ Non-blocking (Netty) │ Blocking                │   │
│  │  Connection sharing   │ Single connection OK │ Needs connection pool  │   │
│  │  Cluster support      │ Excellent            │ Good                    │   │
│  │  Reactive support     │ Yes                  │ No                      │   │
│  │  Memory footprint     │ Higher               │ Lower                   │   │
│  │  Maturity             │ Newer                │ Older, stable           │   │
│  │  Spring Boot default  │ Yes ✓               │ No                      │   │
│  │                                                                          │   │
│  │  RECOMMENDATION:                                                        │   │
│  │  • Use Lettuce (default) for most applications                         │   │
│  │  • Use Jedis if you need simple, blocking operations                   │   │
│  │  • Use Lettuce for reactive applications                               │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```xml
<!-- Use Jedis instead of Lettuce -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-redis</artifactId>
    <exclusions>
        <exclusion>
            <groupId>io.lettuce</groupId>
            <artifactId>lettuce-core</artifactId>
        </exclusion>
    </exclusions>
</dependency>
<dependency>
    <groupId>redis.clients</groupId>
    <artifactId>jedis</artifactId>
</dependency>
```

```java
// Jedis connection pool configuration
@Configuration
public class JedisConfig {
    
    @Bean
    public JedisConnectionFactory jedisConnectionFactory() {
        RedisStandaloneConfiguration config = new RedisStandaloneConfiguration();
        config.setHostName("localhost");
        config.setPort(6379);
        
        JedisPoolConfig poolConfig = new JedisPoolConfig();
        poolConfig.setMaxTotal(50);         // Max connections
        poolConfig.setMaxIdle(20);          // Max idle connections
        poolConfig.setMinIdle(5);           // Min idle connections
        poolConfig.setMaxWaitMillis(2000);  // Max wait for connection
        poolConfig.setTestOnBorrow(true);   // Test connection before use
        poolConfig.setTestWhileIdle(true);  // Test idle connections
        
        JedisClientConfiguration clientConfig = JedisClientConfiguration.builder()
            .usePooling()
            .poolConfig(poolConfig)
            .build();
        
        return new JedisConnectionFactory(config, clientConfig);
    }
}
```

### 12.6 Horizontal Scaling

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS HORIZONTAL SCALING OPTIONS                              │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  1. READ REPLICAS (Read Scaling)                                        │   │
│  │  ────────────────────────────────                                        │   │
│  │                                                                          │   │
│  │      ┌──────────────┐                                                   │   │
│  │      │    Master    │ ← All writes                                      │   │
│  │      └──────┬───────┘                                                   │   │
│  │             │ Replication                                               │   │
│  │       ┌─────┴─────┐                                                     │   │
│  │       ▼           ▼                                                      │   │
│  │  ┌─────────┐ ┌─────────┐                                               │   │
│  │  │Replica 1│ │Replica 2│ ← Reads distributed                           │   │
│  │  └─────────┘ └─────────┘                                               │   │
│  │                                                                          │   │
│  │  2. REDIS CLUSTER (Write + Read Scaling)                                │   │
│  │  ────────────────────────────────────────                                │   │
│  │                                                                          │   │
│  │  Data sharded across multiple masters:                                  │   │
│  │  • Each master handles subset of keys                                   │   │
│  │  • Linear scalability                                                   │   │
│  │  • Automatic failover                                                   │   │
│  │                                                                          │   │
│  │  3. CLIENT-SIDE PARTITIONING                                            │   │
│  │  ───────────────────────────────                                         │   │
│  │                                                                          │   │
│  │  Application routes keys to specific Redis instances:                   │   │
│  │  • Simple to implement                                                  │   │
│  │  • No automatic failover                                                │   │
│  │  • Manual rebalancing needed                                            │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Read from replica configuration
@Configuration
public class ReadReplicaConfig {
    
    @Bean
    public LettuceConnectionFactory redisConnectionFactory() {
        RedisStandaloneConfiguration masterConfig = new RedisStandaloneConfiguration();
        masterConfig.setHostName("redis-master");
        masterConfig.setPort(6379);
        
        // Configure read from replica
        LettuceClientConfiguration clientConfig = LettuceClientConfiguration.builder()
            .readFrom(ReadFrom.REPLICA_PREFERRED)  // Prefer replicas for reads
            .commandTimeout(Duration.ofSeconds(2))
            .build();
        
        return new LettuceConnectionFactory(masterConfig, clientConfig);
    }
}

// ReadFrom options:
// ReadFrom.MASTER           - Always read from master
// ReadFrom.REPLICA          - Always read from replica
// ReadFrom.REPLICA_PREFERRED - Prefer replica, fallback to master
// ReadFrom.NEAREST          - Read from nearest node (lowest latency)
```

---

## 13. When to Use Redis for Caching

### 13.1 Ideal Scenarios for Redis Caching

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS CACHING USE CASES                                       │
│                                                                                  │
│  ✅ SCENARIO 1: MICROSERVICES ARCHITECTURE                                      │
│  ─────────────────────────────────────────                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Multiple services need shared cache                                   │   │
│  │  • Service instances scale independently                                 │   │
│  │  • Consistent cache across all instances                                │   │
│  │  • Session sharing between services                                     │   │
│  │                                                                          │   │
│  │  [Service A (3 instances)] ──┐                                          │   │
│  │  [Service B (5 instances)] ──┼──→ [Redis Cache]                        │   │
│  │  [Service C (2 instances)] ──┘                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ✅ SCENARIO 2: HIGH READ-HEAVY SYSTEMS                                         │
│  ──────────────────────────────────────                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Read:Write ratio > 10:1                                              │   │
│  │  • Product catalogs, user profiles                                      │   │
│  │  • Configuration data                                                   │   │
│  │                                                                          │   │
│  │  Without cache: 10,000 reads/sec → Database overloaded                 │   │
│  │  With Redis:    10,000 reads/sec → 95% from cache, DB relaxed          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ✅ SCENARIO 3: API RESPONSE CACHING                                            │
│  ─────────────────────────────────────                                           │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • External API calls (rate limited, slow)                              │   │
│  │  • Expensive computations                                               │   │
│  │  • Aggregated data                                                      │   │
│  │                                                                          │   │
│  │  Example: Weather API                                                   │   │
│  │    Without cache: $500/month API costs, 500ms latency                  │   │
│  │    With cache:    $50/month API costs, 5ms latency                     │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ✅ SCENARIO 4: SESSION STORAGE                                                 │
│  ────────────────────────────                                                    │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Stateless application servers                                        │   │
│  │  • Load-balanced environments                                           │   │
│  │  • Session persistence across restarts                                  │   │
│  │                                                                          │   │
│  │  Spring Session + Redis:                                                │   │
│  │    @EnableRedisHttpSession                                              │   │
│  │    // Sessions automatically stored in Redis                            │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ✅ SCENARIO 5: RATE LIMITING                                                   │
│  ───────────────────────────                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • API rate limiting per user/IP                                        │   │
│  │  • Throttling concurrent requests                                       │   │
│  │  • Sliding window counters                                              │   │
│  │                                                                          │   │
│  │  INCR rate:user:123:minute                                              │   │
│  │  EXPIRE rate:user:123:minute 60                                         │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ✅ SCENARIO 6: REAL-TIME APPLICATIONS                                          │
│  ────────────────────────────────────                                            │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Leaderboards (Sorted Sets)                                           │   │
│  │  • Real-time counters                                                   │   │
│  │  • Pub/Sub messaging                                                    │   │
│  │  • Online presence indicators                                           │   │
│  │                                                                          │   │
│  │  ZADD leaderboard 1000 "player1" 950 "player2" 900 "player3"           │   │
│  │  ZREVRANGE leaderboard 0 9 WITHSCORES  # Top 10                        │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 13.2 When NOT to Use Redis

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    WHEN NOT TO USE REDIS                                         │
│                                                                                  │
│  ❌ SCENARIO 1: SINGLE INSTANCE APPLICATION                                     │
│  ───────────────────────────────────────────                                     │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Use Caffeine instead:                                                   │   │
│  │  • No network latency                                                   │   │
│  │  • No serialization overhead                                            │   │
│  │  • No additional infrastructure                                         │   │
│  │  • Simpler configuration                                                │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❌ SCENARIO 2: HIGHLY DYNAMIC DATA                                             │
│  ──────────────────────────────────                                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Data changes every second                                            │   │
│  │  • Cache invalidation more frequent than reads                          │   │
│  │  • Real-time stock prices (use streaming instead)                      │   │
│  │                                                                          │   │
│  │  If cache hit rate < 50%, caching may hurt more than help              │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❌ SCENARIO 3: STRONG CONSISTENCY REQUIRED                                     │
│  ──────────────────────────────────────────                                      │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Financial transactions                                               │   │
│  │  • Inventory management (exact counts)                                  │   │
│  │  • Any scenario where stale data causes business issues                │   │
│  │                                                                          │   │
│  │  Caching introduces eventual consistency by nature                      │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❌ SCENARIO 4: VERY LARGE OBJECTS                                              │
│  ─────────────────────────────────                                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Objects > 10MB                                                       │   │
│  │  • Binary files, images                                                 │   │
│  │  • Large reports                                                        │   │
│  │                                                                          │   │
│  │  Use: CDN, object storage (S3), or specialized solutions               │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  ❌ SCENARIO 5: LIMITED BUDGET / SIMPLE NEEDS                                   │
│  ─────────────────────────────────────────────                                   │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  • Small application with low traffic                                   │   │
│  │  • Redis adds operational complexity                                    │   │
│  │  • Monitoring, backups, failover needed                                │   │
│  │                                                                          │   │
│  │  Consider: In-process cache (Caffeine) or managed Redis service        │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 14. Common Issues & Troubleshooting

### 14.1 Serialization Issues

```java
// Problem: ClassNotFoundException during deserialization
// Cause: Class not found or serialVersionUID mismatch

// Solution 1: Use JSON serialization (recommended)
@Bean
public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
    GenericJackson2JsonRedisSerializer serializer = 
        new GenericJackson2JsonRedisSerializer();
    
    RedisCacheConfiguration config = RedisCacheConfiguration.defaultCacheConfig()
        .serializeValuesWith(RedisSerializationContext.SerializationPair
            .fromSerializer(serializer));
    
    return RedisCacheManager.builder(factory)
        .cacheDefaults(config)
        .build();
}

// Solution 2: Always include serialVersionUID
public class ProductDTO implements Serializable {
    private static final long serialVersionUID = 1L;  // Always define this!
    // fields...
}

// Solution 3: Handle deserialization errors gracefully
@Bean
public CacheErrorHandler errorHandler() {
    return new CacheErrorHandler() {
        @Override
        public void handleCacheGetError(RuntimeException e, Cache cache, Object key) {
            if (e.getCause() instanceof SerializationException) {
                log.warn("Serialization error, evicting key: {}", key);
                cache.evict(key);  // Remove corrupted entry
            }
        }
        // other methods...
    };
}
```

### 14.2 Connection Timeout

```java
// Problem: RedisConnectionFailureException, timeout errors

// Solution 1: Configure proper timeouts
@Bean
public LettuceConnectionFactory redisConnectionFactory() {
    RedisStandaloneConfiguration config = new RedisStandaloneConfiguration();
    config.setHostName("localhost");
    config.setPort(6379);
    
    LettuceClientConfiguration clientConfig = LettuceClientConfiguration.builder()
        .commandTimeout(Duration.ofSeconds(5))     // Command timeout
        .shutdownTimeout(Duration.ofMillis(100))   // Shutdown timeout
        .build();
    
    return new LettuceConnectionFactory(config, clientConfig);
}

// Solution 2: Enable connection pooling
spring:
  redis:
    lettuce:
      pool:
        enabled: true
        max-active: 20
        max-idle: 10
        min-idle: 5
        max-wait: 2s        # Max wait for connection from pool

// Solution 3: Health check before operations
@Component
@RequiredArgsConstructor
public class RedisHealthCheck {
    
    private final StringRedisTemplate redisTemplate;
    
    public boolean isRedisAvailable() {
        try {
            String result = redisTemplate.execute((RedisCallback<String>) 
                connection -> connection.ping());
            return "PONG".equals(result);
        } catch (Exception e) {
            return false;
        }
    }
}
```

### 14.3 Memory Overflow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS MEMORY OVERFLOW ISSUES                                  │
│                                                                                  │
│  SYMPTOMS:                                                                       │
│  • OOM (Out of Memory) errors                                                   │
│  • Redis process killed by OS                                                   │
│  • Write operations failing                                                     │
│                                                                                  │
│  SOLUTIONS:                                                                      │
│  ──────────                                                                      │
│                                                                                  │
│  1. Set maxmemory limit:                                                        │
│     CONFIG SET maxmemory 2gb                                                    │
│                                                                                  │
│  2. Configure eviction policy:                                                  │
│     CONFIG SET maxmemory-policy allkeys-lru                                     │
│                                                                                  │
│  3. Monitor memory usage:                                                       │
│     INFO memory                                                                 │
│     MEMORY DOCTOR                                                               │
│                                                                                  │
│  4. Set TTL on all cache keys:                                                  │
│     # Never store cache without TTL                                             │
│     SET key value EX 3600                                                       │
│                                                                                  │
│  5. Review key sizes:                                                           │
│     redis-cli --bigkeys                                                         │
│     MEMORY USAGE key                                                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 14.4 Cache Inconsistency

```java
// Problem: Cache shows stale data after database update

// Solution 1: Update cache after DB update (CachePut)
@CachePut(value = "products", key = "#id")
@Transactional
public ProductDTO updateProduct(Long id, ProductRequest request) {
    Product product = repository.findById(id).orElseThrow();
    // update fields...
    return toDTO(repository.save(product));
}

// Solution 2: Evict cache after DB update (safer)
@CacheEvict(value = "products", key = "#id")
@Transactional
public ProductDTO updateProduct(Long id, ProductRequest request) {
    // update logic...
}

// Solution 3: Use transaction-aware cache manager
@Bean
public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
    return RedisCacheManager.builder(factory)
        .cacheDefaults(config)
        .transactionAware()  // Sync with @Transactional
        .build();
}

// Solution 4: Listen to domain events for cache invalidation
@Component
@RequiredArgsConstructor
public class CacheInvalidationListener {
    
    private final CacheManager cacheManager;
    
    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void handleProductUpdate(ProductUpdatedEvent event) {
        Cache cache = cacheManager.getCache("products");
        if (cache != null) {
            cache.evict(event.getProductId());
        }
    }
}
```

### 14.5 TTL Misconfiguration

```java
// Problem: Cache not expiring, or expiring too quickly

// Solution: Proper TTL configuration per cache
@Bean
public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
    Map<String, RedisCacheConfiguration> configs = new HashMap<>();
    
    // Frequently changing data: short TTL
    configs.put("searchResults", createConfig(Duration.ofMinutes(5)));
    configs.put("activePromotions", createConfig(Duration.ofMinutes(15)));
    
    // Moderately changing data: medium TTL
    configs.put("products", createConfig(Duration.ofMinutes(30)));
    configs.put("userProfiles", createConfig(Duration.ofHours(1)));
    
    // Rarely changing data: long TTL
    configs.put("categories", createConfig(Duration.ofHours(24)));
    configs.put("configurations", createConfig(Duration.ofDays(7)));
    
    return RedisCacheManager.builder(factory)
        .withInitialCacheConfigurations(configs)
        .build();
}

private RedisCacheConfiguration createConfig(Duration ttl) {
    return RedisCacheConfiguration.defaultCacheConfig()
        .entryTtl(ttl);
}

// Debugging: Check TTL of cached keys
// redis-cli TTL "myapp:cache:products::1001"
```

---

## 15. Best Practices

### 15.1 Naming Conventions

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS KEY NAMING CONVENTIONS                                  │
│                                                                                  │
│  RECOMMENDED FORMAT:                                                             │
│  ───────────────────                                                             │
│  {app}:{cache}:{entity}:{identifier}                                            │
│                                                                                  │
│  EXAMPLES:                                                                       │
│  ──────────                                                                      │
│  ✓ myapp:cache:products:1001                                                   │
│  ✓ myapp:cache:users:email:john@mail.com                                       │
│  ✓ myapp:session:abc123def                                                     │
│  ✓ myapp:ratelimit:ip:192.168.1.1                                              │
│  ✓ myapp:lock:product:1001                                                     │
│                                                                                  │
│  AVOID:                                                                          │
│  ───────                                                                         │
│  ✗ product1001            (no namespacing)                                     │
│  ✗ my_app:cache:products  (underscores - use colons)                           │
│  ✗ MYAPP:CACHE:PRODUCTS   (uppercase - harder to read)                         │
│  ✗ very:long:key:name:that:goes:on:forever:and:ever:1001 (too long)           │
│                                                                                  │
│  RULES:                                                                          │
│  ───────                                                                         │
│  • Use colons (:) as separators                                                 │
│  • Use lowercase                                                                 │
│  • Keep keys short but descriptive                                              │
│  • Include version if schema might change: v1:products:1001                    │
│  • Group related keys with common prefix                                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

```java
// Key naming utility
public class CacheKeyBuilder {
    
    private static final String APP_PREFIX = "myapp";
    private static final String CACHE_PREFIX = "cache";
    private static final String SEPARATOR = ":";
    
    public static String build(String entity, Object id) {
        return String.join(SEPARATOR, APP_PREFIX, CACHE_PREFIX, entity, String.valueOf(id));
    }
    
    public static String build(String entity, String field, Object value) {
        return String.join(SEPARATOR, APP_PREFIX, CACHE_PREFIX, entity, field, String.valueOf(value));
    }
    
    // Examples:
    // CacheKeyBuilder.build("products", 1001) → "myapp:cache:products:1001"
    // CacheKeyBuilder.build("users", "email", "john@mail.com") 
    //     → "myapp:cache:users:email:john@mail.com"
}
```

### 15.2 TTL Guidelines

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    TTL GUIDELINES BY DATA TYPE                                   │
│                                                                                  │
│  │ Data Type                    │ Recommended TTL │ Reason                    │ │
│  │──────────────────────────────┼─────────────────┼───────────────────────────│ │
│  │ Static config                │ 24 hours        │ Rarely changes            │ │
│  │ Product catalog              │ 30 min - 1 hour │ Occasional updates        │ │
│  │ User profiles                │ 1 hour          │ Moderate changes          │ │
│  │ Search results               │ 5-15 minutes    │ Dynamic content           │ │
│  │ Session data                 │ 30 minutes      │ Security considerations   │ │
│  │ Rate limit counters          │ 1 minute        │ Short window needed       │ │
│  │ OAuth tokens                 │ Token expiry    │ Security critical         │ │
│  │ OTP codes                    │ 5-10 minutes    │ Short-lived by design     │ │
│  │ API responses (external)     │ Based on API    │ Respect API cache headers │ │
│                                                                                  │
│  FORMULA:                                                                        │
│  ─────────                                                                       │
│  TTL = min(acceptable_staleness, avg_update_frequency * 0.5)                    │
│                                                                                  │
│  ANTI-PATTERNS:                                                                  │
│  ───────────────                                                                 │
│  ✗ Very long TTL (days) without eviction strategy                              │
│  ✗ No TTL at all (memory leak risk)                                            │
│  ✗ Same TTL for all caches (different data, different needs)                   │
│  ✗ TTL shorter than average request processing time                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 15.3 Avoiding Large Object Caching

```java
// Problem: Caching large objects causes memory issues

// Solution 1: Cache references instead of full objects
@Service
public class OrderService {
    
    // BAD: Caching entire order with all items
    // @Cacheable("orders")
    // public Order getOrder(Long id) { ... }
    
    // GOOD: Cache order summary, fetch details on demand
    @Cacheable("orderSummaries")
    public OrderSummaryDTO getOrderSummary(Long id) {
        return OrderSummaryDTO.builder()
            .id(order.getId())
            .status(order.getStatus())
            .totalAmount(order.getTotalAmount())
            .itemCount(order.getItems().size())
            .build();
    }
    
    // Fetch full order only when needed (not cached)
    public Order getOrderDetails(Long id) {
        return orderRepository.findById(id).orElseThrow();
    }
}

// Solution 2: Paginate large collections
@Cacheable(value = "productList", key = "#page + '-' + #size")
public Page<ProductDTO> getProducts(int page, int size) {
    return repository.findAll(PageRequest.of(page, size))
        .map(this::toDTO);
}

// Solution 3: Monitor and alert on large keys
@Scheduled(fixedRate = 3600000)  // Every hour
public void checkLargeKeys() {
    // Use redis-cli --bigkeys in production
    // Or implement custom monitoring
}
```

### 15.4 Versioned Keys

```java
// Use versioned keys when cache schema changes
@Service
@RequiredArgsConstructor
public class VersionedCacheService {
    
    // Cache version - increment when ProductDTO schema changes
    private static final int CACHE_VERSION = 2;
    
    private final StringRedisTemplate redisTemplate;
    private final ObjectMapper objectMapper;
    
    public void cacheProduct(ProductDTO product) {
        String key = buildVersionedKey("product", product.getId());
        redisTemplate.opsForValue().set(key, serialize(product));
    }
    
    public ProductDTO getProduct(Long id) {
        String key = buildVersionedKey("product", id);
        String value = redisTemplate.opsForValue().get(key);
        return value != null ? deserialize(value) : null;
    }
    
    private String buildVersionedKey(String entity, Object id) {
        return String.format("v%d:%s:%s", CACHE_VERSION, entity, id);
        // Example: "v2:product:1001"
    }
    
    // On version change, old keys (v1:...) are ignored
    // They expire naturally via TTL
}

// Alternative: Clear all cache on deployment
@EventListener(ApplicationReadyEvent.class)
public void clearCacheOnStartup() {
    // Only in environments where this is acceptable
    cacheManager.getCacheNames().forEach(name -> 
        cacheManager.getCache(name).clear());
}
```

### 15.5 Monitoring & Alerts

```java
// Prometheus metrics for Redis cache
@Configuration
public class RedisMetricsConfig {
    
    @Bean
    public MeterBinder redisMetrics(RedisConnectionFactory factory) {
        return registry -> {
            // Connection pool metrics
            Gauge.builder("redis.pool.active", factory, 
                f -> getPoolMetric(f, "active"))
                .register(registry);
            
            Gauge.builder("redis.pool.idle", factory, 
                f -> getPoolMetric(f, "idle"))
                .register(registry);
        };
    }
}

// Custom cache metrics
@Aspect
@Component
@RequiredArgsConstructor
public class CacheMetricsAspect {
    
    private final MeterRegistry registry;
    
    @Around("@annotation(cacheable)")
    public Object measureCacheAccess(ProceedingJoinPoint pjp, Cacheable cacheable) 
            throws Throwable {
        String cacheName = cacheable.value()[0];
        Timer.Sample sample = Timer.start(registry);
        
        try {
            Object result = pjp.proceed();
            
            // Record hit/miss (simplified - actual would check cache)
            registry.counter("cache.access", 
                "cache", cacheName, 
                "result", result != null ? "hit" : "miss")
                .increment();
            
            return result;
        } finally {
            sample.stop(Timer.builder("cache.latency")
                .tag("cache", cacheName)
                .register(registry));
        }
    }
}
```

```yaml
# application.yml - Alerting thresholds
management:
  metrics:
    tags:
      application: ${spring.application.name}

# Alertmanager rules (Prometheus)
# alerts:
#   - alert: RedisCacheMissRateHigh
#     expr: rate(cache_misses_total[5m]) / rate(cache_requests_total[5m]) > 0.5
#     for: 5m
#     labels:
#       severity: warning
#   
#   - alert: RedisConnectionPoolExhausted
#     expr: redis_pool_active >= redis_pool_max * 0.9
#     for: 2m
#     labels:
#       severity: critical
```

### 15.6 Production Deployment Recommendations

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    PRODUCTION DEPLOYMENT CHECKLIST                               │
│                                                                                  │
│  ☑ INFRASTRUCTURE                                                               │
│  ─────────────────                                                               │
│  □ Use Redis Cluster or Sentinel for HA                                        │
│  □ Configure maxmemory appropriately                                           │
│  □ Enable persistence (RDB/AOF) if data loss unacceptable                     │
│  □ Set up Redis on dedicated servers (not shared)                             │
│  □ Use SSD storage for persistence                                             │
│  □ Configure proper network security (VPC, firewall)                          │
│                                                                                  │
│  ☑ SECURITY                                                                     │
│  ───────────                                                                     │
│  □ Enable authentication (requirepass)                                         │
│  □ Use SSL/TLS for connections                                                 │
│  □ Disable dangerous commands (FLUSHALL, DEBUG, etc.)                         │
│  □ Use ACLs for fine-grained access control (Redis 6+)                        │
│  □ Run Redis as non-root user                                                  │
│                                                                                  │
│  ☑ MONITORING                                                                   │
│  ────────────                                                                    │
│  □ Set up INFO metrics collection                                              │
│  □ Monitor memory usage and evictions                                          │
│  □ Track hit/miss rates                                                        │
│  □ Alert on connection issues                                                  │
│  □ Monitor slow log                                                            │
│  □ Set up latency monitoring                                                   │
│                                                                                  │
│  ☑ APPLICATION                                                                  │
│  ──────────────                                                                  │
│  □ Configure connection pooling                                                │
│  □ Set appropriate timeouts                                                    │
│  □ Implement circuit breaker for Redis calls                                  │
│  □ Use JSON serialization (not JDK)                                           │
│  □ Handle cache errors gracefully                                             │
│  □ Set TTL on all cache keys                                                  │
│  □ Test failover scenarios                                                     │
│                                                                                  │
│  ☑ BACKUP & RECOVERY                                                           │
│  ─────────────────────                                                           │
│  □ Schedule regular RDB backups                                                │
│  □ Test restore procedures                                                     │
│  □ Document recovery runbook                                                   │
│  □ Implement warmup strategy for cache                                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 16. Redis Caching Conclusion

### 16.1 Why Redis is Powerful for Caching

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    REDIS CACHING: KEY TAKEAWAYS                                  │
│                                                                                  │
│  🚀 PERFORMANCE                                                                  │
│  ─────────────────                                                               │
│  • Sub-millisecond latency (memory-based)                                       │
│  • 100,000+ operations per second                                               │
│  • Reduces database load by 90%+                                                │
│  • Improves user experience significantly                                       │
│                                                                                  │
│  📡 DISTRIBUTED BY DESIGN                                                       │
│  ─────────────────────────                                                       │
│  • Shared cache across all instances                                            │
│  • Native clustering and replication                                            │
│  • Perfect for microservices architecture                                       │
│  • Horizontal scalability                                                       │
│                                                                                  │
│  🛠️ FEATURE-RICH                                                                 │
│  ─────────────────                                                               │
│  • Multiple data structures                                                     │
│  • TTL and eviction policies                                                    │
│  • Pub/Sub for cache invalidation                                               │
│  • Lua scripting for atomic operations                                          │
│  • Transactions support                                                         │
│                                                                                  │
│  💪 BATTLE-TESTED                                                                │
│  ─────────────────                                                               │
│  • Used by Twitter, GitHub, Pinterest, Stack Overflow                          │
│  • Mature ecosystem with excellent tooling                                      │
│  • Active community and development                                             │
│  • Proven at massive scale                                                      │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 16.2 How Spring Boot Simplifies Redis Integration

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    SPRING BOOT + REDIS: INTEGRATION                              │
│                                                                                  │
│  STEP 1: Add Dependency                                                          │
│  ────────────────────────                                                        │
│  implementation 'org.springframework.boot:spring-boot-starter-data-redis'       │
│                                                                                  │
│  STEP 2: Configure                                                               │
│  ──────────────────                                                              │
│  spring.redis.host=localhost                                                    │
│  spring.redis.port=6379                                                         │
│                                                                                  │
│  STEP 3: Enable Caching                                                          │
│  ──────────────────────                                                          │
│  @EnableCaching                                                                  │
│  @Configuration                                                                  │
│  public class CacheConfig { }                                                   │
│                                                                                  │
│  STEP 4: Use Annotations                                                         │
│  ───────────────────────                                                         │
│  @Cacheable("products")                                                          │
│  public Product getProduct(Long id) { ... }                                     │
│                                                                                  │
│  THAT'S IT! Spring Boot handles:                                                │
│  • Connection management                                                         │
│  • Serialization/deserialization                                                │
│  • Cache abstraction                                                            │
│  • Error handling                                                               │
│  • Health checks                                                                │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 16.3 Production-Ready Implementation Guidelines

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    PRODUCTION GUIDELINES SUMMARY                                 │
│                                                                                  │
│  ✓ DESIGN                                                                       │
│  ──────────                                                                      │
│  • Choose Redis when you have multiple instances                               │
│  • Use Caffeine for single-instance applications                               │
│  • Consider L1 (local) + L2 (Redis) hybrid approach                           │
│                                                                                  │
│  ✓ CONFIGURATION                                                                │
│  ───────────────                                                                 │
│  • Always set maxmemory and eviction policy                                    │
│  • Use appropriate TTL per cache                                               │
│  • Configure connection pooling                                                │
│  • Use JSON serialization                                                      │
│                                                                                  │
│  ✓ RESILIENCE                                                                   │
│  ─────────────                                                                   │
│  • Implement graceful degradation                                              │
│  • Use circuit breaker pattern                                                 │
│  • Application should work without cache                                       │
│  • Handle cache errors without failing requests                                │
│                                                                                  │
│  ✓ MONITORING                                                                   │
│  ────────────                                                                    │
│  • Track hit/miss rates (target: >80% hit rate)                               │
│  • Monitor memory usage and evictions                                          │
│  • Alert on connection failures                                                │
│  • Review slow log regularly                                                   │
│                                                                                  │
│  ✓ SECURITY                                                                     │
│  ───────────                                                                     │
│  • Enable authentication                                                        │
│  • Use SSL/TLS in production                                                   │
│  • Never expose Redis to public internet                                       │
│                                                                                  │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│  REMEMBER:                                                                       │
│  ─────────                                                                       │
│  "Redis is a powerful tool, but it's not a silver bullet.                      │
│   Use it where it makes sense, configure it properly,                          │
│   monitor it actively, and always have a fallback plan."                       │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Redis Deep Dive Document Summary

| Section | Topics Covered |
|---------|----------------|
| 8. Introduction | What is Redis, why popular, in-memory concept, cache vs database |
| 9. How Redis Works | Key-value model, TTL, expiration, eviction policies, cache patterns |
| 10. Implementation | Dependencies, configuration, annotations, complete code examples |
| 11. Advanced Techniques | Cluster, Sentinel, locks, stampede prevention, Lua scripts |
| 12. Performance | Memory optimization, serialization, monitoring, scaling |
| 13. When to Use | Ideal scenarios, when NOT to use Redis |
| 14. Troubleshooting | Serialization, timeout, memory, consistency issues |
| 15. Best Practices | Naming, TTL, monitoring, production checklist |
| 16. Conclusion | Summary, Spring Boot integration, production guidelines |

---
- Part-1: Spring Framework Complete Guide
- Part-2: Spring Boot Complete Guide
- Part-3: Spring Security Complete Guide
- Part-4: JPA & Hibernate Complete Guide
- Part-5: Microservices Architecture Part 1
- Part-6: Microservices Architecture Part 2
- Part-7: JUnit & Mockito Testing Guide
- Architect Decision Guide: SpringBoot Microservices

---

*Document Version: 1.0*
*Last Updated: 2024*
