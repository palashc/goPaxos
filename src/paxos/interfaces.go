package paxos

type ProposerInterface interface {
	KV
}

type AcceptorInterface interface {
}

type LearnerInterface interface {
}

type KV interface {
	Get(key string, value *string) error
	Set(value string, ret *bool) error
}
