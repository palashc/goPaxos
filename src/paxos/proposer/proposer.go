package proposer

import (
	"fmt"
	"paxos"
	"paxos/types"
	"sync"
)

const MAX_RETRY = 3

type Proposer struct {
	id              int
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
		id:              id,
	}

	return p
}

func (p *Proposer) Set(value string, ret *bool) error {

	fmt.Printf("[Proposer:Set] Proposer %d got value %s\n", p.id, value)

	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.doConsensus(value, 0)
	if err != nil {
		fmt.Println("[Proposer:Set] Could not achieve consensus")
		*ret = false
		return err
	}
	fmt.Printf("[Proposer:Set] CONSENSUS\n")
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

	fmt.Printf("[Proposer:doConsensus] Performing Consensus for values %s, retry %d\n", value, retry+1)
	if retry == MAX_RETRY {
		return fmt.Errorf("Too many retries")
	}

	proposalNum := p.getNextProposalNumber()
	prepareReq := types.PrepareRequest{proposalNum, value}
	nPromises := 0
	fmt.Printf("[Proposer:doConsensus] Proposal Number: %d\n", proposalNum)

	// Phase-1: PREPARE
	fmt.Printf("[Proposer:doConsensus] PREPARE\n")
	for i, acceptor := range p.acceptors {
		var prepareRes types.PrepareResponse
		var proposalNum int
		err := acceptor.Prepare(prepareReq, &prepareRes)
		if err == nil && prepareRes.Status {
			fmt.Printf("[Proposer:doConsensus] Response from acceptor %d: %+v\n", i, prepareRes)
			nPromises++
			if prepareRes.PrevAccepted {
				if prepareRes.Proposal.N > proposalNum {
					fmt.Printf("[Proposer:doConsensus] acceptor %d had previously accepted larger proposalNum: %d, value: %s\n", i, prepareRes.Proposal.N, prepareRes.Proposal.V)
					value = prepareRes.Proposal.V
				}
			}

		}
	}

	fmt.Printf("[Proposer:doConsensus] Got %d promises\n", nPromises)

	//Check majority promises
	if nPromises >= int(len(p.acceptors)/2)+1 {
		fmt.Printf("[Proposer:doConsensus] ACCEPT\n")
		// Phase-2: Accept
		acceptReq := types.AcceptRequest{prepareReq.N, value}
		var acceptRes types.AcceptResponse
		var nAccepts int

		fmt.Printf("[Proposer:doConsensus] Informing Acceptors\n")
		for _, acceptor := range p.acceptors {
			err := acceptor.Accept(acceptReq, &acceptRes)
			if err == nil && acceptRes.Status {
				nAccepts++
			}
		}
		fmt.Printf("[Proposer:doConsensus] nAccepts: %d\n", nAccepts)
	} else {
		fmt.Printf("[Proposer:doConsensus] RETRY\n")
		p.doConsensus(value, retry+1)
	}
	return nil
}

func (p *Proposer) getNextProposalNumber() int {
	pNum := p.currProposalNum
	p.currProposalNum += p.numProposers
	return pNum
}
