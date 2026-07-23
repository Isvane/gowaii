# Gomen

This repository is my first introduction in learning the Go programming language.

## Quick Start

```console
go run cmd/server/main.go
```

## Project Structure

```text
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── internal
│   ├── api
│   │   ├── handlers.go
│   │   └── middleware.go
│   ├── model
│   │   └── models.go
│   └── repository
│       └── user.go
└── README.md
```

## API Endpoints

| Method | Endpoint | Description | Sample Payload |
|--------|----------|-------------|----------------|
| GET | `/echo/{arg}` | Echoes back the path argument. | — |
| POST | `/user` | Creates a new user. | `{"name": "alice", "age": 25}` |
| GET | `/user/{name}` | Retrieves user info by name. | — |
| PUT | `/user/{name}` | Updates or renames an existing user. | `{"name": "alice_renamed", "age": 26}` |
| DELETE | `/user/{name}` | Removes a user by name. | — |

## Learning Takeaway

| Concept | Implementation |
| :--- | :--- |
| **Routing** | Go 1.22 stdlib routing and path params (`r.PathValue()`) |
| **Concurrency** | Thread-safe state management via `sync.RWMutex` |
| **JSON** | Encoding/decoding structs with `json` tags |
| **Middleware** | Custom HTTP logging wrapper with `slog` |
| **Hardening** | Read/write timeouts and OS signal-based graceful shutdown |
| **Architecture** | Standard `cmd/` vs `internal/` package separation |
| **Testing** | How to test a Go handlers |
