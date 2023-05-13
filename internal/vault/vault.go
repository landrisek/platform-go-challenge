package vault

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// VaultConfig represents the configuration for accessing Vault
type VaultConfig struct {
	Address string
	Token   string
}

// GetSQLCredentials retrieves SQL credentials from Vault
func GetSQLCredentials(vaultConfig VaultConfig, mountPath string) (map[string]string, error) {
	// Create a new Vault client
	client, err := api.NewClient(&api.Config{
		Address: vaultConfig.Address,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Vault client")
	}

	// Set the Vault token
	client.SetToken(vaultConfig.Token)

	// Read the SQL credentials from Vault
	secretPath := fmt.Sprintf("%s/creds/sudo", mountPath)
	secretData, err := client.Logical().Read(secretPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read SQL credentials from Vault: %s", secretPath)
	}
	if secretData == nil {
		return nil, errors.Errorf("no secret data found at path: %s", secretPath)
	}

	// Extract the SQL credentials from the secret data
	credentials := make(map[string]string)
	for key, value := range secretData.Data {
		if stringValue, ok := value.(string); ok {
			credentials[key] = stringValue
		}
	}

	return credentials, nil
}

// Example usage
func main() {
	vaultConfig := VaultConfig{
		Address: "http://vault.example.com:8200",
		Token:   "your_vault_token",
	}
	mountPath := "secret/database"

	credentials, err := GetSQLCredentials(vaultConfig, mountPath)
	if err != nil {
		log.Fatalf("Failed to retrieve SQL credentials from Vault: %v", err)
	}

	// Use the retrieved credentials as needed
	fmt.Println("Username:", credentials["username"])
	fmt.Println("Password:", credentials["password"])
}
