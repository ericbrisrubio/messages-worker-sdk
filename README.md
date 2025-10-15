# Messages Worker SDK

A Go SDK client for interacting with the messages-worker service. This SDK provides a convenient way to submit messages, manage workers, and check service health from Go applications.

## Features

- **Message Operations**: Submit single or bulk messages with different priorities
- **Worker Management**: Monitor and scale workers dynamically
- **Health Checks**: Check service health and availability
- **Error Handling**: Comprehensive error handling with custom error types
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Type Safety**: Strongly typed API with proper validation

## Installation

```bash
go get github.com/ericbrisrubio/messages-worker-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ericbrisrubio/messages-worker-sdk"
)

func main() {
    // Create a client
    client := sdk.NewClientWithDefaults()
    
    ctx := context.Background()
    
    // Submit a message
    resp, err := client.PostHighPriorityMessage(ctx, "pr-123", "https://example.com/callback", map[string]interface{}{
        "pull_request": map[string]interface{}{
            "id": 123,
            "title": "Add new feature",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Message submitted: %s\n", resp.ID)
}
```

## Client Configuration

### Default Configuration

```go
client := sdk.NewClientWithDefaults()
// Uses: http://localhost:8083 with 30s timeout
```

### Custom Configuration

```go
config := &sdk.Config{
    BaseURL: "https://messages-worker.example.com",
    Timeout: 60 * time.Second,
}
client := sdk.NewClient(config)
```

## Message Operations

### Single Message Submission

```go
messageReq := &sdk.MessageRequest{
    ItemID:      "pr-123",
    Priority:    sdk.PriorityHigh,
    Topic:       sdk.TopicPullRequests,
    CallbackURL: "https://example.com/callback",
    ObjectBody: map[string]interface{}{
        "pull_request": map[string]interface{}{
            "id": 123,
            "title": "Add new feature",
        },
    },
}

resp, err := client.PostMessage(ctx, messageReq)
```

### Bulk Message Submission

```go
bulkReq := &sdk.BulkMessageRequest{
    Messages: []sdk.MessageRequest{
        {
            ItemID:      "pr-123",
            Priority:    sdk.PriorityHigh,
            Topic:       sdk.TopicPullRequests,
            CallbackURL: "https://example.com/callback1",
            ObjectBody:  map[string]interface{}{"data": "value1"},
        },
        {
            ItemID:      "pr-124",
            Priority:    sdk.PriorityMedium,
            Topic:       sdk.TopicPullRequests,
            CallbackURL: "https://example.com/callback2",
            ObjectBody:  map[string]interface{}{"data": "value2"},
        },
    },
}

resp, err := client.PostBulkMessages(ctx, bulkReq)
```

### Convenience Methods

```go
// High priority message
resp, err := client.PostHighPriorityMessage(ctx, "pr-123", "https://example.com/callback", data)

// Medium priority message (default)
resp, err := client.PostMessageWithDefaults(ctx, "pr-123", "https://example.com/callback", data)

// Low priority message
resp, err := client.PostLowPriorityMessage(ctx, "pr-123", "https://example.com/callback", data)
```

## Worker Management

### Get Worker Status

```go
status, err := client.GetWorkerStatus(ctx)
fmt.Printf("Total workers: %d\n", status.TotalWorkers)
fmt.Printf("High priority workers: %d\n", status.HighPriority.Count)
```

### Scale Workers

```go
// Add workers
resp, err := client.AddWorkers(ctx, "high", 3)

// Remove workers
resp, err := client.RemoveWorkers(ctx, "medium", 2)

// Scale workers (positive to add, negative to remove)
resp, err := client.ScaleWorkers(ctx, "low", -1)
```

### Remove All Workers

```go
resp, err := client.RemoveAllWorkers(ctx)
fmt.Printf("Removed %d workers\n", resp.TotalRemoved)
```

### Get Worker Counts

```go
// Get count for specific priority
count, err := client.GetWorkerCount(ctx, "high")

// Get total count
total, err := client.GetTotalWorkerCount(ctx)
```

## Health Checks

### Check Service Health

```go
health, err := client.CheckHealth(ctx)
if err != nil {
    log.Printf("Service is unhealthy: %v", err)
} else {
    fmt.Printf("Service is healthy: %s\n", health.Status)
}
```

### Simple Health Check

```go
if client.IsHealthy(ctx) {
    fmt.Println("Service is healthy")
} else {
    fmt.Println("Service is unhealthy")
}
```

### Ping Service

```go
resp, err := client.Ping(ctx)
```

## Error Handling

The SDK provides comprehensive error handling with custom error types:

```go
resp, err := client.PostMessage(ctx, messageReq)
if err != nil {
    if sdk.IsAPIError(err) {
        // Handle API errors (4xx, 5xx responses)
        apiErr := err.(*sdk.APIError)
        fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
    } else {
        // Handle other errors (network, parsing, etc.)
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Error Types

- **APIError**: Errors returned by the API (HTTP 4xx, 5xx)
- **Network errors**: Connection failures, timeouts
- **Validation errors**: Invalid request parameters
- **Parsing errors**: JSON marshaling/unmarshaling failures

## Message Priorities

The SDK supports three priority levels:

- `sdk.PriorityLow`: Low priority messages (30s delay)
- `sdk.PriorityMedium`: Medium priority messages (15s delay)  
- `sdk.PriorityHigh`: High priority messages (5s delay)

## Topics

Currently supported topics:

- `sdk.TopicPullRequests`: Pull request related messages

## Context Support

All SDK methods support `context.Context` for timeouts and cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.PostMessage(ctx, messageReq)
```

## Examples

See the `examples/` directory for complete working examples:

- `examples/main.go`: Comprehensive example showing all SDK features

## API Reference

### Client Methods

#### Message Operations
- `PostMessage(ctx, req)` - Submit a single message
- `PostBulkMessages(ctx, req)` - Submit multiple messages
- `PostMessageWithDefaults(ctx, itemID, callbackURL, objectBody)` - Submit with defaults
- `PostHighPriorityMessage(ctx, itemID, callbackURL, objectBody)` - Submit high priority
- `PostLowPriorityMessage(ctx, itemID, callbackURL, objectBody)` - Submit low priority

#### Worker Management
- `GetWorkerStatus(ctx)` - Get current worker status
- `ScaleWorkers(ctx, priority, count)` - Scale workers (positive/negative count)
- `AddWorkers(ctx, priority, count)` - Add workers
- `RemoveWorkers(ctx, priority, count)` - Remove workers
- `RemoveAllWorkers(ctx)` - Remove all workers
- `GetWorkerCount(ctx, priority)` - Get worker count for priority
- `GetTotalWorkerCount(ctx)` - Get total worker count

#### Health Checks
- `CheckHealth(ctx)` - Check service health
- `IsHealthy(ctx)` - Simple boolean health check
- `Ping(ctx)` - Alias for CheckHealth

### Types

#### Message Types
- `MessageRequest` - Single message request
- `MessageResponse` - Single message response
- `BulkMessageRequest` - Bulk message request
- `BulkMessageResponse` - Bulk message response

#### Worker Types
- `WorkerInfo` - Individual worker information
- `PriorityWorkerInfo` - Worker info for a priority level
- `WorkerStatusResponse` - Complete worker status
- `ScaleWorkersResponse` - Worker scaling response
- `RemoveAllWorkersResponse` - Remove all workers response

#### Health Types
- `HealthResponse` - Health check response

#### Configuration Types
- `Config` - Client configuration
- `APIError` - API error type

## License

This SDK is part of the messages-worker project and follows the same license terms.
