# Input/Output Handling - Master Fast I/O

## Why Fast I/O Matters

**Real Scenario:**
- Problem: Read 1,000,000 integers and print them
- Scanner: **5-10 seconds** ❌ (Time Limit Exceeded!)
- BufferedReader: **0.5 seconds** ✅ (Accepted!)

In competitive programming, **I/O can be your bottleneck**. Let's master fast techniques!

---

## Method 1: Scanner (Beginner-Friendly)

### Basic Usage

```java
import java.util.Scanner;

public class ScannerBasic {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        
        // Read single values
        int n = sc.nextInt();
        long l = sc.nextLong();
        double d = sc.nextDouble();
        String word = sc.next();        // Reads until whitespace
        
        // Common mistake! Need to consume newline after nextInt()
        sc.nextLine();  // Clear buffer
        String line = sc.nextLine();    // Now reads full line
        
        sc.close();
    }
}
```

---

### Reading Arrays

```java
import java.util.Scanner;

public class ScannerArrays {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        
        int n = sc.nextInt();
        int[] arr = new int[n];
        
        // Read n integers
        for (int i = 0; i < n; i++) {
            arr[i] = sc.nextInt();
        }
        
        // Read 2D array (matrix)
        int rows = sc.nextInt();
        int cols = sc.nextInt();
        int[][] matrix = new int[rows][cols];
        
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                matrix[i][j] = sc.nextInt();
            }
        }
        
        sc.close();
    }
}
```

---

### When to Use Scanner

✅ **Good for:**
- Small inputs (n < 10,000)
- Quick testing
- Simple problems
- When clarity > speed

❌ **Avoid for:**
- Large inputs (n > 100,000)
- Multiple test cases with large data
- Time-critical problems

---

## Method 2: BufferedReader (Fast Input)

### Basic Template

```java
import java.io.*;
import java.util.*;

public class BufferedReaderBasic {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        // Read single line
        String line = br.readLine();
        
        // Parse single integer
        int n = Integer.parseInt(br.readLine());
        
        // Parse single long
        long l = Long.parseLong(br.readLine());
        
        // Parse single double
        double d = Double.parseDouble(br.readLine());
        
        br.close();
    }
}
```

---

### Reading Multiple Values from One Line

```java
import java.io.*;

public class ReadMultipleValues {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        // Input: "5 10 15"
        String[] tokens = br.readLine().split(" ");
        int a = Integer.parseInt(tokens[0]);  // 5
        int b = Integer.parseInt(tokens[1]);  // 10
        int c = Integer.parseInt(tokens[2]);  // 15
        
        // Or use StringTokenizer (slightly faster)
        StringTokenizer st = new StringTokenizer(br.readLine());
        int x = Integer.parseInt(st.nextToken());
        int y = Integer.parseInt(st.nextToken());
        int z = Integer.parseInt(st.nextToken());
        
        br.close();
    }
}
```

---

### Reading Arrays (Fast Way)

```java
import java.io.*;

public class FastArrayInput {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        // Read array size
        int n = Integer.parseInt(br.readLine());
        
        // Method 1: Split and parse
        String[] tokens = br.readLine().split(" ");
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = Integer.parseInt(tokens[i]);
        }
        
        // Method 2: StringTokenizer (faster!)
        StringTokenizer st = new StringTokenizer(br.readLine());
        int[] arr2 = new int[n];
        for (int i = 0; i < n; i++) {
            arr2[i] = Integer.parseInt(st.nextToken());
        }
        
        br.close();
    }
}
```

---

### Reading 2D Arrays

```java
import java.io.*;
import java.util.*;

public class Read2DArray {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        String[] firstLine = br.readLine().split(" ");
        int rows = Integer.parseInt(firstLine[0]);
        int cols = Integer.parseInt(firstLine[1]);
        
        int[][] matrix = new int[rows][cols];
        
        for (int i = 0; i < rows; i++) {
            StringTokenizer st = new StringTokenizer(br.readLine());
            for (int j = 0; j < cols; j++) {
                matrix[i][j] = Integer.parseInt(st.nextToken());
            }
        }
        
        br.close();
    }
}
```

