package world

import (
	"github.com/josephburnett/colony2/pkg/protocol"
)

var nextClientId int32 = 0

// Register submits a connected client and gets a unique client id.
func (w *RunningWorld) Register(stream protocol.ColonyService_ColonyServer) ClientId {
	log.Debugf("registering stream %v", stream)
	w.clientMux.Lock()
	defer w.clientMux.Unlock()
	id := atomicAddInt32(&nextClientId, 1)
	w.clients[id] = &client{
		id:     id,
		stream: stream,
	}
	log.Debugf("registered stream %v as client %v", stream, id)
	return id
}

// Unregister deletes a connected client by id and its subscriptions.
func (w *RunningWorld) Unregister(id ClientId) {
	log.Debugf("unregistering client %v", id)
	w.clientMux.Lock()
	defer w.clientsMux.Unlock()
	if colonies, ok := w.clientSubscriptions; ok {
		for c := range colonies {
			if ids, ok := w.clientSubscribers[c]; ok {
				delete(ids, id)
			}
			log.Printf("deleted client %v subscription to %v", id, c)
		}
		delete(w.clientSubscriptions, id)
	}
	if c, ok := w.clients[id]; ok {
		delete(w.clients, id)
	}
	log.Printf("unregistered client %v", id)
}

type client struct {
	id     ClientId
	stream protocol.ColonyService_ColonyServer
}

func (c *client) View(v *protocol.View, msg ...string) error {
	log.Debugf("sending colony %v view to client %v", "?", c.id)
	resp := &protocol.ColonyResp{
		Res: &protocol.ColonyResp_View{
			View: v,
		},
	}
	if len(msg) > 0 {
		resp.Messages = msg
	}
	return c.stream.Send(resp)
}

func (c *client) Update(u *protocol.Update, msg ...string) error {
	resp := &protocol.ColonyResp{
		Res: &protocol.ColonyResp_Update{
			Update: u,
		},
	}
	if len(msg) > 0 {
		resp.Messages = msg
	}
	return c.stream.Send(resp)
}

func (c *client) Error(err error, msg ...string) {
	resp := &protocol.ColonyResp{
		Messages: []string{err.Error()},
	}
	if len(msg) > 0 {
		resp.Messages = msg
	}
	return c.stream.Send(resp)
}

func (c *client) Msg(m string, msg ...string) {
	msg = append(msg, m)
	resp := &protocol.ColonyResp{
		Messages: msg,
	}
	return c.stream.Send(resp)
}
