# Networking Essentials for Engineers
> Interview-ready practical networking concepts every engineer must know

## 📡 1. TCP vs UDP - When to Use What

### TCP (Transmission Control Protocol)
**Use When:** You need reliability, order, and data integrity
```
✅ Web browsing (HTTP/HTTPS)
✅ Email (SMTP, IMAP)
✅ File transfers (FTP, SFTP)
✅ Database connections
```

**Key Features:**
- Connection-oriented (3-way handshake)
- Guarantees delivery and order
- Flow control and congestion control
- Slower but reliable

### UDP (User Datagram Protocol)
**Use When:** Speed > Reliability
```
✅ Video streaming (YouTube, Netflix)
✅ Online gaming
✅ VoIP calls (Zoom, Skype)
✅ DNS lookups
✅ Live broadcasts
```

**Key Features:**
- Connectionless
- No delivery guarantee
- No ordering
- Faster, lower latency

**Interview Tip:** "I'd use TCP for a banking app where data integrity is critical, but UDP for a live sports streaming app where occasional packet loss is acceptable."

---

## 🌐 2. HTTP/HTTPS - The Web's Foundation

### HTTP Methods (REST APIs)
```
GET     - Retrieve data (idempotent, cacheable)
POST    - Create new resource (not idempotent)
PUT     - Update/replace entire resource (idempotent)
PATCH   - Partial update (not always idempotent)
DELETE  - Remove resource (idempotent)
HEAD    - Get headers only (no body)
OPTIONS - Get allowed methods
```

### Status Codes You Must Know
```
2xx - Success
  200 OK
  201 Created
  204 No Content

3xx - Redirection
  301 Moved Permanently
  302 Found (Temporary redirect)
  304 Not Modified (use cache)

4xx - Client Errors
  400 Bad Request
  401 Unauthorized (not authenticated)
  403 Forbidden (authenticated but not authorized)
  404 Not Found
  429 Too Many Requests (rate limiting)

5xx - Server Errors
  500 Internal Server Error
  502 Bad Gateway (reverse proxy got invalid response)
  503 Service Unavailable
  504 Gateway Timeout
```

### HTTPS = HTTP + TLS/SSL
**Why it matters:**
- Encrypts data in transit
- Prevents man-in-the-middle attacks
- Required for modern browsers (Chrome marks HTTP as "Not Secure")
- SEO boost

**How it works:**
1. Client initiates SSL/TLS handshake
2. Server sends SSL certificate
3. Client validates certificate
4. Symmetric key exchange
5. Encrypted communication begins

---

## 🔍 3. DNS - The Internet's Phonebook

### How DNS Resolution Works
```
1. Browser cache
2. OS cache
3. Router cache
4. ISP DNS server (Recursive resolver)
5. Root nameserver
6. TLD nameserver (.com, .org)
7. Authoritative nameserver (actual answer)
```

### DNS Record Types
```
A     - Maps domain to IPv4 (example.com → 192.168.1.1)
AAAA  - Maps domain to IPv6
CNAME - Alias (www.example.com → example.com)
MX    - Mail server
TXT   - Text records (SPF, DKIM, verification)
NS    - Nameserver records
SOA   - Start of Authority (primary DNS info)
```

### Real-World Scenarios

**Q: Why is DNS a common bottleneck?**
```
❌ Problem: Every request needs DNS lookup
✅ Solution: 
   - DNS caching (browser, OS, application level)
   - Use CDN with built-in DNS optimization
   - Consider DNS prefetching: <link rel="dns-prefetch" href="//api.example.com">
```

**Q: What's DNS propagation?**
- Time for DNS changes to spread globally (TTL dependent)
- Can take 24-48 hours
- Lower TTL before migration (e.g., 300s instead of 86400s)

---

## ⚖️ 4. Load Balancing - Distributing Traffic

### Load Balancing Algorithms

**1. Round Robin**
```
Request 1 → Server A
Request 2 → Server B
Request 3 → Server C
Request 4 → Server A (cycle repeats)

✅ Simple, fair distribution
❌ Doesn't consider server capacity or load
```

**2. Least Connections**
```
Routes to server with fewest active connections

✅ Better for long-lived connections
✅ Adapts to server load
```

**3. Weighted Round Robin**
```
Server A (capacity: 8GB) → Weight 2
Server B (capacity: 4GB) → Weight 1

2 requests to A, 1 to B, repeat

✅ Accounts for different server capacities
```

