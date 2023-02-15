package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/dockerfile-addendum/command"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err := command.New().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error())
		stop()
		os.Exit(1)
	}

	stop()
	os.Exit(0)
}
