// Package secrets provides secure handling of API keys and other sensitive information.
package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// EncryptedSecrets represents the structure of the encrypted secrets file
type EncryptedSecrets struct {
	Nonce   string `json:"nonce"`
	Secrets string `json:"secrets"`
}

// LoadEncryptedFromFile loads and decrypts secrets from an encrypted file
// The encryption key should be a 32-byte key (for AES-256)
func (m *Manager) LoadEncryptedFromFile(filePath string, encryptionKey []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("encrypted secrets file does not exist: %s", filePath)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted secrets file: %w", err)
	}

	// Parse the JSON
	var encryptedSecrets EncryptedSecrets
	if err := json.Unmarshal(data, &encryptedSecrets); err != nil {
		return fmt.Errorf("failed to parse encrypted secrets file: %w", err)
	}

	// Decode the nonce
	nonce, err := base64.StdEncoding.DecodeString(encryptedSecrets.Nonce)
	if err != nil {
		return fmt.Errorf("failed to decode nonce: %w", err)
	}

	// Decode the encrypted data
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedSecrets.Secrets)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	// Decrypt the data
	decryptedData, err := decrypt(encryptedData, encryptionKey, nonce)
	if err != nil {
		return fmt.Errorf("failed to decrypt secrets: %w", err)
	}

	// Parse the decrypted JSON
	var secrets map[string]string
	if err := json.Unmarshal(decryptedData, &secrets); err != nil {
		return fmt.Errorf("failed to parse decrypted secrets: %w", err)
	}

	// Add the secrets to our map
	for k, v := range secrets {
		m.secrets[k] = v
	}

	return nil
}

// SaveEncryptedToFile encrypts and saves the current secrets to a file
// The encryption key should be a 32-byte key (for AES-256)
func (m *Manager) SaveEncryptedToFile(filePath string, encryptionKey []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal the secrets to JSON
	secretsJSON, err := json.Marshal(m.secrets)
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, 12) // 96 bits is recommended for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	encryptedData, err := encrypt(secretsJSON, encryptionKey, nonce)
	if err != nil {
		return fmt.Errorf("failed to encrypt secrets: %w", err)
	}

	// Create the encrypted secrets structure
	encryptedSecrets := EncryptedSecrets{
		Nonce:   base64.StdEncoding.EncodeToString(nonce),
		Secrets: base64.StdEncoding.EncodeToString(encryptedData),
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(encryptedSecrets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted secrets: %w", err)
	}

	// Write to the file with restricted permissions
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write encrypted secrets file: %w", err)
	}

	return nil
}

// GenerateEncryptionKey generates a random 32-byte key suitable for AES-256 encryption
func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return key, nil
}

// encrypt encrypts data using AES-GCM
func encrypt(data, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, data, nil), nil
}

// decrypt decrypts data using AES-GCM
func decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
