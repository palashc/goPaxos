package main

import (
	"flag"
	"fmt"
	"log"
	"paxos/config"
	"paxos/randaddr"
	"strings"
)

var (
	nProposers = flag.Int("np", 2, "number of proposers")
	nAcceptors = flag.Int("na", 3, "number of acceptors")
	nLearners  = flag.Int("nl", 1, "number of learners")
	nFrontends = flag.Int("nf", 1, "number of frontends")
	frc        = flag.String("config", config.DefaultConfigPath, "paxos config file")
	ips        = flag.String("ips", "localhost", "comma-seperated list of IP addresses of the set of machines that'll host components")
)

func main() {
	flag.Parse()

	p := randaddr.RandPort()

	paxosConfig := new(config.PaxosConfig)
	paxosConfig.Proposers = make([]string, *nProposers)
	paxosConfig.Acceptors = make([]string, *nAcceptors)
	paxosConfig.Learners = make([]string, *nLearners)
	paxosConfig.Frontends = make([]string, *nFrontends)

	ipAddrs := strings.Split(*ips, ",")
	if nMachine := len(ipAddrs); nMachine > 0 {
		for i := 0; i < *nProposers; i++ {
			host := fmt.Sprintf("%s", ipAddrs[i%nMachine])
			paxosConfig.Proposers[i] = fmt.Sprintf("%s:%d", host, p)
			p++
		}

		for i := 0; i < *nAcceptors; i++ {
			host := fmt.Sprintf("%s", ipAddrs[i%nMachine])
			paxosConfig.Acceptors[i] = fmt.Sprintf("%s:%d", host, p)
			p++
		}

		for i := 0; i < *nLearners; i++ {
			host := fmt.Sprintf("%s", ipAddrs[i%nMachine])
			paxosConfig.Learners[i] = fmt.Sprintf("%s:%d", host, p)
			p++
		}
	}

	fmt.Println(paxosConfig.String())

	if *frc != "" {
		e := paxosConfig.Save(*frc)
		if e != nil {
			log.Fatal(e)
		}
	}
}
