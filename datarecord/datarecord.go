package datarecord

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

type void struct{}
type dataReaderData map[time.Time]map[string][]float32
type dataReader struct {
	dateFormat  string
	dateColumn  int
	pivotColumn int

	columnNames map[string]void
	points      int
	data        dataReaderData
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

	reader.columnNames = make(map[string]void)
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

func (obj *dataReader) ReadDataRecord(data string) {
	record := obj.getDataRecord(data)

	if len(record.pivot) != 0 {
		obj.columnNames[record.pivot] = void{}
	}
	if obj.points < len(record.points) {
		obj.points = len(record.points)
	}

	_, ok := obj.data[record.dateTime]
	if !ok {
		obj.data[record.dateTime] = make(map[string][]float32)
	}
	obj.data[record.dateTime][record.pivot] = record.points
}

func (obj *dataReader) GetColumns() []string {

	columns := make([]string, 0, 10)
	if len(obj.columnNames) == 0 {
		for i := 1; i < obj.points+1; i++ {
			columns = append(columns, fmt.Sprintf("Column %d", i))
		}
	} else {
		for name := range obj.columnNames {
			for i := 1; i < obj.points+1; i++ {
				columns = append(columns, fmt.Sprintf("%s%d", name, i))
			}
		}
	}

	return columns
}

func (obj *dataReader) GetDataRows() []string {

	rows := make([]string, 0, 10)

	columns := maps.Keys(obj.columnNames)
	if len(columns) == 0 {
		columns = append(columns, "")
	}

	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)

	for i := 0; i < obj.points; i++ {
		writer.WriteString(", ")
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
// dateRecord

func (obj *dataReader) getDataRecord(data string) (record dataRecord) {
	scan := bufio.NewScanner(strings.NewReader(data))
	scan.Split(bufio.ScanWords)

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

// func (obj *dataRecord) String() string {
// 	buffer := new(bytes.Buffer)
// 	writer := bufio.NewWriter(buffer)

// 	//[new Date(2314, 2, 16), {"": [24045, 12374]}],

// 	writer.WriteString("[")
// 	writer.WriteString(fmt.Sprintf("new Date(%s)", obj.dateTime.Format("2006, 01, 02, 15, 04, 05")))
// 	for _, point := range obj.points {
// 		writer.WriteString(fmt.Sprintf(", %g", point))
// 	}
// 	writer.WriteString("]")

// 	writer.Flush()
// 	return buffer.String()
// }

func (obj *dataRecord) Columns() int {
	return len(obj.points)
}
