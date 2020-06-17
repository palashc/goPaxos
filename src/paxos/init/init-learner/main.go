package main

import (
	"flag"
	"fmt"
	"log"
	"paxos/config"
	"paxos/learner"
)

var frc = flag.String("conf", config.DefaultConfigPath, "config file")

// if -1, run all workers
var pid = flag.Int("pid", -1, "which learner to run")

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	flag.Parse()

	rc, e := config.LoadConfig(*frc)
	noError(e)

	run := func(i int) {
		if i > len(rc.Learners) {
			noError(fmt.Errorf("index out of range: %d", i))
		}

		lAddr := rc.Learners[i]
		lConfig := rc.NewLearnerConfig(i, learner.NewLearner(i, lAddr))

		log.Printf("monitor serving on %s", lConfig.Addr)

		noError(learner.Serve(lConfig))
	}

	// run all monitors
	if *pid == -1 {
		for i, _ := range rc.Learners {
			go run(i)
		}
	} else {
		run(*pid)
	}

	select {}

}
