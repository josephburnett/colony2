package client

import (
	"sync/atomic"

	"github.com/josephburnett/colony2/pkg/protocol"
)

var nextClientId int32 = 0

// Client represents a unique connection to the Colony Server.
// Multiple Clients may subscribe to a single Colony.
// Clients are stateful and may connect to only one Colony at
// a time.
type Client struct {
	id     int32
	stream protocol.ColonyService_ColonyServer
}

// NewClient provisions a Client with a unique id.
func NewClient(stream protocol.ColonyService_ColonyServer) *Client {
	return &Client{
		id:     atomic.AddInt32(&nextClientId, 1),
		stream: stream,
	}
}

func (c *Client) Send(resp *protocol.ColonyResp) error {
	return c.stream.Send(resp)
}

func (c *Client) Error(err error) {
	c.stream.Send(&protocol.ColonyResp{
		Messages: []string{err.Error()},
	})
}
