package models_test

import (
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"testing"
)

func TestTimeOfDay_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		t       models.TimeOfDay
		args    args
		wantErr bool
		expect  models.TimeOfDay
	}{
		{name: "valid string", t: (models.TimeOfDay)(""), args: args{v: "16:00"}, wantErr: false, expect: (models.TimeOfDay)("16:00")},
		{name: "incorrect format", t: (models.TimeOfDay)(""), args: args{v: "16:00:00"}, wantErr: true},
		{name: "out of range time", t: (models.TimeOfDay)(""), args: args{v: "16:66"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.t != tt.expect {
				t.Errorf("time of day = %s, want %s", tt.t, tt.expect)
			}
		})
	}
}
