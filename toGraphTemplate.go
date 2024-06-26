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
        var columns = new google.visualization.DataTable();
        columns.addColumn('string', 'Title');
        columns.addColumn('number', 'Minimum');
        columns.addColumn('number', 'Maximum');
        columns.addColumn('number', 'Average');
        columns.addColumn('boolean', 'Show');

        columns.addRows([
        {{range $i, $column := .Columns -}}
          ['{{$column.Name}}', {{$column.Minimum}}, {{$column.Maximum}}, {{$column.Average}}, true],
        {{end}}
        ]);


        var data = new google.visualization.DataTable();
        data.addColumn('date', 'Date');
        {{range $column := .Columns -}}
        data.addColumn('number', '{{$column.Name}}');
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

        table = new google.visualization.Table(document.getElementById('table_div'));
        chart = new google.visualization.AnnotationChart(document.getElementById('chart_div'));

        chart.draw(data, {displayAnnotations: true});
        setColumnColors(columns);
        drawTable(table, columns)

        google.visualization.events.addListener(table, 'select', function() {
          var row = table.getSelection()[0].row;
          var value = columns.getValue(row, 4);

          value = !value;
          columns.setValue(row, 4, value);
          if (value) {
              chart.showDataColumns(row)
          } else {
              chart.hideDataColumns(row)
          }

          drawTable(table, columns);
        });

      }

      function setColumnColors(table) {
        var legends = document.getElementsByClassName("legend-dot");

        for (let i = 0; i < legends.length; i++) {
            table.setProperties(i, 0, {'style': 'color: ' + legends[i].style.backgroundColor +'; white-space: nowrap;'});
            table.setProperties(i, 4, {'style': 'color: ' + legends[i].style.backgroundColor +';'});
        }
      }

      function drawTable(table, data) {
        table.draw(data, {showRowNumber: true, allowHtml: true, width: '100%', height: '100%'});
      }

      </script>
  </head>

  <body>
    <div class="dropdown" style='height: 30px;'>
      <span>{{.Title}}</span>
      <div class="dropdown-content">
        <div id='table_div'></div>
      </div>
    </div>

    <div id='chart_div' style='width: 100%; height: calc(100% - 30px);'></div>
  </body>
</html>
`
