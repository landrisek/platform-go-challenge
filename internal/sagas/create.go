package sagas

import (
	"encoding/json"
	"log"

	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/jmoiron/sqlx"
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
func Create(db *sqlx.DB, blacklistAddr string) Orchestrator {
	orchestrator := SagaOrchestrator{}
	orchestrator.AddSaga(NewBlacklistSaga(blacklistAddr))
	orchestrator.AddSaga(NewCreateSaga(db))
	return &orchestrator
}

func (saga *CreateSaga) Run(orchestrator Orchestrator) error {
	var users []models.User
	genericReq := orchestrator.GetRequest()
	err := json.Unmarshal(genericReq.Data, &users)
	if err != nil {
		log.Println("Error on unmarshaling in create saga:", err)
		return err
	}

	var response []models.User

	for _, user := range users {
		respUser := models.User{
			ID:   user.ID,
			Name: user.Name,
		}
		// audiences
		for _, audience := range user.Audiences {
			err := models.CreateAudience(saga.db, audience, user.ID)
			if err != nil {
				log.Println("Error on create audience:", err)
				audience.Error = "Database error on audience"
				respUser.Audiences = append(respUser.Audiences, audience)
			}
		}
		// charts
		for _, chart := range user.Charts {
			err := models.CreateChart(saga.db, chart, user.ID)
			if err != nil {
				log.Println("Error on create chart:", err)
				chart.Error = "Database error on chart"
				respUser.Charts = append(respUser.Charts, chart)
			}
		}
		// insights
		for _, insight := range user.Insights {
			err := models.CreateInsight(saga.db, insight, user.ID)
			if err != nil {
				log.Println("Error on create insight:", err)
				insight.Error = "Database error on insight"
				respUser.Insights = append(respUser.Insights, insight)
			}
		}
		if len(respUser.Audiences) > 0 || len(respUser.Charts) > 0 || len(respUser.Insights) > 0 {
			response = append(response, respUser)
		}
	}

	responseData, err := json.Marshal(response)
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

func (saga *CreateSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}
