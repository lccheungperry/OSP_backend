package service

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

func NewSurveyNotFoundError(err error) error {
	return &ServiceError{
		Code:    "SURVEY_NOT_FOUND",
		Message: fmt.Sprintf("survey not found: %v", err),
	}
}

func NewQuestionNotFoundError(err error) error {
	return &ServiceError{
		Code:    "QUESTION_NOT_FOUND",
		Message: fmt.Sprintf("question not found: %v", err),
	}
}

func NewQuestionNotInSurveyError(questionID, surveyID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "QUESTION_NOT_IN_SURVEY",
		Message: fmt.Sprintf("question %s not found in survey %s", questionID.Hex(), surveyID.Hex()),
	}
}

func NewInvalidSpecificationsError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_SPECIFICATIONS",
		Message: fmt.Sprintf("invalid specifications format for question %s", questionID.Hex()),
	}
}

func NewInvalidOptionsError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_OPTIONS",
		Message: fmt.Sprintf("invalid options format for question %s", questionID.Hex()),
	}
}

func NewInvalidAnswerTypeError(questionID primitive.ObjectID, questionType string) error {
	return &ServiceError{
		Code:    "INVALID_ANSWER_TYPE",
		Message: fmt.Sprintf("invalid answer value type for %s question %s", questionType, questionID.Hex()),
	}
}

func NewInvalidOptionValueError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_OPTION_VALUE",
		Message: fmt.Sprintf("invalid option value for question %s", questionID.Hex()),
	}
}

func NewInvalidScaleValueError(questionID primitive.ObjectID) error {
	return &ServiceError{
		Code:    "INVALID_SCALE_VALUE",
		Message: fmt.Sprintf("invalid scale value for question %s", questionID.Hex()),
	}
}

func NewResponseCreationError(err error) error {
	return &ServiceError{
		Code:    "RESPONSE_CREATION_ERROR",
		Message: fmt.Sprintf("failed to create response: %v", err),
	}
}
