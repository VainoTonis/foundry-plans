package cmd

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseRepositoryIDs(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    []int64
		wantErr bool
	}{
		{name: "valid", value: []interface{}{json.Number("1"), json.Number("42")}, want: []int64{1, 42}},
		{name: "missing", value: nil, wantErr: true},
		{name: "empty", value: []interface{}{}, wantErr: true},
		{name: "string", value: []interface{}{"1"}, wantErr: true},
		{name: "fraction", value: []interface{}{json.Number("1.5")}, wantErr: true},
		{name: "zero", value: []interface{}{json.Number("0")}, wantErr: true},
		{name: "negative", value: []interface{}{json.Number("-1")}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRepositoryIDs(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseRepositoryIDs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseRepositoryIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
