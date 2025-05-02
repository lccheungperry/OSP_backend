package question_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetQuestionHandler(t *testing.T) {
	tests := []struct {
		name          string
		query         *query.GetQuestionQuery
		getByIDFunc   func(ctx context.Context, id string) (*model.Question, error)
		expectedError error
	}{
		{
			name: "successful question retrieval",
			query: &query.GetQuestionQuery{
				ID: "123",
			},
			getByIDFunc: func(ctx context.Context, id string) (*model.Question, error) {
				return &model.Question{
					ID:             primitive.NewObjectID(),
					Title:          "Test Question",
					Format:         model.FormatTextbox,
					Specifications: map[string]interface{}{},
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "question not found",
			query: &query.GetQuestionQuery{
				ID: "123",
			},
			getByIDFunc: func(ctx context.Context, id string) (*model.Question, error) {
				return nil, repository.ErrQuestionNotFound
			},
			expectedError: repository.ErrQuestionNotFound,
		},
		{
			name: "repository error",
			query: &query.GetQuestionQuery{
				ID: "123",
			},
			getByIDFunc: func(ctx context.Context, id string) (*model.Question, error) {
				return nil, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockQuestionRepository{
				GetByIDFunc: tt.getByIDFunc,
			}
			handler := &query.GetQuestionHandler{
				Repository: repo,
			}

			_, err := handler.HandleQuery(context.Background(), tt.query)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestListQuestionsHandler(t *testing.T) {
	tests := []struct {
		name          string
		query         *query.ListQuestionsQuery
		listFunc      func(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error)
		expectedError error
	}{
		{
			name: "successful question listing",
			query: &query.ListQuestionsQuery{
				Filter: &model.QuestionFilter{
					Format: model.FormatTextbox,
					Limit:  10,
					Offset: 0,
				},
			},
			listFunc: func(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
				return []*model.Question{
					{
						ID:             primitive.NewObjectID(),
						Title:          "Test Question 1",
						Format:         model.FormatTextbox,
						Specifications: map[string]interface{}{},
					},
					{
						ID:             primitive.NewObjectID(),
						Title:          "Test Question 2",
						Format:         model.FormatTextbox,
						Specifications: map[string]interface{}{},
					},
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			query: &query.ListQuestionsQuery{
				Filter: &model.QuestionFilter{
					Format: model.FormatTextbox,
					Limit:  10,
					Offset: 0,
				},
			},
			listFunc: func(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
				return nil, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockQuestionRepository{
				ListFunc: tt.listFunc,
			}
			handler := &query.ListQuestionsHandler{
				Repository: repo,
			}

			_, err := handler.HandleQuery(context.Background(), tt.query)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
