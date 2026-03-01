# Problem-Solving Patterns - Your Competitive Programming Toolkit

## Why Patterns Matter?

In competitive programming, you'll see the **same patterns** repeated in different problems. Once you recognize the pattern, you already know 80% of the solution!

Think of patterns like **Lego blocks** - learn the blocks once, build anything!

---

## Pattern Recognition Framework

When you see a problem, ask:

1. **What's the input?** Array? String? Graph? Tree?
2. **What am I looking for?** Maximum? Minimum? Count? Subsequence?
3. **What are the constraints?** Small n? Large n? Sorted? Unsorted?
4. **Does this remind me of something?** Have I seen similar?

---

## Pattern 1: Two Pointers

**When to use:** Array/string problems, finding pairs, reversing, partitioning

**Idea:** Use two pointers moving from different positions

### Example 1: Find Pair with Target Sum (Sorted Array)

**Problem:** Given sorted array, find two numbers that add to target

```java
public class TwoPointers {
    
    // Find pair that sums to target
    public static int[] findPair(int[] arr, int target) {
        int left = 0;
        int right = arr.length - 1;
        
        while (left < right) {
            int sum = arr[left] + arr[right];
            
            if (sum == target) {
                return new int[]{left, right};
            } else if (sum < target) {
                left++;     // Need bigger sum
            } else {
                right--;    // Need smaller sum
            }
        }
        
        return new int[]{-1, -1};  // Not found
    }
    
    public static void main(String[] args) {
        int[] arr = {1, 3, 5, 7, 9, 11};
        int target = 14;
        int[] result = findPair(arr, target);
        System.out.println("Indices: " + result[0] + ", " + result[1]);
        // Output: Indices: 2, 4 (5 + 9 = 14)
    }
}
```

**Why it works:** Sorted array lets us decide which pointer to move!
**Time:** O(n), **Space:** O(1)

---

### Example 2: Reverse Array In-Place

```java
public static void reverseArray(int[] arr) {
    int left = 0;
    int right = arr.length - 1;
    
    while (left < right) {
        // Swap
        int temp = arr[left];
        arr[left] = arr[right];
        arr[right] = temp;
        
        left++;
        right--;
    }
}
```

---

### Example 3: Remove Duplicates from Sorted Array

```java
public static int removeDuplicates(int[] arr) {
    if (arr.length == 0) return 0;
    
    int slow = 0;  // Points to last unique element
    
    for (int fast = 1; fast < arr.length; fast++) {
        if (arr[fast] != arr[slow]) {
            slow++;
            arr[slow] = arr[fast];
        }
    }
    
    return slow + 1;  // Length of unique elements
}

// Example: [1,1,2,2,3,4,4] → [1,2,3,4,...]
```

**Two pointers:** `slow` (write position), `fast` (read position)

---

## Pattern 2: Sliding Window

**When to use:** Subarray/substring problems with consecutive elements

**Idea:** Maintain a window and slide it through the array

### Example 1: Maximum Sum of K Consecutive Elements

**Problem:** Find maximum sum of any k consecutive elements

```java
public class SlidingWindow {
    
    // Maximum sum of k consecutive elements
    public static int maxSumSubarray(int[] arr, int k) {
        int n = arr.length;
        if (n < k) return -1;
        
        // Calculate sum of first window
        int windowSum = 0;
        for (int i = 0; i < k; i++) {
            windowSum += arr[i];
        }
        
        int maxSum = windowSum;
        
        // Slide the window
        for (int i = k; i < n; i++) {
            windowSum += arr[i] - arr[i - k];  // Add new, remove old
            maxSum = Math.max(maxSum, windowSum);
        }
        
        return maxSum;
    }
    
    public static void main(String[] args) {
        int[] arr = {2, 1, 5, 1, 3, 2};
        int k = 3;
        System.out.println("Max sum: " + maxSumSubarray(arr, k));
        // Output: Max sum: 9 (5+1+3)
    }
}
```

**Without sliding window:** O(n×k) - recalculate each window
**With sliding window:** O(n) - reuse previous calculation!

---

### Example 2: Longest Substring Without Repeating Characters

```java
public static int lengthOfLongestSubstring(String s) {
    HashMap<Character, Integer> map = new HashMap<>();
    int maxLength = 0;
    int start = 0;
    
    for (int end = 0; end < s.length(); end++) {
        char c = s.charAt(end);
        
        // If character already in window, move start
        if (map.containsKey(c)) {
            start = Math.max(start, map.get(c) + 1);
        }
        
        map.put(c, end);
        maxLength = Math.max(maxLength, end - start + 1);
    }
    
    return maxLength;
}

// Example: "abcabcbb" → 3 ("abc")
```