**4. IP Hash**
```
hash(client_ip) % server_count = server_index

Same client → Same server (session affinity)

✅ Maintains sessions
❌ Uneven distribution if clients are concentrated
```

**5. Least Response Time**
```
Routes to server with lowest latency + fewest connections

✅ Best performance
❌ More complex to implement
```

### Layer 4 vs Layer 7 Load Balancing

**Layer 4 (Transport Layer) - TCP/UDP**
- Faster (only looks at IP + port)
- No content awareness
- Example: AWS NLB

**Layer 7 (Application Layer) - HTTP**
- Content-based routing (URL, headers, cookies)
- SSL termination
- More CPU intensive
- Example: AWS ALB, nginx

**Interview Question:** "Design Instagram's load balancing"
```
Layer 7 LB → Routes /api/users to User Service
          → Routes /api/posts to Post Service
          → Routes /media/* to Media Service with CDN
```

---

## 🔐 5. Common Ports - Must Know

```
20/21   - FTP (File Transfer)
22      - SSH (Secure Shell)
23      - Telnet (avoid - insecure)
25      - SMTP (Email sending)
53      - DNS
80      - HTTP
443     - HTTPS
3306    - MySQL
5432    - PostgreSQL
6379    - Redis
27017   - MongoDB
8080    - Alternative HTTP (development)
9200    - Elasticsearch
```

**Security Tip:** Only expose necessary ports. Use security groups/firewalls to restrict access.

---

## 🌍 6. IP Addressing & Subnetting (Practical)

### IPv4 Basics
```
Format: xxx.xxx.xxx.xxx (each octet 0-255)
Example: 192.168.1.100

Private IP Ranges (not routable on internet):
  10.0.0.0      - 10.255.255.255    (Class A)
  172.16.0.0    - 172.31.255.255    (Class B)
  192.168.0.0   - 192.168.255.255   (Class C)
```

### CIDR Notation (Quick Reference)
```
/32 - 1 IP        (255.255.255.255)
/24 - 256 IPs     (255.255.255.0)   - Common for small networks
/16 - 65,536 IPs  (255.255.0.0)     - Common for VPCs
/8  - 16M IPs     (255.0.0.0)
```

### Real-World Example: AWS VPC
```
VPC: 10.0.0.0/16 (65,536 IPs)
  
  Subnet 1 (Public):  10.0.1.0/24 (256 IPs) - Web servers
  Subnet 2 (Private): 10.0.2.0/24 (256 IPs) - App servers
  Subnet 3 (Private): 10.0.3.0/24 (256 IPs) - Databases
```

---

## 🚀 7. CDN (Content Delivery Network)

### How CDN Works
```
User in Tokyo requests image from example.com (server in USA)

Without CDN: 150ms latency
With CDN:    15ms latency (served from Tokyo edge server)
```

### What to Cache on CDN
```
✅ Static assets (images, CSS, JS, fonts)
✅ Videos
✅ Static API responses
❌ User-specific data
❌ Frequently changing data
❌ Sensitive information
```

### Popular CDNs
- Cloudflare
- AWS CloudFront
- Akamai
- Fastly

**Interview Scenario:** "How would you optimize image loading for a global news website?"
```
1. Use CDN to serve images from edge locations
2. Implement lazy loading
3. Use responsive images (srcset)
4. WebP format with fallbacks
5. Image compression
6. Cache-Control headers (max-age=31536000 for immutable assets)
```

---

## 🔄 8. Proxy vs Reverse Proxy

### Forward Proxy (Proxy)
```
Client → Proxy → Internet

Use cases:
✅ Hide client IP (anonymity)
✅ Bypass geo-restrictions
✅ Content filtering (corporate networks)
✅ Caching

Example: VPN, corporate proxy
```

### Reverse Proxy
```
Client → Reverse Proxy → Backend Servers

Use cases:
✅ Load balancing
✅ SSL termination
✅ Caching
✅ Security (hide backend architecture)
✅ Compression

Example: nginx, HAProxy, AWS ALB
```

**Key Difference:** 
- Forward proxy represents clients
- Reverse proxy represents servers

---

## 🔌 9. WebSockets vs HTTP

### HTTP (Request-Response)
```
Client: "Give me data" → Server: "Here it is"
(Connection closes)

✅ Simple, stateless
❌ Inefficient for real-time updates (need polling)
```

