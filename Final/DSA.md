Here is a curated list of **15 most frequently asked Two-Pointer problems** in coding interviews. They are categorized by the pointer strategies they employ (Opposite Direction, Fast & Slow, Sliding Window, and Expansion).

For each problem, I have provided a detailed breakdown of the statement, input/output specifications, constraints, and examples—**without including the solution or algorithmic approach.**

---

### 1. Two Sum II - Input Array Is Sorted
- **Category:** Opposite Direction (Left & Right)
- **Problem Statement:** Given a **1-indexed** array of integers `numbers` that is already sorted in non-decreasing order, find two numbers such that they add up to a specific `target` number. Let these two numbers be `numbers[index1]` and `numbers[index2]` where `1 <= index1 < index2 <= numbers.length`. Return the indices of the two numbers, `index1` and `index2`, as an integer array `[index1, index2]` of length 2.
- **Input:** An array of integers `numbers` (sorted ascending) and an integer `target`.
- **Output:** An array of two integers representing the 1-based indices.
- **Constraints:**
    - `2 <= numbers.length <= 3 * 10^4`
    - `-1000 <= numbers[i] <= 1000`
    - `numbers` is sorted in non-decreasing order.
    - `-1000 <= target <= 1000`
    - The tests are generated such that exactly **one** valid solution exists.
- **Example:** `Input: numbers = [2,7,11,15], target = 9` → `Output: [1,2]`

---

### 2. 3Sum
- **Category:** Opposite Direction (with fixed pivot)
- **Problem Statement:** Given an integer array `nums`, return all the triplets `[nums[i], nums[j], nums[k]]` such that `i != j`, `i != k`, and `j != k`, and `nums[i] + nums[j] + nums[k] == 0`. The solution set must not contain duplicate triplets.
- **Input:** A single integer array `nums`.
- **Output:** A 2D list containing all unique triplets that sum to zero.
- **Constraints:**
    - `0 <= nums.length <= 3000`
    - `-10^5 <= nums[i] <= 10^5`
- **Example:** `Input: nums = [-1,0,1,2,-1,-4]` → `Output: [[-1,-1,2],[-1,0,1]]`

---

### 3. Container With Most Water
- **Category:** Opposite Direction (Greedy shift)
- **Problem Statement:** You are given an integer array `height` of length `n`. There are `n` vertical lines drawn such that the two endpoints of the `ith` line are `(i, 0)` and `(i, height[i])`. Find two lines that together with the x-axis form a container, such that the container contains the most water. Return the **maximum amount of water** a container can store. Notice that you may not slant the container.
- **Input:** An array of integers `height`.
- **Output:** An integer representing the maximum area.
- **Constraints:**
    - `n == height.length`
    - `2 <= n <= 10^5`
    - `0 <= height[i] <= 10^4`
- **Example:** `Input: height = [1,8,6,2,5,4,8,3,7]` → `Output: 49`

---

### 4. Valid Palindrome
- **Category:** Opposite Direction (with character skipping)
- **Problem Statement:** A phrase is a **palindrome** if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers. Given a string `s`, return `true` if it is a palindrome, or `false` otherwise.
- **Input:** A string `s`.
- **Output:** A boolean value (`true` or `false`).
- **Constraints:**
    - `1 <= s.length <= 2 * 10^5`
    - `s` consists only of printable ASCII characters.
- **Example:** `Input: s = "A man, a plan, a canal: Panama"` → `Output: true` (Explanation: After normalization, it becomes "amanaplanacanalpanama").

---

### 5. Trapping Rain Water
- **Category:** Opposite Direction (maintaining boundary maxima)
- **Problem Statement:** Given `n` non-negative integers representing an elevation map where the width of each bar is `1`, compute how much water it can trap after raining.
- **Input:** An array of integers `height`.
- **Output:** An integer representing the total units of trapped rainwater.
- **Constraints:**
    - `n == height.length`
    - `1 <= n <= 2 * 10^4`
    - `0 <= height[i] <= 10^5`
- **Example:** `Input: height = [0,1,0,2,1,0,1,3,2,1,2,1]` → `Output: 6`

---

### 6. Remove Duplicates from Sorted Array
- **Category:** Fast & Slow Pointer (In-place overwrite)
- **Problem Statement:** Given an integer array `nums` sorted in non-decreasing order, remove the duplicates **in-place** such that each unique element appears only once. The relative order of the elements should be kept the same. Return the number of unique elements (`k`). Do not allocate extra space for another array; you must do this by modifying the input array in-place with O(1) extra memory.
- **Input:** A sorted integer array `nums`.
- **Output:** An integer `k` (the number of unique elements). The first `k` elements of `nums` should hold the final result.
- **Constraints:**
    - `1 <= nums.length <= 3 * 10^4`
    - `-100 <= nums[i] <= 100`
    - `nums` is sorted in non-decreasing order.
- **Example:** `Input: nums = [0,0,1,1,1,2,2,3,3,4]` → `Output: 5` (The array becomes `[0,1,2,3,4,_,_,_,_,_]`).

---

### 7. Move Zeroes
- **Category:** Fast & Slow Pointer (Partitioning)
- **Problem Statement:** Given an integer array `nums`, move all `0`'s to the end of it while maintaining the relative order of the non-zero elements. You must do this **in-place** without making a copy of the array. Minimize the total number of operations.
- **Input:** An integer array `nums`.
- **Output:** The modified array `nums` (returned as void, but transformed in-place).
- **Constraints:**
    - `1 <= nums.length <= 10^4`
    - `-2^31 <= nums[i] <= 2^31 - 1`
- **Example:** `Input: nums = [0,1,0,3,12]` → `Output: [1,3,12,0,0]`

---

### 8. Remove Element
- **Category:** Fast & Slow Pointer (In-place filtering)
- **Problem Statement:** Given an integer array `nums` and an integer `val`, remove all occurrences of `val` in `nums` **in-place**. The order of the elements may be changed. Then return the number of elements in `nums` which are not equal to `val`. Consider the number of elements not equal to `val` be `k`. To get accepted, you need to do the following: Change the array `nums` such that the first `k` elements of `nums` contain the elements not equal to `val`. The remaining elements beyond `k` are not important.
- **Input:** An integer array `nums` and an integer `val`.
- **Output:** An integer `k`.
- **Constraints:**
    - `0 <= nums.length <= 100`
    - `0 <= nums[i] <= 50`
    - `0 <= val <= 100`
- **Example:** `Input: nums = [3,2,2,3], val = 3` → `Output: 2` (The array becomes `[2,2,_,_]`).

---

### 9. Merge Sorted Array
- **Category:** Backwards Two Pointers (End to Start)
- **Problem Statement:** You are given two integer arrays `nums1` and `nums2`, sorted in non-decreasing order, and two integers `m` and `n`, representing the number of elements in `nums1` and `nums2` respectively. Merge `nums1` and `nums2` into a single array sorted in non-decreasing order. The final sorted array should be stored inside `nums1` (which has a length of `m + n`). The first `m` elements of `nums1` are the elements to be merged, and the last `n` elements are set to `0` and should be ignored. `nums2` has a length of `n`.
- **Input:** Array `nums1`, integer `m`, array `nums2`, integer `n`.
- **Output:** The merged array stored in `nums1` (returned as void).
- **Constraints:**
    - `nums1.length == m + n`
    - `nums2.length == n`
    - `0 <= m, n <= 200`
    - `1 <= m + n <= 200`
    - `-10^9 <= nums1[i], nums2[j] <= 10^9`
