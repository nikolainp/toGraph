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

// DataSeries struct
type DataSeries struct {
	Name string

	Minimum float64
	Maximum float64
	Count   int
	Average float64
	Data    string
}

// DataReader ...
type DataReader interface {
	WithDateFormat(string) DataReader
	WithDateColumn(int) DataReader
	WithPivotColumn(int) DataReader
	WithDelimiter(string) DataReader
	WithColumnNames(string) DataReader
	WithColumnNamesInFirstRow(bool) DataReader

	ReadDataRecord(string)
	GetData() []DataSeries
}

///////////////////////////////////////////////////////

type dataReader struct {
	dateFormat              string
	dateColumn              int
	pivotColumn             int
	delimiter               []byte
	isColumnNamesInFirstRow bool
	columnNames             []string

	data map[string]*dataReaderData
}

type dataReaderData struct {
	statistic columnStatistic
	points    []dataPoint
}

type columnStatistic struct {
	minimum float64
	maximum float64
	sum     float64
	count   int
}

type dataPoint struct {
	timestamp time.Time
	point     float64
	isNull    bool
}

type dataRecord struct {
	timestamp time.Time
	pivot     string
	data      []string
}

///////////////////////////////////////////////////////

var checkErr = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

///////////////////////////////////////////////////////
// dataReader

// GetDataReader ...
func GetDataReader() DataReader {

	obj := new(dataReader)
	obj.data = make(map[string]*dataReaderData)

	return obj
}

func (obj *dataReader) WithDateFormat(dateFormat string) DataReader {
	obj.dateFormat = dateFormat
	return obj
}

func (obj *dataReader) WithDateColumn(column int) DataReader {
	obj.dateColumn = column
	return obj
}

func (obj *dataReader) WithPivotColumn(column int) DataReader {
	obj.pivotColumn = column
	return obj
}

func (obj *dataReader) WithDelimiter(delimiter string) DataReader {
	obj.delimiter = []byte(delimiter)
	return obj
}

func (obj *dataReader) WithColumnNames(data string) DataReader {
	if len(data) == 0 {
		return obj
	}

	obj.columnNames = strings.Split(data, ",")

	return obj
}

func (obj *dataReader) WithColumnNamesInFirstRow(data bool) DataReader {
	obj.isColumnNamesInFirstRow = data
	return obj
}

func (obj *dataReader) ReadDataRecord(data string) {

	if obj.isColumnNamesInFirstRow {
		obj.isColumnNamesInFirstRow = false
		obj.columnNames = strings.Split(data, string(obj.delimiter))
		return
	}

	record := obj.getDataRecord(data)
	for id := range record.data {
		isSingleColumn := len(record.data) == 1
		column := obj.getColumnName(record.pivot, id, isSingleColumn)
		point := newDataPoint(record.timestamp, record.data[id])

		if data, ok := obj.data[column]; ok {
			data.statistic.addDataPoint(point)
			data.points = append(data.points, point)
		} else {
			data = new(dataReaderData)
			data.statistic.addDataPoint(point)
			data.points = make([]dataPoint, 0)
			data.points = append(data.points, point)
			obj.data[column] = data
		}
	}

}

func (obj *dataReader) GetData() []DataSeries {

	series := make([]DataSeries, len(obj.data))
	id := 0

	for dataName, data := range obj.data {
		series[id].Name = dataName

		series[id].Minimum = data.statistic.minimum
		series[id].Maximum = data.statistic.maximum
		series[id].Count = data.statistic.count
		series[id].Average = data.statistic.average()

		sort.Slice(data.points, func(i, j int) bool {
			return data.points[j].timestamp.After(data.points[i].timestamp)
		})

		rows := make([]string, len(data.points))
		for i, point := range data.points {
			rows[i] = point.string()
		}
		series[id].Data = strings.Join(rows, ",")

		id++
	}

	return series
}

///////////////////////////////////////////////////////

func (obj *dataReader) getColumnName(pivotName string, id int, isSingleColumn bool) string {
	var columnName string

	switch {
	case id < len(obj.columnNames) && pivotName == "":
		columnName = obj.columnNames[id]
	case id < len(obj.columnNames) && pivotName != "":
		columnName = fmt.Sprintf("%s %s", pivotName, obj.columnNames[id])
	case pivotName == "":
		columnName = fmt.Sprintf("Column %d", id+1)
	case isSingleColumn:
		columnName = fmt.Sprintf("%s", pivotName)
	default:
		columnName = fmt.Sprintf("%s %d", pivotName, id+1)
	}

	return columnName
}

///////////////////////////////////////////////////////
// columnStatistic

func (obj *columnStatistic) addDataPoint(data dataPoint) {
	if obj.count == 0 {
		obj.minimum = data.float64()
		obj.maximum = data.float64()
	} else {
		obj.minimum = min(obj.minimum, data.float64())
		obj.maximum = max(obj.maximum, data.float64())
	}
	obj.sum += data.float64()
	obj.count += data.count()
}

func (obj *columnStatistic) average() float64 {
	if obj.count == 0 {
		return 0
	}
	return obj.sum / float64(obj.count)
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

	record.data = make([]string, 0, 5)

	for column := 1; scan.Scan(); column++ {
		word := scan.Text()

		if column == obj.dateColumn {
			t, err := time.ParseInLocation(obj.dateFormat, word, time.Local)
			checkErr(err)
			record.timestamp = t

			continue
		}

		if column == obj.pivotColumn {
			record.pivot = word

			continue
		}

		record.data = append(record.data, word)
	}

	return
}

///////////////////////////////////////////////////////
// datePoint

func newDataPoint(dateTime time.Time, data string) (point dataPoint) {
	point.timestamp = dateTime
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
		return fmt.Sprintf("['%s', null]",
			obj.timestamp.Format("2006-01-02 15:04:05"))
	}
	return fmt.Sprintf("['%s', %g]",
		obj.timestamp.Format("2006-01-02 15:04:05"), obj.point)
}

func (obj *dataPoint) float64() float64 {
	if obj.isNull {
		return 0
	}
	return obj.point
}

func (obj *dataPoint) count() int {
	if obj.isNull {
		return 0
	}
	return 1
}
