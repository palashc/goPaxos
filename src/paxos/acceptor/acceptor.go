package acceptor

import (
	"paxos/types"
	"sync"
)

type Acceptor struct {
	id            int
	addr          string
	lock          sync.Mutex
	maxPrepareNum int
	accepted      bool
	acceptNum     int
	acceptValue   string
	value         string
}

func NewAcceptor(id int, addr string) *Acceptor {

	a := &Acceptor{
		id:   id,
		addr: addr,
	}
	return a
}

func (a *Acceptor) Prepare(req types.PrepareRequest, res *types.PrepareResponse) error {
	panic("todo")
}

func (a *Acceptor) Accept(req types.AcceptRequest, res *types.AcceptResponse) error {
	panic("todo")
}
