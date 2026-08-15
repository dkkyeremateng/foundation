# Foundation

A foundation library of reusable Go packages for building web services.

## Overview

Foundation is a library (not a runnable binary) intended to be imported by
other services. It provides four packages:

- **web** — a small HTTP framework extension built on `httptreemux` and
  OpenTelemetry. Provides `App`, `Handler`, middleware wrapping, JSON
  request/response helpers, request context values (trace ID, timestamps),
  and graceful-shutdown signaling.
- **logger** — a convenience wrapper around `zap` producing a SugaredLogger
  with ISO8601 timestamps writing to stdout.
- **keystore** — an in-memory JWT RSA key store. Supports `New`, `NewMap`,
  and `NewFS` (loads `.pem` files rooted in an `fs.FS`, using the filename
  as the key id).
- **validate** — struct/model validation built on `go-playground/validator`
  with English translations, JSON tag names, `FieldErrors` collection, and
  UUID helpers (`GenerateID`, `CheckID`).

## Packages

### web

The `web` package is the entrypoint for HTTP services. `NewApp` creates an
`App` that wraps an `httptreemux` router with OpenTelemetry tracing, and
`App.Handle` registers routes with optional per-route and app-wide
middleware.

```go
shutdown := make(chan os.Signal, 1)
app := web.NewApp(shutdown, mid.Logger(log))

app.Handle(http.MethodGet, "", "/hello/:name", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
    return web.Respond(ctx, w, map[string]string{"hello": web.Param(r, "name")}, http.StatusOK)
})
```

### logger

`logger.New(service)` returns a `*zap.SugaredLogger` configured for
production output to stdout with human-readable ISO8601 timestamps.

### keystore

`keystore` stores RSA private/public keys by key id (`kid`) for JWT
signing and verification. Keys can be added programmatically (`New`,
`NewMap`, `Add`) or loaded from `.pem` files in an `fs.FS` (`NewFS`).

### validate

`validate.Check` validates a struct against its `validate` tags and returns
a `FieldErrors` value with translated, JSON-tag-named error messages.
`GenerateID` and `CheckID` generate and validate UUIDs.

## Requirements

- Go 1.17+

## Install

```sh
go get github.com/dkkyeremateng/foundation
```

## Usage

A minimal web server using the `web` package:

```go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkkyeremateng/foundation/logger"
	"github.com/dkkyeremateng/foundation/web"
)

type greeting struct {
	Name string `json:"name" validate:"required"`
}

func main() {
	log, err := logger.New("EXAMPLE")
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Middleware runs before/after each handler.
	mw := func(next web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Infow("request", "traceID", web.GetTraceID(ctx), "path", r.URL.Path)
			return next(ctx, w, r)
		}
	}

	app := web.NewApp(shutdown, mw)

	app.Handle(http.MethodPost, "", "/greet", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var g greeting
		if err := web.Decode(r, &g); err != nil {
			return web.Respond(ctx, w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		}
		return web.Respond(ctx, w, map[string]string{"message": "hello " + g.Name}, http.StatusOK)
	})

	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Errorw("server", "error", err)
	}
}
```

## Testing

```sh
go test ./...
```
