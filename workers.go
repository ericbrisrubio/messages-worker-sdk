package sdk

import (
	"context"
	"fmt"
	"net/http"
)

// WorkerInfo represents information about a single worker
type WorkerInfo struct {
	ID        string `json:"id"`
	QueueName string `json:"queue_name"`
	Status    string `json:"status"`
	StartedAt string `json:"started_at"`
}

// PriorityWorkerInfo represents worker information for a specific priority
type PriorityWorkerInfo struct {
	Count      int          `json:"count"`
	QueueDepth int          `json:"queue_depth"`
	Workers    []WorkerInfo `json:"workers"`
}

// WorkerStatusResponse represents the response from the worker status endpoint
type WorkerStatusResponse struct {
	TotalWorkers    int                    `json:"total_workers"`
	LowPriority     PriorityWorkerInfo    `json:"low_priority"`
	MediumPriority  PriorityWorkerInfo    `json:"medium_priority"`
	HighPriority    PriorityWorkerInfo    `json:"high_priority"`
	AllWorkers      []WorkerInfo          `json:"all_workers"`
}

// ScaleWorkersRequest represents a request to scale workers
type ScaleWorkersRequest struct {
	Priority string `json:"priority"`
	Count    int    `json:"count"`
}

// ScaleWorkersResponse represents the response from scaling workers
type ScaleWorkersResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Priority string `json:"priority"`
	Count    int    `json:"count"`
	Action   string `json:"action"`
}

// RemoveAllWorkersResponse represents the response from removing all workers
type RemoveAllWorkersResponse struct {
	Status        string   `json:"status"`
	Message       string   `json:"message"`
	TotalRemoved  int      `json:"total_removed"`
	Errors        []string `json:"errors,omitempty"`
}

// GetWorkerStatus returns the current status of all workers
func (c *Client) GetWorkerStatus(ctx context.Context) (*WorkerStatusResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/api/v1/workers/status", nil)
	if err != nil {
		return nil, err
	}

	var statusResp WorkerStatusResponse
	if err := c.parseResponse(resp, &statusResp); err != nil {
		return nil, err
	}

	return &statusResp, nil
}

// ScaleWorkers scales workers for a specific priority queue
func (c *Client) ScaleWorkers(ctx context.Context, priority string, count int) (*ScaleWorkersResponse, error) {
	if priority == "" {
		return nil, fmt.Errorf("priority is required")
	}

	if priority != "low" && priority != "medium" && priority != "high" {
		return nil, fmt.Errorf("priority must be 'low', 'medium', or 'high'")
	}

	if count == 0 {
		return nil, fmt.Errorf("count cannot be 0")
	}

	path := fmt.Sprintf("/api/v1/workers/scale/%s?count=%d", priority, count)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	var scaleResp ScaleWorkersResponse
	if err := c.parseResponse(resp, &scaleResp); err != nil {
		return nil, err
	}

	return &scaleResp, nil
}

// AddWorkers adds workers for a specific priority queue
func (c *Client) AddWorkers(ctx context.Context, priority string, count int) (*ScaleWorkersResponse, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	return c.ScaleWorkers(ctx, priority, count)
}

// RemoveWorkers removes workers for a specific priority queue
func (c *Client) RemoveWorkers(ctx context.Context, priority string, count int) (*ScaleWorkersResponse, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	return c.ScaleWorkers(ctx, priority, -count)
}

// RemoveAllWorkers removes all running workers across all priority queues
func (c *Client) RemoveAllWorkers(ctx context.Context) (*RemoveAllWorkersResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/api/v1/workers/remove-all", nil)
	if err != nil {
		return nil, err
	}

	var removeResp RemoveAllWorkersResponse
	if err := c.parseResponse(resp, &removeResp); err != nil {
		return nil, err
	}

	return &removeResp, nil
}

// GetWorkerCount returns the number of workers for a specific priority
func (c *Client) GetWorkerCount(ctx context.Context, priority string) (int, error) {
	status, err := c.GetWorkerStatus(ctx)
	if err != nil {
		return 0, err
	}

	switch priority {
	case "low":
		return status.LowPriority.Count, nil
	case "medium":
		return status.MediumPriority.Count, nil
	case "high":
		return status.HighPriority.Count, nil
	default:
		return 0, fmt.Errorf("invalid priority: %s", priority)
	}
}

// GetTotalWorkerCount returns the total number of workers across all priorities
func (c *Client) GetTotalWorkerCount(ctx context.Context) (int, error) {
	status, err := c.GetWorkerStatus(ctx)
	if err != nil {
		return 0, err
	}

	return status.TotalWorkers, nil
}

