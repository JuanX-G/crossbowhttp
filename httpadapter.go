package crossbowhttp

import (
	"net/http"

	"github.com/JuanX-G/crossbow"
)

type HttpMarshaller[M any, O any] interface {
	Decode(*http.Request) (M, error)
	Econde(O, error) ([]byte, int, error)
	MapError(error) ([]byte, int)
}

type HttpAdapter[M any, O any] struct {
	server     *crossbow.Server[crossbow.ServerHandler[M, O], M, O]
	marshaller HttpMarshaller[M, O]
}

func (h *HttpAdapter[M, O]) encodeError(w http.ResponseWriter, err error) {
	res, code := h.marshaller.MapError(err)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	w.WriteHeader(code)
	w.Write(res)
}

func (h *HttpAdapter[M, O]) HandleSend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := h.marshaller.Decode(r)

		err = h.server.Send(r.Context(), req)
		if err != nil {
			h.encodeError(w, err)
			return
		}
	}
}

func (h *HttpAdapter[M, O]) HandleCall() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := h.marshaller.Decode(r)
		if err != nil {
			h.encodeError(w, err)
			return
		}

		out, err := h.server.Call(r.Context(), req)

		res, code, err := h.marshaller.Econde(out, err)
		if err != nil {
			h.encodeError(w, err)
		}

		w.WriteHeader(code)
		w.Write(res)
	}
}

type RegisterHttpAdapterParams struct {
	HealthEnc        HealthEncoder
	urlBase          string
	SendMiddleware   chain
	CallMiddleware   chain
	HealthMiddleware chain
}

func RegisterHttpAdapter[M any, O any](mux *http.ServeMux, adapter HttpAdapter[M, O], params RegisterHttpAdapterParams) {
	mux.Handle("POST "+params.urlBase+"/send", params.SendMiddleware.thenFunc(adapter.HandleSend()))
	mux.Handle("POST "+params.urlBase+"/call", params.CallMiddleware.thenFunc(adapter.HandleCall()))
	if params.HealthEnc != nil {
		mux.Handle("GET "+params.urlBase+"/health", params.HealthMiddleware.thenFunc(adapter.HealthHandler(params.HealthEnc)))
	}
}
