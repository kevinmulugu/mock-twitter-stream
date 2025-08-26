# Mock Twitter (X) Streaming API

This repository contains a small Go application that simulates the historical Twitter Streaming API (the `POST /1.1/statuses/filter.json` endpoint). It is intended for learning, local testing, and development when the real API is unavailable or changed.

> Note: The project was developed and tested with Go `go1.23.8` on `darwin/amd64`, but the code is cross-platform and should build on other OS/architectures that support the Go toolchain.

## Overview

The server accepts a `POST` request to `/1.1/statuses/filter.json` with a form field `track` containing a comma-separated list of keywords. It then continuously streams simulated tweet objects (JSON) back to the client — one JSON object per line — until the client closes the connection or the server is shut down.

This mock is intentionally simple and designed for educational use and local development only.

## Features

* Simulates a streaming endpoint that emits newline-delimited JSON tweet objects.
* Accepts comma-separated `track` keywords the same way the historical API did.
* Implements graceful shutdown on `SIGINT`/`SIGTERM`.
* Minimal, dependency-free implementation in Go.

## Requirements

* Go 1.23 or later (tested with `go1.23.8`).
* A terminal / environment capable of making HTTP requests (e.g. `curl`).

## Build

From the repository root:

```bash
# build binary
go build -o mock-twitter-stream
```

Or run directly without building:

```bash
go run .
```

## Run

Default address is `:8080`.

```bash
# run built binary (default)
./mock-twitter-stream

# or run on a specific address e.g. :9000
./mock-twitter-stream -addr :9000

# run with go run
go run . -addr :8080
```

The server logs startup information and will handle `SIGINT` / `SIGTERM` to perform a graceful shutdown with a short timeout.

## API

**Endpoint**

```
POST /1.1/statuses/filter.json
```

**Form parameters**

* `track` (required): comma-separated list of keywords to simulate tracking. Example: `track=golang,dev,open-source`

**Response**

* `Content-Type: application/json`
* The server streams newline-delimited JSON objects. Each object corresponds to a simplified tweet in the form:

```json
{ "text": "Someone just mentioned <keyword>" }
```

The stream does not terminate on its own; it will continue until the client closes the connection or the server shuts down.

## Examples

### Using `curl` (recommended for streaming)

Use the `-N` option to disable buffering so you can see streamed lines as they arrive.

```bash
curl -N -X POST -d "track=golang,dev" http://localhost:8080/1.1/statuses/filter.json
```

You should see a continuous stream of JSON objects, one per line. Press `Ctrl+C` to close the connection from your client.

### Example client in Go (simple)

```go
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	resp, err := http.PostForm("http://localhost:8080/1.1/statuses/filter.json", url.Values{"track": {"golang,dev"}})
	if err != nil {
		fmt.Println("request error:", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
	}
}
```

This client posts the `track` parameter and reads the streaming response line by line.

## Configuration

* `-addr` : HTTP listen address (default `:8080`).

## Behavior and limitations

* The server fakes tweets: each emitted object is a small JSON with a `text` field.
* The stream is infinite by design; it will only stop when the client disconnects or the server is shut down.
* There is no authentication, rate-limiting, or persistence in this mock implementation.
* The code is intended for local development and learning. Do not use this as a production replacement for an official API.

## Dockerfile

A simple Dockerfile is included for containerization:

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o mock-twitter-stream

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/mock-twitter-stream /usr/local/bin/mock-twitter-stream
EXPOSE 8080
CMD ["/usr/local/bin/mock-twitter-stream", "-addr", ":8080"]
```

### Build and run with Docker

```bash
# build the docker image
docker build -t mock-twitter-stream .

# run the container
docker run -p 8080:8080 mock-twitter-stream
```

Once running, the mock server is accessible at `http://localhost:8080/1.1/statuses/filter.json`.

## Contributing

Contributions and improvements are welcome. Suggested areas:

* Add configurable output formats (e.g. full tweet fields).
* Add artificial delays, jitter, or distribution that better mimic real-world traffic.
* Add TLS support and basic authentication for private testing environments.

## License

Specify your preferred license here (e.g., MIT, Apache-2.0). If you have no preference, consider adding an `LICENSE` file.
