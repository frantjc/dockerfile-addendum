package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/dockerfile-addendum/pkg/command"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		sigC        = make(chan os.Signal, 1)
		err         error
	)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		cancel()
	}()

	if err = command.NewRoot().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