### WebSockets (Persistent Connection)
```
Client ↔ Server (bidirectional, persistent)

✅ Real-time, low latency
✅ Less overhead (no repeated handshakes)
❌ More complex to scale
❌ Stateful (connection state management)
```

### When to Use WebSockets
```
✅ Chat applications (Slack, WhatsApp)
✅ Live sports scores
✅ Collaborative editing (Google Docs)
✅ Gaming
✅ Stock tickers

❌ Regular CRUD operations (use REST)
❌ Simple data fetching
```

### Modern Alternative: Server-Sent Events (SSE)
```
One-way server → client updates
Simpler than WebSockets
Auto-reconnect
Use for: Live feeds, notifications
```

---

## 🛡️ 10. SSL/TLS - Encryption Basics

### Symmetric vs Asymmetric Encryption

**Symmetric (Same Key)**
```
Encryption key = Decryption key
Fast
Example: AES

Problem: How to share the key securely?
```

**Asymmetric (Public/Private Key Pair)**
```
Public key (encrypt) ≠ Private key (decrypt)
Slower
Example: RSA

Used for: Initial handshake, key exchange
```

### How HTTPS Handshake Works (Simplified)
```
1. Client: "Hello, I support these cipher suites"
2. Server: "Hello, let's use AES-256. Here's my certificate"
3. Client: Validates certificate, generates session key
4. Client: Encrypts session key with server's public key, sends it
5. Server: Decrypts with private key
6. Both now have the same session key (symmetric encryption begins)
```

**Interview Tip:** "HTTPS uses asymmetric encryption for the handshake (slow but secure), then switches to symmetric encryption for data transfer (fast)."

---

## 📊 11. Network Troubleshooting Commands

### Essential Commands

**1. ping - Test connectivity**
```bash
ping google.com

Checks: Is host reachable?
```

**2. traceroute (tracert on Windows) - Path to destination**
```bash
traceroute google.com

Shows: Each hop from source to destination
Identifies: Where delays/failures occur
```

**3. nslookup/dig - DNS lookup**
```bash
nslookup example.com
dig example.com

Shows: IP address, DNS records
```

**4. netstat - Network connections**
```bash
netstat -an   # All connections
netstat -tulpn # Linux: listening ports

Shows: Active connections, listening ports
```

**5. curl - Test HTTP endpoints**
```bash
curl -I https://example.com  # Headers only
curl -v https://api.example.com/users  # Verbose

Test: APIs, connectivity, response times
```

**6. telnet - Test port connectivity**
```bash
telnet example.com 80

Checks: Is port open and accepting connections?
```

### Real Interview Scenario
**Q: "User in Europe can't access our app, but US users can. How do you troubleshoot?"**
```
1. ping api.example.com (from Europe)
   → Check basic connectivity
   
2. traceroute api.example.com
   → Identify where packets are dropping
   
3. nslookup api.example.com
   → Verify DNS resolves correctly in Europe
   
4. curl -I https://api.example.com
   → Test HTTP layer
   
5. Check CDN/Load balancer logs
   → Geo-routing issue?
   
6. Check firewall/security group rules
   → IP blocking?
```

---

## 🔥 12. Common Interview Questions

### Q1: "What happens when you type google.com in a browser?"

**Complete Answer:**
```
1. Browser checks cache (recently visited?)
2. DNS lookup:
   - Browser cache → OS cache → Router → ISP
   - Resolves google.com → 142.250.185.46
3. TCP connection:
   - 3-way handshake (SYN, SYN-ACK, ACK)
4. TLS handshake (if HTTPS):
   - Certificate verification
   - Key exchange
5. HTTP request:
   - GET / HTTP/1.1
   - Headers (User-Agent, Accept, etc.)
6. Server processes request:
   - Load balancer routes to server
   - Server generates response
7. HTTP response:
   - Status code, headers, body
8. Browser renders:
   - Parses HTML
   - Loads CSS, JS, images (parallel requests)
   - Renders page
9. Connection closes (or kept alive for HTTP/1.1+)
```

### Q2: "How would you design a rate limiter?"

**Approaches:**
```
1. Token Bucket:
   - Bucket refills at fixed rate
   - Request consumes token
   - Reject if no tokens available
   ✅ Allows bursts

2. Fixed Window:
   - Count requests per time window
   - Reset counter at window boundary
   ✅ Simple
   ❌ Burst at window edges

3. Sliding Window:
   - Weighted count of current + previous window
   ✅ Smoother
   ❌ More memory

4. Leaky Bucket:
   - Requests queued, processed at fixed rate
   ✅ Smooth output
   ❌ Can add latency
```

