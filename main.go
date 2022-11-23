package main

import (
	"bytes"
	"fmt"
	"os"
)

const content = `
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := mainE(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		// Only exit non-zero if our initial context has yet to be canceled.
		// Otherwise it's very likely that the error we're seeing is a result of our attempt at graceful shutdown.
		if ctx.Err() == nil {
			os.Exit(1)
		}
	}
}

func mainE(ctx context.Context) error {
	return nil
}
`

func main() {
	if _, err := os.Stat("main.go"); !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "main.go already exists", err)
		os.Exit(1)
	}
	if err := os.WriteFile("main.go", bytes.TrimSpace([]byte(content)), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
