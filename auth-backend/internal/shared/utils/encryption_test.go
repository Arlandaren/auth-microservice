package utils

import (
	"encoding/base64"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	originalText := "This is a secret message."
	keyBase64 := "wRq8N2Nz8D2jblorVtR/4UxEofcU11dRTB3lIy9bB+M="

	if keyBase64 == "" {
		t.Fatalf("KEY environment variable is not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		t.Fatalf("invalid key: %v", err)

	}

	if len(key) != 32 {
		t.Fatalf("invalid key size: %d bytes. Expected 32 bytes for AES-256.", len(key))
	}

	encryptedText, err := Encrypt(originalText, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decryptedText, err := Decrypt(encryptedText, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decryptedText != originalText {
		t.Fatalf("Decrypted text does not match original. Got '%s', expected '%s'", decryptedText, originalText)
	}
}
