package world

import (
	"log"
	"sync/atomic"

	"github.com/josephburnett/colony2/pkg/protocol"
)

var nextClientId int32 = 0

// Register submits a connected client and gets a unique client id.
func (w *RunningWorld) Register(stream protocol.ColonyService_ColonyServer) ClientId {
	log.Printf("registering stream %v", stream)
	w.clientMux.Lock()
	defer w.clientMux.Unlock()
	id := ClientId(atomic.AddInt32(&nextClientId, 1))
	w.clients[id] = &client{
		id:     id,
		stream: stream,
	}
	log.Printf("registered stream %v as client %v", stream, id)
	return id
}

// Unregister deletes a connected client by id and its subscriptions.
func (w *RunningWorld) Unregister(id ClientId) {
	log.Printf("unregistering client %v", id)
	w.clientMux.Lock()
	defer w.clientMux.Unlock()
	if colonies, ok := w.clientSubscriptions[id]; ok {
		for c := range colonies {
			if ids, ok := w.clientSubscribers[c]; ok {
				delete(ids, id)
			}
			log.Printf("deleted client %v subscription to %v", id, c)
		}
		delete(w.clientSubscriptions, id)
	}
	if _, ok := w.clients[id]; ok {
		delete(w.clients, id)
	}
	log.Printf("unregistered client %v", id)
}

type client struct {
	id     ClientId
	stream protocol.ColonyService_ColonyServer
}

func (c *client) View(v *protocol.View, msg ...string) error {
	log.Printf("sending colony %v view to client %v", "?", c.id)
	resp := &protocol.ColonyResp{
		Req: &protocol.ColonyResp_View{
			View: v,
		},
	}
	if len(msg) > 0 {
		resp.Messages = msg
	}
	return c.stream.Send(resp)
}

func (c *client) Update(u *protocol.ViewUpdate, msg ...string) error {
	log.Printf("sending colony %v update to client %v", "?", c.id)
	resp := &protocol.ColonyResp{
		Req: &protocol.ColonyResp_Update{
			Update: u,
		},
	}
	if len(msg) > 0 {
		resp.Messages = msg
	}
	return c.stream.Send(resp)
}

func (c *client) Error(err error, msg ...string) error {
	log.Printf("sending error %q to client %v", err.Error(), c.id)
	resp := &protocol.ColonyResp{
		Messages: []string{err.Error()},
	}
	if len(msg) > 0 {
		resp.Messages = msg
	}
	return c.stream.Send(resp)
}

func (c *client) Msg(m string, msg ...string) error {
	log.Printf("sending messages %q ... to client %v", m, c.id)
	msg = append(msg, m)
	resp := &protocol.ColonyResp{
		Messages: msg,
	}
	return c.stream.Send(resp)
}
