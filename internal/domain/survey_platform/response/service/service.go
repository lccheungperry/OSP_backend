package service

import (
	"context"
	"fmt"

	bus_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	bus_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	platform_error "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/error"
	question_repo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	resp_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	resp_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/query"
	resp_repo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/repository"
	survey_repo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseService struct {
	responseRepo resp_repo.ResponseRepository
	surveyRepo   survey_repo.SurveyRepository
	questionRepo question_repo.QuestionRepository
}

func NewResponseService(responseRepo resp_repo.ResponseRepository, surveyRepo survey_repo.SurveyRepository, questionRepo question_repo.QuestionRepository) *ResponseService {
	return &ResponseService{
		responseRepo: responseRepo,
		surveyRepo:   surveyRepo,
		questionRepo: questionRepo,
	}
}

func (s *ResponseService) HandleCommand(ctx context.Context, cmd bus_command.Command) error {
	switch c := cmd.(type) {
	case *resp_command.SubmitResponseCommand:
		return s.handleSubmitResponse(ctx, c)
	case *resp_command.DeleteResponseCommand:
		return s.DeleteResponse(ctx, c)
	default:
		return platform_error.NewInvalidCommandTypeError(fmt.Sprintf("unknown command type: %T", cmd))
	}
}

func (s *ResponseService) HandleQuery(ctx context.Context, q bus_query.Query) (interface{}, error) {
	switch query := q.(type) {
	case *resp_query.GetResponsesQuery:
		return s.handleGetResponses(ctx, query)
	default:
		return nil, platform_error.NewInvalidQueryTypeError(fmt.Sprintf("unknown query type: %T", q))
	}
}

func (s *ResponseService) handleSubmitResponse(ctx context.Context, cmd *resp_command.SubmitResponseCommand) error {
	// Validate survey exists
	if _, err := s.surveyRepo.FindByID(ctx, cmd.SurveyID); err != nil {
		return platform_error.ErrSurveyNotFound
	}

	// Get question assignments for the survey
	assignments, err := s.surveyRepo.GetQuestionAssignments(ctx, cmd.SurveyID)
	if err != nil {
		return platform_error.ErrQuestionNotFound
	}

	validQuestions := make(map[primitive.ObjectID]bool)
	for _, assignment := range assignments {
		validQuestions[assignment.QuestionID] = true
	}

	for _, answer := range cmd.Answers {
		if !validQuestions[answer.QuestionID] {
			return platform_error.NewQuestionNotInSurveyError(answer.QuestionID, cmd.SurveyID)
		}

		question, err := s.questionRepo.GetByID(ctx, answer.QuestionID.Hex())
		if err != nil {
			return platform_error.ErrQuestionNotFound
		}

		switch question.Format {
		case "multiple_choice":
			specs, ok := question.Specifications.(map[string]interface{})
			if !ok {
				return platform_error.NewInvalidSpecificationsError(answer.QuestionID)
			}
			options, ok := specs["options"].([]interface{})
			if !ok {
				return platform_error.NewInvalidOptionsError(answer.QuestionID)
			}
			found := false
			answerValue, ok := answer.Value.(string)
			if !ok {
				return platform_error.NewInvalidAnswerTypeError(answer.QuestionID, "multiple_choice")
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
				return platform_error.NewInvalidOptionValueError(answer.QuestionID)
			}

		case "likert":
			specs, ok := question.Specifications.(map[string]interface{})
			if !ok {
				return platform_error.NewInvalidSpecificationsError(answer.QuestionID)
			}
			options, ok := specs["options"].([]interface{})
			if !ok {
				return platform_error.NewInvalidOptionsError(answer.QuestionID)
			}
			scale, ok := answer.Value.(float64)
			if !ok {
				return platform_error.NewInvalidAnswerTypeError(answer.QuestionID, "likert")
			}
			if scale < 1 || scale > float64(len(options)) {
				return platform_error.NewInvalidScaleValueError(answer.QuestionID)
			}
		}
	}

	// Create response
	response := &model.Response{
		SurveyID: cmd.SurveyID,
		Answers:  cmd.Answers,
	}

	if err := s.responseRepo.Create(ctx, response); err != nil {
		return platform_error.NewResponseCreationError(err)
	}

	return nil
}

func (s *ResponseService) handleGetResponses(ctx context.Context, q *resp_query.GetResponsesQuery) (interface{}, error) {
	// Validate survey exists
	if _, err := s.surveyRepo.FindByID(ctx, q.SurveyID); err != nil {
		return nil, platform_error.ErrSurveyNotFound
	}

	responses, total, err := s.responseRepo.FindBySurveyID(ctx, q.SurveyID, q.Skip, q.Limit)
	if err != nil {
		return nil, platform_error.NewResponseRetrievalError(err)
	}

	return struct {
		Responses []*model.Response
		Total     int64
	}{
		Responses: responses,
		Total:     total,
	}, nil
}

func (s *ResponseService) DeleteResponse(ctx context.Context, cmd *resp_command.DeleteResponseCommand) error {
	return s.responseRepo.Delete(ctx, cmd.ID)
}
