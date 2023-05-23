package vault

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// VaultConfig represents the configuration for accessing Vault
type VaultConfig struct {
	Address string
	Token   string
	Mount   string
}

type VaultResponse struct {
	Data  VaultData `json:"data"`
	Lease int       `json:"lease_duration"`
}

type VaultData struct {
	Username string
	Password string
}

// HINT: Redundant, left here for broader discussion on in-build cli vs. in-build request
// which are more or less same
func GetSQLCredentialsWithCLI(vaultConfig VaultConfig) (map[string]string, error) {
	// Create a new Vault client
	client, err := api.NewClient(&api.Config{
		Address: vaultConfig.Address,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Vault client")
	}

	// Set the Vault token
	client.SetToken(vaultConfig.Token)

	secretPath := fmt.Sprintf("%s/creds/sudo", vaultConfig.Mount)
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

// GetSQLCredentials retrieves SQL credentials from Vault
// HINT: curl --header "X-Vault-Token: myroot" --request GET  http://vault:8200/v1/mysql_sandbox/creds/sudo
func GetSQLCredentials(vaultConfig VaultConfig) (map[string]string, error) {
	req, err := http.NewRequest(http.MethodGet, vaultConfig.Address+"/v1/"+vaultConfig.Mount+"/creds/sudo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", vaultConfig.Token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var vaultResponse VaultResponse
	err = json.NewDecoder(response.Body).Decode(&vaultResponse)
	if err != nil {
		return nil, err
	}

	// HINT: this will be removed, left here for discussion
	fmt.Println("Username:", vaultResponse.Data.Username)
	fmt.Println("Password:", vaultResponse.Data.Password)
	fmt.Println("Lease Duration:", vaultResponse.Lease)
	return map[string]string{
		"username": vaultResponse.Data.Username,
		"password": vaultResponse.Data.Password,
	}, nil
}
