# Collections Framework - Complete Guide for Interviews & Competitive Programming

## What is the Collections Framework?

The **Java Collections Framework** is a unified architecture for representing and manipulating collections of objects. It provides:
- **Interfaces** (abstract data types): List, Set, Queue, Map
- **Implementations** (concrete classes): ArrayList, HashSet, HashMap, etc.
- **Algorithms** (utility methods): sort, search, shuffle, etc.

**Why it matters:** 90% of interview and competitive programming problems use collections!

---

## Collections Hierarchy

```
                    Collection (Interface)
                         |
        +----------------+----------------+
        |                |                |
      List             Set              Queue
        |                |                |
   ArrayList        HashSet          LinkedList
   LinkedList       TreeSet          PriorityQueue
   Vector           LinkedHashSet    ArrayDeque
   Stack
   
   
                    Map (Interface) - NOT part of Collection
                         |
        +----------------+----------------+
        |                |                |
     HashMap          TreeMap        LinkedHashMap
     Hashtable        
```

---

## 1. List Interface - Ordered Collection (Allows Duplicates)

### Theory

- **Ordered:** Elements maintain insertion order
- **Indexed:** Access by position (0, 1, 2, ...)
- **Duplicates:** Allowed
- **Null:** Usually allowed

---

### ArrayList - Dynamic Array

**Internal Structure:** Resizable array

```java
import java.util.*;

public class ArrayListDemo {
    public static void main(String[] args) {
        // Creation
        ArrayList<Integer> list = new ArrayList<>();           // Default capacity: 10
        ArrayList<Integer> list2 = new ArrayList<>(20);        // Custom capacity
        ArrayList<Integer> list3 = new ArrayList<>(Arrays.asList(1, 2, 3));
        
        // Add elements - O(1) amortized
        list.add(10);                   // [10]
        list.add(20);                   // [10, 20]
        list.add(1, 15);                // [10, 15, 20] - O(n) for insertion
        
        // Access - O(1)
        int val = list.get(0);          // 10
        int last = list.get(list.size() - 1);
        
        // Update - O(1)
        list.set(0, 5);                 // [5, 15, 20]
        
        // Remove - O(n)
        list.remove(1);                 // Remove at index 1: [5, 20]
        list.remove(Integer.valueOf(20)); // Remove value 20: [5]
        
        // Search - O(n)
        boolean contains = list.contains(5);    // true
        int index = list.indexOf(5);            // 0
        
        // Size
        int size = list.size();
        boolean empty = list.isEmpty();
        
        // Iterate
        for (int num : list) {
            System.out.println(num);
        }
        
        // Convert to array
        Integer[] arr = list.toArray(new Integer[0]);
        
        // Sort
        Collections.sort(list);                          // Ascending
        Collections.sort(list, Collections.reverseOrder()); // Descending
        
        // Clear
        list.clear();
    }
}
```

**When to use:**
- ✅ Random access by index
- ✅ Iterating through elements
- ✅ Searching/sorting
- ❌ Frequent insertions/deletions in middle

---

### LinkedList - Doubly Linked List

**Internal Structure:** Chain of nodes with prev/next pointers

```java
import java.util.*;

public class LinkedListDemo {
    public static void main(String[] args) {
        LinkedList<Integer> list = new LinkedList<>();
        
        // Add at ends - O(1)
        list.addFirst(10);              // [10]
        list.addLast(20);               // [10, 20]
        list.add(30);                   // [10, 20, 30] - adds at end
        
        // Access - O(n)
        int first = list.getFirst();    // 10
        int last = list.getLast();      // 30
        int val = list.get(1);          // 20 - O(n)!
        
        // Remove from ends - O(1)
        list.removeFirst();             // [20, 30]
        list.removeLast();              // [20]
        
        // Use as Stack (LIFO)
        list.push(10);                  // [10, 20]
        int popped = list.pop();        // 10, list: [20]
        
        // Use as Queue (FIFO)
        list.offer(30);                 // [20, 30]
        int polled = list.poll();       // 20, list: [30]
    }
}
```

**When to use:**
- ✅ Frequent insertions/deletions at ends
- ✅ Implementing Stack/Queue/Deque
- ❌ Random access by index

---

### ArrayList vs LinkedList - Comparison

