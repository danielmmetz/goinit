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
