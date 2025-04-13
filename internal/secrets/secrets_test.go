package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_BasicOperations(t *testing.T) {
	// Create a new manager
	manager := NewManager()

	// Test Set and Get
	manager.Set("key1", "value1")
	value, exists := manager.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist, but it doesn't")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}

	// Test GetOrDefault with existing key
	value = manager.GetOrDefault("key1", "default")
	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}

	// Test GetOrDefault with non-existing key
	value = manager.GetOrDefault("nonexistent", "default")
	if value != "default" {
		t.Errorf("Expected default, got %s", value)
	}

	// Test GetAll
	manager.Set("key2", "value2")
	allSecrets := manager.GetAll()
	if len(allSecrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(allSecrets))
	}
	if allSecrets["key1"] != "value1" || allSecrets["key2"] != "value2" {
		t.Error("GetAll returned incorrect values")
	}
}

func TestManager_LoadFromEnv(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TEST_SECRET_KEY1", "value1")
	os.Setenv("TEST_SECRET_KEY2", "value2")
	os.Setenv("OTHER_PREFIX_KEY", "other_value")
	defer func() {
		os.Unsetenv("TEST_SECRET_KEY1")
		os.Unsetenv("TEST_SECRET_KEY2")
		os.Unsetenv("OTHER_PREFIX_KEY")
	}()

	// Create a new manager and load from env
	manager := NewManager()
	err := manager.LoadFromEnv("TEST_SECRET_")
	if err != nil {
		t.Errorf("LoadFromEnv failed: %v", err)
	}

	// Verify the secrets were loaded correctly
	value, exists := manager.Get("KEY1")
	if !exists || value != "value1" {
		t.Errorf("Expected KEY1=value1, got exists=%v, value=%s", exists, value)
	}

	value, exists = manager.Get("KEY2")
	if !exists || value != "value2" {
		t.Errorf("Expected KEY2=value2, got exists=%v, value=%s", exists, value)
	}

	// Verify that keys with other prefixes were not loaded
	_, exists = manager.Get("OTHER_PREFIX_KEY")
	if exists {
		t.Error("Key with different prefix should not have been loaded")
	}
}

func TestManager_FileOperations(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "secrets_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file path
	testFilePath := filepath.Join(tempDir, "test_secrets.json")

	// Create a manager with some secrets
	manager := NewManager()
	manager.Set("key1", "value1")
	manager.Set("key2", "value2")

	// Save to file
	err = manager.SaveToFile(testFilePath)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Create a new manager and load from the file
	newManager := NewManager()
	err = newManager.LoadFromFile(testFilePath)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// Verify the secrets were loaded correctly
	value, exists := newManager.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected key1=value1, got exists=%v, value=%s", exists, value)
	}

	value, exists = newManager.Get("key2")
	if !exists || value != "value2" {
		t.Errorf("Expected key2=value2, got exists=%v, value=%s", exists, value)
	}
}

func TestManager_EncryptedFileOperations(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "secrets_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file path
	testFilePath := filepath.Join(tempDir, "test_secrets_encrypted.json")

	// Generate an encryption key
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("Failed to generate encryption key: %v", err)
	}

	// Create a manager with some secrets
	manager := NewManager()
	manager.Set("key1", "value1")
	manager.Set("key2", "value2")

	// Save to encrypted file
	err = manager.SaveEncryptedToFile(testFilePath, key)
	if err != nil {
		t.Fatalf("SaveEncryptedToFile failed: %v", err)
	}

	// Create a new manager and load from the encrypted file
	newManager := NewManager()
	err = newManager.LoadEncryptedFromFile(testFilePath, key)
	if err != nil {
		t.Fatalf("LoadEncryptedFromFile failed: %v", err)
	}

	// Verify the secrets were loaded correctly
	value, exists := newManager.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected key1=value1, got exists=%v, value=%s", exists, value)
	}

	value, exists = newManager.Get("key2")
	if !exists || value != "value2" {
		t.Errorf("Expected key2=value2, got exists=%v, value=%s", exists, value)
	}
}

func TestKeyDerivation(t *testing.T) {
	// Test key derivation with defaults
	password := "test-password"
	key, salt, err := DeriveKeyWithDefaults(password)
	if err != nil {
		t.Fatalf("DeriveKeyWithDefaults failed: %v", err)
	}
	if len(key) != DefaultKeyLength {
		t.Errorf("Expected key length %d, got %d", DefaultKeyLength, len(key))
	}
	if len(salt) != DefaultSaltLength {
		t.Errorf("Expected salt length %d, got %d", DefaultSaltLength, len(salt))
	}

	// Test key derivation with the same password and salt
	key2 := DeriveKeyWithSalt(password, salt)
	if len(key2) != DefaultKeyLength {
		t.Errorf("Expected key length %d, got %d", DefaultKeyLength, len(key2))
	}

	// Keys derived from the same password and salt should be identical
	for i := 0; i < len(key); i++ {
		if key[i] != key2[i] {
			t.Errorf("Keys don't match at position %d: %d != %d", i, key[i], key2[i])
		}
	}

	// Test encoding and decoding
	encodedKey := EncodeKeyToString(key)
	decodedKey, err := DecodeKeyFromString(encodedKey)
	if err != nil {
		t.Fatalf("DecodeKeyFromString failed: %v", err)
	}

	// Original and decoded keys should be identical
	for i := 0; i < len(key); i++ {
		if key[i] != decodedKey[i] {
			t.Errorf("Original and decoded keys don't match at position %d: %d != %d", i, key[i], decodedKey[i])
		}
	}
}
