package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ericbrisrubio/messages-worker-sdk"
)

func main() {
	// Create a client with default configuration
	client := sdk.NewClientWithDefaults()

	// Or create a client with custom configuration
	ctx := context.Background()

	// Example 1: Health Check
	fmt.Println("=== Health Check ===")
	health, err := client.CheckHealth(ctx)
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Printf("Service is healthy: %s\n", health.Status)
	}

	// Example 2: Submit a single message
	fmt.Println("\n=== Single Message Submission ===")
	messageReq := &sdk.MessageRequest{
		ItemID:      "pr-123",
		Priority:    sdk.PriorityHigh,
		Topic:       sdk.TopicPullRequests,
		CallbackURL: "https://httpbin.org/post",
		ObjectBody: map[string]interface{}{
			"pull_request": map[string]interface{}{
				"id":         123,
				"title":      "Add new feature",
				"author":     "john.doe",
				"repository": "myorg/myrepo",
				"status":     "open",
			},
		},
	}

	messageResp, err := client.PostMessage(ctx, messageReq)
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	} else {
		fmt.Printf("Message posted successfully: ID=%s, Status=%s\n", messageResp.ID, messageResp.Status)
	}

	// Example 3: Submit bulk messages
	fmt.Println("\n=== Bulk Message Submission ===")
	bulkReq := &sdk.BulkMessageRequest{
		Messages: []sdk.MessageRequest{
			{
				ItemID:      "pr-124",
				Priority:    sdk.PriorityMedium,
				Topic:       sdk.TopicPullRequests,
				CallbackURL: "https://httpbin.org/post",
				ObjectBody: map[string]interface{}{
					"pull_request": map[string]interface{}{
						"id":         124,
						"title":      "Fix bug in authentication",
						"author":     "jane.smith",
						"repository": "myorg/myrepo",
						"status":     "open",
					},
				},
			},
			{
				ItemID:      "pr-125",
				Priority:    sdk.PriorityLow,
				Topic:       sdk.TopicPullRequests,
				CallbackURL: "https://httpbin.org/post",
				ObjectBody: map[string]interface{}{
					"pull_request": map[string]interface{}{
						"id":         125,
						"title":      "Update documentation",
						"author":     "bob.wilson",
						"repository": "myorg/myrepo",
						"status":     "open",
					},
				},
			},
		},
	}

	bulkResp, err := client.PostBulkMessages(ctx, bulkReq)
	if err != nil {
		log.Printf("Failed to post bulk messages: %v", err)
	} else {
		fmt.Printf("Bulk messages posted successfully: Count=%d, Status=%s\n", bulkResp.Count, bulkResp.Status)
		for i, msg := range bulkResp.Messages {
			fmt.Printf("  Message %d: ID=%s, Priority=%s\n", i+1, msg.ID, msg.Priority)
		}
	}

	// Example 4: Convenience methods
	fmt.Println("\n=== Convenience Methods ===")

	// Post a high priority message
	highPriorityResp, err := client.PostHighPriorityMessage(ctx, "pr-126", "https://httpbin.org/post", map[string]interface{}{
		"urgent": true,
		"type":   "security_fix",
	})
	if err != nil {
		log.Printf("Failed to post high priority message: %v", err)
	} else {
		fmt.Printf("High priority message posted: ID=%s\n", highPriorityResp.ID)
	}

	// Post a low priority message
	lowPriorityResp, err := client.PostLowPriorityMessage(ctx, "pr-127", "https://httpbin.org/post", map[string]interface{}{
		"type": "documentation",
	})
	if err != nil {
		log.Printf("Failed to post low priority message: %v", err)
	} else {
		fmt.Printf("Low priority message posted: ID=%s\n", lowPriorityResp.ID)
	}

	// Example 5: Worker Management
	fmt.Println("\n=== Worker Management ===")

	// Get worker status
	status, err := client.GetWorkerStatus(ctx)
	if err != nil {
		log.Printf("Failed to get worker status: %v", err)
	} else {
		fmt.Printf("Total workers: %d\n", status.TotalWorkers)
		fmt.Printf("Low priority workers: %d\n", status.LowPriority.Count)
		fmt.Printf("Medium priority workers: %d\n", status.MediumPriority.Count)
		fmt.Printf("High priority workers: %d\n", status.HighPriority.Count)
	}

	// Scale workers (add 2 high priority workers)
	scaleResp, err := client.AddWorkers(ctx, "high", 2)
	if err != nil {
		log.Printf("Failed to add workers: %v", err)
	} else {
		fmt.Printf("Workers scaled: %s\n", scaleResp.Message)
	}

	// Get worker count for specific priority
	count, err := client.GetWorkerCount(ctx, "high")
	if err != nil {
		log.Printf("Failed to get worker count: %v", err)
	} else {
		fmt.Printf("High priority worker count: %d\n", count)
	}

	// Example 6: Error Handling
	fmt.Println("\n=== Error Handling ===")

	// Try to post a message with invalid priority
	invalidReq := &sdk.MessageRequest{
		ItemID:      "pr-invalid",
		Priority:    "invalid", // This will cause an error
		Topic:       sdk.TopicPullRequests,
		CallbackURL: "https://httpbin.org/post",
		ObjectBody:  map[string]interface{}{"test": true},
	}

	_, err = client.PostMessage(ctx, invalidReq)
	if err != nil {
		if sdk.IsAPIError(err) {
			fmt.Printf("API Error: %v\n", err)
		} else {
			fmt.Printf("Other Error: %v\n", err)
		}
	}

	fmt.Println("\n=== Example completed ===")
}
