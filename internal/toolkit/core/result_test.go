package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestSetFailed(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "writes_the_error_command",
			msg:  "Failure message",
			want: fmt.Sprintf("::%s::Failure message\n", ErrorCmd),
		},
		{
			name: "escapes_the_failure_message",
			msg:  "Failure \r\n\nmessage\r",
			want: fmt.Sprintf("::%s::Failure %%0D%%0A%%0Amessage%%0D\n", ErrorCmd),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := SetFailed(&b, tt.msg)

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should match expected
		})
	}
}
