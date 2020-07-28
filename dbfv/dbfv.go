package dbfv

import (
	"github.com/ldsec/lattigo/bfv"
	"github.com/ldsec/lattigo/ring"
	"github.com/ldsec/lattigo/utils"
)

type dbfvContext struct {
	params *bfv.Parameters

	// Polynomial degree
	n uint64

	// floor(Q/T) mod each Qi in Montgomery form
	deltaMont []uint64

	// Polynomial contexts
	contextT  *ring.Context
	contextQ  *ring.Context
	contextP  *ring.Context
	contextQP *ring.Context
}

func newDbfvContext(params *bfv.Parameters) *dbfvContext {

	LogN := params.LogN
	n := uint64(1 << LogN)

	contextT, err := ring.NewContextWithParams(n, []uint64{params.T})
	if err != nil {
		panic(err)
	}

	contextQ, err := ring.NewContextWithParams(n, params.Qi)
	if err != nil {
		panic(err)
	}

	contextP, err := ring.NewContextWithParams(n, params.Pi)
	if err != nil {
		panic(err)
	}

	contextQP, err := ring.NewContextWithParams(n, append(params.Qi, params.Pi...))
	if err != nil {
		panic(err)
	}

	deltaMont := bfv.GenLiftParams(contextQ, params.T)

	return &dbfvContext{
		params:    params.Copy(),
		n:         n,
		deltaMont: deltaMont,
		contextT:  contextT,
		contextQ:  contextQ,
		contextP:  contextP,
		contextQP: contextQP,
	}
}

// NewCRPGenerator creates a new deterministic random polynomial generator.
func NewCRPGenerator(params *bfv.Parameters, key []byte) *ring.UniformSampler {
	ctx := newDbfvContext(params)
	prng, err := utils.NewKeyedPRNG(key)
	if err != nil {
		panic(err)
	}
	return ring.NewUniformSampler(prng, ctx.contextQP)
}
