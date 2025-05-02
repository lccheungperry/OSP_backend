package command

import (
	"context"
	"fmt"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
)

type CreateSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *CreateSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) error {
	createCmd, ok := cmd.(*CreateSurveyCommand)
	if !ok {
		return fmt.Errorf("invalid command type: %T", cmd)
	}

	survey := &model.Survey{
		Title: createCmd.Title,
	}

	return h.Repository.Create(ctx, survey)
}

// UpdateSurveyHandler handles the UpdateSurveyCommand
type UpdateSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *UpdateSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) error {
	updateCmd, ok := cmd.(*UpdateSurveyCommand)
	if !ok {
		return fmt.Errorf("invalid command type: %T", cmd)
	}

	survey := &model.Survey{
		ID:    updateCmd.ID,
		Title: updateCmd.Title,
	}

	return h.Repository.Update(ctx, survey)
}

// DeleteSurveyHandler handles the DeleteSurveyCommand
type DeleteSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *DeleteSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) error {
	deleteCmd, ok := cmd.(*DeleteSurveyCommand)
	if !ok {
		return fmt.Errorf("invalid command type: %T", cmd)
	}

	return h.Repository.Delete(ctx, deleteCmd.ID)
}
