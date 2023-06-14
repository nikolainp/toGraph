package datarecord

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type dataRecord struct {
	dateTime time.Time
	points   []float32
}

func GetDataRecord(data string) (record dataRecord) {
	scan := bufio.NewScanner(strings.NewReader(data))
	scan.Split(bufio.ScanWords)
	for scan.Scan() {
		word := scan.Text()

		if record.dateTime.IsZero() {
			t, err := time.ParseInLocation("20060102150405", word, time.Local)
			if err != nil {
				fmt.Print(err)
			}
			record.dateTime = t

			continue
		}

		s, err := strconv.ParseFloat(word, 32)
		if err != nil {
			fmt.Print(err)
		}
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
