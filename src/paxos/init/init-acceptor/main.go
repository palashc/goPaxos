package main

import (
	"flag"
	"fmt"
	"log"
	"paxos"
	"paxos/acceptor"
	"paxos/config"
	"paxos/learner"
)

var frc = flag.String("conf", config.DefaultConfigPath, "config file")

// if -1, run all workers
var pid = flag.Int("pid", -1, "which proposer to run")

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	flag.Parse()

	rc, e := config.LoadConfig(*frc)
	noError(e)

	learners := []paxos.LearnerInterface{}
	for _, lAddr := range rc.Learners {
		learners = append(learners, learner.GetNewLearnerClient(lAddr))
	}

	run := func(i int) {
		if i > len(rc.Acceptors) {
			noError(fmt.Errorf("index out of range: %d", i))
		}

		aAddr := rc.Acceptors[i]
		aConfig := rc.NewAcceptorConfig(i, acceptor.NewAcceptor(i, aAddr, learners))

		log.Printf("monitor serving on %s", aConfig.Addr)

		noError(acceptor.Serve(aConfig))
	}

	// run all monitors
	if *pid == -1 {
		for i, _ := range rc.Acceptors {
			go run(i)
		}
	} else {
		run(*pid)
	}

	select {}

}
