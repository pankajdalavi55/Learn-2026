# AWS Interview Questions by Experience Level

> Comprehensive AWS interview preparation guide organized by experience level with detailed answers

---

## Table of Contents

1. [Junior Level (0-2 Years)](#1-junior-level-0-2-years)
2. [Mid-Level (2-5 Years)](#2-mid-level-2-5-years)
3. [Senior Level (5-8 Years)](#3-senior-level-5-8-years)
4. [Principal/Architect Level (8+ Years)](#4-principalarchitect-level-8-years)
5. [Java + AWS Interview Questions](#5-java--aws-interview-questions)
6. [Scenario-Based Questions](#6-scenario-based-questions)
7. [Quick Reference Cheat Sheet](#7-quick-reference-cheat-sheet)

---

## 1. Junior Level (0-2 Years)

### Core Services

**Q1: What is AWS and what are its main benefits?**

**A:** AWS (Amazon Web Services) is a comprehensive cloud computing platform providing:

- **On-demand resources:** Pay only for what you use
- **Scalability:** Scale up/down based on demand
- **Global reach:** 30+ regions worldwide
- **Security:** Enterprise-grade security features
- **Reliability:** High availability with SLAs

```
Traditional vs Cloud:
─────────────────────

Traditional:              AWS Cloud:
┌─────────────────┐       ┌─────────────────┐
│ Buy servers     │       │ Rent capacity   │
│ Long setup time │  vs   │ Minutes to start│
│ Fixed capacity  │       │ Auto-scale      │
│ CapEx heavy     │       │ OpEx model      │
└─────────────────┘       └─────────────────┘
```

---

**Q2: Explain the difference between EC2 and Lambda.**

**A:**

| Aspect | EC2 | Lambda |
|--------|-----|--------|
| **Type** | Virtual Machine | Serverless Function |
| **Management** | You manage OS, patching | AWS manages everything |
| **Billing** | Per hour/second (running) | Per request + duration |
| **Scaling** | Manual/Auto Scaling Groups | Automatic (1000s concurrent) |
| **Max Runtime** | Unlimited | 15 minutes |
| **Use Case** | Long-running apps, full control | Event-driven, microservices |

```python
# Lambda Example
def lambda_handler(event, context):
    name = event.get('name', 'World')
    return {
        'statusCode': 200,
        'body': f'Hello, {name}!'
    }
```

**Deep Dive - Lambda Execution Model:**

```
Lambda Lifecycle:
─────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Lambda Execution Environment                      │
│                                                                      │
│   COLD START (First invocation or after idle):                      │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │ 1. Download code │ 2. Start runtime │ 3. Run init code       │  │
│   │    (from S3)     │    (JVM, Python) │    (outside handler)   │  │
│   │    ~100-500ms    │    ~100-1000ms   │    Variable            │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                              │                                       │
│                              ▼                                       │
│   WARM START (Reused execution environment):                        │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │ Execute handler only (~1-10ms overhead)                      │  │
│   │ • Variables persist between invocations                      │  │
│   │ • DB connections can be reused                               │  │
│   │ • /tmp storage (512MB-10GB) persists                         │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Cold Start Factors:
• Runtime: Java/C# > Python/Node.js
• Package size: Larger = slower
• VPC: Adds ENI attachment time (~1-2s historically, now faster)
• Memory: More memory = more CPU = faster init
```

**Concurrency Model:**
- **Reserved Concurrency:** Guarantees capacity, limits max
- **Provisioned Concurrency:** Pre-warms instances (eliminates cold starts)
- **Burst Concurrency:** Initial burst of 500-3000 (region dependent)
- **Scaling Rate:** 500 additional instances/minute after burst

**Cost Formula:**
```
Cost = (Requests × $0.20/million) + (GB-seconds × $0.0000166667)
GB-seconds = Memory(GB) × Duration(seconds) × Invocations
```

---

**Q3: What is S3? Explain storage classes.**

**A:** S3 (Simple Storage Service) is object storage with 99.999999999% (11 9's) durability.

```
S3 Storage Classes:
───────────────────

┌────────────────────┬─────────────┬──────────────┬─────────────────────┐
│ Storage Class      │ Availability│ Min Duration │ Use Case            │
├────────────────────┼─────────────┼──────────────┼─────────────────────┤
│ Standard           │ 99.99%      │ None         │ Frequent access     │
│ Intelligent-Tier   │ 99.9%       │ None         │ Unknown patterns    │
│ Standard-IA        │ 99.9%       │ 30 days      │ Infrequent access   │
│ One Zone-IA        │ 99.5%       │ 30 days      │ Recreatable data    │
│ Glacier Instant    │ 99.9%       │ 90 days      │ Archive, quick access│
│ Glacier Flexible   │ 99.99%      │ 90 days      │ Archive, hours      │
│ Glacier Deep       │ 99.99%      │ 180 days     │ 7-10 year retention │
└────────────────────┴─────────────┴──────────────┴─────────────────────┘

Cost (per GB/month): Standard > IA > Glacier
Retrieval time:      Instant ◄─────────────► Hours
```

**Deep Dive - S3 Architecture & Consistency:**

```
S3 Internal Architecture:
─────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                         S3 Request Flow                              │
│                                                                      │
│  Client ───► S3 Frontend ───► Index Tier ───► Storage Nodes        │
│              (Load Balance)   (Metadata)      (Actual data)         │
│                                                                      │
│  • Objects split into chunks, replicated across devices             │
│  • Minimum 3 AZ replication for Standard class                      │
│  • 11 9's durability = lose 1 object per 10 million in 10,000 years │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Consistency Model (Strong Read-After-Write since Dec 2020):
──────────────────────────────────────────────────────────
• PUT new object → Immediately readable with latest data
• PUT overwrite → Immediately readable with latest data  
• DELETE → Immediately consistent (404 on next read)
• LIST → Eventually consistent (may take seconds to reflect changes)
```

**S3 Performance Optimization:**

| Technique | Benefit | When to Use |
|-----------|---------|-------------|
| **Multipart Upload** | Parallel uploads, retry parts | Files > 100MB |
| **S3 Transfer Acceleration** | CloudFront edge locations | Cross-region uploads |
| **Byte-Range Fetches** | Parallel downloads | Large file downloads |
| **Prefix Distribution** | 3,500 PUT/s, 5,500 GET/s per prefix | High-throughput apps |

**Lifecycle Policies - Cost Optimization:**
```json
{
  "Rules": [{
    "ID": "MoveToIA",
    "Status": "Enabled",
    "Transitions": [
      {"Days": 30, "StorageClass": "STANDARD_IA"},
      {"Days": 90, "StorageClass": "GLACIER"},
      {"Days": 365, "StorageClass": "DEEP_ARCHIVE"}
    ],
    "Expiration": {"Days": 730}
  }]
}
```

**S3 Security Layers:**
1. **Bucket Policies:** Resource-based (who can access this bucket)
2. **IAM Policies:** Identity-based (what can this user access)
3. **ACLs:** Legacy, object-level (avoid using)
4. **Block Public Access:** Account/bucket level safety net
5. **VPC Endpoints:** Private access without internet

---

**Q4: What is a VPC and its core components?**

**A:** VPC (Virtual Private Cloud) is your isolated network in AWS.

```
VPC Architecture:
─────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        VPC (10.0.0.0/16)                            │
│                                                                      │
│    ┌─────────────────────────┐    ┌─────────────────────────┐      │
│    │    Public Subnet        │    │    Private Subnet       │      │
│    │    10.0.1.0/24          │    │    10.0.2.0/24          │      │
│    │                         │    │                         │      │
│    │   ┌─────────────────┐   │    │   ┌─────────────────┐   │      │
│    │   │   Web Server    │   │    │   │   Database      │   │      │
│    │   │   (Public IP)   │   │    │   │   (Private IP)  │   │      │
│    │   └────────┬────────┘   │    │   └────────┬────────┘   │      │
│    │            │            │    │            │            │      │
│    └────────────┼────────────┘    └────────────┼────────────┘      │
│                 │                              │                    │
│    ┌────────────▼────────────┐    ┌────────────▼────────────┐      │
│    │   Internet Gateway      │    │   NAT Gateway           │      │
│    │   (Outbound + Inbound)  │    │   (Outbound only)       │      │
│    └─────────────────────────┘    └─────────────────────────┘      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Core Components:
• Subnets: Subdivisions of VPC (public/private)
• Route Tables: Control traffic routing
• Internet Gateway: Connect to internet
• NAT Gateway: Allow private subnet outbound access
• Security Groups: Instance-level firewall (stateful)
• NACLs: Subnet-level firewall (stateless)
```

---

**Q5: What is the difference between Security Groups and NACLs?**

**A:**

| Feature | Security Group | NACL |
|---------|---------------|------|
| **Level** | Instance (ENI) | Subnet |
| **State** | Stateful (return traffic auto-allowed) | Stateless (must define both) |
| **Rules** | Allow only | Allow and Deny |
| **Evaluation** | All rules evaluated | Rules evaluated in order |
| **Default** | Deny all inbound, allow all outbound | Allow all |

```
Security Group (Stateful):
─────────────────────────
Inbound: Allow port 80
Outbound: (Automatically allows response)

NACL (Stateless):
─────────────────
Inbound Rule 100: Allow port 80
Outbound Rule 100: Allow ephemeral ports (1024-65535) ← Must define!
```

---

**Q6: What is IAM? Explain users, groups, roles, and policies.**

**A:** IAM (Identity and Access Management) controls who can access what.

```
IAM Components:
───────────────

┌─────────────────────────────────────────────────────────────────────┐
│                            IAM                                       │
│                                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐│
│  │    Users    │  │   Groups    │  │    Roles    │  │  Policies   ││
│  │             │  │             │  │             │  │             ││
│  │  • People   │  │ • Collection│  │ • Temporary │  │ • JSON docs ││
│  │  • Apps     │  │   of users  │  │   credentials│  │ • Define    ││
│  │  • Long-term│  │ • Share     │  │ • EC2, Lambda│  │   permissions│
│  │   credentials│  │   permissions│  │ • Cross-acct│  │             ││
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Policy Example:
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::my-bucket/*"
  }]
}
```

**Best Practices:**
- Never use root account for daily tasks
- Enable MFA for all users
- Use roles for applications (not access keys)
- Follow principle of least privilege

**Deep Dive - IAM Policy Evaluation Logic:**

```
IAM Policy Evaluation Flow:
───────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                     Request Arrives                                  │
│                          │                                           │
│                          ▼                                           │
│              ┌───────────────────────┐                              │
│              │ Explicit DENY exists? │                              │
│              └───────────┬───────────┘                              │
│                    Yes   │   No                                      │
│                    │     │                                           │
│               ┌────▼────┐│                                           │
│               │  DENY   ││                                           │
│               └─────────┘│                                           │
│                          ▼                                           │
│              ┌───────────────────────┐                              │
│              │   SCP Allows action?  │ (If in AWS Organizations)    │
│              └───────────┬───────────┘                              │
│                    No    │   Yes                                     │
│                    │     │                                           │
│               ┌────▼────┐│                                           │
│               │  DENY   ││                                           │
│               └─────────┘│                                           │
│                          ▼                                           │
│              ┌───────────────────────┐                              │
│              │ Permission Boundary   │ (Intersection)               │
│              │      Allows?          │                              │
│              └───────────┬───────────┘                              │
│                    No    │   Yes                                     │
│                    │     │                                           │
│               ┌────▼────┐│                                           │
│               │  DENY   ││                                           │
│               └─────────┘│                                           │
│                          ▼                                           │
│              ┌───────────────────────┐                              │
│              │ Explicit ALLOW exists?│ (Identity or Resource policy)│
│              └───────────┬───────────┘                              │
│                    No    │   Yes                                     │
│                    │     │                                           │
│               ┌────▼────┐│┌────▼────┐                               │
│               │  DENY   │││  ALLOW  │                               │
│               └─────────┘│└─────────┘                               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Key Rule: Explicit DENY always wins, default is implicit DENY
```

**IAM Policy Types:**

| Policy Type | Attached To | Use Case |
|-------------|-------------|----------|
| **Identity-based** | Users, Groups, Roles | "What can this identity do?" |
| **Resource-based** | S3, SQS, KMS, etc. | "Who can access this resource?" |
| **Permission Boundary** | Users, Roles | Maximum permissions cap |
| **SCP** | AWS Organization OUs | Guardrails for accounts |
| **Session Policy** | Assumed role sessions | Temporary restrictions |

**Cross-Account Access Pattern:**
```
Account A                          Account B
┌─────────────────────┐           ┌─────────────────────┐
│                     │           │                     │
│  IAM User/Role      │           │  Resource (S3)      │
│  + Identity Policy  │──────────►│  + Resource Policy  │
│  "Allow sts:AssumeRole"         │  "Allow Account A"  │
│                     │           │                     │
└─────────────────────┘           └─────────────────────┘

BOTH policies must allow for cross-account access to work!
```

---

**Q7: What is an Availability Zone and Region?**

**A:**

```
AWS Global Infrastructure:
──────────────────────────

Region (e.g., us-east-1)
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│   AZ-a (us-east-1a)      AZ-b (us-east-1b)      AZ-c (us-east-1c)  │
│   ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐│
│   │  Data Center 1  │    │  Data Center 3  │    │  Data Center 5  ││
│   │  Data Center 2  │    │  Data Center 4  │    │  Data Center 6  ││
│   └────────┬────────┘    └────────┬────────┘    └────────┬────────┘│
│            │                      │                      │          │
│            └──────────────────────┴──────────────────────┘          │
│                    Low-latency private links                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

• Region: Geographic area (e.g., US East, EU West)
• Availability Zone: Isolated data center(s) within a region
• Edge Location: CloudFront cache points (400+)

Design for HA: Deploy across multiple AZs
```

---

**Q8: What is RDS? Compare RDS vs self-managed database on EC2.**

**A:** RDS (Relational Database Service) is managed database service.

| Feature | RDS | EC2 Self-Managed |
|---------|-----|------------------|
| **Setup** | Minutes | Hours/Days |
| **Patching** | Automated | Manual |
| **Backups** | Automated, point-in-time | Manual setup |
| **HA/Failover** | Multi-AZ option | Complex setup |
| **Scaling** | One-click | Manual migration |
| **Control** | Limited | Full |
| **Cost** | Higher | Lower (more effort) |

**RDS Engines:** MySQL, PostgreSQL, MariaDB, Oracle, SQL Server, Aurora

**Deep Dive - RDS Architecture:**

```
RDS Multi-AZ Deployment:
────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                           Region                                      │
│                                                                      │
│   Availability Zone A         Availability Zone B                   │
│   ┌────────────────────┐       ┌────────────────────┐              │
│   │  PRIMARY          │       │  STANDBY          │              │
│   │  (Read/Write)     │       │  (Sync replica)   │              │
│   │                   │──sync─►│                   │              │
│   │  ┌─────────────┐   │       │  ┌─────────────┐   │              │
│   │  │   EBS       │   │       │  │   EBS       │   │              │
│   │  │   Volume    │   │       │  │   Volume    │   │              │
│   │  └─────────────┘   │       │  └─────────────┘   │              │
│   └────────────────────┘       └────────────────────┘              │
│                                                                      │
│   Failover: Automatic (60-120 seconds), DNS endpoint unchanged      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Read Replicas (Scaling reads):
──────────────────────────────
• Async replication (eventual consistency)
• Up to 5 read replicas per instance
• Can be in different regions (cross-region replica)
• Can be promoted to standalone database
```

**Aurora vs Standard RDS:**

| Feature | Aurora | Standard RDS |
|---------|--------|-------------|
| **Storage** | Auto-scales to 128TB | Manual provisioning |
| **Replication** | 6 copies across 3 AZs | 1 copy in standby AZ |
| **Failover** | ~30 seconds | 60-120 seconds |
| **Read Replicas** | Up to 15 | Up to 5 |
| **Performance** | 5x MySQL, 3x PostgreSQL | Baseline |
| **Cost** | ~20% more | Baseline |

**Backup Strategies:**
- **Automated Backups:** Daily snapshots + transaction logs (point-in-time recovery)
- **Retention:** 0-35 days (0 disables backups)
- **Manual Snapshots:** Persist until deleted, can share cross-account
- **Backup Window:** Schedule during low traffic (avoid I/O freeze impact)

---

**Q9: What is CloudWatch? What can you monitor?**

**A:** CloudWatch is AWS's monitoring and observability service.

```
CloudWatch Components:
──────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        CloudWatch                                    │
│                                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐│
│  │   Metrics   │  │    Logs     │  │   Alarms    │  │  Dashboards ││
│  │             │  │             │  │             │  │             ││
│  │ • CPU       │  │ • App logs  │  │ • Threshold │  │ • Visualize ││
│  │ • Network   │  │ • System    │  │ • Actions   │  │ • Share     ││
│  │ • Custom    │  │ • Lambda    │  │ • SNS       │  │ • Custom    ││
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Common Metrics:
• EC2: CPUUtilization, NetworkIn/Out, DiskReadOps
• RDS: DatabaseConnections, FreeStorageSpace, ReadLatency
• Lambda: Invocations, Duration, Errors, Throttles
• ALB: RequestCount, TargetResponseTime, HTTPCode_Target_5XX
```

**Deep Dive - CloudWatch Architecture:**

```
CloudWatch Data Flow:
─────────────────────

                    ┌───────────────────────┐
                    │    CloudWatch        │
   Data Sources     │                       │     Actions
   ────────────     │  ┌───────────────┐  │     ───────
                    │  │   Metrics    │  │
   EC2 ───────────►│  │ (Time-series)│  │───► Alarms ──► SNS
   RDS ───────────►│  └───────────────┘  │───► Auto Scaling
   Lambda ────────►│                       │───► EC2 Action
   Custom ────────►│  ┌───────────────┐  │───► Lambda
                    │  │     Logs     │  │
   App Logs ───────►│  │ (Log groups) │  │───► Insights
   VPC Flow ───────►│  └───────────────┘  │───► S3 Export
                    │                       │
                    └───────────────────────┘
```

**Metric Types:**

| Type | Resolution | Retention | Cost |
|------|------------|-----------|------|
| **Standard** | 1 minute (EC2), 5 min (others) | 15 months | Free |
| **High-Resolution** | 1 second | 3 hours at 1s, then aggregated | Extra cost |
| **Custom** | Your choice | Same as standard | Per metric/month |

**CloudWatch Logs Insights Query Language:**
```sql
-- Find errors in Lambda logs
fields @timestamp, @message
| filter @message like /ERROR/
| sort @timestamp desc
| limit 20

-- Calculate error rate
stats count(*) as total,
      sum(strcontains(@message, 'ERROR')) as errors,
      (sum(strcontains(@message, 'ERROR')) / count(*)) * 100 as error_rate
by bin(1h)
```

**Alarm States:** `OK` → `ALARM` → `INSUFFICIENT_DATA`

**Key Concepts:**
- **Namespace:** Container for metrics (e.g., `AWS/EC2`)
- **Dimension:** Name-value pair to identify metric (e.g., `InstanceId=i-123`)
- **Period:** Time length for aggregation (60s, 300s, etc.)
- **Statistic:** Aggregation function (Average, Sum, Min, Max, SampleCount)

---

**Q10: What is Auto Scaling and how does it work?**

**A:** Auto Scaling automatically adjusts capacity based on demand.

```
Auto Scaling Group (ASG):
─────────────────────────

              CloudWatch Alarm
              (CPU > 80%)
                    │
                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Auto Scaling Group                                │
│                                                                      │
│    Min: 2          Desired: 4           Max: 10                     │
│                                                                      │
│    ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐     ┌─────┐ ┌─────┐           │
│    │ EC2 │ │ EC2 │ │ EC2 │ │ EC2 │ --> │ EC2 │ │ EC2 │ (scaling) │
│    └─────┘ └─────┘ └─────┘ └─────┘     └─────┘ └─────┘           │
│                                                                      │
└───────────────────────────────┬─────────────────────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Load Balancer       │
                    └───────────────────────┘

Scaling Policies:
• Target Tracking: Maintain 70% CPU
• Step Scaling: Add 2 if CPU > 80%, add 4 if > 90%
• Scheduled: Scale up at 9 AM, down at 6 PM
```

**Deep Dive - Auto Scaling Internals:**

```
Auto Scaling Decision Process:
──────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                     CloudWatch Metric                                │
│                          │                                           │
│                          ▼                                           │
│                   ┌────────────────┐                                 │
│                   │  Alarm Breach  │                                 │
│                   └────────┬───────┘                                 │
│                            │                                         │
│                            ▼                                         │
│                   ┌────────────────┐                                 │
│                   │ Cooldown Period│ (Default: 300s)                 │
│                   │ Check: Active? │                                 │
│                   └────────┬───────┘                                 │
│                    No     │   Yes                                    │
│                    │      │                                          │
│                    ▼      └──► (Wait, ignore scaling)                │
│           ┌────────────────┐                                        │
│           │ Calculate New  │                                        │
│           │ Desired Capacity│                                        │
│           └────────┬───────┘                                        │
│                    │                                                 │
│                    ▼                                                 │
│           ┌────────────────┐                                        │
│           │ Enforce Min/Max│                                        │
│           └────────┬───────┘                                        │
│                    │                                                 │
│                    ▼                                                 │
│           ┌────────────────┐                                        │
│           │ Launch/Terminate│ instances across AZs                  │
│           └────────────────┘                                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Scaling Policy Comparison:**

| Policy Type | How It Works | Best For |
|-------------|--------------|----------|
| **Target Tracking** | Maintains metric at target (like thermostat) | Steady workloads |
| **Step Scaling** | Adds/removes specific amounts based on alarm severity | Variable traffic spikes |
| **Simple Scaling** | Single adjustment, waits for cooldown | Basic use cases |
| **Scheduled** | Scales at specific times | Predictable patterns |
| **Predictive** | ML-based, forecasts load | Cyclical workloads |

**Key Concepts:**

```
Cooldown Period:
────────────────
• Prevents rapid scale in/out oscillation
• Default: 300 seconds
• During cooldown, ASG ignores new scaling activities
• Set based on instance startup time

Warm Pools (Reduce scale-out time):
───────────────────────────────────
• Pre-initialized instances in Stopped or Running state
• Faster than launching from AMI
• Pay for EBS storage only (Stopped state)

Instance Health Checks:
───────────────────────
• EC2 Status Check: Hardware/software issues
• ELB Health Check: Application-level health
• Custom Health: Via API (SetInstanceHealth)

Termination Policies:
─────────────────────
1. Default: Oldest launch config, then closest to billing hour
2. OldestInstance: Oldest first
3. NewestInstance: Newest first
4. OldestLaunchTemplate: Oldest template version
```

**Launch Template vs Launch Configuration:**
- Launch Configuration: Legacy, immutable, being deprecated
- Launch Template: Versioned, supports mixed instances, spot options

---

## 2. Mid-Level (2-5 Years)

### Architecture & Design

**Q11: Design a highly available 3-tier web application.**

**A:**

```
3-Tier Architecture:
────────────────────

                    ┌─────────────────────────┐
                    │      Route 53           │
                    │   (DNS + Health Check)  │
                    └───────────┬─────────────┘
                                │
                    ┌───────────▼───────────┐
                    │      CloudFront        │
                    │   (CDN + WAF)          │
                    └───────────┬───────────┘
                                │
┌───────────────────────────────┼───────────────────────────────────┐
│                               │                VPC                 │
│                   ┌───────────▼───────────┐                       │
│                   │   Application Load    │                       │
│                   │      Balancer         │                       │
│                   └───────────┬───────────┘                       │
│                               │                                    │
│           ┌───────────────────┼───────────────────┐               │
│           │                   │                   │               │
│   ┌───────▼───────┐   ┌───────▼───────┐   ┌───────▼───────┐     │
│   │   AZ-1a       │   │   AZ-1b       │   │   AZ-1c       │     │
│   │               │   │               │   │               │     │
│   │ ┌───────────┐ │   │ ┌───────────┐ │   │ ┌───────────┐ │     │
│   │ │  Web/App  │ │   │ │  Web/App  │ │   │ │  Web/App  │ │     │
│   │ │  (ASG)    │ │   │ │  (ASG)    │ │   │ │  (ASG)    │ │     │
│   │ └─────┬─────┘ │   │ └─────┬─────┘ │   │ └─────┬─────┘ │     │
│   │       │       │   │       │       │   │       │       │     │
│   │ ┌─────▼─────┐ │   │ ┌─────▼─────┐ │   │ ┌─────▼─────┐ │     │
│   │ │ElastiCache│ │   │ │ElastiCache│ │   │ │ElastiCache│ │     │
│   │ │ (Replica) │ │   │ │ (Primary) │ │   │ │ (Replica) │ │     │
│   │ └───────────┘ │   │ └───────────┘ │   │ └───────────┘ │     │
│   │               │   │               │   │               │     │
│   │ ┌───────────┐ │   │ ┌───────────┐ │   │ ┌───────────┐ │     │
│   │ │  Aurora   │ │   │ │  Aurora   │ │   │ │  Aurora   │ │     │
│   │ │ (Replica) │ │   │ │ (Writer)  │ │   │ │ (Replica) │ │     │
│   │ └───────────┘ │   │ └───────────┘ │   │ └───────────┘ │     │
│   └───────────────┘   └───────────────┘   └───────────────┘     │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘

Key Components:
• Route 53: DNS failover, health checks
• CloudFront: Global CDN, DDoS protection
• ALB: Layer 7 routing, TLS termination
• ASG: Auto-scale web/app tier
• ElastiCache: Session storage, caching
• Aurora: Multi-AZ database with read replicas
```

---

**Q12: Compare Application Load Balancer, Network Load Balancer, and Gateway Load Balancer.**

**A:**

| Feature | ALB | NLB | GWLB |
|---------|-----|-----|------|
| **Layer** | 7 (HTTP/HTTPS) | 4 (TCP/UDP) | 3 (IP) |
| **Latency** | ~400ms | ~100μs | Variable |
| **Throughput** | Millions req/sec | Millions req/sec | Highest |
| **SSL Termination** | Yes | Yes (TLS) | No |
| **Path Routing** | Yes | No | No |
| **WebSocket** | Yes | Yes | No |
| **Use Case** | Web apps, microservices | Gaming, IoT, TCP | Security appliances |

```
Load Balancer Selection:
────────────────────────

HTTP/HTTPS API? ─────Yes────► ALB
       │
       No
       │
Ultra-low latency? ──Yes────► NLB
       │
       No
       │
Virtual appliance? ──Yes────► Gateway LB
```

---

**Q13: What is DynamoDB? When would you use it over RDS?**

**A:** DynamoDB is a fully managed NoSQL database.

```
DynamoDB vs RDS:
────────────────

┌────────────────────────────────────────────────────────────────────┐
│ Criteria          │ DynamoDB              │ RDS                    │
├───────────────────┼───────────────────────┼────────────────────────┤
│ Data Model        │ Key-Value, Document   │ Relational (tables)    │
│ Schema            │ Flexible              │ Fixed                  │
│ Scaling           │ Automatic, unlimited  │ Vertical + Read replica│
│ Transactions      │ Limited (25 items)    │ Full ACID              │
│ Queries           │ Primary/Secondary key │ Complex SQL joins      │
│ Latency           │ Single-digit ms       │ Variable               │
│ Cost Model        │ Per request + storage │ Instance hours         │
│ Use Case          │ High scale, simple    │ Complex queries        │
└───────────────────┴───────────────────────┴────────────────────────┘

Choose DynamoDB when:
• Need single-digit millisecond latency at any scale
• Data access patterns are known (key-based)
• Don't need complex joins
• Need automatic scaling

Choose RDS when:
• Need complex SQL queries and joins
• Data has relationships
• Need ACID transactions across many rows
• Using existing SQL-based application
```

**DynamoDB Key Concepts:**
- **Partition Key:** How data is distributed
- **Sort Key:** Range queries within partition
- **GSI (Global Secondary Index):** Query on non-key attributes
- **LSI (Local Secondary Index):** Alternative sort key

**Deep Dive - DynamoDB Internals:**

```
DynamoDB Partition Architecture:
────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                     DynamoDB Table                                   │
│                                                                      │
│  Partition Key: userId   Sort Key: timestamp                        │
│                                                                      │
│  ┌──────────────────┬──────────────────┬──────────────────┐        │
│  │   Partition A    │   Partition B    │   Partition C    │        │
│  │   Hash: 0-33%    │   Hash: 34-66%   │   Hash: 67-100%  │        │
│  │                  │                  │                  │        │
│  │  userId=U001     │  userId=U002     │  userId=U003     │        │
│  │  userId=U004     │  userId=U005     │  userId=U006     │        │
│  │  ...             │  ...             │  ...             │        │
│  │                  │                  │                  │        │
│  │  Capacity:       │  Capacity:       │  Capacity:       │        │
│  │  1000 WCU        │  1000 WCU        │  1000 WCU        │        │
│  │  3000 RCU        │  3000 RCU        │  3000 RCU        │        │
│  └──────────────────┴──────────────────┴──────────────────┘        │
│                                                                      │
│  Total: 3000 WCU, 9000 RCU (evenly distributed)                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Hot Partition Problem:**
```
BAD Design - Hot Key:
─────────────────────
Partition Key: date (e.g., "2024-01-15")
→ All today's traffic hits ONE partition!
→ Throttling despite available capacity

GOOD Design - Distributed:
──────────────────────────
Partition Key: date#shard (e.g., "2024-01-15#3")
→ Write to random shard (1-10)
→ Read from all shards with parallel queries
```

**Consistency Models:**

| Read Type | Consistency | Cost | Use Case |
|-----------|-------------|------|----------|
| **Eventually Consistent** | May see stale data | 0.5 RCU/4KB | Non-critical reads |
| **Strongly Consistent** | Always latest | 1 RCU/4KB | Financial data |
| **Transactional** | ACID guarantees | 2 RCU/4KB | Multi-item updates |

**Capacity Planning:**
```
Write Capacity Unit (WCU) = 1 write/sec for items up to 1KB
Read Capacity Unit (RCU)  = 1 strongly consistent read/sec for items up to 4KB
                         = 2 eventually consistent reads/sec for items up to 4KB

Example: 100 items/sec, 2KB each, strongly consistent
→ Writes: 100 × ceil(2/1) = 200 WCU
→ Reads: 100 × ceil(2/4) = 100 RCU
```

**DynamoDB Streams:**
- Captures item-level changes (INSERT, MODIFY, DELETE)
- 24-hour retention
- Use cases: Replication, analytics, triggers (Lambda)
- Stream record contains: Keys only, New image, Old image, or Both

---

**Q14: Explain SQS vs SNS vs EventBridge. When to use each?**

**A:**

```
Messaging Services Comparison:
──────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                              SQS                                     │
│                         (Queue - Pull)                               │
│                                                                      │
│  Producer ───► [Message Queue] ◄─── Consumer                        │
│               (Buffering, Retry)   (Polls for messages)             │
│                                                                      │
│  Use: Decoupling, load leveling, guaranteed delivery                │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                              SNS                                     │
│                        (Pub/Sub - Push)                              │
│                                                                      │
│                    ┌─── Subscriber 1 (Lambda)                       │
│  Publisher ───► [Topic] ─── Subscriber 2 (SQS)                      │
│                    └─── Subscriber 3 (HTTP)                          │
│                                                                      │
│  Use: Fan-out, push notifications, multiple consumers               │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                          EventBridge                                 │
│                      (Event Bus - Rules)                             │
│                                                                      │
│  Event ───► [Event Bus] ───► Rule 1 ───► Target (Lambda)            │
│             (Schema)    ───► Rule 2 ───► Target (SQS)               │
│             (Archive)   ───► Rule 3 ───► Target (API)               │
│                                                                      │
│  Use: Event-driven architecture, SaaS integration, scheduling       │
└─────────────────────────────────────────────────────────────────────┘

Decision Matrix:
────────────────
• Need queue with retry? ──────────────► SQS
• Need fan-out to many? ───────────────► SNS
• Need content-based routing? ─────────► EventBridge
• Need SaaS integration? ──────────────► EventBridge
• Need message scheduling? ────────────► EventBridge
```

**Deep Dive - SQS Internals:**

```
SQS Message Lifecycle:
──────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  Producer ──► [Send Message] ──► Queue ──► [Receive] ──► Consumer   │
│                    │                            │                    │
│                    ▼                            ▼                    │
│              Message becomes              Message becomes            │
│              "available"                  "in-flight"                │
│                                                │                     │
│                                    ┌───────────┴───────────┐        │
│                                    │                       │        │
│                                    ▼                       ▼        │
│                              [Delete]              [Visibility       │
│                              Message               Timeout Expires]  │
│                                    │                       │        │
│                                    ▼                       ▼        │
│                              Permanently           Returns to       │
│                              removed               queue (retry)     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**SQS Queue Types:**

| Feature | Standard Queue | FIFO Queue |
|---------|---------------|------------|
| **Throughput** | Unlimited | 300 msg/sec (batching: 3000) |
| **Ordering** | Best-effort | Strict FIFO |
| **Delivery** | At-least-once (duplicates possible) | Exactly-once |
| **Use Case** | High throughput, order not critical | Order matters, deduplication needed |

**Key SQS Concepts:**

```
Visibility Timeout (Default: 30s):
──────────────────────────────────
• Time message is hidden from other consumers after receive
• Should be > processing time
• If processing fails, message reappears in queue

Message Retention (Default: 4 days, Max: 14 days):
──────────────────────────────────────────────────
• How long message stays in queue if not deleted
• After expiry, message is permanently lost

Dead Letter Queue (DLQ):
────────────────────────
• Captures messages that fail processing multiple times
• maxReceiveCount: Number of retries before DLQ
• Essential for debugging and preventing poison pills

Long Polling (Recommended):
───────────────────────────
• WaitTimeSeconds: 1-20 seconds
• Reduces API calls and cost
• Returns as soon as message available or timeout
```

**Deep Dive - SNS Internals:**

```
SNS Delivery Mechanisms:
────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                          SNS Topic                                   │
│                              │                                       │
│         ┌────────────────────┼────────────────────┐                 │
│         │                    │                    │                 │
│         ▼                    ▼                    ▼                 │
│    [Lambda]              [SQS]               [HTTP/S]               │
│    Async invoke         Enqueue              POST request           │
│    Retry: 3x            (reliable)           Retry: 3x              │
│                                                                      │
│    [Email]              [SMS]                [Mobile Push]          │
│    Best effort          Throttled            Platform specific      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Message Filtering:
──────────────────
{
  "store": ["electronics"],
  "price": [{"numeric": [">", 100]}]
}
→ Only matching messages delivered to subscriber
→ Reduces Lambda invocations and costs
```

**Deep Dive - EventBridge:**

```
EventBridge Event Structure:
────────────────────────────
{
  "version": "0",
  "id": "12345-abcde",
  "detail-type": "Order Placed",
  "source": "com.myapp.orders",
  "account": "123456789012",
  "time": "2024-01-15T12:00:00Z",
  "region": "us-east-1",
  "detail": {
    "orderId": "ORD-123",
    "customerId": "CUST-456",
    "total": 99.99
  }
}

Rule Pattern Matching:
──────────────────────
{
  "source": ["com.myapp.orders"],
  "detail-type": ["Order Placed"],
  "detail": {
    "total": [{"numeric": [">", 50]}]
  }
}
→ Routes only orders > $50 to premium processing
```

**EventBridge vs SNS vs SQS - When to Use:**

| Scenario | Best Choice | Why |
|----------|-------------|-----|
| Simple fan-out | SNS | Lowest latency, simple setup |
| Need retry/DLQ | SNS + SQS | SNS fans out, SQS handles failures |
| Content-based routing | EventBridge | Powerful pattern matching |
| Cross-account events | EventBridge | Built-in cross-account support |
| SaaS integration | EventBridge | Native integrations (Shopify, Zendesk) |
| Scheduled jobs | EventBridge | Cron/rate expressions built-in |
| Exactly-once processing | SQS FIFO | Deduplication ID |

---

**Q15: What are the different ways to connect your data center to AWS?**

**A:**

```
Connectivity Options:
─────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  Option 1: Site-to-Site VPN                                         │
│  ─────────────────────────                                          │
│  On-Prem ══════╗                                                    │
│  VPN Device    ║ IPSec Tunnel        ┌─────────────┐               │
│                ║ (over Internet)     │  VPN        │               │
│                ╚═════════════════════│  Gateway    │───► VPC       │
│                                      └─────────────┘               │
│  • Setup: Hours                                                     │
│  • Cost: Low ($0.05/hr)                                            │
│  • Bandwidth: Up to 1.25 Gbps                                      │
│  • Reliability: Internet-dependent                                  │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Option 2: Direct Connect                                           │
│  ────────────────────────                                           │
│  On-Prem ═══════════════════════════════════════════════► VPC      │
│          Dedicated fiber (1/10/100 Gbps)                            │
│                                                                      │
│  • Setup: Weeks to months                                           │
│  • Cost: Higher (port fee + data transfer)                         │
│  • Bandwidth: 1-100 Gbps                                           │
│  • Reliability: High (dedicated)                                    │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Option 3: Direct Connect + VPN (encrypted)                         │
│  ──────────────────────────────────────────                         │
│  On-Prem ═══[Direct Connect]═══╗                                    │
│                                ║ IPSec over DX                      │
│                                ╚════════════════════════► VPC      │
│                                                                      │
│  • Best of both: High bandwidth + encryption                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Comparison:
──────────────────────────────────────────────────────────
Feature          │ VPN          │ Direct Connect
─────────────────┼──────────────┼───────────────────────
Setup Time       │ Hours        │ Weeks/Months
Cost             │ Low          │ Higher
Max Bandwidth    │ 1.25 Gbps    │ 100 Gbps
Latency          │ Variable     │ Consistent
Encryption       │ Built-in     │ Optional (add VPN)
──────────────────────────────────────────────────────────
```

---

**Q16: Explain how you would implement a CI/CD pipeline on AWS.**

**A:**

```
AWS CI/CD Pipeline:
───────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                         CodePipeline                                 │
│                                                                      │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐     │
│  │  Source  │───►│  Build   │───►│  Test    │───►│  Deploy  │     │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘     │
│       │               │               │               │             │
│       ▼               ▼               ▼               ▼             │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐     │
│  │CodeCommit│    │CodeBuild │    │CodeBuild │    │CodeDeploy│     │
│  │ GitHub   │    │          │    │          │    │  ECS     │     │
│  │ S3       │    │ • Compile│    │ • Unit   │    │  Lambda  │     │
│  │          │    │ • Docker │    │ • Integr │    │  EC2     │     │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

# buildspec.yml
version: 0.2
phases:
  install:
    runtime-versions:
      java: corretto17
  pre_build:
    commands:
      - echo "Running pre-build..."
      - mvn clean
  build:
    commands:
      - echo "Building..."
      - mvn package -DskipTests
  post_build:
    commands:
      - echo "Running tests..."
      - mvn test
artifacts:
  files:
    - target/*.jar
    - appspec.yml
```

---

**Q17: What is CloudFormation? Compare with Terraform.**

**A:** Both are Infrastructure as Code (IaC) tools.

| Feature | CloudFormation | Terraform |
|---------|---------------|-----------|
| **Provider** | AWS only | Multi-cloud |
| **Language** | JSON/YAML | HCL |
| **State** | Managed by AWS | Self-managed (S3+DynamoDB) |
| **Drift Detection** | Yes | Yes |
| **Rollback** | Automatic | Manual |
| **Learning Curve** | Lower (AWS native) | Moderate |
| **Community** | AWS Support | Large community, modules |

```yaml
# CloudFormation Example
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-unique-bucket-name
      VersioningConfiguration:
        Status: Enabled
```

```hcl
# Terraform Example
resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-unique-bucket-name"
}

resource "aws_s3_bucket_versioning" "versioning" {
  bucket = aws_s3_bucket.my_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}
```

---

**Q18: What is ECS? Compare with EKS.**

**A:**

```
Container Orchestration Options:
────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                              ECS                                     │
│                  (AWS Native Container Service)                      │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │                         ECS Cluster                              ││
│  │                                                                  ││
│  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           ││
│  │   │   Service   │  │   Service   │  │   Service   │           ││
│  │   │  (Task x 3) │  │  (Task x 2) │  │  (Task x 1) │           ││
│  │   └─────────────┘  └─────────────┘  └─────────────┘           ││
│  │                                                                  ││
│  │   Launch Types:                                                  ││
│  │   • Fargate (Serverless)                                        ││
│  │   • EC2 (Self-managed)                                          ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                              EKS                                     │
│                    (Managed Kubernetes)                              │
└─────────────────────────────────────────────────────────────────────┘

Comparison:
──────────────────────────────────────────────────────────
Feature          │ ECS                │ EKS
─────────────────┼────────────────────┼────────────────────
Complexity       │ Simple             │ Complex
Learning Curve   │ AWS-specific       │ Kubernetes (portable)
Ecosystem        │ AWS only           │ Large K8s ecosystem
Cost             │ Lower              │ $0.10/hr + resources
Portability      │ AWS locked         │ Multi-cloud
Best For         │ AWS-centric teams  │ K8s experience, multi-cloud
──────────────────────────────────────────────────────────
```

**Deep Dive - Container Architecture:**

```
ECS Terminology:
────────────────
• Cluster: Logical grouping of tasks/services
• Task Definition: Blueprint for containers (like Dockerfile + compose)
• Task: Running instance of task definition
• Service: Maintains desired count of tasks, integrates with ALB

EKS Terminology (Kubernetes):
─────────────────────────────
• Cluster: Control plane + worker nodes
• Pod: Smallest deployable unit (1+ containers)
• Deployment: Manages ReplicaSets and rolling updates
• Service: Exposes pods via stable endpoint
• Ingress: HTTP/S routing rules
```

**Fargate vs EC2 Launch Type:**

| Aspect | Fargate | EC2 |
|--------|---------|-----|
| **Management** | Serverless (no instances) | Manage EC2 instances |
| **Pricing** | Per vCPU + memory/second | EC2 instance pricing |
| **Scaling** | Task-level only | Instance + task level |
| **GPU Support** | No | Yes |
| **Spot Support** | Yes (Fargate Spot) | Yes |
| **Best For** | Variable workloads | Predictable, GPU, compliance |

**When to Choose:**
```
Choose ECS when:
────────────────
• Team is new to containers
• AWS-centric infrastructure
• Want simplicity over flexibility
• Using Fargate for serverless containers

Choose EKS when:
────────────────
• Team has Kubernetes experience
• Multi-cloud or hybrid strategy
• Need K8s ecosystem (Helm, Istio, ArgoCD)
• Existing K8s manifests/tools
```

---

**Q19: How do you secure data at rest and in transit on AWS?**

**A:**

```
Encryption Strategy:
────────────────────

At Rest:
─────────────────────────────────────────────────────────────────────

┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    S3       │     │    EBS      │     │    RDS      │
│             │     │             │     │             │
│  SSE-S3     │     │  KMS        │     │  TDE (Oracle)│
│  SSE-KMS    │     │  Encrypted  │     │  KMS        │
│  SSE-C      │     │  AMI        │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       └───────────────────┴───────────────────┘
                           │
                    ┌──────▼──────┐
                    │    KMS      │
                    │ (Key Mgmt)  │
                    │             │
                    │ • CMK       │
                    │ • Rotation  │
                    │ • Policies  │
                    └─────────────┘

In Transit:
─────────────────────────────────────────────────────────────────────

Client ───[HTTPS/TLS]───► ALB ───[TLS]───► EC2
                          │
                    ┌─────▼─────┐
                    │    ACM    │
                    │(Cert Mgmt)│
                    │           │
                    │ • Free    │
                    │ • Auto-   │
                    │   renew   │
                    └───────────┘

Best Practices:
• Enable default encryption on S3 buckets
• Use KMS for key management (not self-managed)
• Enforce TLS 1.2+ minimum
• Use ACM for certificate management
• Enable VPC endpoints for AWS services (stay private)
```

**Deep Dive - Encryption Theory:**

```
KMS Key Hierarchy:
──────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                         KMS                                          │
│                                                                      │
│   ┌───────────────────────────────────────────────────────────┐ │
│   │  Customer Master Key (CMK)                                      │ │
│   │  • Never leaves KMS                                            │ │
│   │  • Used to encrypt Data Keys                                   │ │
│   └────────────────────────────┬──────────────────────────────┘ │
│                                │                                     │
│                                ▼                                     │
│   ┌───────────────────────────────────────────────────────────┐ │
│   │  Data Key (DEK)                                                 │ │
│   │  • Generated per object/volume                                 │ │
│   │  • Stored encrypted alongside data                             │ │
│   │  • Used for actual encryption (AES-256)                        │ │
│   └───────────────────────────────────────────────────────────┘ │
│                                                                      │
│   This is "Envelope Encryption" - Key encrypts Key encrypts Data    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**S3 Encryption Options:**

| Option | Key Management | Who Encrypts | Use Case |
|--------|----------------|--------------|----------|
| **SSE-S3** | AWS managed | S3 | Simplest, no key mgmt |
| **SSE-KMS** | Customer managed (KMS) | S3 | Audit, key rotation control |
| **SSE-C** | Customer provided | S3 | Full key control |
| **Client-side** | Customer managed | Client | End-to-end encryption |

**TLS Best Practices:**

```
TLS Policy Selection (ALB/CloudFront):
──────────────────────────────────────

Strict (Recommended):    TLS 1.2+ only
                         Modern ciphers (AES-GCM)
                         No deprecated protocols

Compatible:              TLS 1.0+ supported
                         Older ciphers available
                         For legacy client support
```

**Zero Trust Security Model:**
- Encrypt everywhere (at rest AND in transit)
- Use VPC endpoints (traffic never leaves AWS network)
- Enable CloudTrail for audit trail
- Use AWS PrivateLink for service access

---

**Q20: What is the difference between Horizontal and Vertical scaling?**

**A:**

```
Scaling Strategies:
───────────────────

Vertical Scaling (Scale Up):
────────────────────────────
                    ┌───────────────┐
    Before:         │   m5.large    │
                    │   2 vCPU      │
                    │   8 GB RAM    │
                    └───────────────┘
                           │
                           ▼ Stop, resize, start
                    ┌───────────────┐
    After:          │   m5.4xlarge  │
                    │   16 vCPU     │
                    │   64 GB RAM   │
                    └───────────────┘

    Pros: Simple, no code changes
    Cons: Downtime, upper limit, single point of failure

Horizontal Scaling (Scale Out):
───────────────────────────────
                    ┌─────────────────────────┐
                    │    Load Balancer        │
                    └───────────┬─────────────┘
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
    ┌────▼────┐           ┌────▼────┐           ┌────▼────┐
    │  EC2    │           │  EC2    │           │  EC2    │
    │ m5.large│           │ m5.large│           │ m5.large│
    └─────────┘           └─────────┘           └─────────┘

    Pros: No downtime, fault tolerant, unlimited scale
    Cons: Application must be stateless, complex

AWS Services for Horizontal Scaling:
• Auto Scaling Groups (EC2)
• ECS/EKS Service Auto Scaling
• Lambda (automatic)
• DynamoDB (automatic)
• Aurora (read replicas)
```

---

## 3. Senior Level (5-8 Years)

### Advanced Architecture

**Q21: Design a multi-region disaster recovery strategy with RPO < 1 minute and RTO < 5 minutes.**

**A:**

```
Multi-Region Active-Active DR:
──────────────────────────────

                    ┌─────────────────────────┐
                    │      Route 53           │
                    │   (Latency/Failover)    │
                    │   Health Checks: 10s    │
                    └───────────┬─────────────┘
                                │
          ┌─────────────────────┴─────────────────────┐
          │                                           │
          ▼                                           ▼
┌─────────────────────────────────┐   ┌─────────────────────────────────┐
│     US-EAST-1 (Primary)          │   │     EU-WEST-1 (Secondary)       │
│                                  │   │                                  │
│  ┌────────────────────────────┐  │   │  ┌────────────────────────────┐  │
│  │   Application (EKS)       │  │   │  │   Application (EKS)       │  │
│  │   Full capacity           │  │   │  │   Full capacity           │  │
│  └────────────┬───────────────┘  │   │  └────────────┬───────────────┘  │
│               │                  │   │               │                  │
│  ┌────────────▼───────────────┐  │   │  ┌────────────▼───────────────┐  │
│  │   Aurora Global Database  │  │   │  │   Aurora Global Database  │  │
│  │   PRIMARY (Writer)        │◄─┼─<1s─┼──►│   SECONDARY (Reader)     │  │
│  │                           │  │ repl │  │   Promotable             │  │
│  └────────────────────────────┘  │   │  └────────────────────────────┘  │
│                                  │   │                                  │
│  ┌────────────────────────────┐  │   │  ┌────────────────────────────┐  │
│  │   DynamoDB Global Table   │◄─┼──async─┼──►│   DynamoDB Global Table │  │
│  └────────────────────────────┘  │   │  └────────────────────────────┘  │
│                                  │   │                                  │
│  ┌────────────────────────────┐  │   │  ┌────────────────────────────┐  │
│  │   ElastiCache Global      │◄─┼──async─┼──►│   ElastiCache Global    │  │
│  └────────────────────────────┘  │   │  └────────────────────────────┘  │
│                                  │   │                                  │
└──────────────────────────────────┘   └──────────────────────────────────┘

Failover Process (Automated):
────────────────────────────────────────────────────────────────

1. Route 53 Health Check Fails (10-30 seconds)
   └── Triggers automatic DNS failover

2. Aurora Failover (managed or detach/promote) (~1 minute)
   └── Secondary becomes primary writer
   └── Application reconnects (connection string update)

3. Traffic flows to secondary region
   └── RTO achieved: ~5 minutes total

Key Components:
• Aurora Global Database: <1 second replication lag
• DynamoDB Global Tables: Automatic multi-region replication
• S3 Cross-Region Replication: Object sync
• Route 53: DNS-based failover with health checks

Cost Consideration: 2x infrastructure cost
```

---

**Q22: How would you implement a zero-downtime deployment strategy?**

**A:**

```
Blue/Green Deployment with CodeDeploy:
──────────────────────────────────────

                         ALB
                          │
           ┌──────────────┴──────────────┐
           │                             │
    ┌──────▼──────┐              ┌──────▼──────┐
    │ Target Group│              │ Target Group│
    │   (Blue)    │              │   (Green)   │
    │   100%      │              │    0%       │
    └──────┬──────┘              └──────┬──────┘
           │                             │
    ┌──────▼──────┐              ┌──────▼──────┐
    │    ECS      │              │    ECS      │
    │  Version 1  │              │  Version 2  │
    │  (Current)  │              │  (New)      │
    └─────────────┘              └─────────────┘

Traffic Shift Options:
──────────────────────
• AllAtOnce: 0% → 100% instantly
• Linear10PercentEvery1Minute: 10% every minute
• Canary10Percent5Minutes: 10% for 5 min, then 100%

# Terraform - CodeDeploy Blue/Green
resource "aws_codedeploy_deployment_group" "main" {
  app_name               = aws_codedeploy_app.main.name
  deployment_group_name  = "production"
  service_role_arn       = aws_iam_role.codedeploy.arn
  deployment_config_name = "CodeDeployDefault.ECSCanary10Percent5Minutes"

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }
    
    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.main.name
    service_name = aws_ecs_service.main.name
  }

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM"]
  }

  alarm_configuration {
    alarms  = [aws_cloudwatch_metric_alarm.error_rate.name]
    enabled = true
  }
}

Rollback Triggers:
• CloudWatch Alarm (error rate > 5%)
• Health check failures
• Manual intervention
```

---

**Q23: Explain how you would design a multi-tenant SaaS architecture on AWS.**

**A:**

```
Multi-Tenant Architecture Options:
──────────────────────────────────

Option 1: Pool Model (Shared Everything)
────────────────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│                         Shared Resources                             │
│                                                                      │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │                    Shared Database                           │  │
│   │                                                              │  │
│   │  ┌────────────────────────────────────────────────────────┐ │  │
│   │  │ tenant_id │ data_column_1 │ data_column_2 │ ...        │ │  │
│   │  ├───────────┼───────────────┼───────────────┼────────────┤ │  │
│   │  │ tenant_A  │ ...           │ ...           │            │ │  │
│   │  │ tenant_B  │ ...           │ ...           │            │ │  │
│   │  │ tenant_C  │ ...           │ ...           │            │ │  │
│   │  └────────────────────────────────────────────────────────┘ │  │
│   │  Row-Level Security: WHERE tenant_id = :current_tenant     │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                                                                      │
│   Pros: Cost-efficient, simple operations                           │
│   Cons: Noisy neighbor, complex isolation                           │
└─────────────────────────────────────────────────────────────────────┘

Option 2: Silo Model (Dedicated Everything)
───────────────────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│   Tenant A Account        Tenant B Account        Tenant C Account  │
│   ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐│
│   │  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  ││
│   │  │    EKS    │  │    │  │    EKS    │  │    │  │    EKS    │  ││
│   │  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  ││
│   │  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  ││
│   │  │  Aurora   │  │    │  │  Aurora   │  │    │  │  Aurora   │  ││
│   │  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  ││
│   └─────────────────┘    └─────────────────┘    └─────────────────┘│
│                                                                      │
│   Pros: Full isolation, compliance                                  │
│   Cons: Expensive, complex operations                               │
└─────────────────────────────────────────────────────────────────────┘

Option 3: Bridge Model (Hybrid)
───────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │                  Shared Compute (EKS)                        │  │
│   │   Namespace: tenant-a    Namespace: tenant-b                │  │
│   │   ResourceQuota: 4 CPU   ResourceQuota: 8 CPU               │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                                                                      │
│   ┌─────────────────┐    ┌─────────────────┐                       │
│   │  Tenant A DB    │    │  Tenant B DB    │   (Dedicated DBs)    │
│   │  (schema/db)    │    │  (schema/db)    │                       │
│   └─────────────────┘    └─────────────────┘                       │
│                                                                      │
│   Pros: Balanced cost and isolation                                 │
│   Cons: Complex implementation                                      │
└─────────────────────────────────────────────────────────────────────┘

Tenant Isolation Techniques:
────────────────────────────
• API Gateway: Usage plans per tenant
• IAM: Tenant-scoped policies using tags
• Database: Row-level security, separate schemas
• Compute: Kubernetes namespaces with quotas
• Network: Security groups, NACLs
```

---

**Q24: How do you optimize costs while maintaining performance on AWS?**

**A:**

```
Cost Optimization Framework:
────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Cost Optimization Pillars                         │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────┐│
│  │ Right-sizing │  │ Purchasing   │  │ Architecture │  │ Visibility││
│  │              │  │ Options      │  │ Optimization │  │          ││
│  │ • Compute   │  │ • Reserved   │  │ • Serverless │  │ • Tags   ││
│  │   Optimizer  │  │ • Savings    │  │ • Caching    │  │ • Budgets││
│  │ • Trusted   │  │   Plans      │  │ • Data       │  │ • Cost   ││
│  │   Advisor    │  │ • Spot       │  │   transfer  │  │   Explorer││
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

1. Right-Sizing (Save 30-50%):
───────────────────────────────
# Check Compute Optimizer recommendations
aws compute-optimizer get-ec2-instance-recommendations \
  --filters "name=Finding,values=OVER_PROVISIONED"

Current: m5.2xlarge (CPU: 15%, Mem: 20%)
Recommended: m5.large (Save $200/month)

2. Purchasing Options:
────────────────────────────────────────────────────────────
                    │ Discount │ Flexibility │ Best For
────────────────────┼──────────┼─────────────┼─────────────
Reserved Instances  │ Up to 72%│ Low         │ Steady workloads
Savings Plans       │ Up to 66%│ High        │ Variable compute
Spot Instances      │ Up to 90%│ Highest     │ Fault-tolerant
────────────────────────────────────────────────────────────

3. Architecture Optimizations:
──────────────────────────────
• S3 Intelligent-Tiering: Auto-move infrequent data
• CloudFront: Reduce data transfer costs
• VPC Endpoints: Avoid NAT Gateway costs ($0.045/GB)
• Graviton (ARM): 20% better price/performance
• Lambda: Pay per invocation vs always-on

4. Quick Wins Checklist:
────────────────────────
□ Delete unattached EBS volumes
□ Release unused Elastic IPs ($3.60/month each)
□ Stop non-production instances after hours
□ Enable S3 lifecycle policies
□ Use CloudWatch Logs retention policies
□ Review cross-region data transfer
```

---

**Q25: Explain the AWS Well-Architected Framework.**

**A:**

```
AWS Well-Architected Framework:
───────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                     Six Pillars                                      │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 1. Operational Excellence                                      │  │
│  │    • Automate operations                                       │  │
│  │    • Make frequent, small changes                              │  │
│  │    • Anticipate failure                                        │  │
│  │    Tools: CloudWatch, Systems Manager, Config                  │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 2. Security                                                    │  │
│  │    • Implement strong identity foundation                      │  │
│  │    • Enable traceability                                       │  │
│  │    • Apply security at all layers                             │  │
│  │    Tools: IAM, GuardDuty, Security Hub, KMS                   │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 3. Reliability                                                 │  │
│  │    • Automatically recover from failure                        │  │
│  │    • Test recovery procedures                                  │  │
│  │    • Scale horizontally                                        │  │
│  │    Tools: Auto Scaling, Route 53, Multi-AZ, Backup            │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 4. Performance Efficiency                                      │  │
│  │    • Use serverless architectures                             │  │
│  │    • Go global in minutes                                      │  │
│  │    • Experiment more often                                     │  │
│  │    Tools: Lambda, CloudFront, Auto Scaling, Elasticache       │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 5. Cost Optimization                                           │  │
│  │    • Adopt consumption model                                   │  │
│  │    • Measure overall efficiency                                │  │
│  │    • Stop spending on data center operations                  │  │
│  │    Tools: Cost Explorer, Budgets, Savings Plans               │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 6. Sustainability (Newest)                                     │  │
│  │    • Understand your impact                                    │  │
│  │    • Maximize utilization                                      │  │
│  │    • Use managed services                                      │  │
│  │    Tools: Graviton, Spot, Right-sizing                        │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Well-Architected Review:
• Use AWS Well-Architected Tool
• Answer questions across all pillars
• Identify high-risk issues
• Create improvement plan
```

---

**Q26: How would you troubleshoot a Lambda function timeout issue?**

**A:**

```
Lambda Troubleshooting Checklist:
─────────────────────────────────

1. Check Configuration:
────────────────────────────────────────────────────────────────
aws lambda get-function-configuration \
  --function-name my-function \
  --query '{Timeout:Timeout,Memory:MemorySize,VPC:VpcConfig}'

Common Issues:
• Timeout too low (default 3s, max 900s)
• Memory too low (also affects CPU)
• VPC cold start (add provisioned concurrency)

2. Analyze X-Ray Traces:
────────────────────────────────────────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│ Request Timeline (X-Ray)                                             │
│                                                                      │
│ Lambda ──────────────────────────────────────────────────────────►  │
│   │                                                                  │
│   ├── Initialization (Cold Start): 500ms ◄── Problem area          │
│   │                                                                  │
│   ├── Handler Start: 200ms                                          │
│   │   │                                                              │
│   │   ├── DynamoDB GetItem: 50ms                                    │
│   │   ├── External API Call: 4500ms ◄── TIMEOUT HERE               │
│   │   └── S3 PutObject: (never reached)                            │
│   │                                                                  │
│   └── Total: 5250ms (Timeout: 5000ms)                               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

3. Solutions by Root Cause:
────────────────────────────────────────────────────────────────

┌─────────────────┬─────────────────────────────────────────────────┐
│ Root Cause      │ Solution                                         │
├─────────────────┼─────────────────────────────────────────────────┤
│ Cold Start      │ • Provisioned Concurrency                       │
│                 │ • Warm-up events (CloudWatch)                   │
│                 │ • Reduce package size                           │
├─────────────────┼─────────────────────────────────────────────────┤
│ VPC Latency     │ • Use VPC endpoints                             │
│                 │ • More ENIs (increase memory)                   │
│                 │ • Hyperplane (automatic now)                    │
├─────────────────┼─────────────────────────────────────────────────┤
│ External Calls  │ • Add timeouts to SDK clients                   │
│                 │ • Async processing (SQS)                        │
│                 │ • Circuit breaker pattern                       │
├─────────────────┼─────────────────────────────────────────────────┤
│ DB Connections  │ • Connection pooling (RDS Proxy)                │
│                 │ • Reuse connections across invocations          │
│                 │ • Use DynamoDB (connectionless)                 │
└─────────────────┴─────────────────────────────────────────────────┘

# Fix: Add timeout to external calls
import requests

response = requests.get(
    'https://api.example.com',
    timeout=2.0  # 2 second timeout
)
```

---

**Q27: How do you implement least privilege access in AWS?**

**A:**

```
Least Privilege Implementation:
───────────────────────────────

1. Start with Zero Permissions:
────────────────────────────────────────────────────────────────
# Bad: Allow all S3 actions on all buckets
{
  "Effect": "Allow",
  "Action": "s3:*",
  "Resource": "*"
}

# Good: Allow specific actions on specific bucket
{
  "Effect": "Allow",
  "Action": [
    "s3:GetObject",
    "s3:PutObject"
  ],
  "Resource": "arn:aws:s3:::my-app-bucket/uploads/*"
}

2. Use IAM Access Analyzer:
────────────────────────────────────────────────────────────────
• Identifies unused permissions
• Generates policy based on CloudTrail activity
• Validates policies against best practices

# Generate policy from activity
aws accessanalyzer start-policy-generation \
  --policy-generation-details '{
    "principalArn": "arn:aws:iam::123456789012:role/MyRole"
  }' \
  --cloud-trail-details '{
    "trailArn": "arn:aws:cloudtrail:us-east-1:123456789012:trail/main",
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-31T00:00:00Z"
  }'

3. Permission Boundaries:
────────────────────────────────────────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  Permission Boundary ─────────────────────────────────────────┐     │
│  (Maximum permissions)                                          │     │
│                                                                 │     │
│     ┌─────────────────────────────────────────────────────┐   │     │
│     │                                                      │   │     │
│     │   IAM Policy (Granted permissions)                  │   │     │
│     │                                                      │   │     │
│     │      ┌────────────────────────────────────────┐     │   │     │
│     │      │                                         │     │   │     │
│     │      │   Effective Permissions                │     │   │     │
│     │      │   (Intersection)                       │     │   │     │
│     │      │                                         │     │   │     │
│     │      └────────────────────────────────────────┘     │   │     │
│     │                                                      │   │     │
│     └─────────────────────────────────────────────────────┘   │     │
│                                                                 │     │
└─────────────────────────────────────────────────────────────────┘     │

4. Service Control Policies (SCPs):
────────────────────────────────────────────────────────────────
# Deny actions outside allowed regions
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Deny",
    "NotAction": [
      "iam:*",
      "organizations:*",
      "support:*"
    ],
    "Resource": "*",
    "Condition": {
      "StringNotEquals": {
        "aws:RequestedRegion": ["us-east-1", "us-west-2"]
      }
    }
  }]
}
```

---

## 4. Principal/Architect Level (8+ Years)

### Architecture Leadership

**Q28: How would you design an enterprise-wide multi-account strategy?**

**A:**

```
Enterprise Multi-Account Architecture:
──────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    AWS Organizations                                 │
│                                                                      │
│  Root (Management Account)                                          │
│  └── Don't deploy workloads here                                    │
│                                                                      │
│  ├── Security OU                                                    │
│  │   ├── Security-Tooling (GuardDuty, Security Hub admin)          │
│  │   ├── Log-Archive (CloudTrail, Config, VPC Flow Logs)           │
│  │   └── Forensics (Incident investigation)                        │
│  │                                                                   │
│  ├── Infrastructure OU                                              │
│  │   ├── Network-Hub (Transit Gateway, DNS, Direct Connect)        │
│  │   ├── Shared-Services (CI/CD, Artifact registry, AMIs)          │
│  │   └── Identity (Identity Center, directory services)            │
│  │                                                                   │
│  ├── Workloads OU                                                   │
│  │   ├── Production OU                                              │
│  │   │   ├── Prod-App-A                                            │
│  │   │   ├── Prod-App-B                                            │
│  │   │   └── Prod-Data                                             │
│  │   │                                                               │
│  │   ├── Non-Production OU                                          │
│  │   │   ├── Staging-App-A                                         │
│  │   │   └── Dev-App-A                                             │
│  │   │                                                               │
│  │   └── Data OU                                                    │
│  │       ├── Data-Lake                                              │
│  │       └── Analytics                                              │
│  │                                                                   │
│  ├── Sandbox OU                                                     │
│  │   └── Developer sandboxes (auto-cleanup)                        │
│  │                                                                   │
│  └── Suspended OU                                                   │
│      └── Quarantine accounts                                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Key Design Decisions:
─────────────────────

1. Account Vending:
   • Control Tower Account Factory
   • Terraform + Service Catalog
   • Automated guardrails via SCPs

2. Network Architecture:
   ┌─────────────────────────────────────────────────────────────────┐
   │                    Transit Gateway                               │
   │                         │                                        │
   │    ┌────────────────────┼────────────────────┐                  │
   │    │                    │                    │                  │
   │    ▼                    ▼                    ▼                  │
   │ ┌────────┐        ┌────────┐          ┌────────┐              │
   │ │ Shared │        │  Prod  │          │  Dev   │              │
   │ │Services│        │  VPC   │          │  VPC   │              │
   │ └────────┘        └────────┘          └────────┘              │
   │                                                                  │
   │ Route Tables: Segment traffic by environment                    │
   │ Inspection VPC: Centralized firewall (optional)                 │
   └─────────────────────────────────────────────────────────────────┘

3. Identity Federation:
   • AWS Identity Center (SSO)
   • Permission Sets mapped to AD groups
   • No long-term credentials

4. Guardrails (SCPs):
   • Deny root user actions
   • Enforce encryption
   • Region restrictions
   • Prevent leaving organization
```

---

**Q29: How would you approach migrating 500+ applications to AWS?**

**A:**

```
Large-Scale Migration Framework:
────────────────────────────────

Phase 1: Mobilize (Weeks 1-4)
─────────────────────────────────────────────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│  • Build migration team                                              │
│  • Establish migration factory                                       │
│  • Set up landing zone (Control Tower)                              │
│  • Define governance, security baseline                             │
│  • Create communication plan                                         │
└─────────────────────────────────────────────────────────────────────┘

Phase 2: Assess (Weeks 5-12)
─────────────────────────────────────────────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│  Discovery Tools:                                                    │
│  • AWS Application Discovery Service                                │
│  • Migration Hub                                                    │
│  • Partner tools (Cloudamize, CAST)                                 │
│                                                                      │
│  Application Portfolio:                                              │
│  ┌────────────────┬───────────────────┬─────────────────────────┐  │
│  │ Application    │ Complexity        │ Migration Strategy      │  │
│  ├────────────────┼───────────────────┼─────────────────────────┤  │
│  │ Legacy Java    │ High              │ Replatform (Containers) │  │
│  │ .NET Apps      │ Medium            │ Rehost (EC2)            │  │
│  │ Custom Apps    │ Medium            │ Rehost → Refactor       │  │
│  │ COTS           │ Low               │ Repurchase (SaaS)       │  │
│  │ Data Warehouse │ High              │ Rearchitect (Redshift)  │  │
│  │ End-of-Life    │ N/A               │ Retire                  │  │
│  └────────────────┴───────────────────┴─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘

Phase 3: Migrate (Waves)
─────────────────────────────────────────────────────────────────────
The 7 Rs of Migration:
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  1. Rehost (Lift & Shift)                                           │
│     └── AWS MGN for servers, DMS for databases                      │
│                                                                      │
│  2. Replatform (Lift, Tinker & Shift)                               │
│     └── Move to RDS, containerize                                   │
│                                                                      │
│  3. Repurchase                                                       │
│     └── Move to SaaS (Salesforce, Workday)                          │
│                                                                      │
│  4. Refactor/Rearchitect                                            │
│     └── Redesign for cloud-native                                   │
│                                                                      │
│  5. Relocate                                                         │
│     └── VMware Cloud on AWS                                         │
│                                                                      │
│  6. Retain                                                           │
│     └── Keep on-premises (hybrid)                                   │
│                                                                      │
│  7. Retire                                                           │
│     └── Decommission                                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Wave Planning:
─────────────────────────────────────────────────────────────────────
Wave 1 (Pilot): 5-10 simple apps
Wave 2-5: Increasing complexity
Wave 6+: Complex, critical apps

Migration Factory Metrics:
• Apps migrated per week: Target 20-50
• Total migration time: 12-18 months
• Cost savings: 30-50% typical
```

---

**Q30: Design an event-driven architecture for processing 1 million events per second.**

**A:**

```
High-Throughput Event Architecture:
───────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Ingestion Layer (1M+ events/sec)                  │
│                                                                      │
│  Option A: Kinesis Data Streams                                     │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │  • 500+ shards (2000 records/sec/shard)                         ││
│  │  • Enhanced fan-out consumers                                    ││
│  │  • On-demand mode for auto-scaling                              ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  Option B: Amazon MSK (Managed Kafka)                               │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │  • Multi-AZ cluster                                             ││
│  │  • Provisioned IOPS storage                                     ││
│  │  • Consumer groups for parallel processing                       ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Processing Layer                                  │
│                                                                      │
│  Real-time Processing:                                              │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │                                                                  ││
│  │   Kinesis Data Streams ───► Kinesis Data Analytics (Flink)     ││
│  │                              │                                   ││
│  │                              ├── Aggregations (per second)      ││
│  │                              ├── Anomaly detection              ││
│  │                              └── Enrichment                     ││
│  │                                                                  ││
│  │   OR                                                             ││
│  │                                                                  ││
│  │   MSK ───► EKS (Flink/Spark Streaming) ───► Results            ││
│  │                                                                  ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  Batch Processing:                                                  │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │   Kinesis Firehose ───► S3 ───► EMR/Glue ───► Redshift        ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Storage Layer                                     │
│                                                                      │
│  Hot Data:                  Warm Data:              Cold Data:      │
│  ┌──────────────┐          ┌──────────────┐       ┌──────────────┐ │
│  │ ElastiCache  │          │    S3        │       │  S3 Glacier  │ │
│  │ (Real-time)  │          │ (Analytics)  │       │  (Archive)   │ │
│  └──────────────┘          └──────────────┘       └──────────────┘ │
│  ┌──────────────┐          ┌──────────────┐                        │
│  │ DynamoDB     │          │  Redshift    │                        │
│  │ (Events)     │          │  (Reporting) │                        │
│  └──────────────┘          └──────────────┘                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Scaling Considerations:
──────────────────────────────────────────────────────────────────────

1. Kinesis:
   • On-Demand: Auto-scales to 200 MB/s ingress
   • Provisioned: Plan for 1 MB/s per shard
   • 500 shards = 500K records/sec × 1KB = 500 MB/s

2. Processing:
   • Kinesis Analytics: Auto-scales
   • EKS/Flink: Horizontal pod autoscaling
   • Checkpoint to S3 for exactly-once

3. DynamoDB:
   • On-Demand: Unlimited scale
   • DAX: Sub-millisecond reads
   • Global Tables: Multi-region

Cost Optimization:
• Kinesis: On-Demand vs Provisioned (predictable better with provisioned)
• EMR: Spot instances for processing
• S3: Intelligent-Tiering
```

---

## 5. Java + AWS Interview Questions

### Java SDK and Integration

**Q31: How do you integrate a Spring Boot application with AWS services?**

**A:**

```java
// 1. Add Spring Cloud AWS dependency
// pom.xml
<dependency>
    <groupId>io.awspring.cloud</groupId>
    <artifactId>spring-cloud-aws-starter</artifactId>
    <version>3.0.0</version>
</dependency>

// 2. Configure AWS credentials (application.yml)
spring:
  cloud:
    aws:
      credentials:
        instance-profile: true  # Use IAM role (recommended)
      region:
        static: us-east-1

// 3. S3 Integration
@Service
public class S3Service {
    
    private final S3Client s3Client;
    
    public S3Service(S3Client s3Client) {
        this.s3Client = s3Client;
    }
    
    public void uploadFile(String bucket, String key, byte[] content) {
        s3Client.putObject(
            PutObjectRequest.builder()
                .bucket(bucket)
                .key(key)
                .build(),
            RequestBody.fromBytes(content)
        );
    }
    
    public byte[] downloadFile(String bucket, String key) {
        ResponseInputStream<GetObjectResponse> response = 
            s3Client.getObject(
                GetObjectRequest.builder()
                    .bucket(bucket)
                    .key(key)
                    .build()
            );
        return response.readAllBytes();
    }
}

// 4. DynamoDB Integration
@DynamoDbBean
public class Order {
    private String orderId;
    private String customerId;
    private Instant createdAt;
    
    @DynamoDbPartitionKey
    public String getOrderId() { return orderId; }
    
    @DynamoDbSortKey
    public String getCustomerId() { return customerId; }
    
    // getters and setters
}

@Repository
public class OrderRepository {
    
    private final DynamoDbEnhancedClient enhancedClient;
    private final DynamoDbTable<Order> orderTable;
    
    public OrderRepository(DynamoDbEnhancedClient enhancedClient) {
        this.enhancedClient = enhancedClient;
        this.orderTable = enhancedClient.table("Orders", 
            TableSchema.fromBean(Order.class));
    }
    
    public void save(Order order) {
        orderTable.putItem(order);
    }
    
    public Order findById(String orderId, String customerId) {
        return orderTable.getItem(
            Key.builder()
                .partitionValue(orderId)
                .sortValue(customerId)
                .build()
        );
    }
}

// 5. SQS Integration
@Service
public class SqsService {
    
    private final SqsTemplate sqsTemplate;
    
    public void sendMessage(String queue, Object message) {
        sqsTemplate.send(queue, message);
    }
    
    @SqsListener("order-queue")
    public void handleMessage(Order order) {
        // Process order
        log.info("Received order: {}", order.getOrderId());
    }
}
```

**Deep Dive - AWS SDK Architecture:**

```
AWS SDK v2 Architecture:
────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│     Your Application                                                 │
│            │                                                         │
│            ▼                                                         │
│     ┌──────────────────────────────────────────────────────┐      │
│     │          Service Client (S3Client, DynamoDbClient)        │      │
│     └────────────────────────┴─────────────────────────────┘      │
│                           │                                         │
│                           ▼                                         │
│     ┌──────────────────────────────────────────────────────┐      │
│     │          SDK Core (Request/Response handling)             │      │
│     │                                                           │      │
│     │  ┌─────────────┐ ┌─────────────┐ ┌────────────────┐    │      │
│     │  │ Credentials │ │   Signer    │ │     Retry      │    │      │
│     │  │  Provider   │ │  (SigV4)    │ │    Strategy    │    │      │
│     │  └─────────────┘ └─────────────┘ └────────────────┘    │      │
│     └────────────────────────┴─────────────────────────────┘      │
│                           │                                         │
│                           ▼                                         │
│     ┌──────────────────────────────────────────────────────┐      │
│     │          HTTP Client (Apache, Netty, URL Connection)      │      │
│     └────────────────────────┴─────────────────────────────┘      │
│                           │                                         │
│                           ▼                                         │
│                       AWS Service                                   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Credential Provider Chain (Default Order):**

| Priority | Provider | Source |
|----------|----------|--------|
| 1 | Environment Variables | `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` |
| 2 | System Properties | `aws.accessKeyId`, `aws.secretAccessKey` |
| 3 | Web Identity Token | OIDC token (EKS) |
| 4 | Profile Credentials | `~/.aws/credentials` |
| 5 | Container Credentials | ECS task role |
| 6 | Instance Profile | EC2 instance role |

**Best Practices:**
- Never hardcode credentials in code
- Use IAM roles for EC2/ECS/Lambda (most secure)
- Use environment variables for local development
- Enable credential caching (default in SDK v2)

**SDK v2 vs SDK v1:**

| Feature | SDK v1 | SDK v2 |
|---------|--------|--------|
| **Package** | `com.amazonaws` | `software.amazon.awssdk` |
| **HTTP Client** | Apache only | Pluggable (Netty, Apache) |
| **Async Support** | Callback-based | CompletableFuture |
| **Performance** | Baseline | 20-30% faster |
| **Modularity** | Monolithic | Service-specific JARs |
| **Non-blocking I/O** | No | Yes (with Netty) |

---

**Q32: How do you handle AWS SDK pagination and retries in Java?**

**A:**

```java
// 1. Pagination with SDK v2
public List<S3Object> listAllObjects(String bucket) {
    ListObjectsV2Request request = ListObjectsV2Request.builder()
        .bucket(bucket)
        .maxKeys(1000)
        .build();
    
    List<S3Object> allObjects = new ArrayList<>();
    
    // Method 1: Manual pagination
    ListObjectsV2Response response;
    do {
        response = s3Client.listObjectsV2(request);
        allObjects.addAll(response.contents());
        
        request = request.toBuilder()
            .continuationToken(response.nextContinuationToken())
            .build();
    } while (response.isTruncated());
    
    return allObjects;
}

// Method 2: Paginator (Recommended)
public List<S3Object> listAllObjectsWithPaginator(String bucket) {
    ListObjectsV2Request request = ListObjectsV2Request.builder()
        .bucket(bucket)
        .build();
    
    return s3Client.listObjectsV2Paginator(request)
        .contents()
        .stream()
        .collect(Collectors.toList());
}

// 2. Custom Retry Configuration
@Configuration
public class AwsConfig {
    
    @Bean
    public S3Client s3Client() {
        return S3Client.builder()
            .region(Region.US_EAST_1)
            .overrideConfiguration(ClientOverrideConfiguration.builder()
                .retryPolicy(RetryPolicy.builder()
                    .numRetries(5)
                    .backoffStrategy(
                        FullJitterBackoffStrategy.builder()
                            .baseDelay(Duration.ofMillis(100))
                            .maxBackoffTime(Duration.ofSeconds(20))
                            .build()
                    )
                    .throttlingBackoffStrategy(
                        EqualJitterBackoffStrategy.builder()
                            .baseDelay(Duration.ofSeconds(1))
                            .maxBackoffTime(Duration.ofSeconds(20))
                            .build()
                    )
                    .retryCondition(RetryCondition.defaultRetryCondition())
                    .build()
                )
                .apiCallTimeout(Duration.ofSeconds(30))
                .apiCallAttemptTimeout(Duration.ofSeconds(10))
                .build()
            )
            .build();
    }
}

// 3. Handling Specific Exceptions
public void handleS3Errors(String bucket, String key) {
    try {
        s3Client.getObject(GetObjectRequest.builder()
            .bucket(bucket)
            .key(key)
            .build());
    } catch (NoSuchKeyException e) {
        log.warn("Object not found: {}/{}", bucket, key);
        throw new NotFoundException("File not found");
    } catch (S3Exception e) {
        if (e.statusCode() == 403) {
            log.error("Access denied to {}/{}", bucket, key);
            throw new AccessDeniedException("No permission");
        }
        throw e;
    } catch (SdkClientException e) {
        log.error("SDK error: {}", e.getMessage());
        throw new ServiceException("AWS service unavailable");
    }
}
```

---

**Q33: How do you implement async operations with AWS SDK in Java?**

**A:**

```java
// 1. Async S3 Client
@Configuration
public class AsyncAwsConfig {
    
    @Bean
    public S3AsyncClient s3AsyncClient() {
        return S3AsyncClient.builder()
            .region(Region.US_EAST_1)
            .build();
    }
}

@Service
public class AsyncS3Service {
    
    private final S3AsyncClient s3AsyncClient;
    
    // Single async operation
    public CompletableFuture<PutObjectResponse> uploadAsync(
            String bucket, String key, byte[] content) {
        return s3AsyncClient.putObject(
            PutObjectRequest.builder()
                .bucket(bucket)
                .key(key)
                .build(),
            AsyncRequestBody.fromBytes(content)
        );
    }
    
    // Multiple parallel operations
    public CompletableFuture<List<S3Object>> uploadMultipleAsync(
            String bucket, List<FileUpload> files) {
        
        List<CompletableFuture<PutObjectResponse>> futures = files.stream()
            .map(file -> uploadAsync(bucket, file.getKey(), file.getContent()))
            .collect(Collectors.toList());
        
        return CompletableFuture.allOf(
            futures.toArray(new CompletableFuture[0])
        ).thenApply(v -> 
            futures.stream()
                .map(CompletableFuture::join)
                .collect(Collectors.toList())
        );
    }
    
    // With timeout and error handling
    public void uploadWithTimeout(String bucket, String key, byte[] content) {
        try {
            PutObjectResponse response = uploadAsync(bucket, key, content)
                .orTimeout(30, TimeUnit.SECONDS)
                .exceptionally(ex -> {
                    log.error("Upload failed: {}", ex.getMessage());
                    throw new CompletionException(ex);
                })
                .join();
            
            log.info("Uploaded: ETag={}", response.eTag());
        } catch (CompletionException e) {
            if (e.getCause() instanceof TimeoutException) {
                throw new ServiceException("Upload timed out");
            }
            throw e;
        }
    }
}

// 2. DynamoDB Async with Reactive Streams
@Service
public class ReactiveDynamoDbService {
    
    private final DynamoDbAsyncClient dynamoDbAsyncClient;
    
    public Mono<Order> findOrderAsync(String orderId) {
        return Mono.fromFuture(() -> 
            dynamoDbAsyncClient.getItem(GetItemRequest.builder()
                .tableName("Orders")
                .key(Map.of("orderId", AttributeValue.builder()
                    .s(orderId)
                    .build()))
                .build())
        ).map(response -> mapToOrder(response.item()));
    }
    
    public Flux<Order> queryOrdersAsync(String customerId) {
        return Flux.from(
            dynamoDbAsyncClient.queryPaginator(QueryRequest.builder()
                .tableName("Orders")
                .indexName("customer-index")
                .keyConditionExpression("customerId = :cid")
                .expressionAttributeValues(Map.of(
                    ":cid", AttributeValue.builder().s(customerId).build()
                ))
                .build())
        ).flatMapIterable(QueryResponse::items)
         .map(this::mapToOrder);
    }
}
```

---

**Q34: How do you implement connection pooling for RDS in a Java application?**

**A:**

```java
// 1. HikariCP Configuration (Recommended)
@Configuration
public class DataSourceConfig {
    
    @Bean
    @ConfigurationProperties("spring.datasource.hikari")
    public HikariConfig hikariConfig() {
        return new HikariConfig();
    }
    
    @Bean
    public DataSource dataSource(HikariConfig hikariConfig) {
        return new HikariDataSource(hikariConfig);
    }
}

// application.yml
spring:
  datasource:
    url: jdbc:postgresql://mydb.cluster-xyz.us-east-1.rds.amazonaws.com:5432/mydb
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
    driver-class-name: org.postgresql.Driver
    hikari:
      pool-name: MyAppPool
      maximum-pool-size: 20        # Max connections
      minimum-idle: 5              # Min idle connections
      idle-timeout: 300000         # 5 minutes
      max-lifetime: 1800000        # 30 minutes
      connection-timeout: 30000    # 30 seconds
      leak-detection-threshold: 60000  # Log leak after 1 min
      validation-timeout: 5000
      connection-test-query: SELECT 1

// 2. IAM Database Authentication
@Configuration
public class IamAuthDataSourceConfig {
    
    @Value("${aws.rds.hostname}")
    private String hostname;
    
    @Value("${aws.rds.port}")
    private int port;
    
    @Value("${aws.rds.database}")
    private String database;
    
    @Value("${aws.rds.username}")
    private String username;
    
    @Bean
    public DataSource dataSource() {
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl(String.format(
            "jdbc:postgresql://%s:%d/%s", hostname, port, database));
        config.setUsername(username);
        config.setMaximumPoolSize(20);
        
        // Generate IAM auth token
        config.setDataSourceProperties(getIamProperties());
        
        return new HikariDataSource(config);
    }
    
    private Properties getIamProperties() {
        Properties props = new Properties();
        
        RdsUtilities rdsUtilities = RdsUtilities.builder()
            .region(Region.US_EAST_1)
            .build();
        
        String authToken = rdsUtilities.generateAuthenticationToken(
            GenerateAuthenticationTokenRequest.builder()
                .hostname(hostname)
                .port(port)
                .username(username)
                .build()
        );
        
        props.setProperty("password", authToken);
        props.setProperty("ssl", "true");
        props.setProperty("sslmode", "verify-full");
        
        return props;
    }
}

// 3. RDS Proxy with Spring Boot
// application.yml for RDS Proxy
spring:
  datasource:
    url: jdbc:postgresql://my-proxy.proxy-xyz.us-east-1.rds.amazonaws.com:5432/mydb
    hikari:
      maximum-pool-size: 50  # Can be higher with RDS Proxy
      # RDS Proxy handles connection pooling, so app pool can be larger
```

**Deep Dive - Connection Pooling Theory:**

```
Connection Pool Architecture:
─────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                       Application                                   │
│                                                                      │
│  Thread 1 ───┐                              ┌───► Connection 1       │
│  Thread 2 ───┤      ┌───────────────┐      ├───► Connection 2       │
│  Thread 3 ───┤      │   HikariCP    │      ├───► Connection 3  ───► RDS
│  Thread 4 ───┤─────►│   (Pool)      │─────►│    (Idle)           │
│  ...      ───┤      │               │      ├───► Connection 4       │
│  Thread N ───┘      └───────────────┘      └───► Connection 5       │
│                                                                      │
│  • Threads wait if all connections in use                           │
│  • Connection reused after release                                   │
│  • Idle connections cleaned up periodically                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Pool Sizing Formula:**
```
Optimal Pool Size = (Core Count * 2) + Effective Spindle Count

For most apps: 10-20 connections is sufficient
Formula for CPU-bound: Pool Size = CPU Cores + 1
Formula for I/O-bound: Pool Size = CPU Cores * (1 + Wait Time / Service Time)
```

**Why RDS Proxy?**

| Problem | Without Proxy | With RDS Proxy |
|---------|---------------|----------------|
| **Lambda + RDS** | Connection explosion | Multiplexes connections |
| **Failover** | Apps must reconnect | Transparent failover |
| **IAM Auth** | Token refresh complexity | Managed by proxy |
| **Idle Connections** | Waste DB resources | Proxy manages lifecycle |

**Connection Lifecycle:**

| Parameter | Purpose | Recommended |
|-----------|---------|-------------|
| `maximum-pool-size` | Max connections | 10-20 for web apps |
| `minimum-idle` | Warm connections ready | Same as max (HikariCP recommends) |
| `idle-timeout` | Close unused connections | 10-30 minutes |
| `max-lifetime` | Force connection refresh | 30 minutes (< DB timeout) |
| `connection-timeout` | Wait for connection | 30 seconds |

---

**Q35: How do you deploy a Java application to AWS Lambda?**

**A:**

```java
// 1. Lambda Handler with AWS SDK v2
public class OrderHandler implements RequestHandler<APIGatewayProxyRequestEvent, 
                                                     APIGatewayProxyResponseEvent> {
    
    private final ObjectMapper objectMapper = new ObjectMapper();
    private final DynamoDbClient dynamoDbClient;
    
    // Initialize outside handler for connection reuse
    public OrderHandler() {
        this.dynamoDbClient = DynamoDbClient.builder()
            .region(Region.US_EAST_1)
            .build();
    }
    
    @Override
    public APIGatewayProxyResponseEvent handleRequest(
            APIGatewayProxyRequestEvent event, Context context) {
        
        try {
            String httpMethod = event.getHttpMethod();
            
            if ("POST".equals(httpMethod)) {
                return createOrder(event, context);
            } else if ("GET".equals(httpMethod)) {
                return getOrder(event, context);
            }
            
            return response(405, "Method not allowed");
            
        } catch (Exception e) {
            context.getLogger().log("Error: " + e.getMessage());
            return response(500, "Internal error");
        }
    }
    
    private APIGatewayProxyResponseEvent createOrder(
            APIGatewayProxyRequestEvent event, Context context) throws Exception {
        
        Order order = objectMapper.readValue(event.getBody(), Order.class);
        order.setOrderId(UUID.randomUUID().toString());
        order.setCreatedAt(Instant.now().toString());
        
        dynamoDbClient.putItem(PutItemRequest.builder()
            .tableName(System.getenv("ORDERS_TABLE"))
            .item(Map.of(
                "orderId", AttributeValue.builder().s(order.getOrderId()).build(),
                "customerId", AttributeValue.builder().s(order.getCustomerId()).build(),
                "createdAt", AttributeValue.builder().s(order.getCreatedAt()).build()
            ))
            .build());
        
        return response(201, objectMapper.writeValueAsString(order));
    }
    
    private APIGatewayProxyResponseEvent response(int statusCode, String body) {
        return new APIGatewayProxyResponseEvent()
            .withStatusCode(statusCode)
            .withHeaders(Map.of("Content-Type", "application/json"))
            .withBody(body);
    }
}

// 2. Spring Cloud Function for Lambda
@SpringBootApplication
public class OrderApplication {
    
    @Bean
    public Function<Order, Order> processOrder(OrderService orderService) {
        return order -> {
            order.setProcessedAt(Instant.now());
            return orderService.save(order);
        };
    }
    
    @Bean
    public Consumer<SQSEvent> handleSqsEvent(OrderService orderService) {
        return event -> {
            event.getRecords().forEach(record -> {
                Order order = parseOrder(record.getBody());
                orderService.process(order);
            });
        };
    }
}

// 3. Best Practices for Java Lambda
/*
Optimization Tips:
─────────────────────────────────────────────────────────────────────

1. Reduce Cold Start:
   • Use Provisioned Concurrency for critical functions
   • Minimize dependencies (use AWS SDK v2 modules)
   • Use GraalVM native-image (experimental)
   • SnapStart for Java 11+ (up to 10x faster cold start)

2. Memory/CPU:
   • More memory = more CPU = faster
   • 1769 MB = 1 full vCPU
   • Test with different memory settings

3. Connection Management:
   • Initialize SDK clients outside handler
   • Reuse connections across invocations
   • Use RDS Proxy for database connections

4. Package Size:
   • Use maven-shade-plugin to create uber-jar
   • Exclude unnecessary dependencies
   • Consider layers for large dependencies
*/

// pom.xml optimizations
<build>
    <plugins>
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-shade-plugin</artifactId>
            <version>3.4.1</version>
            <configuration>
                <createDependencyReducedPom>false</createDependencyReducedPom>
                <minimizeJar>true</minimizeJar>
            </configuration>
            <executions>
                <execution>
                    <phase>package</phase>
                    <goals>
                        <goal>shade</goal>
                    </goals>
                </execution>
            </executions>
        </plugin>
    </plugins>
</build>
```

---

**Q36: How do you implement distributed tracing with X-Ray in a Java microservices application?**

**A:**

```java
// 1. Add X-Ray SDK Dependencies
// pom.xml
<dependency>
    <groupId>com.amazonaws</groupId>
    <artifactId>aws-xray-recorder-sdk-spring</artifactId>
    <version>2.14.0</version>
</dependency>
<dependency>
    <groupId>com.amazonaws</groupId>
    <artifactId>aws-xray-recorder-sdk-aws-sdk-v2</artifactId>
    <version>2.14.0</version>
</dependency>

// 2. X-Ray Configuration
@Configuration
public class XRayConfig {
    
    static {
        // Configure X-Ray sampling rules
        AWSXRayRecorderBuilder builder = AWSXRayRecorderBuilder.standard()
            .withSamplingStrategy(new LocalizedSamplingStrategy(
                XRayConfig.class.getResource("/sampling-rules.json")))
            .withSegmentListener(new SLF4JSegmentListener());
        
        AWSXRay.setGlobalRecorder(builder.build());
    }
    
    @Bean
    public Filter tracingFilter() {
        return new AWSXRayServletFilter("order-service");
    }
}

// sampling-rules.json
{
  "version": 2,
  "default": {
    "fixed_target": 1,
    "rate": 0.1
  },
  "rules": [
    {
      "description": "Health checks",
      "host": "*",
      "http_method": "GET",
      "url_path": "/health",
      "fixed_target": 0,
      "rate": 0
    }
  ]
}

// 3. Instrument HTTP Clients
@Configuration
public class HttpClientConfig {
    
    @Bean
    public RestTemplate restTemplate() {
        RestTemplate restTemplate = new RestTemplate();
        restTemplate.setInterceptors(List.of(
            new XRayClientHttpRequestInterceptor()
        ));
        return restTemplate;
    }
    
    @Bean
    public WebClient webClient() {
        return WebClient.builder()
            .filter(new TracingExchangeFilterFunction())
            .build();
    }
}

// 4. Instrument AWS SDK Clients
@Bean
public S3Client s3Client() {
    return S3Client.builder()
        .overrideConfiguration(ClientOverrideConfiguration.builder()
            .addExecutionInterceptor(new TracingInterceptor())
            .build())
        .build();
}

// 5. Custom Subsegments
@Service
public class OrderService {
    
    public Order processOrder(Order order) {
        // Create custom subsegment
        Subsegment subsegment = AWSXRay.beginSubsegment("ProcessOrder");
        
        try {
            subsegment.putAnnotation("orderId", order.getOrderId());
            subsegment.putMetadata("orderDetails", order);
            
            // Business logic
            validateOrder(order);
            calculateTotal(order);
            saveOrder(order);
            
            return order;
            
        } catch (Exception e) {
            subsegment.addException(e);
            throw e;
        } finally {
            AWSXRay.endSubsegment();
        }
    }
    
    @XRayEnabled  // Annotation-based tracing
    public void validateOrder(Order order) {
        // Automatically traced
    }
}

// 6. Propagate Trace Context in Async Operations
@Service
public class AsyncOrderService {
    
    private final ExecutorService executor;
    
    public AsyncOrderService() {
        // Wrap executor with X-Ray context propagation
        this.executor = new SegmentContextExecutors()
            .wrapExecutorService(Executors.newFixedThreadPool(10));
    }
    
    public CompletableFuture<Order> processAsync(Order order) {
        return CompletableFuture.supplyAsync(() -> {
            // Trace context is propagated
            return processOrder(order);
        }, executor);
    }
}

// 7. X-Ray with SQS
@SqsListener("orders-queue")
public void handleMessage(
        @Payload Order order,
        @Header("AWSTraceHeader") String traceHeader) {
    
    // Propagate trace from SQS message
    TraceHeader header = TraceHeader.fromString(traceHeader);
    Segment segment = AWSXRay.beginSegment("ProcessSqsMessage");
    segment.setTraceId(header.getRootTraceId());
    segment.setParentId(header.getParentId());
    
    try {
        processOrder(order);
    } finally {
        AWSXRay.endSegment();
    }
}
```

```
X-Ray Trace Visualization:
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│ Trace: 1-abc123-def456                         Duration: 245ms      │
│                                                                      │
│ ┌─[order-service]──────────────────────────────────────────────────┐│
│ │                                                                    ││
│ │  ┌─[ProcessOrder]───────────────────────────────────────────────┐ ││
│ │  │ orderId: ORD-123                                               │ ││
│ │  │                                                                 │ ││
│ │  │  ┌─[DynamoDB: GetItem]────────┐                               │ ││
│ │  │  │ Table: Orders              │                               │ ││
│ │  │  │ Duration: 12ms             │                               │ ││
│ │  │  └────────────────────────────┘                               │ ││
│ │  │                                                                 │ ││
│ │  │  ┌─[payment-service]────────────────────────────────────────┐ │ ││
│ │  │  │                                                            │ │ ││
│ │  │  │  ┌─[ProcessPayment]─────────────────────────────────────┐ │ │ ││
│ │  │  │  │  ┌─[Stripe API]───────────┐                          │ │ │ ││
│ │  │  │  │  │ Duration: 180ms        │                          │ │ │ ││
│ │  │  │  │  └────────────────────────┘                          │ │ │ ││
│ │  │  │  └──────────────────────────────────────────────────────┘ │ │ ││
│ │  │  └──────────────────────────────────────────────────────────┘ │ ││
│ │  │                                                                 │ ││
│ │  │  ┌─[DynamoDB: PutItem]────────┐                               │ ││
│ │  │  │ Table: Orders              │                               │ ││
│ │  │  │ Duration: 8ms              │                               │ ││
│ │  │  └────────────────────────────┘                               │ ││
│ │  └─────────────────────────────────────────────────────────────────┘ ││
│ └──────────────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────┘
```

---

**Q37: How do you handle secrets and configuration in a Java application on AWS?**

**A:**

```java
// 1. Secrets Manager Integration
@Configuration
public class SecretsConfig {
    
    private final SecretsManagerClient secretsManagerClient;
    
    public SecretsConfig() {
        this.secretsManagerClient = SecretsManagerClient.builder()
            .region(Region.US_EAST_1)
            .build();
    }
    
    @Bean
    public DatabaseCredentials databaseCredentials() {
        GetSecretValueResponse response = secretsManagerClient.getSecretValue(
            GetSecretValueRequest.builder()
                .secretId("prod/myapp/database")
                .build()
        );
        
        return new ObjectMapper().readValue(
            response.secretString(), 
            DatabaseCredentials.class
        );
    }
}

// 2. Spring Cloud AWS Secrets Manager (Auto-configuration)
// bootstrap.yml
spring:
  cloud:
    aws:
      secretsmanager:
        prefix: /secret
        default-context: application
        name: myapp
        profile-separator: _

// Secrets are automatically loaded as properties:
// /secret/myapp_prod -> spring.datasource.password

// 3. Parameter Store Integration
@Configuration
public class ParameterStoreConfig {
    
    private final SsmClient ssmClient;
    
    public ParameterStoreConfig() {
        this.ssmClient = SsmClient.builder()
            .region(Region.US_EAST_1)
            .build();
    }
    
    @Bean("appConfig")
    public Map<String, String> appConfig() {
        GetParametersByPathResponse response = ssmClient.getParametersByPath(
            GetParametersByPathRequest.builder()
                .path("/myapp/prod/")
                .recursive(true)
                .withDecryption(true)
                .build()
        );
        
        return response.parameters().stream()
            .collect(Collectors.toMap(
                p -> p.name().replace("/myapp/prod/", ""),
                Parameter::value
            ));
    }
}

// 4. Caching Secrets (Avoid API calls on every request)
@Component
public class CachedSecretsManager {
    
    private final SecretsManagerClient client;
    private final LoadingCache<String, String> cache;
    
    public CachedSecretsManager(SecretsManagerClient client) {
        this.client = client;
        this.cache = CacheBuilder.newBuilder()
            .expireAfterWrite(1, TimeUnit.HOURS)
            .build(new CacheLoader<>() {
                @Override
                public String load(String secretId) {
                    return fetchSecret(secretId);
                }
            });
    }
    
    public String getSecret(String secretId) {
        return cache.getUnchecked(secretId);
    }
    
    private String fetchSecret(String secretId) {
        return client.getSecretValue(
            GetSecretValueRequest.builder()
                .secretId(secretId)
                .build()
        ).secretString();
    }
    
    // Handle secret rotation via Lambda trigger
    @Scheduled(fixedRate = 300000) // 5 minutes
    public void refreshCacheIfNeeded() {
        cache.invalidateAll();
    }
}

// 5. Environment-specific Configuration
@Configuration
@Profile("prod")
public class ProdConfig {
    
    @Bean
    public DataSource dataSource(DatabaseCredentials creds) {
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl(creds.getJdbcUrl());
        config.setUsername(creds.getUsername());
        config.setPassword(creds.getPassword());
        config.setMaximumPoolSize(20);
        return new HikariDataSource(config);
    }
}

// 6. External Configuration with AppConfig
@Service
public class FeatureFlagService {
    
    private final AppConfigDataClient appConfigClient;
    private String configurationToken;
    
    public FeatureFlagService() {
        this.appConfigClient = AppConfigDataClient.builder().build();
        initializeSession();
    }
    
    private void initializeSession() {
        StartConfigurationSessionResponse response = 
            appConfigClient.startConfigurationSession(
                StartConfigurationSessionRequest.builder()
                    .applicationIdentifier("myapp")
                    .environmentIdentifier("prod")
                    .configurationProfileIdentifier("feature-flags")
                    .build()
            );
        this.configurationToken = response.initialConfigurationToken();
    }
    
    public FeatureFlags getFeatureFlags() {
        GetLatestConfigurationResponse response = 
            appConfigClient.getLatestConfiguration(
                GetLatestConfigurationRequest.builder()
                    .configurationToken(configurationToken)
                    .build()
            );
        
        this.configurationToken = response.nextPollConfigurationToken();
        
        return new ObjectMapper().readValue(
            response.configuration().asUtf8String(),
            FeatureFlags.class
        );
    }
}
```

---

## 6. Scenario-Based Questions

**Q38: Your application is experiencing slow database queries. How do you troubleshoot and fix it?**

**A:**

**Theory: Database Performance Optimization Pyramid**

```
                    ▲
                   /  \
                  / HA \ (Least Impact)
                 /______\
                /        \
               / Caching  \
              /____________\
             /              \
            /  Read Replicas \
           /__________________\
          /                    \
         /    Query Tuning      \
        /________________________\
       /                          \
      /      Index Optimization    \
     /______________________________\
    /                                \
   /         Schema Design           \ (Most Impact)
  /____________________________________\

Start from bottom (schema) for maximum impact!
```

**Query Performance Analysis Framework:**

| Layer | What to Check | Tools |
|-------|---------------|-------|
| **Application** | N+1 queries, unnecessary fetches | APM, X-Ray |
| **Connection** | Pool exhaustion, timeout | HikariCP metrics |
| **Query** | Execution plan, missing indexes | EXPLAIN ANALYZE |
| **Instance** | CPU, memory, I/O | CloudWatch, Performance Insights |
| **Storage** | IOPS limits, throughput | CloudWatch |

```
Database Performance Troubleshooting:
─────────────────────────────────────

Step 1: Identify the Problem
────────────────────────────────────────────────────────────────────

CloudWatch Metrics to Check:
• DatabaseConnections (hitting max?)
• CPUUtilization (> 80%?)
• ReadLatency / WriteLatency
• FreeableMemory (swapping?)
• DiskQueueDepth (I/O bottleneck?)

# Enable Performance Insights
aws rds modify-db-instance \
  --db-instance-identifier mydb \
  --enable-performance-insights \
  --performance-insights-retention-period 7

Step 2: Analyze Slow Queries
────────────────────────────────────────────────────────────────────

-- PostgreSQL: Check slow queries
SELECT query, calls, mean_time, total_time
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 10;

-- MySQL: Enable slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

Step 3: Common Fixes
────────────────────────────────────────────────────────────────────

┌─────────────────────┬───────────────────────────────────────────────┐
│ Problem             │ Solution                                       │
├─────────────────────┼───────────────────────────────────────────────┤
│ Missing indexes     │ Add indexes based on query patterns           │
│ Full table scans    │ Optimize WHERE clauses, add covering indexes │
│ Too many connections│ Use connection pooling (RDS Proxy, HikariCP)  │
│ CPU bound           │ Scale up instance, optimize queries           │
│ I/O bound           │ Use Provisioned IOPS, optimize queries        │
│ Memory pressure     │ Scale up, tune buffer pool                    │
│ Lock contention     │ Optimize transactions, read replicas          │
└─────────────────────┴───────────────────────────────────────────────┘

Step 4: Implement Caching
────────────────────────────────────────────────────────────────────

Application ───► ElastiCache ───► RDS
               (if miss)

# ElastiCache pattern
public User getUser(String userId) {
    String cacheKey = "user:" + userId;
    
    // Check cache first
    User cached = cache.get(cacheKey, User.class);
    if (cached != null) {
        return cached;
    }
    
    // Cache miss - query database
    User user = userRepository.findById(userId);
    cache.put(cacheKey, user, Duration.ofMinutes(15));
    return user;
}
```

**Caching Strategies Deep Dive:**

| Strategy | Description | Use When |
|----------|-------------|----------|
| **Cache-Aside (Lazy)** | App checks cache, then DB | Read-heavy, tolerates stale data |
| **Write-Through** | Write to cache and DB together | Data consistency critical |
| **Write-Behind** | Write to cache, async to DB | High write throughput |
| **Read-Through** | Cache fetches from DB on miss | Simplify app code |

**Cache Invalidation Patterns:**
```
TTL-Based:
──────────
• Set expiry time (e.g., 15 minutes)
• Simple but may serve stale data

Event-Based:
────────────
DB Update ──► EventBridge ──► Lambda ──► Invalidate Cache
• More complex but always fresh
• Use DynamoDB Streams or RDS Events
```

---

**Q39: You need to process 10 million records from S3 and load them into DynamoDB. How would you design this?**

**A:**

```
Batch Processing Architecture:
──────────────────────────────

Option 1: AWS Glue (Managed ETL)
────────────────────────────────────────────────────────────────────

┌──────────┐     ┌──────────┐     ┌──────────┐
│   S3     │────►│  Glue    │────►│ DynamoDB │
│ (Source) │     │  Job     │     │ (Target) │
└──────────┘     └──────────┘     └──────────┘

# Glue job (PySpark)
import boto3
from pyspark.context import SparkContext
from awsglue.context import GlueContext

sc = SparkContext()
glueContext = GlueContext(sc)

# Read from S3
datasource = glueContext.create_dynamic_frame.from_options(
    connection_type="s3",
    connection_options={"paths": ["s3://bucket/data/"]},
    format="json"
)

# Write to DynamoDB
glueContext.write_dynamic_frame.from_options(
    frame=datasource,
    connection_type="dynamodb",
    connection_options={"dynamodb.output.tableName": "MyTable"}
)

Option 2: Step Functions + Lambda (Serverless)
────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Step Functions Workflow                           │
│                                                                      │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐         │
│  │ List S3      │───►│ Map State    │───►│ Aggregate    │         │
│  │ Objects      │    │ (Parallel)   │    │ Results      │         │
│  └──────────────┘    └──────┬───────┘    └──────────────┘         │
│                             │                                        │
│                    ┌────────┴────────┐                              │
│                    │                 │                              │
│               ┌────▼────┐      ┌────▼────┐                        │
│               │ Lambda  │ ...  │ Lambda  │  (1000 concurrent)     │
│               │ Batch 1 │      │ Batch N │                        │
│               └─────────┘      └─────────┘                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

# Lambda function (Java)
public class BatchProcessor implements RequestHandler<S3Event, String> {
    
    private final DynamoDbClient dynamoDb;
    
    public BatchProcessor() {
        this.dynamoDb = DynamoDbClient.builder().build();
    }
    
    @Override
    public String handleRequest(S3Event event, Context context) {
        String bucket = event.getRecords().get(0).getS3().getBucket().getName();
        String key = event.getRecords().get(0).getS3().getObject().getKey();
        
        // Read from S3
        List<Map<String, AttributeValue>> items = readFromS3(bucket, key);
        
        // Batch write to DynamoDB (25 items per batch)
        Lists.partition(items, 25).forEach(batch -> {
            Map<String, List<WriteRequest>> requestItems = Map.of(
                "MyTable",
                batch.stream()
                    .map(item -> WriteRequest.builder()
                        .putRequest(PutRequest.builder()
                            .item(item)
                            .build())
                        .build())
                    .collect(Collectors.toList())
            );
            
            dynamoDb.batchWriteItem(BatchWriteItemRequest.builder()
                .requestItems(requestItems)
                .build());
        });
        
        return "Processed " + items.size() + " items";
    }
}

Performance Considerations:
────────────────────────────────────────────────────────────────────

• DynamoDB Write Capacity:
  - On-Demand: Unlimited (but costs more)
  - Provisioned: Plan for burst (scale before job)
  - 10M records / 25 per batch = 400K batch writes
  - At 1000 WCU = ~7 hours

• Optimization:
  - Use BatchWriteItem (25 items/request)
  - Enable DynamoDB auto-scaling
  - Use parallel workers (Step Functions Map state)
  - Exponential backoff for throttling
```

---

## 7. Quick Reference Cheat Sheet

### AWS Services Quick Reference

```
Compute:
────────────────────────────────────────────────────────────────────
EC2           Virtual servers
Lambda        Serverless functions (up to 15 min)
Fargate       Serverless containers
ECS           Container orchestration (AWS native)
EKS           Managed Kubernetes

Storage:
────────────────────────────────────────────────────────────────────
S3            Object storage (11 9s durability)
EBS           Block storage for EC2
EFS           Managed NFS
FSx           Managed Windows File Server / Lustre

Database:
────────────────────────────────────────────────────────────────────
RDS           Managed relational (MySQL, PostgreSQL, etc.)
Aurora        High-performance MySQL/PostgreSQL
DynamoDB      Managed NoSQL (key-value, document)
ElastiCache   Managed Redis/Memcached
DocumentDB    MongoDB compatible
Neptune       Graph database

Networking:
────────────────────────────────────────────────────────────────────
VPC           Virtual private cloud
ALB           Layer 7 load balancer
NLB           Layer 4 load balancer
Route 53      DNS + health checks
CloudFront    CDN
API Gateway   API management
Transit GW    VPC connectivity hub

Security:
────────────────────────────────────────────────────────────────────
IAM           Identity and access management
KMS           Key management
Secrets Mgr   Secrets storage with rotation
GuardDuty     Threat detection
Security Hub  Security posture management
WAF           Web application firewall

Integration:
────────────────────────────────────────────────────────────────────
SQS           Message queuing
SNS           Pub/sub messaging
EventBridge   Event bus
Step Functions Workflow orchestration
```

### Common Port Numbers

```
Port    │ Service
────────┼─────────────────────────
22      │ SSH
80      │ HTTP
443     │ HTTPS
3306    │ MySQL/Aurora MySQL
5432    │ PostgreSQL/Aurora PostgreSQL
6379    │ Redis
27017   │ DocumentDB (MongoDB)
```

### AWS CLI Quick Commands

```bash
# EC2
aws ec2 describe-instances --filters "Name=instance-state-name,Values=running"
aws ec2 start-instances --instance-ids i-xxxxx
aws ec2 stop-instances --instance-ids i-xxxxx

# S3
aws s3 ls s3://bucket-name/
aws s3 cp file.txt s3://bucket-name/
aws s3 sync ./local-dir s3://bucket-name/prefix/

# Lambda
aws lambda invoke --function-name my-func output.json
aws lambda update-function-code --function-name my-func --zip-file fileb://code.zip

# CloudWatch
aws logs tail /aws/lambda/my-func --follow
aws cloudwatch get-metric-statistics --namespace AWS/EC2 ...

# DynamoDB
aws dynamodb scan --table-name MyTable --max-items 10
aws dynamodb query --table-name MyTable --key-condition-expression "pk = :pk"
```

---

## Summary

This guide covers AWS interview questions across experience levels:

- **Junior (0-2 years):** Core services, basic concepts
- **Mid-Level (2-5 years):** Architecture design, multiple services
- **Senior (5-8 years):** Complex architectures, optimization
- **Principal (8+ years):** Enterprise strategy, migration
- **Java + AWS:** SDK integration, best practices

**Key Interview Tips:**
1. Always think about scalability, availability, and cost
2. Reference specific AWS services and their trade-offs
3. Use diagrams to explain architectures
4. Mention security at every layer
5. Discuss monitoring and operational excellence

---

*Companion document to the AWS Learning Guide series*