**Window expands** with `end`, **contracts** when duplicate found

---

## Pattern 3: Fast & Slow Pointers (Floyd's Cycle Detection)

**When to use:** Linked list problems, cycle detection, finding middle

**Idea:** Two pointers moving at different speeds

### Example 1: Detect Cycle in Linked List

```java
class ListNode {
    int val;
    ListNode next;
    ListNode(int val) { this.val = val; }
}

public class FastSlowPointers {
    
    public static boolean hasCycle(ListNode head) {
        if (head == null) return false;
        
        ListNode slow = head;
        ListNode fast = head;
        
        while (fast != null && fast.next != null) {
            slow = slow.next;           // Move 1 step
            fast = fast.next.next;      // Move 2 steps
            
            if (slow == fast) {
                return true;  // Cycle detected!
            }
        }
        
        return false;  // No cycle
    }
}
```

**Why it works:** If there's a cycle, fast will eventually catch slow (like runners on a track)!

---

### Example 2: Find Middle of Linked List

```java
public static ListNode findMiddle(ListNode head) {
    ListNode slow = head;
    ListNode fast = head;
    
    while (fast != null && fast.next != null) {
        slow = slow.next;
        fast = fast.next.next;
    }
    
    return slow;  // When fast reaches end, slow is at middle
}
```

---

## Pattern 4: Prefix Sum

**When to use:** Range sum queries, subarray sum problems

**Idea:** Precompute cumulative sums for O(1) range queries

### Example: Range Sum Query

```java
public class PrefixSum {
    
    private int[] prefix;
    
    public PrefixSum(int[] arr) {
        int n = arr.length;
        prefix = new int[n + 1];
        
        // Build prefix sum array
        for (int i = 0; i < n; i++) {
            prefix[i + 1] = prefix[i] + arr[i];
        }
    }
    
    // Get sum from index left to right (inclusive)
    public int rangeSum(int left, int right) {
        return prefix[right + 1] - prefix[left];
    }
    
    public static void main(String[] args) {
        int[] arr = {1, 2, 3, 4, 5};
        PrefixSum ps = new PrefixSum(arr);
        
        System.out.println(ps.rangeSum(1, 3));  // Sum of [2,3,4] = 9
        System.out.println(ps.rangeSum(0, 4));  // Sum of [1,2,3,4,5] = 15
    }
}
```

**Without prefix:** O(n) per query
**With prefix:** O(1) per query (after O(n) preprocessing)

---

## Pattern 5: Hashing (Frequency Counter)

**When to use:** Counting occurrences, finding duplicates, anagrams

**Idea:** Use HashMap to count/track elements

### Example 1: First Non-Repeating Character

```java
public static char firstUnique(String s) {
    HashMap<Character, Integer> freq = new HashMap<>();
    
    // Count frequencies
    for (char c : s.toCharArray()) {
        freq.put(c, freq.getOrDefault(c, 0) + 1);
    }
    
    // Find first with frequency 1
    for (char c : s.toCharArray()) {
        if (freq.get(c) == 1) {
            return c;
        }
    }
    
    return '_';  // None found
}

// Example: "leetcode" → 'l'
```

---

### Example 2: Check if Two Strings are Anagrams

```java
public static boolean areAnagrams(String s1, String s2) {
    if (s1.length() != s2.length()) return false;
    
    HashMap<Character, Integer> freq = new HashMap<>();
    
    // Add characters from s1
    for (char c : s1.toCharArray()) {
        freq.put(c, freq.getOrDefault(c, 0) + 1);
    }
    
    // Subtract characters from s2
    for (char c : s2.toCharArray()) {
        if (!freq.containsKey(c)) return false;
        freq.put(c, freq.get(c) - 1);
        if (freq.get(c) == 0) {
            freq.remove(c);
        }
    }
    
    return freq.isEmpty();
}

// Example: "listen", "silent" → true
```

---

## Pattern 6: Sorting First

**When to use:** When problem becomes easier with sorted data

**Idea:** Sort first (O(n log n)), then solve (often O(n))

### Example: Merge Intervals

