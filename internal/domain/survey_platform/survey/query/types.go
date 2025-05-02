package query

import (
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
)

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
