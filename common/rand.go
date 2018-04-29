package common

import (
	"errors"
	"math/rand"
	"time"
)

type Rate struct {
	r *rand.Rand

	rate int
}

func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func NewRate(rate float32) (*Rate, error) {
	if rate < 0 || rate > 1 {
		return nil, errors.New("rate is not allow!")
	}

	var r = &Rate{}
	r.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	r.rate = int(rate * 10000)
	return r, nil
}

func (r *Rate) Check() bool {
	return r.r.Intn(10000) <= r.rate
}
