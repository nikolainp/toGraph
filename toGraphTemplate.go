package main

//		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}

const graphTemplate = `
<html>
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <style>
        .dropdown {
          position: relative;
          display: inline-block;
        }
        
        .dropdown-content {
          display: none;
          position: absolute;
          background-color: #f9f9f9;
          min-width: 160px;
          box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
          padding: 12px 16px;
          z-index: 1;
        }
        
        .dropdown:hover .dropdown-content {
          display: block;
        }
    </style>

    <title>{{.Title}}</title>

    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type='text/javascript'>
      google.charts.load('current', {'packages':['annotationchart']});
      google.charts.setOnLoadCallback(drawChart);

      let chart;

      function drawChart() {
        var data = new google.visualization.DataTable();
        data.addColumn('date', 'Date');
        {{range $column := .Columns -}}
        data.addColumn('number', '{{$column}}');
        {{end}}
		    // [new Date(2314, 2, 16), 24045, 12374],
        
        data.addRows([
			{{- range $index, $dataRow := .DataRows -}}
			  {{- if (eq $index 0)}}
           {{$dataRow}}
        {{- else}}
          ,{{$dataRow}}
        {{- end}}
			{{- end}}
        ]);

        chart = new google.visualization.AnnotationChart(document.getElementById('chart_div'));

        var options = {
          displayAnnotations: true
        };

        chart.draw(data, options);
      }

      function onChangeColumn(element) {
        if (element.checked) {
            chart.showDataColumns(element.value)
        } else {
            chart.hideDataColumns(element.value)
        }
      }
    </script>
  </head>

  <body>
    <div class="dropdown" style='height: 30px;'>
      <span>My page</span>
      <div class="dropdown-content">
        {{range $i, $column := .Columns -}}
          <input type="checkbox" name="{{$column}}" value="{{$i}}" onchange="onChangeColumn(this)" checked>
            <label for="{{$column}}">{{$column}}</label><br>
        {{end}}
      </div>
    </div>

    <div id='chart_div' style='width: 100%; height: calc(100% - 30px);'></div>
  </body>
</html>
`
