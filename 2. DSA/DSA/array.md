## Longest Substring without repeating characters

# ✅ Problem (one-liner you can say in interview)

> Maintain a sliding window `[l, r]` such that it always contains **unique characters**, and maximize its length.

---

# 🔹 Variant 1: `int[]` (ASCII-optimized)

```
public static int lengthOfLongestSubstring_Array(String s) {
    int n = s.length();
    int l = 0, maxLen = 0;

    // lastSeen[c] = (last index of c) + 1
    int[] lastSeen = new int[128]; // ASCII

    for (int r = 0; r < n; r++) {
        char c = s.charAt(r);

        // Jump left pointer to maintain uniqueness
        l = Math.max(l, lastSeen[c]);

        // Update answer using current window size
        maxLen = Math.max(maxLen, r - l + 1);

        // Record last seen position (+1 to avoid default 0 ambiguity)
        lastSeen[c] = r + 1;
    }

    return maxLen;
}
```

---

## 🧠 How to explain (crisp + strong)

**1. Data structure choice**

- 

I use an `int[128]` for **O(1) direct indexing**  

- 

Works because input is assumed **ASCII**  

---

**2. Key invariant**

> Window `[l, r]` always contains **no duplicate characters**

---

**3. Core idea (important line)**

```
l = Math.max(l, lastSeen[c]);
```

Explain like this:

- 

If character `c` was seen before, I **jump** `l` **directly**  

- 

I use `Math.max` to ensure `l` **never moves backward**  

---

**4. Why store** `r + 1`

- Default array value is `0`
- So:
  - `0` → never seen
  - `>0` → seen before
- Avoids off-by-one bugs

---

**5. Complexity**

- Time: **O(n)**
- Space: **O(1)**

---

# 🔹 Variant 2: `Map<Character, Integer>` (Unicode-safe)

```
public static int lengthOfLongestSubstring_Map(String s) {
    int n = s.length();
    int l = 0, maxLen = 0;

    Map<Character, Integer> lastSeen = new HashMap<>();

    for (int r = 0; r < n; r++) {
        char c = s.charAt(r);

        if (lastSeen.containsKey(c)) {
            l = Math.max(l, lastSeen.get(c));
        }

        maxLen = Math.max(maxLen, r - l + 1);

        lastSeen.put(c, r + 1);
    }

    return maxLen;
}
```

---

## 🧠 How to explain (interview version)

**1. Why Map instead of array**

- 

Supports **Unicode / large character sets**  

- 

More flexible than fixed-size array  

---

**2. Logic is identical**

- 

Map stores:  

→ `character → (last index + 1)`  

- 

Same **jumping left pointer technique**  

---

**3. Key line explanation**

```
l = Math.max(l, lastSeen.get(c));
```

- 

Ensures window remains valid (no duplicates)  

---

**4. Complexity**

- Time: **O(n)**  
- Space: **O(min(n, charset))**

---

# ⚖️ When to use which (say this confidently)


| Scenario                  | Choose           |
| ------------------------- | ---------------- |
| ASCII input               | `int[]` (faster) |
| Unicode / unknown charset | `Map`            |
| Performance-critical      | `int[]`          |


---

# 🔥 Final 30-sec interview summary

> “I use a sliding window with two pointers.  
> I track the last seen position of each character and whenever I see a duplicate, I jump the left pointer instead of moving it step-by-step. This guarantees O(n) time.  
> For ASCII inputs, I use an array for constant-time access; otherwise, I use a HashMap.”

---

# ⚠️ Common pitfalls (mention if asked)

- Forgetting `Math.max(l, …)` → breaks window  
- Storing `r` instead of `r+1` → off-by-one bugs  
- Using `Set` and removing one-by-one → less optimal pattern

