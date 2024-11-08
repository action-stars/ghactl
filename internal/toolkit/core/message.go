package core

import "io"

// AnnotationProperties represents the properties for a GitHub Actions annotation command.
type AnnotationProperties struct {
	Title     string
	File      string
	Column    int
	EndColumn int
	Line      int
	EndLine   int
}

// GetCommandProperties returns the properties for a GitHub Actions annotation command.
func (a *AnnotationProperties) GetCommandProperties() CommandProperties {
	cp := CommandProperties{}

	if len(a.Title) > 0 {
		cp = append(cp, CommandProperty{Key: "title", Value: a.Title})
	}

	if len(a.File) > 0 {
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

// Info sends an info message to the workflow log writer.
func Info(w io.Writer, message string, properties AnnotationProperties) error {
	c, err := NewCommand(InfoCmd, properties.GetCommandProperties(), message)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
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
