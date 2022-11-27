package main

import (
	"context"
	"fmt"
	"github.com/Rostislaved/ethereum-parser/internal/app/adapters/httpAdapter"
	"github.com/Rostislaved/ethereum-parser/internal/app/config"
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"github.com/Rostislaved/ethereum-parser/internal/app/parser"
	"github.com/Rostislaved/ethereum-parser/internal/app/provider"
	"github.com/Rostislaved/ethereum-parser/internal/app/storage/inmemory_storage"
	"github.com/Rostislaved/ethereum-parser/internal/pkg/hexconverter"
	signalListener "github.com/Rostislaved/ethereum-parser/internal/pkg/signal-listener"
	"log"
	"sync"
)

var _ Parser = (*parser.Parser)(nil) // parser implements the interface

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []entity.Transaction
}

func main() {
	config := config.Get()

	hexconverter := hexconverter.New()

	prv := provider.New(config.Provider, hexconverter)

	storage := inMemoryStorage.New()

	parser := parser.New(config.Parser, storage, prv)

	httpAdapter := httpAdapter.New(config.Server, parser)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		parser.Start(ctx)

		fmt.Println("Shutdown parser")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpAdapter.Start()

		fmt.Println("Shutdown web server")
	}()

	signalListener := signalListener.New()

	select {
	case err := <-httpAdapter.Notify():
		log.Println(err)

	case signal := <-signalListener.Notify():
		log.Printf("\nCaught signal: %v. Exiting...\n", signal)

		err := httpAdapter.Shutdown()
		if err != nil {
			log.Println(err)
		}
	}

	cancel()

	wg.Wait()
}
