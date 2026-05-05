package tool

import (
	"testing"

	"github.com/matryer/is"
)

func TestCmd_VersionCheck(t *testing.T) {
	c := &Cmd{}

	tests := []struct {
		name    string
		version string
		spec    string
		want    bool
		wantErr bool
	}{
		{
			name:    "errors_if_version_is_invalid",
			version: "",
			spec:    "*",
			wantErr: true,
		},
		{
			name:    "errors_if_version_spec_is_invalid",
			version: "1.0.0",
			spec:    "abc",
			wantErr: true,
		},
		{
			name:    "returns_true_for_matching_version",
			version: "1.0.0",
			spec:    "^1.0.0",
			want:    true,
		},
		{
			name:    "returns_false_for_non-matching_version",
			version: "2.0.0",
			spec:    "^1.0.0",
			want:    false,
		},
		{
			name:    "matches_exact_version",
			version: "1.2.3",
			spec:    "1.2.3",
			want:    true,
		},
		{
			name:    "matches_wildcard",
			version: "5.6.7",
			spec:    "*",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			valid, err := c.VersionCheck(tt.version, tt.spec)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)            // should not error
			is.Equal(valid, tt.want) // should match expected
		})
	}
}
