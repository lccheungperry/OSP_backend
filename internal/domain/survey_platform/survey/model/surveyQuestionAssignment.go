package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SurveyQuestionAssignment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SurveyID   primitive.ObjectID `bson:"survey_id" json:"survey_id" validate:"required"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id" validate:"required"`
	Order      int                `bson:"order" json:"order" validate:"required,min=1"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