---

## Method 3: Fast Output

### Why StringBuilder?

```java
// ❌ SLOW - Each println is a separate I/O operation
for (int i = 0; i < 100000; i++) {
    System.out.println(i);  // 100,000 I/O calls!
}

// ✅ FAST - Single I/O operation
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 100000; i++) {
    sb.append(i).append("\n");
}
System.out.print(sb);  // 1 I/O call!
```

---

### BufferedWriter for Maximum Speed

```java
import java.io.*;

public class FastOutput {
    public static void main(String[] args) throws IOException {
        BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(System.out));
        
        // Write single values
        bw.write("Hello World\n");
        bw.write(String.valueOf(42) + "\n");
        
        // Write array
        int[] arr = {1, 2, 3, 4, 5};
        for (int num : arr) {
            bw.write(num + " ");
        }
        bw.write("\n");
        
        // CRITICAL: Must flush before closing!
        bw.flush();
        bw.close();
    }
}
```

---

### StringBuilder + BufferedWriter (Fastest!)

```java
import java.io.*;

public class UltraFastOutput {
    public static void main(String[] args) throws IOException {
        BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(System.out));
        StringBuilder sb = new StringBuilder();
        
        // Build output in StringBuilder
        for (int i = 1; i <= 100000; i++) {
            sb.append(i).append(" ");
            
            // Optional: Flush in chunks to avoid memory issues
            if (i % 10000 == 0) {
                bw.write(sb.toString());
                sb.setLength(0);  // Clear StringBuilder
            }
        }
        
        // Write remaining
        bw.write(sb.toString());
        
        bw.flush();
        bw.close();
    }
}
```

---

## Complete Templates for Common Scenarios

### Template 1: Single Test Case

```java
import java.io.*;
import java.util.*;

public class SingleTestCase {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        // Read input
        int n = Integer.parseInt(br.readLine());
        StringTokenizer st = new StringTokenizer(br.readLine());
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = Integer.parseInt(st.nextToken());
        }
        
        // Solve
        int result = solve(arr);
        
        // Output
        System.out.println(result);
        
        br.close();
    }
    
    static int solve(int[] arr) {
        // Your solution here
        return 0;
    }
}
```

---

### Template 2: Multiple Test Cases

```java
import java.io.*;
import java.util.*;

public class MultipleTestCases {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        StringBuilder sb = new StringBuilder();
        
        int t = Integer.parseInt(br.readLine());
        
        while (t-- > 0) {
            // Read input for this test case
            int n = Integer.parseInt(br.readLine());
            StringTokenizer st = new StringTokenizer(br.readLine());
            int[] arr = new int[n];
            for (int i = 0; i < n; i++) {
                arr[i] = Integer.parseInt(st.nextToken());
            }
            
            // Solve
            int result = solve(arr);
            
            // Append to output
            sb.append(result).append("\n");
        }
        
        // Print all at once
        System.out.print(sb);
        
        br.close();
    }
    
    static int solve(int[] arr) {
        // Your solution here
        return 0;
    }
}
```

---

### Template 3: Unknown Number of Test Cases (EOF)

```java
import java.io.*;
import java.util.*;

public class EOFTemplate {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        StringBuilder sb = new StringBuilder();
        
        String line;
        while ((line = br.readLine()) != null) {
            // Process this line
            int n = Integer.parseInt(line);
            
            // Read next line if needed
            String[] tokens = br.readLine().split(" ");
            
            // Solve and append result
            int result = solve(n);
            sb.append(result).append("\n");
        }
        
        System.out.print(sb);
        br.close();
    }
    
    static int solve(int n) {
        return 0;
    }
}
```

---

### Template 4: Grid/Matrix Input

```java
import java.io.*;
import java.util.*;

public class GridTemplate {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        StringTokenizer st = new StringTokenizer(br.readLine());
        int rows = Integer.parseInt(st.nextToken());
        int cols = Integer.parseInt(st.nextToken());
        
        char[][] grid = new char[rows][cols];
        
        for (int i = 0; i < rows; i++) {
            String line = br.readLine();
            for (int j = 0; j < cols; j++) {
                grid[i][j] = line.charAt(j);
            }
        }
        
        // Process grid
        
        br.close();
    }
}
```

