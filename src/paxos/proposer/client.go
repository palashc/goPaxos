package proposer

import (
	"net/rpc"
	"sync"
)

type ProposerClient struct {
	Addr string
	conn *rpc.Client
	lock sync.Mutex
}

func GetNewClient(addr string) *ProposerClient {
	return &ProposerClient{Addr: addr}
}

func (pc *ProposerClient) rpcConn() error {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	if pc.conn != nil {
		return nil
	}

	var err error
	pc.conn, err = rpc.DialHTTP("tcp", pc.Addr)
	if err != nil {
		pc.conn = nil
	}
	return err
}

func (pc *ProposerClient) Get(key string, value *string) error {

	err := pc.rpcConn()
	if err != nil {
		return err
	}

	err = pc.conn.Call("Proposer.Get", key, value)
	if err != nil {
		return err
	}

	return nil
}

func (pc *ProposerClient) Set(value string, ret *bool) error {

	err := pc.rpcConn()
	if err != nil {
		return err
	}

	err = pc.conn.Call("Proposer.Set", value, ret)
	if err != nil {
		return err
	}

	return nil
}
