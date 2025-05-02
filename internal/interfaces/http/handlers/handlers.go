package handlers

import (
	"github.com/gin-gonic/gin"
	ques_service "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/question/service"
	resp_service "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/response/service"
	surv_service "github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/service"
	"github.com/lccheungperry/OSP_backend/internal/infrastructure/mongodb/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	questionRepo := repository.NewMongoQuestionRepository(db)
	responseRepo := repository.NewMongoResponseRepository(db)
	surveyRepo := repository.NewMongoSurveyRepository(db)

	questionService := ques_service.NewQuestionService(questionRepo)
	responseService := resp_service.NewResponseService(responseRepo, surveyRepo, questionRepo)
	surveyService := surv_service.NewSurveyService(surveyRepo)

	questionHandler := NewQuestionHandler(questionService)
	responseHandler := NewResponseHandler(responseService)
	surveyHandler := NewSurveyHandler(surveyService)

	questions := router.Group("/questions")
	{
		questions.POST("", questionHandler.CreateQuestion)
		questions.GET("", questionHandler.ListQuestions)
		questions.GET("/:id", questionHandler.GetQuestion)
		questions.PUT("/:id", questionHandler.UpdateQuestion)
		questions.DELETE("/:id", questionHandler.DeleteQuestion)
	}

	surveys := router.Group("/surveys")
	{
		surveys.POST("", surveyHandler.CreateSurvey)
		surveys.GET("", surveyHandler.ListSurveys)
		surveys.GET("/:id", surveyHandler.GetSurvey)
		surveys.PUT("/:id", surveyHandler.UpdateSurvey)
		surveys.DELETE("/:id", surveyHandler.DeleteSurvey)
		surveys.GET("/token/:token", surveyHandler.GetSurveyByToken)

		surveys.POST("/:id/responses", responseHandler.SubmitResponse)
		surveys.GET("/:id/responses", responseHandler.GetResponses)
	}
}
