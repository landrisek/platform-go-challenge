package sagas

import (
	"encoding/json"
	"log"

	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/jmoiron/sqlx"
)

type UpdateSaga struct {
	db *sqlx.DB
}

func NewUpdateSaga(db *sqlx.DB) *UpdateSaga {
	return &UpdateSaga{
		db: db,
	}
}

// Update perform all necessary steps to update given user assets
func Update(db *sqlx.DB) Orchestrator {
	orchestrator := SagaOrchestrator{}
	orchestrator.AddSaga(NewUpdateSaga(db))
	return &orchestrator
}

func (saga *UpdateSaga) Run(orchestrator Orchestrator) error {
	var users []models.UserSafeUpdate
	genericReq := orchestrator.GetRequest()
	err := json.Unmarshal(genericReq.Data, &users)
	if err != nil {
		log.Println("Error on unmarshaling in update saga:", err)
		return err
	}

	var response []models.UserSafeUpdate

	for _, user := range users {
		respUser := models.UserSafeUpdate{
			ID:   user.ID,
			Name: user.Name,
		}
		// audiences
		for _, audience := range user.Audiences {
			err := models.UpdateAudience(saga.db, audience, user.ID)
			if err != nil {
				log.Println("Error on update audience:", err)
				audience.Error = "Database error on audience"
				respUser.Audiences = append(respUser.Audiences, audience)
			}
		}
		// charts
		for _, chart := range user.Charts {
			err := models.UpdateChart(saga.db, chart, user.ID)
			if err != nil {
				log.Println("Error on update chart:", err)
				chart.Error = "Database error on chart"
				respUser.Charts = append(respUser.Charts, chart)
			}
		}
		// insights
		for _, insight := range user.Insights {
			err := models.UpdateInsight(saga.db, insight, user.ID)
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
		Format: "json",
		Data:   responseData,
	})

	return nil
}

func (saga *UpdateSaga) Retrieve(orchestrator Orchestrator) error {
	return nil
}
