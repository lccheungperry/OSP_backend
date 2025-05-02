package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	platform_error "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/error"
)

func SendErrorResponse(c *gin.Context, status int, errorCode string, details string) {
	c.JSON(status, gin.H{
		"error":   errorCode,
		"details": details,
	})
}

func HandleServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *platform_error.ServiceError:
		switch e.Code {
		case "SURVEY_NOT_FOUND", "QUESTION_NOT_FOUND":
			SendErrorResponse(c, http.StatusNotFound, e.Code, e.Message)
		case "INVALID_SPECIFICATIONS", "INVALID_OPTIONS", "INVALID_ANSWER_TYPE",
			"INVALID_OPTION_VALUE", "INVALID_SCALE_VALUE", "QUESTION_NOT_IN_SURVEY":
			SendErrorResponse(c, http.StatusBadRequest, e.Code, e.Message)
		case "RESPONSE_CREATION_ERROR", "RESPONSE_RETRIEVAL_ERROR":
			SendErrorResponse(c, http.StatusInternalServerError, e.Code, e.Message)
		default:
			SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", e.Message)
		}
	default:
		SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}

func ValidateObjectID(c *gin.Context, id string) (string, bool) {
	if id == "" {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID parameter is required")
		return "", false
	}
	return id, true
}

func ParsePaginationParams(c *gin.Context) (page, pageSize int64, ok bool) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err = strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 64)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize, true
}
