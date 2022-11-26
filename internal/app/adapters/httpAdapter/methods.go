package httpAdapter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (a *HttpAdapter) Start() {
	err := http.ListenAndServe(":8080", a.getMux())
	if err != nil {
		log.Fatalln(err)
	}
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

	writer.Write(transactionsJsonBytes)
}

func (a *HttpAdapter) getStorageInfo(writer http.ResponseWriter, request *http.Request) {
	storageInfo := a.parser.GetStorageInfo()

	storageInfoJsonBytes, err := json.Marshal(storageInfo)
	if err != nil {
		return
	}

	writer.Write(storageInfoJsonBytes)
}
