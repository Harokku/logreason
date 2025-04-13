# LogReason Secrets Management

This package provides a secure way to manage API keys and other sensitive information in the LogReason application.

## Features

- Load secrets from environment variables
- Load secrets from plain JSON files (for development)
- Load secrets from encrypted JSON files (for production)
- Thread-safe access to secrets
- Password-based encryption using industry-standard algorithms
- Flexible API for accessing secrets throughout the application

## Security Considerations

This package implements several security best practices:

1. **Environment Variables**: Preferred method for production as they're not stored on disk
2. **Encryption**: AES-256-GCM for file-based secrets, providing confidentiality and integrity
3. **Key Derivation**: PBKDF2 with SHA-256 for deriving encryption keys from passwords
4. **File Permissions**: Restricted file permissions (0600) for secret files
5. **Thread Safety**: Mutex-protected access to prevent race conditions
6. **Memory Safety**: No logging of actual secret values

## Usage Examples

### Basic Usage

```go
// Create a new secrets manager
secretsManager := secrets.NewManager()

// Load secrets from environment variables
secretsManager.LoadFromEnv("API_KEY_")

// Access a secret
apiKey, exists := secretsManager.Get("SERVICE1")
if exists {
    // Use the API key
}

// Access with default value
apiKey := secretsManager.GetOrDefault("SERVICE2", "default-key")
```

### Loading from Files

```go
// Load from a plain JSON file (development only)
secretsManager.LoadFromFile("config/secrets.json")

// Load from an encrypted file
password := "your-secure-password"
salt := retrieveSaltFromSecureStorage()
key := secrets.DeriveKeyWithSalt(password, salt)
secretsManager.LoadEncryptedFromFile("config/secrets.enc.json", key)
```

### Saving Secrets

```go
// Save to a plain JSON file (development only)
secretsManager.SaveToFile("config/secrets.json")

// Save to an encrypted file
password := "your-secure-password"
key, salt, _ := secrets.DeriveKeyWithDefaults(password)
// Store the salt securely for later decryption
storeSaltSecurely(salt)
secretsManager.SaveEncryptedToFile("config/secrets.enc.json", key)
```

## Best Practices

1. **Environment Variables**: Use environment variables in production whenever possible
2. **Encryption Keys**: Never hardcode encryption keys or passwords
3. **Salt Storage**: Store salt values separately from encrypted data
4. **Key Rotation**: Implement a key rotation policy for long-term deployments
5. **Access Control**: Limit which parts of your application can access secrets
6. **Monitoring**: Log access attempts to sensitive operations (but not the secret values themselves)

## Implementation Details

The package consists of three main components:

1. **Manager**: Core functionality for storing and retrieving secrets
2. **Encryption**: AES-GCM encryption/decryption for file-based secrets
3. **Key Derivation**: PBKDF2 for deriving encryption keys from passwords

## Future Enhancements

- Integration with cloud key management services (AWS KMS, Google Cloud KMS, etc.)
- Hardware security module (HSM) support
- Secret rotation and versioning
- Access auditing and logging