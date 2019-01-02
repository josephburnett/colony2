package world

import (
	"github.com/josephburnett/colony2/pkg/client"
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *RunningWorld) Produce(c *client.Client, req *protocol.Produce) error {
	// TODO: set the production state for the subscribed colony.
	c.Send(&protocol.ColonyResp{Messages: []string{"Produce!"}})
	return nil
}