---

## Real Problem Example

**Problem:** Read N integers, sort them, and print

### Using Scanner (Slow)

```java
import java.util.*;

public class SortSlow {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        
        int n = sc.nextInt();
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = sc.nextInt();
        }
        
        Arrays.sort(arr);
        
        for (int num : arr) {
            System.out.print(num + " ");
        }
        
        sc.close();
    }
}
```

---

### Using BufferedReader (Fast)

```java
import java.io.*;
import java.util.*;

public class SortFast {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        int n = Integer.parseInt(br.readLine());
        StringTokenizer st = new StringTokenizer(br.readLine());
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = Integer.parseInt(st.nextToken());
        }
        
        Arrays.sort(arr);
        
        StringBuilder sb = new StringBuilder();
        for (int num : arr) {
            sb.append(num).append(" ");
        }
        System.out.println(sb);
        
        br.close();
    }
}
```

---

## Performance Comparison

| Method | Input Time (1M integers) | Output Time (1M integers) |
|--------|-------------------------|---------------------------|
| Scanner | ~8-10 seconds | - |
| BufferedReader + split() | ~2-3 seconds | - |
| BufferedReader + StringTokenizer | **~0.5-1 seconds** ✅ | - |
| System.out.println() | - | ~5-8 seconds |
| StringBuilder + print() | - | **~0.3-0.5 seconds** ✅ |
| BufferedWriter | - | **~0.2-0.4 seconds** ✅ |

---

## Common Input Patterns

### Pattern 1: T Test Cases, Each Has N Elements

```
Input:
3          // T test cases
5          // N for test 1
1 2 3 4 5  // Elements
3          // N for test 2
10 20 30   // Elements
4          // N for test 3
5 6 7 8    // Elements
```

```java
import java.io.*;
import java.util.*;

public class Pattern1 {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        StringBuilder sb = new StringBuilder();
        
        int t = Integer.parseInt(br.readLine());
        
        while (t-- > 0) {
            int n = Integer.parseInt(br.readLine());
            StringTokenizer st = new StringTokenizer(br.readLine());
            int[] arr = new int[n];
            for (int i = 0; i < n; i++) {
                arr[i] = Integer.parseInt(st.nextToken());
            }
            
            // Solve
            int result = Arrays.stream(arr).sum();
            sb.append(result).append("\n");
        }
        
        System.out.print(sb);
        br.close();
    }
}
```

---

### Pattern 2: First Line N M, Then N Lines of M Elements

```
Input:
3 4        // N rows, M columns
1 2 3 4
5 6 7 8
9 10 11 12
```

```java
import java.io.*;
import java.util.*;

public class Pattern2 {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        StringTokenizer st = new StringTokenizer(br.readLine());
        int n = Integer.parseInt(st.nextToken());
        int m = Integer.parseInt(st.nextToken());
        
        int[][] matrix = new int[n][m];
        
        for (int i = 0; i < n; i++) {
            st = new StringTokenizer(br.readLine());
            for (int j = 0; j < m; j++) {
                matrix[i][j] = Integer.parseInt(st.nextToken());
            }
        }
        
        br.close();
    }
}
```

---

### Pattern 3: Pairs on Each Line

```
Input:
4          // Number of pairs
1 10
2 20
3 30
4 40
```

```java
import java.io.*;
import java.util.*;

public class Pattern3 {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        
        int n = Integer.parseInt(br.readLine());
        int[][] pairs = new int[n][2];
        
        for (int i = 0; i < n; i++) {
            StringTokenizer st = new StringTokenizer(br.readLine());
            pairs[i][0] = Integer.parseInt(st.nextToken());
            pairs[i][1] = Integer.parseInt(st.nextToken());
        }
        
        br.close();
    }
}
```

---

## Common Pitfalls & Solutions

### Issue 1: Mixing nextInt() and nextLine()

