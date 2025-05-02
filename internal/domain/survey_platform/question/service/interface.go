package service

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
)

type QuestionServiceInterface interface {
	HandleCommand(ctx context.Context, cmd command.Command) error
	HandleQuery(ctx context.Context, q query.Query) (interface{}, error)
}
