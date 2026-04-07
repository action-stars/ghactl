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
	tests := []struct {
		name       string
		cmdType    string
		properties CommandProperties
		message    string
		wantErr    bool
	}{
		{
			name:    "errors_if_there_is_no_command_type",
			cmdType: "",
			wantErr: true,
		},
		{
			name:    "creates_a_new_command_with_just_the_type",
			cmdType: string(StartGroupCmd),
		},
		{
			name:       "creates_a_new_command_with_all_inputs",
			cmdType:    string(StartGroupCmd),
			properties: CommandProperties{{Key: "name", Value: "test"}},
			message:    "message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			c, err := NewCommand(CommandType(tt.cmdType), tt.properties, tt.message)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)                             // should not error
			is.Equal(c.Type, CommandType(tt.cmdType)) // should match type
			is.Equal(c.Properties, tt.properties)     // should match properties
			is.Equal(c.Message, tt.message)           // should match message
		})
	}
}

func TestIssueCommand(t *testing.T) {
	tests := []struct {
		name string
		cmd  Command
		want string
	}{
		{
			name: "writes_a_command_with_only_a_type",
			cmd:  Command{Type: StartGroupCmd},
			want: fmt.Sprintf("::%s::\n", StartGroupCmd),
		},
		{
			name: "writes_a_command_with_a_type_and_message",
			cmd:  Command{Type: StartGroupCmd, Message: "message"},
			want: fmt.Sprintf("::%s::message\n", StartGroupCmd),
		},
		{
			name: "writes_a_command_with_a_type_and_a_property",
			cmd:  Command{Type: StartGroupCmd, Properties: CommandProperties{{Key: "name", Value: "test"}}},
			want: fmt.Sprintf("::%s name=test::\n", StartGroupCmd),
		},
		{
			name: "writes_a_command_with_a_type_and_multiple_properties",
			cmd:  Command{Type: StartGroupCmd, Properties: CommandProperties{{Key: "name", Value: "test"}, {Key: "age", Value: 0}}},
			want: fmt.Sprintf("::%s name=test,age=0::\n", StartGroupCmd),
		},
		{
			name: "writes_a_command_with_a_type_multiple_properties_and_a_message",
			cmd:  Command{Type: StartGroupCmd, Properties: CommandProperties{{Key: "name", Value: "test"}, {Key: "age", Value: 0}}, Message: "message"},
			want: fmt.Sprintf("::%s name=test,age=0::message\n", StartGroupCmd),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := IssueCommand(&b, &tt.cmd)

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should match expected
		})
	}
}

func TestIssueFileCommand(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "errors_if_path_is_invalid",
			path:    "/x/test",
			key:     "key",
			value:   "value",
			wantErr: true,
		},
		{
			name:  "writes_single_line_value_as_heredoc",
			key:   "key",
			value: "value",
		},
		{
			name:  "writes_multi_line_value_as_heredoc",
			key:   "key",
			value: "my\n    value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			p := tt.path
			if p == "" {
				p = filepath.Join(t.TempDir(), "test")
			}

			err := IssueFileCommand(p, tt.key, tt.value)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			data, _ := os.ReadFile(p)

			is.NoErr(err)                                     // should not error
			is.True(strings.Contains(string(data), tt.key))   // should contain key
			is.True(strings.Contains(string(data), tt.value)) // should contain value
			is.True(strings.Contains(string(data), "<<"))     // should use heredoc format
		})
	}
}

func TestIssueCommand_Escaping(t *testing.T) {
	tests := []struct {
		name         string
		cmd          Command
		wantContains []string
		wantMissing  []string
	}{
		{
			name:         "escapes_special_characters_in_message",
			cmd:          Command{Type: DebugCmd, Message: "hello%world\nnewline\rcarriage"},
			wantContains: []string{"%25", "%0A", "%0D"},
			wantMissing:  []string{"%world\n"},
		},
		{
			name:         "escapes_special_characters_in_property_values",
			cmd:          Command{Type: DebugCmd, Properties: CommandProperties{{Key: "file", Value: "path:with,special%chars\nnewline"}}},
			wantContains: []string{"%25", "%0A", "%3A", "%2C"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := IssueCommand(&b, &tt.cmd)

			is.NoErr(err)
			for _, s := range tt.wantContains {
				is.True(strings.Contains(b.String(), s)) // should contain escaped sequence
			}
			for _, s := range tt.wantMissing {
				is.True(!strings.Contains(b.String(), s)) // should not contain raw chars
			}
		})
	}
}
