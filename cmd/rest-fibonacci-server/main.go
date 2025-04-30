package main

import (
	restfibonacciserver "github.com/cristianortiz/fibonacci-app/apps/rest-fibonacci-server"
)

func main() {
	app := restfibonacciserver.NewApp()
	app.Start()
}
