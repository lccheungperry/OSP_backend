package command

import (
	"context"
)

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	HandleCommand(ctx context.Context, cmd Command) (interface{}, error)
}

type CommandBus interface {
	RegisterHandler(commandName string, handler CommandHandler)
	Dispatch(ctx context.Context, cmd Command) (interface{}, error)
}

type InMemoryCommandBus struct {
	handlers map[string]CommandHandler
}

func NewInMemoryCommandBus() *InMemoryCommandBus {
	return &InMemoryCommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

func (b *InMemoryCommandBus) RegisterHandler(commandName string, handler CommandHandler) {
	b.handlers[commandName] = handler
}

func (b *InMemoryCommandBus) Dispatch(ctx context.Context, cmd Command) (interface{}, error) {
	handler, exists := b.handlers[cmd.CommandName()]
	if !exists {
		return nil, ErrCommandHandlerNotFound
	}
	return handler.HandleCommand(ctx, cmd)
}

var ErrCommandHandlerNotFound = NewCommandError("command handler not found")

type CommandError struct {
	message string
}

func NewCommandError(message string) *CommandError {
	return &CommandError{message: message}
}

func (e *CommandError) Error() string {
	return e.message
}
