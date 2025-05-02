package command

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSurveyCommand struct {
	Title string `json:"title" validate:"required,min=3,max=100"`
}

func (c *CreateSurveyCommand) CommandName() string {
	return "create_survey"
}

type UpdateSurveyCommand struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	Title     string             `json:"title" validate:"required,min=3,max=100"`
	Questions []struct {
		QuestionID primitive.ObjectID `json:"questionId" validate:"required"`
		Order      int                `json:"order" validate:"required,min=1"`
	} `json:"questions"`
}

func (c *UpdateSurveyCommand) CommandName() string {
	return "update_survey"
}

type DeleteSurveyCommand struct {
	ID primitive.ObjectID `json:"id" validate:"required"`
}

func (c *DeleteSurveyCommand) CommandName() string {
	return "delete_survey"
}

type CreateAssignmentCommand struct {
	SurveyID   primitive.ObjectID `json:"survey_id" validate:"required"`
	QuestionID primitive.ObjectID `json:"question_id" validate:"required"`
	Order      int                `json:"order" validate:"required,min=1"`
}

func (c *CreateAssignmentCommand) CommandName() string {
	return "create_assignment"
}

type UpdateAssignmentCommand struct {
	ID    primitive.ObjectID `json:"id" validate:"required"`
	Order int                `json:"order" validate:"required,min=1"`
}

func (c *UpdateAssignmentCommand) CommandName() string {
	return "update_assignment"
}

type DeleteAssignmentCommand struct {
	ID primitive.ObjectID `json:"id" validate:"required"`
}

func (c *DeleteAssignmentCommand) CommandName() string {
	return "delete_assignment"
}
