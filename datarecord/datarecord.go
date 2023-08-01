package datarecord

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type dataReader struct {
	dateFormat string
}
type dataRecord struct {
	dateTime time.Time
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
	return
}

func (obj *dataReader) WithDateFormat(dateFormat string) *dataReader {
	obj.dateFormat = dateFormat
	return obj
}

///////////////////////////////////////////////////////
// dateRecord

func (obj *dataReader) GetDataRecord(data string) (record dataRecord) {
	scan := bufio.NewScanner(strings.NewReader(data))
	scan.Split(bufio.ScanWords)
	for scan.Scan() {
		word := scan.Text()

		if record.dateTime.IsZero() {
			t, err := time.ParseInLocation(obj.dateFormat, word, time.Local)
			checkErr(err)
			record.dateTime = t

			continue
		}

		s, err := strconv.ParseFloat(word, 32)
		checkErr(err)
		record.points = append(record.points, float32(s))
	}

	return
}

func (obj *dataRecord) String() string {
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)

	//[new Date(2314, 2, 16), 24045, 12374],

	writer.WriteString("[")
	writer.WriteString(fmt.Sprintf("new Date(%s)", obj.dateTime.Format("2006, 01, 02, 15, 04, 05")))
	for _, point := range obj.points {
		writer.WriteString(fmt.Sprintf(", %g", point))
	}
	writer.WriteString("]")

	writer.Flush()
	return buffer.String()
}

func (obj *dataRecord) Columns() int {
	return len(obj.points)
}
