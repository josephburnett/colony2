package server

import (
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *WorldServer) produce(req *protocol.Produce) ([]string, error) {
	return []string{
		"Produce!",
	}, nil
}
