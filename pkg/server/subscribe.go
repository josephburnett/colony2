package server

import (
	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *WorldServer) subscribe(req *protocol.Subscribe) (*protocol.View, []string, error) {
	return nil, []string{
		"Subscribe!",
	}, nil
}
