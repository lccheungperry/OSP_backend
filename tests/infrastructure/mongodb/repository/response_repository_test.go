package repository

import (
	"context"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	responseRepo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/repository"
	mongoRepo "github.com/lccheungperry/OSP_backend/internal/infrastructure/mongodb/repository"
	"github.com/lccheungperry/OSP_backend/tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoResponseRepository(t *testing.T) {
	db, cleanup := utils.SetupTestDB(t, "test_response_repo")
	defer cleanup()

	repo := mongoRepo.NewMongoResponseRepository(db)

	t.Run("Create and FindByID", func(t *testing.T) {
		utils.CleanCollection(t, db, "responses")
		ctx := context.Background()

		response := &model.Response{
			SurveyID: primitive.NewObjectID(),
			Answers: []model.Answer{
				{
					QuestionID: primitive.NewObjectID(),
					Value:      "Test answer",
				},
			},
		}

		err := repo.Create(ctx, response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.NotZero(t, response.CreatedAt)

		retrieved, err := repo.FindByID(ctx, response.ID)
		require.NoError(t, err)
		assert.Equal(t, response.ID, retrieved.ID)
		assert.Equal(t, response.SurveyID, retrieved.SurveyID)
		assert.Equal(t, response.Answers[0].Value, retrieved.Answers[0].Value)
	})

	t.Run("FindBySurveyID", func(t *testing.T) {
		utils.CleanCollection(t, db, "responses")
		ctx := context.Background()

		surveyID := primitive.NewObjectID()
		responses := []*model.Response{
			{
				SurveyID: surveyID,
				Answers: []model.Answer{
					{
						QuestionID: primitive.NewObjectID(),
						Value:      "Answer 1",
					},
				},
			},
			{
				SurveyID: surveyID,
				Answers: []model.Answer{
					{
						QuestionID: primitive.NewObjectID(),
						Value:      "Answer 2",
					},
				},
			},
		}

		for _, r := range responses {
			err := repo.Create(ctx, r)
			require.NoError(t, err)
		}

		results, total, err := repo.FindBySurveyID(ctx, surveyID, 0, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)

		results, total, err = repo.FindBySurveyID(ctx, surveyID, 2, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 0)
	})

	t.Run("Delete", func(t *testing.T) {
		utils.CleanCollection(t, db, "responses")
		ctx := context.Background()

		response := &model.Response{
			SurveyID: primitive.NewObjectID(),
			Answers: []model.Answer{
				{
					QuestionID: primitive.NewObjectID(),
					Value:      "To be deleted",
				},
			},
		}

		err := repo.Create(ctx, response)
		require.NoError(t, err)

		err = repo.Delete(ctx, response.ID)
		require.NoError(t, err)

		_, err = repo.FindByID(ctx, response.ID)
		assert.ErrorIs(t, err, responseRepo.ErrResponseNotFound)
	})

	t.Run("FindByID non-existent", func(t *testing.T) {
		utils.CleanCollection(t, db, "responses")
		ctx := context.Background()

		_, err := repo.FindByID(ctx, primitive.NewObjectID())
		assert.ErrorIs(t, err, responseRepo.ErrResponseNotFound)
	})

	t.Run("FindBySurveyID non-existent", func(t *testing.T) {
		utils.CleanCollection(t, db, "responses")
		ctx := context.Background()

		results, total, err := repo.FindBySurveyID(ctx, primitive.NewObjectID(), 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, results)
	})
}
