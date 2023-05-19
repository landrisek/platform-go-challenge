package sagas

import (
	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"	
)

type ReadSaga struct {
	db *sqlx.DB
}

func NewReadSaga(db *sqlx.DB) *ReadSaga {
	return &ReadSaga{
		db: db,
	}
}

func (saga *ReadSaga) Run(parent Orchestrator) error {
	orchestrator := parent.(*SagaOrchestrator)
	var userID int
	charts, err :=  models.ReadCharts(saga.db, userID)
	if err != nil {
		return err
	}
	generics := chartsToGenerics(charts)
	orchestrator.Response = generics
	return nil
}

func chartsToGenerics(charts []models.Chart) GenericResponse {
	return GenericResponse{}
}

func (saga *ReadSaga) Retrieve(step *Saga) error {
	return nil
}