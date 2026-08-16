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

// called on transport/adapter layer errors (failing to decode a body, etc.). Writes headers and the response body
func (h *HttpAdapter[M, O]) encodeError(w http.ResponseWriter, err error) {
	res, code := h.marshaller.MapError(err)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	w.WriteHeader(code)
	w.Write(res)
}

// Returns the http handler function for asynchronous sends to the server
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

// Returns the http handler function for synchronous calls to the server
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

// Params for HttpAdapter to http.ServeMux registration
type RegisterHttpAdapterParams struct {
	HealthEnc        HealthEncoder // function that encoedes the server health to the desired output format
	urlBase          string        // base url for handlers. The final url is: urlBase + "/[call | send | health]"
	SendMiddleware   Chain         // middleware chain for the .../send endpoint
	CallMiddleware   Chain         // middleware chain for the .../call endpoint
	HealthMiddleware Chain         // middleware chain for the .../health endpoint
}

// Register a HttpAdapter to a http.ServeMux
func RegisterHttpAdapter[M any, O any](mux *http.ServeMux, adapter HttpAdapter[M, O], params RegisterHttpAdapterParams) {
	mux.Handle("POST "+params.urlBase+"/send", params.SendMiddleware.thenFunc(adapter.HandleSend()))
	mux.Handle("POST "+params.urlBase+"/call", params.CallMiddleware.thenFunc(adapter.HandleCall()))
	if params.HealthEnc != nil {
		mux.Handle("GET "+params.urlBase+"/health", params.HealthMiddleware.thenFunc(adapter.HealthHandler(params.HealthEnc)))
	}
}
