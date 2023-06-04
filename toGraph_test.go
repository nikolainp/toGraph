package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name     string
		sInput   string
		wantSOut string
		wantErr  bool
	}{
		{"Test 1",
			`20121015100000 1 2 3
20121015100030 2 3 4
20121015100100 1 2 3
20121015100130 2 3 4
20121015100200 1 2 3`,
			`{
"hostname": "dummy.rtp.raleigh.ibm.com",
"whenPattern": "MM/dd/yyyy HH:mm:ss",
"metadata": {
	"mdata1": "value 1",
	"mdata2": "value 2",
	"mdata3": "completely arbitrary data goes here"
},
"types": [
	{
	"id": "type1",
	"name": "type 1",
	"fields": [ "field11", "field12", "field13" ]
	}
],
"data": [
	{ "when": "10/15/2012 10:00:00",
	"type1": [ 1, 2, 3 ]
	},
	{ "when": "10/15/2012 10:00:30",
	"type1": [ 2, 3, 4 ]
	},
	{ "when": "10/15/2012 10:01:00",
	"type1": [ 1, 2, 3 ]
	},
	{ "when": "10/15/2012 10:01:30",
	"type1": [ 2, 3, 4 ]
	},
	{ "when": "10/15/2012 10:02:00",
	"type1": [ 1, 2, 3 ]
	}
]
}
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sIn := strings.NewReader(tt.sInput)
			sOut := &bytes.Buffer{}
			if err := run(sIn, sOut); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSOut := sOut.String(); gotSOut != tt.wantSOut {
				t.Errorf("run() = %v, want %v", gotSOut, tt.wantSOut)
			}
		})
	}
}
