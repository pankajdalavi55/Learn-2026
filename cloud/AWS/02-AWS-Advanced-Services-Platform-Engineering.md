# AWS Advanced Services: Platform Engineering Perspective

> **Target Audience**: Senior Software Engineers transitioning to Platform Engineer / Architect roles
> **Prerequisites**: Completed Part 1 (AWS Foundations), hands-on AWS experience

---

## Table of Contents

1. [Advanced Networking](#1-advanced-networking)
2. [Serverless Deep Dive](#2-serverless-deep-dive)
3. [Container Services](#3-container-services)
4. [Observability & Monitoring](#4-observability--monitoring)
5. [Infrastructure as Code](#5-infrastructure-as-code)
6. [Interview Questions & Scenarios](#6-interview-questions--scenarios)

---

## 1. Advanced Networking

### 1.1 Network Architecture Evolution

```
Network Complexity Progression
──────────────────────────────────────────────────────────────────────

Level 1: Single VPC
└── Simple applications, small teams

Level 2: VPC Peering
└── 2-5 VPCs, point-to-point connections

Level 3: Transit Gateway
└── Hub-and-spoke, 10+ VPCs, hybrid connectivity

Level 4: Global Network
└── Multi-region, Direct Connect, complex routing
```

### 1.2 AWS Transit Gateway (TGW)

**The Problem It Solves**:

```
Before Transit Gateway (VPC Peering Mesh):
─────────────────────────────────────────────────────────────────

     VPC-A ─────────── VPC-B
       │ \           / │
       │   \       /   │
       │     \   /     │
       │       X       │      n VPCs = n(n-1)/2 connections
       │     /   \     │      10 VPCs = 45 peering connections!
       │   /       \   │
       │ /           \ │
     VPC-C ─────────── VPC-D

Problems:
• Non-transitive (A→B, B→C doesn't mean A→C)
• Management complexity scales quadratically
• No centralized routing control
```

```
After Transit Gateway (Hub-and-Spoke):
─────────────────────────────────────────────────────────────────

                    VPC-A     VPC-B
                       \       /
                        \     /
                      ┌────────────┐
                      │  Transit   │
     On-Premise ──────│  Gateway   │────── VPN
                      │   (TGW)    │
                      └────────────┘
                        /     \
                       /       \
                    VPC-C     VPC-D

Benefits:
• Centralized routing
• Transitive connectivity
• Scales to 5,000 attachments
• Supports VPN, Direct Connect, peering
```

**TGW Core Concepts**:

| Component | Description |
|-----------|-------------|
| **Attachment** | Connection to VPC, VPN, Direct Connect, or peering |
| **Route Table** | Routing rules for traffic between attachments |
| **Association** | Links attachment to a route table |
| **Propagation** | Automatic route learning from attachments |

**TGW Architecture Example**:

```
Enterprise Transit Gateway Design
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        AWS Transit Gateway                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Route Tables:                                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  │
│  │   Production RT   │  │  Development RT   │  │   Shared Svcs RT │  │
│  │                   │  │                   │  │                   │  │
│  │ 10.1.0.0/16 → A1 │  │ 10.2.0.0/16 → A3 │  │ 0.0.0.0/0 → A5   │  │
│  │ 10.2.0.0/16 → A2 │  │ 10.100.0.0/16→A5 │  │ 10.0.0.0/8 → All │  │
│  │ 10.100.0.0/16→A5 │  │                   │  │                   │  │
│  │ 0.0.0.0/0 → A6   │  │                   │  │                   │  │
│  └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘  │
│           │                     │                      │            │
│  Associations:                                                       │
│  ┌────────▼─────────┐  ┌────────▼─────────┐  ┌────────▼─────────┐  │
│  │  Prod VPCs       │  │  Dev VPCs        │  │  Shared Services  │  │
│  │  A1: Prod-VPC-1  │  │  A3: Dev-VPC-1   │  │  A5: Shared-VPC   │  │
│  │  A2: Prod-VPC-2  │  │  A4: Dev-VPC-2   │  │  A6: Egress-VPC   │  │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘  │
│                                                                      │
│  External Attachments:                                               │
│  A7: VPN to On-Premise                                              │
│  A8: Direct Connect Gateway                                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Segmentation Benefits:
• Prod can't route to Dev directly
• Shared services accessible from both
• Centralized egress through Egress-VPC
• On-premise has controlled access
```

**TGW Terraform Example**:

```hcl
# Transit Gateway
resource "aws_ec2_transit_gateway" "main" {
  description                     = "Main Transit Gateway"
  auto_accept_shared_attachments  = "enable"
  default_route_table_association = "disable"
  default_route_table_propagation = "disable"
  
  tags = {
    Name = "main-tgw"
  }
}

# VPC Attachment
resource "aws_ec2_transit_gateway_vpc_attachment" "prod" {
  subnet_ids         = var.prod_subnet_ids
  transit_gateway_id = aws_ec2_transit_gateway.main.id
  vpc_id             = var.prod_vpc_id
  
  tags = {
    Name = "prod-attachment"
  }
}

# Route Table
resource "aws_ec2_transit_gateway_route_table" "prod" {
  transit_gateway_id = aws_ec2_transit_gateway.main.id
  
  tags = {
    Name = "prod-route-table"
  }
}

# Association
resource "aws_ec2_transit_gateway_route_table_association" "prod" {
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.prod.id
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.prod.id
}

# Static Route
resource "aws_ec2_transit_gateway_route" "to_shared" {
  destination_cidr_block         = "10.100.0.0/16"
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.shared.id
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.prod.id
}
```

### 1.3 AWS Direct Connect

**Overview**: Dedicated, private network connection from on-premise to AWS.

```
Direct Connect Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────┐    ┌──────────────────┐    ┌─────────────────────┐
│   On-Premise    │    │  Direct Connect  │    │        AWS          │
│   Data Center   │    │    Location      │    │                     │
│                 │    │                  │    │                     │
│ ┌─────────────┐ │    │ ┌──────────────┐ │    │  ┌───────────────┐  │
│ │   Router    │─┼────┼─│ Customer or  │─┼────┼──│ Direct Connect│  │
│ │             │ │ X  │ │ Partner      │ │    │  │    Router     │  │
│ └─────────────┘ │ C  │ │ Equipment    │ │    │  └───────┬───────┘  │
│                 │ o  │ └──────────────┘ │    │          │          │
│                 │ n  │                  │    │  ┌───────▼───────┐  │
│                 │ n  │                  │    │  │    Private    │  │
│                 │ e  │                  │    │  │ Virtual Inter-│  │
│                 │ c  │                  │    │  │     face      │──┼──→ VPC
│                 │ t  │                  │    │  └───────────────┘  │
│                 │    │                  │    │                     │
│                 │    │                  │    │  ┌───────────────┐  │
│                 │    │                  │    │  │    Public     │  │
│                 │    │                  │    │  │ Virtual Inter-│──┼──→ S3, DynamoDB
│                 │    │                  │    │  │     face      │  │
│                 │    │                  │    │  └───────────────┘  │
└─────────────────┘    └──────────────────┘    └─────────────────────┘
      │                                                    │
      └──── Physical Cross-Connect (1Gbps / 10Gbps / 100Gbps) ───┘
```

**Connection Types**:

| Type | Speed | Use Case |
|------|-------|----------|
| **Dedicated Connection** | 1/10/100 Gbps | High bandwidth, predictable traffic |
| **Hosted Connection** | 50 Mbps - 10 Gbps | Partner-provided, flexible capacity |

**Virtual Interfaces (VIFs)**:

| VIF Type | Purpose | BGP Requirement |
|----------|---------|-----------------|
| **Private VIF** | Access VPCs via private IPs | Private ASN |
| **Public VIF** | Access AWS public services | Public ASN |
| **Transit VIF** | Connect to Transit Gateway | Private ASN |

**Direct Connect Gateway**:

```
Direct Connect Gateway - Multi-Region Access
──────────────────────────────────────────────────────────────────────

                          ┌─────────────────────┐
                          │  Direct Connect     │
                          │      Gateway        │
                          └──────────┬──────────┘
                                     │
              ┌──────────────────────┼──────────────────────┐
              │                      │                      │
              ▼                      ▼                      ▼
    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
    │  us-east-1 VPC  │    │  eu-west-1 VPC  │    │  ap-south-1 VPC │
    │  10.1.0.0/16    │    │  10.2.0.0/16    │    │  10.3.0.0/16    │
    └─────────────────┘    └─────────────────┘    └─────────────────┘

Benefits:
• Single Direct Connect, multiple regions
• No need for separate connections per region
• Works with Transit Gateway (Transit VIF)
```

**High Availability Design**:

```
Direct Connect HA Architecture
──────────────────────────────────────────────────────────────────────

                      Primary DC Location        Secondary DC Location
                      (Different facility)       (Different facility)
                            │                          │
On-Premise ─────────────────┼──────────────────────────┤
Router 1 ───────────────────┤                          │
                            │                          │
                       ┌────▼────┐                ┌────▼────┐
                       │  DX 1   │                │  DX 2   │
                       │ (1Gbps) │                │ (1Gbps) │
                       └────┬────┘                └────┬────┘
                            │                          │
                       ┌────▼──────────────────────────▼────┐
                       │       Direct Connect Gateway        │
                       └────────────────┬───────────────────┘
                                        │
                                ┌───────▼───────┐
                                │ Transit       │
                                │ Gateway       │
                                └───────────────┘

Failover:
• BGP handles automatic failover
• AS Path prepending for active/passive
• Typical failover: 30-60 seconds
```

### 1.4 AWS Site-to-Site VPN

**When to Use**:
- Quick setup needed (vs. weeks for Direct Connect)
- Backup for Direct Connect
- Lower bandwidth requirements (<1.25 Gbps per tunnel)
- Cost-sensitive scenarios

```
VPN Architecture
──────────────────────────────────────────────────────────────────────

On-Premise                                          AWS
──────────────────                        ──────────────────────────
                                          
┌──────────────┐                          ┌─────────────────────────┐
│              │                          │     Virtual Private     │
│   Customer   │◄─────── Tunnel 1 ───────►│       Gateway (VGW)     │
│   Gateway    │                          │           or            │
│   (Router)   │◄─────── Tunnel 2 ───────►│    Transit Gateway      │
│              │                          │                         │
└──────────────┘                          └───────────┬─────────────┘
       │                                              │
       │                                              │
  On-Premise                                    ┌─────▼─────┐
   Network                                      │    VPC    │
 10.0.0.0/8                                     │ 172.16.0.0│
                                                └───────────┘

Tunnel Details:
• 2 tunnels per VPN connection (HA)
• IPsec encrypted
• Up to 1.25 Gbps per tunnel
• AWS-managed endpoints (no maintenance)
```

**VPN Options Comparison**:

| Feature | VGW-based VPN | TGW-based VPN |
|---------|---------------|---------------|
| Max Bandwidth | 1.25 Gbps | 1.25 Gbps × tunnels |
| ECMP | No | Yes (up to 50 tunnels) |
| Transitive Routing | No | Yes |
| Multiple VPCs | Need per VPC | Single connection |
| Best For | Simple setups | Enterprise scale |

**Accelerated Site-to-Site VPN**:

```
Standard VPN vs Accelerated VPN
──────────────────────────────────────────────────────────────────────

Standard VPN:
On-Premise ────── Public Internet ────── AWS VPN Endpoint
                 (Variable latency)

Accelerated VPN:
On-Premise ──► Edge Location ──► AWS Global Network ──► VPN Endpoint
              (Closest PoP)     (Optimized backbone)

Benefits of Accelerated:
• Uses AWS Global Accelerator network
• More consistent latency
• Better for latency-sensitive applications
• ~$36/month additional cost
```

### 1.5 Network Connectivity Decision Matrix

```
Connectivity Selection Guide
──────────────────────────────────────────────────────────────────────

Start: What are your requirements?
        │
        ├── Need for speed (>1 Gbps)?
        │   │
        │   ├── Yes ────► Direct Connect (Dedicated)
        │   │             └── Add VPN for backup
        │   │
        │   └── No ─────► Site-to-Site VPN
        │                 └── Consider Accelerated for consistency
        │
        ├── Multiple VPCs to connect?
        │   │
        │   ├── <5 VPCs, same region ────► VPC Peering
        │   │
        │   └── 5+ VPCs or multi-region ── Transit Gateway
        │
        └── Need private access to AWS services (S3, DynamoDB)?
            │
            ├── From VPC ────► VPC Endpoints (Gateway or Interface)
            │
            └── From On-Premise ────► Direct Connect Public VIF
                                      or VPC Endpoint + DX/VPN
```

### 1.6 Advanced VPC Features

**VPC Traffic Mirroring**:

```
Traffic Mirroring Architecture
──────────────────────────────────────────────────────────────────────

┌──────────────────────────────────────────────────────────────────┐
│                              VPC                                  │
│                                                                   │
│  ┌─────────────────┐          ┌─────────────────┐                │
│  │   Source ENI    │──────────│   Mirror        │                │
│  │   (EC2/RDS)     │  Copy    │   Target        │                │
│  │                 │  Traffic │   (NLB/ENI)     │                │
│  └─────────────────┘          └────────┬────────┘                │
│                                        │                          │
│                               ┌────────▼────────┐                │
│                               │   IDS/Packet    │                │
│                               │   Analysis      │                │
│                               │   (Suricata)    │                │
│                               └─────────────────┘                │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘

Use Cases:
• Intrusion detection
• Network forensics
• Compliance monitoring
• Threat analysis
```

**VPC Flow Logs Enhanced**:

```python
# Flow Log Record Format (Version 5)
{
    "version": 5,
    "account-id": "123456789012",
    "interface-id": "eni-1234567890abcdef0",
    "srcaddr": "10.0.1.5",
    "dstaddr": "10.0.2.10",
    "srcport": 49152,
    "dstport": 443,
    "protocol": 6,  # TCP
    "packets": 10,
    "bytes": 4200,
    "start": 1620140661,
    "end": 1620140721,
    "action": "ACCEPT",
    "log-status": "OK",
    # New in v5
    "vpc-id": "vpc-12345678",
    "subnet-id": "subnet-12345678",
    "instance-id": "i-1234567890abcdef0",
    "tcp-flags": 2,  # SYN
    "type": "IPv4",
    "pkt-srcaddr": "10.0.1.5",
    "pkt-dstaddr": "10.0.2.10",
    "region": "us-east-1",
    "az-id": "use1-az1",
    "sublocation-type": "",
    "sublocation-id": "",
    "pkt-src-aws-service": "",
    "pkt-dst-aws-service": "AMAZON",
    "flow-direction": "egress",
    "traffic-path": 1
}

# Analysis Query (CloudWatch Logs Insights)
fields @timestamp, srcaddr, dstaddr, dstport, action
| filter action = "REJECT"
| stats count(*) as rejections by srcaddr, dstaddr, dstport
| sort rejections desc
| limit 20
```

**AWS Network Firewall**:

```
Network Firewall Architecture
──────────────────────────────────────────────────────────────────────

                         Internet
                             │
                    ┌────────▼────────┐
                    │  Internet       │
                    │  Gateway        │
                    └────────┬────────┘
                             │
              ┌──────────────▼──────────────┐
              │   Firewall Subnet           │
              │   ┌──────────────────────┐  │
              │   │   AWS Network        │  │
              │   │   Firewall           │  │
              │   │                      │  │
              │   │ • Stateful rules     │  │
              │   │ • Domain filtering   │  │
              │   │ • IPS signatures     │  │
              │   │ • TLS inspection     │  │
              │   └──────────────────────┘  │
              └──────────────┬──────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
  ┌─────▼─────┐        ┌─────▼─────┐        ┌─────▼─────┐
  │ Public    │        │ Private   │        │ Private   │
  │ Subnet    │        │ Subnet    │        │ Subnet    │
  │ (ALB)     │        │ (App)     │        │ (DB)      │
  └───────────┘        └───────────┘        └───────────┘

Rule Types:
• Stateless: 5-tuple match (quick filtering)
• Stateful: Connection tracking, deep inspection
• Domain Lists: Allow/deny by FQDN
• Suricata IPS: Threat signatures
```

### 1.7 Transit Gateway Network Manager

```
Global Network Visualization
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    AWS Network Manager                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Global Network View:                                                │
│                                                                      │
│      ┌──────────┐              ┌──────────┐                         │
│      │ TGW      │──────────────│ TGW      │                         │
│      │ us-east-1│   Peering    │ eu-west-1│                         │
│      └────┬─────┘              └────┬─────┘                         │
│           │                         │                                │
│     ┌─────┴─────┐             ┌─────┴─────┐                         │
│     │ 5 VPCs    │             │ 3 VPCs    │                         │
│     │ 2 VPNs    │             │ 1 DX      │                         │
│     └───────────┘             └───────────┘                         │
│                                                                      │
│  Features:                                                           │
│  • Topology visualization                                           │
│  • Route analysis                                                   │
│  • CloudWatch metrics integration                                   │
│  • SD-WAN integration (third-party devices)                        │
│  • Events and alerts                                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Serverless Deep Dive

### 2.1 Serverless Computing Model

```
Serverless Characteristics
──────────────────────────────────────────────────────────────────────

✓ No server management         ✓ Pay-per-use pricing
✓ Auto-scaling (to zero)       ✓ Built-in high availability
✓ Event-driven execution       ✓ Managed security patches

AWS Serverless Services:
├── Compute: Lambda, Fargate
├── API: API Gateway, AppSync
├── Storage: S3, DynamoDB, Aurora Serverless
├── Integration: EventBridge, SQS, SNS, Step Functions
├── Analytics: Kinesis, Athena
└── Orchestration: Step Functions
```

### 2.2 AWS Lambda Deep Dive

**Execution Model**:

```
Lambda Execution Lifecycle
──────────────────────────────────────────────────────────────────────

Cold Start:
┌─────────────────────────────────────────────────────────────────────┐
│ Download   │  Start      │  Initialize    │  Handler    │          │
│ Code       │  Runtime    │  Extensions    │  Execution  │  Cleanup │
│            │  (Node/Py)  │  + Init Code   │             │          │
└────────────┴─────────────┴────────────────┴─────────────┴──────────┘
     │              │               │
     └──────────────┴───────────────┘
          Cold Start Latency
          (100ms - 10s depending on runtime/size)

Warm Start:
┌─────────────────────────────────────────────────────────────────────┐
│  Handler    │          │
│  Execution  │  Cleanup │
└─────────────┴──────────┘
     Execution context reused!


Execution Context Reuse:
• Runtime initialized once
• DB connections persist
• Global variables maintained
• /tmp directory preserved (512MB - 10GB)
```

**Memory and CPU Allocation**:

```
Lambda Resource Allocation
──────────────────────────────────────────────────────────────────────

Memory (MB)  │  vCPU Equivalent  │  Network Bandwidth
─────────────┼───────────────────┼────────────────────
    128      │     0.083         │       Low
    512      │     0.333         │       Low
   1024      │     0.583         │       Moderate
   1769      │     1.0 vCPU      │       Full
   3538      │     2.0 vCPU      │       Full
   10240     │     6.0 vCPU      │       Full

Key Insight:
At 1,769 MB, you get 1 full vCPU.
CPU scales linearly with memory allocation.
Network bandwidth also scales with memory.
```

**Lambda Optimization Strategies**:

```python
# COLD START OPTIMIZATION

# ❌ Bad: Import inside handler
def handler(event, context):
    import boto3  # Imported every invocation
    import pandas as pd
    # ...

# ✅ Good: Import outside handler
import boto3
import pandas as pd

# Initialize outside handler (reused in warm starts)
dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('my-table')

def handler(event, context):
    # Use pre-initialized resources
    response = table.get_item(Key={'id': event['id']})
    return response
```

```python
# PROVISIONED CONCURRENCY
# For latency-sensitive workloads

# AWS CLI
aws lambda put-provisioned-concurrency-config \
    --function-name my-function \
    --qualifier prod \
    --provisioned-concurrent-executions 100

# Effect: 100 execution environments kept warm
# Cost: ~$0.015/hour per provisioned instance
```

**Lambda Layers**:

```
Lambda Layers Architecture
──────────────────────────────────────────────────────────────────────

┌────────────────────────────────────────────┐
│              Lambda Function                │
│  ┌──────────────────────────────────────┐  │
│  │         Your Code (handler.py)        │  │
│  └──────────────────────────────────────┘  │
│                    │                        │
│  ┌─────────────────▼──────────────────────┐│
│  │              Layers                     ││
│  │  ┌──────────┐ ┌──────────┐ ┌────────┐ ││
│  │  │ Layer 1  │ │ Layer 2  │ │Layer 3 │ ││
│  │  │ (pandas) │ │ (numpy)  │ │(custom)│ ││
│  │  └──────────┘ └──────────┘ └────────┘ ││
│  └─────────────────────────────────────────┘│
│                    │                        │
│  ┌─────────────────▼──────────────────────┐│
│  │           Runtime (Python 3.12)         ││
│  └─────────────────────────────────────────┘│
└────────────────────────────────────────────┘

Benefits:
• Share libraries across functions
• Reduce deployment package size
• Version and update independently
• Up to 5 layers per function
• Total unzipped size: 250 MB
```

**Lambda Event Source Mappings**:

```
Event Source Types
──────────────────────────────────────────────────────────────────────

Synchronous (Push):                Asynchronous (Push):
• API Gateway                      • S3
• ALB                             • SNS
• Lambda Function URLs            • CloudWatch Events
• Cognito                         • EventBridge
                                  • SES
                                  • CloudFormation

Poll-based (Pull):                 Stream-based:
• SQS                             • Kinesis Data Streams
• DynamoDB Streams                • DynamoDB Streams
• Amazon MQ                       • Amazon MSK
• Kafka (self-managed)

Error Handling by Type:
──────────────────────────────────────────────────────────────────────
Sync:    Returns error to caller
Async:   2 retries, then DLQ/destination
Stream:  Retry until success or record expires
SQS:     Return to queue (visibility timeout), then DLQ
```

**Lambda Best Practices**:

```python
# IDEMPOTENCY - Handle duplicate invocations
import hashlib
from aws_lambda_powertools.utilities.idempotency import (
    idempotent,
    DynamoDBPersistenceLayer
)

persistence_layer = DynamoDBPersistenceLayer(table_name="IdempotencyTable")

@idempotent(persistence_store=persistence_layer)
def handler(event, context):
    # This will only execute once per unique event
    return process_order(event['order_id'])


# TIMEOUT HANDLING
import signal

def timeout_handler(signum, frame):
    raise Exception("Function timeout approaching!")

def handler(event, context):
    # Warning 5 seconds before timeout
    remaining_time = context.get_remaining_time_in_millis()
    signal.alarm(int((remaining_time - 5000) / 1000))
    
    try:
        # Your code here
        result = long_running_operation()
        return result
    finally:
        signal.alarm(0)  # Cancel the alarm
```

### 2.3 Amazon API Gateway

**API Gateway Types**:

| Feature | REST API | HTTP API | WebSocket API |
|---------|----------|----------|---------------|
| Cost | $3.50/million | $1.00/million | $1.00/million |
| Latency | ~29ms | ~10ms | N/A |
| Features | Full | Core | Real-time |
| Auth | IAM, Cognito, Lambda | IAM, Cognito, JWT | IAM, Lambda |
| Caching | Yes | No | No |
| WAF | Yes | No | Yes |
| Throttling | Yes | Yes | Yes |

**REST API Architecture**:

```
API Gateway REST API Components
──────────────────────────────────────────────────────────────────────

Client Request Flow:
─────────────────────

Request ──► Method Request ──► Integration Request ──► Backend
                │                     │
                │                     │
         ┌──────┴──────┐       ┌──────┴──────┐
         │ Validation  │       │ Mapping     │
         │ • Query     │       │ Template    │
         │ • Headers   │       │ (VTL)       │
         │ • Body      │       │             │
         └─────────────┘       └─────────────┘

Backend ──► Integration Response ──► Method Response ──► Response
                    │                      │
             ┌──────┴──────┐        ┌──────┴──────┐
             │ Mapping     │        │ Headers     │
             │ Template    │        │ Status Code │
             │             │        │             │
             └─────────────┘        └─────────────┘
```

**API Gateway Integration Types**:

```
Integration Types
──────────────────────────────────────────────────────────────────────

1. Lambda Proxy (Recommended)
   API Gateway ──► Lambda (full request/response passthrough)
   
2. Lambda Custom
   API Gateway ──► Transform ──► Lambda ──► Transform ──► Response
   
3. HTTP Proxy
   API Gateway ──► HTTP Endpoint (passthrough)
   
4. HTTP Custom
   API Gateway ──► Transform ──► HTTP ──► Transform ──► Response
   
5. AWS Service
   API Gateway ──► Direct AWS Service Call (SQS, SNS, Step Functions)
   
6. Mock
   API Gateway ──► Static Response (no backend)
```

**API Gateway Request Validation**:

```json
// OpenAPI 3.0 Specification with Validation
{
  "openapi": "3.0.0",
  "paths": {
    "/orders": {
      "post": {
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CreateOrder"
              }
            }
          }
        },
        "x-amazon-apigateway-request-validator": "all"
      }
    }
  },
  "components": {
    "schemas": {
      "CreateOrder": {
        "type": "object",
        "required": ["customerId", "items"],
        "properties": {
          "customerId": {
            "type": "string",
            "pattern": "^[A-Z]{2}[0-9]{6}$"
          },
          "items": {
            "type": "array",
            "minItems": 1,
            "items": {
              "$ref": "#/components/schemas/OrderItem"
            }
          }
        }
      }
    }
  },
  "x-amazon-apigateway-request-validators": {
    "all": {
      "validateRequestBody": true,
      "validateRequestParameters": true
    }
  }
}
```

**Throttling and Quotas**:

```
API Gateway Throttling Architecture
──────────────────────────────────────────────────────────────────────

Account Level Limits:
├── 10,000 requests/second (soft limit)
└── 5,000 concurrent executions (burstable)

Stage Level:
├── Default: 10,000 req/s
└── Burst: 5,000 requests

Method Level (Override):
├── GET /products: 1,000 req/s
└── POST /orders: 100 req/s

Usage Plans (Per API Key):
┌─────────────────────────────────────────────┐
│ Plan: "Enterprise"                          │
│ ├── Throttle: 1,000 req/s                  │
│ ├── Burst: 500 requests                     │
│ └── Quota: 1,000,000 requests/month        │
├─────────────────────────────────────────────┤
│ Plan: "Free Tier"                           │
│ ├── Throttle: 10 req/s                     │
│ ├── Burst: 10 requests                      │
│ └── Quota: 1,000 requests/day              │
└─────────────────────────────────────────────┘
```

### 2.4 AWS Step Functions

**State Machine Fundamentals**:

```
Step Functions State Types
──────────────────────────────────────────────────────────────────────

Task        │ Do work (Lambda, ECS, Activity, AWS SDK)
Choice      │ Branching logic (if/else)
Parallel    │ Execute branches concurrently
Map         │ Iterate over array items
Wait        │ Delay for time or until timestamp
Pass        │ Transform data or inject static data
Succeed     │ End with success
Fail        │ End with failure
```

**Workflow Types**:

| Feature | Standard | Express |
|---------|----------|---------|
| Duration | Up to 1 year | 5 minutes |
| Execution Rate | 2,000/sec | 100,000/sec |
| Pricing | Per state transition | Per execution + duration |
| History | Full (90 days) | CloudWatch Logs |
| Use Case | Long-running, auditable | High-volume, streaming |

**Complex Workflow Example**:

```json
{
  "Comment": "Order Processing Workflow",
  "StartAt": "ValidateOrder",
  "States": {
    "ValidateOrder": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:validateOrder",
      "Next": "CheckInventory",
      "Catch": [{
        "ErrorEquals": ["ValidationError"],
        "Next": "OrderFailed"
      }]
    },
    
    "CheckInventory": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:checkInventory",
      "Next": "InventoryAvailable?"
    },
    
    "InventoryAvailable?": {
      "Type": "Choice",
      "Choices": [
        {
          "Variable": "$.inventoryStatus",
          "StringEquals": "AVAILABLE",
          "Next": "ProcessPayment"
        },
        {
          "Variable": "$.inventoryStatus",
          "StringEquals": "BACKORDER",
          "Next": "NotifyBackorder"
        }
      ],
      "Default": "OrderFailed"
    },
    
    "ProcessPayment": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:processPayment",
      "Retry": [{
        "ErrorEquals": ["PaymentGatewayTimeout"],
        "IntervalSeconds": 3,
        "MaxAttempts": 3,
        "BackoffRate": 2
      }],
      "Next": "ParallelFulfillment"
    },
    
    "ParallelFulfillment": {
      "Type": "Parallel",
      "Branches": [
        {
          "StartAt": "UpdateInventory",
          "States": {
            "UpdateInventory": {
              "Type": "Task",
              "Resource": "arn:aws:lambda:us-east-1:123456789012:function:updateInventory",
              "End": true
            }
          }
        },
        {
          "StartAt": "SendConfirmation",
          "States": {
            "SendConfirmation": {
              "Type": "Task",
              "Resource": "arn:aws:lambda:us-east-1:123456789012:function:sendEmail",
              "End": true
            }
          }
        },
        {
          "StartAt": "InitiateShipping",
          "States": {
            "InitiateShipping": {
              "Type": "Task",
              "Resource": "arn:aws:lambda:us-east-1:123456789012:function:createShipment",
              "End": true
            }
          }
        }
      ],
      "Next": "OrderComplete"
    },
    
    "NotifyBackorder": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:notifyBackorder",
      "Next": "WaitForInventory"
    },
    
    "WaitForInventory": {
      "Type": "Wait",
      "Seconds": 86400,
      "Next": "CheckInventory"
    },
    
    "OrderComplete": {
      "Type": "Succeed"
    },
    
    "OrderFailed": {
      "Type": "Fail",
      "Error": "OrderProcessingFailed",
      "Cause": "Order could not be processed"
    }
  }
}
```

**Visual Representation**:

```
Order Processing State Machine
──────────────────────────────────────────────────────────────────────

                    ┌──────────────────┐
                    │  ValidateOrder   │
                    └────────┬─────────┘
                             │
                    ┌────────▼─────────┐
                    │  CheckInventory  │◄───────────────┐
                    └────────┬─────────┘                │
                             │                          │
                    ┌────────▼─────────┐                │
               ┌────│InventoryAvailable?│────┐          │
               │    └──────────────────┘    │          │
          AVAILABLE                    BACKORDER       │
               │                            │          │
        ┌──────▼──────┐           ┌─────────▼────────┐ │
        │ProcessPayment│           │ NotifyBackorder  │ │
        └──────┬──────┘           └─────────┬────────┘ │
               │                            │          │
        ┌──────▼──────────────┐   ┌─────────▼────────┐ │
        │ ParallelFulfillment │   │  WaitForInventory│─┘
        │  ├─ UpdateInventory │   │    (24 hours)    │
        │  ├─ SendConfirmation│   └──────────────────┘
        │  └─ InitiateShipping│
        └──────┬──────────────┘
               │
        ┌──────▼──────┐
        │ OrderComplete│
        │   (Succeed)  │
        └─────────────┘
