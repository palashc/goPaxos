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
	fmt.Printf("[Acceptor %d:Prepare] Got prepare request with PN: %d, value: %s\n", a.id, req.N, req.V)

	if req.N > a.maxPrepareNum {
		fmt.Printf("[Acceptor %d:Prepare] This is a larger proposal number\n", a.id)
		a.maxPrepareNum = req.N
		res.Status = true

		if a.accepted {
			fmt.Printf("[Acceptor %d:Prepare] Aleady accepted a previous proposal with PN: %d, valye: %s\n", a.id, a.acceptNum, a.acceptValue)
			res.PrevAccepted = true
			acceptedProposal := types.Proposal{a.acceptNum, a.acceptValue}
			res.Proposal = acceptedProposal
		} else {
			fmt.Printf("[Acceptor %d:Prepare] Did not accept anything previously\n", a.id)
			res.PrevAccepted = false
		}
	} else {
		fmt.Printf("[Acceptor %d:Prepare] Aleady seen larger proposal number\n", a.id)
		res.Status = false
	}

	return nil
}

func (a *Acceptor) Accept(req types.AcceptRequest, res *types.AcceptResponse) error {

	a.lock.Lock()
	defer a.lock.Unlock()
	fmt.Printf("[Acceptor %d:Accept] Got accept request with PN: %d, value: %s\n", a.id, req.N, req.V)

	if req.N >= a.maxPrepareNum {
		fmt.Printf("[Acceptor %d:Accept] Accepting this value!\n", a.id)
		res.Status = true
		a.accepted = true
		a.acceptNum = req.N
		a.acceptValue = req.V
		a.value = req.V

		// notify learners
		fmt.Printf("[Acceptor %d:Accept] Notifying Learners\n", a.id)
		for i, learner := range a.learners {
			var ret bool
			err := learner.Notify(a.value, &ret)
			if err != nil || !ret {
				fmt.Printf("[Acceptor %d:Accept] Could not notify learner %d", a.id, i)
			}
		}
	} else {
		fmt.Printf("[Acceptor %d:Accept] Aleady seen larger proposal number\n", a.id)
		res.Status = false
		res.N = -1
	}

	return nil
}
