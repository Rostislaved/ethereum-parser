package parser

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/config"
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
)

type Parser struct {
	config                 config.Parser
	lastParsedBlockNumber  int64
	subscribedAddressesSet map[string]struct{}
	storage                Storage
	provider               Provider
}

type Storage interface {
	SaveTransaction(address string, transaction entity.Transaction) (err error)
	GetTransactions(address string) (transactions []entity.Transaction, err error)
	GetStorageInfo() (entity.StorageInfo, error)
}

type Provider interface {
	GetLastBlockNumber() (int64, error)
	GetBlockByNumber(number int64, full bool) (block entity.Block, err error)
}

func New(cfg config.Parser, storage Storage, provider Provider) *Parser {
	subscribedAddressesSet := make(map[string]struct{})

	return &Parser{
		config:                 cfg,
		lastParsedBlockNumber:  cfg.InitialBlockNumber,
		subscribedAddressesSet: subscribedAddressesSet,
		storage:                storage,
		provider:               provider,
	}
}
