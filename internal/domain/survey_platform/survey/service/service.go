package service

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	survey_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/command"
	survey_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
)

type SurveyService struct {
	commandBus command.CommandBus
	queryBus   query.QueryBus
}

func NewSurveyService(repo repository.SurveyRepository) SurveyServiceInterface {
	service := &SurveyService{
		commandBus: command.NewInMemoryCommandBus(),
		queryBus:   query.NewInMemoryQueryBus(),
	}

	service.commandBus.RegisterHandler("create_survey", &survey_command.CreateSurveyHandler{Repository: repo})
	service.commandBus.RegisterHandler("update_survey", &survey_command.UpdateSurveyHandler{Repository: repo})
	service.commandBus.RegisterHandler("delete_survey", &survey_command.DeleteSurveyHandler{Repository: repo})

	service.queryBus.RegisterHandler("get_survey", &survey_query.GetSurveyHandler{Repository: repo})
	service.queryBus.RegisterHandler("get_survey_by_token", &survey_query.GetSurveyByTokenHandler{Repository: repo})
	service.queryBus.RegisterHandler("list_surveys", &survey_query.ListSurveysHandler{Repository: repo})

	return service
}

func (s *SurveyService) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	return s.commandBus.Dispatch(ctx, cmd)
}

func (s *SurveyService) HandleQuery(ctx context.Context, query query.Query) (interface{}, error) {
	return s.queryBus.Dispatch(ctx, query)
}
