package service

import (
	"context"
	"time"

	bus_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	bus_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	question_repo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	resp_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	resp_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/query"
	resp_repository "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/repository"
	survey_repo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseService struct {
	repo         resp_repository.ResponseRepository
	surveyRepo   survey_repo.SurveyRepository
	questionRepo question_repo.QuestionRepository
}

func NewResponseService(repo resp_repository.ResponseRepository, surveyRepo survey_repo.SurveyRepository, questionRepo question_repo.QuestionRepository) *ResponseService {
	return &ResponseService{
		repo:         repo,
		surveyRepo:   surveyRepo,
		questionRepo: questionRepo,
	}
}

func (s *ResponseService) SubmitResponse(ctx context.Context, cmd *resp_command.SubmitResponseCommand) (*model.Response, error) {
	_, err := s.surveyRepo.FindByID(ctx, cmd.SurveyID)
	if err != nil {
		return nil, NewSurveyNotFoundError(err)
	}

	assignments, err := s.surveyRepo.GetQuestionAssignments(ctx, cmd.SurveyID)
	if err != nil {
		return nil, NewQuestionNotFoundError(err)
	}

	validQuestions := make(map[primitive.ObjectID]bool)
	for _, assignment := range assignments {
		validQuestions[assignment.QuestionID] = true
	}

	for _, answer := range cmd.Answers {
		if !validQuestions[answer.QuestionID] {
			return nil, NewQuestionNotInSurveyError(answer.QuestionID, cmd.SurveyID)
		}

		question, err := s.questionRepo.GetByID(ctx, answer.QuestionID.Hex())
		if err != nil {
			return nil, NewQuestionNotFoundError(err)
		}

		switch question.Format {
		case "multiple_choice":
			specs, ok := question.Specifications.(map[string]interface{})
			if !ok {
				return nil, NewInvalidSpecificationsError(answer.QuestionID)
			}
			options, ok := specs["options"].([]interface{})
			if !ok {
				return nil, NewInvalidOptionsError(answer.QuestionID)
			}
			found := false
			answerValue, ok := answer.Value.(string)
			if !ok {
				return nil, NewInvalidAnswerTypeError(answer.QuestionID, "multiple_choice")
			}
			for _, opt := range options {
				option, ok := opt.(map[string]interface{})
				if !ok {
					continue
				}
				if option["id"] == answerValue {
					found = true
					break
				}
			}
			if !found {
				return nil, NewInvalidOptionValueError(answer.QuestionID)
			}

		case "likert":
			specs, ok := question.Specifications.(map[string]interface{})
			if !ok {
				return nil, NewInvalidSpecificationsError(answer.QuestionID)
			}
			options, ok := specs["options"].([]interface{})
			if !ok {
				return nil, NewInvalidOptionsError(answer.QuestionID)
			}
			scale, ok := answer.Value.(float64)
			if !ok {
				return nil, NewInvalidAnswerTypeError(answer.QuestionID, "likert")
			}
			if scale < 1 || scale > float64(len(options)) {
				return nil, NewInvalidScaleValueError(answer.QuestionID)
			}
		}
	}

	response := &model.Response{
		SurveyID:  cmd.SurveyID,
		Answers:   cmd.Answers,
		CreatedAt: time.Now(),
	}

	err = s.repo.Create(ctx, response)
	if err != nil {
		return nil, NewResponseCreationError(err)
	}
	return response, nil
}

func (s *ResponseService) GetResponses(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error) {
	_, err := s.surveyRepo.FindByID(ctx, surveyID)
	if err != nil {
		return nil, 0, NewSurveyNotFoundError(err)
	}

	responses, total, err := s.repo.FindBySurveyID(ctx, surveyID, skip, limit)
	if err != nil {
		return nil, 0, NewResponseCreationError(err)
	}

	return responses, total, nil
}

func (s *ResponseService) DeleteResponse(ctx context.Context, cmd *resp_command.DeleteResponseCommand) error {
	return s.repo.Delete(ctx, cmd.ID)
}

func (s *ResponseService) HandleCommand(ctx context.Context, cmd bus_command.Command) error {
	switch c := cmd.(type) {
	case *resp_command.SubmitResponseCommand:
		_, err := s.SubmitResponse(ctx, c)
		return err
	case *resp_command.DeleteResponseCommand:
		return s.DeleteResponse(ctx, c)
	default:
		return nil
	}
}

func (s *ResponseService) HandleQuery(ctx context.Context, q bus_query.Query) (interface{}, error) {
	switch query := q.(type) {
	case *resp_query.GetResponsesQuery:
		responses, total, err := s.repo.FindBySurveyID(ctx, query.SurveyID, query.Skip, query.Limit)
		if err != nil {
			return nil, NewResponseCreationError(err)
		}
		return struct {
			Responses []*model.Response
			Total     int64
		}{
			Responses: responses,
			Total:     total,
		}, nil
	default:
		return nil, nil
	}
}
