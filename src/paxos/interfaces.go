package paxos

import (
	"paxos/types"
)

type ProposerInterface interface {
	KV
}

type AcceptorInterface interface {
	Prepare(req types.PrepareRequest, res *types.PrepareResponse) error
	Accept(req types.AcceptRequest, res *types.AcceptResponse) error
}

type LearnerInterface interface {
}

type KV interface {
	Get(key string, value *string) error
	Set(value string, ret *bool) error
}
