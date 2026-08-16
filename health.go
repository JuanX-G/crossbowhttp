package crossbowhttp

import (
	"net/http"

	"github.com/JuanX-G/crossbow"
)

// Health metrics of a server aggregated into a single struct for encoding
type ServerHealth struct {
	Stats      crossbow.StatsSnapshot
	MailboxLen int
	Terminated bool
}

// function mapping health metrics of a server to a byte blob ready for streaming over http.
// The second return is the http status code
type HealthEncoder = func(ServerHealth) ([]byte, int)

// returns the http handler for fetching the health metrics of the server.
func (h *HttpAdapter[M, O]) HealthHandler(enc HealthEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := ServerHealth{}
		health.Stats = h.server.Stats()
		health.Terminated = h.server.Terminated()
		health.MailboxLen = h.server.MailboxLen()
		health.Terminated = h.server.Terminated()
		res, code := enc(health)
		w.WriteHeader(code)
		w.Write(res)
	}
}
