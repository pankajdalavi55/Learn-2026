# Java Basics for Competitive Programming

## Why This Matters

In competitive programming, you need to write code **fast and correctly**. This guide focuses on Java features you'll use in **90% of problems** - no fluff, just practical essentials.

---

## 1. Input/Output - The Fastest Way

### Basic Input with Scanner (Easy but Slow)

```java
import java.util.Scanner;

public class BasicIO {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        
        // Read different types
        int n = sc.nextInt();           // Read integer
        long l = sc.nextLong();         // Read long
        double d = sc.nextDouble();     // Read double
        String s = sc.next();           // Read single word
        String line = sc.nextLine();    // Read entire line
        
        // Read array
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = sc.nextInt();
        }
        
        sc.close();
    }
}
```

**Problem:** Scanner is SLOW for large inputs (>10^5 elements)

---

### Fast Input with BufferedReader (For Large Inputs)

```java
import java.io.*;
import java.util.*;

public class FastIO {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        // Read single integer
        int n = Integer.parseInt(br.readLine());
        
        // Read multiple integers from one line
        String[] parts = br.readLine().split(" ");
        int a = Integer.parseInt(parts[0]);
        int b = Integer.parseInt(parts[1]);
        
        // Read array efficiently
        String[] tokens = br.readLine().split(" ");
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = Integer.parseInt(tokens[i]);
        }
        
        br.close();
    }
}
```

---

### Fast Output with BufferedWriter

```java
import java.io.*;

public class FastOutput {
    public static void main(String[] args) throws IOException {
        BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(System.out));
        
        // Write output
        bw.write("Hello World\n");
        bw.write(String.valueOf(42) + "\n");
        
        // Must flush!
        bw.flush();
        bw.close();
    }
}
```

**Pro Tip:** Use StringBuilder for building output, then write once!

```java
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 1000; i++) {
    sb.append(i).append(" ");
}
System.out.println(sb);  // Much faster than 1000 print statements!
```

---

## 2. Arrays - Your Best Friend

### Array Basics

```java
public class ArrayBasics {
    public static void main(String[] args) {
        // Declaration and initialization
        int[] arr1 = new int[5];                    // [0, 0, 0, 0, 0]
        int[] arr2 = {1, 2, 3, 4, 5};              // Direct initialization
        int[] arr3 = new int[]{1, 2, 3};           // Explicit size
        
        // 2D Arrays
        int[][] matrix = new int[3][4];             // 3 rows, 4 columns
        int[][] grid = {
            {1, 2, 3},
            {4, 5, 6},
            {7, 8, 9}
        };
        
        // Access
        int val = arr2[0];                          // First element
        int lastVal = arr2[arr2.length - 1];        // Last element
        
        // Print array
        System.out.println(Arrays.toString(arr2));  // [1, 2, 3, 4, 5]
        System.out.println(Arrays.deepToString(grid)); // For 2D arrays
    }
}
```

---

### Essential Array Operations (Arrays Class)

```java
import java.util.Arrays;

public class ArrayOperations {
    public static void main(String[] args) {
        int[] arr = {5, 2, 8, 1, 9};
        
        // 1. Sort (in-place)
        Arrays.sort(arr);                   // [1, 2, 5, 8, 9]
        
        // 2. Sort in descending order (use Integer[] not int[])
        Integer[] arr2 = {5, 2, 8, 1, 9};
        Arrays.sort(arr2, Collections.reverseOrder());  // [9, 8, 5, 2, 1]
        
        // 3. Sort portion of array
        int[] arr3 = {5, 2, 8, 1, 9};
        Arrays.sort(arr3, 1, 4);           // Sort index 1 to 3: [5, 1, 2, 8, 9]
        
        // 4. Binary search (array MUST be sorted!)
        Arrays.sort(arr);
        int index = Arrays.binarySearch(arr, 5);  // Returns index of 5
        
        // 5. Fill array
        int[] arr4 = new int[5];
        Arrays.fill(arr4, 10);             // [10, 10, 10, 10, 10]
        
        // 6. Copy array
        int[] copy = Arrays.copyOf(arr, arr.length);
        int[] partial = Arrays.copyOfRange(arr, 1, 4);  // Copy elements 1 to 3
        
        // 7. Compare arrays
        int[] a = {1, 2, 3};
        int[] b = {1, 2, 3};
        boolean equal = Arrays.equals(a, b);  // true
    }
}
```

