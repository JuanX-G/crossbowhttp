package json

import (
	"encoding/json/v2"
	"fmt"
)

type basicErrorRes struct {
	Error string `json:"error"`
}

func BasicErrorMapper(err error) ([]byte, int) {
	res, err := json.Marshal(err.Error())
	if err != nil {
		res = []byte(fmt.Sprintf("{error: \"%s\"}", err.Error()))
	}
	return res, 500
}
