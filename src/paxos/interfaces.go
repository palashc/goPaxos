package paxos

type Proposer interface {
	KV
}

type Acceptor interface {
}

type Learner interface {
}

type KV interface {
	Get(key string, value *string) error
	Set(value string, ret *bool) error
}
