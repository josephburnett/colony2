package main

import (
	"context"
	"net"
	"time"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/net/ws"
	"github.com/josephburnett/colony2/pkg/protocol"
	"google.golang.org/grpc"
)

func dialer(s string, dt time.Duration) (net.Conn, error) {
	return ws.Dial(s)
}

func main() {
	p1 := dom.Doc.CreateElement("p")
	dom.Body.AppendChild(p1)

	inp := dom.Doc.NewInput("text")
	p1.AppendChild(inp)

	btn := dom.Doc.NewButton("Go!")
	p1.AppendChild(btn)

	ch := make(chan string, 1)
	btn.OnClick(func(_ dom.Event) {
		ch <- inp.Value()
	})

	conn, err := grpc.Dial("ws://localhost:8080/ws", grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cli := protocol.NewHelloServiceClient(conn)

	printMsg := func(s string) {
		p := dom.Doc.CreateElement("p")
		p.SetTextContent(s)
		dom.Body.AppendChild(p)
	}

	stream, err := cli.Hello(context.Background())
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				panic(err)
			}
			printMsg(in.Text)
		}
	}()
	for {
		name := <-ch
		printMsg("say hello to: " + name)

		req := &protocol.HelloReq{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			panic(err)
		}
	}
	dom.Loop()
}
