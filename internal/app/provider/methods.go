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

const (
	BlockNumber      = "eth_blockNumber"
	GetBlockByNumber = "eth_getBlockByNumber"
)

func (p *Provider) GetLastBlockNumber() (int64, error) {

	requestString := jsonrpc.NewRequest(BlockNumber)

	bodyReader := strings.NewReader(requestString)

	request, err := http.NewRequest("POST", p.config.URL, bodyReader)
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
		return 0, fmt.Errorf("got status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var dto blockNumberDTO

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
	numberHex := p.hexconverter.EncodeUint64(uint64(number))

	requestString := jsonrpc.NewRequest(GetBlockByNumber, numberHex, full)

	bodyReader := strings.NewReader(requestString)

	request, err := http.NewRequest("POST", p.config.URL, bodyReader)
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
		return entity.Block{}, fmt.Errorf("got status code: %d", resp.StatusCode)
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
