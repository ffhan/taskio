package main

import (
	"context"
	"github.com/google/uuid"
	"taskio"
	"taskio/api"
)

type server struct {
	api.UnimplementedTaskApiServer
	taskRepository taskio.TaskRepository
}

func NewServer(taskRepository taskio.TaskRepository) *server {
	return &server{taskRepository: taskRepository}
}

func (s *server) SetResult(ctx context.Context, result *api.Result) (*api.Empty, error) {
	id, err := uuid.FromBytes(result.Id)
	if err != nil {
		return nil, err
	}
	return &api.Empty{}, s.taskRepository.StoreResult(ctx, id, result.Data)
}

func (s *server) GetTask(ctx context.Context, request *api.Request) (*api.Task, error) {
	nextTask, err := s.taskRepository.GetNext(ctx, request.Target)
	if err != nil {
		return nil, err
	}
	return &api.Task{
		Id:     nextTask.ID[:],
		Binary: nextTask.Data,
	}, nil
}
