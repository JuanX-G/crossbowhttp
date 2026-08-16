package crossbowhttp

import (
	"net/http"
	"slices"
)

type chain []func(http.Handler) http.Handler

func (c chain) thenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

func (c chain) then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}

func (c *chain) Add(f func(http.Handler) http.Handler) {
	tmp := append(*c, f)
	c = &tmp
}
