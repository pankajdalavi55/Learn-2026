# Comparators & Comparable - Complete Sorting Guide

## Why Sorting Matters in Competitive Programming

**90% of problems** involve sorting at some point:
- Finding min/max elements
- Binary search (requires sorted data)
- Greedy algorithms
- Meeting room problems
- Interval problems

**You MUST master custom sorting!**

---

## The Problem: Sorting Custom Objects

```java
class Student {
    String name;
    int marks;
    
    Student(String name, int marks) {
        this.name = name;
        this.marks = marks;
    }
}

public class SortingProblem {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 78));
        
        // How do we sort this?
        // Collections.sort(students);  // ERROR! Doesn't know how to compare Students
    }
}
```

**Solution:** Use **Comparable** or **Comparator**!

---

## Comparable - Natural Ordering (Inside the Class)

### Theory

- **Interface:** `java.lang.Comparable<T>`
- **Method:** `int compareTo(T other)`
- **Purpose:** Define the "natural" ordering for a class
- **Location:** Inside the class being compared
- **Use:** When there's ONE obvious way to sort

### Return Values

```
compareTo() returns:
  Negative (< 0) : this < other  → this comes BEFORE other
  Zero (0)       : this == other → equal
  Positive (> 0) : this > other  → this comes AFTER other
```

---

### Example 1: Sorting Students by Marks (Ascending)

```java
class Student implements Comparable<Student> {
    String name;
    int marks;
    
    Student(String name, int marks) {
        this.name = name;
        this.marks = marks;
    }
    
    @Override
    public int compareTo(Student other) {
        // Ascending order by marks
        return this.marks - other.marks;
        
        // Same as:
        // if (this.marks < other.marks) return -1;
        // if (this.marks > other.marks) return 1;
        // return 0;
    }
    
    @Override
    public String toString() {
        return name + "(" + marks + ")";
    }
}

public class ComparableDemo {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 78));
        
        Collections.sort(students);  // Uses compareTo()
        
        System.out.println(students);
        // Output: [Charlie(78), Alice(85), Bob(92)] - Ascending
    }
}
```

---

### Example 2: Sorting by Multiple Criteria

```java
class Student implements Comparable<Student> {
    String name;
    int marks;
    
    Student(String name, int marks) {
        this.name = name;
        this.marks = marks;
    }
    
    @Override
    public int compareTo(Student other) {
        // First, compare by marks (descending)
        if (this.marks != other.marks) {
            return other.marks - this.marks;  // Descending!
        }
        // If marks are same, compare by name (ascending)
        return this.name.compareTo(other.name);
    }
    
    @Override
    public String toString() {
        return name + "(" + marks + ")";
    }
}

public class MultiCriteriaComparable {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 85));
        students.add(new Student("David", 92));
        
        Collections.sort(students);
        
        System.out.println(students);
        // Output: [Bob(92), David(92), Alice(85), Charlie(85)]
        // First by marks descending, then by name ascending
    }
}
```

---

### Built-in Comparable Classes

These already implement Comparable:

```java
import java.util.*;

public class BuiltInComparable {
    public static void main(String[] args) {
        
        // Integer - natural order is ascending
        List<Integer> nums = Arrays.asList(5, 2, 8, 1, 9);
        Collections.sort(nums);
        System.out.println(nums);  // [1, 2, 5, 8, 9]
        
        // String - natural order is lexicographic (dictionary)
        List<String> words = Arrays.asList("banana", "apple", "cherry");
        Collections.sort(words);
        System.out.println(words);  // [apple, banana, cherry]
        
        // Character
        List<Character> chars = Arrays.asList('z', 'a', 'm');
        Collections.sort(chars);
        System.out.println(chars);  // [a, m, z]
    }
}
```

---

## Comparator - Custom Ordering (Outside the Class)

### Theory

- **Interface:** `java.util.Comparator<T>`
- **Method:** `int compare(T o1, T o2)`
- **Purpose:** Define custom ordering (can have multiple!)
- **Location:** Separate class or lambda
- **Use:** When you need different ways to sort

### Return Values

```
compare() returns:
  Negative (< 0) : o1 < o2  → o1 comes BEFORE o2
  Zero (0)       : o1 == o2 → equal
  Positive (> 0) : o1 > o2  → o1 comes AFTER o2
```

---

### Example 1: Sorting with Anonymous Class

```java
class Student {
    String name;
    int marks;
    
    Student(String name, int marks) {
        this.name = name;
        this.marks = marks;
    }
    
    @Override
    public String toString() {
        return name + "(" + marks + ")";
    }
}

public class ComparatorAnonymous {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 78));
        
        // Sort by marks ascending
        Collections.sort(students, new Comparator<Student>() {
            @Override
            public int compare(Student s1, Student s2) {
                return s1.marks - s2.marks;
            }
        });
        
        System.out.println(students);
        // [Charlie(78), Alice(85), Bob(92)]
    }
}
```

---

### Example 2: Sorting with Lambda (Java 8+)

