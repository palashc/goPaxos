package proposer

import (
	"net"
	"net/http"
	"net/rpc"
	"paxos/config"
)

// Serve as a backend based on the given configuration
func Serve(b *config.ProposerConfig) error {

	server := rpc.NewServer()
	server.Register(b.Proposer)

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
