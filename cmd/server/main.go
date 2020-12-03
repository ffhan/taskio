package main

import (
	"google.golang.org/grpc"
	"net"
	"taskio"
	"taskio/api"
	"taskio/helloworld/repository"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	taskio.Must(err)

	server := grpc.NewServer()
	taskApiServer := NewServer(repository.NewTaskRepository()) // todo: implement repo

	api.RegisterTaskApiServer(server, taskApiServer)
	taskio.Must(server.Serve(listen))
}
