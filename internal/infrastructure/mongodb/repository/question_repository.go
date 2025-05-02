package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQuestionRepository struct {
	collection *mongo.Collection
}

func NewMongoQuestionRepository(db *mongo.Database) repository.QuestionRepository {
	return &MongoQuestionRepository{
		collection: db.Collection("questions"),
	}
}

func (r *MongoQuestionRepository) Create(ctx context.Context, question *model.Question) error {
	now := time.Now()
	question.CreatedAt = now
	question.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, question)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		question.ID = oid
	}

	return nil
}

func (r *MongoQuestionRepository) Update(ctx context.Context, id string, question *model.Question) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	question.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, question)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return repository.ErrQuestionNotFound
	}

	return nil
}

func (r *MongoQuestionRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrQuestionNotFound
	}

	return nil
}

func (r *MongoQuestionRepository) GetByID(ctx context.Context, id string) (*model.Question, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var question model.Question
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrQuestionNotFound
		}
		return nil, err
	}

	return &question, nil
}

func (r *MongoQuestionRepository) List(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
	query := bson.M{}
	if filter.Format != "" {
		query["format"] = filter.Format
	}

	// Set up pagination options
	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var questions []*model.Question
	if err := cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}
