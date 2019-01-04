package world

import (
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *RunningWorld) Poke(id ClientId, req *protocol.Poke) error {
	w.clientMux.RLock()
	defer w.clientMux.RUnlock()
	if c, ok := w.clients[id]; ok {
		// TODO: poke an ant at the location if owned by subscribed colony.
		c.Msg("Poke!")
	}
	return nil
}
