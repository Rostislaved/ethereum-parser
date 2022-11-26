package main

import (
	"context"
	"github.com/Rostislaved/ethereum-parser/internal/app/adapters/httpAdapter"
	"github.com/Rostislaved/ethereum-parser/internal/app/config"
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"github.com/Rostislaved/ethereum-parser/internal/app/parser"
	"github.com/Rostislaved/ethereum-parser/internal/app/provider"
	"github.com/Rostislaved/ethereum-parser/internal/app/storage/inmemory_storage"
	"github.com/Rostislaved/ethereum-parser/internal/pkg/hexconverter"
	"os"
	"os/signal"
	"syscall"
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
	var initialBlockNumber int64 = 16050110

	cfg := config.Parser{
		IntervalInSecs:            5,
		InitialBlockNumber:        16050110,
		NumberOfFetchingWorkers:   10,
		NumberOfProcessingWorkers: 100,
		NumberOfSavingWorkers:     10,
	}

	hexconverter := hexconverter.New()

	prv := provider.New(hexconverter)

	storage := inMemoryStorage.New()

	parser := parser.New(cfg, storage, prv, initialBlockNumber)

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal)
	signal.Notify(
		signalChan,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-signalChan
		cancel()
	}()

	go parser.Start(ctx)

	httpAdapter := httpAdapter.New(parser)

	httpAdapter.Start()
}
