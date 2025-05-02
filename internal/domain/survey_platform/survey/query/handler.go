package query

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	platform_error "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/error"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DEFAULT_LIMIT = 10

type GetSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *GetSurveyHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	getQuery, ok := q.(*GetSurveyQuery)
	if !ok {
		return nil, platform_error.NewInvalidQueryTypeError(q)
	}

	id, err := primitive.ObjectIDFromHex(getQuery.ID)
	if err != nil {
		return nil, platform_error.NewInvalidIDError(getQuery.ID, err)
	}

	survey, err := h.Repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

// GetSurveyByTokenHandler handles the GetSurveyByTokenQuery
type GetSurveyByTokenHandler struct {
	Repository repository.SurveyRepository
}

func (h *GetSurveyByTokenHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	getQuery, ok := q.(*GetSurveyByTokenQuery)
	if !ok {
		return nil, platform_error.NewInvalidQueryTypeError(q)
	}

	survey, err := h.Repository.FindByToken(ctx, getQuery.Token)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

// ListSurveysHandler handles the ListSurveysQuery
type ListSurveysHandler struct {
	Repository repository.SurveyRepository
}

func (h *ListSurveysHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	listQuery, ok := q.(*ListSurveysQuery)
	if !ok {
		return nil, platform_error.NewInvalidQueryTypeError(q)
	}

	var skip, limit int64
	if listQuery.Filter != nil {
		skip = int64(listQuery.Filter.Offset)
		limit = int64(listQuery.Filter.Limit)
	}

	surveys, total, err := h.Repository.List(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	return struct {
		Surveys []*model.Survey
		Total   int64
	}{
		Surveys: surveys,
		Total:   total,
	}, nil
}
