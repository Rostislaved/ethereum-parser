package parser

import (
	"context"
	"fmt"
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"time"
)

func (p *Parser) Start(ctx context.Context) {
	//p.subscribedAddressesSet["0x06450dee7fd2fb8e39061434babcfc05599a6fb8"] = struct{}{}

	blockNumbersChan := p.getBlocksNumbersChan()

	blocksChan := p.fetchBlocks(blockNumbersChan)

	transactionsChan := p.processBlocks(blocksChan)

	p.saveTransactions(transactionsChan)

}

func (p *Parser) getBlocksNumbersChan() chan int64 {
	blocksNumbersChan := make(chan int64, 10)

	ticker := time.NewTicker(time.Duration(p.config.IntervalInSecs) * time.Second)

	go func() {
		for ; true; <-ticker.C {
			func() {
				lastBlockInBlockchain, err := p.provider.GetLastBlockNumber()
				if err != nil {
					fmt.Println(err)
					return
				}

				if lastBlockInBlockchain <= p.lastParsedBlockNumber {
					return
				}

				for i := p.lastParsedBlockNumber; i < lastBlockInBlockchain; i++ {
					blocksNumbersChan <- i
					p.lastParsedBlockNumber = i // todo not quite correct
				}
			}()
		}
	}()

	return blocksNumbersChan
}

func (p *Parser) fetchBlocks(blocksNumbersChan chan int64) chan entity.Block {
	blocksChan := make(chan entity.Block, 10)

	fmt.Println(1) //todo remove

	for i := 0; i < p.config.NumberOfFetchingWorkers; i++ {
		go func() {
			for blocksNumbers := range blocksNumbersChan {

				func() {
					block, err := p.provider.GetBlockByNumber(blocksNumbers, true)
					if err != nil {
						fmt.Println(err)
						return
					}

					blocksChan <- block
				}()

			}
		}()
	}

	return blocksChan
}

type TransactionWithAddress struct {
	Transaction entity.Transaction
	Address     string
}

func (p *Parser) processBlocks(blocksChan chan entity.Block) chan TransactionWithAddress {
	transactionsChan := make(chan TransactionWithAddress, 10)

	for i := 0; i < p.config.NumberOfProcessingWorkers; i++ {
		go func() {
			for block := range blocksChan {
				for _, transaction := range block.Transactions {
					for address, _ := range p.subscribedAddressesSet {
						if transaction.From == address ||
							transaction.To == address {
							transactionsChan <- TransactionWithAddress{
								Transaction: transaction,
								Address:     address,
							}
						}
					}
				}

			}
		}()
	}

	return transactionsChan
}

func (p *Parser) saveTransactions(transactionsChan chan TransactionWithAddress) {
	for i := 0; i < p.config.NumberOfSavingWorkers; i++ {
		go func() {
			for transactionWithAddress := range transactionsChan {
				err := p.storage.SaveTransaction(transactionWithAddress.Address, transactionWithAddress.Transaction)
				if err != nil {
					return
				}
			}
		}()
	}
}