- **Example:** `Input: nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3` → `Output: [1,2,2,3,5,6]`

---

### 10. Longest Substring Without Repeating Characters
- **Category:** Sliding Window (Fast & Slow with Hash Set)
- **Problem Statement:** Given a string `s`, find the length of the **longest substring** without repeating characters.
- **Input:** A single string `s`.
- **Output:** An integer representing the maximum length.
- **Constraints:**
    - `0 <= s.length <= 5 * 10^4`
    - `s` consists of English letters, digits, symbols, and spaces.
- **Example:** `Input: s = "abcabcbb"` → `Output: 3` (Explanation: The answer is "abc", with the length of 3).

---

### 11. Minimum Window Substring
- **Category:** Sliding Window (Fast & Slow with frequency maps)
- **Problem Statement:** Given two strings `s` and `t` of lengths `m` and `n` respectively, return the **minimum window substring** of `s` such that every character in `t` (including duplicates) is included in the window. If there is no such substring, return the empty string `""`. The testcases will be generated such that the answer is unique.
- **Input:** Two strings, `s` and `t`.
- **Output:** A string (the minimum window).
- **Constraints:**
    - `m == s.length`, `n == t.length`
    - `1 <= m, n <= 10^5`
    - `s` and `t` consist of uppercase and lowercase English letters.
- **Example:** `Input: s = "ADOBECODEBANC", t = "ABC"` → `Output: "BANC"`

---

### 12. Fruit Into Baskets
- **Category:** Sliding Window (Longest subarray with at most 2 distinct types)
- **Problem Statement:** You are visiting a farm that has a row of fruit trees arranged from left to right. The trees are represented by an integer array `fruits` where `fruits[i]` is the type of fruit the `ith` tree produces. You want to collect as much fruit as possible. However, the owner has some strict rules: You only have **two baskets**, and each basket can only hold a **single type** of fruit. There is no limit on the amount of fruit each basket can hold. Starting from any tree, you must pick exactly one fruit from every tree (including the start tree) while moving to the right. The picked fruits must fit in one of the baskets. Return the **maximum number of fruits** you can pick.
- **Input:** An integer array `fruits`.
- **Output:** An integer representing the maximum total fruits.
- **Constraints:**
    - `1 <= fruits.length <= 10^5`
    - `0 <= fruits[i] < fruits.length`
- **Example:** `Input: fruits = [3,3,3,1,2,1,1,2,3,3,4]` → `Output: 5` (Explanation: The longest subarray with at most 2 distinct types is `[1,2,1,1,2]`).

---

### 13. Longest Palindromic Substring
- **Category:** Expand Around Center (Two Pointers moving outward)
- **Problem Statement:** Given a string `s`, return the longest palindromic substring in `s`.
- **Input:** A single string `s`.
- **Output:** A string (the longest palindrome).
- **Constraints:**
    - `1 <= s.length <= 1000`
    - `s` consists of only digits and English letters.
- **Example:** `Input: s = "babad"` → `Output: "bab"` (Note: "aba" is also a valid answer).

---

### 14. 4Sum
- **Category:** Opposite Direction (with nested fixed pivots)
- **Problem Statement:** Given an array `nums` of `n` integers, return an array of all the **unique** quadruplets `[nums[a], nums[b], nums[c], nums[d]]` such that:
    - `0 <= a, b, c, d < n`
    - `a`, `b`, `c`, and `d` are distinct.
    - `nums[a] + nums[b] + nums[c] + nums[d] == target`
    - The solution set must not contain duplicate quadruplets.
- **Input:** An integer array `nums` and an integer `target`.
- **Output:** A 2D list containing all unique quadruplets.
- **Constraints:**
    - `1 <= nums.length <= 200`
    - `-10^9 <= nums[i] <= 10^9`
    - `-10^9 <= target <= 10^9`
- **Example:** `Input: nums = [1,0,-1,0,-2,2], target = 0` → `Output: [[-2,-1,1,2],[-2,0,0,2],[-1,0,0,1]]`

---

### 15. Sort Colors (Dutch National Flag Problem)
- **Category:** Three Pointers (variant of two pointers for partitioning)
- **Problem Statement:** Given an array `nums` with `n` objects colored red, white, or blue, sort them **in-place** so that objects of the same color are adjacent, with the colors in the order red, white, and blue. We will use the integers `0`, `1`, and `2` to represent the color red, white, and blue respectively. You must solve this problem without using the library's sort function.
- **Input:** An integer array `nums` consisting solely of `0`s, `1`s, and `2`s.
- **Output:** The sorted array `nums` (modified in-place).
- **Constraints:**
    - `n == nums.length`
    - `1 <= n <= 300`
    - `nums[i]` is either `0`, `1`, or `2`.
- **Example:** `Input: nums = [2,0,2,1,1,0]` → `Output: [0,0,1,1,2,2]`


---

## Sliding Window
---

Here is a comprehensive list of **25 most frequently asked Sliding Window problems** in coding interviews. They are categorized into **Fixed-Size Windows** (where the window length `k` is given) and **Variable-Size Windows** (where you expand and shrink dynamically to find an optimal condition).

For each problem, I have provided a detailed breakdown of the statement, input/output specifications, constraints, and examples—**strictly without including the solution or algorithmic approach.**

---

### Fixed-Size Window Problems (Window length is explicitly given)

---

#### 1. Maximum Average Subarray I
- **Category:** Fixed-Size Window
- **Problem Statement:** You are given an integer array `nums` consisting of `n` elements, and an integer `k`. Find a contiguous subarray whose **length is equal to `k`** that has the maximum average value. Return the maximum average. Any answer within `10^-5` of the actual answer will be accepted.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** A floating-point number representing the maximum average.
- **Constraints:**
    - `n == nums.length`
    - `1 <= k <= n <= 10^5`
    - `-10^4 <= nums[i] <= 10^4`
- **Example:** `Input: nums = [1,12,-5,-6,50,3], k = 4` → `Output: 12.75000` (Explanation: Window `[12,-5,-6,50]` has average `12.75`).

---

#### 2. Maximum Sum Subarray of Size K
- **Category:** Fixed-Size Window
- **Problem Statement:** Given an array of positive integers `arr` and an integer `k`, find the maximum sum of any contiguous subarray of length exactly `k`.
- **Input:** An integer array `arr` and an integer `k`.
- **Output:** An integer representing the maximum sum.
- **Constraints:**
    - `1 <= k <= arr.length <= 10^5`
    - `1 <= arr[i] <= 10^4`
- **Example:** `Input: arr = [2, 1, 5, 1, 3, 2], k = 3` → `Output: 9` (Explanation: Subarray `[5, 1, 3]` sums to `9`).

---

