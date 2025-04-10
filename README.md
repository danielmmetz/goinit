### goinit

A trivial program to generate a starting point for Go programs.

Generated code wires up a `context.Context` that will be canceled on SIGINT or SIGTERM,
and encourages pushing logic into a `mainE(context.Context) error` function so that
you can simply return errors and handle exiting once.

It can additionally create a go.mod for you and wire-up basic net/http server scaffolding.

```
❯ goinit -h
Usage of goinit:
  -http-server
        include net/http.Server set-up
  -module-name string
        create a go module with this name
```

#### Example

For a vanilla run (without specifying `-http-server`), it generates a file like:
```go
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if err := mainE(ctx, logger); err != nil {
		logger.ErrorContext(ctx, "exiting with error", slog.Any("err", err))
		// Only exit non-zero if our initial context has yet to be canceled.
		// Otherwise it's very likely that the error we're seeing is a result of our attempt at graceful shutdown.
		if ctx.Err() == nil {
			os.Exit(1)
		}
	}
}

func mainE(ctx context.Context, _ *slog.Logger) error {
	return nil
}
```