```java
import java.util.*;

public class ComparatorLambda {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 78));
        
        // Sort by marks ascending (Lambda)
        Collections.sort(students, (s1, s2) -> s1.marks - s2.marks);
        System.out.println(students);
        
        // Sort by marks descending
        Collections.sort(students, (s1, s2) -> s2.marks - s1.marks);
        System.out.println(students);
        
        // Sort by name
        Collections.sort(students, (s1, s2) -> s1.name.compareTo(s2.name));
        System.out.println(students);
    }
}
```

---

### Example 3: Using Comparator Helper Methods (Java 8+)

```java
import java.util.*;

public class ComparatorHelpers {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student("Alice", 85));
        students.add(new Student("Bob", 92));
        students.add(new Student("Charlie", 78));
        
        // comparing() - Extract key and compare
        students.sort(Comparator.comparing(s -> s.marks));
        // OR with method reference
        students.sort(Comparator.comparing(Student::getMarks));
        
        // reversed() - Reverse order
        students.sort(Comparator.comparing(s -> s.marks).reversed());
        
        // comparingInt() - For primitive int (more efficient)
        students.sort(Comparator.comparingInt(s -> s.marks));
        
        // thenComparing() - Multiple criteria
        students.sort(
            Comparator.comparing((Student s) -> s.marks)
                     .reversed()
                     .thenComparing(s -> s.name)
        );
        
        // nullsFirst() / nullsLast() - Handle nulls
        students.sort(Comparator.nullsFirst(Comparator.comparing(s -> s.name)));
    }
}

class Student {
    String name;
    int marks;
    
    Student(String name, int marks) {
        this.name = name;
        this.marks = marks;
    }
    
    public int getMarks() { return marks; }
    public String getName() { return name; }
    
    @Override
    public String toString() {
        return name + "(" + marks + ")";
    }
}
```

---

## Comparable vs Comparator - Complete Comparison

| Aspect | Comparable | Comparator |
|--------|-----------|------------|
| **Package** | java.lang | java.util |
| **Method** | `compareTo(T other)` | `compare(T o1, T o2)` |
| **Location** | Inside the class | Separate class/lambda |
| **Parameters** | 1 (compares with this) | 2 (compares both) |
| **Sorting ways** | ONE natural order | MULTIPLE custom orders |
| **Modify class** | YES - must implement | NO - external |
| **Use case** | Default/natural sorting | Custom/flexible sorting |
| **Example** | String, Integer | Lambda expressions |

---

## When to Use What?

### Use Comparable when:
✅ There's ONE obvious way to sort (e.g., Person by age)
✅ You own the class and can modify it
✅ You want a default natural ordering

### Use Comparator when:
✅ Need multiple ways to sort (by name, age, salary, etc.)
✅ Don't own the class (can't modify)
✅ Need to sort in different ways at different times
✅ One-time custom sorting

---

## Common Sorting Patterns for Competitive Programming

### Pattern 1: Sort 2D Array by First Element

```java
public class Sort2DArray {
    public static void main(String[] args) {
        int[][] pairs = {
            {3, 4},
            {1, 2},
            {5, 6},
            {1, 3}
        };
        
        // Sort by first element ascending
        Arrays.sort(pairs, (a, b) -> a[0] - b[0]);
        
        // Sort by first element ascending, then second descending
        Arrays.sort(pairs, (a, b) -> {
            if (a[0] != b[0]) {
                return a[0] - b[0];
            }
            return b[1] - a[1];
        });
        
        // Print
        for (int[] pair : pairs) {
            System.out.println(Arrays.toString(pair));
        }
    }
}
```

---

### Pattern 2: Sort by String Length

```java
public class SortByLength {
    public static void main(String[] args) {
        String[] words = {"apple", "pie", "banana", "cat"};
        
        // Sort by length ascending
        Arrays.sort(words, (a, b) -> a.length() - b.length());
        System.out.println(Arrays.toString(words));
        // [pie, cat, apple, banana]
        
        // Sort by length descending, then alphabetically
        Arrays.sort(words, (a, b) -> {
            if (a.length() != b.length()) {
                return b.length() - a.length();
            }
            return a.compareTo(b);
        });
        System.out.println(Arrays.toString(words));
        // [banana, apple, cat, pie]
    }
}
```

---

### Pattern 3: Sort by Absolute Value

```java
public class SortByAbsoluteValue {
    public static void main(String[] args) {
        Integer[] nums = {-5, 2, -8, 1, 9, -3};
        
        // Sort by absolute value
        Arrays.sort(nums, (a, b) -> Math.abs(a) - Math.abs(b));
        System.out.println(Arrays.toString(nums));
        // [1, 2, -3, -5, -8, 9]
        
        // Or with Comparator.comparingInt
        Arrays.sort(nums, Comparator.comparingInt(Math::abs));
    }
}
```

---

### Pattern 4: Sort Map by Values