### Q3: "Difference between connection timeout and read timeout?"

```
Connection Timeout:
- Time to establish TCP connection
- Set: 5-10 seconds
- Failure: Server unreachable, network issue

Read Timeout:
- Time waiting for data after connection established
- Set: 30-60 seconds (depends on operation)
- Failure: Server slow, processing issue

Example (Java):
client.setConnectTimeout(10000);  // 10 seconds
client.setReadTimeout(30000);     // 30 seconds
```

### Q4: "How do cookies vs tokens work?"

**Cookies (Traditional)**
```
1. Server sets: Set-Cookie: sessionId=abc123
2. Browser stores cookie
3. Browser automatically sends with every request
4. Server validates sessionId

✅ Automatic, no code needed
❌ CSRF vulnerable
❌ Doesn't work well with mobile apps
```

**Tokens (JWT - Modern)**
```
1. Login: Server returns token
2. Client stores (localStorage/sessionStorage)
3. Client manually includes: Authorization: Bearer <token>
4. Server validates token signature

✅ Stateless
✅ Works across domains
✅ Mobile-friendly
❌ Can't invalidate easily (until expiry)
```

### Q5: "How would you handle 1 million concurrent WebSocket connections?"

```
Challenges:
- Each connection = memory + CPU
- Single server limit: ~10k-100k connections

Solution:
1. Horizontal scaling:
   - Multiple WebSocket servers
   - Load balancer with sticky sessions
   
2. Connection pooling:
   - Not all 1M are actively sending data
   - Efficient event loop (Node.js, Go)
   
3. Message broker:
   - Redis Pub/Sub or Kafka
   - Servers subscribe to channels
   - Broadcast messages across servers
   
4. Database:
   - Don't keep all state in memory
   - Use Redis for active sessions
   
5. Monitoring:
   - Track connection counts
   - Auto-scale based on load

Architecture:
Client → Load Balancer → WS Server 1 ↘
                      → WS Server 2  → Redis/Kafka → Database
                      → WS Server N ↗
```

---

## 🎯 Quick Interview Prep Checklist

Before your interview, make sure you can explain:

- [ ] TCP 3-way handshake
- [ ] HTTP vs HTTPS (and how TLS works)
- [ ] Common HTTP methods and status codes
- [ ] DNS resolution process
- [ ] Load balancing algorithms
- [ ] Difference between Layer 4 and Layer 7 LB
- [ ] When to use TCP vs UDP
- [ ] How CDN improves performance
- [ ] Forward proxy vs reverse proxy
- [ ] WebSockets vs HTTP
- [ ] IP addressing and CIDR notation
- [ ] How to troubleshoot network issues
- [ ] Rate limiting strategies
- [ ] Session management (cookies vs tokens)

---

## 🧰 Practical Tips

### 1. Design Decisions
Always consider:
- **Latency requirements**: Real-time? → WebSockets, UDP
- **Reliability needs**: Banking? → TCP, HTTPS
- **Scale**: Millions of users? → CDN, load balancing, horizontal scaling
- **Geography**: Global users? → CDN, geo-routing
- **Security**: Sensitive data? → HTTPS, encryption, WAF

### 2. Performance Optimization
```
Network level:
✅ Use CDN for static content
✅ Enable HTTP/2 (multiplexing)
✅ Compress responses (gzip, brotli)
✅ Use connection pooling
✅ Implement caching headers
✅ Reduce DNS lookups
✅ Use HTTP keep-alive
```

### 3. Security Best Practices
```
✅ Always use HTTPS
✅ Implement rate limiting
✅ Validate input
✅ Use WAF (Web Application Firewall)
✅ Keep software updated
✅ Principle of least privilege (firewall rules)
✅ Monitor for anomalies
```

---

## 📚 Resources for Deep Dive

- **Tools**: Wireshark (packet analysis), Postman (API testing)
- **Books**: "Computer Networking: A Top-Down Approach"
- **Practice**: Set up nginx, configure load balancer, capture packets
- **Visualize**: Draw diagrams of your explanations

---

**Remember**: In interviews, it's not just about knowing the answer—it's about demonstrating **how you think** through problems. Always clarify requirements, discuss trade-offs, and explain your reasoning.

Good luck! 🚀
