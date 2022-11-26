package jsonrpc

import "encoding/json"

type Param json.RawMessage

type Params []Param

func (m Param) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

func MakeParams(params ...interface{}) (Params, error) {
	if len(params) == 0 {
		return nil, nil
	}

	out := make(Params, len(params))
	for i, param := range params {
		b, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}

		out[i] = Param(b)
	}

	return out, nil
}

type JSONRPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	Id      int    `json:"id"`
}

func (req JSONRPCRequest) ToJson() (string, error) {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
