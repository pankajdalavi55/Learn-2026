# AWS Foundations: An Architect's Perspective

> **Target Audience**: Senior Software Engineers transitioning to Platform Engineer / Architect roles
> **Prerequisites**: Strong programming background, basic understanding of distributed systems, familiarity with web applications

---

## Table of Contents

1. [AWS Global Infrastructure](#1-aws-global-infrastructure)
2. [AWS Identity and Access Management (IAM)](#2-aws-identity-and-access-management-iam)
3. [Virtual Private Cloud (VPC) - Network Foundation](#3-virtual-private-cloud-vpc---network-foundation)
4. [Core Compute Services](#4-core-compute-services)
5. [Storage Services Foundation](#5-storage-services-foundation)
6. [Database Services Overview](#6-database-services-overview)
7. [AWS Well-Architected Framework](#7-aws-well-architected-framework)
8. [Cost Management Fundamentals](#8-cost-management-fundamentals)
9. [Architecture Decision Framework](#9-architecture-decision-framework)
10. [Interview Questions & Scenarios](#10-interview-questions--scenarios)

---

## 1. AWS Global Infrastructure

### 1.1 Understanding the Hierarchy

```
AWS Global Infrastructure
├── Regions (33+ globally)
│   ├── Availability Zones (AZs) - 2-6 per region
│   │   └── Data Centers (1+ per AZ)
│   └── Local Zones (for low-latency edge)
├── Edge Locations (400+)
│   └── CloudFront, Route 53, AWS WAF
└── Wavelength Zones (5G edge)
```

### 1.2 Regions

**Definition**: A geographical area containing multiple isolated data center clusters (AZs).

**Key Characteristics**:
- Completely isolated from other regions
- Data doesn't replicate across regions unless explicitly configured
- Services and pricing vary by region

**Region Selection Criteria** (Architect's Checklist):

| Factor | Consideration |
|--------|---------------|
| **Compliance** | Data sovereignty laws (GDPR, HIPAA) |
| **Latency** | Proximity to end users |
| **Service Availability** | Not all services available in all regions |
| **Cost** | Pricing varies by region (us-east-1 typically cheapest) |
| **Disaster Recovery** | Secondary region for DR strategy |

```python
# Region naming convention
# Format: {geography}-{direction}-{number}
# Examples:
us-east-1      # N. Virginia (oldest, most services)
us-west-2      # Oregon
eu-west-1      # Ireland
ap-south-1     # Mumbai
```

### 1.3 Availability Zones (AZs)

**Definition**: One or more discrete data centers with redundant power, networking, and connectivity.

**Key Characteristics**:
- Physically separated (miles apart)
- Connected via low-latency private fiber
- Independent failure domains
- Synchronous replication feasible (< 2ms latency)

**Architect's Mental Model**:

```
┌─────────────────── Region (us-east-1) ───────────────────┐
│                                                           │
│  ┌─────────┐     ┌─────────┐     ┌─────────┐            │
│  │  AZ-a   │     │  AZ-b   │     │  AZ-c   │            │
│  │         │     │         │     │         │            │
│  │ ┌─────┐ │     │ ┌─────┐ │     │ ┌─────┐ │            │
│  │ │ DC1 │ │     │ │ DC1 │ │     │ │ DC1 │ │            │
│  │ └─────┘ │     │ └─────┘ │     │ └─────┘ │            │
│  │ ┌─────┐ │     │ ┌─────┐ │     │         │            │
│  │ │ DC2 │ │     │ │ DC2 │ │     │         │            │
│  │ └─────┘ │     │ └─────┘ │     │         │            │
│  └────┬────┘     └────┬────┘     └────┬────┘            │
│       │               │               │                  │
│       └───────── Private Fiber ───────┘                  │
│              (< 2ms latency)                             │
└──────────────────────────────────────────────────────────┘
```

### 1.4 Edge Locations & Points of Presence

**Purpose**: Content delivery and DNS resolution at the edge

**Services Using Edge Locations**:
- **CloudFront**: CDN for static/dynamic content
- **Route 53**: DNS service
- **AWS WAF**: Web Application Firewall
- **Lambda@Edge**: Serverless at edge
- **AWS Global Accelerator**: Network layer optimization

### 1.5 Architect's Decision Matrix - Region & AZ

```
Question: How many AZs should I use?

Development/Test    → 1 AZ (cost optimization)
Production (99.9%)  → 2 AZs (standard HA)
Mission Critical    → 3+ AZs (maximum resilience)
Global Users        → Multi-Region with CloudFront
```

---

## 2. AWS Identity and Access Management (IAM)

### 2.1 IAM Core Components

```
IAM Hierarchy
├── AWS Account (Root User)
├── IAM Users (Human identities)
├── IAM Groups (Collection of users)
├── IAM Roles (Assumed identities)
├── IAM Policies (Permission documents)
└── Identity Providers (Federation)
```

### 2.2 The Principal Types

| Principal | Use Case | Best Practice |
|-----------|----------|---------------|
| **Root User** | Initial setup, billing | MFA, never use for daily tasks |
| **IAM User** | Human access with long-term credentials | Use for break-glass scenarios only |
| **IAM Role** | AWS services, cross-account, federation | Preferred for all automation |
| **Federated User** | Enterprise SSO | Use Identity Center (SSO) |

### 2.3 IAM Policies Deep Dive

**Policy Structure**:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowS3ReadAccess",
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::my-bucket",
                "arn:aws:s3:::my-bucket/*"
            ],
            "Condition": {
                "IpAddress": {
                    "aws:SourceIp": "192.168.1.0/24"
                }
            }
        }
    ]
}
```

**Policy Evaluation Logic**:

```
┌─────────────────────────────────────────────────────────┐
│                   Policy Evaluation                      │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  1. Explicit Deny?  ─── YES ──→ DENY                    │
│         │                                                │
│        NO                                                │
│         ↓                                                │
│  2. SCP Allow? (if Org) ─── NO ──→ DENY                 │
│         │                                                │
│        YES                                               │
│         ↓                                                │
│  3. Resource Policy Allow? ─── YES ──→ ALLOW            │
│         │                                                │
│        NO                                                │
│         ↓                                                │
│  4. Identity Policy Allow? ─── NO ──→ DENY              │
│         │                                                │
│        YES                                               │
│         ↓                                                │
│  5. Permission Boundary Allow? ─── NO ──→ DENY          │
│         │                                                │
│        YES                                               │
│         ↓                                                │
│      ALLOW                                               │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 2.4 IAM Roles for Services

**Common Role Trust Policies**:

```json
// EC2 Instance Role
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```

```json
// Cross-Account Role
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::123456789012:root"
            },
            "Action": "sts:AssumeRole",
            "Condition": {
                "StringEquals": {
                    "sts:ExternalId": "unique-external-id"
                }
            }
        }
    ]
}
```

### 2.5 IAM Best Practices for Architects

```
Security Principle          │ Implementation
───────────────────────────┼──────────────────────────────────────
Least Privilege            │ Start with zero permissions, add as needed
Use Roles, Not Users       │ Prefer temporary credentials
Enable MFA                 │ Enforce MFA for console and API (sensitive)
Rotate Credentials         │ Automate key rotation
Use Policy Conditions      │ IP restrictions, MFA requirements
Audit Regularly            │ IAM Access Analyzer, CloudTrail
Separate Environments      │ Different accounts for dev/staging/prod
Use Permission Boundaries  │ Delegate admin with guardrails
```

### 2.6 AWS Organizations & SCPs

**Service Control Policies (SCPs)** - Guardrails at the organization level:

```json
// Prevent leaving organization
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PreventLeaving",
            "Effect": "Deny",
            "Action": "organizations:LeaveOrganization",
            "Resource": "*"
        }
    ]
}
```

```json
// Restrict to specific regions
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "DenyOutsideRegions",
            "Effect": "Deny",
            "NotAction": [
                "iam:*",
                "organizations:*",
                "support:*"
            ],
            "Resource": "*",
            "Condition": {
                "StringNotEquals": {
                    "aws:RequestedRegion": [
                        "us-east-1",
                        "eu-west-1"
                    ]
                }
            }
        }
    ]
}
```

---

## 3. Virtual Private Cloud (VPC) - Network Foundation

### 3.1 VPC Core Concepts

**VPC** = Your isolated network in AWS cloud

```
VPC Architecture Overview
─────────────────────────────────────────────────────────────
│                     VPC (10.0.0.0/16)                      │
│                     65,536 IP addresses                     │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Availability Zone A                      │   │
│  │  ┌─────────────────┐  ┌─────────────────┐           │   │
│  │  │ Public Subnet   │  │ Private Subnet  │           │   │
│  │  │ 10.0.1.0/24     │  │ 10.0.10.0/24    │           │   │
│  │  │ (251 usable)    │  │ (251 usable)    │           │   │
│  │  └─────────────────┘  └─────────────────┘           │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Availability Zone B                      │   │
│  │  ┌─────────────────┐  ┌─────────────────┐           │   │
│  │  │ Public Subnet   │  │ Private Subnet  │           │   │
│  │  │ 10.0.2.0/24     │  │ 10.0.20.0/24    │           │   │
│  │  └─────────────────┘  └─────────────────┘           │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
─────────────────────────────────────────────────────────────
```

### 3.2 CIDR Notation Quick Reference

| CIDR | Subnet Mask | Available IPs | AWS Usable* |
|------|-------------|---------------|-------------|
| /16 | 255.255.0.0 | 65,536 | 65,531 |
| /20 | 255.255.240.0 | 4,096 | 4,091 |
| /24 | 255.255.255.0 | 256 | 251 |
| /28 | 255.255.255.240 | 16 | 11 |

*AWS reserves 5 IPs per subnet: Network, VPC Router, DNS, Future, Broadcast

### 3.3 Subnet Types and Routing

**Public Subnet**:
- Has route to Internet Gateway (IGW)
- Resources can have public IPs
- Use for: Load balancers, Bastion hosts, NAT Gateways

**Private Subnet**:
- No direct internet route
- Outbound internet via NAT Gateway (if needed)
- Use for: Application servers, Databases, Internal services

**Route Table Configuration**:

```
Public Subnet Route Table:
┌──────────────────┬─────────────────┬───────────────┐
│   Destination    │     Target      │    Status     │
├──────────────────┼─────────────────┼───────────────┤
│   10.0.0.0/16    │     local       │    Active     │
│   0.0.0.0/0      │     igw-xxx     │    Active     │
└──────────────────┴─────────────────┴───────────────┘

Private Subnet Route Table:
┌──────────────────┬─────────────────┬───────────────┐
│   Destination    │     Target      │    Status     │
├──────────────────┼─────────────────┼───────────────┤
│   10.0.0.0/16    │     local       │    Active     │
│   0.0.0.0/0      │     nat-xxx     │    Active     │
└──────────────────┴─────────────────┴───────────────┘
```

### 3.4 Security: NACLs vs Security Groups

```
                    NACL                    Security Group
                    ────                    ──────────────
Level:              Subnet                  Instance (ENI)
State:              Stateless               Stateful
Rules:              Allow & Deny            Allow only
Evaluation:         Sequential (by #)       All rules evaluated
Default:            Allow all               Deny all inbound
Return Traffic:     Must be explicit        Automatic
```

**Security Group Best Practices**:

```
┌─────────────────────────────────────────────────────────┐
│                 Security Group Design                    │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Web-SG (ALB)                                           │
│  ├─ Inbound: 443 from 0.0.0.0/0                        │
│  └─ Outbound: All to App-SG                            │
│                                                          │
│  App-SG (EC2/ECS)                                       │
│  ├─ Inbound: 8080 from Web-SG                          │
│  └─ Outbound: 5432 to DB-SG                            │
│                                                          │
│  DB-SG (RDS)                                            │
│  ├─ Inbound: 5432 from App-SG                          │
│  └─ Outbound: None (stateful return)                   │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3.5 VPC Connectivity Options

```
┌────────────────────────────────────────────────────────────────┐
│                    VPC Connectivity Matrix                      │
├─────────────────────┬──────────────────────────────────────────┤
│ Connection Type     │ Use Case                                  │
├─────────────────────┼──────────────────────────────────────────┤
│ Internet Gateway    │ Public internet access (bidirectional)   │
│ NAT Gateway         │ Outbound internet for private subnets    │
│ VPC Peering         │ Connect 2 VPCs (same/different account)  │
│ Transit Gateway     │ Hub-spoke for multiple VPCs              │
│ VPN Connection      │ Encrypted tunnel to on-premise           │
│ Direct Connect      │ Dedicated physical connection            │
│ VPC Endpoints       │ Private access to AWS services           │
│ PrivateLink         │ Private access to 3rd party services     │
└─────────────────────┴──────────────────────────────────────────┘
```

### 3.6 VPC Endpoints Deep Dive

**Gateway Endpoints** (Free):
- S3 and DynamoDB only
- Route table entry
- Regional (same region only)

**Interface Endpoints** (PrivateLink):
- Most AWS services
- ENI in your subnet (costs $)
- Private DNS supported

```
┌─────────────────── VPC ───────────────────┐
│                                            │
│  ┌──────────────────────────────────────┐ │
│  │ Private Subnet                        │ │
│  │                                       │ │
│  │  ┌─────────┐    ┌─────────────────┐  │ │
│  │  │   EC2   │───→│ VPC Endpoint    │  │ │
│  │  └─────────┘    │ (Interface)     │  │ │
│  │                 │ Private IP      │  │ │
│  │                 └────────┬────────┘  │ │
│  └──────────────────────────│───────────┘ │
│                             │              │
└─────────────────────────────│──────────────┘
                              │ AWS PrivateLink
                              ↓
                    ┌─────────────────┐
                    │   AWS Service   │
                    │   (e.g., SQS)   │
                    └─────────────────┘
```

---

## 4. Core Compute Services

### 4.1 Compute Options Overview

```
Compute Spectrum (More Control ←→ Less Control)
──────────────────────────────────────────────────────────────
EC2          │  ECS/EKS      │  Fargate       │  Lambda
─────────────┼───────────────┼────────────────┼──────────────
Full VM      │  Container    │  Serverless    │  Serverless
control      │  orchestration│  containers    │  functions
             │               │                │
Manage:      │  Manage:      │  Manage:       │  Manage:
- OS         │  - Container  │  - Container   │  - Code only
- Patching   │  - Infra*     │    image       │
- Scaling    │               │                │
```

### 4.2 EC2 - Elastic Compute Cloud

**Instance Families**:

| Family | Optimized For | Use Cases |
|--------|---------------|-----------|
| **M** (General) | Balanced | Web servers, small DBs |
| **C** (Compute) | CPU | Batch processing, ML inference |
| **R** (Memory) | RAM | In-memory caching, analytics |
| **I** (Storage) | Local NVMe | NoSQL DBs, data warehousing |
| **G/P** (GPU) | Graphics/ML | Deep learning, video encoding |
| **T** (Burstable) | Variable workloads | Dev/test, low-traffic sites |

**Instance Naming Convention**:
```
m5.xlarge
│ │  │
│ │  └── Size (nano → metal)
│ └───── Generation (higher = newer)
└─────── Family
```

**Sizing Quick Reference**:

| Size | vCPU | Memory | Network |
|------|------|--------|---------|
| nano | 2 | 0.5 GB | Low |
| micro | 2 | 1 GB | Low |
| small | 2 | 2 GB | Low-Mod |
| medium | 2 | 4 GB | Moderate |
| large | 2 | 8 GB | Moderate |
| xlarge | 4 | 16 GB | High |
| 2xlarge | 8 | 32 GB | High |
| 4xlarge | 16 | 64 GB | Very High |

### 4.3 EC2 Purchase Options

```
┌────────────────────────────────────────────────────────────────┐
│                    EC2 Pricing Models                           │
├─────────────────┬──────────────┬───────────────────────────────┤
│    Model        │   Discount   │         Best For               │
├─────────────────┼──────────────┼───────────────────────────────┤
│ On-Demand       │     0%       │ Short-term, unpredictable      │
│ Reserved (1yr)  │   ~40%       │ Steady-state workloads         │
│ Reserved (3yr)  │   ~60%       │ Long-term commitments          │
│ Savings Plans   │   ~40-60%    │ Flexible commitment ($/hr)     │
│ Spot Instances  │   ~90%       │ Fault-tolerant, flexible       │
│ Dedicated Host  │   Premium    │ Licensing, compliance          │
└─────────────────┴──────────────┴───────────────────────────────┘
```

**Spot Instance Architecture Pattern**:

```
┌─────────────────────────────────────────────────────────────┐
│                 Spot Instance Strategy                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   Auto Scaling Group (Mixed Instances)                       │
│   ├── On-Demand Base: 2 instances (guaranteed)              │
│   ├── Spot Instances: 8 instances (cost savings)            │
│   └── Allocation Strategy: capacity-optimized               │
│                                                              │
│   Instance Diversification:                                  │
│   ├── m5.xlarge (primary)                                   │
│   ├── m5a.xlarge (alternative)                              │
│   ├── m4.xlarge (alternative)                               │
│   └── m5d.xlarge (alternative)                              │
│                                                              │
│   Interruption Handling:                                     │
│   ├── 2-minute warning via instance metadata                │
│   ├── Graceful shutdown via lifecycle hooks                 │
│   └── State saved to S3/EFS                                 │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 4.4 EC2 Placement Groups

| Type | Use Case | Constraints |
|------|----------|-------------|
| **Cluster** | HPC, Low latency | Same AZ, same rack |
| **Spread** | Critical instances | Max 7 per AZ |
| **Partition** | Distributed systems (HDFS, Kafka) | Large scale, failure isolation |

### 4.5 Auto Scaling

**Components**:

```
Auto Scaling Architecture
─────────────────────────────────────────────────────────────
                    ┌───────────────────┐
                    │  Launch Template  │
                    │  (AMI, Type, SG)  │
                    └─────────┬─────────┘
                              │
                    ┌─────────▼─────────┐
                    │  Auto Scaling     │
                    │  Group (ASG)      │
                    │  Min: 2, Max: 10  │
                    │  Desired: 4       │
                    └─────────┬─────────┘
                              │
           ┌──────────────────┼──────────────────┐
           │                  │                  │
    ┌──────▼──────┐    ┌──────▼──────┐    ┌──────▼──────┐
    │ Scaling     │    │ Scaling     │    │ Scaling     │
    │ Policy:     │    │ Policy:     │    │ Policy:     │
    │ Target      │    │ Step        │    │ Scheduled   │
    │ Tracking    │    │ Scaling     │    │             │
    └─────────────┘    └─────────────┘    └─────────────┘
```

**Scaling Policy Types**:

```python
# Target Tracking (Recommended - Simple)
{
    "PolicyType": "TargetTrackingScaling",
    "TargetTrackingConfiguration": {
        "PredefinedMetricSpecification": {
            "PredefinedMetricType": "ASGAverageCPUUtilization"
        },
        "TargetValue": 70.0  # Keep CPU at 70%
    }
}

# Step Scaling (More Control)
# CPU > 80% → Add 2 instances
# CPU > 90% → Add 4 instances
# CPU < 40% → Remove 1 instance
```

### 4.6 Elastic Load Balancing (ELB)

**Types Comparison**:

| Feature | ALB (Layer 7) | NLB (Layer 4) | GLB (Layer 3) |
|---------|---------------|---------------|---------------|
| Protocol | HTTP/HTTPS | TCP/UDP/TLS | IP packets |
| Routing | Path, Host, Header | Port | GENEVE |
| Latency | ~400ms added | ~100μs added | Varies |
| Static IP | No (use Global Accelerator) | Yes | Yes |
| WebSocket | Yes | Yes | No |
| Use Case | Web apps, microservices | Gaming, IoT, extreme perf | 3rd party appliances |

**ALB Advanced Routing**:

```yaml
# Path-based routing
/api/*          → API Target Group
/static/*       → Static Content Target Group
/                → Default Target Group

# Host-based routing
api.example.com → API Target Group
www.example.com → Web Target Group

# Header-based routing
X-Custom-Header: mobile  → Mobile Target Group
X-Custom-Header: desktop → Desktop Target Group
```

---

## 5. Storage Services Foundation

### 5.1 Storage Types Overview

```
AWS Storage Spectrum
──────────────────────────────────────────────────────────────
Block Storage     │   File Storage     │   Object Storage
(EBS)             │   (EFS/FSx)        │   (S3)
──────────────────┼────────────────────┼──────────────────────
- Attached to EC2 │   - Shared across  │   - Unlimited scale
- Low latency     │     instances      │   - HTTP access
- Single AZ*      │   - POSIX compliant│   - Eventual consistency
- Databases       │   - NFS/SMB        │   - Static content
                  │                    │   - Data lakes
```

### 5.2 Amazon S3 Deep Dive

**Storage Classes**:

| Class | Availability | Min Duration | Use Case | Retrieval |
|-------|-------------|--------------|----------|-----------|
| Standard | 99.99% | None | Frequently accessed | Instant |
| Intelligent-Tiering | 99.9% | None | Unknown patterns | Instant |
| Standard-IA | 99.9% | 30 days | Infrequent, quick access | Instant |
| One Zone-IA | 99.5% | 30 days | Reproducible data | Instant |
| Glacier Instant | 99.9% | 90 days | Archive, instant access | Instant |
| Glacier Flexible | 99.99% | 90 days | Archive, minutes-hours | 1-12 hours |
| Glacier Deep Archive | 99.99% | 180 days | Long-term archive | 12-48 hours |

**S3 Security Layers**:

```
┌─────────────────────────────────────────────────────────────┐
│                    S3 Security Model                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Block Public Access (Account/Bucket level)              │
│     └── First line of defense, blocks all public access     │
│                                                              │
│  2. Bucket Policy (Resource-based)                          │
│     └── JSON policy attached to bucket                      │
│                                                              │
│  3. IAM Policy (Identity-based)                             │
│     └── Permissions attached to users/roles                 │
│                                                              │
│  4. ACLs (Legacy - Avoid)                                   │
│     └── Object/Bucket level, limited flexibility            │
│                                                              │
│  5. Encryption                                               │
│     ├── SSE-S3: AWS managed keys (default)                  │
│     ├── SSE-KMS: Customer managed keys (auditable)          │
│     └── SSE-C: Customer provided keys                       │
│                                                              │
│  6. VPC Endpoints                                            │
│     └── Private access without internet                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**S3 Bucket Policy Example**:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowVPCAccess",
            "Effect": "Deny",
            "Principal": "*",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::my-bucket",
                "arn:aws:s3:::my-bucket/*"
            ],
            "Condition": {
                "StringNotEquals": {
                    "aws:sourceVpc": "vpc-1234567890"
                }
            }
        }
    ]
}
```

### 5.3 Amazon EBS (Elastic Block Store)

**Volume Types**:

| Type | IOPS | Throughput | Use Case |
|------|------|------------|----------|
| gp3 | 3,000-16,000 | 125-1,000 MB/s | General purpose (recommended) |
| gp2 | 100-16,000 | 128-250 MB/s | Legacy general purpose |
| io2 Block Express | 256,000 | 4,000 MB/s | Mission-critical DBs |
| io1/io2 | 64,000 | 1,000 MB/s | High-performance DBs |
| st1 | 500 | 500 MB/s | Big data, data warehouses |
| sc1 | 250 | 250 MB/s | Cold storage, infrequent access |

**EBS Snapshots**:

```
EBS Snapshot Strategy
─────────────────────────────────────────────────────────────

        EBS Volume               Incremental Snapshots
        (AZ-specific)            (Region-specific, stored in S3)
             │
             │ First snapshot: Full copy
             ▼
        ┌─────────┐
        │   S1    │ (Full - 100GB)
        └────┬────┘
             │ Second: Only changes
             ▼
        ┌─────────┐
        │   S2    │ (Incremental - 5GB changed)
        └────┬────┘
             │
             ▼
        ┌─────────┐
        │   S3    │ (Incremental - 2GB changed)
        └─────────┘

Features:
- Copy to other regions for DR
- Share with other accounts
- Fast Snapshot Restore (FSR) for quick volume creation
- Archive tier for cost savings (90-day minimum)
```

### 5.4 Amazon EFS (Elastic File System)

**When to Use**:
- Shared storage across multiple EC2 instances
- POSIX-compliant file system needed
- Auto-scaling storage requirements
- Content management, web serving, home directories

**Performance Modes**:

| Mode | Latency | Throughput | Use Case |
|------|---------|------------|----------|
| General Purpose | Low (sub-ms) | Lower | Web serving, CMS |
| Max I/O | Higher | Higher | Big data, parallel processing |

**Throughput Modes**:

| Mode | Behavior | Use Case |
|------|----------|----------|
| Bursting | Scales with size | Variable workloads |
| Provisioned | Fixed throughput | Consistent requirements |
| Elastic | Auto-scales with demand | Unpredictable (recommended) |

---

## 6. Database Services Overview

### 6.1 Database Selection Guide

```
Database Selection Decision Tree
──────────────────────────────────────────────────────────────

Start: What's your data model?
        │
        ├── Relational (SQL, ACID) ──────────────────┐
        │                                             ▼
        │                              ┌─────────────────────────┐
        │                              │ Amazon RDS / Aurora     │
        │                              │ MySQL, PostgreSQL,      │
        │                              │ Oracle, SQL Server      │
        │                              └─────────────────────────┘
        │
        ├── Key-Value (Simple, Fast) ──────────────────┐
        │                                               ▼
        │                              ┌─────────────────────────┐
        │                              │ DynamoDB                │
        │                              │ ElastiCache (Redis)     │
        │                              └─────────────────────────┘
        │
        ├── Document (Flexible Schema) ────────────────┐
        │                                               ▼
        │                              ┌─────────────────────────┐
        │                              │ DynamoDB (with sort)    │
        │                              │ DocumentDB              │
        │                              └─────────────────────────┘
        │
        ├── Graph (Relationships) ──────────────────────┐
        │                                                ▼
        │                              ┌─────────────────────────┐
        │                              │ Neptune                 │
        │                              └─────────────────────────┘
        │
        ├── Time-Series ────────────────────────────────┐
        │                                                ▼
        │                              ┌─────────────────────────┐
        │                              │ Timestream              │
        │                              └─────────────────────────┘
        │
        └── In-Memory (Caching, Sessions) ──────────────┐
                                                         ▼
                                       ┌─────────────────────────┐
                                       │ ElastiCache             │
                                       │ (Redis/Memcached)       │
                                       └─────────────────────────┘
```

### 6.2 Amazon RDS (Relational Database Service)

**Key Features**:
- Managed database (backups, patching, scaling)
- Multi-AZ for high availability
- Read Replicas for read scaling
- Automated failover (Multi-AZ)

**Multi-AZ vs Read Replica**:

```
┌────────────────────────────────────────────────────────────────┐
│                Multi-AZ vs Read Replica                         │
├────────────────────────┬───────────────────────────────────────┤
│       Multi-AZ         │         Read Replica                   │
├────────────────────────┼───────────────────────────────────────┤
│ Purpose: HA/DR         │ Purpose: Read scaling                  │
│ Sync replication       │ Async replication                      │
│ Automatic failover     │ Manual promotion                       │
│ Same region only       │ Cross-region possible                  │
│ Not readable           │ Readable                               │
│ Standby in different   │ Can promote to standalone              │
│ AZ, same endpoint      │ Separate endpoint                      │
└────────────────────────┴───────────────────────────────────────┤

Architecture Example:
                                                                  │
    ┌─────────────┐       Sync        ┌─────────────┐            │
    │   Primary   │───────────────────│   Standby   │            │
    │   (AZ-a)    │                   │   (AZ-b)    │            │
    └──────┬──────┘                   └─────────────┘            │
           │                                                      │
           │ Async                                                │
           │                                                      │
    ┌──────▼──────┐                   ┌─────────────┐            │
    │ Read        │                   │ Read        │            │
    │ Replica 1   │                   │ Replica 2   │            │
    │ (us-east-1) │                   │ (eu-west-1) │            │
    └─────────────┘                   └─────────────┘            │
```

### 6.3 Amazon Aurora

**Why Aurora?**:
- 5x throughput of MySQL, 3x of PostgreSQL
- Storage auto-scales up to 128TB
- 6 copies of data across 3 AZs
- Up to 15 read replicas
- Serverless option available

**Aurora Architecture**:

```
Aurora Storage Architecture
──────────────────────────────────────────────────────────────

            ┌─────────────────────────────────┐
            │        Aurora Cluster            │
            │                                  │
            │  ┌──────────┐   ┌──────────┐    │
            │  │  Writer  │   │  Reader  │    │
            │  │ Instance │   │ Instance │    │
            │  └────┬─────┘   └────┬─────┘    │
            │       │              │          │
            └───────┼──────────────┼──────────┘
                    │              │
                    └──────┬───────┘
                           │
    ┌──────────────────────▼──────────────────────┐
    │         Aurora Storage (Shared)              │
    │                                              │
    │  ┌────────────────────────────────────────┐ │
    │  │            AZ-a        AZ-b       AZ-c │ │
    │  │           ┌───┐      ┌───┐      ┌───┐ │ │
    │  │           │ 1 │      │ 2 │      │ 3 │ │ │
    │  │           │ 4 │      │ 5 │      │ 6 │ │ │
    │  │           └───┘      └───┘      └───┘ │ │
    │  │                                        │ │
    │  │    6 copies across 3 AZs               │ │
    │  │    - Tolerates loss of 2 copies (read) │ │
    │  │    - Tolerates loss of 3 copies (write)│ │
    │  └────────────────────────────────────────┘ │
    └─────────────────────────────────────────────┘
```

### 6.4 Amazon DynamoDB

**Key Concepts**:

```
DynamoDB Data Model
──────────────────────────────────────────────────────────────

Table: Users
    │
    ├── Partition Key (PK): UserID
    │   └── Determines physical partition
    │
    └── Sort Key (SK): Timestamp (optional)
        └── Enables range queries within partition

Example Item:
{
    "UserID": "user123",         // Partition Key
    "Timestamp": "2024-01-15",   // Sort Key
    "Name": "John Doe",
    "Email": "john@example.com",
    "Orders": [...]              // Nested document
}

Access Patterns:
- Get single item: PK + SK
- Query: PK + SK conditions
- Scan: Full table (expensive, avoid)
```

**Capacity Modes**:

| Mode | Pricing | Use Case |
|------|---------|----------|
| On-Demand | Pay per request | Unpredictable traffic |
| Provisioned | Pay for capacity | Predictable, consistent |
| Provisioned + Auto Scaling | Hybrid | Variable but bounded |

**DynamoDB Global Tables**:

```
Multi-Region Active-Active
──────────────────────────────────────────────────────────────

        us-east-1                     eu-west-1
        ─────────                     ─────────
        ┌───────────┐                 ┌───────────┐
        │ DynamoDB  │ ←── Async ────→ │ DynamoDB  │
        │  Table    │    Replication  │  Table    │
        └───────────┘                 └───────────┘
              ↑                             ↑
              │                             │
         Application                   Application
         (US Users)                    (EU Users)

Features:
- Sub-second replication
- Conflict resolution (last writer wins)
- No application changes needed
```

---

## 7. AWS Well-Architected Framework

### 7.1 The Six Pillars

```
┌────────────────────────────────────────────────────────────────┐
│              AWS Well-Architected Framework                     │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. OPERATIONAL EXCELLENCE                                      │
│     │ Run and monitor systems to deliver business value         │
│     └─ Automate, evolve, learn from failures                   │
│                                                                 │
│  2. SECURITY                                                    │
│     │ Protect information and systems                           │
│     └─ Defense in depth, least privilege, traceability         │
│                                                                 │
│  3. RELIABILITY                                                 │
│     │ Recover from failures, meet demand                        │
│     └─ Distributed design, auto-recovery, change management    │
│                                                                 │
│  4. PERFORMANCE EFFICIENCY                                      │
│     │ Use resources efficiently                                 │
│     └─ Right-sizing, advanced tech, global reach               │
│                                                                 │
│  5. COST OPTIMIZATION                                           │
│     │ Avoid unnecessary costs                                   │
│     └─ Consumption model, efficiency, measured spending        │
│                                                                 │
│  6. SUSTAINABILITY                                              │
│     │ Minimize environmental impact                             │
│     └─ Understand impact, maximize utilization, reduce waste   │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### 7.2 Key Design Principles

```
Pillar                  │ Key Principles
────────────────────────┼─────────────────────────────────────────
Operational Excellence  │ • IaC (Infrastructure as Code)
                       │ • Small, reversible changes
                       │ • Frequent procedures refinement
                       │ • Anticipate failure
────────────────────────┼─────────────────────────────────────────
Security               │ • Strong identity foundation
                       │ • Enable traceability
                       │ • Apply security at all layers
                       │ • Automate security best practices
                       │ • Protect data in transit and at rest
────────────────────────┼─────────────────────────────────────────
Reliability            │ • Auto-recover from failure
                       │ • Test recovery procedures
                       │ • Scale horizontally
                       │ • Stop guessing capacity
                       │ • Manage change with automation
────────────────────────┼─────────────────────────────────────────
Performance            │ • Democratize advanced technologies
                       │ • Go global in minutes
                       │ • Use serverless architectures
                       │ • Experiment more often
                       │ • Consider mechanical sympathy
────────────────────────┼─────────────────────────────────────────
Cost Optimization      │ • Implement Cloud Financial Management
                       │ • Adopt consumption model
                       │ • Measure overall efficiency
                       │ • Stop spending on undifferentiated work
                       │ • Analyze and attribute expenditure
────────────────────────┼─────────────────────────────────────────
Sustainability         │ • Understand your impact
                       │ • Establish sustainability goals
                       │ • Maximize utilization
                       │ • Anticipate and adopt efficient offerings
                       │ • Reduce downstream impact
```

### 7.3 Applying the Framework

**Question-Based Approach**:

For each pillar, AWS provides specific questions. Example for Reliability:

```
Reliability Questions (Sample)
──────────────────────────────────────────────────────────────

REL 1: How do you manage service quotas and constraints?
       → Monitor usage, request increases proactively

REL 2: How do you plan your network topology?
       → Use multiple AZs, private subnets, redundant connectivity

REL 3: How do you design your workload service architecture?
       → Loosely coupled components, idempotent operations

REL 4: How do you design interactions to prevent failures?
       → Throttling, circuit breakers, fail fast

REL 5: How do you design interactions to mitigate failures?
       → Retries with backoff, idempotency, graceful degradation
```

---

## 8. Cost Management Fundamentals

### 8.1 AWS Pricing Models

```
Cost Components
──────────────────────────────────────────────────────────────

1. COMPUTE
   └── EC2 hours, Lambda invocations, Fargate vCPU-hours

2. STORAGE
   └── GB-month stored, requests, data retrieval

3. DATA TRANSFER
   ├── Inbound: Free
   ├── Outbound to Internet: $0.09/GB (tiered)
   ├── Between Regions: $0.01-0.02/GB
   └── Within AZ: Free (private IP)

4. REQUEST/API CALLS
   └── Per request pricing (S3, API Gateway)
```

### 8.2 Cost Optimization Strategies

```
Strategy                │ Implementation                         │ Savings
────────────────────────┼───────────────────────────────────────┼─────────
Right-sizing           │ Use Compute Optimizer, downsize unused │ 10-30%
Reserved/Savings Plans │ Commit to 1-3 year terms               │ 30-60%
Spot Instances         │ Fault-tolerant workloads               │ Up to 90%
Storage Tiering        │ S3 Lifecycle, EBS right-sizing         │ 20-50%
Cleanup Unused         │ Delete snapshots, stop idle resources  │ Variable
Architectural          │ Serverless, caching, CDN               │ 30-70%
```

### 8.3 Cost Management Tools

| Tool | Purpose |
|------|---------|
| **Cost Explorer** | Visualize and analyze spending |
| **Budgets** | Set alerts for cost/usage thresholds |
| **Cost Anomaly Detection** | ML-powered unusual spending alerts |
| **Compute Optimizer** | Right-sizing recommendations |
| **Trusted Advisor** | Cost optimization checks |
| **Cost & Usage Report** | Detailed billing data for analysis |

---

## 9. Architecture Decision Framework

### 9.1 Common Architecture Patterns

**Three-Tier Web Application**:

```
┌─────────────────────────────────────────────────────────────────┐
│                     Three-Tier Architecture                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│                        ┌─────────────┐                          │
│                        │  Route 53   │                          │
│                        │    (DNS)    │                          │
│                        └──────┬──────┘                          │
│                               │                                  │
│                        ┌──────▼──────┐                          │
│                        │ CloudFront  │                          │
│                        │   (CDN)     │                          │
│                        └──────┬──────┘                          │
│                               │                                  │
│  ┌───────────────────────────┬┴────────────────────────────┐    │
│  │                           │                              │    │
│  │                    ┌──────▼──────┐                       │    │
│  │    PRESENTATION    │     ALB     │                       │    │
│  │       TIER        │ (Public SN)  │                       │    │
│  │                    └──────┬──────┘                       │    │
│  │                           │                              │    │
│  │────────────────────────────────────────────────────────│    │
│  │                           │                              │    │
│  │                 ┌─────────┴─────────┐                   │    │
│  │    APPLICATION  │                   │                   │    │
│  │       TIER     ┌▼──────┐      ┌─────▼─┐                │    │
│  │               │  EC2   │      │  EC2  │                │    │
│  │               │ (ASG)  │      │ (ASG) │                │    │
│  │               │ AZ-a   │      │ AZ-b  │                │    │
│  │               └───┬────┘      └───┬───┘                │    │
│  │                   │              │                      │    │
│  │────────────────────────────────────────────────────────│    │
│  │                   │              │                      │    │
│  │       DATA       ┌▼──────────────▼┐                    │    │
│  │       TIER       │   Aurora       │                    │    │
│  │                  │   (Multi-AZ)   │                    │    │
│  │                  └────────────────┘                    │    │
│  │                                                         │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Event-Driven Architecture**:

```
┌─────────────────────────────────────────────────────────────────┐
│                   Event-Driven Architecture                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Event Producers           Event Router        Event Consumers   │
│  ──────────────           ────────────         ───────────────   │
│                                                                  │
│  ┌──────────┐         ┌────────────────┐       ┌──────────┐     │
│  │   API    │────────→│                │──────→│  Lambda  │     │
│  │ Gateway  │         │                │       │ (Process)│     │
│  └──────────┘         │   EventBridge  │       └──────────┘     │
│                       │       or       │                         │
│  ┌──────────┐         │      SNS       │       ┌──────────┐     │
│  │   S3     │────────→│                │──────→│   SQS    │     │
│  │ (Upload) │         │   (Fan-out)    │       │ (Queue)  │     │
│  └──────────┘         │                │       └────┬─────┘     │
│                       │                │            │            │
│  ┌──────────┐         │                │       ┌────▼─────┐     │
│  │DynamoDB  │────────→│                │──────→│   ECS    │     │
│  │ Streams  │         └────────────────┘       │ (Batch)  │     │
│  └──────────┘                                  └──────────┘     │
│                                                                  │
│  Benefits:                                                       │
│  • Loose coupling between components                            │
│  • Easy to add new consumers                                    │
│  • Built-in retry and DLQ                                       │
│  • Scale each component independently                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Microservices on ECS/EKS**:

```
┌─────────────────────────────────────────────────────────────────┐
│                   Microservices Architecture                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                      API Gateway                         │    │
│  │                   (Rate Limiting, Auth)                  │    │
│  └─────────────────────────┬───────────────────────────────┘    │
│                            │                                     │
│  ┌─────────────────────────▼───────────────────────────────┐    │
│  │                   App Mesh / Service Mesh                │    │
│  │                 (Service Discovery, Traffic)             │    │
│  └─────────────────────────┬───────────────────────────────┘    │
│                            │                                     │
│  ┌──────────┬──────────┬───┴────┬──────────┬──────────┐        │
│  │          │          │        │          │          │         │
│  ▼          ▼          ▼        ▼          ▼          ▼         │
│ ┌────┐    ┌────┐    ┌────┐   ┌────┐    ┌────┐    ┌────┐       │
│ │Svc │    │Svc │    │Svc │   │Svc │    │Svc │    │Svc │       │
│ │ A  │    │ B  │    │ C  │   │ D  │    │ E  │    │ F  │       │
│ └──┬─┘    └──┬─┘    └──┬─┘   └──┬─┘    └──┬─┘    └──┬─┘       │
│    │         │         │        │         │         │          │
│    ▼         ▼         ▼        ▼         ▼         ▼          │
│ ┌────┐    ┌────┐    ┌────┐   ┌────┐    ┌────┐    ┌────┐       │
│ │ DB │    │Cache│    │ DB │   │Queue│   │ S3 │    │ DB │       │
│ └────┘    └────┘    └────┘   └────┘    └────┘    └────┘       │
│                                                                  │
│  Cross-Cutting Concerns:                                        │
│  • X-Ray (Distributed Tracing)                                  │
│  • CloudWatch (Logs, Metrics)                                   │
│  • Secrets Manager (Configuration)                              │
│  • Parameter Store (Non-sensitive config)                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 9.2 Architecture Decision Records (ADRs)

When making architecture decisions, document them:

```markdown
# ADR-001: Database Selection for Order Service

## Status
Accepted

## Context
We need a database for the Order service that handles:
- High write throughput during peak hours (10K orders/min)
- Read queries by customer ID and order date
- No complex joins required

## Decision
Use DynamoDB with:
- Partition Key: CustomerID
- Sort Key: OrderDate
- Global Secondary Index: OrderID (for direct lookups)

## Consequences
### Positive
- Auto-scaling handles traffic spikes
- Single-digit ms latency
- No maintenance overhead

### Negative
- Limited query flexibility
- Higher cost if data access patterns change
- Team needs DynamoDB expertise

## Alternatives Considered
1. Aurora PostgreSQL - Better query flexibility but higher ops overhead
2. DocumentDB - Similar to DynamoDB but less integrated
```

---

## 10. Interview Questions & Scenarios

### 10.1 Foundational Questions

**Q1: Explain the difference between a Region and an Availability Zone.**

**Answer**:
- **Region**: A geographical area with 2+ AZs (e.g., us-east-1). Completely isolated from other regions. Choose based on compliance, latency, and service availability.
- **AZ**: Isolated data center(s) within a region. Connected via low-latency links (<2ms). Designed for fault isolation.
- **Key Point**: Data doesn't replicate between regions automatically but can sync within AZs with low latency.

---

**Q2: How would you implement least privilege access for a team of developers who need to deploy Lambda functions?**

**Answer**:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "LambdaDeploymentPermissions",
            "Effect": "Allow",
            "Action": [
                "lambda:CreateFunction",
                "lambda:UpdateFunctionCode",
                "lambda:UpdateFunctionConfiguration",
                "lambda:GetFunction",
                "lambda:ListFunctions"
            ],
            "Resource": "arn:aws:lambda:us-east-1:123456789012:function:team-*",
            "Condition": {
                "StringEquals": {
                    "aws:RequestTag/team": "${aws:PrincipalTag/team}"
                }
            }
        },
        {
            "Sid": "CloudWatchLogsAccess",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogGroups"
            ],
            "Resource": "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/team-*"
        }
    ]
}
```

Key principles applied:
- Resource-level restrictions (only `team-*` functions)
- Tag-based access control
- Only necessary permissions
- Condition-based restrictions

---

**Q3: What's the difference between Security Groups and NACLs? When would you use each?**

**Answer**:

| Aspect | Security Groups | NACLs |
|--------|----------------|-------|
| Level | Instance (ENI) | Subnet |
| State | Stateful | Stateless |
| Rules | Allow only | Allow & Deny |
| Default | Deny all inbound | Allow all |
| Order | All rules evaluated | Sequential (rule #) |

**Use Cases**:
- **Security Groups**: Primary defense, instance-level control, application-aware rules
- **NACLs**: Subnet-wide blocking, compliance requirements (explicit deny), defense-in-depth

---

### 10.2 Architecture Scenario Questions

**Scenario 1: Design a highly available web application**

**Question**: Design an architecture for an e-commerce website that needs to handle 50,000 concurrent users, with 99.9% availability requirement.

**Answer Framework**:

```
Requirements Analysis:
- 50K concurrent users → ~5,000 req/s (assuming 10 req/user/min)
- 99.9% availability → 8.76 hours downtime/year
- E-commerce → Payment processing, inventory, user sessions

Architecture:

1. DNS & CDN Layer
   - Route 53 with health checks
   - CloudFront for static assets + API acceleration
   - Multi-region failover capability

2. Load Balancing
   - ALB in 2+ AZs
   - Path-based routing (/api, /static, /checkout)
   - WAF for security

3. Application Layer
   - ECS Fargate or EKS (containerized)
   - Auto Scaling (target tracking: 70% CPU)
   - Blue-green deployments

4. Data Layer
   - Aurora PostgreSQL (Multi-AZ) for orders
   - DynamoDB for user sessions (DAX for caching)
   - ElastiCache Redis for product catalog caching

5. Ancillary Services
   - S3 for images/static content
   - SQS for order processing queue
   - Lambda for async tasks (emails, inventory updates)

Cost Optimization:
- Reserved instances for baseline
- Spot instances for batch processing
- S3 Intelligent-Tiering for logs
```

---

**Scenario 2: Data migration strategy**

**Question**: Your company is migrating a 50TB PostgreSQL database from on-premise to AWS. The application can tolerate maximum 4 hours of downtime. What's your approach?

**Answer**:

```
Migration Strategy: AWS DMS with minimal downtime

Phase 1: Setup (Week 1-2)
─────────────────────────
• Provision Aurora PostgreSQL
• Set up Direct Connect or VPN
• Configure AWS DMS replication instance

Phase 2: Initial Load (Week 2-3)
────────────────────────────────
• Full load via DMS (~24-48 hours for 50TB)
• Keep source DB operational
• Validate data in target

Phase 3: CDC (Change Data Capture) - Week 3-4
─────────────────────────────────────────────
• Enable ongoing replication
• Monitor replication lag
• Validate data consistency

Phase 4: Cutover (Planned Window)
─────────────────────────────────
1. Stop application writes (T+0)
2. Wait for replication to complete (T+1 hour)
3. Validate data integrity (T+2 hours)
4. Switch application connection string (T+3 hours)
5. Smoke test and go-live (T+4 hours)

Rollback Plan:
• Keep source DB running for 48 hours
• DMS reverse replication ready
• DNS-based failback
```

---

**Scenario 3: Cost optimization review**

**Question**: You've been asked to reduce AWS costs by 30% for a production workload. How would you approach this?

**Answer**:

```
Step 1: Assessment (Week 1)
─────────────────────────
• Cost Explorer analysis by service, tag, region
• Trusted Advisor cost recommendations
• Compute Optimizer for right-sizing
• Identify unused resources

Step 2: Quick Wins (Week 2-3)
────────────────────────────
• Delete unused EBS volumes/snapshots
• Terminate stopped instances
• Remove unattached Elastic IPs
• Review and clean up old AMIs

Step 3: Right-Sizing (Week 3-4)
─────────────────────────────
• Downsize over-provisioned instances
• gp2 → gp3 migration (20% savings)
• Review RDS instance sizes

Step 4: Commitment Decisions (Week 4-5)
──────────────────────────────────────
• Savings Plans for baseline compute
• Reserved Instances for RDS/ElastiCache
• Spot Instances for batch workloads

Step 5: Architectural Changes (Ongoing)
──────────────────────────────────────
• Implement S3 Lifecycle policies
• Review data transfer patterns
• Consider Graviton instances (up to 40% savings)
• Evaluate serverless options

Expected Savings Breakdown:
• Right-sizing: 10-15%
• Commitments: 30-40%
• Cleanup: 5-10%
• Architecture: 10-20%
```

---

### 10.3 Troubleshooting Scenarios

**Scenario**: EC2 instance in private subnet can't access the internet

**Diagnostic Checklist**:

```
1. NAT Gateway Check
   □ NAT Gateway exists in public subnet
   □ NAT Gateway in same AZ or route to NAT in other AZ
   □ NAT Gateway has Elastic IP

2. Route Table Check
   □ Private subnet using correct route table
   □ Route: 0.0.0.0/0 → nat-xxx

3. Security Group Check
   □ Outbound rule allows traffic (usually all)
   □ Response traffic not blocked

4. NACL Check
   □ Outbound & Inbound rules allow traffic
   □ Ephemeral ports open (1024-65535)

5. VPC Check
   □ DNS resolution enabled
   □ DNS hostnames enabled

Debug Commands:
─────────────────
# Check route table
aws ec2 describe-route-tables --filters "Name=association.subnet-id,Values=subnet-xxx"

# Check NAT Gateway status
aws ec2 describe-nat-gateways --filter "Name=state,Values=available"

# Test connectivity from instance
curl -v https://aws.amazon.com
```

---

## Summary: Key Takeaways for Architects

```
┌────────────────────────────────────────────────────────────────┐
│              Architect's Checklist - AWS Foundation             │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│ □ Global Infrastructure                                         │
│   • Region selection based on compliance, latency, cost        │
│   • Multi-AZ for high availability                             │
│   • Edge locations for global performance                       │
│                                                                 │
│ □ Security (IAM)                                                │
│   • Roles over users                                           │
│   • Least privilege always                                     │
│   • SCPs for organizational guardrails                         │
│                                                                 │
│ □ Networking (VPC)                                              │
│   • Public/private subnet separation                           │
│   • Security groups as primary defense                         │
│   • VPC endpoints for AWS service access                       │
│                                                                 │
│ □ Compute Selection                                             │
│   • EC2 for control, containers for portability                │
│   • Lambda for event-driven, short-duration                    │
│   • Fargate for serverless containers                          │
│                                                                 │
│ □ Storage Strategy                                              │
│   • S3 for objects, EBS for block, EFS for shared files       │
│   • Lifecycle policies for cost optimization                   │
│   • Encryption at rest and in transit                          │
│                                                                 │
│ □ Database Selection                                            │
│   • RDS/Aurora for relational                                  │
│   • DynamoDB for key-value at scale                           │
│   • Caching layer for performance                              │
│                                                                 │
│ □ Well-Architected Review                                       │
│   • Apply six pillars to every design                          │
│   • Document decisions with ADRs                               │
│   • Regular review and optimization                            │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

---

## Next Steps

This completes **Part 1: AWS Foundations**. 

**Part 2** will cover:
- Advanced Networking (Transit Gateway, Direct Connect, VPN)
- Serverless Deep Dive (Lambda, API Gateway, Step Functions)
- Container Services (ECS, EKS, ECR)
- Observability (CloudWatch, X-Ray, CloudTrail)
- Infrastructure as Code (CloudFormation, CDK, Terraform)

**Part 3** will cover:
- Advanced Architecture Patterns
- Migration Strategies (6 Rs)
- Multi-Account Strategy
- DR and Business Continuity
- Performance Optimization
- Security Deep Dive (GuardDuty, Security Hub, Inspector)
- Real-world Case Studies

---

*Last Updated: February 2026*
