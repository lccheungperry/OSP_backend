package service

import (
	"context"

	buscommand "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	busquery "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	surveycommand "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/command"
	surveyquery "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
)

type SurveyService struct {
	commandBus buscommand.CommandBus
	queryBus   busquery.QueryBus
}

func NewSurveyService(repo repository.SurveyRepository) SurveyServiceInterface {
	service := &SurveyService{
		commandBus: buscommand.NewInMemoryCommandBus(),
		queryBus:   busquery.NewInMemoryQueryBus(),
	}

	service.commandBus.RegisterHandler("create_survey", &surveycommand.CreateSurveyHandler{Repository: repo})
	service.commandBus.RegisterHandler("update_survey", &surveycommand.UpdateSurveyHandler{Repository: repo})
	service.commandBus.RegisterHandler("delete_survey", &surveycommand.DeleteSurveyHandler{Repository: repo})

	service.queryBus.RegisterHandler("get_survey", &surveyquery.GetSurveyHandler{Repository: repo})
	service.queryBus.RegisterHandler("get_survey_by_token", &surveyquery.GetSurveyByTokenHandler{Repository: repo})
	service.queryBus.RegisterHandler("list_surveys", &surveyquery.ListSurveysHandler{Repository: repo})

	return service
}

func (s *SurveyService) HandleCommand(ctx context.Context, cmd buscommand.Command) error {
	return s.commandBus.Dispatch(ctx, cmd)
}

func (s *SurveyService) HandleQuery(ctx context.Context, query busquery.Query) (interface{}, error) {
	return s.queryBus.Dispatch(ctx, query)
}
