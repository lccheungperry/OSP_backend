package repository

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SurveyRepository interface {
	Create(ctx context.Context, survey *model.Survey) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Survey, error)
	FindByToken(ctx context.Context, token string) (*model.Survey, error)
	List(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error)
	Update(ctx context.Context, survey *model.Survey) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
