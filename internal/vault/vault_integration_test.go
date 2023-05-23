//go:build integration
// +build integration

package vault

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSQLCredentials(t *testing.T) {
	vaultConfig := VaultConfig{
		Address: strings.Replace(os.Getenv("VAULT_ADDR"), "vault", "localhost", 1),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   os.Getenv("VAULT_MOUNT"),
	}

	// Define the test cases
	testCases := []struct {
		name                string
		expectedCredentials map[string]int
	}{
		{
			name: "Valid credentials",
			expectedCredentials: map[string]int{
				"username": 10,
				"password": 10,
			},
		},
	}

	// Iterate over the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Make the request to retrieve SQL credentials from Vault
			credentials, err := GetSQLCredentials(vaultConfig)
			if err != nil {
				t.Fatalf("Failed to retrieve SQL credentials from Vault: %v", err)
			}

			// Check the retrieved credentials
			assert.Contains(t, credentials, "username", "key username does not exist in the map for test %s", testCase.name)
			assert.Contains(t, credentials, "password", "key password does not exist in the map for test %s", testCase.name)
			assert.Less(t, testCase.expectedCredentials["username"], len(credentials["username"]))
			assert.Less(t, testCase.expectedCredentials["password"], len(credentials["password"]))
		})
	}
}
