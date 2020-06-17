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
	fmt.Printf("[Acceptor:Prepare] Got prepare request with PN: %d, value: %s\n", req.N, req.V)

	if req.N > a.maxPrepareNum {
		fmt.Printf("[Acceptor:Prepare] This is a larger proposal number\n")
		a.maxPrepareNum = req.N
		res.Status = true

		if a.accepted {
			fmt.Printf("[Acceptor:Prepare] Aleady accepted a previous proposal with PN: %d, valye: %s\n", a.acceptNum, a.acceptValue)
			res.PrevAccepted = true
			acceptedProposal := types.Proposal{a.acceptNum, a.acceptValue}
			res.Proposal = acceptedProposal
		} else {
			fmt.Printf("[Acceptor:Prepare] Did not accept anything previously\n")
			res.PrevAccepted = false
		}
	} else {
		fmt.Printf("[Acceptor:Prepare] Aleady seen larger proposal number\n")
		res.Status = false
	}

	return nil
}

func (a *Acceptor) Accept(req types.AcceptRequest, res *types.AcceptResponse) error {

	a.lock.Lock()
	defer a.lock.Unlock()
	fmt.Printf("[Acceptor:Accept] Got accept request with PN: %d, value: %s\n", req.N, req.V)

	if req.N >= a.maxPrepareNum {
		fmt.Printf("[Acceptor:Accept] Accepting this value!\n")
		a.accepted = true
		a.acceptNum = req.N
		a.acceptValue = req.V
		a.value = req.V

		// notify learners
		fmt.Printf("[Acceptor:Accept] Notifying Learners\n")
		for i, learner := range a.learners {
			var ret bool
			err := learner.Notify(a.value, &ret)
			if err != nil || !ret {
				fmt.Println("[Acceptor:Accept] Could not notify learner ", i)
			}
		}
	} else {
		fmt.Printf("[Acceptor:Accept] Aleady seen larger proposal number\n")
		res.Status = false
		res.N = -1
	}

	return nil
}
