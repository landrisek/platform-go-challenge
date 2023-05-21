package sagas

import (
	"encoding/json"
	"log"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/jmoiron/sqlx"
)

type DeleteSaga struct {
	db *sqlx.DB
}

func NewDeleteSaga(db *sqlx.DB) *DeleteSaga {
	return &DeleteSaga{
		db: db,
	}
}

// Delete perform all necessary steps to delete given user assets
func Delete(db *sqlx.DB) Orchestrator {
	orchestrator := SagaOrchestrator{}
	orchestrator.AddSaga(NewDeleteSaga(db))
	return &orchestrator
}

func (saga *DeleteSaga) Run(orchestrator Orchestrator) error {
	var users []models.User
	genericReq := orchestrator.GetRequest()
	err := json.Unmarshal(genericReq.Data, &users)
	if err != nil {
		log.Println("Error on unmarshaling in update saga:", err)
		return err
	}

	var response []models.User

	for _, user := range users {
		respUser := models.User{
			ID: user.ID,
			Name: user.Name,
		}
		// audiences
		for _, audience := range user.Audiences {
			err := models.DeleteAudience(saga.db, audience.ID, user.ID)
			if err != nil {
				log.Println("Error on create audience:", err)
				audience.Error = "Database error on audience"
				respUser.Audiences = append(respUser.Audiences, audience)
			}
		}
		// charts
		for _, chart := range user.Charts {
			err := models.DeleteChart(saga.db, chart.ID, user.ID)
			if err != nil {
				log.Println("Error on create chart:", err)
				chart.Error = "Database error on chart"
				respUser.Charts = append(respUser.Charts, chart)
			}
		}
		// insights
		for _, insight := range user.Insights {
			err := models.DeleteInsight(saga.db, insight.ID, user.ID)
			if err != nil {
				log.Println("Error on update insight:", err)
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
		Format: genericReq.Format,
		Data:   responseData,
	})

	return nil
}

func (saga *DeleteSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}