```java
import java.util.*;

class Interval {
    int start, end;
    Interval(int s, int e) { start = s; end = e; }
}

public class MergeIntervals {
    
    public static List<Interval> merge(List<Interval> intervals) {
        if (intervals.isEmpty()) return intervals;
        
        // Sort by start time
        intervals.sort((a, b) -> a.start - b.start);
        
        List<Interval> merged = new ArrayList<>();
        Interval current = intervals.get(0);
        
        for (int i = 1; i < intervals.size(); i++) {
            Interval next = intervals.get(i);
            
            if (current.end >= next.start) {
                // Merge overlapping intervals
                current.end = Math.max(current.end, next.end);
            } else {
                // No overlap, add current and move to next
                merged.add(current);
                current = next;
            }
        }
        
        merged.add(current);  // Don't forget last interval
        return merged;
    }
}

// Example: [[1,3], [2,6], [8,10], [15,18]]
// Output: [[1,6], [8,10], [15,18]]
```

---

## Pattern 7: Greedy Approach

**When to use:** Optimization problems where local optimal → global optimal

**Idea:** Make the best choice at each step

### Example: Activity Selection

```java
class Activity {
    int start, end;
    Activity(int s, int e) { start = s; end = e; }
}

public static int maxActivities(List<Activity> activities) {
    // Sort by end time (finish earliest activities first)
    activities.sort((a, b) -> a.end - b.end);
    
    int count = 1;
    int lastEnd = activities.get(0).end;
    
    for (int i = 1; i < activities.size(); i++) {
        if (activities.get(i).start >= lastEnd) {
            count++;
            lastEnd = activities.get(i).end;
        }
    }
    
    return count;
}
```

**Greedy choice:** Always pick activity that finishes earliest!

---

## Pattern Recognition Cheat Sheet

| Pattern | Problem Indicators | Time | Space |
|---------|-------------------|------|-------|
| **Two Pointers** | Sorted array, find pair, reverse | O(n) | O(1) |
| **Sliding Window** | Consecutive subarray, "k elements" | O(n) | O(1) |
| **Fast & Slow** | Linked list, cycle detection | O(n) | O(1) |
| **Prefix Sum** | Range queries, subarray sum | O(1) query | O(n) |
| **Hashing** | Count frequency, duplicates | O(n) | O(n) |
| **Sorting First** | Intervals, pairs, anagrams | O(n log n) | O(1) or O(n) |
| **Binary Search** | Sorted data, "find in log time" | O(log n) | O(1) |
| **Recursion/Backtracking** | Generate all, permutations | O(2ⁿ) or O(n!) | O(n) |

---

## Practice: Identify the Pattern!

### Problem 1
**Find if array has duplicates**
<details>
<summary>Pattern</summary>

**Hashing** - Use HashSet to track seen elements
</details>

---

### Problem 2
**Find maximum product of 3 numbers in array**
<details>
<summary>Pattern</summary>

**Sorting First** - Sort and check largest 3 vs (smallest 2 × largest 1)
</details>

---

### Problem 3
**Find longest subarray with sum ≤ k**
<details>
<summary>Pattern</summary>

**Sliding Window** - Expand window while sum ≤ k, contract when > k
</details>

---

### Problem 4
**Check if linked list is palindrome**
<details>
<summary>Pattern</summary>

**Fast & Slow Pointers** - Find middle, reverse second half, compare
</details>

---

## Complete Example: Applying Multiple Patterns

**Problem:** Find the pair in array with smallest difference

```java
public class SmallestDifference {
    
    public static int[] findPair(int[] arr) {
        // Pattern 1: Sort first
        Arrays.sort(arr);
        
        int minDiff = Integer.MAX_VALUE;
        int[] result = new int[2];
        
        // Pattern 2: Two pointers (adjacent elements)
        for (int i = 0; i < arr.length - 1; i++) {
            int diff = arr[i + 1] - arr[i];
            if (diff < minDiff) {
                minDiff = diff;
                result[0] = arr[i];
                result[1] = arr[i + 1];
            }
        }
        
        return result;
    }
    
    public static void main(String[] args) {
        int[] arr = {10, 5, 3, 9, 2, 8};
        int[] pair = findPair(arr);
        System.out.println("Pair: " + pair[0] + ", " + pair[1]);
        // Output: Pair: 8, 9 (difference = 1)
    }
}
```

---

## Next Steps

Now that you know the patterns, it's time to practice!

**Practice Strategy:**
1. Identify the pattern BEFORE coding
2. Write pseudocode first
3. Implement and test
4. Analyze time/space complexity

**Ready to apply these patterns?** Let's move to Arrays!

[← Back: Big O Notation](./BigO-Notation.md) | [Next: Arrays →](../03-Arrays/01-Array-Basics.md)
