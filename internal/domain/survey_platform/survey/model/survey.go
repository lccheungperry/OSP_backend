package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Survey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" validate:"required,min=3,max=100"`
	Token     string             `bson:"token" json:"token" validate:"required,len=5"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
