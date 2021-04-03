package main

import "testing"

func TestGetSequenceNumber(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{
			name: "happy",
			want: 1,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSequenceNumber(); got != tt.want {
				t.Errorf("GetSequenceNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getContentDisposition(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContentDisposition(tt.args.format); got != tt.want {
				t.Errorf("getContentDisposition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_requestStateList(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requestStateList()
			if (err != nil) != tt.wantErr {
				t.Errorf("requestStateList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("requestStateList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
