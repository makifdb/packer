package packer

import (
	"testing"
)

// TODO: Add more tests.

func TestCheck(t *testing.T) {
	type args struct {
		packageName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Wrong package name",
			args: args{"unam"},
			want: false,
		},
		{
			name: "Correct package name",
			args: args{"uname"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(tt.args.packageName); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
