package datarecord

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ColumnStatistic struct {
	Name    string
	Minimum float64
	Maximum float64
	Average float64
}

type dataReaderData map[time.Time]map[string][]dataPoint

type columnStatistic struct {
	minimum float64
	maximum float64
	sum     float64
	count   int
}

type dataColumns struct {
	names     []string
	statistic map[string][]columnStatistic
}

type dataReader struct {
	dateFormat              string
	dateColumn              int
	pivotColumn             int
	delimiter               []byte
	isColumnNamesInFirstRow bool

	columns dataColumns
	points  int
	data    dataReaderData
}

type dataRecord struct {
	dateTime time.Time
	pivot    string
	points   []dataPoint
}

type dataPoint struct {
	point  float64
	isNull bool
}

var checkErr = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

///////////////////////////////////////////////////////
// dataReader

func GetDataReader() (reader dataReader) {

	reader.columns.initialize()
	reader.data = make(dataReaderData)

	return
}

func (obj *dataReader) WithDateFormat(dateFormat string) *dataReader {
	obj.dateFormat = dateFormat
	return obj
}

func (obj *dataReader) WithDateColumn(column int) *dataReader {
	obj.dateColumn = column
	return obj
}

func (obj *dataReader) WithPivotColumn(column int) *dataReader {
	obj.pivotColumn = column
	return obj
}

func (obj *dataReader) WithDelimiter(delimiter string) *dataReader {
	obj.delimiter = []byte(delimiter)
	return obj
}

func (obj *dataReader) WithColumnNames(data string) *dataReader {
	if len(data) == 0 {
		return obj
	}

	obj.columns.names = strings.Split(data, ",")
	for i := range obj.columns.names {
		obj.columns.names[i] = strings.TrimSpace(obj.columns.names[i])
	}

	return obj
}

func (obj *dataReader) WithColumnNamesInFirstRow(data bool) *dataReader {
	obj.isColumnNamesInFirstRow = data
	return obj
}

func (obj *dataReader) ReadDataRecord(data string) {

	if obj.isColumnNamesInFirstRow {
		obj.isColumnNamesInFirstRow = false
		obj.columns.names = strings.Split(data, string(obj.delimiter))
		return
	}

	record := obj.getDataRecord(data)

	obj.columns.addDataRecord(record)
	if obj.points < len(record.points) {
		obj.points = len(record.points)
	}

	_, ok := obj.data[record.dateTime]
	if !ok {
		obj.data[record.dateTime] = make(map[string][]dataPoint)
	}
	obj.data[record.dateTime][record.pivot] = record.points
}

func (obj *dataReader) GetColumns() []ColumnStatistic {
	return obj.columns.getColumnStatistics()
}

func (obj *dataReader) GetDataRows() []string {

	rows := make([]string, 0, 10)

	columns := obj.columns.getPivotColumnNames()

	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)

	for i := 0; i < obj.points; i++ {
		writer.WriteString(", null")
	}
	writer.Flush()
	blankPoints := buffer.String()
	buffer.Reset()

	for dateTime, data := range obj.data {

		writer.WriteString("[")
		writer.WriteString(fmt.Sprintf("new Date(%s)", dateTime.Format("2006, 01, 02, 15, 04, 05")))

		for _, columnName := range columns {
			points, ok := data[columnName]

			if ok {
				for i := range points {
					writer.WriteString(points[i].string())
				}
			} else {
				writer.WriteString(blankPoints)
			}
		}

		writer.WriteString("]")
		writer.Flush()

		rows = append(rows, buffer.String())
		buffer.Reset()
	}

	return rows
}

///////////////////////////////////////////////////////

///////////////////////////////////////////////////////
// dataColumns

func (obj *dataColumns) initialize() {
	obj.statistic = make(map[string][]columnStatistic)
}

func (obj *dataColumns) addDataRecord(data dataRecord) {
	if _, ok := obj.statistic[data.pivot]; !ok {
		obj.statistic[data.pivot] = make([]columnStatistic, len(data.points))
	}

	element := obj.statistic[data.pivot]
	for i, dataPoint := range data.points {
		element[i].addDataPoint(dataPoint.float64())
	}
}

func (obj *dataColumns) getPivotColumnNames() []string {
	columns := make([]string, 0, 10)
	if len(obj.statistic) == 0 {
		columns = append(columns, "")
	} else {
		for name := range obj.statistic {
			columns = append(columns, name)
		}
	}

	sort.Strings(columns)

	return columns
}

func (obj *dataColumns) getColumnStatistics() []ColumnStatistic {
	columns := make([]ColumnStatistic, 0)

	getColumnStatistic := func(name string, data columnStatistic) ColumnStatistic {
		return ColumnStatistic{
			Name:    name,
			Minimum: data.minimum,
			Maximum: data.maximum,
			Average: data.sum / float64(data.count),
		}
	}

	columnNames := obj.getPivotColumnNames()
	for i := range columnNames {
		name := columnNames[i]
		pivotData := obj.statistic[name]

		if len(pivotData) == 1 {
			columns = append(columns, getColumnStatistic(name, pivotData[0]))
		} else {
			for i, data := range pivotData {
				var columnName string

				switch {
				case i < len(obj.names) && name == "":
					columnName = obj.names[i]
				case i < len(obj.names) && name != "":
					columnName = fmt.Sprintf("%s %s", name, obj.names[i])
				case name == "":
					columnName = fmt.Sprintf("Column %d", i+1)
				default:
					columnName = fmt.Sprintf("%s %d", name, i+1)
				}
				columns = append(columns, getColumnStatistic(columnName, data))
			}
		}
	}

	return columns
}

///////////////////////////////////////////////////////
// columnStatistic

func (obj *columnStatistic) addDataPoint(data float64) {
	if obj.count == 0 {
		obj.minimum = data
		obj.maximum = data
	} else {
		obj.minimum = min(obj.minimum, data)
		obj.maximum = max(obj.maximum, data)
	}
	obj.sum += data
	obj.count++
}

///////////////////////////////////////////////////////
// dateRecord

func (obj *dataReader) getDataRecord(data string) (record dataRecord) {
	scan := bufio.NewScanner(strings.NewReader(data))

	onDelimiter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.Index(data, obj.delimiter); i >= 0 {
			return i + len(obj.delimiter), data[:i], nil
		}
		return 0, data, bufio.ErrFinalToken
	}
	scan.Split(onDelimiter) //bufio.ScanWords

	record.points = make([]dataPoint, 0, 5)

	for column := 1; scan.Scan(); column++ {
		word := scan.Text()

		if column == obj.dateColumn {
			t, err := time.ParseInLocation(obj.dateFormat, word, time.Local)
			checkErr(err)
			record.dateTime = t

			continue
		}

		if column == obj.pivotColumn {
			record.pivot = word

			continue
		}

		record.points = append(record.points, newDataPoint(word))
	}

	return
}

///////////////////////////////////////////////////////
// datePoint

func newDataPoint(data string) (point dataPoint) {
	point.setValue(data)
	return
}

func (obj *dataPoint) setValue(data string) {

	if len(strings.TrimSpace(data)) == 0 {
		obj.isNull = true
		return
	}

	s, err := strconv.ParseFloat(data, 32)
	checkErr(err)
	obj.point = s
}

func (obj *dataPoint) string() string {
	if obj.isNull {
		return ", null"
	}
	return fmt.Sprintf(", %g", obj.point)
}

func (obj *dataPoint) float64() float64 {
	if obj.isNull {
		return 0
	}
	return obj.point
}
