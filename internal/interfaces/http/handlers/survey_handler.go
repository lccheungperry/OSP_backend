package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SurveyHandler struct {
	service service.SurveyServiceInterface
}

func NewSurveyHandler(service service.SurveyServiceInterface) *SurveyHandler {
	return &SurveyHandler{service: service}
}

func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required,min=3,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if !ValidateSurveyRequest(c, req.Title, nil) {
		return
	}

	cmd := &command.CreateSurveyCommand{
		Title: req.Title,
	}

	result, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	survey, ok := result.(*model.Survey)
	if !ok {
		SendErrorResponse(c, http.StatusInternalServerError, "INVALID_RESPONSE", "Invalid response type")
		return
	}

	SendSuccessResponse(c, http.StatusCreated, survey)
}

func (h *SurveyHandler) ListSurveys(c *gin.Context) {
	q := &query.ListSurveysQuery{
		Filter: &model.SurveyFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	result, err := h.service.HandleQuery(c.Request.Context(), q)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	SendSuccessResponse(c, http.StatusOK, result)
}

func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	id, ok := ValidateObjectID(c, c.Param("id"))
	if !ok {
		return
	}

	q := &query.GetSurveyQuery{
		ID: id,
	}

	result, err := h.service.HandleQuery(c.Request.Context(), q)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			SendErrorResponse(c, http.StatusNotFound, "SURVEY_NOT_FOUND",
				"Survey with ID '"+id+"' not found")
			return
		}
		HandleServiceError(c, err)
		return
	}

	SendSuccessResponse(c, http.StatusOK, result)
}

func (h *SurveyHandler) UpdateSurvey(c *gin.Context) {
	id, ok := ValidateObjectID(c, c.Param("id"))
	if !ok {
		return
	}

	var req struct {
		Title     string `json:"title" binding:"required,min=3,max=100"`
		Questions []struct {
			QuestionID string `json:"questionId" binding:"required"`
			Order      int    `json:"order" binding:"required,min=1"`
		} `json:"questions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if !ValidateSurveyRequest(c, req.Title, req.Questions) {
		return
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	cmd := &command.UpdateSurveyCommand{
		ID:    objID,
		Title: req.Title,
	}

	cmd.Questions = make([]struct {
		QuestionID primitive.ObjectID `json:"questionId" validate:"required"`
		Order      int                `json:"order" validate:"required,min=1"`
	}, len(req.Questions))
	for i, q := range req.Questions {
		questionID, err := primitive.ObjectIDFromHex(q.QuestionID)
		if err != nil {
			SendErrorResponse(c, http.StatusBadRequest, "INVALID_QUESTION_ID",
				fmt.Sprintf("Question ID '%s' is not a valid ObjectID: %v", q.QuestionID, err))
			return
		}
		cmd.Questions[i] = struct {
			QuestionID primitive.ObjectID `json:"questionId" validate:"required"`
			Order      int                `json:"order" validate:"required,min=1"`
		}{
			QuestionID: questionID,
			Order:      q.Order,
		}
	}

	result, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			SendErrorResponse(c, http.StatusNotFound, "SURVEY_NOT_FOUND",
				"Survey with ID '"+id+"' not found")
			return
		}
		HandleServiceError(c, err)
		return
	}

	survey, ok := result.(*model.Survey)
	if !ok {
		SendErrorResponse(c, http.StatusInternalServerError, "INVALID_RESPONSE", "Invalid response type")
		return
	}

	SendSuccessResponse(c, http.StatusOK, survey)
}

func (h *SurveyHandler) DeleteSurvey(c *gin.Context) {
	id, ok := ValidateObjectID(c, c.Param("id"))
	if !ok {
		return
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	cmd := &command.DeleteSurveyCommand{
		ID: objID,
	}

	_, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			SendErrorResponse(c, http.StatusNotFound, "SURVEY_NOT_FOUND",
				"Survey with ID '"+id+"' not found")
			return
		}
		HandleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SurveyHandler) GetSurveyByToken(c *gin.Context) {
	token := c.Param("token")
	if len(token) != 5 {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_TOKEN", "Token must be exactly 5 characters long")
		return
	}

	q := &query.GetSurveyByTokenQuery{
		Token: token,
	}

	result, err := h.service.HandleQuery(c.Request.Context(), q)
	if err != nil {
		if err == repository.ErrSurveyNotFound {
			SendErrorResponse(c, http.StatusNotFound, "SURVEY_NOT_FOUND",
				"Survey with token '"+token+"' not found")
			return
		}
		HandleServiceError(c, err)
		return
	}

	SendSuccessResponse(c, http.StatusOK, result)
}
