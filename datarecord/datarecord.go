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
	Minimum float32
	Maximum float32
	Average float32
}

type dataReaderData map[time.Time]map[string][]float32

type columnStatistic struct {
	minimum float32
	maximum float32
	sum     float32
	count   int
}

type dataColumns struct {
	names     []string
	statistic map[string][]columnStatistic
}

type dataReader struct {
	dateFormat  string
	dateColumn  int
	pivotColumn int
	delimiter   []byte

	columns dataColumns
	points  int
	data    dataReaderData
}

type dataRecord struct {
	dateTime time.Time
	pivot    string
	points   []float32
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

func (obj *dataReader) ReadDataRecord(data string) {
	record := obj.getDataRecord(data)

	obj.columns.addDataRecord(record)
	if obj.points < len(record.points) {
		obj.points = len(record.points)
	}

	_, ok := obj.data[record.dateTime]
	if !ok {
		obj.data[record.dateTime] = make(map[string][]float32)
	}
	obj.data[record.dateTime][record.pivot] = record.points
}

func (obj *dataReader) GetColumns() []ColumnStatistic {
	return obj.columns.getColumnStatistics()
}

func (obj *dataReader) GetDataRows() []string {

	rows := make([]string, 0, 10)

	columns := obj.columns.getColumnNames()

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
				for _, point := range points {
					writer.WriteString(fmt.Sprintf(", %g", point))
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
		element[i].addDataPoint(dataPoint)
	}

}

func (obj *dataColumns) getColumnNames() []string {
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

	for name, pivotData := range obj.statistic {
		if name == "" {
			name = "Column"
		}
		for i, data := range pivotData {
			columns = append(columns, ColumnStatistic{
				Name:    fmt.Sprintf("%s %d", name, i+1),
				Minimum: data.minimum,
				Maximum: data.maximum,
				Average: data.sum / float32(data.count),
			})
		}
	}

	return columns
}

///////////////////////////////////////////////////////
// columnStatistic

func (obj *columnStatistic) addDataPoint(data float32) {
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

	record.points = make([]float32, 0, 5)

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

		s, err := strconv.ParseFloat(word, 32)
		checkErr(err)
		record.points = append(record.points, float32(s))
	}

	return
}
