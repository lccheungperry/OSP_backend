package handlers

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockSurveyService struct {
	HandleCommandFunc func(ctx context.Context, cmd command.Command) error
	HandleQueryFunc   func(ctx context.Context, q query.Query) (interface{}, error)
}

func (m *MockSurveyService) HandleCommand(ctx context.Context, cmd command.Command) error {
	if m.HandleCommandFunc != nil {
		return m.HandleCommandFunc(ctx, cmd)
	}
	return nil
}

func (m *MockSurveyService) HandleQuery(ctx context.Context, q query.Query) (interface{}, error) {
	if m.HandleQueryFunc != nil {
		return m.HandleQueryFunc(ctx, q)
	}
	return &model.Survey{
		ID:    primitive.NewObjectID(),
		Title: "Test Survey",
	}, nil
}
