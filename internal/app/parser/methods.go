package parser

import (
	"context"
	"fmt"
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"sync"
	"time"
)

func (p *Parser) Start(ctx context.Context) {
	blockNumbersChan := p.getBlocksNumbersChan(ctx)

	blocksChan := p.fetchBlocks(blockNumbersChan)

	transactionsChan := p.processBlocks(blocksChan)

	p.saveTransactions(transactionsChan)

}

func (p *Parser) getBlocksNumbersChan(ctx context.Context) chan int64 {
	blocksNumbersChan := make(chan int64, 10)

	ticker := time.NewTicker(time.Duration(p.config.IntervalInSecs) * time.Second)

	go func() {
		for ; true; <-ticker.C {
			lastBlockInBlockchain, err := p.provider.GetLastBlockNumber()
			if err != nil {
				fmt.Println(err)
				return
			}

			if lastBlockInBlockchain <= p.lastParsedBlockNumber {
				continue
			}

			for i := p.lastParsedBlockNumber; i < lastBlockInBlockchain; i++ {
				select {
				case <-ctx.Done():
					close(blocksNumbersChan)
					fmt.Println("Shutdown getBlocksNumbersChan")

					return

				default:

				}

				blocksNumbersChan <- i
				p.lastParsedBlockNumber = i
			}
		}
	}()

	return blocksNumbersChan
}

func (p *Parser) fetchBlocks(blocksNumbersChan chan int64) chan entity.Block {
	blocksChan := make(chan entity.Block, 10)

	var wg sync.WaitGroup

	for i := 0; i < p.config.NumberOfFetchingWorkers; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			for blocksNumbers := range blocksNumbersChan {
				block, err := p.provider.GetBlockByNumber(blocksNumbers, true)
				if err != nil {
					fmt.Println(err)
					continue
				}

				blocksChan <- block
			}
		}()
	}

	go func() {
		wg.Wait()
		close(blocksChan)
		fmt.Println("Shutdown fetchBlocks")
	}()

	return blocksChan
}

type TransactionWithAddress struct {
	Transaction entity.Transaction
	Address     string
}

func (p *Parser) processBlocks(blocksChan chan entity.Block) chan TransactionWithAddress {
	transactionsChan := make(chan TransactionWithAddress, 10)

	var wg sync.WaitGroup

	for i := 0; i < p.config.NumberOfProcessingWorkers; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			for block := range blocksChan {
				for _, transaction := range block.Transactions {
					for address := range p.subscribedAddressesSet {
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

	go func() {
		wg.Wait()
		close(transactionsChan)
		fmt.Println("Shutdown processBlocks")
	}()

	return transactionsChan
}

func (p *Parser) saveTransactions(transactionsChan chan TransactionWithAddress) {
	var wg sync.WaitGroup

	for i := 0; i < p.config.NumberOfSavingWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for transactionWithAddress := range transactionsChan {
				err := p.storage.SaveTransaction(transactionWithAddress.Address, transactionWithAddress.Transaction)
				if err != nil {
					return
				}
			}
		}()
	}

	wg.Wait()
	fmt.Println("Shutdown saveTransactions")
}
