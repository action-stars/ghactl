package exec

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestExec(t *testing.T) {
	echoCmd, echoArgs := echoCommand()
	falseCmd, falseArgs := falseCommand()

	tests := []struct {
		name         string
		cmd          string
		args         []string
		opts         Options
		wantExitCode int
		wantOutput   bool
		wantErr      bool
	}{
		{
			name:         "executes_a_command_successfully",
			cmd:          echoCmd,
			args:         echoArgs,
			wantExitCode: 0,
			wantOutput:   true,
		},
		{
			name:         "returns_error_for_command_not_found",
			cmd:          "non-existent-command-xyz",
			opts:         Options{Silent: true},
			wantExitCode: -1,
			wantErr:      true,
		},
		{
			name:         "returns_non-zero_exit_code",
			cmd:          falseCmd,
			args:         falseArgs,
			opts:         Options{Silent: true},
			wantExitCode: 1,
			wantErr:      true,
		},
		{
			name:         "silent_mode_suppresses_output",
			cmd:          echoCmd,
			args:         echoArgs,
			opts:         Options{Silent: true},
			wantExitCode: 0,
		},
		{
			name:         "captures_stdout_with_custom_writer",
			cmd:          echoCmd,
			args:         echoArgs,
			wantExitCode: 0,
			wantOutput:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			var stdout bytes.Buffer

			opts := tt.opts
			if !opts.Silent {
				opts.Stdout = &stdout
				opts.Stderr = &bytes.Buffer{}
			}

			exitCode, err := Exec(t.Context(), tt.cmd, tt.args, opts)

			if tt.wantErr {
				is.True(err != nil)                 // should error
				is.Equal(exitCode, tt.wantExitCode) // exit code should match
				return
			}

			is.NoErr(err)                       // should not error
			is.Equal(exitCode, tt.wantExitCode) // exit code should match

			if tt.wantOutput {
				is.True(stdout.Len() > 0) // should have output
			}
		})
	}
}
