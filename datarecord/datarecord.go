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

func (obj dataRecord) String() string {
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)

	writer.WriteString(fmt.Sprintf("{ \"when\": \"%s\",\n", obj.dateTime.Format("01/02/2006 15:04:05")))
	writer.WriteString("\t\"type1\": [ ")
	for i, point := range obj.points {
		if i == 0 {
			writer.WriteString(fmt.Sprintf("%g", point))
		} else {
			writer.WriteString(fmt.Sprintf(", %g", point))
		}
	}
	writer.WriteString(" ]\n")
	writer.WriteString("}")

	writer.Flush()
	return buffer.String()
}
