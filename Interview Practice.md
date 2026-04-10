# Java Collections Framework — Interview Practice (Senior Developer, 6 YOE)

> Answers written from the perspective of a senior Java developer with 6 years of production experience, including real-world scenarios, debugging war stories, and architectural decisions.

---

## 🔥 1. Core Fundamentals (Must Know)

---

### Q: What is the Java Collections Framework (JCF)?

**Answer:**

The Java Collections Framework is a unified architecture introduced in Java 1.2 that provides interfaces, implementations, and algorithms for storing and manipulating groups of objects. Before JCF, we had `Vector`, `Hashtable`, and arrays — all disconnected, inconsistent, and hard to extend.

In production, JCF is the backbone of almost every Java application. For instance, in an order management system I worked on, we used `HashMap` for fast order lookups by ID, `LinkedHashMap` for maintaining insertion-ordered audit trails, `PriorityQueue` for processing orders by priority, and `ConcurrentHashMap` for shared state across threads in the payment gateway.

The framework provides:
- **Interfaces** — `List`, `Set`, `Queue`, `Map`, `Deque`
- **Implementations** — `ArrayList`, `HashMap`, `TreeMap`, `ConcurrentHashMap`, etc.
- **Algorithms** — `Collections.sort()`, `Collections.unmodifiableList()`, `Collections.synchronizedMap()`, etc.

The real power is **programming to interfaces** — I can swap `ArrayList` for `LinkedList` or `HashMap` for `ConcurrentHashMap` without changing the caller code. This saved us during a production incident where we swapped `HashMap` to `ConcurrentHashMap` in a caching layer without touching the service layer.

---

### Q: Difference between Collection vs Collections vs Iterable?

**Answer:**

| Concept | What It Is | Role |
|---------|-----------|------|
| **`Iterable<T>`** | Root interface (`java.lang`) | Enables for-each loop. Any class implementing `Iterable` can be used in `for(T item : collection)`. Defines `iterator()` method. |
| **`Collection<E>`** | Interface (`java.util`) extending `Iterable` | Parent interface for `List`, `Set`, `Queue`. Defines core operations: `add()`, `remove()`, `contains()`, `size()`, `stream()`. **`Map` does NOT extend `Collection`.** |
| **`Collections`** | Utility class (`java.util`) | Static helper methods — `sort()`, `unmodifiableList()`, `synchronizedMap()`, `singletonList()`, `emptyList()`, etc. |

**Real-world usage:**
- I use `Collection<E>` as a method parameter type when I don't care whether the caller passes a `List`, `Set`, or `Queue` — for example, a method that sends notifications to a collection of users.
- `Collections.unmodifiableList()` is something I use extensively when returning internal state from a class — prevents callers from accidentally mutating your data. Got burned once when a downstream service modified a returned list and corrupted our cache.
- `Iterable` — I've implemented custom `Iterable` on a paginated database result wrapper, so teams could iterate over millions of rows with a simple for-each without loading everything into memory.

---

### Q: What are the main interfaces in JCF? (List, Set, Queue, Map)

**Answer:**

| Interface | Duplicates | Order | Null | Primary Use |
|-----------|-----------|-------|------|-------------|
| **`List`** | Yes | Insertion order maintained | Yes (multiple) | Ordered sequence — index-based access |
| **`Set`** | No | Depends on impl (`HashSet` = no order, `LinkedHashSet` = insertion, `TreeSet` = sorted) | `HashSet`/`LinkedHashSet`: one null, `TreeSet`: no null | Unique elements |
| **`Queue`** | Yes | FIFO (or priority-based) | Generally no (except `LinkedList`) | Processing elements in order |
| **`Map`** | Keys: No, Values: Yes | Depends on impl | `HashMap`: one null key, `TreeMap`: no null key | Key-value association |

**Production perspective:**

In a real e-commerce system I've built:
- **`List`** — Order line items (duplicates possible, order matters)
- **`Set`** — Unique coupon codes applied to a cart
- **`Queue`** — Async event processing (order placed → payment → shipping)
- **`Map`** — User session data (sessionId → UserSession), product catalog cache (productId → Product)

---

### Q: Difference between List, Set, Map?

**Answer:**

**`List`** — Think of it as an ordered sequence with an index. You can have duplicates, you can access by position. Use when order matters and duplicates are valid.

```java
List<String> recentSearches = new ArrayList<>();
recentSearches.add("laptop");
recentSearches.add("laptop"); // duplicate allowed — user searched twice
recentSearches.get(0);        // index-based access
```

**`Set`** — A mathematical set. No duplicates. Use when uniqueness is the constraint.

```java
Set<String> activeUserIds = new HashSet<>();
activeUserIds.add("user-123");
activeUserIds.add("user-123"); // ignored — already exists
// Used in production: tracking unique visitors, deduplicating event streams
```

**`Map`** — Key-value pairs. Keys are unique, values can repeat. This is probably the most-used collection in any production system.

```java
Map<String, UserSession> sessions = new ConcurrentHashMap<>();
sessions.put("sess-abc", new UserSession("user-123", Instant.now()));
sessions.get("sess-abc"); // O(1) lookup
// Used in production: caching, configuration, routing tables, session management
```

**Key difference in production thinking:**
- If I'm asked "should we use a List or Set?", I ask: "Do you need duplicates? Do you need index-based access?" If no to both → `Set`.
- If someone uses a `List` and then manually checks `contains()` before adding, that's a red flag — they should use a `Set`. I've refactored this pattern multiple times in code reviews.

---

### Q: Difference between ArrayList vs LinkedList?

**Answer:**

| Feature | ArrayList | LinkedList |
|---------|-----------|------------|
| **Internal structure** | Dynamic array (`Object[]`) | Doubly-linked list (nodes with prev/next pointers) |
| **Random access** | O(1) — direct index calculation | O(n) — must traverse from head/tail |
| **Add at end** | O(1) amortized (occasional resize) | O(1) — just link a new node |
| **Add at middle** | O(n) — shifts elements right | O(1) if you have the node reference, but O(n) to find the position |
| **Memory** | Compact — contiguous array | ~40 bytes overhead per element (prev + next pointers + node object header) |
| **Cache performance** | Excellent — CPU cache-friendly due to contiguity | Poor — nodes scattered in heap, cache misses |
| **Iterator remove** | O(n) due to shifting | O(1) |

**My production stance:**

In 6 years, I've used `LinkedList` in production maybe twice. `ArrayList` wins in almost every real-world scenario because:

1. **CPU cache locality** — Modern CPUs fetch memory in cache lines (64 bytes). `ArrayList`'s contiguous array means sequential access hits the L1/L2 cache. `LinkedList` nodes are scattered across the heap — every `node.next` is a potential cache miss. I've seen benchmark differences of 5-10x for iteration.

2. **Memory overhead** — For a list of 1 million `Integer` objects, `LinkedList` uses roughly 3x more memory than `ArrayList` because each node stores two pointers plus object header overhead.

3. **The "fast insertion" myth** — People say `LinkedList` has O(1) insertion, but that's only if you already have a reference to the node. In practice, you almost always need to find the position first (O(n)), making it no better than `ArrayList`.

**When I've actually used LinkedList:**
- As a `Deque` (double-ended queue) for a sliding window in a rate limiter — frequent add-first/remove-last operations. Even then, `ArrayDeque` is usually better.

---

### Q: Difference between HashSet vs TreeSet vs LinkedHashSet?

**Answer:**

