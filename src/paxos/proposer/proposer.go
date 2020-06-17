package proposer

import (
	"fmt"
	"paxos"
	"paxos/types"
	"sync"
)

const MAX_RETRY = 3

type Proposer struct {
	currProposalNum int
	value           string
	addr            string
	acceptors       []paxos.AcceptorInterface
	numProposers    int
	lock            sync.Mutex
}

func NewProposer(id int, addr string, acceptors []paxos.AcceptorInterface, numProposers int) *Proposer {

	p := &Proposer{
		currProposalNum: id,
		addr:            addr,
		acceptors:       acceptors,
		numProposers:    numProposers,
	}

	return p
}

func (p *Proposer) Set(value string, ret *bool) error {

	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.doConsensus(value, 0)
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

	*value = p.value
	return nil
}

func (p *Proposer) doConsensus(value string, retry int) error {

	if retry == MAX_RETRY {
		return fmt.Errorf("Too many retries")
	}

	prepareReq := types.PrepareRequest{p.getNextProposalNumber(), value}
	nPromises := 0

	// Phase-1: PREPARE
	for _, acceptor := range p.acceptors {
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
		for _, acceptor := range p.acceptors {
			err := acceptor.Accept(acceptReq, &acceptRes)
			if err == nil && acceptRes.Status {
				nAccepts++
			}
		}
	} else {
		p.doConsensus(value, retry+1)
	}
	return nil
}

func (p *Proposer) getNextProposalNumber() int {
	pNum := p.currProposalNum
	p.currProposalNum += p.numProposers
	return pNum
}
