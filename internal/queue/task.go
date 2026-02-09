package queue

import "context"

// Task is a unit of work to be executed in a lane.
type Task struct {
	SessionID string
	Fn        func(ctx context.Context) error
}
