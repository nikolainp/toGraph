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
			want{Configuration{
				DateFormat: "YYYYMMDDHHmmSS",
				DateColumn: 1,
				Delimiter: " ",
				printUsage: false}, errEmptyArgumentList}},
		{"test 1", []string{"programname"},
			want{Configuration{
				DateFormat:  "YYYYMMDDHHmmSS",
				DateColumn:  1,
				Delimiter: " ",
				programName: "programname",
				printUsage:  true}, nil}},
		{"help", []string{"programname", "-h"},
			want{Configuration{
				DateFormat:  "YYYYMMDDHHmmSS",
				DateColumn:  1,
				Delimiter: " ",
				programName: "programname",
				printUsage:  true}, nil}},
		{"dateFormat", []string{"programname", "-t", "YYMMDDHHmm"},
			want{Configuration{
				DateFormat:  "YYMMDDHHmm",
				DateColumn:  1,
				Delimiter: " ",
				programName: "programname",
				printUsage:  true}, nil}},
		{"dateColumn", []string{"programname", "-tc", "2"},
			want{Configuration{
				DateFormat:  "YYYYMMDDHHmmSS",
				DateColumn:  2,
				Delimiter: " ",
				programName: "programname",
				printUsage:  true}, nil}},
		{"test 2", []string{"programname", "what"},
			want{Configuration{
				InputFiles:  []string{"what"},
				DateFormat:  "YYYYMMDDHHmmSS",
				DateColumn:  1,
				Delimiter: " ",
				programName: "programname"}, nil}},
		{"test 3", []string{"programname", "what", "where"},
			want{Configuration{
				InputFiles:  []string{"what", "where"},
				DateFormat:  "YYYYMMDDHHmmSS",
				DateColumn:  1,
				Delimiter: " ",
				programName: "programname"}, nil}},
		{"test 4", []string{"programname", "newLine", "what", "where"},
			want{Configuration{
				InputFiles:  []string{"newLine", "what", "where"},
				DateFormat:  "YYYYMMDDHHmmSS",
				DateColumn:  1,
				Delimiter: " ",
				programName: "programname"}, nil}},
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

func Test_covertDateFormat(t *testing.T) {
	tests := []struct {
		name             string
		commonDateFormat string
		wantGoDateFormat string
	}{
		{"test 1", "YYYYMMDDHHmmSS", "20060102150405"},
		{"test 2", "YYMMDDHHmmSS", "060102150405"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotGoDateFormat := covertDateFormat(tt.commonDateFormat); gotGoDateFormat != tt.wantGoDateFormat {
				t.Errorf("covertDateFormat() = %v, want %v", gotGoDateFormat, tt.wantGoDateFormat)
			}
		})
	}
}
