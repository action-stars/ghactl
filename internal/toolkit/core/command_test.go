package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestNewCommand(t *testing.T) {
	t.Run("errors if there is no command type", func(t *testing.T) {
		is := is.New(t)
		_, err := NewCommand("", nil, "")

		is.True(err != nil) // should error
	})

	t.Run("creates a new command with just the type", func(t *testing.T) {
		is := is.New(t)
		ct := StartGroupCmd
		c, err := NewCommand(ct, nil, "")

		is.NoErr(err)               // should not error
		is.Equal(c.Type, ct)        // should match
		is.Equal(c.Properties, nil) // should match
		is.Equal(c.Message, "")     // should match
	})

	t.Run("creates a new command with all inputs", func(t *testing.T) {
		is := is.New(t)
		ct := StartGroupCmd
		cp := CommandProperties{{Key: "name", Value: "test"}}
		cm := "message"
		c, err := NewCommand(ct, cp, cm)

		is.NoErr(err)              // should not error
		is.Equal(c.Type, ct)       // should match
		is.Equal(c.Properties, cp) // should match
		is.Equal(c.Message, cm)    // should match
	})
}

func TestIssueCommand(t *testing.T) {
	t.Run("writes a command with only a type", func(t *testing.T) {
		is := is.New(t)
		c := Command{Type: StartGroupCmd}

		var b bytes.Buffer
		err := IssueCommand(&b, &c)

		is.NoErr(err)                                         // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::\n", c.Type)) // should match
	})

	t.Run("writes a command with a type and message", func(t *testing.T) {
		is := is.New(t)
		c := Command{Type: StartGroupCmd, Message: "message"}

		var b bytes.Buffer
		err := IssueCommand(&b, &c)

		is.NoErr(err)                                                      // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", c.Type, c.Message)) // should match
	})

	t.Run("writes a command with a type and a property", func(t *testing.T) {
		is := is.New(t)
		c := Command{Type: StartGroupCmd, Properties: CommandProperties{{Key: "name", Value: "test"}}}

		var b bytes.Buffer
		err := IssueCommand(&b, &c)

		is.NoErr(err)                                                                                           // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s %s=%v::\n", c.Type, c.Properties[0].Key, c.Properties[0].Value)) // should match
	})

	t.Run("writes a command with a type and multiple properties", func(t *testing.T) {
		is := is.New(t)
		c := Command{Type: StartGroupCmd, Properties: CommandProperties{{Key: "name", Value: "test"}, {Key: "age", Value: 0}}}

		var b bytes.Buffer
		err := IssueCommand(&b, &c)

		is.NoErr(err)                                                                                                                                             // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s %s=%v,%s=%v::\n", c.Type, c.Properties[0].Key, c.Properties[0].Value, c.Properties[1].Key, c.Properties[1].Value)) // should match
	})

	t.Run("writes a command with a type multiple properties and a message", func(t *testing.T) {
		is := is.New(t)
		c := Command{Type: StartGroupCmd, Properties: CommandProperties{{Key: "name", Value: "test"}, {Key: "age", Value: 0}}, Message: "message"}

		var b bytes.Buffer
		err := IssueCommand(&b, &c)

		is.NoErr(err)                                                                                                                                                          // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s %s=%v,%s=%v::%s\n", c.Type, c.Properties[0].Key, c.Properties[0].Value, c.Properties[1].Key, c.Properties[1].Value, c.Message)) // should match
	})
}

func TestIssueFileCommand(t *testing.T) {
	t.Run("errors if path is invalid", func(t *testing.T) {
		is := is.New(t)

		err := IssueFileCommand("/x/test", "key", "value")

		is.True(err != nil) // should error
	})

	t.Run("writes single line value", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		key := "key"
		value := "value"

		err := IssueFileCommand(p, key, value)

		data, _ := os.ReadFile(p)

		is.NoErr(err)                                              // should not error
		is.Equal(string(data), fmt.Sprintf("%s=%s\n", key, value)) // should match
	})

	t.Run("writes multi line value", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		key := "key"
		value := `my
    value`

		err := IssueFileCommand(p, key, value)

		data, _ := os.ReadFile(p)

		is.NoErr(err)                                  // should not error
		is.True(strings.Contains(string(data), key))   // should contain key
		is.True(strings.Contains(string(data), value)) // should contain value
	})
}
