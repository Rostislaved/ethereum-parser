package provider

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/config"
	"math/big"
	"net/http"
	"time"
)

type Provider struct {
	config       config.Provider
	client       http.Client
	hexconverter hexconverter
}

type hexconverter interface {
	DecodeBig(input string) (*big.Int, error)
	EncodeUint64(input uint64) string
}

func New(config config.Provider, hexconverter hexconverter) *Provider {
	client := http.Client{
		Timeout: time.Duration(config.ClientTimeoutInSecs) * time.Second,
	}

	return &Provider{
		config:       config,
		client:       client,
		hexconverter: hexconverter,
	}
}
