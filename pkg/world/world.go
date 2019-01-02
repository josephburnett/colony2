package world

import (
	"sync"

	"github.com/josephburnett/colony2/pkg/protocol"
)

type RunningWorld struct {

	// Source-of-truth World.
	worldMux sync.RWMutex
	world    *protocol.World

	// Colony-specific views of the World.
	views map[int32]*view
}

func NewRunningWorld() *RunningWorld {
	w := &RunningWorld{
		world: &protocol.World{},
	}
	return w
}
