# Big O Notation - The Complete Practical Guide

## What is Big O? (Simple Explanation)

Big O is a way to describe **how your code slows down** as input gets bigger.

Think of it like this:
- You're making sandwiches
- O(1): Making 1 sandwich takes same time whether you're serving 1 or 1000 people
- O(n): Making sandwiches for n people takes n times as long
- O(n²): You make a sandwich for each person, and each sandwich needs approval from everyone!

---

## The Rules of Big O

### Rule 1: Drop Constants

```java
// Algorithm 1: Two loops
public void algo1(int n) {
    for (int i = 0; i < n; i++) {
        System.out.println(i);
    }
    for (int j = 0; j < n; j++) {
        System.out.println(j);
    }
}
```

**Steps:** n + n = 2n
**Big O:** O(2n) → **O(n)** ✅

We **drop the constant 2** because:
- When n = 1,000,000: difference between n and 2n is tiny
- We care about growth rate, not exact count

---

### Rule 2: Drop Non-Dominant Terms

```java
public void algo2(int n) {
    // First part: O(n²)
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n; j++) {
            System.out.println(i * j);
        }
    }
    
    // Second part: O(n)
    for (int k = 0; k < n; k++) {
        System.out.println(k);
    }
}
```

**Steps:** n² + n
**Big O:** O(n² + n) → **O(n²)** ✅

We **drop the +n** because:
- When n = 1000: n² = 1,000,000 but n = 1,000
- n² dominates completely!

**Dominance order:** O(1) < O(log n) < O(n) < O(n log n) < O(n²) < O(n³) < O(2ⁿ) < O(n!)

---

### Rule 3: Different Inputs Get Different Variables

```java
public void printPairs(int[] arr1, int[] arr2) {
    // First array
    for (int i = 0; i < arr1.length; i++) {
        System.out.println(arr1[i]);
    }
    
    // Second array
    for (int j = 0; j < arr2.length; j++) {
        System.out.println(arr2[j]);
    }
}
```

**Big O:** O(a + b) where a = arr1.length, b = arr2.length
**NOT** O(n) or O(2n)!

---

### Rule 4: Watch for Early Exits

```java
public boolean contains(int[] arr, int target) {
    for (int i = 0; i < arr.length; i++) {
        if (arr[i] == target) {
            return true;    // Early exit!
        }
    }
    return false;
}
```

**Best case:** O(1) - found at first position
**Worst case:** O(n) - found at last position or not found
**Average case:** O(n/2) → O(n)

**We use worst case** for Big O: **O(n)**

---

## Common Big O Complexities (With Real Examples)

### 1. O(1) - Constant Time

**The Flash** - Always instant, no matter what!

```java
// Example 1: Array access
public int getElement(int[] arr, int index) {
    return arr[index];              // Always 1 operation
}

// Example 2: HashMap get
public String getName(HashMap<Integer, String> map, int id) {
    return map.get(id);             // Average O(1)
}

// Example 3: Math formula
public int sumFirstN(int n) {
    return n * (n + 1) / 2;         // Direct formula
}

// Even with multiple operations - still O(1)
public int doStuff(int n) {
    int a = 5;
    int b = 10;
    int c = a + b;
    int d = c * 2;
    return d;                       // 4 operations, but fixed!
}
```

**Key:** Number of operations doesn't depend on input size!

---

### 2. O(log n) - Logarithmic Time

**The Smart Detective** - Cuts problem in half each time

```java
// Example 1: Binary Search
public int binarySearch(int[] arr, int target) {
    int left = 0, right = arr.length - 1;
    
    while (left <= right) {
        int mid = left + (right - left) / 2;
        
        if (arr[mid] == target) {
            return mid;
        } else if (arr[mid] < target) {
            left = mid + 1;         // Cut left half
        } else {
            right = mid - 1;        // Cut right half
        }
    }
    return -1;
}

// Example 2: Finding power (divide & conquer)
public int power(int base, int exp) {
    if (exp == 0) return 1;
    if (exp == 1) return base;
    
    int half = power(base, exp / 2);    // Divide by 2 each time!
    
    if (exp % 2 == 0) {
        return half * half;
    } else {
        return half * half * base;
    }
}
```

