package command

import (
	"context"
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
		return nil, NewInvalidCommandTypeError(cmd)
	}

	survey := &model.Survey{
		Title: createCmd.Title,
	}

	if err := h.Repository.Create(ctx, survey); err != nil {
		return nil, repository.NewCreateError(err)
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
		return nil, NewInvalidCommandTypeError(cmd)
	}

	existingSurvey, err := h.Repository.FindByID(ctx, updateCmd.ID)
	if err != nil {
		return nil, err
	}

	existingSurvey.Title = updateCmd.Title
	existingSurvey.UpdatedAt = time.Now()

	for i, question := range updateCmd.Questions {
		assignment := &model.SurveyQuestionAssignment{
			SurveyID:   updateCmd.ID,
			QuestionID: question.QuestionID,
			Order:      i + 1,
			CreatedAt:  time.Now(),
		}
		if err := h.Repository.CreateQuestionAssignment(ctx, assignment); err != nil {
			return nil, repository.NewAssignmentError(err)
		}
	}

	if err := h.Repository.Update(ctx, existingSurvey); err != nil {
		return nil, repository.NewUpdateError(err)
	}

	return existingSurvey, nil
}

type DeleteSurveyHandler struct {
	Repository repository.SurveyRepository
}

func (h *DeleteSurveyHandler) HandleCommand(ctx context.Context, cmd command.Command) (interface{}, error) {
	deleteCmd, ok := cmd.(*DeleteSurveyCommand)
	if !ok {
		return nil, NewInvalidCommandTypeError(cmd)
	}

	if err := h.Repository.Delete(ctx, deleteCmd.ID); err != nil {
		return nil, repository.NewDeleteError(err)
	}

	return nil, nil
}
