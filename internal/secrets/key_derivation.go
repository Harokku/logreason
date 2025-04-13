// Package secrets provides secure handling of API keys and other sensitive information.
package secrets

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

const (
	// DefaultIterations is the default number of iterations for PBKDF2
	DefaultIterations = 10000
	// DefaultSaltLength is the default length of the salt in bytes
	DefaultSaltLength = 16
	// DefaultKeyLength is the default length of the derived key in bytes (32 bytes = 256 bits for AES-256)
	DefaultKeyLength = 32
)

// DeriveKeyFromPassword derives an encryption key from a password using PBKDF2
// This is useful when you want to use a password to encrypt/decrypt secrets
// rather than managing a raw encryption key
func DeriveKeyFromPassword(password string, salt []byte, iterations, keyLength int) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
}

// GenerateSalt generates a random salt for key derivation
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// EncodeKeyToString encodes a binary key to a base64 string for storage
func EncodeKeyToString(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// DecodeKeyFromString decodes a base64 string back to a binary key
func DecodeKeyFromString(encodedKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedKey)
}

// DeriveKeyWithDefaults derives an encryption key from a password using default parameters
// It generates a new salt and returns both the derived key and the salt
func DeriveKeyWithDefaults(password string) (key, salt []byte, err error) {
	salt, err = GenerateSalt(DefaultSaltLength)
	if err != nil {
		return nil, nil, err
	}

	key = DeriveKeyFromPassword(password, salt, DefaultIterations, DefaultKeyLength)
	return key, salt, nil
}

// DeriveKeyWithSalt derives an encryption key from a password using a provided salt
// This should be used when you have stored the salt from a previous key derivation
func DeriveKeyWithSalt(password string, salt []byte) []byte {
	return DeriveKeyFromPassword(password, salt, DefaultIterations, DefaultKeyLength)
}
