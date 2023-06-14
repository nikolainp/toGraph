package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_processFile(t *testing.T) {
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
			`<html>
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>My page</title>

	<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
	<script type='text/javascript'>
		google.charts.load('current', {'packages':['annotationchart']});
		google.charts.setOnLoadCallback(drawChart);

		function drawChart() {
		var data = new google.visualization.DataTable();
		data.addColumn('date', 'Date');
		data.addColumn('number', 'Column 1');
		data.addColumn('number', 'Column 2');
		data.addColumn('number', 'Column 3');
		
		    // [new Date(2314, 2, 16), 24045, 12374],
		
		data.addRows([
			 [new Date(2012, 10, 15, 10, 00, 00), 1, 2, 3]
			,[new Date(2012, 10, 15, 10, 00, 30), 2, 3, 4]
			,[new Date(2012, 10, 15, 10, 01, 00), 1, 2, 3]
			,[new Date(2012, 10, 15, 10, 01, 30), 2, 3, 4]
			,[new Date(2012, 10, 15, 10, 02, 00), 1, 2, 3]
		]);

		var chart = new google.visualization.AnnotationChart(document.getElementById('chart_div'));

		var options = {
			displayAnnotations: true
		};

		chart.draw(data, options);
		}
	</script>
	</head>

	<body>
	<div id='chart_div' style='width: 100%; height: 100%;'></div>
	</body>
</html>
`,
			false,
		},
	}

	replacer := strings.NewReplacer(" ", "", "\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sIn := strings.NewReader(tt.sInput)
			sOut := &bytes.Buffer{}
			if err := processFile(sIn, sOut); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSOut := sOut.String(); replacer.Replace(gotSOut) != replacer.Replace(tt.wantSOut) {
				t.Errorf("run() = %v, want %v", gotSOut, tt.wantSOut)
			}
		})
	}
}
