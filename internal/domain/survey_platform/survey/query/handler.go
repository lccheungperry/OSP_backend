package query

import (
	"context"
	"fmt"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
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
		return nil, fmt.Errorf("invalid query type: %T", q)
	}

	id, err := primitive.ObjectIDFromHex(getQuery.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid survey ID: %v", err)
	}

	return h.Repository.FindByID(ctx, id)
}

// GetSurveyByTokenHandler handles the GetSurveyByTokenQuery
type GetSurveyByTokenHandler struct {
	Repository repository.SurveyRepository
}

func (h *GetSurveyByTokenHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	getQuery, ok := q.(*GetSurveyByTokenQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type: %T", q)
	}

	return h.Repository.FindByToken(ctx, getQuery.Token)
}

// ListSurveysHandler handles the ListSurveysQuery
type ListSurveysHandler struct {
	Repository repository.SurveyRepository
}

func (h *ListSurveysHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	listQuery, ok := q.(*ListSurveysQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type: %T", q)
	}

	skip := int64(listQuery.Filter.Offset)
	limit := int64(listQuery.Filter.Limit)
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	surveys, total, err := h.Repository.List(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	return struct {
		Surveys []*model.Survey `json:"surveys"`
		Total   int64           `json:"total"`
	}{
		Surveys: surveys,
		Total:   total,
	}, nil
}
