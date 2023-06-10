package statecontext

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func initState() *state {
	obj := new(state)

	obj.signChan = make(chan os.Signal, 1)

	obj.ctx, obj.cancel = context.WithCancel(context.Background())

	signal.Notify(obj.signChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go listenSignal()

	return obj
}

func listenSignal() {
	sig := <-innerState.signChan

	fmt.Printf("%v signal received", sig)

	innerState.cancel()
}
