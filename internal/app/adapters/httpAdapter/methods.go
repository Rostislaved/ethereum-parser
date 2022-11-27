package httpAdapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *HttpAdapter) Start() {
	log.Println("Server started")

	err := a.server.ListenAndServe()

	a.notify <- err

	close(a.notify)
}

func (a *HttpAdapter) Notify() <-chan error {
	return a.notify
}
func (a *HttpAdapter) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.ShutdownTimeout)*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (a *HttpAdapter) getMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/get-current-block", a.getCurrentBlock)
	mux.HandleFunc("/subscribe", a.subscribe)
	mux.HandleFunc("/get-transactions", a.getTransactions)
	mux.HandleFunc("/get-storage-info", a.getStorageInfo)

	return mux
}

func (a *HttpAdapter) getCurrentBlock(writer http.ResponseWriter, request *http.Request) {
	currentBlock := a.parser.GetCurrentBlock()

	currentBlockString := strconv.Itoa(currentBlock)

	fmt.Fprintln(writer, currentBlockString)
}

func (a *HttpAdapter) subscribe(writer http.ResponseWriter, request *http.Request) {
	address := request.URL.Query().Get("address")

	result := a.parser.Subscribe(address)

	var resultString string

	if result {
		resultString = "true"
	} else {
		resultString = "false"
	}

	fmt.Fprintln(writer, resultString)
}

func (a *HttpAdapter) getTransactions(writer http.ResponseWriter, request *http.Request) {
	address := request.URL.Query().Get("address")

	transactions := a.parser.GetTransactions(address)

	transactionsJsonBytes, err := json.Marshal(transactions)
	if err != nil {
		return
	}

	_, err = writer.Write(transactionsJsonBytes)
	if err != nil {
		log.Println(err)
	}
}

func (a *HttpAdapter) getStorageInfo(writer http.ResponseWriter, request *http.Request) {
	storageInfo := a.parser.GetStorageInfo()

	storageInfoJsonBytes, err := json.Marshal(storageInfo)
	if err != nil {
		return
	}

	_, err = writer.Write(storageInfoJsonBytes)
	if err != nil {
		log.Println(err)
	}
}
