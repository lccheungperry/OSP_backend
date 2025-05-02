package service

import (
	"context"

	buscommand "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	busquery "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	questioncommand "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/command"
	questionquery "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
)

type QuestionService struct {
	commandBus buscommand.CommandBus
	queryBus   busquery.QueryBus
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	service := &QuestionService{
		commandBus: buscommand.NewInMemoryCommandBus(),
		queryBus:   busquery.NewInMemoryQueryBus(),
	}

	service.commandBus.RegisterHandler("create_question", &questioncommand.CreateQuestionHandler{Repository: repo})
	service.commandBus.RegisterHandler("update_question", &questioncommand.UpdateQuestionHandler{Repository: repo})
	service.commandBus.RegisterHandler("delete_question", &questioncommand.DeleteQuestionHandler{Repository: repo})

	service.queryBus.RegisterHandler("get_question", &questionquery.GetQuestionHandler{Repository: repo})
	service.queryBus.RegisterHandler("list_questions", &questionquery.ListQuestionsHandler{Repository: repo})

	return service
}

func (s *QuestionService) HandleCommand(ctx context.Context, cmd buscommand.Command) error {
	return s.commandBus.Dispatch(ctx, cmd)
}

func (s *QuestionService) HandleQuery(ctx context.Context, query busquery.Query) (interface{}, error) {
	return s.queryBus.Dispatch(ctx, query)
}
