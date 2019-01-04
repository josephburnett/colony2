package world

import (
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *RunningWorld) Produce(id ClienId, req *protocol.Produce) error {
	w.clientMux.RLock()
	defer w.clientMux.RUnlock()
	if c, ok := w.clients[id]; ok {
		// TODO: set the production state for the subscribed colony.
		c.Msg("Produce!")
		return nil
	}
}
