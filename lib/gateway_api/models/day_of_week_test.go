package models_test

import (
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"testing"
)

func TestDayOfWeek_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name     string
		i        models.DayOfWeek
		args     args
		wantErr  bool
		expected models.DayOfWeek
	}{
		{name: "monday", i: (models.DayOfWeek)(0), args: args{v: 1}, wantErr: false, expected: models.Monday},
		{name: "out of range low", i: (models.DayOfWeek)(0), args: args{v: 0}, wantErr: true},
		{name: "out of range high", i: (models.DayOfWeek)(0), args: args{v: 8}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.i != tt.expected {
				t.Errorf("got day of week = '%v', wanted '%v'", tt.i, tt.expected)
			}
		})
	}
}
