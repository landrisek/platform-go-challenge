package sagas

import (
	"github.com/landrisek/platform-go-challenge/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"	
)

type ReadSaga struct {
	db *sqlx.DB
}

func NewReadSaga(db *sqlx.DB) *Saga {
	return &ReadSaga{
		db: db
	}
}

func (saga *ReadSaga) Run(parent *Orchestrator) error {
	orchestrator := parent.(*SagaOrchestrator)

	charts, err :=  models.GetChartsByUserID(saga.db, orchestrator.userID)
	if err != nil {
		return err
	}
	generics := chartsToGenerics(charts)
	orchestrator.assets = append(orchestrator.assets, generics)
	return nil
}

func chartsToGenerics(charts []models.Chart) []GenericResponse {
	return nil
}

func (saga *ReadSaga) Retrieve(step *Saga) error {
	return nil
}