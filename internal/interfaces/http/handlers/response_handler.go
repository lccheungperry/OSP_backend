package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseHandler struct {
	service service.ResponseServiceInterface
}

func NewResponseHandler(service service.ResponseServiceInterface) *ResponseHandler {
	return &ResponseHandler{service: service}
}

func (h *ResponseHandler) SubmitResponse(c *gin.Context) {
	surveyID, ok := ValidateObjectID(c, c.Param("id"))
	if !ok {
		return
	}

	var req struct {
		Answers []struct {
			QuestionID string      `json:"questionId" binding:"required"`
			Value      interface{} `json:"value" binding:"required"`
		} `json:"answers" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if !ValidateResponseRequest(c, req.Answers) {
		return
	}

	surveyOID, err := primitive.ObjectIDFromHex(surveyID)
	if err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_SURVEY_ID",
			fmt.Sprintf("Survey ID '%s' is not a valid ObjectID: %v", surveyID, err))
		return
	}

	answers := make([]model.Answer, len(req.Answers))
	for i, ans := range req.Answers {
		questionOID, err := primitive.ObjectIDFromHex(ans.QuestionID)
		if err != nil {
			SendErrorResponse(c, http.StatusBadRequest, "INVALID_QUESTION_ID",
				fmt.Sprintf("Question ID '%s' is not a valid ObjectID: %v", ans.QuestionID, err))
			return
		}
		answers[i] = model.Answer{
			QuestionID: questionOID,
			Value:      ans.Value,
		}
	}

	cmd := &command.SubmitResponseCommand{
		SurveyID: surveyOID,
		Answers:  answers,
	}

	result, err := h.service.HandleCommand(c.Request.Context(), cmd)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	response, ok := result.(*model.Response)
	if !ok {
		SendErrorResponse(c, http.StatusInternalServerError, "INVALID_RESPONSE", "Invalid response type")
		return
	}

	SendSuccessResponse(c, http.StatusCreated, response)
}

func (h *ResponseHandler) GetResponses(c *gin.Context) {
	surveyID, ok := ValidateObjectID(c, c.Param("id"))
	if !ok {
		return
	}

	oid, err := primitive.ObjectIDFromHex(surveyID)
	if err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_SURVEY_ID",
			fmt.Sprintf("Survey ID '%s' is not a valid ObjectID: %v", surveyID, err))
		return
	}

	page, pageSize, ok := ParsePaginationParams(c)
	if !ok {
		return
	}

	skip := (page - 1) * pageSize

	q := &query.GetResponsesQuery{
		SurveyID: oid,
		Skip:     skip,
		Limit:    pageSize,
	}

	result, err := h.service.HandleQuery(c.Request.Context(), q)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	response := result.(struct {
		Responses []*model.Response `json:"responses"`
		Total     int64             `json:"total"`
	})

	SendSuccessResponse(c, http.StatusOK, response)
}
