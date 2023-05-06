package main

import (
	"bytes"
	"io"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		sIn io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantSOut string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sOut := &bytes.Buffer{}
			if err := run(tt.args.sIn, sOut); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSOut := sOut.String(); gotSOut != tt.wantSOut {
				t.Errorf("run() = %v, want %v", gotSOut, tt.wantSOut)
			}
		})
	}
}
