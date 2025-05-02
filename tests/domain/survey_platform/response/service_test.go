package response_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockResponseRepository struct {
	CreateFunc         func(ctx context.Context, response *model.Response) error
	DeleteFunc         func(ctx context.Context, id primitive.ObjectID) error
	FindByIDFunc       func(ctx context.Context, id primitive.ObjectID) (*model.Response, error)
	FindBySurveyIDFunc func(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error)
}

func (m *MockResponseRepository) Create(ctx context.Context, response *model.Response) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, response)
	}
	return nil
}

func (m *MockResponseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockResponseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Response, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockResponseRepository) FindBySurveyID(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error) {
	if m.FindBySurveyIDFunc != nil {
		return m.FindBySurveyIDFunc(ctx, surveyID, skip, limit)
	}
	return nil, 0, nil
}

func TestSubmitResponse(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.SubmitResponseCommand
		createFunc    func(ctx context.Context, response *model.Response) error
		expectedError error
	}{
		{
			name: "successful response submission",
			cmd: &command.SubmitResponseCommand{
				SurveyID: primitive.NewObjectID(),
				Answers: []model.Answer{
					{
						QuestionID: primitive.NewObjectID(),
						Value:      "Test Answer",
					},
				},
			},
			createFunc: func(ctx context.Context, response *model.Response) error {
				if len(response.Answers) != 1 {
					return errors.New("unexpected number of answers")
				}
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.SubmitResponseCommand{
				SurveyID: primitive.NewObjectID(),
				Answers: []model.Answer{
					{
						QuestionID: primitive.NewObjectID(),
						Value:      "Test Answer",
					},
				},
			},
			createFunc: func(ctx context.Context, response *model.Response) error {
				return errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockResponseRepository{
				CreateFunc: tt.createFunc,
			}
			svc := service.NewResponseService(repo)

			_, err := svc.SubmitResponse(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("SubmitResponse() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestGetResponses(t *testing.T) {
	tests := []struct {
		name               string
		surveyID           primitive.ObjectID
		page               int64
		pageSize           int64
		findBySurveyIDFunc func(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error)
		expectedError      error
	}{
		{
			name:     "successful response retrieval",
			surveyID: primitive.NewObjectID(),
			page:     1,
			pageSize: 10,
			findBySurveyIDFunc: func(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error) {
				return []*model.Response{
					{
						ID:        primitive.NewObjectID(),
						SurveyID:  surveyID,
						Answers:   []model.Answer{},
						CreatedAt: time.Now(),
					},
				}, 1, nil
			},
			expectedError: nil,
		},
		{
			name:     "repository error",
			surveyID: primitive.NewObjectID(),
			page:     1,
			pageSize: 10,
			findBySurveyIDFunc: func(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error) {
				return nil, 0, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockResponseRepository{
				FindBySurveyIDFunc: tt.findBySurveyIDFunc,
			}
			svc := service.NewResponseService(repo)

			_, _, err := svc.GetResponses(context.Background(), tt.surveyID, tt.page, tt.pageSize)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("GetResponses() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestDeleteResponse(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.DeleteResponseCommand
		deleteFunc    func(ctx context.Context, id primitive.ObjectID) error
		expectedError error
	}{
		{
			name: "successful response deletion",
			cmd: &command.DeleteResponseCommand{
				ID: primitive.NewObjectID(),
			},
			deleteFunc: func(ctx context.Context, id primitive.ObjectID) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.DeleteResponseCommand{
				ID: primitive.NewObjectID(),
			},
			deleteFunc: func(ctx context.Context, id primitive.ObjectID) error {
				return errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockResponseRepository{
				DeleteFunc: tt.deleteFunc,
			}
			svc := service.NewResponseService(repo)

			err := svc.DeleteResponse(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("DeleteResponse() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
