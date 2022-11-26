package parser

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"log"
)

func (p *Parser) GetCurrentBlock() int {
	return int(p.lastParsedBlockNumber) // not cool, but I have to implement the interface
}

func (p *Parser) Subscribe(address string) bool {
	_, ok := p.subscribedAddressesSet[address]
	if !ok {
		p.subscribedAddressesSet[address] = struct{}{}

		return true
	}

	return false // not quite sure why interface has bool return value, so I return this
}

func (p *Parser) GetTransactions(address string) []entity.Transaction {
	transactions, err := p.storage.GetTransactions(address)
	if err != nil {
		log.Println(err)

		return nil
	}

	return transactions
}

func (p *Parser) GetStorageInfo() entity.StorageInfo {
	storageInfo, err := p.storage.GetStorageInfo()
	if err != nil {
		log.Println(err)

		return entity.StorageInfo{}
	}

	return storageInfo
}
