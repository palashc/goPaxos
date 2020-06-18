package main

import (
	"flag"
	"fmt"
	"log"
	"paxos"
	"paxos/config"
	"paxos/learner"
	"paxos/proposer"
	"time"
)

var frc = flag.String("conf", config.DefaultConfigPath, "config file")

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	flag.Parse()

	rc, e := config.LoadConfig(*frc)
	noError(e)

	proposer1 := proposer.GetNewProposerClient(rc.Proposers[0])
	proposer2 := proposer.GetNewProposerClient(rc.Proposers[1])

	learner1 := learner.GetNewLearnerClient(rc.Learners[0])
	learner2 := learner.GetNewLearnerClient(rc.Learners[1])

	myValue1 := "hello_world1"
	myValue2 := "hello_world2"

	propose := func(proposer paxos.ProposerInterface, v string) {
		var ret bool
		err := proposer.Set(v, &ret)
		if err != nil || !ret {
			fmt.Println(err)
			panic(err)
		}
	}

	go propose(proposer1, myValue1)
	//time.Sleep(5 * time.Millisecond)
	go propose(proposer2, myValue2)

	time.Sleep(1 * time.Second)

	var S1 string
	err := learner1.Get("", &S1)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var S2 string
	err = learner2.Get("", &S2)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if S1 == S2 {
		fmt.Printf("PASS with %s\n", S1)
	} else {
		fmt.Printf("FAIL. S1: %s, S2: %s\n", S1, S2)
	}
}
