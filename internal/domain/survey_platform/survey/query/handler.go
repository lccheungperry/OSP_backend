package query

import (
	"context"

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
	getCmd, ok := q.(*GetSurveyQuery)
	if !ok {
		return nil, NewInvalidQueryTypeError(q)
	}

	id, err := primitive.ObjectIDFromHex(getCmd.ID)
	if err != nil {
		return nil, NewInvalidIDError(getCmd.ID, err)
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
	getCmd, ok := q.(*GetSurveyByTokenQuery)
	if !ok {
		return nil, NewInvalidQueryTypeError(q)
	}

	survey, err := h.Repository.FindByToken(ctx, getCmd.Token)
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
	listCmd, ok := q.(*ListSurveysQuery)
	if !ok {
		return nil, NewInvalidQueryTypeError(q)
	}

	skip := int64(listCmd.Filter.Offset)
	limit := int64(listCmd.Filter.Limit)
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
