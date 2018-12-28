package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dennwc/dom/net/ws"
	"github.com/josephburnett/colony2/pkg/protocol"
	"google.golang.org/grpc"
)

//go:generate GOOS=js GOARCH=wasm go build -o app.wasm ./cmd/client.go

func main() {
	s := server{}

	srv := grpc.NewServer()
	protocol.RegisterService(srv, s)

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

type server struct{}

func (server) Hello(stream protocol.HelloService_HelloServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		resp := &protocol.HelloResp{
			Text: fmt.Sprintf("Hello, %s!", in.Name),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
