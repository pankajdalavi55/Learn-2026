# Time & Space Complexity Analysis

## Why Should You Care?

Imagine you wrote two solutions for the same problem. Both give correct output. But one takes **3 seconds** and another takes **3 hours** for large input. Which one would you choose? That's where complexity analysis comes in!

In competitive programming, your solution must be:
- ✅ **Correct** (obviously!)
- ✅ **Fast enough** (within time limits)
- ✅ **Memory efficient** (within memory limits)

---

## Time Complexity - How Long Does It Take?

**Time Complexity** tells you how the runtime grows as input size increases.

### Real Example - Finding Maximum in Array

```java
public class ComplexityDemo {
    
    // Approach 1: Single pass through array
    public static int findMax(int[] arr) {
        int max = arr[0];
        for (int i = 1; i < arr.length; i++) {
            if (arr[i] > max) {
                max = arr[i];
            }
        }
        return max;
    }
    
    public static void main(String[] args) {
        int[] arr = {5, 2, 9, 1, 7};
        System.out.println("Max: " + findMax(arr));
    }
}
```

**Analysis:**
- Loop runs **n times** (where n = array length)
- Each iteration does constant work (1 comparison, maybe 1 assignment)
- Total: **O(n)** - Linear Time

---

## Counting Operations - The Practical Way

Let's count actual operations step by step:

### Example 1: Simple Loop
```java
public void printNumbers(int n) {
    for (int i = 0; i < n; i++) {      // runs n times
        System.out.println(i);          // 1 operation each time
    }
}
```
**Total operations:** n × 1 = **n operations** → **O(n)**

---

### Example 2: Nested Loops
```java
public void printPairs(int n) {
    for (int i = 0; i < n; i++) {           // runs n times
        for (int j = 0; j < n; j++) {       // runs n times for each i
            System.out.println(i + "," + j);
        }
    }
}
```
**Total operations:** n × n = **n² operations** → **O(n²)**

---

### Example 3: Two Separate Loops
```java
public void twoLoops(int n) {
    // First loop
    for (int i = 0; i < n; i++) {
        System.out.println(i);
    }
    
    // Second loop
    for (int j = 0; j < n; j++) {
        System.out.println(j);
    }
}
```
**Total operations:** n + n = **2n operations** → **O(n)**
*(We drop constants, so 2n becomes n)*

---

### Example 4: Loop That Doubles
```java
public void doubling(int n) {
    for (int i = 1; i < n; i = i * 2) {     // i: 1, 2, 4, 8, 16...
        System.out.println(i);
    }
}
```
**How many times?**
- i starts at 1, then 2, 4, 8, 16... until < n
- If n = 16: runs 4 times (1, 2, 4, 8)
- 2^k = n → k = log₂(n)

**Total operations:** **log(n) operations** → **O(log n)**

---

## Common Time Complexities (Best to Worst)

| Complexity | Name | Example | For n=1000 | For n=1,000,000 |
|------------|------|---------|------------|------------------|
| **O(1)** | Constant | Accessing array[5] | 1 | 1 |
| **O(log n)** | Logarithmic | Binary search | ~10 | ~20 |
| **O(n)** | Linear | Single loop | 1,000 | 1,000,000 |
| **O(n log n)** | Linearithmic | Merge sort | ~10,000 | ~20,000,000 |
| **O(n²)** | Quadratic | Nested loops | 1,000,000 | 1,000,000,000,000 |
| **O(2ⁿ)** | Exponential | Recursive fibonacci | HUGE | Impossible! |

---

## Practical Time Limits in Competitive Programming

Most online judges allow **1-2 seconds** of execution time. Here's what you can do:

| Time Complexity | Max Input Size (approx) | When to Use |
|-----------------|-------------------------|-------------|
| O(1) | Any | Direct formulas, array access |
| O(log n) | n ≤ 10^18 | Binary search, tree operations |
| O(n) | n ≤ 10^8 | Single loop over data |
| O(n log n) | n ≤ 10^6 | Sorting, efficient algorithms |
| O(n²) | n ≤ 10^4 | Small inputs, brute force |
| O(n³) | n ≤ 500 | Very small inputs only |
| O(2ⁿ) | n ≤ 20 | Backtracking small cases |

