package taskio

import (
	"context"
	"github.com/google/uuid"
)

type TaskStateCode uint

const (
	Error      TaskStateCode = iota // couldn't execute the task or get the results
	New                             // new issued task, before it's sent to a client
	Dispatched                      // task dispatched to a client, before returned results
	Results                         // returned results from a client, before processing the results
	Done                            // results processed, process is completed
)

type TaskState struct {
	Code    TaskStateCode
	Message string
}

type Task struct {
	ID    uuid.UUID
	Data  []byte
	State TaskState
}

type TaskRepository interface {
	GetNext(ctx context.Context, target string) (Task, error)
	StoreResult(ctx context.Context, id uuid.UUID, result []byte) error
	StoreTask(task Task) error
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
