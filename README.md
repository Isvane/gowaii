# 🙇‍♂️ Gomen

This repository is my first introduction in learning the Go programming language.

What I learned so far:
* **Routing & Path Parameters:** I learned how to use Go's standard `net/http` router (`http.NewServeMux`) to match HTTP methods (`GET`, `POST`, `PUT`, `DELETE`) and pull path parameters directly using `r.PathValue()`.
* **Thread Safety:** Standard Go maps aren't safe for concurrent writes. I protected in-memory user data using a readers-writer mutex (`sync.RWMutex`) with read (`RLock`) and write (`Lock`) locks.
* **JSON:** I learned how to decode request payloads (`json.NewDecoder`) into Go structs using struct tags (`json:"name"`), handle validation, and stream JSON responses using `json.NewEncoder`.
* **HTTP Middleware:** I learned how the HTTP request lifecycle works by wrapping the router in a custom middleware function (`logMiddleware`) that logs incoming requests before passing them to the handlers.
* **Server Configuration & Hardening:** Instead of using default `http.ListenAndServe`, I configured an `http.Server` struct with explicit `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes` limits to prevent slow-loris attacks and resource leaks.
* **Project Structure:** I learned how to clean up my code by moving from a single `main.go` into a standard layout to practice SoC.
