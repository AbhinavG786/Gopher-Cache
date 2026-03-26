# GopherCache

> A small, fast, in-memory cache and pub/sub server written in Go. Inspired by Redis, GopherCache focuses on simplicity and clarity for learning and lightweight caching workloads.

## Features

- In-memory key-value store with simple commands
- Append-only file (AOF) persistence for durability
- Pub/Sub hub for lightweight messaging
- Minimal Go client API
- Small codebase ideal for learning and extension

## Architecture Overview

GopherCache is organized into a few small components (see the repository layout below):

- **Storage/Engine**: core in-memory store and operations are in `internal/engine` (`store.go`).
- **AOF Persistence**: append-only persistence logic lives in `internal/aof` (`aof.go`) and writes incoming commands to `database.aof`.
- **Protocol & Parser**: the server protocol parser is implemented in `internal/protocol` (`parser.go`) to accept a simple text/RESP-like format.
- **Pub/Sub**: `pubsub/hub.go` implements a hub for publish/subscribe messaging between clients.
- **Server**: the server entrypoint is `cmd/server/main.go` which wires the components and exposes a TCP listener.
- **Client**: a minimal Go client is available at `client/client.go` for example usage.

This layout keeps concerns separated and makes it easy to extend each part independently.

## Quickstart — Build & Run

Prerequisites: Go 1.18+ installed.

Build the server:

```bash
go build -o gophercache ./cmd/server
```

Run the server (default: listens on :6380):

```bash
./gophercache
```

You can also run directly with `go run` during development:

```bash
go run ./cmd/server
```

## Using the client

The repository includes a small Go client at `client/client.go`. Example usage:

```go
package main

import (
	"log"
	"time"
	"GopherCache/client"
)

func main() {
	c, err := client.Dial("tcp", ":6380")
	if err != nil { log.Fatal(err) }
	defer c.Close()

	// Set a key
	_, _ = c.Do("SET", "foo", "bar")
	// Get a key
	val, _ := c.Do("GET", "foo")
	log.Println("GET foo ->", val)

	// Simple publish
	_, _ = c.Do("PUBLISH", "channel1", "hello")
	time.Sleep(50 * time.Millisecond)
}
```

If you prefer a raw telnet-like interaction (simple text protocol), connect and type commands:

```bash
nc localhost 6380
SET mykey myvalue
GET mykey
```

## Commands (examples)

- `SET key value` — store a value
- `GET key` — retrieve a value
- `DEL key` — delete a key
- `PUBLISH channel message` — publish a message
- `SUBSCRIBE channel` — subscribe to a channel (client blocks and receives messages)

Refer to the `internal/protocol` and `pubsub` packages for the exact wire semantics.

## Persistence (AOF)

Commands that modify state are appended to `database.aof` by the AOF module in `internal/aof`. On startup the server replays `database.aof` to rebuild in-memory state. AOF provides a simple, crash-resilient persistence mechanism suitable for learning and small deployments.

## Design notes

- The project favors clarity over optimization. Data structures are straightforward and easy to follow.
- Concurrency: the store is safe for concurrent access; see `internal/engine/store.go` for locking strategy.
- Extensibility: new commands, eviction policies, or clustering behavior can be added by extending the engine and protocol parser.

## Running tests

If repository contains tests, run:

```bash
go test ./...
```

## Contributing

Contributions welcome. Open issues for bugs or feature requests, and submit PRs with tests and concise descriptions.

## License

This project is provided as-is for learning and demo purposes. Add a license file if you plan to publish or reuse commercially.
