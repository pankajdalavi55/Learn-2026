# Java Collections Framework – Internals, Performance & Interview Guide

> **Audience:** Staff / Senior Java Engineers (6+ years) preparing for advanced interviews & system design.  
> **Scope:** Internals, performance analysis, concurrency, production insights, and interview-focused content.

---

## Table of Contents

1. [Overview of Java Collections Framework](#1-overview-of-java-collections-framework)
- **[Theory C: Amortized Analysis — Why ArrayList.add() is O(1)](#theory-c-amortized-analysis--why-arraylistadd-is-o1)**
2. [List Implementations – Internals & Performance](#2-list-implementations--internals--performance)
- **[Theory A: Hashing — Fundamentals & Deep Understanding](#theory-a-hashing--fundamentals--deep-understanding)**
3. [Set Implementations – Internals](#3-set-implementations--internals)
    - [Theory B: Red-Black Tree — Understanding the Data Structure](#theory-b-red-black-tree--understanding-the-data-structure) *(inside Section 3, before TreeSet)*
4. [Map Implementations – Deep Dive](#4-map-implementations--deep-dive)
- **[Theory D: Java Memory Model, Object Layout & Reference Types](#theory-d-java-memory-model-object-layout--reference-types)**
5. [Queue & Deque Implementations](#5-queue--deque-implementations)
- **[Theory E: Concurrency Primitives — Foundations for Thread-Safe Collections](#theory-e-concurrency-primitives--foundations-for-thread-safe-collections)**
6. [Concurrency & Collections](#6-concurrency--collections)
- **[Theory F: CPU Cache, Data Locality & Why It Matters for Collections](#theory-f-cpu-cache-data-locality--why-it-matters-for-collections)**
7. [Performance & Memory Considerations](#7-performance--memory-considerations)
- **[Theory G: Collection Design Patterns & Best Practices](#theory-g-collection-design-patterns--best-practices)**
8. [Advanced Topics](#8-advanced-topics)
9. [Staff-Level / Senior-Level Interview Questions](#9-staff-level--senior-level-interview-questions)
    - 9A. Deep Internal Questions
    - 9B. Performance & Debugging Scenarios
    - 9C. System Design Angle

---

## 1. Overview of Java Collections Framework

### 1.1 Evolution of Collections (Java 1.2 → Java 21+)

| Java Version | Key Collections Additions |
|---|---|
| **1.0** | `Vector`, `Hashtable`, `Stack`, `Enumeration` — synchronized, legacy |
| **1.2** | **Collections Framework born:** `ArrayList`, `LinkedList`, `HashMap`, `HashSet`, `TreeMap`, `TreeSet`, `Iterator` |
| **1.4** | `LinkedHashMap`, `LinkedHashSet`, `IdentityHashMap` |
| **1.5** | Generics, `PriorityQueue`, `ConcurrentHashMap`, `CopyOnWriteArrayList`, `CopyOnWriteArraySet`, `EnumSet`, `EnumMap`, `Queue`, `BlockingQueue` |
| **1.6** | `ArrayDeque`, `NavigableSet`, `NavigableMap`, `ConcurrentSkipListMap`, `ConcurrentSkipListSet`, `Deque` |
| **1.7** | Diamond operator `<>`, `TransferQueue`, `LinkedTransferQueue` |
| **1.8** | `Spliterator`, `Stream` integration, `default` methods in `Collection`/`Map` (`forEach`, `removeIf`, `replaceAll`, `compute`, `merge`), HashMap treeification |
| **9** | Factory methods: `List.of()`, `Set.of()`, `Map.of()`, `Map.ofEntries()` — immutable collections |
| **10** | `List.copyOf()`, `Set.copyOf()`, `Map.copyOf()`, `Collectors.toUnmodifiableList/Set/Map` |
| **16** | `Stream.toList()` (unmodifiable) |
| **21** | `SequencedCollection`, `SequencedSet`, `SequencedMap` — unified access to first/last elements |

---

### 1.2 Core Interface Hierarchy (ASCII Diagram)

```
                          ┌────────────┐
                          │  Iterable  │
                          └─────┬──────┘
                                │
                          ┌─────▼──────┐
                          │ Collection │
                          └─────┬──────┘
                  ┌─────────────┼──────────────┐
                  │             │              │
            ┌─────▼────┐ ┌─────▼─────┐  ┌─────▼─────┐
            │   List   │ │    Set    │  │   Queue   │
            └─────┬────┘ └─────┬─────┘  └─────┬─────┘
                  │            │              │
            ┌─────┘      ┌────┴─────┐   ┌────▼─────┐
            │            │          │   │   Deque   │
      Implementations  SortedSet   │   └──────────┘
      - ArrayList      │          │   - ArrayDeque
      - LinkedList  NavigableSet  │   - LinkedList
      - Vector         │          │   - PriorityQueue
      - CopyOnWrite  TreeSet      │   - BlockingQueues
        ArrayList    LinkedHashSet│
                     HashSet      │
                     EnumSet      │
                                  │
                            ┌─────┘
                            │
                      Implementations
                      - HashSet
                      - LinkedHashSet
                      - TreeSet
                      - EnumSet

       ┌──────────┐
       │   Map    │  (NOT part of Collection hierarchy)
       └────┬─────┘
     ┌──────┼────────────┬──────────────┐
     │      │            │              │
  HashMap  Hashtable  SortedMap    WeakHashMap
  │               ┌──────┘
  LinkedHashMap  NavigableMap
  ConcurrentHM   │
  EnumMap       TreeMap
  IdentityHM    ConcurrentSkipListMap

  Java 21+:
  SequencedCollection ─► SequencedSet ─► SortedSet ─► NavigableSet
  SequencedMap ─► SortedMap ─► NavigableMap
```

---

### 1.3 Core Interfaces — Responsibilities

| Interface | Purpose | Key Methods |
|---|---|---|
| **Iterable** | Enables for-each loop | `iterator()`, `forEach()`, `spliterator()` |
| **Collection** | Root of the collection hierarchy | `add`, `remove`, `contains`, `size`, `stream`, `toArray` |
| **List** | Ordered, indexed, duplicates allowed | `get(i)`, `set(i, e)`, `indexOf`, `subList`, `listIterator` |
| **Set** | No duplicates, unordered (unless sorted/linked) | `add`, `contains`, `remove` (all O(1) for Hash-based) |
| **SortedSet** | Set with total ordering | `first()`, `last()`, `headSet()`, `tailSet()`, `subSet()` |
| **NavigableSet** | SortedSet + navigation | `lower()`, `floor()`, `ceiling()`, `higher()`, `pollFirst/Last()` |
| **Queue** | FIFO (typically), head/tail operations | `offer()`, `poll()`, `peek()` (non-throwing), `add/remove/element` (throwing) |
| **Deque** | Double-ended queue | `offerFirst/Last()`, `pollFirst/Last()`, `peekFirst/Last()` |
| **Map** | Key-value pairs, unique keys | `get`, `put`, `containsKey`, `entrySet`, `compute`, `merge` |
| **SortedMap** | Map with sorted keys | `firstKey()`, `lastKey()`, `headMap()`, `tailMap()` |
| **NavigableMap** | SortedMap + navigation | `lowerEntry()`, `floorEntry()`, `ceilingEntry()`, `higherEntry()` |

---

### 1.4 Iterable & Iterator Internals

```java
// Iterator protocol — every Collection supports this
Iterator<String> it = list.iterator();
while (it.hasNext()) {
    String s = it.next();
    if (shouldRemove(s)) it.remove();  // ONLY safe way to remove during iteration
}

// ListIterator — bidirectional, index-aware (List only)
ListIterator<String> lit = list.listIterator(list.size()); // start at end
while (lit.hasPrevious()) {
    String s = lit.previous();
    lit.set(s.toUpperCase());  // Replace current element in-place
}
```

**Fail-fast contract:**
- Internal `modCount` field tracks structural modifications.
- Iterator caches `expectedModCount` at creation.
- On each `next()` / `remove()`, checks `modCount == expectedModCount`.
- If mismatch → throws `ConcurrentModificationException`.
- **Best-effort only**, not guaranteed under concurrent access.

---

### 1.5 Comparable vs Comparator

```
╔══════════════════════════╦═════════════════════════════════════════════════╗
║       Comparable         ║         Comparator                             ║
╠══════════════════════════╬═════════════════════════════════════════════════╣
║ java.lang.Comparable<T>  ║ java.util.Comparator<T>                       ║
║ Object defines its own   ║ External comparison logic                     ║
║   natural ordering       ║   (decoupled from class)                      ║
║ compareTo(T o)           ║ compare(T o1, T o2)                           ║
║ One per class            ║ Many per class (different orderings)           ║
║ Used by: TreeSet,        ║ Used by: Collections.sort(list, comp),         ║
║   TreeMap, sorted()      ║   stream().sorted(comp), TreeSet(comp)        ║
║ Consistent with equals() ║ May differ from equals()                      ║
║   — STRONGLY recommended ║   — document if inconsistent                  ║
╚══════════════════════════╩═════════════════════════════════════════════════╝
```

```java
// Comparable — natural order baked into the class
public record Employee(String name, int age) implements Comparable<Employee> {
    @Override
    public int compareTo(Employee other) {
        return Integer.compare(this.age, other.age);  // natural order by age
    }
}

// Comparator — external, composable, reusable
Comparator<Employee> byName = Comparator.comparing(Employee::name);
Comparator<Employee> byAgeThenName = Comparator.comparingInt(Employee::age)
    .thenComparing(Employee::name);
Comparator<Employee> bySalaryDesc = Comparator.comparing(Employee::salary).reversed();
Comparator<Employee> nullSafe = Comparator.comparing(Employee::name, 
    Comparator.nullsLast(Comparator.naturalOrder()));
```

**Interview pitfall:** If `compareTo()` is inconsistent with `equals()`, `TreeSet` and `TreeMap` behave incorrectly — they use `compareTo == 0` to determine equality, not `equals()`.

---

### 1.6 Legacy Classes — Why Discouraged

| Legacy Class | Modern Replacement | Why Discouraged |
|---|---|---|
| `Vector` | `ArrayList` (+ `Collections.synchronizedList` if needed) | Synchronizes **every** method call → severe contention under load |
| `Hashtable` | `HashMap` or `ConcurrentHashMap` | Global lock on every operation; `null` keys/values prohibited unnecessarily |
| `Stack` | `ArrayDeque` | Extends `Vector` (inheritance tax), allows random access via `get(i)` which breaks stack semantics |
| `Enumeration` | `Iterator` | No `remove()`, no fail-fast, verbose API |
| `Dictionary` | `Map` | Abstract class, not interface; obsolete since Java 1.2 |

```java
// ❌ Don't
Vector<String> v = new Vector<>();
Hashtable<String, Integer> ht = new Hashtable<>();
Stack<Integer> stack = new Stack<>();

// ✅ Do
List<String> list = new ArrayList<>();                              // single-threaded
List<String> syncList = Collections.synchronizedList(new ArrayList<>());  // if sync needed
Map<String, Integer> map = new ConcurrentHashMap<>();               // concurrent access
Deque<Integer> stack = new ArrayDeque<>();                          // LIFO stack
```

**Production insight:** `Vector` and `Hashtable` still appear in legacy codebases and older library APIs (JNDI, JDBC metadata). When interfacing with them, convert immediately:
```java
List<String> modern = new ArrayList<>(legacyVector);
Map<K,V> modern = new HashMap<>(legacyHashtable);
```

---

### 1.7 Summary Comparison Table — All Core Implementations

```
╔══════════════════════╦═══════╦════════╦═════════╦══════════╦════════════════════════════╗
║ Implementation        ║ Order  ║ Nulls  ║ Thread  ║ Dups     ║ Backed By                 ║
║                       ║        ║        ║ Safe    ║          ║                           ║
╠══════════════════════╬═══════╬════════╬═════════╬══════════╬════════════════════════════╣
║ ArrayList            ║ Index  ║ Yes    ║ No      ║ Yes      ║ Object[]                  ║
║ LinkedList           ║ Index  ║ Yes    ║ No      ║ Yes      ║ Doubly-linked nodes       ║
║ CopyOnWriteArrayList ║ Index  ║ Yes    ║ Yes     ║ Yes      ║ Copy-on-write Object[]    ║
╠══════════════════════╬═══════╬════════╬═════════╬══════════╬════════════════════════════╣
║ HashSet              ║ None   ║ 1 null ║ No      ║ No       ║ HashMap                   ║
║ LinkedHashSet        ║ Insert ║ 1 null ║ No      ║ No       ║ LinkedHashMap              ║
║ TreeSet              ║ Sorted ║ No*    ║ No      ║ No       ║ TreeMap (Red-Black Tree)   ║
║ EnumSet              ║ Enum   ║ No     ║ No      ║ No       ║ Bit vector                ║
║ CopyOnWriteArraySet  ║ Insert ║ Yes    ║ Yes     ║ No       ║ CopyOnWriteArrayList      ║
╠══════════════════════╬═══════╬════════╬═════════╬══════════╬════════════════════════════╣
║ HashMap              ║ None   ║ 1 null ║ No      ║ No (key) ║ Array + linked/tree bins  ║
║                       ║        ║ key    ║         ║          ║                           ║
║ LinkedHashMap        ║ Insert ║ 1 null ║ No      ║ No (key) ║ HashMap + doubly-linked   ║
║                       ║ /Access║ key   ║         ║          ║                           ║
║ TreeMap              ║ Sorted ║ No*    ║ No      ║ No (key) ║ Red-Black Tree            ║
║ ConcurrentHashMap    ║ None   ║ No     ║ Yes     ║ No (key) ║ CAS + synchronized bins   ║
║ WeakHashMap          ║ None   ║ 1 null ║ No      ║ No (key) ║ Weak ref keys + ReferenceQ║
║ EnumMap              ║ Enum   ║ No key ║ No      ║ No (key) ║ Object[] indexed by ordinal║
║ IdentityHashMap      ║ None   ║ Yes    ║ No      ║ No (key) ║ Linear probe, == identity ║
╠══════════════════════╬═══════╬════════╬═════════╬══════════╬════════════════════════════╣
║ PriorityQueue        ║ Heap   ║ No     ║ No      ║ Yes      ║ Object[] (binary heap)    ║
║ ArrayDeque           ║ FIFO/  ║ No     ║ No      ║ Yes      ║ Circular Object[]         ║
║                       ║ LIFO   ║        ║         ║          ║                           ║
║ ArrayBlockingQueue   ║ FIFO   ║ No     ║ Yes     ║ Yes      ║ Circular Object[] + locks ║
║ LinkedBlockingQueue  ║ FIFO   ║ No     ║ Yes     ║ Yes      ║ Linked nodes + 2 locks    ║
║ ConcurrentLinkedQueue║ FIFO   ║ No     ║ Yes     ║ Yes      ║ Lock-free linked nodes    ║
╚══════════════════════╩═══════╩════════╩═════════╩══════════╩════════════════════════════╝

* TreeSet/TreeMap: null keys throw NullPointerException (Comparator decides)
```

---

## Theory C: Amortized Analysis — Why ArrayList.add() is O(1)

> Understanding amortized analysis is essential for correctly reasoning about ArrayList, HashMap resize, and StringBuilder performance.

### What is Amortized Analysis?

Amortized analysis computes the **average cost per operation over a sequence of operations**, even when individual operations vary wildly in cost. It's not the same as average-case analysis (which relies on probability/input distribution) — amortized analysis provides a **guaranteed worst-case bound for the total**.

```
Simple Example: ArrayList.add()

  Most adds: O(1) — just store reference at index
  Occasional add: O(N) — resize array (copy all N elements)

  But resize doubles capacity, so the next N adds are O(1) again.
  Total cost of N adds = N (direct) + (1 + 2 + 4 + 8 + ... + N) (copies)
                       = N + ~2N = 3N
  Amortized cost per add = 3N / N = O(1)  ✓
```

### Three Methods of Amortized Analysis

#### 1. Aggregate Method

Compute total cost of N operations, divide by N.

```
ArrayList with growth factor 1.5× (Java's actual factor):

  Capacity sequence: 10 → 15 → 22 → 33 → 49 → 73 → ...
  
  For N = 100 inserts:
    Copies at resize: 10 + 15 + 22 + 33 = 80 copies total
    Direct inserts: 100
    Total work: 180
    Amortized per insert: 180/100 = 1.8 = O(1)

  Growth factor matters:
    Factor 2×:   copies sum to ~2N  → amortized = ~3
    Factor 1.5×: copies sum to ~3N  → amortized = ~4  (slightly worse)
    Factor 1.1×: copies sum to ~10N → amortized = ~11 (poor)
    Factor 1.0×: copies = N per resize → amortized = O(N) (linear growth = no amortization)
```

#### 2. Accounting Method (Banker's Method)

Assign each operation a fixed "amortized cost" that may overcharge cheap operations, building up "credit" to pay for expensive ones.

```
ArrayList.add() with 2× growth:
  Charge each add $3 (amortized cost):
    $1: Store the element
    $1: Save to pay for copying THIS element during next resize
    $1: Save to pay for copying an OLDER element during next resize

  When resize happens:
    Each of the N/2 new elements since last resize has saved $2
    Total savings: N/2 × $2 = $N → exactly enough to copy all N elements!

  Conclusion: $3 per operation covers everything → O(1) amortized ✓
```

#### 3. Potential Method (Physicist's Method)

Define a potential function Φ that measures "stored energy" in the data structure.

```
For ArrayList:
  Φ = 2 × size - capacity

  Normal add: actual cost = 1, ΔΦ = +2
    Amortized = 1 + 2 = 3

  Resize add: actual cost = N+1, capacity doubles from N to 2N
    ΔΦ = (2(N+1) - 2N) - (2N - N) = (2) - (N) = 2 - N
    Amortized = (N+1) + (2 - N) = 3

  Same answer: O(1) amortized ✓
```

### Amortized Analysis in Java Collections

| Operation | Worst Case | Amortized | Explanation |
|---|---|---|---|
| `ArrayList.add()` | O(N) | **O(1)** | Resize doubles capacity |
| `HashMap.put()` | O(N) | **O(1)** | Rehash at threshold, doubles capacity |
| `StringBuilder.append()` | O(N) | **O(1)** | Same resize strategy as ArrayList |
| `ArrayDeque.add()` | O(N) | **O(1)** | Circular array doubles on full |
| `PriorityQueue.add()` | O(log N) | **O(log N)** | Heap sift-up (no resize amortization helps complexity class) |
| `Stack.push()` (Vector) | O(N) | **O(1)** | Same as ArrayList internally |

### Common Interview Misconception

```
❌ "ArrayList.add() is O(1) on average"
   → Wrong terminology. "Average" implies probability distribution.
  
✅ "ArrayList.add() is O(1) amortized"
   → Correct. It's a deterministic guarantee over any sequence.

❌ "ArrayList.add() is always O(1)"
   → Wrong. Individual operations can be O(N) during resize.

✅ "ArrayList.add() is O(1) amortized, O(N) worst-case per single operation"
   → Precise and complete answer.
```

---

## 2. List Implementations – Internals & Performance

### 2.1 ArrayList — The Default Workhorse

#### Internal Structure

```
ArrayList<E>:
  ┌─────────────────────────────────────────────────────┐
  │  Object[] elementData  (the backing array)          │
  │  int size              (logical size, NOT capacity)  │
  ├─────────────────────────────────────────────────────┤
  │  elementData:                                       │
  │  ┌────┬────┬────┬────┬────┬────┬────┬────┬────┬───┐ │
  │  │ e0 │ e1 │ e2 │ e3 │ e4 │null│null│null│null│...│ │
  │  └────┴────┴────┴────┴────┴────┴────┴────┴────┴───┘ │
  │  ◄────── size=5 ──────►◄──── unused capacity ────► │
  │  ◄────────────── elementData.length = 10 ────────► │
  └─────────────────────────────────────────────────────┘
```

#### Internal Array Resizing Mechanism

```java
// Simplified from OpenJDK source — ArrayList.grow():
private Object[] grow(int minCapacity) {
    int oldCapacity = elementData.length;
    if (oldCapacity > 0 || elementData != DEFAULTCAPACITY_EMPTY_ELEMENTDATA) {
        int newCapacity = ArraysSupport.newLength(
            oldCapacity,
            minCapacity - oldCapacity,   // minimum growth
            oldCapacity >> 1             // preferred growth = 50%
        );
        return elementData = Arrays.copyOf(elementData, newCapacity);
    } else {
        return elementData = new Object[Math.max(DEFAULT_CAPACITY, minCapacity)];
    }
}
```

**Key internals:**

| Property | Value |
|---|---|
| Default initial capacity | **10** (allocated on first `add()`, not on construction) |
| Growth factor | **1.5×** (oldCapacity + oldCapacity >> 1) |
| Empty list memory | Just the reference + empty `Object[]` sentinel |
| Max array size | `Integer.MAX_VALUE - 8` (some VMs reserve header words) |
| Copy mechanism | `System.arraycopy()` → native memcpy, very fast |

**Amortized O(1) `add()` analysis:**
- N insertions trigger ~log₁.₅(N) resize operations.
- Each resize copies all current elements: costs 10 + 15 + 22 + 33 + ...
- Total copy cost for N inserts ≈ 3N → amortized **O(1)** per `add()`.
- **Worst single insert:** O(N) when resize happens.

#### Memory Footprint

```
ArrayList object:       16 bytes (header) + 4 (size) + 4 (modCount) + 8 (ref to array) = ~32 bytes
Backing array:          16 bytes (header) + 4 (length) + capacity × 8 (references) + padding
Per-element overhead:   8 bytes (reference) + boxed object overhead if primitive

Example: ArrayList<Integer> with 1000 elements:
  - ArrayList object:     ~32 bytes
  - Object[] (cap 1500):  16 + 1500 × 8 = ~12,016 bytes
  - 1000 Integer objects:  1000 × 16 = 16,000 bytes
  Total: ~28 KB for 1000 ints (vs. 4 KB for int[1000])
  Overhead ratio: 7× vs primitive array
```

#### Fail-Fast Iterator Internals

```java
// Inside ArrayList$Itr:
int cursor;               // index of next element to return
int lastRet = -1;         // index of last element returned
int expectedModCount = modCount;  // snapshot at creation

public E next() {
    checkForComodification();  // if (modCount != expectedModCount) throw CME
    int i = cursor;
    Object[] elementData = ArrayList.this.elementData;
    cursor = i + 1;
    return (E) elementData[lastRet = i];
}
```

---

### 2.2 LinkedList — Doubly-Linked Nodes

#### Internal Structure

```
LinkedList<E>:
  ┌──────────────────────────────────────────────────────────────┐
  │  int size                                                     │
  │  Node<E> first  ──────────┐                                  │
  │  Node<E> last   ──────────┼──┐                               │
  ├────────────────────────────┼──┼──────────────────────────────┤
  │                            │  │                               │
  │  null ◄── prev             ▼  │                               │
  │        ┌──────┐  next  ┌──────┐  next  ┌──────┐  next        │
  │        │ Node │ ─────► │ Node │ ─────► │ Node │ ─────► null  │
  │   null │ item │ ◄───── │ item │ ◄───── │ item │              │
  │        │  =A  │  prev  │  =B  │  prev  │  =C  │◄─────┘      │
  │        └──────┘        └──────┘        └──────┘              │
  │          ▲ first                          ▲ last              │
  └──────────────────────────────────────────────────────────────┘
```

#### Node Structure

```java
// OpenJDK: LinkedList.Node
private static class Node<E> {
    E item;
    Node<E> next;
    Node<E> prev;
    Node(Node<E> prev, E element, Node<E> next) {
        this.item = element;
        this.next = next;
        this.prev = prev;
    }
}
// Memory per node: 16 (object header) + 8 (item ref) + 8 (next ref) + 8 (prev ref) = 40 bytes
// vs. ArrayList: 8 bytes per reference slot
// 5× more memory per element than ArrayList
```

#### Why LinkedList Is Rarely Optimal In Practice

```
"LinkedList: O(1) insert/delete!" — True in theory. False in practice.

Problem 1: O(N) traversal to FIND the position
  - list.get(500000) → walks 500K nodes (no random access)
  - list.add(index, element) → traversal to index = O(N), then insert = O(1)
  - Net: O(N) for indexed insert, not O(1)

Problem 2: Cache locality disaster
  - ArrayList: elements in contiguous memory → CPU cache prefetch works
  - LinkedList: nodes scattered across heap → every node access = cache miss
  - Real-world benchmark: ArrayList iteration is 10-100× faster on modern CPUs

Problem 3: Memory overhead
  - Each node: 40 bytes overhead (prev + next + header)
  - 1M elements: LinkedList = ~40 MB overhead vs ArrayList ≈ ~8 MB (array padding)
  - GC must trace every node individually → longer GC pauses

Problem 4: GC pressure
  - Adding 1M elements = 1M node objects allocated
  - ArrayList: ~20 array copies (log1.5(1M) ≈ 35 resizes)
  - Old arrays immediately eligible for GC, simple

When LinkedList IS appropriate:
  ✓ Frequent add/removeFirst/Last (stack/queue) — O(1) without copy
  ✓ Iterator-based removal during traversal
  ✓ Implementing LRU with manual node manipulation
  ✓ Almost never in production code — ArrayDeque is usually better
```

---

### 2.3 CopyOnWriteArrayList — Thread-Safe Reads

#### Internal Mechanism

```
CopyOnWriteArrayList<E>:
  ┌──────────────────────────────────────────────────────────────┐
  │  volatile Object[] array;   // reads see latest snapshot     │
  │  final ReentrantLock lock;  // writes are serialized         │
  └──────────────────────────────────────────────────────────────┘

Write operation (add/set/remove):
  1. Acquire lock
  2. Copy entire current array to new array
  3. Modify new array
  4. Volatile-write new array reference (publishes to readers)
  5. Release lock

Read operation (get/iterator/size):
  1. Read volatile array reference (one volatile read)
  2. Access element directly — NO locking, NO CAS
  3. Iterator snapshots array at creation time — never throws CME
```

```java
// Simplified add() from OpenJDK:
public boolean add(E e) {
    synchronized (lock) {
        Object[] es = getArray();
        int len = es.length;
        es = Arrays.copyOf(es, len + 1);  // full copy!
        es[len] = e;
        setArray(es);                       // volatile write
        return true;
    }
}

// get() — no locking at all:
public E get(int index) {
    return elementAt(getArray(), index);  // just a volatile read + array access
}
```

#### Memory Trade-offs

```
Operation cost:
  add():    O(N) — copies entire array every time
  set():    O(N) — copies entire array
  get():    O(1) — volatile read of array ref, then direct index
  size():   O(1)
  iterator: O(1) to create, snapshot semantics (stale reads OK)

Memory during write:
  - Original array:      N × 8 bytes (alive until GC)
  - New array:           (N+1) × 8 bytes
  - Peak: 2× array memory during each write

Ideal for:
  ✓ Read-heavy, write-rare (listener lists, configuration, routing tables)
  ✓ Small collections (< 1000 elements)
  ✓ Scenarios where snapshot iteration is acceptable

NOT suitable for:
  ✗ Frequent writes (O(N) copy per write = catastrophic)
  ✗ Large collections (array copy cost + memory)
  ✗ Write-heavy concurrent workloads
```

---

### 2.4 Big-O Complexity Comparison — List Implementations

```
╔═══════════════════════╦════════════════╦════════════════╦══════════════════════╗
║ Operation              ║ ArrayList      ║ LinkedList     ║ CopyOnWriteArrayList ║
╠═══════════════════════╬════════════════╬════════════════╬══════════════════════╣
║ get(index)            ║ O(1)           ║ O(N)           ║ O(1)                 ║
║ set(index, e)         ║ O(1)           ║ O(N)           ║ O(N) (copy)          ║
║ add(e) — end          ║ O(1) amortized ║ O(1)           ║ O(N) (copy)          ║
║ add(index, e) — mid   ║ O(N)           ║ O(N) (traverse)║ O(N) (copy)          ║
║ remove(index)         ║ O(N)           ║ O(N) (traverse)║ O(N) (copy)          ║
║ remove(object)        ║ O(N)           ║ O(N)           ║ O(N) (copy)          ║
║ contains(e)           ║ O(N)           ║ O(N)           ║ O(N)                 ║
║ iterator.next()       ║ O(1)           ║ O(1)           ║ O(1) (snapshot)      ║
║ iterator.remove()     ║ O(N)           ║ O(1)           ║ Not supported        ║
║ addFirst/addLast      ║ O(N) / O(1)*   ║ O(1) / O(1)   ║ O(N) (copy)          ║
║ size()                ║ O(1)           ║ O(1)           ║ O(1)                 ║
╠═══════════════════════╬════════════════╬════════════════╬══════════════════════╣
║ Thread-safe           ║ No             ║ No             ║ Yes                  ║
║ Iterator type         ║ Fail-fast      ║ Fail-fast      ║ Snapshot (fail-safe) ║
║ Memory per element    ║ ~8 bytes       ║ ~40 bytes      ║ ~8 bytes             ║
║ Cache-friendly        ║ Yes            ║ No             ║ Yes                  ║
╚═══════════════════════╩════════════════╩════════════════╩══════════════════════╝

* ArrayList.addFirst: O(N) due to array shift; addLast = O(1) amortized
```

---

### 2.5 When To Use / When NOT To Use

```
USE ArrayList WHEN:
  ✓ Default choice for almost all list use cases
  ✓ Random access by index is needed
  ✓ Iteration is the primary operation (cache-friendly)
  ✓ Additions are mostly at the end
  ✓ Known maximum size → pre-size with new ArrayList<>(capacity)

USE LinkedList WHEN:
  ✓ Primarily used as Deque (addFirst/addLast, removeFirst/removeLast)
  ✓ Frequent removal during iteration (via iterator.remove())
  ✓ Almost never — prefer ArrayDeque for stack/queue behavior

USE CopyOnWriteArrayList WHEN:
  ✓ Read/iterate frequency >> write frequency (100:1 ratio or more)
  ✓ Collection is small (< 1000 elements)
  ✓ Listener lists, observer pattern, configuration lists
  ✓ Need snapshot iteration without external synchronization

PRODUCTION PITFALL — ArrayList pre-sizing:
```

```java
// ❌ Causes ~25 resize operations for 1M elements:
List<String> list = new ArrayList<>();
for (int i = 0; i < 1_000_000; i++) list.add(data[i]);

// ✅ Zero resizes, zero wasted copies:
List<String> list = new ArrayList<>(1_000_000);
for (int i = 0; i < 1_000_000; i++) list.add(data[i]);

// ✅ Even better — if size is known, use streams:
List<String> list = Arrays.stream(data).collect(Collectors.toCollection(
    () -> new ArrayList<>(data.length)));

// ✅ Best for read-only: immutable
List<String> list = List.of(data);  // or Arrays.asList(data)

// PRODUCTION: trimToSize() after bulk loading to reclaim excess capacity
list.trimToSize();  // shrinks elementData to match size
```

---

## Theory A: Hashing — Fundamentals & Deep Understanding

> This theory section provides the foundational concepts needed to understand HashSet, HashMap, ConcurrentHashMap, and related hashing-based collections.

### What is Hashing?

Hashing is the process of converting an input (key) into a fixed-size integer (hash code) that serves as an index into a data structure. The goal is **O(1) average-case** access to data.

```
Key "hello"  →  hashCode()  →  142948  →  bucket index  →  stored at position
                  ↑                          ↑
           deterministic               modular arithmetic
           (same input →               hash % capacity
            same output)               or hash & (cap-1)
```

### Hash Function Properties

| Property | Meaning | Why It Matters |
|---|---|---|
| **Deterministic** | Same input → same output, always | Core correctness requirement |
| **Uniform distribution** | Output spread evenly across range | Minimizes collisions |
| **Avalanche effect** | Small input change → large output change | Prevents clustering |
| **Fast to compute** | O(1) computation | Don't trade access speed for hash quality |

### How Java's hashCode() Works

```java
// Object.hashCode() — default
// Returns identity hash (derived from memory address, or random, JVM-dependent)
Object obj = new Object();
obj.hashCode();  // e.g., 366712642 — not the actual memory address in modern JVMs

// String.hashCode() — well-known formula
// s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1]
// Why 31? It's an odd prime. Compiler optimizes 31*i as (i << 5) - i.
"hello".hashCode();  // 99162322

// Integer.hashCode() — returns the value itself
Integer.valueOf(42).hashCode();  // 42
```

### Hash Collision: The Core Problem

A **collision** occurs when two different keys produce the same bucket index:

```
Key A: hashCode() = 100, bucket = 100 & 15 = 4
Key B: hashCode() = 116, bucket = 116 & 15 = 4   ← COLLISION! Same bucket

  Bucket[4]:
    ┌────────┐    ┌────────┐
    │ Key A  │───►│ Key B  │───► null    (linked list chaining)
    └────────┘    └────────┘

  After 8+ collisions (Java 8+):
    ┌────────┐
    │ Key A  │
    └───┬────┘
       ╱ ╲
    Key B  Key C    (Red-Black Tree — O(log N) instead of O(N))
```

### Collision Resolution Strategies

| Strategy | How It Works | Used In |
|---|---|---|
| **Separate Chaining** | Each bucket holds a linked list (or tree) of entries | Java `HashMap`, `HashSet` |
| **Open Addressing (Linear Probing)** | On collision, check next slot: `(hash + 1) % cap` | `IdentityHashMap`, some 3rd-party libs |
| **Open Addressing (Quadratic Probing)** | Check `(hash + 1²)`, `(hash + 2²)`, ... | Some embedded systems |
| **Double Hashing** | Use second hash function for step size | Rarely in Java |
| **Robin Hood Hashing** | Steal from rich (low probe distance) for poor (high probe distance) | Rust's HashMap, some C++ impls |

Java's `HashMap` uses **separate chaining** with a twist: chains longer than 8 entries become **Red-Black Trees** (treeification) for O(log N) worst case instead of O(N).

### Load Factor — The Space/Time Trade-off

```
Load Factor (λ) = number of entries / number of buckets

  λ = 0.0  →  empty table       →  wasted memory, fast lookups
  λ = 0.75 →  Java default       →  good balance
  λ = 1.0  →  full table         →  high collision probability
  λ > 1.0  →  guaranteed collisions (with chaining, this is okay)

  Why 0.75?
  ┌──────────────────────────────────────────────┐
  │  At λ = 0.75, Poisson distribution predicts: │
  │  P(bucket has 0 entries) = 0.4724             │
  │  P(bucket has 1 entry)   = 0.3543             │
  │  P(bucket has 2 entries) = 0.1329             │
  │  P(bucket has 8+ entries)= 0.00000006         │
  │                                                │
  │  → 8+ collisions virtually never happen       │
  │  → treeification is extremely rare            │
  │  → 25% wasted space is acceptable trade-off   │
  └──────────────────────────────────────────────┘
```

### Java HashMap's Hash Spreading

Java doesn't use `hashCode()` directly. It applies a **secondary hash** to spread bits:

```java
// HashMap.hash() — actual Java source
static final int hash(Object key) {
    int h;
    return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
}

// Why XOR the top 16 bits into the bottom 16?
// Bucket index = hash & (capacity - 1)
// If capacity is small (e.g., 16), only BOTTOM 4 bits matter
// Without spreading: keys differing only in top bits → same bucket!

// Example:
//   hashCode = 0xABCD_0010  →  bucket = 0x0010 & 0xF = 0
//   hashCode = 0xFFFF_0010  →  bucket = 0x0010 & 0xF = 0  ← COLLISION!
// After spreading:
//   hash = 0xABCD_0010 ^ 0x0000_ABCD = 0xABCD_ABDD  →  bucket = 0xD & 0xF = 13
//   hash = 0xFFFF_0010 ^ 0x0000_FFFF = 0xFFFF_FFEF  →  bucket = 0xF & 0xF = 15  ← NO COLLISION!
```

### Rehashing — When & How

```
Threshold = capacity × loadFactor
  ↓
When size > threshold → resize (double the capacity)

  Step 1: Allocate new array of 2× size
  Step 2: For each entry, recalculate bucket:
          newIndex = hash & (newCapacity - 1)
  Step 3: Entry either stays at index `i` or moves to `i + oldCapacity`
  Step 4: Old array becomes garbage → GC collects

  Cost: O(N) per resize, but amortized over N inserts → O(1) per insert

  PRODUCTION TIP: If you know the expected size, pre-size the map!
    new HashMap<>(expectedSize / 0.75 + 1)  →  zero resizes
    Or in Java 19+: HashMap.newHashMap(expectedSize)  →  does the math for you
```

### hashCode() and equals() Contract

```
The Contract (JSR / java.lang.Object JavaDoc):
  ╔═══════════════════════════════════════════════════════════════════════╗
  ║ 1. If a.equals(b) is true  → a.hashCode() == b.hashCode() MUST hold ║
  ║ 2. If a.hashCode() != b.hashCode() → a.equals(b) MUST be false      ║
  ║ 3. If a.hashCode() == b.hashCode() → a.equals(b) MAY be true or not ║
  ║    (hash collisions are allowed)                                      ║
  ║ 4. hashCode() must return the same value within a single execution   ║
  ║    (consistency — unless fields used in equals() change)             ║
  ╚═══════════════════════════════════════════════════════════════════════╝

Breaking the contract:

  Case 1: Override equals() but NOT hashCode()
    a.equals(b) → true, but hashCode differs
    → HashMap puts them in DIFFERENT buckets → can't find the key!

  Case 2: Use mutable fields in hashCode()
    key.setName("new") → hashCode changes
    → HashMap still looks in the OLD bucket → key is "lost"

  Case 3: Not symmetric/transitive in equals()
    → Set operations (contains, remove) become unpredictable
```

```java
// Correct implementation (using Java 7+ Objects utility):
@Override
public int hashCode() {
    return Objects.hash(firstName, lastName, age);  // consistent with equals
}

@Override
public boolean equals(Object o) {
    if (this == o) return true;
    if (o == null || getClass() != o.getClass()) return false;
    Person p = (Person) o;
    return age == p.age 
        && Objects.equals(firstName, p.firstName)
        && Objects.equals(lastName, p.lastName);
}

// Modern best practice: Use Java 16+ records
record Person(String firstName, String lastName, int age) {}
// Records auto-generate equals(), hashCode(), toString() correctly!
```

---

## 3. Set Implementations – Internals

### 3.1 HashSet — Backed by HashMap

#### Internal Structure

```
HashSet<E>:
  ┌──────────────────────────────────────────────────┐
  │  private transient HashMap<E, Object> map;       │
  │  private static final Object PRESENT = new Object(); // dummy value
  └──────────────────────────────────────────────────┘

  HashSet.add(e) → map.put(e, PRESENT);     // key = element, value = dummy
  HashSet.contains(e) → map.containsKey(e);
  HashSet.remove(e) → map.remove(e);

  All complexity is delegated to HashMap internals (see Section 4).
```

**Memory overhead:** Each `HashSet` element occupies a full `HashMap.Node`:
- 32 bytes per node (header + hash + key ref + value ref + next ref)
- Plus the `PRESENT` object: 16 bytes (shared singleton, negligible)
- Plus bucket array: capacity × 8 bytes

#### Hashing Flow — How `add(e)` Works Internally

```
HashSet.add("Hello"):
  │
  ├─ 1. Compute hash:
  │     int h = "Hello".hashCode();        → 69609650
  │     hash = h ^ (h >>> 16);             → spread bits (see Section 4.1)
  │                                          69609650 ^ 1062 = 69610200
  │
  ├─ 2. Find bucket:
  │     int index = hash & (capacity - 1); → 69610200 & 15 = 8 (for cap=16)
  │
  ├─ 3. Check bucket[8]:
  │     ├─ Empty? → Create new Node(hash, "Hello", PRESENT, null)
  │     │            Return true (added)
  │     │
  │     └─ Occupied? → Walk chain:
  │           For each node in chain:
  │             if (node.hash == hash && node.key.equals("Hello"))
  │               → Already exists, return false (not added)
  │           If no match → append new node to chain
  │           If chain length >= TREEIFY_THRESHOLD (8) → treeify bin
  │
  └─ 4. Check load:
        if (++size > threshold)  // threshold = capacity × loadFactor
          → resize() (double capacity, rehash all entries)
```

#### Load Factor & Treeification (Java 8+)

```
╔═══════════════════════╦══════════════════════════════════════════════════╗
║ Parameter              ║ Details                                          ║
╠═══════════════════════╬══════════════════════════════════════════════════╣
║ Default capacity      ║ 16 (always power of 2)                          ║
║ Default load factor   ║ 0.75 (good balance of space vs. time)           ║
║ Threshold             ║ capacity × loadFactor = 12 (initial resize at)  ║
║ Resize factor         ║ 2× (doubles capacity)                           ║
║ TREEIFY_THRESHOLD     ║ 8 — bin converts to Red-Black Tree when chain   ║
║                        ║   length >= 8 AND capacity >= 64                ║
║ UNTREEIFY_THRESHOLD   ║ 6 — tree converts back to list after removal    ║
║ MIN_TREEIFY_CAPACITY  ║ 64 — won't treeify if capacity < 64 (resizes   ║
║                        ║   instead)                                      ║
╚═══════════════════════╩══════════════════════════════════════════════════╝

Collision handling evolution:
  Java 7:  Bucket = singly-linked list always → O(N) worst case per bucket
  Java 8+: Bucket = linked list → Red-Black Tree (when >= 8 nodes)
            → O(log N) worst case per bucket

Why threshold 8 → 6 (not 8 → 8)?
  - Hysteresis prevents thrashing between list/tree
  - If threshold were same: remove one, treeify → add one, untreeify → repeat
  - Gap of 2 eliminates this oscillation

Why convert at 8?
  - Poisson distribution: P(8+ collisions) ≈ 0.00000006 with load factor 0.75
  - If it happens, likely pathological keys → tree gives O(log N) guarantee
```

---

### 3.2 LinkedHashSet — Insertion Order Maintenance

```
LinkedHashSet<E>:
  ┌─────────────────────────────────────────────────────────────────┐
  │  Extends HashSet, backed by LinkedHashMap                       │
  │                                                                 │
  │  Hash Table buckets:     Doubly-linked list (insertion order):  │
  │  ┌───┐                                                         │
  │  │ 0 │──► nodeA ──────────────────────────────────────────┐    │
  │  ├───┤                                                     │    │
  │  │ 1 │──► nodeC      head ──► nodeA ──► nodeB ──► nodeC ──►tail│
  │  ├───┤                 ◄──────── ◄──────── ◄──────── ◄─────    │
  │  │ 2 │                (doubly-linked for insertion order)       │
  │  ├───┤                                                         │
  │  │ 3 │──► nodeB                                                │
  │  └───┘                                                         │
  └─────────────────────────────────────────────────────────────────┘

  Each node has: before/after pointers (linked list) + next pointer (bucket chain)
  Memory per node: ~48 bytes (32 base + 8 before + 8 after)

  Cost of order maintenance:
  - add():    O(1) extra (append to tail of linked list)
  - remove(): O(1) extra (unlink from linked list)
  - Iteration: in insertion order, O(size) — NOT O(capacity) like HashSet
```

**When is insertion order useful?**
- Deterministic iteration for logging, debugging, testing
- Serialization/deserialization where order matters
- Maintaining user-defined ordering (config keys, column order)

---

### Theory B: Red-Black Tree — Understanding the Data Structure

> TreeSet and TreeMap are backed by Red-Black Trees. Understanding this structure is essential for senior-level interviews.

#### What is a Binary Search Tree (BST)?

A tree where each node's left children are smaller and right children are larger:

```
       8         ← root
      / \
     3   10      ← left < parent < right
    / \    \
   1   6   14
      / \  /
     4  7 13

  Lookup "7": 8→3→6→7 ✓  (O(log N) if balanced)
  
  Problem: BSTs can become UNBALANCED:
  Insert sorted data [1, 2, 3, 4, 5]:
     1
      \
       2
        \
         3         ← Degenerates to linked list!
          \        ← All operations become O(N)
           4
            \
             5
```

#### What is a Red-Black Tree?

A **self-balancing BST** that uses node coloring and rotation rules to guarantee O(log N) height:

```
Red-Black Tree Properties (ALL must be maintained):
  ╔══════════════════════════════════════════════════════════════════╗
  ║ 1. Every node is either RED or BLACK                            ║
  ║ 2. Root is always BLACK                                         ║
  ║ 3. Every leaf (NIL/null) is BLACK                               ║
  ║ 4. If a node is RED, both its children must be BLACK            ║
  ║    (no two consecutive red nodes on any path)                   ║
  ║ 5. Every path from root to any NIL leaf has the SAME number    ║
  ║    of BLACK nodes ("black-height" is uniform)                   ║
  ╚══════════════════════════════════════════════════════════════════╝

  These properties guarantee:
    Maximum height ≤ 2 × log₂(N+1)
    → All operations (insert, delete, search) are O(log N) GUARANTEED
```

#### Visual Example

```
  A valid Red-Black Tree (B=Black, R=Red):

          [13B]
         /      \
      [8R]      [17R]
      /   \      /   \
   [1B]  [11B] [15B] [25B]
     \                /
    [6R]           [22R]

  Black-height from root to any NIL = 2 (count only black nodes)
  Path: 13B → 8R → 1B → NIL  = 2 black nodes (13, 1)
  Path: 13B → 17R → 25B → 22R → NIL = 2 black nodes (13, 25)
```

#### Rotations — How Balance is Maintained

When an insertion or deletion violates Red-Black properties, **rotations** fix the tree:

```
Left Rotation (around node X):

    X                Y
   / \              / \
  A   Y    →      X   C
     / \         / \
    B   C       A   B

Right Rotation (around node Y):

      Y            X
     / \          / \
    X   C  →    A   Y
   / \             / \
  A   B           B   C

  Key insight: Rotations are O(1) — only pointer re-assignments
  At most 2 rotations per insert, 3 per delete
```

#### Insertion Algorithm (Simplified)

```
1. Insert as in normal BST
2. Color the new node RED
3. Fix violations ("fixup"):

   Case 1: Uncle is RED
     → Recolor parent & uncle to BLACK, grandparent to RED
     → Move up to grandparent and repeat

   Case 2: Uncle is BLACK, node is inner child (zig-zag)
     → Rotate node's parent in opposite direction
     → Fall through to Case 3

   Case 3: Uncle is BLACK, node is outer child (zig-zig)
     → Rotate grandparent in opposite direction
     → Recolor: parent → BLACK, grandparent → RED
     → Done!

4. Ensure root is BLACK
```

#### Red-Black Tree vs Other Balanced Trees

| Tree Type | Height Guarantee | Rotations (Insert) | Rotations (Delete) | Used In |
|---|---|---|---|---|
| **Red-Black Tree** | ≤ 2·log₂(N+1) | ≤ 2 | ≤ 3 | Java TreeMap/TreeSet, Linux kernel |
| **AVL Tree** | ≤ 1.44·log₂(N) | ≤ 2 | O(log N) | SQLite, read-heavy workloads |
| **B-Tree** | ≤ log_B(N) | N/A (split/merge) | N/A (split/merge) | Databases, filesystems |
| **Skip List** | ~log₂(N) expected | N/A (probabilistic) | N/A | Java ConcurrentSkipListMap |

**Why Java chose Red-Black Trees over AVL Trees:**
- RB trees have **fewer rotations on insert/delete** (≤3 vs O(log N) for AVL).
- Java maps/sets are **read-write balanced** workloads — AVL's stricter balancing (better reads) isn't worth the cost.
- RB trees have slightly **worse read performance** (taller by ~20%) but **significantly better write performance**.

---

### 3.3 TreeSet — Red-Black Tree Internals

#### Structure

```
TreeSet<E>:
  ┌──────────────────────────────────────────┐
  │  private transient NavigableMap<E,Object> m;   // TreeMap underneath
  └──────────────────────────────────────────┘

  TreeSet is to TreeMap as HashSet is to HashMap.
  All elements stored as keys in the backing TreeMap.
```

#### Red-Black Tree Overview

```
Red-Black Tree properties (MUST hold at all times):
  1. Every node is either RED or BLACK
  2. Root is always BLACK
  3. Every null leaf (NIL) is BLACK
  4. If a node is RED, both children must be BLACK (no two consecutive reds)
  5. Every path from root to NIL has the same number of BLACK nodes

This guarantees: height ≤ 2 × log₂(N+1)
                 All operations = O(log N) guaranteed

Example TreeSet{1, 3, 5, 7, 9, 11, 13}:

                    ┌─────────┐
                    │  7 (B)  │
                    └────┬────┘
               ┌─────────┴──────────┐
          ┌────▼────┐          ┌────▼────┐
          │  3 (R)  │          │ 11 (R)  │
          └────┬────┘          └────┬────┘
       ┌───────┴──────┐    ┌───────┴──────┐
  ┌────▼────┐  ┌──────▼┐ ┌▼──────┐ ┌─────▼───┐
  │  1 (B)  │  │ 5 (B) │ │ 9 (B) │ │ 13 (B)  │
  └─────────┘  └───────┘ └───────┘ └─────────┘
```

#### TreeSet Operations

```java
TreeSet<Integer> set = new TreeSet<>();
set.addAll(List.of(5, 10, 15, 20, 25, 30, 35));

// NavigableSet operations — range queries in O(log N):
set.headSet(20);            // [5, 10, 15] — strictly less than 20
set.tailSet(20);            // [20, 25, 30, 35] — ≥ 20
set.subSet(10, 30);         // [10, 15, 20, 25] — [10, 30)
set.subSet(10, true, 30, true); // [10, 15, 20, 25, 30] — inclusive both

// Navigation:
set.floor(18);              // 15 — greatest element ≤ 18
set.ceiling(18);            // 20 — smallest element ≥ 18
set.lower(20);              // 15 — greatest element < 20
set.higher(20);             // 25 — smallest element > 20

// With custom Comparator:
TreeSet<String> caseInsensitive = new TreeSet<>(String.CASE_INSENSITIVE_ORDER);
caseInsensitive.add("Hello");
caseInsensitive.contains("hello"); // true — uses comparator, not equals()
```

---

### 3.4 EnumSet — The Fastest Set

```java
// EnumSet is backed by a single long (up to 64 enum values) or long[]
// Operations are BITWISE — add = OR, remove = AND NOT, contains = AND
// O(1) for ALL operations. Zero boxing. Minimal memory.

enum Permission { READ, WRITE, EXECUTE, DELETE, ADMIN }

EnumSet<Permission> perms = EnumSet.of(Permission.READ, Permission.WRITE);
perms.add(Permission.DELETE);
perms.contains(Permission.ADMIN);  // false

// Internal representation (for enums ≤ 64 values):
// long elements = 0b00000_1_0_1_1 = READ | WRITE | DELETE
//                           │ │ │
//                           │ │ └─ READ   (ordinal 0, bit 0)
//                           │ └─── WRITE  (ordinal 1, bit 1)
//                           └───── DELETE (ordinal 3, bit 3)
```

---

### 3.5 equals() & hashCode() Contract — Deep Dive

```
THE CONTRACT (from java.lang.Object javadoc):

  1. If a.equals(b) → a.hashCode() == b.hashCode()   (MANDATORY)
  2. If a.hashCode() != b.hashCode() → !a.equals(b)  (contrapositive of 1)
  3. If a.hashCode() == b.hashCode() → maybe equals, maybe not (collision OK)
  4. a.equals(a) must be true                          (reflexive)
  5. a.equals(b) → b.equals(a)                        (symmetric)
  6. a.equals(b) && b.equals(c) → a.equals(c)         (transitive)
  7. a.equals(null) must be false                      (null safety)

WHAT BREAKS IF VIOLATED:
```

```java
// ❌ BROKEN: Override equals() without hashCode()
class Employee {
    String id;
    @Override
    public boolean equals(Object o) {
        return o instanceof Employee e && id.equals(e.id);
    }
    // hashCode NOT overridden — uses default Object.hashCode() (memory address)
}

Set<Employee> set = new HashSet<>();
Employee e1 = new Employee("EMP001");
set.add(e1);

Employee e2 = new Employee("EMP001");
set.contains(e2);  // FALSE! e2.hashCode() != e1.hashCode()
                    // despite e1.equals(e2) = true
                    // HashMap looks in WRONG bucket

// ✅ CORRECT:
class Employee {
    String id;
    @Override public boolean equals(Object o) {
        return o instanceof Employee e && id.equals(e.id);
    }
    @Override public int hashCode() {
        return Objects.hash(id);  // consistent with equals
    }
}
```

#### Senior-Level Pitfalls

```java
// PITFALL 1: Mutable field in hashCode → lost in HashMap/HashSet
class MutableKey {
    String name;
    @Override public int hashCode() { return name.hashCode(); }
    @Override public boolean equals(Object o) { /* compares name */ }
}

Set<MutableKey> set = new HashSet<>();
MutableKey key = new MutableKey("A");
set.add(key);           // stored in bucket for hash("A")
key.name = "B";         // hashCode CHANGES
set.contains(key);      // FALSE — looks in bucket for hash("B"), but key is in hash("A") bucket
set.remove(key);        // FALSE — can't find it
// Key is ORPHANED — exists in set but can never be found or removed
// → Memory leak in production

// RULE: All fields used in hashCode()/equals() must be IMMUTABLE (or effectively final)

// PITFALL 2: equals() breaks Liskov Substitution Principle (LSP)
class Point { int x, y; /* equals compares x,y */ }
class ColorPoint extends Point { Color c; /* equals compares x,y,c */ }

Point p = new Point(1, 2);
ColorPoint cp = new ColorPoint(1, 2, RED);
p.equals(cp);   // true  (Point.equals ignores color)
cp.equals(p);   // false (ColorPoint.equals sees p has no color)
// SYMMETRIC violation → HashSet/HashMap behave unpredictably

// SOLUTION: Use composition, not inheritance. Or: use records (auto-generated correct equals/hashCode)
record Point(int x, int y) {}  // correct equals/hashCode/toString for free
```

---

### 3.6 Set Performance Comparison

```
╔═══════════════════════╦═══════════╦═══════════════╦════════════╦══════════╗
║ Operation              ║ HashSet   ║ LinkedHashSet ║ TreeSet    ║ EnumSet  ║
╠═══════════════════════╬═══════════╬═══════════════╬════════════╬══════════╣
║ add(e)                ║ O(1)*     ║ O(1)*         ║ O(log N)   ║ O(1)     ║
║ contains(e)           ║ O(1)*     ║ O(1)*         ║ O(log N)   ║ O(1)     ║
║ remove(e)             ║ O(1)*     ║ O(1)*         ║ O(log N)   ║ O(1)     ║
║ Iteration order       ║ undefined ║ insertion     ║ sorted     ║ natural  ║
║ Iteration cost        ║ O(cap)    ║ O(size)       ║ O(size)    ║ O(size)  ║
║ first() / last()      ║ N/A       ║ N/A           ║ O(log N)   ║ O(1)     ║
║ floor() / ceiling()   ║ N/A       ║ N/A           ║ O(log N)   ║ N/A      ║
║ Memory/element        ║ ~32 bytes ║ ~48 bytes     ║ ~48 bytes  ║ 1 bit    ║
║ Null allowed          ║ 1 null    ║ 1 null        ║ No         ║ No       ║
║ Thread-safe           ║ No        ║ No            ║ No         ║ No       ║
╚═══════════════════════╩═══════════╩═══════════════╩════════════╩══════════╝

* Amortized O(1), O(N) worst-case with pathological hash collisions (O(log N) with treeified bins)
```

## 4. Map Implementations – Deep Dive

### 4.1 HashMap — The Most Important Collection to Understand

#### Internal Data Structure

```
HashMap<K,V>:
  ┌──────────────────────────────────────────────────────────────────────┐
  │  transient Node<K,V>[] table;    // bucket array (power of 2 length) │
  │  int size;                        // number of key-value pairs        │
  │  int threshold;                   // capacity × loadFactor            │
  │  float loadFactor;                // default 0.75                     │
  │  int modCount;                    // structural modification counter  │
  └──────────────────────────────────────────────────────────────────────┘

  table (capacity = 16, loadFactor = 0.75, threshold = 12):
  ┌─────┬──────────────────────────────────────────────────────────────┐
  │  0  │ null                                                          │
  │  1  │ → Node(hash, "key1", val1, null)                             │
  │  2  │ null                                                          │
  │  3  │ → Node → Node → Node  (chain of 3, same bucket)             │
  │  4  │ null                                                          │
  │  5  │ → TreeNode (root of Red-Black Tree, ≥8 collisions)          │
  │  6  │ null                                                          │
  │  7  │ → Node(hash, "key7", val7, null)                             │
  │ ... │ ...                                                           │
  │ 15  │ → Node → Node (chain of 2)                                   │
  └─────┴──────────────────────────────────────────────────────────────┘
```

#### Node Structure

```java
// OpenJDK HashMap.Node:
static class Node<K,V> implements Map.Entry<K,V> {
    final int hash;     // cached hash (avoid recomputation during resize)
    final K key;
    V value;
    Node<K,V> next;     // linked list pointer within bucket
    // Memory: 16 (header) + 4 (hash) + 4 (padding) + 8 (key) + 8 (value) + 8 (next) = 48 bytes
}

// TreeNode (used when bin is treeified):
static final class TreeNode<K,V> extends LinkedHashMap.Entry<K,V> {
    TreeNode<K,V> parent;  // red-black tree links
    TreeNode<K,V> left;
    TreeNode<K,V> right;
    TreeNode<K,V> prev;    // linked list for traversal
    boolean red;
    // Memory: ~104 bytes per node (significant overhead)
}
```

#### Hash Spreading Function — Why XOR with Upper Bits

```java
// OpenJDK: HashMap.hash()
static final int hash(Object key) {
    int h;
    return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
}

// WHY? Because bucket index = hash & (capacity - 1)
// If capacity = 16, index = hash & 0x0000000F → only LOWEST 4 bits matter!
//
// Problem: Many hashCode() implementations have patterns in low bits:
//   Integer.hashCode() = value itself → sequential ints cluster in same buckets
//   String.hashCode() uses polynomial → low bits may have poor distribution
//
// Solution: XOR upper 16 bits into lower 16 bits → spreads information
//
//   Original:  1010 1100 0011 0111  0000 0000 0000 1010
//   >>> 16:    0000 0000 0000 0000  1010 1100 0011 0111
//   XOR:       1010 1100 0011 0111  1010 1100 0011 1101
//                                   ^^^^^^^^^^^^^^^^
//                                   These bits now incorporate upper-bit info
//   & 0xF:                                          1101 = 13 (bucket index)
```

#### Why Capacity Is Always Power of 2

```
Bucket index calculation:
  index = hash & (capacity - 1)      // bitwise AND

This ONLY works correctly when capacity is a power of 2:
  capacity = 16 → binary: 10000
  capacity - 1  → binary: 01111   (all 1s in lower bits = perfect bitmask)

  hash     = ...1010 1101
  & 01111  = ...0000 1101 = 13    ← uniform bucket selection using lower log₂(cap) bits

If capacity were 15 (not power of 2):
  15 - 1 = 14 → binary: 01110   (bit 0 is 0!)
  ANY hash: ...xxxx xxx0 & 01110 → bit 0 ALWAYS 0 → odd buckets NEVER used
  → 50% of buckets wasted → doubled collision rate

Additional benefit of power-of-2:
  - During resize (2× capacity), each node either stays in same bucket
    or moves to bucket (old_index + old_capacity)
  - Only need to check ONE bit: hash & old_capacity == 0 ? same : moved
  - No full rehash needed — just redistribute based on one bit
```

#### Resize Operation — The Costly Event

```
Resize triggers when: size > threshold (capacity × loadFactor)

Resize process (simplified from OpenJDK):
  1. Create new Node[] of 2× capacity
  2. For each bucket in old table:
     a. If single node → recompute index: hash & (newCap - 1), place directly
     b. If linked list → split into "low" and "high" lists:
        - Low:  (hash & oldCap) == 0 → stays at index i
        - High: (hash & oldCap) != 0 → moves to index (i + oldCap)
     c. If tree → split similarly; if resulting list < UNTREEIFY_THRESHOLD → back to list
  3. Replace table reference

Cost:
  - O(N) — must visit every entry
  - Allocates new array (2× old size)
  - Old array → GC
  - ALL entries get new bucket assignment
  - Can cause latency spike in low-latency systems
```

```
Resize example (capacity 4 → 8):

Old table (cap=4):                New table (cap=8):
  [0] → A(hash=0) → D(hash=8)     [0] → A(hash=0)
  [1] → B(hash=1)                  [1] → B(hash=1)
  [2] → null                       [2] → null
  [3] → C(hash=3)                  [3] → C(hash=3)
                                    [4] → D(hash=8)   ← moved: hash & oldCap(4) = 8 & 4 = 4 ≠ 0
                                    [5] → null
                                    [6] → null
                                    [7] → null

Node D: hash=8, oldIndex = 8 & 3 = 0, newIndex = 8 & 7 = 0... 
Wait: hash & oldCap = 8 & 4 = non-zero → moves to 0 + 4 = 4 ✓
```

#### Treeification Threshold Details

```
TREEIFY_THRESHOLD = 8
  - When a bucket's linked list grows to 8 nodes, it converts to a Red-Black Tree
  - ONLY if table capacity >= MIN_TREEIFY_CAPACITY (64)
  - If capacity < 64, resize instead (spreading fixes most collisions)

UNTREEIFY_THRESHOLD = 6
  - When tree node count drops below 6 during resize split, convert back to list
  - Gap of 2 prevents thrashing

Why 8?
  - Under random hashing with loadFactor 0.75, probability of a bin having:
    0 elements: 0.4724
    1 element:  0.3543
    2 elements: 0.1329
    3 elements: 0.0332
    4 elements: 0.0062
    5 elements: 0.0009
    6 elements: 0.0001
    7 elements: 0.00001
    8 elements: 0.000001  ← extremely rare
  - If 8+ elements in a bin, it's likely a hash DoS attack or broken hashCode()
  - Tree guarantees O(log N) even under adversarial conditions
```

---

### 4.2 ConcurrentHashMap — Lock-Free Reads, Fine-Grained Writes

#### Evolution: Pre-Java 8 vs Java 8+

```
╔═══════════════════════╦══════════════════════╦════════════════════════════════╗
║ Aspect                 ║ Pre-Java 8 (Segment) ║ Java 8+ (CAS + sync bins)    ║
╠═══════════════════════╬══════════════════════╬════════════════════════════════╣
║ Locking granularity   ║ Segment (groups of   ║ Per-bin (each bucket has its  ║
║                        ║   buckets), default  ║   own lock via synchronized   ║
║                        ║   16 segments         ║   on first node in bin)      ║
║ Read mechanism         ║ volatile reads,      ║ volatile reads, Unsafe/      ║
║                        ║   no locking         ║   VarHandle, NO locking      ║
║ Write mechanism        ║ Segment.lock()       ║ CAS for empty bin, sync on   ║
║                        ║   (ReentrantLock)    ║   first node for occupied bin ║
║ Concurrency level     ║ Fixed at construction║ Dynamic — scales to # bins    ║
║ Treeification         ║ No                   ║ Yes (same as HashMap: ≥8)     ║
║ size() accuracy       ║ Locks all segments   ║ Approximate (baseCount +      ║
║                        ║                      ║   CounterCell[] for low       ║
║                        ║                      ║   contention)                 ║
║ computeIfAbsent       ║ N/A                  ║ Atomic — CAS + sync           ║
║ Null keys/values      ║ No                   ║ No (by design — ambiguity)    ║
╚═══════════════════════╩══════════════════════╩════════════════════════════════╝
```

#### Java 8+ Internal Mechanism

```
ConcurrentHashMap — put(key, value):

  1. Compute hash = spread(key.hashCode())
     → Same spreading as HashMap + additionally: (h ^ (h >>> 16)) & 0x7fffffff
     → Ensures positive hash (negative used for special nodes: forwarding, tree)

  2. Find bucket: tabAt(tab, i = (n - 1) & hash)
     → Unsafe/VarHandle volatile read of table[i]

  3. If bucket is EMPTY:
     → CAS (Compare-And-Swap) to install new Node
     → No lock needed! If CAS fails (another thread beat us), retry loop

  4. If bucket is OCCUPIED:
     → synchronized (firstNodeOfBin) {    // lock ONLY this bin
           // Walk chain or tree
           // Insert/update node
       }

  5. If chain length >= TREEIFY_THRESHOLD:
     → treeifyBin() under synchronization

  6. addCount() → increment size:
     → Uses LongAdder-style Counter (baseCount + CounterCell[])
     → Very low contention for size tracking

  Read operations (get, containsKey):
     → ZERO locking
     → volatile read of table, then volatile read of node
     → Memory visibility guaranteed by Java Memory Model
```

```java
// computeIfAbsent internals — ATOMIC guarantee
ConcurrentHashMap<String, List<Order>> orderCache = new ConcurrentHashMap<>();

// Thread-safe atomic insertion:
List<Order> orders = orderCache.computeIfAbsent("customer-123", key -> {
    return loadOrdersFromDB(key);  // Called AT MOST ONCE per key
});
// ⚠️ WARNING: The mapping function must be SHORT and FAST
//    It runs INSIDE the synchronized block on the bin
//    A slow function (DB call, HTTP) blocks other threads accessing same bin!
//    Alternative: use putIfAbsent with pre-computed value, or async with CF
```

#### ConcurrentHashMap — Why No Nulls?

```
HashMap allows 1 null key and null values.
ConcurrentHashMap bans BOTH. Why?

Ambiguity under concurrency:
  map.get(key) returns null → Does it mean:
    (a) Key exists with null value?
    (b) Key doesn't exist?

  In HashMap: use containsKey() to distinguish → safe (single-threaded)
  In ConcurrentHashMap: between get() and containsKey(), another thread
    may have removed/added the key → TOCTOU race condition
  
  Doug Lea: "The main reason that nulls aren't allowed in ConcurrentMaps
   is that ambiguities that may be just barely tolerable in non-concurrent
   maps can't be accommodated."
```

#### Weakly Consistent Iterators

```
ConcurrentHashMap iterators:
  - Do NOT throw ConcurrentModificationException
  - Reflect the state of the map at some point at or since iterator creation
  - May or may not reflect concurrent modifications
  - Guaranteed to traverse elements as they existed upon construction
  - Each element returned AT MOST ONCE (no duplicates)

This is a DESIGN TRADE-OFF:
  Strong consistency → requires global lock during iteration → kills concurrency
  Weak consistency → lock-free iteration → high throughput
```

---

### 4.3 LinkedHashMap — LRU Cache Pattern

#### Internal Structure

```
LinkedHashMap<K,V> extends HashMap<K,V>:
  ┌──────────────────────────────────────────────────────────────────┐
  │  HashMap buckets (for O(1) lookup)                               │
  │  + Doubly-linked list threading through ALL entries              │
  │                                                                  │
  │  Has two modes:                                                  │
  │    accessOrder = false (default): linked list = insertion order   │
  │    accessOrder = true:  linked list = access order (LRU → MRU)  │
  │                                                                  │
  │  Entry extends HashMap.Node, adds:                               │
  │    Entry<K,V> before, after;   // doubly-linked list pointers   │
  └──────────────────────────────────────────────────────────────────┘

  accessOrder = true:
  After get("C") or put("C", newVal):
    Before: head ↔ A ↔ B ↔ C ↔ D ↔ tail
    After:  head ↔ A ↔ B ↔ D ↔ C ↔ tail   (C moved to end = most recently used)
```

#### LRU Cache with LinkedHashMap

```java
// Classic LRU cache — override removeEldestEntry
public class LRUCache<K, V> extends LinkedHashMap<K, V> {
    private final int maxSize;
    
    public LRUCache(int maxSize) {
        super(
            maxSize + 1,    // initial capacity (avoid immediate resize)
            0.75f,          // load factor
            true            // accessOrder = true (CRITICAL for LRU)
        );
        this.maxSize = maxSize;
    }
    
    @Override
    protected boolean removeEldestEntry(Map.Entry<K, V> eldest) {
        return size() > maxSize;  // auto-evict LRU when over capacity
    }
}

// Usage:
LRUCache<String, UserProfile> cache = new LRUCache<>(1000);
cache.put("user-1", profile1);
cache.get("user-1");             // moves to MRU position
cache.put("user-1001", profile); // if size > 1000, LRU entry auto-removed

// Thread-safe version:
Map<String, UserProfile> syncCache = Collections.synchronizedMap(new LRUCache<>(1000));
// Better: use Caffeine or Guava Cache for production LRU

// removeEldestEntry is called after EVERY put():
//   put() → addEntry → if (removeEldestEntry(eldest)) remove eldest
```

---

### 4.4 TreeMap — Sorted Map with Red-Black Tree

```
TreeMap<K,V>:
  - Backed by a Red-Black Tree (same as TreeSet)
  - Keys sorted by natural order (Comparable) or provided Comparator
  - All operations: O(log N)
  - Implements NavigableMap → range queries, floor/ceiling, etc.

  Internal structure:
  ┌─────────────────────────────────────────────────────────┐
  │  Entry<K,V> root;        // root of red-black tree       │
  │  Comparator<? super K> comparator; // null = natural order│
  │  int size;                                                │
  │  int modCount;                                            │
  └─────────────────────────────────────────────────────────┘

  Entry node:
  ┌─────────────────────────┐
  │  K key                   │
  │  V value                 │
  │  Entry left, right, parent│
  │  boolean color (RED/BLACK)│
  └─────────────────────────┘
  Memory: ~56 bytes per entry
```

```java
TreeMap<LocalDate, BigDecimal> dailySales = new TreeMap<>();
dailySales.put(LocalDate.of(2025, 1, 15), new BigDecimal("5000"));
dailySales.put(LocalDate.of(2025, 1, 20), new BigDecimal("7500"));
dailySales.put(LocalDate.of(2025, 2, 1),  new BigDecimal("6000"));
dailySales.put(LocalDate.of(2025, 2, 10), new BigDecimal("8000"));

// Range query — sales in January:
SortedMap<LocalDate, BigDecimal> janSales = dailySales.subMap(
    LocalDate.of(2025, 1, 1),   // inclusive
    LocalDate.of(2025, 2, 1)    // exclusive
);

// Navigation:
Map.Entry<LocalDate, BigDecimal> latest = dailySales.lastEntry();
Map.Entry<LocalDate, BigDecimal> beforeFeb = 
    dailySales.lowerEntry(LocalDate.of(2025, 2, 1));  // Jan 20 entry
```

---

### 4.5 WeakHashMap — GC-Friendly Mapping

```
WeakHashMap<K,V>:
  ┌─────────────────────────────────────────────────────────────────┐
  │  Keys are wrapped in WeakReference<K>                           │
  │  When key has no strong references elsewhere → GC collects key  │
  │  On next access to map, expired entries are silently cleaned up  │
  │                                                                 │
  │  Internal: ReferenceQueue<Object> queue; // GC enqueues expired │
  │  Before each operation: expungeStaleEntries() drains queue      │
  └─────────────────────────────────────────────────────────────────┘

  Lifecycle:
  ┌───────────────┐   ┌──────────┐   ┌───────────┐   ┌────────────┐
  │ Strong ref to │──►│ WeakRef  │   │ GC runs,  │   │ Next map   │
  │ key exists    │   │ in map   │   │ collects  │   │ operation  │
  │ → entry lives │   │ is alive │   │ weak key  │   │ cleans up  │
  └───────────────┘   └──────────┘   └───────────┘   └────────────┘
       active              active        collected       cleaned
```

```java
// Use case: Metadata cache that auto-cleans when objects are GC'd
WeakHashMap<ClassLoader, Map<String, Class<?>>> classCache = new WeakHashMap<>();

// When a ClassLoader is unloaded (no strong references):
//   → Weak key becomes eligible for GC
//   → Entry automatically removed from cache
//   → No memory leak from stale ClassLoader metadata

// ⚠️ CAUTION: String literals and small Integer/Long are interned
//   → They always have strong references → NEVER get GC'd from WeakHashMap
WeakHashMap<String, String> map = new WeakHashMap<>();
map.put("hello", "world");  // "hello" is interned → NEVER collected
map.put(new String("hello"), "world"); // new String → eligible for GC
```

---

### 4.6 Map Concurrency Comparison

```
╔══════════════════════╦════════════╦════════════════╦═══════════════╦══════════════╗
║ Feature               ║ HashMap    ║ ConcurrentHM   ║ Hashtable     ║ synced Map   ║
╠══════════════════════╬════════════╬════════════════╬═══════════════╬══════════════╣
║ Thread-safe          ║ No         ║ Yes             ║ Yes           ║ Yes          ║
║ Lock granularity     ║ N/A        ║ Per-bin (CAS +  ║ Entire map    ║ Entire map   ║
║                       ║            ║   sync)        ║ (every method)║ (wrapper)    ║
║ Null key             ║ 1 allowed  ║ Not allowed     ║ Not allowed   ║ 1 allowed    ║
║ Null value           ║ Allowed    ║ Not allowed     ║ Not allowed   ║ Allowed      ║
║ Iterator             ║ Fail-fast  ║ Weakly          ║ Fail-fast     ║ Fail-fast    ║
║                       ║            ║   consistent    ║               ║              ║
║ Atomic operations    ║ No         ║ Yes (compute,   ║ No            ║ No           ║
║                       ║            ║ merge, putIf)  ║               ║              ║
║ Read throughput      ║ Highest    ║ Near HashMap    ║ Low           ║ Low          ║
║   (multi-threaded)   ║ (unsafe)   ║ (no locks)     ║ (global lock) ║ (global lock)║
║ Write throughput     ║ N/A        ║ High            ║ Low           ║ Low          ║
║   (multi-threaded)   ║ (unsafe)   ║ (per-bin lock) ║ (global lock) ║ (global lock)║
╚══════════════════════╩════════════╩════════════════╩═══════════════╩══════════════╝

synced Map = Collections.synchronizedMap(new HashMap<>())
```

---

### 4.7 Real Production Debugging Examples

```java
// INCIDENT 1: CPU spike due to hash collision attack (pre-Java 8)
// Symptom: 100% CPU on web server processing POST requests
// Root cause: Attacker sent thousands of keys with SAME hashCode → all in one bucket
//   → O(N²) lookup in single linked-list bucket
// Fix: Java 8 treeification limits damage to O(N log N)
//   + Use randomized hash seed (HashMap.hash() already adds entropy)
//   + WAF → limit request parameter count

// INCIDENT 2: Memory leak from non-removed ConcurrentHashMap entries
// Symptom: OldGen growing unboundedly, eventually OOM after days
ConcurrentHashMap<String, SessionData> sessions = new ConcurrentHashMap<>();
// Sessions were added but never removed after user logout
// Fix: Use computeIfPresent with TTL check, or Caffeine with expireAfterAccess

// INCIDENT 3: Deadlock in computeIfAbsent
ConcurrentHashMap<String, ConcurrentHashMap<String, String>> nested = new ConcurrentHashMap<>();
nested.computeIfAbsent("outer", k -> {
    // THIS DEADLOCKS if "outer" maps to same bin as some other operation
    // computeIfAbsent holds synchronized on the bin while executing the function
    nested.computeIfAbsent("inner", k2 -> "value");  // tries to lock another bin
    return new ConcurrentHashMap<>();
});
// Fix: NEVER access the same ConcurrentHashMap inside its own computeIfAbsent
//   Alternative: use putIfAbsent with pre-created value

// INCIDENT 4: IdentityHashMap confusion in serialization
// IdentityHashMap uses == instead of equals() for key comparison
// A deserialize → new object → == fails even if equals() returns true
// Fix: Use HashMap unless you specifically need reference identity semantics
```

---

## Theory D: Java Memory Model, Object Layout & Reference Types

> Understanding how objects are laid out in memory and how references work is critical for reasoning about collection memory overhead, GC behavior, and WeakHashMap/SoftReference patterns.

### Object Memory Layout in HotSpot JVM

Every Java object has a header, then fields, then optional padding:

```
  64-bit JVM with Compressed Oops (default for heap < 32GB):
  ┌─────────────────────────────────────────────────┐
  │  Mark Word              (8 bytes)             │  ← hash, GC age, lock state
  │  Klass Pointer           (4 bytes compressed)  │  ← pointer to class metadata
  │  [Array Length]          (4 bytes, arrays only) │
  │  Instance Fields         (varies)              │
  │  Padding                 (to 8-byte boundary)  │
  └─────────────────────────────────────────────────┘

  Minimum object size: 16 bytes (12-byte header + 4-byte padding)
  
  Example: Integer object
  ┌────────────────────────┐
  │ Mark Word     (8 bytes) │
  │ Klass Ptr     (4 bytes) │
  │ int value     (4 bytes) │
  └────────────────────────┘
  Total: 16 bytes for one int (4× overhead vs primitive)

  Example: String object (Java 9+ compact strings)
  ┌──────────────────────────────┐
  │ Mark Word         (8 bytes)  │
  │ Klass Ptr         (4 bytes)  │
  │ byte[] value ref  (4 bytes)  │  ← reference to backing array
  │ int hash           (4 bytes)  │  ← cached hashCode
  │ byte coder         (1 byte)  │  ← LATIN1 or UTF16
  │ padding            (3 bytes)  │
  └──────────────────────────────┘
  String object: 24 bytes + byte[] array overhead  
```

### Why This Matters for Collections

```
  HashMap<Integer, String> with 1 million entries:

  Per entry:
    HashMap.Node object: 32 bytes (header + hash + key ref + value ref + next ref)
    Integer key object:  16 bytes
    String value object: ~24 bytes + byte[] for chars
    ───────────────────────────────────
    Minimum per entry:   72+ bytes

  1M entries: ~72 MB just for the map structure
  Plus bucket array: 1M * 1.33 (load factor inv.) * 4 bytes = ~5.3 MB
  Total: ~77 MB for 1M int→String mappings

  Primitive int→String without boxing: ~30 MB (Eclipse Collections IntObjectHashMap)
  Savings: 60%+ memory reduction
```

### Java Reference Types (Strong → Phantom)

| Reference Type | GC Behavior | Collected When | Use Case |
|---|---|---|---|
| **Strong** (`T ref`) | Never collected while reachable | Only when unreachable | Default. All normal variables |
| **Soft** (`SoftReference<T>`) | Collected when memory is low | Before OutOfMemoryError | Memory-sensitive caches |
| **Weak** (`WeakReference<T>`) | Collected at next GC cycle | When only weakly reachable | WeakHashMap, canonicalization |
| **Phantom** (`PhantomReference<T>`) | Never accessible via `.get()` | After finalization | Resource cleanup (replacing finalizers) |

```
  Reference strength hierarchy:
  
  Strong > Soft > Weak > Phantom
  
  GC decision tree:
  ┌──────────────────────────────────────────────┐
  │ Object has strong reference?                    │
  │   YES → KEEP (never collected)                  │
  │   NO → Has soft reference?                       │
  │          YES → KEEP unless memory pressure       │
  │          NO → Has weak reference?                 │
  │                 YES → COLLECT at next GC          │
  │                 NO → Has phantom reference?        │
  │                        YES → Enqueue, then COLLECT│
  │                        NO → COLLECT immediately   │
  └──────────────────────────────────────────────┘
```

```java
// WeakReference example — how WeakHashMap works internally:
Object key = new Object();
WeakReference<Object> weakRef = new WeakReference<>(key);

System.out.println(weakRef.get());  // → the object
key = null;  // remove strong reference
System.gc();  // hint to GC
System.out.println(weakRef.get());  // → null (collected!)

// This is exactly what WeakHashMap does:
// Keys are wrapped in WeakReference
// When key has no more strong references → GC collects it
// WeakHashMap.expungeStaleEntries() removes the mapping

// SoftReference for caching:
SoftReference<byte[]> cache = new SoftReference<>(loadLargeData());
byte[] data = cache.get();
if (data == null) {
    data = loadLargeData();  // reload, it was GC'd under memory pressure
    cache = new SoftReference<>(data);
}
```

### Mark Word Deep Dive (What's Stored)

```
  Mark Word (64-bit JVM):
  ┌─────────────────────────────────────────────────┐
  │ State          │  Bits Layout                     │
  ├────────────────┼────────────────────────────────┤
  │ Unlocked       │ [hashCode:31][age:4][biased:0][01]│
  │ Biased Lock    │ [threadId:54][epoch:2][age:4][1][01]│
  │ Lightweight    │ [ptr to lock record :62][00]      │
  │ Heavyweight    │ [ptr to monitor   :62][10]      │
  │ GC Marked      │ [forwarding address:62][11]      │
  └────────────────┴────────────────────────────────┘

  This is why:
  - hashCode() is "free" for identity hashing (stored in header)
  - synchronized works per-object (lock state in header)
  - GC can track object age for generational collection (4-bit age)
  - Object minimum size = 16 bytes (header alone = 12, padded to 16)
```

---

## 5. Queue & Deque Implementations

### 5.1 PriorityQueue — Binary Heap Internals

#### Structure

```
PriorityQueue<E>:
  ┌──────────────────────────────────────────────────────────────┐
  │  Object[] queue;          // binary heap stored in array      │
  │  int size;                                                    │
  │  Comparator<? super E> comparator;  // null = natural order   │
  └──────────────────────────────────────────────────────────────┘

  Binary min-heap property: parent ≤ both children (for min-heap)

  Array representation of heap:
    Parent of node at index i:  (i - 1) >>> 1   (i.e., (i-1)/2)
    Left child of index i:      2*i + 1
    Right child of index i:     2*i + 2

  Example: PriorityQueue containing {1, 3, 5, 7, 9, 11, 13}

  Logical tree:              Array layout:
        1                    ┌───┬───┬───┬───┬───┬────┬────┬─────┐
       / \                   │ 1 │ 3 │ 5 │ 7 │ 9 │ 11 │ 13 │ ... │
      3   5                  └───┴───┴───┴───┴───┴────┴────┴─────┘
     / \ / \                  [0] [1] [2] [3] [4]  [5]  [6]
    7  9 11 13

  offer(element):
    1. Add to end of array (index = size)
    2. Sift-up: while element < parent, swap with parent
    3. O(log N) — height of tree

  poll() → removes and returns minimum:
    1. Save root (index 0) as result
    2. Move last element to root
    3. Sift-down: while element > smaller child, swap with smaller child
    4. O(log N)

  peek():
    1. Return queue[0]   — O(1)
```

```java
// Common usage: Top-K elements, scheduling, Dijkstra's algorithm
PriorityQueue<Integer> minHeap = new PriorityQueue<>();  // min at head
PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Comparator.reverseOrder());

// Custom priority:
record Task(String name, int priority, Instant deadline) {}
PriorityQueue<Task> taskQueue = new PriorityQueue<>(
    Comparator.comparingInt(Task::priority)
              .thenComparing(Task::deadline)
);

// ⚠️ PriorityQueue is NOT sorted — only guarantees heap property
// iterator() does NOT return elements in priority order!
// To get sorted: poll() repeatedly (O(N log N) total)
// PriorityQueue is NOT thread-safe → use PriorityBlockingQueue for concurrency
```

#### Complexity

| Operation | Complexity | Notes |
|---|---|---|
| `offer(e)` / `add(e)` | O(log N) | Sift-up |
| `poll()` / `remove()` | O(log N) | Sift-down |
| `peek()` / `element()` | O(1) | Array index 0 |
| `remove(Object)` | O(N) | Linear scan + sift |
| `contains(Object)` | O(N) | Linear scan |
| Growth | 50% if small (<64), else 2× | `grow()` similar to ArrayList |

---

### 5.2 ArrayDeque — Circular Array

#### Internal Structure

```
ArrayDeque<E>:
  ┌───────────────────────────────────────────────────────────────────┐
  │  Object[] elements;   // circular buffer                          │
  │  int head;            // index of first element                   │
  │  int tail;            // index AFTER last element                 │
  └───────────────────────────────────────────────────────────────────┘

  Circular array (capacity = 8):
                    head=2        tail=6
                      ↓             ↓
  ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┐
  │ null │ null │  A   │  B   │  C   │  D   │ null │ null │
  └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┘
    [0]    [1]    [2]    [3]    [4]    [5]    [6]    [7]

  After addFirst("X"):   head wraps around: head = (head - 1) & (length - 1)
                    head=1
                      ↓                            tail=6
  ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┐
  │ null │  X   │  A   │  B   │  C   │  D   │ null │ null │
  └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┘

  After addLast("Y"):                              tail=7
  ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┐
  │ null │  X   │  A   │  B   │  C   │  D   │  Y   │ null │
  └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┘

  Wrap-around: head and tail use bitwise AND: index & (length - 1)
  → Same power-of-2 trick as HashMap
```

#### Why ArrayDeque Is Better Than Stack and LinkedList

```
╔══════════════════════╦════════════════╦═══════════════╦═══════════════════╗
║ Property              ║ ArrayDeque     ║ Stack (Vector) ║ LinkedList       ║
╠══════════════════════╬════════════════╬═══════════════╬═══════════════════╣
║ Push/Pop             ║ O(1) amortized ║ O(1) amortized║ O(1)             ║
║ Synchronized         ║ No             ║ Yes (always!)  ║ No               ║
║ Memory/element       ║ ~8 bytes       ║ ~8 bytes + sync║ ~40 bytes (node) ║
║ Cache-friendly       ║ Yes (array)    ║ Yes (array)    ║ No (scattered)   ║
║ Random access        ║ No (by design) ║ Yes (breaks    ║ O(N) — slow      ║
║                       ║                ║   stack API)   ║                   ║
║ Null elements        ║ No             ║ Yes            ║ Yes              ║
║ Iterator order       ║ Head → Tail    ║ Bottom → Top   ║ First → Last     ║
╚══════════════════════╩════════════════╩═══════════════╩═══════════════════╝

Javadoc: "This class is likely to be faster than Stack when used as a stack,
          and faster than LinkedList when used as a queue."
```

---

### 5.3 BlockingQueue Implementations — Producer/Consumer

#### Overview Comparison

```
╔═══════════════════════════╦═════════════╦═══════════╦════════════════════════════════╗
║ Implementation             ║ Bounded?    ║ Ordering  ║ Internal Structure             ║
╠═══════════════════════════╬═════════════╬═══════════╬════════════════════════════════╣
║ ArrayBlockingQueue        ║ Yes (fixed) ║ FIFO      ║ Circular array + single lock   ║
║ LinkedBlockingQueue       ║ Optional    ║ FIFO      ║ Linked nodes + TWO locks       ║
║ PriorityBlockingQueue     ║ Unbounded   ║ Priority  ║ Binary heap + single lock      ║
║ DelayQueue                ║ Unbounded   ║ Delay     ║ PriorityQueue + leader thread  ║
║ SynchronousQueue          ║ No capacity ║ N/A       ║ Transfer — no storage          ║
║ LinkedTransferQueue       ║ Unbounded   ║ FIFO      ║ Dual queue (lock-free)         ║
╚═══════════════════════════╩═════════════╩═══════════╩════════════════════════════════╝
```

#### ArrayBlockingQueue vs LinkedBlockingQueue

```
ArrayBlockingQueue:
  ┌──────────────────────────────────────────────────────┐
  │  Object[] items;         // fixed-size circular array │
  │  int takeIndex, putIndex; // head and tail pointers   │
  │  int count;                                           │
  │  ReentrantLock lock;     // SINGLE lock for both ops  │
  │  Condition notEmpty;     // signals waiting consumers │
  │  Condition notFull;      // signals waiting producers │
  └──────────────────────────────────────────────────────┘
  → Single lock = producers and consumers contend with each other
  → Fixed capacity → predictable memory usage

LinkedBlockingQueue:
  ┌──────────────────────────────────────────────────────┐
  │  Node<E> head, last;                                  │
  │  AtomicInteger count;                                 │
  │  ReentrantLock takeLock;  // separate lock for takes  │
  │  ReentrantLock putLock;   // separate lock for puts   │
  │  Condition notEmpty, notFull;                         │
  └──────────────────────────────────────────────────────┘
  → TWO locks = producers and consumers can operate concurrently
  → Higher throughput when both putting and taking simultaneously
  → Default capacity: Integer.MAX_VALUE (effectively unbounded — dangerous!)
```

#### Producer-Consumer Pattern

```java
// Classic producer-consumer with BlockingQueue:
BlockingQueue<Task> queue = new ArrayBlockingQueue<>(1000);

// Producer thread:
Runnable producer = () -> {
    while (running) {
        Task task = generateTask();
        queue.put(task);       // BLOCKS if queue is full — back-pressure!
    }
};

// Consumer thread:
Runnable consumer = () -> {
    while (running) {
        Task task = queue.take();  // BLOCKS if queue is empty — waits for work
        process(task);
    }
};

// Graceful shutdown pattern:
// Producer: queue.put(POISON_PILL);
// Consumer: if (task == POISON_PILL) break;

// ⚠️ Common production pitfall:
// Using unbounded LinkedBlockingQueue as default → OOM under burst load
// ThreadPoolExecutor default: new LinkedBlockingQueue<>() — unbounded!
// Fix: Always specify capacity: new LinkedBlockingQueue<>(10_000)
```

#### DelayQueue — Scheduled Task Processing

```java
record ScheduledTask(String name, Instant executeAt) implements Delayed {
    @Override
    public long getDelay(TimeUnit unit) {
        return unit.convert(Duration.between(Instant.now(), executeAt));
    }
    @Override
    public int compareTo(Delayed other) {
        return Long.compare(this.getDelay(TimeUnit.MILLISECONDS), 
                           other.getDelay(TimeUnit.MILLISECONDS));
    }
}

DelayQueue<ScheduledTask> scheduler = new DelayQueue<>();
scheduler.put(new ScheduledTask("retry-email", Instant.now().plusSeconds(30)));
scheduler.put(new ScheduledTask("cleanup-temp", Instant.now().plusMinutes(5)));

// Consumer: blocks until task's delay has expired
ScheduledTask task = scheduler.take();  // waits until earliest task is ready
execute(task);
```

#### SynchronousQueue — Direct Handoff

```
SynchronousQueue:
  - Has ZERO capacity — not really a "queue"
  - put() blocks until another thread calls take()
  - Direct handoff between producer and consumer
  - Used by Executors.newCachedThreadPool() internally
  
  Producer: put(task) ──BLOCKS──► Consumer: take() receives task immediately
  
  Two modes:
    fair = false (default): LIFO stack (faster, unfair)
    fair = true:  FIFO queue (guaranteed ordering, slightly slower)
```

---

### 5.4 Queue/Deque Performance Comparison

```
╔═══════════════════════════╦═══════════╦═════════════╦═══════════════╦═══════════╗
║ Operation                  ║ ArrayDeque║ PriorityQ   ║ ArrayBlockQ   ║ LinkedBQ  ║
╠═══════════════════════════╬═══════════╬═════════════╬═══════════════╬═══════════╣
║ offer / add (tail)        ║ O(1)*     ║ O(log N)    ║ O(1)          ║ O(1)      ║
║ poll / remove (head)      ║ O(1)      ║ O(log N)    ║ O(1)          ║ O(1)      ║
║ peek                      ║ O(1)      ║ O(1)        ║ O(1)          ║ O(1)      ║
║ remove(Object)            ║ O(N)      ║ O(N)        ║ O(N)          ║ O(N)      ║
║ contains(Object)          ║ O(N)      ║ O(N)        ║ O(N)          ║ O(N)      ║
║ Thread-safe               ║ No        ║ No          ║ Yes           ║ Yes       ║
║ Blocking                  ║ No        ║ No          ║ Yes           ║ Yes       ║
║ Lock scheme               ║ N/A       ║ N/A         ║ 1 lock        ║ 2 locks   ║
║ Memory/element            ║ ~8 bytes  ║ ~8 bytes    ║ ~8 bytes      ║ ~40 bytes ║
║ Bounded                   ║ No        ║ No          ║ Yes (fixed)   ║ Optional  ║
╚═══════════════════════════╩═══════════╩═════════════╩═══════════════╩═══════════╝

* O(1) amortized, O(N) worst case during resize
```

---

## Theory E: Concurrency Primitives — Foundations for Thread-Safe Collections

> Understanding these primitives is essential before studying ConcurrentHashMap, BlockingQueue, CopyOnWriteArrayList, and Section 6.

### Thread Safety — What Does It Actually Mean?

A class is **thread-safe** if it behaves correctly when accessed from multiple threads simultaneously, **with no additional synchronization by the caller**.

```
  Three conditions for a data race:
  1. Two or more threads access the same variable
  2. At least one access is a write
  3. No synchronization orders the accesses

  All three must be true for a race condition. Eliminate ANY one to be safe.
```

### Mutual Exclusion: synchronized & Locks

```java
// 1. Intrinsic lock (monitor) via synchronized
synchronized (lockObject) {
    // only ONE thread can be here at a time for this lockObject
    map.put(key, value);
}
// Lock is released automatically (even on exception)

// 2. ReentrantLock — more flexible
ReentrantLock lock = new ReentrantLock();
lock.lock();
try {
    map.put(key, value);
} finally {
    lock.unlock();  // MUST be in finally!
}

// Advantages of ReentrantLock over synchronized:
// - tryLock() → non-blocking attempt
// - lockInterruptibly() → respond to interruption
// - Condition variables (await/signal vs wait/notify)
// - Fair lock option (FIFO ordering)
```

### Visibility: volatile & happens-before

```
  Problem: Without synchronization, Thread B may never see Thread A's write

  Thread A:  running = true      │  Thread B:  while (running) { ... }
                                  │
  CPU Cache A: running = true   │  CPU Cache B: running = false  ← STALE!
                                  │
  Main Memory: running = true (eventually) but B may never read it

  Solution: volatile keyword
  volatile boolean running = true;
  → Write to volatile → immediately flushed to main memory
  → Read of volatile → always reads from main memory
  → Establishes happens-before relationship
```

```
  Java Memory Model (JMM) Happens-Before Rules:
  ╔════════════════════════════════════════════════════════════════════╗
  ║ 1. Program order: actions in a thread happen in sequence       ║
  ║ 2. Monitor lock: unlock HB subsequent lock of same monitor    ║
  ║ 3. Volatile: write HB subsequent read of same variable        ║
  ║ 4. Thread start: Thread.start() HB any action in new thread  ║
  ║ 5. Thread join: all actions in thread HB return from join()   ║
  ║ 6. Transitivity: if A HB B and B HB C, then A HB C           ║
  ╚════════════════════════════════════════════════════════════════════╝
  
  HB = Happens-Before
  → This is how ConcurrentHashMap guarantees visibility WITHOUT global locks:
    - volatile reads for table reference
    - CAS operations imply memory barriers
    - synchronized on individual bins (tree operations)
```

### CAS (Compare-And-Swap) — Lock-Free Foundation

```
  CAS is a CPU-level atomic instruction:
  
  CAS(address, expectedValue, newValue):
    atomically:
      if *address == expectedValue:
          *address = newValue
          return true (success)
      else:
          return false (someone else changed it)

  Retry loop (optimistic locking):
  do {
      oldValue = read(address)
      newValue = compute(oldValue)
  } while (!CAS(address, oldValue, newValue))  // retry if someone beat us
```

```java
// Java provides CAS via atomic classes:
AtomicInteger counter = new AtomicInteger(0);

// Thread-safe increment without locks:
counter.incrementAndGet();  // internally uses CAS loop

// Manual CAS:
int current, next;
do {
    current = counter.get();
    next = current + 1;
} while (!counter.compareAndSet(current, next));

// CAS vs Locks:
// CAS:  No thread blocking, no context switch, no deadlock possible
//       But: high contention → spin-wait wastes CPU ("CAS storm")
// Lock: Thread sleeps on contention (saves CPU)
//       But: context switch cost, deadlock risk, priority inversion
```

### How Collections Use These Primitives

| Collection | Synchronization Mechanism | Details |
|---|---|---|
| `Collections.synchronizedMap()` | Single intrinsic lock | `synchronized(mutex)` on every method |
| `Hashtable` | `synchronized` on every method | Essentially same as above |
| `ConcurrentHashMap` (Java 8+) | CAS + per-bin `synchronized` | Reads are lock-free (volatile), writes lock only the bin |
| `CopyOnWriteArrayList` | `ReentrantLock` + volatile array | Lock on write, volatile for read visibility |
| `ConcurrentLinkedQueue` | Lock-free CAS | Michael-Scott queue algorithm |
| `LinkedBlockingQueue` | Two `ReentrantLock`s | Separate locks for head (take) and tail (put) |
| `ArrayBlockingQueue` | Single `ReentrantLock` | One lock for both put and take |

### ABA Problem (Advanced CAS Pitfall)

```
  Thread 1: read value A, prepare to CAS(A → C)
  Thread 2: changes A → B → A  (changed and changed back!)
  Thread 1: CAS succeeds (sees A), but the "A" is not the SAME A!

  Solution: AtomicStampedReference<T>
    CAS checks BOTH value AND stamp (version number)
    Even if value reverts to A, stamp 1 ≠ stamp 3

  In practice: Java's ConcurrentHashMap avoids ABA issues because
  node objects are never reused — new nodes are always allocated.
```

---

## 6. Concurrency & Collections

### 6.1 Fail-Fast vs Fail-Safe (Weakly Consistent) Iterators

```
╔═══════════════════════╦══════════════════════════════╦═══════════════════════════════╗
║ Property               ║ Fail-Fast                    ║ Weakly Consistent (Fail-Safe) ║
╠═══════════════════════╬══════════════════════════════╬═══════════════════════════════╣
║ Throws CME?           ║ Yes — on structural mod      ║ Never                         ║
║ Detection             ║ modCount != expectedModCount║ N/A                           ║
║ Data snapshot          ║ No — sees live structure     ║ Yes — snapshot or eventual    ║
║ Used by               ║ ArrayList, HashMap, HashSet, ║ ConcurrentHashMap,            ║
║                        ║ LinkedList, TreeMap          ║ CopyOnWriteArrayList,         ║
║                        ║                              ║ ConcurrentLinkedQueue          ║
║ Guarantees            ║ "Best effort" only — not     ║ Elements returned at most     ║
║                        ║ guaranteed under races       ║ once, no CME, eventual data   ║
║ Concurrent safe?      ║ No — designed for detection  ║ Yes — designed for concurrency║
╚═══════════════════════╩══════════════════════════════╩═══════════════════════════════╝
```

```java
// Fail-fast in action:
List<String> list = new ArrayList<>(List.of("A", "B", "C"));
for (String s : list) {
    if (s.equals("B")) list.remove(s);  // ConcurrentModificationException!
}

// ✅ Fix 1: Use Iterator.remove()
Iterator<String> it = list.iterator();
while (it.hasNext()) {
    if (it.next().equals("B")) it.remove();  // safe
}

// ✅ Fix 2: Use removeIf (Java 8+)
list.removeIf(s -> s.equals("B"));

// ✅ Fix 3: Use CopyOnWriteArrayList (if concurrent)
List<String> cowList = new CopyOnWriteArrayList<>(List.of("A", "B", "C"));
for (String s : cowList) {
    cowList.remove(s);  // No CME — iterator uses snapshot
}
```

---

### 6.2 Synchronized Wrappers vs Concurrent Collections

```
Synchronized wrappers (Collections.synchronizedXxx):
  ┌──────────────────────────────────────────────────────────────────┐
  │  Every method call acquires the same mutex (the wrapper object)  │
  │  Collections.synchronizedMap(new HashMap<>());                   │
  │                                                                  │
  │  Problems:                                                       │
  │  1. Global lock → only ONE thread can access at a time           │
  │  2. Compound operations are NOT atomic:                          │
  │     if (!map.containsKey(k)) map.put(k, v);  // RACE CONDITION! │
  │  3. Iteration requires manual synchronization:                   │
  │     synchronized (map) {                                         │
  │         for (var entry : map.entrySet()) { ... }  // must hold lock │
  │     }                                                            │
  │  4. Throughput collapses under contention                        │
  └──────────────────────────────────────────────────────────────────┘

Concurrent collections (java.util.concurrent):
  ┌──────────────────────────────────────────────────────────────────┐
  │  Fine-grained locking or lock-free algorithms                    │
  │  Atomic compound operations (computeIfAbsent, merge, etc.)      │
  │  Weakly consistent iterators (no CME, no external sync needed)  │
  │  Much higher throughput under contention                         │
  └──────────────────────────────────────────────────────────────────┘
```

```java
// ❌ Synchronized wrapper — compound operation race:
Map<String, List<Order>> syncMap = Collections.synchronizedMap(new HashMap<>());

// Thread 1 and Thread 2 simultaneously:
if (!syncMap.containsKey("customer-1")) {     // TOCTOU: check-then-act race
    syncMap.put("customer-1", new ArrayList<>());
}
syncMap.get("customer-1").add(newOrder);

// ✅ ConcurrentHashMap — atomic compound operation:
ConcurrentHashMap<String, List<Order>> concMap = new ConcurrentHashMap<>();
concMap.computeIfAbsent("customer-1", k -> new CopyOnWriteArrayList<>()).add(newOrder);
// computeIfAbsent is ATOMIC — no race condition
```

---

### 6.3 When To Use Which Concurrent Collection

```
╔═══════════════════════════════════╦═══════════════════════════════════════════════╗
║ Requirement                        ║ Recommended Collection                        ║
╠═══════════════════════════════════╬═══════════════════════════════════════════════╣
║ Concurrent key-value store        ║ ConcurrentHashMap                             ║
║ Concurrent sorted map             ║ ConcurrentSkipListMap                         ║
║ Read-heavy, write-rare list       ║ CopyOnWriteArrayList                          ║
║ Read-heavy, write-rare set        ║ CopyOnWriteArraySet                           ║
║ Concurrent sorted set             ║ ConcurrentSkipListSet                         ║
║ Concurrent FIFO queue             ║ ConcurrentLinkedQueue (unbounded, lock-free)  ║
║ Bounded blocking queue            ║ ArrayBlockingQueue (single lock)              ║
║ High-throughput blocking queue    ║ LinkedBlockingQueue (two locks: put + take)   ║
║ Direct handoff (no buffering)     ║ SynchronousQueue                              ║
║ Delayed task scheduling           ║ DelayQueue                                    ║
║ Priority-based blocking queue     ║ PriorityBlockingQueue                         ║
║ Transfer queue                    ║ LinkedTransferQueue                           ║
╚═══════════════════════════════════╩═══════════════════════════════════════════════╝
```

---

### 6.4 Lock-Free & CAS Techniques

```
CAS (Compare-And-Swap):
  ┌──────────────────────────────────────────────────────────────┐
  │  Atomic hardware instruction:                                │
  │  compareAndSwap(memoryLocation, expectedValue, newValue)     │
  │                                                              │
  │  If current value == expected:                               │
  │    → Set to newValue, return true                            │
  │  Else:                                                       │
  │    → Do nothing, return false (another thread won)           │
  │    → Caller retries in a loop (spin)                         │
  │                                                              │
  │  No locking needed — hence "lock-free"                       │
  │  Used in: ConcurrentHashMap (empty bins), AtomicInteger,     │
  │           ConcurrentLinkedQueue, LongAdder                   │
  └──────────────────────────────────────────────────────────────┘
```

```java
// CAS in ConcurrentHashMap — inserting into empty bin:
// Pseudo-code from OpenJDK:
if (tabAt(tab, i) == null) {
    if (casTabAt(tab, i, null, new Node<>(hash, key, value))) {
        break;  // CAS succeeded — node installed without any lock
    }
    // CAS failed — another thread installed something → retry loop
}

// ConcurrentLinkedQueue — lock-free offer():
// Uses CAS to atomically update tail.next pointer
// Multiple threads can offer() simultaneously — CAS resolves contention
// No thread ever blocks — they just retry

// AtomicInteger — lock-free increment:
AtomicInteger counter = new AtomicInteger(0);
counter.incrementAndGet();  // internally:
//   do { current = get(); } while (!compareAndSet(current, current + 1));
```

---

### 6.5 Copy-on-Write Trade-offs

```
Writes:                              Reads:
  ┌─────────────┐                      ┌─────────────┐
  │ acquire lock │                      │ NO lock     │
  │ copy array   │                      │ volatile    │
  │ modify copy  │                      │   read of   │
  │ volatile     │                      │   array ref │
  │   write ref  │                      │ direct      │
  │ release lock │                      │   index     │
  └─────────────┘                      └─────────────┘
  Cost: O(N)                            Cost: O(1)
  Memory: 2× during write              Memory: stable
  Frequency: rare                       Frequency: dominant

Trade-offs:
  ✓ Readers NEVER block — highest possible read throughput
  ✓ Snapshot iteration — iterate without worrying about modifications
  ✓ No ConcurrentModificationException
  ✗ Each write copies entire array → O(N) per write
  ✗ During write: 2× memory (old array + new array alive simultaneously)
  ✗ Write latency proportional to collection size
  ✗ Stale reads during iteration (snapshot = stale data)

Rule of thumb:
  Read/Write ratio > 100:1 → Copy-on-Write
  Write-heavy → ConcurrentHashMap / ConcurrentLinkedQueue
```

---

### 6.6 False Sharing Concerns

```
Cache line: CPU fetches memory in 64-byte cache lines (typical x86)

False sharing:
  ┌────────────────────────────── Cache Line (64 bytes) ──────────────────────────────┐
  │  counterA (Thread 1 writes)  │  counterB (Thread 2 writes)  │  padding...         │
  └──────────────────────────────┴──────────────────────────────┴─────────────────────┘
  
  Thread 1 writes counterA → invalidates ENTIRE cache line for Thread 2
  Thread 2 writes counterB → invalidates ENTIRE cache line for Thread 1
  → Both cores continuously invalidate each other's cache → performance collapse

Java mitigation:
  @jdk.internal.vm.annotation.Contended  // pads fields to avoid false sharing
  
  LongAdder uses this:
    Cell[] cells;  // each Cell is padded to its own cache line
    // Thread 1 increments cells[0], Thread 2 increments cells[1]
    // Different cache lines → no false sharing → near-linear scalability
    // sum() = base + Σ(cells[i])

  ConcurrentHashMap.CounterCell uses @Contended for its size-tracking cells

In collection design:
  - ConcurrentHashMap: volatile reads of table array → sequential consistency penalty
  - LinkedBlockingQueue: two separate locks (putLock, takeLock) on different fields
    → reduces false sharing between producer and consumer hot paths
```

---

### 6.7 Contention & Scalability

```
Scalability under contention (threads: 1 → 64):

  Throughput (ops/sec, relative)

  100% ┤
       │  ★ ConcurrentHashMap ────────── scales ~linearly with readers
   80% ┤  ★ ConcurrentSkipListMap ────── scales well for sorted needs
       │  ★ CopyOnWriteArrayList ─────── near-perfect read scaling
   60% ┤
       │
   40% ┤
       │  ▲ LinkedBlockingQueue ──────── decent (2-lock scheme)
   20% ┤  ▲ ArrayBlockingQueue ──────── moderate (single lock)
       │
    0% ┤  ■ synchronizedMap ─────────── collapses under contention
       │  ■ Hashtable ───────────────── collapses under contention
       └──┬──────┬──────┬──────┬──────
          1      8     16     32    64  threads

Key insight:
  - Global lock (Hashtable, synchronizedMap): throughput DECREASES with more threads
    (serialization + context switching overhead)
  - Fine-grained lock (ConcurrentHashMap): throughput INCREASES with more threads
    (independent locks on different bins allow true parallelism)
  - Lock-free (ConcurrentLinkedQueue): best theoretical scalability
    (no locks at all, just CAS retries under contention)
```

---

## Theory F: CPU Cache, Data Locality & Why It Matters for Collections

> The performance difference between ArrayList and LinkedList isn't just about Big-O — it's largely about how CPU caches work. This section explains why.

### The Memory Hierarchy

```
  Speed & Size (approximate, modern hardware):
  
  ┌──────────────┐
  │  CPU Register  │  ~0.3 ns    │  ~1 KB       │  Fastest, compiler-managed
  ├──────────────┤
  │  L1 Cache      │  ~1 ns      │  32-64 KB    │  Per-core, split I/D
  ├──────────────┤
  │  L2 Cache      │  ~4 ns      │  256 KB-1 MB │  Per-core
  ├──────────────┤
  │  L3 Cache      │  ~12 ns     │  8-64 MB     │  Shared across cores
  ├──────────────┤
  │  Main RAM      │  ~60-100 ns │  8-512 GB    │  DRAM
  ├──────────────┤
  │  SSD/NVMe      │  ~10-100 µs │  256 GB-8 TB │  1000× slower than RAM
  ├──────────────┤
  │  HDD           │  ~5-10 ms   │  1-20 TB     │  100,000× slower than RAM
  └──────────────┘

  Key insight: L1 cache hit is 60-100× faster than a main memory access!
```

### Cache Lines — The Unit of Transfer

```
  CPU never reads a single byte from RAM. It reads a CACHE LINE (typically 64 bytes).

  When you access array[5], the CPU loads:
    array[0..7] (if each element is 8 bytes) into one cache line

  Next access to array[6]?
    → Already in cache! FREE! ("cache hit")

  Access linked list node.next?
    → Node could be ANYWHERE in memory
    → Likely NOT in cache ("cache miss")
    → Must wait ~100 ns for RAM fetch
```

### Spatial Locality vs Temporal Locality

```
  Spatial Locality: Nearby memory addresses are likely to be accessed soon
    ✓ Arrays (elements are contiguous) → perfect spatial locality
    ✗ Linked structures (nodes scattered in heap) → poor spatial locality

  Temporal Locality: Recently accessed data is likely accessed again
    ✓ Loop variables, counters, hot fields
    ✗ One-time-use data

  Collections ranked by spatial locality:
  
  Best  ▐███████████▌  int[]          (primitive, contiguous)
        ▐██████████ ▌  ArrayList      (Object[], contiguous refs)
        ▐█████████  ▌  ArrayDeque     (Object[], contiguous refs)
        ▐███████    ▌  HashMap        (array of bins, but nodes scattered)
        ▐█████      ▌  TreeMap        (node-based, scattered)
        ▐███        ▌  LinkedHashMap  (node + linked list, scattered)
  Worst ▐██         ▌  LinkedList     (nodes scattered everywhere)
```

### Real-World Impact: ArrayList vs LinkedList Iteration

```
  Iterating 1 million elements:

  ArrayList (contiguous array):
    [elem0][elem1][elem2][elem3][elem4][elem5][elem6][elem7]...
    └─────────── one cache line ──────────┘
    → Read 1 cache line, get 8 elements for free
    → Sequential memory access → CPU prefetcher predicts next reads
    → Total cache misses: ~125,000 (1M refs × 8 bytes / 64 bytes per line)

  LinkedList (scattered nodes):
    [node0] ─► [node7832] ─► [node291] ─► [node50123] ─► ...
    ↑           ↑             ↑            ↑
    random      random        random       random memory locations
    → Almost every node access = cache miss
    → CPU prefetcher can't predict next address
    → Total cache misses: ~1,000,000 (worst case, one per node)

  Result: ArrayList iteration is 5-10× FASTER than LinkedList
  even though both are theoretically O(N)!
```

### Prefetching — Why Sequential Access is King

```
  Modern CPUs have a hardware prefetcher that detects sequential access patterns:

  Access pattern: addr, addr+64, addr+128, addr+192, ...
  Prefetcher: "I see a stride-64 pattern! I'll pre-load the NEXT cache lines
              before you ask for them."

  Result: For sequential array traversal, data is already in L1/L2
  cache by the time you need it → effective latency ≈ 0!

  For pointer chasing (linked structures):
  Each node.next is a random address → prefetcher CANNOT predict
  → Every access waits for full RAM latency (~100 ns)
```

### False Sharing — Cache Line Contention

```
  Two variables on the SAME cache line, accessed by DIFFERENT threads:

  Cache Line (64 bytes):
  ┌────────────────────────────────────────────────┐
  │  counterA (Thread 1)  │  counterB (Thread 2)          │
  └────────────────────────────────────────────────┘

  Thread 1 writes counterA → ENTIRE cache line invalidated in Thread 2's cache
  Thread 2 writes counterB → ENTIRE cache line invalidated in Thread 1's cache
  → Both threads constantly invalidating each other's cache = SLOW!

  Solution: Padding to ensure each variable gets its own cache line
  Java 8+: @Contended annotation (JDK internal)
  ConcurrentHashMap uses CounterCell[] with @Contended to avoid this
```

### Practical Guidelines for Collection Performance

```
  Rule 1: Prefer array-backed over node-based structures
    ArrayList > LinkedList (always, for iteration)
    ArrayDeque > LinkedList (always, for stack/queue)
    HashMap with good distribution > TreeMap (unless ordering needed)

  Rule 2: Access data sequentially when possible
    for (int i = 0; i < list.size(); i++) → sequential access ✓
    for (int i = list.size()-1; i >= 0; i--) → still good (prefetcher handles backward)
    random.get(random.nextInt(size)) → random access, poor locality ✗

  Rule 3: Keep hot data small and close together
    Struct-of-Arrays > Array-of-Structs for batch processing
    Smaller objects = more per cache line = better throughput

  Rule 4: Beware of pointer chasing
    HashMap lookup: array[bucket] → node.next → node.next = 2-3 cache misses
    Open-addressing maps: array[bucket] → array[bucket+1] = 0-1 cache misses
```

---

## 7. Performance & Memory Considerations

### 7.1 Big-O Comparison — All Major Implementations

```
╔═══════════════════╦═══════╦═════════╦══════════╦════════╦═══════════╦════════════╗
║ Operation          ║ Array ║ Linked  ║ HashMap  ║ TreeMap║ PriorityQ ║ ArrayDeque ║
║                    ║ List  ║ List    ║ /HashSet ║/TreeSet║ (heap)    ║            ║
╠═══════════════════╬═══════╬═════════╬══════════╬════════╬═══════════╬════════════╣
║ get(index)        ║ O(1)  ║ O(N)    ║  N/A     ║  N/A   ║  N/A      ║  N/A       ║
║ get(key)          ║ N/A   ║ N/A     ║ O(1)*    ║O(log N)║  N/A      ║  N/A       ║
║ add / put         ║ O(1)† ║ O(1)‡   ║ O(1)*    ║O(log N)║ O(log N)  ║ O(1)†      ║
║ add(index) / mid  ║ O(N)  ║ O(N)    ║  N/A     ║  N/A   ║  N/A      ║  N/A       ║
║ remove(index)     ║ O(N)  ║ O(N)    ║  N/A     ║  N/A   ║  N/A      ║  N/A       ║
║ remove(key/obj)   ║ O(N)  ║ O(N)    ║ O(1)*    ║O(log N)║ O(N)      ║ O(N)       ║
║ contains          ║ O(N)  ║ O(N)    ║ O(1)*    ║O(log N)║ O(N)      ║ O(N)       ║
║ peek / min / max  ║ N/A   ║ N/A     ║  N/A     ║O(log N)║ O(1)      ║ O(1)       ║
║ poll / dequeue    ║ N/A   ║ O(1)    ║  N/A     ║O(log N)║ O(log N)  ║ O(1)       ║
║ Iterator.next()   ║ O(1)  ║ O(1)    ║ O(1) avg ║ O(1)   ║ O(1)      ║ O(1)       ║
╠═══════════════════╬═══════╬═════════╬══════════╬════════╬═══════════╬════════════╣
║ Iteration total   ║ O(N)  ║ O(N)    ║ O(cap)§  ║ O(N)   ║ O(N)      ║ O(N)       ║
║ sort()            ║O(Nlog)║O(Nlog)  ║  N/A     ║ sorted ║  N/A      ║  N/A       ║
╚═══════════════════╩═══════╩═════════╩══════════╩════════╩═══════════╩════════════╝

*  O(1) amortized, O(log N) worst case (treeified bin post-Java 8)
†  O(1) amortized, O(N) worst case during resize
‡  O(1) at head/tail, O(N) for traversal to middle
§  HashSet/HashMap: iteration cost = O(capacity), not O(size)
   → If capacity = 10000 but size = 10, iterator still visits all empty buckets
   → LinkedHashMap/LinkedHashSet: O(size) for iteration (linked list bypass)
```

---

### 7.2 Memory Overhead Per Data Structure

```
Bytes per element (approximate, 64-bit JVM, compressed oops):

╔══════════════════════════╦═════════════════╦══════════════════════════════════╗
║ Data Structure            ║ Bytes/Element   ║ Breakdown                        ║
╠══════════════════════════╬═════════════════╬══════════════════════════════════╣
║ int[] (primitive array)  ║      4 bytes    ║ raw value, no overhead           ║
║ Integer[] (boxed)        ║     20 bytes    ║ 16 (header) + 4 (int value)      ║
║ ArrayList<Integer>       ║   ~28 bytes     ║ 8 (ref in array) + 20 (Integer)  ║
║ LinkedList<Integer>      ║   ~60 bytes     ║ 40 (Node) + 20 (Integer)         ║
╠══════════════════════════╬═════════════════╬══════════════════════════════════╣
║ HashSet<String>          ║   ~80 bytes     ║ 48 (Node) + ~32 (avg String)     ║
║ LinkedHashSet<String>    ║   ~96 bytes     ║ 64 (Entry) + ~32 (avg String)    ║
║ TreeSet<String>          ║   ~88 bytes     ║ 56 (TreeEntry) + ~32 (String)    ║
║ EnumSet                  ║    ~1 bit       ║ bit vector (1 long for ≤64 enums)║
╠══════════════════════════╬═════════════════╬══════════════════════════════════╣
║ HashMap<String,String>   ║   ~120 bytes    ║ 48 (Node) + 32 (key) + 32 (val) ║
║ TreeMap<String,String>   ║   ~128 bytes    ║ 56 (Entry) + 32 (key) + 32 (val)║
║ ConcurrentHashMap        ║   ~120 bytes    ║ similar to HashMap per entry     ║
╠══════════════════════════╬═════════════════╬══════════════════════════════════╣
║ PriorityQueue<Integer>   ║   ~28 bytes     ║ 8 (ref in array) + 20 (Integer)  ║
║ ArrayDeque<Integer>      ║   ~28 bytes     ║ 8 (ref in array) + 20 (Integer)  ║
╚══════════════════════════╩═════════════════╩══════════════════════════════════╝

Plus fixed overhead per collection:
  ArrayList:       ~32 bytes (object header + elementData ref + size + modCount)
  HashMap:         ~48 bytes (table ref + size + threshold + loadFactor + modCount)
  TreeMap:         ~48 bytes (root + comparator + size + modCount)
  ConcurrentHM:    ~64 bytes (+ CounterCell[] grows with contention)
```

---

### 7.3 Cache Locality Impact

```
Modern CPU cache hierarchy:
  L1 cache:  ~32 KB,  ~1 ns    (per core)
  L2 cache:  ~256 KB, ~4 ns    (per core)
  L3 cache:  ~8 MB,   ~12 ns   (shared across cores)
  Main RAM:            ~100 ns  (60-100× slower than L1)

Array-backed structures (ArrayList, ArrayDeque, PriorityQueue):
  ┌────┬────┬────┬────┬────┬────┬────┬────┐
  │ e0 │ e1 │ e2 │ e3 │ e4 │ e5 │ e6 │ e7 │  ← contiguous memory
  └────┴────┴────┴────┴────┴────┴────┴────┘
  → CPU prefetcher predicts sequential access → elements pre-loaded into cache
  → Iteration: mostly L1/L2 cache hits → FAST

Node-based structures (LinkedList, TreeMap, HashMap chains):
  ┌──────┐      ┌──────┐      ┌──────┐
  │ Node │ --→  │ Node │ --→  │ Node │
  │ @0x1 │      │ @0x9 │      │ @0x5 │  ← scattered across heap
  └──────┘      └──────┘      └──────┘
  → Each node access = potential cache miss → main memory fetch (~100 ns)
  → Iteration: dominated by cache misses → SLOW

Practical impact:
  ArrayList sequential iteration:  ~1 billion elements/sec
  LinkedList sequential iteration: ~100 million elements/sec
  → 10× difference purely from cache behavior (same O(N) complexity!)

HashMap iteration:
  - Visits bucket array sequentially → cache-friendly for array
  - Following chains/trees → cache-unfriendly for nodes
  - LinkedHashMap iteration: follows linked list → slightly worse cache behavior
    but O(size) instead of O(capacity)
```

---

### 7.4 GC Impact

```
How collections affect garbage collection:

1. Object count:
   - LinkedList (1M entries): 1M Node objects + 1M element objects = 2M objects
   - ArrayList (1M entries):  1 array + 1M element objects = ~1M objects
   - GC must trace every live object → more objects = longer GC pauses

2. Object graph depth:
   - ArrayList: root → array → elements (depth 2)
   - LinkedList: root → node → node → node → ... (depth N)
   - TreeMap: root → entry → entry → entry (depth log N)
   - Deep graphs increase GC mark phase time (stack depth for tracing)

3. Promotion patterns:
   - Short-lived collections: allocated in Young Gen, collected quickly (minor GC)
   - Long-lived collections: promoted to Old Gen → full GC needed to collect
   - Large arrays: may be allocated directly in Old Gen (humongous allocation in G1)

4. HashMap resize GC pressure:
   - Resize creates new Node[] (2× old capacity)
   - Old array becomes garbage immediately after resize
   - If old array is in Old Gen → triggers GC card table updates
   - For large maps (>100K entries), resize can cause GC spikes

5. WeakHashMap + GC:
   - Keys are WeakReferences → collected if no strong refs
   - expungeStaleEntries() runs on EVERY map operation
   - GC behavior directly affects map contents → non-deterministic size()

Mitigation strategies:
  ✓ Pre-size collections: new HashMap<>(expectedSize * 4 / 3 + 1)
  ✓ Use primitive-specialized collections (Eclipse Collections, HPPC, Koloboke)
  ✓ Reuse collections instead of creating new ones
  ✓ Use off-heap storage for very large datasets (Chronicle Map, MapDB)
  ✓ Trim after bulk load: arrayList.trimToSize(), or copy to exact-sized array
```

---

### 7.5 Tuning HashMap Capacity & Load Factor

```java
// Scenario: Building a lookup map for 10,000 known items

// ❌ Default — causes ~15 resize operations:
Map<String, Config> map = new HashMap<>();          // capacity 16, load 0.75
// 12, 24, 48, 96, 192, 384, 768, 1536, 3072, 6144, 12288 → resizes at each threshold

// ✅ Pre-sized — zero resize operations:
Map<String, Config> map = new HashMap<>(10000 * 4 / 3 + 1);
// capacity = 13334 → rounds up to 16384 (next power of 2)
// threshold = 16384 × 0.75 = 12288 > 10000 → no resize ever

// ✅ Guava — handles the math:
Map<String, Config> map = Maps.newHashMapWithExpectedSize(10000);

// Load factor trade-offs:
//   loadFactor = 0.5:  more memory, fewer collisions, faster lookups
//   loadFactor = 0.75: default — good balance
//   loadFactor = 1.0:  less memory, more collisions, slower lookups
//   loadFactor = 2.0:  very compact, many collisions → use TreeMap instead

// For KNOWN STATIC DATA (lookup table, constants):
Map<String, String> lookup = Map.of("K1", "V1", "K2", "V2", ...);
// Map.of uses optimized internal implementations:
//   0-2 entries: field-based (no array at all)
//   3+ entries:  hash table with probe sequence
//   Immutable → no resize ever, minimal memory
```

---

### 7.6 Large-Scale System Considerations

```
When dealing with millions of entries:

╔══════════════════════════════════════╦══════════════════════════════════════════╗
║ Concern                              ║ Mitigation                               ║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ OOM from large HashMap              ║ Use off-heap: Chronicle Map, MapDB       ║
║                                      ║ Or partition across JVMs (Redis, Hazelcast)║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ HashMap resize spike                ║ Pre-size aggressively                    ║
║ (copies millions of entries)        ║ Or use ConcurrentHashMap (incremental    ║
║                                      ║   resize via forwarding nodes)           ║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ GC pressure from boxed primitives   ║ Eclipse Collections: IntObjectHashMap    ║
║                                      ║ Koloboke: HashIntIntMap                  ║
║                                      ║ HPPC: IntIntHashMap (no boxing)          ║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ Sorting millions of elements        ║ Arrays.parallelSort() (fork-join merge)  ║
║                                      ║ External merge sort for disk-based data  ║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ Memory layout inefficiency          ║ Consider array-of-structs vs struct-of-  ║
║ (pointers to objects everywhere)    ║   arrays layout for cache efficiency     ║
║                                      ║ Project Valhalla (value types) will fix  ║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ Concurrent access at scale          ║ ConcurrentHashMap (not sync wrapper)     ║
║                                      ║ Striped locks for custom structures      ║
║                                      ║ Thread-local accumulation + periodic merge║
╠══════════════════════════════════════╬══════════════════════════════════════════╣
║ Iteration over huge collections     ║ Stream.parallel() for CPU-bound          ║
║                                      ║ Batched processing (limit/skip)          ║
║                                      ║ Cursor/pagination for external data      ║
╚══════════════════════════════════════╩══════════════════════════════════════════╝
```

---

## Theory G: Collection Design Patterns & Best Practices

> This section covers recurring patterns, anti-patterns, and idiomatic Java practices for working with collections in production code.

### Pattern 1: Choosing the Right Collection (Decision Flowchart)

```
  Need a collection?
  │
  ├─ Need key-value pairs?
  │   ├─ Need sorted keys?          → TreeMap
  │   ├─ Need insertion order?      → LinkedHashMap
  │   ├─ Need thread-safety?
  │   │   ├─ High concurrency?      → ConcurrentHashMap
  │   │   └─ Low concurrency?       → Collections.synchronizedMap()
  │   ├─ Need weak keys (GC)?      → WeakHashMap
  │   └─ Default (fastest)          → HashMap
  │
  ├─ Need unique elements only?
  │   ├─ Need sorted?               → TreeSet
  │   ├─ Need insertion order?      → LinkedHashSet
  │   ├─ Enum values?               → EnumSet
  │   └─ Default (fastest)          → HashSet
  │
  ├─ Need ordered / indexed access?
  │   ├─ Frequent random access?    → ArrayList
  │   ├─ Frequent insert/remove at head? → ArrayDeque
  │   ├─ Read-heavy, shared?        → CopyOnWriteArrayList
  │   └─ Default                    → ArrayList
  │
  └─ Need FIFO / LIFO / Priority?
      ├─ FIFO queue?                → ArrayDeque (or LinkedBlockingQueue)
      ├─ LIFO stack?                → ArrayDeque (NOT Stack class)
      ├─ Priority ordering?         → PriorityQueue
      ├─ Producer-consumer?         → BlockingQueue (ArrayBQ or LinkedBQ)
      └─ Timed scheduling?          → DelayQueue
```

### Pattern 2: Defensive Copying

```java
// PROBLEM: Returning internal collection exposes mutable state
public class UserService {
    private final List<User> users = new ArrayList<>();
    
    // ❌ BAD: Caller can modify our internal list!
    public List<User> getUsers() { return users; }
    
    // ✅ GOOD: Return unmodifiable view
    public List<User> getUsers() { return Collections.unmodifiableList(users); }
    
    // ✅ BETTER (Java 10+): Return immutable copy
    public List<User> getUsers() { return List.copyOf(users); }
    
    // ✅ BEST for API: Return stream (lazy, no copy needed)
    public Stream<User> users() { return users.stream(); }
}

// PROBLEM: Accepting mutable collection as parameter
// ❌ Caller can modify after passing!
public void setUsers(List<User> users) { this.users = users; }

// ✅ GOOD: Defensive copy on input
public void setUsers(List<User> users) { this.users = new ArrayList<>(users); }
```

### Pattern 3: Multimap, Multiset, and Computed Collections

```java
// PROBLEM: Map<K, List<V>> is tedious to manage
// ❌ Verbose and error-prone:
Map<String, List<Order>> ordersByCustomer = new HashMap<>();
ordersByCustomer.computeIfAbsent("alice", k -> new ArrayList<>()).add(order);

// ✅ Java 8+ computeIfAbsent makes it tolerable:
map.computeIfAbsent(key, k -> new ArrayList<>()).add(value);

// ✅ Guava Multimap (cleaner API):
ListMultimap<String, Order> orders = ArrayListMultimap.create();
orders.put("alice", order1);
orders.put("alice", order2);
List<Order> aliceOrders = orders.get("alice");  // [order1, order2]

// COUNTING PATTERN (frequency map):
// ✅ Java 8+ merge:
Map<String, Integer> freq = new HashMap<>();
for (String word : words) {
    freq.merge(word, 1, Integer::sum);
}

// ✅ Streams (more idiomatic):
Map<String, Long> freq = Arrays.stream(words)
    .collect(Collectors.groupingBy(Function.identity(), Collectors.counting()));
```

### Pattern 4: The Builder Pattern for Complex Collections

```java
// Building complex immutable structures:
// ❌ Mutable then wrap:
Map<String, List<String>> map = new HashMap<>();
map.put("fruits", Arrays.asList("apple", "banana"));
map.put("vegs", Arrays.asList("carrot", "pea"));
Map<String, List<String>> immutable = Collections.unmodifiableMap(map);

// ✅ Java 9+ factory methods:
Map<String, List<String>> immutable = Map.of(
    "fruits", List.of("apple", "banana"),
    "vegs", List.of("carrot", "pea")
);

// ✅ Guava ImmutableMap builder for larger maps:
ImmutableMap<String, List<String>> map = ImmutableMap.<String, List<String>>builder()
    .put("fruits", ImmutableList.of("apple", "banana"))
    .put("vegs", ImmutableList.of("carrot", "pea"))
    .buildOrThrow();
```

### Anti-Pattern Gallery

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `new LinkedList<>()` for general use | Poor cache locality, 40 bytes/node | Use `ArrayList` or `ArrayDeque` |
| `new Vector<>()` or `new Stack<>()` | Synchronized on every operation, legacy API | Use `ArrayList` + explicit sync, or `ArrayDeque` |
| `Collections.synchronizedMap(new HashMap<>())` at scale | Global lock, TOCTOU races | Use `ConcurrentHashMap` |
| `HashMap` with mutable keys | Key moves to different bucket after mutation → lost | Use immutable keys or `record` types |
| `new HashMap<>()` for known-size data | Multiple unnecessary resizes | `new HashMap<>(expectedSize * 4/3 + 1)` |
| Iterating `Map.keySet()` then calling `get()` | Two lookups per entry: O(2N) hash operations | Use `entrySet()` → single lookup per entry |
| `list.contains()` in a hot loop | O(N) per call → O(N²) total | Convert to `HashSet` for O(1) lookup |
| Returning `null` for empty collections | Forces callers to null-check everywhere | Return `Collections.emptyList()` or `List.of()` |

### The `Collections` Utility Class — Key Methods

```java
// Sorting and searching:
Collections.sort(list);                      // natural order
Collections.sort(list, comparator);          // custom order
Collections.binarySearch(list, key);         // O(log N) — list MUST be sorted!

// Thread-safe wrappers:
Collections.synchronizedList(list);          // wraps with synchronized
Collections.synchronizedMap(map);

// Immutable empty collections (singleton, no allocation):
Collections.emptyList();                     // same empty instance every time
Collections.emptySet();
Collections.emptyMap();

// Singleton collections:
Collections.singletonList(item);             // fixed-size list with 1 element
Collections.singleton(item);                 // fixed-size set with 1 element

// Utility:
Collections.frequency(collection, element);  // count occurrences
Collections.disjoint(c1, c2);               // true if no common elements
Collections.unmodifiableList(list);          // read-only view
Collections.reverse(list);                   // in-place reverse
Collections.shuffle(list);                   // random permutation
```

---

## 8. Advanced Topics

### 8.1 Immutable Collections

```
╔══════════════════════════════╦═══════════════════════════╦═══════════════════════════════╗
║ Method                        ║ Returns                   ║ Notes                         ║
╠══════════════════════════════╬═══════════════════════════╬═══════════════════════════════╣
║ Collections.unmodifiable      ║ Unmodifiable VIEW         ║ Backed by original — changes  ║
║   List/Set/Map(original)     ║ of the original           ║   to original reflect here!   ║
╠══════════════════════════════╬═══════════════════════════╬═══════════════════════════════╣
║ List.of() / Set.of()        ║ Truly IMMUTABLE copy      ║ No null elements allowed.     ║
║ Map.of()  (Java 9+)         ║ (not backed by anything)  ║ Set.of: no duplicates         ║
╠══════════════════════════════╬═══════════════════════════╬═══════════════════════════════╣
║ List.copyOf() / Set.copyOf()║ Immutable copy of input   ║ Returns same reference if     ║
║ Map.copyOf()  (Java 10+)    ║                           ║ input is already immutable    ║
╠══════════════════════════════╬═══════════════════════════╬═══════════════════════════════╣
║ Collectors.toUnmodifiable    ║ Immutable from stream     ║ toUnmodifiableList(),         ║
║   List/Set/Map (Java 10+)   ║   pipeline                ║ toUnmodifiableSet(),          ║
║                               ║                           ║ toUnmodifiableMap()           ║
╠══════════════════════════════╬═══════════════════════════╬═══════════════════════════════╣
║ Stream.toList() (Java 16+)  ║ Unmodifiable list         ║ Shorthand for                 ║
║                               ║                           ║ collect(toUnmodifiableList()) ║
╚══════════════════════════════╩═══════════════════════════╩═══════════════════════════════╝
```

```java
// ⚠️ CRITICAL DIFFERENCE: unmodifiable VIEW vs immutable COPY

List<String> original = new ArrayList<>(List.of("A", "B", "C"));

// Unmodifiable VIEW — original changes bleed through:
List<String> view = Collections.unmodifiableList(original);
original.add("D");
System.out.println(view); // [A, B, C, D] ← view changed!

// Immutable COPY — completely detached:
List<String> copy = List.copyOf(original);
original.add("E");
System.out.println(copy); // [A, B, C, D] ← unaffected

// Stream.toList() — also creates unmodifiable copy:
List<String> streamList = original.stream().filter(s -> s.length() > 0).toList();
// streamList.add("X"); → UnsupportedOperationException
```

#### Internal Implementation Optimization (Java 9+)

```java
// List.of() uses specialized implementations based on size:
List.of()                    // → List0 (empty singleton, no fields)
List.of("A")                 // → List1 (single field, no array)
List.of("A", "B")            // → List2 (two fields, no array)
List.of("A", "B", "C", ...)  // → ListN (Object[] array — minimal memory)

// Set.of() uses a probe-sequence hash table (not HashMap!):
// → Smaller memory footprint than HashSet
// → Iteration order is randomized per-run (security against DoS)

// Map.of() similarly uses specialized MapN with interleaved key-value array:
// → Object[] = {key0, val0, key1, val1, ...}
// → Much more cache-friendly than HashMap's Node objects
```

---

### 8.2 Spliterator

```
Spliterator<T>: "Splittable Iterator" — designed for parallel processing
  ┌─────────────────────────────────────────────────────────────────┐
  │  boolean tryAdvance(Consumer<? super T> action)                 │
  │     → Process one element (like Iterator.next + hasNext)        │
  │                                                                 │
  │  Spliterator<T> trySplit()                                      │
  │     → Split into two halves for parallel processing             │
  │     → Returns spliterator for first half, this covers second    │
  │     → Returns null if can't split further                       │
  │                                                                 │
  │  long estimateSize()                                            │
  │     → Remaining element count estimate                          │
  │                                                                 │
  │  int characteristics()                                          │
  │     → Bit flags describing the source (ORDERED, DISTINCT,       │
  │       SORTED, SIZED, NONNULL, IMMUTABLE, CONCURRENT, SUBSIZED) │
  └─────────────────────────────────────────────────────────────────┘

How parallel streams use Spliterator:

  Original collection: [1, 2, 3, 4, 5, 6, 7, 8]
                              │
                       trySplit()
                       ┌──────┴──────┐
                  [1,2,3,4]      [5,6,7,8]
                    │                 │
              trySplit()         trySplit()
              ┌────┴────┐     ┌────┴────┐
           [1,2]     [3,4]  [5,6]    [7,8]
             │         │      │         │
          Thread1   Thread2  Thread3  Thread4
             │         │      │         │
             ▼         ▼      ▼         ▼
           process   process process  process
             │         │      │         │
             └─────────┴──────┴─────────┘
                        │
                     combine
```

```java
// Characteristics matter for optimization:
// SIZED + SUBSIZED → parallel framework knows exact partition boundaries
//   ArrayList: SIZED, ORDERED, SUBSIZED → excellent parallel split
//   HashSet:   SIZED, DISTINCT → decent parallel split
//   TreeSet:   SORTED, ORDERED, SIZED, DISTINCT → sorted parallel split
//   LinkedList: ORDERED, SIZED → poor parallel split (sequential traversal)

// Custom Spliterator for a database ResultSet:
class ResultSetSpliterator implements Spliterator<Row> {
    private final ResultSet rs;
    private int remaining;
    
    @Override
    public boolean tryAdvance(Consumer<? super Row> action) {
        if (rs.next()) { action.accept(mapRow(rs)); remaining--; return true; }
        return false;
    }
    
    @Override
    public Spliterator<Row> trySplit() { return null; } // can't split a ResultSet
    
    @Override
    public long estimateSize() { return remaining; }
    
    @Override
    public int characteristics() { return ORDERED | NONNULL; }
}
```

---

### 8.3 Parallel Stream Performance Caveats

```
╔════════════════════════════════╦══════════════════════════════════════════════╗
║ Use parallel WHEN              ║ AVOID parallel WHEN                          ║
╠════════════════════════════════╬══════════════════════════════════════════════╣
║ CPU-intensive computation     ║ I/O-bound operations (DB, HTTP, file)        ║
║ Large datasets (>10K elements)║ Small datasets (<10K elements)               ║
║ Stateless operations          ║ Stateful operations (shared mutable state)   ║
║ Source: ArrayList, arrays     ║ Source: LinkedList, Stream.iterate()          ║
║   (good Spliterator split)   ║   (poor split → sequential bottleneck)       ║
║ Independent element processing║ Order-dependent processing                   ║
║ No contention in reduction    ║ Synchronized/contended collectors            ║
╚════════════════════════════════╩══════════════════════════════════════════════╝
```

```java
// ❌ DANGER: parallel() on I/O operations → blocks ForkJoinPool.commonPool()
products.parallelStream()
    .map(p -> httpClient.fetchPrice(p))  // blocks common pool threads!
    .toList();
// Common pool: default = Runtime.getRuntime().availableProcessors() - 1 threads
// ALL parallel streams in your JVM share this pool → one slow I/O starves everyone

// ✅ FIX 1: Use custom ForkJoinPool for parallel I/O:
ForkJoinPool customPool = new ForkJoinPool(20);
List<Price> prices = customPool.submit(() -> 
    products.parallelStream().map(p -> httpClient.fetchPrice(p)).toList()
).join();

// ✅ FIX 2: Use CompletableFuture + custom executor (better for I/O):
ExecutorService io = Executors.newFixedThreadPool(20);
List<CompletableFuture<Price>> futures = products.stream()
    .map(p -> CompletableFuture.supplyAsync(() -> httpClient.fetchPrice(p), io))
    .toList();
List<Price> prices = futures.stream().map(CompletableFuture::join).toList();

// ❌ DANGER: accumulation into shared mutable state
List<String> results = new ArrayList<>(); // NOT thread-safe!
stream.parallel().forEach(s -> results.add(s.toUpperCase())); // DATA RACE!

// ✅ FIX: Use collect (thread-safe accumulation built-in):
List<String> results = stream.parallel()
    .map(String::toUpperCase)
    .collect(Collectors.toList()); // or .toList()
```

#### Spliterator Quality by Collection Type

```
Parallel efficiency (best → worst):

1. Arrays / ArrayList          → BEST: random-access split → perfect halves
   Spliterator: SIZED, SUBSIZED, ORDERED
   Split: O(1) — just divide index range

2. HashSet / HashMap.keySet()  → GOOD: bucket-based splitting
   Spliterator: SIZED, DISTINCT
   Split: divide bucket range

3. TreeSet / TreeMap            → MODERATE: requires tree traversal
   Spliterator: SORTED, ORDERED, SIZED, DISTINCT
   Split: balanced but involves tree navigation

4. LinkedList                   → POOR: must traverse to find midpoint
   Spliterator: SIZED, ORDERED
   Split: O(N) traversal to split point → sequential bottleneck

5. Stream.iterate()            → WORST: inherently sequential
   Spliterator: ORDERED only
   Split: effectively can't split → parallel brings zero benefit
```

---

### 8.4 Custom Collection Implementation Guidelines

```java
// When might you need a custom collection?
// - Domain-specific invariants (e.g., unique by a field, max-size, sorted by multiple criteria)
// - Performance-critical hot path (avoid boxing, ensure locality)
// - Missing from JDK (e.g., MultiMap, BiMap, CircularBuffer)

// APPROACH 1: Extend AbstractList/AbstractSet/AbstractMap (recommended)
// → Only need to implement a few abstract methods
// → get(index) + size() for AbstractList → all other methods derived

public class ImmutableRangeList extends AbstractList<Integer> {
    private final int start, end;  // inclusive start, exclusive end
    
    public ImmutableRangeList(int start, int end) {
        this.start = start;
        this.end = end;
    }
    
    @Override public Integer get(int index) {
        Objects.checkIndex(index, size());
        return start + index;
    }
    
    @Override public int size() { return end - start; }
    
    @Override public boolean contains(Object o) {
        if (o instanceof Integer i) return i >= start && i < end;
        return false;
    }
    // get(), size() → gives you: iterator(), indexOf(), subList(), stream(), etc. FOR FREE
}

// APPROACH 2: Composition over inheritance
public class BoundedSet<E> implements Set<E> {
    private final Set<E> delegate;
    private final int maxSize;
    
    public BoundedSet(int maxSize) {
        this.delegate = new LinkedHashSet<>();
        this.maxSize = maxSize;
    }
    
    @Override public boolean add(E e) {
        if (delegate.size() >= maxSize) 
            throw new IllegalStateException("Max size " + maxSize + " reached");
        return delegate.add(e);
    }
    
    // Delegate all other methods to `delegate`
    @Override public int size() { return delegate.size(); }
    @Override public boolean contains(Object o) { return delegate.contains(o); }
    // ... etc.
}

// APPROACH 3: Guava's ForwardingCollection wrappers
// → Less boilerplate than manual delegation
public class LoggingMap<K,V> extends ForwardingMap<K,V> {
    private final Map<K,V> delegate = new HashMap<>();
    @Override protected Map<K,V> delegate() { return delegate; }
    
    @Override public V put(K key, V value) {
        log.debug("put({}, {})", key, value);
        return super.put(key, value);
    }
}
```

---

### 8.5 Java 21+ SequencedCollection

```java
// Java 21 introduced SequencedCollection, SequencedSet, SequencedMap
// Unified access to first/last elements across ordered collections

// Before Java 21 — inconsistent API:
list.get(0);                    list.get(list.size() - 1);
deque.getFirst();               deque.getLast();
sortedSet.first();              sortedSet.last();
linkedHashSet.iterator().next(); // no easy way to get last!

// After Java 21 — unified:
SequencedCollection<E>:
  E getFirst();
  E getLast();
  void addFirst(E);
  void addLast(E);
  E removeFirst();
  E removeLast();
  SequencedCollection<E> reversed();  // reverse-order view

// All of these now implement SequencedCollection:
// List, Deque, LinkedHashSet, SortedSet, SequencedSet

// Similarly, SequencedMap provides:
SequencedMap<K,V>:
  Map.Entry<K,V> firstEntry();
  Map.Entry<K,V> lastEntry();
  Map.Entry<K,V> pollFirstEntry();
  Map.Entry<K,V> pollLastEntry();
  SequencedMap<K,V> reversed();

// Implemented by: LinkedHashMap, TreeMap, ConcurrentSkipListMap
```

## 9. Staff-Level / Senior-Level Interview Questions

### 9A. Deep Internal Questions

---

#### Q1: Why is HashMap capacity always a power of 2?

**Expected Senior Answer:**

Bucket index is computed as `hash & (capacity - 1)` — a bitwise AND. When capacity is a power of 2, `capacity - 1` produces a bitmask of all 1s in the lower bits (e.g., 16 → `0b01111`), giving uniform distribution across all buckets.

If capacity were not a power of 2 (e.g., 15 → `0b01110`), bit 0 would always be 0, meaning odd-numbered buckets are never used — effectively halving the table and doubling collision rates.

**Additional benefits:**
- **Resize optimization:** During doubling (N → 2N), each entry either stays at index `i` or moves to `i + oldCapacity`. Determined by checking a single bit: `hash & oldCapacity`. No full rehash needed — just one bit check per entry.
- **Compiler optimization:** `hash & (cap - 1)` is a single machine instruction, whereas `hash % cap` requires integer division (much slower).

**Common mistakes:**
- Saying "modulo is used for bucket index" — it's bitwise AND, not modulo.
- Not knowing about the resize optimization.

**Follow-up probes:**
- "If you set initial capacity to 13, what does HashMap actually do?" → Rounds up to 16 (`tableSizeFor()` finds next power of 2).
- "How does the hash spreading function `hash(key)` relate to this?"

---

#### Q2: Explain the treeification thresholds — why 8 and 6?

**Expected Senior Answer:**

Under normal operation with load factor 0.75, the probability of a single bucket accumulating K entries follows a Poisson distribution:

```
P(k) = (λ^k × e^-λ) / k!   where λ ≈ 0.5 (expected entries per bucket)

P(0) = 0.60653    P(4) = 0.00155
P(1) = 0.30327    P(5) = 0.00016
P(2) = 0.07582    P(6) = 0.00001
P(3) = 0.01264    P(7) = 0.000001
                   P(8) = 0.00000006
```

The probability of 8+ entries in one bucket is ~6 × 10⁻⁸ — nearly impossible with good hash functions. If it happens, it indicates either:
1. **Pathological/adversarial keys** (hash DoS attack)
2. **Broken hashCode()** implementation

Converting to a Red-Black Tree limits worst-case from O(N) to O(log N), providing attack resistance.

**Why untreeify at 6 (not 8)?**
Hysteresis — a gap of 2 prevents thrashing. If both thresholds were 8, repeatedly adding and removing the 8th element would cause continuous conversion between list and tree.

**Why not treeify until capacity ≥ 64 (`MIN_TREEIFY_CAPACITY`)?**
With small tables, collisions are better resolved by resizing (doubling capacity to spread entries) rather than treeifying a single bin. Treeification has higher memory overhead (~104 bytes per TreeNode vs ~48 bytes per Node).

**Follow-up probes:**
- "The keys need to implement `Comparable` for treeification — what if they don't?"
  → HashMap uses `System.identityHashCode()` as tiebreaker for tree ordering.
- "What's the amortized cost of treeification?"

---

#### Q3: How does ConcurrentHashMap avoid global locks?

**Expected Senior Answer:**

**Java 8+ architecture (current):**

1. **Empty bin insertion → CAS (lock-free):**
   - `tabAt(tab, i)` reads bucket using volatile semantics.
   - If null, `casTabAt(tab, i, null, newNode)` atomically installs the node.
   - No lock is ever acquired — pure CAS retry loop.

2. **Non-empty bin → synchronized on first node:**
   - `synchronized (f)` where `f` is the first node in the bucket.
   - Only the specific bin is locked — other bins are unaffected.
   - This is equivalent to lock striping with N stripes (one per bucket).

3. **Read operations → NO locks at all:**
   - `get()` uses volatile reads of `table` and `Node.val`/`Node.next`.
   - Java Memory Model ensures happens-before relationship via volatile.

4. **Size tracking → LongAdder-style CAS cells:**
   - `baseCount` for low-contention updates.
   - `CounterCell[]` array when CAS on `baseCount` fails.
   - Each thread uses a different cell → spreads contention.

5. **Resize → incremental transfer:**
   - Uses `ForwardingNode` sentinels.
   - Multiple threads can cooperatively transfer buckets during resize.
   - Readers encountering `ForwardingNode` are redirected to the new table.

**Common mistakes:**
- Describing the old pre-Java 8 Segment architecture (16 fixed segments with ReentrantLock).
- Saying "reads use CAS" — reads don't use CAS, they use volatile reads.
- Thinking `computeIfAbsent` is lock-free — it's not; it synchronizes on the bin.

**Follow-up probes:**
- "What's the concurrency difference between `putIfAbsent` and `computeIfAbsent`?"
  → `putIfAbsent` evaluates value before lock; `computeIfAbsent` runs function under lock.
- "Can you deadlock with ConcurrentHashMap?" → Yes, if computeIfAbsent's function accesses the same CHM.

---

#### Q4: What happens during HashMap `resize()`?

**Expected Senior Answer:**

```
1. Allocate new Node[] of 2× capacity.
2. For each bucket in old table:
   a. Single node → recalculate index: hash & (newCap - 1), place in new table.
   b. Linked list → split into two lists using (hash & oldCap) bit:
      - 0 → "lo" list, stays at index i
      - 1 → "hi" list, moves to index i + oldCap
      Preserves relative order (important for consistent iteration in LinkedHashMap).
   c. Tree bin → split similarly; if resulting subtree size < UNTREEIFY_THRESHOLD (6),
      convert back to linked list.
3. Replace table reference.
4. Old array becomes garbage (GC will collect).

Cost: O(N) — visits every entry.
Memory: briefly 2× (old + new table coexist).
Latency impact: single put() that triggers resize pays the ENTIRE O(N) cost.
```

**Critical insight for production:**
For a HashMap with 10M entries, resize copies 10M nodes. At ~50ns per node (memory access + pointer update), that's ~500ms of stop-the-world processing on the thread that triggered resize. In a latency-sensitive system, this is catastrophic.

**Mitigation:**
- Pre-size: `new HashMap<>((int)(expectedSize / 0.75) + 1)`
- Use ConcurrentHashMap even in single-threaded code if resize latency matters — it does incremental resize across multiple operations.

---

#### Q5: Memory visibility in concurrent collections — how is it guaranteed?

**Expected Senior Answer:**

Java Memory Model (JMM) guarantees:

**ConcurrentHashMap:**
- `table` array is `volatile` → reading the table reference establishes a happens-before edge.
- Node fields `val` and `next` are `volatile` → writes by putter are visible to getter.
- Synchronized bin access → monitor release happens-before monitor acquire.
- Combination provides: "a successful `put()` happens-before a subsequent `get()` that sees the value."

**CopyOnWriteArrayList:**
- Internal array is `volatile` → writing new array reference (after copy + modification) happens-before reading it.
- This is why reads need zero synchronization — the volatile write of the array reference is the publication fence.

**BlockingQueue:**
- `put()` happens-before corresponding `take()` receives the element.
- Guaranteed by the internal `ReentrantLock` release (put) → acquire (take) chain.

**Common mistake:**
- Thinking `synchronized` only provides mutual exclusion. It also establishes memory visibility — everything done inside a synchronized block is visible to the next thread entering a synchronized block on the same monitor.

---

#### Q6: HashMap vs Hashtable vs ConcurrentHashMap — internals comparison

**Expected Senior Answer:**

| Aspect | HashMap | Hashtable | ConcurrentHashMap |
|---|---|---|---|
| Lock | None | `synchronized` on every method | CAS for empty bins, `synchronized` per bin |
| Null key | 1 allowed (bucket 0) | Not allowed (NPE) | Not allowed (ambiguity) |
| Null value | Allowed | Not allowed | Not allowed |
| Hash function | `(h = hashCode()) ^ (h >>> 16)` | `(hashCode() & 0x7FFFFFFF) % capacity` | `(h ^ (h >>> 16)) & 0x7FFFFFFF` |
| Capacity | Always power of 2 | Any odd number preferred | Power of 2 |
| Resize | 2× | 2× + 1 | 2×, incremental (ForwardingNode) |
| Treeification | Yes (Java 8+) | No | Yes (Java 8+) |
| Iterator | Fail-fast | Fail-fast (Enumerator is not) | Weakly consistent |
| Legacy | No | Yes (Java 1.0) | No (Java 1.5+) |

**Why `null` not allowed in ConcurrentHashMap:**
`map.get(key)` returns `null` → ambiguous: key absent or value is null? In single-threaded HashMap, use `containsKey()` to disambiguate. In concurrent setting, another thread may change the state between `get()` and `containsKey()` → TOCTOU race → unsafe.

---

### 9B. Performance & Debugging Scenarios

---

#### Q7: Diagnose high CPU caused by hash collisions

**Scenario:** Production web app's CPU spikes to 100% on POST requests processing user-supplied data. Thread dump shows all threads stuck in `HashMap.get()`.

**Expected Senior Answer:**

**Root Cause Analysis:**
1. Application uses user-supplied strings (e.g., JSON keys, form parameters) as HashMap keys.
2. Attacker crafts keys with identical `hashCode()` values → all entries fall into one bucket.
3. Pre-Java 8: bucket is a linked list → `get()` degrades to O(N) per lookup.
4. With thousands of colliding keys: O(N²) total → CPU saturated.

**Diagnosis steps:**
```java
// 1. Thread dump — look for:
"http-thread-123" RUNNABLE
  java.util.HashMap.getNode(HashMap.java:574)
  java.util.HashMap.get(HashMap.java:556)
  com.app.MyServlet.doPost(MyServlet.java:42)
// Multiple threads stuck in HashMap.getNode → collision chain traversal

// 2. Verify: Check bucket distribution
// (In debugging/test environment)
HashMap<String, Object> map = ...; // the problematic map
Field tableField = HashMap.class.getDeclaredField("table");
tableField.setAccessible(true);
Object[] table = (Object[]) tableField.get(map);
for (int i = 0; i < table.length; i++) {
    int chainLen = 0;
    for (Object node = table[i]; node != null; ) {
        chainLen++;
        Field next = node.getClass().getDeclaredField("next");
        next.setAccessible(true);
        node = next.get(node);
    }
    if (chainLen > 10) System.out.println("Bucket " + i + ": " + chainLen + " entries");
}
```

**Fixes:**
1. **Upgrade to Java 8+** → treeification limits worst case to O(N log N).
2. **Limit request parameters** at WAF/framework level (max 100 params).
3. **Use ConcurrentHashMap** which has additional hash spreading.
4. **Randomize hash seed** (JDK already does mixing via `hash()` function).
5. **Don't use user-controlled strings as HashMap keys** without sanitization.

**Follow-up:** "How would the behavior differ on Java 8 vs Java 7?" → Java 8 treeifies at 8 entries per bin, converting O(N) chain to O(log N) tree. The attack is mitigated but not eliminated (O(N log N) total instead of O(N²)).

---

#### Q8: Debug ConcurrentModificationException in production

**Scenario:** Intermittent `ConcurrentModificationException` in a microservice. Stack trace points to a `forEach` over an `ArrayList`. The code has no explicit threading.

**Expected Senior Answer:**

**Diagnosis approach:**

```java
// Stack trace:
java.util.ConcurrentModificationException
  at java.util.ArrayList$Itr.checkForComodification(ArrayList.java:911)
  at java.util.ArrayList$Itr.next(ArrayList.java:861)
  at com.app.OrderProcessor.processOrders(OrderProcessor.java:55)

// Common causes:
// 1. Modification during enhanced for-loop:
for (Order order : orders) {
    if (order.isExpired()) orders.remove(order); // CME!
}

// 2. Callback/listener modifying the collection being iterated:
for (Order order : orders) {
    orderService.process(order); // internally calls listener that adds to `orders`
}

// 3. Another thread modifies (even if not "explicit" threading):
//    - Spring @Async method
//    - CompletableFuture callback
//    - Timer/ScheduledExecutor task
//    - Servlet thread reusing shared collection
```

**Fixes (ordered by recommendation):**

```java
// FIX 1: Use removeIf (Java 8+):
orders.removeIf(Order::isExpired);

// FIX 2: Use Iterator.remove():
Iterator<Order> it = orders.iterator();
while (it.hasNext()) {
    if (it.next().isExpired()) it.remove();
}

// FIX 3: Collect-then-modify:
List<Order> expired = orders.stream().filter(Order::isExpired).toList();
orders.removeAll(expired);

// FIX 4: If concurrent access: use concurrent collection:
List<Order> orders = new CopyOnWriteArrayList<>(); // if read-heavy
// OR: ConcurrentLinkedDeque / synchronized block if write-heavy

// FIX 5: If concurrent modification by callback:
// Iterate over a snapshot:
for (Order order : List.copyOf(orders)) {  // snapshot — safe even if orders is modified
    orderService.process(order);
}
```

**Key insight:** `ConcurrentModificationException` is NOT about threads — it's about structural modification during iteration. Single-threaded code triggers it too.

---

#### Q9: Production memory bloat from collections

**Scenario:** Heap dump shows 2 GB occupied by `HashMap$Node` objects in a service with only 500K logical entries. Expected usage: ~200 MB.

**Expected Senior Answer:**

**Root cause investigation:**

```
Step 1: Analyze heap dump (MAT / VisualVM)
  - Dominator tree: Which HashMap holds the most nodes?
  - Check retained size of each collection

Common causes:
  a) HashMap never shrinks — entries removed but capacity stays at peak:
     HashMap grew to 1M entries, now has 500K, but capacity is still 2M
     → 2M × 48 bytes (Node[]) = ~96 MB just for the empty bucket array
     Fix: Periodically recreate: new HashMap<>(currentMap)

  b) Boxed primitives:
     HashMap<Integer, Long> with 500K entries:
     → 500K × 48 (Node) + 500K × 16 (Integer key) + 500K × 24 (Long value)
     = 44 MB nodes + 8 MB keys + 12 MB values = 64 MB
     Fix: Use primitive-specialized map (Eclipse Collections IntLongHashMap ≈ 12 MB)

  c) Excessive load factor / capacity:
     new HashMap<>(4_000_000, 0.5f) → 8M buckets × 8 bytes = 64 MB just for array
     Even with 500K entries: 7.5M empty slots

  d) Memory leaks:
     Entries never removed (session data, caches without eviction)
     WeakHashMap with strong references in values pointing back to keys

  e) String duplication in keys:
     500K distinct strings × 50 chars avg = ~50 MB in String objects
     Fix: String.intern() for high-overlap keys, or use symbol table
```

**Production fix pattern:**
```java
// Periodic compaction for long-lived maps:
if (map.size() < map.capacity() / 4) {
    map = new HashMap<>(map);  // rebuild with right-sized capacity
}

// Better: Use Caffeine with maximumSize for caches:
Cache<String, UserData> cache = Caffeine.newBuilder()
    .maximumSize(100_000)
    .expireAfterAccess(Duration.ofMinutes(30))
    .build();
```

---

#### Q10: Design and implement an LRU cache in Java

**Expected Senior Answer:**

```java
// APPROACH 1: LinkedHashMap (simplest — sufficient for most interviews)
public class LRUCache<K, V> extends LinkedHashMap<K, V> {
    private final int capacity;
    
    public LRUCache(int capacity) {
        super(capacity + 1, 0.75f, true); // accessOrder = true
        this.capacity = capacity;
    }
    
    @Override
    protected boolean removeEldestEntry(Map.Entry<K, V> eldest) {
        return size() > capacity;
    }
}

// APPROACH 2: ConcurrentHashMap + ConcurrentLinkedDeque (thread-safe, O(1))
public class ConcurrentLRUCache<K, V> {
    private final int capacity;
    private final ConcurrentHashMap<K, V> map;
    private final ConcurrentLinkedDeque<K> order;
    
    public ConcurrentLRUCache(int capacity) {
        this.capacity = capacity;
        this.map = new ConcurrentHashMap<>(capacity);
        this.order = new ConcurrentLinkedDeque<>();
    }
    
    public V get(K key) {
        V value = map.get(key);
        if (value != null) {
            order.remove(key);    // O(N) — weakness of this approach
            order.addLast(key);
        }
        return value;
    }
    
    public void put(K key, V value) {
        if (map.containsKey(key)) {
            order.remove(key);
        } else if (map.size() >= capacity) {
            K evicted = order.pollFirst();
            if (evicted != null) map.remove(evicted);
        }
        map.put(key, value);
        order.addLast(key);
    }
}

// APPROACH 3: Production — use Caffeine (battle-tested LRU/LFU cache)
Cache<Key, Value> cache = Caffeine.newBuilder()
    .maximumSize(10_000)
    .expireAfterWrite(Duration.ofMinutes(5))
    .recordStats()
    .build();
// Uses Window TinyLfu eviction — better hit rate than LRU
// Lock-free reads, O(1) amortized writes
```

**Follow-up probes:**
- "What eviction policy is better than LRU?" → LFU, ARC, W-TinyLfu (Caffeine's algorithm)
- "How would you add TTL-based expiration?" → `ScheduledExecutorService` + DelayQueue, or use Caffeine
- "How to handle cache stampede?" → `computeIfAbsent` (blocks concurrent loads for same key) or `AsyncCache`

---

### 9C. System Design Angle

---

#### Q11: Design an in-memory caching layer with eviction

**Expected Senior Answer:**

```
Requirements clarification:
  - Expected key-value pairs: ~10M entries
  - Read/write ratio: ~100:1
  - Eviction: LRU or TTL-based
  - Thread-safe: high concurrent reads
  - Metrics: hit rate, eviction count

Architecture:
  ┌─────────────────────────────────────────────────────────┐
  │                  CacheManager                            │
  │                                                         │
  │  ┌──────────────────────────────────────────────────┐   │
  │  │  Partitioned Cache (16-256 segments)              │   │
  │  │                                                    │   │
  │  │  Segment[0]: ConcurrentHashMap + LRU eviction     │   │
  │  │  Segment[1]: ConcurrentHashMap + LRU eviction     │   │
  │  │  ...                                               │   │
  │  │  Segment[N]: ConcurrentHashMap + LRU eviction     │   │
  │  │                                                    │   │
  │  │  Key → Segment mapping: hash(key) & (N-1)        │   │
  │  └──────────────────────────────────────────────────┘   │
  │                                                         │
  │  ┌─────────────┐  ┌────────────────┐  ┌──────────────┐ │
  │  │  Expiration  │  │  Size Monitor  │  │   Metrics    │ │
  │  │  Scheduler   │  │  (background)  │  │  (hit/miss)  │ │
  │  └─────────────┘  └────────────────┘  └──────────────┘ │
  └─────────────────────────────────────────────────────────┘
```

**Collection choices and rationale:**

| Component | Collection | Why |
|---|---|---|
| Per-segment store | `ConcurrentHashMap` | Lock-free reads, per-bin lock writes |
| LRU tracking | `ConcurrentLinkedDeque<K>` or access-time field per entry | Track access order |
| Expiration | `DelayQueue<Expirable>` or background `ScheduledExecutor` | Time-based eviction |
| Metrics counters | `LongAdder` | Contention-free counting |

```java
// Simplified segment-based cache:
public class SegmentedCache<K, V> {
    private static final int SEGMENTS = 16;
    private final Segment<K, V>[] segments;
    
    @SuppressWarnings("unchecked")
    public SegmentedCache(int maxSizePerSegment) {
        segments = new Segment[SEGMENTS];
        for (int i = 0; i < SEGMENTS; i++) {
            segments[i] = new Segment<>(maxSizePerSegment);
        }
    }
    
    private Segment<K, V> segmentFor(K key) {
        int hash = key.hashCode();
        return segments[(hash ^ (hash >>> 16)) & (SEGMENTS - 1)];
    }
    
    public V get(K key) { return segmentFor(key).get(key); }
    public void put(K key, V value) { segmentFor(key).put(key, value); }
    
    static class Segment<K, V> {
        private final LinkedHashMap<K, V> map; // accessOrder=true for LRU
        private final int maxSize;
        
        Segment(int maxSize) {
            this.maxSize = maxSize;
            this.map = new LinkedHashMap<>(maxSize, 0.75f, true) {
                @Override protected boolean removeEldestEntry(Map.Entry<K, V> e) {
                    return size() > maxSize;
                }
            };
        }
        
        synchronized V get(K key) { return map.get(key); }
        synchronized void put(K key, V value) { map.put(key, value); }
    }
}

// Production answer: Use Caffeine
// → Window TinyLfu eviction (near-optimal hit rate)
// → Lock-free reads, O(1) amortized writes
// → Built-in TTL, refresh, stats, async loading
```

**Follow-up:** "How would you add distributed eviction notifications?" → JMS/Kafka events on eviction, or use Hazelcast/Redis pub-sub for cache coherence.

---

#### Q12: Handling millions of keys — collection choice and tuning

**Scenario:** Real-time analytics system processes 50M events/day. Need to maintain rolling counters per user (10M users), with fast increment and periodic aggregation.

**Expected Senior Answer:**

```
Analysis:
  - 10M users × counter object → memory is critical
  - High write rate → contention on counters
  - Periodic read for aggregation → bulk scan

Option 1: ConcurrentHashMap<UserId, LongAdder>
  ✅ Lock-free reads and writes
  ✅ LongAdder: contention-free increments (CAS cells)
  ❌ 10M entries × ~120 bytes/entry + 10M × ~72 bytes (LongAdder) ≈ 1.8 GB
  ❌ LongAdder per user is wasteful for low-frequency users

Option 2: ConcurrentHashMap<UserId, AtomicLong>
  ✅ Lower memory: ~120 + 24 bytes per entry ≈ 1.4 GB
  ✅ Atomic increment, lock-free
  ❌ AtomicLong under contention → CAS retries
  ✅ Good enough if per-user contention is low (likely, with 10M users)

Option 3: Primitive-specialized map (Eclipse Collections, Koloboke)
  If UserId is a long:
  ✅ LongLongHashMap: ~24 bytes per entry → 10M × 24 = 240 MB
  ✅ No boxing overhead
  ❌ Not thread-safe by default → need synchronized access or segments

Option 4: Off-heap storage (Chronicle Map)
  ✅ Minimal GC impact (off-heap memory, not traced by GC)
  ✅ Memory-mapped → survives restarts
  ✅ Can handle 100M+ entries
  ❌ Serialization overhead, more complex API

Recommended architecture:
  ┌──────────────────────────────────────────────┐
  │  Per-thread: ThreadLocal<HashMap<Long, long>>│  ← batch events locally
  │  Periodically: merge into central counters    │
  │  Central: ConcurrentHashMap<Long, AtomicLong> │  ← or off-heap
  │  Aggregation: entrySet().parallelStream()    │
  └──────────────────────────────────────────────┘
  
  Thread-local accumulation → periodic merge drastically reduces contention
  10M per-user merges per minute vs 50M individual increments per day
```

**Follow-up:** "How would you persist these counters?" → Periodic snapshot to RocksDB or Redis. Or use Chronicle Map with persistence.

---

#### Q13: Avoiding GC pressure with collections

**Expected Senior Answer:**

```
Major GC pressure sources from collections:

1. Object creation:
   - Boxing: int → Integer creates objects → GC
   - Node/Entry allocation: HashMap.put() creates Node objects → GC
   - Resize: ArrayList/HashMap creates new arrays → old ones become garbage

2. Object retention:
   - Long-lived collections in Old Gen → only Full GC can reclaim
   - References from Old Gen to Young Gen → card table writes → GC overhead

3. Object graph complexity:
   - LinkedList: N nodes → GC traces N objects
   - HashMap: N entries + N nodes + bucket array → complex graph

Strategies:

a) Avoid boxing (biggest single win):
   ❌ HashMap<Integer, Integer> → 3 objects per entry (Node, Integer key, Integer value)
   ✅ IntIntHashMap (Eclipse) → zero objects per entry (primitive arrays)

b) Object pooling for entries:
   ❌ Create 10M Entry objects → 10M allocations
   ✅ Array-based storage: keys[] + values[] parallel arrays

c) Pre-size + trimToSize:
   ❌ ArrayList default → 25+ resize → 25 garbage arrays
   ✅ new ArrayList<>(expected) → 0 resizes → 0 garbage

d) Off-heap storage:
   ✅ Chronicle Map, MapDB → data stored outside JVM heap
   ✅ DirectByteBuffer for manual management
   ❌ More complex, must manage serialization

e) Reuse collections instead of creating:
   ❌ In hot loop: List<T> result = new ArrayList<>(); → GC per call
   ✅ ThreadLocal<ArrayList<T>> with .clear() between uses

f) Use array-backed structures over node-based:
   ✅ ArrayList over LinkedList (1 array vs N nodes)
   ✅ ArrayDeque over LinkedList (1 array vs N nodes)
   ✅ Open-addressing hash map vs chaining hash map (fewer objects)
```

```java
// Before: Every getUserOrders() creates 3 garbage objects per call
public List<Order> getUserOrders(String userId) {
    List<Order> orders = new ArrayList<>();  // new object
    // ... fill from DB
    return orders;
}

// After: Reuse with ThreadLocal (hot path optimization)
private static final ThreadLocal<ArrayList<Order>> BUFFER = 
    ThreadLocal.withInitial(ArrayList::new);

public List<Order> getUserOrders(String userId) {
    ArrayList<Order> buffer = BUFFER.get();
    buffer.clear();
    // ... fill from DB into buffer
    return List.copyOf(buffer);  // immutable snapshot, buffer reused
}
```

---

#### Q14: Collection choice for a high-frequency trading (HFT) system

**Expected Senior Answer:**

```
HFT constraints:
  - Latency: <10 μs per operation (microseconds, not milliseconds)
  - Throughput: millions of messages per second
  - GC: ZERO GC pauses (or <1ms if unavoidable)
  - Determinism: consistent latency, no outliers

Collection choices:
  ╔════════════════════════════╦══════════════════════════════════════════╗
  ║ Use Case                    ║ Collection Choice                        ║
  ╠════════════════════════════╬══════════════════════════════════════════╣
  ║ Order book (price→qty)     ║ Primitive-specialized TreeMap or         ║
  ║                              ║ manually sorted array (fixed size)      ║
  ║                              ║ NEVER java.util.TreeMap (object alloc)  ║
  ╠════════════════════════════╬══════════════════════════════════════════╣
  ║ Symbol lookup table        ║ Open-addressing hash map with primitive  ║
  ║                              ║ keys (Koloboke/HPPC) or perfect hashing ║
  ╠════════════════════════════╬══════════════════════════════════════════╣
  ║ Message queue              ║ Ring buffer (Disruptor pattern),         ║
  ║                              ║ NOT BlockingQueue (lock contention)     ║
  ╠════════════════════════════╬══════════════════════════════════════════╣
  ║ Connection pool            ║ Pre-allocated array, NOT ArrayList       ║
  ║                              ║ (control exact object count)            ║
  ╠════════════════════════════╬══════════════════════════════════════════╣
  ║ Event log                  ║ Off-heap circular buffer                 ║
  ║                              ║ (zero GC, memory-mapped)               ║
  ╚════════════════════════════╩══════════════════════════════════════════╝

Key principles:
  1. Pre-allocate everything at startup → zero runtime allocation
  2. Use primitive types → zero boxing → zero GC
  3. Use arrays over linked structures → cache-friendly
  4. Avoid locks → use CAS or single-threaded design
  5. Memory-map for persistence → avoid serialization overhead
  6. Use Epsilon GC (Java 11+) or Shenandoah/ZGC if allocation is unavoidable
```

---

#### Q15: Scaling a concurrent map across a distributed system

**Expected Senior Answer:**

```
Single-JVM ConcurrentHashMap → need to scale beyond one machine:

Architecture options:

  Option A: Embedded distributed cache (Hazelcast IMap, Apache Ignite)
  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
  │   JVM Node 1 │  │   JVM Node 2 │  │   JVM Node 3 │
  │  ┌─────────┐ │  │  ┌─────────┐ │  │  ┌─────────┐ │
  │  │ IMap    │ │──│──│ IMap    │ │──│──│ IMap    │ │
  │  │ shard A │ │  │  │ shard B │ │  │  │ shard C │ │
  │  └─────────┘ │  │  └─────────┘ │  │  └─────────┘ │
  └──────────────┘  └──────────────┘  └──────────────┘
  ✅ ConcurrentMap API → easy migration from CHM
  ✅ Automatic partitioning + replication
  ❌ Network latency on remote partitions
  ❌ Split-brain risk (need quorum config)

  Option B: External distributed cache (Redis, Memcached)
  ┌──────────┐     ┌──────────────────┐
  │  App JVM  │────►│  Redis Cluster   │
  │  (client) │     │  (sharded, HA)   │
  └──────────┘     └──────────────────┘
  ✅ Battle-tested at scale
  ✅ Persistence options (RDB, AOF)
  ❌ Serialization overhead (network roundtrip)
  ❌ Not a native Java collection interface

  Option C: Near-cache + distributed backend (hybrid)
  ┌──────────────────────────────────┐
  │  JVM (Caffeine L1 cache)         │
  │    ↓ miss                         │
  │  Redis / Hazelcast (L2 cache)    │
  │    ↓ miss                         │
  │  Database (source of truth)      │
  └──────────────────────────────────┘
  ✅ Sub-microsecond for hot data (L1 hit)
  ✅ Scales via distributed L2
  ❌ Cache coherence challenge (L1 may serve stale data)
  ❌ Complexity of invalidation

Partitioning strategy:
  - Consistent hashing → minimize redistribution on node add/remove
  - Virtual nodes → uniform load distribution
  - Partition key = hash(key) mod N partitions
  
Cache coherence for L1:
  - TTL-based (accept staleness up to N seconds)
  - Event-driven: Redis Pub/Sub or Kafka topic for invalidation events
  - Versioned reads: value includes version → check on use
```

**Follow-up:** 
- "How do you handle hot keys?" → Local aggregation, probabilistic early expiration, key splitting.
- "How to handle cache stampede in distributed setting?" → Distributed lock (Redisson), request coalescing, stale-while-revalidate.

---

> **End of Guide**  
> Sections 1–8 cover technical foundations. Section 9 maps directly to senior/staff interview expectations.  
> For each answer: demonstrate internals knowledge, production awareness, and trade-off analysis.
