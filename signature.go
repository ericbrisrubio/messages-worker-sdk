package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// VerifyCallbackSignature verifies the callback signature using HMAC-SHA256
func VerifySignatureSha256(payload []byte, signature string, webhookSecret string) bool {
	if signature == "" {
		return false
	}

	// Remove "sha256=" prefix if present
	signature = strings.TrimPrefix(signature, "sha256=")

	// Create HMAC-SHA256 hash
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)

	// Convert to hex string
	expectedSignature := hex.EncodeToString(expectedMAC)

	// Compare signatures using constant time comparison
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// GetSignatureHeader returns the header name used for callback signatures
func GetSignatureHeader() string {
	return "X-Callback-Signature"
}
