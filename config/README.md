# Configuration Directory

This directory contains configuration files for the LogReason application.

## Secret Management

### secret.json

The `secret.json` file is an example configuration file that demonstrates how to store API keys and other sensitive information for development purposes. In a production environment, it's recommended to use environment variables or encrypted files instead.

**Format:**
```json
{
  "KEY_NAME1": "secret_value1",
  "KEY_NAME2": "secret_value2"
}
```

**Example Usage in Code:**
```go
secretsManager := secrets.NewManager()
if err := secretsManager.LoadFromFile("config/secret.json"); err != nil {
    log.Printf("Warning: Failed to load secrets from file: %v", err)
}

// Access a secret
apiKey, exists := secretsManager.Get("API_KEY_SERVICE1")
```

### Security Considerations

1. **Never commit real secrets to version control**
   - The example `secret.json` file contains dummy values only
   - Add `config/secret.json` to your `.gitignore` file (already done in this project)

2. **For production environments:**
   - Use environment variables when possible
   - Or use encrypted files with `secrets.enc.json`
   - Store encryption keys and salts separately from your application

3. **File permissions:**
   - The secrets package automatically sets restrictive permissions (0600) on secret files

## Creating an Encrypted Secrets File

To create an encrypted version of your secrets:

```go
// Generate a key from a password
password := "your-secure-password"
key, salt, _ := secrets.DeriveKeyWithDefaults(password)

// Store the salt securely for later decryption
fmt.Printf("Generated salt (save this): %s\n", secrets.EncodeKeyToString(salt))

// Save the secrets to an encrypted file
secretsManager.SaveEncryptedToFile("config/secrets.enc.json", key)
```

For more details, see the documentation in the `internal/secrets` package.