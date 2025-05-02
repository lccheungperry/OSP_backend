package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	bus_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	platform_error "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/error"
	question_model "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	question_repo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	resp_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	resp_model "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
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

func (s *ResponseService) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	switch c := cmd.(type) {
	case *resp_command.SubmitResponseCommand:
		return s.handleSubmitResponse(ctx, c)
	case *resp_command.DeleteResponseCommand:
		return nil, s.DeleteResponse(ctx, c)
	default:
		return nil, platform_error.NewInvalidCommandTypeError(cmd)
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
		if str, err := convertToString(answerValue); err == nil {
			answerStr = str
		} else {
			return platform_error.NewInvalidAnswerTypeError(questionID, "multiple_choice")
		}
	}

	validOptions := make(map[string]bool)
	for _, opt := range options {
		optMap, ok := opt.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := optMap["id"].(string); ok {
			validOptions[id] = true
		}
	}

	if !validOptions[answerStr] {
		return platform_error.NewInvalidOptionValueError(questionID)
	}

	return nil
}

func convertToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	case json.Number:
		return v.String(), nil
	default:
		if bytes, err := json.Marshal(value); err == nil {
			return string(bytes), nil
		}
		return "", fmt.Errorf("unable to convert %T to string", value)
	}
}

func convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case json.Number:
		return v.Float64()
	default:
		if str, ok := value.(string); ok {
			if num, err := strconv.ParseFloat(str, 64); err == nil {
				return num, nil
			}
		}
		return 0, fmt.Errorf("unable to convert %T to float64", value)
	}
}

func convertToMap(value interface{}) (map[string]interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case primitive.M:
		return map[string]interface{}(v), nil
	case primitive.D:
		return v.Map(), nil
	default:
		if bytes, err := json.Marshal(value); err == nil {
			var m map[string]interface{}
			if err := json.Unmarshal(bytes, &m); err == nil {
				return m, nil
			}
		}
		return nil, fmt.Errorf("failed to convert %T to map", value)
	}
}

func convertToArray(value interface{}) ([]interface{}, error) {
	switch v := value.(type) {
	case []interface{}:
		return v, nil
	case primitive.A:
		return []interface{}(v), nil
	default:
		if bytes, err := json.Marshal(value); err == nil {
			var arr []interface{}
			if err := json.Unmarshal(bytes, &arr); err == nil {
				return arr, nil
			}
		}
		return nil, fmt.Errorf("failed to convert %T to array", value)
	}
}

func validateLikertScale(scale float64, options []interface{}, questionID primitive.ObjectID) error {
	maxScale := float64(len(options))
	if maxScale == 0 {
		return platform_error.NewInvalidScaleValueError(questionID, scale, maxScale)
	}

	if scale < 1 || scale > maxScale {
		return platform_error.NewInvalidScaleValueError(questionID, scale, maxScale)
	}

	return nil
}

func validateLikertAnswer(answerValue interface{}, options []interface{}, questionID primitive.ObjectID) error {
	scale, err := convertToFloat64(answerValue)
	if err != nil {
		return platform_error.NewInvalidAnswerTypeError(questionID, "likert")
	}

	return validateLikertScale(scale, options, questionID)
}

func getOptionsFromSpecs(specs interface{}, questionID primitive.ObjectID) ([]interface{}, error) {
	if likertSpecs, ok := specs.(question_model.LikertSpecification); ok {
		return convertLikertSpecsToOptions(likertSpecs), nil
	}

	specsMap, err := convertToMap(specs)
	if err != nil {
		return nil, platform_error.NewInvalidOptionsError(questionID)
	}

	optionsRaw, ok := specsMap["options"]
	if !ok {
		return nil, platform_error.NewInvalidOptionsError(questionID)
	}

	options, err := convertToArray(optionsRaw)
	if err != nil {
		return nil, platform_error.NewInvalidOptionsError(questionID)
	}

	return options, nil
}

func convertLikertSpecsToOptions(specs question_model.LikertSpecification) []interface{} {
	options := make([]interface{}, len(specs.Options))
	for i, opt := range specs.Options {
		options[i] = map[string]interface{}{
			"id":    opt.ID,
			"label": opt.Label,
			"scale": opt.Scale,
		}
	}
	return options
}

func (s *ResponseService) handleSubmitResponse(ctx context.Context, cmd *resp_command.SubmitResponseCommand) (interface{}, error) {
	survey, err := s.surveyRepo.FindByID(ctx, cmd.SurveyID)
	if err != nil {
		return nil, platform_error.ErrSurveyNotFound
	}

	assignments, err := s.surveyRepo.GetQuestionAssignments(ctx, cmd.SurveyID)
	if err != nil {
		return nil, platform_error.ErrQuestionNotFound
	}

	response := &resp_model.Response{
		ID:        primitive.NewObjectID(),
		SurveyID:  survey.ID,
		Answers:   cmd.Answers,
		CreatedAt: time.Now(),
	}

	validQuestions := make(map[primitive.ObjectID]bool)
	for _, assignment := range assignments {
		validQuestions[assignment.QuestionID] = true
	}

	for _, answer := range cmd.Answers {
		if !validQuestions[answer.QuestionID] {
			return nil, platform_error.NewQuestionNotInSurveyError(answer.QuestionID, cmd.SurveyID)
		}

		question, err := s.questionRepo.GetByID(ctx, answer.QuestionID.Hex())
		if err != nil {
			return nil, platform_error.ErrQuestionNotFound
		}

		options, err := getOptionsFromSpecs(question.Specifications, answer.QuestionID)
		if err != nil {
			return nil, err
		}

		switch question.Format {
		case "multiple_choice":
			if err := validateMultipleChoiceAnswer(answer.Value, options, answer.QuestionID); err != nil {
				return nil, err
			}
		case "likert":
			if err := validateLikertAnswer(answer.Value, options, answer.QuestionID); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid question format: %s", question.Format)
		}
	}

	if err := s.responseRepo.Create(ctx, response); err != nil {
		return nil, platform_error.NewResponseCreationError(err)
	}

	return response, nil
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
		Responses []*resp_model.Response `json:"responses"`
		Total     int64                  `json:"total"`
	}{
		Responses: responses,
		Total:     total,
	}, nil
}

func (s *ResponseService) DeleteResponse(ctx context.Context, cmd *resp_command.DeleteResponseCommand) error {
	return s.responseRepo.Delete(ctx, cmd.ID)
}
