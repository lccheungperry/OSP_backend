package survey

import (
	"context"
	"errors"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockSurveyRepository struct {
	CreateFunc      func(ctx context.Context, survey *model.Survey) error
	UpdateFunc      func(ctx context.Context, survey *model.Survey) error
	DeleteFunc      func(ctx context.Context, id primitive.ObjectID) error
	FindByIDFunc    func(ctx context.Context, id primitive.ObjectID) (*model.Survey, error)
	FindByTokenFunc func(ctx context.Context, token string) (*model.Survey, error)
	ListFunc        func(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error)
}

func (m *MockSurveyRepository) Create(ctx context.Context, survey *model.Survey) error {
	if survey == nil {
		return errors.New("survey cannot be nil")
	}
	if survey.Title == "" {
		return errors.New("survey title cannot be empty")
	}
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, survey)
	}

	survey.ID = primitive.NewObjectID()
	survey.CreatedAt = time.Now()
	survey.UpdatedAt = time.Now()
	return nil
}

func (m *MockSurveyRepository) Update(ctx context.Context, survey *model.Survey) error {
	if survey == nil {
		return errors.New("survey cannot be nil")
	}
	if survey.ID.IsZero() {
		return errors.New("survey ID cannot be zero")
	}
	if survey.Title == "" {
		return errors.New("survey title cannot be empty")
	}
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, survey)
	}

	survey.UpdatedAt = time.Now()
	return nil
}

func (m *MockSurveyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if id.IsZero() {
		return errors.New("survey ID cannot be zero")
	}
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	return nil
}

func (m *MockSurveyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Survey, error) {
	if id.IsZero() {
		return nil, errors.New("survey ID cannot be zero")
	}
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}

	return &model.Survey{
		ID:        id,
		Title:     "Default Survey",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockSurveyRepository) FindByToken(ctx context.Context, token string) (*model.Survey, error) {
	if token == "" {
		return nil, errors.New("survey token cannot be empty")
	}
	if m.FindByTokenFunc != nil {
		return m.FindByTokenFunc(ctx, token)
	}

	return &model.Survey{
		ID:        primitive.NewObjectID(),
		Title:     "Default Survey",
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockSurveyRepository) List(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error) {
	if skip < 0 {
		return nil, 0, errors.New("skip cannot be negative")
	}
	if limit < 0 {
		return nil, 0, errors.New("limit cannot be negative")
	}
	if m.ListFunc != nil {
		return m.ListFunc(ctx, skip, limit)
	}

	return []*model.Survey{
		{
			ID:        primitive.NewObjectID(),
			Title:     "Default Survey 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Title:     "Default Survey 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, 2, nil
}
