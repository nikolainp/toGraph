package datarecord

import (
	"reflect"
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
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "", []dataPoint{{point: 1}, {point: 2}, {point: 3}}}},
		{"test 2", args{"20121015100130 2 3 4"},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), "", []dataPoint{{point: 2}, {point: 3}, {point: 4}}}},
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
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "first", []dataPoint{{point: 1}, {point: 2}, {point: 3}}}},
		{"test 2", args{"20121015100130 second 2 3 4"},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), "second", []dataPoint{{point: 2}, {point: 3}, {point: 4}}}},
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
		want dataReader
		data string
	}{
		{"test 1", dataReader{
			dateColumn:  1,
			dateFormat:  "20060102150405",
			pivotColumn: 0,
			delimiter:   []byte{' '},

			points: 3,
			data:   dataReaderData{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local): {"": []dataPoint{{point: 1}, {point: 2}, {point: 3}}}},
		},
			"20121015100130 1 2 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := GetDataReader()

			got.dateColumn = tt.want.dateColumn
			got.dateFormat = tt.want.dateFormat
			got.pivotColumn = tt.want.pivotColumn
			got.delimiter = tt.want.delimiter

			got.ReadDataRecord(tt.data)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("GetDataRecord():\n got  %v\n want %v", got, tt.want)
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func Test_dataColumns_addDataRecord(t *testing.T) {
	tests := []struct {
		name string
		obj  dataColumns
		data dataRecord
		want dataColumns
	}{
		{
			"test 1",
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{}},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "first", []dataPoint{{point: 1}, {point: 2}, {point: 3}}},
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"first": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
		},
		{
			"test 2",
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{}},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "", []dataPoint{{point: 1}, {point: 2}, {point: 3}}},
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
		},
		{
			"test 3",
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"first": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
			dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), "first", []dataPoint{{point: 10}, {point: -2}, {point: 3}}},
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"first": {{1, 10, 11, 2}, {-2, 2, 0, 2}, {3, 3, 6, 2}}}},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			//tt.obj.initialize()
			tt.obj.addDataRecord(tt.data)
			if !reflect.DeepEqual(tt.obj, tt.want) {
				t.Errorf("GetDataRecord():\n got  %v\n want %v", tt.obj, tt.want)
			}
		})
	}
}
func Test_dataColumns_GetColumns(t *testing.T) {
	tests := []struct {
		name string
		obj  dataColumns
		want []string
	}{
		{
			"test 1", dataColumns{statistic: map[string][]columnStatistic{"first": {{}, {}, {}}, "second": {{}, {}, {}}}},
			[]string{"first", "second"},
		},
		{
			"test 2", dataColumns{statistic: map[string][]columnStatistic{"": {{}, {}, {}}}},
			[]string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.getPivotColumnNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataReader.GetColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataColumns_getColumnStatistics(t *testing.T) {
	tests := []struct {
		name string
		obj  dataColumns
		want []ColumnStatistic
	}{
		{
			"test 1",
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"first": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
			[]ColumnStatistic{{"first 1", 1, 1, 1}, {"first 2", 2, 2, 2}, {"first 3", 3, 3, 3}},
		},
		{
			"test 2",
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
			[]ColumnStatistic{{"Column 1", 1, 1, 1}, {"Column 2", 2, 2, 2}, {"Column 3", 3, 3, 3}},
		},
		{
			"test 3",
			dataColumns{names: []string{}, statistic: map[string][]columnStatistic{"first": {{1, 1, 1, 1}}}},
			[]ColumnStatistic{{"first", 1, 1, 1}},
		},
		{
			"test 4",
			dataColumns{names: []string{"A","B","C"}, statistic: map[string][]columnStatistic{"": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
			[]ColumnStatistic{{"A", 1, 1, 1}, {"B", 2, 2, 2}, {"C", 3, 3, 3}},
		},
		{
			"test 5",
			dataColumns{names: []string{"A","B","C"}, statistic: map[string][]columnStatistic{"first": {{1, 1, 1, 1}, {2, 2, 2, 1}, {3, 3, 3, 1}}}},
			[]ColumnStatistic{{"first A", 1, 1, 1}, {"first B", 2, 2, 2}, {"first C", 3, 3, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.getColumnStatistics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataColumns.getColumnStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}
