package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ericbrisrubio/messages-worker-sdk"
)

// CallbackHandler handles incoming callbacks from messages-worker
type CallbackHandler struct {
	secret string
}

// NewCallbackHandler creates a new callback handler
func NewCallbackHandler(secret string) *CallbackHandler {
	return &CallbackHandler{
		secret: secret,
	}
}

// HandleCallback processes incoming callback requests
func (h *CallbackHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify signature if secret is configured
	if h.secret != "" {
		signature := r.Header.Get(sdk.GetSignatureHeader())
		valid := sdk.VerifySignatureSha256(body, signature, h.secret)
		if !valid {
			log.Printf("Invalid signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
		log.Printf("Signature verified successfully")
	}

	// Parse the message
	var message map[string]interface{}
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("Failed to decode message: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process the message
	log.Printf("Received callback message: %+v", message)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Callback processed successfully",
	})
}

func main() {
	// Get callback secret from environment
	secret := os.Getenv("CALLBACK_SECRET")
	if secret == "" {
		log.Println("Warning: CALLBACK_SECRET not set, callbacks will not be verified")
	}

	// Create callback handler
	handler := NewCallbackHandler(secret)

	// Setup routes
	http.HandleFunc("/callback", handler.HandleCallback)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting callback server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
