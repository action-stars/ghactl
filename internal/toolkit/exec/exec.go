package exec

import (
	"io"
	"os"
	"os/exec"
	"time"
)

// ExecOptions are options for the Exec function.
type ExecOptions struct {
	// Dir is the working directory of the command.
	// If empty, the current process's working directory is used.
	Dir string
	// Env is the environment to use for the command.
	// If nil, the current process's environment is used.
	Env []string
	// Silent is a flag to suppress output from the command.
	Silent bool
	// Stdout is the writer to write the output of the command to.
	// If nil, the output is written to os.Stdout.
	Stdout io.Writer
	// Stderr is the writer to write the error output of the command to.
	// If nil, the error output is written to os.Stderr.
	Stderr io.Writer
	// Stdin is the reader to read the input of the command from.
	Stdin io.Reader
	// WaitDelay is the time to wait for the I/O pipes to be closed.
	WaitDelay time.Duration
}

// Exec executes a command and returns the exit code.
// If cmd is
func Exec(cmdLine string, args []string, opts ExecOptions) (int, error) {
	cmd := exec.Command(cmdLine, args...)

	if len(opts.Dir) > 0 {
		cmd.Dir = opts.Dir
	}

	if len(opts.Env) > 0 {
		cmd.Env = opts.Env
	}

	if opts.Silent {
		cmd.Stdout = nil
		cmd.Stderr = nil
	} else {
		if opts.Stdout != nil {
			cmd.Stdout = opts.Stdout
		} else {
			cmd.Stdout = os.Stdout
		}
		if opts.Stderr != nil {
			cmd.Stderr = opts.Stderr
		} else {
			cmd.Stderr = os.Stderr
		}
	}

	if opts.Stdin != nil {
		cmd.Stdin = opts.Stdin
	}

	if opts.WaitDelay > 0 {
		cmd.WaitDelay = opts.WaitDelay
	}

	err := cmd.Run()
	if err != nil {
		cmd.ProcessState.ExitCode()
		return cmd.ProcessState.ExitCode(), err
	}

	return 0, nil
}
