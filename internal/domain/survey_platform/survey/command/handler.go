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

func (h *CreateSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	createCmd, ok := cmd.(*CreateSurveyCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type: %T", cmd)
	}

	survey := &model.Survey{
		Title: createCmd.Title,
	}

	if err := h.Repository.Create(ctx, survey); err != nil {
		return nil, err
	}

	return survey, nil
}

// UpdateSurveyHandler handles the UpdateSurveyCommand
type UpdateSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *UpdateSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	updateCmd, ok := cmd.(*UpdateSurveyCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type: %T", cmd)
	}

	survey := &model.Survey{
		ID:    updateCmd.ID,
		Title: updateCmd.Title,
	}

	if err := h.Repository.Update(ctx, survey); err != nil {
		return nil, err
	}

	return survey, nil
}

// DeleteSurveyHandler handles the DeleteSurveyCommand
type DeleteSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *DeleteSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	deleteCmd, ok := cmd.(*DeleteSurveyCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type: %T", cmd)
	}

	if err := h.Repository.Delete(ctx, deleteCmd.ID); err != nil {
		return nil, err
	}

	return nil, nil
}
