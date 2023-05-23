package sagas

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BlacklistSaga struct {
	blacklistAddr string
}

func NewBlacklistSaga(blacklistAddr string) Saga {
	return &BlacklistSaga{
		blacklistAddr: blacklistAddr,
	}
}

func (saga *BlacklistSaga) Run(orchestrator Orchestrator) error {
	genericReq := orchestrator.GetRequest()
	// Send the HTTP request
	resp, err := http.Post(saga.blacklistAddr, "application/json", bytes.NewBuffer(genericReq.Data))
	if err != nil {
		return fmt.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %v", err)
	}

	// Process the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error response: %s", resp.Status)
	}

	genericReq.Data = respBody
	// blacklisted request is set for following processing
	orchestrator.SetRequest(genericReq)

	return nil
}

func (saga *BlacklistSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}
