package response_test

import (
	"context"
	"errors"
	"testing"
	"time"

	question_model "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/service"
	survey_model "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Survey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type SurveyQuestionAssignment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	SurveyID   primitive.ObjectID `bson:"survey_id"`
	QuestionID primitive.ObjectID `bson:"question_id"`
	Order      int                `bson:"order"`
}

type Question struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Title          string             `bson:"title"`
	Format         string             `bson:"format"`
	Specifications interface{}        `bson:"specifications"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

type QuestionFilter struct {
	Format string
}

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

type MockSurveyRepository struct {
	FindByIDFunc                 func(ctx context.Context, id primitive.ObjectID) (*survey_model.Survey, error)
	GetQuestionAssignmentsFunc   func(ctx context.Context, surveyID primitive.ObjectID) ([]*survey_model.SurveyQuestionAssignment, error)
	CreateFunc                   func(ctx context.Context, survey *survey_model.Survey) error
	FindByTokenFunc              func(ctx context.Context, token string) (*survey_model.Survey, error)
	ListFunc                     func(ctx context.Context, skip, limit int64) ([]*survey_model.Survey, int64, error)
	UpdateFunc                   func(ctx context.Context, survey *survey_model.Survey) error
	DeleteFunc                   func(ctx context.Context, id primitive.ObjectID) error
	CreateQuestionAssignmentFunc func(ctx context.Context, assignment *survey_model.SurveyQuestionAssignment) error
}

func (m *MockSurveyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*survey_model.Survey, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockSurveyRepository) GetQuestionAssignments(ctx context.Context, surveyID primitive.ObjectID) ([]*survey_model.SurveyQuestionAssignment, error) {
	if m.GetQuestionAssignmentsFunc != nil {
		return m.GetQuestionAssignmentsFunc(ctx, surveyID)
	}
	return nil, nil
}

func (m *MockSurveyRepository) Create(ctx context.Context, survey *survey_model.Survey) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, survey)
	}
	return nil
}

func (m *MockSurveyRepository) FindByToken(ctx context.Context, token string) (*survey_model.Survey, error) {
	if m.FindByTokenFunc != nil {
		return m.FindByTokenFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockSurveyRepository) List(ctx context.Context, skip, limit int64) ([]*survey_model.Survey, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, skip, limit)
	}
	return nil, 0, nil
}

func (m *MockSurveyRepository) Update(ctx context.Context, survey *survey_model.Survey) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, survey)
	}
	return nil
}

func (m *MockSurveyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockSurveyRepository) CreateQuestionAssignment(ctx context.Context, assignment *survey_model.SurveyQuestionAssignment) error {
	if m.CreateQuestionAssignmentFunc != nil {
		return m.CreateQuestionAssignmentFunc(ctx, assignment)
	}
	return nil
}

type MockQuestionRepository struct {
	GetByIDFunc func(ctx context.Context, id string) (*question_model.Question, error)
	CreateFunc  func(ctx context.Context, question *question_model.Question) error
	UpdateFunc  func(ctx context.Context, id string, question *question_model.Question) error
	DeleteFunc  func(ctx context.Context, id string) error
	ListFunc    func(ctx context.Context, filter *question_model.QuestionFilter) ([]*question_model.Question, error)
}

func (m *MockQuestionRepository) GetByID(ctx context.Context, id string) (*question_model.Question, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockQuestionRepository) Create(ctx context.Context, question *question_model.Question) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, question)
	}
	return nil
}

func (m *MockQuestionRepository) Update(ctx context.Context, id string, question *question_model.Question) error {
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

func (m *MockQuestionRepository) List(ctx context.Context, filter *question_model.QuestionFilter) ([]*question_model.Question, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter)
	}
	return nil, nil
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
			responseRepo := &MockResponseRepository{
				CreateFunc: tt.createFunc,
			}
			surveyRepo := &MockSurveyRepository{
				FindByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*survey_model.Survey, error) {
					return &survey_model.Survey{
						ID: tt.cmd.SurveyID,
					}, nil
				},
				GetQuestionAssignmentsFunc: func(ctx context.Context, surveyID primitive.ObjectID) ([]*survey_model.SurveyQuestionAssignment, error) {
					assignments := make([]*survey_model.SurveyQuestionAssignment, len(tt.cmd.Answers))
					for i, answer := range tt.cmd.Answers {
						assignments[i] = &survey_model.SurveyQuestionAssignment{
							SurveyID:   surveyID,
							QuestionID: answer.QuestionID,
							Order:      i + 1,
						}
					}
					return assignments, nil
				},
			}
			questionRepo := &MockQuestionRepository{
				GetByIDFunc: func(ctx context.Context, id string) (*question_model.Question, error) {
					return &question_model.Question{
						ID:     primitive.NewObjectID(),
						Format: "text",
					}, nil
				},
			}
			svc := service.NewResponseService(responseRepo, surveyRepo, questionRepo)

			err := svc.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestGetResponses(t *testing.T) {
	tests := []struct {
		name               string
		query              *query.GetResponsesQuery
		findBySurveyIDFunc func(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error)
		expectedError      error
	}{
		{
			name: "successful response retrieval",
			query: &query.GetResponsesQuery{
				SurveyID: primitive.NewObjectID(),
				Skip:     0,
				Limit:    10,
			},
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
			name: "repository error",
			query: &query.GetResponsesQuery{
				SurveyID: primitive.NewObjectID(),
				Skip:     0,
				Limit:    10,
			},
			findBySurveyIDFunc: func(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error) {
				return nil, 0, errors.New("repository error")
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRepo := &MockResponseRepository{
				FindBySurveyIDFunc: tt.findBySurveyIDFunc,
			}
			surveyRepo := &MockSurveyRepository{
				FindByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*survey_model.Survey, error) {
					return &survey_model.Survey{}, nil
				},
			}
			questionRepo := &MockQuestionRepository{}
			svc := service.NewResponseService(responseRepo, surveyRepo, questionRepo)

			result, err := svc.HandleQuery(context.Background(), tt.query)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("HandleQuery() error = %v, expectedError %v", err, tt.expectedError)
			}
			if err == nil {
				response := result.(struct {
					Responses []*model.Response
					Total     int64
				})
				if len(response.Responses) != 1 {
					t.Errorf("Expected 1 response, got %d", len(response.Responses))
				}
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
			responseRepo := &MockResponseRepository{
				DeleteFunc: tt.deleteFunc,
			}
			surveyRepo := &MockSurveyRepository{}
			questionRepo := &MockQuestionRepository{}
			svc := service.NewResponseService(responseRepo, surveyRepo, questionRepo)

			err := svc.HandleCommand(context.Background(), tt.cmd)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("DeleteResponse() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