```java
import java.util.*;
import java.util.stream.*;

public class SortMapByValue {
    public static void main(String[] args) {
        Map<String, Integer> scores = new HashMap<>();
        scores.put("Alice", 85);
        scores.put("Bob", 92);
        scores.put("Charlie", 78);
        scores.put("David", 92);
        
        // Convert to list of entries and sort
        List<Map.Entry<String, Integer>> list = new ArrayList<>(scores.entrySet());
        list.sort((e1, e2) -> e2.getValue() - e1.getValue());  // Descending
        
        // Or create sorted LinkedHashMap
        Map<String, Integer> sortedMap = list.stream()
            .collect(Collectors.toMap(
                Map.Entry::getKey,
                Map.Entry::getValue,
                (e1, e2) -> e1,
                LinkedHashMap::new
            ));
        
        System.out.println(sortedMap);
        // {Bob=92, David=92, Alice=85, Charlie=78}
    }
}
```

---

### Pattern 5: Sort with Custom Object (Intervals)

```java
import java.util.*;

class Interval {
    int start, end;
    
    Interval(int start, int end) {
        this.start = start;
        this.end = end;
    }
    
    @Override
    public String toString() {
        return "[" + start + "," + end + "]";
    }
}

public class SortIntervals {
    public static void main(String[] args) {
        List<Interval> intervals = new ArrayList<>();
        intervals.add(new Interval(1, 3));
        intervals.add(new Interval(2, 6));
        intervals.add(new Interval(8, 10));
        intervals.add(new Interval(15, 18));
        
        // Sort by start time
        intervals.sort((a, b) -> a.start - b.start);
        
        // Sort by end time
        intervals.sort((a, b) -> a.end - b.end);
        
        // Sort by start, then by end if start is same
        intervals.sort((a, b) -> {
            if (a.start != b.start) {
                return a.start - b.start;
            }
            return a.end - b.end;
        });
        
        System.out.println(intervals);
    }
}
```

---

### Pattern 6: Reverse Order

```java
import java.util.*;

public class ReverseOrder {
    public static void main(String[] args) {
        
        // Method 1: Collections.reverseOrder()
        Integer[] nums = {5, 2, 8, 1, 9};
        Arrays.sort(nums, Collections.reverseOrder());
        System.out.println(Arrays.toString(nums));
        // [9, 8, 5, 2, 1]
        
        // Method 2: Lambda
        Arrays.sort(nums, (a, b) -> b - a);
        
        // Method 3: Comparator.reverseOrder()
        Arrays.sort(nums, Comparator.reverseOrder());
        
        // For primitives (int[]), need to convert
        int[] primitives = {5, 2, 8, 1, 9};
        Arrays.sort(primitives);  // Sort ascending first
        // Reverse manually
        for (int i = 0; i < primitives.length / 2; i++) {
            int temp = primitives[i];
            primitives[i] = primitives[primitives.length - 1 - i];
            primitives[primitives.length - 1 - i] = temp;
        }
    }
}
```

---

## Priority Queue with Custom Comparator

```java
import java.util.*;

public class PriorityQueueComparator {
    public static void main(String[] args) {
        
        // Min heap (default) - smallest element first
        PriorityQueue<Integer> minHeap = new PriorityQueue<>();
        minHeap.offer(5);
        minHeap.offer(2);
        minHeap.offer(8);
        System.out.println(minHeap.poll());  // 2
        
        // Max heap - largest element first
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>((a, b) -> b - a);
        // OR
        PriorityQueue<Integer> maxHeap2 = new PriorityQueue<>(Collections.reverseOrder());
        maxHeap.offer(5);
        maxHeap.offer(2);
        maxHeap.offer(8);
        System.out.println(maxHeap.poll());  // 8
        
        // Custom objects
        PriorityQueue<Student> pq = new PriorityQueue<>((a, b) -> b.marks - a.marks);
        pq.offer(new Student("Alice", 85));
        pq.offer(new Student("Bob", 92));
        pq.offer(new Student("Charlie", 78));
        
        System.out.println(pq.poll());  // Bob(92) - highest marks
    }
}
```

---

## TreeSet/TreeMap with Comparator

```java
import java.util.*;

public class TreeComparator {
    public static void main(String[] args) {
        
        // TreeSet with custom comparator
        TreeSet<Integer> set = new TreeSet<>((a, b) -> b - a);  // Descending
        set.add(5);
        set.add(2);
        set.add(8);
        System.out.println(set);  // [8, 5, 2]
        
        // TreeMap with custom comparator (sort by keys)
        TreeMap<String, Integer> map = new TreeMap<>((a, b) -> b.compareTo(a));
        map.put("Alice", 85);
        map.put("Bob", 92);
        map.put("Charlie", 78);
        System.out.println(map);  // {Charlie=78, Bob=92, Alice=85} - reverse alphabetical
    }
}
```

---

## Common Pitfalls & Solutions

### Pitfall 1: Integer Overflow in Subtraction

