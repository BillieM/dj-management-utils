package helpers_test

import (
	"testing"

	"github.com/billiem/seren-management/pkg/helpers"
)

func TestContainsNonEmptyString(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "empty slice",
			args: []string{},
			want: false,
		},
		{
			name: "all empty strings",
			args: []string{"", "", ""},
			want: false,
		},
		{
			name: "one non-empty string",
			args: []string{"", "hello", ""},
			want: true,
		},
		{
			name: "multiple non-empty strings",
			args: []string{"", "hello", "world"},
			want: true,
		},
		{
			name: "only non-empty string",
			args: []string{"hello"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.ContainsNonEmptyString(tt.args); got != tt.want {
				t.Errorf("ContainsNonEmptyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64ArrayToJoinedString(t *testing.T) {
	tests := []struct {
		name string
		args []int64
		want string
	}{
		{
			name: "empty slice",
			args: []int64{},
			want: "",
		},
		{
			name: "single element",
			args: []int64{1},
			want: "1",
		},
		{
			name: "multiple elements",
			args: []int64{1, 2, 3},
			want: "1,2,3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.Int64ArrayToJoinedString(tt.args); got != tt.want {
				t.Errorf("Int64ArrayToJoinedString() = %v, want %v", got, tt.want)
			}
		})
	}
}
