package main

import (
	"flag"
	"fmt"
	"log"
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

	proposer := proposer.GetNewProposerClient(rc.Proposers[0])
	learner := learner.GetNewLearnerClient(rc.Learners[0])

	myValue := "hello_world"

	var ret bool
	err := proposer.Set(myValue, &ret)
	if err != nil || !ret {
		fmt.Println(err)
		panic(err)
	}

	time.Sleep(1 * time.Second)

	var S string
	err = learner.Get("", &S)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if S == myValue {
		fmt.Println("PASS.")
	} else {
		fmt.Printf("FAIL. Expected: %s, Actual: %s\n", myValue, S)
	}
}
