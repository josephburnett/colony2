package world

import (
	"sync"

	"github.com/josephburnett/colony2/pkg/protocol"
)

type ColonyId int32
type ClientId int32

// RunningWorld implements the Colony game logic around a proto World.
type RunningWorld struct {

	// Source-of-truth proto World.
	world    *protocol.World
	worldMux sync.RWMutex

	// Colony-specific views of the World.
	views    map[ColonyId]*view
	viewsMux sync.RWMutex

	// Registered clients.
	clients             map[ClientId]*client
	clientSubscriptions map[ClientId]map[ColonyId]bool
	clientSubscribers   map[ColonyId]map[ClientId]bool
	clientMux           sync.RWMutex
}

// NewRunningWorld starts a new, empty RunningWorld.
func NewRunningWorld() *RunningWorld {
	w := &RunningWorld{
		world: &protocol.World{},
	}
	return w
}
