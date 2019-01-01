package server

import (
	"sync"

	"github.com/josephburnett/colony2/pkg/protocol"
)

type WorldServer struct {
	worldMux sync.RWMutex
	world    *protocol.World
}

func NewWorldServer() *WorldServer {
	w := &WorldServer{
		world: &protocol.World{},
	}
	return w
}

func (w *WorldServer) Request(req *protocol.ColonyReq) (*protocol.ColonyResp, error) {
	if subscribe := req.GetSubscribe(); subscribe != nil {
		view, msg, err := w.subscribe(subscribe)
		if err != nil {
			return nil, err
		}
		return &protocol.ColonyResp{
			Req: &protocol.ColonyResp_View{
				View: view,
			},
			Messages: msg,
		}, nil
	}
	if produce := req.GetProduce(); produce != nil {
		msg, err := w.produce(produce)
		if err != nil {
			return nil, err
		}
		return &protocol.ColonyResp{
			Messages: msg,
		}, nil
	}
	if poke := req.GetPoke(); poke != nil {
		msg, err := w.poke(poke)
		if err != nil {
			return nil, err
		}
		return &protocol.ColonyResp{
			Messages: msg,
		}, nil
	}
	return nil, nil
}
