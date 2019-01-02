package main

import (
	"log"
	"net/http"

	"github.com/dennwc/dom/net/ws"
	"github.com/josephburnett/colony2/pkg/protocol"
	"github.com/josephburnett/colony2/pkg/server"
	"github.com/josephburnett/colony2/pkg/world"
	"google.golang.org/grpc"
)

//go:generate GOOS=js GOARCH=wasm go build -o app.wasm ./cmd/client.go

const (
	listenAt = "ws://localhost:8080/ws"
)

func main() {

	// Create a new, empty World.
	world := world.NewRunningWorld()

	// Setup a Colony GRPC over web socket server.
	s := server.Server{
		World: world,
	}
	srv := grpc.NewServer()
	protocol.RegisterColonyServiceServer(srv, s)

	// Serve the client code.
	handler := http.FileServer(http.Dir("./resources"))
	lis, err := ws.Listen(listenAt, handler)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	log.Printf("listening %s", listenAt)
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}
