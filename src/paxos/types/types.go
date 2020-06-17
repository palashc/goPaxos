package types

type PrepareRequest struct {
	N int
	V string
}

type PrepareResponse struct {
	Status       bool
	PrevAccepted bool
	Proposal     PrepareRequest
}

type AcceptRequest struct {
	N int
	V string
}

type AcceptResponse struct {
	Status bool
	N      int
}
