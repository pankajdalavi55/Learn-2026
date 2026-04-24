# Complete System Design: News Feed / Timeline System (Production-Ready)

> **Complexity Level:** Intermediate to Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** Facebook News Feed, Twitter Timeline, Instagram Feed, LinkedIn Feed

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Fan-out Strategies](#7-deep-dive-fan-out-strategies)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a news feed system like Facebook or Twitter that shows a personalized feed of posts from friends/followed users."

**Interviewer's Perspective:**  
This is a top-tier system design problem used to assess:
- **Fan-out strategies** — push vs pull vs hybrid
- **Ranking algorithms** — chronological vs ML-based relevance
- **Caching at scale** — billions of pre-computed feed entries
- **Real-time updates** — WebSocket/SSE for live content
- **Write amplification** — the celebrity problem and how to mitigate it
- **Data modeling** — social graphs, feed materialization, denormalization

This problem separates senior candidates from the rest because it has no single "correct" answer — the tradeoffs are the answer.

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I dive into the design, I'd like to clarify the scope. Can I ask a few questions?"

**Interviewer:** "Absolutely, go ahead."

**Candidate:** "What types of content can users post? Just text, or also images, videos, links?"

**Interviewer:** "Support text, images, and videos. Media can be handled by a separate media pipeline — focus on the feed logic."

**Candidate:** "For the feed itself, should it be purely chronological (newest first), or ranked by relevance like Facebook?"

**Interviewer:** "Start with chronological, but design it so we can layer on a ranking algorithm later. Assume we'll add ranking in V2."

**Candidate:** "What's the social model? Bidirectional friendships like Facebook, or unidirectional follows like Twitter?"

**Interviewer:** "Unidirectional follows — like Twitter. User A can follow User B without B following back."

**Candidate:** "Should the feed update in real-time, or is it okay if new posts appear on the next refresh?"

**Interviewer:** "Eventual consistency is fine. A 1-2 second delay between posting and feed appearance is acceptable. Real-time push is a nice-to-have."

**Candidate:** "Any engagement features? Likes, comments, shares?"

**Interviewer:** "Yes — likes and comments. Users should see counts on feed items, but the detailed comment thread is a separate page. Don't over-design the comment system."

**Candidate:** "Perfect. Let me summarize what I've gathered."

### 2.1 Functional Requirements

| # | Requirement | Description |
|---|------------|-------------|
| FR-1 | Create Post | Users create posts with text, images, or video |
| FR-2 | View Feed | Users see a personalized feed of posts from followed users |
| FR-3 | Follow/Unfollow | Users can follow/unfollow other users (unidirectional) |
| FR-4 | Like/Comment | Users can like and comment on posts |
| FR-5 | Feed Ordering | Feed sorted chronologically (V1), by relevance (V2) |
| FR-6 | Pagination | Infinite scroll with cursor-based pagination |

### 2.2 Non-Functional Requirements

| # | Requirement | Target |
|---|------------|--------|
| NFR-1 | Feed Latency | < 500ms for feed load (p99) |
| NFR-2 | Scale | 500M daily active users |
| NFR-3 | Consistency | Eventual (1-2 second lag acceptable) |
| NFR-4 | Availability | 99.99% uptime (availability > consistency) |
| NFR-5 | Durability | Posts must never be lost |
| NFR-6 | Throughput | Handle 58K+ feed reads/sec, 580+ post writes/sec |

### 2.3 Scale Parameters

| Parameter | Value |
|-----------|-------|
| Daily Active Users (DAU) | 500 million |
| Average follows per user | 200 |
| Feed checks per user per day | 10 |
| New posts per day | 50 million |
| Average post size (text + metadata) | ~1 KB |
| Media per post (images/video) | Stored separately on object storage |

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me work through the numbers to understand the scale we're dealing with."

```
Feed Reads:
  500M DAU × 10 feed checks/day = 5 billion reads/day
  5B / 86,400 sec = ~58,000 reads/sec
  Peak (3× average) = ~174,000 reads/sec

Post Writes:
  50M new posts/day
  50M / 86,400 sec = ~580 writes/sec
  Peak (5× average) = ~2,900 writes/sec

Read:Write Ratio = 5B / 50M = 100:1 (extremely read-heavy)
```

### 3.2 Fan-out Estimation

```
Average followers per poster:
  Each post fans out to ~200 followers on average

Total fan-out deliveries per day:
  50M posts × 200 followers = 10 billion feed deliveries/day
  10B / 86,400 = ~115,740 fan-out writes/sec

Celebrity impact (worst case):
  1 post from user with 100M followers = 100M fan-out writes
  At 10 KB per entry = 1 TB for a single celebrity post
```

### 3.3 Storage Estimation

```
Post Storage (1 year):
  50M posts/day × 365 days × 1 KB = ~18.25 TB/year (text + metadata only)
  Media (images/video) = 10× text = ~180 TB/year (stored on S3/CDN)

Feed Cache (materialized feeds):
  500M users × 20 cached posts × 200 bytes per entry = ~2 TB
  This fits comfortably in a Redis cluster

Social Graph:
  500M users × 200 follows = 100 billion edges
  Each edge ~50 bytes = ~5 TB
```

### 3.4 Bandwidth Estimation

```
Feed reads:
  58K reads/sec × 20 posts × 1 KB = ~1.16 GB/sec (metadata only)
  With media thumbnails: ~5-10 GB/sec (served from CDN)

Post writes:
  580 writes/sec × 1 KB = ~580 KB/sec (text)
  Media uploads handled separately via presigned S3 URLs
```

**Candidate:** "The key insight here is that this system is overwhelmingly read-heavy — 100:1 ratio. The fan-out problem is the dominant cost: 10 billion deliveries per day. This tells me we need aggressive caching and a smart fan-out strategy."

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                    │
│         Mobile App (iOS/Android)  |  Web App  |  Third-Party APIs       │
└────────────────────────┬────────────────────────────────────────────────┘
                         │
                    ┌────▼─────┐
                    │   CDN    │  (Media, static assets, thumbnails)
                    └────┬─────┘
                         │
                 ┌───────▼─────────┐
                 │   API Gateway   │  Auth, rate limiting, routing
                 │  (Kong / Envoy) │
                 └───────┬─────────┘
                         │
      ┌──────────────────┼──────────────────────┐
      │                  │                      │
 ┌────▼─────┐     ┌──────▼──────┐       ┌──────▼──────┐
 │  Post    │     │   Feed      │       │   Social    │
 │  Service │     │   Service   │       │   Graph     │
 │          │     │             │       │   Service   │
 └────┬─────┘     └──────┬──────┘       └──────┬──────┘
      │                  │                      │
      │           ┌──────▼──────┐               │
      │           │  Ranking    │               │
      │           │  Service    │               │
      │           └──────┬──────┘               │
      │                  │                      │
      ▼                  ▼                      ▼
 ┌─────────────────────────────────────────────────────┐
 │                   MESSAGE BUS (Kafka)                │
 │  [post-events]  [fan-out-tasks]  [feed-updates]     │
 └────────────┬──────────────┬──────────────┬──────────┘
              │              │              │
       ┌──────▼──────┐      │       ┌──────▼──────┐
       │  Fan-out    │      │       │ Notification│
       │  Service    │      │       │  Service    │
       │  (Workers)  │      │       └─────────────┘
       └──────┬──────┘      │
              │              │
              ▼              ▼
 ┌─────────────────────────────────────────────────────┐
 │                    DATA LAYER                        │
 │                                                      │
 │  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
 │  │  Posts   │  │  Social  │  │    Feed Cache     │  │
 │  │  (MySQL) │  │  Graph   │  │    (Redis         │  │
 │  │          │  │  (MySQL) │  │     Cluster)      │  │
 │  └──────────┘  └──────────┘  └───────────────────┘  │
 │                                                      │
 │  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
 │  │  Media   │  │  Feed    │  │   User Activity   │  │
 │  │  (S3 +   │  │  Store   │  │   (Redis /        │  │
 │  │   CDN)   │  │(Cassandra│  │    DynamoDB)      │  │
 │  └──────────┘  └──────────┘  └───────────────────┘  │
 └─────────────────────────────────────────────────────┘
```

### 4.2 API Design

**Candidate:** "Here are the core APIs. I'll use RESTful design with cursor-based pagination."

```
POST   /api/v1/posts                    — Create a new post
GET    /api/v1/posts/{postId}           — Get a single post
DELETE /api/v1/posts/{postId}           — Delete a post

GET    /api/v1/feed?cursor=&size=20     — Get personalized feed
GET    /api/v1/feed/refresh             — Check for new posts since last fetch

POST   /api/v1/users/{userId}/follow    — Follow a user
DELETE /api/v1/users/{userId}/follow    — Unfollow a user
GET    /api/v1/users/{userId}/followers — List followers
GET    /api/v1/users/{userId}/following — List following

POST   /api/v1/posts/{postId}/like      — Like a post
DELETE /api/v1/posts/{postId}/like      — Unlike a post
POST   /api/v1/posts/{postId}/comments  — Add a comment
```

### 4.3 Data Flow: Creating a Post

```
User creates post
       │
       ▼
┌──────────────┐     ┌─────────────┐
│ Post Service │────▶│ Posts DB    │  1. Persist post
└──────┬───────┘     └─────────────┘
       │
       ▼
┌──────────────┐     ┌─────────────┐
│ Media Service│────▶│ S3 + CDN   │  2. Upload media (async)
└──────┬───────┘     └─────────────┘
       │
       ▼
┌──────────────┐     ┌─────────────────────────────┐
│ Kafka Topic  │────▶│ Fan-out Service (Workers)   │  3. Fan-out to followers
│ [post-events]│     └──────────────┬──────────────┘
└──────────────┘                    │
                                    ▼
                             ┌──────────────┐
                             │ Feed Cache   │  4. Write to each follower's cache
                             │ (Redis)      │
                             └──────────────┘
```

### 4.4 Data Flow: Reading the Feed

```
User opens feed
       │
       ▼
┌──────────────┐     ┌──────────────┐
│ Feed Service │────▶│ Feed Cache   │  1. Fetch pre-computed feed from Redis
└──────┬───────┘     │ (Redis)      │
       │             └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│ Ranking      │────▶│ ML Model /   │  2. Re-rank posts (optional, V2)
│ Service      │     │ Feature Store│
└──────┬───────┘     └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│ Post Service │────▶│ Posts DB /   │  3. Hydrate post details (author, media)
│ (Hydration)  │     │ Cache        │
└──────┬───────┘     └──────────────┘
       │
       ▼
   Return feed
   to client
```

---

## 5. Deep Dive: Core Components

### 5.1 Post Service

**Candidate:** "The Post Service handles CRUD operations for posts. It's the write entry point."

**Responsibilities:**
- Create, read, update, delete posts
- Validate content (text length, media formats)
- Generate presigned URLs for media uploads
- Publish post events to Kafka for downstream processing

```javascript
// Post Service — Express.js handler
class PostService {
  async createPost(userId, content, mediaIds) {
    const post = {
      id: generateSnowflakeId(),
      authorId: userId,
      content: content,
      mediaUrls: mediaIds.map(id => this.mediaService.getUrl(id)),
      createdAt: Date.now(),
      likeCount: 0,
      commentCount: 0,
    };

    // Persist to primary database
    await this.postRepository.save(post);

    // Invalidate author's profile cache
    await this.cache.del(`user:${userId}:posts`);

    // Publish event for fan-out
    await this.kafka.publish('post-events', {
      type: 'POST_CREATED',
      postId: post.id,
      authorId: userId,
      createdAt: post.createdAt,
    });

    return post;
  }

  async getPost(postId) {
    // Check cache first
    const cached = await this.cache.get(`post:${postId}`);
    if (cached) return JSON.parse(cached);

    const post = await this.postRepository.findById(postId);
    if (post) {
      await this.cache.setex(`post:${postId}`, 3600, JSON.stringify(post));
    }
    return post;
  }
}
```

### 5.2 Social Graph Service

**Candidate:** "The Social Graph Service manages follow relationships. This is critical because it determines who sees whose posts."

**Data Structure Decision:**

| Approach | Pros | Cons | Use When |
|----------|------|------|----------|
| Adjacency List (DB rows) | Simple, flexible queries | Slow for graph traversals | Follow counts < 10K |
| Adjacency Matrix | O(1) lookup | Space: O(n²), impractical | Tiny user base |
| Graph Database (Neo4j) | Native graph queries | Operational complexity | Complex social features |
| **Redis Sets** | Fast lookups, intersection | Memory cost | Hot path lookups |

**Candidate:** "I'll use MySQL for the source of truth and Redis sets for hot-path lookups."

```javascript
class SocialGraphService {
  async follow(followerId, followeeId) {
    // Write to MySQL (source of truth)
    await this.db.query(
      'INSERT IGNORE INTO user_follows (follower_id, followee_id, created_at) VALUES (?, ?, NOW())',
      [followerId, followeeId]
    );

    // Update Redis for fast lookups
    await this.redis.sadd(`following:${followerId}`, followeeId);
    await this.redis.sadd(`followers:${followeeId}`, followerId);

    // Update follower counts
    await this.redis.incr(`user:${followerId}:following_count`);
    await this.redis.incr(`user:${followeeId}:follower_count`);
  }

  async unfollow(followerId, followeeId) {
    await this.db.query(
      'DELETE FROM user_follows WHERE follower_id = ? AND followee_id = ?',
      [followerId, followeeId]
    );

    await this.redis.srem(`following:${followerId}`, followeeId);
    await this.redis.srem(`followers:${followeeId}`, followerId);

    await this.redis.decr(`user:${followerId}:following_count`);
    await this.redis.decr(`user:${followeeId}:follower_count`);
  }

  async getFollowers(userId, cursor, limit = 100) {
    // For fan-out: need all followers — paginate from DB
    return this.db.query(
      `SELECT follower_id FROM user_follows
       WHERE followee_id = ? AND follower_id > ?
       ORDER BY follower_id LIMIT ?`,
      [userId, cursor || 0, limit]
    );
  }

  async getFollowerCount(userId) {
    return this.redis.get(`user:${userId}:follower_count`);
  }
}
```

### 5.3 Fan-out Service

**Candidate:** "This is the heart of the system. When a post is created, the fan-out service delivers it to all followers' feeds. I'll explain the three strategies in detail in the next section."

### 5.4 Feed Service

**Candidate:** "The Feed Service is what the client calls to get the user's feed. It aggregates from the cache, optionally merges in celebrity posts, and hydrates the response."

```javascript
class FeedService {
  async getFeed(userId, cursor, size = 20) {
    // Step 1: Get pre-computed feed entries from Redis
    const feedKey = `feed:${userId}`;
    const start = cursor ? await this.findCursorIndex(feedKey, cursor) : 0;
    const entries = await this.redis.zrevrange(
      feedKey, start, start + size - 1, 'WITHSCORES'
    );

    if (entries.length === 0) {
      // Cache miss — fall back to pull-based feed generation
      return this.generateFeedOnRead(userId, size);
    }

    // Step 2: Hydrate post details
    const postIds = entries.map(e => e.member);
    const posts = await this.postService.getPostsBatch(postIds);

    // Step 3: Merge celebrity posts (hybrid approach)
    const celebrityPosts = await this.getCelebrityPosts(userId);
    const merged = this.mergeSorted(posts, celebrityPosts);

    // Step 4: Build response with next cursor
    const nextCursor = merged.length > 0
      ? merged[merged.length - 1].id
      : null;

    return {
      posts: merged.slice(0, size),
      nextCursor,
      hasMore: merged.length > size,
    };
  }

  async generateFeedOnRead(userId, size) {
    // Pull model fallback: fetch latest posts from all followed users
    const following = await this.socialGraph.getFollowing(userId);

    const postLists = await Promise.all(
      following.map(uid =>
        this.postService.getRecentPosts(uid, 5)
      )
    );

    // Merge-sort all post lists by timestamp
    return this.mergeKSorted(postLists, size);
  }
}
```

---

## 6. Deep Dive: Database Design

### 6.1 Schema Design

**Candidate:** "Let me walk through the key tables."

#### Posts Table (MySQL / PostgreSQL)

```sql
CREATE TABLE posts (
    id              BIGINT PRIMARY KEY,          -- Snowflake ID (time-sortable)
    author_id       BIGINT NOT NULL,
    content         TEXT,                         -- Text content (max 5000 chars)
    media_urls      JSON,                         -- Array of media URLs
    post_type       ENUM('text','image','video','link') DEFAULT 'text',
    like_count      INT DEFAULT 0,
    comment_count   INT DEFAULT 0,
    share_count     INT DEFAULT 0,
    is_deleted      BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_author_created (author_id, created_at DESC),
    INDEX idx_created (created_at DESC)
) ENGINE=InnoDB;
```

#### User Follows Table (MySQL / PostgreSQL)

```sql
CREATE TABLE user_follows (
    follower_id     BIGINT NOT NULL,
    followee_id     BIGINT NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (follower_id, followee_id),
    INDEX idx_followee (followee_id, follower_id),
    INDEX idx_followee_created (followee_id, created_at)
) ENGINE=InnoDB;
```

#### Feed Items Table (Cassandra — Materialized Feed)

```sql
-- Cassandra CQL: Wide-column store, partitioned by user_id
-- Each user's feed is a single partition with posts sorted by time
CREATE TABLE feed_items (
    user_id     BIGINT,
    post_id     BIGINT,
    author_id   BIGINT,
    score       DOUBLE,           -- Ranking score (timestamp for V1)
    created_at  TIMESTAMP,

    PRIMARY KEY (user_id, created_at, post_id)
) WITH CLUSTERING ORDER BY (created_at DESC, post_id DESC)
  AND default_time_to_live = 2592000;  -- 30-day TTL
```

#### Likes Table

```sql
CREATE TABLE post_likes (
    post_id     BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (post_id, user_id),
    INDEX idx_user_likes (user_id, created_at DESC)
) ENGINE=InnoDB;
```

### 6.2 Database Selection Rationale

```
┌───────────────┬──────────────────┬─────────────────────────────────────┐
│ Data Type     │ Database         │ Reasoning                           │
├───────────────┼──────────────────┼─────────────────────────────────────┤
│ Posts         │ MySQL (sharded)  │ ACID, relational queries,           │
│               │                  │ mature tooling                      │
├───────────────┼──────────────────┼─────────────────────────────────────┤
│ Social Graph  │ MySQL + Redis    │ MySQL = source of truth,            │
│               │                  │ Redis = fast lookups                │
├───────────────┼──────────────────┼─────────────────────────────────────┤
│ Feed Store    │ Cassandra        │ Wide-column, partition by user_id,  │
│ (materialized)│                  │ time-sorted, handles massive writes │
├───────────────┼──────────────────┼─────────────────────────────────────┤
│ Feed Cache    │ Redis Cluster    │ In-memory, sorted sets (ZSET),      │
│ (hot feeds)   │                  │ sub-ms reads, LRU eviction          │
├───────────────┼──────────────────┼─────────────────────────────────────┤
│ Media         │ S3 + CloudFront  │ Object storage + CDN for            │
│               │                  │ low-latency media delivery          │
├───────────────┼──────────────────┼─────────────────────────────────────┤
│ Counters      │ Redis            │ Atomic INCR/DECR for like/comment   │
│ (likes, etc.) │                  │ counts, async flush to MySQL        │
└───────────────┴──────────────────┴─────────────────────────────────────┘
```

### 6.3 Sharding Strategy

**Candidate:** "For posts, I'll shard by `author_id` so all posts from a single user live on one shard — this keeps 'get all posts by user X' efficient."

```
Sharding Plan:
  Posts DB:     Shard by author_id (consistent hashing, 256 virtual shards)
  Social Graph: Shard by follower_id (fan-out reads are by followee, so add GSI)
  Feed Store:   Partition by user_id (Cassandra handles this natively)
  Feed Cache:   Shard by user_id across Redis Cluster nodes
```

---

## 7. Deep Dive: Fan-out Strategies

**Candidate:** "This is the most important design decision for a news feed system. There are three approaches, and the right choice depends on the user distribution."

### 7.1 Fan-out on Write (Push Model)

**Concept:** When a user publishes a post, immediately push it into every follower's feed.

```
User A posts
     │
     ▼
┌──────────┐    ┌──────────────┐    ┌──────────────────────────┐
│ Post     │───▶│ Kafka Topic  │───▶│ Fan-out Workers (100+)   │
│ Service  │    │ [post-events]│    │                          │
└──────────┘    └──────────────┘    │ For each follower of A:  │
                                    │   → Write to feed cache  │
                                    │   → Write to feed store  │
                                    └────────────┬─────────────┘
                                                 │
                           ┌─────────────────────┼─────────────────────┐
                           │                     │                     │
                    ┌──────▼──────┐       ┌──────▼──────┐      ┌──────▼──────┐
                    │ Follower B  │       │ Follower C  │      │ Follower D  │
                    │ Feed Cache  │       │ Feed Cache  │      │ Feed Cache  │
                    └─────────────┘       └─────────────┘      └─────────────┘
```

**Implementation:**

```python
# Fan-out Worker (Python) — processes post events from Kafka
class FanoutWorker:
    def __init__(self, redis_client, social_graph_client, cassandra_session):
        self.redis = redis_client
        self.social_graph = social_graph_client
        self.cassandra = cassandra_session

    def process_post_event(self, event):
        post_id = event['post_id']
        author_id = event['author_id']
        created_at = event['created_at']

        # Paginate through all followers
        cursor = None
        while True:
            followers, cursor = self.social_graph.get_followers(
                author_id, cursor=cursor, limit=1000
            )

            if not followers:
                break

            # Batch write to Redis feed caches
            pipeline = self.redis.pipeline()
            for follower_id in followers:
                feed_key = f"feed:{follower_id}"
                pipeline.zadd(feed_key, {post_id: created_at})
                # Keep only the latest 800 posts per feed
                pipeline.zremrangebyrank(feed_key, 0, -801)
            pipeline.execute()

            # Batch write to Cassandra (durable feed store)
            batch = BatchStatement()
            for follower_id in followers:
                batch.add(self.insert_feed_stmt, (
                    follower_id, post_id, author_id, created_at
                ))
            self.cassandra.execute(batch)

            if cursor is None:
                break
```

**Pros and Cons:**

| Pros | Cons |
|------|------|
| Feed reads are instant (pre-computed) | **Celebrity problem**: user with 100M followers = 100M writes per post |
| Simple read path | High write amplification (10B writes/day) |
| Predictable read latency | Wasted work for inactive users (most never open the app) |
| Easy to implement | Fan-out lag during bursts |

### 7.2 Fan-out on Read (Pull Model)

**Concept:** Don't pre-compute feeds. When a user opens their feed, fetch the latest posts from all users they follow and merge them in real-time.

```
User B opens feed
     │
     ▼
┌──────────┐    ┌──────────────────────────────────────────────┐
│ Feed     │───▶│ For each user B follows (A, C, D, E...):     │
│ Service  │    │   → Fetch latest N posts from each           │
│          │    │   → Merge-sort by timestamp                  │
│          │    │   → Return top 20                            │
└──────────┘    └──────────────────────────────────────────────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
   ┌──────▼──────┐ ┌──────▼──────┐ ┌───────▼─────┐
   │ User A's    │ │ User C's    │ │ User D's    │
   │ posts cache │ │ posts cache │ │ posts cache │
   └─────────────┘ └─────────────┘ └─────────────┘
```

**Implementation:**

```javascript
// Pull-based feed generation
class PullBasedFeedService {
  async getFeed(userId, size = 20) {
    // Get all users this person follows
    const following = await this.socialGraph.getFollowing(userId);

    // Fetch recent posts from each followed user (parallel)
    const postLists = await Promise.all(
      following.map(followeeId =>
        this.postCache.getRecentPosts(followeeId, 50) // top 50 per user
      )
    );

    // K-way merge sort by timestamp (descending)
    const merged = this.kWayMerge(postLists, size);

    return merged;
  }

  kWayMerge(lists, k) {
    // Min-heap based merge of K sorted lists
    const heap = new MinHeap((a, b) => b.createdAt - a.createdAt);

    // Seed heap with first element of each list
    for (let i = 0; i < lists.length; i++) {
      if (lists[i].length > 0) {
        heap.push({ post: lists[i][0], listIdx: i, postIdx: 0 });
      }
    }

    const result = [];
    while (result.length < k && heap.size() > 0) {
      const { post, listIdx, postIdx } = heap.pop();
      result.push(post);

      if (postIdx + 1 < lists[listIdx].length) {
        heap.push({
          post: lists[listIdx][postIdx + 1],
          listIdx,
          postIdx: postIdx + 1,
        });
      }
    }

    return result;
  }
}
```

**Pros and Cons:**

| Pros | Cons |
|------|------|
| No write amplification | Slow reads (merge 200+ lists at read time) |
| No celebrity problem | High read latency (200+ cache lookups per request) |
| Always fresh data | Expensive computation at read time |
| No wasted work for inactive users | Hard to apply ranking (need all candidates first) |

### 7.3 Hybrid Approach (Facebook/Twitter Solution)

**Candidate:** "In practice, Facebook and Twitter use a hybrid approach. This is what I'd recommend."

**Key Insight:** Treat users differently based on their follower count.

```
┌──────────────────────────────────────────────────────────────┐
│                    HYBRID FAN-OUT STRATEGY                    │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  Normal Users (< 10K followers):                             │
│    → Fan-out on WRITE (push to all followers)                │
│    → Followers get instant feed updates                      │
│                                                              │
│  Celebrities (>= 10K followers):                             │
│    → Do NOT fan out on write                                 │
│    → Their posts go to a "celebrity posts" index             │
│    → Merged into feeds at READ time                          │
│                                                              │
│  Feed Read = Pre-computed feed + Real-time celebrity merge   │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

**Architecture:**

```
User A (normal) posts          User X (celebrity) posts
       │                              │
       ▼                              ▼
┌──────────────┐              ┌──────────────────┐
│ Fan-out on   │              │ Celebrity Post   │
│ Write        │              │ Index (Redis)    │
│ → push to    │              │ → Stored by      │
│   all 200    │              │   author_id      │
│   followers  │              │ → NOT fanned out │
└──────┬───────┘              └────────┬─────────┘
       │                               │
       ▼                               │
┌──────────────┐                       │
│ Follower's   │                       │
│ Pre-computed │◄──────────────────────┘  (merged at read time)
│ Feed Cache   │
└──────────────┘
```

**Implementation:**

```python
# Hybrid Fan-out Strategy
CELEBRITY_THRESHOLD = 10_000

class HybridFanoutService:
    def process_post(self, event):
        author_id = event['author_id']
        post_id = event['post_id']
        created_at = event['created_at']
        follower_count = self.social_graph.get_follower_count(author_id)

        if follower_count >= CELEBRITY_THRESHOLD:
            # Celebrity: store in celebrity index, skip fan-out
            self.redis.zadd(
                f"celebrity_posts:{author_id}",
                {post_id: created_at}
            )
            self.redis.zremrangebyrank(
                f"celebrity_posts:{author_id}", 0, -101
            )
        else:
            # Normal user: fan-out on write
            self._fanout_to_followers(author_id, post_id, created_at)

    def _fanout_to_followers(self, author_id, post_id, created_at):
        cursor = None
        while True:
            followers, cursor = self.social_graph.get_followers(
                author_id, cursor=cursor, limit=1000
            )
            if not followers:
                break

            pipeline = self.redis.pipeline()
            for follower_id in followers:
                pipeline.zadd(f"feed:{follower_id}", {post_id: created_at})
                pipeline.zremrangebyrank(f"feed:{follower_id}", 0, -801)
            pipeline.execute()

            if cursor is None:
                break


class HybridFeedService:
    def get_feed(self, user_id, cursor=None, size=20):
        # Step 1: Fetch pre-computed feed (from fan-out on write)
        pre_computed = self.redis.zrevrangebyscore(
            f"feed:{user_id}",
            max=cursor or '+inf',
            min='-inf',
            start=0,
            num=size,
            withscores=True
        )

        # Step 2: Get celebrity posts for users this person follows
        following = self.social_graph.get_following(user_id)
        celebrities = [
            uid for uid in following
            if self.social_graph.get_follower_count(uid) >= CELEBRITY_THRESHOLD
        ]

        celebrity_posts = []
        for celeb_id in celebrities:
            posts = self.redis.zrevrangebyscore(
                f"celebrity_posts:{celeb_id}",
                max=cursor or '+inf',
                min='-inf',
                start=0,
                num=5,
                withscores=True
            )
            celebrity_posts.extend(posts)

        # Step 3: Merge and sort by score (timestamp)
        all_posts = list(pre_computed) + celebrity_posts
        all_posts.sort(key=lambda x: x[1], reverse=True)

        # Step 4: Hydrate and return top N
        top_posts = all_posts[:size]
        post_ids = [p[0] for p in top_posts]
        hydrated = self.post_service.get_posts_batch(post_ids)

        next_cursor = top_posts[-1][1] if top_posts else None
        return {
            'posts': hydrated,
            'next_cursor': next_cursor,
            'has_more': len(all_posts) > size
        }
```

### 7.4 Strategy Comparison

```
┌────────────────────┬────────────────┬────────────────┬────────────────┐
│ Metric             │ Fan-out Write  │ Fan-out Read   │ Hybrid         │
├────────────────────┼────────────────┼────────────────┼────────────────┤
│ Feed Read Latency  │ < 10ms         │ 200-500ms      │ < 50ms         │
│ Post Write Latency │ High (fan-out) │ < 10ms         │ Medium         │
│ Celebrity Impact   │ Catastrophic   │ None           │ Handled        │
│ Storage Cost       │ Very High      │ Low            │ Medium         │
│ Implementation     │ Simple         │ Moderate       │ Complex        │
│ Feed Freshness     │ 1-5 sec lag    │ Real-time      │ 1-5 sec lag    │
│ Used By            │ Twitter (old)  │ (Rare alone)   │ Facebook, X    │
└────────────────────┴────────────────┴────────────────┴────────────────┘
```

### 7.5 Feed Ranking

**Candidate:** "For V2, we'd add ML-based ranking. Here's a simplified version of how Facebook's EdgeRank-like system works."

**EdgeRank Formula (Simplified):**
```
Score = Σ (Affinity × Weight × Decay)

Where:
  Affinity  = How close the viewer is to the post author
              (interaction history: likes, comments, message frequency)
  Weight    = Content type weight (video > image > text > link)
  Decay     = Time decay: 1 / (1 + α × hours_since_posted)
```

**Implementation:**

```python
class FeedRanker:
    def __init__(self, feature_store, ml_model):
        self.feature_store = feature_store
        self.ml_model = ml_model

    def rank_feed(self, user_id, candidate_posts):
        features_batch = []

        for post in candidate_posts:
            features = self._extract_features(user_id, post)
            features_batch.append(features)

        # ML model predicts engagement probability
        scores = self.ml_model.predict_batch(features_batch)

        # Combine ML score with time decay
        ranked = []
        for post, score in zip(candidate_posts, scores):
            hours_old = (time.time() - post['created_at']) / 3600
            time_decay = 1.0 / (1.0 + 0.1 * hours_old)
            final_score = 0.7 * score + 0.3 * time_decay
            ranked.append((post, final_score))

        ranked.sort(key=lambda x: x[1], reverse=True)
        return [post for post, _ in ranked]

    def _extract_features(self, user_id, post):
        author_id = post['author_id']

        # Affinity: interaction history between viewer and author
        affinity = self.feature_store.get_affinity(user_id, author_id)

        # Content type weight
        type_weights = {'video': 1.5, 'image': 1.2, 'text': 1.0, 'link': 0.8}
        content_weight = type_weights.get(post['post_type'], 1.0)

        # Post engagement signals
        like_rate = post['like_count'] / max(post.get('impression_count', 1), 1)
        comment_rate = post['comment_count'] / max(post.get('impression_count', 1), 1)

        return {
            'affinity': affinity,
            'content_weight': content_weight,
            'like_rate': like_rate,
            'comment_rate': comment_rate,
            'author_follower_count': post.get('author_follower_count', 0),
            'post_age_hours': (time.time() - post['created_at']) / 3600,
            'has_media': 1 if post.get('media_urls') else 0,
        }
```

---

## 8. Scaling Strategies

### 8.1 Feed Cache Sharding

```
┌────────────────────────────────────────────────────────────────┐
│                    Redis Cluster (Feed Cache)                  │
│                                                                │
│  Shard by user_id using consistent hashing (CRC16 % 16384)    │
│                                                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │ Shard 0  │  │ Shard 1  │  │ Shard 2  │  │ Shard N  │      │
│  │ Users    │  │ Users    │  │ Users    │  │ Users    │      │
│  │ 0-31M    │  │ 31M-62M  │  │ 62M-93M  │  │ 469M-500M│      │
│  │          │  │          │  │          │  │          │      │
│  │ 64 GB    │  │ 64 GB    │  │ 64 GB    │  │ 64 GB    │      │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
│                                                                │
│  Total: 16 shards × 64 GB = 1 TB (with replicas: 3 TB)       │
└────────────────────────────────────────────────────────────────┘
```

### 8.2 Hot/Cold Storage Separation

```
┌────────────────────────────────────────────────────────┐
│                   HOT / COLD TIERING                   │
├────────────────────────────────────────────────────────┤
│                                                        │
│  HOT (Last 7 days):                                    │
│    → Redis Cluster (in-memory)                         │
│    → Feed items for active users                       │
│    → ~2 TB, sub-ms latency                             │
│                                                        │
│  WARM (7-30 days):                                     │
│    → Cassandra (SSD-backed)                            │
│    → Full materialized feed history                    │
│    → ~50 TB, 5-20ms latency                            │
│                                                        │
│  COLD (30+ days):                                      │
│    → S3 / HDFS                                         │
│    → Archived feeds for compliance and analytics       │
│    → ~500 TB, seconds latency                          │
│                                                        │
│  Policy: If Redis cache miss → read from Cassandra     │
│          → backfill Redis for that user                 │
│                                                        │
└────────────────────────────────────────────────────────┘
```

### 8.3 Fan-out Worker Scaling with Kafka

```
┌─────────────────────────────────────────────────────────────┐
│                 Kafka Fan-out Pipeline                       │
│                                                             │
│  Topic: post-events (128 partitions)                        │
│                                                             │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐     ┌─────────┐    │
│  │ Part 0  │  │ Part 1  │  │ Part 2  │ ... │ Part 127│    │
│  └────┬────┘  └────┬────┘  └────┬────┘     └────┬────┘    │
│       │            │            │                │         │
│       ▼            ▼            ▼                ▼         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐     ┌─────────┐    │
│  │Worker 0 │  │Worker 1 │  │Worker 2 │ ... │Worker N │    │
│  │         │  │         │  │         │     │         │    │
│  │ Fan out │  │ Fan out │  │ Fan out │     │ Fan out │    │
│  │ to Redis│  │ to Redis│  │ to Redis│     │ to Redis│    │
│  └─────────┘  └─────────┘  └─────────┘     └─────────┘    │
│                                                             │
│  Partition key: author_id                                   │
│  Consumer group: fanout-workers                             │
│  Auto-scaling: scale workers based on consumer lag          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 8.4 Social Graph Caching

```python
# Multi-layer caching for social graph
class CachedSocialGraph:
    def __init__(self, local_cache, redis, mysql):
        self.local_cache = local_cache  # LRU, 10K entries, 60s TTL
        self.redis = redis
        self.mysql = mysql

    def get_follower_count(self, user_id):
        # L1: Process-local cache
        count = self.local_cache.get(f"fc:{user_id}")
        if count is not None:
            return count

        # L2: Redis
        count = self.redis.get(f"user:{user_id}:follower_count")
        if count is not None:
            self.local_cache.set(f"fc:{user_id}", int(count))
            return int(count)

        # L3: MySQL (cold path)
        count = self.mysql.query(
            "SELECT COUNT(*) FROM user_follows WHERE followee_id = %s",
            (user_id,)
        )
        self.redis.setex(f"user:{user_id}:follower_count", 300, count)
        self.local_cache.set(f"fc:{user_id}", count)
        return count
```

### 8.5 Media Pipeline

```
┌──────────┐    ┌────────────────┐    ┌──────────────┐    ┌─────────┐
│ Client   │───▶│ Presigned URL  │───▶│ S3 Upload    │───▶│ Lambda  │
│ uploads  │    │ from API       │    │ (original)   │    │ trigger │
│ media    │    └────────────────┘    └──────────────┘    └────┬────┘
└──────────┘                                                   │
                                                               ▼
                                                    ┌──────────────────┐
                                                    │ Image Processing │
                                                    │ (resize, thumb,  │
                                                    │  compress, EXIF  │
                                                    │  strip)          │
                                                    └────────┬─────────┘
                                                             │
                                                             ▼
                                                    ┌──────────────────┐
                                                    │ CDN (CloudFront) │
                                                    │ thumbnail.jpg    │
                                                    │ medium.jpg       │
                                                    │ original.jpg     │
                                                    └──────────────────┘
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Failure Matrix

| # | Failure Scenario | Impact | Mitigation |
|---|-----------------|--------|------------|
| 1 | **Fan-out service lag** | Feeds become stale; new posts don't appear for minutes | Monitor Kafka consumer lag; auto-scale workers; set SLA alert at 30s lag |
| 2 | **Cache stampede on cold start** | Thousands of cache misses hit Cassandra simultaneously | Use request coalescing (singleflight pattern); probabilistic early expiration |
| 3 | **Social graph inconsistency** | User unfollows someone but still sees their posts | Async graph update + eventual feed cleanup; hard delete on next feed load |
| 4 | **Celebrity post thundering herd** | Celebrity posts 100M feed reads at once, spiking read load | Celebrity posts served from dedicated hot cache; rate-limit fan-out |
| 5 | **Feed cache eviction** | Inactive user's feed evicted; first load is slow | Fall back to pull-based generation; backfill cache asynchronously |
| 6 | **Post DB failure** | Cannot hydrate feed items | Read replica failover; serve degraded feed (show cached metadata only) |
| 7 | **Kafka broker failure** | Fan-out stalls | Multi-broker cluster with replication factor 3; dead letter queue for failures |
| 8 | **Redis cluster node failure** | Partial feed cache loss | Redis Cluster automatic failover to replicas; client retries on MOVED/ASK |

### 9.2 Cache Stampede Prevention

```javascript
// Singleflight pattern: coalesce concurrent cache misses for the same key
class SingleFlight {
  constructor() {
    this.inFlight = new Map();
  }

  async do(key, fetchFn) {
    if (this.inFlight.has(key)) {
      return this.inFlight.get(key);
    }

    const promise = fetchFn()
      .finally(() => this.inFlight.delete(key));

    this.inFlight.set(key, promise);
    return promise;
  }
}

// Usage in Feed Service
const singleFlight = new SingleFlight();

async function getFeed(userId) {
  const cacheKey = `feed:${userId}`;
  const cached = await redis.get(cacheKey);
  if (cached) return JSON.parse(cached);

  // Only one concurrent request per userId hits the database
  return singleFlight.do(cacheKey, async () => {
    const feed = await generateFeedFromDB(userId);
    await redis.setex(cacheKey, 300, JSON.stringify(feed));
    return feed;
  });
}
```

### 9.3 Graceful Degradation

```
┌──────────────────────────────────────────────────────────┐
│                DEGRADATION LADDER                         │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Level 0 (Normal):                                       │
│    Full ranked feed + celebrity merge + real-time updates │
│                                                          │
│  Level 1 (Ranking degraded):                             │
│    Chronological feed only (skip ML ranking)             │
│    → Triggered when: Ranking service latency > 200ms     │
│                                                          │
│  Level 2 (Cache-only):                                   │
│    Serve whatever is in Redis, no DB fallback            │
│    → Triggered when: Cassandra latency > 500ms           │
│                                                          │
│  Level 3 (Static fallback):                              │
│    Serve trending/popular posts (same for all users)     │
│    → Triggered when: Redis cluster partially down        │
│                                                          │
│  Level 4 (Service unavailable):                          │
│    Return cached client-side feed + "Service degraded"   │
│    → Triggered when: Multiple critical systems down      │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

| Category | Metric | Target | Alert Threshold |
|----------|--------|--------|-----------------|
| **Feed Freshness** | Time from post creation to feed appearance | < 5s (p99) | > 30s |
| **Feed Latency** | Feed API response time (p50 / p99) | 100ms / 500ms | p99 > 1s |
| **Cache Hit Rate** | Feed cache hit ratio | > 95% | < 85% |
| **Fan-out Lag** | Kafka consumer lag (messages behind) | < 10K | > 100K |
| **Fan-out Throughput** | Feed deliveries per second | ~115K/sec | Drop > 50% |
| **Post Write Latency** | Time to persist a new post | < 50ms | > 200ms |
| **Error Rate** | 5xx responses on feed endpoint | < 0.01% | > 0.1% |
| **Cache Eviction Rate** | Redis keys evicted per minute | < 100/min | > 1000/min |

### 10.2 Grafana Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────────┐
│                    NEWS FEED SYSTEM DASHBOARD                        │
├─────────────────────────────┬───────────────────────────────────────┤
│                             │                                       │
│  Feed Read Latency (p50/99) │  Feed Reads/sec (QPS)                │
│  ┌───────────────────────┐  │  ┌─────────────────────────────────┐  │
│  │  ──────────── p50     │  │  │            ▄▄▄                  │  │
│  │  - - - - - - - p99    │  │  │          ▄█████▄                │  │
│  │  50ms ──────          │  │  │  58K ──▄████████▄──             │  │
│  │  400ms - - - -        │  │  │       ██████████████            │  │
│  └───────────────────────┘  │  └─────────────────────────────────┘  │
│                             │                                       │
├─────────────────────────────┼───────────────────────────────────────┤
│                             │                                       │
│  Fan-out Kafka Lag          │  Cache Hit Rate                       │
│  ┌───────────────────────┐  │  ┌─────────────────────────────────┐  │
│  │ 5K ──────────         │  │  │  97% ─────────────────          │  │
│  │                       │  │  │                                 │  │
│  │  ▲ spike = concern    │  │  │  85% - - - - alert threshold    │  │
│  └───────────────────────┘  │  └─────────────────────────────────┘  │
│                             │                                       │
├─────────────────────────────┼───────────────────────────────────────┤
│                             │                                       │
│  Feed Freshness (post→feed) │  Error Rate (5xx)                    │
│  ┌───────────────────────┐  │  ┌─────────────────────────────────┐  │
│  │  2s ──────────        │  │  │  0.005% ──────────              │  │
│  │  30s - - - alert      │  │  │  0.1%   - - - - alert          │  │
│  └───────────────────────┘  │  └─────────────────────────────────┘  │
│                             │                                       │
└─────────────────────────────┴───────────────────────────────────────┘
```

### 10.3 Alerting Rules

```yaml
# Prometheus alerting rules
groups:
  - name: news-feed-alerts
    rules:
      - alert: FeedLatencyHigh
        expr: histogram_quantile(0.99, rate(feed_read_latency_seconds_bucket[5m])) > 1.0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Feed p99 latency exceeds 1 second"

      - alert: FanoutLagCritical
        expr: kafka_consumer_lag{topic="post-events", group="fanout-workers"} > 100000
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Fan-out Kafka lag exceeds 100K messages"

      - alert: CacheHitRateLow
        expr: rate(feed_cache_hits[5m]) / (rate(feed_cache_hits[5m]) + rate(feed_cache_misses[5m])) < 0.85
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Feed cache hit rate dropped below 85%"

      - alert: FeedFreshnessLag
        expr: histogram_quantile(0.99, rate(feed_freshness_seconds_bucket[5m])) > 30
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Posts taking >30s to appear in feeds"

      - alert: FanoutWorkerErrorRate
        expr: rate(fanout_worker_errors_total[5m]) / rate(fanout_worker_processed_total[5m]) > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Fan-out worker error rate exceeds 1%"
```

### 10.4 Distributed Tracing

```
Feed Request Trace (example):
─────────────────────────────────────────────────────────
│ API Gateway           │ 2ms    │ auth + routing       │
│  └─ Feed Service      │ 45ms   │ orchestration        │
│      ├─ Redis GET      │ 3ms    │ feed cache lookup    │
│      ├─ Celebrity Merge│ 12ms   │ 3 celebrity lookups  │
│      ├─ Ranking Service│ 15ms   │ ML scoring           │
│      └─ Post Hydration │ 10ms   │ batch post lookup    │
│          └─ Redis MGET │ 5ms    │ 20 post details      │
─────────────────────────────────────────────────────────
Total: 47ms (well within 500ms SLA)
```

---

## 11. Advanced Features

### 11.1 Infinite Scroll (Cursor-Based Pagination)

**Candidate:** "Offset-based pagination breaks when new posts are inserted mid-scroll. Cursor-based pagination solves this."

```javascript
// Cursor-based pagination using the last seen post's timestamp
async function getFeedPage(userId, cursor, size = 20) {
  const feedKey = `feed:${userId}`;

  if (cursor) {
    // Fetch posts older than cursor (which is a timestamp)
    const entries = await redis.zrevrangebyscore(
      feedKey,
      `(${cursor}`,  // exclusive upper bound
      '-inf',
      'LIMIT', 0, size + 1
    );

    const hasMore = entries.length > size;
    const posts = entries.slice(0, size);
    const nextCursor = posts.length > 0
      ? posts[posts.length - 1].score
      : null;

    return { posts, nextCursor, hasMore };
  } else {
    // First page: get latest N posts
    const entries = await redis.zrevrange(feedKey, 0, size);
    const nextCursor = entries.length > 0
      ? entries[entries.length - 1].score
      : null;

    return {
      posts: entries.slice(0, size),
      nextCursor,
      hasMore: entries.length > size,
    };
  }
}
```

### 11.2 "Seen" Tracking and New Post Indicator

```python
class SeenTracker:
    def mark_seen(self, user_id, post_ids):
        """Track which posts a user has seen (for deduplication and 'new' badge)."""
        key = f"seen:{user_id}"
        pipeline = self.redis.pipeline()
        for post_id in post_ids:
            pipeline.sadd(key, post_id)
        pipeline.expire(key, 7 * 86400)  # 7-day TTL
        pipeline.execute()

    def get_unseen_count(self, user_id):
        """Return count of new posts since last feed open."""
        last_seen_ts = self.redis.get(f"last_feed_open:{user_id}")
        if not last_seen_ts:
            return 0

        feed_key = f"feed:{user_id}"
        return self.redis.zcount(feed_key, f"({last_seen_ts}", "+inf")

    def record_feed_open(self, user_id):
        self.redis.set(f"last_feed_open:{user_id}", time.time())
```

### 11.3 Real-Time Feed Updates (WebSocket/SSE)

```javascript
// Server-Sent Events for real-time feed updates
class FeedSSEService {
  constructor(redisSubscriber) {
    this.connections = new Map(); // userId → SSE response
    this.redisSubscriber = redisSubscriber;
  }

  subscribe(userId, res) {
    res.writeHead(200, {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      'Connection': 'keep-alive',
    });

    this.connections.set(userId, res);

    // Subscribe to user's feed channel
    this.redisSubscriber.subscribe(`feed-updates:${userId}`);

    res.on('close', () => {
      this.connections.delete(userId);
      this.redisSubscriber.unsubscribe(`feed-updates:${userId}`);
    });
  }

  // Called by fan-out workers after writing to feed cache
  notifyNewPost(userId, postPreview) {
    const conn = this.connections.get(userId);
    if (conn) {
      conn.write(`event: new-post\n`);
      conn.write(`data: ${JSON.stringify(postPreview)}\n\n`);
    }
  }
}

// Client-side handling
const eventSource = new EventSource('/api/v1/feed/stream');
eventSource.addEventListener('new-post', (event) => {
  const post = JSON.parse(event.data);
  // Show "N new posts" banner instead of auto-inserting (avoids layout shift)
  showNewPostsBanner(post);
});
```

### 11.4 Content Filtering & Moderation

```
Post Created
     │
     ▼
┌──────────────────┐
│ Moderation       │
│ Pipeline (async) │
├──────────────────┤
│ 1. Spam filter   │ → ML classifier (text + image)
│ 2. NSFW detect   │ → Image classification model
│ 3. Hate speech   │ → NLP model (multi-language)
│ 4. Policy check  │ → Rules engine (links, keywords)
│ 5. Human review  │ → Queue for borderline cases
└────────┬─────────┘
         │
    ┌────▼─────┐
    │ Decision │
    ├──────────┤
    │ PASS     │ → Allow in feeds
    │ SOFT_BAN │ → Reduce distribution (lower rank score)
    │ BLOCK    │ → Remove from all feeds, notify author
    │ REVIEW   │ → Hold for human moderator
    └──────────┘
```

### 11.5 Trending Topics

```python
class TrendingService:
    def update_trending(self, post):
        """Extract hashtags and update trending counters."""
        hashtags = self.extract_hashtags(post['content'])
        now = int(time.time())
        window = now - (now % 3600)  # Hourly window

        pipeline = self.redis.pipeline()
        for tag in hashtags:
            # Sliding window: increment count for current hour
            pipeline.zincrby(f"trending:{window}", 1, tag)
            pipeline.expire(f"trending:{window}", 86400)
        pipeline.execute()

    def get_trending(self, top_n=10):
        """Aggregate last 6 hours of hashtag counts."""
        now = int(time.time())
        keys = [f"trending:{now - (now % 3600) - i * 3600}" for i in range(6)]

        # Weighted union: recent hours count more
        weights = [6, 5, 4, 3, 2, 1]
        self.redis.zunionstore("trending:aggregated", dict(zip(keys, weights)))
        return self.redis.zrevrange("trending:aggregated", 0, top_n - 1, withscores=True)
```

### 11.6 Story / Ephemeral Posts

```
Stories are posts with a 24-hour TTL:
  - Stored in Redis with EXPIRE set to 86400 seconds
  - Separate feed: GET /api/v1/stories
  - Rendered in a horizontal carousel above the main feed
  - Not indexed in Cassandra (no permanent storage)

Key design:
  feed:stories:{userId} → Redis ZSET with 24h TTL
  Each story entry also has its own TTL
  Stories from followed users merged at read time
```

---

## 12. Interview Q&A

### Q1: How do you handle the celebrity problem (user with 100M followers)?

**Candidate:** "The celebrity problem is the biggest challenge in feed systems. If a user with 100M followers posts, fan-out on write would mean 100M writes — this takes minutes and creates enormous write amplification.

**Solution: Hybrid approach.** I classify users into two tiers:
- **Normal users** (< 10K followers): Fan-out on write. Their posts are pushed to all followers' feed caches immediately.
- **Celebrities** (>= 10K followers): No fan-out. Their posts are stored in a dedicated celebrity post index.

When a user reads their feed, the Feed Service:
1. Reads the pre-computed feed from Redis (normal user posts).
2. Checks which celebrities the user follows.
3. Fetches recent posts from those celebrities' indices.
4. Merges and sorts both sets by time.

The threshold (10K) is configurable. The key insight is that a small fraction of users (< 0.1%) are celebrities, so the read-time merge adds minimal overhead — maybe 3-5 extra Redis lookups per feed read. The write savings are massive."

---

### Q2: Fan-out on write vs fan-out on read — when would you choose each?

**Candidate:** "The choice depends on the read/write ratio and the follower distribution.

**Fan-out on write** is better when:
- The system is read-heavy (100:1 ratio, like ours).
- Most users have a moderate follower count (< 10K).
- Feed latency requirements are strict (< 100ms).
- You can tolerate slightly stale feeds (1-5 second delay).

**Fan-out on read** is better when:
- Users follow a very large number of accounts (merging is bounded).
- Post freshness is critical (real-time feeds).
- You have a low read/write ratio.
- Write amplification cost exceeds compute cost.

**Hybrid** is the industry standard for any system at scale because the user distribution is always power-law: a few celebrities with millions of followers, and millions of normal users with < 1K followers.

In an interview, I'd always propose hybrid and explain why, then discuss the tradeoffs of the two extremes."

---

### Q3: How do you rank posts in the feed? Walk me through a ranking algorithm.

**Candidate:** "Feed ranking is essentially a recommendation problem. Here's the pipeline:

**Step 1 — Candidate Generation:**
Pull the top 500 candidate posts from the user's pre-computed feed and celebrity posts.

**Step 2 — Feature Extraction (per candidate):**
- **Affinity**: How often the viewer interacts with the post author (likes, comments, profile visits). Computed offline and stored in a feature store.
- **Content signals**: Post type (video > image > text), media quality score, text length.
- **Engagement signals**: Like rate, comment rate, share rate relative to impressions.
- **Time decay**: Exponential decay based on post age.
- **Social proof**: How many mutual friends engaged with this post.

**Step 3 — Scoring:**
An ML model (typically gradient-boosted trees or a lightweight neural net) predicts the probability of the viewer engaging with each post. The model is trained offline on historical engagement data.

**Step 4 — Post-processing:**
- Diversity injection: avoid showing 5 posts from the same author in a row.
- Content-type mixing: ensure a mix of text, images, and videos.
- Anti-echo-chamber rules: occasionally inject posts from outside the user's usual bubble.

A simplified scoring formula:
```
Score = 0.4 × P(like) + 0.3 × P(comment) + 0.2 × P(share) + 0.1 × P(click)
```

Facebook's original EdgeRank was `Σ(Affinity × Weight × Decay)`, but modern systems use deep learning models."

---

### Q4: How do you handle feed pagination for infinite scroll?

**Candidate:** "I use **cursor-based pagination** instead of offset-based pagination.

**Why not offset?** If I request page 2 with `offset=20` and a new post is inserted at the top while the user is reading page 1, the first post on page 2 would be a duplicate of the last post on page 1.

**Cursor-based approach:**
- The cursor is the timestamp (or score) of the last post the client saw.
- The next page request says 'give me 20 posts older than this timestamp.'
- In Redis: `ZREVRANGEBYSCORE feed:{userId} ({cursor} -inf LIMIT 0 20`
- This is stable regardless of new insertions.

**Edge cases:**
- Posts with identical timestamps: use a composite cursor `{timestamp}:{postId}`.
- Feed updates while scrolling: new posts accumulate at the top. The client can show a 'N new posts' banner to avoid disorienting the user.
- Very old cursor (user scrolls deep): may require switching from Redis to Cassandra for older feed items."

---

### Q5: How would you implement real-time feed updates?

**Candidate:** "I'd use **Server-Sent Events (SSE)** rather than WebSockets for this use case.

**Why SSE over WebSocket?**
- Feed updates are server → client only (unidirectional).
- SSE is simpler to implement and works through HTTP (no upgrade needed).
- Automatic reconnection built into the browser API.
- WebSocket would be overkill for push-only updates.

**Architecture:**
1. When the fan-out worker writes a post to a user's feed cache, it also publishes to a Redis Pub/Sub channel: `feed-updates:{userId}`.
2. The SSE server subscribes to this channel for connected users.
3. It pushes a lightweight notification: `{postId, authorName, preview}`.
4. The client shows a 'New posts available' banner (not auto-insert, to avoid scroll jumping).

**Scaling concern:** With 500M DAU, we can't keep SSE connections open for everyone. The solution:
- Only users actively viewing their feed get SSE connections.
- Connections are time-limited (auto-close after 5 minutes of inactivity).
- Use sticky sessions or a connection registry so the right SSE server gets the notification."

---

### Q6: How do you handle feed consistency when a user unfollows someone?

**Candidate:** "This is an interesting eventual consistency challenge. When User A unfollows User B:

**Immediate actions:**
1. Remove the follow relationship from MySQL and Redis.
2. The client-side immediately filters out B's posts from the local feed view.

**Async cleanup (best-effort):**
3. Publish an `UNFOLLOW` event to Kafka.
4. A cleanup worker removes B's posts from A's feed cache in Redis.
5. In Cassandra, B's posts are eventually cleaned up (or just expire via TTL).

**Why not synchronous cleanup?**
- A's feed might have hundreds of B's posts across months.
- Scanning and deleting them synchronously would be slow and could block the unfollow action.
- From the user's perspective, the client-side filter makes the change appear instant.

**Edge case — re-follow:**
- If A unfollows and re-follows B within the cleanup window, the cleanup worker should check the current follow state before deleting. This prevents accidentally removing posts the user should now see again."

---

### Q7: How would you implement content moderation in the feed pipeline?

**Candidate:** "Content moderation needs to happen at two stages:

**Stage 1 — Write-time moderation (inline, before fan-out):**
- Run the post through a fast ML classifier (< 50ms) that catches obvious spam, NSFW content, and hate speech.
- If confidence is high (> 95%), block immediately and notify the author.
- If confidence is borderline (50-95%), allow the post but flag it for human review and reduce its fan-out priority.

**Stage 2 — Async moderation pipeline:**
- After the post is published, a deeper analysis runs asynchronously.
- Image/video analysis for NSFW content, violence, misinformation.
- If flagged, the post is removed from all feeds it was already fanned out to — this requires publishing a `POST_REMOVED` event that triggers deletion from Redis and Cassandra feeds.

**Stage 3 — Read-time filtering:**
- Users can set content preferences (e.g., hide political posts).
- A lightweight filter at read time removes posts matching the user's block list.

**Trade-off:** Blocking every post for moderation before fan-out adds latency and blocks legitimate content. The industry standard is to fan out immediately and remove post-hoc if moderation flags it — accept that some users may briefly see violating content."

---

### Q8: How do you handle feed generation for a new user with no history?

**Candidate:** "The cold start problem. A new user has no follow graph, no engagement history, and no pre-computed feed.

**Solution — multi-stage onboarding feed:**

1. **Interest-based onboarding:** During signup, ask the user to select topics of interest. Use these to seed initial recommendations.

2. **Popular/trending feed:** For the first few sessions, the feed is populated with:
   - Trending posts in the user's selected interest categories.
   - Posts from 'suggested' accounts (high-quality content creators).
   - Geo-local popular content.

3. **Suggested follows:** Prominently recommend accounts to follow based on:
   - The user's contact list (with permission).
   - Popular accounts in selected interest areas.
   - Accounts followed by users with similar demographics.

4. **Gradual transition:** As the user follows accounts and engages with content, the personalized feed gradually takes over from the trending/suggested feed.

The key metric for a new user is **time to first meaningful feed** — we want it under 5 seconds after account creation. The trending feed should be pre-computed and cacheable so there's zero personalization latency for new users."

---

## 13. Production Checklist

### Pre-Launch

| # | Task | Status |
|---|------|--------|
| 1 | Load test fan-out pipeline at 2× expected peak (116K writes/sec) | ☐ |
| 2 | Load test feed reads at 3× expected peak (174K reads/sec) | ☐ |
| 3 | Verify Redis Cluster failover (kill a master, confirm auto-promotion) | ☐ |
| 4 | Verify Kafka consumer rebalancing (add/remove workers) | ☐ |
| 5 | Set up feed freshness monitoring (post → feed appearance lag) | ☐ |
| 6 | Implement circuit breakers on ranking service (fallback to chronological) | ☐ |
| 7 | Test cache stampede protection (singleflight pattern) | ☐ |
| 8 | Configure CDN for media delivery (cache headers, invalidation) | ☐ |
| 9 | Set up dead letter queue for failed fan-out events | ☐ |
| 10 | Implement rate limiting on post creation (anti-spam) | ☐ |

### Day 1

| # | Task | Status |
|---|------|--------|
| 1 | Monitor Kafka consumer lag dashboards (alert on > 50K lag) | ☐ |
| 2 | Watch feed latency p99 (target: < 500ms) | ☐ |
| 3 | Monitor Redis memory usage and eviction rates | ☐ |
| 4 | Verify content moderation pipeline is processing in < 60s | ☐ |
| 5 | Check error rates on all services (< 0.01%) | ☐ |

### Week 1

| # | Task | Status |
|---|------|--------|
| 1 | Review fan-out worker scaling policies (are auto-scaling thresholds correct?) | ☐ |
| 2 | Analyze cache hit rates and adjust TTLs | ☐ |
| 3 | Review slow query logs on Posts DB and Social Graph DB | ☐ |
| 4 | A/B test feed ranking model vs chronological (engagement metrics) | ☐ |
| 5 | Identify top 1% users by follower count — verify they're classified as celebrities | ☐ |

### Month 1

| # | Task | Status |
|---|------|--------|
| 1 | Analyze hot partition keys in Cassandra (rebalance if needed) | ☐ |
| 2 | Review storage costs (S3, Cassandra, Redis) — optimize retention policies | ☐ |
| 3 | Train and deploy feed ranking ML model V2 based on real engagement data | ☐ |
| 4 | Implement real-time feed updates (SSE) for active users | ☐ |
| 5 | Plan multi-region deployment for < 100ms global feed latency | ☐ |
| 6 | Conduct chaos engineering (kill services, simulate network partitions) | ☐ |

---

## Summary

### System Overview

| Dimension | Decision |
|-----------|----------|
| **Core Strategy** | Hybrid fan-out (push for normal users, pull for celebrities) |
| **Feed Storage** | Redis Cluster (hot) + Cassandra (warm) + S3 (cold) |
| **Post Storage** | MySQL sharded by author_id |
| **Social Graph** | MySQL + Redis sets for fast lookups |
| **Message Bus** | Kafka (128 partitions, fan-out workers as consumers) |
| **Media Pipeline** | S3 + Lambda processing + CloudFront CDN |
| **Ranking** | ML-based scoring with time decay (V2) |
| **Real-time Updates** | Server-Sent Events for active sessions |
| **Pagination** | Cursor-based (timestamp + post_id) |

### Scalability Path

```
┌────────────────────────────────────────────────────────────────────┐
│                      SCALABILITY ROADMAP                           │
├──────────┬─────────────────────────────────────────────────────────┤
│          │                                                         │
│  V1      │  Single region, chronological feed, fan-out on write   │
│  (MVP)   │  MySQL + Redis, basic Kafka pipeline                   │
│          │  Target: 1M DAU                                         │
│          │                                                         │
├──────────┼─────────────────────────────────────────────────────────┤
│          │                                                         │
│  V2      │  Hybrid fan-out (celebrity threshold), ML ranking      │
│  (Scale) │  Add Cassandra for durable feed store                  │
│          │  Sharded MySQL, Redis Cluster                           │
│          │  Target: 50M DAU                                        │
│          │                                                         │
├──────────┼─────────────────────────────────────────────────────────┤
│          │                                                         │
│  V3      │  Multi-region deployment, global feed consistency      │
│  (Global)│  Real-time updates (SSE), advanced content moderation  │
│          │  Feature store for ML ranking, A/B testing framework   │
│          │  Target: 500M DAU                                       │
│          │                                                         │
├──────────┼─────────────────────────────────────────────────────────┤
│          │                                                         │
│  V4      │  Graph-based social recommendations                    │
│  (AI)    │  Deep learning ranking (transformer models)            │
│          │  Personalized notification timing                      │
│          │  Cross-platform feed unification                       │
│          │  Target: 1B+ DAU                                        │
│          │                                                         │
└──────────┴─────────────────────────────────────────────────────────┘
```

### Key Interview Talking Points

1. **Always propose hybrid fan-out** — it shows you understand the tradeoffs.
2. **Quantify the celebrity problem** — "100M followers × 1 post = 100M writes."
3. **Cache is king** — feed reads must come from Redis, never the database.
4. **Eventual consistency is acceptable** — feeds don't need strong consistency.
5. **Ranking is a separate concern** — design the feed pipeline first, add ranking as a layer.
6. **Cursor-based pagination** — never use offset for infinite scroll.
7. **Graceful degradation** — have a plan for when ranking/cache/fanout fails.
8. **Media is separate** — don't mix media storage with feed logic.

---

> **Interview Duration Guide:**  
> - Requirements + Scale: 10 min  
> - High-Level Design: 10 min  
> - Fan-out Deep Dive: 15 min (this is where you differentiate)  
> - Database + Caching: 10 min  
> - Q&A: 15 min
