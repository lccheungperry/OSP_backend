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

// Helper function to convert LikertSpecification to BSON
func convertLikertSpecsToBSON(specs model.LikertSpecification) bson.M {
	options := make([]bson.M, len(specs.Options))
	for i, opt := range specs.Options {
		options[i] = bson.M{
			"label": opt.Label,
			"scale": opt.Scale,
		}
	}
	return bson.M{
		"options": options,
	}
}

// Helper function to convert MultipleChoiceSpecification to BSON
func convertMultipleChoiceSpecsToBSON(specs model.MultipleChoiceSpecification) bson.M {
	options := make([]bson.M, len(specs.Options))
	for i, opt := range specs.Options {
		options[i] = bson.M{
			"label": opt.Label,
		}
	}
	return bson.M{
		"options": options,
	}
}

// Helper function to convert BSON to LikertSpecification
func convertBSONToLikertSpecs(specs bson.M) model.LikertSpecification {
	if options, ok := specs["options"].(primitive.A); ok {
		likertOptions := make([]model.LikertOption, len(options))
		for i, opt := range options {
			if optMap, ok := opt.(bson.M); ok {
				likertOptions[i] = model.LikertOption{
					Label: optMap["label"].(string),
					Scale: int(optMap["scale"].(int32)),
				}
			}
		}
		return model.LikertSpecification{
			Options: likertOptions,
		}
	}
	return model.LikertSpecification{}
}

// Helper function to convert BSON to MultipleChoiceSpecification
func convertBSONToMultipleChoiceSpecs(specs bson.M) model.MultipleChoiceSpecification {
	if options, ok := specs["options"].(primitive.A); ok {
		mcOptions := make([]model.Option, len(options))
		for i, opt := range options {
			if optMap, ok := opt.(bson.M); ok {
				mcOptions[i] = model.Option{
					Label: optMap["label"].(string),
				}
			}
		}
		return model.MultipleChoiceSpecification{
			Options: mcOptions,
		}
	}
	return model.MultipleChoiceSpecification{}
}

func (r *MongoQuestionRepository) Create(ctx context.Context, question *model.Question) error {
	now := time.Now()
	question.CreatedAt = now
	question.UpdatedAt = now

	// Convert specifications to BSON
	switch question.Format {
	case model.FormatLikert:
		if likertSpecs, ok := question.Specifications.(model.LikertSpecification); ok {
			question.Specifications = convertLikertSpecsToBSON(likertSpecs)
		}
	case model.FormatMultipleChoice:
		if mcSpecs, ok := question.Specifications.(model.MultipleChoiceSpecification); ok {
			question.Specifications = convertMultipleChoiceSpecsToBSON(mcSpecs)
		}
	}

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

	// Convert specifications to BSON
	switch question.Format {
	case model.FormatLikert:
		if likertSpecs, ok := question.Specifications.(model.LikertSpecification); ok {
			question.Specifications = convertLikertSpecsToBSON(likertSpecs)
		}
	case model.FormatMultipleChoice:
		if mcSpecs, ok := question.Specifications.(model.MultipleChoiceSpecification); ok {
			question.Specifications = convertMultipleChoiceSpecsToBSON(mcSpecs)
		}
	}

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

	var raw bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&raw)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrQuestionNotFound
		}
		return nil, err
	}

	// Convert raw BSON to Question model
	question := &model.Question{
		ID:        raw["_id"].(primitive.ObjectID),
		Title:     raw["title"].(string),
		Format:    model.QuestionFormat(raw["format"].(string)),
		CreatedAt: raw["created_at"].(primitive.DateTime).Time(),
		UpdatedAt: raw["updated_at"].(primitive.DateTime).Time(),
	}

	// Convert specifications based on format
	if specs, ok := raw["specifications"].(bson.M); ok {
		switch question.Format {
		case model.FormatLikert:
			question.Specifications = convertBSONToLikertSpecs(specs)
		case model.FormatMultipleChoice:
			question.Specifications = convertBSONToMultipleChoiceSpecs(specs)
		}
	}

	return question, nil
}

func (r *MongoQuestionRepository) List(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
	query := bson.M{}
	if filter.Format != "" {
		query["format"] = filter.Format
	}

	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var questions []*model.Question
	if err := cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	for _, question := range questions {
		if specs, ok := question.Specifications.(map[string]interface{}); ok {
			switch question.Format {
			case model.FormatLikert:
				if options, ok := specs["options"].([]interface{}); ok {
					likertOptions := make([]model.LikertOption, len(options))
					for i, opt := range options {
						if optMap, ok := opt.(map[string]interface{}); ok {
							likertOptions[i] = model.LikertOption{
								Label: optMap["label"].(string),
								Scale: int(optMap["scale"].(float64)),
							}
						}
					}
					question.Specifications = model.LikertSpecification{
						Options: likertOptions,
					}
				}
			case model.FormatMultipleChoice:
				if options, ok := specs["options"].([]interface{}); ok {
					mcOptions := make([]model.Option, len(options))
					for i, opt := range options {
						if optMap, ok := opt.(map[string]interface{}); ok {
							mcOptions[i] = model.Option{
								Label: optMap["label"].(string),
							}
						}
					}
					question.Specifications = model.MultipleChoiceSpecification{
						Options: mcOptions,
					}
				}
			}
		}
	}

	return questions, nil
}
