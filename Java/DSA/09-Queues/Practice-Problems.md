# Queue Practice Problems

## Table of Contents
1. [Easy Problems](#easy-problems)
2. [Medium Problems](#medium-problems)
3. [Hard Problems](#hard-problems)
4. [Topic-Wise Problems](#topic-wise-problems)
5. [Company-Wise Problems](#company-wise-problems)
6. [Problem Solutions](#problem-solutions)

---

## Easy Problems

### 1. Implement Queue using Stacks ⭐⭐⭐⭐⭐
**LeetCode 232** | **Difficulty: Easy** | **Must Solve**

**Problem**: Implement a FIFO queue using only two stacks.

**Operations**:
- `void push(int x)`: Pushes element x to back
- `int pop()`: Removes element from front
- `int peek()`: Gets front element
- `boolean empty()`: Returns true if empty

**Solution**:
```java
class MyQueue {
    private Deque<Integer> s1;  // For push
    private Deque<Integer> s2;  // For pop
    
    public MyQueue() {
        s1 = new ArrayDeque<>();
        s2 = new ArrayDeque<>();
    }
    
    public void push(int x) {
        s1.push(x);
    }
    
    public int pop() {
        peek();  // Ensure s2 has elements
        return s2.pop();
    }
    
    public int peek() {
        if (s2.isEmpty()) {
            while (!s1.isEmpty()) {
                s2.push(s1.pop());
            }
        }
        return s2.peek();
    }
    
    public boolean empty() {
        return s1.isEmpty() && s2.isEmpty();
    }
}
```

**Time**: Push O(1), Pop/Peek O(1) amortized
**Space**: O(n)

---

### 2. Implement Stack using Queues ⭐⭐⭐⭐
**LeetCode 225** | **Difficulty: Easy**

**Problem**: Implement a LIFO stack using only queues.

**Solution 1: One Queue (Push O(n))**:
```java
class MyStack {
    private Queue<Integer> queue;
    
    public MyStack() {
        queue = new LinkedList<>();
    }
    
    public void push(int x) {
        queue.offer(x);
        int size = queue.size();
        
        // Rotate queue to make new element front
        for (int i = 1; i < size; i++) {
            queue.offer(queue.poll());
        }
    }
    
    public int pop() {
        return queue.poll();
    }
    
    public int top() {
        return queue.peek();
    }
    
    public boolean empty() {
        return queue.isEmpty();
    }
}
```

**Time**: Push O(n), Pop/Top O(1)
**Space**: O(n)

---

### 3. Design Circular Queue ⭐⭐⭐⭐⭐
**LeetCode 622** | **Difficulty: Medium** | **Must Solve**

**Problem**: Design circular queue with fixed size.

**Solution**:
```java
class MyCircularQueue {
    private int[] arr;
    private int front;
    private int rear;
    private int size;
    private int capacity;
    
    public MyCircularQueue(int k) {
        this.capacity = k;
        this.arr = new int[k];
        this.front = 0;
        this.rear = -1;
        this.size = 0;
    }
    
    public boolean enQueue(int value) {
        if (isFull()) return false;
        
        rear = (rear + 1) % capacity;
        arr[rear] = value;
        size++;
        return true;
    }
    
    public boolean deQueue() {
        if (isEmpty()) return false;
        
        front = (front + 1) % capacity;
        size--;
        return true;
    }
    
    public int Front() {
        return isEmpty() ? -1 : arr[front];
    }
    
    public int Rear() {
        return isEmpty() ? -1 : arr[rear];
    }
    
    public boolean isEmpty() {
        return size == 0;
    }
    
    public boolean isFull() {
        return size == capacity;
    }
}
```

**Time**: All operations O(1)
**Space**: O(k)

---

### 4. Number of Recent Calls ⭐⭐⭐
**LeetCode 933** | **Difficulty: Easy**

**Problem**: Count requests in last 3000 milliseconds.

**Solution**:
```java
class RecentCounter {
    private Queue<Integer> queue;
    
    public RecentCounter() {
        queue = new LinkedList<>();
    }
    
    public int ping(int t) {
        queue.offer(t);
        
        // Remove requests outside 3000ms window
        while (queue.peek() < t - 3000) {
            queue.poll();
        }
        
        return queue.size();
    }
}
```

**Time**: O(1) amortized
**Space**: O(n)

---

### 5. Moving Average from Data Stream ⭐⭐⭐
**LeetCode 346** | **Difficulty: Easy**

**Problem**: Calculate moving average of last k values.

**Solution**:
```java
class MovingAverage {
    private Queue<Integer> queue;
    private int size;
    private double sum;
    
    public MovingAverage(int size) {
        this.queue = new LinkedList<>();
        this.size = size;
        this.sum = 0;
    }
    
    public double next(int val) {
        queue.offer(val);
        sum += val;
        
        if (queue.size() > size) {
            sum -= queue.poll();
        }
        
        return sum / queue.size();
    }
}
```

**Time**: O(1)
**Space**: O(k)

---

### 6. Last Stone Weight ⭐⭐⭐
**LeetCode 1046** | **Difficulty: Easy**

**Problem**: Smash two heaviest stones until ≤1 remains.

**Solution**:
```java
public class LastStoneWeight {
    public int lastStoneWeight(int[] stones) {
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
        
        for (int stone : stones) {
            maxHeap.offer(stone);
        }
        
        while (maxHeap.size() > 1) {
            int first = maxHeap.poll();
            int second = maxHeap.poll();
            
            if (first != second) {
                maxHeap.offer(first - second);
            }
        }
        
        return maxHeap.isEmpty() ? 0 : maxHeap.peek();
    }
}
```

**Time**: O(n log n)
**Space**: O(n)

---

### 7. Kth Largest Element in Stream ⭐⭐⭐⭐
**LeetCode 703** | **Difficulty: Easy**

**Problem**: Design class to find kth largest in stream.

**Solution**:
```java
class KthLargest {
    private PriorityQueue<Integer> minHeap;
    private int k;
    
    public KthLargest(int k, int[] nums) {
        this.k = k;
        this.minHeap = new PriorityQueue<>();
        
        for (int num : nums) {
            add(num);
        }
    }
    
    public int add(int val) {
        minHeap.offer(val);
        
        if (minHeap.size() > k) {
            minHeap.poll();
        }
        
        return minHeap.peek();
    }
}
```

**Time**: O(log k) per add
**Space**: O(k)

---

## Medium Problems

### 8. Design Circular Deque ⭐⭐⭐⭐⭐
**LeetCode 641** | **Difficulty: Medium** | **Must Solve**

**Problem**: Design double-ended circular queue.

**Solution**:
```java
class MyCircularDeque {
    private int[] arr;
    private int front;
    private int rear;
    private int size;
    private int capacity;
    
    public MyCircularDeque(int k) {
        this.capacity = k;
        this.arr = new int[k];
        this.front = 0;
        this.rear = k - 1;
        this.size = 0;
    }
    
    public boolean insertFront(int value) {
        if (isFull()) return false;
        
        front = (front - 1 + capacity) % capacity;
        arr[front] = value;
        size++;
        return true;
    }
    
    public boolean insertLast(int value) {
        if (isFull()) return false;
        
        rear = (rear + 1) % capacity;
        arr[rear] = value;
        size++;
        return true;
    }
    
    public boolean deleteFront() {
        if (isEmpty()) return false;
        
        front = (front + 1) % capacity;
        size--;
        return true;
    }
    
    public boolean deleteLast() {
        if (isEmpty()) return false;
        
        rear = (rear - 1 + capacity) % capacity;
        size--;
        return true;
    }
    
    public int getFront() {
        return isEmpty() ? -1 : arr[front];
    }
    
    public int getRear() {
        return isEmpty() ? -1 : arr[rear];
    }
    
    public boolean isEmpty() {
        return size == 0;
    }
    
    public boolean isFull() {
        return size == capacity;
    }
}
```

**Time**: All operations O(1)
**Space**: O(k)

---

### 9. Binary Tree Level Order Traversal ⭐⭐⭐⭐⭐
**LeetCode 102** | **Difficulty: Medium** | **Must Solve**

**Problem**: Return level order traversal of binary tree.

**Solution**:
```java
public class LevelOrder {
    public List<List<Integer>> levelOrder(TreeNode root) {
        List<List<Integer>> result = new ArrayList<>();
        if (root == null) return result;
        
        Queue<TreeNode> queue = new LinkedList<>();
        queue.offer(root);
        
        while (!queue.isEmpty()) {
            int levelSize = queue.size();
            List<Integer> level = new ArrayList<>();
            
            for (int i = 0; i < levelSize; i++) {
                TreeNode node = queue.poll();
                level.add(node.val);
                
                if (node.left != null) queue.offer(node.left);
                if (node.right != null) queue.offer(node.right);
            }
            
            result.add(level);
        }
        
        return result;
    }
}
```

**Time**: O(n)
**Space**: O(w) where w is max width

---

### 10. Sliding Window Maximum ⭐⭐⭐⭐⭐
**LeetCode 239** | **Difficulty: Hard** | **Must Solve**

**Problem**: Find max in each sliding window of size k.

**Solution**:
```java
public class SlidingWindowMaximum {
    public int[] maxSlidingWindow(int[] nums, int k) {
        int n = nums.length;
        int[] result = new int[n - k + 1];
        Deque<Integer> deque = new ArrayDeque<>();
        
        for (int i = 0; i < n; i++) {
            // Remove outside window
            while (!deque.isEmpty() && deque.peekFirst() < i - k + 1) {
                deque.pollFirst();
            }
            
            // Remove smaller elements
            while (!deque.isEmpty() && nums[deque.peekLast()] < nums[i]) {
                deque.pollLast();
            }
            
            deque.offerLast(i);
            
            if (i >= k - 1) {
                result[i - k + 1] = nums[deque.peekFirst()];
            }
        }
        
        return result;
    }
}
```

**Time**: O(n)
**Space**: O(k)

---

### 11. Kth Largest Element in Array ⭐⭐⭐⭐⭐
**LeetCode 215** | **Difficulty: Medium** | **Must Solve**

**Problem**: Find kth largest element.

**Solution**:
```java
public class KthLargest {
    public int findKthLargest(int[] nums, int k) {
        PriorityQueue<Integer> minHeap = new PriorityQueue<>();
        
        for (int num : nums) {
            minHeap.offer(num);
            
            if (minHeap.size() > k) {
                minHeap.poll();
            }
        }
        
        return minHeap.peek();
    }
}
```

**Time**: O(n log k)
**Space**: O(k)

---

### 12. Top K Frequent Elements ⭐⭐⭐⭐⭐
**LeetCode 347** | **Difficulty: Medium** | **Must Solve**

**Problem**: Find k most frequent elements.

**Solution**:
```java
public class TopKFrequent {
    public int[] topKFrequent(int[] nums, int k) {
        Map<Integer, Integer> freqMap = new HashMap<>();
        for (int num : nums) {
            freqMap.put(num, freqMap.getOrDefault(num, 0) + 1);
        }
        
        PriorityQueue<Integer> minHeap = new PriorityQueue<>(
            (a, b) -> freqMap.get(a) - freqMap.get(b)
        );
        
        for (int num : freqMap.keySet()) {
            minHeap.offer(num);
            
            if (minHeap.size() > k) {
                minHeap.poll();
            }
        }
        
        int[] result = new int[k];
        for (int i = 0; i < k; i++) {
            result[i] = minHeap.poll();
        }
        
        return result;
    }
}
```

**Time**: O(n log k)
**Space**: O(n)

---

### 13. Rotting Oranges ⭐⭐⭐⭐
**LeetCode 994** | **Difficulty: Medium**

**Problem**: Find minimum time for all oranges to rot.

**Solution**:
```java
public class RottingOranges {
    public int orangesRotting(int[][] grid) {
        int m = grid.length;
        int n = grid[0].length;
        Queue<int[]> queue = new LinkedList<>();
        int fresh = 0;
        
        // Add all rotten oranges to queue
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == 2) {
                    queue.offer(new int[]{i, j});
                } else if (grid[i][j] == 1) {
                    fresh++;
                }
            }
        }
        
        if (fresh == 0) return 0;
        
        int[][] dirs = {{0, 1}, {1, 0}, {0, -1}, {-1, 0}};
        int minutes = 0;
        
        while (!queue.isEmpty()) {
            int size = queue.size();
            
            for (int i = 0; i < size; i++) {
                int[] cell = queue.poll();
                int row = cell[0];
                int col = cell[1];
                
                for (int[] dir : dirs) {
                    int newRow = row + dir[0];
                    int newCol = col + dir[1];
                    
                    if (newRow >= 0 && newRow < m && newCol >= 0 && newCol < n
                        && grid[newRow][newCol] == 1) {
                        grid[newRow][newCol] = 2;
                        queue.offer(new int[]{newRow, newCol});
                        fresh--;
                    }
                }
            }
            
            minutes++;
        }
        
        return fresh == 0 ? minutes - 1 : -1;
    }
}
```

**Time**: O(m × n)
**Space**: O(m × n)

---

### 14. K Closest Points to Origin ⭐⭐⭐⭐
**LeetCode 973** | **Difficulty: Medium**

**Problem**: Find k closest points to origin.

**Solution**:
```java
public class KClosestPoints {
    public int[][] kClosest(int[][] points, int k) {
        PriorityQueue<int[]> maxHeap = new PriorityQueue<>(
            (a, b) -> (b[0]*b[0] + b[1]*b[1]) - (a[0]*a[0] + a[1]*a[1])
        );
        
        for (int[] point : points) {
            maxHeap.offer(point);
            
            if (maxHeap.size() > k) {
                maxHeap.poll();
            }
        }
        
        int[][] result = new int[k][2];
        for (int i = 0; i < k; i++) {
            result[i] = maxHeap.poll();
        }
        
        return result;
    }
}
```

**Time**: O(n log k)
**Space**: O(k)

---

### 15. Reorganize String ⭐⭐⭐⭐
**LeetCode 767** | **Difficulty: Medium**

**Problem**: Rearrange string so no adjacent chars are same.

**Solution**:
```java
public class ReorganizeString {
    public String reorganizeString(String s) {
        Map<Character, Integer> freqMap = new HashMap<>();
        for (char c : s.toCharArray()) {
            freqMap.put(c, freqMap.getOrDefault(c, 0) + 1);
        }
        
        PriorityQueue<Character> maxHeap = new PriorityQueue<>(
            (a, b) -> freqMap.get(b) - freqMap.get(a)
        );
        maxHeap.addAll(freqMap.keySet());
        
        StringBuilder result = new StringBuilder();
        
        while (maxHeap.size() >= 2) {
            char first = maxHeap.poll();
            char second = maxHeap.poll();
            
            result.append(first);
            result.append(second);
            
            freqMap.put(first, freqMap.get(first) - 1);
            freqMap.put(second, freqMap.get(second) - 1);
            
            if (freqMap.get(first) > 0) maxHeap.offer(first);
            if (freqMap.get(second) > 0) maxHeap.offer(second);
        }
        
        if (!maxHeap.isEmpty()) {
            char last = maxHeap.poll();
            if (freqMap.get(last) > 1) return "";
            result.append(last);
        }
        
        return result.toString();
    }
}
```

**Time**: O(n log k)
**Space**: O(k)

---

### 16. Task Scheduler ⭐⭐⭐⭐
**LeetCode 621** | **Difficulty: Medium**

**Problem**: Schedule tasks with cooling interval.

**Solution**:
```java
public class TaskScheduler {
    public int leastInterval(char[] tasks, int n) {
        int[] freq = new int[26];
        for (char task : tasks) {
            freq[task - 'A']++;
        }
        
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
        for (int f : freq) {
            if (f > 0) maxHeap.offer(f);
        }
        
        int time = 0;
        
        while (!maxHeap.isEmpty()) {
            List<Integer> temp = new ArrayList<>();
            
            for (int i = 0; i <= n; i++) {
                if (!maxHeap.isEmpty()) {
                    int count = maxHeap.poll();
                    if (count > 1) {
                        temp.add(count - 1);
                    }
                }
                time++;
                
                if (maxHeap.isEmpty() && temp.isEmpty()) {
                    break;
                }
            }
            
            for (int count : temp) {
                maxHeap.offer(count);
            }
        }
        
        return time;
    }
}
```

**Time**: O(n)
**Space**: O(1)

---

### 17. Longest Continuous Subarray With Limit ⭐⭐⭐⭐
**LeetCode 1438** | **Difficulty: Medium**

**Problem**: Find longest subarray where |max - min| ≤ limit.

**Solution**:
```java
public class LongestSubarrayWithLimit {
    public int longestSubarray(int[] nums, int limit) {
        Deque<Integer> maxDeque = new ArrayDeque<>();
        Deque<Integer> minDeque = new ArrayDeque<>();
        
        int left = 0;
        int maxLen = 0;
        
        for (int right = 0; right < nums.length; right++) {
            while (!maxDeque.isEmpty() && nums[maxDeque.peekLast()] < nums[right]) {
                maxDeque.pollLast();
            }
            maxDeque.offerLast(right);
            
            while (!minDeque.isEmpty() && nums[minDeque.peekLast()] > nums[right]) {
                minDeque.pollLast();
            }
            minDeque.offerLast(right);
            
            while (nums[maxDeque.peekFirst()] - nums[minDeque.peekFirst()] > limit) {
                left++;
                
                if (maxDeque.peekFirst() < left) maxDeque.pollFirst();
                if (minDeque.peekFirst() < left) minDeque.pollFirst();
            }
            
            maxLen = Math.max(maxLen, right - left + 1);
        }
        
        return maxLen;
    }
}
```

**Time**: O(n)
**Space**: O(n)

---

### 18. Open the Lock ⭐⭐⭐
**LeetCode 752** | **Difficulty: Medium**

**Problem**: Find minimum turns to open lock avoiding deadends.

**Solution**:
```java
public class OpenLock {
    public int openLock(String[] deadends, String target) {
        Set<String> dead = new HashSet<>(Arrays.asList(deadends));
        Set<String> visited = new HashSet<>();
        Queue<String> queue = new LinkedList<>();
        
        String start = "0000";
        if (dead.contains(start)) return -1;
        
        queue.offer(start);
        visited.add(start);
        int moves = 0;
        
        while (!queue.isEmpty()) {
            int size = queue.size();
            
            for (int i = 0; i < size; i++) {
                String current = queue.poll();
                
                if (current.equals(target)) {
                    return moves;
                }
                
                for (int j = 0; j < 4; j++) {
                    for (int delta : new int[]{-1, 1}) {
                        String next = turn(current, j, delta);
                        
                        if (!visited.contains(next) && !dead.contains(next)) {
                            queue.offer(next);
                            visited.add(next);
                        }
                    }
                }
            }
            
            moves++;
        }
        
        return -1;
    }
    
    private String turn(String s, int pos, int delta) {
        char[] arr = s.toCharArray();
        arr[pos] = (char)((arr[pos] - '0' + delta + 10) % 10 + '0');
        return new String(arr);
    }
}
```

**Time**: O(10^4)
**Space**: O(10^4)

---

### 19. Perfect Squares ⭐⭐⭐
**LeetCode 279** | **Difficulty: Medium**

**Problem**: Find minimum perfect squares that sum to n.

**Solution**:
```java
public class PerfectSquares {
    public int numSquares(int n) {
        Queue<Integer> queue = new LinkedList<>();
        Set<Integer> visited = new HashSet<>();
        
        queue.offer(n);
        visited.add(n);
        int level = 0;
        
        while (!queue.isEmpty()) {
            int size = queue.size();
            level++;
            
            for (int i = 0; i < size; i++) {
                int current = queue.poll();
                
                for (int j = 1; j * j <= current; j++) {
                    int next = current - j * j;
                    
                    if (next == 0) {
                        return level;
                    }
                    
                    if (!visited.contains(next)) {
                        queue.offer(next);
                        visited.add(next);
                    }
                }
            }
        }
        
        return level;
    }
}
```

**Time**: O(n√n)
**Space**: O(n)

---

### 20. Jump Game VI ⭐⭐⭐
**LeetCode 1696** | **Difficulty: Medium**

**Problem**: Maximum score jumping at most k steps.

**Solution**:
```java
public class JumpGameVI {
    public int maxResult(int[] nums, int k) {
        int n = nums.length;
        int[] dp = new int[n];
        dp[0] = nums[0];
        
        Deque<Integer> deque = new ArrayDeque<>();
        deque.offer(0);
        
        for (int i = 1; i < n; i++) {
            while (!deque.isEmpty() && deque.peekFirst() < i - k) {
                deque.pollFirst();
            }
            
            dp[i] = nums[i] + dp[deque.peekFirst()];
            
            while (!deque.isEmpty() && dp[deque.peekLast()] <= dp[i]) {
                deque.pollLast();
            }
            
            deque.offerLast(i);
        }
        
        return dp[n - 1];
    }
}
```

**Time**: O(n)
**Space**: O(n)

---

## Hard Problems

### 21. Find Median from Data Stream ⭐⭐⭐⭐⭐
**LeetCode 295** | **Difficulty: Hard** | **Must Solve**

**Problem**: Support addNum and findMedian operations.

**Solution**:
```java
class MedianFinder {
    private PriorityQueue<Integer> maxHeap;  // Lower half
    private PriorityQueue<Integer> minHeap;  // Upper half
    
    public MedianFinder() {
        maxHeap = new PriorityQueue<>(Collections.reverseOrder());
        minHeap = new PriorityQueue<>();
    }
    
    public void addNum(int num) {
        maxHeap.offer(num);
        minHeap.offer(maxHeap.poll());
        
        if (maxHeap.size() < minHeap.size()) {
            maxHeap.offer(minHeap.poll());
        }
    }
    
    public double findMedian() {
        if (maxHeap.size() > minHeap.size()) {
            return maxHeap.peek();
        } else {
            return (maxHeap.peek() + minHeap.peek()) / 2.0;
        }
    }
}
```

**Time**: addNum O(log n), findMedian O(1)
**Space**: O(n)

---

### 22. Merge K Sorted Lists ⭐⭐⭐⭐⭐
**LeetCode 23** | **Difficulty: Hard** | **Must Solve**

**Problem**: Merge k sorted linked lists.

**Solution**:
```java
public class MergeKSortedLists {
    public ListNode mergeKLists(ListNode[] lists) {
        if (lists == null || lists.length == 0) return null;
        
        PriorityQueue<ListNode> minHeap = new PriorityQueue<>(
            (a, b) -> a.val - b.val
        );
        
        for (ListNode head : lists) {
            if (head != null) {
                minHeap.offer(head);
            }
        }
        
        ListNode dummy = new ListNode(0);
        ListNode current = dummy;
        
        while (!minHeap.isEmpty()) {
            ListNode node = minHeap.poll();
            current.next = node;
            current = current.next;
            
            if (node.next != null) {
                minHeap.offer(node.next);
            }
        }
        
        return dummy.next;
    }
}
```

**Time**: O(n log k)
**Space**: O(k)

---

### 23. Shortest Subarray with Sum at Least K ⭐⭐⭐⭐⭐
**LeetCode 862** | **Difficulty: Hard** | **Must Solve**

**Problem**: Find shortest subarray with sum ≥ k.

**Solution**:
```java
public class ShortestSubarrayWithSum {
    public int shortestSubarray(int[] nums, int k) {
        int n = nums.length;
        long[] prefixSum = new long[n + 1];
        
        for (int i = 0; i < n; i++) {
            prefixSum[i + 1] = prefixSum[i] + nums[i];
        }
        
        Deque<Integer> deque = new ArrayDeque<>();
        int minLen = Integer.MAX_VALUE;
        
        for (int i = 0; i <= n; i++) {
            while (!deque.isEmpty() && prefixSum[i] - prefixSum[deque.peekFirst()] >= k) {
                minLen = Math.min(minLen, i - deque.pollFirst());
            }
            
            while (!deque.isEmpty() && prefixSum[i] <= prefixSum[deque.peekLast()]) {
                deque.pollLast();
            }
            
            deque.offerLast(i);
        }
        
        return minLen == Integer.MAX_VALUE ? -1 : minLen;
    }
}
```

**Time**: O(n)
**Space**: O(n)

---

### 24. Sliding Window Median ⭐⭐⭐⭐
**LeetCode 480** | **Difficulty: Hard**

**Problem**: Find median in each sliding window.

**Solution**:
```java
public class SlidingWindowMedian {
    public double[] medianSlidingWindow(int[] nums, int k) {
        double[] result = new double[nums.length - k + 1];
        
        PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
        PriorityQueue<Integer> minHeap = new PriorityQueue<>();
        
        for (int i = 0; i < nums.length; i++) {
            if (maxHeap.isEmpty() || nums[i] <= maxHeap.peek()) {
                maxHeap.offer(nums[i]);
            } else {
                minHeap.offer(nums[i]);
            }
            
            if (maxHeap.size() > minHeap.size() + 1) {
                minHeap.offer(maxHeap.poll());
            } else if (minHeap.size() > maxHeap.size()) {
                maxHeap.offer(minHeap.poll());
            }
            
            if (i >= k) {
                int toRemove = nums[i - k];
                if (toRemove <= maxHeap.peek()) {
                    maxHeap.remove(toRemove);
                } else {
                    minHeap.remove(toRemove);
                }
                
                if (maxHeap.size() > minHeap.size() + 1) {
                    minHeap.offer(maxHeap.poll());
                } else if (minHeap.size() > maxHeap.size()) {
                    maxHeap.offer(minHeap.poll());
                }
            }
            
            if (i >= k - 1) {
                if (k % 2 == 0) {
                    result[i - k + 1] = ((long)maxHeap.peek() + (long)minHeap.peek()) / 2.0;
                } else {
                    result[i - k + 1] = maxHeap.peek();
                }
            }
        }
        
        return result;
    }
}
```

**Time**: O(n k)
**Space**: O(k)

---

### 25. Smallest Range K Lists ⭐⭐⭐⭐
**LeetCode 632** | **Difficulty: Hard**

**Problem**: Find smallest range containing element from each list.

**Solution**:
```java
public class SmallestRange {
    public int[] smallestRange(List<List<Integer>> nums) {
        PriorityQueue<int[]> minHeap = new PriorityQueue<>(
            (a, b) -> a[0] - b[0]
        );
        
        int currentMax = Integer.MIN_VALUE;
        
        for (int i = 0; i < nums.size(); i++) {
            int value = nums.get(i).get(0);
            minHeap.offer(new int[]{value, i, 0});
            currentMax = Math.max(currentMax, value);
        }
        
        int rangeStart = 0;
        int rangeEnd = Integer.MAX_VALUE;
        
        while (minHeap.size() == nums.size()) {
            int[] current = minHeap.poll();
            int currentMin = current[0];
            int listIdx = current[1];
            int elementIdx = current[2];
            
            if (currentMax - currentMin < rangeEnd - rangeStart) {
                rangeStart = currentMin;
                rangeEnd = currentMax;
            }
            
            if (elementIdx + 1 < nums.get(listIdx).size()) {
                int nextValue = nums.get(listIdx).get(elementIdx + 1);
                minHeap.offer(new int[]{nextValue, listIdx, elementIdx + 1});
                currentMax = Math.max(currentMax, nextValue);
            }
        }
        
        return new int[]{rangeStart, rangeEnd};
    }
}
```

**Time**: O(n log k)
**Space**: O(k)

---

### 26. Constrained Subsequence Sum ⭐⭐⭐
**LeetCode 1425** | **Difficulty: Hard**

**Problem**: Maximum sum subsequence with constraint.

**Solution**:
```java
public class ConstrainedSubsequenceSum {
    public int constrainedSubsetSum(int[] nums, int k) {
        Deque<Integer> deque = new ArrayDeque<>();
        int[] dp = new int[nums.length];
        
        for (int i = 0; i < nums.length; i++) {
            while (!deque.isEmpty() && deque.peekFirst() < i - k) {
                deque.pollFirst();
            }
            
            dp[i] = nums[i];
            if (!deque.isEmpty()) {
                dp[i] = Math.max(dp[i], nums[i] + dp[deque.peekFirst()]);
            }
            
            while (!deque.isEmpty() && dp[deque.peekLast()] < dp[i]) {
                deque.pollLast();
            }
            
            if (dp[i] > 0) {
                deque.offerLast(i);
            }
        }
        
        int maxSum = Integer.MIN_VALUE;
        for (int val : dp) {
            maxSum = Math.max(maxSum, val);
        }
        
        return maxSum;
    }
}
```

**Time**: O(n)
**Space**: O(n)

---

### 27. IPO ⭐⭐⭐
**LeetCode 502** | **Difficulty: Hard**

**Problem**: Maximize capital by selecting projects.

**Solution**:
```java
public class IPO {
    public int findMaximizedCapital(int k, int w, int[] profits, int[] capital) {
        PriorityQueue<int[]> minCapital = new PriorityQueue<>(
            (a, b) -> a[0] - b[0]
        );
        
        PriorityQueue<Integer> maxProfit = new PriorityQueue<>(Collections.reverseOrder());
        
        for (int i = 0; i < profits.length; i++) {
            minCapital.offer(new int[]{capital[i], profits[i]});
        }
        
        for (int i = 0; i < k; i++) {
            while (!minCapital.isEmpty() && minCapital.peek()[0] <= w) {
                maxProfit.offer(minCapital.poll()[1]);
            }
            
            if (maxProfit.isEmpty()) {
                break;
            }
            
            w += maxProfit.poll();
        }
        
        return w;
    }
}
```

**Time**: O(n log n + k log n)
**Space**: O(n)

---

### 28. Trapping Rain Water II ⭐⭐⭐
**LeetCode 407** | **Difficulty: Hard**

**Problem**: Calculate trapped rain water in 2D.

**Solution**:
```java
public class TrappingRainWater2D {
    public int trapRainWater(int[][] heights) {
        if (heights == null || heights.length == 0) return 0;
        
        int m = heights.length;
        int n = heights[0].length;
        
        PriorityQueue<int[]> minHeap = new PriorityQueue<>(
            (a, b) -> a[2] - b[2]
        );
        
        boolean[][] visited = new boolean[m][n];
        
        // Add border cells
        for (int i = 0; i < m; i++) {
            minHeap.offer(new int[]{i, 0, heights[i][0]});
            minHeap.offer(new int[]{i, n - 1, heights[i][n - 1]});
            visited[i][0] = true;
            visited[i][n - 1] = true;
        }
        
        for (int j = 1; j < n - 1; j++) {
            minHeap.offer(new int[]{0, j, heights[0][j]});
            minHeap.offer(new int[]{m - 1, j, heights[m - 1][j]});
            visited[0][j] = true;
            visited[m - 1][j] = true;
        }
        
        int[][] dirs = {{0, 1}, {1, 0}, {0, -1}, {-1, 0}};
        int water = 0;
        
        while (!minHeap.isEmpty()) {
            int[] cell = minHeap.poll();
            int row = cell[0];
            int col = cell[1];
            int height = cell[2];
            
            for (int[] dir : dirs) {
                int newRow = row + dir[0];
                int newCol = col + dir[1];
                
                if (newRow >= 0 && newRow < m && newCol >= 0 && newCol < n
                    && !visited[newRow][newCol]) {
                    
                    water += Math.max(0, height - heights[newRow][newCol]);
                    minHeap.offer(new int[]{
                        newRow, newCol, 
                        Math.max(height, heights[newRow][newCol])
                    });
                    visited[newRow][newCol] = true;
                }
            }
        }
        
        return water;
    }
}
```

**Time**: O(m×n log(m×n))
**Space**: O(m×n)

---

## Topic-Wise Problems

### Queue Basics
1. ✅ Implement Queue using Stacks (LC 232) - Easy
2. ✅ Implement Stack using Queues (LC 225) - Easy
3. ✅ Number of Recent Calls (LC 933) - Easy
4. Binary Tree Level Order Traversal (LC 102) - Medium
5. Rotting Oranges (LC 994) - Medium
6. Open the Lock (LC 752) - Medium
7. Perfect Squares (LC 279) - Medium

### Circular Queue
1. ✅ Design Circular Queue (LC 622) - Medium
2. Design Hit Counter (LC 362) - Medium
3. Moving Average from Data Stream (LC 346) - Easy
4. Gas Station (LC 134) - Medium

### Deque
1. ✅ Design Circular Deque (LC 641) - Medium
2. ✅ Sliding Window Maximum (LC 239) - Hard
3. Longest Continuous Subarray With Limit (LC 1438) - Medium
4. Jump Game VI (LC 1696) - Medium
5. Shortest Subarray with Sum ≥ K (LC 862) - Hard
6. Constrained Subsequence Sum (LC 1425) - Hard

### Priority Queue
1. ✅ Kth Largest Element (LC 215) - Medium
2. ✅ Top K Frequent Elements (LC 347) - Medium
3. ✅ K Closest Points (LC 973) - Medium
4. ✅ Find Median from Data Stream (LC 295) - Hard
5. ✅ Merge K Sorted Lists (LC 23) - Hard
6. Reorganize String (LC 767) - Medium
7. Task Scheduler (LC 621) - Medium
8. Sliding Window Median (LC 480) - Hard
9. Smallest Range K Lists (LC 632) - Hard
10. IPO (LC 502) - Hard

---

## Company-Wise Problems

### Amazon ⭐⭐⭐⭐⭐
1. Binary Tree Level Order Traversal (LC 102)
2. Kth Largest Element (LC 215)
3. Top K Frequent Elements (LC 347)
4. Merge K Sorted Lists (LC 23)
5. Sliding Window Maximum (LC 239)
6. Rotting Oranges (LC 994)
7. K Closest Points (LC 973)

### Google ⭐⭐⭐⭐⭐
1. Sliding Window Maximum (LC 239)
2. Find Median from Data Stream (LC 295)
3. Merge K Sorted Lists (LC 23)
4. Task Scheduler (LC 621)
5. Smallest Range K Lists (LC 632)
6. IPO (LC 502)

### Microsoft ⭐⭐⭐⭐
1. Implement Queue using Stacks (LC 232)
2. Design Circular Queue (LC 622)
3. Kth Largest Element (LC 215)
4. Reorganize String (LC 767)
5. Open the Lock (LC 752)

### Facebook/Meta ⭐⭐⭐⭐
1. Binary Tree Level Order Traversal (LC 102)
2. Top K Frequent Elements (LC 347)
3. K Closest Points (LC 973)
4. Task Scheduler (LC 621)
5. Sliding Window Median (LC 480)

### Apple ⭐⭐⭐
1. Design Circular Deque (LC 641)
2. Kth Largest Element in Stream (LC 703)
3. Last Stone Weight (LC 1046)
4. Perfect Squares (LC 279)

---

## Problem Solutions

### Strategy Guide

**1. Queue Basics Problems**:
- Use `Queue<>` interface with `LinkedList` or `ArrayDeque`
- BFS pattern: Level-by-level processing
- Multi-source BFS: Add all sources initially

**2. Circular Queue Problems**:
- Use modulo arithmetic: `(index + 1) % capacity`
- Track `front`, `rear`, and `size`
- Handle wrap-around carefully

**3. Deque Problems**:
- **Monotonic Deque**: Maintain increasing/decreasing order
- Remove from both ends for optimization
- Store indices instead of values for window problems

**4. Priority Queue Problems**:
- **Kth Largest**: Use min heap of size k
- **Kth Smallest**: Use max heap of size k
- **Median**: Use two heaps (max + min)
- **Merge K**: Use min heap with k elements

### Common Patterns

**Pattern 1: BFS Level Order**
```java
Queue<Node> queue = new LinkedList<>();
queue.offer(root);

while (!queue.isEmpty()) {
    int levelSize = queue.size();
    
    for (int i = 0; i < levelSize; i++) {
        Node node = queue.poll();
        // Process node
        // Add children
    }
}
```

**Pattern 2: Sliding Window with Deque**
```java
Deque<Integer> deque = new ArrayDeque<>();

for (int i = 0; i < n; i++) {
    // Remove outside window
    while (!deque.isEmpty() && deque.peekFirst() < i - k + 1) {
        deque.pollFirst();
    }
    
    // Maintain monotonic property
    while (!deque.isEmpty() && condition) {
        deque.pollLast();
    }
    
    deque.offerLast(i);
}
```

**Pattern 3: Top K with Min Heap**
```java
PriorityQueue<Integer> minHeap = new PriorityQueue<>();

for (int num : nums) {
    minHeap.offer(num);
    
    if (minHeap.size() > k) {
        minHeap.poll();
    }
}

return minHeap.peek();  // Kth largest
```

**Pattern 4: Two Heaps for Median**
```java
PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
PriorityQueue<Integer> minHeap = new PriorityQueue<>();

// Add to maxHeap first
maxHeap.offer(num);

// Balance
minHeap.offer(maxHeap.poll());

if (maxHeap.size() < minHeap.size()) {
    maxHeap.offer(minHeap.poll());
}
```

---

## Summary

### Must Solve Problems (Top 10)

**Priority Queue**:
1. ✅ Kth Largest Element (LC 215) - ⭐⭐⭐⭐⭐
2. ✅ Top K Frequent (LC 347) - ⭐⭐⭐⭐⭐
3. ✅ Find Median from Stream (LC 295) - ⭐⭐⭐⭐⭐
4. ✅ Merge K Sorted Lists (LC 23) - ⭐⭐⭐⭐⭐

**Deque**:
5. ✅ Sliding Window Maximum (LC 239) - ⭐⭐⭐⭐⭐
6. ✅ Shortest Subarray Sum ≥ K (LC 862) - ⭐⭐⭐⭐⭐

**Queue Basics**:
7. ✅ Binary Tree Level Order (LC 102) - ⭐⭐⭐⭐⭐
8. ✅ Implement Queue using Stacks (LC 232) - ⭐⭐⭐⭐⭐

**Circular Queue**:
9. ✅ Design Circular Queue (LC 622) - ⭐⭐⭐⭐⭐
10. ✅ Design Circular Deque (LC 641) - ⭐⭐⭐⭐⭐

### Complexity Reference

| Data Structure | Insert | Delete | Peek | Space |
|----------------|--------|--------|------|-------|
| Queue | O(1) | O(1) | O(1) | O(n) |
| Circular Queue | O(1) | O(1) | O(1) | O(k) |
| Deque | O(1) | O(1) | O(1) | O(n) |
| Priority Queue | O(log n) | O(log n) | O(1) | O(n) |

### Study Plan

**Week 1: Queue Basics**
- Day 1-2: Implement Queue using Stacks, Stack using Queues
- Day 3-4: Design Circular Queue, Number of Recent Calls
- Day 5-6: Binary Tree Level Order, Rotting Oranges
- Day 7: Review and practice

**Week 2: Deque & Priority Queue**
- Day 1-2: Sliding Window Maximum, Design Circular Deque
- Day 3-4: Kth Largest, Top K Frequent
- Day 5-6: Find Median from Stream, Merge K Sorted Lists
- Day 7: Review and practice

**Week 3: Advanced Problems**
- Day 1-2: Shortest Subarray Sum ≥ K, Task Scheduler
- Day 3-4: Sliding Window Median, Smallest Range
- Day 5-6: Constrained Subsequence Sum, IPO
- Day 7: Mock interview practice

---

**Previous**: [9.4 Priority Queue](9.4-Priority-Queue.md)

**Next**: Move to [10. Hashing](../10-Hashing/01-HashMap-Basics.md)

**Remember**: Master the patterns - BFS, Monotonic Deque, Top K, and Two Heaps!
