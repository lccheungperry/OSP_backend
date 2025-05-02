package service

import (
	"context"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseService struct {
	repo repository.ResponseRepository
}

func NewResponseService(repo repository.ResponseRepository) *ResponseService {
	return &ResponseService{repo: repo}
}

func (s *ResponseService) SubmitResponse(ctx context.Context, cmd *command.SubmitResponseCommand) (*model.Response, error) {
	answers := make([]model.Answer, len(cmd.Answers))
	for i, cmdAnswer := range cmd.Answers {
		answers[i] = model.Answer{
			QuestionID: cmdAnswer.QuestionID,
			Value:      cmdAnswer.Value,
		}
	}

	response := &model.Response{
		SurveyID:  cmd.SurveyID,
		Answers:   answers,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ResponseService) GetResponses(ctx context.Context, surveyID primitive.ObjectID, page, pageSize int64) ([]*model.Response, int64, error) {
	skip := (page - 1) * pageSize
	return s.repo.FindBySurveyID(ctx, surveyID, skip, pageSize)
}

func (s *ResponseService) DeleteResponse(ctx context.Context, cmd *command.DeleteResponseCommand) error {
	return s.repo.Delete(ctx, cmd.ID)
}
