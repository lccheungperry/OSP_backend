package question_test

import (
	"context"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/tests/domain/survey_platform/question"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateQuestionHandler(t *testing.T) {
	repo := question.NewMockQuestionRepository()
	handler := &command.CreateQuestionHandler{Repository: repo}

	cmd := &command.CreateQuestionCommand{
		Title:          "Test Question",
		Format:         "textbox",
		Specifications: map[string]interface{}{},
	}

	result, err := handler.HandleCommand(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	question, ok := result.(*model.Question)
	assert.True(t, ok)
	assert.Equal(t, "Test Question", question.Title)
	assert.Equal(t, model.QuestionFormat("textbox"), question.Format)
}

func TestUpdateQuestionHandler(t *testing.T) {
	repo := question.NewMockQuestionRepository()
	handler := &command.UpdateQuestionHandler{Repository: repo}

	id := primitive.NewObjectID()
	cmd := &command.UpdateQuestionCommand{
		ID:             id.Hex(),
		Title:          "Updated Question",
		Format:         "multiple_choice",
		Specifications: map[string]interface{}{},
	}

	result, err := handler.HandleCommand(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	question, ok := result.(*model.Question)
	assert.True(t, ok)
	assert.Equal(t, id, question.ID)
	assert.Equal(t, "Updated Question", question.Title)
	assert.Equal(t, model.QuestionFormat("multiple_choice"), question.Format)
}

func TestDeleteQuestionHandler(t *testing.T) {
	repo := question.NewMockQuestionRepository()
	handler := &command.DeleteQuestionHandler{Repository: repo}

	id := primitive.NewObjectID()
	cmd := &command.DeleteQuestionCommand{
		ID: id.Hex(),
	}

	result, err := handler.HandleCommand(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Nil(t, result)
}
