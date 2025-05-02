package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionFormat string

const (
	FormatTextbox        QuestionFormat = "textbox"
	FormatMultipleChoice QuestionFormat = "multiple_choice"
	FormatLikert         QuestionFormat = "likert"
)

type Question struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title" validate:"required,min=3,max=500"`
	Format         QuestionFormat     `bson:"format" json:"format" validate:"required,oneof=textbox multiple_choice likert"`
	Specifications interface{}        `bson:"specifications" json:"specifications" validate:"required"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type Option struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Label string             `json:"label" validate:"required"`
}

type MultipleChoiceSpecification struct {
	Options []Option `json:"options" validate:"required,min=2"`
}

type LikertOption struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Label string             `json:"label" validate:"required"`
	Scale int                `json:"scale" validate:"required"`
}

type LikertSpecification struct {
	Options []LikertOption `json:"options" validate:"required,min=2"`
}

type QuestionFilter struct {
	Format QuestionFormat `json:"format"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
