package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const preamble = `
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)
`

const preambleWithHTTPServer = `
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)
`

const root = `
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
`

const mainE = `
func mainE(ctx context.Context) error {
	return nil
}
`

const mainEWithHTTPServer = `
func mainE(ctx context.Context) error {
	s := http.Server{Addr: "localhost:8080"}
	var eg errgroup.Group
	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return s.Shutdown(ctx)
	})

	switch err := s.ListenAndServe(); err {
	case http.ErrServerClosed:
		return eg.Wait()
	default:
		return err
	}
}
`

func main() {
	var moduleName string
	var httpServer bool
	flag.StringVar(&moduleName, "module-name", "", "create a go module with this name")
	flag.BoolVar(&httpServer, "http-server", false, "include net/http.Server set-up")
	flag.Parse()

	if _, err := os.Stat("main.go"); !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "main.go already exists", err)
		os.Exit(1)
	}

	if moduleName != "" {
		cmd := exec.Command("go", "mod", "init", moduleName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error executing %s: %s\n", cmd.String(), err.Error())
			os.Exit(1)
		}
	}

	if httpServer {
		cmd := exec.Command("go", "get", "golang.org/x/sync")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error executing %s: %s\n", cmd.String(), err.Error())
			os.Exit(1)
		}
	}

	var content []byte
	switch {
	case httpServer:
		content = []byte(preambleWithHTTPServer)
		content = append(content, []byte(root)...)
		content = append(content, []byte(mainEWithHTTPServer)...)
	default:
		content = []byte(preamble)
		content = append(content, []byte(root)...)
		content = append(content, []byte(mainE)...)
	}

	if err := os.WriteFile("main.go", bytes.TrimSpace([]byte(content)), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
