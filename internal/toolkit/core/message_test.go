package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestDebug(t *testing.T) {
	t.Run("writes the debug message command", func(t *testing.T) {
		is := is.New(t)
		m := "message"

		var b bytes.Buffer
		err := Debug(&b, m)

		is.NoErr(err)                                                // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", DebugCmd, m)) // should be equal
	})
}

func TestError(t *testing.T) {
	t.Run("writes the error log message command", func(t *testing.T) {
		is := is.New(t)
		m := "message"

		var b bytes.Buffer
		err := Error(&b, m, AnnotationProperties{})

		is.NoErr(err)                                                // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", ErrorCmd, m)) // should be equal
	})

	t.Run("writes the error log message command with an annotation", func(t *testing.T) {
		is := is.New(t)
		m := "message"
		a := AnnotationProperties{
			Title:     "title",
			File:      "file",
			Column:    1,
			EndColumn: 2,
			Line:      3,
			EndLine:   4,
		}

		var b bytes.Buffer
		err := Error(&b, m, a)

		is.NoErr(err)                                                                                                                                                                   // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s title=%s,file=%s,col=%d,endColumn=%d,line=%d,endLine=%d::%s\n", ErrorCmd, a.Title, a.File, a.Column, a.EndColumn, a.Line, a.EndLine, m)) // should be equal
	})
}

func TestWarning(t *testing.T) {
	t.Run("writes the warning log message command", func(t *testing.T) {
		is := is.New(t)
		m := "message"

		var b bytes.Buffer
		err := Warning(&b, m, AnnotationProperties{})

		is.NoErr(err)                                                  // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", WarningCmd, m)) // should be equal
	})

	t.Run("writes the warning log message command with an annotation", func(t *testing.T) {
		is := is.New(t)
		m := "message"
		a := AnnotationProperties{
			Title:     "title",
			File:      "file",
			Column:    1,
			EndColumn: 2,
			Line:      3,
			EndLine:   4,
		}

		var b bytes.Buffer
		err := Warning(&b, m, a)

		is.NoErr(err)                                                                                                                                                                     // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s title=%s,file=%s,col=%d,endColumn=%d,line=%d,endLine=%d::%s\n", WarningCmd, a.Title, a.File, a.Column, a.EndColumn, a.Line, a.EndLine, m)) // should be equal
	})
}

func TestNotice(t *testing.T) {
	t.Run("writes the notice log message command", func(t *testing.T) {
		is := is.New(t)
		m := "message"

		var b bytes.Buffer
		err := Notice(&b, m, AnnotationProperties{})

		is.NoErr(err)                                                 // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", NoticeCmd, m)) // should be equal
	})

	t.Run("writes the notice log message command with an annotation", func(t *testing.T) {
		is := is.New(t)
		m := "message"
		a := AnnotationProperties{
			Title:     "title",
			File:      "file",
			Column:    1,
			EndColumn: 2,
			Line:      3,
			EndLine:   4,
		}

		var b bytes.Buffer
		err := Notice(&b, m, a)

		is.NoErr(err)                                                                                                                                                                    // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s title=%s,file=%s,col=%d,endColumn=%d,line=%d,endLine=%d::%s\n", NoticeCmd, a.Title, a.File, a.Column, a.EndColumn, a.Line, a.EndLine, m)) // should be equal
	})
}

func TestInfo(t *testing.T) {
	t.Run("writes the info log message command", func(t *testing.T) {
		is := is.New(t)
		m := "message"

		var b bytes.Buffer
		err := Info(&b, m, AnnotationProperties{})

		is.NoErr(err)                                               // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", InfoCmd, m)) // should be equal
	})

	t.Run("writes the info log message command with an annotation", func(t *testing.T) {
		is := is.New(t)
		m := "message"
		a := AnnotationProperties{
			Title:     "title",
			File:      "file",
			Column:    1,
			EndColumn: 2,
			Line:      3,
			EndLine:   4,
		}

		var b bytes.Buffer
		err := Info(&b, m, a)

		is.NoErr(err)                                                                                                                                                                  // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s title=%s,file=%s,col=%d,endColumn=%d,line=%d,endLine=%d::%s\n", InfoCmd, a.Title, a.File, a.Column, a.EndColumn, a.Line, a.EndLine, m)) // should be equal
	})
}

func TestStartGroup(t *testing.T) {
	t.Run("writes the start log group command", func(t *testing.T) {
		is := is.New(t)
		n := "test"

		var b bytes.Buffer
		err := StartGroup(&b, n)

		is.NoErr(err)                                                     // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", StartGroupCmd, n)) // should be equal
	})
}

func TestEndGroup(t *testing.T) {
	t.Run("writes the end log group command", func(t *testing.T) {
		is := is.New(t)

		var b bytes.Buffer
		err := EndGroup(&b)

		is.NoErr(err)                                              // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::\n", EndGroupCmd)) // should be equal
	})
}
