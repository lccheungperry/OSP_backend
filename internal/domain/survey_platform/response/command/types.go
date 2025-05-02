package command

import (
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitResponseCommand struct {
	SurveyID primitive.ObjectID `json:"survey_id" validate:"required"`
	Answers  []model.Answer     `json:"answers" validate:"required,min=1"`
}

func (c *SubmitResponseCommand) CommandName() string {
	return "submit_response"
}

type DeleteResponseCommand struct {
	ID primitive.ObjectID `json:"id" validate:"required"`
}

func (c *DeleteResponseCommand) CommandName() string {
	return "delete_response"
}
