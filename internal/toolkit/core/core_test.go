package core

import (
	"testing"

	"github.com/matryer/is"
)

func TestIsDebug(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{
			name: "returns_false_if_env_is_unset",
			env:  "",
			want: false,
		},
		{
			name: "returns_true_if_env_is_set",
			env:  "1",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(runnerDebug, tt.env)

			result := IsDebug()

			is.Equal(result, tt.want) // should match expected
		})
	}
}
