package survey

import (
	"context"
	"errors"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSurveyService(t *testing.T) {
	t.Run("CreateSurvey", func(t *testing.T) {
		tests := []struct {
			name          string
			cmd           *command.CreateSurveyCommand
			createFunc    func(ctx context.Context, survey *model.Survey) error
			expectedError error
		}{
			{
				name: "successful survey creation",
				cmd: &command.CreateSurveyCommand{
					Title: "Test Survey",
				},
				createFunc: func(ctx context.Context, survey *model.Survey) error {
					if survey.Title != "Test Survey" {
						return errors.New("unexpected survey title")
					}
					return nil
				},
				expectedError: nil,
			},
			{
				name: "repository error",
				cmd: &command.CreateSurveyCommand{
					Title: "Test Survey",
				},
				createFunc: func(ctx context.Context, survey *model.Survey) error {
					return errors.New("repository error")
				},
				expectedError: errors.New("repository error"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := &MockSurveyRepository{
					CreateFunc: tt.createFunc,
				}
				svc := service.NewSurveyService(repo)

				result, err := svc.HandleCommand(context.Background(), tt.cmd)
				if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
					t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
				}
				if err == nil && result == nil {
					t.Error("HandleCommand() result is nil")
				}
			})
		}
	})

	t.Run("GetSurvey", func(t *testing.T) {
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := &MockSurveyRepository{
					FindByIDFunc: tt.findByIDFunc,
				}
				svc := service.NewSurveyService(repo)

				_, err := svc.HandleQuery(context.Background(), tt.query)
				if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
					t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
				}
			})
		}
	})

	t.Run("ListSurveys", func(t *testing.T) {
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
				svc := service.NewSurveyService(repo)

				result, err := svc.HandleQuery(context.Background(), tt.query)
				if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
					t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
				}
				if err == nil {
					response := result.(struct {
						Surveys []*model.Survey
						Total   int64
					})
					if len(response.Surveys) != 2 || response.Total != 2 {
						t.Errorf("HandleQuery() got = %v items, total %v, want 2 items, total 2", len(response.Surveys), response.Total)
					}
				}
			})
		}
	})
}
