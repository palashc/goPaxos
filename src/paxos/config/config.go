package config

import (
	"encoding/json"
	"fmt"
	"os"
	"paxos"
)

var DefaultConfigPath = "paxos_config.conf"

type ProposerConfig struct {
	Addr     string
	Proposer paxos.ProposerInterface
	Ready    chan bool
}

type AcceptorConfig struct {
	Addr     string
	Acceptor paxos.AcceptorInterface
	Ready    chan bool
}

type LearnerConfig struct {
	Addr    string
	Learner paxos.LearnerInterface
	Ready   chan bool
}

type PaxosConfig struct {
	Frontends []string
	Proposers []string
	Acceptors []string
	Learners  []string
}

func (pc *PaxosConfig) NewProposerConfig(i int, p paxos.ProposerInterface) *ProposerConfig {
	ret := new(ProposerConfig)
	ret.Addr = pc.Proposers[i]
	ret.Proposer = p
	ret.Ready = make(chan bool, 1)
	return ret
}

func (pc *PaxosConfig) NewAcceptorConfig(i int, a paxos.AcceptorInterface) *AcceptorConfig {
	ret := new(AcceptorConfig)
	ret.Addr = pc.Acceptors[i]
	ret.Acceptor = a
	ret.Ready = make(chan bool, 1)
	return ret
}

func (pc *PaxosConfig) NewLearnerConfig(i int, a paxos.LearnerInterface) *LearnerConfig {
	ret := new(LearnerConfig)
	ret.Addr = pc.Learners[i]
	ret.Learner = a
	ret.Ready = make(chan bool, 1)
	return ret
}

func (pc *PaxosConfig) Save(p string) error {
	b := pc.marshal()

	fout, e := os.Create(p)
	if e != nil {
		return e
	}

	_, e = fout.Write(b)
	if e != nil {
		return e
	}

	_, e = fmt.Fprintln(fout)
	if e != nil {
		return e
	}

	return fout.Close()
}
func (pc *PaxosConfig) Write(p string) (*os.File, error) {
	b := pc.marshal()

	fout, e := os.Create(p)
	if e != nil {
		return nil, e
	}

	_, e = fout.Write(b)
	if e != nil {
		return nil, e
	}

	return fout, nil
}

func (pc *PaxosConfig) String() string {
	b := pc.marshal()
	return string(b)
}

func LoadConfig(p string) (*PaxosConfig, error) {
	fin, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer fin.Close()

	ret := new(PaxosConfig)
	e = json.NewDecoder(fin).Decode(ret)
	if e != nil {
		return nil, e
	}

	return ret, nil
}

func (pc *PaxosConfig) marshal() []byte {
	b, e := json.MarshalIndent(pc, "", "    ")
	if e != nil {
		panic(e)
	}

	return b
}
