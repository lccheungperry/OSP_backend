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

func validateMultipleChoiceAnswer(answerValue interface{}, options []interface{}, questionID primitive.ObjectID) error {
	answerStr, ok := answerValue.(string)
	if !ok {
		return platform_error.NewInvalidAnswerTypeError(questionID, "multiple_choice")
	}

	for _, opt := range options {
		optMap, ok := opt.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := optMap["id"].(string); ok && id == answerStr {
			return nil
		}
	}
	return platform_error.NewInvalidOptionValueError(questionID)
}

func validateLikertAnswer(answerValue interface{}, options []interface{}, questionID primitive.ObjectID) error {
	var scale float64
	switch v := answerValue.(type) {
	case float64:
		scale = v
	case int32:
		scale = float64(v)
	case int:
		scale = float64(v)
	default:
		return platform_error.NewInvalidAnswerTypeError(questionID, "likert")
	}

	maxScale := 0.0
	for _, opt := range options {
		if optDoc, ok := opt.(primitive.D); ok {
			var scaleValue float64
			for _, elem := range optDoc {
				if elem.Key == "scale" {
					if value, ok := elem.Value.(float64); ok {
						scaleValue = value
					}
				}
			}
			if scaleValue > maxScale {
				maxScale = scaleValue
			}
		}
	}

	if maxScale == 0 {
		return platform_error.NewInvalidScaleValueError(questionID, scale, maxScale)
	}

	if scale < 1 || scale > maxScale {
		return platform_error.NewInvalidScaleValueError(questionID, scale, maxScale)
	}

	return nil
}

func getOptionsFromSpecs(specs interface{}, questionID primitive.ObjectID) ([]interface{}, error) {
	switch v := specs.(type) {
	case primitive.D:
		for _, elem := range v {
			if elem.Key == "options" {
				if options, ok := elem.Value.(primitive.A); ok {
					return []interface{}(options), nil
				}
			}
		}
	case []interface{}:
		return v, nil
	case map[string]interface{}:
		if options, ok := v["options"].([]interface{}); ok {
			return options, nil
		}
	}
	return nil, platform_error.NewInvalidOptionsError(questionID)
}

func (s *ResponseService) handleSubmitResponse(ctx context.Context, cmd *resp_command.SubmitResponseCommand) error {
	if _, err := s.surveyRepo.FindByID(ctx, cmd.SurveyID); err != nil {
		return platform_error.ErrSurveyNotFound
	}

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

		options, err := getOptionsFromSpecs(question.Specifications, answer.QuestionID)
		if err != nil {
			return err
		}

		switch question.Format {
		case "multiple_choice":
			if err := validateMultipleChoiceAnswer(answer.Value, options, answer.QuestionID); err != nil {
				return err
			}
		case "likert":
			if err := validateLikertAnswer(answer.Value, options, answer.QuestionID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid question format: %s", question.Format)
		}
	}

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
