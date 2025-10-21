package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVerifyCallbackSignature(t *testing.T) {
	secret := "test-secret"
	payload := `{"id":"test-id","item_id":"pr-123","priority":"high"}`

	tests := []struct {
		name          string
		payload       string
		signature     string
		secret        string
		expected      bool
		expectedError bool
	}{
		{
			name:          "valid signature",
			payload:       payload,
			signature:     generateTestSignature(payload, secret),
			secret:        secret,
			expected:      true,
			expectedError: false,
		},
		{
			name:          "missing signature",
			payload:       payload,
			signature:     "",
			secret:        secret,
			expected:      false,
			expectedError: true,
		},
		{
			name:          "invalid signature",
			payload:       payload,
			signature:     "invalidhash1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			secret:        secret,
			expected:      false,
			expectedError: false,
		},
		{
			name:          "empty secret",
			payload:       payload,
			signature:     generateTestSignature(payload, secret),
			secret:        "",
			expected:      false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifySignatureSha256([]byte(tt.payload), tt.signature, tt.secret)

			if tt.expectedError {
				if result == true {
					t.Errorf("VerifyCallbackSignature() expected error, got nil")
				}
				return
			}

			if result != tt.expected {
				t.Errorf("VerifyCallbackSignature() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGetSignatureHeader(t *testing.T) {
	expected := "X-Callback-Signature"
	result := GetSignatureHeader()
	if result != expected {
		t.Errorf("GetSignatureHeader() = %v, expected %v", result, expected)
	}
}

// Helper function to generate test signatures
func generateTestSignature(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
