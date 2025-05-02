package query

import (
	"fmt"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
)

type QueryError struct {
	Code    string
	Message string
}

func (e *QueryError) Error() string {
	return e.Message
}

func NewInvalidQueryTypeError(q interface{}) error {
	return &QueryError{
		Code:    "INVALID_QUERY_TYPE",
		Message: fmt.Sprintf("invalid query type: %T", q),
	}
}

func NewInvalidIDError(id string, err error) error {
	return &QueryError{
		Code:    "INVALID_ID",
		Message: fmt.Sprintf("invalid ID format: %s: %v", id, err),
	}
}

type GetSurveyQuery struct {
	ID string `json:"id" validate:"required"`
}

func (q *GetSurveyQuery) QueryName() string {
	return "get_survey"
}

type GetSurveyByTokenQuery struct {
	Token string `json:"token" validate:"required,len=5"`
}

func (q *GetSurveyByTokenQuery) QueryName() string {
	return "get_survey_by_token"
}

type ListSurveysQuery struct {
	Filter *model.SurveyFilter `json:"filter"`
}

func (q *ListSurveysQuery) QueryName() string {
	return "list_surveys"
}
