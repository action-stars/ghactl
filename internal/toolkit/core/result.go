package core

import "io"

// SetFailed sets the action status to failed by logging an error message.
// The caller is responsible for exiting the process with an appropriate exit code.
func SetFailed(w io.Writer, message string) error {
	return Error(w, message, AnnotationProperties{})
}