---

## 3. Strings - Immutable & Powerful

### String Basics

```java
public class StringBasics {
    public static void main(String[] args) {
        String s = "Hello";
        
        // Length
        int len = s.length();                   // 5
        
        // Access character
        char ch = s.charAt(0);                  // 'H'
        char last = s.charAt(s.length() - 1);   // 'o'
        
        // Substring
        String sub1 = s.substring(1);           // "ello"
        String sub2 = s.substring(1, 4);        // "ell" (start inclusive, end exclusive)
        
        // Concatenation
        String s2 = s + " World";               // "Hello World"
        String s3 = s.concat(" World");         // "Hello World"
        
        // Comparison
        boolean eq1 = s.equals("Hello");        // true (content)
        boolean eq2 = (s == "Hello");           // DON'T USE! (reference comparison)
        int cmp = s.compareTo("Hello");         // 0 (lexicographic)
        
        // Search
        int idx = s.indexOf('l');               // 2 (first occurrence)
        int lastIdx = s.lastIndexOf('l');       // 3 (last occurrence)
        boolean contains = s.contains("ell");   // true
        
        // Case conversion
        String upper = s.toUpperCase();         // "HELLO"
        String lower = s.toLowerCase();         // "hello"
        
        // Split
        String sentence = "one two three";
        String[] words = sentence.split(" ");   // ["one", "two", "three"]
        
        // Convert to char array
        char[] chars = s.toCharArray();         // ['H', 'e', 'l', 'l', 'o']
    }
}
```

---

### StringBuilder - For String Manipulation

**CRITICAL:** Strings are immutable! Concatenation creates new objects.

```java
// ❌ SLOW - Creates n new String objects!
String result = "";
for (int i = 0; i < 1000; i++) {
    result += i;  // Each += creates new String!
}

// ✅ FAST - Single mutable object
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 1000; i++) {
    sb.append(i);
}
String result = sb.toString();
```

**StringBuilder Operations:**

```java
public class StringBuilderDemo {
    public static void main(String[] args) {
        StringBuilder sb = new StringBuilder();
        
        // Append
        sb.append("Hello");
        sb.append(" ").append("World");         // Chaining
        sb.append(123);                         // Append number
        
        // Insert
        sb.insert(5, ",");                      // "Hello, World123"
        
        // Delete
        sb.delete(5, 6);                        // Remove comma: "Hello World123"
        sb.deleteCharAt(5);                     // Remove at index 5
        
        // Replace
        sb.replace(0, 5, "Hi");                 // "Hi World123"
        
        // Reverse
        sb.reverse();                           // "321dlroW iH"
        
        // Set character
        sb.setCharAt(0, 'X');
        
        // Length and capacity
        int len = sb.length();
        int cap = sb.capacity();
        
        // Convert to String
        String result = sb.toString();
    }
}
```

---

## 4. Collections Framework - Must Know

### ArrayList - Dynamic Array

```java
import java.util.*;

public class ArrayListDemo {
    public static void main(String[] args) {
        // Creation
        ArrayList<Integer> list = new ArrayList<>();
        ArrayList<String> list2 = new ArrayList<>(10);  // Initial capacity
        
        // Add elements
        list.add(10);                       // [10]
        list.add(20);                       // [10, 20]
        list.add(1, 15);                    // [10, 15, 20] (insert at index 1)
        
        // Access
        int val = list.get(0);              // 10
        list.set(0, 5);                     // Change to 5: [5, 15, 20]
        
        // Size
        int size = list.size();             // 3
        boolean empty = list.isEmpty();     // false
        
        // Remove
        list.remove(1);                     // Remove at index 1: [5, 20]
        list.remove(Integer.valueOf(20));   // Remove value 20: [5]
        
        // Search
        boolean contains = list.contains(5);    // true
        int index = list.indexOf(5);            // 0
        
        // Iterate
        for (int num : list) {
            System.out.println(num);
        }
        
        // Convert to array
        Integer[] arr = list.toArray(new Integer[0]);
        
        // Sort
        Collections.sort(list);                         // Ascending
        Collections.sort(list, Collections.reverseOrder()); // Descending
    }
}
```

