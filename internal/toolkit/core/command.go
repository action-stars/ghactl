package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/action-stars/ghactl/internal/util"
)

// CommandType represents a GitHub Actions command type.
type CommandType string

const (
	StartGroupCmd    CommandType = "group"
	EndGroupCmd      CommandType = "endgroup"
	DebugCmd         CommandType = "debug"
	ErrorCmd         CommandType = "error"
	WarningCmd       CommandType = "warning"
	InfoCmd          CommandType = "info"
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

// CommandProperties represents a map of properties for a GitHub Actions command.
type CommandProperties []CommandProperty

// Command represents a GitHub Actions command.
type Command struct {
	Type       CommandType
	Properties CommandProperties
	Message    string
}

// String prints the Command as a string.
func (c *Command) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("::%s", c.Type))
	if len(c.Properties) > 0 {
		sb.WriteString(" ")
	}

	sep := ""
	for _, v := range c.Properties {
		sb.WriteString(fmt.Sprintf("%s%s=%v", sep, v.Key, v.Value))
		sep = ","
	}
	sb.WriteString("::")

	if len(c.Message) > 0 {
		sb.WriteString(c.Message)
	}

	return sb.String()
}

// NewCommand creates a new GitHub Actions command.
func NewCommand(t CommandType, p CommandProperties, m string) (*Command, error) {
	if len(t) == 0 {
		return nil, fmt.Errorf("command type is required")
	}

	c := Command{
		Type:       t,
		Properties: p,
		Message:    m,
	}

	return &c, nil
}

// IssueCommand issues a GitHub Actions command.
func IssueCommand(w io.Writer, c *Command) error {
	_, err := fmt.Fprintln(w, c.String())
	if err != nil {
		return err
	}
	return nil
}

// IssueFileCommand writes a key-value pair to a file.
// If the value contains a newline character, it will be written as a heredoc using a random delimiter.
func IssueFileCommand(p, key, value string) error {
	var bytes []byte

	if strings.Contains(value, "\n") {
		del := util.GenerateRandomString(10)
		bytes = []byte(fmt.Sprintf("%s<<%s\n%s\n%[2]s\n", key, del, value))
	} else {
		bytes = []byte(fmt.Sprintf("%s=%s\n", key, value))
	}

	return util.WriteFile(p, bytes)
}
