package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// OnCloseSignal call f on receiving close signal
func OnCloseSignal(f func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-ch
		f()
	}()
}

// CancelOnExitContext return a context , call cancel func on exit
func CancelOnExitContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	OnCloseSignal(cancel)
	return ctx
}
