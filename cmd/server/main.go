package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
		makeRequest()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type VaultResponse struct {
	Data  VaultData `json:"data"`
	Lease int       `json:"lease_duration"`
}

type VaultData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//curl --header "X-Vault-Token: myroot" --request GET  http://localhost:8200/v1/database/creds/sudo
func makeRequest() {
	fmt.Println("-------makeRequest()-------")
	req, err := http.NewRequest(http.MethodGet, "http://vault:8200/v1/database/creds/sudo", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Vault-Token", "myroot")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var vaultResponse VaultResponse
	err = json.NewDecoder(response.Body).Decode(&vaultResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Username:", vaultResponse.Data.Username)
	fmt.Println("Password:", vaultResponse.Data.Password)
	fmt.Println("Lease Duration:", vaultResponse.Lease)
}