package sagas

import (
	"fmt"

	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"	
)

type CreateSaga struct {
	db *sqlx.DB
}

type CreateUserRequest struct {
	Users map[string]models.User `json:"users"`
}

func NewCreateSaga(db *sqlx.DB) Saga {
	return &CreateSaga{
		db: db,
	}
}

func (saga *CreateSaga) Run(orchestrator Orchestrator, genericReq GenericRequest) error {
	fmt.Println("----------todo: create assets--------------")
	fmt.Println(genericReq)
	return nil
}

func (saga *CreateSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}