#### 3. Number of Sub-arrays of Size K and Average Greater than or Equal to Threshold
- **Category:** Fixed-Size Window
- **Problem Statement:** Given an array of integers `arr` and two integers `k` and `threshold`, return the number of sub-arrays of size `k` whose average is greater than or equal to `threshold`.
- **Input:** An integer array `arr`, an integer `k`, and an integer `threshold`.
- **Output:** An integer representing the count.
- **Constraints:**
    - `1 <= arr.length <= 10^5`
    - `1 <= k <= arr.length`
    - `0 <= arr[i] <= 10^4`
    - `0 <= threshold <= 10^4`
- **Example:** `Input: arr = [2,2,2,2,5,5,5,8], k = 3, threshold = 4` → `Output: 3` (Explanation: Subarrays `[2,5,5]`, `[5,5,5]`, and `[5,5,8]` have averages `4`, `5`, and `6` respectively).

---

#### 4. Maximum Number of Vowels in a Substring of Given Length
- **Category:** Fixed-Size Window
- **Problem Statement:** Given a string `s` and an integer `k`, return the maximum number of vowel letters (`'a'`, `'e'`, `'i'`, `'o'`, `'u'`) in any substring of `s` with length exactly `k`.
- **Input:** A string `s` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= s.length <= 10^5`
    - `s` consists of lowercase English letters.
    - `1 <= k <= s.length`
- **Example:** `Input: s = "abciiidef", k = 3` → `Output: 3` (Explanation: The substring `"iii"` contains 3 vowels).

---

#### 5. Permutation in String
- **Category:** Fixed-Size Window (Character Frequency Matching)
- **Problem Statement:** Given two strings `s1` and `s2`, return `true` if `s2` contains a permutation of `s1`, or `false` otherwise. In other words, return `true` if one of `s1`'s permutations is the substring of `s2`.
- **Input:** Two strings, `s1` and `s2`.
- **Output:** A boolean (`true` or `false`).
- **Constraints:**
    - `1 <= s1.length, s2.length <= 10^4`
    - `s1` and `s2` consist of lowercase English letters.
- **Example:** `Input: s1 = "ab", s2 = "eidbaooo"` → `Output: true` (Explanation: `s2` contains `"ba"` which is a permutation of `"ab"`).

---

#### 6. Find All Anagrams in a String
- **Category:** Fixed-Size Window (Character Frequency Matching)
- **Problem Statement:** Given two strings `s` and `p`, return an array of all the start indices of `p`'s anagrams in `s`. You may return the answer in any order. An anagram is a word or phrase formed by rearranging the letters of a different word, typically using all the original letters exactly once.
- **Input:** A string `s` and a string `p`.
- **Output:** An integer array containing the start indices.
- **Constraints:**
    - `1 <= s.length, p.length <= 3 * 10^4`
    - `s` and `p` consist of lowercase English letters.
- **Example:** `Input: s = "cbaebabacd", p = "abc"` → `Output: [0,6]` (Explanation: Substrings starting at index `0` (`"cba"`) and index `6` (`"bac"`) are anagrams).

---

#### 7. Substrings of Size Three with Distinct Characters
- **Category:** Fixed-Size Window
- **Problem Statement:** A string is **good** if there are no repeated characters. Given a string `s`, return the number of **good** substrings of length exactly **three** in `s`. Note that if there are multiple occurrences of the same substring, every occurrence should be counted.
- **Input:** A string `s`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= s.length <= 100`
    - `s` consists of lowercase English letters.
- **Example:** `Input: s = "xyzzaz"` → `Output: 1` (Explanation: Only the substring `"xyz"` has distinct characters; `"yzz"` and `"zza"` do not).

---

#### 8. Grumpy Bookstore Owner
- **Category:** Fixed-Size Window (with optimization)
- **Problem Statement:** There is a bookstore owner that has a store open for `n` minutes. Every minute, some number of customers enter the store. You are given an integer array `customers` of length `n` where `customers[i]` is the number of customers entering at the start of the `i-th` minute, and an integer array `grumpy` where `grumpy[i]` is `1` if the owner is grumpy at the `i-th` minute, and `0` otherwise. When the owner is grumpy, the customers at that minute are not satisfied. Otherwise, they are satisfied. The owner has a secret technique to keep themselves **not grumpy** for `minutes` consecutive minutes, but they can only use it once. Return the maximum number of customers that can be satisfied throughout the day.
- **Input:** An integer array `customers`, an integer array `grumpy`, and an integer `minutes`.
- **Output:** An integer.
- **Constraints:**
    - `n == customers.length == grumpy.length`
    - `1 <= minutes <= n <= 2 * 10^4`
    - `0 <= customers[i] <= 1000`
    - `grumpy[i]` is either `0` or `1`.
- **Example:** `Input: customers = [1,0,1,2,1,1,7,5], grumpy = [0,1,0,1,0,1,0,1], minutes = 3` → `Output: 16`

---

#### 9. K Radius Subarray Averages
- **Category:** Fixed-Size Window (with middle element)
- **Problem Statement:** You are given a 0-indexed array `nums` of `n` integers, and an integer `k`. The **k-radius average** for a subarray of `nums` centered at some index `i` with the radius `k` is the average of all elements in `nums` between the indices `i - k` and `i + k` (inclusive). If there are less than `k` elements before or after the index `i`, then the k-radius average is `-1`. Build and return an array `avgs` of length `n` where `avgs[i]` is the k-radius average for the subarray centered at index `i`.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer array `avgs` (floor of the average).
- **Constraints:**
    - `n == nums.length`
    - `1 <= n <= 10^5`
    - `0 <= k <= n`
    - `0 <= nums[i] <= 10^5`
- **Example:** `Input: nums = [7,4,3,9,1,8,5,2,6], k = 3` → `Output: [-1,-1,-1,5,4,4,-1,-1,-1]`

---

### Variable-Size Window Problems (Window expands/shrinks to meet a condition)

---

#### 10. Longest Substring Without Repeating Characters
- **Category:** Variable-Size Window
- **Problem Statement:** Given a string `s`, find the length of the **longest substring** without repeating characters.
- **Input:** A single string `s`.
- **Output:** An integer representing the maximum length.
- **Constraints:**
    - `0 <= s.length <= 5 * 10^4`
    - `s` consists of English letters, digits, symbols, and spaces.
- **Example:** `Input: s = "abcabcbb"` → `Output: 3` (Explanation: The answer is `"abc"`, with length 3).

---

