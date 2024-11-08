package core

import (
	"fmt"
	"io"

	"github.com/action-stars/ghactl/internal/util"
)

// AddMatcher sends an add matcher command to the workflow writer.
// The problem matcher is defined by the path to the problem matcher file.
func AddMatcher(w io.Writer, p string) error {
	exists, err := util.FileExists(p)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("problem matcher file %s does not exist", p)
	}

	c, err := NewCommand(AddMatcherCmd, nil, p)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// RemoveMatcher sends a remove matcher command to the workflow writer.
// The owner is defined in the problem matcher file that was added.
func RemoveMatcher(w io.Writer, owner string) error {
	c, err := NewCommand(RemoveMatcherCmd, CommandProperties{{Key: "owner", Value: owner}}, "")
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}
