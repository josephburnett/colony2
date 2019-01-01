package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dennwc/dom/net/ws"
	"github.com/josephburnett/colony2/pkg/protocol"
	"google.golang.org/grpc"
)

//go:generate GOOS=js GOARCH=wasm go build -o app.wasm ./cmd/client.go

func main() {

	view := &protocol.View{
		World: &protocol.World{
			Objects: map[int32]*protocol.World_ObjectRow{
				0: &protocol.World_ObjectRow{
					Columns: map[int32]*protocol.Object{
						0: &protocol.Object{
							Type: protocol.Object_WORKER,
						},
					},
				},
			},
		},
		XMin: -10,
		XMax: 10,
		YMin: -10,
		YMax: 10,
	}
	fmt.Println(view.Render())

	s := server{}

	srv := grpc.NewServer()
	protocol.RegisterHelloServiceServer(srv, s)

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
	stopCh := make(chan struct{})
	defer close(stopCh)
	go func() {
		tickerCh := time.NewTicker(30 * time.Second).C
		for {
			select {
			case <-stopCh:
				return
			case <-tickerCh:
				resp := &protocol.HelloResp{
					Text: "Hello?",
				}
				if err := stream.Send(resp); err != nil {
					return
				}
			}
		}
	}()
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