#### 11. Longest Repeating Character Replacement
- **Category:** Variable-Size Window
- **Problem Statement:** You are given a string `s` and an integer `k`. You can choose any character of the string and change it to any other uppercase English character. You may perform this operation at most `k` times. Return the length of the **longest substring** containing the same letter you can get after performing the above operations.
- **Input:** A string `s` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= s.length <= 10^5`
    - `s` consists of only uppercase English letters.
    - `0 <= k <= s.length`
- **Example:** `Input: s = "ABAB", k = 2` → `Output: 4` (Explanation: Replace the two `'B'`s with `'A'`s or vice versa).

---

#### 12. Minimum Window Substring
- **Category:** Variable-Size Window (Hard)
- **Problem Statement:** Given two strings `s` and `t` of lengths `m` and `n` respectively, return the **minimum window substring** of `s` such that every character in `t` (including duplicates) is included in the window. If there is no such substring, return the empty string `""`. The testcases will be generated such that the answer is unique.
- **Input:** Two strings, `s` and `t`.
- **Output:** A string (the minimum window).
- **Constraints:**
    - `m == s.length`, `n == t.length`
    - `1 <= m, n <= 10^5`
    - `s` and `t` consist of uppercase and lowercase English letters.
- **Example:** `Input: s = "ADOBECODEBANC", t = "ABC"` → `Output: "BANC"`

---

#### 13. Minimum Size Subarray Sum
- **Category:** Variable-Size Window
- **Problem Statement:** Given an array of positive integers `nums` and a positive integer `target`, return the **minimal length** of a contiguous subarray `[nums[l], nums[l+1], ..., nums[r-1], nums[r]]` of which the sum is greater than or equal to `target`. If there is no such subarray, return `0` instead.
- **Input:** An integer `target` and an array of integers `nums`.
- **Output:** An integer representing the minimal length.
- **Constraints:**
    - `1 <= target <= 10^9`
    - `1 <= nums.length <= 10^5`
    - `1 <= nums[i] <= 10^4`
- **Example:** `Input: target = 7, nums = [2,3,1,2,4,3]` → `Output: 2` (Explanation: `[4,3]` has length 2).

---

#### 14. Fruit Into Baskets
- **Category:** Variable-Size Window (At most 2 distinct types)
- **Problem Statement:** You are visiting a farm that has a row of fruit trees represented by an integer array `fruits` where `fruits[i]` is the type of fruit. You have **two baskets**, and each basket can only hold a **single type** of fruit. There is no limit on the amount. Starting from any tree, you must pick exactly one fruit from every tree while moving right, but you can only pick fruits that fit in one of the two baskets. Return the **maximum number of fruits** you can pick.
- **Input:** An integer array `fruits`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= fruits.length <= 10^5`
    - `0 <= fruits[i] < fruits.length`
- **Example:** `Input: fruits = [3,3,3,1,2,1,1,2,3,3,4]` → `Output: 5` (Longest subarray with ≤2 distinct types is `[1,2,1,1,2]`).

---

#### 15. Longest Substring with At Most K Distinct Characters
- **Category:** Variable-Size Window
- **Problem Statement:** Given a string `s` and an integer `k`, return the length of the **longest substring** of `s` that contains at most `k` **distinct** characters.
- **Input:** A string `s` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= s.length <= 5 * 10^4`
    - `s` consists of English letters and digits.
    - `0 <= k <= s.length`
- **Example:** `Input: s = "eceba", k = 2` → `Output: 3` (Substring is `"ece"`).

---

#### 16. Max Consecutive Ones III
- **Category:** Variable-Size Window (At most K zeros)
- **Problem Statement:** Given a binary array `nums` and an integer `k`, return the maximum number of consecutive `1`'s in the array if you can flip at most `k` `0`'s to `1`'s.
- **Input:** An integer array `nums` (containing only `0`s and `1`s) and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `nums[i]` is `0` or `1`.
    - `0 <= k <= nums.length`
- **Example:** `Input: nums = [1,1,1,0,0,0,1,1,1,1,0], k = 2` → `Output: 6`

---

#### 17. Subarray Product Less Than K
- **Category:** Variable-Size Window (Counting subarrays)
- **Problem Statement:** Given an array of positive integers `nums` and an integer `k`, return the number of contiguous subarrays where the product of all the elements in the subarray is strictly less than `k`.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 3 * 10^4`
    - `1 <= nums[i] <= 1000`
    - `0 <= k <= 10^6`
- **Example:** `Input: nums = [10,5,2,6], k = 100` → `Output: 8`

---

#### 18. Count Number of Nice Subarrays
- **Category:** Variable-Size Window (Exactly K odd numbers)
- **Problem Statement:** Given an array of integers `nums` and an integer `k`. A continuous subarray is called **nice** if there are exactly `k` odd numbers in it. Return the number of **nice** subarrays.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 50000`
    - `1 <= nums[i] <= 10^5`
    - `1 <= k <= nums.length`
- **Example:** `Input: nums = [1,1,2,1,1], k = 3` → `Output: 2` (Explanation: Subarrays `[1,1,2,1]` and `[1,2,1,1]`).

---

#### 19. Binary Subarrays With Sum
- **Category:** Variable-Size Window (Exactly sum K)
- **Problem Statement:** Given a binary array `nums` (containing only `0`s and `1`s) and an integer `goal`, return the number of non-empty **subarrays** with a sum equal to `goal`.
- **Input:** An integer array `nums` and an integer `goal`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 3 * 10^4`
    - `nums[i]` is `0` or `1`.
    - `0 <= goal <= nums.length`
- **Example:** `Input: nums = [1,0,1,0,1], goal = 2` → `Output: 4`

---

#### 20. Longest Subarray of 1's After Deleting One Element
- **Category:** Variable-Size Window (At most one zero)
- **Problem Statement:** Given a binary array `nums`, you should delete one element from it. Return the size of the longest non-empty subarray containing only `1`'s in the resulting array. Return `0` if there is no such subarray.
- **Input:** A binary integer array `nums`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `nums[i]` is either `0` or `1`.
- **Example:** `Input: nums = [1,1,0,1]` → `Output: 3` (Explanation: Deleting the `0` leaves `[1,1,1]`).

---

#### 21. Frequency of the Most Frequent Element
- **Category:** Variable-Size Window (with Sorting)
- **Problem Statement:** The **frequency** of an element is the number of times it occurs in an array. You are given an integer array `nums` and an integer `k`. In one operation, you can choose an index of `nums` and increment the element at that index by `1`. Return the **maximum possible frequency** of an element after performing **at most** `k` operations.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `1 <= nums[i] <= 10^5`
    - `1 <= k <= 10^5`
- **Example:** `Input: nums = [1,2,4], k = 5` → `Output: 3` (Explanation: Increment `1` to `4` (cost 3) and `2` to `4` (cost 2), total cost 5, making all elements `4`).

---

### Advanced / Hard Variable-Size Window Problems (often requiring auxiliary data structures)

---

#### 22. Shortest Subarray with Sum at Least K
- **Category:** Variable-Size Window (requires Monotonic Deque due to negative numbers)
- **Problem Statement:** Given an integer array `nums` and an integer `k`, return the **length of the shortest non-empty subarray** of `nums` with a sum of at least `k`. If there is no such subarray, return `-1`. A subarray is a contiguous part of an array.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `-10^5 <= nums[i] <= 10^5`
    - `1 <= k <= 10^9`
- **Example:** `Input: nums = [2,-1,2], k = 3` → `Output: 3` (Explanation: The subarray `[2,-1,2]` sums to 3).

---

#### 23. Longest Continuous Subarray With Absolute Diff Less Than or Equal to Limit
- **Category:** Variable-Size Window (requires Monotonic Deques for min/max)
- **Problem Statement:** Given an array of integers `nums` and an integer `limit`, return the size of the **longest non-empty subarray** such that the absolute difference between the maximum and minimum element in the subarray is less than or equal to `limit`.
- **Input:** An integer array `nums` and an integer `limit`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `0 <= nums[i] <= 10^9`
    - `0 <= limit <= 10^9`
