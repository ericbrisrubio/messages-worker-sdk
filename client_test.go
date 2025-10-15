package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// Test with default config
	client := NewClientWithDefaults()
	if client.baseURL != "http://localhost:8083" {
		t.Errorf("Expected baseURL to be 'http://localhost:8083', got '%s'", client.baseURL)
	}
	if client.timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", client.timeout)
	}

	// Test with custom config
	config := &Config{
		BaseURL: "https://example.com",
		Timeout: 60 * time.Second,
	}
	client = NewClient(config)
	if client.baseURL != "https://example.com" {
		t.Errorf("Expected baseURL to be 'https://example.com', got '%s'", client.baseURL)
	}
	if client.timeout != 60*time.Second {
		t.Errorf("Expected timeout to be 60s, got %v", client.timeout)
	}

	// Test with nil config
	client = NewClient(nil)
	if client.baseURL != "http://localhost:8083" {
		t.Errorf("Expected baseURL to be 'http://localhost:8083' with nil config, got '%s'", client.baseURL)
	}
}

func TestCheckHealth(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected path '/health', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected method 'GET', got '%s'", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Create client with test server URL
	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	// Test health check
	ctx := context.Background()
	health, err := client.CheckHealth(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if health.Status != "OK" {
		t.Errorf("Expected status 'OK', got '%s'", health.Status)
	}
}

func TestCheckHealthError(t *testing.T) {
	// Create a test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Service Unavailable"))
	}))
	defer server.Close()

	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	ctx := context.Background()
	_, err := client.CheckHealth(ctx)
	if err == nil {
		t.Error("Expected error for 500 status, got nil")
	}

	if !IsAPIError(err) {
		t.Error("Expected APIError, got different error type")
	}

	apiErr := err.(*APIError)
	if apiErr.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", apiErr.StatusCode)
	}
}

