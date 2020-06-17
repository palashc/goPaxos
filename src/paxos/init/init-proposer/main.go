package main

import (
	"flag"
	"fmt"
	"log"
	"paxos"
	"paxos/acceptor"
	"paxos/config"
	"paxos/proposer"
)

var frc = flag.String("conf", config.DefaultConfigPath, "config file")

// if -1, run all workers
var workerId = flag.Int("pid", -1, "which proposer to run")

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	flag.Parse()

	rc, e := config.LoadConfig(*frc)
	noError(e)

	acceptorClients := []paxos.AcceptorInterface{}
	for _, acceptorAddr := range rc.Acceptors {
		acceptorClients = append(acceptorClients, acceptor.GetNewAcceptorClient(acceptorAddr))
	}

	run := func(i int) {
		if i > len(rc.Proposers) {
			noError(fmt.Errorf("index out of range: %d", i))
		}

		pAddr := rc.Proposers[i]
		pConfig := rc.NewProposerConfig(i, proposer.NewProposer(i, pAddr, acceptorClients, len(rc.Proposers)))

		log.Printf("monitor serving on %s", pConfig.Addr)

		noError(proposer.Serve(pConfig))
	}

	// run all monitors
	if *pid == -1 {
		for i, _ := range rc.Proposers {
			go run(i)
		}
	} else {
		run(*pid)
	}

	select {}

}
