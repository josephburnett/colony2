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
	subscribeBtn := dom.Doc.NewButton("Subscribe.")
	produceBtn := dom.Doc.NewButton("Produce.")
	pokeBtn := dom.Doc.NewButton("Poke.")
	p1.AppendChild(subscribeBtn)
	p1.AppendChild(produceBtn)
	p1.AppendChild(pokeBtn)

	ch := make(chan *protocol.ColonyReq, 1)
	subscribeBtn.OnClick(func(_ dom.Event) {
		ch <- &protocol.ColonyReq{
			Req: &protocol.ColonyReq_Subscribe{
				Subscribe: &protocol.Subscribe{},
			},
		}
	})
	produceBtn.OnClick(func(_ dom.Event) {
		ch <- &protocol.ColonyReq{
			Req: &protocol.ColonyReq_Produce{
				Produce: &protocol.Produce{},
			},
		}
	})
	pokeBtn.OnClick(func(_ dom.Event) {
		ch <- &protocol.ColonyReq{
			Req: &protocol.ColonyReq_Poke{
				Poke: &protocol.Poke{},
			},
		}
	})

	conn, err := grpc.Dial("ws://localhost:8080/ws", grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cli := protocol.NewColonyServiceClient(conn)

	printMsg := func(s string) {
		p := dom.Doc.CreateElement("p")
		p.SetTextContent(s)
		dom.Body.AppendChild(p)
	}

	stream, err := cli.Colony(context.Background())
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				panic(err)
			}
			if in.Messages != nil {
				for _, msg := range in.Messages {
					printMsg(msg)
				}
			}

		}
	}()
	for {
		req := <-ch
		if err := stream.Send(req); err != nil {
			panic(err)
		}
	}
	dom.Loop()
}
