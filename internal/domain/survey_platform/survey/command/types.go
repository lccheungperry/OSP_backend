package command

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommandError represents errors that can occur during command handling
type CommandError struct {
	Code    string
	Message string
}

func (e *CommandError) Error() string {
	return e.Message
}

// NewInvalidCommandTypeError creates a new error for invalid command types
func NewInvalidCommandTypeError(cmd interface{}) error {
	return &CommandError{
		Code:    "INVALID_COMMAND_TYPE",
		Message: fmt.Sprintf("invalid command type: %T", cmd),
	}
}

// NewAssignmentError creates a new error for assignment operations
func NewAssignmentError(err error) error {
	return &CommandError{
		Code:    "ASSIGNMENT_ERROR",
		Message: fmt.Sprintf("failed to create question assignment: %v", err),
	}
}

// NewInvalidIDError creates a new error for invalid IDs
func NewInvalidIDError(id string, err error) error {
	return &CommandError{
		Code:    "INVALID_ID",
		Message: fmt.Sprintf("invalid ID format: %s: %v", id, err),
	}
}

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
