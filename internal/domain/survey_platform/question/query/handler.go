package query

import (
	"context"
	"fmt"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
)

var ErrInvalidQuery = fmt.Errorf("invalid query")

type QueryHandler interface {
	HandleQuery(ctx context.Context, query interface{}) (interface{}, error)
}

type GetQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *GetQuestionHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	getQuery, ok := q.(*GetQuestionQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type: %T", q)
	}

	return h.Repository.GetByID(ctx, getQuery.ID)
}

type ListQuestionsHandler struct {
	Repository repository.QuestionRepository
}

func (h *ListQuestionsHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	listQuery, ok := q.(*ListQuestionsQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type: %T", q)
	}

	filter := &model.QuestionFilter{
		Format: listQuery.Filter.Format,
		Limit:  listQuery.Filter.Limit,
		Offset: listQuery.Filter.Offset,
	}

	return h.Repository.List(ctx, filter)
}
