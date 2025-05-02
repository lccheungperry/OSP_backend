package question

import (
	"context"
	"errors"

	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrQuestionNotFound = errors.New("question not found")

type MockQuestionRepository struct {
	CreateFunc  func(ctx context.Context, question *model.Question) error
	UpdateFunc  func(ctx context.Context, id string, question *model.Question) error
	DeleteFunc  func(ctx context.Context, id string) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Question, error)
	ListFunc    func(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error)
}

func NewMockQuestionRepository() *MockQuestionRepository {
	return &MockQuestionRepository{
		CreateFunc: func(ctx context.Context, question *model.Question) error {
			question.ID = primitive.NewObjectID()
			return nil
		},
		UpdateFunc: func(ctx context.Context, id string, question *model.Question) error {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return err
			}
			question.ID = objID
			return nil
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
		GetByIDFunc: func(ctx context.Context, id string) (*model.Question, error) {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, err
			}
			return &model.Question{
				ID:             objID,
				Title:          "Test Question",
				Format:         model.QuestionFormat("textbox"),
				Specifications: map[string]interface{}{},
			}, nil
		},
		ListFunc: func(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
			return []*model.Question{}, nil
		},
	}
}

func (r *MockQuestionRepository) Create(ctx context.Context, question *model.Question) error {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, question)
	}
	return errors.New("CreateFunc not implemented")
}

func (r *MockQuestionRepository) Update(ctx context.Context, id string, question *model.Question) error {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, question)
	}
	return errors.New("UpdateFunc not implemented")
}

func (r *MockQuestionRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return errors.New("DeleteFunc not implemented")
}

func (r *MockQuestionRepository) GetByID(ctx context.Context, id string) (*model.Question, error) {
	if r.GetByIDFunc != nil {
		return r.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("GetByIDFunc not implemented")
}

func (r *MockQuestionRepository) List(ctx context.Context, filter *model.QuestionFilter) ([]*model.Question, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx, filter)
	}
	return nil, errors.New("ListFunc not implemented")
}