| Operation | ArrayList | LinkedList | Winner |
|-----------|-----------|------------|--------|
| **Get by index** | O(1) | O(n) | ArrayList ✅ |
| **Add at end** | O(1) amortized | O(1) | Tie ✅ |
| **Add at beginning** | O(n) | O(1) | LinkedList ✅ |
| **Add in middle** | O(n) | O(n) | Tie |
| **Remove from end** | O(1) | O(1) | Tie ✅ |
| **Remove from beginning** | O(n) | O(1) | LinkedList ✅ |
| **Search** | O(n) | O(n) | Tie |
| **Memory overhead** | Low | High (prev/next pointers) | ArrayList ✅ |
| **Cache locality** | Excellent | Poor | ArrayList ✅ |

**Interview Tip:** In 95% of cases, use ArrayList. LinkedList only when you need O(1) operations at both ends.

---

## 2. Set Interface - Unique Elements

### Theory

- **Unordered** (except TreeSet/LinkedHashSet)
- **No duplicates:** Duplicate elements rejected
- **No indexing:** Can't access by position
- **Null:** HashSet allows 1 null, TreeSet doesn't allow null

---

### HashSet - Unordered Unique Elements

**Internal Structure:** HashMap (uses keys only, values are dummy)

```java
import java.util.*;

public class HashSetDemo {
    public static void main(String[] args) {
        HashSet<Integer> set = new HashSet<>();
        
        // Add - O(1) average
        set.add(10);                    // true
        set.add(20);                    // true
        set.add(10);                    // false (duplicate)
        // Set: {10, 20} - no guaranteed order
        
        // Contains - O(1) average
        boolean has = set.contains(10); // true
        
        // Remove - O(1) average
        set.remove(10);                 // true
        set.remove(100);                // false (not present)
        
        // Size
        int size = set.size();          // 1
        
        // Iterate (order not guaranteed)
        for (int num : set) {
            System.out.println(num);
        }
        
        // Convert List to Set (remove duplicates)
        List<Integer> listWithDups = Arrays.asList(1, 2, 2, 3, 3, 3);
        Set<Integer> unique = new HashSet<>(listWithDups); // {1, 2, 3}
    }
}
```

---

### LinkedHashSet - Maintains Insertion Order

```java
import java.util.*;

public class LinkedHashSetDemo {
    public static void main(String[] args) {
        LinkedHashSet<Integer> set = new LinkedHashSet<>();
        
        set.add(30);
        set.add(10);
        set.add(20);
        
        // Iterates in insertion order: 30, 10, 20
        for (int num : set) {
            System.out.println(num);
        }
    }
}
```

**When to use:** Need unique elements AND maintain insertion order

---

### TreeSet - Sorted Unique Elements

**Internal Structure:** Red-Black Tree (self-balancing BST)

```java
import java.util.*;

public class TreeSetDemo {
    public static void main(String[] args) {
        TreeSet<Integer> set = new TreeSet<>();
        
        // Add - O(log n)
        set.add(30);
        set.add(10);
        set.add(20);
        set.add(40);
        // Stored as: [10, 20, 30, 40] - always sorted!
        
        // Contains - O(log n)
        boolean has = set.contains(20);
        
        // Remove - O(log n)
        set.remove(20);
        
        // Range operations
        int first = set.first();                // 10 - smallest
        int last = set.last();                  // 40 - largest
        
        // Ceiling/Floor
        Integer ceil = set.ceiling(25);         // 30 (smallest >= 25)
        Integer floor = set.floor(25);          // 20 (largest <= 25)
        
        // Higher/Lower
        Integer higher = set.higher(20);        // 30 (strictly greater)
        Integer lower = set.lower(30);          // 20 (strictly smaller)
        
        // Subset operations
        SortedSet<Integer> headSet = set.headSet(30);   // [10, 20]
        SortedSet<Integer> tailSet = set.tailSet(20);   // [20, 30, 40]
        SortedSet<Integer> subSet = set.subSet(10, 40); // [10, 20, 30]
        
        // Poll (remove and return)
        int pollFirst = set.pollFirst();        // 10
        int pollLast = set.pollLast();          // 40
    }
}
```

---

### Set Comparison Table

