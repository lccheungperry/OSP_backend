package query

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetResponsesQuery struct {
	SurveyID primitive.ObjectID `json:"survey_id" validate:"required"`
	Skip     int64              `json:"skip" validate:"min=0"`
	Limit    int64              `json:"limit" validate:"min=1,max=100"`
}

func (q *GetResponsesQuery) QueryName() string {
	return "get_responses"
}
