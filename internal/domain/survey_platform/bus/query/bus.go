package query

import (
	"context"
)

type Query interface {
	QueryName() string
}

type QueryHandler interface {
	HandleQuery(ctx context.Context, query Query) (interface{}, error)
}

type QueryBus interface {
	RegisterHandler(queryName string, handler QueryHandler)
	Dispatch(ctx context.Context, query Query) (interface{}, error)
}

type InMemoryQueryBus struct {
	handlers map[string]QueryHandler
}

func NewInMemoryQueryBus() *InMemoryQueryBus {
	return &InMemoryQueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

func (b *InMemoryQueryBus) RegisterHandler(queryName string, handler QueryHandler) {
	b.handlers[queryName] = handler
}

func (b *InMemoryQueryBus) Dispatch(ctx context.Context, query Query) (interface{}, error) {
	handler, exists := b.handlers[query.QueryName()]
	if !exists {
		return nil, nil
	}
	return handler.HandleQuery(ctx, query)
}

var ErrQueryHandlerNotFound = NewQueryError("query handler not found")

type QueryError struct {
	message string
}

func NewQueryError(message string) *QueryError {
	return &QueryError{message: message}
}

func (e *QueryError) Error() string {
	return e.message
}
