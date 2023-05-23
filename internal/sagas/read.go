package sagas

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type ReadSaga struct {
	db *sqlx.DB
	// HINT: This is for future implementation of cache
	client *redis.Client
}

func NewReadSaga(db *sqlx.DB, client *redis.Client) *ReadSaga {
	return &ReadSaga{
		db:     db,
		client: client,
	}
}

// Read perform all necessary steps to read given user assets
func Read(db *sqlx.DB, client *redis.Client) Orchestrator {
	orchestrator := SagaOrchestrator{}
	orchestrator.AddSaga(NewReadSaga(db, client))
	return &orchestrator
}

func (saga *ReadSaga) Run(orchestrator Orchestrator) error {
	genericReq := orchestrator.GetRequest()

	var users []models.User

	err := json.Unmarshal(genericReq.Data, &users)
	if err != nil {
		log.Println("Error on unmarshaling in read saga:", err)
		return err
	}

	if len(users) > 1 {
		// HINT: more users will be easy to implement, we are just meeting acceptance criteria
		return fmt.Errorf("Update assets of more users is not supported")
	}

	user := users[0]
	audiences, err := models.ReadAudiences(saga.db, user.ID)
	if err != nil {
		return err
	}
	user.Audiences = audiences

	charts, err := models.ReadCharts(saga.db, user.ID)
	if err != nil {
		return err
	}
	user.Charts = charts

	insights, err := models.ReadInsights(saga.db, user.ID)
	if err != nil {
		return err
	}
	user.Insights = insights

	responseData, err := json.Marshal([]models.User{user})
	if err != nil {
		log.Println("Error on marshaling error response:", err)
		return err
	}

	orchestrator.SetResponse(GenericResponse{
		Format: "json",
		Data:   responseData,
	})

	return nil
}

func (saga *ReadSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}
