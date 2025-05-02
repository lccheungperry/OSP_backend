package repository

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseRepository interface {
	Create(ctx context.Context, response *model.Response) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Response, error)
	FindBySurveyID(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
