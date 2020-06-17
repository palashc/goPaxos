package learner

import (
	"net/rpc"
	"sync"
)

type LearnerClient struct {
	Addr string
	conn *rpc.Client
	lock sync.Mutex
}

func GetNewLearnerClient(addr string) *LearnerClient {
	return &LearnerClient{Addr: addr}
}

func (lc *LearnerClient) rpcConn() error {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	if lc.conn != nil {
		return nil
	}

	var err error
	lc.conn, err = rpc.DialHTTP("tcp", lc.Addr)
	if err != nil {
		lc.conn = nil
	}
	return err
}

func (lc *LearnerClient) Notify(value string, ret *bool) error {

	err := lc.rpcConn()
	if err != nil {
		return err
	}

	err = lc.conn.Call("Learner.Notify", value, ret)
	if err != nil {
		return err
	}

	return nil
}
