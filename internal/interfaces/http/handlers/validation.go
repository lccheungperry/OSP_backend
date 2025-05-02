package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateSurveyRequest(c *gin.Context, title string, questions []struct {
	QuestionID string `json:"questionId" binding:"required"`
	Order      int    `json:"order" binding:"required,min=1"`
}) bool {
	if title == "" {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Title is required")
		return false
	}

	for _, q := range questions {
		if _, err := primitive.ObjectIDFromHex(q.QuestionID); err != nil {
			SendErrorResponse(c, http.StatusBadRequest, "INVALID_QUESTION_ID",
				"Invalid question ID format: "+q.QuestionID)
			return false
		}
	}

	return true
}

func ValidateResponseRequest(c *gin.Context, answers []struct {
	QuestionID string      `json:"questionId" binding:"required"`
	Value      interface{} `json:"value" binding:"required"`
}) bool {
	if len(answers) == 0 {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "At least one answer is required")
		return false
	}

	for _, ans := range answers {
		if _, err := primitive.ObjectIDFromHex(ans.QuestionID); err != nil {
			SendErrorResponse(c, http.StatusBadRequest, "INVALID_QUESTION_ID",
				"Invalid question ID format: "+ans.QuestionID)
			return false
		}
	}

	return true
}