```java
// ❌ WRONG - Can cause overflow!
Arrays.sort(array, (a, b) -> a - b);

// If a = Integer.MAX_VALUE and b = -1
// a - b = overflow (becomes negative!)

// ✅ CORRECT - Use Integer.compare()
Arrays.sort(array, (a, b) -> Integer.compare(a, b));

// OR
Arrays.sort(array, Integer::compare);
```

---

### Pitfall 2: Comparing Primitives vs Objects

```java
// ❌ Can't use reverseOrder() on primitive array
int[] nums = {5, 2, 8};
// Arrays.sort(nums, Collections.reverseOrder());  // ERROR!

// ✅ Use Integer[] (wrapper class)
Integer[] nums2 = {5, 2, 8};
Arrays.sort(nums2, Collections.reverseOrder());  // OK!
```

---

### Pitfall 3: Null Values

```java
List<String> list = Arrays.asList("apple", null, "banana", null, "cherry");

// ❌ NullPointerException!
// list.sort((a, b) -> a.compareTo(b));

// ✅ Handle nulls
list.sort(Comparator.nullsFirst(String::compareTo));
// OR
list.sort(Comparator.nullsLast(String::compareTo));
```

---

### Pitfall 4: Modifying Objects During Sort

```java
// ❌ NEVER modify objects while sorting!
List<Student> students = new ArrayList<>();
students.sort((a, b) -> {
    a.marks = 100;  // DON'T DO THIS!
    return a.marks - b.marks;
});

// ✅ Only read values during comparison
students.sort((a, b) -> a.marks - b.marks);
```

---

## Top 15 Interview Questions

### Q1: What's the difference between Comparable and Comparator?

**Answer:**

| Comparable | Comparator |
|-----------|------------|
| Inside class | Outside class |
| compareTo(T other) | compare(T o1, T o2) |
| Natural ordering | Custom ordering |
| One way to sort | Multiple ways |
| Modify class required | No modification needed |

---

### Q2: Can a class implement both Comparable and Comparator?

**Answer:** Yes! But it's unusual. Class can implement Comparable for natural ordering and use Comparator for custom orderings.

```java
class Student implements Comparable<Student>, Comparator<Student> {
    String name;
    int marks;
    
    // Comparable - natural order by marks
    public int compareTo(Student other) {
        return this.marks - other.marks;
    }
    
    // Comparator - custom order
    public int compare(Student s1, Student s2) {
        return s1.name.compareTo(s2.name);
    }
}
```

---

### Q3: What happens if compareTo() is inconsistent with equals()?

**Answer:** Violates general contract. Can cause unexpected behavior in sorted collections (TreeSet, TreeMap).

```java
// BAD: compareTo() says equal but equals() says not equal
class Bad implements Comparable<Bad> {
    int value;
    
    public int compareTo(Bad other) {
        return 0;  // Always equal
    }
    
    public boolean equals(Object obj) {
        return false;  // Never equal
    }
}

// TreeSet uses compareTo(), so all elements appear equal
TreeSet<Bad> set = new TreeSet<>();
set.add(new Bad());
set.add(new Bad());
System.out.println(set.size());  // 1 (should be 2)
```

---

### Q4: How to sort in descending order?

**Answer:**

```java
// Method 1: Collections.reverseOrder()
Collections.sort(list, Collections.reverseOrder());

// Method 2: Lambda
Collections.sort(list, (a, b) -> b - a);

// Method 3: Comparator.reverseOrder()
list.sort(Comparator.reverseOrder());

// Method 4: reversed()
list.sort(Comparator.naturalOrder().reversed());
```

---

### Q5: Why use Integer.compare() instead of subtraction?

**Answer:** To avoid integer overflow!

```java
int a = Integer.MAX_VALUE;
int b = -1;

// ❌ Overflow: a - b becomes negative!
int wrong = a - b;  // Negative due to overflow

// ✅ Correct
int correct = Integer.compare(a, b);  // Positive
```

---

### Q6: Can we sort a List in place?

**Answer:** Yes! Multiple ways:

```java
List<Integer> list = new ArrayList<>(Arrays.asList(5, 2, 8, 1));

// Method 1: Collections.sort()
Collections.sort(list);

// Method 2: List.sort() (Java 8+)
list.sort(Comparator.naturalOrder());

// Method 3: Stream (creates new list)
List<Integer> sorted = list.stream().sorted().collect(Collectors.toList());
```

---

### Q7: How to break ties in sorting?

**Answer:** Use thenComparing():

```java
students.sort(
    Comparator.comparing((Student s) -> s.marks)
             .reversed()  // Primary: marks descending
             .thenComparing(s -> s.name)  // Tie-breaker: name ascending
);
```

---

### Q8: Can we use Comparator with TreeSet?

**Answer:** Yes! TreeSet uses comparator for ordering:

```java
// Descending order
TreeSet<Integer> set = new TreeSet<>((a, b) -> b - a);
set.add(5);
set.add(2);
set.add(8);
System.out.println(set);  // [8, 5, 2]
```

---

### Q9: What's the time complexity of sorting?

