# Lambda



Def

> AWS Lambda is a serverless compute service that executes code in response to events without requiring developers to provision or manage servers. AWS automatically handles infrastructure, scaling, availability, patching, and capacity management. You pay only for the actual execution time of your code.



> AWS Lambda is a serverless compute service that allows us to run code without provisioning or managing servers. Instead of maintaining EC2 instances, we upload our code, configure an event source, and AWS automatically executes it when the event occurs. Lambda scales automatically based on demand, and billing is based only on the number of invocations and execution duration.
>
> In production, Lambda is commonly used for event-driven workloads such as image processing after S3 uploads, sending emails from SQS messages, scheduled jobs using EventBridge, webhook processing, and lightweight APIs behind API Gateway. It integrates well with services like S3, DynamoDB, SNS, and SQS.
>
> One important consideration is the cold start, especially for Java applications, where the first invocation may take longer because the runtime and JVM need to initialize. For long-running services or complex Spring Boot microservices that require stable low latency, EC2, ECS, or EKS are often more appropriate. Lambda is best suited for short-lived, stateless, event-driven workloads that benefit from automatic scaling and pay-per-use pricing.



The first request may take longer because AWS has to:

- Allocate infrastructure
- Start the runtime
- Load your code
- Initialize dependencies

#### Use Case 1: Image Processing (Very Common)

This is probably the most common Lambda interview example.

###### Problem

Users upload profile pictures.

Without Lambda

```
User
   |
Spring Boot
   |
Resize Image
   |
Save Image
```

Problems:

-  API becomes slow 
-  High CPU usage 
-  Large uploads can block request threads 

Instead

```
User
   |
Spring Boot
   |
Upload Original Image
   |
S3
   |
S3 Event
   |
Lambda
   |
Resize Image
   |
Save Thumbnail
```

Now your API responds in milliseconds.

---

#### Spring Boot Code

Upload original image.

```
@PostMapping("/upload")
public String upload(MultipartFile file){

    String key = UUID.randomUUID()+".jpg";

    s3Client.putObject(
            PutObjectRequest.builder()
                    .bucket("profile-images")
                    .key(key)
                    .build(),
            RequestBody.fromBytes(file.getBytes()));

    return key;
}
```

Notice

No resize logic.

No thumbnail generation.

---

### Lambda Code

```
public class ImageResizeLambda implements RequestHandler<S3Event,String>{

    @Override
    public String handleRequest(S3Event event, Context context){
        S3EventNotificationRecord record = event.getRecords().get(0);
        
        String bucket = record.getS3().getBucket().getName();

        String key = record.getS3().getObject().getKey();

        // Download image

        // Resize

        // Upload thumbnail

        return "done";
    }
}
```

---

Production benefit

-  API latency decreases 
-  Image processing scales automatically 
-  Failures are isolated

#### Production Perspective

> In production, I use Lambda for event-driven and asynchronous workloads rather than core business APIs. For example, in an e-commerce application, the Order Service saves the order and publishes an event to SQS. A Lambda function consumes the message, generates the invoice, uploads it to S3, sends an email using Amazon SES, and updates the order status. This keeps the API response fast while allowing independent scaling and automatic retries. Similarly, we use Lambda for image resizing after S3 uploads, scheduled cleanup jobs triggered by EventBridge, webhook processing through API Gateway, and log processing from CloudWatch. For long-running REST APIs or complex microservices, I prefer Spring Boot on ECS or EKS, while Lambda is ideal for short-lived, stateless, event-driven tasks because it scales automatically and follows a pay-per-use model.



#### ✅ Lambda

For

- Image processing
- Thumbnail generation
- PDF generation
- Email sending
- Notification
- Scheduled cleanup
- Event processing
- Webhooks
- Log analysis
- Data transformation
- ETL pipelines



## SQS

> Amazon SQS is a fully managed message queue that enables asynchronous communication between services. Instead of making synchronous calls, a producer publishes a message to an SQS queue, and one or more consumers process it independently. This decouples services, improves scalability, and makes the system more resilient to downstream failures.
>
> In a production e-commerce application, after an order is saved in the database, the Order Service publishes a message to an SQS queue rather than directly sending emails or generating invoices. Background workers or Lambda functions consume those messages to send confirmation emails, generate PDFs, update analytics, or notify other systems. This keeps the API response time low while allowing consumers to scale independently.
>
> SQS provides important reliability features such as visibility timeout, automatic retries, and Dead Letter Queues for messages that repeatedly fail processing. For most workloads we use Standard queues because they offer very high throughput. When strict ordering and deduplication are required, such as in certain financial workflows, we use FIFO queues.
>
> In short, SQS is the preferred choice for reliable, asynchronous background processing in AWS-based microservice architectures.

#### When do we use SQS?

Common production use cases:

- Order processing
- Email sending
- SMS notifications
- Invoice generation
- Image processing
- Video conversion
- Payment reconciliation
- Audit logging
- Inventory updates
- Background jobs

