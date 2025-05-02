package service

import (
	"context"
	"time"

	bus_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	bus_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	resp_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	resp_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseService struct {
	repo repository.ResponseRepository
}

func NewResponseService(repo repository.ResponseRepository) *ResponseService {
	return &ResponseService{repo: repo}
}

func (s *ResponseService) SubmitResponse(ctx context.Context, cmd *resp_command.SubmitResponseCommand) (*model.Response, error) {
	answers := make([]model.Answer, len(cmd.Answers))
	for i, cmdAnswer := range cmd.Answers {
		answers[i] = model.Answer{
			QuestionID: cmdAnswer.QuestionID,
			Value:      cmdAnswer.Value,
		}
	}

	response := &model.Response{
		SurveyID:  cmd.SurveyID,
		Answers:   answers,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *ResponseService) GetResponses(ctx context.Context, surveyID primitive.ObjectID, page, pageSize int64) ([]*model.Response, int64, error) {
	skip := (page - 1) * pageSize
	return s.repo.FindBySurveyID(ctx, surveyID, skip, pageSize)
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
			return nil, err
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
