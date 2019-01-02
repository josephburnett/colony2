package world

import (
	"github.com/josephburnett/colony2/pkg/client"
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *RunningWorld) Poke(c *client.Client, req *protocol.Poke) error {
	// TODO: poke an ant at the location if owned by subscribed colony.
	c.Send(&protocol.ColonyResp{Messages: []string{"Poke!"}})
	return nil
}
