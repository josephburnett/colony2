package server

import (
	"io"

	"github.com/josephburnett/colony2/pkg/client"
	"github.com/josephburnett/colony2/pkg/protocol"
	"github.com/josephburnett/colony2/pkg/world"
)

var _ protocol.ColonyServiceServer = Server{}

// Server knows how to dispatch Colony messages to a RunningWorld.
// Implements the ColonyServiceServer interface.
type Server struct {
	World *world.RunningWorld
}

func (s Server) Colony(stream protocol.ColonyService_ColonyServer) error {

	// New Client to represent this unique connection.
	c := client.NewClient(stream)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if subscribe := req.GetSubscribe(); subscribe != nil {
			err := s.World.Subscribe(c, subscribe)
			if err != nil {
				c.Error(err)
				return nil
			}
			return nil
		}
		if produce := req.GetProduce(); produce != nil {
			err := s.World.Produce(c, produce)
			if err != nil {
				c.Error(err)
				return nil
			}
			return nil
		}
		if poke := req.GetPoke(); poke != nil {
			err := s.World.Poke(c, poke)
			if err != nil {
				c.Error(err)
				return nil
			}
			return nil
		}
	}
}
