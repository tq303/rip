package state

import (
	"context"
	"os"
	"os/signal"
)

var C context.Context
var cancel context.CancelFunc

func init() {
	C, cancel = context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()
}

func Get() context.Context {
	return C
}
