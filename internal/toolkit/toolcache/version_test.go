package toolcache

import (
	"testing"

	"github.com/matryer/is"
)

func TestCheckVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		spec    string
		want    bool
		wantErr bool
	}{
		{
			name:    "errors_if_version_is_empty",
			version: "",
			spec:    "*",
			wantErr: true,
		},
		{
			name:    "errors_if_version_is_invalid",
			version: "a.b.c",
			spec:    "*",
			wantErr: true,
		},
		{
			name:    "errors_if_version_has_v_prefix",
			version: "v1.0.0",
			spec:    "1.0.0",
			wantErr: true,
		},
		{
			name:    "errors_if_version_spec_is_empty",
			version: "1.0.0",
			spec:    "",
			wantErr: true,
		},
		{
			name:    "errors_if_version_spec_is_invalid",
			version: "1.0.0",
			spec:    "test",
			wantErr: true,
		},
		{
			name:    "matches_explicit_version",
			version: "1.0.0",
			spec:    "1.0.0",
			want:    true,
		},
		{
			name:    "blocks_invalid_explicit_version",
			version: "1.0.0",
			spec:    "1.0.1",
			want:    false,
		},
		{
			name:    "matches_patch_version_constraint",
			version: "1.0.1",
			spec:    "~1.0.0",
			want:    true,
		},
		{
			name:    "blocks_invalid_patch_version_constraint",
			version: "1.0.1",
			spec:    "~1.1.0",
			want:    false,
		},
		{
			name:    "matches_minor_version_constraint",
			version: "1.1.1",
			spec:    "^1.0.0",
			want:    true,
		},
		{
			name:    "blocks_invalid_minor_version_constraint",
			version: "1.1.1",
			spec:    "^2.0.0",
			want:    false,
		},
		{
			name:    "matches_greater_than_version_constraint",
			version: "1.1.1",
			spec:    ">=1.0.0",
			want:    true,
		},
		{
			name:    "blocks_invalid_greater_than_version_constraint",
			version: "1.1.1",
			spec:    ">=1.2.0",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			valid, err := CheckVersion(tt.version, tt.spec)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)            // should not error
			is.Equal(valid, tt.want) // should match expected
		})
	}
}
