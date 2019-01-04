package server

import (
	"io"

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

	// Register this connection as a unique client.
	id := s.World.Register(stream)
	defer s.World.Unregister(id)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if subscribe := req.GetSubscribe(); subscribe != nil {
			if err := s.World.Subscribe(id, subscribe); err != nil {
				return err
			}
		}
		if produce := req.GetProduce(); produce != nil {
			if err := s.World.Produce(id, produce); err != nil {
				return err
			}
		}
		if poke := req.GetPoke(); poke != nil {
			if err := s.World.Poke(id, poke); err != nil {
				return err
			}
		}
	}
}