```

### 2.5 Amazon EventBridge

**Event-Driven Architecture with EventBridge**:

```
EventBridge Architecture
──────────────────────────────────────────────────────────────────────

Event Sources                   Event Bus              Targets
─────────────                   ─────────              ───────

┌──────────────┐             ┌──────────────┐     ┌──────────────┐
│ AWS Services │────────────►│              │────►│   Lambda     │
│ (S3, EC2...) │             │              │     └──────────────┘
└──────────────┘             │              │     ┌──────────────┐
                             │   Default    │────►│  Step Funcs  │
┌──────────────┐             │   Event Bus  │     └──────────────┘
│ Custom Apps  │────────────►│              │     ┌──────────────┐
│ (PutEvents)  │             │              │────►│    SQS       │
└──────────────┘             │              │     └──────────────┘
                             └──────────────┘     ┌──────────────┐
┌──────────────┐                                  │    SNS       │
│ SaaS Partners│             ┌──────────────┐     └──────────────┘
│ (Zendesk,    │────────────►│  Partner     │     ┌──────────────┐
│  Datadog...) │             │  Event Bus   │────►│  API Gateway │
└──────────────┘             └──────────────┘     └──────────────┘

                             ┌──────────────┐     ┌──────────────┐
                             │   Custom     │────►│ Cross-Account│
                             │  Event Bus   │     │   Event Bus  │
                             └──────────────┘     └──────────────┘
