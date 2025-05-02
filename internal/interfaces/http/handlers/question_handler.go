package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/repository"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/service"
)

type QuestionHandler struct {
	service service.QuestionServiceInterface
}

func NewQuestionHandler(service service.QuestionServiceInterface) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req struct {
		Title          string      `json:"title" binding:"required,min=3,max=500"`
		Format         string      `json:"format" binding:"required,oneof=textbox multiple_choice likert"`
		Specifications interface{} `json:"specifications" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	var specifications interface{}
	switch req.Format {
	case "likert":
		var likertSpecs struct {
			Options []struct {
				Label string `json:"label"`
				Scale int    `json:"scale"`
			} `json:"options"`
		}
		if specBytes, err := json.Marshal(req.Specifications); err == nil {
			if err := json.Unmarshal(specBytes, &likertSpecs); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid Likert specifications",
					"details": err.Error(),
				})
				return
			}
			options := make([]model.LikertOption, len(likertSpecs.Options))
			for i, opt := range likertSpecs.Options {
				options[i] = model.LikertOption{
					Label: opt.Label,
					Scale: opt.Scale,
				}
			}
			specifications = model.LikertSpecification{
				Options: options,
			}
		}
	default:
		specifications = req.Specifications
	}

	cmd := &command.CreateQuestionCommand{
		Title:          req.Title,
		Format:         req.Format,
		Specifications: specifications,
	}

	result, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create question",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *QuestionHandler) ListQuestions(c *gin.Context) {
	format := c.Query("format")
	var questionFormat model.QuestionFormat
	if format != "" {
		questionFormat = model.QuestionFormat(format)
	}

	q := &query.ListQuestionsQuery{
		Filter: &model.QuestionFilter{
			Format: questionFormat,
			Limit:  10,
			Offset: 0,
		},
	}

	result, err := h.service.HandleQuery(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list questions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	id := c.Param("id")
	q := &query.GetQuestionQuery{
		ID: id,
	}

	result, err := h.service.HandleQuery(c.Request.Context(), q)
	if err != nil {
		if err == repository.ErrQuestionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Question not found",
				"details": "Question with ID '" + id + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get question",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title          string      `json:"title" binding:"required,min=3,max=500"`
		Format         string      `json:"format" binding:"required,oneof=textbox multiple_choice likert"`
		Specifications interface{} `json:"specifications" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	var specifications interface{}
	switch req.Format {
	case "likert":
		var likertSpecs struct {
			Options []struct {
				Label string `json:"label"`
				Scale int    `json:"scale"`
			} `json:"options"`
		}
		if specBytes, err := json.Marshal(req.Specifications); err == nil {
			if err := json.Unmarshal(specBytes, &likertSpecs); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid Likert specifications",
					"details": err.Error(),
				})
				return
			}
			options := make([]model.LikertOption, len(likertSpecs.Options))
			for i, opt := range likertSpecs.Options {
				options[i] = model.LikertOption{
					Label: opt.Label,
					Scale: opt.Scale,
				}
			}
			specifications = model.LikertSpecification{
				Options: options,
			}
		}
	default:
		specifications = req.Specifications
	}

	cmd := &command.UpdateQuestionCommand{
		ID:             id,
		Title:          req.Title,
		Format:         req.Format,
		Specifications: specifications,
	}

	_, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		if err == repository.ErrQuestionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Question not found",
				"details": "Question with ID '" + id + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update question",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	cmd := &command.DeleteQuestionCommand{
		ID: id,
	}

	_, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		if err == repository.ErrQuestionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Question not found",
				"details": "Question with ID '" + id + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete question",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
