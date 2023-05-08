package datarecord

import (
	"bufio"
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
