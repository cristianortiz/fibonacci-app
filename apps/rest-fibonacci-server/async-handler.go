package restfibonacciserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type AsyncResponse struct {
	RequestId        string `json:"requestid"`
	FibonacciNumbers []int  `json:"fibonacciNumbers"`
	EndOfResponse    bool   `json:"endOfResponse"`
}

type AsyncStore struct {
	mu             sync.Mutex
	current        int
	requestedRange int
	numbers        []int
}

func NewAsyncStore(requestedRange int) *AsyncStore {
	return &AsyncStore{requestedRange: requestedRange}
}

func (ns *AsyncStore) Write(number, current int) {
	ns.mu.Lock()
	ns.current = current
	ns.numbers = append(ns.numbers, number)
	defer ns.mu.Unlock()
}

func (ns *AsyncStore) Read() ([]int, int, int) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	readNumbers := make([]int, len(ns.numbers))
	copy(readNumbers, ns.numbers)
	ns.numbers = []int{}
	return readNumbers, ns.current, ns.requestedRange
}

// calculates async fibo series based on client reuqest, the server also streams the results to clients as the
// becomes avalaible
func (a *App) fibonacciAsyncHandler(w http.ResponseWriter, r *http.Request) {
	//extract the value of the number from the request and perform error checking
	vars := mux.Vars(r)
	number := vars["number"]
	numFibonacci, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	headers := r.Header
	//extract request-id from http headers, helps to maintain the state of the ongoing calculations
	//for each client
	reqId := headers.Get("request-id")
	if strings.TrimSpace(reqId) == "" {
		http.Error(w, "no request id in request", http.StatusBadRequest)
		return
	}
	//checks if an AsyncStore object exists for the given request-id
	//if not, creates a new one
	if _, ok := a.asyncStores[reqId]; !ok {
		fmt.Println("creating new store for reqId", reqId)
		a.asyncStores[reqId] = NewAsyncStore(numFibonacci)
		//go routine launched to perform the fibo calculations concurrently
		go a.fibAsync(numFibonacci, reqId)
	}
	//reads results from AsyncStore, numbersNow, is already computed numbers, current  is last fibo number computed
	// , 'requested' target range
	numbersNow, current, requested := a.asyncStores[reqId].Read()
	fmt.Printf("read fibs reqId %s till current %d and numbers are: %v\n", reqId, current, numbersNow)
	end := false
	//if current number equals tje requested change, the fibo calculation is complete and endOfResponse will be
	//set to true in response object
	if current == requested {
		end = true
		delete(a.asyncStores, reqId)
	}

	response := AsyncResponse{
		RequestId:        reqId,
		FibonacciNumbers: numbersNow,
		EndOfResponse:    end,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}
func (a *App) fibAsync(n int, reqId string) {
	for i := range n {
		fmt.Printf("for %s computing and writing fib of %d\n", reqId, i)
		a.asyncStores[reqId].Write(fib(i), i)
	}

}
