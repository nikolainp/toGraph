package datarecord

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/go-test/deep"
)

func TestGetDataRecord(t *testing.T) {
	var reader dataReader
	reader.dateFormat = "20060102150405"
	reader.dateColumn = 1
	reader.delimiter = []byte{' '}

	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want dataRecord
	}{
		{"test 1", args{"20121015100100 1 2 3"},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "", []string{"1", "2", "3"}}},
		{"test 2", args{"20121015100130 2 3 4"},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), "", []string{"2", "3", "4"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reader.getDataRecord(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDataRecordPivot(t *testing.T) {
	var reader dataReader
	reader.dateFormat = "20060102150405"
	reader.dateColumn = 1
	reader.pivotColumn = 2
	reader.delimiter = []byte{' '}

	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want dataRecord
	}{
		{"test 1", args{"20121015100100 first 1 2 3"},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "first", []string{"1", "2", "3"}}},
		{"test 2", args{"20121015100130 second 2 3 4"},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), "second", []string{"2", "3", "4"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reader.getDataRecord(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataReader_ReadDataRecord(t *testing.T) {
	tests := []struct {
		name string
		want *dataReader
		data string
	}{
		{"test 1",
			&dataReader{
				dateColumn:  1,
				dateFormat:  "20060102150405",
				pivotColumn: 0,
				delimiter:   []byte{' '},

				data: map[string]*dataReaderData{
					"Column 1": {
						statistic: columnStatistic{1, 1, 1, 1},
						data:      []dataPoint{{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 1, false}},
					},
					"Column 2": {
						statistic: columnStatistic{2, 2, 2, 1},
						data:      []dataPoint{{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 2, false}},
					},
					"Column 3": {
						statistic: columnStatistic{3, 3, 3, 1},
						data:      []dataPoint{{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 3, false}},
					},
				},
			},
			"20121015100130 1 2 3",
		},
	}

	//{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local): {"": []dataPoint{{point: 1}, {point: 2}, {point: 3}}}			},
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := GetDataReader()

			got.WithDateColumn(tt.want.dateColumn)
			got.WithDateFormat(tt.want.dateFormat)
			got.WithPivotColumn(tt.want.pivotColumn)
			got.WithDelimiter(string(tt.want.delimiter))

			got.ReadDataRecord(tt.data)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("GetDataRecord():\n got  %v\n want %v", got, tt.want)
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func Test_dataReader_GetData(t *testing.T) {
	tests := []struct {
		name string
		obj  dataReader
		want []DataSeries
	}{
		{"test 1", dataReader{
			dateColumn:  1,
			dateFormat:  "20060102150405",
			pivotColumn: 0,
			delimiter:   []byte{' '},

			data: map[string]*dataReaderData{
				"Column 1": {
					statistic: columnStatistic{1, 1, 1, 1},
					data:      []dataPoint{{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 1, false}},
				},
			},
		},
			[]DataSeries{
				{"Column 1", 1, 1, 1, 1, "['2012-10-15 10:01:30', 1]"},
			},
		},
		{"test 2", dataReader{
			dateColumn:  1,
			dateFormat:  "20060102150405",
			pivotColumn: 0,
			delimiter:   []byte{' '},

			data: map[string]*dataReaderData{
				"Column 1": {
					statistic: columnStatistic{1, 3, 6, 3},
					data: []dataPoint{
						{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 1, false},
						{time.Date(2012, time.October, 15, 10, 1, 40, 0, time.Local), 2, false},
						{time.Date(2012, time.October, 15, 10, 1, 50, 0, time.Local), 3, false},
					},
				},
			},
		},
			[]DataSeries{
				{"Column 1", 1, 3, 3, 2, "['2012-10-15 10:01:30', 1],['2012-10-15 10:01:40', 2],['2012-10-15 10:01:50', 3]"},
			},
		},
		{"test 3", dataReader{
			dateColumn:  1,
			dateFormat:  "20060102150405",
			pivotColumn: 0,
			delimiter:   []byte{' '},

			data: map[string]*dataReaderData{
				"Column 1": {
					statistic: columnStatistic{1, 3, 6, 3},
					data: []dataPoint{
						{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 1, false},
						{time.Date(2012, time.October, 15, 10, 1, 40, 0, time.Local), 2, false},
						{time.Date(2012, time.October, 15, 10, 1, 50, 0, time.Local), 3, false},
					},
				},
				"Column 2": {
					statistic: columnStatistic{1, 3, 6, 3},
					data: []dataPoint{
						{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), 4, false},
						{time.Date(2012, time.October, 15, 10, 1, 40, 0, time.Local), 5, false},
						{time.Date(2012, time.October, 15, 10, 1, 50, 0, time.Local), 6, false},
					},
				},
			},
		},
			[]DataSeries{
				{"Column 1", 1, 3, 3, 2, "['2012-10-15 10:01:30', 1],['2012-10-15 10:01:40', 2],['2012-10-15 10:01:50', 3]"},
				{"Column 2", 1, 3, 3, 2, "['2012-10-15 10:01:30', 4],['2012-10-15 10:01:40', 5],['2012-10-15 10:01:50', 6]"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.GetData()
			sort.Slice(got, func(i, j int) bool { return got[i].Name < got[j].Name })
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataReader.GetDataRows() = %v, want %v", got, tt.want)
			}
		})
	}
}
