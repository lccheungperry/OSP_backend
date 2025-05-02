package repository

import (
	"context"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	surveyRepo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	mongoRepo "github.com/lccheungperry/OSP_backend/internal/infrastructure/mongodb/repository"
	"github.com/lccheungperry/OSP_backend/tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoSurveyRepository(t *testing.T) {
	db, cleanup := utils.SetupTestDB(t, "test_survey_repo")
	defer cleanup()

	repo := mongoRepo.NewMongoSurveyRepository(db)

	t.Run("Create and FindByID", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		survey := &model.Survey{
			Title: "Test Survey",
		}

		err := repo.Create(ctx, survey)
		require.NoError(t, err)
		assert.NotEmpty(t, survey.ID)
		assert.NotEmpty(t, survey.Token)
		assert.NotZero(t, survey.CreatedAt)

		retrieved, err := repo.FindByID(ctx, survey.ID)
		require.NoError(t, err)
		assert.Equal(t, survey.ID, retrieved.ID)
		assert.Equal(t, survey.Title, retrieved.Title)
		assert.Equal(t, survey.Token, retrieved.Token)
	})

	t.Run("FindByToken", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		survey := &model.Survey{
			Title: "Test Survey",
		}

		err := repo.Create(ctx, survey)
		require.NoError(t, err)

		retrieved, err := repo.FindByToken(ctx, survey.Token)
		require.NoError(t, err)
		assert.Equal(t, survey.ID, retrieved.ID)
		assert.Equal(t, survey.Title, retrieved.Title)
		assert.Equal(t, survey.Token, retrieved.Token)
	})

	t.Run("Update", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		survey := &model.Survey{
			Title: "Original Title",
		}

		err := repo.Create(ctx, survey)
		require.NoError(t, err)

		survey.Title = "Updated Title"
		err = repo.Update(ctx, survey)
		require.NoError(t, err)

		retrieved, err := repo.FindByID(ctx, survey.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", retrieved.Title)
	})

	t.Run("Delete", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		survey := &model.Survey{
			Title: "To be deleted",
		}

		err := repo.Create(ctx, survey)
		require.NoError(t, err)

		err = repo.Delete(ctx, survey.ID)
		require.NoError(t, err)

		_, err = repo.FindByID(ctx, survey.ID)
		assert.ErrorIs(t, err, surveyRepo.ErrSurveyNotFound)
	})

	t.Run("List", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		surveys := []*model.Survey{
			{Title: "Survey 1"},
			{Title: "Survey 2"},
			{Title: "Survey 3"},
		}

		for _, s := range surveys {
			err := repo.Create(ctx, s)
			require.NoError(t, err)
		}

		results, total, err := repo.List(ctx, 0, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 2)

		results, total, err = repo.List(ctx, 2, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 1)
	})

	t.Run("FindByID non-existent", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		_, err := repo.FindByID(ctx, primitive.NewObjectID())
		assert.ErrorIs(t, err, surveyRepo.ErrSurveyNotFound)
	})

	t.Run("FindByToken non-existent", func(t *testing.T) {
		utils.CleanCollection(t, db, "surveys")
		ctx := context.Background()

		_, err := repo.FindByToken(ctx, "non-existent-token")
		assert.ErrorIs(t, err, surveyRepo.ErrSurveyNotFound)
	})
}
