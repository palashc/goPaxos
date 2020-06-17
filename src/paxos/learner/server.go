package learner

import (
	"net"
	"net/http"
	"net/rpc"
	"paxos/config"
)

func Serve(b *config.LearnerConfig) error {

	server := rpc.NewServer()
	server.Register(b.Learner)

	listener, err := net.Listen("tcp", b.Addr)
	if err != nil {
		if b.Ready != nil {
			b.Ready <- false
		}
		return err
	}

	if b.Ready != nil {
		b.Ready <- true
	}

	return http.Serve(listener, server)
}
