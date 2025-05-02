package command

type CreateQuestionCommand struct {
	Title          string      `json:"title" validate:"required,min=3,max=500"`
	Format         string      `json:"format" validate:"required,oneof=textbox multiple_choice likert"`
	Specifications interface{} `json:"specifications" validate:"required"`
}

func (c *CreateQuestionCommand) CommandName() string {
	return "create_question"
}

type UpdateQuestionCommand struct {
	ID             string      `json:"id" validate:"required"`
	Title          string      `json:"title" validate:"required,min=3,max=500"`
	Format         string      `json:"format" validate:"required,oneof=textbox multiple_choice likert"`
	Specifications interface{} `json:"specifications" validate:"required"`
}

func (c *UpdateQuestionCommand) CommandName() string {
	return "update_question"
}

type DeleteQuestionCommand struct {
	ID string `json:"id" validate:"required"`
}

func (c *DeleteQuestionCommand) CommandName() string {
	return "delete_question"
}
