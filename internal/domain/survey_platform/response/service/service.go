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
)

type ResponseService struct {
	repo repository.ResponseRepository
}

func NewResponseService(repo repository.ResponseRepository) *ResponseService {
	return &ResponseService{repo: repo}
}

func (s *ResponseService) HandleCommand(ctx context.Context, cmd bus_command.Command) error {
	switch c := cmd.(type) {
	case *resp_command.SubmitResponseCommand:
		answers := make([]model.Answer, len(c.Answers))
		for i, cmdAnswer := range c.Answers {
			answers[i] = model.Answer{
				QuestionID: cmdAnswer.QuestionID,
				Value:      cmdAnswer.Value,
			}
		}

		response := &model.Response{
			SurveyID:  c.SurveyID,
			Answers:   answers,
			CreatedAt: time.Now(),
		}

		return s.repo.Create(ctx, response)
	case *resp_command.DeleteResponseCommand:
		return s.repo.Delete(ctx, c.ID)
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
			Responses []*model.Response `json:"responses"`
			Total     int64             `json:"total"`
		}{
			Responses: responses,
			Total:     total,
		}, nil
	default:
		return nil, nil
	}
}
