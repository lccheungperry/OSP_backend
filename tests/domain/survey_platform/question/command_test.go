package question_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
)

type MockQuestionRepository struct {
	CreateFunc  func(ctx context.Context, question *model.Question) error
	UpdateFunc  func(ctx context.Context, id string, question *model.Question) error
	DeleteFunc  func(ctx context.Context, id string) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Question, error)
	ListFunc    func(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error)
}

func (m *MockQuestionRepository) Create(ctx context.Context, question *model.Question) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, question)
	}
	return nil
}

func (m *MockQuestionRepository) Update(ctx context.Context, id string, question *model.Question) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, question)
	}
	return nil
}

func (m *MockQuestionRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockQuestionRepository) GetByID(ctx context.Context, id string) (*model.Question, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockQuestionRepository) List(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter)
	}
	return nil, nil
}

func TestCreateQuestionHandler(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.CreateQuestionCommand
		createFunc    func(ctx context.Context, question *model.Question) error
		expectedError error
	}{
		{
			name: "successful question creation",
			cmd: &command.CreateQuestionCommand{
				Title:          "Test Question",
				Format:         "textbox",
				Specifications: map[string]interface{}{},
			},
			createFunc: func(ctx context.Context, question *model.Question) error {
				if question.Title != "Test Question" {
					return errors.New("unexpected question title")
				}
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.CreateQuestionCommand{
				Title:          "Test Question",
				Format:         "textbox",
				Specifications: map[string]interface{}{},
			},
			createFunc: func(ctx context.Context, question *model.Question) error {
				return errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockQuestionRepository{
				CreateFunc: tt.createFunc,
			}
			handler := &command.CreateQuestionHandler{
				Repository: repo,
			}

			err := handler.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestUpdateQuestionHandler(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.UpdateQuestionCommand
		updateFunc    func(ctx context.Context, id string, question *model.Question) error
		expectedError error
	}{
		{
			name: "successful question update",
			cmd: &command.UpdateQuestionCommand{
				ID:             "123",
				Title:          "Updated Question",
				Format:         "textbox",
				Specifications: map[string]interface{}{},
			},
			updateFunc: func(ctx context.Context, id string, question *model.Question) error {
				if question.Title != "Updated Question" {
					return errors.New("unexpected question title")
				}
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.UpdateQuestionCommand{
				ID:             "123",
				Title:          "Updated Question",
				Format:         "textbox",
				Specifications: map[string]interface{}{},
			},
			updateFunc: func(ctx context.Context, id string, question *model.Question) error {
				return errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockQuestionRepository{
				UpdateFunc: tt.updateFunc,
			}
			handler := &command.UpdateQuestionHandler{
				Repository: repo,
			}

			err := handler.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestDeleteQuestionHandler(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.DeleteQuestionCommand
		deleteFunc    func(ctx context.Context, id string) error
		expectedError error
	}{
		{
			name: "successful question deletion",
			cmd: &command.DeleteQuestionCommand{
				ID: "123",
			},
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.DeleteQuestionCommand{
				ID: "123",
			},
			deleteFunc: func(ctx context.Context, id string) error {
				return errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockQuestionRepository{
				DeleteFunc: tt.deleteFunc,
			}
			handler := &command.DeleteQuestionHandler{
				Repository: repo,
			}

			err := handler.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