**Why O(log n)?**
- n = 1024 → 10 steps (2^10 = 1024)
- Each step cuts problem in half

---

### 3. O(n) - Linear Time

**The Thorough Inspector** - Check everything once

```java
// Example 1: Find maximum
public int findMax(int[] arr) {
    int max = arr[0];
    for (int i = 1; i < arr.length; i++) {
        if (arr[i] > max) {
            max = arr[i];
        }
    }
    return max;
}

// Example 2: Count occurrences
public int countOccurrences(int[] arr, int target) {
    int count = 0;
    for (int num : arr) {
        if (num == target) {
            count++;
        }
    }
    return count;
}

// Example 3: Reverse array
public void reverseArray(int[] arr) {
    int left = 0, right = arr.length - 1;
    while (left < right) {
        int temp = arr[left];
        arr[left] = arr[right];
        arr[right] = temp;
        left++;
        right--;
    }
}
```

---

### 4. O(n log n) - Linearithmic Time

**The Organized Manager** - Efficient sorting/divide-conquer

```java
// Example: Merge Sort
public void mergeSort(int[] arr, int left, int right) {
    if (left < right) {
        int mid = (left + right) / 2;
        
        mergeSort(arr, left, mid);      // Divide
        mergeSort(arr, mid + 1, right); // Divide
        merge(arr, left, mid, right);   // Conquer (O(n))
    }
}

private void merge(int[] arr, int left, int mid, int right) {
    // Merging takes O(n) time
    // Called log(n) times (tree depth)
    // Total: O(n log n)
}
```

**Common in:**
- Merge Sort
- Quick Sort (average case)
- Heap Sort

---

### 5. O(n²) - Quadratic Time

**The Double Checker** - Compare everything with everything

```java
// Example 1: Bubble Sort
public void bubbleSort(int[] arr) {
    for (int i = 0; i < arr.length; i++) {          // n times
        for (int j = 0; j < arr.length - 1; j++) {  // n times
            if (arr[j] > arr[j + 1]) {
                // Swap
                int temp = arr[j];
                arr[j] = arr[j + 1];
                arr[j + 1] = temp;
            }
        }
    }
}

// Example 2: Find all pairs
public void printAllPairs(int[] arr) {
    for (int i = 0; i < arr.length; i++) {
        for (int j = 0; j < arr.length; j++) {
            System.out.println(arr[i] + ", " + arr[j]);
        }
    }
}

// Example 3: Find duplicates (naive)
public boolean hasDuplicates(int[] arr) {
    for (int i = 0; i < arr.length; i++) {
        for (int j = i + 1; j < arr.length; j++) {
            if (arr[i] == arr[j]) {
                return true;
            }
        }
    }
    return false;
}
```

---

### 6. O(2ⁿ) - Exponential Time

**The Exploding Tree** - Avoid if possible!

```java
// Example: Fibonacci (naive recursion)
public int fibonacci(int n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);  // 2 recursive calls!
}

// Example: Generate all subsets
public void generateSubsets(int[] arr, int index, List<Integer> current) {
    if (index == arr.length) {
        // Process subset
        return;
    }
    
    // Don't include current element
    generateSubsets(arr, index + 1, current);
    
    // Include current element
    current.add(arr[index]);
    generateSubsets(arr, index + 1, current);
    current.remove(current.size() - 1);
}
```

**Why so slow?**
- Each element: 2 choices (include or exclude)
- n elements: 2 × 2 × 2... (n times) = 2ⁿ possibilities

---

## Practice: Analyze These Functions

### Challenge 1
```java
public void function1(int n) {
    int sum = 0;
    for (int i = 0; i < n; i++) {
        sum += i;
    }
    System.out.println(sum);
}
```
<details>
<summary>Answer</summary>

**O(n)** - Single loop running n times
</details>

---

### Challenge 2
```java
public void function2(int n) {
    for (int i = 0; i < n; i = i + 2) {
        System.out.println(i);
    }
}
```
<details>
<summary>Answer</summary>