---

### HashMap - Key-Value Pairs

```java
import java.util.*;

public class HashMapDemo {
    public static void main(String[] args) {
        HashMap<String, Integer> map = new HashMap<>();
        
        // Add/Update
        map.put("Alice", 25);
        map.put("Bob", 30);
        map.put("Alice", 26);           // Updates Alice's value
        
        // Get
        int age = map.get("Alice");     // 26
        int def = map.getOrDefault("Charlie", 0);  // 0 (not found)
        
        // Check existence
        boolean hasKey = map.containsKey("Bob");    // true
        boolean hasVal = map.containsValue(30);     // true
        
        // Remove
        map.remove("Bob");
        
        // Size
        int size = map.size();
        
        // Iterate through keys
        for (String key : map.keySet()) {
            System.out.println(key + ": " + map.get(key));
        }
        
        // Iterate through entries (better!)
        for (Map.Entry<String, Integer> entry : map.entrySet()) {
            System.out.println(entry.getKey() + ": " + entry.getValue());
        }
        
        // Get all keys/values
        Set<String> keys = map.keySet();
        Collection<Integer> values = map.values();
    }
}
```

**Common Pattern - Frequency Counter:**

```java
public static HashMap<Character, Integer> charFrequency(String s) {
    HashMap<Character, Integer> freq = new HashMap<>();
    for (char c : s.toCharArray()) {
        freq.put(c, freq.getOrDefault(c, 0) + 1);
    }
    return freq;
}
```

---

### HashSet - Unique Elements

```java
import java.util.*;

public class HashSetDemo {
    public static void main(String[] args) {
        HashSet<Integer> set = new HashSet<>();
        
        // Add
        set.add(10);
        set.add(20);
        set.add(10);            // Duplicate ignored
        
        // Contains
        boolean has = set.contains(10);     // true
        
        // Remove
        set.remove(10);
        
        // Size
        int size = set.size();              // 1
        
        // Iterate
        for (int num : set) {
            System.out.println(num);
        }
        
        // Convert ArrayList to HashSet (remove duplicates)
        ArrayList<Integer> list = new ArrayList<>(Arrays.asList(1, 2, 2, 3, 3, 3));
        HashSet<Integer> unique = new HashSet<>(list);  // {1, 2, 3}
    }
}
```

---

### TreeMap & TreeSet - Sorted Collections

```java
import java.util.*;

public class TreeCollections {
    public static void main(String[] args) {
        // TreeSet - Sorted unique elements
        TreeSet<Integer> treeSet = new TreeSet<>();
        treeSet.add(5);
        treeSet.add(2);
        treeSet.add(8);
        treeSet.add(1);
        // Stored as: [1, 2, 5, 8] (automatically sorted)
        
        int first = treeSet.first();        // 1
        int last = treeSet.last();          // 8
        int ceil = treeSet.ceiling(3);      // 5 (smallest >= 3)
        int floor = treeSet.floor(3);       // 2 (largest <= 3)
        
        // TreeMap - Sorted keys
        TreeMap<Integer, String> treeMap = new TreeMap<>();
        treeMap.put(3, "three");
        treeMap.put(1, "one");
        treeMap.put(2, "two");
        // Keys stored sorted: {1=one, 2=two, 3=three}
        
        int firstKey = treeMap.firstKey();          // 1
        int lastKey = treeMap.lastKey();            // 3
        Map.Entry<Integer, String> firstEntry = treeMap.firstEntry();
    }
}
```

**When to use:**
- **HashMap/HashSet:** Fast (O(1) average), unordered
- **TreeMap/TreeSet:** Sorted (O(log n)), ordered operations

---

## 5. Stack & Queue

### Stack (LIFO - Last In First Out)

