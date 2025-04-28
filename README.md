# gRPC Server Side Streaming App Fibonacci Sequence

Building a Fibonacci service that streams the Fibonacci sequence from the server to the client, using two approaches:

- REST Polling service
- gRPC Server-Side Streaming

## REST Fibonacci Server

The `apps/rest-fibonacci-server` folder contains the implementation of a REST-based Fibonacci service. This service provides two endpoints for calculating Fibonacci sequences:

1. **Synchronous Fibonacci Calculation** (`/fibonacci/sync/{number}`): Calculates the Fibonacci sequence up to the given number synchronously and returns the result along with the time taken for the computation.
2. **Asynchronous Fibonacci Calculation** (`/fibonacci/async/{number}`): Demonstrates REST polling by calculating the Fibonacci sequence asynchronously. The server streams partial results to the client as they become available.

### File Descriptions

#### `rest-fibonacci-server.go`

This file contains the main application logic for the REST Fibonacci server. It defines the `App` struct, which manages the server's state, including a map of asynchronous stores (`asyncStores`) for handling ongoing asynchronous Fibonacci calculations.

Key points:

- The `NewApp` function initializes a new `App` instance with an empty `asyncStores` map.
- The `Start` method sets up the REST server using the Gorilla Mux router. It registers middleware for generating request IDs and defines the routes for synchronous and asynchronous Fibonacci calculations.

#### `sync-handler.go`

This file implements the synchronous Fibonacci calculation endpoint (`/fibonacci/sync/{number}`).

Key points:

- The `fibonacciSyncHandler` extracts the number from the URL, validates it, and calculates the Fibonacci sequence up to the given number using a simple recursive function (`fib`).
- The time taken for the computation is measured and included in the response.
- The response is returned as JSON, containing the Fibonacci sequence and the time taken.

#### `async-handler.go`

This file implements the asynchronous Fibonacci calculation endpoint (`/fibonacci/async/{number}`).

Key points:

- The `fibonacciAsyncHandler` extracts the number from the URL and validates it. It also checks for a `request-id` in the HTTP headers to maintain the state of ongoing calculations for each client.
- If no `AsyncStore` exists for the given `request-id`, a new one is created, and a goroutine is launched to perform the Fibonacci calculations asynchronously.
- The server streams partial results to the client by reading from the `AsyncStore`. If the calculation is complete, the `AsyncStore` is deleted.
- The `AsyncStore` struct is used to manage the state of asynchronous calculations, including the current progress and the calculated numbers.

### Additional Notes

- The `generateRequestIDMiddleware` in `rest-fibonacci-server.go` ensures that each request has a unique `request-id`, which is critical for managing asynchronous calculations.
- The `fib` function in `sync-handler.go` is a simple recursive implementation of the Fibonacci sequence. While functional, it is not optimized for large inputs and may cause performance issues for high numbers.
- The asynchronous implementation demonstrates how to handle long-running computations in a REST API by using polling and partial responses.

This REST service is designed to showcase the differences between synchronous and asynchronous approaches to handling computationally intensive tasks like generating Fibonacci sequences.