```

**Event Pattern Matching**:

```json
// Match EC2 instance state changes
{
  "source": ["aws.ec2"],
  "detail-type": ["EC2 Instance State-change Notification"],
  "detail": {
    "state": ["running", "stopped", "terminated"]
  }
}

// Match custom events with content filtering
{
  "source": ["com.mycompany.orders"],
  "detail-type": ["OrderCreated"],
  "detail": {
    "orderValue": [{
      "numeric": [">=", 1000]
    }],
    "customerType": ["enterprise", "premium"],
    "region": [{
      "prefix": "us-"
    }]
  }
}

// Complex filtering with $or
{
  "source": ["com.mycompany.orders"],
  "$or": [
    { "detail": { "priority": ["high"] } },
    { "detail": { "orderValue": [{ "numeric": [">=", 5000] }] } }
  ]
}
```

**EventBridge Pipes**:

```
EventBridge Pipes - Point-to-Point Integration
──────────────────────────────────────────────────────────────────────

┌─────────┐    ┌──────────┐    ┌─────────────┐    ┌──────────┐
│ Source  │───►│ Filtering │───►│ Enrichment  │───►│  Target  │
│         │    │           │    │ (Optional)  │    │          │
└─────────┘    └──────────┘    └─────────────┘    └──────────┘

Sources:              Enrichment:           Targets:
• DynamoDB Streams    • Lambda              • Lambda
• Kinesis             • API Gateway         • Step Functions
• SQS                 • Step Functions      • API Gateway
• Amazon MQ           • API Destination     • Kinesis
• Kafka (MSK)                               • SQS, SNS
                                            • EventBridge

Use Case Example:
DynamoDB ──► Filter(INSERT) ──► Lambda(Enrich) ──► Step Functions
```

### 2.6 Serverless Patterns

**Pattern 1: API + Lambda + DynamoDB**

```
REST API Pattern
──────────────────────────────────────────────────────────────────────

┌──────────┐    ┌──────────────┐    ┌──────────┐    ┌──────────┐
│  Client  │───►│ API Gateway  │───►│  Lambda  │───►│ DynamoDB │
└──────────┘    │  (REST/HTTP) │    │          │    │          │
                └──────────────┘    └──────────┘    └──────────┘

Best Practices:
• Use HTTP API for cost (70% cheaper)
• Lambda function per HTTP method/resource
• Single-table DynamoDB design
• API Gateway caching for read-heavy
```

**Pattern 2: Event Processing Pipeline**

```
Event Processing Pattern
──────────────────────────────────────────────────────────────────────

┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│   S3    │───►│  Lambda │───►│   SQS   │───►│ Lambda  │
│ Upload  │    │ Trigger │    │  Queue  │    │ Process │
└─────────┘    └─────────┘    └─────────┘    └────┬────┘
                                                  │
                                    ┌─────────────▼─────────┐
                                    │        DynamoDB       │
                                    │     │    Kinesis      │
                                    │     │    S3 Output    │
                                    └─────────────────────────┘

Best Practices:
• SQS between Lambda functions for reliability
• Dead Letter Queue for failed messages
• Batch processing (10 messages per invocation)
• Idempotency for at-least-once delivery
```

**Pattern 3: Saga Pattern with Step Functions**

```
Saga Pattern for Distributed Transactions
──────────────────────────────────────────────────────────────────────

                  ┌──────────────────────────────────────┐
                  │         Step Functions               │
                  │                                      │
Forward Path:     │  ┌────────┐   ┌────────┐   ┌──────┐│
                  │  │Reserve │──►│ Charge │──►│ Ship ││
                  │  │Inventory│   │Payment │   │Order ││
                  │  └────┬───┘   └───┬────┘   └──┬───┘│
                  │       │           │           │     │
Compensating      │       ▼           ▼           ▼     │
Transactions:     │  ┌────────┐   ┌────────┐   ┌──────┐│
(on failure)      │  │Release │   │ Refund │   │Cancel││
                  │  │Inventory│   │Payment │   │Ship  ││
                  │  └────────┘   └────────┘   └──────┘│
                  │                                      │
                  └──────────────────────────────────────┘

Key Principles:
• Each step has a compensating action
• Steps are idempotent
• State machine orchestrates flow
• Failures trigger compensation
```

---

## 3. Container Services

### 3.1 Container Orchestration Options

```
AWS Container Spectrum
──────────────────────────────────────────────────────────────────────

More Control ◄──────────────────────────────────────► Less Management
                                                      
EKS on EC2      │  ECS on EC2      │  EKS Fargate   │  ECS Fargate
────────────────┼──────────────────┼────────────────┼───────────────
Full Kubernetes │  AWS-native      │  Kubernetes    │  Simplest
control         │  orchestration   │  without nodes │  container
                │                  │                │  platform
Manage: Nodes,  │  Manage: EC2     │  Manage: Pods  │  Manage:
K8s upgrades,   │  capacity        │  only          │  Tasks only
networking      │                  │                │
```

### 3.2 Amazon ECS (Elastic Container Service)

**ECS Architecture**:

```
ECS Core Components
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                          ECS Cluster                                 │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                        Service                               │   │
│  │                                                               │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │   │
│  │  │    Task      │  │    Task      │  │    Task      │       │   │
│  │  │  ┌────────┐  │  │  ┌────────┐  │  │  ┌────────┐  │       │   │
│  │  │  │Container│  │  │  │Container│  │  │  │Container│  │       │   │
│  │  │  └────────┘  │  │  └────────┘  │  └────────┘  │       │   │
│  │  │  ┌────────┐  │  │  ┌────────┐  │  │  ┌────────┐  │       │   │
│  │  │  │Sidecar │  │  │  │Sidecar │  │  │  │Sidecar │  │       │   │
│  │  │  └────────┘  │  │  └────────┘  │  │  └────────┘  │       │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘       │   │
│  │                                                               │   │
│  │  Task Definition: Container image, CPU, Memory, Ports,       │   │
│  │                   Environment, IAM Role, Logging             │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  Capacity Providers:                                                 │
│  ├── Fargate (Serverless)                                          │
│  ├── Fargate Spot (Cost savings)                                   │
│  └── EC2 Auto Scaling Group (Self-managed)                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Terminology:
• Cluster: Logical grouping of tasks/services
• Service: Maintains desired count of tasks
• Task: Running instance of task definition
• Task Definition: Blueprint (like docker-compose)
• Container: Individual container within a task
```

**Task Definition Example**:

```json
{
  "family": "web-app",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::123456789012:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "web",
      "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/web-app:latest",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "NODE_ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:123456789012:secret:db-url"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/web-app",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "web"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      }
    },
    {
      "name": "datadog-agent",
      "image": "datadog/agent:latest",
      "essential": false,
      "environment": [
        {"name": "DD_API_KEY", "value": "xxx"},
        {"name": "ECS_FARGATE", "value": "true"}
      ]
    }
  ]
}
```

**ECS Service Types**:

```
ECS Service Scheduling
──────────────────────────────────────────────────────────────────────

