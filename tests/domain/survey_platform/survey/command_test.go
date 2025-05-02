package survey

import (
	"context"
	"errors"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateSurveyHandler(t *testing.T) {
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
			handler := &command.CreateSurveyHandler{
				Repository: repo,
			}

			_, err := handler.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestUpdateSurveyHandler(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.UpdateSurveyCommand
		updateFunc    func(ctx context.Context, survey *model.Survey) error
		expectedError error
	}{
		{
			name: "successful survey update",
			cmd: &command.UpdateSurveyCommand{
				ID:    primitive.NewObjectID(),
				Title: "Updated Survey",
			},
			updateFunc: func(ctx context.Context, survey *model.Survey) error {
				if survey.Title != "Updated Survey" {
					return errors.New("unexpected survey title")
				}
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.UpdateSurveyCommand{
				ID:    primitive.NewObjectID(),
				Title: "Updated Survey",
			},
			updateFunc: func(ctx context.Context, survey *model.Survey) error {
				return errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockSurveyRepository{
				UpdateFunc: tt.updateFunc,
			}
			handler := &command.UpdateSurveyHandler{
				Repository: repo,
			}

			_, err := handler.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestDeleteSurveyHandler(t *testing.T) {
	tests := []struct {
		name          string
		cmd           *command.DeleteSurveyCommand
		deleteFunc    func(ctx context.Context, id primitive.ObjectID) error
		expectedError error
	}{
		{
			name: "successful survey deletion",
			cmd: &command.DeleteSurveyCommand{
				ID: primitive.NewObjectID(),
			},
			deleteFunc: func(ctx context.Context, id primitive.ObjectID) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			cmd: &command.DeleteSurveyCommand{
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
			repo := &MockSurveyRepository{
				DeleteFunc: tt.deleteFunc,
			}
			handler := &command.DeleteSurveyHandler{
				Repository: repo,
			}

			_, err := handler.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
