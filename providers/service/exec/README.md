# Exec

The `exec` service is a wrapper around the Go [os/exec](https://pkg.go.dev/os/exec) package in the standard library.

This implementation will build a simple process table in memory and leverage Linux pipes for reading `stdout` and `stderr`.