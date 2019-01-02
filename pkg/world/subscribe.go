package world

import (
	"github.com/josephburnett/colony2/pkg/client"
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *RunningWorld) Subscribe(c *client.Client, req *protocol.Subscribe) error {
	// TODO: subscribe client to updates from a new View (or an existing one)
	c.Send(&protocol.ColonyResp{Messages: []string{"Subscribe!"}})
	return nil
}
