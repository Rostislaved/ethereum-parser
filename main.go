package main

import "github.com/Rostislaved/ethereum-parser/entity"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []entity.Transaction
}
