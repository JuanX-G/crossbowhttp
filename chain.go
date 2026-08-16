package crossbowhttp

import (
	"net/http"
	"slices"
)

// A simple chain of middlewares to be put in front of an endpoint
type Chain []func(http.Handler) http.Handler

func (c Chain) thenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

func (c Chain) then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}

func (c *Chain) Add(f func(http.Handler) http.Handler) {
	tmp := append(*c, f)
	c = &tmp
}
