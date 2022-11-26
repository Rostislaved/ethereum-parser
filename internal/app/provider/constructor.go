package provider

import (
	"math/big"
	"net/http"
	"time"
)

type Provider struct {
	client       http.Client
	url          string
	hexconverter hexconverter
}

type hexconverter interface {
	DecodeBig(input string) (*big.Int, error)
	EncodeUint64(input uint64) string
}

func New(hexconverter hexconverter) *Provider {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://cloudflare-eth.com"

	return &Provider{
		client:       client,
		url:          url,
		hexconverter: hexconverter,
	}
}