REPLICA Service:                    DAEMON Service:
• Maintains desired task count      • One task per EC2 instance
• Distributes across AZs            • Monitoring agents, log shippers
• Auto-replaces unhealthy tasks     • Not supported on Fargate

Service Configuration:
┌─────────────────────────────────────────────────────────────────────┐
│  service:                                                           │
│    desiredCount: 3                                                  │
│    deploymentConfiguration:                                         │
│      minimumHealthyPercent: 50    # At least 50% during deploy     │
│      maximumPercent: 200          # Can scale to 200% during deploy│
│    deploymentController:                                            │
│      type: ECS | CODE_DEPLOY | EXTERNAL                            │
│    placementStrategies:                                             │
│      - type: spread                                                 │
│        field: attribute:ecs.availability-zone                      │
│    capacityProviderStrategy:                                        │
│      - capacityProvider: FARGATE                                   │
│        weight: 1                                                   │
│        base: 2                    # Guaranteed Fargate tasks       │
│      - capacityProvider: FARGATE_SPOT                              │
│        weight: 4                  # 4x more Spot tasks             │
└─────────────────────────────────────────────────────────────────────┘
```

**ECS Auto Scaling**:

```
ECS Service Auto Scaling
──────────────────────────────────────────────────────────────────────

                    ┌────────────────────────────────┐
                    │     Application Auto Scaling   │
                    └───────────────┬────────────────┘
                                    │
         ┌──────────────────────────┼────────────────────────────┐
         │                          │                            │
         ▼                          ▼                            ▼
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│ Target Tracking │      │  Step Scaling   │      │   Scheduled     │
│                 │      │                 │      │                 │
│ CPU = 70%       │      │ CPU > 80% +2    │      │ Scale up at 9am │
│ Memory = 80%    │      │ CPU < 40% -1    │      │ Scale down 6pm  │
│ Request/target  │      │                 │      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘

# Terraform Example
resource "aws_appautoscaling_target" "ecs" {
  max_capacity       = 10
  min_capacity       = 2
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.main.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "ecs_cpu" {
  name               = "cpu-auto-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.ecs.resource_id
  scalable_dimension = aws_appautoscaling_target.ecs.scalable_dimension
  service_namespace  = aws_appautoscaling_target.ecs.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value = 70
  }
}
```

### 3.3 Amazon EKS (Elastic Kubernetes Service)

**EKS Architecture**:

```
EKS Architecture Overview
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                         EKS Cluster                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────┐                        │
│  │         Control Plane (AWS Managed)     │                        │
│  │                                          │                        │
│  │  ┌────────────┐  ┌────────────┐         │                        │
│  │  │ API Server │  │   etcd     │         │  Runs across 3 AZs    │
│  │  └────────────┘  └────────────┘         │  Automatic upgrades    │
│  │  ┌────────────┐  ┌────────────┐         │  Managed HA           │
│  │  │ Scheduler  │  │ Controller │         │                        │
│  │  │            │  │  Manager   │         │                        │
│  │  └────────────┘  └────────────┘         │                        │
│  │                                          │                        │
│  └─────────────────────────────────────────┘                        │
│                         │                                            │
│                    ┌────┴────┐                                       │
│                    │ EKS API │                                       │
│                    └────┬────┘                                       │
│                         │                                            │
│  Customer VPC           │                                            │
│  ┌──────────────────────┼──────────────────────────────────────┐   │
│  │                      │                                       │   │
│  │  ┌───────────────────▼───────────────────┐                  │   │
│  │  │           Data Plane Options          │                  │   │
│  │  ├───────────────────────────────────────┤                  │   │
│  │  │                                        │                  │   │
│  │  │  Managed Node Groups     │  Fargate   │                  │   │
│  │  │  (EC2 with ASG)          │  Profiles  │                  │   │
│  │  │                          │            │                  │   │
│  │  │  Self-Managed Nodes      │            │                  │   │
│  │  │  (Your EC2 + AMI)        │            │                  │   │
│  │  │                          │            │                  │   │
│  │  └───────────────────────────────────────┘                  │   │
│  │                                                              │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Node Group Types**:

| Type | Management | Use Case |
|------|------------|----------|
| **Managed Node Groups** | AWS manages ASG, updates | Standard workloads |
| **Self-Managed Nodes** | You manage everything | Custom AMIs, GPU |
| **Fargate** | Serverless, no nodes | Pod-centric, variable |
| **Karpenter** | Auto-provisioning | Flexible, cost-optimized |

**EKS Networking**:

```
EKS Networking Options
──────────────────────────────────────────────────────────────────────

VPC CNI Plugin (Default):
• Each pod gets VPC IP address
• Native VPC networking
• Security groups for pods
• Limitation: IP address exhaustion

┌─────────────────────────────────────────────────────────────────────┐
│                              VPC                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                          Subnet                               │   │
│  │                                                               │   │
│  │  Node (10.0.1.10)                                            │   │
│  │  ├── Pod-1 (10.0.1.11)  ← VPC IP                            │   │
│  │  ├── Pod-2 (10.0.1.12)  ← VPC IP                            │   │
│  │  └── Pod-3 (10.0.1.13)  ← VPC IP                            │   │
│  │                                                               │   │
│  └───────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘

IP Prefix Delegation (Recommended):
• Assigns /28 prefix to nodes
• 16 IPs per prefix
• Better IP utilization

Secondary CIDR / Custom Networking:
• Use separate CIDR for pods
• Preserve VPC IP space
```

**EKS Add-ons**:

```
Essential EKS Add-ons
──────────────────────────────────────────────────────────────────────

1. CoreDNS                 - Kubernetes DNS
2. kube-proxy             - Network proxy
3. VPC CNI                - Pod networking
4. AWS Load Balancer Ctrl - ALB/NLB ingress
5. EBS CSI Driver         - EBS volumes
6. EFS CSI Driver         - EFS volumes
7. Cluster Autoscaler     - Node scaling (or Karpenter)
8. Metrics Server         - Resource metrics

# Install AWS Load Balancer Controller
helm repo add eks https://aws.github.io/eks-charts
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=my-cluster \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller
```

**EKS IAM Integration (IRSA)**:

```yaml
# IAM Roles for Service Accounts (IRSA)

# 1. Create IAM Role with trust policy
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/XXXXXXXXXX"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "oidc.eks.us-east-1.amazonaws.com/id/XXXXXXXXXX:sub": "system:serviceaccount:default:my-service-account",
          "oidc.eks.us-east-1.amazonaws.com/id/XXXXXXXXXX:aud": "sts.amazonaws.com"
        }
      }
    }
  ]
}

# 2. Create Kubernetes Service Account
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account
  namespace: default
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/my-pod-role

# 3. Use in Pod
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: my-service-account
  containers:
    - name: my-container
      image: my-image
      # AWS SDK automatically uses IRSA credentials
```

### 3.4 Amazon ECR (Elastic Container Registry)

**ECR Features**:

