package msgraph

import (
	"fmt"
	"testing"
	"time"
)

func TestResponseStatus_Equal(t *testing.T) {

	testTime := time.Now()
	type args struct {
		other ResponseStatus
	}
	tests := []struct {
		name string
		s    ResponseStatus
		args args
		want bool
	}{
		{
			name: "Equal",
			s:    ResponseStatus{Response: "accepted", Time: testTime},
			args: args{other: ResponseStatus{Response: "confirmed", Time: testTime}},
			want: true,
		}, {
			name: "non-equal response",
			s:    ResponseStatus{Response: "declined", Time: testTime},
			args: args{other: ResponseStatus{Response: "accepted", Time: testTime}},
			want: false,
		}, {
			name: "non-equal time",
			s:    ResponseStatus{Response: "accepted", Time: testTime.Add(time.Minute)},
			args: args{other: ResponseStatus{Response: "accepted", Time: testTime}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Equal(tt.args.other); got != tt.want {
				t.Errorf("ResponseStatus.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseStatus_UnmarshalJSON(t *testing.T) {
	fmt.Println(time.Now().Format(time.RFC3339Nano))
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       *ResponseStatus
		args    args
		wantErr bool
	}{
		{
			name:    "all good",
			s:       &ResponseStatus{},
			args:    args{data: []byte("{\"response\" : \"accepted\", \"time\" : \"" + time.Now().Format(time.RFC3339Nano) + "\"}")},
			wantErr: false,
		}, {
			name:    "invalid time format",
			s:       &ResponseStatus{},
			args:    args{data: []byte("{\"response\" : \"accepted\", \"time\" : \"" + time.Now().Format(time.RFC1123) + "\"}")},
			wantErr: true,
		}, {
			name:    "response-field empty",
			s:       &ResponseStatus{},
			args:    args{data: []byte("{\"time\" : \"" + time.Now().Format(time.RFC1123) + "\"}")},
			wantErr: true,
		}, {
			name:    "time-field empty",
			s:       &ResponseStatus{},
			args:    args{data: []byte("{\"response\" : \"accepted\"}")},
			wantErr: true,
		}, {
			name:    "invalid json",
			s:       &ResponseStatus{},
			args:    args{data: []byte("{\"response\"  \"accepted\", \"time\" : \"" + time.Now().Format(time.RFC3339Nano) + "\"}")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ResponseStatus.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