| Feature | HashSet | LinkedHashSet | TreeSet |
|---------|---------|---------------|---------|
| **Ordering** | No order | Insertion order | Sorted (natural/comparator) |
| **Add/Remove/Contains** | O(1) avg | O(1) avg | O(log n) |
| **Null elements** | 1 null allowed | 1 null allowed | ❌ No null |
| **Memory** | Low | Medium | High |
| **Performance** | Fastest | Medium | Slowest |
| **Use case** | General purpose | Order matters | Need sorted data |

---

## 3. Queue Interface - FIFO Processing

### Theory

- **FIFO:** First In, First Out
- **Operations:** offer (add), poll (remove), peek (view front)

---

### LinkedList as Queue

```java
import java.util.*;

public class QueueDemo {
    public static void main(String[] args) {
        Queue<Integer> queue = new LinkedList<>();
        
        // Add to rear - O(1)
        queue.offer(10);        // or add()
        queue.offer(20);
        queue.offer(30);
        // Queue: [10, 20, 30] (10 at front)
        
        // Peek at front - O(1)
        int front = queue.peek();       // 10 (doesn't remove)
        
        // Remove from front - O(1)
        int removed = queue.poll();     // 10, queue: [20, 30]
        
        // Size
        int size = queue.size();
        boolean empty = queue.isEmpty();
    }
}
```

**offer vs add:** offer returns false if full, add throws exception

---

### PriorityQueue - Min Heap (Priority Queue)

**Internal Structure:** Binary Heap (array-based)

```java
import java.util.*;

public class PriorityQueueDemo {
    public static void main(String[] args) {
        // Min Heap (default) - smallest element first
        PriorityQueue<Integer> minHeap = new PriorityQueue<>();
        
        minHeap.offer(30);
        minHeap.offer(10);
        minHeap.offer(20);
        
        System.out.println(minHeap.poll());  // 10 (smallest)
        System.out.println(minHeap.poll());  // 20
        System.out.println(minHeap.poll());  // 30
        
        // Max Heap - largest element first
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
        
        maxHeap.offer(30);
        maxHeap.offer(10);
        maxHeap.offer(20);
        
        System.out.println(maxHeap.poll());  // 30 (largest)
        System.out.println(maxHeap.poll());  // 20
        System.out.println(maxHeap.poll());  // 10
        
        // Custom comparator - sort by absolute value
        PriorityQueue<Integer> customPQ = new PriorityQueue<>((a, b) -> Math.abs(a) - Math.abs(b));
        
        customPQ.offer(-30);
        customPQ.offer(10);
        customPQ.offer(-20);
        
        System.out.println(customPQ.poll());  // 10 (|10| = 10 is smallest)
    }
}
```

**Time Complexity:**
- offer (insert): O(log n)
- poll (remove min/max): O(log n)
- peek: O(1)

**Common Interview Use Cases:**
- Top K elements
- Kth largest/smallest
- Merge K sorted lists
- Meeting room scheduling

---

### Deque (Double-Ended Queue)

```java
import java.util.*;

public class DequeDemo {
    public static void main(String[] args) {
        Deque<Integer> deque = new ArrayDeque<>();
        
        // Add to front
        deque.addFirst(20);     // [20]
        deque.addFirst(10);     // [10, 20]
        
        // Add to rear
        deque.addLast(30);      // [10, 20, 30]
        
        // Remove from front
        int front = deque.removeFirst();    // 10, [20, 30]
        
        // Remove from rear
        int rear = deque.removeLast();      // 30, [20]
        
        // Peek both ends
        int peekFront = deque.peekFirst();
        int peekRear = deque.peekLast();
        
        // Use as Stack
        deque.push(100);            // Add to front
        int popped = deque.pop();   // Remove from front
    }
}
```

**ArrayDeque vs LinkedList:**
- ArrayDeque: Faster, less memory, better cache locality ✅
- LinkedList: Slightly slower, more memory

---

## 4. Map Interface - Key-Value Pairs

### Theory

- **Key-Value pairs:** Each key maps to a value
- **Unique keys:** Duplicate keys override value
- **Not part of Collection hierarchy**

---

### HashMap - Unordered Key-Value

**Internal Structure:** Array of buckets (linked lists/trees)

