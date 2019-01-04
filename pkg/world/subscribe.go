package world

import (
	"fmt"
	"log"

	"github.com/josephburnett/colony2/pkg/protocol"
)

func (w *RunningWorld) Subscribe(id ClientId, req *protocol.Subscribe) error {
	log.Printf("subscribing client %v to colony %v", id, req.Id)

	// Get the client.
	w.clientMux.Lock()
	defer w.clientMux.Unlock()
	client, ok := w.clients[id]
	if !ok {
		return fmt.Errorf("unknown client id %v", id)
	}

	// Get the colony.
	colonyId := ColonyId(req.Id)
	w.worldMux.RLock()
	defer w.worldMux.RUnlock()
	if _, ok := w.world.Colonies[int32(colonyId)]; !ok {
		client.Error(fmt.Errorf("unknown colony %v", colonyId))
		return nil
	}

	// Subscribe the client to the colony.
	subscriptions, ok := w.clientSubscriptions[id]
	if !ok {
		subscriptions = make(map[ColonyId]bool)
		w.clientSubscriptions[id] = subscriptions
	}
	subscriptions[colonyId] = true

	// And add the client as a subscriber.
	subscribers, ok := w.clientSubscribers[colonyId]
	if !ok {
		subscribers = make(map[ClientId]bool)
		w.clientSubscribers[colonyId] = subscribers
	}
	subscribers[id] = true

	log.Printf("subscribed client %v to colony %v", id, req.Id)
	return nil
}
