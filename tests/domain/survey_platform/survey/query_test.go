package survey

import (
	"context"
	"errors"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetSurveyHandler(t *testing.T) {
	tests := []struct {
		name          string
		query         *query.GetSurveyQuery
		findByIDFunc  func(ctx context.Context, id primitive.ObjectID) (*model.Survey, error)
		expectedError error
	}{
		{
			name: "successful survey retrieval",
			query: &query.GetSurveyQuery{
				ID: primitive.NewObjectID().Hex(),
			},
			findByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*model.Survey, error) {
				return &model.Survey{
					ID:    id,
					Title: "Test Survey",
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "invalid survey ID",
			query: &query.GetSurveyQuery{
				ID: "invalid-id",
			},
			findByIDFunc:  nil,
			expectedError: errors.New("invalid survey ID"),
		},
		{
			name: "repository error",
			query: &query.GetSurveyQuery{
				ID: primitive.NewObjectID().Hex(),
			},
			findByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*model.Survey, error) {
				return nil, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockSurveyRepository{
				FindByIDFunc: tt.findByIDFunc,
			}
			handler := &query.GetSurveyHandler{
				Repository: repo,
			}

			_, err := handler.HandleQuery(context.Background(), tt.query)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestGetSurveyByTokenHandler(t *testing.T) {
	tests := []struct {
		name            string
		query           *query.GetSurveyByTokenQuery
		findByTokenFunc func(ctx context.Context, token string) (*model.Survey, error)
		expectedError   error
	}{
		{
			name: "successful survey retrieval by token",
			query: &query.GetSurveyByTokenQuery{
				Token: "abcde",
			},
			findByTokenFunc: func(ctx context.Context, token string) (*model.Survey, error) {
				return &model.Survey{
					ID:    primitive.NewObjectID(),
					Title: "Test Survey",
					Token: token,
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			query: &query.GetSurveyByTokenQuery{
				Token: "abcde",
			},
			findByTokenFunc: func(ctx context.Context, token string) (*model.Survey, error) {
				return nil, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockSurveyRepository{
				FindByTokenFunc: tt.findByTokenFunc,
			}
			handler := &query.GetSurveyByTokenHandler{
				Repository: repo,
			}

			_, err := handler.HandleQuery(context.Background(), tt.query)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestListSurveysHandler(t *testing.T) {
	tests := []struct {
		name          string
		query         *query.ListSurveysQuery
		listFunc      func(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error)
		expectedError error
	}{
		{
			name: "successful survey listing",
			query: &query.ListSurveysQuery{
				Filter: &model.SurveyFilter{
					Limit:  10,
					Offset: 0,
				},
			},
			listFunc: func(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error) {
				return []*model.Survey{
					{
						ID:    primitive.NewObjectID(),
						Title: "Test Survey 1",
					},
					{
						ID:    primitive.NewObjectID(),
						Title: "Test Survey 2",
					},
				}, 2, nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			query: &query.ListSurveysQuery{
				Filter: &model.SurveyFilter{
					Limit:  10,
					Offset: 0,
				},
			},
			listFunc: func(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error) {
				return nil, 0, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockSurveyRepository{
				ListFunc: tt.listFunc,
			}
			handler := &query.ListSurveysHandler{
				Repository: repo,
			}

			result, err := handler.HandleQuery(context.Background(), tt.query)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
			}
			if err == nil {
				response := result.(struct {
					Surveys []*model.Survey `json:"surveys"`
					Total   int64           `json:"total"`
				})
				if len(response.Surveys) != 2 || response.Total != 2 {
					t.Errorf("HandleQuery() got = %v items, total %v, want 2 items, total 2", len(response.Surveys), response.Total)
				}
			}
		})
	}
}
