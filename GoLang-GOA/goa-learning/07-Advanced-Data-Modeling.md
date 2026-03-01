# Part 7: Advanced Data Modeling

> **Goal:** Master advanced Goa DSL data modeling - Result Types, Views, Type Reuse, Polymorphism, Inheritance Patterns, and Struct Composition

---

## 📚 Table of Contents

1. [Data Modeling Overview](#data-modeling-overview)
2. [Result Types](#result-types)
3. [Views](#views)
4. [Type Reuse](#type-reuse)
5. [Polymorphism](#polymorphism)
6. [Inheritance Patterns](#inheritance-patterns)
7. [Struct Composition](#struct-composition)
8. [Advanced Patterns](#advanced-patterns)
9. [Complete Examples](#complete-examples)
10. [Summary](#summary)
11. [Knowledge Check](#knowledge-check)

---

## 🎯 Data Modeling Overview

### Why Advanced Data Modeling?

In real-world APIs, simple types aren't enough. You need:

```
┌─────────────────────────────────────────────────────────────────┐
│                 DATA MODELING CHALLENGES                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. DIFFERENT VIEWS                                             │
│     Same data, different contexts                               │
│     • List view (summary) vs Detail view (full)                 │
│     • Public view vs Admin view                                 │
│     • Mobile view vs Desktop view                               │
│                                                                 │
│  2. CODE REUSE                                                  │
│     Don't repeat yourself                                       │
│     • Common fields across types                                │
│     • Shared validation rules                                   │
│     • Consistent patterns                                       │
│                                                                 │
│  3. FLEXIBLE RESPONSES                                          │
│     Multiple return possibilities                               │
│     • Success or Error                                          │
│     • Different response structures                             │
│     • Conditional fields                                        │
│                                                                 │
│  4. TYPE RELATIONSHIPS                                          │
│     Model real-world entities                                   │
│     • Base types with variants                                  │
│     • Shared behaviors                                          │
│     • Type hierarchies                                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Goa Data Modeling Constructs

| Construct | Purpose | Use Case |
|-----------|---------|----------|
| `Type` | Define data structures | Basic type definition |
| `ResultType` | Define response structures with views | API responses |
| `View` | Define projections of types | Different response shapes |
| `Extend` | Add fields to existing types | Type extension |
| `Reference` | Reuse attribute definitions | Code reuse |
| `ArrayOf` | Collections | Lists |
| `MapOf` | Key-value structures | Dictionaries |

---

## 📤 Result Types

### What are Result Types?

Result Types are special types designed for API responses. They support:
- Multiple **views** (projections)
- Built-in **links** support
- **Content type** specification

```go
// design/types.go
package design

import . "goa.design/goa/v3/dsl"

// Basic Type vs Result Type
var _ = Type("CreateUserPayload", func() {
    // Regular type - used for inputs
    Attribute("name", String)
    Attribute("email", String)
    Required("name", "email")
})

var _ = ResultType("application/vnd.user", func() {
    // Result type - used for outputs
    // Has a media type identifier
    Description("A user resource")
    
    Attributes(func() {
        Attribute("id", Int64, "User ID")
        Attribute("name", String, "User name")
        Attribute("email", String, "Email address")
        Attribute("created_at", String, "Creation timestamp", func() {
            Format(FormatDateTime)
        })
        Attribute("updated_at", String, "Last update timestamp", func() {
            Format(FormatDateTime)
        })
        Attribute("profile", UserProfile, "User profile details")
        Attribute("settings", UserSettings, "User settings")
    })
    
    // Views define different projections
    View("default", func() {
        Attribute("id")
        Attribute("name")
        Attribute("email")
        Attribute("created_at")
        Attribute("updated_at")
        Attribute("profile")
        Attribute("settings")
    })
    
    View("tiny", func() {
        Attribute("id")
        Attribute("name")
    })
    
    View("summary", func() {
        Attribute("id")
        Attribute("name")
        Attribute("email")
    })
    
    Required("id", "name", "email")
})
```

### Result Type Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    RESULT TYPE STRUCTURE                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ResultType("application/vnd.user")                             │
│       │                                                         │
│       ├── TypeName: "User"                                      │
│       │                                                         │
│       ├── ContentType: "application/vnd.user"                   │
│       │                                                         │
│       ├── Attributes                                            │
│       │       ├── id (Int64)                                    │
│       │       ├── name (String)                                 │
│       │       ├── email (String)                                │
│       │       ├── created_at (String)                           │
│       │       ├── profile (UserProfile)                         │
│       │       └── settings (UserSettings)                       │
│       │                                                         │
│       └── Views                                                 │
│               ├── "default" ─── all attributes                  │
│               ├── "tiny" ────── id, name only                   │
│               └── "summary" ─── id, name, email                 │
│                                                                 │
│  Generated Go Structs:                                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  type User struct {                                      │   │
│  │      ID        *int64                                    │   │
│  │      Name      *string                                   │   │
│  │      Email     *string                                   │   │
│  │      CreatedAt *string                                   │   │
│  │      ...                                                 │   │
│  │  }                                                       │   │
│  │                                                          │   │
│  │  type UserTiny struct {                                  │   │
│  │      ID   *int64                                         │   │
│  │      Name *string                                        │   │
│  │  }                                                       │   │
│  │                                                          │   │
│  │  type UserSummary struct {                               │   │
│  │      ID    *int64                                        │   │
│  │      Name  *string                                       │   │
│  │      Email *string                                       │   │
│  │  }                                                       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Named Result Types

```go
// Give the result type an explicit Go type name
var User = ResultType("application/vnd.user", "User", func() {
    Description("User resource")
    
    Attributes(func() {
        Attribute("id", Int64)
        Attribute("name", String)
        Attribute("email", String)
    })
    
    Required("id", "name", "email")
})
```

### Collection Result Types

```go
// Define a collection result type
var Users = CollectionOf(User, func() {
    // Views for the collection inherit from the element type
    View("default")
    View("tiny")
    View("summary")
})

// Or create a custom collection result type
var UserList = ResultType("application/vnd.user-list", func() {
    Description("List of users with pagination")
    
    Attributes(func() {
        Attribute("users", CollectionOf(User), "List of users")
        Attribute("total", Int64, "Total count")
        Attribute("page", Int32, "Current page")
        Attribute("per_page", Int32, "Items per page")
        Attribute("total_pages", Int32, "Total pages")
    })
    
    View("default", func() {
        Attribute("users", func() {
            View("summary") // Use summary view for nested users
        })
        Attribute("total")
        Attribute("page")
        Attribute("per_page")
        Attribute("total_pages")
    })
    
    View("full", func() {
        Attribute("users") // Uses default view for nested users
        Attribute("total")
        Attribute("page")
        Attribute("per_page")
        Attribute("total_pages")
    })
    
    Required("users", "total", "page", "per_page")
})
```

### Using Result Types in Methods

```go
var _ = Service("users", func() {
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int64, "User ID")
            Required("id")
        })
        
        // Default view
        Result(User)
        
        HTTP(func() {
            GET("/users/{id}")
            Response(StatusOK)
        })
    })
    
    Method("get_summary", func() {
        Payload(func() {
            Attribute("id", Int64)
            Required("id")
        })
        
        // Specific view
        Result(User, func() {
            View("summary")
        })
        
        HTTP(func() {
            GET("/users/{id}/summary")
            Response(StatusOK)
        })
    })
    
    Method("list", func() {
        Payload(func() {
            Attribute("page", Int32, func() {
                Default(1)
                Minimum(1)
            })
            Attribute("per_page", Int32, func() {
                Default(20)
                Minimum(1)
                Maximum(100)
            })
        })
        
        Result(UserList)
        
        HTTP(func() {
            GET("/users")
            Param("page")
            Param("per_page")
            Response(StatusOK)
        })
    })
})
```

---

## 👁️ Views

### Understanding Views

Views are **projections** of a Result Type - they define which attributes are included in different contexts.

```
┌─────────────────────────────────────────────────────────────────┐
│                      VIEWS CONCEPT                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Full Type Attributes:                                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  id | name | email | phone | address | profile | created │  │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Views (Projections):                                           │
│                                                                 │
│  "default" view ─────────────────────────────────────────────   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  id | name | email | phone | address | profile | created │  │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  "summary" view ──────────────────────────                      │
│  ┌─────────────────────────────┐                               │
│  │  id | name | email          │                               │
│  └─────────────────────────────┘                               │
│                                                                 │
│  "tiny" view ─────────                                          │
│  ┌─────────────┐                                               │
│  │  id | name  │                                               │
│  └─────────────┘                                               │
│                                                                 │
│  "admin" view ────────────────────────────────────────────────  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  id | name | email | phone | address | profile | created │  │
│  │  + internal_id | audit_log | permissions                 │  │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Defining Views

```go
var _ = ResultType("application/vnd.article", func() {
    Description("An article resource")
    
    Attributes(func() {
        // Basic info
        Attribute("id", Int64, "Article ID")
        Attribute("title", String, "Article title")
        Attribute("slug", String, "URL slug")
        Attribute("excerpt", String, "Short excerpt")
        
        // Content
        Attribute("content", String, "Full article content")
        Attribute("content_html", String, "Rendered HTML content")
        
        // Metadata
        Attribute("author", User, "Article author")
        Attribute("category", String, "Category")
        Attribute("tags", ArrayOf(String), "Tags")
        
        // Timestamps
        Attribute("created_at", String, func() {
            Format(FormatDateTime)
        })
        Attribute("updated_at", String, func() {
            Format(FormatDateTime)
        })
        Attribute("published_at", String, func() {
            Format(FormatDateTime)
        })
        
        // Stats
        Attribute("view_count", Int64, "View count")
        Attribute("like_count", Int64, "Like count")
        Attribute("comment_count", Int64, "Comment count")
    })
    
    // Default view - full article detail
    View("default", func() {
        Attribute("id")
        Attribute("title")
        Attribute("slug")
        Attribute("content_html")
        Attribute("author", func() {
            View("summary") // Use summary view for nested author
        })
        Attribute("category")
        Attribute("tags")
        Attribute("created_at")
        Attribute("updated_at")
        Attribute("published_at")
        Attribute("view_count")
        Attribute("like_count")
        Attribute("comment_count")
    })
    
    // List view - for article listings
    View("list", func() {
        Attribute("id")
        Attribute("title")
        Attribute("slug")
        Attribute("excerpt")
        Attribute("author", func() {
            View("tiny") // Just author name
        })
        Attribute("category")
        Attribute("published_at")
        Attribute("view_count")
    })
    
    // Card view - for previews
    View("card", func() {
        Attribute("id")
        Attribute("title")
        Attribute("slug")
        Attribute("excerpt")
        Attribute("category")
        Attribute("published_at")
    })
    
    // Edit view - for editors
    View("edit", func() {
        Attribute("id")
        Attribute("title")
        Attribute("slug")
        Attribute("content") // Raw content, not HTML
        Attribute("category")
        Attribute("tags")
    })
    
    Required("id", "title", "slug")
})
```

### Dynamic View Selection

```go
var _ = Service("articles", func() {
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int64)
            // Allow client to request specific view
            Attribute("view", String, func() {
                Enum("default", "list", "card", "edit")
                Default("default")
            })
            Required("id")
        })
        
        Result(Article)
        
        HTTP(func() {
            GET("/articles/{id}")
            Param("view")
            Response(StatusOK)
        })
    })
    
    Method("list", func() {
        Payload(ListArticlesPayload)
        
        // Collections use list view by default
        Result(CollectionOf(Article), func() {
            View("list")
        })
        
        HTTP(func() {
            GET("/articles")
            Response(StatusOK)
        })
    })
})
```

### Service Implementation with Views

```go
// service/articles.go
package service

import (
    "context"
    
    articles "myproject/gen/articles"
)

type articlesService struct {
    repo Repository
}

func (s *articlesService) Get(ctx context.Context, p *articles.GetPayload) (*articles.Article, error) {
    article, err := s.repo.FindByID(ctx, p.ID)
    if err != nil {
        return nil, err
    }
    
    // Convert to result type
    result := &articles.Article{
        ID:           &article.ID,
        Title:        &article.Title,
        Slug:         &article.Slug,
        Excerpt:      &article.Excerpt,
        Content:      &article.Content,
        ContentHTML:  &article.ContentHTML,
        Category:     &article.Category,
        Tags:         article.Tags,
        CreatedAt:    &article.CreatedAt,
        UpdatedAt:    &article.UpdatedAt,
        PublishedAt:  article.PublishedAt,
        ViewCount:    &article.ViewCount,
        LikeCount:    &article.LikeCount,
        CommentCount: &article.CommentCount,
    }
    
    // Author with appropriate view
    if article.Author != nil {
        result.Author = &articles.User{
            ID:   &article.Author.ID,
            Name: &article.Author.Name,
            // Include email only if view needs it
        }
    }
    
    return result, nil
}
```

### View Selection in Generated Code

```go
// gen/articles/service.go (generated)
type Article struct {
    ID           *int64
    Title        *string
    Slug         *string
    // ... all fields
}

// gen/articles/views/views.go (generated)
type ArticleView struct {
    Projected *ArticleViewProjected
    View      string
}

type ArticleViewProjected struct {
    ID           *int64
    Title        *string
    // ... based on view
}

// View rendering happens during encoding
func (v *ArticleView) Render() *ArticleViewProjected {
    // Returns only fields for the selected view
}
```

---

## 🔄 Type Reuse

### The DRY Principle in Goa

Don't Repeat Yourself - reuse type definitions across your API.

### Using Reference

`Reference` copies attribute definitions from one type to another:

```go
// design/types.go
package design

import . "goa.design/goa/v3/dsl"

// Base type with common fields
var BaseEntity = Type("BaseEntity", func() {
    Attribute("id", Int64, "Unique identifier")
    Attribute("created_at", String, "Creation timestamp", func() {
        Format(FormatDateTime)
    })
    Attribute("updated_at", String, "Last update timestamp", func() {
        Format(FormatDateTime)
    })
    Attribute("created_by", String, "Creator user ID")
    Attribute("updated_by", String, "Last editor user ID")
})

// User type references BaseEntity
var User = Type("User", func() {
    // Copy all attributes from BaseEntity
    Reference(BaseEntity)
    
    // User-specific attributes
    Attribute("name", String, "User name")
    Attribute("email", String, "Email address")
    Attribute("role", String, func() {
        Enum("user", "admin", "moderator")
    })
    
    Required("id", "name", "email")
})

// Product type also references BaseEntity
var Product = Type("Product", func() {
    Reference(BaseEntity)
    
    Attribute("name", String, "Product name")
    Attribute("description", String, "Product description")
    Attribute("price", Float64, "Price")
    Attribute("sku", String, "Stock keeping unit")
    
    Required("id", "name", "price")
})

// Order references BaseEntity
var Order = Type("Order", func() {
    Reference(BaseEntity)
    
    Attribute("user_id", Int64, "Customer ID")
    Attribute("items", ArrayOf(OrderItem), "Order items")
    Attribute("total", Float64, "Order total")
    Attribute("status", String, func() {
        Enum("pending", "processing", "shipped", "delivered", "cancelled")
    })
    
    Required("id", "user_id", "items", "total", "status")
})
```

### Selective Reference

```go
// Only reference specific attributes
var UserProfile = Type("UserProfile", func() {
    // Reference only id from BaseEntity
    Attribute("id", Int64, func() {
        Reference(BaseEntity)
    })
    
    Attribute("bio", String)
    Attribute("avatar_url", String)
    Attribute("website", String)
})
```

### Extend for Type Extension

`Extend` creates a new type based on an existing one:

```go
// Base address type
var Address = Type("Address", func() {
    Attribute("street", String)
    Attribute("city", String)
    Attribute("state", String)
    Attribute("postal_code", String)
    Attribute("country", String, func() {
        Default("US")
    })
    
    Required("street", "city", "country")
})

// Extended address with additional fields
var DetailedAddress = Type("DetailedAddress", func() {
    // Extend copies all attributes and adds new ones
    Extend(Address)
    
    Attribute("unit", String, "Apartment/Suite number")
    Attribute("building_name", String)
    Attribute("latitude", Float64)
    Attribute("longitude", Float64)
    Attribute("delivery_instructions", String)
})

// Billing address with payment info
var BillingAddress = Type("BillingAddress", func() {
    Extend(Address)
    
    Attribute("is_default", Boolean, func() {
        Default(false)
    })
    Attribute("label", String, "Address label (Home, Work, etc.)")
})
```

### Reusable Attribute Patterns

```go
// design/attributes.go
package design

import . "goa.design/goa/v3/dsl"

// Reusable attribute definitions using Meta
var EmailAttribute = func() {
    Attribute("email", String, "Email address", func() {
        Format(FormatEmail)
        Example("user@example.com")
    })
}

var PhoneAttribute = func() {
    Attribute("phone", String, "Phone number", func() {
        Pattern(`^\+?[1-9]\d{1,14}$`)
        Example("+1234567890")
    })
}

var PaginationAttributes = func() {
    Attribute("page", Int32, "Page number", func() {
        Minimum(1)
        Default(1)
    })
    Attribute("per_page", Int32, "Items per page", func() {
        Minimum(1)
        Maximum(100)
        Default(20)
    })
}

var TimestampAttributes = func() {
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
        Description("Creation timestamp")
        Example("2024-01-15T10:30:00Z")
    })
    Attribute("updated_at", String, func() {
        Format(FormatDateTime)
        Description("Last update timestamp")
        Example("2024-01-15T10:30:00Z")
    })
}

// Usage in types
var Contact = Type("Contact", func() {
    Attribute("id", Int64)
    Attribute("name", String)
    EmailAttribute()        // Reuse email
    PhoneAttribute()        // Reuse phone
    TimestampAttributes()   // Reuse timestamps
    
    Required("id", "name", "email")
})

var ListPayload = Type("ListPayload", func() {
    PaginationAttributes()  // Reuse pagination
    Attribute("sort", String, func() {
        Enum("created_at", "updated_at", "name")
        Default("created_at")
    })
    Attribute("order", String, func() {
        Enum("asc", "desc")
        Default("desc")
    })
})
```

### Type Aliases

```go
// Create type aliases for clarity
var UserID = Type("UserID", Int64, func() {
    Description("User identifier")
    Minimum(1)
})

var ProductID = Type("ProductID", Int64, func() {
    Description("Product identifier")
    Minimum(1)
})

var Money = Type("Money", func() {
    Description("Monetary value")
    Attribute("amount", Float64, "Amount in dollars")
    Attribute("currency", String, func() {
        Enum("USD", "EUR", "GBP", "JPY")
        Default("USD")
    })
    Required("amount")
})

// Usage
var Order = Type("Order", func() {
    Attribute("id", Int64)
    Attribute("user_id", UserID)
    Attribute("total", Money)
    Attribute("items", ArrayOf(OrderItem))
})
```

---

## 🔀 Polymorphism

### Understanding Polymorphism in APIs

Polymorphism allows different types to be used interchangeably through a common interface.

```
┌─────────────────────────────────────────────────────────────────┐
│                    API POLYMORPHISM                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Scenario: Notification System                                  │
│                                                                 │
│  All notifications share:                                       │
│    - id                                                         │
│    - type                                                       │
│    - created_at                                                 │
│    - read                                                       │
│                                                                 │
│  But differ in details:                                         │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ EmailNotif      │  │ SMSNotif        │  │ PushNotif       │ │
│  ├─────────────────┤  ├─────────────────┤  ├─────────────────┤ │
│  │ + subject       │  │ + phone_number  │  │ + device_token  │ │
│  │ + body          │  │ + message       │  │ + title         │ │
│  │ + recipient     │  │ + sender_id     │  │ + body          │ │
│  │ + from_address  │  │                 │  │ + badge_count   │ │
│  │ + attachments   │  │                 │  │ + data          │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                                 │
│  API returns: { notifications: [mixed types] }                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### OneOf Pattern (Discriminated Union)

```go
// design/notifications.go
package design

import . "goa.design/goa/v3/dsl"

// Base notification attributes
var NotificationBase = Type("NotificationBase", func() {
    Attribute("id", Int64, "Notification ID")
    Attribute("type", String, "Notification type", func() {
        Enum("email", "sms", "push", "webhook")
    })
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("read", Boolean, func() {
        Default(false)
    })
    Attribute("user_id", Int64, "Target user ID")
    
    Required("id", "type", "user_id")
})

// Email-specific notification
var EmailNotification = Type("EmailNotification", func() {
    Extend(NotificationBase)
    
    Attribute("subject", String, "Email subject")
    Attribute("body", String, "Email body (HTML)")
    Attribute("recipient", String, "Email recipient", func() {
        Format(FormatEmail)
    })
    Attribute("from_address", String, func() {
        Format(FormatEmail)
    })
    Attribute("attachments", ArrayOf(String), "Attachment URLs")
    
    Required("subject", "body", "recipient")
})

// SMS notification
var SMSNotification = Type("SMSNotification", func() {
    Extend(NotificationBase)
    
    Attribute("phone_number", String, "Recipient phone")
    Attribute("message", String, "SMS message", func() {
        MaxLength(160)
    })
    Attribute("sender_id", String, "Sender ID")
    
    Required("phone_number", "message")
})

// Push notification
var PushNotification = Type("PushNotification", func() {
    Extend(NotificationBase)
    
    Attribute("device_token", String, "Device push token")
    Attribute("title", String, "Notification title")
    Attribute("body", String, "Notification body")
    Attribute("badge_count", Int32, "Badge number")
    Attribute("data", MapOf(String, Any), "Custom data payload")
    
    Required("device_token", "title", "body")
})

// OneOf for polymorphic notification
var Notification = Type("Notification", func() {
    Description("A notification (email, SMS, or push)")
    
    OneOf("notification", func() {
        Attribute("email", EmailNotification, "Email notification")
        Attribute("sms", SMSNotification, "SMS notification")
        Attribute("push", PushNotification, "Push notification")
    })
})
```

### Service Using Polymorphic Types

```go
var _ = Service("notifications", func() {
    Description("Notification management service")
    
    Method("list", func() {
        Payload(func() {
            Attribute("user_id", Int64)
            Attribute("type", String, func() {
                Enum("email", "sms", "push", "all")
                Default("all")
            })
            PaginationAttributes()
            Required("user_id")
        })
        
        Result(func() {
            Attribute("notifications", ArrayOf(Notification))
            Attribute("total", Int64)
            Required("notifications", "total")
        })
        
        HTTP(func() {
            GET("/users/{user_id}/notifications")
            Param("type")
            Param("page")
            Param("per_page")
            Response(StatusOK)
        })
    })
    
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int64)
            Required("id")
        })
        
        Result(Notification)
        
        HTTP(func() {
            GET("/notifications/{id}")
            Response(StatusOK)
        })
    })
    
    Method("send", func() {
        Payload(Notification) // Can be any notification type
        Result(Notification)
        
        HTTP(func() {
            POST("/notifications")
            Response(StatusCreated)
        })
    })
})
```

### Type Switching in Implementation

```go
// service/notifications.go
package service

import (
    "context"
    "fmt"
    
    notifications "myproject/gen/notifications"
)

type notificationsService struct {
    emailSender EmailSender
    smsSender   SMSSender
    pushSender  PushSender
}

func (s *notificationsService) Send(ctx context.Context, p *notifications.Notification) (*notifications.Notification, error) {
    // Type switch on the OneOf field
    switch {
    case p.Email != nil:
        return s.sendEmail(ctx, p.Email)
    case p.SMS != nil:
        return s.sendSMS(ctx, p.SMS)
    case p.Push != nil:
        return s.sendPush(ctx, p.Push)
    default:
        return nil, fmt.Errorf("unknown notification type")
    }
}

func (s *notificationsService) sendEmail(ctx context.Context, n *notifications.EmailNotification) (*notifications.Notification, error) {
    // Send email
    err := s.emailSender.Send(ctx, n.Recipient, n.Subject, n.Body)
    if err != nil {
        return nil, err
    }
    
    // Save to database and return
    return &notifications.Notification{
        Email: n,
    }, nil
}
```

### Discriminator Field Pattern

```go
// Alternative: using a discriminator field
var Event = Type("Event", func() {
    Description("An event with polymorphic data")
    
    // Discriminator field
    Attribute("type", String, "Event type", func() {
        Enum("user.created", "user.updated", "order.placed", "order.shipped")
    })
    
    // Common fields
    Attribute("id", String, "Event ID")
    Attribute("timestamp", String, func() {
        Format(FormatDateTime)
    })
    Attribute("source", String, "Event source")
    
    // Polymorphic data field
    Attribute("data", Any, "Event-specific data")
    
    Required("type", "id", "timestamp", "data")
})

// In service implementation
func (s *svc) HandleEvent(ctx context.Context, event *events.Event) error {
    switch *event.Type {
    case "user.created":
        var data UserCreatedData
        if err := mapstructure.Decode(event.Data, &data); err != nil {
            return err
        }
        return s.handleUserCreated(ctx, &data)
        
    case "order.placed":
        var data OrderPlacedData
        if err := mapstructure.Decode(event.Data, &data); err != nil {
            return err
        }
        return s.handleOrderPlaced(ctx, &data)
        
    default:
        return fmt.Errorf("unknown event type: %s", *event.Type)
    }
}
```

---

## 🧬 Inheritance Patterns

### Simulating Inheritance in Goa

Go doesn't have traditional inheritance, but Goa provides patterns to achieve similar results.

### Base Type Pattern

```go
// design/inheritance.go
package design

import . "goa.design/goa/v3/dsl"

// ============================================
// PATTERN 1: Base Type with Extend
// ============================================

// Abstract base - defines common structure
var Vehicle = Type("Vehicle", func() {
    Description("Base vehicle type")
    
    Attribute("id", Int64, "Vehicle ID")
    Attribute("type", String, "Vehicle type", func() {
        Enum("car", "truck", "motorcycle", "bicycle")
    })
    Attribute("brand", String, "Brand name")
    Attribute("model", String, "Model name")
    Attribute("year", Int32, "Manufacturing year", func() {
        Minimum(1900)
        Maximum(2100)
    })
    Attribute("color", String, "Color")
    Attribute("vin", String, "Vehicle identification number")
    
    Required("id", "type", "brand", "model", "year")
})

// Concrete types extend base
var Car = Type("Car", func() {
    Extend(Vehicle)
    
    Attribute("num_doors", Int32, "Number of doors", func() {
        Enum(2, 4, 5)
    })
    Attribute("fuel_type", String, func() {
        Enum("gasoline", "diesel", "electric", "hybrid")
    })
    Attribute("transmission", String, func() {
        Enum("manual", "automatic", "cvt")
    })
    Attribute("trunk_capacity", Float64, "Trunk capacity in liters")
})

var Truck = Type("Truck", func() {
    Extend(Vehicle)
    
    Attribute("cargo_capacity", Float64, "Cargo capacity in kg")
    Attribute("bed_length", Float64, "Bed length in feet")
    Attribute("towing_capacity", Float64, "Towing capacity in kg")
    Attribute("num_axles", Int32, "Number of axles")
})

var Motorcycle = Type("Motorcycle", func() {
    Extend(Vehicle)
    
    Attribute("engine_cc", Int32, "Engine displacement in CC")
    Attribute("bike_type", String, func() {
        Enum("sport", "cruiser", "touring", "adventure", "standard")
    })
    Attribute("has_sidecar", Boolean, func() {
        Default(false)
    })
})
```

### Interface Pattern

```go
// ============================================
// PATTERN 2: Interface Pattern
// ============================================

// Interface definition - what all payment methods must have
var PaymentMethodInterface = Type("PaymentMethodInterface", func() {
    Description("Common payment method fields")
    
    Attribute("id", Int64, "Payment method ID")
    Attribute("type", String, "Payment method type", func() {
        Enum("credit_card", "bank_account", "paypal", "crypto")
    })
    Attribute("is_default", Boolean, "Is default payment method")
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
})

// Implementations
var CreditCard = Type("CreditCard", func() {
    Reference(PaymentMethodInterface) // Implements interface
    
    // Interface fields (required)
    Attribute("id", Int64)
    Attribute("type", String, func() {
        Default("credit_card")
    })
    Attribute("is_default", Boolean)
    Attribute("created_at", String)
    
    // Implementation-specific fields
    Attribute("card_number", String, "Masked card number")
    Attribute("card_brand", String, func() {
        Enum("visa", "mastercard", "amex", "discover")
    })
    Attribute("exp_month", Int32, func() {
        Minimum(1)
        Maximum(12)
    })
    Attribute("exp_year", Int32)
    Attribute("last_four", String, func() {
        Pattern(`^\d{4}$`)
    })
    Attribute("billing_address", Address)
})

var BankAccount = Type("BankAccount", func() {
    Reference(PaymentMethodInterface)
    
    Attribute("id", Int64)
    Attribute("type", String, func() {
        Default("bank_account")
    })
    Attribute("is_default", Boolean)
    Attribute("created_at", String)
    
    // Bank-specific fields
    Attribute("bank_name", String)
    Attribute("account_type", String, func() {
        Enum("checking", "savings")
    })
    Attribute("routing_number", String)
    Attribute("account_number_last_four", String)
})

var PayPalAccount = Type("PayPalAccount", func() {
    Reference(PaymentMethodInterface)
    
    Attribute("id", Int64)
    Attribute("type", String, func() {
        Default("paypal")
    })
    Attribute("is_default", Boolean)
    Attribute("created_at", String)
    
    // PayPal-specific
    Attribute("email", String, func() {
        Format(FormatEmail)
    })
    Attribute("payer_id", String)
})
```

### Hierarchy Pattern

```go
// ============================================
// PATTERN 3: Type Hierarchy
// ============================================

// Level 1: Root
var Content = Type("Content", func() {
    Description("Base content type")
    
    Attribute("id", Int64)
    Attribute("title", String)
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("author_id", Int64)
    Attribute("status", String, func() {
        Enum("draft", "published", "archived")
    })
})

// Level 2: Branches
var TextContent = Type("TextContent", func() {
    Extend(Content)
    
    Attribute("body", String, "Text content body")
    Attribute("word_count", Int32)
    Attribute("reading_time_minutes", Int32)
})

var MediaContent = Type("MediaContent", func() {
    Extend(Content)
    
    Attribute("url", String, "Media URL")
    Attribute("file_size", Int64, "File size in bytes")
    Attribute("mime_type", String, "MIME type")
    Attribute("duration_seconds", Int32, "Duration for video/audio")
})

// Level 3: Leaves
var Article = Type("Article", func() {
    Extend(TextContent)
    
    Attribute("category", String)
    Attribute("tags", ArrayOf(String))
    Attribute("featured_image", String)
    Attribute("excerpt", String)
})

var BlogPost = Type("BlogPost", func() {
    Extend(TextContent)
    
    Attribute("slug", String)
    Attribute("comments_enabled", Boolean)
    Attribute("pinned", Boolean)
})

var Video = Type("Video", func() {
    Extend(MediaContent)
    
    Attribute("thumbnail_url", String)
    Attribute("resolution", String, func() {
        Enum("480p", "720p", "1080p", "4k")
    })
    Attribute("captions_url", String)
})

var Podcast = Type("Podcast", func() {
    Extend(MediaContent)
    
    Attribute("episode_number", Int32)
    Attribute("season_number", Int32)
    Attribute("transcript_url", String)
})
```

### Inheritance Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    TYPE HIERARCHY                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                        Content                                  │
│                           │                                     │
│              ┌────────────┴────────────┐                        │
│              │                         │                        │
│        TextContent               MediaContent                   │
│              │                         │                        │
│       ┌──────┴──────┐           ┌──────┴──────┐                 │
│       │             │           │             │                 │
│    Article     BlogPost      Video       Podcast                │
│                                                                 │
│  Each level adds specific attributes:                           │
│                                                                 │
│  Content:      id, title, created_at, author_id, status        │
│                      │                                          │
│  TextContent:        + body, word_count, reading_time          │
│                      │                                          │
│  Article:            + category, tags, featured_image          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🧱 Struct Composition

### Composition Over Inheritance

Go favors composition over inheritance. Goa supports this through embedded types and attribute grouping.

```go
// design/composition.go
package design

import . "goa.design/goa/v3/dsl"

// ============================================
// Composable Components
// ============================================

// Audit information component
var Auditable = Type("Auditable", func() {
    Description("Audit trail fields")
    
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("created_by", String)
    Attribute("updated_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("updated_by", String)
    Attribute("version", Int32, "Optimistic locking version")
})

// Soft delete component
var SoftDeletable = Type("SoftDeletable", func() {
    Description("Soft delete fields")
    
    Attribute("deleted", Boolean, func() {
        Default(false)
    })
    Attribute("deleted_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("deleted_by", String)
})

// Taggable component
var Taggable = Type("Taggable", func() {
    Description("Tagging fields")
    
    Attribute("tags", ArrayOf(String), "Tags")
    Attribute("labels", MapOf(String, String), "Key-value labels")
})

// Ownable component
var Ownable = Type("Ownable", func() {
    Description("Ownership fields")
    
    Attribute("owner_id", Int64, "Owner user ID")
    Attribute("owner_type", String, "Owner type", func() {
        Enum("user", "organization", "team")
    })
    Attribute("shared_with", ArrayOf(Int64), "User IDs with access")
})

// ============================================
// Composed Types
// ============================================

// Document composes multiple components
var Document = Type("Document", func() {
    Description("A document with full audit, tagging, and ownership")
    
    // Identity
    Attribute("id", Int64)
    Attribute("name", String)
    Attribute("content", String)
    Attribute("mime_type", String)
    Attribute("size", Int64)
    
    // Compose: Audit
    Attribute("audit", Auditable, "Audit information")
    
    // Compose: Soft Delete
    Attribute("deletion", SoftDeletable, "Deletion information")
    
    // Compose: Tags
    Attribute("tagging", Taggable, "Tagging information")
    
    // Compose: Ownership
    Attribute("ownership", Ownable, "Ownership information")
    
    Required("id", "name")
})
```

### Flat Composition Pattern

```go
// Flatten composed attributes for simpler API
var DocumentFlat = Type("DocumentFlat", func() {
    Description("Document with flattened structure")
    
    // Core fields
    Attribute("id", Int64)
    Attribute("name", String)
    Attribute("content", String)
    
    // Inline audit fields
    Reference(Auditable)
    Attribute("created_at", String)
    Attribute("created_by", String)
    Attribute("updated_at", String)
    Attribute("updated_by", String)
    Attribute("version", Int32)
    
    // Inline soft delete fields
    Reference(SoftDeletable)
    Attribute("deleted", Boolean)
    Attribute("deleted_at", String)
    Attribute("deleted_by", String)
    
    // Inline taggable fields
    Reference(Taggable)
    Attribute("tags", ArrayOf(String))
    Attribute("labels", MapOf(String, String))
    
    // Inline ownable fields
    Reference(Ownable)
    Attribute("owner_id", Int64)
    Attribute("owner_type", String)
    Attribute("shared_with", ArrayOf(Int64))
    
    Required("id", "name")
})
```

### Mixin Pattern

```go
// design/mixins.go
package design

import . "goa.design/goa/v3/dsl"

// Mixin functions add common attribute sets
func WithID() {
    Attribute("id", Int64, "Unique identifier", func() {
        Minimum(1)
    })
}

func WithTimestamps() {
    Attribute("created_at", String, "Creation timestamp", func() {
        Format(FormatDateTime)
    })
    Attribute("updated_at", String, "Last update timestamp", func() {
        Format(FormatDateTime)
    })
}

func WithSoftDelete() {
    Attribute("deleted", Boolean, "Is deleted", func() {
        Default(false)
    })
    Attribute("deleted_at", String, "Deletion timestamp", func() {
        Format(FormatDateTime)
    })
}

func WithPagination() {
    Attribute("page", Int32, "Page number", func() {
        Minimum(1)
        Default(1)
    })
    Attribute("per_page", Int32, "Items per page", func() {
        Minimum(1)
        Maximum(100)
        Default(20)
    })
}

func WithSort(allowedFields ...string) {
    Attribute("sort_by", String, "Sort field", func() {
        Enum(allowedFields...)
        Default(allowedFields[0])
    })
    Attribute("sort_order", String, "Sort order", func() {
        Enum("asc", "desc")
        Default("asc")
    })
}

func WithSearch() {
    Attribute("query", String, "Search query")
    Attribute("search_in", ArrayOf(String), "Fields to search in")
}

// Usage
var Project = Type("Project", func() {
    WithID()
    WithTimestamps()
    WithSoftDelete()
    
    Attribute("name", String)
    Attribute("description", String)
    Attribute("status", String, func() {
        Enum("active", "inactive", "archived")
    })
    
    Required("id", "name")
})

var ListProjectsPayload = Type("ListProjectsPayload", func() {
    WithPagination()
    WithSort("name", "created_at", "updated_at")
    WithSearch()
    
    Attribute("status", String, func() {
        Enum("active", "inactive", "archived", "all")
        Default("all")
    })
})
```

### Component Assembly Pattern

```go
// design/components.go
package design

import . "goa.design/goa/v3/dsl"

// Component interfaces (as type definitions)
var Identifiable = Type("Identifiable", func() {
    Attribute("id", Int64)
})

var Named = Type("Named", func() {
    Attribute("name", String)
    Attribute("slug", String)
})

var Describable = Type("Describable", func() {
    Attribute("description", String)
    Attribute("summary", String, func() {
        MaxLength(255)
    })
})

var Publishable = Type("Publishable", func() {
    Attribute("status", String, func() {
        Enum("draft", "review", "published", "archived")
    })
    Attribute("published_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("published_by", Int64)
})

var Commentable = Type("Commentable", func() {
    Attribute("comments_enabled", Boolean, func() {
        Default(true)
    })
    Attribute("comment_count", Int64)
})

var Likeable = Type("Likeable", func() {
    Attribute("likes_count", Int64)
    Attribute("liked_by_current_user", Boolean)
})

// Assemble types from components
var Post = Type("Post", func() {
    Description("A blog post - assembled from components")
    
    // Include component attributes
    Reference(Identifiable)
    Attribute("id", Int64)
    
    Reference(Named)
    Attribute("name", String)
    Attribute("slug", String)
    
    Reference(Describable)
    Attribute("description", String)
    Attribute("summary", String)
    
    Reference(Publishable)
    Attribute("status", String)
    Attribute("published_at", String)
    Attribute("published_by", Int64)
    
    Reference(Commentable)
    Attribute("comments_enabled", Boolean)
    Attribute("comment_count", Int64)
    
    Reference(Likeable)
    Attribute("likes_count", Int64)
    Attribute("liked_by_current_user", Boolean)
    
    // Post-specific
    Attribute("content", String)
    Attribute("author_id", Int64)
    Attribute("category_id", Int64)
    Attribute("featured_image_url", String)
    
    Required("id", "name", "slug", "content", "author_id")
})
```

---

## 🎨 Advanced Patterns

### Generic Collection Response

```go
// design/generics.go
package design

import . "goa.design/goa/v3/dsl"

// Generic paginated response builder
func PaginatedResult(itemType interface{}, name string) {
    ResultType("application/vnd."+name+"-list", func() {
        Description("Paginated list of " + name)
        
        Attributes(func() {
            Attribute("items", ArrayOf(itemType), "List of items")
            Attribute("pagination", PaginationMeta, "Pagination metadata")
        })
        
        View("default", func() {
            Attribute("items")
            Attribute("pagination")
        })
        
        Required("items", "pagination")
    })
}

var PaginationMeta = Type("PaginationMeta", func() {
    Description("Pagination metadata")
    
    Attribute("page", Int32, "Current page")
    Attribute("per_page", Int32, "Items per page")
    Attribute("total_items", Int64, "Total item count")
    Attribute("total_pages", Int32, "Total page count")
    Attribute("has_next", Boolean, "Has next page")
    Attribute("has_prev", Boolean, "Has previous page")
    
    Required("page", "per_page", "total_items", "total_pages")
})

// Pre-defined paginated results
var UserList = ResultType("application/vnd.user-list", func() {
    Description("Paginated list of users")
    
    Attributes(func() {
        Attribute("items", CollectionOf(User), "Users")
        Attribute("pagination", PaginationMeta)
    })
    
    View("default", func() {
        Attribute("items", func() {
            View("summary")
        })
        Attribute("pagination")
    })
    
    View("full", func() {
        Attribute("items")
        Attribute("pagination")
    })
    
    Required("items", "pagination")
})
```

### Recursive Types

```go
// design/recursive.go
package design

import . "goa.design/goa/v3/dsl"

// Self-referential type for tree structures
var Category = Type("Category", func() {
    Description("A hierarchical category")
    
    Attribute("id", Int64, "Category ID")
    Attribute("name", String, "Category name")
    Attribute("slug", String, "URL slug")
    Attribute("description", String, "Category description")
    Attribute("parent_id", Int64, "Parent category ID (null for root)")
    
    // Self-reference for children
    Attribute("children", ArrayOf("Category"), "Child categories")
    
    // Computed fields
    Attribute("depth", Int32, "Depth in tree (0 = root)")
    Attribute("path", String, "Full path (e.g., 'electronics/computers/laptops')")
    Attribute("child_count", Int32, "Number of direct children")
    
    Required("id", "name", "slug")
})

// Comment with nested replies
var Comment = Type("Comment", func() {
    Description("A comment with nested replies")
    
    Attribute("id", Int64)
    Attribute("content", String)
    Attribute("author_id", Int64)
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    
    Attribute("parent_id", Int64, "Parent comment ID for replies")
    Attribute("replies", ArrayOf("Comment"), "Nested replies")
    Attribute("reply_count", Int32)
    
    Required("id", "content", "author_id")
})

// Menu with nested items
var MenuItem = Type("MenuItem", func() {
    Description("A menu item with sub-items")
    
    Attribute("id", Int64)
    Attribute("label", String, "Display label")
    Attribute("url", String, "Link URL")
    Attribute("icon", String, "Icon class/name")
    Attribute("target", String, func() {
        Enum("_self", "_blank")
        Default("_self")
    })
    
    Attribute("children", ArrayOf("MenuItem"), "Sub-menu items")
    Attribute("parent_id", Int64)
    Attribute("order", Int32, "Display order")
    
    Required("id", "label")
})
```

### Versioned Types

```go
// design/versioned.go
package design

import . "goa.design/goa/v3/dsl"

// V1 API types
var UserV1 = Type("UserV1", func() {
    Description("User - API v1")
    
    Attribute("id", Int64)
    Attribute("username", String)
    Attribute("email", String)
    
    Required("id", "username", "email")
})

// V2 API types - evolved schema
var UserV2 = Type("UserV2", func() {
    Description("User - API v2")
    
    Attribute("id", String, "Changed to UUID", func() {
        Format(FormatUUID)
    })
    Attribute("username", String)
    Attribute("email", String, func() {
        Format(FormatEmail)
    })
    
    // New fields in V2
    Attribute("profile", UserProfileV2)
    Attribute("preferences", UserPreferencesV2)
    
    // Deprecated fields (still present for compatibility)
    Attribute("name", String, "Deprecated: Use profile.display_name instead")
    
    Required("id", "username", "email")
})

// Services for different API versions
var _ = Service("users_v1", func() {
    HTTP(func() {
        Path("/v1/users")
    })
    
    Method("list", func() {
        Result(ArrayOf(UserV1))
        HTTP(func() {
            GET("/")
        })
    })
})

var _ = Service("users_v2", func() {
    HTTP(func() {
        Path("/v2/users")
    })
    
    Method("list", func() {
        Result(ArrayOf(UserV2))
        HTTP(func() {
            GET("/")
        })
    })
})
```

### Conditional Fields

```go
// design/conditional.go
package design

import . "goa.design/goa/v3/dsl"

// Response with conditional fields based on user role
var OrderResponse = ResultType("application/vnd.order", func() {
    Description("Order with role-based field visibility")
    
    Attributes(func() {
        // Public fields
        Attribute("id", Int64)
        Attribute("status", String)
        Attribute("total", Float64)
        Attribute("created_at", String)
        
        // Customer fields
        Attribute("items", ArrayOf(OrderItem))
        Attribute("shipping_address", Address)
        Attribute("tracking_number", String)
        
        // Admin-only fields
        Attribute("customer_id", Int64)
        Attribute("profit_margin", Float64)
        Attribute("internal_notes", String)
        Attribute("cost_breakdown", CostBreakdown)
        
        // Support fields
        Attribute("support_history", ArrayOf(SupportTicket))
    })
    
    // Public view - minimal info
    View("public", func() {
        Attribute("id")
        Attribute("status")
        Attribute("total")
        Attribute("created_at")
    })
    
    // Customer view - their order info
    View("customer", func() {
        Attribute("id")
        Attribute("status")
        Attribute("total")
        Attribute("items")
        Attribute("shipping_address")
        Attribute("tracking_number")
        Attribute("created_at")
    })
    
    // Support view - for support staff
    View("support", func() {
        Attribute("id")
        Attribute("status")
        Attribute("total")
        Attribute("items")
        Attribute("customer_id")
        Attribute("support_history")
        Attribute("created_at")
    })
    
    // Admin view - everything
    View("admin", func() {
        Attribute("id")
        Attribute("status")
        Attribute("total")
        Attribute("items")
        Attribute("shipping_address")
        Attribute("tracking_number")
        Attribute("customer_id")
        Attribute("profit_margin")
        Attribute("internal_notes")
        Attribute("cost_breakdown")
        Attribute("created_at")
    })
    
    // Default = customer view
    View("default", func() {
        Attribute("id")
        Attribute("status")
        Attribute("total")
        Attribute("items")
        Attribute("created_at")
    })
})
```

---

## 📦 Complete Examples

### E-Commerce Type System

```go
// design/ecommerce.go
package design

import . "goa.design/goa/v3/dsl"

// ============================================
// Base Components
// ============================================

func WithEntityID() {
    Attribute("id", Int64, "Unique identifier", func() {
        Minimum(1)
        Example(12345)
    })
}

func WithAuditFields() {
    Attribute("created_at", String, func() {
        Format(FormatDateTime)
    })
    Attribute("updated_at", String, func() {
        Format(FormatDateTime)
    })
}

// ============================================
// Core Types
// ============================================

var Money = Type("Money", func() {
    Description("Monetary value with currency")
    
    Attribute("amount", Float64, "Amount in smallest unit", func() {
        Minimum(0)
        Example(1999)
    })
    Attribute("currency", String, "ISO 4217 currency code", func() {
        Enum("USD", "EUR", "GBP", "JPY", "CAD", "AUD")
        Default("USD")
        Example("USD")
    })
    Attribute("formatted", String, "Formatted display string", func() {
        Example("$19.99")
    })
    
    Required("amount", "currency")
})

var Address = Type("Address", func() {
    Description("Physical address")
    
    Attribute("line1", String, "Street address", func() {
        MinLength(1)
        MaxLength(100)
    })
    Attribute("line2", String, "Apt, suite, etc.", func() {
        MaxLength(100)
    })
    Attribute("city", String, func() {
        MinLength(1)
        MaxLength(50)
    })
    Attribute("state", String, func() {
        MaxLength(50)
    })
    Attribute("postal_code", String, func() {
        Pattern(`^[A-Z0-9\s-]{3,10}$`)
    })
    Attribute("country", String, "ISO 3166-1 alpha-2", func() {
        Pattern(`^[A-Z]{2}$`)
        Default("US")
    })
    
    Required("line1", "city", "country")
})

// ============================================
// Product Types
// ============================================

var Product = ResultType("application/vnd.product", func() {
    Description("A product")
    
    Attributes(func() {
        WithEntityID()
        
        Attribute("name", String, "Product name", func() {
            MinLength(1)
            MaxLength(200)
        })
        Attribute("slug", String, "URL-friendly identifier")
        Attribute("description", String, "Full description")
        Attribute("short_description", String, func() {
            MaxLength(500)
        })
        
        Attribute("sku", String, "Stock keeping unit")
        Attribute("barcode", String, "UPC/EAN barcode")
        
        Attribute("price", Money, "Current price")
        Attribute("compare_at_price", Money, "Original price for sales")
        Attribute("cost", Money, "Cost price (admin only)")
        
        Attribute("status", String, func() {
            Enum("draft", "active", "archived")
            Default("draft")
        })
        
        Attribute("inventory_quantity", Int32)
        Attribute("inventory_policy", String, func() {
            Enum("deny", "continue")
            Default("deny")
        })
        
        Attribute("categories", ArrayOf(Category))
        Attribute("tags", ArrayOf(String))
        Attribute("images", ArrayOf(ProductImage))
        Attribute("variants", ArrayOf(ProductVariant))
        
        Attribute("weight", Float64, "Weight in grams")
        Attribute("dimensions", Dimensions)
        
        WithAuditFields()
    })
    
    View("default", func() {
        Attribute("id")
        Attribute("name")
        Attribute("slug")
        Attribute("description")
        Attribute("price")
        Attribute("compare_at_price")
        Attribute("images")
        Attribute("variants")
        Attribute("status")
    })
    
    View("card", func() {
        Attribute("id")
        Attribute("name")
        Attribute("slug")
        Attribute("short_description")
        Attribute("price")
        Attribute("compare_at_price")
        Attribute("images", func() {
            // Only first image
        })
    })
    
    View("admin", func() {
        Attribute("id")
        Attribute("name")
        Attribute("slug")
        Attribute("description")
        Attribute("sku")
        Attribute("barcode")
        Attribute("price")
        Attribute("compare_at_price")
        Attribute("cost")
        Attribute("status")
        Attribute("inventory_quantity")
        Attribute("inventory_policy")
        Attribute("categories")
        Attribute("tags")
        Attribute("images")
        Attribute("variants")
        Attribute("weight")
        Attribute("dimensions")
        Attribute("created_at")
        Attribute("updated_at")
    })
    
    Required("id", "name", "slug", "price", "status")
})

var ProductImage = Type("ProductImage", func() {
    Attribute("id", Int64)
    Attribute("url", String, func() {
        Format(FormatURI)
    })
    Attribute("alt_text", String)
    Attribute("position", Int32)
    Attribute("variant_ids", ArrayOf(Int64), "Associated variant IDs")
    
    Required("id", "url")
})

var ProductVariant = Type("ProductVariant", func() {
    WithEntityID()
    
    Attribute("product_id", Int64)
    Attribute("name", String, "Variant name (e.g., 'Large / Red')")
    Attribute("sku", String)
    Attribute("price", Money)
    Attribute("inventory_quantity", Int32)
    
    Attribute("option1", String, "First option value")
    Attribute("option2", String, "Second option value")
    Attribute("option3", String, "Third option value")
    
    Attribute("image_id", Int64)
    Attribute("weight", Float64)
    
    Required("id", "product_id", "name", "price")
})

var Dimensions = Type("Dimensions", func() {
    Attribute("length", Float64, "Length in cm")
    Attribute("width", Float64, "Width in cm")
    Attribute("height", Float64, "Height in cm")
})

// ============================================
// Order Types
// ============================================

var Order = ResultType("application/vnd.order", func() {
    Description("A customer order")
    
    Attributes(func() {
        WithEntityID()
        
        Attribute("order_number", String, "Human-readable order number")
        Attribute("customer_id", Int64)
        Attribute("customer", Customer)
        
        Attribute("status", String, func() {
            Enum("pending", "confirmed", "processing", "shipped", "delivered", "cancelled", "refunded")
        })
        Attribute("fulfillment_status", String, func() {
            Enum("unfulfilled", "partial", "fulfilled")
        })
        Attribute("payment_status", String, func() {
            Enum("pending", "authorized", "paid", "partially_refunded", "refunded", "voided")
        })
        
        Attribute("line_items", ArrayOf(LineItem))
        
        Attribute("subtotal", Money)
        Attribute("shipping_total", Money)
        Attribute("tax_total", Money)
        Attribute("discount_total", Money)
        Attribute("total", Money)
        
        Attribute("shipping_address", Address)
        Attribute("billing_address", Address)
        
        Attribute("shipping_method", ShippingMethod)
        Attribute("tracking_number", String)
        Attribute("tracking_url", String)
        
        Attribute("notes", String, "Customer notes")
        Attribute("internal_notes", String, "Internal staff notes")
        
        Attribute("placed_at", String, func() {
            Format(FormatDateTime)
        })
        Attribute("shipped_at", String, func() {
            Format(FormatDateTime)
        })
        Attribute("delivered_at", String, func() {
            Format(FormatDateTime)
        })
        
        WithAuditFields()
    })
    
    View("default", func() {
        Attribute("id")
        Attribute("order_number")
        Attribute("status")
        Attribute("fulfillment_status")
        Attribute("payment_status")
        Attribute("line_items")
        Attribute("subtotal")
        Attribute("shipping_total")
        Attribute("tax_total")
        Attribute("discount_total")
        Attribute("total")
        Attribute("shipping_address")
        Attribute("tracking_number")
        Attribute("tracking_url")
        Attribute("placed_at")
    })
    
    View("summary", func() {
        Attribute("id")
        Attribute("order_number")
        Attribute("status")
        Attribute("total")
        Attribute("placed_at")
    })
    
    View("admin", func() {
        Attribute("id")
        Attribute("order_number")
        Attribute("customer_id")
        Attribute("customer", func() {
            View("summary")
        })
        Attribute("status")
        Attribute("fulfillment_status")
        Attribute("payment_status")
        Attribute("line_items")
        Attribute("subtotal")
        Attribute("shipping_total")
        Attribute("tax_total")
        Attribute("discount_total")
        Attribute("total")
        Attribute("shipping_address")
        Attribute("billing_address")
        Attribute("shipping_method")
        Attribute("tracking_number")
        Attribute("notes")
        Attribute("internal_notes")
        Attribute("placed_at")
        Attribute("shipped_at")
        Attribute("delivered_at")
        Attribute("created_at")
        Attribute("updated_at")
    })
    
    Required("id", "order_number", "status", "total")
})

var LineItem = Type("LineItem", func() {
    WithEntityID()
    
    Attribute("product_id", Int64)
    Attribute("variant_id", Int64)
    Attribute("product_name", String)
    Attribute("variant_name", String)
    Attribute("sku", String)
    Attribute("image_url", String)
    
    Attribute("quantity", Int32)
    Attribute("unit_price", Money)
    Attribute("total", Money)
    
    Attribute("discount_amount", Money)
    Attribute("tax_amount", Money)
    
    Required("id", "product_id", "product_name", "quantity", "unit_price", "total")
})
```

### Service Using Complete Type System

```go
// design/services.go
package design

import . "goa.design/goa/v3/dsl"

var _ = API("ecommerce", func() {
    Title("E-Commerce API")
    Version("1.0")
})

var _ = Service("products", func() {
    Description("Product management")
    
    HTTP(func() {
        Path("/products")
    })
    
    Method("list", func() {
        Description("List products with filtering")
        
        Payload(func() {
            Attribute("category_id", Int64)
            Attribute("status", String, func() {
                Enum("draft", "active", "archived", "all")
                Default("active")
            })
            Attribute("min_price", Float64)
            Attribute("max_price", Float64)
            Attribute("search", String)
            Attribute("page", Int32, func() {
                Default(1)
                Minimum(1)
            })
            Attribute("per_page", Int32, func() {
                Default(20)
                Minimum(1)
                Maximum(100)
            })
        })
        
        Result(func() {
            Attribute("products", CollectionOf(Product))
            Attribute("pagination", PaginationMeta)
            Required("products", "pagination")
        })
        
        HTTP(func() {
            GET("/")
            Param("category_id")
            Param("status")
            Param("min_price")
            Param("max_price")
            Param("search")
            Param("page")
            Param("per_page")
            Response(StatusOK)
        })
    })
    
    Method("get", func() {
        Payload(func() {
            Attribute("id", Int64)
            Required("id")
        })
        
        Result(Product)
        
        Error("not_found", ErrorResult, "Product not found")
        
        HTTP(func() {
            GET("/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
        })
    })
    
    Method("create", func() {
        Security(JWTAuth, func() {
            Scope("admin:products")
        })
        
        Payload(func() {
            Token("token", String)
            Attribute("product", CreateProductPayload)
            Required("product")
        })
        
        Result(Product, func() {
            View("admin")
        })
        
        HTTP(func() {
            POST("/")
            Body("product")
            Response(StatusCreated)
        })
    })
})

var _ = Service("orders", func() {
    Description("Order management")
    
    HTTP(func() {
        Path("/orders")
    })
    
    Method("list", func() {
        Description("List orders - view depends on user role")
        
        Security(JWTAuth)
        
        Payload(func() {
            Token("token", String)
            Attribute("status", String)
            Attribute("page", Int32, func() { Default(1) })
            Attribute("per_page", Int32, func() { Default(20) })
        })
        
        Result(func() {
            Attribute("orders", CollectionOf(Order))
            Attribute("pagination", PaginationMeta)
            Required("orders", "pagination")
        })
        
        HTTP(func() {
            GET("/")
            Param("status")
            Param("page")
            Param("per_page")
            Response(StatusOK)
        })
    })
    
    Method("get", func() {
        Security(JWTAuth)
        
        Payload(func() {
            Token("token", String)
            Attribute("id", Int64)
            Required("id")
        })
        
        // View selected based on user role in implementation
        Result(Order)
        
        Error("not_found", ErrorResult)
        Error("forbidden", ErrorResult)
        
        HTTP(func() {
            GET("/{id}")
            Response(StatusOK)
            Response("not_found", StatusNotFound)
            Response("forbidden", StatusForbidden)
        })
    })
})
```

---

## 📝 Summary

### Key Concepts

| Concept | Purpose | DSL |
|---------|---------|-----|
| **ResultType** | Typed API responses with views | `ResultType(mediaType, func)` |
| **View** | Projections of types | `View(name, func)` |
| **Reference** | Copy attributes from type | `Reference(type)` |
| **Extend** | Create subtype | `Extend(type)` |
| **OneOf** | Polymorphic union | `OneOf(name, func)` |
| **ArrayOf** | Collections | `ArrayOf(type)` |
| **MapOf** | Key-value maps | `MapOf(keyType, valueType)` |

### Best Practices

1. **Use Result Types for responses** - Take advantage of views
2. **Define reusable components** - Mixins, base types
3. **Use composition over inheritance** - More flexible
4. **Keep views focused** - Each view for a specific use case
5. **Version types carefully** - Plan for evolution
6. **Document with descriptions** - Self-documenting API

### Common Patterns

- **Base Entity Pattern**: Common ID, timestamps
- **Mixin Pattern**: Reusable attribute sets
- **View Pattern**: Role-based field visibility
- **Composition Pattern**: Build from components
- **Polymorphism Pattern**: OneOf for variants

---

## 📋 Knowledge Check

Before proceeding, ensure you can:

- [ ] Create Result Types with multiple views
- [ ] Use views for different response shapes
- [ ] Reuse types with Reference and Extend
- [ ] Create polymorphic types with OneOf
- [ ] Implement inheritance patterns
- [ ] Compose types from reusable components
- [ ] Define recursive/self-referential types
- [ ] Create paginated collection responses
- [ ] Apply views conditionally in services

---

## 🔗 Quick Reference Links

- [Goa DSL Reference](https://goa.design/reference/goa/v3/dsl/)
- [Result Types](https://goa.design/design/types/)
- [Views Documentation](https://goa.design/design/views/)
- [Type Reuse](https://goa.design/design/overview/)

---

> **Next Up:** Part 8 - Testing & Deployment (Unit Tests, Integration Tests, Mocking, Docker, Kubernetes)
