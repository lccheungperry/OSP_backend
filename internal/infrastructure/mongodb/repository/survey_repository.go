package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSurveyRepository struct {
	collection *mongo.Collection
}

func NewMongoSurveyRepository(db *mongo.Database) repository.SurveyRepository {
	return &MongoSurveyRepository{
		collection: db.Collection("surveys"),
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:5], nil
}

func (r *MongoSurveyRepository) Create(ctx context.Context, survey *model.Survey) error {
	now := time.Now()
	survey.CreatedAt = now
	survey.UpdatedAt = now

	token, err := generateToken()
	if err != nil {
		return err
	}
	survey.Token = token

	result, err := r.collection.InsertOne(ctx, survey)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		survey.ID = oid
	}

	return nil
}

func (r *MongoSurveyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Survey, error) {
	var survey model.Survey
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&survey)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrSurveyNotFound
		}
		return nil, err
	}
	return &survey, nil
}

func (r *MongoSurveyRepository) FindByToken(ctx context.Context, token string) (*model.Survey, error) {
	var survey model.Survey
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&survey)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrSurveyNotFound
		}
		return nil, err
	}
	return &survey, nil
}

func (r *MongoSurveyRepository) List(ctx context.Context, skip, limit int64) ([]*model.Survey, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var surveys []*model.Survey
	if err := cursor.All(ctx, &surveys); err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
}

func (r *MongoSurveyRepository) Update(ctx context.Context, survey *model.Survey) error {
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": survey.ID}, survey)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return repository.ErrSurveyNotFound
	}

	return nil
}

func (r *MongoSurveyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrSurveyNotFound
	}

	return nil
}
