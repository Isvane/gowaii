# 🍙 Gohan

This repository is my first introduction in learning the Go programming language.

What I learned so far:
*   I learned how to use Go's standard `net/http` router (`http.NewServeMux`) to pull parameters directly from a URL path using `r.PathValue()`.
*   I learned that standard Go maps aren't safe for simultaneous writes. I fixed this by using a readers-writer mutex (`sync.RWMutex`) to lock and unlock the map when saving or reading data.
*   I figured out how the HTTP lifecycle works by wrapping a standard router inside a custom logging function (`middlewareTest`) that runs code before passing the request along.
*   Instead of just using `http.ListenAndServe`, I learned how to configure a custom `http.Server` struct with explicit timeouts (`ReadTimeout` and `WriteTimeout`) to keep the server stable.
