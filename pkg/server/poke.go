package server

import (
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *WorldServer) poke(req *protocol.Poke) ([]string, error) {
	return []string{
		"Poke!",
	}, nil
}