**Answer:**
- Arrays.sort() / Collections.sort(): **O(n log n)** (Dual-Pivot Quicksort for primitives, Timsort for objects)
- TreeSet/TreeMap insertion: **O(log n)** per element, so n insertions = **O(n log n)**

---

### Q10: How to sort by multiple fields?

**Answer:**

```java
// Method 1: Manual if-else
students.sort((a, b) -> {
    if (a.marks != b.marks) {
        return b.marks - a.marks;  // Descending marks
    }
    return a.name.compareTo(b.name);  // Ascending name
});

// Method 2: thenComparing()
students.sort(
    Comparator.comparingInt((Student s) -> s.marks)
             .reversed()
             .thenComparing(s -> s.name)
);
```

---

### Q11: Can primitives be sorted in descending order directly?

**Answer:** No! Must use wrapper classes:

```java
// ❌ Doesn't work
int[] nums = {5, 2, 8};
// Arrays.sort(nums, Collections.reverseOrder());  // ERROR

// ✅ Use Integer[]
Integer[] nums2 = {5, 2, 8};
Arrays.sort(nums2, Collections.reverseOrder());
```

---

### Q12: What's the difference between sorted() and sort()?

**Answer:**

```java
List<Integer> list = Arrays.asList(5, 2, 8);

// sort() - In-place, void return
list.sort(Comparator.naturalOrder());

// sorted() - Stream operation, returns new stream
List<Integer> newList = list.stream()
                            .sorted()
                            .collect(Collectors.toList());
```

---

### Q13: How to make a class sortable?

**Answer:**

```java
// Option 1: Implement Comparable
class Student implements Comparable<Student> {
    int marks;
    
    public int compareTo(Student other) {
        return Integer.compare(this.marks, other.marks);
    }
}

// Option 2: Provide Comparator when sorting
Collections.sort(students, (a, b) -> a.marks - b.marks);
```

---

### Q14: Can we chain multiple Comparators?

**Answer:** Yes, using thenComparing():

```java
Comparator<Student> comp = Comparator
    .comparingInt((Student s) -> s.marks).reversed()
    .thenComparing(s -> s.name)
    .thenComparingInt(s -> s.age);
```

---

### Q15: What's the difference between stable and unstable sort?

**Answer:**

**Stable sort:** Maintains relative order of equal elements
**Unstable sort:** May change relative order of equal elements

Java's:
- **Arrays.sort()** for primitives: Unstable (Quicksort)
- **Arrays.sort()** for objects: **Stable** (Timsort)
- **Collections.sort()**: **Stable** (Timsort)

```java
// Stable sort example
class Student {
    String name;
    int marks;
}

List<Student> students = Arrays.asList(
    new Student("Alice", 85),
    new Student("Bob", 85),
    new Student("Charlie", 85)
);

// Sort by marks (all equal) - stable sort maintains order
Collections.sort(students, (a, b) -> a.marks - b.marks);
// Order: Alice, Bob, Charlie (maintained)
```

---

## Quick Reference Cheat Sheet

### Sorting Syntax Quick Guide

```java
// Arrays
Arrays.sort(array);                                    // Ascending
Arrays.sort(array, Collections.reverseOrder());        // Descending (wrapper only)
Arrays.sort(array, (a, b) -> a[0] - b[0]);            // Custom

// Lists
Collections.sort(list);                                // Ascending
Collections.sort(list, Collections.reverseOrder());    // Descending
list.sort((a, b) -> a - b);                           // Lambda

// Stream
list.stream().sorted().collect(Collectors.toList());
list.stream().sorted(Comparator.reverseOrder()).collect(Collectors.toList());
```

---

### Common Comparator Patterns

```java
// By field
Comparator.comparing(Student::getName)

// By primitive field
Comparator.comparingInt(Student::getMarks)

// Reversed
Comparator.comparing(Student::getMarks).reversed()

// Multiple criteria
Comparator.comparing(Student::getMarks)
         .thenComparing(Student::getName)

// Handle nulls
Comparator.nullsFirst(Comparator.naturalOrder())
```

---

## Practice Problems with Solutions

### Problem 1: Meeting Rooms - Can Attend All Meetings?

**Problem:** Given an array of meeting time intervals, determine if a person could attend all meetings.

**Example:** `[[0,30],[5,10],[15,20]]` → `false` (conflicts)

```java
import java.util.*;

class Interval {
    int start, end;
    Interval(int start, int end) {
        this.start = start;
        this.end = end;
    }
}

public class MeetingRooms {
    
    public static boolean canAttendMeetings(Interval[] intervals) {
        if (intervals == null || intervals.length == 0) {
            return true;
        }
        
        // Sort by start time
        Arrays.sort(intervals, (a, b) -> a.start - b.start);
        
        // Check for overlaps
        for (int i = 1; i < intervals.length; i++) {
            if (intervals[i].start < intervals[i - 1].end) {
                return false;  // Overlap found
            }
        }
        
        return true;
    }
    
    public static void main(String[] args) {
        Interval[] meetings1 = {
            new Interval(0, 30),
            new Interval(5, 10),
            new Interval(15, 20)
        };
        System.out.println(canAttendMeetings(meetings1));  // false
        
        Interval[] meetings2 = {
            new Interval(7, 10),
            new Interval(2, 4),
            new Interval(11, 15)
        };
        System.out.println(canAttendMeetings(meetings2));  // true
    }
}
```

