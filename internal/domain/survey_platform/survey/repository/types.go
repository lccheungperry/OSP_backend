package repository

import "fmt"

type RepositoryError struct {
	Code    string
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}

func NewAssignmentError(err error) error {
	return &RepositoryError{
		Code:    "ASSIGNMENT_ERROR",
		Message: fmt.Sprintf("failed to handle assignment: %v", err),
	}
}

func NewDeleteError(err error) error {
	return &RepositoryError{
		Code:    "DELETE_ERROR",
		Message: fmt.Sprintf("failed to delete: %v", err),
	}
}

func NewCreateError(err error) error {
	return &RepositoryError{
		Code:    "CREATE_ERROR",
		Message: fmt.Sprintf("failed to create: %v", err),
	}
}

func NewUpdateError(err error) error {
	return &RepositoryError{
		Code:    "UPDATE_ERROR",
		Message: fmt.Sprintf("failed to update: %v", err),
	}
}
