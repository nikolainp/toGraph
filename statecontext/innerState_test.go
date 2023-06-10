package statecontext

import (
	"reflect"
	"testing"
)

func Test_initState(t *testing.T) {
	tests := []struct {
		name string
		want *state
	}{
		{
			"init all",
			new(state),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initState(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("initState() = %v, want init", got)
			}
		})
		//syscall.Kill(syscall.Getpid(), syscall.SIGINT) // stop waiting signal
	}
}
