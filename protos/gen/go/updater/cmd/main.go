package main

import (
	"context"
	pbn "obsidianGoNaive/protos/gen/go/notes"
	pbu "obsidianGoNaive/protos/gen/go/updater"
	"obsidianGoNaive/protos/gen/go/updater/use_case"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	server := grpc.NewServer()
	pbu.RegisterUpdaterServer(server, &use_case.UpdaterService{})

	conn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pbn.NewNotesClient(conn)
}
