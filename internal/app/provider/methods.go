package provider

import (
	"encoding/json"
	"fmt"
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"github.com/Rostislaved/ethereum-parser/internal/pkg/jsonrpc"
	"io"
	"net/http"
	"strings"
)

func (p *Provider) GetLastBlockNumber() (int64, error) {
	method := "eth_blockNumber"

	jsonrpcRequest := jsonrpc.JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  nil,
		Id:      83,
	}

	jsonString, err := jsonrpcRequest.ToJson()
	if err != nil {
		return 0, err
	}

	bodyReader := strings.NewReader(jsonString)

	request, err := http.NewRequest("POST", p.url, bodyReader)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.Do(request)
	if err != nil {
		return 0, err
	}

	defer func() {
		errC := resp.Body.Close()
		if errC != nil {
			if err == nil {
				err = errC
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("got status code: %s", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var dto jsonrpcDTO

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return 0, err
	}

	result := dto.Result

	quantity, err := p.hexconverter.DecodeBig(result)
	if err != nil {
		return 0, err
	}

	return quantity.Int64(), nil
}

func (p *Provider) GetBlockByNumber(number int64, full bool) (block entity.Block, err error) {
	method := "eth_getBlockByNumber"

	numberHex := p.hexconverter.EncodeUint64(uint64(number))

	params, err := jsonrpc.MakeParams(numberHex, full)
	if err != nil {
		return entity.Block{}, err
	}

	jsonrpcRequest := jsonrpc.JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      83,
	}

	jsonString, err := jsonrpcRequest.ToJson()
	if err != nil {
		return entity.Block{}, err
	}

	bodyReader := strings.NewReader(jsonString)

	request, err := http.NewRequest("POST", p.url, bodyReader)
	if err != nil {
		return entity.Block{}, err
	}

	resp, err := p.client.Do(request)
	if err != nil {
		return entity.Block{}, err
	}

	defer func() {
		errC := resp.Body.Close()
		if errC != nil {
			if err == nil {
				err = errC
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return entity.Block{}, fmt.Errorf("got status code: %s", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.Block{}, err
	}

	var dto blockDTO

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return entity.Block{}, err
	}

	block = dto.Result

	return block, nil
}

type blockDTO struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  entity.Block `json:"result"`
	Id      int          `json:"id"`
}
