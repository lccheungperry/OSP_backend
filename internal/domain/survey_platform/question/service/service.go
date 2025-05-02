package service

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	question_command "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/command"
	question_query "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
)

type QuestionService struct {
	commandBus command.CommandBus
	queryBus   query.QueryBus
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	service := &QuestionService{
		commandBus: command.NewInMemoryCommandBus(),
		queryBus:   query.NewInMemoryQueryBus(),
	}

	service.commandBus.RegisterHandler("create_question", &question_command.CreateQuestionHandler{Repository: repo})
	service.commandBus.RegisterHandler("update_question", &question_command.UpdateQuestionHandler{Repository: repo})
	service.commandBus.RegisterHandler("delete_question", &question_command.DeleteQuestionHandler{Repository: repo})

	service.queryBus.RegisterHandler("get_question", &question_query.GetQuestionHandler{Repository: repo})
	service.queryBus.RegisterHandler("list_questions", &question_query.ListQuestionsHandler{Repository: repo})

	return service
}

func (s *QuestionService) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	return s.commandBus.Dispatch(ctx, cmd)
}

func (s *QuestionService) HandleQuery(ctx context.Context, query query.Query) (interface{}, error) {
	return s.queryBus.Dispatch(ctx, query)
}
