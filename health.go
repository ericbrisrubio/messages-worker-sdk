package sdk

import (
	"context"
	"fmt"
	"net/http"
)

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Status string `json:"status"`
}

// CheckHealth checks if the messages-worker service is healthy
func (c *Client) CheckHealth(ctx context.Context) (*HealthResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/health", nil)
	if err != nil {
		return nil, err
	}

	// The health endpoint returns "OK" as plain text, not JSON
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("health check failed with status %d", resp.StatusCode),
		}
	}

	return &HealthResponse{
		Status: "OK",
	}, nil
}

// IsHealthy returns true if the service is healthy, false otherwise
func (c *Client) IsHealthy(ctx context.Context) bool {
	_, err := c.CheckHealth(ctx)
	return err == nil
}

// Ping is an alias for CheckHealth for convenience
func (c *Client) Ping(ctx context.Context) (*HealthResponse, error) {
	return c.CheckHealth(ctx)
}

