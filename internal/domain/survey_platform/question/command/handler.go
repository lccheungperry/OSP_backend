package command

import (
	"context"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	platform_error "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/error"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidCommand = platform_error.NewInvalidCommandError()

type CommandHandler interface {
	HandleCommand(ctx context.Context, cmd interface{}) (interface{}, error)
}

type CreateQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *CreateQuestionHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	createCmd, ok := cmd.(*CreateQuestionCommand)
	if !ok {
		return nil, platform_error.NewInvalidCommandTypeError(cmd)
	}

	question := &model.Question{
		Title:          createCmd.Title,
		Format:         model.QuestionFormat(createCmd.Format),
		Specifications: createCmd.Specifications,
	}

	if err := h.Repository.Create(ctx, question); err != nil {
		return nil, platform_error.NewCreateError(err)
	}

	return question, nil
}

type UpdateQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *UpdateQuestionHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	updateCmd, ok := cmd.(*UpdateQuestionCommand)
	if !ok {
		return nil, platform_error.NewInvalidCommandTypeError(cmd)
	}

	id, err := primitive.ObjectIDFromHex(updateCmd.ID)
	if err != nil {
		return nil, platform_error.NewInvalidQuestionIDError(err)
	}

	question := &model.Question{
		ID:             id,
		Title:          updateCmd.Title,
		Format:         model.QuestionFormat(updateCmd.Format),
		Specifications: updateCmd.Specifications,
	}

	if err := h.Repository.Update(ctx, updateCmd.ID, question); err != nil {
		return nil, platform_error.NewUpdateError(err)
	}

	return question, nil
}

type DeleteQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *DeleteQuestionHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	deleteCmd, ok := cmd.(*DeleteQuestionCommand)
	if !ok {
		return nil, platform_error.NewInvalidCommandTypeError(cmd)
	}

	if err := h.Repository.Delete(ctx, deleteCmd.ID); err != nil {
		return nil, platform_error.NewDeleteError(err)
	}

	return nil, nil
}
