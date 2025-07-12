package action

import "context"

// Executor is the interface for executing actions
type Executor interface {
	// Execute runs the action
	Execute(ctx context.Context) error

	// String returns a string representation of the action
	String() string
}