**O(n)** - Loop runs n/2 times, which is still O(n)
</details>

---

### Challenge 3
```java
public void function3(int n) {
    for (int i = 1; i < n; i = i * 2) {
        for (int j = 0; j < n; j++) {
            System.out.println(i + " " + j);
        }
    }
}
```
<details>
<summary>Answer</summary>

**O(n log n)**
- Outer loop: log(n) times (doubling)
- Inner loop: n times
- Total: log(n) × n = n log n
</details>

---

### Challenge 4
```java
public void function4(int n) {
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < i; j++) {
            System.out.println(i + " " + j);
        }
    }
}
```
<details>
<summary>Answer</summary>

**O(n²)**
- When i=0: 0 iterations
- When i=1: 1 iteration
- When i=2: 2 iterations
- ...
- Total: 0+1+2+...+(n-1) = n(n-1)/2 ≈ n²
</details>

---

### Challenge 5
```java
public int function5(int n) {
    if (n <= 1) return 1;
    return function5(n - 1) + function5(n - 1);
}
```
<details>
<summary>Answer</summary>

**O(2ⁿ)**
- Each call makes 2 recursive calls
- Tree depth is n
- Total nodes: 2ⁿ
</details>

---

## Space Complexity - Big O for Memory

### O(1) Space - Constant
```java
public int sum(int a, int b) {
    int result = a + b;     // Only few variables
    return result;
}
```

### O(n) Space - Linear
```java
public int[] doubleArray(int[] arr) {
    int[] result = new int[arr.length];  // New array of size n
    for (int i = 0; i < arr.length; i++) {
        result[i] = arr[i] * 2;
    }
    return result;
}
```

### O(n) Space - Recursion Stack
```java
public int factorial(int n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1);  // n recursive calls on stack
}
```

---

## Real Interview Question - Optimize This!

**Problem:** Check if array has duplicate

### Solution 1: Brute Force - O(n²) time, O(1) space
```java
public boolean hasDuplicate1(int[] arr) {
    for (int i = 0; i < arr.length; i++) {
        for (int j = i + 1; j < arr.length; j++) {
            if (arr[i] == arr[j]) return true;
        }
    }
    return false;
}
```

### Solution 2: Using HashSet - O(n) time, O(n) space
```java
public boolean hasDuplicate2(int[] arr) {
    HashSet<Integer> seen = new HashSet<>();
    for (int num : arr) {
        if (seen.contains(num)) return true;
        seen.add(num);
    }
    return false;
}
```

### Solution 3: Sorting First - O(n log n) time, O(1) space
```java
public boolean hasDuplicate3(int[] arr) {
    Arrays.sort(arr);  // O(n log n)
    for (int i = 0; i < arr.length - 1; i++) {
        if (arr[i] == arr[i + 1]) return true;
    }
    return false;
}
```

**Which to choose?** Depends on constraints!
- Small array + limited memory → Solution 1 or 3
- Large array + plenty memory → Solution 2
- Can't modify array → Solution 2

---

## Quick Reference Table

| Big O | Name | Example Operations |
|-------|------|-------------------|
| O(1) | Constant | Array access, hash table get, math formula |
| O(log n) | Logarithmic | Binary search, balanced tree operations |
| O(n) | Linear | Loop through array, linear search |
| O(n log n) | Linearithmic | Efficient sorting (merge/quick/heap sort) |
| O(n²) | Quadratic | Nested loops, bubble sort |
| O(n³) | Cubic | Triple nested loops |
| O(2ⁿ) | Exponential | Recursive fibonacci, generate all subsets |
| O(n!) | Factorial | Generate all permutations |

---

## Key Takeaways

✅ **Big O describes growth rate** - not exact operations
✅ **Drop constants and non-dominant terms** - O(2n + 5) → O(n)
✅ **Use worst case** - unless specified otherwise
✅ **Different inputs = different variables** - O(a + b), not O(n)
✅ **Space matters too** - sometimes trade time for space or vice versa

---

[← Back: Complexity Analysis](./Complexity-Analysis.md) | [Next: Problem Solving Patterns →](./Problem-Solving-Patterns.md)