**Time Complexity:** O(n log n) - Sorting
**Space Complexity:** O(1)

---

### Problem 2: Top K Frequent Elements

**Problem:** Given an array of integers, return the k most frequent elements.

**Example:** `[1,1,1,2,2,3], k=2` → `[1,2]`

```java
import java.util.*;

public class TopKFrequent {
    
    public static int[] topKFrequent(int[] nums, int k) {
        // Step 1: Count frequencies
        Map<Integer, Integer> freq = new HashMap<>();
        for (int num : nums) {
            freq.put(num, freq.getOrDefault(num, 0) + 1);
        }
        
        // Step 2: Use min heap with custom comparator
        // Keep only k elements, min frequency at top
        PriorityQueue<Integer> minHeap = new PriorityQueue<>(
            (a, b) -> freq.get(a) - freq.get(b)
        );
        
        for (int num : freq.keySet()) {
            minHeap.offer(num);
            if (minHeap.size() > k) {
                minHeap.poll();  // Remove least frequent
            }
        }
        
        // Step 3: Extract result
        int[] result = new int[k];
        for (int i = 0; i < k; i++) {
            result[i] = minHeap.poll();
        }
        
        return result;
    }
    
    public static void main(String[] args) {
        int[] nums = {1, 1, 1, 2, 2, 3};
        int k = 2;
        int[] result = topKFrequent(nums, k);
        System.out.println(Arrays.toString(result));  // [2, 1] or [1, 2]
    }
}
```

**Time Complexity:** O(n log k) - Heap operations
**Space Complexity:** O(n) - HashMap

---

### Problem 3: Merge Intervals

**Problem:** Given intervals, merge all overlapping intervals.

**Example:** `[[1,3],[2,6],[8,10],[15,18]]` → `[[1,6],[8,10],[15,18]]`

```java
import java.util.*;

public class MergeIntervals {
    
    public static int[][] merge(int[][] intervals) {
        if (intervals.length <= 1) {
            return intervals;
        }
        
        // Sort by start time
        Arrays.sort(intervals, (a, b) -> a[0] - b[0]);
        
        List<int[]> merged = new ArrayList<>();
        int[] current = intervals[0];
        
        for (int i = 1; i < intervals.length; i++) {
            if (current[1] >= intervals[i][0]) {
                // Overlapping - merge
                current[1] = Math.max(current[1], intervals[i][1]);
            } else {
                // No overlap - add current and move to next
                merged.add(current);
                current = intervals[i];
            }
        }
        
        // Don't forget last interval
        merged.add(current);
        
        return merged.toArray(new int[merged.size()][]);
    }
    
    public static void main(String[] args) {
        int[][] intervals = {
            {1, 3},
            {2, 6},
            {8, 10},
            {15, 18}
        };
        
        int[][] result = merge(intervals);
        for (int[] interval : result) {
            System.out.println(Arrays.toString(interval));
        }
        // Output: [1, 6], [8, 10], [15, 18]
    }
}
```

**Time Complexity:** O(n log n) - Sorting
**Space Complexity:** O(n) - Result list

---

### Problem 4: Largest Number

**Problem:** Given a list of non-negative integers, arrange them to form the largest number.

**Example:** `[3, 30, 34, 5, 9]` → `"9534330"`

```java
import java.util.*;

public class LargestNumber {
    
    public static String largestNumber(int[] nums) {
        // Convert to strings
        String[] strs = new String[nums.length];
        for (int i = 0; i < nums.length; i++) {
            strs[i] = String.valueOf(nums[i]);
        }
        
        // Custom comparator: compare concatenated results
        // If "a+b" > "b+a", then a should come before b
        Arrays.sort(strs, (a, b) -> (b + a).compareTo(a + b));
        
        // Edge case: all zeros
        if (strs[0].equals("0")) {
            return "0";
        }
        
        // Build result
        StringBuilder sb = new StringBuilder();
        for (String s : strs) {
            sb.append(s);
        }
        
        return sb.toString();
    }
    
    public static void main(String[] args) {
        int[] nums1 = {3, 30, 34, 5, 9};
        System.out.println(largestNumber(nums1));  // "9534330"
        
        int[] nums2 = {10, 2};
        System.out.println(largestNumber(nums2));  // "210"
        
        int[] nums3 = {0, 0};
        System.out.println(largestNumber(nums3));  // "0"
    }
}
```

**Time Complexity:** O(n log n) - Sorting with string comparison
**Space Complexity:** O(n) - String array

