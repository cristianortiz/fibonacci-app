# gRPC Server Side Streaming App Fibonacci Sequence

Building a Fibonacci service that streams the Fibonacci sequence from the server to the client, using two approaches:

- REST Polling service
- gRPC Server-Side Streaming

## REST Fibonacci Server

The `apps/rest-fibonacci-server` folder contains the implementation of a REST-based Fibonacci service. This service provides two endpoints for calculating Fibonacci sequences:

1. **Synchronous Fibonacci Calculation** (`/fibonacci/sync/{number}`): Calculates the Fibonacci sequence up to the given number synchronously and returns the result along with the time taken for the computation.
2. **Asynchronous Fibonacci Calculation** (`/fibonacci/async/{number}`): Demonstrates REST polling by calculating the Fibonacci sequence asynchronously. The server streams partial results to the client as they become available.

### What is the Fibonacci Sequence?

The Fibonacci sequence is a series of numbers where each number is the sum of the two preceding ones, starting from 0 and 1. Mathematically, it is defined as:

- `F(0) = 0`
- `F(1) = 1`
- `F(n) = F(n-1) + F(n-2)` for `n > 1`

For example, the first few numbers in the Fibonacci sequence are: `0, 1, 1, 2, 3, 5, 8, 13, 21, ...`.

This sequence is often used as a programming exercise because it can be implemented in various ways (e.g., recursion, iteration, dynamic programming) and demonstrates concepts like performance optimization and algorithm design.

### What is REST Polling?

REST polling is a technique used in REST APIs to handle long-running operations. Instead of waiting for the server to complete a task and return the result in a single response, the client repeatedly sends requests to the server to check the status of the operation or retrieve partial results. This approach is useful when:

- The operation takes a long time to complete.
- The client needs to display progress or intermediate results.

In this application, REST polling is demonstrated by the asynchronous Fibonacci endpoint (`/fibonacci/async/{number}`), where the server streams partial results to the client as they are calculated.

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

#### `async-handler.go`

This file implements the asynchronous Fibonacci calculation endpoint (`/fibonacci/async/{number}`).

Key points:

- The `fibonacciAsyncHandler` is the core of the asynchronous Fibonacci calculation. It performs the following steps:

  1. **Extracts and validates the input**: The number to calculate is extracted from the URL, converted to an integer, and validated. If invalid, a `400 Bad Request` error is returned.
  2. **Manages the `request-id`**: The `request-id` is extracted from the HTTP headers to uniquely identify the client's request. If missing, a `400 Bad Request` error is returned.
  3. **Creates or retrieves an `AsyncStore`**: If no `AsyncStore` exists for the given `request-id`, a new one is created. A goroutine is launched to perform the Fibonacci calculations asynchronously.
  4. **Reads partial results**: The server reads the numbers calculated so far from the `AsyncStore`. It also retrieves the current index and the target range.
  5. **Determines completion**: If the current index matches the target range, the calculation is complete, and the `AsyncStore` is deleted.
  6. **Builds and sends the response**: A JSON response is constructed with the calculated numbers, the `request-id`, and a flag indicating whether the calculation is complete. This response is sent back to the client.

- The `AsyncStore` struct is used to manage the state of asynchronous calculations, including the current progress and the calculated numbers. It ensures thread safety using a mutex.
- The `fibAsync` method is executed in a goroutine to perform the Fibonacci calculations concurrently. It writes the results to the `AsyncStore` as they are computed.

### Additional Notes

- The `generateRequestIDMiddleware` in `rest-fibonacci-server.go` ensures that each request has a unique `request-id`, which is critical for managing asynchronous calculations.
- The `fib` function in `sync-handler.go` is a simple recursive implementation of the Fibonacci sequence. While functional, it is not optimized for large inputs and may cause performance issues for high numbers.
- The asynchronous implementation demonstrates how to handle long-running computations in a REST API by using polling and partial responses.

This REST service is designed to showcase the differences between synchronous and asynchronous approaches to handling computationally intensive tasks like generating Fibonacci sequences.
