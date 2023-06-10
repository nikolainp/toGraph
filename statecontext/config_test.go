package statecontext

import (
	"reflect"
	"testing"
)

func Test_readCommandLineArguments(t *testing.T) {
	type want struct {
		config Config
		err    error
	}
	tests := []struct {
		name string
		args []string
		want want
	}{
		{"test 1", []string{},
			want{Config{printUsage: true}, nil}},
		{"help", []string{"-h"},
			want{Config{printUsage: true}, nil}},
		{"test 2", []string{"what"},
			want{Config{NewLineRegex: "", SearchLineRegex: "what", SearchPath: ""}, nil}},
		{"test 3", []string{"what", "where"},
			want{Config{NewLineRegex: "", SearchLineRegex: "what", SearchPath: "where"}, nil}},
		{"newLine", []string{"-l", "newLine", "what", "where"},
			want{Config{NewLineRegex: "newLine", SearchLineRegex: "what", SearchPath: "where"}, nil}},
		{"logFile", []string{"-o", "logFile", "what", "where"},
			want{Config{SearchLineRegex: "what", SearchPath: "where", LogOutputPath: "logFile"}, nil}},
		{"test 6", []string{"arg1", "arg2", "arg3"},
			want{Config{printUsage: true}, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf Config
			_, err := readCommandLineArguments(&conf, tt.args)
			if err != tt.want.err {
				t.Errorf("configure() error = %v, wantErr %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(conf, tt.want.config) {
				t.Errorf("configure() config = %v, want %v", conf, tt.want.config)
			}
		})
	}
}
