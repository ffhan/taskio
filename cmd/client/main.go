package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"runtime"
	"taskio"
	"taskio/api"
	"taskio/container"
	"time"
)

func run(client api.TaskApiClient) {
	const defaultTimeout = 5 * time.Second

	backoff := time.Minute

	target := fmt.Sprintf("%s %s", runtime.GOARCH, runtime.GOOS)

	for {
		timeout, _ := context.WithTimeout(context.Background(), defaultTimeout)
		task, err := client.GetTask(timeout, &api.Request{Target: target})
		if err != nil {
			logrus.Errorf("cannot get a task: %v", err)
			time.Sleep(backoff)
			if backoff < 30*time.Minute {
				backoff *= 2
			}
			continue
		}
		id, err := uuid.FromBytes(task.Id)
		if err != nil {
			logrus.Errorf("cannot parse the UUID: %v", err)
			if _, err := client.SetResult(timeout, &api.Result{
				Id:      task.Id,
				Data:    nil,
				Success: false,
			}); err != nil {
				logrus.Errorf("cannot return the result: %v", err)
				continue
			}
		}
		fmt.Printf("got a task %s\n", id.String())
		var buf bytes.Buffer
		container.Run(bytes.NewBuffer(task.Binary), &buf)

		timeout, _ = context.WithTimeout(context.Background(), defaultTimeout)
		if _, err := client.SetResult(timeout, &api.Result{
			Id:      task.Id,
			Data:    buf.Bytes(),
			Success: true,
		}); err != nil {
			logrus.Errorf("cannot return the result: %v", err)
			continue
		}

		time.Sleep(30 * time.Second)
	}
}

func runClient() {
	dial, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock())
	taskio.Must(err)

	client := api.NewTaskApiClient(dial)

	run(client)
}

func main() {
	if os.Getpid() == 1 {
		container.RunChild()
		return
	}
	runClient()
}