func TestPostMessage(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/messages" {
			t.Errorf("Expected path '/api/v1/messages', got '%s'", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected method 'POST', got '%s'", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var req MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Validate request
		if req.ItemID != "test-123" {
			t.Errorf("Expected ItemID 'test-123', got '%s'", req.ItemID)
		}
		if req.Priority != PriorityHigh {
			t.Errorf("Expected Priority 'high', got '%s'", req.Priority)
		}

		// Return response
		response := MessageResponse{
			ID:       "msg-123",
			Status:   "published",
			ItemID:   req.ItemID,
			Priority: req.Priority,
			Topic:    req.Topic,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	// Test message submission
	ctx := context.Background()
	req := &MessageRequest{
		ItemID:      "test-123",
		Priority:    PriorityHigh,
		Topic:       TopicPullRequests,
		CallbackURL: "https://example.com/callback",
		ObjectBody:  map[string]interface{}{"test": true},
	}

	resp, err := client.PostMessage(ctx, req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.ID != "msg-123" {
		t.Errorf("Expected ID 'msg-123', got '%s'", resp.ID)
	}
	if resp.Status != "published" {
		t.Errorf("Expected Status 'published', got '%s'", resp.Status)
	}
}

func TestPostMessageValidation(t *testing.T) {
	client := NewClientWithDefaults()
	ctx := context.Background()

	// Test nil request
	_, err := client.PostMessage(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}
	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Expected error message about nil request, got: %v", err)
	}
}

func TestPostBulkMessages(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/messages/bulk" {
			t.Errorf("Expected path '/api/v1/messages/bulk', got '%s'", r.URL.Path)
		}

		// Parse request body
		var req BulkMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if len(req.Messages) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(req.Messages))
		}

		// Return response
		response := BulkMessageResponse{
			Status:   "published",
			Count:    2,
			Messages: []MessageResponse{
				{ID: "msg-1", Status: "published", ItemID: "test-1", Priority: PriorityHigh, Topic: TopicPullRequests},
				{ID: "msg-2", Status: "published", ItemID: "test-2", Priority: PriorityMedium, Topic: TopicPullRequests},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	// Test bulk message submission
	ctx := context.Background()
	req := &BulkMessageRequest{
		Messages: []MessageRequest{
			{
				ItemID:      "test-1",
				Priority:    PriorityHigh,
				Topic:       TopicPullRequests,
				CallbackURL: "https://example.com/callback1",
				ObjectBody:  map[string]interface{}{"test": 1},
			},
			{
				ItemID:      "test-2",
				Priority:    PriorityMedium,
				Topic:       TopicPullRequests,
				CallbackURL: "https://example.com/callback2",
				ObjectBody:  map[string]interface{}{"test": 2},
			},
		},
	}

	resp, err := client.PostBulkMessages(ctx, req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Count != 2 {
		t.Errorf("Expected count 2, got %d", resp.Count)
	}
	if len(resp.Messages) != 2 {
		t.Errorf("Expected 2 messages in response, got %d", len(resp.Messages))
	}
}

func TestGetWorkerStatus(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/workers/status" {
			t.Errorf("Expected path '/api/v1/workers/status', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected method 'GET', got '%s'", r.Method)
		}

		// Return mock worker status
		status := WorkerStatusResponse{
			TotalWorkers: 6,
			LowPriority: PriorityWorkerInfo{
				Count:      2,
				QueueDepth: 5,
				Workers: []WorkerInfo{
					{ID: "low-1", QueueName: "low_priority", Status: "running", StartedAt: "2024-01-01T00:00:00Z"},
				},
			},
			MediumPriority: PriorityWorkerInfo{
				Count:      2,
				QueueDepth: 0,
				Workers:    []WorkerInfo{},
			},
			HighPriority: PriorityWorkerInfo{
				Count:      2,
				QueueDepth: 3,
				Workers:    []WorkerInfo{},
			},
			AllWorkers: []WorkerInfo{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}))
	defer server.Close()

	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	ctx := context.Background()
	status, err := client.GetWorkerStatus(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if status.TotalWorkers != 6 {
		t.Errorf("Expected total workers 6, got %d", status.TotalWorkers)
	}
	if status.LowPriority.Count != 2 {
		t.Errorf("Expected low priority workers 2, got %d", status.LowPriority.Count)
	}
}

func TestScaleWorkers(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/v1/workers/scale/") {
			t.Errorf("Expected path to start with '/api/v1/workers/scale/', got '%s'", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected method 'POST', got '%s'", r.Method)
		}

		// Extract priority from path
		parts := strings.Split(r.URL.Path, "/")
		priority := parts[len(parts)-1]
		if priority != "high" {
			t.Errorf("Expected priority 'high', got '%s'", priority)
		}

		// Extract count from query
		count := r.URL.Query().Get("count")
		if count != "2" {
			t.Errorf("Expected count '2', got '%s'", count)
		}

		// Return response
		response := ScaleWorkersResponse{
			Status:   "success",
			Message:  "Successfully added 2 high priority workers",
			Priority: "high",
			Count:    2,
			Action:   "added",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	ctx := context.Background()
	resp, err := client.ScaleWorkers(ctx, "high", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Priority != "high" {
		t.Errorf("Expected priority 'high', got '%s'", resp.Priority)
	}
	if resp.Count != 2 {
		t.Errorf("Expected count 2, got %d", resp.Count)
	}
	if resp.Action != "added" {
		t.Errorf("Expected action 'added', got '%s'", resp.Action)
	}
}

func TestScaleWorkersValidation(t *testing.T) {
	client := NewClientWithDefaults()
	ctx := context.Background()

	// Test empty priority
	_, err := client.ScaleWorkers(ctx, "", 1)
	if err == nil {
		t.Error("Expected error for empty priority, got nil")
	}

	// Test invalid priority
	_, err = client.ScaleWorkers(ctx, "invalid", 1)
	if err == nil {
		t.Error("Expected error for invalid priority, got nil")
	}

	// Test zero count
	_, err = client.ScaleWorkers(ctx, "high", 0)
	if err == nil {
		t.Error("Expected error for zero count, got nil")
	}
}

func TestIsHealthy(t *testing.T) {
	// Test healthy service
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewClient(config)

	ctx := context.Background()
	if !client.IsHealthy(ctx) {
		t.Error("Expected service to be healthy")
	}

	// Test unhealthy service
	unhealthyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer unhealthyServer.Close()

	unhealthyConfig := &Config{
		BaseURL: unhealthyServer.URL,
		Timeout: 5 * time.Second,
	}
	unhealthyClient := NewClient(unhealthyConfig)

	if unhealthyClient.IsHealthy(ctx) {
		t.Error("Expected service to be unhealthy")
	}
}
