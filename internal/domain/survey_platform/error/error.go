package error

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

var (
	ErrSurveyNotFound   = fmt.Errorf("survey not found")
	ErrQuestionNotFound = fmt.Errorf("question not found")
)

func NewInvalidCommandError() error {
	return &ServiceError{
		Code:    "INVALID_COMMAND",
		Message: "invalid command",
	}
}

func NewInvalidCommandTypeError(cmd interface{}) error {
	return &ServiceError{
		Code:    "INVALID_COMMAND_TYPE",
		Message: fmt.Sprintf("invalid command type: %T", cmd),
	}
}

func NewInvalidQuestionIDError(err error) error {
	return &ServiceError{
		Code:    "INVALID_QUESTION_ID",
		Message: fmt.Sprintf("invalid question ID: %v", err),
	}
}

func NewAssignmentError(err error) error {
	return &ServiceError{
		Code:    "ASSIGNMENT_ERROR",
		Message: fmt.Sprintf("failed to create question assignment: %v", err),
	}
}

func NewInvalidIDError(id string, err error) error {
	return &ServiceError{
		Code:    "INVALID_ID",
		Message: fmt.Sprintf("invalid ID format: %s: %v", id, err),
	}
}

func NewInvalidQueryTypeError(q interface{}) error {
	return &ServiceError{
		Code:    "INVALID_QUERY_TYPE",
		Message: fmt.Sprintf("invalid query type: %T", q),
	}
}

func NewCreateError(err error) error {
	return &ServiceError{
		Code:    "CREATE_ERROR",
		Message: fmt.Sprintf("failed to create: %v", err),
	}
}

func NewUpdateError(err error) error {
	return &ServiceError{
		Code:    "UPDATE_ERROR",
		Message: fmt.Sprintf("failed to update: %v", err),
	}
}

func NewDeleteError(err error) error {
	return &ServiceError{
		Code:    "DELETE_ERROR",
		Message: fmt.Sprintf("failed to delete: %v", err),
	}
}

func NewInvalidSpecificationsError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_SPECIFICATIONS",
		Message: fmt.Sprintf("invalid specifications for question %s", questionID.Hex()),
	}
}

func NewInvalidOptionsError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_OPTIONS",
		Message: fmt.Sprintf("invalid options for question %s", questionID.Hex()),
	}
}

func NewInvalidAnswerTypeError(questionID primitive.ObjectID, expectedType string) error {
	return &ServiceError{
		Code:    "INVALID_ANSWER_TYPE",
		Message: fmt.Sprintf("invalid answer type for question %s, expected %s", questionID.Hex(), expectedType),
	}
}

func NewInvalidOptionValueError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_OPTION_VALUE",
		Message: fmt.Sprintf("invalid option value for question %s", questionID.Hex()),
	}
}

func NewInvalidScaleValueError(questionID primitive.ObjectID, scale float64, maxScale interface{}) error {
	var maxScaleStr string
	switch v := maxScale.(type) {
	case float64:
		maxScaleStr = fmt.Sprintf("%v", v)
	case int:
		maxScaleStr = fmt.Sprintf("%d", v)
	case int32:
		maxScaleStr = fmt.Sprintf("%d", v)
	case int64:
		maxScaleStr = fmt.Sprintf("%d", v)
	default:
		maxScaleStr = fmt.Sprintf("%v", v)
	}
	return &ServiceError{
		Code:    "INVALID_SCALE_VALUE",
		Message: fmt.Sprintf("invalid scale value: %v (max: %s) for question %s", scale, maxScaleStr, questionID.Hex()),
	}
}

func NewQuestionNotInSurveyError(questionID, surveyID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "QUESTION_NOT_IN_SURVEY",
		Message: fmt.Sprintf("question %s is not assigned to survey %s", questionID.Hex(), surveyID.Hex()),
	}
}

func NewResponseCreationError(err error) error {
	return &ServiceError{
		Code:    "RESPONSE_CREATION_ERROR",
		Message: fmt.Sprintf("failed to create response: %v", err),
	}
}

func NewResponseRetrievalError(err error) error {
	return &ServiceError{
		Code:    "RESPONSE_RETRIEVAL_ERROR",
		Message: fmt.Sprintf("failed to retrieve responses: %v", err),
	}
}
