package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidCommand = errors.New("invalid command")

type CommandHandler interface {
	HandleCommand(ctx context.Context, cmd interface{}) error
}

type CreateQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *CreateQuestionHandler) HandleCommand(ctx context.Context, cmd command.Command) error {
	createCmd, ok := cmd.(*CreateQuestionCommand)
	if !ok {
		return fmt.Errorf("invalid command type: %T", cmd)
	}

	question := &model.Question{
		Title:          createCmd.Title,
		Format:         model.QuestionFormat(createCmd.Format),
		Specifications: createCmd.Specifications,
	}

	return h.Repository.Create(ctx, question)
}

type UpdateQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *UpdateQuestionHandler) HandleCommand(ctx context.Context, cmd command.Command) error {
	updateCmd, ok := cmd.(*UpdateQuestionCommand)
	if !ok {
		return fmt.Errorf("invalid command type: %T", cmd)
	}

	id, err := primitive.ObjectIDFromHex(updateCmd.ID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	question := &model.Question{
		ID:             id,
		Title:          updateCmd.Title,
		Format:         model.QuestionFormat(updateCmd.Format),
		Specifications: updateCmd.Specifications,
	}

	return h.Repository.Update(ctx, updateCmd.ID, question)
}

type DeleteQuestionHandler struct {
	Repository repository.QuestionRepository
}

func (h *DeleteQuestionHandler) HandleCommand(ctx context.Context, cmd command.Command) error {
	deleteCmd, ok := cmd.(*DeleteQuestionCommand)
	if !ok {
		return fmt.Errorf("invalid command type: %T", cmd)
	}

	return h.Repository.Delete(ctx, deleteCmd.ID)
}