**Why this comparator works:**
- For `"3"` and `"30"`: `"330"` vs `"303"` → `"330"` is larger, so `"3"` comes first
- For `"9"` and `"5"`: `"95"` vs `"59"` → `"95"` is larger, so `"9"` comes first

---

### Problem 5: Sort Colors (Dutch National Flag)

**Problem:** Sort array with only 0s, 1s, and 2s in-place.

**Example:** `[2,0,2,1,1,0]` → `[0,0,1,1,2,2]`

```java
import java.util.*;

public class SortColors {
    
    // Three-way partitioning (Dutch National Flag Algorithm)
    public static void sortColors(int[] nums) {
        int low = 0;      // Boundary for 0s
        int mid = 0;      // Current element
        int high = nums.length - 1;  // Boundary for 2s
        
        while (mid <= high) {
            if (nums[mid] == 0) {
                // Swap with low boundary
                swap(nums, low, mid);
                low++;
                mid++;
            } else if (nums[mid] == 1) {
                // Already in correct place
                mid++;
            } else {  // nums[mid] == 2
                // Swap with high boundary
                swap(nums, mid, high);
                high--;
                // Don't increment mid (need to check swapped element)
            }
        }
    }
    
    private static void swap(int[] nums, int i, int j) {
        int temp = nums[i];
        nums[i] = nums[j];
        nums[j] = temp;
    }
    
    // Alternative: Using counting sort
    public static void sortColorsCount(int[] nums) {
        int[] count = new int[3];
        
        // Count occurrences
        for (int num : nums) {
            count[num]++;
        }
        
        // Overwrite array
        int index = 0;
        for (int i = 0; i < 3; i++) {
            for (int j = 0; j < count[i]; j++) {
                nums[index++] = i;
            }
        }
    }
    
    public static void main(String[] args) {
        int[] nums = {2, 0, 2, 1, 1, 0};
        sortColors(nums);
        System.out.println(Arrays.toString(nums));  // [0, 0, 1, 1, 2, 2]
    }
}
```

**Time Complexity:** O(n) - Single pass
**Space Complexity:** O(1) - In-place

---

### Problem 6: Kth Largest Element in Array

**Problem:** Find the kth largest element in an unsorted array.

**Example:** `[3,2,1,5,6,4], k=2` → `5`

```java
import java.util.*;

public class KthLargest {
    
    // Method 1: Min Heap (Best for streaming data)
    public static int findKthLargest(int[] nums, int k) {
        // Min heap of size k
        PriorityQueue<Integer> minHeap = new PriorityQueue<>();
        
        for (int num : nums) {
            minHeap.offer(num);
            if (minHeap.size() > k) {
                minHeap.poll();  // Remove smallest
            }
        }
        
        return minHeap.peek();  // Kth largest at top
    }
    
    // Method 2: Max Heap (Simple but less efficient)
    public static int findKthLargestMaxHeap(int[] nums, int k) {
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>(
            Collections.reverseOrder()
        );
        
        for (int num : nums) {
            maxHeap.offer(num);
        }
        
        // Poll k-1 times
        for (int i = 0; i < k - 1; i++) {
            maxHeap.poll();
        }
        
        return maxHeap.poll();
    }
    
    // Method 3: Sorting (Simplest)
    public static int findKthLargestSort(int[] nums, int k) {
        Arrays.sort(nums);
        return nums[nums.length - k];
    }
    
    // Method 4: QuickSelect (Most efficient - O(n) average)
    public static int findKthLargestQuickSelect(int[] nums, int k) {
        return quickSelect(nums, 0, nums.length - 1, nums.length - k);
    }
    
    private static int quickSelect(int[] nums, int left, int right, int kIndex) {
        if (left == right) return nums[left];
        
        int pivotIndex = partition(nums, left, right);
        
        if (kIndex == pivotIndex) {
            return nums[kIndex];
        } else if (kIndex < pivotIndex) {
            return quickSelect(nums, left, pivotIndex - 1, kIndex);
        } else {
            return quickSelect(nums, pivotIndex + 1, right, kIndex);
        }
    }
    
    private static int partition(int[] nums, int left, int right) {
        int pivot = nums[right];
        int i = left;
        
        for (int j = left; j < right; j++) {
            if (nums[j] <= pivot) {
                swap(nums, i, j);
                i++;
            }
        }
        
        swap(nums, i, right);
        return i;
    }
    
    private static void swap(int[] nums, int i, int j) {
        int temp = nums[i];
        nums[i] = nums[j];
        nums[j] = temp;
    }
    
    public static void main(String[] args) {
        int[] nums = {3, 2, 1, 5, 6, 4};
        int k = 2;
        
        System.out.println(findKthLargest(nums, k));  // 5
        System.out.println(findKthLargestMaxHeap(nums, k));  // 5
        System.out.println(findKthLargestSort(nums, k));  // 5
        System.out.println(findKthLargestQuickSelect(nums, k));  // 5
    }
}
```

