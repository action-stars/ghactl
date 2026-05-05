package core

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"
)

// CommandType represents a GitHub Actions command type.
type CommandType string

const (
	StartGroupCmd    CommandType = "group"
	EndGroupCmd      CommandType = "endgroup"
	DebugCmd         CommandType = "debug"
	ErrorCmd         CommandType = "error"
	WarningCmd       CommandType = "warning"
	NoticeCmd        CommandType = "notice"
	MaskCmd          CommandType = "add-mask"
	AddMatcherCmd    CommandType = "add-matcher"
	RemoveMatcherCmd CommandType = "remove-matcher"
)

// CommandProperty represents a key-value pair for a GitHub Actions command property.
type CommandProperty struct {
	Key   string
	Value any
}

// CommandProperties represents a slice of properties for a GitHub Actions command.
type CommandProperties []CommandProperty

// Command represents a GitHub Actions workflow command with a type, properties, and message.
type Command struct {
	Type       CommandType
	Properties CommandProperties
	Message    string
}

// escapeData escapes special characters in a command message.
func escapeData(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}

// escapeProperty escapes special characters in a command property value.
func escapeProperty(s string) string {
	s = escapeData(s)
	s = strings.ReplaceAll(s, ":", "%3A")
	s = strings.ReplaceAll(s, ",", "%2C")
	return s
}

// String prints the Command as a string.
func (c *Command) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "::%s", c.Type)
	if len(c.Properties) > 0 {
		sb.WriteString(" ")
	}

	sep := ""
	for _, v := range c.Properties {
		fmt.Fprintf(&sb, "%s%s=%s", sep, v.Key, escapeProperty(fmt.Sprintf("%v", v.Value)))
		sep = ","
	}
	sb.WriteString("::")
	sb.WriteString(escapeData(c.Message))

	return sb.String()
}

// NewCommand creates a new GitHub Actions command.
// It returns an error if the command type is empty.
func NewCommand(t CommandType, p CommandProperties, m string) (*Command, error) {
	if t == "" {
		return nil, fmt.Errorf("command type is required")
	}

	c := Command{
		Type:       t,
		Properties: p,
		Message:    m,
	}

	return &c, nil
}

// IssueCommand issues a GitHub Actions command by writing it to the provided writer.
func IssueCommand(w io.Writer, c *Command) error {
	_, err := fmt.Fprintln(w, c.String())
	return err
}

// generateDelimiter generates a unique delimiter for file commands.
func generateDelimiter() string {
	b := make([]byte, 16)
	// rand.Read from crypto/rand always returns len(b) bytes and a nil error on
	// supported platforms; a panic here indicates a broken runtime.
	_, _ = rand.Read(b)
	return fmt.Sprintf("ghadelimiter_%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// IssueFileCommand writes a key-value pair to a GitHub Actions environment file.
// It always uses heredoc format with a unique delimiter to prevent injection.
func IssueFileCommand(p, key, value string) error {
	del := generateDelimiter()

	if strings.Contains(key, del) {
		return fmt.Errorf("unexpected input: name should not contain the delimiter %q", del)
	}
	if strings.Contains(value, del) {
		return fmt.Errorf("unexpected input: value should not contain the delimiter %q", del)
	}

	bytes := fmt.Appendf(nil, "%s<<%s\n%s\n%[2]s\n", key, del, value)

	return writeFile(p, bytes)
}

// writeFile appends bytes to a file, creating it if it doesn't exist.
func writeFile(name string, value []byte) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(value)
	return err
}
