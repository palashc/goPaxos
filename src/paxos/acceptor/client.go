package acceptor

import (
	"net/rpc"
	"paxos/types"
	"sync"
)

type AcceptorClient struct {
	Addr string
	conn *rpc.Client
	lock sync.Mutex
}

func GetNewAcceptorClient(addr string) *AcceptorClient {
	return &AcceptorClient{Addr: addr}
}

func (ac *AcceptorClient) rpcConn() error {
	ac.lock.Lock()
	defer ac.lock.Unlock()

	if ac.conn != nil {
		return nil
	}

	var err error
	ac.conn, err = rpc.DialHTTP("tcp", ac.Addr)
	if err != nil {
		ac.conn = nil
	}
	return err
}

func (ac *AcceptorClient) Prepare(req types.PrepareRequest, res *types.PrepareResponse) error {

	err := ac.rpcConn()
	if err != nil {
		return err
	}

	err = ac.conn.Call("Acceptor.Prepare", req, res)
	if err != nil {
		return err
	}

	return nil
}

func (ac *AcceptorClient) Accept(req types.AcceptRequest, res *types.AcceptResponse) error {

	err := ac.rpcConn()
	if err != nil {
		return err
	}

	err = ac.conn.Call("Acceptor.Accept", req, res)
	if err != nil {
		return err
	}

	return nil
}
