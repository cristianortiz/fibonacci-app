package grpcfibonacciserver

import (
	grpcfibonacciserver "github.com/cristianortiz/fibonacci-app/apps/grpc-fibonacci-server"
)

func main() {
	app := grpcfibonacciserver.NewApp()
	app.Start()
}
