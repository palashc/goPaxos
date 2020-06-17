package acceptor

import (
	"fmt"
	"paxos"
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
	learners      []paxos.LearnerInterface
}

func NewAcceptor(id int, addr string, learners []paxos.LearnerInterface) *Acceptor {

	a := &Acceptor{
		id:            id,
		addr:          addr,
		maxPrepareNum: -1,
		accepted:      false,
		learners:      learners,
	}
	return a
}

func (a *Acceptor) Prepare(req types.PrepareRequest, res *types.PrepareResponse) error {

	a.lock.Lock()
	defer a.lock.Unlock()

	if req.N > a.maxPrepareNum {
		a.maxPrepareNum = req.N
		res.Status = true
		if a.accepted {
			res.PrevAccepted = true
			acceptedProposal := types.Proposal{a.acceptNum, a.acceptValue}
			res.Proposal = acceptedProposal
		} else {
			res.PrevAccepted = false
		}
	} else {
		res.Status = false
	}
	return nil
}

func (a *Acceptor) Accept(req types.AcceptRequest, res *types.AcceptResponse) error {

	a.lock.Lock()
	defer a.lock.Unlock()

	if req.N >= a.maxPrepareNum {
		a.accepted = true
		a.acceptNum = req.N
		a.acceptValue = req.V
		a.value = req.V

		// notify learners
		for i, learner := range a.learners {
			var ret bool
			err := learner.Notify(a.value, &ret)
			if err != nil || !ret {
				fmt.Println("[Acceptor:Accept] Could not notify learner ", i)
			}
		}
	} else {
		res.Status = false
		res.N = -1
	}

	return nil
}
