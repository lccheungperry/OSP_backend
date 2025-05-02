package repository

import (
	"context"
	"errors"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
)

var ErrQuestionNotFound = errors.New("question not found")

type QuestionRepository interface {
	Create(ctx context.Context, question *model.Question) error
	Update(ctx context.Context, id string, question *model.Question) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Question, error)
	List(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error)
}
