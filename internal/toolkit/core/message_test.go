package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestDebug(t *testing.T) {
	is := is.New(t)

	var b bytes.Buffer
	err := Debug(&b, "message")

	is.NoErr(err)                                                  // should not error
	is.Equal(b.String(), fmt.Sprintf("::%s::message\n", DebugCmd)) // should be equal
}

func TestError(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		ann  AnnotationProperties
		want string
	}{
		{
			name: "writes_the_error_log_message_command",
			msg:  "message",
			ann:  AnnotationProperties{},
			want: fmt.Sprintf("::%s::message\n", ErrorCmd),
		},
		{
			name: "writes_the_error_log_message_command_with_an_annotation",
			msg:  "message",
			ann:  AnnotationProperties{Title: "title", File: "file", Column: 1, EndColumn: 2, Line: 3, EndLine: 4},
			want: fmt.Sprintf("::%s title=title,file=file,col=1,endColumn=2,line=3,endLine=4::message\n", ErrorCmd),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := Error(&b, tt.msg, tt.ann)

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should be equal
		})
	}
}

func TestWarning(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		ann  AnnotationProperties
		want string
	}{
		{
			name: "writes_the_warning_log_message_command",
			msg:  "message",
			ann:  AnnotationProperties{},
			want: fmt.Sprintf("::%s::message\n", WarningCmd),
		},
		{
			name: "writes_the_warning_log_message_command_with_an_annotation",
			msg:  "message",
			ann:  AnnotationProperties{Title: "title", File: "file", Column: 1, EndColumn: 2, Line: 3, EndLine: 4},
			want: fmt.Sprintf("::%s title=title,file=file,col=1,endColumn=2,line=3,endLine=4::message\n", WarningCmd),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := Warning(&b, tt.msg, tt.ann)

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should be equal
		})
	}
}

func TestNotice(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		ann  AnnotationProperties
		want string
	}{
		{
			name: "writes_the_notice_log_message_command",
			msg:  "message",
			ann:  AnnotationProperties{},
			want: fmt.Sprintf("::%s::message\n", NoticeCmd),
		},
		{
			name: "writes_the_notice_log_message_command_with_an_annotation",
			msg:  "message",
			ann:  AnnotationProperties{Title: "title", File: "file", Column: 1, EndColumn: 2, Line: 3, EndLine: 4},
			want: fmt.Sprintf("::%s title=title,file=file,col=1,endColumn=2,line=3,endLine=4::message\n", NoticeCmd),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := Notice(&b, tt.msg, tt.ann)

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should be equal
		})
	}
}

func TestInfo(t *testing.T) {
	is := is.New(t)

	var b bytes.Buffer
	err := Info(&b, "message")

	is.NoErr(err)                     // should not error
	is.Equal(b.String(), "message\n") // should be direct output
}

func TestStartGroup(t *testing.T) {
	is := is.New(t)

	var b bytes.Buffer
	err := StartGroup(&b, "test")

	is.NoErr(err)                                                    // should not error
	is.Equal(b.String(), fmt.Sprintf("::%s::test\n", StartGroupCmd)) // should be equal
}

func TestEndGroup(t *testing.T) {
	is := is.New(t)

	var b bytes.Buffer
	err := EndGroup(&b)

	is.NoErr(err)                                              // should not error
	is.Equal(b.String(), fmt.Sprintf("::%s::\n", EndGroupCmd)) // should be equal
}