```java
import java.util.*;

public class StackDemo {
    public static void main(String[] args) {
        Stack<Integer> stack = new Stack<>();
        
        // Push
        stack.push(10);
        stack.push(20);
        stack.push(30);     // [10, 20, 30] (30 on top)
        
        // Peek (look at top without removing)
        int top = stack.peek();         // 30
        
        // Pop (remove and return top)
        int removed = stack.pop();      // 30, stack: [10, 20]
        
        // Check if empty
        boolean empty = stack.isEmpty();
        
        // Size
        int size = stack.size();
    }
}
```

**Common Use:** Parentheses matching, undo operations, DFS

---

### Queue (FIFO - First In First Out)

```java
import java.util.*;

public class QueueDemo {
    public static void main(String[] args) {
        Queue<Integer> queue = new LinkedList<>();
        
        // Add to back
        queue.offer(10);    // or add()
        queue.offer(20);
        queue.offer(30);    // [10, 20, 30] (10 at front)
        
        // Peek (look at front without removing)
        int front = queue.peek();       // 10
        
        // Remove from front
        int removed = queue.poll();     // 10, queue: [20, 30]
        
        // Check if empty
        boolean empty = queue.isEmpty();
        
        // Size
        int size = queue.size();
    }
}
```

**Common Use:** BFS, level-order traversal, process scheduling

---

### Deque (Double-Ended Queue)

```java
import java.util.*;

public class DequeDemo {
    public static void main(String[] args) {
        Deque<Integer> deque = new ArrayDeque<>();
        
        // Add to front
        deque.addFirst(10);     // [10]
        deque.addFirst(5);      // [5, 10]
        
        // Add to back
        deque.addLast(20);      // [5, 10, 20]
        
        // Remove from front
        int front = deque.removeFirst();    // 5, deque: [10, 20]
        
        // Remove from back
        int back = deque.removeLast();      // 20, deque: [10]
        
        // Peek both ends
        int peekFront = deque.peekFirst();
        int peekBack = deque.peekLast();
    }
}
```

**Use as Stack:** `addFirst()`, `removeFirst()`
**Use as Queue:** `addLast()`, `removeFirst()`

---

## 6. PriorityQueue (Heap)

```java
import java.util.*;

public class PriorityQueueDemo {
    public static void main(String[] args) {
        // Min Heap (default)
        PriorityQueue<Integer> minHeap = new PriorityQueue<>();
        minHeap.offer(5);
        minHeap.offer(2);
        minHeap.offer(8);
        minHeap.offer(1);
        
        System.out.println(minHeap.poll());  // 1 (smallest)
        System.out.println(minHeap.poll());  // 2
        
        // Max Heap
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
        maxHeap.offer(5);
        maxHeap.offer(2);
        maxHeap.offer(8);
        
        System.out.println(maxHeap.poll());  // 8 (largest)
        
        // Custom comparator
        PriorityQueue<int[]> pq = new PriorityQueue<>((a, b) -> a[0] - b[0]);
        pq.offer(new int[]{3, 100});
        pq.offer(new int[]{1, 200});
        pq.offer(new int[]{2, 300});
        
        int[] min = pq.poll();  // [1, 200] (smallest first element)
    }
}
```

---

## 7. Sorting & Comparators

### Sorting Primitives

```java
import java.util.*;

public class SortingPrimitives {
    public static void main(String[] args) {
        int[] arr = {5, 2, 8, 1, 9};
        Arrays.sort(arr);           // [1, 2, 5, 8, 9]
        
        // Can't directly sort primitives in descending order
        // Must use Integer[] wrapper
    }
}
```

---

### Sorting Objects with Comparators

```java
import java.util.*;

class Student {
    String name;
    int marks;
    
    Student(String name, int marks) {
        this.name = name;
        this.marks = marks;
    }
}

public class ComparatorDemo {
    public static void main(String[] args) {
        ArrayList<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 78));
        
        // Sort by marks (ascending)
        Collections.sort(students, (a, b) -> a.marks - b.marks);
        
        // Sort by marks (descending)
        Collections.sort(students, (a, b) -> b.marks - a.marks);
        
        // Sort by name
        Collections.sort(students, (a, b) -> a.name.compareTo(b.name));
        
        // Multiple criteria: first by marks desc, then by name asc
        Collections.sort(students, (a, b) -> {
            if (a.marks != b.marks) {
                return b.marks - a.marks;  // Descending marks
            }
            return a.name.compareTo(b.name);  // Ascending name
        });
    }
}
```

