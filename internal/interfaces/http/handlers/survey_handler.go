package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/service"
)

// SurveyHandler handles HTTP requests for survey operations
type SurveyHandler struct {
	service *service.SurveyService
}

// NewSurveyHandler creates a new survey handler
func NewSurveyHandler(service *service.SurveyService) *SurveyHandler {
	return &SurveyHandler{service: service}
}

// CreateSurvey handles the creation of a new survey
func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	var req model.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	survey, err := h.service.CreateSurvey(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create survey",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, survey)
}

// ListSurveys handles listing all surveys
func (h *SurveyHandler) ListSurveys(c *gin.Context) {
	// TODO: Implement pagination parameters
	surveys, total, err := h.service.ListSurveys(c.Request.Context(), 1, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list surveys",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"surveys": surveys,
		"total":   total,
	})
}

// GetSurvey handles getting a survey by ID
func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	id := c.Param("id")
	survey, err := h.service.GetSurvey(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Survey not found",
				"details": "Survey with ID '" + id + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get survey",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, survey)
}

// UpdateSurvey handles updating a survey
func (h *SurveyHandler) UpdateSurvey(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	survey, err := h.service.UpdateSurvey(c.Request.Context(), id, &req)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Survey not found",
				"details": "Survey with ID '" + id + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update survey",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, survey)
}

// DeleteSurvey handles deleting a survey
func (h *SurveyHandler) DeleteSurvey(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSurvey(c.Request.Context(), id); err != nil {
		if err == repository.ErrSurveyNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Survey not found",
				"details": "Survey with ID '" + id + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete survey",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSurveyByToken handles getting a survey by token
func (h *SurveyHandler) GetSurveyByToken(c *gin.Context) {
	token := c.Param("token")
	survey, err := h.service.GetSurveyByToken(c.Request.Context(), token)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Survey not found",
				"details": "Survey with token '" + token + "' not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get survey",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, survey)
}
