package sagas

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"	
)

type CreateSaga struct {
	db *sqlx.DB
}

func NewCreateSaga(db *sqlx.DB) Saga {
	return &CreateSaga{
		db: db,
	}
}

func (saga *CreateSaga) Run(orchestrator Orchestrator, genericReq GenericRequest) error {
	genericReq.Data
	fmt.Println("----------todo: create assets--------------")
	return nil
}

func (saga *CreateSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}