```java
import java.util.*;

public class HashMapDemo {
    public static void main(String[] args) {
        HashMap<String, Integer> map = new HashMap<>();
        
        // Put - O(1) average
        map.put("Alice", 25);
        map.put("Bob", 30);
        map.put("Alice", 26);       // Overwrites Alice's value
        
        // Get - O(1) average
        int age = map.get("Alice");             // 26
        int def = map.getOrDefault("Charlie", 0); // 0 (not found)
        
        // Contains - O(1) average
        boolean hasKey = map.containsKey("Bob");    // true
        boolean hasValue = map.containsValue(30);   // true
        
        // Remove - O(1) average
        map.remove("Bob");
        
        // Size
        int size = map.size();
        boolean empty = map.isEmpty();
        
        // Iterate through keys
        for (String key : map.keySet()) {
            System.out.println(key + " = " + map.get(key));
        }
        
        // Iterate through entries (BETTER - single lookup)
        for (Map.Entry<String, Integer> entry : map.entrySet()) {
            System.out.println(entry.getKey() + " = " + entry.getValue());
        }
        
        // Get all keys/values
        Set<String> keys = map.keySet();
        Collection<Integer> values = map.values();
        
        // putIfAbsent - only put if key doesn't exist
        map.putIfAbsent("David", 35);
        
        // compute/merge - advanced operations
        map.merge("Alice", 1, Integer::sum);  // Alice's age + 1
    }
}
```

---

### Common HashMap Patterns for Interviews

#### Pattern 1: Frequency Counter

```java
public static Map<Character, Integer> charFrequency(String s) {
    Map<Character, Integer> freq = new HashMap<>();
    for (char c : s.toCharArray()) {
        freq.put(c, freq.getOrDefault(c, 0) + 1);
    }
    return freq;
}

// Example: "hello" → {h=1, e=1, l=2, o=1}
```

---

#### Pattern 2: Group Anagrams

```java
public List<List<String>> groupAnagrams(String[] strs) {
    Map<String, List<String>> map = new HashMap<>();
    
    for (String s : strs) {
        char[] chars = s.toCharArray();
        Arrays.sort(chars);
        String key = new String(chars);
        
        map.putIfAbsent(key, new ArrayList<>());
        map.get(key).add(s);
    }
    
    return new ArrayList<>(map.values());
}

// Example: ["eat","tea","tan","ate","nat","bat"]
// → [["eat","tea","ate"], ["tan","nat"], ["bat"]]
```

---

#### Pattern 3: Two Sum

```java
public int[] twoSum(int[] nums, int target) {
    Map<Integer, Integer> map = new HashMap<>();
    
    for (int i = 0; i < nums.length; i++) {
        int complement = target - nums[i];
        if (map.containsKey(complement)) {
            return new int[]{map.get(complement), i};
        }
        map.put(nums[i], i);
    }
    
    return new int[]{};
}
```

---

### LinkedHashMap - Maintains Insertion Order

```java
import java.util.*;

public class LinkedHashMapDemo {
    public static void main(String[] args) {
        LinkedHashMap<String, Integer> map = new LinkedHashMap<>();
        
        map.put("C", 3);
        map.put("A", 1);
        map.put("B", 2);
        
        // Iterates in insertion order: C, A, B
        for (String key : map.keySet()) {
            System.out.println(key + " = " + map.get(key));
        }
        
        // LRU Cache implementation
        LinkedHashMap<Integer, Integer> lruCache = new LinkedHashMap<>(16, 0.75f, true) {
            protected boolean removeEldestEntry(Map.Entry eldest) {
                return size() > 3;  // Max size 3
            }
        };
    }
}
```

---

### TreeMap - Sorted Keys

**Internal Structure:** Red-Black Tree

