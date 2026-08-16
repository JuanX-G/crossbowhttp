package json

import (
	"encoding/json/v2"
	"net/http"

	"github.com/JuanX-G/crossbowhttp"
)

type JsonMarshaller[M any, O any] struct {
	JsonErrorMapper func(error) ([]byte, int)
}

func (j *JsonMarshaller[M, O]) Decode(r *http.Request) (M, error) {
	var val M
	err := json.UnmarshalRead(r.Body, &val)
	return val, err
}

func (j *JsonMarshaller[M, O]) Encode(val O, err error) ([]byte, int, error) {
	if err != nil {
		res, code := j.JsonErrorMapper(err)
		return res, code, nil
	}
	res, err := json.Marshal(val)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}
	return res, http.StatusOK, nil
}

func (j *JsonMarshaller[M, O]) MapError(err error) ([]byte, int) {
	return j.JsonErrorMapper(err)
}

type ServerHealthJson struct {
	Fails      uint64 `json:"fails"`
	Panics     uint64 `json:"panics"`
	MailboxLen int    `json:"mailbox_length"`
	Terminated bool   `json:"terminated"`
}

func JsonHealthEncoder(health crossbowhttp.ServerHealth) ([]byte, int) {
	healthStruct := ServerHealthJson{
		Fails:      health.Stats.Failures,
		Panics:     health.Stats.Panics,
		MailboxLen: health.MailboxLen,
		Terminated: health.Terminated,
	}
	res, err := json.Marshal(healthStruct)
	if err != nil {
		res, _ = BasicErrorMapper(err)
		return res, http.StatusInternalServerError
	}
	return res, http.StatusOK
}
