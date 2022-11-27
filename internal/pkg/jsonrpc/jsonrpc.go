package jsonrpc

import (
	"encoding/json"
)

type param json.RawMessage

type params []param

func (m param) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

func MakeParams(params ...interface{}) (params, error) {
	if len(params) == 0 {
		return nil, nil
	}

	out := make([]param, len(params))
	for i, p := range params {
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}

		out[i] = param(b)
	}

	return out, nil
}

type jsonRPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  params `json:"params"`
	Id      int    `json:"id"`
}

func NewRequest(method string, params ...interface{}) (requestString string) {
	p, err := MakeParams(params...)
	if err != nil {
		return
	}

	jsonrpcRequest := jsonRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  p,
		Id:      83,
	}

	jsonString, err := jsonrpcRequest.toJSON()
	if err != nil {
		return
	}

	return jsonString
}

func (req jsonRPCRequest) toJSON() (string, error) {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
