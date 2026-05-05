package exec

import "runtime"

func echoCommand() (string, []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", "echo", "hello"}
	}
	return "echo", []string{"hello"}
}

func falseCommand() (string, []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", "exit", "1"}
	}
	return "false", nil
}