```
ECR Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                            ECR Registry                              │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  Repository: my-app                                          │   │
│  │  ┌─────────────────────────────────────────────────────────┐│   │
│  │  │  Image: my-app:v1.0.0                                   ││   │
│  │  │  ├── Manifest (sha256:abc123)                          ││   │
│  │  │  ├── Config Layer                                       ││   │
│  │  │  └── Image Layers (cached/shared)                      ││   │
│  │  └─────────────────────────────────────────────────────────┘│   │
│  │  ┌─────────────────────────────────────────────────────────┐│   │
│  │  │  Image: my-app:v1.1.0                                   ││   │
│  │  │  └── Shares base layers with v1.0.0                    ││   │
│  │  └─────────────────────────────────────────────────────────┘│   │
│  │                                                              │   │
│  │  Features:                                                   │   │
│  │  • Vulnerability scanning (on push / continuous)            │   │
│  │  • Image signing (with AWS Signer)                          │   │
│  │  • Lifecycle policies                                        │   │
│  │  • Cross-region replication                                  │   │
│  │  • Cross-account access                                      │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ECR Public: Public container registry (like Docker Hub)            │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**ECR Lifecycle Policy**:

```json
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Keep last 10 production images",
      "selection": {
        "tagStatus": "tagged",
        "tagPrefixList": ["prod"],
        "countType": "imageCountMoreThan",
        "countNumber": 10
      },
      "action": {
        "type": "expire"
      }
    },
    {
      "rulePriority": 2,
      "description": "Delete untagged images older than 1 day",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 1
      },
      "action": {
        "type": "expire"
      }
    },
    {
      "rulePriority": 3,
      "description": "Keep last 5 dev images",
      "selection": {
        "tagStatus": "tagged",
        "tagPrefixList": ["dev"],
        "countType": "imageCountMoreThan",
        "countNumber": 5
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
```

### 3.5 Container Deployment Strategies

```
Deployment Strategies Comparison
──────────────────────────────────────────────────────────────────────

Rolling Update (Default):
──────────────────────────
v1 ████████████
       ↓ Replace 25% at a time
v1 ████████░░░░
v2 ░░░░░░░░████
       ↓
v2 ████████████

Pros: Zero downtime, gradual
Cons: Mixed versions during deploy


Blue-Green:
───────────
┌──────────┐     ┌──────────┐
│  Blue    │     │  Green   │
│  (v1)    │     │  (v2)    │
│  Active  │     │  Standby │
└────┬─────┘     └────┬─────┘
     │                │
     └───── LB ───────┘
           │
    Switch traffic instantly

Pros: Instant rollback, full testing
Cons: Double resources during deploy


Canary:
───────
Traffic:  90% ──────────────────▶  v1 (stable)
          10% ──────────────────▶  v2 (canary)

Gradually increase v2 traffic: 10% → 25% → 50% → 100%

Pros: Risk mitigation, real user testing
Cons: Complex traffic management


ECS + CodeDeploy (Blue-Green):
──────────────────────────────
┌─────────────────────────────────────────────────────────────────────┐
│                           ALB                                        │
│  ┌──────────────────────┐     ┌──────────────────────┐             │
│  │ Target Group Blue    │     │ Target Group Green   │             │
│  │ (Production Traffic) │     │ (Test Traffic)       │             │
│  └──────────┬───────────┘     └──────────┬───────────┘             │
│             │                             │                          │
│             ▼                             ▼                          │
│  ┌──────────────────────┐     ┌──────────────────────┐             │
│  │ ECS Service v1       │     │ ECS Service v2       │             │
│  │ (3 tasks)            │     │ (3 tasks)            │             │
│  └──────────────────────┘     └──────────────────────┘             │
│                                                                      │
│  Traffic shift options:                                             │
│  • AllAtOnce                                                        │
│  • Linear10PercentEvery1Minute                                      │
│  • Canary10Percent5Minutes                                          │
└─────────────────────────────────────────────────────────────────────┘
```

### 3.6 ECS vs EKS Decision Matrix

```
When to Choose What
──────────────────────────────────────────────────────────────────────

Choose ECS When:
✓ Team is new to containers
✓ AWS-native integration preferred
✓ Simpler operational model needed
✓ Tight IAM/CloudWatch integration required
✓ No existing Kubernetes expertise
✓ Smaller scale deployments

Choose EKS When:
✓ Existing Kubernetes expertise/workloads
✓ Multi-cloud/hybrid portability needed
✓ Rich ecosystem (Helm, operators) required
✓ Complex scheduling requirements
✓ Service mesh (Istio, Linkerd) needed
✓ Stronger community/third-party tooling

Choose Fargate When:
✓ Want to avoid node management
✓ Variable/unpredictable workloads
✓ Small to medium scale
✓ Development/test environments
✓ Batch processing with quick scale-out

Choose EC2 When:
✓ Need GPU instances
✓ Specific instance type requirements
✓ Cost optimization with Reserved Instances
✓ Need local storage (NVMe)
✓ Large, stable workloads
```

---

## 4. Observability & Monitoring

### 4.1 Observability Pillars on AWS

```
Three Pillars of Observability
──────────────────────────────────────────────────────────────────────

         METRICS                  LOGS                   TRACES
         ───────                  ────                   ──────
    CloudWatch Metrics      CloudWatch Logs           AWS X-Ray
    
    • Time-series data      • Event records           • Request flow
    • Aggregated values     • Searchable text         • Service map
    • Dashboards            • Patterns/Insights       • Latency analysis
    • Alarms                • Subscriptions           • Error tracking
    
         │                        │                        │
         └────────────────────────┼────────────────────────┘
                                  │
                        ┌─────────▼─────────┐
                        │  CloudWatch       │
                        │  ServiceLens      │
                        │  (Unified View)   │
                        └───────────────────┘
```

### 4.2 Amazon CloudWatch Deep Dive

**CloudWatch Metrics**:

```
CloudWatch Metrics Hierarchy
──────────────────────────────────────────────────────────────────────

Namespace (AWS/EC2, Custom)
    └── Metric Name (CPUUtilization)
        └── Dimensions (InstanceId=i-xxx, AutoScalingGroupName=xxx)
            └── Data Points (Timestamp, Value, Unit)

Standard Resolution: 1-minute (60 seconds)
High Resolution:     1-second (for custom metrics)

Retention:
• < 60 seconds:  3 hours
• 60 seconds:    15 days
• 300 seconds:   63 days
• 3600 seconds:  455 days (15 months)
```

**Custom Metrics Example**:

```python
import boto3
from datetime import datetime

cloudwatch = boto3.client('cloudwatch')

# Put custom metric
cloudwatch.put_metric_data(
    Namespace='MyApplication',
    MetricData=[
        {
            'MetricName': 'OrdersProcessed',
            'Dimensions': [
                {'Name': 'Environment', 'Value': 'Production'},
                {'Name': 'Service', 'Value': 'OrderProcessor'}
            ],
            'Timestamp': datetime.utcnow(),
            'Value': 150,
            'Unit': 'Count',
            'StorageResolution': 1  # High resolution (1 second)
        },
        {
            'MetricName': 'ProcessingLatency',
            'Dimensions': [
                {'Name': 'Environment', 'Value': 'Production'},
                {'Name': 'Service', 'Value': 'OrderProcessor'}
            ],
            'Timestamp': datetime.utcnow(),
            'Value': 245.5,
            'Unit': 'Milliseconds'
        }
    ]
)

# Using Embedded Metric Format (EMF) - Recommended for Lambda
import json

def handler(event, context):
    # This automatically creates CloudWatch metrics
    print(json.dumps({
        "_aws": {
            "Timestamp": int(datetime.utcnow().timestamp() * 1000),
            "CloudWatchMetrics": [{
                "Namespace": "MyApplication",
                "Dimensions": [["Service", "Operation"]],
                "Metrics": [
                    {"Name": "ProcessingTime", "Unit": "Milliseconds"},
                    {"Name": "RecordsProcessed", "Unit": "Count"}
                ]
            }]
        },
        "Service": "OrderService",
        "Operation": "CreateOrder",
        "ProcessingTime": 125,
        "RecordsProcessed": 10
    }))
```

**CloudWatch Alarms**:

```
CloudWatch Alarm Architecture
──────────────────────────────────────────────────────────────────────

Alarm States:
• OK          - Metric within threshold
• ALARM       - Metric breached threshold
• INSUFFICIENT_DATA - Not enough data points

Alarm Types:
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  1. Metric Alarm                                                    │
│     └── Single metric threshold                                     │
│                                                                      │
│  2. Composite Alarm                                                 │
│     └── Multiple alarms with AND/OR logic                          │
│                                                                      │
│  3. Anomaly Detection Alarm                                         │
│     └── ML-based dynamic threshold                                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Alarm Example (Terraform):
─────────────────────────

resource "aws_cloudwatch_metric_alarm" "high_cpu" {
  alarm_name          = "high-cpu-utilization"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "CPU utilization exceeds 80%"
  
  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.main.name
  }
  
  alarm_actions = [
    aws_sns_topic.alerts.arn,
    aws_autoscaling_policy.scale_up.arn
  ]
  
  ok_actions = [
    aws_sns_topic.alerts.arn
  ]
  
  treat_missing_data = "notBreaching"
}

# Composite Alarm
resource "aws_cloudwatch_composite_alarm" "service_health" {
  alarm_name = "service-health-composite"
  
  alarm_rule = "ALARM(${aws_cloudwatch_metric_alarm.high_cpu.alarm_name}) AND ALARM(${aws_cloudwatch_metric_alarm.high_memory.alarm_name})"
  
  alarm_actions = [aws_sns_topic.critical.arn]
}
```

**CloudWatch Logs**:

```
CloudWatch Logs Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        Log Group                                     │
│                        /aws/lambda/my-function                       │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  Log Stream: 2024/01/15/[$LATEST]abc123                      │   │
│  │  ┌─────────────────────────────────────────────────────────┐│   │
│  │  │ Log Event: {"timestamp": "...", "message": "..."}       ││   │
│  │  │ Log Event: {"timestamp": "...", "message": "..."}       ││   │
│  │  └─────────────────────────────────────────────────────────┘│   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  Settings:                                                           │
│  ├── Retention: 1 day to 10 years (or never expire)                │
│  ├── Encryption: AWS-managed or CMK                                 │
│  └── Metric Filters: Extract metrics from log patterns             │
│                                                                      │
│  Outputs:                                                            │
│  ├── Subscription Filter → Lambda, Kinesis, OpenSearch              │
│  ├── Export to S3 (batch)                                           │
│  └── CloudWatch Logs Insights (query)                               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**CloudWatch Logs Insights**:

```sql
-- Find errors with context
fields @timestamp, @message, @logStream
| filter @message like /ERROR|Exception/
| sort @timestamp desc
| limit 100

-- Analyze Lambda cold starts
fields @timestamp, @message, @duration, @billedDuration
| filter @type = "REPORT"
| filter @message like /Init Duration/
| parse @message /Init Duration: (?<initDuration>[\d.]+)/
| stats count(*) as coldStarts, 
        avg(initDuration) as avgInitDuration,
        max(initDuration) as maxInitDuration
  by bin(1h)

-- API Gateway latency analysis
fields @timestamp, @message
| filter @message like /API Gateway/
| parse @message /latency: (?<latency>[\d.]+)/
| stats percentile(latency, 50) as p50,
        percentile(latency, 95) as p95,
        percentile(latency, 99) as p99
  by bin(5m)

-- Error rate calculation
fields @timestamp, @message
| stats count(*) as total,
        sum(case when @message like /ERROR/ then 1 else 0 end) as errors
  by bin(1h)
| display total, errors, (errors * 100.0 / total) as errorRate

-- Find slow database queries
fields @timestamp, @message
| filter @message like /SELECT|INSERT|UPDATE|DELETE/
| parse @message /duration: (?<duration>[\d.]+)ms/
| filter duration > 1000
| sort duration desc
| limit 50
```

**Metric Filters**:

```terraform
# Extract metrics from logs
resource "aws_cloudwatch_log_metric_filter" "error_count" {
  name           = "ErrorCount"
  log_group_name = "/aws/lambda/my-function"
  pattern        = "[timestamp, requestId, level=ERROR, ...]"
  
  metric_transformation {
    name          = "ErrorCount"
    namespace     = "MyApplication"
    value         = "1"
    default_value = "0"
  }
}

# JSON pattern matching
resource "aws_cloudwatch_log_metric_filter" "latency" {
  name           = "ResponseLatency"
  log_group_name = "/aws/apigateway/my-api"
  pattern        = "{ $.latency > 0 }"
  
  metric_transformation {
    name      = "ResponseLatency"
    namespace = "MyApplication"
    value     = "$.latency"
    unit      = "Milliseconds"
  }
}
```

### 4.3 AWS X-Ray

**Distributed Tracing Concepts**:

```
X-Ray Trace Structure
──────────────────────────────────────────────────────────────────────

Trace (entire request journey)
│
├── Segment: API Gateway
│   ├── Subsegment: Request validation
│   └── Subsegment: Lambda invocation
│
├── Segment: Lambda Function
│   ├── Subsegment: Initialization
│   ├── Subsegment: Handler execution
│   │   ├── Subsegment: DynamoDB GetItem
│   │   └── Subsegment: S3 PutObject
│   └── Subsegment: Response
│
└── Segment: DynamoDB
    └── Subsegment: GetItem operation

Annotations: Key-value pairs for filtering (indexed)
Metadata: Non-indexed data for context
```

**X-Ray Integration Example**:

```python
# Lambda with X-Ray
from aws_xray_sdk.core import xray_recorder
from aws_xray_sdk.core import patch_all

# Automatically instrument boto3, requests, etc.
patch_all()

@xray_recorder.capture('process_order')
def process_order(order_id):
    # Add annotation for filtering
    xray_recorder.put_annotation('order_id', order_id)
    xray_recorder.put_annotation('customer_tier', 'premium')
    
    # Add metadata for debugging
    xray_recorder.put_metadata('order_details', {
        'items': 5,
        'total': 299.99
    })
    
    # Subsegments for detailed tracing
    with xray_recorder.in_subsegment('validate_order'):
        validate(order_id)
    
    with xray_recorder.in_subsegment('charge_payment'):
        charge_payment(order_id)
    
    with xray_recorder.in_subsegment('update_inventory'):
        update_inventory(order_id)

def handler(event, context):
    order_id = event['order_id']
    return process_order(order_id)
```

**X-Ray Service Map**:

```
X-Ray Service Map Visualization
──────────────────────────────────────────────────────────────────────

         ┌─────────────┐
         │   Client    │
         │ Requests/s  │
         │ Avg latency │
         └──────┬──────┘
                │
         ┌──────▼──────┐
         │ API Gateway │──────────┐
         │   500 req/s │          │
         │   45ms avg  │          │ Errors shown
         └──────┬──────┘          │ in red
                │                 │
         ┌──────▼──────┐          │
         │   Lambda    │          │
         │   Order     │◄─────────┘
         │   Handler   │
         └──────┬──────┘
                │
     ┌──────────┼──────────┐
     │          │          │
┌────▼────┐ ┌───▼───┐ ┌────▼────┐
│DynamoDB │ │   S3  │ │  SQS    │
│  25ms   │ │  10ms │ │  5ms    │
└─────────┘ └───────┘ └─────────┘

Features:
• Latency distribution
• Error rates
• Request counts
• Service dependencies
• Drill-down to traces
```

### 4.4 AWS CloudTrail

**CloudTrail Event Types**:

```
CloudTrail Event Categories
──────────────────────────────────────────────────────────────────────

1. Management Events (Control Plane)
   └── API calls that manage resources
   └── Examples: CreateBucket, RunInstances, CreateUser
   └── Default: Enabled

2. Data Events (Data Plane)
   └── API calls on resources
   └── Examples: GetObject, PutItem, Invoke
   └── Default: Disabled (high volume)

3. Insights Events
   └── Unusual API activity detection
   └── ML-based anomaly detection
   └── Examples: Spike in TerminateInstances calls

Organization Trail:
└── Single trail for all accounts in AWS Organization
└── Centralized logging to management account S3
```

**CloudTrail Event Structure**:

```json
{
    "eventVersion": "1.08",
    "userIdentity": {
        "type": "AssumedRole",
        "principalId": "AROAXXXXXXXXXXXXXXXXX:i-1234567890abcdef0",
        "arn": "arn:aws:sts::123456789012:assumed-role/EC2Role/i-1234567890abcdef0",
        "accountId": "123456789012",
        "sessionContext": {
            "sessionIssuer": {
                "type": "Role",
                "principalId": "AROAXXXXXXXXXXXXXXXXX",
                "arn": "arn:aws:iam::123456789012:role/EC2Role",
                "accountId": "123456789012",
                "userName": "EC2Role"
            }
        }
    },
    "eventTime": "2024-01-15T10:30:45Z",
    "eventSource": "s3.amazonaws.com",
    "eventName": "GetObject",
    "awsRegion": "us-east-1",
    "sourceIPAddress": "10.0.1.50",
    "userAgent": "aws-sdk-python/1.26.0",
    "requestParameters": {
        "bucketName": "my-bucket",
        "key": "sensitive/data.json"
    },
    "responseElements": null,
    "requestID": "XXXXXXXXXX",
    "eventID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "readOnly": true,
    "resources": [
        {
            "type": "AWS::S3::Object",
            "ARN": "arn:aws:s3:::my-bucket/sensitive/data.json"
        }
    ],
    "eventType": "AwsApiCall",
    "managementEvent": false,
    "recipientAccountId": "123456789012"
}
```

**CloudTrail Lake (Querying)**:

```sql
-- Find all root user activity
SELECT eventTime, eventName, sourceIPAddress, userAgent
FROM cloudtrail_events
WHERE userIdentity.type = 'Root'
  AND eventTime > '2024-01-01'
ORDER BY eventTime DESC

-- Security: Find access from unusual IPs
SELECT eventTime, eventName, userIdentity.arn, sourceIPAddress
FROM cloudtrail_events
WHERE sourceIPAddress NOT LIKE '10.%'
  AND sourceIPAddress NOT LIKE '172.16.%'
  AND eventTime > '2024-01-14'

-- Audit: Track IAM changes
SELECT eventTime, eventName, userIdentity.arn, requestParameters
FROM cloudtrail_events
WHERE eventSource = 'iam.amazonaws.com'
  AND eventName LIKE '%User%' OR eventName LIKE '%Role%' OR eventName LIKE '%Policy%'
ORDER BY eventTime DESC
LIMIT 100
```

### 4.5 Observability Architecture Pattern

```
Enterprise Observability Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                      Data Sources                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │   EC2   │ │   ECS   │ │   EKS   │ │ Lambda  │ │   RDS   │       │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │
│       │           │           │           │           │             │
└───────┼───────────┼───────────┼───────────┼───────────┼─────────────┘
        │           │           │           │           │
        └───────────┴───────────┴───────────┴───────────┘
                              │
┌─────────────────────────────┼───────────────────────────────────────┐
│                             │  Collection Layer                      │
├─────────────────────────────┼───────────────────────────────────────┤
│                             │                                        │
│  ┌──────────────────────────▼──────────────────────────────────┐   │
│  │    CloudWatch Agent / ADOT (AWS Distro for OpenTelemetry)   │   │
│  │    ├── Metrics → CloudWatch Metrics                          │   │
│  │    ├── Logs → CloudWatch Logs                                │   │
│  │    └── Traces → X-Ray                                        │   │
│  └────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────┼───────────────────────────────────────┐
│                             │  Processing & Storage                  │
├─────────────────────────────┼───────────────────────────────────────┤
│                             │                                        │
│  ┌───────────────┐  ┌──────▼──────┐  ┌───────────────┐             │
│  │ CloudWatch    │  │ Kinesis     │  │ CloudWatch     │             │
│  │ Logs Insights │  │ Data        │  │ Metrics        │             │
│  │               │  │ Firehose    │  │ (Time Series)  │             │
│  └───────┬───────┘  └──────┬──────┘  └───────┬───────┘             │
│          │                 │                  │                      │
│          │          ┌──────▼──────┐          │                      │
│          │          │ S3 (Archive)│          │                      │
│          │          │ Athena Query│          │                      │
│          │          └─────────────┘          │                      │
│          │                                   │                      │
└──────────┼───────────────────────────────────┼──────────────────────┘
           │                                   │
┌──────────┼───────────────────────────────────┼──────────────────────┐
│          │  Visualization & Alerting         │                      │
├──────────┼───────────────────────────────────┼──────────────────────┤
│          │                                   │                      │
│  ┌───────▼────────────────────────────────────▼───────┐            │
│  │              CloudWatch Dashboards                  │            │
│  │              ┌─────────────────────────────────┐   │            │
│  │              │  ┌─────┐ ┌─────┐ ┌─────┐ ┌────┐│   │            │
│  │              │  │Graph│ │Logs │ │Trace│ │Alrm││   │            │
│  │              │  └─────┘ └─────┘ └─────┘ └────┘│   │            │
│  │              └─────────────────────────────────┘   │            │
│  └────────────────────────────────────────────────────┘            │
│                               │                                     │
│                        ┌──────▼──────┐                              │
│                        │ SNS/Lambda  │──► PagerDuty, Slack         │
│                        │ Alerting    │                              │
│                        └─────────────┘                              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 4.6 Amazon Managed Grafana & Prometheus

```
Managed Observability Stack
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  Amazon Managed Prometheus (AMP)                                     │
│  ├── PromQL queries                                                 │
│  ├── Kubernetes metrics (via remote_write)                          │
│  ├── 150-day retention                                              │
│  └── Multi-AZ HA                                                    │
│                                                                      │
│  Amazon Managed Grafana (AMG)                                        │
│  ├── AWS data source plugins                                        │
│  │   ├── CloudWatch                                                 │
│  │   ├── X-Ray                                                      │
│  │   ├── Prometheus                                                 │
│  │   └── OpenSearch                                                 │
│  ├── SSO integration                                                │
│  └── Managed dashboards                                             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Use Case: EKS Monitoring
────────────────────────

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    EKS       │────►│ Prometheus   │────►│   Grafana    │
│  Cluster     │     │ (AMP)        │     │   (AMG)      │
└──────────────┘     └──────────────┘     └──────────────┘
       │
       │ ADOT Collector
       │ (OpenTelemetry)
       │
┌──────▼──────┐
│   X-Ray     │
│  (Traces)   │
└─────────────┘
```

---

## 5. Infrastructure as Code

### 5.1 IaC Tools Comparison

```
Infrastructure as Code Landscape
──────────────────────────────────────────────────────────────────────

              AWS-Native                    Multi-Cloud
              ──────────                    ───────────
              
Declarative   CloudFormation               Terraform
              (JSON/YAML)                  (HCL)
              
Imperative    AWS CDK                      Pulumi
              (TypeScript, Python, etc.)   (TypeScript, Python, etc.)
              
              SAM (Serverless)             CrossPlane
              (CloudFormation extension)   (Kubernetes-native)

Selection Criteria:
┌─────────────────┬──────────────────────────────────────────────────┐
│ Choose          │ When                                              │
├─────────────────┼──────────────────────────────────────────────────┤
│ CloudFormation  │ AWS-only, enterprise compliance, GovCloud        │
│ AWS CDK         │ Developer-friendly, complex logic, AWS-focused   │
│ Terraform       │ Multi-cloud, mature ecosystem, team expertise    │
│ Pulumi          │ General-purpose languages, existing SDK skills   │
└─────────────────┴──────────────────────────────────────────────────┘
```

### 5.2 AWS CloudFormation

**Template Structure**:

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Description: 'Production VPC with ECS Cluster'

# Input parameters
Parameters:
  Environment:
    Type: String
    Default: prod
    AllowedValues: [dev, staging, prod]
  VpcCidr:
    Type: String
    Default: '10.0.0.0/16'
    
# Conditional resource creation
Conditions:
  IsProd: !Equals [!Ref Environment, 'prod']
  
# Mappings for environment-specific values
Mappings:
  RegionAMI:
    us-east-1:
      HVM64: ami-0c55b159cbfafe1f0
    us-west-2:
      HVM64: ami-0a1b2c3d4e5f67890

# Resources to create
Resources:
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: !Ref VpcCidr
      EnableDnsHostnames: true
      EnableDnsSupport: true
      Tags:
        - Key: Name
          Value: !Sub '${Environment}-vpc'
        - Key: Environment
          Value: !Ref Environment
          
  PublicSubnet1:
    Type: AWS::EC2::Subnet
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [0, !Cidr [!Ref VpcCidr, 4, 8]]
      AvailabilityZone: !Select [0, !GetAZs '']
      MapPublicIpOnLaunch: true
      Tags:
        - Key: Name
          Value: !Sub '${Environment}-public-1'
          
  # Conditional resource
  ProductionAlarm:
    Type: AWS::CloudWatch::Alarm
    Condition: IsProd
    Properties:
      AlarmName: !Sub '${Environment}-high-cpu'
      MetricName: CPUUtilization
      Namespace: AWS/EC2
      Threshold: 80
      ComparisonOperator: GreaterThanThreshold
      EvaluationPeriods: 2
      Period: 300
      Statistic: Average

# Output values for cross-stack references
Outputs:
  VpcId:
    Description: VPC ID
    Value: !Ref VPC
    Export:
      Name: !Sub '${Environment}-VpcId'
      
  PublicSubnet1Id:
    Description: Public Subnet 1 ID
    Value: !Ref PublicSubnet1
    Export:
      Name: !Sub '${Environment}-PublicSubnet1Id'
```

**CloudFormation Intrinsic Functions**:

```yaml
# Reference Functions
!Ref ResourceName                    # Returns resource ID
!GetAtt Resource.Attribute           # Get resource attribute

# String Functions
!Sub '${AWS::StackName}-resource'    # String substitution
!Join ['-', [prefix, !Ref Env]]      # Join strings
!Split ['/', !Ref S3Path]            # Split string

# Conditional Functions
!If [ConditionName, TrueValue, FalseValue]
!Not [!Equals [!Ref Env, 'prod']]
!And [!Condition A, !Condition B]
!Or [!Condition A, !Condition B]

# List Functions
!Select [0, !GetAZs '']              # Select from list
!Cidr [!Ref VpcCidr, 4, 8]           # Generate CIDR blocks

# Import/Export
!ImportValue ExportedValueName       # Cross-stack reference

# Base64
!Base64 |
  #!/bin/bash
  yum update -y
```

**Nested Stacks Pattern**:

```yaml
# Parent Stack
Resources:
  NetworkStack:
    Type: AWS::CloudFormation::Stack
    Properties:
      TemplateURL: https://s3.amazonaws.com/mybucket/network.yaml
      Parameters:
        Environment: !Ref Environment
        
  DatabaseStack:
    Type: AWS::CloudFormation::Stack
    DependsOn: NetworkStack
    Properties:
      TemplateURL: https://s3.amazonaws.com/mybucket/database.yaml
      Parameters:
        VpcId: !GetAtt NetworkStack.Outputs.VpcId
        SubnetIds: !GetAtt NetworkStack.Outputs.PrivateSubnetIds
        
  ApplicationStack:
    Type: AWS::CloudFormation::Stack
    DependsOn: [NetworkStack, DatabaseStack]
    Properties:
      TemplateURL: https://s3.amazonaws.com/mybucket/application.yaml
      Parameters:
        VpcId: !GetAtt NetworkStack.Outputs.VpcId
        DatabaseEndpoint: !GetAtt DatabaseStack.Outputs.Endpoint
```

**StackSets for Multi-Account/Region**:

```
CloudFormation StackSets
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Management Account                                │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                     StackSet                                 │   │
│  │                     (Template)                               │   │
│  └───────────────────────────┬─────────────────────────────────┘   │
│                              │                                       │
└──────────────────────────────┼───────────────────────────────────────┘
                               │
          ┌────────────────────┼────────────────────┐
          │                    │                    │
          ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Account A       │  │ Account B       │  │ Account C       │
│ ├── us-east-1   │  │ ├── us-east-1   │  │ ├── us-east-1   │
│ │   (Stack)     │  │ │   (Stack)     │  │ │   (Stack)     │
│ └── eu-west-1   │  │ └── eu-west-1   │  │ └── eu-west-1   │
│     (Stack)     │  │     (Stack)     │  │     (Stack)     │
└─────────────────┘  └─────────────────┘  └─────────────────┘

Use Cases:
• Deploy security baselines across all accounts
• Set up logging/monitoring in all regions
• Compliance controls via SCPs
```

### 5.3 AWS CDK (Cloud Development Kit)

**CDK Concepts**:

```
CDK Architecture
──────────────────────────────────────────────────────────────────────

Your Code (TypeScript/Python/...)
            │
            ▼
    ┌─────────────────┐
    │   CDK App       │
    │                 │
    │  ┌───────────┐  │
    │  │  Stack 1  │  │
    │  │ ┌───────┐ │  │
    │  │ │Construct│ │  │  L1: CloudFormation resources (Cfn*)
    │  │ │ (L1-L3)│ │  │  L2: Higher-level abstractions
    │  │ └───────┘ │  │  L3: Patterns (complete solutions)
    │  └───────────┘  │
    │  ┌───────────┐  │
    │  │  Stack 2  │  │
    │  └───────────┘  │
    └────────┬────────┘
             │
             │ cdk synth
             ▼
    CloudFormation Template
             │
             │ cdk deploy
             ▼
    AWS Resources
```

**CDK Application Example (TypeScript)**:

```typescript
import * as cdk from 'aws-cdk-lib';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ecs_patterns from 'aws-cdk-lib/aws-ecs-patterns';
import * as rds from 'aws-cdk-lib/aws-rds';
import { Construct } from 'constructs';

// Props interface for reusable stack
interface ApplicationStackProps extends cdk.StackProps {
  environment: string;
  vpcCidr?: string;
}

export class ApplicationStack extends cdk.Stack {
  public readonly vpc: ec2.Vpc;
  public readonly cluster: ecs.Cluster;
  
  constructor(scope: Construct, id: string, props: ApplicationStackProps) {
    super(scope, id, props);

    // L2 Construct: VPC with best practices
    this.vpc = new ec2.Vpc(this, 'VPC', {
      ipAddresses: ec2.IpAddresses.cidr(props.vpcCidr || '10.0.0.0/16'),
      maxAzs: 3,
      natGateways: props.environment === 'prod' ? 3 : 1,
      subnetConfiguration: [
        {
          name: 'Public',
          subnetType: ec2.SubnetType.PUBLIC,
          cidrMask: 24,
        },
        {
          name: 'Private',
          subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
          cidrMask: 24,
        },
        {
          name: 'Isolated',
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
          cidrMask: 24,
        }
      ]
    });

    // ECS Cluster
    this.cluster = new ecs.Cluster(this, 'Cluster', {
      vpc: this.vpc,
      containerInsights: true,
    });

    // Aurora Serverless v2 Database
    const database = new rds.DatabaseCluster(this, 'Database', {
      engine: rds.DatabaseClusterEngine.auroraPostgres({
        version: rds.AuroraPostgresEngineVersion.VER_15_4,
      }),
      serverlessV2MinCapacity: 0.5,
      serverlessV2MaxCapacity: 8,
      writer: rds.ClusterInstance.serverlessV2('writer'),
      readers: [
        rds.ClusterInstance.serverlessV2('reader', {
          scaleWithWriter: true,
        }),
      ],
      vpc: this.vpc,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
      },
    });

    // L3 Construct: Fargate Service with ALB (Pattern)
    const service = new ecs_patterns.ApplicationLoadBalancedFargateService(
      this,
      'WebService',
      {
        cluster: this.cluster,
        cpu: 512,
        memoryLimitMiB: 1024,
        desiredCount: props.environment === 'prod' ? 3 : 1,
        taskImageOptions: {
          image: ecs.ContainerImage.fromAsset('./src/web'),
          containerPort: 8080,
          environment: {
            NODE_ENV: props.environment,
            DATABASE_HOST: database.clusterEndpoint.hostname,
          },
          secrets: {
            DATABASE_PASSWORD: ecs.Secret.fromSecretsManager(
              database.secret!,
              'password'
            ),
          },
        },
        publicLoadBalancer: true,
        circuitBreaker: { rollback: true },
      }
    );

    // Auto-scaling
    const scaling = service.service.autoScaleTaskCount({
      minCapacity: props.environment === 'prod' ? 3 : 1,
      maxCapacity: 20,
    });

    scaling.scaleOnCpuUtilization('CpuScaling', {
      targetUtilizationPercent: 70,
    });

    scaling.scaleOnRequestCount('RequestScaling', {
      requestsPerTarget: 1000,
      targetGroup: service.targetGroup,
    });

    // Allow service to connect to database
    database.connections.allowDefaultPortFrom(service.service);

    // Outputs
    new cdk.CfnOutput(this, 'LoadBalancerDNS', {
      value: service.loadBalancer.loadBalancerDnsName,
    });
  }
}

// App entry point
const app = new cdk.App();

new ApplicationStack(app, 'DevStack', {
  environment: 'dev',
  env: { account: '123456789012', region: 'us-east-1' },
});

new ApplicationStack(app, 'ProdStack', {
  environment: 'prod',
  vpcCidr: '10.1.0.0/16',
  env: { account: '987654321098', region: 'us-east-1' },
});
```

**CDK Custom Constructs**:

```typescript
// Reusable construct for secure S3 bucket
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as kms from 'aws-cdk-lib/aws-kms';

export interface SecureBucketProps {
  bucketName?: string;
  expirationDays?: number;
}

export class SecureBucket extends Construct {
  public readonly bucket: s3.Bucket;
  public readonly key: kms.Key;

  constructor(scope: Construct, id: string, props: SecureBucketProps = {}) {
    super(scope, id);

    // CMK for encryption
    this.key = new kms.Key(this, 'Key', {
      enableKeyRotation: true,
      description: `Key for ${id}`,
    });

    // Secure bucket with all best practices
    this.bucket = new s3.Bucket(this, 'Bucket', {
      bucketName: props.bucketName,
      encryption: s3.BucketEncryption.KMS,
      encryptionKey: this.key,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
      enforceSSL: true,
      versioned: true,
      lifecycleRules: [
        {
          id: 'ExpireOldVersions',
          noncurrentVersionExpiration: cdk.Duration.days(
            props.expirationDays || 90
          ),
        },
        {
          id: 'IntelligentTiering',
          transitions: [
            {
              storageClass: s3.StorageClass.INTELLIGENT_TIERING,
              transitionAfter: cdk.Duration.days(30),
            },
          ],
        },
      ],
      serverAccessLogsBucket: this.createAccessLogsBucket(),
    });
  }

  private createAccessLogsBucket(): s3.Bucket {
    return new s3.Bucket(this, 'AccessLogsBucket', {
      encryption: s3.BucketEncryption.S3_MANAGED,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
      lifecycleRules: [
        {
          expiration: cdk.Duration.days(365),
        },
      ],
    });
  }
}
```

### 5.4 Terraform for AWS

**Terraform Project Structure**:

```
terraform-project/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars
│   ├── staging/
│   └── prod/
├── modules/
│   ├── vpc/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── ecs/
│   └── rds/
└── shared/
    ├── providers.tf
    └── backend.tf
```

**Terraform Module Example**:

```hcl
# modules/vpc/main.tf

variable "environment" {
  type        = string
  description = "Environment name"
}

variable "vpc_cidr" {
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  type    = list(string)
  default = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

locals {
  public_subnets  = [for i, az in var.availability_zones : cidrsubnet(var.vpc_cidr, 8, i)]
  private_subnets = [for i, az in var.availability_zones : cidrsubnet(var.vpc_cidr, 8, i + 10)]
}

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "${var.environment}-vpc"
    Environment = var.environment
  }
}

resource "aws_subnet" "public" {
  count                   = length(var.availability_zones)
  vpc_id                  = aws_vpc.main.id
  cidr_block              = local.public_subnets[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name        = "${var.environment}-public-${count.index + 1}"
    Environment = var.environment
    Type        = "public"
  }
}

resource "aws_subnet" "private" {
  count             = length(var.availability_zones)
  vpc_id            = aws_vpc.main.id
  cidr_block        = local.private_subnets[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = {
    Name        = "${var.environment}-private-${count.index + 1}"
    Environment = var.environment
    Type        = "private"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "${var.environment}-igw"
    Environment = var.environment
  }
}

resource "aws_nat_gateway" "main" {
  count         = var.environment == "prod" ? length(var.availability_zones) : 1
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id

  tags = {
    Name        = "${var.environment}-nat-${count.index + 1}"
    Environment = var.environment
  }
}

resource "aws_eip" "nat" {
  count  = var.environment == "prod" ? length(var.availability_zones) : 1
  domain = "vpc"

  tags = {
    Name        = "${var.environment}-nat-eip-${count.index + 1}"
    Environment = var.environment
  }
}

# Route tables
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name        = "${var.environment}-public-rt"
    Environment = var.environment
  }
}

resource "aws_route_table" "private" {
  count  = length(var.availability_zones)
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = var.environment == "prod" ? aws_nat_gateway.main[count.index].id : aws_nat_gateway.main[0].id
  }

  tags = {
    Name        = "${var.environment}-private-rt-${count.index + 1}"
    Environment = var.environment
  }
}

resource "aws_route_table_association" "public" {
  count          = length(var.availability_zones)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count          = length(var.availability_zones)
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# Outputs
output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_ids" {
  value = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  value = aws_subnet.private[*].id
}
```

**Terraform State Management**:

```hcl
# backend.tf - Remote state with locking

terraform {
  backend "s3" {
    bucket         = "my-terraform-state-bucket"
    key            = "environments/prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}

# State locking table
resource "aws_dynamodb_table" "terraform_lock" {
  name         = "terraform-state-lock"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
```

**Terraform Workspaces vs Environments**:

```
Workspace Pattern (Simple):
──────────────────────────────────────────────────────────────────────

terraform workspace new dev
terraform workspace new staging
terraform workspace new prod

# In code:
locals {
  environment = terraform.workspace
  
  instance_count = {
    dev     = 1
    staging = 2
    prod    = 3
  }
}

resource "aws_instance" "app" {
  count = local.instance_count[local.environment]
  # ...
}


Directory Pattern (Recommended for production):
──────────────────────────────────────────────────────────────────────

environments/
├── dev/
│   └── terraform.tfvars
├── staging/
│   └── terraform.tfvars
└── prod/
    └── terraform.tfvars

# Each environment has its own state file
# Better isolation and blast radius control
```

### 5.5 IaC Best Practices

```
Infrastructure as Code Best Practices
──────────────────────────────────────────────────────────────────────

1. VERSION CONTROL
   ├── All IaC in Git
   ├── Pull request reviews
   ├── Branch protection
   └── Semantic versioning for modules

2. STATE MANAGEMENT
   ├── Remote state (S3, Terraform Cloud)
   ├── State locking (DynamoDB)
   ├── Encrypt state at rest
   └── Never commit state files

3. SECRETS HANDLING
   ├── Use AWS Secrets Manager / Parameter Store
   ├── Never hardcode secrets
   ├── Reference secrets by ARN
   └── Rotate secrets regularly

4. TESTING
   ├── Validate syntax (terraform validate, cfn-lint)
   ├── Security scanning (checkov, tfsec, cfn-nag)
   ├── Integration tests (terratest)
   └── Cost estimation (Infracost)

5. CI/CD PIPELINE
   ┌─────────────────────────────────────────────────────────────┐
   │  Commit → Lint → Security Scan → Plan → Review → Apply     │
   └─────────────────────────────────────────────────────────────┘

6. TAGGING STRATEGY
   ├── Environment
   ├── Owner/Team
   ├── Cost Center
   ├── Application
   └── Managed-by: terraform/cloudformation

7. MODULARIZATION
   ├── DRY (Don't Repeat Yourself)
   ├── Reusable modules
   ├── Version pinning
   └── Document interfaces
```

**Security Scanning Example**:

```yaml
# GitHub Actions workflow for Terraform
name: Terraform CI

on:
  pull_request:
    paths:
      - 'terraform/**'

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
        
      - name: Terraform Format Check
        run: terraform fmt -check -recursive
        
      - name: Terraform Init
        run: terraform init
        
      - name: Terraform Validate
        run: terraform validate
        
      - name: tfsec Security Scan
        uses: aquasecurity/tfsec-action@v1.0.0
        
      - name: Checkov Security Scan
        uses: bridgecrewio/checkov-action@master
        with:
          directory: terraform/
          framework: terraform
          
      - name: Infracost
        uses: infracost/actions/setup@v2
        
      - name: Generate Infracost Report
        run: |
          infracost breakdown --path=. --format=json > infracost.json
          infracost comment github --path=infracost.json
          
      - name: Terraform Plan
        run: terraform plan -out=tfplan
        
      - name: Upload Plan
        uses: actions/upload-artifact@v3
        with:
          name: tfplan
          path: tfplan
```

---

## 6. Interview Questions & Scenarios

### 6.1 Advanced Networking Questions

**Q1: How would you design a network for 100+ VPCs across multiple accounts and regions?**

**Answer**:

```
Architecture: Hub-and-Spoke with Transit Gateway

1. Regional Transit Gateways
   • One TGW per region in a central networking account
   • Share via RAM to spoke accounts

2. Multi-Region Connectivity
   • TGW peering between regions for global reach
   • Or use AWS Cloud WAN for automated global mesh

3. Segmentation via Route Tables
   • Separate route tables for Prod, Non-Prod, Shared Services
   • Control traffic flow between environments

4. On-Premise Connectivity
   • Direct Connect to primary region
   • VPN backup to secondary region
   • Advertise routes via BGP

5. Centralized Egress
   • Single VPC with NAT Gateways/Firewalls
   • All outbound traffic routes through inspection

Benefits:
• Centralized management
• Cost optimization (shared Direct Connect)
• Security inspection at chokepoint
• Simplified troubleshooting
```

---

**Q2: Your application in a private subnet can't reach S3. Walk through your troubleshooting process.**

**Answer**:

```
Troubleshooting Checklist:

1. VPC Endpoint Check
   □ Is there a Gateway Endpoint for S3?
   □ Is it associated with the route table?
   
2. Route Table Check
   □ Route for S3 prefix list → vpce-xxx
   □ If NAT: Route 0.0.0.0/0 → nat-xxx
   
3. Security Group Check
   □ Outbound rule allows HTTPS (443)
   □ For Gateway Endpoint: allows S3 prefix list
   
4. NACL Check
   □ Allows outbound 443
   □ Allows inbound ephemeral ports (1024-65535)
   
5. S3 Bucket Policy Check
   □ Not denying VPC/VPCE access
   □ Check aws:sourceVpc condition
   
6. IAM Permissions Check
   □ Instance role has s3:* permissions
   □ Not blocked by SCP
   
7. DNS Resolution Check
   □ s3.amazonaws.com resolves correctly
   □ VPC has enableDnsSupport=true

Diagnostic Commands:
aws ec2 describe-vpc-endpoints --filters Name=vpc-id,Values=vpc-xxx
aws ec2 describe-route-tables --route-table-ids rtb-xxx
nslookup s3.amazonaws.com
curl -v https://s3.amazonaws.com
```

---

### 6.2 Serverless Scenario Questions

**Q3: Design a serverless video processing pipeline that handles 10,000 uploads/hour.**

**Answer**:

```
Architecture:
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│   S3 (Upload)                                                        │
│       │                                                              │
│       │ S3 Event Notification                                       │
│       ▼                                                              │
│   ┌─────────┐                                                        │
│   │   SQS   │ ← Decoupling + Batching                               │
│   │  Queue  │   (1,000 messages buffer)                             │
│   └────┬────┘                                                        │
│        │                                                             │
│        │ Event Source Mapping (batch: 10)                           │
│        ▼                                                             │
│   ┌───────────────┐                                                  │
│   │    Lambda     │ ← Start MediaConvert Job                        │
│   │  (Initiator)  │   Concurrency: 100                              │
│   └───────┬───────┘                                                  │
│           │                                                          │
│           ▼                                                          │
│   ┌───────────────┐                                                  │
│   │ MediaConvert  │ ← Managed transcoding service                   │
│   │   (Jobs)      │   Parallel processing                           │
│   └───────┬───────┘                                                  │
│           │                                                          │
│           │ CloudWatch Event (Job Complete)                         │
│           ▼                                                          │
│   ┌───────────────┐                                                  │
│   │  EventBridge  │ ← Route to post-processing                      │
│   └───────┬───────┘                                                  │
│           │                                                          │
│     ┌─────┴─────┐                                                    │
│     ▼           ▼                                                    │
│ ┌─────────┐ ┌─────────┐                                             │
│ │ Lambda  │ │ Lambda  │ ← Update DB, notifications                  │
│ │(Metadata)│ │(Notify) │                                            │
│ └────┬────┘ └────┬────┘                                             │
│      │           │                                                   │
│      ▼           ▼                                                   │
│   DynamoDB     SNS                                                   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Key Design Decisions:
1. SQS buffer absorbs burst traffic
2. Lambda reserved concurrency prevents throttling
3. MediaConvert handles heavy lifting (no Lambda timeouts)
4. EventBridge decouples completion handling
5. DLQ for failed processing

Scaling:
• 10,000 uploads/hour = ~3/second
• SQS batching: 300 Lambda invocations/hour
• MediaConvert: On-demand queues auto-scale
```

---

**Q4: How do you handle cold starts for a latency-sensitive API?**

**Answer**:

```
Cold Start Mitigation Strategies:

1. PROVISIONED CONCURRENCY (Most Effective)
   • Keep N environments warm
   • Cost: ~$0.015/hour per concurrent execution
   • Best for predictable traffic patterns

2. ARCHITECTURE OPTIMIZATION
   • Minimize package size
   • Use arm64 (Graviton) - faster cold starts
   • Lazy load dependencies
   • Initialize outside handler

3. RUNTIME SELECTION
   Fastest → Slowest:
   • Go, Rust, C# AOT
   • Python, Node.js
   • Java, .NET

4. CODE OPTIMIZATION
   // Bad
   def handler(event, context):
       import pandas  # 3-5s import time
   
   // Good  
   import pandas  # Reused in warm starts
   df_processor = Processor()  # Pre-initialized

5. KEEP-WARM STRATEGY (If Provisioned Concurrency too expensive)
   • CloudWatch Events every 5 minutes
   • Only effective for single-concurrency scenarios

6. HTTP API vs REST API
   • HTTP API: ~10ms overhead
   • REST API: ~30ms overhead

Recommended Architecture:
┌─────────────────────────────────────────┐
│ CloudFront → HTTP API → Lambda          │
│                (Provisioned: 50)        │
│                                         │
│ + Application Auto Scaling for         │
│   Provisioned Concurrency              │
└─────────────────────────────────────────┘
```

---

### 6.3 Container Platform Questions

**Q5: ECS vs EKS - what factors drive your decision?**

**Answer**:

```
Decision Matrix:

CHOOSE ECS WHEN:
┌─────────────────────────────────────────────────────────────────────┐
│ ✓ AWS-native integration is priority                                │
│ ✓ Team lacks Kubernetes expertise                                   │
│ ✓ Simpler operational model preferred                               │
│ ✓ Tight IAM/CloudWatch integration needed                          │
│ ✓ Fewer services, straightforward networking                       │
│ ✓ Cost optimization (no control plane cost)                        │
└─────────────────────────────────────────────────────────────────────┘

CHOOSE EKS WHEN:
┌─────────────────────────────────────────────────────────────────────┐
│ ✓ Existing Kubernetes expertise/workloads                          │
│ ✓ Multi-cloud/hybrid cloud strategy                                │
│ ✓ Need Kubernetes ecosystem (Helm, Operators, CRDs)                │
│ ✓ Complex scheduling requirements                                   │
│ ✓ Service mesh (Istio, Linkerd) needed                             │
│ ✓ Strong open-source community tooling                             │
│ ✓ GitOps workflows (ArgoCD, Flux)                                  │
└─────────────────────────────────────────────────────────────────────┘

COST COMPARISON (100 tasks/pods, Fargate):
─────────────────────────────────────────────────────────────────────
ECS:  $0/month control plane + compute
EKS:  $73/month control plane + compute

OPERATIONAL COMPARISON:
─────────────────────────────────────────────────────────────────────
ECS:  AWS manages orchestrator
EKS:  AWS manages control plane, you manage upgrades, add-ons

My Recommendation:
• New containerization journey → ECS (simpler)
• Existing Kubernetes workloads → EKS
• Multi-cloud requirement → EKS
• Serverless preference → Either with Fargate
```

---

### 6.4 Observability Questions

**Q6: Design a monitoring strategy for a microservices application with 50+ services.**

**Answer**:

```
Observability Architecture:
──────────────────────────────────────────────────────────────────────

METRICS (CloudWatch + Prometheus)
├── Golden Signals per Service:
│   • Latency (p50, p95, p99)
│   • Traffic (requests/second)
│   • Errors (error rate %)
│   • Saturation (CPU, Memory)
│
├── Business Metrics:
│   • Orders processed
│   • Revenue generated
│   • User sign-ups
│
└── Infrastructure Metrics:
    • Node/container health
    • Network I/O
    • Disk utilization

LOGS (CloudWatch Logs + OpenSearch)
├── Structured Logging (JSON):
│   {
│     "timestamp": "...",
│     "service": "order-service",
│     "trace_id": "abc123",
│     "level": "ERROR",
│     "message": "...",
│     "metadata": {...}
│   }
│
├── Log Aggregation:
│   • CloudWatch Logs per service
│   • Subscription filter → OpenSearch
│   • 30-day hot, S3 archive
│
└── Log Insights Dashboards:
    • Error patterns
    • Latency outliers

TRACES (X-Ray / OpenTelemetry)
├── Distributed Tracing:
│   • Trace ID propagation across services
│   • Segment per service
│   • Subsegments for external calls
│
├── Service Map:
│   • Dependency visualization
│   • Error rate by service
│   • Latency distribution
│
└── Sampling Strategy:
    • 100% for errors
    • 5% for successful requests
    • Dynamic rate adjustment

ALERTING STRATEGY:
┌─────────────────────────────────────────────────────────────────────┐
│ Severity   │ Condition              │ Response Time │ Channel      │
├────────────┼────────────────────────┼───────────────┼──────────────┤
│ Critical   │ Service down           │ 5 min         │ PagerDuty    │
│ Warning    │ Error rate > 1%        │ 15 min        │ Slack        │
│ Info       │ Latency p99 > SLA      │ 1 hour        │ Ticket       │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 6.5 IaC Scenario Questions

**Q7: How would you structure Terraform for a multi-account AWS environment?**

**Answer**:

```
Directory Structure:
──────────────────────────────────────────────────────────────────────

terraform/
├── modules/                    # Reusable modules
│   ├── vpc/
│   ├── eks/
│   ├── rds/
│   └── security-baseline/
│
├── accounts/                   # Account-level config
│   ├── management/
│   │   ├── organizations.tf
│   │   └── sso.tf
│   ├── shared-services/
│   │   ├── network-hub.tf
│   │   └── dns.tf
│   ├── security/
│   │   ├── guardduty.tf
│   │   └── config-rules.tf
│   └── workloads/
│       ├── dev/
│       ├── staging/
│       └── prod/
│
├── governance/                 # Account-wide policies
│   ├── scps/
│   ├── tag-policies/
│   └── iam-boundaries/
│
└── pipelines/                  # CI/CD definitions
    ├── account-setup.yml
    └── workload-deploy.yml


State Management:
──────────────────────────────────────────────────────────────────────

Option 1: Separate State per Account/Region
accounts/dev/us-east-1/terraform.tfstate
accounts/dev/eu-west-1/terraform.tfstate
accounts/prod/us-east-1/terraform.tfstate

Option 2: Terraform Cloud / Enterprise
• Workspaces per account-region
• VCS-driven workflows
• Policy as Code (Sentinel)


Access Pattern:
──────────────────────────────────────────────────────────────────────

                 ┌───────────────────┐
                 │  Automation Acct  │
                 │  (GitHub Actions) │
                 └─────────┬─────────┘
                           │ assume role
            ┌──────────────┼──────────────┐
            ▼              ▼              ▼
       ┌─────────┐   ┌─────────┐   ┌─────────┐
       │   Dev   │   │ Staging │   │  Prod   │
       │ Account │   │ Account │   │ Account │
       └─────────┘   └─────────┘   └─────────┘

# Trust policy in each account
{
  "Principal": {
    "AWS": "arn:aws:iam::AUTOMATION_ACCOUNT:role/terraform-role"
  },
  "Action": "sts:AssumeRole",
  "Condition": {
    "StringEquals": {
      "aws:PrincipalTag/team": "platform"
    }
  }
}
```

---

**Q8: How do you handle secrets in Infrastructure as Code?**

**Answer**:

```
Secrets Management Strategies:
──────────────────────────────────────────────────────────────────────

1. AWS SECRETS MANAGER (Recommended)
   ┌─────────────────────────────────────────────────────────────┐
   │ # Terraform Reference                                       │
   │ data "aws_secretsmanager_secret_version" "db" {            │
   │   secret_id = "prod/database/credentials"                   │
   │ }                                                           │
   │                                                             │
   │ resource "aws_rds_cluster" "main" {                        │
   │   master_password = jsondecode(                            │
   │     data.aws_secretsmanager_secret_version.db.secret_string│
   │   )["password"]                                            │
   │ }                                                           │
   └─────────────────────────────────────────────────────────────┘

2. SSM PARAMETER STORE (Cost-effective)
   ┌─────────────────────────────────────────────────────────────┐
   │ data "aws_ssm_parameter" "db_password" {                   │
   │   name            = "/prod/db/password"                    │
   │   with_decryption = true                                   │
   │ }                                                           │
   └─────────────────────────────────────────────────────────────┘

3. NEVER DO THIS:
   ✗ Hardcode secrets in .tf files
   ✗ Store in terraform.tfvars (even gitignored)
   ✗ Pass via environment variables in CI/CD logs
   ✗ Store in state file unencrypted

4. BEST PRACTICES:
   ├── Pre-create secrets outside Terraform
   ├── Reference by ARN/name in IaC
   ├── Enable automatic rotation
   ├── Use IAM conditions for access
   └── Audit access via CloudTrail

5. CI/CD SECRET INJECTION:
   ┌─────────────────────────────────────────────────────────────┐
   │ # GitHub Actions with OIDC                                  │
   │ - name: Configure AWS credentials                           │
   │   uses: aws-actions/configure-aws-credentials@v4           │
   │   with:                                                     │
   │     role-to-assume: arn:aws:iam::123456789012:role/GHA     │
   │     aws-region: us-east-1                                   │
   │     # No secrets stored in GitHub!                         │
   └─────────────────────────────────────────────────────────────┘
```

---

## Summary: Key Takeaways for Part 2

```
┌────────────────────────────────────────────────────────────────────┐
│           Platform Engineer's Checklist - Part 2                   │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ □ Advanced Networking                                              │
│   • Transit Gateway for multi-VPC at scale                        │
│   • Direct Connect + VPN for hybrid connectivity                  │
│   • Network Firewall for centralized inspection                   │
│                                                                     │
│ □ Serverless                                                       │
│   • Lambda optimization (layers, provisioned concurrency)         │
│   • API Gateway selection (HTTP vs REST)                          │
│   • Step Functions for workflow orchestration                     │
│   • EventBridge for event-driven architecture                     │
│                                                                     │
│ □ Containers                                                       │
│   • ECS for AWS-native, EKS for Kubernetes ecosystem             │
│   • Fargate for serverless, EC2 for control                       │
│   • IRSA for secure pod-level AWS access                         │
│   • GitOps deployment patterns                                    │
│                                                                     │
│ □ Observability                                                    │
│   • CloudWatch (Metrics, Logs, Alarms)                           │
│   • X-Ray for distributed tracing                                 │
│   • CloudTrail for audit and compliance                          │
│   • Golden signals: Latency, Traffic, Errors, Saturation         │
│                                                                     │
│ □ Infrastructure as Code                                           │
│   • CloudFormation for AWS-native compliance                      │
│   • CDK for developer-friendly abstractions                       │
│   • Terraform for multi-cloud and mature ecosystem               │
│   • Security scanning in CI/CD pipeline                          │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

---

## Next Steps

This completes **Part 2: Advanced Services & Platform Engineering**.

**Part 3** will cover:
- Multi-Account Strategy & AWS Organizations
- Migration Strategies (6 Rs, AWS Migration Hub)
- Disaster Recovery & Business Continuity
- Security Deep Dive (GuardDuty, Security Hub, IAM Access Analyzer)
- Cost Optimization at Scale
- Real-world Architecture Case Studies
- DevOps & CI/CD on AWS

---

*Last Updated: February 2026*

