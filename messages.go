package sdk

import (
	"context"
	"fmt"
	"net/http"
)

// Priority represents the priority level of a message
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Topic represents the topic/category of a message
type Topic string

const (
	TopicPullRequests Topic = "pullrequests"
)

// MessageRequest represents a single message request
type MessageRequest struct {
	ItemID      string      `json:"item_id"`
	Priority    Priority    `json:"priority"`
	Topic       Topic       `json:"topic"`
	CallbackURL string      `json:"callback_url"`
	ObjectBody  interface{} `json:"object_body"`
}

// MessageResponse represents the response for a single message
type MessageResponse struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	ItemID   string   `json:"itemId"`
	Priority Priority `json:"priority"`
	Topic    Topic    `json:"topic"`
}

// BulkMessageRequest represents a request to post multiple messages
type BulkMessageRequest struct {
	Messages []MessageRequest `json:"messages"`
}

// BulkMessageResponse represents the response for bulk messages
type BulkMessageResponse struct {
	Status   string            `json:"status"`
	Count    int               `json:"count"`
	Messages []MessageResponse `json:"messages"`
}

// PostMessage submits a single message for processing
func (c *Client) PostMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("message request cannot be nil")
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/api/v1/messages", req)
	if err != nil {
		return nil, err
	}

	var messageResp MessageResponse
	if err := c.parseResponse(resp, &messageResp); err != nil {
		return nil, err
	}

	return &messageResp, nil
}

// PostBulkMessages submits multiple messages for processing
func (c *Client) PostBulkMessages(ctx context.Context, req *BulkMessageRequest) (*BulkMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("bulk message request cannot be nil")
	}

	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/api/v1/messages/bulk", req)
	if err != nil {
		return nil, err
	}

	var bulkResp BulkMessageResponse
	if err := c.parseResponse(resp, &bulkResp); err != nil {
		return nil, err
	}

	return &bulkResp, nil
}

// PostMessageWithDefaults creates a message request with default values and submits it
func (c *Client) PostMessageWithDefaults(ctx context.Context, itemID string, callbackURL string, objectBody interface{}) (*MessageResponse, error) {
	req := &MessageRequest{
		ItemID:      itemID,
		Priority:    PriorityMedium,
		Topic:       TopicPullRequests,
		CallbackURL: callbackURL,
		ObjectBody:  objectBody,
	}

	return c.PostMessage(ctx, req)
}

// PostHighPriorityMessage submits a high priority message
func (c *Client) PostHighPriorityMessage(ctx context.Context, itemID string, callbackURL string, objectBody interface{}) (*MessageResponse, error) {
	req := &MessageRequest{
		ItemID:      itemID,
		Priority:    PriorityHigh,
		Topic:       TopicPullRequests,
		CallbackURL: callbackURL,
		ObjectBody:  objectBody,
	}

	return c.PostMessage(ctx, req)
}

// PostLowPriorityMessage submits a low priority message
func (c *Client) PostLowPriorityMessage(ctx context.Context, itemID string, callbackURL string, objectBody interface{}) (*MessageResponse, error) {
	req := &MessageRequest{
		ItemID:      itemID,
		Priority:    PriorityLow,
		Topic:       TopicPullRequests,
		CallbackURL: callbackURL,
		ObjectBody:  objectBody,
	}

	return c.PostMessage(ctx, req)
}