```java
// ❌ WRONG
Scanner sc = new Scanner(System.in);
int n = sc.nextInt();
String line = sc.nextLine();  // Gets empty string!

// ✅ CORRECT
Scanner sc = new Scanner(System.in);
int n = sc.nextInt();
sc.nextLine();  // Consume newline
String line = sc.nextLine();  // Now gets the line
```

---

### Issue 2: Not Flushing BufferedWriter

```java
// ❌ WRONG - Output might not appear!
BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(System.out));
bw.write("Hello");
bw.close();

// ✅ CORRECT
BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(System.out));
bw.write("Hello");
bw.flush();  // Force output!
bw.close();
```

---

### Issue 3: StringTokenizer on Empty Line

```java
// ❌ Can throw exception on empty line
StringTokenizer st = new StringTokenizer("");
st.nextToken();  // NoSuchElementException!

// ✅ Check first
String line = br.readLine();
if (line.trim().isEmpty()) {
    // Handle empty line
} else {
    StringTokenizer st = new StringTokenizer(line);
}
```

---

### Issue 4: Memory with Large StringBuilder

```java
// ❌ Can cause OutOfMemoryError for huge outputs
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 10000000; i++) {
    sb.append(i).append("\n");
}
System.out.print(sb);

// ✅ Flush periodically
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 10000000; i++) {
    sb.append(i).append("\n");
    if (i % 100000 == 0) {
        System.out.print(sb);
        sb.setLength(0);  // Clear
    }
}
System.out.print(sb);
```

---

## Ultimate Fast I/O Template

```java
import java.io.*;
import java.util.*;

public class FastIO {
    static class Reader {
        BufferedReader br;
        StringTokenizer st;
        
        public Reader() {
            br = new BufferedReader(new InputStreamReader(System.in));
        }
        
        String next() {
            while (st == null || !st.hasMoreElements()) {
                try {
                    st = new StringTokenizer(br.readLine());
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
            return st.nextToken();
        }
        
        int nextInt() {
            return Integer.parseInt(next());
        }
        
        long nextLong() {
            return Long.parseLong(next());
        }
        
        double nextDouble() {
            return Double.parseDouble(next());
        }
        
        String nextLine() {
            String str = "";
            try {
                str = br.readLine();
            } catch (IOException e) {
                e.printStackTrace();
            }
            return str;
        }
    }
    
    public static void main(String[] args) {
        Reader sc = new Reader();
        
        int n = sc.nextInt();
        int[] arr = new int[n];
        for (int i = 0; i < n; i++) {
            arr[i] = sc.nextInt();
        }
        
        // Now use like Scanner but with BufferedReader speed!
        
        // Output
        StringBuilder sb = new StringBuilder();
        for (int num : arr) {
            sb.append(num).append(" ");
        }
        System.out.println(sb);
    }
}
```

---

## Decision Guide

**Choose Scanner when:**
- Input size < 10,000
- Practicing/learning
- Problem is not time-critical

**Choose BufferedReader when:**
- Input size > 100,000
- Multiple test cases
- Time limit is tight
- Competitive programming contest

**Choose StringBuilder when:**
- Printing > 1,000 values
- Multiple test cases
- Building complex output

**Choose BufferedWriter when:**
- Absolutely maximum speed needed
- Printing > 100,000 values

---

## Practice Problems

Try these on LeetCode/HackerRank:
1. Read array and find sum
2. Read matrix and find max element
3. Process multiple test cases
4. Read until EOF
5. Large input sorting (n = 1,000,000)

---

## Key Takeaways

✅ **Scanner is slow** - Only for small inputs
✅ **BufferedReader + StringTokenizer** - Fast input
✅ **StringBuilder** - Fast output building
✅ **BufferedWriter** - Maximum speed output
✅ **Always flush** - BufferedWriter needs flush()
✅ **Build then print** - Don't print in loop
✅ **Use templates** - Copy-paste for contests

---

[← Back: Java Refresher](./Java-Refresher.md) | [Next: Collections Framework →](./Collections-Framework.md)
