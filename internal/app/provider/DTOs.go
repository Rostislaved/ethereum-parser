package provider

import "github.com/Rostislaved/ethereum-parser/internal/app/entity"

type blockNumberDTO struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type blockDTO struct {
	Id      int          `json:"id"`
	Jsonrpc string       `json:"jsonrpc"`
	Result  entity.Block `json:"result"`
}
