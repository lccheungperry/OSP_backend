package repository

import (
	"context"
	"testing"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	questionRepo "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	mongoRepo "github.com/lccheungperry/OSP_backend/internal/infrastructure/mongodb/repository"
	"github.com/lccheungperry/OSP_backend/tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoQuestionRepository(t *testing.T) {
	db, cleanup := utils.SetupTestDB(t, "test_question_repo")
	defer cleanup()

	repo := mongoRepo.NewMongoQuestionRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		question := &model.Question{
			Title:  "Test Question",
			Format: model.FormatTextbox,
			Specifications: model.MultipleChoiceSpecification{
				Options: []model.Option{
					{Label: "Option 1"},
					{Label: "Option 2"},
				},
			},
		}

		err := repo.Create(ctx, question)
		require.NoError(t, err)
		assert.NotEmpty(t, question.ID)
		assert.NotZero(t, question.CreatedAt)
		assert.NotZero(t, question.UpdatedAt)

		retrieved, err := repo.GetByID(ctx, question.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, question.ID, retrieved.ID)
		assert.Equal(t, question.Title, retrieved.Title)
		assert.Equal(t, question.Format, retrieved.Format)
	})

	t.Run("Update", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		question := &model.Question{
			Title:  "Original Title",
			Format: model.FormatTextbox,
			Specifications: model.MultipleChoiceSpecification{
				Options: []model.Option{
					{Label: "Option 1"},
					{Label: "Option 2"},
				},
			},
		}

		err := repo.Create(ctx, question)
		require.NoError(t, err)

		question.Title = "Updated Title"
		err = repo.Update(ctx, question.ID.Hex(), question)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, question.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", retrieved.Title)
	})

	t.Run("Delete", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		question := &model.Question{
			Title:  "To be deleted",
			Format: model.FormatTextbox,
			Specifications: model.MultipleChoiceSpecification{
				Options: []model.Option{
					{Label: "Option 1"},
					{Label: "Option 2"},
				},
			},
		}

		err := repo.Create(ctx, question)
		require.NoError(t, err)

		err = repo.Delete(ctx, question.ID.Hex())
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, question.ID.Hex())
		assert.ErrorIs(t, err, questionRepo.ErrQuestionNotFound)
	})

	t.Run("List with filter", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		questions := []*model.Question{
			{
				Title:  "Textbox Question",
				Format: model.FormatTextbox,
				Specifications: model.MultipleChoiceSpecification{
					Options: []model.Option{
						{Label: "Option 1"},
						{Label: "Option 2"},
					},
				},
			},
			{
				Title:  "Multiple Choice Question",
				Format: model.FormatMultipleChoice,
				Specifications: model.MultipleChoiceSpecification{
					Options: []model.Option{
						{Label: "Option 1"},
						{Label: "Option 2"},
					},
				},
			},
			{
				Title:  "Likert Question",
				Format: model.FormatLikert,
				Specifications: model.LikertSpecification{
					Options: []model.LikertOption{
						{Label: "Strongly Disagree", Scale: 1},
						{Label: "Strongly Agree", Scale: 5},
					},
				},
			},
		}

		for _, q := range questions {
			err := repo.Create(ctx, q)
			require.NoError(t, err)
		}

		filter := &model.QuestionFilter{
			Format: model.FormatMultipleChoice,
		}
		results, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, model.FormatMultipleChoice, results[0].Format)

		filter = &model.QuestionFilter{
			Limit:  2,
			Offset: 0,
		}
		results, err = repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("GetByID with invalid ID", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		_, err := repo.GetByID(ctx, "invalid_id")
		assert.Error(t, err)
	})

	t.Run("Update non-existent question", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		nonExistentID := primitive.NewObjectID().Hex()
		question := &model.Question{
			Title:  "Non-existent",
			Format: model.FormatTextbox,
			Specifications: model.MultipleChoiceSpecification{
				Options: []model.Option{
					{Label: "Option 1"},
					{Label: "Option 2"},
				},
			},
		}

		err := repo.Update(ctx, nonExistentID, question)
		assert.ErrorIs(t, err, questionRepo.ErrQuestionNotFound)
	})

	t.Run("Delete non-existent question", func(t *testing.T) {
		utils.CleanCollection(t, db, "questions")
		ctx := context.Background()

		nonExistentID := primitive.NewObjectID().Hex()
		err := repo.Delete(ctx, nonExistentID)
		assert.ErrorIs(t, err, questionRepo.ErrQuestionNotFound)
	})
}
