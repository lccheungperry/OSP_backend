package command

import (
	"context"
	"fmt"
	"time"

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

	// Get existing survey
	existingSurvey, err := h.Repository.FindByID(ctx, updateCmd.ID)
	if err != nil {
		return nil, err
	}

	// Update only the fields that need to be changed
	existingSurvey.Title = updateCmd.Title
	existingSurvey.UpdatedAt = time.Now()

	if err := h.Repository.Update(ctx, existingSurvey); err != nil {
		return nil, err
	}

	return existingSurvey, nil
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