**Time Complexity:** 
- Min Heap: O(n log k)
- Max Heap: O(n log n)
- Sorting: O(n log n)
- QuickSelect: O(n) average, O(n²) worst

**Space Complexity:** O(k) for heap, O(1) for QuickSelect

---

### Problem 7: Sort Characters by Frequency

**Problem:** Sort characters in a string by frequency (descending).

**Example:** `"tree"` → `"eert"` or `"eetr"` (both valid)

```java
import java.util.*;

public class SortCharactersByFrequency {
    
    // Method 1: Using PriorityQueue
    public static String frequencySort(String s) {
        // Count frequencies
        Map<Character, Integer> freq = new HashMap<>();
        for (char c : s.toCharArray()) {
            freq.put(c, freq.getOrDefault(c, 0) + 1);
        }
        
        // Max heap by frequency
        PriorityQueue<Character> maxHeap = new PriorityQueue<>(
            (a, b) -> freq.get(b) - freq.get(a)
        );
        maxHeap.addAll(freq.keySet());
        
        // Build result
        StringBuilder sb = new StringBuilder();
        while (!maxHeap.isEmpty()) {
            char c = maxHeap.poll();
            int count = freq.get(c);
            for (int i = 0; i < count; i++) {
                sb.append(c);
            }
        }
        
        return sb.toString();
    }
    
    // Method 2: Using Sorting
    public static String frequencySortArray(String s) {
        // Count frequencies
        Map<Character, Integer> freq = new HashMap<>();
        for (char c : s.toCharArray()) {
            freq.put(c, freq.getOrDefault(c, 0) + 1);
        }
        
        // Convert to list and sort
        List<Character> chars = new ArrayList<>(freq.keySet());
        chars.sort((a, b) -> freq.get(b) - freq.get(a));
        
        // Build result
        StringBuilder sb = new StringBuilder();
        for (char c : chars) {
            int count = freq.get(c);
            for (int i = 0; i < count; i++) {
                sb.append(c);
            }
        }
        
        return sb.toString();
    }
    
    // Method 3: Bucket Sort (Most efficient for frequency sorting)
    public static String frequencySortBucket(String s) {
        // Count frequencies
        Map<Character, Integer> freq = new HashMap<>();
        for (char c : s.toCharArray()) {
            freq.put(c, freq.getOrDefault(c, 0) + 1);
        }
        
        // Create buckets (index = frequency)
        List<Character>[] buckets = new List[s.length() + 1];
        for (char c : freq.keySet()) {
            int f = freq.get(c);
            if (buckets[f] == null) {
                buckets[f] = new ArrayList<>();
            }
            buckets[f].add(c);
        }
        
        // Build result from high frequency to low
        StringBuilder sb = new StringBuilder();
        for (int i = buckets.length - 1; i >= 0; i--) {
            if (buckets[i] != null) {
                for (char c : buckets[i]) {
                    for (int j = 0; j < i; j++) {
                        sb.append(c);
                    }
                }
            }
        }
        
        return sb.toString();
    }
    
    public static void main(String[] args) {
        String s = "tree";
        System.out.println(frequencySort(s));        // "eert" or "eetr"
        System.out.println(frequencySortArray(s));   // "eert" or "eetr"
        System.out.println(frequencySortBucket(s));  // "eert" or "eetr"
        
        String s2 = "cccaaa";
        System.out.println(frequencySort(s2));       // "cccaaa" or "aaaccc"
        
        String s3 = "Aabb";
        System.out.println(frequencySort(s3));       // "bbAa" or "bbaA"
    }
}
```

**Time Complexity:** 
- PriorityQueue: O(n log k) where k = unique characters
- Sorting: O(k log k)
- Bucket Sort: O(n)

**Space Complexity:** O(n)

---

## Comparison of Solutions

| Problem | Best Approach | Time | Space |
|---------|--------------|------|-------|
| Meeting Rooms | Sort by start | O(n log n) | O(1) |
| Top K Frequent | Min Heap | O(n log k) | O(n) |
| Merge Intervals | Sort + Merge | O(n log n) | O(n) |
| Largest Number | Custom Comparator | O(n log n) | O(n) |
| Sort Colors | Dutch Flag | O(n) | O(1) |
| Kth Largest | QuickSelect or Min Heap | O(n) avg | O(k) |
| Sort by Frequency | Bucket Sort | O(n) | O(n) |

---

## Key Takeaways

✅ **Comparable** - Natural ordering, inside class, one way
✅ **Comparator** - Custom ordering, outside class, multiple ways
✅ **Use Integer.compare()** - Avoid overflow, not subtraction
✅ **Lambda > Anonymous class** - More concise
✅ **thenComparing()** - For tie-breaking
✅ **Primitives limitation** - Can't use Comparator directly
✅ **Collections.sort() is stable** - Maintains order of equals
✅ **Practice custom comparators** - 90% of problems need it!

---

[← Back: Java 8-17 Features](./Java8-to-17-Features.md) | [Next: Arrays →](../03-Arrays/01-Array-Basics.md)