- **Example:** `Input: nums = [8,2,4,7], limit = 4` → `Output: 2` (Explanation: `[2,4]` and `[4,7]` are valid, length 2).

---

#### 24. Replace the Substring for Balanced String
- **Category:** Variable-Size Window (Hard)
- **Problem Statement:** A string is called **balanced** if each of the characters `'Q'`, `'W'`, `'E'`, and `'R'` appears exactly `n/4` times where `n` is the length of the string. You are given a string `s` of length `n` that is guaranteed to be a multiple of `4`. In one operation, you can replace a **substring** with any other string of the same length that contains only the above characters. Return the **minimum length** of the substring that you can replace to make `s` balanced.
- **Input:** A string `s`.
- **Output:** An integer.
- **Constraints:**
    - `n == s.length`
    - `4 <= n <= 10^5`
    - `n` is a multiple of `4`.
    - `s` consists only of `'Q'`, `'W'`, `'E'`, and `'R'`.
- **Example:** `Input: s = "WWEQERQW"` → `Output: 3` (Explanation: Replacing `"EQE"` with `"RQQ"` balances the string).

---

#### 25. Count Subarrays with Fixed Bounds
- **Category:** Variable-Size Window (Hard – tracking min/max positions)
- **Problem Statement:** You are given an integer array `nums` and two integers `minK` and `maxK`. A **fixed-bound subarray** of `nums` is a subarray that satisfies the following conditions:
    - The **minimum** value in the subarray is equal to `minK`.
    - The **maximum** value in the subarray is equal to `maxK`.
    Return the **number of fixed-bound subarrays**.
- **Input:** An integer array `nums`, an integer `minK`, and an integer `maxK`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `1 <= nums[i], minK, maxK <= 10^6`
- **Example:** `Input: nums = [1,3,5,2,7,5], minK = 1, maxK = 5` → `Output: 2` (Explanation: The subarrays `[1,3,5]` and `[1,3,5,2]` satisfy the conditions).

---
---

Here is a comprehensive list of **15 most frequently asked Prefix Sum problems** in coding interviews. They are categorized by complexity, ranging from basic 1D/2D prefix arrays to advanced combinations with hash maps and modulo arithmetic.

For each problem, I have provided a detailed breakdown of the statement, input/output specifications, constraints, and examples—**strictly without including the solution or algorithmic approach.**

---

### Basic 1D Prefix Sum Problems

---

#### 1. Range Sum Query - Immutable
- **Category:** Basic 1D Prefix Sum
- **Problem Statement:** Given an integer array `nums`, handle multiple queries of the following type:
    - Calculate the **sum** of the elements of `nums` between indices `left` and `right` **inclusive** where `left <= right`.
    - Implement the `NumArray` class:
        - `NumArray(int[] nums)` initializes the object with the integer array `nums`.
        - `int sumRange(int left, int right)` returns the sum of the elements of `nums` from index `left` to index `right` inclusive.
- **Input:** An integer array `nums` for initialization, followed by multiple query pairs `(left, right)`.
- **Output:** An integer for each `sumRange` call, representing the sum.
- **Constraints:**
    - `1 <= nums.length <= 10^4`
    - `-10^5 <= nums[i] <= 10^5`
    - `0 <= left <= right < nums.length`
    - At most `10^4` calls will be made to `sumRange`.
- **Example:** `Input: nums = [-2, 0, 3, -5, 2, -1]`, `sumRange(0, 2)` → `Output: 1` (Explanation: `-2 + 0 + 3 = 1`).

---

#### 2. Running Sum of 1d Array
- **Category:** Basic 1D Prefix Sum (In-place)
- **Problem Statement:** Given an array `nums`, define a running sum of an array as `runningSum[i] = sum(nums[0]...nums[i])`. Return the running sum of `nums`.
- **Input:** An integer array `nums`.
- **Output:** An integer array of the same length containing the running sums.
- **Constraints:**
    - `1 <= nums.length <= 1000`
    - `-10^6 <= nums[i] <= 10^6`
- **Example:** `Input: nums = [1,2,3,4]` → `Output: [1,3,6,10]`

---

#### 3. Find Pivot Index
- **Category:** Basic 1D Prefix Sum (Left/Right equality)
- **Problem Statement:** Given an array of integers `nums`, calculate the **pivot index** of this array. The pivot index is the index where the sum of all the numbers **strictly** to the left of the index is equal to the sum of all the numbers **strictly** to the right of the index. If the index is on the left edge of the array, then the left sum is `0` because there are no elements to the left. This also applies to the right edge. Return the **leftmost pivot index**. If no such index exists, return `-1`.
- **Input:** An integer array `nums`.
- **Output:** An integer (the pivot index or `-1`).
- **Constraints:**
    - `1 <= nums.length <= 10^4`
    - `-1000 <= nums[i] <= 1000`
- **Example:** `Input: nums = [1,7,3,6,5,6]` → `Output: 3` (Explanation: Left sum = `1+7+3 = 11`, Right sum = `5+6 = 11`).

---

### 2D Prefix Sum Problems

---

#### 4. Range Sum Query 2D - Immutable
- **Category:** 2D Prefix Sum
- **Problem Statement:** Given a 2D matrix `matrix`, handle multiple queries of the following type:
    - Calculate the **sum** of the elements of `matrix` inside the rectangle defined by its upper left corner `(row1, col1)` and lower right corner `(row2, col2)`.
    - Implement the `NumMatrix` class:
        - `NumMatrix(int[][] matrix)` initializes the object with the integer matrix.
        - `int sumRegion(int row1, int col1, int row2, int col2)` returns the sum of the elements of the matrix inside the rectangle.
- **Input:** A 2D integer matrix for initialization, followed by multiple query coordinates.
- **Output:** An integer for each `sumRegion` call, representing the rectangular sum.
- **Constraints:**
    - `m == matrix.length`, `n == matrix[i].length`
    - `1 <= m, n <= 200`
    - `-10^4 <= matrix[i][j] <= 10^4`
    - `0 <= row1 <= row2 < m`
    - `0 <= col1 <= col2 < n`
    - At most `10^4` calls will be made to `sumRegion`.
- **Example:** `Input: matrix = [[3,0,1,4,2],[5,6,3,2,1],[1,2,0,1,5],[4,1,0,1,7],[1,0,3,0,5]]`, `sumRegion(2,1,4,3)` → `Output: 8` (sum of rectangle from row 2 to 4, col 1 to 3).

---

#### 5. Number of Submatrices That Sum to Target
- **Category:** 2D Prefix Sum + Hash Map
- **Problem Statement:** Given a `matrix` and a `target`, return the number of non-empty submatrices that sum to the target. A submatrix is formed by selecting any contiguous rows and any contiguous columns.
- **Input:** A 2D integer array `matrix` and an integer `target`.
- **Output:** An integer representing the count.
- **Constraints:**
    - `1 <= matrix.length <= 100`
    - `1 <= matrix[0].length <= 100`
    - `-1000 <= matrix[i][j] <= 1000`
    - `-10^8 <= target <= 10^8`
