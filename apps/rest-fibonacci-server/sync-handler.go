package restfibonacciserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type SyncResponse struct {
	TimeTaken        string `json:"timeTaken"`
	FibonacciNumbers []int  `json:"fibonacciNumbers"`
}

// calculates sync fibo series and time taken to do it
func (a *App) fibonacciSyncHandler(w http.ResponseWriter, r *http.Request) {
	//extract the value of the number from the url path using mux
	vars := mux.Vars(r)
	number := vars["number"]
	//converts the mumber to an int and error checking
	numFibonacci, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//init int array whit the given number to perform the fibo series
	fibonacciNumbers := make([]int, numFibonacci)
	//records the time to perform the calculation
	now := time.Now()
	//calculates fibo numbers up to the given number, it iterates through the numbers
	//from 0 to numFibonacci, calling fib func to calculate each fibo number
	for i := range numFibonacci {
		fibonacciNumbers[i] = fib(i)
	}
	timeTaken := time.Since(now).Seconds()
	//response wuth the calcultad fibo series and the time taken
	response := SyncResponse{
		TimeTaken:        fmt.Sprintf("%f seconds", timeTaken),
		FibonacciNumbers: fibonacciNumbers,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// simple recursive implementation of fibonacci sequence calculations
func fib(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}

}
