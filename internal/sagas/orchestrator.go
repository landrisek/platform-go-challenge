package sagas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//GenericRequest is used to specify a generic request to the REST endpoints
type GenericRequest struct {
	Format  RequestResponseFormat `json:"format"`
	Data    json.RawMessage       `json:"data"`
	Context context.Context       `json:"-,omitempty"`
}

//GenericResponse is used to respond to an generic response to the REST endpoints
type GenericResponse struct {
	Format RequestResponseFormat `json:"format"`
	Data   json.RawMessage       `json:"data"`
}

type RequestResponseFormat string

type Orchestrator interface {
	AddSaga(Saga) error
	Run(GenericRequest) (GenericResponse, error)
}

type Saga interface {
	Run(Orchestrator, GenericRequest) error
	Retrieve(Orchestrator) error
}

type SagaOrchestrator struct {
	db       *sqlx.DB
	Response GenericResponse
	steps    []Saga
}

func NewOrchestrator(db *sqlx.DB) Orchestrator {
	return &SagaOrchestrator{
		db: db,
	}
}

func (saga *SagaOrchestrator) AddSaga(step Saga) error {
	saga.steps = append(saga.steps, step)
	return nil
}

type RetrievableError struct {

}

func (e RetrievableError) Error() string {
	return "Retrievable error happened"
}

func (orchestrator *SagaOrchestrator) Run(request GenericRequest) (GenericResponse, error) {
	for i:=0;i<len(orchestrator.steps);i++ {
		step := orchestrator.steps[i]
		if err := step.Run(orchestrator, request); err != nil {
			// saga consider all return error for non-retrievable by default
			if _, ok := err.(RetrievableError); !ok {
				// call all previous sagas to retrieve
				for j:=i;j>=0;j-- {
					if sErr := step.Retrieve(orchestrator); sErr != nil {
						err = fmt.Errorf("after failed: %w, also retrieval failed: %v", err, sErr)					
					}
				}
				return GenericResponse{}, err
			}

		}
	}
	return orchestrator.Response, nil
}

// Create perform all necessary steps to proper create given user assets
func Create(db *sqlx.DB) Orchestrator {
	orchestrator := SagaOrchestrator{}
	orchestrator.AddSaga(NewCreateSaga(db))
	return &orchestrator
}