- **Example:** `Input: matrix = [[0,1,0],[1,1,1],[0,1,0]], target = 0` → `Output: 4` (Explanation: There are 4 submatrices that sum to 0: four 1x1 submatrices containing `0`).

---

### Prefix Sum + Hash Map (Subarray Sum Problems)

---

#### 6. Subarray Sum Equals K
- **Category:** Prefix Sum + Hash Map (Classic)
- **Problem Statement:** Given an array of integers `nums` and an integer `k`, return the total number of subarrays whose sum equals to `k`. A subarray is a contiguous non-empty sequence of elements within the array.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer (the count).
- **Constraints:**
    - `1 <= nums.length <= 2 * 10^4`
    - `-1000 <= nums[i] <= 1000`
    - `-10^7 <= k <= 10^7`
- **Example:** `Input: nums = [1,1,1], k = 2` → `Output: 2` (Explanation: Subarrays `[1,1]` at indices 0-1 and 1-2).

---

#### 7. Contiguous Array
- **Category:** Prefix Sum + Hash Map (Treat 0 as -1)
- **Problem Statement:** Given a binary array `nums` (containing only `0` and `1`), return the maximum length of a contiguous subarray with an **equal number** of `0` and `1`.
- **Input:** A binary integer array `nums`.
- **Output:** An integer (the maximum length).
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `nums[i]` is either `0` or `1`.
- **Example:** `Input: nums = [0,1,0]` → `Output: 2` (Explanation: `[0,1]` or `[1,0]` both have equal 0s and 1s).

---

#### 8. Subarray Sums Divisible by K
- **Category:** Prefix Sum + Modulo + Hash Map
- **Problem Statement:** Given an integer array `nums` and an integer `k`, return the number of non-empty subarrays that have a sum divisible by `k`. A subarray is a contiguous part of an array.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 3 * 10^4`
    - `-10^4 <= nums[i] <= 10^4`
    - `2 <= k <= 10^4`
- **Example:** `Input: nums = [4,5,0,-2,-3,1], k = 5` → `Output: 7` (Explanation: There are 7 subarrays with sums divisible by 5).

---

#### 9. Maximum Size Subarray Sum Equals k
- **Category:** Prefix Sum + Hash Map (Earliest occurrence)
- **Problem Statement:** Given an integer array `nums` and an integer `k`, return the **maximum length** of a contiguous subarray that sums to `k`. If there is no such subarray, return `0`. (Note: This is a classic interview problem, commonly asked by Facebook/Google).
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** An integer (the maximum length).
- **Constraints:**
    - `1 <= nums.length <= 2 * 10^4`
    - `-10^4 <= nums[i] <= 10^4`
    - `-10^9 <= k <= 10^9`
- **Example:** `Input: nums = [1, -1, 5, -2, 3], k = 3` → `Output: 4` (Explanation: Subarray `[1, -1, 5, -2]` sums to 3 and has length 4).

---

#### 10. Make Sum Divisible by P
- **Category:** Prefix Sum + Modulo + Hash Map (Advanced)
- **Problem Statement:** Given an array of positive integers `nums`, remove the **smallest** subarray (contiguous) such that the sum of the remaining elements is divisible by `p`. It is **not** allowed to remove the whole array. Return the length of the smallest subarray that you need to remove, or `-1` if it's impossible. If the total sum is already divisible by `p`, you should return `0`.
- **Input:** An integer array `nums` and an integer `p`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `1 <= nums[i] <= 10^9`
    - `1 <= p <= 10^9`
- **Example:** `Input: nums = [3,1,4,2], p = 6` → `Output: 1` (Explanation: Sum of array is `10`, need to remove a subarray with sum `4` (e.g., `[4]`) to make remaining sum `6` divisible by `6`).

---

### Prefix/Suffix Product Problems

---

#### 11. Product of Array Except Self
- **Category:** Prefix Product + Suffix Product
- **Problem Statement:** Given an integer array `nums`, return an array `answer` such that `answer[i]` is equal to the product of all the elements of `nums` except `nums[i]`. The product of any prefix or suffix of `nums` is **guaranteed** to fit in a 32-bit integer. You must write an algorithm that runs in `O(n)` time and **without using the division operation**.
- **Input:** An integer array `nums`.
- **Output:** An integer array `answer`.
- **Constraints:**
    - `2 <= nums.length <= 10^5`
    - `-30 <= nums[i] <= 30`
    - The product of any prefix or suffix fits in a 32-bit integer.
- **Example:** `Input: nums = [1,2,3,4]` → `Output: [24,12,8,6]`

---

#### 12. Maximum Subarray Sum with One Deletion
- **Category:** Prefix Sum + Suffix Sum (or DP, but commonly solved with prefix/suffix)
- **Problem Statement:** Given an array of integers, return the maximum sum for a **non-empty** subarray (contiguous) with **at most one** element deletion. In other words, you want to choose a subarray and optionally delete one element from it so that the sum of the remaining elements is maximized. Return the maximum sum.
- **Input:** An integer array `nums`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `-10^4 <= nums[i] <= 10^4`
- **Example:** `Input: nums = [1,-2,0,3]` → `Output: 4` (Explanation: Choose `[1,-2,0,3]` and delete `-2`, sum = `1+0+3 = 4`).

---

### Advanced / Applications of Prefix Sum

---

#### 13. Longest Well-Performing Interval
- **Category:** Prefix Sum + Hash Map (Treat >8 as +1, <=8 as -1)
- **Problem Statement:** We are given `hours`, a list of the number of hours worked per day for a given employee. A day is considered a **tiring day** if and only if the number of hours worked is **strictly greater than 8**. A **well-performing interval** is an interval of days for which the number of tiring days is **strictly larger** than the number of non-tiring days. Return the length of the longest well-performing interval.
- **Input:** An integer array `hours`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= hours.length <= 10^4`
    - `0 <= hours[i] <= 16`
- **Example:** `Input: hours = [9,9,6,0,6,6,9]` → `Output: 3` (Explanation: The longest interval is `[9,9,6]` or `[9,6,9]` where tiring days > non-tiring days).

---