**Rule of Thumb:** Your code should do about **10^8 operations per second**.

---

## Real Problem - Sum of Array

Let's solve a real problem with different approaches:

**Problem:** Find sum of all elements in array

### Approach 1: Using Loop - O(n)
```java
public static int sumArray(int[] arr) {
    int sum = 0;
    for (int i = 0; i < arr.length; i++) {      // runs n times
        sum += arr[i];                           // 1 operation
    }
    return sum;
}
// Time: O(n), Space: O(1)
```

### Approach 2: Using Recursion - O(n)
```java
public static int sumRecursive(int[] arr, int index) {
    if (index == arr.length) {
        return 0;
    }
    return arr[index] + sumRecursive(arr, index + 1);
}
// Time: O(n), Space: O(n) - due to recursion stack
```

**Which is better?** Approach 1! Same time complexity but better space.

---

## Space Complexity - How Much Memory?

Space complexity measures **extra memory** your algorithm uses.

### Example 1: Constant Space - O(1)
```java
public int sum(int n) {
    int total = 0;              // 1 variable
    for (int i = 1; i <= n; i++) {
        total += i;
    }
    return total;
}
// Only uses fixed variables (total, i) regardless of n
```

---

### Example 2: Linear Space - O(n)
```java
public int[] createArray(int n) {
    int[] arr = new int[n];     // Array of size n
    for (int i = 0; i < n; i++) {
        arr[i] = i;
    }
    return arr;
}
// Creates array of size n
```

---

### Example 3: Recursion Space - O(n)
```java
public int factorial(int n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1);
}
// Each recursive call adds to call stack
// Maximum depth = n, so space = O(n)
```

---

## Practical Exercise - Analyze These!

Try to figure out time and space complexity:

### Exercise 1
```java
public void mystery1(int n) {
    for (int i = 0; i < n; i++) {
        for (int j = i; j < n; j++) {
            System.out.println(i + " " + j);
        }
    }
}
```
<details>
<summary>Click for answer</summary>

**Time:** O(n²)
- Outer loop: n times
- Inner loop: (n-i) times
- Total: n + (n-1) + (n-2) + ... + 1 = n(n+1)/2 ≈ n²

**Space:** O(1) - only uses i, j
</details>

---

### Exercise 2
```java
public boolean isPrime(int n) {
    for (int i = 2; i * i <= n; i++) {
        if (n % i == 0) {
            return false;
        }
    }
    return true;
}
```
<details>
<summary>Click for answer</summary>

**Time:** O(√n)
- Loop runs while i² ≤ n
- So i goes from 2 to √n

**Space:** O(1) - only variable i
</details>

---

### Exercise 3
```java
public int fibonacci(int n) {
    if (n <= 1) return n;
    return fibonacci(n-1) + fibonacci(n-2);
}
```
<details>
<summary>Click for answer</summary>

**Time:** O(2ⁿ) - exponential!
- Each call makes 2 more calls
- Creates a tree of depth n

**Space:** O(n)
- Maximum recursion depth is n
- Call stack can go n levels deep
</details>

---

## Quick Decision Guide

When you see a problem, ask yourself:

1. **What's the input size?**
   - n ≤ 20 → O(2ⁿ) or O(n!) might work
   - n ≤ 500 → O(n³) is okay
   - n ≤ 10,000 → O(n²) is fine
   - n ≤ 1,000,000 → Need O(n log n) or better
   - n ≤ 100,000,000 → Must be O(n) or O(log n)

2. **What complexity do I need?**
   - Look at constraints in problem
   - Calculate: Can my solution handle max input in 1-2 seconds?

---

## Key Takeaways

✅ **Always analyze before coding** - Think about complexity first
✅ **Count the loops** - Nested loops multiply, sequential loops add
✅ **Check constraints** - Problem tells you what complexity you need
✅ **Space matters too** - Don't create huge arrays if not needed
✅ **Test with large input** - Does it finish in time?

---

## What's Next?

Now that you understand complexity analysis, let's dive deeper into **Big O Notation** to formalize these concepts!

[Next: Big O Notation →](./BigO-Notation.md)
