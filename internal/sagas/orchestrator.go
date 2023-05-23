package sagas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// GenericRequest is used to specify a generic request to the REST endpoints
type GenericRequest struct {
	Format  RequestResponseFormat `json:"format"`
	Data    json.RawMessage       `json:"data"`
	Context context.Context       `json:"-,omitempty"`
}

// GenericResponse is used to respond to an generic response to the REST endpoints
type GenericResponse struct {
	Format RequestResponseFormat `json:"format"`
	Data   json.RawMessage       `json:"data"`
}

type RequestResponseFormat string

type Orchestrator interface {
	AddSaga(Saga) error
	Run(GenericRequest) (GenericResponse, error)
	GetResponse() GenericResponse
	SetResponse(GenericResponse)
	GetRequest() GenericRequest
	SetRequest(GenericRequest)
}

type Saga interface {
	Run(Orchestrator) error
	Retrieve(Orchestrator) error
}

type SagaOrchestrator struct {
	db       *sqlx.DB
	response GenericResponse
	request  GenericRequest
	steps    []Saga
}

func NewOrchestrator(db *sqlx.DB) Orchestrator {
	return &SagaOrchestrator{
		db:       db,
		response: GenericResponse{},
	}
}

func (orchestrator *SagaOrchestrator) AddSaga(step Saga) error {
	orchestrator.steps = append(orchestrator.steps, step)
	return nil
}

func (orchestrator *SagaOrchestrator) Run(genericReq GenericRequest) (GenericResponse, error) {
	orchestrator.request = genericReq
	for i := 0; i < len(orchestrator.steps); i++ {
		step := orchestrator.steps[i]
		if err := step.Run(orchestrator); err != nil {
			// saga consider all return error for non-retrievable by default
			if _, ok := err.(RetrievableError); !ok {
				// call all previous sagas to retrieve
				for j := i; j >= 0; j-- {
					if sErr := step.Retrieve(orchestrator); sErr != nil {
						err = fmt.Errorf("after failed: %w, also retrieval failed: %v", err, sErr)
					}
				}
				return orchestrator.response, err
			}
		}
	}
	return orchestrator.response, nil
}

func (orchestrator *SagaOrchestrator) GetResponse() GenericResponse {
	return orchestrator.response
}

func (orchestrator *SagaOrchestrator) SetResponse(response GenericResponse) {
	orchestrator.response = response
}

func (orchestrator *SagaOrchestrator) SetRequest(request GenericRequest) {
	orchestrator.request = request
}

func (orchestrator *SagaOrchestrator) GetRequest() GenericRequest {
	return orchestrator.request
}