#### 14. Maximum Points You Can Obtain from Cards
- **Category:** Prefix Sum (or Sliding Window from ends)
- **Problem Statement:** There are several cards arranged in a row, and each card has an associated number of points. The points are given in the integer array `cardPoints`. In one step, you can take one card from the beginning or from the end of the row. You have to take exactly `k` cards. Your score is the sum of the points of the cards you have taken. Return the **maximum score** you can obtain.
- **Input:** An integer array `cardPoints` and an integer `k`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= cardPoints.length <= 10^5`
    - `1 <= cardPoints[i] <= 10^4`
    - `1 <= k <= cardPoints.length`
- **Example:** `Input: cardPoints = [1,2,3,4,5,6,1], k = 3` → `Output: 12` (Explanation: Take `[1,6,5]` from ends – pick left 1, right 6, right 5 → sum 12).

---

#### 15. Minimum Operations to Reduce X to Zero
- **Category:** Prefix Sum + Hash Map (Transform to longest subarray sum)
- **Problem Statement:** You are given an integer array `nums` and an integer `x`. In one operation, you can remove the leftmost or the rightmost element from the array and subtract its value from `x`. Note that this modifies the array for future operations. Return the **minimum number of operations** required to reduce `x` to **exactly `0`**. If it is impossible, return `-1`.
- **Input:** An integer array `nums` and an integer `x`.
- **Output:** An integer.
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `1 <= nums[i] <= 10^4`
    - `1 <= x <= 10^9`
- **Example:** `Input: nums = [1,1,4,2,3], x = 5` → `Output: 2` (Explanation: Remove `[3]` from the right and `[1,1]` from the left? Actually optimal: remove `[3]` and `[2]` from right → sum = 5 in 2 operations).

---
---
Here is a comprehensive list of **15 most frequently asked Array Rotation problems** in coding interviews. They are categorized into **1D Array Rotations**, **Searching in Rotated Arrays**, **2D Matrix Rotations**, and **Mathematical/Query-Based Rotations**.

For each problem, I have provided a detailed breakdown of the statement, input/output specifications, constraints, and examples—**strictly without including the solution or algorithmic approach.**

---

### 1D Array Rotation (Left / Right Shifts)

---

#### 1. Rotate Array (Right by K Steps)
- **Category:** 1D Rotation (Right)
- **Problem Statement:** Given an integer array `nums` and an integer `k`, rotate the array to the **right** by `k` steps, where `k` is non-negative. Right rotation means moving each element `k` positions to the end, and elements that go past the last position wrap around to the front.
- **Input:** An integer array `nums` and an integer `k`.
- **Output:** The modified array `nums` (in-place transformation; the function should return `void` but the array is mutated).
- **Constraints:**
    - `1 <= nums.length <= 10^5`
    - `-2^31 <= nums[i] <= 2^31 - 1`
    - `0 <= k <= 10^5`
- **Example:** `Input: nums = [1,2,3,4,5,6,7], k = 3` → `Output: [5,6,7,1,2,3,4]`

---

#### 2. Rotate Array (Left by K Steps)
- **Category:** 1D Rotation (Left)
- **Problem Statement:** Given an integer array `arr` and an integer `d`, rotate the array to the **left** by `d` steps. Left rotation means moving the first `d` elements to the end of the array, preserving their relative order. You must perform the rotation in `O(1)` extra space.
- **Input:** An integer array `arr` and an integer `d`.
- **Output:** The rotated array `arr` (modified in-place).
- **Constraints:**
    - `1 <= arr.length <= 10^5`
    - `0 <= d <= arr.length`
    - `arr[i]` fits in a 32-bit integer.
- **Example:** `Input: arr = [1,2,3,4,5], d = 2` → `Output: [3,4,5,1,2]`

---

#### 3. Rotate Array by One (Multiple Times)
- **Category:** 1D Rotation (Step-by-step)
- **Problem Statement:** Given an array `arr` of size `n` and a positive integer `k`, perform `k` **sequential** left rotations by one position each. Unlike a direct rotation, this problem specifically tests the naive step-by-step approach (often asked in beginner interviews to contrast with efficient methods). Return the array after all rotations.
- **Input:** An integer array `arr`, an integer `k` (number of rotations).
- **Output:** The final array after `k` rotations.
- **Constraints:**
    - `1 <= n <= 10^3` (usually smaller to allow O(n*k) but still a valid test).
    - `1 <= k <= 10^3`
- **Example:** `Input: arr = [1,2,3,4,5], k = 2` → `Output: [3,4,5,1,2]`

---

### Searching and Finding in Rotated Sorted Arrays

---

#### 4. Search in Rotated Sorted Array (Unique Elements)
- **Category:** Binary Search on Rotated Array
- **Problem Statement:** There is an integer array `nums` sorted in ascending order (with **distinct** values). Prior to being passed to your function, `nums` is **rotated** at an unknown pivot index. For example, `[0,1,2,4,5,6,7]` might become `[4,5,6,7,0,1,2]`. Given the array `nums` after the rotation and an integer `target`, return the index of `target` if it is in `nums`, or `-1` if it is not. You must write an algorithm with `O(log n)` runtime complexity.
- **Input:** An integer array `nums` (rotated sorted, unique) and an integer `target`.
- **Output:** An integer (the index or `-1`).
- **Constraints:**
    - `1 <= nums.length <= 5000`
    - `-10^4 <= nums[i] <= 10^4`
    - All values of `nums` are **unique**.
    - `nums` is guaranteed to be rotated at some pivot.
- **Example:** `Input: nums = [4,5,6,7,0,1,2], target = 0` → `Output: 4`

---

#### 5. Search in Rotated Sorted Array II (With Duplicates)
- **Category:** Binary Search on Rotated Array (with duplicates)
- **Problem Statement:** This problem is the same as the above, but the array `nums` may contain **duplicates**. Given the rotated array and a `target`, return `true` if `target` exists, or `false` otherwise. The presence of duplicates increases the worst-case complexity to `O(n)` but you should still aim for `O(log n)` in the average case.
- **Input:** An integer array `nums` (rotated sorted, may contain duplicates) and an integer `target`.
- **Output:** A boolean (`true` or `false`).
- **Constraints:**
    - `1 <= nums.length <= 5000`
    - `-10^4 <= nums[i] <= 10^4`
    - `nums` is guaranteed to be rotated at some pivot.
- **Example:** `Input: nums = [2,5,6,0,0,1,2], target = 0` → `Output: true`; `target = 3` → `Output: false`

---

#### 6. Find Minimum in Rotated Sorted Array (Unique Elements)
- **Category:** Binary Search for Minimum
- **Problem Statement:** Suppose an array of length `n` sorted in ascending order is **rotated** between `1` and `n` times. For example, the array `nums = [0,1,2,4,5,6,7]` might become `[4,5,6,7,0,1,2]`. Given the rotated array of **unique** elements, return the **minimum element** of this array. You must write an algorithm that runs in `O(log n)` time.
- **Input:** An integer array `nums` (rotated sorted, unique).
- **Output:** An integer (the minimum value).
- **Constraints:**
    - `1 <= nums.length <= 5000`
    - `-5000 <= nums[i] <= 5000`
    - All the integers of `nums` are **unique**.
    - `nums` is sorted and rotated.
- **Example:** `Input: nums = [3,4,5,1,2]` → `Output: 1`

---

#### 7. Find Minimum in Rotated Sorted Array II (With Duplicates)
- **Category:** Binary Search for Minimum (with duplicates)
- **Problem Statement:** This is the same as the previous problem, but the array `nums` may contain **duplicates**. Find the minimum element in the rotated sorted array. Even with duplicates, your algorithm should still attempt to run in `O(log n)` time on average, though worst-case can be `O(n)`.
- **Input:** An integer array `nums` (rotated sorted, may contain duplicates).
- **Output:** An integer (the minimum value).
- **Constraints:**
    - `1 <= nums.length <= 5000`
    - `-5000 <= nums[i] <= 5000`
    - `nums` is sorted and rotated.
- **Example:** `Input: nums = [2,2,2,0,1]` → `Output: 0`

---

#### 8. Find Pivot Index / Count of Rotations
- **Category:** Binary Search for Pivot
- **Problem Statement:** Given a rotated sorted array `arr` of **distinct** elements, find the **pivot index**. The pivot index is the index where the largest element resides, or equivalently, the number of times the original sorted array has been rotated. For example, in `[4,5,6,7,0,1,2]`, the pivot is at index `3` (value `7`), meaning the array has been rotated `4` times (or `1` time, depending on definition—typically, it's the index of the minimum element, but this problem asks for the count of rotations, i.e., the index of the minimum element). Return the number of rotations.
- **Input:** An integer array `arr` (rotated sorted, distinct).
- **Output:** An integer (the rotation count, i.e., the index of the minimum element).
- **Constraints:**
    - `1 <= arr.length <= 10^5`
    - `arr` consists of distinct integers.
- **Example:** `Input: arr = [15, 18, 2, 3, 6, 12]` → `Output: 2` (because the array was rotated twice to the left from `[2,3,6,12,15,18]`).

---

### Verification / Structure Checks

---

#### 9. Check if Array Is Sorted and Rotated
- **Category:** Verification (Single traversal)
- **Problem Statement:** Given an array `nums`, return `true` if the array was originally sorted in non-decreasing order, then rotated some number of positions (including zero). Otherwise, return `false`. Note that there may be **duplicates** in the array.
- **Input:** An integer array `nums`.
- **Output:** A boolean (`true` or `false`).
- **Constraints:**
    - `1 <= nums.length <= 100`
    - `1 <= nums[i] <= 100`
- **Example:** `Input: nums = [3,4,5,1,2]` → `Output: true`; `Input: nums = [2,1,3,4]` → `Output: false` (Because it cannot be obtained by rotating a sorted array).

---

### 2D Array / Matrix Rotations

---

#### 10. Rotate Image (90 Degrees Clockwise)
- **Category:** 2D Matrix Rotation (In-place)
- **Problem Statement:** You are given an `n x n` 2D `matrix` representing an image, rotate the image by **90 degrees (clockwise)**. You have to rotate the image **in-place**, which means you have to modify the input 2D matrix directly. **DO NOT** allocate another 2D matrix and do the rotation.
- **Input:** A 2D integer array `matrix` of size `n x n`.
- **Output:** The modified `matrix` (rotated in-place, return `void`).
- **Constraints:**
    - `n == matrix.length == matrix[i].length`
    - `1 <= n <= 20`
    - `-1000 <= matrix[i][j] <= 1000`
- **Example:** `Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]` → `Output: [[7,4,1],[8,5,2],[9,6,3]]`

---

#### 11. Rotate a Matrix by 90 Degrees Counter-Clockwise
- **Category:** 2D Matrix Rotation (In-place)
- **Problem Statement:** Given an `n x n` integer matrix, rotate it **90 degrees counter-clockwise** (anti-clockwise) in-place. Perform the rotation without using any extra space for another matrix.
- **Input:** A 2D integer array `matrix` of size `n x n`.
- **Output:** The modified `matrix` (rotated in-place).
- **Constraints:**
    - `1 <= n <= 20`
    - `-10^3 <= matrix[i][j] <= 10^3`
- **Example:** `Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]` → `Output: [[3,6,9],[2,5,8],[1,4,7]]`

---

#### 12. Rotate a Matrix by 180 Degrees
- **Category:** 2D Matrix Rotation (In-place)
- **Problem Statement:** Given an `n x n` matrix, rotate it by **180 degrees** in the clockwise direction. Perform this transformation in-place.
- **Input:** A 2D integer array `matrix` of size `n x n`.
- **Output:** The modified `matrix` (rotated in-place).
- **Constraints:**
    - `1 <= n <= 20`
    - `-10^4 <= matrix[i][j] <= 10^4`
- **Example:** `Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]` → `Output: [[9,8,7],[6,5,4],[3,2,1]]`

---

### Mathematical / Query-Based Rotations

---

#### 13. Maximum Sum of i * arr[i] Among All Rotations
- **Category:** Rotation Sum Optimization (DP / Math)
- **Problem Statement:** Given an array `arr` of `n` integers, consider all its **left rotations**. For each rotation, compute the sum of `i * arr[i]` for `i = 0` to `n-1`. Return the **maximum** such sum achievable over all possible rotations. You are allowed to rotate the array any number of times (including zero).
- **Input:** An integer array `arr`.
- **Output:** An integer (the maximum possible sum).
- **Constraints:**
    - `1 <= arr.length <= 10^5`
    - `-100 <= arr[i] <= 100`
- **Example:** `Input: arr = [8, 3, 1, 2]` → `Output: 29` (Explanation: The rotation `[3, 1, 2, 8]` gives `0*3 + 1*1 + 2*2 + 3*8 = 29`, which is maximum).

---

#### 14. Find a Pair with Given Sum in a Rotated Sorted Array
- **Category:** Two-Pointer on Rotated Array
- **Problem Statement:** Given a sorted and **rotated** array `arr` of **distinct** elements, and an integer `x`, find whether there exists a pair of elements in the array whose sum is exactly `x`. You must do this in `O(n)` time without using extra space.
- **Input:** An integer array `arr` (rotated sorted, distinct) and an integer `x`.
- **Output:** A boolean (`true` if such a pair exists, else `false`).
- **Constraints:**
    - `2 <= arr.length <= 10^5`
    - `-10^6 <= arr[i], x <= 10^6`
    - Elements are distinct.
- **Example:** `Input: arr = [11, 15, 6, 8, 9, 10], x = 16` → `Output: true` (Explanation: Pair `(6, 10)` sums to 16).

---

#### 15. Queries on Left and Right Circular Rotation
- **Category:** Rotation Query Processing (Prefix Sum)
- **Problem Statement:** You are given an array `arr` of size `n`. You are also given `q` queries of two types:
    - **Type 1 (L, R):** Left rotate the array by `L` positions, then find the sum of elements between indices `L` and `R` (inclusive) on the **rotated** array.
    - **Type 2 (L, R):** Right rotate the array by `L` positions, then find the sum of elements between indices `L` and `R` (inclusive) on the **rotated** array.
    - **Note:** The rotations are cumulative; each query is performed on the current state of the array. Return the sum for each query.
- **Input:** An integer array `arr`, an integer `q` (number of queries), followed by `q` queries each containing a string `type` ("left" or "right"), an integer `k` (rotation amount), and two integers `L` and `R`.
- **Output:** An integer for each query (the sum over the specified range after the rotation).
- **Constraints:**
    - `1 <= n, q <= 10^5`
    - `1 <= arr[i] <= 10^4`
    - `0 <= k <= 10^5`
    - `0 <= L <= R < n`
- **Example:** `Input: arr = [1, 2, 3, 4, 5], queries = [("left", 2, 0, 2)]` → After left rotation by 2, array becomes `[3,4,5,1,2]`. Sum from `0` to `2` is `3+4+5 = 12` → `Output: 12`
---
---

