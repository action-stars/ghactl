package core

import (
	"fmt"
	"io"
)

// AnnotationProperties represents the properties for a GitHub Actions annotation command.
type AnnotationProperties struct {
	Title     string
	File      string
	Column    int
	EndColumn int
	Line      int
	EndLine   int
}

// GetCommandProperties returns the annotation properties as CommandProperties
// suitable for use with a GitHub Actions annotation command.
func (a *AnnotationProperties) GetCommandProperties() CommandProperties {
	cp := CommandProperties{}

	if a.Title != "" {
		cp = append(cp, CommandProperty{Key: "title", Value: a.Title})
	}

	if a.File != "" {
		cp = append(cp, CommandProperty{Key: "file", Value: a.File})
	}

	if a.Column > 0 {
		cp = append(cp, CommandProperty{Key: "col", Value: a.Column})
	}

	if a.EndColumn > 0 {
		cp = append(cp, CommandProperty{Key: "endColumn", Value: a.EndColumn})
	}

	if a.Line > 0 {
		cp = append(cp, CommandProperty{Key: "line", Value: a.Line})
	}

	if a.EndLine > 0 {
		cp = append(cp, CommandProperty{Key: "endLine", Value: a.EndLine})
	}

	return cp
}

// Debug sends a debug message to the workflow log writer.
func Debug(w io.Writer, message string) error {
	c, err := NewCommand(DebugCmd, nil, message)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// Error sends an error message to the workflow log writer.
func Error(w io.Writer, message string, properties AnnotationProperties) error {
	c, err := NewCommand(ErrorCmd, properties.GetCommandProperties(), message)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// Warning sends an warning message to the workflow log writer.
func Warning(w io.Writer, message string, properties AnnotationProperties) error {
	c, err := NewCommand(WarningCmd, properties.GetCommandProperties(), message)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// Notice sends an notice message to the workflow log writer.
func Notice(w io.Writer, message string, properties AnnotationProperties) error {
	c, err := NewCommand(NoticeCmd, properties.GetCommandProperties(), message)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// Info writes a message directly to the workflow log writer.
// Unlike other log functions, info does not use a workflow command.
func Info(w io.Writer, message string) error {
	_, err := fmt.Fprintln(w, message)
	return err
}

// StartGroup starts a new group of messages in the workflow log writer.
func StartGroup(w io.Writer, name string) error {
	c, err := NewCommand(StartGroupCmd, nil, name)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// EndGroup ends the current group of messages in the workflow log writer.
func EndGroup(w io.Writer) error {
	c, err := NewCommand(EndGroupCmd, nil, "")
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}
