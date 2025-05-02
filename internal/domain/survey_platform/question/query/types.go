package query

import "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"

type GetQuestionQuery struct {
	ID string `json:"id" validate:"required"`
}

func (q *GetQuestionQuery) QueryName() string {
	return "get_question"
}

type ListQuestionsQuery struct {
	Filter *model.QuestionFilter `json:"filter"`
}

func (q *ListQuestionsQuery) QueryName() string {
	return "list_questions"
}
