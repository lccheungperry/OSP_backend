package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/command"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/bus/query"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/model"
	"github.com/lccheungperry/OSP_backend/internal/domain/survey_platform/survey/repository"
	"github.com/lccheungperry/OSP_backend/internal/interfaces/http/handlers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestSurveyHandler(t *testing.T) {
	t.Run("CreateSurvey", func(t *testing.T) {
		tests := []struct {
			name           string
			requestBody    interface{}
			expectedStatus int
			handleCommand  func(ctx context.Context, cmd command.Command) error
		}{
			{
				name: "successful survey creation",
				requestBody: map[string]interface{}{
					"title": "Test Survey",
				},
				expectedStatus: http.StatusCreated,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return nil
				},
			},
			{
				name: "invalid request body",
				requestBody: map[string]interface{}{
					"title": "", // Empty title
				},
				expectedStatus: http.StatusBadRequest,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return nil
				},
			},
			{
				name: "repository error",
				requestBody: map[string]interface{}{
					"title": "Test Survey",
				},
				expectedStatus: http.StatusInternalServerError,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return errors.New("repository error")
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := setupRouter()
				mockService := &MockSurveyService{
					HandleCommandFunc: tt.handleCommand,
				}
				handler := handlers.NewSurveyHandler(mockService)

				router.POST("/api/surveys", handler.CreateSurvey)

				body, _ := json.Marshal(tt.requestBody)
				req, _ := http.NewRequest(http.MethodPost, "/api/surveys", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("CreateSurvey() status = %v, want %v", w.Code, tt.expectedStatus)
				}
			})
		}
	})

	t.Run("GetSurvey", func(t *testing.T) {
		tests := []struct {
			name           string
			surveyID       string
			expectedStatus int
			handleQuery    func(ctx context.Context, q query.Query) (interface{}, error)
		}{
			{
				name:           "successful survey retrieval",
				surveyID:       primitive.NewObjectID().Hex(),
				expectedStatus: http.StatusOK,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return &model.Survey{
						ID:    primitive.NewObjectID(),
						Title: "Test Survey",
					}, nil
				},
			},
			{
				name:           "invalid survey ID",
				surveyID:       "invalid-id",
				expectedStatus: http.StatusBadRequest,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return nil, errors.New("invalid survey ID")
				},
			},
			{
				name:           "survey not found",
				surveyID:       primitive.NewObjectID().Hex(),
				expectedStatus: http.StatusNotFound,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return nil, repository.ErrSurveyNotFound
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := setupRouter()
				mockService := &MockSurveyService{
					HandleQueryFunc: tt.handleQuery,
				}
				handler := handlers.NewSurveyHandler(mockService)

				router.GET("/api/surveys/:id", handler.GetSurvey)

				req, _ := http.NewRequest(http.MethodGet, "/api/surveys/"+tt.surveyID, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("GetSurvey() status = %v, want %v", w.Code, tt.expectedStatus)
				}
			})
		}
	})

	t.Run("ListSurveys", func(t *testing.T) {
		tests := []struct {
			name           string
			expectedStatus int
			handleQuery    func(ctx context.Context, q query.Query) (interface{}, error)
		}{
			{
				name:           "successful survey listing",
				expectedStatus: http.StatusOK,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return struct {
						Surveys []*model.Survey `json:"surveys"`
						Total   int64           `json:"total"`
					}{
						Surveys: []*model.Survey{
							{
								ID:    primitive.NewObjectID(),
								Title: "Test Survey 1",
							},
							{
								ID:    primitive.NewObjectID(),
								Title: "Test Survey 2",
							},
						},
						Total: 2,
					}, nil
				},
			},
			{
				name:           "repository error",
				expectedStatus: http.StatusInternalServerError,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return nil, errors.New("repository error")
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := setupRouter()
				mockService := &MockSurveyService{
					HandleQueryFunc: tt.handleQuery,
				}
				handler := handlers.NewSurveyHandler(mockService)

				router.GET("/api/surveys", handler.ListSurveys)

				req, _ := http.NewRequest(http.MethodGet, "/api/surveys", nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("ListSurveys() status = %v, want %v", w.Code, tt.expectedStatus)
				}
			})
		}
	})

	t.Run("UpdateSurvey", func(t *testing.T) {
		tests := []struct {
			name           string
			surveyID       string
			requestBody    interface{}
			expectedStatus int
			handleCommand  func(ctx context.Context, cmd command.Command) error
		}{
			{
				name:     "successful survey update",
				surveyID: primitive.NewObjectID().Hex(),
				requestBody: map[string]interface{}{
					"title": "Updated Survey",
				},
				expectedStatus: http.StatusOK,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return nil
				},
			},
			{
				name:     "invalid survey ID",
				surveyID: "invalid-id",
				requestBody: map[string]interface{}{
					"title": "Updated Survey",
				},
				expectedStatus: http.StatusBadRequest,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return errors.New("invalid survey ID")
				},
			},
			{
				name:     "survey not found",
				surveyID: primitive.NewObjectID().Hex(),
				requestBody: map[string]interface{}{
					"title": "Updated Survey",
				},
				expectedStatus: http.StatusNotFound,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return repository.ErrSurveyNotFound
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := setupRouter()
				mockService := &MockSurveyService{
					HandleCommandFunc: tt.handleCommand,
				}
				handler := handlers.NewSurveyHandler(mockService)

				router.PUT("/api/surveys/:id", handler.UpdateSurvey)

				body, _ := json.Marshal(tt.requestBody)
				req, _ := http.NewRequest(http.MethodPut, "/api/surveys/"+tt.surveyID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("UpdateSurvey() status = %v, want %v", w.Code, tt.expectedStatus)
				}
			})
		}
	})

	t.Run("DeleteSurvey", func(t *testing.T) {
		tests := []struct {
			name           string
			surveyID       string
			expectedStatus int
			handleCommand  func(ctx context.Context, cmd command.Command) error
		}{
			{
				name:           "successful survey deletion",
				surveyID:       primitive.NewObjectID().Hex(),
				expectedStatus: http.StatusNoContent,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return nil
				},
			},
			{
				name:           "invalid survey ID",
				surveyID:       "invalid-id",
				expectedStatus: http.StatusBadRequest,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return errors.New("invalid survey ID")
				},
			},
			{
				name:           "survey not found",
				surveyID:       primitive.NewObjectID().Hex(),
				expectedStatus: http.StatusNotFound,
				handleCommand: func(ctx context.Context, cmd command.Command) error {
					return repository.ErrSurveyNotFound
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := setupRouter()
				mockService := &MockSurveyService{
					HandleCommandFunc: tt.handleCommand,
				}
				handler := handlers.NewSurveyHandler(mockService)

				router.DELETE("/api/surveys/:id", handler.DeleteSurvey)

				req, _ := http.NewRequest(http.MethodDelete, "/api/surveys/"+tt.surveyID, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("DeleteSurvey() status = %v, want %v", w.Code, tt.expectedStatus)
				}
			})
		}
	})

	t.Run("GetSurveyByToken", func(t *testing.T) {
		tests := []struct {
			name           string
			token          string
			expectedStatus int
			handleQuery    func(ctx context.Context, q query.Query) (interface{}, error)
		}{
			{
				name:           "successful survey retrieval by token",
				token:          "abcde",
				expectedStatus: http.StatusOK,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return &model.Survey{
						ID:    primitive.NewObjectID(),
						Title: "Test Survey",
						Token: "abcde",
					}, nil
				},
			},
			{
				name:           "invalid token",
				token:          "invalid",
				expectedStatus: http.StatusBadRequest,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return nil, errors.New("invalid token")
				},
			},
			{
				name:           "survey not found",
				token:          "abcde",
				expectedStatus: http.StatusNotFound,
				handleQuery: func(ctx context.Context, q query.Query) (interface{}, error) {
					return nil, repository.ErrSurveyNotFound
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := setupRouter()
				mockService := &MockSurveyService{
					HandleQueryFunc: tt.handleQuery,
				}
				handler := handlers.NewSurveyHandler(mockService)

				router.GET("/api/surveys/token/:token", handler.GetSurveyByToken)

				req, _ := http.NewRequest(http.MethodGet, "/api/surveys/token/"+tt.token, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("GetSurveyByToken() status = %v, want %v", w.Code, tt.expectedStatus)
				}
			})
		}
	})
}
