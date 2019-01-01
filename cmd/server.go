package main

import (
	"io"
	"log"
	"net/http"

	"github.com/dennwc/dom/net/ws"
	"github.com/josephburnett/colony2/pkg/protocol"
	"github.com/josephburnett/colony2/pkg/server"
	"google.golang.org/grpc"
)

//go:generate GOOS=js GOARCH=wasm go build -o app.wasm ./cmd/client.go

var world *server.WorldServer

func main() {

	s := webServer{}
	world = server.NewWorldServer()

	srv := grpc.NewServer()
	protocol.RegisterColonyServiceServer(srv, s)

	const host = "localhost:8080"

	handler := http.FileServer(http.Dir("./resources"))
	lis, err := ws.Listen("ws://"+host+"/ws", handler)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	log.Printf("listening on http://%s", host)
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}

type webServer struct{}

func (webServer) Colony(stream protocol.ColonyService_ColonyServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		resp, err := world.Request(in)
		if err != nil {
			return err
		}
		if resp == nil {
			continue
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
