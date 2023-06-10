package statecontext

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	NewLineRegex    string
	SearchLineRegex string
	SearchPath      string
	LogOutputPath   string

	printUsage      bool
}

type state struct {
	ctx    context.Context
	cancel context.CancelFunc

	signChan  chan os.Signal
	workersWg sync.WaitGroup

	config Config
}

var once sync.Once
var innerState *state = nil

// Public

func InitState() {
	once.Do(func() {
		innerState = initState()
	})
}

func Configure(args ...string) (*Config) {
	InitState()
	innerState.config.configure(args...)

	return &innerState.config
}

func Go(work func()) {
	innerState.workersWg.Add(1)
	go func() {
		defer innerState.workersWg.Done()

		work()
	}()
}

func Wait() {
	innerState.workersWg.Wait()
}

func Done() <-chan struct{} {
	return innerState.ctx.Done()
}

func IsDone() bool {
	select {
    case _, ok := <- Done():
        return !ok
    default:
        return false
    }
}

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		innerState.cancel()
	}
}
