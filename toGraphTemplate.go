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

        chart = new google.visualization.AnnotationChart(document.getElementById('chart_div'));

        var options = {
          displayAnnotations: true
        };

        chart.draw(data, options);
        setColumnColors();
      }

      function setColumnColors() {
        var legends = document.getElementsByClassName("legend-dot");
        var columns = document.getElementsByClassName("columnLabel");

        for (let i = 0; i < columns.length; i++) {
            columns[i].style.color = legends[i].style.backgroundColor;
        }
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
      <span>{{.Title}}</span>
      <div class="dropdown-content">
        <table>
          <tr><td>&nbsp;</td><td>name</td><td>minimum</td><td>maximum</td><td>average</td></tr>
          {{range $i, $column := .Columns -}}
          <tr>
            <td><input type="checkbox" name="{{$column.Name}}" value="{{$i}}" onchange="onChangeColumn(this)" checked></td>
            <td nowrap><label class="columnLabel" for="{{$column.Name}}">{{$column.Name}}</label></td>
            <td nowrap>{{$column.Minimum}}</td>
            <td nowrap>{{$column.Maximum}}</td>
            <td nowrap>{{$column.Average}}</td>
          </tr>
          {{end}}
        </table>
      </div>
    </div>

    <div id='chart_div' style='width: 100%; height: calc(100% - 30px);'></div>
  </body>
</html>
`
