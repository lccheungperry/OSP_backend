package question

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockQuestionService struct {
	HandleCommandFunc func(ctx context.Context, cmd command.Command) error
	HandleQueryFunc   func(ctx context.Context, q query.Query) (interface{}, error)
}

func (m *MockQuestionService) HandleCommand(ctx context.Context, cmd command.Command) error {
	if m.HandleCommandFunc != nil {
		return m.HandleCommandFunc(ctx, cmd)
	}
	return nil
}

func (m *MockQuestionService) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	if m.HandleQueryFunc != nil {
		return m.HandleQueryFunc(ctx, q)
	}
	return &model.Question{
		ID:    primitive.NewObjectID(),
		Title: "Test Question",
	}, nil
}
