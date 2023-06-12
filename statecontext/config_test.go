package statecontext

import (
	"reflect"
	"testing"
)

func Test_readCommandLineArguments(t *testing.T) {
	type want struct {
		config Configuration
		err    error
	}
	tests := []struct {
		name string
		args []string
		want want
	}{
		{"empty argument list", []string{},
			want{Configuration{printUsage: false}, errEmptyArgumentList}},
		{"test 1", []string{"programname"},
			want{Configuration{programName: "programname", printUsage: true}, nil}},
		{"help", []string{"programname", "-h"},
			want{Configuration{programName: "programname", printUsage: true}, nil}},
		{"test 2", []string{"programname", "what"},
			want{Configuration{InputFiles: []string{"what"}, programName: "programname"}, nil}},
		{"test 3", []string{"programname", "what", "where"},
			want{Configuration{InputFiles: []string{"what", "where"}, programName: "programname"}, nil}},
		{"test 4", []string{"programname", "newLine", "what", "where"},
			want{Configuration{InputFiles: []string{"newLine", "what", "where"}, programName: "programname"}, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf Configuration
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