| Feature | HashSet | TreeSet | LinkedHashSet |
|---------|---------|---------|---------------|
| **Backing structure** | `HashMap` | `TreeMap` (Red-Black Tree) | `LinkedHashMap` |
| **Order** | No guaranteed order | Sorted (natural or Comparator) | Insertion order |
| **Null** | Allows one null | No null (throws NPE — can't compare null) | Allows one null |
| **add/remove/contains** | O(1) average | O(log n) | O(1) average |
| **Use case** | Fast unique checks | When you need sorted unique elements | When you need unique + insertion order |

**Production usage:**

```java
// HashSet — deduplicating event IDs in a message consumer
Set<String> processedEventIds = new HashSet<>(10_000);

// TreeSet — leaderboard scores (auto-sorted)
TreeSet<PlayerScore> leaderboard = new TreeSet<>(Comparator.comparingInt(PlayerScore::score).reversed());
leaderboard.first(); // top scorer

// LinkedHashSet — maintaining unique tags in user-entered order
Set<String> tags = new LinkedHashSet<>();
tags.add("java");
tags.add("spring");
tags.add("java"); // ignored, order preserved: [java, spring]
```

**Real incident:** We had a bug where API response JSON had set elements in random order across different calls (using `HashSet`). A downstream consumer was doing string comparison on the entire JSON body for caching. Switching to `LinkedHashSet` fixed the inconsistency. Lesson: if your set is ever serialized, use `LinkedHashSet` for deterministic output.

---

### Q: Difference between HashMap vs TreeMap vs LinkedHashMap?

**Answer:**

| Feature | HashMap | TreeMap | LinkedHashMap |
|---------|---------|---------|---------------|
| **Internal structure** | Array of buckets + linked list/tree | Red-Black Tree | HashMap + doubly-linked list |
| **Order** | No order guaranteed | Sorted by key (natural or Comparator) | Insertion order (or access order) |
| **Null key** | One null key allowed | No null key (NPE) | One null key allowed |
| **get/put** | O(1) average | O(log n) | O(1) average |
| **Use case** | General-purpose fast lookup | Range queries, sorted iteration | Insertion-ordered or LRU cache |

**Production scenarios:**

```java
// HashMap — 90% of use cases: caching, lookups, configuration
Map<String, Config> configCache = new HashMap<>();

// TreeMap — finding all keys in a range (time-series data)
TreeMap<Instant, MetricValue> metrics = new TreeMap<>();
// All metrics between 2pm and 3pm:
metrics.subMap(start, end);

// LinkedHashMap — LRU cache (access-order mode)
Map<String, Object> lruCache = new LinkedHashMap<>(16, 0.75f, true) {
    @Override
    protected boolean removeEldestEntry(Map.Entry<String, Object> eldest) {
        return size() > MAX_CACHE_SIZE;
    }
};
```

I built a simple LRU config cache using `LinkedHashMap` with `accessOrder=true` in a microservice that couldn't justify adding Redis. It handled 10K RPM for a year with zero issues. Only when we scaled beyond that did we move to Redis.

---

### Q: What is Comparable vs Comparator?

**Answer:**

| Aspect | Comparable | Comparator |
|--------|-----------|------------|
| **Package** | `java.lang` | `java.util` |
| **Method** | `compareTo(T o)` | `compare(T o1, T o2)` |
| **Defines** | Natural ordering (single) | External/custom ordering (multiple possible) |
| **Modifies class** | Yes — class must implement it | No — separate strategy object |
| **Use** | `Collections.sort(list)` | `Collections.sort(list, comparator)` |

**Production approach:**

I implement `Comparable` for the one natural ordering that makes business sense — e.g., `Order` sorted by `createdAt`. Then I create `Comparator` instances for all other sort needs.

```java
public class Order implements Comparable<Order> {
    private Instant createdAt;
    private BigDecimal amount;
    private Priority priority;

    @Override
    public int compareTo(Order other) {
        return this.createdAt.compareTo(other.createdAt); // natural order: by time
    }
}

// Multiple comparators for different views
Comparator<Order> byAmount = Comparator.comparing(Order::getAmount).reversed();
Comparator<Order> byPriority = Comparator.comparing(Order::getPriority)
                                         .thenComparing(Order::getCreatedAt);
```

**Java 8+ Comparator is incredibly powerful:**
```java
Comparator.comparing(Employee::getDepartment)
          .thenComparing(Employee::getSalary, Comparator.reverseOrder())
          .thenComparing(Employee::getName);
```

**Gotcha I've hit:** Never return `o1.value - o2.value` in a comparator for `int` — it overflows for extreme values (e.g., `Integer.MIN_VALUE - 1`). Always use `Integer.compare(a, b)`.

---

### Q: What is equals() and hashCode() contract?

**Answer:**

The contract is:

1. **If `a.equals(b)` is true, then `a.hashCode() == b.hashCode()` must be true.**
2. If `a.hashCode() == b.hashCode()`, `a.equals(b)` may or may not be true (hash collision).
3. `equals()` must be reflexive, symmetric, transitive, consistent, and `a.equals(null)` must return false.

**Why this matters in production:**

`HashMap`, `HashSet`, `LinkedHashMap` — all hash-based collections depend on this contract. When you call `map.get(key)`:
1. Compute `key.hashCode()` → find the bucket index
2. In that bucket, compare each entry using `key.equals(entry.key)` to find the exact match

If `hashCode()` is inconsistent with `equals()`, objects that are "equal" will land in different buckets, and you'll never find them.

```java
public class OrderId {
    private final String region;
    private final long sequence;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        OrderId that = (OrderId) o;
        return sequence == that.sequence && Objects.equals(region, that.region);
    }

    @Override
    public int hashCode() {
        return Objects.hash(region, sequence);
    }
}
```

**Best practice:** Use `@EqualsAndHashCode` from Lombok or IDE generation and include only the fields that define business identity — not mutable state, not audit fields.

---

### Q: What happens if you override equals() but not hashCode()?

**Answer:**

You break the contract, and hash-based collections will malfunction silently.

```java
public class UserId {
    private String id;

    @Override
    public boolean equals(Object o) {
        if (o instanceof UserId) return this.id.equals(((UserId) o).id);
        return false;
    }
    // hashCode() NOT overridden — uses default Object.hashCode() (memory address based)
}

Map<UserId, String> map = new HashMap<>();
UserId key1 = new UserId("user-123");
map.put(key1, "Alice");

UserId key2 = new UserId("user-123");
map.get(key2); // returns null! Even though key1.equals(key2) is true
```

**What happens internally:**
- `key1.hashCode()` returns (say) 12345 → stored in bucket 5
- `key2.hashCode()` returns (say) 67890 → looks in bucket 12
- HashMap never even checks `equals()` because it's looking in the wrong bucket

**Real production incident:** A colleague created a custom `CacheKey` class, overrode `equals()` but forgot `hashCode()`. The cache hit rate dropped to near 0% because every lookup created a new object with a different default `hashCode()`. We were essentially writing to the cache but never reading from it. Took two hours to find because functionally the app "worked" — it just bypassed the cache entirely and hit the database every time, causing a slow performance degradation under load.

---

## ⚙️ 2. Internal Working (Very Important)

---

### HashMap Deep Dive

---

### Q: How does HashMap work internally?

**Answer:**

`HashMap` internally uses an **array of `Node<K,V>`** (called the "table" or "bucket array"). Each `Node` contains:
- `int hash` — cached hash value
- `K key`
- `V value`
- `Node<K,V> next` — pointer to next node (for collision chaining)

**Put operation (`map.put(key, value)`):**

```
1. Compute hash:  hash = spread(key.hashCode())
                  // spread() mixes higher bits into lower bits to reduce collisions
                  // h ^ (h >>> 16)

2. Find bucket:   index = hash & (table.length - 1)
                  // Bitwise AND instead of modulo — works because capacity is power of 2

3. Check bucket:
   - If empty → create new Node, place it there
   - If occupied → traverse the linked list / tree:
     a. If key matches (hash == hash && (key == k || key.equals(k))) → replace value
     b. If no match → append new node at the end
     c. If chain length >= TREEIFY_THRESHOLD (8) AND table.length >= 64 → convert to Red-Black Tree

4. After insertion:
   - Increment size
   - If size > threshold (capacity × loadFactor) → resize (double the array, rehash)
```

**Get operation (`map.get(key)`):**

```
1. Compute hash → find bucket index
2. If bucket is empty → return null
3. Check first node — if hash matches AND equals() → return value
4. If it's a TreeNode → tree search O(log n)
5. Else → traverse linked list O(n in bucket)
```

**Visual:**
```
table[] (length = 16)
┌──────┬──────┬──────┬──────┬──────┬───────┬──────┐
│ [0]  │ [1]  │ [2]  │ [3]  │ [4]  │ ...   │ [15] │
│ null │ Node │ null │ Node │ null │       │ null │
└──────┴──┬───┴──────┴──┬───┴──────┴───────┴──────┘
          │             │
          ▼             ▼
        (A,1)         (C,3)
          │
          ▼
        (B,2)  ← collision chain
```

---

### Q: What is hashing?

**Answer:**

Hashing is the process of converting an object into a fixed-size integer (hash code) that can be used as an index into an array. The goal is to distribute objects uniformly across buckets.

In Java:
- `Object.hashCode()` returns an `int` (32-bit) — roughly 4 billion possible values
- `HashMap` doesn't use this directly — it applies a **spread function**: `h ^ (h >>> 16)` which XORs the upper 16 bits with the lower 16 bits. This ensures that even if your hashCode only varies in the high bits, the variation still affects the bucket index.

**Why spread matters:**

If your table size is 16, the bucket index is `hash & 15` (only the last 4 bits matter). Without spreading, two keys whose hashCodes differ only in the upper bits would collide. The spread function is cheap (one XOR, one shift) and significantly reduces collisions in practice.

```java
static final int hash(Object key) {
    int h;
    return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
}
```

**Production insight:** I always make sure custom key classes have a good hashCode distribution. `Objects.hash(field1, field2)` uses 31-based polynomial hashing which is decent. For high-performance scenarios with millions of entries, I've profiled bucket distribution to ensure we're not getting hotspots.

---

### Q: What is a hash collision and how is it handled?

**Answer:**

A hash collision occurs when two different keys produce the same bucket index: `hash(key1) & (n-1) == hash(key2) & (n-1)`.

**Collision resolution in Java's HashMap:**

**Before Java 8:** Chaining with linked list only.
- Colliding entries form a singly linked list in the bucket.
- Worst case: all keys land in one bucket → O(n) lookup.

**Java 8+:** Chaining with linked list + treeification.
- Starts as linked list.
- When a single bucket has **8+ nodes** (TREEIFY_THRESHOLD) AND table capacity is **≥ 64** → converts to a balanced Red-Black Tree.
- When tree shrinks to **6 or fewer** (UNTREEIFY_THRESHOLD) → converts back to linked list.
- Worst case with tree: O(log n) per bucket instead of O(n).

**Why 8 and not 5 or 10?**
Under random hash distribution, the probability of 8 or more entries in a single bucket follows a Poisson distribution and is approximately 0.00000006. So treeification is a safety net for pathological cases (bad hashCode implementations or hash-flooding attacks), not a normal operation.

**Real scenario:** We once had a HashMap used to index objects by a composite key where the `hashCode()` was poorly implemented (returning the same hash for many keys). The map degraded from O(1) to O(n) lookups under load. After Java 8, the tree fallback saved us from complete disaster, but we still fixed the root cause — the hashCode implementation.

---

### Q: What is a bucket / bin in HashMap?

**Answer:**

A bucket (or bin) is a single slot in HashMap's internal array (`Node<K,V>[] table`). Each bucket can hold:
- `null` — empty, no entry at this index
- A single `Node` — one key-value pair
- A linked list of `Node`s — multiple entries that hash to the same index (collision chain)
- A `TreeNode` (Red-Black Tree) — when collisions exceed the treeify threshold (Java 8+)

The bucket index is calculated as: `index = hash & (table.length - 1)`

Since table length is always a power of 2, `(table.length - 1)` gives a bitmask. For table size 16: `hash & 0x0F` → only the last 4 bits determine the bucket.

---

### Q: What is load factor and threshold?

**Answer:**

**Load factor** — A float value (default `0.75f`) that controls the trade-off between space and time.

- `load factor = number of entries / number of buckets`
- Higher load factor → more entries per bucket → more collisions → slower lookups, but less memory
- Lower load factor → fewer collisions → faster lookups, but more memory wasted on empty buckets

**Threshold** — The actual trigger point for resizing: `threshold = capacity × loadFactor`

- Default: `16 × 0.75 = 12`. When the 13th entry is added, the table resizes to 32 buckets.

**When to change the default:**

```java
// If you know you'll store ~1000 entries, pre-size to avoid repeated resizing:
Map<String, Object> map = new HashMap<>(1400); // 1000 / 0.75 ≈ 1334, rounds up to 2048

// If memory is tight and you can tolerate slightly slower lookups:
Map<String, Object> map = new HashMap<>(16, 0.9f);

// If you need maximum speed and have memory to spare:
Map<String, Object> map = new HashMap<>(16, 0.5f);
```

**Production tip:** In hot paths, I always pre-size HashMaps when I know the approximate entry count. The default capacity of 16 means inserting 10,000 entries causes **~10 resizes** (16→32→64→128→256→512→1024→2048→4096→8192→16384), each copying the entire table. Pre-sizing with `new HashMap<>(expectedSize * 4 / 3 + 1)` avoids this entirely.

---

### Q: What is rehashing / resizing?

**Answer:**

When `size > threshold`, HashMap doubles its array and rehashes all existing entries into the new array.

**Process:**
1. Create a new `Node[]` array of double the size (e.g., 16 → 32)
2. For each entry in the old array, recalculate its new bucket index using the new capacity
3. Move each entry to its new position in the new array

**Java 8 optimization:** Instead of recalculating `hash % newCapacity` for every entry, Java 8 uses a clever bit trick. Since capacity doubles, each entry either stays at the same index or moves to `oldIndex + oldCapacity`. It checks a single bit: `(hash & oldCapacity) == 0` → stays, else → moves. This avoids re-computing modulo and also preserves the linked list order.

**Performance impact:**
- Resizing is O(n) — every entry must be visited
- During resize, all other operations on the map are blocked (single-threaded context)
- In production, unexpected resizing under load causes latency spikes

**Real-world lesson:** We had a p99 latency spike in a payment service traced to HashMap resizing. The map was being populated in a loop from a database result set (~50K rows), starting with default capacity 16. That's 12 resizes. Pre-sizing the map cut our p99 by 40ms in that code path.

---

### Q: What is treeification in HashMap (Java 8+)?

**Answer:**

Treeification is the process of converting a linked list in a HashMap bucket into a balanced Red-Black Tree when the chain length exceeds the `TREEIFY_THRESHOLD` (8) AND the table size is at least `MIN_TREEIFY_CAPACITY` (64).

**Why both conditions?**
- If the table is small (< 64 buckets), collisions are more likely just because there aren't enough buckets. In this case, HashMap prefers to **resize** the table (more buckets → better distribution) rather than treeify.
- Treeification is a last resort for when collisions persist even with a large table.

**The process:**
```
LinkedList (chain length ≥ 8, table ≥ 64):
  Node → Node → Node → Node → Node → Node → Node → Node
                    ↓ treeify
  Converts to TreeNode (Red-Black Tree):
          Node4
         /     \
      Node2    Node6
      /  \     /  \
   Node1 Node3 Node5 Node7
                       \
                      Node8
```

**Untreeification** happens when tree size drops to ≤ 6 (UNTREEIFY_THRESHOLD), converting back to a linked list to avoid the overhead of maintaining tree structure for small chains.

**Why the gap (8 to treeify, 6 to untreeify)?** To prevent thrashing — without hysteresis, a bucket hovering around 7-8 elements would constantly convert back and forth.

---

### Q: Why is Red-Black Tree used in HashMap?

**Answer:**

Red-Black Tree is chosen because it provides **guaranteed O(log n) operations** (search, insert, delete) with **cheaper rebalancing** than AVL trees.

**Why not AVL Tree?**
- AVL trees are more strictly balanced (height difference ≤ 1), giving slightly faster lookups.
- But AVL requires more rotations during insert/delete.
- HashMap's treeified buckets need frequent insertions and deletions (as entries are added/removed/rehashed). Red-Black Trees allow up to 2x height difference between subtrees, resulting in fewer rotations — better insert/delete performance.

**Why not just a sorted array?**
- Sorted array gives O(log n) search via binary search, but O(n) insertion.
- Red-Black Tree gives O(log n) for all three operations.

**In practice**, treeification is rare. With a decent hashCode and load factor of 0.75, you'll almost never see 8+ collisions in a bucket. It's primarily a defense against **hash-flooding DoS attacks** where an attacker crafts inputs that all hash to the same bucket, degrading HashMap to O(n). With treeification, the worst case becomes O(log n).

---

### Q: What is the time complexity of HashMap operations?

**Answer:**

| Operation | Average Case | Worst Case (pre-Java 8) | Worst Case (Java 8+) |
|-----------|-------------|------------------------|----------------------|
| `put(key, value)` | O(1) | O(n) | O(log n) |
| `get(key)` | O(1) | O(n) | O(log n) |
| `remove(key)` | O(1) | O(n) | O(log n) |
| `containsKey(key)` | O(1) | O(n) | O(log n) |
| `containsValue(value)` | O(n) | O(n) | O(n) |
| resize | O(n) | O(n) | O(n) |

**The "O(1)" asterisk:** It's O(1) *amortized* and *assuming good hash distribution*. The constant factor includes: computing hashCode, spread function, array access, and potentially 1-2 equality checks. In practice, for well-designed keys, this is extremely fast — sub-microsecond.

**`containsValue()` is always O(n)** because values aren't indexed. I've seen junior developers use `containsValue()` in a loop — that's O(n²). If you need value-based lookups, maintain a reverse map or use a BiMap from Guava.

---

### ArrayList Internals

---

### Q: How does ArrayList grow dynamically?

**Answer:**

ArrayList internally uses `Object[] elementData`. When you add an element and the array is full:

1. Calculate new capacity: `newCapacity = oldCapacity + (oldCapacity >> 1)` — **1.5x growth**
2. Create a new array of `newCapacity` using `Arrays.copyOf()` (which calls `System.arraycopy()`)
3. Copy all elements from old array to new array
4. Old array becomes eligible for GC
5. Add the new element

```
Initial:   [_, _, _, _, _, _, _, _, _, _]  capacity=10, size=0
After 10:  [a, b, c, d, e, f, g, h, i, j]  capacity=10, size=10
Add 11th → resize:
           [a, b, c, d, e, f, g, h, i, j, k, _, _, _, _]  capacity=15, size=11
```

**Why 1.5x and not 2x?**
- 2x growth (used by `Vector`) wastes more memory. With 1.5x, after repeated doublings, the unused trailing space is smaller.
- `Vector` uses `2x` which is more wasteful; this is one reason `ArrayList` replaced it.

**Amortized O(1):** While a single resize is O(n), if you add n elements, the total copy operations across all resizes sum to approximately 3n. So per-element cost is O(1) amortized.

---

### Q: What is the default capacity?

**Answer:**

- **Default initial capacity:** 10 (when you use `new ArrayList<>()`)
- But actually, `new ArrayList<>()` creates an **empty array** (`{}`) initially.
- The first `add()` call triggers the first allocation to capacity 10.
- This is a Java 8+ optimization — empty ArrayLists (common in frameworks with lazy init) consume zero heap for the backing array.

```java
// No-arg constructor — starts truly empty, grows to 10 on first add
List<String> list1 = new ArrayList<>();

// Pre-sized — allocates array of 500 immediately
List<String> list2 = new ArrayList<>(500);

// From collection — sized to match
List<String> list3 = new ArrayList<>(existingCollection);
```

**Production best practice:** Always pre-size when you know the approximate element count.
```java
// Bad — causes multiple resizes
List<OrderDTO> orders = new ArrayList<>();
for (Order order : dbResults) { // 5000 results
    orders.add(toDTO(order));
}

// Good — one allocation, no resizes
List<OrderDTO> orders = new ArrayList<>(dbResults.size());
```

---

### Q: What is the cost of resizing?

**Answer:**

Each resize involves:
1. **Memory allocation** — Allocating a new, larger array on the heap
2. **Array copy** — `System.arraycopy()` (native method, very optimized — uses CPU memory-copy instructions)
3. **GC pressure** — The old array becomes garbage, contributing to GC pauses

**Concrete numbers:** Starting from default capacity 10, adding 1 million elements triggers approximately 44 resizes (10 → 15 → 22 → 33 → ... → ~1.5M). The total elements copied across all resizes is roughly 3 million (about 3× the final size).

**In high-throughput systems**, this matters:
- Each resize briefly doubles memory usage (old + new array exist simultaneously)
- Large array allocations (> ~8KB) may go directly to old generation in some GC configurations, affecting GC patterns
- `System.arraycopy` is fast but still causes a pause proportional to array size

**Mitigation:** `new ArrayList<>(expectedSize)`. In one project, we reduced young-gen GC pauses by 15% just by pre-sizing collections in a hot data-processing pipeline.

---

### LinkedList Internals

---

### Q: How does LinkedList work internally?

**Answer:**

Java's `LinkedList` is a **doubly-linked list**. Each element is wrapped in a `Node` object:

```java
private static class Node<E> {
    E item;
    Node<E> next;
    Node<E> prev;
}
```

The `LinkedList` maintains:
- `Node<E> first` — pointer to the first node
- `Node<E> last` — pointer to the last node
- `int size` — count of elements

**Operations:**
- `addFirst(e)` / `addLast(e)` — O(1), just update first/last pointers
- `get(index)` — O(n), traverses from whichever end is closer (optimizes by checking `index < size/2`)
- `add(index, element)` — O(n) to find the position, O(1) to insert
- `remove(Object)` — O(n) to find, O(1) to unlink

**Memory layout:**
```
[first] → Node(prev=null, item=A, next=→) ⇄ Node(prev=←, item=B, next=→) ⇄ Node(prev=←, item=C, next=null) ← [last]
```

`LinkedList` also implements `Deque`, so it supports stack and queue operations. But `ArrayDeque` outperforms it for both.

---

### Q: Why is insertion fast but access slow?

**Answer:**

**Fast insertion (at a known position):** Inserting a node into a doubly-linked list only requires updating 4 pointers — O(1):
```
Before: ... ⇄ [A] ⇄ [C] ⇄ ...
Insert B between A and C:
  B.prev = A
  B.next = C
  A.next = B
  C.prev = B
After:  ... ⇄ [A] ⇄ [B] ⇄ [C] ⇄ ...
```

No elements need to be shifted (unlike ArrayList where inserting at position `i` shifts `n-i` elements).

**Slow access:** There's no index-based random access. To reach element at index `k`, you must traverse node by node from either `first` (if `k < size/2`) or `last` (if `k >= size/2`). Getting element at index 500 in a 1000-element list means following 500 `next` pointers.

**The hidden cost:** Even insertion by index — `list.add(index, element)` — is O(n) because you first need O(n) traversal to find the position. The insertion itself is O(1), but finding where to insert dominates.

**CPU cache penalty:** Each node is a separate heap object, allocated at different memory addresses. Traversing the list means following pointers to random memory locations — every step is potentially a CPU cache miss. This is why even O(n) traversal in LinkedList is significantly slower than O(n) traversal in ArrayList in practice.

---

## 🚫 3. Fail-Fast & Iteration

---

### Q: What is fail-fast vs fail-safe iterator?

**Answer:**

| Aspect | Fail-Fast | Fail-Safe |
|--------|----------|-----------|
| **Behavior** | Throws `ConcurrentModificationException` if collection is modified during iteration | Does not throw — works on a copy or a snapshot |
| **Mechanism** | Checks `modCount` against expected count | Operates on a separate data structure |
| **Collections** | `ArrayList`, `HashMap`, `HashSet`, `LinkedList` — all standard collections | `CopyOnWriteArrayList`, `ConcurrentHashMap`, `ConcurrentSkipListMap` |
| **Overhead** | None — just a counter check | Memory (copy) or complexity (weakly consistent view) |
| **Guarantee** | Best-effort detection, not guaranteed under all concurrency scenarios | Never throws CME |

**Important nuance:** `ConcurrentHashMap`'s iterator is technically **weakly consistent**, not fail-safe in the traditional sense. It reflects some (but not necessarily all) modifications made after the iterator was created. It doesn't work on a full copy like `CopyOnWriteArrayList`.

**Production relevance:**

```java
// This WILL throw ConcurrentModificationException
List<Order> orders = new ArrayList<>(getOrders());
for (Order order : orders) {
    if (order.isExpired()) {
        orders.remove(order); // Structural modification during iteration
    }
}

// Safe alternatives:
// 1. Iterator.remove()
Iterator<Order> it = orders.iterator();
while (it.hasNext()) {
    if (it.next().isExpired()) it.remove();
}

// 2. removeIf (Java 8+) — cleanest
orders.removeIf(Order::isExpired);

// 3. CopyOnWriteArrayList — for concurrent scenarios
List<Order> orders = new CopyOnWriteArrayList<>(getOrders());
```

---

### Q: What is ConcurrentModificationException?

**Answer:**

`ConcurrentModificationException` (CME) is thrown when a collection detects that it was structurally modified (add, remove, clear) while being iterated — but NOT through the iterator's own methods.

**Key points:**
- It's NOT limited to multi-threaded scenarios. A single thread can trigger it by modifying a collection inside a for-each loop.
- It's a **best-effort** mechanism — not guaranteed to detect all concurrent modifications (the Javadoc explicitly states this).
- It extends `RuntimeException` (unchecked).

**Common triggers:**

```java
// 1. Single-threaded — modifying inside for-each
for (String item : list) {
    if (item.equals("x")) list.remove(item); // CME
}

// 2. Multi-threaded — one thread iterates, another modifies
// Thread 1:
for (String item : sharedList) { process(item); }
// Thread 2:
sharedList.add("new"); // CME in Thread 1

// 3. Stream + modification
list.stream().forEach(item -> {
    if (item.equals("x")) list.remove(item); // CME
});
```

**The fix depends on the scenario:**
- Single-threaded: Use `Iterator.remove()` or `removeIf()`
- Multi-threaded: Use `ConcurrentHashMap`, `CopyOnWriteArrayList`, or external synchronization

---

### Q: How does modCount work?

**Answer:**

`modCount` is a `transient int` field in `AbstractList` (and similar base classes) that tracks the number of **structural modifications** — operations that change the size of the collection (add, remove, clear), not value updates.

**Mechanism:**
1. Every structural modification increments `modCount`
2. When an iterator is created, it captures the current `modCount` as `expectedModCount`
3. On each `next()` or `remove()` call, the iterator checks: `modCount != expectedModCount`
4. If they differ → throw `ConcurrentModificationException`

```java
// Simplified from ArrayList source:
private class Itr implements Iterator<E> {
    int expectedModCount = modCount; // captured at iterator creation

    public E next() {
        checkForComodification();
        // ... return element
    }

    final void checkForComodification() {
        if (modCount != expectedModCount)
            throw new ConcurrentModificationException();
    }

    public void remove() {
        // ... remove element
        expectedModCount = modCount; // sync after legal modification
    }
}
```

**Why `Iterator.remove()` doesn't throw:** Because after removing via the iterator, it synchronizes `expectedModCount = modCount`. The iterator "knows" about the modification.

---

### Q: Difference between Iterator, ListIterator, and Enumeration?

**Answer:**

| Feature | Iterator | ListIterator | Enumeration |
|---------|----------|-------------|-------------|
| **Introduced** | Java 1.2 | Java 1.2 | Java 1.0 (legacy) |
| **Direction** | Forward only | Forward and backward | Forward only |
| **Available on** | All `Collection` types | Only `List` | Legacy: `Vector`, `Hashtable` |
| **Can remove** | Yes — `remove()` | Yes — `remove()`, `add()`, `set()` | No |
| **Fail-fast** | Yes | Yes | No |
| **Methods** | `hasNext()`, `next()`, `remove()` | All of Iterator + `hasPrevious()`, `previous()`, `nextIndex()`, `previousIndex()`, `add()`, `set()` | `hasMoreElements()`, `nextElement()` |

**When I use each:**

```java
// Iterator — 90% of cases, simple forward traversal with optional removal
Iterator<String> it = set.iterator();
while (it.hasNext()) {
    if (shouldRemove(it.next())) it.remove();
}

// ListIterator — when I need bidirectional traversal or in-place modification
ListIterator<String> lit = list.listIterator(list.size()); // start at end
while (lit.hasPrevious()) {
    String item = lit.previous();
    lit.set(item.toUpperCase()); // modify in place
}

// Enumeration — only when dealing with legacy APIs
// e.g., reading servlet request headers:
Enumeration<String> headers = request.getHeaderNames();
while (headers.hasMoreElements()) {
    String name = headers.nextElement();
}
```

**In modern Java (8+),** I rarely use explicit iterators. I use `forEach`, streams, or `removeIf()`. Explicit iterators are reserved for cases where I need `remove()` during iteration or bidirectional traversal.

---

## 🔁 4. Concurrency & Thread Safety (Very Important for Senior Roles)

---

### Q: Difference between HashMap vs Hashtable?

**Answer:**

| Feature | HashMap | Hashtable |
|---------|---------|-----------|
| **Thread safety** | Not synchronized | All methods are `synchronized` |
| **Null support** | One null key, multiple null values | No null key or value (NPE) |
| **Performance** | Fast (no lock overhead) | Slow (method-level synchronization) |
| **Iterator** | Fail-fast | Fail-safe (Enumerator), fail-fast (Iterator) |
| **Introduced** | Java 1.2 | Java 1.0 (legacy) |
| **Superclass** | `AbstractMap` | `Dictionary` (obsolete) |
| **Recommended** | Yes | No — use `ConcurrentHashMap` instead |

**Why Hashtable is dead in modern Java:**

`Hashtable` synchronizes every method (`get`, `put`, `size`, even `toString`), which means:
1. **Coarse-grained locking** — Only one thread can access the map at a time, even for reads
2. **No compound atomicity** — `if (!map.containsKey(k)) map.put(k, v)` is NOT atomic even with Hashtable — another thread can sneak in between the two calls
3. **No null support** — an unnecessary restriction

```java
// Don't do this
Map<String, String> cache = new Hashtable<>();

// Do this for thread safety
Map<String, String> cache = new ConcurrentHashMap<>();

// Or if you need a synchronized wrapper for some reason
Map<String, String> cache = Collections.synchronizedMap(new HashMap<>());
```

I've refactored `Hashtable` to `ConcurrentHashMap` in legacy codebases and seen throughput improve 3-5x under concurrent access because `ConcurrentHashMap` allows concurrent reads and segment-level (or bucket-level in Java 8+) writes.

---

### Q: Why is HashMap not thread-safe?

**Answer:**

HashMap has no synchronization — no locks, no volatile fields, no CAS operations. Multiple threads accessing a HashMap concurrently can cause:

**1. Lost updates (race condition on put):**
```
Thread 1: put("key", value1) — computes bucket 5, about to insert
Thread 2: put("key2", value2) — computes bucket 5, inserts
Thread 1: inserts — overwrites Thread 2's entry (same bucket, same slot)
```

**2. Infinite loop (pre-Java 8 — during resize):**
In Java 7, HashMap used head-insertion for linked list buckets during resize. Two threads resizing simultaneously could create a circular linked list. A subsequent `get()` would loop forever, pinning the CPU at 100%.

**3. Stale reads:** Without memory barriers, Thread 1 may not see Thread 2's modifications due to CPU caching / JMM visibility rules.

**4. Corrupted size:** `size++` is not atomic (read-increment-write). Concurrent puts can produce an incorrect `size`, causing unexpected resizing behavior.

**Real production story:** In a pre-Java 8 service, we had a thread running at 100% CPU. Thread dump showed it stuck inside `HashMap.get()` — infinite loop caused by concurrent modification during resize. The fix was switching to `ConcurrentHashMap`. This is one of the classic Java concurrency bugs.

---

### Q: What issues occur in concurrent HashMap (pre-Java 8)? Infinite loop problem?

**Answer:**

The infamous **infinite loop** in Java 7's `HashMap.resize()`:

**Root cause:** Java 7 used **head insertion** when transferring entries during resize. In a linked list `A → B → C`, after resize with head insertion, the order reverses: `C → B → A`.

**The deadly sequence:**
```
Initial state (Bucket 3): A → B → null

Thread 1 starts resize:
  - Reads: e = A, next = B
  - Gets suspended by OS

Thread 2 completes resize:
  - Bucket X: B → A → null (reversed due to head insertion)

Thread 1 resumes:
  - Processes A: inserts A at head → A → null
  - Processes B (which now points to A due to Thread 2's resize):
    B → A → null
  - Processes A again (because B.next = A):
    A → B → A → B → ... (circular reference!)

Any subsequent get() hitting this bucket enters an infinite loop.
```

**Java 8 fix:** Changed to **tail insertion** — entries maintain their relative order during resize. This eliminates the circular reference. However, `HashMap` is still NOT thread-safe in Java 8 — you can still get lost updates and stale reads. The infinite loop specific bug is just gone.

**Bottom line:** Never use `HashMap` across threads. Period. Use `ConcurrentHashMap`.

---

### Q: How does ConcurrentHashMap work internally?

**Answer:**

**Java 7 — Segment-based locking:**
- The map was divided into 16 `Segment`s (by default), each being a mini `Hashtable` with its own lock.
- A put operation only locks the specific segment → 16 threads could write concurrently (to different segments).
- Reads were mostly lock-free (using `volatile` for the value field).
- Limitation: Max concurrency was bounded by the number of segments (default 16).

**Java 8+ — CAS + synchronized (per-bucket locking):**

The architecture was completely rewritten:
- Uses a single `Node<K,V>[]` table (like HashMap)
- **No segments** — locking is per-bucket (per array slot)
- Uses **CAS (Compare-And-Swap)** for the first insertion into an empty bucket
- Uses **`synchronized` on the first node** of a bucket for subsequent insertions to that bucket
- Reads are completely lock-free (using `volatile` reads and `Unsafe` operations)
- Supports treeification (same as HashMap — 8+ collisions → Red-Black Tree)

**Put operation (Java 8+):**
```
1. Compute hash → find bucket index
2. If bucket is empty:
   → CAS to insert the first node (no locking)
3. If bucket is occupied:
   → synchronized(firstNodeOfBucket) {
       traverse chain/tree, insert/update
     }
4. After insert, check if treeification needed
5. addCount() — atomically increment size using CAS + counter cells (LongAdder-like)
```

**Why this is genius:**
- Empty buckets (very common) use CAS — zero contention
- Occupied buckets lock only that bucket — two puts to different buckets never block each other
- Reads never block
- `size()` is approximate during concurrent modification (uses distributed counters for scalability)

**Compound operations (Java 8+):**
```java
// Atomic — replaces the check-then-act antipattern
cache.computeIfAbsent(key, k -> expensiveComputation(k));
cache.merge(key, 1, Integer::sum); // atomic increment
cache.compute(key, (k, v) -> v == null ? 1 : v + 1);
```

---

### Q: What is segment locking (Java 7) vs CAS + synchronized (Java 8)?

**Answer:**

**Java 7 Segments:**
```
ConcurrentHashMap
├── Segment[0] (ReentrantLock) → HashEntry[] → chains
├── Segment[1] (ReentrantLock) → HashEntry[] → chains
├── ...
└── Segment[15] (ReentrantLock) → HashEntry[] → chains

Max concurrency = 16 (concurrencyLevel)
```
- Each segment is an independent hash table with its own lock
- Put: hash key → find segment → lock segment → put in segment's array
- Memory overhead: 16 lock objects + 16 arrays + extra indirection

**Java 8 CAS + synchronized:**
```
ConcurrentHashMap
└── Node[] table
    ├── [0] null          → CAS for first insert
    ├── [1] Node → Node   → synchronized(first node) for updates
    ├── [2] null          → CAS
    ├── ...
    └── [n] TreeNode      → synchronized(root) for tree operations

Max concurrency = table.length (thousands of buckets)
```
- Granularity is per-bucket, not per-segment
- No pre-allocated lock objects — uses the Node itself as the monitor
- CAS for empty-bucket insertions (common case, very fast)
- Much better scalability under high concurrency

**Performance difference:** In benchmarks with 64+ threads, Java 8's `ConcurrentHashMap` is 2-4x faster than Java 7's due to finer-grained locking and less contention.

---

### Q: What is CopyOnWriteArrayList?

**Answer:**

`CopyOnWriteArrayList` (COWAL) creates a **new copy of the entire underlying array** on every write operation (add, set, remove). Reads operate on the current snapshot without any locking.

```java
// Internal mechanics (simplified):
public boolean add(E e) {
    synchronized (lock) {
        Object[] current = getArray();
        Object[] newArray = Arrays.copyOf(current, current.length + 1);
        newArray[current.length] = e;
        setArray(newArray); // volatile write — visible to all threads
        return true;
    }
}

public E get(int index) {
    return (E) getArray()[index]; // no lock, volatile read
}
```

**When to use:**
- Read-heavy, write-rare scenarios
- Small lists
- When you need snapshot iteration (iterator never throws CME)

**Production use case:** Event listener lists. In our event bus, listeners were registered at startup (rare writes) and dispatched to on every event (frequent reads). `CopyOnWriteArrayList` was perfect — no locking on the hot path (dispatch), and listener registration was rare and could afford the copy cost.

```java
private final List<EventListener> listeners = new CopyOnWriteArrayList<>();

public void addListener(EventListener listener) {
    listeners.add(listener); // rare, copies array — OK
}

public void fireEvent(Event event) {
    for (EventListener listener : listeners) { // no lock, snapshot iteration
        listener.onEvent(event);
    }
}
```

**When NOT to use:** If writes are frequent. Copying a 10,000-element array on every add is expensive. I've seen this misused in a queue scenario — terrible performance.

---

### Q: When to use Collections.synchronizedList vs concurrent collections?

**Answer:**

| Aspect | `Collections.synchronizedList()` | `CopyOnWriteArrayList` / Concurrent collections |
|--------|----------------------------------|------------------------------------------------|
| **Mechanism** | Wraps with `synchronized(mutex)` on every method | Lock-free reads / copy-on-write / CAS |
| **Iteration safety** | Must manually synchronize on the list during iteration | Iterator is safe (snapshot / weakly consistent) |
| **Performance** | Single lock → bottleneck under contention | Designed for concurrent access |
| **Compound operations** | Not atomic (check-then-act still unsafe) | Provides atomic compound operations |

**`Collections.synchronizedList()` pitfalls:**

```java
List<String> syncList = Collections.synchronizedList(new ArrayList<>());

// This is NOT safe — two separate synchronized calls:
if (!syncList.contains(item)) {  // lock, release
    syncList.add(item);           // lock, release — another thread might add in between
}

// Iteration MUST be manually synchronized:
synchronized (syncList) { // required!
    for (String s : syncList) {
        process(s);
    }
}
// Forgetting this synchronized block → ConcurrentModificationException
```

**My recommendation:**
- Need thread-safe `List` with mostly reads? → `CopyOnWriteArrayList`
- Need thread-safe `Map`? → `ConcurrentHashMap` (always, no exceptions)
- Need thread-safe `Queue`? → `ConcurrentLinkedQueue` or `BlockingQueue` implementations
- `Collections.synchronizedXxx()` → Only when you need a quick-and-dirty wrapper and understand the limitations. I avoid it in production code.

---

## ⚡ 5. Performance & Complexity

---

### Q: Time complexity of ArrayList, LinkedList, HashMap, TreeMap?

**Answer:**

**ArrayList:**

| Operation | Time Complexity |
|-----------|----------------|
| `get(index)` | O(1) |
| `add(element)` (end) | O(1) amortized |
| `add(index, element)` (middle) | O(n) — shifts elements |
| `remove(index)` | O(n) — shifts elements |
| `contains(object)` | O(n) — linear scan |
| `set(index, element)` | O(1) |

**LinkedList:**

| Operation | Time Complexity |
|-----------|----------------|
| `get(index)` | O(n) — traversal |
| `addFirst()` / `addLast()` | O(1) |
| `add(index, element)` | O(n) — find position + O(1) insert |
| `removeFirst()` / `removeLast()` | O(1) |
| `contains(object)` | O(n) |

**HashMap:**

| Operation | Average | Worst (Java 8+) |
|-----------|---------|-----------------|
| `put` / `get` / `remove` | O(1) | O(log n) per bucket |
| `containsKey` | O(1) | O(log n) |
| `containsValue` | O(n) | O(n) |

**TreeMap:**

| Operation | Time Complexity |
|-----------|----------------|
| `put` / `get` / `remove` | O(log n) — always |
| `firstKey()` / `lastKey()` | O(log n) |
| `subMap` / `headMap` / `tailMap` | O(log n) to find bounds, O(k) to iterate k elements |

---

### Q: Which collection is best for: Fast lookup, Frequent insert/delete, Ordered data?

**Answer:**

**Fast lookup by key:** `HashMap` — O(1) average. For 10 million entries, lookup time is essentially the same as for 100 entries (assuming good hash distribution). This is our go-to for caching, indexing, and any key-value scenario.

**Fast lookup by value/condition:** No single collection does this well. Use two maps (forward + reverse) or a database. In production, I maintain dual maps when I need bidirectional lookup:
```java
Map<String, User> byId = new HashMap<>();
Map<String, User> byEmail = new HashMap<>();
```

**Frequent insert/delete at ends:** `ArrayDeque` — O(1) for both ends, cache-friendly. Use for stacks, queues, and sliding windows.

**Frequent insert/delete in the middle (sorted order):** `TreeMap` / `TreeSet` — O(log n) insert, delete, and search. Used for leaderboards, time-series windowing, scheduled tasks.

**Ordered data (insertion order):** `ArrayList` (indexed) or `LinkedHashMap` / `LinkedHashSet` (for map/set with insertion order).

**Ordered data (sorted):** `TreeMap` / `TreeSet` — Red-Black Tree maintains sorted order automatically.

**Production decision framework:**
```
Need key-value? → Map
  Sorted keys? → TreeMap
  Insertion order? → LinkedHashMap
  Thread-safe? → ConcurrentHashMap
  Default? → HashMap

Need unique elements? → Set
  Sorted? → TreeSet
  Insertion order? → LinkedHashSet
  Default? → HashSet

Need ordered sequence? → List
  Random access? → ArrayList (99% of the time)
  Queue/Stack behavior? → ArrayDeque
```

---

### Q: Difference between O(1) vs O(log n) collections?

**Answer:**

| Entries (n) | O(1) ops | O(log n) ops | O(n) ops |
|-------------|----------|-------------|----------|
| 100 | 1 | ~7 | 100 |
| 10,000 | 1 | ~14 | 10,000 |
| 1,000,000 | 1 | ~20 | 1,000,000 |
| 100,000,000 | 1 | ~27 | 100,000,000 |

**O(1) — `HashMap`, `HashSet`:**
- Constant time regardless of size
- But: higher constant factor (hash computation, potential collision handling)
- No ordering guarantee

**O(log n) — `TreeMap`, `TreeSet`:**
- Grows very slowly — even 100M entries only need ~27 comparisons
- Provides sorted order, range queries, `firstKey()`, `lastKey()`, `subMap()`
- Each operation involves tree traversal and comparisons

**When O(log n) beats O(1) in practice:**
- When you need sorted iteration (HashMap requires sorting: O(n log n) vs TreeMap: O(n))
- When you need range queries (`subMap`, `headMap`, `tailMap` — impossible with HashMap)
- When key objects have expensive `hashCode()` but cheap `compareTo()`

**Real example:** In a rate limiter, I used `TreeMap<Instant, Integer>` to store request timestamps. To count requests in the last minute: `treeMap.subMap(now.minusMinutes(1), now).size()` — O(log n + k) where k is the result size. With `HashMap`, this would require scanning all entries — O(n).

---

## 🌳 6. Sorted & Ordered Collections

---

### Q: Difference between HashMap vs TreeMap?

**Answer:**

| Feature | HashMap | TreeMap |
|---------|---------|--------|
| **Internal structure** | Hash table (array + linked list/tree) | Red-Black Tree |
| **Ordering** | No guaranteed order | Sorted by key (natural or Comparator) |
| **Null key** | One null key allowed | No null key (NPE — can't compare null) |
| **Performance** | O(1) average for get/put | O(log n) for get/put |
| **Range queries** | Not supported | `subMap()`, `headMap()`, `tailMap()`, `firstKey()`, `lastKey()` |
| **Memory** | Less per entry (Node: key, value, hash, next) | More per entry (TreeNode: key, value, left, right, parent, color) |
| **Implements** | `Map` | `NavigableMap` → `SortedMap` → `Map` |

**Decision criteria:**
- Default choice → `HashMap`
- Need sorted keys → `TreeMap`
- Need range operations → `TreeMap`
- Need first/last element → `TreeMap`
- Need maximum throughput → `HashMap`

```java
// HashMap — fast user lookup
Map<String, User> userCache = new HashMap<>();

// TreeMap — time-based event store with range queries
TreeMap<LocalDateTime, Event> events = new TreeMap<>();
// Get all events today:
events.subMap(today.atStartOfDay(), today.plusDays(1).atStartOfDay());
// Get most recent event:
events.lastEntry();
```

---

### Q: What is NavigableMap / SortedMap?

**Answer:**

**`SortedMap<K,V>`** — Guarantees keys are in sorted order. Provides:
- `firstKey()`, `lastKey()`
- `headMap(toKey)`, `tailMap(fromKey)`, `subMap(from, to)`
- `comparator()` — returns the Comparator, or null if natural ordering

**`NavigableMap<K,V>`** (extends `SortedMap`, added in Java 6) — Adds navigation methods:
- `lowerEntry(key)` / `lowerKey(key)` — strictly less than
- `floorEntry(key)` / `floorKey(key)` — less than or equal
- `ceilingEntry(key)` / `ceilingKey(key)` — greater than or equal
- `higherEntry(key)` / `higherKey(key)` — strictly greater than
- `pollFirstEntry()` / `pollLastEntry()` — retrieve and remove
- `descendingMap()` — reverse-order view

**Production use case — finding the nearest config threshold:**
```java
NavigableMap<Integer, String> pricingTiers = new TreeMap<>();
pricingTiers.put(0, "free");
pricingTiers.put(100, "basic");
pricingTiers.put(500, "premium");
pricingTiers.put(2000, "enterprise");

int userUsage = 350;
String tier = pricingTiers.floorEntry(userUsage).getValue(); // "basic"
```

---

### Q: How does TreeMap maintain order?

**Answer:**

TreeMap uses a **Red-Black Tree** — a self-balancing binary search tree. Every entry is a tree node with left child, right child, parent, and a color (red or black).

**Properties:**
1. Every node is either red or black
2. Root is always black
3. No two consecutive red nodes (red node's children must be black)
4. Every path from root to null leaf has the same number of black nodes

After every insertion or deletion, the tree performs **rotations** (left/right) and **color flips** to maintain these properties. This guarantees the tree height is at most `2 × log₂(n+1)`, ensuring O(log n) operations.

**Ordering depends on:**
1. **Natural ordering** (key implements `Comparable`) — used by default
2. **Custom `Comparator`** — provided at construction

```java
// Natural ordering — keys must be Comparable
TreeMap<String, Integer> map = new TreeMap<>(); // alphabetical

// Custom comparator — reverse order
TreeMap<String, Integer> map = new TreeMap<>(Comparator.reverseOrder());

// Complex comparator
TreeMap<Employee, String> map = new TreeMap<>(
    Comparator.comparing(Employee::getDepartment)
              .thenComparing(Employee::getName)
);
```

**In-order traversal** of the Red-Black Tree gives sorted output, which is why `TreeMap.entrySet()` iteration is always in key order.

---

### Q: What happens if keys are not Comparable?

**Answer:**

If you don't provide a `Comparator` at construction time, TreeMap assumes keys implement `Comparable`. If they don't:

```java
TreeMap<MyObject, String> map = new TreeMap<>();
map.put(new MyObject(), "value");
// Throws ClassCastException: MyObject cannot be cast to java.lang.Comparable
```

The exception happens at runtime, not compile time (because `TreeMap` accepts any `K`).

**Solutions:**

```java
// Option 1: Make the key class implement Comparable
public class MyObject implements Comparable<MyObject> {
    @Override
    public int compareTo(MyObject other) { ... }
}

// Option 2: Provide a Comparator (preferred — more flexible)
TreeMap<MyObject, String> map = new TreeMap<>(
    Comparator.comparing(MyObject::getName)
);
```

**Subtle gotcha:** Even if the first key implements `Comparable`, inserting a second key of a different type will throw `ClassCastException` during comparison. The first `put()` on an empty TreeMap doesn't compare anything (Java 8+), so it silently succeeds even if the key isn't `Comparable`. The second `put()` triggers the comparison and fails.

---

## 🧠 7. Advanced Concepts (Senior-Level)

---

### Q: What is an immutable collection?

**Answer:**

An immutable collection is one that **cannot be modified after creation** — no add, remove, set, or clear. Any attempt throws `UnsupportedOperationException`.

**Key distinction — truly immutable:**
- The collection itself is immutable
- It does NOT contain any references to the original mutable collection
- No one else can modify it through a back-door reference

```java
// Java 9+ — truly immutable
List<String> immutable = List.of("a", "b", "c");
immutable.add("d"); // UnsupportedOperationException
```

**Why immutability matters in production:**
1. **Thread safety for free** — immutable objects need no synchronization
2. **Safe sharing** — pass to any method without worry about mutation
3. **Cache-friendly** — can be cached and reused without defensive copies
4. **Predictable** — no surprising side effects

I make collections immutable at API boundaries. If a service method returns a list, it should be immutable to prevent callers from accidentally corrupting internal state.

---

### Q: How to create immutable collections in Java?

**Answer:**

```java
// 1. Java 9+ factory methods (RECOMMENDED)
List<String> list = List.of("a", "b", "c");
Set<String> set = Set.of("x", "y", "z");
Map<String, Integer> map = Map.of("a", 1, "b", 2);
Map<String, Integer> map2 = Map.ofEntries(
    Map.entry("a", 1),
    Map.entry("b", 2)
);

// 2. Java 10+ — copy of existing collection
List<String> copy = List.copyOf(mutableList);
Set<String> copy2 = Set.copyOf(mutableSet);

// 3. Java 10+ Collectors
List<String> collected = stream.collect(Collectors.toUnmodifiableList());

// 4. Java 16+
List<String> collected = stream.toList(); // returns unmodifiable list

// 5. Guava (pre-Java 9)
ImmutableList<String> list = ImmutableList.of("a", "b");
ImmutableMap<String, Integer> map = ImmutableMap.builder()
    .put("a", 1)
    .put("b", 2)
    .build();

// 6. Collections.unmodifiableList (NOT truly immutable — see next question)
List<String> unmodifiable = Collections.unmodifiableList(mutableList);
```

**Important:** `List.of()` and friends reject `null` elements — they throw `NullPointerException`. If you need nulls (which you probably shouldn't), you can't use these factory methods.

---

### Q: What is unmodifiable vs immutable?

**Answer:**

This is a critical distinction:

**Unmodifiable** — A wrapper that prevents modification through the wrapper reference, but the underlying collection can still be modified through the original reference.

```java
List<String> original = new ArrayList<>(Arrays.asList("a", "b"));
List<String> unmodifiable = Collections.unmodifiableList(original);

unmodifiable.add("c"); // UnsupportedOperationException — good
original.add("c");     // works! And now unmodifiable also shows "c"
System.out.println(unmodifiable); // [a, b, c] — SURPRISE!
```

**Immutable** — No one can modify the collection. No back-door reference exists.

```java
List<String> immutable = List.of("a", "b"); // no underlying mutable list exists
// OR
List<String> immutable = List.copyOf(original); // defensive copy — detached from original
```

**Production rule:** I never return `Collections.unmodifiableList(this.internalList)` from a getter — that's a false sense of security. The caller can't mutate it, but internal code still can, and the wrapper will reflect those changes. Instead:
```java
// Return a truly detached immutable copy
public List<String> getItems() {
    return List.copyOf(this.items);
}
```

---

### Q: What is WeakHashMap?

**Answer:**

`WeakHashMap` stores keys as **weak references**. When a key has no strong references remaining, the garbage collector can collect it, and the entry is automatically removed from the map.

```java
WeakHashMap<Object, String> cache = new WeakHashMap<>();
Object key = new Object();
cache.put(key, "value");

System.out.println(cache.size()); // 1

key = null; // no more strong references to the key
System.gc(); // hint to GC

// After GC runs:
System.out.println(cache.size()); // 0 — entry was automatically removed
```

**How it works internally:** Keys are wrapped in `WeakReference` objects. The GC enqueues collected weak references into a `ReferenceQueue`. On every `get()`, `put()`, or `size()` call, `WeakHashMap` polls this queue and removes expired entries.

**Production use case — class metadata cache:**
```java
// Cache computed metadata about classes
// When a ClassLoader is unloaded (e.g., in app servers during redeploy),
// the Class objects become weakly reachable, and entries are cleaned up automatically
private static final WeakHashMap<Class<?>, Metadata> metadataCache = new WeakHashMap<>();
```

**Gotchas:**
- NOT thread-safe — wrap with `Collections.synchronizedMap()` if needed
- `String` literals and small `Integer`s are always strongly referenced (interned), so they'll never be collected — don't use them as keys if you expect cleanup
- Only keys are weak — values are strong references. If a value strongly references its key, the key will never be collected (memory leak)

---

### Q: What is IdentityHashMap?

**Answer:**

`IdentityHashMap` uses **reference equality (`==`)** instead of `equals()` for comparing keys. Two keys are considered the same only if they are the exact same object in memory.

```java
IdentityHashMap<String, Integer> map = new IdentityHashMap<>();
String s1 = new String("hello");
String s2 = new String("hello");

map.put(s1, 1);
map.put(s2, 2);

System.out.println(map.size()); // 2! Even though s1.equals(s2) is true
// In a normal HashMap, size would be 1
```

**Use cases:**
1. **Serialization frameworks** — tracking which objects have already been serialized (by identity, not equality) to handle circular references
2. **Object graph traversal** — visited-set that distinguishes between two equal but distinct objects
3. **Proxy/wrapper tracking** — maintaining metadata about specific object instances

**Internally:** Uses a flat `Object[]` array with linear probing (not chaining). Keys and values are stored adjacently in the array: `[key1, val1, key2, val2, ...]`. This is more cache-friendly than Node-based structures.

---

### Q: What is EnumMap?

**Answer:**

`EnumMap` is a specialized `Map` implementation for enum keys. It's backed by a simple array indexed by `enum.ordinal()`.

```java
enum Status { PENDING, ACTIVE, SUSPENDED, CLOSED }

EnumMap<Status, List<User>> usersByStatus = new EnumMap<>(Status.class);
usersByStatus.put(Status.ACTIVE, activeUsers);
usersByStatus.put(Status.PENDING, pendingUsers);
```

**Why it's excellent:**
- **Blazing fast** — `get(key)` is `array[key.ordinal()]` — direct array access, O(1) with minimal constant
- **Compact** — array size = number of enum constants, no hashing overhead, no wasted buckets
- **Maintains natural enum order** (declaration order) during iteration
- **Null values** allowed (but not null keys)
- **Type-safe** — constructor requires the enum class

**Performance vs HashMap:** `EnumMap` is roughly 2-4x faster than `HashMap` for enum keys because there's no hashing, no collision handling, and perfect array locality.

**Production use:** I use `EnumMap` everywhere I have enum keys — feature flags by environment, config by region, handlers by event type. It's one of those "free performance wins" that I always apply.

```java
// Dispatch handlers by event type
private static final EnumMap<EventType, EventHandler> HANDLERS = new EnumMap<>(EventType.class);
static {
    HANDLERS.put(EventType.ORDER_CREATED, new OrderCreatedHandler());
    HANDLERS.put(EventType.PAYMENT_RECEIVED, new PaymentHandler());
}
```

---

### Q: What is BlockingQueue?

**Answer:**

`BlockingQueue` is a `Queue` that supports **blocking operations** — it waits (blocks the thread) when you try to take from an empty queue or put into a full queue.

| Method Type | Throws Exception | Returns Special Value | Blocks | Times Out |
|------------|-----------------|----------------------|--------|-----------|
| **Insert** | `add(e)` | `offer(e)` → false | `put(e)` | `offer(e, timeout, unit)` |
| **Remove** | `remove()` | `poll()` → null | `take()` | `poll(timeout, unit)` |
| **Examine** | `element()` | `peek()` → null | N/A | N/A |

**Implementations:**

| Class | Bounded | Backing | Use Case |
|-------|---------|---------|----------|
| `ArrayBlockingQueue` | Yes (fixed) | Circular array | Fixed-size producer-consumer buffer |
| `LinkedBlockingQueue` | Optional | Linked nodes | Unbounded/bounded work queues |
| `PriorityBlockingQueue` | No (unbounded) | Heap | Priority-based processing |
| `SynchronousQueue` | Zero capacity | None | Direct hand-off between threads |
| `DelayQueue` | No | Heap | Scheduled/delayed tasks |

**Production use — producer-consumer pattern:**

```java
BlockingQueue<Task> taskQueue = new ArrayBlockingQueue<>(1000);

// Producer threads
taskQueue.put(new Task(orderId)); // blocks if queue is full — backpressure!

// Consumer threads
Task task = taskQueue.take(); // blocks if queue is empty — waits for work
process(task);
```

This is the backbone of most async processing in Java. `ThreadPoolExecutor` itself uses a `BlockingQueue` internally. I've used `ArrayBlockingQueue` to implement backpressure in a message processing pipeline — when the queue fills up, producers slow down naturally instead of overwhelming the system.

---

### Q: Difference between Queue vs Deque?

**Answer:**

| Feature | Queue | Deque (Double-Ended Queue) |
|---------|-------|---------------------------|
| **Insert** | Only at tail | At head or tail |
| **Remove** | Only from head | From head or tail |
| **Usage** | FIFO processing | FIFO, LIFO (stack), or both ends |
| **Key methods** | `offer()`, `poll()`, `peek()` | `offerFirst()`, `offerLast()`, `pollFirst()`, `pollLast()` |

`Deque` extends `Queue` and adds operations at both ends.

```java
// Queue — FIFO task processing
Queue<Task> tasks = new LinkedList<>();
tasks.offer(task1); // enqueue
tasks.poll();        // dequeue from front

// Deque as a stack (LIFO) — better than java.util.Stack
Deque<String> stack = new ArrayDeque<>();
stack.push("a");    // addFirst
stack.pop();         // removeFirst

// Deque as a sliding window
Deque<Integer> window = new ArrayDeque<>();
window.addLast(newElement);      // add to end
if (window.size() > windowSize) {
    window.removeFirst();         // remove from front
}
```

**Best implementation:** `ArrayDeque` — faster than `LinkedList` for both stack and queue operations due to cache locality. The only time I use `LinkedList` as a Deque is when I also need `null` elements (ArrayDeque doesn't allow null).

---

### Q: What is PriorityQueue?

**Answer:**

`PriorityQueue` is a **min-heap** based queue where elements are dequeued in priority order (smallest first by natural ordering, or by a custom Comparator).

```java
// Min-heap (default) — smallest element comes out first
PriorityQueue<Integer> minHeap = new PriorityQueue<>();
minHeap.offer(30);
minHeap.offer(10);
minHeap.offer(20);
minHeap.poll(); // 10 (smallest)
minHeap.poll(); // 20

// Max-heap
PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Comparator.reverseOrder());

// Custom priority
PriorityQueue<Task> taskQueue = new PriorityQueue<>(
    Comparator.comparing(Task::getPriority)
              .thenComparing(Task::getCreatedAt)
);
```

**Key characteristics:**
- **Not sorted** — only guarantees the head is the min/max. Internal order is a heap (partially ordered).
- **O(log n)** for offer and poll
- **O(1)** for peek
- **O(n)** for remove(Object) and contains()
- **Not thread-safe** — use `PriorityBlockingQueue` for concurrent access
- **No null** elements allowed
- **Unbounded** — grows dynamically like ArrayList

**Production use — Top-K problem:**
```java
// Find top 5 highest-scoring users from a stream of millions
PriorityQueue<User> topK = new PriorityQueue<>(
    Comparator.comparingInt(User::getScore) // min-heap
);
for (User user : allUsers) {
    topK.offer(user);
    if (topK.size() > 5) {
        topK.poll(); // remove lowest score
    }
}
// topK now contains the 5 highest-scoring users
// O(n log k) — much better than sorting the entire list O(n log n)
```

---

## 🧪 8. Practical Coding / Scenario Questions

---

### Q: How to remove duplicates from a list?

**Answer:**

```java
// Method 1: Using LinkedHashSet (preserves insertion order) — my go-to
List<String> list = new ArrayList<>(Arrays.asList("a", "b", "a", "c", "b"));
List<String> unique = new ArrayList<>(new LinkedHashSet<>(list));
// [a, b, c]

// Method 2: Using Stream (Java 8+)
List<String> unique = list.stream()
    .distinct()
    .collect(Collectors.toList());

// Method 3: If order doesn't matter
Set<String> unique = new HashSet<>(list);

// Method 4: Remove duplicates by a specific field (e.g., deduplicate users by email)
List<User> uniqueUsers = users.stream()
    .collect(Collectors.collectingAndThen(
        Collectors.toMap(User::getEmail, Function.identity(), (a, b) -> a, LinkedHashMap::new),
        m -> new ArrayList<>(m.values())
    ));

// Simpler with a seen-set:
Set<String> seen = new HashSet<>();
List<User> uniqueUsers = users.stream()
    .filter(u -> seen.add(u.getEmail()))
    .collect(Collectors.toList());
```

**Production note:** Method 4 (dedup by field) comes up constantly. In one project, we received duplicate webhook events and needed to deduplicate by event ID while preserving order. The `seen.add()` filter pattern is clean and efficient.

---

### Q: How to find frequency of elements?

**Answer:**

```java
List<String> words = Arrays.asList("apple", "banana", "apple", "cherry", "banana", "apple");

// Method 1: Stream with Collectors.groupingBy + counting (RECOMMENDED)
Map<String, Long> frequency = words.stream()
    .collect(Collectors.groupingBy(Function.identity(), Collectors.counting()));
// {apple=3, banana=2, cherry=1}

// Method 2: merge() — elegant and works with any map
Map<String, Integer> freq = new HashMap<>();
for (String word : words) {
    freq.merge(word, 1, Integer::sum);
}

// Method 3: compute() — more verbose but flexible
Map<String, Integer> freq = new HashMap<>();
for (String word : words) {
    freq.compute(word, (k, v) -> v == null ? 1 : v + 1);
}

// Method 4: getOrDefault — old-school
Map<String, Integer> freq = new HashMap<>();
for (String word : words) {
    freq.put(word, freq.getOrDefault(word, 0) + 1);
}

// Find most frequent element:
String mostFrequent = frequency.entrySet().stream()
    .max(Map.Entry.comparingByValue())
    .map(Map.Entry::getKey)
    .orElse(null);
```

I prefer `merge()` for imperative code and `groupingBy` + `counting()` for stream-based code.

---

### Q: How to sort objects using Comparator?

**Answer:**

```java
public class Employee {
    private String name;
    private String department;
    private int salary;
    private LocalDate joinDate;
}

// Single field sort
employees.sort(Comparator.comparing(Employee::getSalary));

// Reverse order
employees.sort(Comparator.comparing(Employee::getSalary).reversed());

// Multi-field sort (department ASC, then salary DESC, then name ASC)
employees.sort(
    Comparator.comparing(Employee::getDepartment)
              .thenComparing(Employee::getSalary, Comparator.reverseOrder())
              .thenComparing(Employee::getName)
);

// Null-safe sorting
employees.sort(
    Comparator.comparing(Employee::getDepartment, 
                         Comparator.nullsLast(Comparator.naturalOrder()))
);

// Sort with stream (returns new list, original unchanged)
List<Employee> sorted = employees.stream()
    .sorted(Comparator.comparing(Employee::getSalary).reversed())
    .collect(Collectors.toList());
```

**Production tip:** I define common comparators as static constants in the entity class or a utility class:
```java
public class Employee {
    public static final Comparator<Employee> BY_SALARY_DESC = 
        Comparator.comparing(Employee::getSalary).reversed();
    public static final Comparator<Employee> BY_DEPT_THEN_NAME = 
        Comparator.comparing(Employee::getDepartment)
                  .thenComparing(Employee::getName);
}
```

---

### Q: How to find top K elements?

**Answer:**

```java
// Method 1: PriorityQueue (min-heap) — O(n log k), best for streaming/large data
public static <T> List<T> topK(Collection<T> items, int k, Comparator<T> comparator) {
    PriorityQueue<T> minHeap = new PriorityQueue<>(k, comparator);
    for (T item : items) {
        minHeap.offer(item);
        if (minHeap.size() > k) {
            minHeap.poll(); // remove smallest
        }
    }
    List<T> result = new ArrayList<>(minHeap);
    result.sort(comparator.reversed()); // optional: sort in descending order
    return result;
}

// Usage: top 5 highest salary employees
List<Employee> top5 = topK(employees, 5, Comparator.comparingInt(Employee::getSalary));

// Method 2: Stream — O(n log n), simpler but sorts everything
List<Employee> top5 = employees.stream()
    .sorted(Comparator.comparingInt(Employee::getSalary).reversed())
    .limit(5)
    .collect(Collectors.toList());

// Method 3: TreeMap/TreeSet with bounded size — good for maintaining a running top-K
TreeSet<Employee> topK = new TreeSet<>(Comparator.comparingInt(Employee::getSalary));
for (Employee e : employees) {
    topK.add(e);
    if (topK.size() > k) {
        topK.pollFirst(); // remove lowest
    }
}
```

**When to use which:**
- **PriorityQueue** — best for large data sets, streaming data, or when n >> k. Memory: O(k).
- **Stream sort + limit** — cleanest code, acceptable when n is manageable. Memory: O(n).
- **TreeSet** — when you also need to query "is this element in the top K?" during processing.

---

### Q: How to implement LRU cache using collections?

**Answer:**

```java
// Method 1: LinkedHashMap with accessOrder=true (simplest, production-ready for single-threaded)
public class LRUCache<K, V> extends LinkedHashMap<K, V> {
    private final int maxSize;

    public LRUCache(int maxSize) {
        super(maxSize * 4 / 3 + 1, 0.75f, true); // accessOrder=true
        this.maxSize = maxSize;
    }

    @Override
    protected boolean removeEldestEntry(Map.Entry<K, V> eldest) {
        return size() > maxSize;
    }
}

// Usage
LRUCache<String, UserProfile> cache = new LRUCache<>(1000);
cache.put("user-1", profile1);
cache.get("user-1"); // moves to tail (most recently used)
cache.put("user-2", profile2);
// When size > 1000, least recently accessed entry is evicted
```

```java
// Method 2: Thread-safe LRU using ConcurrentHashMap + ConcurrentLinkedDeque
// (for production multi-threaded scenarios)
public class ConcurrentLRUCache<K, V> {
    private final int maxSize;
    private final ConcurrentHashMap<K, V> map;
    private final Deque<K> accessOrder = new ConcurrentLinkedDeque<>();

    public ConcurrentLRUCache(int maxSize) {
        this.maxSize = maxSize;
        this.map = new ConcurrentHashMap<>(maxSize);
    }

    public V get(K key) {
        V value = map.get(key);
        if (value != null) {
            accessOrder.remove(key);
            accessOrder.addLast(key);
        }
        return value;
    }

    public void put(K key, V value) {
        if (map.containsKey(key)) {
            accessOrder.remove(key);
        } else if (map.size() >= maxSize) {
            K eldest = accessOrder.pollFirst();
            if (eldest != null) map.remove(eldest);
        }
        map.put(key, value);
        accessOrder.addLast(key);
    }
}
```

**Real production choice:** For serious caching needs, I use **Caffeine** library (successor to Guava Cache). It provides a high-performance, near-optimal LRU/W-TinyLFU cache with O(1) operations, async loading, stats, and eviction listeners. But understanding the LinkedHashMap approach is essential for interviews and for cases where you can't add dependencies.

---

### Q: Reverse a list efficiently?

**Answer:**

```java
// Method 1: Collections.reverse() — in-place, O(n)
List<String> list = new ArrayList<>(Arrays.asList("a", "b", "c", "d"));
Collections.reverse(list);
// [d, c, b, a]

// Method 2: Stream — returns new list, original unchanged
List<String> reversed = IntStream.rangeClosed(1, list.size())
    .mapToObj(i -> list.get(list.size() - i))
    .collect(Collectors.toList());

// Method 3: ListIterator — manual reverse into new list
List<String> reversed = new ArrayList<>(list.size());
ListIterator<String> it = list.listIterator(list.size());
while (it.hasPrevious()) {
    reversed.add(it.previous());
}

// Method 4: Java 21+ SequencedCollection
List<String> reversed = list.reversed(); // returns a reversed VIEW
```

**Production recommendation:** `Collections.reverse()` for in-place. If you need a reversed view without copying, Java 21's `reversed()` is ideal. For streams, just `Collections.reverse()` on the collected result — it's cleaner than reverse-indexing in a stream.

---

### Q: Merge two sorted lists?

**Answer:**

```java
// Two-pointer merge — O(n + m), same as merge step in merge sort
public static <T extends Comparable<T>> List<T> mergeSorted(List<T> list1, List<T> list2) {
    List<T> result = new ArrayList<>(list1.size() + list2.size());
    int i = 0, j = 0;

    while (i < list1.size() && j < list2.size()) {
        if (list1.get(i).compareTo(list2.get(j)) <= 0) {
            result.add(list1.get(i++));
        } else {
            result.add(list2.get(j++));
        }
    }

    while (i < list1.size()) result.add(list1.get(i++));
    while (j < list2.size()) result.add(list2.get(j++));

    return result;
}
```

**Alternative (simpler but slower):**
```java
List<Integer> merged = new ArrayList<>(list1);
merged.addAll(list2);
Collections.sort(merged); // O((n+m) log(n+m)) — worse than two-pointer
```

Use the two-pointer approach in production when dealing with pre-sorted data (e.g., merging sorted pages from a database or merging sorted files in external sort).

---

### Q: Detect cycle in LinkedList?

**Answer:**

**Floyd's Tortoise and Hare algorithm — O(n) time, O(1) space:**

```java
public boolean hasCycle(ListNode head) {
    if (head == null || head.next == null) return false;

    ListNode slow = head;       // moves 1 step
    ListNode fast = head;       // moves 2 steps

    while (fast != null && fast.next != null) {
        slow = slow.next;
        fast = fast.next.next;
        if (slow == fast) return true; // they meet → cycle exists
    }
    return false; // fast reached end → no cycle
}

// To find the START of the cycle:
public ListNode detectCycleStart(ListNode head) {
    ListNode slow = head, fast = head;

    while (fast != null && fast.next != null) {
        slow = slow.next;
        fast = fast.next.next;
        if (slow == fast) {
            // Reset slow to head, move both at same speed
            slow = head;
            while (slow != fast) {
                slow = slow.next;
                fast = fast.next;
            }
            return slow; // cycle start
        }
    }
    return null;
}
```

**Why this works:** When the fast pointer moves 2x faster, if there's a cycle, they will eventually meet inside the cycle. The math proves that after meeting, the distance from the head to the cycle start equals the distance from the meeting point to the cycle start (going around the cycle).

---

### Q: Find first non-repeating character?

**Answer:**

```java
// Method 1: LinkedHashMap — preserves insertion order, O(n)
public static Character firstNonRepeating(String str) {
    Map<Character, Integer> freq = new LinkedHashMap<>();
    for (char c : str.toCharArray()) {
        freq.merge(c, 1, Integer::sum);
    }
    for (Map.Entry<Character, Integer> entry : freq.entrySet()) {
        if (entry.getValue() == 1) return entry.getKey();
    }
    return null;
}

// Method 2: Array-based (for lowercase ASCII only — faster)
public static char firstNonRepeating(String str) {
    int[] freq = new int[26];
    for (char c : str.toCharArray()) freq[c - 'a']++;
    for (char c : str.toCharArray()) {
        if (freq[c - 'a'] == 1) return c;
    }
    return '\0';
}

// Method 3: Stream (Java 8+)
Character result = str.chars()
    .mapToObj(c -> (char) c)
    .collect(Collectors.groupingBy(Function.identity(), LinkedHashMap::new, Collectors.counting()))
    .entrySet().stream()
    .filter(e -> e.getValue() == 1)
    .map(Map.Entry::getKey)
    .findFirst()
    .orElse(null);
```

**Why `LinkedHashMap`?** A regular `HashMap` doesn't preserve insertion order, so iterating over it won't give you the *first* non-repeating character. `LinkedHashMap` iterates in insertion order, so the first entry with count 1 is our answer.

---

## 🏗️ 9. Design-Level Questions (Very Important for Senior Engineers)

---

### Q: Which collection would you choose for: High throughput system? Low latency system?

**Answer:**

**High throughput (maximize operations per second):**
- **`ConcurrentHashMap`** for shared state — bucket-level locking allows massive parallelism
- **`ArrayDeque`** for per-thread queues — cache-friendly, no lock contention
- **Pre-sized `ArrayList`** for batch processing — minimize allocations and GC
- **`EnumMap`** when keys are enums — array-indexed, zero hashing overhead
- Avoid: `Collections.synchronizedMap()` (global lock kills throughput), `TreeMap` (O(log n) adds up at scale)

**Low latency (minimize worst-case response time):**
- **`HashMap` with pre-sized capacity** — avoid resize spikes during request processing
- **`CopyOnWriteArrayList`** for read-heavy config/listener lists — zero-cost reads
- **Avoid `TreeMap` in hot paths** — O(log n) with pointer chasing → unpredictable latency due to cache misses
- **Avoid large `ArrayList` resizing** in request scope — causes GC pauses
- **Object pooling** for frequently created collections — reduces allocation pressure

**Real example:** In a payment gateway handling 50K TPS:
- Session data: `ConcurrentHashMap` with initial capacity set to expected concurrent sessions
- Per-request processing: pre-sized `ArrayList` or reusable thread-local lists
- Config/feature flags: `CopyOnWriteArrayList` (updated once per minute, read on every request)
- We profiled and found that 15% of young-gen GC was from short-lived `ArrayList` instances in request handlers. Switching to pre-sized thread-local lists cut that significantly.

---

### Q: How would you design: In-memory cache?

**Answer:**

**Requirements-first approach:**

1. **Basic LRU cache (single-threaded):**
```java
LinkedHashMap<K, V> cache = new LinkedHashMap<>(capacity, 0.75f, true) {
    protected boolean removeEldestEntry(Map.Entry<K, V> eldest) {
        return size() > capacity;
    }
};
```

2. **Production-grade in-memory cache (what I'd actually build/use):**

**Option A: Caffeine (recommended)**
```java
Cache<String, UserProfile> cache = Caffeine.newBuilder()
    .maximumSize(10_000)
    .expireAfterWrite(Duration.ofMinutes(15))
    .expireAfterAccess(Duration.ofMinutes(5))
    .recordStats()
    .removalListener((key, value, cause) -> log.info("Evicted: {} reason: {}", key, cause))
    .build();
```

**Option B: Custom implementation (when library isn't an option)**

Key design decisions:
- **Eviction policy:** LRU (LinkedHashMap), LFU, or TTL-based
- **Concurrency:** `ConcurrentHashMap` + striped locks for LRU ordering
- **Bounded size:** Hard cap to prevent OOM
- **Expiration:** Background thread with `DelayQueue` or timestamp checks on access
- **Statistics:** Hit rate, miss rate, eviction count for monitoring

```java
public class ProductionCache<K, V> {
    private final ConcurrentHashMap<K, CacheEntry<V>> store;
    private final int maxSize;
    private final Duration ttl;
    private final ScheduledExecutorService cleaner;

    record CacheEntry<V>(V value, Instant expiresAt) {
        boolean isExpired() { return Instant.now().isAfter(expiresAt); }
    }

    public V get(K key) {
        CacheEntry<V> entry = store.get(key);
        if (entry == null) return null;
        if (entry.isExpired()) {
            store.remove(key);
            return null;
        }
        return entry.value();
    }

    public void put(K key, V value) {
        if (store.size() >= maxSize) {
            evict(); // remove expired entries, then oldest if needed
        }
        store.put(key, new CacheEntry<>(value, Instant.now().plus(ttl)));
    }
}
```

**Design considerations I raise in interviews:**
- **Memory pressure:** Cache can grow and cause OOM → always set max size, consider `SoftReference` values
- **Thundering herd:** When a popular cache entry expires, 100 threads all try to recompute → use `computeIfAbsent` or single-flight pattern
- **Serialization:** If caching across JVMs, need to think about serialization cost
- **Monitoring:** Without hit/miss rate monitoring, you're flying blind

---

### Q: How would you design: Rate limiter?

**Answer:**

**Sliding window counter using `TreeMap`:**

```java
public class RateLimiter {
    private final int maxRequests;
    private final Duration window;
    private final ConcurrentHashMap<String, TreeMap<Long, Integer>> clientWindows = new ConcurrentHashMap<>();

    public RateLimiter(int maxRequests, Duration window) {
        this.maxRequests = maxRequests;
        this.window = window;
    }

    public synchronized boolean allowRequest(String clientId) {
        long now = System.currentTimeMillis();
        long windowStart = now - window.toMillis();

        TreeMap<Long, Integer> timestamps = clientWindows.computeIfAbsent(clientId, k -> new TreeMap<>());

        // Remove entries outside the window
        timestamps.headMap(windowStart).clear();

        // Count requests in the window
        int count = timestamps.values().stream().mapToInt(Integer::intValue).sum();

        if (count >= maxRequests) {
            return false;
        }

        timestamps.merge(now, 1, Integer::sum);
        return true;
    }
}
```

**Token bucket using `ArrayDeque`:**

```java
public class TokenBucketRateLimiter {
    private final int maxTokens;
    private final double refillRatePerMs;
    private double tokens;
    private long lastRefillTimestamp;

    public TokenBucketRateLimiter(int maxTokens, int refillPerSecond) {
        this.maxTokens = maxTokens;
        this.refillRatePerMs = refillPerSecond / 1000.0;
        this.tokens = maxTokens;
        this.lastRefillTimestamp = System.currentTimeMillis();
    }

    public synchronized boolean tryConsume() {
        refill();
        if (tokens >= 1) {
            tokens--;
            return true;
        }
        return false;
    }

    private void refill() {
        long now = System.currentTimeMillis();
        double newTokens = (now - lastRefillTimestamp) * refillRatePerMs;
        tokens = Math.min(maxTokens, tokens + newTokens);
        lastRefillTimestamp = now;
    }
}
```

**Collections choice rationale:**
- `TreeMap` — gives natural windowing with `headMap()` / `subMap()` for sliding window
- `ConcurrentHashMap` — client-level isolation, thread-safe per-client state
- For distributed rate limiting, I'd use Redis with Lua scripts, but the data structure concepts are the same

---

### Q: When NOT to use HashMap?

**Answer:**

1. **When you need sorted/ordered keys** → Use `TreeMap`
2. **When keys are enums** → Use `EnumMap` (faster and more memory-efficient)
3. **When you need concurrent access** → Use `ConcurrentHashMap`
4. **When key objects have mutable fields used in hashCode/equals** → HashMap will "lose" entries when hash changes after insertion
5. **When you need identity-based comparison** → Use `IdentityHashMap`
6. **When keys are weakly held for GC purposes** → Use `WeakHashMap`
7. **When memory is extremely constrained** → HashMap has ~48 bytes overhead per entry. For millions of simple int→int mappings, use specialized libraries like Eclipse Collections `IntIntHashMap` or arrays
8. **When data is too large for heap** → Use off-heap stores or disk-based solutions (MapDB, RocksDB)

**Story:** We had a service storing 50M entries in a `HashMap<Long, Long>`. Each `Entry` object + boxed `Long` keys/values consumed ~100 bytes. Total: ~5GB just for the map. Switching to Eclipse Collections' `LongLongHashMap` (primitive-specialized) cut it to ~800MB.

---

### Q: How to handle large data sets in memory?

**Answer:**

**Strategies I've applied in production:**

1. **Streaming/pagination** — Don't load everything. Use `Stream`, database cursors, or paginated queries:
```java
try (Stream<Order> orders = orderRepo.streamAll()) {
    orders.filter(Order::isActive)
          .forEach(this::process);
}
```

2. **Primitive-specialized collections** — Avoid autoboxing overhead:
```java
// Instead of HashMap<Integer, Integer> (100 bytes/entry)
IntIntHashMap map = new IntIntHashMap(); // ~16 bytes/entry
```

3. **Off-heap storage** — For very large data that shouldn't pressure GC:
   - Chronicle Map — off-heap concurrent map
   - MapDB — disk/off-heap backed collections

4. **Bloom filters** — For set-membership checks on huge datasets:
```java
BloomFilter<String> filter = BloomFilter.create(Funnels.stringFunnel(), 10_000_000, 0.01);
filter.put("item");
filter.mightContain("item"); // fast probabilistic check
```

5. **Partitioning/sharding** — Split data across multiple maps or nodes

6. **Compact representations** — Use `byte[]` or `ByteBuffer` instead of rich objects when possible

---

### Q: How do collections impact GC and memory usage?

**Answer:**

**Memory overhead per collection entry:**

| Collection | Overhead per Entry (approx.) |
|-----------|------------------------------|
| `ArrayList` | 4 bytes (object reference) + amortized spare capacity |
| `LinkedList` | ~48 bytes (Node: item ref + prev + next + header) |
| `HashMap` | ~48-64 bytes (Node: key ref + value ref + hash + next + header) |
| `TreeMap` | ~64-80 bytes (TreeNode: key + value + left + right + parent + color + header) |
| `HashSet` | Same as HashMap (it wraps HashMap internally) |

**GC impact:**

1. **Large arrays** (ArrayList, HashMap's table): Allocated in old-gen if > ~8KB. Resizing creates garbage in old-gen → full GC pauses.

2. **Linked structures** (LinkedList, TreeMap): Each node is a separate heap object → more objects for GC to track, longer GC scan times, poor locality.

3. **Short-lived collections:** Creating and discarding collections in hot loops generates enormous young-gen garbage → frequent minor GCs.

4. **Finalizers/cleaners on cached objects:** If cached objects have finalizers, they survive an extra GC cycle → memory pressure.

**Production optimizations I've applied:**
- Pre-size collections to avoid resize garbage
- Reuse collections via `clear()` instead of creating new ones in loops
- Use `ArrayList` over `LinkedList` — fewer objects, less GC scan work
- Monitor collection sizes in production (expose as metrics) — a growing map often indicates a memory leak
- Use WeakHashMap or Caffeine with size bounds for caches — unbounded caches are the #1 cause of OOM in Java services I've debugged

---

## ⚠️ 10. Tricky / Frequently Asked Edge Cases

---

### Q: Can HashMap have null key? Null values?

**Answer:**

**Yes to both.**
- `HashMap` allows **one** null key and **multiple** null values.
- Null key always maps to bucket index 0.
- Null values are valid and retrievable.

```java
Map<String, String> map = new HashMap<>();
map.put(null, "nullKeyValue");     // allowed
map.put("key1", null);              // allowed
map.put("key2", null);              // allowed

map.get(null);    // "nullKeyValue"
map.get("key1");  // null — but is it "no mapping" or "mapped to null"?
```

**The ambiguity problem:**
```java
map.get("unknownKey"); // null — key doesn't exist
map.get("key1");       // null — key exists but value is null

// Use containsKey to disambiguate:
map.containsKey("unknownKey"); // false
map.containsKey("key1");       // true
```

**Production stance:** I avoid null keys and values in maps. They create ambiguity and are bug magnets. Use `Optional` or a sentinel value instead. If you use null values, always use `containsKey()` for existence checks, never rely on `get() != null`.

---

### Q: Why does HashMap allow only one null key?

**Answer:**

Because keys must be unique (by `equals()`), and `null.equals(null)` is conceptually true (handled specially in the code). So there can only be one null key, just like there can only be one of any other key.

**Internally:** HashMap special-cases null keys:
```java
static final int hash(Object key) {
    int h;
    return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
}
```

Null key always hashes to 0, always goes to bucket 0. When you `put(null, value)` again, it finds the existing null key in bucket 0 and replaces the value — same as any duplicate key.

---

### Q: Can TreeMap have null key?

**Answer:**

**No.** `TreeMap` throws `NullPointerException` if you try to insert a null key.

**Why:** TreeMap uses comparisons (`compareTo()` or `Comparator.compare()`) to position keys in the Red-Black Tree. You can't compare null with anything — `null.compareTo(x)` throws NPE, and most Comparators don't handle null.

```java
TreeMap<String, String> map = new TreeMap<>();
map.put(null, "value"); // NullPointerException
```

**Workaround** (if you really need it — which you shouldn't):
```java
TreeMap<String, String> map = new TreeMap<>(Comparator.nullsFirst(Comparator.naturalOrder()));
map.put(null, "value"); // works — nullsFirst handles null comparison
```

**Note:** Null values are fine in TreeMap — only null keys are problematic.

---

### Q: What happens if hashCode changes after insertion?

**Answer:**

The entry becomes **unreachable** — effectively a memory leak within the map.

```java
class MutableKey {
    String value;
    
    public int hashCode() { return value.hashCode(); }
    public boolean equals(Object o) {
        return o instanceof MutableKey && ((MutableKey) o).value.equals(this.value);
    }
}

Map<MutableKey, String> map = new HashMap<>();
MutableKey key = new MutableKey();
key.value = "original";
map.put(key, "data");         // stored in bucket based on "original".hashCode()

key.value = "mutated";        // hashCode changes!
map.get(key);                 // null! Looks in bucket for "mutated".hashCode() — wrong bucket
map.containsKey(key);         // false!
map.size();                   // 1 — entry still exists, just unreachable

// Even creating a "matching" key won't find it:
MutableKey key2 = new MutableKey();
key2.value = "original";
map.get(key2);                // null! (unless key2 happens to land in the same bucket AND
                              // the original entry's stored hash matches)
```

**The entry is orphaned:**
- It exists in the map, consuming memory
- It can't be retrieved, updated, or removed by key
- It contributes to `size()` and affects resize threshold
- Only `clear()` or iterating over `entrySet()` can reach it

**Production rule:** **Map keys must be immutable.** This is why `String`, `Integer`, and other immutable types are ideal keys. If you must use a mutable object as a key, never modify the fields that participate in `hashCode()` and `equals()` after insertion.

---

### Q: Why is String commonly used as HashMap key?

**Answer:**

`String` is the ideal HashMap key for several reasons:

1. **Immutable** — `hashCode()` can never change after insertion, preventing the orphaned-entry problem.

2. **Cached hashCode** — `String` computes its hashCode lazily and caches it:
```java
// Inside String class:
private int hash; // cached, default 0
public int hashCode() {
    int h = hash;
    if (h == 0 && !hashIsZero) {
        h = computeHashCode(); // computed once
        if (h == 0) hashIsZero = true;
        else hash = h;
    }
    return h;
}
```
Repeated `hashCode()` calls (during resize, collision resolution) are free after the first call.

3. **Well-distributed hashCode** — String's hash algorithm (31-based polynomial) provides good distribution across buckets for typical string data.

4. **Interning** — String literals and interned strings enable `==` comparison (used as a fast-path before `equals()` in HashMap's node comparison).

5. **Ubiquitous** — Almost every domain uses string identifiers (user IDs, order numbers, API keys, config keys).

**Production tip:** For maps with many identical-prefix string keys (like URLs or fully-qualified class names), consider that the hash distribution might not be ideal because `String.hashCode()` weighs earlier characters more. In extreme cases, a custom hash function can help — but I've rarely needed to go that far.

---

## 💡 Bonus (Modern Java)

---

### Q: What are List.of(), Set.of(), Map.of()?

**Answer:**

Introduced in **Java 9**, these are factory methods for creating **truly immutable** collections:

```java
// Immutable List
List<String> list = List.of("a", "b", "c");

// Immutable Set
Set<String> set = Set.of("x", "y", "z");

// Immutable Map (up to 10 entries)
Map<String, Integer> map = Map.of("a", 1, "b", 2, "c", 3);

// Immutable Map (any number of entries)
Map<String, Integer> map = Map.ofEntries(
    Map.entry("a", 1),
    Map.entry("b", 2),
    Map.entry("c", 3)
);
```

**Properties:**
- **Immutable** — `add()`, `remove()`, `set()`, `put()` → `UnsupportedOperationException`
- **No nulls** — null elements/keys/values → `NullPointerException`
- **No duplicates** — duplicate elements in `Set.of()` or duplicate keys in `Map.of()` → `IllegalArgumentException`
- **Iteration order:** `List.of()` preserves order. `Set.of()` and `Map.of()` have **unspecified** iteration order (may vary across JVM runs)
- **Compact** — internally uses specialized implementations (field-based for small sizes, array-based for larger)
- **Serializable** — safe for RMI and other serialization use cases

**Java 10 additions:**
```java
List<String> copy = List.copyOf(mutableList);   // immutable copy
Set<String> copy = Set.copyOf(mutableSet);
Map<K, V> copy = Map.copyOf(mutableMap);

// Collectors
list.stream().collect(Collectors.toUnmodifiableList());
```

**Java 16:**
```java
List<String> list = stream.toList(); // unmodifiable, slightly more lenient than List.of()
```

---

### Q: Difference between Arrays.asList() vs List.of()?

**Answer:**

| Feature | `Arrays.asList()` | `List.of()` (Java 9+) |
|---------|-------------------|----------------------|
| **Mutability** | Fixed-size but `set()` works | Fully immutable |
| **Null elements** | Allowed | Not allowed (NPE) |
| **Backed by array** | Yes — changes to array reflect in list | No — independent copy |
| **add()/remove()** | `UnsupportedOperationException` | `UnsupportedOperationException` |
| **set(index, value)** | Works! Modifies the backing array | `UnsupportedOperationException` |
| **Serializable** | Yes | Yes |

```java
// Arrays.asList — backed by the original array
String[] arr = {"a", "b", "c"};
List<String> list1 = Arrays.asList(arr);
arr[0] = "z";
System.out.println(list1.get(0)); // "z" — change reflected!
list1.set(1, "y");                // works — modifies arr[1] too
list1.add("d");                   // UnsupportedOperationException

// List.of — truly immutable, independent
List<String> list2 = List.of("a", "b", "c");
list2.set(0, "z");   // UnsupportedOperationException
list2.add("d");       // UnsupportedOperationException
```

**Common pitfall:**
```java
// This does NOT create a list of 3 ints:
List<int[]> oops = Arrays.asList(new int[]{1, 2, 3});
// It creates a List containing ONE element: the int[] array
// Because Arrays.asList(T...) treats int[] as a single T

// For primitives, use:
List<Integer> correct = List.of(1, 2, 3);
```

**My production preference:**
- Need a mutable list from an array? → `new ArrayList<>(Arrays.asList(arr))`
- Need an immutable list? → `List.of(...)` or `List.copyOf(collection)`
- `Arrays.asList()` directly? → Only in tests or throwaway code. The half-mutable behavior is confusing.

---

### Q: Streams + Collections integration?

**Answer:**

Java 8 bridged Streams and Collections, making functional-style data processing a first-class citizen.

**Key integration points:**

```java
// 1. Collection → Stream
List<Order> orders = getOrders();
Stream<Order> stream = orders.stream();           // sequential
Stream<Order> parallel = orders.parallelStream();  // parallel

// 2. Stream → Collection
List<Order> filtered = orders.stream()
    .filter(o -> o.getAmount().compareTo(BigDecimal.valueOf(100)) > 0)
    .collect(Collectors.toList());

// 3. Stream → Map (extremely common in production)
Map<String, Order> orderById = orders.stream()
    .collect(Collectors.toMap(Order::getId, Function.identity()));

// Handle duplicates:
Map<String, Order> latest = orders.stream()
    .collect(Collectors.toMap(Order::getId, Function.identity(), 
             (existing, replacement) -> replacement)); // keep last

// 4. Stream → Grouped Map
Map<Status, List<Order>> byStatus = orders.stream()
    .collect(Collectors.groupingBy(Order::getStatus));

// 5. Stream → Partitioned
Map<Boolean, List<Order>> partitioned = orders.stream()
    .collect(Collectors.partitioningBy(o -> o.getAmount().compareTo(threshold) > 0));

// 6. Default methods added to Collection (Java 8)
orders.forEach(this::process);                         // iteration
orders.removeIf(Order::isCancelled);                   // conditional removal
orders.sort(Comparator.comparing(Order::getCreatedAt)); // in-place sort
orders.replaceAll(o -> o.withStatus(Status.ARCHIVED));  // transform all

// 7. Default methods added to Map (Java 8)
map.computeIfAbsent(key, k -> expensiveLoad(k));       // lazy init
map.merge(key, 1, Integer::sum);                        // atomic update
map.forEach((k, v) -> log.info("{} -> {}", k, v));     // iteration
map.replaceAll((k, v) -> v.toUpperCase());              // transform values
```

**Production patterns I use daily:**

```java
// Transform list of entities to DTOs
List<OrderDTO> dtos = orders.stream()
    .map(orderMapper::toDTO)
    .collect(Collectors.toList());

// Build lookup map from DB results
Map<String, User> userCache = userRepo.findAll().stream()
    .collect(Collectors.toMap(User::getId, Function.identity()));

// Aggregate statistics
DoubleSummaryStatistics stats = orders.stream()
    .mapToDouble(o -> o.getAmount().doubleValue())
    .summaryStatistics();
// stats.getAverage(), stats.getMax(), stats.getCount()

// Flatten nested collections
List<LineItem> allItems = orders.stream()
    .flatMap(o -> o.getLineItems().stream())
    .collect(Collectors.toList());
```

**Parallel streams warning:** Don't blindly use `parallelStream()`. It uses the common `ForkJoinPool` (shared across your entire application). A slow stream operation can starve other parallel streams. In production, I use parallel streams only for CPU-bound operations on large datasets (10K+ elements) and always benchmark first. For I/O-bound operations, use a dedicated `ExecutorService` instead.

---

> **Final note:** These answers reflect 6 years of building, debugging, and optimizing Java production systems. The key differentiator at senior level isn't just knowing the API — it's understanding the internals, knowing when the default choice breaks down, and having war stories about what goes wrong in production. Always think about thread safety, memory footprint, GC impact, and cache locality when choosing a collection.