> SQS does not automatically know whether my business logic has succeeded. When a consumer receives a message, SQS makes it temporarily invisible using the visibility timeout but does not delete it. After my application finishes processing successfully, it explicitly calls the `DeleteMessage` API using the message's `ReceiptHandle`. If processing fails or the application crashes before calling `DeleteMessage`, the visibility timeout expires and SQS makes the message visible again so another consumer can retry it. This design provides **at-least-once delivery**, so consumers should be **idempotent** because the same message may be delivered more than once.



- **SQS = One producer → One consumer (work queue)**
- **SNS = One producer → Many consumers (publish/subscribe)**

## SNS

> Amazon SNS (Simple Notification Service) is a fully managed publish-subscribe messaging service. A publisher sends a message to an SNS Topic, and SNS automatically delivers that message to all subscribed endpoints such as SQS queues, Lambda functions, HTTP endpoints, email, SMS, or mobile push notifications.

#### Supported Subscribers

SNS can publish to

- SQS
- Lambda
- HTTP/HTTPS
- Email
- SMS
- Mobile Push Notifications

#### Why SNS + SQS Together? This is used heavily in production.

SNS broadcasts.

SQS guarantees reliable processing.

> Amazon SNS is a fully managed publish-subscribe messaging service used to broadcast events to multiple subscribers. Instead of directly calling several downstream services, a producer publishes a message to an SNS topic, and SNS delivers that message to all subscribed endpoints such as SQS queues, Lambda functions, HTTP endpoints, email, or SMS.
>
> In a production e-commerce system, after an order is successfully created, the Order Service publishes an `OrderCreated` event to an SNS topic. Different services subscribe independently—for example, the Email Service sends a confirmation email, the Inventory Service reserves stock, the Analytics Service records business metrics, and the Loyalty Service awards reward points. The Order Service remains completely unaware of these consumers, making the architecture loosely coupled and easy to extend.
>
> In practice, we often combine SNS with SQS. SNS performs the fan-out by broadcasting the event to multiple SQS queues, and each consumer processes messages from its own queue with retries and Dead Letter Queue support. This isolates failures—for example, if the Email Service is down, only the email queue is affected while inventory and analytics continue processing normally.



**Q: Why not connect SNS directly to Lambda or HTTP endpoints?**

A strong answer is:

> Direct SNS subscriptions are fine for lightweight consumers. However, for critical business workflows, we usually place an SQS queue between SNS and the consumer. The queue provides durability, buffering during traffic spikes, retries, visibility timeouts, and Dead Letter Queues. This prevents message loss and isolates slow or temporarily unavailable consumers, making the system more resilient.



### SQS vs SNS

SQS and SNS solve different messaging problems. SQS is used for asynchronous work distribution where one message is processed by one consumer. It acts as a durable queue that stores messages until they are successfully processed. SNS, on the other hand, is a publish-subscribe service used to broadcast the same event to multiple subscribers. A producer publishes a message to an SNS topic, and SNS delivers that message to all subscribed endpoints such as SQS queues, Lambda functions, email, SMS, or HTTP endpoints.

In production, we often use them together. SNS performs the fan-out by broadcasting an event to multiple SQS queues, and each service consumes messages from its own queue independently. This creates a loosely coupled, scalable, and fault-tolerant architecture.





### Decoupling with AWS SQS

> Decoupling with Amazon SQS means removing the direct dependency between two services by introducing a message queue. Instead of one service calling another synchronously and waiting for a response, it publishes a message to SQS. The consumer service processes that message independently at its own pace.
>
> For example, in an e-commerce application, after an order is created, the Order Service saves the order and publishes an `OrderCreated` message to an SQS queue. The Email Service later consumes the message and sends the confirmation email. If the Email Service is temporarily unavailable, the order is still created successfully because the message remains in the queue until it is processed. This improves fault tolerance, absorbs traffic spikes, and allows producers and consumers to scale independently. That's why SQS is commonly used to decouple microservices in production systems.





# Interview Ready Answer (3–4 minutes)

> Amazon ECS and Amazon EKS are both container orchestration services used to deploy and manage Docker containers in production. ECS is AWS's native orchestration platform, while EKS is AWS's managed Kubernetes service.
>
> In ECS, we define a Task Definition that specifies the container image, CPU, memory, ports, and environment variables. ECS Services ensure the desired number of tasks are always running and support auto scaling, rolling deployments, and health checks. ECS is generally simpler to configure and is a good choice for applications that run exclusively on AWS.
>
> EKS provides a managed Kubernetes control plane. Applications are deployed using Kubernetes resources such as Deployments, Services, and Pods. EKS is preferred when an organization already uses Kubernetes or requires Kubernetes features and portability across cloud providers.
>
> In both cases, the typical workflow is to build a Docker image, push it to Amazon ECR, and deploy it through ECS or EKS. The orchestrator then handles scheduling, scaling, health monitoring, and rolling updates, allowing teams to focus on the application rather than managing individual containers.