```java
import java.util.*;

public class TreeMapDemo {
    public static void main(String[] args) {
        TreeMap<Integer, String> map = new TreeMap<>();
        
        map.put(3, "three");
        map.put(1, "one");
        map.put(2, "two");
        map.put(5, "five");
        
        // Keys always sorted: {1=one, 2=two, 3=three, 5=five}
        
        // First/Last - O(log n)
        Map.Entry<Integer, String> first = map.firstEntry();   // 1=one
        Map.Entry<Integer, String> last = map.lastEntry();     // 5=five
        
        int firstKey = map.firstKey();      // 1
        int lastKey = map.lastKey();        // 5
        
        // Ceiling/Floor
        Integer ceilKey = map.ceilingKey(4);    // 5 (smallest key >= 4)
        Integer floorKey = map.floorKey(4);     // 3 (largest key <= 4)
        
        // Higher/Lower
        Integer higher = map.higherKey(2);      // 3
        Integer lower = map.lowerKey(3);        // 2
        
        // Range views
        SortedMap<Integer, String> headMap = map.headMap(3);    // keys < 3
        SortedMap<Integer, String> tailMap = map.tailMap(3);    // keys >= 3
        SortedMap<Integer, String> subMap = map.subMap(2, 5);   // 2 <= keys < 5
    }
}
```

---

### Map Comparison Table

| Feature | HashMap | LinkedHashMap | TreeMap | Hashtable |
|---------|---------|---------------|---------|-----------|
| **Ordering** | No order | Insertion order | Sorted keys | No order |
| **Get/Put/Remove** | O(1) avg | O(1) avg | O(log n) | O(1) avg |
| **Null key** | 1 allowed | 1 allowed | ❌ Not allowed | ❌ Not allowed |
| **Null values** | Allowed | Allowed | Allowed | ❌ Not allowed |
| **Thread-safe** | ❌ No | ❌ No | ❌ No | ✅ Yes (legacy) |
| **Performance** | Fastest | Medium | Slowest | Slow (synchronized) |
| **Use case** | General | Order matters | Sorted data | Legacy code only |

**Interview Note:** Never use Hashtable - use ConcurrentHashMap for thread-safety instead!

---

## Collections Utility Class

### Sorting

```java
import java.util.*;

public class CollectionsSort {
    public static void main(String[] args) {
        List<Integer> list = new ArrayList<>(Arrays.asList(5, 2, 8, 1, 9));
        
        // Sort ascending
        Collections.sort(list);                 // [1, 2, 5, 8, 9]
        
        // Sort descending
        Collections.sort(list, Collections.reverseOrder());  // [9, 8, 5, 2, 1]
        
        // Custom comparator
        Collections.sort(list, (a, b) -> b - a);  // Descending
        
        // Reverse
        Collections.reverse(list);
        
        // Shuffle
        Collections.shuffle(list);
    }
}
```

---

### Searching

```java
List<Integer> list = Arrays.asList(1, 3, 5, 7, 9);

// Binary search (list must be sorted!)
int index = Collections.binarySearch(list, 5);  // 2

// If not found, returns -(insertion point) - 1
int notFound = Collections.binarySearch(list, 6);  // -4 (would insert at index 3)
```

---

### Min/Max

```java
List<Integer> list = Arrays.asList(5, 2, 8, 1, 9);

int min = Collections.min(list);    // 1
int max = Collections.max(list);    // 9

// Custom comparator
int minAbs = Collections.min(list, (a, b) -> Math.abs(a) - Math.abs(b));
```

---

### Frequency

```java
List<Integer> list = Arrays.asList(1, 2, 2, 3, 3, 3);

int count = Collections.frequency(list, 3);  // 3
```

---

## Top 20 Interview Questions

### Q1: What's the difference between ArrayList and LinkedList?

**Answer:**
- **ArrayList:** Array-based, O(1) random access, O(n) insertion at beginning
- **LinkedList:** Node-based, O(n) random access, O(1) insertion at ends
- **Use ArrayList** 95% of the time due to better cache locality

---

### Q2: How does HashMap work internally?

**Answer:**
- Uses array of buckets (default size 16)
- Hash function: `index = hashCode() % array_length`
- **Collision handling:** Linked list (Java 7), Red-Black tree if >8 elements (Java 8+)
- **Load factor:** 0.75 (resize when 75% full)
- **Capacity:** Doubles on resize

**Hash collision example:**
```
"Aa".hashCode() == "BB".hashCode()  // Same hash, different strings
```

---

### Q3: Why is String/Integer immutable good for HashMap keys?

**Answer:**
- **hashCode() stays constant** - can't change after insertion
- If mutable: change key → wrong bucket → can't find value!
- Immutability ensures consistency

---

