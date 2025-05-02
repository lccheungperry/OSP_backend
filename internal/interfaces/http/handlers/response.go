package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendSuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func FormatSurveyResponse(survey interface{}) gin.H {
	if s, ok := survey.(struct {
		ID        primitive.ObjectID
		Title     string
		Token     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}); ok {
		return gin.H{
			"id":        s.ID.Hex(),
			"title":     s.Title,
			"token":     s.Token,
			"createdAt": s.CreatedAt.Format(time.RFC3339),
			"updatedAt": s.UpdatedAt.Format(time.RFC3339),
		}
	}
	return nil
}

func FormatResponseResponse(response interface{}) gin.H {
	if r, ok := response.(struct {
		ID        primitive.ObjectID
		SurveyID  primitive.ObjectID
		Answers   []interface{}
		CreatedAt time.Time
	}); ok {
		return gin.H{
			"id":        r.ID.Hex(),
			"surveyId":  r.SurveyID.Hex(),
			"answers":   r.Answers,
			"createdAt": r.CreatedAt.Format(time.RFC3339),
		}
	}
	return nil
}

func FormatPaginationResponse(items interface{}, total int64, page, pageSize int64) gin.H {
	return gin.H{
		"items":    items,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}
}
