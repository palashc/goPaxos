package proposer

import (
	"fmt"
	"paxos"
	"sync"
	"types"
)

type Proposer struct {
	currProposalNum int
	value           string
	addr            string
	acceptors       []paxos.Acceptor
	numProposers    int
	lock            sync.Mutex
}

func (p *Proposer) Set(value string, ret *bool) error {

	//lock using per-key lock
	p.lock.Lock()
	defer p.lock.Unlock()

	err = p.doConsensus(value)
	if err != nil {
		fmt.Println("[Proposer:Set] Could not achieve consensus")
		*ret = false
		return err
	}

	p.value = value
	*ret = true
	return nil
}

func (p *Proposer) Get(key string, value *string) error {

	p.lock.Lock()
	defer p.lock.Unlock()

	*value = v
	return nil
}

func (p *Proposer) doConsensus(value string) error {

	prepareReq := types.PrepareRequest{p.getNextProposalNumber(), valueStr}
	nPromises := 0

	// Phase-1: PREPARE
	for _, acceptor := range acceptors {
		var prepareRes types.PrepareResponse
		var proposalNum int
		err := acceptor.Prepare(prepareReq, &prepareRes)
		if err == nil && prepareRes.Status {
			nPromises++
			if prepareRes.PrevAccepted {
				if prepareRes.Proposal.N > proposalNum {
					value = prepareRes.Proposal.V
				}
			}

		}

	}

	//Check majority promises
	if nPromises >= int(len(p.acceptors)/2)+1 {
		// Phase-2: Accept
		acceptReq := types.AcceptRequest{prepareReq.N, value}
		var acceptRes types.AcceptResponse
		var nAccepts int
		for _, acceptor := range acceptors {
			err = acceptor.Accept(acceptReq, &acceptRes)
			if err == nil && acceptRes.Status {
				nAccepts++
			}
		}
	} else {

	}

}

func (p *Proposer) getNextProposalNumber() int {
	pNum := p.currProposalNum
	p.currProposalNum += numProposers
	return pNum
}
