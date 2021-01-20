package models_test

import (
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"testing"
)

func TestDate_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		t       models.Date
		args    args
		wantErr bool
		expect  models.Date
	}{
		{name: "valid string", t: (models.Date)(""), args: args{v: "01-01-2020"}, wantErr: false, expect: (models.Date)("01-01-2020")},
		{name: "incorrect format", t: (models.Date)(""), args: args{v: "01-1-15-2020"}, wantErr: true},
		{name: "out of range time", t: (models.Date)(""), args: args{v: "01-13-2020"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.t != tt.expect {
				t.Errorf("date = %s, want %s", tt.t, tt.expect)
			}
		})
	}
}