### Q4: HashMap vs HashTable vs ConcurrentHashMap?

| Feature | HashMap | Hashtable | ConcurrentHashMap |
|---------|---------|-----------|-------------------|
| **Thread-safe** | ❌ No | ✅ Yes | ✅ Yes |
| **Null key/value** | ✅ 1 null key | ❌ None | ❌ None |
| **Performance** | Fast | Slow (full lock) | Fast (segment locking) |
| **Legacy** | Modern | ⚠️ Legacy | Modern |

**Use:** HashMap (single-thread), ConcurrentHashMap (multi-thread)

---

### Q5: What's the time complexity of HashSet operations?

**Answer:**
- add/remove/contains: **O(1) average**, O(n) worst case
- Worst case when many hash collisions
- Uses HashMap internally (value is dummy object)

---

### Q6: How to make ArrayList thread-safe?

**Answer:**
```java
// Option 1: Collections.synchronizedList
List<Integer> syncList = Collections.synchronizedList(new ArrayList<>());

// Option 2: CopyOnWriteArrayList (for read-heavy workloads)
List<Integer> cowList = new CopyOnWriteArrayList<>();
```

---

### Q7: What's fail-fast vs fail-safe iteration?

**Answer:**

**Fail-fast (most collections):**
```java
List<Integer> list = new ArrayList<>(Arrays.asList(1, 2, 3));
for (int num : list) {
    list.remove(Integer.valueOf(num));  // ConcurrentModificationException!
}
```

**Fail-safe (concurrent collections):**
```java
CopyOnWriteArrayList<Integer> list = new CopyOnWriteArrayList<>(Arrays.asList(1, 2, 3));
for (int num : list) {
    list.remove(Integer.valueOf(num));  // Works! (operates on copy)
}
```

---

### Q8: Difference between Iterator and ListIterator?

| Feature | Iterator | ListIterator |
|---------|----------|--------------|
| **Direction** | Forward only | Bidirectional |
| **Methods** | next(), hasNext(), remove() | previous(), hasPrevious(), add(), set() |
| **Works on** | All collections | Lists only |

```java
ListIterator<Integer> it = list.listIterator();
while (it.hasNext()) {
    it.next();
}
while (it.hasPrevious()) {
    it.previous();
}
```

---

### Q9: How to sort a HashMap by values?

```java
public static Map<String, Integer> sortByValue(Map<String, Integer> map) {
    List<Map.Entry<String, Integer>> list = new ArrayList<>(map.entrySet());
    list.sort(Map.Entry.comparingByValue());
    
    Map<String, Integer> result = new LinkedHashMap<>();
    for (Map.Entry<String, Integer> entry : list) {
        result.put(entry.getKey(), entry.getValue());
    }
    return result;
}
```

---

### Q10: What's the difference between Comparable and Comparator?

**Comparable:** Natural ordering (within the class)
```java
class Student implements Comparable<Student> {
    int marks;
    
    public int compareTo(Student other) {
        return this.marks - other.marks;  // Ascending by marks
    }
}

Collections.sort(students);  // Uses compareTo
```

**Comparator:** Custom ordering (external)
```java
Collections.sort(students, (s1, s2) -> s2.marks - s1.marks);  // Descending
```

---

### Q11: When to use ArrayList vs HashSet?

- **ArrayList:** Order matters, duplicates allowed, need indexing
- **HashSet:** Only unique elements, order doesn't matter, fast lookup

---

### Q12: What's initial capacity and load factor in HashMap?

- **Initial capacity:** 16 (default)
- **Load factor:** 0.75 (resize at 75% full)
- **Resize:** Doubles capacity (16 → 32 → 64 ...)

```java
HashMap<String, Integer> map = new HashMap<>(32, 0.8f);  // Custom
```

---

### Q13: How to create immutable collection?

```java
// Java 9+
List<Integer> immutable = List.of(1, 2, 3);

// Java 8
List<Integer> immutable = Collections.unmodifiableList(
    new ArrayList<>(Arrays.asList(1, 2, 3))
);

// Trying to modify throws UnsupportedOperationException
```

---

### Q14: What's the difference between poll() and remove() in Queue?

- **poll():** Returns null if queue empty
- **remove():** Throws NoSuchElementException if empty

