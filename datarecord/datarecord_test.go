package datarecord

import (
	"reflect"
	"testing"
	"time"
)

func TestGetDataRecord(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want dataRecord
	}{
		{"test 1", args{"20121015100100 1 2 3"}, dataRecord{time.Date(2012, time.October, 15, 10, 1, 0, 0, time.Local), []float32{1, 2, 3}}},
		{"test 2", args{"20121015100130 2 3 4"}, dataRecord{time.Date(2012, time.October, 15, 10, 1, 30, 0, time.Local), []float32{2, 3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDataRecord(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
