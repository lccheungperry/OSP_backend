package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoResponseRepository struct {
	collection *mongo.Collection
}

func NewMongoResponseRepository(db *mongo.Database) repository.ResponseRepository {
	return &MongoResponseRepository{
		collection: db.Collection("responses"),
	}
}

func (r *MongoResponseRepository) Create(ctx context.Context, response *model.Response) error {
	response.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, response)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		response.ID = oid
	}

	return nil
}

func (r *MongoResponseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Response, error) {
	var response model.Response
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&response)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrResponseNotFound
		}
		return nil, err
	}
	return &response, nil
}

func (r *MongoResponseRepository) FindBySurveyID(ctx context.Context, surveyID primitive.ObjectID, skip, limit int64) ([]*model.Response, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{"survey_id": surveyID})
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"survey_id": surveyID}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var responses []*model.Response
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, 0, err
	}

	return responses, total, nil
}

func (r *MongoResponseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrResponseNotFound
	}

	return nil
}
