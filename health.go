package crossbowhttp

import (
	"net/http"

	"github.com/JuanX-G/crossbow"
)

type ServerHealth struct {
	Stats      crossbow.StatsSnapshot
	MailboxLen int
	Terminated bool
}

type HealthEncoder = func(ServerHealth) ([]byte, int)

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
