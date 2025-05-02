# OSP Backend API Documentation

## Base URL
All API endpoints are prefixed with `/api`

## Authentication
Currently, the API is public and does not require authentication. Future versions may implement JWT-based authentication.

## Survey Management

### Create Survey
`POST /api/surveys`

Request:
```json
{
  "title": "Customer Satisfaction Survey",
  "questions": [
    {
      "questionId": "question_id_1",
      "order": 1
    },
    {
      "questionId": "question_id_2",
      "order": 2
    }
  ]
}
```

Response (201 Created):
```json
{
  "id": "survey_id",
  "title": "Customer Satisfaction Survey",
  "token": "abc12",
  "createdAt": "2024-03-20T10:00:00Z"
}
```

### List Surveys
`GET /api/surveys`

Response (200 OK):
```json
{
  "surveys": [
    {
      "id": "survey_id",
      "title": "Customer Satisfaction Survey",
      "token": "abc12",
      "createdAt": "2024-03-20T10:00:00Z"
    }
  ],
  "total": 1
}
```

### Get Survey by ID
`GET /api/surveys/{id}`

Response (200 OK):
```json
{
  "id": "survey_id",
  "title": "Customer Satisfaction Survey",
  "token": "abc12",
  "questions": [
    {
      "id": "question_id_1",
      "title": "How satisfied are you?",
      "format": "likert",
      "specifications": {
        "scale": ["Very Dissatisfied", "Dissatisfied", "Neutral", "Satisfied", "Very Satisfied"]
      }
    }
  ],
  "createdAt": "2024-03-20T10:00:00Z",
  "updatedAt": "2024-03-20T10:00:00Z"
}
```

### Get Survey by Token
`GET /api/surveys/token/{token}`

Response (200 OK):
```json
{
  "id": "survey_id",
  "title": "Customer Satisfaction Survey",
  "token": "abc12",
  "questions": [
    {
      "id": "question_id_1",
      "title": "How satisfied are you?",
      "format": "likert",
      "specifications": {
        "scale": ["Very Dissatisfied", "Dissatisfied", "Neutral", "Satisfied", "Very Satisfied"]
      }
    }
  ]
}
```

### Update Survey
`PUT /api/surveys/{id}`

Request:
```json
{
  "title": "Updated Survey Title",
  "questions": [
    {
      "questionId": "question_id_1",
      "order": 1
    }
  ]
}
```

Response (200 OK):
```json
{
  "id": "survey_id",
  "title": "Updated Survey Title",
  "token": "abc12",
  "updatedAt": "2024-03-20T11:00:00Z"
}
```

### Delete Survey
`DELETE /api/surveys/{id}`

Response (204 No Content)

## Response Management

### Submit Response
`POST /api/surveys/{surveyId}/responses`

Request:
```json
{
  "answers": [
    {
      "questionId": "question_id_1",
      "value": "Satisfied"
    }
  ]
}
```

Response (201 Created):
```json
{
  "id": "response_id",
  "surveyId": "survey_id",
  "createdAt": "2024-03-20T12:00:00Z"
}
```

### Get Survey Responses
`GET /api/surveys/{surveyId}/responses`

Response (200 OK):
```json
{
  "responses": [
    {
      "id": "response_id",
      "surveyId": "survey_id",
      "answers": [
        {
          "questionId": "question_id_1",
          "value": "Satisfied"
        }
      ],
      "createdAt": "2024-03-20T12:00:00Z"
    }
  ],
  "total": 1
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request format",
  "details": "Field 'title' is required"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found",
  "details": "Survey with ID 'survey_id' not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "details": "An unexpected error occurred"
}
```
