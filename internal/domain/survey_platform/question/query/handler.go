package query

import (
	"context"
	"fmt"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	platform_error "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/error"
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
		return nil, platform_error.NewInvalidQueryTypeError(q)
	}

	question, err := h.Repository.GetByID(ctx, getQuery.ID)
	if err != nil {
		return nil, err
	}

	return question, nil
}

type ListQuestionsHandler struct {
	Repository repository.QuestionRepository
}

func (h *ListQuestionsHandler) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	listQuery, ok := q.(*ListQuestionsQuery)
	if !ok {
		return nil, platform_error.NewInvalidQueryTypeError(q)
	}

	filter := &model.QuestionFilter{
		Format: listQuery.Filter.Format,
		Limit:  listQuery.Filter.Limit,
		Offset: listQuery.Filter.Offset,
	}

	questions, err := h.Repository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return questions, nil
}
