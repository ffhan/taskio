package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"taskio"
)

type repo struct {
	task taskio.Task
}

func NewTaskRepository() *repo {
	binary, err := ioutil.ReadFile("helloworld/helloworld")
	taskio.Must(err)

	return &repo{task: taskio.Task{
		ID:   uuid.New(),
		Data: binary,
		State: taskio.TaskState{
			Code:    taskio.New,
			Message: "",
		},
	}}
}

func (r *repo) GetNext(ctx context.Context, target string) (taskio.Task, error) {
	r.task.ID = uuid.New()
	return r.task, nil
}

func (r *repo) StoreResult(ctx context.Context, id uuid.UUID, result []byte) error {
	fmt.Printf("stored result %s for %s\n", string(result), id)
	return nil
}

func (r *repo) StoreTask(task taskio.Task) error {
	r.task = task
	return nil
}
