package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id" validate:"required"`
	Value      interface{}        `bson:"value" json:"value" validate:"required"`
}

type Response struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SurveyID  primitive.ObjectID `bson:"survey_id" json:"survey_id" validate:"required"`
	Answers   []Answer           `bson:"answers" json:"answers" validate:"required,min=1"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type ResponseFilter struct {
	SurveyID primitive.ObjectID `json:"survey_id"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}