```java
Queue<Integer> q = new LinkedList<>();
Integer val1 = q.poll();    // null
Integer val2 = q.remove();  // Exception!
```

---

### Q15: How does TreeSet maintain order?

- Uses **Red-Black Tree** (self-balancing BST)
- Elements must be Comparable or use Comparator
- Maintains sorted order automatically

---

### Q16: Can we store null in collections?

| Collection | Null Support |
|------------|--------------|
| ArrayList | ✅ Multiple nulls |
| LinkedList | ✅ Multiple nulls |
| HashSet | ✅ One null |
| TreeSet | ❌ No null (NPE) |
| HashMap key | ✅ One null |
| HashMap value | ✅ Multiple nulls |
| TreeMap | ❌ No null key |

---

### Q17: What's the difference between Collection and Collections?

- **Collection:** Interface (root of hierarchy)
- **Collections:** Utility class (sort, reverse, etc.)

---

### Q18: How to find duplicates in array?

```java
public static List<Integer> findDuplicates(int[] arr) {
    Set<Integer> seen = new HashSet<>();
    Set<Integer> duplicates = new HashSet<>();
    
    for (int num : arr) {
        if (!seen.add(num)) {  // add() returns false if already present
            duplicates.add(num);
        }
    }
    
    return new ArrayList<>(duplicates);
}
```

---

### Q19: What's PriorityQueue and its use cases?

- **Min Heap** (default): Smallest element first
- **O(log n)** insert/remove, O(1) peek
- **Use cases:** Top K problems, median finding, merge K sorted lists

---

### Q20: Explain CopyOnWriteArrayList

- **Thread-safe** variant of ArrayList
- **Copy on write:** Creates new copy on modification
- **Use case:** Read-heavy, write-rare scenarios
- **Trade-off:** High memory, slow writes for thread-safety

---

## Performance Cheat Sheet

### Time Complexity Summary

| Operation | ArrayList | LinkedList | HashSet | TreeSet | HashMap | TreeMap |
|-----------|-----------|------------|---------|---------|---------|---------|
| **Add** | O(1) amortized | O(1) at ends | O(1) avg | O(log n) | O(1) avg | O(log n) |
| **Get by index** | O(1) | O(n) | N/A | N/A | N/A | N/A |
| **Get by key** | N/A | N/A | O(1) avg | O(log n) | O(1) avg | O(log n) |
| **Contains** | O(n) | O(n) | O(1) avg | O(log n) | O(1) avg | O(log n) |
| **Remove** | O(n) | O(1) at ends | O(1) avg | O(log n) | O(1) avg | O(log n) |
| **Iterate** | O(n) | O(n) | O(n) | O(n) | O(n) | O(n) |

---

## Decision Tree: Which Collection to Use?

```
Need key-value pairs?
├─ Yes → Map
│  ├─ Need sorted keys? → TreeMap
│  ├─ Need insertion order? → LinkedHashMap
│  └─ Otherwise → HashMap
│
└─ No → Collection
   ├─ Need unique elements?
   │  ├─ Yes → Set
   │  │  ├─ Need sorted? → TreeSet
   │  │  ├─ Need insertion order? → LinkedHashSet
   │  │  └─ Otherwise → HashSet
   │  │
   │  └─ No → List/Queue
   │     ├─ Need indexing? → ArrayList
   │     ├─ Need priority? → PriorityQueue
   │     ├─ Need both ends? → ArrayDeque
   │     └─ Otherwise → ArrayList (default)
```

---

## Key Takeaways

✅ **ArrayList** - Default choice for lists (95% of cases)
✅ **HashMap** - Default choice for key-value (O(1) operations)
✅ **HashSet** - Unique elements, fast lookup
✅ **TreeMap/TreeSet** - When you need sorted data
✅ **PriorityQueue** - Min/Max heap operations
✅ **ArrayDeque** - Stack/Queue/Deque operations
✅ **LinkedHashMap/LinkedHashSet** - When insertion order matters

❌ **Never use Vector or Hashtable** (legacy, use ArrayList/HashMap)
❌ **LinkedList rarely needed** (ArrayList usually better)

---

[← Back: I/O Handling](./IO-Handling.md) | [Next: Comparators & Comparable →](./Comparators.md)
