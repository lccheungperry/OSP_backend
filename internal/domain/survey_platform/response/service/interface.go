package service

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseServiceInterface interface {
	HandleCommand(ctx context.Context, cmd command.Command) error
	HandleQuery(ctx context.Context, q query.Query) (interface{}, error)
	GetResponses(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error)
}
