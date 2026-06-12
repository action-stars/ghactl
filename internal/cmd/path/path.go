package path

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"

	"github.com/action-stars/ghactl/internal/toolkit/core"
)

// Cmd provides the action logic for path subcommands.
type Cmd struct{}

// New returns the fully-wired "path" CLI command tree.
func New() *cli.Command {
	c := &Cmd{}

	return &cli.Command{
		Name:  "path",
		Usage: "Manage GitHub Actions PATH entries.",
		Commands: []*cli.Command{
			c.addCommand(),
		},
	}
}

// Add appends a path entry to the GitHub Actions PATH file and prepends to current PATH.
func (c *Cmd) Add(value string) error {
	return core.AddPath(value)
}

func (c *Cmd) addCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a path entry.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Path entry to add.",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			value := cmd.String("path")

			slog.Debug("Adding path entry.", slog.String("path", value))

			if err := c.Add(value); err != nil {
				return cli.Exit(err, 1)
			}

			slog.Debug("Path entry added.", slog.String("path", value))
			return nil
		},
	}
}
