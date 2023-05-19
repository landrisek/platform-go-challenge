package sagas

import (
	"encoding/json"
	"log"

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

// Create perform all necessary steps to proper create given user assets
func Create(db *sqlx.DB) Orchestrator {
	orchestrator := SagaOrchestrator{}
	orchestrator.AddSaga(NewCreateSaga(db))
	return &orchestrator
}

func (saga *CreateSaga) Run(orchestrator Orchestrator, genericReq GenericRequest) error {
	var users []models.User
	err := json.Unmarshal(genericReq.Data, &users)
	if err != nil {
		log.Println("Error on umarshaling in create saga:", err)
		return err
	}

	var errors []AssetError

	for _, user := range users {
		// charts
		for _, chart := range user.Charts {
			err := models.CreateChart(saga.db, chart, user.ID)
			if err != nil {
				log.Println("Error on create chart:", err)
				errors = append(errors, AssetError{
					Description: chart.Description,
					// we cannot expose actual internal error to external api
					Message: "Database error on chart",
				})
			}
		}
		// insights
		for _, insight := range user.Insights {
			err := models.CreateInsight(saga.db, insight, user.ID)
			if err != nil {
				log.Println("Error on create insight:", err)
				errors = append(errors, AssetError{
					Description: insight.Description,
					Message: "Database error on insight",
				})
			}
		}
		// audiences
		for _, audience := range user.Audiences {
			err := models.CreateAudience(saga.db, audience, user.ID)
			if err != nil {
				log.Println("Error on create audience:", err)
				errors = append(errors, AssetError{
					Description: audience.Description,
					Message:"Database error on audience",
				})
			}
		}
	}

	if len(errors) > 0 {
		response := ErrorResponse{
			Errors: errors,
		}
		responseData, err := json.Marshal(response)
		if err != nil {
			log.Println("Error on marshaling error response:", err)
			return err
		}

		orchestrator.SetResponse(GenericResponse{
			Format: genericReq.Format,
			Data:   responseData,
		})
	}

	return nil
}

func (saga *CreateSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}