# Taskio

![logo](assets/logo_small.png)

A toy client-server system that distributes tasks to clients and runs them in containers - completely isolated from the
environment.

**The main purpose would be to enable running sensitive code on untrusted clients (with their knowledge) to utilise
their computing power and get the compute results back safely.**

Tasks are basically binaries targeted for the specific client architecture and OS - the client cannot see the binary at
all - and the binary cannot see the clients OS.

### How to run

Run these from the repository root.

#### Client

`go run cmd/client/main.go` or `go build -o client cmd/client/main.go && ./client`

#### Server

`go run cmd/server/main.go` or `go build -o server ./cmd/server && ./server`