**Common Comparator Patterns:**

```java
// Sort 2D array by first element
int[][] pairs = {{3, 4}, {1, 2}, {5, 6}};
Arrays.sort(pairs, (a, b) -> a[0] - b[0]);

// Sort strings by length
String[] words = {"apple", "pie", "banana"};
Arrays.sort(words, (a, b) -> a.length() - b.length());

// Sort in reverse
Arrays.sort(arr, Collections.reverseOrder());
```

---

## 8. Math Utilities

```java
public class MathUtils {
    public static void main(String[] args) {
        // Absolute value
        int abs = Math.abs(-5);             // 5
        
        // Max/Min
        int max = Math.max(10, 20);         // 20
        int min = Math.min(10, 20);         // 10
        
        // Power
        double pow = Math.pow(2, 10);       // 1024.0
        
        // Square root
        double sqrt = Math.sqrt(16);        // 4.0
        
        // Ceiling/Floor
        double ceil = Math.ceil(4.2);       // 5.0
        double floor = Math.floor(4.8);     // 4.0
        
        // Random
        double rand = Math.random();        // [0.0, 1.0)
        int randInt = (int)(Math.random() * 100);  // [0, 99]
        
        // GCD (using recursion)
        int gcd = gcd(12, 18);              // 6
    }
    
    public static int gcd(int a, int b) {
        return b == 0 ? a : gcd(b, a % b);
    }
}
```

---

## 9. Common Pitfalls & Tips

### ❌ Integer Overflow

```java
int a = 1000000;
int b = 1000000;
int product = a * b;            // OVERFLOW! Result is negative

// ✅ Fix: Use long
long product = (long)a * b;     // Correct
```

---

### ❌ Array Index Out of Bounds

```java
int[] arr = new int[5];
int val = arr[5];               // ERROR! Valid indices: 0-4

// ✅ Always check: i < arr.length
```

---

### ❌ Comparing Strings with ==

```java
String a = new String("hello");
String b = new String("hello");
if (a == b) { }                 // FALSE! Compares references

// ✅ Use .equals()
if (a.equals(b)) { }            // TRUE! Compares content
```

---

### ❌ Modifying Collection While Iterating

```java
ArrayList<Integer> list = new ArrayList<>(Arrays.asList(1, 2, 3, 4));
for (int num : list) {
    list.remove(Integer.valueOf(num));  // ConcurrentModificationException!
}

// ✅ Use Iterator
Iterator<Integer> it = list.iterator();
while (it.hasNext()) {
    it.next();
    it.remove();  // Safe removal
}
```

---

## 10. Quick Reference Template

```java
import java.io.*;
import java.util.*;

public class Solution {
    public static void main(String[] args) throws IOException {
        // Fast Input
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        // Read integers
        int n = Integer.parseInt(br.readLine());
        String[] tokens = br.readLine().split(" ");
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = Integer.parseInt(tokens[i]);
        }
        
        // Your solution here
        
        // Fast Output
        StringBuilder sb = new StringBuilder();
        for (int num : arr) {
            sb.append(num).append(" ");
        }
        System.out.println(sb);
    }
}
```

---

## Key Takeaways

✅ **Use BufferedReader for large inputs** - Scanner is slow
✅ **Use StringBuilder for string building** - String concatenation is slow
✅ **Know your collections** - HashMap, HashSet, ArrayList, PriorityQueue
✅ **Master Arrays class** - sort, binarySearch, fill, copyOf
✅ **Watch for integer overflow** - use long when needed
✅ **Use .equals() for strings** - not ==
✅ **Practice comparators** - sorting custom objects is common

---

[Next: Input/Output Handling →](./IO-Handling.md)
